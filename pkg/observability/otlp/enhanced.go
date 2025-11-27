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

	"github.com/yourusername/golang/pkg/sampling"
)

// EnhancedOTLP 增强的 OTLP 集成
// 提供采样、追踪、指标的完整支持
type EnhancedOTLP struct {
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
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
	traceExporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(cfg.Insecure),
	)
	if err != nil {
		return nil, err
	}

	// 创建追踪提供者（带采样）
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(newSamplerWrapper(sampler)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 创建指标导出器
	metricExporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
		otlpmetricgrpc.WithInsecure(cfg.Insecure),
	)
	if err != nil {
		return nil, err
	}

	// 创建指标提供者
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
			sdkmetric.WithInterval(10*time.Second),
		)),
	)

	otel.SetMeterProvider(mp)

	return &EnhancedOTLP{
		tracerProvider: tp,
		meterProvider:  mp,
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
			Decision:   sdktrace.RecordAndSample,
			Tracestate: params.Tracestate,
		}
	}
	return sdktrace.SamplingResult{
		Decision:   sdktrace.Drop,
		Tracestate: params.Tracestate,
	}
}

func (s *samplerWrapper) Description() string {
	return "Framework Sampler Wrapper"
}
