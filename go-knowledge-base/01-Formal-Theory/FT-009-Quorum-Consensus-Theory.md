# FT-009: Quorum 共识理论 (Quorum Consensus Theory)

> **维度**: Formal Theory
> **级别**: S (15+ KB)
> **标签**: #quorum #consensus #distributed-systems #majority
> **权威来源**: [Weighted Voting for Replicated Data](https://dl.acm.org/doi/10.1145/358699.358703) - Gifford (1979)

---

## 核心概念

Quorum 是一组满足特定条件的节点集合，用于在分布式系统中安全地进行决策。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Quorum System                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  节点集合: N = {N1, N2, N3, N4, N5}                                          │
│                                                                              │
│  读 Quorum (R): {N1, N2, N3}                                                │
│  写 Quorum (W): {N3, N4, N5}                                                │
│                                                                              │
│  关键约束: R ∩ W ≠ ∅  (读写 Quorum 必须相交)                                  │
│                                                                              │
│  ┌───────────┐         ┌───────────┐                                        │
│  │     R     │         │     W     │                                        │
│  │  ┌─────┐  │         │  ┌─────┐  │                                        │
│  │  │ N1  │  │         │  │ N3  │◄─┼──── 交集: N3                           │
│  │  │ N2  │  │         │  │ N4  │  │                                        │
│  │  │ N3  │◄─┼─────────┼──┤ N5  │  │                                        │
│  │  └─────┘  │         │  └─────┘  │                                        │
│  └───────────┘         └───────────┘                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Quorum 类型

### 1. 简单多数 (Majority)

```
N = 5 节点
Quorum 大小 = ⌊N/2⌋ + 1 = 3

读 Quorum: 任意 3 个节点
写 Quorum: 任意 3 个节点

交集保证: |R| + |W| > N
         3 + 3 = 6 > 5 ✓

可用性: 最多容忍 (N-1)/2 = 2 个节点故障
```

### 2. 带权 Quorum (Weighted Voting)

```go
// 不同节点有不同权重
type Node struct {
    ID     string
    Weight int
}

// Quorum 要求总权重 > 50%
nodes := []Node{
    {"A", 3},  // 重要节点
    {"B", 2},
    {"C", 2},
    {"D", 1},  // 次要节点
}
// 总权重 = 8
// Quorum 需要 > 4

// 示例: {A, B} = 5 > 4, 形成 Quorum
```

### 3. 网格 Quorum (Grid Quorum)

```
将节点组织成网格，Quorum = 整行 + 整列

    1   2   3
   ┌───┬───┬───┐
1  │ A │ B │ C │
   ├───┼───┼───┤
2  │ D │ E │ F │
   ├───┼───┼───┤
3  │ G │ H │ I │
   └───┴───┴───┘

Quorum: {A, B, C} ∪ {C, F, I} = {A, B, C, F, I}
大小: 2√N - 1 = 5 (对于 N=9)

任意两个 Quorum 相交 (在网格的交叉点)
```

---

## 数学证明

### 安全性证明

```
定理: 如果 R ∩ W ≠ ∅，则任何读操作都能看到最近的写操作

证明:
1. 设 W1 是写操作 Quorum
2. 设 R 是后续读操作 Quorum
3. 因为 R ∩ W1 ≠ ∅，存在节点 n ∈ R ∩ W1
4. 节点 n 参与了 W1，存储了写入值
5. 读操作从 R 中读取，必然包含 n
6. 因此读操作能看到写入值

一致性保证: 读不会看到比已确认写更旧的值
```

---

## 应用示例

### Dynamo 风格读写

```go
package quorum

import (
    "context"
    "errors"
)

type DynamoStore struct {
    nodes     []Node
    readQuorum  int // R
    writeQuorum int // W
}

// Get 读取数据 (读取 R 个节点，返回最新版本)
func (s *DynamoStore) Get(ctx context.Context, key string) (Value, error) {
    responses := make(chan NodeResponse, len(s.nodes))

    // 并发读取所有节点
    for _, node := range s.nodes {
        go func(n Node) {
            val, err := n.Read(ctx, key)
            responses <- NodeResponse{Node: n, Value: val, Err: err}
        }(node)
    }

    // 等待 R 个成功响应
    var values []Value
    successCount := 0
    for resp := range responses {
        if resp.Err == nil {
            values = append(values, resp.Value)
            successCount++
            if successCount >= s.readQuorum {
                break
            }
        }
    }

    if successCount < s.readQuorum {
        return Value{}, errors.New("insufficient nodes available")
    }

    // 返回最新版本 (基于向量时钟)
    return s.resolveConflict(values), nil
}

// Put 写入数据 (写入 W 个节点)
func (s *DynamoStore) Put(ctx context.Context, key string, value Value) error {
    acks := make(chan error, len(s.nodes))

    // 并发写入
    for _, node := range s.nodes {
        go func(n Node) {
            acks <- n.Write(ctx, key, value)
        }(node)
    }

    // 等待 W 个确认
    successCount := 0
    for err := range acks {
        if err == nil {
            successCount++
            if successCount >= s.writeQuorum {
                return nil
            }
        }
    }

    return errors.New("write failed: insufficient acks")
}

// 配置示例: N=3, W=2, R=2 (平衡配置)
// 配置示例: N=3, W=1, R=3 (高写可用)
// 配置示例: N=3, W=3, R=1 (强一致性)
```

---

## 参考文献

1. [Weighted Voting for Replicated Data](https://dl.acm.org/doi/10.1145/358699.358703) - David K. Gifford
2. [Dynamo: Amazon's Highly Available Key-value Store](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf)
3. [Quorum Systems](https://www.cs.cornell.edu/home/rvr/papers/QSaA.pdf)
