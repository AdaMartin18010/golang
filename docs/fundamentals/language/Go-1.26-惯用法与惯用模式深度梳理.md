# Go 1.26 惯用法与惯用模式深度梳理

**版本**: Go 1.26
**更新日期**: 2026-03-08
**性质**: 实战语义深度分析

---

## 目录

- [Go 1.26 惯用法与惯用模式深度梳理](#go-126-惯用法与惯用模式深度梳理)
  - [目录](#目录)
  - [1. 类型系统的惯用模式](#1-类型系统的惯用模式)
    - [1.1 零值模式 (Zero Value Pattern)](#11-零值模式-zero-value-pattern)
    - [1.2 类型别名 vs 类型定义的惯用选择](#12-类型别名-vs-类型定义的惯用选择)
    - [1.3 嵌入式字段的惯用组合](#13-嵌入式字段的惯用组合)
  - [2. 控制流的惯用法](#2-控制流的惯用法)
    - [2.1 if 语句的惯用紧凑模式](#21-if-语句的惯用紧凑模式)
    - [2.2 for 循环的惯用变体](#22-for-循环的惯用变体)
    - [2.3 switch 的惯用模式](#23-switch-的惯用模式)
    - [2.4 defer 的惯用模式](#24-defer-的惯用模式)
  - [3. 函数与方法的惯用模式](#3-函数与方法的惯用模式)
    - [3.1 函数选项模式 (Functional Options)](#31-函数选项模式-functional-options)
    - [3.2 方法接收者的选择](#32-方法接收者的选择)
    - [3.3 闭包的惯用模式](#33-闭包的惯用模式)
  - [4. 接口的惯用设计](#4-接口的惯用设计)
    - [4.1 小接口原则](#41-小接口原则)
    - [4.2 接口断言的惯用模式](#42-接口断言的惯用模式)
    - [4.3 空接口的惯用约束](#43-空接口的惯用约束)
  - [5. 并发的惯用模式](#5-并发的惯用模式)
    - [5.1 Pipeline 模式](#51-pipeline-模式)
    - [5.2 Worker Pool 模式](#52-worker-pool-模式)
    - [5.3 Context 的惯用模式](#53-context-的惯用模式)
    - [5.4 Select 的高级惯用模式](#54-select-的高级惯用模式)
  - [6. 错误处理的惯用法](#6-错误处理的惯用法)
    - [6.1 错误包装链](#61-错误包装链)
    - [6.2 Sentinel 错误模式](#62-sentinel-错误模式)
    - [6.3 自定义错误类型](#63-自定义错误类型)
  - [7. 内存管理的惯用模式](#7-内存管理的惯用模式)
    - [7.1 Sync.Pool 的对象复用](#71-syncpool-的对象复用)
    - [7.2 切片预分配](#72-切片预分配)
  - [8. 反射的惯用场景](#8-反射的惯用场景)
    - [8.1 结构体标签解析](#81-结构体标签解析)
    - [8.2 类型注册模式](#82-类型注册模式)
  - [9. 泛型的惯用模式](#9-泛型的惯用模式)
    - [9.1 约束设计](#91-约束设计)
    - [9.2 泛型数据结构](#92-泛型数据结构)
    - [9.3 泛型函数组合](#93-泛型函数组合)
  - [总结](#总结)
    - [核心惯用原则](#核心惯用原则)
    - [避免的反模式](#避免的反模式)

---

## 1. 类型系统的惯用模式

### 1.1 零值模式 (Zero Value Pattern)

**语义**: Go的零值是类型安全的默认值，利用零值可以简化初始化逻辑。

```go
// 惯用法: 利用零值避免显式初始化
var (
    m   map[string]int      // nil map, 读安全，写需make
    s   []int               // nil slice, len=0, cap=0
    ch  chan int            // nil channel, 阻塞操作
    p   *User               // nil pointer
    i   interface{}         // nil interface
    fn  func()              // nil function
)

// 实战模式: 延迟初始化
func (c *Cache) Get(key string) (Value, bool) {
    // 零值检查避免nil panic
    if c == nil {
        return nil, false
    }
    // 延迟初始化
    if c.data == nil {
        c.data = make(map[string]Value)
    }
    return c.data[key]
}
```

**惯用场景**:

- 配置结构体的可选字段
- 延迟初始化的缓存
- 函数选项模式中的默认值

### 1.2 类型别名 vs 类型定义的惯用选择

```go
// 类型定义 - 创建新类型 (推荐用于领域建模)
type UserID int64
type Email string

// 方法绑定强化类型语义
func (id UserID) Valid() bool {
    return id > 0
}

func (e Email) Domain() string {
    parts := strings.Split(string(e), "@")
    if len(parts) != 2 {
        return ""
    }
    return parts[1]
}

// 类型别名 - 兼容现有代码 (仅用于重构/迁移)
type OldAPI = NewAPI  // 完全等价，仅在迁移期使用
```

**决策树**:

- 需要方法？→ 类型定义
- 需要不同零值语义？→ 类型定义
- 仅仅是代码迁移？→ 类型别名

### 1.3 嵌入式字段的惯用组合

```go
// 惯用模式: 组合优于继承
type Server struct {
    http.Server                    // 匿名嵌入: 方法提升
    logger    *zap.Logger          // 命名嵌入: 显式组合
    config    Config               // 值嵌入: 避免nil检查
}

// 初始化惯用法
func NewServer(cfg Config) *Server {
    return &Server{
        Server: http.Server{
            Addr:    cfg.Addr,
            Handler: cfg.Handler,
        },
        logger: cfg.Logger,
        config: cfg,  // 值复制，线程安全
    }
}

// 方法重写惯用法
func (s *Server) ListenAndServe() error {
    s.logger.Info("starting server", zap.String("addr", s.Addr))
    return s.Server.ListenAndServe()  // 显式调用嵌入方法
}
```

**陷阱避免**:

```go
// ❌ 错误: 嵌入指针
 type Bad struct {
     *http.Server  // 可能导致nil panic
 }

// ✅ 正确: 嵌入值或使用命名嵌入
 type Good struct {
     server http.Server  // 值嵌入，零值安全
 }
```

---

## 2. 控制流的惯用法

### 2.1 if 语句的惯用紧凑模式

```go
// 惯用法: 简短语句 + 条件
if err := doSomething(); err != nil {
    return fmt.Errorf("failed: %w", err)
}
// err 作用域到此结束

// 惯用法: 错误值预声明模式
var err error
if resp, err = client.Do(req); err != nil {
    return err
}
// resp 在后续代码中可用

// 惯用法: 多重条件短路
if user != nil && user.Profile != nil && user.Profile.Email != "" {
    // 安全检查链
}

// Go 1.26: 泛型辅助函数避免重复
func If[T any](cond bool, a, b T) T {
    if cond {
        return a
    }
    return b
}

// 使用
max := If(a > b, a, b)
```

### 2.2 for 循环的惯用变体

```go
// 惯用法 1: 无限循环 + break
for {
    item, ok := queue.Dequeue()
    if !ok {
        break
    }
    process(item)
}

// 惯用法 2: 条件循环
for hasMore() {
    processNext()
}

// 惯用法 3: range + 索引/值选择
for i := range slice {        // 只需要索引
    slice[i] *= 2
}

for _, v := range slice {     // 只需要值
    sum += v
}

for i, v := range slice {     // 两者都需要
    result[i] = transform(v)
}

// 惯用法 4: map range + 删除
for k := range m {
    if shouldDelete(k) {
        delete(m, k)  // 删除是安全的
    }
}

// 惯用法 5: 字符串 range (rune 遍历)
for i, r := range "Hello, 世界" {
    // i: 字节索引, r: rune (Unicode码点)
    fmt.Printf("%d: %c\n", i, r)
}
```

### 2.3 switch 的惯用模式

```go
// 惯用法 1: 类型断言 switch
func handleEvent(e Event) {
    switch v := e.(type) {
    case *UserCreated:
        handleUserCreated(v)
    case *OrderPlaced:
        handleOrderPlaced(v)
    case nil:
        log.Println("nil event")
    default:
        log.Printf("unknown event: %T", v)
    }
}

// 惯用法 2: 表达式 switch (替代 if-else 链)
func grade(score int) string {
    switch {
    case score >= 90:
        return "A"
    case score >= 80:
        return "B"
    case score >= 70:
        return "C"
    case score >= 60:
        return "D"
    default:
        return "F"
    }
}

// 惯用法 3: 枚举模式 + switch
 type Status int
 const (
     StatusPending Status = iota
     StatusProcessing
     StatusCompleted
     StatusFailed
 )

func (s Status) String() string {
    switch s {
    case StatusPending:
        return "pending"
    case StatusProcessing:
        return "processing"
    case StatusCompleted:
        return "completed"
    case StatusFailed:
        return "failed"
    default:
        return fmt.Sprintf("Status(%d)", s)
    }
}
```

### 2.4 defer 的惯用模式

```go
// 惯用法 1: 资源清理栈 (LIFO 顺序)
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // 第一个defer，最后执行

    r, err := gzip.NewReader(f)
    if err != nil {
        return err
    }
    defer r.Close()  // 第二个defer，先执行

    return process(r)
}

// 惯用法 2: 参数预计算
func trace(msg string) func() {
    start := time.Now()
    log.Printf("enter %s", msg)
    return func() {
        log.Printf("exit %s (%s)", msg, time.Since(start))
    }
}

func example() {
    defer trace("example")()  // trace立即执行，返回的函数延迟执行
    // do work...
}

// 惯用法 3: 错误修改
func doSomething() (err error) {
    f, err := os.Create("file.txt")
    if err != nil {
        return err
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = cerr  // 修改命名返回值
        }
    }()
    return writeData(f)
}

// 惯用法 4: panic 恢复
func safeCall(fn func()) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("recovered: %v", r)
        }
    }()
    fn()
}
```

---

## 3. 函数与方法的惯用模式

### 3.1 函数选项模式 (Functional Options)

```go
// 惯用法: 可扩展的配置选项
 type Server struct {
     addr     string
     timeout  time.Duration
     tls      *tls.Config
     logger   *zap.Logger
 }

 type Option func(*Server)

 func WithAddress(addr string) Option {
     return func(s *Server) {
         s.addr = addr
     }
 }

 func WithTimeout(d time.Duration) Option {
     return func(s *Server) {
         s.timeout = d
     }
 }

 func WithTLS(cfg *tls.Config) Option {
     return func(s *Server) {
         s.tls = cfg
     }
 }

 func WithLogger(l *zap.Logger) Option {
     return func(s *Server) {
         s.logger = l
     }
 }

 // 构造函数
 func NewServer(opts ...Option) *Server {
     s := &Server{
         addr:    ":8080",
         timeout: 30 * time.Second,
         logger:  zap.NewNop(),
     }
     for _, opt := range opts {
         opt(s)
     }
     return s
 }

 // 使用
 srv := NewServer(
     WithAddress(":9090"),
     WithTimeout(60*time.Second),
     WithLogger(logger),
 )
```

**优势**:

- 向后兼容添加选项
- 自文档化 API
- 默认值支持
- 可组合

### 3.2 方法接收者的选择

```go
// 值接收者: 不可变语义
 type Point struct{ X, Y float64 }

 func (p Point) Distance(q Point) float64 {
     return math.Hypot(q.X-p.X, q.Y-p.Y)
 }

// 指针接收者: 可变语义或大型结构
 type Buffer struct {
     data []byte
     off  int
 }

 func (b *Buffer) Write(p []byte) (n int, err error) {
     b.data = append(b.data, p...)
     return len(p), nil
 }

// 混合模式: 读用值，写用指针
 type Counter struct {
     count int64
 }

 func (c Counter) Value() int64 {    // 值接收: 读操作
     return c.count
 }

 func (c *Counter) Inc() {           // 指针接收: 写操作
     atomic.AddInt64(&c.count, 1)
 }
```

**决策规则**:

- 需要修改接收者？→ 指针
- 包含 mutex/sync 字段？→ 指针
- 大型结构 (>100 bytes)？→ 指针
- 需要 nil 接收者？→ 指针
- 否则 → 值

### 3.3 闭包的惯用模式

```go
// 惯用法 1: 工厂函数
 func makeMultiplier(factor int) func(int) int {
     return func(x int) int {
         return x * factor  // 捕获 factor
     }
 }

 double := makeMultiplier(2)
 triple := makeMultiplier(3)

// 惯用法 2: 中间件/装饰器
 func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
     return func(w http.ResponseWriter, r *http.Request) {
         start := time.Now()
         next(w, r)
         log.Printf("%s %s %s", r.Method, r.URL, time.Since(start))
     }
 }

// 惯用法 3: 状态保持
 func NewLimiter(rate int) func() bool {
     tokens := rate
     last := time.Now()

     return func() bool {
         now := time.Now()
         elapsed := now.Sub(last)
         last = now

         tokens += int(elapsed.Seconds()) * rate
         if tokens > rate {
             tokens = rate
         }

         if tokens > 0 {
             tokens--
             return true
         }
         return false
     }
 }
```

---

## 4. 接口的惯用设计

### 4.1 小接口原则

```go
// 惯用法: 小接口，大能力
 type Reader interface {
     Read(p []byte) (n int, err error)
 }

 type Writer interface {
     Write(p []byte) (n int, err error)
 }

 type Closer interface {
     Close() error
 }

 // 通过组合小接口构建大接口
 type ReadWriter interface {
     Reader
     Writer
 }

 type ReadWriteCloser interface {
     Reader
     Writer
     Closer
 }
```

### 4.2 接口断言的惯用模式

```go
// 惯用法 1: 类型断言 + ok 模式
func process(w io.Writer) {
    if f, ok := w.(*os.File); ok {
        // 是文件，可以获取文件名
        fmt.Println("writing to file:", f.Name())
    }
    // 否则继续使用 Writer 接口
}

// 惯用法 2: 接口升级模式
func optimizedWrite(w io.Writer, p []byte) (int, error) {
    // 检查是否有更高效的实现
    if rw, ok := w.(io.ReaderFrom); ok {
        return rw.ReadFrom(bytes.NewReader(p))
    }
    return w.Write(p)
}

// 惯用法 3: 能力检测
func closeIfPossible(c interface{}) error {
    if closer, ok := c.(io.Closer); ok {
        return closer.Close()
    }
    return nil
}
```

### 4.3 空接口的惯用约束

```go
// 避免空接口，使用约束
 func ProcessAny(v interface{}) {  // ❌ 避免
     // 无类型信息
 }

// 使用类型约束
 func Process[T any](v T) {  // ✅ 更好
     // 有类型信息，编译器可优化
 }

// 使用具体接口约束
 func Stringify(w io.Writer, v interface{}) {  // ✅ 特定场景
     fmt.Fprint(w, v)
 }
```

---

## 5. 并发的惯用模式

### 5.1 Pipeline 模式

```go
// 惯用法: 扇出-扇入流水线
 func generator(nums ...int) <-chan int {
     out := make(chan int)
     go func() {
         defer close(out)
         for _, n := range nums {
             out <- n
         }
     }()
     return out
 }

 func square(in <-chan int) <-chan int {
     out := make(chan int)
     go func() {
         defer close(out)
         for n := range in {
             out <- n * n
         }
     }()
     return out
 }

 // 扇出: 多个 goroutine 处理
 func merge(cs ...<-chan int) <-chan int {
     var wg sync.WaitGroup
     out := make(chan int)

     output := func(c <-chan int) {
         defer wg.Done()
         for n := range c {
             out <- n
         }
     }

     wg.Add(len(cs))
     for _, c := range cs {
         go output(c)
     }

     go func() {
         wg.Wait()
         close(out)
     }()

     return out
 }

 // 使用
 in := generator(1, 2, 3, 4, 5)

 // 扇出: 两个 square 处理器
 c1 := square(in)
 c2 := square(in)

 // 扇入: 合并结果
 for n := range merge(c1, c2) {
     fmt.Println(n)
 }
```

### 5.2 Worker Pool 模式

```go
// 惯用法: 固定大小的 worker pool
 type Pool struct {
     workers int
     jobs    chan Job
     results chan Result
     wg      sync.WaitGroup
 }

 func NewPool(workers int) *Pool {
     return &Pool{
         workers: workers,
         jobs:    make(chan Job),
         results: make(chan Result, workers),
     }
 }

 func (p *Pool) Start(ctx context.Context) {
     for i := 0; i < p.workers; i++ {
         p.wg.Add(1)
         go p.worker(ctx, i)
     }
 }

 func (p *Pool) worker(ctx context.Context, id int) {
     defer p.wg.Done()
     for {
         select {
         case job, ok := <-p.jobs:
             if !ok {
                 return
             }
             p.results <- job.Process()
         case <-ctx.Done():
             return
         }
     }
 }

 func (p *Pool) Submit(j Job) { p.jobs <- j }
 func (p *Pool) Results() <-chan Result { return p.results }

 func (p *Pool) Stop() {
     close(p.jobs)
     p.wg.Wait()
     close(p.results)
 }
```

### 5.3 Context 的惯用模式

```go
// 惯用法 1: 超时控制
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()

 req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
 resp, err := client.Do(req)

// 惯用法 2: 值传递 (仅用于请求上下文)
 type contextKey string
 const userKey contextKey = "user"

 func WithUser(ctx context.Context, user *User) context.Context {
     return context.WithValue(ctx, userKey, user)
 }

 func UserFromContext(ctx context.Context) (*User, bool) {
     user, ok := ctx.Value(userKey).(*User)
     return user, ok
 }

// 惯用法 3: 取消传播
 func parent(ctx context.Context) {
     ctx, cancel := context.WithCancel(ctx)
     defer cancel()

     go child(ctx)

     // 条件触发取消
     if shouldStop() {
         cancel()  // 所有子 goroutine 收到取消信号
     }
 }

 func child(ctx context.Context) {
     for {
         select {
         case <-ctx.Done():
             return
         default:
             // do work
         }
     }
 }
```

### 5.4 Select 的高级惯用模式

```go
// 惯用法 1: 默认分支防止阻塞
 select {
 case v := <-ch:
     process(v)
 default:
     // 立即执行，不阻塞
 }

// 惯用法 2: 超时等待
 select {
 case v := <-ch:
     process(v)
 case <-time.After(5 * time.Second):
     // 超时处理
 }

// 惯用法 3: 随机选择多个就绪 channel
 for {
     select {
     case a := <-chA:
         processA(a)
     case b := <-chB:
         processB(b)
     case c := <-chC:
         processC(c)
     }
 }

// 惯用法 4: 优雅关闭
 func (s *Server) Run() {
     for {
         select {
         case req := <-s.requests:
             s.handle(req)
         case <-s.quit:
             // 优雅关闭
             s.shutdown()
             return
         }
     }
 }
```

---

## 6. 错误处理的惯用法

### 6.1 错误包装链

```go
// 惯用法: 构建错误上下文链
 func readConfig(path string) (*Config, error) {
     data, err := os.ReadFile(path)
     if err != nil {
         return nil, fmt.Errorf("reading config %q: %w", path, err)
     }

     var cfg Config
     if err := json.Unmarshal(data, &cfg); err != nil {
         return nil, fmt.Errorf("parsing config %q: %w", path, err)
     }

     if err := cfg.Validate(); err != nil {
         return nil, fmt.Errorf("invalid config %q: %w", path, err)
     }

     return &cfg, nil
 }

// 错误检查惯用法
 if err := readConfig("config.json"); err != nil {
     // 检查原始错误
     if errors.Is(err, os.ErrNotExist) {
         // 配置文件不存在
     }

     // 检查特定错误类型
     var cfgErr *ConfigError
     if errors.As(err, &cfgErr) {
         // 配置错误
     }
 }
```

### 6.2 Sentinel 错误模式

```go
// 惯用法: 预定义错误变量
 var (
     ErrNotFound     = errors.New("not found")
     ErrAlreadyExist = errors.New("already exists")
     ErrUnauthorized = errors.New("unauthorized")
 )

 type Store struct {
     data map[string]Value
 }

 func (s *Store) Get(key string) (Value, error) {
     v, ok := s.data[key]
     if !ok {
         return nil, fmt.Errorf("key %q: %w", key, ErrNotFound)
     }
     return v, nil
 }

// 错误检查
 if err := store.Get("key"); err != nil {
     if errors.Is(err, ErrNotFound) {
         // 处理不存在
     }
 }
```

### 6.3 自定义错误类型

```go
// 惯用法: 结构化错误信息
 type ValidationError struct {
     Field   string
     Value   interface{}
     Message string
 }

 func (e *ValidationError) Error() string {
     return fmt.Sprintf("validation failed for %s: %v - %s",
         e.Field, e.Value, e.Message)
 }

 func (e *ValidationError) Unwrap() error {
     // 支持错误链
     return ErrInvalidInput
 }

// 使用
 if age < 0 {
     return &ValidationError{
         Field:   "age",
         Value:   age,
         Message: "must be non-negative",
     }
 }
```

---

## 7. 内存管理的惯用模式

### 7.1 Sync.Pool 的对象复用

```go
// 惯用法: 高频率分配的对象池
 type BufferPool struct {
     pool sync.Pool
 }

 func NewBufferPool(size int) *BufferPool {
     return &BufferPool{
         pool: sync.Pool{
             New: func() interface{} {
                 return make([]byte, 0, size)
             },
         },
     }
 }

 func (p *BufferPool) Get() []byte {
     return p.pool.Get().([]byte)
 }

 func (p *BufferPool) Put(b []byte) {
     // 重置但保留容量
     b = b[:0]
     p.pool.Put(b)
 }

// 使用
 var pool = NewBufferPool(4096)

 func process(w io.Writer, data []byte) {
     buf := pool.Get()
     defer pool.Put(buf)

     // 使用 buf 处理数据
     buf = append(buf, data...)
     w.Write(buf)
 }
```

### 7.2 切片预分配

```go
// 惯用法: 已知大小时预分配
 func filter(items []Item, predicate func(Item) bool) []Item {
     // 最坏情况: 全部匹配
     result := make([]Item, 0, len(items))

     for _, item := range items {
         if predicate(item) {
             result = append(result, item)
         }
     }

     return result
 }

// 惯用法: map 预分配
 func groupBy(items []Item) map[string][]Item {
     // 预估大小
     groups := make(map[string][]Item, len(items))

     for _, item := range items {
         key := item.Category
         groups[key] = append(groups[key], item)
     }

     return groups
 }
```

---

## 8. 反射的惯用场景

### 8.1 结构体标签解析

```go
// 惯用法: 自定义标签解析
 type Config struct {
     Host string `env:"HOST" default:"localhost"`
     Port int    `env:"PORT" default:"8080"`
 }

 func LoadFromEnv(cfg interface{}) error {
     v := reflect.ValueOf(cfg).Elem()
     t := v.Type()

     for i := 0; i < t.NumField(); i++ {
         field := t.Field(i)
         envKey := field.Tag.Get("env")
         defaultVal := field.Tag.Get("default")

         val := os.Getenv(envKey)
         if val == "" {
             val = defaultVal
         }

         // 设置字段值
         // ...
     }
     return nil
 }
```

### 8.2 类型注册模式

```go
// 惯用法: 类型工厂
 var types = make(map[string]reflect.Type)

 func Register(name string, prototype interface{}) {
     types[name] = reflect.TypeOf(prototype).Elem()
 }

 func New(name string) (interface{}, error) {
     t, ok := types[name]
     if !ok {
         return nil, fmt.Errorf("unknown type: %s", name)
     }
     return reflect.New(t).Interface(), nil
 }

// 使用
 type User struct{}

 func init() {
     Register("user", (*User)(nil))
 }
```

---

## 9. 泛型的惯用模式

### 9.1 约束设计

```go
// 惯用法: 最小约束原则
 // 只需要 String 方法
 type Stringer interface {
     String() string
 }

 func ToString[T Stringer](v T) string {
     return v.String()
 }

// 可比较类型
 func Contains[T comparable](s []T, v T) bool {
     for _, vs := range s {
         if v == vs {
             return true
         }
     }
     return false
 }

// 数值类型约束
 type Number interface {
     ~int | ~int8 | ~int16 | ~int32 | ~int64 |
     ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
     ~float32 | ~float64
 }

 func Sum[T Number](vals []T) T {
     var sum T
     for _, v := range vals {
         sum += v
     }
     return sum
 }
```

### 9.2 泛型数据结构

```go
// 惯用法: 类型安全的栈
type Stack[T any] struct {
    items []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{}
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    v := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

// 使用
intStack := NewStack[int]()
intStack.Push(1)
intStack.Push(2)
v, _ := intStack.Pop()  // v = 2, 类型为 int

strStack := NewStack[string]()
strStack.Push("hello")
```

### 9.3 泛型函数组合

```go
// 惯用法: Map/Filter/Reduce
func Map[T, R any](s []T, fn func(T) R) []R {
    result := make([]R, len(s))
    for i, v := range s {
        result[i] = fn(v)
    }
    return result
}

func Filter[T any](s []T, fn func(T) bool) []T {
    result := make([]T, 0, len(s))
    for _, v := range s {
        if fn(v) {
            result = append(result, v)
        }
    }
    return result
}

func Reduce[T, R any](s []T, init R, fn func(R, T) R) R {
    result := init
    for _, v := range s {
        result = fn(result, v)
    }
    return result
}

// 使用
numbers := []int{1, 2, 3, 4, 5}

// 映射
squares := Map(numbers, func(n int) int { return n * n })

// 过滤
evens := Filter(numbers, func(n int) bool { return n%2 == 0 })

// 归约
sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n })
```

---

## 总结

### 核心惯用原则

1. **利用零值**: 减少显式初始化，提高代码简洁性
2. **小接口**: 接口越小，组合越灵活
3. **组合优于继承**: 通过嵌入实现代码复用
4. **显式优于隐式**: 错误处理、类型转换都要显式
5. **并发安全**: 共享数据必须通过channel或sync
6. **最小约束**: 泛型约束越小，适用性越广

### 避免的反模式

- ❌ 空接口滥用
- ❌ 过度使用反射
- ❌ 忽略错误返回值
- ❌ 全局可变状态
- ❌ 复杂的嵌套控制流

---

**文档版本**: 1.0
**最后更新**: 2026-03-08
