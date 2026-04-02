# FT-004: 分布式系统理论基础 (Distributed Systems Fundamentals: CAP, BASE, ACID)

> **维度**: Formal Theory
> **级别**: S (20+ KB)
> **标签**: #cap-theorem #base #acid #distributed-systems #consistency
> **权威来源**: [CAP Twelve Years Later](https://sites.cs.ucsb.edu/~rich/class/cs293b-cloud/papers/brewer-cap.pdf), [Harvest, Yield, and Scalable Tolerant Systems](https://cs.uwaterloo.ca/~brecht/courses/854-Emerging-2009/readings/harvest-yield.pdf), [Designing Data-Intensive Applications](https://dataintensive.net/)

---

## 一致性模型谱系

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Consistency Model Spectrum                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Strong ───────────────────────────────────────────► Eventual               │
│                                                                             │
│  Linearizable ──► Sequential ──► Causal ──► Session ──► Eventual            │
│       │              │            │           │            │                │
│       │              │            │           │            │                │
│    最强一致       顺序一致      因果一致    会话一致      最终一致             │
│    性能最差       性能较差      性能中等    性能较好      性能最好             │
│                                                                             │
│  适用场景:                                                                   │
│  • Linearizable: 金融交易、库存扣减                                           │
│  • Sequential:   聊天应用、协作编辑                                           │
│  • Causal:       社交网络、评论系统                                           │
│  • Eventual:     CDN、DNS、缓存                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ACID vs BASE

### ACID (传统数据库)

$$
\begin{aligned}
&\text{Atomicity (原子性): } \\
&\quad \forall t \in \text{Transactions}: t.commit() \lor t.rollback() \\
&\quad \text{没有部分提交的状态} \\
\\
&\text{Consistency (一致性): } \\
&\quad \forall t: \text{valid}(state_{pre}) \land t.execute() \Rightarrow \text{valid}(state_{post}) \\
&\quad \text{约束不被违反} \\
\\
&\text{Isolation (隔离性): } \\
&\quad \forall t_1, t_2: t_1 \parallel t_2 \Rightarrow \exists s \in \text{SerialSchedule}: result(t_1, t_2) = result(s) \\
&\quad \text{并发等价于串行} \\
\\
&\text{Durability (持久性): } \\
&\quad \forall t: t.commit() \Rightarrow \neg lost(t) \\
&\quad \text{已提交不丢失}
\end{aligned}
$$

### BASE (分布式系统)

| 特性 | 含义 | 与 ACID 对比 |
|------|------|-------------|
| **B**asically **A**vailable | 基本可用 | 牺牲部分一致性换取可用性 |
| **S**oft state | 软状态 | 允许中间状态，不要求实时一致 |
| **E**ventually consistent | 最终一致 | 不保证立即一致，但保证最终一致 |

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        ACID vs BASE Trade-offs                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ACID Systems                              BASE Systems                      │
│  ────────────                              ───────────                       │
│                                                                              │
│  • Strong consistency                      • High availability               │
│  • Isolated transactions                   • Partition tolerance             │
│  • Complex coordination                    • Simple replication              │
│  • Lower throughput                        • Higher throughput               │
│                                                                              │
│  Examples:                                 Examples:                         │
│  • PostgreSQL, MySQL (InnoDB)              • Cassandra, DynamoDB, MongoDB    │
│  • Single-node or small cluster            • Large-scale distributed         │
│                                                                              │
│  When to choose:                                                            │
│  • Financial transactions                  • Social media feeds              │
│  • Inventory management                    • User preferences                │
│  • Order processing                        • Analytics, logging              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## CAP 定理

### 形式化定义

$$
\text{For a distributed data store, it is impossible to simultaneously provide:} \\
\begin{aligned}
&\text{C - Consistency: } \forall n_i, n_j \in \text{Nodes}: read(n_i, k) = read(n_j, k) \text{ after write}(k, v) \\
&\text{A - Availability: } \forall n \in \text{Nodes}: \neg failed(n) \Rightarrow responds(n) \\
&\text{P - Partition Tolerance: } \forall p \in \text{Partitions}: system\ continues \\
\\
&\text{Theorem: } \neg(C \land A \land P) \\
&\text{Corollary: In presence of partitions, choose CP or AP}
\end{aligned}
$$

### CAP 权衡决策树

```
                        Network Partition?
                              │
                    ┌─────────┴─────────┐
                    │                   │
                   Yes                  No
                    │                   │
            ┌───────┴───────┐           │
            │               │           │
        Consistency?    Availability?   │
            │               │           │
      ┌─────┴─────┐   ┌─────┴─────┐     │
      │           │   │           │     │
     Yes          No  Yes         No    │
      │           │   │           │     │
      ▼           ▼   ▼           ▼     ▼
    ┌──────┐   ┌──────┐    ┌──────┐   ┌──────┐
    │  CP  │   │  ??? │    │  AP  │   │  CA  │
    │      │   │      │    │      │   │      │
    │ HBase│   │Not   │    │Cassa-│   │Single│
    │ Redis│   │Valid │    │ndra  │   │Node  │
    │(cluster)│  │      │    │Dynamo│   │      │
    └──────┘   └──────┘    └──────┘   └──────┘
```

### CP 系统示例

```go
// CP 系统：优先一致性，分区时拒绝服务
package cp

import (
    "context"
    "errors"
    "sync"
    "time"
)

// ConsistentStore CP 存储实现
type ConsistentStore struct {
    mu    sync.RWMutex
    data  map[string]string

    // Raft 共识层
    raft  *RaftNode

    // 配置
    majority int
}

func (s *ConsistentStore) Write(ctx context.Context, key, value string) error {
    // 需要多数派确认
    proposal := &WriteProposal{Key: key, Value: value}

    // 提交到 Raft
    future := s.raft.Apply(proposal, 5*time.Second)
    if err := future.Error(); err != nil {
        // 无法达到多数派，拒绝写入
        return errors.New("write rejected: unable to reach majority")
    }

    return nil
}

func (s *ConsistentStore) Read(ctx context.Context, key string) (string, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    // 检查是否为主节点
    if s.raft.State() != Leader {
        return "", errors.New("not leader: read from leader only")
    }

    value, ok := s.data[key]
    if !ok {
        return "", errors.New("key not found")
    }

    return value, nil
}

// 分区处理：如果失去多数派，停止服务
func (s *ConsistentStore) onPartitionDetected() {
    if s.raft.GetClusterSize()/2+1 > s.raft.GetActiveNodes() {
        s.enterDegradedMode()
    }
}

func (s *ConsistentStore) enterDegradedMode() {
    // CP 系统：拒绝所有写请求
    s.raft.SetReadOnly(true)
    log.Error("Entering degraded mode: partition detected, consistency prioritized")
}
```

### AP 系统示例

```go
// AP 系统：优先可用性，允许临时不一致
package ap

import (
    "context"
    "sync"
    "time"
)

// AvailableStore AP 存储实现
type AvailableStore struct {
    mu   sync.RWMutex
    data map[string]*DataVersion

    // 副本管理
    replicas []string

    // 冲突解决
    resolver ConflictResolver
}

type DataVersion struct {
    Value     string
    Timestamp int64
    VectorClock VectorClock
}

func (s *AvailableStore) Write(ctx context.Context, key, value string) error {
    version := &DataVersion{
        Value:     value,
        Timestamp: time.Now().UnixNano(),
        VectorClock: s.incrementClock(),
    }

    // 异步写入所有副本
    for _, replica := range s.replicas {
        go s.writeToReplica(replica, key, version)
    }

    // 立即返回（不等待副本确认）
    s.mu.Lock()
    s.data[key] = version
    s.mu.Unlock()

    return nil
}

func (s *AvailableStore) Read(ctx context.Context, key string) (string, error) {
    s.mu.RLock()
    localVersion := s.data[key]
    s.mu.RUnlock()

    // 从多个副本读取
    versions := s.readFromReplicas(key, 3) // R=3
    versions = append(versions, localVersion)

    // 冲突解决
    resolved := s.resolver.Resolve(versions)

    return resolved.Value, nil
}

// 分区处理：继续服务，标记为潜在不一致
func (s *AvailableStore) onPartitionDetected() {
    s.setHintedHandoffEnabled(true)
    log.Warn("Partition detected: continuing in available mode, conflicts may occur")
}
```

---

## PACELC 定理

CAP 的扩展：如果有分区 (P)，在可用性 (A) 和一致性 (C) 之间选择；
否则 (E)，在延迟 (L) 和一致性 (C) 之间选择。

$$
\text{PACELC: } \underbrace{PA/EL}_{\text{延迟优先}} \text{ vs } \underbrace{PC/EC}_{\text{一致优先}}
$$

| 系统 | 类型 | 说明 |
|------|------|------|
| DynamoDB | PA/EL | 默认最终一致，可配置强一致 |
| MongoDB | PA/EC | 主从复制，可配置写入确认级别 |
| Cassandra | PA/EL | 可调一致性级别 (ONE, QUORUM, ALL) |
| Bigtable | PC/EC | 强一致，但延迟较高 |

---

## 实践决策框架

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Choosing the Right Consistency Model                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Q1: Can you tolerate stale reads?                                          │
│      ├─ No → Use strong consistency (CP)                                    │
│      └─ Yes → Continue to Q2                                                │
│                                                                              │
│  Q2: What is the conflict resolution strategy?                              │
│      ├─ Last-write-wins → Use AP with timestamps                            │
│      ├─ Merge functions → Use CRDTs                                         │
│      ├─ User intervention → Use MVCC with versions                          │
│      └─ Avoid conflicts → Use domain partitioning                           │
│                                                                              │
│  Q3: What is the latency requirement?                                       │
│      ├─ < 10ms local → Use in-memory with async replication                 │
│      ├─ < 100ms regional → Use quorum replication                           │
│      └─ < 1s global → Use eventual consistency with conflict resolution     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 参考文献

1. [CAP Twelve Years Later: How the "Rules" Have Changed](https://sites.cs.ucsb.edu/~rich/class/cs293b-cloud/papers/brewer-cap.pdf) - Eric Brewer, 2012
2. [Harvest, Yield, and Scalable Tolerant Systems](https://cs.uwaterloo.ca/~brecht/courses/854-Emerging-2009/readings/harvest-yield.pdf) - Fox & Brewer
3. [Designing Data-Intensive Applications](https://dataintensive.net/) - Martin Kleppmann
4. [Consistency Models in Non-Relational Databases](https://www.vldb.org/pvldb/vol12/p1323-kleppmann.pdf) - Kleppmann et al.
5. [Eventually Consistent - Revisited](https://www.allthingsdistributed.com/2008/12/eventually_consistent.html) - Werner Vogels
