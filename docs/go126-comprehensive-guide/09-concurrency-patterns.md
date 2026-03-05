# 第九章：并发与并行模式

> Go 并发编程的高级模式和最佳实践

---

## 9.1 基础并发模式

### 9.1.1 Generator 模式

```go
// 返回只读 channel 的函数模式
func Generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// 使用
for n := range Generator(1, 2, 3, 4, 5) {
    fmt.Println(n)
}

// 带取消的 Generator
func GeneratorCtx(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

### 9.1.2 Pipeline 模式

```go
// 管道模式：数据流处理

// Stage 1: 生成数据
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Stage 2: 平方
func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Stage 3: 过滤奇数
func filterOdd(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

// 组合管道
func main() {
    c := gen(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    out := filterOdd(sq(c))

    for n := range out {
        fmt.Println(n) // 4, 16, 36, 64, 100
    }
}
```

### 9.1.3 Fan-Out / Fan-In 模式

```go
// 扇出：多个 goroutine 读取同一输入
func fanOut(in <-chan int, n int) []<-chan int {
    outs := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        out := make(chan int)
        outs[i] = out
        go func() {
            defer close(out)
            for v := range in {
                out <- process(v)
            }
        }()
    }
    return outs
}

// 扇入：合并多个 channel
func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func(c <-chan int) {
        defer wg.Done()
        for v := range c {
            out <- v
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

// 使用
func main() {
    in := gen(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

    // 扇出到 3 个处理器
    outs := fanOut(in, 3)

    // 扇入合并结果
    for result := range fanIn(outs...) {
        fmt.Println(result)
    }
}
```

---

## 9.2 同步模式

### 9.2.1 Worker Pool 模式

```go
type Task func()

// WorkerPool 管理一组工作 goroutine
type WorkerPool struct {
    tasks   chan Task
    wg      sync.WaitGroup
    size    int
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewWorkerPool(size int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        tasks:  make(chan Task),
        size:   size,
        ctx:    ctx,
        cancel: cancel,
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.size; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()
    for {
        select {
        case task, ok := <-p.tasks:
            if !ok {
                return
            }
            task()
        case <-p.ctx.Done():
            return
        }
    }
}

func (p *WorkerPool) Submit(task Task) {
    select {
    case p.tasks <- task:
    case <-p.ctx.Done():
    }
}

func (p *WorkerPool) Stop() {
    p.cancel()
    close(p.tasks)
    p.wg.Wait()
}

// 使用
func main() {
    pool := NewWorkerPool(4)
    pool.Start()

    for i := 0; i < 100; i++ {
        n := i
        pool.Submit(func() {
            fmt.Printf("Processing task %d\n", n)
            time.Sleep(100 * time.Millisecond)
        })
    }

    pool.Stop()
}
```

### 9.2.2 信号量模式

```go
// 带权重的信号量
type Semaphore struct {
    sem chan struct{}
}

func NewSemaphore(n int) *Semaphore {
    return &Semaphore{
        sem: make(chan struct{}, n),
    }
}

func (s *Semaphore) Acquire(ctx context.Context) error {
    select {
    case s.sem <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (s *Semaphore) Release() {
    select {
    case <-s.sem:
    default:
        panic("semaphore: release without acquire")
    }
}

// 使用：限制并发数
func processItems(items []Item) {
    sem := NewSemaphore(10) // 最多 10 个并发
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        go func(it Item) {
            defer wg.Done()

            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()

            if err := sem.Acquire(ctx); err != nil {
                log.Printf("Failed to acquire semaphore: %v", err)
                return
            }
            defer sem.Release()

            process(it)
        }(item)
    }

    wg.Wait()
}
```

### 9.2.3 屏障模式 (Barrier)

```go
// 等待所有 goroutine 到达某个点
type Barrier struct {
    n      int
    count  int
    mutex  sync.Mutex
    cond   *sync.Cond
}

func NewBarrier(n int) *Barrier {
    b := &Barrier{n: n}
    b.cond = sync.NewCond(&b.mutex)
    return b
}

func (b *Barrier) Wait() {
    b.mutex.Lock()
    b.count++
    if b.count == b.n {
        b.count = 0
        b.cond.Broadcast()
    } else {
        b.cond.Wait()
    }
    b.mutex.Unlock()
}

// 使用
func main() {
    barrier := NewBarrier(3)

    for i := 0; i < 3; i++ {
        go func(id int) {
            fmt.Printf("Worker %d: Phase 1\n", id)
            time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
            barrier.Wait()

            fmt.Printf("Worker %d: Phase 2\n", id)
            time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
            barrier.Wait()

            fmt.Printf("Worker %d: Phase 3\n", id)
        }(i)
    }

    time.Sleep(2 * time.Second)
}
```

---

## 9.3 超时与取消模式

### 9.3.1 Timeout 模式

```go
// 函数超时包装器
func WithTimeout(fn func() error, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    errChan := make(chan error, 1)
    go func() {
        errChan <- fn()
    }()

    select {
    case err := <-errChan:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

// 带超时的 HTTP 请求
func fetchWithTimeout(url string, timeout time.Duration) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

// 慢操作超时
func slowOperation(ctx context.Context) (Result, error) {
    resultChan := make(chan Result, 1)

    go func() {
        resultChan <- doSlowWork()
    }()

    select {
    case result := <-resultChan:
        return result, nil
    case <-ctx.Done():
        return Result{}, ctx.Err()
    }
}
```

### 9.3.2 优雅关闭模式

```go
// GracefulShutdown 等待所有任务完成或超时
func GracefulShutdown(shutdownFunc func(context.Context) error, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    done := make(chan struct{})
    var err error

    go func() {
        err = shutdownFunc(ctx)
        close(done)
    }()

    select {
    case <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

// HTTP 服务器优雅关闭
func runServer() {
    srv := &http.Server{Addr: ":8080"}

    // 启动服务器
    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited")
}
```

---

## 9.4 高级并发模式

### 9.4.1 Circuit Breaker 熔断器

```go
type State int

const (
    StateClosed State = iota    // 正常
    StateOpen                   // 熔断
    StateHalfOpen               // 半开
)

type CircuitBreaker struct {
    failureThreshold int
    successThreshold int
    timeout          time.Duration

    state        State
    failures     int
    successes    int
    lastFailure  time.Time
    mutex        sync.Mutex
}

func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failureThreshold: failureThreshold,
        successThreshold: successThreshold,
        timeout:          timeout,
        state:            StateClosed,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.canExecute() {
        return errors.New("circuit breaker is open")
    }

    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
            cb.failures = 0
            cb.successes = 0
            return true
        }
        return false
    case StateHalfOpen:
        return true
    }
    return false
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

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
        cb.successes++
        if cb.successes >= cb.successThreshold {
            cb.state = StateClosed
            cb.failures = 0
            cb.successes = 0
        }
    }
}

func (cb *CircuitBreaker) onFailure() {
    cb.failures++
    cb.lastFailure = time.Now()

    switch cb.state {
    case StateClosed:
        if cb.failures >= cb.failureThreshold {
            cb.state = StateOpen
        }
    case StateHalfOpen:
        cb.state = StateOpen
    }
}
```

### 9.4.2 Rate Limiter 限流器

```go
// Token Bucket 限流器
type TokenBucket struct {
    rate       float64    // 每秒产生令牌数
    capacity   int        // 桶容量
    tokens     float64    // 当前令牌数
    lastUpdate time.Time
    mutex      sync.Mutex
}

func NewTokenBucket(rate float64, capacity int) *TokenBucket {
    return &TokenBucket{
        rate:       rate,
        capacity:   capacity,
        tokens:     float64(capacity),
        lastUpdate: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    return tb.AllowN(1)
}

func (tb *TokenBucket) AllowN(n int) bool {
    tb.mutex.Lock()
    defer tb.mutex.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastUpdate).Seconds()
    tb.lastUpdate = now

    // 添加新令牌
    tb.tokens += elapsed * tb.rate
    if tb.tokens > float64(tb.capacity) {
        tb.tokens = float64(tb.capacity)
    }

    // 检查是否足够
    if tb.tokens >= float64(n) {
        tb.tokens -= float64(n)
        return true
    }
    return false
}

// 等待获取令牌
func (tb *TokenBucket) Wait(ctx context.Context) error {
    for {
        if tb.Allow() {
            return nil
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(time.Millisecond * 10):
            // 继续尝试
        }
    }
}
```

### 9.4.3 SingleFlight 防止缓存击穿

```go
type call struct {
    wg  sync.WaitGroup
    val interface{}
    err error
}

type Group struct {
    mutex sync.Mutex
    m     map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
    g.mutex.Lock()
    if g.m == nil {
        g.m = make(map[string]*call)
    }
    if c, ok := g.m[key]; ok {
        g.mutex.Unlock()
        c.wg.Wait()
        return c.val, c.err
    }
    c := new(call)
    c.wg.Add(1)
    g.m[key] = c
    g.mutex.Unlock()

    c.val, c.err = fn()
    c.wg.Done()

    g.mutex.Lock()
    delete(g.m, key)
    g.mutex.Unlock()

    return c.val, c.err
}

// 使用
type Cache struct {
    data  map[string]string
    mutex sync.RWMutex
    group singleflight.Group
}

func (c *Cache) Get(key string) (string, error) {
    c.mutex.RLock()
    if val, ok := c.data[key]; ok {
        c.mutex.RUnlock()
        return val, nil
    }
    c.mutex.RUnlock()

    // 防止缓存击穿
    val, err := c.group.Do(key, func() (interface{}, error) {
        // 从数据库加载
        data, err := loadFromDB(key)
        if err != nil {
            return nil, err
        }

        c.mutex.Lock()
        c.data[key] = data
        c.mutex.Unlock()

        return data, nil
    })

    if err != nil {
        return "", err
    }
    return val.(string), nil
}
```

---

## 9.5 上下文传播模式

```go
// Context 值传播的最佳实践

// 使用私有类型作为 key，避免冲突
type contextKey string

const userIDKey contextKey = "userID"
const requestIDKey contextKey = "requestID"

// WithUserID 添加上下文
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext 提取
func UserIDFromContext(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}

// 中间件中使用
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := validateToken(r.Header.Get("Authorization"))
        ctx := WithUserID(r.Context(), userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Handler 中使用
func GetProfile(w http.ResponseWriter, r *http.Request) {
    userID, ok := UserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    // 使用 userID...
}
```

---

## 9.6 模式选择决策树

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    并发模式选择决策树                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  需要保护共享状态？                                                          │
│       │                                                                     │
│       ├── Yes ──▶ 访问频率？                                                │
│       │              │                                                      │
│       │              ├── 高频读/低频写 ──▶ RWMutex                         │
│       │              │                                                      │
│       │              ├── 高频写 ──▶ Mutex                                  │
│       │              │                                                      │
│       │              └── 计数/标志 ──▶ atomic                              │
│       │                                                                     │
│       └── No ──▶ 需要协调 goroutine？                                       │
│                     │                                                       │
│                     ├── Yes ──▶ 协调方式？                                  │
│                     │              │                                        │
│                     │              ├── 数据传递 ──▶ Channel                 │
│                     │              │                                        │
│                     │              ├── 一对多通知 ──▶ sync.Cond             │
│                     │              │                                        │
│                     │              └── 等待组完成 ──▶ WaitGroup             │
│                     │                                                       │
│                     └── No ──▶ 使用独立状态（每个 goroutine 一份）           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*本章涵盖了 Go 并发编程的核心模式和最佳实践，是构建高性能、可靠并发系统的基础。*
