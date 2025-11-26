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
	ErrEventBusStopped = errors.New("event bus is stopped")
	// ErrSubscriptionNotFound 订阅未找到
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

// Event 事件接口
type Event interface {
	Type() string
	Data() interface{}
	Timestamp() time.Time
}

// BaseEvent 基础事件实现
type BaseEvent struct {
	eventType string
	data      interface{}
	timestamp time.Time
	metadata  map[string]interface{}
}

// NewEvent 创建事件
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
func (e *BaseEvent) SetMetadata(key string, value interface{}) {
	e.metadata[key] = value
}

// GetMetadata 获取元数据
func (e *BaseEvent) GetMetadata(key string) (interface{}, bool) {
	value, ok := e.metadata[key]
	return value, ok
}

// Handler 事件处理器函数类型
type Handler func(ctx context.Context, event Event) error

// Filter 事件过滤器函数类型
type Filter func(event Event) bool

// Subscription 事件订阅
type Subscription struct {
	ID        string
	EventType string
	Handler   Handler
	Filter    Filter
	CreatedAt time.Time
}

// EventBus 事件总线
type EventBus struct {
	subscriptions map[string][]*Subscription
	subIndex      map[string]*Subscription
	eventChan     chan Event
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	bufferSize    int
	metrics       Metrics
	nextSubID     int64
}

// Metrics 事件总线指标
type Metrics struct {
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
		subscriptions: make(map[string][]*Subscription),
		subIndex:      make(map[string]*Subscription),
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
func (eb *EventBus) Subscribe(eventType string, handler Handler) (string, error) {
	return eb.SubscribeWithFilter(eventType, handler, nil)
}

// SubscribeWithFilter 订阅事件（带过滤器）
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
func (eb *EventBus) Clear() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscriptions = make(map[string][]*Subscription)
	eb.subIndex = make(map[string]*Subscription)

	eb.metrics.mu.Lock()
	eb.metrics.ActiveSubscriptions = 0
	eb.metrics.mu.Unlock()
}
