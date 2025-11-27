package control

import (
	"context"
	"sync"
	"time"
)

// Controller 控制器接口
// 提供细粒度的控制和动态配置能力
type Controller interface {
	// Enable 启用功能
	Enable(name string) error

	// Disable 禁用功能
	Disable(name string) error

	// IsEnabled 检查功能是否启用
	IsEnabled(name string) bool

	// SetConfig 设置配置
	SetConfig(name string, config interface{}) error

	// GetConfig 获取配置
	GetConfig(name string) (interface{}, error)

	// Watch 监听配置变化
	Watch(name string, callback func(interface{})) error

	// Unwatch 取消监听
	Unwatch(name string) error
}

// FeatureController 功能控制器
// 提供功能开关和配置管理
type FeatureController struct {
	mu       sync.RWMutex
	features map[string]*Feature
	watchers map[string][]func(interface{})
}

// Feature 功能定义
type Feature struct {
	Name        string
	Enabled     bool
	Config      interface{}
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewFeatureController 创建功能控制器
func NewFeatureController() Controller {
	return &FeatureController{
		features: make(map[string]*Feature),
		watchers: make(map[string][]func(interface{})),
	}
}

// Register 注册功能
func (fc *FeatureController) Register(name, description string, enabled bool, config interface{}) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.features[name] = &Feature{
		Name:        name,
		Enabled:     enabled,
		Config:      config,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (fc *FeatureController) Enable(name string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	feature, exists := fc.features[name]
	if !exists {
		return ErrFeatureNotFound
	}

	feature.Enabled = true
	feature.UpdatedAt = time.Now()

	// 通知监听者
	fc.notifyWatchers(name, feature.Config)

	return nil
}

func (fc *FeatureController) Disable(name string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	feature, exists := fc.features[name]
	if !exists {
		return ErrFeatureNotFound
	}

	feature.Enabled = false
	feature.UpdatedAt = time.Now()

	// 通知监听者
	fc.notifyWatchers(name, feature.Config)

	return nil
}

func (fc *FeatureController) IsEnabled(name string) bool {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	feature, exists := fc.features[name]
	if !exists {
		return false
	}

	return feature.Enabled
}

func (fc *FeatureController) SetConfig(name string, config interface{}) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	feature, exists := fc.features[name]
	if !exists {
		return ErrFeatureNotFound
	}

	feature.Config = config
	feature.UpdatedAt = time.Now()

	// 通知监听者
	fc.notifyWatchers(name, config)

	return nil
}

func (fc *FeatureController) GetConfig(name string) (interface{}, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	feature, exists := fc.features[name]
	if !exists {
		return nil, ErrFeatureNotFound
	}

	return feature.Config, nil
}

func (fc *FeatureController) Watch(name string, callback func(interface{})) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if _, exists := fc.features[name]; !exists {
		return ErrFeatureNotFound
	}

	fc.watchers[name] = append(fc.watchers[name], callback)
	return nil
}

func (fc *FeatureController) Unwatch(name string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	delete(fc.watchers, name)
	return nil
}

// notifyWatchers 通知监听者
func (fc *FeatureController) notifyWatchers(name string, config interface{}) {
	watchers := fc.watchers[name]
	for _, callback := range watchers {
		go callback(config) // 异步通知
	}
}

// RateController 速率控制器
// 提供细粒度的速率控制
type RateController struct {
	mu      sync.RWMutex
	limits  map[string]*RateLimit
	enabled map[string]bool
}

// RateLimit 速率限制
type RateLimit struct {
	Name      string
	MaxRate   float64
	Window    time.Duration
	Current   float64
	LastReset time.Time
	Enabled   bool
}

// NewRateController 创建速率控制器
func NewRateController() *RateController {
	return &RateController{
		limits:  make(map[string]*RateLimit),
		enabled: make(map[string]bool),
	}
}

// SetRateLimit 设置速率限制
func (rc *RateController) SetRateLimit(name string, maxRate float64, window time.Duration) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.limits[name] = &RateLimit{
		Name:      name,
		MaxRate:   maxRate,
		Window:    window,
		Current:   0,
		LastReset: time.Now(),
		Enabled:   true,
	}
	rc.enabled[name] = true
}

// Allow 检查是否允许操作
func (rc *RateController) Allow(name string) bool {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	limit, exists := rc.limits[name]
	if !exists || !limit.Enabled {
		return true
	}

	now := time.Now()
	if now.Sub(limit.LastReset) >= limit.Window {
		limit.Current = 0
		limit.LastReset = now
	}

	if limit.Current >= limit.MaxRate {
		return false
	}

	limit.Current++
	return true
}

// Enable 启用速率限制
func (rc *RateController) Enable(name string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if limit, exists := rc.limits[name]; exists {
		limit.Enabled = true
		rc.enabled[name] = true
	}
}

// Disable 禁用速率限制
func (rc *RateController) Disable(name string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if limit, exists := rc.limits[name]; exists {
		limit.Enabled = false
		rc.enabled[name] = false
	}
}

// CircuitController 熔断器控制器
// 提供细粒度的熔断控制
type CircuitController struct {
	mu        sync.RWMutex
	circuits  map[string]*Circuit
}

// Circuit 熔断器
type Circuit struct {
	Name          string
	State         CircuitState
	FailureCount  int64
	SuccessCount  int64
	FailureThreshold int64
	SuccessThreshold int64
	Timeout       time.Duration
	LastFailure   time.Time
	LastSuccess   time.Time
}

// CircuitState 熔断器状态
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"   // 关闭（正常）
	CircuitStateOpen     CircuitState = "open"     // 打开（熔断）
	CircuitStateHalfOpen CircuitState = "half-open" // 半开（尝试恢复）
)

// NewCircuitController 创建熔断器控制器
func NewCircuitController() *CircuitController {
	return &CircuitController{
		circuits: make(map[string]*Circuit),
	}
}

// RegisterCircuit 注册熔断器
func (cc *CircuitController) RegisterCircuit(name string, failureThreshold, successThreshold int64, timeout time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.circuits[name] = &Circuit{
		Name:            name,
		State:           CircuitStateClosed,
		FailureThreshold: failureThreshold,
		SuccessThreshold: successThreshold,
		Timeout:         timeout,
	}
}

// RecordSuccess 记录成功
func (cc *CircuitController) RecordSuccess(name string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	circuit, exists := cc.circuits[name]
	if !exists {
		return
	}

	circuit.SuccessCount++
	circuit.LastSuccess = time.Now()

	if circuit.State == CircuitStateHalfOpen {
		if circuit.SuccessCount >= circuit.SuccessThreshold {
			circuit.State = CircuitStateClosed
			circuit.FailureCount = 0
			circuit.SuccessCount = 0
		}
	}
}

// RecordFailure 记录失败
func (cc *CircuitController) RecordFailure(name string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	circuit, exists := cc.circuits[name]
	if !exists {
		return
	}

	circuit.FailureCount++
	circuit.LastFailure = time.Now()

	if circuit.State == CircuitStateClosed {
		if circuit.FailureCount >= circuit.FailureThreshold {
			circuit.State = CircuitStateOpen
		}
	} else if circuit.State == CircuitStateHalfOpen {
		circuit.State = CircuitStateOpen
	}
}

// IsOpen 检查是否熔断
func (cc *CircuitController) IsOpen(name string) bool {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	circuit, exists := cc.circuits[name]
	if !exists {
		return false
	}

	// 检查是否可以尝试恢复
	if circuit.State == CircuitStateOpen {
		if time.Since(circuit.LastFailure) >= circuit.Timeout {
			circuit.State = CircuitStateHalfOpen
			circuit.SuccessCount = 0
		}
	}

	return circuit.State == CircuitStateOpen
}

// Allow 检查是否允许操作
func (cc *CircuitController) Allow(name string) bool {
	return !cc.IsOpen(name)
}
