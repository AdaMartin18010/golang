# FT-007: Multi-Paxos 的形式化理论与实践 (Multi-Paxos: Formal Theory & Practice)

> **维度**: Formal Theory
> **级别**: S (22+ KB)
> **标签**: #multi-paxos #consensus #log-replication #distributed-systems #optimization
> **权威来源**:
>
> - [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf) - Lamport (2001)
> - [Chubby: The Lock Service](https://research.google/pubs/chubby-the-lock-service-for-loosely-coupled-distributed-systems/) - Burrows (2006)
> - [Paxos Made Live](https://research.google/pubs/paxos-made-live-an-engineering-perspective/) - Chandra et al. (2007)
> - [Raft: Understandable Consensus](https://raft.github.io/raft.pdf) - Ongaro & Ousterhout (2014)
> - [Flexible Paxos](https://arxiv.org/abs/1608.06696) - Howard et al. (2016)

---

## 1. 从 Paxos 到 Multi-Paxos

### 1.1 基础 Paxos 的局限

**问题 1: 每个值都需要两阶段**

基础 Paxos 流程（每个值）：

```
Client → Proposer: Request
Proposer → Acceptors: Phase 1 (Prepare)
Acceptors → Proposer: Phase 1 (Promise)
Proposer → Acceptors: Phase 2 (AcceptRequest)
Acceptors → Proposer: Phase 2 (Accepted)
Proposer → Client: Response
```

延迟: **4 RTT** (Client-Proposer 往返 + Paxos 两阶段)

**问题 2: 并发冲突**
多个 Proposer 同时尝试提出值，导致 ballot 号竞争，可能活锁。

### 1.2 Multi-Paxos 核心思想

**定义 1.1 (Multi-Paxos)**
Multi-Paxos 是 Paxos 的优化变体，通过选举**稳定 Leader** 来：

1. 跳过 Phase 1（对连续提案复用 Promise）
2. 批量提交多个值（日志复制）
3. 将延迟从 4 RTT 降低到 **2 RTT**

**定理 1.1 (Multi-Paxos 优化)**

在稳定 Leader 场景下，Multi-Paxos 的**均摊延迟**为：

$$L_{MultiPaxos} = \frac{2 \cdot RTT_{steady} + 4 \cdot RTT_{leader\_change}}{n_{requests}}$$

其中 $n_{requests}$ 是连续请求数。

---

## 2. Multi-Paxos 形式化

### 2.1 系统模型

**定义 2.1 (Multi-Paxos 系统)**
系统 $\mathcal{M}$ 是八元组 $\langle \Pi, L, \mathcal{I}, \mathcal{V}, \mathcal{L}, Q, T, B \rangle$：

- $\Pi = \{p_1, ..., p_n\}$: 进程集合
- $L \in \Pi \cup \{\bot\}$: Leader
- $\mathcal{I} = \mathbb{N}^+$: 日志索引集合
- $\mathcal{V}$: 值域
- $\mathcal{L}: \Pi \times \mathcal{I} \rightarrow (\mathbb{N} \times \mathcal{V}) \cup \{\bot\}$: 日志
- $Q \subseteq 2^\Pi$: 法定人数集合
- $T \in \mathbb{N}$: 当前任期 (term/epoch)
- $B: \Pi \rightarrow \mathbb{N}$: 每个进程的最后 ballot

### 2.2 日志结构

**定义 2.2 (日志条目)**

$$\text{LogEntry} = \langle \text{index}: \mathcal{I}, \text{term}: \mathbb{N}, \text{value}: \mathcal{V}, \text{state}: \{\text{Pending}, \text{Accepted}, \text{Committed}\} \rangle$$

**定义 2.3 (日志完整性)**

$$\forall p_i, p_j \in \Pi, \forall k \in \mathcal{I}:$$
$$(\mathcal{L}(p_i)[k].\text{term} = \mathcal{L}(p_j)[k].\text{term}) \Rightarrow \mathcal{L}(p_i)[k].\text{value} = \mathcal{L}(p_j)[k].\text{value}$$

### 2.3 Leader 状态

**定义 2.4 (Leader 状态变量)**

| 变量 | 类型 | 说明 |
|------|------|------|
| $isLeader$ | $\mathbb{B}$ | 是否自认为是 Leader |
| $leaderTerm$ | $\mathbb{N}$ | Leader 任期 |
| $nextIndex$ | $\Pi \rightarrow \mathcal{I}$ | 每个 follower 的下一个发送索引 |
| $matchIndex$ | $\Pi \rightarrow \mathcal{I}$ | 每个 follower 已复制的最高索引 |
| $commitIndex$ | $\mathcal{I}$ | 已提交的最高索引 |
| $prepared$ | $\mathbb{B}$ | 是否已完成 Phase 1 |

### 2.4 状态转换

**转换 1: Leader 选举 (仅 Leader 变更时执行)**

```
Phase 1a: Leader → Acceptors: Prepare(T)
Phase 1b: Acceptors → Leader: Promise(T, maxAccepted)

条件: 收到多数派 Promise
动作:
    - 设置 prepared = true
    - 对每个索引 i，选择最高 term 的已接受值
    - 恢复未提交的日志条目
```

**转换 2: 日志复制 (稳定状态下)**

```
条件: prepared = true ∧ isLeader = true
动作:
    Leader → Followers: AppendEntries(T, prevIndex, prevTerm, entries[], commitIndex)
    Followers → Leader: Success/Failure

    如果收到多数派 Success:
        commitIndex = max(commitIndex, 成功复制的索引)
```

**转换 3: 提交确认**

```
条件: entry[i].state = Accepted ∧ 复制到多数派
动作: entry[i].state = Committed
```

---

## 3. 正确性证明

### 3.1 安全属性

**定理 3.1 (日志一致性)**
如果日志条目 $e$ 在索引 $i$ 被提交，则所有更高任期的 Leader 在索引 $i$ 都有 $e$。

*证明*:

1. 设 $e$ 在任期 $T$ 被提交于索引 $i$
2. 提交意味着 $e$ 被复制到多数派 $Q_{commit}$
3. 新 Leader $L'$ 必须获得多数派 $Q_{elect}$ 的选票
4. 由 Quorum 交集：$Q_{commit} \cap Q_{elect} \neq \emptyset$
5. 设 $p \in Q_{commit} \cap Q_{elect}$
6. $p$ 投票给 $L'$ 的条件：$L'$ 的日志至少和 $p$ 一样新
7. $p$ 有 $e$ 在索引 $i$，所以 $L'$ 也有 $e$
8. Leader 从不覆盖已提交条目

$\square$

**定理 3.2 (状态机安全)**
如果节点应用了索引 $i$ 的条目，则没有其他节点会在 $i$ 应用不同条目。

*证明*:
直接由定理 3.1 和提交定义得出。$\square$

### 3.2 活性属性

**定理 3.3 (Multi-Paxos 活性)**
在部分同步网络中，如果 Leader 稳定且多数派存活，则日志会无限增长。

*证明概要*:

1. Leader 无需执行 Phase 1（已 prepared）
2. 每条日志条目只需 Phase 2
3. 在同步期，消息在有限时间内送达
4. Leader 收到多数派确认即提交
5. 因此日志无限增长

$\square$

---

## 4. 多元表征

### 4.1 Multi-Paxos 概念地图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Multi-Paxos Concept Network                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Multi-Paxos                                  │   │
│  └───────────────────────────────┬─────────────────────────────────────┘   │
│                                  │                                         │
│              ┌───────────────────┼───────────────────┐                     │
│              ▼                   ▼                   ▼                     │
│       ┌───────────┐       ┌───────────┐       ┌───────────┐               │
│       │   Leader  │       │   Log     │       │   Quorum  │               │
│       │  Election │       │ Replication│      │  Consensus│               │
│       └─────┬─────┘       └─────┬─────┘       └─────┬─────┘               │
│             │                   │                   │                       │
│       ┌─────┴─────┐       ┌─────┴─────┐     ┌─────┴─────┐                 │
│       │           │       │           │     │           │                 │
│       │ • Phase 1 │       │ • Phase 2 │     │ • Prepare │                 │
│       │   once    │       │   skip 1  │     │ • Accept  │                 │
│       │ • Stable  │       │ • Batch   │     │ • Majority│                 │
│       │   Leader  │       │ • Pipeline│     │           │                 │
│       └───────────┘       └───────────┘     └───────────┘                 │
│                                                                              │
│  Optimization Hierarchy:                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ Level 1: Skip Phase 1 (Stable Leader)                               │   │
│  │ Level 2: Pipeline Requests                                          │   │
│  │ Level 3: Batch Multiple Entries                                     │   │
│  │ Level 4: Leader Leases (避免 Phase 1)                               │   │
│  │ Level 5: Flexible Paxos (优化 Quorum)                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Trade-offs:                                                                 │
│  ├── Latency: 4 RTT → 2 RTT (10x improvement with batching)                │
│  ├── Complexity: Higher than basic Paxos                                   │
│  └── Leader Bottleneck: All writes go through Leader                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 优化决策树

```
优化 Multi-Paxos 性能?
│
├── 延迟敏感? (需要最低延迟)
│   ├── 是 → 启用 Pipeline
│   │       └── Leader 无需等待确认即可发送下一个
│   │
│   └── 否 → 标准模式
│
├── 高吞吐需求?
│   ├── 是 → 启用 Batch
│   │       └── 合并多个请求为一个 AppendEntries
│   │
│   └── 否 → 单条模式
│
├── Leader 变更频繁?
│   ├── 是 → Leader Lease
│   │       └── 延长 Leader 任期，减少选举
│   │
│   └── 否 → 标准选举
│
├── 读操作多?
│   ├── 是 → Read Index / Lease Read
│   │       └── Leader 本地处理读，不复制
│   │
│   └── 否 → 标准读 (通过日志)
│
└── 跨机房部署?
    ├── 是 → Witness / Flexible Paxos
    │       └── 减少跨机房消息
    │
    └── 否 → 标准 Quorum
```

### 4.3 Paxos 变体对比矩阵

| 属性 | Basic Paxos | Multi-Paxos | Fast Paxos | Flexible Paxos | EPaxos |
|------|-------------|-------------|------------|----------------|--------|
| **正常情况延迟** | 2 RTT | 1 RTT | 1 RTT | 1-2 RTT | 1 RTT |
| **Leader 变更** | 2 RTT | 2 RTT | 2 RTT | 2 RTT | 1 RTT |
| **消息复杂度** | $O(n^2)$ | $O(n)$ | $O(n)$ | $O(n)$ | $O(n)$ |
| **Leader 瓶颈** | 无 | 有 | 有 | 有 | 无 |
| **冲突处理** | 重试 | 序列化 | 经典 Quorum | 经典 | 依赖图 |
| **实现复杂度** | 极高 | 高 | 高 | 高 | 极高 |
| **优化目标** | 基础 | 吞吐 | 延迟 | 灵活 | 无 Leader |
| **工业应用** | 理论 | Chubby, Spanner | - | - | - |

### 4.4 时序图: Multi-Paxos 流程

```
时间 →

场景 1: Leader 变更 (需要 Phase 1)

Node A                Node B                Node C                Node D
(New Leader)          (Follower)            (Follower)            (Follower)
   │                      │                      │                      │
   │ Phase 1a: Prepare(T=2)                    │                      │
   │─────────────────────►│─────────────────────►│                      │
   │                      │                      │                      │
   │ Phase 1b: Promise(T=2, maxAccepted)       │                      │
   │◄─────────────────────│◄─────────────────────│                      │
   │                      │                      │                      │
   │ (收到多数派，成为 Leader)                                            │
   │                      │                      │                      │
   │ Phase 2a: AppendEntries(T=2, entries)     │                      │
   │─────────────────────►│─────────────────────►│─────────────────────►│
   │                      │                      │                      │
   │ Phase 2b: Success    │◄─────────────────────│◄─────────────────────│
   │◄─────────────────────│                      │                      │
   │                      │                      │                      │

──────────────────────────────────────────────────────────────────────────────

场景 2: 稳定状态 (跳过 Phase 1)

Client                Leader                Followers
   │                      │                      │
   │ Request(v1)          │                      │
   │─────────────────────►│                      │
   │                      │                      │
   │                      │ AppendEntries(v1)    │
   │                      │─────────────────────►│
   │                      │─────────────────────►│
   │                      │                      │
   │                      │ Success              │
   │                      │◄─────────────────────│
   │                      │                      │
   │ Response             │                      │
   │◄─────────────────────│                      │
   │                      │                      │
   │ Request(v2)          │                      │
   │─────────────────────►│                      │
   │                      │ (无需 Phase 1!)      │
   │                      │                      │
   │                      │ AppendEntries(v2)    │
   │                      │─────────────────────►│
   │                      │─────────────────────►│
   │                      │                      │
   │ Response             │                      │
   │◄─────────────────────│                      │
   │                      │                      │

延迟: 2 RTT (Client→Leader + Leader→Followers→Leader)
```

---

## 5. TLA+ 形式化规约

```tla
------------------------------- MODULE MultiPaxos -----------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS
    Server,             \* 服务器集合
    Value,              \* 值域
    Term,               \* 任期集合
    Index,              \* 日志索引
    Nil

VARIABLES
    currentTerm,        \* 每个服务器的当前任期
    state,              \* 服务器状态 (Leader/Follower)
    log,                \* 日志条目: [term, value]
    commitIndex,        \* 已提交的最高索引

    \* Leader 状态
    nextIndex,          \* 对每个服务器，下一个发送的日志索引
    matchIndex,         \* 对每个服务器，已知的最高匹配索引
    prepared,           \* 是否已完成 Phase 1

    \* 消息
    msgs

vars == <<currentTerm, state, log, commitIndex, nextIndex, matchIndex,
          prepared, msgs>>

-----------------------------------------------------------------------------
\* 辅助定义

LogEntry == [term: Term, value: Value]

Quorum == {Q \in SUBSET Server: Cardinality(Q) * 2 > Cardinality(Server)}

\* 最后日志索引
LastLogIndex(s) == Len(log[s])

\* 最后日志任期
LastLogTerm(s) == IF LastLogIndex(s) = 0 THEN 0
                  ELSE log[s][LastLogIndex(s)].term

\* 日志是否至少一样新
LogIsAtLeastAsNewAs(s1, s2) ==
    LET lastTerm1 == LastLogTerm(s1)
        lastTerm2 == LastLogTerm(s2)
    IN /\ lastTerm1 > lastTerm2
       \/ /\ lastTerm1 = lastTerm2
          /\ LastLogIndex(s1) >= LastLogIndex(s2)

-----------------------------------------------------------------------------
\* Leader 选举 (Phase 1)

\* Leader 发送 Prepare
SendPrepare(i) ==
    /\ state[i] = "Leader"
    /\ ~prepared[i]
    /\ msgs' = msgs \cup {[type |-> "Prepare",
                          term |-> currentTerm[i],
                          from |-> i]}
    /\ UNCHANGED <<currentTerm, state, log, commitIndex,
                   nextIndex, matchIndex, prepared>>

\* Follower 响应 Promise
RecvPromise(i) ==
    \E m \in msgs:
        /\ m.type = "Prepare"
        /\ m.term >= currentTerm[i]
        /\ currentTerm' = [currentTerm EXCEPT ![i] = m.term]
        /\ state' = [state EXCEPT ![i] = "Follower"]
        /\ msgs' = msgs \cup {[type |-> "Promise",
                              term |-> m.term,
                              from |-> i,
                              log |-> log[i]]}
        /\ UNCHANGED <<log, commitIndex, nextIndex, matchIndex, prepared>>

\* Leader 处理 Promise，成为 prepared
BecomePrepared(i) ==
    /\ state[i] = "Leader"
    /\ ~prepared[i]
    /\ \E Q \in Quorum:
        LET promises == {m \in msgs: m.type = "Promise" /\ m.term = currentTerm[i]}
        IN /\ \A s \in Q: \E m \in promises: m.from = s
           /\ prepared' = [prepared EXCEPT ![i] = TRUE]
           /\ nextIndex' = [nextIndex EXCEPT ![i] =
               [j \in Server |-> Max({m.log[j]: m \in promises}) + 1]]
           /\ matchIndex' = [matchIndex EXCEPT ![i] =
               [j \in Server |-> 0]]
    /\ UNCHANGED <<currentTerm, state, log, commitIndex, msgs>>

-----------------------------------------------------------------------------
\* 日志复制 (Phase 2)

\* Leader 发送 AppendEntries
SendAppendEntries(i, j) ==
    /\ state[i] = "Leader"
    /\ prepared[i]
    /\ j \in Server \\ {i}
    /\ LET prevLogIndex == nextIndex[i][j] - 1
           prevLogTerm == IF prevLogIndex = 0 THEN 0
                          ELSE log[i][prevLogIndex].term
           entries == SubSeq(log[i], nextIndex[i][j], Len(log[i]))
       IN msgs' = msgs \cup {[type |-> "AppendEntries",
                              term |-> currentTerm[i],
                              from |-> i,
                              to |-> j,
                              prevLogIndex |-> prevLogIndex,
                              prevLogTerm |-> prevLogTerm,
                              entries |-> entries,
                              commitIndex |-> commitIndex[i]]}
    /\ UNCHANGED <<currentTerm, state, log, commitIndex,
                   nextIndex, matchIndex, prepared>>

\* Follower 处理 AppendEntries
RecvAppendEntries(i) ==
    \E m \in msgs:
        /\ m.type = "AppendEntries"
        /\ m.to = i
        /\ m.term >= currentTerm[i]
        /\ currentTerm' = [currentTerm EXCEPT ![i] = m.term]
        /\ state' = [state EXCEPT ![i] = "Follower"]
        /\ LET logOk == m.prevLogIndex = 0
                     \/ /\ m.prevLogIndex <= Len(log[i])
                        /\ log[i][m.prevLogIndex].term = m.prevLogTerm
           IN /\ logOk
              /\ log' = [log EXCEPT ![i] =
                  SubSeq(log[i], 1, m.prevLogIndex) \o m.entries]
              /\ commitIndex' = [commitIndex EXCEPT ![i] =
                  Min({m.commitIndex, Len(log'[i])})]
        /\ msgs' = msgs \cup {[type |-> "AppendResponse",
                              term |-> currentTerm[i],
                              from |-> i,
                              to |-> m.from,
                              success |-> logOk]}
        /\ UNCHANGED <<nextIndex, matchIndex, prepared>>

\* Leader 处理成功响应
HandleAppendSuccess(i) ==
    \E m \in msgs:
        /\ m.type = "AppendResponse"
        /\ m.to = i
        /\ m.term = currentTerm[i]
        /\ m.success
        /\ matchIndex' = [matchIndex EXCEPT ![i][m.from] =
            Max({@, m.prevLogIndex + Len(m.entries)})]
        /\ nextIndex' = [nextIndex EXCEPT ![i][m.from] =
            matchIndex'[i][m.from] + 1]
        /\ UNCHANGED <<currentTerm, state, log, commitIndex,
                       prepared, msgs>>

\* Leader 提交日志
AdvanceCommitIndex(i) ==
    /\ state[i] = "Leader"
    /\ prepared[i]
    /\ \E idx \in {commitIndex[i] + 1 .. Len(log[i])}:
        /\ log[i][idx].term = currentTerm[i]
        /\ \E Q \in Quorum:
            \A s \in Q: matchIndex[i][s] >= idx
        /\ commitIndex' = [commitIndex EXCEPT ![i] = idx]
    /\ UNCHANGED <<currentTerm, state, log, nextIndex,
                   matchIndex, prepared, msgs>>

-----------------------------------------------------------------------------
\* 安全属性

\* 状态机安全: 已提交的日志不会被覆盖
StateMachineSafety ==
    \A s1, s2 \in Server:
        \A idx \in 1 .. Min({Len(log[s1]), Len(log[s2])}):
            /\ idx <= commitIndex[s1]
            /\ idx <= commitIndex[s2]
            /\ log[s1][idx].term = log[s2][idx].term
            => log[s1][idx].value = log[s2][idx].value

=============================================================================
```

---

## 6. Go 实现

### 6.1 Multi-Paxos 核心实现

```go
package multipaxos

import (
    "context"
    "errors"
    "sync"
    "sync/atomic"
    "time"
)

// LogEntry 日志条目
type LogEntry struct {
    Index   uint64
    Term    uint64
    Command interface{}
}

// NodeState 节点状态
type NodeState int

const (
    Follower NodeState = iota
    Candidate
    Leader
)

// MultiPaxos 实现 Multi-Paxos 算法
type MultiPaxos struct {
    // 持久化状态
    currentTerm uint64
    votedFor    string
    log         []*LogEntry

    // 易失性状态
    state       NodeState
    commitIndex uint64
    lastApplied uint64

    // Leader 状态
    nextIndex   map[string]uint64
    matchIndex  map[string]uint64
    prepared    bool

    // 配置
    id          string
    peers       []string
    quorum      int

    mu          sync.RWMutex

    // 通道
    proposeCh   chan *ProposeRequest
    applyCh     chan *LogEntry

    // 网络
    transport   Transport

    ctx         context.Context
    cancel      context.CancelFunc
}

// Transport 定义网络传输接口
type Transport interface {
    Send(to string, msg Message) error
    Broadcast(msg Message) error
}

// Message 定义消息类型
type Message struct {
    Type         string // "Prepare", "Promise", "AppendEntries", "AppendResponse"
    From         string
    To           string
    Term         uint64

    // Prepare/Promise
    LastLogIndex uint64
    LastLogTerm  uint64

    // AppendEntries
    PrevLogIndex uint64
    PrevLogTerm  uint64
    Entries      []*LogEntry
    LeaderCommit uint64

    // Response
    Success      bool
}

// ProposeRequest 提案请求
type ProposeRequest struct {
    Command interface{}
    Result  chan error
}

// NewMultiPaxos 创建 Multi-Paxos 实例
func NewMultiPaxos(id string, peers []string, transport Transport) *MultiPaxos {
    ctx, cancel := context.WithCancel(context.Background())

    m := &MultiPaxos{
        id:         id,
        peers:      peers,
        quorum:     len(peers)/2 + 1,
        state:      Follower,
        log:        make([]*LogEntry, 0),
        nextIndex:  make(map[string]uint64),
        matchIndex: make(map[string]uint64),
        proposeCh:  make(chan *ProposeRequest, 100),
        applyCh:    make(chan *LogEntry, 100),
        transport:  transport,
        ctx:        ctx,
        cancel:     cancel,
    }

    // 启动主循环
    go m.run()

    return m
}

// Propose 提出命令
func (m *MultiPaxos) Propose(ctx context.Context, cmd interface{}) error {
    req := &ProposeRequest{
        Command: cmd,
        Result:  make(chan error, 1),
    }

    select {
    case m.proposeCh <- req:
    case <-ctx.Done():
        return ctx.Err()
    }

    select {
    case err := <-req.Result:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}

// run 主循环
func (m *MultiPaxos) run() {
    ticker := time.NewTicker(50 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-m.ctx.Done():
            return

        case req := <-m.proposeCh:
            m.handlePropose(req)

        case <-ticker.C:
            m.tick()
        }
    }
}

// handlePropose 处理提案
func (m *MultiPaxos) handlePropose(req *ProposeRequest) {
    m.mu.RLock()
    state := m.state
    prepared := m.prepared
    term := m.currentTerm
    m.mu.RUnlock()

    if state != Leader {
        req.Result <- errors.New("not leader")
        return
    }

    // 如果未 prepared，先执行 Phase 1
    if !prepared {
        if err := m.runPhase1(); err != nil {
            req.Result <- err
            return
        }
    }

    // 添加到本地日志
    m.mu.Lock()
    entry := &LogEntry{
        Index:   uint64(len(m.log)) + 1,
        Term:    term,
        Command: req.Command,
    }
    m.log = append(m.log, entry)
    m.mu.Unlock()

    // 复制到 followers
    go m.replicateEntry(entry, req.Result)
}

// runPhase1 执行 Phase 1
func (m *MultiPaxos) runPhase1() error {
    // 发送 Prepare
    msg := Message{
        Type: "Prepare",
        From: m.id,
        Term: m.currentTerm,
    }

    if err := m.transport.Broadcast(msg); err != nil {
        return err
    }

    // 等待 Promise (简化实现)
    // 实际实现需要异步收集响应

    m.mu.Lock()
    m.prepared = true
    for _, peer := range m.peers {
        m.nextIndex[peer] = uint64(len(m.log)) + 1
        m.matchIndex[peer] = 0
    }
    m.mu.Unlock()

    return nil
}

// replicateEntry 复制条目到 followers
func (m *MultiPaxos) replicateEntry(entry *LogEntry, result chan error) {
    m.mu.RLock()
    term := m.currentTerm
    commitIndex := m.commitIndex
    m.mu.RUnlock()

    // 发送 AppendEntries
    for _, peer := range m.peers {
        msg := Message{
            Type:         "AppendEntries",
            From:         m.id,
            To:           peer,
            Term:         term,
            PrevLogIndex: entry.Index - 1,
            Entries:      []*LogEntry{entry},
            LeaderCommit: commitIndex,
        }

        go m.transport.Send(peer, msg)
    }

    // 简化：假设成功
    // 实际实现需要收集多数派确认
    result <- nil
}

// HandleMessage 处理接收到的消息
func (m *MultiPaxos) HandleMessage(msg Message) error {
    switch msg.Type {
    case "Prepare":
        return m.handlePrepare(msg)
    case "Promise":
        return m.handlePromise(msg)
    case "AppendEntries":
        return m.handleAppendEntries(msg)
    case "AppendResponse":
        return m.handleAppendResponse(msg)
    default:
        return errors.New("unknown message type")
    }
}

// handlePrepare 处理 Prepare
func (m *MultiPaxos) handlePrepare(msg Message) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if msg.Term < m.currentTerm {
        return nil // 忽略旧 term
    }

    m.currentTerm = msg.Term
    m.state = Follower
    m.votedFor = ""

    // 发送 Promise
    response := Message{
        Type:         "Promise",
        From:         m.id,
        To:           msg.From,
        Term:         m.currentTerm,
        LastLogIndex: uint64(len(m.log)),
        LastLogTerm:  m.getLastLogTerm(),
    }

    return m.transport.Send(msg.From, response)
}

// handleAppendEntries 处理 AppendEntries
func (m *MultiPaxos) handleAppendEntries(msg Message) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if msg.Term < m.currentTerm {
        return nil
    }

    m.currentTerm = msg.Term
    m.state = Follower

    // 日志一致性检查
    if msg.PrevLogIndex > 0 {
        if msg.PrevLogIndex > uint64(len(m.log)) {
            // 日志不匹配
            return m.sendAppendResponse(msg.From, false)
        }
        if m.log[msg.PrevLogIndex-1].Term != msg.PrevLogTerm {
            // 日志不匹配
            return m.sendAppendResponse(msg.From, false)
        }
    }

    // 追加日志
    for i, entry := range msg.Entries {
        idx := msg.PrevLogIndex + uint64(i) + 1
        if idx <= uint64(len(m.log)) {
            if m.log[idx-1].Term != entry.Term {
                // 删除冲突条目及其后的所有条目
                m.log = m.log[:idx-1]
                m.log = append(m.log, entry)
            }
        } else {
            m.log = append(m.log, entry)
        }
    }

    // 更新 commitIndex
    if msg.LeaderCommit > m.commitIndex {
        m.commitIndex = min(msg.LeaderCommit, uint64(len(m.log)))
    }

    return m.sendAppendResponse(msg.From, true)
}

// sendAppendResponse 发送 AppendResponse
func (m *MultiPaxos) sendAppendResponse(to string, success bool) error {
    msg := Message{
        Type:    "AppendResponse",
        From:    m.id,
        To:      to,
        Term:    m.currentTerm,
        Success: success,
    }
    return m.transport.Send(to, msg)
}

// tick 定时任务
func (m *MultiPaxos) tick() {
    // 简化：实际的 Leader 选举超时处理
}

// 辅助函数
func (m *MultiPaxos) getLastLogTerm() uint64 {
    if len(m.log) == 0 {
        return 0
    }
    return m.log[len(m.log)-1].Term
}

func min(a, b uint64) uint64 {
    if a < b {
        return a
    }
    return b
}
```

---

## 7. 学术参考文献

### 7.1 核心论文

1. **Lamport, L. (2001)**. Paxos Made Simple. *ACM SIGACT News*, 32(4), 18-25.
   - Multi-Paxos 的基础

2. **Burrows, M. (2006)**. The Chubby Lock Service for Loosely-Coupled Distributed Systems. *OSDI*.
   - Google 的 Multi-Paxos 实现

3. **Chandra, T. D., Griesemer, R., & Redstone, J. (2007)**. Paxos Made Live: An Engineering Perspective. *PODC*.
   - 工程实践经验

### 7.2 优化研究

1. **Howard, H., et al. (2016)**. Flexible Paxos: Quorum Intersection Revisited. *OPODIS*.
   - Flexible Paxos 优化

2. **Moraru, I., Andersen, D. G., & Kaminsky, M. (2013)**. There is More Consensus in Egalitarian Parliaments. *SOSP*.
   - EPaxos: 无 Leader 优化

---

## 8. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Multi-Paxos Toolkit                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心记忆锚点: "稳定 Leader + 跳过 Phase 1"                                  │
│  ├── 首请求: 4 RTT (需要 Leader 选举 + Phase 1 + Phase 2)                    │
│  ├── 后续请求: 2 RTT (仅 Phase 2)                                            │
│  └── 批处理: 1 RTT 均摊 (多条日志一起复制)                                    │
│                                                                              │
│  关键优化:                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ 1. Pipeline: Leader 无需等待确认即可发送下一个                          │    │
│  │ 2. Batching: 合并多个请求为一个 AppendEntries                           │    │
│  │ 3. Leader Lease: 延长 Leader 任期，减少选举                              │    │
│  │ 4. Read Index: Leader 本地处理读请求                                    │    │
│  │ 5. Checksum: 减少日志比较开销                                           │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  常见误区:                                                                   │
│  ❌ "Multi-Paxos 就是 Raft" → Raft 是 Multi-Paxos 的重新设计                 │
│  ❌ "Phase 1 总是可以跳过" → Leader 变更时必须执行                           │
│  ❌ "Leader 永远不会变" → 网络分区时会触发 Leader 变更                       │
│  ❌ "所有节点都有完整日志" → 落后节点需要日志追赶                            │
│                                                                              │
│  性能公式:                                                                   │
│  ├── 延迟 = 2 × RTT (稳定状态)                                              │
│  ├── 吞吐 = BatchSize / RTT                                                │
│  └── Leader 瓶颈 = n × (CPU + Network)                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*文档生成时间: 2026-04-02*
*维护者: Distributed Systems Knowledge Base*
*版本: S-Level (22+ KB)*
