# FT-012: 因果一致性的形式化理论 (Causal Consistency: Formal Theory)

> **维度**: Formal Theory
> **级别**: S (21+ KB)
> **标签**: #causal-consistency #vector-clocks #happens-before #eventual-consistency
> **权威来源**:
>
> - [Time, Clocks, and the Ordering of Events](https://amturing.acm.org/bib/lamport_1978_time.pdf) - Lamport (1978)
> - [Causal Memory: Definitions, Implementation, and Programming](https://www.vs.inf.ethz.ch/publ/papers/caumatechreport.pdf) - Ahamad et al. (1995)
> - [COPS: The Scalable Causal Consistency Platform](https://www.cs.cmu.edu/~dga/papers/cops-sosp2011.pdf) - Lloyd et al. (2011)
> - [Bolt-on Causal Consistency](https://www.cs.cmu.edu/~pavlo/courses/fall2013/static/papers/bailis2013bolton.pdf) - Bailis et al. (2013)
> - [The Complexity of Transactional Causal Consistency](https://arxiv.org/abs/1503.07687) - Brutschy et al. (2017)

---

## 1. 因果一致性的形式化定义

### 1.1 Happens-Before 关系

**定义 1.1 (Happens-Before $\prec$ - Lamport 1978)**

$$
\begin{aligned}
&\text{(程序序)}: &&o_1, o_2 \text{ 在同一进程且 } o_1 \text{ 先于 } o_2 \Rightarrow o_1 \prec o_2 \\
&\text{(读-从)}: &&o_1 = \text{write}(x, v) \land o_2 = \text{read}(x) \rightarrow v \Rightarrow o_1 \prec o_2 \\
&\text{(传递性)}: &&o_1 \prec o_2 \land o_2 \prec o_3 \Rightarrow o_1 \prec o_3
\end{aligned}
$$

**定义 1.2 (并发操作)**

$$o_1 \parallel o_2 \Leftrightarrow \neg(o_1 \prec o_2) \land \neg(o_2 \prec o_1)$$

### 1.2 因果一致性定义

**定义 1.3 (因果一致性 - Ahamad et al. 1995)**

一个执行是因果一致的，如果：

1. **因果序保持**: 如果 $o_1 \prec o_2$，则所有进程看到的 $o_1$ 都在 $o_2$ 之前
2. **写收敛**: 所有进程最终看到相同的写顺序
3. **读值正确**: 读操作返回因果序中最近的写

形式化：

$$\text{CausalConsistency} \equiv \forall p_i, \forall o_1, o_2:$$
$$(o_1 \prec o_2 \Rightarrow \text{visible}_{p_i}(o_1) \text{ before } \text{visible}_{p_i}(o_2))$$

### 1.3 与相关模型的关系

**定理 1.1 (蕴含层次)**

$$\text{Linearizable} \Rightarrow \text{Sequential} \Rightarrow \text{Causal} \Rightarrow \text{Eventual}$$

*证明概要*:

- Linearizable $\Rightarrow$ Sequential: 线性化序保持程序序
- Sequential $\Rightarrow$ Causal: 全局序包含因果序
- Causal $\Rightarrow$ Eventual: 因果一致包含收敛

$\square$

**定理 1.2 (与顺序一致性的关系)**

$$\text{Causal} \not\Rightarrow \text{Sequential} \land \text{Sequential} \not\Rightarrow \text{Causal}$$

两者是正交的：

- 因果一致性不保证程序序中的读-写顺序
- 顺序一致性不保证因果关系

---

## 2. 向量时钟实现

### 2.1 向量时钟定义

**定义 2.1 (向量时钟)**
对于 $n$ 个进程的系统中，向量时钟是 $n$ 维向量：

$$VC: \Pi \rightarrow \mathbb{N}^n$$

其中 $VC(p_i)[j]$ 表示进程 $p_i$ 知道的 $p_j$ 的事件数。

**定义 2.2 (向量时钟比较)**

$$
\begin{aligned}
VC_1 \leq VC_2 &\Leftrightarrow \forall i: VC_1[i] \leq VC_2[i] \\
VC_1 < VC_2 &\Leftrightarrow VC_1 \leq VC_2 \land VC_1 \neq VC_2 \\
VC_1 \parallel VC_2 &\Leftrightarrow \neg(VC_1 \leq VC_2) \land \neg(VC_2 \leq VC_1)
\end{aligned}
$$

**定理 2.1 (向量时钟与 Happens-Before)**

$$e_1 \prec e_2 \Leftrightarrow VC(e_1) < VC(e_2)$$

$$e_1 \parallel e_2 \Leftrightarrow VC(e_1) \parallel VC(e_2)$$

### 2.2 向量时钟算法

**算法 1: 更新规则**

```
// 本地事件
VC[i] += 1

// 发送消息
VC[i] += 1
send(message, VC)

// 接收消息
VC[i] += 1
VC = max(VC, received_VC)  // 逐分量取最大值
```

**定理 2.2 (向量时钟正确性)**
上述算法正确维护 Happens-Before 关系。

*证明*:

1. **程序序**: 本地事件递增保证了同一进程的因果关系
2. **读-从**: 消息传递合并向量时钟保证了跨进程因果关系
3. **传递性**: max 操作保证了传递闭包

$\square$

### 2.3 因果一致性实现

**COPS 算法 (Lloyd et al. 2011)**

**写操作**:

1. 客户端发送写请求到本地数据中心
2. 分配版本号 (依赖向量)
3. 本地提交，异步复制

**读操作**:

1. 读取本地副本
2. 如果依赖未满足，阻塞或返回旧版本
3. 版本选择确保因果一致性

**依赖追踪**:

- 每个值携带其因果依赖
- 读取值时同时获取其依赖
- 确保依赖在读取前可见

---

## 3. 正确性证明

### 3.1 因果一致性条件

**定理 3.1 (因果一致性判定)**
一个执行是因果一致的当且仅当：

1. **本地顺序**: 每个进程的操作按程序序可见
2. **写-写因果**: 因果相关的写按相同顺序被所有进程看到
3. **读-写因果**: 进程读取的值 $v$ 后，必然看到 $v$ 的所有因果前驱

*证明*:

$(\Rightarrow)$: 因果一致性定义直接蕴含这三个条件。

$(\Leftarrow)$:

- 本地顺序保证程序序保持
- 写-写因果保证因果序的全局一致
- 读-写因果保证因果前驱的可见性

$\square$

### 3.2 收敛性

**定理 3.2 (因果一致性收敛)**
在因果一致的系统中，如果停止写入，所有副本最终一致。

*证明*:

1. 因果一致性要求写收敛
2. 所有进程最终看到相同的写顺序
3. 停止写入后，无新值产生
4. 因此所有副本收敛到相同状态

$\square$

---

## 4. 多元表征

### 4.1 因果一致性概念地图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Causal Consistency Concept Network                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                    ┌─────────────────────────┐                              │
│                    │    Causal Consistency   │                              │
│                    │    (Happens-Before)     │                              │
│                    └───────────┬─────────────┘                              │
│                                │                                            │
│              ┌─────────────────┼─────────────────┐                         │
│              ▼                 ▼                 ▼                         │
│       ┌───────────┐    ┌───────────┐    ┌───────────┐                     │
│       │ Program   │    │  Read-    │    │   Write   │                     │
│       │ Order     │    │  From     │    │  Order    │                     │
│       │ Preserved │    │  Relation │    │  Converges│                     │
│       └─────┬─────┘    └─────┬─────┘    └─────┬─────┘                     │
│             │                │                │                            │
│       ┌─────┴─────┐    ┌─────┴─────┐  ┌─────┴─────┐                       │
│       │           │    │           │  │           │                       │
│       │ Per-      │    │ Read sees │  │ All procs │                       │
│       │ process   │    │ the write │  │ see same  │                       │
│       │ causal    │    │ it reads  │  │ write     │                       │
│       │ order     │    │ from      │  │ order     │                       │
│       └───────────┘    └───────────┘  └───────────┘                       │
│                                                                              │
│  Vector Clock Implementation:                                                │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                                                                     │    │
│  │   P1: [1,0,0] ──► [2,0,0] ─────────► [3,1,0]                       │    │
│  │              local       send to P2      receive from P2            │    │
│  │                    m                     merge(max)                 │    │
│  │                                                                     │    │
│  │   P2: [0,0,0] ─────────► [2,1,0] ───► [2,2,0]                      │    │
│  │              receive     local                                           │    │
│  │              merge(max)                                                  │    │
│  │                                                                     │    │
│  │   Causality: VC(a) < VC(b) ⟺ a happens-before b                     │    │
│  │   Concurrency: VC(a) || VC(b) ⟺ a and b are concurrent              │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  Key Properties:                                                             │
│  ├── Available under partitions (AP systems)                                 │
│  ├── Convergent (all replicas eventually agree)                              │
│  ├── Local property (can be verified per-object)                             │
│  └── Suitable for social networking, collaborative editing                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 一致性模型决策树

```
设计分布式系统?
│
├── 需要实时顺序?
│   ├── 是 → Linearizable / Sequential
│   │       └── 配置服务、锁服务
│   │
│   └── 否 → 继续
│
├── 需要因果顺序?
│   ├── 是 → Causal Consistency
│   │       ├── 需要会话保证?
│   │       │   ├── 是 → Read-Your-Writes + Monotonic-Reads
│   │       │   └── 否 → Pure Causal
│   │       │
│   │       └── 应用场景:
│   │           ├── 社交网络 → 发帖/评论因果可见
│   │           ├── 协作编辑 → 操作因果排序
│   │           └── 消息系统 → 消息因果传递
│   │
│   └── 否 → 继续
│
├── 需要单调读?
│   ├── 是 → Monotonic Reads
│   │       └── 用户会话内保证
│   │
│   └── 否 → 继续
│
└── 最终一致足够?
    └── 是 → Eventual Consistency
        └── DNS、CDN、日志分析
```

### 4.3 因果系统对比矩阵

| 属性 | COPS | ChainReaction | Eiger | GentleRain | Cassandra (Causal) |
|------|------|---------------|-------|------------|-------------------|
| **依赖追踪** | 显式向量 | 隐式链 | 显式向量 | 物理时钟 | 向量时钟 |
| **读延迟** | 本地 | 可能远程 | 本地 | 本地 | 本地 |
| **写延迟** | 本地 | 链式传播 | 本地 | 本地 | 本地 |
| **复制** | 异步 | 同步链 | 异步 | 异步 | 异步 |
| **冲突解决** | 最后写胜出 | 版本向量 | 最后写胜出 | 时间戳 | 向量时钟 |
| **可见性** | 显式依赖 | 链式 | 显式依赖 | 时间戳 | 向量时钟 |
| **可扩展性** | 高 | 中 | 高 | 高 | 高 |
| **应用场景** | 全球存储 | 缓存 | 键值存储 | 全球存储 | 通用 |

### 4.4 向量时钟示例图

```
时间 →

进程 P1:  [1,0,0] ──► [2,0,0] ──────────────────────► [3,1,0]
          event a     event b (send m1 to P2)         event c (recv m2)

          m1: [2,0,0]
          ─────────────────────────────────────────────►


          ◄─────────────────────────────────────────────
          m2: [2,1,0]

进程 P2:  [0,0,0] ────────────────► [2,1,0] ──► [2,2,0] ──► [3,2,0]
                         (recv m1)   event d   event e    (send m2)
                         merge:
                         max([0,0,0], [2,0,0])
                         + local inc

进程 P3:  [0,0,0] ───────────────────────────────────────────────► [0,0,1]
                                                                   event f

Causality Analysis:
────────────────────
a: [1,0,0] ◄──── b: [2,0,0] ◄──── c: [3,1,0]
                    │
                    ▼
              d: [2,1,0] ◄──── e: [2,2,0]
                                    │
                                    ▼
                              c: [3,1,0] (recv m2: max([3,1,0],[2,2,0])=[3,2,0]?)

Correction: c's VC after receiving m2:
  c receives m2 ([2,2,0]) with current [3,1,0]
  merge: max([3,1,0], [2,2,0]) = [3,2,0]
  local inc: [3,3,0]

Concurrency:
────────────────────
f: [0,0,1]  ||  c: [3,3,0]
  (neither [0,0,1] ≤ [3,3,0] nor [3,3,0] ≤ [0,0,1])
  → f and c are concurrent!

This means P3's event f is not causally related to P1's event c.
```

---

## 5. TLA+ 形式化规约

```tla
------------------------------- MODULE CausalConsistency ----------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS
    Process,            \* 进程集合
    Location,           \* 存储位置
    Value,              \* 值域
    Nil

VARIABLES
    operations,         \* 操作集合
    happensBefore,      \* Happens-Before 关系
    vectorClocks,       \* 向量时钟
    visibleOrder,       \* 每个进程看到的顺序
    store               \* 存储状态

causalVars == <<operations, happensBefore, vectorClocks, visibleOrder, store>>

-----------------------------------------------------------------------------
\* Happens-Before 关系定义

\* 1. 程序序: 同一进程的操作顺序
ProgramOrder(op1, op2) ==
    /\ op1.proc = op2.proc
    /\ op1.timestamp < op2.timestamp

\* 2. 读-从关系: 读从写读取
ReadFrom(readOp, writeOp) ==
    /\ readOp.type = "read"
    /\ writeOp.type = "write"
    /\ readOp.loc = writeOp.loc
    /\ readOp.value = writeOp.value

\* 3. Happens-Before 是传递闭包
HappensBefore(op1, op2) ==
    \/ ProgramOrder(op1, op2)
    \/ ReadFrom(op1, op2)
    \/ \E op3: HappensBefore(op1, op3) /\ HappensBefore(op3, op2)

-----------------------------------------------------------------------------
\* 向量时钟操作

\* 比较两个向量时钟
VCLess(vc1, vc2) ==
    /\ \A p \in Process: vc1[p] <= vc2[p]
    /\ \E p \in Process: vc1[p] < vc2[p]

VCConcurrent(vc1, vc2) ==
    /\ ~VCLess(vc1, vc2)
    /\ ~VCLess(vc2, vc1)

\* 向量时钟与 Happens-Before 对应
VCConsistent ==
    \A op1, op2 \in operations:
        HappensBefore(op1, op2) => VCLess(vectorClocks[op1], vectorClocks[op2])

-----------------------------------------------------------------------------
\* 因果一致性条件

\* 1. 因果序保持
PreservesCausalOrder ==
    \A p \in Process:
        \A op1, op2 \in operations:
            HappensBefore(op1, op2) =>
                InVisibleOrder(p, op1, op2)

\* 辅助: op1 在进程 p 的可见顺序中先于 op2
InVisibleOrder(p, op1, op2) ==
    \E i, j \in 1..Len(visibleOrder[p]):
        /\ i < j
        /\ visibleOrder[p][i] = op1
        /\ visibleOrder[p][j] = op2

\* 2. 写收敛
WriteConvergence ==
    \A w1, w2 \in {op \in operations: op.type = "write"}:
        \A p1, p2 \in Process:
            /\ HappensBefore(w1, w2) \/ HappensBefore(w2, w1)
            => OrderConsistent(p1, p2, w1, w2)

\* 辅助: 两个进程对两个写的顺序一致
OrderConsistent(p1, p2, w1, w2) ==
    (InVisibleOrder(p1, w1, w2) => InVisibleOrder(p2, w1, w2))
    /\ (InVisibleOrder(p1, w2, w1) => InVisibleOrder(p2, w2, w1))

\* 3. 读值正确
ReadValueCorrect ==
    \A r \in {op \in operations: op.type = "read"}:
        LET writes == {w \in operations:
                        w.type = "write" /\
                        w.loc = r.loc /\
                        HappensBefore(w, r)}
        IN writes = {} => r.value = Nil
           \/ \E w \in writes:
               /\ r.value = w.value
               /\ ~\E w2 \in writes:
                   w2 # w /\ HappensBefore(w, w2) /\ HappensBefore(w2, r)

-----------------------------------------------------------------------------
\* 因果一致性判定

IsCausallyConsistent ==
    /\ PreservesCausalOrder
    /\ WriteConvergence
    /\ ReadValueCorrect
    /\ VCConsistent

=============================================================================
```

---

## 6. Go 实现

### 6.1 向量时钟实现

```go
package causal

import (
    "fmt"
    "sync"
)

// VectorClock 向量时钟
type VectorClock map[string]uint64

// NewVectorClock 创建新向量时钟
func NewVectorClock(processes []string) VectorClock {
    vc := make(VectorClock)
    for _, p := range processes {
        vc[p] = 0
    }
    return vc
}

// Increment 增加本地时钟
func (vc VectorClock) Increment(process string) {
    vc[process]++
}

// Merge 合并另一个向量时钟
func (vc VectorClock) Merge(other VectorClock) {
    for process, timestamp := range other {
        if vc[process] < timestamp {
            vc[process] = timestamp
        }
    }
}

// Less 检查是否小于另一个向量时钟
func (vc VectorClock) Less(other VectorClock) bool {
    lessOrEqual := true
    strictlyLess := false

    for process, timestamp := range vc {
        if timestamp > other[process] {
            lessOrEqual = false
            break
        }
        if timestamp < other[process] {
            strictlyLess = true
        }
    }

    return lessOrEqual && strictlyLess
}

// Concurrent 检查是否与另一个向量时钟并发
func (vc VectorClock) Concurrent(other VectorClock) bool {
    return !vc.Less(other) && !other.Less(vc)
}

// Copy 复制向量时钟
func (vc VectorClock) Copy() VectorClock {
    copy := make(VectorClock)
    for k, v := range vc {
        copy[k] = v
    }
    return copy
}

func (vc VectorClock) String() string {
    return fmt.Sprintf("%v", map[string]uint64(vc))
}

// CausalStore 因果一致性存储
type CausalStore struct {
    id      string
    data    map[string]*VersionedValue
    vc      VectorClock
    peers   []string
    mu      sync.RWMutex
}

// VersionedValue 带版本号的值
type VersionedValue struct {
    Value     interface{}
    Version   VectorClock
    DependsOn []VectorClock // 因果依赖
}

// NewCausalStore 创建因果存储
func NewCausalStore(id string, peers []string) *CausalStore {
    allProcesses := append([]string{id}, peers...)
    return &CausalStore{
        id:    id,
        data:  make(map[string]*VersionedValue),
        vc:    NewVectorClock(allProcesses),
        peers: peers,
    }
}

// Write 写入数据
func (s *CausalStore) Write(key string, value interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // 增加本地时钟
    s.vc.Increment(s.id)

    // 创建版本化值
    versioned := &VersionedValue{
        Value:     value,
        Version:   s.vc.Copy(),
        DependsOn: []VectorClock{},
    }

    s.data[key] = versioned
}

// Read 读取数据
func (s *CausalStore) Read(key string) (interface{}, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    val, exists := s.data[key]
    if !exists {
        return nil, false
    }

    // 检查依赖是否满足
    for _, dep := range val.DependsOn {
        if !s.dependencySatisfied(dep) {
            // 依赖未满足，可能需要等待或返回旧版本
            return nil, false
        }
    }

    return val.Value, true
}

// dependencySatisfied 检查依赖是否满足
func (s *CausalStore) dependencySatisfied(dep VectorClock) bool {
    // 简化: 检查本地时钟是否覆盖依赖
    for process, timestamp := range dep {
        if s.vc[process] < timestamp {
            return false
        }
    }
    return true
}

// Replicate 复制数据到远程
func (s *CausalStore) Replicate(remote *CausalStore) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    remote.mu.Lock()
    defer remote.mu.Unlock()

    // 合并向量时钟
    remote.vc.Merge(s.vc)

    // 复制数据
    for key, value := range s.data {
        remote.data[key] = &VersionedValue{
            Value:     value.Value,
            Version:   value.Version.Copy(),
            DependsOn: copyDepends(value.DependsOn),
        }
    }
}

func copyDepends(deps []VectorClock) []VectorClock {
    result := make([]VectorClock, len(deps))
    for i, dep := range deps {
        result[i] = dep.Copy()
    }
    return result
}

// GetVectorClock 获取当前向量时钟
func (s *CausalStore) GetVectorClock() VectorClock {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.vc.Copy()
}
```

### 6.2 因果一致性验证

```go
package causal

import (
    "testing"
    "time"
)

func TestVectorClockComparison(t *testing.T) {
    vc1 := VectorClock{"p1": 1, "p2": 0, "p3": 0}
    vc2 := VectorClock{"p1": 2, "p2": 1, "p3": 0}
    vc3 := VectorClock{"p1": 1, "p2": 0, "p3": 1}

    // vc1 < vc2
    if !vc1.Less(vc2) {
        t.Error("vc1 should be less than vc2")
    }

    // vc1 || vc3 (concurrent)
    if !vc1.Concurrent(vc3) {
        t.Error("vc1 and vc3 should be concurrent")
    }
}

func TestCausalStore(t *testing.T) {
    store1 := NewCausalStore("p1", []string{"p2"})
    store2 := NewCausalStore("p2", []string{"p1"})

    // P1 写入
    store1.Write("key1", "value1")

    // 复制到 P2
    store1.Replicate(store2)

    // P2 应该能看到值
    val, ok := store2.Read("key1")
    if !ok || val != "value1" {
        t.Error("Replication failed")
    }

    // P2 写入
    store2.Write("key2", "value2")

    // P1 应该看不到 (尚未复制)
    _, ok = store1.Read("key2")
    if ok {
        t.Error("Should not see un-replicated value")
    }
}

func TestHappensBefore(t *testing.T) {
    // 模拟:
    // P1: write(x) ──► read(y)
    // P2: write(y) (由 P1 的读触发)

    store1 := NewCausalStore("p1", []string{"p2"})
    store2 := NewCausalStore("p2", []string{"p1"})

    // P1 写入 x
    store1.Write("x", "x1")

    // 复制到 P2
    store1.Replicate(store2)

    // P2 读取 x 后写入 y
    if val, ok := store2.Read("x"); ok {
        store2.Write("y", "y1-triggered-by-"+val.(string))
    }

    // P2 的 y 写入依赖于 P1 的 x 写入
    // 复制回 P1
    store2.Replicate(store1)

    // 验证因果链
    vc1 := store1.GetVectorClock()
    if vc1["p1"] == 0 || vc1["p2"] == 0 {
        t.Error("Causal chain not established")
    }
}
```

---

## 7. 学术参考文献

### 7.1 经典论文

1. **Lamport, L. (1978)**. Time, Clocks, and the Ordering of Events in a Distributed System. *Communications of the ACM*, 21(7), 558-565.
   - Happens-Before 关系的奠基性工作

2. **Ahamad, M., Neiger, G., Burns, J. E., Kohli, P., & Hutto, P. W. (1995)**. Causal Memory: Definitions, Implementation, and Programming. *Distributed Computing*, 9(1), 37-49.
   - 因果一致性的形式化定义

### 7.2 现代系统

1. **Lloyd, W., et al. (2011)**. Don't Settle for Eventual: Scalable Causal Consistency for Wide-Area Storage with COPS. *SOSP*.
   - COPS 系统

2. **Bailis, P., et al. (2013)**. Bolt-on Causal Consistency. *SIGMOD*.
   - 在现有存储上实现因果一致性

3. **Du, J., et al. (2013)**. Generalized Fence-Preserving Causal Consistency. *EuroSys*.
   - GentleRain 算法

---

## 8. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Causal Consistency Toolkit                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心记忆锚点: "Happens-Before 保持"                                         │
│  ├── 程序序保持                                                              │
│  ├── 读-从关系保持                                                           │
│  └── 传递闭包保持                                                            │
│                                                                              │
│  向量时钟是因果一致性的实现基础:                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ VC(a) < VC(b)  ⟺  a happens-before b                                │    │
│  │ VC(a) || VC(b) ⟺  a and b are concurrent                            │    │
│  │                                                                      │    │
│  │ 实现:                                                                │    │
│  │   - 本地事件: VC[i]++                                                │    │
│  │   - 发送消息: 附带 VC                                                │    │
│  │   - 接收消息: VC = max(VC, received_VC); VC[i]++                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  应用场景:                                                                   │
│  ├── 社交网络: 发帖/评论的因果可见性                                          │
│  ├── 协作编辑: 操作的因果排序                                                 │
│  ├── 消息系统: 消息的因果传递                                                 │
│  └── 电商系统: 下单→支付的因果链                                              │
│                                                                              │
│  常见误区:                                                                   │
│  ❌ "因果一致 = 顺序一致" → 因果一致更弱，不保证程序序                       │
│  ❌ "向量时钟开销大" → 现代系统使用压缩和垃圾回收                             │
│  ❌ "因果一致不需要协调" → 跨数据中心仍需协调                                 │
│  ❌ "所有应用都需要因果一致" → 有时最终一致足够                               │
│                                                                              │
│  与其他模型对比:                                                              │
│  ├── Linearizable: 需要实时序，太强                                          │
│  ├── Sequential: 需要全局程序序，太强                                        │
│  ├── Causal: 只需要因果序，刚刚好                                            │
│  └── Eventual: 无保证，太弱                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*文档生成时间: 2026-04-02*
*维护者: Distributed Systems Knowledge Base*
*版本: S-Level (21+ KB)*
