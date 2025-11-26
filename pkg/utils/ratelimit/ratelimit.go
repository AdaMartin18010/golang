package ratelimit

import (
	"sync"
	"time"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow() bool
	Wait()
	Reset()
}

// TokenBucket 令牌桶限流器
type TokenBucket struct {
	capacity     int64         // 桶容量
	tokens       int64         // 当前令牌数
	refillRate   int64         // 每秒补充的令牌数
	lastRefill   time.Time     // 上次补充时间
	mu           sync.Mutex    // 互斥锁
}

// NewTokenBucket 创建令牌桶限流器
func NewTokenBucket(capacity, refillRate int64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许通过
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	// 补充令牌
	tb.refill()
	
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// Wait 等待直到允许通过
func (tb *TokenBucket) Wait() {
	for !tb.Allow() {
		time.Sleep(10 * time.Millisecond)
	}
}

// Reset 重置限流器
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens = tb.capacity
	tb.lastRefill = time.Now()
}

// refill 补充令牌
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int64(elapsed.Seconds()) * tb.refillRate
	
	if tokensToAdd > 0 {
		tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
		tb.lastRefill = now
	}
}

// LeakyBucket 漏桶限流器
type LeakyBucket struct {
	capacity     int64         // 桶容量
	level        int64         // 当前水位
	leakRate     int64         // 每秒漏出速率
	lastLeak     time.Time     // 上次漏水时间
	mu           sync.Mutex    // 互斥锁
}

// NewLeakyBucket 创建漏桶限流器
func NewLeakyBucket(capacity, leakRate int64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		level:    0,
		leakRate: leakRate,
		lastLeak: time.Now(),
	}
}

// Allow 检查是否允许通过
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	// 漏水
	lb.leak()
	
	if lb.level < lb.capacity {
		lb.level++
		return true
	}
	return false
}

// Wait 等待直到允许通过
func (lb *LeakyBucket) Wait() {
	for !lb.Allow() {
		time.Sleep(10 * time.Millisecond)
	}
}

// Reset 重置限流器
func (lb *LeakyBucket) Reset() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.level = 0
	lb.lastLeak = time.Now()
}

// leak 漏水
func (lb *LeakyBucket) leak() {
	now := time.Now()
	elapsed := now.Sub(lb.lastLeak)
	leaked := int64(elapsed.Seconds()) * lb.leakRate
	
	if leaked > 0 {
		lb.level = max(0, lb.level-leaked)
		lb.lastLeak = now
	}
}

// SlidingWindow 滑动窗口限流器
type SlidingWindow struct {
	windowSize   time.Duration // 窗口大小
	maxRequests  int64         // 最大请求数
	requests     []time.Time   // 请求时间戳
	mu           sync.Mutex    // 互斥锁
}

// NewSlidingWindow 创建滑动窗口限流器
func NewSlidingWindow(windowSize time.Duration, maxRequests int64) *SlidingWindow {
	return &SlidingWindow{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		requests:    make([]time.Time, 0),
	}
}

// Allow 检查是否允许通过
func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-sw.windowSize)
	
	// 移除窗口外的请求
	validRequests := make([]time.Time, 0)
	for _, req := range sw.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	sw.requests = validRequests
	
	// 检查是否超过限制
	if int64(len(sw.requests)) < sw.maxRequests {
		sw.requests = append(sw.requests, now)
		return true
	}
	return false
}

// Wait 等待直到允许通过
func (sw *SlidingWindow) Wait() {
	for !sw.Allow() {
		time.Sleep(10 * time.Millisecond)
	}
}

// Reset 重置限流器
func (sw *SlidingWindow) Reset() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.requests = sw.requests[:0]
}

// FixedWindow 固定窗口限流器
type FixedWindow struct {
	windowSize   time.Duration // 窗口大小
	maxRequests  int64         // 最大请求数
	requests     int64         // 当前窗口请求数
	windowStart  time.Time     // 窗口开始时间
	mu           sync.Mutex    // 互斥锁
}

// NewFixedWindow 创建固定窗口限流器
func NewFixedWindow(windowSize time.Duration, maxRequests int64) *FixedWindow {
	return &FixedWindow{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		requests:    0,
		windowStart: time.Now(),
	}
}

// Allow 检查是否允许通过
func (fw *FixedWindow) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	
	now := time.Now()
	
	// 检查是否需要重置窗口
	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.requests = 0
		fw.windowStart = now
	}
	
	// 检查是否超过限制
	if fw.requests < fw.maxRequests {
		fw.requests++
		return true
	}
	return false
}

// Wait 等待直到允许通过
func (fw *FixedWindow) Wait() {
	for !fw.Allow() {
		time.Sleep(10 * time.Millisecond)
	}
}

// Reset 重置限流器
func (fw *FixedWindow) Reset() {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.requests = 0
	fw.windowStart = time.Now()
}

// 辅助函数
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

