# Go sync 包内部实现 (Go sync Package Internals)

> **分类**: 语言设计
> **标签**: #sync #mutex #rwmutex #waitgroup #source-code
> **参考**: Go 1.24 src/sync, src/sync/atomic, Linux futex

---

## sync.Mutex 内部实现

### 状态定义

```go
// src/sync/mutex.go
type Mutex struct {
    state int32  // 状态：0=unlocked, 1=locked, 其他=contended
    sema  uint32 // 信号量，用于阻塞/唤醒
}

const (
    mutexLocked      = 1 << iota // 1: 锁被持有
    mutexWoken                   // 2: 有 goroutine 被唤醒
    mutexStarving                // 4: 锁处于饥饿模式
    mutexWaiterShift = iota      // 等待者计数位移
)

// 最大等待者数（29位）
const mutexMaxWaiters = 1<<29 - 1
```

### Lock 实现

```go
func (m *Mutex) Lock() {
    // 快速路径：无竞争时直接获取锁
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        if race.Enabled {
            race.Acquire(unsafe.Pointer(m))
        }
        return
    }
    // 慢速路径：有竞争
    m.lockSlow()
}

func (m *Mutex) lockSlow() {
    var waitStartTime int64
    starving := false // 当前 goroutine 是否处于饥饿状态
    awoke := false    // 当前 goroutine 是否被唤醒
    iter := 0         // 自旋迭代次数
    old := m.state    // 保存旧状态

    for {
        // 1. 自旋阶段：尝试在不阻塞的情况下获取锁
        // 条件：锁被持有但未处于饥饿模式，且可以自旋
        if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
            // 尝试设置 mutexWoken 标志，通知 Unlock 不需要唤醒其他 waiter
            if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
                atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
                awoke = true
            }
            runtime_doSpin()  // 自旋
            iter++
            old = m.state
            continue
        }

        // 2. 准备阻塞
        new := old
        // 不要设置 starving 标志，我们还不确定是否饥饿
        if old&mutexStarving == 0 {
            new |= mutexLocked  // 尝试获取锁
        }
        if old&(mutexLocked|mutexStarving) != 0 {
            new += 1 << mutexWaiterShift  // 增加等待者计数
        }
        if starving && old&mutexLocked != 0 {
            new |= mutexStarving  // 设置饥饿模式
        }
        if awoke {
            // 被唤醒的 goroutine 必须是 mutexWoken 状态
            if new&mutexWoken == 0 {
                throw("sync: inconsistent mutex state")
            }
            new &^= mutexWoken  // 清除 woken 标志
        }

        // CAS 更新状态
        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            // 成功获取锁
            if old&(mutexLocked|mutexStarving) == 0 {
                break  // 快速路径成功
            }

            // 需要排队等待
            queueLifo := waitStartTime != 0
            if waitStartTime == 0 {
                waitStartTime = runtime_nanotime()
            }
            // 阻塞在信号量上
            runtime_SemacquireMutex(&m.sema, queueLifo, 1)

            // 被唤醒后检查是否饥饿
            starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
            old = m.state

            // 饥饿模式下直接获取锁
            if old&mutexStarving != 0 {
                if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
                    throw("sync: inconsistent mutex state")
                }
                // 获取锁并减少等待者计数
                delta := int32(mutexLocked - 1<<mutexWaiterShift)
                if !starving || old>>mutexWaiterShift == 1 {
                    delta -= mutexStarving  // 清除饥饿模式
                }
                atomic.AddInt32(&m.state, delta)
                break
            }

            awoke = true
            iter = 0
        } else {
            old = m.state
        }
    }

    if race.Enabled {
        race.Acquire(unsafe.Pointer(m))
    }
}
```

### Unlock 实现

```go
func (m *Mutex) Unlock() {
    if race.Enabled {
        race.Release(unsafe.Pointer(m))
    }

    // 快速路径：直接解锁
    new := atomic.AddInt32(&m.state, -mutexLocked)
    if new != 0 {
        // 慢速路径：有等待者或处于特殊状态
        m.unlockSlow(new)
    }
}

func (m *Mutex) unlockSlow(new int32) {
    if (new+mutexLocked)&mutexLocked == 0 {
        throw("sync: unlock of unlocked mutex")
    }

    if new&mutexStarving == 0 {
        // 正常模式
        old := new
        for {
            // 如果没有等待者，或已经有 goroutine 被唤醒/获取锁，直接返回
            if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0 {
                return
            }
            // 唤醒一个等待者
            new = (old - 1<<mutexWaiterShift) | mutexWoken
            if atomic.CompareAndSwapInt32(&m.state, old, new) {
                runtime_Semrelease(&m.sema, false, 1)
                return
            }
            old = m.state
        }
    } else {
        // 饥饿模式：直接唤醒下一个等待者
        runtime_Semrelease(&m.sema, true, 1)
    }
}
```

---

## sync.RWMutex 实现

### 结构设计

```go
// src/sync/rwmutex.go
type RWMutex struct {
    w           Mutex        // 写锁互斥锁
    writerSem   uint32       // 写者信号量
    readerSem   uint32       // 读者信号量
    readerCount atomic.Int32 // 读者计数（包括等待的和活跃的）
    readerWait  atomic.Int32 // 等待的读者数量
}

const rwmutexMaxReaders = 1 << 30
```

### RLock 实现

```go
func (rw *RWMutex) RLock() {
    if race.Enabled {
        race.Disable()
    }

    // 增加读者计数
    if rw.readerCount.Add(1) < 0 {
        // 有写者在等待，需要等待
        runtime_SemacquireMutex(&rw.readerSem, false, 0)
    }

    if race.Enabled {
        race.Enable()
        race.Acquire(unsafe.Pointer(&rw.readerSem))
    }
}
```

### Lock 实现

```go
func (rw *RWMutex) Lock() {
    if race.Enabled {
        race.Disable()
    }

    // 1. 获取写锁互斥锁
    rw.w.Lock()

    // 2. 标记有写者在等待（减少 readerCount 一个大数）
    // 这会使得新的 RLock 阻塞
    r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders

    // 3. 如果有活跃的读者，等待它们完成
    if r != 0 && rw.readerWait.Add(r) != 0 {
        runtime_SemacquireMutex(&rw.writerSem, false, 0)
    }

    if race.Enabled {
        race.Enable()
        race.Acquire(unsafe.Pointer(&rw.writerSem))
    }
}
```

### RUnlock 实现

```go
func (rw *RWMutex) RUnlock() {
    if race.Enabled {
        race.ReleaseMerge(unsafe.Pointer(&rw.writerSem))
    }

    // 减少读者计数
    if rw.readerCount.Add(-1) < 0 {
        // 有写者在等待
        rw.rUnlockSlow(0)
    }
}

func (rw *RWMutex) rUnlockSlow(r int32) {
    // 减少等待的读者计数
    if rw.readerWait.Add(-1) == 0 {
        // 最后一个读者，唤醒写者
        runtime_Semrelease(&rw.writerSem, false, 1)
    }
}
```

---

## sync.WaitGroup 实现

### 结构定义

```go
// src/sync/waitgroup.go
type WaitGroup struct {
    noCopy noCopy

    // 高 32 位：计数器
    // 低 32 位：等待者数量
    state atomic.Uint64
    sema  uint32
}
```

### Add/Done/Wait 实现

```go
func (wg *WaitGroup) Add(delta int) {
    if race.Enabled {
        if delta < 0 {
            race.ReleaseMerge(unsafe.Pointer(wg))
        }
        race.Disable()
        defer race.Enable()
    }

    // 获取当前状态
    state := wg.state.Add(uint64(delta) << 32)
    v := int32(state >> 32) // 计数器
    w := uint32(state)       // 等待者数量

    if v < 0 {
        panic("sync: negative WaitGroup counter")
    }

    if v > 0 || w == 0 {
        // 还有活跃任务，或无等待者
        return
    }

    // v == 0 且 w > 0：所有任务完成，唤醒所有等待者
    if wg.state.CompareAndSwap(state, 0) {
        for ; w != 0; w-- {
            runtime_Semrelease(&wg.sema, false, 0)
        }
    }
}

func (wg *WaitGroup) Done() {
    wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
    if race.Enabled {
        race.Disable()
    }

    for {
        state := wg.state.Load()
        v := int32(state >> 32)
        if v == 0 {
            // 计数器为 0，直接返回
            if race.Enabled {
                race.Enable()
                race.Acquire(unsafe.Pointer(wg))
            }
            return
        }

        // 增加等待者计数
        if wg.state.CompareAndSwap(state, state+1) {
            // 阻塞等待
            runtime_Semacquire(&wg.sema)
            if wg.state.Load() != 0 {
                throw("sync: WaitGroup is reused before previous Wait has returned")
            }
            if race.Enabled {
                race.Enable()
                race.Acquire(unsafe.Pointer(wg))
            }
            return
        }
    }
}
```

---

## sync.Pool 实现

### 结构设计

```go
// src/sync/pool.go
type Pool struct {
    noCopy noCopy

    local     unsafe.Pointer // [P]poolLocal
    localSize uintptr        // local 数组大小

    victim     unsafe.Pointer // [P]poolLocal (上一次 GC 的数据)
    victimSize uintptr

    New func() any
}

type poolLocal struct {
    poolLocalInternal

    // 防止 false sharing
    pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}

type poolLocalInternal struct {
    private any       // 只能由当前 P 使用
    shared  poolChain // 本地共享队列
}
```

### Get/Put 实现

```go
func (p *Pool) Get() any {
    if race.Enabled {
        race.Disable()
    }

    l, pid := p.pin()
    x := l.private
    l.private = nil

    if x == nil {
        // 尝试从本地共享队列获取
        x, _ = l.shared.popHead()
        if x == nil {
            // 从其他 P 偷取
            x = p.getSlow(pid)
        }
    }

    runtime_procUnpin()
    if race.Enabled {
        race.Enable()
        if x != nil {
            race.Acquire(poolRaceAddr(x))
        }
    }

    if x == nil && p.New != nil {
        x = p.New()
    }
    return x
}

func (p *Pool) getSlow(pid int) any {
    // 尝试从其他 P 的共享队列偷取
    size := runtime_LoadAcqUintptr(&p.localSize)
    locals := p.local

    // 尝试从 victim 获取
    for i := 0; i < int(size); i++ {
        l := indexLocal(locals, (pid+i+1)%int(size))
        if x, _ := l.shared.popTail(); x != nil {
            return x
        }
    }

    // 尝试从 victim 获取
    size = runtime_LoadAcqUintptr(&p.victimSize)
    if uintptr(pid) >= size {
        return nil
    }
    locals = p.victim
    l := indexLocal(locals, pid)
    if x := l.private; x != nil {
        l.private = nil
        return x
    }
    for i := 0; i < int(size); i++ {
        l := indexLocal(locals, (pid+i)%int(size))
        if x, _ := l.shared.popTail(); x != nil {
            return x
        }
    }

    return nil
}

func (p *Pool) Put(x any) {
    if x == nil {
        return
    }

    if race.Enabled {
        race.Disable()
    }

    l, _ := p.pin()
    if l.private == nil {
        l.private = x
    } else {
        l.shared.pushHead(x)
    }
    runtime_procUnpin()

    if race.Enabled {
        race.Enable()
    }
}
```

---

## sync.Once 实现

### 结构定义

```go
// src/sync/once.go
type Once struct {
    done uint32
    m    Mutex
}
```

### Do 实现

```go
func (o *Once) Do(f func()) {
    // 快速路径：已经执行过
    if atomic.LoadUint32(&o.done) == 1 {
        return
    }
    // 慢速路径
    o.doSlow(f)
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()

    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```

---

## sync.Cond 实现

### 结构定义

```go
// src/sync/cond.go
type Cond struct {
    noCopy noCopy

    L Locker  // 关联的锁

    notify  notifyList
    checker copyChecker
}
```

### Wait/Signal/Broadcast 实现

```go
func (c *Cond) Wait() {
    c.checker.check()

    // 增加到等待列表
    t := runtime_notifyListAdd(&c.notify)

    // 释放锁
    c.L.Unlock()

    // 等待信号
    runtime_notifyListWait(&c.notify, t)

    // 重新获取锁
    c.L.Lock()
}

func (c *Cond) Signal() {
    c.checker.check()
    runtime_notifyListNotifyOne(&c.notify)
}

func (c *Cond) Broadcast() {
    c.checker.check()
    runtime_notifyListNotifyAll(&c.notify)
}
```

---

## 性能对比

```
┌─────────────────────────────────────────────────────────────────┐
│                    sync 包性能对比 (ns/op)                       │
├─────────────────────────────────────────────────────────────────┤
│  操作              │  Mutex    │ RWMutex   │ Atomic    │ Channel │
├─────────────────────────────────────────────────────────────────┤
│  Lock/Unlock       │    15     │    15     │     -     │    -    │
│  RLock/RUnlock     │     -     │     8     │     -     │    -    │
│  AddInt64          │     -     │     -     │     7     │    -    │
│  Send/Recv         │     -     │     -     │     -     │   100   │
│  WaitGroup.Add/Done│     -     │     -     │    25     │    -    │
└─────────────────────────────────────────────────────────────────┘

注：数据来自 go1.24 linux/amd64，仅供参考
```

---

## 最佳实践

```go
// 1. 不要拷贝 sync 类型
func badPractice() {
    var m sync.Mutex
    m2 := m  // 错误：拷贝了 Mutex
    _ = m2
}

// 2. 使用 defer Unlock
func goodPractice() {
    var m sync.Mutex
    m.Lock()
    defer m.Unlock()
    // 业务逻辑
}

// 3. RWMutex 偏好读场景
func rwBestPractice() {
    var rw sync.RWMutex
    data := make(map[string]int)

    // 读（并发安全，可并行）
    rw.RLock()
    v := data["key"]
    rw.RUnlock()

    // 写（互斥）
    rw.Lock()
    data["key"] = v + 1
    rw.Unlock()
}

// 4. Pool 用于减少 GC 压力
var bufPool = sync.Pool{
    New: func() any {
        return make([]byte, 1024)
    },
}

func usePool() {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf)
    // 使用 buf
}
```
