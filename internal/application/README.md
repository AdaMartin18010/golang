# Application Layer (应用层)

Clean Architecture 的应用层，提供框架层面的应用服务抽象和通用应用模式。

## ⚠️ 重要说明

**本框架是框架性代码，不包含具体业务逻辑**：

- ✅ 提供通用的应用模式（Command、Query、Event、DTO 等）
- ✅ 用户可以通过这些模式实现自己的应用服务
- ❌ **不包含具体业务用例**（如用户管理、订单管理等）
- ❌ **不包含具体业务逻辑**

**用户使用框架时**：

1. 使用框架提供的应用模式（Command、Query、Event）
2. 实现自己的应用服务
3. 参考 `examples/framework-usage/` 中的示例代码

**注意**：如果存在 `user/`、`order/` 等目录，这些是**示例代码**，仅用于展示框架的使用方式，不是框架的核心部分。

## 结构

```text
application/
├── patterns/      # 通用应用模式
│   ├── command.go    # 命令模式
│   ├── query.go      # 查询模式
│   ├── event.go      # 事件模式
│   └── dto.go        # DTO 基类
└── workflow/      # 工作流模式（Temporal）
    └── context.go    # 工作流上下文
```

## 规则

- ✅ 只能导入 domain 层
- ✅ 提供框架层面的应用模式抽象
- ✅ 不包含具体业务用例
- ❌ 不能导入 infrastructure 或 interfaces 层
- ❌ 不包含具体业务应用服务

## 框架应用模式

### Command 模式

框架提供命令模式的抽象，用于处理写操作。

```go
// Command 命令接口（示例）
type Command interface {
    Execute(ctx context.Context) error
}

// CommandHandler 命令处理器接口（示例）
type CommandHandler[T Command] interface {
    Handle(ctx context.Context, cmd T) error
}
```

### Query 模式

框架提供查询模式的抽象，用于处理读操作。

```go
// Query 查询接口（示例）
type Query interface {
    Execute(ctx context.Context) (interface{}, error)
}

// QueryHandler 查询处理器接口（示例）
type QueryHandler[T Query, R any] interface {
    Handle(ctx context.Context, query T) (R, error)
}
```

### Event 模式

框架提供事件模式的抽象，用于处理领域事件。

```go
// Event 事件接口（示例）
type Event interface {
    Type() string
    Data() interface{}
    Timestamp() time.Time
}

// EventHandler 事件处理器接口（示例）
type EventHandler[T Event] interface {
    Handle(ctx context.Context, event T) error
}
```

### DTO 基类

框架提供 DTO 的基类和工具函数。

```go
// DTO 数据传输对象基类（示例）
type DTO struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ToDTO 转换函数接口（示例）
type ToDTO[T any] interface {
    ToDTO() T
}
```

## 用户如何使用

### 1. 定义自己的应用服务

用户在自己的项目中定义应用服务：

```go
// 用户项目中的应用服务
package application

type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error)
    GetUser(ctx context.Context, id string) (*UserDTO, error)
    // ...
}

type userService struct {
    userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}
```

### 2. 使用命令/查询模式

```go
// 用户项目中使用命令模式
type CreateUserCommand struct {
    Email string
    Name  string
}

type CreateUserCommandHandler struct {
    userRepo domain.UserRepository
}

func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
    user := domain.NewUser(cmd.Email, cmd.Name)
    return h.userRepo.Create(ctx, user)
}
```

### 3. 使用工作流模式

```go
// 用户项目中使用 Temporal 工作流
func UserWorkflow(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
    // 工作流逻辑
}
```

## 设计原则

1. **模式抽象**: 提供通用的应用模式抽象，不包含具体业务
2. **用例编排**: 应用层负责协调领域对象完成业务用例
3. **DTO 转换**: 应用层负责领域对象和 DTO 之间的转换
4. **事务边界**: 应用服务方法通常是一个事务边界

## 相关资源

- [应用服务示例](../../examples/application-service/) - 应用服务实现示例
- [命令/查询模式示例](../../examples/cqrs/) - CQRS 模式示例
- [工作流示例](../../examples/workflow/) - Temporal 工作流示例
