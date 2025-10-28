# Go可观测性实践

Go可观测性实践完整指南，涵盖日志管理、指标监控、链路追踪和告警管理。

---

## 📚 核心内容

1. **日志管理**
   - 日志级别
   - 结构化日志
   - 日志聚合 (ELK)

2. **指标监控**
   - Prometheus
   - Grafana
   - 指标类型

3. **链路追踪**
   - OpenTelemetry
   - Jaeger
   - 分布式追踪

4. **告警管理**
   - 告警规则
   - 告警通知
   - On-call

---

## 🚀 Prometheus示例

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequests)
}

func handler(w http.ResponseWriter, r *http.Request) {
    httpRequests.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
    // 处理请求...
}
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**最后更新**: 2025-10-28

