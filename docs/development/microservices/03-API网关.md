# API网关

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.23+

---
## 📋 目录

- [API网关](#api网关)
  - [📚 理论分析](#理论分析)
  - [💻 代码示例](#代码示例)
  - [🎯 最佳实践](#最佳实践)
  - [🔍 常见问题](#常见问题)
  - [📚 扩展阅读](#扩展阅读)
  - [✅ 落地检查清单](#落地检查清单)
  - [🧪 本地快速验证](#本地快速验证)

---

## 📚 理论分析

### API网关作用

API网关是微服务架构中的统一入口，负责：

- **统一接入**: 为所有微服务提供统一的访问入口
- **路由转发**: 根据请求路径将流量转发到对应的后端服务
- **协议转换**: 支持HTTP、gRPC、WebSocket等多种协议
- **服务治理**: 提供限流、熔断、监控等服务治理功能

### 核心功能

| 功能类别 | 具体功能 | 说明 |
|----------|----------|------|
| 路由转发 | 路径匹配 | 根据URL路径匹配后端服务 |
| 负载均衡 | 轮询、加权轮询 | 在多个服务实例间分发请求 |
| 认证授权 | JWT、OAuth2 | 统一身份认证和权限控制 |
| 限流熔断 | 令牌桶、滑动窗口 | 保护后端服务不被过载 |
| 监控日志 | 请求追踪、性能监控 | 全链路监控和日志收集 |
| 协议转换 | HTTP/gRPC/WebSocket | 支持多种通信协议 |

### Go语言实现方案

- **Kong**: 基于Nginx的API网关
- **Traefik**: 云原生API网关
- **Ambassador**: Kubernetes原生API网关
- **自研网关**: 基于Gin/Echo等框架实现

## 💻 代码示例

### 基础网关实现

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

type Gateway struct {
    routes map[string]*Route
    client *http.Client
}

type Route struct {
    Name        string
    Path        string
    Methods     []string
    BackendURLs []string
    Middleware  []gin.HandlerFunc
}

type ServiceRegistry struct {
    services map[string][]string
}

func NewGateway() *Gateway {
    return &Gateway{
        routes: make(map[string]*Route),
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (g *Gateway) AddRoute(route *Route) {
    g.routes[route.Path] = route
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 查找匹配的路由
    route := g.findRoute(r.URL.Path)
    if route == nil {
        http.NotFound(w, r)
        return
    }

    // 检查HTTP方法
    if !g.isMethodAllowed(route, r.Method) {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // 选择后端服务
    backendURL := g.selectBackend(route)
    if backendURL == "" {
        http.Error(w, "No Backend Available", http.StatusServiceUnavailable)
        return
    }

    // 转发请求
    g.proxyRequest(w, r, backendURL)
}

func (g *Gateway) findRoute(path string) *Route {
    for pattern, route := range g.routes {
        if strings.HasPrefix(path, pattern) {
            return route
        }
    }
    return nil
}

func (g *Gateway) isMethodAllowed(route *Route, method string) bool {
    for _, allowedMethod := range route.Methods {
        if allowedMethod == method || allowedMethod == "*" {
            return true
        }
    }
    return false
}

func (g *Gateway) selectBackend(route *Route) string {
    if len(route.BackendURLs) == 0 {
        return ""
    }

    // 简单的轮询负载均衡
    // 实际应用中可以使用更复杂的算法
    return route.BackendURLs[0]
}

func (g *Gateway) proxyRequest(w http.ResponseWriter, r *http.Request, backendURL string) {
    // 解析后端URL
    target, err := url.Parse(backendURL)
    if err != nil {
        http.Error(w, "Invalid Backend URL", http.StatusInternalServerError)
        return
    }

    // 创建反向代理
    proxy := httputil.NewSingleHostReverseProxy(target)

    // 修改请求
    r.URL.Host = target.Host
    r.URL.Scheme = target.Scheme
    r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
    r.Host = target.Host

    // 转发请求
    proxy.ServeHTTP(w, r)
}

// 使用示例
func main() {
    gateway := NewGateway()

    // 添加路由
    gateway.AddRoute(&Route{
        Name:        "user-service",
        Path:        "/api/users",
        Methods:     []string{"GET", "POST", "PUT", "DELETE"},
        BackendURLs: []string{"http://localhost:8081", "http://localhost:8082"},
    })

    gateway.AddRoute(&Route{
        Name:        "order-service",
        Path:        "/api/orders",
        Methods:     []string{"GET", "POST"},
        BackendURLs: []string{"http://localhost:8083"},
    })

    // 启动网关
    server := &http.Server{
        Addr:    ":8080",
        Handler: gateway,
    }

    log.Println("API网关启动在 :8080")
    log.Fatal(server.ListenAndServe())
}
```

### 路由与负载均衡

```go
package main

import (
    "math/rand"
    "sync"
    "time"
)

type LoadBalancer interface {
    Select(backends []string) string
}

type RoundRobinBalancer struct {
    current int
    Mutex   sync.Mutex
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
    return &RoundRobinBalancer{}
}

func (rr *RoundRobinBalancer) Select(backends []string) string {
    if len(backends) == 0 {
        return ""
    }

    rr.Mutex.Lock()
    defer rr.Mutex.Unlock()

    backend := backends[rr.current]
    rr.current = (rr.current + 1) % len(backends)

    return backend
}

type WeightedRoundRobinBalancer struct {
    weights map[string]int
    current map[string]int
    Mutex   sync.Mutex
}

func NewWeightedRoundRobinBalancer() *WeightedRoundRobinBalancer {
    return &WeightedRoundRobinBalancer{
        weights: make(map[string]int),
        current: make(map[string]int),
    }
}

func (wrr *WeightedRoundRobinBalancer) SetWeight(backend string, weight int) {
    wrr.Mutex.Lock()
    defer wrr.Mutex.Unlock()

    wrr.weights[backend] = weight
    wrr.current[backend] = 0
}

func (wrr *WeightedRoundRobinBalancer) Select(backends []string) string {
    if len(backends) == 0 {
        return ""
    }

    wrr.Mutex.Lock()
    defer wrr.Mutex.Unlock()

    // 找到权重最大的后端
    maxWeight := -1
    selectedBackend := backends[0]

    for _, backend := range backends {
        weight := wrr.weights[backend]
        if weight <= 0 {
            weight = 1 // 默认权重
        }

        wrr.current[backend] += weight

        if wrr.current[backend] > maxWeight {
            maxWeight = wrr.current[backend]
            selectedBackend = backend
        }
    }

    wrr.current[selectedBackend] -= maxWeight

    return selectedBackend
}

type RandomBalancer struct {
    rand *rand.Rand
}

func NewRandomBalancer() *RandomBalancer {
    return &RandomBalancer{
        rand: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (rb *RandomBalancer) Select(backends []string) string {
    if len(backends) == 0 {
        return ""
    }

    return backends[rb.rand.Intn(len(backends))]
}

// 增强的网关实现
type EnhancedGateway struct {
    routes       map[string]*Route
    balancers    map[string]LoadBalancer
    client       *http.Client
}

func NewEnhancedGateway() *EnhancedGateway {
    return &EnhancedGateway{
        routes:    make(map[string]*Route),
        balancers: make(map[string]LoadBalancer),
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (eg *EnhancedGateway) AddRoute(route *Route, balancer LoadBalancer) {
    eg.routes[route.Path] = route
    eg.balancers[route.Name] = balancer
}

func (eg *EnhancedGateway) selectBackend(route *Route) string {
    balancer, exists := eg.balancers[route.Name]
    if !exists {
        // 默认使用轮询
        balancer = NewRoundRobinBalancer()
        eg.balancers[route.Name] = balancer
    }

    return balancer.Select(route.BackendURLs)
}

// 使用示例
func main() {
    gateway := NewEnhancedGateway()

    // 用户服务使用加权轮询
    userBalancer := NewWeightedRoundRobinBalancer()
    userBalancer.SetWeight("http://localhost:8081", 3)
    userBalancer.SetWeight("http://localhost:8082", 1)

    gateway.AddRoute(&Route{
        Name:        "user-service",
        Path:        "/api/users",
        Methods:     []string{"GET", "POST", "PUT", "DELETE"},
        BackendURLs: []string{"http://localhost:8081", "http://localhost:8082"},
    }, userBalancer)

    // 订单服务使用随机负载均衡
    orderBalancer := NewRandomBalancer()
    gateway.AddRoute(&Route{
        Name:        "order-service",
        Path:        "/api/orders",
        Methods:     []string{"GET", "POST"},
        BackendURLs: []string{"http://localhost:8083", "http://localhost:8084"},
    }, orderBalancer)

    // 启动网关
    server := &http.Server{
        Addr:    ":8080",
        Handler: gateway,
    }

    log.Println("增强API网关启动在 :8080")
    log.Fatal(server.ListenAndServe())
}
```

### 认证与授权

```go
package main

import (
    "crypto/rsa"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
    publicKey *rsa.PublicKey
}

type Claims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

func NewAuthMiddleware(publicKeyPath string) (*AuthMiddleware, error) {
    keyBytes, err := ioutil.ReadFile(publicKeyPath)
    if err != nil {
        return nil, err
    }

    publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
    if err != nil {
        return nil, err
    }

    return &AuthMiddleware{publicKey: publicKey}, nil
}

func (am *AuthMiddleware) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return am.publicKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

func (am *AuthMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 获取Authorization头
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            // 检查Bearer前缀
            parts := strings.SplitN(authHeader, " ", 2)
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
                return
            }

            // 验证token
            claims, err := am.ValidateToken(parts[1])
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            // 将用户信息添加到请求上下文
            ctx := Context.WithValue(r.Context(), "user", claims)
            r = r.WithContext(ctx)

            next.ServeHTTP(w, r)
        })
    }
}

func (am *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := r.Context().Value("user").(*Claims)
            if !ok {
                http.Error(w, "User not authenticated", http.StatusUnauthorized)
                return
            }

            // 检查用户角色
            hasRole := false
            for _, userRole := range user.Roles {
                if userRole == role {
                    hasRole = true
                    break
                }
            }

            if !hasRole {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// 使用示例
func main() {
    // 创建认证中间件
    authMiddleware, err := NewAuthMiddleware("public.pem")
    if err != nil {
        log.Fatal(err)
    }

    // 创建路由
    mux := http.NewServeMux()

    // 公开路由
    mux.HandleFunc("/api/public", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Public endpoint"))
    })

    // 需要认证的路由
    protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := r.Context().Value("user").(*Claims)
        w.Write([]byte(fmt.Sprintf("Hello, %s!", user.Username)))
    })

    mux.Handle("/api/protected", authMiddleware.Middleware()(protectedHandler))

    // 需要管理员角色的路由
    adminHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Admin only endpoint"))
    })

    adminRoute := authMiddleware.Middleware()(adminHandler)
    adminRoute = authMiddleware.RequireRole("admin")(adminRoute)
    mux.Handle("/api/admin", adminRoute)

    // 启动服务器
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    log.Println("认证网关启动在 :8080")
    log.Fatal(server.ListenAndServe())
}
```

### 限流与熔断

```go
package main

import (
    "Context"
    "fmt"
    "sync"
    "time"
)

// 令牌桶限流器
type TokenBucket struct {
    capacity     int
    tokens       int
    refillRate   int
    lastRefill   time.Time
    Mutex        sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
    return &TokenBucket{
        capacity:   capacity,
        tokens:     capacity,
        refillRate: refillRate,
        lastRefill: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.Mutex.Lock()
    defer tb.Mutex.Unlock()

    // 补充令牌
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill)
    tokensToAdd := int(elapsed.Seconds()) * tb.refillRate

    if tokensToAdd > 0 {
        tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
        tb.lastRefill = now
    }

    // 检查是否有可用令牌
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }

    return false
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// 滑动窗口限流器
type SlidingWindow struct {
    windowSize time.Duration
    requests   []time.Time
    Mutex      sync.Mutex
    maxRequests int
}

func NewSlidingWindow(windowSize time.Duration, maxRequests int) *SlidingWindow {
    return &SlidingWindow{
        windowSize:   windowSize,
        maxRequests:  maxRequests,
        requests:     make([]time.Time, 0),
    }
}

func (sw *SlidingWindow) Allow() bool {
    sw.Mutex.Lock()
    defer sw.Mutex.Unlock()

    now := time.Now()
    cutoff := now.Add(-sw.windowSize)

    // 移除过期的请求
    var validRequests []time.Time
    for _, reqTime := range sw.requests {
        if reqTime.After(cutoff) {
            validRequests = append(validRequests, reqTime)
        }
    }
    sw.requests = validRequests

    // 检查是否超过限制
    if len(sw.requests) >= sw.maxRequests {
        return false
    }

    // 记录当前请求
    sw.requests = append(sw.requests, now)
    return true
}

// 熔断器
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
    state       string // "closed", "open", "half-open"
    Mutex       sync.Mutex
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       "closed",
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.Mutex.Lock()
    defer cb.Mutex.Unlock()

    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    }

    err := fn()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }

        return err
    }

    // 成功调用，重置状态
    cb.failures = 0
    cb.state = "closed"

    return nil
}

// 限流中间件
type RateLimitMiddleware struct {
    limiter *TokenBucket
}

func NewRateLimitMiddleware(limiter *TokenBucket) *RateLimitMiddleware {
    return &RateLimitMiddleware{limiter: limiter}
}

func (rlm *RateLimitMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !rlm.limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// 熔断中间件
type CircuitBreakerMiddleware struct {
    breaker *CircuitBreaker
}

func NewCircuitBreakerMiddleware(breaker *CircuitBreaker) *CircuitBreakerMiddleware {
    return &CircuitBreakerMiddleware{breaker: breaker}
}

func (cbm *CircuitBreakerMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            err := cbm.breaker.Call(func() error {
                // 创建一个响应写入器来捕获状态码
                rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
                next.ServeHTTP(rw, r)

                if rw.statusCode >= 500 {
                    return fmt.Errorf("server error: %d", rw.statusCode)
                }

                return nil
            })

            if err != nil {
                http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
            }
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// 使用示例
func main() {
    // 创建限流器
    rateLimiter := NewTokenBucket(100, 10) // 容量100，每秒补充10个令牌
    rateLimitMiddleware := NewRateLimitMiddleware(rateLimiter)

    // 创建熔断器
    circuitBreaker := NewCircuitBreaker(5, 30*time.Second) // 5次失败后熔断30秒
    circuitBreakerMiddleware := NewCircuitBreakerMiddleware(circuitBreaker)

    // 创建路由
    mux := http.NewServeMux()
    mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Test endpoint"))
    })

    // 应用中间件
    handler := rateLimitMiddleware.Middleware()(mux)
    handler = circuitBreakerMiddleware.Middleware()(handler)

    // 启动服务器
    server := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }

    log.Println("限流熔断网关启动在 :8080")
    log.Fatal(server.ListenAndServe())
}
```

## 🎯 最佳实践

1. **高可用设计**: 网关本身应该是高可用的，避免单点故障
2. **性能优化**: 使用连接池、缓存等技术提升性能
3. **监控告警**: 全面监控网关的请求量、响应时间、错误率等指标
4. **配置管理**: 支持动态配置更新，无需重启服务
5. **安全防护**: 实现完整的认证授权、限流熔断等安全机制

## 🔍 常见问题

1. **性能瓶颈**: 网关可能成为性能瓶颈，需要优化
2. **配置复杂**: 路由配置可能变得复杂，需要良好的管理工具
3. **故障传播**: 网关故障会影响所有服务，需要高可用设计
4. **版本兼容**: 后端服务升级时的兼容性处理

## 📚 扩展阅读

- [Kong官方文档](https://docs.konghq.com/)
- [Traefik官方文档](https://doc.traefik.io/traefik/)
- [Ambassador官方文档](https://www.getambassador.io/docs/)
- [API网关设计模式](https://microservices.io/patterns/apigateway.html)

## ✅ 落地检查清单

- 高可用：多副本/多可用区部署；无状态/会话粘滞策略明确
- 路由管理：动态下发与热更新；变更审计与回滚；灰度/金丝雀规则
- 认证授权：OIDC/JWT 验签与过期处理；权限模型与细粒度策略
- 治理策略：超时、重试（含抖动）、熔断、降级、限流全链覆盖
- 观测性：RED 指标、错误分级、分布式追踪（traceparent）透传
- 安全：mTLS、TLS 版本与Cipher 套件、WAF/防护规则、速率限制
- 性能：连接池、Keep-Alive、压缩与缓存；基准压测与容量规划

## 🧪 本地快速验证
