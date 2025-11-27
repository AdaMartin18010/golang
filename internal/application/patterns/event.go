package patterns

import (
	"context"
	"time"
)

// Event 事件接口（框架抽象）
// 用于处理领域事件
type Event interface {
	// Type 返回事件类型
	Type() string

	// Data 返回事件数据
	Data() interface{}

	// Timestamp 返回事件时间戳
	Timestamp() time.Time
}

// EventHandler 事件处理器接口（框架抽象）
type EventHandler[T Event] interface {
	// Handle 处理事件
	Handle(ctx context.Context, event T) error
}

// BaseEvent 基础事件实现
type BaseEvent struct {
	eventType string
	data      interface{}
	timestamp time.Time
}

// NewBaseEvent 创建基础事件
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
