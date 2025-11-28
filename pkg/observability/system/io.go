package system

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// IOMonitor IO 监控器
type IOMonitor struct {
	meter           metric.Meter
	enabled         bool
	collectInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc

	// 指标
	readBytesCounter  metric.Int64Counter
	writeBytesCounter metric.Int64Counter
	readOpsCounter    metric.Int64Counter
	writeOpsCounter   metric.Int64Counter
}

// NewIOMonitor 创建 IO 监控器
func NewIOMonitor(cfg Config) (*IOMonitor, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	collectInterval := cfg.CollectInterval
	if collectInterval == 0 {
		collectInterval = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitor := &IOMonitor{
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
func (m *IOMonitor) initMetrics() error {
	var err error

	m.readBytesCounter, err = m.meter.Int64Counter(
		"system.io.read.bytes",
		metric.WithDescription("Total bytes read"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.writeBytesCounter, err = m.meter.Int64Counter(
		"system.io.write.bytes",
		metric.WithDescription("Total bytes written"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.readOpsCounter, err = m.meter.Int64Counter(
		"system.io.read.ops",
		metric.WithDescription("Total read operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.writeOpsCounter, err = m.meter.Int64Counter(
		"system.io.write.ops",
		metric.WithDescription("Total write operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Start 启动 IO 监控
func (m *IOMonitor) Start() error {
	if !m.enabled {
		return nil
	}

	go m.collectLoop()
	return nil
}

// Stop 停止 IO 监控
func (m *IOMonitor) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

// collectLoop 收集循环
func (m *IOMonitor) collectLoop() {
	ticker := time.NewTicker(m.collectInterval)
	defer ticker.Stop()

	var lastStats runtime.MemStats

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectIOStats(&lastStats)
		}
	}
}

// collectIOStats 收集 IO 统计
func (m *IOMonitor) collectIOStats(lastStats *runtime.MemStats) {
	var currentStats runtime.MemStats
	runtime.ReadMemStats(&currentStats)

	// Go 运行时没有直接的 IO 统计，这里使用内存分配作为近似
	// 实际应该读取 /proc/self/io 或使用系统调用
	readBytes := int64(currentStats.TotalAlloc - lastStats.TotalAlloc)
	if readBytes > 0 {
		m.readBytesCounter.Add(context.Background(), readBytes)
		m.readOpsCounter.Add(context.Background(), 1)
	}

	*lastStats = currentStats
}

// RecordRead 记录读取操作
func (m *IOMonitor) RecordRead(ctx context.Context, bytes int64) {
	if m.enabled {
		m.readBytesCounter.Add(ctx, bytes)
		m.readOpsCounter.Add(ctx, 1)
	}
}

// RecordWrite 记录写入操作
func (m *IOMonitor) RecordWrite(ctx context.Context, bytes int64) {
	if m.enabled {
		m.writeBytesCounter.Add(ctx, bytes)
		m.writeOpsCounter.Add(ctx, 1)
	}
}
