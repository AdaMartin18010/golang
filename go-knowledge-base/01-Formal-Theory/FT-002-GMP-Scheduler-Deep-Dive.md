# FT-002: GMP 调度器深度分析 (GMP Scheduler Deep Dive)

> **维度**: Formal Theory
> **级别**: S (30+ KB)
> **标签**: #scheduler #gmp #goroutine #os-thread
> **权威来源**: [Go Scheduler](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html), [Go Runtime Source](https://github.com/golang/go/tree/master/src/runtime)

---

## GMP 模型概述

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          GMP Scheduler Model                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  G (Goroutine)    M (Machine/OS Thread)    P (Processor/Logical CPU)        │
│  ─────────────    ─────────────────────    ──────────────────────────        │
│                                                                              │
│  ┌─────────┐      ┌─────────┐              ┌─────────┐                     │
│  │  Stack  │      │  Stack  │              │ RunQ    │                     │
│  │  ~2KB   │      │  ~8KB   │              │ (Local) │                     │
│  │  grow   │      │  fixed  │              │         │                     │
│  └────┬────┘      └────┬────┘              └────┬────┘                     │
│       │                │                        │                          │
│       └────────────────┴────────────────────────┘                          │
│                        │                                                   │
│                   Global RunQ                                               │
│                   Work Stealing                                             │
│                   Network Poller                                            │
│                                                                              │
│  G count: 100k+    M count: ~CPU cores       P count: GOMAXPROCS            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心数据结构

### G (Goroutine)

```go
// src/runtime/runtime2.go

type g struct {
    stack       stack          // 栈信息
    sched       gobuf          // 调度上下文
    goid        int64          // 唯一 ID
    atomicstatus uint32        // 状态
    m            *m             // 绑定的 M
    p            uintptr        // 绑定的 P

    // 调度相关
    schedlink    guintptr       // 调度链表
    waitsince    int64          // 等待时间
    waitreason   waitReason     // 等待原因

    // 栈增长
    stackguard0  uintptr
    stackguard1  uintptr
}

// G 状态
const (
    _Gidle = iota      // 刚分配
    _Grunnable        // 可运行
    _Grunning         // 运行中
    _Gsyscall         // 系统调用
    _Gwaiting         // 等待（阻塞）
    _Gdead            // 已退出
    _Gcopystack       // 栈拷贝中
    _Gpreempted       // 被抢占
)
```

### M (Machine)

```go
type m struct {
    g0           *g          // 调度 goroutine
    curg         *g          // 当前运行的 goroutine
    p            puintptr    // 绑定的 P
    nextp        puintptr    // 下一个 P

    // 系统调用
    syscalltick  uint32
    syscallsp    uintptr
    syscallpc    uintptr

    // 线程 ID
    procid       uint64

    // 信号处理
    gsignal      *g
    sigmask      sigset
}
```

### P (Processor)

```go
type p struct {
    id           int32
    status       uint32

    // 运行队列
    runqhead     uint32
    runqtail     uint32
    runq         [256]guintptr  // 本地队列（循环数组）

    // 下一个 G
    runnext      guintptr

    // 空闲 G 列表
    gFree        struct {
        gList
        n int32
    }

    // 内存分配
    mcache       *mcache

    // 统计
    stats        struct {
        // ...
    }
}
```

---

## 调度循环 (Scheduler Loop)

```go
// src/runtime/proc.go

// The main scheduler loop
func schedule() {
    _g_ := getg()

    // 获取当前 P
    _p_ := _g_.m.p.ptr()

    // 1. 尝试从本地队列获取
    if gp, inheritTime := runqget(_p_); gp != nil {
        execute(gp, inheritTime)
        return
    }

    // 2. 尝试从全局队列获取
    if gp := globrunqget(_p_, 0); gp != nil {
        execute(gp, false)
        return
    }

    // 3. 网络轮询（非阻塞）
    if netpollinited() && netpollWaiters.Load() > 0 {
        if list := netpoll(0); !list.empty() {
            gp := list.pop()
            injectglist(&list)
            execute(gp, false)
            return
        }
    }

    // 4. Work Stealing
    if gp := stealWork(now); gp != nil {
        execute(gp, false)
        return
    }

    // 5. 没有工作，进入空闲
    stopm()
}
```

---

## Work Stealing 算法

```go
// src/runtime/proc.go

func stealWork(now int64) *g {
    _p_ := getg().m.p.ptr()

    // 随机顺序遍历其他 P
    for i := 0; i < 4; i++ {
        // 1. 尝试偷取 runnext（下一个要运行的 G）
        if p2 := allp[stealOrder[_p_.id][i]]; p2 != _p_ {
            if gp := runqsteal(_p_, p2); gp != nil {
                return gp
            }
        }
    }

    // 2. 尝试从全局队列偷取
    if gp := globrunqget(_p_, 1); gp != nil {
        return gp
    }

    // 3. 尝试从网络轮询器偷取
    if netpollinited() && atomic.Load(&netpollWaiters) > 0 {
        if list := netpoll(0); !list.empty() {
            gp := list.pop()
            injectglist(&list)
            return gp
        }
    }

    return nil
}

// 偷取半个队列
func runqsteal(_p_, p2 *p) *g {
    // 偷取 p2.runq 的一半
    n := runqgrabsize(p2)
    batch := runqgrab(p2, &p2.runq, n)

    // 将偷取的放入 _p_ 的队列
    for i := 0; i < n; i++ {
        runqput(_p_, batch[i], false)
    }

    return batch[n-1]  // 返回最后一个直接运行
}
```

---

## 系统调用处理

```go
// 进入系统调用
func entersyscall() {
    _g_ := getg()
    _p_ := _g_.m.p.ptr()

    // 保存状态
    save(_g_)

    // 释放 P，让其他 M 可以运行
    _g_.m.p = 0
    atomic.Store(&_p_.status, _Psyscall)

    // 如果有空闲 P，唤醒其他 M
    if sched.npidle < sched.nmspinning {
        wakep()
    }
}

// 退出系统调用
func exitsyscall() {
    _g_ := getg()
    _p_ := _g_.m.oldp.ptr()

    // 尝试重新获取原来的 P
    if cas(&_p_.status, _Psyscall, _Pidle) {
        // 获取成功
        acquirep(_p_)
    } else {
        // P 已被偷走，获取新的 P
        _p_ = acquirep1()
    }

    _g_.m.p.set(_p_)
}
```

---

## 抢占式调度

```go
// src/runtime/preempt.go

// 信号驱动的抢占
func preemptM(mp *m) {
    // 发送 SIGURG 信号
    signalM(mp, sigPreempt)
}

// 信号处理
func doSigPreempt(gp *g, ctxt *sigctxt) {
    // 检查是否可以安全抢占
    if canPreemptM(gp.m) {
        // 设置抢占标志
        gp.preempt = true

        // 在函数调用边界检查并执行抢占
        // 实际抢占发生在 morestack 或同步点
    }
}
```

### 抢占触发点

1. **函数调用**: 检查 `gp.preempt`
2. **系统调用**: 返回时检查
3. **GC**: STW 前抢占所有 G

---

## 性能优化

### 调度延迟分解

| 阶段 | 时间量级 | 优化手段 |
|------|---------|---------|
| 本地队列获取 | ~10ns | 无锁设计 |
| 全局队列获取 | ~100ns | 批量获取 |
| Work Stealing | ~1μs | 随机化减少冲突 |
| 创建新 M | ~10μs | M 池化 |
| 栈增长 | ~100μs | 分段栈 |

### 调度亲和性

```go
// 保持 G 在同一个 P 上运行
func lockOSThread() {
    _g_ := getg()
    _g_.m.lockedg.set(_g_)
    _g_.lockedm.set(_g_.m)
}

// 解除绑定
func unlockOSThread() {
    _g_ := getg()
    _g_.m.lockedg = 0
    _g_.lockedm = 0
}
```

---

## 与 OS 调度器对比

| 特性 | Go Scheduler | Linux CFS |
|------|-------------|-----------|
| 调度单位 | Goroutine (~2KB) | Thread (~8MB) |
| 切换开销 | ~200ns | ~1-2μs |
| 调度策略 | Work Stealing + M:N | CFS (红黑树) |
| 抢占 | 协作式 + 信号 | 时间片 |
| 核心数支持 | 数百万 G | 数千线程 |

---

## 参考文献

1. [Scheduling In Go - Part II](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html) - William Kennedy
2. [Go Runtime Scheduler](https://github.com/golang/go/tree/master/src/runtime) - Go Source
3. [Analysis of the Go Runtime Scheduler](http://www1.cs.columbia.edu/~aho/cs6998/reports/12-12-11_DeshpandeSponslerWeiss_GO.pdf) - Columbia University
4. [The Linux Scheduler](https://docs.kernel.org/scheduler/sched-design-CFS.html) - Kernel Documentation
