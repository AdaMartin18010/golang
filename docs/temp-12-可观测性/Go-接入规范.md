
## 接入原则

- 默认开启 Trace 与 Metrics；Logs 采用结构化（slog）并注入 trace_id/span_id。
- 统一环境变量：
  - `OTEL_SERVICE_NAME`、`DEPLOY_ENV`（必填）
  - `OTEL_EXPORTER_OTLP_ENDPOINT`（默认 `http://localhost:4317` gRPC）
  - `OTEL_TRACES_SAMPLER`（默认 `parentbased_traceidratio`）与 `OTEL_TRACES_SAMPLER_ARG`（默认 `0.05`）
- 版本与依赖：`go.opentelemetry.io/otel`、`go.opentelemetry.io/contrib`（otelhttp 等）。

## 初始化模板

```go
package obs

import (
    "context"
    "os"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Shutdown func(context.Context) error

func Init(serviceName, env string) (Shutdown, Shutdown, error) {
    res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceName(serviceName),
        semconv.DeploymentEnvironment(env),
        semconv.ServiceVersion(os.Getenv("OTEL_SERVICE_VERSION")),
    ))

    // Trace
    texp, err := otlptracegrpc.New(context.Background())
    if err != nil { return nil, nil, err }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(texp),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)

    // Metrics（可按需启用）
    mexp, err := otlpmetricgrpc.New(context.Background())
    if err != nil { return tp.Shutdown, nil, err }
    mp := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mexp)),
        sdkmetric.WithResource(res),
    )
    otel.SetMeterProvider(mp)

    return tp.Shutdown, mp.Shutdown, nil
}
```

## HTTP 服务端/客户端埋点

```go
// 服务端
mux := http.NewServeMux()
mux.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("ok"))
}))
wrapped := otelhttp.NewHandler(mux, "http.server")
http.ListenAndServe(":8080", wrapped)

// 客户端
client := http.Client{ Transport: otelhttp.NewTransport(http.DefaultTransport) }
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://svc/hello", nil)
client.Do(req)
```

## 数据库/队列

- 优先使用 OTel 已有 instrumentation（`go.opentelemetry.io/contrib`）。
- 未覆盖时以 Span 包裹关键操作；记录 `db.system`/`messaging.system` 等语义属性。

## 日志关联（slog）

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
ctx, span := otel.Tracer("svc").Start(context.Background(), "op")
defer span.End()
sc := trace.SpanFromContext(ctx).SpanContext()
logger.Info("op done", "trace_id", sc.TraceID().String(), "span_id", sc.SpanID().String())
```

## 采样与标签治理

- 默认比例采样 5%；关键端点配置 100%。
- 限制标签基数：用户/订单号等高基数不作为标签；改写入日志正文或事件。

## 健康检查

- 进程启动后验证：
  - 查看 Collector 接收成功计数
  - Grafana 内是否可见 Trace 与 Metrics
  - 日志是否带 `trace_id`

---

## Go 应用接入规范（OpenTelemetry）

slug: observability-otel-go-guideline
topic: observability
level: guideline
goVersion: 1.21+
lastReviewed: 2025-09-15
owner: core-team
status: active

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

## 日志关联1（slog）

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
