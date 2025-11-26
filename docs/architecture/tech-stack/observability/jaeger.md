# 1. 📊 Jaeger 深度解析

> **简介**: 本文档详细阐述了 Jaeger 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 📊 Jaeger 深度解析](#1--jaeger-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 配置 OpenTelemetry 导出到 Jaeger](#131-配置-opentelemetry-导出到-jaeger)
    - [1.3.2 在 Jaeger UI 中查看追踪](#132-在-jaeger-ui-中查看追踪)
    - [1.3.3 查询追踪数据](#133-查询追踪数据)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 追踪设计最佳实践](#141-追踪设计最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Jaeger 是什么？**

Jaeger 是一个开源的分布式追踪系统。

**核心特性**:

- ✅ **分布式追踪**: 完整的分布式追踪支持
- ✅ **可视化**: 直观的追踪可视化
- ✅ **查询**: 强大的查询功能
- ✅ **集成**: 与 OpenTelemetry 集成良好

---

## 1.2 选型论证

**为什么选择 Jaeger？**

**论证矩阵**:

| 评估维度 | 权重 | Jaeger | Zipkin | Datadog APM | New Relic | 说明 |
|---------|------|--------|--------|-------------|-----------|------|
| **OpenTelemetry 集成** | 35% | 10 | 9 | 8 | 7 | Jaeger 与 OTLP 集成最好 |
| **可视化能力** | 25% | 10 | 8 | 9 | 9 | Jaeger UI 功能完善 |
| **性能** | 20% | 9 | 8 | 8 | 8 | Jaeger 性能优秀 |
| **成本** | 15% | 10 | 10 | 3 | 3 | Jaeger 开源免费 |
| **社区生态** | 5% | 9 | 8 | 7 | 7 | Jaeger 社区活跃 |
| **加权总分** | - | **9.60** | 8.70 | 7.20 | 7.00 | Jaeger 得分最高 |

**核心优势**:

1. **OpenTelemetry 集成（权重 35%）**:
   - 原生支持 OTLP 协议，与 OpenTelemetry 集成完美
   - 支持 gRPC 和 HTTP 两种传输方式
   - 与 OpenTelemetry Collector 集成良好

2. **可视化能力（权重 25%）**:
   - 直观的追踪可视化，支持时间线视图
   - 强大的查询功能，支持多维度查询
   - 支持追踪对比和性能分析

3. **成本（权重 15%）**:
   - 完全开源免费，无授权成本
   - 自托管，数据可控
   - 适合中小型项目

**为什么不选择其他追踪系统？**

1. **Zipkin**:
   - ✅ 轻量级，易于部署
   - ❌ 与 OpenTelemetry 集成不如 Jaeger
   - ❌ 可视化能力不如 Jaeger
   - ❌ 查询功能不如 Jaeger 强大

2. **Datadog APM**:
   - ✅ 功能完善，SaaS 服务
   - ❌ 成本高，不适合中小型项目
   - ❌ 数据存储在第三方
   - ❌ 与 OpenTelemetry 集成不如 Jaeger

3. **New Relic**:
   - ✅ APM 功能强大
   - ❌ 成本高
   - ❌ 数据存储在第三方
   - ❌ 与 OpenTelemetry 集成不如 Jaeger

---

## 1.3 实际应用

### 1.3.1 配置 OpenTelemetry 导出到 Jaeger

**配置导出器**:

```go
// 配置 Jaeger 导出器
import (
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
)

func InitJaegerTracing(ctx context.Context, endpoint string) (*trace.TracerProvider, error) {
    // 创建 Jaeger 导出器
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
    if err != nil {
        return nil, err
    }

    // 创建 TracerProvider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("golang-service"),
        )),
    )

    otel.SetTracerProvider(tp)
    return tp, nil
}
```

### 1.3.2 在 Jaeger UI 中查看追踪

**查看追踪步骤**:

1. 访问 Jaeger UI: `http://localhost:16686`
2. 选择服务: 选择 "golang-service"
3. 查看追踪: 查看请求的完整追踪链路
4. 分析性能: 分析每个 Span 的执行时间

### 1.3.3 查询追踪数据

**查询示例**:

```go
// 通过 Jaeger API 查询追踪
import "github.com/jaegertracing/jaeger-client-go"

func QueryTraces(ctx context.Context, serviceName string, startTime, endTime time.Time) ([]*jaeger.Trace, error) {
    // 查询追踪数据
    // 实现查询逻辑
    return traces, nil
}
```

---

## 1.4 最佳实践

### 1.4.1 追踪设计最佳实践

**为什么需要良好的追踪设计？**

良好的追踪设计可以提高问题排查效率，便于性能分析。

**追踪设计原则**:

1. **Span 命名**: 使用清晰的 Span 名称
2. **属性设置**: 设置有意义的属性
3. **采样策略**: 根据环境配置合适的采样率
4. **错误记录**: 记录详细的错误信息

**实际应用示例**:

```go
// 追踪设计最佳实践
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "user.create")
    defer span.End()

    // 设置属性
    span.SetAttributes(
        attribute.String("user.email", req.Email),
        attribute.String("user.name", req.Name),
        attribute.String("http.method", r.Method),
        attribute.String("http.path", r.URL.Path),
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
3. **采样策略**: 生产环境使用较低的采样率，开发环境使用 100% 采样
4. **错误记录**: 使用 `span.RecordError()` 记录错误，便于问题排查

---

## 📚 扩展阅读

- [Jaeger 官方文档](https://www.jaegertracing.io/)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Jaeger 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
