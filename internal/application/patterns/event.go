package patterns

import (
	"context"
	"time"
)

// Event 事件接口（框架抽象）
//
// 设计原理：
// 1. 领域事件（Domain Event）表示领域中的重要业务事件
// 2. 事件是不可变的（immutable），一旦创建就不能修改
// 3. 事件用于解耦业务逻辑，支持异步处理和事件溯源
// 4. 事件应该包含事件发生时的完整上下文信息
//
// 架构位置：
// - 定义：Application Layer (internal/application/patterns/)
// - 使用：Application Layer 和 Infrastructure Layer
//
// 领域事件原理：
// 1. 事件驱动架构（Event-Driven Architecture）
//    - 业务操作产生事件
//    - 事件被发布到事件总线
//    - 事件处理器异步处理事件
// 2. 解耦和扩展性
//    - 业务逻辑与事件处理解耦
//    - 可以轻松添加新的事件处理器
//    - 支持多个处理器处理同一事件
// 3. 事件溯源（Event Sourcing）
//    - 可以记录所有业务事件
//    - 可以通过重放事件重建状态
//
// 使用场景：
// 1. 业务操作完成后需要通知其他模块
// 2. 需要异步处理的业务逻辑
// 3. 需要记录业务操作历史
// 4. 需要实现最终一致性
//
// 示例：
//   // 定义领域事件
//   type UserCreatedEvent struct {
//       patterns.BaseEvent
//       UserID  string
//       Email   string
//       CreatedAt time.Time
//   }
//
//   func NewUserCreatedEvent(user *domain.User) *UserCreatedEvent {
//       return &UserCreatedEvent{
//           BaseEvent: *patterns.NewBaseEvent("user.created", user),
//           UserID:    user.ID,
//           Email:     user.Email,
//           CreatedAt: user.CreatedAt,
//       }
//   }
//
//   // 在命令处理器中发布事件
//   func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
//       user := domain.NewUser(cmd.Email, cmd.Name)
//       if err := h.userRepo.Create(ctx, user); err != nil {
//           return err
//       }
//
//       // 发布领域事件
//       event := NewUserCreatedEvent(user)
//       return h.eventBus.Publish(ctx, event)
//   }
//
//   // 定义事件处理器
//   type SendWelcomeEmailHandler struct {
//       emailService EmailService
//   }
//
//   func (h *SendWelcomeEmailHandler) Handle(ctx context.Context, event *UserCreatedEvent) error {
//       return h.emailService.SendWelcomeEmail(ctx, event.Email)
//   }
//
// 注意事项：
// - 事件应该是不可变的
// - 事件应该包含完整的上下文信息
// - 事件处理应该是幂等的
// - 事件处理失败应该重试或记录
//
// 用户需要根据业务需求定义具体的事件，例如：
//   type OrderPlacedEvent struct {
//       patterns.BaseEvent
//       OrderID string
//       UserID  string
//       TotalAmount float64
//   }
type Event interface {
	// Type 返回事件类型
	//
	// 返回：
	//   - string: 事件类型，用于路由和识别事件
	//
	// 事件类型命名规范：
	//   - 使用点分隔的命名空间，如 "user.created"
	//   - 使用过去时态，表示已经发生的事件
	//   - 使用小写字母和点分隔符
	//
	// 示例：
	//   - "user.created"
	//   - "order.placed"
	//   - "payment.completed"
	Type() string

	// Data 返回事件数据
	//
	// 返回：
	//   - interface{}: 事件数据，可以是领域对象或 DTO
	//
	// 数据要求：
	//   - 应该包含事件发生时的完整上下文
	//   - 应该包含足够的信息供事件处理器使用
	//   - 应该避免包含敏感信息（如密码）
	//
	// 示例：
	//   - 用户创建事件：包含用户信息
	//   - 订单创建事件：包含订单信息和用户信息
	Data() interface{}

	// Timestamp 返回事件时间戳
	//
	// 返回：
	//   - time.Time: 事件发生的时间
	//
	// 时间戳用途：
	//   - 记录事件发生的时间
	//   - 用于事件排序和重放
	//   - 用于审计和日志
	Timestamp() time.Time
}

// EventHandler 事件处理器接口（框架抽象）
//
// 设计原理：
// 1. 事件处理器负责处理特定类型的事件
// 2. 使用泛型 T 确保类型安全
// 3. 一个事件可以有多个处理器
// 4. 事件处理应该是异步的
//
// 职责：
// 1. 处理特定类型的事件
// 2. 执行业务逻辑（如发送邮件、更新缓存）
// 3. 处理错误和重试
//
// 示例：
//   type SendWelcomeEmailHandler struct {
//       emailService EmailService
//   }
//
//   func NewSendWelcomeEmailHandler(emailService EmailService) *SendWelcomeEmailHandler {
//       return &SendWelcomeEmailHandler{
//           emailService: emailService,
//       }
//   }
//
//   func (h *SendWelcomeEmailHandler) Handle(ctx context.Context, event *UserCreatedEvent) error {
//       return h.emailService.SendWelcomeEmail(ctx, event.Email)
//   }
//
// 注意事项：
// - 事件处理器应该是无状态的
// - 应该通过依赖注入获取依赖
// - 应该处理错误和重试
// - 应该考虑幂等性
type EventHandler[T Event] interface {
	// Handle 处理事件
	//
	// 参数：
	//   - ctx: 上下文
	//   - event: 要处理的事件，类型为 T
	//
	// 返回：
	//   - error: 处理失败时返回错误
	//
	// 处理流程：
	// 1. 验证事件数据
	// 2. 执行业务逻辑
	// 3. 处理错误和重试
	// 4. 返回结果
	//
	// 错误处理：
	// - 应该返回明确的错误信息
	// - 应该支持重试机制
	// - 应该记录错误日志
	// - 应该考虑死信队列（Dead Letter Queue）
	Handle(ctx context.Context, event T) error
}

// BaseEvent 基础事件实现
//
// 设计原理：
// 1. 提供事件接口的默认实现
// 2. 用户可以直接使用或继承
// 3. 包含事件的基本属性
//
// 使用方式：
// 1. 直接使用：适用于简单事件
// 2. 嵌入使用：适用于复杂事件（推荐）
//
// 示例：
//   // 方式1：直接使用
//   event := patterns.NewBaseEvent("user.created", user)
//
//   // 方式2：嵌入使用（推荐）
//   type UserCreatedEvent struct {
//       patterns.BaseEvent
//       UserID string
//       Email  string
//   }
//
//   func NewUserCreatedEvent(user *domain.User) *UserCreatedEvent {
//       return &UserCreatedEvent{
//           BaseEvent: *patterns.NewBaseEvent("user.created", user),
//           UserID:    user.ID,
//           Email:     user.Email,
//       }
//   }
type BaseEvent struct {
	// eventType 事件类型
	eventType string

	// data 事件数据
	data interface{}

	// timestamp 事件时间戳
	timestamp time.Time
}

// NewBaseEvent 创建基础事件
//
// 参数：
//   - eventType: 事件类型，如 "user.created"
//   - data: 事件数据，可以是领域对象或 DTO
//
// 返回：
//   - *BaseEvent: 创建的基础事件
//
// 示例：
//   event := patterns.NewBaseEvent("user.created", user)
func NewBaseEvent(eventType string, data interface{}) *BaseEvent {
	return &BaseEvent{
		eventType: eventType,
		data:      data,
		timestamp: time.Now(),
	}
}

// Type 返回事件类型
func (e *BaseEvent) Type() string {
	return e.eventType
}

// Data 返回事件数据
func (e *BaseEvent) Data() interface{} {
	return e.data
}

// Timestamp 返回事件时间戳
func (e *BaseEvent) Timestamp() time.Time {
	return e.timestamp
}
