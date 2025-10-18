# 并发编程示例

## 说明

此目录包含Go并发编程的各种模式和最佳实践示例。

## 示例列表

### 1. Pipeline - 管道模式

**目录**: `pipeline_example/`

**功能**:
- 生成器模式
- 管道阶段组合
- 扇入扇出（Fan-In/Fan-Out）
- 带缓冲的管道
- 错误处理

**运行**:
```bash
cd pipeline_example
go run main.go
```

---

### 2. Worker Pool - 工作池模式

**目录**: `worker_pool_example/`

**功能**:
- Worker池实现
- 任务队列管理
- 结果收集
- 优雅关闭

**运行**:
```bash
cd worker_pool_example
go run main.go
```

---

### 3. Concurrency Test - 并发测试

**文件**: `concurrency_test.go`

**功能**:
- 并发安全测试
- 竞态条件检测
- 性能基准测试

**运行**:
```bash
# 运行测试
go test -v

# 竞态检测
go test -race

# 基准测试
go test -bench=.
```

## 并发模式总结

### 1. Pipeline（管道）

**适用场景**:
- 数据处理流水线
- 多阶段转换
- 流式处理

**示例**:
```go
input := generator(10)
squared := square(input)
filtered := filterEven(squared)
output := printer(filtered)
```

---

### 2. Fan-Out/Fan-In（扇出/扇入）

**适用场景**:
- 并行处理
- 负载分散
- 结果聚合

**示例**:
```go
input := generator(100)
workers := fanOut(input, 5)  // 5个worker
merged := merge(workers...)  // 合并结果
```

---

### 3. Worker Pool（工作池）

**适用场景**:
- 限制并发数
- 资源池管理
- 任务调度

**示例**:
```go
pool := NewWorkerPool(4)
pool.Start()
pool.AddJob(job)
pool.Stop()
```

---

## 并发最佳实践

### 1. 避免竞态条件

```go
// ❌ 错误：无保护的共享状态
var count int
go func() { count++ }()

// ✅ 正确：使用互斥锁
var mu sync.Mutex
var count int
go func() {
    mu.Lock()
    count++
    mu.Unlock()
}()

// ✅ 更好：使用atomic
var count atomic.Int64
go func() { count.Add(1) }()
```

---

### 2. 避免Goroutine泄漏

```go
// ❌ 错误：无退出机制
go func() {
    for {
        doWork()
    }
}()

// ✅ 正确：使用context
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            doWork()
        }
    }
}
```

---

### 3. 正确使用WaitGroup

```go
// ❌ 错误：Add在goroutine内
var wg sync.WaitGroup
go func() {
    wg.Add(1)  // 可能太晚
    defer wg.Done()
}()

// ✅ 正确：Add在启动前
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
}()
```

---

### 4. Channel使用

```go
// 无缓冲channel：同步
ch := make(chan int)

// 缓冲channel：异步
ch := make(chan int, 100)

// 单向channel：类型安全
func send(ch chan<- int) { ch <- 1 }
func receive(ch <-chan int) { <-ch }
```

---

## 性能测试

### 运行竞态检测

```bash
go test -race ./...
```

### 运行基准测试

```bash
go test -bench=. -benchmem
```

### 生成性能剖析

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

---

## 调试工具

### 1. Go Trace

```go
import "runtime/trace"

f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()

// 你的代码
```

查看trace：
```bash
go tool trace trace.out
```

### 2. pprof

```go
import _ "net/http/pprof"

go func() {
    http.ListenAndServe("localhost:6060", nil)
}()
```

访问：http://localhost:6060/debug/pprof/

---

## 常见陷阱

### 1. 循环变量捕获

```go
// ❌ 错误
for _, v := range items {
    go func() {
        process(v)  // 所有goroutine看到最后一个v
    }()
}

// ✅ 正确
for _, v := range items {
    v := v  // 创建副本
    go func() {
        process(v)
    }()
}
```

### 2. Channel死锁

```go
// ❌ 错误
ch := make(chan int)
ch <- 42  // 阻塞，没有接收者

// ✅ 正确
ch := make(chan int, 1)  // 带缓冲
ch <- 42
```

### 3. 忘记关闭Channel

```go
// ❌ 错误
ch := make(chan int)
go producer(ch)
for v := range ch {  // 永远不会结束
    process(v)
}

// ✅ 正确
ch := make(chan int)
go func() {
    producer(ch)
    close(ch)  // 关闭channel
}()
for v := range ch {
    process(v)
}
```

---

## 更多资源

### 文档
- [Go并发编程](../../docs/03-并发编程/)
- [并发模式](../../docs/04-设计模式/)
- [性能优化](../../docs/05-性能优化/)

### 书籍推荐
- Go Concurrency Patterns（官方博客）
- Concurrency in Go（Katherine Cox-Buday）

---

**创建**: 2025年10月18日  
**Go版本**: 1.21+  
**状态**: 生产就绪
