package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// 事件总线 - Event Bus
// =============================================================================

// EventType 事件类型
type EventType string

const (
	EventTypeAgentStarted       EventType = "agent.started"
	EventTypeAgentStopped       EventType = "agent.stopped"
	EventTypeProcessingStarted  EventType = "processing.started"
	EventTypeProcessingSuccess  EventType = "processing.success"
	EventTypeProcessingFailed   EventType = "processing.failed"
	EventTypeDecisionMade       EventType = "decision.made"
	EventTypeLearningCompleted  EventType = "learning.completed"
	EventTypePluginRegistered   EventType = "plugin.registered"
	EventTypePluginUnregistered EventType = "plugin.unregistered"
	EventTypeCustom             EventType = "custom"
)

// Event 事件结构
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Source    string                 `json:"source"`    // 事件来源
	Data      interface{}            `json:"data"`      // 事件数据
	Metadata  map[string]interface{} `json:"metadata"`  // 元数据
	Timestamp time.Time              `json:"timestamp"` // 时间戳
}

// EventHandler 事件处理器函数类型
type EventHandler func(ctx context.Context, event Event) error

// EventSubscription 事件订阅
type EventSubscription struct {
	ID           string
	EventType    EventType
	Handler      EventHandler
	Filter       EventFilter // 事件过滤器
	CreatedAt    time.Time
	handledCount int64
}

// EventFilter 事件过滤器函数类型
type EventFilter func(event Event) bool

// EventBus 事件总线
type EventBus struct {
	subscriptions map[EventType][]*EventSubscription
	subIndex      map[string]*EventSubscription // 订阅ID索引
	mu            sync.RWMutex
	eventChan     chan Event
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	bufferSize    int
	metrics       EventBusMetrics
	nextSubID     int64 // 订阅ID计数器
}

// EventBusMetrics 事件总线指标
type EventBusMetrics struct {
	TotalEvents         int64
	HandledEvents       int64
	FailedEvents        int64
	DroppedEvents       int64
	ActiveSubscriptions int
	mu                  sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus(bufferSize int) *EventBus {
	if bufferSize <= 0 {
		bufferSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventBus{
		subscriptions: make(map[EventType][]*EventSubscription),
		subIndex:      make(map[string]*EventSubscription),
		eventChan:     make(chan Event, bufferSize),
		ctx:           ctx,
		cancel:        cancel,
		bufferSize:    bufferSize,
	}
}

// Start 启动事件总线
func (eb *EventBus) Start() error {
	eb.wg.Add(1)
	go eb.processEvents()
	return nil
}

// Stop 停止事件总线
func (eb *EventBus) Stop() error {
	eb.cancel()
	eb.wg.Wait()
	close(eb.eventChan)
	return nil
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) (string, error) {
	return eb.SubscribeWithFilter(eventType, handler, nil)
}

// SubscribeWithFilter 订阅事件（带过滤器）
func (eb *EventBus) SubscribeWithFilter(eventType EventType, handler EventHandler, filter EventFilter) (string, error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// 使用原子递增的计数器确保ID唯一性
	eb.nextSubID++
	sub := &EventSubscription{
		ID:        fmt.Sprintf("sub_%d", eb.nextSubID),
		EventType: eventType,
		Handler:   handler,
		Filter:    filter,
		CreatedAt: time.Now(),
	}

	eb.subscriptions[eventType] = append(eb.subscriptions[eventType], sub)
	eb.subIndex[sub.ID] = sub
	eb.metrics.ActiveSubscriptions++

	return sub.ID, nil
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(subscriptionID string) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	sub, exists := eb.subIndex[subscriptionID]
	if !exists {
		return fmt.Errorf("subscription %s not found", subscriptionID)
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
	eb.metrics.ActiveSubscriptions--

	return nil
}

// Publish 发布事件
func (eb *EventBus) Publish(event Event) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	select {
	case eb.eventChan <- event:
		eb.metrics.mu.Lock()
		eb.metrics.TotalEvents++
		eb.metrics.mu.Unlock()
		return nil
	case <-eb.ctx.Done():
		return fmt.Errorf("event bus is stopped")
	default:
		// 缓冲区满，丢弃事件
		eb.metrics.mu.Lock()
		eb.metrics.DroppedEvents++
		eb.metrics.mu.Unlock()
		return fmt.Errorf("event buffer is full, event dropped")
	}
}

// PublishAsync 异步发布事件（不阻塞）
func (eb *EventBus) PublishAsync(event Event) {
	go func() {
		_ = eb.Publish(event)
	}()
}

// processEvents 处理事件（内部goroutine）
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
func (eb *EventBus) handleEvent(event Event) {
	eb.mu.RLock()
	subs := eb.subscriptions[event.Type]
	eb.mu.RUnlock()

	for _, sub := range subs {
		// 应用过滤器
		if sub.Filter != nil && !sub.Filter(event) {
			continue
		}

		// 异步执行处理器
		go func(subscription *EventSubscription) {
			ctx, cancel := context.WithTimeout(eb.ctx, 5*time.Second)
			defer cancel()

			if err := subscription.Handler(ctx, event); err != nil {
				eb.metrics.mu.Lock()
				eb.metrics.FailedEvents++
				eb.metrics.mu.Unlock()
			} else {
				subscription.handledCount++
				eb.metrics.mu.Lock()
				eb.metrics.HandledEvents++
				eb.metrics.mu.Unlock()
			}
		}(sub)
	}
}

// GetMetrics 获取指标
func (eb *EventBus) GetMetrics() EventBusMetrics {
	eb.metrics.mu.RLock()
	defer eb.metrics.mu.RUnlock()

	return EventBusMetrics{
		TotalEvents:         eb.metrics.TotalEvents,
		HandledEvents:       eb.metrics.HandledEvents,
		FailedEvents:        eb.metrics.FailedEvents,
		DroppedEvents:       eb.metrics.DroppedEvents,
		ActiveSubscriptions: eb.metrics.ActiveSubscriptions,
	}
}

// ListSubscriptions 列出所有订阅
func (eb *EventBus) ListSubscriptions() []EventSubscription {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	result := make([]EventSubscription, 0, len(eb.subIndex))
	for _, sub := range eb.subIndex {
		result = append(result, *sub)
	}

	return result
}

// Clear 清空所有订阅
func (eb *EventBus) Clear() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscriptions = make(map[EventType][]*EventSubscription)
	eb.subIndex = make(map[string]*EventSubscription)
	eb.metrics.ActiveSubscriptions = 0
}
