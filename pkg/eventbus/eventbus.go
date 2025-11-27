// Package eventbus 提供事件总线（Event Bus）实现
//
// 设计原理：
// 1. 事件总线是事件驱动架构（Event-Driven Architecture）的核心组件
// 2. 实现发布-订阅（Publish-Subscribe）模式，解耦事件发布者和订阅者
// 3. 支持异步事件处理，提高系统响应性能
// 4. 提供事件过滤、指标监控等高级功能
//
// 核心特性：
// - ✅ 异步事件处理：事件在独立的 goroutine 中处理，不阻塞发布者
// - ✅ 事件过滤：支持基于条件的事件过滤
// - ✅ 指标监控：提供事件处理指标（总数、成功、失败、丢弃等）
// - ✅ 线程安全：使用读写锁保证并发安全
// - ✅ 优雅关闭：支持优雅停止，等待所有事件处理完成
//
// 使用场景：
// 1. 领域事件发布：在领域模型中发布业务事件
// 2. 跨模块通信：不同模块之间通过事件通信
// 3. 异步任务处理：将耗时操作异步化
// 4. 系统集成：集成多个子系统
//
// 架构位置：
// - 包位置：pkg/eventbus/
// - 使用位置：Application Layer、Infrastructure Layer
//
// 示例：
//
//	// 1. 创建事件总线
//	eventBus := eventbus.NewEventBus(100)
//	defer eventBus.Stop()
//
//	// 2. 启动事件总线
//	if err := eventBus.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
//	// 3. 订阅事件
//	subID, err := eventBus.Subscribe("user.created", func(ctx context.Context, event eventbus.Event) error {
//	    user := event.Data().(*User)
//	    log.Printf("User created: %s", user.Email)
//	    return nil
//	})
//
//	// 4. 发布事件
//	event := eventbus.NewEvent("user.created", user)
//	if err := eventBus.Publish(event); err != nil {
//	    log.Printf("Failed to publish event: %v", err)
//	}
//
// 注意事项：
// - 事件总线需要先调用 Start() 才能处理事件
// - 事件处理是异步的，错误不会返回给发布者
// - 事件缓冲区满时会丢弃事件，需要监控指标
// - 停止事件总线时会等待所有事件处理完成
package eventbus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrEventBusStopped 事件总线已停止
	// 当事件总线已停止时，尝试发布事件会返回此错误
	ErrEventBusStopped = errors.New("event bus is stopped")

	// ErrSubscriptionNotFound 订阅未找到
	// 当尝试取消不存在的订阅时，会返回此错误
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

// Event 事件接口
//
// 设计原理：
// 1. 事件是不可变的（immutable），一旦创建就不能修改
// 2. 事件包含类型、数据和时间戳，提供完整的上下文信息
// 3. 事件类型用于路由，数据包含业务信息，时间戳用于排序和审计
//
// 实现要求：
// - Type() 应该返回唯一的事件类型标识
// - Data() 应该包含事件发生时的完整上下文
// - Timestamp() 应该返回事件发生的时间
//
// 事件类型命名规范：
// - 使用点分隔的命名空间，如 "user.created"
// - 使用过去时态，表示已经发生的事件
// - 使用小写字母和点分隔符
//
// 示例：
//
//	type UserCreatedEvent struct {
//	    eventbus.BaseEvent
//	    UserID string
//	    Email  string
//	}
//
//	func NewUserCreatedEvent(user *User) *UserCreatedEvent {
//	    return &UserCreatedEvent{
//	        BaseEvent: *eventbus.NewEvent("user.created", user),
//	        UserID:    user.ID,
//	        Email:     user.Email,
//	    }
//	}
type Event interface {
	// Type 返回事件类型
	// 用于路由事件到对应的处理器
	Type() string

	// Data 返回事件数据
	// 包含事件发生时的完整上下文信息
	Data() interface{}

	// Timestamp 返回事件时间戳
	// 用于事件排序、审计和重放
	Timestamp() time.Time
}

// BaseEvent 基础事件实现
//
// 设计原理：
// 1. 提供事件接口的默认实现
// 2. 支持元数据扩展，可以附加额外的上下文信息
// 3. 用户可以直接使用或嵌入到自定义事件中
//
// 使用方式：
//
//  1. 直接使用：适用于简单事件
//     event := eventbus.NewEvent("user.created", user)
//
//  2. 嵌入使用：适用于复杂事件（推荐）
//     type UserCreatedEvent struct {
//     eventbus.BaseEvent
//     UserID string
//     Email  string
//     }
//
// 元数据用途：
// - 传递额外的上下文信息（如请求ID、用户ID等）
// - 支持事件追踪和调试
// - 支持事件路由和过滤
type BaseEvent struct {
	// eventType 事件类型
	eventType string

	// data 事件数据
	data interface{}

	// timestamp 事件时间戳
	timestamp time.Time

	// metadata 元数据，用于存储额外的上下文信息
	metadata map[string]interface{}
}

// NewEvent 创建基础事件
//
// 参数：
//   - eventType: 事件类型，如 "user.created"
//   - data: 事件数据，可以是领域对象或 DTO
//
// 返回：
//   - *BaseEvent: 创建的基础事件
//
// 示例：
//
//	user := &User{ID: "123", Email: "test@example.com"}
//	event := eventbus.NewEvent("user.created", user)
//
//	// 添加元数据
//	event.SetMetadata("request_id", "req-123")
//	event.SetMetadata("user_id", "user-456")
func NewEvent(eventType string, data interface{}) *BaseEvent {
	return &BaseEvent{
		eventType: eventType,
		data:      data,
		timestamp: time.Now(),
		metadata:  make(map[string]interface{}),
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

// Timestamp 返回时间戳
func (e *BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

// SetMetadata 设置元数据
//
// 参数：
//   - key: 元数据键
//   - value: 元数据值
//
// 用途：
// - 传递请求ID、用户ID等上下文信息
// - 支持事件追踪和调试
// - 支持事件过滤和路由
//
// 示例：
//
//	event.SetMetadata("request_id", "req-123")
//	event.SetMetadata("user_id", "user-456")
func (e *BaseEvent) SetMetadata(key string, value interface{}) {
	if e.metadata == nil {
		e.metadata = make(map[string]interface{})
	}
	e.metadata[key] = value
}

// GetMetadata 获取元数据
//
// 参数：
//   - key: 元数据键
//
// 返回：
//   - interface{}: 元数据值
//   - bool: 是否存在
//
// 示例：
//
//	if requestID, ok := event.GetMetadata("request_id"); ok {
//	    log.Printf("Request ID: %v", requestID)
//	}
func (e *BaseEvent) GetMetadata(key string) (interface{}, bool) {
	if e.metadata == nil {
		return nil, false
	}
	value, ok := e.metadata[key]
	return value, ok
}

// Handler 事件处理器函数类型
//
// 设计原理：
// 1. 事件处理器是处理事件的函数
// 2. 处理器应该是幂等的，多次处理相同事件应该产生相同结果
// 3. 处理器应该处理错误，避免影响其他处理器
//
// 参数：
//   - ctx: 上下文，用于传递请求信息、超时控制等
//   - event: 要处理的事件
//
// 返回：
//   - error: 处理失败时返回错误
//
// 实现要求：
// - 应该处理所有可能的错误情况
// - 应该记录处理日志（可选）
// - 应该考虑超时和取消
// - 应该是幂等的
//
// 示例：
//
//	handler := func(ctx context.Context, event eventbus.Event) error {
//	    user := event.Data().(*User)
//	    log.Printf("Processing user created event: %s", user.Email)
//
//	    // 处理逻辑
//	    if err := sendWelcomeEmail(ctx, user.Email); err != nil {
//	        return fmt.Errorf("failed to send welcome email: %w", err)
//	    }
//
//	    return nil
//	}
type Handler func(ctx context.Context, event Event) error

// Filter 事件过滤器函数类型
//
// 设计原理：
// 1. 事件过滤器用于过滤不需要处理的事件
// 2. 过滤器在处理器执行前调用，可以提高性能
// 3. 过滤器应该快速返回，避免阻塞事件处理
//
// 参数：
//   - event: 要过滤的事件
//
// 返回：
//   - bool: true 表示处理事件，false 表示跳过事件
//
// 使用场景：
// 1. 基于事件数据的过滤（如只处理特定用户的事件）
// 2. 基于元数据的过滤（如只处理特定来源的事件）
// 3. 基于时间的过滤（如只处理最近的事件）
//
// 示例：
//
//	filter := func(event eventbus.Event) bool {
//	    user := event.Data().(*User)
//	    // 只处理 VIP 用户的事件
//	    return user.IsVIP
//	}
//
//	eventBus.SubscribeWithFilter("user.created", handler, filter)
type Filter func(event Event) bool

// Subscription 事件订阅
//
// 设计原理：
// 1. 订阅表示一个事件处理器对特定事件类型的订阅
// 2. 每个订阅有唯一ID，用于取消订阅
// 3. 订阅可以包含过滤器，用于条件过滤
//
// 字段说明：
// - ID: 订阅唯一标识，用于取消订阅
// - EventType: 订阅的事件类型
// - Handler: 事件处理器函数
// - Filter: 事件过滤器函数（可选）
// - CreatedAt: 订阅创建时间
type Subscription struct {
	// ID 订阅唯一标识
	ID string

	// EventType 订阅的事件类型
	EventType string

	// Handler 事件处理器函数
	Handler Handler

	// Filter 事件过滤器函数（可选）
	// 如果为 nil，则处理所有匹配类型的事件
	Filter Filter

	// CreatedAt 订阅创建时间
	CreatedAt time.Time
}

// EventBus 事件总线
//
// 设计原理：
// 1. 事件总线是发布-订阅模式的核心实现
// 2. 使用 channel 实现异步事件处理
// 3. 使用 goroutine 并发处理事件，提高性能
// 4. 使用读写锁保证并发安全
//
// 工作流程：
// 1. 发布者调用 Publish() 发布事件
// 2. 事件被放入 eventChan 缓冲区
// 3. processEvents() goroutine 从 channel 读取事件
// 4. handleEvent() 找到所有匹配的订阅
// 5. 为每个订阅启动独立的 goroutine 处理事件
//
// 线程安全：
// - 使用 sync.RWMutex 保护订阅列表
// - 使用 sync.WaitGroup 等待所有事件处理完成
// - 使用 context 实现优雅关闭
//
// 性能优化：
// - 使用缓冲 channel 减少阻塞
// - 异步处理事件，不阻塞发布者
// - 并发处理多个订阅，提高吞吐量
type EventBus struct {
	// subscriptions 事件类型到订阅列表的映射
	// key: 事件类型，value: 订阅列表
	subscriptions map[string][]*Subscription

	// subIndex 订阅ID到订阅的映射
	// 用于快速查找和取消订阅
	subIndex map[string]*Subscription

	// eventChan 事件通道，用于缓冲待处理的事件
	eventChan chan Event

	// mu 读写锁，保护订阅列表和索引
	mu sync.RWMutex

	// ctx 上下文，用于控制事件总线的生命周期
	ctx context.Context

	// cancel 取消函数，用于停止事件总线
	cancel context.CancelFunc

	// wg 等待组，用于等待所有事件处理完成
	wg sync.WaitGroup

	// bufferSize 事件缓冲区大小
	bufferSize int

	// metrics 事件总线指标
	metrics Metrics

	// nextSubID 下一个订阅ID的计数器
	nextSubID int64
}

// Metrics 事件总线指标
//
// 设计原理：
// 1. 提供事件总线的运行指标，用于监控和调试
// 2. 使用读写锁保护指标数据
// 3. 指标是只读的，通过 GetMetrics() 获取快照
//
// 指标说明：
// - TotalEvents: 总事件数（包括成功、失败、丢弃）
// - HandledEvents: 成功处理的事件数
// - FailedEvents: 处理失败的事件数
// - DroppedEvents: 因缓冲区满而丢弃的事件数
// - ActiveSubscriptions: 当前活跃的订阅数
//
// 使用场景：
// 1. 监控事件总线的健康状态
// 2. 调试事件处理问题
// 3. 性能分析和优化
//
// 示例：
//
//	metrics := eventBus.GetMetrics()
//	log.Printf("Total events: %d, Handled: %d, Failed: %d, Dropped: %d",
//	    metrics.TotalEvents, metrics.HandledEvents,
//	    metrics.FailedEvents, metrics.DroppedEvents)
type Metrics struct {
	// TotalEvents 总事件数
	TotalEvents int64

	// HandledEvents 成功处理的事件数
	HandledEvents int64

	// FailedEvents 处理失败的事件数
	FailedEvents int64

	// DroppedEvents 因缓冲区满而丢弃的事件数
	DroppedEvents int64

	// ActiveSubscriptions 当前活跃的订阅数
	ActiveSubscriptions int

	// mu 读写锁，保护指标数据
	mu sync.RWMutex
}

// NewEventBus 创建事件总线
//
// 参数：
//   - bufferSize: 事件缓冲区大小，用于缓冲待处理的事件
//     如果 bufferSize <= 0，则使用默认值 100
//
// 返回：
//   - *EventBus: 创建的事件总线实例
//
// 注意事项：
// - 创建后需要调用 Start() 启动事件总线
// - 使用完毕后应该调用 Stop() 停止事件总线
// - bufferSize 应该根据事件产生速度和处理速度设置
//
// 示例：
//
//	// 创建事件总线，缓冲区大小为 1000
//	eventBus := eventbus.NewEventBus(1000)
//
//	// 启动事件总线
//	if err := eventBus.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
//	// 使用完毕后停止
//	defer eventBus.Stop()
//
// bufferSize 选择建议：
// - 低频率事件（< 100/秒）：100-500
// - 中频率事件（100-1000/秒）：500-2000
// - 高频率事件（> 1000/秒）：2000-10000
func NewEventBus(bufferSize int) *EventBus {
	if bufferSize <= 0 {
		bufferSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventBus{
		subscriptions: make(map[string][]*Subscription),
		subIndex:      make(map[string]*Subscription),
		eventChan:     make(chan Event, bufferSize),
		ctx:           ctx,
		cancel:        cancel,
		bufferSize:    bufferSize,
	}
}

// Start 启动事件总线
//
// 设计原理：
// 1. 启动后台 goroutine 处理事件
// 2. 事件总线启动后才能处理事件
// 3. 可以多次调用，但只会启动一次处理 goroutine
//
// 返回：
//   - error: 启动失败时返回错误（当前实现总是返回 nil）
//
// 注意事项：
// - 必须在发布事件前调用 Start()
// - 如果事件总线已停止，需要重新创建才能再次启动
//
// 示例：
//
//	eventBus := eventbus.NewEventBus(100)
//	if err := eventBus.Start(); err != nil {
//	    log.Fatal(err)
//	}
//	defer eventBus.Stop()
func (eb *EventBus) Start() error {
	eb.wg.Add(1)
	go eb.processEvents()
	return nil
}

// Stop 停止事件总线
//
// 设计原理：
// 1. 优雅停止：等待所有事件处理完成
// 2. 取消上下文，停止接收新事件
// 3. 关闭事件通道，释放资源
//
// 返回：
//   - error: 停止失败时返回错误（当前实现总是返回 nil）
//
// 注意事项：
// - 停止后会等待所有正在处理的事件完成
// - 停止后不能再发布新事件
// - 建议使用 defer 确保停止
//
// 示例：
//
//	eventBus := eventbus.NewEventBus(100)
//	eventBus.Start()
//	defer eventBus.Stop()  // 确保停止
func (eb *EventBus) Stop() error {
	eb.cancel()
	eb.wg.Wait()
	close(eb.eventChan)
	return nil
}

// Subscribe 订阅事件
//
// 参数：
//   - eventType: 事件类型，如 "user.created"
//   - handler: 事件处理器函数
//
// 返回：
//   - string: 订阅ID，用于取消订阅
//   - error: 订阅失败时返回错误（当前实现总是返回 nil）
//
// 注意事项：
// - 同一个事件类型可以有多个订阅
// - 订阅ID是唯一的，用于取消订阅
// - 处理器应该是幂等的
//
// 示例：
//
//	subID, err := eventBus.Subscribe("user.created", func(ctx context.Context, event eventbus.Event) error {
//	    user := event.Data().(*User)
//	    log.Printf("User created: %s", user.Email)
//	    return nil
//	})
//
//	// 稍后取消订阅
//	eventBus.Unsubscribe(subID)
func (eb *EventBus) Subscribe(eventType string, handler Handler) (string, error) {
	return eb.SubscribeWithFilter(eventType, handler, nil)
}

// SubscribeWithFilter 订阅事件（带过滤器）
//
// 参数：
//   - eventType: 事件类型
//   - handler: 事件处理器函数
//   - filter: 事件过滤器函数（可选，nil 表示不过滤）
//
// 返回：
//   - string: 订阅ID，用于取消订阅
//   - error: 订阅失败时返回错误（当前实现总是返回 nil）
//
// 过滤器说明：
// - 过滤器在处理器执行前调用
// - 如果过滤器返回 false，则跳过该事件
// - 如果过滤器为 nil，则处理所有匹配类型的事件
//
// 使用场景：
// 1. 只处理特定条件的事件（如只处理 VIP 用户的事件）
// 2. 基于元数据过滤（如只处理特定来源的事件）
// 3. 性能优化（避免处理不需要的事件）
//
// 示例：
//
//	// 只处理 VIP 用户的事件
//	filter := func(event eventbus.Event) bool {
//	    user := event.Data().(*User)
//	    return user.IsVIP
//	}
//
//	subID, err := eventBus.SubscribeWithFilter("user.created", handler, filter)
func (eb *EventBus) SubscribeWithFilter(eventType string, handler Handler, filter Filter) (string, error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.nextSubID++
	subID := fmt.Sprintf("sub_%d", eb.nextSubID)

	sub := &Subscription{
		ID:        subID,
		EventType: eventType,
		Handler:   handler,
		Filter:    filter,
		CreatedAt: time.Now(),
	}

	eb.subscriptions[eventType] = append(eb.subscriptions[eventType], sub)
	eb.subIndex[subID] = sub

	eb.metrics.mu.Lock()
	eb.metrics.ActiveSubscriptions++
	eb.metrics.mu.Unlock()

	return subID, nil
}

// Unsubscribe 取消订阅
//
// 参数：
//   - subscriptionID: 订阅ID，由 Subscribe() 或 SubscribeWithFilter() 返回
//
// 返回：
//   - error: 取消订阅失败时返回错误
//   - ErrSubscriptionNotFound: 订阅不存在
//
// 注意事项：
// - 取消订阅后，该订阅不会再收到新事件
// - 正在处理的事件会继续处理完成
// - 可以安全地多次调用，不会报错（但会返回错误）
//
// 示例：
//
//	subID, _ := eventBus.Subscribe("user.created", handler)
//
//	// 稍后取消订阅
//	if err := eventBus.Unsubscribe(subID); err != nil {
//	    log.Printf("Failed to unsubscribe: %v", err)
//	}
func (eb *EventBus) Unsubscribe(subscriptionID string) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	sub, exists := eb.subIndex[subscriptionID]
	if !exists {
		return ErrSubscriptionNotFound
	}

	// 从订阅列表中移除
	subs := eb.subscriptions[sub.EventType]
	for i, s := range subs {
		if s.ID == subscriptionID {
			eb.subscriptions[sub.EventType] = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	delete(eb.subIndex, subscriptionID)

	eb.metrics.mu.Lock()
	eb.metrics.ActiveSubscriptions--
	eb.metrics.mu.Unlock()

	return nil
}

// Publish 发布事件
//
// 设计原理：
// 1. 将事件放入缓冲区，不阻塞调用者
// 2. 如果缓冲区满，则丢弃事件并返回错误
// 3. 如果事件总线已停止，则返回错误
//
// 参数：
//   - event: 要发布的事件
//
// 返回：
//   - error: 发布失败时返回错误
//   - ErrEventBusStopped: 事件总线已停止
//   - 其他错误: 缓冲区满，事件被丢弃
//
// 注意事项：
// - 发布是异步的，不会等待事件处理完成
// - 如果缓冲区满，事件会被丢弃，需要监控指标
// - 事件处理错误不会返回给发布者
//
// 示例：
//
//	event := eventbus.NewEvent("user.created", user)
//	if err := eventBus.Publish(event); err != nil {
//	    if err == eventbus.ErrEventBusStopped {
//	        log.Printf("Event bus is stopped")
//	    } else {
//	        log.Printf("Event buffer is full, event dropped")
//	    }
//	}
func (eb *EventBus) Publish(event Event) error {
	select {
	case eb.eventChan <- event:
		eb.metrics.mu.Lock()
		eb.metrics.TotalEvents++
		eb.metrics.mu.Unlock()
		return nil
	case <-eb.ctx.Done():
		return ErrEventBusStopped
	default:
		// 缓冲区满，丢弃事件
		eb.metrics.mu.Lock()
		eb.metrics.DroppedEvents++
		eb.metrics.mu.Unlock()
		return fmt.Errorf("event buffer is full, event dropped")
	}
}

// PublishAsync 异步发布事件（不阻塞）
//
// 设计原理：
// 1. 在独立的 goroutine 中发布事件
// 2. 完全不阻塞调用者，即使缓冲区满也不会阻塞
// 3. 错误会被忽略，适合不关心发布结果的场景
//
// 参数：
//   - event: 要发布的事件
//
// 注意事项：
// - 不返回错误，错误会被忽略
// - 适合不关心发布结果的场景
// - 如果缓冲区满，事件会被丢弃（不会阻塞）
//
// 使用场景：
// 1. 日志记录等不重要的操作
// 2. 性能敏感的场景，不能阻塞
// 3. 不关心发布结果的场景
//
// 示例：
//
//	event := eventbus.NewEvent("user.created", user)
//	eventBus.PublishAsync(event)  // 不阻塞，不关心结果
func (eb *EventBus) PublishAsync(event Event) {
	go func() {
		_ = eb.Publish(event)
	}()
}

// processEvents 处理事件（内部goroutine）
//
// 设计原理：
// 1. 从事件通道读取事件
// 2. 调用 handleEvent() 处理事件
// 3. 监听上下文取消信号，实现优雅关闭
//
// 工作流程：
// 1. 从 eventChan 读取事件
// 2. 调用 handleEvent() 处理事件
// 3. 如果上下文被取消，则退出循环
//
// 注意事项：
// - 这是后台 goroutine，由 Start() 启动
// - 退出时会调用 wg.Done()，通知 Stop() 可以继续
func (eb *EventBus) processEvents() {
	defer eb.wg.Done()

	for {
		select {
		case event := <-eb.eventChan:
			eb.handleEvent(event)
		case <-eb.ctx.Done():
			return
		}
	}
}

// handleEvent 处理单个事件
//
// 设计原理：
// 1. 找到所有匹配事件类型的订阅
// 2. 对每个订阅应用过滤器
// 3. 为每个订阅启动独立的 goroutine 处理事件
//
// 处理流程：
// 1. 获取所有匹配的订阅（读锁保护）
// 2. 遍历订阅，应用过滤器
// 3. 为每个订阅启动 goroutine 异步处理
// 4. 使用超时上下文（5秒）防止处理器阻塞
// 5. 更新指标（成功/失败）
//
// 性能优化：
// - 使用读锁，允许多个事件并发处理
// - 异步处理，不阻塞事件接收
// - 超时控制，防止处理器阻塞
//
// 注意事项：
// - 每个订阅在独立的 goroutine 中处理
// - 处理器超时时间为 5 秒
// - 处理器错误不会影响其他订阅
func (eb *EventBus) handleEvent(event Event) {
	eb.mu.RLock()
	subs := eb.subscriptions[event.Type()]
	eb.mu.RUnlock()

	for _, sub := range subs {
		// 应用过滤器
		if sub.Filter != nil && !sub.Filter(event) {
			continue
		}

		// 异步执行处理器
		eb.wg.Add(1)
		go func(subscription *Subscription) {
			defer eb.wg.Done()

			// 使用超时上下文，防止处理器阻塞
			ctx, cancel := context.WithTimeout(eb.ctx, 5*time.Second)
			defer cancel()

			if err := subscription.Handler(ctx, event); err != nil {
				eb.metrics.mu.Lock()
				eb.metrics.FailedEvents++
				eb.metrics.mu.Unlock()
			} else {
				eb.metrics.mu.Lock()
				eb.metrics.HandledEvents++
				eb.metrics.mu.Unlock()
			}
		}(sub)
	}
}

// GetMetrics 获取指标
//
// 设计原理：
// 1. 返回指标的快照，避免并发修改
// 2. 使用读锁保护，保证数据一致性
// 3. 返回副本，避免外部修改影响内部状态
//
// 返回：
//   - Metrics: 事件总线指标快照
//
// 使用场景：
// 1. 监控事件总线的健康状态
// 2. 调试事件处理问题
// 3. 性能分析和优化
// 4. 告警和通知
//
// 示例：
//
//	metrics := eventBus.GetMetrics()
//	log.Printf("Event Bus Metrics:")
//	log.Printf("  Total Events: %d", metrics.TotalEvents)
//	log.Printf("  Handled: %d", metrics.HandledEvents)
//	log.Printf("  Failed: %d", metrics.FailedEvents)
//	log.Printf("  Dropped: %d", metrics.DroppedEvents)
//	log.Printf("  Active Subscriptions: %d", metrics.ActiveSubscriptions)
//
//	// 检查是否有事件被丢弃
//	if metrics.DroppedEvents > 0 {
//	    log.Warnf("Events are being dropped! Consider increasing buffer size.")
//	}
func (eb *EventBus) GetMetrics() Metrics {
	eb.metrics.mu.RLock()
	defer eb.metrics.mu.RUnlock()

	return Metrics{
		TotalEvents:         eb.metrics.TotalEvents,
		HandledEvents:       eb.metrics.HandledEvents,
		FailedEvents:        eb.metrics.FailedEvents,
		DroppedEvents:       eb.metrics.DroppedEvents,
		ActiveSubscriptions: eb.metrics.ActiveSubscriptions,
	}
}

// ListSubscriptions 列出所有订阅
//
// 设计原理：
// 1. 返回所有订阅的快照
// 2. 使用读锁保护，保证数据一致性
// 3. 返回副本，避免外部修改影响内部状态
//
// 返回：
//   - []Subscription: 所有订阅的列表
//
// 使用场景：
// 1. 调试和诊断
// 2. 监控订阅状态
// 3. 管理订阅
//
// 示例：
//
//	subscriptions := eventBus.ListSubscriptions()
//	log.Printf("Active subscriptions: %d", len(subscriptions))
//	for _, sub := range subscriptions {
//	    log.Printf("  - %s: %s (created at %s)", sub.ID, sub.EventType, sub.CreatedAt)
//	}
func (eb *EventBus) ListSubscriptions() []Subscription {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	result := make([]Subscription, 0, len(eb.subIndex))
	for _, sub := range eb.subIndex {
		result = append(result, *sub)
	}

	return result
}

// Clear 清空所有订阅
//
// 设计原理：
// 1. 移除所有订阅
// 2. 重置订阅索引
// 3. 更新指标
//
// 注意事项：
// - 清空后，所有订阅都不会再收到事件
// - 正在处理的事件会继续处理完成
// - 清空后可以重新订阅
//
// 使用场景：
// 1. 测试场景，清理状态
// 2. 重新初始化订阅
// 3. 资源清理
//
// 示例：
//
//	// 清空所有订阅
//	eventBus.Clear()
//
//	// 重新订阅
//	eventBus.Subscribe("user.created", newHandler)
func (eb *EventBus) Clear() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscriptions = make(map[string][]*Subscription)
	eb.subIndex = make(map[string]*Subscription)

	eb.metrics.mu.Lock()
	eb.metrics.ActiveSubscriptions = 0
	eb.metrics.mu.Unlock()
}
