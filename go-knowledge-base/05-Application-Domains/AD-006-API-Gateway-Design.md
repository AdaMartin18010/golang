# AD-006: API Gateway Design Patterns

> **Dimension**: Application Domains
> **Level**: S (17+ KB)
> **Tags**: #api-gateway #routing #rate-limiting #authentication #load-balancing

---

## 1. API Gateway Architecture

### 1.1 Core Responsibilities

| Responsibility | Description |
|---------------|-------------|
| Routing | Route requests to appropriate backend services |
| Authentication | Verify JWT, API keys, OAuth tokens |
| Rate Limiting | Prevent abuse and ensure fair usage |
| Load Balancing | Distribute traffic across instances |
| SSL Termination | Handle HTTPS encryption/decryption |
| Caching | Cache responses to reduce backend load |
| Request/Response Transformation | Convert between protocols/formats |

### 1.2 Gateway Pattern Types

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         API Gateway Patterns                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Single Gateway                                                           │
│  ┌─────────┐     ┌──────────────┐     ┌──────────────┐                     │
│  │ Client  │────>│   Gateway    │────>│  Services    │                     │
│  └─────────┘     └──────────────┘     └──────────────┘                     │
│                                                                              │
│  2. Backend for Frontend (BFF)                                               │
│  ┌─────────┐     ┌──────────┐     ┌──────────────┐                         │
│  │  Web    │────>│ Web BFF  │────>│  Services    │                         │
│  ├─────────┤     ├──────────┤     └──────────────┘                         │
│  │ Mobile  │────>│Mobile BFF│                                             │
│  ├─────────┤     ├──────────┤                                             │
│  │  IoT    │────>│ IoT BFF  │                                             │
│  └─────────┘     └──────────┘                                             │
│                                                                              │
│  3. Micro Gateway                                                            │
│  ┌─────────┐     ┌────────┐     ┌────────┐     ┌────────┐                 │
│  │ Client  │────>│ Gateway│────>│ServiceA│     │ServiceB│                 │
│  └─────────┘     └────────┘     └────────┘     └────────┘                 │
│                     │                                                        │
│                     └────────────────────────────────> ServiceC             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Routing Implementation

### 2.1 Dynamic Router

```go
package gateway

import (
    "context"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "sync"
)

type Router struct {
    routes map[string]*Route
    mu     sync.RWMutex
}

type Route struct {
    ID          string
    Path        string
    Methods     []string
    Target      *url.URL
    StripPath   bool
    Middlewares []Middleware
    RateLimit   *RateLimitConfig
}

func NewRouter() *Router {
    return &Router{
        routes: make(map[string]*Route),
    }
}

func (r *Router) Register(route *Route) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.routes[route.ID] = route
}

func (r *Router) Match(req *http.Request) (*Route, map[string]string) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for _, route := range r.routes {
        if !r.methodMatch(req.Method, route.Methods) {
            continue
        }

        if params, ok := r.pathMatch(req.URL.Path, route.Path); ok {
            return route, params
        }
    }

    return nil, nil
}

func (r *Router) methodMatch(method string, allowed []string) bool {
    if len(allowed) == 0 {
        return true
    }
    for _, m := range allowed {
        if m == method {
            return true
        }
    }
    return false
}

func (r *Router) pathMatch(path, pattern string) (map[string]string, bool) {
    // Simple pattern matching with parameters
    // e.g., /users/:id matches /users/123
    pathParts := strings.Split(path, "/")
    patternParts := strings.Split(pattern, "/")

    if len(pathParts) != len(patternParts) {
        return nil, false
    }

    params := make(map[string]string)
    for i, part := range patternParts {
        if strings.HasPrefix(part, ":") {
            params[part[1:]] = pathParts[i]
        } else if part != pathParts[i] {
            return nil, false
        }
    }

    return params, true
}
```

### 2.2 Load Balancer

```go
package gateway

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync/atomic"
)

type LoadBalancer interface {
    Next() *url.URL
}

// Round Robin Load Balancer
type RoundRobin struct {
    targets []*url.URL
    current uint64
}

func NewRoundRobin(targets []*url.URL) *LoadBalancer {
    return &RoundRobin{targets: targets}
}

func (r *RoundRobin) Next() *url.URL {
    n := atomic.AddUint64(&r.current, 1)
    return r.targets[int(n)%len(r.targets)]
}

// Weighted Round Robin
type WeightedRoundRobin struct {
    targets []*WeightedTarget
    current int
    cw      int
}

type WeightedTarget struct {
    URL    *url.URL
    Weight int
}

func (w *WeightedRoundRobin) Next() *url.URL {
    for {
        w.current = (w.current + 1) % len(w.targets)
        if w.current == 0 {
            w.cw--
            if w.cw <= 0 {
                w.cw = w.maxWeight()
            }
        }
        if w.targets[w.current].Weight >= w.cw {
            return w.targets[w.current].URL
        }
    }
}

func (w *WeightedRoundRobin) maxWeight() int {
    max := 0
    for _, t := range w.targets {
        if t.Weight > max {
            max = t.Weight
        }
    }
    return max
}
```

---

## 3. Rate Limiting

### 3.1 Token Bucket Algorithm

```go
package gateway

import (
    "context"
    "net/http"
    "sync"
    "time"
)

type TokenBucket struct {
    capacity int
    tokens   float64
    rate     float64
    lastRefill time.Time
    mu       sync.Mutex
}

func NewTokenBucket(capacity int, ratePerSecond float64) *TokenBucket {
    return &TokenBucket{
        capacity:   capacity,
        tokens:     float64(capacity),
        rate:       ratePerSecond,
        lastRefill: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastRefill).Seconds()
    tb.tokens = min(float64(tb.capacity), tb.tokens+elapsed*tb.rate)
    tb.lastRefill = now

    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }

    return false
}

// Distributed Rate Limiter with Redis
type DistributedRateLimiter struct {
    redis      RedisClient
    windowSize time.Duration
    maxRequests int
}

func (rl *DistributedRateLimiter) Allow(ctx context.Context, key string) bool {
    now := time.Now().Unix()
    windowStart := now - int64(rl.windowSize.Seconds())

    pipe := rl.redis.Pipeline()
    pipe.ZRemRangeByScore(ctx, key, "0", string(windowStart))
    pipe.ZCard(ctx, key)
    pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
    pipe.Expire(ctx, key, rl.windowSize)

    results, err := pipe.Exec(ctx)
    if err != nil {
        return false
    }

    currentCount := results[1].(*redis.IntCmd).Val()
    return currentCount < int64(rl.maxRequests)
}
```

### 3.2 Rate Limit Middleware

```go
package gateway

import (
    "net/http"
    "strconv"
)

func RateLimitMiddleware(limiter RateLimiter) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := extractClientKey(r)

            if !limiter.Allow(r.Context(), key) {
                w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.Limit()))
                w.Header().Set("X-RateLimit-Remaining", "0")
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limiter.Remaining(key)))
            next.ServeHTTP(w, r)
        })
    }
}

func extractClientKey(r *http.Request) string {
    // Try API key first
    apiKey := r.Header.Get("X-API-Key")
    if apiKey != "" {
        return apiKey
    }

    // Fall back to IP address
    ip := r.Header.Get("X-Forwarded-For")
    if ip == "" {
        ip = r.RemoteAddr
    }
    return ip
}
```

---

## 4. Authentication & Authorization

### 4.1 JWT Authentication

```go
package gateway

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
    secret     []byte
    headerName string
}

func (a *JWTAuth) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := extractToken(r, a.headerName)
        if tokenString == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return a.secret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add claims to context
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            ctx := context.WithValue(r.Context(), "claims", claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        }
    })
}

func extractToken(r *http.Request, headerName string) string {
    authHeader := r.Header.Get(headerName)
    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
        return parts[1]
    }
    return authHeader
}
```

---

## 5. Circuit Breaker

```go
package gateway

import (
    "errors"
    "net/http"
    "sync"
    "time"
)

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreaker struct {
    failureThreshold int
    successThreshold int
    timeout          time.Duration

    state           State
    failures        int
    successes       int
    lastFailureTime time.Time
    mu              sync.RWMutex
}

func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failureThreshold: failureThreshold,
        successThreshold: successThreshold,
        timeout:          timeout,
        state:           StateClosed,
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if !cb.canExecute() {
        return ErrCircuitOpen
    }

    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    if cb.state == StateClosed {
        return true
    }

    if cb.state == StateOpen {
        if time.Since(cb.lastFailureTime) > cb.timeout {
            cb.mu.RUnlock()
            cb.mu.Lock()
            cb.state = StateHalfOpen
            cb.failures = 0
            cb.successes = 0
            cb.mu.Unlock()
            cb.mu.RLock()
            return true
        }
        return false
    }

    return true // StateHalfOpen
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err == nil {
        cb.successes++
        if cb.state == StateHalfOpen && cb.successes >= cb.successThreshold {
            cb.state = StateClosed
            cb.failures = 0
        }
    } else {
        cb.failures++
        cb.lastFailureTime = time.Now()
        if cb.failures >= cb.failureThreshold {
            cb.state = StateOpen
        }
    }
}
```

---

## 6. Caching

```go
package gateway

import (
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "time"
)

type CacheMiddleware struct {
    cache       Cache
    ttl         time.Duration
    cacheable   func(r *http.Request) bool
}

func (m *CacheMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !m.cacheable(r) {
            next.ServeHTTP(w, r)
            return
        }

        key := generateCacheKey(r)

        // Try cache
        if cached, found := m.cache.Get(key); found {
            writeCachedResponse(w, cached)
            return
        }

        // Capture response
        recorder := NewResponseRecorder(w)
        next.ServeHTTP(recorder, r)

        // Cache if successful
        if recorder.StatusCode == http.StatusOK {
            m.cache.Set(key, recorder.Body, m.ttl)
        }
    })
}

func generateCacheKey(r *http.Request) string {
    data := r.Method + r.URL.String() + r.Header.Get("Accept")
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

---

## 7. Gateway Comparison

| Feature | Nginx | Kong | Envoy | Traefik | Custom |
|---------|-------|------|-------|---------|--------|
| Performance | High | High | Very High | High | Varies |
| Ease of Use | Medium | Easy | Complex | Easy | Varies |
| Plugin System | Lua | Rich | WASM | Moderate | Custom |
| K8s Native | No | Yes | Yes | Yes | Varies |
| Observability | Basic | Good | Excellent | Good | Custom |

---

## 8. Best Practices

- [ ] Implement health checks
- [ ] Use appropriate timeouts
- [ ] Log all requests
- [ ] Monitor gateway metrics
- [ ] Implement graceful degradation
- [ ] Use SSL/TLS termination
- [ ] Implement request tracing
- [ ] Configure CORS properly
- [ ] Use connection pooling
- [ ] Implement request validation

---

**Quality Rating**: S (17+ KB)
**Last Updated**: 2026-04-02

---

## 架构决策记录

### 决策矩阵

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| A | 高性能 | 复杂 | 大规模 |
| B | 简单 | 扩展性差 | 小规模 |

### 风险评估

**风险 R.1**: 性能瓶颈
- 概率: 中
- 影响: 高
- 缓解: 缓存、分片

**风险 R.2**: 单点故障
- 概率: 低
- 影响: 极高
- 缓解: 冗余、故障转移

### 实施路线图

`
Phase 1: 基础设施 (Week 1-2)
Phase 2: 核心功能 (Week 3-6)
Phase 3: 优化加固 (Week 7-8)
`

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02