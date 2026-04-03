# EC-051: Metrics Collection Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #metrics #prometheus #opentelemetry #observability #timeseries #monitoring
> **Authoritative Sources**:
>
> - [Google SRE Book - Monitoring](https://sre.google/sre-book/monitoring-distributed-systems/) - Google (2017)
> - [Prometheus: Up & Running](https://www.oreilly.com/library/view/prometheus-up/9781492034141/) - Brian Brazil (2018)
> - [OpenTelemetry Metrics](https://opentelemetry.io/docs/concepts/signals/metrics/) - CNCF (2024)
> - [The RED Method](https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/) - Weaveworks (2015)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Metrics Domain)**
Let $\mathcal{M}$ be the set of observable measurements from system $S$ over time $T$:

- Time series: $m(t) \in \mathcal{M}$ where $t \in T$
- Dimensions: $d(m) = \{label_1, label_2, ..., label_k\}$ providing cardinality

**Metric Types:**

| Type | Mathematical Definition | Use Case |
|------|------------------------|----------|
| **Counter** | $C: T \to \mathbb{N}$, monotonic increasing | Requests served, errors |
| **Gauge** | $G: T \to \mathbb{R}$, arbitrary | Temperature, queue depth |
| **Histogram** | $H: T \to \mathcal{P}(\mathbb{R})$ | Request latency distribution |
| **Summary** | $S: T \to \{quantiles\}$ | Pre-calculated percentiles |

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Cardinality Explosion** | $\prod_{i=1}^{k} |label_i| \to \infty$ | Storage and query degradation |
| **Time Resolution** | $\Delta t_{sample} < \Delta t_{event}$ | Events may be missed |
| **Resource Overhead** | $overhead(collection) < \theta_{budget}$ | Trade-off with granularity |
| **Retention Limits** | $|\mathcal{M}_{stored}| \leq capacity$ | Aggregation and downsampling |

### 1.2 Problem Statement

**Problem 1.1 (Observability Coverage)**
Given system $S$ with observability requirements $\mathcal{R} = \{R_1, R_2, ..., R_m\}$, select metric set $\mathcal{M}$ such that:

$$\forall R_i \in \mathcal{R}: \exists \mathcal{M}_i \subseteq \mathcal{M}: satisfies(\mathcal{M}_i, R_i)$$

**Key Challenges:**

1. **Signal Selection**: Choose metrics that provide actionable insights
2. **Cardinality Management**: Prevent label explosion
3. **Collection Efficiency**: Minimize overhead while maintaining coverage
4. **Temporal Aggregation**: Balance resolution with retention
5. **Correlation**: Link metrics to traces and logs

### 1.3 Formal Requirements Specification

**Requirement 1.1 (Completeness)**
$$\forall failure\_mode: detectable(failure\_mode) \Rightarrow \exists m \in \mathcal{M}: alert(m)$$

**Requirement 1.2 (Cardinality Bound)**
$$\forall m \in \mathcal{M}: cardinality(m) < C_{max}$$

---

## 2. Solution Architecture

### 2.1 Metrics Taxonomy (RED, USE, Golden Signals)

**The RED Method (Request-oriented):**

- **Rate**: Requests per second
- **Errors**: Error rate
- **Duration**: Request latency

**The USE Method (Resource-oriented):**

- **Utilization**: % time busy
- **Saturation**: Queue depth / work backlog
- **Errors**: Error count

**The Four Golden Signals:**

- Latency, Traffic, Errors, Saturation

---

## 3. Visual Representations

### 3.1 Metrics Collection Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    METRICS COLLECTION ARCHITECTURE                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  APPLICATION LAYER                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  Service A          Service B          Service C                     │   │
│  │  ┌─────────┐       ┌─────────┐       ┌─────────┐                    │   │
│  │  │ Metrics │       │ Metrics │       │ Metrics │                    │   │
│  │  │ Counter │       │ Gauge   │       │Histogram│                    │   │
│  │  │ requests│       │active   │       │duration │                    │   │
│  │  │  total  │       │connections     │         │                    │   │
│  │  └────┬────┘       └────┬────┘       └────┬────┘                    │   │
│  │       │                 │                 │                          │   │
│  └───────┼─────────────────┼─────────────────┼──────────────────────────┘   │
│          │                 │                 │                               │
│          └─────────────────┴─────────────────┘                               │
│                        │                                                     │
│                        │ /metrics endpoint                                   │
│                        ▼                                                     │
│  COLLECTION LAYER                                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Prometheus Server                            │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │   Service   │  │   Service   │  │   Service   │                  │   │
│  │  │  Discovery  │  │   Scraping  │  │   Storage   │                  │   │
│  │  │             │──►│   15s int.  │──►│  TSDB       │                  │   │
│  │  │• Kubernetes │  │• Pull model │  │• 15d local  │                  │   │
│  │  │  pods       │  │• HTTP       │  │• Compression│                  │   │
│  │  │• Consul     │  │             │  │• WAL        │                  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                  │   │
│  │                                                                      │   │
│  │  Alerting Rules:                                                     │   │
│  │  • error_rate > 1%                                                   │   │
│  │  • latency_p99 > 500ms                                               │   │
│  │  • cpu_utilization > 80%                                             │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                   │                                         │
│          ┌────────────────────────┼────────────────────────┐               │
│          │                        │                        │               │
│          ▼                        ▼                        ▼               │
│  ┌───────────────┐      ┌───────────────┐      ┌───────────────────┐       │
│  │    Alert      │      │    Grafana    │      │  Remote Storage   │       │
│  │   Manager     │      │               │      │                   │       │
│  │               │      │  Dashboards   │      │  • Thanos         │       │
│  │ • PagerDuty   │      │  • Panels     │      │  • Cortex         │       │
│  │ • Slack       │      │  • Alerts     │      │  • InfluxDB       │       │
│  │ • Webhook     │      │               │      │  • Long-term      │       │
│  └───────────────┘      └───────────────┘      └───────────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Histogram and Quantile Computation

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HISTOGRAM BUCKETING AND QUANTILES                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Raw Latency Measurements (ms):                                             │
│  23, 45, 67, 89, 123, 156, 178, 234, 267, 289, 345, 378, 456, 567, 678     │
│                                                                             │
│  Histogram Buckets (exponential):                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Bucket         │  Count  │  Cumulative  │  Visual Representation   │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │  le=0.005       │    0    │      0       │                          │   │
│  │  le=0.01        │    0    │      0       │                          │   │
│  │  le=0.025       │    0    │      0       │                          │   │
│  │  le=0.05        │    0    │      0       │                          │   │
│  │  le=0.1         │    0    │      0       │                          │   │
│  │  le=0.25        │    4    │      4       │ ████                     │   │
│  │  le=0.5         │   11    │     11       │ ███████████              │   │
│  │  le=1.0         │   15    │     15       │ ███████████████          │   │
│  │  le=2.5         │   15    │     15       │ ███████████████          │   │
│  │  le=5.0         │   15    │     15       │ ███████████████          │   │
│  │  le=10.0        │   15    │     15       │ ███████████████          │   │
│  │  +Inf           │   15    │     15       │ ███████████████          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  Quantile Estimation (histogram_quantile):                                  │
│                                                                             │
│  P50 (median):  ~250ms                                                      │
│  P90:           ~400ms                                                      │
│  P95:           ~500ms                                                      │
│  P99:           ~650ms                                                      │
│                                                                             │
│  Note: These are estimates based on bucket boundaries, not exact values     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Cardinality and Label Management

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    METRIC CARDINALITY MANAGEMENT                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  PROBLEMATIC (High Cardinality):                                            │
│  ───────────────────────────────                                            │
│                                                                             │
│  http_requests_total{                                                      │
│    method="GET",                                                            │
│    path="/api/users/12345",    ◄── Each user ID creates new series         │
│    status="200",                                                            │
│    user_id="12345"            ◄── High cardinality label                   │
│  }                                                                          │
│                                                                             │
│  Cardinality = 4 methods × 100K users × 5 statuses = 2,000,000 series       │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  OPTIMIZED (Low Cardinality):                                               │
│  ────────────────────────────                                               │
│                                                                             │
│  http_requests_total{                                                       │
│    method="GET",                                                            │
│    path="/api/users/:id",      ◄── Parameterized path                      │
│    status="200"                                                             │
│  }                                                                          │
│                                                                             │
│  Cardinality = 4 methods × ~100 paths × 5 statuses = 2,000 series           │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  ALTERNATIVE: Use separate metric for user-level tracking                   │
│                                                                             │
│  user_request_counter{                                                      │
│    user_id="12345"            ◄── Only if needed for billing/analytics     │
│  }                                                                          │
│                                                                             │
│  This metric goes to analytics pipeline, not Prometheus                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 OpenTelemetry Metrics Setup

```go
package metrics

import (
 "context"
 "fmt"
 "time"

 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/exporters/prometheus"
 "go.opentelemetry.io/otel/metric"
 sdkmetric "go.opentelemetry.io/otel/sdk/metric"
 "go.opentelemetry.io/otel/sdk/resource"
 semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Config holds metrics configuration
type Config struct {
 ServiceName    string
 ServiceVersion string
 Environment    string
 ExportInterval time.Duration
}

// MetricsProvider manages OpenTelemetry metrics
type MetricsProvider struct {
 provider *sdkmetric.MeterProvider
 meter    metric.Meter
}

// New creates a new metrics provider with Prometheus exporter
func New(cfg Config) (*MetricsProvider, error) {
 // Create Prometheus exporter
 exporter, err := prometheus.New()
 if err != nil {
  return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
 }

 // Create resource
 res, err := resource.New(context.Background(),
  resource.WithAttributes(
   semconv.ServiceName(cfg.ServiceName),
   semconv.ServiceVersion(cfg.ServiceVersion),
   attribute.String("environment", cfg.Environment),
  ),
 )
 if err != nil {
  return nil, fmt.Errorf("failed to create resource: %w", err)
 }

 // Create meter provider
 provider := sdkmetric.NewMeterProvider(
  sdkmetric.WithReader(exporter),
  sdkmetric.WithResource(res),
 )

 return &MetricsProvider{
  provider: provider,
  meter:    provider.Meter(cfg.ServiceName),
 }, nil
}

// Shutdown gracefully shuts down the provider
func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
 return mp.provider.Shutdown(ctx)
}

// Meter returns the metric meter
func (mp *MetricsProvider) Meter() metric.Meter {
 return mp.meter
}

// HTTPMetrics holds HTTP-specific metrics
type HTTPMetrics struct {
 RequestsTotal   metric.Int64Counter
 RequestDuration metric.Float64Histogram
 RequestSize     metric.Int64Histogram
 ResponseSize    metric.Int64Histogram
 ActiveRequests  metric.Int64UpDownCounter
}

// NewHTTPMetrics creates HTTP metrics
func NewHTTPMetrics(meter metric.Meter) (*HTTPMetrics, error) {
 requestsTotal, err := meter.Int64Counter(
  "http_requests_total",
  metric.WithDescription("Total number of HTTP requests"),
 )
 if err != nil {
  return nil, err
 }

 requestDuration, err := meter.Float64Histogram(
  "http_request_duration_seconds",
  metric.WithDescription("HTTP request duration in seconds"),
 )
 if err != nil {
  return nil, err
 }

 requestSize, err := meter.Int64Histogram(
  "http_request_size_bytes",
  metric.WithDescription("HTTP request size in bytes"),
 )
 if err != nil {
  return nil, err
 }

 responseSize, err := meter.Int64Histogram(
  "http_response_size_bytes",
  metric.WithDescription("HTTP response size in bytes"),
 )
 if err != nil {
  return nil, err
 }

 activeRequests, err := meter.Int64UpDownCounter(
  "http_active_requests",
  metric.WithDescription("Number of active HTTP requests"),
 )
 if err != nil {
  return nil, err
 }

 return &HTTPMetrics{
  RequestsTotal:   requestsTotal,
  RequestDuration: requestDuration,
  RequestSize:     requestSize,
  ResponseSize:    responseSize,
  ActiveRequests:  activeRequests,
 }, nil
}

// RecordRequest records an HTTP request
func (hm *HTTPMetrics) RecordRequest(
 ctx context.Context,
 method, path, status string,
 duration time.Duration,
 requestSize, responseSize int64,
) {
 attrs := []attribute.KeyValue{
  attribute.String("method", method),
  attribute.String("path", path),
  attribute.String("status", status),
 }

 hm.RequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
 hm.RequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
 hm.RequestSize.Record(ctx, requestSize, metric.WithAttributes(attrs...))
 hm.ResponseSize.Record(ctx, responseSize, metric.WithAttributes(attrs...))
}

// DatabaseMetrics holds database-specific metrics
type DatabaseMetrics struct {
 QueriesTotal    metric.Int64Counter
 QueryDuration   metric.Float64Histogram
 ConnectionsOpen metric.Int64UpDownCounter
 ConnectionWait  metric.Float64Histogram
}

// NewDatabaseMetrics creates database metrics
func NewDatabaseMetrics(meter metric.Meter) (*DatabaseMetrics, error) {
 queriesTotal, err := meter.Int64Counter(
  "db_queries_total",
  metric.WithDescription("Total number of database queries"),
 )
 if err != nil {
  return nil, err
 }

 queryDuration, err := meter.Float64Histogram(
  "db_query_duration_seconds",
  metric.WithDescription("Database query duration in seconds"),
 )
 if err != nil {
  return nil, err
 }

 connectionsOpen, err := meter.Int64UpDownCounter(
  "db_connections_open",
  metric.WithDescription("Number of open database connections"),
 )
 if err != nil {
  return nil, err
 }

 connectionWait, err := meter.Float64Histogram(
  "db_connection_wait_seconds",
  metric.WithDescription("Time spent waiting for a connection"),
 )
 if err != nil {
  return nil, err
 }

 return &DatabaseMetrics{
  QueriesTotal:    queriesTotal,
  QueryDuration:   queryDuration,
  ConnectionsOpen: connectionsOpen,
  ConnectionWait:  connectionWait,
 }, nil
}
```

### 4.2 Prometheus Direct Integration

```go
package metrics

import (
 "net/http"
 "time"

 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMetrics holds Prometheus metrics
type PrometheusMetrics struct {
 Registry *prometheus.Registry

 // HTTP metrics
 RequestsTotal   *prometheus.CounterVec
 RequestDuration *prometheus.HistogramVec

 // Business metrics
 OrdersProcessed *prometheus.CounterVec
 OrderValue      *prometheus.HistogramVec
}

// NewPrometheusMetrics creates Prometheus metrics
func NewPrometheusMetrics() *PrometheusMetrics {
 reg := prometheus.NewRegistry()

 requestsTotal := prometheus.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_requests_total",
   Help: "Total number of HTTP requests",
  },
  []string{"method", "path", "status"},
 )

 requestDuration := prometheus.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_duration_seconds",
   Help:    "HTTP request duration in seconds",
   Buckets: prometheus.DefBuckets,
  },
  []string{"method", "path"},
 )

 ordersProcessed := prometheus.NewCounterVec(
  prometheus.CounterOpts{
   Name: "orders_processed_total",
   Help: "Total number of orders processed",
  },
  []string{"status", "region"},
 )

 orderValue := prometheus.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "order_value_dollars",
   Help:    "Order value distribution",
   Buckets: []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
  },
  []string{"region"},
 )

 // Register metrics
 reg.MustRegister(requestsTotal, requestDuration, ordersProcessed, orderValue)

 return &PrometheusMetrics{
  Registry:        reg,
  RequestsTotal:   requestsTotal,
  RequestDuration: requestDuration,
  OrdersProcessed: ordersProcessed,
  OrderValue:      orderValue,
 }
}

// Handler returns HTTP handler for metrics endpoint
func (pm *PrometheusMetrics) Handler() http.Handler {
 return promhttp.HandlerFor(pm.Registry, promhttp.HandlerOpts{})
}

// RecordHTTPRequest records an HTTP request
func (pm *PrometheusMetrics) RecordHTTPRequest(method, path, status string, duration time.Duration) {
 pm.RequestsTotal.WithLabelValues(method, path, status).Inc()
 pm.RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// RecordOrder records an order
func (pm *PrometheusMetrics) RecordOrder(status, region string, value float64) {
 pm.OrdersProcessed.WithLabelValues(status, region).Inc()
 pm.OrderValue.WithLabelValues(region).Observe(value)
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Failure Taxonomy

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Cardinality Explosion** | Memory exhaustion | Series count alert | Label limits + Cardinality enforcement |
| **Metric Loss** | Blind spots | Scrape failure | Queue buffering + Retries |
| **Timestamp Drift** | Query anomalies | Clock skew | NTP + Monotonic clocks |
| **Hot Shard** | Query timeout | Query latency | Sharding + Rebalancing |
| **OOM on Query** | Service unavailability | Memory spike | Query limits + Timeout |
| **Label Leak** | PII exposure | Pattern scanning | Label allowlists |

---

## 6. Semantic Trade-off Analysis

### 6.1 Metric Type Selection

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    METRIC TYPE SELECTION GUIDE                               │
├─────────────────────┬─────────────────┬─────────────────┬───────────────────┤
│     Use Case        │     Counter     │     Gauge       │    Histogram      │
├─────────────────────┼─────────────────┼─────────────────┼───────────────────┤
│ Requests served     │ ✅ Perfect      │ ❌ Wrong        │ ❌ Wrong          │
│ Active connections  │ ❌ Wrong        │ ✅ Perfect      │ ❌ Wrong          │
│ Request latency     │ ❌ Wrong        │ ⚠️  Limited     │ ✅ Perfect        │
│ Temperature         │ ❌ Wrong        │ ✅ Perfect      │ ❌ Wrong          │
│ Queue depth         │ ❌ Wrong        │ ✅ Perfect      │ ❌ Wrong          │
│ Error rate          │ ✅ Can derive   │ ⚠️  Complex     │ ❌ Wrong          │
└─────────────────────┴─────────────────┴─────────────────┴───────────────────┘
```

---

## 7. References

1. Beyer, B., et al. (2016). *Site Reliability Engineering*. O'Reilly Media.
2. Brazil, B. (2018). *Prometheus: Up & Running*. O'Reilly Media.
3. Wilkinson, T. (2015). The RED Method. Weaveworks Blog.
4. Gregg, B. (2017). *Systems Performance*. Addison-Wesley.
5. OpenTelemetry Project. (2024). *Metrics Specification*. opentelemetry.io.
