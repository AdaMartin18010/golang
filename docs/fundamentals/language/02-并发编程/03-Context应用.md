# Context应用

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

---

## 📋 目录

- [1. Context简介](#1-context简介)
  - [什么是Context](#什么是context)
  - [为什么需要Context](#为什么需要context)
- [2. 创建Context](#2-创建context)
  - [Background和TODO](#background和todo)
  - [WithCancel](#withcancel)
  - [WithTimeout](#withtimeout)
  - [WithDeadline](#withdeadline)
- [3. 超时控制](#3-超时控制)
  - [HTTP请求超时](#http请求超时)
  - [数据库查询超时](#数据库查询超时)
- [4. 取消传播](#4-取消传播)
  - [级联取消](#级联取消)
  - [优雅关闭](#优雅关闭)
- [5. 值传递](#5-值传递)
  - [WithValue](#withvalue)
  - [类型安全的值传递](#类型安全的值传递)
- [6. 实战应用](#6-实战应用)
  - [HTTP服务器中间件](#http服务器中间件)
  - [并行任务处理](#并行任务处理)
- [7. 最佳实践](#7-最佳实践)
  - [1. Context作为第一个参数](#1-context作为第一个参数)
  - [2. 不要存储Context](#2-不要存储context)
  - [3. 总是defer cancel()](#3-总是defer-cancel)
  - [4. 检查Context错误](#4-检查context错误)
  - [5. 不要传递nil Context](#5-不要传递nil-context)
  - [6. Context值只用于请求作用域数据](#6-context值只用于请求作用域数据)
- [🔗 相关资源](#相关资源)

## 1. Context简介

### 什么是Context

**Context** 是Go中管理请求生命周期的标准方式：

- 传递请求作用域的值
- 传递取消信号
- 传递截止时间/超时
- 跨API边界和进程

### 为什么需要Context

```go
// ❌ 没有Context：难以取消
func process() {
    for {
        doWork()
        // 如何停止？
    }
}

// ✅ 有Context：可以优雅取消
func process(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return  // 收到取消信号
        default:
            doWork()
        }
    }
}
```

---

## 2. 创建Context

### Background和TODO

```go
import "context"

// Background: 根Context，通常用于main、init和测试
ctx := context.Background()

// TODO: 不确定使用哪个Context时的占位符
ctx := context.TODO()
```

---

### WithCancel

```go
// 创建可取消的Context
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // 确保释放资源

go func() {
    <-ctx.Done()
    fmt.Println("Context canceled")
}()

// 取消Context
cancel()
```

**完整示例**:

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Worker stopped:", ctx.Err())
                return
            default:
                fmt.Println("Working...")
                time.Sleep(1 * time.Second)
            }
        }
    }()
    
    time.Sleep(3 * time.Second)
    cancel()  // 取消Context
    time.Sleep(1 * time.Second)
}
```

---

### WithTimeout

```go
// 创建带超时的Context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case result := <-doWork(ctx):
    fmt.Println("Result:", result)
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())
}
```

**完整示例**:

```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    go func() {
        select {
        case <-time.After(3 * time.Second):
            fmt.Println("Work done")
        case <-ctx.Done():
            fmt.Println("Timeout:", ctx.Err())
        }
    }()
    
    time.Sleep(3 * time.Second)
}
```

---

### WithDeadline

```go
// 创建带截止时间的Context
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

select {
case <-doWork(ctx):
    fmt.Println("Done")
case <-ctx.Done():
    fmt.Println("Deadline exceeded:", ctx.Err())
}
```

---

## 3. 超时控制

### HTTP请求超时

```go
func fetchURL(ctx context.Context, url string) (string, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    return string(body), err
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    body, err := fetchURL(ctx, "https://example.com")
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Request timeout")
        } else {
            fmt.Println("Error:", err)
        }
        return
    }
    
    fmt.Println("Body:", body)
}
```

---

### 数据库查询超时

```go
func queryDatabase(ctx context.Context, query string) ([]User, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    
    return users, rows.Err()
}
```

---

## 4. 取消传播

### 级联取消

```go
func main() {
    // 父Context
    parentCtx, parentCancel := context.WithCancel(context.Background())
    defer parentCancel()
    
    // 子Context1
    childCtx1, cancel1 := context.WithCancel(parentCtx)
    defer cancel1()
    
    // 子Context2
    childCtx2, cancel2 := context.WithCancel(parentCtx)
    defer cancel2()
    
    go worker(childCtx1, "Worker1")
    go worker(childCtx2, "Worker2")
    
    time.Sleep(2 * time.Second)
    
    // 取消父Context，所有子Context也会被取消
    parentCancel()
    
    time.Sleep(1 * time.Second)
}

func worker(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s stopped\n", name)
            return
        default:
            fmt.Printf("%s working...\n", name)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

---

### 优雅关闭

```go
func server(ctx context.Context) {
    srv := &http.Server{Addr: ":8080"}
    
    go func() {
        <-ctx.Done()
        
        // 优雅关闭，最多等待5秒
        shutdownCtx, cancel := context.WithTimeout(
            context.Background(),
            5*time.Second,
        )
        defer cancel()
        
        if err := srv.Shutdown(shutdownCtx); err != nil {
            log.Println("Server shutdown error:", err)
        }
    }()
    
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go server(ctx)
    
    // 等待中断信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    cancel()  // 取消Context，触发优雅关闭
    time.Sleep(6 * time.Second)
}
```

---

## 5. 值传递

### WithValue

```go
type key string

func main() {
    ctx := context.WithValue(context.Background(), key("userID"), 123)
    ctx = context.WithValue(ctx, key("requestID"), "abc-123")
    
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    userID := ctx.Value(key("userID")).(int)
    requestID := ctx.Value(key("requestID")).(string)
    
    fmt.Printf("UserID: %d, RequestID: %s\n", userID, requestID)
}
```

---

### 类型安全的值传递

```go
type userIDKey struct{}
type requestIDKey struct{}

func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey{}, userID)
}

func GetUserID(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey{}).(int)
    return userID, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey{}, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
    requestID, ok := ctx.Value(requestIDKey{}).(string)
    return requestID, ok
}

// 使用
func main() {
    ctx := context.Background()
    ctx = WithUserID(ctx, 123)
    ctx = WithRequestID(ctx, "abc-123")
    
    if userID, ok := GetUserID(ctx); ok {
        fmt.Println("UserID:", userID)
    }
}
```

---

## 6. 实战应用

### HTTP服务器中间件

```go
func withRequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        ctx := WithRequestID(r.Context(), requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    requestID, _ := GetRequestID(r.Context())
    log.Printf("RequestID: %s", requestID)
    
    // 处理请求
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", handler)
    
    http.ListenAndServe(":8080", withRequestID(mux))
}
```

---

### 并行任务处理

```go
func processItems(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)
    
    for _, item := range items {
        item := item  // 捕获变量
        g.Go(func() error {
            return processItem(ctx, item)
        })
    }
    
    return g.Wait()
}

func processItem(ctx context.Context, item Item) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // 处理item
        return nil
    }
}
```

---

## 7. 最佳实践

### 1. Context作为第一个参数

```go
// ✅ 推荐
func doWork(ctx context.Context, arg string) error {
    // ...
}

// ❌ 不推荐
func doWork(arg string, ctx context.Context) error {
    // ...
}
```

---

### 2. 不要存储Context

```go
// ❌ 不要存储Context
type Server struct {
    ctx context.Context  // 错误
}

// ✅ 通过参数传递
type Server struct {
    // 其他字段
}

func (s *Server) Handle(ctx context.Context) {
    // 使用ctx
}
```

---

### 3. 总是defer cancel()

```go
// ✅ 推荐
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // 确保释放资源

// ❌ 忘记调用cancel会导致资源泄漏
```

---

### 4. 检查Context错误

```go
// ✅ 推荐
select {
case <-ctx.Done():
    if ctx.Err() == context.Canceled {
        return errors.New("operation was canceled")
    }
    if ctx.Err() == context.DeadlineExceeded {
        return errors.New("operation timed out")
    }
default:
    // 继续执行
}
```

---

### 5. 不要传递nil Context

```go
// ❌ 不要传递nil
doWork(nil)

// ✅ 使用context.TODO()
doWork(context.TODO())

// ✅ 或使用context.Background()
doWork(context.Background())
```

---

### 6. Context值只用于请求作用域数据

```go
// ✅ 适合：请求ID、用户ID、追踪ID
ctx = WithRequestID(ctx, "abc-123")

// ❌ 不适合：可选参数、配置
// 不要用Context传递函数的可选参数
```

---

## 🔗 相关资源

- [Goroutine基础](./01-Goroutine基础.md)
- [Channel详解](./02-Channel详解.md)
- [并发模式](./05-并发模式.md)

---

**最后更新**: 2025-10-29  
**Go版本**: 1.25.3
