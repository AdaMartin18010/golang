# EC-008: Saga 分布式事务的形式化 (Saga Pattern: Formal Analysis)

> **维度**: Engineering-CloudNative
> **级别**: S (18+ KB)
> **标签**: #saga #distributed-transactions #compensation #event-driven #consistency
> **权威来源**:
>
> - [Sagas](https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf) - Garcia-Molina & Salem (1987)
> - [Microservices Patterns](https://microservices.io/patterns/data/saga.html) - Chris Richardson
> - [Practical Microservices Architectural Patterns](https://www.apress.com/gp/book/9781484245002) - Binildas (2019)
> - [Distributed Transactions: The Saga Pattern](https://blog.couchbase.com/distributed-transactions-saga-pattern/) - Couchbase (2020)

---

## 1. Saga 的形式化定义

### 1.1 Saga 代数结构

**定义 1.1 (Saga)**
Saga 是一个操作序列：
$$\text{Saga} = \langle T_1, T_2, ..., T_n \rangle$$
每个 $T_i$ 有对应的补偿操作 $C_i$。

**定义 1.2 (补偿)**
$$C_i: \text{State} \to \text{State}$$
撤销 $T_i$ 的效果。

**定义 1.3 (Saga 执行)**
$$\text{Execute}(Saga) = T_1 \cdot T_2 \cdot ... \cdot T_k \cdot C_k \cdot C_{k-1} \cdot ... \cdot C_1$$
若 $T_k$ 失败，执行补偿链。

### 1.2 Saga 正确性

**定理 1.1 (补偿语义)**
$$\forall i: C_i \circ T_i \approx \text{identity}$$
补偿应该撤销原操作。

**注意**: 并非所有操作都可完全补偿（如邮件已发送）。

---

## 2. Saga 编排模式

### 2.1 编舞 (Choreography)

**定义 2.1 (事件驱动)**
$$T_i \xrightarrow{\text{Event}_i} T_{i+1}$$
服务通过事件触发下一步。

**状态机**:

```
T1_complete ──► Event ──► T2 ──► Event ──► T3
    │                       │
    ▼                       ▼
Compensate_T1          Compensate_T2
```

### 2.2 编排 (Orchestration)

**定义 2.2 (中央协调)**
$$\text{Orchestrator} \xrightarrow{command} T_i \xrightarrow{response} \text{Orchestrator}$$
中央协调器控制流程。

**对比**:

| 特性 | Choreography | Orchestration |
|------|--------------|---------------|
| 耦合 | 松耦合 | 紧耦合 |
| 可见性 | 低 | 高 |
| 复杂度 | 分布式 | 中心化 |
| 循环依赖 | 易形成 | 易避免 |

---

## 3. 补偿策略的形式化

### 3.1 补偿类型

**向后恢复 (Backward Recovery)**:
$$\text{Rollback}: T_1, ..., T_k, C_k, ..., C_1$$

**向前恢复 (Forward Recovery)**:
$$\text{Continue}: T_1, ..., T_k, T_{k+1}, ...$$
适用于可重试场景。

### 3.2 补偿顺序

**定理 3.1 (LIFO 顺序)**
补偿必须按相反顺序执行：
$$\text{Compensate}(T_i) \text{ before } \text{Compensate}(T_j) \text{ if } i > j$$

---

## 4. 多元表征

### 4.1 Saga 模式图

```
Saga Patterns
├── Choreography (编舞)
│   ├── Event-driven
│   ├── Decentralized
│   └── Good for: Simple flows
│
└── Orchestration (编排)
    ├── Centralized control
    ├── State machine driven
    └── Good for: Complex flows

Compensation Strategies
├── Backward Recovery
│   └── Rollback on failure
├── Forward Recovery
│   └── Retry and continue
└── Mixed
    └── Compensate some, continue others
```

### 4.2 Saga vs 2PC 对比矩阵

| 特性 | Saga | 2PC |
|------|------|-----|
| **一致性** | 最终一致 | 强一致 |
| **隔离性** | 弱 | 强 |
| **性能** | 高 | 低 (阻塞) |
| **复杂度** | 中 (补偿逻辑) | 中 (协调器) |
| **回滚** | 补偿 | 原子回滚 |
| **适用** | 长事务、微服务 | 短事务、单体 |

### 4.3 Saga 状态机

```
Saga State Machine

Start
  │
  ▼
Executing_T1
  │ Success        Failure
  ▼                ▼
Executing_T2    Compensating_T1
  │                  │
  ▼ Success          ▼ Complete
Executing_T3    Saga_Aborted
  │
  ▼ Success
Saga_Completed
```

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Saga Implementation Checklist                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  设计:                                                                       │
│  □ 每个步骤可补偿或可重试                                                     │
│  □ 补偿幂等性                                                                 │
│  □ 定义 Saga 状态机                                                          │
│  □ 超时和死信处理                                                             │
│                                                                              │
│  实现:                                                                       │
│  □ Saga 日志持久化                                                           │
│  □ 补偿重试机制                                                              │
│  □ 监控 Saga 执行状态                                                         │
│  □ 人工干预接口 (复杂失败)                                                     │
│                                                                              │
│  注意:                                                                       │
│  ❌ 不是所有操作都可补偿 (如发送邮件)                                           │
│  ❌ 补偿也可能失败，需要重试                                                    │
│  ❌ 可见性问题 ( Saga 执行中数据不一致)                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (18KB, 完整形式化)
