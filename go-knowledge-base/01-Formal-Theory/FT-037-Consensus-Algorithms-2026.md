# FT-037-Consensus-Algorithms-2026

> **Dimension**: 01-Formal-Theory  
> **Status**: S-Level  
> **Created**: 2026-04-03  
> **Version**: 2026 (Raft, Paxos, PBFT, HotStuff, Snow)  
> **Size**: >20KB 

---

## 1. 共识问题定义

### 1.1 形式化定义

**FLP不可能性结果**: 在异步系统中，即使只有一个故障进程，也不存在确定性共识算法。

**共识属性**:
```
设系统有n个进程，最多f个故障

终止性 (Termination):
  ∀p ∈ 正确进程: ◇(p决定某个值)

一致性 (Agreement):
  ∀p,q ∈ 正确进程: p的决定 = q的决定

有效性 (Validity):
  如果所有正确进程提议v，则决定v
```

### 1.2 CAP定理

```
CAP Theorem:
┌─────────────────────────────────────────┐
│  在分布式系统中，最多同时满足两项:      │
│                                         │
│  C - Consistency (一致性)               │
│  A - Availability (可用性)              │
│  P - Partition Tolerance (分区容错性)   │
│                                         │
│  必须选择CP或AP                         │
└─────────────────────────────────────────┘
```

---

## 2. Raft共识算法

### 2.1 概述

Raft是一种易于理解的共识算法，用于管理复制日志。

**状态机**:
```
Follower ──► Candidate ──► Leader
   ▲          │              │
   └──────────┴──────────────┘
   (发现更高任期或新Leader)
```

### 2.2 角色定义

| 角色 | 职责 |
|------|------|
| Leader | 处理所有客户端请求，复制日志 |
| Follower | 被动响应，复制Leader日志 |
| Candidate | 选举期间的临时状态 |

### 2.3 选举机制

```go
type RaftNode struct {
    id          int
    currentTerm int
    votedFor    int
    log         []LogEntry
    commitIndex int
    lastApplied int
    
    // Leader状态
    nextIndex   map[int]int
    matchIndex  map[int]int
    
    state       NodeState
    electionTimer *time.Timer
}

func (n *RaftNode) startElection() {
    n.state = Candidate
    n.currentTerm++
    n.votedFor = n.id
    votes := 1
    
    // 向所有节点请求投票
    for peer := range n.peers {
        go func(peer int) {
            args := RequestVoteArgs{
                Term:         n.currentTerm,
                CandidateId:  n.id,
                LastLogIndex: len(n.log) - 1,
                LastLogTerm:  n.log[len(n.log)-1].Term,
            }
            
            var reply RequestVoteReply
            if n.sendRequestVote(peer, &args, &reply) {
                n.handleVoteReply(&reply, &votes)
            }
        }(peer)
    }
}

func (n *RaftNode) handleVoteReply(reply *RequestVoteReply, votes *int) {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    if reply.Term > n.currentTerm {
        n.becomeFollower(reply.Term)
        return
    }
    
    if reply.VoteGranted && n.state == Candidate {
        *votes++
        if *votes > len(n.peers)/2 {
            n.becomeLeader()
        }
    }
}
```

### 2.4 日志复制

```go
func (n *RaftNode) appendEntries(entries []LogEntry) bool {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    // 找到日志匹配点
    for i, entry := range entries {
        if n.log[n.nextIndex[n.id]+i].Term != entry.Term {
            // 删除冲突条目及之后所有条目
            n.log = n.log[:n.nextIndex[n.id]+i]
            break
        }
    }
    
    // 追加新条目
    n.log = append(n.log, entries...)
    
    return true
}

func (n *RaftNode) replicateLog() {
    for peer := range n.peers {
        go func(peer int) {
            entries := n.log[n.nextIndex[peer]:]
            args := AppendEntriesArgs{
                Term:         n.currentTerm,
                LeaderId:     n.id,
                PrevLogIndex: n.nextIndex[peer] - 1,
                PrevLogTerm:  n.log[n.nextIndex[peer]-1].Term,
                Entries:      entries,
                LeaderCommit: n.commitIndex,
            }
            
            var reply AppendEntriesReply
            if n.sendAppendEntries(peer, &args, &reply) {
                n.handleAppendReply(peer, &reply)
            }
        }(peer)
    }
}
```

### 2.5 安全性证明

**选举安全**: 每个任期最多一个Leader
```
证明:
1. 节点只投票给任期 >= 当前任期的候选人
2. 节点在一个任期内只投一票
3. 需要多数派(n/2+1)才能成为Leader
4. 两个不同节点不能同时获得多数派投票
```

**日志匹配**: 如果两个日志条目有相同索引和任期，则它们存储相同命令
```
证明:
1. Leader只在其当前任期追加条目
2. 条目通过AppendEntries复制到Follower
3. PrevLogIndex/PrevLogTerm确保日志连续性
4. 冲突条目被删除后重新复制
```

---

## 3. Multi-Paxos

### 3.1 基本Paxos

**角色**:
- Proposer: 提出值
- Acceptor: 接受或拒绝提议
- Learner: 学习已决定的值

**两阶段协议**:
```
Phase 1 (Prepare):
  Proposer ──Prepare(n)──► Acceptor
  Acceptor ──Promise(n, v)◄── (如果n > 已承诺的最高编号)

Phase 2 (Accept):
  Proposer ──Accept(n, v)──► Acceptor  
  Acceptor ──Accepted(n, v)◄── (如果n == 已承诺的编号)
```

### 3.2 Multi-Paxos优化

```
Leader选举后，跳过Phase 1:

第1个值: Prepare → Promise → Accept → Accepted
第2+个值: Accept → Accepted (直接执行Phase 2)

效率提升: 2RTT → 1RTT
```

### 3.3 Go实现

```go
type PaxosNode struct {
    id        int
    seqNum    int
    highestPromised int
    acceptedValue   interface{}
    acceptedSeq     int
}

func (n *PaxosNode) Prepare(seqNum int) (interface{}, bool) {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    if seqNum <= n.highestPromised {
        return nil, false
    }
    
    n.highestPromised = seqNum
    return n.acceptedValue, true
}

func (n *PaxosNode) Accept(seqNum int, value interface{}) bool {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    if seqNum < n.highestPromised {
        return false
    }
    
    n.acceptedSeq = seqNum
    n.acceptedValue = value
    return true
}
```

---

## 4. PBFT (实用拜占庭容错)

### 4.1 概述

容忍拜占庭故障(恶意节点)的共识算法。

**容错能力**: n ≥ 3f + 1 (n个节点，f个拜占庭节点)

### 4.2 三阶段协议

```
Request → Pre-Prepare → Prepare → Commit → Reply

1. Client发送Request给Primary
2. Primary广播Pre-Prepare (分配序列号)
3. 所有节点广播Prepare (验证)
4. 收到2f个Prepare后广播Commit
5. 收到2f+1个Commit后执行
6. 发送Reply给Client
```

### 4.3 视图变更

```
当Primary故障时:
1. 节点启动view-change计时器
2. 超时后广播ViewChange消息
3. 新Primary收集ViewChange
4. 新Primary广播NewView
5. 继续正常操作
```

---

## 5. HotStuff

### 5.1 概述

Facebook Diem(Libra)使用的共识算法。

**特点**:
- 线性通信复杂度
- 视图变更简单高效
- 响应式(事件驱动而非超时驱动)

### 5.2 链式结构

```
┌────────┐    ┌────────┐    ┌────────┐
│ Block3 │───►│ Block2 │───►│ Block1 │───► Genesis
│ QC: 2f+1│    │ QC: 2f+1│    │ QC: 2f+1│
└────────┘    └────────┘    └────────┘

QC (Quorum Certificate): 2f+1个投票
每个区块包含对父区块的QC
```

### 5.3 三阶段投票

```
1. NewView: 广播新视图到Leader
2. Prepare: Leader提议区块，节点投票
3. PreCommit: 收到2f+1票后形成QC，广播
4. Commit: 收到PreCommit QC后再次投票
5. Decide: 收到Commit QC后区块确认
```

---

## 6. Avalanche/Snow协议

### 6.1 概述

新型概率性共识协议，用于Avalanche区块链。

**特点**:
- 亚秒级最终性
- 高吞吐 (>4500 TPS)
- 轻量级

### 6.2 Snowball算法

```
参数:
- k: 采样大小 (通常20)
- α: 多数派阈值 (通常14)
- β: 决策置信度 (通常20)

算法:
1. 节点有颜色(初始未着色)
2. 重复β轮:
   a. 随机采样k个节点查询
   b. 如果≥α个节点有相同颜色c，接受c
   c. 否则保持当前颜色
3. β轮后决定最终颜色
```

### 6.3 形式化分析

**安全性**: 拜占庭节点比例 < 50%
**活性**: 网络同步时保证进展
**最终性**: O(log n)轮后以高概率达成

---

## 7. 算法对比

### 7.1 特性对比

| 算法 | 容错类型 | 节点数 | 消息复杂度 | 延迟 | 使用案例 |
|------|---------|--------|-----------|------|---------|
| Paxos | 崩溃故障 | 2f+1 | O(n²) | 2RTT | Chubby |
| Raft | 崩溃故障 | 2f+1 | O(n) | 1RTT | etcd, TiKV |
| PBFT | 拜占庭 | 3f+1 | O(n²) | 3RTT | Hyperledger |
| HotStuff | 拜占庭 | 3f+1 | O(n) | 3RTT | Diem |
| Snow | 拜占庭 | 任意 | O(kn) | 亚秒 | Avalanche |

### 7.2 选择指南

```
崩溃故障场景:
  简单实现 → Raft
  理论正确性 → Paxos/Multi-Paxos

拜占庭故障场景:
  经典方案 → PBFT
  高性能 → HotStuff
  大规模 → Snow/Avalanche
```

---

## 8. 形式化验证

### 8.1 TLA+规格

```tla
--------------------------- MODULE Raft ---------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Servers, MaxTerm

VARIABLES currentTerm, state, votedFor, log, commitIndex

typeInvariant ==
  /\ currentTerm \in [Servers -> 0..MaxTerm]
  /\ state \in [Servers -> {"Follower", "Candidate", "Leader"}]
  /\ votedFor \in [Servers -> Servers \cup {Nil}]
  /\ log \in [Servers -> Seq([term: Nat, value: Nat])]
  /\ commitIndex \in [Servers -> Nat]

Init ==
  /\ currentTerm = [s \in Servers |-> 0]
  /\ state = [s \in Servers |-> "Follower"]
  /\ votedFor = [s \in Servers |-> Nil]
  /\ log = [s \in Servers |-> <<>>]
  /\ commitIndex = [s \in Servers |-> 0]

BecomeCandidate(s) ==
  /\ state[s] \in {"Follower", "Candidate"}
  /\ state' = [state EXCEPT ![s] = "Candidate"]
  /\ currentTerm' = [currentTerm EXCEPT ![s] = @ + 1]
  /\ votedFor' = [votedFor EXCEPT ![s] = s]
  /\ UNCHANGED <<log, commitIndex>>

BecomeLeader(s) ==
  /\ state[s] = "Candidate"
  /\ Cardinality({t \in Servers : votedFor[t] = s}) > Cardinality(Servers) \div 2
  /\ state' = [state EXCEPT ![s] = "Leader"]
  /\ UNCHANGED <<currentTerm, votedFor, log, commitIndex>>

Safety ==
  \* 已提交的日志条目不会被覆盖
  \* 所有节点在同一索引提交相同值
  \A s, t \in Servers :
    \A i \in 1..commitIndex[s] :
      i <= commitIndex[t] => log[s][i] = log[t][i]
============================================================================
```

---

## 9. 参考文献

1. "In Search of an Understandable Consensus Algorithm" (Raft)
2. "The Part-Time Parliament" (Paxos)
3. "Practical Byzantine Fault Tolerance" (PBFT)
4. "HotStuff: BFT Consensus in the Lens of Blockchain"
5. "Snowflake to Avalanche"
6. "FLP Impossibility Result"

---

*Last Updated: 2026-04-03*
