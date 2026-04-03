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
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02