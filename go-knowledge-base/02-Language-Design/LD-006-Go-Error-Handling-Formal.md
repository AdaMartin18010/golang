# LD-006: Go 错误处理的形式化理论与实践 (Go Error Handling: Formal Theory & Practice)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #error-handling #error-wrapping #sentinel-errors #go1.13
> **权威来源**:
>
> - [Error Handling and Go](https://go.dev/blog/error-handling-and-go) - Go Authors
> - [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) - Damien Neil
> - [Failure Handling in Distributed Systems](https://dl.acm.org/doi/10.1145/3335772.3336773) - SOSP 2019
> - [Clean Code](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350884) - Robert C. Martin

---

## 1. 形式化基础

### 1.1 错误处理的理论基础

**定义 1.1 (错误)**
错误是程序执行过程中偏离预期行为的任何事件。

**定义 1.2 (故障、错误、失效)**

- **故障 (Fault)**: 系统中的缺陷
- **错误 (Error)**: 故障的激活状态
- **失效 (Failure)**: 观察到的服务偏离

**定理 1.1 (错误传播)**
若组件 $A$ 调用组件 $B$，$B$ 的错误可能导致 $A$ 失效，除非 $A$ 正确处理 $B$ 的错误。

$$\text{Fault}_B \to \text{Error}_B \xrightarrow{\text{handle}} \text{No Failure}_A$$
$$\text{Fault}_B \to \text{Error}_B \xrightarrow{\text{no handle}} \text{Failure}_A$$

### 1.2 Go 错误处理哲学

**公理 1.1 (显式错误检查)**
错误必须显式处理，不可静默忽略。

**公理 1.2 (错误即值)**
错误是普通的值，非异常控制流。

---

## 2. Go 错误接口的形式化

### 2.1 错误接口定义

**定义 2.1 (error 接口)**

```go
type error interface {
    Error() string
}
```

形式化：
$$\text{error} = \{ \text{Error}(): \text{string} \}$$

**定义 2.2 (错误相等)**
两个错误相等当且仅当：
$$e_1 = e_2 \Leftrightarrow e_1.\text{Error}() = e_2.\text{Error}()$$

注意：这仅是字符串相等，语义可能不等价。

### 2.2 错误包装 (Go 1.13+)

**定义 2.3 (包装错误)**

```go
func Wrap(err error, msg string) error {
    return fmt.Errorf("%s: %w", msg, err)
}
```

**定义 2.4 (解包接口)**

```go
type unwrapper interface {
    Unwrap() error
}
```

**定理 2.1 (错误链)**
包装错误形成链表结构：
$$e_n \to e_{n-1} \to ... \to e_1 \to \text{nil}$$

其中 $e_i$ 是第 $i$ 层包装的错误。

---

## 3. 错误处理模式

### 3.1 哨兵错误 (Sentinel Errors)

**定义 3.1 (哨兵错误)**
预定义的错误变量，用于错误类型判断：

```go
var ErrNotFound = errors.New("not found")
var ErrPermissionDenied = errors.New("permission denied")
```

**定理 3.1 (哨兵错误比较)**
$$\text{err} == \text{ErrNotFound} \Leftrightarrow \text{err 是 NotFound 类型}$$

### 3.2 自定义错误类型

**定义 3.2 (结构体错误)**

```go
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
}
```

**定理 3.2 (类型断言判断)**
$$\text{err}.(*\text{NotFoundError}) \text{ ok} \Leftrightarrow \text{err 是 NotFoundError 类型}$$

### 3.3 错误处理对比矩阵

| 模式 | 适用场景 | 性能 | 可扩展性 | 版本要求 |
|------|----------|------|----------|----------|
| 哨兵错误 | 简单错误类型判断 | 快 | 低 | 全版本 |
| 类型断言 | 需要错误上下文 | 中 | 中 | 全版本 |
| errors.Is | 包装错误检查 | 中 | 高 | Go 1.13+ |
| errors.As | 提取特定错误类型 | 中 | 高 | Go 1.13+ |
| fmt.Errorf %w | 添加上下文 | 中 | 高 | Go 1.13+ |

---

## 4. 多元表征

### 4.1 错误处理决策树

```
函数返回错误?
│
├── 否 → 正常执行
│
└── 是 → 错误类型?
    │
    ├── 可恢复?
    │   ├── 是 → 重试/降级
    │   │       │
    │   │       ├── 临时错误?
    │   │       │   └── 是 → 指数退避重试
    │   │       │
    │   │       └── 降级方案?
    │   │           └── 是 → 返回默认值/缓存
    │   │
    │   └── 否 → 包装错误向上传播
    │           fmt.Errorf("context: %w", err)
    │
    ├── 已知类型?
    │   ├── 是 → 特定处理
    │   │       │
    │   │       ├── errors.Is(err, ErrNotFound)
    │   │       │   └── 返回 404
    │   │       │
    │   │       ├── errors.Is(err, ErrPermissionDenied)
    │   │       │   └── 返回 403
    │   │       │
    │   │       └── var nf *NotFoundError
    │   │           errors.As(err, &nf)
    │   │           └── 使用 nf.Resource, nf.ID
    │   │
    │   └── 否 → 记录并包装
    │           log.Printf("unexpected: %v", err)
    │           return fmt.Errorf("operation failed: %w", err)
    │
    └── 严重错误?
        └── 是 → panic (极少使用)
                └── 仅用于不可恢复状态: 初始化失败、不变量违反
```

### 4.2 错误包装层次图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Error Wrapping Hierarchy                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Level 4: 用户展示层                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ "无法加载用户资料"                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│                                    ▼ %w                                     │
│  Level 3: 业务逻辑层                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ "查询用户失败: 用户ID=123"                                              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│                                    ▼ %w                                     │
│  Level 2: 数据访问层                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ "数据库查询失败: SELECT * FROM users"                                   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│                                    ▼ %w                                     │
│  Level 1: 底层错误                                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ sql.ErrConnDone (哨兵错误) 或 *net.OpError                             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  遍历: for err != nil { err = errors.Unwrap(err) }                          │
│  检查: if errors.Is(err, sql.ErrConnDone)                                   │
│  提取: var opErr *net.OpError; errors.As(err, &opErr)                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 错误处理策略对比

| 策略 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| **立即返回** | 简单、无嵌套 | 丢失上下文 | 简单函数 |
| **包装返回** | 保留上下文链 | 堆分配 | 标准做法 |
| **聚合错误** | 报告所有问题 | 复杂 | 验证场景 |
| **重试+退避** | 处理临时故障 | 延迟增加 | 网络调用 |
| **降级** | 保持可用 | 功能受限 | 高可用系统 |

### 4.4 Go vs 其他语言错误处理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Error Handling Paradigms Comparison                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Go: 显式错误返回值                                                          │
│  ───────────────────────────────────────────────────────────────────────     │
│  result, err := operation()                                                 │
│  if err != nil {                                                            │
│      return fmt.Errorf("context: %w", err)                                  │
│  }                                                                           │
│  ✓ 显式、可预测                                                              │
│  ✗ 代码冗长                                                                  │
│                                                                              │
│  Java/C#: 异常                                                               │
│  ───────────────────────────────────────────────────────────────────────     │
│  try {                                                                       │
│      result = operation();                                                  │
│  } catch (SpecificException e) {                                            │
│      throw new WrapperException("context", e);                              │
│  }                                                                           │
│  ✓ 错误处理与业务分离                                                         │
│  ✗ 隐藏控制流、性能开销                                                       │
│                                                                              │
│  Rust: Result 类型                                                           │
│  ───────────────────────────────────────────────────────────────────────     │
│  let result = operation()?;  // 传播                                         │
│  // 或                                                                       │
│  match operation() {                                                        │
│      Ok(v) => v,                                                            │
│      Err(e) => return Err(e.into()),                                        │
│  }                                                                           │
│  ✓ 类型安全、零开销                                                           │
│  ✗ 语法学习曲线                                                               │
│                                                                              │
│  Zig: Error Union                                                            │
│  ───────────────────────────────────────────────────────────────────────     │
│  const result = try operation();                                            │
│  ✓ 简洁、显式                                                                │
│  ✗ 生态系统成熟度                                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. 错误处理最佳实践

### 5.1 错误包装原则

**原则 1: 添加上下文，不重复**

```go
// 坏
return fmt.Errorf("error: %v", err)

// 好
return fmt.Errorf("查询用户 %d 失败: %w", userID, err)
```

**原则 2: 只在关键边界包装**

- 每个公共 API 边界
- 领域边界（如 repository → service）
- 不包装内部辅助函数的错误

### 5.2 完整错误处理示例

```go
package user

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"
)

// Sentinel errors
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrDuplicateEmail    = errors.New("email already exists")
    ErrInvalidInput      = errors.New("invalid input")
)

// Custom error type for detailed errors
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Repository interface
type Repository interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    Create(ctx context.Context, user *User) error
}

// Service implements business logic
type Service struct {
    repo   Repository
    cache  Cache
    logger Logger
}

// GetUser retrieves user with caching and error handling
func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
    // Try cache first
    if user, err := s.cache.Get(ctx, id); err == nil {
        return user, nil
    }

    // Query database with retry
    var user *User
    var err error

    for attempt := 0; attempt < 3; attempt++ {
        user, err = s.repo.GetByID(ctx, id)
        if err == nil {
            break
        }

        // Check if it's a transient error
        var netErr net.Error
        if errors.As(err, &netErr) && netErr.Temporary() {
            time.Sleep(time.Duration(attempt) * time.Second)
            continue
        }

        // Non-retryable error
        break
    }

    if err != nil {
        // Check specific error types
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("%w: id=%d", ErrUserNotFound, id)
        }

        // Log unexpected errors
        s.logger.Error("failed to get user",
            "id", id,
            "error", err,
            "attempts", 3)

        return nil, fmt.Errorf("database error: %w", err)
    }

    // Update cache (best effort)
    if cacheErr := s.cache.Set(ctx, id, user); cacheErr != nil {
        s.logger.Warn("failed to update cache", "error", cacheErr)
    }

    return user, nil
}

// CreateUser validates and creates a new user
func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // Validation
    if err := validateCreateInput(input); err != nil {
        return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
    }

    user := &User{
        Email:     input.Email,
        Name:      input.Name,
        CreatedAt: time.Now(),
    }

    if err := s.repo.Create(ctx, user); err != nil {
        // Check for duplicate key
        if isDuplicateKeyError(err) {
            return nil, fmt.Errorf("%w: %s", ErrDuplicateEmail, input.Email)
        }
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

func validateCreateInput(input CreateUserInput) error {
    if input.Email == "" {
        return &ValidationError{Field: "email", Message: "required"}
    }
    if input.Name == "" {
        return &ValidationError{Field: "name", Message: "required"}
    }
    return nil
}

func isDuplicateKeyError(err error) bool {
    // Database-specific check
    var pgErr *pq.Error
    if errors.As(err, &pgErr) {
        return pgErr.Code == "23505" // unique_violation
    }
    return false
}
```

---

## 6. 关系网络

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Go Error Handling Context                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  语言对比                                                                    │
│  ├── C: 错误码返回值 (errno)                                                │
│  ├── Java: Checked/Unchecked Exceptions                                     │
│  ├── Rust: Result<T, E> + ? operator                                        │
│  ├── Zig: Error unions                                                      │
│  └── Swift: throws / try                                                    │
│                                                                              │
│  Go 演进                                                                     │
│  ├── Go 1.0-1.12: 基础 error 接口                                           │
│  ├── Go 1.13: errors.Is, errors.As, %w                                      │
│  ├── Go 1.20: errors.Join (多错误聚合)                                       │
│  └── 未来: 可能的 try 语法糖 (讨论中)                                        │
│                                                                              │
│  相关模式                                                                    │
│  ├── Circuit Breaker (断路器)                                               │
│  ├── Retry with Backoff                                                     │
│  ├── Timeout and Cancellation (context)                                     │
│  └── Bulkhead (舱壁隔离)                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Go Error Handling Toolkit                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心原则                                                                    │
│  ═══════════════════════════════════════════════════════════════════════     │
│  1. 错误是值 - 不是控制流异常                                                │
│  2. 显式检查 - 不忽略错误返回值                                               │
│  3. 上下文包装 - 添加上下文，不重复                                           │
│  4. 哨兵 + 自定义类型 - 灵活的错误识别                                        │
│                                                                              │
│  检查清单:                                                                   │
│  □ 每个错误返回值都检查了吗?                                                   │
│  □ 错误信息是否清晰有用?                                                      │
│  □ 是否使用了 %w 包装底层错误?                                                │
│  □ 是否检查了特定的哨兵错误?                                                   │
│  □ 临时错误是否实现了重试?                                                    │
│  □ 是否记录了意外的错误?                                                      │
│                                                                              │
│  错误设计决策:                                                               │
│                                                                              │
│  使用哨兵错误 (var ErrX = errors.New(...)) 当:                               │
│  • 错误类型简单，不需要额外上下文                                              │
│  • 需要跨包比较错误类型                                                        │
│                                                                              │
│  使用自定义类型 (type MyError struct{...}) 当:                               │
│  • 错误需要携带额外上下文                                                      │
│  • 调用方可能需要根据错误字段做不同处理                                        │
│                                                                              │
│  包装策略:                                                                   │
│  • 公共 API 边界: 必须包装，添加上下文                                         │
│  • 内部函数: 通常直接返回                                                      │
│  • 数据库/外部服务: 包装并可能重试                                             │
│                                                                              │
│  常见反模式:                                                                 │
│  ❌ if err != nil { return nil, err } // 丢失上下文                           │
│  ❌ log.Printf("%v", err); return err // 双重记录                             │
│  ❌ panic(err) // 不必要的 panic                                             │
│  ❌ _ = operation() // 忽略错误!                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16KB)
**完成日期**: 2026-04-02
