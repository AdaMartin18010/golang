package patterns

import "context"

// Command 命令接口（框架抽象）
//
// 设计原理：
// 1. 命令模式（Command Pattern）是 CQRS（Command Query Responsibility Segregation）的核心
// 2. 命令用于表示写操作（Create、Update、Delete），不返回数据
// 3. 命令是不可变的，表示一个业务操作的意图
// 4. 命令应该包含执行操作所需的所有信息
//
// 架构位置：
// - 定义：Application Layer (internal/application/patterns/)
// - 使用：Application Layer 和 Interfaces Layer
//
// CQRS 原理：
// - Command（命令）：写操作，改变系统状态
// - Query（查询）：读操作，不改变系统状态
// - 分离读写操作，提高系统的可扩展性和性能
//
// 使用场景：
// 1. 创建、更新、删除操作
// 2. 需要事务保证的操作
// 3. 需要记录操作日志的操作
// 4. 需要触发领域事件的操作
//
// 示例：
//   // 定义命令
//   type CreateUserCommand struct {
//       Email string
//       Name  string
//   }
//
//   func (c CreateUserCommand) Execute(ctx context.Context) error {
//       // 命令本身不包含执行逻辑，只是数据载体
//       return nil
//   }
//
//   // 定义命令处理器
//   type CreateUserCommandHandler struct {
//       userRepo domain.UserRepository
//       eventBus eventbus.EventBus
//   }
//
//   func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
//       // 1. 验证命令
//       if err := validateCreateUserCommand(cmd); err != nil {
//           return err
//       }
//
//       // 2. 创建领域实体
//       user := domain.NewUser(cmd.Email, cmd.Name)
//
//       // 3. 保存实体
//       if err := h.userRepo.Create(ctx, user); err != nil {
//           return err
//       }
//
//       // 4. 发布领域事件
//       event := eventbus.NewEvent("user.created", user)
//       return h.eventBus.Publish(ctx, event)
//   }
//
// 注意事项：
// - 命令应该是不可变的（immutable）
// - 命令应该包含验证逻辑（可选）
// - 命令处理器应该处理事务
// - 命令处理器应该处理错误和异常
//
// 与 Query 的区别：
// - Command：写操作，改变状态，不返回数据
// - Query：读操作，不改变状态，返回数据
//
// 用户需要根据业务需求定义具体的命令，例如：
//   type UpdateUserCommand struct {
//       ID   string
//       Name string
//   }
//
//   func (c UpdateUserCommand) Execute(ctx context.Context) error {
//       return nil
//   }
type Command interface {
	// Execute 执行命令
	//
	// 注意：这个方法通常不包含实际的执行逻辑，只是满足接口要求
	// 实际的执行逻辑应该在 CommandHandler 中实现
	//
	// 参数：
	//   - ctx: 上下文，用于传递请求信息、超时控制等
	//
	// 返回：
	//   - error: 执行失败时返回错误
	//
	// 实现建议：
	//   - 可以在这里添加命令级别的验证
	//   - 可以返回 nil，实际逻辑在 Handler 中
	Execute(ctx context.Context) error
}

// CommandHandler 命令处理器接口（框架抽象）
//
// 设计原理：
// 1. 命令处理器负责执行命令的实际逻辑
// 2. 使用泛型 T 确保类型安全
// 3. 命令处理器应该处理事务、错误处理和事件发布
//
// 职责：
// 1. 验证命令参数
// 2. 协调领域对象完成业务操作
// 3. 管理事务
// 4. 发布领域事件
// 5. 处理错误和异常
//
// 示例：
//   type CreateUserCommandHandler struct {
//       userRepo domain.UserRepository
//       eventBus eventbus.EventBus
//   }
//
//   func NewCreateUserCommandHandler(
//       userRepo domain.UserRepository,
//       eventBus eventbus.EventBus,
//   ) *CreateUserCommandHandler {
//       return &CreateUserCommandHandler{
//           userRepo: userRepo,
//           eventBus: eventBus,
//       }
//   }
//
//   func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
//       // 实现命令处理逻辑
//   }
//
// 注意事项：
// - 命令处理器应该是无状态的
// - 应该通过依赖注入获取依赖
// - 应该处理事务边界
// - 应该处理并发冲突
type CommandHandler[T Command] interface {
	// Handle 处理命令
	//
	// 参数：
	//   - ctx: 上下文
	//   - cmd: 要处理的命令，类型为 T
	//
	// 返回：
	//   - error: 处理失败时返回错误
	//
	// 处理流程：
	// 1. 验证命令参数
	// 2. 执行业务逻辑（协调领域对象）
	// 3. 保存变更（通过仓储）
	// 4. 发布领域事件（可选）
	// 5. 返回结果
	//
	// 错误处理：
	// - 应该返回明确的错误信息
	// - 应该区分业务错误和系统错误
	// - 应该记录错误日志
	Handle(ctx context.Context, cmd T) error
}

// CommandResult 命令执行结果
//
// 设计原理：
// 1. 命令通常不返回数据（只返回错误）
// 2. 但在某些场景下，可能需要返回执行结果
// 3. CommandResult 提供了标准化的结果格式
//
// 使用场景：
// 1. 需要返回创建实体的 ID
// 2. 需要返回操作状态信息
// 3. 需要返回操作影响的数据
//
// 示例：
//   func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*CommandResult, error) {
//       user := domain.NewUser(cmd.Email, cmd.Name)
//       if err := h.userRepo.Create(ctx, user); err != nil {
//           return nil, err
//       }
//
//       return &CommandResult{
//           Success: true,
//           Message: "User created successfully",
//           Data:    user.ID,
//       }, nil
//   }
//
// 注意事项：
// - 大多数情况下，命令不需要返回结果
// - 如果需要返回结果，应该使用 CommandResult
// - 结果应该只包含必要的信息
type CommandResult struct {
	// Success 操作是否成功
	Success bool

	// Message 操作结果消息
	Message string

	// Data 操作返回的数据（可选）
	// 例如：创建的实体 ID、更新的记录数等
	Data interface{}
}
