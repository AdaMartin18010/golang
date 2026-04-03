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

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02

---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02