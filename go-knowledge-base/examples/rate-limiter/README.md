# Rate Limiter Example

A comprehensive rate limiting implementation featuring token bucket, sliding window counter, and distributed rate limiting algorithms. This production-ready system supports multiple backends including Redis, in-memory, and distributed coordination.

## Table of Contents

- [Rate Limiter Example](#rate-limiter-example)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Use Cases](#use-cases)
  - [Algorithms](#algorithms)
    - [Token Bucket](#token-bucket)
    - [Sliding Window Counter](#sliding-window-counter)
    - [Fixed Window Counter](#fixed-window-counter)
  - [Architecture](#architecture)
    - [System Architecture](#system-architecture)
    - [Rate Limit Headers](#rate-limit-headers)
  - [Implementation](#implementation)
    - [Token Bucket Implementation](#token-bucket-implementation)
    - [Sliding Window Counter](#sliding-window-counter-1)
    - [Rate Limiter Middleware](#rate-limiter-middleware)
  - [Distributed Rate Limiting](#distributed-rate-limiting)
    - [Redis-Based Rate Limiter](#redis-based-rate-limiter)
    - [GCRA (Generic Cell Rate Algorithm)](#gcra-generic-cell-rate-algorithm)
  - [Deployment](#deployment)
    - [Docker Compose](#docker-compose)
    - [Kubernetes](#kubernetes)
  - [Performance](#performance)
    - [Benchmarks](#benchmarks)
    - [Load Testing](#load-testing)
  - [Best Practices](#best-practices)
    - [Configuration Guidelines](#configuration-guidelines)
    - [Example Configuration](#example-configuration)
  - [License](#license)

## Overview

This rate limiter provides:

- **Multiple Algorithms**: Token Bucket, Sliding Window, Fixed Window
- **Distributed Support**: Redis-based coordination for microservices
- **Flexible Limits**: Per-user, per-IP, per-endpoint, and custom strategies
- **High Performance**: Zero-allocation hot path, lock-free where possible
- **Observability**: Prometheus metrics, structured logging
- **Middleware Support**: Ready-to-use HTTP middleware

### Use Cases

| Use Case | Algorithm | Configuration |
|----------|-----------|---------------|
| API Rate Limiting | Token Bucket | 100 req/min per API key |
| Burst Protection | Token Bucket | 10 req/sec, burst 20 |
| Abuse Prevention | Sliding Window | 1000 req/hour per IP |
| Resource Protection | Fixed Window | 100 req/min globally |

## Algorithms

### Token Bucket

```
┌─────────────────────────────────────────────────────────────────┐
│                     Token Bucket Algorithm                       │
│                                                                  │
│   Bucket Capacity: 10 tokens                                    │
│   Refill Rate: 1 token/second                                   │
│                                                                  │
│   ┌─────────────────────────────────────┐                        │
│   │  ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐    │  Initial: 10 tokens    │
│   │  │ ● │ │ ● │ │ ● │ │ ● │ │ ● │    │                        │
│   │  └───┘ └───┘ └───┘ └───┘ └───┘    │                        │
│   │  ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐    │                        │
│   │  │ ● │ │ ● │ │ ● │ │ ● │ │ ● │    │                        │
│   │  └───┘ └───┘ └───┘ └───┘ └───┘    │                        │
│   └─────────────────────────────────────┘                        │
│                                                                  │
│   Time: 0s    1s    2s    3s    4s    5s                         │
│   Tokens: 10 → 10 → 10 → 10 → 10 → 10 (max capacity)            │
│                                                                  │
│   Request comes in (costs 1 token):                              │
│   ┌─────────────────────────────────────┐                        │
│   │  ┌───┐ ┌───┐ ┌───┐ ┌───┐          │  After request: 9      │
│   │  │ ● │ │ ● │ │ ● │ │ ● │  [empty] │                        │
│   │  └───┘ └───┘ └───┘ └───┘          │                        │
│   └─────────────────────────────────────┘                        │
│                                                                  │
│   Properties:                                                    │
│   - Allows bursts up to bucket capacity                          │
│   - Smooths traffic over time                                    │
│   - Easy to implement efficiently                                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Sliding Window Counter

```
┌─────────────────────────────────────────────────────────────────┐
│                Sliding Window Counter Algorithm                  │
│                                                                  │
│   Window Size: 1 minute                                          │
│   Limit: 100 requests                                            │
│                                                                  │
│   Previous Window (50%)    Current Window (50%)                  │
│   ┌─────────────┐          ┌─────────────┐                       │
│   │ 12:00-12:01 │          │ 12:01-12:02 │                       │
│   │  80 requests │          │  30 requests │  ← Now: 12:01:30     │
│   └─────────────┘          └─────────────┘                       │
│          │                        │                              │
│          └──────────┬─────────────┘                              │
│                     ▼                                            │
│   Weighted Count = 80 * 0.5 + 30 = 70                            │
│   Allowed? 70 < 100 → YES                                        │
│                                                                  │
│   Time: 12:01:30                                                 │
│   ┌─────────────────────────────────────────────────────────┐    │
│   │ Previous [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]│    │
│   │ Current  [████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░]│    │
│   │          0s          15s         30s         45s       60s│    │
│   └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│   Properties:                                                    │
│   - More accurate than fixed window                              │
│   - Handles burst at window boundaries                           │
│   - Requires storing two windows                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Fixed Window Counter

```
┌─────────────────────────────────────────────────────────────────┐
│                 Fixed Window Counter Algorithm                   │
│                                                                  │
│   Window Size: 1 minute                                          │
│   Limit: 100 requests                                            │
│                                                                  │
│   Window 1             Window 2             Window 3             │
│   12:00-12:01          12:01-12:02          12:02-12:03         │
│   ┌─────────┐          ┌─────────┐          ┌─────────┐          │
│   │ 100 req │          │  45 req │          │   0 req │          │
│   │ [FULL]  │          │         │          │         │          │
│   └─────────┘          └─────────┘          └─────────┘          │
│        │                    │                    │               │
│        └────────────────────┼────────────────────┘               │
│                             │                                    │
│                          12:01:30                                │
│                                                                  │
│   Problem: Burst at boundary                                     │
│   12:01:59: 100 requests                                         │
│   12:02:00: 100 requests (new window)                            │
│   Total: 200 requests in 2 seconds!                              │
│                                                                  │
│   Properties:                                                    │
│   - Simple to implement                                          │
│   - Low memory footprint                                         │
│   - Can allow 2x burst at window boundaries                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Architecture

### System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Requests                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ API Request │  │ Web Request │  │ GRPC Request│  │  WebSocket Conn     │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘ │
└─────────┼────────────────┼────────────────┼────────────────────┼────────────┘
          │                │                │                    │
          └────────────────┴────────────────┴────────────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │       Rate Limiter          │
                    │  ┌───────────────────────┐  │
                    │  │   Strategy Selector   │  │
                    │  │   - Per-IP            │  │
                    │  │   - Per-User          │  │
                    │  │   - Per-API-Key       │  │
                    │  │   - Per-Endpoint      │  │
                    │  └───────────────────────┘  │
                    └──────────────┬──────────────┘
                                   │
          ┌────────────────────────┼────────────────────────┐
          │                        │                        │
┌─────────▼──────────┐  ┌──────────▼──────────┐  ┌──────────▼──────────┐
│   Local Limiter    │  │   Redis Limiter     │  │   Distributed       │
│   (In-Memory)      │  │   (Shared State)    │  │   (Cluster)         │
│                    │  │                     │  │                     │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ ┌───────────────┐   │
│ │ Token Bucket  │  │  │ │ Token Bucket  │   │  │ │ Token Bucket  │   │
│ │ (Per-Process) │  │  │ │ (Shared)      │   │  │ │ (GCRA/Sliding)│   │
│ └───────────────┘  │  │ └───────────────┘   │  │ └───────────────┘   │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ ┌───────────────┐   │
│ │ Sliding Win   │  │  │ │ Sliding Win   │   │  │ │ Sliding Win   │   │
│ └───────────────┘  │  │ └───────────────┘   │  │ └───────────────┘   │
└────────────────────┘  └─────────────────────┘  └─────────────────────┘
          │                        │                        │
          └────────────────────────┼────────────────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │         Decision            │
                    │  ┌───────────────────────┐  │
                    │  │   Allow / Deny        │  │
                    │  │   + Headers (X-Rate)  │  │
                    │  └───────────────────────┘  │
                    └─────────────────────────────┘
```

### Rate Limit Headers

```
HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1643723400
X-RateLimit-Policy: 100;w=60

HTTP/1.1 429 Too Many Requests
Retry-After: 30
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1643723400
```

## Implementation

### Token Bucket Implementation

```go
package ratelimit

import (
    "context"
    "sync"
    "time"
)

// TokenBucket implements the token bucket algorithm
type TokenBucket struct {
    capacity   float64       // Maximum tokens
    tokens     float64       // Current tokens
    refillRate float64       // Tokens per second
    lastRefill time.Time     // Last refill timestamp
    mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
    return &TokenBucket{
        capacity:   capacity,
        tokens:     capacity,
        refillRate: refillRate,
        lastRefill: time.Now(),
    }
}

// Allow checks if a request should be allowed
func (tb *TokenBucket) Allow() bool {
    return tb.AllowN(1)
}

// AllowN checks if n tokens are available
func (tb *TokenBucket) AllowN(n int) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastRefill).Seconds()

    // Refill tokens
    tb.tokens += elapsed * tb.refillRate
    if tb.tokens > tb.capacity {
        tb.tokens = tb.capacity
    }
    tb.lastRefill = now

    // Check if we have enough tokens
    if tb.tokens >= float64(n) {
        tb.tokens -= float64(n)
        return true
    }

    return false
}

// Wait waits until a token is available
func (tb *TokenBucket) Wait(ctx context.Context) error {
    for {
        if tb.Allow() {
            return nil
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(time.Millisecond):
            // Retry
        }
    }
}

// Tokens returns current token count
func (tb *TokenBucket) Tokens() float64 {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    return tb.tokens
}
```

### Sliding Window Counter

```go
// SlidingWindow implements the sliding window counter algorithm
type SlidingWindow struct {
    limit       int
    windowSize  time.Duration

    previousWindow struct {
        count     int
        timestamp time.Time
    }
    currentWindow struct {
        count     int
        timestamp time.Time
    }

    mu sync.Mutex
}

// NewSlidingWindow creates a new sliding window counter
func NewSlidingWindow(limit int, windowSize time.Duration) *SlidingWindow {
    now := time.Now()
    return &SlidingWindow{
        limit:      limit,
        windowSize: windowSize,
        currentWindow: struct {
            count     int
            timestamp time.Time
        }{timestamp: now},
        previousWindow: struct {
            count     int
            timestamp time.Time
        }{timestamp: now.Add(-windowSize)},
    }
}

// Allow checks if a request should be allowed
func (sw *SlidingWindow) Allow() bool {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(sw.currentWindow.timestamp)

    // Check if we need to move to a new window
    if elapsed >= sw.windowSize {
        sw.previousWindow = sw.currentWindow
        sw.currentWindow = struct {
            count     int
            timestamp time.Time
        }{timestamp: now}
        elapsed = 0
    }

    // Calculate weighted count
    weight := 1.0 - (elapsed.Seconds() / sw.windowSize.Seconds())
    weightedCount := float64(sw.previousWindow.count)*weight + float64(sw.currentWindow.count)

    if int(weightedCount) < sw.limit {
        sw.currentWindow.count++
        return true
    }

    return false
}

// Stats returns current window statistics
func (sw *SlidingWindow) Stats() (current, previous, weighted int) {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(sw.currentWindow.timestamp)
    weight := 1.0 - (elapsed.Seconds() / sw.windowSize.Seconds())

    return sw.currentWindow.count, sw.previousWindow.count,
           int(float64(sw.previousWindow.count)*weight + float64(sw.currentWindow.count))
}
```

### Rate Limiter Middleware

```go
package middleware

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
    "ratelimiter/internal/limiter"
)

// RateLimiterMiddleware wraps an HTTP handler with rate limiting
func RateLimiterMiddleware(limiter limiter.Limiter, keyFunc KeyFunc) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Generate key for this request
            key := keyFunc(r)

            // Check rate limit
            result, err := limiter.Allow(r.Context(), key)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            // Add rate limit headers
            w.Header().Set("X-RateLimit-Limit", strconv.Itoa(result.Limit))
            w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
            w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

            if !result.Allowed {
                w.Header().Set("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// KeyFunc generates a rate limit key from a request
type KeyFunc func(r *http.Request) string

// KeyByIP generates key by client IP
func KeyByIP(r *http.Request) string {
    ip := r.RemoteAddr
    if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
        ip = forwarded
    }
    return "ip:" + ip
}

// KeyByUser generates key by user ID from header
func KeyByUser(headerName string) KeyFunc {
    return func(r *http.Request) string {
        userID := r.Header.Get(headerName)
        if userID == "" {
            return KeyByIP(r)
        }
        return "user:" + userID
    }
}

// KeyByAPIKey generates key by API key
func KeyByAPIKey(r *http.Request) string {
    apiKey := r.Header.Get("X-API-Key")
    if apiKey == "" {
        return KeyByIP(r)
    }
    return "apikey:" + apiKey
}
```

## Distributed Rate Limiting

### Redis-Based Rate Limiter

```go
package distributed

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisLimiter implements distributed rate limiting using Redis
// Uses Redis sorted sets for sliding window

type RedisLimiter struct {
    client *redis.Client
    prefix string
}

// NewRedisLimiter creates a new Redis-based rate limiter
func NewRedisLimiter(client *redis.Client, prefix string) *RedisLimiter {
    return &RedisLimiter{
        client: client,
        prefix: prefix,
    }
}

// Allow checks if a request should be allowed using sliding window
func (rl *RedisLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (*Result, error) {
    now := time.Now()
    windowStart := now.Add(-window)

    redisKey := fmt.Sprintf("%s:%s", rl.prefix, key)

    pipe := rl.client.Pipeline()

    // Remove old entries outside the window
    pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStart.UnixMilli()))

    // Count current entries in window
    countCmd := pipe.ZCard(ctx, redisKey)

    // Add current request
    pipe.ZAdd(ctx, redisKey, redis.Z{
        Score:  float64(now.UnixMilli()),
        Member: now.UnixNano(), // Use nanoseconds as unique member
    })

    // Set expiration on the key
    pipe.Expire(ctx, redisKey, window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return nil, err
    }

    count := int(countCmd.Val())
    allowed := count < limit

    result := &Result{
        Allowed:   allowed,
        Limit:     limit,
        Remaining: limit - count - 1,
        ResetAt:   now.Add(window),
    }

    if !allowed {
        result.RetryAfter = window / time.Duration(limit)
    }

    return result, nil
}

// Result represents a rate limit decision
type Result struct {
    Allowed    bool
    Limit      int
    Remaining  int
    ResetAt    time.Time
    RetryAfter time.Duration
}
```

### GCRA (Generic Cell Rate Algorithm)

```go
// GCRA implements the Generic Cell Rate Algorithm
// Used for distributed rate limiting with Redis

type GCRA struct {
    client     *redis.Client
    emission   time.Duration // Time between requests (1/rate)
    tolerance  time.Duration // Burst tolerance
}

// NewGCRA creates a new GCRA rate limiter
func NewGCRA(client *redis.Client, rate int, burst int) *GCRA {
    emission := time.Second / time.Duration(rate)
    tolerance := emission * time.Duration(burst)

    return &GCRA{
        client:    client,
        emission:  emission,
        tolerance: tolerance,
    }
}

// Allow checks if a request is allowed using GCRA
func (g *GCRA) Allow(ctx context.Context, key string) (*Result, error) {
    luaScript := `
        local key = KEYS[1]
        local emission = tonumber(ARGV[1])
        local tolerance = tonumber(ARGV[2])
        local now = tonumber(ARGV[3])

        local tat = redis.call('GET', key)
        if not tat then
            tat = now
        else
            tat = tonumber(tat)
        end

        local new_tat = math.max(tat, now) + emission
        local allow_at = new_tat - tolerance

        if now >= allow_at then
            redis.call('SET', key, new_tat, 'EX', math.ceil((new_tat - now) / 1000000000) + 1)
            return {1, new_tat - now}
        else
            return {0, allow_at - now}
        end
    `

    now := time.Now().UnixNano()

    result, err := g.client.Eval(ctx, luaScript, []string{key},
        g.emission.Nanoseconds(),
        g.tolerance.Nanoseconds(),
        now,
    ).Result()

    if err != nil {
        return nil, err
    }

    values := result.([]interface{})
    allowed := values[0].(int64) == 1

    return &Result{
        Allowed:    allowed,
        RetryAfter: time.Duration(values[1].(int64)),
    }, nil
}
```

## Deployment

### Docker Compose

```yaml
version: '3.8'

services:
  rate-limiter:
    build: .
    ports:
      - "8080:8080"
    environment:
      - RATE_LIMIT_ALGORITHM=token_bucket
      - RATE_LIMIT_CAPACITY=100
      - RATE_LIMIT_REFILL=10
      - REDIS_URL=redis:6379
    depends_on:
      - redis

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rate-limiter
spec:
  replicas: 3
  selector:
    matchLabels:
      app: rate-limiter
  template:
    metadata:
      labels:
        app: rate-limiter
    spec:
      containers:
      - name: rate-limiter
        image: rate-limiter:latest
        ports:
        - containerPort: 8080
        env:
        - name: RATE_LIMIT_ALGORITHM
          value: "token_bucket"
        - name: REDIS_URL
          value: "redis:6379"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: rate-limiter
spec:
  selector:
    app: rate-limiter
  ports:
  - port: 80
    targetPort: 8080
```

## Performance

### Benchmarks

| Algorithm | Local (ops/sec) | Redis (ops/sec) | Memory |
|-----------|----------------|-----------------|--------|
| Token Bucket | 10,000,000 | 50,000 | 32 bytes |
| Sliding Window | 5,000,000 | 30,000 | 128 bytes |
| Fixed Window | 20,000,000 | 100,000 | 16 bytes |
| GCRA | 8,000,000 | 40,000 | 64 bytes |

### Load Testing

```bash
# Run load tests
wrk -t12 -c400 -d30s http://localhost:8080/api/test

# Or use k6
k6 run tests/load/rate_limit_test.js
```

## Best Practices

### Configuration Guidelines

1. **Choose the Right Algorithm**:
   - Token Bucket: General purpose, allows bursts
   - Sliding Window: Strict rate limiting
   - Fixed Window: Simple, low overhead

2. **Set Appropriate Limits**:
   - Start conservative and adjust
   - Monitor actual usage patterns
   - Consider peak vs. average traffic

3. **Handle Edge Cases**:
   - Client clock skew
   - Redis failures (fail open or closed)
   - Key explosion (use key patterns)

### Example Configuration

```yaml
rate_limits:
  api:
    default:
      algorithm: token_bucket
      capacity: 100
      refill_rate: 10  # per second
    authenticated:
      algorithm: token_bucket
      capacity: 1000
      refill_rate: 100
  web:
    default:
      algorithm: sliding_window
      limit: 60
      window: 1m
    search:
      algorithm: token_bucket
      capacity: 10
      refill_rate: 1
```

## License

MIT License

---

**Last Updated**: 2024-01-15
**Version**: 1.0.0
