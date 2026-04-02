# LD-004: Go 调度器的形式化理论 (Go Scheduler: Formal Theory)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #scheduler #formal-semantics #concurrency #operating-systems #m-n-threading
> **权威来源**:
>
> - [Scheduling Multithreaded Computations by Work Stealing](https://dl.acm.org/doi/10.1145/324133.324234) - Blumofe & Leiserson
> - [Go Scheduler Design](https://go.dev/s/go11sched) - Dmitry Vyukov
> - [The Linux Scheduler](https://www.kernel.org/doc/html/latest/scheduler/) - Linux Kernel
> - [Cilk-5](https://dl.acm.org/doi/10.1145/277651.277685) - Blumofe et al.

---

## 1. 形式化基础

### 1.1 调度理论基础

**定义 1.1 (调度问题)**
给定任务集合 T 和资源集合 R，找到映射 S: T × Time → R 满足约束：

```
∀t ∈ T: S(t) ∈ R
∀t1, t2 ∈ T: S(t1) = S(t2) ⟹ t1 ≠ t2 at same time
```

**定义 1.2 (调度目标)**

```
最小化: makespan = max(Ci)  // 完成时间
最小化: Σ(Ci - Ai)          // 平均响应时间
最大化: throughput = |T| / makespan
最大化: utilization = busy time / total time
```

**定理 1.1 (调度最优性)**
对于抢占式调度，最短剩余时间优先（SRTF）是最优的（最小化平均响应时间）。

*证明*：
考虑任意非 SRTF 调度，找到第一个与 SRTF 不同的调度点。
交换两个任务的执行顺序可以减少总等待时间。
重复此过程可转化为 SRTF 调度且不会增加平均响应时间。

### 1.2 线程模型

**定义 1.3 (用户线程 vs 内核线程)**

```
1:1 模型: 每个用户线程对应一个内核线程
  例: Java, C++, Rust
  优点: 简单，内核调度公平
  缺点: 线程切换开销大，数量有限

M:1 模型: 多个用户线程复用一个内核线程
  例: Python asyncio, Node.js
  优点: 轻量级切换，高并发
  缺点: 无法利用多核，阻塞调用影响所有用户线程

M:N 模型: M 个用户线程映射到 N 个内核线程
  例: Go, Erlang, Haskell
  优点: 轻量级 + 多核支持
  缺点: 复杂的用户态调度器
```

**定理 1.2 (M:N 优势)**
M:N 线程模型可以在用户态完成快速上下文切换，避免内核态开销。

*证明*：

- 用户态上下文切换只需保存/恢复寄存器，约 10-20 条指令
- 内核态切换需要特权级切换，TLB 刷新，约 1000+ 周期
- Go goroutine 切换约 200ns，线程切换约 1-2μs

---

## 2. Go 调度器形式化

### 2.1 GMP 形式定义

**定义 2.1 (Goroutine G)**

```
G = (id, state, stack, fn, context, m, p)

where:
  id ∈ ℕ: unique identifier
  state ∈ {idle, runnable, running, waiting, dead}
  stack = (lo, hi) ∈ Addr × Addr: stack boundaries
  fn: entry function pointer
  context = (pc, sp, bp): saved registers
  m: M* | nil: bound OS thread
  p: P* | nil: bound processor
```

**定义 2.2 (Machine M)**

```
M = (id, g0, curg, p, tls, status)

where:
  id ∈ ℕ: thread identifier
  g0: G - scheduler goroutine
  curg: G* | nil - currently running G
  p: P* | nil - bound P
  tls: [uintptr] - thread-local storage
  status ∈ {idle, running, syscall, spinning}
```

**定义 2.3 (Processor P)**

```
P = (id, status, m, runq, runnext, mcache)

where:
  id ∈ [0, GOMAXPROCS): processor identifier
  status ∈ {idle, running, syscall, gcstop}
  m: M* | nil - bound M
  runq: Queue<G> - local runnable queue, |runq| ≤ 256
  runnext: G* | nil - high priority next G
  mcache: Cache - memory allocator cache
```

### 2.2 状态转换系统

**G 状态转换图**

```
                         create
                    ┌─────────────┐
                    ▼             │
    ┌───────────[idle]        [dead]◄───────┐
    │               │             ▲         │
    │               ▼             │         │ complete
    │         [runnable]──────────┘         │
    │               │   schedule            │
    │    block      ▼                       │
    └──────────[running]────────────────────┘
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
   [waiting]   [syscall]   [copystack]
        │           │           │
        └───────────┴───────────┘
                    │
                    ▼ wakeup/return
              [runnable]
```

**形式化规则**

```
G_idle ──create(f)──► G_runnable(id, f)
G_runnable ──schedule(p)──► G_running(p)
G_running ──block(channel)──► G_waiting(channel)
G_running ──syscall──► G_syscall
G_waiting ──wakeup──► G_runnable
G_syscall ──return──► G_runnable
G_running ──complete──► G_dead
G_running ──preempt──► G_runnable
```

### 2.3 调度不变式

**定理 2.1 (执行不变式)**
在任何时刻：

```
∀m ∈ M: m.curg.state = running ⟹ m.p ≠ nil
∀g ∈ G: g.state = running ⟹ ∃!m: m.curg = g
∀p ∈ P: p.status = running ⟹ ∃!m: p.m = m
∀p ∈ P: |p.runq| = p.runqtail - p.runqhead (mod 256)
```

**定理 2.2 (队列一致性)**

```
GlobalQueue ∪ ⋃(p.runq for p in P) = {g | g.state = runnable}
```

---

## 3. 工作窃取算法形式化

### 3.1 算法定义

**定义 3.1 (工作窃取)**
当一个 P 空闲时，从其他 P 偷取 G 执行。

**算法 3.1 (随机工作窃取)**

```
function stealWork(p_i):
    for attempt in 1..numP:
        j = random(numP)  // 或确定性顺序
        p_j = allp[j]
        if p_j ≠ p_i and |p_j.runq| > 0:
            // 偷取一半的 G
            stolen = floor(|p_j.runq| / 2)
            for k in 1..stolen:
                g = dequeue(p_j.runq)
                enqueue(p_i.runq, g)
            return stolen
    return 0
```

**算法 3.2 (确定性工作窃取)**

```
function stealWorkOrdered(p_i):
    for offset in 1..numP:
        j = (i + offset) mod numP
        p_j = allp[j]
        if |p_j.runq| > 0:
            return stealHalf(p_i, p_j)
    return null
```

**定理 3.1 (工作窃取效率)**
在 P 个处理器上，工作窃取算法的期望窃取次数为 O(P × S)，其中 S 是串行执行时间。

*证明* (Blumofe & Leiserson, 1999):

1. 每个窃取操作至少执行一个任务单位
2. 总工作量为 W，关键路径长度为 S
3. 由 Brent 定理，并行时间为 Tp ≤ W/P + S
4. 窃取次数上限为 O(P × S)

### 3.2 负载均衡

**定义 3.2 (负载均衡)**
令 L_i = |p_i.runq| + (p_i.runnext ≠ nil ? 1 : 0)
系统负载均衡当: max(L_i) / avg(L) ≤ threshold

**定理 3.2 (Go 调度器负载均衡)**
工作窃取保证：

```
∀t: max_i |p_i.runq| ≤ min_i |p_i.runq| + 2
```

*证明*：

1. 当 P_i 负载比 P_j 多超过 2 个 G 时，P_j 会偷取一半
2. 偷取后，P_i 剩下 ⌈n/2⌉，P_j 得到 ⌊n/2⌋
3. 差值最多为 1
4. 考虑 runnext，差值最多为 2

---

## 4. 抢占形式化

### 4.1 协作式抢占

**定义 4.1 (安全点)**
函数调用和循环回边是安全点，可以插入抢占检查。

**抢占检查**

```
function preemptCheck(g):
    if g.preemptFlag and atSafePoint(g.pc):
        saveContext(g)
        g.state = runnable
        enqueue(g.p.runq, g)
        schedule()
```

**定理 4.1 (协作式抢占正确性)**
协作式抢占保证程序状态一致性。

*证明*：

- 抢占只在安全点发生
- 安全点处所有寄存器状态已保存到栈
- GC 可以找到所有根引用
- 状态转换是原子的

### 4.2 信号抢占

**定义 4.2 (异步抢占)**
使用 SIGURG 信号强制目标线程保存上下文。

**信号处理**

```
handler(SIGURG):
    if targetG.canPreempt:
        // 检查在安全点
        if isAtSafePoint(targetG.pc):
            injectCall(asyncPreempt, targetG)
        else:
            // 标记抢占请求
            targetG.preemptStop = true
```

**定理 4.2 (信号抢占安全性)**
异步抢占只在安全点生效，保证程序一致性。

---

## 5. 性能模型

### 5.1 调度延迟

**定义 5.1 (调度延迟)**

```
D = T_start - T_submit
其中 T_submit 是 G 变为 runnable 的时间
      T_start 是 G 开始执行的时间
```

**定理 5.1 (延迟上界)**
在 P 个处理器上，Goroutine 的期望调度延迟为：

```
E[D] ≤ O(T_steal × log P / log log P)
```

*证明*：

1. 最坏情况：G 进入全局队列，所有 P 忙碌
2. 每个 P 完成当前 G 后尝试偷取
3. 偷取全局队列的概率为 1/P
4. 期望等待 P 次偷取尝试
5. 每次偷取 O(log P) 时间

### 5.2 吞吐量模型

**定义 5.2 (吞吐量)**

```
Throughput = CompletedTasks / Time
```

**定理 5.2 (吞吐量下界)**
对于计算密集型任务，Go 调度器达到最优吞吐量的常数因子内：

```
Throughput ≥ c × Optimal, where c ≈ 0.5
```

---

## 6. 多元表征

### 6.1 调度算法对比

| 特性 | 1:1 线程 | M:N 协程 | Go GMP |
|------|----------|----------|--------|
| 上下文切换 | ~1μs | ~100ns | ~200ns |
| 内存开销 | ~1MB | ~4KB | ~2KB |
| 调度策略 | 内核 CFS | 用户决定 | 工作窃取 |
| 负载均衡 | 内核 | 用户 | 分布式偷取 |
| 抢占 | 内核 | 协作 | 强制+协作 |
| 最大并发 | ~10K | ~1M | ~1M |
| 系统调用 | 阻塞 | 非阻塞 | P 移交 |

### 6.2 调度决策图

```
Goroutine 状态变化
│
├── 创建 (go func())
│   ├── 当前 P.runnext 空闲?
│   │   └── 是 → 放入 runnext（立即执行）
│   │
│   ├── 当前 P.runq 未满?
│   │   └── 是 → 放入 runq
│   │
│   ├── 有空闲 P 且有空闲 M?
│   │   └── 是 → 唤醒 M 执行新 G
│   │
│   └── 放入全局队列
│
├── 阻塞 (channel, mutex, sleep)
│   ├── 保存 G 状态
│   ├── G.state = waiting
│   ├── 切换到 g0 栈
│   └── schedule() 找下一个 G
│
├── 唤醒 (channel recv, unlock)
│   ├── 尝试放入 P.runnext
│   ├── 满则放入 P.runq
│   └── 满则放入全局队列
│
├── 系统调用
│   ├── 保存状态
│   ├── M 释放 P
│   ├── P 可被其他 M 接管
│   └── 返回时尝试重获 P
│
└── 完成
    ├── G.state = dead
    ├── 放回 G 空闲列表
    └── schedule() 找下一个 G
```

### 6.3 性能权衡图

```
低延迟 ◄─────────────────────────────────────► 高吞吐
       │              Go Scheduler             │
       │           (平衡设计目标)               │
       │                                        │
       │  asyncio (Node.js)  Thread Pool (Java) │
       │  低吞吐，高并发       高吞吐，高延迟    │
       │                                        │
       └────────────────────────────────────────┘
              低并发                         高并发
```

### 6.4 调度器结构图

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go Scheduler Architecture                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Global State                                                    │
│  ├── allp [GOMAXPROCS]*P     // 所有 P                          │
│  ├── allm *m                  // 所有 M 链表                     │
│  ├── sched {                                                     │
│  │   ├── lock mutex                                                │
│  │   ├── midle *m              // 空闲 M 列表                    │
│  │   ├── nmidle int32          // 空闲 M 数量                    │
│  │   ├── runqhead guintptr     // 全局队列头                     │
│  │   ├── runqtail guintptr     // 全局队列尾                     │
│  │   └── runqsize int32        // 全局队列大小                   │
│  │   }                                                           │
│  └── gcwaiting bool            // GC 等待                        │
│                                                                  │
│  Per-P State                                                     │
│  ├── runq [256]guintptr        // 本地队列                       │
│  ├── runqhead uint32                                           │
│  ├── runqtail uint32                                           │
│  ├── runnext guintptr          // 下一个优先 G                   │
│  └── mcache *mcache            // 内存缓存                       │
│                                                                  │
│  Per-M State                                                     │
│  ├── g0 *g                     // 调度 goroutine                 │
│  ├── curg *g                   // 当前 G                         │
│  ├── p puintptr                // 绑定的 P                       │
│  └── spinning bool             // 正在寻找工作                   │
│                                                                  │
│  Goroutine State                                                 │
│  ├── stack                     // 栈边界                         │
│  ├── sched gobuf               // 保存的寄存器                   │
│  ├── status int32              // 状态                           │
│  ├── m *m                      // 绑定的 M                       │
│  └── p uintptr                 // 绑定的 P                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 7. 代码示例

### 7.1 调度器分析

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/trace"
    "os"
    "sync"
    "time"
)

func main() {
    // 开启执行追踪
    f, _ := os.Create("trace.out")
    defer f.Close()
    trace.Start(f)
    defer trace.Stop()

    // 设置 GOMAXPROCS
    runtime.GOMAXPROCS(4)
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

    // 创建大量 goroutine
    var wg sync.WaitGroup
    start := time.Now()

    for i := 0; i < 100000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // 模拟工作
            sum := 0
            for i := 0; i < 1000; i++ {
                sum += i
            }
            _ = sum
        }()
    }

    wg.Wait()
    fmt.Printf("Time: %v\n", time.Since(start))
}
```

### 7.2 调度控制

```go
package main

import (
    "runtime"
    "time"
)

func main() {
    // 让出时间片
    go func() {
        for {
            // 处理工作...

            // 显式让出
            runtime.Gosched()
        }
    }()

    // 锁定到 OS 线程
    go func() {
        runtime.LockOSThread()
        defer runtime.UnlockOSThread()

        // 这段代码在同一个 OS 线程执行
        // 适用于需要线程本地存储的场景
        // 如某些图形库、实时音频处理
    }()

    time.Sleep(time.Second)
}
```

### 7.3 调度性能测试

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "sync/atomic"
    "time"
)

func main() {
    const n = 1000000

    // 测试 goroutine 创建开销
    start := time.Now()
    var wg sync.WaitGroup
    wg.Add(n)

    for i := 0; i < n; i++ {
        go func() {
            wg.Done()
        }()
    }
    wg.Wait()

    elapsed := time.Since(start)
    fmt.Printf("Created %d goroutines in %v\n", n, elapsed)
    fmt.Printf("Per goroutine: %v\n", elapsed/n)

    // 测试上下文切换
    runtime.GOMAXPROCS(1) // 强制单核

    var counter int64
    done := make(chan bool)

    for i := 0; i < 2; i++ {
        go func() {
            for j := 0; j < 1000000; j++ {
                atomic.AddInt64(&counter, 1)
                runtime.Gosched() // 强制切换
            }
            done <- true
        }()
    }

    <-done
    <-done

    fmt.Printf("Context switches: %d\n", counter)
}
```

---

## 8. 关系网络

```
Go Scheduler Theory
├── Thread Models
│   ├── 1:1 (Kernel threads: Java, C++)
│   ├── M:1 (Green threads: Python asyncio)
│   └── M:N (Hybrid: Go, Erlang)
├── Scheduling Algorithms
│   ├── Round Robin
│   ├── Priority Queue
│   ├── CFS (Linux)
│   └── Work Stealing (Go, Cilk)
├── Preemption
│   ├── Cooperative (Python, Ruby)
│   ├── Time-slice (Java)
│   └── Signal-based (Go 1.14+)
└── Synchronization
    ├── Spin Lock
    ├── Mutex
    └── Futex
```

---

## 9. 参考文献

1. Blumofe, R. D. & Leiserson, C. E. Scheduling Multithreaded Computations by Work Stealing.
2. Vyukov, D. Go Scheduler Design Doc.
3. Brent, R. P. The Parallel Evaluation of General Arithmetic Expressions.

---

## 10. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────┐
│                  Go Scheduler Toolkit                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  核心概念                                                        │
│  ═══════════════════════════════════════════════════════════    │
│  • G (Goroutine): 用户级线程，轻量级，~2KB 栈                    │
│  • M (Machine): OS 线程，执行 G 的载体                           │
│  • P (Processor): 逻辑处理器，调度和资源管理                       │
│                                                                  │
│  调度原则                                                        │
│  • M 必须绑定 P 才能执行 G                                       │
│  • P 的本地队列优先于全局队列                                    │
│  • 空闲 P 从其他 P 偷取工作                                      │
│  • 系统调用时 M 释放 P                                          │
│                                                                  │
│  性能提示                                                        │
│  □ 避免阻塞操作（使用 channel 而非 mutex）                       │
│  □ 适当设置 GOMAXPROCS                                          │
│  □ 大量 goroutine 考虑 worker pool                               │
│  □ 使用 sync.Pool 减少分配                                       │
│                                                                  │
│  调试工具                                                        │
│  • GODEBUG=schedtrace=X                                         │
│  • runtime/trace                                                 │
│  • go tool trace                                                 │
│                                                                  │
│  常见反模式                                                      │
│  ❌ 创建过多 goroutine 不控制并发度                              │
│  ❌ 在热路径使用 LockOSThread                                    │
│  ❌ 长时间运行的 goroutine 不退出                                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16KB)
**完成日期**: 2026-04-02
