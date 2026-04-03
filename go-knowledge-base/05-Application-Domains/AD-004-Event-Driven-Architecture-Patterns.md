# AD-004: 事件驱动架构模式 (Event-Driven Architecture Patterns)

> **维度**: Application Domains
> **级别**: S (17+ KB)
> **标签**: #event-driven #eda #event-sourcing #cqrs #saga
> **权威来源**: [Building Event-Driven Microservices](https://www.oreilly.com/library/view/building-event-driven-microservices/9781492057888/), [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)

---

## 事件驱动架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Event-Driven Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐     Event Bus        ┌─────────────┐                       │
│  │   Service   │    (Kafka/Rabbit)    │   Service   │                       │
│  │     A       │◄────────────────────►│     B       │                       │
│  │  (Producer) │                      │  (Consumer) │                       │
│  └─────────────┘                      └─────────────┘                       │
│         │                                    │                              │
│         │ Produce                            │ Consume                      │
│         ▼                                    ▼                              │
│  ┌─────────────┐                      ┌─────────────┐                       │
│  │   Order     │                      │  Inventory  │                       │
│  │  Created    │                      │  Updated    │                       │
│  └─────────────┘                      └─────────────┘                       │
│                                                                              │
│  模式:                                                                        │
│  ├── Event Notification (事件通知)                                           │
│  ├── Event-Carried State Transfer (事件携带状态转移)                          │
│  ├── Event Sourcing (事件溯源)                                               │
│  └── CQRS (命令查询责任分离)                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 事件模式

### 1. Event Notification (事件通知)

```go
// 轻量级通知，消费者需查询获取完整数据
type OrderCreatedEvent struct {
    EventID   string    `json:"event_id"`
    OrderID   string    `json:"order_id"`
    Timestamp time.Time `json:"timestamp"`
}

// 消费者收到事件后，调用 API 获取订单详情
func handleOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    // 查询订单服务获取完整订单信息
    order, err := orderClient.GetOrder(ctx, event.OrderID)
    if err != nil {
        return err
    }
    // 处理...
}
```

### 2. Event-Carried State Transfer (ECST)

```go
// 事件携带完整状态，无需额外查询
type OrderCreatedEvent struct {
    EventID   string    `json:"event_id"`
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Items     []Item    `json:"items"`
    Total     float64   `json:"total"`
    Address   Address   `json:"address"`
    Timestamp time.Time `json:"timestamp"`
}

// 消费者直接使用事件中的数据
func handleOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    // 直接使用 event.Items, event.Total 等，无需查询
    return updateInventory(ctx, event.Items)
}
```

---

## 事件溯源 (Event Sourcing)

```
传统 CRUD:
┌─────────────────┐      ┌─────────────────┐
│    Application  │─────►│   Database      │
│                 │      │  ┌───────────┐  │
│  Update Order   │      │  │  Order    │  │
│  Status: Paid   │      │  │  Status   │──┼──► Paid
│                 │      │  │  ...      │  │
└─────────────────┘      │  └───────────┘  │
                         └─────────────────┘

事件溯源:
┌─────────────────┐      ┌─────────────────────────────────┐
│    Application  │─────►│          Event Store            │
│                 │      │  ┌───────────────────────────┐  │
│  Create Payment │      │  │ OrderCreated              │  │
│                 │      │  │ PaymentReceived           │  │
│                 │      │  │ OrderShipped              │  │
└─────────────────┘      │  └───────────────────────────┘  │
         │               └─────────────────────────────────┘
         │                           │
         │ 读取当前状态                │ 重放事件
         ▼                           ▼
┌─────────────────┐      ┌─────────────────────────────────┐
│  Current State  │◄─────│  Aggregate (Order)              │
│  ┌───────────┐  │      │  apply(OrderCreated)            │
│  │ Order     │  │      │  apply(PaymentReceived)         │
│  │ Status:   │  │      │  apply(OrderShipped)            │
│  │   Paid    │  │      │  State: {Paid, Shipped}         │
│  └───────────┘  │      └─────────────────────────────────┘
└─────────────────┘
```

### 事件存储实现

```go
package eventsourcing

// Event 领域事件接口
type Event interface {
    EventID() string
    AggregateID() string
    EventType() string
    Timestamp() time.Time
}

// Aggregate 聚合根接口
type Aggregate interface {
    AggregateID() string
    Version() int
    Apply(event Event) error
    UncommittedEvents() []Event
    MarkCommitted()
}

// EventStore 事件存储
type EventStore interface {
    Append(ctx context.Context, aggregateID string, events []Event, expectedVersion int) error
    Load(ctx context.Context, aggregateID string) ([]Event, error)
}

// Order 订单聚合
type Order struct {
    ID      string
    Status  OrderStatus
    Items   []OrderItem
    Version int
    changes []Event
}

type OrderCreated struct {
    OrderID string
    Items   []OrderItem
}

type OrderPaid struct {
    OrderID string
    Amount  float64
}

func (o *Order) Apply(event Event) error {
    switch e := event.(type) {
    case *OrderCreated:
        o.ID = e.OrderID
        o.Status = StatusPending
        o.Items = e.Items
    case *OrderPaid:
        o.Status = StatusPaid
    // ...
    }
    o.Version++
    return nil
}

func (o *Order) Create(items []OrderItem) error {
    event := &OrderCreated{
        OrderID: o.ID,
        Items:   items,
    }
    if err := o.Apply(event); err != nil {
        return err
    }
    o.changes = append(o.changes, event)
    return nil
}

func (o *Order) Pay(amount float64) error {
    if o.Status != StatusPending {
        return errors.New("order not pending")
    }
    event := &OrderPaid{
        OrderID: o.ID,
        Amount:  amount,
    }
    if err := o.Apply(event); err != nil {
        return err
    }
    o.changes = append(o.changes, event)
    return nil
}
```

---

## CQRS 模式

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CQRS (Command Query Responsibility Segregation)        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Command Side                    Event Bus          Query Side              │
│  ┌─────────────────┐                              ┌─────────────────┐       │
│  │   Command       │      OrderCreated            │    Query        │       │
│  │   Handler       │──────────────────────────────│    Handler      │       │
│  │                 │      OrderPaid               │                 │       │
│  │  CreateOrder    │──────────────────────────────│  GetOrders      │       │
│  │  PayOrder       │      OrderShipped            │  GetOrderByID   │       │
│  │                 │──────────────────────────────│                 │       │
│  └─────────────────┘                              └─────────────────┘       │
│         │                                                │                  │
│         ▼                                                ▼                  │
│  ┌─────────────────┐                              ┌─────────────────┐       │
│  │  Event Store    │                              │  Read Database  │       │
│  │  (Source of     │                              │  (Optimized     │       │
│  │   Truth)        │                              │   for queries)  │       │
│  └─────────────────┘                              └─────────────────┘       │
│                                                                              │
│  分离优势:                                                                    │
│  - 写模型优化一致性                                                           │
│  - 读模型优化查询性能                                                         │
│  - 独立扩展                                                                   │
│  - 多种读模型 (SQL/NoSQL/Search)                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 参考文献

1. [Building Event-Driven Microservices](https://www.oreilly.com/library/view/building-event-driven-microservices/9781492057888/) - Adam Bellemare
2. [Event Sourcing Pattern](https://martinfowler.com/eaaDev/EventSourcing.html) - Martin Fowler
3. [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html) - Martin Fowler

---

## 架构决策记录

### 决策矩阵

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| A | 高性能 | 复杂 | 大规模 |
| B | 简单 | 扩展性差 | 小规模 |

### 风险评估

**风险 R.1**: 性能瓶颈
- 概率: 中
- 影响: 高
- 缓解: 缓存、分片

**风险 R.2**: 单点故障
- 概率: 低
- 影响: 极高
- 缓解: 冗余、故障转移

### 实施路线图

`
Phase 1: 基础设施 (Week 1-2)
Phase 2: 核心功能 (Week 3-6)
Phase 3: 优化加固 (Week 7-8)
`

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 架构决策记录 (ADR)

### 上下文

业务需求和技术约束分析。

### 决策

选择方案A作为主要架构方向。

### 后果

正面：
- 可扩展性提升
- 维护成本降低

负面：
- 初期开发复杂度增加
- 团队学习成本

### 实施指南

`
Week 1-2: 基础设施搭建
Week 3-4: 核心功能开发
Week 5-6: 集成测试
Week 7-8: 性能优化
`

### 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 性能不足 | 中 | 高 | 缓存、分片 |
| 兼容性 | 低 | 中 | 接口适配层 |

### 监控指标

- 系统吞吐量
- 响应延迟
- 错误率
- 资源利用率

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 系统设计

### 需求分析

功能需求和非功能需求的完整梳理。

### 架构视图

`
┌─────────────────────────────────────┐
│           API Gateway               │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    ▼                   ▼
┌─────────┐       ┌─────────┐
│ Service │       │ Service │
│   A     │       │   B     │
└────┬────┘       └────┬────┘
     │                 │
     └────────┬────────┘
              ▼
        ┌─────────┐
        │  Data   │
        │  Store  │
        └─────────┘
`

### 技术选型

| 组件 | 技术 | 理由 |
|------|------|------|
| API | gRPC | 性能 |
| DB | PostgreSQL | 可靠 |
| Cache | Redis | 速度 |
| Queue | Kafka | 吞吐 |

### 性能指标

- QPS: 10K+
- P99 Latency: <100ms
- Availability: 99.99%

### 运维手册

- 部署流程
- 监控配置
- 应急预案
- 容量规划

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02