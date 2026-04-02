# 事件驱动架构 (Event-Driven Architecture)

> **分类**: 成熟应用领域

---

## 事件总线

```go
type EventBus struct {
    subscribers map[string][]chan Event
    mu          sync.RWMutex
}

type Event struct {
    Type    string
    Payload interface{}
    Time    time.Time
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]chan Event),
    }
}

func (eb *EventBus) Subscribe(eventType string, ch chan Event) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
}

func (eb *EventBus) Publish(event Event) {
    eb.mu.RLock()
    defer eb.mu.RUnlock()

    for _, ch := range eb.subscribers[event.Type] {
        go func(c chan Event) {
            c <- event
        }(ch)
    }
}
```

---

## CQRS 模式

```go
// 命令端
type CommandHandler struct {
    eventStore EventStore
}

func (h *CommandHandler) CreateOrder(cmd CreateOrderCommand) error {
    order := NewOrder(cmd.ID, cmd.Items)

    // 保存事件
    events := order.UncommittedEvents()
    if err := h.eventStore.Save(events); err != nil {
        return err
    }

    return nil
}

// 查询端
type QueryHandler struct {
    readModel ReadModel
}

func (h *QueryHandler) GetOrder(query GetOrderQuery) (*OrderView, error) {
    return h.readModel.Find(query.ID)
}
```

---

## 事件溯源

```go
type EventStore interface {
    Save(events []Event) error
    GetStream(aggregateID string) ([]Event, error)
}

type Aggregate interface {
    Apply(event Event)
    UncommittedEvents() []Event
}

// 重放事件
func ReplayEvents(store EventStore, aggregateID string, aggregate Aggregate) error {
    events, err := store.GetStream(aggregateID)
    if err != nil {
        return err
    }

    for _, event := range events {
        aggregate.Apply(event)
    }

    return nil
}
```

---

## Saga 编排

```go
type Orchestrator struct {
    saga Saga
    bus  EventBus
}

func (o *Orchestrator) Start(order Order) {
    // 发送第一个命令
    o.bus.Publish(Event{
        Type:    "RESERVE_INVENTORY",
        Payload: order,
    })
}

func (o *Orchestrator) HandleEvent(event Event) {
    switch event.Type {
    case "INVENTORY_RESERVED":
        o.bus.Publish(Event{
            Type:    "PROCESS_PAYMENT",
            Payload: event.Payload,
        })
    case "PAYMENT_FAILED":
        o.bus.Publish(Event{
            Type:    "RELEASE_INVENTORY",
            Payload: event.Payload,
        })
    }
}
```
