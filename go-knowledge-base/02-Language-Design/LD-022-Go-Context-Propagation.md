# LD-022: Go 上下文传播机制 (Go Context Propagation)

> **维度**: Language Design
> **级别**: S (17+ KB)
> **标签**: #context #cancellation #timeout #deadline #propagation #request-scoped
> **权威来源**:
>
> - [context Package](https://github.com/golang/go/tree/master/src/context) - Go Authors
> - [Go Concurrency Patterns: Context](https://go.dev/blog/context) - Sameer Ajmani
> - [Context Best Practices](https://rakyll.org/context/) - rakyll

---

## 1. Context 设计原理

### 1.1 核心概念

```
┌─────────────────────────────────────────────────────────────┐
│                      Context Tree                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│                         root                                 │
│                          │                                   │
│                    background()                              │
│                          │                                   │
│             ┌────────────┼────────────┐                     │
│             │            │            │                     │
│             ▼            ▼            ▼                     │
│         ctx1          ctx2         ctx3                     │
│       (timeout)    (cancel)     (values)                    │
│             │            │            │                     │
│       ┌─────┘            │            ├─────┐               │
│       │                  │            │     │               │
│       ▼                  ▼            ▼     ▼               │
│     ctx4               ctx5        ctx6  ctx7              │
│   (value)           (deadline)                                │
│                                                              │
│  特性:                                                        │
│  - 树形结构，父节点取消传播到子节点                              │
│  - 不可变，派生创建新 Context                                   │
│  - 线程安全，可被多个 goroutine 同时访问                         │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 接口定义

```go
// src/context/context.go

// Context 接口
type Context interface {
    // Deadline 返回超时时间，ok=false 表示没有设置
    Deadline() (deadline time.Time, ok bool)
    
    // Done 返回一个 channel，当 Context 被取消或超时时关闭
    Done() <-chan struct{}
    
    // Err 返回 Context 结束的原因
    Err() error
    
    // Value 根据 key 获取存储的值
    Value(key any) any
}

// 预定义错误
var (
    Canceled         = errors.New("context canceled")
    DeadlineExceeded = errors.New("context deadline exceeded")
)
```

---

## 2. Context 实现

### 2.1 空 Context

```go
// emptyCtx 是 background 和 todo 的实现
type emptyCtx int

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
    return
}

func (emptyCtx) Done() <-chan struct{} {
    return nil
}

func (emptyCtx) Err() error {
    return nil
}

func (emptyCtx) Value(key any) any {
    return nil
}

var (
    background = new(emptyCtx)
    todo       = new(emptyCtx)
)

func Background() Context {
    return background
}

func TODO() Context {
    return todo
}
```

### 2.2 取消 Context

```go
// cancelCtx 是可取消的 Context
type cancelCtx struct {
    Context
    
    mu       sync.Mutex
    done     atomic.Value      // chan struct{}，懒加载
    children map[canceler]struct{} // 子节点集合
    err      error             // 取消原因
}

// canceler 接口
type canceler interface {
    cancel(removeFromParent bool, err error)
    Done() <-chan struct{}
}

func (c *cancelCtx) Done() <-chan struct{} {
    // 懒加载 done channel
    d := c.done.Load()
    if d != nil {
        return d.(chan struct{})
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    
    d = c.done.Load()
    if d == nil {
        d = make(chan struct{})
        c.done.Store(d)
    }
    return d.(chan struct{})
}

func (c *cancelCtx) Err() error {
    c.mu.Lock()
    err := c.err
    c.mu.Unlock()
    return err
}

func (c *cancelCtx) cancel(removeFromParent bool, err error) {
    if err == nil {
        panic("context: internal error: missing cancel error")
    }
    
    c.mu.Lock()
    if c.err != nil {
        c.mu.Unlock()
        return // 已取消
    }
    
    c.err = err
    
    // 关闭 done channel
    d, _ := c.done.Load().(chan struct{})
    if d == nil {
        c.done.Store(closedchan)
    } else {
        close(d)
    }
    
    // 递归取消子节点
    for child := range c.children {
        child.cancel(false, err)
    }
    c.children = nil
    c.mu.Unlock()
    
    // 从父节点移除
    if removeFromParent {
        removeChild(c.Context, c)
    }
}

// WithCancel 创建可取消的 Context
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    c := newCancelCtx(parent)
    propagateCancel(parent, c)
    return c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) *cancelCtx {
    return &cancelCtx{Context: parent}
}

// propagateCancel 建立父子关系
func propagateCancel(parent Context, child canceler) {
    done := parent.Done()
    if done == nil {
        // 父节点永远不会取消
        return
    }
    
    select {
    case <-done:
        // 父节点已取消
        child.cancel(false, parent.Err())
        return
    default:
    }
    
    // 将 child 添加到父节点的 children 集合
    if p, ok := parentCancelCtx(parent); ok {
        p.mu.Lock()
        if p.err != nil {
            // 父节点已取消
            p.mu.Unlock()
            child.cancel(false, p.err)
        } else {
            if p.children == nil {
                p.children = make(map[canceler]struct{})
            }
            p.children[child] = struct{}{}
            p.mu.Unlock()
        }
    } else {
        // 父节点实现了自定义的 cancel 机制
        // 启动 goroutine 监听
        go func() {
            select {
            case <-parent.Done():
                child.cancel(false, parent.Err())
            case <-child.Done():
            }
        }()
    }
}
```

### 2.3 超时 Context

```go
// timerCtx 带超时功能的 Context
type timerCtx struct {
    cancelCtx
    timer *time.Timer // 定时器
    deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
    return c.deadline, true
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
    // 停止定时器
    c.cancelCtx.cancel(false, err)
    if removeFromParent {
        removeChild(c.cancelCtx.Context, c)
    }
    c.mu.Lock()
    if c.timer != nil {
        c.timer.Stop()
        c.timer = nil
    }
    c.mu.Unlock()
}

// WithDeadline 创建带截止时间的 Context
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    
    // 检查父节点的 deadline
    if cur, ok := parent.Deadline(); ok && cur.Before(d) {
        // 父节点 deadline 更早，直接继承
        return WithCancel(parent)
    }
    
    c := &timerCtx{
        cancelCtx: newCancelCtx(parent),
        deadline:  d,
    }
    propagateCancel(parent, c)
    
    dur := time.Until(d)
    if dur <= 0 {
        c.cancel(true, DeadlineExceeded)
        return c, func() { c.cancel(false, Canceled) }
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.err == nil {
        c.timer = time.AfterFunc(dur, func() {
            c.cancel(true, DeadlineExceeded)
        })
    }
    
    return c, func() { c.cancel(true, Canceled) }
}

// WithTimeout 创建带超时的 Context
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    return WithDeadline(parent, time.Now().Add(timeout))
}
```

### 2.4 Value Context

```go
// valueCtx 存储键值对的 Context
type valueCtx struct {
    Context
    key, val any
}

func (c *valueCtx) Value(key any) any {
    if c.key == key {
        return c.val
    }
    return c.Context.Value(key)
}

// WithValue 创建带值的 Context
func WithValue(parent Context, key, val any) Context {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    if key == nil {
        panic("nil key")
    }
    if !reflectlite.TypeOf(key).Comparable() {
        panic("key is not comparable")
    }
    return &valueCtx{parent, key, val}
}
```

---

## 3. 内存分配模式

### 3.1 Context 分配分析

```go
// 每个派生操作都会分配
func analyzeAllocations() {
    ctx := context.Background()
    
    // 分配 cancelCtx (约 48 bytes)
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    // 分配 timerCtx (约 72 bytes + Timer)
    ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // 分配 valueCtx (约 32 bytes)
    ctx = context.WithValue(ctx, "key", "value")
    
    // 深度嵌套时的内存使用
    for i := 0; i < 100; i++ {
        ctx = context.WithValue(ctx, i, i) // O(n) 查找
    }
}
```

### 3.2 减少分配的技巧

```go
// 使用结构体作为 key 避免字符串分配
type contextKey struct {
    name string
}

var (
    userKey = &contextKey{"user"}
    traceKey = &contextKey{"trace_id"}
)

func SetUser(ctx context.Context, user *User) context.Context {
    return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) *User {
    user, _ := ctx.Value(userKey).(*User)
    return user
}

// 预创建常用 Context
var (
    shortTimeout = 5 * time.Second
    longTimeout  = 30 * time.Second
)

func withShortTimeout(parent context.Context) (context.Context, context.CancelFunc) {
    return context.WithTimeout(parent, shortTimeout)
}
```

---

## 4. 并发安全分析

### 4.1 线程安全保证

```go
// Context 的所有操作都是线程安全的

// 1. Done() 返回的 channel 可以被多个 goroutine 监听
func multiListener(ctx context.Context) {
    for i := 0; i < 10; i++ {
        go func() {
            <-ctx.Done()
            println("cancelled!")
        }()
    }
}

// 2. Value() 可以被并发读取
func concurrentRead(ctx context.Context) {
    for i := 0; i < 10; i++ {
        go func() {
            v := ctx.Value("key")
            _ = v
        }()
    }
}

// 3. cancel() 是幂等的，可以安全调用多次
func safeCancel(cancel context.CancelFunc) {
    cancel() // 第一次
    cancel() // 第二次，安全
}
```

### 4.2 竞态条件防护

```go
// 错误：在 goroutine 之间共享 cancel 函数可能导致竞态
func wrongPattern() {
    ctx, cancel := context.WithCancel(context.Background())
    
    for i := 0; i < 10; i++ {
        go func() {
            // 错误：多个 goroutine 可能同时调用 cancel
            if someCondition() {
                cancel()
            }
        }()
    }
}

// 正确：使用 sync.Once 或原子操作
func correctPattern() {
    ctx, cancel := context.WithCancel(context.Background())
    var once sync.Once
    
    for i := 0; i < 10; i++ {
        go func() {
            if someCondition() {
                once.Do(cancel) // 只执行一次
            }
        }()
    }
}
```

---

## 5. 性能优化

### 5.1 最佳实践

```go
// 1. 尽早取消，避免资源泄漏
func processRequest(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel() // 确保取消
    
    // 处理请求...
    return nil
}

// 2. 不要在 map 中存储 Context
// ❌ 错误
type BadRequest struct {
    ctx context.Context
    data []byte
}

// ✅ 正确：作为第一个参数传递
func GoodHandler(ctx context.Context, req *Request) error {
    // ...
}

// 3. 传递 nil Context 时处理
func safeHandler(ctx context.Context) {
    if ctx == nil {
        ctx = context.Background()
    }
    // ...
}

// 4. 避免过度嵌套
func flatContexts() {
    ctx := context.Background()
    
    // 合并多个 value
    ctx = withMetadata(ctx, Metadata{
        TraceID: "xxx",
        UserID:  "yyy",
        SpanID:  "zzz",
    })
    
    // 而不是：
    // ctx = context.WithValue(ctx, "trace_id", "xxx")
    // ctx = context.WithValue(ctx, "user_id", "yyy")
    // ctx = context.WithValue(ctx, "span_id", "zzz")
}
```

### 5.2 基准测试

```go
func BenchmarkWithCancel(b *testing.B) {
    ctx := context.Background()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, cancel := context.WithCancel(ctx)
        cancel()
    }
}

func BenchmarkWithValue(b *testing.B) {
    ctx := context.Background()
    key := struct{}{}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx = context.WithValue(ctx, key, i)
    }
}

func BenchmarkValueLookup(b *testing.B) {
    ctx := context.Background()
    // 构建深度为 100 的 Context 链
    for i := 0; i < 100; i++ {
        ctx = context.WithValue(ctx, i, i)
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx.Value(50)
    }
}

// 典型结果 (Go 1.21)
// BenchmarkWithCancel-8    20000000    65 ns/op    48 B/op    1 allocs/op
// BenchmarkWithValue-8     30000000    45 ns/op    32 B/op    1 allocs/op
// BenchmarkValueLookup-8   50000000    28 ns/op     0 B/op    0 allocs/op
```

---

## 6. 视觉表征

### 6.1 Context 树结构

```
Background()
    │
    ├── WithTimeout(5s)
    │       │
    │       ├── WithValue("user", u1)
    │       │       └── HTTP Handler 1
    │       │
    │       └── WithValue("user", u2)
    │               └── HTTP Handler 2
    │
    └── WithCancel
            │
            ├── WithValue("trace", t1)
            │       └── gRPC Call 1
            │
            └── WithDeadline(tomorrow)
                    └── Background Job
```

### 6.2 取消传播流程

```
Parent Context
      │ cancel()
      ▼
┌─────────────┐     ┌─────────────┐
│  close(done)│────►│  Child 1    │
└─────────────┘     │  (收到信号)  │
                    └──────┬──────┘
                           │ cancel()
                    ┌──────┴──────┐
                    │  GrandChild │
                    └─────────────┘

┌─────────────┐
│  Child 2    │
│  (同样收到)  │
└─────────────┘
```

### 6.3 使用决策树

```
需要 Context?
│
├── 没有请求作用域需求?
│   └── 使用 context.Background()
│
├── 需要取消信号?
│   ├── 手动触发? → WithCancel
│   └── 超时触发? → WithTimeout / WithDeadline
│
├── 需要存储数据?
│   └── WithValue (仅用于请求作用域数据)
│
└── 组合使用?
    └── ctx, cancel := WithTimeout(parent, timeout)
        ctx = WithValue(ctx, key, value)
```

---

## 7. 完整代码示例

### 7.1 HTTP 服务器集成

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

// 中间件：添加上下文信息
func contextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 添加请求 ID
        reqID := generateRequestID()
        ctx = withRequestID(ctx, reqID)
        
        // 添加超时
        ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
        defer cancel()
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 中间件：日志
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        reqID := getRequestID(r.Context())
        
        next.ServeHTTP(w, r)
        
        fmt.Printf("[%s] %s %s %v\n",
            reqID,
            r.Method,
            r.URL.Path,
            time.Since(start),
        )
    })
}

// 处理函数
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // 检查取消
    if err := ctx.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusRequestTimeout)
        return
    }
    
    // 带超时的数据库查询
    user, err := getUserWithContext(ctx, "123")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "User: %+v\n", user)
}

// 带上下文的数据库查询
func getUserWithContext(ctx context.Context, id string) (*User, error) {
    // 使用 ctx 创建带超时的查询
    type result struct {
        user *User
        err  error
    }
    
    done := make(chan result, 1)
    
    go func() {
        user, err := db.GetUser(id)
        done <- result{user, err}
    }()
    
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case r := <-done:
        return r.user, r.err
    }
}

// 键类型（避免冲突）
type requestIDKey struct{}

func withRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey{}, id)
}

func getRequestID(ctx context.Context) string {
    id, _ := ctx.Value(requestIDKey{}).(string)
    return id
}

type User struct {
    ID   string
    Name string
}

var db = &mockDB{}

type mockDB struct{}

func (m *mockDB) GetUser(id string) (*User, error) {
    time.Sleep(100 * time.Millisecond)
    return &User{ID: id, Name: "Test"}, nil
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)
    
    // 应用中间件
    handler := contextMiddleware(loggingMiddleware(mux))
    
    fmt.Println("Server started on :8080")
    http.ListenAndServe(":8080", handler)
}
```

### 7.2 Pipeline 模式

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// 带取消的 pipeline
func pipeline(ctx context.Context, inputs []int) <-chan int {
    out := make(chan int)
    
    go func() {
        defer close(out)
        
        for _, n := range inputs {
            select {
            case <-ctx.Done():
                fmt.Println("Pipeline cancelled")
                return
            case out <- n * 2:
            }
        }
    }()
    
    return out
}

// 带超时的处理
func processWithTimeout(inputs []int, timeout time.Duration) ([]int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    out := pipeline(ctx, inputs)
    
    var results []int
    for n := range out {
        results = append(results, n)
    }
    
    if err := ctx.Err(); err != nil {
        return results, err
    }
    
    return results, nil
}

func main() {
    inputs := []int{1, 2, 3, 4, 5}
    
    // 成功场景
    results, err := processWithTimeout(inputs, 5*time.Second)
    fmt.Printf("Results: %v, Error: %v\n", results, err)
    
    // 超时场景
    results, err = processWithTimeout(inputs, 1*time.Nanosecond)
    fmt.Printf("Results: %v, Error: %v\n", results, err)
}
```

---

**质量评级**: S (17KB)
**完成日期**: 2026-04-02
