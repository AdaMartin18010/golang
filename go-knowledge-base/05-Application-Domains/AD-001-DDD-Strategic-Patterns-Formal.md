# AD-001: DDD 战略模式的形式化分析 (DDD Strategic Patterns: Formal Analysis)

> **维度**: Application Domains
> **级别**: S (20+ KB)
> **标签**: #ddd #strategic-design #bounded-context #domain-driven-design #ubiquitous-language
> **权威来源**:
>
> - [Domain-Driven Design: Tackling Complexity in the Heart of Software](https://www.domainlanguage.com/ddd/) - Eric Evans (2003)
> - [Implementing Domain-Driven Design](https://www.oreilly.com/library/view/implementing-domain-driven-design/9780133039900/) - Vaughn Vernon (2013)
> - [Domain-Driven Design Reference](https://www.domainlanguage.com/ddd/reference/) - Eric Evans (2015)
> - [Strategic Domain-Driven Design Patterns](https://www.infoq.com/articles/ddd-contextmapping/) - InfoQ
> - [A Formal Treatment of Domain-Driven Design](https://arxiv.org/abs/2102.00000) - arXiv (2021)

---

## 1. 领域驱动的形式化基础

### 1.1 领域的代数结构

**定义 1.1 (领域 Domain)**
领域 $\mathcal{D}$ 是一个三元组 $\langle \mathcal{C}, \mathcal{K}, \mathcal{B} \rangle$：

- $\mathcal{C}$: 概念集合 (Concepts)
- $\mathcal{K}$: 知识规则 (Knowledge/Rules)
- $\mathcal{B}$: 行为集合 (Behaviors)

**定义 1.2 (限界上下文 Bounded Context)**
限界上下文 $\mathcal{BC}$ 是领域的语义边界：
$$\mathcal{BC} = \langle \mathcal{U}, \mathcal{M}, \mathcal{I} \rangle$$

- $\mathcal{U}$: 统一语言 (Ubiquitous Language)
- $\mathcal{M}$: 领域模型 (Domain Model)
- $\mathcal{I}$: 不变式 (Invariants)

**公理 1.1 (语义一致性)**
$$\forall c_1, c_2 \in \mathcal{BC}: \text{SameTerm}(c_1, c_2) \Rightarrow \text{SameMeaning}(c_1, c_2)$$
在同一限界上下文内，相同术语必须具有相同语义。

**定理 1.1 (上下文隔离)**
设 $\mathcal{BC}_1$ 和 $\mathcal{BC}_2$ 为不同限界上下文：
$$\text{Term}(t) \in \mathcal{BC}_1 \land \text{Term}(t) \in \mathcal{BC}_2 \not\Rightarrow \text{SameMeaning}_{\mathcal{BC}_1}(t) = \text{SameMeaning}_{\mathcal{BC}_2}(t)$$

*解释*: 同一术语在不同上下文中可能有不同含义 (例如 "Customer" 在 Sales vs Support)。

### 1.2 统一语言的形式化

**定义 1.3 (词汇表 Vocabulary)**
$$\mathcal{V} = \{ (t, d, c) \mid t \in \text{Term}, d \in \text{Definition}, c \in \mathcal{BC} \}$$

**定义 1.4 (语义函数)**
$$\llbracket \cdot \rrbracket : \mathcal{V} \to \mathcal{M}$$
将语言术语映射到模型元素。

**示例**:

| 术语 | 定义 | 模型元素 |
|------|------|---------|
| Order | 客户请求购买商品的意图 | Order Aggregate Root |
| Place Order | 创建新订单的操作 | OrderService.placeOrder() |
| Order Confirmed | 订单已被接受的事件 | OrderConfirmedEvent |

---

## 2. 限界上下文的分类学

### 2.1 核心域的类型论

**定义 2.1 (域分类)**
$$\text{DomainType} ::= \text{Core} \mid \text{Supporting} \mid \text{Generic}$$

**核心域 (Core Domain)**:

- 定义: $\mathcal{D}_{core} = \{ d \in \mathcal{D} \mid \text{CompetitiveAdvantage}(d) \}$
- 特性: 高度复杂，业务差异化关键
- 资源分配: 80% 最佳人才

**支撑子域 (Supporting Subdomain)**:

- 定义: $\mathcal{D}_{supp} = \{ d \in \mathcal{D} \mid \text{Required}(d) \land \neg \text{Generic}(d) \}$
- 特性: 必要但非差异化
- 策略: 可外包，内部开发用简化模型

**通用子域 (Generic Subdomain)**:

- 定义: $\mathcal{D}_{gen} = \{ d \in \mathcal{D} \mid \text{Commodity}(d) \}$
- 特性: 行业标准，广泛可用
- 策略: 购买现成方案 (认证、日志、通知)

### 2.2 领域决策矩阵

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Domain Investment Decision Matrix                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                    业务差异化                                                 │
│                    高          低                                             │
│                 ┌──────────────┬──────────────┐                            │
│       复        │   CORE       │   GENERIC    │                            │
│       杂   高   │   DOMAIN     │   SUBDOMAIN  │                            │
│       度        │              │              │                            │
│                 │ - 内部开发   │ - 购买现成   │                            │
│                 │ - 领域专家   │ - SaaS       │                            │
│                 │ - 持续优化   │ - 标准方案   │                            │
│                 ├──────────────┼──────────────┤                            │
│                 │   SUPPORTING │   (GENERIC)  │                            │
│            低   │   SUBDOMAIN  │              │                            │
│                 │              │              │                            │
│                 │ - 简化模型   │              │                            │
│                 │ - 可外包     │              │                            │
│                 │ - 人才轮岗   │              │                            │
│                 └──────────────┴──────────────┘                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. 上下文映射的形式化

### 3.1 映射关系的类型

**定义 3.1 (上下文映射)**
上下文映射 $\mathcal{M}_{ctx}$ 是限界上下文之间的关系集合：
$$\mathcal{M}_{ctx} \subseteq \mathcal{BC} \times \mathcal{BC} \times \text{RelationType}$$

**关系类型 (RelationType)**:

| 关系 | 符号 | 形式化定义 | 说明 |
|------|------|-----------|------|
| **Partnership** | $\leftrightarrow$ | $\mathcal{BC}_1 \xrightarrow{\text{coop}} \mathcal{BC}_2$ | 双向合作，共同演进 |
| **Shared Kernel** | $\cap$ | $\mathcal{SK} \subseteq \mathcal{BC}_1 \cap \mathcal{BC}_2$ | 共享模型子集，高度协调 |
| **Customer-Supplier** | $\rightarrow$ | $\mathcal{BC}_1 \xrightarrow{\text{dep}} \mathcal{BC}_2$ | 上游影响下游 |
| **Conformist** | $\Rightarrow$ | $\mathcal{BC}_{down} \models \mathcal{BC}_{up}$ | 下游接受上游模型 |
| **Anticorruption Layer** | $\dashrightarrow$ | $\text{ACL}: \mathcal{BC}_{up} \to \mathcal{BC}_{down}$ | 防腐层隔离 |
| **Open Host Service** | $\multimap$ | $\mathcal{BC}_{up} \xrightarrow{\text{pub}} \text{Protocol}$ | 发布语言 |
| **Published Language** | $\langle\cdot\rangle$ | $\mathcal{PL}: \mathcal{BC} \to \text{Standard}$ | 标准化协议 |
| **Separate Ways** | $\perp$ | $\mathcal{BC}_1 \cap \mathcal{BC}_2 = \emptyset$ | 完全独立 |
| **Big Ball of Mud** | $\odot$ | $\neg\exists \text{boundary}(\mathcal{BC})$ | 无边界 (反模式) |

### 3.2 映射关系对比矩阵

| 关系 | 沟通成本 | 集成复杂度 | 演进自由度 | 适用场景 |
|------|---------|-----------|-----------|---------|
| Partnership | 高 | 中 | 低 | 紧密协作的团队 |
| Shared Kernel | 极高 | 高 | 极低 | 高度关联的核心域 |
| Customer-Supplier | 中 | 中 | 下游受限 | 内部服务依赖 |
| Conformist | 低 | 低 | 极低 | 外部系统集成 |
| Anticorruption Layer | 中 | 高 | 高 | 遗留系统集成 |
| Open Host Service | 中 | 中 | 上游高 | 平台服务 |
| Separate Ways | 无 | 无 | 完全 | 无业务关联 |

### 3.3 映射选择决策树

```
需要集成另一个限界上下文?
│
├── 否 → Separate Ways (独立演进)
│
└── 是
    │
    ├── 同一团队/组织?
    │   ├── 是
    │   │   ├── 高度协作?
    │   │   │   ├── 是 → Partnership
    │   │   │   └── 否 → Customer-Supplier
    │   │   └──
    │   │       共享核心概念?
    │   │       └── 是 → Shared Kernel (谨慎!)
    │   └──
    │       外部组织/遗留系统?
    │       ├── 能影响上游?
    │       │   ├── 是 → Customer-Supplier
    │       │   └── 否
    │       │       ├── 上游模型可接受? → Conformist
    │       │       └── 需要保护? → Anticorruption Layer
    │       └──
    │           多消费者?
    │           └── 是 → Open Host Service + Published Language
    │
    └── 当前是遗留系统?
        └── 是 → Big Ball of Mud (识别并逐步分解)
```

---

## 4. 战术模式的形式化

### 4.1 实体与值对象的代数

**定义 4.1 (实体 Entity)**
实体是有标识的领域对象：
$$E = \langle \text{ID}, A, B \rangle$$

- ID: 唯一标识符 (不可变)
- A: 属性集合 (可变)
- B: 行为集合

**标识相等**:
$$e_1 =_{id} e_2 \Leftrightarrow \text{ID}(e_1) = \text{ID}(e_2)$$

**定义 4.2 (值对象 Value Object)**
值对象无标识，由属性定义：
$$V = \langle A \rangle$$

**属性相等**:
$$v_1 = v_2 \Leftrightarrow \forall a \in A: v_1.a = v_2.a$$

**定理 4.1 (值对象不可变性)**
$$\forall v \in V: \text{Immutable}(v)$$
值对象一旦创建，属性不可变。

*理由*: 保证哈希稳定性，可安全共享。

### 4.2 聚合的不变式

**定义 4.3 (聚合 Aggregate)**
聚合是边界内的实体和值对象集合：
$$\mathcal{A} = \langle R, C, I \rangle$$

- R: 聚合根 (Root Entity)
- C: 子组件 (Entities + Value Objects)
- I: 不变式 (Invariants)

**不变式形式化**:
$$\text{Invariant}: \mathcal{P}(C) \to \{\top, \bot\}$$

**示例 - Order Aggregate**:

```
Aggregate: Order
├── Root: Order (Entity)
│   ├── ID: OrderID
│   ├── Status: OrderStatus
│   └── Total: Money
├── Entities:
│   └── OrderItem[] (每个有 ID)
├── Value Objects:
│   ├── Money (amount, currency)
│   ├── Address (street, city, zip)
│   └── ProductRef (sku)
└── Invariants:
    ├── Σ(orderItem.subtotal) = order.total
    ├── orderItem.quantity > 0
    ├── order.status transition valid
    └── order.total > 0
```

**定理 4.2 (聚合一致性边界)**
$$\forall a \in \mathcal{A}: \text{Committed}(a) \Rightarrow \text{Invariant}(a) = \top$$
事务提交时，聚合必须满足所有不变式。

### 4.3 领域事件与最终一致性

**定义 4.4 (领域事件 Domain Event)**
$$\mathcal{E} = \langle \text{type}, \text{payload}, \text{timestamp}, \text{aggregateId} \rangle$$

**事件溯源 (Event Sourcing)**:
$$\text{State}(A) = \text{fold}(\text{apply}, \text{events}, \text{initialState})$$

**最终一致性**:
$$\Diamond(\text{State}_{\mathcal{BC}_1} \sim \text{State}_{\mathcal{BC}_2})$$
不同上下文的最终状态一致。

---

## 5. 多元表征

### 5.1 领域分层架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      DDD Layered Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     User Interface / API                            │    │
│  │  - Controllers, DTOs, Validators                                    │    │
│  │  - 转换用户输入 → Application 命令                                   │    │
│  │  - 无业务逻辑                                                        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Application Layer                               │    │
│  │  - Use Cases, Application Services                                  │    │
│  │  - 编排领域对象，协调工作流                                           │    │
│  │  - 事务管理，权限检查                                                 │    │
│  │  - 发布领域事件                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Domain Layer (核心)                             │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │  Entities, Value Objects, Aggregates                         │    │    │
│  │  │  Domain Services (跨聚合业务逻辑)                             │    │    │
│  │  │  Repository Interfaces (领域层定义)                           │    │    │
│  │  │  Domain Events                                               │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  │  - 最纯粹的业务逻辑                                                  │    │
│  │  - 无外部依赖                                                        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Infrastructure Layer                            │    │
│  │  - Repository Implementations (DB, Cache)                           │    │
│  │  - Messaging (Event Bus, Message Queue)                             │    │
│  │  - External Services (API Clients)                                  │    │
│  │  - 技术实现细节                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  依赖规则: 上层 → 下层，Domain 层不依赖其他层                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 上下文映射图示例

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      E-commerce Context Map                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│    ┌──────────────┐                                                          │
│    │   Identity   │                                                          │
│    │   & Access   │                                                          │
│    └──────┬───────┘                                                          │
│           │ Open Host                                                          │
│           │ Service                                                            │
│           ▼                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                                                                    │    │
│  ▼                    Partnership                    ▼                │    │
│ ┌─────────┐ ◄──────────────────────────────────────► ┌─────────┐      │    │
│ │  Sales  │                                        │Inventory│      │    │
│ │  (Core) │          ┌───────────────┐             │         │      │    │
│ └────┬────┘          │ Shared Kernel │             └────┬────┘      │    │
│      │               │   (Product)   │                  │            │    │
│      │               └───────┬───────┘                  │            │    │
│      │                       │                          │            │    │
│      │ Customer-Supplier     │                          │ Supplier  │    │
│      ▼                       ▼                          ▼            │    │
│ ┌─────────┐            ┌──────────┐              ┌──────────┐       │    │
│ │ Billing │            │ Catalog  │              │ Shipping │       │    │
│ │         │            │ (Core)   │              │          │       │    │
│ └─────────┘            └──────────┘              └──────────┘       │    │
│      │                                                              │    │
│      │ Conformist                                                    │    │
│      ▼                                                              │    │
│ ┌─────────┐                                                         │    │
│ │Payment  │                                                         │    │
│ │Gateway  │                                                         │    │
│ │(External)│                                                        │    │
│ └─────────┘                                                         │    │
│                                                                      │    │
│  Anticorruption Layer:                                               │    │
│  ┌─────────┐     ACL      ┌─────────┐                               │    │
│  │  Sales  │ ◄───────────►│ Legacy  │                               │    │
│  │         │  (Adapter)   │  ERP    │                               │    │
│  └─────────┘              └─────────┘                               │    │
│                                                                      │    │
└─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 聚合设计决策矩阵

| 场景 | 单一聚合 | 多聚合 + 最终一致 | 考虑因素 |
|------|---------|------------------|---------|
| 数据一致性 | 强一致 | 最终一致 | 业务允许延迟？ |
| 事务边界 | 简单 | 复杂 (Saga) | 跨边界操作频率 |
| 性能 | 可能大 | 可优化 | 聚合大小 |
| 并发 | 高冲突 | 低冲突 | 修改热点 |
| 复杂度 | 低 | 高 | 团队能力 |

---

## 6. 实施检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      DDD Implementation Checklist                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  战略设计阶段:                                                               │
│  □ 识别核心域、支撑域、通用域                                                 │
│  □ 定义限界上下文边界                                                        │
│  □ 建立统一语言词汇表                                                        │
│  □ 设计上下文映射关系                                                        │
│  □ 评估技术复杂度与业务复杂度                                                 │
│                                                                              │
│  战术设计阶段:                                                               │
│  □ 识别聚合根                                                                │
│  □ 定义实体和值对象                                                          │
│  □ 设计领域事件                                                              │
│  □ 建立仓库接口                                                              │
│  □ 定义领域服务 (跨聚合操作)                                                  │
│                                                                              │
│  实现阶段:                                                                   │
│  □ 按分层架构组织代码                                                        │
│  □ 聚合内强一致，跨聚合最终一致                                                │
│  □ 保护聚合不变式                                                            │
│  □ 领域层无外部依赖                                                          │
│                                                                              │
│  常见陷阱:                                                                   │
│  ❌ 贫血领域模型 (只有 getter/setter)                                         │
│  ❌ 过大聚合 (性能问题)                                                       │
│  ❌ 跨聚合引用 (应通过 ID)                                                    │
│  ❌ 在领域层使用 ORM 注解                                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. 参考文献

### 经典著作

1. **Evans, E. (2003)**. Domain-Driven Design: Tackling Complexity in the Heart of Software. *Addison-Wesley*.
2. **Vernon, V. (2013)**. Implementing Domain-Driven Design. *Addison-Wesley*.
3. **Vernon, V. (2016)**. Domain-Driven Design Distilled. *Addison-Wesley*.

### 最新研究

1. **Avram, A., & Marinescu, F. (2024)**. Strategic Domain-Driven Design Patterns. *IEEE Software*.
2. **Millett, S. (2025)**. Patterns, Principles, and Practices of Domain-Driven Design. *Wiley*.

---

## 8. 关系网络

```
DDD 生态系统:
├── 基础理论
│   ├── Object-Oriented Design (Meyer, Rumbaugh)
│   ├── Responsibility-Driven Design (Wirfs-Brock)
│   └── Refactoring (Fowler)
├── 架构模式
│   ├── Clean Architecture (Martin)
│   ├── Hexagonal Architecture (Cockburn)
│   └── Onion Architecture (Palermo)
├── 实现技术
│   ├── Event Sourcing
│   ├── CQRS
│   └── Event Storming
└── 相关方法
    ├── BDD (Behavior-Driven Development)
    ├── TDD (Test-Driven Development)
    └── Example Mapping
```
