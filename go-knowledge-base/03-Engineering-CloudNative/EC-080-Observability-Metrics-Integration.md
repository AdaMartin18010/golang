# 可观测性与指标集成 (Observability & Metrics Integration)

> **分类**: 工程与云原生
> **标签**: #observability #metrics #prometheus #grafana
> **参考**: OpenTelemetry Metrics, Prometheus Best Practices

---

## 指标架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Observability Metrics Architecture                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Application Metrics                              │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  Counter    │  │   Gauge     │  │ Histogram   │  │  Summary    │ │   │
│  │  │  (tasks_    │  │  (active_   │  │ (duration_  │  │  (request_  │ │   │
│  │  │  processed) │  │  workers)   │  │  seconds)   │  │  size)      │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │   RED       │  │   USE       │  │  Custom     │                  │   │
│  │  │ (Rate/Err/  │  │(Util/Sat/  │  │  Business   │                  │   │
│  │  │  Duration)  │  │  Errors)    │  │  Metrics    │                  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Metric Collection                                 │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  Prometheus │  │ OpenTelemetry│  │   StatsD    │  │   Custom    │ │   │
│  │  │  Client     │  │   SDK       │  │  Client     │  │   Bridge    │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Visualization & Alerting                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Grafana   │  │  Prometheus │  │   Jaeger    │  │   PagerDuty │ │   │
│  │  │  Dashboard  │  │   Alerts    │  │   Traces    │  │   Alerts    │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Prometheus 指标实现

```go
package metrics

import (
    "context"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// TaskMetrics 任务相关指标
type TaskMetrics struct {
    // 计数器
    TasksSubmitted   prometheus.Counter
    TasksCompleted   prometheus.Counter
    TasksFailed      prometheus.Counter
    TasksRetried     prometheus.Counter

    // 仪表盘
    ActiveWorkers    prometheus.Gauge
    PendingTasks     prometheus.Gauge
    QueueDepth       prometheus.Gauge

    // 直方图
    TaskDuration     prometheus.Histogram
    TaskQueueTime    prometheus.Histogram

    // 带标签的计数器
    TasksByType      *prometheus.CounterVec
    ErrorsByType     *prometheus.CounterVec
}

// NewTaskMetrics 创建任务指标
func NewTaskMetrics() *TaskMetrics {
    return &TaskMetrics{
        TasksSubmitted: promauto.NewCounter(prometheus.CounterOpts{
            Name: "scheduler_tasks_submitted_total",
            Help: "Total number of tasks submitted",
        }),
        TasksCompleted: promauto.NewCounter(prometheus.CounterOpts{
            Name: "scheduler_tasks_completed_total",
            Help: "Total number of tasks completed",
        }),
        TasksFailed: promauto.NewCounter(prometheus.CounterOpts{
            Name: "scheduler_tasks_failed_total",
            Help: "Total number of tasks failed",
        }),
        TasksRetried: promauto.NewCounter(prometheus.CounterOpts{
            Name: "scheduler_tasks_retried_total",
            Help: "Total number of task retries",
        }),
        ActiveWorkers: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "scheduler_active_workers",
            Help: "Number of active workers",
        }),
        PendingTasks: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "scheduler_pending_tasks",
            Help: "Number of pending tasks",
        }),
        QueueDepth: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "scheduler_queue_depth",
            Help: "Current queue depth",
        }),
        TaskDuration: promauto.NewHistogram(prometheus.HistogramOpts{
            Name:    "scheduler_task_duration_seconds",
            Help:    "Task execution duration in seconds",
            Buckets: prometheus.DefBuckets,
        }),
        TaskQueueTime: promauto.NewHistogram(prometheus.HistogramOpts{
            Name:    "scheduler_task_queue_time_seconds",
            Help:    "Time spent waiting in queue",
            Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
        }),
        TasksByType: promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "scheduler_tasks_by_type_total",
            Help: "Total number of tasks by type",
        }, []string{"type", "status"}),
        ErrorsByType: promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "scheduler_errors_total",
            Help: "Total number of errors by type",
        }, []string{"type"}),
    }
}

// RecordTaskSubmitted 记录任务提交
func (m *TaskMetrics) RecordTaskSubmitted() {
    m.TasksSubmitted.Inc()
}

// RecordTaskCompleted 记录任务完成
func (m *TaskMetrics) RecordTaskCompleted(duration time.Duration) {
    m.TasksCompleted.Inc()
    m.TaskDuration.Observe(duration.Seconds())
}

// RecordTaskFailed 记录任务失败
func (m *TaskMetrics) RecordTaskFailed(taskType string) {
    m.TasksFailed.Inc()
    m.TasksByType.WithLabelValues(taskType, "failed").Inc()
}

// RecordTaskRetry 记录任务重试
func (m *TaskMetrics) RecordTaskRetry() {
    m.TasksRetried.Inc()
}

// SetActiveWorkers 设置活跃工作线程数
func (m *TaskMetrics) SetActiveWorkers(n int) {
    m.ActiveWorkers.Set(float64(n))
}

// SetQueueDepth 设置队列深度
func (m *TaskMetrics) SetQueueDepth(n int) {
    m.QueueDepth.Set(float64(n))
}

// HTTPMetrics HTTP 指标
type HTTPMetrics struct {
    RequestsTotal   *prometheus.CounterVec
    RequestDuration *prometheus.HistogramVec
    RequestSize     *prometheus.SummaryVec
    ResponseSize    *prometheus.SummaryVec
    ActiveRequests  prometheus.Gauge
}

// NewHTTPMetrics 创建 HTTP 指标
func NewHTTPMetrics() *HTTPMetrics {
    return &HTTPMetrics{
        RequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        }, []string{"method", "path", "status"}),
        RequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        }, []string{"method", "path"}),
        RequestSize: promauto.NewSummaryVec(prometheus.SummaryOpts{
            Name:       "http_request_size_bytes",
            Help:       "HTTP request size in bytes",
            Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
        }, []string{"method", "path"}),
        ResponseSize: promauto.NewSummaryVec(prometheus.SummaryOpts{
            Name:       "http_response_size_bytes",
            Help:       "HTTP response size in bytes",
            Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
        }, []string{"method", "path"}),
        ActiveRequests: promauto.NewGauge(prometheus.GaugeOpts{
            Name: "http_active_requests",
            Help: "Number of active HTTP requests",
        }),
    }
}

// MetricsMiddleware HTTP 指标中间件
func (m *HTTPMetrics) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        m.ActiveRequests.Inc()
        defer m.ActiveRequests.Dec()

        // 包装 ResponseWriter 以捕获状态码和大小
        rw := &responseRecorder{ResponseWriter: w}

        next.ServeHTTP(rw, r)

        duration := time.Since(start)

        // 记录指标
        m.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, string(rw.statusCode)).Inc()
        m.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
        m.ResponseSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(rw.written))
    })
}

// StartMetricsServer 启动指标服务器
func StartMetricsServer(addr string) *http.Server {
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())

    server := &http.Server{
        Addr:    addr,
        Handler: mux,
    }

    go server.ListenAndServe()

    return server
}
```

---

## OpenTelemetry 集成

```go
package metrics

import (
    "context"
    "time"

    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/metric/global"
    "go.opentelemetry.io/otel/sdk/metric/export/prometheus"
    controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
    processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
    "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// OTelMetrics OpenTelemetry 指标
type OTelMetrics struct {
    meter        metric.Meter

    // 计数器
    taskCounter  metric.Int64Counter
    errorCounter metric.Int64Counter

    // 仪表盘
    workerGauge  metric.Int64UpDownCounter
    queueGauge   metric.Int64ObservableGauge

    // 直方图
    durationHist metric.Float64Histogram
}

// NewOTelMetrics 创建 OpenTelemetry 指标
func NewOTelMetrics() (*OTelMetrics, error) {
    meter := global.Meter("scheduler")

    taskCounter, err := meter.Int64Counter(
        "scheduler.tasks",
        metric.WithDescription("Number of tasks"),
    )
    if err != nil {
        return nil, err
    }

    errorCounter, err := meter.Int64Counter(
        "scheduler.errors",
        metric.WithDescription("Number of errors"),
    )
    if err != nil {
        return nil, err
    }

    workerGauge, err := meter.Int64UpDownCounter(
        "scheduler.workers",
        metric.WithDescription("Number of workers"),
    )
    if err != nil {
        return nil, err
    }

    queueGauge, err := meter.Int64ObservableGauge(
        "scheduler.queue.depth",
        metric.WithDescription("Queue depth"),
    )
    if err != nil {
        return nil, err
    }

    durationHist, err := meter.Float64Histogram(
        "scheduler.task.duration",
        metric.WithDescription("Task duration"),
    )
    if err != nil {
        return nil, err
    }

    return &OTelMetrics{
        meter:        meter,
        taskCounter:  taskCounter,
        errorCounter: errorCounter,
        workerGauge:  workerGauge,
        queueGauge:   queueGauge,
        durationHist: durationHist,
    }, nil
}

// RecordTask 记录任务
func (m *OTelMetrics) RecordTask(ctx context.Context, taskType string, status string) {
    m.taskCounter.Add(ctx, 1,
        attribute.String("type", taskType),
        attribute.String("status", status),
    )
}

// RecordError 记录错误
func (m *OTelMetrics) RecordError(ctx context.Context, errorType string) {
    m.errorCounter.Add(ctx, 1,
        attribute.String("type", errorType),
    )
}

// RecordDuration 记录持续时间
func (m *OTelMetrics) RecordDuration(ctx context.Context, duration time.Duration) {
    m.durationHist.Record(ctx, duration.Seconds())
}

// IncWorkers 增加工作线程数
func (m *OTelMetrics) IncWorkers(ctx context.Context) {
    m.workerGauge.Add(ctx, 1)
}

// DecWorkers 减少工作线程数
func (m *OTelMetrics) DecWorkers(ctx context.Context) {
    m.workerGauge.Add(ctx, -1)
}

// SetupOTelPrometheusBridge 设置 OTel-Prometheus 桥接
func SetupOTelPrometheusBridge() (*prometheus.Exporter, error) {
    config := prometheus.Config{}

    ctrl := controller.New(
        processor.New(
            simple.NewWithHistogramDistribution(),
            config,
        ),
    )

    exporter, err := prometheus.New(config, ctrl)
    if err != nil {
        return nil, err
    }

    global.SetMeterProvider(ctrl)

    return exporter, nil
}
```

---

## 追踪集成

```go
package metrics

import (
    "context"
    "fmt"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

// Tracing 追踪工具
type Tracing struct {
    tracer trace.Tracer
}

// NewTracing 创建追踪工具
func NewTracing(tracerName string) *Tracing {
    return &Tracing{
        tracer: otel.Tracer(tracerName),
    }
}

// StartSpan 开始 Span
func (t *Tracing) StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
    return t.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// TraceTask 追踪任务执行
func (t *Tracing) TraceTask(ctx context.Context, taskType, taskID string, fn func(context.Context) error) error {
    ctx, span := t.StartSpan(ctx, "execute-task",
        attribute.String("task.type", taskType),
        attribute.String("task.id", taskID),
    )
    defer span.End()

    err := fn(ctx)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    } else {
        span.SetStatus(codes.Ok, "success")
    }

    return err
}

// TraceHTTP 追踪 HTTP 请求
func (t *Tracing) TraceHTTP(ctx context.Context, method, path string, fn func(context.Context) error) error {
    ctx, span := t.StartSpan(ctx, fmt.Sprintf("HTTP %s %s", method, path),
        attribute.String("http.method", method),
        attribute.String("http.path", path),
    )
    defer span.End()

    return fn(ctx)
}

// TraceDatabase 追踪数据库操作
func (t *Tracing) TraceDatabase(ctx context.Context, operation, table string, fn func(context.Context) error) error {
    ctx, span := t.StartSpan(ctx, fmt.Sprintf("DB %s", operation),
        attribute.String("db.operation", operation),
        attribute.String("db.table", table),
    )
    defer span.End()

    return fn(ctx)
}
```

---

## 健康检查指标

```go
package metrics

import (
    "context"
    "sync"
    "time"
)

// HealthChecker 健康检查器
type HealthChecker struct {
    checks map[string]HealthCheck
    mu     sync.RWMutex
}

// HealthCheck 健康检查函数
type HealthCheck func(ctx context.Context) error

// NewHealthChecker 创建健康检查器
func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        checks: make(map[string]HealthCheck),
    }
}

// Register 注册健康检查
func (h *HealthChecker) Register(name string, check HealthCheck) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.checks[name] = check
}

// CheckAll 执行所有健康检查
func (h *HealthChecker) CheckAll(ctx context.Context) map[string]error {
    h.mu.RLock()
    checks := make(map[string]HealthCheck, len(h.checks))
    for k, v := range h.checks {
        checks[k] = v
    }
    h.mu.RUnlock()

    results := make(map[string]error)

    for name, check := range checks {
        ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
        results[name] = check(ctx)
        cancel()
    }

    return results
}

// IsHealthy 是否健康
func (h *HealthChecker) IsHealthy(ctx context.Context) bool {
    results := h.CheckAll(ctx)
    for _, err := range results {
        if err != nil {
            return false
        }
    }
    return true
}

// CustomMetrics 自定义业务指标
type CustomMetrics struct {
    // 可以根据业务需求定义
}

// RecordBusinessMetric 记录业务指标
func (c *CustomMetrics) RecordBusinessMetric(name string, value float64, labels map[string]string) {
    // 实现业务指标记录
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "metrics"
)

func main() {
    // 创建指标
    taskMetrics := metrics.NewTaskMetrics()
    httpMetrics := metrics.NewHTTPMetrics()

    // 启动指标服务器
    metrics.StartMetricsServer(":9090")

    // 创建 HTTP 服务器
    mux := http.NewServeMux()
    mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
        // 模拟任务处理
        taskMetrics.RecordTaskSubmitted()

        start := time.Now()
        time.Sleep(100 * time.Millisecond)

        taskMetrics.RecordTaskCompleted(time.Since(start))

        w.Write([]byte("Task processed"))
    })

    // 应用指标中间件
    handler := httpMetrics.Middleware(mux)

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", handler)
}
```
