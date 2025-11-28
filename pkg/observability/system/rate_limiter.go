package system

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// RateLimiter 限流器监控
type RateLimiter struct {
	meter           metric.Meter
	enabled         bool
	limit           int64 // 每秒限制
	window          time.Duration
	mu              sync.RWMutex
	currentRequests int64
	lastReset       time.Time

	// 指标
	rateLimitCounter    metric.Int64Counter
	rateLimitRejectedCounter metric.Int64Counter
	rateLimitRemainingGauge  metric.Int64ObservableGauge
}

// RateLimiterConfig 限流器配置
type RateLimiterConfig struct {
	Meter   metric.Meter
	Enabled bool
	Limit   int64         // 每秒限制
	Window  time.Duration // 时间窗口
}

// NewRateLimiter 创建限流器监控
func NewRateLimiter(cfg RateLimiterConfig) (*RateLimiter, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	window := cfg.Window
	if window == 0 {
		window = 1 * time.Second
	}

	limiter := &RateLimiter{
		meter:        cfg.Meter,
		enabled:      cfg.Enabled,
		limit:        cfg.Limit,
		window:       window,
		lastReset:    time.Now(),
	}

	// 初始化指标
	if err := limiter.initMetrics(); err != nil {
		return nil, fmt.Errorf("failed to init metrics: %w", err)
	}

	return limiter, nil
}

// initMetrics 初始化指标
func (rl *RateLimiter) initMetrics() error {
	var err error

	rl.rateLimitCounter, err = rl.meter.Int64Counter(
		"rate_limit.requests",
		metric.WithDescription("Total number of rate-limited requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	rl.rateLimitRejectedCounter, err = rl.meter.Int64Counter(
		"rate_limit.rejected",
		metric.WithDescription("Number of rejected requests due to rate limit"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	rl.rateLimitRemainingGauge, err = rl.meter.Int64ObservableGauge(
		"rate_limit.remaining",
		metric.WithDescription("Remaining requests in current window"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	// 注册可观察指标回调
	_, err = rl.meter.RegisterCallback(rl.collectMetrics, rl.rateLimitRemainingGauge)
	if err != nil {
		return err
	}

	return nil
}

// collectMetrics 收集指标（可观察指标回调）
func (rl *RateLimiter) collectMetrics(ctx context.Context, obs metric.Observer) error {
	rl.mu.RLock()
	remaining := rl.limit - rl.currentRequests
	if remaining < 0 {
		remaining = 0
	}
	rl.mu.RUnlock()

	obs.ObserveInt64(rl.rateLimitRemainingGauge, remaining)
	return nil
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ctx context.Context) bool {
	if !rl.enabled {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 检查是否需要重置窗口
	now := time.Now()
	if now.Sub(rl.lastReset) >= rl.window {
		rl.currentRequests = 0
		rl.lastReset = now
	}

	// 检查是否超过限制
	if rl.currentRequests >= rl.limit {
		rl.rateLimitRejectedCounter.Add(ctx, 1)
		return false
	}

	rl.currentRequests++
	rl.rateLimitCounter.Add(ctx, 1)
	return true
}

// GetRemaining 获取剩余请求数
func (rl *RateLimiter) GetRemaining() int64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	remaining := rl.limit - rl.currentRequests
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetLimit 获取限制
func (rl *RateLimiter) GetLimit() int64 {
	return rl.limit
}

// UpdateLimit 更新限制
func (rl *RateLimiter) UpdateLimit(limit int64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limit = limit
}
