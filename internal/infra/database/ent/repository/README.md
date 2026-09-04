# Ent Repository 实现

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 概述

本包提供基于 Ent ORM 的通用仓储实现，实现框架定义的 `Repository` 接口。

---

## 🎯 设计原理

### 1. 架构位置

- **位置**: Infrastructure Layer (`internal/infrastructure/database/ent/repository/`)
- **职责**: Ent 仓储实现
- **依赖**: Ent ORM、Domain Layer 接口

### 2. 设计原则

- **依赖倒置**: 实现领域层定义的接口，不依赖具体实现
- **泛型设计**: 使用泛型提供类型安全的实现
- **实体转换**: 在领域实体和 Ent 实体之间进行转换
- **事务支持**: 提供事务管理功能

---

## 📦 组件说明

### BaseRepository

`BaseRepository` 是基础仓储实现，提供通用的 CRUD 操作框架。

**类型参数**:
- `T`: 领域实体类型（Domain Entity）
- `E`: Ent 实体类型（Ent Entity）

**功能**:
- ✅ 提供 CRUD 操作框架
- ✅ 实体转换支持
- ✅ 事务管理
- ✅ 错误处理

**注意**: 这是一个基础实现，用户需要在具体的仓储中实现业务特定的方法。

---

## 🚀 使用指南

### 1. 创建具体的仓储实现

```go
package repository

import (
    "context"

    "entgo.io/ent"
    "github.com/yourusername/golang/internal/domain/user"
    "github.com/yourusername/golang/internal/infrastructure/database/ent"
    "github.com/yourusername/golang/internal/infrastructure/database/ent/repository"
    entuser "github.com/yourusername/golang/internal/infrastructure/database/ent/user"
)

// UserRepository 用户仓储实现
type UserRepository struct {
    *repository.BaseRepository[*user.User, *entuser.User]
    client *ent.Client
}

// NewUserRepository 创建用户仓储
func NewUserRepository(client *ent.Client) *UserRepository {
    return &UserRepository{
        BaseRepository: repository.NewBaseRepository(
            client,
            toDomainUser,  // Ent 实体转领域实体
            toEntUser,     // 领域实体转 Ent 实体
            getUserID,     // 获取用户 ID
            setUserID,     // 设置用户 ID
        ),
        client: client,
    }
}

// toDomainUser 将 Ent 用户实体转换为领域用户实体
func toDomainUser(entUser *entuser.User) (*user.User, error) {
    return &user.User{
        ID:        entUser.ID,
        Email:     entUser.Email,
        Name:      entUser.Name,
        CreatedAt: entUser.CreatedAt,
        UpdatedAt: entUser.UpdatedAt,
    }, nil
}

// toEntUser 将领域用户实体转换为 Ent 用户实体
func toEntUser(domainUser *user.User) (*entuser.User, error) {
    builder := entuser.Create().
        SetEmail(domainUser.Email).
        SetName(domainUser.Name)

    if domainUser.ID != "" {
        builder.SetID(domainUser.ID)
    }

    return builder, nil
}

// getUserID 获取用户 ID
func getUserID(domainUser *user.User) (string, error) {
    return domainUser.ID, nil
}

// setUserID 设置用户 ID
func setUserID(domainUser *user.User, id string) error {
    domainUser.ID = id
    return nil
}

// Create 创建用户（重写基础方法）
func (r *UserRepository) Create(ctx context.Context, entity *user.User) error {
    entUser, err := toEntUser(entity)
    if err != nil {
        return err
    }

    created, err := r.client.User.Create().
        SetEmail(entity.Email).
        SetName(entity.Name).
        Save(ctx)
    if err != nil {
        return handleEntError(err)
    }

    // 设置 ID 和时间戳
    entity.ID = created.ID
    entity.CreatedAt = created.CreatedAt
    entity.UpdatedAt = created.UpdatedAt

    return nil
}

// FindByID 根据 ID 查找用户（重写基础方法）
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
    entUser, err := r.client.User.Get(ctx, id)
    if err != nil {
        return nil, handleEntError(err)
    }

    return toDomainUser(entUser)
}

// FindByEmail 根据邮箱查找用户（业务特定方法）
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    entUser, err := r.client.User.Query().
        Where(entuser.EmailEQ(email)).
        Only(ctx)
    if err != nil {
        return nil, handleEntError(err)
    }

    return toDomainUser(entUser)
}
```

### 2. 使用仓储

```go
// 创建 Ent 客户端
client, err := ent.NewClientFromConfig(ctx, ...)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 创建仓储
userRepo := repository.NewUserRepository(client)

// 使用仓储
user := &user.User{
    Email: "test@example.com",
    Name:  "Test User",
}

err = userRepo.Create(ctx, user)
if err != nil {
    log.Fatal(err)
}

found, err := userRepo.FindByID(ctx, user.ID)
if err != nil {
    log.Fatal(err)
}
```

### 3. 使用事务

```go
// 在事务中执行多个操作
err := userRepo.WithTx(ctx, func(tx *ent.Tx) error {
    // 创建用户
    user := &user.User{Email: "test@example.com", Name: "Test"}
    if err := userRepo.Create(ctx, user); err != nil {
        return err
    }

    // 更新用户
    user.Name = "Updated"
    if err := userRepo.Update(ctx, user); err != nil {
        return err
    }

    return nil
})
```

---

## 📚 相关文档

- Ent ORM 文档
- [仓储接口定义](../../../../internal/domain/interfaces/repository.go)
- [框架数据库抽象](../../../../../pkg/database/README.md)

---

**最后更新**: 2025-01-XX
