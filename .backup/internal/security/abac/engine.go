// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件实现了 ABAC 引擎，是 ABAC 系统的核心组件。
//
// 引擎功能：
//   - 策略管理：添加、更新、删除、查询策略
//   - 访问评估：评估访问请求，决定是否允许
//   - 批量评估：支持批量请求评估
//   - 统计信息：提供引擎运行统计
//   - 缓存集成：支持策略缓存提升性能
//
// 使用示例：
//
//	// 创建 ABAC 引擎
//	engine := abac.NewEngine(
//	    abac.WithDefaultEffect(abac.Deny),
//	    abac.WithCache(abac.DefaultCacheConfig()),
//	)
//
//	// 添加策略
//	policy := abac.Policy{
//	    ID:       "policy-001",
//	    Name:     "Allow managers to edit documents during business hours",
//	    Effect:   abac.Allow,
//	    Priority: 100,
//	    Rules: abac.And(
//	        abac.SubjectHasRole("manager"),
//	        abac.ResourceTypeIs("document"),
//	        abac.ActionIs("edit"),
//	        abac.EnvironmentAttributeEquals("connection", "internal"),
//	    ),
//	    Enabled: true,
//	}
//
//	if err := engine.AddPolicy(policy); err != nil {
//	    log.Fatal(err)
//	}
//
//	// 评估访问请求
//	result := engine.Evaluate(ctx, abac.Request{
//	    Subject: abac.Subject{
//	        ID:    "user1",
//	        Roles: []string{"manager"},
//	    },
//	    Resource: abac.Resource{
//	        Type:  "document",
//	        Owner: "user1",
//	    },
//	    Action: abac.Action{Name: "edit"},
//	    Environment: abac.Environment{
//	        Connection: "internal",
//	    },
//	})
//
//	if result.Allowed {
//	    // 允许访问
//	} else {
//	    // 拒绝访问，查看原因
//	    log.Printf("Access denied: %s", result.Reason)
//	}
package abac

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Result 表示评估结果
type Result struct {
	Allowed       bool     `json:"allowed"`
	Decision      Effect   `json:"decision"`
	MatchedPolicy *Policy `json:"matched_policy,omitempty"`
	Reason        string   `json:"reason,omitempty"`
	Errors        []error  `json:"errors,omitempty"`
	Cached        bool     `json:"cached"` // 是否来自缓存
}

// IsAllowed 检查是否允许访问
//
// 返回：
//   - bool: 如果允许访问返回 true
func (r Result) IsAllowed() bool {
	return r.Allowed && r.Decision == Allow
}

// IsDenied 检查是否拒绝访问
//
// 返回：
//   - bool: 如果拒绝访问返回 true
func (r Result) IsDenied() bool {
	return !r.Allowed || r.Decision == Deny
}

// Engine 是 ABAC 引擎的核心结构
//
// 负责管理策略集合并执行访问控制评估
type Engine struct {
	policies      map[string]Policy
	defaultEffect Effect
	mu            sync.RWMutex
	cache         Cache
	cacheEnabled  bool
	onPolicyChange func(event PolicyChangeEvent)
	stats         EngineStats
}

// PolicyChangeEvent 表示策略变更事件
type PolicyChangeEvent struct {
	Type     string  // "add", "update", "remove"
	PolicyID string
	Policy   *Policy
}

// EngineOption 是引擎配置选项
type EngineOption func(*Engine)

// WithDefaultEffect 设置默认效果
//
// 当没有策略匹配时的默认行为
//
// 参数：
//   - effect: 默认效果（Allow 或 Deny）
//
// 返回：
//   - EngineOption: 配置选项
func WithDefaultEffect(effect Effect) EngineOption {
	return func(e *Engine) {
		e.defaultEffect = effect
	}
}

// WithCache 启用缓存
//
// 参数：
//   - config: 缓存配置
//
// 返回：
//   - EngineOption: 配置选项
func WithCache(config CacheConfig) EngineOption {
	return func(e *Engine) {
		e.cache = NewCache(config)
		e.cacheEnabled = true
	}
}

// WithPolicyChangeCallback 设置策略变更回调
//
// 当策略发生变更时触发回调函数
//
// 参数：
//   - callback: 回调函数
//
// 返回：
//   - EngineOption: 配置选项
func WithPolicyChangeCallback(callback func(event PolicyChangeEvent)) EngineOption {
	return func(e *Engine) {
		e.onPolicyChange = callback
	}
}

// NewEngine 创建新的 ABAC 引擎
//
// 示例：
//
//	// 基础用法
//	engine := abac.NewEngine()
//
//	// 带选项的用法
//	engine := abac.NewEngine(
//	    abac.WithDefaultEffect(abac.Deny),
//	    abac.WithCache(abac.DefaultCacheConfig()),
//	)
//
// 参数：
//   - opts: 引擎配置选项
//
// 返回：
//   - *Engine: ABAC 引擎实例
func NewEngine(opts ...EngineOption) *Engine {
	engine := &Engine{
		policies:      make(map[string]Policy),
		defaultEffect: Deny, // 默认拒绝
		cacheEnabled:  false,
		stats:         EngineStats{StartedAt: time.Now().Unix()},
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine
}

// AddPolicy 添加策略
//
// 如果策略ID已存在，返回错误
//
// 示例：
//
//	policy := abac.Policy{
//	    ID:       "policy-001",
//	    Name:     "Allow admin access",
//	    Effect:   abac.Allow,
//	    Rules:    abac.SubjectHasRole("admin"),
//	    Enabled:  true,
//	}
//
//	if err := engine.AddPolicy(policy); err != nil {
//	    log.Fatal(err)
//	}
//
// 参数：
//   - policy: 要添加的策略
//
// 返回：
//   - error: 如果添加失败返回错误
func (e *Engine) AddPolicy(policy Policy) error {
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("invalid policy: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.policies[policy.ID]; exists {
		return fmt.Errorf("policy with ID %s already exists", policy.ID)
	}

	e.policies[policy.ID] = policy
	e.stats.TotalPolicies++
	if policy.Enabled {
		e.stats.EnabledPolicies++
	}

	// 触发回调
	if e.onPolicyChange != nil {
		e.onPolicyChange(PolicyChangeEvent{
			Type:     "add",
			PolicyID: policy.ID,
			Policy:   &policy,
		})
	}

	// 使缓存失效
	if e.cacheEnabled && e.cache != nil {
		e.cache.Clear()
	}

	return nil
}

// UpdatePolicy 更新策略
//
// 如果策略不存在，返回错误
//
// 参数：
//   - policy: 要更新的策略
//
// 返回：
//   - error: 如果更新失败返回错误
func (e *Engine) UpdatePolicy(policy Policy) error {
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("invalid policy: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	oldPolicy, exists := e.policies[policy.ID]
	if !exists {
		return fmt.Errorf("policy with ID %s not found", policy.ID)
	}

	// 更新统计
	if oldPolicy.Enabled != policy.Enabled {
		if policy.Enabled {
			e.stats.EnabledPolicies++
		} else {
			e.stats.EnabledPolicies--
		}
	}

	e.policies[policy.ID] = policy

	// 触发回调
	if e.onPolicyChange != nil {
		e.onPolicyChange(PolicyChangeEvent{
			Type:     "update",
			PolicyID: policy.ID,
			Policy:   &policy,
		})
	}

	// 使缓存失效
	if e.cacheEnabled && e.cache != nil {
		e.cache.Clear()
	}

	return nil
}

// RemovePolicy 删除策略
//
// 示例：
//
//	if err := engine.RemovePolicy("policy-001"); err != nil {
//	    log.Fatal(err)
//	}
//
// 参数：
//   - policyID: 策略ID
//
// 返回：
//   - error: 如果删除失败返回错误
func (e *Engine) RemovePolicy(policyID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	policy, exists := e.policies[policyID]
	if !exists {
		return fmt.Errorf("policy with ID %s not found", policyID)
	}

	if policy.Enabled {
		e.stats.EnabledPolicies--
	}
	e.stats.TotalPolicies--
	delete(e.policies, policyID)

	// 触发回调
	if e.onPolicyChange != nil {
		e.onPolicyChange(PolicyChangeEvent{
			Type:     "remove",
			PolicyID: policyID,
			Policy:   &policy,
		})
	}

	// 使缓存失效
	if e.cacheEnabled && e.cache != nil {
		e.cache.Clear()
	}

	return nil
}

// GetPolicy 获取策略
//
// 示例：
//
//	policy, err := engine.GetPolicy("policy-001")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// 参数：
//   - policyID: 策略ID
//
// 返回：
//   - *Policy: 策略指针
//   - error: 如果策略不存在返回错误
func (e *Engine) GetPolicy(policyID string) (*Policy, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	policy, exists := e.policies[policyID]
	if !exists {
		return nil, fmt.Errorf("policy with ID %s not found", policyID)
	}

	return &policy, nil
}

// ListPolicies 列出所有策略
//
// 返回按优先级排序的策略列表
//
// 返回：
//   - []Policy: 策略列表
func (e *Engine) ListPolicies() []Policy {
	e.mu.RLock()
	defer e.mu.RUnlock()

	policies := make([]Policy, 0, len(e.policies))
	for _, policy := range e.policies {
		policies = append(policies, policy)
	}

	// 按优先级降序排序
	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Priority > policies[j].Priority
	})

	return policies
}

// ClearPolicies 清空所有策略
func (e *Engine) ClearPolicies() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.policies = make(map[string]Policy)
	e.stats.TotalPolicies = 0
	e.stats.EnabledPolicies = 0

	// 使缓存失效
	if e.cacheEnabled && e.cache != nil {
		e.cache.Clear()
	}
}

// Evaluate 评估访问请求
//
// 按优先级从高到低依次评估策略，第一个匹配的策略决定结果
//
// 参数：
//   - ctx: 上下文
//   - req: 访问请求
//
// 返回：
//   - Result: 评估结果
//
// 示例：
//
//	result := engine.Evaluate(ctx, abac.Request{
//	    Subject: abac.Subject{
//	        ID:    "user1",
//	        Roles: []string{"manager"},
//	    },
//	    Resource: abac.Resource{
//	        Type:  "document",
//	        Owner: "user1",
//	    },
//	    Action: abac.Action{Name: "edit"},
//	})
//
//	if result.IsAllowed() {
//	    // 允许访问
//	}
func (e *Engine) Evaluate(ctx context.Context, req Request) Result {
	start := time.Now()
	defer func() {
		e.mu.Lock()
		e.stats.TotalEvaluations++
		e.stats.TotalEvaluationTime += time.Since(start).Milliseconds()
		e.mu.Unlock()
	}()

	policies := e.ListPolicies()

	var errors []error

	for _, policy := range policies {
		// 跳过未启用的策略
		if !policy.Enabled {
			continue
		}

		// 评估策略
		matched, err := policy.Match(ctx, req)
		if err != nil {
			errors = append(errors, fmt.Errorf("policy %s evaluation error: %w", policy.ID, err))
			continue
		}

		if matched {
			return Result{
				Allowed:       policy.Effect == Allow,
				Decision:      policy.Effect,
				MatchedPolicy: &policy,
				Reason:        fmt.Sprintf("Matched policy: %s (%s)", policy.Name, policy.ID),
				Errors:        errors,
			}
		}
	}

	// 没有匹配的策略，返回默认效果
	return Result{
		Allowed:  e.defaultEffect == Allow,
		Decision: e.defaultEffect,
		Reason:   "No matching policy found",
		Errors:   errors,
	}
}

// EvaluateWithReason 评估访问请求并返回详细原因
//
// 返回所有策略的评估结果，用于调试
//
// 参数：
//   - ctx: 上下文
//   - req: 访问请求
//
// 返回：
//   - Result: 评估结果
//   - map[string]bool: 每个策略的匹配结果
func (e *Engine) EvaluateWithReason(ctx context.Context, req Request) (Result, map[string]bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	policyResults := make(map[string]bool)
	var errors []error

	// 按优先级排序
	policies := make([]Policy, 0, len(e.policies))
	for _, policy := range e.policies {
		policies = append(policies, policy)
	}
	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Priority > policies[j].Priority
	})

	for _, policy := range policies {
		if !policy.Enabled {
			policyResults[policy.ID] = false
			continue
		}

		matched, err := policy.Match(ctx, req)
		if err != nil {
			errors = append(errors, fmt.Errorf("policy %s: %w", policy.ID, err))
			policyResults[policy.ID] = false
			continue
		}

		policyResults[policy.ID] = matched

		if matched {
			return Result{
				Allowed:       policy.Effect == Allow,
				Decision:      policy.Effect,
				MatchedPolicy: &policy,
				Reason:        fmt.Sprintf("Matched policy: %s (%s)", policy.Name, policy.ID),
				Errors:        errors,
			}, policyResults
		}
	}

	return Result{
		Allowed:  e.defaultEffect == Allow,
		Decision: e.defaultEffect,
		Reason:   "No matching policy found",
		Errors:   errors,
	}, policyResults
}

// BatchEvaluate 批量评估多个请求
//
// 返回每个请求的结果
//
// 参数：
//   - ctx: 上下文
//   - requests: 访问请求列表
//
// 返回：
//   - []Result: 每个请求的结果
func (e *Engine) BatchEvaluate(ctx context.Context, requests []Request) []Result {
	results := make([]Result, len(requests))
	for i, req := range requests {
		results[i] = e.Evaluate(ctx, req)
	}
	return results
}

// EngineStats 引擎统计信息
type EngineStats struct {
	TotalPolicies       int    `json:"total_policies"`
	EnabledPolicies     int    `json:"enabled_policies"`
	TotalEvaluations    uint64 `json:"total_evaluations"`
	TotalEvaluationTime int64  `json:"total_evaluation_time_ms"` // 总评估时间（毫秒）
	StartedAt           int64  `json:"started_at"`               // 启动时间戳
}

// GetStats 获取引擎统计信息
//
// 返回：
//   - EngineStats: 引擎统计信息
func (e *Engine) GetStats() EngineStats {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.stats
}

// Validate 验证引擎配置
//
// 检查策略的有效性
//
// 返回：
//   - error: 如果配置无效返回错误
func (e *Engine) Validate() error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for id, policy := range e.policies {
		if err := policy.Validate(); err != nil {
			return fmt.Errorf("policy %s validation error: %w", id, err)
		}
	}

	return nil
}

// Clone 创建引擎的深拷贝
//
// 用于测试和备份
//
// 返回：
//   - *Engine: 引擎的深拷贝
func (e *Engine) Clone() *Engine {
	e.mu.RLock()
	defer e.mu.RUnlock()

	clone := NewEngine(WithDefaultEffect(e.defaultEffect))

	for _, policy := range e.policies {
		clone.policies[policy.ID] = policy
	}

	return clone
}

// EnableCache 启用缓存
//
// 参数：
//   - config: 缓存配置
func (e *Engine) EnableCache(config CacheConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.cache = NewCache(config)
	e.cacheEnabled = true
}

// DisableCache 禁用缓存
func (e *Engine) DisableCache() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.cacheEnabled = false
}

// GetCacheStats 获取缓存统计
//
// 返回：
//   - CacheStats: 缓存统计信息
//   - bool: 是否启用了缓存
func (e *Engine) GetCacheStats() (CacheStats, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.cacheEnabled || e.cache == nil {
		return CacheStats{}, false
	}

	return e.cache.Stats(), true
}

// IsAllowed 是 Evaluate 的便捷方法，只返回布尔值
//
// 参数：
//   - ctx: 上下文
//   - req: 访问请求
//
// 返回：
//   - bool: 如果允许访问返回 true
func (e *Engine) IsAllowed(ctx context.Context, req Request) bool {
	result := e.Evaluate(ctx, req)
	return result.IsAllowed()
}
