# LD-021: Go Sync 包深度剖析 (Go Sync Package Deep Dive)

> **维度**: Language Design
> **级别**: S (19+ KB)
> **标签**: #sync #concurrency #mutex #atomic #waitgroup #pool #once
> **权威来源**:
>
> - [sync Package](https://github.com/golang/go/tree/master/src/sync) - Go Authors
> - [Go Memory Model](https://go.dev/ref/mem) - Go Authors
> - [The Go Programming Language](https://www.gopl.io/) - Donovan & Kernighan

---

## 1. Sync 包架构

### 1.1 组件概览

```
┌─────────────────────────────────────────────────────────────┐
│                      sync/                                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  互斥锁                                                      │
│  ├── Mutex        - 基本互斥锁                               │
│  └── RWMutex      - 读写锁                                   │
│                                                              │
│  同步原语                                                    │
│  ├── WaitGroup    - 等待组                                   │
│  ├── Once         - 一次性执行                               │
│  ├── Pool         - 对象池                                   │
│  └── Cond         - 条件变量                                 │
│                                                              │
│  原子操作 (sync/atomic)                                      │
│  ├── Add/Sub      - 增减操作                                 │
│  ├── CompareAndSwap - CAS                                   │
│  └── Load/Store   - 读写操作                                 │
│                                                              │
│  Map (sync.Map)                                              │
│  └── 并发安全的 map                                          │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心原则

**原则 1: 零值可用**

```go
// 所有 sync 类型都可以直接使用零值
var mu sync.Mutex        // 未锁定
var wg sync.WaitGroup    // 计数为 0
var once sync.Once       // 未执行
var pool sync.Pool       // 空池
```

**原则 2: 禁止复制**

```go
// 编译时检测复制
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Mutex struct {
    _ noCopy  // 嵌入检测
    // ...
}
```

---

## 2. Mutex 实现

### 2.1 状态分析

```go
// src/sync/mutex.go

type Mutex struct {
    state int32
    sema  uint32
}

const (
    mutexLocked = 1 << iota      // 1 - 锁定状态
    mutexWoken                   // 2 - 被唤醒
    mutexStarving                // 4 - 饥饿模式
    mutexWaiterShift = iota      // 3 - waiter 计数位移
)

// 正常模式 vs 饥饿模式
//
// 正常模式：waiter 按 FIFO 排队，但新到达的 goroutine 有更高概率获取锁
//          （因为它们已经在运行）
//
// 饥饿模式：waiter 严格按 FIFO 获取锁，新到达的 goroutine 不尝试获取
//          防止长尾延迟
```

### 2.2 Lock 实现

```go
func (m *Mutex) Lock() {
    // 快速路径：无竞争时直接获取
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        return
    }

    // 慢速路径：需要排队
    m.lockSlow()
}

func (m *Mutex) lockSlow() {
    var waitStartTime int64
    starving := false
    awoke := false
    iter := 0
    old := m.state

    for {
        // 饥饿模式或锁已持有，自旋
        if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
            // 主动自旋，不加入队列
            if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
                atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
                awoke = true
            }
            runtime_doSpin()  // runtime 实现的自旋
            iter++
            old = m.state
            continue
        }

        // 计算新状态
        new := old
        if old&mutexStarving == 0 {
            new |= mutexLocked  // 尝试获取锁
        }
        if old&(mutexLocked|mutexStarving) != 0 {
            new += 1 << mutexWaiterShift  // 增加 waiter 计数
        }
        if starving && old&mutexLocked != 0 {
            new |= mutexStarving  // 进入饥饿模式
        }
        if awoke {
            new &^= mutexWoken  // 清除唤醒标志
        }

        if atomic.CompareAndSwapInt32(&m.state, old, new) {
            if old&(mutexLocked|mutexStarving) == 0 {
                break  // 获取成功
            }

            // 排队等待
            queueStart := false
            if waitStartTime == 0 {
                queueStart = true
                waitStartTime = runtime_nanotime()
            }

            runtime_SemacquireMutex(&m.sema, queueStart, starving)

            // 检查是否进入饥饿模式
            starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
            old = m.state

            if old&mutexStarving != 0 {
                // 饥饿模式：锁直接传递给我们
                // 修正等待者计数
                delta := int32(mutexLocked - 1<<mutexWaiterShift)
                if !starving || old>>mutexWaiterShift == 1 {
                    delta -= mutexStarving  // 退出饥饿模式
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
}
```

### 2.3 Unlock 实现

```go
func (m *Mutex) Unlock() {
    // 快速路径：直接解锁
    new := atomic.AddInt32(&m.state, -mutexLocked)
    if new != 0 {
        m.unlockSlow(new)
    }
}

func (m *Mutex) unlockSlow(new int32) {
    if (new+mutexLocked)&mutexLocked == 0 {
        throw("sync: unlock of unlocked mutex")
    }

    if new&mutexStarving == 0 {
        old := new
        for {
            // 没有等待者，或者已经有人被唤醒或获取了锁
            if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
                return
            }
            // 唤醒一个等待者
            new = (old - 1<<mutexWaiterShift) | mutexWoken
            if atomic.CompareAndSwapInt32(&m.state, old, new) {
                runtime_Semrelease(&m.sema, false)
                return
            }
            old = m.state
        }
    } else {
        // 饥饿模式：直接传递给下一个等待者
        runtime_Semrelease(&m.sema, true)
    }
}
```

---

## 3. RWMutex 实现

### 3.1 设计原理

```go
// src/sync/rwmutex.go

type RWMutex struct {
    w           Mutex        // 写锁
    writerSem   uint32       // 写等待信号
    readerSem   uint32       // 读等待信号
    readerCount atomic.Int32 // 读锁计数
    readerWait  atomic.Int32 // 等待写锁的读锁数量
}

// readerCount 编码：
// 正数：当前持有读锁的 goroutine 数
// 负数：有写者在等待 (-rwmutexMaxReaders + n 读者)
const rwmutexMaxReaders = 1 << 30
```

### 3.2 RLock/RUnlock

```go
func (rw *RWMutex) RLock() {
    // 增加读者计数
    if rw.readerCount.Add(1) < 0 {
        // 有写者在等待，需要排队
        runtime_Semacquire(&rw.readerSem)
    }
}

func (rw *RWMutex) RUnlock() {
    // 减少读者计数
    if r := rw.readerCount.Add(-1); r < 0 {
        rw.rUnlockSlow(r)
    }
}

func (rw *RWMutex) rUnlockSlow(r int32) {
    // 减少等待写的读者计数
    if rw.readerWait.Add(-1) == 0 {
        // 最后一个读者，唤醒写者
        runtime_Semrelease(&rw.writerSem, false)
    }
}
```

### 3.3 Lock/Unlock

```go
func (rw *RWMutex) Lock() {
    // 获取写互斥锁
    rw.w.Lock()

    // 标记有写者在等待
    // 阻止新读者
    r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders

    // 等待现有读者完成
    if r != 0 && rw.readerWait.Add(r) != 0 {
        runtime_Semacquire(&rw.writerSem)
    }
}

func (rw *RWMutex) Unlock() {
    // 标记无写者等待
    r := rw.readerCount.Add(rwmutexMaxReaders)

    // 唤醒等待的读者
    for i := 0; i < int(r); i++ {
        runtime_Semrelease(&rw.readerSem, false)
    }

    // 释放写互斥锁
    rw.w.Unlock()
}
```

---

## 4. WaitGroup 实现

### 4.1 数据结构

```go
// src/sync/waitgroup.go

type WaitGroup struct {
    noCopy noCopy

    state atomic.Uint64 // 高32位：计数器，低32位：等待者数
    sema  uint32        // 信号量
}

// state 布局:
// +-----------------+-------------+
// |  counter (32)   | waiters (32)|
// +-----------------+-------------+
// 64                 32            0
```

### 4.2 Add/Done/Wait

```go
func (wg *WaitGroup) Add(delta int) {
    // 获取当前状态
    state := wg.state.Add(uint64(delta) << 32)

    v := int32(state >> 32)     // 计数器
    w := uint32(state)          // 等待者数

    if v < 0 {
        panic("sync: negative WaitGroup counter")
    }

    // v == delta 且 w > 0 意味着之前计数器为 0，现在有 Add 操作
    // 这是竞态条件
    if v > 0 || w == 0 {
        return
    }

    // 计数器为 0，唤醒所有等待者
    if wg.state.Load() != state {
        panic("sync: WaitGroup misuse: Add called concurrently with Wait")
    }

    // 重置状态并唤醒
    wg.state.Store(0)
    for ; w != 0; w-- {
        runtime_Semrelease(&wg.sema, false)
    }
}

func (wg *WaitGroup) Done() {
    wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
    for {
        state := wg.state.Load()
        v := int32(state >> 32)
        if v == 0 {
            // 计数器为 0，无需等待
            return
        }

        // 增加等待者计数
        if wg.state.CompareAndSwap(state, state+1) {
            // 等待信号
            runtime_Semacquire(&wg.sema)
            return
        }
    }
}
```

---

## 5. Once 实现

### 5.1 演进历史

```go
// Go 1.x 实现 (慢路径使用互斥锁)
type Once struct {
    m    Mutex
    done uint32
}

// Go 1.18+ 实现 (使用原子操作优化)
type Once struct {
    done atomic.Uint32
    m    Mutex
}
```

### 5.2 当前实现

```go
// src/sync/once.go

type Once struct {
    done atomic.Uint32
    m    Mutex
}

func (o *Once) Do(f func()) {
    // 快速路径：检查是否已执行
    if o.done.Load() == 1 {
        return
    }

    // 慢速路径
    o.doSlow(f)
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()

    // 双重检查
    if o.done.Load() == 0 {
        defer o.done.Store(1)
        f()
    }
}
```

---

## 6. Pool 实现

### 6.1 核心设计

```go
// src/sync/pool.go

type Pool struct {
    noCopy noCopy

    local     unsafe.Pointer // [P]poolLocal 数组
    localSize uintptr        // 数组大小

    victim     unsafe.Pointer // 上一轮 GC 的 local
    victimSize uintptr        // victim 大小

    New func() any
}

type poolLocal struct {
    poolLocalInternal

    // 防止 false sharing
    pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}

type poolLocalInternal struct {
    private any      // 私有对象（仅当前 P 可访问）
    shared  poolChain // 共享列表（LFIFO）
}
```

### 6.2 Get/Put 实现

```go
func (p *Pool) Get() any {
    // 禁止抢占，确保获取当前 P
    l, pid := p.pin()
    x := l.private
    l.private = nil

    if x == nil {
        // 尝试从共享列表获取
        x, _ = l.shared.popHead()
        if x == nil {
            x = p.getSlow(pid)
        }
    }
    runtime_procUnpin()

    if x == nil && p.New != nil {
        x = p.New()
    }
    return x
}

func (p *Pool) getSlow(pid int) any {
    size := runtime_LoadAcquintptr(&p.localSize)
    locals := p.local

    // 尝试从其他 P 偷取
    for i := 0; i < int(size); i++ {
        l := indexLocal(locals, (pid+i+1)%int(size))
        if x, _ := l.shared.popTail(); x != nil {
            return x
        }
    }

    // 尝试从 victim 获取
    size = runtime_LoadAcquintptr(&p.victimSize)
    if uintptr(pid) >= size {
        return nil
    }
    l := indexLocal(p.victim, pid)
    if x := l.private; x != nil {
        l.private = nil
        return x
    }
    if x, _ := l.shared.popTail(); x != nil {
        return x
    }

    return nil
}

func (p *Pool) Put(x any) {
    if x == nil {
        return
    }

    l, _ := p.pin()
    if l.private == nil {
        l.private = x
    } else {
        l.shared.pushHead(x)
    }
    runtime_procUnpin()
}
```

### 6.3 GC 处理

```go
// 每个 GC 周期调用
func poolCleanup() {
    // 丢弃所有 victim
    for _, p := range oldPools {
        p.victim = nil
        p.victimSize = 0
    }

    // local 变为 victim
    for _, p := range allPools {
        p.victim = p.local
        p.victimSize = p.localSize
        p.local = nil
        p.localSize = 0
    }

    oldPools, allPools = allPools, nil
}
```

---

## 7. Map 实现

### 7.1 设计原理

```go
// src/sync/map.go

// Map 针对两种场景优化：
// 1. 写一次，读多次（entry 只会递增更新）
// 2. 多个 goroutine 读写不相交的键集

type Map struct {
    mu Mutex

    // 只读部分（原子访问）
    read atomic.Pointer[readOnly]

    // 脏数据（需要 mu 保护）
    dirty map[any]*entry

    // 统计
    misses int
}

type readOnly struct {
    m       map[any]*entry
    amended bool // dirty 包含 read 中没有的键
}

type entry struct {
    p atomic.Pointer[any] // *any 或 expunged
}

var expunged = new(any) // 标记已删除
```

### 7.2 Load/Store

```go
func (m *Map) Load(key any) (value any, ok bool) {
    read := m.loadRead()
    e, ok := read.m[key]
    if !ok && read.amended {
        // read 中没有，检查 dirty
        m.mu.Lock()
        read = m.loadRead()
        e, ok = read.m[key]
        if !ok && read.amended {
            e, ok = m.dirty[key]
            m.missLocked()
        }
        m.mu.Unlock()
    }
    if !ok {
        return nil, false
    }
    return e.load()
}

func (e *entry) load() (value any, ok bool) {
    p := e.p.Load()
    if p == nil || p == expunged {
        return nil, false
    }
    return *p, true
}

func (m *Map) Store(key, value any) {
    read := m.loadRead()
    if e, ok := read.m[key]; ok && e.tryStore(&value) {
        return
    }

    m.mu.Lock()
    read = m.loadRead()
    if e, ok := read.m[key]; ok {
        if e.unexpungeLocked() {
            // 之前是 expunged，需要写入 dirty
            m.dirty[key] = e
        }
        e.storeLocked(&value)
    } else if e, ok := m.dirty[key]; ok {
        e.storeLocked(&value)
    } else {
        if !read.amended {
            // 首次写入新键，需要 dirty
            m.dirtyLocked()
            m.read.Store(&readOnly{m: read.m, amended: true})
        }
        m.dirty[key] = newEntry(value)
    }
    m.mu.Unlock()
}
```

---

## 8. 原子操作

### 8.1 内存序

```go
// sync/atomic 保证顺序一致性

var flag int32
var data string

func write() {
    data = "hello"
    atomic.StoreInt32(&flag, 1)
}

func read() {
    if atomic.LoadInt32(&flag) == 1 {
        // 保证看到 data = "hello"
        println(data)
    }
}
```

### 8.2 实现原理

```go
// src/sync/atomic/doc.go

// 所有操作编译为硬件原子指令
// - x86: LOCK 前缀
// - ARM: LDREX/STREX
// - RISC-V: AMO 指令

// 原子类型 (Go 1.19+)
type Int32 struct { v int32 }
type Int64 struct { v int64 }
type Uint32 struct { v uint32 }
type Uint64 struct { v uint64 }
type Uintptr struct { v uintptr }
type Pointer[T any] struct { v *T }
```

---

## 9. 性能特征与基准测试

### 9.1 基准测试

```go
func BenchmarkMutex(b *testing.B) {
    var mu sync.Mutex
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mu.Lock()
            mu.Unlock()
        }
    })
}

func BenchmarkRWMutex(b *testing.B) {
    var mu sync.RWMutex
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mu.RLock()
            mu.RUnlock()
        }
    })
}

func BenchmarkPool(b *testing.B) {
    pool := sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)
        },
    }
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            buf := pool.Get()
            pool.Put(buf)
        }
    })
}

// 典型结果 (Go 1.21, 8核)
// BenchmarkMutex-8      100000000    15 ns/op    0 allocs/op
// BenchmarkRWMutex-8     80000000    18 ns/op    0 allocs/op  (读)
// BenchmarkPool-8       200000000     6 ns/op    0 allocs/op
```

### 9.2 性能对比

| 原语 | 无竞争延迟 | 竞争延迟 | 内存开销 |
|------|-----------|---------|---------|
| sync.Mutex | ~10ns | 100-500ns | 8 bytes |
| sync.RWMutex | ~15ns (读) | 200-800ns | 24 bytes |
| sync/atomic | ~5ns | ~5ns | 4-8 bytes |
| sync.Map | ~50ns | ~100ns | 高 |
| channel | ~100ns | ~200ns | 高 |

---

## 10. 视觉表征

### 10.1 Mutex 状态转换

```
        ┌─────────────┐
        │   Unlocked   │
        │   (state=0)  │
        └──────┬──────┘
               │ Lock (无竞争)
               │ CAS(0, 1)
               ▼
        ┌─────────────┐
        │   Locked     │◄─────────┐
        │  (state=1)   │          │
        └──────┬──────┘          │
               │                 │
      Lock     │ Unlock          │ Lock (竞争)
      (竞争)   │                 │
               ▼                 │
        ┌─────────────┐         │
        │  等待队列    │─────────┘
        │  (waiters)  │  唤醒
        └─────────────┘
```

### 10.2 RWMutex 并发模型

```
时间 ───────────────────────────────────────►

Reader 1: [RLock------------------------RUnlock]
Reader 2:    [RLock--------------RUnlock]
Reader 3:           [RLock-----------------RUnlock]
Writer 1:                  [Lock------------------Unlock]
                                  ▲
                                  │
                          写者等待所有读者
```

### 10.3 Pool 双队列设计

```
┌─────────────────────────────────────────┐
│              sync.Pool                   │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────────┐    ┌─────────────┐    │
│  │    local    │    │   victim    │    │
│  │  (当前周期)  │    │  (上一周期)  │    │
│  ├─────────────┤    ├─────────────┤    │
│  │ P0: private │    │ P0: private │    │
│  │    shared   │    │    shared   │    │
│  ├─────────────┤    ├─────────────┤    │
│  │ P1: private │    │ P1: private │    │
│  │    shared   │    │    shared   │    │
│  ├─────────────┤    ├─────────────┤    │
│  │ P2: private │    │ P2: private │    │
│  │    shared   │    │    shared   │    │
│  └─────────────┘    └─────────────┘    │
│                                         │
│  GC 时: local → victim, victim → nil    │
└─────────────────────────────────────────┘
```

---

## 11. 完整代码示例

### 11.1 并发安全的缓存

```go
package main

import (
    "sync"
    "time"
)

// 并发安全的 LRU 缓存
type LRUCache struct {
    mu       sync.RWMutex
    items    map[string]*cacheItem
    capacity int
}

type cacheItem struct {
    value      interface{}
    expiration int64
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        items:    make(map[string]*cacheItem),
        capacity: capacity,
    }
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, ok := c.items[key]
    if !ok {
        return nil, false
    }

    // 检查过期
    if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
        return nil, false
    }

    return item.value, true
}

func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    expiration := int64(0)
    if ttl > 0 {
        expiration = time.Now().Add(ttl).UnixNano()
    }

    c.items[key] = &cacheItem{
        value:      value,
        expiration: expiration,
    }
}

// 使用 sync.Map 的版本
type SyncMapCache struct {
    items sync.Map
}

func (c *SyncMapCache) Get(key string) (interface{}, bool) {
    value, ok := c.items.Load(key)
    if !ok {
        return nil, false
    }

    item := value.(*cacheItem)
    if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
        c.items.Delete(key)
        return nil, false
    }

    return item.value, true
}

func (c *SyncMapCache) Set(key string, value interface{}, ttl time.Duration) {
    expiration := int64(0)
    if ttl > 0 {
        expiration = time.Now().Add(ttl).UnixNano()
    }

    c.items.Store(key, &cacheItem{
        value:      value,
        expiration: expiration,
    })
}
```

---

**质量评级**: S (19KB)
**完成日期**: 2026-04-02
