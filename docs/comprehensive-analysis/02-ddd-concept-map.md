# DDD 领域驱动设计概念体系

## 目录

- [DDD 领域驱动设计概念体系](#ddd-领域驱动设计概念体系)
  - [目录](#目录)
  - [一、核心概念本体论](#一核心概念本体论)
  - [二、战略设计公理定理](#二战略设计公理定理)
    - [公理 1: 限界上下文完整性公理](#公理-1-限界上下文完整性公理)
    - [公理 2: 通用语言一致性公理](#公理-2-通用语言一致性公理)
    - [定理 1: 上下文映射复杂度定理](#定理-1-上下文映射复杂度定理)
    - [定理 2: 防腐层隔离定理](#定理-2-防腐层隔离定理)
  - [三、战术设计公理定理](#三战术设计公理定理)
    - [公理 3: 实体身份连续性公理](#公理-3-实体身份连续性公理)
    - [公理 4: 值对象不可变性公理](#公理-4-值对象不可变性公理)
    - [定理 3: 聚合一致性定理](#定理-3-聚合一致性定理)
    - [定理 4: 领域事件溯源定理](#定理-4-领域事件溯源定理)
  - [四、概念关系属性图](#四概念关系属性图)
  - [五、实体 vs 值对象 决策树](#五实体-vs-值对象-决策树)
  - [六、聚合设计规则与示例](#六聚合设计规则与示例)
    - [规则 1: 聚合根唯一入口](#规则-1-聚合根唯一入口)
    - [规则 2: 事务边界 = 聚合边界](#规则-2-事务边界--聚合边界)
    - [规则 3: 聚合间通过ID引用](#规则-3-聚合间通过id引用)
  - [七、领域服务 vs 应用服务](#七领域服务-vs-应用服务)
  - [八、规格模式 (Specification Pattern) 详解](#八规格模式-specification-pattern-详解)
  - [九、场景适用矩阵](#九场景适用矩阵)
  - [十、常见反模式与解决方案](#十常见反模式与解决方案)
    - [反模式 1: 贫血领域模型](#反模式-1-贫血领域模型)
    - [反模式 2: 万能聚合](#反模式-2-万能聚合)
    - [反模式 3: 领域层依赖基础设施](#反模式-3-领域层依赖基础设施)

## 一、核心概念本体论

```text
DDD (Domain-Driven Design)
├── 战略设计 (Strategic Design)
│   ├── 限界上下文 (Bounded Context)
│   │   ├── 定义: 领域模型的边界
│   │   ├── 属性: 语言一致性、模型完整性
│   │   └── 关系: 上下文映射
│   ├── 上下文映射 (Context Mapping)
│   │   ├── 合作关系 (Partnership)
│   │   ├── 共享内核 (Shared Kernel)
│   │   ├── 客户-供应商 (Customer-Supplier)
│   │   ├── 遵奉者 (Conformist)
│   │   ├── 防腐层 (Anti-Corruption Layer)
│   │   ├── 开放主机服务 (Open Host Service)
│   │   └── 发布语言 (Published Language)
│   └── 通用语言 (Ubiquitous Language)
│       ├── 定义: 团队共享的术语
│       ├── 属性: 精确、一致、无处不在
│       └── 反例: 技术术语混入业务语言
│
└── 战术设计 (Tactical Design)
    ├── 实体 (Entity)
    │   ├── 定义: 有唯一标识的对象
    │   ├── 属性: ID、连续性、可变性
    │   └── 示例: 用户、订单、商品
    ├── 值对象 (Value Object)
    │   ├── 定义: 无标识，由属性决定
    │   ├── 属性: 不可变、可替换、相等性
    │   └── 示例: 地址、金额、时间段
    ├── 聚合 (Aggregate)
    │   ├── 定义: 一致性边界内的实体组
    │   ├── 属性: 根实体、边界、事务一致性
    │   └── 规则: 外部只能通过根访问
    ├── 领域服务 (Domain Service)
    │   ├── 定义: 不适合放在实体/值对象的行为
    │   ├── 属性: 无状态、表达业务逻辑
    │   └── 示例: 转账服务、推荐服务
    ├── 领域事件 (Domain Event)
    │   ├── 定义: 领域中发生的业务事件
    │   ├── 属性: 时态性、不可变、可订阅
    │   └── 示例: OrderCreated, PaymentReceived
    ├── 仓储 (Repository)
    │   ├── 定义: 聚合的持久化抽象
    │   ├── 属性: 集合语义、持久化透明
    │   └── 反例: 暴露数据库细节
    ├── 工厂 (Factory)
    │   ├── 定义: 复杂对象的创建逻辑
    │   └── 示例: 订单工厂创建复杂订单
    └── 规格模式 (Specification Pattern)
        ├── 定义: 业务规则的谓词封装
        ├── 属性: 可组合、可复用、可测试
        └── 操作: AND, OR, NOT
```

## 二、战略设计公理定理

### 公理 1: 限界上下文完整性公理

```text
定义: 限界上下文内的模型必须自洽
数学表达: ∀m ∈ Model(BC), Consistent(m) ∧ Complete(m)
约束: 跨上下文的模型必须通过显式映射
```

### 公理 2: 通用语言一致性公理

```
定义: 同一限界上下文内必须使用一致的领域语言
数学表达: ∀t ∈ Terms(BC), ∀u ∈ Users(BC), Meaning(t,u) = Constant
推论: 技术人员和业务人员使用相同术语
```

### 定理 1: 上下文映射复杂度定理

```
条件: 系统包含 n 个限界上下文
证明:
  1. 每对上下文间可能存在映射关系
  2. 映射关系数 ≤ n(n-1)/2
  3. 复杂度 O(n²)
结论: 限界上下文数量应控制在合理范围 (建议 3-7 个)
```

### 定理 2: 防腐层隔离定理

```
条件: 外部系统与核心领域通过防腐层交互
证明:
  1. 防腐层转换外部模型到内部模型
  2. 外部变化被限制在防腐层
  3. 核心领域不受外部影响
结论: 外部系统变化不影响核心业务逻辑
```

## 三、战术设计公理定理

### 公理 3: 实体身份连续性公理

```
定义: 实体身份在生命周期内保持不变
数学表达: ∀e ∈ Entity, ∀t₁,t₂ ∈ Time, ID(e,t₁) = ID(e,t₂)
约束: 实体的属性可以改变，身份不能
```

### 公理 4: 值对象不可变性公理

```
定义: 值对象创建后不可修改
数学表达: ∀vo ∈ ValueObject, Immutable(vo)
推论: 值对象修改 = 创建新实例
```

### 定理 3: 聚合一致性定理

```
条件: 聚合边界内的所有对象强一致
证明:
  1. 聚合根控制所有成员访问
  2. 业务规则在聚合内强制执行
  3. 事务边界与聚合边界对齐
结论: 聚合内数据始终满足业务规则
```

### 定理 4: 领域事件溯源定理

```
条件: 系统通过领域事件记录所有变更
证明:
  1. 当前状态 = Fold(初始状态, 事件列表)
  2. 事件不可变 append-only
  3. 因此可以重建任意历史状态
结论: 完整审计追踪 + 时间旅行调试
```

## 四、概念关系属性图

```
┌─────────────────────────────────────────────────────────────────────┐
│                      限界上下文 (Bounded Context)                    │
│                                                                     │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │
│   │   订单上下文 │    │   支付上下文 │    │   库存上下文 │            │
│   │             │    │             │    │             │            │
│   │  Order (AR) │◄───│ Payment (AR)│    │ Stock (AR)  │            │
│   │  OrderItem  │    │             │    │             │            │
│   │  Address(VO)│    │             │    │             │            │
│   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘            │
│          │                  │                  │                    │
│          │     ACL          │     ACL          │                    │
│          ▼                  ▼                  ▼                    │
│   ┌────────────────────────────────────────────────────┐           │
│   │              防腐层 (Anti-Corruption Layer)         │           │
│   └────────────────────────────────────────────────────┘           │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

AR = Aggregate Root, VO = Value Object, ACL = Anti-Corruption Layer
```

## 五、实体 vs 值对象 决策树

```
判断领域概念是实体还是值对象
│
├─ 是否需要唯一标识?
│   ├─ 是 (如用户ID、订单号) ──────────────► 实体 (Entity)
│   │   ├─ 属性可以变化
│   │   ├─ 有生命周期
│   │   └─ 通过ID相等
│   │
│   └─ 否 ────────────────────────────────► 值对象 (Value Object)
│       ├─ 属性决定相等性
│       ├─ 不可变
│       └─ 可替换
│
├─ 是否需要在数据库中独立表?
│   ├─ 是 ────────────────────────────────► 实体
│   └─ 否 (嵌入其他表) ───────────────────► 值对象
│
└─ 业务概念是否强调连续性?
    ├─ 是 (如用户历史订单) ───────────────► 实体
    └─ 否 (如金额、地址) ─────────────────► 值对象
```

## 六、聚合设计规则与示例

### 规则 1: 聚合根唯一入口

```go
// ✅ 正确：通过聚合根访问成员
order := orderRepo.FindByID(orderID)
order.AddItem(productID, quantity)  // 通过根操作
orderRepo.Save(order)

// ❌ 错误：直接操作聚合成员
item := orderItemRepo.FindByID(itemID)  // 错误！
item.UpdateQuantity(5)
```

### 规则 2: 事务边界 = 聚合边界

```go
// ✅ 正确：单个聚合内的事务
type Order struct {
    ID      string
    Items   []OrderItem
    Status  OrderStatus
}

func (o *Order) Confirm() error {
    // 所有变更在一个事务中
    o.Status = Confirmed
    for _, item := range o.Items {
        item.ReserveStock()
    }
    return nil
}
```

### 规则 3: 聚合间通过ID引用

```go
// ✅ 正确：聚合间松耦合
type Order struct {
    ID       string
    UserID   string    // 只存ID，不存User对象
    Items    []OrderItem
}

// ❌ 错误：直接引用其他聚合
type Order struct {
    ID   string
    User *User   // 错误！跨聚合引用
    Items []OrderItem
}
```

## 七、领域服务 vs 应用服务

```
┌────────────────────────────────────────────────────────────────┐
│                    服务类型对比                                 │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│   领域服务 (Domain Service)        应用服务 (Application Service)│
│   ───────────────────────          ───────────────────────────  │
│                                                                │
│   • 包含业务逻辑                   • 协调用例流程               │
│   • 无状态                         • 事务边界                   │
│   • 表达通用语言                   • 调用领域服务               │
│   • 属于领域层                     • 属于应用层                 │
│                                                                │
│   示例:                            示例:                       │
│   type TransferService struct {}   type OrderAppService struct{}│
│                                                                │
│   func (s *TransferService)        func (s *OrderAppService)    │
│       Transfer(from, to, amount)       CreateOrder(cmd) error { │
│       // 业务规则: 检查余额           // 1. 验证输入            │
│       // 业务规则: 检查限额           // 2. 开启事务            │
│       // 业务规则: 计算手续费         // 3. 调用领域服务         │
│       // ...                          // 4. 保存聚合            │
│                                         // 5. 发布事件          │
│                                         // 6. 提交事务          │
│                                    }                           │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

## 八、规格模式 (Specification Pattern) 详解

```go
// 规格接口
type Specification interface {
    IsSatisfiedBy(candidate interface{}) bool
    And(other Specification) Specification
    Or(other Specification) Specification
    Not() Specification
}

// 具体规格：库存充足
type InStockSpecification struct {
    inventory Inventory
}

func (s *InStockSpecification) IsSatisfiedBy(product interface{}) bool {
    p := product.(Product)
    return s.inventory.GetStock(p.ID) > 0
}

// 具体规格：价格范围内
type PriceRangeSpecification struct {
    Min, Max Money
}

func (s *PriceRangeSpecification) IsSatisfiedBy(product interface{}) bool {
    p := product.(Product)
    return p.Price.GreaterThanOrEqual(s.Min) &&
           p.Price.LessThanOrEqual(s.Max)
}

// 组合使用
func main() {
    inStock := &InStockSpecification{inventory}
    affordable := &PriceRangeSpecification{Min: 0, Max: 1000}

    // 可购买的商品 = 有库存 AND 价格合适
    available := inStock.And(affordable)

    products := productRepo.FindAll()
    for _, p := range products {
        if available.IsSatisfiedBy(p) {
            fmt.Println("可购买:", p.Name)
        }
    }
}
```

## 九、场景适用矩阵

| 场景 | DDD 适用度 | 原因 | 反模式 |
|------|-----------|------|--------|
| 复杂业务系统 | ⭐⭐⭐⭐⭐ | 业务逻辑复杂，需要领域建模 | 贫血模型 |
| 简单 CRUD | ⭐⭐ | 过度设计，增加复杂度 | 事务脚本 |
| 微服务架构 | ⭐⭐⭐⭐⭐ | 限界上下文天然对应服务边界 | 分布式单体 |
| 遗留系统重构 | ⭐⭐⭐⭐ | 防腐层隔离旧系统 | 大爆炸重写 |
| 数据密集型 | ⭐⭐⭐ | 关注点不同 | 误用聚合 |
| 算法密集型 | ⭐⭐ | 不是 DDD 强项 | 强行套用 |

## 十、常见反模式与解决方案

### 反模式 1: 贫血领域模型

```go
// ❌ 贫血模型：只有数据，没有行为
type Order struct {
    ID     string
    UserID string
    Status string
    Items  []OrderItem
}

// 所有业务逻辑在服务中
func (s *OrderService) ConfirmOrder(orderID string) {
    order := s.repo.Find(orderID)
    order.Status = "confirmed"  // 直接修改状态，无验证
    s.repo.Save(order)
}

// ✅ 充血模型：数据与行为封装
type Order struct {
    ID     string
    UserID string
    status OrderStatus
    Items  []OrderItem
}

func (o *Order) Confirm() error {
    if o.status != Pending {
        return errors.New("只能确认待处理订单")
    }
    if len(o.Items) == 0 {
        return errors.New("空订单不能确认")
    }
    o.status = Confirmed
    o.AddEvent(OrderConfirmedEvent{OrderID: o.ID})
    return nil
}
```

### 反模式 2: 万能聚合

```go
// ❌ 错误：聚合过大
type User struct {
    ID       string
    Profile  Profile
    Orders   []Order      // 错误！应该通过ID引用
    Addresses []Address
    Settings UserSettings
    // ... 更多
}

// ✅ 正确：小聚合，通过ID引用
type User struct {
    ID      string
    Profile Profile
}

type Order struct {
    ID     string
    UserID string  // 只存ID
    // ...
}
```

### 反模式 3: 领域层依赖基础设施

```go
// ❌ 错误：领域层依赖具体实现
package domain

import "github.com/some/orm"  // 错误！

type Order struct {
    orm.Model  // 错误！
}

// ✅ 正确：领域层纯洁，依赖倒置
package domain

// 领域层只定义接口
type OrderRepository interface {
    FindByID(id string) (*Order, error)
    Save(order *Order) error
}

// 实现放在基础设施层
package infrastructure

type OrderRepoImpl struct {
    db *sql.DB
}

func (r *OrderRepoImpl) FindByID(id string) (*domain.Order, error) {
    // 具体实现
}
```

---

**参考来源**:

- Domain-Driven Design: Tackling Complexity in the Heart of Software - Eric Evans, 2003
- Implementing Domain-Driven Design - Vaughn Vernon, 2013
- Patterns, Principles, and Practices of Domain-Driven Design - Scott Millett, 2015
