# Clean Architecture 核心概念体系

## 目录

- [Clean Architecture 核心概念体系](#clean-architecture-核心概念体系)
  - [目录](#目录)
  - [一、核心公理与定理](#一核心公理与定理)
    - [公理 1: 依赖方向公理 (Dependency Rule Axiom)](#公理-1-依赖方向公理-dependency-rule-axiom)
    - [公理 2: 抽象层次公理 (Abstraction Level Axiom)](#公理-2-抽象层次公理-abstraction-level-axiom)
    - [定理 1: 可测试性定理 (Testability Theorem)](#定理-1-可测试性定理-testability-theorem)
    - [定理 2: 框架独立性定理 (Framework Independence Theorem)](#定理-2-框架独立性定理-framework-independence-theorem)
  - [二、概念本体论 (Ontology)](#二概念本体论-ontology)
  - [三、概念关系属性图](#三概念关系属性图)
  - [四、决策推理树](#四决策推理树)
  - [五、属性定义与约束](#五属性定义与约束)
  - [六、示例与反例](#六示例与反例)
    - [示例 1: 正确的依赖方向](#示例-1-正确的依赖方向)
    - [反例 1: 错误的依赖方向](#反例-1-错误的依赖方向)
    - [示例 2: 正确的边界跨越](#示例-2-正确的边界跨越)
    - [反例 2: 绕过边界](#反例-2-绕过边界)
  - [七、与相关架构的关系](#七与相关架构的关系)
  - [八、证明树：为什么 Clean Architecture 有效](#八证明树为什么-clean-architecture-有效)
  - [九、应用场景矩阵](#九应用场景矩阵)
  - [十、形式化验证](#十形式化验证)

## 一、核心公理与定理

### 公理 1: 依赖方向公理 (Dependency Rule Axiom)

```
定义: 源代码依赖关系只能向内指向更高级的策略
数学表达: ∀A,B ∈ Layers, Dependency(A,B) → Level(A) > Level(B)
推论: 外层可以依赖内层，内层不能依赖外层
```

### 公理 2: 抽象层次公理 (Abstraction Level Axiom)

```
定义: 越靠近核心的层，抽象程度越高
数学表达: ∀L ∈ Layers, Abstraction(L) ∝ 1/Distance(L, Core)
推论: Entities > Use Cases > Interface Adapters > Frameworks
```

### 定理 1: 可测试性定理 (Testability Theorem)

```
条件: 业务逻辑不依赖框架、UI、数据库
证明:
  1. 业务逻辑在 Entities 和 Use Cases 层
  2. 这两层没有外部依赖（根据公理1）
  3. 因此可以用 mock/stub 替代所有依赖
  4. 所以业务逻辑可以独立测试
结论: Testability(Business Logic) = 100%
```

### 定理 2: 框架独立性定理 (Framework Independence Theorem)

```
条件: 框架代码只在最外层
证明:
  1. 框架代码在 Frameworks & Drivers 层
  2. 业务逻辑在内部三层
  3. 内部层不依赖外部层（根据公理1）
  4. 因此更换框架不影响业务逻辑
结论: Framework(业务逻辑) = ∅
```

## 二、概念本体论 (Ontology)

```
Clean Architecture
├── 核心原则
│   ├── 依赖方向: 向内指向策略，向外指向细节
│   ├── 分层隔离: Entities > Use Cases > Adapters > Frameworks
│   └── 边界保护: 依赖反转 + 接口隔离
├── 组件类型
│   ├── Entities: 企业级业务规则，最抽象，最稳定
│   ├── Use Cases: 应用业务规则，编排实体，定义流程
│   ├── Interface Adapters: 格式转换，MVC/MVP，数据映射
│   └── Frameworks: 具体实现，最易变，最外层
├── 设计原则
│   ├── SOLID: SRP, OCP, LSP, ISP, DIP
│   └── 其他: DRY, KISS, YAGNI
└── 质量属性
    ├── 可测试性
    ├── 可维护性
    ├── 可扩展性
    └── 可移植性
```

## 三、概念关系属性图

```
                    ┌─────────────────────────────────────┐
                    │        Frameworks & Drivers         │
                    │   (Web, DB, External Services)      │
                    └──────────────┬──────────────────────┘
                                   │ 依赖
                                   ▼
                    ┌─────────────────────────────────────┐
                    │      Interface Adapters             │
                    │  (Controllers, Presenters, Gateways)│
                    └──────────────┬──────────────────────┘
                                   │ 依赖
                                   ▼
                    ┌─────────────────────────────────────┐
                    │         Use Cases                   │
                    │    (Application Business Rules)     │
                    └──────────────┬──────────────────────┘
                                   │ 依赖
                                   ▼
                    ┌─────────────────────────────────────┐
                    │           Entities                  │
                    │      (Enterprise Business Rules)    │
                    └─────────────────────────────────────┘
```

## 四、决策推理树

```
Clean Architecture 决策
│
├─ 业务逻辑复杂度?
│   ├─ 简单 (CRUD为主) ────────► 简化分层 (应用三层架构)
│   └─ 复杂 (复杂业务规则) ────► 完整分层 (四层架构+DDD)
│
├─ 技术栈稳定性?
│   ├─ 稳定 ─────────────────► 内层抽象 (定义清晰接口)
│   └─ 不稳定 ───────────────► 外层隔离 (适配器模式)
│
├─ 团队规模?
│   ├─ 小团队 ───────────────► 轻量架构 (聚焦核心)
│   └─ 大团队 ───────────────► 严格分层 (边界清晰)
│
└─ 质量要求?
    ├─ 高要求 ───────────────► 完整测试 (单元+集成+E2E)
    └─ 一般要求 ─────────────► 核心测试 (业务逻辑优先)
```

## 五、属性定义与约束

| 属性 | 定义 | 值域 | 约束 |
|------|------|------|------|
| **稳定性** | 组件变化频率的倒数 | [0,1] | Entities > Use Cases > Adapters > Frameworks |
| **抽象度** | 与具体实现无关的程度 | [0,1] | 与层深度成正比 |
| **依赖方向** | 源代码依赖的方向 | {in, out} | 必须为 in |
| **测试覆盖率** | 被测试覆盖的代码比例 | [0,100] | Entities 层 ≥ 90% |
| **接口数量** | 定义的稳定接口数 | ℕ | 内层 > 外层 |

## 六、示例与反例

### 示例 1: 正确的依赖方向

```go
// ✅ 正确：Use Cases 依赖 Entities，不依赖框架
package usecases

import "domain/entities"  // 向内依赖

type CreateOrderUseCase struct {
    orderRepo entities.OrderRepository  // 依赖接口，不是实现
}

func (uc *CreateOrderUseCase) Execute(input CreateOrderInput) error {
    order := entities.NewOrder(input.UserID, input.Items)
    return uc.orderRepo.Save(order)
}
```

### 反例 1: 错误的依赖方向

```go
// ❌ 错误：Entities 依赖外部框架
package entities

import "github.com/gin-gonic/gin"  // 错误！实体依赖框架

type Order struct {
    // ...
}

func (o *Order) ToJSON(c *gin.Context) {  // 错误！实体知道 HTTP
    c.JSON(200, o)
}
```

### 示例 2: 正确的边界跨越

```go
// ✅ 正确：通过 DTO 跨越边界
package controllers

import "usecases"

type OrderController struct {
    createOrderUC *usecases.CreateOrderUseCase
}

func (ctrl *OrderController) Create(w http.ResponseWriter, r *http.Request) {
    // 1. 适配外部格式到内部 DTO
    var req OrderRequest
    json.NewDecoder(r.Body).Decode(&req)

    input := usecases.CreateOrderInput{
        UserID: req.UserID,
        Items:  adaptItems(req.Items),
    }

    // 2. 调用用例
    err := ctrl.createOrderUC.Execute(input)

    // 3. 适配内部结果到外部格式
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(201)
}
```

### 反例 2: 绕过边界

```go
// ❌ 错误：Controller 直接操作数据库
package controllers

import "database/sql"

type OrderController struct {
    db *sql.DB  // 错误！跨越了两层
}

func (ctrl *OrderController) Create(w http.ResponseWriter, r *http.Request) {
    // 错误！业务逻辑在 Controller 中
    _, err := ctrl.db.Exec("INSERT INTO orders ...")
    // ...
}
```

## 七、与相关架构的关系

```
┌─────────────────────────────────────────────────────────────┐
│                     架构谱系关系                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   Clean Architecture ────────┐                              │
│   (Robert C. Martin)         │                              │
│                              ├── 共同原则：依赖向内          │
│   Hexagonal Architecture ────┤                              │
│   (Alistair Cockburn)        │                              │
│                              ├── 侧重点：端口与适配器        │
│   Onion Architecture ────────┤                              │
│   (Jeffrey Palermo)          │                              │
│                              ├── 侧重点：领域为中心          │
│   DDD Layered ───────────────┘                              │
│   (Eric Evans)                                              │
│                              侧重点：业务复杂度              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## 八、证明树：为什么 Clean Architecture 有效

```
目标: Clean Architecture 能够降低软件复杂度
│
├─ 子目标 1: 分离关注点
│   ├─ 前提: 不同变化速率的代码应该分离
│   ├─ 前提: 业务规则变化慢，框架变化快
│   └─ 结论: 分层隔离满足不同变化速率
│
├─ 子目标 2: 提高可测试性
│   ├─ 前提: 测试需要控制所有依赖
│   ├─ 前提: 依赖向内使内部层不依赖外部
│   └─ 结论: 可以用 mock 替代外部依赖
│
├─ 子目标 3: 延迟技术决策
│   ├─ 前提: 早期决策往往基于不完整信息
│   ├─ 前提: 业务逻辑不依赖具体技术
│   └─ 结论: 可以在后期选择最合适的技术
│
└─ 子目标 4: 支持团队并行
    ├─ 前提: 清晰的接口定义减少协调成本
    ├─ 前提: 层间通过接口交互
    └─ 结论: 不同团队可以独立开发不同层
```

## 九、应用场景矩阵

| 场景 | 适用度 | 原因 | 反例场景 |
|------|--------|------|----------|
| 复杂业务系统 | ⭐⭐⭐⭐⭐ | 业务规则多，需要隔离 | 简单 CRUD |
| 长期维护项目 | ⭐⭐⭐⭐⭐ | 框架会更新，业务稳定 | 短期原型 |
| 多团队项目 | ⭐⭐⭐⭐⭐ | 边界清晰，并行开发 | 单人项目 |
| 微服务 | ⭐⭐⭐⭐ | 服务内部分层 | 无业务逻辑的网关 |
| 嵌入式系统 | ⭐⭐⭐ | 资源限制可能需要简化 | 资源极受限 |
| 脚本工具 | ⭐⭐ | 过度设计 | 一次性脚本 |

## 十、形式化验证

```
定理: Clean Architecture 保证业务逻辑可独立测试

证明:
设:
  - L = {Entities, UseCases, Adapters, Frameworks}
  - Dependency(x,y): x 依赖 y
  - Independent(x): x 可以独立测试

已知:
  1. ∀a ∈ Entities, ∀b ∈ Layers, Dependency(a,b) → b ∈ Entities
     (Entities 只依赖 Entities)
  2. Independent(x) ↔ ∀y, Dependency(x,y) → y 可被 mock

证明 Entities 可独立测试:
  取任意 e ∈ Entities
  根据已知 1, e 依赖的都在 Entities
  Entities 中只有接口定义，没有外部依赖
  因此 Entities 的依赖可以被 mock
  根据已知 2, Independent(e) 成立

结论: Entities 可以独立测试 ∎
```

---

**参考来源**:

- Clean Architecture - Robert C. Martin, 2017
- Domain-Driven Design - Eric Evans, 2003
- Implementing Domain-Driven Design - Vaughn Vernon, 2013
- Get Your Hands Dirty on Clean Architecture - Tom Hombergs, 2019
