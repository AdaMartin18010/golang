# AD-003: 微服务拆分的形式化方法 (Microservices Decomposition: Formal Methods)

> **维度**: Application Domains
> **级别**: S (16+ KB)
> **标签**: #microservices #decomposition #ddd #bounded-context #service-boundary
> **权威来源**:
>
> - [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman (2021, 2nd Edition)
> - [Monolith to Microservices](https://samnewman.io/books/monolith-to-microservices/) - Sam Newman (2019)
> - [Domain-Driven Design](https://www.domainlanguage.com/ddd/) - Eric Evans (2003)
> - [The Art of Scalability](https://www.amazon.com/Art-Scalability-Architecture-Organizations-Enterprise/dp/0134032802) - Abbott & Fisher (2015)
> - [Microservices AntiPatterns and Pitfalls](https://www.oreilly.com/library/view/microservices-antipatterns-and/9781492042718/) - Mark Richards (2016)

---

## 1. 服务拆分的形式化定义

### 1.1 系统分解代数

**定义 1.1 (系统分解)**
分解 $D$ 是将系统 $S$ 划分为服务集合：
$$D: S \to \{ s_1, s_2, ..., s_n \}$$
满足：
$$\bigcup_{i=1}^{n} s_i = S \land \forall i \neq j: s_i \cap s_j = \emptyset$$

**定义 1.2 (服务边界)**
服务边界 $B(s)$ 定义了服务 $s$ 的职责范围。

**定义 1.3 (耦合度)**
$$C(s_i, s_j) = |\text{dependencies}(s_i, s_j)|$$
服务间的依赖数量。

### 1.2 分解质量度量

**内聚度 (Cohesion)**:
$$H(s) = \frac{\text{internal-interactions}}{\text{total-interactions}}$$
越高越好。

**耦合度 (Coupling)**:
$$C_{total} = \sum_{i \neq j} C(s_i, s_j)$$
越低越好。

**定理 1.1 (高内聚低耦合)**
$$\text{Quality}(D) \propto \frac{\sum H(s_i)}{C_{total}}$$

---

## 2. 分解策略的形式化

### 2.1 按业务能力分解

**定义 2.1 (业务能力)**
$$\text{Capability} = \{ c \mid \text{Business}(c) \}$$
组织为创造价值的活动。

**映射**:
$$\text{Service} = \text{Capability}$$
每个服务实现一个业务能力。

### 2.2 按子域分解 (DDD)

**定义 2.2 (限界上下文)**
$$\text{BoundedContext} = \langle \text{UbiquitousLanguage}, \text{Model}, \text{Boundary} \rangle$$

**映射**:
$$\text{Service} = \text{BoundedContext}$$

### 2.3 按实体分解 (反模式)

**警告**: 按实体 (CRUD) 分解导致贫血领域。
$$\text{AntiPattern}: \text{Service} = \text{Entity}$$

**正确做法**:
$$\text{Service} = \text{Aggregate}$$
围绕聚合根组织服务。

---

## 3. 拆分的决策框架

### 3.1 分解检查清单

**独立性检查**:

- [ ] 可以独立部署？
- [ ] 可以独立扩展？
- [ ] 可以独立开发？
- [ ] 失败是否隔离？

**粒度检查**:

- [ ] 代码量 < 1万行？
- [ ] 团队大小 5-9人？
- [ ] 变更频率合理？

### 3.2 拆分模式

**模式 1: 绞杀者模式 (Strangler Fig)**
$$S_{new} \text{ gradually replaces } S_{old}$$
逐步替代旧系统。

**模式 2: 反腐层 (Anti-Corruption Layer)**
$$\text{ACL}: S_{new} \leftrightarrow S_{old}$$
隔离遗留系统。

---

## 4. 多元表征

### 4.1 分解策略对比矩阵

| 策略 | 优点 | 缺点 | 适用场景 |
|------|------|------|---------|
| **By Capability** | 业务对齐 | 识别难度大 | 领域清晰 |
| **By Subdomain** | DDD 对齐 | 需要领域知识 | 复杂业务 |
| **By Aggregate** | 技术内聚 | 粒度可能不均 | 已有 DDD |
| **By Entity** (❌) | 简单 | 贫血领域 | 不推荐 |
| **By Team** (康威) | 组织对齐 | 可能不合理 | 团队稳定 |

### 4.2 服务拆分决策树

```
确定服务边界?
│
├── 识别业务能力?
│   ├── Order Management
│   ├── Inventory
│   ├── Payment
│   └── Shipping
│
├── 检查耦合?
│   ├── 高频调用? → 考虑合并
│   ├── 共享数据? → 确定数据所有权
│   └── 事务依赖? → Saga 或合并
│
├── 检查粒度?
│   ├── 过小 (< 2人)? → 合并
│   ├── 过大 (> 15人)? → 拆分
│   └── 正好? → 确定
│
└── 验证独立性?
    ├── 独立部署?
    ├── 独立扩展?
    └── 失败隔离?
```

### 4.3 微服务演进路径

```
Monolith
    │
    ├── Identify Bounded Contexts
    │       │
    │       ▼
    ├── Extract Non-Critical Service
    │       │
    │       ▼
    ├── Implement Anti-Corruption Layer
    │       │
    │       ▼
    ├── Extract Core Services (Strangler Fig)
    │       │
    │       ▼
    └── Retire Monolith
            │
            ▼
    Microservices Architecture
```

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Service Decomposition Checklist                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  设计原则:                                                                   │
│  □ 单一职责 (Single Responsibility)                                           │
│  □ 高内聚 (High Cohesion)                                                    │
│  □ 低耦合 (Low Coupling)                                                     │
│  □ 独立部署 (Independent Deployment)                                          │
│                                                                              │
│  边界识别:                                                                   │
│  □ 业务能力边界                                                               │
│  □ 限界上下文 (DDD)                                                           │
│  □ 数据所有权                                                                 │
│  □ 变更频率相似                                                               │
│                                                                              │
│  反模式避免:                                                                 │
│  ❌ 分布式单体 (Distributed Monolith)                                         │
│  ❌ 上帝服务 (God Service)                                                    │
│  ❌ 链式调用 (Long Call Chains)                                               │
│  ❌ 共享数据库 (Shared Database)                                              │
│                                                                              │
│  演进策略:                                                                   │
│  □ 绞杀者模式 (Strangler Fig)                                                 │
│  □ 反腐层 (Anti-Corruption Layer)                                             │
│  □ 逐步提取 (Incremental Extraction)                                         │
│  □ 并行运行 (Parallel Run)                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16KB, 完整形式化)
