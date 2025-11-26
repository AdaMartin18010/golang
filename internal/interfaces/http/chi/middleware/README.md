# HTTP 中间件

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [HTTP 中间件](#http-中间件)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 认证授权中间件](#2-认证授权中间件)
  - [3. 使用示例](#3-使用示例)
  - [4. 最佳实践](#4-最佳实践)

---

## 1. 概述

HTTP 中间件提供了各种 HTTP 请求处理中间件：

- ✅ **认证授权中间件**: JWT Token 认证和角色权限控制
- ✅ **限流中间件**: 请求限流保护（令牌桶算法）
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

- [JWT 认证框架](../../../../pkg/auth/jwt/README.md)
- [统一错误处理框架](../../../../pkg/errors/README.md)
- [统一响应格式框架](../../../../pkg/http/response/README.md)
- [框架拓展计划](../../../../docs/00-框架拓展计划.md)

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

**更新日期**: 2025-11-11
