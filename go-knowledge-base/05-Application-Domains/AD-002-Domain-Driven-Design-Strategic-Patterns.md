# AD-002: 领域驱动设计战略模式 (Domain-Driven Design Strategic Patterns)

> **维度**: Application Domains
> **级别**: S (25+ KB)
> **标签**: #ddd #domain-driven-design #bounded-context #strategic-design
> **权威来源**: [Domain-Driven Design](https://www.domainlanguage.com/ddd/) - Eric Evans, [Implementing DDD](https://www.amazon.com/Implementing-Domain-Driven-Design-Vaughn-Vernon/dp/0321834577) - Vaughn Vernon

---

## DDD 核心概念

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Domain-Driven Design Overview                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Problem Space                    Solution Space                             │
│  ─────────────                    ─────────────                              │
│                                                                              │
│  ┌─────────────┐                  ┌─────────────┐                           │
│  │   Domain    │                  │  Bounded    │                           │
│  │  (业务领域)  │─────────────────►│  Context    │                           │
│  └─────────────┘                  │  (限界上下文)│                           │
│                                   └──────┬──────┘                           │
│                                          │                                   │
│                                   ┌──────┴──────┐                           │
│                                   │  Subdomain  │                           │
│                                   │  (子域)     │                           │
│                                   └─────────────┘                           │
│                                                                              │
│  Core Domain: 核心竞争力，最复杂，投入最多资源                                  │
│  Supporting Subdomain: 支持核心，可能外包或使用现成方案                          │
│  Generic Subdomain: 通用功能，使用现成方案                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 限界上下文 (Bounded Context)

### 定义

**限界上下文是语义一致性的边界。在同一个限界上下文内，领域模型是一致的；跨上下文则需要显式映射。**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Bounded Contexts Example                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  电商系统:                                                                    │
│                                                                              │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐         │
│  │  Sales Context  │    │ Inventory Ctx   │    │  Shipping Ctx   │         │
│  │                 │    │                 │    │                 │         │
│  │  Product        │    │  Product        │    │  Product        │         │
│  │  - price        │    │  - quantity     │    │  - weight       │         │
│  │  - promotion    │    │  - location     │    │  - dimensions   │         │
│  │                 │    │                 │    │                 │         │
│  │  Order          │    │  Stock          │    │  Shipment       │         │
│  │  - totalAmount  │    │  - available    │    │  - trackingNo   │         │
│  │  - discount     │    │  - reserved     │    │  - carrier      │         │
│  └────────┬────────┘    └────────┬────────┘    └────────┬────────┘         │
│           │                      │                      │                  │
│           └──────────────────────┼──────────────────────┘                  │
│                                  │                                          │
│                           ┌──────▼──────┐                                   │
│                           │  Shared     │                                   │
│                           │  Kernel     │                                   │
│                           │  (共享内核)  │                                   │
│                           └─────────────┘                                   │
│                                                                              │
│  注意：同一个概念在不同上下文中有不同含义：                                    │
│  • Sales.Product.price: 售价（含促销）                                        │
│  • Inventory.Product.quantity: 库存数量                                       │
│  • Shipping.Product.weight: 重量（用于计算运费）                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 上下文映射 (Context Mapping)

| 关系类型 | 说明 | 使用场景 |
|---------|------|---------|
| **Partnership** | 双向依赖，共同演进 | 紧密合作的团队 |
| **Shared Kernel** | 共享子域，同时修改 | 核心概念共享 |
| **Customer-Supplier** | 上游优先，下游适配 | 有依赖关系 |
| **Conformist** | 下游完全接受上游模型 | 无法影响上游 |
| **Anticorruption Layer** | 防腐层隔离 | 遗留系统集成 |
| **Open Host Service** | 发布语言供外部使用 | 开放平台 |
| **Published Language** | 明确发布的共享语言 | 跨组织 |
| **Separate Ways** | 完全独立 | 无共享需求 |

---

## 战术设计模式

### 实体 (Entity)

```go
// 有唯一标识，生命周期中状态可变
type Order struct {
    ID        OrderID       // 唯一标识
    CustomerID CustomerID
    Items     []OrderItem   // 值对象
    Status    OrderStatus   // 状态
    Total     Money         // 值对象
    CreatedAt time.Time
}

func (o *Order) Confirm() error {
    if o.Status != Pending {
        return errors.New("order already confirmed")
    }
    o.Status = Confirmed
    o.AddDomainEvent(OrderConfirmed{OrderID: o.ID})
    return nil
}
```

### 值对象 (Value Object)

```go
// 无标识，不可变，通过属性相等性比较
type Money struct {
    Amount   decimal.Decimal
    Currency string
}

func (m Money) Add(other Money) (Money, error) {
    if m.Currency != other.Currency {
        return Money{}, errors.New("currency mismatch")
    }
    return Money{
        Amount:   m.Amount.Add(other.Amount),
        Currency: m.Currency,
    }, nil
}

// 相等性比较
func (m Money) Equals(other Money) bool {
    return m.Amount.Equal(other.Amount) && m.Currency == other.Currency
}

// 其他值对象
type Address struct {
    Street  string
    City    string
    ZipCode string
    Country string
}
```

### 聚合 (Aggregate)

```go
// 一致性边界，外部只能通过根实体引用
type OrderAggregate struct {
    Order  *Order        // 聚合根
    Items  []OrderItem   // 内部实体
    Payment *Payment     // 内部实体
}

// 聚合根提供所有操作
func (oa *OrderAggregate) AddItem(product Product, quantity int) error {
    if oa.Order.Status != Pending {
        return errors.New("cannot modify confirmed order")
    }

    item := OrderItem{
        ProductID: product.ID,
        Quantity:  quantity,
        Price:     product.Price,
    }

    oa.Items = append(oa.Items, item)
    oa.recalculateTotal()

    return nil
}

func (oa *OrderAggregate) Confirm() error {
    // 业务规则验证
    if len(oa.Items) == 0 {
        return errors.New("order must have at least one item")
    }

    if oa.Payment == nil {
        return errors.New("payment required")
    }

    return oa.Order.Confirm()
}
```

### 领域事件 (Domain Event)

```go
// 记录领域中的重要发生的事情
type DomainEvent interface {
    EventID() string
    OccurredAt() time.Time
}

type OrderConfirmed struct {
    EventID_   string
    OccurredAt_ time.Time
    OrderID    OrderID
    CustomerID CustomerID
    Total      Money
}

func (e OrderConfirmed) EventID() string     { return e.EventID_ }
func (e OrderConfirmed) OccurredAt() time.Time { return e.OccurredAt_ }

// 聚合发布事件
func (o *Order) AddDomainEvent(event DomainEvent) {
    o.events = append(o.events, event)
}

// 应用服务处理事件
func (s *OrderService) ConfirmOrder(orderID OrderID) error {
    order, err := s.repo.Get(orderID)
    if err != nil {
        return err
    }

    if err := order.Confirm(); err != nil {
        return err
    }

    // 保存聚合和事件
    if err := s.repo.Save(order); err != nil {
        return err
    }

    // 发布事件
    for _, event := range order.PullDomainEvents() {
        s.eventBus.Publish(event)
    }

    return nil
}
```

---

## 分层架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          DDD Layered Architecture                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          User Interface                             │   │
│  │  ─────────────────────────────────────────────────────────────────  │   │
│  │  • HTTP/REST Handlers                                               │   │
│  │  • gRPC Services                                                    │   │
│  │  • CLI Commands                                                     │   │
│  │  • Event Consumers                                                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Application Layer                           │   │
│  │  ─────────────────────────────────────────────────────────────────  │   │
│  │  • Application Services (Orchestration)                             │   │
│  │  • DTOs (Data Transfer Objects)                                     │   │
│  │  • Use Cases                                                        │   │
│  │  • Transaction Management                                           │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                           Domain Layer                              │   │
│  │  ─────────────────────────────────────────────────────────────────  │   │
│  │  • Entities                                                         │   │
│  │  • Value Objects                                                    │   │
│  │  • Aggregates                                                       │   │
│  │  • Domain Events                                                    │   │
│  │  • Domain Services                                                  │   │
│  │  • Repositories (Interfaces)                                        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Infrastructure Layer                           │   │
│  │  ─────────────────────────────────────────────────────────────────  │   │
│  │  • Repository Implementations (DB, Cache)                           │   │
│  │  • Message Bus (Kafka, RabbitMQ)                                    │   │
│  │  • External Services Clients                                        │   │
│  │  • Configuration                                                    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 微服务拆分指导

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Monolith to Microservices with DDD                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Step 1: Identify Bounded Contexts                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Monolithic E-commerce Application                                  │   │
│  │                                                                     │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐      │   │
│  │  │  User   │ │  Order  │ │ Payment │ │Inventory│ │ Shipping│      │   │
│  │  │ Profile │ │ Process │ │ Process │ │  Mgmt   │ │ Process │      │   │
│  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘      │   │
│  │       │           │           │           │           │            │   │
│  │       └───────────┴───────────┴───────────┴───────────┘            │   │
│  │                           Shared DB                                  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                   │                                          │
│                                   ▼                                          │
│  Step 2: Define Services                                                     │
│  ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌───────────────┐   │
│  │ User Service  │ │ Order Service │ │Payment Service│ │Inventory Serv │   │
│  │               │ │               │ │               │ │               │   │
│  │ • User Profile│ │ • Order Mgmt  │ │ • Payment     │ │ • Stock Mgmt  │   │
│  │ • Auth        │ │ • Order Query │ │ • Refund      │ │ • Warehouse   │   │
│  └───────────────┘ └───────────────┘ └───────────────┘ └───────────────┘   │
│                                                                              │
│  Step 3: Define Integration Patterns                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │  Order Service ─────Event Bus─────► Inventory Service               │   │
│  │       │                               │                             │   │
│  │       │ OrderCreated                  │ ReserveStock                │   │
│  │       │                               │                             │   │
│  │       └──────► Payment Service        │                             │   │
│  │                ProcessPayment         │                             │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 参考文献

1. [Domain-Driven Design](https://www.domainlanguage.com/ddd/) - Eric Evans
2. [Implementing Domain-Driven Design](https://www.amazon.com/Implementing-Domain-Driven-Design-Vaughn-Vernon/dp/0321834577) - Vaughn Vernon
3. [Domain-Driven Design Reference](https://www.domainlanguage.com/wp-content/uploads/2016/05/DDD_Reference_2015-03.pdf) - 速查手册
4. [Patterns, Principles, and Practices of Domain-Driven Design](https://www.amazon.com/Patterns-Principles-Practices-Domain-Driven-Design/dp/1118714709) - Scott Millett
