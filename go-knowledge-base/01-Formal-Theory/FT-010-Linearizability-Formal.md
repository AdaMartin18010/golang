# FT-010: 线性一致性的形式化理论 (Linearizability: Formal Theory)

> **维度**: Formal Theory
> **级别**: S (21+ KB)
> **标签**: #linearizability #consistency-models #formal-methods #concurrent-programming
> **权威来源**:
>
> - [Linearizability: A Correctness Condition for Concurrent Objects](https://dl.acm.org/doi/10.1145/78969.78972) - Herlihy & Wing (1990)
> - [How to Make a Multiprocessor Computer That Correctly Executes Multiprocess Programs](https://ieeexplore.ieee.org/document/1675439) - Lamport (1979)
> - [On the Correctness of Database Systems Under Weak Consistency](https://dl.acm.org/doi/10.1145/3035918.3064037) - Cerone & Gotsman (2018)
> - [Principles of Eventual Consistency](https://www.microsoft.com/en-us/research/publication/principles-of-eventual-consistency/) - Burckhardt (2014)
> - [Consistency in Non-Transactional Distributed Storage Systems](https://dl.acm.org/doi/10.1145/2926965) - Viotti & Vukolić (2016)

---

## 1. 线性一致性的形式化定义

### 1.1 基本模型

**定义 1.1 (并发系统)**
一个并发系统 $\mathcal{C}$ 由进程集合 $\Pi = \{p_1, ..., p_n\}$ 和共享对象集合 $\mathcal{O}$ 组成。

**定义 1.2 (操作)**
每个操作 $op$ 是一个二元组：

$$op = \langle \text{invocation}, \text{response} \rangle$$

其中：

- $\text{invocation} = (o, \text{method}, \text{args})$ 在时刻 $t_{\text{inv}}$ 发生
- $\text{response} = (\text{result})$ 在时刻 $t_{\text{res}}$ 发生

**定义 1.3 (实时顺序)**

$$op_1 <_{\text{realtime}} op_2 \Leftrightarrow t_{\text{res}}(op_1) < t_{\text{inv}}(op_2)$$

### 1.2 线性一致性定义

**定义 1.4 (线性一致性 - Herlihy & Wing 1990)**

一个并发执行 $E$ 是线性化的（linearizable），如果存在：

1. **全序 $<$**: 扩展了实时顺序 $<_\text{realtime}$
2. **线性化点**: 每个操作在调用和返回之间的某个瞬间原子执行
3. **顺序正确性**: 按全序 $<$ 执行的结果与串行执行一致

形式化：

$$\text{Linearizable}(E) \equiv \exists <: \forall op_1, op_2:$$
$$(op_1 <_\text{realtime} op_2 \Rightarrow op_1 < op_2) \land \text{SequentialCorrectness}(<)$$

**定义 1.5 (线性化点)**

对于每个操作 $op$，存在线性化点 $t_{\text{lin}}$ 满足：

$$t_{\text{inv}}(op) \leq t_{\text{lin}}(op) \leq t_{\text{res}}(op)$$

在 $t_{\text{lin}}$，操作看起来瞬间完成。

### 1.3 历史与规范

**定义 1.6 (历史)**
历史 $H$ 是操作的集合加上偏序关系：

$$H = (E, <_H)$$

其中 $E$ 是事件集合，$<_H$ 是进程序和实时序的传递闭包。

**定义 1.7 (顺序规范)**
对象 $O$ 的顺序规范 $S_O$ 是允许的串行历史集合。

**定理 1.1 (线性化判定)**
历史 $H$ 是线性化的当且仅当：

$$\exists H' \in \text{Sequential}: H' \text{ extends } H \land H' \in S_O$$

---

## 2. 线性一致性与并发

### 2.1 与串行化的关系

**定义 2.1 (可串行化)**
一个交错执行是可串行化的，如果它等价于某个串行执行：

$$\exists \text{串行调度 } S: \text{执行结果}(E) = \text{执行结果}(S)$$

**定理 2.1 (线性化 vs 串行化)**

$$\text{Linearizable} \subset \text{Serializable}$$

线性一致性是更强的条件，要求：

1. 实时顺序的尊重
2. 每个操作的原子性点

**区别**:

| 属性 | 可串行化 | 线性一致性 |
|------|----------|-----------|
| **顺序约束** | 无 | 尊重实时序 |
| **事务** | 多操作 | 单操作 |
| **实时** | 不保证 | 保证 |
| **应用场景** | 数据库事务 | 并发对象 |

### 2.2 本地性与组合性

**定理 2.2 (本地性 - Locality Property)**
历史 $H$ 是线性化的当且仅当每个对象的历史 $H|_O$ 是线性化的。

*证明概要*:

$(\Rightarrow)$: 显然，全局线性化点的限制是对象的线性化点。

$(\Leftarrow)$:

- 对每个对象 $O_i$，存在线性化 $<_{O_i}$
- 构造全局序：如果 $op_1$ 在 $O_i$ 中先于 $op_2$，或 $op_1 <_\text{realtime} op_2$，则 $op_1 < op_2$
- 由对象线性化的存在性，此序是良定义的

$\square$

**重要性**: 本地性允许独立验证每个对象的正确性。

### 2.3 非阻塞性

**定理 2.3 (非阻塞性 - Non-blocking Property)**
线性一致性不阻塞：

$$\text{pending}(op) \not\Rightarrow \text{block}(\text{complete}(op))$$

即未完成的操作不会阻止其他操作的完成。

---

## 3. 线性一致性的形式化证明

### 3.1 证明方法

**方法 1: 线性化点构造**
为每个操作分配线性化点，验证：

1. 点在调用和返回之间
2. 按线性化点排序的执行满足规范

**方法 2: 交换论证**
通过交换相邻的独立操作，将并发历史转换为串行历史。

**方法 3: 模拟关系**
证明实现与规范之间的模拟关系。

### 3.2 计数器示例

**规范**: 计数器支持 $inc()$ 和 $get()$ 操作。

**定理 3.1 (计数器线性化)**
使用 CAS (Compare-And-Swap) 实现的计数器是线性化的。

*证明*:

实现：

```go
func (c *Counter) Inc() {
    for {
        old := c.value.Load()
        new := old + 1
        if c.value.CompareAndSwap(old, new) {
            return
        }
    }
}
```

线性化点：CAS 成功的瞬间。

验证：

1. 线性化点在调用（读取 old）和返回（CAS 成功）之间
2. 每个成功的 CAS 对应一个原子增量
3. 按 CAS 成功时间排序，执行等价于串行计数器

$\square$

### 3.3 队列示例

**规范**: FIFO 队列支持 $enq(v)$ 和 $deq()$ 操作。

**定理 3.2 (Michael-Scott 队列)**
Michael-Scott 无锁队列是线性化的。

*证明概要*:

**Enq 线性化点**: CAS 成功修改 $tail.next$ 的瞬间。

**Deq 线性化点**: CAS 成功修改 $head$ 的瞬间。

**验证**:

- Enq 线性化点确定了元素的入队顺序
- Deq 线性化点确定了元素的出队顺序
- 队列的链接结构保证 FIFO
- 无元素丢失或重复

$\square$

---

## 4. 多元表征

### 4.1 线性一致性概念地图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Linearizability Concept Network                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                    ┌─────────────────────────┐                              │
│                    │    Linearizability      │                              │
│                    │  (Strongest Consistency)│                              │
│                    └───────────┬─────────────┘                              │
│                                │                                            │
│              ┌─────────────────┼─────────────────┐                         │
│              ▼                 ▼                 ▼                         │
│       ┌───────────┐    ┌───────────┐    ┌───────────┐                     │
│       │ Real-Time │    │ Atomic    │    │ Sequential│                     │
│       │  Order    │    │  Point    │    │ Spec      │                     │
│       └─────┬─────┘    └─────┬─────┘    └─────┬─────┘                     │
│             │                │                │                            │
│       ┌─────┴─────┐    ┌─────┴─────┐  ┌─────┴─────┐                       │
│       │           │    │           │  │           │                       │
│       │ op1 ends  │    │ Invocation│  │ Every op  │                       │
│       │ before    │    │ ───► Lin  │  │ matches   │                       │
│       │ op2 starts│    │ Point ───►│  │ serial    │                       │
│       │ => op1 <  │    │ Response  │  │ execution │                       │
│       │    op2    │    │           │  │           │                       │
│       └───────────┘    └───────────┘  └───────────┘                       │
│                                                                              │
│  Key Properties:                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ Locality: System is linearizable iff each object is linearizable   │    │
│  │ Non-blocking: Pending ops don't block other ops                     │    │
│  │ Composability: Linearizable objects compose into linearizable sys   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  Consistency Hierarchy:                                                      │
│                                                                              │
│       Linearizable ────────────────────────────────────────────────►      │
│            │                                                                │
│            ▼                                                                │
│       Sequential ──────────────────────────────────────────────────►      │
│            │                                                                │
│            ▼                                                                │
│       Causal ──────────────────────────────────────────────────────►      │
│            │                                                                │
│            ▼                                                                │
│       Eventual ────────────────────────────────────────────────────►      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 一致性模型决策树

```
选择一致性模型?
│
├── 需要实时顺序保证?
│   ├── 是 → Linearizable
│   │       └── 所有操作看起来在调用和返回之间瞬间完成
│   │
│   └── 否 → 继续
│
├── 需要程序顺序保证?
│   ├── 是 → Sequential
│   │       └── 每个进程的操作按程序序执行
│   │
│   └── 否 → 继续
│
├── 需要因果顺序保证?
│   ├── 是 → Causal
│   │       └── 因果相关的操作顺序一致
│   │
│   └── 否 → 继续
│
├── 需要单调读保证?
│   ├── 是 → PRAM / Read-Your-Writes
│   │
│   └── 否 → 继续
│
└── 最终一致足够?
    └── 是 → Eventual
        └── 如果停止更新，最终所有副本一致
```

### 4.3 一致性模型对比矩阵

| 属性 | Linearizable | Sequential | Causal | PRAM | Eventual |
|------|--------------|------------|--------|------|----------|
| **实时序** | ✓ | - | - | - | - |
| **程序序** | ✓ | ✓ | - | ✓ | - |
| **因果序** | ✓ | ✓ | ✓ | - | - |
| **收敛** | ✓ | ✓ | ✓ | ✓ | ✓ |
| **本地性** | ✓ | - | ✓ | ✓ | ✓ |
| **非阻塞** | ✓ | ✓ | ✓ | ✓ | ✓ |
| **典型延迟** | 高 | 中 | 低 | 低 | 最低 |
| **协调开销** | 高 | 中 | 低 | 低 | 无 |
| **实现示例** | Java ConcurrentHashMap | 数据库串行化 | COPS | - | DNS |

### 4.4 线性化点示例图

```
时间 →

进程 P1:    [======= op1 =======]
                    ▲
                    │ 线性化点 (Lin Point)
                    ▼
进程 P2:        [======= op2 =======]
                        ▲
                        │ 线性化点
                        ▼
进程 P3:            [======= op3 =======]
                            ▲
                            │ 线性化点
                            ▼

实时序: op1 < op2 < op3

线性化执行序列:
    op1 ──► op2 ──► op3

注意: 即使 op2 和 op3 在时间上重叠，
      线性一致性要求它们有一个全序关系

示例: 计数器 (初始值 = 0)

P1: increment()
P2: increment()
P3: get()

重叠执行:
    P1: [======== inc() ========]
             ▲ Lin Point (value=1)
    P2:     [======== inc() ========]
                 ▲ Lin Point (value=2)
    P3:             [==== get() ====]
                         ▲ Lin Point (returns 2)

可能的线性化结果:
- 如果 P3 的 get 在 P1 和 P2 之后: 返回 2 ✓
- 如果 P3 的 get 在 P1 之后、P2 之前: 返回 1 ✓
- 如果 P3 的 get 在两个之前: 返回 0 ✓

非线性化结果 (不可能):
- 返回 3 (没有 3 个 increment)
- P1 先执行但 P2 的 inc 先完成，get 返回 1 但 P2 的 inc 已在 P3 开始之前完成
```

---

## 5. TLA+ 形式化规约

```tla
------------------------------- MODULE Linearizability ------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS
    Process,            \* 进程集合
    Value,              \* 值域
    Object,             \* 对象集合
    Nil

VARIABLES
    operations,         \* 操作集合
    invocationTime,     \* 调用时间
    responseTime,       \* 返回时间
    linearizationPoint, \* 线性化点
    state               \* 对象状态

linVars == <<operations, invocationTime, responseTime, linearizationPoint, state>>

-----------------------------------------------------------------------------
\* 操作定义

Operation == [proc: Process, obj: Object, op: {"read", "write"},
              val: Value, status: {"pending", "complete"}]

-----------------------------------------------------------------------------
\* 实时顺序

RealTimeOrder(op1, op2) ==
    responseTime[op1] < invocationTime[op2]

\* 实时顺序的传递闭包
RealTimeBefore(op1, op2) ==
    RealTimeOrder(op1, op2)
    \/ \E op3: RealTimeOrder(op1, op3) /\ RealTimeBefore(op3, op2)

-----------------------------------------------------------------------------
\* 线性化条件

\* 1. 线性化点在调用和返回之间
ValidLinearizationPoint(op) ==
    /\ linearizationPoint[op] >= invocationTime[op]
    /\ linearizationPoint[op] <= responseTime[op]

\* 2. 线性化点定义全序
LinearizationOrder(op1, op2) ==
    linearizationPoint[op1] < linearizationPoint[op2]

\* 3. 线性化序扩展实时序
ExtendsRealTime ==
    \A op1, op2 \in operations:
        RealTimeBefore(op1, op2) => LinearizationOrder(op1, op2)

-----------------------------------------------------------------------------
\* 顺序规范 (寄存器)

\* 读操作返回最近的写
RegisterSequentialSpec ==
    \A r \in {op \in operations: op.op = "read"}:
        LET writes == {w \in operations: w.op = "write" /\
                       LinearizationOrder(w, r)}
        IN writes = {} => r.val = Nil
           \/ \E w \in writes:
               /\ r.val = w.val
               /\ ~\E w2 \in writes:
                   w2 # w /\ LinearizationOrder(w, w2) /\ LinearizationOrder(w2, r)

-----------------------------------------------------------------------------
\* 线性化判定

IsLinearizable ==
    /\ \A op \in operations: ValidLinearizationPoint(op)
    /\ ExtendsRealTime
    /\ RegisterSequentialSpec

=============================================================================
```

---

## 6. Go 实现

### 6.1 线性化测试框架

```go
package linearizability

import (
    "fmt"
    "sort"
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

// Event 表示操作事件
type Event struct {
    Process   int
    Op        string // "call" or "return"
    Name      string // operation name
    Value     interface{}
    Timestamp int64
}

// History 操作历史
type History []Event

// LinearizationChecker 线性化检查器
type LinearizationChecker struct {
    Model Model
}

// Model 定义对象模型
type Model struct {
    Init func() interface{}
    Step func(state interface{}, input interface{}, output interface{}) (bool, interface{})
}

// Check 检查历史是否线性化
func (c *LinearizationChecker) Check(history History) bool {
    // 提取操作
    ops := c.extractOperations(history)

    // 尝试找到线性化点
    return c.checkLinearization(ops, c.Model.Init())
}

// extractOperations 从事件中提取操作
func (c *LinearizationChecker) extractOperations(history History) []Operation {
    var ops []Operation
    calls := make(map[int]Event)

    for _, e := range history {
        if e.Op == "call" {
            calls[e.Process] = e
        } else {
            if call, ok := calls[e.Process]; ok {
                ops = append(ops, Operation{
                    Process: call.Process,
                    Name:    call.Name,
                    Input:   call.Value,
                    Output:  e.Value,
                    Start:   call.Timestamp,
                    End:     e.Timestamp,
                })
                delete(calls, e.Process)
            }
        }
    }

    return ops
}

// Operation 表示完整操作
type Operation struct {
    Process int
    Name    string
    Input   interface{}
    Output  interface{}
    Start   int64
    End     int64
}

// checkLinearization 递归检查线性化
func (c *LinearizationChecker) checkLinearization(ops []Operation, state interface{}) bool {
    if len(ops) == 0 {
        return true
    }

    // 尝试每个可能的下一个操作
    for i, op := range ops {
        ok, newState := c.Model.Step(state, op.Input, op.Output)
        if ok {
            // 移除第 i 个操作
            remaining := append(ops[:i], ops[i+1:]...)
            if c.checkLinearization(remaining, newState) {
                return true
            }
        }
    }

    return false
}

// 计数器模型
var CounterModel = Model{
    Init: func() interface{} { return 0 },
    Step: func(state interface{}, input interface{}, output interface{}) (bool, interface{}) {
        counter := state.(int)

        switch input.(string) {
        case "inc":
            return true, counter + 1
        case "get":
            if output.(int) == counter {
                return true, counter
            }
            return false, counter
        default:
            return false, counter
        }
    },
}

// 寄存器模型
var RegisterModel = Model{
    Init: func() interface{} { return nil },
    Step: func(state interface{}, input interface{}, output interface{}) (bool, interface{}) {
        reg := state

        switch inp := input.(type) {
        case WriteOp:
            return true, inp.Value
        case ReadOp:
            if output == reg {
                return true, reg
            }
            return false, reg
        default:
            return false, reg
        }
    },
}

type WriteOp struct{ Value interface{} }
type ReadOp struct{}

// ConcurrentCounter 并发计数器实现
type ConcurrentCounter struct {
    value int64
}

func (c *ConcurrentCounter) Inc() {
    atomic.AddInt64(&c.value, 1)
}

func (c *ConcurrentCounter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}

// 线性化测试
func TestCounterLinearizability(t *testing.T) {
    checker := &LinearizationChecker{Model: CounterModel}

    // 执行并发操作并记录历史
    var history History
    var mu sync.Mutex
    var wg sync.WaitGroup

    counter := &ConcurrentCounter{}

    // 多个 goroutine 并发操作
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            // Inc
            start := time.Now().UnixNano()
            counter.Inc()
            end := time.Now().UnixNano()

            mu.Lock()
            history = append(history,
                Event{Process: id, Op: "call", Name: "inc", Timestamp: start},
                Event{Process: id, Op: "return", Name: "inc", Timestamp: end},
            )
            mu.Unlock()
        }(i)
    }

    wg.Wait()

    // 检查线性化
    if !checker.Check(history) {
        t.Error("Counter is not linearizable")
    }
}
```

### 6.2 线性化数据结构实现

```go
package linearizability

import (
    "sync"
    "sync/atomic"
)

// LinearizableQueue 线性化队列接口
type LinearizableQueue interface {
    Enqueue(v interface{})
    Dequeue() (interface{}, bool)
}

// LockFreeQueue Michael-Scott 无锁队列
type LockFreeQueue struct {
    head unsafe.Pointer
    tail unsafe.Pointer
}

type node struct {
    value interface{}
    next  unsafe.Pointer
}

func NewLockFreeQueue() *LockFreeQueue {
    dummy := &node{}
    return &LockFreeQueue{
        head: unsafe.Pointer(dummy),
        tail: unsafe.Pointer(dummy),
    }
}

func (q *LockFreeQueue) Enqueue(v interface{}) {
    newNode := &node{value: v}

    for {
        tail := (*node)(atomic.LoadPointer(&q.tail))
        next := (*node)(atomic.LoadPointer(&tail.next))

        // 检查 tail 是否仍然是最新的
        if tail == (*node)(atomic.LoadPointer(&q.tail)) {
            if next == nil {
                // 尝试链接新节点
                if atomic.CompareAndSwapPointer(&tail.next, unsafe.Pointer(next), unsafe.Pointer(newNode)) {
                    // 尝试更新 tail
                    atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(newNode))
                    return
                }
            } else {
                // 帮助更新 tail
                atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
            }
        }
    }
}

func (q *LockFreeQueue) Dequeue() (interface{}, bool) {
    for {
        head := (*node)(atomic.LoadPointer(&q.head))
        tail := (*node)(atomic.LoadPointer(&q.tail))
        next := (*node)(atomic.LoadPointer(&head.next))

        // 检查 head 是否仍然是最新的
        if head == (*node)(atomic.LoadPointer(&q.head)) {
            if head == tail {
                // 队列可能为空
                if next == nil {
                    return nil, false
                }
                // 帮助更新 tail
                atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
            } else {
                value := next.value
                // 尝试更新 head
                if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
                    return value, true
                }
            }
        }
    }
}
```

---

## 7. 学术参考文献

### 7.1 经典论文

1. **Herlihy, M. P., & Wing, J. M. (1990)**. Linearizability: A Correctness Condition for Concurrent Objects. *ACM Transactions on Programming Languages and Systems*, 12(3), 463-492.
   - 线性一致性的奠基性工作

2. **Lamport, L. (1979)**. How to Make a Multiprocessor Computer That Correctly Executes Multiprocess Programs. *IEEE Transactions on Computers*, C-28(9), 690-691.
   - 顺序一致性的早期形式化

### 7.2 现代研究

1. **Burckhardt, S. (2014)**. Principles of Eventual Consistency. *Foundations and Trends in Programming Languages*, 1(1-2), 1-150.
   - 最终一致性的综合理论

2. **Viotti, P., & Vukolić, M. (2016)**. Consistency in Non-Transactional Distributed Storage Systems. *ACM Computing Surveys*, 49(1), 1-34.
   - 一致性模型的综述

3. **Cerone, A., & Gotsman, A. (2018)**. Analysing Snapshot Isolation. *Journal of the ACM*, 65(2), 1-41.
   - 快照隔离的形式化分析

---

## 8. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Linearizability Toolkit                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心记忆锚点: "每个操作有一个瞬间完成的点"                                   │
│  ├── 线性化点在调用和返回之间                                                │
│  ├── 按线性化点排序的执行是串行的                                            │
│  └── 串行执行满足对象规范                                                    │
│                                                                              │
│  验证方法:                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ 1. 为每个操作分配线性化点                                             │    │
│  │ 2. 确保点在调用和返回之间                                             │    │
│  │ 3. 验证按线性化点排序的执行满足顺序规范                               │    │
│  │ 4. 验证线性化序扩展实时序                                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  关键性质:                                                                   │
│  ├── 本地性: 系统线性化 ⟺ 每个对象线性化                                   │
│  ├── 非阻塞性: 未完成操作不阻塞其他操作                                     │
│  └── 组合性: 线性化对象组合成线性化系统                                     │
│                                                                              │
│  常见误区:                                                                   │
│  ❌ "Linearizability = Serializability" → 线性化是更强的条件               │
│  ❌ "线性化太慢" → 无锁算法可以实现线性化                                   │
│  ❌ "线性化只适用于单对象" → 本地性允许组合                                 │
│  ❌ "所有并发结构都应线性化" → 有时较弱一致性更合适                         │
│                                                                              │
│  实现技巧:                                                                   │
│  ├── 使用原子操作 (CAS, FAA)                                              │
│  ├── 线性化点通常是 CAS 成功的瞬间                                          │
│  ├── 使用互斥锁 (粗粒度，但正确)                                            │
│  └── 参考无锁数据结构 (Michael-Scott queue, Treiber stack)                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Learning Resources

### Academic Papers

1. **Herlihy, M. P., & Wing, J. M.** (1990). Linearizability: A Correctness Condition for Concurrent Objects. *ACM Transactions on Programming Languages and Systems*, 12(3), 463-492. DOI: [10.1145/78969.78972](https://doi.org/10.1145/78969.78972)
2. **Herlihy, M.** (1991). Wait-Free Synchronization. *ACM TOPLAS*, 13(1), 124-149. DOI: [10.1145/114005.102808](https://doi.org/10.1145/114005.102808)
3. **Herlihy, M., & Shavit, N.** (2011). *The Art of Multiprocessor Programming* (Revised). Morgan Kaufmann.
4. **Dongol, B., & Groves, L.** (2015). Contextual Refinement of the Michael-Scott Queue. *Formal Aspects of Computing*, 27(4), 651-669. DOI: [10.1007/s00165-015-0335-3](https://doi.org/10.1007/s00165-015-0335-3)

### Video Tutorials

1. **Maurice Herlihy.** (2013). [The Art of Multiprocessor Programming](https://www.youtube.com/watch?v=KyFVA4Spcgg). Microsoft Research.
2. **MIT 6.852.** (2014). [Distributed Algorithms](https://www.youtube.com/watch?v=VpmSDT1mY5E). Lecture: Linearizability.
3. **Kavya Joshi.** (2018). [Understanding Lock-Free Algorithms](https://www.youtube.com/watch?v=UOgBFL_S6R8). Strange Loop.
4. **Martin Kleppmann.** (2019). [Isolation Levels and Linearizability](https://www.youtube.com/watchv=plA1GAnFtLA). QCon.

### Book References

1. **Herlihy, M., & Shavit, N.** (2012). *The Art of Multiprocessor Programming* (Chapters 3, 9, 16). Morgan Kaufmann.
2. **Scott, M. L.** (2013). *Shared-Memory Synchronization* (Chapters 4-5). Morgan & Claypool.
3. **Lynch, N. A.** (1996). *Distributed Algorithms* (Chapter 13). Morgan Kaufmann.
4. **Taubenfeld, G.** (2006). *Synchronization Algorithms and Concurrent Programming*. Pearson.

### Online Courses

1. **MIT 6.852.** [Distributed Algorithms](https://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-852j-distributed-algorithms-fall-2009/) - Full course.
2. **Coursera.** [Parallel Programming](https://www.coursera.org/learn/parprog1) - EPFL.
3. **Udacity.** [High Performance Computing](https://www.udacity.com/course/high-performance-computing--ud281) - Concurrency.
4. **edX.** [Hardware/Software Interface](https://www.edx.org/course/computer-systems) - Synchronization.

### GitHub Repositories

1. [mpherlihy/TAOMP](https://github.com/) - Code from "The Art of Multiprocessor Programming".
2. [kavjeydev/lockfree](https://github.com/) - Lock-free data structures examples.
3. [golang/go](https://github.com/golang/go/tree/master/src/sync/atomic) - Go atomic operations.
4. [uber-go/atomic](https://github.com/uber-go/atomic) - Enhanced atomic operations.

### Conference Talks

1. **Maurice Herlihy.** (2012). *Linearizability and Beyond*. PODC Keynote.
2. **Nir Shavit.** (2011). *Data Structures in the Multicore Age*. ESA.
3. **Tim Harris.** (2010). *Transactional Memory*. Microsoft Research.
4. **Kavya Joshi.** (2017). *Understanding the Linux Kernel CPU Scheduler*. GOTO Chicago.

---

*文档生成时间: 2026-04-02*
*维护者: Distributed Systems Knowledge Base*
*版本: S-Level (21+ KB)*
