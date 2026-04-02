# COMPARISON: Raft vs Paxos 共识算法对比

> **维度**: Formal Theory
> **级别**: S (16+ KB)
> **标签**: #consensus #raft #paxos #distributed-systems

---

## 核心对比

| 特性 | Raft | Paxos | Multi-Paxos |
|------|------|-------|-------------|
| **目标** | 可理解的共识 | 理论正确性 | 实用共识 |
| **分解** | 3子问题 (选举/日志/安全) | 单阶段 | 单阶段+领导者优化 |
| **领导者** | 强领导者 | 无 | 有 (优化后) |
| **选主** | 显式 Heartbeat | 隐式 | 显式 |
| **日志** | 连续复制 | 无序接受 | 连续复制 |
| **理解难度** | 简单 | 困难 | 中等 |
| **实现难度** | 中等 | 困难 | 中等 |
| **消息延迟** | 2 RTT (正常) | 2 RTT | 1 RTT (优化后) |
| **适用场景** | 教学、新系统 | 理论验证 | 生产系统 |

---

## 架构对比

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Raft Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                         Raft Node                                    │   │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────────────┐    │   │
│  │  │ Leader        │  │   Log         │  │      State Machine    │    │   │
│  │  │ Election      │  │   Replication │  │                       │    │   │
│  │  │               │  │               │  │  ┌─────────────────┐  │    │   │
│  │  │ ┌───────────┐ │  │ ┌───────────┐ │  │  │ Applied Entries │  │    │   │
│  │  │ │ Term      │ │  │ │ Index 1   │ │  │  │                 │  │    │   │
│  │  │ │ Vote      │ │  │ │ Index 2   │ │  │  │  User Data      │  │    │   │
│  │  │ │ Heartbeat │ │  │ │ Index 3   │ │  │  │                 │  │    │   │
│  │  │ └───────────┘ │  │ └───────────┘ │  │  └─────────────────┘  │    │   │
│  │  └───────────────┘  └───────────────┘  └───────────────────────┘    │   │
│  │                                                                      │   │
│  │  Safety: commitIndex ≥ matchIndex on majority                       │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Commit Rule: 条目在多数节点复制后才能提交                                      │
│  Leader Election: 超时+随机退避                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              Paxos Architecture                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Paxos Instance                                  │   │
│  │                                                                      │   │
│  │  Phase 1: Prepare                                                    │   │
│  │    ┌─────────┐              ┌─────────┐                             │   │
│  │    │ Proposer│─────────────►│Acceptors│  "Can I propose value?"     │   │
│  │    │         │◄─────────────│ (N/2+1) │  Promise(no lower proposal) │   │
│  │    └─────────┘              └─────────┘                             │   │
│  │                                                                      │   │
│  │  Phase 2: Accept                                                   │   │
│  │    ┌─────────┐              ┌─────────┐                             │   │
│  │    │ Proposer│─────────────►│Acceptors│  "Accept this value"        │   │
│  │    │         │◄─────────────│ (N/2+1) │  Accepted                   │   │
│  │    └─────────┘              └─────────┘                             │   │
│  │                                                                      │   │
│  │  Note: 每个值独立运行 Paxos，无连续日志概念                               │   │
│  │        Multi-Paxos 添加 Leader 优化连续值                                 │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 算法流程对比

### Raft 领导者选举

```go
// 简化示意
func (n *Node) startElection() {
    n.currentTerm++
    n.votedFor = n.id
    votes := 1

    // 向所有节点发送 RequestVote
    for _, peer := range n.peers {
        go func(p Peer) {
            reply := p.RequestVote(n.currentTerm, n.id, n.lastLogIndex, n.lastLogTerm)
            if reply.VoteGranted {
                votes++
                if votes > len(n.peers)/2 {
                    n.becomeLeader()
                }
            }
        }(peer)
    }
}

// 选举超时处理
timeout := random(150ms, 300ms)  // 随机退避
select {
case <-n.heartbeatCh:
    // 收到心跳，保持 follower
case <-time.After(timeout):
    n.startElection()  // 超时，发起选举
}
```

### Paxos Prepare Phase

```go
// 简化示意
func (p *Proposer) prepare(value Value) bool {
    proposalNum := p.generateProposalNumber()
    promises := 0

    for _, acceptor := range p.acceptors {
        promise := acceptor.Prepare(proposalNum)
        if promise.Ok {
            promises++
            // 学习已接受的值
            if promise.AcceptedValue != nil {
                value = promise.AcceptedValue
            }
        }
    }

    return promises > len(p.acceptors)/2
}

func (p *Proposer) accept(proposalNum int64, value Value) bool {
    accepts := 0
    for _, acceptor := range p.acceptors {
        if acceptor.Accept(proposalNum, value) {
            accepts++
        }
    }
    return accepts > len(p.acceptors)/2
}
```

---

## 生产应用

| 系统 | 算法 | 说明 |
|------|------|------|
| etcd | Raft | CoreOS 开发，Kubernetes 默认存储 |
| Consul | Raft | HashiCorp 服务发现 |
| TiKV | Raft | PingCAP 分布式 KV |
| Chubby | Paxos | Google 锁服务 |
| Spanner | Paxos | Google 全局数据库 |

---

## 选择建议

```
选择 Raft 如果你:
- 需要易于理解和维护的代码
- 团队对分布式系统经验有限
- 构建新的分布式系统
- 需要活跃的社区支持

选择 Paxos/Multi-Paxos 如果你:
- 需要极致的性能优化
- 已有 Paxos 经验
- 需要定制特殊场景
- 研究理论验证
```

---

## 参考文献

1. [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf) - Diego Ongaro, John Ousterhout
2. [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf) - Leslie Lamport
3. [Raft vs Paxos](https://web.stanford.edu/~ouster/cgi-bin/papers/raft-atc14) - 原始论文对比
