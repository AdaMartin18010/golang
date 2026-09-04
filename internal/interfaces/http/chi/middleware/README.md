# HTTP 中间件

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26.2

---

## 📋 目录

- [HTTP 中间件](#http-中间件)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 认证授权中间件](#2-认证授权中间件)
    - [2.1 功能特性](#21-功能特性)
    - [2.2 配置](#22-配置)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 角色权限控制](#32-角色权限控制)
    - [3.3 在Handler中使用用户信息](#33-在handler中使用用户信息)
    - [3.4 可选认证](#34-可选认证)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)
  - [6. 请求追踪中间件](#6-请求追踪中间件)
    - [6.1 功能特性](#61-功能特性)
    - [6.2 使用示例](#62-使用示例)
  - [7. 性能监控中间件](#7-性能监控中间件)
    - [7.1 功能特性](#71-功能特性)
    - [7.2 使用示例](#72-使用示例)
  - [8. 恢复中间件](#8-恢复中间件)
    - [8.1 功能特性](#81-功能特性)
    - [8.2 使用示例](#82-使用示例)
  - [9. CORS中间件](#9-cors中间件)
    - [9.1 功能特性](#91-功能特性)
    - [9.2 使用示例](#92-使用示例)
  - [10. 限流中间件](#10-限流中间件)
    - [10.1 功能特性](#101-功能特性)
    - [10.2 限流算法](#102-限流算法)
    - [10.3 使用示例](#103-使用示例)
    - [10.4 Redis 分布式限流](#104-redis-分布式限流)

---

## 1. 概述

HTTP 中间件提供了各种 HTTP 请求处理中间件：

- ✅ **认证授权中间件**: JWT Token 认证和角色权限控制
- ✅ **限流中间件**: 请求限流保护（支持多种算法：令牌桶、滑动窗口、漏桶，支持 Redis 分布式限流）
- ✅ **熔断器中间件**: 服务熔断保护（三种状态）
- ✅ **请求追踪中间件**: 请求链路追踪（基于OpenTelemetry）
- ✅ **性能监控中间件**: 请求性能监控和指标收集
- ✅ **恢复中间件**: Panic恢复和错误处理
- ✅ **CORS中间件**: 跨域资源共享支持

---

## 2. 认证授权中间件

### 2.1 功能特性

- JWT Token 验证
- 用户信息注入到 Context
- 角色权限验证
- 可选认证支持
- 路径跳过支持

### 2.2 配置

```go
type AuthConfig struct {
    JWT          *jwt.JWT
    SkipPaths    []string // 跳过认证的路径
    OptionalAuth bool     // 是否可选认证
}
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/pkg/auth/jwt"
)

// 创建JWT管理器
jwtManager, _ := jwt.NewJWT(jwt.Config{
    SecretKey:      "your-secret-key",
    SigningMethod:  "HS256",
    AccessTokenTTL: 15 * time.Minute,
})

// 创建路由
r := chi.NewRouter()

// 添加认证中间件
r.Use(middleware.AuthMiddleware(middleware.AuthConfig{
    JWT:       jwtManager,
    SkipPaths: []string{"/public", "/health"},
}))

// 受保护的路由
r.Get("/users", getUserHandler)
r.Post("/users", createUserHandler)
```

### 3.2 角色权限控制

```go
// 要求admin角色
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireRole("admin"))
    r.Delete("/users/{id}", deleteUserHandler)
})

// 要求任一角色
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAnyRole("admin", "moderator"))
    r.Put("/users/{id}", updateUserHandler)
})

// 要求所有角色
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAllRoles("admin", "superuser"))
    r.Delete("/users/{id}", deleteUserHandler)
})
```

### 3.3 在Handler中使用用户信息

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    // 从context获取用户信息
    userID := middleware.GetUserID(r.Context())
    username := middleware.GetUsername(r.Context())
    roles := middleware.GetRoles(r.Context())

    // 使用用户信息
    user, err := userService.GetUser(r.Context(), userID)
    // ...
}
```

### 3.4 可选认证

```go
r.Use(middleware.AuthMiddleware(middleware.AuthConfig{
    JWT:          jwtManager,
    OptionalAuth: true, // 可选认证，不强制
}))
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用中间件**: 在路由级别使用认证中间件
2. **跳过公开路径**: 为公开API设置跳过路径
3. **角色验证**: 使用角色中间件保护敏感操作
4. **从Context获取**: 从Context获取用户信息而不是从Token解析
5. **错误处理**: 使用统一的错误响应格式

### 4.2 DON'Ts ❌

1. **不要跳过所有路径**: 只跳过真正公开的路径
2. **不要硬编码角色**: 使用配置或数据库管理角色
3. **不要暴露敏感信息**: 错误消息不要暴露内部细节
4. **不要重复验证**: 在中间件中验证后不要在Handler中重复验证

---

## 5. 相关资源

- JWT 认证框架
- 统一错误处理框架
- 统一响应格式框架
- 框架拓展计划

---

## 6. 请求追踪中间件

### 6.1 功能特性

- OpenTelemetry 集成
- 自动提取和传播追踪上下文
- 请求和响应属性记录
- 追踪ID注入到响应头
- 路径跳过支持

### 6.2 使用示例

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
)

r := chi.NewRouter()
r.Use(middleware.TracingMiddleware(middleware.TracingConfig{
    TracerName:     "api-server",
    ServiceName:    "user-service",
    ServiceVersion: "v1.0.0",
    SkipPaths:      []string{"/health", "/metrics"},
    AddRequestID:   true,
    AddUserID:      true,
}))
```

---

## 7. 性能监控中间件

### 7.1 功能特性

- 请求计数统计
- 请求耗时统计
- 错误计数统计
- 活跃请求数统计
- 性能指标响应头

### 7.2 使用示例

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
)

metrics := middleware.NewMetrics()
r := chi.NewRouter()
r.Use(middleware.MetricsMiddleware(metrics))

// 指标查询端点
r.Get("/metrics", middleware.MetricsHandler(metrics))
```

---

## 8. 恢复中间件

### 8.1 功能特性

- Panic恢复
- 堆栈信息记录
- 错误响应格式化
- 可配置堆栈大小

### 8.2 使用示例

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/pkg/logger"
)

log := logger.NewLogger(slog.LevelError)
r.Use(middleware.RecoveryMiddleware(middleware.RecoveryConfig{
    Logger:    log,
    StackAll:  true,
    StackSize: 8192,
}))
```

---

## 9. CORS中间件

### 9.1 功能特性

- 可配置的允许源
- 支持预检请求
- 凭证支持
- 自定义请求头和响应头

### 9.2 使用示例

```go
r.Use(middleware.CORSMiddleware(middleware.CORSConfig{
    AllowedOrigins:   []string{"http://localhost:3000", "https://example.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           3600,
}))
```

---

## 10. 采样中间件

### 10.1 功能特性

- ✅ **多种采样策略**: 支持概率采样、速率限制采样、自适应采样等
- ✅ **路径跳过**: 支持跳过特定路径的采样（如 /health、/metrics）
- ✅ **采样决策传递**: 将采样决策添加到上下文，供后续中间件使用
- ✅ **响应头信息**: 可选在响应头中添加采样决策信息

### 10.2 使用示例

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/pkg/sampling"
)

// 创建采样器
sampler, _ := sampling.NewProbabilisticSampler(0.5)

// 配置采样中间件
router.Use(middleware.SamplingMiddleware(middleware.SamplingConfig{
    Sampler:             sampler,
    SkipPaths:           []string{"/health", "/metrics"},
    AddSamplingDecision: true,
}))

// 在处理器中使用采样决策
func MyHandler(w http.ResponseWriter, r *http.Request) {
    if middleware.IsSampled(r.Context()) {
        // 只有被采样的请求才记录详细日志
        log.Debug("Detailed request info", ...)
    }
}
```

---

## 11. 数据转换中间件

### 11.1 功能特性

- ✅ **请求数据转换**: 自动转换请求数据格式（JSON、Form 等）
- ✅ **响应数据转换**: 自动转换响应数据格式
- ✅ **多种格式支持**: 支持 JSON、XML、Form 等多种格式
- ✅ **上下文传递**: 将转换后的数据添加到上下文

### 11.2 使用示例

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
)

// 配置数据转换中间件
router.Use(middleware.ConverterMiddleware(middleware.ConverterConfig{
    EnableRequestConversion:  true,
    EnableResponseConversion: true,
    RequestFormats:           []string{"json", "form"},
    ResponseFormats:          []string{"json"},
    DefaultResponseFormat:    "json",
}))

// 在处理器中获取转换后的数据
func MyHandler(w http.ResponseWriter, r *http.Request) {
    data := middleware.GetRequestData(r.Context())
    if data != nil {
        // 使用转换后的数据
    }
}
```

---

## 12. 精细控制中间件

### 12.1 功能特性

- ✅ **功能开关**: 基于框架的功能控制器，动态启用/禁用功能
- ✅ **速率控制**: 基于框架的速率控制器，细粒度限流
- ✅ **熔断器**: 基于框架的熔断器控制器，自动熔断和恢复
- ✅ **路径配置**: 支持按路径配置不同的控制策略
- ✅ **上下文传递**: 将控制器添加到上下文，供后续使用

### 12.2 使用示例

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/pkg/control"
)

// 创建控制器
featureController := control.NewFeatureController()
rateController := control.NewRateController()
circuitController := control.NewCircuitController()

// 注册功能
featureController.Register("experimental-feature", "Experimental feature", true, nil)

// 注册速率限制
rateController.SetRateLimit("user-api", 100.0, time.Second)

// 注册熔断器
circuitController.RegisterCircuit("external-api", 10, 5, 30*time.Second)

// 配置精细控制中间件
router.Use(middleware.ControlMiddleware(middleware.ControlConfig{
    FeatureController: featureController,
    RateController:    rateController,
    CircuitController: circuitController,
    FeatureFlags: map[string]string{
        "/api/v1/experimental": "experimental-feature",
    },
    RateLimits: map[string]string{
        "/api/v1/users": "user-api",
    },
    CircuitBreakers: map[string]string{
        "/api/v1/external": "external-api",
    },
    SkipPaths: []string{"/health", "/metrics"},
}))

// 在处理器中使用功能开关
func MyHandler(w http.ResponseWriter, r *http.Request) {
    if middleware.GetFeatureFlag(r.Context(), "experimental-feature") {
        // 执行实验性功能
    }
}
```

---

## 13. 反射/自解释中间件

### 13.1 功能特性

- ✅ **元数据信息**: 在响应头中添加请求元数据信息
- ✅ **自描述功能**: 支持自描述功能
- ✅ **反射检查器**: 提供反射检查器供后续使用
- ✅ **路径配置**: 支持按路径配置元数据

### 13.2 使用示例

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
)

// 配置反射中间件
router.Use(middleware.ReflectMiddleware(middleware.ReflectConfig{
    EnableMetadata:     true,
    EnableSelfDescribe: true,
    SkipPaths:          []string{"/health", "/metrics"},
}))

// 在处理器中使用反射检查器
func MyHandler(w http.ResponseWriter, r *http.Request) {
    inspector := middleware.GetInspector(r.Context())
    if inspector != nil {
        metadata := inspector.InspectType(myStruct)
        // 使用元数据...
    }
}
```

---

## 14. 限流中间件

### 10.1 功能特性

- ✅ **多种限流算法**: 支持令牌桶、滑动窗口、漏桶三种算法
- ✅ **分布式限流**: 支持 Redis 分布式限流（适用于多实例部署）
- ✅ **灵活配置**: 可配置限流速率、突发容量、时间窗口
- ✅ **路径跳过**: 支持跳过特定路径的限流
- ✅ **自定义键生成**: 支持自定义限流键生成函数（默认基于 IP）

### 10.2 限流算法

#### 令牌桶算法 (Token Bucket)

- **特点**: 允许突发流量，适合需要处理突发请求的场景
- **适用场景**: API 限流、用户请求限流
- **优势**: 平滑处理突发流量

#### 滑动窗口算法 (Sliding Window)

- **特点**: 精确控制时间窗口内的请求数
- **适用场景**: 需要精确限流的场景
- **优势**: 更精确的限流控制

#### 漏桶算法 (Leaky Bucket)

- **特点**: 以固定速率处理请求，平滑输出
- **适用场景**: 需要平滑输出流量的场景
- **优势**: 输出速率恒定

### 10.3 使用示例

#### 基本使用（令牌桶算法）

```go
r.Use(middleware.RateLimitMiddleware(middleware.RateLimitConfig{
    RequestsPerSecond: 100,
    Burst:             200,
    Window:            time.Second,
    Algorithm:         middleware.AlgorithmTokenBucket,
}))
```

#### 滑动窗口算法

```go
r.Use(middleware.RateLimitMiddleware(middleware.RateLimitConfig{
    RequestsPerSecond: 100,
    Window:            time.Second,
    Algorithm:         middleware.AlgorithmSlidingWindow,
}))
```

#### 自定义限流键和跳过路径

```go
r.Use(middleware.RateLimitMiddleware(middleware.RateLimitConfig{
    RequestsPerSecond: 100,
    Burst:             200,
    Algorithm:         middleware.AlgorithmTokenBucket,
    KeyFunc: func(r *http.Request) string {
        userID := r.Header.Get("X-User-ID")
        if userID != "" {
            return "user:" + userID
        }
        return r.RemoteAddr
    },
    SkipPaths: []string{"/health", "/metrics"},
}))
```

### 10.4 Redis 分布式限流

当使用多个服务实例时，需要使用 Redis 进行分布式限流。需要实现 `RedisClient` 接口。

---

**更新日期**: 2025-11-11
