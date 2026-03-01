// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// ABAC (Attribute-Based Access Control) 基于属性的访问控制模型，
// 通过评估主体（Subject）、资源（Resource）、操作（Action）和环境（Environment）
// 的属性来决定访问是否被允许。
//
// 与 RBAC 的区别：
//   - RBAC 基于预定义的角色进行访问控制
//   - ABAC 基于动态属性进行细粒度访问控制
//
// 核心组件：
//   - Engine: ABAC 引擎，管理策略并执行评估
//   - Policy: 策略，定义访问控制规则
//   - Rule: 规则，决定策略是否匹配
//   - Condition: 条件，评估复杂的属性表达式
//
// 快速开始：
//
//	// 创建 ABAC 引擎
//	engine := abac.NewEngine()
//
//	// 添加策略
//	engine.AddPolicy(abac.Policy{
//	    ID:       "policy-1",
//	    Name:     "Allow managers to edit documents during business hours",
//	    Effect:   abac.Allow,
//	    Priority: 100,
//	    Rules: abac.And(
//	        abac.SubjectHasRole("manager"),
//	        abac.ResourceTypeIs("document"),
//	        abac.ActionIs("edit"),
//	    ),
//	    Enabled: true,
//	})
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
//	})
//
//	if result.Allowed {
//	    // 允许访问
//	} else {
//	    // 拒绝访问
//	}
package abac

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Result 表示评估结果
type Result struct {
	Allowed     bool     `json:"allowed"`
	Decision    Effect   `json:"decision"`
	MatchedPolicy *Policy `json:"matched_policy,omitempty"`
	Reason      string   `json:"reason,omitempty"`
	Errors      []error  `json:"errors,omitempty"`
}

// IsAllowed 检查是否允许访问
func (r Result) IsAllowed() bool {
	return r.Allowed && r.Decision == Allow
}

// IsDenied 检查是否拒绝访问
func (r Result) IsDenied() bool {
	return !r.Allowed || r.Decision == Deny
}

// Engine 是 ABAC 引擎的核心结构
//
// 负责管理策略集合并执行访问控制评估
type Engine struct {
	policies     map[string]Policy
	defaultEffect Effect
	mu           sync.RWMutex
	// 可选：策略变更回调
	onPolicyChange func(event PolicyChangeEvent)
}

// PolicyChangeEvent 表示策略变更事件
type PolicyChangeEvent struct {
	Type     string // "add", "update", "remove"
	PolicyID string
	Policy   *Policy
}

// EngineOption 是引擎配置选项
type EngineOption func(*Engine)

// WithDefaultEffect 设置默认效果
//
// 当没有策略匹配时的默认行为
func WithDefaultEffect(effect Effect) EngineOption {
	return func(e *Engine) {
		e.defaultEffect = effect
	}
}

// WithPolicyChangeCallback 设置策略变更回调
func WithPolicyChangeCallback(callback func(event PolicyChangeEvent)) EngineOption {
	return func(e *Engine) {
		e.onPolicyChange = callback
	}
}

// NewEngine 创建新的 ABAC 引擎
//
// 示例：
//
//	engine := abac.NewEngine()
//
//	// 或者带选项
//	engine := abac.NewEngine(
//	    abac.WithDefaultEffect(abac.Deny),
//	)
func NewEngine(opts ...EngineOption) *Engine {
	engine := &Engine{
		policies:      make(map[string]Policy),
		defaultEffect: Deny, // 默认拒绝
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
//	    ID:       "policy-1",
//	    Name:     "Allow admin access",
//	    Effect:   abac.Allow,
//	    Rules:    abac.SubjectHasRole("admin"),
//	    Enabled:  true,
//	}
//
//	if err := engine.AddPolicy(policy); err != nil {
//	    log.Fatal(err)
//	}
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

	// 触发回调
	if e.onPolicyChange != nil {
		e.onPolicyChange(PolicyChangeEvent{
			Type:     "add",
			PolicyID: policy.ID,
			Policy:   &policy,
		})
	}

	return nil
}

// UpdatePolicy 更新策略
//
// 如果策略不存在，返回错误
func (e *Engine) UpdatePolicy(policy Policy) error {
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("invalid policy: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.policies[policy.ID]; !exists {
		return fmt.Errorf("policy with ID %s not found", policy.ID)
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

	return nil
}

// RemovePolicy 删除策略
//
// 示例：
//
//	if err := engine.RemovePolicy("policy-1"); err != nil {
//	    log.Fatal(err)
//	}
func (e *Engine) RemovePolicy(policyID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	policy, exists := e.policies[policyID]
	if !exists {
		return fmt.Errorf("policy with ID %s not found", policyID)
	}

	delete(e.policies, policyID)

	// 触发回调
	if e.onPolicyChange != nil {
		e.onPolicyChange(PolicyChangeEvent{
			Type:     "remove",
			PolicyID: policyID,
			Policy:   &policy,
		})
	}

	return nil
}

// GetPolicy 获取策略
//
// 示例：
//
//	policy, err := engine.GetPolicy("policy-1")
//	if err != nil {
//	    log.Fatal(err)
//	}
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
		Allowed:    e.defaultEffect == Allow,
		Decision:   e.defaultEffect,
		Reason:     "No matching policy found",
		Errors:     errors,
	}
}

// EvaluateWithReason 评估访问请求并返回详细原因
//
// 返回所有策略的评估结果，用于调试
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
		Allowed:    e.defaultEffect == Allow,
		Decision:   e.defaultEffect,
		Reason:     "No matching policy found",
		Errors:     errors,
	}, policyResults
}

// BatchEvaluate 批量评估多个请求
//
// 返回每个请求的结果
func (e *Engine) BatchEvaluate(ctx context.Context, requests []Request) []Result {
	results := make([]Result, len(requests))
	for i, req := range requests {
		results[i] = e.Evaluate(ctx, req)
	}
	return results
}

// GetStats 获取引擎统计信息
//
// 返回策略数量等统计信息
func (e *Engine) GetStats() EngineStats {
	e.mu.RLock()
	defer e.mu.RUnlock()

	stats := EngineStats{
		TotalPolicies: len(e.policies),
		EnabledPolicies: 0,
	}

	for _, policy := range e.policies {
		if policy.Enabled {
			stats.EnabledPolicies++
		}
	}

	return stats
}

// EngineStats 引擎统计信息
type EngineStats struct {
	TotalPolicies   int `json:"total_policies"`
	EnabledPolicies int `json:"enabled_policies"`
}

// Validate 验证引擎配置
//
// 检查策略的有效性
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
func (e *Engine) Clone() *Engine {
	e.mu.RLock()
	defer e.mu.RUnlock()

	clone := NewEngine(WithDefaultEffect(e.defaultEffect))

	for _, policy := range e.policies {
		clone.policies[policy.ID] = policy
	}

	return clone
}
