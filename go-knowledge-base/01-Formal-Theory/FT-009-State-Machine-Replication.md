# FT-009: 复制状态机理论 (State Machine Replication)

> **维度**: Formal Theory
> **级别**: S (15+ KB)
> **标签**: #state-machine-replication #smr #deterministic #replication
> **权威来源**: [Implementing Fault-Tolerant Services](https://www.cs.cornell.edu/home/rvr/papers/osdi04.pdf) - Schneider

---

## 核心概念

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      State Machine Replication (SMR)                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心思想：                                                                   │
│  如果所有副本从相同初始状态开始，                                              │
│  按相同顺序执行相同操作，                                                      │
│  那么最终状态必然一致。                                                        │
│                                                                              │
│  ┌─────────────┐     命令序列      ┌─────────────┐                          │
│  │   Client    │───► [C1,C2,C3] ───►│  Replica 1  │                          │
│  └─────────────┘                   │  (State S1) │                          │
│                                    └──────┬──────┘                          │
│                                           │ Apply(C1,C2,C3)                  │
│                                           ▼                                 │
│                                    ┌─────────────┐                          │
│  ┌─────────────┐     相同序列      │  Replica 2  │                          │
│  │   Client    │───► [C1,C2,C3] ───►│  (State S2) │                          │
│  └─────────────┘                   └──────┬──────┘                          │
│                                           │ Apply(C1,C2,C3)                  │
│                                           ▼                                 │
│                                    ┌─────────────┐                          │
│                                    │  Final State│                          │
│                                    │   S1 == S2  │                          │
│                                    └─────────────┘                          │
│                                                                              │
│  要求：                                                                       │
│  1. 确定性 (Deterministic): 相同输入产生相同输出                              │
│  2. 全序 (Total Order): 所有副本执行顺序一致                                  │
│  3. 容错 (Fault Tolerance): 容忍 f 个故障节点                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 形式化定义

$$
\begin{aligned}
&\text{State Machine: } M = (S, S_0, C, O, T) \\
&\text{where:} \\
&\quad S: \text{状态集合} \\
&\quad S_0 \in S: \text{初始状态} \\
&\quad C: \text{命令集合} \\
&\quad O: \text{输出集合} \\
&\quad T: S \times C \rightarrow S \times O \text{ (状态转换函数)} \\
\\
&\text{Replication:} \\
&\quad \forall i, j \in \text{Replicas}: \\
&\quad \quad \text{if } \pi_i = \pi_j \text{ (相同命令序列)} \\
&\quad \quad \text{then } state_i = state_j \\
\\
&\text{Safety:} \\
&\quad \text{所有正确副本达成一致} \\
\\
&\text{Liveness:} \\
&\quad \text{最终所有命令被执行}
\end{aligned}
$$

---

## 关键组件

### 1. 共识模块 (Consensus)

```go
// 使用 Raft/Paxos 保证命令顺序一致
type ConsensusModule interface {
    Propose(cmd Command) (Index, error)
    GetLog(index Index) (Command, error)
}
```

### 2. 状态机 (State Machine)

```go
type StateMachine interface {
    Apply(cmd Command) (Result, error)
    Snapshot() (Snapshot, error)
    Restore(snap Snapshot) error
}
```

### 3. 复制管理器

```go
type ReplicationManager struct {
    consensus ConsensusModule
    stateMachine StateMachine
    log []LogEntry
    lastApplied Index
}

func (rm *ReplicationManager) Execute(cmd Command) (Result, error) {
    // 1. 提交到共识模块
    index, err := rm.consensus.Propose(cmd)
    if err != nil {
        return nil, err
    }

    // 2. 等待提交
    for rm.lastApplied < index {
        entry := rm.consensus.GetLog(rm.lastApplied + 1)

        // 3. 应用到状态机
        result, err := rm.stateMachine.Apply(entry.Command)
        if err != nil {
            return nil, err
        }

        rm.lastApplied++

        if rm.lastApplied == index {
            return result, nil
        }
    }

    return nil, errors.New("unexpected state")
}
```

---

## 示例：Key-Value Store 复制

```go
type KVStateMachine struct {
    data map[string]string
}

func (kv *KVStateMachine) Apply(cmd Command) (Result, error) {
    switch cmd.Op {
    case "SET":
        kv.data[cmd.Key] = cmd.Value
        return "OK", nil
    case "GET":
        return kv.data[cmd.Key], nil
    case "DELETE":
        delete(kv.data, cmd.Key)
        return "OK", nil
    default:
        return nil, fmt.Errorf("unknown operation: %s", cmd.Op)
    }
}

// 确定性保证
// - 相同命令序列产生相同状态
// - 无随机性，无外部依赖
// - 时间戳使用逻辑时钟
```

---

## 与 Primary-Backup 对比

| 特性 | State Machine Replication | Primary-Backup |
|------|---------------------------|----------------|
| 复制粒度 | 操作/命令 | 状态/数据 |
| 带宽 | 低（只复制命令） | 高（复制状态） |
| CPU | 高（所有副本执行） | 低（只有主执行） |
| 适用 | CPU 便宜，带宽贵 | CPU 贵，带宽便宜 |
| 延迟 | 取决于共识 | 取决于网络 |

---

## 参考文献

1. [Implementing Fault-Tolerant Services Using the State Machine Approach](https://www.cs.cornell.edu/home/rvr/papers/osdi04.pdf) - Schneider
2. [Raft: In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)
3. [Paxos Made Live](https://research.google/pubs/paxos-made-live-an-engineering-perspective/)
