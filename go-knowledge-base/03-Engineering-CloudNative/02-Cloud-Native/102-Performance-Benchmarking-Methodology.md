# 性能基准测试方法论 (Performance Benchmarking Methodology)

> **分类**: 工程与云原生
> **标签**: #benchmarking #performance #testing #methodology
> **参考**: Google Benchmark, JMH, Performance Testing Best Practices

---

## 性能测试框架

```go
package benchmark

import (
    "context"
    "fmt"
    "math"
    "runtime"
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

// BenchmarkConfig 基准测试配置
type BenchmarkConfig struct {
    Name            string
    Duration        time.Duration
    WarmupDuration  time.Duration
    Concurrency     int
    RateLimit       int // 每秒请求数，0表示无限制

    // 统计配置
    LatencyPercentiles []float64 // 如 []float64{0.5, 0.95, 0.99}
}

// BenchmarkResult 基准测试结果
type BenchmarkResult struct {
    Config          BenchmarkConfig

    // 吞吐量
    TotalRequests   int64
    RequestsPerSec  float64

    // 延迟统计
    LatencyStats    LatencyStatistics

    // 错误统计
    ErrorCount      int64
    ErrorRate       float64

    // 资源使用
    MaxMemory       uint64
    AvgCPUUsage     float64

    // 自定义指标
    CustomMetrics   map[string]float64
}

// LatencyStatistics 延迟统计
type LatencyStatistics struct {
    Min       time.Duration
    Max       time.Duration
    Mean      time.Duration
    StdDev    time.Duration

    Percentiles map[float64]time.Duration
}

// LatencyRecorder 延迟记录器
type LatencyRecorder struct {
    values   []time.Duration
    capacity int
    mu       sync.RWMutex
}

// NewLatencyRecorder 创建延迟记录器
func NewLatencyRecorder(capacity int) *LatencyRecorder {
    return &LatencyRecorder{
        values:   make([]time.Duration, 0, capacity),
        capacity: capacity,
    }
}

// Record 记录延迟
func (lr *LatencyRecorder) Record(d time.Duration) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if len(lr.values) >= lr.capacity {
        // 随机采样丢弃
        idx := int(time.Now().UnixNano() % int64(len(lr.values)))
        lr.values[idx] = d
    } else {
        lr.values = append(lr.values, d)
    }
}

// Statistics 计算统计
func (lr *LatencyRecorder) Statistics(percentiles []float64) LatencyStatistics {
    lr.mu.RLock()
    values := make([]time.Duration, len(lr.values))
    copy(values, lr.values)
    lr.mu.RUnlock()

    if len(values) == 0 {
        return LatencyStatistics{}
    }

    // 排序
    sortDurations(values)

    stats := LatencyStatistics{
        Min:         values[0],
        Max:         values[len(values)-1],
        Percentiles: make(map[float64]time.Duration),
    }

    // 计算平均值
    var sum time.Duration
    for _, v := range values {
        sum += v
    }
    stats.Mean = sum / time.Duration(len(values))

    // 计算标准差
    var variance float64
    meanFloat := float64(stats.Mean)
    for _, v := range values {
        diff := float64(v) - meanFloat
        variance += diff * diff
    }
    variance /= float64(len(values))
    stats.StdDev = time.Duration(math.Sqrt(variance))

    // 计算百分位数
    for _, p := range percentiles {
        idx := int(float64(len(values)-1) * p)
        stats.Percentiles[p] = values[idx]
    }

    return stats
}

func sortDurations(durations []time.Duration) {
    // 简单的快速排序实现
    if len(durations) <= 1 {
        return
    }

    pivot := durations[len(durations)/2]
    left, right := 0, len(durations)-1

    for left <= right {
        for durations[left] < pivot {
            left++
        }
        for durations[right] > pivot {
            right--
        }
        if left <= right {
            durations[left], durations[right] = durations[right], durations[left]
            left++
            right--
        }
    }

    sortDurations(durations[:right+1])
    sortDurations(durations[left:])
}

// BenchmarkRunner 基准测试运行器
type BenchmarkRunner struct {
    config   BenchmarkConfig
    recorder *LatencyRecorder

    totalRequests int64
    errorCount    int64

    startTime    time.Time
    warmupEnd    time.Time
    testEnd      time.Time

    stopCh chan struct{}
}

// NewBenchmarkRunner 创建基准测试运行器
func NewBenchmarkRunner(config BenchmarkConfig) *BenchmarkRunner {
    return &BenchmarkRunner{
        config:   config,
        recorder: NewLatencyRecorder(100000),
        stopCh:   make(chan struct{}),
    }
}

// Run 执行基准测试
func (br *BenchmarkRunner) Run(operation func() error) (*BenchmarkResult, error) {
    br.startTime = time.Now()
    br.warmupEnd = br.startTime.Add(br.config.WarmupDuration)
    br.testEnd = br.warmupEnd.Add(br.config.Duration)

    // 启动工作线程
    var wg sync.WaitGroup
    rateLimiter := NewRateLimiter(br.config.RateLimit)

    for i := 0; i < br.config.Concurrency; i++ {
        wg.Add(1)
        go br.worker(&wg, operation, rateLimiter)
    }

    // 等待测试结束
    time.Sleep(br.config.WarmupDuration + br.config.Duration)
    close(br.stopCh)

    wg.Wait()

    // 收集结果
    result := &BenchmarkResult{
        Config:        br.config,
        TotalRequests: atomic.LoadInt64(&br.totalRequests),
        ErrorCount:    atomic.LoadInt64(&br.errorCount),
        LatencyStats:  br.recorder.Statistics(br.config.LatencyPercentiles),
        CustomMetrics: make(map[string]float64),
    }

    // 计算吞吐量
    elapsed := br.config.Duration.Seconds()
    if elapsed > 0 {
        result.RequestsPerSec = float64(result.TotalRequests) / elapsed
    }

    // 计算错误率
    if result.TotalRequests > 0 {
        result.ErrorRate = float64(result.ErrorCount) / float64(result.TotalRequests) * 100
    }

    // 获取内存统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    result.MaxMemory = m.Sys

    return result, nil
}

func (br *BenchmarkRunner) worker(wg *sync.WaitGroup, operation func() error, rateLimiter *RateLimiter) {
    defer wg.Done()

    for {
        select {
        case <-br.stopCh:
            return
        default:
        }

        // 速率限制
        if rateLimiter != nil {
            rateLimiter.Wait()
        }

        // 检查是否在预热阶段
        now := time.Now()
        isWarmup := now.Before(br.warmupEnd)

        // 执行操作
        start := time.Now()
        err := operation()
        latency := time.Since(start)

        if !isWarmup {
            atomic.AddInt64(&br.totalRequests, 1)

            if err != nil {
                atomic.AddInt64(&br.errorCount, 1)
            }

            br.recorder.Record(latency)
        }
    }
}

// RateLimiter 速率限制器
type RateLimiter struct {
    rate      int
    interval  time.Duration
    tokens    chan struct{}
    stopCh    chan struct{}
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate int) *RateLimiter {
    if rate <= 0 {
        return nil
    }

    rl := &RateLimiter{
        rate:     rate,
        interval: time.Second / time.Duration(rate),
        tokens:   make(chan struct{}, rate),
        stopCh:   make(chan struct{}),
    }

    go rl.refill()
    return rl
}

func (rl *RateLimiter) refill() {
    ticker := time.NewTicker(rl.interval)
    defer ticker.Stop()

    for {
        select {
        case <-rl.stopCh:
            return
        case <-ticker.C:
            select {
            case rl.tokens <- struct{}{}:
            default:
            }
        }
    }
}

// Wait 等待令牌
func (rl *RateLimiter) Wait() {
    <-rl.tokens
}

// Stop 停止
func (rl *RateLimiter) Stop() {
    close(rl.stopCh)
}

// ReportGenerator 报告生成器
type ReportGenerator struct {
    results []*BenchmarkResult
}

// NewReportGenerator 创建报告生成器
func NewReportGenerator() *ReportGenerator {
    return &ReportGenerator{
        results: make([]*BenchmarkResult, 0),
    }
}

// AddResult 添加结果
func (rg *ReportGenerator) AddResult(result *BenchmarkResult) {
    rg.results = append(rg.results, result)
}

// GenerateMarkdown 生成Markdown报告
func (rg *ReportGenerator) GenerateMarkdown() string {
    report := "# Benchmark Report\n\n"
    report += "## Summary\n\n"
    report += "| Test | RPS | P50 | P95 | P99 | Error Rate |\n"
    report += "|------|-----|-----|-----|-----|------------|\n"

    for _, r := range rg.results {
        report += fmt.Sprintf("| %s | %.2f | %s | %s | %s | %.2f%% |\n",
            r.Config.Name,
            r.RequestsPerSec,
            r.LatencyStats.Percentiles[0.50],
            r.LatencyStats.Percentiles[0.95],
            r.LatencyStats.Percentiles[0.99],
            r.ErrorRate,
        )
    }

    return report
}

// CompareResults 比较结果
func (rg *ReportGenerator) CompareResults(baseline, current *BenchmarkResult) *ComparisonResult {
    comp := &ComparisonResult{}

    // 吞吐量变化
    comp.RPSChange = (current.RequestsPerSec - baseline.RequestsPerSec) / baseline.RequestsPerSec * 100

    // 延迟变化
    for p, baselineLatency := range baseline.LatencyStats.Percentiles {
        if currentLatency, ok := current.LatencyStats.Percentiles[p]; ok {
            change := float64(currentLatency-baselineLatency) / float64(baselineLatency) * 100
            comp.LatencyChanges[p] = change
        }
    }

    // 判定是否回归
    comp.IsRegression = comp.RPSChange < -10 // RPS下降超过10%

    return comp
}

// ComparisonResult 比较结果
type ComparisonResult struct {
    RPSChange      float64
    LatencyChanges map[float64]float64
    IsRegression   bool
}
```

---

## 使用示例

```go
package main

import (
    "fmt"
    "math/rand"
    "time"

    "benchmark"
)

func main() {
    // 创建基准测试配置
    config := benchmark.BenchmarkConfig{
        Name:               "TaskEnqueue",
        Duration:           30 * time.Second,
        WarmupDuration:     5 * time.Second,
        Concurrency:        100,
        RateLimit:          0, // 无限制
        LatencyPercentiles: []float64{0.50, 0.95, 0.99},
    }

    // 创建运行器
    runner := benchmark.NewBenchmarkRunner(config)

    // 定义测试操作
    operation := func() error {
        // 模拟任务入队操作
        time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)

        // 偶尔模拟错误
        if rand.Float32() < 0.001 {
            return fmt.Errorf("random error")
        }

        return nil
    }

    // 运行测试
    result, err := runner.Run(operation)
    if err != nil {
        panic(err)
    }

    // 打印结果
    fmt.Printf("Benchmark: %s\n", result.Config.Name)
    fmt.Printf("Total Requests: %d\n", result.TotalRequests)
    fmt.Printf("Requests/sec: %.2f\n", result.RequestsPerSec)
    fmt.Printf("Error Rate: %.4f%%\n", result.ErrorRate)
    fmt.Printf("Latency P50: %s\n", result.LatencyStats.Percentiles[0.50])
    fmt.Printf("Latency P95: %s\n", result.LatencyStats.Percentiles[0.95])
    fmt.Printf("Latency P99: %s\n", result.LatencyStats.Percentiles[0.99])
    fmt.Printf("Max Memory: %d MB\n", result.MaxMemory/1024/1024)
}
```
