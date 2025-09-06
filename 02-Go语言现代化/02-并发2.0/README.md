# Go并发编程2.0 - 现代化并发模式

<!-- TOC START -->
- [Go并发编程2.0 - 现代化并发模式](#go并发编程20---现代化并发模式)
  - [1.1 概述](#11-概述)
  - [1.2 核心特性](#12-核心特性)
  - [1.3 现代化并发模式](#13-现代化并发模式)
  - [1.4 性能优化](#14-性能优化)
  - [1.5 最佳实践](#15-最佳实践)
  - [1.6 实际应用](#16-实际应用)
<!-- TOC END -->

## 1.1 概述

Go并发编程2.0代表了Go语言并发编程的现代化演进，结合了最新的语言特性和最佳实践，为开发者提供了更高效、更安全的并发编程体验。

## 1.2 核心特性

### 1.2.1 泛型并发模式

```go
// 泛型Worker Pool
type WorkerPool[T any] struct {
    workers    int
    jobQueue   chan Job[T]
    resultChan chan Result[T]
    wg         sync.WaitGroup
}

type Job[T any] struct {
    ID   string
    Data T
    Process func(T) (T, error)
}

type Result[T any] struct {
    JobID string
    Data  T
    Error error
}

func NewWorkerPool[T any](workers int) *WorkerPool[T] {
    return &WorkerPool[T]{
        workers:    workers,
        jobQueue:   make(chan Job[T], 100),
        resultChan: make(chan Result[T], 100),
    }
}

func (wp *WorkerPool[T]) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx)
    }
}

func (wp *WorkerPool[T]) worker(ctx context.Context) {
    defer wp.wg.Done()
    
    for {
        select {
        case job := <-wp.jobQueue:
            result, err := job.Process(job.Data)
            wp.resultChan <- Result[T]{
                JobID: job.ID,
                Data:  result,
                Error: err,
            }
        case <-ctx.Done():
            return
        }
    }
}
```

### 1.2.2 结构化并发

```go
// 结构化并发管理器
type StructuredConcurrency struct {
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

func NewStructuredConcurrency() *StructuredConcurrency {
    ctx, cancel := context.WithCancel(context.Background())
    return &StructuredConcurrency{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (sc *StructuredConcurrency) Go(fn func(context.Context) error) {
    sc.wg.Add(1)
    go func() {
        defer sc.wg.Done()
        if err := fn(sc.ctx); err != nil {
            sc.cancel() // 取消所有其他goroutine
        }
    }()
}

func (sc *StructuredConcurrency) Wait() error {
    sc.wg.Wait()
    return sc.ctx.Err()
}

func (sc *StructuredConcurrency) Shutdown() {
    sc.cancel()
    sc.wg.Wait()
}
```

### 1.2.3 响应式并发

```go
// 响应式事件流
type EventStream[T any] struct {
    events chan T
    done   chan struct{}
    mu     sync.RWMutex
    subs   []chan T
}

func NewEventStream[T any]() *EventStream[T] {
    return &EventStream[T]{
        events: make(chan T, 100),
        done:   make(chan struct{}),
    }
}

func (es *EventStream[T]) Publish(event T) {
    select {
    case es.events <- event:
    case <-es.done:
    }
}

func (es *EventStream[T]) Subscribe() <-chan T {
    es.mu.Lock()
    defer es.mu.Unlock()
    
    ch := make(chan T, 10)
    es.subs = append(es.subs, ch)
    return ch
}

func (es *EventStream[T]) Start() {
    go func() {
        defer es.close()
        for {
            select {
            case event := <-es.events:
                es.broadcast(event)
            case <-es.done:
                return
            }
        }
    }()
}

func (es *EventStream[T]) broadcast(event T) {
    es.mu.RLock()
    defer es.mu.RUnlock()
    
    for _, sub := range es.subs {
        select {
        case sub <- event:
        default:
            // 订阅者处理不过来，跳过
        }
    }
}
```

## 1.3 现代化并发模式

### 1.3.1 管道模式2.0

```go
// 类型安全的管道
type Pipeline[T, U any] struct {
    stages []func(context.Context, <-chan T) <-chan U
}

func NewPipeline[T, U any]() *Pipeline[T, U] {
    return &Pipeline[T, U]{}
}

func (p *Pipeline[T, U]) AddStage(fn func(context.Context, <-chan T) <-chan U) *Pipeline[T, U] {
    p.stages = append(p.stages, fn)
    return p
}

func (p *Pipeline[T, U]) Execute(ctx context.Context, input <-chan T) <-chan U {
    current := input
    for _, stage := range p.stages {
        current = stage(ctx, current)
    }
    return current
}

// 使用示例
func ExamplePipeline() {
    pipeline := NewPipeline[int, string]().
        AddStage(func(ctx context.Context, input <-chan int) <-chan int {
            output := make(chan int)
            go func() {
                defer close(output)
                for n := range input {
                    select {
                    case output <- n * 2:
                    case <-ctx.Done():
                        return
                    }
                }
            }()
            return output
        }).
        AddStage(func(ctx context.Context, input <-chan int) <-chan string {
            output := make(chan string)
            go func() {
                defer close(output)
                for n := range input {
                    select {
                    case output <- fmt.Sprintf("Result: %d", n):
                    case <-ctx.Done():
                        return
                    }
                }
            }()
            return output
        })
    
    ctx := context.Background()
    input := make(chan int, 10)
    
    // 启动管道
    result := pipeline.Execute(ctx, input)
    
    // 发送数据
    go func() {
        defer close(input)
        for i := 1; i <= 5; i++ {
            input <- i
        }
    }()
    
    // 收集结果
    for res := range result {
        fmt.Println(res)
    }
}
```

### 1.3.2 背压控制

```go
// 自适应背压控制器
type BackpressureController struct {
    maxConcurrency int
    currentLoad    int32
    queue          chan func()
    mu             sync.RWMutex
}

func NewBackpressureController(maxConcurrency int) *BackpressureController {
    return &BackpressureController{
        maxConcurrency: maxConcurrency,
        queue:          make(chan func(), 1000),
    }
}

func (bpc *BackpressureController) Execute(fn func()) error {
    if atomic.LoadInt32(&bpc.currentLoad) >= int32(bpc.maxConcurrency) {
        // 背压：将任务加入队列
        select {
        case bpc.queue <- fn:
            return nil
        default:
            return errors.New("queue full, backpressure applied")
        }
    }
    
    atomic.AddInt32(&bpc.currentLoad, 1)
    go func() {
        defer atomic.AddInt32(&bpc.currentLoad, -1)
        fn()
    }()
    
    return nil
}

func (bpc *BackpressureController) ProcessQueue() {
    for fn := range bpc.queue {
        if atomic.LoadInt32(&bpc.currentLoad) < int32(bpc.maxConcurrency) {
            bpc.Execute(fn)
        } else {
            // 重新入队
            select {
            case bpc.queue <- fn:
            default:
                // 队列满，丢弃任务
            }
        }
    }
}
```

## 1.4 性能优化

### 1.4.1 无锁并发

```go
// 无锁环形缓冲区
type LockFreeRingBuffer[T any] struct {
    buffer []T
    mask   uint64
    head   uint64
    tail   uint64
}

func NewLockFreeRingBuffer[T any](size int) *LockFreeRingBuffer[T] {
    if size&(size-1) != 0 {
        panic("size must be power of 2")
    }
    
    return &LockFreeRingBuffer[T]{
        buffer: make([]T, size),
        mask:   uint64(size - 1),
    }
}

func (rb *LockFreeRingBuffer[T]) Push(item T) bool {
    head := atomic.LoadUint64(&rb.head)
    tail := atomic.LoadUint64(&rb.tail)
    
    if head-tail >= uint64(len(rb.buffer)) {
        return false // 缓冲区满
    }
    
    rb.buffer[head&rb.mask] = item
    atomic.StoreUint64(&rb.head, head+1)
    return true
}

func (rb *LockFreeRingBuffer[T]) Pop() (T, bool) {
    var zero T
    tail := atomic.LoadUint64(&rb.tail)
    head := atomic.LoadUint64(&rb.head)
    
    if tail >= head {
        return zero, false // 缓冲区空
    }
    
    item := rb.buffer[tail&rb.mask]
    atomic.StoreUint64(&rb.tail, tail+1)
    return item, true
}
```

### 1.4.2 内存池优化

```go
// 高性能对象池
type ObjectPool[T any] struct {
    pool    sync.Pool
    factory func() T
    reset   func(T)
}

func NewObjectPool[T any](factory func() T, reset func(T)) *ObjectPool[T] {
    return &ObjectPool[T]{
        factory: factory,
        reset:   reset,
        pool: sync.Pool{
            New: func() interface{} {
                return factory()
            },
        },
    }
}

func (op *ObjectPool[T]) Get() T {
    return op.pool.Get().(T)
}

func (op *ObjectPool[T]) Put(obj T) {
    if op.reset != nil {
        op.reset(obj)
    }
    op.pool.Put(obj)
}

// 使用示例
type Worker struct {
    ID   int
    Data []byte
}

func (w *Worker) Reset() {
    w.Data = w.Data[:0] // 重用切片
}

func ExampleObjectPool() {
    pool := NewObjectPool(
        func() *Worker {
            return &Worker{
                Data: make([]byte, 0, 1024),
            }
        },
        func(w *Worker) {
            w.Reset()
        },
    )
    
    // 使用对象池
    worker := pool.Get()
    defer pool.Put(worker)
    
    // 使用worker...
}
```

## 1.5 最佳实践

### 1.5.1 错误处理

```go
// 并发安全的错误收集器
type ErrorCollector struct {
    errors []error
    mu     sync.Mutex
}

func (ec *ErrorCollector) Add(err error) {
    if err != nil {
        ec.mu.Lock()
        ec.errors = append(ec.errors, err)
        ec.mu.Unlock()
    }
}

func (ec *ErrorCollector) Errors() []error {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    return append([]error(nil), ec.errors...)
}

func (ec *ErrorCollector) HasErrors() bool {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    return len(ec.errors) > 0
}
```

### 1.5.2 监控和指标

```go
// 并发指标收集器
type ConcurrencyMetrics struct {
    ActiveGoroutines int64
    CompletedTasks   int64
    FailedTasks      int64
    QueueSize        int64
}

func (cm *ConcurrencyMetrics) IncrementActive() {
    atomic.AddInt64(&cm.ActiveGoroutines, 1)
}

func (cm *ConcurrencyMetrics) DecrementActive() {
    atomic.AddInt64(&cm.ActiveGoroutines, -1)
}

func (cm *ConcurrencyMetrics) IncrementCompleted() {
    atomic.AddInt64(&cm.CompletedTasks, 1)
}

func (cm *ConcurrencyMetrics) IncrementFailed() {
    atomic.AddInt64(&cm.FailedTasks, 1)
}

func (cm *ConcurrencyMetrics) SetQueueSize(size int64) {
    atomic.StoreInt64(&cm.QueueSize, size)
}
```

## 1.6 实际应用

### 1.6.1 高并发HTTP服务

```go
// 高并发HTTP处理器
type ConcurrentHandler struct {
    workerPool *WorkerPool[HTTPRequest]
    metrics    *ConcurrencyMetrics
}

type HTTPRequest struct {
    ID      string
    Request *http.Request
    Writer  http.ResponseWriter
    Process func(*http.Request) (*http.Response, error)
}

func (ch *ConcurrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ch.metrics.IncrementActive()
    defer ch.metrics.DecrementActive()
    
    req := HTTPRequest{
        ID:      generateID(),
        Request: r,
        Writer:  w,
        Process: ch.processRequest,
    }
    
    if err := ch.workerPool.Execute(req); err != nil {
        ch.metrics.IncrementFailed()
        http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
        return
    }
    
    ch.metrics.IncrementCompleted()
}
```

### 1.6.2 实时数据处理

```go
// 实时数据处理器
type RealTimeProcessor[T any] struct {
    inputStream  *EventStream[T]
    processors   []func(T) T
    outputStream *EventStream[T]
    metrics      *ConcurrencyMetrics
}

func NewRealTimeProcessor[T any]() *RealTimeProcessor[T] {
    return &RealTimeProcessor[T]{
        inputStream:  NewEventStream[T](),
        outputStream: NewEventStream[T](),
        metrics:      &ConcurrencyMetrics{},
    }
}

func (rtp *RealTimeProcessor[T]) AddProcessor(fn func(T) T) {
    rtp.processors = append(rtp.processors, fn)
}

func (rtp *RealTimeProcessor[T]) Start(ctx context.Context) {
    rtp.inputStream.Start()
    rtp.outputStream.Start()
    
    go rtp.process(ctx)
}

func (rtp *RealTimeProcessor[T]) process(ctx context.Context) {
    for event := range rtp.inputStream.Subscribe() {
        rtp.metrics.IncrementActive()
        
        go func(e T) {
            defer rtp.metrics.DecrementActive()
            
            result := e
            for _, processor := range rtp.processors {
                result = processor(result)
            }
            
            rtp.outputStream.Publish(result)
            rtp.metrics.IncrementCompleted()
        }(event)
    }
}
```

---

**总结**: Go并发编程2.0通过泛型、结构化并发、响应式编程等现代化模式，为开发者提供了更安全、更高效的并发编程体验。这些模式不仅提高了代码的可维护性，还显著提升了系统的性能和可靠性。
