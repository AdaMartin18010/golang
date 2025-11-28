package system

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// APMMonitor 应用性能监控器
// 提供应用级别的性能监控
type APMMonitor struct {
	meter           metric.Meter
	tracer          trace.Tracer
	enabled         bool
	collectInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc

	// 指标
	requestDurationHistogram metric.Float64Histogram
	requestCountCounter      metric.Int64Counter
	errorCountCounter        metric.Int64Counter
	activeConnectionsGauge   metric.Int64ObservableGauge
	throughputGauge          metric.Float64ObservableGauge
}

// APMConfig APM 配置
type APMConfig struct {
	Meter           metric.Meter
	Tracer          trace.Tracer
	Enabled         bool
	CollectInterval time.Duration
}

// NewAPMMonitor 创建 APM 监控器
func NewAPMMonitor(cfg APMConfig) (*APMMonitor, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	collectInterval := cfg.CollectInterval
	if collectInterval == 0 {
		collectInterval = 10 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitor := &APMMonitor{
		meter:           cfg.Meter,
		tracer:          cfg.Tracer,
		enabled:         cfg.Enabled,
		collectInterval: collectInterval,
		ctx:             ctx,
		cancel:          cancel,
	}

	// 初始化指标
	if err := monitor.initMetrics(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to init metrics: %w", err)
	}

	return monitor, nil
}

// initMetrics 初始化指标
func (m *APMMonitor) initMetrics() error {
	var err error

	m.requestDurationHistogram, err = m.meter.Float64Histogram(
		"apm.request.duration",
		metric.WithDescription("Request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	m.requestCountCounter, err = m.meter.Int64Counter(
		"apm.request.count",
		metric.WithDescription("Total number of requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.errorCountCounter, err = m.meter.Int64Counter(
		"apm.error.count",
		metric.WithDescription("Total number of errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.activeConnectionsGauge, err = m.meter.Int64ObservableGauge(
		"apm.connections.active",
		metric.WithDescription("Number of active connections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.throughputGauge, err = m.meter.Float64ObservableGauge(
		"apm.throughput",
		metric.WithDescription("Requests per second"),
		metric.WithUnit("1/s"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Start 启动 APM 监控
func (m *APMMonitor) Start() error {
	if !m.enabled {
		return nil
	}

	// 注册可观察指标回调
	_, err := m.meter.RegisterCallback(m.collectMetrics, m.activeConnectionsGauge, m.throughputGauge)
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	return nil
}

// Stop 停止 APM 监控
func (m *APMMonitor) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

// collectMetrics 收集指标（可观察指标回调）
func (m *APMMonitor) collectMetrics(ctx context.Context, obs metric.Observer) error {
	// 可以在这里收集应用级别的指标
	// 当前为占位实现
	return nil
}

// RecordRequest 记录请求
func (m *APMMonitor) RecordRequest(ctx context.Context, duration time.Duration, statusCode int, attrs ...attribute.KeyValue) {
	if !m.enabled {
		return
	}

	m.requestCountCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.requestDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	if statusCode >= 400 {
		m.errorCountCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordError 记录错误
func (m *APMMonitor) RecordError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	if !m.enabled {
		return
	}

	allAttrs := append(attrs, attribute.String("error.type", fmt.Sprintf("%T", err)))
	m.errorCountCounter.Add(ctx, 1, metric.WithAttributes(allAttrs...))
}

// StartSpan 开始追踪 Span
func (m *APMMonitor) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if m.tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}
	return m.tracer.Start(ctx, name, opts...)
}
