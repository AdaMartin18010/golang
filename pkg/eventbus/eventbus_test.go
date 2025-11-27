package eventbus

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestEventBus_Subscribe(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	var received bool
	handler := func(ctx context.Context, event Event) error {
		received = true
		return nil
	}

	_, err := eb.Subscribe("test.event", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	event := NewEvent("test.event", "test data")
	err = eb.Publish(event)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// 等待处理
	time.Sleep(100 * time.Millisecond)

	if !received {
		t.Error("Expected event to be received")
	}
}

func TestEventBus_Unsubscribe(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	subID, err := eb.Subscribe("test.event", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	err = eb.Unsubscribe(subID)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// 尝试取消不存在的订阅
	err = eb.Unsubscribe("nonexistent")
	if err != ErrSubscriptionNotFound {
		t.Errorf("Expected ErrSubscriptionNotFound, got %v", err)
	}
}

func TestEventBus_WithFilter(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	var received bool
	handler := func(ctx context.Context, event Event) error {
		received = true
		return nil
	}

	filter := func(event Event) bool {
		if data, ok := event.Data().(string); ok {
			return data == "filtered"
		}
		return false
	}

	_, err := eb.SubscribeWithFilter("test.event", handler, filter)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 发布应该被过滤的事件
	event1 := NewEvent("test.event", "not filtered")
	eb.Publish(event1)
	time.Sleep(50 * time.Millisecond)

	if received {
		t.Error("Expected event to be filtered")
	}

	// 发布应该通过的事件
	event2 := NewEvent("test.event", "filtered")
	eb.Publish(event2)
	time.Sleep(50 * time.Millisecond)

	if !received {
		t.Error("Expected event to pass filter")
	}
}

func TestEventBus_Metrics(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()
	defer eb.Stop()

	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	_, err := eb.Subscribe("test.event", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	event := NewEvent("test.event", "test")
	eb.Publish(event)
	time.Sleep(50 * time.Millisecond)

	metrics := eb.GetMetrics()
	if metrics.TotalEvents != 1 {
		t.Errorf("Expected 1 total event, got %d", metrics.TotalEvents)
	}
	if metrics.HandledEvents != 1 {
		t.Errorf("Expected 1 handled event, got %d", metrics.HandledEvents)
	}
}

func TestEventBus_Concurrent(t *testing.T) {
	eb := NewEventBus(100)
	eb.Start()
	defer eb.Stop()

	var mu sync.Mutex
	count := 0

	handler := func(ctx context.Context, event Event) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	}

	_, err := eb.Subscribe("test.event", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 并发发布事件
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			event := NewEvent("test.event", "test")
			eb.Publish(event)
		}()
	}

	wg.Wait()
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	if count != 10 {
		t.Errorf("Expected 10 events, got %d", count)
	}
	mu.Unlock()
}

func TestEventBus_Stop(t *testing.T) {
	eb := NewEventBus(10)
	eb.Start()

	handler := func(ctx context.Context, event Event) error {
		return nil
	}

	_, err := eb.Subscribe("test.event", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	err = eb.Stop()
	if err != nil {
		t.Fatalf("Failed to stop: %v", err)
	}

	// 停止后发布应该失败
	event := NewEvent("test.event", "test")
	err = eb.Publish(event)
	if err != ErrEventBusStopped {
		t.Errorf("Expected ErrEventBusStopped, got %v", err)
	}
}

func TestBaseEvent(t *testing.T) {
	event := NewEvent("test.event", "test data")

	if event.Type() != "test.event" {
		t.Errorf("Expected type 'test.event', got '%s'", event.Type())
	}

	if event.Data() != "test data" {
		t.Errorf("Expected data 'test data', got '%v'", event.Data())
	}

	event.SetMetadata("key", "value")
	value, ok := event.GetMetadata("key")
	if !ok || value != "value" {
		t.Errorf("Expected metadata value 'value', got '%v'", value)
	}
}
