# 领域模型设计

> **简介**: 本文档介绍本项目的领域模型设计，包括实体、仓储接口、领域服务和领域错误的定义。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [领域模型设计](#领域模型设计)
  - [📋 目录](#-目录)
  - [1. 📚 概述](#1--概述)
    - [1.1 为什么需要领域模型？](#11-为什么需要领域模型)
    - [1.2 DDD 核心概念](#12-ddd-核心概念)
  - [2. 👤 用户领域（User Domain）](#2--用户领域user-domain)
    - [2.1 实体（Entity）](#21-实体entity)
      - [2.1.1 实体定义](#211-实体定义)
      - [2.1.2 业务规则封装](#212-业务规则封装)
      - [2.1.3 实体设计原则](#213-实体设计原则)
    - [2.2 值对象（Value Object）](#22-值对象value-object)
      - [2.2.1 值对象定义](#221-值对象定义)
      - [2.2.2 为什么需要值对象？](#222-为什么需要值对象)
      - [2.2.3 值对象设计原则](#223-值对象设计原则)
      - [2.2.4 实际应用示例](#224-实际应用示例)
    - [2.3 聚合（Aggregate）](#23-聚合aggregate)
      - [2.3.1 聚合定义](#231-聚合定义)
      - [2.3.2 为什么需要聚合？](#232-为什么需要聚合)
      - [2.3.3 聚合设计原则](#233-聚合设计原则)
      - [2.3.4 实际应用示例](#234-实际应用示例)
    - [2.4 领域事件（Domain Event）](#24-领域事件domain-event)
      - [2.4.1 领域事件定义](#241-领域事件定义)
      - [2.4.2 为什么需要领域事件？](#242-为什么需要领域事件)
      - [2.4.3 领域事件设计原则](#243-领域事件设计原则)
      - [2.4.4 实际应用示例](#244-实际应用示例)
    - [2.5 仓储接口（Repository Interface）](#25-仓储接口repository-interface)
      - [2.2.1 接口定义](#221-接口定义)
      - [2.5.2 为什么接口定义在领域层？](#252-为什么接口定义在领域层)
      - [2.5.3 接口设计原则](#253-接口设计原则)
    - [2.6 领域服务（Domain Service）](#26-领域服务domain-service)
      - [2.6.1 领域服务定义](#261-领域服务定义)
      - [2.6.2 领域服务 vs 应用服务](#262-领域服务-vs-应用服务)
    - [2.7 领域错误（Domain Errors）](#27-领域错误domain-errors)
      - [2.7.1 错误定义](#271-错误定义)
      - [2.7.2 为什么错误定义在领域层？](#272-为什么错误定义在领域层)
      - [2.7.3 错误处理模式](#273-错误处理模式)
  - [3. 🎯 设计原则与论证](#3--设计原则与论证)
    - [3.1 实体独立性原则](#31-实体独立性原则)
      - [3.1.1 原则说明](#311-原则说明)
      - [3.1.2 实际应用](#312-实际应用)
    - [3.2 接口定义原则](#32-接口定义原则)
      - [3.2.1 原则说明](#321-原则说明)
      - [3.2.2 实际应用](#322-实际应用)
    - [3.3 业务规则封装原则](#33-业务规则封装原则)
      - [3.3.1 原则说明](#331-原则说明)
      - [3.3.2 实际应用](#332-实际应用)
    - [3.4 错误处理原则](#34-错误处理原则)
      - [3.4.1 原则说明](#341-原则说明)
      - [3.4.2 实际应用](#342-实际应用)
  - [4. 📚 扩展阅读](#4--扩展阅读)
    - [架构相关](#架构相关)
    - [外部资源](#外部资源)

---

## 1. 📚 概述

领域模型是 Clean Architecture 的核心，包含业务实体、业务规则和领域逻辑。本项目的领域模型遵循 Domain-Driven Design (DDD) 原则。

### 1.1 为什么需要领域模型？

**论证**:

1. **业务表达**:
   - 代码直接表达业务概念
   - 业务专家可以理解代码
   - 减少业务和技术之间的鸿沟

2. **业务规则集中**:
   - 业务规则封装在领域模型中
   - 业务逻辑集中管理
   - 避免业务规则分散

3. **可维护性**:
   - 业务逻辑稳定，技术实现可以变化
   - 业务变更影响范围小
   - 便于长期维护

### 1.2 DDD 核心概念

**实体 (Entity)**:

- 有唯一标识的对象
- 状态可以变化
- 通过标识判断相等性

**值对象 (Value Object)**:

- 没有唯一标识的对象
- 状态不可变
- 通过值判断相等性
- 可以包含业务规则

**聚合 (Aggregate)**:

- 一组相关对象的集合
- 有聚合根（Aggregate Root）
- 保证一致性边界
- 通过聚合根访问内部对象

**领域事件 (Domain Event)**:

- 表示领域中的重要事件
- 不可变
- 包含事件发生的时间和上下文

**仓储 (Repository)**:

- 抽象数据访问
- 提供领域对象持久化接口
- 隐藏数据访问细节

---

## 2. 👤 用户领域（User Domain）

### 2.1 实体（Entity）

#### 2.1.1 实体定义

```go
// internal/domain/user/entity.go
package user

import (
    "errors"
    "time"
)

// User 领域实体
type User struct {
    ID        string    // 唯一标识
    Email     string    // 邮箱（可以提取为 Email 值对象）
    Name      string    // 名称
    CreatedAt time.Time // 创建时间
    UpdatedAt time.Time // 更新时间
}
```

#### 2.1.2 业务规则封装

**为什么业务规则要封装在实体中？**

**论证**:

- **业务逻辑集中**: 业务规则集中在实体中，便于管理
- **避免贫血模型**: 实体包含行为，不仅仅是数据
- **可测试性**: 可以独立测试业务规则

**实际应用**:

```go
// 业务规则：验证邮箱格式
func (u *User) ValidateEmail() error {
    if !isValidEmail(u.Email) {
        return ErrInvalidEmail
    }
    return nil
}

// 业务规则：验证用户名称
func (u *User) ValidateName() error {
    if len(u.Name) < 2 {
        return ErrInvalidName
    }
    if len(u.Name) > 50 {
        return ErrInvalidName
    }
    return nil
}

// 业务规则：验证用户完整性
func (u *User) Validate() error {
    if err := u.ValidateEmail(); err != nil {
        return err
    }
    if err := u.ValidateName(); err != nil {
        return err
    }
    return nil
}
```

#### 2.1.3 实体设计原则

**唯一标识原则**:

- 每个实体必须有唯一标识
- 通过标识判断相等性
- 标识不可变

**业务规则封装原则**:

- 业务规则封装在实体中
- 实体负责自己的有效性
- 避免贫血模型

**不可变值对象原则**:

- 值对象应该是不可变的
- 通过值判断相等性
- 可以提取为独立的值对象

---

### 2.2 值对象（Value Object）

#### 2.2.1 值对象定义

**什么是值对象？**

值对象是没有唯一标识的对象，通过其属性值来判断相等性。值对象是不可变的，一旦创建就不能修改。

**值对象特征**:

- **无唯一标识**: 没有 ID，通过值判断相等性
- **不可变**: 创建后不能修改
- **自包含**: 包含验证逻辑和业务规则
- **可替换**: 可以整体替换

**值对象示例**:

```go
// internal/domain/user/value_object.go
package user

import (
    "errors"
    "regexp"
)

// Email 邮箱值对象
type Email struct {
    value string
}

// NewEmail 创建邮箱值对象
func NewEmail(value string) (Email, error) {
    if !isValidEmail(value) {
        return Email{}, ErrInvalidEmail
    }
    return Email{value: value}, nil
}

// Value 获取邮箱值
func (e Email) Value() string {
    return e.value
}

// Equals 判断相等性
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}
```

#### 2.2.2 为什么需要值对象？

**论证**:

1. **类型安全**:
   - 使用值对象可以避免原始类型（如 string）的误用
   - 编译时检查类型错误
   - 提高代码可读性

2. **业务规则封装**:
   - 值对象可以封装验证逻辑
   - 确保值的有效性
   - 集中管理业务规则

3. **可维护性**:
   - 值对象的逻辑集中在一个地方
   - 修改业务规则时只需修改值对象
   - 减少代码重复

4. **语义清晰**:
   - 值对象名称表达业务概念
   - 代码更易理解
   - 减少注释需求

#### 2.2.3 值对象设计原则

**不可变性原则**:

- 值对象创建后不能修改
- 需要修改时创建新对象
- 避免副作用

**自包含原则**:

- 值对象包含验证逻辑
- 值对象包含业务规则
- 值对象是自包含的

**相等性判断原则**:

- 通过值判断相等性
- 所有属性都相等才相等
- 实现 `Equals` 方法

#### 2.2.4 实际应用示例

**Email 值对象**:

```go
// Email 值对象：封装邮箱验证逻辑
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !isValidEmail(value) {
        return Email{}, ErrInvalidEmail
    }
    return Email{value: value}, nil
}

func (e Email) Value() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}
```

**Money 值对象**:

```go
// Money 值对象：封装金额和货币
type Money struct {
    amount   int64  // 金额（以分为单位）
    currency string // 货币代码（如 "USD", "CNY"）
}

func NewMoney(amount int64, currency string) (Money, error) {
    if amount < 0 {
        return Money{}, errors.New("amount cannot be negative")
    }
    if !isValidCurrency(currency) {
        return Money{}, errors.New("invalid currency")
    }
    return Money{amount: amount, currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    return NewMoney(m.amount+other.amount, m.currency)
}

func (m Money) Equals(other Money) bool {
    return m.amount == other.amount && m.currency == other.currency
}
```

**Address 值对象**:

```go
// Address 值对象：封装地址信息
type Address struct {
    street   string
    city     string
    state    string
    zipCode  string
    country  string
}

func NewAddress(street, city, state, zipCode, country string) (Address, error) {
    if street == "" || city == "" || country == "" {
        return Address{}, errors.New("address fields cannot be empty")
    }
    return Address{
        street:  street,
        city:    city,
        state:   state,
        zipCode: zipCode,
        country: country,
    }, nil
}

func (a Address) Equals(other Address) bool {
    return a.street == other.street &&
        a.city == other.city &&
        a.state == other.state &&
        a.zipCode == other.zipCode &&
        a.country == other.country
}
```

**在实体中使用值对象**:

```go
// 在 User 实体中使用 Email 值对象
type User struct {
    ID        string
    Email     Email  // 使用值对象而不是 string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 创建用户时使用值对象
func NewUser(email Email, name string) (*User, error) {
    // Email 值对象已经验证过，不需要再次验证
    return &User{
        ID:        generateID(),
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

---

### 2.3 聚合（Aggregate）

#### 2.3.1 聚合定义

**什么是聚合？**

聚合是一组相关对象的集合，有一个聚合根（Aggregate Root）。聚合保证一致性边界，外部只能通过聚合根访问聚合内的对象。

**聚合特征**:

- **聚合根**: 聚合有一个唯一的聚合根
- **一致性边界**: 聚合保证内部对象的一致性
- **事务边界**: 聚合是事务的边界
- **不可变引用**: 聚合内的对象引用不可变

**聚合示例**:

```go
// internal/domain/order/aggregate.go
package order

import (
    "time"
    "github.com/yourusername/golang/internal/domain/user"
)

// Order 订单聚合根
type Order struct {
    ID          string
    UserID      string
    Items       []OrderItem  // 聚合内的实体
    Status      OrderStatus
    TotalAmount Money        // 值对象
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// OrderItem 订单项（聚合内的实体）
type OrderItem struct {
    ProductID string
    Quantity  int
    Price     Money
}

// AddItem 添加订单项（通过聚合根操作）
func (o *Order) AddItem(productID string, quantity int, price Money) error {
    if o.Status != OrderStatusDraft {
        return ErrOrderCannotBeModified
    }

    item := OrderItem{
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
    }

    o.Items = append(o.Items, item)
    o.TotalAmount = o.calculateTotal()
    o.UpdatedAt = time.Now()

    return nil
}

// calculateTotal 计算总金额（聚合内部方法）
func (o *Order) calculateTotal() Money {
    total := Money{amount: 0, currency: "CNY"}
    for _, item := range o.Items {
        itemTotal := Money{
            amount:   item.Price.amount * int64(item.Quantity),
            currency: item.Price.currency,
        }
        total, _ = total.Add(itemTotal)
    }
    return total
}
```

#### 2.3.2 为什么需要聚合？

**论证**:

1. **一致性保证**:
   - 聚合保证内部对象的一致性
   - 通过聚合根控制修改
   - 避免数据不一致

2. **事务边界**:
   - 聚合是事务的边界
   - 一个事务只能修改一个聚合
   - 简化事务管理

3. **封装性**:
   - 聚合封装内部实现
   - 外部只能通过聚合根访问
   - 保护内部对象

4. **性能优化**:
   - 聚合可以控制加载范围
   - 避免加载不必要的数据
   - 提升性能

#### 2.3.3 聚合设计原则

**聚合根原则**:

- 每个聚合只有一个聚合根
- 外部只能通过聚合根访问聚合
- 聚合根负责维护一致性

**一致性边界原则**:

- 聚合保证内部对象的一致性
- 跨聚合的一致性通过最终一致性保证
- 避免跨聚合的直接引用

**事务边界原则**:

- 一个事务只能修改一个聚合
- 跨聚合操作使用领域事件
- 避免分布式事务

#### 2.3.4 实际应用示例

**Order 聚合**:

```go
// Order 聚合：订单和订单项
type Order struct {
    ID          string
    UserID      string
    Items       []OrderItem
    Status      OrderStatus
    TotalAmount Money
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 通过聚合根添加订单项
func (o *Order) AddItem(productID string, quantity int, price Money) error {
    // 业务规则：只有草稿状态的订单可以修改
    if o.Status != OrderStatusDraft {
        return ErrOrderCannotBeModified
    }

    // 业务规则：数量必须大于 0
    if quantity <= 0 {
        return ErrInvalidQuantity
    }

    item := OrderItem{
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
    }

    o.Items = append(o.Items, item)
    o.TotalAmount = o.calculateTotal()
    o.UpdatedAt = time.Now()

    return nil
}

// 通过聚合根提交订单
func (o *Order) Submit() error {
    if o.Status != OrderStatusDraft {
        return ErrOrderCannotBeSubmitted
    }

    if len(o.Items) == 0 {
        return ErrOrderHasNoItems
    }

    o.Status = OrderStatusSubmitted
    o.UpdatedAt = time.Now()

    return nil
}
```

---

### 2.4 领域事件（Domain Event）

#### 2.4.1 领域事件定义

**什么是领域事件？**

领域事件表示领域中发生的重要事件。领域事件是不可变的，包含事件发生的时间和上下文信息。

**领域事件特征**:

- **不可变**: 事件一旦创建就不能修改
- **时间戳**: 包含事件发生的时间
- **上下文**: 包含事件发生的上下文信息
- **业务含义**: 表达业务概念

**领域事件示例**:

```go
// internal/domain/user/event.go
package user

import "time"

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
    UserID    string
    Email     string
    Name      string
    CreatedAt time.Time
}

// UserUpdatedEvent 用户更新事件
type UserUpdatedEvent struct {
    UserID    string
    Email     string
    Name      string
    UpdatedAt time.Time
}

// UserDeletedEvent 用户删除事件
type UserDeletedEvent struct {
    UserID    string
    DeletedAt time.Time
}
```

#### 2.4.2 为什么需要领域事件？

**论证**:

1. **解耦**:
   - 领域事件解耦不同聚合
   - 避免直接依赖
   - 提高系统灵活性

2. **最终一致性**:
   - 通过领域事件实现最终一致性
   - 避免分布式事务
   - 提高系统性能

3. **可追溯性**:
   - 领域事件记录重要业务事件
   - 便于审计和追溯
   - 支持事件溯源

4. **扩展性**:
   - 通过领域事件扩展功能
   - 新功能可以订阅事件
   - 不影响现有代码

#### 2.4.3 领域事件设计原则

**不可变性原则**:

- 领域事件创建后不能修改
- 事件是历史记录
- 保持事件的真实性

**业务含义原则**:

- 领域事件表达业务概念
- 事件名称清晰明确
- 事件包含必要的上下文

**时间戳原则**:

- 领域事件包含时间戳
- 记录事件发生的时间
- 支持事件排序

#### 2.4.4 实际应用示例

**在实体中发布事件**:

```go
// User 实体：发布领域事件
type User struct {
    ID        string
    Email     Email
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time

    // 领域事件列表
    events []DomainEvent
}

// DomainEvent 领域事件接口
type DomainEvent interface {
    EventType() string
    OccurredAt() time.Time
}

// NewUser 创建用户并发布事件
func NewUser(email Email, name string) (*User, error) {
    user := &User{
        ID:        generateID(),
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        events:    []DomainEvent{},
    }

    // 发布用户创建事件
    user.addEvent(UserCreatedEvent{
        UserID:    user.ID,
        Email:     user.Email.Value(),
        Name:      user.Name,
        CreatedAt: user.CreatedAt,
    })

    return user, nil
}

// addEvent 添加领域事件
func (u *User) addEvent(event DomainEvent) {
    u.events = append(u.events, event)
}

// GetEvents 获取领域事件
func (u *User) GetEvents() []DomainEvent {
    return u.events
}

// ClearEvents 清除领域事件（发布后清除）
func (u *User) ClearEvents() {
    u.events = []DomainEvent{}
}
```

**在应用层处理事件**:

```go
// 在应用层发布领域事件
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
    email, err := user.NewEmail(req.Email)
    if err != nil {
        return nil, err
    }

    domainUser, err := user.NewUser(email, req.Name)
    if err != nil {
        return nil, err
    }

    // 保存用户
    if err := s.repo.Create(ctx, domainUser); err != nil {
        return nil, err
    }

    // 发布领域事件
    events := domainUser.GetEvents()
    for _, event := range events {
        s.eventPublisher.Publish(ctx, event)
    }

    // 清除事件
    domainUser.ClearEvents()

    return toUserDTO(domainUser), nil
}
```

---

### 2.5 仓储接口（Repository Interface）

#### 2.2.1 接口定义

```go
// internal/domain/user/repository.go
package user

import "context"

// Repository 仓储接口：定义在领域层
type Repository interface {
    // 创建用户
    Create(ctx context.Context, user *User) error

    // 根据 ID 查找用户
    FindByID(ctx context.Context, id string) (*User, error)

    // 根据邮箱查找用户
    FindByEmail(ctx context.Context, email string) (*User, error)

    // 更新用户
    Update(ctx context.Context, user *User) error

    // 删除用户
    Delete(ctx context.Context, id string) error

    // 列出用户
    List(ctx context.Context, limit, offset int) ([]*User, error)
}
```

#### 2.5.2 为什么接口定义在领域层？

**论证**:

1. **业务需求驱动**:
   - 接口表达业务需求
   - 不包含技术细节
   - 稳定，不轻易变化

2. **依赖倒置**:
   - 应用层依赖领域层接口
   - 基础设施层实现领域层接口
   - 符合依赖倒置原则

3. **可测试性**:
   - 可以 Mock 接口进行测试
   - 不需要真实的数据库
   - 测试速度快

#### 2.5.3 接口设计原则

**接口粒度适中**:

- 不要过于细粒度（增加复杂度）
- 不要过于粗粒度（降低灵活性）
- 平衡灵活性和复杂度

**接口稳定性**:

- 接口应该稳定，不轻易变化
- 变化时考虑向后兼容
- 使用版本控制管理接口变化

### 2.6 领域服务（Domain Service）

#### 2.6.1 领域服务定义

**什么时候需要领域服务？**

**论证**:

- **跨聚合操作**: 操作涉及多个聚合
- **无状态操作**: 操作不涉及状态变化
- **业务规则**: 复杂的业务规则

**实际应用**:

```go
// 领域服务：处理跨聚合的业务规则
type DomainService interface {
    // 验证邮箱格式
    ValidateEmail(email string) bool

    // 检查邮箱是否唯一（需要访问 Repository）
    IsEmailUnique(ctx context.Context, email string) (bool, error)

    // 生成用户 ID
    GenerateUserID() string
}

// 领域服务实现
type domainService struct {
    repo Repository
}

func NewDomainService(repo Repository) DomainService {
    return &domainService{repo: repo}
}

func (s *domainService) ValidateEmail(email string) bool {
    // 邮箱格式验证逻辑
    return isValidEmail(email)
}

func (s *domainService) IsEmailUnique(ctx context.Context, email string) (bool, error) {
    user, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return true, nil  // 用户不存在，邮箱唯一
        }
        return false, err
    }
    return user == nil, nil  // 用户存在，邮箱不唯一
}
```

#### 2.6.2 领域服务 vs 应用服务

**领域服务**:

- 处理领域逻辑
- 可以访问 Repository
- 无状态

**应用服务**:

- 处理用例编排
- 协调领域对象
- 管理事务边界

**区别**:

- 领域服务: 处理领域逻辑，属于 Domain Layer
- 应用服务: 处理用例编排，属于 Application Layer

### 2.7 领域错误（Domain Errors）

#### 2.7.1 错误定义

```go
// internal/domain/user/errors.go
package user

import "errors"

// 领域错误：定义在领域层
var (
    // ErrUserNotFound 用户不存在
    ErrUserNotFound = errors.New("user not found")

    // ErrUserAlreadyExists 用户已存在
    ErrUserAlreadyExists = errors.New("user already exists")

    // ErrInvalidEmail 无效邮箱
    ErrInvalidEmail = errors.New("invalid email")

    // ErrInvalidName 无效名称
    ErrInvalidName = errors.New("invalid name")
)
```

#### 2.7.2 为什么错误定义在领域层？

**论证**:

1. **业务错误表达**:
   - 错误表达业务概念
   - 不包含技术细节
   - 业务专家可以理解

2. **统一错误处理**:
   - 统一的错误定义
   - 便于错误处理
   - 便于错误转换

3. **可测试性**:
   - 可以测试错误场景
   - 错误信息清晰
   - 便于问题排查

#### 2.7.3 错误处理模式

**错误定义**:

- 使用 `errors.New()` 定义错误
- 错误信息清晰
- 错误类型明确

**错误使用**:

- 在领域层返回领域错误
- 在应用层处理领域错误
- 在接口层转换错误格式

---

## 3. 🎯 设计原则与论证

### 3.1 实体独立性原则

#### 3.1.1 原则说明

**定义**: 实体不依赖外部框架，只包含业务逻辑。

**论证**:

1. **业务逻辑独立性**:
   - 业务逻辑是系统的核心价值
   - 业务逻辑应该独立于技术实现
   - 技术实现可以变化，但业务逻辑应该稳定

2. **可测试性**:
   - 不依赖外部框架，可以独立测试
   - 可以使用简单的单元测试
   - 测试速度快，隔离性好

3. **可移植性**:
   - 可以轻松移植到其他技术栈
   - 可以轻松替换技术实现
   - 业务逻辑不受技术限制

#### 3.1.2 实际应用

```go
// ✅ 正确：实体不依赖外部框架
type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// ❌ 错误：实体依赖外部框架
import "entgo.io/ent"

type User struct {
    ent.Schema  // 依赖 Ent 框架
    ID          string
    Email       string
}
```

### 3.2 接口定义原则

#### 3.2.1 原则说明

**定义**: 仓储接口在领域层定义，实现在基础设施层。

**论证**:

1. **业务需求驱动**:
   - 接口表达业务需求
   - 不包含技术细节
   - 稳定，不轻易变化

2. **依赖倒置**:
   - 应用层依赖领域层接口
   - 基础设施层实现领域层接口
   - 符合依赖倒置原则

3. **可测试性**:
   - 可以 Mock 接口进行测试
   - 不需要真实的数据库
   - 测试速度快

#### 3.2.2 实际应用

```go
// Domain Layer: 定义接口
type Repository interface {
    Create(ctx context.Context, user *User) error
}

// Infrastructure Layer: 实现接口
type EntUserRepository struct {
    client *ent.Client
}

func (r *EntUserRepository) Create(ctx context.Context, user *user.User) error {
    // 实现接口
}
```

### 3.3 业务规则封装原则

#### 3.3.1 原则说明

**定义**: 业务规则封装在实体和领域服务中，确保业务逻辑集中管理。

**论证**:

1. **业务逻辑集中**:
   - 业务规则集中在领域模型中
   - 避免业务规则分散
   - 便于管理和维护

2. **避免贫血模型**:
   - 实体包含行为，不仅仅是数据
   - 业务规则封装在实体中
   - 符合面向对象设计

3. **可测试性**:
   - 可以独立测试业务规则
   - 测试范围明确
   - 测试用例清晰

#### 3.3.2 实际应用

```go
// ✅ 正确：业务规则封装在实体中
type User struct {
    Email string
}

func (u *User) ValidateEmail() error {
    if !isValidEmail(u.Email) {
        return ErrInvalidEmail
    }
    return nil
}

// ❌ 错误：业务规则分散在应用层
type User struct {
    Email string
}

// 在应用层验证
func CreateUser(email string) error {
    if !isValidEmail(email) {  // 业务规则在应用层
        return ErrInvalidEmail
    }
    // ...
}
```

### 3.4 错误处理原则

#### 3.4.1 原则说明

**定义**: 领域错误在领域层定义，提供统一的错误处理机制。

**论证**:

1. **业务错误表达**:
   - 错误表达业务概念
   - 不包含技术细节
   - 业务专家可以理解

2. **统一错误处理**:
   - 统一的错误定义
   - 便于错误处理
   - 便于错误转换

3. **可测试性**:
   - 可以测试错误场景
   - 错误信息清晰
   - 便于问题排查

#### 3.4.2 实际应用

```go
// Domain Layer: 定义领域错误
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
)

// Application Layer: 使用领域错误
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) error {
    existing, err := s.repo.FindByEmail(ctx, req.Email)
    if err == nil && existing != nil {
        return user.ErrUserAlreadyExists  // 返回领域错误
    }
    // ...
}

// Interfaces Layer: 转换错误格式
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    err := h.service.CreateUser(ctx, req)
    if err != nil {
        if errors.Is(err, user.ErrUserAlreadyExists) {
            Error(w, http.StatusConflict, err)  // 转换为 HTTP 错误
        }
    }
}
```

---

## 4. 📚 扩展阅读

### 架构相关

- [Clean Architecture](./clean-architecture.md) - 架构设计详解
- [工作流架构设计](./workflow.md) - 工作流集成
- [架构知识图谱](./00-知识图谱.md) - 架构知识图谱
- [概念定义体系](./00-概念定义体系.md) - 概念定义
- [第三方技术栈深度解析](./00-第三方技术栈深度解析.md) - 第三方技术栈解析

### 外部资源

- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html) - DDD 概念
- [DDD 模式](https://www.domainlanguage.com/ddd/) - DDD 模式详解

---

> 📚 **简介**
> 本文深入探讨项目的领域模型设计，系统讲解实体、仓储接口、领域服务和领域错误的定义。通过本文，您将全面掌握领域模型的设计原则和实践方法。
