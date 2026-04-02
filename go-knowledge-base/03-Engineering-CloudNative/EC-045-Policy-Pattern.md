# EC-045: Policy Pattern (策略/政策模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #policy #strategy #rules-engine #business-logic
> **权威来源**:
>
> - [Strategy Pattern](https://en.wikipedia.org/wiki/Strategy_pattern) - Gang of Four
> - [Policy Pattern in DDD](https://domainlanguage.com/ddd/) - Eric Evans
> - [Specification Pattern](https://en.wikipedia.org/wiki/Specification_pattern) - Evans/Fowler

---

## 1. 模式形式化定义

### 1.2 问题定义

**问题陈述**: 在领域模型中，如何封装和隔离经常变化的业务规则或策略，使系统能够灵活地组合和切换不同的策略实现？

**硬编码规则的问题**:

```
问题: 策略逻辑硬编码在领域对象中
┌─────────────────────────────────────────────────────────────────────────┐
│                    Hardcoded Policy Anti-Pattern                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  type Order struct {                                                    │
│      Items []Item                                                       │
│      Total float64                                                      │
│      CustomerType string  // "REGULAR", "VIP", "ENTERPRISE"            │
│  }                                                                      │
│                                                                         │
│  func (o *Order) CalculateDiscount() float64 {                          │
│      // 硬编码的折扣策略                                                 │
│      if o.CustomerType == "VIP" {                                      │
│          return o.Total * 0.2  // VIP 20% 折扣                          │
│      } else if o.CustomerType == "ENTERPRISE" {                        │
│          return o.Total * 0.3  // 企业 30% 折扣                         │
│      }                                                                  │
│      return 0  // 普通客户无折扣                                         │
│  }                                                                      │
│                                                                         │
│  问题:                                                                  │
│  • 添加新策略需要修改 Order 类                                           │
│  • 违反开闭原则                                                          │
│  • 策略逻辑分散在各个地方                                                 │
│  • 难以测试（需要创建完整的 Order）                                       │
│  • 无法运行时动态切换策略                                                 │
│                                                                         │
│  新需求: "双十一期间所有客户 15% 折扣"                                     │
│  → 需要再次修改 Order.CalculateDiscount()                                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**形式化描述**:

```
给定: 领域行为 B，其具体实现可能变化
给定: 策略集合 S = {s₁, s₂, ..., sₙ}，每个 s 实现 B 的不同变体
约束:
  - 策略切换不影响客户端代码
  - 新策略可以动态添加
  - 策略可以组合
目标: 设计机制使得: client ──uses──► Policy ──delegates──► ConcreteStrategy
```

### 1.2 解决方案形式化

**定义 1.1 (策略模式)**
策略模式定义一系列算法，将它们封装起来，并且使它们可以互相替换：

```
策略接口 P:
  P = ⟨Execute, Validate, Name⟩

具体策略 Sᵢ:
  Sᵢ implements P
  Sᵢ.Execute(context) → result

上下文 C 使用策略:
  C.policy = Sᵢ
  C.Execute() → policy.Execute(context)
```

**定义 1.2 (组合策略)**

```
策略可以组合:
  CompositePolicy = [S₁, S₂, ..., Sₙ]
  Execute(Composite) = fold(Execute, context, [S₁, S₂, ..., Sₙ])

顺序执行:
  result₀ = context
  result₁ = S₁.Execute(result₀)
  result₂ = S₂.Execute(result₁)
  ...
  resultₙ = Sₙ.Execute(resultₙ₋₁)
```

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Policy/Strategy Architecture                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                        Context                                   │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │                   Order (Aggregate)                        │  │   │
│  │  │                                                            │  │   │
│  │  │  func (o *Order) ApplyDiscount() {                        │  │   │
│  │  │      discount := o.pricingPolicy.Calculate(o)             │  │   │
│  │  │      o.Total -= discount                                  │  │   │
│  │  │  }                                                        │  │   │
│  │  │                                                            │  │   │
│  │  │  // 不需要知道具体策略                                      │  │   │
│  │  │  // 只需知道 PricingPolicy 接口                             │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  │                            │                                     │   │
│  │                            │ uses                                │   │
│  │                            ▼                                     │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │                PricingPolicy Interface                     │  │   │
│  │  │                                                            │  │   │
│  │  │  type PricingPolicy interface {                           │  │   │
│  │  │      Calculate(order *Order) Money                        │  │   │
│  │  │      Name() string                                        │  │   │
│  │  │  }                                                        │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  │              ┌───────────────┬────────────────┐                  │   │
│  │              │               │                │                  │   │
│  │              ▼               ▼                ▼                  │   │
│  │  ┌─────────────────┐ ┌──────────────┐ ┌──────────────────┐     │   │
│  │  │  VIPPolicy      │ │ SeasonalSale │ │ EnterprisePolicy │     │   │
│  │  │                 │ │    Policy    │ │                  │     │   │
│  │  │  return total   │ │              │ │  return total    │     │   │
│  │  │  * 0.2          │ │  return      │ │  * 0.3           │     │   │
│  │  │                 │ │  total * 0.15│ │                  │     │   │
│  │  └─────────────────┘ └──────────────┘ └──────────────────┘     │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  组合策略示例:                                                           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                 CompositePolicy (Chain)                          │   │
│  │                                                                  │   │
│  │  1. BaseDiscountPolicy:     计算基础折扣（VIP/企业）               │   │
│  │  2. SeasonalPolicy:         应用季节性折扣                        │   │
│  │  3. CouponPolicy:           应用优惠券                            │   │
│  │  4. MinimumPricePolicy:     确保不低于最低价                      │   │
│  │                                                                  │   │
│  │  总折扣 = 依次应用所有策略                                          │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  策略优势:                                                               │
│  • 新策略无需修改现有代码                                                 │
│  • 策略可以独立测试                                                       │
│  • 可以运行时动态选择策略                                                  │
│  • 策略可以组合使用                                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心策略实现

```go
// policy/core.go
package policy

import (
    "context"
    "fmt"
)

// Policy 策略接口
type Policy interface {
    // Name 策略名称
    Name() string

    // Execute 执行策略
    Execute(ctx context.Context, input interface{}) (interface{}, error)

    // Validate 验证输入
    Validate(input interface{}) error

    // IsApplicable 检查是否适用于输入
    IsApplicable(input interface{}) bool
}

// PolicyResult 策略执行结果
type PolicyResult struct {
    PolicyName string
    Success    bool
    Data       interface{}
    Error      error
    Applied    bool // 是否实际应用
}

// CompositePolicy 组合策略
type CompositePolicy struct {
    name     string
    policies []Policy
    mode     CompositeMode
}

// CompositeMode 组合模式
type CompositeMode int

const (
    ModeChain CompositeMode = iota   // 链式执行，一个结果传给下一个
    ModeFirstMatch                   // 第一个匹配的执行
    ModeAll                          // 全部执行，返回所有结果
)

// NewCompositePolicy 创建组合策略
func NewCompositePolicy(name string, mode CompositeMode, policies ...Policy) *CompositePolicy {
    return &CompositePolicy{
        name:     name,
        policies: policies,
        mode:     mode,
    }
}

// Name 策略名称
func (c *CompositePolicy) Name() string {
    return c.name
}

// Execute 执行组合策略
func (c *CompositePolicy) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    switch c.mode {
    case ModeChain:
        return c.executeChain(ctx, input)
    case ModeFirstMatch:
        return c.executeFirstMatch(ctx, input)
    case ModeAll:
        return c.executeAll(ctx, input)
    default:
        return nil, fmt.Errorf("unknown composite mode: %d", c.mode)
    }
}

func (c *CompositePolicy) executeChain(ctx context.Context, input interface{}) (interface{}, error) {
    result := input
    for _, policy := range c.policies {
        if !policy.IsApplicable(result) {
            continue
        }

        r, err := policy.Execute(ctx, result)
        if err != nil {
            return nil, fmt.Errorf("policy %s failed: %w", policy.Name(), err)
        }
        result = r
    }
    return result, nil
}

func (c *CompositePolicy) executeFirstMatch(ctx context.Context, input interface{}) (interface{}, error) {
    for _, policy := range c.policies {
        if policy.IsApplicable(input) {
            return policy.Execute(ctx, input)
        }
    }
    return input, nil // 没有匹配的策略，返回原值
}

func (c *CompositePolicy) executeAll(ctx context.Context, input interface{}) (interface{}, error) {
    results := make([]PolicyResult, 0, len(c.policies))

    for _, policy := range c.policies {
        result := PolicyResult{
            PolicyName: policy.Name(),
            Applied:    policy.IsApplicable(input),
        }

        if result.Applied {
            data, err := policy.Execute(ctx, input)
            result.Success = err == nil
            result.Data = data
            result.Error = err
        }

        results = append(results, result)
    }

    return results, nil
}

// Validate 验证
func (c *CompositePolicy) Validate(input interface{}) error {
    for _, policy := range c.policies {
        if err := policy.Validate(input); err != nil {
            return fmt.Errorf("policy %s validation failed: %w", policy.Name(), err)
        }
    }
    return nil
}

// IsApplicable 是否适用
func (c *CompositePolicy) IsApplicable(input interface{}) bool {
    for _, policy := range c.policies {
        if policy.IsApplicable(input) {
            return true
        }
    }
    return false
}

// PolicyRegistry 策略注册表
type PolicyRegistry struct {
    policies map[string]Policy
}

// NewPolicyRegistry 创建策略注册表
func NewPolicyRegistry() *PolicyRegistry {
    return &PolicyRegistry{
        policies: make(map[string]Policy),
    }
}

// Register 注册策略
func (r *PolicyRegistry) Register(policy Policy) {
    r.policies[policy.Name()] = policy
}

// Get 获取策略
func (r *PolicyRegistry) Get(name string) (Policy, bool) {
    p, exists := r.policies[name]
    return p, exists
}

// List 列出所有策略
func (r *PolicyRegistry) List() []Policy {
    result := make([]Policy, 0, len(r.policies))
    for _, p := range r.policies {
        result = append(result, p)
    }
    return result
}
```

### 2.2 定价策略实现

```go
// policy/pricing_policy.go
package policy

import (
    "context"
    "errors"
    "fmt"
    "time"
)

// Order 订单
type Order struct {
    ID         string
    Customer   Customer
    Items      []OrderItem
    Total      Money
    CreatedAt  time.Time
    CouponCode string
}

// Customer 客户
type Customer struct {
    ID   string
    Type CustomerType // Regular, VIP, Enterprise
}

// CustomerType 客户类型
type CustomerType string

const (
    CustomerTypeRegular    CustomerType = "REGULAR"
    CustomerTypeVIP        CustomerType = "VIP"
    CustomerTypeEnterprise CustomerType = "ENTERPRISE"
)

// OrderItem 订单项
type OrderItem struct {
    ProductID string
    Quantity  int
    Price     Money
}

// Money 金额
type Money struct {
    Amount   float64
    Currency string
}

// Add 相加
func (m Money) Add(other Money) Money {
    return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}
}

// Multiply 相乘
func (m Money) Multiply(factor float64) Money {
    return Money{Amount: m.Amount * factor, Currency: m.Currency}
}

// DiscountResult 折扣结果
type DiscountResult struct {
    OriginalTotal Money
    Discount      Money
    FinalTotal    Money
    AppliedRules  []string
}

// VIPDiscountPolicy VIP 折扣策略
type VIPDiscountPolicy struct {
    discountRate float64
}

// NewVIPDiscountPolicy 创建 VIP 折扣策略
func NewVIPDiscountPolicy(rate float64) *VIPDiscountPolicy {
    return &VIPDiscountPolicy{discountRate: rate}
}

func (p *VIPDiscountPolicy) Name() string {
    return "VIPDiscount"
}

func (p *VIPDiscountPolicy) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    order, ok := input.(*Order)
    if !ok {
        return nil, errors.New("input must be *Order")
    }

    discount := order.Total.Multiply(p.discountRate)
    finalTotal := Money{Amount: order.Total.Amount - discount.Amount, Currency: order.Total.Currency}

    return &DiscountResult{
        OriginalTotal: order.Total,
        Discount:      discount,
        FinalTotal:    finalTotal,
        AppliedRules:  []string{"VIP_20%_OFF"},
    }, nil
}

func (p *VIPDiscountPolicy) Validate(input interface{}) error {
    _, ok := input.(*Order)
    if !ok {
        return errors.New("input must be *Order")
    }
    return nil
}

func (p *VIPDiscountPolicy) IsApplicable(input interface{}) bool {
    order, ok := input.(*Order)
    if !ok {
        return false
    }
    return order.Customer.Type == CustomerTypeVIP
}

// EnterpriseDiscountPolicy 企业折扣策略
type EnterpriseDiscountPolicy struct {
    discountRate float64
}

// NewEnterpriseDiscountPolicy 创建企业折扣策略
func NewEnterpriseDiscountPolicy(rate float64) *EnterpriseDiscountPolicy {
    return &EnterpriseDiscountPolicy{discountRate: rate}
}

func (p *EnterpriseDiscountPolicy) Name() string {
    return "EnterpriseDiscount"
}

func (p *EnterpriseDiscountPolicy) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    order, ok := input.(*Order)
    if !ok {
        return nil, errors.New("input must be *Order")
    }

    discount := order.Total.Multiply(p.discountRate)
    finalTotal := Money{Amount: order.Total.Amount - discount.Amount, Currency: order.Total.Currency}

    return &DiscountResult{
        OriginalTotal: order.Total,
        Discount:      discount,
        FinalTotal:    finalTotal,
        AppliedRules:  []string{"ENTERPRISE_30%_OFF"},
    }, nil
}

func (p *EnterpriseDiscountPolicy) Validate(input interface{}) error {
    _, ok := input.(*Order)
    if !ok {
        return errors.New("input must be *Order")
    }
    return nil
}

func (p *EnterpriseDiscountPolicy) IsApplicable(input interface{}) bool {
    order, ok := input.(*Order)
    if !ok {
        return false
    }
    return order.Customer.Type == CustomerTypeEnterprise
}

// SeasonalDiscountPolicy 季节性折扣策略
type SeasonalDiscountPolicy struct {
    startDate    time.Time
    endDate      time.Time
    discountRate float64
    name         string
}

// NewSeasonalDiscountPolicy 创建季节性折扣策略
func NewSeasonalDiscountPolicy(name string, start, end time.Time, rate float64) *SeasonalDiscountPolicy {
    return &SeasonalDiscountPolicy{
        startDate:    start,
        endDate:      end,
        discountRate: rate,
        name:         name,
    }
}

func (p *SeasonalDiscountPolicy) Name() string {
    return p.name
}

func (p *SeasonalDiscountPolicy) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    result, ok := input.(*DiscountResult)
    if !ok {
        order, ok := input.(*Order)
        if !ok {
            return nil, errors.New("input must be *Order or *DiscountResult")
        }
        result = &DiscountResult{
            OriginalTotal: order.Total,
            FinalTotal:    order.Total,
        }
    }

    discount := result.FinalTotal.Multiply(p.discountRate)
    result.FinalTotal = Money{Amount: result.FinalTotal.Amount - discount.Amount, Currency: result.FinalTotal.Currency}
    result.Discount = result.Discount.Add(discount)
    result.AppliedRules = append(result.AppliedRules, fmt.Sprintf("SEASONAL_%s", p.name))

    return result, nil
}

func (p *SeasonalDiscountPolicy) Validate(input interface{}) error {
    return nil
}

func (p *SeasonalDiscountPolicy) IsApplicable(input interface{}) bool {
    var orderTime time.Time

    switch v := input.(type) {
    case *Order:
        orderTime = v.CreatedAt
    case *DiscountResult:
        // 默认当前时间，实际应该从上下文获取
        orderTime = time.Now()
    default:
        return false
    }

    return orderTime.After(p.startDate) && orderTime.Before(p.endDate)
}

// MinimumPricePolicy 最低价格策略
type MinimumPricePolicy struct {
    minimumAmount float64
}

// NewMinimumPricePolicy 创建最低价格策略
func NewMinimumPricePolicy(min float64) *MinimumPricePolicy {
    return &MinimumPricePolicy{minimumAmount: min}
}

func (p *MinimumPricePolicy) Name() string {
    return "MinimumPrice"
}

func (p *MinimumPricePolicy) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    result, ok := input.(*DiscountResult)
    if !ok {
        return nil, errors.New("input must be *DiscountResult")
    }

    if result.FinalTotal.Amount < p.minimumAmount {
        adjustment := p.minimumAmount - result.FinalTotal.Amount
        result.FinalTotal.Amount = p.minimumAmount
        result.Discount.Amount -= adjustment // 减少折扣
        result.AppliedRules = append(result.AppliedRules, "MINIMUM_PRICE_ADJUSTED")
    }

    return result, nil
}

func (p *MinimumPricePolicy) Validate(input interface{}) error {
    return nil
}

func (p *MinimumPricePolicy) IsApplicable(input interface{}) bool {
    _, ok := input.(*DiscountResult)
    return ok
}

// PricingEngine 定价引擎
type PricingEngine struct {
    policy Policy
}

// NewPricingEngine 创建定价引擎
func NewPricingEngine(policy Policy) *PricingEngine {
    return &PricingEngine{policy: policy}
}

// Calculate 计算价格
func (e *PricingEngine) Calculate(ctx context.Context, order *Order) (*DiscountResult, error) {
    result, err := e.policy.Execute(ctx, order)
    if err != nil {
        return nil, err
    }

    return result.(*DiscountResult), nil
}

// SetPolicy 切换策略
func (e *PricingEngine) SetPolicy(policy Policy) {
    e.policy = policy
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// policy/pricing_policy_test.go
package policy

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestVIPDiscountPolicy(t *testing.T) {
    policy := NewVIPDiscountPolicy(0.2) // 20% 折扣

    vipOrder := &Order{
        ID: "order-001",
        Customer: Customer{ID: "cust-001", Type: CustomerTypeVIP},
        Total: Money{Amount: 100.0, Currency: "USD"},
    }

    regularOrder := &Order{
        ID: "order-002",
        Customer: Customer{ID: "cust-002", Type: CustomerTypeRegular},
        Total: Money{Amount: 100.0, Currency: "USD"},
    }

    t.Run("applicable to VIP", func(t *testing.T) {
        assert.True(t, policy.IsApplicable(vipOrder))
    })

    t.Run("not applicable to regular", func(t *testing.T) {
        assert.False(t, policy.IsApplicable(regularOrder))
    })

    t.Run("calculate discount", func(t *testing.T) {
        result, err := policy.Execute(context.Background(), vipOrder)
        require.NoError(t, err)

        discountResult := result.(*DiscountResult)
        assert.Equal(t, 100.0, discountResult.OriginalTotal.Amount)
        assert.Equal(t, 20.0, discountResult.Discount.Amount)
        assert.Equal(t, 80.0, discountResult.FinalTotal.Amount)
        assert.Contains(t, discountResult.AppliedRules, "VIP_20%_OFF")
    })
}

func TestCompositePolicy_Chain(t *testing.T) {
    // 创建组合策略: VIP 折扣 + 季节性折扣 + 最低价格限制
    vipPolicy := NewVIPDiscountPolicy(0.2)

    seasonalPolicy := NewSeasonalDiscountPolicy(
        "SummerSale",
        time.Now().Add(-24*time.Hour),
        time.Now().Add(24*time.Hour),
        0.1, // 额外 10%
    )

    minPricePolicy := NewMinimumPricePolicy(70)

    composite := NewCompositePolicy("VIPSummer", ModeChain, vipPolicy, seasonalPolicy, minPricePolicy)

    order := &Order{
        ID:        "order-001",
        Customer:  Customer{ID: "cust-001", Type: CustomerTypeVIP},
        Total:     Money{Amount: 100.0, Currency: "USD"},
        CreatedAt: time.Now(),
    }

    result, err := composite.Execute(context.Background(), order)
    require.NoError(t, err)

    discountResult := result.(*DiscountResult)

    // VIP 20% 折扣: 100 -> 80
    // 季节性 10% 折扣: 80 -> 72
    // 最低价格 70，所以最终 72
    assert.Equal(t, 100.0, discountResult.OriginalTotal.Amount)
    assert.Equal(t, 72.0, discountResult.FinalTotal.Amount)
    assert.Contains(t, discountResult.AppliedRules, "VIP_20%_OFF")
    assert.Contains(t, discountResult.AppliedRules, "SEASONAL_SummerSale")
}

func TestCompositePolicy_FirstMatch(t *testing.T) {
    vipPolicy := NewVIPDiscountPolicy(0.2)
    enterprisePolicy := NewEnterpriseDiscountPolicy(0.3)

    composite := NewCompositePolicy("TieredDiscount", ModeFirstMatch, enterprisePolicy, vipPolicy)

    // VIP 客户
    vipOrder := &Order{
        Customer: Customer{Type: CustomerTypeVIP},
        Total:    Money{Amount: 100.0},
    }

    result, err := composite.Execute(context.Background(), vipOrder)
    require.NoError(t, err)

    discountResult := result.(*DiscountResult)
    // VIP 匹配第二个策略（VIP 20%）
    assert.Equal(t, 80.0, discountResult.FinalTotal.Amount)

    // 企业客户
    enterpriseOrder := &Order{
        Customer: Customer{Type: CustomerTypeEnterprise},
        Total:    Money{Amount: 100.0},
    }

    result, err = composite.Execute(context.Background(), enterpriseOrder)
    require.NoError(t, err)

    discountResult = result.(*DiscountResult)
    // 企业匹配第一个策略（企业 30%）
    assert.Equal(t, 70.0, discountResult.FinalTotal.Amount)
}

func TestPricingEngine(t *testing.T) {
    policy := NewVIPDiscountPolicy(0.2)
    engine := NewPricingEngine(policy)

    order := &Order{
        Customer: Customer{Type: CustomerTypeVIP},
        Total:    Money{Amount: 100.0, Currency: "USD"},
    }

    result, err := engine.Calculate(context.Background(), order)
    require.NoError(t, err)
    assert.Equal(t, 80.0, result.FinalTotal.Amount)

    // 切换策略
    engine.SetPolicy(NewEnterpriseDiscountPolicy(0.3))
    order.Customer.Type = CustomerTypeEnterprise

    result, err = engine.Calculate(context.Background(), order)
    require.NoError(t, err)
    assert.Equal(t, 70.0, result.FinalTotal.Amount)
}

func TestPolicyRegistry(t *testing.T) {
    registry := NewPolicyRegistry()

    vipPolicy := NewVIPDiscountPolicy(0.2)
    seasonalPolicy := NewSeasonalDiscountPolicy("Holiday", time.Now(), time.Now().Add(24*time.Hour), 0.15)

    registry.Register(vipPolicy)
    registry.Register(seasonalPolicy)

    p, found := registry.Get("VIPDiscount")
    assert.True(t, found)
    assert.Equal(t, "VIPDiscount", p.Name())

    policies := registry.List()
    assert.Len(t, policies, 2)
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Specification 模式的关系

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Policy + Specification Integration                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Specification 决定策略是否适用，Policy 执行策略逻辑:                       │
│                                                                         │
│  ┌─────────────────────────┐                                            │
│  │    Specification        │    IsSatisfiedBy(order) ?                  │
│  │  ┌─────────────────┐    │    - VIP customer?                         │
│  │  │ Amount > $100   │    │    - Order amount > threshold?             │
│  │  │ AND             │───►│    - Specific products?                    │
│  │  │ IsVIP(customer) │    │                                            │
│  │  └─────────────────┘    │                                            │
│  └─────────────────────────┘                                            │
│            │                                                            │
│            │ yes                                                        │
│            ▼                                                            │
│  ┌─────────────────────────┐                                            │
│  │    Policy               │    Execute(order)                          │
│  │  ┌─────────────────┐    │    - Calculate discount                    │
│  │  │ Apply 20% off   │───►│    - Apply rules                           │
│  │  └─────────────────┘    │    - Return result                         │
│  └─────────────────────────┘                                            │
│                                                                         │
│  组合使用:                                                               │
│  spec := NewAndSpec(AmountSpec(100), VIPCustomerSpec())                  │
│  policy := NewConditionalPolicy("BigSpenderVIP", spec, DiscountPolicy(0.25))│
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 5. 决策标准

### 5.1 何时使用 Policy 模式

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Policy Pattern Decision Tree                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  有多种算法实现同一行为？ ─────是────► 使用 Policy                         │
│       │                                                                 │
│       否                                                                │
│       │                                                                 │
│       ▼                                                                 │
│  算法需要运行时切换？ ────────是────► 使用 Policy                         │
│       │                                                                 │
│       否                                                                │
│       │                                                                 │
│       ▼                                                                 │
│  算法需要独立测试？ ─────────是────► 使用 Policy                          │
│       │                                                                 │
│       否                                                                │
│       │                                                                 │
│       ▼                                                                 │
│  使用简单函数或方法                                                       │
│                                                                         │
│  常见 Policy 应用场景:                                                   │
│  • 定价策略（折扣、促销）                                                 │
│  • 审批流程                                                              │
│  • 路由规则                                                              │
│  • 验证规则                                                              │
│  • 评分算法                                                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Policy Implementation Checklist                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 定义策略接口                                                          │
│  □ 识别策略变体                                                          │
│  □ 确定策略选择机制（配置、运行时、规则引擎）                               │
│  □ 考虑策略组合需求                                                        │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现具体策略                                                          │
│  □ 实现策略注册表（可选）                                                  │
│  □ 实现策略组合（如需要）                                                  │
│  □ 实现上下文使用策略                                                      │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 避免策略过于细粒度（过度设计）                                         │
│  ❌ 避免策略相互依赖                                                       │
│  ❌ 避免策略执行副作用                                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>17KB, 完整形式化 + Go 实现 + 测试)

**相关文档**:

- [EC-016-CQRS-Pattern-Formal.md](./EC-016-CQRS-Pattern-Formal.md)
- [EC-040-Aggregate-Pattern.md](./EC-040-Aggregate-Pattern.md)
