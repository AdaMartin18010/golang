# Go内存管理与GC详解

> 深入Go内存分配器和垃圾回收器实现

---

## 一、内存分配器

### 1.1 TCMalloc架构

```text
Go内存分配器基于TCMalloc:
────────────────────────────────────────

三层架构:
├─ mcache: 每个P的本地缓存 (无锁)
├─ mcentral: 全局中心缓存 (按size class)
└─ mheap: 全局堆 (向OS申请)

对象大小分类:
├─ Tiny: < 16字节
├─ Small: 16字节 ~ 32KB
└─ Large: > 32KB

Size Class:
├─ 67个size class
├─ 8, 16, 24, 32, 48, 64, 80, 96, 112, 128...
└─ 每个size class对应固定大小的对象
```

### 1.2 小对象分配

```text
小对象分配流程:
────────────────────────────────────────

1. 计算size class
2. 从mcache获取mspan
3. mspan中有空闲对象则分配
4. mcache空则从mcentral获取
5. mcentral空则从mheap分配mspan
6. mheap空则向OS申请内存

代码分析:
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    // 获取当前P的mcache
    mp := acquirem()
    c := mp.mcache

    var x unsafe.Pointer
    noscan := typ == nil || typ.PtrBytes == 0

    if size <= maxSmallSize {
        if noscan && size < maxTinySize {
            // Tiny allocator
            off := c.tinyoffset
            if off+size <= maxTinySize && c.tiny != 0 {
                x = unsafe.Pointer(c.tiny + off)
                c.tinyoffset = off + size
                c.tinyAllocs++
                mp.mallocing = 0
                releasem(mp)
                return x
            }
            // 分配新的tiny块
        }

        // 计算size class
        var sizeclass uint8
        if size <= smallSizeMax-8 {
            sizeclass = size_to_class8[(size+smallSizeDiv-1)/smallSizeDiv]
        } else {
            sizeclass = size_to_class128[(size-smallSizeMax+largeSizeDiv-1)/largeSizeDiv]
        }

        size = uintptr(class_to_size[sizeclass])
        spc := makeSpanClass(sizeclass, noscan)
        span := c.alloc[spc]
        v := nextFreeFast(span)
        if v == 0 {
            v, span, shouldhelpgc = c.nextFree(spc)
        }
        x = unsafe.Pointer(v)
    } else {
        // 大对象分配
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

## 二、垃圾回收器

### 2.1 三色标记算法

```text
三色标记:
────────────────────────────────────────
白色: 未访问 (可能是垃圾)
灰色: 访问过，但子节点未完全访问
黑色: 已完全访问，保留

标记过程:
1. 所有对象初始为白色
2. GC Roots标记为灰色
3. 取出灰色对象，标记为黑色
4. 将其引用的白色对象标记为灰色
5. 重复3-4直到没有灰色对象
6. 白色对象即为垃圾

写屏障 (Write Barrier):
确保并发标记期间对象图一致性

代码分析:
// 标记阶段
func gcDrain(gcw *gcWork, flags gcDrainFlags) {
    gp := getg().m.curg

    for !(preemptible && gp.preempt) {
        // 从gcw获取对象
        var b uintptr
        if blocking {
            b = gcw.get()
        } else {
            b = gcw.tryGetFast()
            if b == 0 {
                b = gcw.tryGet()
            }
        }

        if b == 0 {
            // 工作窃取
            b = trySteal()
        }

        if b == 0 {
            break
        }

        // 扫描对象
        scanobject(b, gcw)
    }
}
```

### 2.2 GC触发与控制

```text
GC触发条件:
────────────────────────────────────────

1. 内存分配达到阈值
   next_gc = heap_marked * (1 + GOGC/100)

2. 手动触发: runtime.GC()

3. 系统监控触发: 超过2分钟未GC

GC调优:
├─ GOGC: 默认100，增大减少GC频率
├─ GOMEMLIMIT: 软内存限制
└─ runtime.SetGCPercent()

代码示例:
// 调整GC频率
func gcTuning() {
    // 减少GC频率 (牺牲内存)
    debug.SetGCPercent(200)

    // 设置内存限制
    debug.SetMemoryLimit(10 << 30) // 10GB

    // 手动触发GC
    runtime.GC()

    // 释放内存给OS
    debug.FreeOSMemory()
}

// GC统计
func gcStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("Alloc = %v MiB\n", bToMb(m.Alloc))
    fmt.Printf("TotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
    fmt.Printf("Sys = %v MiB\n", bToMb(m.Sys))
    fmt.Printf("NumGC = %v\n", m.NumGC)
    fmt.Printf("PauseNs = %v ms\n", m.PauseNs[(m.NumGC+255)%256]/1e6)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
```

### 2.3 GC优化技术

```text
减少GC压力:
────────────────────────────────────────

1. 减少堆分配
   - 使用对象池
   - 栈分配
   - 预分配

2. 减少指针数量
   - 值类型优于指针
   - 避免不必要的引用

3. 分代GC (Go未实现但可模拟)
   - 对象生命周期分离

代码示例:
// 对象池减少分配
var bufPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func processWithPool() {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf)

    // 使用buf...
}

// 值类型减少指针
// 不良: 大量指针
type Node struct {
    Value *int
    Next  *Node
}

// 优化: 值类型
type NodeOpt struct {
    Value int
    Next  int // 索引而非指针
}
```

---

## 三、内存泄漏排查

### 3.1 常见内存泄漏

```text
Go内存泄漏类型:
────────────────────────────────────────

1. Goroutine泄漏
   - 无退出条件的goroutine
   - 阻塞在channel

2. 全局变量累积
   - 无限制增长的map/slice

3. 未释放的资源
   - 未close的文件
   - 未cancel的context

排查工具:
├─ pprof heap
├─ pprof goroutine
├─ runtime.ReadMemStats
└─ Go 1.26 goroutine leak检测

代码示例:
// Goroutine泄漏检测
type GoroutineMonitor struct {
    baseline int
}

func (m *GoroutineMonitor) Start() {
    m.baseline = runtime.NumGoroutine()
}

func (m *GoroutineMonitor) Check() error {
    time.Sleep(100 * time.Millisecond) // 等待稳定
    current := runtime.NumGoroutine()
    if current > m.baseline+10 {
        return fmt.Errorf("goroutine leak: %d -> %d", m.baseline, current)
    }
    return nil
}

// 使用
func TestNoLeak(t *testing.T) {
    monitor := &GoroutineMonitor{}
    monitor.Start()
    defer func() {
        if err := monitor.Check(); err != nil {
            t.Error(err)
        }
    }()

    // 测试代码...
}
```

### 3.2 pprof内存分析

```text
内存分析:
────────────────────────────────────────

生成heap profile:
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

分析命令:
go tool pprof http://localhost:6060/debug/pprof/heap
go tool pprof -inuse_space heap.out
go tool pprof -alloc_space heap.out

常用命令:
├─ top: 查看占用最高的函数
├─ list FuncName: 查看函数详情
├─ web: 生成调用图
└─ svg: 生成SVG图形

代码示例:
// 程序内获取profile
func writeHeapProfile(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    runtime.GC() // 先GC，获取准确数据
    return pprof.WriteHeapProfile(f)
}
```

---

*本章深入剖析了Go内存管理和垃圾回收机制，涵盖内存分配、GC算法、内存泄漏排查等核心内容。*
