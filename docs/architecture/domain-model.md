# 领域模型设计

## 用户领域（User Domain）

### 实体（Entity）

```go
type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 仓储接口（Repository Interface）

```go
type Repository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*User, error)
}
```

### 领域服务（Domain Service）

```go
type DomainService interface {
    ValidateEmail(email string) bool
    IsEmailUnique(ctx context.Context, email string) (bool, error)
}
```

### 领域错误（Domain Errors）

- `ErrUserNotFound` - 用户不存在
- `ErrUserAlreadyExists` - 用户已存在
- `ErrInvalidEmail` - 无效邮箱

---

## 设计原则

1. **实体独立性** - 实体不依赖外部框架
2. **接口定义** - 仓储接口在领域层定义
3. **业务规则** - 业务规则封装在实体和领域服务中
4. **错误处理** - 领域错误在领域层定义
