# FT-010: 分布式时间、时钟与排序 (Time, Clocks and Ordering)

> **维度**: Formal Theory
> **级别**: S (16+ KB)
> **标签**: #logical-clocks #vector-clocks #happens-before #distributed-systems
> **权威来源**: [Time, Clocks, and the Ordering of Events in a Distributed System](https://lamport.azurewebsites.net/pubs/time-clocks.pdf) - Leslie Lamport (1978)

---

## 物理时钟问题

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Physical Clock Limitations                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  节点 A                     节点 B                                          │
│  ┌──────┐                 ┌──────┐                                          │
│  │T1=10 │──Event a───────►│T1=12 │                                          │
│  │      │                 │      │                                          │
│  │T2=20 │◄───Event b─────│T2=15 │  问题: 事件 b 发生在 a 之后，            │
│  │      │                 │      │        但物理时钟显示 T2(a) > T2(b)       │
│  └──────┘                 └──────┘                                          │
│                                                                              │
│  原因:                                                                       │
│  - 时钟漂移 (Clock Drift)                                                    │
│  - 网络延迟不确定性                                                          │
│  - 无法全局同步                                                              │
│                                                                              │
│  解决: 逻辑时钟 (Logical Clocks)                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Happens-Before 关系 (→)

### 定义

```
a → b (a happens-before b) 当且仅当:

1. 同一进程: 若 a 在 b 之前执行
   P: a ──► b  =>  a → b

2. 发送与接收: 若 a 是发送消息，b 是接收同一消息
   P1: a (send m) ──────► P2: b (receive m)  =>  a → b

3. 传递性: 若 a → b 且 b → c，则 a → c
   a → b → c  =>  a → c

4. 并发: 若 a ↛ b 且 b ↛ a，则 a || b (a 与 b 并发)
```

---

## Lamport 逻辑时钟

### 规则

```
每个进程 Pi 维护本地计数器 Ci:

1. 执行本地事件: Ci := Ci + 1

2. 发送消息 m:
   Ci := Ci + 1
   附加时间戳 C(m) = Ci

3. 接收消息 m:
   Ci := max(Ci, C(m)) + 1

示例:
P1: C1=1(a) ──m1──► P2: C2=2(b) ──m2──► P3: C2=4(d)
         │                      │
         │                      │
         ▼                      ▼
    C1=2(c)               C2=3(c)
```

### 局限性

```
Lamport 时钟可以判断: 若 C(a) < C(b)，则可能发生 a → b
但无法判断: 若 C(a) < C(b)，是否确实 a → b (可能是并发)

需要: Vector Clocks (向量时钟)
```

---

## 向量时钟 (Vector Clocks)

### 定义

```
每个进程 Pi 维护向量 Vi[1..n]:
- Vi[i]: Pi 的事件计数
- Vi[j]: Pi 感知的 Pj 的事件计数
```

### 规则

```
1. 本地事件: Vi[i] += 1

2. 发送消息 m:
   Vi[i] += 1
   附加向量 V(m) = Vi

3. 接收消息 m (来自 Pj):
   for k := 1 to n:
       Vi[k] = max(Vi[k], V(m)[k])
   Vi[i] += 1
```

### 比较规则

```
V(a) = V(b)  =>  a 与 b 是同一事件

V(a) < V(b)  =>  a → b (a 确定发生在 b 之前)
对所有 k: V(a)[k] ≤ V(b)[k] 且 存在 j: V(a)[j] < V(b)[j]

V(a) || V(b)  =>  a || b (并发)
存在 i, j: V(a)[i] < V(b)[i] 且 V(a)[j] > V(b)[j]
```

---

## Go 实现

```go
package clock

import (
    "encoding/json"
    "fmt"
)

// VectorClock 向量时钟
type VectorClock map[string]uint64

// New 创建向量时钟
func New(nodeID string) VectorClock {
    vc := make(VectorClock)
    vc[nodeID] = 0
    return vc
}

// Increment 递增本地时钟
func (vc VectorClock) Increment(nodeID string) {
    vc[nodeID]++
}

// Merge 合并另一个向量时钟
func (vc VectorClock) Merge(other VectorClock) {
    for node, time := range other {
        if vc[node] < time {
            vc[node] = time
        }
    }
}

// HappensBefore 判断 happens-before 关系
func (vc VectorClock) HappensBefore(other VectorClock) bool {
    dominates := false
    for node, time := range other {
        if vc[node] > time {
            return false // 不满足小于等于
        }
        if vc[node] < time {
            dominates = true // 至少有一个严格小于
        }
    }
    return dominates
}

// Concurrent 判断并发关系
func (vc VectorClock) Concurrent(other VectorClock) bool {
    return !vc.HappensBefore(other) && !other.HappensBefore(vc)
}

// String 字符串表示
func (vc VectorClock) String() string {
    data, _ := json.Marshal(vc)
    return string(data)
}

// 使用示例
func Example() {
    // P1 创建时钟
    vc1 := New("P1")
    vc1.Increment("P1") // P1 执行事件 a

    // P1 发送消息给 P2
    msgClock := vc1 // 复制时钟

    // P2 接收消息
    vc2 := New("P2")
    vc2.Merge(msgClock)
    vc2.Increment("P2") // P2 执行事件 b

    // 判断关系
    fmt.Println(vc1.HappensBefore(vc2)) // true: a → b
    fmt.Println(vc2.HappensBefore(vc1)) // false
}
```

---

## 应用场景

| 场景 | 时钟类型 | 用途 |
|------|---------|------|
| 因果一致性 | Vector Clock | 判断事件因果 |
| 版本向量 | Vector Clock | 检测冲突 |
| 乐观锁 | Lamport Clock | 简单排序 |
| Dynamo | Vector Clock | 冲突解决 |
| CockroachDB | Hybrid Logical Clock | 全局排序 |

---

## 参考文献

1. [Time, Clocks, and the Ordering of Events in a Distributed System](https://lamport.azurewebsites.net/pubs/time-clocks.pdf) - Leslie Lamport
2. [Dynamo: Amazon's Highly Available Key-value Store](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf)
