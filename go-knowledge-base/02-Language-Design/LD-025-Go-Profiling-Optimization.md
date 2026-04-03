# LD-025: Go 性能剖析与优化 (Go Profiling & Optimization)

> **维度**: Language Design
> **级别**: S (19+ KB)
> **标签**: #profiling #pprof #optimization #performance #gc #memory #cpu
> **权威来源**:
>
> - [pprof Package](https://github.com/google/pprof) - Google
> - [Go Diagnostics](https://go.dev/doc/diagnostics) - Go Authors
> - [Go Performance Book](https://github.com/dgryski/go-perfbook) - Damian Gryski

---

## 1. 性能分析工具链

### 1.1 工具概览

```
┌─────────────────────────────────────────────────────────────┐
│                   Go Profiling Tools                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  运行时内置                                                   │
│  ├── net/http/pprof  - HTTP 接口                             │
│  ├── runtime/pprof   - 程序化接口                            │
│  └── runtime/trace   - 执行追踪                              │
│                                                              │
│  分析类型                                                     │
│  ├── CPU Profile     - CPU 使用分析                          │
│  ├── Memory Profile  - 内存分配分析                          │
│  ├── Block Profile   - 阻塞分析                              │
│  ├── Mutex Profile   - 锁竞争分析                            │
│  ├── Goroutine       - Goroutine 分析                        │
│  └── Trace           - 执行时间线                            │
│                                                              │
│  可视化工具                                                   │
│  ├── go tool pprof   - 命令行交互                            │
│  ├── pprof web UI    - 浏览器可视化                          │
│  └── flamegraph      - 火焰图                                │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 启用 Profiling

```go
// 方式 1: HTTP 接口 (推荐用于服务)
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // ...
}

// 方式 2: 程序化控制
import "runtime/pprof"

func profileCPU() {
    f, _ := os.Create("cpu.prof")
    defer f.Close()

    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // 执行代码...
}

func profileMemory() {
    f, _ := os.Create("mem.prof")
    defer f.Close()

    runtime.GC() // 获取准确数据
    if err := pprof.WriteHeapProfile(f); err != nil {
        log.Fatal(err)
    }
}
```

---

## 2. CPU Profiling

### 2.1 采样原理

```go
// CPU profiler 通过 SIGPROF 信号采样
// 默认 100Hz (每 10ms 采样一次)
// 开销约 5%

// 设置采样率
runtime.SetCPUProfileRate(100) // Hz
```

### 2.2 分析 CPU Profile

```bash
# 收集 CPU profile
go test -cpuprofile=cpu.prof -bench=.

# 或运行时收集
curl http://localhost:6060/debug/pprof/cpu?seconds=30 > cpu.prof

# 交互式分析
go tool pprof cpu.prof

# 常用命令
(pprof) top           # 显示消耗最高的函数
(pprof) top 20        # 显示前 20 个
(pprof) list FunctionName  # 显示函数源码级分析
(pprof) web           # 生成 SVG 火焰图
(pprof) png           # 生成 PNG 图像
```

### 2.3 火焰图解读

```
// 火焰图结构（从下往上）
main.main
  └── service.Process
        └── database.Query
              └── conn.QueryContext
                    └── driver.Query
                      └── network.Read

// 宽度 = 采样次数 = 耗时比例
// 颜色 = 函数所属包（无关性能）
```

---

## 3. 内存 Profiling

### 3.1 内存分配追踪

```go
// 内存 profile 类型

// 1. 已分配内存 (inuse_space)
//    当前存活的对象

// 2. 累计分配内存 (alloc_space)
//    程序启动以来的总分配

// 3. 对象数量 (inuse_objects/alloc_objects)
//    分配的对象个数
```

### 3.2 分析内存使用

```bash
# 收集内存 profile
go test -memprofile=mem.prof -bench=.
curl http://localhost:6060/debug/pprof/heap > heap.prof

# 分析
go tool pprof heap.prof

# 显示当前分配的内存
(pprof) top
(pprof) list FunctionName

# 显示累计分配（包含已释放）
go tool pprof -alloc_space heap.prof
```

### 3.3 内存泄漏检测

```go
// 1. 使用 runtime.ReadMemStats
func logMemoryStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
    log.Printf("TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
    log.Printf("Sys = %v MiB", bToMb(m.Sys))
    log.Printf("NumGC = %v", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

// 2. 比较两个时间点的 heap profile
func detectLeak() {
    f1, _ := os.Create("heap1.prof")
    pprof.WriteHeapProfile(f1)
    f1.Close()

    // 运行一段时间...
    time.Sleep(5 * time.Minute)

    f2, _ := os.Create("heap2.prof")
    pprof.WriteHeapProfile(f2)
    f2.Close()

    // 分析: go tool pprof -base heap1.prof heap2.prof
}
```

---

## 4. Goroutine 分析

### 4.1 Goroutine 泄露检测

```go
// 获取 goroutine profile
curl http://localhost:6060/debug/pprof/goroutine?debug=1

// debug=1: 文本格式，显示每个 goroutine 的栈
// debug=2: 包含更多细节

// 查找泄露模式
// 1. 持续增长的数量
// 2. 大量阻塞在 chan 操作的 goroutine
// 3. 大量等待 sync.Mutex 的 goroutine
```

### 4.2 Goroutine 分析示例

```bash
# 获取 goroutine 信息
curl http://localhost:6060/debug/pprof/goroutine?debug=1

# 输出示例:
# goroutine profile: total 100
# 50 @ 0x1032a20 0x1044e5b 0x1044e29 ...
# #   0x1044e28   main.processJob+0x48    /app/main.go:50
# #   0x1045a32   main.worker+0x32        /app/main.go:30
# #
# 30 @ 0x1032a20 0x1044e5b ...
# #   0x1045b12   database.query+0x112    /app/db.go:25

# 分析
go tool pprof goroutine.prof
(pprof) top
(pprof) traces
```

---

## 5. 阻塞和 Mutex Profiling

### 5.1 阻塞分析

```go
// 启用阻塞分析 (默认关闭，开销较高)
import "runtime"

func init() {
    runtime.SetBlockProfileRate(1) // 采样率：每 1 纳秒阻塞就记录
}
```

```bash
# 收集阻塞 profile
curl http://localhost:6060/debug/pprof/block > block.prof

# 分析
go tool pprof block.prof
(pprof) top
# 显示阻塞时间最长的操作
```

### 5.2 Mutex 竞争分析

```go
// 启用 mutex 分析 (Go 1.8+)
import "runtime"

func init() {
    runtime.SetMutexProfileFraction(1) // 1/1 采样
}
```

```bash
# 收集 mutex profile
curl http://localhost:6060/debug/pprof/mutex > mutex.prof

# 分析
go tool pprof mutex.prof
```

---

## 6. Execution Tracing

### 6.1 Trace 收集

```go
import "runtime/trace"

func main() {
    f, _ := os.Create("trace.out")
    defer f.Close()

    trace.Start(f)
    defer trace.Stop()

    // 执行代码...
}
```

```bash
# 或运行时收集
curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out

# 分析
go tool trace trace.out
# 会在浏览器中打开可视化界面
```

### 6.2 Trace 信息

```
Trace 视图显示：
- Goroutine 调度时间线
- 网络阻塞
- 系统调用
- GC 事件
- Heap 变化
```

---

## 7. 性能优化技巧

### 7.1 内存优化

```go
// 1. 预分配切片容量
func inefficient() []int {
    var result []int  // 容量为 0
    for i := 0; i < 1000; i++ {
        result = append(result, i) // 多次扩容
    }
    return result
}

func efficient() []int {
    result := make([]int, 0, 1000) // 预分配
    for i := 0; i < 1000; i++ {
        result = append(result, i)
    }
    return result
}

// 2. 复用对象（sync.Pool）
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func process() {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    // 使用 buf...
}

// 3. 避免不必要的字符串转换
func inefficientParse(data []byte) {
    s := string(data) // 分配
    for i := 0; i < len(s); i++ {
        // 处理 s[i]
    }
}

func efficientParse(data []byte) {
    for i := 0; i < len(data); i++ {
        // 直接处理 data[i]
    }
}
```

### 7.2 CPU 优化

```go
// 1. 减少锁竞争
type Counter struct {
    mu    sync.Mutex
    count int64
}

// 优化：使用原子操作
type AtomicCounter struct {
    count int64
}

func (c *AtomicCounter) Incr() {
    atomic.AddInt64(&c.count, 1)
}

// 2. 避免反射
// 慢
func slow(data interface{}) {
    v := reflect.ValueOf(data)
    // ...
}

// 快：类型断言
func fast(data interface{}) {
    if d, ok := data.(*MyType); ok {
        // ...
    }
}

// 3. 内联小函数
// go:noinline 可以禁用内联
// 编译器自动内联小函数
```

### 7.3 GC 优化

```go
// 1. 减少堆分配
func stackAlloc() {
    var arr [1024]byte // 栈分配
    _ = arr
}

func heapAlloc() []byte {
    arr := make([]byte, 1024) // 堆分配
    return arr
}

// 2. 调整 GC 目标
// GOGC=100 默认值（堆增长 100% 触发 GC）
// GOGC=200 降低 GC 频率，但增加内存使用
// GOGC=off 关闭 GC（仅手动触发）

// 3. 手动触发 GC
import "runtime"

runtime.GC()          // 立即触发
runtime.SetGCPercent(-1) // 禁用自动 GC
```

---

## 8. 基准测试优化

### 8.1 编写好的基准测试

```go
func BenchmarkSort(b *testing.B) {
    // 准备数据（不计入计时）
    data := generateLargeData()

    b.ResetTimer() // 重置计时器

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        // 重置状态
        copyData := make([]int, len(data))
        copy(copyData, data)
        b.StartTimer()

        sort.Ints(copyData)
    }
}

// 比较基准测试
func BenchmarkCompare(b *testing.B) {
    b.Run("naive", benchmarkNaive)
    b.Run("optimized", benchmarkOptimized)
}
```

### 8.2 使用 benchstat

```bash
# 运行多次获取稳定结果
go test -bench=. -count=5 > old.txt

# 优化后
go test -bench=. -count=5 > new.txt

# 比较
benchstat old.txt new.txt

# 输出：
# name        old time/op    new time/op    delta
# Benchmark   100ns ± 2%     50ns ± 1%    -50.00%  (p=0.008 n=5+5)
```

---

## 9. 视觉表征

### 9.1 Profiling 数据流

```
Running Application
       │
       ├──► CPU Samples ──┐
       │                   │
       ├──► Heap Alloc ────┼──► pprof Profile ──► Analysis
       │                   │      (protobuf)
       ├──► Goroutines ────┤
       │                   │
       ├──► Block Events ──┤
       │                   │
       └──► Trace Events ──┘         │
                                     ▼
                           ┌─────────────────┐
                           │  go tool pprof  │
                           │  pprof web UI   │
                           │  flamegraph     │
                           └─────────────────┘
```

### 9.2 优化流程

```
性能问题?
│
├── 高 CPU 使用?
│   ├── CPU profile → 找出热点函数
│   ├── 减少算法复杂度
   ├── 避免反射
   └── 使用并发
│
├── 高内存使用?
│   ├── Heap profile → 找出分配热点
│   ├── 预分配容量
│   ├── 使用 sync.Pool
   └── 减少逃逸到堆
│
├── GC 频繁?
│   ├── 减少堆分配
│   ├── 复用对象
│   └── 调整 GOGC
│
└── 高延迟?
    ├── Trace → 找出阻塞点
    ├── Block/Mutex profile
    └── 优化锁粒度
```

### 9.3 性能优化检查清单

```
□ 使用 -race 检查竞态条件
□ 使用 -benchmem 分析内存分配
□ CPU profile 确认热点函数
□ Heap profile 检查内存泄漏
□ Trace 分析调度延迟
□ 比较优化前后的 benchmark
□ 生产环境使用 pprof 持续监控
```

---

## 10. 完整代码示例

### 10.1 集成 Profiling 的服务

```go
package main

import (
    "expvar"
    "log"
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "time"
)

func main() {
    // 自定义 metrics
    expvar.Publish("goroutines", expvar.Func(func() interface{} {
        return runtime.NumGoroutine()
    }))

    expvar.Publish("memory", expvar.Func(func() interface{} {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        return map[string]interface{}{
            "alloc":      m.Alloc,
            "total_alloc": m.TotalAlloc,
            "sys":        m.Sys,
            "num_gc":     m.NumGC,
        }
    }))

    // 业务路由
    http.HandleFunc("/", handler)

    // 启动服务器
    log.Println("Server started on :8080")
    log.Println("pprof available at /debug/pprof/")
    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 模拟工作
    data := make([]byte, 1024*1024) // 1MB
    for i := range data {
        data[i] = byte(i)
    }

    w.Write([]byte("Done"))
}
```

### 10.2 性能测试脚本

```bash
#!/bin/bash
# profile.sh - 性能测试脚本

APP_URL="http://localhost:8080"
PPROF_URL="http://localhost:8080/debug/pprof"

echo "=== 收集性能数据 ==="

# 1. 预热
echo "Warming up..."
for i in {1..100}; do
    curl -s "$APP_URL/" > /dev/null
done

# 2. CPU Profile
echo "Collecting CPU profile..."
curl -s "$PPROF_URL/cpu?seconds=30" > cpu.prof
go tool pprof -pdf cpu.prof > cpu.pdf
echo "CPU profile saved to cpu.pdf"

# 3. Heap Profile
echo "Collecting heap profile..."
curl -s "$PPROF_URL/heap" > heap.prof
go tool pprof -pdf heap.prof > heap.pdf
echo "Heap profile saved to heap.pdf"

# 4. Goroutine Profile
echo "Collecting goroutine profile..."
curl -s "$PPROF_URL/goroutine" > goroutine.prof
echo "Goroutine profile saved"

# 5. Trace
echo "Collecting trace..."
curl -s "$PPROF_URL/trace?seconds=5" > trace.out
echo "Trace saved to trace.out"

echo "=== 分析完成 ==="
echo "查看结果:"
echo "  CPU: go tool pprof cpu.prof"
echo "  Heap: go tool pprof heap.prof"
echo "  Trace: go tool trace trace.out"
```

---

**质量评级**: S (19KB)
**完成日期**: 2026-04-02
