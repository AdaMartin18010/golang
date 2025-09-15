---
title: Go 应用接入规范（OpenTelemetry）
slug: observability-otel-go-guideline
topic: observability
level: guideline
goVersion: 1.21+
lastReviewed: 2025-09-15
owner: core-team
status: active
---

## 适用范围

- 面向 Go 服务与 CLI 工具的最小 OTel 接入规范，覆盖 Trace、Metrics、Logs 关联。

## 依赖与版本

```text
go.opentelemetry.io/otel
go.opentelemetry.io/otel/sdk
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc
go.opentelemetry.io/otel/metric
```

## 初始化与关停

```go
package obs

import (
    "context"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/sdk/resource"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Shutdown func(ctx context.Context) error

type Provider struct {
    TraceShutdown  Shutdown
    MetricShutdown Shutdown
}

func Init(serviceName, env string) (*Provider, error) {
    res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceName(serviceName),
        semconv.DeploymentEnvironment(env),
    ))

    // Trace exporter
    trExp, err := otlptracegrpc.New(context.Background())
    if err != nil { return nil, err }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithResource(res),
        sdktrace.WithBatcher(trExp),
    )
    otel.SetTracerProvider(tp)

    // Metric exporter
    mExp, err := otlpmetricgrpc.New(context.Background())
    if err != nil { return nil, err }
    mp := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mExp, sdkmetric.WithInterval(10*time.Second))),
        sdkmetric.WithResource(res),
    )
    metric.SetMeterProvider(mp)

    return &Provider{
        TraceShutdown:  tp.Shutdown,
        MetricShutdown: mp.Shutdown,
    }, nil
}

func (p *Provider) Shutdown(ctx context.Context) error {
    if p == nil { return nil }
    if p.MetricShutdown != nil { _ = p.MetricShutdown(ctx) }
    if p.TraceShutdown != nil { return p.TraceShutdown(ctx) }
    return nil
}
```

## HTTP 与数据库自动埋点（建议）

- HTTP：`go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`
- gRPC：`go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc`
- SQL：`go.opentelemetry.io/contrib/instrumentation/database/sql/otelsql`

## 日志关联（slog）

统一输出 `trace_id`、`span_id` 字段，建议封装 `ContextHandler` 自动提取：

```go
// 伪代码：在 slog.Handler 的 Handle 方法中从 ctx 抽取 TraceID/SpanID 并附加到 attrs
```

## 配置约定

- 环境变量优先：`OTEL_EXPORTER_OTLP_ENDPOINT`、`OTEL_SERVICE_NAME`、`OTEL_RESOURCE_ATTRIBUTES`、`OTEL_TRACES_SAMPLER` 等。
- 缺省采样 1-5%，关键链路 100%。

## 验收清单

- 启动时完成 OTel 初始化，无报错；关闭时优雅退出。
- HTTP 入口、下游调用与 DB 操作自动生成 span；日志带 trace_id。
- 指标至少包含 RED（Rate/Errors/Duration）三类核心指标。


