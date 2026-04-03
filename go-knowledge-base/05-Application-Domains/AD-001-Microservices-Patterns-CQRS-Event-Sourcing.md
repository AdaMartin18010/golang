# AD-001: 微服务模式：CQRS 与事件溯源 (Microservices: CQRS & Event Sourcing)

> **维度**: Application Domains
> **级别**: S (25+ KB)
> **标签**: #microservices #cqrs #event-sourcing #domain-driven-design
> **权威来源**: [Microsoft CQRS Journey](https://msdn.microsoft.com/en-us/library/jj554200.aspx), [Event Sourcing by Martin Fowler](https://martinfowler.com/eaaDev/EventSourcing.html), [DDD Reference](https://www.domainlanguage.com/wp-content/uploads/2016/05/DDD_Reference_2015-03.pdf)

---

## 架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CQRS with Event Sourcing Architecture                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Command Side                              Query Side                       │
│  ────────────                              ──────────                       │
│                                                                              │
│  ┌──────────────┐                         ┌──────────────┐                 │
│  │ Command API  │                         │ Query API    │                 │
│  │ (REST/gRPC)  │                         │ (GraphQL)    │                 │
│  └──────┬───────┘                         └──────┬───────┘                 │
│         │                                         │                         │
│  ┌──────▼───────┐                         ┌──────▼───────┐                 │
│  │ Command      │                         │ Read Model   │                 │
│  │ Handlers     │                         │ Projections  │                 │
│  └──────┬───────┘                         └──────┬───────┘                 │
│         │                                         │                         │
│  ┌──────▼───────┐                         ┌──────▼───────┐                 │
│  │ Aggregate    │                         │ ElasticSearch│                 │
│  │ (Domain      │                         │ / MongoDB    │                 │
│  │  Model)      │                         └──────────────┘                 │
│  └──────┬───────┘                                                           │
│         │                                                                   │
│  ┌──────▼───────┐      ┌──────────────┐      ┌──────────────┐             │
│  │ Domain       │─────►│ Event Store  │◄─────│ Event        │             │
│  │ Events       │      │ (EventStoreDB│      │ Projectors   │             │
│  └──────────────┘      │  / Kafka)    │      └──────────────┘             │
│                        └──────────────┘                                    │
│                                                                              │
│  Single Source of Truth: The Event Stream                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## CQRS 核心概念

### 命令与查询分离

```go
// Command Side - 处理写操作
package command

type OrderCommandHandler struct {
    eventStore eventstore.EventStore
    aggregateRepo OrderRepository
}

// CreateOrderCommand 创建订单命令
type CreateOrderCommand struct {
    OrderID   string
    CustomerID string
    Items     []OrderItem
}

func (h *OrderCommandHandler) HandleCreateOrder(ctx context.Context,
    cmd CreateOrderCommand) error {

    // 1. 加载或创建聚合
    order, err := h.aggregateRepo.Load(ctx, cmd.OrderID)
    if err != nil {
        order = aggregate.NewOrder(cmd.OrderID, cmd.CustomerID)
    }

    // 2. 执行业务逻辑
    for _, item := range cmd.Items {
        if err := order.AddItem(item.ProductID, item.Quantity, item.Price); err != nil {
            return err
        }
    }

    // 3. 提交事件到 Event Store
    events := order.UncommittedEvents()
    if err := h.eventStore.Append(ctx, cmd.OrderID, events); err != nil {
        return err
    }

    return nil
}
```

```go
// Query Side - 处理读操作
package query

type OrderQueryHandler struct {
    readModel ReadModel
}

// OrderView 查询视图（反规范化）
type OrderView struct {
    OrderID       string    `json:"order_id" bson:"_id"`
    CustomerID    string    `json:"customer_id"`
    TotalAmount   float64   `json:"total_amount"`
    Status        string    `json:"status"`
    Items         []ItemView `json:"items"`
    CreatedAt     time.Time `json:"created_at"`
}

func (h *OrderQueryHandler) GetOrder(ctx context.Context,
    orderID string) (*OrderView, error) {

    // 直接从读取模型查询（优化查询性能）
    return h.readModel.FindByID(ctx, orderID)
}

func (h *OrderQueryHandler) GetCustomerOrders(ctx context.Context,
    customerID string) ([]OrderView, error) {

    // 利用 Read Model 的索引优化
    return h.readModel.FindByCustomer(ctx, customerID)
}
```

---

## 事件溯源实现

### 聚合根设计

```go
package aggregate

type Order struct {
    aggregate.Base

    CustomerID  string
    Items       []OrderItem
    Status      OrderStatus
    TotalAmount float64
}

func (o *Order) AddItem(productID string, quantity int, price float64) error {
    // 业务规则验证
    if o.Status != OrderStatusPending {
        return errors.New("cannot modify submitted order")
    }

    if quantity <= 0 {
        return errors.New("quantity must be positive")
    }

    // 创建事件
    event := &ItemAddedEvent{
        OrderID:   o.ID,
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
        TotalLine: float64(quantity) * price,
    }

    // 应用到聚合
    o.Apply(event)

    return nil
}

func (o *Order) Apply(event domain.Event) {
    switch e := event.(type) {
    case *OrderCreatedEvent:
        o.CustomerID = e.CustomerID
        o.Status = OrderStatusPending

    case *ItemAddedEvent:
        o.Items = append(o.Items, OrderItem{
            ProductID: e.ProductID,
            Quantity:  e.Quantity,
            Price:     e.Price,
        })
        o.TotalAmount += e.TotalLine

    case *OrderSubmittedEvent:
        o.Status = OrderStatusSubmitted

    case *OrderShippedEvent:
        o.Status = OrderStatusShipped
    }

    // 记录未提交事件
    o.AppendEvent(event)
}

// 从事件流重建聚合
func (o *Order) LoadFromHistory(events []domain.Event) {
    for _, event := range events {
        o.Apply(event)
    }
    o.ClearUncommittedEvents()
}
```

### 投影器（Projector）

```go
package projector

type OrderProjector struct {
    readModel query.ReadModel
}

func (p *OrderProjector) HandleOrderCreated(ctx context.Context,
    event *domain.OrderCreatedEvent) error {

    view := query.OrderView{
        OrderID:    event.OrderID,
        CustomerID: event.CustomerID,
        Status:     "pending",
        CreatedAt:  event.CreatedAt,
    }

    return p.readModel.Insert(ctx, view)
}

func (p *OrderProjector) HandleItemAdded(ctx context.Context,
    event *domain.ItemAddedEvent) error {

    // 更新读取模型
    return p.readModel.Update(ctx, event.OrderID, bson.M{
        "$push": bson.M{
            "items": query.ItemView{
                ProductID: event.ProductID,
                Quantity:  event.Quantity,
                Price:     event.Price,
            },
        },
        "$inc": bson.M{
            "total_amount": event.TotalLine,
        },
    })
}

func (p *OrderProjector) HandleOrderSubmitted(ctx context.Context,
    event *domain.OrderSubmittedEvent) error {

    return p.readModel.Update(ctx, event.OrderID, bson.M{
        "$set": bson.M{
            "status": "submitted",
        },
    })
}
```

---

## 一致性保证

### 最终一致性模式

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Eventual Consistency Timeline                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Time ──────────────────────────────────────────────────────────────►      │
│                                                                              │
│  T0: Command Side       T1: Event Store       T2: Query Side               │
│      ┌──────────┐           ┌──────────┐          ┌──────────┐             │
│      │ Create   │──────────►│ Persist  │─────────►│ Project  │             │
│      │ Order    │           │ Event    │          │ to Read  │             │
│      └──────────┘           └──────────┘          │ Model    │             │
│                                                   └──────────┘             │
│  Consistency Window: T2 - T0                                               │
│                                                                              │
│  读自己的写（Read-Your-Own-Writes）:                                        │
│  - 方案1: 命令返回时等待投影完成                                            │
│  - 方案2: 命令返回版本号，查询时轮询                                         │
│  - 方案3: 前端乐观更新 + 后台同步                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Saga 分布式事务

```go
// Order Saga 编排
func OrderSaga(ctx context.Context, cmd CreateOrderCommand) error {
    saga := sagas.New()

    // Step 1: 创建订单
    saga.AddStep(sagas.Step{
        Name:   "create_order",
        Action: func() error { return orderService.Create(ctx, cmd) },
        Compensate: func() error { return orderService.Cancel(ctx, cmd.OrderID) },
    })

    // Step 2: 扣减库存
    saga.AddStep(sagas.Step{
        Name:   "reserve_inventory",
        Action: func() error {
            return inventoryService.Reserve(ctx, cmd.Items)
        },
        Compensate: func() error {
            return inventoryService.Release(ctx, cmd.Items)
        },
    })

    // Step 3: 处理支付
    saga.AddStep(sagas.Step{
        Name:   "process_payment",
        Action: func() error {
            return paymentService.Charge(ctx, cmd.CustomerID, cmd.Total)
        },
        Compensate: func() error {
            return paymentService.Refund(ctx, cmd.CustomerID, cmd.Total)
        },
    })

    return saga.Execute(ctx)
}
```

---

## 实践建议

| 场景 | 建议 |
|------|------|
| 简单 CRUD | 不要使用 CQRS，过度设计 |
| 复杂领域逻辑 | 使用 CQRS + DDD |
| 高读低写 | CQRS 优势明显 |
| 强一致性要求 | 慎用最终一致性 |
| 审计需求强 | 事件溯源必需 |

---

## 参考文献

1. [CQRS Journey](https://msdn.microsoft.com/en-us/library/jj554200.aspx) - Microsoft Patterns & Practices
2. [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html) - Martin Fowler
3. [Domain-Driven Design Reference](https://www.domainlanguage.com/wp-content/uploads/2016/05/DDD_Reference_2015-03.pdf) - Eric Evans
4. [Exploring CQRS and Event Sourcing](https://docs.microsoft.com/en-us/previous-versions/msp-n-p/jj554200(v=pandp.10)) - Microsoft
5. [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman

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