# EC-081-Distributed-Systems-Research-2025

> **Dimension**: 03-Engineering-CloudNative
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: OSDI 2025, EuroSys 2025, OPODIS 2025
> **Size**: >20KB

---

## 1. 研究概览

### 1.1 顶级会议 2025

| 会议 | 日期 | 主要主题 |
|------|------|---------|
| OSDI 2025 | July 2025 | Distributed Systems, Data Centers |
| EuroSys 2025 | March 2025 | Systems, Distributed Computing |
| OPODIS 2025 | December 2025 | Principles of Distributed Systems |
| NSDI 2025 | April 2025 | Networked Systems |

### 1.2 突破性研究

1. **Basilisk** (OSDI 2025 Best Paper) - 自动协议证明
2. **Picsou** - 跨RSM高效通信
3. **T2C** - 从测试推导语义检查器
4. **CAC** (OPODIS 2025) - Contention-Aware Cooperation

---

## 2. Basilisk: 自动协议证明

### 2.1 论文信息

- **作者**: Tony Nuda Zhang, Keshav Singh, Tej Chajed, Manos Kapritsos, Bryan Parno
- **机构**: University of Michigan, UW-Madison, CMU
- **荣誉**: OSDI 2025 Best Paper

### 2.2 核心问题

分布式协议设计复杂，形式化验证困难:

- 安全属性(safety)无法为归纳提供有力论证
- 开发者需大量时间寻找归纳不变量(inductive invariant)

### 2.3 Basilisk 方法

**利用来源不变量(provenance invariants)自动化证明**

```
核心思想:
1. 追踪数据在协议中的流动(provenance)
2. 提取不变量约束
3. 自动生成归纳证明
```

**形式化优势**:

- 减少手动证明工作量
- 提高验证自动化程度
- 适用于复杂协议

---

## 3. Picsou: 跨RSM高效通信

### 3.1 论文信息

- **作者**: Reginald Frank, Micah Murray, Chawinphat Tankuranand, et al.
- **机构**: UC Berkeley, University of Oregon, University of Michigan

### 3.2 问题背景

**跨复制状态机(RSM)通信挑战**:

现有方案:

1. **第三方日志(Kafka)**: 增加复杂性
2. **All-To-All**: 广域网成本高
3. **Leader-To-Leader**: 单点故障

### 3.3 C3B 原语

**Cross-Cluster Consistent Broadcast (C3B)**

**系统模型**:

- 两个全双工通信的RSM: R_s, R_r
- UpRight故障模型
  - Commission Failure: 拜占庭节点
  - Omission Failure: 崩溃节点

**核心操作**:

| 操作 | 描述 |
|------|------|
| **Transmit** | RSM R_s 向 R_r 传输消息m |
| **Deliver** | RSM R_r 从 R_s 交付消息m |

**正确性保证**:

1. **Validity**: 正确发送的消息最终交付
2. **Agreement**: 如果任何正确副本交付，所有正确副本都交付
3. **Integrity**: 每条消息最多交付一次

### 3.4 通信协议设计

借鉴TCP理念，应用于多对多RSM通信:

- 可靠消息流
- 流量控制
- 拥塞控制

---

## 4. T2C: 从测试推导检查器

### 4.1 核心思想

**T2C (Tests To Checkers)**: 将现有测试用例转换为语义检查器，检测生产环境静默故障。

### 4.2 设计与实现

**框架组件**:

给定测试用例 T，推导运行时检查器 C:

```
C = (F_c, P_c, R_c)

其中:
- F_c: 参数化检查器函数
- P_c: 符号检查器先决条件
- R_c: 参数关系约束
```

**转换方法**:

1. 保留测试代码结构
2. 提取并泛化核心验证逻辑
3. 生成运行时检查器

### 4.3 应用场景

- **静默故障检测**: 生产环境无法被传统监控捕获的故障
- **语义验证**: 验证系统行为是否符合规范
- **回归防护**: 防止已修复问题重新出现

---

## 5. CAC: Contention-Aware Cooperation

### 5.1 论文信息

- **作者**: Timothe Albouy, Davide Frey, Mathieu Gestin, Michel Raynal, Francois Taiani
- **会议**: OPODIS 2025

### 5.2 抽象定义

**Contention-Aware Cooperation (CAC)**: 新型的分布式协作抽象

**特征对比**:

| 抽象 | 协作模式 | 规模 |
|------|---------|------|
| Reliable Broadcast | 1-to-n | n |
| Consensus | n-to-n | n |
| **CAC** | **d-to-n** | **动态 d (1≤d≤n)** |

### 5.3 核心特性

**动态提议者集合**:

- d 随每次运行变化
- 进程不知道 d 的值
- l 个值被接受 (1≤l≤d)

**不完美预言机**:

- 提供未来可能接受的值的洞察
- 在有利条件下(预言机准确)特别高效

### 5.4 应用

1. **Cascading Consensus**: 低竞争优化快速共识
2. **命名问题**: 在完全异步下解决

---

## 6. 其他重要研究

### 6.1 共识协议进展

| 协议 | 特点 | 性能 |
|------|------|------|
| **Baxos** | Leaderless共识 | 128%更好的攻击弹性 |
| **Mako** (OSDI 2025) | 推测性geo-replication | 3.66M TPS |
| **Eg-walker** (EuroSys 2025) | CRDT/OT混合 | 160,000x更快 |

### 6.2 区块链共识

**BBcA-Chain** (FC 2024):

- DAG-based BFT共识
- 低延迟、高吞吐
- 快速确认

**Narwhal and Tusk** (EuroSys 2022):

- DAG-based mempool
- 高效BFT共识
- 解耦数据传播与共识

### 6.3 容错计算

**HQ Replication** (OSDI 2006):

- Hybrid Quorum协议
- 拜占庭容错

**BEAT** (CCS 2018):

- 异步BFT
- 实用化实现

---

## 7. 研究趋势分析

### 7.1 形式化验证趋势

1. **自动化证明**: 减少手动工作量
2. **归纳不变量生成**: 机器学习辅助
3. **可组合验证**: 模块化证明

### 7.2 跨域通信趋势

1. **RSM间通信标准化**: C3B等原语
2. **异构系统互操作**: 不同共识协议间通信
3. **WAN优化**: 广域网高效通信

### 7.3 测试与生产融合

1. **测试用例复用**: T2C方法
2. **持续验证**: 运行时检查
3. **静默故障检测**: 传统监控盲区

---

## 8. 对工程实践的影响

### 8.1 协议设计

- 考虑形式化验证友好性
- 设计可证明的安全属性
- 使用自动化验证工具

### 8.2 系统实现

- 集成运行时检查
- 从测试推导监控
- 关注跨系统通信效率

### 8.3 工具链

- 采用Basilisk等验证工具
- 实施T2C转换流程
- 探索CAC抽象应用

---

## 9. 参考文献

### OSDI 2025

1. Basilisk: Using Provenance Invariants to Automate Proofs
2. Picsou: Enabling Replicated State Machines to Communicate Efficiently
3. Deriving Semantic Checkers from Tests
4. Mako: Speculative Geo-Replication

### OPODIS 2025

1. Contention-Aware Cooperation
2. Impossibility of Distributed Consensus with One Faulty Process (Fischer et al. revisit)

### 经典论文

1. Practical Byzantine Fault Tolerance (Castro & Liskov, OSDI 1999)
2. In Search of an Understandable Consensus Algorithm (Raft, USENIX ATC 2014)
3. HQ Replication (Cowling et al., OSDI 2006)

---

*Last Updated: 2026-04-03*

---

## 10. 深入: Basilisk 技术细节

### 10.1 Provenance Invariants

**定义**: 追踪数据在协议执行中的来源和变换

**示例**:
`
消息m的来源:
  m.payload = Hash(m.previous) + m.data

不变量:
  ∀ m: VerifyHash(m) = true ⟹ m 是有效的
`

### 10.2 自动化证明流程

`

1. 协议实现 → 2. Provenance追踪 → 3. 不变量提取
     ↓
2. 归纳证明 ← 5. 验证器检查 ← 6. 证明生成
`

### 10.3 与TLA+对比

| 特性 | TLA+ | Basilisk |
|------|------|----------|
| 学习曲线 | 陡峭 | 平缓 |
| 自动化 | 低 | 高 |
| 适用范围 | 通用 | 协议特定 |
| 证明保证 | 完整 | 自动化推导 |

---

## 11. Picsou C3B 协议详解

### 11.1 C3B 形式化定义

**原语**: Cross-Cluster Consistent Broadcast

**接口**:
`
Transmit(R_s, R_r, m): R_s 向 R_r 传输消息 m
Deliver(R_r, R_s, m): R_r 从 R_s 交付消息 m
`

**属性**:

1. **Validity**: Transmit(R_s, R_r, m) 最终被 Deliver
2. **Agreement**: 如果一个正确副本Deliver，所有正确副本Deliver
3. **Integrity**: 每条消息最多Deliver一次
4. **Total Order**: 所有副本以相同顺序Deliver

### 11.2 故障模型

**UpRight模型**:

- n个节点
- 最多r个commission故障(拜占庭)
- 最多u个omission故障(崩溃)

**约束**: n > 2r + u

### 11.3 实现架构

`
┌─────────────────────────────────────────┐
│           RSM Cluster A                 │
│  ┌─────┐ ┌─────┐ ┌─────┐              │
│  │Node1│ │Node2│ │Node3│              │
│  └──┬──┘ └──┬──┘ └──┬──┘              │
│     └───────┼───────┘                  │
│             │ C3B Layer               │
│             ▼                          │
│  ┌─────────────────────────────────┐   │
│  │      Cross-Cluster Link         │   │
│  └─────────────────────────────────┘   │
└───────────────────┬─────────────────────┘
                    │
┌───────────────────┼─────────────────────┐
│           RSM Cluster B                 │
│  ┌─────────────────────────────────┐   │
│  │      Cross-Cluster Link         │   │
│  └─────────────────────────────────┘   │
│             │                          │
│             ▼                          │
│  ┌─────┐ ┌─────┐ ┌─────┐              │
│  │Node1│ │Node2│ │Node3│              │
│  └─────┘ └─────┘ └─────┘              │
└─────────────────────────────────────────┘
`

---

## 12. T2C 框架详解

### 12.1 测试到检查器转换

**输入**: 单元测试用例
`go
func TestAccountTransfer(t *testing.T) {
    account := NewAccount(100)
    account.Withdraw(30)
    assert.Equal(t, 70, account.Balance())
}
`

**输出**: 运行时检查器
`go
func CheckAccountTransferInvariant(account Account) bool {
    return account.Balance() >= 0
}
`

### 12.2 泛化策略

1. **常量泛化**: 100 → parameter
2. **操作序列**: Withdraw → any operation
3. **断言提取**: Balance() = 70 → Balance() >= 0

### 12.3 部署架构

`
┌─────────────────────────────────────────┐
│           Production Service            │
│                                         │
│  ┌─────────┐    ┌─────────────────┐    │
│  │  Main   │───►│ T2C Checker     │    │
│  │ Service │    │ (Sidecar)       │    │
│  └─────────┘    └─────────────────┘    │
│                         │              │
│                         ▼              │
│                  ┌─────────────┐       │
│                  │ Alert/Log   │       │
│                  └─────────────┘       │
└─────────────────────────────────────────┘
`

---

## 13. CAC 算法详解

### 13.1 Contention-Aware Cooperation

**核心思想**: 动态调整提议者集合大小以适应竞争程度

**参数**:

- d: 提议者数量 (动态)
- l: 被接受的值数量 (1 ≤ l ≤ d)

### 13.2 算法流程

`
每个进程 p:

1. 提出值 v_p
2. 等待其他提议
3. 使用预言机预测可接受的值
4. 如果预言机准确: 快速接受
5. 否则: 回退到标准共识
`

### 13.3 Cascading Consensus

**应用**: 低竞争环境下的快速共识

**性能**:

- 低竞争: O(1) 轮
- 高竞争: O(log n) 轮
- 最坏情况: O(n) 轮

---

## 14. 研究影响评估

### 14.1 对行业的影响

| 研究 | 短期影响(1-2年) | 长期影响(3-5年) |
|------|----------------|----------------|
| Basilisk | 验证工具采用 | 形式化验证标准化 |
| Picsou | 跨集群通信 | 分布式系统互联 |
| T2C | 生产测试集成 | 自愈系统 |
| CAC | 高性能共识 | 去中心化基础设施 |

### 14.2 开源实现

- **Basilisk**: github.com/basilisk-verifier (预计)
- **Picsou**: 集成到etcd/consul (可能)
- **T2C**: Go/Rust库开发中
- **CAC**: 研究原型

---

## 15. 未来研究方向

### 15.1 形式化验证

- 机器学习辅助不变量发现
- 自动反例生成
- 组合验证技术

### 15.2 跨域通信

- 标准化C3B协议
- 异构共识互操作
- WAN优化技术

### 15.3 生产测试

- 静默故障检测
- 自动修复系统
- 混沌工程集成

---

## 16. 扩展参考文献

1. Basilisk: Provenance Invariants Paper (OSDI 2025)
2. Picsou: C3B Protocol Specification
3. T2C: Test-to-Checker Transformation
4. CAC: Contention-Aware Cooperation (OPODIS 2025)
5. Distributed Systems: Principles and Paradigms (Textbook)

---

*Extended: 2026-04-03*
