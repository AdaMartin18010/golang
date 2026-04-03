# FT-009: 分布式事务的形式化理论 (Distributed Transactions: Formal Theory)

> **维度**: Formal Theory
> **级别**: S (23+ KB)
> **标签**: #distributed-transactions #2pc #3pc #saga #acid #consensus
> **权威来源**:
>
> - [Atomicity vs. Idempotence](https://dl.acm.org/doi/10.1145/1267960.1267961) - Gray (1981)
> - [Implementing Fault-Tolerant Services](https://dl.acm.org/doi/10.1145/2980.357399) - Oki & Liskov (1988)
> - [Consensus on Transaction Commit](https://dl.acm.org/doi/10.1145/235685.235699) - Gray & Lamport (2006)
> - [Sagas](https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf) - Garcia-Molina & Salem (1987)
> - [Calm Theorem](https://rise.cs.berkeley.edu/wp-content/uploads/2019/06/calm-conjecture.pdf) - Ameloot et al. (2013)

---

## 1. 分布式事务基础

### 1.1 ACID 属性的形式化

**定义 1.1 (分布式事务)**
分布式事务 $T$ 是操作序列跨越多个节点的原子工作单元：

$$T = \langle (p_1, o_1), (p_2, o_2), ..., (p_n, o_n) \rangle$$

其中 $p_i \in \Pi$ 是节点，$o_i$ 是操作。

**定义 1.2 (ACID 属性)**

$$
\begin{aligned}
&\text{Atomicity (原子性)}: &&T \text{ 要么全部执行，要么全部不执行} \\
&\text{Consistency (一致性)}: &&T \text{ 执行后系统处于有效状态} \\
&\text{Isolation (隔离性)}: &&\forall T_1, T_2: \text{并发执行等价于某种串行执行} \\
&\text{Durability (持久性)}: &&T \text{ 提交后，结果永久保存}
\end{aligned}
$$

**形式化原子性**:

$$\text{Atomicity}(T) \equiv (\forall o \in T: \text{executed}(o)) \oplus (\forall o \in T: \neg\text{executed}(o))$$

**形式化隔离性 (可串行化)**:

$$\exists \text{串行调度 } S: \text{并发执行} \equiv S$$

### 1.2 事务状态机

**定义 1.3 (事务状态)**

$$
\text{TransactionState} ::= \text{Active} \mid \text{Prepared} \mid \text{Committed} \mid \text{Aborted}
$$

**状态转换**:

$$
\begin{aligned}
&\text{Active} \xrightarrow{\text{Prepare}} \text{Prepared} \\
&\text{Prepared} \xrightarrow{\text{Commit}} \text{Committed} \\
&\text{Active} \xrightarrow{\text{Abort}} \text{Aborted} \\
&\text{Prepared} \xrightarrow{\text{Abort}} \text{Aborted}
\end{aligned}
$$

---

## 2. 两阶段提交 (2PC) 形式化

### 2.1 协议描述

**参与者角色**:

- **Coordinator**: 协调者，管理事务提交流程
- **Cohort**: 参与者，执行本地事务

**定义 2.1 (2PC 消息)**

$$
\begin{aligned}
M_{2PC} ::= \quad &\text{Prepare} \quad &&\text{// Phase 1: 询问准备} \\
    \mid \quad &\text{VoteYes} \quad &&\text{// 可以提交} \\
    \mid \quad &\text{VoteNo} \quad &&\text{// 不能提交} \\
    \mid \quad &\text{Commit} \quad &&\text{// Phase 2: 提交} \\
    \mid \quad &\text{Abort} \quad &&\text{// Phase 2: 中止} \\
    \mid \quad &\text{ACK} \quad &&\text{// 确认}
\end{aligned}
$$

### 2.2 状态转换系统

**Phase 1 (投票阶段)**:

```
Coordinator: INIT ──Prepare──► WAIT
                    │
                    ▼
Cohorts:    INIT ──VoteYes/No──► READY
```

**Phase 2 (决定阶段)**:

```
Coordinator: WAIT ──(all Yes)──► COMMIT ──Commit──► COMPLETE
                      │
                      └──(any No)──► ABORT ──Abort──► COMPLETE
```

### 2.3 正确性证明

**定理 2.1 (2PC 原子性)**
所有参与者最终要么全部提交，要么全部中止。

*证明*:

**情况 1**: Coordinator 决定 COMMIT

- 只有当收到所有 VoteYes 时才决定 COMMIT
- 发送 Commit 消息给所有参与者
- 所有参与者收到 Commit 后执行提交

**情况 2**: Coordinator 决定 ABORT

- 收到任何 VoteNo 或超时即决定 ABORT
- 发送 Abort 消息给所有参与者
- 所有参与者收到 Abort 后执行回滚

**一致性保证**:

- Coordinator 的日志记录决定
- 参与者收到决定前保持 Prepared 状态（持有锁）
- 任何故障恢复后根据 Coordinator 决定继续

$\square$

**定理 2.2 (2PC 阻塞问题)**
如果 Coordinator 在 Phase 2 故障，参与者可能无限阻塞。

*证明*:

参与者在 Prepared 状态时：

- 已投票 Yes
- 持有资源锁
- 等待 Coordinator 决定
- Coordinator 故障后，无法知道最终决定
- 不能单方面提交（可能其他参与者已中止）
- 不能单方面中止（可能 Coordinator 已提交）

$\square$

---

## 3. 三阶段提交 (3PC) 形式化

### 3.1 解决阻塞问题

**核心思想**: 引入预提交 (PreCommit) 状态，使参与者可以在协调者故障时独立决定。

**定义 3.1 (3PC 状态)**

$$
\text{3PCState} ::= \text{INIT} \mid \text{READY} \mid \text{PRECOMMIT} \mid \text{COMMIT} \mid \text{ABORT}
$$

### 3.2 协议流程

**Phase 1**: CanCommit? (类似 2PC Phase 1)

**Phase 2**: PreCommit

- Coordinator 收到所有 Yes 后发送 PreCommit
- 参与者进入 PRECOMMIT 状态

**Phase 3**: DoCommit

- Coordinator 收到所有 ACK 后发送 Commit
- 参与者提交并释放锁

### 3.3 非阻塞性质

**定理 3.1 (3PC 非阻塞性)**
在同步网络假设下（消息延迟有界），3PC 不会阻塞。

*证明*:

参与者超时处理：

1. **在 READY 状态超时**（未收到 PreCommit）:
   - 可以安全中止
   - Coordinator 不可能已决定提交

2. **在 PRECOMMIT 状态超时**（未收到 Commit）:
   - 联系其他参与者
   - 如果有人已提交，则提交
   - 如果协调者存活但无响应，可以提交（因为 PreCommit 意味着协调者已收到所有 Yes）

$\square$

**注意**: 3PC 的非阻塞性依赖于同步网络假设。在异步网络中，无法区分慢网络和故障，仍然可能阻塞。

---

## 4. Saga 模式形式化

### 4.1 长事务问题

**问题**: 2PC/3PC 在长时间事务中：

- 持有锁时间长，影响并发
- 协调者故障风险增加
- 不适合跨服务/微服务架构

### 4.2 Saga 形式化

**定义 4.1 (Saga)**
Saga 是一系列本地事务 $L_1, L_2, ..., L_n$，每个本地事务 $L_i$ 有对应的补偿事务 $C_i$：

$$\text{Saga} = \langle (L_1, C_1), (L_2, C_2), ..., (L_n, C_n) \rangle$$

**定义 4.2 (Saga 执行语义)**

$$
\text{Execute}(\text{Saga}) =
\begin{cases}
L_1 \circ L_2 \circ ... \circ L_n & \text{if all succeed} \\
L_1 \circ ... \circ L_k \circ C_k \circ ... \circ C_1 & \text{if } L_{k+1} \text{ fails}
\end{cases}
$$

**定理 4.1 (Saga 补偿性质)**
正确实现的补偿事务应满足：

$$C_i \circ L_i \equiv \text{Identity}$$

即补偿事务应该撤销本地事务的效果。

### 4.3 Saga 与 ACID

| 属性 | 2PC | Saga |
|------|-----|------|
| **原子性** | 强 (All-or-nothing) | 弱 (最终一致) |
| **一致性** | 强 | 补偿期间不一致 |
| **隔离性** | 可串行化 | 无隔离性 |
| **持久性** | 强 | 强 |

---

## 5. 多元表征

### 5.1 分布式事务概念地图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Distributed Transactions Network                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                         ┌─────────────────┐                                 │
│                         │  ACID Properties│                                 │
│                         └────────┬────────┘                                 │
│                                  │                                          │
│              ┌───────────────────┼───────────────────┐                      │
│              ▼                   ▼                   ▼                      │
│       ┌───────────┐       ┌───────────┐       ┌───────────┐                │
│       │ Atomicity │       │Consistency│       │ Isolation │                │
│       └─────┬─────┘       └─────┬─────┘       └─────┬─────┘                │
│             │                   │                   │                       │
│             ▼                   ▼                   ▼                       │
│       ┌─────────────────────────────────────────────────────┐              │
│       │            Distributed Commit Protocols              │              │
│       ├─────────────────────────────────────────────────────┤              │
│       │                                                      │              │
│       │  2PC ─────────────────────────────────────────────►  │              │
│       │  ├── Phase 1: Voting                                 │              │
│       │  ├── Phase 2: Decision                               │              │
│       │  └── Blocking Problem                                │              │
│       │                                                      │              │
│       │  3PC ─────────────────────────────────────────────►  │              │
│       │  ├── Phase 1: CanCommit?                             │              │
│       │  ├── Phase 2: PreCommit                              │              │
│       │  ├── Phase 3: DoCommit                               │              │
│       │  └── Non-blocking (sync network)                     │              │
│       │                                                      │              │
│       │  Saga ────────────────────────────────────────────►  │              │
│       │  ├── Local Transactions                              │              │
│       │  ├── Compensation                                    │              │
│       │  └── Eventual Consistency                            │              │
│       │                                                      │              │
│       └─────────────────────────────────────────────────────┘              │
│                                                                              │
│  Trade-offs:                                                                 │
│  ├── 2PC: Strong consistency, blocking, short transactions                   │
│  ├── 3PC: Better availability, complex, needs sync network                   │
│  └── Saga: Best availability, eventual consistency, microservices            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 事务协议决策树

```
需要分布式事务?
│
├── 事务执行时间
│   ├── 短事务 (< 秒级)?
│   │   ├── 是 → 2PC
│   │   │       └── 强一致性，标准实现
│   │   │
│   │   └── 否 → Saga / TCC
│   │           └── 长事务，避免长时间锁定
│   │
│   └── 长事务 (分钟/小时)?
│       └── Saga / Event Sourcing
│
├── 一致性要求
│   ├── 强一致必需?
│   │   ├── 是 → 2PC / 3PC
│   │   │       └── 考虑阻塞问题的缓解方案
│   │   │
│   │   └── 否 → Saga
│   │           └── 最终一致可接受
│   │
│   └── 最终一致足够?
│       └── Saga / 异步消息
│
├── 架构风格
│   ├── 单体应用?
│   │   └── 2PC (数据库原生支持)
│   │
│   ├── 微服务?
│   │   └── Saga (避免服务耦合)
│   │
│   └── 混合架构?
│       └── Saga + 内部 2PC
│
└── 失败处理
    ├── 可自动补偿?
    │   ├── 是 → Saga
    │   └── 否 → 2PC / 人工介入
    │
    └── 需要重试?
        └── Saga + 重试策略
```

### 5.3 协议对比矩阵

| 属性 | 2PC | 3PC | Saga | TCC |
|------|-----|-----|------|-----|
| **一致性** | 强 | 强 | 最终 | 最终 |
| **阻塞性** | 阻塞 | 非阻塞* | 非阻塞 | 非阻塞 |
| **协调者故障** | 可能阻塞 | 可恢复 | 可恢复 | 可恢复 |
| **消息复杂度** | $O(n)$ | $O(n)$ | $O(n)$ | $O(n)$ |
| **延迟** | 2 RTT | 3 RTT | 异步 | 异步 |
| **实现复杂度** | 低 | 高 | 中 | 高 |
| **适用场景** | 短事务 | 高可用 | 微服务 | 资源预留 |
| **持久化要求** | 协调者日志 | 状态机 | 事件日志 | 预留记录 |

*3PC 的非阻塞性依赖于同步网络假设

### 5.4 2PC 时序图

```
时间 →

Coordinator              Cohort 1              Cohort 2              Cohort 3
   │                       │                     │                     │
   │──────Prepare(T1)─────►│                     │                     │
   │──────Prepare(T2)──────┼────────────────────►│                     │
   │──────Prepare(T3)──────┼─────────────────────┼────────────────────►│
   │                       │                     │                     │
   │◄─────VoteYes(T1)──────│                     │                     │
   │◄─────VoteYes(T2)──────┼─────────────────────│                     │
   │◄─────VoteYes(T3)──────┼─────────────────────┼─────────────────────│
   │                       │                     │                     │
   │ (收到所有 Yes, 写决定日志)                                           │
   │                       │                     │                     │
   │──────Commit(T1)──────►│                     │                     │
   │──────Commit(T2)───────┼────────────────────►│                     │
   │──────Commit(T3)───────┼─────────────────────┼────────────────────►│
   │                       │                     │                     │
   │                       │ (执行提交)           │                     │
   │                       │                     │                     │
   │◄──────ACK(T1)─────────│                     │                     │
   │◄──────ACK(T2)─────────┼─────────────────────│                     │
   │◄──────ACK(T3)─────────┼─────────────────────┼─────────────────────│
   │                       │                     │                     │
   │ (收到所有 ACK, 事务完成)                                              │
   ▼                       ▼                     ▼                     ▼

失败场景:

场景 1: Cohort 2 投票 No
Coordinator ──Abort(T1,T2,T3)──► (所有参与者回滚)

场景 2: Coordinator 在 Commit 前故障
Cohorts 等待超时 → 联系 Coordinator → 根据日志恢复
若无法联系 → 参与者保持 Prepared 状态 (阻塞)

场景 3: Coordinator 在发送部分 Commit 后故障
已收到 Commit 的参与者提交
未收到 Commit 的参与者等待 (阻塞)
需要 Coordinator 恢复后继续
```

---

## 6. TLA+ 形式化规约

```tla
------------------------------- MODULE TwoPhaseCommit -------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS
    Coordinator,        \* 协调者
    Cohort,             \* 参与者集合
    Value,              \* 事务值域
    Nil

VARIABLES
    cState,             \* 协调者状态
    cVotes,             \* 协调者收集的投票
    cDecision,          \* 协调者决定

    hState,             \* 参与者状态
    hVote,              \* 参与者投票
    hDecision,          \* 参与者收到的决定

    msgs                \* 消息集合

tpcVars == <<cState, cVotes, cDecision, hState, hVote, hDecision, msgs>>

-----------------------------------------------------------------------------
\* 类型定义

CState == {"init", "waiting", "committed", "aborted"}
HState == {"working", "prepared", "committed", "aborted"}

TypeOK ==
    /\ cState \in CState
    /\ cVotes \in [Cohort -> {"", "yes", "no"}]
    /\ cDecision \in {"", "commit", "abort"}
    /\ hState \in [Cohort -> HState]
    /\ hVote \in [Cohort -> {"", "yes", "no"}]
    /\ hDecision \in [Cohort -> {"", "commit", "abort"}]

-----------------------------------------------------------------------------
\* 消息类型

PrepareMsg == [type: {"prepare"}]
VoteMsg == [type: {"vote"}, cohort: Cohort, vote: {"yes", "no"}]
DecisionMsg == [type: {"decision"}, decision: {"commit", "abort"}]
AckMsg == [type: {"ack"}, cohort: Cohort]

Message == PrepareMsg \cup VoteMsg \cup DecisionMsg \cup AckMsg

-----------------------------------------------------------------------------
\* 协调者动作

\* Phase 1: 发送 Prepare
C_SendPrepare ==
    /\ cState = "init"
    /\ msgs' = msgs \cup {[type |-> "prepare"]}
    /\ cState' = "waiting"
    /\ UNCHANGED <<cVotes, cDecision, hState, hVote, hDecision>>

\* 协调者收集投票
C_ReceiveVote ==
    /\ cState = "waiting"
    /\ \E m \in msgs:
        /\ m.type = "vote"
        /\ cVotes[m.cohort] = ""
        /\ cVotes' = [cVotes EXCEPT ![m.cohort] = m.vote]
    /\ UNCHANGED <<cState, cDecision, hState, hVote, hDecision, msgs>>

\* Phase 2a: 决定提交
C_DecideCommit ==
    /\ cState = "waiting"
    /\ \A h \in Cohort: cVotes[h] = "yes"
    /\ cState' = "committed"
    /\ cDecision' = "commit"
    /\ msgs' = msgs \cup {[type |-> "decision", decision |-> "commit"]}
    /\ UNCHANGED <<cVotes, hState, hVote, hDecision>>

\* Phase 2b: 决定中止
C_DecideAbort ==
    /\ cState = "waiting"
    /\ \E h \in Cohort: cVotes[h] = "no"
    /\ cState' = "aborted"
    /\ cDecision' = "abort"
    /\ msgs' = msgs \cup {[type |-> "decision", decision |-> "abort"]}
    /\ UNCHANGED <<cVotes, hState, hVote, hDecision>>

-----------------------------------------------------------------------------
\* 参与者动作

\* 收到 Prepare，投票 Yes
H_ReceivePrepareAndVoteYes(h) ==
    /\ hState[h] = "working"
    /\ \E m \in msgs: m.type = "prepare"
    /\ hState' = [hState EXCEPT ![h] = "prepared"]
    /\ hVote' = [hVote EXCEPT ![h] = "yes"]
    /\ msgs' = msgs \cup {[type |-> "vote", cohort |-> h, vote |-> "yes"]}
    /\ UNCHANGED <<cState, cVotes, cDecision, hDecision>>

\* 收到 Prepare，投票 No
H_ReceivePrepareAndVoteNo(h) ==
    /\ hState[h] = "working"
    /\ \E m \in msgs: m.type = "prepare"
    /\ hState' = [hState EXCEPT ![h] = "aborted"]
    /\ hVote' = [hVote EXCEPT ![h] = "no"]
    /\ msgs' = msgs \cup {[type |-> "vote", cohort |-> h, vote |-> "no"]}
    /\ UNCHANGED <<cState, cVotes, cDecision, hDecision>>

\* 收到 Commit 决定
H_ReceiveCommit(h) ==
    /\ hState[h] = "prepared"
    /\ \E m \in msgs:
        /\ m.type = "decision"
        /\ m.decision = "commit"
    /\ hState' = [hState EXCEPT ![h] = "committed"]
    /\ hDecision' = [hDecision EXCEPT ![h] = "commit"]
    /\ msgs' = msgs \cup {[type |-> "ack", cohort |-> h]}
    /\ UNCHANGED <<cState, cVotes, cDecision, hVote>>

\* 收到 Abort 决定
H_ReceiveAbort(h) ==
    /\ hState[h] \in {"working", "prepared"}
    /\ \E m \in msgs:
        /\ m.type = "decision"
        /\ m.decision = "abort"
    /\ hState' = [hState EXCEPT ![h] = "aborted"]
    /\ hDecision' = [hDecision EXCEPT ![h] = "abort"]
    /\ msgs' = msgs \cup {[type |-> "ack", cohort |-> h]}
    /\ UNCHANGED <<cState, cVotes, cDecision, hVote>>

-----------------------------------------------------------------------------
\* 安全属性

\* 原子性: 所有参与者最终状态一致
Atomicity ==
    \A h1, h2 \in Cohort:
        (hState[h1] \in {"committed", "aborted"} /\
         hState[h2] \in {"committed", "aborted"})
            => hState[h1] = hState[h2]

\* 一致性: 如果协调者决定提交，所有投票都是 Yes
Consistency ==
    cDecision = "commit" =>
        \A h \in Cohort: hVote[h] = "yes"

=============================================================================
```

---

## 7. Go 实现

### 7.1 2PC 协调者实现

```go
package dtx

import (
    "context"
    "errors"
    "fmt"
    "sync"
)

// TransactionStatus 事务状态
type TransactionStatus int

const (
    StatusPending TransactionStatus = iota
    StatusPrepared
    StatusCommitted
    StatusAborted
)

// Participant 参与者接口
type Participant interface {
    ID() string
    Prepare(ctx context.Context, txID string, operation interface{}) (bool, error)
    Commit(ctx context.Context, txID string) error
    Abort(ctx context.Context, txID string) error
}

// Coordinator 2PC 协调者
type Coordinator struct {
    id          string
    participants []Participant

    // 事务状态
    mu           sync.RWMutex
    transactions map[string]*TransactionRecord

    // 持久化
    log          TransactionLog
}

// TransactionRecord 事务记录
type TransactionRecord struct {
    ID           string
    Status       TransactionStatus
    Participants map[string]bool // participantID -> votedYes
    Operation    interface{}
}

// TransactionLog 事务日志接口 (用于故障恢复)
type TransactionLog interface {
    Write(record *TransactionRecord) error
    Read(txID string) (*TransactionRecord, error)
}

// NewCoordinator 创建协调者
func NewCoordinator(id string, log TransactionLog) *Coordinator {
    return &Coordinator{
        id:           id,
        transactions: make(map[string]*TransactionRecord),
        log:          log,
    }
}

// RegisterParticipant 注册参与者
func (c *Coordinator) RegisterParticipant(p Participant) {
    c.participants = append(c.participants, p)
}

// ExecuteTransaction 执行事务
func (c *Coordinator) ExecuteTransaction(ctx context.Context, txID string, operation interface{}) error {
    // 创建事务记录
    record := &TransactionRecord{
        ID:           txID,
        Status:       StatusPending,
        Participants: make(map[string]bool),
        Operation:    operation,
    }

    c.mu.Lock()
    c.transactions[txID] = record
    c.mu.Unlock()

    // 持久化
    if err := c.log.Write(record); err != nil {
        return fmt.Errorf("failed to log transaction: %w", err)
    }

    // Phase 1: Prepare
    votes, err := c.phase1Prepare(ctx, txID, operation)
    if err != nil {
        return c.abortTransaction(ctx, txID)
    }

    // 检查投票结果
    allYes := true
    for _, votedYes := range votes {
        if !votedYes {
            allYes = false
            break
        }
    }

    // Phase 2: Commit or Abort
    if allYes {
        return c.phase2Commit(ctx, txID)
    }
    return c.phase2Abort(ctx, txID)
}

// phase1Prepare Phase 1: 询问参与者
func (c *Coordinator) phase1Prepare(ctx context.Context, txID string, operation interface{}) (map[string]bool, error) {
    votes := make(map[string]bool)
    var mu sync.Mutex
    var wg sync.WaitGroup

    // 并行询问所有参与者
    for _, p := range c.participants {
        wg.Add(1)
        go func(participant Participant) {
            defer wg.Done()

            votedYes, err := participant.Prepare(ctx, txID, operation)

            mu.Lock()
            defer mu.Unlock()

            votes[participant.ID()] = votedYes && err == nil
        }(p)
    }

    wg.Wait()

    // 更新记录
    c.mu.Lock()
    if record, exists := c.transactions[txID]; exists {
        record.Participants = votes
        record.Status = StatusPrepared
        c.log.Write(record)
    }
    c.mu.Unlock()

    return votes, nil
}

// phase2Commit Phase 2a: 提交
func (c *Coordinator) phase2Commit(ctx context.Context, txID string) error {
    // 记录决定
    c.mu.Lock()
    record := c.transactions[txID]
    record.Status = StatusCommitted
    c.log.Write(record)
    c.mu.Unlock()

    // 通知所有参与者提交
    var wg sync.WaitGroup
    for _, p := range c.participants {
        wg.Add(1)
        go func(participant Participant) {
            defer wg.Done()
            participant.Commit(ctx, txID)
        }(p)
    }
    wg.Wait()

    return nil
}

// phase2Abort Phase 2b: 中止
func (c *Coordinator) phase2Abort(ctx context.Context, txID string) error {
    return c.abortTransaction(ctx, txID)
}

// abortTransaction 中止事务
func (c *Coordinator) abortTransaction(ctx context.Context, txID string) error {
    c.mu.Lock()
    record := c.transactions[txID]
    record.Status = StatusAborted
    c.log.Write(record)
    c.mu.Unlock()

    // 通知所有参与者中止
    var wg sync.WaitGroup
    for _, p := range c.participants {
        wg.Add(1)
        go func(participant Participant) {
            defer wg.Done()
            participant.Abort(ctx, txID)
        }(p)
    }
    wg.Wait()

    return errors.New("transaction aborted")
}

// Recover 故障恢复
func (c *Coordinator) Recover(ctx context.Context) error {
    // 读取所有未完成的事务
    // 根据日志状态决定继续提交或中止
    // ...
    return nil
}
```

### 7.2 Saga 实现

```go
package saga

import (
    "context"
    "fmt"
)

// Step Saga 步骤
type Step struct {
    Name        string
    Action      func(ctx context.Context) error
    Compensation func(ctx context.Context) error
}

// Saga 分布式事务
type Saga struct {
    Name  string
    Steps []Step
}

// Result 执行结果
type Result struct {
    Success       bool
    CompletedStep int
    Error         error
}

// Execute 执行 Saga
func (s *Saga) Execute(ctx context.Context) Result {
    executed := make([]int, 0, len(s.Steps))

    for i, step := range s.Steps {
        if err := step.Action(ctx); err != nil {
            // 执行补偿
            for j := len(executed) - 1; j >= 0; j-- {
                completedStep := s.Steps[executed[j]]
                if completedStep.Compensation != nil {
                    completedStep.Compensation(ctx)
                }
            }

            return Result{
                Success:       false,
                CompletedStep: i,
                Error:         err,
            }
        }

        executed = append(executed, i)
    }

    return Result{
        Success:       true,
        CompletedStep: len(s.Steps),
    }
}

// SagaOrchestrator Saga 编排器
type SagaOrchestrator struct {
    sagas map[string]*Saga
}

// NewSagaOrchestrator 创建编排器
func NewSagaOrchestrator() *SagaOrchestrator {
    return &SagaOrchestrator{
        sagas: make(map[string]*Saga),
    }
}

// Register 注册 Saga
func (o *SagaOrchestrator) Register(saga *Saga) {
    o.sagas[saga.Name] = saga
}

// Execute 执行指定 Saga
func (o *SagaOrchestrator) Execute(ctx context.Context, sagaName string) Result {
    saga, exists := o.sagas[sagaName]
    if !exists {
        return Result{
            Success: false,
            Error:   fmt.Errorf("saga %s not found", sagaName),
        }
    }

    return saga.Execute(ctx)
}
```

---

## 8. 学术参考文献

### 8.1 经典论文

1. **Gray, J. (1978)**. Notes on Data Base Operating Systems. *LNCS 60*.
   - 两阶段提交的原始描述

2. **Skeen, D. (1981)**. Nonblocking Commit Protocols. *ACM SIGMOD*.
   - 非阻塞提交协议

3. **Garcia-Molina, H., & Salem, K. (1987)**. Sagas. *ACM SIGMOD*.
   - Saga 模式，长事务的补偿方法

### 8.2 现代研究

1. **Gray, J., & Lamport, L. (2006)**. Consensus on Transaction Commit. *ACM Transactions on Database Systems*, 31(1), 133-160.
   - Paxos Commit：使用共识代替 2PC

2. **Thomson, A., et al. (2012)**. Calvin: Fast Distributed Transactions for Partitioned Database Systems. *SIGMOD*.
   - 确定性分布式事务

3. **Harding, R., et al. (2017)**. An Evaluation of Distributed Concurrency Control. *VLDB*.
   - 分布式并发控制的系统评估

---

## 9. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Distributed Transactions Toolkit                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心记忆锚点: "2PC = 共识问题"                                              │
│  ├── 2PC 的协调者等价于共识的 Leader                                         │
│  ├── 2PC 的投票等价于共识的 Promise                                          │
│  └── 2PC 的阻塞等价于共识的活性失败                                          │
│                                                                              │
│  协议选择决策:                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ 需要强一致 + 短事务 → 2PC                                             │    │
│  │ 需要高可用 + 可恢复   → 3PC (但需同步网络)                             │    │
│  │ 微服务 + 长事务      → Saga                                           │    │
│  │ 资源预留场景         → TCC                                            │    │
│  │ 最大可用性          → 异步消息 + 最终一致                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  常见误区:                                                                   │
│  ❌ "2PC 太慢不能用" → 优化后可达万级 TPS                                   │
│  ❌ "Saga 没有事务" → Saga 提供最终一致性的事务语义                          │
│  ❌ "3PC 完全解决阻塞" → 异步网络中仍然可能阻塞                              │
│  ❌ "分布式事务应该避免" → 很多时候是必需的                                  │
│                                                                              │
│  关键公式:                                                                   │
│  ├── 2PC 延迟 = 2 × RTT (最优)                                             │
│  ├── 3PC 延迟 = 3 × RTT (最优)                                             │
│  ├── Saga 延迟 = Σ(本地事务延迟)                                           │
│  └── 协调者故障恢复时间 = 检测时间 + 日志读取 + 通知参与者                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Learning Resources

### Academic Papers

1. **Gray, J., & Reuter, A.** (1993). *Transaction Processing: Concepts and Techniques*. Morgan Kaufmann. ISBN: 978-1558601901
2. **Mohan, C., et al.** (1992). ARIES: A Transaction Recovery Method Supporting Fine-Granularity Locking and Partial Rollbacks. *ACM TODS*, 17(1), 94-162. DOI: [10.1145/128765.128770](https://doi.org/10.1145/128765.128770)
3. **Garcia-Molina, H., & Salem, K.** (1987). Sagas. *ACM SIGMOD*, 249-259. DOI: [10.1145/38713.38742](https://doi.org/10.1145/38713.38742)
4. **Lampson, B., & Sturgis, H.** (1976). Crash Recovery in a Distributed Data Storage System. *Xerox PARC Technical Report*.

### Video Tutorials

1. **Martin Kleppmann.** (2018). [Distributed Transactions](https://www.youtube.com/watch?v=5ZjhNTM8XU8). QCon London.
2. **Pat Helland.** (2016). [ACID 2.0](https://www.youtube.com/watch?v=7HHPm1X4NGE). GOTO Conference.
3. **Cockroach Labs.** (2020). [Distributed Transactions at Scale](https://www.youtube.com/watch?v=tgTOvO6e35w). CockroachDB Tech Talk.
4. **ByteByteGo.** (2022). [Saga Pattern and Distributed Transactions](https://www.youtube.com/watch?v=0h0X8L15e1w). System Design.

### Book References

1. **Kleppmann, M.** (2017). *Designing Data-Intensive Applications* (Chapter 7: Transactions). O'Reilly Media.
2. **Bernstein, P. A., & Newcomer, E.** (2009). *Principles of Transaction Processing* (2nd ed.). Morgan Kaufmann.
3. **Weikum, G., & Vossen, G.** (2001). *Transactional Information Systems*. Morgan Kaufmann.
4. **Elmasri, R., & Navathe, S.** (2016). *Fundamentals of Database Systems* (Chapter 21). Pearson.

### Online Courses

1. **Coursera.** [Database Systems](https://www.coursera.org/learn/database-systems) - Transaction management.
2. **edX.** [Distributed Systems by TU Delft](https://www.edx.org/professional-certificate/delftx-cloud-computing) - Transactions module.
3. **Udemy.** [Database Design and Management](https://www.udemy.com/topic/database-design/) - ACID properties.
4. **CMU Database Group.** [Intro to Database Systems](https://www.youtube.com/playlist?list=PLSE8ODhjZXjbohkNB1QrY7J5ak8v6P9) - Lecture 16-18.

### GitHub Repositories

1. [dgraph-io/dgraph](https://github.com/dgraph-io/dgraph) - Distributed graph database with transactions.
2. [cockroachdb/cockroach](https://github.com/cockroachdb/cockroach) - Distributed SQL with serializable default.
3. [pingcap/tidb](https://github.com/pingcap/tidb) - MySQL-compatible distributed database.
4. [saga-go/saga](https://github.com/saga-go/saga) - Saga pattern implementation in Go.

### Conference Talks

1. **Pat Helland.** (2007). *Life Beyond Distributed Transactions*. CIDR.
2. **James Corbett.** (2012). *Spanner: Google's Globally-Distributed Database*. OSDI.
3. **Rebecca Taft.** (2020). *CockroachDB: The Resilient Geo-Distributed SQL Database*. SIGMOD.
4. **Peter Bailis.** (2014). *Coordination Avoidance in Database Systems*. VLDB.

---

*文档生成时间: 2026-04-02*
*维护者: Distributed Systems Knowledge Base*
*版本: S-Level (23+ KB)*
