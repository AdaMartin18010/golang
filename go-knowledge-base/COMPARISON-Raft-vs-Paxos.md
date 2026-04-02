# 对比分析：Raft vs Paxos (Comparison: Raft vs Paxos)

> **分类**: 形式理论对比
> **相关文档**: FT-003, FT-007

---

## 核心对比

| 特性 | Raft | Paxos |
|------|------|-------|
| **设计目标** | 可理解性 | 正确性证明 |
| **领导者** | 强 Leader | 多 Proposer |
| **日志复制** | 原生支持 | Multi-Paxos 扩展 |
| **理解难度** | 相对容易 | 困难 |
| **工业实现** | etcd, Consul, TiKV | Chubby, ZooKeeper |
| **成员变更** | Joint Consensus | 复杂 |
| **性能** | 高吞吐 | 类似 |

---

## 算法流程对比

### Raft (3 阶段)

```
1. 选举 (Election)
   Follower ──timeout──► Candidate ──vote──► Leader

2. 日志复制 (Log Replication)
   Leader ──AppendEntries──► Followers
   Leader waits for majority ack

3. 安全性 (Safety)
   - Leader Completeness
   - State Machine Safety
```

### Paxos (2 阶段基础)

```
1. Prepare Phase
   Proposer ──Prepare(n)──► Acceptors
            ◄──Promise────┘

2. Accept Phase
   Proposer ──Accept(n, v)──► Acceptors
            ◄──Accepted─────┘

Multi-Paxos (优化):
   - 选 Leader 后跳过 Prepare
   - 直接 Accept
```

---

## 选择建议

| 场景 | 推荐 |
|------|------|
| 新系统开发 | Raft |
| 需要形式化证明 | Paxos |
| 快速实现 | Raft |
| 遗留系统兼容 | Paxos |
| 教学目的 | Raft |
| 理论研究 | Paxos |

---

## 代码复杂度对比

| 组件 | Raft LOC | Paxos LOC |
|------|----------|-----------|
| 选举 | ~200 | ~150 |
| 日志复制 | ~300 | ~400 |
| 成员变更 | ~150 | ~500 |
| 快照 | ~200 | ~300 |
| **总计** | **~850** | **~1350** |

---

## 共识算法家族

```
Consensus Algorithms
├── Paxos Family
│   ├── Basic Paxos
│   ├── Multi-Paxos
│   ├── Fast Paxos
│   └── Flexible Paxos
├── Raft Family
│   ├── Standard Raft
│   ├── Pre-vote Raft
│   └── Checksum Raft
└── Others
    ├── Zab (ZooKeeper)
    ├── Viewstamped Replication
    └── PBFT (Byzantine)
```

---

## 参考文献

1. [Raft Paper](https://raft.github.io/raft.pdf)
2. [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf)
3. [Paxos vs Raft](https://www.cockroachlabs.com/blog/paxos-vs-raft/)
