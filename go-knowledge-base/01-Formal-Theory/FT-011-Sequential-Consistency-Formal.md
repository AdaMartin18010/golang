# FT-011: 顺序一致性的形式化理论 (Sequential Consistency: Formal Theory)

> **维度**: Formal Theory
> **级别**: S (20+ KB)
> **标签**: #sequential-consistency #consistency-models #memory-models #multiprocessors
> **权威来源**:
>
> - [How to Make a Multiprocessor Computer](https://ieeexplore.ieee.org/document/1675439) - Lamport (1979)
> - [A Better x86 Memory Model: x86-TSO](https://www.cl.cam.ac.uk/~pes20/weakmemory/x86tso.pdf) - Sewell et al. (2010)
> - [The Java Memory Model](https://dl.acm.org/doi/10.1145/1040305.1040336) - Manson et al. (2005)
> - [Understanding POWER Multiprocessors](https://dl.acm.org/doi/10.1145/2248487.1950392) - Sarkar et al. (2011)
> - [Modular Relaxed Dependencies](https://arxiv.org/abs/1608.05599) - Alglave et al. (2016)

---

## 1. 顺序一致性的形式化定义

### 1.1 基本模型

**定义 1.1 (共享内存系统)**
共享内存系统 $\mathcal{S}$ 由：

- 进程集合 $\Pi = \{p_1, p_2, ..., p_n\}$
- 共享内存位置集合 $L = \{x_1, x_2, ..., x_m\}$
- 操作集合 $O = \{\text{read}, \text{write}\} \times L \times V$

**定义 1.2 (程序序)**
进程 $p_i$ 的程序序 $<_i$ 是操作在进程内的发生顺序：

$$o_1 <_i o_2 \Leftrightarrow o_1 \text{ 在 } p_i \text{ 中先于 } o_2 \text{ 执行}$$

### 1.2 顺序一致性定义

**定义 1.3 (顺序一致性 - Lamport 1979)**

一个并发执行是顺序一致的，如果：

1. **全局序存在**: 存在一个所有操作的全局顺序 $<$
2. **程序序保持**: 每个进程的操作按程序序出现在全局序中
3. **读值正确**: 每个读操作返回全局序中最近的写操作的值

形式化：

$$\text{SequentialConsistency}(E) \equiv \exists <:$$
$$(\forall p_i \in \Pi: <_i \subseteq <) \land (\forall r \in \text{Reads}: \text{value}(r) = \text{last-write}_<(r))$$

**与线性一致性的区别**:

| 属性 | 顺序一致性 | 线性一致性 |
|------|-----------|-----------|
| **实时序** | 不要求 | 要求 |
| **程序序** | 要求 | 要求 |
| **全局序** | 存在 | 存在并扩展实时序 |
| **强度** | 较弱 | 较强 |

### 1.3 形式化比较

**定理 1.1 (蕴含关系)**

$$\text{Linearizable} \Rightarrow \text{SequentiallyConsistent}$$

*证明*:
线性一致性的全局序扩展实时序，因此必然保持程序序。$\square$

**定理 1.2 (非本地性)**
顺序一致性**不具有本地性**。

*证明*:
考虑两个寄存器 $x$ 和 $y$，每个单独是顺序一致的，但组合后可能不是。

**例子**:

```
P1: write(x, 1) ──► write(y, 1)
P2: read(y)=? ──► read(x)=?
```

如果 $x$ 和 $y$ 分别顺序一致，可能看到：

- $read(y) = 1$, $read(x) = 0$ (如果全局序不同)

这违反了直觉上的因果关系。

$\square$

---

## 2. 内存模型与顺序一致性

### 2.1 处理器内存模型

**x86-TSO (Total Store Order)**
x86 提供接近顺序一致性的内存模型：

- 所有核心看到相同的写顺序
- 读可以重排序到写之前（Store Buffer）
- 使用 MFENCE 保证顺序

**ARM/POWER 弱内存模型**
提供更弱的保证：

- 不同核心可能看到不同的写顺序
- 读和写都可以重排序
- 需要显式屏障指令

### 2.2 顺序一致性的实现

**方法 1: 顺序执行**
简单但低效：

$$\text{禁用所有重排序和优化}$$

**方法 2: 内存屏障**
在关键位置插入屏障：

```
write(x, 1)
MFENCE  // 确保写完成
read(y)
```

**方法 3: 锁同步**
使用互斥锁保护访问：

```
lock(mutex)
write(x, 1)
read(y)
unlock(mutex)
```

### 2.3 编译器重排序

**问题**: 编译器优化可能重排序代码：

```c
// 源代码
x = 1;
y = 2;

// 优化后 (可能)
y = 2;
x = 1;
```

**解决方案**:

- `volatile` 关键字 (C/C++/Java)
- 内存序标记 (C++11 memory_order)
- `synchronized` (Java)

---

## 3. 正确性证明框架

### 3.1 Litmus Tests

**Litmus Test**: 用于区分不同内存模型的小型程序。

**测试 1: Store Buffering (SB)**

```
初始: x = 0, y = 0

P1:          P2:
  x = 1        y = 1
  r1 = y       r2 = x

问题: 是否可能 r1 = 0 且 r2 = 0?
```

**顺序一致性**: 不可能
**x86-TSO**: 可能 (Store Buffer 延迟)
**ARM/POWER**: 可能

**测试 2: Message Passing (MP)**

```
初始: x = 0, flag = 0

P1:          P2:
  x = 1        r1 = flag
  flag = 1     r2 = x

问题: 是否可能 r1 = 1 且 r2 = 0?
```

**顺序一致性**: 不可能
**弱内存模型**: 可能

### 3.2 证明技术

**技术 1: 交换论证**
通过交换相邻的独立操作，将执行转换为顺序执行。

**技术 2: 图分析**
构建程序依赖图，检查环的存在。

**技术 3: 模型检验**
使用工具（如 CDSChecker, herd）自动验证。

---

## 4. 多元表征

### 4.1 顺序一致性概念地图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Sequential Consistency Concept Network                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                    ┌─────────────────────────┐                              │
│                    │  Sequential Consistency │                              │
│                    │    (Lamport 1979)       │                              │
│                    └───────────┬─────────────┘                              │
│                                │                                            │
│              ┌─────────────────┼─────────────────┐                         │
│              ▼                 ▼                 ▼                         │
│       ┌───────────┐    ┌───────────┐    ┌───────────┐                     │
│       │Program    │    │ Global    │    │   Read    │                     │
│       │ Order     │    │  Order    │    │   Value   │                     │
│       │ Preserved │    │  Exists   │    │  Correct  │                     │
│       └─────┬─────┘    └─────┬─────┘    └─────┬─────┘                     │
│             │                │                │                            │
│       ┌─────┴─────┐    ┌─────┴─────┐  ┌─────┴─────┐                       │
│       │           │    │           │  │           │                       │
│       │ Per-      │    │ All ops   │  │ Read sees │                       │
│       │ process   │    │ totally   │  │ last write│                       │
│       │ sequential│    │ ordered   │  │ in order  │                       │
│       └───────────┘    └───────────┘  └───────────┘                       │
│                                                                              │
│  Memory Models Comparison:                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ Sequential Consistency ◄── x86-TSO ◄── ARM/POWER                     │    │
│  │     (Strong)              (Medium)       (Weak)                     │    │
│  │                                                                      │    │
│  │ 所有核心看到相同顺序      写全局有序     允许不同核心看到不同顺序       │    │
│  │ 禁止所有重排序           读可重排序     读写都可重排序                │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  Key Properties:                                                             │
│  ├── Program order preserved within each process                             │
│  ├── No local property (unlike linearizability)                              │
│  ├── Requires hardware/compiler cooperation                                  │
│  └── Strong but expensive to implement                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 内存模型决策树

```
选择内存模型?
│
├── 需要强一致性保证?
│   ├── 是 → Sequential Consistency
│   │       └── 所有操作按程序序执行
│   │
│   └── 否 → 继续
│
├── 目标平台
│   ├── x86/x64?
│   │   └── x86-TSO (接近 SC)
│   │       └── 使用 MFENCE 保证顺序
│   │
│   ├── ARM/POWER?
│   │   └── 弱内存模型
│   │       └── 需要显式屏障指令
│   │
│   └── 跨平台?
│       └── 使用语言级抽象
│           └── C++ memory_order, Java volatile
│
├── 性能要求
│   ├── 最高性能?
│   │   └── 弱内存模型 + 最少屏障
│   │
│   └── 安全优先?
│       └── 顺序一致性 + 保守同步
│
└── 编程语言
    ├── C/C++?
    │   └── memory_order_seq_cst (默认)
    │       └── memory_order_release/acquire (优化)
    │
    ├── Java?
    │   └── volatile / synchronized
    │       └── happens-before 规则
    │
    └── Go?
        └── happens-before (channel, mutex)
```

### 4.3 内存模型对比矩阵

| 属性 | Sequential | x86-TSO | ARM | C++ Relaxed |
|------|------------|---------|-----|-------------|
| **程序序** | ✓ | ✓ | ✓ | ✓ |
| **写-写顺序** | ✓ | ✓ | - | - |
| **读-写顺序** | ✓ | - | - | - |
| **读-读顺序** | ✓ | ✓ | - | - |
| **全局写序** | ✓ | ✓ | - | - |
| **Store Buffer** | - | ✓ | ✓ | ✓ |
| **指令重排序** | - | 读可前移 | 读写都可 | 完全重排 |
| **屏障开销** | 高 | 中 | 高 | 无 |
| **编程难度** | 低 | 中 | 高 | 极高 |

### 4.4 重排序示例图

```
顺序一致性执行:

P1: write(x, 1) ──► write(y, 2)
P2:                 read(x) ──► read(y)

全局序必须保持程序序:
  write(x,1) < write(y,2)   (P1 的程序序)
  read(x) < read(y)          (P2 的程序序)

可能的合法全局序:
  1. write(x,1), write(y,2), read(x), read(y)  → read(x)=1, read(y)=2
  2. write(x,1), read(x), write(y,2), read(y)  → read(x)=1, read(y)=2
  3. read(x), read(y), write(x,1), write(y,2)  → read(x)=0, read(y)=0
  4. read(x), write(x,1), read(y), write(y,2)  → read(x)=0, read(y)=0

非法 (违反程序序):
  write(y,2), write(x,1), ...  (违反 P1 的程序序)
  read(y), read(x), ...        (违反 P2 的程序序)

────────────────────────────────────────────────────────────────────────

x86-TSO 执行 (允许的重排序):

P1: write(x, 1)    (进入 Store Buffer)
    read(y)        (可以先执行!)

Store Buffer 延迟可见性:
  P1 的写对其他核心不是立即可见的

可能的结果:
  P1: write(x,1), read(y)=0
  P2: write(y,1), read(x)=0

  两个读都看到旧值!

解决方案: MFENCE
  P1: write(x,1), MFENCE, read(y)
  MFENCE 排空 Store Buffer，确保写完成
```

---

## 5. TLA+ 形式化规约

```tla
------------------------------- MODULE SequentialConsistency ------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS
    Process,            \* 进程集合
    Location,           \* 内存位置
    Value,              \* 值域
    Nil

VARIABLES
    operations,         \* 操作集合
    programOrder,       \* 程序序关系
    globalOrder,        \* 全局序关系
    memoryState         \* 内存状态

scVars == <<operations, programOrder, globalOrder, memoryState>>

-----------------------------------------------------------------------------
\* 操作定义

Operation == [proc: Process, type: {"read", "write"}, loc: Location, val: Value]

-----------------------------------------------------------------------------
\* 程序序: 同一进程内的操作顺序

ProgramOrder(op1, op2) ==
    /\ op1.proc = op2.proc
    /\ programOrder[<<op1, op2>>]

\* 程序序的传递闭包
ProgramOrderBefore(op1, op2) ==
    ProgramOrder(op1, op2)
    \/ \E op3: ProgramOrder(op1, op3) /\ ProgramOrderBefore(op3, op2)

-----------------------------------------------------------------------------
\* 顺序一致性条件

\* 1. 全局序是全序 (满足传递性、反对称性)
IsTotalOrder ==
    \A op1, op2 \in operations:
        op1 # op2 =>
            (globalOrder[<<op1, op2>>] \/ globalOrder[<<op2, op1>>])

\* 2. 全局序扩展程序序
ExtendsProgramOrder ==
    \A op1, op2 \in operations:
        ProgramOrderBefore(op1, op2) => globalOrder[<<op1, op2>>]

\* 3. 读值正确: 读看到全局序中最近的写
ReadValueCorrect ==
    \A r \in {op \in operations: op.type = "read"}:
        LET writes == {w \in operations:
                        w.type = "write" /\
                        w.loc = r.loc /\
                        globalOrder[<<w, r>>]}
        IN writes = {} => r.val = Nil
           \/ \E w \in writes:
               /\ r.val = w.val
               /\ ~\E w2 \in writes:
                   w2 # w /\ globalOrder[<<w, w2>>] /\ globalOrder[<<w2, r>>]

-----------------------------------------------------------------------------
\* 顺序一致性判定

IsSequentiallyConsistent ==
    /\ IsTotalOrder
    /\ ExtendsProgramOrder
    /\ ReadValueCorrect

=============================================================================
```

---

## 6. Go 实现与示例

### 6.1 顺序一致性测试

```go
package sequential

import (
    "sync"
    "sync/atomic"
    "testing"
)

// TestStoreBuffering 测试 Store Buffering 情况
func TestStoreBuffering(t *testing.T) {
    for i := 0; i < 10000; i++ {
        var x, y int32
        var r1, r2 int32

        var wg sync.WaitGroup
        wg.Add(2)

        go func() {
            defer wg.Done()
            atomic.StoreInt32(&x, 1)
            r1 = atomic.LoadInt32(&y)
        }()

        go func() {
            defer wg.Done()
            atomic.StoreInt32(&y, 1)
            r2 = atomic.LoadInt32(&x)
        }()

        wg.Wait()

        // 在顺序一致性下，不可能 r1=0 且 r2=0
        // 但在 x86-TSO 下是可能的
        if r1 == 0 && r2 == 0 {
            t.Logf("Store buffering observed: r1=%d, r2=%d", r1, r2)
        }
    }
}

// TestMessagePassing 测试消息传递模式
func TestMessagePassing(t *testing.T) {
    for i := 0; i < 10000; i++ {
        var x, flag int32
        var r1, r2 int32

        var wg sync.WaitGroup
        wg.Add(2)

        go func() {
            defer wg.Done()
            atomic.StoreInt32(&x, 1)
            atomic.StoreInt32(&flag, 1)
        }()

        go func() {
            defer wg.Done()
            r1 = atomic.LoadInt32(&flag)
            r2 = atomic.LoadInt32(&x)
        }()

        wg.Wait()

        // 如果看到 flag=1，应该也能看到 x=1
        if r1 == 1 && r2 == 0 {
            t.Errorf("Message passing violation: r1=%d, r2=%d", r1, r2)
        }
    }
}

// SequentiallyConsistentCounter 顺序一致性计数器
type SequentiallyConsistentCounter struct {
    value int64
}

func (c *SequentiallyConsistentCounter) Inc() {
    // atomic.AddInt64 在 x86 上提供顺序一致性
    atomic.AddInt64(&c.value, 1)
}

func (c *SequentiallyConsistentCounter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}

// TestSequentialConsistency 测试顺序一致性
func TestSequentialConsistency(t *testing.T) {
    counter := &SequentiallyConsistentCounter{}

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                counter.Inc()
            }
        }()
    }

    wg.Wait()

    expected := int64(100 * 1000)
    if counter.Get() != expected {
        t.Errorf("Expected %d, got %d", expected, counter.Get())
    }
}
```

### 6.2 内存屏障使用

```go
package sequential

import (
    "runtime"
    "sync/atomic"
)

// MemoryBarrier 内存屏障
type MemoryBarrier struct{}

// FullFence 完全屏障 (x86: MFENCE)
func (m *MemoryBarrier) FullFence() {
    // Go 的 atomic 操作自动包含必要的屏障
    // 在 C/C++ 中，这会对应 __sync_synchronize() 或 std::atomic_thread_fence()
    runtime.Gosched() // 仅作示例，实际屏障由 atomic 操作保证
}

// SafeMessagePassing 安全的 MP 模式
type SafeMessagePassing struct {
    data  int32
    ready int32
}

func (mp *SafeMessagePassing) Write(value int32) {
    atomic.StoreInt32(&mp.data, value)
    // StoreStore 屏障确保 data 的写在 ready 之前可见
    atomic.StoreInt32(&mp.ready, 1)
}

func (mp *SafeMessagePassing) Read() (int32, bool) {
    if atomic.LoadInt32(&mp.ready) == 1 {
        // LoadLoad 屏障确保 ready=1 后能看到 data
        return atomic.LoadInt32(&mp.data), true
    }
    return 0, false
}
```

---

## 7. 学术参考文献

### 7.1 经典论文

1. **Lamport, L. (1979)**. How to Make a Multiprocessor Computer That Correctly Executes Multiprocess Programs. *IEEE Transactions on Computers*, 28(9), 690-691.
   - 顺序一致性的首次形式化定义

2. **Adve, S. V., & Gharachorloo, K. (1996)**. Shared Memory Consistency Models: A Tutorial. *IEEE Computer*, 29(12), 66-76.
   - 内存一致性模型的综合教程

### 7.2 现代研究

1. **Sewell, P., et al. (2010)**. x86-TSO: A Rigorous and Usable Programmer's Model for x86 Multiprocessors. *Communications of the ACM*, 53(7), 89-97.
   - x86 内存模型的形式化

2. **Alglave, J., et al. (2014)**. Herding Cats: Modelling, Simulation, Testing, and Data-mining for Weak Memory. *ACM TOPLAS*, 36(2), 1-74.
   - 弱内存模型的工具和方法

---

## 8. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Sequential Consistency Toolkit                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心记忆锚点: "全局序 + 程序序保持"                                         │
│  ├── 所有操作有一个全局执行顺序                                              │
│  ├── 每个进程的操作按程序序出现在全局序中                                    │
│  └── 读操作返回全局序中最近的写                                              │
│                                                                              │
│  与 Linearizability 的区别:                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ Linearizability: 全局序必须扩展实时序                                 │    │
│  │ Sequential Consistency: 全局序只需保持程序序                          │    │
│  │                                                                      │    │
│  │ 结果: SC 比 Linearizability 更弱，但实现更高效                         │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  实现策略:                                                                   │
│  ├── 禁用所有优化 (简单但慢)                                                 │
│  ├── 内存屏障 (在关键位置插入)                                               │
│  └── 锁同步 (隐式屏障)                                                       │
│                                                                              │
│  常见误区:                                                                   │
│  ❌ "SC = 没有重排序" → SC 允许不同进程的操作重排序                         │
│  ❌ "SC 可以自动保证" → 需要硬件和编译器配合                                │
│  ❌ "volatile = SC" → volatile 只保证可见性，不保证顺序                     │
│  ❌ "SC 总是最好的" → 性能代价大，有时弱一致性更合适                         │
│                                                                              │
│  关键公式:                                                                   │
│  ├── 合法执行 = 存在全局序 < 使得:                                         │
│  │   ∀p: program_order(p) ⊆ <                                            │
│  │   ∀r: value(r) = last_write(r, <)                                      │
│  └── 非本地性: 每个对象是 SC 不保证系统是 SC                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*文档生成时间: 2026-04-02*
*维护者: Distributed Systems Knowledge Base*
*版本: S-Level (20+ KB)*
