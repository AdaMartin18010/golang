# EC-038: Command Query Responsibility Segregation (CQRS)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #cqrs #read-model #write-model #event-sourcing
> **权威来源**:
>
> - [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html) - Martin Fowler
> - [CQRS Documents](https://cqrs.files.wordpress.com/2010/11/cqrs_documents.pdf) - Greg Young
> - [Implementing Domain-Driven Design](https://www.oreilly.com/library/view/implementing-domain-driven-design/9780133039900/) - Vaughn Vernon

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在使用统一模型处理读写操作时，由于读写需求差异巨大（读需要高效查询，写需要业务规则验证），导致模型复杂度增加、性能下降，如何解决？

**读写需求差异**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Read vs Write Requirements                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  写操作 (Commands)                        读操作 (Queries)               │
│  ────────────────────────                 ────────────────────────      │
│  • 验证业务规则                            • 高性能查询                   │
│  • 维护数据一致性                          • 复杂过滤和排序               │
│  • 触发领域事件                            • 聚合和统计                   │
│  • 事务边界清晰                            • 多表关联                     │
│  • 更新频率低                              • 读取频率高                   │
│  • 并发冲突处理                            • 最终一致性可接受             │
│                                                                         │
│  统一模型的问题:                                                          │
│  • 为读优化（添加索引、反规范化）影响写性能                                 │
│  • 为写优化（强一致性、验证）导致读复杂                                    │
│  • 领域模型暴露给查询，破坏封装                                           │
│  • 大聚合根加载全部数据，即使只需要一部分                                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**形式化描述**:

```
给定: 模型 M，读操作集合 R，写操作集合 W
约束:
  - R 和 W 有不同性能需求
  - R 和 W 有不同一致性需求
  - R 和 W 有不同数据结构需求
目标: 将 M 分解为 M_read 和 M_write，使得:
  - M_write 优化写入性能和业务规则
  - M_read 优化查询性能
  - M_read 和 M_write 保持最终一致
```

### 1.2 解决方案形式化

**定义 1.1 (CQRS)**
CQRS 将数据存储的读取和写入分离为两个独立的模型：

- **Command Model（写模型）**: 处理命令，执行业务逻辑，验证规则
- **Query Model（读模型）**: 优化查询，物化视图，可反规范化

**形式化表示**:

```
分离:
  Command: C × M_write → M_write × Events
  Query:   Q × M_read → Result

同步:
  Events → Projection → M_read

一致性:
  ∀t: M_write(t) →◇ M_read(t + Δt)
```

**定义 1.2 (投影)**
投影是将写模型产生的事件转换为读模型的过程：

```
Projection: Event* × State → State

读取模型重建:
  M_read = fold(apply_event, initial_state, events)
```

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CQRS Architecture                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                        Command Side                              │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │   Command   │───►│  Command    │───►│    Write Model      │  │   │
│  │  │   API       │    │  Handler    │    │    (Domain Model)   │  │   │
│  │  └─────────────┘    └──────┬──────┘    └──────────┬──────────┘  │   │
│  │                            │                       │             │   │
│  │                            │                       ▼             │   │
│  │                            │              ┌─────────────────┐   │   │
│  │                            │              │  Event Store    │   │   │
│  │                            │              │  (Append-only)  │   │   │
│  │                            │              └────────┬────────┘   │   │
│  │                            │                       │             │   │
│  │                            └───────────────────────┘             │   │
│  │                                    Publish Events                 │   │
│  └────────────────────────────────────┬────────────────────────────┘   │
│                                       │                                 │
│                                       ▼                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Event Bus / Message Queue                    │   │
│  │                                                                  │   │
│  │         ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │   │
│  │         │  Projection │  │  Projection │  │  Projection │       │   │
│  │         │  Handler 1  │  │  Handler 2  │  │  Handler N  │       │   │
│  │         └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │   │
│  └────────────────┼────────────────┼────────────────┼──────────────┘   │
│                   │                │                │                   │
│                   ▼                ▼                ▼                   │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                         Query Side                               │   │
│  │                                                                  │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │   │
│  │  │ Read Model  │  │ Read Model  │  │ Read Model  │              │   │
│  │  │ (MongoDB)   │  │ (Elastic)   │  │  (Redis)    │              │   │
│  │  │             │  │             │  │             │              │   │
│  │  │ • List view │  │ • Search    │  │ • Cache     │              │   │
│  │  │ • Detail    │  │ • Analytics │  │ • Hot data  │              │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │   │
│  │         └────────────────┼────────────────┘                      │   │
│  │                          │                                       │   │
│  │  ┌───────────────────────┴───────────────────────────────────┐  │   │
│  │  │                    Query API                                 │  │   │
│  │  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │  │   │
│  │  │  │   Query     │───►│   Query     │───►│  Response   │     │  │   │
│  │  │  │   Handler   │    │   Processor │    │  (DTO)      │     │  │   │
│  │  │  └─────────────┘    └─────────────┘    └─────────────┘     │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  关键特点:                                                               │
│  • 写模型和读模型完全分离                                                │
│  • 读模型是物化视图，可反规范化                                            │
│  • 事件驱动的最终一致性                                                   │
│  • 读模型可独立扩展和优化                                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心 CQRS 实现

```go
// cqrs/core.go
package cqrs

import (
    "context"
    "fmt"
    "reflect"
    "time"
)

// Command 命令接口
type Command interface {
    CommandName() string
    AggregateID() string
}

// CommandHandler 命令处理器接口
type CommandHandler interface {
    Handle(ctx context.Context, cmd Command) error
}

// Query 查询接口
type Query interface {
    QueryName() string
}

// QueryHandler 查询处理器接口
type QueryHandler interface {
    Handle(ctx context.Context, query Query) (interface{}, error)
}

// Event 事件接口
type Event interface {
    EventName() string
    AggregateID() string
    OccurredAt() time.Time
}

// EventHandler 事件处理器接口
type EventHandler interface {
    Handle(ctx context.Context, event Event) error
}

// EventStore 事件存储接口
type EventStore interface {
    Append(ctx context.Context, aggregateID string, events []Event, expectedVersion int) error
    GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
    GetAllEvents(ctx context.Context, afterPosition int64, batchSize int) ([]Event, int64, error)
}

// ReadModel 读模型接口
type ReadModel interface {
    // Update 更新读模型
    Update(ctx context.Context, event Event) error

    // Query 执行查询
    Query(ctx context.Context, query Query) (interface{}, error)
}

// Bus 消息总线接口
type Bus interface {
    // SendCommand 发送命令
    SendCommand(ctx context.Context, cmd Command) error

    // SendQuery 发送查询
    SendQuery(ctx context.Context, query Query) (interface{}, error)

    // PublishEvent 发布事件
    PublishEvent(ctx context.Context, event Event) error

    // RegisterCommandHandler 注册命令处理器
    RegisterCommandHandler(cmdType reflect.Type, handler CommandHandler)

    // RegisterQueryHandler 注册查询处理器
    RegisterQueryHandler(queryType reflect.Type, handler QueryHandler)

    // RegisterEventHandler 注册事件处理器
    RegisterEventHandler(eventType reflect.Type, handler EventHandler)
}

// CQRS 核心结构
type CQRS struct {
    commandHandlers map[string]CommandHandler
    queryHandlers   map[string]QueryHandler
    eventHandlers   map[string][]EventHandler
    eventStore      EventStore
    bus             Bus
}

// NewCQRS 创建 CQRS 实例
func NewCQRS(eventStore EventStore, bus Bus) *CQRS {
    return &CQRS{
        commandHandlers: make(map[string]CommandHandler),
        queryHandlers:   make(map[string]QueryHandler),
        eventHandlers:   make(map[string][]EventHandler),
        eventStore:      eventStore,
        bus:             bus,
    }
}

// RegisterCommandHandler 注册命令处理器
func (c *CQRS) RegisterCommandHandler(cmd Command, handler CommandHandler) {
    c.commandHandlers[cmd.CommandName()] = handler
    c.bus.RegisterCommandHandler(reflect.TypeOf(cmd), handler)
}

// RegisterQueryHandler 注册查询处理器
func (c *CQRS) RegisterQueryHandler(query Query, handler QueryHandler) {
    c.queryHandlers[query.QueryName()] = handler
    c.bus.RegisterQueryHandler(reflect.TypeOf(query), handler)
}

// RegisterEventHandler 注册事件处理器
func (c *CQRS) RegisterEventHandler(event Event, handler EventHandler) {
    eventName := event.EventName()
    c.eventHandlers[eventName] = append(c.eventHandlers[eventName], handler)
    c.bus.RegisterEventHandler(reflect.TypeOf(event), handler)
}

// ExecuteCommand 执行命令
func (c *CQRS) ExecuteCommand(ctx context.Context, cmd Command) error {
    handler, exists := c.commandHandlers[cmd.CommandName()]
    if !exists {
        return fmt.Errorf("no handler for command: %s", cmd.CommandName())
    }

    return handler.Handle(ctx, cmd)
}

// ExecuteQuery 执行查询
func (c *CQRS) ExecuteQuery(ctx context.Context, query Query) (interface{}, error) {
    handler, exists := c.queryHandlers[query.QueryName()]
    if !exists {
        return nil, fmt.Errorf("no handler for query: %s", query.QueryName())
    }

    return handler.Handle(ctx, query)
}

// DispatchEvent 分发事件
func (c *CQRS) DispatchEvent(ctx context.Context, event Event) error {
    handlers := c.eventHandlers[event.EventName()]

    for _, handler := range handlers {
        if err := handler.Handle(ctx, event); err != nil {
            // 记录但继续处理其他处理器
            fmt.Printf("event handler failed: %v\n", err)
        }
    }

    return nil
}
```

### 2.2 订单聚合根示例

```go
// cqrs/order_aggregate.go
package cqrs

import (
    "context"
    "fmt"
    "time"
)

// OrderAggregate 订单聚合根
type OrderAggregate struct {
    ID         string
    CustomerID string
    Items      []OrderItem
    Total      float64
    Status     string
    Version    int
    Changes    []Event
}

// OrderItem 订单项
type OrderItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

// CreateOrderCommand 创建订单命令
type CreateOrderCommand struct {
    OrderID    string
    CustomerID string
    Items      []OrderItem
}

func (c CreateOrderCommand) CommandName() string { return "CreateOrder" }
func (c CreateOrderCommand) AggregateID() string { return c.OrderID }

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
    OrderID    string
    CustomerID string
    Items      []OrderItem
    Total      float64
    Timestamp  time.Time
}

func (e OrderCreatedEvent) EventName() string    { return "OrderCreated" }
func (e OrderCreatedEvent) AggregateID() string  { return e.OrderID }
func (e OrderCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

// CreateOrderHandler 创建订单命令处理器
type CreateOrderHandler struct {
    eventStore EventStore
}

// NewCreateOrderHandler 创建处理器
func NewCreateOrderHandler(eventStore EventStore) *CreateOrderHandler {
    return &CreateOrderHandler{eventStore: eventStore}
}

// Handle 处理命令
func (h *CreateOrderHandler) Handle(ctx context.Context, cmd Command) error {
    createCmd, ok := cmd.(CreateOrderCommand)
    if !ok {
        return fmt.Errorf("invalid command type")
    }

    // 计算总价
    var total float64
    for _, item := range createCmd.Items {
        total += item.Price * float64(item.Quantity)
    }

    // 创建事件
    event := OrderCreatedEvent{
        OrderID:    createCmd.OrderID,
        CustomerID: createCmd.CustomerID,
        Items:      createCmd.Items,
        Total:      total,
        Timestamp:  time.Now(),
    }

    // 保存到事件存储
    if err := h.eventStore.Append(ctx, createCmd.OrderID, []Event{event}, 0); err != nil {
        return fmt.Errorf("failed to append event: %w", err)
    }

    return nil
}

// OrderReadModel 订单读模型
type OrderReadModel struct {
    ID         string    `json:"id" db:"id"`
    CustomerID string    `json:"customer_id" db:"customer_id"`
    Total      float64   `json:"total" db:"total"`
    Status     string    `json:"status" db:"status"`
    ItemCount  int       `json:"item_count" db:"item_count"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// OrderProjection 订单投影
type OrderProjection struct {
    repository OrderReadRepository
}

// OrderReadRepository 订单读模型仓库
type OrderReadRepository interface {
    Save(ctx context.Context, order *OrderReadModel) error
    GetByID(ctx context.Context, id string) (*OrderReadModel, error)
    GetByCustomer(ctx context.Context, customerID string) ([]*OrderReadModel, error)
}

// NewOrderProjection 创建投影
func NewOrderProjection(repo OrderReadRepository) *OrderProjection {
    return &OrderProjection{repository: repo}
}

// Handle 处理事件
func (p *OrderProjection) Handle(ctx context.Context, event Event) error {
    switch e := event.(type) {
    case OrderCreatedEvent:
        return p.handleOrderCreated(ctx, e)
    default:
        return nil
    }
}

func (p *OrderProjection) handleOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    order := &OrderReadModel{
        ID:         event.OrderID,
        CustomerID: event.CustomerID,
        Total:      event.Total,
        Status:     "PENDING",
        ItemCount:  len(event.Items),
        CreatedAt:  event.Timestamp,
    }

    return p.repository.Save(ctx, order)
}

// GetCustomerOrdersQuery 获取客户订单查询
type GetCustomerOrdersQuery struct {
    CustomerID string
}

func (q GetCustomerOrdersQuery) QueryName() string { return "GetCustomerOrders" }

// GetCustomerOrdersHandler 查询处理器
type GetCustomerOrdersHandler struct {
    repository OrderReadRepository
}

// NewGetCustomerOrdersHandler 创建查询处理器
func NewGetCustomerOrdersHandler(repo OrderReadRepository) *GetCustomerOrdersHandler {
    return &GetCustomerOrdersHandler{repository: repo}
}

// Handle 处理查询
func (h *GetCustomerOrdersHandler) Handle(ctx context.Context, query Query) (interface{}, error) {
    q, ok := query.(GetCustomerOrdersQuery)
    if !ok {
        return nil, fmt.Errorf("invalid query type")
    }

    return h.repository.GetByCustomer(ctx, q.CustomerID)
}
```

### 2.3 内存事件存储实现

```go
// cqrs/memory_event_store.go
package cqrs

import (
    "context"
    "fmt"
    "sync"
)

// StoredEvent 存储的事件
type StoredEvent struct {
    AggregateID string
    Event       Event
    Version     int
}

// MemoryEventStore 内存事件存储
type MemoryEventStore struct {
    events map[string][]StoredEvent
    mu     sync.RWMutex
}

// NewMemoryEventStore 创建内存事件存储
func NewMemoryEventStore() *MemoryEventStore {
    return &MemoryEventStore{
        events: make(map[string][]StoredEvent),
    }
}

// Append 追加事件
func (s *MemoryEventStore) Append(ctx context.Context, aggregateID string, events []Event, expectedVersion int) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    existing := s.events[aggregateID]

    if len(existing) != expectedVersion {
        return fmt.Errorf("concurrency conflict: expected version %d, found %d", expectedVersion, len(existing))
    }

    for i, event := range events {
        stored := StoredEvent{
            AggregateID: aggregateID,
            Event:       event,
            Version:     expectedVersion + i + 1,
        }
        s.events[aggregateID] = append(s.events[aggregateID], stored)
    }

    return nil
}

// GetEvents 获取事件
func (s *MemoryEventStore) GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var events []Event
    for _, stored := range s.events[aggregateID] {
        if stored.Version >= fromVersion {
            events = append(events, stored.Event)
        }
    }

    return events, nil
}

// GetAllEvents 获取所有事件
func (s *MemoryEventStore) GetAllEvents(ctx context.Context, afterPosition int64, batchSize int) ([]Event, int64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var events []Event
    var position int64

    for _, storedEvents := range s.events {
        for _, stored := range storedEvents {
            events = append(events, stored.Event)
            position++
            if len(events) >= batchSize {
                return events, position, nil
            }
        }
    }

    return events, position, nil
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// cqrs/core_test.go
package cqrs

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCQRS_CommandHandling(t *testing.T) {
    eventStore := NewMemoryEventStore()
    bus := NewInMemoryBus()
    cqrs := NewCQRS(eventStore, bus)

    // 注册命令处理器
    handler := NewCreateOrderHandler(eventStore)
    cqrs.RegisterCommandHandler(CreateOrderCommand{}, handler)

    // 执行命令
    cmd := CreateOrderCommand{
        OrderID:    "order-001",
        CustomerID: "customer-001",
        Items: []OrderItem{
            {ProductID: "product-001", Quantity: 2, Price: 10.0},
        },
    }

    err := cqrs.ExecuteCommand(context.Background(), cmd)
    require.NoError(t, err)

    // 验证事件已存储
    events, err := eventStore.GetEvents(context.Background(), "order-001", 0)
    require.NoError(t, err)
    assert.Len(t, events, 1)

    createdEvent, ok := events[0].(OrderCreatedEvent)
    require.True(t, ok)
    assert.Equal(t, "order-001", createdEvent.OrderID)
    assert.Equal(t, 20.0, createdEvent.Total)
}

func TestCQRS_QueryHandling(t *testing.T) {
    // 创建读模型仓库
    repo := NewMemoryOrderReadRepository()

    // 预填充数据
    repo.Save(context.Background(), &OrderReadModel{
        ID:         "order-001",
        CustomerID: "customer-001",
        Total:      100.0,
        Status:     "PENDING",
    })

    eventStore := NewMemoryEventStore()
    bus := NewInMemoryBus()
    cqrs := NewCQRS(eventStore, bus)

    // 注册查询处理器
    queryHandler := NewGetCustomerOrdersHandler(repo)
    cqrs.RegisterQueryHandler(GetCustomerOrdersQuery{}, queryHandler)

    // 执行查询
    query := GetCustomerOrdersQuery{CustomerID: "customer-001"}
    result, err := cqrs.ExecuteQuery(context.Background(), query)

    require.NoError(t, err)
    orders, ok := result.([]*OrderReadModel)
    require.True(t, ok)
    assert.Len(t, orders, 1)
    assert.Equal(t, "order-001", orders[0].ID)
}

// MemoryOrderReadRepository 内存读模型仓库
type MemoryOrderReadRepository struct {
    orders map[string]*OrderReadModel
}

func NewMemoryOrderReadRepository() *MemoryOrderReadRepository {
    return &MemoryOrderReadRepository{
        orders: make(map[string]*OrderReadModel),
    }
}

func (r *MemoryOrderReadRepository) Save(ctx context.Context, order *OrderReadModel) error {
    r.orders[order.ID] = order
    return nil
}

func (r *MemoryOrderReadRepository) GetByID(ctx context.Context, id string) (*OrderReadModel, error) {
    return r.orders[id], nil
}

func (r *MemoryOrderReadRepository) GetByCustomer(ctx context.Context, customerID string) ([]*OrderReadModel, error) {
    var result []*OrderReadModel
    for _, order := range r.orders {
        if order.CustomerID == customerID {
            result = append(result, order)
        }
    }
    return result, nil
}

// InMemoryBus 内存总线
type InMemoryBus struct{}

func NewInMemoryBus() *InMemoryBus {
    return &InMemoryBus{}
}

func (b *InMemoryBus) SendCommand(ctx context.Context, cmd Command) error { return nil }
func (b *InMemoryBus) SendQuery(ctx context.Context, query Query) (interface{}, error) { return nil, nil }
func (b *InMemoryBus) PublishEvent(ctx context.Context, event Event) error { return nil }
func (b *InMemoryBus) RegisterCommandHandler(cmdType interface{}, handler CommandHandler) {}
func (b *InMemoryBus) RegisterQueryHandler(queryType interface{}, handler QueryHandler) {}
func (b *InMemoryBus) RegisterEventHandler(eventType interface{}, handler EventHandler) {}
```

---

## 4. 与其他模式的集成

### 4.1 与 Event Sourcing 的关系

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CQRS with Event Sourcing                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  CQRS 和 Event Sourcing 经常一起使用，但不是必须:                           │
│                                                                         │
│  CQRS + Event Sourcing:                                                 │
│  ┌─────────────┐    Events    ┌─────────────────┐                      │
│  │   Command   │─────────────►│   Event Store   │                      │
│  │   Handler   │              │  (Source of     │                      │
│  │             │              │   Truth)        │                      │
│  └─────────────┘              └────────┬────────┘                      │
│                                        │                                │
│                        ┌───────────────┴───────────────┐                │
│                        ▼                               ▼                │
│              ┌─────────────────┐           ┌─────────────────┐         │
│              │  Aggregate      │           │  Projections    │         │
│              │  Reconstruction │           │  (Read Models)  │         │
│              └─────────────────┘           └─────────────────┘         │
│                                                                         │
│  CQRS without Event Sourcing:                                           │
│  ┌─────────────┐    Update    ┌─────────────────┐                      │
│  │   Command   │─────────────►│    Write DB     │                      │
│  │   Handler   │              │  (Relational)   │                      │
│  └─────────────┘              └─────────────────┘                      │
│                                          │                              │
│                                          │ Sync                         │
│                                          ▼                              │
│                               ┌─────────────────┐                      │
│                               │    Read DB      │                      │
│                               │  (Optimized)    │                      │
│                               └─────────────────┘                      │
│                                                                         │
│  选择建议:                                                               │
│  • 需要完整审计日志 ──► Event Sourcing                                  │
│  • 需要复杂业务逻辑 ──► Event Sourcing                                  │
│  • 简单 CRUD 应用 ──► CQRS only                                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 与 Materialized View 的关系

| CQRS | Materialized View |
|------|-------------------|
| 完整的架构模式 | 数据库特性 |
| 包含命令端 | 仅针对查询优化 |
| 事件驱动同步 | 定时刷新或触发器 |
| 支持多个读模型 | 通常单个视图 |

---

## 5. 决策标准

### 5.1 何时使用 CQRS

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CQRS Decision Tree                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  读/写比例 > 10:1 ? ──────────是────────► CQRS 候选                      │
│       │                                                                  │
│       否                                                                 │
│       │                                                                  │
│       ▼                                                                  │
│  读写需要不同数据模型？ ────是────────► CQRS 候选                         │
│       │                                                                  │
│       否                                                                 │
│       │                                                                  │
│       ▼                                                                  │
│  需要独立扩展读写？ ────────是────────► CQRS 候选                         │
│       │                                                                  │
│       否                                                                 │
│       │                                                                  │
│       ▼                                                                  │
│  读模型需要物化视图？ ──────是────────► CQRS 候选                         │
│       │                                                                  │
│       否                                                                 │
│       │                                                                  │
│       ▼                                                                  │
│  传统架构足够                                                            │
│                                                                         │
│  警告: CQRS 增加系统复杂度，只在必要时使用                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CQRS Implementation Checklist                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 明确定义命令和查询的边界                                              │
│  □ 设计事件契约（如使用 Event Sourcing）                                 │
│  □ 确定读模型的物化策略                                                  │
│  □ 规划最终一致性窗口                                                    │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现命令处理器（聚合根）                                              │
│  □ 实现事件存储（如果使用 Event Sourcing）                               │
│  □ 实现投影处理器                                                        │
│  □ 优化读模型查询性能                                                    │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 不要对所有应用使用 CQRS（过度设计）                                    │
│  ❌ 不要让读模型依赖写模型的一致性                                         │
│  ❌ 不要在投影中执行业务逻辑                                               │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>18KB, 完整形式化 + Go 实现 + 测试)

**相关文档**:

- [EC-016-CQRS-Pattern-Formal.md](./EC-016-CQRS-Pattern-Formal.md)
- [EC-015-Event-Sourcing-Formal.md](./EC-015-Event-Sourcing-Formal.md)
- [EC-039-Domain-Event-Pattern.md](./EC-039-Domain-Event-Pattern.md)
