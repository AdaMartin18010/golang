# 基准测试 (Benchmarking)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #benchmarking #performance #testing #optimization
> **权威来源**:
>
> - [Package testing](https://pkg.go.dev/testing) - Go Official
> - [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) - Go Perf

---

## 1. 形式化定义

### 1.1 基准测试模型

**定义 1.1 (基准测试)**
$$\text{Benchmark} = \langle f, n, t_{total}, m_{alloc} \rangle$$

其中：

- $f$: 被测函数
- $n$: 迭代次数
- $t_{total}$: 总执行时间
- $m_{alloc}$: 内存分配量

**定义 1.2 (性能指标)**
$$\text{Throughput} = \frac{n}{t_{total}} \quad \text{(ops/sec)}$$

$$\text{Latency} = \frac{t_{total}}{n} \quad \text{(ns/op)}$$

**定义 1.3 (统计置信度)**
$$\text{ConfidenceInterval} = \bar{x} \pm z \cdot \frac{\sigma}{\sqrt{n}}$$

### 1.2 性能回归检测

**定理 1.1 (性能回归判定)**
$$\text{Regression} = \frac{\text{Latency}_{new} - \text{Latency}_{baseline}}{\text{Latency}_{baseline}} > \theta$$

其中 $\theta$ 是回归阈值（通常 5-10%）

---

## 2. Go 基准测试详解

### 2.1 基础用法

```go
// 基本基准测试
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

// 子基准测试
func BenchmarkFibonacciSizes(b *testing.B) {
    sizes := []int{10, 20, 30, 40}
    for _, n := range sizes {
        b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Fibonacci(n)
            }
        })
    }
}

// 并行基准测试
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ProcessItem()
        }
    })
}
```

### 2.2 内存分析

```go
// 内存分配分析
func BenchmarkWithAllocs(b *testing.B) {
    b.ReportAllocs()  // 报告内存分配

    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}

// 自定义内存统计
func BenchmarkCustomMetrics(b *testing.B) {
    var m1, m2 runtime.MemStats

    runtime.GC()
    runtime.ReadMemStats(&m1)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ProcessData()
    }
    b.StopTimer()

    runtime.ReadMemStats(&m2)

    b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "bytes/op")
    b.ReportMetric(float64(m2.NumGC-m1.NumGC), "gc/op")
}
```

### 2.3 高级技巧

```go
// 重置计时器
func BenchmarkWithSetup(b *testing.B) {
    // 准备阶段（不计时）
    data := setupLargeData()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Process(data)
    }

    b.StopTimer()
    cleanup(data)
}

// 控制并行度
func BenchmarkParallelControl(b *testing.B) {
    // GOMAXPROCS 控制并行度
    procs := []int{1, 2, 4, 8, 16}

    for _, p := range procs {
        b.Run(fmt.Sprintf("procs=%d", p), func(b *testing.B) {
            runtime.GOMAXPROCS(p)
            b.RunParallel(func(pb *testing.PB) {
                for pb.Next() {
                    Work()
                }
            })
        })
    }
}
```

---

## 3. 统计分析与比较

### 3.1 使用 benchstat

```bash
# 运行多次获取稳定结果
go test -bench=. -count=5 > old.txt

# 修改代码后再次运行
go test -bench=. -count=5 > new.txt

# 比较结果
benchstat old.txt new.txt
```

**输出示例**:

```
name        old time/op    new time/op    delta
Fibonacci   1.23µs ± 2%    1.15µs ± 3%   -6.50%  (p=0.008 n=5+5)

name        old alloc/op   new alloc/op   delta
Fibonacci     0.00B          0.00B          ~     (all equal)
```

### 3.2 置信区间分析

```go
package main

import (
    "fmt"
    "math"
    "sort"
)

// BenchmarkResult 基准测试结果
type BenchmarkResult struct {
    Name      string
    Samples   []float64  // 每次运行的结果 (ns/op)
    Iterations int
}

// Statistics 统计信息
type Statistics struct {
    Mean      float64
    Median    float64
    StdDev    float64
    Min       float64
    Max       float64
    P95       float64
    P99       float64
    CI95Lower float64
    CI95Upper float64
}

// CalculateStats 计算统计信息
func (br *BenchmarkResult) CalculateStats() Statistics {
    if len(br.Samples) == 0 {
        return Statistics{}
    }

    // 排序用于计算百分位数
    sorted := make([]float64, len(br.Samples))
    copy(sorted, br.Samples)
    sort.Float64s(sorted)

    // 均值
    var sum float64
    for _, v := range sorted {
        sum += v
    }
    mean := sum / float64(len(sorted))

    // 中位数
    median := calculateMedian(sorted)

    // 标准差
    variance := 0.0
    for _, v := range sorted {
        variance += math.Pow(v-mean, 2)
    }
    stdDev := math.Sqrt(variance / float64(len(sorted)))

    // 百分位数
    p95 := calculatePercentile(sorted, 0.95)
    p99 := calculatePercentile(sorted, 0.99)

    // 95% 置信区间 (使用 t-分布近似)
    n := float64(len(sorted))
    margin := 1.96 * stdDev / math.Sqrt(n)

    return Statistics{
        Mean:      mean,
        Median:    median,
        StdDev:    stdDev,
        Min:       sorted[0],
        Max:       sorted[len(sorted)-1],
        P95:       p95,
        P99:       p99,
        CI95Lower: mean - margin,
        CI95Upper: mean + margin,
    }
}

func calculateMedian(sorted []float64) float64 {
    n := len(sorted)
    if n%2 == 0 {
        return (sorted[n/2-1] + sorted[n/2]) / 2
    }
    return sorted[n/2]
}

func calculatePercentile(sorted []float64, p float64) float64 {
    index := p * float64(len(sorted)-1)
    lower := int(math.Floor(index))
    upper := int(math.Ceil(index))

    if lower == upper {
        return sorted[lower]
    }

    weight := index - float64(lower)
    return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// Compare 比较两个基准测试结果
func (br *BenchmarkResult) Compare(other *BenchmarkResult) ComparisonResult {
    stats1 := br.CalculateStats()
    stats2 := other.CalculateStats()

    delta := ((stats2.Mean - stats1.Mean) / stats1.Mean) * 100

    // 检查置信区间是否重叠
    significant := stats1.CI95Upper < stats2.CI95Lower ||
                   stats2.CI95Upper < stats1.CI95Lower

    return ComparisonResult{
        Baseline:    stats1,
        Current:     stats2,
        DeltaPercent: delta,
        Significant: significant,
        Regression:  delta > 5.0, // 5% 阈值
    }
}

type ComparisonResult struct {
    Baseline     Statistics
    Current      Statistics
    DeltaPercent float64
    Significant  bool
    Regression   bool
}

func (cr *ComparisonResult) String() string {
    direction := "improvement"
    if cr.DeltaPercent > 0 {
        direction = "regression"
    }

    sig := ""
    if cr.Significant {
        sig = " (significant)"
    }

    return fmt.Sprintf("%.2f%% %s%s", cr.DeltaPercent, direction, sig)
}
```

---

## 4. 基准测试模式

### 4.1 表格驱动基准测试

```go
func BenchmarkVariousInputs(b *testing.B) {
    benchmarks := []struct {
        name  string
        input int
    }{
        {"small", 10},
        {"medium", 100},
        {"large", 1000},
        {"xlarge", 10000},
    }

    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Process(bm.input)
            }
        })
    }
}
```

### 4.2 Map 性能基准

```go
func BenchmarkMapOperations(b *testing.B) {
    sizes := []int{100, 1000, 10000, 100000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            m := make(map[int]string, size)
            for i := 0; i < size; i++ {
                m[i] = fmt.Sprintf("value-%d", i)
            }

            b.Run("read", func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    _ = m[i%size]
                }
            })

            b.Run("write", func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    m[i] = fmt.Sprintf("new-%d", i)
                }
            })
        })
    }
}
```

### 4.3 并发基准

```go
func BenchmarkConcurrencyScenarios(b *testing.B) {
    scenarios := []struct {
        name       string
        goroutines int
        bufferSize int
    }{
        {"low-concurrency", 10, 100},
        {"medium-concurrency", 100, 1000},
        {"high-concurrency", 1000, 10000},
    }

    for _, s := range scenarios {
        b.Run(s.name, func(b *testing.B) {
            ch := make(chan int, s.bufferSize)
            var wg sync.WaitGroup

            // 启动消费者
            for i := 0; i < s.goroutines; i++ {
                wg.Add(1)
                go func() {
                    defer wg.Done()
                    for range ch {
                        // 处理数据
                    }
                }()
            }

            b.ResetTimer()

            // 生产者
            for i := 0; i < b.N; i++ {
                ch <- i
            }
            close(ch)
            wg.Wait()
        })
    }
}
```

---

## 5. 性能优化验证

### 5.1 优化前后对比框架

```go
// Optimizer 性能优化验证器
type Optimizer struct {
    baseline  BenchmarkResult
    optimized BenchmarkResult
}

// ValidateOptimization 验证优化效果
func (o *Optimizer) ValidateOptimization(minImprovement float64) (bool, error) {
    comparison := o.baseline.Compare(&o.optimized)

    if !comparison.Significant {
        return false, fmt.Errorf("result not statistically significant")
    }

    if comparison.DeltaPercent > -minImprovement {
        return false, fmt.Errorf("improvement %.2f%% below threshold %.2f%%",
            -comparison.DeltaPercent, minImprovement)
    }

    return true, nil
}

// ProfileGuidedOptimization 基于性能分析的优化
func ProfileGuidedOptimization(b *testing.B, fn func()) {
    // CPU Profiling
    cpuFile, _ := os.Create("cpu.prof")
    defer cpuFile.Close()
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()

    // Memory Profiling
    var m1, m2 runtime.MemStats
    runtime.ReadMemStats(&m1)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fn()
    }
    b.StopTimer()

    runtime.ReadMemStats(&m2)

    // 记录内存分析
    memFile, _ := os.Create("mem.prof")
    defer memFile.Close()
    pprof.WriteHeapProfile(memFile)

    b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "bytes/op")
}
```

---

## 6. 可视化与报告

### 6.1 性能趋势图

```mermaid
xychart-beta
    title "性能优化趋势"
    x-axis ["baseline", "opt-1", "opt-2", "opt-3", "final"]
    y-axis "ns/op" 0 --> 1000
    line [850, 720, 540, 380, 290]
```

### 6.2 性能对比矩阵

| 场景 | Baseline | Optimized | 提升 | 显著性 |
|------|----------|-----------|------|--------|
| 小数据 | 120ns | 95ns | 20.8% | ✅ |
| 中数据 | 450ns | 320ns | 28.9% | ✅ |
| 大数据 | 1200ns | 890ns | 25.8% | ✅ |
| 并发 | 850ns | 560ns | 34.1% | ✅ |

---

## 7. 最佳实践

### 7.1 Do's and Don'ts

| ✅ Do | ❌ Don't |
|-------|----------|
| 使用 `-count=5` 或更多 | 只运行一次就下结论 |
| 检查内存分配 | 只关注 CPU 时间 |
| 使用 `b.ResetTimer()` | 包含初始化时间 |
| 报告统计显著性 | 忽略置信区间 |
| 在相同硬件上比较 | 在不同机器上比较 |

### 7.2 基准测试清单

```markdown
## 基准测试审查清单

- [ ] 测试函数以 `Benchmark` 开头
- [ ] 使用 `b.N` 控制迭代
- [ ] 重置计时器排除初始化
- [ ] 报告内存分配 (`b.ReportAllocs()`)
- [ ] 运行多次获取稳定结果
- [ ] 使用 benchstat 比较结果
- [ ] 检查统计显著性
```

---

## 8. 工具链集成

### 8.1 CI/CD 集成

```yaml
# .github/workflows/benchmark.yml
name: Benchmark
on: [push, pull_request]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run Benchmarks
        run: go test -bench=. -benchmem -count=5 > benchmark.txt

      - name: Compare with main
        run: |
          git checkout main
          go test -bench=. -benchmem -count=5 > main.txt
          benchstat main.txt benchmark.txt
```

---

**质量评级**: S (15+ KB, 完整统计方法, Go 实现, 可视化)

**相关文档**:

- [性能优化](./02-Optimization.md)
- [内存分析](../05-Memory-Leak-Detection.md)
- [逃逸分析](../07-Escape-Analysis.md)
