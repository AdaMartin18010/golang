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

// RedisClient 定义了 Redis 客户端接口，用于分布式限流。
//
// 功能说明：
// - 避免直接依赖具体的 Redis 客户端库
// - 支持依赖注入和测试
// - 提供 Redis 有序集合操作（用于滑动窗口算法）
//
// 接口方法：
// - Pipeline: 创建 Pipeline 用于批量操作
// - ZRemRangeByScore: 移除有序集合中指定分数范围的成员
// - ZCard: 获取有序集合的成员数量
// - ZAdd: 向有序集合添加成员
// - Expire: 设置键的过期时间
// - Exec: 执行 Pipeline 中的命令
//
// 使用场景：
// - 分布式限流（多实例共享限流状态）
// - 跨服务限流
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

// RateLimitAlgorithm 是限流算法的类型定义。
//
// 支持的算法：
// - AlgorithmTokenBucket: 令牌桶算法
//   特点：允许突发流量，平滑限流
//   适用场景：需要允许突发请求的场景
// - AlgorithmSlidingWindow: 滑动窗口算法
//   特点：精确控制时间窗口内的请求数
//   适用场景：需要精确限流的场景
// - AlgorithmLeakyBucket: 漏桶算法
//   特点：平滑输出，限制突发流量
//   适用场景：需要平滑输出的场景
type RateLimitAlgorithm string

const (
	AlgorithmTokenBucket   RateLimitAlgorithm = "token_bucket"   // 令牌桶算法
	AlgorithmSlidingWindow RateLimitAlgorithm = "sliding_window" // 滑动窗口算法
	AlgorithmLeakyBucket   RateLimitAlgorithm = "leaky_bucket"   // 漏桶算法
)

// RateLimitConfig 是限流中间件的配置结构。
//
// 功能说明：
// - 配置限流算法和参数
// - 支持内存限流和分布式限流（Redis）
// - 支持自定义限流键生成和限流处理
//
// 字段说明：
// - RequestsPerSecond: 每秒允许的请求数（默认：100）
// - Burst: 突发请求数（令牌桶和漏桶算法使用，默认：等于 RequestsPerSecond）
// - Window: 时间窗口（滑动窗口算法使用，默认：1秒）
// - Algorithm: 限流算法（默认：令牌桶）
// - KeyFunc: 限流键生成函数（默认：基于 IP 地址）
// - SkipPaths: 跳过限流的路径列表
// - OnLimitExceeded: 限流时的处理函数（默认：返回 429 错误）
// - RedisClient: Redis 客户端（可选，用于分布式限流）
// - RedisKeyPrefix: Redis 键前缀（默认：ratelimit）
//
// 使用示例：
//
//	config := middleware.RateLimitConfig{
//	    RequestsPerSecond: 100,
//	    Burst:             50,
//	    Algorithm:         middleware.AlgorithmTokenBucket,
//	    SkipPaths:         []string{"/health", "/metrics"},
//	}
//	router.Use(middleware.RateLimitMiddleware(config))
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

// TokenBucket 是令牌桶限流器的实现。
//
// 算法原理：
// - 桶中存放令牌，请求需要消耗令牌
// - 按固定速率向桶中添加令牌
// - 如果桶中有令牌，允许请求；否则拒绝请求
//
// 特点：
// - 允许突发流量（桶满时可以快速处理多个请求）
// - 平滑限流（令牌按固定速率添加）
// - 适合需要允许突发请求的场景
//
// 字段说明：
// - capacity: 桶的容量（最大令牌数）
// - tokens: 当前令牌数
// - refillRate: 填充速率（令牌/秒）
// - lastRefill: 上次填充时间
// - mu: 互斥锁（保证并发安全）
type TokenBucket struct {
	capacity   int       // 容量
	tokens     float64   // 当前令牌数
	refillRate float64   // 填充速率（令牌/秒）
	lastRefill time.Time // 上次填充时间
	mu         sync.Mutex
}

// NewTokenBucket 创建并初始化令牌桶。
//
// 参数：
// - capacity: 桶的容量（最大令牌数）
// - refillRate: 填充速率（令牌/秒）
//
// 返回：
// - *TokenBucket: 配置好的令牌桶实例
//
// 使用示例：
//
//	bucket := NewTokenBucket(100, 10.0) // 容量100，每秒填充10个令牌
func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity), // 初始时桶是满的
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求通过。
//
// 功能说明：
// - 根据时间差计算应该填充的令牌数
// - 如果桶中有令牌，消耗一个令牌并返回 true
// - 如果桶中没有令牌，返回 false
//
// 返回：
// - bool: 如果允许请求返回 true，否则返回 false
//
// 工作流程：
// 1. 计算自上次填充以来的时间差
// 2. 根据时间差和填充速率计算应填充的令牌数
// 3. 更新令牌数（不超过容量）
// 4. 如果有令牌，消耗一个并返回 true
// 5. 否则返回 false
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

// RateLimitMiddleware 创建限流中间件。
//
// 功能说明：
// - 根据配置的算法创建限流中间件
// - 支持内存限流和分布式限流（Redis）
// - 支持多种限流算法（令牌桶、滑动窗口、漏桶）
//
// 工作流程：
// 1. 检查路径是否在跳过列表中
// 2. 生成限流键（默认基于 IP 地址）
// 3. 根据算法检查是否允许请求
// 4. 如果限流，调用 OnLimitExceeded 处理函数
// 5. 如果允许，继续处理请求
//
// 分布式限流：
// - 如果配置了 RedisClient，使用 Redis 进行分布式限流
// - 多个服务实例共享限流状态
// - 使用滑动窗口算法实现
//
// 内存限流：
// - 如果未配置 RedisClient，使用内存限流
// - 每个服务实例独立限流
// - 支持三种算法：令牌桶、滑动窗口、漏桶
//
// 参数：
// - config: 限流配置
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	// 内存限流（令牌桶）
//	config := middleware.RateLimitConfig{
//	    RequestsPerSecond: 100,
//	    Burst:             50,
//	    Algorithm:         middleware.AlgorithmTokenBucket,
//	}
//	router.Use(middleware.RateLimitMiddleware(config))
//
//	// 分布式限流（Redis）
//	config := middleware.RateLimitConfig{
//	    RequestsPerSecond: 100,
//	    RedisClient:       redisClient,
//	    RedisKeyPrefix:    "ratelimit",
//	}
//	router.Use(middleware.RateLimitMiddleware(config))
//
// 注意事项：
// - 内存限流适用于单实例部署
// - 分布式限流适用于多实例部署
// - Redis 限流失败时会降级为允许请求（fail-open）
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

// defaultKeyFunc 是默认的限流键生成函数（基于 IP 地址）。
//
// 功能说明：
// - 根据请求的 IP 地址生成限流键
// - 支持代理场景（X-Forwarded-For header）
//
// 提取顺序：
// 1. X-Forwarded-For header（如果存在，说明请求经过代理）
// 2. RemoteAddr（直接连接的客户端地址）
//
// 参数：
// - r: HTTP 请求
//
// 返回：
// - string: 限流键（IP 地址）
//
// 使用场景：
// - 按 IP 地址限流
// - 防止单个 IP 地址的恶意请求
func defaultKeyFunc(r *http.Request) string {
	// 优先使用X-Forwarded-For（代理场景）
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	// 使用RemoteAddr（直接连接）
	return r.RemoteAddr
}

// defaultOnLimitExceeded 是默认的限流处理函数。
//
// 功能说明：
// - 当请求被限流时调用
// - 返回 429 Too Many Requests 错误
//
// 参数：
// - w: HTTP 响应写入器
// - r: HTTP 请求
//
// 使用场景：
// - 限流中间件的默认处理
// - 可以自定义以提供更详细的错误信息
func defaultOnLimitExceeded(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusTooManyRequests,
		errors.NewRateLimitError("rate limit exceeded"))
}

// shouldSkipRateLimit 检查指定路径是否应该跳过限流。
//
// 功能说明：
// - 检查路径是否在跳过列表中
// - 支持精确匹配
//
// 匹配规则：
// - 精确匹配：路径完全相等
// - 带斜杠匹配：路径相等或路径+"/"
//
// 参数：
// - path: 要检查的路径
// - skipPaths: 跳过限流的路径列表
//
// 返回：
// - bool: 如果应该跳过限流返回 true，否则返回 false
func shouldSkipRateLimit(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || path == skipPath+"/" {
			return true
		}
	}
	return false
}

// SlidingWindow 是滑动窗口限流器的实现。
//
// 算法原理：
// - 维护一个时间窗口内的请求时间戳列表
// - 每次请求时，移除窗口外的请求
// - 如果窗口内的请求数小于限制，允许请求；否则拒绝
//
// 特点：
// - 精确控制时间窗口内的请求数
// - 适合需要精确限流的场景
// - 内存占用与请求数相关
//
// 字段说明：
// - requests: 请求时间戳列表
// - limit: 时间窗口内的最大请求数
// - window: 时间窗口大小
// - mu: 互斥锁（保证并发安全）
type SlidingWindow struct {
	requests []time.Time   // 请求时间戳列表
	limit    int           // 限制数量
	window   time.Duration // 时间窗口
	mu       sync.Mutex
}

// NewSlidingWindow 创建并初始化滑动窗口限流器。
//
// 参数：
// - limit: 时间窗口内的最大请求数
// - window: 时间窗口大小
//
// 返回：
// - *SlidingWindow: 配置好的滑动窗口限流器实例
//
// 使用示例：
//
//	limiter := NewSlidingWindow(100, time.Second) // 每秒最多100个请求
func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		requests: make([]time.Time, 0, limit),
		limit:    limit,
		window:   window,
	}
}

// Allow 检查是否允许请求通过。
//
// 功能说明：
// - 移除时间窗口外的请求时间戳
// - 如果窗口内的请求数小于限制，添加当前请求并返回 true
// - 如果窗口内的请求数已达到限制，返回 false
//
// 返回：
// - bool: 如果允许请求返回 true，否则返回 false
//
// 工作流程：
// 1. 计算窗口的起始时间（当前时间 - 窗口大小）
// 2. 移除窗口外的请求时间戳
// 3. 检查窗口内的请求数是否小于限制
// 4. 如果小于限制，添加当前请求并返回 true
// 5. 否则返回 false
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

// LeakyBucket 是漏桶限流器的实现。
//
// 算法原理：
// - 请求进入桶中（增加水量）
// - 按固定速率从桶中漏出（减少水量）
// - 如果桶未满，允许请求；否则拒绝请求
//
// 特点：
// - 平滑输出，限制突发流量
// - 输出速率恒定
// - 适合需要平滑输出的场景
//
// 字段说明：
// - capacity: 桶的容量（最大水量）
// - leakRate: 漏出速率（请求/秒）
// - water: 当前水量
// - lastLeak: 上次漏水时间
// - mu: 互斥锁（保证并发安全）
type LeakyBucket struct {
	capacity int       // 容量
	leakRate float64   // 漏出速率（请求/秒）
	water    float64   // 当前水量
	lastLeak time.Time // 上次漏水时间
	mu       sync.Mutex
}

// NewLeakyBucket 创建并初始化漏桶限流器。
//
// 参数：
// - capacity: 桶的容量（最大水量）
// - leakRate: 漏出速率（请求/秒）
//
// 返回：
// - *LeakyBucket: 配置好的漏桶限流器实例
//
// 使用示例：
//
//	bucket := NewLeakyBucket(100, 10.0) // 容量100，每秒漏出10个请求
func NewLeakyBucket(capacity int, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		water:    0, // 初始时桶是空的
		lastLeak: time.Now(),
	}
}

// Allow 检查是否允许请求通过。
//
// 功能说明：
// - 根据时间差计算应该漏出的水量
// - 如果桶未满，增加水量并返回 true
// - 如果桶已满，返回 false
//
// 返回：
// - bool: 如果允许请求返回 true，否则返回 false
//
// 工作流程：
// 1. 计算自上次漏水以来的时间差
// 2. 根据时间差和漏出速率计算应漏出的水量
// 3. 更新水量（不能小于0）
// 4. 如果桶未满，增加水量并返回 true
// 5. 否则返回 false
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

// RedisRateLimiter 是 Redis 分布式限流器的实现。
//
// 功能说明：
// - 使用 Redis 有序集合实现分布式限流
// - 多个服务实例共享限流状态
// - 使用滑动窗口算法
//
// 实现原理：
// - 使用 Redis 有序集合存储请求时间戳
// - 使用时间戳作为分数，请求标识作为成员
// - 通过 ZRemRangeByScore 移除窗口外的请求
// - 通过 ZCard 统计窗口内的请求数
//
// 字段说明：
// - client: Redis 客户端
// - keyPrefix: Redis 键前缀
// - limit: 时间窗口内的最大请求数
// - window: 时间窗口大小
type RedisRateLimiter struct {
	client    RedisClient
	keyPrefix string
	limit     int
	window    time.Duration
}

// NewRedisRateLimiter 创建并初始化 Redis 分布式限流器。
//
// 参数：
// - client: Redis 客户端
// - keyPrefix: Redis 键前缀
// - limit: 时间窗口内的最大请求数
// - window: 时间窗口大小
//
// 返回：
// - *RedisRateLimiter: 配置好的 Redis 限流器实例
//
// 使用示例：
//
//	limiter := NewRedisRateLimiter(redisClient, "ratelimit", 100, time.Second)
func NewRedisRateLimiter(client RedisClient, keyPrefix string, limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:    client,
		keyPrefix: keyPrefix,
		limit:     limit,
		window:    window,
	}
}

// Allow 检查是否允许请求（使用滑动窗口算法）。
//
// 功能说明：
// - 使用 Redis Pipeline 批量执行命令
// - 移除时间窗口外的请求记录
// - 统计窗口内的请求数
// - 添加当前请求记录
// - 设置键的过期时间
//
// 参数：
// - ctx: 上下文
// - key: 限流键（如 IP 地址）
//
// 返回：
// - bool: 如果允许请求返回 true，否则返回 false
// - error: 如果 Redis 操作失败，返回错误
//
// 工作流程：
// 1. 构建完整的 Redis 键（prefix:key）
// 2. 计算窗口的起始时间
// 3. 使用 Pipeline 批量执行：
//    - 移除窗口外的记录（ZRemRangeByScore）
//    - 统计窗口内的请求数（ZCard）
//    - 添加当前请求（ZAdd）
//    - 设置过期时间（Expire）
// 4. 执行 Pipeline
// 5. 检查请求数是否小于限制
//
// 注意事项：
// - 使用 Pipeline 减少网络往返
// - 设置过期时间避免键无限增长
// - 错误时返回 false 和错误信息
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
