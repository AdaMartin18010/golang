# 可观测性

> **分类**: 工程与云原生

---

## 三大支柱

| 类型 | 工具 | 用途 |
|------|------|------|
| 日志 | ELK, Loki | 事件记录 |
| 指标 | Prometheus | 性能监控 |
| 追踪 | Jaeger | 请求链路 |

---

## 结构化日志

```go
import "go.uber.org/zap"

logger.Info("request processed",
    zap.String("method", "GET"),
    zap.Int("status", 200),
    zap.Duration("latency", 45*time.Millisecond),
)
```

---

## Prometheus 指标

```go
var requestCounter = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total HTTP requests",
    },
    []string{"method", "endpoint", "status"},
)

http.Handle("/metrics", promhttp.Handler())
```

---

## 分布式追踪

```go
var tracer = otel.Tracer("my-service")

func processOrder(ctx context.Context, orderID string) error {
    ctx, span := tracer.Start(ctx, "processOrder",
        trace.WithAttributes(attribute.String("order.id", orderID)),
    )
    defer span.End()

    return saveOrder(ctx, orderID)
}
```
