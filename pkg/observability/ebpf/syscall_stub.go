//go:build !linux
// +build !linux

package ebpf

// 非 Linux 系统的 stub 实现
// 这些函数在非 Linux 平台上不会实际加载 eBPF 程序

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// SyscallTracerConfig 系统调用追踪器配置
type SyscallTracerConfig struct {
	Tracer  trace.Tracer
	Meter   metric.Meter
	Enabled bool
}

// SyscallTracer 系统调用追踪器（非 Linux 系统的 stub 实现）
type SyscallTracer struct {
	config SyscallTracerConfig
}

// NewSyscallTracer 创建系统调用追踪器（stub）
func NewSyscallTracer(config SyscallTracerConfig) (*SyscallTracer, error) {
	if !config.Enabled {
		return &SyscallTracer{config: config}, nil
	}
	return nil, errors.New("eBPF syscall tracer is only supported on Linux")
}

// Start 开始追踪（stub）
func (t *SyscallTracer) Start(ctx context.Context) error {
	return errors.New("eBPF syscall tracer is only supported on Linux")
}

// Stop 停止追踪（stub）
func (t *SyscallTracer) Stop() error {
	return nil
}

// IsEnabled 检查是否启用
func (t *SyscallTracer) IsEnabled() bool {
	return t.config.Enabled
}

// GetStats 获取统计信息（stub）
func (t *SyscallTracer) GetStats() map[string]interface{} {
	return map[string]interface{}{}
}

// GetSyscallStats 获取系统调用统计（stub）
func (t *SyscallTracer) GetSyscallStats(ctx context.Context) (map[uint64]uint64, error) {
	return nil, errors.New("eBPF syscall tracer is only supported on Linux")
}
