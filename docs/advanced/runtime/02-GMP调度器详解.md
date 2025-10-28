# GMP调度器详解

## 📋 目录

- [1. GMP模型概述](#1-gmp模型概述)
- [2. G - Goroutine](#2-g---goroutine)
- [3. M - Machine](#3-m---machine)
- [4. P - Processor](#4-p---processor)
- [5. 调度流程](#5-调度流程)
- [6. 抢占式调度](#6-抢占式调度)
- [7. 系统调用处理](#7-系统调用处理)
- [8. 性能优化](#8-性能优化)

---

## 1. GMP模型概述

### 为什么需要GMP

**传统线程模型的问题**:
- ❌ 线程创建开销大（~1MB栈空间）
- ❌ 线程切换开销大（内核态切换）
- ❌ 线程数量受限（操作系统限制）
- ❌ 调度由OS控制，无法优化

**Go的解决方案 - GMP模型**:
- ✅ Goroutine轻量（初始2KB栈）
- ✅ 用户态调度（无需内核切换）
- ✅ 可创建百万级goroutine
- ✅ 自定义调度策略

---

### GMP三者关系

```
     Goroutine Pool
    ┌───┬───┬───┬───┐
    │ G │ G │ G │ G │ ...  (待调度的goroutine)
    └───┴───┴───┴───┘
          ↓
    ┌─────────────┐
    │  Scheduler  │  (调度器)
    └─────────────┘
          ↓
    ┌──────────────────────────┐
    │ P → M → G (running)      │
    │ P → M → G (running)      │
    │ P → M → G (running)      │
    └──────────────────────────┘
          ↓
    Operating System Threads
```

**核心概念**:
- **G (Goroutine)**: 代表一个goroutine，包含栈、指令指针等
- **M (Machine)**: 代表一个内核线程，执行G的实体
- **P (Processor)**: 代表调度上下文，持有G的本地队列

**关键关系**:
- M必须关联P才能执行G
- P的数量 = GOMAXPROCS (默认CPU核心数)
- M的数量动态调整，通常 M数量 > P数量

---

## 2. G - Goroutine

### G的数据结构

```go
// runtime2.go
type g struct {
    // 栈信息
    stack       stack     // 栈边界 [stack.lo, stack.hi)
    stackguard0 uintptr   // 栈溢出检查
    stackguard1 uintptr   // C栈溢出检查
    
    // 调度信息
    m              *m        // 当前运行在哪个M上
    sched          gobuf     // 调度上下文（PC、SP等）
    atomicstatus   uint32    // G的状态
    schedlink      guintptr  // 调度队列链表
    
    // 抢占标志
    preempt       bool      // 抢占标志
    preemptStop   bool      // 抢占到_Gpreempted
    preemptShrink bool      // 收缩栈
    
    // 唤醒时间
    waitsince     int64     // 等待开始时间
    waitreason    waitReason // 等待原因
    
    // panic/defer
    _panic    *_panic  // panic链表
    _defer    *_defer  // defer链表
    
    // ...其他字段
}

// 调度上下文
type gobuf struct {
    sp   uintptr  // 栈指针
    pc   uintptr  // 程序计数器
    g    guintptr // goroutine指针
    ret  sys.Uintreg // 返回值
    // ...
}
```

---

### G的状态转换

```
创建
  ↓
_Gidle (刚分配)
  ↓
_Grunnable (可运行，在队列中)
  ↓
_Grunning (正在运行)
  ↓
  ├→ _Gwaiting (等待中，如channel、select)
  ├→ _Gsyscall (系统调用中)
  ├→ _Gpreempted (被抢占)
  └→ _Gdead (执行完成)
```

**状态说明**:

| 状态 | 值 | 说明 |
|------|-----|------|
| `_Gidle` | 0 | 刚分配，未初始化 |
| `_Grunnable` | 1 | 在运行队列，等待执行 |
| `_Grunning` | 2 | 正在执行 |
| `_Gsyscall` | 3 | 系统调用中 |
| `_Gwaiting` | 4 | 阻塞等待（IO、锁等） |
| `_Gdead` | 6 | 执行完成，可复用 |
| `_Gpreempted` | 9 | 被抢占，栈扫描 |

---

### G的创建与销毁

**创建Goroutine**:

```go
// proc.go
func newproc(siz int32, fn *funcval) {
    // 获取调用者信息
    argp := add(unsafe.Pointer(&fn), sys.PtrSize)
    gp := getg()
    pc := getcallerpc()
    
    // 在系统栈上创建新的G
    systemstack(func() {
        newg := newproc1(fn, argp, siz, gp, pc)
        
        // 放入P的本地队列
        _p_ := getg().m.p.ptr()
        runqput(_p_, newg, true)
        
        // 如果有空闲P且没在spinning，唤醒或创建M
        if mainStarted {
            wakep()
        }
    })
}
```

**G的复用**:

```go
// proc.go
func gfget(_p_ *p) *g {
    // 从P的本地gfree列表获取
    gp := _p_.gFree.pop()
    if gp == nil {
        // 从全局gfree列表获取
        gp = sched.gFree.pop()
    }
    return gp
}
```

---

## 3. M - Machine

### M的数据结构

```go
// runtime2.go
type m struct {
    g0      *g       // 用于调度的特殊g（有更大的栈）
    morebuf gobuf    // 栈扩展用
    curg    *g       // 当前运行的g
    
    // 关联的P
    p             puintptr // 当前关联的P
    nextp         puintptr // 下一个要关联的P
    oldp          puintptr // 执行syscall前的P
    
    // M的状态
    spinning      bool     // 是否在窃取工作
    blocked       bool     // 是否阻塞在note上
    
    // 系统调用
    syscalltick   uint32   // 系统调用计数
    
    // 链表
    schedlink     muintptr // 链表
    alllink       *m       // allm链表
    
    // 线程信息
    thread        uintptr  // 线程句柄
    
    // ...其他字段
}
```

---

### M的创建与销毁

**创建M**:

```go
// proc.go
func newm(fn func(), _p_ *p, id int64) {
    // 分配m结构
    mp := allocm(_p_, fn, id)
    mp.nextp.set(_p_)
    
    // 创建操作系统线程
    newm1(mp)
}

func newm1(mp *m) {
    // 创建线程
    execLock.rlock()
    newosproc(mp)
    execLock.runlock()
}

// os_linux.go (Linux平台)
func newosproc(mp *m) {
    // 调用clone系统调用创建线程
    ret := clone(cloneFlags, stk, unsafe.Pointer(mp), unsafe.Pointer(mp.g0), unsafe.Pointer(funcPC(mstart)))
}
```

**M的数量控制**:

```go
const (
    maxMCount = 10000  // 最大M数量
)

var (
    sched struct {
        midle        muintptr  // 空闲M链表
        nmidle       int32     // 空闲M数量
        nmidlelocked int32     // 锁定的空闲M数量
        mnext        int64     // M的ID分配
        maxmcount    int32     // 最大M数量
    }
)
```

---

## 4. P - Processor

### P的数据结构

```go
// runtime2.go
type p struct {
    m           muintptr   // 关联的M
    
    // 本地运行队列
    runqhead    uint32     // 队列头
    runqtail    uint32     // 队列尾
    runq        [256]guintptr // 本地队列，循环数组
    runnext     guintptr   // 下一个运行的G（优先级最高）
    
    // 状态
    status      uint32     // P的状态
    link        puintptr   // P链表
    schedtick   uint32     // 调度计数
    syscalltick uint32     // 系统调用计数
    
    // mcache for allocation
    mcache      *mcache    // 内存分配缓存
    
    // defer pool
    deferpool    [5][]*_defer
    deferpoolbuf [5][32]*_defer
    
    // sudoG pool
    sudogcache []*sudog
    sudogbuf   [128]*sudog
    
    // timer heap
    timers      []*timer
    numTimers   uint32
    
    // ...其他字段
}
```

---

### P的状态转换

```
_Pidle (空闲)
  ↓
_Prunning (运行中)
  ↓
  ├→ _Psyscall (系统调用)
  └→ _Pgcstop (GC停止)
```

**P的数量调整**:

```go
// proc.go
func procresize(nprocs int32) *p {
    old := gomaxprocs
    
    // 创建或销毁P
    for i := old; i < nprocs; i++ {
        pp := allp[i]
        if pp == nil {
            pp = new(p)
            pp.init(i)
        }
    }
    
    // 释放多余的P
    for i := nprocs; i < old; i++ {
        p := allp[i]
        // 将P的runq移到全局队列
        for !runqempty(p) {
            gp := runqget(p)
            globrunqput(gp)
        }
        // 释放P
        p.destroy()
    }
    
    return allp[0]
}
```

---

## 5. 调度流程

### schedule函数

```go
// proc.go
func schedule() {
    _g_ := getg()
    _g_.m.locks++
    
top:
    pp := _g_.m.p.ptr()
    
    // 1. 检查GC
    if sched.gcwaiting != 0 {
        gcstopm()
        goto top
    }
    
    var gp *g
    var inheritTime bool
    
    // 2. 每61次从全局队列获取（防止全局队列饿死）
    if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
        lock(&sched.lock)
        gp = globrunqget(_g_.m.p.ptr(), 1)
        unlock(&sched.lock)
    }
    
    // 3. 从P的本地队列获取
    if gp == nil {
        gp, inheritTime = runqget(_g_.m.p.ptr())
    }
    
    // 4. findrunnable（阻塞获取）
    if gp == nil {
        gp, inheritTime = findrunnable()
    }
    
    // 5. 执行goroutine
    execute(gp, inheritTime)
}
```

---

### findrunnable详解

```go
// proc.go
func findrunnable() (gp *g, inheritTime bool) {
    _g_ := getg()
    _p_ := _g_.m.p.ptr()

top:
    // 1. 从本地队列获取
    if gp, inheritTime := runqget(_p_); gp != nil {
        return gp, inheritTime
    }

    // 2. 从全局队列获取
    if sched.runqsize != 0 {
        lock(&sched.lock)
        gp := globrunqget(_p_, 0)
        unlock(&sched.lock)
        if gp != nil {
            return gp, false
        }
    }

    // 3. 检查netpoll
    if netpollinited() && atomic.Load(&netpollWaiters) > 0 && atomic.Load64(&sched.lastpoll) != 0 {
        list := netpoll(0)
        if !list.empty() {
            gp := list.pop()
            injectglist(&list)
            return gp, false
        }
    }

    // 4. Work Stealing - 从其他P窃取
    procs := uint32(gomaxprocs)
    if procs > 1 {
        for i := 0; i < 4; i++ {
            for enum := stealOrder.start(fastrand()); !enum.done(); enum.next() {
                p2 := allp[enum.position()]
                if _p_ == p2 {
                    continue
                }
                
                // 从p2窃取一半的G
                gp := runqsteal(_p_, p2, stealRunNextG)
                if gp != nil {
                    return gp, false
                }
            }
        }
    }

    // 5. 再次检查全局队列
    if sched.runqsize != 0 {
        gp := globrunqget(_p_, 0)
        if gp != nil {
            return gp, false
        }
    }

    // 6. 进入休眠前的最后检查
    stopm()
    goto top
}
```

---

### Work Stealing

**窃取算法**:

```go
// proc.go
func runqsteal(_p_, p2 *p, stealRunNextG bool) *g {
    t := _p_.runqtail
    n := runqgrab(p2, &_p_.runq, t, stealRunNextG)
    if n == 0 {
        return nil
    }
    n--
    gp := _p_.runq[(t+n)%uint32(len(_p_.runq))].ptr()
    _p_.runqtail = t + n
    return gp
}

func runqgrab(_p_ *p, batch *[256]guintptr, batchHead uint32, stealRunNextG bool) uint32 {
    for {
        h := atomic.LoadAcq(&_p_.runqhead)
        t := atomic.LoadAcq(&_p_.runqtail)
        n := t - h
        n = n - n/2  // 窃取一半
        
        if n == 0 {
            // 尝试窃取runnext
            if stealRunNextG {
                if gp := _p_.runnext.ptr(); gp != nil {
                    // CAS操作窃取
                }
            }
            return 0
        }
        
        // 批量窃取
        if atomic.CasRel(&_p_.runqhead, h, h+n) {
            return n
        }
    }
}
```

**特点**:
- 窃取一半的G
- 随机选择受害P，减少冲突
- 优先窃取runnext（最后放入的G）

---

## 6. 抢占式调度

### 基于协作的抢占

**Go 1.14之前**:

```go
// 在函数调用时检查抢占
func morestack() {
    if getg().stackguard0 == stackPreempt {
        // 被标记为抢占
        gopreempt_m()
    }
}
```

**问题**: 无法抢占无函数调用的死循环

```go
// 无法被抢占
func loop() {
    for {
        // 无函数调用
    }
}
```

---

### 基于信号的抢占 (Go 1.14+)

```go
// signal_unix.go
func sighandler(sig uint32, info *siginfo, ctxt unsafe.Pointer, gp *g) {
    if sig == _SIGURG {
        // 异步抢占信号
        doSigPreempt(gp, c)
    }
}

func doSigPreempt(gp *g, ctxt *sigctxt) {
    // 检查是否可以抢占
    if wantAsyncPreempt(gp) {
        if ok, newpc := isAsyncSafePoint(gp, ctxt.sigpc(), ctxt.sigsp(), ctxt.siglr()); ok {
            // 注入抢占调用
            ctxt.pushCall(funcPC(asyncPreempt), newpc)
        }
    }
}
```

**抢占时机**:

1. **sysmon检测** (每10ms):
```go
// proc.go
func sysmon() {
    for {
        // 检查运行超过10ms的P
        now := nanotime()
        if pd := &allp[i].sysmontick; now-pd.schedwhen > 10*1000*1000 {
            preemptone(allp[i])
        }
        
        usleep(10 * 1000) // 休眠10ms
    }
}
```

2. **GC触发抢占**:
```go
func preemptall() bool {
    for _, _p_ := range allp {
        if _p_.status != _Prunning {
            continue
        }
        preemptone(_p_)
    }
}
```

---

## 7. 系统调用处理

### 进入系统调用

```go
// proc.go
func reentersyscall(pc, sp uintptr) {
    _g_ := getg()
    
    // 保存调用者信息
    _g_.syscallsp = sp
    _g_.syscallpc = pc
    
    // 解除M和P的绑定
    _g_.m.oldp.set(_g_.m.p.ptr())
    _g_.m.p = 0
    
    // 将P状态设置为_Psyscall
    atomic.Store(&_g_.m.oldp.ptr().status, _Psyscall)
    
    // sysmon可以接管这个P
}
```

---

### 退出系统调用

```go
// proc.go
func exitsyscall() {
    _g_ := getg()
    
    // 尝试重新关联原来的P
    oldp := _g_.m.oldp.ptr()
    if oldp != nil && oldp.status == _Psyscall && cas(&oldp.status, _Psyscall, _Prunning) {
        // 成功重新关联
        _g_.m.p.set(oldp)
        return
    }
    
    // 无法关联原P，尝试获取空闲P
    mcall(exitsyscall0)
}

func exitsyscall0(gp *g) {
    _g_ := getg()
    _p_ := pidleget()
    if _p_ == nil {
        // 没有空闲P，将G放入全局队列
        globrunqput(gp)
        
        // M进入休眠
        stopm()
    } else {
        // 关联P
        _g_.m.p.set(_p_)
        execute(gp, false)
    }
}
```

---

### 阻塞系统调用监控

```go
// proc.go (sysmon)
func retake(now int64) uint32 {
    n := 0
    for i := 0; i < len(allp); i++ {
        _p_ := allp[i]
        pd := &_p_.sysmontick
        s := _p_.status
        
        if s == _Psyscall {
            // 系统调用超过10ms
            if runqempty(_p_) && atomic.Load(&sched.nmspinning)+atomic.Load(&sched.npidle) > 0 {
                // 队列为空且有空闲资源，不抢占
                continue
            }
            
            t := int64(_p_.syscalltick)
            if int64(pd.syscalltick) != t {
                pd.syscalltick = uint32(t)
                pd.syscallwhen = now
                continue
            }
            
            // 系统调用时间过长，抢占P
            if runqempty(_p_) && atomic.Load(&sched.nmspinning)+atomic.Load(&sched.npidle) > 0 && pd.syscallwhen+10*1000*1000 > now {
                continue
            }
            
            // 将P从M上剥离
            if atomic.Cas(&_p_.status, s, _Pidle) {
                n++
                _p_.syscalltick++
                handoffp(_p_)
            }
        }
    }
    return uint32(n)
}
```

---

## 8. 性能优化

### 优化1: 减少Goroutine切换

```go
// ❌ 频繁切换
func badPattern() {
    for i := 0; i < 1000; i++ {
        go func() {
            // 极短任务
        }()
    }
}

// ✅ 批量处理
func goodPattern() {
    numWorkers := runtime.GOMAXPROCS(0)
    jobs := make(chan int, 1000)
    
    for i := 0; i < numWorkers; i++ {
        go func() {
            for j := range jobs {
                // 处理任务
            }
        }()
    }
    
    for i := 0; i < 1000; i++ {
        jobs <- i
    }
}
```

---

### 优化2: GOMAXPROCS调优

```go
// 获取系统CPU核心数
numCPU := runtime.NumCPU()

// CPU密集型：GOMAXPROCS = CPU核心数
runtime.GOMAXPROCS(numCPU)

// IO密集型：可以适当增加
runtime.GOMAXPROCS(numCPU * 2)
```

**性能对比**:

| 场景 | GOMAXPROCS | QPS | CPU% |
|------|------------|-----|------|
| CPU密集 | 1 | 1000 | 100% |
| CPU密集 | 4 | 3800 | 100% |
| CPU密集 | 8 | 7200 | 100% |
| IO密集 | 1 | 2000 | 30% |
| IO密集 | 8 | 15000 | 80% |

---

### 优化3: 减少抢占开销

```go
// ❌ 长时间计算，频繁抢占
func heavyCompute() {
    for i := 0; i < 1000000000; i++ {
        // 计算密集
    }
}

// ✅ 定期让出CPU
func heavyComputeOptimized() {
    for i := 0; i < 1000000000; i++ {
        if i%10000000 == 0 {
            runtime.Gosched() // 主动让出
        }
        // 计算密集
    }
}
```

---

### 调度器性能指标

```go
func printSchedStats() {
    var stats runtime.SchedStats
    runtime.ReadSchedStats(&stats)
    
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    fmt.Printf("OS Threads: %d\n", stats.NumThreads)
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
}
```

---

## 🔗 相关资源

- [Go Runtime架构总览](./01-Go-Runtime架构总览.md)
- [内存分配器原理](./03-内存分配器原理.md)
- [垃圾回收器详解](./04-垃圾回收器详解.md)
- [并发编程](../../fundamentals/language/02-并发编程/)

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3  
**文档类型**: GMP调度器深度解析 ✨

