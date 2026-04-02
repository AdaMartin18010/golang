# 错误处理模式 (Error Handling Patterns)

> **分类**: 工程与云原生
> **标签**: #error-handling #patterns #best-practices

---

## 错误包装

```go
import "fmt"

// 简单包装
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// 添加上下文
if err != nil {
    return fmt.Errorf("process user %d: %w", userID, err)
}

// 多层包装
if err != nil {
    err = fmt.Errorf("database query: %w", err)
    err = fmt.Errorf("fetch user %d: %w", userID, err)
    err = fmt.Errorf("handle request: %w", err)
    return err
}
```

---

## 错误类型判断

### errors.Is

```go
import "errors"

if errors.Is(err, sql.ErrNoRows) {
    // 处理未找到
}

// 自定义错误
var ErrUserNotFound = errors.New("user not found")

if errors.Is(err, ErrUserNotFound) {
    // ...
}
```

### errors.As

```go
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) {
    // 访问具体字段
    fmt.Println(pgErr.Code)
    fmt.Println(pgErr.Message)
}
```

---

## 哨兵错误模式

```go
package mypkg

import "errors"

// 定义哨兵错误
var (
    ErrNotFound    = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInternal     = errors.New("internal error")
)

// 使用
func FindUser(id string) (*User, error) {
    user, err := db.Query(id)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user %s: %w", id, ErrNotFound)
    }
    return user, err
}
```

---

## 结构化错误

```go
type AppError struct {
    Code    string
    Message string
    Details map[string]interface{}
    Cause   error
}

func (e *AppError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Cause
}

// 创建错误
func NewNotFoundError(resource string, id string) *AppError {
    return &AppError{
        Code:    "NOT_FOUND",
        Message: fmt.Sprintf("%s %s not found", resource, id),
        Details: map[string]interface{}{
            "resource": resource,
            "id":       id,
        },
    }
}

// HTTP 响应
func (e *AppError) HTTPStatus() int {
    switch e.Code {
    case "NOT_FOUND":
        return http.StatusNotFound
    case "INVALID_INPUT":
        return http.StatusBadRequest
    case "UNAUTHORIZED":
        return http.StatusUnauthorized
    default:
        return http.StatusInternalServerError
    }
}
```

---

## 错误聚合

```go
type MultiError struct {
    errors []error
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.errors = append(m.errors, err)
    }
}

func (m *MultiError) Error() string {
    var msgs []string
    for _, err := range m.errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (m *MultiError) HasErrors() bool {
    return len(m.errors) > 0
}

// 使用
func validate(user *User) error {
    var m MultiError

    if user.Name == "" {
        m.Add(errors.New("name is required"))
    }
    if user.Email == "" {
        m.Add(errors.New("email is required"))
    }
    if !isValidEmail(user.Email) {
        m.Add(errors.New("email is invalid"))
    }

    if m.HasErrors() {
        return &m
    }
    return nil
}
```

---

## 重试模式

```go
func Retry(fn func() error, attempts int, delay time.Duration) error {
    var err error

    for i := 0; i < attempts; i++ {
        if err = fn(); err == nil {
            return nil
        }

        // 不要重试特定错误
        if errors.Is(err, ErrInvalidInput) {
            return err
        }

        if i < attempts-1 {
            time.Sleep(delay * time.Duration(i+1))
        }
    }

    return fmt.Errorf("after %d attempts: %w", attempts, err)
}
```

---

## 降级模式

```go
func GetDataWithFallback(ctx context.Context, key string) (Data, error) {
    // 先尝试缓存
    data, err := cache.Get(ctx, key)
    if err == nil {
        return data, nil
    }

    // 缓存未命中，尝试数据库
    data, err = db.Get(ctx, key)
    if err == nil {
        // 回填缓存
        cache.Set(ctx, key, data, time.Hour)
        return data, nil
    }

    // 数据库失败，使用默认值
    return DefaultData(), nil
}
```
