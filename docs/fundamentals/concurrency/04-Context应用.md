# Context应用

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

## 📋 目录

- [Context应用](#context应用)
  - [📋 目录](#-目录)
  - [1. 📖 概念介绍](#1--概念介绍)
  - [2. 🎯 核心知识点](#2--核心知识点)
    - [1. Context的设计理念](#1-context的设计理念)
      - [Context接口定义](#context接口定义)
      - [Context的核心原则](#context的核心原则)
    - [2. 四种Context类型](#2-四种context类型)
      - [Background和TODO](#background和todo)
      - [WithCancel](#withcancel)
      - [WithTimeout](#withtimeout)
      - [WithDeadline](#withdeadline)
      - [WithValue](#withvalue)
    - [3. 超时控制实战](#3-超时控制实战)
      - [HTTP请求超时](#http请求超时)
      - [数据库查询超时](#数据库查询超时)
    - [4. 取消信号传播](#4-取消信号传播)
      - [父子Context取消传播](#父子context取消传播)
      - [多层Goroutine取消](#多层goroutine取消)
    - [5. 值传递最佳实践](#5-值传递最佳实践)
      - [正确的值传递](#正确的值传递)
      - [错误的值传递](#错误的值传递)
    - [6. Context在HTTP中的应用](#6-context在http中的应用)
      - [HTTP服务器中的Context](#http服务器中的context)
      - [HTTP客户端中的Context](#http客户端中的context)
  - [🏗️ 实战案例](#️-实战案例)
    - [案例：Pipeline with Context](#案例pipeline-with-context)
  - [⚠️ 常见问题](#️-常见问题)
    - [Q1: Context应该在什么时候取消？](#q1-context应该在什么时候取消)
    - [Q2: Context.Value应该存储什么？](#q2-contextvalue应该存储什么)
    - [Q3: Context会泄漏吗？](#q3-context会泄漏吗)
    - [Q4: 如何测试使用Context的代码？](#q4-如何测试使用context的代码)
  - [📚 相关资源](#-相关资源)
    - [下一步学习](#下一步学习)
    - [推荐阅读](#推荐阅读)

## 1. 📖 概念介绍

Context是Go 1.7引入的标准库包，用于在Goroutine之间传递取消信号、超时控制和请求范围的值。它是构建健壮并发程序的重要工具。

---

## 2. 🎯 核心知识点

### 1. Context的设计理念

#### Context接口定义

```go
type Context interface {
    // Deadline返回context的过期时间
    Deadline() (deadline time.Time, ok bool)

    // Done返回一个channel，当context被取消或过期时关闭
    Done() <-chan struct{}

    // Err在Done channel关闭后返回错误原因
    Err() error

    // Value返回context关联的key对应的值
    Value(key interface{}) interface{}
}
```

#### Context的核心原则

```go
package main

import (
    "context"
    "fmt"
)

/*
Context设计原则：
1. 不要存储Context，而是显式传递
2. Context作为函数的第一个参数，命名为ctx
3. 不要传递nil Context，使用context.TODO()
4. Context只传递请求相关的值，不传递可选参数
5. Context是不可变的（immutable）
*/

// ✅ 正确示例
func doSomething(ctx context.Context, arg string) error {
    // ctx作为第一个参数
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        fmt.Println("Working on:", arg)
        return nil
    }
}

// ❌ 错误示例
type Worker struct {
    ctx context.Context // 不要存储Context
}

func main() {
    ctx := context.Background()
    doSomething(ctx, "task1")
}
```

---

### 2. 四种Context类型

#### Background和TODO

```go
package main

import (
    "context"
    "fmt"
)

func contextRoots() {
    // Background：根Context，永不取消，通常在main、init、测试中使用
    ctx1 := context.Background()
    fmt.Printf("Background: %v\n", ctx1)

    // TODO：当不确定使用哪个Context时使用（临时占位）
    ctx2 := context.TODO()
    fmt.Printf("TODO: %v\n", ctx2)
}

func main() {
    contextRoots()
}
```

#### WithCancel

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func withCancelExample() {
    // 创建可取消的Context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // 确保释放资源

    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine cancelled:", ctx.Err())
                return
            default:
                fmt.Println("Working...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()

    // 2秒后取消
    time.Sleep(2 * time.Second)
    cancel()

    time.Sleep(1 * time.Second)
}

func main() {
    withCancelExample()
}
```

#### WithTimeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func withTimeoutExample() {
    // 创建带超时的Context（3秒后自动取消）
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Timeout:", ctx.Err())
                return
            default:
                fmt.Println("Processing...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()

    time.Sleep(5 * time.Second)
}

func main() {
    withTimeoutExample()
}
```

#### WithDeadline

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func withDeadlineExample() {
    // 创建有截止时间的Context
    deadline := time.Now().Add(2 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()

    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Deadline reached:", ctx.Err())
                return
            default:
                fmt.Println("Working until deadline...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()

    time.Sleep(3 * time.Second)
}

func main() {
    withDeadlineExample()
}
```

#### WithValue

```go
package main

import (
    "context"
    "fmt"
)

// 定义类型化的key，避免冲突
type contextKey string

const (
    userIDKey contextKey = "userID"
    traceIDKey contextKey = "traceID"
)

func withValueExample() {
    // 创建带值的Context
    ctx := context.WithValue(context.Background(), userIDKey, "12345")
    ctx = context.WithValue(ctx, traceIDKey, "trace-abc")

    // 读取值
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    // 类型断言获取值
    if userID, ok := ctx.Value(userIDKey).(string); ok {
        fmt.Printf("Processing request for user: %s\n", userID)
    }

    if traceID, ok := ctx.Value(traceIDKey).(string); ok {
        fmt.Printf("Trace ID: %s\n", traceID)
    }

    // 调用其他函数，传递context
    doWork(ctx)
}

func doWork(ctx context.Context) {
    userID := ctx.Value(userIDKey)
    fmt.Printf("DoWork for user: %v\n", userID)
}

func main() {
    withValueExample()
}
```

---

### 3. 超时控制实战

#### HTTP请求超时

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
    // 创建带超时的Context
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // 创建带Context的HTTP请求
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }

    // 执行请求
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

func main() {
    result, err := fetchWithTimeout("https://httpbin.org/delay/2", 3*time.Second)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Result length: %d\n", len(result))
}
```

#### 数据库查询超时

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "time"
)

func queryWithTimeout(db *sql.DB) error {
    // 3秒超时
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // 使用Context执行查询
    rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE age > ?", 18)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var name string
        var age int

        if err := rows.Scan(&id, &name, &age); err != nil {
            return err
        }

        fmt.Printf("User: %d, %s, %d\n", id, name, age)
    }

    return rows.Err()
}

// 示例函数（实际使用需要真实数据库连接）
func databaseExample() {
    // db, _ := sql.Open("mysql", "user:password@/dbname")
    // queryWithTimeout(db)
    fmt.Println("Database query with timeout example")
}

func main() {
    databaseExample()
}
```

---

### 4. 取消信号传播

#### 父子Context取消传播

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func cancelPropagation() {
    // 创建根Context
    parent, parentCancel := context.WithCancel(context.Background())
    defer parentCancel()

    // 创建子Context
    child1, child1Cancel := context.WithCancel(parent)
    defer child1Cancel()

    child2, child2Cancel := context.WithCancel(parent)
    defer child2Cancel()

    // 子Goroutine 1
    go func() {
        <-child1.Done()
        fmt.Println("Child1 cancelled:", child1.Err())
    }()

    // 子Goroutine 2
    go func() {
        <-child2.Done()
        fmt.Println("Child2 cancelled:", child2.Err())
    }()

    time.Sleep(1 * time.Second)

    // 取消父Context会自动取消所有子Context
    fmt.Println("Cancelling parent...")
    parentCancel()

    time.Sleep(1 * time.Second)
}

func main() {
    cancelPropagation()
}
```

#### 多层Goroutine取消

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: cancelled\n", name)
            return
        default:
            fmt.Printf("%s: working...\n", name)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func supervisor(ctx context.Context, name string) {
    // 创建子Context
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // 启动多个worker
    for i := 0; i < 3; i++ {
        go worker(ctx, fmt.Sprintf("%s-worker-%d", name, i))
    }

    // 等待取消信号
    <-ctx.Done()
    fmt.Printf("%s: shutting down workers...\n", name)
    cancel() // 取消所有worker
    time.Sleep(1 * time.Second)
}

func multiLayerCancellation() {
    ctx, cancel := context.WithCancel(context.Background())

    go supervisor(ctx, "Supervisor-A")
    go supervisor(ctx, "Supervisor-B")

    time.Sleep(2 * time.Second)
    fmt.Println("Main: cancelling all...")
    cancel()

    time.Sleep(2 * time.Second)
}

func main() {
    multiLayerCancellation()
}
```

---

### 5. 值传递最佳实践

#### 正确的值传递

```go
package main

import (
    "context"
    "fmt"
)

// 定义类型化的key
type requestKey string

const (
    requestIDKey requestKey = "requestID"
    userKey      requestKey = "user"
)

// User结构体
type User struct {
    ID   string
    Name string
}

// ✅ 正确：只传递请求相关的值
func goodPractice() {
    ctx := context.Background()
    ctx = context.WithValue(ctx, requestIDKey, "req-123")
    ctx = context.WithValue(ctx, userKey, User{ID: "u1", Name: "Alice"})

    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    requestID := ctx.Value(requestIDKey).(string)
    user := ctx.Value(userKey).(User)

    fmt.Printf("Processing request %s for user %s\n", requestID, user.Name)

    // 传递给其他函数
    logRequest(ctx)
}

func logRequest(ctx context.Context) {
    requestID, ok := ctx.Value(requestIDKey).(string)
    if !ok {
        fmt.Println("No request ID in context")
        return
    }
    fmt.Printf("Logging request: %s\n", requestID)
}

func main() {
    goodPractice()
}
```

#### 错误的值传递

```go
package main

import (
    "context"
    "fmt"
)

// ❌ 错误：不要传递可选参数或配置
type Config struct {
    MaxRetries int
    Timeout    int
}

func badPractice(ctx context.Context) {
    // ❌ 不要这样做
    config := ctx.Value("config").(Config)
    fmt.Printf("Config: %+v\n", config)
}

// ✅ 正确：显式传递配置参数
func goodPractice(ctx context.Context, config Config) {
    fmt.Printf("Config: %+v\n", config)
}

func main() {
    // 配置应该显式传递，不要放在Context中
    config := Config{MaxRetries: 3, Timeout: 5}
    goodPractice(context.Background(), config)
}
```

---

### 6. Context在HTTP中的应用

#### HTTP服务器中的Context

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // HTTP请求自带Context
    ctx := r.Context()

    // 添加请求ID
    requestID := r.Header.Get("X-Request-ID")
    if requestID == "" {
        requestID = "generated-id"
    }
    ctx = context.WithValue(ctx, "requestID", requestID)

    // 模拟长时间处理
    select {
    case <-time.After(5 * time.Second):
        fmt.Fprintf(w, "Request completed: %s\n", requestID)
    case <-ctx.Done():
        // 客户端断开连接
        fmt.Printf("Request cancelled: %s, error: %v\n", requestID, ctx.Err())
        http.Error(w, "Request cancelled", 499)
    }
}

func httpServerExample() {
    http.HandleFunc("/", handler)
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}

// 取消注释以运行服务器
// func main() {
//     httpServerExample()
// }
```

#### HTTP客户端中的Context

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func httpClientWithContext() {
    // 创建带超时的Context
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // 创建请求
    req, _ := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/3", nil)

    // 执行请求
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("Request failed:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response: %d bytes\n", len(body))
}

func main() {
    httpClientWithContext()
}
```

---

## 🏗️ 实战案例

### 案例：Pipeline with Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func generator(ctx context.Context, nums ...int) <-chan int {
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

func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func pipelineExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // 构建pipeline
    ch := generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    ch = square(ctx, ch)

    // 消费结果
    for n := range ch {
        fmt.Println(n)
        time.Sleep(500 * time.Millisecond)
    }
}

func main() {
    pipelineExample()
}
```

---

## ⚠️ 常见问题

### Q1: Context应该在什么时候取消？

- 任务完成后立即取消
- 使用defer cancel()确保释放资源
- 超时后自动取消

### Q2: Context.Value应该存储什么？

- ✅ 请求范围的值（requestID、traceID、用户信息）
- ❌ 可选参数、配置、业务数据

### Q3: Context会泄漏吗？

- 如果不调用cancel，会导致资源泄漏
- 使用defer cancel()确保释放
- 父Context取消会自动清理子Context

### Q4: 如何测试使用Context的代码？

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    err := doWork(ctx)
    if err != context.DeadlineExceeded {
        t.Errorf("Expected timeout, got %v", err)
    }
}
```

---

## 📚 相关资源

### 下一步学习

- [05-并发模式](./05-并发模式.md)
- [HTTP服务器](../../development/web/03-HTTP服务器.md)

### 推荐阅读

- [Go Blog - Context](https://go.dev/blog/context)
- [Context Package Doc](https://pkg.go.dev/context)

---

**最后更新**: 2025-10-29
**作者**: Documentation Team
