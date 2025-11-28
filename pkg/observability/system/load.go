package system

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// LoadMonitor 负载监控器
// 监控系统负载、请求速率、并发数等
type LoadMonitor struct {
	meter           metric.Meter
	enabled         bool
	collectInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc

	// 指标
	loadAverageGauge    metric.Float64ObservableGauge
	requestRateCounter  metric.Int64Counter
	concurrentRequestsGauge metric.Int64ObservableGauge
	queueLengthGauge   metric.Int64ObservableGauge
}

// NewLoadMonitor 创建负载监控器
func NewLoadMonitor(cfg Config) (*LoadMonitor, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	collectInterval := cfg.CollectInterval
	if collectInterval == 0 {
		collectInterval = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitor := &LoadMonitor{
		meter:           cfg.Meter,
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
func (m *LoadMonitor) initMetrics() error {
	var err error

	m.loadAverageGauge, err = m.meter.Float64ObservableGauge(
		"system.load.average",
		metric.WithDescription("System load average"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.requestRateCounter, err = m.meter.Int64Counter(
		"system.load.request.rate",
		metric.WithDescription("Request rate per second"),
		metric.WithUnit("1/s"),
	)
	if err != nil {
		return err
	}

	m.concurrentRequestsGauge, err = m.meter.Int64ObservableGauge(
		"system.load.concurrent.requests",
		metric.WithDescription("Number of concurrent requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.queueLengthGauge, err = m.meter.Int64ObservableGauge(
		"system.load.queue.length",
		metric.WithDescription("Queue length"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Start 启动负载监控
func (m *LoadMonitor) Start() error {
	if !m.enabled {
		return nil
	}

	// 注册可观察指标回调
	_, err := m.meter.RegisterCallback(m.collectMetrics, m.loadAverageGauge, m.concurrentRequestsGauge, m.queueLengthGauge)
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	go m.collectLoop()
	return nil
}

// Stop 停止负载监控
func (m *LoadMonitor) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

// collectLoop 收集循环
func (m *LoadMonitor) collectLoop() {
	ticker := time.NewTicker(m.collectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			// 可以在这里读取 /proc/loadavg 获取系统负载
		}
	}
}

// collectMetrics 收集指标（可观察指标回调）
func (m *LoadMonitor) collectMetrics(ctx context.Context, obs metric.Observer) error {
	// 获取系统负载（简化实现）
	loadAvg := m.getLoadAverage()
	obs.ObserveFloat64(m.loadAverageGauge, loadAvg)

	// 基于 Goroutine 数量估算并发请求数
	concurrent := int64(runtime.NumGoroutine())
	obs.ObserveInt64(m.concurrentRequestsGauge, concurrent)

	// 队列长度（简化实现，基于 Goroutine 数量）
	queueLength := concurrent
	obs.ObserveInt64(m.queueLengthGauge, queueLength)

	return nil
}

// getLoadAverage 获取系统负载（简化实现）
func (m *LoadMonitor) getLoadAverage() float64 {
	// 在 Linux 上可以读取 /proc/loadavg
	// 当前使用基于 Goroutine 的估算
	numGoroutines := runtime.NumGoroutine()
	cpus := runtime.NumCPU()
	
	// 简单的负载估算：Goroutine 数量 / CPU 核心数
	load := float64(numGoroutines) / float64(cpus)
	return load
}

// RecordRequest 记录请求
func (m *LoadMonitor) RecordRequest(ctx context.Context) {
	if m.enabled {
		m.requestRateCounter.Add(ctx, 1)
	}
}
