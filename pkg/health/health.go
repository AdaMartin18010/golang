package health

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrCheckNotFound 检查未找到
	ErrCheckNotFound = errors.New("check not found")
)

// Status 健康状态
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// Check 健康检查接口
type Check interface {
	// Name 返回检查名称
	Name() string
	// Check 执行健康检查
	Check(ctx context.Context) Result
}

// Result 检查结果
type Result struct {
	Status    Status                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// HealthChecker 健康检查器
type HealthChecker struct {
	checks map[string]Check
	mu     sync.RWMutex
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]Check),
	}
}

// Register 注册健康检查
func (hc *HealthChecker) Register(check Check) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[check.Name()] = check
}

// Unregister 注销健康检查
func (hc *HealthChecker) Unregister(name string) error {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if _, exists := hc.checks[name]; !exists {
		return ErrCheckNotFound
	}

	delete(hc.checks, name)
	return nil
}

// Check 执行所有健康检查
func (hc *HealthChecker) Check(ctx context.Context) map[string]Result {
	hc.mu.RLock()
	checks := make([]Check, 0, len(hc.checks))
	for _, check := range hc.checks {
		checks = append(checks, check)
	}
	hc.mu.RUnlock()

	results := make(map[string]Result)
	for _, check := range checks {
		results[check.Name()] = check.Check(ctx)
	}

	return results
}

// OverallStatus 获取整体健康状态
func (hc *HealthChecker) OverallStatus(ctx context.Context) Status {
	results := hc.Check(ctx)

	if len(results) == 0 {
		return StatusHealthy
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, result := range results {
		switch result.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}
	if hasDegraded {
		return StatusDegraded
	}

	return StatusHealthy
}

// SimpleCheck 简单健康检查
type SimpleCheck struct {
	name    string
	checkFn func(ctx context.Context) error
}

// NewSimpleCheck 创建简单健康检查
func NewSimpleCheck(name string, checkFn func(ctx context.Context) error) *SimpleCheck {
	return &SimpleCheck{
		name:    name,
		checkFn: checkFn,
	}
}

// Name 返回检查名称
func (sc *SimpleCheck) Name() string {
	return sc.name
}

// Check 执行健康检查
func (sc *SimpleCheck) Check(ctx context.Context) Result {
	err := sc.checkFn(ctx)
	if err != nil {
		return Result{
			Status:    StatusUnhealthy,
			Message:   err.Error(),
			Timestamp: time.Now(),
		}
	}

	return Result{
		Status:    StatusHealthy,
		Timestamp: time.Now(),
	}
}

// TimeoutCheck 带超时的健康检查
type TimeoutCheck struct {
	name    string
	timeout time.Duration
	check   Check
}

// NewTimeoutCheck 创建带超时的健康检查
func NewTimeoutCheck(name string, timeout time.Duration, check Check) *TimeoutCheck {
	return &TimeoutCheck{
		name:    name,
		timeout: timeout,
		check:   check,
	}
}

// Name 返回检查名称
func (tc *TimeoutCheck) Name() string {
	return tc.name
}

// Check 执行健康检查
func (tc *TimeoutCheck) Check(ctx context.Context) Result {
	ctx, cancel := context.WithTimeout(ctx, tc.timeout)
	defer cancel()

	done := make(chan Result, 1)
	go func() {
		done <- tc.check.Check(ctx)
	}()

	select {
	case result := <-done:
		return result
	case <-ctx.Done():
		return Result{
			Status:    StatusUnhealthy,
			Message:   "health check timeout",
			Timestamp: time.Now(),
		}
	}
}

// PeriodicCheck 定期健康检查
type PeriodicCheck struct {
	name     string
	interval time.Duration
	check    Check
	mu       sync.RWMutex
	cached   Result
	lastRun  time.Time
}

// NewPeriodicCheck 创建定期健康检查
func NewPeriodicCheck(name string, interval time.Duration, check Check) *PeriodicCheck {
	pc := &PeriodicCheck{
		name:     name,
		interval: interval,
		check:    check,
	}

	// 立即执行一次
	pc.cached = check.Check(context.Background())
	pc.lastRun = time.Now()

	return pc
}

// Name 返回检查名称
func (pc *PeriodicCheck) Name() string {
	return pc.name
}

// Check 执行健康检查
func (pc *PeriodicCheck) Check(ctx context.Context) Result {
	pc.mu.RLock()
	lastRun := pc.lastRun
	cached := pc.cached
	pc.mu.RUnlock()

	// 如果缓存未过期，返回缓存结果
	if time.Since(lastRun) < pc.interval {
		return cached
	}

	// 执行新的检查
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// 双重检查
	if time.Since(pc.lastRun) < pc.interval {
		return pc.cached
	}

	pc.cached = pc.check.Check(ctx)
	pc.lastRun = time.Now()

	return pc.cached
}

// AggregateCheck 聚合健康检查
type AggregateCheck struct {
	name   string
	checks []Check
}

// NewAggregateCheck 创建聚合健康检查
func NewAggregateCheck(name string, checks ...Check) *AggregateCheck {
	return &AggregateCheck{
		name:   name,
		checks: checks,
	}
}

// Name 返回检查名称
func (ac *AggregateCheck) Name() string {
	return ac.name
}

// Check 执行健康检查
func (ac *AggregateCheck) Check(ctx context.Context) Result {
	results := make([]Result, 0, len(ac.checks))
	for _, check := range ac.checks {
		results = append(results, check.Check(ctx))
	}

	hasUnhealthy := false
	hasDegraded := false
	messages := make([]string, 0)

	for _, result := range results {
		switch result.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
			if result.Message != "" {
				messages = append(messages, result.Message)
			}
		case StatusDegraded:
			hasDegraded = true
		}
	}

	status := StatusHealthy
	if hasUnhealthy {
		status = StatusUnhealthy
	} else if hasDegraded {
		status = StatusDegraded
	}

	message := ""
	if len(messages) > 0 {
		message = messages[0] // 简化：只返回第一个错误消息
	}

	return Result{
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"checks": len(results),
		},
	}
}
