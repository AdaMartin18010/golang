# 内存泄漏检测 (Memory Leak Detection)

> **分类**: 工程与云原生

---

## 常见内存泄漏

### 1. Goroutine 泄漏

```go
// ❌ 错误: 发送者阻塞
func bad() {
    ch := make(chan int)
    go func() {
        ch <- 42  // 无人接收，永久阻塞
    }()
}

// ✅ 正确: 使用缓冲或 select
func good() {
    ch := make(chan int, 1)  // 缓冲
    go func() {
        ch <- 42
    }()
}
```

### 2. Timer 未停止

```go
// ❌ 错误
timer := time.NewTimer(time.Hour)
// 如果提前返回，timer 继续运行

// ✅ 正确
timer := time.NewTimer(time.Hour)
defer timer.Stop()
```

### 3. 全局引用

```go
// ❌ 错误: 全局缓存无限增长
var cache = map[string][]byte{}

func store(key string, data []byte) {
    cache[key] = data  // 永不清理
}

// ✅ 正确: LRU 缓存
import lru "github.com/hashicorp/golang-lru"

cache, _ := lru.New(1000)  // 限制大小
```

---

## 检测方法

### pprof

```go
import _ "net/http/pprof"

go http.ListenAndServe("localhost:6060", nil)
```

```bash
# 查看堆分配
go tool pprof http://localhost:6060/debug/pprof/heap

# 对比两个时间点的堆
# T1
curl http://localhost:6060/debug/pprof/heap > heap.1
# T2
curl http://localhost:6060/debug/pprof/heap > heap.2

# 比较
go tool pprof -base heap.1 heap.2
```

### Goroutine 泄漏检测

```go
// 使用 goleak
import "go.uber.org/goleak"

func TestFunction(t *testing.T) {
    defer goleak.VerifyNone(t)

    // 测试代码
}
```

---

## 监控

```go
// 导出内存指标
var (
    heapAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "go_heap_alloc_bytes",
        Help: "Heap allocation",
    })
    goroutines = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "go_goroutines",
        Help: "Number of goroutines",
    })
)

func recordMetrics() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    heapAlloc.Set(float64(m.HeapAlloc))
    goroutines.Set(float64(runtime.NumGoroutine()))
}
```
