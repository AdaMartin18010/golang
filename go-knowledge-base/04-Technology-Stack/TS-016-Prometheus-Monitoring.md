# TS-016: Prometheus Monitoring - Metrics Collection & Alerting

> **维度**: Technology Stack
> **级别**: S (16+ KB)
> **标签**: #prometheus #monitoring #metrics #alerting #observability
> **权威来源**:
>
> - [Prometheus Documentation](https://prometheus.io/docs/introduction/overview/) - Prometheus.io
> - [Prometheus: Up & Running](https://www.oreilly.com/library/view/prometheus-up/9781492034148/) - O'Reilly Media
> - [Prometheus Best Practices](https://prometheus.io/docs/practices/) - Prometheus.io

---

## 1. Prometheus Architecture

### 1.1 Core Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Prometheus Monitoring Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Prometheus Server                                   │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Retrieval (Scraping)                          │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │  │  │
│  │  │  │ HTTP GET    │  │ HTTP GET    │  │ HTTP GET    │             │  │  │
│  │  │  │ /metrics    │  │ /metrics    │  │ /metrics    │             │  │  │
│  │  │  │ (Target 1)  │  │ (Target 2)  │  │ (Target N)  │             │  │  │
│  │  │  │ every 15s   │  │ every 15s   │  │ every 15s   │             │  │  │
│  │  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘             │  │  │
│  │  │         └────────────────┼─────────────────┘                    │  │  │
│  │  │                          │                                      │  │  │
│  │  │                          ▼                                      │  │  │
│  │  │  ┌──────────────────────────────────────────────────────────┐  │  │  │
│  │  │  │                    Parse & Expose Formats                 │  │  │  │
│  │  │  │                                                           │  │  │  │
│  │  │  │  • Prometheus text format (default)                       │  │  │  │
│  │  │  │  • OpenMetrics                                            │  │  │  │
│  │  │  │  • Protocol Buffers (legacy)                              │  │  │  │
│  │  │  └──────────────────────────────────────────────────────────┘  │  │  │
│  │  │                          │                                      │  │  │
│  │  │                          ▼                                      │  │  │
│  │  │  ┌──────────────────────────────────────────────────────────┐  │  │  │
│  │  │  │                    Service Discovery                      │  │  │  │
│  │  │  │                                                           │  │  │  │
│  │  │  │  Static:  File-based target configuration                 │  │  │  │
│  │  │  │  Dynamic: Kubernetes, Consul, EC2, Azure, GCE, etc.      │  │  │  │
│  │  │  │                                                           │  │  │  │
│  │  │  │  Relabel configs: Filter and modify targets               │  │  │  │
│  │  │  └──────────────────────────────────────────────────────────┘  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Storage (TSDB)                                │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌──────────────────────────────────────────────────────────┐  │  │  │
│  │  │  │                    Time Series Database                    │  │  │  │
│  │  │  │                                                           │  │  │  │
│  │  │  │  On-disk:                                                 │  │  │  │
│  │  │  │  • Chunk files (compressed samples)                       │  │  │  │
│  │  │  │  • Index files (inverted index for labels)               │  │  │  │
│  │  │  │  • Tombstones (deleted series)                           │  │  │  │
│  │  │  │  • WAL (write-ahead log)                                 │  │  │  │
│  │  │  │                                                           │  │  │  │
│  │  │  │  In-memory:                                               │  │  │  │
│  │  │  │  • Head block (recent samples, mutable)                  │  │  │  │
│  │  │  │  • Series index (label → series ID mapping)              │  │  │  │
│  │  │  └──────────────────────────────────────────────────────────┘  │  │  │
│  │  │                                                                  │  │  │
│  │  │  Data Model:                                                   │  │  │
│  │  │  • Metric name + Labels = Time Series                         │  │  │
│  │  │  • Timestamp + Value = Sample                                 │  │  │
│  │  │                                                                  │  │  │
│  │  │  Example:                                                      │  │  │
│  │  │  http_requests_total{method="GET",status="200"}                │  │  │
│  │  │  ├── 1704067200, 100                                          │  │  │
│  │  │  ├── 1704067215, 150                                          │  │  │
│  │  │  └── 1704067230, 175                                          │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    PromQL Engine                                 │  │  │
│  │  │                                                                  │  │  │
│  │  │  Query Processing:                                              │  │  │
│  │  │  1. Parse PromQL expression                                     │  │  │
│  │  │  2. Resolve label matchers (use index to find series)          │  │  │
│  │  │  3. Fetch chunks from disk/memory                               │  │  │
│  │  │  4. Apply functions/aggregations                                │  │  │
│  │  │  5. Return results                                              │  │  │
│  │  │                                                                  │  │  │
│  │  │  Example: sum(rate(http_requests_total[5m])) by (status)       │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Alert Manager Integration                     │  │  │
│  │  │                                                                  │  │  │
│  │  │  Recording Rules:  Pre-compute expensive queries                │  │  │
│  │  │  Alerting Rules:   Define alert conditions                      │  │  │
│  │  │                                                                  │  │  │
│  │  │  Evaluation: Every evaluation_interval (default 15s)            │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Pushgateway (for short-lived jobs)                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Batch Job ──► Push metrics ──► Pushgateway ──► Prometheus scrape     │  │
│  │  (short-lived)   (before exit)   (acts as proxy)                      │  │
│  │                                                                        │  │
│  │  Use case: CI/CD pipelines, cron jobs, batch processing               │  │
│  │  Best practice: Use service discovery instead when possible           │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Federation (for scaling)                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌───────────────┐        ┌───────────────┐                           │  │
│  │  │ Prometheus A  │───────►│ Prometheus C  │  (Federation)             │  │
│  │  │ (dc1 metrics) │        │ (Global view) │                           │  │
│  │  └───────────────┘        └───────┬───────┘                           │  │
│  │  ┌───────────────┐                │                                   │  │
│  │  │ Prometheus B  │───────────────┘                                   │  │
│  │  │ (dc2 metrics) │                                                    │  │
│  │  └───────────────┘                                                    │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Metric Types

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Prometheus Metric Types                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ 1. Counter                                                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  • Monotonically increasing (can only go up or reset to 0)           │  │
│  │  • Use for: total requests, total errors, bytes sent/received         │  │
│  │                                                                        │  │
│  │  Example: http_requests_total                                         │  │
│  │  ┌─────────────────────────────────────────────────────────────┐     │  │
│  │  │  Value                                                      │     │  │
│  │  │    ▲                                                        │     │  │
│  │  │    │     ╱╲                                                │     │  │
│  │  │    │    ╱  ╲╱╲                                             │     │  │
│  │  │    │   ╱      ╲╱╲                                          │     │  │
│  │  │    │  ╱          ╲___                                      │     │  │
│  │  │    │╱                  ╲                                   │     │  │
│  │  │    └──────────────────────────► Time                      │     │  │
│  │  │         Resets at restart                                   │     │  │
│  │  └─────────────────────────────────────────────────────────────┘     │  │
│  │                                                                        │  │
│  │  PromQL: rate(http_requests_total[5m])  (gives requests/second)       │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ 2. Gauge                                                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  • Can go up and down                                                 │  │
│  │  • Use for: temperature, memory usage, queue size, current connections│  │
│  │                                                                        │  │
│  │  Example: memory_usage_bytes                                          │  │
│  │  ┌─────────────────────────────────────────────────────────────┐     │  │
│  │  │  Value                                                      │     │  │
│  │  │    ▲     ╱╲  ╱╲                                             │     │  │
│  │  │    │    ╱  ╲╱  ╲                                            │     │  │
│  │  │    │   ╱        ╲___                                        │     │  │
│  │  │    │  ╱              ╲                                      │     │  │
│  │  │    │ ╱                ╲                                     │     │  │
│  │  │    └──────────────────────────► Time                      │     │  │
│  │  │         Goes up and down                                    │     │  │
│  │  └─────────────────────────────────────────────────────────────┘     │  │
│  │                                                                        │  │
│  │  PromQL: memory_usage_bytes / 1024 / 1024  (convert to MB)            │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ 3. Histogram                                                          │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  • Samples observations into configurable buckets                     │  │
│  │  • Use for: request duration, response size                           │  │
│  │  • Counts total observations + sum of values + per-bucket counts      │  │
│  │                                                                        │  │
│  │  Example: http_request_duration_seconds{bucket="le=0.1"}              │  │
│  │           http_request_duration_seconds{bucket="le=0.5"}              │  │
│  │           http_request_duration_seconds{bucket="le=1.0"}              │  │
│  │           http_request_duration_seconds{bucket="le=+Inf"}             │  │
│  │           http_request_duration_seconds_sum                           │  │
│  │           http_request_duration_seconds_count                         │  │
│  │                                                                        │  │
│  │  Buckets:                                                              │  │
│  │  ┌─────────────────────────────────────────────────────────────┐     │  │
│  │  │  Count                                                      │     │  │
│  │  │    │    ┌──┐                                               │     │  │
│  │  │    │    │  │  ┌──────┐                                     │     │  │
│  │  │    │    │  │  │      │  ┌────┐                             │     │  │
│  │  │    │    │  │  │      │  │    │  ┌──┐                       │     │  │
│  │  │    └────┴──┴──┴──────┴──┴────┴──┴──┴──────────►           │     │  │
│  │  │         0.1   0.5     1.0    5.0    +Inf  (seconds)       │     │  │
│  │  └─────────────────────────────────────────────────────────────┘     │  │
│  │                                                                        │  │
│  │  PromQL: histogram_quantile(0.95,                                   │  │
│  │            rate(http_request_duration_seconds_bucket[5m]))            │  │
│  │            (calculates 95th percentile)                               │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │ 4. Summary (legacy, use Histogram instead)                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  • Similar to histogram but calculates quantiles on client side       │  │
│  │  • Cannot be aggregated across instances                              │  │
│  │  • More expensive client-side                                         │  │
│  │  • Use histogram_quantile() with histograms instead                   │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go Implementation

```go
package prometheus

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics 指标集合
type Metrics struct {
    // Counter
    RequestsTotal   *prometheus.CounterVec
    ErrorsTotal     *prometheus.CounterVec

    // Gauge
    ActiveConnections prometheus.Gauge
    QueueSize         prometheus.Gauge

    // Histogram
    RequestDuration *prometheus.HistogramVec
    RequestSize     *prometheus.HistogramVec
}

// NewMetrics 创建指标
func NewMetrics() *Metrics {
    return &Metrics{
        RequestsTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "endpoint", "status"},
        ),
        ErrorsTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_errors_total",
                Help: "Total number of HTTP errors",
            },
            []string{"method", "endpoint", "error_type"},
        ),
        ActiveConnections: promauto.NewGauge(
            prometheus.GaugeOpts{
                Name: "active_connections",
                Help: "Number of active connections",
            },
        ),
        QueueSize: promauto.NewGauge(
            prometheus.GaugeOpts{
                Name: "queue_size",
                Help: "Current queue size",
            },
        ),
        RequestDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "http_request_duration_seconds",
                Help:    "HTTP request duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "endpoint"},
        ),
        RequestSize: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "http_request_size_bytes",
                Help:    "HTTP request size in bytes",
                Buckets: prometheus.ExponentialBuckets(100, 10, 8),
            },
            []string{"method", "endpoint"},
        ),
    }
}

// RecordRequest 记录请求
func (m *Metrics) RecordRequest(method, endpoint, status string, duration time.Duration) {
    m.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()
    m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordError 记录错误
func (m *Metrics) RecordError(method, endpoint, errorType string) {
    m.ErrorsTotal.WithLabelValues(method, endpoint, errorType).Inc()
}

// SetActiveConnections 设置活跃连接数
func (m *Metrics) SetActiveConnections(n float64) {
    m.ActiveConnections.Set(n)
}

// StartMetricsServer 启动指标服务
func StartMetricsServer(addr string) error {
    http.Handle("/metrics", promhttp.Handler())
    return http.ListenAndServe(addr, nil)
}
```

---

## 3. Configuration Best Practices

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'production'
    replica: '{{.ExternalURL}}'

# 告警规则
rule_files:
  - /etc/prometheus/rules/*.yml

# 抓取配置
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance

  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true

# 远程存储 (可选)
remote_write:
  - url: "http://cortex:9009/api/prom/push"
    queue_config:
      max_samples_per_send: 1000
      max_shards: 200

# 告警管理器
alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']
```

---

## 4. Visual Representations

### Alert Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Prometheus Alerting Pipeline                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Alert Rule Evaluation                                                  │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │  groups:                                                             │   │
│  │    - name: api_alerts                                                │   │
│  │      rules:                                                          │   │
│  │        - alert: HighErrorRate                                        │   │
│  │          expr: |                                                     │   │
│  │            sum(rate(http_requests_total{status=~"5.."}[5m]))         │   │
│  │            /                                                         │   │
│  │            sum(rate(http_requests_total[5m])) > 0.05                 │   │
│  │          for: 5m                                                     │   │
│  │          labels:                                                     │   │
│  │            severity: critical                                        │   │
│  │          annotations:                                                │   │
│  │            summary: "High error rate detected"                       │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                │                                             │
│                                ▼                                             │
│  2. Alert Firing                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │  Prometheus ──► Alert (FIRING) ──► Alert Manager                     │   │
│  │  (evaluation)     for: 5m satisfied    (routing)                     │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                │                                             │
│                                ▼                                             │
│  3. Alert Manager Routing                                                  │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │  Route: severity=critical ──► PagerDuty                              │   │
│  │  Route: severity=warning  ──► Slack                                  │   │
│  │  Route: team=backend      ──► Backend On-Call                        │   │
│  │                                                                      │   │
│  │  Grouping:                                                           │   │
│  │  • alertname=HighErrorRate, region=us-east ──► Single notification   │   │
│  │                                                                      │   │
│  │  Inhibition:                                                         │   │
│  │  • If ClusterDown firing, suppress NodeDown alerts                   │   │
│  │                                                                      │   │
│  │  Silencing:                                                          │   │
│  │  • Scheduled maintenance windows                                     │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. References

1. **Prometheus Documentation** (2024). prometheus.io/docs
2. **Turnbull, J.** (2018). Prometheus: Up & Running. O'Reilly Media.
3. **Prometheus Community** (2024). prometheus/client_golang.

---

*Document Version: 1.0 | Last Updated: 2024*
