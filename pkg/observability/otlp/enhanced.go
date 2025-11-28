package otlp

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/yourusername/golang/pkg/sampling"
)

// EnhancedOTLP 增强的 OTLP 集成
// 提供采样、追踪、指标的完整支持
// 注意：日志导出器需要 OpenTelemetry 日志 SDK 正式发布后实现
type EnhancedOTLP struct {
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	logExporter    LogExporter // 日志导出器（占位实现）
	sampler        sampling.Sampler
	resource       *resource.Resource
}

// Config 配置
type Config struct {
	ServiceName    string
	ServiceVersion string
	Endpoint       string
	Insecure       bool
	Sampler        sampling.Sampler
	SampleRate     float64
	// MetricInterval 指标导出间隔（默认：10秒）
	MetricInterval time.Duration
	// TraceBatchTimeout 追踪批处理超时（默认：5秒）
	TraceBatchTimeout time.Duration
	// TraceBatchSize 追踪批处理大小（默认：512）
	TraceBatchSize int
}

// NewEnhancedOTLP 创建增强的 OTLP 集成
func NewEnhancedOTLP(cfg Config) (*EnhancedOTLP, error) {
	// 创建资源
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	// 创建采样器
	sampler := cfg.Sampler
	if sampler == nil {
		if cfg.SampleRate > 0 {
			sampler, _ = sampling.NewProbabilisticSampler(cfg.SampleRate)
		} else {
			sampler = sampling.NewAlwaysSampler()
		}
	}

	// 创建追踪导出器
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	// 生产环境应使用 TLS：
	// opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(&tls.Config{})))
	traceExporter, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// 配置追踪批处理选项
	batchTimeout := cfg.TraceBatchTimeout
	if batchTimeout == 0 {
		batchTimeout = 5 * time.Second
	}
	batchSize := cfg.TraceBatchSize
	if batchSize == 0 {
		batchSize = 512
	}

	// 创建追踪提供者（带采样和批处理配置）
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBatchTimeout(batchTimeout),
			sdktrace.WithMaxExportBatchSize(batchSize),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(newSamplerWrapper(sampler)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 创建指标导出器
	metricOpts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
	}
	if cfg.Insecure {
		metricOpts = append(metricOpts, otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
	}
	metricExporter, err := otlpmetricgrpc.New(context.Background(), metricOpts...)
	if err != nil {
		return nil, err
	}

	// 配置指标导出间隔
	metricInterval := cfg.MetricInterval
	if metricInterval == 0 {
		metricInterval = 10 * time.Second
	}

	// 创建指标提供者（带周期性读取器）
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
			sdkmetric.WithInterval(metricInterval),
		)),
	)
	otel.SetMeterProvider(mp)

	// 创建日志导出器（占位实现）
	// 注意：OpenTelemetry 日志导出器（otlploggrpc）可能尚未正式发布
	// 当前使用占位实现，等待官方发布后替换
	logExporter, err := NewLogExporter(LogExporterConfig{
		Endpoint: cfg.Endpoint,
		Insecure: cfg.Insecure,
	})
	if err != nil {
		// 如果创建失败，使用占位实现（禁用状态）
		// 这是预期的行为，因为当前日志导出器尚未正式实现
		logExporter = NewPlaceholderLogExporter()
	}

	return &EnhancedOTLP{
		tracerProvider: tp,
		meterProvider:  mp,
		logExporter:    logExporter,
		sampler:        sampler,
		resource:       res,
	}, nil
}

// Shutdown 关闭
func (e *EnhancedOTLP) Shutdown(ctx context.Context) error {
	if err := e.tracerProvider.Shutdown(ctx); err != nil {
		return err
	}
	if err := e.meterProvider.Shutdown(ctx); err != nil {
		return err
	}
	if e.logExporter != nil {
		if err := e.logExporter.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Tracer 获取追踪器
func (e *EnhancedOTLP) Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// Meter 获取指标器
func (e *EnhancedOTLP) Meter(name string) metric.Meter {
	return otel.Meter(name)
}

// LogExporter 获取日志导出器
// 当前返回占位实现，等待官方发布后返回实际实现
func (e *EnhancedOTLP) LogExporter() LogExporter {
	return e.logExporter
}

// GetResource 获取资源信息
func (e *EnhancedOTLP) GetResource() *resource.Resource {
	return e.resource
}

// ShouldSample 检查是否应该采样
func (e *EnhancedOTLP) ShouldSample(ctx context.Context) bool {
	return e.sampler.ShouldSample(ctx)
}

// UpdateSampleRate 更新采样率
func (e *EnhancedOTLP) UpdateSampleRate(rate float64) error {
	return e.sampler.UpdateRate(rate)
}

// samplerWrapper 采样器包装器
// 将框架采样器适配到 OpenTelemetry 采样器
type samplerWrapper struct {
	sampler sampling.Sampler
}

func newSamplerWrapper(sampler sampling.Sampler) sdktrace.Sampler {
	return &samplerWrapper{sampler: sampler}
}

func (s *samplerWrapper) ShouldSample(params sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if s.sampler.ShouldSample(params.ParentContext) {
		return sdktrace.SamplingResult{
			Decision: sdktrace.RecordAndSample,
		}
	}
	return sdktrace.SamplingResult{
		Decision: sdktrace.Drop,
	}
}

func (s *samplerWrapper) Description() string {
	return "Framework Sampler Wrapper"
}
