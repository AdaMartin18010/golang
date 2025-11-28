package system

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Monitor 系统资源监控器
// 提供 CPU、内存、IO、网络等系统资源的监控
type Monitor struct {
	meter           metric.Meter
	enabled         bool
	collectInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
	linuxCPUMonitor interface{} // Linux 平台的精确 CPU 监控（*LinuxCPUMonitor，仅在 Linux 上可用）

	// 指标
	cpuUsageGauge      metric.Float64ObservableGauge
	memoryUsageGauge   metric.Int64ObservableGauge
	memoryTotalGauge   metric.Int64ObservableGauge
	goroutineGauge     metric.Int64ObservableGauge
	gcCountCounter      metric.Int64Counter
	gcDurationHistogram metric.Float64Histogram
}

// Config 监控器配置
type Config struct {
	Meter           metric.Meter
	Enabled         bool
	CollectInterval time.Duration // 收集间隔（默认：5秒）
}

// NewMonitor 创建系统资源监控器
func NewMonitor(cfg Config) (*Monitor, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	collectInterval := cfg.CollectInterval
	if collectInterval == 0 {
		collectInterval = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitor := &Monitor{
		meter:           cfg.Meter,
		enabled:         cfg.Enabled,
		collectInterval: collectInterval,
		ctx:             ctx,
		cancel:          cancel,
	}

	// 在 Linux 上初始化精确 CPU 监控
	// 注意：NewLinuxCPUMonitor 仅在 Linux 上可用（通过 build tag）
	if runtime.GOOS == "linux" {
		monitor.linuxCPUMonitor = initLinuxCPUMonitor()
	}

	// 初始化指标
	if err := monitor.initMetrics(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to init metrics: %w", err)
	}

	return monitor, nil
}

// initMetrics 初始化指标
func (m *Monitor) initMetrics() error {
	var err error

	// CPU 使用率
	m.cpuUsageGauge, err = m.meter.Float64ObservableGauge(
		"system.cpu.usage",
		metric.WithDescription("CPU usage percentage (0-100)"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return err
	}

	// 内存使用量
	m.memoryUsageGauge, err = m.meter.Int64ObservableGauge(
		"system.memory.usage",
		metric.WithDescription("Memory usage in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	// 内存总量
	m.memoryTotalGauge, err = m.meter.Int64ObservableGauge(
		"system.memory.total",
		metric.WithDescription("Total memory in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	// Goroutine 数量
	m.goroutineGauge, err = m.meter.Int64ObservableGauge(
		"system.goroutines",
		metric.WithDescription("Number of goroutines"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	// GC 次数
	m.gcCountCounter, err = m.meter.Int64Counter(
		"system.gc.count",
		metric.WithDescription("Number of GC cycles"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	// GC 持续时间
	m.gcDurationHistogram, err = m.meter.Float64Histogram(
		"system.gc.duration",
		metric.WithDescription("GC duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Start 启动监控器
func (m *Monitor) Start() error {
	if !m.enabled {
		return nil
	}

	// 注册可观察指标的回调
	_, err := m.meter.RegisterCallback(m.collectMetrics, m.cpuUsageGauge, m.memoryUsageGauge, m.memoryTotalGauge, m.goroutineGauge)
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	// 启动后台收集协程
	go m.collectLoop()

	return nil
}

// Stop 停止监控器
func (m *Monitor) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

// collectLoop 收集循环
func (m *Monitor) collectLoop() {
	ticker := time.NewTicker(m.collectInterval)
	defer ticker.Stop()

	var lastGCStats runtime.MemStats
	runtime.ReadMemStats(&lastGCStats)

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectGCStats(&lastGCStats)
		}
	}
}

// collectMetrics 收集指标（可观察指标回调）
func (m *Monitor) collectMetrics(ctx context.Context, obs metric.Observer) error {
	// CPU 使用率（简化实现，实际需要更精确的计算）
	cpuUsage := m.getCPUUsage()

	// 内存统计
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 记录指标
	obs.ObserveFloat64(m.cpuUsageGauge, cpuUsage)
	obs.ObserveInt64(m.memoryUsageGauge, int64(memStats.Alloc))
	obs.ObserveInt64(m.memoryTotalGauge, int64(memStats.Sys))
	obs.ObserveInt64(m.goroutineGauge, int64(runtime.NumGoroutine()))

	return nil
}

// collectGCStats 收集 GC 统计
func (m *Monitor) collectGCStats(lastStats *runtime.MemStats) {
	var currentStats runtime.MemStats
	runtime.ReadMemStats(&currentStats)

	// 计算 GC 次数增量
	gcCount := int64(currentStats.NumGC - lastStats.NumGC)
	if gcCount > 0 {
		m.gcCountCounter.Add(context.Background(), gcCount)
	}

	// 记录 GC 持续时间
	if currentStats.NumGC > lastStats.NumGC {
		// 获取最近的 GC 持续时间
		lastGCIndex := (currentStats.NumGC - 1) % uint32(len(currentStats.PauseNs))
		if lastGCIndex < uint32(len(currentStats.PauseNs)) {
			pauseNs := currentStats.PauseNs[lastGCIndex]
			m.gcDurationHistogram.Record(context.Background(), float64(pauseNs)/1e9) // 转换为秒
		}
	}

	*lastStats = currentStats
}

// getCPUUsage 获取 CPU 使用率
// 在 Linux 上使用精确方法，其他平台使用简化实现
func (m *Monitor) getCPUUsage() float64 {
	// 尝试使用 Linux 精确方法
	if m.linuxCPUMonitor != nil && runtime.GOOS == "linux" {
		if usage := getLinuxCPUUsage(m.linuxCPUMonitor); usage >= 0 {
			return usage
		}
	}

	// 回退到简化实现：基于 Goroutine 数量和系统负载估算
	numGoroutines := runtime.NumGoroutine()
	
	// 简单的启发式：基于 Goroutine 数量估算 CPU 使用率
	// 这不是精确的，但可以作为近似值
	usage := float64(numGoroutines) * 0.1
	if usage > 100 {
		usage = 100
	}
	return usage
}

// GetMemoryStats 获取内存统计
func (m *Monitor) GetMemoryStats() MemoryStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return MemoryStats{
		Alloc:      memStats.Alloc,
		TotalAlloc: memStats.TotalAlloc,
		Sys:        memStats.Sys,
		NumGC:      memStats.NumGC,
		HeapAlloc:  memStats.HeapAlloc,
		HeapSys:    memStats.HeapSys,
	}
}

// MemoryStats 内存统计
type MemoryStats struct {
	Alloc      uint64 // 当前分配的内存
	TotalAlloc uint64 // 累计分配的内存
	Sys        uint64 // 系统内存
	NumGC      uint32 // GC 次数
	HeapAlloc  uint64 // 堆内存分配
	HeapSys    uint64 // 堆内存系统
}

// GetGoroutineCount 获取 Goroutine 数量
func (m *Monitor) GetGoroutineCount() int {
	return runtime.NumGoroutine()
}
