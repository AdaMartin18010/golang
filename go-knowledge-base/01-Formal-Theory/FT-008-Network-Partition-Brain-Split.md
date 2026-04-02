# FT-008: 网络分区与脑裂处理 (Network Partition & Brain Split Handling)

> **维度**: Formal Theory
> **级别**: S (18+ KB)
> **标签**: #network-partition #brain-split #split-brain #quorum
> **权威来源**: [Jepsen Tests](https://jepsen.io/), [CAP Theorem](https://sites.cs.ucsb.edu/~rich/class/cs293b-cloud/papers/brewer-cap.pdf)

---

## 网络分区类型

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Network Partition Types                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 简单分区 (Simple Partition)                                              │
│                                                                              │
│     [A]────[B]    [C]────[D]                                                │
│                                                                              │
│     A-B 互通, C-D 互通, A/B 与 C/D 不通                                      │
│                                                                              │
│  2. 非对称分区 (Asymmetric Partition)                                        │
│                                                                              │
│     [A]────►[B]    [C]                                                      │
│         ╲   │                                                                │
│          ╲  ▼                                                                │
│           [D]                                                                │
│                                                                              │
│     A 可以发送给 B, 但 B 无法回复 A                                          │
│                                                                              │
│  3. 延迟分区 (Latency Partition)                                             │
│                                                                              │
│     消息延迟 > timeout, 导致误判为分区                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 脑裂问题 (Split Brain)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Split Brain Scenario                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  分区前:                              分区后:                                │
│                                                                              │
│  ┌─────────┐                         ┌─────────┐    ╱╲    ┌─────────┐      │
│  │ Node A  │◄──────leader──────────►│ Node A  │◄───╱  ╲──►│ Node B  │      │
│  │ (Leader)│                         │ (Leader)│   分区   │ (New    │      │
│  └────┬────┘                         └────┬────┘          │  Leader)│      │
│       │                                   │               └────┬────┘      │
│       │                                   │                    │           │
│  ┌────┴────┐                         ┌────┴────┐          ┌────┴────┐      │
│  │ Node B  │                         │ Node C  │          │ Node D  │      │
│  │ Node C  │                         │         │          │         │      │
│  │ Node D  │                         └─────────┘          └─────────┘      │
│  └─────────┘                                                                │
│                                                                              │
│  问题：                                                                       │
│  • 两个分区都认为自己是 majority                                             │
│  • 可能选举出两个 leader                                                      │
│  • 导致数据不一致 (divergence)                                                │
│                                                                              │
│  解决：Quorum 机制 (多数派)                                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Quorum 机制

### 定义

$$
\begin{aligned}
&\text{Quorum: 读写操作需要的最小节点数} \\
&\text{Write Quorum: } W > \frac{N}{2} \\
&\text{Read Quorum: } R > \frac{N}{2} \\
&\text{Constraint: } W + R > N \\
\\
&\text{保证: 任何写入的 Quorum 和读取的 Quorum 至少有一个共同节点} \\
&\Rightarrow \text{读取一定能看到最新的写入}
\end{aligned}
$$

### Go 实现

```go
package quorum

import (
    "context"
    "errors"
    "sync"
)

// QuorumStore Quorum-based storage
type QuorumStore struct {
    nodes     []Node
    writeQuorum int
    readQuorum  int
}

func NewQuorumStore(nodes []Node) *QuorumStore {
    n := len(nodes)
    return &QuorumStore{
        nodes:       nodes,
        writeQuorum: n/2 + 1,  // 多数派
        readQuorum:  n/2 + 1,
    }
}

func (s *QuorumStore) Write(ctx context.Context, key, value string) error {
    acks := 0
    var mu sync.Mutex
    var lastErr error

    // 并行写入所有节点
    var wg sync.WaitGroup
    for _, node := range s.nodes {
        wg.Add(1)
        go func(n Node) {
            defer wg.Done()

            if err := n.Put(ctx, key, value); err != nil {
                mu.Lock()
                lastErr = err
                mu.Unlock()
                return
            }

            mu.Lock()
            acks++
            mu.Unlock()
        }(node)
    }

    // 等待达到 Quorum
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        if acks >= s.writeQuorum {
            return nil
        }
        return errors.New("write quorum not reached")
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (s *QuorumStore) Read(ctx context.Context, key string) (string, error) {
    responses := make(chan string, len(s.nodes))

    // 并行读取
    for _, node := range s.nodes {
        go func(n Node) {
            val, err := n.Get(ctx, key)
            if err == nil {
                responses <- val
            }
        }(node)
    }

    // 收集响应
    values := make(map[string]int)
    for i := 0; i < s.readQuorum; i++ {
        select {
        case val := <-responses:
            values[val]++
            if values[val] >= s.readQuorum {
                return val, nil
            }
        case <-ctx.Done():
            return "", ctx.Err()
        }
    }

    return "", errors.New("read quorum not reached")
}
```

---

## 分区检测与处理

### 检测机制

```go
// 心跳检测
type PartitionDetector struct {
    nodes     []string
    heartbeat map[string]time.Time
    timeout   time.Duration
}

func (d *PartitionDetector) CheckPartition() []string {
    now := time.Now()
    var partitioned []string

    for node, lastSeen := range d.heartbeat {
        if now.Sub(lastSeen) > d.timeout {
            partitioned = append(partitioned, node)
        }
    }

    return partitioned
}

// 误判处理：使用多个独立网络路径
func (d *PartitionDetector) CheckWithMultiplePaths(node string) bool {
    paths := []string{"tcp", "udp", "icmp"}
    success := 0

    for _, path := range paths {
        if d.ping(node, path) {
            success++
        }
    }

    // 多数路径成功则认为可达
    return success >= len(paths)/2+1
}
```

### 处理策略

| 策略 | 行为 | 适用场景 |
|------|------|---------|
| **Fail Fast** | 立即返回错误 | 金融交易 |
| **Degrade** | 只读模式 | 内容系统 |
| **Wait** | 等待分区恢复 | 短期分区 |
| **Merge** | 手动合并数据 | 长期分区后 |

---

## 参考文献

1. [Jepsen Tests](https://jepsen.io/) - Distributed Systems Safety Analysis
2. [The Part-Time Parliament](https://lamport.azurewebsites.net/pubs/lamport-paxos.pdf) - Lamport
3. [Dynamo: Amazon's Highly Available Key-value Store](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf)
