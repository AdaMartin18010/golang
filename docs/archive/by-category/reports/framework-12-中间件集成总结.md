# 中间件集成总结

> **版本**: v1.0
> **日期**: 2025-01-XX
> **状态**: ✅ 已完成

---

## 📋 概述

本文档总结框架核心能力在 HTTP 中间件中的集成情况，展示如何将框架的各种能力应用到 HTTP 请求处理流程中。

---

## ✅ 已集成的中间件

### 1. 采样中间件 ✅

**文件**: `internal/interfaces/http/chi/middleware/sampling.go`

**功能**:

- ✅ 集成框架的采样机制
- ✅ 支持多种采样策略（概率、速率限制、自适应等）
- ✅ 路径跳过功能
- ✅ 采样决策传递到上下文
- ✅ 响应头信息

**使用示例**:

```go
sampler, _ := sampling.NewProbabilisticSampler(0.5)
router.Use(middleware.SamplingMiddleware(middleware.SamplingConfig{
    Sampler:             sampler,
    SkipPaths:           []string{"/health", "/metrics"},
    AddSamplingDecision: true,
}))
```

---

### 2. 数据转换中间件 ✅

**文件**: `internal/interfaces/http/chi/middleware/converter.go`

**功能**:

- ✅ 集成框架的数据转换工具
- ✅ 自动转换请求数据格式（JSON、Form 等）
- ✅ 将转换后的数据添加到上下文

**使用示例**:

```go
router.Use(middleware.ConverterMiddleware(middleware.ConverterConfig{
    EnableRequestConversion:  true,
    EnableResponseConversion: true,
    DefaultResponseFormat:    "json",
}))
```

---

### 3. 精细控制中间件 ✅

**文件**: `internal/interfaces/http/chi/middleware/control.go`

**功能**:

- ✅ 集成框架的功能控制器
- ✅ 集成框架的速率控制器
- ✅ 集成框架的熔断器控制器
- ✅ 按路径配置不同的控制策略
- ✅ 上下文传递

**使用示例**:

```go
featureController := control.NewFeatureController()
rateController := control.NewRateController()
circuitController := control.NewCircuitController()

router.Use(middleware.ControlMiddleware(middleware.ControlConfig{
    FeatureController: featureController,
    RateController:    rateController,
    CircuitController: circuitController,
    FeatureFlags: map[string]string{
        "/api/v1/experimental": "experimental-feature",
    },
}))
```

---

### 4. 反射/自解释中间件 ✅

**文件**: `internal/interfaces/http/chi/middleware/reflect.go`

**功能**:

- ✅ 集成框架的反射能力
- ✅ 在响应头中添加元数据信息
- ✅ 提供反射检查器供后续使用
- ✅ 自描述功能

**使用示例**:

```go
router.Use(middleware.ReflectMiddleware(middleware.ReflectConfig{
    EnableMetadata:     true,
    EnableSelfDescribe: true,
    SkipPaths:          []string{"/health"},
}))
```

---

## 🔗 中间件与框架能力的对应关系

| 中间件 | 框架能力 | 集成状态 |
|--------|---------|---------|
| 采样中间件 | `pkg/sampling` | ✅ |
| 数据转换中间件 | `pkg/converter` | ✅ |
| 精细控制中间件 | `pkg/control` | ✅ |
| 反射/自解释中间件 | `pkg/reflect` | ✅ |
| 追踪中间件 | `pkg/tracing` + OpenTelemetry | ✅ |
| 限流中间件 | 内置实现 | ✅ |
| 熔断器中间件 | 内置实现 | ✅ |

---

## 📊 中间件执行顺序建议

建议的中间件执行顺序：

1. **RequestID** - 生成请求ID
2. **RealIP** - 获取真实IP
3. **采样中间件** - 决定是否采样（影响后续中间件的行为）
4. **追踪中间件** - OpenTelemetry 追踪（需要 RequestID）
5. **反射中间件** - 添加元数据信息
6. **数据转换中间件** - 转换请求数据
7. **精细控制中间件** - 功能开关、速率控制、熔断器
8. **限流中间件** - 请求限流
9. **认证中间件** - 身份验证
10. **日志中间件** - 请求日志（需要 RequestID 和 Tracing）
11. **恢复中间件** - Panic 恢复（保护所有后续处理）
12. **超时中间件** - 请求超时
13. **CORS 中间件** - 跨域支持（最后执行，处理响应头）

---

## 🎯 完整使用示例

```go
package main

import (
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/pkg/control"
    "github.com/yourusername/golang/pkg/sampling"
)

func main() {
    r := chi.NewRouter()

    // 1. 基础中间件
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)

    // 2. 采样中间件
    sampler, _ := sampling.NewProbabilisticSampler(0.5)
    r.Use(middleware.SamplingMiddleware(middleware.SamplingConfig{
        Sampler:             sampler,
        SkipPaths:           []string{"/health", "/metrics"},
        AddSamplingDecision: true,
    }))

    // 3. 追踪中间件
    r.Use(middleware.TracingMiddleware(middleware.TracingConfig{
        ServiceName:    "my-service",
        ServiceVersion: "v1.0.0",
        SkipPaths:      []string{"/health", "/metrics"},
    }))

    // 4. 反射中间件
    r.Use(middleware.ReflectMiddleware(middleware.ReflectConfig{
        EnableMetadata:     true,
        EnableSelfDescribe: true,
        SkipPaths:          []string{"/health", "/metrics"},
    }))

    // 5. 数据转换中间件
    r.Use(middleware.ConverterMiddleware(middleware.ConverterConfig{
        EnableRequestConversion:  true,
        EnableResponseConversion: true,
        DefaultResponseFormat:    "json",
    }))

    // 6. 精细控制中间件
    featureController := control.NewFeatureController()
    rateController := control.NewRateController()
    circuitController := control.NewCircuitController()

    featureController.Register("experimental-feature", "Experimental feature", true, nil)
    rateController.SetRateLimit("user-api", 100.0, time.Second)
    circuitController.RegisterCircuit("external-api", 10, 5, 30*time.Second)

    r.Use(middleware.ControlMiddleware(middleware.ControlConfig{
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

    // 7. 其他中间件
    r.Use(middleware.LoggingMiddleware)
    r.Use(middleware.RecoveryMiddleware)
    r.Use(middleware.TimeoutMiddleware(60 * time.Second))
    r.Use(middleware.CORSMiddleware)

    // 路由...
}
```

---

## 📚 相关文档

- [HTTP 中间件文档](../../internal/interfaces/http/chi/middleware/README.md)
- [框架核心能力总结](07-框架核心能力总结.md)
- [核心能力使用示例](08-核心能力使用示例.md)

---

**最后更新**: 2025-01-XX
