# 1. 📊 Prometheus 深度解析

> **简介**: 本文档详细阐述了 Prometheus 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 📊 Prometheus 深度解析](#1--prometheus-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 指标定义和暴露](#131-指标定义和暴露)
    - [1.3.2 在代码中使用指标](#132-在代码中使用指标)
    - [1.3.3 告警规则配置](#133-告警规则配置)
    - [1.3.4 服务发现配置](#134-服务发现配置)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 指标设计最佳实践](#141-指标设计最佳实践)
    - [1.4.2 告警规则最佳实践](#142-告警规则最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Prometheus 是什么？**

Prometheus 是一个开源的监控和告警系统。

**核心特性**:

- ✅ **指标收集**: 拉取模型收集指标
- ✅ **查询语言**: PromQL 查询语言
- ✅ **告警**: 支持告警规则
- ✅ **可视化**: 支持 Grafana 集成

---

## 1.2 选型论证

**为什么选择 Prometheus？**

**论证矩阵**:

| 评估维度 | 权重 | Prometheus | InfluxDB | Graphite | Datadog | 说明 |
|---------|------|-----------|----------|----------|---------|------|
| **标准兼容** | 30% | 10 | 7 | 6 | 8 | Prometheus 是事实标准 |
| **性能** | 25% | 9 | 10 | 7 | 8 | Prometheus 性能优秀 |
| **生态集成** | 20% | 10 | 7 | 6 | 9 | Prometheus 生态最丰富 |
| **学习成本** | 15% | 8 | 7 | 8 | 6 | Prometheus 学习曲线适中 |
| **成本** | 10% | 10 | 8 | 9 | 3 | Prometheus 开源免费 |
| **加权总分** | - | **9.30** | 7.90 | 7.05 | 7.20 | Prometheus 得分最高 |

**核心优势**:

1. **标准兼容性（权重 30%）**:
   - Prometheus 是事实上的监控标准，PromQL 被广泛采用
   - 与 OpenTelemetry 集成良好，支持 OTLP 协议
   - 与 Grafana 集成完美，生态成熟

2. **生态集成（权重 20%）**:
   - 丰富的 Exporter 生态，支持各种系统监控
   - 与 Kubernetes 集成良好，是云原生监控标准
   - 社区活跃，文档完善

3. **成本（权重 10%）**:
   - 完全开源免费，无授权成本
   - 自托管，数据可控
   - 适合中小型项目

**为什么不选择其他监控系统？**

1. **InfluxDB**:
   - ✅ 时序数据库性能优秀
   - ❌ 监控生态不如 Prometheus 丰富
   - ❌ 学习成本较高
   - ❌ 与 OpenTelemetry 集成不如 Prometheus

2. **Graphite**:
   - ✅ 简单易用
   - ❌ 功能相对简单
   - ❌ 生态不如 Prometheus 丰富
   - ❌ 性能不如 Prometheus

3. **Datadog**:
   - ✅ 功能完善，SaaS 服务
   - ❌ 成本高，不适合中小型项目
   - ❌ 数据存储在第三方，隐私性差
   - ❌ 依赖外部服务，可用性受限于服务商

---

## 1.3 实际应用

### 1.3.1 指标定义和暴露

**定义指标**:

```go
// internal/infrastructure/observability/prometheus.go
package observability

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP 请求计数器
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    // HTTP 请求持续时间直方图
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path", "status"},
    )

    // 活跃连接数
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )
)
```

**暴露指标端点**:

```go
// 暴露 Prometheus 指标
import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *chi.Mux {
    r := chi.NewRouter()

    // Prometheus 指标端点
    r.Handle("/metrics", promhttp.Handler())

    // 其他路由
    r.Route("/api/v1", func(r chi.Router) {
        r.Mount("/users", userRoutes())
    })

    return r
}
```

### 1.3.2 在代码中使用指标

**在 Handler 中使用指标**:

```go
// 在 Handler 中使用指标
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)

    // 记录指标
    duration := time.Since(start).Seconds()
    status := "success"
    if err != nil {
        status = "error"
    }

    httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
    httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)

    // 处理响应
    if err != nil {
        Error(w, http.StatusInternalServerError, err)
        return
    }
    Success(w, http.StatusCreated, user)
}
```

### 1.3.3 告警规则配置

**告警规则配置**:

```yaml
# configs/prometheus/alerts.yml
groups:
  - name: http_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status="error"}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors per second"

      - alert: HighLatency
        expr: histogram_quantile(0.95, http_request_duration_seconds) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "95th percentile latency is {{ $value }} seconds"
```

### 1.3.4 服务发现配置

**服务发现配置**:

```yaml
# configs/prometheus/prometheus.yml
scrape_configs:
  - job_name: 'golang-service'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
    scrape_timeout: 10s
```

---

## 1.4 最佳实践

### 1.4.1 指标设计最佳实践

**为什么需要良好的指标设计？**

良好的指标设计可以提高监控的有效性，便于问题排查和性能优化。

**指标设计原则**:

1. **指标命名**: 使用清晰的、有意义的指标名称
2. **标签选择**: 选择有意义的标签，避免高基数标签
3. **指标类型**: 根据场景选择合适的指标类型（Counter、Gauge、Histogram）
4. **单位统一**: 使用统一的单位

**实际应用示例**:

```go
// 指标设计最佳实践
var (
    // Counter: 累计值，只增不减
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"}, // 低基数标签
    )

    // Gauge: 当前值，可增可减
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )

    // Histogram: 分布值
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{0.1, 0.5, 1, 2, 5, 10}, // 自定义桶
        },
        []string{"method", "path"},
    )
)
```

**最佳实践要点**:

1. **指标命名**: 使用标准的命名规范（如 `http_requests_total`）
2. **标签选择**: 避免高基数标签（如用户 ID），使用低基数标签（如状态码）
3. **指标类型**: 根据场景选择合适的指标类型
4. **单位统一**: 使用统一的单位（如秒、字节）

### 1.4.2 告警规则最佳实践

**为什么需要告警规则？**

告警规则可以帮助及时发现系统问题，减少故障影响。

**告警规则设计原则**:

1. **告警阈值**: 设置合理的告警阈值
2. **告警持续时间**: 设置合理的告警持续时间，避免误报
3. **告警级别**: 区分不同级别的告警（critical、warning）
4. **告警信息**: 提供清晰的告警信息

**实际应用示例**:

```yaml
# 告警规则最佳实践
groups:
  - name: service_alerts
    rules:
      # 错误率告警
      - alert: HighErrorRate
        expr: rate(http_requests_total{status="error"}[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }}"

      # 延迟告警
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "95th percentile latency is {{ $value }} seconds"
```

**最佳实践要点**:

1. **告警阈值**: 根据历史数据设置合理的告警阈值
2. **告警持续时间**: 设置合理的持续时间，避免瞬时波动触发告警
3. **告警级别**: 区分不同级别的告警，优先处理 critical 告警
4. **告警信息**: 提供清晰的告警信息，便于快速定位问题

---

## 📚 扩展阅读

- [Prometheus 官方文档](https://prometheus.io/)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Prometheus 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
