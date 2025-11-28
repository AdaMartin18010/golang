package system

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// DiskMonitor 磁盘监控器
type DiskMonitor struct {
	meter           metric.Meter
	enabled         bool
	collectInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc

	// 指标
	diskUsageGauge      metric.Int64ObservableGauge
	diskTotalGauge      metric.Int64ObservableGauge
	diskAvailableGauge  metric.Int64ObservableGauge
	diskReadBytesCounter  metric.Int64Counter
	diskWriteBytesCounter metric.Int64Counter
}

// NewDiskMonitor 创建磁盘监控器
func NewDiskMonitor(cfg Config) (*DiskMonitor, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	collectInterval := cfg.CollectInterval
	if collectInterval == 0 {
		collectInterval = 10 * time.Second // 磁盘监控间隔可以更长
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitor := &DiskMonitor{
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
func (m *DiskMonitor) initMetrics() error {
	var err error

	m.diskUsageGauge, err = m.meter.Int64ObservableGauge(
		"system.disk.usage",
		metric.WithDescription("Disk usage in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.diskTotalGauge, err = m.meter.Int64ObservableGauge(
		"system.disk.total",
		metric.WithDescription("Total disk space in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.diskAvailableGauge, err = m.meter.Int64ObservableGauge(
		"system.disk.available",
		metric.WithDescription("Available disk space in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.diskReadBytesCounter, err = m.meter.Int64Counter(
		"system.disk.read.bytes",
		metric.WithDescription("Total bytes read from disk"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.diskWriteBytesCounter, err = m.meter.Int64Counter(
		"system.disk.write.bytes",
		metric.WithDescription("Total bytes written to disk"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Start 启动磁盘监控
func (m *DiskMonitor) Start() error {
	if !m.enabled {
		return nil
	}

	// 注册可观察指标回调
	_, err := m.meter.RegisterCallback(m.collectMetrics, m.diskUsageGauge, m.diskTotalGauge, m.diskAvailableGauge)
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	go m.collectLoop()
	return nil
}

// Stop 停止磁盘监控
func (m *DiskMonitor) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

// collectLoop 收集循环
func (m *DiskMonitor) collectLoop() {
	ticker := time.NewTicker(m.collectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			// 可以在这里读取 /proc/diskstats 获取磁盘 IO 统计
			// 当前为占位实现
		}
	}
}

// collectMetrics 收集指标（可观察指标回调）
func (m *DiskMonitor) collectMetrics(ctx context.Context, obs metric.Observer) error {
	if runtime.GOOS == "windows" {
		return m.collectMetricsWindows(ctx, obs)
	}
	return m.collectMetricsUnix(ctx, obs)
}


// RecordRead 记录磁盘读取
func (m *DiskMonitor) RecordRead(ctx context.Context, bytes int64) {
	if m.enabled {
		m.diskReadBytesCounter.Add(ctx, bytes)
	}
}

// RecordWrite 记录磁盘写入
func (m *DiskMonitor) RecordWrite(ctx context.Context, bytes int64) {
	if m.enabled {
		m.diskWriteBytesCounter.Add(ctx, bytes)
	}
}
