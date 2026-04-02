# FT-013: 拜占庭容错与 PBFT 算法 (Byzantine Fault Tolerance & PBFT)

> **维度**: Formal Theory
> **级别**: S (17+ KB)
> **标签**: #bft #pbft #byzantine-faults #distributed-consensus
> **权威来源**: [Practical Byzantine Fault Tolerance](http://pmg.csail.mit.edu/papers/osdi99.pdf) - Castro & Liskov

---

## 拜占庭将军问题

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Byzantine Generals Problem                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  场景: 拜占庭军队围攻城市，将军们通过信使通信                                    │
│                                                                              │
│  挑战:                                                                       │
│  1. 叛徒可能发送虚假消息                                                      │
│  2. 信使可能被截获/篡改                                                       │
│  3. 需要忠诚将军达成一致的决策                                                │
│                                                                              │
│  图示:                                                                       │
│         ┌─────────┐                                                          │
│         │ Commander │──► "Attack" ──► ┌─────────┐                          │
│         │ (忠诚)    │──► "Attack" ──► │ General │ 叛徒?                    │
│         └─────────┘  ──► "Retreat"──► │   B     │                          │
│                                       └─────────┘                          │
│                                                                              │
│  拜占庭容错条件:                                                              │
│  - 系统可容忍 f 个拜占庭节点                                                  │
│  - 需要至少 3f + 1 个节点                                                     │
│  - 即: 叛徒不超过总数的 1/3                                                   │
│                                                                              │
│  比较:                                                                        │
│  - 崩溃容错 (Crash Fault): 节点停止响应，共需要 2f+1 个节点 (多数派)          │
│  - 拜占庭容错 (Byzantine Fault): 节点可能撒谎，需要 3f+1 个节点               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## PBFT 算法

### 三阶段协议

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      PBFT Three-Phase Protocol                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Client ──Request──► Primary (Leader)                                       │
│                              │                                               │
│  Phase 1: PRE-PREPARE                                                        │
│                              │                                               │
│  Primary ──PRE-PREPARE(seq, digest)──► All Replicas                         │
│                              │                                               │
│  Phase 2: PREPARE                                                           │
│                              │                                               │
│  Each Replica ──PREPARE(seq, digest)──► All Replicas                        │
│                              │                                               │
│  Phase 3: COMMIT                                                            │
│                              │                                               │
│  Replica (收到 2f PREPAREs) ──COMMIT(seq)──► All Replicas                   │
│                              │                                               │
│  Execution (收到 2f+1 COMMITs)                                              │
│                              │                                               │
│  Replicas ──Reply──► Client                                                 │
│                                                                              │
│  视图更换 (View Change): 当主节点故障时触发                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### PBFT 消息类型

```go
package pbft

import (
    "crypto/sha256"
    "encoding/hex"
    "time"
)

// MessageType PBFT 消息类型
type MessageType int

const (
    REQUEST MessageType = iota
    PRE_PREPARE
    PREPARE
    COMMIT
    REPLY
    VIEW_CHANGE
    NEW_VIEW
    CHECKPOINT
)

// Message PBFT 消息
type Message struct {
    Type       MessageType
    View       int           // 当前视图号
    Sequence   int           // 序列号
    Digest     string        // 请求摘要
    Request    *Request      // 原始请求
    ReplicaID  int           // 发送者 ID
    Signature  []byte        // 数字签名
    Timestamp  time.Time
}

// Request 客户端请求
type Request struct {
    ClientID  string
    Timestamp int64
    Operation []byte
}

// Digest 计算请求摘要
func Digest(req *Request) string {
    h := sha256.New()
    h.Write([]byte(req.ClientID))
    h.Write([]byte(req.Operation))
    return hex.EncodeToString(h.Sum(nil))
}

// PBFTNode 节点
type PBFTNode struct {
    ID          int
    View        int
    Sequence    int
    IsPrimary   bool
   Replicas    []*Replica

    // 状态
    prePrepares map[int]*Message  // sequence -> pre-prepare
    prepares    map[int][]*Message // sequence -> prepares
    commits     map[int][]*Message // sequence -> commits

    // 检查点
    checkpoints map[int]Checkpoint

    // 请求日志
    log         map[int]*Request
}

// Replica 对等节点信息
type Replica struct {
    ID      int
    Address string
    PubKey  []byte
}

// Checkpoint 检查点
type Checkpoint struct {
    Sequence int
    Digest   string
    Proofs   []*Message
}
```

### 核心逻辑

```go
// HandleRequest 处理客户端请求 (主节点)
func (n *PBFTNode) HandleRequest(req *Request) {
    if !n.IsPrimary {
        return // 转发给主节点
    }

    n.Sequence++
    digest := Digest(req)

    // 发送 PRE-PREPARE
    prePrepare := &Message{
        Type:     PRE_PREPARE,
        View:     n.View,
        Sequence: n.Sequence,
        Digest:   digest,
        Request:  req,
        ReplicaID: n.ID,
    }

    n.broadcast(prePrepare)
}

// HandlePrePrepare 处理 PRE-PREPARE 消息
func (n *PBFTNode) HandlePrePrepare(msg *Message) {
    // 验证: 序列号、视图号、摘要
    if !n.validPrePrepare(msg) {
        return
    }

    n.prePrepares[msg.Sequence] = msg

    // 发送 PREPARE
    prepare := &Message{
        Type:     PREPARE,
        View:     msg.View,
        Sequence: msg.Sequence,
        Digest:   msg.Digest,
        ReplicaID: n.ID,
    }

    n.broadcast(prepare)
}

// HandlePrepare 处理 PREPARE 消息
func (n *PBFTNode) HandlePrepare(msg *Message) {
    n.prepares[msg.Sequence] = append(n.prepares[msg.Sequence], msg)

    // 收到 2f 个 PREPARE (包括自己的)，进入 prepared 状态
    if len(n.prepares[msg.Sequence]) >= 2*n.f() {
        // 发送 COMMIT
        commit := &Message{
            Type:     COMMIT,
            View:     msg.View,
            Sequence: msg.Sequence,
            Digest:   msg.Digest,
            ReplicaID: n.ID,
        }
        n.broadcast(commit)
    }
}

// HandleCommit 处理 COMMIT 消息
func (n *PBFTNode) HandleCommit(msg *Message) {
    n.commits[msg.Sequence] = append(n.commits[msg.Sequence], msg)

    // 收到 2f+1 个 COMMIT，执行请求
    if len(n.commits[msg.Sequence]) >= 2*n.f()+1 {
        n.execute(msg.Sequence)
    }
}

// execute 执行请求
func (n *PBFTNode) execute(sequence int) {
    req := n.log[sequence]
    // 执行业务逻辑
    result := n.apply(req.Operation)

    // 发送 REPLY 给客户端
    reply := &Message{
        Type:      REPLY,
        View:      n.View,
        Sequence:  sequence,
        ReplicaID: n.ID,
        Result:    result,
    }

    n.sendToClient(req.ClientID, reply)
}

// f 计算最大容错数
func (n *PBFTNode) f() int {
    return (len(n.Replicas) - 1) / 3
}
```

---

## 性能与优化

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      PBFT Performance & Optimizations                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  基础 PBFT:                                                                  │
│  - 通信复杂度: O(n²) 每请求                                                   │
│  - 延迟: 3 轮网络往返                                                         │
│  - 吞吐量: 受限                                                               │
│                                                                              │
│  优化方案:                                                                   │
│                                                                              │
│  1. 批量处理 (Batching)                                                      │
│     - 主节点收集多个请求，一次性处理                                          │
│     - 吞吐量提升 10x+                                                        │
│                                                                              │
│  2. 流水线 (Pipelining)                                                      │
│     - 不等待当前请求完成，继续下一个                                          │
│     - 提高吞吐量，保持延迟                                                    │
│                                                                              │
│  3. 投机执行 (Speculative Execution)                                         │
│     - 收到 PRE-PREPARE 后立即执行，如果失败则回滚                              │
│     - 降低延迟                                                               │
│                                                                              │
│  4. 检查点与垃圾回收                                                          │
│     - 定期生成检查点，清理旧日志                                              │
│     - 状态传输给落后节点                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 应用场景

| 系统 | BFT 算法 | 用途 |
|------|----------|------|
| Hyperledger Fabric | PBFT/SBFT | 联盟链共识 |
| Tendermint | BFT + PoS | Cosmos 区块链 |
| HotStuff | 改进 BFT | Facebook Libra |
| Algorand | BA* | 公链共识 |

---

## 参考文献

1. [Practical Byzantine Fault Tolerance](http://pmg.csail.mit.edu/papers/osdi99.pdf) - Castro & Liskov
2. [Byzantine Fault Tolerance](https://en.wikipedia.org/wiki/Byzantine_fault)
3. [The Byzantine Generals Problem](https://lamport.azurewebsites.net/pubs/byz.pdf) - Lamport et al.
