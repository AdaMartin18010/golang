# Go并发模式深度解析

> 基于CSP形式化模型的并发设计模式与工程实践

---

## 一、并发模式的CSP基础

### 1.1 模式的形式化定义

```
并发模式 = (Process Structure, Communication Protocol, Synchronization Strategy)

CSP视角:
├─ Process: Goroutine的独立控制流
├─ Communication: Channel事件同步
└─ Composition: 并行组合算子

形式化验证目标:
├─ 死锁自由: 无循环等待
├─ 活锁自由: 持续进展保证
├─ 安全性: 不变式保持
└─ 活性: 期望事件最终发生
```

### 1.2 并发模式分类学

```
按通信模式分类:
────────────────────────────────────────
同步模式:
├── Rendezvous (会合)        ── 直接channel通信
├── Barrier (屏障)           ── 多goroutine同步点
└── Semaphore (信号量)       ── 容量控制的channel

异步模式:
├── Worker Pool              ── 任务分发处理
├── Pipeline                 ── 数据流阶段处理
├── Fan-out/Fan-in           ── 并行分发与聚合
└── Pub/Sub                  ── 广播通知

控制模式:
├── Context Cancellation     ── 级联取消传播
├── Timeout                  ── 时间约束
├── Rate Limiting            ── 流量控制
└── Circuit Breaker          ── 故障隔离
```

---

## 二、基础同步模式

### 2.1 Rendezvous (会合模式)

```
CSP形式化:
────────────────────────────────────────
P = a → P'   (进程P执行事件a后继续为P')
Q = a → Q'   (进程Q执行事件a后继续为Q')
P ∥ Q = a → (P' ∥ Q')  (并行组合要求同步执行a)

Go实现:
ch := make(chan struct{})  // 无缓冲channel

go func() {  // Process P
    // 准备阶段
    ch <- struct{}{}  // 发送事件
    // 继续P'
}()

// Process Q
<-ch  // 接收事件 (同步点)
// 继续Q'

语义保证:
发送者发送 ≺ 接收者接收 (Happens-before)
双方在该点同步后才继续

代码示例:
// 会合模式: 确保goroutine准备就绪
func rendezvousPattern() {
    ready := make(chan struct{})
    done := make(chan struct{})

    go func() {
        fmt.Println("Worker: 初始化中...")
        time.Sleep(time.Second)
        fmt.Println("Worker: 准备就绪")
        close(ready)  // 信号: 准备完成

        // 等待主goroutine确认
        <-done
        fmt.Println("Worker: 开始工作")
    }()

    <-ready  // 会合点: 等待worker就绪
    fmt.Println("Main: Worker已就绪")
    close(done)  // 确认
    time.Sleep(time.Second)
}

// 反例: 无同步导致竞态
func noRendezvous() {
    var result int

    go func() {
        result = 42  // 写操作
    }()

    // 可能读到0或42，无保证
    fmt.Println(result)
}
```

### 2.2 Barrier (屏障模式)

```
形式化:
Barrier(n) = 等待n个参与者到达后才放行

Go实现:
type Barrier struct {
    count int
    mu    sync.Mutex
    cond  *sync.Cond
}

func NewBarrier(n int) *Barrier {
    b := &Barrier{count: n}
    b.cond = sync.NewCond(&b.mu)
    return b
}

func (b *Barrier) Wait() {
    b.mu.Lock()
    b.count--
    if b.count == 0 {
        b.cond.Broadcast()  // 最后一个到达，放行所有
    } else {
        b.cond.Wait()       // 等待其他参与者
    }
    b.mu.Unlock()
}

应用:
├─ 并行计算阶段同步
├─ 批量处理等待
└─ 测试协调

代码示例:
// 并行计算阶段同步
func parallelPhases() {
    const workers = 4
    barrier := NewBarrier(workers)

    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            // 阶段1
            fmt.Printf("Worker %d: 阶段1开始\n", id)
            time.Sleep(time.Duration(id*100) * time.Millisecond)
            fmt.Printf("Worker %d: 阶段1完成\n", id)

            barrier.Wait()  // 等待所有worker完成阶段1

            // 阶段2 (所有worker同时开始)
            fmt.Printf("Worker %d: 阶段2开始\n", id)
        }(i)
    }
    wg.Wait()
}

// 使用sync.WaitGroup的简化屏障
func waitGroupBarrier() {
    const workers = 4

    var phase1 sync.WaitGroup
    var phase2 sync.WaitGroup

    phase1.Add(workers)
    phase2.Add(workers)

    for i := 0; i < workers; i++ {
        go func(id int) {
            // 阶段1
            time.Sleep(time.Duration(id) * time.Millisecond)
            phase1.Done()
            phase1.Wait()  // 屏障

            // 阶段2
            fmt.Printf("Worker %d: 阶段2\n", id)
            phase2.Done()
        }(i)
    }
    phase2.Wait()
}
```

---

## 三、Pipeline模式

### 3.1 阶段组合的形式化

```
Pipeline代数:
────────────────────────────────────────
Stage = Input → Output  (类型转换)
Pipeline = Stage₁ ∘ Stage₂ ∘ ... ∘ Stageₙ

性质:
├─ 结合律: (A ∘ B) ∘ C = A ∘ (B ∘ C)
├─ 可并行: 各阶段并发执行
└─ 背压传播: 下游慢速自动限流上游

Go实现框架:
func Stage[T, R any](in <-chan T, fn func(T) R) <-chan R {
    out := make(chan R)
    go func() {
        defer close(out)
        for v := range in {
            select {
            case out <- fn(v):
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

类型安全Pipeline:
gen := Generator[int]()           // chan int
sq := Stage(gen, square)          // chan int → chan int
filter := Stage(sq, isEven)       // chan int → chan bool
out := Stage(filter, toString)    // chan bool → chan string
```

### 3.2 完整Pipeline示例

```
Pipeline代码示例:
────────────────────────────────────────

// 生成器阶段
func Generator(start, end int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := start; i < end; i++ {
            out <- i
        }
    }()
    return out
}

// 平方阶段
func Square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// 过滤阶段
func Filter(in <-chan int, predicate func(int) bool) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if predicate(n) {
                out <- n
            }
        }
    }()
    return out
}

// 消费阶段
func Consumer(in <-chan int, process func(int)) {
    for n := range in {
        process(n)
    }
}

// 组合Pipeline
func pipelineExample() {
    // pipeline: generate -> square -> filter(even) -> print
    gen := Generator(1, 100)
    sq := Square(gen)
    even := Filter(sq, func(n int) bool { return n%2 == 0 })

    Consumer(even, func(n int) {
        fmt.Println(n)
    })
}

// 带缓冲的Pipeline (提高吞吐量)
func BufferedStage[T, R any](in <-chan T, fn func(T) R, bufSize int) <-chan R {
    out := make(chan R, bufSize)
    go func() {
        defer close(out)
        for v := range in {
            out <- fn(v)
        }
    }()
    return out
}
```

### 3.3 错误传播机制

```
错误处理的形式化:
────────────────────────────────────────
Stage = Input → (Output, Error)

Go实现 (结果包含错误):
type Result[T any] struct {
    Value T
    Error error
}

func StageWithError[T, R any](
    in <-chan T,
    fn func(T) (R, error),
) <-chan Result[R] {
    out := make(chan Result[R])
    go func() {
        defer close(out)
        for v := range in {
            r, err := fn(v)
            select {
            case out <- Result[R]{r, err}:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

模式选择:
├─ 快速失败: 错误时关闭pipeline
├─ 错误传递: 结果包含错误继续处理
└─ 错误聚合: 收集所有错误后报告

代码示例:
// 快速失败模式
func pipelineFastFail() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    gen := GeneratorWithContext(ctx, 1, 100)
    processed := ProcessWithError(gen, func(n int) (int, error) {
        if n == 50 {
            return 0, fmt.Errorf("error at %d", n)
        }
        return n * 2, nil
    })

    for result := range processed {
        if result.Error != nil {
            cancel()  // 快速失败
            log.Fatal(result.Error)
        }
        fmt.Println(result.Value)
    }
}

// 错误聚合模式
func pipelineErrorAggregation() {
    gen := Generator(1, 100)
    processed := ProcessWithError(gen, riskyOperation)

    var errors []error
    for result := range processed {
        if result.Error != nil {
            errors = append(errors, result.Error)
            continue
        }
        fmt.Println(result.Value)
    }

    if len(errors) > 0 {
        log.Printf("处理完成，%d 个错误", len(errors))
    }
}
```

---

## 四、Worker Pool模式

### 4.1 基础Worker Pool

```
Worker Pool 的形式化:
────────────────────────────────────────
Pool = (Workers, JobQueue, ResultQueue)
Workers = {w₁, w₂, ..., wₙ}  where n = capacity

不变式:
├─ 0 ≤ active_workers ≤ capacity
├─ JobQueue 满时 Submit 阻塞或失败
└─ ResultQueue 消费慢时结果缓冲或丢弃

Go实现:
type Pool struct {
    workers int
    jobs    chan Job
    results chan Result
    wg      sync.WaitGroup
}

func NewPool(workers int) *Pool {
    p := &Pool{
        workers: workers,
        jobs:    make(chan Job),
        results: make(chan Result, workers),
    }
    for i := 0; i < workers; i++ {
        p.wg.Add(1)
        go p.worker()
    }
    return p
}

func (p *Pool) worker() {
    defer p.wg.Done()
    for job := range p.jobs {
        p.results <- job.Process()
    }
}

终止协议:
1. Close(jobs) - 无新任务
2. WaitGroup.Wait() - 等待完成
3. Close(results) - 清理资源
```

### 4.2 完整Worker Pool实现

```
完整代码实现:
────────────────────────────────────────

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID int
    Value string
    Error error
}

type WorkerPool struct {
    numWorkers int
    jobs       chan Job
    results    chan Result
    wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    return &WorkerPool{
        numWorkers: numWorkers,
        jobs:       make(chan Job, numWorkers*2),
        results:    make(chan Result, numWorkers*2),
    }
}

func (wp *WorkerPool) Start(ctx context.Context, processor func(Job) Result) {
    for i := 0; i < wp.numWorkers; i++ {
        wp.wg.Add(1)
        go func(id int) {
            defer wp.wg.Done()
            for {
                select {
                case job, ok := <-wp.jobs:
                    if !ok {
                        return
                    }
                    wp.results <- processor(job)
                case <-ctx.Done():
                    return
                }
            }
        }(i)
    }

    // 结果收集goroutine
    go func() {
        wp.wg.Wait()
        close(wp.results)
    }()
}

func (wp *WorkerPool) Submit(job Job) bool {
    select {
    case wp.jobs <- job:
        return true
    default:
        return false  // 队列满
    }
}

func (wp *WorkerPool) Results() <-chan Result {
    return wp.results
}

func (wp *WorkerPool) Stop() {
    close(wp.jobs)
}

// 使用示例
func workerPoolExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    pool := NewWorkerPool(4)

    pool.Start(ctx, func(job Job) Result {
        // 处理任务
        time.Sleep(100 * time.Millisecond)
        return Result{
            JobID: job.ID,
            Value: fmt.Sprintf("Processed: %s", job.Data),
        }
    })

    // 提交任务
    go func() {
        for i := 0; i < 20; i++ {
            pool.Submit(Job{
                ID:   i,
                Data: fmt.Sprintf("data-%d", i),
            })
        }
        pool.Stop()
    }()

    // 收集结果
    for result := range pool.Results() {
        fmt.Println(result.Value)
    }
}
```

### 4.3 动态Worker Pool

```
动态调整Worker数量:
────────────────────────────────────────

type DynamicPool struct {
    minWorkers int
    maxWorkers int
    jobs       chan Job
    results    chan Result
    active     int32
    mu         sync.Mutex
}

func (dp *DynamicPool) adjustWorkers() {
    queueLen := len(dp.jobs)
    active := int(atomic.LoadInt32(&dp.active))

    // 队列压力大，增加worker
    if queueLen > active*2 && active < dp.maxWorkers {
        dp.mu.Lock()
        if int(atomic.LoadInt32(&dp.active)) < dp.maxWorkers {
            atomic.AddInt32(&dp.active, 1)
            go dp.worker()
        }
        dp.mu.Unlock()
    }

    // 队列空闲，减少worker (可选)
    // ...
}

func (dp *DynamicPool) worker() {
    defer atomic.AddInt32(&dp.active, -1)

    idleTimer := time.NewTimer(30 * time.Second)
    defer idleTimer.Stop()

    for {
        select {
        case job, ok := <-dp.jobs:
            if !ok {
                return
            }
            idleTimer.Reset(30 * time.Second)
            dp.results <- job.Process()

        case <-idleTimer.C:
            // 空闲超时，worker退出
            return
        }
    }
}
```

---

## 五、Fan-out/Fan-in模式

### 5.1 并行分发的形式化

```
Fan-out:
────────────────────────────────────────
Input → [Distributor] → Output₁, Output₂, ..., Outputₙ

Distributor逻辑:
├─ 轮询: round-robin分配
├─ 哈希: 相同key到同一worker
└─ 随机: 负载均衡

Fan-in:
Output₁ →
Output₂ → [Merger] → Combined Output
...    →
Outputₙ →

CSP表示:
FanOut = ∥ᵢ (Pᵢ(Workerᵢ))
FanIn = Merge(∪ᵢ Outputᵢ)
```

### 5.2 Go实现

```
Fan-out实现:
func FanOut[T any](in <-chan T, n int) []<-chan T {
    outs := make([]<-chan T, n)
    for i := 0; i < n; i++ {
        ch := make(chan T)
        outs[i] = ch
        go func() {
            defer close(ch)
            for v := range in {
                ch <- v
            }
        }()
    }
    return outs
}

Fan-in实现:
func FanIn[T any](ins ...<-chan T) <-chan T {
    out := make(chan T)
    var wg sync.WaitGroup
    wg.Add(len(ins))

    for _, in := range ins {
        go func(ch <-chan T) {
            defer wg.Done()
            for v := range ch {
                out <- v
            }
        }(in)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

完整示例:
func fanOutFanInExample() {
    // 输入
    input := make(chan int)
    go func() {
        defer close(input)
        for i := 0; i < 100; i++ {
            input <- i
        }
    }()

    // Fan-out到4个worker
    workers := make([]<-chan int, 4)
    for i := 0; i < 4; i++ {
        workers[i] = process(input, i)
    }

    // Fan-in合并结果
    merged := FanIn(workers...)

    // 消费
    for result := range merged {
        fmt.Println(result)
    }
}

func process(in <-chan int, workerID int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            // 模拟处理
            time.Sleep(10 * time.Millisecond)
            out <- n * n
        }
    }()
    return out
}

// 有序Fan-in (保持顺序)
func OrderedFanIn[T any](ins ...<-chan T) <-chan T {
    out := make(chan T)

    go func() {
        defer close(out)

        // 为每个输入channel创建反射case
        cases := make([]reflect.SelectCase, len(ins))
        for i, ch := range ins {
            cases[i] = reflect.SelectCase{
                Dir:  reflect.SelectRecv,
                Chan: reflect.ValueOf(ch),
            }
        }

        remaining := len(ins)
        for remaining > 0 {
            chosen, value, ok := reflect.Select(cases)
            if !ok {
                // 此channel已关闭
                cases[chosen].Chan = reflect.ValueOf(nil)
                remaining--
                continue
            }
            out <- value.Interface().(T)
        }
    }()

    return out
}
```

---

## 六、Context模式

### 6.1 取消传播的形式化

```
Context作为信号:
────────────────────────────────────────
Context = (ValueMap, CancelChannel, Deadline)

取消传播规则:
Parent.Cancel() → ∀Child: Child.Cancel()

CSP视角:
Cancel = 广播事件到所有监听goroutine

Go实现语义:
ctx, cancel := context.WithCancel(parent)
├─ cancel(): 关闭ctx.Done() channel
├─ 衍生context继承取消信号
└─ 形成取消树
```

### 6.2 Context完整示例

```
Context使用模式:
────────────────────────────────────────

// 1. 超时控制
func operationWithTimeout() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    select {
    case result := <-doWork(ctx):
        return result
    case <-ctx.Done():
        return ctx.Err()  // context.DeadlineExceeded
    }
}

// 2. 手动取消
func operationWithCancel() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        time.Sleep(2 * time.Second)
        cancel()  // 触发取消
    }()

    <-ctx.Done()
    fmt.Println("取消原因:", ctx.Err())
}

// 3. 级联取消
func cascadeCancel() {
    rootCtx, rootCancel := context.WithCancel(context.Background())
    defer rootCancel()

    child1Ctx, _ := context.WithCancel(rootCtx)
    child2Ctx, _ := context.WithTimeout(child1Ctx, 10*time.Second)

    rootCancel()  // 取消root，child1和child2都被取消

    fmt.Println("root:", rootCtx.Err())
    fmt.Println("child1:", child1Ctx.Err())
    fmt.Println("child2:", child2Ctx.Err())
}

// 4. 传递值
func contextWithValue() {
    type userIDKey struct{}
    type traceIDKey struct{}

    ctx := context.Background()
    ctx = context.WithValue(ctx, userIDKey{}, "user-123")
    ctx = context.WithValue(ctx, traceIDKey{}, "trace-456")

    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    type userIDKey struct{}
    userID := ctx.Value(userIDKey{}).(string)
    fmt.Println("处理用户:", userID)
}

// 5. 实际应用: HTTP服务器
func httpServerWithContext() {
    http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        data, err := fetchData(ctx)
        if err != nil {
            if err == context.DeadlineExceeded {
                http.Error(w, "Timeout", http.StatusGatewayTimeout)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(data)
    })

    server := &http.Server{
        Addr:         ":8080",
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    server.ListenAndServe()
}

func fetchData(ctx context.Context) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com/data", nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

---

## 七、高级并发模式

### 7.1 Circuit Breaker (熔断器)

```
状态机模型:
────────────────────────────────────────
States: Closed → Open → Half-Open → Closed

Closed: 正常通过，记录失败率
Open: 快速失败，拒绝请求
Half-Open: 试探性放行，检测恢复

形式化:
FailureRate = Failures / Total
If FailureRate > Threshold: Closed → Open
After Timeout: Open → Half-Open
If ProbeSuccess: Half-Open → Closed
If ProbeFail: Half-Open → Open

Go实现:
type CircuitBreaker struct {
    state       State
    failures    int
    successes   int
    threshold   int
    timeout     time.Duration
    lastFailure time.Time
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.allow() {
        return ErrCircuitOpen
    }

    err := fn()
    cb.recordResult(err)
    return err
}

完整实现:
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreaker struct {
    name        string
    maxFailures int
    timeout     time.Duration

    mu          sync.Mutex
    state       State
    failures    int
    lastFailure time.Time
    halfOpenReqs int
}

func NewCircuitBreaker(name string, maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        name:        name,
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       StateClosed,
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    state := cb.state
    cb.mu.Unlock()

    switch state {
    case StateOpen:
        if cb.shouldAttemptReset() {
            cb.setState(StateHalfOpen)
        } else {
            return fmt.Errorf("circuit breaker open")
        }

    case StateHalfOpen:
        // 限制半开状态的请求数
        cb.mu.Lock()
        if cb.halfOpenReqs >= 1 {
            cb.mu.Unlock()
            return fmt.Errorf("circuit breaker half-open, too many requests")
        }
        cb.halfOpenReqs++
        cb.mu.Unlock()
    }

    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err == nil {
        cb.onSuccess()
    } else {
        cb.onFailure()
    }
}

func (cb *CircuitBreaker) onSuccess() {
    switch cb.state {
    case StateClosed:
        cb.failures = 0
    case StateHalfOpen:
        cb.setState(StateClosed)
        cb.failures = 0
        cb.halfOpenReqs = 0
    }
}

func (cb *CircuitBreaker) onFailure() {
    cb.failures++
    cb.lastFailure = time.Now()

    switch cb.state {
    case StateClosed:
        if cb.failures >= cb.maxFailures {
            cb.setState(StateOpen)
        }
    case StateHalfOpen:
        cb.setState(StateOpen)
        cb.halfOpenReqs = 0
    }
}

func (cb *CircuitBreaker) setState(state State) {
    fmt.Printf("CircuitBreaker %s: %v -> %v\n", cb.name, cb.state, state)
    cb.state = state
}

func (cb *CircuitBreaker) shouldAttemptReset() bool {
    return time.Since(cb.lastFailure) > cb.timeout
}

// 使用示例
func circuitBreakerExample() {
    cb := NewCircuitBreaker("api", 3, 5*time.Second)

    for i := 0; i < 10; i++ {
        err := cb.Call(func() error {
            // 模拟API调用
            if i < 5 {
                return fmt.Errorf("connection failed")
            }
            fmt.Println("Request succeeded")
            return nil
        })

        if err != nil {
            fmt.Printf("Request %d failed: %v\n", i, err)
        }
    }
}
```

### 7.2 Rate Limiter (限流器)

```
限流算法形式化:
────────────────────────────────────────

令牌桶 (Token Bucket):
Bucket = (Tokens, Capacity, RefillRate)
允许请求 ⟺ Tokens ≥ RequestCost
请求后: Tokens -= RequestCost
定时: Tokens = min(Capacity, Tokens + RefillRate × Δt)

漏桶 (Leaky Bucket):
Bucket = (Queue, LeakRate)
到达请求入队，队列满则拒绝
定时漏出处理

Go实现 (golang.org/x/time/rate):
lim := rate.NewLimiter(rate.Every(100*time.Millisecond), 10)
if err := lim.Wait(ctx); err != nil {
    // 限流或取消
}

自定义实现:
type TokenBucket struct {
    capacity int
    tokens   int
    rate     time.Duration
    mu       sync.Mutex
    last     time.Time
}

func NewTokenBucket(capacity int, rate time.Duration) *TokenBucket {
    return &TokenBucket{
        capacity: capacity,
        tokens:   capacity,
        rate:     rate,
        last:     time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    // 补充令牌
    now := time.Now()
    elapsed := now.Sub(tb.last)
    tokensToAdd := int(elapsed / tb.rate)

    if tokensToAdd > 0 {
        tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
        tb.last = now
    }

    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    return false
}

func (tb *TokenBucket) Wait(ctx context.Context) error {
    for {
        if tb.Allow() {
            return nil
        }

        select {
        case <-time.After(tb.rate):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

---

## 八、模式验证

### 8.1 死锁检测

```
死锁条件 (Coffman):
────────────────────────────────────────
1. 互斥: 资源独占
2. 占有并等待: 持有资源同时等待
3. 不可抢占: 资源不能被强制释放
4. 循环等待: 进程间形成循环等待链

Go中预防:
├─ 避免循环等待: 统一加锁顺序
├─ 使用Timeout: select + time.After
├─ 避免持有锁时启动goroutine
└─ 使用context取消

Go 1.26增强:
runtime.SetGoroutineLeakCallback 检测潜在死锁

代码示例:
// 潜在死锁
func potentialDeadlock() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        ch1 <- 1
        ch2 <- 2
    }()

    <-ch2  // 等待
    <-ch1  // 顺序错误，可能死锁
}

// 正确顺序
func correctOrder() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        ch1 <- 1
        ch2 <- 2
    }()

    <-ch1  // 先接收ch1
    <-ch2  // 再接收ch2
}

// 使用select避免死锁
func avoidWithSelect() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        ch1 <- 1
    }()

    go func() {
        ch2 <- 2
    }()

    for i := 0; i < 2; i++ {
        select {
        case v := <-ch1:
            fmt.Println("ch1:", v)
        case v := <-ch2:
            fmt.Println("ch2:", v)
        }
    }
}
```

### 8.2 模式正确性检查表

```
Worker Pool检查:
├─ [ ] Jobs channel关闭后worker退出
├─ [ ] WaitGroup正确等待所有worker
├─ [ ] 错误处理不阻塞worker
└─ [ ] Context取消传播

Pipeline检查:
├─ [ ] 所有阶段关闭output
├─ [ ] 错误时取消或传播
├─ [ ] 无goroutine泄露
└─ [ ] 背压处理

Fan-out/Fan-in检查:
├─ [ ] Fan-out distributor不泄露
├─ [ ] Fan-in等待所有input关闭
├─ [ ] Output正确关闭
└─ [ ] 无序结果可接受
```

---

*本章基于CSP形式化模型，提供了丰富的并发模式代码示例、反例对比和工程实践。*
