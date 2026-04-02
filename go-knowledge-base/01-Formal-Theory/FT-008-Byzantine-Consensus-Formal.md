# FT-008: 拜占庭共识的形式化理论 (Byzantine Consensus: Formal Theory)

> **维度**: Formal Theory
> **级别**: S (24+ KB)
> **标签**: #byzantine-fault-tolerance #pbft #consensus #blockchain #formal-verification
> **权威来源**:
>
> - [Practical Byzantine Fault Tolerance](http://pmg.csail.mit.edu/papers/osdi99.pdf) - Castro & Liskov (1999)
> - [The Byzantine Generals Problem](https://dl.acm.org/doi/10.1145/357172.357176) - Lamport, Shostak, Pease (1982)
> - [HotStuff: BFT Consensus in the Lens of Blockchain](https://arxiv.org/abs/1803.05069) - Yin et al. (2018)
> - [Tendermint: Byzantine Fault Tolerance](https://tendermint.com/static/docs/tendermint.pdf) - Kwon (2014)
> - [The Latest Gossip on BFT Consensus](https://arxiv.org/abs/1807.04938) - Buchman et al. (2018)

---

## 1. 拜占庭故障模型

### 1.1 故障分类

**定义 1.1 (拜占庭故障)**
拜占庭故障进程可能表现出**任意行为**：

$$\text{Byzantine}(p) \Rightarrow \forall o \in \text{Outputs}: p \text{ may output } o$$

包括：

- 停止响应
- 发送错误消息
- 发送矛盾消息给不同节点
- 与其他故障节点串通

**定义 1.2 (故障层次)**

```
Byzantine (任意行为) ───────────────────────── f < n/3
    │
    ├── Authentication-detectable Byzantine ── f < n/2 (带签名)
    │
    ├── Performance (性能故障) ─────────────── 可恢复
    │
    ├── Omission (遗漏故障) ────────────────── 重传机制
    │
    ├── Crash-Recovery (崩溃恢复) ──────────── f < n/2 + 持久化
    │
    └── Crash-Stop (崩溃停止) ──────────────── f < n/2
```

### 1.2 拜占庭将军问题

**定义 1.3 (拜占庭将军问题)**
$n$ 个将军需要就共同行动计划达成一致，其中 $f$ 个可能是叛徒：

$$
\begin{aligned}
&\text{IC1 (一致性)}: &&\text{所有忠诚将军采取相同行动} \\
&\text{IC2 (有效性)}: &&\text{如果指挥官忠诚，所有忠诚将军遵循其命令}
\end{aligned}
$$

**定理 1.1 (拜占庭容错下界 - Lamport et al. 1982)**
在同步网络中，拜占庭共识需要：

$$n \geq 3f + 1$$

*证明*:

假设 $n = 3$，$f = 1$。

1. 设忠诚节点 $L_1$, $L_2$，拜占庭节点 $B$
2. 指挥官是 $B$，发送矛盾消息：$v$ 给 $L_1$，$v'$ 给 $L_2$
3. $L_1$ 和 $L_2$ 交换观察：
   - $L_1$ 看到 ($v$, $v'$)
   - $L_2$ 看到 ($v'$, $v$)
4. $L_1$ 无法区分：
   - 场景 A: $B$ 是叛徒，$v$ 是真命令
   - 场景 B: $L_2$ 是叛徒，$v'$ 是真命令
5. 因此 $L_1$ 无法决定

$\square$

**定理 1.2 (异步网络下界)**
在异步网络中，即使 $n \geq 3f + 1$，确定性拜占庭共识也不可能 (FLP 扩展)。

---

## 2. PBFT 算法形式化

### 2.1 系统模型

**定义 2.1 (PBFT 系统)**
系统 $\mathcal{B}$ 是七元组 $\langle N, C, f, T, S, V, \Sigma \rangle$：

- $N = \{r_1, ..., r_n\}$: 副本集合，$n = 3f + 1$
- $C = \{c_1, c_2, ...\}$: 客户端集合
- $f$: 最大拜占庭故障数
- $T \in \mathbb{N}$: 视图 (view) 编号
- $S \in \{\text{Normal}, \text{ViewChange}\}$: 系统状态
- $V \subseteq N$: 视图中的主节点
- $\Sigma$: 数字签名方案

### 2.2 状态机

**定义 2.2 (副本状态)**

$$\text{ReplicaState} = \langle v, n, h, C, P, Q \rangle$$

- $v$: 当前视图编号
- $n$: 下一个序列号
- $h$: 低水标记 (已稳定检查点)
- $C$: 检查点集合
- $P$: prepared 证书集合
- $Q$: committed 证书集合

### 2.3 消息类型

**定义 2.3 (PBFT 消息)**

$$
\begin{aligned}
M ::= \quad &\langle \text{REQUEST}, o, t, c \rangle_{\sigma_c} &&\text{// 客户端请求} \\
    \mid \quad &\langle \text{PRE-PREPARE}, v, n, d, m \rangle_{\sigma_p} &&\text{// 主节点提案} \\
    \mid \quad &\langle \text{PREPARE}, v, n, d, i \rangle_{\sigma_i} &&\text{// 准备投票} \\
    \mid \quad &\langle \text{COMMIT}, v, n, d, i \rangle_{\sigma_i} &&\text{// 提交投票} \\
    \mid \quad &\langle \text{REPLY}, v, t, c, i, r \rangle_{\sigma_i} &&\text{// 执行结果} \\
    \mid \quad &\langle \text{VIEW-CHANGE}, ... \rangle &&\text{// 视图变更}
\end{aligned}
$$

其中 $d = H(m)$ 是消息摘要。

### 2.4 三阶段协议

**阶段 1: Pre-Prepare (主节点序列化)**

$$
\frac{\text{Primary}(p, v) \land \text{REQUEST}(m) \land n \not\in P}{\text{broadcast PRE-PREPARE}(v, n, H(m), m)}
$$

**阶段 2: Prepare (副本验证)**

$$
\frac{\text{PRE-PREPARE}(v, n, d, m)_{\sigma_p} \land \text{valid}(m)}{\text{broadcast PREPARE}(v, n, d, i)_{\sigma_i}}
$$

**阶段 3: Commit (共识达成)**

$$
\frac{|\{\text{PREPARE}(v, n, d, *)\}| \geq 2f}{\text{broadcast COMMIT}(v, n, d, i)_{\sigma_i}}
$$

**执行条件**:

$$
|\{\text{COMMIT}(v, n, d, *)\}| \geq 2f + 1 \Rightarrow \text{execute}(m)$$

---

## 3. 正确性证明

### 3.1 安全属性

**定理 3.1 (PBFT 一致性)**
所有非故障副本以相同顺序执行相同请求。

*证明*:

**引理 3.1 (Pre-Prepare 唯一性)**:
对于给定 $(v, n)$，最多一个非故障副本可以 pre-prepare 消息 $m$。

*证明*: 主节点是唯一的，它只能为一个序列号分配一个 digest。$\square$

**引理 3.2 (Prepare 传播)**:
如果非故障副本 $i$ 对 $(v, n, d)$ prepared，则至少有 $f + 1$ 个非故障副本也 prepared。

*证明*:
- Prepared 需要 $2f$ 个 prepare 消息（包括自己）
- 其中最多 $f$ 个来自拜占庭副本
- 因此至少 $f + 1$ 个来自非故障副本

$\square$

**主证明**:

假设副本 $i$ 和 $j$ 对 $(v, n)$ 执行不同请求 $m_i \neq m_j$。

1. $i$ 执行 $\Rightarrow i$ 收到 $2f + 1$ 个 commit
2. 其中至少 $f + 1$ 个来自非故障副本
3. 这些非故障副本必须已 prepared $(v, n, d_i)$
4. 由引理 3.2，$f + 1$ 个非故障副本 prepared 意味着至少 $2f + 1$ 个总副本 prepared
5. 同理，$j$ 执行意味着至少 $2f + 1$ 个副本 prepared $(v, n, d_j)$
6. 两个集合交集至少 $(2f + 1) + (2f + 1) - (3f + 1) = f + 1$ 个副本
7. 这些副本 prepared 了两个不同的 digest，矛盾

$\square$

### 3.2 活性属性

**定理 3.2 (PBFT 活性)**
在同步期内，客户端请求最终被所有非故障副本执行。

*证明概要*:

1. 客户端发送请求到所有副本
2. 至少 $2f + 1$ 个副本（包括主节点）收到请求
3. 主节点正常时，发起三阶段协议
4. 非故障副本响应，形成共识
5. 如果主节点故障，触发视图变更
6. 新主节点继续处理

$\square$

---

## 4. 多元表征

### 4.1 拜占庭共识概念地图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Byzantine Consensus Network                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                    ┌─────────────────────────┐                              │
│                    │   Byzantine Generals    │                              │
│                    │       Problem           │                              │
│                    └───────────┬─────────────┘                              │
│                                │                                            │
│              ┌─────────────────┼─────────────────┐                         │
│              ▼                 ▼                 ▼                         │
│       ┌───────────┐    ┌───────────┐    ┌───────────┐                     │
│       │  Synchronous│   │ Partially │    │ Asynchronous│                   │
│       │    n≥3f+1   │    │  Sync     │    │ Impossible  │                  │
│       └─────┬───────┘    └─────┬─────┘    └─────┬─────┘                   │
│             │                  │                │                          │
│       ┌─────┴───────┐    ┌─────┴───────┐      │                          │
│       │             │    │             │      │                          │
│       │ • PBFT      │    │ • Tendermint│      │                          │
│       │ • HotStuff  │    │ • Casper    │      │                          │
│       │ • Zyzzyva   │    │ • Algorand  │      │                          │
│       └─────────────┘    └─────────────┘      │                          │
│                                               │                          │
│  PBFT Protocol:                               │                          │
│  ┌─────────────────────────────────────────────────────────────────┐      │
│  │                                                                  │      │
│  │  Client ──REQUEST──► Replicas                                   │      │
│  │                          │                                      │      │
│  │                          ▼                                      │      │
│  │  ┌───────────────────────────────────────────────────────────┐  │      │
│  │  │ Phase 1: Pre-Prepare (Primary assigns sequence number)    │  │      │
│  │  │ Phase 2: Prepare (Replicas validate and vote)             │  │      │
│  │  │ Phase 3: Commit (Quorum of 2f+1 commits)                  │  │      │
│  │  │ Phase 4: Execute (Apply to state machine)                 │  │      │
│  │  └───────────────────────────────────────────────────────────┘  │      │
│  │                          │                                      │      │
│  │                          ▼                                      │      │
│  │  Client ◄──REPLY──── Replicas                                   │      │
│  │                                                                  │      │
│  └─────────────────────────────────────────────────────────────────┘      │
│                                                                              │
│  Key Properties:                                                             │
│  ├── Safety: n ≥ 3f + 1 (任何两个 2f+1 集合相交)                             │
│  ├── Liveness: View change protocol for primary replacement                  │
│  └── Communication: O(n²) messages per request                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 BFT 算法决策树

```
需要拜占庭容错?
│
├── 网络环境
│   ├── 同步网络 (已知延迟上界)
│   │   ├── 是 → PBFT / HotStuff
│   │   │       └── 确定性共识，最优性能
│   │   │
│   │   └── 否 → 检查部分同步假设
│   │
│   └── 部分同步 (最终同步)
│       ├── 是 → Tendermint / Casper FFG
│       │       └── 基于区块的 BFT，适合区块链
│       │
│       └── 否 → 异步网络
│               └── 需要随机化 (HoneyBadgerBFT)
│
├── 性能要求
│   ├── 高吞吐 + 低延迟?
│   │   ├── 是 → HotStuff (线性通信)
│   │   │       └── O(n) 消息复杂度
│   │   │
│   │   └── 否 → PBFT (经典)
│   │           └── O(n²) 消息复杂度
│   │
│   └── 可扩展性优先 (大规模网络)?
│       ├── 是 → Algorand / HoneyBadgerBFT
│       │       └── 随机化 + 加密排序
│       │
│       └── 否 → 经典 BFT
│
└── 应用场景
    ├── 联盟链/许可链 → PBFT, Tendermint
    ├── 公链/加密货币 → HotStuff (Libra), Tendermint (Cosmos)
    ├── 分布式存储 → Zyzzyva (优化读)
    └── 状态机复制 → BFT-SMaRt, UpRight
```

### 4.3 BFT 算法对比矩阵

| 属性 | PBFT | HotStuff | Tendermint | Zyzzyva | HoneyBadgerBFT |
|------|------|----------|------------|---------|----------------|
| **消息复杂度** | $O(n^2)$ | $O(n)$ | $O(n^2)$ | $O(n)$ | $O(n^2)$ |
| **延迟** | 3 RTT | 3 RTT | 3 RTT | 1.5 RTT | 异步 |
| **视图变更** | $O(n^3)$ | $O(n)$ | $O(n^2)$ | $O(n)$ | 无需 |
| **网络假设** | 部分同步 | 部分同步 | 部分同步 | 同步 | 异步 |
| **确定性** | 是 | 是 | 是 | 是 | 否 (随机化) |
| **可扩展性** | 低 (< 20) | 中 (< 100) | 中 (< 100) | 低 (< 20) | 高 (1000+) |
| **实现复杂度** | 高 | 中 | 中 | 高 | 高 |
| **应用场景** | 联盟链 | 公链 | 跨链 | 存储 | 大规模 |

### 4.4 PBFT 时序图

```
时间 →

Client       Primary (Replica 0)    Replica 1    Replica 2    Replica 3 (Byzantine)
   │               │                   │             │              │
   │───────────────REQUEST────────────►│────────────►│─────────────►│
   │               │                   │             │              │
   │               │──────────PRE-PREPARE────────────►│─────────────►│
   │               │  (v, n, d, m)_σp                 │              │
   │               │                   │             │              │
   │               │◄─────────PREPARE─────────────────│              │
   │               │  (v, n, d, 1)_σ1  │             │              │
   │               │◄────────────────────────PREPARE──│              │
   │               │  (v, n, d, 2)_σ2  │             │              │
   │               │                   │             │              │
   │               │ (收集到 2f=2 个 PREPARE, prepared)
   │               │                   │             │              │
   │               │──────────COMMIT─────────────────►│─────────────►│
   │               │                   │             │              │
   │               │◄─────────COMMIT──────────────────│              │
   │               │◄────────────────────────COMMIT───│              │
   │               │                   │             │              │
   │               │ (收集到 2f+1=3 个 COMMIT, committed-local)
   │               │                   │             │              │
   │               │ (执行请求, 更新状态机)            │              │
   │               │                   │             │              │
   │◄──────────────REPLY───────────────│─────────────│              │
   │  (v, t, c, 0, r)_σ0              │             │              │
   │               │                   │             │              │

消息计数: 1 (Request) + 3 (Pre-Prepare) + 6 (Prepare) + 6 (Commit) + 1 (Reply) = 17
故障容忍: f=1, n=4, 需要 3f+1=4 个副本
```

---

## 5. TLA+ 形式化规约

```tla
------------------------------- MODULE PBFT -------------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS
    Replica,            \* 副本集合
    Client,             \* 客户端集合
    Value,              \* 请求值域
    View,               \* 视图编号
    SeqNum,             \* 序列号
    MaxFaulty,          \* 最大故障数 f
    Nil

ASSUME
    /\ Cardinality(Replica) = 3 * MaxFaulty + 1
    /\ MaxFaulty >= 1

VARIABLES
    view,               \* 每个副本的当前视图
    seqNum,             \* 下一个序列号
    log,                \* 请求日志
    prepared,           \* prepared 证书集合
    committed,          \* committed 证书集合
    state,              \* 副本状态 (Normal/ViewChange)
    faulty,             \* 拜占庭故障副本集合
    msgs                \* 消息集合

vars == <<view, seqNum, log, prepared, committed, state, faulty, msgs>>

-----------------------------------------------------------------------------
\* 辅助定义

Quorum == {Q \in SUBSET Replica: Cardinality(Q) >= 2 * MaxFaulty + 1}

Primary(v) == CHOOSE r \in Replica: (v % Cardinality(Replica)) =
    (CHOOSE i \in 0..Cardinality(Replicas)-1:
        r = Replica[(v + i) % Cardinality(Replica)])

-----------------------------------------------------------------------------
\* 消息类型

RequestMsg == [type: {"REQUEST"}, client: Client, value: Value,
               timestamp: Nat, signature: Nat]

PrePrepareMsg == [type: {"PRE-PREPARE"}, view: View, seqNum: SeqNum,
                  digest: Nat, request: RequestMsg, from: Replica]

PrepareMsg == [type: {"PREPARE"}, view: View, seqNum: SeqNum,
               digest: Nat, from: Replica]

CommitMsg == [type: {"COMMIT"}, view: View, seqNum: SeqNum,
              digest: Nat, from: Replica]

ReplyMsg == [type: {"REPLY"}, view: View, timestamp: Nat, client: Client,
             replica: Replica, result: Value]

Message == RequestMsg \cup PrePrepareMsg \cup PrepareMsg \cup CommitMsg \cup ReplyMsg

-----------------------------------------------------------------------------
\* 动作定义

\* 客户端发送请求
ClientRequest(c) ==
    /\ msgs' = msgs \cup {[type |-> "REQUEST", client |-> c,
                          value |-> CHOOSE v \in Value: TRUE,
                          timestamp |-> seqNum, signature |-> c]}
    /\ UNCHANGED <<view, seqNum, log, prepared, committed, state, faulty>>

\* 主节点发送 Pre-Prepare
SendPrePrepare(r) ==
    /\ state[r] = "Normal"
    /\ Primary(view[r]) = r
    /\ \E m \in msgs:
        /\ m.type = "REQUEST"
        /\ ~\E pp \in msgs: pp.type = "PRE-PREPARE" /\ pp.request = m
        /\ msgs' = msgs \cup {[type |-> "PRE-PREPARE",
                              view |-> view[r],
                              seqNum |-> seqNum[r],
                              digest |-> Hash(m),
                              request |-> m,
                              from |-> r]}
        /\ seqNum' = [seqNum EXCEPT ![r] = @ + 1]
    /\ UNCHANGED <<view, log, prepared, committed, state, faulty>>

\* 副本发送 Prepare
SendPrepare(r) ==
    /\ state[r] = "Normal"
    /\ \E pp \in msgs:
        /\ pp.type = "PRE-PREPARE"
        /\ pp.view = view[r]
        /\ ValidRequest(pp.request)
        /\ ~\E p \in msgs: p.type = "PREPARE" /\ p.from = r /\
                         p.seqNum = pp.seqNum
        /\ msgs' = msgs \cup {[type |-> "PREPARE",
                              view |-> view[r],
                              seqNum |-> pp.seqNum,
                              digest |-> pp.digest,
                              from |-> r]}
    /\ UNCHANGED <<view, seqNum, log, prepared, committed, state, faulty>>

\* 副本检查 Prepared 条件
CheckPrepared(r) ==
    /\ state[r] = "Normal"
    /\ \E pp \in msgs:
        /\ pp.type = "PRE-PREPARE"
        /\ pp.view = view[r]
        /\ LET prepares == {p \in msgs: p.type = "PREPARE" /\
                            p.view = view[r] /\
                            p.seqNum = pp.seqNum /\
                            p.digest = pp.digest}
           IN Cardinality(prepares) >= 2 * MaxFaulty
        /\ prepared' = [prepared EXCEPT ![r] = @ \cup {pp.seqNum}]
    /\ UNCHANGED <<view, seqNum, log, committed, state, faulty, msgs>>

\* 副本发送 Commit
SendCommit(r) ==
    /\ state[r] = "Normal"
    /\ \E sn \in prepared[r]:
        /\ ~\E c \in msgs: c.type = "COMMIT" /\ c.from = r /\ c.seqNum = sn
        /\ msgs' = msgs \cup {[type |-> "COMMIT",
                              view |-> view[r],
                              seqNum |-> sn,
                              digest |-> LogDigest(r, sn),
                              from |-> r]}
    /\ UNCHANGED <<view, seqNum, log, prepared, committed, state, faulty>>

\* 副本检查 Committed 条件
CheckCommitted(r) ==
    /\ state[r] = "Normal"
    /\ \E sn \in SeqNum:
        /\ sn \in prepared[r]
        /\ LET commits == {c \in msgs: c.type = "COMMIT" /\
                           c.view = view[r] /\ c.seqNum = sn}
           IN Cardinality(commits) >= 2 * MaxFaulty + 1
        /\ committed' = [committed EXCEPT ![r] = @ \cup {sn}]
        /\ log' = [log EXCEPT ![r] = Append(@, GetRequest(view[r], sn))]
    /\ UNCHANGED <<view, seqNum, prepared, state, faulty, msgs>>

-----------------------------------------------------------------------------
\* 安全属性

\* 一致性: 所有非故障副本执行相同请求序列
Consistency ==
    \A r1, r2 \in Replica \\ faulty:
        \A i \in 1..Min({Len(log[r1]), Len(log[r2])}):
            log[r1][i] = log[r2][i]

\* 有效性: 执行的请求必须被客户端提出
Validity ==
    \A r \in Replica \\ faulty:
        \A i \in 1..Len(log[r]):
            \E m \in msgs: m.type = "REQUEST" /\ m = log[r][i]

=============================================================================
```

---

## 6. Go 实现

### 6.1 PBFT 核心实现

```go
package pbft

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "sync"
)

// Message 表示 PBFT 消息
type Message struct {
    Type     string // "REQUEST", "PRE-PREPARE", "PREPARE", "COMMIT", "REPLY"
    View     uint64
    SeqNum   uint64
    Digest   string

    // 请求相关
    Client    string
    Timestamp int64
    Value     interface{}

    // 签名
    From      string
    Signature []byte
}

// Replica 表示 PBFT 副本节点
type Replica struct {
    ID       string
    IsFaulty bool // 模拟拜占庭故障

    // 状态
    view       uint64
    seqNum     uint64
    state      string // "Normal", "ViewChange"

    // 日志
    log        []*Message

    // 证书集合
    prePrepares map[uint64]*Message  // seqNum -> PrePrepare
    prepares    map[uint64]map[string]*Message  // seqNum -> replicaID -> Prepare
    commits     map[uint64]map[string]*Message  // seqNum -> replicaID -> Commit

    // 状态追踪
    prepared   map[uint64]bool
    committed  map[uint64]bool
    executed   map[uint64]bool

    // 配置
    peers      []string
    f          int // 最大故障数
    quorum     int // 2f + 1

    mu         sync.RWMutex

    // 网络
    transport  Transport

    // 应用回调
    applyFunc  func(value interface{}) (interface{}, error)
}

// Transport 网络传输接口
type Transport interface {
    Broadcast(msg *Message) error
    Send(to string, msg *Message) error
}

// NewReplica 创建新的副本
func NewReplica(id string, peers []string, f int, transport Transport, applyFunc func(interface{}) (interface{}, error)) *Replica {
    return &Replica{
        ID:          id,
        view:        0,
        seqNum:      0,
        state:       "Normal",
        log:         make([]*Message, 0),
        prePrepares: make(map[uint64]*Message),
        prepares:    make(map[uint64]map[string]*Message),
        commits:     make(map[uint64]map[string]*Message),
        prepared:    make(map[uint64]bool),
        committed:   make(map[uint64]bool),
        executed:    make(map[uint64]bool),
        peers:       peers,
        f:           f,
        quorum:      2*f + 1,
        transport:   transport,
        applyFunc:   applyFunc,
    }
}

// IsPrimary 检查是否是主节点
func (r *Replica) IsPrimary() bool {
    primaryIndex := r.view % uint64(len(r.peers)+1)
    return r.ID == r.peers[primaryIndex]
}

// HandleMessage 处理消息
func (r *Replica) HandleMessage(msg *Message) error {
    if r.IsFaulty {
        // 拜占庭节点可能忽略或发送错误消息
        return nil
    }

    switch msg.Type {
    case "REQUEST":
        return r.handleRequest(msg)
    case "PRE-PREPARE":
        return r.handlePrePrepare(msg)
    case "PREPARE":
        return r.handlePrepare(msg)
    case "COMMIT":
        return r.handleCommit(msg)
    default:
        return errors.New("unknown message type")
    }
}

// handleRequest 处理客户端请求 (仅主节点)
func (r *Replica) handleRequest(msg *Message) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if !r.IsPrimary() {
        // 转发给主节点
        return nil
    }

    if r.state != "Normal" {
        return errors.New("not in normal state")
    }

    // 分配序列号
    r.seqNum++
    seqNum := r.seqNum

    // 计算摘要
    digest := r.digest(msg)

    // 创建 Pre-Prepare
    prePrepare := &Message{
        Type:      "PRE-PREPARE",
        View:      r.view,
        SeqNum:    seqNum,
        Digest:    digest,
        Client:    msg.Client,
        Timestamp: msg.Timestamp,
        Value:     msg.Value,
        From:      r.ID,
    }

    // 存储
    r.prePrepares[seqNum] = prePrepare
    r.prepares[seqNum] = make(map[string]*Message)
    r.commits[seqNum] = make(map[string]*Message)

    // 广播
    return r.transport.Broadcast(prePrepare)
}

// handlePrePrepare 处理 Pre-Prepare 消息
func (r *Replica) handlePrePrepare(msg *Message) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 验证
    if msg.View != r.view {
        return nil
    }

    if !r.validDigest(msg) {
        return nil
    }

    // 存储
    r.prePrepares[msg.SeqNum] = msg
    if r.prepares[msg.SeqNum] == nil {
        r.prepares[msg.SeqNum] = make(map[string]*Message)
    }
    if r.commits[msg.SeqNum] == nil {
        r.commits[msg.SeqNum] = make(map[string]*Message)
    }

    // 发送 Prepare
    prepare := &Message{
        Type:     "PREPARE",
        View:     r.view,
        SeqNum:   msg.SeqNum,
        Digest:   msg.Digest,
        From:     r.ID,
    }

    // 记录自己的 Prepare
    r.prepares[msg.SeqNum][r.ID] = prepare

    return r.transport.Broadcast(prepare)
}

// handlePrepare 处理 Prepare 消息
func (r *Replica) handlePrepare(msg *Message) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if msg.View != r.view {
        return nil
    }

    seqNum := msg.SeqNum

    // 存储
    if r.prepares[seqNum] == nil {
        r.prepares[seqNum] = make(map[string]*Message)
    }
    r.prepares[seqNum][msg.From] = msg

    // 检查 Prepared 条件
    if r.prepared[seqNum] {
        return nil
    }

    // 需要 Pre-Prepare + 2f 个 Prepare
    prePrepare := r.prePrepares[seqNum]
    if prePrepare == nil || prePrepare.Digest != msg.Digest {
        return nil
    }

    prepareCount := 0
    for _, p := range r.prepares[seqNum] {
        if p.Digest == msg.Digest {
            prepareCount++
        }
    }

    if prepareCount >= 2*r.f {
        r.prepared[seqNum] = true

        // 发送 Commit
        commit := &Message{
            Type:     "COMMIT",
            View:     r.view,
            SeqNum:   seqNum,
            Digest:   msg.Digest,
            From:     r.ID,
        }

        r.commits[seqNum][r.ID] = commit
        r.mu.Unlock()
        err := r.transport.Broadcast(commit)
        r.mu.Lock()
        return err
    }

    return nil
}

// handleCommit 处理 Commit 消息
func (r *Replica) handleCommit(msg *Message) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if msg.View != r.view {
        return nil
    }

    seqNum := msg.SeqNum

    // 存储
    if r.commits[seqNum] == nil {
        r.commits[seqNum] = make(map[string]*Message)
    }
    r.commits[seqNum][msg.From] = msg

    // 检查 Committed 条件
    if r.committed[seqNum] && r.executed[seqNum] {
        return nil
    }

    commitCount := 0
    for _, c := range r.commits[seqNum] {
        if c.Digest == msg.Digest {
            commitCount++
        }
    }

    if commitCount >= r.quorum {
        r.committed[seqNum] = true

        // 执行
        prePrepare := r.prePrepares[seqNum]
        if prePrepare != nil && !r.executed[seqNum] {
            r.executed[seqNum] = true

            // 添加到日志
            r.log = append(r.log, prePrepare)

            // 应用到状态机
            r.mu.Unlock()
            result, _ := r.applyFunc(prePrepare.Value)
            r.mu.Lock()

            // 如果是主节点，发送回复
            if r.IsPrimary() {
                reply := &Message{
                    Type:      "REPLY",
                    View:      r.view,
                    Timestamp: prePrepare.Timestamp,
                    Client:    prePrepare.Client,
                    Value:     result,
                    From:      r.ID,
                }
                r.mu.Unlock()
                r.transport.Send(prePrepare.Client, reply)
                r.mu.Lock()
            }
        }
    }

    return nil
}

// digest 计算消息摘要
func (r *Replica) digest(msg *Message) string {
    data := msg.Client + string(msg.Timestamp) + msg.Value.(string)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}

// validDigest 验证摘要
func (r *Replica) validDigest(msg *Message) bool {
    // 实际实现需要验证请求内容的摘要
    return true
}
```

---

## 7. 学术参考文献

### 7.1 经典论文

1. **Lamport, L., Shostak, R., & Pease, M. (1982)**. The Byzantine Generals Problem. *ACM Transactions on Programming Languages and Systems*, 4(3), 382-401.
   - 拜占庭将军问题的奠基性工作

2. **Castro, M., & Liskov, B. (1999)**. Practical Byzantine Fault Tolerance. *OSDI*.
   - PBFT 算法，第一个实用的 BFT 共识

3. **Castro, M., & Liskov, B. (2002)**. Practical Byzantine Fault Tolerance and Proactive Recovery. *ACM Transactions on Computer Systems*, 20(4), 398-461.
   - PBFT 的完整描述和主动恢复

### 7.2 现代 BFT 研究

1. **Yin, M., et al. (2019)**. HotStuff: BFT Consensus in the Lens of Blockchain. *arXiv:1803.05069*.
   - HotStuff: 线性通信复杂度的 BFT

2. **Buchman, K., Kwon, J., & Milosevic, Z. (2018)**. The Latest Gossip on BFT Consensus. *arXiv:1807.04938*.
   - Tendermint 共识算法

3. **Abraham, I., et al. (2017)**. Solida: A Blockchain Protocol Based on Reconfigurable Byzantine Consensus. *OPODIS*.
   - 可重配置 BFT 共识

---

## 8. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Byzantine Consensus Toolkit                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心记忆锚点: "n ≥ 3f + 1"                                                  │
│  ├── 需要 3f+1 个副本容忍 f 个拜占庭故障                                     │
│  ├── 任何两个 2f+1 集合必定相交 (f+1 个诚实节点)                             │
│  └── 相交节点保证诚实节点看到相同的消息                                      │
│                                                                              │
│  PBFT 三阶段:                                                                │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ Phase 1: Pre-Prepare                                                │    │
│  │   主节点分配序列号，确保请求的全序                                      │    │
│  │                                                                     │    │
│  │ Phase 2: Prepare                                                    │    │
│  │   副本投票，确保 2f+1 个节点 prepared 相同请求                         │    │
│  │                                                                     │    │
│  │ Phase 3: Commit                                                     │    │
│  │   副本确认提交，确保 2f+1 个节点 committed                             │    │
│  │                                                                     │    │
│  │ Execute: 应用请求到状态机                                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  常见误区:                                                                   │
│  ❌ "BFT 太慢不能用" → PBFT 可达 10k+ TPS                                   │
│  ❌ "BFT 只能用于区块链" → 传统系统也适用                                    │
│  ❌ "n=3f+1 太浪费" → 这是理论下界，无法更优                                │
│  ❌ "BFT 假设太强" → 公链/云环境都需要考虑拜占庭行为                          │
│                                                                              │
│  关键公式:                                                                   │
│  ├── 副本数: n ≥ 3f + 1                                                    │
│  ├── Quorum: 2f + 1                                                        │
│  ├── 消息复杂度: O(n²) 每请求 (PBFT)                                        │
│  └── 优化后: O(n) (HotStuff)                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*文档生成时间: 2026-04-02*
*维护者: Distributed Systems Knowledge Base*
*版本: S-Level (24+ KB)*
