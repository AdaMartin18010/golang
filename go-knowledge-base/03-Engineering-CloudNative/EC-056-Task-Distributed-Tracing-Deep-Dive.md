# 任务分布式追踪深入剖析 (Task Distributed Tracing Deep Dive)

> **分类**: 工程与云原生
> **标签**: #distributed-tracing #opentelemetry #observability #deep-dive

---

## 追踪模型架构

```go
// OpenTelemetry 完整追踪实现
package tracing

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

// TracerProvider 配置
type TracerConfig struct {
    ServiceName    string
    ServiceVersion string
    Environment    string
    ExporterType   string  // jaeger, zipkin, otlp
    SamplingRate   float64
}

func InitTracer(config TracerConfig) (*TracerProvider, error) {
    // 创建资源
    res, _ := resource.Merge(
        resource.Default(),
        resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(config.ServiceName),
            semconv.ServiceVersion(config.ServiceVersion),
            attribute.String("environment", config.Environment),
        ),
    )

    // 配置采样
    sampler := sdktrace.ParentBased(
        sdktrace.TraceIDRatioBased(config.SamplingRate),
    )

    // 创建导出器
    exporter, err := createExporter(config.ExporterType)
    if err != nil {
        return nil, err
    }

    // 创建 provider
    provider := sdktrace.NewTracerProvider(
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sampler),
        sdktrace.WithBatcher(exporter),
    )

    otel.SetTracerProvider(provider)

    return provider, nil
}
```

---

## 任务生命周期追踪

```go
// 完整任务追踪器
type TaskTracer struct {
    tracer trace.Tracer
}

func (tt *TaskTracer) TraceTaskLifecycle(ctx context.Context, task *Task) (context.Context, func()) {
    // 创建根 span
    ctx, span := tt.tracer.Start(ctx, "task.execute",
        trace.WithSpanKind(trace.SpanKindInternal),
        trace.WithAttributes(
            attribute.String("task.id", task.ID),
            attribute.String("task.type", task.Type),
            attribute.String("task.name", task.Name),
            attribute.Int("task.priority", task.Priority),
        ),
    )

    // 记录任务开始事件
    span.AddEvent("task.scheduled", trace.WithAttributes(
        attribute.String("scheduler", task.SchedulerID),
        attribute.Time("scheduled_at", task.ScheduledAt),
    ))

    // 返回清理函数
    return ctx, func() {
        span.End()
    }
}

func (tt *TaskTracer) TraceTaskStage(ctx context.Context, stage string, fn func(context.Context) error) error {
    ctx, span := tt.tracer.Start(ctx, fmt.Sprintf("task.stage.%s", stage),
        trace.WithAttributes(
            attribute.String("stage.name", stage),
        ),
    )
    defer span.End()

    // 记录开始时间
    start := time.Now()

    err := fn(ctx)

    // 记录结果
    duration := time.Since(start)
    span.SetAttributes(attribute.Int64("stage.duration_ms", duration.Milliseconds()))

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    }

    return err
}

func (tt *TaskTracer) TraceRetry(ctx context.Context, attempt int, err error) {
    span := trace.SpanFromContext(ctx)

    span.AddEvent("task.retry", trace.WithAttributes(
        attribute.Int("retry.attempt", attempt),
        attribute.String("retry.error", err.Error()),
    ))

    span.SetAttributes(attribute.Int("task.retry_count", attempt))
}
```

---

## 链路关联分析

```go
// 链路分析器
type TraceAnalyzer struct {
    storage TraceStorage
}

// 查找慢链路
func (ta *TraceAnalyzer) FindSlowTraces(ctx context.Context, threshold time.Duration, limit int) ([]TraceSummary, error) {
    return ta.storage.Query(ctx, TraceQuery{
        MinDuration: threshold,
        Limit:       limit,
        OrderBy:     "duration DESC",
    })
}

// 错误链路分析
func (ta *TraceAnalyzer) AnalyzeErrorTraces(ctx context.Context, window time.Duration) (*ErrorAnalysis, error) {
    // 获取错误链路
    errorTraces, _ := ta.storage.Query(ctx, TraceQuery{
        Status:    "error",
        StartTime: time.Now().Add(-window),
    })

    analysis := &ErrorAnalysis{
        TotalErrors: len(errorTraces),
        Patterns:    make(map[string]int),
    }

    for _, trace := range errorTraces {
        // 分析错误模式
        pattern := ta.classifyError(trace)
        analysis.Patterns[pattern]++

        // 提取错误来源
        source := ta.findErrorSource(trace)
        analysis.TopSources[source]++
    }

    return analysis, nil
}

func (ta *TraceAnalyzer) classifyError(t TraceSummary) string {
    // 基于错误类型和位置分类
    if t.ErrorType == "timeout" {
        return "timeout"
    }

    if strings.Contains(t.ErrorMessage, "connection refused") {
        return "connectivity"
    }

    if strings.Contains(t.ErrorMessage, "rate limit") {
        return "rate_limit"
    }

    return "other"
}

// 链路性能热力图
func (ta *TraceAnalyzer) GenerateHeatmap(ctx context.Context, service string, hours int) (*Heatmap, error) {
    // 按时间段聚合
    buckets := make(map[string][]time.Duration)

    traces, _ := ta.storage.Query(ctx, TraceQuery{
        Service:   service,
        StartTime: time.Now().Add(-time.Duration(hours) * time.Hour),
    })

    for _, trace := range traces {
        bucket := trace.StartTime.Format("2006-01-02 15:00")
        buckets[bucket] = append(buckets[bucket], trace.Duration)
    }

    // 计算每个 bucket 的 P99
    heatmap := &Heatmap{}
    for bucket, durations := range buckets {
        p99 := calculateP99(durations)
        heatmap.Cells = append(heatmap.Cells, HeatmapCell{
            Time:     bucket,
            Duration: p99,
            Color:    ta.durationToColor(p99),
        })
    }

    return heatmap, nil
}
```

---

## 采样策略高级配置

```go
// 自适应采样
type AdaptiveSampler struct {
    baseRate     float64
    slowThreshold time.Duration
    errorRate    float64
    mu           sync.RWMutex
}

func (as *AdaptiveSampler) ShouldSample(params sdktrace.SamplingParameters) sdktrace.SamplingResult {
    // 始终采样慢请求
    if as.isSlowOperation(params.Name) {
        return sdktrace.SamplingResult{
            Decision: sdktrace.RecordAndSample,
        }
    }

    // 始终采样错误
    if as.isErrorOperation(params.Name) {
        return sdktrace.SamplingResult{
            Decision: sdktrace.RecordAndSample,
        }
    }

    // 基于特征采样
    as.mu.RLock()
    rate := as.baseRate
    as.mu.RUnlock()

    // 根据优先级调整采样率
    if priority := as.getOperationPriority(params.Name); priority > 5 {
        rate = min(1.0, rate*2)
    }

    if rand.Float64() < rate {
        return sdktrace.SamplingResult{
            Decision: sdktrace.RecordAndSample,
        }
    }

    return sdktrace.SamplingResult{
        Decision: sdktrace.Drop,
    }
}

// 尾部采样 (Tail-based Sampling)
type TailSampler struct {
    buffer  *TraceBuffer
    policies []SamplingPolicy
}

type SamplingPolicy interface {
    ShouldSample(spans []ReadOnlySpan) bool
}

func (ts *TailSampler) ProcessCompletedTrace(traceID string, spans []ReadOnlySpan) {
    // 等待所有 span 到达
    completeSpans := ts.buffer.WaitForComplete(traceID, 5*time.Second)

    // 应用策略
    shouldSample := false
    for _, policy := range ts.policies {
        if policy.ShouldSample(completeSpans) {
            shouldSample = true
            break
        }
    }

    if shouldSample {
        ts.export(completeSpans)
    }
}

// 错误采样策略
type ErrorPolicy struct{}

func (ep *ErrorPolicy) ShouldSample(spans []ReadOnlySpan) bool {
    for _, span := range spans {
        if span.Status().Code == codes.Error {
            return true
        }
    }
    return false
}

// 慢请求采样策略
type SlowPolicy struct {
    threshold time.Duration
}

func (sp *SlowPolicy) ShouldSample(spans []ReadOnlySpan) bool {
    for _, span := range spans {
        if span.EndTime().Sub(span.StartTime()) > sp.threshold {
            return true
        }
    }
    return false
}
```
