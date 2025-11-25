package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
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

	// TODO: 实现 OTLP metrics 导出器
	// 暂时使用空的 reader，后续添加 OTLP exporter
	// 创建指标提供者
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		// metric.WithReader(metric.NewPeriodicReader(exporter)),
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
