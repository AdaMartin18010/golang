# 分布式追踪 (Distributed Tracing)

> **分类**: 工程与云原生
> **标签**: #tracing #opentelemetry #observability

---

## OpenTelemetry 集成

### 初始化 Tracer

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://localhost:14268/api/traces"),
    ))
    if err != nil {
        return nil, err
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
        )),
    )

    otel.SetTracerProvider(tp)
    return tp, nil
}
```

---

## Span 管理

### 创建 Span

```go
tracer := otel.Tracer("my-service")

func ProcessOrder(ctx context.Context, orderID string) error {
    // 创建 span
    ctx, span := tracer.Start(ctx, "process-order",
        trace.WithAttributes(
            attribute.String("order.id", orderID),
        ),
    )
    defer span.End()

    // 业务逻辑
    if err := validateOrder(ctx, orderID); err != nil {
        span.RecordError(err)
        return err
    }

    span.SetStatus(codes.Ok, "order processed")
    return nil
}
```

### 嵌套 Span

```go
func ParentOperation(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "parent")
    defer span.End()

    // 子操作会自动成为子 span
    ChildOperation(ctx)
}

func ChildOperation(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "child")
    defer span.End()

    // 记录事件
    span.AddEvent("processing", trace.WithAttributes(
        attribute.Int("item.count", 10),
    ))
}
```

---

## 上下文传播

### HTTP 传播

```go
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

// 服务端
handler := otelhttp.NewHandler(http.HandlerFunc(myHandler), "server-handler")

// 客户端
transport := otelhttp.NewTransport(http.DefaultTransport)
client := &http.Client{Transport: transport}
```

### gRPC 传播

```go
import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// 服务端
s := grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
)

// 客户端
conn, _ := grpc.Dial(addr,
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
)
```

---

## 与日志集成

```go
func LogWithTrace(ctx context.Context, msg string) {
    span := trace.SpanFromContext(ctx)
    spanContext := span.SpanContext()

    logger.Info(msg,
        zap.String("trace_id", spanContext.TraceID().String()),
        zap.String("span_id", spanContext.SpanID().String()),
    )
}
```

---

## 采样策略

```go
// 概率采样
tp := sdktrace.NewTracerProvider(
    sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),  // 10% 采样
)

// 父级采样
tp := sdktrace.NewTracerProvider(
    sdktrace.WithSampler(sdktrace.ParentBased(
        sdktrace.TraceIDRatioBased(0.1),
    )),
)
```
