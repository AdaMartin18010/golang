# 1. 📊 OpenTelemetry 深度解析

> **简介**: 本文档详细阐述了 OpenTelemetry 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 📊 OpenTelemetry 深度解析](#1--opentelemetry-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 追踪集成](#131-追踪集成)
    - [1.3.2 指标收集](#132-指标收集)
    - [1.3.3 日志集成](#133-日志集成)
    - [1.3.4 上下文传播](#134-上下文传播)
    - [1.3.5 采样策略配置](#135-采样策略配置)
    - [1.3.6 资源属性配置](#136-资源属性配置)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 追踪最佳实践](#141-追踪最佳实践)
    - [1.4.2 指标最佳实践](#142-指标最佳实践)
    - [1.4.3 日志最佳实践](#143-日志最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**OpenTelemetry 是什么？**

OpenTelemetry 是一个厂商中立的开源可观测性框架，提供统一的 API、SDK 和数据格式。

**核心特性**:

- ✅ **统一标准**: 统一的 API 和数据格式
- ✅ **跨语言**: 支持多种编程语言
- ✅ **可插拔**: 支持多种后端（Jaeger, Prometheus 等）
- ✅ **追踪、指标、日志**: 支持三种可观测性数据

---

## 1.2 选型论证

**为什么选择 OpenTelemetry？**

**论证矩阵**:

| 评估维度 | 权重 | OpenTelemetry | Prometheus | Jaeger | Zipkin | 说明 |
|---------|------|---------------|------------|--------|--------|------|
| **功能完整性** | 30% | 10 | 5 | 5 | 5 | OpenTelemetry 支持三支柱 |
| **标准兼容** | 25% | 10 | 6 | 7 | 7 | OpenTelemetry 是行业标准 |
| **后端选择** | 20% | 10 | 7 | 6 | 6 | OpenTelemetry 后端灵活 |
| **集成复杂度** | 15% | 8 | 8 | 7 | 7 | OpenTelemetry 集成简单 |
| **社区支持** | 10% | 9 | 10 | 8 | 7 | OpenTelemetry 社区活跃 |
| **加权总分** | - | **9.35** | 6.80 | 6.50 | 6.30 | OpenTelemetry 得分最高 |

**核心优势**:

1. **功能完整性（权重 30%）**:
   - 支持追踪、指标、日志三大支柱
   - 统一的可观测性解决方案
   - 避免多套系统集成

2. **标准兼容（权重 25%）**:
   - 行业标准，统一接口
   - 可以轻松切换后端
   - 未来兼容性好

3. **后端选择（权重 20%）**:
   - 可以导出到多种后端（Prometheus, Jaeger, Zipkin 等）
   - 不锁定特定厂商
   - 灵活的后端选择

**为什么不选择其他可观测性方案？**

1. **Prometheus**:
   - ✅ 监控标准，功能强大
   - ❌ 只支持指标，不支持追踪和日志
   - ❌ 需要与其他工具集成

2. **Jaeger**:
   - ✅ 分布式追踪功能强大
   - ❌ 只支持追踪，不支持指标和日志
   - ❌ 需要与其他工具集成

3. **Zipkin**:
   - ✅ 轻量级，易于部署
   - ❌ 只支持追踪，功能有限
   - ❌ 与 OpenTelemetry 集成不如 Jaeger

**详细论证请参考**: [技术对比矩阵](../../00-对比矩阵.md#44-选型决策论证)

---

## 1.3 实际应用

### 1.3.1 追踪集成

**初始化追踪**:

```go
// internal/infrastructure/observability/tracing.go
package observability

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracing(ctx context.Context, endpoint string) (*trace.TracerProvider, error) {
    // 创建资源
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String("golang-service"),
            semconv.ServiceVersionKey.String("1.0.0"),
        ),
    )
    if err != nil {
        return nil, err
    }

    // 创建导出器
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(endpoint),
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }

    // 创建 TracerProvider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(res),
        trace.WithSampler(trace.AlwaysSample()),
    )

    // 设置为全局 TracerProvider
    otel.SetTracerProvider(tp)

    return tp, nil
}
```

**在代码中使用追踪**:

```go
// 在 HTTP Handler 中使用
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "user.create")
    defer span.End()

    span.SetAttributes(
        attribute.String("user.email", req.Email),
        attribute.String("user.name", req.Name),
    )

    user, err := h.service.CreateUser(ctx, req)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        Error(w, http.StatusInternalServerError, err)
        return
    }

    span.SetAttributes(attribute.String("user.id", user.ID))
    Success(w, http.StatusCreated, user)
}

// 在数据库操作中使用
func (r *UserRepository) Create(ctx context.Context, user *User) error {
    ctx, span := tracer.Start(ctx, "db.user.create")
    defer span.End()

    span.SetAttributes(
        attribute.String("db.system", "postgresql"),
        attribute.String("db.operation", "create"),
    )

    // 执行数据库操作
    err := r.client.User.Create().
        SetEmail(user.Email).
        SetName(user.Name).
        Exec(ctx)

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    return nil
}
```

### 1.3.2 指标收集

**初始化指标**:

```go
// internal/infrastructure/observability/metrics.go
package observability

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
    "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/resource"
)

func InitMetrics(ctx context.Context, endpoint string) (*metric.MeterProvider, error) {
    // 创建资源
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String("golang-service"),
        ),
    )
    if err != nil {
        return nil, err
    }

    // 创建导出器
    exporter, err := otlpmetricgrpc.New(ctx,
        otlpmetricgrpc.WithEndpoint(endpoint),
        otlpmetricgrpc.WithInsecure(),
    )
    if err != nil {
        return nil, err
    }

    // 创建 MeterProvider
    mp := metric.NewMeterProvider(
        metric.WithReader(metric.NewPeriodicReader(exporter)),
        metric.WithResource(res),
    )

    // 设置为全局 MeterProvider
    otel.SetMeterProvider(mp)

    return mp, nil
}
```

**使用指标**:

```go
// 定义指标
var (
    requestCounter metric.Int64Counter
    requestDuration metric.Float64Histogram
)

func init() {
    meter := otel.Meter("golang-service")

    requestCounter, _ = meter.Int64Counter(
        "http_requests_total",
        metric.WithDescription("Total number of HTTP requests"),
    )

    requestDuration, _ = meter.Float64Histogram(
        "http_request_duration_seconds",
        metric.WithDescription("HTTP request duration in seconds"),
    )
}

// 在 Handler 中使用指标
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // 增加计数器
    requestCounter.Add(r.Context(), 1,
        attribute.String("method", r.Method),
        attribute.String("path", r.URL.Path),
    )

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)

    // 记录持续时间
    duration := time.Since(start).Seconds()
    requestDuration.Record(r.Context(), duration,
        attribute.String("method", r.Method),
        attribute.String("path", r.URL.Path),
        attribute.String("status", getStatus(err)),
    )

    // 处理响应
    if err != nil {
        Error(w, http.StatusInternalServerError, err)
        return
    }
    Success(w, http.StatusCreated, user)
}
```

### 1.3.3 日志集成

**结构化日志集成**:

```go
// internal/infrastructure/observability/logging.go
package observability

import (
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
    "log/slog"
)

// LogHandler 集成 OpenTelemetry 的日志处理器
type LogHandler struct {
    handler slog.Handler
}

func NewLogHandler(handler slog.Handler) *LogHandler {
    return &LogHandler{handler: handler}
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
    // 从上下文获取追踪信息
    span := trace.SpanFromContext(ctx)
    if span.IsRecording() {
        spanCtx := span.SpanContext()
        if spanCtx.IsValid() {
            r.AddAttrs(
                slog.String("trace_id", spanCtx.TraceID().String()),
                slog.String("span_id", spanCtx.SpanID().String()),
            )
        }
    }

    return h.handler.Handle(ctx, r)
}

// 使用示例
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    logger := slog.Default().With(
        "method", r.Method,
        "path", r.URL.Path,
    )

    logger.InfoContext(r.Context(), "Creating user",
        "email", req.Email,
        "name", req.Name,
    )

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)

    if err != nil {
        logger.ErrorContext(r.Context(), "Failed to create user",
            "error", err,
        )
        Error(w, http.StatusInternalServerError, err)
        return
    }

    logger.InfoContext(r.Context(), "User created successfully",
        "user_id", user.ID,
    )
    Success(w, http.StatusCreated, user)
}
```

### 1.3.4 上下文传播

**上下文传播示例**:

```go
// 在 HTTP 请求中传播上下文
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 从 HTTP Header 提取追踪信息
        ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

        // 创建新的 Span
        ctx, span := tracer.Start(ctx, r.URL.Path)
        defer span.End()

        // 将上下文传递给下一个 Handler
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 在 gRPC 中传播上下文
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // 从 gRPC Metadata 提取追踪信息
    md, _ := metadata.FromIncomingContext(ctx)
    ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.NewGRPCFieldsCarrier(md))

    // 创建新的 Span
    ctx, span := tracer.Start(ctx, info.FullMethod)
    defer span.End()

    // 调用处理函数
    resp, err := handler(ctx, req)

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    }

    return resp, err
}
```

### 1.3.5 采样策略配置

**采样策略配置**:

```go
// 配置采样策略
func NewTracerProviderWithSampling(ctx context.Context, endpoint string, sampleRate float64) (*trace.TracerProvider, error) {
    // 创建采样器
    sampler := trace.TraceIDRatioBased(sampleRate)

    // 或者使用父采样器
    sampler = trace.ParentBased(sampler)

    // 创建 TracerProvider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(res),
        trace.WithSampler(sampler),
    )

    return tp, nil
}

// 使用示例
// 生产环境：采样率 10%
tp, _ := NewTracerProviderWithSampling(ctx, endpoint, 0.1)

// 开发环境：采样率 100%
tp, _ := NewTracerProviderWithSampling(ctx, endpoint, 1.0)
```

### 1.3.6 资源属性配置

**资源属性配置**:

```go
// 配置资源属性
func NewResource(ctx context.Context) (*resource.Resource, error) {
    return resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String("golang-service"),
            semconv.ServiceVersionKey.String("1.0.0"),
            semconv.ServiceNamespaceKey.String("production"),
            attribute.String("environment", "production"),
            attribute.String("deployment.environment", "production"),
        ),
        resource.WithProcessRuntimeDescription(),
        resource.WithProcessRuntimeName(),
        resource.WithProcessRuntimeVersion(),
    )
}
```

---

## 1.4 最佳实践

### 1.4.1 追踪最佳实践

**为什么需要追踪？**

追踪可以帮助理解请求在分布式系统中的完整调用链，识别性能瓶颈和故障点。

**追踪最佳实践**:

1. **Span 命名**: 使用清晰的 Span 名称，如 "user.create"、"db.query"
2. **属性设置**: 设置有意义的属性，如用户 ID、请求参数
3. **错误记录**: 使用 `span.RecordError()` 记录错误
4. **采样策略**: 根据环境配置合适的采样率

**实际应用示例**:

```go
// 追踪最佳实践
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "user.create")
    defer span.End()

    // 设置属性
    span.SetAttributes(
        attribute.String("user.email", req.Email),
        attribute.String("user.name", req.Name),
    )

    // 业务逻辑
    user, err := h.service.CreateUser(ctx, req)

    if err != nil {
        // 记录错误
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        Error(w, http.StatusInternalServerError, err)
        return
    }

    // 设置成功属性
    span.SetAttributes(attribute.String("user.id", user.ID))
    span.SetStatus(codes.Ok, "User created successfully")
    Success(w, http.StatusCreated, user)
}
```

**最佳实践要点**:

1. **Span 命名**: 使用清晰的、有意义的 Span 名称
2. **属性设置**: 设置有助于调试和监控的属性
3. **错误记录**: 使用 `span.RecordError()` 记录错误
4. **采样策略**: 生产环境使用较低的采样率，开发环境使用 100% 采样

### 1.4.2 指标最佳实践

**为什么需要指标？**

指标可以帮助监控系统性能、识别趋势和异常。

**指标最佳实践**:

1. **指标命名**: 使用标准的指标命名规范
2. **标签选择**: 选择有意义的标签，避免高基数标签
3. **指标类型**: 根据场景选择合适的指标类型（Counter、Gauge、Histogram）
4. **聚合策略**: 合理设置聚合策略

**实际应用示例**:

```go
// 指标最佳实践
var (
    httpRequestsTotal = metric.Int64Counter(
        "http_requests_total",
        metric.WithDescription("Total number of HTTP requests"),
    )

    httpRequestDuration = metric.Float64Histogram(
        "http_request_duration_seconds",
        metric.WithDescription("HTTP request duration in seconds"),
        metric.WithUnit("s"),
    )
)

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // 增加计数器
    httpRequestsTotal.Add(r.Context(), 1,
        attribute.String("method", r.Method),
        attribute.String("path", r.URL.Path),
        attribute.String("status", "pending"),
    )

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)

    // 记录持续时间
    duration := time.Since(start).Seconds()
    status := "success"
    if err != nil {
        status = "error"
    }

    httpRequestDuration.Record(r.Context(), duration,
        attribute.String("method", r.Method),
        attribute.String("path", r.URL.Path),
        attribute.String("status", status),
    )

    // 处理响应
    if err != nil {
        Error(w, http.StatusInternalServerError, err)
        return
    }
    Success(w, http.StatusCreated, user)
}
```

**最佳实践要点**:

1. **指标命名**: 使用标准的指标命名规范（如 Prometheus 规范）
2. **标签选择**: 选择有意义的标签，避免高基数标签（如用户 ID）
3. **指标类型**: 根据场景选择合适的指标类型
4. **单位设置**: 为指标设置合适的单位

### 1.4.3 日志最佳实践

**为什么需要结构化日志？**

结构化日志可以提高日志的可读性和可查询性，便于日志分析和问题排查。

**日志最佳实践**:

1. **结构化日志**: 使用结构化日志格式（JSON）
2. **日志级别**: 合理使用日志级别（DEBUG、INFO、WARN、ERROR）
3. **上下文信息**: 在日志中包含追踪信息（TraceID、SpanID）
4. **敏感信息**: 避免在日志中记录敏感信息

**实际应用示例**:

```go
// 日志最佳实践
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    logger := slog.Default().With(
        "method", r.Method,
        "path", r.URL.Path,
        "request_id", middleware.GetReqID(r.Context()),
    )

    // 从上下文获取追踪信息
    span := trace.SpanFromContext(r.Context())
    if span.SpanContext().IsValid() {
        logger = logger.With(
            "trace_id", span.SpanContext().TraceID().String(),
            "span_id", span.SpanContext().SpanID().String(),
        )
    }

    logger.InfoContext(r.Context(), "Creating user",
        "email", req.Email,
        "name", req.Name,
    )

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)

    if err != nil {
        logger.ErrorContext(r.Context(), "Failed to create user",
            "error", err.Error(),
        )
        Error(w, http.StatusInternalServerError, err)
        return
    }

    logger.InfoContext(r.Context(), "User created successfully",
        "user_id", user.ID,
    )
    Success(w, http.StatusCreated, user)
}
```

**最佳实践要点**:

1. **结构化日志**: 使用结构化日志格式，便于日志分析
2. **日志级别**: 合理使用日志级别，避免过度日志
3. **上下文信息**: 在日志中包含追踪信息，便于关联
4. **敏感信息**: 避免在日志中记录敏感信息（如密码、Token）

---

## 📚 扩展阅读

- [OpenTelemetry 官方文档](https://opentelemetry.io/)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 OpenTelemetry 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
