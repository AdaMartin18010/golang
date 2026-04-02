# OpenTelemetry 分布式追踪生产实践 (OpenTelemetry Distributed Tracing Production Guide)

> **分类**: 工程与云原生
> **标签**: #opentelemetry #distributed-tracing #observability #production
> **参考**: OpenTelemetry Go SDK v1.24+, W3C Trace Context

---

## 生产级 SDK 配置

```go
package telemetry

import (
    "context"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

// TracerConfig 追踪器配置
type TracerConfig struct {
    ServiceName       string
    ServiceVersion    string
    Environment       string
    OTLPEndpoint      string
    OTLPInsecure      bool
    OTLPHeaders       map[string]string
    SampleRate        float64
    MaxQueueSize      int
    BatchTimeout      time.Duration
    ExportTimeout     time.Duration
    MaxExportBatchSize int
}

// InitTracerProvider 初始化生产级 TracerProvider
func InitTracerProvider(ctx context.Context, cfg TracerConfig) (*sdktrace.TracerProvider, error) {
    // 1. 创建资源
    res, err := resource.New(ctx,
        resource.WithFromEnv(),
        resource.WithProcess(),
        resource.WithTelemetrySDK(),
        resource.WithHost(),
        resource.WithAttributes(
            semconv.ServiceName(cfg.ServiceName),
            semconv.ServiceVersion(cfg.ServiceVersion),
            attribute.String("environment", cfg.Environment),
            attribute.String("host.id", getHostID()),
            attribute.String("host.name", getHostname()),
        ),
    )
    if err != nil {
        return nil, err
    }

    // 2. 创建 OTLP Exporter
    exporter, err := createOTLPExporter(ctx, cfg)
    if err != nil {
        return nil, err
    }

    // 3. 配置采样策略
    sampler := createSampler(cfg.SampleRate)

    // 4. 配置 Span Processor
    batchProcessor := sdktrace.NewBatchSpanProcessor(exporter,
        sdktrace.WithMaxQueueSize(cfg.MaxQueueSize),
        sdktrace.WithBatchTimeout(cfg.BatchTimeout),
        sdktrace.WithExportTimeout(cfg.ExportTimeout),
        sdktrace.WithMaxExportBatchSize(cfg.MaxExportBatchSize),
    )

    // 5. 创建 TracerProvider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sampler),
        sdktrace.WithSpanProcessor(batchProcessor),
    )

    // 6. 设置为全局 Provider
    otel.SetTracerProvider(tp)

    // 7. 配置传播器 (W3C Trace Context + Baggage)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return tp, nil
}

// createOTLPExporter 创建 OTLP gRPC Exporter
func createOTLPExporter(ctx context.Context, cfg TracerConfig) (sdktrace.SpanExporter, error) {
    opts := []otlptracegrpc.Option{
        otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
        otlptracegrpc.WithTimeout(cfg.ExportTimeout),
    }

    if cfg.OTLPInsecure {
        opts = append(opts, otlptracegrpc.WithInsecure())
    } else {
        creds := credentials.NewClientTLSFromCert(nil, "")
        opts = append(opts, otlptracegrpc.WithTLSCredentials(creds))
    }

    if len(cfg.OTLPHeaders) > 0 {
        opts = append(opts, otlptracegrpc.WithHeaders(cfg.OTLPHeaders))
    }

    // 自定义 gRPC 连接选项
    opts = append(opts, otlptracegrpc.WithDialOption(
        grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
    ))

    return otlptracegrpc.New(ctx, opts...)
}

// createSampler 创建采样器
func createSampler(rate float64) sdktrace.Sampler {
    // 父采样优先 + 按比例采样
    return sdktrace.ParentBased(
        sdktrace.TraceIDRatioBased(rate),
        sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
        sdktrace.WithRemoteParentNotSampled(sdktrace.NeverSample()),
        sdktrace.WithLocalParentSampled(sdktrace.AlwaysSample()),
        sdktrace.WithLocalParentNotSampled(sdktrace.TraceIDRatioBased(rate)),
    )
}
```

---

## HTTP/gRPC 自动埋点

```go
// HTTP 服务埋点
import (
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func setupHTTPServer() *http.Server {
    // 包装 Handler 自动创建 Span
    handler := otelhttp.NewHandler(
        http.DefaultServeMux,
        "http-server",
        otelhttp.WithPublicEndpoint(),
        otelhttp.WithServerName("api-gateway"),
        otelhttp.WithSpanOptions(
            trace.WithAttributes(
                attribute.String("http.server_name", "api-gateway"),
            ),
        ),
    )

    return &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }
}

// HTTP 客户端埋点
func createHTTPClient() *http.Client {
    return &http.Client{
        Transport: otelhttp.NewTransport(
            http.DefaultTransport,
            otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
                return otelhttp.DefaultClientTrace(ctx)
            }),
        ),
        Timeout: 30 * time.Second,
    }
}

// gRPC 服务端埋点
import (
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    "google.golang.org/grpc"
)

func setupGRPCServer() *grpc.Server {
    server := grpc.NewServer(
        grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(
            otelgrpc.WithPropagators(propagation.TraceContext{}),
        )),
        grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
    )

    return server
}

// gRPC 客户端埋点
func createGRPCClient(target string) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(target,
        grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
        grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
    )
    return conn, err
}
```

---

## 消息队列上下文传播

```go
// Kafka 消息传播
package kafka

import (
    "context"
    "github.com/segmentio/kafka-go"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/trace"
)

type KafkaCarrier struct {
    Headers []kafka.Header
}

func (c *KafkaCarrier) Get(key string) string {
    for _, h := range c.Headers {
        if h.Key == key {
            return string(h.Value)
        }
    }
    return ""
}

func (c *KafkaCarrier) Set(key, value string) {
    // 删除已存在的同名 header
    for i, h := range c.Headers {
        if h.Key == key {
            c.Headers = append(c.Headers[:i], c.Headers[i+1:]...)
            break
        }
    }
    c.Headers = append(c.Headers, kafka.Header{
        Key:   key,
        Value: []byte(value),
    })
}

func (c *KafkaCarrier) Keys() []string {
    keys := make([]string, len(c.Headers))
    for i, h := range c.Headers {
        keys[i] = h.Key
    }
    return keys
}

// InjectTracingHeaders 注入追踪上下文到 Kafka 消息
func InjectTracingHeaders(ctx context.Context, msg *kafka.Message) {
    carrier := &KafkaCarrier{Headers: msg.Headers}
    propagator := otel.GetTextMapPropagator()
    propagator.Inject(ctx, carrier)
    msg.Headers = carrier.Headers
}

// ExtractTracingContext 从 Kafka 消息提取追踪上下文
func ExtractTracingContext(ctx context.Context, msg *kafka.Message) context.Context {
    carrier := &KafkaCarrier{Headers: msg.Headers}
    propagator := otel.GetTextMapPropagator()
    return propagator.Extract(ctx, carrier)
}

// 生产者包装
func (p *Producer) Produce(ctx context.Context, topic string, key, value []byte) error {
    tracer := otel.Tracer("kafka-producer")
    ctx, span := tracer.Start(ctx, "kafka.produce",
        trace.WithAttributes(
            attribute.String("messaging.system", "kafka"),
            attribute.String("messaging.destination", topic),
            attribute.String("messaging.destination_kind", "topic"),
        ),
    )
    defer span.End()

    msg := kafka.Message{
        Topic: topic,
        Key:   key,
        Value: value,
    }

    // 注入追踪上下文
    InjectTracingHeaders(ctx, &msg)

    return p.writer.WriteMessages(ctx, msg)
}

// 消费者处理
func (c *Consumer) ProcessMessage(ctx context.Context, msg kafka.Message) error {
    // 提取追踪上下文
    ctx = ExtractTracingContext(ctx, msg)

    tracer := otel.Tracer("kafka-consumer")
    ctx, span := tracer.Start(ctx, "kafka.consume",
        trace.WithAttributes(
            attribute.String("messaging.system", "kafka"),
            attribute.String("messaging.source", msg.Topic),
            attribute.String("messaging.operation", "receive"),
            attribute.String("messaging.kafka.partition", strconv.Itoa(msg.Partition)),
            attribute.Int64("messaging.kafka.offset", msg.Offset),
        ),
        trace.WithSpanKind(trace.SpanKindConsumer),
    )
    defer span.End()

    // 处理消息
    return c.handler(ctx, msg)
}
```

---

## 数据库追踪埋点

```go
// SQL 数据库追踪
package sqltrace

import (
    "context"
    "database/sql"
    "database/sql/driver"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

// TracedDB 包装 sql.DB
type TracedDB struct {
    *sql.DB
    tracer trace.Tracer
}

// QueryContext 追踪查询
func (db *TracedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    ctx, span := db.tracer.Start(ctx, "db.query",
        trace.WithAttributes(
            attribute.String("db.system", db.detectDBSystem()),
            attribute.String("db.statement", query),
            attribute.String("db.operation", "SELECT"),
        ),
    )
    defer span.End()

    rows, err := db.DB.QueryContext(ctx, query, args...)
    if err != nil {
        span.RecordError(err)
    }

    return rows, err
}

// ExecContext 追踪执行
func (db *TracedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    ctx, span := db.tracer.Start(ctx, "db.exec",
        trace.WithAttributes(
            attribute.String("db.system", db.detectDBSystem()),
            attribute.String("db.statement", query),
            attribute.String("db.operation", detectOperation(query)),
        ),
    )
    defer span.End()

    result, err := db.DB.ExecContext(ctx, query, args...)
    if err != nil {
        span.RecordError(err)
    }

    return result, err
}

// Redis 追踪
package redistrace

import (
    "github.com/redis/go-redis/v9"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

type TracedRedisClient struct {
    *redis.Client
    tracer trace.Tracer
}

func (c *TracedRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
    ctx, span := c.tracer.Start(ctx, "redis.get",
        trace.WithAttributes(
            attribute.String("db.system", "redis"),
            attribute.String("db.operation", "GET"),
            attribute.String("db.redis.key", key),
        ),
    )
    defer span.End()

    cmd := c.Client.Get(ctx, key)
    if err := cmd.Err(); err != nil && err != redis.Nil {
        span.RecordError(err)
    }

    return cmd
}

func (c *TracedRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
    ctx, span := c.tracer.Start(ctx, "redis.set",
        trace.WithAttributes(
            attribute.String("db.system", "redis"),
            attribute.String("db.operation", "SET"),
            attribute.String("db.redis.key", key),
        ),
    )
    defer span.End()

    cmd := c.Client.Set(ctx, key, value, expiration)
    if err := cmd.Err(); err != nil {
        span.RecordError(err)
    }

    return cmd
}
```

---

## 高级采样策略

```go
// 尾部采样实现 (Tail-based Sampling)

type TailSampler struct {
    storage     SpanStorage
    maxWaitTime time.Duration
    policies    []SamplingPolicy
}

type SamplingPolicy interface {
    ShouldSample(traceID string, spans []ReadableSpan) bool
}

// ErrorPolicy 错误采样
func (p *ErrorPolicy) ShouldSample(traceID string, spans []ReadableSpan) bool {
    for _, span := range spans {
        if span.Status().Code == codes.Error {
            return true
        }
        for _, event := range span.Events() {
            if event.Name == "exception" {
                return true
            }
        }
    }
    return false
}

// LatencyPolicy 延迟采样
func (p *LatencyPolicy) ShouldSample(traceID string, spans []ReadableSpan) bool {
    var maxDuration time.Duration
    for _, span := range spans {
        duration := span.EndTime().Sub(span.StartTime())
        if duration > maxDuration {
            maxDuration = duration
        }
    }
    return maxDuration > p.Threshold
}

// AttributePolicy 属性匹配采样
func (p *AttributePolicy) ShouldSample(traceID string, spans []ReadableSpan) bool {
    for _, span := range spans {
        for _, attr := range span.Attributes() {
            if attr.Key == p.Key && attr.Value.AsString() == p.Value {
                return true
            }
        }
    }
    return false
}

// ProbabilityPolicy 概率采样
func (p *ProbabilityPolicy) ShouldSample(traceID string, spans []ReadableSpan) bool {
    // 使用 traceID 的哈希值确保一致性
    hash := fnv64a(traceID)
    return float64(hash%10000)/10000 < p.Rate
}

// 自适应采样
func (p *AdaptivePolicy) ShouldSample(traceID string, spans []ReadableSpan) bool {
    // 根据当前吞吐量动态调整采样率
    currentRate := p.getCurrentRate()

    // 高优先级请求始终采样
    for _, span := range spans {
        for _, attr := range span.Attributes() {
            if attr.Key == "http.route" {
                if isCriticalRoute(attr.Value.AsString()) {
                    return true
                }
            }
        }
    }

    hash := fnv64a(traceID)
    return float64(hash%10000)/10000 < currentRate
}
```

---

## 日志与追踪关联

```go
// 结构化日志关联 TraceID
package logging

import (
    "context"
    "go.opentelemetry.io/otel/trace"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// LoggerWithTrace 添加追踪信息到日志
type LoggerWithTrace struct {
    logger *zap.Logger
}

func (l *LoggerWithTrace) Info(ctx context.Context, msg string, fields ...zapcore.Field) {
    fields = append(fields, extractTraceFields(ctx)...)
    l.logger.Info(msg, fields...)
}

func (l *LoggerWithTrace) Error(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
    fields = append(fields, zap.Error(err))
    fields = append(fields, extractTraceFields(ctx)...)
    l.logger.Error(msg, fields...)
}

func extractTraceFields(ctx context.Context) []zapcore.Field {
    spanContext := trace.SpanContextFromContext(ctx)
    if !spanContext.IsValid() {
        return nil
    }

    return []zapcore.Field{
        zap.String("trace_id", spanContext.TraceID().String()),
        zap.String("span_id", spanContext.SpanID().String()),
        zap.Bool("trace_sampled", spanContext.IsSampled()),
    }
}

// Logstash/ELK 格式
func (l *LoggerWithTrace) WithELKFormat(ctx context.Context, msg string, fields map[string]interface{}) {
    spanContext := trace.SpanContextFromContext(ctx)

    logEntry := map[string]interface{}{
        "message":     msg,
        "timestamp":   time.Now().UTC(),
        "level":       "info",
        "service":     l.serviceName,
        "trace_id":    spanContext.TraceID().String(),
        "span_id":     spanContext.SpanID().String(),
        "fields":      fields,
    }

    // 输出 JSON
    json.NewEncoder(os.Stdout).Encode(logEntry)
}
```

---

## 性能优化与监控

```go
// Tracer 性能监控

type TracerMetrics struct {
    spansCreated     prometheus.Counter
    spansExported    prometheus.Counter
    spansDropped     prometheus.Counter
    exportDuration   prometheus.Histogram
    queueSize        prometheus.Gauge
}

func NewTracerMetrics() *TracerMetrics {
    return &TracerMetrics{
        spansCreated: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "otel_spans_created_total",
            Help: "Total number of spans created",
        }),
        exportDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "otel_export_duration_seconds",
            Help:    "Duration of export operations",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
        }),
        queueSize: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "otel_span_queue_size",
            Help: "Current size of the span queue",
        }),
    }
}

// 内存限制
func WithMemoryLimit(limit int) sdktrace.TracerProviderOption {
    return sdktrace.WithSpanProcessor(
        sdktrace.NewBatchSpanProcessor(
            nil,
            sdktrace.WithMaxQueueSize(limit),
        ),
    )
}

// 批处理优化
func optimizedBatchProcessor(exporter sdktrace.SpanExporter) sdktrace.SpanProcessor {
    return sdktrace.NewBatchSpanProcessor(exporter,
        sdktrace.WithBatchTimeout(100*time.Millisecond),
        sdktrace.WithExportTimeout(30*time.Second),
        sdktrace.WithMaxQueueSize(2048),
        sdktrace.WithMaxExportBatchSize(512),
    )
}
```
