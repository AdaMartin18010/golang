# LD-011: Go 垃圾回收算法与内存管理 (Go GC Algorithm & Memory Management)

> **维度**: Language Design
> **级别**: S (20+ KB)
> **标签**: #garbage-collection #tricolor #concurrent-gc #write-barrier #memory-management #tri-color
> **权威来源**:
>
> - [On-the-fly Garbage Collection](https://dl.acm.org/doi/10.1145/359580.359587) - Dijkstra et al. (1978)
> - [Go GC Guide](https://go.dev/doc/gc-guide) - Go Authors
> - [Concurrent Garbage Collection](https://www.cs.cmu.edu/~fp/courses/15411-f14/lectures/23-gc.pdf) - CMU 15-411
> - [Go 1.5 GC](https://go.dev/s/go15gc) - Rick Hudson (2015)
> - [The Garbage Collection Handbook](https://gchandbook.org/) - Jones et al. (2012)

---

## 1. 形式化基础

### 1.1 垃圾回收理论

**定义 1.1 (可达性)**
对象 $o$ 从根集合 $R$ 可达，当且仅当存在引用链：

$$\text{reachable}(o) \Leftrightarrow \exists r \in R: r \to^* o$$

**定义 1.2 (垃圾)**
垃圾是不可达对象的集合：

$$\text{Garbage} = \{ o \in \text{Heap} \mid \neg \text{reachable}(o) \}$$

**定义 1.3 (根集合)**

$$R = \text{Globals} \cup \text{Stacks} \cup \text{Registers}$$

**定理 1.1 (GC 安全性)**
垃圾回收器不会回收可达对象：

$$\forall o: \text{collected}(o) \Rightarrow o \in \text{Garbage}$$

*证明*：

1. GC 从 $R$ 开始标记所有可达对象
2. 只有未被标记的对象才会被回收
3. 因此回收对象必定不可达

### 1.2 三色标记-清除

**定义 1.4 (三色抽象)**

$$\text{Color} = \{ \text{White}, \text{Grey}, \text{Black} \}$$

- **White**: 待检查对象（候选垃圾）
- **Grey**: 已发现但未扫描子引用的对象
- **Black**: 已完全扫描的对象

**定义 1.5 (三色不变式)**

$$\forall b \in \text{Black}, w \in \text{White}: \neg(b \to w)$$

黑色对象不直接引用白色对象。

**定理 1.2 (三色不变式保持)**
若维持三色不变式，当 Grey 集合为空时，White 对象即为垃圾。

*证明*：

1. 所有从根可达对象要么是 Black（已扫描），要么是 Grey（待扫描）
2. Grey 为空意味着所有可达对象都是 Black
3. 根据不变式，Black 对象不引用 White 对象
4. 因此 White 对象不可达，是垃圾

---

## 2. 并发 GC 算法

### 2.1 并发标记算法

**算法 2.1 (并发三色标记)**

```
1. 初始化:
   - 所有对象标记为 White
   - 根对象标记为 Grey，加入 worklist

2. 标记阶段 (与 Mutator 并发):
   while worklist 非空:
       obj = worklist.pop()
       obj.color = Black
       for ref in obj.refs:
           if ref.color == White:
               ref.color = Grey
               worklist.push(ref)

3. 重扫描 (STW):
   - 处理标记期间被修改的对象

4. 清除阶段:
   - 回收所有 White 对象
```

### 2.2 写屏障 (Write Barrier)

**定义 2.1 (指针更新)**
当 mutator 修改指针时：

$$\text{slot} \to \text{ptr}$$

**写屏障规则 (Dijkstra)**

```
writeBarrier(slot, ptr):
    if gcPhase == MARK && slot is BLACK && ptr is WHITE:
        shade(ptr)    // 标记 ptr 为 Grey
    *slot = ptr
```

**定义 2.2 (混合写屏障)**
Go 使用 Yuasa 风格的删除屏障 + Dijkstra 插入屏障：

```
writeBarrier(slot, ptr):
    // 删除屏障: 标记旧值
    if gcPhase == MARK:
        shade(*slot)

    // 插入屏障: 标记新值
    if gcPhase == MARK && ptr != nil:
        shade(ptr)

    *slot = ptr
```

**定理 2.1 (混合屏障正确性)**
混合屏障保证不会漏标对象，无需重新扫描栈。

### 2.3 GC 阶段详解

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go GC Cycle                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Off                                                                        │
│   │                                                                         │
│   ▼ 堆达到阈值                                                              │
│  ┌─────────────────┐                                                        │
│  │ Sweep Termination│  STW                                                  │
│  │ (完成上一次清扫) │  • 等待所有清扫完成                                    │
│  └────────┬────────┘                                                        │
│           │                                                                 │
│           ▼                                                                 │
│  ┌─────────────────┐                                                        │
│  │    Mark Start   │  STW                                                  │
│  │   (开始标记)     │  • 启用写屏障                                          │
│  │                 │  • 扫描根对象 (goroutine 栈)                            │
│  └────────┬────────┘  • 加入 Grey 集合                                      │
│           │                                                                 │
│           ▼                                                                 │
│  ┌─────────────────┐                                                        │
│  │    Mark         │  并发                                                  │
│  │   (并发标记)     │  • 从 Grey 集合扫描对象                                │
│  │                 │  • 标记子对象为 Grey                                    │
│  │                 │  • 当前对象为 Black                                     │
│  │                 │  • Mutator 并发运行，写屏障保护                          │
│  └────────┬────────┘                                                        │
│           │                                                                 │
│           ▼ Grey 为空                                                        │
│  ┌─────────────────┐                                                        │
│  │  Mark Termination│ STW                                                   │
│  │   (标记终止)     │ • 停止所有 Mutator                                     │
│  │                 │ • 刷新写屏障缓冲区                                      │
│  │                 │ • 统计 GC 数据                                          │
│  └────────┬────────┘                                                        │
│           │                                                                 │
│           ▼                                                                 │
│  ┌─────────────────┐                                                        │
│  │     Sweep       │  并发                                                  │
│  │    (并发清扫)    │  • 逐 span 清扫                                         │
│  │                 │  • 回收 White 对象到空闲列表                             │
│  │                 │  • Mutator 可分配新内存                                   │
│  └─────────────────┘                                                        │
│                                                                              │
│  循环: 下一轮 GC 触发                                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. 内存分配器

### 3.1 分级分配

**定义 3.1 (Span 结构)**

```go
type mspan struct {
    next      *mspan        // 链表下一个
    prev      *mspan        // 链表上一个
    startAddr uintptr       // 起始地址
    npages    uintptr       // 页数
    freeindex uintptr       // 下一个空闲索引
    allocCache uint64       // 分配缓存位图
    allocBits  *gcBits      // 分配位图
    gcmarkBits *gcBits      // GC 标记位图
}
```

**定义 3.2 (Size Class)**

| Class | Size | Max Waste | Objects/Page |
|-------|------|-----------|--------------|
| 1 | 8B | 87.50% | 1024 |
| 2 | 16B | 46.67% | 512 |
| 3 | 32B | 46.67% | 256 |
| ... | ... | ... | ... |
| 67 | 32KB | 0% | 1 |

**定义 3.3 (分配路径)**

```
Tiny (≤16B):
    • 合并小对象
    • 无锁快速分配

Small (16B - 32KB):
    • 按 size class 分配
    • mcache → mcentral → mheap

Large (>32KB):
    • 直接从 mheap 分配
    • 整数页数
```

### 3.2 分配器架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Go Memory Allocator Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Per-P Cache (mcache)                                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Tiny allocator (≤16B)                                              │    │
│  │  ├── tiny  uintptr    // 当前 tiny 块指针                           │    │
│  │  ├── tinyoffset uintptr // 当前偏移                                 │    │
│  │  └── tinyAllocs uintptr // 计数                                     │    │
│  │                                                                     │    │
│  │  Span caches (67 size classes)                                      │    │
│  │  ├── alloc [numSpanClasses]*mspan  // 本地 span 缓存               │    │
│  │  └── stackcache [numStackClasses]*mspan // 栈缓存                   │    │
│  │                                                                     │    │
│  │  分配流程 (Small objects):                                          │    │
│  │  1. 计算 size class                                                 │    │
│  │  2. 从 mcache.alloc[class] 分配 (无锁)                               │    │
│  │  3. 若空，从 mcentral 补充                                           │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                  │                                           │
│                                  ▼ (mcache 不足)                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Central Cache (mcentral)                                           │    │
│  │  ├── lock mutex                                                     │    │
│  │  ├── nonempty mSpanList  // 非空 span 链表                          │    │
│  │  └── empty mSpanList     // 空 span 链表                            │    │
│  │                                                                     │    │
│  │  补充流程:                                                           │    │
│  │  1. 锁定 mcentral                                                   │    │
│  │  2. 从 nonempty 获取 span                                           │    │
│  │  3. 若空，从 mheap 分配                                              │    │
│  │  4. 解锁，返回 span                                                  │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                  │                                           │
│                                  ▼ (mcentral 不足)                           │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Heap (mheap)                                                       │    │
│  │  ├── lock mutex                                                     │    │
│  │  ├── free [maxMHeapList]mSpanList  // 空闲 span 列表                 │    │
│  │  ├── freelarge mSpanList           // 大 span 列表                   │    │
│  │  ├── arenas [1 << arenaL1Bits]*[1 << arenaL2Bits]*heapArena          │    │
│  │  └── curArena struct {                                             │    │
│  │      base, end uintptr  // 当前 arena 范围                          │    │
│  │  }                                                                  │    │
│  │                                                                     │    │
│  │  分配流程:                                                           │    │
│  │  1. 从 free 列表查找合适 span                                        │    │
│  │  2. 若无，从操作系统分配新 arena                                     │    │
│  │  3. 分割 span 返回给 mcentral                                        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. GC 触发与调优

### 4.1 触发条件

**定义 4.1 (GC 触发堆大小)**

$$H_T = \frac{GOGC}{100} \times H_L$$

其中 $H_L$ 是上一次 GC 后的存活堆大小。

**定义 4.2 (目标 CPU 使用率)**

$$\text{Target CPU} = \frac{GOGC}{GOGC + 100}$$

默认 GOGC=100 时，GC 目标使用 50% CPU。

### 4.2 内存限制 (Go 1.19+)

**定义 4.3 (软内存限制)**

```go
// runtime/debug.SetMemoryLimit
// 当内存接近限制时，GC 会更激进
// 超过限制时，GC 会更频繁
```

---

## 5. 多元表征

### 5.1 GC 算法演进对比

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Go GC Algorithm Evolution                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Go 1.0: 停止世界标记-清除                                                   │
│  ├── 简单实现                                                                │
│  ├── 长停顿 (100ms+)                                                         │
│  └── 无并发                                                                  │
│                                                                              │
│  Go 1.3: 并行清扫                                                            │
│  ├── 清扫阶段与 Mutator 并发                                                 │
│  └── 减少停顿时间                                                            │
│                                                                              │
│  Go 1.5: 并发三色标记                                                        │
│  ├── 写屏障引入                                                              │
│  ├── 目标停顿 < 10ms                                                         │
│  └── 大部分工作与 Mutator 并发                                               │
│                                                                              │
│  Go 1.8: 亚毫秒 STW                                                          │
│  ├── 优化根对象扫描                                                          │
│  ├── 目标停顿 < 100μs                                                        │
│  └── 实际通常 < 50μs                                                         │
│                                                                              │
│  Go 1.14+: 异步抢占                                                          │
│  ├── 解决长时间循环导致的 GC 延迟                                            │
│  └── 信号驱动的异步抢占                                                      │
│                                                                              │
│  Go 1.19+: 软内存限制                                                        │
│  ├── SetMemoryLimit API                                                      │
│  └── 更精细的内存控制                                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 GC 决策树

```
遇到 GC 相关问题?
│
├── 高停顿时间?
│   ├── 检查 GC trace: GODEBUG=gctrace=1
│   ├── 分析 heap profile
│   ├── 检查 finalizers (延迟回收)
│   └── 考虑减少堆分配
│
├── 高 CPU 使用?
│   ├── 增加 GOGC (牺牲内存换 CPU)
│   │   └── GOGC=200 或更高
│   ├── 减少不必要的 allocation
│   └── 使用 sync.Pool 重用对象
│
├── 内存使用过高?
│   ├── 检查存活对象: runtime.ReadMemStats
│   ├── 使用 pprof heap 分析
│   ├── 检查 goroutine 泄漏
│   └── 减少 GOGC (增加 GC 频率)
│
└── 需要更细粒度控制?
    ├── 使用 runtime.SetGCPercent
    ├── 使用 runtime/debug.SetMemoryLimit
    └── 使用 runtime.GC() 手动触发 (不推荐常规使用)
```

### 5.3 内存分配决策图

```
需要分配内存?
│
├── 对象大小?
│   ├── ≤ 16 bytes → Tiny allocator
│   │   └── 合并分配，减少碎片
│   │
│   ├── 16 bytes - 32 KB → Size class allocation
│   │   ├── 计算 size class
│   │   ├── mcache 分配 (无锁)
│   │   ├── mcentral 补充 (有锁)
│   │   └── mheap 分配新 span
│   │
│   └── > 32 KB → Large allocation
│       └── 直接从 mheap 分配
│
├── 是否需要零值?
│   ├── make/new 自动零值
│   └── 手动分配需 memset
│
├── 性能关键路径?
│   ├── 是
│   │   ├── 使用 sync.Pool
│   │   ├── 栈分配 (逃逸分析)
│   │   └── 预分配 slice/map 容量
│   └── 否
│
└── 对象生命周期?
    ├── 短期 → 栈分配 (优先)
    ├── 中期 → 堆分配
    └── 长期 → 考虑对象池
```

### 5.4 GC 可视化

```
三色标记过程示例:

初始状态:                    标记根对象:
  Roots → A → B → C           Roots(黑) → A(灰) → B(白) → C(白)
          ↓                             ↓
          D                             D(白)

所有节点白色

标记 A:                      标记 B:
  Roots(黑) → A(黑) → B(灰) → C(白)    Roots(黑) → A(黑) → B(黑) → C(灰)
              ↓                                    ↓
              D(灰)                                D(灰)

标记 C:                      完成 (Grey 为空):
  Roots(黑) → A(黑) → B(黑) → C(黑)    Roots(黑) → A(黑) → B(黑) → C(黑)
              ↓                                    ↓
              D(黑)                                D(黑)

白色集合为空，无垃圾可回收
```

---

## 6. 代码示例与基准测试

### 6.1 GC 控制与监控

```go
package gc

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "time"
)

// GC 调优示例
func GCTuningExample() {
    // 设置 GC 目标百分比
    old := debug.SetGCPercent(100) // 默认 100
    fmt.Printf("Old GOGC: %d\n", old)

    // 设置内存限制 (Go 1.19+)
    debug.SetMemoryLimit(1 << 30) // 1GB

    // 手动触发 GC
    runtime.GC()

    // 读取 GC 统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("HeapSys: %d MB\n", m.HeapSys/1024/1024)
    fmt.Printf("NumGC: %d\n", m.NumGC)
    fmt.Printf("Last GC: %v\n", time.Unix(0, int64(m.LastGC)))

    // GC 暂停时间
    fmt.Printf("PauseNs (last 5): ")
    for i := 0; i < 5 && i < m.NumGC; i++ {
        idx := (m.NumGC - uint32(i) + 255) % 256
        fmt.Printf("%d ", m.PauseNs[idx]/1000)
    }
    fmt.Println("μs")
}

// 对象池模式
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func ProcessWithPool() {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)

    // 使用 buf...
    for i := range buf {
        buf[i] = byte(i % 256)
    }
}

// 预分配减少 GC 压力
func PreallocateExample(n int) [][]int {
    // 预分配外层 slice
    result := make([][]int, 0, n)

    for i := 0; i < n; i++ {
        // 预分配内层 slice
        row := make([]int, 0, 100)
        for j := 0; j < 100; j++ {
            row = append(row, j)
        }
        result = append(result, row)
    }

    return result
}
```

### 6.2 性能基准测试

```go
package gc_test

import (
    "runtime"
    "testing"
)

// 基准测试: 分配开销
func BenchmarkAllocationSmall(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = make([]byte, 64)
    }
}

func BenchmarkAllocationLarge(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = make([]byte, 64*1024)
    }
}

// 基准测试: 对象池
func BenchmarkWithPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gc.ProcessWithPool()
    }
}

func BenchmarkWithoutPool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        buf := make([]byte, 1024)
        for j := range buf {
            buf[j] = byte(j % 256)
        }
    }
}

// 基准测试: 预分配 vs 动态增长
func BenchmarkPreallocatedSlice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        s := make([]int, 0, 1000)
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

func BenchmarkDynamicSlice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var s []int
        for j := 0; j < 1000; j++ {
            s = append(s, j)
        }
    }
}

// 基准测试: GC 影响
func BenchmarkGCImpact(b *testing.B) {
    // 创建大量临时对象
    for i := 0; i < b.N; i++ {
        data := make([]*int, 1000)
        for j := range data {
            v := j
            data[j] = &v
        }
        // data 在此丢弃，需要 GC 回收
    }
}

// 基准测试: 逃逸分析
func BenchmarkStackAllocation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 可能分配在栈上
        x := struct{ a, b int }{i, i}
        _ = x.a + x.b
    }
}

func BenchmarkHeapAllocation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // 强制分配在堆上
        x := &struct{ a, b int }{i, i}
        _ = x.a + x.b
    }
}

// 测试 GC 触发
func TestGCTrigger(t *testing.T) {
    var m1, m2 runtime.MemStats

    runtime.ReadMemStats(&m1)

    // 分配大量内存
    data := make([][]byte, 1000)
    for i := range data {
        data[i] = make([]byte, 1024*1024) // 1MB each
    }

    runtime.ReadMemStats(&m2)

    t.Logf("HeapAlloc before: %d MB", m1.HeapAlloc/1024/1024)
    t.Logf("HeapAlloc after: %d MB", m2.HeapAlloc/1024/1024)
    t.Logf("NumGC: %d -> %d", m1.NumGC, m2.NumGC)
}
```

---

## 7. 关系网络

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Go GC Context                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  GC 算法家族                                                                 │
│  ├── Mark-Sweep (McCarthy, 1960)                                            │
│  ├── Reference Counting (Collins, 1960)                                     │
│  ├── Copying (Fenichel, 1969)                                               │
│  ├── Generational (Lieberman, 1983)                                         │
│  ├── Tri-color (Dijkstra, 1978)                                             │
│  ├── Concurrent GC (Boehm, 1993)                                            │
│  └── Region-based GC (G1, ZGC, Shenandoah)                                  │
│                                                                              │
│  内存分配策略                                                                │
│  ├── TCMalloc (Google)                                                      │
│  ├── jemalloc (FreeBSD/Facebook)                                            │
│  ├── ptmalloc (glibc)                                                       │
│  └── mimalloc (Microsoft)                                                   │
│                                                                              │
│  Go 演进                                                                     │
│  ├── Go 1.0: 停止世界                                                       │
│  ├── Go 1.3: 并行清扫                                                       │
│  ├── Go 1.5: 并发三色标记                                                   │
│  ├── Go 1.8: 亚毫秒 STW                                                     │
│  ├── Go 1.14: 异步抢占                                                      │
│  └── Go 1.19: 软内存限制                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. 参考文献

### 经典 GC 文献

1. **McCarthy, J. (1960)**. Recursive Functions of Symbolic Expressions. *CACM*.
2. **Dijkstra, E.W. et al. (1978)**. On-the-fly Garbage Collection: An Exercise in Cooperation. *CACM*.
3. **Jones, R. & Lins, R. (1996)**. Garbage Collection: Algorithms for Automatic Dynamic Memory Management.
4. **Jones, R. et al. (2012)**. The Garbage Collection Handbook. *CRC Press*.

### Go GC 相关

1. **Hudson, R. (2015)**. Go 1.5 Concurrent Garbage Collector.
2. **Go Authors**. Go GC Guide.
3. **Go Authors**. runtime/mgc.go

---

**质量评级**: S (20+ KB)
**完成日期**: 2026-04-02
