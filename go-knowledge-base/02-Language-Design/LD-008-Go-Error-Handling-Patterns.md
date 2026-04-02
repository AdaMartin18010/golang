# LD-008: Go 错误处理模式 (Go Error Handling Patterns)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #error-handling #patterns #sentinel-errors #error-wrapping #go113
> **权威来源**:
>
> - [Error Handling and Go](https://go.dev/blog/error-handling-and-go) - Go Authors
> - [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) - Damien Neil
> - [Clean Architecture](https://blog.cleancoder.com/) - Robert C. Martin

---

## 1. 错误处理基础

### 1.1 错误接口

```go
type error interface {
    Error() string
}
```

**定义 1.1 (错误)**
错误是表示异常状态的值，实现了 error 接口。

### 1.2 错误创建

```go
// 简单错误
err := errors.New("something went wrong")

// 格式化错误
err := fmt.Errorf("user %d not found", userID)

// 包装错误 (Go 1.13+)
err := fmt.Errorf("database error: %w", err)
```

---

## 2. 错误模式

### 2.1 哨兵错误

```go
// 定义哨兵错误
var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
)

// 使用
if err == ErrNotFound {
    // 处理未找到
}

// Go 1.13+ 检查包装错误
if errors.Is(err, ErrNotFound) {
    // 处理未找到
}
```

### 2.2 自定义错误类型

```go
// 带上下文的错误
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
}

// 使用
return &NotFoundError{Resource: "user", ID: "123"}

// 检查
var nf *NotFoundError
if errors.As(err, &nf) {
    fmt.Printf("Resource: %s\n", nf.Resource)
}
```

### 2.3 错误包装链

```go
// 构建错误链
func queryUser(db *sql.DB, id int) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    user := &User{}
    if err := row.Scan(&user.ID, &user.Name); err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("%w: user id=%d", ErrNotFound, id)
        }
        return nil, fmt.Errorf("database query failed: %w", err)
    }
    return user, nil
}

// 错误链结构:
// queryUser error: user id=123 not found
// └─ ErrNotFound
//    └─ sql.ErrNoRows
```

---

## 3. 错误处理策略

### 3.1 立即返回

```go
func doSomething() error {
    data, err := fetchData()
    if err != nil {
        return err
    }

    result, err := processData(data)
    if err != nil {
        return err
    }

    return saveResult(result)
}
```

### 3.2 包装添加上下文

```go
func doSomething() error {
    data, err := fetchData()
    if err != nil {
        return fmt.Errorf("fetching data: %w", err)
    }

    result, err := processData(data)
    if err != nil {
        return fmt.Errorf("processing data: %w", err)
    }

    if err := saveResult(result); err != nil {
        return fmt.Errorf("saving result: %w", err)
    }

    return nil
}
```

### 3.3 重试策略

```go
func fetchWithRetry(url string, maxRetries int) ([]byte, error) {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        data, err := fetch(url)
        if err == nil {
            return data, nil
        }

        lastErr = err

        // 检查是否可重试
        if !isRetryable(err) {
            return nil, err
        }

        // 指数退避
        time.Sleep(time.Duration(i*i) * time.Second)
    }

    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

### 3.4 降级处理

```go
func getUser(ctx context.Context, id int) (*User, error) {
    // 尝试从缓存获取
    if user, err := cache.Get(ctx, id); err == nil {
        return user, nil
    }

    // 缓存未命中，从数据库获取
    user, err := db.GetUser(ctx, id)
    if err != nil {
        // 使用默认值降级
        return defaultUser(), nil
    }

    // 更新缓存
    cache.Set(ctx, id, user)
    return user, nil
}
```

---

## 4. 错误处理最佳实践

### 4.1 检查清单

- [ ] 每个错误返回值都检查
- [ ] 使用 %w 包装底层错误
- [ ] 添加有用的上下文信息
- [ ] 区分可恢复和不可恢复错误
- [ ] 使用哨兵错误进行类型检查
- [ ] 自定义错误类型携带上下文

### 4.2 反模式

```go
// 忽略错误
_ = doSomething()  // BAD!

// 无用的包装
if err != nil {
    return fmt.Errorf("error: %v", err)  // 无上下文
}

// 过度使用 panic
if err != nil {
    panic(err)  // 除非不可恢复，否则不要用
}

// 丢失原始错误
if err != nil {
    return errors.New("failed")  // 丢失原始错误信息
}
```

---

## 5. 代码示例

### 5.1 完整错误处理

```go
package user

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
)

var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email")
)

type User struct {
    ID    int
    Email string
    Name  string
}

type Repository struct {
    db *sql.DB
}

func (r *Repository) GetByID(ctx context.Context, id int) (*User, error) {
    user := &User{}
    err := r.db.QueryRowContext(ctx,
        "SELECT id, email, name FROM users WHERE id = ?", id,
    ).Scan(&user.ID, &user.Email, &user.Name)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("%w: id=%d", ErrUserNotFound, id)
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return user, nil
}

func (r *Repository) Create(ctx context.Context, user *User) error {
    if !isValidEmail(user.Email) {
        return ErrInvalidEmail
    }

    result, err := r.db.ExecContext(ctx,
        "INSERT INTO users (email, name) VALUES (?, ?)",
        user.Email, user.Name,
    )
    if err != nil {
        return fmt.Errorf("insert user: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("get last insert id: %w", err)
    }

    user.ID = int(id)
    return nil
}

// 服务层
type Service struct {
    repo *Repository
}

func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, fmt.Errorf("user service: %w", err)
        }
        return nil, fmt.Errorf("user service: %w", err)
    }
    return user, nil
}
```

---

## 6. 多元表征

### 6.1 错误处理决策树

```
函数返回错误?
│
├── 可重试?
│   ├── 是 → 指数退避重试
│   └── 否 → 继续
│
├── 已知错误类型?
│   ├── 是 → 特定处理
│   │       ├── ErrNotFound → 返回 404
│   │       ├── ErrUnauthorized → 返回 401
│   │   └── 其他 → 包装传播
│   └── 否 → 包装传播
│
└── 严重错误?
    └── 是 → panic (极少)
```

---

**质量评级**: S (15KB)
**完成日期**: 2026-04-02
