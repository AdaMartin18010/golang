package patterns

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockEventData 测试用的事件数据
type MockEventData struct {
	ID   string
	Name string
}

// ==================== BaseEvent 测试 ====================

// TestBaseEvent 测试 BaseEvent
func TestBaseEvent(t *testing.T) {
	t.Run("创建事件", func(t *testing.T) {
		data := &MockEventData{ID: "123", Name: "Test"}
		event := NewBaseEvent("test.event", data)

		assert.Equal(t, "test.event", event.Type())
		assert.Equal(t, data, event.Data())
		assert.WithinDuration(t, time.Now(), event.Timestamp(), time.Second)
	})

	t.Run("空数据事件", func(t *testing.T) {
		event := NewBaseEvent("empty.event", nil)

		assert.Equal(t, "empty.event", event.Type())
		assert.Nil(t, event.Data())
	})

	t.Run("复杂数据事件", func(t *testing.T) {
		data := map[string]interface{}{
			"id":     "complex-123",
			"name":   "Complex Event",
			"count":  42,
			"active": true,
			"nested": map[string]string{"key": "value"},
		}
		event := NewBaseEvent("complex.event", data)

		assert.Equal(t, "complex.event", event.Type())
		assert.Equal(t, data, event.Data())
	})

	t.Run("空类型事件", func(t *testing.T) {
		event := NewBaseEvent("", "data")

		assert.Empty(t, event.Type())
		assert.Equal(t, "data", event.Data())
	})

	t.Run("字符串数据", func(t *testing.T) {
		event := NewBaseEvent("string.event", "simple string data")

		assert.Equal(t, "string.event", event.Type())
		assert.Equal(t, "simple string data", event.Data())
	})

	t.Run("数字数据", func(t *testing.T) {
		event := NewBaseEvent("number.event", 42)

		assert.Equal(t, "number.event", event.Type())
		assert.Equal(t, 42, event.Data())
	})

	t.Run("切片数据", func(t *testing.T) {
		data := []string{"item1", "item2", "item3"}
		event := NewBaseEvent("slice.event", data)

		assert.Equal(t, "slice.event", event.Type())
		assert.Equal(t, data, event.Data())
	})
}

// TestBaseEventImmutability 测试事件不可变性（语义上）
func TestBaseEventImmutability(t *testing.T) {
	data := &MockEventData{ID: "123", Name: "Original"}
	event := NewBaseEvent("test.event", data)

	// 原始数据修改不应该影响事件（因为是指针，实际上会修改，但这是设计选择）
	// 这里我们测试时间戳是不变的
	timestamp := event.Timestamp()
	time.Sleep(10 * time.Millisecond)

	// 时间戳应该保持不变
	assert.Equal(t, timestamp, event.Timestamp())
}

// TestBaseEventTimestamp 测试事件时间戳
func TestBaseEventTimestamp(t *testing.T) {
	before := time.Now()
	event := NewBaseEvent("test.event", "data")
	after := time.Now()

	// 时间戳应该在 before 和 after 之间
	assert.True(t, event.Timestamp().After(before) || event.Timestamp().Equal(before))
	assert.True(t, event.Timestamp().Before(after) || event.Timestamp().Equal(after))
}

// TestBaseEventMultipleCreations 测试创建多个事件
func TestBaseEventMultipleCreations(t *testing.T) {
	events := make([]*BaseEvent, 10)
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Millisecond) // 确保时间戳不同
		events[i] = NewBaseEvent("test.event", i)
	}

	// 验证每个事件都有唯一的时间戳（至少是递增的）
	for i := 1; i < len(events); i++ {
		assert.True(t, events[i].Timestamp().After(events[i-1].Timestamp()) ||
			events[i].Timestamp().Equal(events[i-1].Timestamp()),
			"时间戳应该是非递减的")
	}
}

// ==================== MockEvent 测试 ====================

// MockEvent 实现 Event 接口的测试事件
type MockEvent struct {
	eventType string
	data      interface{}
	timestamp time.Time
}

func (m *MockEvent) Type() string {
	return m.eventType
}

func (m *MockEvent) Data() interface{} {
	return m.data
}

func (m *MockEvent) Timestamp() time.Time {
	return m.timestamp
}

// TestMockEvent 测试 MockEvent 实现
func TestMockEvent(t *testing.T) {
	now := time.Now()
	event := &MockEvent{
		eventType: "mock.event",
		data:      map[string]string{"key": "value"},
		timestamp: now,
	}

	assert.Equal(t, "mock.event", event.Type())
	assert.Equal(t, map[string]string{"key": "value"}, event.Data())
	assert.Equal(t, now, event.Timestamp())

	// 验证实现了 Event 接口
	var _ Event = event
}

// TestMockEventWithDifferentTypes 测试不同类型的事件
func TestMockEventWithDifferentTypes(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		data      interface{}
	}{
		{
			name:      "用户创建事件",
			eventType: "user.created",
			data:      map[string]string{"userId": "123", "email": "test@example.com"},
		},
		{
			name:      "订单完成事件",
			eventType: "order.completed",
			data:      map[string]interface{}{"orderId": "456", "amount": 99.99},
		},
		{
			name:      "支付成功事件",
			eventType: "payment.success",
			data:      map[string]string{"paymentId": "789", "status": "success"},
		},
		{
			name:      "系统通知事件",
			eventType: "system.notification",
			data:      map[string]string{"message": "系统维护通知"},
		},
		{
			name:      "空数据事件",
			eventType: "empty.data",
			data:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &MockEvent{
				eventType: tt.eventType,
				data:      tt.data,
				timestamp: time.Now(),
			}

			assert.Equal(t, tt.eventType, event.Type())
			assert.Equal(t, tt.data, event.Data())
			assert.NotZero(t, event.Timestamp())
		})
	}
}

// ==================== MockEventHandler 测试 ====================

// MockEventHandler 测试用的事件处理器
type MockEventHandler struct {
	handledEvents []Event
	returnError   error
}

func (m *MockEventHandler) Handle(ctx context.Context, event Event) error {
	m.handledEvents = append(m.handledEvents, event)
	return m.returnError
}

// TestEventHandler 测试事件处理器
func TestEventHandler(t *testing.T) {
	t.Run("处理成功", func(t *testing.T) {
		handler := &MockEventHandler{returnError: nil}
		event := &MockEvent{
			eventType: "test.event",
			data:      map[string]string{"key": "value"},
			timestamp: time.Now(),
		}
		ctx := context.Background()

		err := handler.Handle(ctx, event)

		require.NoError(t, err)
		assert.Len(t, handler.handledEvents, 1)
		assert.Equal(t, event, handler.handledEvents[0])
	})

	t.Run("处理失败", func(t *testing.T) {
		handler := &MockEventHandler{returnError: errors.New("处理失败")}
		event := &MockEvent{
			eventType: "test.event",
			data:      nil,
			timestamp: time.Now(),
		}
		ctx := context.Background()

		err := handler.Handle(ctx, event)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "处理失败")
	})

	t.Run("处理多个事件", func(t *testing.T) {
		handler := &MockEventHandler{returnError: nil}
		ctx := context.Background()

		for i := 0; i < 5; i++ {
			event := &MockEvent{
				eventType: "test.event",
				data:      map[string]int{"index": i},
				timestamp: time.Now(),
			}
			err := handler.Handle(ctx, event)
			require.NoError(t, err)
		}

		assert.Len(t, handler.handledEvents, 5)
	})

	t.Run("处理不同类型事件", func(t *testing.T) {
		handler := &MockEventHandler{returnError: nil}
		ctx := context.Background()

		eventTypes := []string{"user.created", "user.updated", "user.deleted"}
		for _, eventType := range eventTypes {
			event := NewBaseEvent(eventType, "data")
			err := handler.Handle(ctx, event)
			require.NoError(t, err)
		}

		assert.Len(t, handler.handledEvents, 3)
		assert.Equal(t, "user.created", handler.handledEvents[0].Type())
		assert.Equal(t, "user.updated", handler.handledEvents[1].Type())
		assert.Equal(t, "user.deleted", handler.handledEvents[2].Type())
	})
}

// TestEventHandlerInterface 测试 EventHandler 接口
func TestEventHandlerInterface(t *testing.T) {
	t.Run("接口实现", func(t *testing.T) {
		// 验证 MockEventHandler 实现了 EventHandler 接口
		var _ EventHandler[Event] = &MockEventHandler{}
	})

	t.Run("泛型处理", func(t *testing.T) {
		handler := &MockEventHandler{returnError: nil}
		ctx := context.Background()
		event := NewBaseEvent("test.event", "data")

		var h EventHandler[Event] = handler
		err := h.Handle(ctx, event)

		require.NoError(t, err)
	})
}

// ==================== Event 类型测试 ====================

// TestEventTypes 测试不同类型的事件
func TestEventTypes(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		data      interface{}
	}{
		{
			name:      "用户创建事件",
			eventType: "user.created",
			data:      map[string]string{"userId": "123", "email": "test@example.com"},
		},
		{
			name:      "订单完成事件",
			eventType: "order.completed",
			data:      map[string]interface{}{"orderId": "456", "amount": 99.99},
		},
		{
			name:      "支付成功事件",
			eventType: "payment.success",
			data:      map[string]string{"paymentId": "789", "status": "success"},
		},
		{
			name:      "系统通知事件",
			eventType: "system.notification",
			data:      map[string]string{"message": "系统维护通知"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewBaseEvent(tt.eventType, tt.data)

			assert.Equal(t, tt.eventType, event.Type())
			assert.Equal(t, tt.data, event.Data())
			assert.NotZero(t, event.Timestamp())
		})
	}
}

// TestEventTimestampOrdering 测试事件时间戳顺序
func TestEventTimestampOrdering(t *testing.T) {
	// 创建多个事件，验证时间戳递增
	var timestamps []time.Time

	for i := 0; i < 5; i++ {
		event := NewBaseEvent("test.event", map[string]int{"index": i})
		timestamps = append(timestamps, event.Timestamp())
		time.Sleep(5 * time.Millisecond) // 确保时间戳不同
	}

	// 验证时间戳是非递减的
	for i := 1; i < len(timestamps); i++ {
		assert.True(t, timestamps[i].After(timestamps[i-1]) || timestamps[i].Equal(timestamps[i-1]),
			"时间戳应该是递增的")
	}
}

// ==================== 上下文测试 ====================

// TestEventWithContext 测试带上下文的事件处理
func TestEventWithContext(t *testing.T) {
	t.Run("正常上下文", func(t *testing.T) {
		ctx := context.Background()
		handler := &MockEventHandler{returnError: nil}
		event := NewBaseEvent("test.event", "data")

		err := handler.Handle(ctx, event)
		require.NoError(t, err)
	})

	t.Run("取消的上下文", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		handler := &MockEventHandler{returnError: nil}
		event := NewBaseEvent("test.event", "data")

		// 这里 handler 没有检查上下文，所以不会返回错误
		// 实际实现应该检查 ctx.Err()
		err := handler.Handle(ctx, event)
		require.NoError(t, err)
	})

	t.Run("带值的上下文", func(t *testing.T) {
		type contextKey string
		key := contextKey("request-id")
		ctx := context.WithValue(context.Background(), key, "req-123")

		handler := &MockEventHandler{returnError: nil}
		event := NewBaseEvent("test.event", "data")

		err := handler.Handle(ctx, event)
		require.NoError(t, err)
	})

	t.Run("超时上下文", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()

		handler := &MockEventHandler{returnError: nil}
		event := NewBaseEvent("test.event", "data")

		err := handler.Handle(ctx, event)
		require.NoError(t, err)
	})
}

// ==================== 具体领域事件测试 ====================

// UserCreatedEvent 用户创建事件示例
type UserCreatedEvent struct {
	BaseEvent
	UserID string
	Email  string
}

// TestUserCreatedEvent 测试具体的领域事件
func TestUserCreatedEvent(t *testing.T) {
	now := time.Now()
	baseEvent := NewBaseEvent("user.created", nil)

	event := &UserCreatedEvent{
		BaseEvent: *baseEvent,
		UserID:    "user-123",
		Email:     "test@example.com",
	}

	assert.Equal(t, "user.created", event.Type())
	assert.Equal(t, "user-123", event.UserID)
	assert.Equal(t, "test@example.com", event.Email)
	assert.WithinDuration(t, now, event.Timestamp(), time.Second)
}

// OrderPlacedEvent 订单创建事件示例
type OrderPlacedEvent struct {
	BaseEvent
	OrderID string
	UserID  string
	Amount  float64
}

// TestOrderPlacedEvent 测试订单事件
func TestOrderPlacedEvent(t *testing.T) {
	baseEvent := NewBaseEvent("order.placed", nil)

	event := &OrderPlacedEvent{
		BaseEvent: *baseEvent,
		OrderID:   "order-456",
		UserID:    "user-123",
		Amount:    199.99,
	}

	assert.Equal(t, "order.placed", event.Type())
	assert.Equal(t, "order-456", event.OrderID)
	assert.Equal(t, "user-123", event.UserID)
	assert.Equal(t, 199.99, event.Amount)
}

// PaymentCompletedEvent 支付完成事件示例
type PaymentCompletedEvent struct {
	BaseEvent
	PaymentID string
	OrderID   string
	Status    string
}

// TestPaymentCompletedEvent 测试支付事件
func TestPaymentCompletedEvent(t *testing.T) {
	baseEvent := NewBaseEvent("payment.completed", nil)

	event := &PaymentCompletedEvent{
		BaseEvent: *baseEvent,
		PaymentID: "pay-789",
		OrderID:   "order-456",
		Status:    "success",
	}

	assert.Equal(t, "payment.completed", event.Type())
	assert.Equal(t, "pay-789", event.PaymentID)
	assert.Equal(t, "success", event.Status)
}

// ==================== 复杂场景测试 ====================

// TestEventWithComplexData 测试复杂数据的事件
func TestEventWithComplexData(t *testing.T) {
	type Address struct {
		Street  string
		City    string
		Country string
	}

	type UserProfile struct {
		Name    string
		Email   string
		Address Address
	}

	profile := UserProfile{
		Name:  "John Doe",
		Email: "john@example.com",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
	}

	event := NewBaseEvent("user.profile.updated", profile)

	data, ok := event.Data().(UserProfile)
	require.True(t, ok)
	assert.Equal(t, "John Doe", data.Name)
	assert.Equal(t, "New York", data.Address.City)
}

// TestEventWithSliceData 测试切片数据的事件
func TestEventWithSliceData(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	event := NewBaseEvent("items.created", items)

	data, ok := event.Data().([]string)
	require.True(t, ok)
	assert.Len(t, data, 3)
	assert.Contains(t, data, "item2")
}

// TestEventWithMapData 测试 Map 数据的事件
func TestEventWithMapData(t *testing.T) {
	data := map[string]interface{}{
		"id":       "123",
		"name":     "Test",
		"count":    42,
		"active":   true,
		"tags":     []string{"tag1", "tag2"},
		"metadata": map[string]string{"key": "value"},
	}

	event := NewBaseEvent("complex.event", data)

	result, ok := event.Data().(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "123", result["id"])
	assert.Equal(t, 42, result["count"])
	assert.Equal(t, true, result["active"])
}

// ==================== 性能测试 ====================

// BenchmarkBaseEvent_Create BaseEvent 创建性能测试
func BenchmarkBaseEvent_Create(b *testing.B) {
	data := &MockEventData{ID: "123", Name: "Test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBaseEvent("test.event", data)
	}
}

// BenchmarkEventHandler_Handle Event Handler 性能测试
func BenchmarkEventHandler_Handle(b *testing.B) {
	ctx := context.Background()
	handler := &MockEventHandler{returnError: nil}
	event := NewBaseEvent("test.event", "data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(ctx, event)
	}
}

// BenchmarkEventHandler_HandleParallel 并行事件处理性能测试
func BenchmarkEventHandler_HandleParallel(b *testing.B) {
	ctx := context.Background()
	handler := &MockEventHandler{returnError: nil}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			event := NewBaseEvent("test.event", "data")
			_ = handler.Handle(ctx, event)
		}
	})
}

// BenchmarkEventTimestamp 时间戳获取性能测试
func BenchmarkEventTimestamp(b *testing.B) {
	event := NewBaseEvent("test.event", "data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = event.Timestamp()
	}
}

// ==================== 示例函数 ====================

// ExampleBaseEvent BaseEvent 使用示例
func ExampleBaseEvent() {
	// 创建基础事件
	event := NewBaseEvent("user.created", map[string]string{
		"userId": "123",
		"email":  "test@example.com",
	})

	_ = event.Type()
	_ = event.Data()
	_ = event.Timestamp()
	// Output:
}

// ExampleEventHandler EventHandler 使用示例
func ExampleEventHandler() {
	handler := &MockEventHandler{returnError: nil}
	event := NewBaseEvent("test.event", "data")
	ctx := context.Background()

	_ = handler.Handle(ctx, event)
	// Output:
}

// ExampleUserCreatedEvent UserCreatedEvent 使用示例
func ExampleUserCreatedEvent() {
	baseEvent := NewBaseEvent("user.created", nil)
	event := &UserCreatedEvent{
		BaseEvent: *baseEvent,
		UserID:    "user-123",
		Email:     "test@example.com",
	}

	_ = event.Type()
	// Output:
}
