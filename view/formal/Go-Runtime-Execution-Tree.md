# Go-Runtime-Execution-Tree: Go运行时GMP数据控制执行树 (2025深度版)

## 1. Go程序执行生命周期树

```
[go run main.go / ./executable]
    │
    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ 1. 操作系统加载器                                                        │
│    • 解析ELF/Mach-O/PE格式                                              │
│    • 分配虚拟地址空间                                                    │
│    • 设置初始栈                                                          │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ 2. 运行时初始化 (rt0_go)                                                 │
│    • 汇编入口: runtime/asm_amd64.s:_rt0_amd64                            │
│    • 设置g0栈 (16KB-48KB，固定大小，用于调度器)                          │
│    • 初始化m0 (主OS线程)                                                 │
│    • 调用runtime·check (架构检查)                                        │
│    • 调用runtime·args (解析命令行参数)                                    │
│    • 调用runtime·osinit (OS初始化，读取CPU核心数)                         │
│    • 调用runtime·schedinit (调度器初始化)                                 │
│         - 初始化全局调度器schedt                                         │
│         - 创建allp数组 (GOMAXPROCS个P)                                   │
│         - 设置m0与p0绑定                                                 │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ 3. 创建main goroutine                                                    │
│    • newproc(runtime·main)                                               │
│    • 将main goroutine放入p0本地队列                                      │
│    • mstart (启动m0执行调度循环)                                          │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ 4. 程序执行阶段                                                          │
│    • runtime·main:                                                       │
│         - 调用runtime·init (包初始化，执行所有init函数)                    │
│         - 调用main·main (用户main函数)                                    │
│    • GMP调度器持续工作                                                    │
│    • GC后台 goroutine运行                                                │
│    • sysmon后台监控线程运行                                               │
└─────────────────────────────────┬───────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────────┐
│ 5. 程序退出                                                              │
│    • main·main返回                                                       │
│    • 执行defer链 (LIFO顺序)                                              │
│    • 调用runtime·exit                                                    │
│    • 调用exit系统调用                                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

## 2. GMP调度器形式化模型

### 2.1 GMP状态空间定义

```
GMP状态机: M = (G, M, P, schedt, →)

G = {g₀, g₁, g₂, ...}  (Goroutine集合)
每个g ∈ G包含:
  - stack: 栈内存 (2KB起始，可增长)
  - status: {_Gidle, _Grunnable, _Grunning, _Gsyscall, _Gwaiting, _Gdead}
  - m: 当前绑定的M指针 (可为nil)
  - sched: 保存的寄存器上下文 (SP, PC, etc.)

M = {m₀, m₁, m₂, ...}  (OS线程集合)
每个m ∈ M包含:
  - g0: 调度栈goroutine
  - curg: 当前运行的用户goroutine
  - p: 当前绑定的P指针 (可为nil)
  - lockedg: 锁定的goroutine (G锁定到M)

P = {p₀, p₁, ..., p_{GOMAXPROCS-1}}  (逻辑处理器)
每个p ∈ P包含:
  - runq: [256]guintptr  (本地运行队列，循环数组)
  - runqhead, runqtail: 队列头尾指针
  - runnext: gpointer (快速路径，下一个运行的G)
  - gFree: gList (空闲G缓存)
  - mcache: *mcache (内存分配缓存)
  - status: {_Pidle, _Prunning, _Psyscall, _Pgcstop}

schedt (全局调度器):
  - runq: gQueue (全局运行队列)
  - runqsize: int32
  - pidle: pList (空闲P列表)
  - midle: mList (空闲M列表)
  - nmspinning: int32 (自旋M计数)
  - lock: mutex (保护共享状态)
```

### 2.2 调度循环形式化

```
调度主循环 (runtime·schedule):

func schedule() {
    // _g_ = g0 (调度goroutine)

    // Step 1: 尝试从本地队列获取
    gp := runqget(_p_)
    if gp != nil {
        goto execute
    }

    // Step 2: 每61次检查全局队列
    if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
        gp = globrunqget(_p_, 1)
        if gp != nil {
            goto execute
        }
    }

    // Step 3: 工作窃取
    gp = findrunnable()  // 阻塞直到找到可运行G

execute:
    execute(gp)  // 切换到gp执行
}

形式化转移关系:
  ⟨G, M, P, schedt⟩ ─schedule()→ ⟨G', M', P', schedt'⟩

其中状态转换:
  - runqget: 从p.runq取出G，原子操作
  - globrunqget: 锁定sched.lock，转移G到本地队列
  - findrunnable: 复杂搜索，可能阻塞(M自旋或休眠)
```

### 2.3 工作窃取算法决策树

```
[findrunnable() 入口]
    │
    ▼
[检查本地runnext插槽]
    │
    ├── 有G ──► 返回该G
    │
    └── 无 ───► [检查本地队列]
                    │
                    ├── 非空 ──► 返回队列头部G
                    │
                    └── 空 ───► [检查全局队列]
                                    │
                                    ├── 非空 ──► 转移一批到本地，返回一个
                                    │             (转移数量: n = min(len(global)/GOMAXPROCS + 1, len(local)/2))
                                    │
                                    └── 空 ───► [工作窃取]
                                                    │
                                                    ├── 随机选择其他P
                                                    │   ├── 成功窃取 ──► 返回窃取的G
                                                    │   └── 失败 ─────► [尝试其他P]
                                                    │
                                                    ├── 所有P都无工作 ──► [检查netpoller]
                                                    │                        │
                                                    │                        ├── 有就绪网络事件 ──► 返回对应G
                                                    │                        │
                                                    │                        └── 无 ──► [休眠或自旋决策]
                                                    │                                        │
                                                    │                                        ├── 有自旋M ──► M休眠(park)
                                                    │                                        │
                                                    │                                        └── 无自旋M ──► M自旋等待
                                                    │
                                                    └── 窃取成功 ──► 返回

窃取规则:
  - 从 victim P 的 runq 窃取一半G
  - 原子操作，保证线程安全
  - 如果runq为空，尝试窃取runnext
```

## 3. Channel操作执行树

### 3.1 Channel数据结构

```
type hchan struct {
    qcount   uint           // 队列中元素总数
    dataqsiz uint           // 循环队列大小 (0表示无缓冲)
    buf      unsafe.Pointer // 指向循环队列数组
    elemsize uint16         // 元素大小
    closed   uint32         // 关闭标志 (0=开放, 1=关闭)

    // 等待队列
    sendq    waitq          // 发送者等待队列 (sudog链表)
    recvq    waitq          // 接收者等待队列 (sudog链表)

    // 锁保护
    lock mutex
}

type sudog struct {
    // 代表等待的goroutine
    g          *g
    isSelect   bool
    elem       unsafe.Pointer // 数据元素指针
    c          *hchan          // 指向channel
    ...
}
```

### 3.2 发送操作详细流程

```
[ch <- v] 操作:
    │
    ▼
[调用 runtime·chansend1 或 chansend]
    │
    ▼
[检查ch是否为nil]
    │
    ├── 是nil ──► 永久阻塞 (panic如果是close后的channel)
    │
    └── 非nil ──► [加锁: lock(&hchan.lock)]
                        │
                        ▼
                [检查channel是否已关闭]
                        │
                        ├── 已关闭 ──► 解锁，panic("send on closed channel")
                        │
                        └── 未关闭 ──► [检查recvq是否有等待者]
                                            │
                                            ├── 有等待者 ──► [直接传递]
                                            │                  │
                                            │                  ▼
                                            │             [从recvq取出sudog]
                                            │                  │
                                            │                  ▼
                                            │             [复制v到sudog.elem]
                                            │                  │
                                            │                  ▼
                                            │             [唤醒接收者goready(sg.g)]
                                            │                  │
                                            │                  ▼
                                            │             [解锁，返回true]
                                            │
                                            └── 无等待者 ──► [检查缓冲区]
                                                                │
                                                                ├── 缓冲区有空间 ──► [入队]
                                                                │                      │
                                                                │                      ▼
                                                                │                 [复制到buf]
                                                                │                      │
                                                                │                      ▼
                                                                │                 [递增qcount]
                                                                │                      │
                                                                │                      ▼
                                                                │                 [解锁，返回true]
                                                                │
                                                                └── 缓冲区满 ──► [阻塞发送者]
                                                                                    │
                                                                                    ▼
                                                                               [创建sudog]
                                                                                    │
                                                                                    ▼
                                                                               [sudog.g = 当前g]
                                                                                    │
                                                                                    ▼
                                                                               [sudog.elem = &v]
                                                                                    │
                                                                                    ▼
                                                                               [加入sendq链表]
                                                                                    │
                                                                                    ▼
                                                                               [goparkunlock(&ch.lock)]
                                                                               (休眠，解锁channel)
                                                                                    │
                                                                                    ▼
                                                                               [被唤醒后继续]
                                                                                    │
                                                                                    ▼
                                                                               [返回成功]
```

### 3.3 Select语句执行树

```
[select { case c1<-v1: ... case c2:=<-c3: ... default: ... }]
    │
    ▼
[编译器转换]
    │
    ├── 无default ──► 阻塞select
    │   └── 调用 runtime·selectgo
    │
    └── 有default ──► 非阻塞select
        └── 调用 runtime·selectnbexec

[selectgo 执行流程]:
    │
    ▼
[收集所有case]
    │
    ▼
[随机打乱case顺序 (保证公平性)]
    │
    ▼
[第一轮: 非阻塞尝试]
    │
    ▼
[遍历每个case]
    │
    ├── 发送case ──► [尝试chansend]
    │   ├── 成功 ──► [执行case body]
    │   └── 失败 ──► [继续下一个case]
    │
    ├── 接收case ──► [尝试chanrecv]
    │   ├── 成功 ──► [执行case body]
    │   └── 失败 ──► [继续下一个case]
    │
    └── 遍历完成 ──► [无成功case?]
                            │
                            ├── 有default ──► [执行default]
                            │
                            └── 无default ──► [第二轮: 阻塞等待]
                                                    │
                                                    ▼
                                            [为所有channel创建sudog]
                                                    │
                                                    ▼
                                            [加入各channel的等待队列]
                                                    │
                                                    ▼
                                            [gopark (休眠)]
                                                    │
                                                    ▼
                                            [被其中一个channel唤醒]
                                                    │
                                                    ▼
                                            [从其他channel移除sudog]
                                                    │
                                                    ▼
                                            [执行成功的case]

关键实现细节:
  - sudog被所有涉及的channel共享
  - 唤醒后需要清理其他channel的等待队列
  - isSelect标志用于特殊处理
```

## 4. 系统调用处理决策树

```
[goroutine执行系统调用]
    │
    ▼
[调用entersyscall]
    │
    ▼
[保存goroutine状态]
    │
    ▼
[设置g.status = _Gsyscall]
    │
    ▼
[设置m.curg = nil]
    │
    ▼
[释放P: releasep()]
    │
    ├── 将P状态设为_Psyscall
    │
    ├── 将P从M分离
    │
    └── 如果P.runq非空 ──► [唤醒空闲M或创建新M执行P]

[系统调用执行中]
    │ (OS内核处理)
    ▼

[系统调用返回]
    │
    ▼
[调用exitsyscall]
    │
    ▼
[尝试重新获取原来的P]
    │
    ├── P可用且idle ──► [acquirep] ──► [继续执行]
    │
    └── P不可用 ──► [获取其他空闲P]
                        │
                        ├── 成功 ──► [继续执行]
                        │
                        └── 失败 ──► [将G放入全局队列]
                                            │
                                            ▼
                                    [M休眠等待工作]
                                            │
                                            ▼
                                    [被sysmon或调度器唤醒]

sysmon角色:
  - 独立线程，周期运行(10-20ms)
  - 检测syscall超时的M
  - 长时间syscall后抢占P
  - 触发GC
  - 强制抢占长时间运行的G
```

## 5. 形式化对应关系

### 5.1 Go ↔ 理论模型对应

```
Go概念                    理论对应                      形式化关系
─────────────────────────────────────────────────────────────────────────────
Goroutine (G)        ⟷  CSP进程 / Actor             近似对应，但有调度器管理
                     ⟷  λ演算中的轻量级线程

Channel (缓冲)       ⟷  CSP通道 + 队列              Go的缓冲channel超越CSP同步语义
                     ⟷  Actor mailbox (每个channel) 异步语义

Channel (无缓冲)     ⟷  CSP同步通道                 直接对应

Select               ⟷  CSP外部选择□ + 内部选择⊓    混合语义，Go的select是随机公平选择

Mutex/RWMutex        ⟷  π-calculus限制              底层同步原语，非CSP概念

WaitGroup            ⟷  进程组合join                类似并行组合的同步终止

Context              ⟷  监督树(简化)                类似Actor的容错信号传递

GMP调度器            ⟷  无直接对应                   实现层面的优化，非语义层面

work stealing        ⟷  任务窃取调度               负载均衡策略
```

### 5.2 执行迹 ↔ SOS推导

```
执行迹示例:
  tr = [ch<-1, <-ch, go f(), g()完成]

对应的SOS推导链:
  ⟨ch<-1, σ, Ξ⟩ ─τ→ ⟨0, σ', Ξ'⟩    (发送动作)
  ⟨<-ch, σ', Ξ'⟩ ─1→ ⟨0, σ'', Ξ''⟩ (接收值1)
  ⟨go f(), σ'', Ξ''⟩ ─τ→ ⟨0, σ'', Ξ''∪{f()}⟩ (创建goroutine)
  ...

观察等价:
  外部观察只能看到:
    - channel通信的值
    - IO操作
    - 最终的程序结果

  内部τ动作(调度、G切换)对外部不可见

  因此: 多个不同内部执行路径可能迹等价
```

---

**参考文献**:

1. Go Runtime源码: runtime/proc.go, runtime/chan.go
2. "Understanding the Go Runtime" (2026 internals-for-interns.com)
3. A. Gerrand: "Go's work-stealing scheduler"
4. D. Vyukov: Original GMP design documents
5. Kogan et al.: "MCS: A Work-Stealing Scheduler for Go"

**关联文档**:

- [formal/Go-CSP-Formal](../formal/Go-CSP-Formal.md)
- [TH1-Process-Calculi-MindMap](../01-core-theory/TH1-Process-Calculi-MindMap.md)
