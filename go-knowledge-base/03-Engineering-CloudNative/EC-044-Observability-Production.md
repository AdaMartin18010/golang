# EC-044: 可观测性生产实践 (Observability in Production)

> **维度**: Engineering CloudNative
> **级别**: S (20+ KB)
> **标签**: #observability #metrics #logging #tracing #monitoring
> **相关**: EC-006, EC-032, EC-080

---

## 整合说明

本文档整合：

- `06-Distributed-Tracing.md` (已重命名为 EC-006)
- `22-Context-Aware-Logging.md` (5.8 KB)
- `26-Task-Monitoring-Alerting.md` (7.3 KB)
- `32-Task-Observability.md` (5.9 KB)
- `56-Task-Distributed-Tracing-Deep-Dive.md` (8.5 KB)
- `60-OpenTelemetry-Distributed-Tracing-Production.md` (18 KB)
- `80-Observability-Metrics-Integration.md` (20 KB)

---

## 三大支柱

```
┌─────────────────────────────────────────────────────────────────┐
│                     Observability Pillars                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│    Metrics          Logs            Traces                      │
│    ───────          ────            ──────                      │
│                                                                  │
│    ┌────────┐      ┌────────┐      ┌────────┐                 │
│    │Counter │      │Structured│     │Span    │                 │
│    │Gauge   │      │Text    │      │Context │                 │
│    │Histogram│     │JSON    │      │Trace   │                 │
│    └────────┘      └────────┘      └────────┘                 │
│                                                                  │
│    When?            What?           Where?                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 指标 (Metrics)

```go
package metrics

import (
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
)

// 任务调度指标
var (
 TasksSubmitted = promauto.NewCounterVec(prometheus.CounterOpts{
  Name: "tasks_submitted_total",
  Help: "Total number of tasks submitted",
 }, []string{"type", "priority"})

 TasksCompleted = promauto.NewCounterVec(prometheus.CounterOpts{
  Name: "tasks_completed_total",
  Help: "Total number of tasks completed",
 }, []string{"type", "status"})

 TaskDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
  Name:    "task_duration_seconds",
  Help:    "Task execution duration",
  Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
 }, []string{"type"})

 QueueDepth = promauto.NewGaugeVec(prometheus.GaugeOpts{
  Name: "queue_depth",
  Help: "Current queue depth",
 }, []string{"queue_name"})

 ActiveWorkers = promauto.NewGauge(prometheus.GaugeOpts{
  Name: "active_workers",
  Help: "Number of active workers",
 })
)

// 使用示例
func (s *Scheduler) Submit(task *Task) error {
 TasksSubmitted.WithLabelValues(task.Type, string(task.Priority)).Inc()

 start := time.Now()
 err := s.doSubmit(task)

 TaskDuration.WithLabelValues(task.Type).Observe(time.Since(start).Seconds())

 return err
}
```

---

## 日志 (Logging)

```go
package logging

import (
 "context"
 "go.uber.org/zap"
)

// ContextualLogger 带上下文的日志记录器
type ContextualLogger struct {
 logger *zap.Logger
}

func (l *ContextualLogger) WithContext(ctx context.Context) *zap.Logger {
 fields := []zap.Field{}

 if traceID := GetTraceID(ctx); traceID != "" {
  fields = append(fields, zap.String("trace_id", traceID))
 }

 if spanID := GetSpanID(ctx); spanID != "" {
  fields = append(fields, zap.String("span_id", spanID))
 }

 if taskID := GetTaskID(ctx); taskID != "" {
  fields = append(fields, zap.String("task_id", taskID))
 }

 return l.logger.With(fields...)
}

// 结构化日志示例
func (e *TaskExecutor) Execute(ctx context.Context, task *Task) {
 logger := logging.WithContext(ctx)

 logger.Info("task execution started",
  zap.String("task_type", task.Type),
  zap.Int("priority", int(task.Priority)),
  zap.String("worker_id", task.WorkerID),
 )

 defer func() {
  if err := recover(); err != nil {
   logger.Error("task execution panicked",
    zap.Any("error", err),
    zap.Stack("stacktrace"),
   )
  }
 }()

 err := e.runTask(ctx, task)

 if err != nil {
  logger.Error("task execution failed",
   zap.Error(err),
   zap.Duration("duration", time.Since(start)),
  )
 } else {
  logger.Info("task execution completed",
   zap.Duration("duration", time.Since(start)),
  )
 }
}
```

---

## 追踪 (Tracing)

```go
package tracing

import (
 "context"
 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("task-scheduler")

// StartSpan 开始追踪跨度
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
 ctx, span := tracer.Start(ctx, name)
 span.SetAttributes(attrs...)
 return ctx, span
}

// 任务执行追踪
func (e *TaskExecutor) Execute(ctx context.Context, task *Task) error {
 ctx, span := StartSpan(ctx, "task.execute",
  attribute.String("task.id", task.ID),
  attribute.String("task.type", task.Type),
  attribute.Int("task.priority", int(task.Priority)),
 )
 defer span.End()

 // 预处理
 ctx, preprocessSpan := StartSpan(ctx, "task.preprocess")
 if err := e.preprocess(ctx, task); err != nil {
  preprocessSpan.RecordError(err)
  span.SetStatus(codes.Error, "preprocess failed")
  return err
 }
 preprocessSpan.End()

 // 执行
 ctx, execSpan := StartSpan(ctx, "task.run")
 err := e.run(ctx, task)
 execSpan.End()

 if err != nil {
  span.RecordError(err)
  span.SetStatus(codes.Error, err.Error())
 }

 return err
}
```

---

## 告警规则

```yaml
# Prometheus Alert Rules
groups:
  - name: task-scheduler
    rules:
      - alert: HighTaskFailureRate
        expr: |
          (
            sum(rate(tasks_completed_total{status="failed"}[5m]))
            /
            sum(rate(tasks_completed_total[5m]))
          ) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High task failure rate"
          description: "Task failure rate is above 10%"

      - alert: QueueDepthHigh
        expr: queue_depth > 1000
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Queue depth too high"
          description: "Queue {{ $labels.queue_name }} has {{ $value }} pending tasks"

      - alert: WorkerDown
        expr: active_workers < 3
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Insufficient workers"
```

---

## Dashboard

```json
{
  "title": "Task Scheduler Dashboard",
  "panels": [
    {
      "title": "Task Rate",
      "targets": [
        {
          "expr": "rate(tasks_submitted_total[5m])",
          "legendFormat": "Submitted"
        },
        {
          "expr": "rate(tasks_completed_total[5m])",
          "legendFormat": "Completed"
        }
      ]
    },
    {
      "title": "Queue Depth",
      "targets": [
        {
          "expr": "queue_depth",
          "legendFormat": "{{ queue_name }}"
        }
      ]
    },
    {
      "title": "Task Duration",
      "targets": [
        {
          "expr": "histogram_quantile(0.95, rate(task_duration_seconds_bucket[5m]))",
          "legendFormat": "p95"
        },
        {
          "expr": "histogram_quantile(0.99, rate(task_duration_seconds_bucket[5m]))",
          "legendFormat": "p99"
        }
      ]
    }
  ]
}
```
