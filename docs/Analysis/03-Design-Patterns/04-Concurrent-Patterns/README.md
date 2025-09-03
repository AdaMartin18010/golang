# 3.4.1 并发模式分析

<!-- TOC START -->
- [3.4.1 并发模式分析](#341-并发模式分析)
  - [3.4.1.1 目录](#3411-目录)
  - [3.4.1.2 1. 并发基础](#3412-1-并发基础)
    - [3.4.1.2.1 并发与并行](#34121-并发与并行)
    - [3.4.1.2.2 并发模型](#34122-并发模型)
  - [3.4.1.3 2. Golang并发原语](#3413-2-golang并发原语)
    - [3.4.1.3.1 Goroutine](#34131-goroutine)
    - [3.4.1.3.2 Channel](#34132-channel)
    - [3.4.1.3.3 Select语句](#34133-select语句)
  - [3.4.1.4 3. 并发设计模式](#3414-3-并发设计模式)
    - [3.4.1.4.1 Worker Pool模式](#34141-worker-pool模式)
    - [3.4.1.4.2 Pipeline模式](#34142-pipeline模式)
    - [3.4.1.4.3 Fan-Out/Fan-In模式](#34143-fan-outfan-in模式)
  - [3.4.1.5 4. 同步机制](#3415-4-同步机制)
    - [3.4.1.5.1 Mutex](#34151-mutex)
    - [3.4.1.5.2 WaitGroup](#34152-waitgroup)
    - [3.4.1.5.3 Once](#34153-once)
  - [3.4.1.6 5. 并发数据结构](#3416-5-并发数据结构)
    - [3.4.1.6.1 并发Map](#34161-并发map)
    - [3.4.1.6.2 并发队列](#34162-并发队列)
  - [3.4.1.7 6. 性能优化](#3417-6-性能优化)
    - [3.4.1.7.1 内存池](#34171-内存池)
    - [3.4.1.7.2 工作窃取](#34172-工作窃取)
  - [3.4.1.8 7. 错误处理](#3418-7-错误处理)
    - [3.4.1.8.1 错误传播](#34181-错误传播)
    - [3.4.1.8.2 超时控制](#34182-超时控制)
  - [3.4.1.9 8. 最佳实践](#3419-8-最佳实践)
    - [3.4.1.9.1 避免竞态条件](#34191-避免竞态条件)
    - [3.4.1.9.2 避免死锁](#34192-避免死锁)
    - [3.4.1.9.3 资源管理](#34193-资源管理)
  - [3.4.1.10 9. 案例分析](#34110-9-案例分析)
    - [3.4.1.10.1 Web服务器](#341101-web服务器)
    - [3.4.1.10.2 数据处理管道](#341102-数据处理管道)
  - [3.4.1.11 参考资料](#34111-参考资料)
<!-- TOC END -->

## 3.4.1.1 目录

## 3.4.1.2 1. 并发基础

### 3.4.1.2.1 并发与并行

**定义 1.1** (并发): 并发是指多个任务在时间上重叠执行的能力。
**定义 1.2** (并行): 并行是指多个任务同时执行的能力。

**形式化表示**：
设 $T = \{t_1, t_2, ..., t_n\}$ 为任务集合，$S(t)$ 为任务 $t$ 的执行时间区间，则：

- **并发**: $\exists t_i, t_j \in T: S(t_i) \cap S(t_j) \neq \emptyset$
- **并行**: $\forall t_i, t_j \in T: S(t_i) \cap S(t_j) \neq \emptyset$

### 3.4.1.2.2 并发模型

**CSP模型** (Communicating Sequential Processes):

```go
// CSP模型示例
func producer(ch chan<- int) {
    for i := 0; i < 10; i++ {
        ch <- i // 发送数据
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

func consumer(ch <-chan int) {
    for value := range ch {
        fmt.Printf("Received: %d\n", value)
    }
}

func main() {
    ch := make(chan int)
    go producer(ch)
    go consumer(ch)
    time.Sleep(2 * time.Second)
}
```

## 3.4.1.3 2. Golang并发原语

### 3.4.1.3.1 Goroutine

**定义 2.1** (Goroutine): Goroutine是Go语言的轻量级线程，由Go运行时管理。

**数学表示**：
设 $G$ 为Goroutine集合，$M$ 为系统线程集合，则：
$$|G| \gg |M|$$

```go
// Goroutine基础用法
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("worker %d processing job %d\n", id, j)
        time.Sleep(time.Second)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // 启动3个worker
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // 发送任务
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)
    
    // 收集结果
    for a := 1; a <= 9; a++ {
        <-results
    }
}
```

### 3.4.1.3.2 Channel

**定义 2.2** (Channel): Channel是Golang中用于Goroutine间通信的管道。

**类型定义**：

```go
// Channel类型
type Channel[T any] interface {
    Send(T) error
    Receive() (T, error)
    Close() error
}

// 实现示例
type bufferedChannel[T any] struct {
    buffer chan T
    closed bool
    mutex  sync.RWMutex
}

func NewBufferedChannel[T any](capacity int) Channel[T] {
    return &bufferedChannel[T]{
        buffer: make(chan T, capacity),
    }
}

func (c *bufferedChannel[T]) Send(value T) error {
    c.mutex.RLock()
    if c.closed {
        c.mutex.RUnlock()
        return errors.New("channel closed")
    }
    c.mutex.RUnlock()
    
    select {
    case c.buffer <- value:
        return nil
    default:
        return errors.New("channel full")
    }
}

func (c *bufferedChannel[T]) Receive() (T, error) {
    var zero T
    select {
    case value, ok := <-c.buffer:
        if !ok {
            return zero, errors.New("channel closed")
        }
        return value, nil
    default:
        return zero, errors.New("no data available")
    }
}

func (c *bufferedChannel[T]) Close() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if c.closed {
        return errors.New("channel already closed")
    }
    
    c.closed = true
    close(c.buffer)
    return nil
}
```

### 3.4.1.3.3 Select语句

**定义 2.3** (Select): Select语句用于在多个Channel操作中进行非阻塞选择。

```go
// Select模式
func fanIn(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    
    go func() {
        defer close(out)
        
        for {
            select {
            case x := <-ch1:
                out <- x
            case x := <-ch2:
                out <- x
            case <-time.After(1 * time.Second):
                return
            }
        }
    }()
    
    return out
}

// 超时模式
func timeoutExample() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case result := <-ch:
        fmt.Printf("Received: %s\n", result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

## 3.4.1.4 3. 并发设计模式

### 3.4.1.4.1 Worker Pool模式

**定义 3.1** (Worker Pool): Worker Pool模式使用固定数量的Goroutine处理任务队列。

**数学分析**：
设 $W$ 为Worker数量，$T$ 为任务数量，$P$ 为处理时间，则：

- 总处理时间: $O(\frac{T \cdot P}{W})$
- 内存使用: $O(W + Q)$，其中 $Q$ 为队列大小

```go
// WorkerPool 工作池
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    resultQueue chan Result
    wg         sync.WaitGroup
}

type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID  int
    Data   interface{}
    Error  error
}

func NewWorkerPool(workers int, queueSize int) *WorkerPool {
    return &WorkerPool{
        workers:     workers,
        jobQueue:    make(chan Job, queueSize),
        resultQueue: make(chan Result, queueSize),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for job := range wp.jobQueue {
        result := wp.processJob(job)
        wp.resultQueue <- result
    }
}

func (wp *WorkerPool) processJob(job Job) Result {
    // 模拟处理时间
    time.Sleep(100 * time.Millisecond)
    
    return Result{
        JobID: job.ID,
        Data:  fmt.Sprintf("Processed job %d", job.ID),
    }
}

func (wp *WorkerPool) Submit(job Job) {
    wp.jobQueue <- job
}

func (wp *WorkerPool) Results() <-chan Result {
    return wp.resultQueue
}

func (wp *WorkerPool) Stop() {
    close(wp.jobQueue)
    wp.wg.Wait()
    close(wp.resultQueue)
}

// 使用示例
func main() {
    pool := NewWorkerPool(4, 100)
    pool.Start()
    
    // 提交任务
    for i := 0; i < 10; i++ {
        pool.Submit(Job{ID: i, Data: fmt.Sprintf("Task %d", i)})
    }
    
    // 收集结果
    go func() {
        for result := range pool.Results() {
            fmt.Printf("Result: %+v\n", result)
        }
    }()
    
    pool.Stop()
}
```

### 3.4.1.4.2 Pipeline模式

**定义 3.2** (Pipeline): Pipeline模式将复杂任务分解为多个阶段，每个阶段处理数据并传递给下一阶段。

```go
// Pipeline 管道模式
type Pipeline struct {
    stages []Stage
}

type Stage func(<-chan interface{}) <-chan interface{}

func NewPipeline(stages ...Stage) *Pipeline {
    return &Pipeline{stages: stages}
}

func (p *Pipeline) Execute(input <-chan interface{}) <-chan interface{} {
    current := input
    
    for _, stage := range p.stages {
        current = stage(current)
    }
    
    return current
}

// 具体阶段实现
func Stage1(input <-chan interface{}) <-chan interface{} {
    output := make(chan interface{})
    
    go func() {
        defer close(output)
        for value := range input {
            // 处理逻辑
            output <- fmt.Sprintf("Stage1: %v", value)
        }
    }()
    
    return output
}

func Stage2(input <-chan interface{}) <-chan interface{} {
    output := make(chan interface{})
    
    go func() {
        defer close(output)
        for value := range input {
            // 处理逻辑
            output <- fmt.Sprintf("Stage2: %v", value)
        }
    }()
    
    return output
}

// 使用示例
func main() {
    input := make(chan interface{})
    
    pipeline := NewPipeline(Stage1, Stage2)
    output := pipeline.Execute(input)
    
    // 发送输入
    go func() {
        for i := 0; i < 5; i++ {
            input <- i
        }
        close(input)
    }()
    
    // 接收输出
    for result := range output {
        fmt.Println(result)
    }
}
```

### 3.4.1.4.3 Fan-Out/Fan-In模式

**定义 3.3** (Fan-Out): 将输入分发到多个Goroutine处理。
**定义 3.4** (Fan-In): 将多个Goroutine的输出合并到一个Channel。

```go
// FanOut 分发模式
func FanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    
    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output
        
        go func(id int) {
            defer close(output)
            for value := range input {
                // 模拟处理
                time.Sleep(100 * time.Millisecond)
                output <- value * value
            }
        }(i)
    }
    
    return outputs
}

// FanIn 合并模式
func FanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup
    
    for _, input := range inputs {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for value := range ch {
                output <- value
            }
        }(input)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}

// 使用示例
func main() {
    input := make(chan int)
    
    // Fan-Out
    outputs := FanOut(input, 3)
    
    // Fan-In
    result := FanIn(outputs...)
    
    // 发送数据
    go func() {
        for i := 0; i < 10; i++ {
            input <- i
        }
        close(input)
    }()
    
    // 接收结果
    for value := range result {
        fmt.Printf("Result: %d\n", value)
    }
}
```

## 3.4.1.5 4. 同步机制

### 3.4.1.5.1 Mutex

**定义 4.1** (Mutex): 互斥锁用于保护共享资源，确保同一时间只有一个Goroutine访问。

```go
// SafeCounter 线程安全计数器
type SafeCounter struct {
    value int
    mutex sync.RWMutex
}

func (sc *SafeCounter) Increment() {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()
    sc.value++
}

func (sc *SafeCounter) Decrement() {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()
    sc.value--
}

func (sc *SafeCounter) GetValue() int {
    sc.mutex.RLock()
    defer sc.mutex.RUnlock()
    return sc.value
}

// 读写锁优化
type OptimizedCounter struct {
    value int
    mutex sync.RWMutex
}

func (oc *OptimizedCounter) Increment() {
    oc.mutex.Lock()
    defer oc.mutex.Unlock()
    oc.value++
}

func (oc *OptimizedCounter) GetValue() int {
    oc.mutex.RLock()
    defer oc.mutex.RUnlock()
    return oc.value
}
```

### 3.4.1.5.2 WaitGroup

**定义 4.2** (WaitGroup): WaitGroup用于等待一组Goroutine完成。

```go
// WaitGroup使用示例
func processItems(items []string) {
    var wg sync.WaitGroup
    
    for _, item := range items {
        item := item // 创建副本
        wg.Add(1)
        go func(item string) {
            defer wg.Done()
            processItem(item)
        }(item)
    }
    
    wg.Wait()
    fmt.Println("All items processed")
}

func processItem(item string) {
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Processed: %s\n", item)
}
```

### 3.4.1.5.3 Once

**定义 4.3** (Once): Once确保某个函数只执行一次。

```go
// Singleton模式使用Once
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            data: "Initialized",
        }
    })
    return instance
}
```

## 3.4.1.6 5. 并发数据结构

### 3.4.1.6.1 并发Map

```go
// ConcurrentMap 并发安全的Map
type ConcurrentMap[K comparable, V any] struct {
    data map[K]V
    mutex sync.RWMutex
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
    return &ConcurrentMap[K, V]{
        data: make(map[K]V),
    }
}

func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    cm.data[key] = value
}

func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    value, exists := cm.data[key]
    return value, exists
}

func (cm *ConcurrentMap[K, V]) Delete(key K) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    delete(cm.data, key)
}

func (cm *ConcurrentMap[K, V]) Range(f func(K, V) bool) {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    for key, value := range cm.data {
        if !f(key, value) {
            break
        }
    }
}
```

### 3.4.1.6.2 并发队列

```go
// ConcurrentQueue 并发安全队列
type ConcurrentQueue[T any] struct {
    data []T
    mutex sync.Mutex
    cond  *sync.Cond
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
    q := &ConcurrentQueue[T]{}
    q.cond = sync.NewCond(&q.mutex)
    return q
}

func (cq *ConcurrentQueue[T]) Enqueue(item T) {
    cq.mutex.Lock()
    defer cq.mutex.Unlock()
    
    cq.data = append(cq.data, item)
    cq.cond.Signal()
}

func (cq *ConcurrentQueue[T]) Dequeue() (T, bool) {
    cq.mutex.Lock()
    defer cq.mutex.Unlock()
    
    for len(cq.data) == 0 {
        cq.cond.Wait()
    }
    
    item := cq.data[0]
    cq.data = cq.data[1:]
    return item, true
}

func (cq *ConcurrentQueue[T]) Size() int {
    cq.mutex.Lock()
    defer cq.mutex.Unlock()
    return len(cq.data)
}
```

## 3.4.1.7 6. 性能优化

### 3.4.1.7.1 内存池

```go
// ObjectPool 对象池
type ObjectPool[T any] struct {
    pool chan T
    new  func() T
    reset func(T)
}

func NewObjectPool[T any](size int, newFunc func() T, resetFunc func(T)) *ObjectPool[T] {
    return &ObjectPool[T]{
        pool:  make(chan T, size),
        new:   newFunc,
        reset: resetFunc,
    }
}

func (op *ObjectPool[T]) Get() T {
    select {
    case obj := <-op.pool:
        return obj
    default:
        return op.new()
    }
}

func (op *ObjectPool[T]) Put(obj T) {
    if op.reset != nil {
        op.reset(obj)
    }
    
    select {
    case op.pool <- obj:
    default:
        // 池已满，丢弃对象
    }
}

// 使用示例
type Buffer struct {
    data []byte
}

func NewBuffer() Buffer {
    return Buffer{data: make([]byte, 1024)}
}

func ResetBuffer(buf Buffer) {
    buf.data = buf.data[:0]
}

func main() {
    pool := NewObjectPool(10, NewBuffer, ResetBuffer)
    
    // 获取对象
    buf := pool.Get()
    
    // 使用对象
    buf.data = append(buf.data, "Hello"...)
    
    // 归还对象
    pool.Put(buf)
}
```

### 3.4.1.7.2 工作窃取

```go
// WorkStealingScheduler 工作窃取调度器
type WorkStealingScheduler struct {
    workers []*Worker
    queues  []*Deque
}

type Worker struct {
    id    int
    queue *Deque
    steal func(int) interface{}
}

type Deque struct {
    data []interface{}
    mutex sync.Mutex
}

func (d *Deque) PushFront(item interface{}) {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    d.data = append([]interface{}{item}, d.data...)
}

func (d *Deque) PopFront() (interface{}, bool) {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    
    if len(d.data) == 0 {
        return nil, false
    }
    
    item := d.data[0]
    d.data = d.data[1:]
    return item, true
}

func (d *Deque) PopBack() (interface{}, bool) {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    
    if len(d.data) == 0 {
        return nil, false
    }
    
    item := d.data[len(d.data)-1]
    d.data = d.data[:len(d.data)-1]
    return item, true
}
```

## 3.4.1.8 7. 错误处理

### 3.4.1.8.1 错误传播

```go
// ErrorGroup 错误组
func processWithErrorGroup(items []string) error {
    var eg errgroup.Group
    
    for _, item := range items {
        item := item // 创建副本
        eg.Go(func() error {
            return processItemWithError(item)
        })
    }
    
    return eg.Wait()
}

func processItemWithError(item string) error {
    if item == "error" {
        return fmt.Errorf("processing error for item: %s", item)
    }
    
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Processed: %s\n", item)
    return nil
}
```

### 3.4.1.8.2 超时控制

```go
// TimeoutWrapper 超时包装器
func WithTimeout[T any](fn func() (T, error), timeout time.Duration) (T, error) {
    var zero T
    
    done := make(chan struct{})
    var result T
    var err error
    
    go func() {
        result, err = fn()
        close(done)
    }()
    
    select {
    case <-done:
        return result, err
    case <-time.After(timeout):
        return zero, errors.New("operation timeout")
    }
}

// 使用示例
func main() {
    result, err := WithTimeout(func() (string, error) {
        time.Sleep(2 * time.Second)
        return "success", nil
    }, 1*time.Second)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %s\n", result)
    }
}
```

## 3.4.1.9 8. 最佳实践

### 3.4.1.9.1 避免竞态条件

```go
// 错误示例：竞态条件
var counter int

func increment() {
    counter++ // 竞态条件
}

// 正确示例：使用互斥锁
type SafeCounter struct {
    value int
    mutex sync.Mutex
}

func (sc *SafeCounter) Increment() {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()
    sc.value++
}
```

### 3.4.1.9.2 避免死锁

```go
// 死锁示例
func deadlockExample() {
    var mu1, mu2 sync.Mutex
    
    go func() {
        mu1.Lock()
        time.Sleep(100 * time.Millisecond)
        mu2.Lock()
        // ...
        mu2.Unlock()
        mu1.Unlock()
    }()
    
    mu2.Lock()
    time.Sleep(100 * time.Millisecond)
    mu1.Lock()
    // ...
    mu1.Unlock()
    mu2.Unlock()
}

// 避免死锁：固定顺序
func safeExample() {
    var mu1, mu2 sync.Mutex
    
    go func() {
        mu1.Lock()
        defer mu1.Unlock()
        time.Sleep(100 * time.Millisecond)
        mu2.Lock()
        defer mu2.Unlock()
        // ...
    }()
    
    mu1.Lock()
    defer mu1.Unlock()
    time.Sleep(100 * time.Millisecond)
    mu2.Lock()
    defer mu2.Unlock()
    // ...
}
```

### 3.4.1.9.3 资源管理

```go
// 资源管理示例
type ResourceManager struct {
    resources chan *Resource
    maxSize   int
}

type Resource struct {
    ID   int
    Data string
}

func NewResourceManager(maxSize int) *ResourceManager {
    rm := &ResourceManager{
        resources: make(chan *Resource, maxSize),
        maxSize:   maxSize,
    }
    
    // 预分配资源
    for i := 0; i < maxSize; i++ {
        rm.resources <- &Resource{ID: i, Data: fmt.Sprintf("Resource %d", i)}
    }
    
    return rm
}

func (rm *ResourceManager) Acquire() *Resource {
    return <-rm.resources
}

func (rm *ResourceManager) Release(resource *Resource) {
    select {
    case rm.resources <- resource:
    default:
        // 池已满，丢弃资源
    }
}
```

## 3.4.1.10 9. 案例分析

### 3.4.1.10.1 Web服务器

```go
// ConcurrentWebServer 并发Web服务器
type ConcurrentWebServer struct {
    listener net.Listener
    handler  http.Handler
    pool     *WorkerPool
}

func NewConcurrentWebServer(addr string, workers int) (*ConcurrentWebServer, error) {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return nil, err
    }
    
    return &ConcurrentWebServer{
        listener: listener,
        handler:  http.DefaultServeMux,
        pool:     NewWorkerPool(workers, 1000),
    }, nil
}

func (s *ConcurrentWebServer) Start() {
    s.pool.Start()
    
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            log.Printf("Accept error: %v", err)
            continue
        }
        
        s.pool.Submit(Job{
            ID:   time.Now().UnixNano(),
            Data: conn,
        })
    }
}

func (s *ConcurrentWebServer) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // 处理HTTP请求
    server := &http.Server{
        Handler: s.handler,
    }
    
    server.Serve(&singleConnListener{conn: conn})
}

type singleConnListener struct {
    conn net.Conn
    done bool
}

func (l *singleConnListener) Accept() (net.Conn, error) {
    if l.done {
        return nil, net.ErrClosed
    }
    l.done = true
    return l.conn, nil
}

func (l *singleConnListener) Close() error {
    return l.conn.Close()
}

func (l *singleConnListener) Addr() net.Addr {
    return l.conn.LocalAddr()
}
```

### 3.4.1.10.2 数据处理管道

```go
// DataProcessingPipeline 数据处理管道
type DataProcessingPipeline struct {
    input  <-chan Data
    output chan<- Result
    stages []ProcessingStage
}

type Data struct {
    ID   int
    Value string
}

type Result struct {
    DataID int
    ProcessedValue string
    Error error
}

type ProcessingStage func(<-chan Data) <-chan Data

func NewDataProcessingPipeline(input <-chan Data, output chan<- Result, stages ...ProcessingStage) *DataProcessingPipeline {
    return &DataProcessingPipeline{
        input:  input,
        output: output,
        stages: stages,
    }
}

func (p *DataProcessingPipeline) Start() {
    current := p.input
    
    // 应用所有阶段
    for _, stage := range p.stages {
        current = stage(current)
    }
    
    // 收集结果
    go func() {
        defer close(p.output)
        for data := range current {
            result := Result{
                DataID: data.ID,
                ProcessedValue: data.Value,
            }
            p.output <- result
        }
    }()
}

// 具体处理阶段
func ValidationStage(input <-chan Data) <-chan Data {
    output := make(chan Data)
    
    go func() {
        defer close(output)
        for data := range input {
            if len(data.Value) > 0 {
                output <- data
            }
        }
    }()
    
    return output
}

func TransformationStage(input <-chan Data) <-chan Data {
    output := make(chan Data)
    
    go func() {
        defer close(output)
        for data := range input {
            data.Value = strings.ToUpper(data.Value)
            output <- data
        }
    }()
    
    return output
}
```

---

## 3.4.1.11 参考资料

1. [Go并发编程指南](https://golang.org/doc/effective_go.html#concurrency)
2. [CSP模型论文](https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf)
3. [并发编程模式](https://en.wikipedia.org/wiki/Concurrent_programming)
4. [Go内存模型](https://golang.org/ref/mem)
5. [并发数据结构](https://en.wikipedia.org/wiki/Concurrent_data_structure)

---

*本文档涵盖了Golang并发编程的核心概念、设计模式和最佳实践，为构建高性能并发应用提供指导。*
