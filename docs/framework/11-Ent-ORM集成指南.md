# Ent ORM 集成指南

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 概述

本文档介绍如何在框架中使用 Ent ORM 进行数据访问，包括 Schema 定义、客户端创建、仓储实现等。

---

## 🎯 Ent ORM 简介

Ent 是 Facebook 开源的 Go 语言实体框架，提供：

- ✅ **类型安全**: 编译时类型检查
- ✅ **代码生成**: 自动生成类型安全的查询代码
- ✅ **迁移支持**: 自动数据库迁移
- ✅ **关系管理**: 支持复杂的关系映射
- ✅ **性能优化**: 高效的查询和更新

---

## 📦 框架集成

### 1. Ent 客户端

框架提供了便捷的 Ent 客户端创建函数：

```go
import (
    "context"
    "github.com/yourusername/golang/internal/infrastructure/database/ent"
)

// 创建 Ent 客户端
client, err := ent.NewClientFromConfig(
    ctx,
    "localhost",    // host
    "5432",         // port
    "postgres",     // user
    "password",     // password
    "mydb",         // dbname
    "disable",      // sslmode
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### 2. 数据库迁移

```go
// 运行数据库迁移
if err := client.Migrate(ctx); err != nil {
    log.Fatalf("Failed to run migrations: %v", err)
}
```

---

## 🚀 使用指南

### 1. 定义 Ent Schema

在用户项目中定义 Ent Schema：

```go
// schema/user.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/edge"
    "time"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("id").Unique().Immutable(),
        field.String("email").Unique().NotEmpty(),
        field.String("name").NotEmpty(),
        field.Time("created_at").Default(time.Now).Immutable(),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        // 定义关系
    }
}
```

### 2. 生成 Ent 代码

```bash
# 生成 Ent 代码
go generate ./ent
```

### 3. 实现仓储

使用框架提供的 `BaseRepository` 实现仓储：

```go
package repository

import (
    "context"
    "github.com/yourusername/golang/internal/domain/user"
    "github.com/yourusername/golang/internal/infrastructure/database/ent"
    "github.com/yourusername/golang/internal/infrastructure/database/ent/repository"
    entuser "github.com/yourusername/golang/internal/infrastructure/database/ent/user"
)

type UserRepository struct {
    *repository.BaseRepository[*user.User, *entuser.User]
    client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
    return &UserRepository{
        BaseRepository: repository.NewBaseRepository(
            client,
            toDomainUser,
            toEntUser,
            getUserID,
            setUserID,
        ),
        client: client,
    }
}

// 实现实体转换方法
func toDomainUser(entUser *entuser.User) (*user.User, error) {
    return &user.User{
        ID:        entUser.ID,
        Email:     entUser.Email,
        Name:      entUser.Name,
        CreatedAt: entUser.CreatedAt,
        UpdatedAt: entUser.UpdatedAt,
    }, nil
}

// 实现 CRUD 方法
func (r *UserRepository) Create(ctx context.Context, entity *user.User) error {
    created, err := r.client.User.Create().
        SetEmail(entity.Email).
        SetName(entity.Name).
        Save(ctx)
    if err != nil {
        return handleEntError(err)
    }

    entity.ID = created.ID
    entity.CreatedAt = created.CreatedAt
    entity.UpdatedAt = created.UpdatedAt

    return nil
}
```

### 4. 使用仓储

```go
// 在应用服务中使用仓储
type UserService struct {
    repo user.Repository
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*user.User, error) {
    u := &user.User{
        Email: email,
        Name:  name,
    }

    if err := s.repo.Create(ctx, u); err != nil {
        return nil, err
    }

    return u, nil
}
```

---

## 🔧 高级功能

### 1. 复杂查询

```go
// 使用 Ent 的查询构建器
users, err := client.User.Query().
    Where(
        user.EmailContains("@example.com"),
        user.CreatedAtGT(time.Now().AddDate(0, -1, 0)),
    ).
    Order(ent.Desc(user.FieldCreatedAt)).
    Limit(10).
    All(ctx)
```

### 2. 关系查询

```go
// 查询用户及其关联数据
user, err := client.User.Query().
    WithPosts().  // 预加载关联的 posts
    WithProfile(). // 预加载关联的 profile
    Where(user.IDEQ(userID)).
    Only(ctx)
```

### 3. 事务

```go
// 使用仓储的事务方法
err := userRepo.WithTx(ctx, func(tx *ent.Tx) error {
    // 在事务中执行多个操作
    user := &user.User{Email: "test@example.com"}
    if err := userRepo.Create(ctx, user); err != nil {
        return err
    }

    // 其他操作...
    return nil
})
```

---

## 📚 相关文档

- [Ent 官方文档](https://entgo.io/)
- [框架仓储接口](../../internal/domain/interfaces/repository.go)
- [Ent Repository 实现](../../internal/infrastructure/database/ent/repository/README.md)
- [框架数据库抽象](../../pkg/database/README.md)

---

**最后更新**: 2025-01-XX
