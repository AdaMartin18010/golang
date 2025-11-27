# 框架使用完整指南

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 概述

本文档提供框架使用的完整指南，从快速开始到高级用法，帮助用户快速上手并充分利用框架的各种能力。

---

## 🚀 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行示例

```bash
# 完整示例
cd examples/framework-usage/complete
go run main.go

# Wire 依赖注入示例
cd examples/framework-usage/wire-example
go generate
go run .
```

---

## 📚 使用场景

### 场景 1: 基础 HTTP 服务

创建一个基础的 HTTP 服务，使用框架的核心能力：

```go
package main

import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/pkg/sampling"
)

func main() {
    r := chi.NewRouter()

    // 配置采样中间件
    sampler, _ := sampling.NewProbabilisticSampler(0.5)
    r.Use(middleware.SamplingMiddleware(middleware.SamplingConfig{
        Sampler: sampler,
    }))

    // 配置追踪中间件
    r.Use(middleware.TracingMiddleware(middleware.TracingConfig{
        ServiceName: "my-service",
    }))

    // 配置路由...
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    http.ListenAndServe(":8080", r)
}
```

### 场景 2: 使用数据库

```go
import "github.com/yourusername/golang/pkg/database"

// 创建数据库连接
db, err := database.NewDatabase(database.Config{
    Driver: database.DriverPostgreSQL,
    DSN:    "postgres://user:password@localhost/dbname?sslmode=disable",
})
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// 使用数据库
rows, err := db.Query(ctx, "SELECT * FROM users")
```

### 场景 3: 使用 Ent ORM

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/database/ent"
    "github.com/yourusername/golang/internal/infrastructure/database/ent/repository"
)

// 创建 Ent 客户端
client, err := ent.NewClientFromConfig(ctx, "localhost", "5432", "user", "password", "dbname", "disable")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 创建仓储
userRepo := repository.NewUserRepository(client)

// 使用仓储
user := &user.User{Email: "test@example.com", Name: "Test"}
err = userRepo.Create(ctx, user)
```

### 场景 4: 使用精细控制

```go
import "github.com/yourusername/golang/pkg/control"

// 创建控制器
featureController := control.NewFeatureController()
rateController := control.NewRateController()
circuitController := control.NewCircuitController()

// 注册功能
featureController.Register("feature-a", "Feature A", true, nil)

// 设置速率限制
rateController.SetRateLimit("api", 100.0, time.Second)

// 注册熔断器
circuitController.RegisterCircuit("external-api", 10, 5, 30*time.Second)

// 在中间件中使用
r.Use(middleware.ControlMiddleware(middleware.ControlConfig{
    FeatureController: featureController,
    RateController:    rateController,
    CircuitController: circuitController,
}))
```

---

## 🎯 最佳实践

### 1. 中间件顺序

建议的中间件执行顺序：

1. RequestID
2. RealIP
3. 采样中间件
4. 追踪中间件
5. 反射中间件
6. 数据转换中间件
7. 精细控制中间件
8. 限流中间件
9. 认证中间件
10. 日志中间件
11. 恢复中间件
12. 超时中间件
13. CORS 中间件

### 2. 错误处理

```go
import "github.com/yourusername/golang/pkg/errors"

// 创建错误
err := errors.NewNotFoundError("user", "123")

// 添加详细信息
err = err.WithDetails("field", "email").WithTraceID("trace-123")

// 检查错误类型
if appErr, ok := err.(*errors.AppError); ok {
    log.Printf("Error code: %s, HTTP status: %d", appErr.Code, appErr.HTTPStatusCode())
}
```

### 3. 日志记录

```go
import "github.com/yourusername/golang/internal/framework/logger"

// 创建日志器
log := logger.NewLogger(logger.LevelInfo)

// 使用日志器
log.Info("Request processed", "path", r.URL.Path, "status", 200)
```

---

## 📚 相关文档

- [快速开始指南](05-快速开始指南.md)
- [最佳实践指南](06-最佳实践指南.md)
- [核心能力使用示例](08-核心能力使用示例.md)
- [框架能力完整集成示例](13-框架能力完整集成示例.md)

---

**最后更新**: 2025-01-XX
