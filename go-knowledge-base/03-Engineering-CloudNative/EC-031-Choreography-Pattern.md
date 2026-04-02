# EC-031: Choreography Pattern (编舞模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #choreography #event-driven #decentralized #saga #microservices
> **权威来源**:
>
> - [Choreography Pattern](https://microservices.io/patterns/data/saga.html) - Chris Richardson
> - [Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/) - Hohpe & Woolf
> - [Designing Event-Driven Systems](https://www.oreilly.com/library/view/designing-event-driven-systems/9781492038252/) - Ben Stopford
> - [Building Microservices](https://www.oreilly.com/library/view/building-microservices-2nd/9781492034018/) - Sam Newman

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在分布式微服务架构中，如何协调跨多个服务的业务事务而不引入单点故障和紧耦合？

**形式化描述**:

```
给定: 服务集合 S = {S₁, S₂, ..., Sₙ}
给定: 业务事务 T = {t₁, t₂, ..., tₘ}，其中每个 tᵢ 由某个 Sⱼ 执行
约束:
  - 避免中央协调器（防止单点故障）
  - 最小化服务间耦合
  - 保证最终一致性
目标: 找到协调函数 C: T × S → Event，使得事务原子性在分布式环境下得以保持
```

**反模式**:

- 同步编排：服务直接调用其他服务，形成调用链
- 共享数据库：多个服务直接访问同一数据库
- 分布式事务（2PC）：使用两阶段提交，阻塞且难以扩展

### 1.2 解决方案形式化

**定义 1.1 (编舞模式)**
编舞是一种去中心化的协作模式，其中每个服务：

1. 执行本地事务
2. 发布领域事件
3. 订阅相关事件
4. 响应事件执行后续操作

**形式化表示**:

```
设服务 Sᵢ 执行操作 Oᵢ:
  Sᵢ: Oᵢ → Eventᵢ
  Sᵢ₊₁: Subscribe(Eventᵢ) → Oᵢ₊₁ → Eventᵢ₊₁

事务执行流程:
  Start → O₁ → Event₁ → O₂ → Event₂ → ... → Oₙ → Eventₙ → Complete
```

**定义 1.2 (事件契约)**
事件是服务间通信的唯一契约：

```go
type DomainEvent interface {
    EventID()      string
    EventType()    string
    AggregateID()  string
    OccurredAt()   time.Time
    Payload()      interface{}
}
```

### 1.3 状态机模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Choreography State Machine                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌─────────┐     Event1      ┌─────────┐     Event2      ┌─────────┐  │
│   │  Idle   │ ───────────────►│ Active  │ ───────────────►│Processing│ │
│   └─────────┘                 └────┬────┘                 └────┬────┘  │
│        ▲                           │                          │       │
│        │                           ▼                          ▼       │
│        │                     ┌─────────┐                 ┌─────────┐   │
│        └─────────────────────│ Waiting │◄────────────────│  Done   │   │
│                              │  Event  │                 └────┬────┘   │
│                              └─────────┘                      │        │
│                                   ▲                           │        │
│                                   └───────────────────────────┘        │
│                                                                         │
│   Failure Transitions:                                                  │
│   ┌─────────┐    Compensate    ┌─────────┐                             │
│   │ Failed  │◄─────────────────│ Active  │                             │
│   └────┬────┘                  └─────────┘                             │
│        │                                                               │
│        ▼                                                               │
│   ┌─────────┐                                                          │
│   │Compensated│                                                         │
│   └─────────┘                                                          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心组件实现

```go
// choreography/core.go
package choreography

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
    GetID() string
    GetType() string
    GetAggregateID() string
    GetTimestamp() time.Time
    GetPayload() []byte
}

// BaseEvent 事件基础结构
type BaseEvent struct {
    ID          string          `json:"id"`
    Type        string          `json:"type"`
    AggregateID string          `json:"aggregate_id"`
    Timestamp   time.Time       `json:"timestamp"`
    Payload     json.RawMessage `json:"payload"`
    Metadata    map[string]string `json:"metadata"`
}

func (e BaseEvent) GetID() string          { return e.ID }
func (e BaseEvent) GetType() string        { return e.Type }
func (e BaseEvent) GetAggregateID() string { return e.AggregateID }
func (e BaseEvent) GetTimestamp() time.Time { return e.Timestamp }
func (e BaseEvent) GetPayload() []byte     { return e.Payload }

// EventBus 事件总线接口
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) error

// ChoreographyStep 编舞步骤
type ChoreographyStep struct {
    Name            string
    ServiceName     string
    Execute         func(ctx context.Context, event Event) (Event, error)
    Compensate      func(ctx context.Context, event Event) error
    SubscribedEvents []string
    PublishedEvent   string
}

// StepResult 步骤执行结果
type StepResult struct {
    StepName  string
    Success   bool
    Event     Event
    Error     error
    Timestamp time.Time
}

// SagaInstance Saga 实例
type SagaInstance struct {
    ID          string
    Type        string
    Status      SagaStatus
    Steps       []StepResult
    StartedAt   time.Time
    CompletedAt *time.Time
    mu          sync.RWMutex
}

// SagaStatus Saga 状态
type SagaStatus int

const (
    SagaStatusPending SagaStatus = iota
    SagaStatusRunning
    SagaStatusCompleted
    SagaStatusCompensating
    SagaStatusCompensated
    SagaStatusFailed
)

func (s SagaStatus) String() string {
    switch s {
    case SagaStatusPending:
        return "PENDING"
    case SagaStatusRunning:
        return "RUNNING"
    case SagaStatusCompleted:
        return "COMPLETED"
    case SagaStatusCompensating:
        return "COMPENSATING"
    case SagaStatusCompensated:
        return "COMPENSATED"
    case SagaStatusFailed:
        return "FAILED"
    default:
        return "UNKNOWN"
    }
}

// SagaStore Saga 存储接口
type SagaStore interface {
    Save(ctx context.Context, saga *SagaInstance) error
    Load(ctx context.Context, sagaID string) (*SagaInstance, error)
    UpdateStatus(ctx context.Context, sagaID string, status SagaStatus) error
    AppendStep(ctx context.Context, sagaID string, result StepResult) error
}

// Choreographer 编舞协调器
type Choreographer struct {
    eventBus    EventBus
    sagaStore   SagaStore
    steps       map[string]*ChoreographyStep
    handlers    map[string][]EventHandler
    mu          sync.RWMutex
    logger      Logger
}

// Logger 日志接口
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// NewChoreographer 创建编舞器
func NewChoreographer(eventBus EventBus, sagaStore SagaStore, logger Logger) *Choreographer {
    c := &Choreographer{
        eventBus:  eventBus,
        sagaStore: sagaStore,
        steps:     make(map[string]*ChoreographyStep),
        handlers:  make(map[string][]EventHandler),
        logger:    logger,
    }
    return c
}

// RegisterStep 注册编舞步骤
func (c *Choreographer) RegisterStep(step *ChoreographyStep) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if _, exists := c.steps[step.Name]; exists {
        return fmt.Errorf("step %s already registered", step.Name)
    }

    c.steps[step.Name] = step

    // 为每个订阅的事件注册处理器
    for _, eventType := range step.SubscribedEvents {
        handler := c.createHandler(step)
        if err := c.eventBus.Subscribe(eventType, handler); err != nil {
            return fmt.Errorf("failed to subscribe to %s: %w", eventType, err)
        }
        c.handlers[eventType] = append(c.handlers[eventType], handler)
    }

    c.logger.Info("registered choreography step",
        Field{"step", step.Name},
        Field{"service", step.ServiceName},
    )

    return nil
}

func (c *Choreographer) createHandler(step *ChoreographyStep) EventHandler {
    return func(ctx context.Context, event Event) error {
        sagaID := event.GetMetadata("saga_id")
        if sagaID == "" {
            sagaID = uuid.New().String()
        }

        // 加载或创建 Saga 实例
        saga, err := c.sagaStore.Load(ctx, sagaID)
        if err != nil {
            saga = &SagaInstance{
                ID:        sagaID,
                Type:      event.GetMetadata("saga_type"),
                Status:    SagaStatusRunning,
                Steps:     []StepResult{},
                StartedAt: time.Now(),
            }
            if err := c.sagaStore.Save(ctx, saga); err != nil {
                return fmt.Errorf("failed to create saga: %w", err)
            }
        }

        // 执行步骤
        c.logger.Info("executing choreography step",
            Field{"saga_id", sagaID},
            Field{"step", step.Name},
            Field{"event_type", event.GetType()},
        )

        newEvent, err := step.Execute(ctx, event)

        result := StepResult{
            StepName:  step.Name,
            Success:   err == nil,
            Event:     newEvent,
            Error:     err,
            Timestamp: time.Now(),
        }

        if err := c.sagaStore.AppendStep(ctx, sagaID, result); err != nil {
            c.logger.Error("failed to append step result",
                Field{"saga_id", sagaID},
                Field{"error", err},
            )
        }

        if err != nil {
            // 执行失败，触发补偿
            c.logger.Error("step execution failed, initiating compensation",
                Field{"saga_id", sagaID},
                Field{"step", step.Name},
                Field{"error", err},
            )
            return c.compensate(ctx, saga, step)
        }

        // 发布新事件
        if newEvent != nil {
            if err := c.eventBus.Publish(ctx, newEvent); err != nil {
                c.logger.Error("failed to publish event",
                    Field{"saga_id", sagaID},
                    Field{"error", err},
                )
                return err
            }
        }

        return nil
    }
}

func (c *Choreographer) compensate(ctx context.Context, saga *SagaInstance, failedStep *ChoreographyStep) error {
    c.sagaStore.UpdateStatus(ctx, saga.ID, SagaStatusCompensating)

    // 按相反顺序执行补偿
    for i := len(saga.Steps) - 1; i >= 0; i-- {
        step := saga.Steps[i]
        if step.Success && step.StepName != failedStep.Name {
            if s, exists := c.steps[step.StepName]; exists && s.Compensate != nil {
                if err := s.Compensate(ctx, step.Event); err != nil {
                    c.logger.Error("compensation failed",
                        Field{"saga_id", saga.ID},
                        Field{"step", step.StepName},
                        Field{"error", err},
                    )
                    // 补偿失败需要人工介入
                    return fmt.Errorf("compensation failed for step %s: %w", step.StepName, err)
                }
            }
        }
    }

    c.sagaStore.UpdateStatus(ctx, saga.ID, SagaStatusCompensated)
    return nil
}

func (e BaseEvent) GetMetadata(key string) string {
    if e.Metadata == nil {
        return ""
    }
    return e.Metadata[key]
}
```

### 2.2 事件总线实现（基于内存）

```go
// choreography/memory_bus.go
package choreography

import (
    "context"
    "fmt"
    "sync"
)

// MemoryEventBus 内存事件总线实现
type MemoryEventBus struct {
    subscribers map[string][]EventHandler
    mu          sync.RWMutex
}

// NewMemoryEventBus 创建内存事件总线
func NewMemoryEventBus() *MemoryEventBus {
    return &MemoryEventBus{
        subscribers: make(map[string][]EventHandler),
    }
}

// Publish 发布事件
func (b *MemoryEventBus) Publish(ctx context.Context, event Event) error {
    b.mu.RLock()
    handlers := b.subscribers[event.GetType()]
    b.mu.RUnlock()

    for _, handler := range handlers {
        if err := handler(ctx, event); err != nil {
            return fmt.Errorf("handler failed: %w", err)
        }
    }

    return nil
}

// Subscribe 订阅事件
func (b *MemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.subscribers[eventType] = append(b.subscribers[eventType], handler)
    return nil
}

// Unsubscribe 取消订阅
func (b *MemoryEventBus) Unsubscribe(eventType string, handler EventHandler) error {
    // 简化实现：实际应该比较函数指针
    return nil
}
```

### 2.3 Saga 存储实现（基于内存）

```go
// choreography/memory_store.go
package choreography

import (
    "context"
    "fmt"
    "sync"
)

// MemorySagaStore 内存 Saga 存储
type MemorySagaStore struct {
    sagas map[string]*SagaInstance
    mu    sync.RWMutex
}

// NewMemorySagaStore 创建内存存储
func NewMemorySagaStore() *MemorySagaStore {
    return &MemorySagaStore{
        sagas: make(map[string]*SagaInstance),
    }
}

// Save 保存 Saga
func (s *MemorySagaStore) Save(ctx context.Context, saga *SagaInstance) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.sagas[saga.ID] = saga
    return nil
}

// Load 加载 Saga
func (s *MemorySagaStore) Load(ctx context.Context, sagaID string) (*SagaInstance, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    saga, exists := s.sagas[sagaID]
    if !exists {
        return nil, fmt.Errorf("saga not found: %s", sagaID)
    }

    return saga, nil
}

// UpdateStatus 更新状态
func (s *MemorySagaStore) UpdateStatus(ctx context.Context, sagaID string, status SagaStatus) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    saga, exists := s.sagas[sagaID]
    if !exists {
        return fmt.Errorf("saga not found: %s", sagaID)
    }

    saga.Status = status
    return nil
}

// AppendStep 追加步骤结果
func (s *MemorySagaStore) AppendStep(ctx context.Context, sagaID string, result StepResult) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    saga, exists := s.sagas[sagaID]
    if !exists {
        return fmt.Errorf("saga not found: %s", sagaID)
    }

    saga.Steps = append(saga.Steps, result)
    return nil
}
```

### 2.4 订单处理示例

```go
// examples/order/order_choreography.go
package order

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "go-knowledge-base/choreography"
)

// 事件定义
const (
    EventOrderCreated      = "order.created"
    EventInventoryReserved = "inventory.reserved"
    EventPaymentProcessed  = "payment.processed"
    EventOrderConfirmed    = "order.confirmed"
    EventOrderCancelled    = "order.cancelled"
)

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
    OrderID    string    `json:"order_id"`
    CustomerID string    `json:"customer_id"`
    Items      []Item    `json:"items"`
    Total      float64   `json:"total"`
    CreatedAt  time.Time `json:"created_at"`
}

// InventoryReservedEvent 库存预留事件
type InventoryReservedEvent struct {
    OrderID     string            `json:"order_id"`
    Reservations []Reservation    `json:"reservations"`
    ReservedAt  time.Time         `json:"reserved_at"`
}

// Reservation 库存预留
type Reservation struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

// PaymentProcessedEvent 支付处理事件
type PaymentProcessedEvent struct {
    OrderID       string    `json:"order_id"`
    PaymentID     string    `json:"payment_id"`
    Amount        float64   `json:"amount"`
    Status        string    `json:"status"`
    ProcessedAt   time.Time `json:"processed_at"`
}

// Item 订单项
type Item struct {
    ProductID string  `json:"product_id"`
    Name      string  `json:"name"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// OrderService 订单服务
type OrderService struct {
    choreographer *choreography.Choreographer
    orders        map[string]*Order
}

// Order 订单
type Order struct {
    ID         string
    CustomerID string
    Items      []Item
    Total      float64
    Status     string
    CreatedAt  time.Time
}

// NewOrderService 创建订单服务
func NewOrderService(ch *choreography.Choreographer) *OrderService {
    return &OrderService{
        choreographer: ch,
        orders:        make(map[string]*Order),
    }
}

// RegisterOrderSteps 注册订单编舞步骤
func (s *OrderService) RegisterOrderSteps() error {
    // 步骤1: 订单创建（起始步骤）
    step1 := &choreography.ChoreographyStep{
        Name:             "create_order",
        ServiceName:      "order-service",
        SubscribedEvents: []string{"order.create.request"},
        PublishedEvent:   EventOrderCreated,
        Execute: func(ctx context.Context, event choreography.Event) (choreography.Event, error) {
            var req OrderCreatedEvent
            if err := json.Unmarshal(event.GetPayload(), &req); err != nil {
                return nil, err
            }

            // 创建订单
            order := &Order{
                ID:         req.OrderID,
                CustomerID: req.CustomerID,
                Items:      req.Items,
                Total:      req.Total,
                Status:     "PENDING",
                CreatedAt:  time.Now(),
            }
            s.orders[order.ID] = order

            // 创建订单创建事件
            payload, _ := json.Marshal(OrderCreatedEvent{
                OrderID:    order.ID,
                CustomerID: order.CustomerID,
                Items:      order.Items,
                Total:      order.Total,
                CreatedAt:  order.CreatedAt,
            })

            return choreography.BaseEvent{
                ID:          generateID(),
                Type:        EventOrderCreated,
                AggregateID: order.ID,
                Timestamp:   time.Now(),
                Payload:     payload,
                Metadata: map[string]string{
                    "saga_id":   event.GetMetadata("saga_id"),
                    "saga_type": "order_processing",
                },
            }, nil
        },
        Compensate: func(ctx context.Context, event choreography.Event) error {
            var data OrderCreatedEvent
            if err := json.Unmarshal(event.GetPayload(), &data); err != nil {
                return err
            }
            // 取消订单
            if order, exists := s.orders[data.OrderID]; exists {
                order.Status = "CANCELLED"
            }
            return nil
        },
    }

    if err := s.choreographer.RegisterStep(step1); err != nil {
        return err
    }

    return nil
}

// InventoryService 库存服务
type InventoryService struct {
    choreographer *choreography.Choreographer
    stock         map[string]int
    reservations  map[string][]Reservation
}

// NewInventoryService 创建库存服务
func NewInventoryService(ch *choreography.Choreographer) *InventoryService {
    return &InventoryService{
        choreographer: ch,
        stock: map[string]int{
            "PROD-001": 100,
            "PROD-002": 50,
            "PROD-003": 200,
        },
        reservations: make(map[string][]Reservation),
    }
}

// RegisterInventorySteps 注册库存步骤
func (s *InventoryService) RegisterInventorySteps() error {
    step := &choreography.ChoreographyStep{
        Name:             "reserve_inventory",
        ServiceName:      "inventory-service",
        SubscribedEvents: []string{EventOrderCreated},
        PublishedEvent:   EventInventoryReserved,
        Execute: func(ctx context.Context, event choreography.Event) (choreography.Event, error) {
            var order OrderCreatedEvent
            if err := json.Unmarshal(event.GetPayload(), &order); err != nil {
                return nil, err
            }

            // 预留库存
            var reservations []Reservation
            for _, item := range order.Items {
                available, exists := s.stock[item.ProductID]
                if !exists || available < item.Quantity {
                    return nil, fmt.Errorf("insufficient stock for %s", item.ProductID)
                }
                s.stock[item.ProductID] -= item.Quantity
                reservations = append(reservations, Reservation{
                    ProductID: item.ProductID,
                    Quantity:  item.Quantity,
                })
            }

            s.reservations[order.OrderID] = reservations

            payload, _ := json.Marshal(InventoryReservedEvent{
                OrderID:      order.OrderID,
                Reservations: reservations,
                ReservedAt:   time.Now(),
            })

            return choreography.BaseEvent{
                ID:          generateID(),
                Type:        EventInventoryReserved,
                AggregateID: order.OrderID,
                Timestamp:   time.Now(),
                Payload:     payload,
                Metadata: map[string]string{
                    "saga_id": event.GetMetadata("saga_id"),
                },
            }, nil
        },
        Compensate: func(ctx context.Context, event choreography.Event) error {
            var data InventoryReservedEvent
            if err := json.Unmarshal(event.GetPayload(), &data); err != nil {
                return err
            }

            // 释放库存
            for _, res := range data.Reservations {
                s.stock[res.ProductID] += res.Quantity
            }
            delete(s.reservations, data.OrderID)

            return nil
        },
    }

    return s.choreographer.RegisterStep(step)
}

// PaymentService 支付服务
type PaymentService struct {
    choreographer *choreography.Choreographer
    payments      map[string]*Payment
}

// Payment 支付记录
type Payment struct {
    ID        string
    OrderID   string
    Amount    float64
    Status    string
    CreatedAt time.Time
}

// NewPaymentService 创建支付服务
func NewPaymentService(ch *choreography.Choreographer) *PaymentService {
    return &PaymentService{
        choreographer: ch,
        payments:    make(map[string]*Payment),
    }
}

// RegisterPaymentSteps 注册支付步骤
func (s *PaymentService) RegisterPaymentSteps() error {
    step := &choreography.ChoreographyStep{
        Name:             "process_payment",
        ServiceName:      "payment-service",
        SubscribedEvents: []string{EventInventoryReserved},
        PublishedEvent:   EventPaymentProcessed,
        Execute: func(ctx context.Context, event choreography.Event) (choreography.Event, error) {
            var invEvent InventoryReservedEvent
            if err := json.Unmarshal(event.GetPayload(), &invEvent); err != nil {
                return nil, err
            }

            // 处理支付（简化版）
            payment := &Payment{
                ID:        generateID(),
                OrderID:   invEvent.OrderID,
                Status:    "COMPLETED",
                CreatedAt: time.Now(),
            }
            s.payments[payment.ID] = payment

            payload, _ := json.Marshal(PaymentProcessedEvent{
                OrderID:     payment.OrderID,
                PaymentID:   payment.ID,
                Status:      payment.Status,
                ProcessedAt: time.Now(),
            })

            return choreography.BaseEvent{
                ID:          generateID(),
                Type:        EventPaymentProcessed,
                AggregateID: invEvent.OrderID,
                Timestamp:   time.Now(),
                Payload:     payload,
                Metadata: map[string]string{
                    "saga_id": event.GetMetadata("saga_id"),
                },
            }, nil
        },
        Compensate: func(ctx context.Context, event choreography.Event) error {
            var data PaymentProcessedEvent
            if err := json.Unmarshal(event.GetPayload(), &data); err != nil {
                return err
            }

            // 退款处理
            if payment, exists := s.payments[data.PaymentID]; exists {
                payment.Status = "REFUNDED"
            }

            return nil
        },
    }

    return s.choreographer.RegisterStep(step)
}

func generateID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// choreography/core_test.go
package choreography

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...Field)  {}
func (m *mockLogger) Error(msg string, fields ...Field) {}
func (m *mockLogger) Debug(msg string, fields ...Field) {}

func TestChoreographer_RegisterStep(t *testing.T) {
    bus := NewMemoryEventBus()
    store := NewMemorySagaStore()
    logger := &mockLogger{}
    choreographer := NewChoreographer(bus, store, logger)

    step := &ChoreographyStep{
        Name:             "test_step",
        ServiceName:      "test-service",
        SubscribedEvents: []string{"test.event"},
        PublishedEvent:   "test.completed",
        Execute: func(ctx context.Context, event Event) (Event, error) {
            return BaseEvent{
                ID:        "test-id",
                Type:      "test.completed",
                Timestamp: time.Now(),
            }, nil
        },
    }

    err := choreographer.RegisterStep(step)
    require.NoError(t, err)

    // 验证步骤已注册
    assert.Contains(t, choreographer.steps, "test_step")
}

func TestChoreography_Execute(t *testing.T) {
    bus := NewMemoryEventBus()
    store := NewMemorySagaStore()
    logger := &mockLogger{}
    choreographer := NewChoreographer(bus, store, logger)

    var executed bool
    step := &ChoreographyStep{
        Name:             "execute_test",
        ServiceName:      "test-service",
        SubscribedEvents: []string{"trigger.event"},
        PublishedEvent:   "completed.event",
        Execute: func(ctx context.Context, event Event) (Event, error) {
            executed = true
            return BaseEvent{
                ID:          "result-id",
                Type:        "completed.event",
                Timestamp:   time.Now(),
                AggregateID: "test-aggregate",
            }, nil
        },
    }

    require.NoError(t, choreographer.RegisterStep(step))

    // 触发事件
    triggerEvent := BaseEvent{
        ID:          "trigger-id",
        Type:        "trigger.event",
        Timestamp:   time.Now(),
        AggregateID: "test-aggregate",
        Metadata:    map[string]string{"saga_id": "test-saga"},
    }

    err := bus.Publish(context.Background(), triggerEvent)
    require.NoError(t, err)

    assert.True(t, executed, "step should have been executed")
}

func TestChoreography_Compensation(t *testing.T) {
    bus := NewMemoryEventBus()
    store := NewMemorySagaStore()
    logger := &mockLogger{}
    choreographer := NewChoreographer(bus, store, logger)

    var compensated bool
    step := &ChoreographyStep{
        Name:             "failing_step",
        ServiceName:      "test-service",
        SubscribedEvents: []string{"trigger.event"},
        Execute: func(ctx context.Context, event Event) (Event, error) {
            return nil, assert.AnError
        },
        Compensate: func(ctx context.Context, event Event) error {
            compensated = true
            return nil
        },
    }

    require.NoError(t, choreographer.RegisterStep(step))

    // 先保存一个已执行的 saga
    saga := &SagaInstance{
        ID:        "test-saga",
        Type:      "test",
        Status:    SagaStatusRunning,
        StartedAt: time.Now(),
    }
    require.NoError(t, store.Save(context.Background(), saga))

    // 触发事件
    triggerEvent := BaseEvent{
        ID:          "trigger-id",
        Type:        "trigger.event",
        Timestamp:   time.Now(),
        AggregateID: "test-aggregate",
        Metadata:    map[string]string{"saga_id": "test-saga"},
    }

    // 执行应该失败并触发补偿
    _ = bus.Publish(context.Background(), triggerEvent)

    assert.True(t, compensated, "compensation should have been triggered")
}

func TestSagaInstance_Status(t *testing.T) {
    tests := []struct {
        status   SagaStatus
        expected string
    }{
        {SagaStatusPending, "PENDING"},
        {SagaStatusRunning, "RUNNING"},
        {SagaStatusCompleted, "COMPLETED"},
        {SagaStatusCompensating, "COMPENSATING"},
        {SagaStatusCompensated, "COMPENSATED"},
        {SagaStatusFailed, "FAILED"},
        {SagaStatus(999), "UNKNOWN"},
    }

    for _, tt := range tests {
        t.Run(tt.expected, func(t *testing.T) {
            assert.Equal(t, tt.expected, tt.status.String())
        })
    }
}

func TestMemoryEventBus(t *testing.T) {
    bus := NewMemoryEventBus()

    var received bool
    handler := func(ctx context.Context, event Event) error {
        received = true
        return nil
    }

    require.NoError(t, bus.Subscribe("test.event", handler))

    event := BaseEvent{
        ID:        "test-id",
        Type:      "test.event",
        Timestamp: time.Now(),
    }

    err := bus.Publish(context.Background(), event)
    require.NoError(t, err)
    assert.True(t, received)
}
```

### 3.2 集成测试

```go
// examples/order/integration_test.go
package order

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go-knowledge-base/choreography"
)

func TestOrderChoreography_Success(t *testing.T) {
    // 设置测试环境
    bus := choreography.NewMemoryEventBus()
    store := choreography.NewMemorySagaStore()
    logger := &mockLogger{}
    choreographer := choreography.NewChoreographer(bus, store, logger)

    // 创建并注册服务
    orderService := NewOrderService(choreographer)
    inventoryService := NewInventoryService(choreographer)
    paymentService := NewPaymentService(choreographer)

    require.NoError(t, orderService.RegisterOrderSteps())
    require.NoError(t, inventoryService.RegisterInventorySteps())
    require.NoError(t, paymentService.RegisterPaymentSteps())

    // 创建订单请求
    orderReq := OrderCreatedEvent{
        OrderID:    "ORDER-001",
        CustomerID: "CUST-001",
        Items: []Item{
            {ProductID: "PROD-001", Name: "Product 1", Quantity: 2, Price: 29.99},
        },
        Total: 59.98,
    }

    payload, _ := json.Marshal(orderReq)
    event := choreography.BaseEvent{
        ID:        "req-001",
        Type:      "order.create.request",
        Timestamp: time.Now(),
        Payload:   payload,
        Metadata:  map[string]string{"saga_id": "saga-001", "saga_type": "order_processing"},
    }

    // 执行编舞
    err := bus.Publish(context.Background(), event)
    require.NoError(t, err)

    // 验证订单创建
    order, exists := orderService.orders["ORDER-001"]
    require.True(t, exists)
    assert.Equal(t, "PENDING", order.Status)

    // 验证库存预留
    reservations, exists := inventoryService.reservations["ORDER-001"]
    require.True(t, exists)
    assert.Len(t, reservations, 1)

    // 验证支付创建
    var paymentFound bool
    for _, p := range paymentService.payments {
        if p.OrderID == "ORDER-001" {
            paymentFound = true
            break
        }
    }
    assert.True(t, paymentFound)
}

func TestOrderChoreography_Compensation(t *testing.T) {
    // 设置测试环境
    bus := choreography.NewMemoryEventBus()
    store := choreography.NewMemorySagaStore()
    logger := &mockLogger{}
    choreographer := choreography.NewChoreographer(bus, store, logger)

    // 创建服务（库存不足以触发补偿）
    orderService := NewOrderService(choreographer)
    inventoryService := NewInventoryService(choreographer)

    // 设置极低库存
    inventoryService.stock["PROD-001"] = 0

    require.NoError(t, orderService.RegisterOrderSteps())
    require.NoError(t, inventoryService.RegisterInventorySteps())

    initialStock := inventoryService.stock["PROD-001"]

    // 创建订单请求
    orderReq := OrderCreatedEvent{
        OrderID:    "ORDER-002",
        CustomerID: "CUST-001",
        Items: []Item{
            {ProductID: "PROD-001", Name: "Product 1", Quantity: 1, Price: 29.99},
        },
        Total: 29.99,
    }

    payload, _ := json.Marshal(orderReq)
    event := choreography.BaseEvent{
        ID:        "req-002",
        Type:      "order.create.request",
        Timestamp: time.Now(),
        Payload:   payload,
        Metadata:  map[string]string{"saga_id": "saga-002", "saga_type": "order_processing"},
    }

    // 执行编舞（应该失败并触发补偿）
    _ = bus.Publish(context.Background(), event)

    // 验证订单已被取消
    order, exists := orderService.orders["ORDER-002"]
    require.True(t, exists)
    assert.Equal(t, "CANCELLED", order.Status)

    // 验证库存未变化
    assert.Equal(t, initialStock, inventoryService.stock["PROD-001"])
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...choreography.Field)  {}
func (m *mockLogger) Error(msg string, fields ...choreography.Field) {}
func (m *mockLogger) Debug(msg string, fields ...choreography.Field) {}
```

---

## 4. 与其他模式的集成

### 4.1 与 Saga 模式的关系

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Choreography vs Orchestration                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Choreography (编舞)                      Orchestration (编排)           │
│  ┌─────────────────────────┐            ┌─────────────────────────┐    │
│  │  ┌─────┐  ┌─────┐      │            │     ┌─────────────┐     │    │
│  │  │ Svc1│──►│ Svc2│      │            │     │ Orchestrator│     │    │
│  │  └──┬──┘  └─────┘      │            │     └──────┬──────┘     │    │
│  │     │ Event             │            │            │            │    │
│  │     ▼                   │            │     ┌──────┴──────┐     │    │
│  │  ┌─────┐               │            │     ▼             ▼     │    │
│  │  │ Svc3│               │            │  ┌─────┐        ┌─────┐  │    │
│  │  └─────┘               │            │  │ Svc1│        │ Svc2│  │    │
│  │                         │            │  └──┬──┘        └──┬──┘  │    │
│  │  Characteristics:       │            │     │              │     │    │
│  │  • 去中心化               │            │     └──────────────┘     │    │
│  │  • 松耦合                │            │  Command/Response         │    │
│  │  • 事件驱动              │            │                           │    │
│  │  • 适合简单流程           │            │  Characteristics:         │    │
│  │                         │            │  • 中央协调                 │    │
│  │  Trade-offs:            │            │  • 紧耦合                   │    │
│  │  + 高可用性              │            │  • 命令驱动                 │    │
│  │  + 独立扩展              │            │  • 适合复杂流程              │    │
│  │  - 流程可见性低           │            │                           │    │
│  │  - 容易形成循环依赖        │            │  Trade-offs:              │    │
│  │                         │            │  + 流程可见性高             │    │
│  │                         │            │  + 易于调试                 │    │
│  │                         │            │  - 单点故障风险             │    │
│  │                         │            │  - 协调器成为瓶颈           │    │
│  └─────────────────────────┘            └─────────────────────────┘    │
│                                                                         │
│  Decision Matrix:                                                       │
│  ┌────────────────┬─────────────────┬────────────────────────────────┐ │
│  │ Criteria       │ Choreography    │ Orchestration                  │ │
│  ├────────────────┼─────────────────┼────────────────────────────────┤ │
│  │ 服务数量        │ < 5 services    │ > 5 services                   │ │
│  │ 流程复杂度      │ 线性流程         │ 条件分支、循环                  │ │
│  │ 失败处理        │ 简单补偿         │ 复杂重试策略                    │ │
│  │ 团队结构        │ 自治团队         │ 平台团队                        │ │
│  │ 可观察性要求     │ 可接受延迟       │ 实时追踪                        │ │
│  └────────────────┴─────────────────┴────────────────────────────────┘ │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 与 Event Sourcing 的集成

```go
// choreography/event_sourcing_integration.go
package choreography

import (
    "context"
    "encoding/json"
    "time"
)

// EventSourcedSagaStore 事件溯源 Saga 存储
type EventSourcedSagaStore struct {
    eventStore EventStore
}

// EventStore 事件存储接口
type EventStore interface {
    Append(ctx context.Context, streamID string, events []Event) error
    Read(ctx context.Context, streamID string, fromVersion int) ([]Event, error)
}

// NewEventSourcedSagaStore 创建事件溯源存储
func NewEventSourcedSagaStore(eventStore EventStore) *EventSourcedSagaStore {
    return &EventSourcedSagaStore{eventStore: eventStore}
}

// Save 保存 Saga（作为事件）
func (s *EventSourcedSagaStore) Save(ctx context.Context, saga *SagaInstance) error {
    event := SagaStartedEvent{
        SagaID:    saga.ID,
        SagaType:  saga.Type,
        StartedAt: saga.StartedAt,
    }

    payload, _ := json.Marshal(event)

    e := BaseEvent{
        ID:          generateID(),
        Type:        "saga.started",
        AggregateID: saga.ID,
        Timestamp:   time.Now(),
        Payload:     payload,
    }

    return s.eventStore.Append(ctx, saga.ID, []Event{e})
}

// SagaStartedEvent Saga 开始事件
type SagaStartedEvent struct {
    SagaID    string    `json:"saga_id"`
    SagaType  string    `json:"saga_type"`
    StartedAt time.Time `json:"started_at"`
}
```

### 4.3 与 CQRS 的集成

```go
// choreography/cqrs_integration.go
package choreography

// CQRSReadModel CQRS 读模型接口
type CQRSReadModel interface {
    Update(ctx context.Context, event Event) error
    Query(ctx context.Context, query interface{}) (interface{}, error)
}

// SagaReadModel Saga 读模型
type SagaReadModel struct {
    store SagaStore
}

// NewSagaReadModel 创建 Saga 读模型
func NewSagaReadModel(store SagaStore) *SagaReadModel {
    return &SagaReadModel{store: store}
}

// GetSagaStatus 获取 Saga 状态（优化读操作）
func (r *SagaReadModel) GetSagaStatus(ctx context.Context, sagaID string) (string, error) {
    saga, err := r.store.Load(ctx, sagaID)
    if err != nil {
        return "", err
    }
    return saga.Status.String(), nil
}

// GetActiveSagas 获取活跃 Saga 列表
func (r *SagaReadModel) GetActiveSagas(ctx context.Context) ([]*SagaInstance, error) {
    // 实际实现会查询专门的读模型数据库（如 Redis、Elasticsearch）
    // 这里简化为直接查询
    return nil, nil
}
```

---

## 5. 决策标准

### 5.1 何时使用编舞模式

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Choreography Decision Tree                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  开始                                                                    │
│   │                                                                     │
│   ▼                                                                     │
│  ┌─────────────────────┐                                               │
│  │ 需要协调多个服务？    │───否──► 不需要分布式事务                        │
│  └──────────┬──────────┘                                               │
│             │是                                                         │
│             ▼                                                           │
│  ┌─────────────────────┐                                               │
│  │ 需要强一致性？       │───是──► 考虑 2PC / TCC                         │
│  └──────────┬──────────┘                                               │
│             │否                                                         │
│             ▼                                                           │
│  ┌─────────────────────┐                                               │
│  │ 最终一致性可接受？   │───否──► 重新评估架构需求                         │
│  └──────────┬──────────┘                                               │
│             │是                                                         │
│             ▼                                                           │
│  ┌─────────────────────┐                                               │
│  │ 流程复杂度如何？     │                                               │
│  └──────────┬──────────┘                                               │
│             │                                                           │
│      ┌──────┴──────┐                                                    │
│      ▼             ▼                                                    │
│    简单/线性      复杂/条件分支                                           │
│      │             │                                                    │
│      ▼             ▼                                                    │
│  ┌─────────┐   ┌─────────┐                                              │
│  │Choreography│  │Orchestration│                                          │
│  │ 编舞模式   │   │ 编排模式    │                                          │
│  └─────────┘   └─────────┘                                              │
│      │             │                                                    │
│      ▼             ▼                                                    │
│  事件驱动         中央协调                                               │
│  松耦合          紧耦合                                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 选型矩阵

| 决策因素 | 使用 Choreography | 使用 Orchestration |
|---------|------------------|-------------------|
| **服务数量** | < 5 个服务 | > 5 个服务 |
| **流程复杂度** | 线性流程，少量步骤 | 复杂分支，多步骤 |
| **团队结构** | 自治团队，领域边界清晰 | 集中式开发团队 |
| **可观察性需求** | 可接受异步延迟 | 需要实时监控 |
| **事务边界** | 清晰的服务边界 | 模糊的边界 |
| **失败处理** | 简单补偿逻辑 | 复杂重试策略 |
| **演化需求** | 独立演化 | 协调演化 |

### 5.3 反模式警告

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Choreography Anti-Patterns                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ⚠️ 循环依赖                                                            │
│     Service A ──► Service B ──► Service C ──► Service A                │
│     解决方案: 引入 Saga 协调器或重新划分领域                              │
│                                                                         │
│  ⚠️ 事件爆炸                                                            │
│     每个操作产生过多细粒度事件                                           │
│     解决方案: 批量事件、领域事件聚合                                      │
│                                                                         │
│  ⚠️ 分布式单体                                                          │
│     表面上微服务，实际上紧耦合                                            │
│     解决方案: 明确服务边界，避免共享数据库                                 │
│                                                                         │
│  ⚠️ 上帝事件                                                            │
│     超大事件包含所有数据                                                  │
│     解决方案: 事件只包含 ID，通过 API 获取详情                            │
│                                                                         │
│  ⚠️ 缺少补偿                                                            │
│     只实现成功路径，忽略失败处理                                          │
│     解决方案: 每个步骤都要有补偿策略                                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 6. 生产环境考虑

### 6.1 可观察性

```go
// choreography/observability.go
package choreography

import (
    "context"
    "time"
)

// SagaMetrics Saga 指标
type SagaMetrics struct {
    Started     Counter
    Completed   Counter
    Failed      Counter
    Compensated Counter
    Duration    Histogram
}

// Counter 计数器接口
type Counter interface {
    Inc()
    Add(delta float64)
}

// Histogram 直方图接口
type Histogram interface {
    Observe(value float64)
}

// InstrumentedChoreographer 带指标的装饰器
type InstrumentedChoreographer struct {
    choreographer *Choreographer
    metrics       *SagaMetrics
}

// NewInstrumentedChoreographer 创建带指标的编舞器
func NewInstrumentedChoreographer(ch *Choreographer, metrics *SagaMetrics) *InstrumentedChoreographer {
    return &InstrumentedChoreographer{
        choreographer: ch,
        metrics:       metrics,
    }
}

// RegisterStep 注册带指标的步骤
func (ic *InstrumentedChoreographer) RegisterStep(step *ChoreographyStep) error {
    instrumentedStep := &ChoreographyStep{
        Name:             step.Name,
        ServiceName:      step.ServiceName,
        SubscribedEvents: step.SubscribedEvents,
        PublishedEvent:   step.PublishedEvent,
        Execute: func(ctx context.Context, event Event) (Event, error) {
            start := time.Now()
            result, err := step.Execute(ctx, event)
            ic.metrics.Duration.Observe(time.Since(start).Seconds())

            if err != nil {
                ic.metrics.Failed.Inc()
            } else {
                ic.metrics.Completed.Inc()
            }

            return result, err
        },
        Compensate: step.Compensate,
    }

    return ic.choreographer.RegisterStep(instrumentedStep)
}
```

### 6.2 性能优化

```go
// choreography/performance.go
package choreography

import (
    "context"
    "sync"
)

// BatchEventBus 批量事件总线
type BatchEventBus struct {
    inner     EventBus
    buffer    []Event
    mu        sync.Mutex
    batchSize int
    flushChan chan struct{}
}

// NewBatchEventBus 创建批量事件总线
func NewBatchEventBus(inner EventBus, batchSize int) *BatchEventBus {
    bus := &BatchEventBus{
        inner:     inner,
        batchSize: batchSize,
        flushChan: make(chan struct{}),
    }
    go bus.flushLoop()
    return bus
}

// Publish 批量发布
func (b *BatchEventBus) Publish(ctx context.Context, event Event) error {
    b.mu.Lock()
    b.buffer = append(b.buffer, event)
    shouldFlush := len(b.buffer) >= b.batchSize
    b.mu.Unlock()

    if shouldFlush {
        return b.Flush(ctx)
    }

    return nil
}

// Flush 强制刷新
func (b *BatchEventBus) Flush(ctx context.Context) error {
    b.mu.Lock()
    events := b.buffer
    b.buffer = nil
    b.mu.Unlock()

    for _, event := range events {
        if err := b.inner.Publish(ctx, event); err != nil {
            return err
        }
    }

    return nil
}

func (b *BatchEventBus) flushLoop() {
    // 定期刷新逻辑
}

func (b *BatchEventBus) Subscribe(eventType string, handler EventHandler) error {
    return b.inner.Subscribe(eventType, handler)
}

func (b *BatchEventBus) Unsubscribe(eventType string, handler EventHandler) error {
    return b.inner.Unsubscribe(eventType, handler)
}
```

---

## 7. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Choreography Implementation Checklist                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 定义清晰的领域事件契约                                                │
│  □ 绘制事件流图，识别循环依赖                                             │
│  □ 为每个步骤定义补偿操作                                                │
│  □ 确定 Saga 边界和超时策略                                              │
│  □ 设计事件版本策略（向后兼容）                                           │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现幂等的事件处理器                                                  │
│  □ 配置死信队列处理失败事件                                              │
│  □ 实现 Saga 状态持久化                                                  │
│  □ 添加分布式追踪（Correlation ID）                                      │
│  □ 实现事件模式验证                                                      │
│                                                                         │
│  测试阶段:                                                               │
│  □ 单元测试每个步骤                                                      │
│  □ 集成测试完整流程                                                      │
│  □ 混沌测试（模拟服务故障）                                              │
│  □ 负载测试事件吞吐量                                                    │
│  □ 补偿路径测试                                                          │
│                                                                         │
│  运维阶段:                                                               │
│  □ 配置 Saga 监控仪表板                                                  │
│  □ 设置告警（长时间运行 Saga）                                           │
│  □ 实现人工干预接口                                                      │
│  □ 定期审查事件序列图                                                    │
│  □ 维护事件契约文档                                                      │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 不要在事件中包含敏感信息                                              │
│  ❌ 不要依赖事件顺序（使用因果向量）                                      │
│  ❌ 不要忽略补偿失败的情况                                                │
│  ❌ 不要在事件处理器中执行长时间操作                                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>20KB, 完整形式化 + Go 实现 + 测试 + 决策标准)

**相关文档**:

- [EC-032-Orchestration-Pattern.md](./EC-032-Orchestration-Pattern.md)
- [EC-033-Transactional-Outbox.md](./EC-033-Transactional-Outbox.md)
- [EC-039-Domain-Event-Pattern.md](./EC-039-Domain-Event-Pattern.md)
- [EC-090-Task-Compensation-Saga-Pattern.md](./EC-090-Task-Compensation-Saga-Pattern.md)
