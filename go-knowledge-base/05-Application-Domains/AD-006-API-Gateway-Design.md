# AD-006: API 网关设计模式 (API Gateway Design Patterns)

> **维度**: Application Domains
> **级别**: S (17+ KB)
> **标签**: #api-gateway #edge-service #routing #security #rate-limiting
> **权威来源**: [API Gateway Pattern](https://microservices.io/patterns/apigateway.html), [Building Microservices](https://samnewman.io/books/building_microservices/)

---

## API 网关架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      API Gateway Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Client Apps                                                                 │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                                     │
│  │  Web    │  │ Mobile  │  │ 3rd Party│                                     │
│  │  App    │  │  Apps   │  │  API    │                                     │
│  └────┬────┘  └────┬────┘  └────┬────┘                                     │
│       │            │            │                                           │
│       └────────────┴────────────┘                                           │
│                    │                                                        │
│                    ▼                                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      API Gateway                                    │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │   │
│  │  │   Routing   │  │   Auth      │  │   Rate      │  │  Request   │ │   │
│  │  │   & LB      │  │   (JWT/OAuth)│  │   Limiting  │  │  Transform │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘ │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │   │
│  │  │  Circuit    │  │   Caching   │  │  Logging &  │  │   SSL/     │ │   │
│  │  │  Breaker    │  │             │  │  Metrics    │  │   TLS      │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘ │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                    │                                                        │
│       ┌────────────┼────────────┬────────────┐                             │
│       ▼            ▼            ▼            ▼                             │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐                        │
│  │ Service │  │ Service │  │ Service │  │ Service │                        │
│  │    A    │  │    B    │  │    C    │  │    D    │                        │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心功能

### 1. 路由与负载均衡

```go
package gateway

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "sync"
    "time"
)

// Route 路由配置
type Route struct {
    ID          string
    PathPrefix  string
    StripPrefix bool
    TargetURL   *url.URL
    Methods     []string
    Middlewares []Middleware
}

// Router 路由器
type Router struct {
    routes map[string]*Route
    mu     sync.RWMutex
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

func (r *Router) Match(req *http.Request) (*Route, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for _, route := range r.routes {
        if !r.methodMatch(req.Method, route.Methods) {
            continue
        }
        if strings.HasPrefix(req.URL.Path, route.PathPrefix) {
            return route, true
        }
    }
    return nil, false
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

// ReverseProxy 反向代理
func (r *Router) ReverseProxy(route *Route) *httputil.ReverseProxy {
    director := func(req *http.Request) {
        target := route.TargetURL

        req.URL.Scheme = target.Scheme
        req.URL.Host = target.Host

        if route.StripPrefix {
            req.URL.Path = strings.TrimPrefix(req.URL.Path, route.PathPrefix)
        }

        req.Header.Set("X-Forwarded-Host", req.Host)
        req.Header.Set("X-Forwarded-For", req.RemoteAddr)
    }

    return &httputil.ReverseProxy{
        Director: director,
        ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
        },
    }
}
```

### 2. 认证与授权

```go
package gateway

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware JWT 认证中间件
type AuthMiddleware struct {
    secretKey []byte
    validator TokenValidator
}

type TokenValidator interface {
    Validate(token string) (*Claims, error)
}

type Claims struct {
    UserID   string   `json:"sub"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

func (a *AuthMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
            return
        }

        claims, err := a.validator.Validate(parts[1])
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // 将用户信息注入上下文
        ctx := context.WithValue(r.Context(), "claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// RBACMiddleware 基于角色的访问控制
type RBACMiddleware struct {
    requiredRoles []string
}

func (r *RBACMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        claims, ok := req.Context().Value("claims").(*Claims)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if !hasAnyRole(claims.Roles, r.requiredRoles) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, req)
    })
}

func hasAnyRole(userRoles, requiredRoles []string) bool {
    roleSet := make(map[string]bool)
    for _, r := range userRoles {
        roleSet[r] = true
    }
    for _, r := range requiredRoles {
        if roleSet[r] {
            return true
        }
    }
    return false
}
```

### 3. 限流

```go
package gateway

import (
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

// RateLimiter 令牌桶限流器
type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    burst,
    }
}

func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
    rl.mu.RLock()
    limiter, exists := rl.limiters[key]
    rl.mu.RUnlock()

    if exists {
        return limiter
    }

    rl.mu.Lock()
    defer rl.mu.Unlock()

    // 双重检查
    if limiter, exists := rl.limiters[key]; exists {
        return limiter
    }

    limiter = rate.NewLimiter(rl.rate, rl.burst)
    rl.limiters[key] = limiter
    return limiter
}

func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 按 IP 限流
        key := r.RemoteAddr

        // 或按用户限流
        if claims, ok := r.Context().Value("claims").(*Claims); ok {
            key = claims.UserID
        }

        limiter := rl.GetLimiter(key)
        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## 网关实现对比

| 特性 | Nginx | Kong | Envoy | Spring Cloud Gateway | Traefik |
|------|-------|------|-------|---------------------|---------|
| 语言 | C/Lua | Lua | C++ | Java | Go |
| 配置 | 文件 | Admin API | xDS | Java Config | 动态 |
| 插件 | Lua | 丰富 | WASM | Spring | 中等 |
| K8s 集成 | 一般 | 好 | 原生 | 好 | 原生 |
| 性能 | 高 | 高 | 极高 | 中 | 高 |
| 服务发现 | 需配置 | 支持 | 原生 | 支持 | 原生 |

---

## 生产建议

| 关注点 | 建议 |
|--------|------|
| 高可用 | 多实例 + 负载均衡器 |
| 缓存 | 静态资源 CDN，API 响应缓存 |
| 超时 | 设置合理的上下游超时 |
| 熔断 | 集成断路器防止级联故障 |
| 监控 | 延迟、错误率、流量指标 |

---

## 参考文献

1. [API Gateway Pattern](https://microservices.io/patterns/apigateway.html)
2. [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman
3. [Kong Documentation](https://docs.konghq.com/)
