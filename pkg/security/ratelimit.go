package security

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrRateLimitExceeded 速率限制超出
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// RateLimiter 速率限制器
type RateLimiter struct {
	limit     int
	window    time.Duration
	requests  map[string][]time.Time
	mu        sync.RWMutex
	cleanup   *time.Ticker
	stopCleanup chan struct{}
}

// RateLimiterConfig 速率限制器配置
type RateLimiterConfig struct {
	Limit  int           // 限制数量
	Window time.Duration // 时间窗口
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		limit:     config.Limit,
		window:    config.Window,
		requests:  make(map[string][]time.Time),
		stopCleanup: make(chan struct{}),
	}

	// 启动清理协程
	rl.cleanup = time.NewTicker(1 * time.Minute)
	go rl.cleanupExpired()

	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// 获取或创建请求记录
	requests, exists := rl.requests[key]
	if !exists {
		requests = make([]time.Time, 0)
	}

	// 清理过期请求
	validRequests := make([]time.Time, 0)
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// 检查是否超过限制
	if len(validRequests) >= rl.limit {
		return false, ErrRateLimitExceeded
	}

	// 添加新请求
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests

	return true, nil
}

// Check 检查是否允许请求（不记录）
func (rl *RateLimiter) Check(ctx context.Context, key string) (bool, error) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	requests, exists := rl.requests[key]
	if !exists {
		return true, nil
	}

	// 统计有效请求
	count := 0
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			count++
		}
	}

	return count < rl.limit, nil
}

// Reset 重置指定 key 的请求记录
func (rl *RateLimiter) Reset(ctx context.Context, key string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.requests, key)
	return nil
}

// GetRemaining 获取剩余请求次数
func (rl *RateLimiter) GetRemaining(ctx context.Context, key string) (int, error) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	requests, exists := rl.requests[key]
	if !exists {
		return rl.limit, nil
	}

	// 统计有效请求
	count := 0
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			count++
		}
	}

	remaining := rl.limit - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// cleanupExpired 清理过期请求记录
func (rl *RateLimiter) cleanupExpired() {
	for {
		select {
		case <-rl.cleanup.C:
			rl.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-rl.window * 2) // 清理两倍窗口时间前的记录

			for key, requests := range rl.requests {
				validRequests := make([]time.Time, 0)
				for _, reqTime := range requests {
					if reqTime.After(cutoff) {
						validRequests = append(validRequests, reqTime)
					}
				}

				if len(validRequests) == 0 {
					delete(rl.requests, key)
				} else {
					rl.requests[key] = validRequests
				}
			}
			rl.mu.Unlock()

		case <-rl.stopCleanup:
			rl.cleanup.Stop()
			return
		}
	}
}

// Shutdown 关闭速率限制器
func (rl *RateLimiter) Shutdown(ctx context.Context) error {
	close(rl.stopCleanup)
	return nil
}

// IPRateLimiter IP 速率限制器
type IPRateLimiter struct {
	*RateLimiter
}

// NewIPRateLimiter 创建 IP 速率限制器
func NewIPRateLimiter(config RateLimiterConfig) *IPRateLimiter {
	return &IPRateLimiter{
		RateLimiter: NewRateLimiter(config),
	}
}

// AllowIP 检查 IP 是否允许请求
func (rl *IPRateLimiter) AllowIP(ctx context.Context, ip string) (bool, error) {
	return rl.Allow(ctx, "ip:"+ip)
}

// UserRateLimiter 用户速率限制器
type UserRateLimiter struct {
	*RateLimiter
}

// NewUserRateLimiter 创建用户速率限制器
func NewUserRateLimiter(config RateLimiterConfig) *UserRateLimiter {
	return &UserRateLimiter{
		RateLimiter: NewRateLimiter(config),
	}
}

// AllowUser 检查用户是否允许请求
func (rl *UserRateLimiter) AllowUser(ctx context.Context, userID string) (bool, error) {
	return rl.Allow(ctx, "user:"+userID)
}

// EndpointRateLimiter 端点速率限制器
type EndpointRateLimiter struct {
	*RateLimiter
}

// NewEndpointRateLimiter 创建端点速率限制器
func NewEndpointRateLimiter(config RateLimiterConfig) *EndpointRateLimiter {
	return &EndpointRateLimiter{
		RateLimiter: NewRateLimiter(config),
	}
}

// AllowEndpoint 检查端点是否允许请求
func (rl *EndpointRateLimiter) AllowEndpoint(ctx context.Context, endpoint string) (bool, error) {
	return rl.Allow(ctx, "endpoint:"+endpoint)
}
