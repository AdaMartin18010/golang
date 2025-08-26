package testing_system

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// PerformanceBenchmark 性能基准测试
type PerformanceBenchmark struct {
	Name        string
	Description string
	Run         func(ctx context.Context) (BenchmarkResult, error)
	Timeout     time.Duration
	Iterations  int
	Warmup      int
	Threshold   float64
	Baseline    *BenchmarkResult
	mu          sync.RWMutex
}

// BenchmarkResult 基准测试结果
type BenchmarkResult struct {
	Name           string            `json:"name"`
	Duration       time.Duration     `json:"duration"`
	Operations     int64             `json:"operations"`
	Throughput     float64           `json:"throughput"`
	MemoryUsage    MemoryUsage       `json:"memory_usage"`
	CPUUsage       CPUUsage          `json:"cpu_usage"`
	Iterations     int               `json:"iterations"`
	MinDuration    time.Duration     `json:"min_duration"`
	MaxDuration    time.Duration     `json:"max_duration"`
	AvgDuration    time.Duration     `json:"avg_duration"`
	StdDev         float64           `json:"std_dev"`
	Percentiles    map[int]time.Duration `json:"percentiles"`
	Timestamp      time.Time         `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// MemoryUsage 内存使用情况
type MemoryUsage struct {
	Allocated uint64  `json:"allocated"`
	Total     uint64  `json:"total"`
	Heap      uint64  `json:"heap"`
	Stack     uint64  `json:"stack"`
	GC        uint64  `json:"gc"`
}

// CPUUsage CPU使用情况
type CPUUsage struct {
	UserTime   time.Duration `json:"user_time"`
	SystemTime time.Duration `json:"system_time"`
	IdleTime   time.Duration `json:"idle_time"`
	Usage      float64       `json:"usage_percent"`
}

// PerformanceMonitor 性能监控器
type PerformanceMonitor struct {
	benchmarks map[string]*PerformanceBenchmark
	results    []BenchmarkResult
	detector   *RegressionDetector
	config     *PerformanceConfig
	mu         sync.RWMutex
}

// PerformanceConfig 性能测试配置
type PerformanceConfig struct {
	DefaultTimeout   time.Duration
	DefaultIterations int
	DefaultWarmup    int
	RegressionThreshold float64
	OutputDir         string
	ReportFormat      string
	EnableProfiling   bool
	ProfilingDir      string
}

// RegressionDetector 回归检测器
type RegressionDetector struct {
	baseline    map[string]BenchmarkResult
	current     map[string]BenchmarkResult
	threshold   float64
	alerts      []RegressionAlert
	mu          sync.RWMutex
}

// RegressionAlert 回归告警
type RegressionAlert struct {
	BenchmarkName string    `json:"benchmark_name"`
	Severity      string    `json:"severity"`
	Message       string    `json:"message"`
	Baseline      float64   `json:"baseline"`
	Current       float64   `json:"current"`
	Degradation   float64   `json:"degradation"`
	Timestamp     time.Time `json:"timestamp"`
}

// NewPerformanceBenchmark 创建性能基准测试
func NewPerformanceBenchmark(name, description string, run func(ctx context.Context) (BenchmarkResult, error)) *PerformanceBenchmark {
	return &PerformanceBenchmark{
		Name:        name,
		Description: description,
		Run:         run,
		Timeout:     60 * time.Second,
		Iterations:  100,
		Warmup:      10,
		Threshold:   0.1, // 10%性能下降阈值
		Metadata:    make(map[string]interface{}),
	}
}

// SetBaseline 设置基准结果
func (pb *PerformanceBenchmark) SetBaseline(baseline BenchmarkResult) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.Baseline = &baseline
}

// RunBenchmark 运行基准测试
func (pb *PerformanceBenchmark) RunBenchmark(ctx context.Context) (BenchmarkResult, error) {
	pb.mu.RLock()
	defer pb.mu.RUnlock()

	// 创建带超时的上下文
	benchCtx, cancel := context.WithTimeout(ctx, pb.Timeout)
	defer cancel()

	// 预热
	for i := 0; i < pb.Warmup; i++ {
		if _, err := pb.Run(benchCtx); err != nil {
			return BenchmarkResult{}, fmt.Errorf("warmup failed: %w", err)
		}
	}

	// 执行基准测试
	var durations []time.Duration
	var totalOperations int64
	var totalMemory MemoryUsage
	var totalCPU CPUUsage

	for i := 0; i < pb.Iterations; i++ {
		result, err := pb.Run(benchCtx)
		if err != nil {
			return BenchmarkResult{}, fmt.Errorf("benchmark iteration %d failed: %w", i, err)
		}

		durations = append(durations, result.Duration)
		totalOperations += result.Operations
		totalMemory.Allocated += result.MemoryUsage.Allocated
		totalMemory.Total += result.MemoryUsage.Total
		totalMemory.Heap += result.MemoryUsage.Heap
		totalMemory.Stack += result.MemoryUsage.Stack
		totalMemory.GC += result.MemoryUsage.GC
		totalCPU.UserTime += result.CPUUsage.UserTime
		totalCPU.SystemTime += result.CPUUsage.SystemTime
		totalCPU.IdleTime += result.CPUUsage.IdleTime
	}

	// 计算统计信息
	avgDuration := calculateAverageDuration(durations)
	minDuration := calculateMinDuration(durations)
	maxDuration := calculateMaxDuration(durations)
	stdDev := calculateStdDev(durations, avgDuration)
	percentiles := calculatePercentiles(durations)

	// 计算吞吐量
	throughput := float64(totalOperations) / avgDuration.Seconds()

	// 计算平均资源使用
	avgMemory := MemoryUsage{
		Allocated: totalMemory.Allocated / int64(pb.Iterations),
		Total:     totalMemory.Total / int64(pb.Iterations),
		Heap:      totalMemory.Heap / int64(pb.Iterations),
		Stack:     totalMemory.Stack / int64(pb.Iterations),
		GC:        totalMemory.GC / int64(pb.Iterations),
	}

	avgCPU := CPUUsage{
		UserTime:   totalCPU.UserTime / time.Duration(pb.Iterations),
		SystemTime: totalCPU.SystemTime / time.Duration(pb.Iterations),
		IdleTime:   totalCPU.IdleTime / time.Duration(pb.Iterations),
		Usage:      totalCPU.Usage / float64(pb.Iterations),
	}

	return BenchmarkResult{
		Name:        pb.Name,
		Duration:    avgDuration,
		Operations:  totalOperations,
		Throughput:  throughput,
		MemoryUsage: avgMemory,
		CPUUsage:    avgCPU,
		Iterations:  pb.Iterations,
		MinDuration: minDuration,
		MaxDuration: maxDuration,
		AvgDuration: avgDuration,
		StdDev:      stdDev,
		Percentiles: percentiles,
		Timestamp:   time.Now(),
		Metadata:    pb.Metadata,
	}, nil
}

// NewPerformanceMonitor 创建性能监控器
func NewPerformanceMonitor(config *PerformanceConfig) *PerformanceMonitor {
	if config == nil {
		config = &PerformanceConfig{
			DefaultTimeout:        60 * time.Second,
			DefaultIterations:     100,
			DefaultWarmup:         10,
			RegressionThreshold:   0.1,
			OutputDir:             "./performance-results",
			ReportFormat:          "json",
			EnableProfiling:       false,
			ProfilingDir:          "./profiles",
		}
	}

	return &PerformanceMonitor{
		benchmarks: make(map[string]*PerformanceBenchmark),
		results:    make([]BenchmarkResult, 0),
		detector:   NewRegressionDetector(config.RegressionThreshold),
		config:     config,
	}
}

// RegisterBenchmark 注册基准测试
func (pm *PerformanceMonitor) RegisterBenchmark(benchmark *PerformanceBenchmark) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.benchmarks[benchmark.Name] = benchmark
}

// RunBenchmark 运行指定基准测试
func (pm *PerformanceMonitor) RunBenchmark(ctx context.Context, name string) (BenchmarkResult, error) {
	pm.mu.RLock()
	benchmark, exists := pm.benchmarks[name]
	pm.mu.RUnlock()

	if !exists {
		return BenchmarkResult{}, fmt.Errorf("benchmark '%s' not found", name)
	}

	result, err := benchmark.RunBenchmark(ctx)
	if err != nil {
		return result, err
	}

	// 保存结果
	pm.mu.Lock()
	pm.results = append(pm.results, result)
	pm.mu.Unlock()

	// 检测回归
	pm.detector.DetectRegression(name, result)

	return result, nil
}

// RunAllBenchmarks 运行所有基准测试
func (pm *PerformanceMonitor) RunAllBenchmarks(ctx context.Context) (map[string]BenchmarkResult, error) {
	pm.mu.RLock()
	benchmarkNames := make([]string, 0, len(pm.benchmarks))
	for name := range pm.benchmarks {
		benchmarkNames = append(benchmarkNames, name)
	}
	pm.mu.RUnlock()

	results := make(map[string]BenchmarkResult)

	for _, name := range benchmarkNames {
		result, err := pm.RunBenchmark(ctx, name)
		if err != nil {
			return results, fmt.Errorf("failed to run benchmark '%s': %w", name, err)
		}
		results[name] = result
	}

	return results, nil
}

// GetResults 获取所有结果
func (pm *PerformanceMonitor) GetResults() []BenchmarkResult {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	results := make([]BenchmarkResult, len(pm.results))
	copy(results, pm.results)
	return results
}

// GetRegressionAlerts 获取回归告警
func (pm *PerformanceMonitor) GetRegressionAlerts() []RegressionAlert {
	return pm.detector.GetAlerts()
}

// NewRegressionDetector 创建回归检测器
func NewRegressionDetector(threshold float64) *RegressionDetector {
	return &RegressionDetector{
		baseline:  make(map[string]BenchmarkResult),
		current:   make(map[string]BenchmarkResult),
		threshold: threshold,
		alerts:    make([]RegressionAlert, 0),
	}
}

// SetBaseline 设置基准结果
func (rd *RegressionDetector) SetBaseline(name string, baseline BenchmarkResult) {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	rd.baseline[name] = baseline
}

// DetectRegression 检测性能回归
func (rd *RegressionDetector) DetectRegression(name string, current BenchmarkResult) {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	rd.current[name] = current

	baseline, exists := rd.baseline[name]
	if !exists {
		return
	}

	// 计算性能变化
	baselineThroughput := baseline.Throughput
	currentThroughput := current.Throughput

	if baselineThroughput > 0 {
		degradation := (baselineThroughput - currentThroughput) / baselineThroughput

		if degradation > rd.threshold {
			// 检测到性能回归
			alert := RegressionAlert{
				BenchmarkName: name,
				Severity:      "high",
				Message:       fmt.Sprintf("Performance degradation detected: %.2f%%", degradation*100),
				Baseline:      baselineThroughput,
				Current:       currentThroughput,
				Degradation:   degradation,
				Timestamp:     time.Now(),
			}

			rd.alerts = append(rd.alerts, alert)
		}
	}
}

// GetAlerts 获取告警
func (rd *RegressionDetector) GetAlerts() []RegressionAlert {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	
	alerts := make([]RegressionAlert, len(rd.alerts))
	copy(alerts, rd.alerts)
	return alerts
}

// ClearAlerts 清除告警
func (rd *RegressionDetector) ClearAlerts() {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	rd.alerts = make([]RegressionAlert, 0)
}

// 辅助函数

func calculateAverageDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

func calculateMinDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	min := durations[0]
	for _, d := range durations[1:] {
		if d < min {
			min = d
		}
	}
	return min
}

func calculateMaxDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	max := durations[0]
	for _, d := range durations[1:] {
		if d > max {
			max = d
		}
	}
	return max
}

func calculateStdDev(durations []time.Duration, mean time.Duration) float64 {
	if len(durations) == 0 {
		return 0
	}

	var sum float64
	meanNs := float64(mean.Nanoseconds())

	for _, d := range durations {
		diff := float64(d.Nanoseconds()) - meanNs
		sum += diff * diff
	}

	variance := sum / float64(len(durations))
	return math.Sqrt(variance)
}

func calculatePercentiles(durations []time.Duration) map[int]time.Duration {
	if len(durations) == 0 {
		return make(map[int]time.Duration)
	}

	// 排序
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	
	// 简单的冒泡排序
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	percentiles := make(map[int]time.Duration)
	percentiles[50] = sorted[len(sorted)/2] // 中位数
	percentiles[90] = sorted[int(float64(len(sorted))*0.9)]
	percentiles[95] = sorted[int(float64(len(sorted))*0.95)]
	percentiles[99] = sorted[int(float64(len(sorted))*0.99)]

	return percentiles
}
