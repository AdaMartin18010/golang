# LD-023: Go 错误处理模式详解 (Go Error Handling Patterns Deep Dive)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #error-handling #errors #wrapping #sentinel #custom-errors #go113
> **权威来源**:
>
> - [Error Handling and Go](https://go.dev/blog/error-handling-and-go) - Go Authors
> - [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) - Damien Neil
> - [Error Value FAQ](https://github.com/golang/go/wiki/ErrorValueFAQ) - Go Wiki

---

## 1. 错误接口与基础

### 1.1 error 接口

```go
// 内置 error 接口
type error interface {
    Error() string
}

// 最简单实现
func New(text string) error {
    return &errorString{text}
}

type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

### 1.2 错误创建方式

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    // 方式 1: errors.New (静态错误)
    err1 := errors.New("something went wrong")
    
    // 方式 2: fmt.Errorf (格式化错误)
    err2 := fmt.Errorf("user %d not found", 42)
    
    // 方式 3: 格式化 + 包装 (Go 1.13+)
    err3 := fmt.Errorf("database error: %w", err1)
    
    // 方式 4: 自定义错误类型
    err4 := &NotFoundError{Resource: "user", ID: "42"}
}

// 自定义错误类型
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
}
```

---

## 2. Go 1.13+ 错误包装

### 2.1 包装与展开

```go
// src/errors/wrap.go

// 包装错误接口
type wrapper interface {
    Unwrap() error
}

// errors.Is 检查错误链中是否存在目标错误
func Is(err, target error) bool {
    if target == nil {
        return err == target
    }
    
    isComparable := reflectlite.TypeOf(target).Comparable()
    
    for {
        if isComparable && err == target {
            return true
        }
        if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
            return true
        }
        
        // 展开错误链
        switch x := err.(type) {
        case interface{ Unwrap() error }:
            err = x.Unwrap()
            if err == nil {
                return false
            }
        case interface{ Unwrap() []error }:
            for _, err := range x.Unwrap() {
                if Is(err, target) {
                    return true
                }
            }
            return false
        default:
            return false
        }
    }
}

// errors.As 将错误转换为特定类型
func As(err error, target any) bool {
    if target == nil {
        panic("errors: target cannot be nil")
    }
    
    val := reflectlite.ValueOf(target)
    typ := val.Type()
    
    for err != nil {
        if reflectlite.TypeOf(err).AssignableTo(typ) {
            val.Elem().Set(reflectlite.ValueOf(err))
            return true
        }
        
        if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
            return true
        }
        
        switch x := err.(type) {
        case interface{ Unwrap() error }:
            err = x.Unwrap()
        case interface{ Unwrap() []error }:
            for _, err := range x.Unwrap() {
                if As(err, target) {
                    return true
                }
            }
            return false
        default:
            return false
        }
    }
    return false
}
```

### 2.2 哨兵错误模式

```go
package main

import (
    "errors"
    "fmt"
    "io"
)

// 定义哨兵错误
var (
    ErrNotFound     = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
    ErrTimeout      = errors.New("operation timed out")
    ErrConflict     = errors.New("resource conflict")
)

// 使用示例
func fetchUser(id string) error {
    // 模拟未找到
    return fmt.Errorf("%w: user id=%s", ErrNotFound, id)
}

func main() {
    err := fetchUser("123")
    
    // 检查特定错误
    if errors.Is(err, ErrNotFound) {
        fmt.Println("User not found, handle appropriately")
    }
    
    // 检查标准库错误
    if errors.Is(err, io.EOF) {
        // ...
    }
}
```

---

## 3. 自定义错误类型

### 3.1 完整错误类型实现

```go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "time"
)

// 应用错误接口
type AppError interface {
    error
    Code() string
    HTTPStatus() int
    Details() map[string]interface{}
}

// 基础应用错误
type appError struct {
    code       string
    message    string
    status     int
    details    map[string]interface{}
    wrappedErr error
    timestamp  time.Time
    stack      string
}

func (e *appError) Error() string {
    if e.wrappedErr != nil {
        return fmt.Sprintf("[%s] %s: %v", e.code, e.message, e.wrappedErr)
    }
    return fmt.Sprintf("[%s] %s", e.code, e.message)
}

func (e *appError) Code() string {
    return e.code
}

func (e *appError) HTTPStatus() int {
    return e.status
}

func (e *appError) Details() map[string]interface{} {
    return e.details
}

func (e *appError) Unwrap() error {
    return e.wrappedErr
}

func (e *appError) MarshalJSON() ([]byte, error) {
    return json.Marshal(map[string]interface{}{
        "code":      e.code,
        "message":   e.message,
        "status":    e.status,
        "details":   e.details,
        "timestamp": e.timestamp,
    })
}

// 错误构建器
type ErrorBuilder struct {
    err *appError
}

func NewError(code, message string) *ErrorBuilder {
    return &ErrorBuilder{
        err: &appError{
            code:      code,
            message:   message,
            status:    http.StatusInternalServerError,
            timestamp: time.Now(),
        },
    }
}

func (b *ErrorBuilder) WithStatus(status int) *ErrorBuilder {
    b.err.status = status
    return b
}

func (b *ErrorBuilder) WithDetails(details map[string]interface{}) *ErrorBuilder {
    b.err.details = details
    return b
}

func (b *ErrorBuilder) Wrap(err error) *ErrorBuilder {
    b.err.wrappedErr = err
    return b
}

func (b *ErrorBuilder) Error() error {
    return b.err
}

// 使用示例
func processOrder(orderID string) error {
    // 一些处理...
    
    return NewError("ORDER_NOT_FOUND", "Order not found").
        WithStatus(http.StatusNotFound).
        WithDetails(map[string]interface{}{
            "order_id": orderID,
            "retry_after": 60,
        }).
        Error()
}
```

### 3.2 错误链处理

```go
// 多层错误包装
func serviceLayer() error {
    if err := repositoryLayer(); err != nil {
        return fmt.Errorf("service failed: %w", err)
    }
    return nil
}

func repositoryLayer() error {
    if err := dbLayer(); err != nil {
        return fmt.Errorf("repository failed: %w", err)
    }
    return nil
}

func dbLayer() error {
    return fmt.Errorf("%w: user id=123", ErrNotFound)
}

// 错误展开
func inspectError() {
    err := serviceLayer()
    
    // 打印完整错误链
    for err != nil {
        fmt.Println(err)
        
        if wrapped, ok := err.(interface{ Unwrap() error }); ok {
            err = wrapped.Unwrap()
        } else {
            break
        }
    }
    // 输出:
    // service failed: repository failed: resource not found: user id=123
    // repository failed: resource not found: user id=123
    // resource not found: user id=123
}
```

---

## 4. 错误处理模式

### 4.1 策略模式

```go
// 错误处理策略
type ErrorHandler func(error) error

// 重试策略
func WithRetry(maxRetries int, delay time.Duration) ErrorHandler {
    return func(err error) error {
        var lastErr error
        for i := 0; i < maxRetries; i++ {
            if lastErr == nil {
                return nil
            }
            time.Sleep(delay * time.Duration(i+1))
        }
        return fmt.Errorf("max retries exceeded: %w", lastErr)
    }
}

// 降级策略
func WithFallback(fallback func() error) ErrorHandler {
    return func(err error) error {
        if err == nil {
            return nil
        }
        return fallback()
    }
}

// 使用
func fetchDataWithStrategy() error {
    err := fetchFromPrimary()
    if err != nil {
        handler := WithFallback(func() error {
            return fetchFromCache()
        })
        return handler(err)
    }
    return nil
}
```

### 4.2 错误聚合

```go
// Go 1.20+ 错误列表
type joinError struct {
    errs []error
}

func (e *joinError) Error() string {
    var b []byte
    for i, err := range e.errs {
        if i > 0 {
            b = append(b, '\n')
        }
        b = append(b, err.Error()...)
    }
    return string(b)
}

func (e *joinError) Unwrap() []error {
    return e.errs
}

// errors.Join (Go 1.20+)
func Join(errs ...error) error {
    // ...
}

// 并行错误处理
func processBatch(items []Item) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(items))
    
    for _, item := range items {
        wg.Add(1)
        go func(it Item) {
            defer wg.Done()
            if err := process(it); err != nil {
                errChan <- fmt.Errorf("item %s: %w", it.ID, err)
            }
        }(item)
    }
    
    wg.Wait()
    close(errChan)
    
    var errs []error
    for err := range errChan {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return errors.Join(errs...)
    }
    return nil
}
```

---

## 5. 性能优化

### 5.1 错误创建开销

```go
// 预分配哨兵错误（推荐）
var (
    ErrNotFound = errors.New("not found")
    ErrTimeout  = errors.New("timeout")
)

// 动态错误（有分配开销）
func dynamicError(id int) error {
    return fmt.Errorf("item %d not found", id) // 分配
}

// 预格式化错误（优化）
var errorTemplates = map[int]error{
    404: errors.New("not found"),
    500: errors.New("internal error"),
}

// 零分配错误检查
func isNotFound(err error) bool {
    return errors.Is(err, ErrNotFound) // 无分配
}
```

### 5.2 栈追踪

```go
// 带栈追踪的错误
type stackError struct {
    msg   string
    stack []uintptr
}

func NewStackError(msg string) error {
    return &stackError{
        msg:   msg,
        stack: callers(),
    }
}

func callers() []uintptr {
    var pcs [32]uintptr
    n := runtime.Callers(3, pcs[:])
    return pcs[:n]
}

func (e *stackError) Error() string {
    return e.msg
}

func (e *stackError) StackTrace() string {
    var buf strings.Builder
    frames := runtime.CallersFrames(e.stack)
    for {
        frame, more := frames.Next()
        fmt.Fprintf(&buf, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
        if !more {
            break
        }
    }
    return buf.String()
}
```

---

## 6. 最佳实践与反模式

### 6.1 最佳实践

```go
// ✅ 1. 使用 %w 包装错误
func fetchUser(id string) error {
    if err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user); err != nil {
        return fmt.Errorf("fetch user %s: %w", id, err)
    }
    return nil
}

// ✅ 2. 使用哨兵错误
if errors.Is(err, sql.ErrNoRows) {
    return ErrNotFound
}

// ✅ 3. 使用 As 转换错误类型
var dbErr *pq.Error
if errors.As(err, &dbErr) {
    // 处理 PostgreSQL 特定错误
}

// ✅ 4. 早期返回
func process() error {
    if err := step1(); err != nil {
        return err
    }
    if err := step2(); err != nil {
        return err
    }
    return step3()
}

// ✅ 5. 集中错误处理
func handler(w http.ResponseWriter, r *http.Request) {
    if err := doWork(); err != nil {
        handleError(w, err)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func handleError(w http.ResponseWriter, err error) {
    var appErr AppError
    if errors.As(err, &appErr) {
        w.WriteHeader(appErr.HTTPStatus())
        json.NewEncoder(w).Encode(appErr)
        return
    }
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
```

### 6.2 反模式

```go
// ❌ 1. 忽略错误
_ = doSomething() // 危险!

// ❌ 2. 丢失原始错误信息
if err != nil {
    return errors.New("failed") // 丢失了原始 err
}

// ❌ 3. 使用 panic 处理错误
if err != nil {
    panic(err) // 除非不可恢复，否则不要用
}

// ❌ 4. 错误信息大写开头（不符合惯例）
return errors.New("Something went wrong") // 应该是 "something went wrong"

// ❌ 5. 过度包装
if err != nil {
    return fmt.Errorf("layer1: %w", fmt.Errorf("layer2: %w", fmt.Errorf("layer3: %w", err)))
}

// ❌ 6. 在错误中使用敏感信息
return fmt.Errorf("user %s with password %s failed login", username, password)
```

---

## 7. 视觉表征

### 7.1 错误链结构

```
HTTP Handler Error
       │
       ▼
┌─────────────────┐
│ "request failed │
│  to process"    │
└────────┬────────┘
         │ Unwrap()
         ▼
┌─────────────────┐
│ "service layer  │
│  error"         │
└────────┬────────┘
         │ Unwrap()
         ▼
┌─────────────────┐
│ "repository     │
│  error"         │
└────────┬────────┘
         │ Unwrap()
         ▼
┌─────────────────┐
│ sql.ErrNoRows   │  ◄── errors.Is() 可检测
└─────────────────┘
```

### 7.2 错误处理决策树

```
收到错误?
│
├── nil? → 继续正常流程
│
└── 非 nil?
    │
    ├── 可恢复?
    │   ├── 重试? → 指数退避重试
    │   ├── 降级? → 使用缓存/默认值
    │   └── 记录并继续
    │
    └── 不可恢复?
        │
        ├── 客户端错误? → 4xx 响应
        └── 服务器错误? → 5xx 响应 + 告警
```

---

**质量评级**: S (16KB)
**完成日期**: 2026-04-02
