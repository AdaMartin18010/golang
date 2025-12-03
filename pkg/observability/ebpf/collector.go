package ebpf

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Collector eBPF 收集器
// 提供基于 eBPF 的系统级可观测性数据收集
// 使用 Cilium eBPF 库实现真正的内核级监控
type Collector struct {
	tracer                     trace.Tracer
	meter                      metric.Meter
	enabled                    bool
	ctx                        context.Context
	cancel                     context.CancelFunc
	collectInterval            time.Duration
	enableSyscallTracking      bool
	enableNetworkMonitoring    bool
	enablePerformanceProfiling bool

	// 子追踪器
	syscallTracer *SyscallTracer
	networkTracer *NetworkTracer

	// 指标（保留用于兼容性）
	syscallCounter          metric.Int64Counter
	networkPacketCounter    metric.Int64Counter
	syscallLatencyHistogram metric.Float64Histogram
}

// Config 配置
type Config struct {
	Tracer  trace.Tracer
	Meter   metric.Meter
	Enabled bool
	// CollectInterval 收集间隔（默认：5秒）
	CollectInterval time.Duration
	// EnableSyscallTracking 是否启用系统调用追踪
	EnableSyscallTracking bool
	// EnableNetworkMonitoring 是否启用网络监控
	EnableNetworkMonitoring bool
	// EnablePerformanceProfiling 是否启用性能分析
	EnablePerformanceProfiling bool
}

// NewCollector 创建 eBPF 收集器
// 使用 Cilium eBPF 库实现真正的系统级监控
func NewCollector(cfg Config) (*Collector, error) {
	if !cfg.Enabled {
		return &Collector{enabled: false}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	collectInterval := cfg.CollectInterval
	if collectInterval == 0 {
		collectInterval = 5 * time.Second
	}

	collector := &Collector{
		tracer:                     cfg.Tracer,
		meter:                      cfg.Meter,
		enabled:                    true,
		ctx:                        ctx,
		cancel:                     cancel,
		collectInterval:            collectInterval,
		enableSyscallTracking:      cfg.EnableSyscallTracking,
		enableNetworkMonitoring:    cfg.EnableNetworkMonitoring,
		enablePerformanceProfiling: cfg.EnablePerformanceProfiling,
	}

	// 创建系统调用追踪器
	if cfg.EnableSyscallTracking {
		syscallTracer, err := NewSyscallTracer(SyscallTracerConfig{
			Tracer:  cfg.Tracer,
			Meter:   cfg.Meter,
			Enabled: true,
		})
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create syscall tracer: %w", err)
		}
		collector.syscallTracer = syscallTracer
	}

	// 创建网络追踪器
	if cfg.EnableNetworkMonitoring {
		networkTracer, err := NewNetworkTracer(NetworkTracerConfig{
			Tracer:        cfg.Tracer,
			Meter:         cfg.Meter,
			Enabled:       true,
			TrackInbound:  true,
			TrackOutbound: true,
		})
		if err != nil {
			cancel()
			if collector.syscallTracer != nil {
				collector.syscallTracer.Stop()
			}
			return nil, fmt.Errorf("failed to create network tracer: %w", err)
		}
		collector.networkTracer = networkTracer
	}

	return collector, nil
}

// Start 启动收集器
// 启动所有子追踪器
func (c *Collector) Start() error {
	if !c.enabled {
		return nil
	}

	// 启动系统调用追踪器
	if c.syscallTracer != nil {
		if err := c.syscallTracer.Start(); err != nil {
			return fmt.Errorf("failed to start syscall tracer: %w", err)
		}
	}

	// 启动网络追踪器
	if c.networkTracer != nil {
		if err := c.networkTracer.Start(); err != nil {
			// 清理已启动的追踪器
			if c.syscallTracer != nil {
				c.syscallTracer.Stop()
			}
			return fmt.Errorf("failed to start network tracer: %w", err)
		}
	}

	return nil
}

// GetSyscallTracer 获取系统调用追踪器
func (c *Collector) GetSyscallTracer() *SyscallTracer {
	return c.syscallTracer
}

// GetNetworkTracer 获取网络追踪器
func (c *Collector) GetNetworkTracer() *NetworkTracer {
	return c.networkTracer
}

// Stop 停止收集器
// 停止所有子追踪器并清理资源
func (c *Collector) Stop() error {
	if !c.enabled {
		return nil
	}

	if c.cancel != nil {
		c.cancel()
	}

	// 停止系统调用追踪器
	if c.syscallTracer != nil {
		if err := c.syscallTracer.Stop(); err != nil {
			return fmt.Errorf("failed to stop syscall tracer: %w", err)
		}
	}

	// 停止网络追踪器
	if c.networkTracer != nil {
		if err := c.networkTracer.Stop(); err != nil {
			return fmt.Errorf("failed to stop network tracer: %w", err)
		}
	}

	return nil
}

// IsEnabled 检查是否启用
func (c *Collector) IsEnabled() bool {
	return c.enabled
}

// Enable 启用收集器
func (c *Collector) Enable() {
	c.enabled = true
}

// Disable 禁用收集器
func (c *Collector) Disable() {
	c.enabled = false
	if c.syscallTracer != nil {
		c.syscallTracer.Disable()
	}
}
