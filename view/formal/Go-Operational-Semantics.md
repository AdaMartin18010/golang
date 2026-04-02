# Go操作语义：Goroutine调度与执行

## 执行树：Goroutine生命周期与调度

```
                    Goroutine生命周期
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
     创建(go)           运行               结束
        │                │                  │
        ▼                ▼                  ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
| 分配G结构体   │  | 调度到P队列   │  | 正常返回      │
| 初始化栈(2KB)│  | M(线程)绑定  │  | panic/recover │
| 加入P本地队列│  | 上下文切换   │  | 栈扩容/收缩   │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┼─────────────────┘
                         ▼
            ┌─────────────────────────┐
            │    GMP调度模型          │
            │  G(Goroutine)           │
            │  M(OS Thread)           │
            │  P(Logical Processor)   │
            │                         │
            │  P的本地队列 + 全局队列   │
            │  Work Stealing机制      │
            └─────────────────────────┘
                         │
                         ▼
            ┌─────────────────────────┐
            │    调度决策点            │
            │  • 阻塞系统调用          │
            │  • Channel操作          │
            │  • time.Sleep          │
            │  • 手动yield            │
            │  • 时间片耗尽           │
            └─────────────────────────┘
```

---

## 1. 形式化定义

### 1.1 GMP模型

**定义 1.1.1** (GMP抽象).
$$
\text{Scheduler} = (G, M, P, \text{runq}, \text{globalq})
$$
其中：

- **G**: Goroutine集合，每个 **g ∈ G** 有状态 **{idle, runnable, running, blocked, dead}**
- **M**: 机器线程(OS线程)集合
- **P**: 逻辑处理器集合，**|P| = GOMAXPROCS**
- **runq: P → G/*** : 每个P的本地可运行队列
- **globalq**: 全局可运行队列

**定义 1.1.2** (Goroutine状态).
$$
\text{State}_g \in \{\text{idle}, \text{runnable}, \text{running}, \text{blocked}, \text{dead}\}
$$

**BNF (Goroutine相关操作)**:

```
<goroutine_op>    ::= "go" <function_call>                  (* 创建 *)
                    | "runtime.Gosched" "()"               (* 让出 *)
                    | "runtime.Goexit" "()"                (* 退出 *)
                    | <channel_op>                         (* 通道操作 *)
                    | <mutex_op>                           (* 锁操作 *)
                    | <select_stmt>                        (* 选择 *)

<channel_op>      ::= <channel> "<-" <value>               (* 发送 *)
                    | <variable> "<-" <channel>            (* 接收 *)
                    | "close" "(" <channel> ")"            (* 关闭 *)
                    | <range_over_channel>                 (* 遍历 *)

<select_stmt>     ::= "select" "{" <select_case>+ "}"
<select_case>     ::= "case" <communication> ":" <statements>
                    | "default" ":" <statements>
<communication>   ::= <channel_op> | <assignment_from_chan>

<mutex_op>        ::= <mutex> "." "Lock" "()"
                    | <mutex> "." "Unlock" "()"
                    | <rwmutex> "." "RLock" "()"
                    | <rwmutex> "." "RUnlock" "()"
```

---

## 2. Goroutine调度语义

### 2.1 创建语义

**规则 2.1.1** (go语句).
$$
\frac{\text{go } f() \text{ in } g_{parent} \quad g_{new} \text{ fresh}}{\langle \text{go } f(), g_{parent}, P \rangle \longrightarrow \langle \epsilon, g_{parent}, P[\text{runq} := \text{runq}(P) \cup \{g_{new}\}] \rangle}
$$
其中：

- **g_new.stack = 2KB** (初始栈)
- **g_new.state = runnable**
- **g_new.parent = g_parent**

### 2.2 调度决策

**规则 2.2.1** (本地调度).
$$
\frac{g \in \text{runq}(P) \quad m \text{ bound to } P}{\langle P, m \rangle \xrightarrow{\text{schedule}} \langle P[g \mapsto \text{running}], m \rangle \circ g}
$$

**规则 2.2.2** (Work Stealing).
$$
\frac{\text{runq}(P_i) = \emptyset \quad g \in \text{runq}(P_j) \quad i \neq j}{P_i \xrightarrow{\text{steal}} P_i[\text{runq} := \{g\}] \circ P_j[\text{runq} := \text{runq}(P_j) \setminus \{g\}]}
$$

**规则 2.2.3** (全局队列获取).
$$
\frac{\text{runq}(P) = \emptyset \quad g \in \text{globalq}}{P \xrightarrow{\text{getglobal}} P[\text{runq} := \{g\}] \circ \text{globalq} := \text{globalq} \setminus \{g\}}
$$

### 2.3 状态转换

**BNF (状态机)**:

```
<g_state>         ::= "idle" | "runnable" | "running" | "blocked" | "dead"

<state_transition> ::= "idle" --"create"--> "runnable"
                     | "runnable" --"schedule"--> "running"
                     | "running" --"yield"--> "runnable"
                     | "running" --"block"--> "blocked"
                     | "running" --"complete"--> "dead"
                     | "blocked" --"unblock"--> "runnable"
                     | "blocked" --"timeout"--> "runnable"
                     | "running" --"panic"--> "dead"
                     | "running" --"systemcall"--> "blocked"
```

---

## 3. Channel操作语义

### 3.1 通道抽象

**定义 3.1.1** (Channel状态).
Channel **ch = (buf, cap, sendq, recvq, closed)**：

- **buf**: 循环缓冲区
- **cap**: 容量（0表示无缓冲）
- **sendq**: 等待发送的Goroutine队列
- **recvq**: 等待接收的Goroutine队列
- **closed**: 是否已关闭

### 3.2 发送操作

**规则 3.2.1** (直接发送 - 有等待接收者).
$$
\frac{ch.recvq \neq \emptyset \quad g_{recv} = \text{dequeue}(ch.recvq)}{\langle ch \leftarrow v, g_{sender} \rangle \longrightarrow \langle g_{recv}.\text{val} := v, g_{sender} \rangle \circ g_{recv}.\text{state} := \text{runnable}}
$$

**规则 3.2.2** (缓冲写入 - 有空间).
$$
\frac{ch.recvq = \emptyset \quad |ch.buf| < ch.cap}{\langle ch \leftarrow v, g \rangle \longrightarrow \langle ch.buf := ch.buf \cup \{v\}, g \rangle}
$$

**规则 3.2.3** (阻塞发送 - 无空间).
$$
\frac{ch.recvq = \emptyset \quad |ch.buf| = ch.cap}{\langle ch \leftarrow v, g \rangle \longrightarrow \langle ch.sendq := ch.sendq \cup \{g\}, g.\text{state} := \text{blocked} \rangle}
$$

### 3.3 接收操作

**规则 3.3.1** (直接从缓冲区接收).
$$
\frac{ch.buf \neq \emptyset}{\langle x \leftarrow ch, g \rangle \longrightarrow \langle x := \text{dequeue}(ch.buf), g \rangle}
$$

**规则 3.3.2** (从发送者直接接收 - 无缓冲或空缓冲).
$$
\frac{ch.buf = \emptyset \quad ch.sendq \neq \emptyset \quad g_{send} = \text{dequeue}(ch.sendq)}{\langle x \leftarrow ch, g_{recv} \rangle \longrightarrow \langle x := g_{send}.\text{val}, g_{send}.\text{state} := \text{runnable} \rangle}
$$

### 3.4 关闭操作

**规则 3.4.1** (关闭Channel).
$$
\frac{\neg ch.closed}{\langle \text{close}(ch), g \rangle \longrightarrow \langle ch.closed := \text{true} \rangle \circ \forall g' \in ch.recvq. g'.\text{state} := \text{runnable}}
$$

---

## 4. Select语句语义

### 4.1 非确定性选择

**定义 4.1.1** (Select状态).
Select **sel = (cases, default?)**，其中每个 **case** 关联一个通信操作。

**算法 4.1.2** (Select执行).

```
function ExecuteSelect(sel, g):
    ready ← {c ∈ sel.cases | CanProceedWithoutBlock(c)}

    if ready ≠ ∅:
        // 随机选择一个就绪case
        c ← RandomSelect(ready)
        Execute(c)
    else if sel.default ≠ nil:
        Execute(sel.default)
    else:
        // 所有case阻塞，注册等待
        for c ∈ sel.cases:
            RegisterWait(c, g)
        g.state ← blocked
```

---

## 5. 系统级论证

### 5.1 调度公平性

**定理 5.1.1** (Goroutine饥饿自由).
在有限Goroutine假设下，每个可运行的Goroutine最终被执行。

**证明概要**:

1. 每个P定期从全局队列获取Goroutine
2. Work stealing确保空闲P从忙碌P窃取
3. 没有Goroutine被无限期跳过 ∎

### 5.2 Channel安全性

**定理 5.2.1** (Channel类型安全).
正确类型的Channel操作不会导致运行时类型错误。

**定理 5.2.2** (无死锁检测).
Go不保证检测所有死锁，但保证无数据竞争程序有定义行为。

---

## 6. 多维属性矩阵

| 特性 | Goroutine | Channel | Mutex | Select |
|------|-----------|---------|-------|--------|
| 创建开销 | 2KB栈 | 内存分配 | 内存分配 | 编译时 |
| 切换开销 | ~200ns | 阻塞时切换 | 阻塞时切换 | case评估 |
| 公平性 | 是 | FIFO | 非公平 | 随机 |
| 可组合性 | 高 | 高 | 中 | 高 |
| 死锁风险 | 无 | 有 | 有 | 有 |

---

## 参考文献

1. Go Authors. (2024). *Go Scheduler Internals*. go.dev.
2. Go Authors. (2024). *Go Memory Model*. go.dev/ref/mem.
3. Vyukov, D. *Scalable Go Scheduler Design*. (GMP设计文档)
4. Hoare, C.A.R. (1978). *Communicating Sequential Processes*.

---

*文档版本: 2025年 | 理论深度: L4 | 形式化等级: 完整SOS*
