package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/http/response"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerSecond int           // 每秒请求数
	Burst             int           // 突发请求数
	Window            time.Duration // 时间窗口
	KeyFunc           func(*http.Request) string // 限流键生成函数
	SkipPaths         []string      // 跳过限流的路径
	OnLimitExceeded   func(http.ResponseWriter, *http.Request) // 限流时的处理函数
}

// TokenBucket 令牌桶
type TokenBucket struct {
	capacity     int       // 容量
	tokens       float64   // 当前令牌数
	refillRate   float64   // 填充速率（令牌/秒）
	lastRefill   time.Time // 上次填充时间
	mu           sync.Mutex
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
	if config.KeyFunc == nil {
		config.KeyFunc = defaultKeyFunc
	}
	if config.OnLimitExceeded == nil {
		config.OnLimitExceeded = defaultOnLimitExceeded
	}

	// 创建令牌桶映射
	buckets := make(map[string]*TokenBucket)
	var mu sync.RWMutex

	refillRate := float64(config.RequestsPerSecond) / config.Window.Seconds()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否跳过限流
			if shouldSkipRateLimit(r.URL.Path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// 生成限流键
			key := config.KeyFunc(r)

			// 获取或创建令牌桶
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

			// 检查是否允许请求
			if !bucket.Allow() {
				config.OnLimitExceeded(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
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

// min 返回两个浮点数中的较小值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
