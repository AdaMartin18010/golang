# 领域驱动设计模式 (DDD Patterns)

> **分类**: 成熟应用领域  
> **标签**: #ddd #domain-driven-design #architecture

---

## 实体 (Entity)

```go
// 有唯一标识的对象
type Order struct {
    id        OrderID
    items     []OrderItem
    status    OrderStatus
    createdAt time.Time
}

type OrderID string

func NewOrder(id OrderID) *Order {
    return &Order{
        id:        id,
        status:    OrderStatusPending,
        createdAt: time.Now(),
    }
}

func (o *Order) ID() OrderID {
    return o.id
}

func (o *Order) AddItem(product Product, quantity int) error {
    if quantity <= 0 {
        return ErrInvalidQuantity
    }
    
    o.items = append(o.items, OrderItem{
        Product:  product,
        Quantity: quantity,
    })
    
    return nil
}

func (o *Order) Pay() error {
    if o.status != OrderStatusPending {
        return ErrInvalidStatus
    }
    
    o.status = OrderStatusPaid
    return nil
}
```

---

## 值对象 (Value Object)

```go
// 无标识，不可变
type Money struct {
    amount   decimal.Decimal
    currency string
}

func NewMoney(amount decimal.Decimal, currency string) Money {
    return Money{
        amount:   amount,
        currency: strings.ToUpper(currency),
    }
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, ErrCurrencyMismatch
    }
    return NewMoney(m.amount.Add(other.amount), m.currency), nil
}

func (m Money) Equals(other Money) bool {
    return m.currency == other.currency && m.amount.Equal(other.amount)
}

// 使用
price := NewMoney(decimal.NewFromFloat(100), "USD")
discount := NewMoney(decimal.NewFromFloat(10), "USD")
total, _ := price.Sub(discount)
```

---

## 聚合根 (Aggregate Root)

```go
type OrderAggregate struct {
    Order
    events []DomainEvent
}

func (o *OrderAggregate) Pay(payment Payment) error {
    if err := o.Order.Pay(); err != nil {
        return err
    }
    
    // 记录领域事件
    o.events = append(o.events, OrderPaidEvent{
        OrderID:   o.ID(),
        Amount:    payment.Amount,
        PaidAt:    time.Now(),
    })
    
    return nil
}

func (o *OrderAggregate) UncommittedEvents() []DomainEvent {
    return o.events
}

func (o *OrderAggregate) MarkCommitted() {
    o.events = nil
}
```

---

## 仓储 (Repository)

```go
// 接口定义在领域层
type OrderRepository interface {
    Save(ctx context.Context, order *OrderAggregate) error
    FindByID(ctx context.Context, id OrderID) (*OrderAggregate, error)
    FindByUser(ctx context.Context, userID UserID) ([]*OrderAggregate, error)
    Update(ctx context.Context, order *OrderAggregate) error
    Delete(ctx context.Context, id OrderID) error
}

// 实现基础设置层
type PostgresOrderRepository struct {
    db *sql.DB
}

func (r *PostgresOrderRepository) Save(ctx context.Context, order *OrderAggregate) error {
    // 保存到数据库
    // ...
    
    // 发布领域事件
    for _, event := range order.UncommittedEvents() {
        r.eventBus.Publish(event)
    }
    order.MarkCommitted()
    
    return nil
}
```

---

## 领域服务

```go
// 跨实体的业务逻辑
type PricingService struct {
    discountRepo DiscountRepository
}

func (s *PricingService) CalculatePrice(ctx context.Context, order *Order) (Money, error) {
    var total Money
    
    for _, item := range order.Items {
        itemPrice := item.Product.Price.Multiply(decimal.NewFromInt(int64(item.Quantity)))
        total, _ = total.Add(itemPrice)
    }
    
    // 应用折扣
    discount, err := s.discountRepo.FindApplicable(ctx, order)
    if err == nil {
        total, _ = total.Sub(discount.Amount)
    }
    
    return total, nil
}
```

---

## 应用服务

```go
type OrderApplicationService struct {
    orderRepo   OrderRepository
    pricingSvc  *PricingService
    paymentSvc  PaymentService
}

func (s *OrderApplicationService) PlaceOrder(ctx context.Context, cmd PlaceOrderCommand) error {
    // 1. 创建订单
    order := NewOrder(cmd.OrderID)
    
    for _, item := range cmd.Items {
        order.AddItem(item.Product, item.Quantity)
    }
    
    // 2. 计算价格
    total, err := s.pricingSvc.CalculatePrice(ctx, order)
    if err != nil {
        return err
    }
    
    // 3. 处理支付
    if err := s.paymentSvc.Process(ctx, total); err != nil {
        return err
    }
    
    // 4. 保存
    return s.orderRepo.Save(ctx, order)
}
```
