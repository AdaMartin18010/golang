package middleware

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/http/response"
)

// RedisClient 定义 Redis 客户端接口，避免直接依赖
type RedisClient interface {
	Pipeline() RedisPipeline
	ZRemRangeByScore(ctx context.Context, key, min, max string) *RedisIntCmd
	ZCard(ctx context.Context, key string) *RedisIntCmd
	ZAdd(ctx context.Context, key string, members ...interface{}) *RedisIntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *RedisBoolCmd
	Exec(ctx context.Context) ([]RedisCmder, error)
}

// RedisPipeline Redis Pipeline 接口
type RedisPipeline interface {
	ZRemRangeByScore(ctx context.Context, key, min, max string) *RedisIntCmd
	ZCard(ctx context.Context, key string) *RedisIntCmd
	ZAdd(ctx context.Context, key string, members ...interface{}) *RedisIntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *RedisBoolCmd
	Exec(ctx context.Context) ([]RedisCmder, error)
}

// RedisCmder Redis 命令接口
type RedisCmder interface {
	Err() error
}

// RedisIntCmd Redis 整数命令
type RedisIntCmd struct {
	val int64
	err error
}

func (cmd *RedisIntCmd) Val() int64 { return cmd.val }
func (cmd *RedisIntCmd) Err() error { return cmd.err }

// RedisBoolCmd Redis 布尔命令
type RedisBoolCmd struct {
	val bool
	err error
}

func (cmd *RedisBoolCmd) Val() bool  { return cmd.val }
func (cmd *RedisBoolCmd) Err() error { return cmd.err }

// RedisZ Redis 有序集合成员
type RedisZ struct {
	Score  float64
	Member interface{}
}

// RateLimitAlgorithm 限流算法类型
type RateLimitAlgorithm string

const (
	AlgorithmTokenBucket   RateLimitAlgorithm = "token_bucket"   // 令牌桶算法
	AlgorithmSlidingWindow RateLimitAlgorithm = "sliding_window" // 滑动窗口算法
	AlgorithmLeakyBucket   RateLimitAlgorithm = "leaky_bucket"   // 漏桶算法
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerSecond int                                      // 每秒请求数
	Burst             int                                      // 突发请求数
	Window            time.Duration                            // 时间窗口
	Algorithm         RateLimitAlgorithm                       // 限流算法
	KeyFunc           func(*http.Request) string               // 限流键生成函数
	SkipPaths         []string                                 // 跳过限流的路径
	OnLimitExceeded   func(http.ResponseWriter, *http.Request) // 限流时的处理函数
	// Redis 配置（用于分布式限流）
	RedisClient    RedisClient // Redis 客户端（可选，用于分布式限流）
	RedisKeyPrefix string      // Redis 键前缀
}

// TokenBucket 令牌桶
type TokenBucket struct {
	capacity   int       // 容量
	tokens     float64   // 当前令牌数
	refillRate float64   // 填充速率（令牌/秒）
	lastRefill time.Time // 上次填充时间
	mu         sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	// 填充令牌
	tb.tokens = min(float64(tb.capacity), tb.tokens+elapsed*tb.refillRate)
	tb.lastRefill = now

	// 检查是否有足够的令牌
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config RateLimitConfig) func(http.Handler) http.Handler {
	// 默认配置
	if config.RequestsPerSecond == 0 {
		config.RequestsPerSecond = 100
	}
	if config.Burst == 0 {
		config.Burst = config.RequestsPerSecond
	}
	if config.Window == 0 {
		config.Window = time.Second
	}
	if config.Algorithm == "" {
		config.Algorithm = AlgorithmTokenBucket
	}
	if config.KeyFunc == nil {
		config.KeyFunc = defaultKeyFunc
	}
	if config.OnLimitExceeded == nil {
		config.OnLimitExceeded = defaultOnLimitExceeded
	}
	if config.RedisKeyPrefix == "" {
		config.RedisKeyPrefix = "ratelimit"
	}

	// 如果使用 Redis，创建 Redis 限流器
	if config.RedisClient != nil {
		redisLimiter := NewRedisRateLimiter(
			config.RedisClient,
			config.RedisKeyPrefix,
			config.RequestsPerSecond,
			config.Window,
		)

		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 检查是否跳过限流
				if shouldSkipRateLimit(r.URL.Path, config.SkipPaths) {
					next.ServeHTTP(w, r)
					return
				}

				// 生成限流键
				key := config.KeyFunc(r)

				// 使用 Redis 限流
				allowed, err := redisLimiter.Allow(r.Context(), key)
				if err != nil {
					// Redis 错误时允许请求通过（降级策略）
					next.ServeHTTP(w, r)
					return
				}

				if !allowed {
					config.OnLimitExceeded(w, r)
					return
				}

				next.ServeHTTP(w, r)
			})
		}
	}

	// 内存限流器
	var limiters map[string]interface{}
	var mu sync.RWMutex

	switch config.Algorithm {
	case AlgorithmSlidingWindow:
		limiters = make(map[string]interface{})

		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if shouldSkipRateLimit(r.URL.Path, config.SkipPaths) {
					next.ServeHTTP(w, r)
					return
				}

				key := config.KeyFunc(r)

				mu.RLock()
				limiter, exists := limiters[key]
				mu.RUnlock()

				if !exists {
					mu.Lock()
					limiter, exists = limiters[key]
					if !exists {
						limiter = NewSlidingWindow(config.RequestsPerSecond, config.Window)
						limiters[key] = limiter
					}
					mu.Unlock()
				}

				if !limiter.(*SlidingWindow).Allow() {
					config.OnLimitExceeded(w, r)
					return
				}

				next.ServeHTTP(w, r)
			})
		}

	case AlgorithmLeakyBucket:
		limiters = make(map[string]interface{})
		leakRate := float64(config.RequestsPerSecond) / config.Window.Seconds()

		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if shouldSkipRateLimit(r.URL.Path, config.SkipPaths) {
					next.ServeHTTP(w, r)
					return
				}

				key := config.KeyFunc(r)

				mu.RLock()
				limiter, exists := limiters[key]
				mu.RUnlock()

				if !exists {
					mu.Lock()
					limiter, exists = limiters[key]
					if !exists {
						limiter = NewLeakyBucket(config.Burst, leakRate)
						limiters[key] = limiter
					}
					mu.Unlock()
				}

				if !limiter.(*LeakyBucket).Allow() {
					config.OnLimitExceeded(w, r)
					return
				}

				next.ServeHTTP(w, r)
			})
		}

	default: // AlgorithmTokenBucket
		buckets := make(map[string]*TokenBucket)
		refillRate := float64(config.RequestsPerSecond) / config.Window.Seconds()

		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if shouldSkipRateLimit(r.URL.Path, config.SkipPaths) {
					next.ServeHTTP(w, r)
					return
				}

				key := config.KeyFunc(r)

				mu.RLock()
				bucket, exists := buckets[key]
				mu.RUnlock()

				if !exists {
					mu.Lock()
					bucket, exists = buckets[key]
					if !exists {
						bucket = NewTokenBucket(config.Burst, refillRate)
						buckets[key] = bucket
					}
					mu.Unlock()
				}

				if !bucket.Allow() {
					config.OnLimitExceeded(w, r)
					return
				}

				next.ServeHTTP(w, r)
			})
		}
	}
}

// defaultKeyFunc 默认限流键生成函数（基于IP）
func defaultKeyFunc(r *http.Request) string {
	// 优先使用X-Forwarded-For
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	// 使用RemoteAddr
	return r.RemoteAddr
}

// defaultOnLimitExceeded 默认限流处理函数
func defaultOnLimitExceeded(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusTooManyRequests,
		errors.NewRateLimitError("rate limit exceeded"))
}

// shouldSkipRateLimit 检查是否应该跳过限流
func shouldSkipRateLimit(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || path == skipPath+"/" {
			return true
		}
	}
	return false
}

// SlidingWindow 滑动窗口限流器
type SlidingWindow struct {
	requests []time.Time   // 请求时间戳列表
	limit    int           // 限制数量
	window   time.Duration // 时间窗口
	mu       sync.Mutex
}

// NewSlidingWindow 创建滑动窗口限流器
func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		requests: make([]time.Time, 0, limit),
		limit:    limit,
		window:   window,
	}
}

// Allow 检查是否允许请求
func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	// 移除窗口外的请求
	validRequests := sw.requests[:0]
	for _, reqTime := range sw.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	sw.requests = validRequests

	// 检查是否超过限制
	if len(sw.requests) >= sw.limit {
		return false
	}

	// 添加当前请求
	sw.requests = append(sw.requests, now)
	return true
}

// LeakyBucket 漏桶限流器
type LeakyBucket struct {
	capacity int       // 容量
	leakRate float64   // 漏出速率（请求/秒）
	water    float64   // 当前水量
	lastLeak time.Time // 上次漏水时间
	mu       sync.Mutex
}

// NewLeakyBucket 创建漏桶限流器
func NewLeakyBucket(capacity int, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		water:    0,
		lastLeak: time.Now(),
	}
}

// Allow 检查是否允许请求
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeak).Seconds()

	// 漏水
	lb.water = max(0, lb.water-elapsed*lb.leakRate)
	lb.lastLeak = now

	// 检查是否有容量
	if lb.water < float64(lb.capacity) {
		lb.water += 1.0
		return true
	}

	return false
}

// RedisRateLimiter Redis 分布式限流器
type RedisRateLimiter struct {
	client    RedisClient
	keyPrefix string
	limit     int
	window    time.Duration
}

// NewRedisRateLimiter 创建 Redis 限流器
func NewRedisRateLimiter(client RedisClient, keyPrefix string, limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:    client,
		keyPrefix: keyPrefix,
		limit:     limit,
		window:    window,
	}
}

// Allow 检查是否允许请求（使用滑动窗口）
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	fullKey := rl.keyPrefix + ":" + key
	now := time.Now()
	windowStart := now.Add(-rl.window)

	pipe := rl.client.Pipeline()
	// 移除窗口外的记录
	pipe.ZRemRangeByScore(ctx, fullKey, "0", formatTime(windowStart))
	// 统计当前窗口内的请求数
	count := pipe.ZCard(ctx, fullKey)
	// 添加当前请求
	pipe.ZAdd(ctx, fullKey, RedisZ{
		Score:  float64(now.UnixNano()),
		Member: now.String(),
	})
	// 设置过期时间
	pipe.Expire(ctx, fullKey, rl.window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return count.Val() < int64(rl.limit), nil
}

// formatTime 格式化时间为 Redis 分数
func formatTime(t time.Time) string {
	return strconv.FormatInt(t.UnixNano(), 10)
}

// min 返回两个浮点数中的较小值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// max 返回两个浮点数中的较大值
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
