//go:build !linux
// +build !linux

package ebpf

// 非 Linux 系统的 stub 实现
// 这些函数在非 Linux 平台上不会实际加载 eBPF 程序

import (
	"context"
	"errors"
	"net"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// NetworkTracerConfig 网络追踪器配置
type NetworkTracerConfig struct {
	Tracer        trace.Tracer
	Meter         metric.Meter
	Enabled       bool
	TrackInbound  bool
	TrackOutbound bool
}

// ConnectionDetail 连接详情
type ConnectionDetail struct {
	SourceIP      net.IP
	SourcePort    uint16
	DestIP        net.IP
	DestPort      uint16
	State         string
	BytesSent     uint64
	BytesReceived uint64
}

// NetworkTracer 网络追踪器（非 Linux 系统的 stub 实现）
type NetworkTracer struct {
	config NetworkTracerConfig
}

// NewNetworkTracer 创建网络追踪器（stub）
func NewNetworkTracer(config NetworkTracerConfig) (*NetworkTracer, error) {
	if !config.Enabled {
		return &NetworkTracer{config: config}, nil
	}
	return nil, errors.New("eBPF network tracer is only supported on Linux")
}

// Start 开始追踪（stub）
func (t *NetworkTracer) Start(ctx context.Context) error {
	return errors.New("eBPF network tracer is only supported on Linux")
}

// Stop 停止追踪（stub）
func (t *NetworkTracer) Stop() error {
	return nil
}

// IsEnabled 检查是否启用
func (t *NetworkTracer) IsEnabled() bool {
	return t.config.Enabled
}

// GetStats 获取统计信息（stub）
func (t *NetworkTracer) GetStats() map[string]interface{} {
	return map[string]interface{}{}
}

// GetActiveConnections 获取活跃连接数（stub）
func (t *NetworkTracer) GetActiveConnections(ctx context.Context) (int64, error) {
	return 0, errors.New("eBPF network tracer is only supported on Linux")
}

// GetConnectionStats 获取网络连接统计（stub）
func (t *NetworkTracer) GetConnectionStats(ctx context.Context) (map[uint32]uint64, error) {
	return nil, errors.New("eBPF network tracer is only supported on Linux")
}

// GetConnectionDetails 获取连接详细信息（stub）
func (t *NetworkTracer) GetConnectionDetails(ctx context.Context) ([]ConnectionDetail, error) {
	return nil, errors.New("eBPF network tracer is only supported on Linux")
}
