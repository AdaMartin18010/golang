# 限流与节流 (Rate Limiting & Throttling)

> **分类**: 工程与云原生
> **标签**: #rate-limiting #throttling #token-bucket #leaky-bucket
> **参考**: Token Bucket Algorithm, Rate Limiter Patterns

---

## 限流算法

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Rate Limiting Algorithms                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Token Bucket (令牌桶)                                                   │
│                                                                              │
│     ┌─────────────┐                                                         │
│     │   Bucket    │  capacity: 100 tokens                                  │
│     │  ┌───┬───┐  │  rate: 10 tokens/second                                │
│     │  │███│   │  │                                                         │
│     │  │███│   │  │  Request ──► Take Token ──► Process                   │
│     │  │███│   │  │             (if available)                             │
│     │  └───┴───┘  │                                                         │
│     └─────────────┘                                                         │
│                                                                              │
│  2. Leaky Bucket (漏桶)                                                     │
│                                                                              │
│     ┌─────────────┐                                                         │
│     │   Bucket    │  capacity: 100 requests                                │
│     │  ┌───────┐  │  leak rate: 10 req/sec                                 │
│     │  │▓▓▓▓▓▓▓│  │                                                         │
│     │  │▓▓▓▓▓▓▓│  │  Request ──► Add to Queue ──► Process (leak)          │
│     │  │▓▓▓▓▓▓▓│  │             (if queue not full)                        │
│     │  └───────┘  │                                                         │
│     └─────────────┘                                                         │
│                                                                              │
│  3. Fixed Window (固定窗口)                                                  │
│                                                                              │
│     Window: [00:00:00 - 00:00:59]  limit: 100                               │
│     Window: [00:01:00 - 00:01:59]  limit: 100                               │
│                                                                              │
│     ┌───┬───┬───┐                                                           │
│     │███│███│░░░│  ███ = used, ░░░ = available                             │
│     └───┴───┴───┘                                                           │
│                                                                              │
│  4. Sliding Window Log (滑动窗口日志)                                        │
│                                                                              │
│     Current time: 12:00:30                                                   │
│     Window: [11:59:30 - 12:00:30]  limit: 100                               │
│                                                                              │
│     ┌─────────────────────────────────────┐                                  │
│     │ 11:59:35  11:59:45  11:59:55  12:00:05 │  (timestamps of requests)     │
│     └─────────────────────────────────────┘                                  │
│                                                                              │
│  5. Sliding Window Counter (滑动窗口计数)                                    │
│                                                                              │
│     Current window: 60% into [12:00:00 - 12:01:00]                          │
│     Previous window: [11:59:00 - 12:00:00]                                  │
│                                                                              │
│     Count = (prev_count * (1 - 0.6)) + curr_count                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心实现

```go
package ratelimit

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Limiter 限流器接口
type Limiter interface {
    Allow() bool
    AllowN(n int) bool
    Wait(ctx context.Context) error
    WaitN(ctx context.Context, n int) error
}

// TokenBucket 令牌桶限流器
type TokenBucket struct {
    capacity   int64
    tokens     int64
    fillRate   float64 // tokens per second
    lastFill   time.Time

    mu         sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity int64, fillRate float64) *TokenBucket {
    return &TokenBucket{
        capacity: capacity,
        tokens:   capacity,
        fillRate: fillRate,
        lastFill: time.Now(),
    }
}

// Allow 获取一个令牌
func (tb *TokenBucket) Allow() bool {
    return tb.AllowN(1)
}

// AllowN 获取 N 个令牌
func (tb *TokenBucket) AllowN(n int) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    tb.refill()

    if tb.tokens >= int64(n) {
        tb.tokens -= int64(n)
        return true
    }

    return false
}

// Wait 等待获取令牌
func (tb *TokenBucket) Wait(ctx context.Context) error {
    return tb.WaitN(ctx, 1)
}

// WaitN 等待获取 N 个令牌
func (tb *TokenBucket) WaitN(ctx context.Context, n int) error {
    for {
        if tb.AllowN(n) {
            return nil
        }

        // 计算等待时间
        tb.mu.Lock()
        need := int64(n) - tb.tokens
        waitTime := time.Duration(float64(need) / tb.fillRate * float64(time.Second))
        tb.mu.Unlock()

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(waitTime):
        }
    }
}

func (tb *TokenBucket) refill() {
    now := time.Now()
    elapsed := now.Sub(tb.lastFill).Seconds()

    tokensToAdd := int64(elapsed * tb.fillRate)
    if tokensToAdd > 0 {
        tb.tokens += tokensToAdd
        if tb.tokens > tb.capacity {
            tb.tokens = tb.capacity
        }
        tb.lastFill = now
    }
}

// LeakyBucket 漏桶限流器
type LeakyBucket struct {
    capacity   int64
    queue      int64
    leakRate   float64 // requests per second
    lastLeak   time.Time

    mu         sync.Mutex
}

// NewLeakyBucket 创建漏桶
func NewLeakyBucket(capacity int64, leakRate float64) *LeakyBucket {
    return &LeakyBucket{
        capacity: capacity,
        leakRate: leakRate,
        lastLeak: time.Now(),
    }
}

// Allow 尝试添加请求
func (lb *LeakyBucket) Allow() bool {
    return lb.AllowN(1)
}

// AllowN 尝试添加 N 个请求
func (lb *LeakyBucket) AllowN(n int) bool {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    lb.leak()

    if lb.queue+int64(n) <= lb.capacity {
        lb.queue += int64(n)
        return true
    }

    return false
}

// Wait 等待添加请求
func (lb *LeakyBucket) Wait(ctx context.Context) error {
    return lb.WaitN(ctx, 1)
}

// WaitN 等待添加 N 个请求
func (lb *LeakyBucket) WaitN(ctx context.Context, n int) error {
    for {
        if lb.AllowN(n) {
            return nil
        }

        lb.mu.Lock()
        space := lb.capacity - lb.queue
        need := int64(n) - space
        waitTime := time.Duration(float64(need) / lb.leakRate * float64(time.Second))
        lb.mu.Unlock()

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(waitTime):
        }
    }
}

func (lb *LeakyBucket) leak() {
    now := time.Now()
    elapsed := now.Sub(lb.lastLeak).Seconds()

    leaked := int64(elapsed * lb.leakRate)
    if leaked > 0 {
        lb.queue -= leaked
        if lb.queue < 0 {
            lb.queue = 0
        }
        lb.lastLeak = now
    }
}
```

---

## 分布式限流（Redis）

```go
package ratelimit

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisRateLimiter Redis 分布式限流器
type RedisRateLimiter struct {
    client   *redis.Client
    key      string
    rate     int64   // 每秒请求数
    burst    int64   // 突发容量
    window   time.Duration
}

// NewRedisRateLimiter 创建 Redis 限流器
func NewRedisRateLimiter(client *redis.Client, key string, rate int64, burst int64) *RedisRateLimiter {
    return &RedisRateLimiter{
        client: client,
        key:    key,
        rate:   rate,
        burst:  burst,
        window: time.Second,
    }
}

// Allow 滑动窗口限流
func (rl *RedisRateLimiter) Allow(ctx context.Context) bool {
    now := time.Now().UnixMilli()
    windowStart := now - rl.window.Milliseconds()

    pipe := rl.client.Pipeline()

    // 移除窗口外的请求记录
    pipe.ZRemRangeByScore(ctx, rl.key, "0", fmt.Sprintf("%d", windowStart))

    // 获取当前窗口内的请求数
    countCmd := pipe.ZCard(ctx, rl.key)

    // 添加当前请求
    pipe.ZAdd(ctx, rl.key, redis.Z{
        Score:  float64(now),
        Member: fmt.Sprintf("%d-%d", now, rand.Int()),
    })

    // 设置过期时间
    pipe.Expire(ctx, rl.key, rl.window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false
    }

    count := countCmd.Val()
    return count < rl.rate
}

// AllowN 令牌桶限流（Redis 实现）
func (rl *RedisRateLimiter) AllowN(ctx context.Context, n int64) bool {
    script := `
        local key = KEYS[1]
        local rate = tonumber(ARGV[1])
        local capacity = tonumber(ARGV[2])
        local now = tonumber(ARGV[3])
        local requested = tonumber(ARGV[4])

        local fill_time = capacity / rate
        local ttl = math.floor(fill_time * 2)

        local last_tokens = redis.call("get", key)
        if last_tokens == false then
            last_tokens = capacity
        end

        local last_updated = redis.call("get", key .. ":last_updated")
        if last_updated == false then
            last_updated = 0
        end

        local delta = math.max(0, now - tonumber(last_updated))
        local filled_tokens = math.min(capacity, tonumber(last_tokens) + (delta * rate))
        local allowed = filled_tokens >= requested
        local new_tokens = filled_tokens

        if allowed then
            new_tokens = filled_tokens - requested
        end

        redis.call("setex", key, ttl, new_tokens)
        redis.call("setex", key .. ":last_updated", ttl, now)

        return allowed
    `

    now := time.Now().Unix()
    result, err := rl.client.Eval(ctx, script, []string{rl.key},
        rl.rate, rl.burst, now, n).Result()

    if err != nil {
        return false
    }

    return result.(int64) == 1
}

// FixedWindow 固定窗口限流
func (rl *RedisRateLimiter) FixedWindow(ctx context.Context) bool {
    window := time.Now().Truncate(time.Second).Unix()
    key := fmt.Sprintf("%s:%d", rl.key, window)

    current, err := rl.client.Incr(ctx, key).Result()
    if err != nil {
        return false
    }

    if current == 1 {
        rl.client.Expire(ctx, key, time.Second)
    }

    return current <= rl.rate
}
```

---

## HTTP 中间件

```go
package ratelimit

import (
    "net/http"
    "strconv"
    "strings"
    "sync"
)

// Middleware HTTP 限流中间件
func Middleware(limiter Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                w.Header().Set("Retry-After", "1")
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// PerClientLimiter 按客户端限流
type PerClientLimiter struct {
    limiters map[string]Limiter
    factory  func() Limiter
    mu       sync.RWMutex
}

func NewPerClientLimiter(factory func() Limiter) *PerClientLimiter {
    return &PerClientLimiter{
        limiters: make(map[string]Limiter),
        factory:  factory,
    }
}

func (p *PerClientLimiter) GetLimiter(clientID string) Limiter {
    p.mu.RLock()
    limiter, ok := p.limiters[clientID]
    p.mu.RUnlock()

    if ok {
        return limiter
    }

    p.mu.Lock()
    defer p.mu.Unlock()

    limiter, ok = p.limiters[clientID]
    if ok {
        return limiter
    }

    limiter = p.factory()
    p.limiters[clientID] = limiter
    return limiter
}

func (p *PerClientLimiter) Allow(clientID string) bool {
    return p.GetLimiter(clientID).Allow()
}

// PerClientMiddleware 按客户端限流中间件
func PerClientMiddleware(factory func() Limiter, keyFunc func(*http.Request) string) func(http.Handler) http.Handler {
    limiters := NewPerClientLimiter(factory)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            clientID := keyFunc(r)
            if !limiters.Allow(clientID) {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// ExtractClientIP 提取客户端 IP
func ExtractClientIP(r *http.Request) string {
    // 检查 X-Forwarded-For
    xff := r.Header.Get("X-Forwarded-For")
    if xff != "" {
        parts := strings.Split(xff, ",")
        return strings.TrimSpace(parts[0])
    }

    // 检查 X-Real-IP
    xri := r.Header.Get("X-Real-Ip")
    if xri != "" {
        return xri
    }

    // 使用 RemoteAddr
    ip := r.RemoteAddr
    if idx := strings.LastIndex(ip, ":"); idx != -1 {
        ip = ip[:idx]
    }
    return ip
}
```

---

## 自适应限流

```go
package ratelimit

import (
    "math"
    "sync"
    "sync/atomic"
    "time"
)

// AdaptiveLimiter 自适应限流器
type AdaptiveLimiter struct {
    // 配置
    minRate     float64
    maxRate     float64
    targetRTT   time.Duration

    // 状态
    currentRate  int64 // 当前限流值
    requestCount int64
    errorCount   int64
    totalRTT     int64

    mu           sync.RWMutex
    lastAdjust   time.Time
}

// NewAdaptiveLimiter 创建自适应限流器
func NewAdaptiveLimiter(minRate, maxRate float64, targetRTT time.Duration) *AdaptiveLimiter {
    return &AdaptiveLimiter{
        minRate:     minRate,
        maxRate:     maxRate,
        targetRTT:   targetRTT,
        currentRate: int64(maxRate),
        lastAdjust:  time.Now(),
    }
}

// Allow 检查是否允许请求
func (al *AdaptiveLimiter) Allow() bool {
    // 简单实现：基于当前速率限制
    // 实际应使用令牌桶等算法
    atomic.AddInt64(&al.requestCount, 1)
    return true
}

// RecordResult 记录请求结果
func (al *AdaptiveLimiter) RecordResult(rtt time.Duration, err error) {
    atomic.AddInt64(&al.totalRTT, int64(rtt))
    if err != nil {
        atomic.AddInt64(&al.errorCount, 1)
    }

    // 定期调整
    if time.Since(al.lastAdjust) > time.Second {
        al.adjust()
    }
}

func (al *AdaptiveLimiter) adjust() {
    al.mu.Lock()
    defer al.mu.Unlock()

    requests := atomic.SwapInt64(&al.requestCount, 0)
    errors := atomic.SwapInt64(&al.errorCount, 0)
    totalRTT := atomic.SwapInt64(&al.totalRTT, 0)

    if requests == 0 {
        return
    }

    errorRate := float64(errors) / float64(requests)
    avgRTT := time.Duration(totalRTT / requests)

    currentRate := float64(atomic.LoadInt64(&al.currentRate))

    // 根据错误率和延迟调整
    if errorRate > 0.1 || avgRTT > al.targetRTT*2 {
        // 降低速率
        currentRate *= 0.9
    } else if errorRate < 0.01 && avgRTT < al.targetRTT {
        // 提高速率
        currentRate *= 1.1
    }

    // 限制范围
    currentRate = math.Max(al.minRate, math.Min(al.maxRate, currentRate))
    atomic.StoreInt64(&al.currentRate, int64(currentRate))

    al.lastAdjust = time.Now()
}

func (al *AdaptiveLimiter) CurrentRate() float64 {
    return float64(atomic.LoadInt64(&al.currentRate))
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "ratelimit"
)

func main() {
    // 创建令牌桶限流器: 100/秒，突发 200
    limiter := ratelimit.NewTokenBucket(200, 100)

    // 全局限流中间件
    http.Handle("/", ratelimit.Middleware(limiter)(http.HandlerFunc(handler)))

    // 按客户端限流
    perClientLimiter := ratelimit.PerClientMiddleware(
        func() ratelimit.Limiter {
            return ratelimit.NewTokenBucket(20, 10)
        },
        ratelimit.ExtractClientIP,
    )

    http.Handle("/api/", perClientLimiter(http.HandlerFunc(apiHandler)))

    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "API")
}

func exampleUsage() {
    // 令牌桶
    tb := ratelimit.NewTokenBucket(100, 10)

    // 同步限流
    if tb.Allow() {
        // 处理请求
    }

    // 阻塞等待
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    tb.Wait(ctx)

    // 漏桶
    lb := ratelimit.NewLeakyBucket(100, 10)
    lb.Allow()
}
```
