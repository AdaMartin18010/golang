# Go并发模式实战手册

> 从经典模式到现代Go并发最佳实践

---

## 一、并发基础模式

### 1.1 Fan-Out/Fan-In

```text
模式说明：
────────────────────────────────────────

Fan-Out：
- 一个输入，多个处理者
- 工作分发到多个goroutine并行处理
- 提高吞吐量

Fan-In：
- 多个输入，一个输出
- 合并多个channel的结果
- 集中处理结果

实际应用：
────────────────────────────────────────

场景：处理大量图片

// 单线程处理太慢
for _, img := range images {
    process(img)  // 串行处理
}

// Fan-Out并行处理
func processImages(images []Image) []Result {
    // Stage 1: Fan-Out
    input := make(chan Image)
    go func() {
        defer close(input)
        for _, img := range images {
            input <- img
        }
    }()

    // 启动多个worker
    numWorkers := runtime.NumCPU()
    results := make([]<-chan Result, numWorkers)

    for i := 0; i < numWorkers; i++ {
        results[i] = worker(input)
    }

    // Stage 2: Fan-In
    return merge(results)
}

func worker(input <-chan Image) <-chan Result {
    output := make(chan Result)
    go func() {
        defer close(output)
        for img := range input {
            output <- process(img)
        }
    }()
    return output
}

func merge(channels []<-chan Result) []Result {
    var wg sync.WaitGroup
    wg.Add(len(channels))

    resultChan := make(chan Result)

    // 从每个channel读取
    for _, ch := range channels {
        go func(c <-chan Result) {
            defer wg.Done()
            for r := range c {
                resultChan <- r
            }
        }(ch)
    }

    // 等待所有channel关闭后关闭resultChan
    go func() {
        wg.Wait()
        close(resultChan)
    }()

    // 收集结果
    var results []Result
    for r := range resultChan {
        results = append(results, r)
    }
    return results
}

性能对比：
────────────────────────────────────────

假设：100张图片，每张处理100ms

串行：100 × 100ms = 10s

Fan-Out (8 workers)：
理论：100 × 100ms / 8 = 1.25s
实际：~1.3s (开销)

加速比：约7.7倍
```

### 1.2 Pipeline模式

```text
流水线概念：
────────────────────────────────────────

数据像流水线一样，经过多个处理阶段：

[Source] → [Stage1] → [Stage2] → [Stage3] → [Sink]

每个阶段：
- 从输入channel读取
- 处理数据
- 发送到输出channel
- 每个阶段可以并行

实现示例：
────────────────────────────────────────

场景：文本处理流水线
1. 读取文件生成行
2. 分割单词
3. 统计词频
4. 输出结果

func generateLines(files []string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for _, file := range files {
            f, _ := os.Open(file)
            scanner := bufio.NewScanner(f)
            for scanner.Scan() {
                out <- scanner.Text()
            }
            f.Close()
        }
    }()
    return out
}

func splitWords(in <-chan string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for line := range in {
            words := strings.Fields(line)
            for _, word := range words {
                out <- strings.ToLower(word)
            }
        }
    }()
    return out
}

func countWords(in <-chan string) <-chan map[string]int {
    out := make(chan map[string]int)
    go func() {
        defer close(out)
        counts := make(map[string]int)
        for word := range in {
            counts[word]++
        }
        out <- counts
    }()
    return out
}

// 组合流水线
func pipeline(files []string) map[string]int {
    c1 := generateLines(files)
    c2 := splitWords(c1)
    c3 := countWords(c2)
    return <-c3
}

优化：并行Stage：
────────────────────────────────────────

func splitWordsParallel(in <-chan string, numWorkers int) <-chan string {
    out := make(chan string)
    var wg sync.WaitGroup
    wg.Add(numWorkers)

    for i := 0; i < numWorkers; i++ {
        go func() {
            defer wg.Done()
            for line := range in {
                words := strings.Fields(line)
                for _, word := range words {
                    out <- strings.ToLower(word)
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

---

## 二、同步模式

### 2.1 使用Context控制

```text
Context的价值：
────────────────────────────────────────

- 超时控制
- 取消信号
- 传递请求范围值
- 控制goroutine生命周期

超时控制：
────────────────────────────────────────

func fetchWithTimeout(url string) (*http.Response, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    return http.DefaultClient.Do(req)
}

取消级联：
────────────────────────────────────────

func processTasks(ctx context.Context, tasks []Task) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    errChan := make(chan error, 1)

    for _, task := range tasks {
        go func(t Task) {
            if err := t.Execute(ctx); err != nil {
                select {
                case errChan <- err:
                    cancel()  // 取消其他任务
                case <-ctx.Done():
                }
            }
        }(task)
    }

    select {
    case err := <-errChan:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

传递元数据：
────────────────────────────────────────

func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateRequestID()
        ctx := context.WithValue(r.Context(), "requestID", requestID)

        log.Printf("[%s] Request started", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
        log.Printf("[%s] Request completed", requestID)
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value("requestID").(string)
    // 使用requestID记录日志
}
```

### 2.2 ErrGroup模式

```text
ErrGroup介绍：
────────────────────────────────────────

golang.org/x/sync/errgroup

- 启动多个goroutine
- 等待全部完成
- 返回第一个错误
- 支持取消

基础使用：
────────────────────────────────────────

import "golang.org/x/sync/errgroup"

func fetchAll(urls []string) error {
    g, _ := errgroup.WithContext(context.Background())

    for _, url := range urls {
        url := url  // 捕获循环变量
        g.Go(func() error {
            resp, err := http.Get(url)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            // 处理响应
            return nil
        })
    }

    return g.Wait()
}

限制并发数：
────────────────────────────────────────

func fetchWithLimit(urls []string, limit int) error {
    g, ctx := errgroup.WithContext(context.Background())
    g.SetLimit(limit)  // 限制并发数

    for _, url := range urls {
        url := url
        g.Go(func() error {
            req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            return nil
        })
    }

    return g.Wait()
}

收集结果：
────────────────────────────────────────

type Result struct {
    URL  string
    Size int
}

func fetchAllWithResults(urls []string) ([]Result, error) {
    g, _ := errgroup.WithContext(context.Background())

    results := make([]Result, len(urls))
    var mu sync.Mutex

    for i, url := range urls {
        i, url := i, url
        g.Go(func() error {
            resp, err := http.Get(url)
            if err != nil {
                return err
            }
            defer resp.Body.Close()

            body, _ := io.ReadAll(resp.Body)

            mu.Lock()
            results[i] = Result{URL: url, Size: len(body)}
            mu.Unlock()

            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

---

## 三、资源管理

### 3.1 连接池模式

```text
为什么需要连接池：
────────────────────────────────────────

资源开销：
- 创建连接：TCP握手、TLS、认证
- 每次请求新建连接：高延迟、资源浪费

连接池优势：
- 复用连接，减少开销
- 限制最大连接数
- 管理连接生命周期

实现示例：
────────────────────────────────────────

type Pool struct {
    factory func() (net.Conn, error)
    maxSize int

    mu      sync.Mutex
    conns   chan net.Conn
    count   int
}

func NewPool(factory func() (net.Conn, error), maxSize int) *Pool {
    return &Pool{
        factory: factory,
        maxSize: maxSize,
        conns:   make(chan net.Conn, maxSize),
    }
}

func (p *Pool) Get() (net.Conn, error) {
    select {
    case conn := <-p.conns:
        return conn, nil
    default:
        p.mu.Lock()
        if p.count < p.maxSize {
            p.count++
            p.mu.Unlock()
            return p.factory()
        }
        p.mu.Unlock()
        // 等待可用连接
        return <-p.conns, nil
    }
}

func (p *Pool) Put(conn net.Conn) {
    if conn == nil {
        return
    }
    select {
    case p.conns <- conn:
    default:
        conn.Close()
        p.mu.Lock()
        p.count--
        p.mu.Unlock()
    }
}

使用：
────────────────────────────────────────

pool := NewPool(func() (net.Conn, error) {
    return net.Dial("tcp", "localhost:8080")
}, 10)

conn, _ := pool.Get()
defer pool.Put(conn)
// 使用连接
```

### 3.2 限流器模式

```text
Token Bucket限流：
────────────────────────────────────────

type TokenBucket struct {
    rate       float64    // 每秒产生token数
    capacity   int        // bucket容量
    tokens     float64    // 当前token数
    lastUpdate time.Time
    mu         sync.Mutex
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
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastUpdate).Seconds()
    tb.tokens += elapsed * tb.rate
    if tb.tokens > float64(tb.capacity) {
        tb.tokens = float64(tb.capacity)
    }
    tb.lastUpdate = now

    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }
    return false
}

使用x/time/rate：
────────────────────────────────────────

import "golang.org/x/time/rate"

// 每秒10个请求，突发20个
limiter := rate.NewLimiter(10, 20)

func handler(w http.ResponseWriter, r *http.Request) {
    if !limiter.Allow() {
        http.Error(w, "Rate limit exceeded", 429)
        return
    }
    // 处理请求
}

等待限流：
────────────────────────────────────────

ctx := context.Background()
err := limiter.Wait(ctx)  // 阻塞直到允许
if err != nil {
    return err
}
// 执行操作
```

---

*本章提供了Go并发模式的实战指南，从基础到高级模式。*
