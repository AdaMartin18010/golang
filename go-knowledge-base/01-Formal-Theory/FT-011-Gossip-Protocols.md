# FT-011: Gossip 协议与流行病算法 (Gossip Protocols & Epidemic Algorithms)

> **维度**: Formal Theory
> **级别**: S (16+ KB)
> **标签**: #gossip #epidemic-algorithms #membership #distributed-systems
> **权威来源**: [Efficient Reconciliation and Flow Control for Anti-Entropy Protocols](https://www.cs.cornell.edu/home/rvr/papers/flowgossip.pdf)

---

## 核心概念

Gossip 协议灵感来自流行病传播，通过随机通信实现信息快速扩散。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Gossip Propagation Models                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. SI Model (Susceptible-Infected)                                         │
│     - 节点一旦被感染，永远保持感染状态                                        │
│     - 适用于：配置传播、状态同步                                              │
│                                                                              │
│  2. SIR Model (Susceptible-Infected-Removed)                                │
│     - 感染后可以被移除（免疫）                                                │
│     - 适用于：谣言传播、一次性通知                                            │
│                                                                              │
│  3. SIS Model (Susceptible-Infected-Susceptible)                            │
│     - 感染后可以恢复为易感状态                                                │
│     - 适用于：周期性同步、心跳检测                                            │
│                                                                              │
│  传播过程:                                                                   │
│                                                                              │
│  Round 1:              Round 2:              Round 3:                        │
│  ┌─────┐               ┌─────┐               ┌─────┐                        │
│  │  A  │───► B          │  A  │───► C         │  A  │───► D                 │
│  │     │               │     │───► D         │     │                        │
│  └─────┘               └─────┘               └─────┘                        │
│  (infected)            B ──► E               B ──► F                        │
│                        C ──► F               C ──► G                        │
│                                                                              │
│  理论传播速度: O(log N) 轮达到所有节点                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Gossip 协议类型

### 1. 反熵协议 (Anti-Entropy)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Anti-Entropy Gossip                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  目标: 保证最终一致性，修复数据差异                                            │
│                                                                              │
│  流程:                                                                       │
│  ┌─────────┐              ┌─────────┐                                       │
│  │ Node A  │───Digest───►│ Node B  │                                       │
│  │         │              │         │                                       │
│  │ {K1: v3,│              │ {K1: v3,│                                       │
│  │  K2: v1,│              │  K2: v2,│                                       │
│  │  K3: v5}│              │  K3: v5}│                                       │
│  └─────────┘              └────┬────┘                                       │
│                                │                                            │
│                                │ Compare                                     │
│                                ▼                                            │
│                         K2: A(v1) < B(v2)                                   │
│                         K3: A(v5) = B(v5)                                   │
│                                │                                            │
│  ┌─────────┐              ┌────┴────┐                                       │
│  │ Node A  │◄──K2: v2────│ Node B  │                                       │
│  │         │              │         │                                       │
│  │ {K1: v3,│              │         │                                       │
│  │  K2: v2,│              │         │                                       │
│  │  K3: v5}│              │         │                                       │
│  └─────────┘              └─────────┘                                       │
│                                                                              │
│  Digest 结构: Merkle Tree / Bloom Filter / Checksum 列表                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2. 谣言传播 (Rumor Mongering)

```go
package gossip

import (
    "math/rand"
    "sync"
    "time"
)

// Rumor 谣言消息
type Rumor struct {
    ID        string
    Payload   []byte
    Timestamp time.Time
    Hops      int
}

// Node Gossip 节点
type Node struct {
    ID       string
    peers    []string
    rumors   map[string]*Rumor
    hot      map[string]bool // 热谣言（还在传播）
    mu       sync.RWMutex
    fanout   int   // 每轮传播的节点数
    rounds   int   // 传播轮数阈值
}

// NewNode 创建节点
func NewNode(id string, fanout, rounds int) *Node {
    return &Node{
        ID:     id,
        peers:  make([]string, 0),
        rumors: make(map[string]*Rumor),
        hot:    make(map[string]bool),
        fanout: fanout,
        rounds: rounds,
    }
}

// AddPeer 添加对等节点
func (n *Node) AddPeer(peerID string) {
    n.peers = append(n.peers, peerID)
}

// SpreadRumor 传播新谣言
func (n *Node) SpreadRumor(rumor *Rumor) {
    n.mu.Lock()
    defer n.mu.Unlock()

    rumor.Hops = 0
    n.rumors[rumor.ID] = rumor
    n.hot[rumor.ID] = true

    // 开始传播
    go n.gossip(rumor.ID, 0)
}

// gossip 执行 gossip 传播
func (n *Node) gossip(rumorID string, round int) {
    if round >= n.rounds {
        // 达到轮数阈值，停止传播
        n.mu.Lock()
        delete(n.hot, rumorID)
        n.mu.Unlock()
        return
    }

    // 随机选择 fanout 个节点
    targets := n.selectRandomPeers(n.fanout)

    for _, target := range targets {
        go func(peerID string) {
            // 发送谣言
            n.sendRumor(peerID, rumorID)
        }(target)
    }

    // 延迟后下一轮
    time.Sleep(time.Second)
    n.gossip(rumorID, round+1)
}

// OnRumorReceived 收到谣言
func (n *Node) OnRumorReceived(rumor *Rumor) bool {
    n.mu.Lock()
    defer n.mu.Unlock()

    // 已知道？
    if _, exists := n.rumors[rumor.ID]; exists {
        return false // 已知道，不再传播
    }

    // 新谣言
    rumor.Hops++
    n.rumors[rumor.ID] = rumor
    n.hot[rumor.ID] = true

    // 继续传播
    go n.gossip(rumor.ID, 0)
    return true
}

// selectRandomPeers 随机选择节点
func (n *Node) selectRandomPeers(count int) []string {
    if len(n.peers) <= count {
        return n.peers
    }

    // Fisher-Yates 洗牌
    perm := rand.Perm(len(n.peers))
    result := make([]string, count)
    for i := 0; i < count; i++ {
        result[i] = n.peers[perm[i]]
    }
    return result
}

func (n *Node) sendRumor(peerID, rumorID string) {
    // 实际网络发送逻辑
}
```

---

## 应用场景

| 系统 | 协议 | 用途 |
|------|------|------|
| Cassandra | Gossip | 集群成员发现、故障检测 |
| Consul | Serf (Gossip) | 成员列表、事件广播 |
| Redis Cluster | Gossip | 节点发现、配置传播 |
| Bitcoin | Gossip | 交易传播、区块同步 |
| Dynamo | Gossip | 成员检测、 Merkle 树同步 |

---

## 数学分析

```
传播轮数分析:

设:
- N = 节点总数
- k = fanout (每轮传播的节点数)
- p = 节点存活概率

期望值:
- 达到所有节点所需轮数: O(log N / log k)
- 消息复杂度: O(N log N) (每个节点接收多次)
- 带宽: 每个节点 O(log N) 条消息

概率保证:
- 消息丢失容错: 即使 50% 消息丢失，仍可传播到 99% 节点
- 节点故障容错: 可容忍 N/2 节点故障
```

---

## 参考文献

1. [Efficient Reconciliation and Flow Control for Anti-Entropy Protocols](https://www.cs.cornell.edu/home/rvr/papers/flowgossip.pdf)
2. [Gossip Protocols](https://www.cs.cornell.edu/home/rvr/papers/gossip.pdf)
3. [SWIM: Scalable Weakly-consistent Infection-style Process Group Membership Protocol](https://www.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf)
