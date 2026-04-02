# EC-039: Domain Event Pattern (领域事件模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #domain-event #event-driven #ddd #loose-coupling
> **权威来源**:
>
> - [Domain Event](https://martinfowler.com/eaaDev/DomainEvent.html) - Martin Fowler
> - [Implementing Domain-Driven Design](https://www.oreilly.com/library/view/implementing-domain-driven-design/9780133039900/) - Vaughn Vernon
> - [Domain-Driven Design](https://domainlanguage.com/ddd/) - Eric Evans

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在领域驱动设计中，如何捕获和传达领域中发生的重要业务事件，使系统的不同部分能够以松耦合的方式响应这些变化？

**形式化描述**:

```
给定: 领域模型 M 包含聚合根 {A₁, A₂, ..., Aₙ}
给定: 业务操作集合 O 作用于聚合根
问题: 如何在 Aᵢ 发生重要变化时，通知相关方而不引入紧耦合？

约束:
  - 聚合根之间不直接引用
  - 业务规则跨越聚合边界时需要协调
  - 其他子域或外部系统需要知道领域变化
```

**传统方法的局限性**:

```
紧耦合方式（不推荐）:
  OrderService.createOrder() {
    order.save()
    inventoryService.decreaseStock()  // 直接调用，紧耦合
    notificationService.sendEmail()   // 直接调用，紧耦合
    analyticsService.recordEvent()    // 直接调用，紧耦合
  }

问题:
  • 订单服务知道所有下游服务
  • 添加新功能需要修改订单服务
  • 一个下游失败影响订单创建
  • 难以测试
```

### 1.2 解决方案形式化

**定义 1.1 (领域事件)**
领域事件是领域中发生的有意义的离散事件，包含：

- **身份标识**: 事件的唯一标识
- **事件类型**: 事件的业务类型
- **时间戳**: 事件发生的时间
- **聚合标识**: 产生事件的聚合根 ID
- **负载**: 事件的业务数据

**形式化表示**:

```
领域事件 E:
  E = ⟨id, type, aggregate_id, occurred_at, payload, metadata⟩

产生:
  execute(A, op) → A' × [E₁, E₂, ..., Eₙ]

消费:
  on(Eᵢ) → side_effect

特性:
  - 不可变性: ∀E: immutable(E)
  - 时序性: E₁.occurred_at < E₂.occurred_at → E₁ happened before E₂
  - 因果性: E₁ → E₂ 表示 E₁ 导致 E₂
```

**定义 1.2 (事件发布与订阅)**

```
发布-订阅模型:
  Publisher: A → Event → EventBus
  Subscriber: EventBus → Event → Handler

解耦:
  Publisher ⊀ Subscriber
  Subscriber ⊀ Publisher
```

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Domain Event Architecture                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Command Side                                 │   │
│  │                                                                  │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │   Command   │───►│  Aggregate  │───►│   Domain Events     │  │   │
│  │  │   Handler   │    │  Root       │    │   (In-Memory)       │  │   │
│  │  │             │    │             │    │                     │  │   │
│  │  │ CreateOrder │    │ Order       │    │ • OrderCreated      │  │   │
│  │  │ PayOrder    │    │ • create()  │    │ • OrderPaid         │  │   │
│  │  │ ShipOrder   │    │ • pay()     │    │ • OrderShipped      │  │   │
│  │  │             │    │ • ship()    │    │ • OrderCancelled    │  │   │
│  │  └─────────────┘    └──────┬──────┘    └──────────┬──────────┘  │   │
│  │                            │                       │             │   │
│  │                            │ State Changes         │ Publish     │   │
│  └────────────────────────────┼───────────────────────┼─────────────┘   │
│                               │                       │                 │
│  ┌────────────────────────────┴───────────────────────┴─────────────┐   │
│  │                       Event Bus / Message Queue                     │   │
│  │                                                                    │   │
│  │         ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │   │
│  │         │  Topic:     │  │  Topic:     │  │  Topic:     │         │   │
│  │         │  orders     │  │  payments   │  │  shipping   │         │   │
│  │         └──────┬──────┘  └──────┬──────┘  └──────┬──────┘         │   │
│  └────────────────┼────────────────┼────────────────┼────────────────┘   │
│                   │                │                │                     │
│  ┌────────────────┼────────────────┼────────────────┼────────────────┐   │
│  │                ▼                ▼                ▼                │   │
│  │           ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │   │
│  │           │ Inventory   │  │ Analytics   │  │ Notification│       │   │
│  │           │ Handler     │  │ Handler     │  │ Handler     │       │   │
│  │           │             │  │             │  │             │       │   │
│  │           │ ReserveStock│  │ RecordEvent │  │ SendEmail   │       │   │
│  │           └─────────────┘  └─────────────┘  └─────────────┘       │   │
│  │                                                                    │   │
│  │                     Event Handlers (Subscribers)                   │   │
│  └────────────────────────────────────────────────────────────────────┘   │
│                                                                           │
│  关键特性:                                                                 │
│  • 聚合根产生事件但不关心谁消费                                             │
│  • 事件处理器独立演化                                                       │
│  • 支持最终一致性                                                          │
│  • 便于审计和追踪                                                          │
│                                                                           │
└───────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心领域事件实现

```go
// domainevent/core.go
package domainevent

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/google/uuid"
)

// Event 领域事件接口
type Event interface {
    EventID() string
    EventType() string
    AggregateID() string
    AggregateType() string
    OccurredAt() time.Time
    Payload() interface{}
    Metadata() map[string]string
}

// BaseEvent 事件基础结构
type BaseEvent struct {
    ID            string                 `json:"id"`
    Type          string                 `json:"type"`
    AggID         string                 `json:"aggregate_id"`
    AggType       string                 `json:"aggregate_type"`
    Timestamp     time.Time              `json:"timestamp"`
    EventPayload  interface{}            `json:"payload"`
    EventMetadata map[string]string      `json:"metadata"`
}

func (e BaseEvent) EventID() string            { return e.ID }
func (e BaseEvent) EventType() string          { return e.Type }
func (e BaseEvent) AggregateID() string        { return e.AggID }
func (e BaseEvent) AggregateType() string      { return e.AggType }
func (e BaseEvent) OccurredAt() time.Time      { return e.Timestamp }
func (e BaseEvent) Payload() interface{}       { return e.EventPayload }
func (e BaseEvent) Metadata() map[string]string { return e.EventMetadata }

// EventPublisher 事件发布者接口
type EventPublisher interface {
    Publish(ctx context.Context, events ...Event) error
}

// EventSubscriber 事件订阅者接口
type EventSubscriber interface {
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) error

// EventBus 事件总线接口
type EventBus interface {
    EventPublisher
    EventSubscriber
}

// AggregateRoot 聚合根接口
type AggregateRoot interface {
    ID() string
    Type() string
    Version() int
    ApplyEvent(event Event) error
    UncommittedEvents() []Event
    ClearUncommittedEvents()
}

// AggregateBase 聚合根基础
type AggregateBase struct {
    id                 string
    version            int
    uncommittedEvents  []Event
    mu                 sync.RWMutex
}

// NewAggregateBase 创建聚合根基础
func NewAggregateBase(id string) *AggregateBase {
    return &AggregateBase{
        id:                id,
        version:           0,
        uncommittedEvents: make([]Event, 0),
    }
}

func (a *AggregateBase) ID() string      { return a.id }
func (a *AggregateBase) Version() int    { return a.version }

func (a *AggregateBase) UncommittedEvents() []Event {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return append([]Event{}, a.uncommittedEvents...)
}

func (a *AggregateBase) ClearUncommittedEvents() {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.uncommittedEvents = a.uncommittedEvents[:0]
}

// RecordEvent 记录事件
func (a *AggregateBase) RecordEvent(eventType string, aggregateType string, payload interface{}, metadata map[string]string) Event {
    a.mu.Lock()
    defer a.mu.Unlock()

    event := BaseEvent{
        ID:            uuid.New().String(),
        Type:          eventType,
        AggID:         a.id,
        AggType:       aggregateType,
        Timestamp:     time.Now(),
        EventPayload:  payload,
        EventMetadata: metadata,
    }

    a.uncommittedEvents = append(a.uncommittedEvents, event)
    a.version++

    return event
}

// IncrementVersion 增加版本
func (a *AggregateBase) IncrementVersion() {
    a.version++
}

// InMemoryEventBus 内存事件总线
type InMemoryEventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

// NewInMemoryEventBus 创建内存事件总线
func NewInMemoryEventBus() *InMemoryEventBus {
    return &InMemoryEventBus{
        handlers: make(map[string][]EventHandler),
    }
}

// Publish 发布事件
func (b *InMemoryEventBus) Publish(ctx context.Context, events ...Event) error {
    for _, event := range events {
        b.mu.RLock()
        handlers := b.handlers[event.EventType()]
        b.mu.RUnlock()

        for _, handler := range handlers {
            if err := handler(ctx, event); err != nil {
                // 记录但不中断其他处理器
                fmt.Printf("event handler failed: %v\n", err)
            }
        }
    }
    return nil
}

// Subscribe 订阅事件
func (b *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.handlers[eventType] = append(b.handlers[eventType], handler)
    return nil
}

// Unsubscribe 取消订阅
func (b *InMemoryEventBus) Unsubscribe(eventType string, handler EventHandler) error {
    // 简化实现
    return nil
}

// EventStore 事件存储接口
type EventStore interface {
    Save(ctx context.Context, aggregateID string, events []Event, expectedVersion int) error
    Get(ctx context.Context, aggregateID string) ([]Event, error)
}

// Repository 仓储接口
type Repository interface {
    Save(ctx context.Context, aggregate AggregateRoot) error
    Get(ctx context.Context, id string) (AggregateRoot, error)
}

// EventSourcedRepository 事件溯源仓储
type EventSourcedRepository struct {
    eventStore EventStore
    factory    func() AggregateRoot
}

// NewEventSourcedRepository 创建事件溯源仓储
func NewEventSourcedRepository(eventStore EventStore, factory func() AggregateRoot) *Repository {
    repo := &EventSourcedRepository{
        eventStore: eventStore,
        factory:    factory,
    }
    var r Repository = repo
    return &r
}

// Save 保存聚合根
func (r *EventSourcedRepository) Save(ctx context.Context, aggregate AggregateRoot) error {
    events := aggregate.UncommittedEvents()
    if len(events) == 0 {
        return nil
    }

    if err := r.eventStore.Save(ctx, aggregate.ID(), events, aggregate.Version()-len(events)); err != nil {
        return fmt.Errorf("failed to save events: %w", err)
    }

    aggregate.ClearUncommittedEvents()
    return nil
}

// Get 获取聚合根
func (r *EventSourcedRepository) Get(ctx context.Context, id string) (AggregateRoot, error) {
    events, err := r.eventStore.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    if len(events) == 0 {
        return nil, fmt.Errorf("aggregate not found: %s", id)
    }

    aggregate := r.factory()
    for _, event := range events {
        if err := aggregate.ApplyEvent(event); err != nil {
            return nil, fmt.Errorf("failed to apply event: %w", err)
        }
    }

    return aggregate, nil
}
```

### 2.2 订单聚合实现

```go
// domainevent/order_aggregate.go
package domainevent

import (
    "encoding/json"
    "fmt"
    "time"
)

// Order 订单聚合根
type Order struct {
    *AggregateBase
    CustomerID string
    Items      []OrderItem
    Total      float64
    Status     OrderStatus
    Address    ShippingAddress
}

// OrderStatus 订单状态
type OrderStatus string

const (
    OrderStatusPending    OrderStatus = "PENDING"
    OrderStatusPaid       OrderStatus = "PAID"
    OrderStatusShipped    OrderStatus = "SHIPPED"
    OrderStatusDelivered  OrderStatus = "DELIVERED"
    OrderStatusCancelled  OrderStatus = "CANCELLED"
)

// OrderItem 订单项
type OrderItem struct {
    ProductID string  `json:"product_id"`
    Name      string  `json:"name"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// ShippingAddress 配送地址
type ShippingAddress struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
    ZipCode string `json:"zip_code"`
}

// OrderCreatedPayload 订单创建负载
type OrderCreatedPayload struct {
    CustomerID string          `json:"customer_id"`
    Items      []OrderItem     `json:"items"`
    Total      float64         `json:"total"`
    Address    ShippingAddress `json:"address"`
}

// OrderPaidPayload 订单支付负载
type OrderPaidPayload struct {
    PaymentID string  `json:"payment_id"`
    Amount    float64 `json:"amount"`
    Method    string  `json:"method"`
}

// NewOrder 创建新订单
func NewOrder(id string, customerID string, items []OrderItem, address ShippingAddress) *Order {
    order := &Order{
        AggregateBase: NewAggregateBase(id),
        CustomerID:    customerID,
        Items:         items,
        Address:       address,
        Status:        OrderStatusPending,
    }

    // 计算总价
    var total float64
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }
    order.Total = total

    // 记录 OrderCreated 事件
    order.RecordEvent("OrderCreated", "Order", OrderCreatedPayload{
        CustomerID: customerID,
        Items:      items,
        Total:      total,
        Address:    address,
    }, map[string]string{
        "correlation_id": id,
    })

    return order
}

// Type 返回聚合类型
func (o *Order) Type() string {
    return "Order"
}

// Pay 支付订单
func (o *Order) Pay(paymentID string, method string) error {
    if o.Status != OrderStatusPending {
        return fmt.Errorf("cannot pay order with status %s", o.Status)
    }

    o.RecordEvent("OrderPaid", "Order", OrderPaidPayload{
        PaymentID: paymentID,
        Amount:    o.Total,
        Method:    method,
    }, map[string]string{
        "order_id": o.ID(),
    })

    o.Status = OrderStatusPaid
    return nil
}

// Ship 发货
func (o *Order) Ship(trackingNumber string) error {
    if o.Status != OrderStatusPaid {
        return fmt.Errorf("cannot ship order with status %s", o.Status)
    }

    o.RecordEvent("OrderShipped", "Order", map[string]string{
        "tracking_number": trackingNumber,
        "shipped_at":      time.Now().Format(time.RFC3339),
    }, nil)

    o.Status = OrderStatusShipped
    return nil
}

// ApplyEvent 应用事件
func (o *Order) ApplyEvent(event Event) error {
    payload, err := json.Marshal(event.Payload())
    if err != nil {
        return err
    }

    switch event.EventType() {
    case "OrderCreated":
        var data OrderCreatedPayload
        if err := json.Unmarshal(payload, &data); err != nil {
            return err
        }
        o.CustomerID = data.CustomerID
        o.Items = data.Items
        o.Total = data.Total
        o.Address = data.Address
        o.Status = OrderStatusPending

    case "OrderPaid":
        o.Status = OrderStatusPaid

    case "OrderShipped":
        o.Status = OrderStatusShipped

    default:
        return fmt.Errorf("unknown event type: %s", event.EventType())
    }

    o.IncrementVersion()
    return nil
}

// OrderEventHandler 订单事件处理器
type OrderEventHandler struct {
    eventBus EventBus
}

// NewOrderEventHandler 创建处理器
func NewOrderEventHandler(eventBus EventBus) *OrderEventHandler {
    return &OrderEventHandler{eventBus: eventBus}
}

// RegisterHandlers 注册处理器
func (h *OrderEventHandler) RegisterHandlers() {
    h.eventBus.Subscribe("OrderCreated", h.handleOrderCreated)
    h.eventBus.Subscribe("OrderPaid", h.handleOrderPaid)
    h.eventBus.Subscribe("OrderShipped", h.handleOrderShipped)
}

func (h *OrderEventHandler) handleOrderCreated(ctx context.Context, event Event) error {
    fmt.Printf("Handling OrderCreated: %s\n", event.AggregateID())
    // 这里可以执行库存预留、发送通知等操作
    return nil
}

func (h *OrderEventHandler) handleOrderPaid(ctx context.Context, event Event) error {
    fmt.Printf("Handling OrderPaid: %s\n", event.AggregateID())
    // 这里可以执行订单确认、发票生成等操作
    return nil
}

func (h *OrderEventHandler) handleOrderShipped(ctx context.Context, event Event) error {
    fmt.Printf("Handling OrderShipped: %s\n", event.AggregateID())
    // 这里可以执行物流跟踪、客户通知等操作
    return nil
}
```

### 2.3 事件存储实现

```go
// domainevent/memory_event_store.go
package domainevent

import (
    "context"
    "fmt"
    "sync"
)

// MemoryEventStore 内存事件存储
type MemoryEventStore struct {
    events map[string][]StoredEvent
    mu     sync.RWMutex
}

// StoredEvent 存储的事件
type StoredEvent struct {
    Event   Event
    Version int
}

// NewMemoryEventStore 创建内存事件存储
func NewMemoryEventStore() *MemoryEventStore {
    return &MemoryEventStore{
        events: make(map[string][]StoredEvent),
    }
}

// Save 保存事件
func (s *MemoryEventStore) Save(ctx context.Context, aggregateID string, events []Event, expectedVersion int) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    existing := s.events[aggregateID]
    if len(existing) != expectedVersion {
        return fmt.Errorf("concurrency conflict: expected version %d, found %d", expectedVersion, len(existing))
    }

    for i, event := range events {
        stored := StoredEvent{
            Event:   event,
            Version: expectedVersion + i + 1,
        }
        s.events[aggregateID] = append(s.events[aggregateID], stored)
    }

    return nil
}

// Get 获取事件
func (s *MemoryEventStore) Get(ctx context.Context, aggregateID string) ([]Event, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var events []Event
    for _, stored := range s.events[aggregateID] {
        events = append(events, stored.Event)
    }

    return events, nil
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// domainevent/order_test.go
package domainevent

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
    items := []OrderItem{
        {ProductID: "PROD-1", Name: "Product 1", Quantity: 2, Price: 10.0},
        {ProductID: "PROD-2", Name: "Product 2", Quantity: 1, Price: 20.0},
    }

    address := ShippingAddress{
        Street:  "123 Main St",
        City:    "New York",
        Country: "USA",
        ZipCode: "10001",
    }

    order := NewOrder("order-001", "customer-001", items, address)

    require.NotNil(t, order)
    assert.Equal(t, "order-001", order.ID())
    assert.Equal(t, "customer-001", order.CustomerID)
    assert.Equal(t, 40.0, order.Total) // 2*10 + 1*20
    assert.Equal(t, OrderStatusPending, order.Status)

    // 验证事件已记录
    events := order.UncommittedEvents()
    require.Len(t, events, 1)
    assert.Equal(t, "OrderCreated", events[0].EventType())
    assert.Equal(t, "Order", events[0].AggregateType())
}

func TestOrder_Pay(t *testing.T) {
    order := createTestOrder()

    err := order.Pay("payment-001", "credit_card")
    require.NoError(t, err)

    assert.Equal(t, OrderStatusPaid, order.Status)

    // 验证事件已记录
    events := order.UncommittedEvents()
    require.Len(t, events, 2) // OrderCreated + OrderPaid
    assert.Equal(t, "OrderPaid", events[1].EventType())
}

func TestOrder_Pay_InvalidStatus(t *testing.T) {
    order := createTestOrder()

    // 先支付
    order.Pay("payment-001", "credit_card")
    order.ClearUncommittedEvents()

    // 再次支付应该失败
    err := order.Pay("payment-002", "credit_card")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot pay order")
}

func TestOrder_ApplyEvent(t *testing.T) {
    order := &Order{
        AggregateBase: NewAggregateBase("order-002"),
    }

    // 创建模拟事件
    event := BaseEvent{
        ID:        "evt-001",
        Type:      "OrderCreated",
        AggID:     "order-002",
        AggType:   "Order",
        Timestamp: time.Now(),
        EventPayload: OrderCreatedPayload{
            CustomerID: "customer-002",
            Items: []OrderItem{
                {ProductID: "PROD-1", Quantity: 1, Price: 50.0},
            },
            Total: 50.0,
            Address: ShippingAddress{
                Street: "456 Oak St",
                City:   "Boston",
            },
        },
    }

    err := order.ApplyEvent(event)
    require.NoError(t, err)

    assert.Equal(t, "customer-002", order.CustomerID)
    assert.Equal(t, 50.0, order.Total)
    assert.Equal(t, OrderStatusPending, order.Status)
}

func TestMemoryEventStore_SaveAndGet(t *testing.T) {
    store := NewMemoryEventStore()
    ctx := context.Background()

    order := createTestOrder()
    events := order.UncommittedEvents()

    // 保存事件
    err := store.Save(ctx, order.ID(), events, 0)
    require.NoError(t, err)

    // 获取事件
    retrieved, err := store.Get(ctx, order.ID())
    require.NoError(t, err)
    assert.Len(t, retrieved, 1)
    assert.Equal(t, "OrderCreated", retrieved[0].EventType())
}

func TestMemoryEventStore_ConcurrencyConflict(t *testing.T) {
    store := NewMemoryEventStore()
    ctx := context.Background()

    // 先保存一些事件
    store.Save(ctx, "order-003", []Event{
        BaseEvent{ID: "evt-1", Type: "OrderCreated"},
    }, 0)

    // 尝试用错误的版本保存
    err := store.Save(ctx, "order-003", []Event{
        BaseEvent{ID: "evt-2", Type: "OrderPaid"},
    }, 0)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "concurrency conflict")
}

func createTestOrder() *Order {
    items := []OrderItem{
        {ProductID: "PROD-1", Quantity: 2, Price: 10.0},
    }
    address := ShippingAddress{Street: "123 Main St", City: "NYC"}
    return NewOrder("order-001", "customer-001", items, address)
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Event Sourcing 的关系

```
┌─────────────────────────────────────────────────────────────────────────┐
│           Domain Events vs Event Sourcing                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Domain Events:                                                         │
│  • 通知机制："发生了什么"                                                │
│  • 可以单独使用                                                          │
│  • 事件可以是派生的（计算得到）                                           │
│  • 可以丢失旧的（如果不需要审计）                                          │
│                                                                         │
│  Event Sourcing:                                                        │
│  • 状态存储：状态 = fold(events)                                         │
│  • 事件是唯一的状态来源                                                  │
│  • 所有事件必须持久化                                                    │
│  • 不能删除事件                                                          │
│                                                                         │
│  关系:                                                                   │
│  Event Sourcing 使用 Domain Events 作为其核心机制                         │
│  Domain Events 可以不使用 Event Sourcing                                  │
│                                                                         │
│  Domain Events only:                                                    │
│  ┌─────────────┐     Event      ┌─────────────┐                        │
│  │   State     │───────────────►│  Handler    │                        │
│  │   (DB)      │                │  (Side FX)  │                        │
│  └─────────────┘                └─────────────┘                        │
│                                                                         │
│  Event Sourcing:                                                        │
│  ┌─────────────┐     Event      ┌─────────────┐                        │
│  │   Events    │───────────────►│   State     │                        │
│  │   (Store)   │◄───────────────│  (Runtime)  │                        │
│  └─────────────┘   Rebuild      └─────────────┘                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 与 Saga 模式的集成

```go
// domainevent/saga_integration.go
package domainevent

import (
    "context"
    "fmt"
)

// SagaOrchestrator Saga 协调器
type SagaOrchestrator struct {
    eventBus    EventBus
    commands    map[string]SagaCommand
    compensations map[string]SagaCommand
}

// SagaCommand Saga 命令
type SagaCommand func(ctx context.Context, payload interface{}) error

// NewSagaOrchestrator 创建 Saga 协调器
func NewSagaOrchestrator(eventBus EventBus) *SagaOrchestrator {
    return &SagaOrchestrator{
        eventBus:      eventBus,
        commands:      make(map[string]SagaCommand),
        compensations: make(map[string]SagaCommand),
    }
}

// RegisterStep 注册 Saga 步骤
func (s *SagaOrchestrator) RegisterStep(stepName string, cmd SagaCommand, compensation SagaCommand) {
    s.commands[stepName] = cmd
    s.compensations[stepName] = compensation
}

// StartSaga 启动 Saga
func (s *SagaOrchestrator) StartSaga(ctx context.Context, sagaID string, steps []string, initialPayload interface{}) error {
    // 订阅完成事件
    for _, step := range steps {
        eventType := fmt.Sprintf("%sCompleted", step)
        s.eventBus.Subscribe(eventType, func(ctx context.Context, event Event) error {
            return s.handleStepCompleted(ctx, sagaID, steps, event)
        })
    }

    // 执行第一步
    if len(steps) > 0 {
        return s.commands[steps[0]](ctx, initialPayload)
    }

    return nil
}

func (s *SagaOrchestrator) handleStepCompleted(ctx context.Context, sagaID string, steps []string, event Event) error {
    // 找到下一步并执行
    // 简化实现
    return nil
}
```

---

## 5. 决策标准

### 5.1 何时使用领域事件

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Domain Event Decision Tree                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  需要通知其他子域业务变化？ ──是──► 使用领域事件                           │
│       │                                                                  │
│       否                                                                 │
│       │                                                                  │
│       ▼                                                                  │
│  需要支持事件溯源？ ────────是──► 使用领域事件                             │
│       │                                                                  │
│       否                                                                 │
│       │                                                                  │
│       ▼                                                                  │
│  需要审计和追踪？ ────────是──► 使用领域事件                              │
│       │                                                                  │
│       否                                                                 │
│       │                                                                  │
│       ▼                                                                  │
│  不需要领域事件                                                          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Domain Event Checklist                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 识别领域中的重要业务事件                                              │
│  □ 设计事件契约（包含必要的信息）                                         │
│  □ 确定事件粒度和命名规范                                                 │
│  □ 规划事件版本策略                                                       │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现聚合根的事件产生逻辑                                               │
│  □ 实现事件处理器（幂等）                                                 │
│  □ 配置事件总线                                                           │
│  □ 实现事件持久化（如需要）                                               │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 不要将领域事件与集成事件混淆                                           │
│  ❌ 不要在事件中包含过多数据（参考 API 模式）                               │
│  ❌ 不要让事件处理器执行同步操作                                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>17KB, 完整形式化 + Go 实现 + 测试)

**相关文档**:

- [EC-015-Event-Sourcing-Formal.md](./EC-015-Event-Sourcing-Formal.md)
- [EC-008-Saga-Pattern-Formal.md](./EC-008-Saga-Pattern-Formal.md)
- [EC-040-Aggregate-Pattern.md](./EC-040-Aggregate-Pattern.md)
