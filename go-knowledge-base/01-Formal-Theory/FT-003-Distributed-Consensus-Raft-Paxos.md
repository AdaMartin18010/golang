# FT-003: 分布式共识：Raft 与 Paxos (Distributed Consensus: Raft & Paxos)

> **维度**: Formal Theory
> **级别**: S (30+ KB)
> **标签**: #consensus #raft #paxos #distributed-systems
> **权威来源**: [Raft Paper](https://raft.github.io/raft.pdf), [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf), [The Part-Time Parliament](https://lamport.azurewebsites.net/pubs/lamport-paxos.pdf)

---

## 共识问题定义

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Distributed Consensus Problem                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  问题：如何让分布式系统中的多个节点就某个值达成一致？                          │
│                                                                              │
│  要求（Safety & Liveness）：                                                  │
│  ────────────────────────────                                                 │
│                                                                              │
│  Safety（安全性）：                                                           │
│  ├─ 所有正确节点最终同意同一个值                                              │
│  ├─ 一旦决定，不可更改                                                        │
│  └─ 只有被提议的值才能被决定                                                  │
│                                                                              │
│  Liveness（活性）：                                                           │
│  ├─ 最终一定能做出决定（部分节点故障时）                                       │
│  └─ 非拜占庭故障（节点不会撒谎）                                              │
│                                                                              │
│  FLP 不可能结果：                                                             │
│  在异步网络中，即使只有一个节点故障，也不存在确定性共识算法。                    │
│  解决：使用超时（partial synchrony）或随机化。                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Paxos 算法

### 核心角色

| 角色 | 职责 |
|------|------|
| Proposer | 提议者，提出值 |
| Acceptor | 接受者，决定是否接受提议 |
| Learner | 学习者，学习最终决定的值 |

### 两阶段协议

```
Phase 1: Prepare
────────────────
Proposer ──► Acceptor: Prepare(n)
            n = 提案编号（全局唯一，单调递增）

Acceptor ──► Proposer: Promise(n, (v, n_prev))
            承诺：不再接受编号小于 n 的提案
            返回：已接受的最高编号的提案

Phase 2: Accept
────────────────
Proposer ──► Acceptor: Accept(n, v)
            v = 要提议的值（如果 Phase 1 已有接受值，则用该值）

Acceptor ──► Proposer: Accepted(n, v)
            接受该提案

Majority: 只要多数 Acceptor 接受，即达成共识
```

### Paxos 完整 Go 实现

```go
package paxos

import (
    "sync"
    "sync/atomic"
)

// Message types
type MessageType int

const (
    Prepare MsgType = iota
    Promise
    Accept
    Accepted
)

type Message struct {
    Type      MessageType
    From      int
    ProposalN int64
    Value     interface{}
    PrevN     int64  // Promise 中返回的先前提案号
    PrevValue interface{}
}

// Acceptor 实现 Paxos Acceptor
type Acceptor struct {
    id int

    // 已承诺的最高提案号
    promisedN int64

    // 已接受的提案
    acceptedN int64
    acceptedV interface{}

    mu sync.RWMutex
}

func (a *Acceptor) HandlePrepare(msg Message) Message {
    a.mu.Lock()
    defer a.mu.Unlock()

    if msg.ProposalN > a.promisedN {
        a.promisedN = msg.ProposalN

        return Message{
            Type:      Promise,
            From:      a.id,
            ProposalN: msg.ProposalN,
            PrevN:     a.acceptedN,
            PrevValue: a.acceptedV,
        }
    }

    // 已承诺更高编号的提案，拒绝
    return Message{
        Type:      Promise,
        From:      a.id,
        ProposalN: a.promisedN, // 返回已承诺的编号
    }
}

func (a *Acceptor) HandleAccept(msg Message) Message {
    a.mu.Lock()
    defer a.mu.Unlock()

    if msg.ProposalN >= a.promisedN {
        a.acceptedN = msg.ProposalN
        a.acceptedV = msg.Value

        return Message{
            Type:      Accepted,
            From:      a.id,
            ProposalN: msg.ProposalN,
            Value:     msg.Value,
        }
    }

    // 已承诺更高编号的提案，拒绝
    return Message{
        Type:      Accepted,
        From:      a.id,
        ProposalN: -1, // 拒绝标记
    }
}
```

---

## Raft 算法

Raft 将共识问题分解为三个子问题：

1. **Leader Election**（领导者选举）
2. **Log Replication**（日志复制）
3. **Safety**（安全性）

### 状态机

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Raft Node States                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                    ┌─────────────┐                                          │
│         ┌─────────│  Follower   │◄────────┐                                │
│         │         └──────┬──────┘         │                                │
│  Timeout│ election       │                │Heartbeat                      │
│  expired│ timeout        │Discover leader │from leader                    │
│         │                │or new term     │                               │
│         ▼                ▼                │                                │
│    ┌─────────┐     ┌──────────┐           │                                │
│    │ Candidate│────►│  Leader  │───────────┘                                │
│    └─────────┘ win  └──────────┘                                            │
│     votes    election                                                       │
│                                                                              │
│  Election Timeout: 150-300ms (随机，避免活锁)                                │
│  Heartbeat Interval: ~50ms                                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 日志复制

```go
package raft

// LogEntry 日志条目
type LogEntry struct {
    Index   int
    Term    int
    Command interface{}
}

// Raft 状态
type Raft struct {
    // 持久化状态
    CurrentTerm int
    VotedFor    int
    Log         []LogEntry

    // 易失状态
    CommitIndex int
    LastApplied int

    // 领导者状态（仅 Leader 使用）
    NextIndex  []int  // 每个 follower 下一个要发送的日志索引
    MatchIndex []int  // 每个 follower 已复制的最高日志索引

    // 角色
    State NodeState

    // 配置
    Me          int
    Peers       []int
    ApplyCh     chan ApplyMsg
}

// AppendEntries RPC (Leader 调用)
type AppendEntriesArgs struct {
    Term         int
    LeaderId     int
    PrevLogIndex int        // 前一个日志索引
    PrevLogTerm  int        // 前一个日志任期
    Entries      []LogEntry // 要复制的日志（空表示心跳）
    LeaderCommit int        // Leader 的 CommitIndex
}

type AppendEntriesReply struct {
    Term    int
    Success bool
    // 优化：快速回退
    ConflictIndex int  // 冲突日志索引
    ConflictTerm  int  // 冲突日志任期
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
    rf.mu.Lock()
    defer rf.mu.Unlock()

    reply.Term = rf.CurrentTerm

    // 1. 如果 term < currentTerm，拒绝
    if args.Term < rf.CurrentTerm {
        reply.Success = false
        return
    }

    // 2. 如果 term > currentTerm，转为 follower
    if args.Term > rf.CurrentTerm {
        rf.CurrentTerm = args.Term
        rf.VotedFor = -1
        rf.State = Follower
    }

    // 重置选举超时
    rf.resetElectionTimer()

    // 3. 检查日志匹配
    if args.PrevLogIndex >= len(rf.Log) {
        reply.Success = false
        reply.ConflictIndex = len(rf.Log)
        reply.ConflictTerm = -1
        return
    }

    if rf.Log[args.PrevLogIndex].Term != args.PrevLogTerm {
        reply.Success = false
        reply.ConflictTerm = rf.Log[args.PrevLogIndex].Term
        // 找到该 term 的第一个索引
        for i := args.PrevLogIndex; i >= 0; i-- {
            if rf.Log[i].Term != reply.ConflictTerm {
                reply.ConflictIndex = i + 1
                break
            }
        }
        return
    }

    // 4. 复制日志
    for i, entry := range args.Entries {
        index := args.PrevLogIndex + 1 + i
        if index >= len(rf.Log) {
            rf.Log = append(rf.Log, entry)
        } else if rf.Log[index].Term != entry.Term {
            // 冲突，截断并追加
            rf.Log = rf.Log[:index]
            rf.Log = append(rf.Log, entry)
        }
    }

    // 5. 更新 commitIndex
    if args.LeaderCommit > rf.CommitIndex {
        rf.CommitIndex = min(args.LeaderCommit, len(rf.Log)-1)
        rf.apply()
    }

    reply.Success = true
}
```

---

## Raft vs Paxos 对比

| 特性 | Paxos | Raft |
|------|-------|------|
| 理解难度 | 困难 | 相对容易 |
| 实现复杂度 | 高 | 中等 |
| Leader 选举 | 隐式 | 显式 |
| 日志复制 | Multi-Paxos 扩展 | 原生支持 |
| 成员变更 | 复杂 | 联合共识（Joint Consensus） |
| 工业应用 | Chubby, ZooKeeper | etcd, Consul, TiKV |

---

## 形式化验证

```
Raft Safety Properties:
───────────────────────
1. Election Safety: 每个任期最多一个 leader
   ∀t: |{n | n.state = Leader ∧ n.term = t}| ≤ 1

2. Leader Append-Only: Leader 从不覆盖或删除日志条目
   ∀n: n.state = Leader ⇒ monotonic(n.log)

3. Log Matching: 如果两个日志条目有相同的 index 和 term，
   则它们存储相同的命令，且之前所有日志也相同

4. Leader Completeness: 如果一条日志已提交，
   则它出现在所有未来任期的 leader 的日志中

5. State Machine Safety: 如果节点应用了 index i 的日志到状态机，
   则没有其他节点会在 index i 应用不同的命令
```

---

## 参考文献

1. [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf) - Ongaro & Ousterhout
2. [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf) - Lamport
3. [The Part-Time Parliament](https://lamport.azurewebsites.net/pubs/lamport-paxos.pdf) - Lamport
4. [Raft Consensus Algorithm](https://raft.github.io/) - 官方可视化
5. [Consensus: Bridging Theory and Practice](https://web.stanford.edu/~ouster/cgi-bin/papers/OngaroPhD.pdf) - Ongaro PhD Thesis
