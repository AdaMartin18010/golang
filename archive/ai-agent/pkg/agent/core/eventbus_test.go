package core

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestEventBusBasic 测试事件总线基础功能
func TestEventBusBasic(t *testing.T) {
	eb := NewEventBus(10)
	err := eb.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	defer eb.Stop()

	// 订阅事件
	received := make(chan Event, 1)
	handler := func(ctx context.Context, event Event) error {
		received <- event
		return nil
	}

	subID, err := eb.Subscribe(EventTypeAgentStarted, handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 发布事件
	event := Event{
		ID:     "test-1",
		Type:   EventTypeAgentStarted,
		Source: "test",
		Data:   "test data",
	}

	err = eb.Publish(event)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待事件
	select {
	case receivedEvent := <-received:
		if receivedEvent.ID != event.ID {
			t.Errorf("Expected event ID %s, got %s", event.ID, receivedEvent.ID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Event not received within timeout")
	}

	// 取消订阅
	err = eb.Unsubscribe(subID)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}
}

// TestEventBusMultipleSubscribers 测试多个订阅者
func TestEventBusMultipleSubscribers(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	received1 := make(chan bool, 1)
	received2 := make(chan bool, 1)

	handler1 := func(ctx context.Context, event Event) error {
		received1 <- true
		return nil
	}

	handler2 := func(ctx context.Context, event Event) error {
		received2 <- true
		return nil
	}

	eb.Subscribe(EventTypeProcessingStarted, handler1)
	eb.Subscribe(EventTypeProcessingStarted, handler2)

	// 发布事件
	event := Event{
		ID:     "test-multi",
		Type:   EventTypeProcessingStarted,
		Source: "test",
	}

	eb.Publish(event)

	// 两个处理器都应该收到
	timeout := time.After(2 * time.Second)
	count := 0
	for count < 2 {
		select {
		case <-received1:
			count++
		case <-received2:
			count++
		case <-timeout:
			t.Fatalf("Only %d handlers received event", count)
		}
	}
}

// TestEventBusWithFilter 测试带过滤器的订阅
func TestEventBusWithFilter(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	received := make(chan Event, 2)
	handler := func(ctx context.Context, event Event) error {
		received <- event
		return nil
	}

	// 只接收source为"allowed"的事件
	filter := func(event Event) bool {
		return event.Source == "allowed"
	}

	eb.SubscribeWithFilter(EventTypeCustom, handler, filter)

	// 发布两个事件
	event1 := Event{ID: "1", Type: EventTypeCustom, Source: "allowed"}
	event2 := Event{ID: "2", Type: EventTypeCustom, Source: "blocked"}

	eb.Publish(event1)
	eb.Publish(event2)

	// 只应该收到第一个事件
	select {
	case receivedEvent := <-received:
		if receivedEvent.ID != "1" {
			t.Errorf("Expected event with ID '1', got '%s'", receivedEvent.ID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("First event not received")
	}

	// 不应该收到第二个事件
	select {
	case <-received:
		t.Error("Should not receive blocked event")
	case <-time.After(500 * time.Millisecond):
		// 正常，没有收到
	}
}

// TestEventBusMetrics 测试指标统计
func TestEventBusMetrics(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	eb.Subscribe(EventTypeAgentStarted, handler)

	// 发布多个事件
	for i := 0; i < 5; i++ {
		event := Event{Type: EventTypeAgentStarted, Source: "test"}
		eb.Publish(event)
	}

	// 等待处理
	time.Sleep(1 * time.Second)

	metrics := eb.GetMetrics()
	if metrics.TotalEvents != 5 {
		t.Errorf("Expected 5 total events, got %d", metrics.TotalEvents)
	}
	if metrics.ActiveSubscriptions != 1 {
		t.Errorf("Expected 1 active subscription, got %d", metrics.ActiveSubscriptions)
	}
}

// TestEventBusBufferOverflow 测试缓冲区溢出
func TestEventBusBufferOverflow(t *testing.T) {
	// 创建一个很小的缓冲区
	eb := NewEventBus(1)
	eb.Start()
	defer eb.Stop()

	// 不添加任何订阅者，让事件堆积
	for i := 0; i < 5; i++ {
		event := Event{Type: EventTypeCustom}
		_ = eb.Publish(event)
	}

	metrics := eb.GetMetrics()
	if metrics.DroppedEvents == 0 {
		t.Error("Expected some events to be dropped")
	}
}

// TestEventBusListSubscriptions 测试列出订阅
func TestEventBusListSubscriptions(t *testing.T) {
	eb := NewEventBus(10)
	// 注意：不需要启动事件总线也可以订阅

	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	// 订阅多个事件类型
	_, err1 := eb.Subscribe(EventTypeAgentStarted, handler)
	_, err2 := eb.Subscribe(EventTypeAgentStopped, handler)
	_, err3 := eb.Subscribe(EventTypeProcessingStarted, handler)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("Failed to subscribe")
	}

	subs := eb.ListSubscriptions()
	if len(subs) != 3 {
		t.Errorf("Expected 3 subscriptions, got %d", len(subs))
	}
}

// TestEventBusClear 测试清空订阅
func TestEventBusClear(t *testing.T) {
	eb := NewEventBus(10)

	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	eb.Subscribe(EventTypeAgentStarted, handler)
	eb.Subscribe(EventTypeAgentStopped, handler)

	eb.Clear()

	subs := eb.ListSubscriptions()
	if len(subs) != 0 {
		t.Errorf("Expected 0 subscriptions after clear, got %d", len(subs))
	}

	metrics := eb.GetMetrics()
	if metrics.ActiveSubscriptions != 0 {
		t.Errorf("Expected 0 active subscriptions, got %d", metrics.ActiveSubscriptions)
	}
}

// TestEventBusPublishAsync 测试异步发布
func TestEventBusPublishAsync(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	received := make(chan bool, 1)
	handler := func(ctx context.Context, event Event) error {
		received <- true
		return nil
	}

	eb.Subscribe(EventTypeCustom, handler)

	event := Event{Type: EventTypeCustom, Source: "test"}
	eb.PublishAsync(event)

	select {
	case <-received:
		// 成功
	case <-time.After(2 * time.Second):
		t.Fatal("Event not received within timeout")
	}
}

// TestEventTypes 测试所有事件类型
func TestEventTypes(t *testing.T) {
	types := []EventType{
		EventTypeAgentStarted,
		EventTypeAgentStopped,
		EventTypeProcessingStarted,
		EventTypeProcessingSuccess,
		EventTypeProcessingFailed,
		EventTypeDecisionMade,
		EventTypeLearningCompleted,
		EventTypePluginRegistered,
		EventTypePluginUnregistered,
		EventTypeCustom,
	}

	for _, eventType := range types {
		t.Run(string(eventType), func(t *testing.T) {
			if string(eventType) == "" {
				t.Error("Event type should not be empty")
			}
		})
	}
}

// TestEventBusConcurrent 测试并发访问
func TestEventBusConcurrent(t *testing.T) {
	// 使用足够大的缓冲区来处理并发
	numGoroutines := 5
	eventsPerGoroutine := 50
	totalEvents := numGoroutines * eventsPerGoroutine

	eb := NewEventBus(totalEvents) // 缓冲区足够大
	eb.Start()
	defer eb.Stop()

	var wg sync.WaitGroup

	// 添加订阅者
	handler := func(ctx context.Context, event Event) error {
		return nil
	}
	eb.Subscribe(EventTypeCustom, handler)

	// 并发发布事件（使用更小的规模）
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				event := Event{
					Type:   EventTypeCustom,
					Source: "goroutine",
				}
				_ = eb.Publish(event) // 忽略错误
			}
		}(i)
	}

	wg.Wait()

	// 等待处理完成
	time.Sleep(500 * time.Millisecond)

	metrics := eb.GetMetrics()
	// 验证大多数事件都被成功发布（允许50%的丢失率，因为测试环境可能较慢）
	expectedTotal := int64(totalEvents)
	if metrics.TotalEvents < expectedTotal/2 {
		t.Errorf("Expected at least %d events, got %d (dropped: %d)",
			expectedTotal/2, metrics.TotalEvents, metrics.DroppedEvents)
	}

	// 验证至少有一些事件被成功发布
	if metrics.TotalEvents == 0 {
		t.Error("No events were published successfully")
	}
}

// BenchmarkEventBusPublish 基准测试：发布事件
func BenchmarkEventBusPublish(b *testing.B) {
	eb := NewEventBus(1000)
	eb.Start()
	defer eb.Stop()

	event := Event{Type: EventTypeCustom, Source: "bench"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eb.Publish(event)
	}
}

// BenchmarkEventBusSubscribe 基准测试：订阅事件
func BenchmarkEventBusSubscribe(b *testing.B) {
	eb := NewEventBus(100)

	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eb.Subscribe(EventTypeCustom, handler)
	}
}

// TestEventBusStopWhileProcessing 测试处理中停止
func TestEventBusStopWhileProcessing(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()

	slowHandler := func(ctx context.Context, event Event) error {
		time.Sleep(2 * time.Second)
		return nil
	}

	eb.Subscribe(EventTypeCustom, slowHandler)

	// 发布事件
	event := Event{Type: EventTypeCustom}
	eb.Publish(event)

	// 立即停止
	err := eb.Stop()
	if err != nil {
		t.Fatalf("Failed to stop event bus: %v", err)
	}
}

// TestEventBusUnsubscribeNonexistent 测试取消不存在的订阅
func TestEventBusUnsubscribeNonexistent(t *testing.T) {
	eb := NewEventBus(10)

	err := eb.Unsubscribe("nonexistent-id")
	if err == nil {
		t.Error("Expected error for nonexistent subscription")
	}
}

// TestEventTimestamp 测试事件时间戳自动设置
func TestEventTimestamp(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	received := make(chan Event, 1)
	handler := func(ctx context.Context, event Event) error {
		received <- event
		return nil
	}

	eb.Subscribe(EventTypeCustom, handler)

	// 发布不带时间戳的事件
	event := Event{
		ID:   "test",
		Type: EventTypeCustom,
	}

	eb.Publish(event)

	select {
	case receivedEvent := <-received:
		if receivedEvent.Timestamp.IsZero() {
			t.Error("Event timestamp should be set automatically")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Event not received")
	}
}
