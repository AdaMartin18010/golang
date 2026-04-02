# Go错误处理模式：形式化分析与最佳实践

> **版本**: 2026.04.01 | **Go版本**: 1.13-1.26.1 (errors包演进) | **模式**: 工程实践
> **关联**: [Go-1.26.1-Comprehensive.md](./Go-1.26.1-Comprehensive.md)

---

## 目录

- [Go错误处理模式：形式化分析与最佳实践](#go错误处理模式形式化分析与最佳实践)
  - [目录](#目录)
  - [1. Go错误处理哲学](#1-go错误处理哲学)
    - [1.1 显式优于隐式](#11-显式优于隐式)
    - [1.2 错误传播原则](#12-错误传播原则)
  - [2. 错误类型系统](#2-错误类型系统)
    - [2.1 error接口](#21-error接口)
    - [2.2 错误创建方式](#22-错误创建方式)
    - [2.3 错误语义分类](#23-错误语义分类)
  - [3. 错误包装与链](#3-错误包装与链)
    - [3.1 错误链（Go 1.13+）](#31-错误链go-113)
    - [3.2 错误包装形式化](#32-错误包装形式化)
    - [3.3 错误遍历](#33-错误遍历)
  - [4. 错误判断与处理策略](#4-错误判断与处理策略)
    - [4.1 策略模式](#41-策略模式)
    - [4.2 决策树](#42-决策树)
    - [4.3 sentinel模式](#43-sentinel模式)
  - [5. Panic与恢复模式](#5-panic与恢复模式)
    - [5.1 Panic语义](#51-panic语义)
    - [5.2 Recover模式](#52-recover模式)
    - [5.3 Panic vs Error](#53-panic-vs-error)
  - [6. 领域特定错误](#6-领域特定错误)
    - [6.1 HTTP错误映射](#61-http错误映射)
    - [6.2 多错误聚合](#62-多错误聚合)
    - [6.3 错误码设计](#63-错误码设计)
  - [7. 错误处理形式化](#7-错误处理形式化)
    - [7.1 错误单子（概念）](#71-错误单子概念)
    - [7.2 错误传播代数](#72-错误传播代数)
    - [7.3 可靠性定理](#73-可靠性定理)
  - [8. 反模式与陷阱](#8-反模式与陷阱)
    - [8.1 常见反模式](#81-常见反模式)
    - [8.2 错误检查疲劳](#82-错误检查疲劳)
  - [9. 性能考量](#9-性能考量)
    - [9.1 错误创建开销](#91-错误创建开销)
    - [9.2 错误比较性能](#92-错误比较性能)
  - [10. 未来演进](#10-未来演进)
    - [10.1 可能改进](#101-可能改进)
    - [10.2 错误处理风格](#102-错误处理风格)
  - [关联文档](#关联文档)

---

## 1. Go错误处理哲学

### 1.1 显式优于隐式

Go坚持**显式错误检查**原则：

```go
// Go风格：显式检查
result, err := DoSomething()
if err != nil {
    return err
}

// 对比异常风格（Go不推荐）
// try {
//     result = DoSomething()
// } catch (Exception e) {
//     throw e
// }
```

**哲学基础**:

- 错误是值（Errors are values）
- 调用者必须显式处理
- 代码路径清晰可见

### 1.2 错误传播原则

| 原则 | 描述 | 实践 |
|------|------|------|
| **快速失败** | 尽早返回错误 | 参数验证前置 |
| **不吞没** | 不忽略错误 | 每个err都有处理 |
| **添加上下文** | 包装错误信息 | `fmt.Errorf("...: %w", err)` |
| **类型区分** | 区分可恢复/不可恢复 | 类型断言判断 |

---

## 2. 错误类型系统

### 2.1 error接口

**定义**:

```go
type error interface {
    Error() string
}
```

**形式化**:

$$
\forall e. \; e : error \iff e \text{ implements } \text{Error}() : string
$$

### 2.2 错误创建方式

```go
// 方式1: errors.New（静态错误）
var ErrNotFound = errors.New("resource not found")

// 方式2: fmt.Errorf（动态错误）
err := fmt.Errorf("user %d not found", id)

// 方式3: 自定义类型（结构化错误）
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}
```

### 2.3 错误语义分类

| 类别 | 语义 | 处理策略 |
|------|------|---------|
| **输入错误** | 无效输入 | 返回400，不重试 |
| **瞬时错误** | 临时失败 | 重试，指数退避 |
| **状态错误** | 不一致状态 | 清理，返回500 |
| **权限错误** | 未授权 | 返回403/401 |
| **不存在** | 资源缺失 | 返回404 |

---

## 3. 错误包装与链

### 3.1 错误链（Go 1.13+）

```go
// 包装错误
if err != nil {
    return fmt.Errorf("database query failed: %w", err)
}

// 形成错误链
// sql.ErrNoRows
// ↓ wrap
// "database query failed: sql: no rows in result set"
// ↓ wrap
// "user lookup failed: database query failed: ..."
```

### 3.2 错误包装形式化

**定义 3.1 (错误链)**:

$$
e = \text{wrap}(m, e_{inner}) \Rightarrow \text{Unwrap}(e) = e_{inner}
$$

**链结构**:

```
e_n (最外层)
  ↓ Unwrap
e_{n-1}
  ↓ Unwrap
...
e_1 (根错误)
```

### 3.3 错误遍历

```go
// 遍历错误链
for err != nil {
    fmt.Println(err)
    err = errors.Unwrap(err)
}

// 使用errors.Is判断
if errors.Is(err, sql.ErrNoRows) {
    // 处理不存在
}

// 使用errors.As提取
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println(valErr.Field)
}
```

**形式化**:

$$
\text{Is}(e, target) \iff e = target \lor \exists e'. \text{Unwrap}(e) = e' \land \text{Is}(e', target)
$$

---

## 4. 错误判断与处理策略

### 4.1 策略模式

```go
// 策略1: 立即返回
if err != nil {
    return err
}

// 策略2: 包装返回
if err != nil {
    return fmt.Errorf("operation X: %w", err)
}

// 策略3: 转换错误类型
if err != nil {
    return &DomainError{Op: "X", Inner: err}
}

// 策略4: 重试
if err != nil {
    if retryable(err) {
        return retry()
    }
    return err
}

// 策略5: 降级
if err != nil {
    log.Error(err)
    return defaultValue, nil  // 降级处理
}
```

### 4.2 决策树

```
错误发生
├── 可重试？
│   ├── 是 → 重试（指数退避）
│   └── 否 → 继续
├── 可降级？
│   ├── 是 → 返回默认值
│   └── 否 → 继续
├── 需包装？
│   ├── 是 → 添加上下文
│   └── 否 → 原样返回
└── 致命？
    ├── 是 → panic
    └── 否 → 返回错误
```

### 4.3 sentinel模式

```go
// 定义哨兵错误
var (
    ErrNotFound    = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
    ErrTimeout      = errors.New("timeout")
)

// 使用
if errors.Is(err, ErrNotFound) {
    http.Error(w, "Not Found", 404)
    return
}
```

---

## 5. Panic与恢复模式

### 5.1 Panic语义

**定义 5.1 (Panic)**:

```go
// panic触发
panic(v interface{})

// 传播：沿调用栈向上展开
// 直到被recover捕获或程序终止
```

**形式化**:

$$
\text{panic}(v) \Rightarrow \text{unwind-stack}(v) \text{ until } \text{recover}()
$$

### 5.2 Recover模式

```go
// 防御性编程
func SafeCall() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
            log.Error("recovered from panic", r)
        }
    }()

    RiskyOperation()
    return nil
}
```

**正确使用场景**:

| 场景 | 使用 | 原因 |
|------|------|------|
| HTTP处理 | ✅ | 防止单个请求崩溃服务 |
| 后台任务 | ✅ | 保持服务运行 |
| 库代码 | ❌ | 让调用者决定 |
| 初始化 | ❌ | 应直接失败 |

### 5.3 Panic vs Error

| 维度 | Panic | Error |
|------|-------|-------|
| **语义** | 不可恢复 | 可处理 |
| **控制流** | 非本地跳转 | 正常返回 |
| **性能** | 昂贵（栈展开） | 廉价 |
| **使用频率** | 极少 | 常规 |

---

## 6. 领域特定错误

### 6.1 HTTP错误映射

```go
var errorToHTTPStatus = map[error]int{
    ErrNotFound:     404,
    ErrInvalidInput: 400,
    ErrUnauthorized: 401,
    ErrForbidden:    403,
    ErrTimeout:      504,
}

func HTTPStatus(err error) int {
    for e, status := range errorToHTTPStatus {
        if errors.Is(err, e) {
            return status
        }
    }
    return 500
}
```

### 6.2 多错误聚合

```go
// 验证多个字段
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    var msgs []string
    for _, e := range m.Errors {
        msgs = append(msgs, e.Error())
    }
    return strings.Join(msgs, "; ")
}

func Validate(user *User) error {
    var errs []error

    if user.Name == "" {
        errs = append(errs, &ValidationError{Field: "name", Message: "required"})
    }

    if user.Email == "" {
        errs = append(errs, &ValidationError{Field: "email", Message: "required"})
    }

    if len(errs) > 0 {
        return &MultiError{Errors: errs}
    }
    return nil
}
```

### 6.3 错误码设计

```go
type CodedError struct {
    Code    string
    Message string
    Details map[string]any
}

var (
    ErrUserNotFound = &CodedError{
        Code:    "USER_001",
        Message: "User not found",
    }
)
```

---

## 7. 错误处理形式化

### 7.1 错误单子（概念）

虽然Go不使用单子，但可类比：

```go
// 类似bind操作
func Bind[T, R any](v T, err error, f func(T) (R, error)) (R, error) {
    if err != nil {
        var zero R
        return zero, err
    }
    return f(v)
}

// 使用
result, err := Bind(GetUser(), err, func(u User) (Profile, error) {
    return GetProfile(u.ID)
})
```

### 7.2 错误传播代数

**操作定义**:

$$
\text{return}(e) = e
$$

$$
\text{wrap}(m, e) = m : e
$$

$$
\text{chain}(e_1, e_2) = \begin{cases}
e_1 & \text{if } e_1 \neq nil \\
e_2 & \text{otherwise}
\end{cases}
$$

### 7.3 可靠性定理

**定理 7.1**: 正确处理错误链保证故障可追溯。

$$
\forall e. \text{Chain}(e) \Rightarrow \exists \{e_i\}. e_n = e \land e_1 = \text{root}
$$

---

## 8. 反模式与陷阱

### 8.1 常见反模式

```go
// ❌ 忽略错误
_ = DoSomething()  // 错误被忽略！

// ✅ 显式忽略（有理由）
_ = DoSomething() // nolint:errcheck - 已知安全

// ❌ 过度包装
return fmt.Errorf("error: %w",
    fmt.Errorf("error: %w",
        fmt.Errorf("error: %w", err)))

// ❌ 丢失原始错误
if err != nil {
    return errors.New("failed")  // 原始err丢失
}

// ✅ 正确包装
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### 8.2 错误检查疲劳

```go
// 过多重复检查
v1, err := f1()
if err != nil {
    return err
}
v2, err := f2()
if err != nil {
    return err
}
// ... 重复

// 改进：使用辅助函数
func run(fn func() error) error {
    return fn()
}
```

---

## 9. 性能考量

### 9.1 错误创建开销

```go
// 基准测试
func BenchmarkErrorCreation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = errors.New("simple error")
    }
}

// 结果：~50ns/op
```

**优化**: 复用错误实例

```go
// 全局复用
var ErrNotFound = errors.New("not found")

// 对比每次创建
func Get(id string) (*Item, error) {
    if !exists(id) {
        return nil, ErrNotFound  // 复用
    }
    // ...
}
```

### 9.2 错误比较性能

```go
// errors.Is（遍历链）vs 直接比较

// 快：直接比较
if err == ErrNotFound { ... }

// 慢：遍历链
if errors.Is(err, ErrNotFound) { ... }
```

---

## 10. 未来演进

### 10.1 可能改进

| 提案 | 状态 | 影响 |
|------|------|------|
| try/catch语法 | 讨论中 | 有争议 |
| ?操作符 | 被拒绝 | 社区分歧 |
| 错误追溯增强 | 可能 | 提升调试 |
| 静态错误检查 | 工具支持 | linter增强 |

### 10.2 错误处理风格

```go
// 当前Go
result, err := Do()
if err != nil {
    return err
}

// 可能的未来（非官方）
result := Do() ?
// 或
result := Do() or return err
```

**社区立场**: 当前显式风格仍占主导。

---

## 关联文档

- [Go-1.26.1-Comprehensive.md](./Go-1.26.1-Comprehensive.md)
- [Go-Testing-Framework-Formalization.md](./Go-Testing-Framework-Formalization.md)

---

*文档版本: 2026-04-01 | Go风格: 显式错误处理 | 哲学: Errors are values*
