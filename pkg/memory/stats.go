package memory

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// =============================================================================
// 内存统计工具 - Memory Statistics Tools
// =============================================================================

// MemoryStats 内存统计信息
type MemoryStats struct {
	// 堆内存统计
	Alloc        uint64 // 当前分配的字节数
	TotalAlloc   uint64 // 累计分配的字节数
	Sys          uint64 // 从系统获取的字节数
	Lookups      uint64 // 指针查找次数
	Mallocs      uint64 // 分配次数
	Frees        uint64 // 释放次数
	HeapAlloc    uint64 // 堆上分配的字节数
	HeapSys      uint64 // 从系统获取的堆字节数
	HeapIdle     uint64 // 空闲堆字节数
	HeapInuse    uint64 // 使用中的堆字节数
	HeapReleased uint64 // 释放给OS的字节数
	HeapObjects  uint64 // 堆上分配的对象数

	// GC统计
	NumGC         uint32        // GC运行次数
	PauseTotal    time.Duration // GC暂停总时间
	LastGC        time.Time     // 上次GC时间
	GCCPUFraction float64       // GC占用的CPU比例

	// 计算字段
	AllocRate float64 // 分配速率 (MB/s)
	GCRate    float64 // GC频率 (次/秒)
	HeapUsage float64 // 堆使用率 (%)
}

// MemoryMonitor 内存监控器
type MemoryMonitor struct {
	mu         sync.RWMutex
	samples    []MemorySample
	maxSamples int
	startTime  time.Time
	lastStats  runtime.MemStats
}

// MemorySample 内存采样点
type MemorySample struct {
	Timestamp time.Time
	Stats     MemoryStats
}

// NewMemoryMonitor 创建内存监控器
func NewMemoryMonitor(maxSamples int) *MemoryMonitor {
	return &MemoryMonitor{
		samples:    make([]MemorySample, 0, maxSamples),
		maxSamples: maxSamples,
		startTime:  time.Now(),
	}
}

// Collect 收集当前内存统计
func (mm *MemoryMonitor) Collect() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := MemoryStats{
		Alloc:         m.Alloc,
		TotalAlloc:    m.TotalAlloc,
		Sys:           m.Sys,
		Lookups:       m.Lookups,
		Mallocs:       m.Mallocs,
		Frees:         m.Frees,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapIdle:      m.HeapIdle,
		HeapInuse:     m.HeapInuse,
		HeapReleased:  m.HeapReleased,
		HeapObjects:   m.HeapObjects,
		NumGC:         m.NumGC,
		PauseTotal:    time.Duration(m.PauseTotalNs),
		LastGC:        time.Unix(0, int64(m.LastGC)),
		GCCPUFraction: m.GCCPUFraction,
	}

	// 计算派生指标
	duration := time.Since(mm.startTime).Seconds()
	if duration > 0 {
		stats.AllocRate = float64(m.TotalAlloc) / duration / 1024 / 1024 // MB/s
		stats.GCRate = float64(m.NumGC) / duration                       // 次/秒
	}

	if m.HeapSys > 0 {
		stats.HeapUsage = float64(m.HeapInuse) / float64(m.HeapSys) * 100
	}

	mm.mu.Lock()
	defer mm.mu.Unlock()

	// 添加样本
	sample := MemorySample{
		Timestamp: time.Now(),
		Stats:     stats,
	}

	mm.samples = append(mm.samples, sample)

	// 限制样本数量
	if len(mm.samples) > mm.maxSamples {
		mm.samples = mm.samples[1:]
	}

	mm.lastStats = m

	return stats
}

// GetSamples 获取所有样本
func (mm *MemoryMonitor) GetSamples() []MemorySample {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	// 返回副本
	samples := make([]MemorySample, len(mm.samples))
	copy(samples, mm.samples)
	return samples
}

// GetTrend 获取内存趋势
func (mm *MemoryMonitor) GetTrend() MemoryTrend {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if len(mm.samples) == 0 {
		return MemoryTrend{}
	}

	var trend MemoryTrend
	trend.Samples = len(mm.samples)

	// 计算平均值和趋势
	var totalAlloc, totalHeapInuse uint64
	var maxAlloc, maxHeapInuse uint64

	for _, sample := range mm.samples {
		totalAlloc += sample.Stats.Alloc
		totalHeapInuse += sample.Stats.HeapInuse

		if sample.Stats.Alloc > maxAlloc {
			maxAlloc = sample.Stats.Alloc
		}
		if sample.Stats.HeapInuse > maxHeapInuse {
			maxHeapInuse = sample.Stats.HeapInuse
		}
	}

	trend.AvgAlloc = totalAlloc / uint64(len(mm.samples))
	trend.AvgHeapInuse = totalHeapInuse / uint64(len(mm.samples))
	trend.MaxAlloc = maxAlloc
	trend.MaxHeapInuse = maxHeapInuse

	// 计算增长率
	if len(mm.samples) >= 2 {
		first := mm.samples[0].Stats
		last := mm.samples[len(mm.samples)-1].Stats

		duration := last.LastGC.Sub(first.LastGC).Seconds()
		if duration > 0 {
			allocGrowth := float64(last.Alloc) - float64(first.Alloc)
			trend.AllocGrowthRate = allocGrowth / duration / 1024 / 1024 // MB/s
		}
	}

	return trend
}

// MemoryTrend 内存趋势
type MemoryTrend struct {
	Samples         int
	AvgAlloc        uint64
	AvgHeapInuse    uint64
	MaxAlloc        uint64
	MaxHeapInuse    uint64
	AllocGrowthRate float64 // MB/s
}

// Reset 重置监控器
func (mm *MemoryMonitor) Reset() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.samples = mm.samples[:0]
	mm.startTime = time.Now()
}

// =============================================================================
// 内存分析器 - Memory Profiler
// =============================================================================

// MemoryProfiler 内存分析器
type MemoryProfiler struct {
	monitor     *MemoryMonitor
	interval    time.Duration
	stopChan    chan struct{}
	running     bool
	mu          sync.Mutex
	onThreshold func(MemoryStats)
	threshold   uint64
}

// NewMemoryProfiler 创建内存分析器
func NewMemoryProfiler(interval time.Duration) *MemoryProfiler {
	return &MemoryProfiler{
		monitor:  NewMemoryMonitor(1000), // 保留最近1000个样本
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start 开始分析
func (mp *MemoryProfiler) Start() {
	mp.mu.Lock()
	if mp.running {
		mp.mu.Unlock()
		return
	}
	mp.running = true
	mp.mu.Unlock()

	go mp.run()
}

// Stop 停止分析
func (mp *MemoryProfiler) Stop() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if !mp.running {
		return
	}

	close(mp.stopChan)
	mp.running = false
}

// run 运行分析循环
func (mp *MemoryProfiler) run() {
	ticker := time.NewTicker(mp.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := mp.monitor.Collect()

			// 检查阈值
			if mp.onThreshold != nil && mp.threshold > 0 {
				if stats.HeapInuse > mp.threshold {
					mp.onThreshold(stats)
				}
			}

		case <-mp.stopChan:
			return
		}
	}
}

// SetThreshold 设置阈值回调
func (mp *MemoryProfiler) SetThreshold(threshold uint64, callback func(MemoryStats)) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.threshold = threshold
	mp.onThreshold = callback
}

// GetMonitor 获取监控器
func (mp *MemoryProfiler) GetMonitor() *MemoryMonitor {
	return mp.monitor
}

// Report 生成报告
func (mp *MemoryProfiler) Report() string {
	stats := mp.monitor.Collect()
	trend := mp.monitor.GetTrend()

	report := fmt.Sprintf(`
Memory Profile Report
=====================

Current Stats:
  Allocated:     %s
  Heap In Use:   %s
  Heap System:   %s
  Heap Objects:  %d
  Heap Usage:    %.2f%%

GC Stats:
  NumGC:         %d
  Pause Total:   %v
  GC CPU:        %.2f%%
  Last GC:       %v

Performance:
  Alloc Rate:    %.2f MB/s
  GC Rate:       %.2f /s

Trend (last %d samples):
  Avg Alloc:     %s
  Max Alloc:     %s
  Growth Rate:   %.2f MB/s
`,
		formatBytes(stats.Alloc),
		formatBytes(stats.HeapInuse),
		formatBytes(stats.HeapSys),
		stats.HeapObjects,
		stats.HeapUsage,
		stats.NumGC,
		stats.PauseTotal,
		stats.GCCPUFraction*100,
		stats.LastGC.Format("15:04:05"),
		stats.AllocRate,
		stats.GCRate,
		trend.Samples,
		formatBytes(trend.AvgAlloc),
		formatBytes(trend.MaxAlloc),
		trend.AllocGrowthRate,
	)

	return report
}

// formatBytes 格式化字节数
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// =============================================================================
// 全局分析器实例
// =============================================================================

var (
	// GlobalProfiler 全局内存分析器
	GlobalProfiler = NewMemoryProfiler(1 * time.Second)
)

// StartProfiling 开始全局内存分析
func StartProfiling(interval time.Duration) {
	if interval > 0 {
		GlobalProfiler.interval = interval
	}
	GlobalProfiler.Start()
}

// StopProfiling 停止全局内存分析
func StopProfiling() {
	GlobalProfiler.Stop()
}

// GetMemoryReport 获取内存报告
func GetMemoryReport() string {
	return GlobalProfiler.Report()
}

// CollectMemoryStats 收集当前内存统计
func CollectMemoryStats() MemoryStats {
	return GlobalProfiler.monitor.Collect()
}
