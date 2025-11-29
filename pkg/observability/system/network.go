package system

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// NetworkMonitor 网络监控器
type NetworkMonitor struct {
	meter           metric.Meter
	enabled         bool
	collectInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc

	// 指标
	bytesSentCounter   metric.Int64Counter
	bytesRecvCounter   metric.Int64Counter
	packetsSentCounter metric.Int64Counter
	packetsRecvCounter metric.Int64Counter
	connectionsGauge   metric.Int64ObservableGauge
}

// NewNetworkMonitor 创建网络监控器
func NewNetworkMonitor(cfg Config) (*NetworkMonitor, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	collectInterval := cfg.CollectInterval
	if collectInterval == 0 {
		collectInterval = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitor := &NetworkMonitor{
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
func (m *NetworkMonitor) initMetrics() error {
	var err error

	m.bytesSentCounter, err = m.meter.Int64Counter(
		"system.network.bytes.sent",
		metric.WithDescription("Total bytes sent"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.bytesRecvCounter, err = m.meter.Int64Counter(
		"system.network.bytes.received",
		metric.WithDescription("Total bytes received"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	m.packetsSentCounter, err = m.meter.Int64Counter(
		"system.network.packets.sent",
		metric.WithDescription("Total packets sent"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.packetsRecvCounter, err = m.meter.Int64Counter(
		"system.network.packets.received",
		metric.WithDescription("Total packets received"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	m.connectionsGauge, err = m.meter.Int64ObservableGauge(
		"system.network.connections",
		metric.WithDescription("Number of network connections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Start 启动网络监控
func (m *NetworkMonitor) Start() error {
	if !m.enabled {
		return nil
	}

	// 注册可观察指标回调
	_, err := m.meter.RegisterCallback(m.collectMetrics, m.connectionsGauge)
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	go m.collectLoop()
	return nil
}

// Stop 停止网络监控
func (m *NetworkMonitor) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

// collectLoop 收集循环
func (m *NetworkMonitor) collectLoop() {
	ticker := time.NewTicker(m.collectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			// 可以在这里读取 /proc/net/sockstat 或使用 netlink
			// 当前为占位实现
		}
	}
}

// collectMetrics 收集指标（可观察指标回调）
func (m *NetworkMonitor) collectMetrics(ctx context.Context, obs metric.Observer) error {
	// 获取连接数（简化实现）
	connCount := m.getConnectionCount()
	obs.ObserveInt64(m.connectionsGauge, connCount)
	return nil
}

// getConnectionCount 获取连接数（简化实现）
func (m *NetworkMonitor) getConnectionCount() int64 {
	// 简化实现：尝试统计本地监听端口
	// 实际应该读取 /proc/net/sockstat 或使用 netlink
	count := int64(0)
	
	// 尝试统计常见的监听端口
	ports := []string{"80", "443", "8080", "9090"}
	for _, port := range ports {
		if conn, err := net.Listen("tcp", ":"+port); err == nil {
			conn.Close()
			count++
		}
	}
	
	return count
}

// RecordBytesSent 记录发送的字节数
func (m *NetworkMonitor) RecordBytesSent(ctx context.Context, bytes int64) {
	if m.enabled {
		m.bytesSentCounter.Add(ctx, bytes)
		m.packetsSentCounter.Add(ctx, 1)
	}
}

// RecordBytesReceived 记录接收的字节数
func (m *NetworkMonitor) RecordBytesReceived(ctx context.Context, bytes int64) {
	if m.enabled {
		m.bytesRecvCounter.Add(ctx, bytes)
		m.packetsRecvCounter.Add(ctx, 1)
	}
}
