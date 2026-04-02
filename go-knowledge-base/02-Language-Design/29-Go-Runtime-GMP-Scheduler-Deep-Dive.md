# Go Runtime GMP 调度器深度剖析 (Go Runtime GMP Scheduler Deep Dive)

> **分类**: 语言设计
> **标签**: #runtime #scheduler #GMP #goroutine
> **参考**: Go 1.21-1.24 Runtime, src/runtime/proc.go, src/runtime/runtime2.go

---

## GMP 模型架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Go Runtime GMP Model                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Global Runtime (schedt)                       │   │
│  │  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐   │   │
│  │  │   Global RunQ    │  │    Idle List     │  │   GC Work Queue  │   │   │
│  │  │   (Lock-free)    │  │   (M & P pools)  │  │                  │   │   │
│  │  └──────────────────┘  └──────────────────┘  └──────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────┼─────────────────────────────────────┐   │
│  │                                 ▼                                     │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │   │
│  │  │      P0     │◄──►│      P1     │◄──►│      P2     │  ...         │   │
│  │  │ (Processor) │    │ (Processor) │    │ (Processor) │              │   │
│  │  ├─────────────┤    ├─────────────┤    ├─────────────┤              │   │
│  │  │ Local RunQ  │    │ Local RunQ  │    │ Local RunQ  │              │   │
│  │  │  (256 max)  │    │  (256 max)  │    │  (256 max)  │              │   │
│  │  ├─────────────┤    ├─────────────┤    ├─────────────┤              │   │
│  │  │   mcache    │    │   mcache    │    │   mcache    │              │   │
│  │  ├─────────────┤    ├─────────────┤    ├─────────────┤              │   │
│  │  │    runnext  │    │    runnext  │    │    runnext  │              │   │
│  │  │  (高优先级)  │    │  (高优先级)  │    │  (高优先级)  │              │   │
│  │  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘              │   │
│  │         │                  │                  │                       │   │
│  │         └──────────────────┼──────────────────┘                       │   │
│  │                            │                                          │   │
│  │         ┌──────────────────┼──────────────────┐                       │   │
│  │         ▼                  ▼                  ▼                       │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │   │
│  │  │      M0     │    │      M1     │    │      M2     │              │   │
│  │  │   (OS Thread)│    │   (OS Thread)│    │   (OS Thread)│              │   │
│  │  ├─────────────┤    ├─────────────┤    ├─────────────┤              │   │
│  │  │     g0      │    │     g0      │    │     g0      │              │   │
│  │  │ (调度栈 2KB) │    │ (调度栈 2KB) │    │ (调度栈 2KB) │              │   │
│  │  ├─────────────┤    ├─────────────┤    ├─────────────┤              │   │
│  │  │  curg (->G) │    │  curg (->G) │    │  curg (->G) │              │   │
│  │  └─────────────┘    └─────────────┘    └─────────────┘              │   │
│  │         │                  │                  │                       │   │
│  │         ▼                  ▼                  ▼                       │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │   │
│  │  │      G      │    │      G      │    │      G      │              │   │
│  │  │ (Goroutine) │    │ (Goroutine) │    │ (Goroutine) │              │   │
│  │  ├─────────────┤    ├─────────────┤    ├─────────────┤              │   │
│  │  │  stack(2KB+)│    │  stack(2KB+)│    │  stack(2KB+)│              │   │
│  │  │  sched(上下文)│    │  sched(上下文)│    │  sched(上下文)│              │   │
│  │  │  status     │    │  status     │    │  status     │              │   │
│  │  └─────────────┘    └─────────────┘    └─────────────┘              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  G = Goroutine (用户协程)  M = Machine (OS 线程)  P = Processor (逻辑处理器)  │
│  P 的数量 = GOMAXPROCS (默认 CPU 核心数)                                      │
│  M 的数量动态增长，最大 10000 (默认)                                          │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## G (Goroutine) 结构详解

```go
// runtime/runtime2.go
type g struct {
    // 栈
    stack       stack   // 栈范围 [stack.lo, stack.hi)
    stackguard0 uintptr // 栈溢出检查
    stackguard1 uintptr // 栈溢出检查 (C 代码)

    // 调度相关
    m              *m      // 当前运行的 M
    sched          gobuf   // 保存的寄存器上下文
    atomicstatus   uint32  // 原子操作的状态
    goid           uint64  // 唯一 ID

    // 链接
    schedlink      guintptr // 调度链表

    // 等待原因
    waitsince      int64   // 开始等待时间
    waitreason     waitReason // 等待原因

    // 抢占
    preempt        bool    // 抢占标志
    preemptStop    bool    // 抢占时停止
    preemptShrink  bool    // 抢占时收缩栈

    // GC 相关
    gcscandone     bool    // GC 扫描完成
    gcscanvalid    bool    // GC 扫描有效

    // 锁
    locks          int32   // 持有的锁数量

    // 其他
    paniconfault   bool    // panic 而非 fault

    // 嵌套调用深度
    nesting        int32   // 统计用途

    // 性能分析
    labels         unsafe.Pointer // 性能分析标签
    timer          *timer         // 计时器

    // 同步等待
    selectDone     uint32  // select 完成标记

    // 协程分组
    coroexit       bool    // 协程退出
    corosleepg     guintptr // 睡眠的 G
}

// 保存的寄存器状态
type gobuf struct {
    sp   uintptr        // 栈指针
    pc   uintptr        // 程序计数器
    g    guintptr       // G 指针
    ctxt unsafe.Pointer // 上下文
    ret  uintptr        // 返回值
    lr   uintptr        // 链接寄存器 (ARM)
    bp   uintptr        // 帧指针 (x86)
}
```

### Goroutine 状态机

```
┌─────────┐   create    ┌─────────┐   schedule   ┌─────────┐
│  NONE   │ ──────────► │ _WAITING│ ───────────► │_RUNNABLE│
└─────────┘             └─────────┘              └────┬────┘
                                                      │
       ┌──────────────────────────────────────────────┘
       │ execute
       ▼
  ┌─────────┐   block    ┌─────────┐   wakeup    ┌─────────┐
  │ _RUNNING │ ─────────► │ _WAITING │ ──────────► │_RUNNABLE│
  └────┬────┘            └─────────┘              └─────────┘
       │
       │ complete
       ▼
  ┌─────────┐
  │  _DEAD  │
  └─────────┘
```

---

## M (Machine) 结构详解

```go
// runtime/runtime2.go
type m struct {
    g0      *g     // 带有调度栈的 G
    morebuf gobuf  // 更多栈时的 gobuf

    // TLS
    tls           [tlsSlots]uintptr // 线程本地存储

    // 当前运行的 G
    curg          *g       // 当前用户 G
    caughtsig     guintptr // 信号处理的 G

    // P 关联
    p             puintptr // 关联的 P
    nextp         puintptr // 下一个 P
    oldp          puintptr // 系统调用前的 P

    // 线程 ID
    id            int64

    // 信号处理
    sighold       uint32   // 信号持有标志

    // 调度锁
    locks         int32    // 持有的锁
    dying         int32    // 正在退出

    // 性能分析
    profilehz     int32    // 分析频率

    // 辅助 GC
    helpgc        int32    // 帮助 GC

    // 旋转状态
    spinning      bool     // 是否在自旋寻找工作

    // 阻塞状态
    blocked       bool     // 是否在阻塞

    // 系统调用
    syscalltick   uint32   // 系统调用计数
    syscallsp     uintptr  // 系统调用时的栈指针
    syscallpc     uintptr  // 系统调用时的 PC

    // 链接
    schedlink     muintptr // 调度链表

    // 锁定时运行的 G
    lockedg       guintptr

    // 创建者
    createstack   [32]uintptr // 创建时的栈跟踪

    // 线程
    thread        uintptr // 线程句柄

    // freelist
    freewait      uint32  // 等待释放

    // 快速随机数
    fastrand      uint64

    // 锁计时
    lockrets      int32   // 锁重试次数

    // 信号队列
    sigmask       sigset  // 信号掩码
}
```

---

## P (Processor) 结构详解

```go
// runtime/runtime2.go
type p struct {
    id          int32    // P ID
    status      uint32   // P 状态
    link        puintptr // 空闲链表

    // 调度
    schedtick   uint32   // 调度计数
    syscalltick uint32   // 系统调用计数
    sysmontick  sysmontick // sysmon 监控

    // M 关联
    m           muintptr // 关联的 M

    // 本地运行队列
    runqhead    uint32   // 队列头
    runqtail    uint32   // 队列尾
    runq        [256]guintptr // 本地 G 队列

    // 下一个运行的 G
    runnext     guintptr // 高优先级 G

    // 空闲 G 列表
    gFree       struct {
        gList  // 列表
        n      int32 // 数量
    }

    // 内存分配
    mcache      *mcache  // 内存缓存

    // 页面缓存
    pcache      pageCache // 页面缓存

    // 计时器
    timersLock  mutex
    timers      []*timer  // 计时器堆
    numTimers   uint32    // 计时器数量
    deletedTimers uint32  // 删除的计时器

    // GC 相关
    gcAssistTime      int64  // GC 辅助时间
    gcFractionalMarkTime int64 // 分数标记时间

    // 统计
    gcw             gcWork  // GC 工作
    wbBuf           writeBarrierBuf // 写屏障缓冲

    // 随机数
    fastrandseed    uint64

    // 串行执行
    runSafePointFn  uint32 // 运行安全点函数

    // 统计
    statsSeq        uint32 // 统计序列

    // 计时器调整
    timer0When      uint64
    timerModifiedEarliest uint64
}

// P 状态
const (
    _Pidle       = iota // 空闲
    _Prunning           // 运行中
    _Psyscall           // 系统调用中
    _Pgcstop            // GC 停止
    _Pdead              // 死亡
)
```

---

## 调度循环详解

```go
// runtime/proc.go

// schedule 是调度器的主循环
func schedule() {
    _g_ := getg()

    if _g_.m.locks != 0 {
        throw("schedule: holding locks")
    }

    if _g_.m.lockedg != 0 {
        // 执行锁定的 G
        stoplockedm()
        execute(_g_.m.lockedg.ptr(), false)
    }

    // 顶层循环
    top:
    pp := _g_.m.p.ptr()

    // 检查 GC
    if sched.gcwaiting != 0 {
        gcstopm()
        goto top
    }

    // 检查安全点
    if pp.runSafePointFn != 0 {
        runSafePointFn()
    }

    // 获取可运行的 G
    var gp *g
    var inheritTime bool

    if gp == nil {
        // 1. 检查 runnext（高优先级）
        if gp, inheritTime = runqget(pp); gp != nil {
            // 从本地队列获取成功
        }
    }

    if gp == nil {
        // 2. 从全局队列获取（每 61 次调度）
        gp = globrunqget(pp, 0)
    }

    if gp == nil {
        // 3. 从网络轮询器获取
        if netpollinited() && atomic.Load(&netpollWaiters) > 0 && atomic.Load64(&sched.lastpoll) != 0 {
            // 非阻塞网络轮询
            if list := netpoll(0); !list.empty() {
                gp = list.pop()
                injectglist(&list)
            }
        }
    }

    if gp == nil {
        // 4. 从其他 P 偷取
        gp, inheritTime = findrunnable()
    }

    if gp == nil {
        // 没有工作，进入空闲
        idle()
        goto top
    }

    // 执行 G
    execute(gp, inheritTime)
}

// findrunnable 寻找可运行的 G
func findrunnable() (gp *g, inheritTime bool) {
    _g_ := getg()

    // 自旋锁检查
    if !_g_.m.spinning {
        // 开始自旋
        _g_.m.spinning = true
        sched.nmspinning.Add(1)
    }

    // 尝试偷取
    for i := 0; i < 4; i++ {
        // 随机选择 P
        for _, p := range stealOrder.randomOrder(_g_.m.p.ptr().id) {
            // 尝试偷取一半的 G
            if gp := runqsteal(_g_.m.p.ptr(), p); gp != nil {
                return gp, false
            }

            // 尝试偷取 runnext
            if gp := runqgrab(p, _g_.m.p.ptr(), 1); gp != nil {
                return gp, false
            }
        }
    }

    // 检查全局队列
    if gp := globrunqget(_g_.m.p.ptr(), 0); gp != nil {
        return gp, false
    }

    // 检查网络轮询（阻塞）
    if netpollinited() && atomic.Load(&netpollWaiters) > 0 {
        if gp := netpoll(-1); gp != nil {
            return gp, false
        }
    }

    // 停止自旋
    if _g_.m.spinning {
        _g_.m.spinning = false
        sched.nmspinning.Add(-1)
    }

    return nil, false
}
```

---

## Work Stealing 算法

```go
// runqsteal 从其他 P 偷取 G
func runqsteal(_p_, p2 *p, stealRunNextG bool) *g {
    // 尝试获取 p2 的运行队列锁
    t := atomic.LoadAcq(&p2.runqtail)

    // 计算可偷取的数量（一半）
    n := t - atomic.Load(&p2.runqhead)
    n = n - n/2

    if n == 0 {
        // 队列为空，尝试偷取 runnext
        if stealRunNextG {
            // 尝试 CAS 偷取 runnext
            if next := p2.runnext; next != 0 && p2.status == _Prunning {
                if atomic.CasRel(&p2.runnext, next, 0) {
                    return next.ptr()
                }
            }
        }
        return nil
    }

    // 偷取 n 个 G
    return runqgrab(p2, _p_, int(n))
}

// 偷取顺序（避免热点）
var stealOrder randomOrder

type randomOrder struct {
    count    uint32
    coprimes []uint32
}

func (ord *randomOrder) reset(count uint32) {
    ord.count = count
    ord.coprimes = ord.coprimes[:0]
    for i := uint32(1); i <= count; i++ {
        if gcd(i, count) == 1 {
            ord.coprimes = append(ord.coprimes, i)
        }
    }
}

func (ord *randomOrder) randomOrder(id uint32) []uint32 {
    // 基于 id 生成伪随机顺序
    // 确保每个 P 有不同的偷取顺序，减少冲突
}
```

---

## 系统调用处理

```go
// entersyscall 进入系统调用
func entersyscall() {
    reentersyscall(getcallerpc(), getcallersp())
}

func reentersyscall(pc, sp uintptr) {
    _g_ := getg()

    // 保存状态
    save(pc, sp)
    _g_.syscallsp = sp
    _g_.syscallpc = pc

    // 释放 P
    _p_ := releasep()

    // 更新 P 状态为系统调用
    atomic.Store(&_p_.status, _Psyscall)

    // 如果调度器等待，唤醒 sysmon
    if sched.sysmonwait.Load() {
        sched.sysmonwait.Store(false)
        notewakeup(&sched.sysmonnote)
    }
}

// exitsyscall 退出系统调用
func exitsyscall() {
    _g_ := getg()
    _g_.m.locks++

    // 尝试获取原来的 P
    _p_ := _g_.m.oldp.ptr()
    if _p_ != nil && _p_.status == _Psyscall && atomic.Cas(&_p_.status, _Psyscall, _Pidle) {
        // 获取 P 成功
        if acquirep(_p_) {
            _g_.m.locks--
            return
        }
    }

    // 获取失败，需要绑定新的 P
    _g_.m.locks--

    // 检查是否有空闲 P
    if sched.pidle.Load() != 0 {
        exitsyscall0(_g_)
        return
    }

    // 没有空闲 P，放入全局队列
    mcall(exitsyscall0)
}
```

---

## Sysmon 监控线程

```go
// sysmon 监控线程
func sysmon() {
    // 初始化
    lock(&sched.lock)
    sched.nmsys++
    checkdead()
    unlock(&sched.lock)

    lasttrace := int64(0)
    nscvg := int64(0)

    for {
        // 延迟（20μs - 10ms）
        delay := uint32(20)
        if sched.gcwaiting != 0 {
            delay = 10 * 1000 // 10ms
        }
        usleep(delay)

        // 检查死锁
        if sched.gcwaiting != 0 || sched.npidle.Load() == uint32(gomaxprocs) {
            lock(&sched.lock)
            if atomic.Load(&sched.npidle) == uint32(gomaxprocs) && sched.runqhead == sched.runqtail {
                // 可能是死锁或没有工作
                checkdead()
            }
            unlock(&sched.lock)
        }

        // 抢占长时间运行的 G
        if preemptenabled() {
            for _, _p_ := range allp {
                if _p_ == nil {
                    continue
                }

                // 检查是否需要抢占
                if _p_.status == _Prunning {
                    if gp := _p_.m.ptr().curg; gp != nil && gp.preempt {
                        // 发送抢占信号
                        preemptone(_p_)
                    }
                }
            }
        }

        // 轮询网络
        if lastpoll != 0 && lastpoll+10*1000*1000 > nanotime() {
            // 检查是否有就绪的网络连接
            gp := netpoll(0)
            if gp != nil {
                injectglist(&gp)
            }
        }

        // 强制 GC
        if t := (gcTrigger{kind: gcTriggerTime, now: nanotime()}); t.test() {
            gcStart(t)
        }

        // 释放物理内存
        if debug.scavenge > 0 && atomic.Load64(&lastscavenge)+int64(5*60*1e9) < nanotime() {
            // 5分钟没有释放，尝试释放
            mheap_.scavenge(uintptr(physPageSize), true)
        }
    }
}
```

---

## 调度器性能数据

| 指标 | 数值 | 说明 |
|------|------|------|
| G 创建时间 | ~200ns | 包含栈分配 |
| G 切换时间 | ~50-100ns | 保存/恢复寄存器 |
| 上下文切换 | ~1.5μs | 涉及 M 切换 |
| 最大 G 数量 | 无限制 | 受内存限制 |
| 典型 G 栈 | 2KB-1GB | 动态增长 |
| M 最大数量 | 10,000 | 可配置 |
| P 数量 | GOMAXPROCS | 默认 CPU 核心数 |

---

## 调试与观察

```bash
# 调度器跟踪
GODEBUG=schedtrace=1000 ./program

# 输出示例：
# SCHED 0ms: gomaxprocs=8 idleprocs=5 threads=4 spinningthreads=0
#   idlethreads=2 runqueue=0 [0 0 0 0 0 0 0 0]

# 详细调度跟踪
GODEBUG=schedtrace=1000,scheddetail=1 ./program

# 生成执行跟踪
go test -trace trace.out
go tool trace trace.out
```
