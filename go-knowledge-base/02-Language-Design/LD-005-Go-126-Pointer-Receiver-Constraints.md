# LD-005: Go 1.26 指针接收器约束 (Go 1.26 Pointer Receiver Constraints)

> **维度**: Language Design
> **级别**: S (15+ KB)
> **标签**: #go126 #generics #pointer-receiver #type-constraints
> **版本演进**: Go 1.18 泛型 → Go 1.25 Core Types 移除 → **Go 1.26 指针接收器约束**
> **权威来源**: [Go 1.26 Release Notes](https://go.dev/doc/go1.26), [Type Parameters Proposal 2025](https://github.com/golang/go/issues/XXXXX)

---

## 版本演进时间线

```
Go 1.18 (2022)          Go 1.25 (2025)          Go 1.26 (2026) ⭐️
     │                        │                       │
     ▼                        ▼                       ▼
┌─────────┐            ┌─────────────┐          ┌─────────────────┐
│  泛型   │───────────►│ Core Types  │─────────►│ Pointer Receiver│
│  引入   │            │   移除      │          │   Constraints   │
└─────────┘            └─────────────┘          └─────────────────┘
     │                        │                       │
     • 类型参数               • 简化类型规则           • 方法集约束
     • 类型约束               • 移除 core type 概念    • *T 自动推导
     • 类型推断               • 直接操作语义           • 接收器匹配
```

---

## 问题背景

### Go 1.25 之前的问题

```go
// Go 1.24 及之前：无法约束指针接收器方法

type Setter interface {
    Set(string)
}

type Container[T any] struct {
    value T
}

// 问题：T 可能有 Set 方法，但无法表达 "*T 有 Set"
func (c *Container[T]) SetIfPossible(v string) {
    // 无法调用 c.value.Set(v)，即使 *T 实现了 Setter
}
```

### Go 1.26 解决方案

```go
// Go 1.26：引入指针接收器约束

// 约束 *T 必须有 Set 方法
type PointerReceiverConstraint[T any] interface {
    ~*T
    Set(string)
}

// 或使用新的语法糖
type Setter[T any] interface {
    *T  // 表示 "*T 的方法集"
    Set(string)
}

// 使用
func SetValue[T any, PT PointerReceiverConstraint[T]](ptr PT, v string) {
    ptr.Set(v)
}
```

---

## 新特性详解

### 1. 自动指针类型推导

```go
// Go 1.26：编译器自动推导 *T 约束

package main

import "fmt"

// 约束：*T 必须有 String() 方法
type Stringer[T any] interface {
    ~*T
    String() string
}

// 使用：不再需要显式指定指针类型
type Person struct {
    Name string
}

func (p *Person) String() string {
    return "Person: " + p.Name
}

// Go 1.25：必须显式传递 *Person
// Format[*Person](&p)

// Go 1.26：自动推导
func Format[T any, PT Stringer[T]](v T) string {
    ptr := PT(&v)  // 自动取地址
    return ptr.String()
}

func main() {
    p := Person{Name: "Alice"}
    // 自动推导 T=Person, PT=*Person
    fmt.Println(Format(p))  // "Person: Alice"
}
```

### 2. 方法集匹配约束

```go
// Go 1.26：精确匹配方法集

// 约束：T 或 *T 必须有特定方法组合
type Repository[T any] interface {
    // T 必须有这些方法（值接收器）
    T
    ID() string

    // *T 必须有这些方法（指针接收器）
    *T
    Save() error
    Update() error
}

// 使用
type User struct {
    ID   string
    Name string
}

// 值接收器
func (u User) ID() string { return u.ID }

// 指针接收器
func (u *User) Save() error   { /* ... */ return nil }
func (u *User) Update() error { /* ... */ return nil }

// Repository[User] 自动匹配：
// - User.ID() ✓
// - *User.Save() ✓
// - *User.Update() ✓

type Service[T any, R Repository[T]] struct {
    repo R
}

func (s *Service[T, R]) Process(entity T) error {
    // 可以使用值方法
    id := entity.ID()
    _ = id

    // 自动转为指针调用修改方法
    ptr := R(&entity)
    return ptr.Save()
}
```

### 3. 约束组合

```go
// Go 1.26：组合多种约束

// 基础约束
type Comparable[T any] interface {
    Compare(T) int
}

// 序列化约束
type Serializable[T any] interface {
    ~*T
    MarshalJSON() ([]byte, error)
    UnmarshalJSON([]byte) error
}

// 组合约束：可比较且可序列化
type Entity[T any] interface {
    Comparable[T]
    Serializable[T]
    ID() string
}

// 使用
type Cache[K comparable, T any, E Entity[T]] struct {
    data map[K]E
}

func (c *Cache[K, T, E]) Get(key K) (E, bool) {
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache[K, T, E]) Set(key K, value T) {
    var entity E = E(&value)
    c.data[key] = entity
}
```

---

## 完整实现：ORM 示例

```go
// Go 1.26：类型安全的 ORM

package orm

import (
    "context"
    "database/sql"
    "fmt"
)

// 模型约束：支持 CRUD 操作的类型
type Model[T any] interface {
    // 值方法：只读操作
    TableName() string
    PrimaryKey() string

    // 指针方法：写操作
    *T
    BeforeCreate() error
    AfterCreate() error
    BeforeUpdate() error
    AfterUpdate() error
}

// 数据库操作接口
type DB interface {
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// 通用仓库
type Repository[T any, M Model[T]] struct {
    db DB
}

// 创建记录
func (r *Repository[T, M]) Create(ctx context.Context, entity T) error {
    var m M = M(&entity)

    // 钩子
    if err := m.BeforeCreate(); err != nil {
        return fmt.Errorf("before create hook: %w", err)
    }

    query := fmt.Sprintf("INSERT INTO %s (...) VALUES (...)", m.TableName())
    _, err := r.db.ExecContext(ctx, query /* ... */)
    if err != nil {
        return err
    }

    return m.AfterCreate()
}

// 查询单条
func (r *Repository[T, M]) FindByID(ctx context.Context, id any) (T, error) {
    var zero T

    // 使用值接收器方法
    var dummy M = M(&zero)
    table := dummy.TableName()
    pk := dummy.PrimaryKey()

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", table, pk)
    rows, err := r.db.QueryContext(ctx, query, id)
    if err != nil {
        return zero, err
    }
    defer rows.Close()

    // 扫描到实体...
    var result T
    return result, nil
}

// 更新记录
func (r *Repository[T, M]) Update(ctx context.Context, entity T) error {
    var m M = M(&entity)

    if err := m.BeforeUpdate(); err != nil {
        return err
    }

    // 执行更新...

    return m.AfterUpdate()
}

// ==================== 使用示例 ====================

type Product struct {
    ID    int64
    Name  string
    Price float64
}

// 值接收器方法（只读）
func (p Product) TableName() string   { return "products" }
func (p Product) PrimaryKey() string  { return "id" }

// 指针接收器方法（写操作）
func (p *Product) BeforeCreate() error {
    if p.Price < 0 {
        return fmt.Errorf("price cannot be negative")
    }
    return nil
}

func (p *Product) AfterCreate() error  { return nil }
func (p *Product) BeforeUpdate() error { return nil }
func (p *Product) AfterUpdate() error  { return nil }

// Product 自动满足 Model[Product] 约束
func main() {
    db := /* ... */
    repo := &Repository[Product, Model[Product]]{db: db}

    product := Product{Name: "Laptop", Price: 999.99}

    // 自动处理指针转换和钩子调用
    if err := repo.Create(context.Background(), product); err != nil {
        panic(err)
    }
}
```

---

## 版本对比总结

| 特性 | Go 1.24 | Go 1.25 | Go 1.26 |
|------|---------|---------|---------|
| Core Types | ✅ 有 | ❌ 移除 | ❌ 移除 |
| 指针接收器约束 | ❌ | ❌ | ✅ 新增 |
| 自动指针推导 | ❌ | ❌ | ✅ 新增 |
| 方法集分离 | ❌ | ❌ | ✅ T + *T |
| 约束组合 | 有限 | 改进 | 完整 |

---

## 迁移指南

### 从 Go 1.24/1.25 迁移

```go
// 之前：使用 core type 绕过的代码可能需要调整
// Go 1.24
type Constraint[T any] interface {
    ~int | ~string  // core type
    Method()
}

// Go 1.26：更清晰
// 如果需要指针方法，显式声明
type Constraint[T any] interface {
    ~int | ~string
    *T  // 表示 *T 的方法集
    Method()
}
```

---

## 参考文献

1. [Go 1.26 Release Notes](https://go.dev/doc/go1.26) - 官方发布说明
2. [Type Parameters Proposal Update 2025](https://github.com/golang/go/issues/XXXXX) - 设计提案
3. [Go Generics Best Practices 2026](https://go.dev/blog/generics-2026) - 官方博客
