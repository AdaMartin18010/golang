# EC-042: Entity Pattern (实体模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #entity #identity #ddd #lifecycle
> **权威来源**:
>
> - [Entity](https://martinfowler.com/bliki/EvansClassification.html) - Martin Fowler
> - [Domain-Driven Design](https://domainlanguage.com/ddd/) - Eric Evans
> - [Implementing Domain-Driven Design](https://www.oreilly.com/library/view/implementing-domain-driven-design/9780133039900/) - Vaughn Vernon

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在领域模型中，如何表示具有独立生命周期、概念标识的业务对象，即使属性变化也能保持身份连续性？

**形式化描述**:

```
给定: 业务概念集合 C = {C₁, C₂, ..., Cₙ}
给定: 某些概念具有:
  - 独立生命周期
  - 需要跟踪状态变化历史
  - 多个实例可能有相同属性但代表不同事物

区分:
  Entity: 概念标识决定对象身份
  Value Object: 属性集合决定对象身份
```

**示例**:

```
Customer 是 Entity:
  - 即使更改了姓名、地址，还是同一个 Customer
  - cust-001 永远是 cust-001
  - 需要跟踪其订单历史

Address 是 Value Object:
  - "123 Main St, NYC" 就是 "123 Main St, NYC"
  - 改变内容就是不同的地址
  - 不需要跟踪地址的历史（除非特殊需求）
```

### 1.2 解决方案形式化

**定义 1.1 (实体)**
实体是具有概念标识的领域对象，其身份不依赖于属性：

```
实体 E:
  E = ⟨ID, Attributes, Behavior, Lifecycle⟩

相等性:
  E₁ = E₂ ⟺ E₁.ID = E₂.ID

特性:
  - 连续性: 属性变化时 ID 不变
  - 可变性: 属性可以修改
  - 可追溯: 生命周期可追踪
  - 引用性: 可被其他对象引用
```

**定义 1.2 (实体标识)**

```
标识 ID 的特性:
  - 唯一性: ∀E₁, E₂: E₁.ID ≠ E₂.ID (E₁ ≠ E₂)
  - 不变性: ID 一旦分配永不改变
  - 无业务含义: ID 不应包含业务信息（如 cust-001 ✗，UUID ✓）
  - 技术生成: 通常由系统生成而非用户输入
```

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Entity Lifecycle                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐   Create    ┌─────────┐   Modify    ┌─────────┐           │
│  │  None   │────────────►│ Active  │────────────►│ Active  │           │
│  └─────────┘             └────┬────┘             └────┬────┘           │
│                               │                        │               │
│                               │ Archive                │ Disable       │
│                               ▼                        ▼               │
│                         ┌─────────┐              ┌─────────┐           │
│                         │Archived │              │Disabled │           │
│                         └────┬────┘              └────┬────┘           │
│                              │                        │                │
│                              │ Delete                 │ Delete         │
│                              ▼                        ▼                │
│                         ┌─────────┐              ┌─────────┐           │
│                         │ Deleted │              │ Deleted │           │
│                         └─────────┘              └─────────┘           │
│                                                                         │
│  Entity Identity 贯穿整个生命周期:                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Customer ID: 550e8400-e29b-41d4-a716-446655440000              │   │
│  │                                                                 │   │
│  │  Version 1: Name="John", Email="john@old.com"                   │   │
│  │  Version 2: Name="John", Email="john@new.com"  ← 修改属性，ID不变 │   │
│  │  Version 3: Status="Disabled"                   ← 状态变更，ID不变 │   │
│  │  Version 4: Status="Archived"                   ← 归档，ID仍然不变 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  关键特性:                                                               │
│  • 修改属性不改变身份                                                    │
│  • 状态转换不改变身份                                                    │
│  • 可与其他实体建立关联                                                   │
│  • 可被多个聚合引用（通过ID）                                             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心实体实现

```go
// entity/core.go
package entity

import (
    "context"
    "fmt"
    "reflect"
    "time"

    "github.com/google/uuid"
)

// Entity 实体接口
type Entity interface {
    // Identity 获取实体标识
    Identity() Identity

    // Equals 基于ID的相等比较
    Equals(other Entity) bool

    // SameIdentityAs 是否是同一实体（可能状态不同）
    SameIdentityAs(other Entity) bool

    // Version 获取版本号（乐观锁）
    Version() int

    // CreatedAt 创建时间
    CreatedAt() time.Time

    // UpdatedAt 更新时间
    UpdatedAt() time.Time
}

// Identity 标识接口
type Identity interface {
    // ID 返回标识值
    ID() string

    // String 字符串表示
    String() string

    // Empty 是否为空标识
    Empty() bool
}

// UUIDIdentity UUID 标识实现
type UUIDIdentity struct {
    value string
}

// NewUUIDIdentity 创建新 UUID 标识
func NewUUIDIdentity() UUIDIdentity {
    return UUIDIdentity{value: uuid.New().String()}
}

// ParseUUIDIdentity 从字符串解析
func ParseUUIDIdentity(id string) (UUIDIdentity, error) {
    if id == "" {
        return UUIDIdentity{}, fmt.Errorf("id cannot be empty")
    }
    // 可以添加 UUID 格式验证
    return UUIDIdentity{value: id}, nil
}

func (i UUIDIdentity) ID() string    { return i.value }
func (i UUIDIdentity) String() string { return i.value }
func (i UUIDIdentity) Empty() bool    { return i.value == "" }

// Equals 标识相等
func (i UUIDIdentity) Equals(other Identity) bool {
    if other == nil {
        return false
    }
    return i.value == other.ID()
}

// EntityBase 实体基础
type EntityBase struct {
    id        UUIDIdentity
    version   int
    createdAt time.Time
    updatedAt time.Time
}

// NewEntityBase 创建实体基础
func NewEntityBase() EntityBase {
    now := time.Now()
    return EntityBase{
        id:        NewUUIDIdentity(),
        version:   1,
        createdAt: now,
        updatedAt: now,
    }
}

// NewEntityBaseWithID 使用指定ID创建
func NewEntityBaseWithID(id string) (EntityBase, error) {
    identity, err := ParseUUIDIdentity(id)
    if err != nil {
        return EntityBase{}, err
    }
    now := time.Now()
    return EntityBase{
        id:        identity,
        version:   1,
        createdAt: now,
        updatedAt: now,
    }, nil
}

func (e *EntityBase) Identity() Identity   { return e.id }
func (e *EntityBase) Version() int         { return e.version }
func (e *EntityBase) CreatedAt() time.Time { return e.createdAt }
func (e *EntityBase) UpdatedAt() time.Time { return e.updatedAt }

// Equals 基于ID的相等
func (e *EntityBase) Equals(other Entity) bool {
    if other == nil {
        return false
    }
    return e.id.Equals(other.Identity())
}

// SameIdentityAs 同一实体判断
func (e *EntityBase) SameIdentityAs(other Entity) bool {
    return e.Equals(other)
}

// MarkModified 标记为已修改
func (e *EntityBase) MarkModified() {
    e.updatedAt = time.Now()
    e.version++
}

// Repository 实体仓储接口
type Repository interface {
    // FindByID 根据ID查找
    FindByID(ctx context.Context, id Identity) (Entity, error)

    // FindAll 查找所有
    FindAll(ctx context.Context) ([]Entity, error)

    // Save 保存实体
    Save(ctx context.Context, entity Entity) error

    // Delete 删除实体
    Delete(ctx context.Context, id Identity) error

    // Exists 判断是否存在
    Exists(ctx context.Context, id Identity) (bool, error)
}

// Specification 规格模式接口
type Specification interface {
    IsSatisfiedBy(entity Entity) bool
    And(other Specification) Specification
    Or(other Specification) Specification
    Not() Specification
}

// BaseSpecification 基础规格
type BaseSpecification struct {
    predicate func(Entity) bool
}

// IsSatisfiedBy 判断是否满足
func (s BaseSpecification) IsSatisfiedBy(entity Entity) bool {
    if s.predicate == nil {
        return true
    }
    return s.predicate(entity)
}

// And 与操作
func (s BaseSpecification) And(other Specification) Specification {
    return BaseSpecification{
        predicate: func(e Entity) bool {
            return s.IsSatisfiedBy(e) && other.IsSatisfiedBy(e)
        },
    }
}

// Or 或操作
func (s BaseSpecification) Or(other Specification) Specification {
    return BaseSpecification{
        predicate: func(e Entity) bool {
            return s.IsSatisfiedBy(e) || other.IsSatisfiedBy(e)
        },
    }
}

// Not 非操作
func (s BaseSpecification) Not() Specification {
    return BaseSpecification{
        predicate: func(e Entity) bool {
            return !s.IsSatisfiedBy(e)
        },
    }
}
```

### 2.2 客户实体实现

```go
// entity/customer.go
package entity

import (
    "errors"
    "fmt"
    "strings"
    "time"
)

// Customer 客户实体
type Customer struct {
    *EntityBase
    name      string
    email     string
    phone     string
    status    CustomerStatus
    addresses []Address
    preferences CustomerPreferences
}

// CustomerStatus 客户状态
type CustomerStatus int

const (
    CustomerStatusActive CustomerStatus = iota
    CustomerStatusInactive
    CustomerStatusSuspended
    CustomerStatusDeleted
)

func (s CustomerStatus) String() string {
    names := []string{"Active", "Inactive", "Suspended", "Deleted"}
    if int(s) < len(names) {
        return names[s]
    }
    return "Unknown"
}

// Address 地址值对象
type Address struct {
    Street  string
    City    string
    State   string
    ZipCode string
    Country string
    Type    string // billing, shipping
}

// CustomerPreferences 客户偏好值对象
type CustomerPreferences struct {
    Newsletter bool
    Currency   string
    Language   string
}

// NewCustomer 创建新客户
func NewCustomer(name, email string) (*Customer, error) {
    customer := &Customer{
        EntityBase:  &NewEntityBase(),
        name:        strings.TrimSpace(name),
        email:       strings.TrimSpace(strings.ToLower(email)),
        status:      CustomerStatusActive,
        addresses:   make([]Address, 0),
        preferences: CustomerPreferences{Currency: "USD", Language: "en"},
    }

    if err := customer.Validate(); err != nil {
        return nil, err
    }

    return customer, nil
}

// NewCustomerWithID 使用指定ID创建客户（用于重建）
func NewCustomerWithID(id string, name, email string) (*Customer, error) {
    base, err := NewEntityBaseWithID(id)
    if err != nil {
        return nil, err
    }

    customer := &Customer{
        EntityBase:  &base,
        name:        strings.TrimSpace(name),
        email:       strings.TrimSpace(strings.ToLower(email)),
        status:      CustomerStatusActive,
        addresses:   make([]Address, 0),
        preferences: CustomerPreferences{Currency: "USD", Language: "en"},
    }

    if err := customer.Validate(); err != nil {
        return nil, err
    }

    return customer, nil
}

// Name 获取姓名
func (c *Customer) Name() string { return c.name }

// Email 获取邮箱
func (c *Customer) Email() string { return c.email }

// Phone 获取电话
func (c *Customer) Phone() string { return c.phone }

// Status 获取状态
func (c *Customer) Status() CustomerStatus { return c.status }

// Addresses 获取地址
func (c *Customer) Addresses() []Address {
    result := make([]Address, len(c.addresses))
    copy(result, c.addresses)
    return result
}

// Preferences 获取偏好
func (c *Customer) Preferences() CustomerPreferences {
    return c.preferences
}

// ChangeName 更改姓名（身份不变）
func (c *Customer) ChangeName(newName string) error {
    trimmed := strings.TrimSpace(newName)
    if trimmed == "" {
        return errors.New("name cannot be empty")
    }

    c.name = trimmed
    c.MarkModified()
    return nil
}

// ChangeEmail 更改邮箱（身份不变）
func (c *Customer) ChangeEmail(newEmail string) error {
    trimmed := strings.TrimSpace(strings.ToLower(newEmail))
    if trimmed == "" {
        return errors.New("email cannot be empty")
    }
    if !strings.Contains(trimmed, "@") {
        return errors.New("invalid email format")
    }

    c.email = trimmed
    c.MarkModified()
    return nil
}

// UpdatePhone 更新电话（身份不变）
func (c *Customer) UpdatePhone(phone string) {
    c.phone = strings.TrimSpace(phone)
    c.MarkModified()
}

// AddAddress 添加地址
func (c *Customer) AddAddress(address Address) {
    c.addresses = append(c.addresses, address)
    c.MarkModified()
}

// RemoveAddress 移除地址
func (c *Customer) RemoveAddress(index int) error {
    if index < 0 || index >= len(c.addresses) {
        return errors.New("invalid address index")
    }

    c.addresses = append(c.addresses[:index], c.addresses[index+1:]...)
    c.MarkModified()
    return nil
}

// UpdatePreferences 更新偏好
func (c *Customer) UpdatePreferences(prefs CustomerPreferences) {
    c.preferences = prefs
    c.MarkModified()
}

// Suspend 暂停账户
func (c *Customer) Suspend(reason string) error {
    if c.status == CustomerStatusDeleted {
        return errors.New("cannot suspend deleted customer")
    }
    c.status = CustomerStatusSuspended
    c.MarkModified()
    return nil
}

// Reactivate 重新激活
func (c *Customer) Reactivate() error {
    if c.status == CustomerStatusDeleted {
        return errors.New("cannot reactivate deleted customer")
    }
    c.status = CustomerStatusActive
    c.MarkModified()
    return nil
}

// Delete 删除（软删除）
func (c *Customer) Delete() error {
    c.status = CustomerStatusDeleted
    c.MarkModified()
    return nil
}

// IsActive 是否活跃
func (c *Customer) IsActive() bool {
    return c.status == CustomerStatusActive
}

// Validate 验证
func (c *Customer) Validate() error {
    if c.name == "" {
        return errors.New("customer name is required")
    }
    if c.email == "" {
        return errors.New("customer email is required")
    }
    return nil
}

// String 字符串表示
func (c *Customer) String() string {
    return fmt.Sprintf("Customer{id=%s, name=%s, email=%s, status=%s}",
        c.Identity().ID(), c.name, c.email, c.status)
}

// CustomerRepository 客户仓储接口
type CustomerRepository interface {
    Repository
    FindByEmail(ctx context.Context, email string) (*Customer, error)
    FindActiveCustomers(ctx context.Context) ([]*Customer, error)
}

// ActiveCustomerSpecification 活跃客户规格
type ActiveCustomerSpecification struct{}

func (s ActiveCustomerSpecification) IsSatisfiedBy(e Entity) bool {
    customer, ok := e.(*Customer)
    if !ok {
        return false
    }
    return customer.IsActive()
}

func (s ActiveCustomerSpecification) And(other Specification) Specification {
    return BaseSpecification{}.And(other)
}

func (s ActiveCustomerSpecification) Or(other Specification) Specification {
    return BaseSpecification{}.Or(other)
}

func (s ActiveCustomerSpecification) Not() Specification {
    return BaseSpecification{}.Not()
}
```

### 2.3 产品实体实现

```go
// entity/product.go
package entity

import (
    "errors"
    "fmt"
    "strings"
    "time"
)

// Product 产品实体
type Product struct {
    *EntityBase
    sku         string
    name        string
    description string
    price       Money
    categoryID  string
    inventory   int
    status      ProductStatus
}

// ProductStatus 产品状态
type ProductStatus int

const (
    ProductStatusDraft ProductStatus = iota
    ProductStatusActive
    ProductStatusDiscontinued
    ProductStatusOutOfStock
)

func (s ProductStatus) String() string {
    names := []string{"Draft", "Active", "Discontinued", "OutOfStock"}
    if int(s) < len(names) {
        return names[s]
    }
    return "Unknown"
}

// Money 金额值对象
type Money struct {
    Amount   float64
    Currency string
}

// NewProduct 创建新产品
func NewProduct(sku, name string, price Money, categoryID string) (*Product, error) {
    product := &Product{
        EntityBase:  &NewEntityBase(),
        sku:         strings.TrimSpace(strings.ToUpper(sku)),
        name:        strings.TrimSpace(name),
        price:       price,
        categoryID:  categoryID,
        inventory:   0,
        status:      ProductStatusDraft,
    }

    if err := product.Validate(); err != nil {
        return nil, err
    }

    return product, nil
}

// SKU 获取SKU
func (p *Product) SKU() string { return p.sku }

// Name 获取名称
func (p *Product) Name() string { return p.name }

// Description 获取描述
func (p *Product) Description() string { return p.description }

// Price 获取价格
func (p *Product) Price() Money { return p.price }

// CategoryID 获取分类ID
func (p *Product) CategoryID() string { return p.categoryID }

// Inventory 获取库存
func (p *Product) Inventory() int { return p.inventory }

// Status 获取状态
func (p *Product) Status() ProductStatus { return p.status }

// IsAvailable 是否可购买
func (p *Product) IsAvailable() bool {
    return p.status == ProductStatusActive && p.inventory > 0
}

// UpdateName 更新名称
func (p *Product) UpdateName(name string) error {
    trimmed := strings.TrimSpace(name)
    if trimmed == "" {
        return errors.New("product name cannot be empty")
    }
    p.name = trimmed
    p.MarkModified()
    return nil
}

// UpdatePrice 更新价格
func (p *Product) UpdatePrice(price Money) error {
    if price.Amount < 0 {
        return errors.New("price cannot be negative")
    }
    p.price = price
    p.MarkModified()
    return nil
}

// UpdateInventory 更新库存
func (p *Product) UpdateInventory(quantity int) error {
    if quantity < 0 {
        return errors.New("inventory cannot be negative")
    }
    p.inventory = quantity

    // 自动更新状态
    if p.status == ProductStatusOutOfStock && quantity > 0 {
        p.status = ProductStatusActive
    } else if quantity == 0 && p.status == ProductStatusActive {
        p.status = ProductStatusOutOfStock
    }

    p.MarkModified()
    return nil
}

// Activate 激活产品
func (p *Product) Activate() error {
    if p.status == ProductStatusDiscontinued {
        return errors.New("cannot activate discontinued product")
    }
    p.status = ProductStatusActive
    p.MarkModified()
    return nil
}

// Discontinue 停产
func (p *Product) Discontinue() {
    p.status = ProductStatusDiscontinued
    p.MarkModified()
}

// Validate 验证
func (p *Product) Validate() error {
    if p.sku == "" {
        return errors.New("SKU is required")
    }
    if p.name == "" {
        return errors.New("product name is required")
    }
    if p.price.Amount < 0 {
        return errors.New("price cannot be negative")
    }
    return nil
}

// String 字符串表示
func (p *Product) String() string {
    return fmt.Sprintf("Product{id=%s, sku=%s, name=%s, price=%.2f %s}",
        p.Identity().ID(), p.sku, p.name, p.price.Amount, p.price.Currency)
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// entity/customer_test.go
package entity

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewCustomer(t *testing.T) {
    customer, err := NewCustomer("John Doe", "john@example.com")

    require.NoError(t, err)
    assert.NotNil(t, customer)
    assert.NotEmpty(t, customer.Identity().ID())
    assert.Equal(t, "John Doe", customer.Name())
    assert.Equal(t, "john@example.com", customer.Email())
    assert.Equal(t, CustomerStatusActive, customer.Status())
    assert.True(t, customer.IsActive())
}

func TestNewCustomer_Validation(t *testing.T) {
    _, err := NewCustomer("", "john@example.com")
    assert.Error(t, err)

    _, err = NewCustomer("John", "")
    assert.Error(t, err)
}

func TestCustomer_ChangeName(t *testing.T) {
    customer, _ := NewCustomer("John Doe", "john@example.com")
    originalVersion := customer.Version()

    err := customer.ChangeName("Jane Doe")

    require.NoError(t, err)
    assert.Equal(t, "Jane Doe", customer.Name())
    assert.Equal(t, originalVersion+1, customer.Version())
    assert.True(t, customer.UpdatedAt().After(customer.CreatedAt()))
}

func TestCustomer_ChangeEmail(t *testing.T) {
    customer, _ := NewCustomer("John Doe", "john@example.com")

    err := customer.ChangeEmail("JOHN@EXAMPLE.COM")

    require.NoError(t, err)
    assert.Equal(t, "john@example.com", customer.Email()) // 转换为小写
}

func TestCustomer_SameIdentity(t *testing.T) {
    customer1, _ := NewCustomerWithID("cust-001", "John", "john@example.com")
    customer2, _ := NewCustomerWithID("cust-001", "Jane", "jane@example.com")
    customer3, _ := NewCustomer("Different", "different@example.com")

    assert.True(t, customer1.SameIdentityAs(customer2))
    assert.False(t, customer1.SameIdentityAs(customer3))

    // 不同属性但相同ID
    assert.Equal(t, "John", customer1.Name())
    assert.Equal(t, "Jane", customer2.Name())
}

func TestCustomer_SuspendAndReactivate(t *testing.T) {
    customer, _ := NewCustomer("John Doe", "john@example.com")

    err := customer.Suspend("violation")
    require.NoError(t, err)
    assert.Equal(t, CustomerStatusSuspended, customer.Status())
    assert.False(t, customer.IsActive())

    err = customer.Reactivate()
    require.NoError(t, err)
    assert.Equal(t, CustomerStatusActive, customer.Status())
}

func TestCustomer_AddressManagement(t *testing.T) {
    customer, _ := NewCustomer("John Doe", "john@example.com")

    address := Address{
        Street:  "123 Main St",
        City:    "NYC",
        State:   "NY",
        ZipCode: "10001",
        Country: "USA",
        Type:    "shipping",
    }

    customer.AddAddress(address)
    assert.Len(t, customer.Addresses(), 1)

    err := customer.RemoveAddress(0)
    require.NoError(t, err)
    assert.Len(t, customer.Addresses(), 0)
}

func TestActiveCustomerSpecification(t *testing.T) {
    activeCustomer, _ := NewCustomer("Active", "active@example.com")
    suspendedCustomer, _ := NewCustomer("Suspended", "suspended@example.com")
    suspendedCustomer.Suspend("test")

    spec := ActiveCustomerSpecification{}

    assert.True(t, spec.IsSatisfiedBy(activeCustomer))
    assert.False(t, spec.IsSatisfiedBy(suspendedCustomer))
    assert.False(t, spec.IsSatisfiedBy(&Product{}))
}

func TestProduct(t *testing.T) {
    t.Run("creation", func(t *testing.T) {
        product, err := NewProduct("SKU-001", "Test Product", Money{Amount: 99.99, Currency: "USD"}, "cat-001")

        require.NoError(t, err)
        assert.Equal(t, "SKU-001", product.SKU())
        assert.Equal(t, "Test Product", product.Name())
        assert.Equal(t, ProductStatusDraft, product.Status())
    })

    t.Run("inventory update", func(t *testing.T) {
        product, _ := NewProduct("SKU-001", "Test", Money{Amount: 10}, "cat")
        product.Activate()

        product.UpdateInventory(100)
        assert.Equal(t, 100, product.Inventory())
        assert.Equal(t, ProductStatusActive, product.Status())
        assert.True(t, product.IsAvailable())
    })

    t.Run("out of stock", func(t *testing.T) {
        product, _ := NewProduct("SKU-001", "Test", Money{Amount: 10}, "cat")
        product.Activate()
        product.UpdateInventory(100)

        product.UpdateInventory(0)
        assert.Equal(t, ProductStatusOutOfStock, product.Status())
        assert.False(t, product.IsAvailable())
    })
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Repository 模式的关系

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Entity + Repository Pattern                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Repository 负责实体的生命周期管理:                                        │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     CustomerRepository                           │   │
│  │                                                                  │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │   FindByID  │    │    Save     │    │      Delete         │  │   │
│  │  │             │    │             │    │                     │  │   │
│  │  │ Load from   │    │ Persist     │    │ Remove from         │  │   │
│  │  │ database    │    │ changes     │    │ database            │  │   │
│  │  │ Reconstruct │    │ Optimistic  │    │ (or soft delete)    │  │   │
│  │  │ entity      │    │ locking     │    │                     │  │   │
│  │  └──────┬──────┘    └──────┬──────┘    └──────────┬──────────┘  │   │
│  │         │                  │                      │             │   │
│  │         └──────────────────┼──────────────────────┘             │   │
│  │                            │                                    │   │
│  │                            ▼                                    │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │                   Customer (Entity)                        │  │   │
│  │  │                                                            │  │   │
│  │  │  - ID: permanent identity                                  │  │   │
│  │  │  - Mutable attributes                                      │  │   │
│  │  │  - Version: optimistic locking                             │  │   │
│  │  │  - Lifecycle tracking                                      │  │   │
│  │  │                                                            │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  关键区别:                                                               │
│  • Entity: 领域概念，有业务含义                                           │
│  • Repository: 技术概念，负责持久化                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 5. 决策标准

### 5.1 何时建模为实体

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Entity Design Decision                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  概念需要跟踪生命周期？ ──────是────► Entity                              │
│       │                                                                 │
│       否                                                                │
│       │                                                                 │
│       ▼                                                                 │
│  概念需要唯一标识？ ─────────是────► Entity                              │
│       │                                                                 │
│       否                                                                │
│       │                                                                 │
│       ▼                                                                 │
│  多个实例可有相同属性？ ──────是────► Entity（需要ID区分）                  │
│       │                                                                 │
│       否                                                                │
│       │                                                                 │
│       ▼                                                                 │
│  考虑使用 Value Object                                                  │
│                                                                         │
│  示例:                                                                  │
│  • Customer, Order, Product, Shipment  → Entity                         │
│  • Money, Address, DateRange, Color    → Value Object                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Entity Implementation Checklist                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 确定实体的唯一标识策略                                                 │
│  □ 定义实体的属性和行为                                                   │
│  □ 定义实体的生命周期状态                                                 │
│  □ 识别与其他实体的关系                                                   │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现 Identity 接口                                                     │
│  □ 实现基于ID的相等比较                                                   │
│  □ 实现乐观锁（版本号）                                                   │
│  □ 实现业务行为和验证                                                     │
│  □ 实现 Repository 接口                                                   │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 不要给实体设置业务相关的ID（如 cust-001）                              │
│  ❌ 不要修改实体的ID                                                       │
│  ❌ 不要在实体中包含过多逻辑（使用领域服务）                                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>17KB, 完整形式化 + Go 实现 + 测试)

**相关文档**:

- [EC-041-Value-Object-Pattern.md](./EC-041-Value-Object-Pattern.md)
- [EC-040-Aggregate-Pattern.md](./EC-040-Aggregate-Pattern.md)
- [EC-043-Repository-Pattern.md](./EC-043-Repository-Pattern.md)
