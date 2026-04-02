# Observability-Driven Development (ODD)

> **分类**: 工程与云原生
> **标签**: #observability #odd #monitoring #telemetry #sre
> **参考**: Google SRE, OpenTelemetry, Site Reliability Engineering

---

## 1. Formal Definition

### 1.1 What is Observability-Driven Development?

Observability-Driven Development (ODD) is a software engineering methodology that treats observability as a first-class citizen throughout the entire software development lifecycle. It extends Test-Driven Development (TDD) by asserting that a system component is not complete until it is observable in production.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Observability-Driven Development Cycle                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐             │
│   │  Design  │───→│Implement │───→│ Observe  │───→│ Validate │             │
│   │          │    │ + Instrument│  │          │    │          │             │
│   └──────────┘    └──────────┘    └────┬─────┘    └────┬─────┘             │
│        ↑                               │               │                   │
│        │                               ↓               │                   │
│        │                          ┌──────────┐         │                   │
│        └──────────────────────────│  Learn   │←────────┘                   │
│                                   │  & Adapt │                             │
│                                   └──────────┘                             │
│                                                                             │
│   Key Principle: "If you can't observe it, you can't validate it"          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 The Three Pillars of Observability

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Three Pillars of Observability                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐            │
│   │     METRICS     │  │      LOGS       │  │     TRACES      │            │
│   │                 │  │                 │  │                 │            │
│   │  Quantitative   │  │   Qualitative   │  │   Transaction   │            │
│   │  measurements   │  │   event records │  │   lifecycle     │            │
│   │  over time      │  │   with context  │  │   across svcs   │            │
│   │                 │  │                 │  │                 │            │
│   │  • Counters     │  │  • Structured   │  │  • Spans        │            │
│   │  • Gauges       │  │  • Timestamped  │  │  • Relationships│            │
│   │  • Histograms   │  │  • Contextual   │  │  • Service map  │            │
│   │                 │  │                 │  │                 │            │
│   │  Question:      │  │  Question:      │  │  Question:      │            │
│   │  "What is the   │  │  "Why is it     │  │  "Where is the  │            │
│   │   error rate?"  │  │   happening?"   │  │   bottleneck?"  │            │
│   └────────┬────────┘  └────────┬────────┘  └────────┬────────┘            │
│            │                    │                    │                     │
│            └────────────────────┼────────────────────┘                     │
│                                 ↓                                         │
│                    ┌─────────────────────────┐                            │
│                    │      CORRELATION        │                            │
│                    │   (TraceID, SpanID,     │                            │
│                    │    Timestamp, Service)  │                            │
│                    └─────────────────────────┘                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 Structured Logging Pattern

```go
package observability

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "sync"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// LogLevel represents the severity of a log entry
type LogLevel int

const (
    DebugLevel LogLevel = iota
    InfoLevel
    WarnLevel
    ErrorLevel
    FatalLevel
)

// Logger is a production-ready structured logger
type Logger struct {
    zapLogger *zap.Logger
    service   string
    version   string
    mu        sync.RWMutex
    fields    map[string]interface{}
}

// LoggerConfig configures the logger
type LoggerConfig struct {
    Service     string
    Version     string
    Environment string
    Level       LogLevel
    Output      io.Writer
    EnableJSON  bool
}

// NewLogger creates a production-ready logger
func NewLogger(cfg LoggerConfig) (*Logger, error) {
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "timestamp",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        FunctionKey:    zapcore.OmitKey,
        MessageKey:     "message",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    level := zapcore.InfoLevel
    switch cfg.Level {
    case DebugLevel:
        level = zapcore.DebugLevel
    case WarnLevel:
        level = zapcore.WarnLevel
    case ErrorLevel:
        level = zapcore.ErrorLevel
    case FatalLevel:
        level = zapcore.FatalLevel
    }

    var encoder zapcore.Encoder
    if cfg.EnableJSON {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }

    output := cfg.Output
    if output == nil {
        output = os.Stdout
    }

    core := zapcore.NewCore(encoder, zapcore.AddSync(output), level)
    zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

    // Add default fields
    zapLogger = zapLogger.With(
        zap.String("service", cfg.Service),
        zap.String("version", cfg.Version),
        zap.String("environment", cfg.Environment),
    )

    return &Logger{
        zapLogger: zapLogger,
        service:   cfg.Service,
        version:   cfg.Version,
        fields:    make(map[string]interface{}),
    }, nil
}

// WithContext creates a logger with context fields
func (l *Logger) WithContext(ctx context.Context) *Logger {
    // Extract trace context
    traceID, _ := ctx.Value("trace_id").(string)
    spanID, _ := ctx.Value("span_id").(string)

    newLogger := &Logger{
        zapLogger: l.zapLogger,
        service:   l.service,
        version:   l.version,
        fields:    make(map[string]interface{}),
    }

    // Copy existing fields
    l.mu.RLock()
    for k, v := range l.fields {
        newLogger.fields[k] = v
    }
    l.mu.RUnlock()

    // Add trace context
    if traceID != "" {
        newLogger.fields["trace_id"] = traceID
    }
    if spanID != "" {
        newLogger.fields["span_id"] = spanID
    }

    newLogger.zapLogger = l.zapLogger.With(
        zap.String("trace_id", traceID),
        zap.String("span_id", spanID),
    )

    return newLogger
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
    newLogger := &Logger{
        zapLogger: l.zapLogger.With(zap.Any(key, value)),
        service:   l.service,
        version:   l.version,
        fields:    make(map[string]interface{}),
    }

    l.mu.RLock()
    for k, v := range l.fields {
        newLogger.fields[k] = v
    }
    l.mu.RUnlock()

    newLogger.fields[key] = value
    return newLogger
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...zap.Field) {
    l.zapLogger.Info(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...zap.Field) {
    allFields := append([]zap.Field{zap.Error(err)}, fields...)
    l.zapLogger.Error(msg, allFields...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
    l.zapLogger.Debug(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
    l.zapLogger.Warn(msg, fields...)
}

// Fatal logs a fatal message
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
    l.zapLogger.Fatal(msg, fields...)
}

// StructuredLogEntry represents a structured log entry
type StructuredLogEntry struct {
    Timestamp   time.Time              `json:"timestamp"`
    Level       string                 `json:"level"`
    Service     string                 `json:"service"`
    Version     string                 `json:"version"`
    TraceID     string                 `json:"trace_id,omitempty"`
    SpanID      string                 `json:"span_id,omitempty"`
    Message     string                 `json:"message"`
    Error       string                 `json:"error,omitempty"`
    Fields      map[string]interface{} `json:"fields,omitempty"`
    Caller      string                 `json:"caller"`
    Duration    time.Duration          `json:"duration,omitempty"`
}

// LogRequest logs an HTTP request with structured fields
func (l *Logger) LogRequest(reqID, method, path, clientIP string, statusCode int, duration time.Duration, err error) {
    entry := StructuredLogEntry{
        Timestamp: time.Now().UTC(),
        Level:     "info",
        Service:   l.service,
        Version:   l.version,
        Message:   fmt.Sprintf("HTTP %s %s", method, path),
        Duration:  duration,
        Fields: map[string]interface{}{
            "request_id":  reqID,
            "method":      method,
            "path":        path,
            "client_ip":   clientIP,
            "status_code": statusCode,
        },
    }

    if err != nil {
        entry.Level = "error"
        entry.Error = err.Error()
    }

    if statusCode >= 400 {
        entry.Level = "warn"
    }
    if statusCode >= 500 {
        entry.Level = "error"
    }

    data, _ := json.Marshal(entry)
    fmt.Println(string(data))
}
```

### 2.2 Metrics Collection Pattern

```go
package observability

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsCollector provides production-ready metrics collection
type MetricsCollector struct {
    namespace   string
    subsystem   string
    counters    map[string]prometheus.Counter
    gauges      map[string]prometheus.Gauge
    histograms  map[string]prometheus.Histogram
    summaries   map[string]prometheus.Summary
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(namespace, subsystem string) *MetricsCollector {
    return &MetricsCollector{
        namespace:  namespace,
        subsystem:  subsystem,
        counters:   make(map[string]prometheus.Counter),
        gauges:     make(map[string]prometheus.Gauge),
        histograms: make(map[string]prometheus.Histogram),
        summaries:  make(map[string]prometheus.Summary),
    }
}

// RegisterCounter registers a new counter metric
func (m *MetricsCollector) RegisterCounter(name, help string, labels ...string) {
    counter := promauto.NewCounterVec(prometheus.CounterOpts{
        Namespace: m.namespace,
        Subsystem: m.subsystem,
        Name:      name,
        Help:      help,
    }, labels)

    m.counters[name] = counter.WithLabelValues()
}

// RegisterHistogram registers a new histogram metric
func (m *MetricsCollector) RegisterHistogram(name, help string, buckets []float64, labels ...string) {
    if buckets == nil {
        buckets = prometheus.DefBuckets
    }

    histogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
        Namespace: m.namespace,
        Subsystem: m.subsystem,
        Name:      name,
        Help:      help,
        Buckets:   buckets,
    }, labels)

    m.histograms[name] = histogram.WithLabelValues()
}

// RegisterGauge registers a new gauge metric
func (m *MetricsCollector) RegisterGauge(name, help string, labels ...string) {
    gauge := promauto.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: m.namespace,
        Subsystem: m.subsystem,
        Name:      name,
        Help:      help,
    }, labels)

    m.gauges[name] = gauge.WithLabelValues()
}

// IncrementCounter increments a counter with labels
func (m *MetricsCollector) IncrementCounter(name string, labelValues ...string) {
    if counterVec, ok := prometheus.DefaultRegisterer.(*prometheus.Registry); ok {
        // Handle with labels
        if vec, err := m.getCounterVec(name); err == nil {
            vec.WithLabelValues(labelValues...).Inc()
        }
    }
}

// getCounterVec retrieves a counter vector by name
func (m *MetricsCollector) getCounterVec(name string) (*prometheus.CounterVec, error) {
    // This is a simplified version - in production use a proper registry
    return nil, fmt.Errorf("not implemented")
}

// RecordDuration records a duration in a histogram
func (m *MetricsCollector) RecordDuration(name string, duration time.Duration, labelValues ...string) {
    if hist, exists := m.histograms[name]; exists {
        hist.Observe(duration.Seconds())
    }
}

// SetGauge sets a gauge value
func (m *MetricsCollector) SetGauge(name string, value float64, labelValues ...string) {
    if gauge, exists := m.gauges[name]; exists {
        gauge.Set(value)
    }
}

// HTTPMiddleware returns an HTTP middleware for metrics collection
func (m *MetricsCollector) HTTPMiddleware(next http.Handler) http.Handler {
    // Register standard HTTP metrics
    requestDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
        Namespace: m.namespace,
        Subsystem: "http",
        Name:      "request_duration_seconds",
        Help:      "HTTP request duration in seconds",
        Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
    }, []string{"method", "path", "status"})

    requestCount := promauto.NewCounterVec(prometheus.CounterOpts{
        Namespace: m.namespace,
        Subsystem: "http",
        Name:      "requests_total",
        Help:      "Total HTTP requests",
    }, []string{"method", "path", "status"})

    requestSize := promauto.NewHistogramVec(prometheus.HistogramOpts{
        Namespace: m.namespace,
        Subsystem: "http",
        Name:      "request_size_bytes",
        Help:      "HTTP request size in bytes",
        Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
    }, []string{"method", "path"})

    activeRequests := promauto.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: m.namespace,
        Subsystem: "http",
        Name:      "active_requests",
        Help:      "Number of active HTTP requests",
    }, []string{"method", "path"})

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        path := r.URL.Path
        method := r.Method

        // Track active requests
        activeRequests.WithLabelValues(method, path).Inc()
        defer activeRequests.WithLabelValues(method, path).Dec()

        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(wrapped, r)

        duration := time.Since(start)
        status := fmt.Sprintf("%d", wrapped.statusCode)

        requestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())
        requestCount.WithLabelValues(method, path, status).Inc()
        requestSize.WithLabelValues(method, path).Observe(float64(r.ContentLength))
    })
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// MetricsServer starts an HTTP server for Prometheus scraping
func (m *MetricsCollector) MetricsServer(addr string) *http.Server {
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())

    // Health check endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"healthy"}`))
    })

    return &http.Server{
        Addr:         addr,
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
}

// Timer is a helper for timing operations
type Timer struct {
    start    time.Time
    hist     prometheus.Observer
    labels   []string
}

// NewTimer creates a new timer
func NewTimer(hist prometheus.Observer) *Timer {
    return &Timer{
        start: time.Now(),
        hist:  hist,
    }
}

// ObserveDuration observes the duration and returns it
func (t *Timer) ObserveDuration() time.Duration {
    duration := time.Since(t.start)
    if t.hist != nil {
        t.hist.Observe(duration.Seconds())
    }
    return duration
}
```

### 2.3 Distributed Tracing Pattern

```go
package observability

import (
    "context"
    "fmt"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
)

// TracerConfig configures the tracer
type TracerConfig struct {
    ServiceName    string
    ServiceVersion string
    Environment    string
    Endpoint       string
    ExporterType   string // "jaeger", "otlp", "stdout"
    SamplingRate   float64
}

// TracerProvider wraps OpenTelemetry tracer provider
type TracerProvider struct {
    provider *sdktrace.TracerProvider
    service  string
}

// NewTracerProvider creates a new tracer provider
func NewTracerProvider(cfg TracerConfig) (*TracerProvider, error) {
    var exporter sdktrace.SpanExporter
    var err error

    switch cfg.ExporterType {
    case "jaeger":
        exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(
            jaeger.WithEndpoint(cfg.Endpoint),
        ))
    case "otlp":
        client := otlptracegrpc.NewClient(
            otlptracegrpc.WithEndpoint(cfg.Endpoint),
        )
        exporter, err = otlptrace.New(context.Background(), client)
    default:
        // Use stdout exporter for debugging
        exporter, err = NewStdoutExporter()
    }

    if err != nil {
        return nil, fmt.Errorf("failed to create exporter: %w", err)
    }

    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String(cfg.ServiceName),
            semconv.ServiceVersionKey.String(cfg.ServiceVersion),
            semconv.DeploymentEnvironmentKey.String(cfg.Environment),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }

    // Configure sampling
    sampler := sdktrace.TraceIDRatioBased(cfg.SamplingRate)
    if cfg.SamplingRate >= 1.0 {
        sampler = sdktrace.AlwaysSample()
    } else if cfg.SamplingRate <= 0 {
        sampler = sdktrace.NeverSample()
    }

    provider := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sampler),
    )

    // Set as global provider
    otel.SetTracerProvider(provider)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return &TracerProvider{
        provider: provider,
        service:  cfg.ServiceName,
    }, nil
}

// Tracer returns a tracer for the service
func (tp *TracerProvider) Tracer(name string) trace.Tracer {
    return tp.provider.Tracer(name)
}

// Shutdown gracefully shuts down the provider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
    return tp.provider.Shutdown(ctx)
}

// SpanContext wraps OpenTelemetry span for easier use
type SpanContext struct {
    context.Context
    span trace.Span
}

// StartSpan starts a new span
func StartSpan(ctx context.Context, tracer trace.Tracer, name string, opts ...trace.SpanStartOption) (SpanContext, func()) {
    ctx, span := tracer.Start(ctx, name, opts...)
    return SpanContext{Context: ctx, span: span}, func() { span.End() }
}

// SetAttributes sets attributes on the span
func (sc SpanContext) SetAttributes(attrs ...attribute.KeyValue) {
    sc.span.SetAttributes(attrs...)
}

// RecordError records an error on the span
func (sc SpanContext) RecordError(err error, opts ...trace.EventOption) {
    sc.span.RecordError(err, opts...)
    sc.span.SetStatus(codes.Error, err.Error())
}

// AddEvent adds an event to the span
func (sc SpanContext) AddEvent(name string, attrs ...attribute.KeyValue) {
    sc.span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetStatus sets the span status
func (sc SpanContext) SetStatus(code codes.Code, description string) {
    sc.span.SetStatus(code, description)
}

// TraceID returns the trace ID
func (sc SpanContext) TraceID() string {
    return sc.span.SpanContext().TraceID().String()
}

// SpanID returns the span ID
func (sc SpanContext) SpanID() string {
    return sc.span.SpanContext().SpanID().String()
}

// StdoutExporter is a simple exporter for debugging
type StdoutExporter struct{}

// NewStdoutExporter creates a new stdout exporter
func NewStdoutExporter() (*StdoutExporter, error) {
    return &StdoutExporter{}, nil
}

// ExportSpans exports spans to stdout
func (e *StdoutExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
    for _, span := range spans {
        fmt.Printf("[TRACE] %s - %s (%s)\n", span.Name(), span.SpanContext().TraceID(), span.EndTime().Sub(span.StartTime()))
    }
    return nil
}

// Shutdown implements the exporter interface
func (e *StdoutExporter) Shutdown(ctx context.Context) error {
    return nil
}

// TracedFunction executes a function within a span
func TracedFunction(ctx context.Context, tracer trace.Tracer, operation string, fn func(context.Context) error) error {
    spanCtx, end := StartSpan(ctx, tracer, operation)
    defer end()

    spanCtx.SetAttributes(
        attribute.String("operation", operation),
    )

    if err := fn(spanCtx.Context); err != nil {
        spanCtx.RecordError(err)
        return err
    }

    spanCtx.SetStatus(codes.Ok, "success")
    return nil
}

// ExtractContext extracts trace context from carrier
func ExtractContext(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
    return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// InjectContext injects trace context into carrier
func InjectContext(ctx context.Context, carrier propagation.TextMapCarrier) {
    otel.GetTextMapPropagator().Inject(ctx, carrier)
}
```

### 2.4 Health Check Pattern

```go
package observability

import (
    "context"
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

// HealthStatus represents the health status of a component
type HealthStatus string

const (
    HealthStatusHealthy   HealthStatus = "healthy"
    HealthStatusDegraded  HealthStatus = "degraded"
    HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// HealthChecker defines the interface for health checks
type HealthChecker interface {
    Name() string
    Check(ctx context.Context) HealthCheckResult
}

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
    Name      string                 `json:"name"`
    Status    HealthStatus           `json:"status"`
    Message   string                 `json:"message,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
    Duration  time.Duration          `json:"duration_ms"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// HealthCheckManager manages multiple health checks
type HealthCheckManager struct {
    checkers map[string]HealthChecker
    mu       sync.RWMutex
    timeout  time.Duration
}

// NewHealthCheckManager creates a new health check manager
func NewHealthCheckManager(timeout time.Duration) *HealthCheckManager {
    if timeout == 0 {
        timeout = 5 * time.Second
    }
    return &HealthCheckManager{
        checkers: make(map[string]HealthChecker),
        timeout:  timeout,
    }
}

// Register registers a health checker
func (h *HealthCheckManager) Register(checker HealthChecker) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.checkers[checker.Name()] = checker
}

// Unregister unregisters a health checker
func (h *HealthCheckManager) Unregister(name string) {
    h.mu.Lock()
    defer h.mu.Unlock()
    delete(h.checkers, name)
}

// CheckAll runs all health checks and returns aggregated results
func (h *HealthCheckManager) CheckAll(ctx context.Context) OverallHealth {
    h.mu.RLock()
    checkers := make([]HealthChecker, 0, len(h.checkers))
    for _, c := range h.checkers {
        checkers = append(checkers, c)
    }
    h.mu.RUnlock()

    results := make([]HealthCheckResult, 0, len(checkers))
    overallStatus := HealthStatusHealthy

    checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
    defer cancel()

    for _, checker := range checkers {
        start := time.Now()
        result := checker.Check(checkCtx)
        result.Duration = time.Since(start)
        results = append(results, result)

        // Aggregate status
        if result.Status == HealthStatusUnhealthy {
            overallStatus = HealthStatusUnhealthy
        } else if result.Status == HealthStatusDegraded && overallStatus == HealthStatusHealthy {
            overallStatus = HealthStatusDegraded
        }
    }

    return OverallHealth{
        Status:    overallStatus,
        Timestamp: time.Now().UTC(),
        Checks:    results,
    }
}

// OverallHealth represents the overall health status
type OverallHealth struct {
    Status    HealthStatus         `json:"status"`
    Timestamp time.Time            `json:"timestamp"`
    Checks    []HealthCheckResult  `json:"checks"`
}

// HTTPHandler returns an HTTP handler for health checks
func (h *HealthCheckManager) HTTPHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        health := h.CheckAll(r.Context())

        statusCode := http.StatusOK
        if health.Status == HealthStatusDegraded {
            statusCode = http.StatusServiceUnavailable
        } else if health.Status == HealthStatusUnhealthy {
            statusCode = http.StatusServiceUnavailable
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(health)
    }
}

// DatabaseHealthChecker checks database health
type DatabaseHealthChecker struct {
    name     string
    pingFunc func(ctx context.Context) error
}

// NewDatabaseHealthChecker creates a database health checker
func NewDatabaseHealthChecker(name string, pingFunc func(ctx context.Context) error) *DatabaseHealthChecker {
    return &DatabaseHealthChecker{
        name:     name,
        pingFunc: pingFunc,
    }
}

// Name returns the checker name
func (d *DatabaseHealthChecker) Name() string {
    return d.name
}

// Check performs the health check
func (d *DatabaseHealthChecker) Check(ctx context.Context) HealthCheckResult {
    result := HealthCheckResult{
        Name:      d.name,
        Timestamp: time.Now().UTC(),
        Metadata:  make(map[string]interface{}),
    }

    if err := d.pingFunc(ctx); err != nil {
        result.Status = HealthStatusUnhealthy
        result.Message = err.Error()
    } else {
        result.Status = HealthStatusHealthy
        result.Message = "connected"
    }

    return result
}

// CacheHealthChecker checks cache health
type CacheHealthChecker struct {
    name     string
    pingFunc func(ctx context.Context) error
}

// NewCacheHealthChecker creates a cache health checker
func NewCacheHealthChecker(name string, pingFunc func(ctx context.Context) error) *CacheHealthChecker {
    return &CacheHealthChecker{
        name:     name,
        pingFunc: pingFunc,
    }
}

// Name returns the checker name
func (c *CacheHealthChecker) Name() string {
    return c.name
}

// Check performs the health check
func (c *CacheHealthChecker) Check(ctx context.Context) HealthCheckResult {
    result := HealthCheckResult{
        Name:      c.name,
        Timestamp: time.Now().UTC(),
        Metadata:  make(map[string]interface{}),
    }

    if err := c.pingFunc(ctx); err != nil {
        result.Status = HealthStatusUnhealthy
        result.Message = err.Error()
    } else {
        result.Status = HealthStatusHealthy
        result.Message = "connected"
    }

    return result
}

// CompositeHealthChecker runs multiple checks and aggregates results
type CompositeHealthChecker struct {
    name     string
    checkers []HealthChecker
}

// NewCompositeHealthChecker creates a composite health checker
func NewCompositeHealthChecker(name string, checkers ...HealthChecker) *CompositeHealthChecker {
    return &CompositeHealthChecker{
        name:     name,
        checkers: checkers,
    }
}

// Name returns the checker name
func (c *CompositeHealthChecker) Name() string {
    return c.name
}

// Check performs all checks
func (c *CompositeHealthChecker) Check(ctx context.Context) HealthCheckResult {
    result := HealthCheckResult{
        Name:      c.name,
        Status:    HealthStatusHealthy,
        Timestamp: time.Now().UTC(),
        Metadata:  make(map[string]interface{}),
    }

    for _, checker := range c.checkers {
        subResult := checker.Check(ctx)
        if subResult.Status == HealthStatusUnhealthy {
            result.Status = HealthStatusUnhealthy
            result.Message = subResult.Message
            break
        } else if subResult.Status == HealthStatusDegraded && result.Status == HealthStatusHealthy {
            result.Status = HealthStatusDegraded
        }
    }

    if result.Status == HealthStatusHealthy {
        result.Message = "all checks passed"
    }

    return result
}
```

---

## 3. Production-Ready Configurations

### 3.1 Kubernetes Deployment with Observability

```yaml
# observability-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: observability-service
  labels:
    app: observability-service
    version: v1.0.0
spec:
  replicas: 3
  selector:
    matchLabels:
      app: observability-service
  template:
    metadata:
      labels:
        app: observability-service
        version: v1.0.0
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
        prometheus.io/path: "/metrics"
    spec:
      containers:
      - name: app
        image: observability-service:v1.0.0
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: SERVICE_NAME
          value: "observability-service"
        - name: SERVICE_VERSION
          value: "1.0.0"
        - name: ENVIRONMENT
          value: "production"
        - name: LOG_LEVEL
          value: "info"
        - name: TRACE_ENDPOINT
          value: "jaeger-collector:14268"
        - name: TRACE_SAMPLING_RATE
          value: "0.1"
        - name: METRICS_PORT
          value: "9090"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        startupProbe:
          httpGet:
            path: /health/startup
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 30
        volumeMounts:
        - name: tmp
          mountPath: /tmp
      volumes:
      - name: tmp
        emptyDir: {}
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - observability-service
              topologyKey: kubernetes.io/hostname

---
apiVersion: v1
kind: Service
metadata:
  name: observability-service
  labels:
    app: observability-service
spec:
  selector:
    app: observability-service
  ports:
  - port: 8080
    targetPort: 8080
    name: http
  - port: 9090
    targetPort: 9090
    name: metrics
  type: ClusterIP

---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: observability-service-monitor
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app: observability-service
  endpoints:
  - port: metrics
    interval: 15s
    path: /metrics
    scrapeTimeout: 10s
```

### 3.2 OpenTelemetry Collector Configuration

```yaml
# otel-collector-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: otel-collector-config
data:
  otel-collector-config.yaml: |
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318

      prometheus:
        config:
          scrape_configs:
            - job_name: 'kubernetes-pods'
              kubernetes_sd_configs:
                - role: pod
              relabel_configs:
                - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
                  action: keep
                  regex: true
                - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
                  action: replace
                  target_label: __metrics_path__
                  regex: (.+)
                - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
                  action: replace
                  regex: ([^:]+)(?::\d+)?;(\d+)
                  replacement: $1:$2
                  target_label: __address__

    processors:
      batch:
        timeout: 1s
        send_batch_size: 1024

      resource:
        attributes:
          - key: environment
            value: production
            action: upsert
          - key: cluster
            value: production-cluster
            action: upsert

      memory_limiter:
        limit_mib: 512
        spike_limit_mib: 128
        check_interval: 5s

      tail_sampling:
        decision_wait: 10s
        num_traces: 100
        expected_new_traces_per_sec: 1000
        policies:
          - name: errors
            type: status_code
            status_code: {status_codes: [ERROR]}
          - name: latency
            type: latency
            latency: {threshold_ms: 1000}

    exporters:
      jaeger:
        endpoint: jaeger-collector:14250
        tls:
          insecure: true

      prometheusremotewrite:
        endpoint: http://prometheus:9090/api/v1/write

      logging:
        verbosity: detailed

      otlp/elastic:
        endpoint: elasticsearch:9200
        tls:
          insecure: true

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [memory_limiter, tail_sampling, batch, resource]
          exporters: [jaeger, logging]

        metrics:
          receivers: [otlp, prometheus]
          processors: [memory_limiter, batch, resource]
          exporters: [prometheusremotewrite, logging]

        logs:
          receivers: [otlp]
          processors: [memory_limiter, batch, resource]
          exporters: [logging]

      extensions: [health_check, pprof, zpages]

  extensions:
    health_check:
      endpoint: 0.0.0.0:13133
    pprof:
      endpoint: 0.0.0.0:1777
    zpages:
      endpoint: 0.0.0.0:55679
```

---

## 4. Security Considerations

### 4.1 Observability Security Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Observability Security Considerations                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Data Type          │  Risk Level  │  Mitigation Strategy                   │
├─────────────────────┼──────────────┼────────────────────────────────────────┤
│  Logs               │              │                                        │
│  ├─ PII Data        │   CRITICAL   │  • Data masking/redaction             │
│  ├─ Credentials     │   CRITICAL   │  • Structured logging patterns        │
│  ├─ Error Details   │   HIGH       │  • Log sanitization                   │
│  └─ IP Addresses    │   MEDIUM     │  • Retention policies                 │
│                                                                             │
│  Metrics            │              │                                        │
│  ├─ Business Data   │   MEDIUM     │  • Label cardinality limits           │
│  ├─ System Info     │   LOW        │  • Authentication on scrape           │
│  └─ Network Info    │   MEDIUM     │  • TLS encryption                     │
│                                                                             │
│  Traces             │              │                                        │
│  ├─ Request Body    │   HIGH       │  • Sampling without body capture      │
│  ├─ Headers         │   HIGH       │  • Header filtering                   │
│  └─ Database Queries│   CRITICAL   │  • Query parameterization checks      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Secure Log Redaction Implementation

```go
package observability

import (
    "regexp"
    "strings"
)

// SensitiveDataRedactor redacts sensitive information from logs
type SensitiveDataRedactor struct {
    patterns map[string]*regexp.Regexp
    mask     string
}

// NewSensitiveDataRedactor creates a new redactor
func NewSensitiveDataRedactor() *SensitiveDataRedactor {
    return &SensitiveDataRedactor{
        mask: "[REDACTED]",
        patterns: map[string]*regexp.Regexp{
            "credit_card":    regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),
            "ssn":            regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
            "email":          regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
            "phone":          regexp.MustCompile(`\b\d{3}-\d{3}-\d{4}\b`),
            "api_key":        regexp.MustCompile(`(?i)(api[_-]?key\s*[:=]\s*)["']?[a-zA-Z0-9]{32,}["']?`),
            "password":       regexp.MustCompile(`(?i)(password\s*[:=]\s*)["']?[^\s"']+["']?`),
            "token":          regexp.MustCompile(`(?i)(token\s*[:=]\s*)["']?[a-zA-Z0-9_-]{20,}["']?`),
            "bearer_token":   regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9_-]{20,}`),
            "auth_header":    regexp.MustCompile(`(?i)authorization\s*[:=]\s*["']?[^"']+["']?`),
        },
    }
}

// Redact redacts sensitive data from a string
func (r *SensitiveDataRedactor) Redact(input string) string {
    result := input
    for name, pattern := range r.patterns {
        result = pattern.ReplaceAllString(result, r.mask+"<"+name+">")
    }
    return result
}

// RedactMap redacts sensitive data from a map
func (r *SensitiveDataRedactor) RedactMap(data map[string]interface{}) map[string]interface{} {
    redacted := make(map[string]interface{})
    sensitiveKeys := []string{
        "password", "secret", "token", "api_key", "apikey",
        "credential", "auth", "authorization", "private_key",
        "credit_card", "ssn", "social_security",
    }

    for key, value := range data {
        lowerKey := strings.ToLower(key)
        isSensitive := false
        for _, sk := range sensitiveKeys {
            if strings.Contains(lowerKey, sk) {
                isSensitive = true
                break
            }
        }

        if isSensitive {
            redacted[key] = r.mask
        } else if str, ok := value.(string); ok {
            redacted[key] = r.Redact(str)
        } else {
            redacted[key] = value
        }
    }

    return redacted
}
```

---

## 5. Compliance Requirements

### 5.1 Compliance Mapping

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Observability Compliance Requirements                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Regulation    │  Requirement          │  Implementation                      │
├────────────────┼───────────────────────┼──────────────────────────────────────┤
│                                                                             │
│  GDPR          │  Data Minimization    │  • Only collect necessary data        │
│  (EU)          │  (Art. 5)             │  • Define retention periods           │
│                │                       │  • Implement right to erasure         │
│                ├───────────────────────┼──────────────────────────────────────┤
│                │  Consent Management   │  • Document lawful basis              │
│                │  (Art. 6)             │  • Audit logging of access            │
│                ├───────────────────────┼──────────────────────────────────────┤
│                │  Security             │  • Encryption at rest/transit         │
│                │  (Art. 32)            │  • Access controls                    │
│                                                                             │
│  SOC 2         │  CC6.1 - Logical      │  • Role-based access to logs          │
│  Type II       │  Access Security      │  • Audit trails for data access       │
│                ├───────────────────────┼──────────────────────────────────────┤
│                │  CC7.2 - System       │  • Comprehensive monitoring           │
│                │  Monitoring           │  • Alerting on anomalies              │
│                ├───────────────────────┼──────────────────────────────────────┤
│                │  CC7.3 - Incident     │  • Incident response procedures       │
│                │  Detection            │  • Automated alerting                 │
│                                                                             │
│  HIPAA         │  §164.312(b) - Audit  │  • Access logs for PHI                │
│  (Healthcare)  │  Controls             │  • Immutable audit trails             │
│                ├───────────────────────┼──────────────────────────────────────┤
│                │  §164.312(c)(1) -     │  • Integrity monitoring               │
│                │  Integrity            │  • Tamper-evident logs                │
│                                                                             │
│  PCI DSS       │  Req. 10 - Logging    │  • Log all access to cardholder data  │
│  (Payment)     │  and Monitoring       │  • Centralized log management         │
│                ├───────────────────────┼──────────────────────────────────────┤
│                │  Req. 11.4 - IDS/IPS  │  • Intrusion detection                │
│                │                       │  • Real-time alerting                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Audit Trail Implementation

```go
package observability

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "time"
)

// AuditEvent represents a compliance audit event
type AuditEvent struct {
    EventID       string                 `json:"event_id"`
    Timestamp     time.Time              `json:"timestamp"`
    EventType     string                 `json:"event_type"`
    Actor         AuditActor             `json:"actor"`
    Resource      AuditResource          `json:"resource"`
    Action        string                 `json:"action"`
    Result        string                 `json:"result"`
    Details       map[string]interface{} `json:"details,omitempty"`
    PreviousHash  string                 `json:"previous_hash,omitempty"`
    EventHash     string                 `json:"event_hash"`
}

// AuditActor represents the entity performing the action
type AuditActor struct {
    ID       string   `json:"id"`
    Type     string   `json:"type"` // user, service, system
    Email    string   `json:"email,omitempty"`
    IP       string   `json:"ip,omitempty"`
    UserAgent string  `json:"user_agent,omitempty"`
}

// AuditResource represents the resource being accessed
type AuditResource struct {
    Type   string `json:"type"`
    ID     string `json:"id"`
    Name   string `json:"name,omitempty"`
}

// AuditLogger provides tamper-evident audit logging
type AuditLogger struct {
    store         AuditStore
    lastHash      string
    redactor      *SensitiveDataRedactor
}

// AuditStore defines the interface for audit storage
type AuditStore interface {
    Append(ctx context.Context, event *AuditEvent) error
    GetLastHash(ctx context.Context) (string, error)
    Query(ctx context.Context, filter AuditFilter) ([]*AuditEvent, error)
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(store AuditStore) (*AuditLogger, error) {
    lastHash, err := store.GetLastHash(context.Background())
    if err != nil {
        lastHash = ""
    }

    return &AuditLogger{
        store:    store,
        lastHash: lastHash,
        redactor: NewSensitiveDataRedactor(),
    }, nil
}

// Log records an audit event
func (a *AuditLogger) Log(ctx context.Context, eventType string, actor AuditActor, resource AuditResource, action, result string, details map[string]interface{}) (*AuditEvent, error) {
    // Redact sensitive data
    details = a.redactor.RedactMap(details)

    event := &AuditEvent{
        EventID:      generateEventID(),
        Timestamp:    time.Now().UTC(),
        EventType:    eventType,
        Actor:        actor,
        Resource:     resource,
        Action:       action,
        Result:       result,
        Details:      details,
        PreviousHash: a.lastHash,
    }

    // Calculate event hash
    event.EventHash = a.calculateHash(event)
    a.lastHash = event.EventHash

    if err := a.store.Append(ctx, event); err != nil {
        return nil, fmt.Errorf("failed to append audit event: %w", err)
    }

    return event, nil
}

// calculateHash creates a tamper-evident hash of the event
func (a *AuditLogger) calculateHash(event *AuditEvent) string {
    data, _ := json.Marshal(struct {
        EventID      string
        Timestamp    time.Time
        EventType    string
        Actor        AuditActor
        Resource     AuditResource
        Action       string
        Result       string
        Details      map[string]interface{}
        PreviousHash string
    }{
        EventID:      event.EventID,
        Timestamp:    event.Timestamp,
        EventType:    event.EventType,
        Actor:        event.Actor,
        Resource:     event.Resource,
        Action:       event.Action,
        Result:       event.Result,
        Details:      event.Details,
        PreviousHash: event.PreviousHash,
    })

    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}

// VerifyIntegrity verifies the integrity of the audit chain
func (a *AuditLogger) VerifyIntegrity(ctx context.Context, filter AuditFilter) (bool, error) {
    events, err := a.store.Query(ctx, filter)
    if err != nil {
        return false, err
    }

    for i, event := range events {
        // Verify event hash
        expectedHash := a.calculateHash(event)
        if expectedHash != event.EventHash {
            return false, fmt.Errorf("hash mismatch at event %d (ID: %s)", i, event.EventID)
        }

        // Verify chain integrity (except first event)
        if i > 0 {
            expectedPrevHash := events[i-1].EventHash
            if event.PreviousHash != expectedPrevHash {
                return false, fmt.Errorf("chain broken at event %d (ID: %s)", i, event.EventID)
            }
        }
    }

    return true, nil
}

// AuditFilter defines filtering criteria for audit queries
type AuditFilter struct {
    StartTime    *time.Time
    EndTime      *time.Time
    EventTypes   []string
    ActorID      string
    ResourceType string
    ResourceID   string
    Actions      []string
    Limit        int
}

// generateEventID generates a unique event ID
func generateEventID() string {
    return fmt.Sprintf("evt_%d_%s", time.Now().UnixNano(), generateRandomString(8))
}

// generateRandomString generates a random string (simplified)
func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}
```

---

## 6. Decision Matrices

### 6.1 Observability Tool Selection Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Observability Tool Selection Matrix                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Criteria               │ Weight │ Prometheus │ Datadog │ New Relic │ Grafana │
│  ───────────────────────┼────────┼────────────┼─────────┼───────────┼─────────│
│  Cost                   │  20%   │    ★★★★★   │  ★★☆☆☆  │  ★★☆☆☆   │  ★★★★☆  │
│  Ease of Setup          │  15%   │    ★★★☆☆   │  ★★★★★  │  ★★★★☆   │  ★★★★☆  │
│  Scalability            │  20%   │    ★★★★★   │  ★★★★☆  │  ★★★★☆   │  ★★★★★  │
│  Customization          │  15%   │    ★★★★★   │  ★★★☆☆  │  ★★★☆☆   │  ★★★★★  │
│  Alerting               │  15%   │    ★★★★☆   │  ★★★★★  │  ★★★★★   │  ★★★★☆  │
│  Support                │  10%   │    ★★★☆☆   │  ★★★★★  │  ★★★★★   │  ★★★☆☆  │
│  Integration            │   5%   │    ★★★★☆   │  ★★★★★  │  ★★★★★   │  ★★★★★  │
│  ───────────────────────┼────────┼────────────┼─────────┼───────────┼─────────│
│  WEIGHTED SCORE         │  100%  │    4.35    │  3.65   │  3.45     │  4.15   │
│                                                                             │
│  Recommendation:                                                          │
│  • Open Source / Budget Conscious: Prometheus + Grafana                    │
│  • Enterprise / Quick Setup: Datadog or New Relic                          │
│  • Maximum Flexibility: Prometheus + Grafana                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Sampling Strategy Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Sampling Strategy Decision Matrix                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Traffic Volume  │  Criticality  │  Budget    │  Recommended Strategy       │
├──────────────────┼───────────────┼────────────┼─────────────────────────────│
│  Low (<100 RPM)  │  High         │  Low       │  100% Sampling              │
│  Low (<100 RPM)  │  Medium       │  Low       │  100% Sampling              │
│  Low (<100 RPM)  │  Low          │  Low       │  50% Sampling               │
│  ────────────────┼───────────────┼────────────┼─────────────────────────────│
│  Medium (100-10K │  High         │  Medium    │  Tail-based + Error 100%    │
│  RPM)            │  Medium       │  Medium    │  10% Head-based             │
│                  │  Low          │  Medium    │  1% Head-based              │
│  ────────────────┼───────────────┼────────────┼─────────────────────────────│
│  High (>10K RPM) │  High         │  High      │  Tail-based + Errors        │
│                  │  Medium       │  High      │  0.1% Head-based            │
│                  │  Low          │  High      │  Adaptive (0.01-1%)         │
│  ────────────────┼───────────────┼────────────┼─────────────────────────────│
│  High (>10K RPM) │  Any          │  Low       │  Adaptive with limits       │
│                                                                             │
│  Strategy Definitions:                                                    │
│  • Head-based: Sample at the start of the trace                            │
│  • Tail-based: Sample after seeing all spans (keeps errors/slow)           │
│  • Adaptive: Dynamically adjust based on load/error rate                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Alert Severity Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Alert Severity Classification                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Severity  │  Response Time  │  Notification    │  Escalation                │
├────────────┼─────────────────┼──────────────────┼────────────────────────────│
│  P0 -      │  Immediate      │  Page + SMS +    │  Auto-escalate to          │
│  Critical  │  (<5 min)       │  Call + Slack    │  manager after 10 min      │
│            │                 │                  │  • All hands on deck       │
│            │                 │                  │  • War room if needed      │
│  ──────────┼─────────────────┼──────────────────┼────────────────────────────│
│  P1 -      │  <15 min        │  Page + Slack    │  Escalate to on-call       │
│  High      │                 │                  │  lead after 30 min         │
│            │                 │                  │  • Customer impact         │
│            │                 │                  │  • Degraded performance    │
│  ──────────┼─────────────────┼──────────────────┼────────────────────────────│
│  P2 -      │  <1 hour        │  Slack + Email   │  Next business day         │
│  Medium    │                 │                  │  • No customer impact      │
│            │                 │                  │  • Capacity concerns       │
│  ──────────┼─────────────────┼──────────────────┼────────────────────────────│
│  P3 -      │  <24 hours      │  Email only      │  Weekly review             │
│  Low       │                 │                  │  • Tracking issues         │
│            │                 │                  │  • Optimization items      │
│  ──────────┼─────────────────┼──────────────────┼────────────────────────────│
│  P4 -      │  Next sprint    │  Ticket created  │  Backlog prioritization    │
│  Info      │                 │                  │  • Trends analysis         │
│            │                 │                  │  • Planning input          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Observability-Driven Development Best Practices             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. DESIGN PHASE                                                            │
│     ✓ Define SLOs (Service Level Objectives) before implementation          │
│     ✓ Identify critical user journeys to trace                              │
│     ✓ Plan metric cardinality carefully                                     │
│     ✓ Design log structure for searchability                                │
│                                                                             │
│  2. IMPLEMENTATION PHASE                                                    │
│     ✓ Instrument as you code, not after                                     │
│     ✓ Use structured logging exclusively                                    │
│     ✓ Include trace context in all async operations                         │
│     ✓ Add health checks for all dependencies                                │
│     ✓ Implement graceful shutdown with telemetry flush                      │
│                                                                             │
│  3. DEPLOYMENT PHASE                                                        │
│     ✓ Canary deployments with metric validation                             │
│     ✓ Automated rollback on SLO violation                                   │
│     ✓ Feature flags with observability hooks                                │
│     ✓ Gradual traffic shifting with monitoring                              │
│                                                                             │
│  4. OPERATIONS PHASE                                                        │
│     ✓ Dashboards for every critical service                                 │
│     ✓ Alert on symptoms, not causes                                         │
│     ✓ Regular review of alert fatigue                                       │
│     ✓ Continuous refinement of SLOs                                         │
│                                                                             │
│  5. CONTINUOUS IMPROVEMENT                                                  │
│     ✓ Post-incident review of observability gaps                            │
│     ✓ Quarterly SLO review and adjustment                                   │
│     ✓ Regular sampling rate optimization                                    │
│     ✓ Cost vs. granularity trade-off analysis                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Google SRE Book - Monitoring Distributed Systems
2. OpenTelemetry Specification
3. Prometheus Best Practices
4. Distributed Systems Observability - Cindy Sridharan
5. The Art of Monitoring - James Turnbull
