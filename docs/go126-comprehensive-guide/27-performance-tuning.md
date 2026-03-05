# Go性能调优实战

> 基于实际场景的性能分析、优化和调优技术

---

## 一、性能分析基础

### 1.1 pprof工具链

```text
性能分析类型:
────────────────────────────────────────

CPU Profile:
├─ 采样CPU使用时间
├─ 默认100Hz采样
└─ 找出热点函数

Heap Profile:
├─ 内存分配分析
├─ inuse_space: 当前占用
└─ alloc_space: 累计分配

Goroutine Profile:
└─ goroutine堆栈分析

Block Profile:
└─ 阻塞分析

Mutex Profile:
└─ 锁竞争分析

Trace:
└─ 执行追踪

代码示例:
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // 应用代码...
}

// 分析命令
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof http://localhost:6060/debug/pprof/heap
go tool pprof http://localhost:6060/debug/pprof/goroutine
go tool trace trace.out
```

### 1.2 Benchmark编写

```text
基准测试:
────────────────────────────────────────

基本格式:
func BenchmarkXxx(b *testing.B) {
    // 准备工作

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 测试代码
    }
}

高级功能:
├─ b.Parallel(): 并行测试
├─ b.ReportAllocs(): 报告内存分配
├─ b.SetBytes(n): 设置处理字节数
└─ b.StopTimer()/StartTimer()

代码示例:
func BenchmarkProcess(b *testing.B) {
    data := generateData(1000)

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        process(data)
    }
}

func BenchmarkProcessParallel(b *testing.B) {
    data := generateData(1000)

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            process(data)
        }
    })
}

// 比较不同实现
func BenchmarkConcat(b *testing.B) {
    b.Run("Plus", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = "hello" + " " + "world"
        }
    })

    b.Run("Sprintf", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = fmt.Sprintf("%s %s", "hello", "world")
        }
    })

    b.Run("Builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var b strings.Builder
            b.WriteString("hello")
            b.WriteString(" ")
            b.WriteString("world")
            _ = b.String()
        }
    })
}
```

---

## 二、CPU优化

### 2.1 热点函数优化

```text
优化策略:
────────────────────────────────────────

1. 算法优化
   - 降低时间复杂度
   - 使用更合适的数据结构

2. 内联
   - 小函数内联
   - 减少调用开销

3. 避免重复计算
   - 缓存结果
   - 循环外提计算

代码示例:
// 不良: O(n²)
func findDuplicatesSlow(items []int) []int {
    var result []int
    for i := 0; i < len(items); i++ {
        for j := i + 1; j < len(items); j++ {
            if items[i] == items[j] {
                result = append(result, items[i])
            }
        }
    }
    return result
}

// 优化: O(n)
func findDuplicatesFast(items []int) []int {
    seen := make(map[int]bool)
    var result []int
    for _, item := range items {
        if seen[item] {
            result = append(result, item)
        } else {
            seen[item] = true
        }
    }
    return result
}

// 缓存计算结果
func fibSlow(n int) int {
    if n < 2 {
        return n
    }
    return fibSlow(n-1) + fibSlow(n-2)
}

func fibFast(n int) int {
    if n < 2 {
        return n
    }
    cache := make([]int, n+1)
    cache[1] = 1
    for i := 2; i <= n; i++ {
        cache[i] = cache[i-1] + cache[i-2]
    }
    return cache[n]
}
```

### 2.2 并行化

```text
并行优化:
────────────────────────────────────────

适用场景:
├─ 独立数据并行处理
├─ CPU密集型任务
└─ 大量计算

代码示例:
// 串行处理
func processSerial(items []Item) []Result {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = process(item)
    }
    return results
}

// 并行处理
func processParallel(items []Item) []Result {
    results := make([]Result, len(items))
    var wg sync.WaitGroup

    workers := runtime.NumCPU()
    chunkSize := (len(items) + workers - 1) / workers

    for w := 0; w < workers; w++ {
        start := w * chunkSize
        end := start + chunkSize
        if end > len(items) {
            end = len(items)
        }

        wg.Add(1)
        go func(s, e int) {
            defer wg.Done()
            for i := s; i < e; i++ {
                results[i] = process(items[i])
            }
        }(start, end)
    }

    wg.Wait()
    return results
}

// 使用Worker Pool
func processWithWorkerPool(items []Item, workers int) []Result {
    results := make([]Result, len(items))

    jobs := make(chan int, len(items))
    for i := range items {
        jobs <- i
    }
    close(jobs)

    var wg sync.WaitGroup
    for w := 0; w < workers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := range jobs {
                results[i] = process(items[i])
            }
        }()
    }

    wg.Wait()
    return results
}
```

---

## 三、内存优化

### 3.1 减少分配

```text
内存优化策略:
────────────────────────────────────────

1. 对象复用
2. 预分配
3. 避免不必要的拷贝
4. 使用值类型

代码示例:
// 对象池
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func processWithPool(data []byte) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf[:cap(buf)])

    copy(buf, data)
    // 处理...
}

// 预分配slice
func buildList(n int) []int {
    // 不良: 多次分配
    var result []int
    for i := 0; i < n; i++ {
        result = append(result, i)
    }

    // 优化: 预分配
    result = make([]int, 0, n)
    for i := 0; i < n; i++ {
        result = append(result, i)
    }

    return result
}

// 避免字符串拷贝
func processString(s string) {
    // 不良: 转换为[]byte会拷贝
    b := []byte(s)

    // 优化: 使用unsafe (小心使用)
    // b := *(*[]byte)(unsafe.Pointer(&s))

    _ = b
}
```

### 3.2 GC调优

```text
GC调优策略:
────────────────────────────────────────

1. 减少堆大小
2. 调整GOGC
3. 使用GOMEMLIMIT

代码示例:
// 监控GC
func monitorGC() {
    var m1, m2 runtime.MemStats

    runtime.ReadMemStats(&m1)
    // 执行操作...
    runtime.ReadMemStats(&m2)

    fmt.Printf("GC次数: %d -> %d\n", m1.NumGC, m2.NumGC)
    fmt.Printf("堆内存: %d MB -> %d MB\n",
        m1.HeapAlloc/1024/1024,
        m2.HeapAlloc/1024/1024)
}

// 调整GC参数
func tuneGC() {
    // 减少GC频率 (增加内存使用)
    debug.SetGCPercent(200)

    // 设置软内存限制
    debug.SetMemoryLimit(4 << 30) // 4GB
}
```

---

## 四、延迟优化

### 4.1 P99延迟优化

```text
延迟优化策略:
────────────────────────────────────────

1. 避免全局锁
2. 使用分片减少竞争
3. 异步处理
4. 超时控制

代码示例:
// 不良: 全局锁
var (
    globalMu sync.Mutex
    globalData map[string]int
)

func getDataBad(key string) int {
    globalMu.Lock()
    defer globalMu.Unlock()
    return globalData[key]
}

// 优化: 分片锁
type ShardedMap struct {
    shards [32]*shard
}

type shard struct {
    mu   sync.RWMutex
    data map[string]int
}

func (sm *ShardedMap) getShard(key string) *shard {
    hash := fnv32(key)
    return sm.shards[hash%32]
}

func (sm *ShardedMap) Get(key string) int {
    shard := sm.getShard(key)
    shard.mu.RLock()
    defer shard.mu.RUnlock()
    return shard.data[key]
}

// 异步处理
func processAsync(data []byte) chan Result {
    result := make(chan Result, 1)
    go func() {
        result <- process(data)
    }()
    return result
}
```

---

## 五、实战案例

### 5.1 HTTP服务优化

```text
HTTP服务性能优化:
────────────────────────────────────────

优化点:
├─ 连接池
├─ 超时设置
├─ 响应压缩
├─ 连接复用
└─ 优雅关闭

代码示例:
func optimizedServer() *http.Server {
    return &http.Server{
        Addr:         ":8080",
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,

        // 连接数限制
        ConnState: func(conn net.Conn, state http.ConnState) {
            switch state {
            case http.StateNew:
                atomic.AddInt32(&activeConnections, 1)
            case http.StateClosed:
                atomic.AddInt32(&activeConnections, -1)
            }
        },
    }
}

// 中间件优化
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 超时控制
        ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
        defer cancel()
        r = r.WithContext(ctx)

        // 响应压缩
        w.Header().Set("Content-Encoding", "gzip")

        next.ServeHTTP(w, r)
    })
}
```

---

*本章提供了Go性能调优的实战技术，涵盖CPU、内存、延迟等多维度优化。*
