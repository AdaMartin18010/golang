# EC-040: Aggregate Pattern (聚合模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #aggregate #ddd #consistency-boundary #transaction
> **权威来源**:
>
> - [Aggregate Pattern](https://martinfowler.com/bliki/DDD_Aggregate.html) - Martin Fowler
> - [Domain-Driven Design](https://domainlanguage.com/ddd/) - Eric Evans
> - [Implementing Domain-Driven Design](https://www.oreilly.com/library/view/implementing-domain-driven-design/9780133039900/) - Vaughn Vernon

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在复杂领域模型中，如何界定一致性边界，确保业务规则的完整性，同时保持模型的可理解性和性能？

**形式化描述**:

```
给定: 领域模型 M = {E₁, E₂, ..., Eₙ}，其中 E 是实体
给定: 业务规则集合 R = {r₁, r₂, ..., rₘ}，每个规则涉及特定实体
约束:
  - 每个事务只能修改一个一致性边界内的数据
  - 大聚合影响性能
  - 分布式事务难以扩展
目标: 找到最优聚合划分，使得：
  - 业务规则完整性最大化
  - 聚合大小合理
  - 支持可扩展性
```

**大聚合的问题**:

```
反模式: 上帝聚合 (God Aggregate)
┌─────────────────────────────────────────────────────────────────────────┐
│                    Order (God Aggregate)                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Order                                                                  │
│  ├── OrderItems (100+)                                                  │
│  ├── Customer (完整信息)                                                 │
│  ├── PaymentInfo (历史记录)                                              │
│  ├── ShippingInfo (跟踪信息)                                             │
│  ├── Invoices (多个)                                                     │
│  ├── Returns (历史)                                                      │
│  └── Reviews (客户评价)                                                  │
│                                                                         │
│  问题:                                                                   │
│  • 加载整个聚合消耗大量内存                                              │
│  • 并发修改冲突频繁（版本号频繁变化）                                       │
│  • 事务边界过大，影响性能                                                  │
│  • 团队无法独立工作                                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.2 解决方案形式化

**定义 1.1 (聚合)**
聚合是一组相关对象的集合，被视为一个单一的数据修改单元：

```
Aggregate A:
  A = ⟨Root, Entities, ValueObjects, Invariants⟩

其中:
  - Root: 聚合根实体，唯一外部引用点
  - Entities: 聚合内的实体集合
  - ValueObjects: 值对象集合
  - Invariants: 必须始终保持的业务规则
```

**一致性边界**:

```
事务边界:
  ∀transaction T: modifies(T) ⊆ Aᵢ

引用规则:
  external_reference(Aᵢ) = Root(Aᵢ)
  ¬(∃E ∈ Aⱼ, j≠i: direct_reference(E) from Aᵢ)
```

**定义 1.2 (聚合根)**
聚合根是聚合的入口点：

- 具有全局唯一标识
- 负责维护聚合不变量
- 控制聚合内实体的生命周期
- 外部只能通过聚合根引用聚合

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Aggregate Architecture                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                      Order Aggregate                             │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │                   Order (Aggregate Root)                   │  │   │
│  │  │  - ID: OrderID (全局唯一)                                   │  │   │
│  │  │  - CustomerID (引用其他聚合的ID)                            │  │   │
│  │  │  - Status                                                   │  │   │
│  │  │  - Total                                                    │  │   │
│  │  │  - ShippingAddress                                          │  │   │
│  │  │                                                             │  │   │
│  │  │  不变量:                                                     │  │   │
│  │  │  • Total = Σ(OrderItem.Price × OrderItem.Quantity)          │  │   │
│  │  │  • Status 转换必须符合状态机                                 │  │   │
│  │  │  • 已发货订单不能取消                                         │  │   │
│  │  └──────────────────┬──────────────────────────────────────────┘  │   │
│  │                     │                                             │   │
│  │                     │ 包含 (By Reference)                          │   │
│  │                     ▼                                             │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │                  OrderItem (Entity)                        │  │   │
│  │  │  - ID: OrderItemID (局部唯一，聚合内唯一)                    │  │   │
│  │  │  - ProductID (引用其他聚合)                                  │  │   │
│  │  │  - Quantity                                                │  │   │
│  │  │  - Price                                                   │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  │                     │                                             │   │
│  │                     │ 包含 (By Value)                              │   │
│  │                     ▼                                             │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │              Money (Value Object)                          │  │   │
│  │  │  - Amount                                                  │  │   │
│  │  │  - Currency                                                │  │   │
│  │  │  特性: 不可变，值相等即相等                                   │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  聚合间引用（通过ID）:                                                     │
│  ┌─────────────┐      ID      ┌─────────────┐      ID      ┌─────────┐ │
│  │   Order     │─────────────►│  Customer   │              │ Product │ │
│  │ (Aggregate) │              │ (Aggregate) │◄─────────────│         │ │
│  └─────────────┘              └─────────────┘              └─────────┘ │
│                                                                         │
│  规则:                                                                   │
│  • 事务边界 = 聚合边界                                                    │
│  • 聚合内强一致性，聚合间最终一致性                                          │
│  • 删除聚合根同时删除聚合内所有对象                                          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心聚合实现

```go
// aggregate/core.go
package aggregate

import (
    "context"
    "fmt"
    "reflect"
    "sync"
)

// Entity 实体接口
type Entity interface {
    ID() string
    Equals(other Entity) bool
}

// ValueObject 值对象接口
type ValueObject interface {
    Equals(other ValueObject) bool
    Clone() ValueObject
}

// AggregateRoot 聚合根接口
type AggregateRoot interface {
    Entity
    Version() int
    ApplyEvent(event DomainEvent) error
    UncommittedEvents() []DomainEvent
    ClearUncommittedEvents()
    Validate() error
}

// DomainEvent 领域事件接口
type DomainEvent interface {
    EventType() string
    AggregateID() string
    OccurredAt() interface{}
}

// AggregateBase 聚合根基础
type AggregateBase struct {
    id          string
    version     int
    events      []DomainEvent
    mu          sync.RWMutex
}

// NewAggregateBase 创建聚合根基础
func NewAggregateBase(id string) *AggregateBase {
    return &AggregateBase{
        id:      id,
        version: 0,
        events:  make([]DomainEvent, 0),
    }
}

func (a *AggregateBase) ID() string       { return a.id }
func (a *AggregateBase) Version() int     { return a.version }

func (a *AggregateBase) UncommittedEvents() []DomainEvent {
    a.mu.RLock()
    defer a.mu.RUnlock()
    result := make([]DomainEvent, len(a.events))
    copy(result, a.events)
    return result
}

func (a *AggregateBase) ClearUncommittedEvents() {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.events = a.events[:0]
}

// RecordEvent 记录事件
func (a *AggregateBase) RecordEvent(event DomainEvent) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.events = append(a.events, event)
}

// IncrementVersion 增加版本
func (a *AggregateBase) IncrementVersion() {
    a.version++
}

// Repository 仓储接口
type Repository interface {
    Save(ctx context.Context, aggregate AggregateRoot) error
    Get(ctx context.Context, id string) (AggregateRoot, error)
    Exists(ctx context.Context, id string) (bool, error)
}

// UnitOfWork 工作单元
type UnitOfWork struct {
    aggregates map[string]AggregateRoot
    repository Repository
    mu         sync.Mutex
}

// NewUnitOfWork 创建工作单元
func NewUnitOfWork(repo Repository) *UnitOfWork {
    return &UnitOfWork{
        aggregates: make(map[string]AggregateRoot),
        repository: repo,
    }
}

// Register 注册聚合
func (u *UnitOfWork) Register(aggregate AggregateRoot) {
    u.mu.Lock()
    defer u.mu.Unlock()
    u.aggregates[aggregate.ID()] = aggregate
}

// Commit 提交变更
func (u *UnitOfWork) Commit(ctx context.Context) error {
    u.mu.Lock()
    defer u.mu.Unlock()

    for _, aggregate := range u.aggregates {
        if err := aggregate.Validate(); err != nil {
            return fmt.Errorf("validation failed for aggregate %s: %w", aggregate.ID(), err)
        }

        if err := u.repository.Save(ctx, aggregate); err != nil {
            return fmt.Errorf("failed to save aggregate %s: %w", aggregate.ID(), err)
        }

        aggregate.ClearUncommittedEvents()
    }

    u.aggregates = make(map[string]AggregateRoot)
    return nil
}

// Rollback 回滚
func (u *UnitOfWork) Rollback() {
    u.mu.Lock()
    defer u.mu.Unlock()
    u.aggregates = make(map[string]AggregateRoot)
}
```

### 2.2 订单聚合实现

```go
// aggregate/order.go
package aggregate

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/google/uuid"
)

// Order 订单聚合
type Order struct {
    *AggregateBase
    customerID string
    items      []*OrderItem
    total      Money
    status     OrderStatus
    address    Address
}

// OrderStatus 订单状态
type OrderStatus int

const (
    OrderStatusPending OrderStatus = iota
    OrderStatusPaid
    OrderStatusShipped
    OrderStatusDelivered
    OrderStatusCancelled
)

func (s OrderStatus) String() string {
    names := []string{"Pending", "Paid", "Shipped", "Delivered", "Cancelled"}
    if int(s) < len(names) {
        return names[s]
    }
    return "Unknown"
}

// CanTransitionTo 检查状态转换是否有效
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
    transitions := map[OrderStatus][]OrderStatus{
        OrderStatusPending:   {OrderStatusPaid, OrderStatusCancelled},
        OrderStatusPaid:      {OrderStatusShipped, OrderStatusCancelled},
        OrderStatusShipped:   {OrderStatusDelivered},
        OrderStatusDelivered: {},
        OrderStatusCancelled: {},
    }

    allowed, exists := transitions[s]
    if !exists {
        return false
    }

    for _, status := range allowed {
        if status == target {
            return true
        }
    }
    return false
}

// OrderItem 订单项
type OrderItem struct {
    productID string
    name      string
    quantity  int
    price     Money
}

// Money 值对象
type Money struct {
    amount   float64
    currency string
}

// Equals 值对象相等比较
func (m Money) Equals(other ValueObject) bool {
    o, ok := other.(Money)
    if !ok {
        return false
    }
    return m.amount == o.amount && m.currency == o.currency
}

// Clone 克隆值对象
func (m Money) Clone() ValueObject {
    return Money{amount: m.amount, currency: m.currency}
}

func (m Money) Amount() float64   { return m.amount }
func (m Money) Currency() string  { return m.currency }

// Add 金额相加
func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

// Multiply 金额相乘
func (m Money) Multiply(factor int) Money {
    return Money{amount: m.amount * float64(factor), currency: m.currency}
}

// Address 值对象
type Address struct {
    street  string
    city    string
    zipCode string
    country string
}

func (a Address) Equals(other ValueObject) bool {
    o, ok := other.(Address)
    if !ok {
        return false
    }
    return a.street == o.street && a.city == o.city &&
           a.zipCode == o.zipCode && a.country == o.country
}

func (a Address) Clone() ValueObject {
    return Address{street: a.street, city: a.city, zipCode: a.zipCode, country: a.country}
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
    OrderID    string    `json:"order_id"`
    CustomerID string    `json:"customer_id"`
    Total      float64   `json:"total"`
    Currency   string    `json:"currency"`
    Timestamp  time.Time `json:"timestamp"`
}

func (e OrderCreatedEvent) EventType() string    { return "OrderCreated" }
func (e OrderCreatedEvent) AggregateID() string  { return e.OrderID }
func (e OrderCreatedEvent) OccurredAt() interface{} { return e.Timestamp }

// OrderItemAddedEvent 订单项添加事件
type OrderItemAddedEvent struct {
    OrderID   string  `json:"order_id"`
    ProductID string  `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

func (e OrderItemAddedEvent) EventType() string    { return "OrderItemAdded" }
func (e OrderItemAddedEvent) AggregateID() string  { return e.OrderID }
func (e OrderItemAddedEvent) OccurredAt() interface{} { return time.Now() }

// OrderStatusChangedEvent 订单状态变更事件
type OrderStatusChangedEvent struct {
    OrderID   string    `json:"order_id"`
    OldStatus string    `json:"old_status"`
    NewStatus string    `json:"new_status"`
    Timestamp time.Time `json:"timestamp"`
}

func (e OrderStatusChangedEvent) EventType() string    { return "OrderStatusChanged" }
func (e OrderStatusChangedEvent) AggregateID() string  { return e.OrderID }
func (e OrderStatusChangedEvent) OccurredAt() interface{} { return e.Timestamp }

// NewOrder 创建新订单
func NewOrder(customerID string, address Address) *Order {
    order := &Order{
        AggregateBase: NewAggregateBase(uuid.New().String()),
        customerID:    customerID,
        items:         make([]*OrderItem, 0),
        total:         Money{amount: 0, currency: "USD"},
        status:        OrderStatusPending,
        address:       address,
    }

    order.RecordEvent(OrderCreatedEvent{
        OrderID:    order.ID(),
        CustomerID: customerID,
        Total:      0,
        Currency:   "USD",
        Timestamp:  time.Now(),
    })

    return order
}

// CustomerID 获取客户ID
func (o *Order) CustomerID() string { return o.customerID }

// Total 获取总价
func (o *Order) Total() Money { return o.total }

// Status 获取状态
func (o *Order) Status() OrderStatus { return o.status }

// Items 获取订单项
func (o *Order) Items() []*OrderItem {
    result := make([]*OrderItem, len(o.items))
    for i, item := range o.items {
        result[i] = &OrderItem{
            productID: item.productID,
            name:      item.name,
            quantity:  item.quantity,
            price:     item.price.Clone().(Money),
        }
    }
    return result
}

// AddItem 添加订单项
func (o *Order) AddItem(productID string, name string, quantity int, price Money) error {
    if o.status != OrderStatusPending {
        return errors.New("cannot modify items for non-pending order")
    }

    if quantity <= 0 {
        return errors.New("quantity must be positive")
    }

    // 检查是否已存在相同商品
    for _, item := range o.items {
        if item.productID == productID {
            return errors.New("item already exists, use UpdateItem instead")
        }
    }

    item := &OrderItem{
        productID: productID,
        name:      name,
        quantity:  quantity,
        price:     price,
    }

    o.items = append(o.items, item)
    o.recalculateTotal()

    o.RecordEvent(OrderItemAddedEvent{
        OrderID:   o.ID(),
        ProductID: productID,
        Quantity:  quantity,
        Price:     price.Amount(),
    })

    return nil
}

// recalculateTotal 重新计算总价（不变量）
func (o *Order) recalculateTotal() {
    var total Money
    for _, item := range o.items {
        itemTotal := item.price.Multiply(item.quantity)
        newTotal, _ := total.Add(itemTotal)
        total = newTotal
    }
    o.total = total
}

// ChangeStatus 变更状态
func (o *Order) ChangeStatus(newStatus OrderStatus) error {
    if !o.status.CanTransitionTo(newStatus) {
        return fmt.Errorf("invalid status transition from %s to %s", o.status, newStatus)
    }

    oldStatus := o.status
    o.status = newStatus

    o.RecordEvent(OrderStatusChangedEvent{
        OrderID:   o.ID(),
        OldStatus: oldStatus.String(),
        NewStatus: newStatus.String(),
        Timestamp: time.Now(),
    })

    return nil
}

// Pay 支付订单
func (o *Order) Pay() error {
    if o.total.Amount() <= 0 {
        return errors.New("cannot pay for empty order")
    }
    return o.ChangeStatus(OrderStatusPaid)
}

// Ship 发货
func (o *Order) Ship() error {
    return o.ChangeStatus(OrderStatusShipped)
}

// Cancel 取消订单
func (o *Order) Cancel() error {
    if o.status == OrderStatusShipped || o.status == OrderStatusDelivered {
        return errors.New("cannot cancel shipped or delivered order")
    }
    return o.ChangeStatus(OrderStatusCancelled)
}

// Validate 验证聚合不变量
func (o *Order) Validate() error {
    if strings.TrimSpace(o.customerID) == "" {
        return errors.New("customer ID is required")
    }

    if len(o.items) == 0 && o.status != OrderStatusPending {
        return errors.New("order must have items")
    }

    // 验证总价不变量
    var calculatedTotal Money
    for _, item := range o.items {
        itemTotal := item.price.Multiply(item.quantity)
        newTotal, _ := calculatedTotal.Add(itemTotal)
        calculatedTotal = newTotal
    }

    if !o.total.Equals(calculatedTotal) {
        return fmt.Errorf("total invariant violated: expected %v, got %v", calculatedTotal, o.total)
    }

    return nil
}

// ApplyEvent 应用事件（事件溯源）
func (o *Order) ApplyEvent(event DomainEvent) error {
    payload, _ := json.Marshal(event)

    switch e := event.(type) {
    case OrderCreatedEvent:
        o.AggregateBase = NewAggregateBase(e.OrderID)
        o.customerID = e.CustomerID
        o.total = Money{amount: e.Total, currency: e.Currency}
        o.status = OrderStatusPending
        o.items = make([]*OrderItem, 0)

    case OrderItemAddedEvent:
        item := &OrderItem{
            productID: e.ProductID,
            quantity:  e.Quantity,
            price:     Money{amount: e.Price, currency: o.total.Currency()},
        }
        o.items = append(o.items, item)
        o.recalculateTotal()

    case OrderStatusChangedEvent:
        switch e.NewStatus {
        case "Paid":
            o.status = OrderStatusPaid
        case "Shipped":
            o.status = OrderStatusShipped
        case "Delivered":
            o.status = OrderStatusDelivered
        case "Cancelled":
            o.status = OrderStatusCancelled
        }
    }

    o.IncrementVersion()
    return nil
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// aggregate/order_test.go
package aggregate

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
    address := Address{street: "123 Main St", city: "NYC", country: "USA"}
    order := NewOrder("customer-001", address)

    require.NotNil(t, order)
    assert.NotEmpty(t, order.ID())
    assert.Equal(t, "customer-001", order.CustomerID())
    assert.Equal(t, OrderStatusPending, order.Status())
    assert.Equal(t, 0.0, order.Total().Amount())

    // 验证事件已记录
    events := order.UncommittedEvents()
    require.Len(t, events, 1)
    assert.Equal(t, "OrderCreated", events[0].EventType())
}

func TestOrder_AddItem(t *testing.T) {
    address := Address{street: "123 Main St", city: "NYC", country: "USA"}
    order := NewOrder("customer-001", address)
    order.ClearUncommittedEvents()

    price := Money{amount: 10.0, currency: "USD"}
    err := order.AddItem("PROD-001", "Product 1", 2, price)

    require.NoError(t, err)
    assert.Len(t, order.Items(), 1)
    assert.Equal(t, 20.0, order.Total().Amount())

    // 验证事件已记录
    events := order.UncommittedEvents()
    require.Len(t, events, 1)
    assert.Equal(t, "OrderItemAdded", events[0].EventType())
}

func TestOrder_AddItem_InvalidQuantity(t *testing.T) {
    order := NewOrder("customer-001", Address{})

    err := order.AddItem("PROD-001", "Product 1", 0, Money{amount: 10.0})

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "quantity must be positive")
}

func TestOrder_StatusTransitions(t *testing.T) {
    tests := []struct {
        name      string
        initial   OrderStatus
        target    OrderStatus
        wantError bool
    }{
        {"Pending to Paid", OrderStatusPending, OrderStatusPaid, false},
        {"Pending to Cancelled", OrderStatusPending, OrderStatusCancelled, false},
        {"Paid to Shipped", OrderStatusPaid, OrderStatusShipped, false},
        {"Paid to Cancelled", OrderStatusPaid, OrderStatusCancelled, false},
        {"Shipped to Delivered", OrderStatusShipped, OrderStatusDelivered, false},
        {"Delivered to Cancelled", OrderStatusDelivered, OrderStatusCancelled, true},
        {"Cancelled to Paid", OrderStatusCancelled, OrderStatusPaid, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            order := NewOrder("customer-001", Address{})
            order.status = tt.initial // 直接设置状态进行测试

            err := order.ChangeStatus(tt.target)

            if tt.wantError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.target, order.Status())
            }
        })
    }
}

func TestOrder_Validate(t *testing.T) {
    t.Run("valid order", func(t *testing.T) {
        order := NewOrder("customer-001", Address{})
        order.AddItem("PROD-001", "Product 1", 2, Money{amount: 10.0, currency: "USD"})

        err := order.Validate()
        assert.NoError(t, err)
    })

    t.Run("missing customer ID", func(t *testing.T) {
        order := NewOrder("", Address{})

        err := order.Validate()
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "customer ID is required")
    })

    t.Run("invariant violated", func(t *testing.T) {
        order := NewOrder("customer-001", Address{})
        order.AddItem("PROD-001", "Product 1", 2, Money{amount: 10.0, currency: "USD"})
        order.total = Money{amount: 999.0, currency: "USD"} // 手动破坏不变量

        err := order.Validate()
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "total invariant violated")
    })
}

func TestMoney_Operations(t *testing.T) {
    m1 := Money{amount: 10.0, currency: "USD"}
    m2 := Money{amount: 20.0, currency: "USD"}

    t.Run("addition", func(t *testing.T) {
        result, err := m1.Add(m2)
        require.NoError(t, err)
        assert.Equal(t, 30.0, result.Amount())
    })

    t.Run("addition different currencies", func(t *testing.T) {
        m3 := Money{amount: 10.0, currency: "EUR"}
        _, err := m1.Add(m3)
        assert.Error(t, err)
    })

    t.Run("multiplication", func(t *testing.T) {
        result := m1.Multiply(3)
        assert.Equal(t, 30.0, result.Amount())
    })

    t.Run("equals", func(t *testing.T) {
        m3 := Money{amount: 10.0, currency: "USD"}
        assert.True(t, m1.Equals(m3))
        assert.False(t, m1.Equals(m2))
    })
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Repository 模式的关系

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Aggregate + Repository Pattern                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Repository 负责聚合的持久化，聚合负责业务逻辑:                             │
│                                                                         │
│  ┌─────────────┐         ┌─────────────┐         ┌─────────────────┐   │
│  │ Application │────────►│ Repository  │────────►│   Data Store    │   │
│  │   Service   │         │             │         │                 │   │
│  └─────────────┘         └──────┬──────┘         └─────────────────┘   │
│                                 │                                       │
│                                 │ Aggregate                             │
│                                 ▼                                       │
│                          ┌─────────────┐                               │
│                          │   Order     │                               │
│                          │ (Aggregate) │                               │
│                          │             │                               │
│                          │ - Business  │                               │
│                          │   Logic     │                               │
│                          │ - Invariants│                               │
│                          │ - Events    │                               │
│                          └─────────────┘                               │
│                                                                         │
│  Repository 职责:                                                        │
│  • 保存聚合（整个一致性边界）                                             │
│  • 通过聚合根ID加载聚合                                                   │
│  • 封装持久化细节                                                         │
│                                                                         │
│  聚合职责:                                                               │
│  • 封装业务逻辑                                                           │
│  • 维护不变量                                                             │
│  • 产生领域事件                                                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 5. 决策标准

### 5.1 聚合设计原则

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Aggregate Design Rules                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. 事务边界规则                                                          │
│     一个事务只修改一个聚合                                                │
│     如果需要跨聚合事务，使用 Saga 模式                                     │
│                                                                         │
│  2. 小聚合原则                                                            │
│     聚合应尽量小，只包含必须一起修改的实体                                   │
│     大聚合影响性能和并发                                                    │
│                                                                         │
│  3. 通过 ID 引用其他聚合                                                   │
│     Order.CustomerID ✓                                                   │
│     Order.Customer ✗ (直接引用)                                           │
│                                                                         │
│  4. 聚合内强一致，聚合间最终一致                                            │
│     使用领域事件实现聚合间同步                                              │
│                                                                         │
│  5. 根实体保护内部对象                                                      │
│     外部只能通过聚合根访问内部对象                                          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Aggregate Design Checklist                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 识别业务不变量                                                        │
│  □ 确定哪些对象必须一起修改                                               │
│  □ 识别聚合根（具有全局标识）                                             │
│  □ 检查聚合大小（实体数量 < 100）                                         │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现聚合根和内部实体                                                   │
│  □ 封装聚合边界                                                           │
│  □ 实现业务规则和不变量检查                                               │
│  □ 产生领域事件                                                           │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 避免大聚合                                                            │
│  ❌ 避免直接引用其他聚合                                                  │
│  ❌ 避免跨聚合事务                                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>17KB, 完整形式化 + Go 实现 + 测试)

**相关文档**:

- [EC-043-Repository-Pattern.md](./EC-043-Repository-Pattern.md)
- [EC-039-Domain-Event-Pattern.md](./EC-039-Domain-Event-Pattern.md)
- [EC-042-Entity-Pattern.md](./EC-042-Entity-Pattern.md)
- [EC-041-Value-Object-Pattern.md](./EC-041-Value-Object-Pattern.md)
