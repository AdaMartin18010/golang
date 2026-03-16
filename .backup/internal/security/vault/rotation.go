// Package vault provides HashiCorp Vault integration for secure secret management.
//
// rotation.go 提供密钥轮换功能，包括：
// 1. 自动密钥轮换
// 2. 轮换策略管理
// 3. 轮换历史记录
// 4. 手动轮换触发
//
// 设计原则：
// 1. 支持基于时间和基于事件的轮换
// 2. 可配置的轮换策略
// 3. 轮换失败处理和重试
// 4. 轮换历史记录和审计
//
// 使用场景：
// - 定期轮换数据库凭据
// - 定期轮换 API 密钥
// - 紧急手动轮换
// - 合规性要求的密钥轮换
package vault

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
)

// RotationPolicy 定义密钥轮换策略
type RotationPolicy struct {
	// ID 是策略唯一标识
	ID string `json:"id"`
	// Name 是策略名称
	Name string `json:"name"`
	// Path 是 Vault 中密钥的路径
	Path string `json:"path"`
	// Interval 是轮换间隔时间
	Interval time.Duration `json:"interval"`
	// GracePeriod 是宽限期（旧密钥仍然有效的时间）
	GracePeriod time.Duration `json:"grace_period"`
	// Enabled 表示策略是否启用
	Enabled bool `json:"enabled"`
	// AutoRotate 表示是否自动轮换
	AutoRotate bool `json:"auto_rotate"`
	// RetainVersions 是保留的版本数量（0 表示不限制）
	RetainVersions int `json:"retain_versions"`
	// PreRotateHook 是轮换前执行的钩子函数（可选）
	PreRotateHook func(ctx context.Context, path string) error `json:"-"`
	// PostRotateHook 是轮换后执行的钩子函数（可选）
	PostRotateHook func(ctx context.Context, path string, newVersion int) error `json:"-"`
}

// RotationHistory 记录密钥轮换历史
type RotationHistory struct {
	// ID 是记录唯一标识
	ID string `json:"id"`
	// PolicyID 是关联的策略 ID
	PolicyID string `json:"policy_id"`
	// Path 是密钥路径
	Path string `json:"path"`
	// OldVersion 是旧版本号
	OldVersion int `json:"old_version"`
	// NewVersion 是新版本号
	NewVersion int `json:"new_version"`
	// RotatedAt 是轮换时间
	RotatedAt time.Time `json:"rotated_at"`
	// RotatedBy 是执行者（系统或用户）
	RotatedBy string `json:"rotated_by"`
	// Status 是轮换状态（success, failed, pending）
	Status string `json:"status"`
	// ErrorMessage 是错误信息（如果失败）
	ErrorMessage string `json:"error_message,omitempty"`
	// Metadata 是额外元数据
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RotationService 是密钥轮换服务接口。
//
// 提供密钥轮换的完整生命周期管理，包括：
// - 轮换策略管理
// - 自动轮换调度
// - 轮换历史记录
// - 手动轮换触发
type RotationService interface {
	// RegisterPolicy 注册密钥轮换策略
	RegisterPolicy(policy *RotationPolicy) error

	// UnregisterPolicy 注销密钥轮换策略
	UnregisterPolicy(policyID string) error

	// GetPolicy 获取轮换策略
	GetPolicy(policyID string) (*RotationPolicy, error)

	// ListPolicies 列出所有轮换策略
	ListPolicies() []*RotationPolicy

	// Rotate 手动触发密钥轮换
	Rotate(ctx context.Context, path string) error

	// RotateWithPolicy 使用指定策略轮换密钥
	RotateWithPolicy(ctx context.Context, policyID string) error

	// StartAutoRotation 启动自动轮换调度
	StartAutoRotation() error

	// StopAutoRotation 停止自动轮换调度
	StopAutoRotation() error

	// IsAutoRotationRunning 检查自动轮换是否正在运行
	IsAutoRotationRunning() bool

	// GetRotationHistory 获取轮换历史记录
	GetRotationHistory(path string) []RotationHistory

	// GetAllRotationHistory 获取所有轮换历史记录
	GetAllRotationHistory() []RotationHistory

	// CleanupOldVersions 清理旧版本密钥
	CleanupOldVersions(ctx context.Context, path string, keepVersions int) error
}

// rotationService 是 RotationService 的实现
type rotationService struct {
	client   *api.Client
	kv       KVManager
	policies map[string]*RotationPolicy
	history  []RotationHistory
	mu       sync.RWMutex
	running  bool
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// newRotationService 创建新的轮换服务
func newRotationService(client *api.Client, kv KVManager) RotationService {
	return &rotationService{
		client:   client,
		kv:       kv,
		policies: make(map[string]*RotationPolicy),
		history:  make([]RotationHistory, 0),
		stopChan: make(chan struct{}),
	}
}

// RegisterPolicy 注册轮换策略
func (r *rotationService) RegisterPolicy(policy *RotationPolicy) error {
	if policy == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	if policy.ID == "" {
		policy.ID = generatePolicyID()
	}

	if policy.Path == "" {
		return fmt.Errorf("policy path cannot be empty")
	}

	if policy.Interval <= 0 {
		policy.Interval = 24 * time.Hour // 默认每天轮换
	}

	if policy.RetainVersions < 0 {
		policy.RetainVersions = 0
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.policies[policy.ID] = policy
	return nil
}

// UnregisterPolicy 注销轮换策略
func (r *rotationService) UnregisterPolicy(policyID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.policies[policyID]; !exists {
		return fmt.Errorf("policy not found: %s", policyID)
	}

	delete(r.policies, policyID)
	return nil
}

// GetPolicy 获取轮换策略
func (r *rotationService) GetPolicy(policyID string) (*RotationPolicy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	policy, exists := r.policies[policyID]
	if !exists {
		return nil, fmt.Errorf("policy not found: %s", policyID)
	}

	// 返回副本，防止外部修改
	policyCopy := *policy
	return &policyCopy, nil
}

// ListPolicies 列出所有轮换策略
func (r *rotationService) ListPolicies() []*RotationPolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	policies := make([]*RotationPolicy, 0, len(r.policies))
	for _, policy := range r.policies {
		policyCopy := *policy
		policies = append(policies, &policyCopy)
	}

	return policies
}

// Rotate 手动触发密钥轮换
func (r *rotationService) Rotate(ctx context.Context, path string) error {
	return r.rotateSecret(ctx, path, "manual", nil)
}

// RotateWithPolicy 使用指定策略轮换密钥
func (r *rotationService) RotateWithPolicy(ctx context.Context, policyID string) error {
	policy, err := r.GetPolicy(policyID)
	if err != nil {
		return err
	}

	if !policy.Enabled {
		return fmt.Errorf("policy is disabled: %s", policyID)
	}

	return r.rotateSecret(ctx, policy.Path, policyID, policy)
}

// rotateSecret 执行密钥轮换
func (r *rotationService) rotateSecret(ctx context.Context, path, triggeredBy string, policy *RotationPolicy) error {
	// 获取当前密钥信息
	oldSecret, err := r.kv.Get(ctx, path)
	if err != nil {
		r.addHistory(path, 0, 0, triggeredBy, "failed", err.Error())
		return fmt.Errorf("failed to get secret for rotation: %w", err)
	}

	oldVersion := oldSecret.Version
	if oldVersion == 0 {
		oldVersion = 1
	}

	// 执行轮换前钩子
	if policy != nil && policy.PreRotateHook != nil {
		if err := policy.PreRotateHook(ctx, path); err != nil {
			r.addHistory(path, oldVersion, 0, triggeredBy, "failed", fmt.Sprintf("pre-rotate hook failed: %v", err))
			return fmt.Errorf("pre-rotate hook failed: %w", err)
		}
	}

	// 创建新版本（通过更新相同数据）
	if err := r.kv.Put(ctx, path, oldSecret.Data); err != nil {
		r.addHistory(path, oldVersion, 0, triggeredBy, "failed", err.Error())
		return fmt.Errorf("failed to rotate secret: %w", err)
	}

	// 获取新版本号
	newSecret, err := r.kv.Get(ctx, path)
	if err != nil {
		r.addHistory(path, oldVersion, 0, triggeredBy, "failed", err.Error())
		return fmt.Errorf("failed to get rotated secret: %w", err)
	}

	newVersion := newSecret.Version

	// 执行轮换后钩子
	if policy != nil && policy.PostRotateHook != nil {
		if err := policy.PostRotateHook(ctx, path, newVersion); err != nil {
			// 记录警告但不影响轮换结果
			fmt.Printf("post-rotate hook failed: %v\n", err)
		}
	}

	// 清理旧版本
	if policy != nil && policy.RetainVersions > 0 {
		if err := r.CleanupOldVersions(ctx, path, policy.RetainVersions); err != nil {
			// 记录警告但不影响轮换结果
			fmt.Printf("cleanup old versions failed: %v\n", err)
		}
	}

	// 记录历史
	r.addHistory(path, oldVersion, newVersion, triggeredBy, "success", "")

	return nil
}

// addHistory 添加轮换历史记录
func (r *rotationService) addHistory(path string, oldVersion, newVersion int, rotatedBy, status, errorMsg string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	history := RotationHistory{
		ID:           generateHistoryID(),
		Path:         path,
		OldVersion:   oldVersion,
		NewVersion:   newVersion,
		RotatedAt:    time.Now(),
		RotatedBy:    rotatedBy,
		Status:       status,
		ErrorMessage: errorMsg,
	}

	r.history = append(r.history, history)

	// 限制历史记录数量
	if len(r.history) > 1000 {
		r.history = r.history[len(r.history)-1000:]
	}
}

// StartAutoRotation 启动自动轮换调度
func (r *rotationService) StartAutoRotation() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return fmt.Errorf("auto rotation is already running")
	}

	r.running = true
	r.stopChan = make(chan struct{})

	// 启动调度器
	r.wg.Add(1)
	go r.scheduler()

	return nil
}

// StopAutoRotation 停止自动轮换调度
func (r *rotationService) StopAutoRotation() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return nil
	}

	r.running = false
	close(r.stopChan)

	r.mu.Unlock()
	r.wg.Wait()
	r.mu.Lock()

	return nil
}

// IsAutoRotationRunning 检查自动轮换是否正在运行
func (r *rotationService) IsAutoRotationRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

// scheduler 自动轮换调度器
func (r *rotationService) scheduler() {
	defer r.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.checkAndRotate()
		}
	}
}

// checkAndRotate 检查并执行轮换
func (r *rotationService) checkAndRotate() {
	r.mu.RLock()
	policies := make([]*RotationPolicy, 0, len(r.policies))
	for _, policy := range r.policies {
		if policy.Enabled && policy.AutoRotate {
			policyCopy := *policy
			policies = append(policies, &policyCopy)
		}
	}
	r.mu.RUnlock()

	now := time.Now()
	for _, policy := range policies {
		// 检查是否需要轮换
		lastRotation := r.getLastRotationTime(policy.Path)
		if now.Sub(lastRotation) >= policy.Interval {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			if err := r.RotateWithPolicy(ctx, policy.ID); err != nil {
				fmt.Printf("auto rotation failed for %s: %v\n", policy.Path, err)
			}
			cancel()
		}
	}
}

// getLastRotationTime 获取最后一次轮换时间
func (r *rotationService) getLastRotationTime(path string) time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := len(r.history) - 1; i >= 0; i-- {
		if r.history[i].Path == path && r.history[i].Status == "success" {
			return r.history[i].RotatedAt
		}
	}

	return time.Time{}
}

// GetRotationHistory 获取指定路径的轮换历史
func (r *rotationService) GetRotationHistory(path string) []RotationHistory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	history := make([]RotationHistory, 0)
	for _, h := range r.history {
		if h.Path == path {
			history = append(history, h)
		}
	}

	return history
}

// GetAllRotationHistory 获取所有轮换历史
func (r *rotationService) GetAllRotationHistory() []RotationHistory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	history := make([]RotationHistory, len(r.history))
	copy(history, r.history)

	return history
}

// CleanupOldVersions 清理旧版本密钥
func (r *rotationService) CleanupOldVersions(ctx context.Context, path string, keepVersions int) error {
	if !r.kv.IsV2() {
		return fmt.Errorf("cleanup old versions is only supported in KV v2")
	}

	versions, err := r.kv.GetVersions(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to get versions for cleanup: %w", err)
	}

	// 过滤出有效的非销毁版本
	validVersions := make([]int, 0)
	for _, v := range versions {
		if !v.Destroyed && v.DeletionTime == nil {
			validVersions = append(validVersions, v.Version)
		}
	}

	// 如果版本数不超过保留数量，不需要清理
	if len(validVersions) <= keepVersions {
		return nil
	}

	// 删除旧版本（保留最近的 keepVersions 个）
	versionsToDelete := validVersions[:len(validVersions)-keepVersions]
	for _, version := range versionsToDelete {
		if err := r.kv.DestroyVersion(ctx, path, version); err != nil {
			return fmt.Errorf("failed to destroy version %d: %w", version, err)
		}
	}

	return nil
}

// generatePolicyID 生成策略 ID
func generatePolicyID() string {
	return fmt.Sprintf("policy-%d", time.Now().UnixNano())
}

// generateHistoryID 生成历史记录 ID
func generateHistoryID() string {
	return fmt.Sprintf("hist-%d", time.Now().UnixNano())
}
