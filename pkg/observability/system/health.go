package system

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Healthy     bool
	Timestamp   time.Time
	MemoryUsage float64 // 内存使用率（0-100）
	CPUUsage    float64 // CPU 使用率（0-100）
	Goroutines  int
	GC          int
	Message     string
}

// HealthChecker 健康检查器
type HealthChecker struct {
	monitor      *Monitor
	thresholds   HealthThresholds
	lastCheck    time.Time
	checkInterval time.Duration
}

// HealthThresholds 健康阈值
type HealthThresholds struct {
	MaxMemoryUsage float64 // 最大内存使用率（0-100）
	MaxCPUUsage    float64 // 最大 CPU 使用率（0-100）
	MaxGoroutines  int     // 最大 Goroutine 数量
	MinGCInterval time.Duration // 最小 GC 间隔
}

// DefaultHealthThresholds 返回默认健康阈值
func DefaultHealthThresholds() HealthThresholds {
	return HealthThresholds{
		MaxMemoryUsage: 90.0,  // 90% 内存使用率
		MaxCPUUsage:    95.0,  // 95% CPU 使用率
		MaxGoroutines:  10000, // 10000 个 Goroutine
		MinGCInterval:  1 * time.Second,
	}
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(monitor *Monitor, thresholds HealthThresholds) *HealthChecker {
	return &HealthChecker{
		monitor:       monitor,
		thresholds:    thresholds,
		checkInterval: 5 * time.Second,
	}
}

// Check 执行健康检查
func (hc *HealthChecker) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Timestamp: time.Now(),
		Healthy:   true,
	}

	// 获取内存统计
	memStats := hc.monitor.GetMemoryStats()
	totalMem := memStats.Sys
	if totalMem > 0 {
		status.MemoryUsage = float64(memStats.Alloc) / float64(totalMem) * 100.0
	}

	// 获取 Goroutine 数量
	status.Goroutines = hc.monitor.GetGoroutineCount()
	status.GC = int(memStats.NumGC)

	// 检查内存使用率
	if status.MemoryUsage > hc.thresholds.MaxMemoryUsage {
		status.Healthy = false
		status.Message = fmt.Sprintf("memory usage too high: %.2f%%", status.MemoryUsage)
		return status
	}

	// 检查 Goroutine 数量
	if status.Goroutines > hc.thresholds.MaxGoroutines {
		status.Healthy = false
		status.Message = fmt.Sprintf("too many goroutines: %d", status.Goroutines)
		return status
	}

	// 检查 CPU 使用率（简化实现）
	// 实际应该使用更精确的方法
	var memStats2 runtime.MemStats
	runtime.ReadMemStats(&memStats2)
	// 基于 Goroutine 数量估算 CPU 使用率
	status.CPUUsage = float64(status.Goroutines) * 0.1
	if status.CPUUsage > 100 {
		status.CPUUsage = 100
	}

	if status.CPUUsage > hc.thresholds.MaxCPUUsage {
		status.Healthy = false
		status.Message = fmt.Sprintf("CPU usage too high: %.2f%%", status.CPUUsage)
		return status
	}

	status.Message = "healthy"
	return status
}

// CheckPeriodically 定期执行健康检查
func (hc *HealthChecker) CheckPeriodically(ctx context.Context, callback func(HealthStatus)) {
	ticker := time.NewTicker(hc.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			status := hc.Check(ctx)
			if callback != nil {
				callback(status)
			}
		}
	}
}

// IsHealthy 检查是否健康
func (hc *HealthChecker) IsHealthy(ctx context.Context) bool {
	status := hc.Check(ctx)
	return status.Healthy
}
