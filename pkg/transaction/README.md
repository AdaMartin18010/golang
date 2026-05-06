# 事务管理框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26.2

---

## 📋 目录

- [事务管理框架](#事务管理框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心功能](#2-核心功能)
    - [2.1 Transaction 接口](#21-transaction-接口)
    - [2.2 Manager 接口](#22-manager-接口)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 手动管理事务](#32-手动管理事务)
    - [3.3 在Repository中使用](#33-在repository中使用)
    - [3.4 嵌套事务（使用现有事务）](#34-嵌套事务使用现有事务)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

事务管理框架提供了统一的事务管理功能：

- ✅ **事务接口**: 统一的事务接口
- ✅ **事务管理器**: 支持多种数据库的事务管理
- ✅ **Context集成**: 事务与Context集成
- ✅ **自动回滚**: 错误时自动回滚
- ✅ **嵌套事务支持**: 支持嵌套事务（通过context）

---

## 2. 核心功能

### 2.1 Transaction 接口

```go
type Transaction interface {
    Commit() error
    Rollback() error
    GetTx() interface{}
}
```

### 2.2 Manager 接口

```go
type Manager interface {
    Begin(ctx context.Context) (Transaction, error)
    Get(ctx context.Context) (Transaction, error)
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
    WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "context"
    "database/sql"
    "github.com/yourusername/golang/pkg/transaction"
)

// 创建事务管理器
db, _ := sql.Open("postgres", "...")
manager := transaction.NewSQLTransactionManager(db)

// 在事务中执行操作
err := manager.WithTransaction(ctx, func(ctx context.Context) error {
    // 获取事务
    tx, ok := transaction.GetSQLTx(ctx)
    if !ok {
        return errors.New("transaction not found")
    }

    // 执行数据库操作
    _, err := tx.Exec("INSERT INTO users (name) VALUES ($1)", "John")
    if err != nil {
        return err
    }

    _, err = tx.Exec("INSERT INTO orders (user_id) VALUES ($1)", userID)
    return err
})

if err != nil {
    // 事务已自动回滚
    log.Printf("Transaction failed: %v", err)
}
```

### 3.2 手动管理事务

```go
// 开始事务
tx, err := manager.Begin(ctx)
if err != nil {
    return err
}

// 将事务添加到context
ctx = context.WithValue(ctx, transactionKey{}, tx)

// 执行操作
sqlTx, _ := transaction.GetSQLTx(ctx)
_, err = sqlTx.Exec("INSERT INTO users (name) VALUES ($1)", "John")
if err != nil {
    tx.Rollback()
    return err
}

// 提交事务
err = tx.Commit()
if err != nil {
    return err
}
```

### 3.3 在Repository中使用

```go
type UserRepository struct {
    manager transaction.Manager
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
    return r.manager.WithTransaction(ctx, func(ctx context.Context) error {
        tx, ok := transaction.GetSQLTx(ctx)
        if !ok {
            return errors.New("transaction not found")
        }

        // 使用事务执行操作
        _, err := tx.Exec(
            "INSERT INTO users (name, email) VALUES ($1, $2)",
            user.Name, user.Email,
        )
        return err
    })
}
```

### 3.4 嵌套事务（使用现有事务）

```go
func (s *Service) CreateUserWithOrders(ctx context.Context, user *User, orders []Order) error {
    return s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
        // 创建用户
        if err := s.userRepo.CreateUser(ctx, user); err != nil {
            return err
        }

        // 创建订单（使用同一个事务）
        for _, order := range orders {
            if err := s.orderRepo.CreateOrder(ctx, user.ID, order); err != nil {
                return err // 自动回滚所有操作
            }
        }

        return nil
    })
}
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用WithTransaction**: 优先使用WithTransaction自动管理事务
2. **错误处理**: 确保返回错误以触发回滚
3. **Context传递**: 始终传递包含事务的context
4. **幂等操作**: 确保Rollback是幂等的

### 4.2 DON'Ts ❌

1. **不要忘记错误**: 忘记返回错误会导致事务不会回滚
2. **不要手动提交**: 在WithTransaction中不要手动提交
3. **不要忽略错误**: 始终检查事务操作的错误
4. **不要跨goroutine**: 事务不应该跨goroutine使用

---

## 5. 相关资源

- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
