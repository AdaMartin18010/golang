# FT-007: 拜占庭容错 (Byzantine Fault Tolerance)

> **维度**: Formal Theory
> **级别**: S (20+ KB)
> **标签**: #bft #pbft #byzantine-generals #consensus
> **权威来源**: [Practical Byzantine Fault Tolerance](http://pmg.csail.mit.edu/papers/osdi99.pdf) - Castro & Liskov, [The Byzantine Generals Problem](https://lamport.azurewebsites.net/pubs/byz.pdf) - Lamport

---

## 拜占庭将军问题

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     The Byzantine Generals Problem                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  场景：拜占庭军队围攻城市，将军们通过信使通信                                   │
│                                                                              │
│  问题：                                                                        │
│  ├─ 部分将军可能是叛徒（发送矛盾信息）                                          │
│  ├─ 信使可能被截获或替换                                                       │
│  └─ 忠诚将军需要达成一致的决策                                                  │
│                                                                              │
│  拜占庭容错：系统在 f 个叛徒存在时仍能正常工作                                   │
│                                                                              │
│  最小节点数：n ≥ 3f + 1                                                        │
│  • 3f+1 个节点可以容忍 f 个拜占庭节点                                          │
│  • 少于 3f+1 时，无法区分叛徒的欺骗                                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 故障类型对比

| 故障类型 | 行为 | 例子 | 容忍方案 |
|---------|------|------|---------|
| **崩溃故障** (Crash) | 停止响应 | 进程崩溃 | Raft/Paxos |
| **omission 故障** | 偶尔丢消息 | 网络丢包 | 重试机制 |
| **拜占庭故障** | 任意行为，包括恶意 | 黑客攻击、软件bug | BFT算法 |

---

## PBFT (Practical Byzantine Fault Tolerance)

### 算法概述

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          PBFT Protocol                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Client ──► Primary ──► Replicas                                             │
│                                                                              │
│  Phase 1: REQUEST                                                            │
│  Client ──request──► Primary                                                 │
│                                                                              │
│  Phase 2: PRE-PREPARE                                                        │
│  Primary ──pre-prepare──► All Replicas                                       │
│  (sequence number, digest)                                                   │
│                                                                              │
│  Phase 3: PREPARE                                                            │
│  Each Replica ──prepare──► All Replicas                                      │
│                                                                              │
│  Phase 4: COMMIT                                                             │
│  Each Replica ──commit──► All Replicas                                       │
│  (收到 2f+1 个 prepare 后)                                                     │
│                                                                              │
│  Phase 5: REPLY                                                              │
│  Replicas ──reply──► Client                                                  │
│  (收到 2f+1 个 commit 后执行)                                                  │
│                                                                              │
│  View Change (主节点故障时):                                                  │
│  1. 超时触发 view change                                                      │
│  2. 选举新主节点                                                              │
│  3. 同步状态                                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 安全性证明

$$
\begin{aligned}
&\text{Theorem: PBFT ensures safety with } n \geq 3f + 1 \\
&\text{Proof Sketch:} \\
&1. \text{Two non-faulty nodes accept different values} \\
&   \Rightarrow \text{Requires } 2f + 1 \text{ nodes for each value} \\
&   \Rightarrow \text{Total } 4f + 2 > 3f + 1 \text{ (contradiction)} \\
&2. \text{Non-faulty nodes agree on sequence} \\
&   \Rightarrow \text{Pre-prepare phase establishes ordering} \\
&   \Rightarrow \text{Prepare/Certify ensures agreement}
\end{aligned}
$$

---

## 现代 BFT 算法

### HotStuff (Libra/Diem)

- 线性通信复杂度 (O(n))
- 链式结构 (类似区块链)
- 视图变更优化

### Tendermint (Cosmos)

- BFT + PoS 结合
- 两阶段提交 (简化 PBFT)
- 即时最终性

### 性能对比

| 算法 | 通信轮次 | 复杂度 | 应用场景 |
|------|---------|--------|---------|
| PBFT | 3 | O(n²) | 许可链 |
| HotStuff | 3 | O(n) | Libra |
| Tendermint | 2 | O(n²) | Cosmos |
| Raft | 2 | O(n) | 非拜占庭 |

---

## BFT vs CFT

| 特性 | BFT (PBFT) | CFT (Raft) |
|------|-----------|-----------|
| 容错类型 | 任意故障 | 崩溃故障 |
| 最小节点 | 3f+1 | 2f+1 |
| 性能 | 较低 (~10K TPS) | 较高 (~100K TPS) |
| 应用场景 | 公链、跨组织 | 企业内网 |
| 实现复杂度 | 高 | 低 |

---

## 参考文献

1. [Practical Byzantine Fault Tolerance](http://pmg.csail.mit.edu/papers/osdi99.pdf) - Castro & Liskov, 1999
2. [The Byzantine Generals Problem](https://lamport.azurewebsites.net/pubs/byz.pdf) - Lamport, 1982
3. [HotStuff: BFT Consensus in the Lens of Blockchain](https://arxiv.org/abs/1803.05069) - Yin et al.
4. [Tendermint: Consensus without Mining](https://tendermint.com/static/docs/tendermint.pdf)
