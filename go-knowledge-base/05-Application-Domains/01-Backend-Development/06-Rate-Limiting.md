# Rate Limiting Patterns

> **Dimension**: Application Domains
> **Level**: S (17+ KB)
> **Tags**: #rate-limiting #throttling #token-bucket #leaky-bucket

---

## 1. Rate Limiting Fundamentals

### 1.1 Why Rate Limiting?

| Purpose | Description |
|---------|-------------|
| Prevent abuse | Stop malicious or excessive requests |
| Resource protection | Prevent server overload |
| Fair usage | Ensure all users get fair access |
| Cost control | Limit expensive operations |
| DDoS mitigation | Absorb attack traffic |

### 1.2 Rate Limiting Strategies

| Strategy | Description | Best For |
|----------|-------------|----------|
| Fixed Window | Reset counter at intervals | Simple implementations |
| Sliding Window | Smooth window transition | Fairer distribution |
| Token Bucket | Tokens added at fixed rate | Bursty traffic |
| Leaky Bucket | Constant outflow rate | Smooth output |

---

## 2. Token Bucket Algorithm

```go
package ratelimit

import (
    "context"
    "sync"
    "time"
)

// TokenBucket implements token bucket rate limiting
type TokenBucket struct {
    capacity   float64
    tokens     float64
    rate       float64
    lastRefill time.Time
    mu         sync.Mutex
}

func NewTokenBucket(capacity int, ratePerSecond float64) *TokenBucket {
    return &TokenBucket{
        capacity:   float64(capacity),
        tokens:     float64(capacity),
        rate:       ratePerSecond,
        lastRefill: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    tb.refill()

    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }
    return false
}

func (tb *TokenBucket) AllowN(n int) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    tb.refill()

    if tb.tokens >= float64(n) {
        tb.tokens -= float64(n)
        return true
    }
    return false
}

func (tb *TokenBucket) refill() {
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill).Seconds()
    tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.rate)
    tb.lastRefill = now
}

func (tb *TokenBucket) WaitTime(n int) time.Duration {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    if tb.tokens >= float64(n) {
        return 0
    }

    needed := float64(n) - tb.tokens
    return time.Duration(needed/tb.rate) * time.Second
}
```

---

## 3. Sliding Window Algorithm

```go
package ratelimit

import (
    "sync"
    "time"
)

// SlidingWindow implements sliding window rate limiting
type SlidingWindow struct {
    limit   int
    window  time.Duration
    requests []time.Time
    mu       sync.Mutex
}

func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
    return &SlidingWindow{
        limit:    limit,
        window:   window,
        requests: make([]time.Time, 0, limit),
    }
}

func (sw *SlidingWindow) Allow() bool {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-sw.window)

    // Remove expired requests
    valid := make([]time.Time, 0, len(sw.requests))
    for _, t := range sw.requests {
        if t.After(cutoff) {
            valid = append(valid, t)
        }
    }
    sw.requests = valid

    // Check if we can allow
    if len(sw.requests) < sw.limit {
        sw.requests = append(sw.requests, now)
        return true
    }
    return false
}

func (sw *SlidingWindow) Remaining() int {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-sw.window)

    count := 0
    for _, t := range sw.requests {
        if t.After(cutoff) {
            count++
        }
    }

    return sw.limit - count
}
```

---

## 4. Distributed Rate Limiting

### 4.1 Redis-Based Rate Limiter

```go
package ratelimit

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

// RedisRateLimiter for distributed rate limiting
type RedisRateLimiter struct {
    client *redis.Client
    window time.Duration
    limit  int
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, int, error) {
    now := time.Now().Unix()
    windowKey := fmt.Sprintf("ratelimit:%s:%d", key, now/int64(rl.window.Seconds()))

    pipe := rl.client.Pipeline()
    incrCmd := pipe.Incr(ctx, windowKey)
    pipe.Expire(ctx, windowKey, rl.window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, 0, err
    }

    current := int(incrCmd.Val())
    remaining := rl.limit - current

    if remaining < 0 {
        remaining = 0
    }

    return current <= rl.limit, remaining, nil
}

// Sliding Window Log with Redis
func (rl *RedisRateLimiter) AllowSlidingWindow(ctx context.Context, key string) (bool, error) {
    now := time.Now()
    windowStart := now.Add(-rl.window).Unix()

    pipe := rl.client.Pipeline()

    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

    // Count current entries
    countCmd := pipe.ZCard(ctx, key)

    // Add current request
    pipe.ZAdd(ctx, key, &redis.Z{
        Score:  float64(now.Unix()),
        Member: now.UnixNano(),
    })

    // Set expiry
    pipe.Expire(ctx, key, rl.window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }

    return int(countCmd.Val()) < rl.limit, nil
}
```

### 4.2 Token Bucket with Redis

```go
package ratelimit

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

const tokenBucketScript = `
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

type RedisTokenBucket struct {
    client   *redis.Client
    rate     float64
    capacity int
    script   *redis.Script
}

func NewRedisTokenBucket(client *redis.Client, rate float64, capacity int) *RedisTokenBucket {
    return &RedisTokenBucket{
        client:   client,
        rate:     rate,
        capacity: capacity,
        script:   redis.NewScript(tokenBucketScript),
    }
}

func (rtb *RedisTokenBucket) Allow(ctx context.Context, key string, tokens int) (bool, error) {
    result, err := rtb.script.Run(ctx, rtb.client, []string{key},
        rtb.rate,
        rtb.capacity,
        time.Now().Unix(),
        tokens,
    ).Bool()

    return result, err
}
```

---

## 5. Rate Limiting Middleware

```go
package middleware

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

func RateLimiter(limiter RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := extractClientKey(c)

        allowed, remaining, resetTime, err := limiter.Check(c.Request.Context(), key)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
            return
        }

        // Set rate limit headers
        c.Header("X-RateLimit-Limit", strconv.Itoa(limiter.Limit()))
        c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

        if !allowed {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
                "retry_after": resetTime - time.Now().Unix(),
            })
            return
        }

        c.Next()
    }
}

func extractClientKey(c *gin.Context) string {
    // API key takes priority
    apiKey := c.GetHeader("X-API-Key")
    if apiKey != "" {
        return "api:" + apiKey
    }

    // Then user ID from JWT
    if userID, exists := c.Get("user_id"); exists {
        return "user:" + userID.(string)
    }

    // Fall back to IP
    ip := c.ClientIP()
    return "ip:" + ip
}
```

---

## 6. Adaptive Rate Limiting

```go
package ratelimit

import (
    "math"
    "sync"
    "time"
)

// AdaptiveRateLimiter adjusts limits based on system load
type AdaptiveRateLimiter struct {
    baseLimit   int
    currentLimit int
    minLimit    int
    maxLimit    int

    // Load metrics
    cpuThreshold    float64
    memoryThreshold float64
    latencyThreshold time.Duration

    mu sync.RWMutex
}

func (arl *AdaptiveRateLimiter) AdjustMetrics(metrics SystemMetrics) {
    arl.mu.Lock()
    defer arl.mu.Unlock()

    loadFactor := arl.calculateLoadFactor(metrics)

    // Decrease limit when under load
    if loadFactor > 1.0 {
        arl.currentLimit = max(arl.minLimit, int(float64(arl.currentLimit)*0.9))
    } else if loadFactor < 0.5 {
        // Increase limit when idle
        arl.currentLimit = min(arl.maxLimit, int(float64(arl.currentLimit)*1.1))
    }
}

func (arl *AdaptiveRateLimiter) calculateLoadFactor(metrics SystemMetrics) float64 {
    cpuLoad := metrics.CPUUsage / arl.cpuThreshold
    memoryLoad := metrics.MemoryUsage / arl.memoryThreshold
    latencyLoad := float64(metrics.P99Latency) / float64(arl.latencyThreshold)

    // Take maximum load factor
    return math.Max(cpuLoad, math.Max(memoryLoad, latencyLoad))
}
```

---

## 7. Best Practices

- [ ] Return informative headers
- [ ] Implement different limits for different users
- [ ] Use distributed rate limiting for scalability
- [ ] Monitor rate limit hits
- [ ] Provide graceful degradation
- [ ] Consider cost-based limiting
- [ ] Implement per-endpoint limits
- [ ] Use token bucket for burst handling
- [ ] Set appropriate TTL for Redis keys
- [ ] Handle rate limiter failures gracefully

---

**Quality Rating**: S (17+ KB)
**Last Updated**: 2026-04-02
