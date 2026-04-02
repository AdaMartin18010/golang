# AD-005: DDD 战术设计模式 (DDD Tactical Design Patterns)

> **维度**: Application Domains
> **级别**: S (17+ KB)
> **标签**: #ddd #tactical-patterns #aggregate #entity #value-object
> **权威来源**: [Domain-Driven Design](https://www.domainlanguage.com/ddd/) - Eric Evans, [Implementing DDD](https://www.oreilly.com/library/view/implementing-domain-driven-design/9780133039900/) - Vaughn Vernon

---

## 战术模式概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      DDD Tactical Patterns                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Aggregate                                   │    │
│  │                    (Consistency Boundary)                           │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │                      Order (Root)                            │    │    │
│  │  │  - ID: OrderID                                               │    │    │
│  │  │  - Status                                                    │    │    │
│  │  │  - Total                                                     │    │    │
│  │  │  ┌─────────────────────────────────────────────────────────┐ │    │    │
│  │  │  │  OrderItem (Entity)        ShippingAddress (VO)       │ │    │    │
│  │  │  │  - ID: ItemID              - Street                    │ │    │    │
│  │  │  │  - Product                 - City                      │ │    │    │
│  │  │  │  - Quantity                - ZipCode                   │ │    │    │
│  │  │  │  - Price                   - Country                   │ │    │    │
│  │  │  └─────────────────────────────────────────────────────────┘ │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  │         ▲                                                           │    │
│  │         │ Repository                                                │    │
│  │         ▼                                                           │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │              OrderRepository (Interface)                    │    │    │
│  │  │  - Save(order *Order) error                                 │    │    │
│  │  │  - FindByID(id OrderID) (*Order, error)                     │    │    │
│  │  │  - FindByCustomer(customerID CustomerID) ([]*Order, error)  │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Domain Services                                 │    │
│  │  - PricingService                                                   │    │
│  │  - PaymentService                                                   │    │
│  │  - NotificationService                                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Domain Events                                   │    │
│  │  - OrderCreated                                                     │    │
│  │  - OrderPaid                                                        │    │
│  │  - OrderShipped                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心模式详解

### 1. Entity (实体)

```go
package order

import (
    "github.com/google/uuid"
    "time"
)

// Order 实体 - 有唯一标识，生命周期内状态可变
type Order struct {
    id        OrderID
    customerID CustomerID
    items     []OrderItem
    status    OrderStatus
    total     Money
    createdAt time.Time
    updatedAt time.Time

    events []DomainEvent // 领域事件
}

// OrderID 值对象
type OrderID string

func NewOrderID() OrderID {
    return OrderID(uuid.New().String())
}

// 实体通过 ID 标识，而非属性
func (o *Order) ID() OrderID {
    return o.id
}

// 业务方法 - 封装领域逻辑
func (o *Order) AddItem(product Product, quantity int, price Money) error {
    if o.status != StatusPending {
        return ErrOrderCannotBeModified
    }

    item := NewOrderItem(product.ID(), quantity, price)
    o.items = append(o.items, item)
    o.recalculateTotal()
    o.updatedAt = time.Now()

    return nil
}

func (o *Order) Pay(payment Payment) error {
    if o.status != StatusPending {
        return ErrInvalidOrderStatus
    }

    if payment.Amount().LessThan(o.total) {
        return ErrInsufficientPayment
    }

    o.status = StatusPaid
    o.updatedAt = time.Now()

    // 记录领域事件
    o.events = append(o.events, OrderPaidEvent{
        OrderID: o.id,
        Amount:  o.total,
        PaidAt:  time.Now(),
    })

    return nil
}

func (o *Order) recalculateTotal() {
    total := NewMoney(0, o.total.Currency())
    for _, item := range o.items {
        total = total.Add(item.Price().Multiply(item.Quantity()))
    }
    o.total = total
}
```

### 2. Value Object (值对象)

```go
package order

import "fmt"

// Money 值对象 - 无标识，不可变，通过属性比较
type Money struct {
    amount   int64  // 以最小单位存储 (分)
    currency string
}

func NewMoney(amount float64, currency string) Money {
    return Money{
        amount:   int64(amount * 100),
        currency: currency,
    }
}

func (m Money) Amount() float64 {
    return float64(m.amount) / 100
}

func (m Money) Currency() string {
    return m.currency
}

// 不可变：返回新实例
func (m Money) Add(other Money) Money {
    if m.currency != other.currency {
        panic("cannot add different currencies")
    }
    return Money{
        amount:   m.amount + other.amount,
        currency: m.currency,
    }
}

func (m Money) Multiply(n int) Money {
    return Money{
        amount:   m.amount * int64(n),
        currency: m.currency,
    }
}

func (m Money) LessThan(other Money) bool {
    return m.amount < other.amount
}

// 值相等比较
func (m Money) Equals(other Money) bool {
    return m.amount == other.amount && m.currency == other.currency
}

func (m Money) String() string {
    return fmt.Sprintf("%.2f %s", m.Amount(), m.currency)
}
```

### 3. Aggregate (聚合)

```go
// Order 聚合根 - 维护一致性边界
// - Order 是聚合根
// - OrderItem 是内部实体
// - 外部只能通过 Order 操作 OrderItem

type Order struct {
    id     OrderID
    items  []OrderItem  // 内部实体，不直接暴露
    status OrderStatus
}

// 通过聚合根方法操作内部实体
func (o *Order) RemoveItem(itemID OrderItemID) error {
    if o.status != StatusPending {
        return ErrOrderCannotBeModified
    }

    for i, item := range o.items {
        if item.ID() == itemID {
            o.items = append(o.items[:i], o.items[i+1:]...)
            o.recalculateTotal()
            return nil
        }
    }
    return ErrItemNotFound
}

// 不变式 (Invariants) 保护
func (o *Order) validate() error {
    if len(o.items) == 0 {
        return ErrOrderMustHaveItems
    }
    return nil
}
```

### 4. Domain Event (领域事件)

```go
package order

import "time"

// DomainEvent 接口
type DomainEvent interface {
    EventID() string
    AggregateID() string
    EventType() string
    OccurredAt() time.Time
}

// OrderCreatedEvent 订单已创建
type OrderCreatedEvent struct {
    eventID     string
    orderID     OrderID
    customerID  CustomerID
    items       []OrderItem
    total       Money
    occurredAt  time.Time
}

func (e OrderCreatedEvent) EventID() string      { return e.eventID }
func (e OrderCreatedEvent) AggregateID() string  { return string(e.orderID) }
func (e OrderCreatedEvent) EventType() string    { return "order.created" }
func (e OrderCreatedEvent) OccurredAt() time.Time { return e.occurredAt }

// OrderPaidEvent 订单已支付
type OrderPaidEvent struct {
    eventID    string
    orderID    OrderID
    amount     Money
    paidAt     time.Time
    occurredAt time.Time
}

// 事件发布
func (o *Order) Events() []DomainEvent {
    return o.events
}

func (o *Order) ClearEvents() {
    o.events = nil
}
```

### 5. Repository (仓储)

```go
package order

import "context"

// Repository 接口 - 在领域层定义
type Repository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id OrderID) (*Order, error)
    FindByCustomer(ctx context.Context, customerID CustomerID) ([]*Order, error)
    Update(ctx context.Context, order *Order) error
}

// 基础设施层实现
type PostgresRepository struct {
    db *sql.DB
}

func (r *PostgresRepository) Save(ctx context.Context, order *Order) error {
    // 实现保存逻辑
    // 同时保存领域事件到 Outbox 表
    return nil
}
```

---

## 模式关系

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Pattern Relationships                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Bounded Context                             │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │                      Aggregate                                │    │    │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                     │    │    │
│  │  │  │ Entity  │  │ Entity  │  │ Value   │                     │    │    │
│  │  │  │ (Root)  │  │         │  │ Object  │                     │    │    │
│  │  │  └────┬────┘  └─────────┘  └─────────┘                     │    │    │
│  │  │       │                                                    │    │    │
│  │  │  ┌────┴────┐                                               │    │    │
│  │  │  │ Repository (Interface)                                  │    │    │
│  │  │  └─────────┘                                               │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  │                                                                      │    │
│  │  DomainService ──► 跨聚合业务逻辑                                    │    │
│  │  DomainEvent   ──► 领域事件                                          │    │
│  │  Factory       ──► 复杂对象创建                                       │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  分层架构:                                                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Application Layer (Use Cases, DTOs)                               │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │  Domain Layer (Aggregates, Entities, Value Objects, Services)      │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │  Infrastructure Layer (Repositories, Messaging, External APIs)     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 参考文献

1. [Domain-Driven Design](https://www.domainlanguage.com/ddd/) - Eric Evans
2. [Implementing Domain-Driven Design](https://www.oreilly.com/library/view/implementing-domain-driven-design/9780133039900/) - Vaughn Vernon
3. [DDD Reference](https://www.domainlanguage.com/ddd/reference/) - Eric Evans
