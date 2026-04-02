# LD-008: Go 错误处理模式 (Go Error Handling Patterns)

> **维度**: Language Design
> **级别**: S (18+ KB)
> **标签**: #go-error #error-handling #error-wrapping #sentinel-errors
> **权威来源**: [Error handling and Go](https://go.dev/blog/error-handling-and-go), [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)

---

## 错误处理哲学

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Go Error Handling Philosophy                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  显式 > 隐式                                                                  │
│  简单 > 复杂                                                                  │
│  错误是值，不是异常                                                            │
│                                                                              │
│  比较:                                                                        │
│                                                                              │
│  Go (显式):                                                                  │
│  f, err := os.Open("file")                                                   │
│  if err != nil {                                                             │
│      return fmt.Errorf("open file: %w", err)                                 │
│  }                                                                           │
│                                                                              │
│  Java (异常):                                                                │
│  try {                                                                       │
│      File f = new File("file");                                              │
│  } catch (IOException e) {                                                   │
│      throw new RuntimeException("open file", e);                             │
│  }                                                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 错误处理模式

### 1. Sentinel Errors (哨兵错误)

```go
package mypkg

import "errors"

// 预定义错误
var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")
var ErrInvalidInput = errors.New("invalid input")

// 使用
func GetUser(id string) (*User, error) {
    user, ok := db[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

// 检查
user, err := GetUser("123")
if errors.Is(err, ErrNotFound) {
    // 处理未找到
}
```

### 2. Error Wrapping (错误包装)

```go
// Go 1.13+ 错误包装
import "fmt"

func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open %s: %w", path, err)
    }
    defer f.Close()

    if err := process(f); err != nil {
        return fmt.Errorf("process %s: %w", path, err)
    }
    return nil
}

// 解包检查
if errors.Is(err, os.ErrNotExist) {
    // 文件不存在
}

// 获取原始错误
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Path:", pathErr.Path)
}
```

### 3. 自定义错误类型

```go
// 结构化错误
type APIError struct {
    Status  int    // HTTP状态码
    Code    string // 错误代码
    Message string // 用户可读消息
    Detail  string // 调试详情
}

func (e *APIError) Error() string {
    return fmt.Sprintf("[%d:%s] %s", e.Status, e.Code, e.Message)
}

func (e *APIError) Unwrap() error {
    // 返回底层错误
    return nil
}

// 使用
return &APIError{
    Status:  http.StatusNotFound,
    Code:    "USER_NOT_FOUND",
    Message: "用户不存在",
    Detail:  fmt.Sprintf("user id: %s", id),
}
```

### 4. 错误码设计

```go
// 统一错误码体系
const (
    ErrCodeOK              = 0
    ErrCodeInvalidParam    = 10001
    ErrCodeUnauthorized    = 10002
    ErrCodeNotFound        = 10004
    ErrCodeInternal        = 50001
    ErrCodeServiceUnavail  = 50003
)

type CodedError struct {
    Code    int
    Message string
    Err     error
}

func (e *CodedError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *CodedError) Unwrap() error {
    return e.Err
}
```

---

## 错误处理最佳实践

### ✅ 推荐

```go
// 1. 一次检查，立即处理
if err != nil {
    return err
}

// 2. 添加上下文
if err != nil {
    return fmt.Errorf("processing order %s: %w", orderID, err)
}

// 3. 使用哨兵错误
var ErrRetryable = errors.New("retryable error")

if errors.Is(err, ErrRetryable) {
    // 重试逻辑
}

// 4. 结构化日志
log.Error().
    Err(err).
    Str("order_id", orderID).
    Msg("failed to process order")
```

### ❌ 避免

```go
// 1. 忽略错误
_ = doSomething() // 危险！

// 2. 过度包装
return fmt.Errorf("error: %w", fmt.Errorf("error: %w", err))

// 3. 使用 panic
if err != nil {
    panic(err) // 不要这样做
}

// 4. 字符串比较
if err.Error() == "not found" { // 不可靠
}
```

---

## 参考文献

1. [Error handling and Go](https://go.dev/blog/error-handling-and-go)
2. [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
3. [pkg/errors](https://github.com/pkg/errors) (历史包)
