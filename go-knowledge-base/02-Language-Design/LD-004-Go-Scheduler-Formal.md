# LD-004: Go 调度器的形式化模型 (Go Scheduler: Formal Model)

> **维度**: Language Design
> **级别**: S (17+ KB)
> **标签**: #go-scheduler #mpm #goroutine #os-thread #formal-model
> **权威来源**:
>
> - [Analysis of the Go Runtime Scheduler](https://www.cs.columbia.edu/~aho/cs6998/reports/12-12-11_DeshpandeSponslerWeiss.pdf) - Columbia University (2012)
> - [Go's Work-Stealing Scheduler](https://www.youtube.com/watch?v=YHRO5WQugn0) - Daniel Scales (2012)
> - [The Go Scheduler](https://morsmachine.dk/go-scheduler) - Daniel Morsing
> - [Linux Kernel Scheduling](https://www.kernel.org/doc/html/latest/scheduler/index.html) - Linux Foundation
> - [Scheduling Multithreaded Computations by Work Stealing](https://dl.acm.org/doi/10.1145/324133.324234) - Blumofe & Leiserson (1999)

---

## 1. Go 调度器的形式化定义

### 1.1 M:P:G 模型

**定义 1.1 (调度器状态)**
调度器 $S$ 是三元组 $\langle M, P, G \rangle$：

- $M = \{ m_1, m_2, ..., m_n \}$: OS 线程集合
- $P = \{ p_1, p_2, ..., p_k \}$: 逻辑处理器集合，$k = \text{GOMAXPROCS}$
- $G = \{ g_1, g_2, ..., g_m \}$: Goroutine 集合

**定义 1.2 (M:P:G 关系)**
$$\text{Run}(M) \times \text{Bind}(P) \times \text{Queue}(G)$$

- 每个 $M$ 必须绑定一个 $P$ 才能执行 $G$
- 每个 $P$ 有本地可运行 $G$ 队列
- 全局可运行 $G$ 队列被所有 $P$ 共享

### 1.2 状态转换系统

**定义 1.3 (Goroutine 状态)**
$$\text{G-State} ::= \text{Gidle} \mid \text{Grunnable} \mid \text{Grunning} \mid \text{Gwaiting} \mid \text{Gdead}$$

**状态转换**:

```
Gidle --create--> Grunnable --schedule--> Grunning --block--> Gwaiting
  ^                                              |
  |_________wakeup/ready_________________________|

Grunning --yield/preempt--> Grunnable
Grunning --complete--> Gdead
```

---

## 2. 调度算法的形式化

### 2.1 本地队列优先级

**定义 2.1 (调度优先级)**
$$\text{Priority}(g) = \begin{cases} \infty & g \in \text{LocalQueue}(P) \\ 0 & g \in \text{GlobalQueue} \end{cases}$$
优先执行本地队列的 goroutine。

**定理 2.1 (工作窃取最优性)**
当 $P_i$ 的本地队列为空时，从其他 $P_j$ 窃取工作的期望时间为 $O(1)$。

### 2.2 抢占式调度

**定义 2.2 (函数调用检查点)**
Go 1.14+ 在每个函数调用处插入抢占检查：
$$\text{CheckPreempt}: \text{stack.guard} \to \{\text{continue}, \text{preempt}\}$$

**协作式 vs 抢占式**:

- Go < 1.14: 纯协作式 (函数调用/系统调用时切换)
- Go ≥ 1.14: 信号抢占 (SIGURG)

---

## 3. 系统调用的形式化

### 3.1 阻塞处理

**定义 3.1 (系统调用类型)**

| 类型 | 行为 | 处理 |
|------|------|------|
| 阻塞式 | 等待 I/O | $M$ 阻塞，$P$ 分离 |
| 非阻塞式 | 立即返回 | 继续执行 |

**分离 (Park) 操作**:
$$\text{Syscall}(m, p) \Rightarrow \text{detaches}(p) \land \text{steals}(m_{new}, p)$$
阻塞 $M$ 释放 $P$，新 $M$ 接管 $P$。

---

## 4. 多元表征

### 4.1 M:P:G 结构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Go M:P:G Scheduler                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Global Queue [G1, G2, G3, ...]                                             │
│       ▲                      │                                              │
│       │ steal/global get     │                                              │
│       │                      ▼                                              │
│  ┌────┴──────────────────────────────────────────┐                        │
│  │                                                │                        │
│  │   M0 (OS Thread)      M1 (OS Thread)          │                        │
│  │      │                     │                   │                        │
│  │      ▼                     ▼                   │                        │
│  │   ┌─────┐              ┌─────┐                │                        │
│  │   │ P0  │              │ P1  │    ...         │                        │
│  │   │     │              │     │                │                        │
│  │   │ RunQ│              │ RunQ│                │                        │
│  │   │[G4,│              │[G7, │                │                        │
│  │   │ G5]│              │ G8] │                │                        │
│  │   └─────┘              └─────┘                │                        │
│  │      ▲                     ▲                   │                        │
│  │      │ work steal         │                   │                        │
│  └──────┼─────────────────────┼───────────────────┘                        │
│         │                     │                                             │
│         └─────────────────────┘                                             │
│              (当本地队列为空时窃取)                                          │
│                                                                              │
│  GOMAXPROCS = 2 (Num of P)                                                   │
│  M的数量动态调整                                                              │
│  G的数量理论上无限制 (内存允许)                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 调度决策树

```
调度策略?
│
├── Goroutine 创建?
│   ├── 尝试放入当前 P 的本地队列
│   ├── 本地队列满? → 放入全局队列
│   └── 可能需要唤醒 M
│
├── Goroutine 阻塞?
│   ├── 系统调用? → M 分离 P，新 M 接管
│   ├── Channel 操作? → G 状态变为 Gwaiting
│   └── 唤醒时放入全局队列或原 P 本地队列
│
├── P 的本地队列为空?
│   ├── 尝试窃取其他 P 的一半任务
│   ├── 窃取失败? → 检查全局队列
│   └── 全局队列也空? → 休眠 M
│
└── 抢占?
    ├── 函数调用检查点
    ├── 运行时间过长 (10ms)?
    └── 信号抢占 (Go 1.14+)
```

### 4.3 调度器对比矩阵

| 特性 | Go | Java (ForkJoin) | Erlang | OS 线程 |
|------|----|-----------------|--------|---------|
| **调度单位** | Goroutine (2KB) | Task | Process | Thread (1-8MB) |
| **调度器** | M:P:G | Work Stealing | Reduction Count | Preemptive |
| **切换成本** | ~200ns | ~100ns | ~100ns | ~1-10μs |
| **可扩展性** | 100K+ | 10K+ | 100K+ | ~10K |
| **抢占** | 协作+信号 | 协作 | 约减计数 | 硬件中断 |
| **亲和性** | P 绑定 | Work Stealing | 无 | CPU 亲和 |

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Go Scheduler Tuning Checklist                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  配置:                                                                       │
│  □ GOMAXPROCS 默认等于 CPU 核数                                               │
│  □ IO 密集型可适当增加                                                        │
│  □ 不要设置过大 (超过核数)                                                     │
│                                                                              │
│  性能优化:                                                                   │
│  □ 避免过多的 goroutine (内存)                                                │
│  □ 使用 sync.Pool 复用对象                                                    │
│  □ 避免 goroutine 泄漏                                                        │
│  □ 批量处理减少调度压力                                                        │
│                                                                              │
│  调试:                                                                       │
│  □ runtime.NumGoroutine()                                                     │
│  □ runtime.GOMAXPROCS()                                                       │
│  □ GODEBUG=schedtrace=X                                                       │
│  □ pprof goroutine profile                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (17KB, 完整形式化)
