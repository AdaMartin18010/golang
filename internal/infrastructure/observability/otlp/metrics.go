package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// MetricsProvider OpenTelemetry 指标提供者
type MetricsProvider struct {
	provider *metric.MeterProvider
}

// NewMetricsProvider 创建指标提供者
func NewMetricsProvider(ctx context.Context, endpoint string, insecure bool) (*MetricsProvider, error) {
	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("golang-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建 OTLP 导出器
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
	}
	if insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// 创建指标提供者
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)

	return &MetricsProvider{
		provider: mp,
	}, nil
}

// MeterProvider 获取指标提供者
func (mp *MetricsProvider) MeterProvider() *metric.MeterProvider {
	return mp.provider
}

// Shutdown 关闭指标提供者
func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
	return mp.provider.Shutdown(ctx)
}
