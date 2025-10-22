package patterns

import (
	"context"
	"sync"
	"time"
)

// RateLimiter 令牌桶限流器
type RateLimiter struct {
	rate       int       // 每秒产生的令牌数
	capacity   int       // 桶容量
	tokens     int       // 当前令牌数
	lastRefill time.Time // 上次填充时间
	mu         sync.Mutex
}

// NewRateLimiter 创建限流器
// rate: 每秒产生的令牌数
// capacity: 桶的最大容量
func NewRateLimiter(rate, capacity int) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求通过
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// Wait 等待直到可以获取令牌
func (rl *RateLimiter) Wait() {
	for {
		if rl.Allow() {
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// WaitWithContext 使用Context等待令牌
func (rl *RateLimiter) WaitWithContext(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if rl.Allow() {
				return nil
			}
		}
	}
}

// refill 填充令牌
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// 计算应该产生的令牌数
	tokensToAdd := int(elapsed.Seconds() * float64(rl.rate))

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.capacity {
			rl.tokens = rl.capacity
		}
		rl.lastRefill = now
	}
}

// LeakyBucket 漏桶限流器
type LeakyBucket struct {
	capacity int           // 桶容量
	rate     time.Duration // 漏水速率
	water    int           // 当前水量
	lastLeak time.Time     // 上次漏水时间
	mu       sync.Mutex
}

// NewLeakyBucket 创建漏桶限流器
func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		rate:     rate,
		water:    0,
		lastLeak: time.Now(),
	}
}

// Allow 检查是否允许请求
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.leak()

	if lb.water < lb.capacity {
		lb.water++
		return true
	}
	return false
}

// leak 执行漏水
func (lb *LeakyBucket) leak() {
	now := time.Now()
	elapsed := now.Sub(lb.lastLeak)

	// 计算应该漏出的水量
	leaked := int(elapsed / lb.rate)

	if leaked > 0 {
		lb.water -= leaked
		if lb.water < 0 {
			lb.water = 0
		}
		lb.lastLeak = now
	}
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	limit    int           // 窗口内最大请求数
	window   time.Duration // 窗口大小
	requests []time.Time   // 请求时间戳
	mu       sync.Mutex
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:    limit,
		window:   window,
		requests: make([]time.Time, 0),
	}
}

// Allow 检查是否允许请求
func (sw *SlidingWindowLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	// 移除过期的请求
	validRequests := make([]time.Time, 0)
	for _, reqTime := range sw.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	sw.requests = validRequests

	// 检查是否超过限制
	if len(sw.requests) < sw.limit {
		sw.requests = append(sw.requests, now)
		return true
	}
	return false
}

// Reset 重置限流器
func (sw *SlidingWindowLimiter) Reset() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.requests = make([]time.Time, 0)
}
