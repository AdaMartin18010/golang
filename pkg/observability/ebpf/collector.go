package ebpf

import (
	"context"
	"fmt"
	"runtime"
	"sync"
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

	// 运行时状态
	started    bool
	startMutex sync.Mutex
	wg         sync.WaitGroup
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

	// 非 Linux 系统不支持 eBPF
	if runtime.GOOS != "linux" {
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

	c.startMutex.Lock()
	defer c.startMutex.Unlock()

	if c.started {
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

	c.started = true

	// 启动定期收集任务
	c.wg.Add(1)
	go c.collectLoop()

	return nil
}

// collectLoop 定期收集循环
func (c *Collector) collectLoop() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.collectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.collect()
		}
	}
}

// collect 执行一次数据收集
func (c *Collector) collect() {
	// 这里可以添加定期收集的逻辑
	// 例如：从 eBPF maps 读取聚合数据、更新指标等
}

// GetSyscallTracer 获取系统调用追踪器
func (c *Collector) GetSyscallTracer() *SyscallTracer {
	return c.syscallTracer
}

// GetNetworkTracer 获取网络追踪器
func (c *Collector) GetNetworkTracer() *NetworkTracer {
	return c.networkTracer
}

// GetSyscallStats 获取系统调用统计
func (c *Collector) GetSyscallStats(ctx context.Context) (map[uint64]uint64, error) {
	if c.syscallTracer == nil {
		return nil, nil
	}
	return c.syscallTracer.GetSyscallStats(ctx)
}

// GetActiveConnections 获取活跃连接数
func (c *Collector) GetActiveConnections(ctx context.Context) (int64, error) {
	if c.networkTracer == nil {
		return 0, nil
	}
	return c.networkTracer.GetActiveConnections(ctx)
}

// GetConnectionStats 获取网络连接统计
func (c *Collector) GetConnectionStats(ctx context.Context) (map[uint32]uint64, error) {
	if c.networkTracer == nil {
		return nil, nil
	}
	return c.networkTracer.GetConnectionStats(ctx)
}

// GetConnectionDetails 获取连接详细信息
func (c *Collector) GetConnectionDetails(ctx context.Context) ([]ConnectionDetail, error) {
	if c.networkTracer == nil {
		return nil, nil
	}
	return c.networkTracer.GetConnectionDetails(ctx)
}

// Stop 停止收集器
// 停止所有子追踪器并清理资源
func (c *Collector) Stop() error {
	if !c.enabled {
		return nil
	}

	c.startMutex.Lock()
	defer c.startMutex.Unlock()

	if !c.started {
		return nil
	}

	// 取消上下文，停止收集循环
	if c.cancel != nil {
		c.cancel()
	}

	// 等待收集循环结束
	c.wg.Wait()

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

	c.started = false
	return nil
}

// IsEnabled 检查是否启用
func (c *Collector) IsEnabled() bool {
	return c.enabled
}

// IsStarted 检查是否已启动
func (c *Collector) IsStarted() bool {
	c.startMutex.Lock()
	defer c.startMutex.Unlock()
	return c.started
}

// Enable 启用收集器
func (c *Collector) Enable() {
	c.enabled = true
}

// Disable 禁用收集器
func (c *Collector) Disable() {
	c.enabled = false
	if c.syscallTracer != nil {
		c.syscallTracer.Stop()
	}
	if c.networkTracer != nil {
		c.networkTracer.Stop()
	}
}

// Close 关闭收集器（别名方法）
func (c *Collector) Close() error {
	return c.Stop()
}

// Status 返回收集器状态
type Status struct {
	Enabled                 bool          `json:"enabled"`
	Started                 bool          `json:"started"`
	SyscallTrackingEnabled  bool          `json:"syscall_tracking_enabled"`
	NetworkMonitoringEnabled bool         `json:"network_monitoring_enabled"`
	SyscallTracerEnabled    bool          `json:"syscall_tracer_enabled"`
	NetworkTracerEnabled    bool          `json:"network_tracer_enabled"`
	CollectInterval         time.Duration `json:"collect_interval"`
}

// GetStatus 获取收集器状态
func (c *Collector) GetStatus() Status {
	c.startMutex.Lock()
	defer c.startMutex.Unlock()

	return Status{
		Enabled:                  c.enabled,
		Started:                  c.started,
		SyscallTrackingEnabled:   c.enableSyscallTracking,
		NetworkMonitoringEnabled: c.enableNetworkMonitoring,
		SyscallTracerEnabled:     c.syscallTracer != nil && c.syscallTracer.IsEnabled(),
		NetworkTracerEnabled:     c.networkTracer != nil && c.networkTracer.IsEnabled(),
		CollectInterval:          c.collectInterval,
	}
}
