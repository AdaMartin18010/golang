# FT-006: 向量时钟与逻辑时间 (Vector Clocks & Logical Time)

> **维度**: Formal Theory
> **级别**: S (20+ KB)
> **标签**: #vector-clocks #logical-time #lamport-timestamp #causality
> **权威来源**: [Time, Clocks, and the Ordering of Events](https://lamport.azurewebsites.net/pubs/time-clocks.pdf) - Lamport, 1978, [Dynamo Paper](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf)

---

## 逻辑时间基础

### 物理时钟的问题

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Physical Clock Limitations                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Node A (时钟快 100ms)                Node B (时钟慢 100ms)                 │
│       │                                    │                                 │
│       │ T=100: 发送消息                   │                                  │
│       │────────────────────────────────────►│ T=0: 接收消息                  │
│       │                                    │                                 │
│       │                                    │ 接收时间 < 发送时间？           │
│       │                                    │ 违反因果律！                     │
│                                                                              │
│  问题：                                                                      │
│  1. 时钟漂移 (Clock Drift)                                                   │
│  2. 网络延迟不确定性                                                         │
│  3. 无法确定事件因果关系                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Lamport 时间戳

### 定义

$$
\begin{aligned}
&\text{Lamport Timestamp Rules:} \\
&1. \text{Local event: } C_i := C_i + 1 \\
&2. \text{Send message: } C_i := C_i + 1; \text{ timestamp} = C_i \\
&3. \text{Receive message: } C_i := \max(C_i, \text{timestamp}_{msg}) + 1 \\
\\
&\text{Happens-Before Relation:} \\
&a \rightarrow b \Rightarrow C(a) < C(b) \\
&\text{But: } C(a) < C(b) \nRightarrow a \rightarrow b \text{ (并发事件可能)} \\
\end{aligned}
$$

### Go 实现

```go
package logicaltime

import (
    "sync"
)

// LamportClock Lamport 逻辑时钟
type LamportClock struct {
    mu sync.Mutex
    timestamp int64
}

// Tick 本地事件，时钟前进
func (c *LamportClock) Tick() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.timestamp++
    return c.timestamp
}

// Send 发送消息，返回时间戳
func (c *LamportClock) Send() int64 {
    return c.Tick()
}

// Receive 接收消息，更新时钟
func (c *LamportClock) Receive(msgTimestamp int64) int64 {
    c.mu.Lock()
    defer c.mu.Unlock()

    if msgTimestamp > c.timestamp {
        c.timestamp = msgTimestamp
    }
    c.timestamp++
    return c.timestamp
}

// Current 获取当前时间戳
func (c *LamportClock) Current() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.timestamp
}
```

---

## 向量时钟 (Vector Clocks)

### 定义

$$
\begin{aligned}
&\text{Vector Clock } V_i = [v_{i1}, v_{i2}, ..., v_{in}] \\
&\text{where } v_{ij} = \text{node } i \text{'s knowledge of node } j \text{'s events} \\
\\
&\text{Rules:} \\
&1. \text{Local event: } V_i[i]++ \\
&2. \text{Send: } V_i[i]++; \text{ send } V_i \\
&3. \text{Receive: } V_i[k] = \max(V_i[k], V_{msg}[k]) \forall k; V_i[i]++ \\
\\
&\text{Comparison:} \\
&V_a \leq V_b \Leftrightarrow \forall i: V_a[i] \leq V_b[i] \\
&V_a < V_b \Leftrightarrow V_a \leq V_b \land \exists i: V_a[i] < V_b[i] \\
&V_a \parallel V_b \Leftrightarrow \neg(V_a < V_b) \land \neg(V_b < V_a) \\
\end{aligned}
$$

### Go 实现

```go
package logicaltime

import (
    "fmt"
    "sync"
)

// VectorClock 向量时钟
type VectorClock struct {
    mu     sync.RWMutex
    nodeID int
    vector map[int]int64  // nodeID -> timestamp
}

// NewVectorClock 创建向量时钟
func NewVectorClock(nodeID int, numNodes int) *VectorClock {
    vc := &VectorClock{
        nodeID: nodeID,
        vector: make(map[int]int64, numNodes),
    }
    for i := 0; i < numNodes; i++ {
        vc.vector[i] = 0
    }
    return vc
}

// Tick 本地事件
func (vc *VectorClock) Tick() {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    vc.vector[vc.nodeID]++
}

// Send 发送消息前调用，返回时钟副本
func (vc *VectorClock) Send() map[int]int64 {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    vc.vector[vc.nodeID]++

    // 返回副本
    copy := make(map[int]int64, len(vc.vector))
    for k, v := range vc.vector {
        copy[k] = v
    }
    return copy
}

// Receive 接收消息
func (vc *VectorClock) Receive(other map[int]int64) {
    vc.mu.Lock()
    defer vc.mu.Unlock()

    // 逐元素取最大值
    for node, timestamp := range other {
        if timestamp > vc.vector[node] {
            vc.vector[node] = timestamp
        }
    }
    vc.vector[vc.nodeID]++
}

// Compare 比较两个向量时钟
// Returns: -1 if vc < other, 0 if concurrent, 1 if vc > other
func (vc *VectorClock) Compare(other map[int]int64) int {
    vc.mu.RLock()
    defer vc.mu.RUnlock()

    less := false
    greater := false

    for node, v1 := range vc.vector {
        v2, exists := other[node]
        if !exists {
            // other 没有该节点信息，视为 0
            v2 = 0
        }

        if v1 < v2 {
            less = true
        } else if v1 > v2 {
            greater = true
        }
    }

    // 检查 other 中有但 vc 中没有的节点
    for node, v2 := range other {
        if _, exists := vc.vector[node]; !exists {
            if v2 > 0 {
                less = true
            }
        }
    }

    if less && !greater {
        return -1  // vc < other
    } else if !less && greater {
        return 1   // vc > other
    } else if !less && !greater {
        return 0   // equal
    } else {
        return 0   // concurrent (不可比较)
    }
}

// String 字符串表示
func (vc *VectorClock) String() string {
    vc.mu.RLock()
    defer vc.mu.RUnlock()
    return fmt.Sprintf("VC[%d]: %v", vc.nodeID, vc.vector)
}

// GetCopy 获取时钟副本
func (vc *VectorClock) GetCopy() map[int]int64 {
    vc.mu.RLock()
    defer vc.mu.RUnlock()

    copy := make(map[int]int64, len(vc.vector))
    for k, v := range vc.vector {
        copy[k] = v
    }
    return copy
}
```

---

## 向量时钟应用：冲突检测

```go
// 版本向量（Dynamo, Cassandra）
type VersionedValue struct {
    Value     string
    VectorClock VectorClock
}

// 解决冲突
func ResolveConflict(values []VersionedValue) (VersionedValue, bool) {
    if len(values) == 0 {
        return VersionedValue{}, false
    }

    if len(values) == 1 {
        return values[0], false
    }

    // 检查因果关系
    causallyRelated := false
    for i := 0; i < len(values); i++ {
        for j := i + 1; j < len(values); j++ {
            cmp := values[i].VectorClock.Compare(values[j].VectorClock.GetCopy())
            if cmp != 0 {
                causallyRelated = true
            }

            if cmp < 0 {  // i < j
                return values[j], false  // j 是更新的版本
            } else if cmp > 0 {  // i > j
                return values[i], false  // i 是更新的版本
            }
        }
    }

    // 并发冲突，需要业务层解决
    return VersionedValue{}, true
}
```

---

## 向量时钟的优化

### 1. 版本向量剪枝

```
问题：向量时钟随节点数增长，存储开销大

解决方案：
1. 只保留最近 N 个版本
2. 合并旧的向量时钟
3. 使用 Merkle Tree 压缩
```

### 2. Dynamo 的简化向量时钟

```
Dynamo 使用 (node, counter) 对列表代替完整向量：
- 当节点数 > 阈值时，移除最旧的条目
- 可能导致假阳性冲突（可以接受）
```

---

## 时间戳对比

| 特性 | Lamport | Vector Clocks | Version Vectors |
|------|---------|---------------|-----------------|
| 精度 | 偏序 | 全序（可检测并发） | 全序 |
| 空间 | O(1) | O(n) | O(n) |
| 检测并发 | ❌ | ✅ | ✅ |
| 适用场景 | 简单顺序 | 冲突检测 | 版本控制 |

---

## 参考文献

1. [Time, Clocks, and the Ordering of Events in a Distributed System](https://lamport.azurewebsites.net/pubs/time-clocks.pdf) - Lamport, 1978
2. [Dynamo: Amazon's Highly Available Key-value Store](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) - DeCandia et al., 2007
3. [Why Vector Clocks are Easy](https://riak.com/posts/technical/why-vector-clocks-are-easy/) - Riak
4. [Version Vectors are Easy](http://gsd.di.uminho.pt/members/cbm/ps/itc2008.pdf) - Almeida et al.
