# Clean Architecture 2025-2026 最佳实践补充

> **文档类型**: 补充更新 (Supplementary Update)
> **更新日期**: 2026-04-01
> **基于**: 网络最新权威资料 + 项目实际代码分析

---

## 1. 接口设计最佳实践 (2026)

### 1.1 小接口原则 (Small Interface Principle)

**反模式**: 胖接口 (Fat Interface)

```go
// ❌ 不推荐: 包含过多方法的接口
type UserRepository interface {
    Create(user User) error
    GetByID(id string) (User, error)
    GetByEmail(email string) (User, error)
    Update(user User) error
    Delete(id string) error
    List(offset, limit int) ([]User, error)
    Count() (int, error)
    ValidatePassword(user User, password string) bool
    UpdatePassword(user User, password string) error
    // ... 更多方法
}
```

**推荐**: 细粒度接口拆分

```go
// ✅ 推荐: 小而专注的接口
type UserReader interface {
    GetByID(id string) (User, error)
    GetByEmail(email string) (User, error)
    List(offset, limit int) ([]User, error)
}

type UserWriter interface {
    Create(user User) error
    Update(user User) error
    Delete(id string) error
}

type UserAuthenticator interface {
    ValidatePassword(user User, password string) bool
    UpdatePassword(user User, password string) error
}

// 组合使用
type UserRepository interface {
    UserReader
    UserWriter
}
```

### 1.2 接口命名规范

**反模式**: 实现类后缀 `Impl`

```go
// ❌ 不推荐
 type UserRepositoryImpl struct{}
 type UserServiceImpl struct{}
```

**推荐**: 描述性实现名称

```go
// ✅ 推荐
type PostgresUserRepository struct{}
type InMemoryUserRepository struct{}
type RedisCachedUserRepository struct{}
```

---

## 2. 测试金字塔策略 (2026)

### 2.1 测试比例目标

| 测试层级 | 覆盖率目标 | 执行时间 | 工具推荐 |
|----------|-----------|----------|----------|
| **单元测试** | 70% | < 1秒 | 标准 `testing` 包 |
| **集成测试** | 20% | < 30秒 | Testcontainers |
| **E2E 测试** | 10% | 分钟级 | Playwright / 自定义 |

### 2.2 分层测试策略

**Domain Layer - 纯单元测试**

```go
func TestOrder_AddItem(t *testing.T) {
    // 无需 mock，测试纯业务逻辑
    order := NewOrder(CustomerID.New())

    err := order.AddItem(product, 2)

    require.NoError(t, err)
    assert.Equal(t, 2, order.Items[0].Quantity)
    assert.Equal(t, OrderStatusPending, order.Status)
}
```

**Use Case Layer - 集成测试**

```go
func TestUserUseCase_CreateUser(t *testing.T) {
    // 使用内存仓库或 testcontainer
    ctx := context.Background()
    repo := NewInMemoryUserRepository() // 或 testcontainer
    uc := NewUserUseCase(repo)

    user, err := uc.Create(ctx, CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    })

    require.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

**Repository Layer - 集成测试**

```go
func TestPostgresUserRepository_Create(t *testing.T) {
    // 使用 testcontainer 测试真实数据库
    ctx := context.Background()
    container := setupPostgresContainer(t)
    defer container.Terminate(ctx)

    repo := NewPostgresUserRepository(container.ConnectionString())

    err := repo.Create(ctx, &User{Email: "test@example.com"})

    require.NoError(t, err)
}
```

### 2.3 本项目测试优化成果

| 模块 | 优化前 | 优化后 | 方法 |
|------|--------|--------|------|
| Redis 缓存 | 62s (外部依赖) | 1.2s | miniredis |
| MQTT | 0% | 17.5% | 提取可测试函数 |
| NATS | 16.7% | 35.9% | 提取 marshalPayload |
| Redis | 44.4% | 94.4% | miniredis + 单元测试 |

---

## 3. Go 1.26 特性在 Clean Architecture 中的应用

### 3.1 errors.AsType - 类型安全错误处理

**旧方式**:

```go
var myErr *MyError
if errors.As(err, &myErr) {
    // 需要预先声明变量
    handle(myErr)
}
```

**Go 1.26 新方式**:

```go
// Domain Layer - 错误定义
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Application Layer - 类型安全处理
if valErr, ok := errors.AsType[*ValidationError](err); ok {
    // 直接在 if 中使用，无需预声明
    log.Printf("Validation failed: field=%s, msg=%s", valErr.Field, valErr.Message)
    return ErrInvalidInput
}
```

### 3.2 new() 表达式 - 简化指针创建

**Use Case 层 DTO 创建**:

```go
// 旧方式
timeout := 30 * time.Second
config := &Config{
    Timeout: &timeout,  // 需要中间变量
}

// Go 1.26 新方式
config := &Config{
    Timeout: new(30 * time.Second),  // 直接表达式
}
```

**Domain Entity 创建**:

```go
// 创建带可选字段的实体
func NewUser(email string, age *int) *User {
    return &User{
        ID:    uuid.New().String(),
        Email: email,
        Age:   age,  // 可能为 nil
    }
}

// 调用
user := NewUser("test@example.com", new(25))  // 直接使用 new()
```

### 3.3 递归类型约束 - 通用 Repository 接口

```go
// Domain Layer - 可比较实体的通用接口
type Entity[E Entity[E]] interface {
    GetID() string
    SetID(id string)
    Clone() E  // 返回自身类型
}

// 具体实体实现
type User struct {
    ID    string
    Email string
}

func (u User) GetID() string { return u.ID }
func (u User) SetID(id string) { u.ID = id }
func (u User) Clone() User {
    return User{ID: u.ID, Email: u.Email}
}

// Repository 接口
type Repository[T Entity[T]] interface {
    Create(ctx context.Context, entity T) error
    GetByID(ctx context.Context, id string) (T, error)
    Update(ctx context.Context, entity T) error
}

// 具体实现
type PostgresRepository[T Entity[T]] struct {
    db *sql.DB
}

func (r *PostgresRepository[T]) Create(ctx context.Context, entity T) error {
    // 通用实现
}
```

---

## 4. 项目架构改进建议

### 4.1 当前架构评估

```
✅ 已正确实现:
- 四层架构清晰分离
- Domain 层无外部依赖
- 依赖方向向内
- 使用接口进行依赖倒置

⚠️ 可改进点:
- 部分接口可以更细粒度
- 测试金字塔可进一步优化
- 部分 Repository 可使用泛型基类
```

### 4.2 具体改进建议

1. **接口细化**:
   - 将 `Cache` 接口拆分为 `CacheReader` / `CacheWriter`
   - 将 `MessagePublisher` 拆分为 `Publisher` / `Subscriber`

2. **泛型 Repository 基类**:
   - 创建 `internal/infra/database/base_repository.go`
   - 实现通用的 CRUD 操作
   - 具体 Repository 嵌入基类

3. **错误处理统一**:
   - 使用 `errors.AsType` 替代 `errors.As`
   - 统一 Domain 错误类型定义

---

## 5. 常见反模式 checklist

| 反模式 | 检测方法 | 修复建议 |
|--------|----------|----------|
| 贫血模型 | Domain 只有 getter/setter | 添加业务方法到 Entity |
| 循环依赖 | 包互相导入 | 提取共享类型到独立包 |
| 基础设施泄漏 | Handler 直接使用 DB | 通过 Use Case 中转 |
| 跨层调用 | Handler 直接调用 Repository | 必须通过 Use Case |
| 同步类型复制 | sync.Mutex 值传递 | 使用指针接收者 |
| 胖接口 | 接口 > 5 个方法 | 拆分为小接口 |

---

## 6. 参考资料

1. [Clean Architecture in Go 2026 - Reintech](https://reintech.io/blog/go-clean-architecture-2026)
2. [Practical Guide to Clean Architecture - CyberAgent](https://developers.cyberagent.co.jp/blog/practical-clean-architecture-go/)
3. [Go Project Structure 2026 - Dasroot](https://dasroot.net/go-project-structure-2026/)
4. [Is Clean Architecture Overengineering? - Three Dots Tech](https://threedots.tech/post/is-clean-architecture-overengineering/)
