# 框架能力完整集成示例

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 概述

本文档提供一个完整的示例，展示如何在一个 HTTP 服务中集成框架的所有核心能力。

---

## 🎯 完整示例

### 1. 初始化框架能力

```go
package main

import (
    "context"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/yourusername/golang/internal/interfaces/http/chi"
    chiMiddleware "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/pkg/control"
    "github.com/yourusername/golang/pkg/database"
    "github.com/yourusername/golang/pkg/observability/otlp"
    "github.com/yourusername/golang/pkg/sampling"
    "github.com/yourusername/golang/pkg/tracing"
)

func main() {
    ctx := context.Background()

    // 1. 初始化可观测性（OTLP）
    sampler, _ := sampling.NewProbabilisticSampler(0.5)
    otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
        ServiceName:    "my-service",
        ServiceVersion: "v1.0.0",
        Endpoint:       "localhost:4317",
        Insecure:       true,
        Sampler:        sampler,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer otlpClient.Shutdown(ctx)

    // 2. 初始化追踪器
    tracer := tracing.NewTracer("my-service")

    // 3. 初始化数据库
    db, err := database.NewDatabase(database.Config{
        Driver:       database.DriverPostgreSQL,
        DSN:          "postgres://user:password@localhost/dbname?sslmode=disable",
        MaxOpenConns: 25,
        MaxIdleConns: 5,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 4. 初始化精细控制
    featureController := control.NewFeatureController()
    rateController := control.NewRateController()
    circuitController := control.NewCircuitController()

    // 注册功能
    featureController.Register("experimental-feature", "Experimental feature", true, nil)
    rateController.SetRateLimit("user-api", 100.0, time.Second)
    circuitController.RegisterCircuit("external-api", 10, 5, 30*time.Second)

    // 5. 创建 HTTP 路由器
    r := chi.NewRouter()

    // 6. 配置中间件（按顺序）
    setupMiddleware(r, sampler, featureController, rateController, circuitController)

    // 7. 配置路由
    setupRoutes(r, db, tracer)

    // 8. 启动服务器
    http.ListenAndServe(":8080", r)
}
```

### 2. 配置中间件

```go
func setupMiddleware(
    r *chi.Mux,
    sampler sampling.Sampler,
    featureController *control.FeatureController,
    rateController *control.RateController,
    circuitController *control.CircuitController,
) {
    // 基础中间件
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)

    // 采样中间件
    r.Use(chiMiddleware.SamplingMiddleware(chiMiddleware.SamplingConfig{
        Sampler:             sampler,
        SkipPaths:           []string{"/health", "/metrics"},
        AddSamplingDecision: true,
    }))

    // 追踪中间件
    r.Use(chiMiddleware.TracingMiddleware(chiMiddleware.TracingConfig{
        ServiceName:    "my-service",
        ServiceVersion: "v1.0.0",
        SkipPaths:      []string{"/health", "/metrics"},
        AddRequestID:   true,
    }))

    // 反射中间件
    r.Use(chiMiddleware.ReflectMiddleware(chiMiddleware.ReflectConfig{
        EnableMetadata:     true,
        EnableSelfDescribe: true,
        SkipPaths:          []string{"/health", "/metrics"},
    }))

    // 数据转换中间件
    r.Use(chiMiddleware.ConverterMiddleware(chiMiddleware.ConverterConfig{
        EnableRequestConversion:  true,
        EnableResponseConversion: true,
        DefaultResponseFormat:    "json",
    }))

    // 精细控制中间件
    r.Use(chiMiddleware.ControlMiddleware(chiMiddleware.ControlConfig{
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

    // 其他中间件
    r.Use(chiMiddleware.LoggingMiddleware)
    r.Use(chiMiddleware.RecoveryMiddleware)
    r.Use(chiMiddleware.TimeoutMiddleware(60 * time.Second))
    r.Use(chiMiddleware.CORSMiddleware)
}
```

### 3. 配置路由和处理器

```go
func setupRoutes(r *chi.Mux, db database.Database, tracer *tracing.Tracer) {
    // 健康检查
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // API 路由
    r.Route("/api/v1", func(r chi.Router) {
        // 用户相关路由
        r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
            // 开始追踪
            ctx, span := tracer.StartSpan(r.Context(), "get-users")
            defer span.End()

            // 检查采样决策
            if chiMiddleware.IsSampled(ctx) {
                // 只有被采样的请求才记录详细日志
                log.Debug("Detailed request info", ...)
            }

            // 使用数据转换
            data := chiMiddleware.GetRequestData(ctx)
            if data != nil {
                // 使用转换后的数据
            }

            // 使用反射检查器
            inspector := chiMiddleware.GetInspector(ctx)
            if inspector != nil {
                metadata := inspector.InspectType(userStruct)
                // 使用元数据...
            }

            // 执行数据库操作
            rows, err := db.Query(ctx, "SELECT * FROM users")
            if err != nil {
                tracer.LocateError(ctx, err, map[string]interface{}{
                    "endpoint": "/api/v1/users",
                })
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            defer rows.Close()

            // 处理结果...
        })
    })
}
```

---

## 📊 能力集成矩阵

| 框架能力 | 中间件集成 | 直接使用 | 状态 |
|---------|-----------|---------|------|
| 采样机制 | ✅ 采样中间件 | ✅ | ✅ |
| 数据转换 | ✅ 数据转换中间件 | ✅ | ✅ |
| 精细控制 | ✅ 精细控制中间件 | ✅ | ✅ |
| 反射/自解释 | ✅ 反射中间件 | ✅ | ✅ |
| 追踪定位 | ✅ 追踪中间件 | ✅ | ✅ |
| 数据库抽象 | - | ✅ | ✅ |
| OTLP 集成 | - | ✅ | ✅ |
| eBPF 支持 | - | ✅ | ✅ |

---

## 🎉 总结

框架的所有核心能力都已经：

1. ✅ **实现完成** - 所有核心能力都已实现
2. ✅ **中间件集成** - 关键能力已集成到 HTTP 中间件
3. ✅ **文档完善** - 提供了完整的使用文档和示例
4. ✅ **测试覆盖** - 核心能力都有单元测试

框架现在是一个功能完整、可观测、可控制、可自解释的现代化 Go 框架！

---

**最后更新**: 2025-01-XX
