# Go运行时内部实现

> 深入Go运行时源码的调度、内存分配和系统调用机制

---

## 一、运行时启动流程

### 1.1 程序入口到main

```text
Go程序启动流程:
────────────────────────────────────────

1. 操作系统加载器
   ↓
2. _rt0_amd64_linux (汇编入口)
   ↓
3. runtime.rt0_go
   ↓
4. runtime.osinit
   ↓
5. runtime.schedinit
   ↓
6. runtime.newproc (创建main goroutine)
   ↓
7. runtime.mstart
   ↓
8. runtime.main
   ↓
9. main.main (用户代码)

代码分析:
// src/runtime/asm_amd64.s
TEXT _rt0_amd64_linux(SB),NOSPLIT,$-8
    LEAQ    argv+0(FP), SI  // argv
    MOVQ    argc-8(FP), DI  // argc
    MOVQ    $runtime·rt0_go(SB), AX
    JMP     AX

// src/runtime/proc.go
func rt0_go() {
    // 初始化stack
    // 初始化CPU特性检测
    // 调用osinit
    // 调用schedinit
    // 创建main goroutine
    // 调用mstart
}
```

### 1.2 运行时初始化

```text
schedinit函数:
────────────────────────────────────────

初始化内容:
├─ 全局变量
├─ 内存分配器
├─ GC参数
├─ P的创建
└─ 信号处理

代码分析:
func schedinit() {
    // 获取CPU核心数
    ncpu := getproccount()

    // 设置GOMAXPROCS
    sched.maxmcount = 10000

    // 调整GOMAXPROCS
    if procresize(procs) != nil {
        throw("unknown runnable goroutine during bootstrap")
    }
}
```

---

## 二、Goroutine调度详解

### 2.1 G结构体

```text
G结构体定义:
────────────────────────────────────────

type g struct {
    stack       stack    // goroutine栈
    stackguard0 uintptr  // 栈溢出检查
    stackguard1 uintptr

    _panic       *_panic // panic链表
    _defer       *_defer // defer链表

    m            *m      // 当前绑定的M
    sched        gobuf   // 调度信息

    syscallsp    uintptr // 系统调用时的SP
    syscallpc    uintptr // 系统调用时的PC

    stktopsp     uintptr // 栈顶SP
    param        unsafe.Pointer // 唤醒参数

    atomicstatus uint32  // 原子状态
    stackLock    uint32  // 栈锁
    goid         int64   // goroutine ID

    schedlink    guintptr // 调度链表
    waitsince    int64   // 等待开始时间
    waitreason   waitReason // 等待原因

    preempt      bool    // 抢占标志
    preemptStop  bool    // 抢占停止
    preemptShrink bool   // 抢占收缩

    asyncSafePoint bool

    paniconfault bool
    gcscandone   bool
    throwsplit   bool

    lockedm      muintptr // 锁定的M
    sig          uint32
    writebuf     []byte
    sigcode0     uintptr
    sigcode1     uintptr
    sigpc        uintptr

    gopc          uintptr // 创建者PC
    ancestors     *[]ancestorInfo
    startpc       uintptr // 起始PC
    racectx       uintptr
    waiting       *sudog
    cgoCtxt       []uintptr
    labels        unsafe.Pointer
    timer         *timer
}

G状态:
const (
    _Gidle = iota      // 刚分配，未初始化
    _Grunnable        // 可运行，在runq中
    _Grunning         // 正在运行
    _Gsyscall         // 系统调用中
    _Gwaiting         // 等待中
    _Gdead            // 已死亡
    _Genqueue_runnable // 全局runq
    _Gcopystack       // 栈拷贝中
    _Gpreempted       // 被抢占
    _Gscan            // GC扫描中
)
```

### 2.2 M结构体

```text
M结构体定义:
────────────────────────────────────────

type m struct {
    g0      *g      // G0，调度栈
    morebuf gobuf   // 更多栈信息
    divmod  uint32  // div/mod常量

    procid  uint64  // 进程ID
    gsignal *g      // 信号处理G

    goSigStack gsignalStack // 信号栈
    sigmask    sigset       // 信号掩码

    tls      [tlsSlots]uintptr // 线程本地存储

    mstartfn      func()     // 启动函数
    curg          *g         // 当前运行的G
    caughtsig     guintptr   // 捕获的信号G
    p             puintptr   // 绑定的P
    nextp         puintptr   // 下一个P
    oldp          puintptr   // 旧的P

    id            int64
    mallocing     int32      // 内存分配中
    throwing      throwType  // panic中
    preemptoff    string     // 抢占关闭原因
    locks         int32      // 锁计数
    dying         int32
    profilehz     int32

    spinning      bool       // 自旋中
    blocked       bool       // 阻塞中
    newSigstack   bool
    printlock     int8
    incgo         bool       // 在C代码中
    freeWait      atomic.Uint32
    fastrand      uint64
    needextram    bool
    traceback     uint8

    ncgocall      uint64     // cgo调用次数
    ncgo          int32      // cgo goroutine数

    pinner        *pinner

    libraryByPC   map[uintptr]*moduledata
}
```

### 2.3 P结构体

```text
P结构体定义:
────────────────────────────────────────

type p struct {
    id          int32
    status      uint32     // P状态
    link        puintptr   // 空闲P链表
    schedtick   uint32     // 调度计数
    syscalltick uint32     // 系统调用计数
    sysmontick  sysmontick // sysmon监控

    m           muintptr   // 绑定的M
    mcache      *mcache    // 内存分配缓存

    pcache      pageCache  // 页缓存
    raceprocctx uintptr

    deferpool    []*_defer // defer对象池
    deferpoolbuf [32]*_defer

    goidcache    uint64     // G ID缓存
    goidcacheend uint64

    runqhead uint32         // 本地runq头
    runqtail uint32         // 本地runq尾
    runq     [256]guintptr  // 本地runq

    runnext guintptr        // 下一个运行的G

    gFree struct {
        gList
        n int32
    }

    sudogcache []*sudog
    sudogbuf   [128]*sudog

    mspancache struct {
        buf [128]*mspan
    }

    pCache pageCache

    timer0When uint64
    timerModifiedEarliest uint64

    gcAssistTime         int64
    gcFractionalMarkTime int64

    gcMarkWorkerMode gcMarkWorkerMode
    gcMarkWorkerStartTime int64

    gcw gcWork

    wbBuf wbBuf

    runSafePointFn uint32

    statsSeq atomic.Uint32

    timersLock mutex
    timers []*timer

    numTimers atomic.Uint32
    deletedTimers atomic.Uint32
    timerRaceCtx uintptr

    maxStackScanDelta int64

    scannedStackSize uint64
    scannedStacks    uint64

    preempt bool
}
```

---

## 三、调度器核心算法

### 3.1 调度循环

```text
schedule函数:
────────────────────────────────────────

func schedule() {
    _g_ := getg()

    var gp *g
    var inheritTime bool

    // 调试相关
    if goyield.active && goyield.m == _g_.m {
        // ...
    }

    top:
    // 检查GC标记
    if sched.gcwaiting != 0 {
        gcstopm()
        goto top
    }

    // 检查P是否被抢占
    if pp.runSafePointFn != 0 {
        runSafePointFn()
    }

    // 获取待运行G
    if gp == nil {
        // 每61次从全局runq获取
        if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
            lock(&sched.lock)
            gp = globrunqget(_g_.m.p.ptr(), 1)
            unlock(&sched.lock)
        }
    }

    // 从本地runq获取
    if gp == nil {
        gp, inheritTime = runqget(_g_.m.p.ptr())
    }

    // 全局runq
    if gp == nil {
        gp, inheritTime = findrunnable()
    }

    // 执行G
    execute(gp, inheritTime)
}
```

### 3.2 工作窃取

```text
findrunnable函数:
────────────────────────────────────────

func findrunnable() (gp *g, inheritTime bool) {
    _g_ := getg()

    // 本地P
    _p_ := _g_.m.p.ptr()

    // 1. 检查runnext
    if gp, inheritTime = runqget(_p_); gp != nil {
        return gp, inheritTime
    }

    // 2. 检查全局runq
    if sched.runqsize != 0 {
        lock(&sched.lock)
        gp = globrunqget(_p_, 0)
        unlock(&sched.lock)
        if gp != nil {
            return gp, false
        }
    }

    // 3. 网络轮询
    if netpollinited() && atomic.Load(&netpollWaiters) > 0 &&
       atomic.Xchg64(&sched.lastpoll, 0) != 0 {
        list := netpoll(0)
        lock(&sched.lock)
        gp = globrunqget(_p_, 0)
        unlock(&sched.lock)
        injectglist(&list)
        if gp != nil {
            return gp, false
        }
    }

    // 4. 从其他P窃取
    procs := uint32(gomaxprocs)
    for i := uint32(0); i < procs; i++ {
        p := allp[_p_.id+i%procs]
        if p != _p_ && !idlepMask.read(p.id) {
            if gp := runqsteal(_p_, p, &stealRunNextG); gp != nil {
                return gp, false
            }
        }
    }

    // 5. 再次检查全局runq
    // ...

    // 6. 网络轮询（阻塞）
    // ...
}

工作窃取算法:
func runqsteal(_p_, p2 *p, stealRunNextG *bool) *g {
    // 从p2的runq尾部窃取
    t := p2.runqtail
    n := t - p2.runqhead

    if n == 0 {
        // 尝试窃取runnext
        if *stealRunNextG {
            if gp := p2.runnext.ptr(); gp != nil {
                // ...
            }
        }
        return nil
    }

    // 窃取一半
    n = n - n/2

    // 将偷来的G放入本地runq
    for i := uint32(0); i < n; i++ {
        g := p2.runq[(p2.runqhead+i)%uint32(len(p2.runq))]
        p2.runq[(_p_.runqtail+i)%uint32(len(_p_.runq))] = g
    }

    return _p_.runq[_p_.runqtail%n].ptr()
}
```

### 3.3 系统调用处理

```text
系统调用调度:
────────────────────────────────────────

当G执行系统调用时:
1. 保存G状态
2. G进入_Gsyscall状态
3. M释放P (P变为_Pidle)
4. M进入系统调用
5. 系统调用返回
6. M尝试重新获取P
7. 如果P被占用，G放入全局runq

entersyscall函数:
func entersyscall() {
    reentersyscall(getcallerpc(), getcallersp())
}

func reentersyscall(pc, sp uintptr) {
    _g_ := getg()

    // 保存现场
    save(pc, sp)
    _g_.syscallsp = sp
    _g_.syscallpc = pc

    // 改变状态
    casgstatus(_g_, _Grunning, _Gsyscall)

    // 释放P
    _p_.m = 0
    _g_.m.oldp.set(_p_)
    _g_.m.p = 0
    atomic.Store(&_p_.status, _Psyscall)

    // 调度锁
    if sched.locks != 0 {
        return
    }

    // 防止M被抢占
    _g_.m.locks++
}

exitsyscall函数:
func exitsyscall() {
    _g_ := getg()
    _g_.m.locks--

    // 尝试获取原来的P
    if exitsyscallfast() {
        // 快速路径成功
        return
    }

    // 慢路径: 获取新P或休眠
    mcall(exitsyscall0)
}

func exitsyscallfast() bool {
    _g_ := getg()
    oldp := _g_.m.oldp.ptr()

    // 尝试获取旧P
    if oldp != nil && oldp.status == _Psyscall &&
       atomic.Cas(&oldp.status, _Psyscall, _Pidle) {
        // 绑定P
        wirep(oldp)
        return true
    }

    return false
}
```

---

## 四、内存分配器

### 4.1 内存分配层级

```text
Go内存分配架构:
────────────────────────────────────────

层级结构:
├─ TCMalloc风格分配器
├─ mheap: 全局堆
├─ mcentral: 中心缓存 (按size class)
├─ mcache: P本地缓存
└─ mspan: 内存管理基本单元

Size Class:
├─ 67个size class (8字节 ~ 32KB)
├─ 每个size class有对应的mspan
└─ 大对象 (>32KB) 直接分配

代码分析:
type mheap struct {
    lock      mutex
    free      mTreap      // 空闲span树
    scav      mTreap      // 已回收span树
    allspans  []*mspan    // 所有mspan

    sweepgen  uint32      // 清扫代数
    sweepdone uint32      // 清扫完成

    central   [numSpanClasses]struct {
        mcentral mcentral
    }

    spanalloc             fixalloc
    cachealloc            fixalloc
    treapalloc            fixalloc

    arenas                [1 << arenaL1Bits]*[1 << arenaL2Bits]*heapArena

    heapArenaAlloc        linearAlloc
    heapArenaAllocHuge    linearAlloc

    arenaHints            *arenaHint
}

type mcentral struct {
    spanclass spanClass
    partial   [2]mSpanList // 部分空闲span
    full      [2]mSpanList // 满span
}

type mcache struct {
    tiny             uintptr
    tinyoffset       uintptr
    tinyAllocs       uintptr
    alloc            [numSpanClasses]*mspan
    stackcache       [stackNumClasses]stackfreelist
    flushGen         uint32
}
```

### 4.2 分配流程

```text
内存分配流程:
────────────────────────────────────────

小对象 (< 16字节):
1. 从mcache.tiny分配
2. 无锁，极快

小对象 (16字节 ~ 32KB):
1. 计算size class
2. 从mcache.alloc获取mspan
3. 如有空闲，直接分配
4. 如mcache为空，从mcentral获取
5. 如mcentral为空，从mheap分配mspan
6. 如mheap不足，向OS申请

大对象 (> 32KB):
1. 直接计算页数
2. 从mheap分配
3. 或向OS申请

分配函数:
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    // 获取当前P的mcache
    c := gomcache()

    // 小对象
    if size <= maxSmallSize {
        if noscan && size < maxTinySize {
            // tiny分配器
            off := c.tinyoffset
            if off+size <= maxTinySize && c.tiny != 0 {
                // 快速路径
                x = unsafe.Pointer(c.tiny + off)
                c.tinyoffset = off + size
                c.tinyAllocs++
                mp.mallocing = 0
                releasem(mp)
                return x
            }
            // 使用新tiny块
        }

        // mcache分配
        var sizeclass uint8
        if size <= smallSizeMax-8 {
            sizeclass = size_to_class8[(size+smallSizeDiv-1)/smallSizeDiv]
        } else {
            sizeclass = size_to_class128[(size-smallSizeMax+largeSizeDiv-1)/largeSizeDiv]
        }

        spc := makeSpanClass(sizeclass, noscan)
        span := c.alloc[spc]
        v := nextFreeFast(span)
        if v == 0 {
            v, span, shouldhelpgc = c.nextFree(spc)
        }
        x = unsafe.Pointer(v)
    } else {
        // 大对象
        var s *mspan
        shouldhelpgc = true
        systemstack(func() {
            s = largeAlloc(size, needzero, noscan)
        })
        s.freeindex = 1
        s.allocCount = 1
        x = unsafe.Pointer(s.base())
    }

    return x
}
```

---

## 五、垃圾回收器

### 5.1 GC算法演进

```text
Go GC发展历程:
────────────────────────────────────────

Go 1.0: 标记-清除，STW
Go 1.3: 并行标记
Go 1.5: 三色并发标记
Go 1.6: 1ms以下STW
Go 1.8: 亚毫秒STW
Go 1.12: 清扫器优化
Go 1.14: 页分配器优化
Go 1.19: 软内存限制
Go 1.26: Green Tea GC (延迟优化)

三色标记:
├─ 白色: 未访问
├─ 灰色: 访问中
└─ 黑色: 已访问，子节点已处理

写屏障:
确保并发标记期间对象图一致性
```

### 5.2 GC触发与控制

```text
GC触发条件:
────────────────────────────────────────

1. 内存分配达到阈值
   当前堆大小 * GOGC/100

2. 手动触发
   runtime.GC()

3. 系统监控触发
   超过2分钟未GC

代码分析:
func gcStart(trigger gcTrigger) {
    // 检查触发条件
    if !trigger.test() {
        return
    }

    // 开始GC循环
    for {
        // 切换GC阶段
        switch gcphase {
        case _GCoff:
            // 开始标记
            gcMarkStartup()
        case _GCmark:
            // 完成标记，开始清扫
            gcMarkDone()
        case _GCmarktermination:
            // GC结束
            gcSweep()
        }
    }
}

GC调优参数:
├─ GOGC: 目标GC频率 (默认100)
├─ GOMEMLIMIT: 软内存限制
└─ runtime.SetGCPercent()

代码示例:
// 调整GC频率
func gcTuning() {
    // 更少的GC，更多内存使用
    debug.SetGCPercent(200)

    // 设置内存限制 (Go 1.19+)
    debug.SetMemoryLimit(10 << 30)  // 10GB

    // 强制GC
    runtime.GC()

    // 释放内存给OS
    debug.FreeOSMemory()
}

// GC统计
func gcStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("HeapSys: %d MB\n", m.HeapSys/1024/1024)
    fmt.Printf("NumGC: %d\n", m.NumGC)
    fmt.Printf("PauseNs: %d ns\n", m.PauseNs[(m.NumGC+255)%256])
    fmt.Printf("GCCPUFraction: %f\n", m.GCCPUFraction)
}
```

---

## 六、系统调用与网络轮询

### 6.1 网络轮询器

```text
网络轮询架构:
────────────────────────────────────────

epoll/kqueue/IOCP封装
全局轮询器 + P本地轮询

结构:
type pollDesc struct {
    link *pollDesc
    lock mutex
    fd   uintptr

    rg   uintptr
    wg   uintptr

    user uint32
    rseq uintptr
    wseq uintptr

    closing bool
}

轮询流程:
1. goroutine执行网络IO
2. 添加到轮询器
3. G阻塞，M切换
4. 网络事件到达
5. 唤醒对应G
6. G重新调度执行

代码示例:
// 内部使用
func netpollinit()
func netpollopen(fd uintptr, pd *pollDesc) int32
func netpollclose(fd uintptr) int32
func netpoll(delta int64) gList
```

### 6.2 定时器实现

```text
Go定时器:
────────────────────────────────────────

层级实现:
├─ 64个timeBuckets
├─ 每个bucket一个堆
└─ 最小堆管理

代码分析:
type timer struct {
    pp       uintptr
    when     int64
    period   int64
    f        func(any, uintptr)
    arg      any
    seq      uintptr
    nextwhen int64
    status   uint32
}

定时器优化 (Go 1.14+):
├─ 每个P独立管理定时器
├─ 减少锁竞争
└─ 更好的扩展性
```

---

*本章深入剖析了Go运行时内部实现，涵盖调度器、内存分配、垃圾回收、系统调用等核心机制。*
