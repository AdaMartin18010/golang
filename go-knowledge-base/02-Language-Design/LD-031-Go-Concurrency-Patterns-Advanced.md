# LD-031-Go-Concurrency-Patterns-Advanced

> **Dimension**: 02-Language-Design
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: Go 1.26 Advanced Concurrency
> **Size**: >20KB

---

## 1. Go并发基础回顾

### 1.1 CSP模型

Go基于CSP(Communicating Sequential Processes)模型：

```
"不要通过共享内存来通信，而是通过通信来共享内存"

              Channel
  Goroutine ◄────────────► Goroutine
      │                        │
      └──────────┬─────────────┘
                 │
            同步/通信
```

### 1.2 Goroutine调度

```
Go Scheduler (GMP模型):

G - Goroutine: 用户态轻量线程 (~2KB初始栈)
M - Machine: OS线程
P - Processor: 逻辑处理器 (GOMAXPROCS)

┌─────────────────────────────────────────┐
│           Go Runtime Scheduler          │
├─────────────────────────────────────────┤
│                                         │
│   Global Queue                          │
│   [G1, G2, G3, ...]                     │
│        │                                │
│        ▼                                │
│   ┌─────────┐  ┌─────────┐  ┌─────────┐│
│   │    P    │  │    P    │  │    P    ││
│   │ ┌─────┐ │  │ ┌─────┐ │  │ ┌─────┐ ││
│   │ │Local│ │  │ │Local│ │  │ │Local│ ││
│   │ │Queue│ │  │ │Queue│ │  │ │Queue│ ││
│   │ └──┬──┘ │  │ └──┬──┘ │  │ └──┬──┘ ││
│   │    │    │  │    │    │  │    │    ││
│   └────┼────┘  └────┼────┘  └────┼────┘│
│        │            │            │     │
│   ┌────┴────────────┴────────────┴────┐│
│   │         M (OS Threads)            ││
│   │    ┌───┐  ┌───┐  ┌───┐  ┌───┐    ││
│   │    │M1 │  │M2 │  │M3 │  │M4 │    ││
│   │    └───┘  └───┘  └───┘  └───┘    ││
│   └───────────────────────────────────┘│
│                                         │
│   Work Stealing: 空闲P从其他P偷取G     │
└─────────────────────────────────────────┘
```

---

## 2. 高级Channel模式

### 2.1 Pipeline模式

```go
// 多阶段数据处理流水线
func PipelineDemo() {
    // Stage 1: 生成数据
    generator := func(nums ...int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for _, n := range nums {
                out <- n
            }
        }()
        return out
    }

    // Stage 2: 平方
    square := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                out <- n * n
            }
        }()
        return out
    }

    // Stage 3: 过滤偶数
    filter := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                if n%2 == 0 {
                    out <- n
                }
            }
        }()
        return out
    }

    // 连接管道
    nums := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squares := square(nums)
    evens := filter(squares)

    // 消费结果
    for n := range evens {
        fmt.Println(n)  // 4, 16, 36, 64, 100
    }
}
```

### 2.2 Fan-Out/Fan-In模式

```go
// 多个worker并行处理，结果汇总
func FanOutFanInDemo() {
    // 任务生成器
    jobs := make(chan int, 100)
    go func() {
        for i := 1; i <= 100; i++ {
            jobs <- i
        }
        close(jobs)
    }()

    // Fan-Out: 启动多个worker
    numWorkers := 5
    results := make([]<-chan int, numWorkers)

    for i := 0; i < numWorkers; i++ {
        results[i] = worker(i, jobs)
    }

    // Fan-In: 合并结果
    final := merge(results...)

    // 收集结果
    var sum int
    for r := range final {
        sum += r
    }
    fmt.Printf("Total sum: %d\n", sum)
}

func worker(id int, jobs <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for job := range jobs {
            // 模拟处理
            result := job * job
            fmt.Printf("Worker %d processed job %d\n", id, job)
            out <- result
        }
    }()
    return out
}

func merge(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            out <- n
        }
    }

    wg.Add(len(channels))
    for _, c := range channels {
        go output(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### 2.3 优雅关闭模式

```go
// 使用context实现优雅关闭
type Server struct {
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

func NewServer() *Server {
    ctx, cancel := context.WithCancel(context.Background())
    return &Server{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (s *Server) Start() {
    // 启动多个服务组件
    s.wg.Add(3)

    go s.httpServer()
    go s.workerPool()
    go s.metricsReporter()
}

func (s *Server) httpServer() {
    defer s.wg.Done()

    srv := &http.Server{Addr: ":8080"}

    go func() {
        <-s.ctx.Done()  // 等待关闭信号
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        srv.Shutdown(shutdownCtx)
    }()

    srv.ListenAndServe()
}

func (s *Server) workerPool() {
    defer s.wg.Done()

    jobs := make(chan Job, 100)

    // 启动worker
    for i := 0; i < 10; i++ {
        go s.worker(jobs)
    }

    // 接收任务直到context取消
    for {
        select {
        case job := <-s.jobSource():
            jobs <- job
        case <-s.ctx.Done():
            close(jobs)
            return
        }
    }
}

func (s *Server) worker(jobs <-chan Job) {
    for job := range jobs {
        s.processJob(job)
    }
}

func (s *Server) Stop() {
    s.cancel()   // 发送关闭信号
    s.wg.Wait()  // 等待所有组件完成
}
```

---

## 3. 同步原语高级用法

### 3.1 sync.Map优化场景

```go
// sync.Map适用于以下场景:
// 1. 只写一次，读很多次
// 2. 多个goroutine并发读写不同key

// 缓存实现
type Cache struct {
    data sync.Map
    ttl  time.Duration
}

type cacheItem struct {
    value      interface{}
    expiration time.Time
}

func (c *Cache) Get(key string) (interface{}, bool) {
    item, ok := c.data.Load(key)
    if !ok {
        return nil, false
    }

    ci := item.(cacheItem)
    if time.Now().After(ci.expiration) {
        c.data.Delete(key)
        return nil, false
    }

    return ci.value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.data.Store(key, cacheItem{
        value:      value,
        expiration: time.Now().Add(c.ttl),
    })
}

// 定期清理
func (c *Cache) StartCleanup() {
    ticker := time.NewTicker(c.ttl)
    go func() {
        for range ticker.C {
            c.cleanup()
        }
    }()
}

func (c *Cache) cleanup() {
    now := time.Now()
    c.data.Range(func(key, value interface{}) bool {
        item := value.(cacheItem)
        if now.After(item.expiration) {
            c.data.Delete(key)
        }
        return true
    })
}
```

### 3.2 sync.Pool高级用法

```go
// 对象池，减少GC压力
type BufferPool struct {
    pool sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, 1024)
            },
        },
    }
}

func (p *BufferPool) Get() []byte {
    buf := p.pool.Get().([]byte)
    return buf[:0]  // 重置但保留容量
}

func (p *BufferPool) Put(buf []byte) {
    if cap(buf) > 1024*1024 {
        // 不回收过大的buffer，避免内存泄漏
        return
    }
    p.pool.Put(buf)
}

// 使用示例
var pool = NewBufferPool()

func ProcessData(data []byte) []byte {
    buf := pool.Get()
    defer pool.Put(buf)

    // 处理数据...
    buf = append(buf, data...)

    result := make([]byte, len(buf))
    copy(result, buf)
    return result
}
```

### 3.3 条件变量模式

```go
// 任务队列，支持批量处理
type BatchQueue struct {
    mu       sync.Mutex
    cond     *sync.Cond
    items    []Item
    maxSize  int
    maxWait  time.Duration
}

func NewBatchQueue(maxSize int, maxWait time.Duration) *BatchQueue {
    bq := &BatchQueue{
        items:   make([]Item, 0, maxSize),
        maxSize: maxSize,
        maxWait: maxWait,
    }
    bq.cond = sync.NewCond(&bq.mu)
    return bq
}

func (bq *BatchQueue) Enqueue(item Item) {
    bq.mu.Lock()
    defer bq.mu.Unlock()

    bq.items = append(bq.items, item)

    if len(bq.items) >= bq.maxSize {
        bq.cond.Broadcast()  // 通知等待的batch processor
    }
}

func (bq *BatchQueue) GetBatch() []Item {
    bq.mu.Lock()
    defer bq.mu.Unlock()

    // 等待直到有数据或超时
    for len(bq.items) == 0 {
        bq.cond.Wait()
    }

    // 复制batch
    batch := make([]Item, len(bq.items))
    copy(batch, bq.items)
    bq.items = bq.items[:0]

    return batch
}

func (bq *BatchQueue) StartProcessor(process func([]Item)) {
    go func() {
        timer := time.NewTimer(bq.maxWait)
        defer timer.Stop()

        for {
            select {
            case <-timer.C:
                batch := bq.GetBatch()
                if len(batch) > 0 {
                    process(batch)
                }
                timer.Reset(bq.maxWait)
            }
        }
    }()
}
```

---

## 4. Context高级模式

### 4.1 可取消的并行任务

```go
// 并行执行多个任务，任一失败则全部取消
func ParallelWithCancel(tasks []Task) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    errChan := make(chan error, len(tasks))
    var wg sync.WaitGroup

    for _, task := range tasks {
        wg.Add(1)
        go func(t Task) {
            defer wg.Done()

            if err := t.Execute(ctx); err != nil {
                errChan <- err
                cancel()  // 取消其他任务
            }
        }(task)
    }

    // 等待所有任务完成
    go func() {
        wg.Wait()
        close(errChan)
    }()

    // 返回第一个错误
    for err := range errChan {
        return err
    }

    return nil
}
```

### 4.2 带超时的级联调用

```go
func CascadeCall(ctx context.Context) error {
    // 总体超时
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // 第一步: 获取token (2秒超时)
    tokenCtx, cancel1 := context.WithTimeout(ctx, 2*time.Second)
    token, err := getToken(tokenCtx)
    cancel1()
    if err != nil {
        return err
    }

    // 第二步: 获取数据 (5秒超时)
    dataCtx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
    data, err := fetchData(dataCtx, token)
    cancel2()
    if err != nil {
        return err
    }

    // 第三步: 处理数据 (剩余时间)
    result, err := processData(ctx, data)
    _ = result

    return err
}
```

### 4.3 Context值传递最佳实践

```go
// 使用私有类型作为key，避免冲突
type contextKey string

const (
    requestIDKey contextKey = "request_id"
    userIDKey    contextKey = "user_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}

// 中间件中使用
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }

        ctx := WithRequestID(r.Context(), requestID)
        w.Header().Set("X-Request-ID", requestID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## 5. 并发模式最佳实践

### 5.1 防goroutine泄漏

```go
// 错误示例 - 可能泄漏
func Leaky() {
    ch := make(chan int)
    go func() {
        ch <- doWork()  // 如果无人接收，永久阻塞
    }()
    // 如果这里提前返回，goroutine泄漏
}

// 正确做法 - 使用buffered channel或select
func NonLeaky() {
    ch := make(chan int, 1)  // buffered
    go func() {
        ch <- doWork()
    }()
}

// 或使用select + done channel
func NonLeakyWithDone(done <-chan struct{}) {
    ch := make(chan int)
    go func() {
        select {
        case ch <- doWork():
        case <-done:
            // 清理资源
            return
        }
    }()
}
```

### 5.2 优雅处理panic

```go
func SafeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Recovered from panic: %v\n%s", r, debug.Stack())
            }
        }()
        fn()
    }()
}

// 使用
SafeGo(func() {
    mightPanic()
})
```

### 5.3 性能优化建议

| 场景 | 建议 |
|------|------|
| 高频创建goroutine | 使用worker pool |
| 大量channel操作 | 考虑batch处理 |
| 共享数据读多写少 | 使用sync.Map |
| 临时对象频繁分配 | 使用sync.Pool |
| 复杂状态同步 | 考虑使用mutex而非channel |

---

## 6. 参考文献

1. "Concurrency in Go" - Katherine Cox-Buday
2. Go Memory Model
3. Go Scheduler Internals
4. "Advanced Go Concurrency Patterns" - Sameer Ajmani

---

*Last Updated: 2026-04-03*
