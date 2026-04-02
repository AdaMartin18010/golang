# FT-012: CRDT 无冲突复制数据类型 (Conflict-Free Replicated Data Types)

> **维度**: Formal Theory
> **级别**: S (17+ KB)
> **标签**: #crdt #eventual-consistency #collaborative-editing #distributed-data
> **权威来源**: [A comprehensive study of Convergent and Commutative Replicated Data Types](https://hal.inria.fr/file/index/docid/555588/filename/techreport.pdf)

---

## 核心概念

CRDT 是一种数据结构，在网络分区时允许多个副本独立更新，无需协调即可自动解决冲突。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CRDT Core Properties                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 强最终一致性 (Strong Eventual Consistency)                               │
│     - 无协调更新                                                             │
│     - 所有更新最终传播到所有副本                                              │
│     - 副本间无冲突（数学保证）                                                │
│                                                                              │
│  2. 合并性质                                                                  │
│     ┌─────────┐      ┌─────────┐                                            │
│     │ Node A  │      │ Node B  │                                            │
│     │  + {x}  │      │  + {y}  │                                            │
│     └────┬────┘      └────┬────┘                                            │
│          │                │                                                  │
│          └───────┬────────┘                                                  │
│                  ▼                                                           │
│            ┌─────────┐                                                       │
│            │ {x, y}  │  自动合并，无冲突                                      │
│            └─────────┘                                                       │
│                                                                              │
│  两种类型:                                                                   │
│  - CmRDT (Commutative/Op-based): 交换律保证                                   │
│  - CvRDT (Convergent/State-based): 合并函数保证                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## CRDT 类型

### 1. G-Counter (增长计数器)

```go
package crdt

import (
    "encoding/json"
    "sync"
)

// GCounter 只增计数器
type GCounter struct {
    mu     sync.RWMutex
    id     string
    counts map[string]uint64 // 每个节点的本地计数
}

// NewGCounter 创建计数器
func NewGCounter(id string) *GCounter {
    return &GCounter{
        id:     id,
        counts: make(map[string]uint64),
    }
}

// Increment 增加计数
func (c *GCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.counts[c.id]++
}

// Value 获取总值
func (c *GCounter) Value() uint64 {
    c.mu.RLock()
    defer c.mu.RUnlock()

    var total uint64
    for _, v := range c.counts {
        total += v
    }
    return total
}

// Merge 合并另一个 GCounter (取每个节点的最大值)
func (c *GCounter) Merge(other *GCounter) {
    c.mu.Lock()
    defer c.mu.Unlock()

    for node, count := range other.counts {
        if count > c.counts[node] {
            c.counts[node] = count
        }
    }
}

// State 获取状态 (用于网络传输)
func (c *GCounter) State() map[string]uint64 {
    c.mu.RLock()
    defer c.mu.RUnlock()

    state := make(map[string]uint64)
    for k, v := range c.counts {
        state[k] = v
    }
    return state
}
```

### 2. PN-Counter (可增可减计数器)

```go
// PNCounter 可增可减计数器
type PNCounter struct {
    increments *GCounter
    decrements *GCounter
}

// NewPNCounter 创建 PN-Counter
func NewPNCounter(id string) *PNCounter {
    return &PNCounter{
        increments: NewGCounter(id),
        decrements: NewGCounter(id),
    }
}

// Increment 增加
func (c *PNCounter) Increment() {
    c.increments.Increment()
}

// Decrement 减少
func (c *PNCounter) Decrement() {
    c.decrements.Increment()
}

// Value 获取净值
func (c *PNCounter) Value() int64 {
    return int64(c.increments.Value()) - int64(c.decrements.Value())
}

// Merge 合并
func (c *PNCounter) Merge(other *PNCounter) {
    c.increments.Merge(other.increments)
    c.decrements.Merge(other.decrements)
}
```

### 3. G-Set (只增集合)

```go
// GSet 只增集合
type GSet[T comparable] struct {
    mu   sync.RWMutex
    data map[T]struct{}
}

// NewGSet 创建集合
func NewGSet[T comparable]() *GSet[T] {
    return &GSet[T]{
        data: make(map[T]struct{}),
    }
}

// Add 添加元素
func (s *GSet[T]) Add(value T) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[value] = struct{}{}
}

// Contains 检查元素
func (s *GSet[T]) Contains(value T) bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    _, ok := s.data[value]
    return ok
}

// Values 获取所有值
func (s *GSet[T]) Values() []T {
    s.mu.RLock()
    defer s.mu.RUnlock()

    result := make([]T, 0, len(s.data))
    for v := range s.data {
        result = append(result, v)
    }
    return result
}

// Merge 合并
func (s *GSet[T]) Merge(other *GSet[T]) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for v := range other.data {
        s.data[v] = struct{}{}
    }
}
```

### 4. LWW-Element-Set (最后写入胜利集合)

```go
import "time"

// LWWSet 最后写入胜利集合 (支持删除)
type LWWSet[T comparable] struct {
    mu      sync.RWMutex
    adds    map[T]int64 // 添加时间戳
    removes map[T]int64 // 删除时间戳
}

// Add 添加元素
func (s *LWWSet[T]) Add(value T) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.adds[value] = time.Now().UnixNano()
}

// Remove 删除元素
func (s *LWWSet[T]) Remove(value T) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.removes[value] = time.Now().UnixNano()
}

// Contains 检查元素存在
func (s *LWWSet[T]) Contains(value T) bool {
    s.mu.RLock()
    defer s.mu.RUnlock()

    addTime, added := s.adds[value]
    removeTime, removed := s.removes[value]

    if !added {
        return false
    }
    if !removed {
        return true
    }
    return addTime > removeTime
}

// Merge 合并 (取每个操作的最大时间戳)
func (s *LWWSet[T]) Merge(other *LWWSet[T]) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for v, t := range other.adds {
        if ct, ok := s.adds[v]; !ok || t > ct {
            s.adds[v] = t
        }
    }
    for v, t := range other.removes {
        if ct, ok := s.removes[v]; !ok || t > ct {
            s.removes[v] = t
        }
    }
}
```

---

## 应用场景

| 系统 | 用途 | CRDT 类型 |
|------|------|----------|
| Redis | 计数器 | G-Counter |
| Riak | 购物车 | OR-Set |
| Figma | 协同编辑 | 序列 CRDT |
| Notion | 文档协同 | Yjs CRDT |
| Apple Notes | 多端同步 | LWW-Element-Set |

---

## 数学保证

```
CRDT 必须满足的性质:

1. 交换律 (Commutativity)
   A ⊔ B = B ⊔ A

2. 结合律 (Associativity)
   (A ⊔ B) ⊔ C = A ⊔ (B ⊔ C)

3. 幂等律 (Idempotence) [可选但推荐]
   A ⊔ A = A

其中 ⊔ 表示合并操作

强最终一致性定理:
如果所有更新都传播到所有副本，则所有副本状态收敛到相同值。
```

---

## 参考文献

1. [A comprehensive study of Convergent and Commutative Replicated Data Types](https://hal.inria.fr/file/index/docid/555588/filename/techreport.pdf)
2. [Conflict-free Replicated Data Types](https://arxiv.org/abs/1805.06358)
3. [CRDT.tech](https://crdt.tech/)
