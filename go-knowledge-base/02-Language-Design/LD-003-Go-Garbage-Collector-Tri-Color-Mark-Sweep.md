# LD-003: Go 三色标记-清除垃圾回收器详解 (Go Tri-Color Mark-Sweep GC Deep Dive)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #gc #tricolor #marksweep #concurrent #memory #runtime
> **权威来源**:
>
> - [Go GC Implementation](https://github.com/golang/go/tree/master/src/runtime/mgc.go) - Go Authors
> - [Tri-color Marking](https://en.wikipedia.org/wiki/Tracing_garbage_collection) - Wikipedia
> - [Dijkstra GC](https://dl.acm.org/doi/10.1145/359580.359587) - Dijkstra et al.
> - [Go GC Guide](https://go.dev/doc/gc-guide) - Go Authors

---

## 1. 三色标记-清除基础

### 1.1 算法定义

**定义 1.1 (三色抽象)**
三色标记算法将对象分为三种颜色状态：

```
白色 (White): 尚未访问的对象，候选垃圾
灰色 (Grey):  已访问但子对象未完全访问的对象
黑色 (Black): 已完全访问的对象，保留对象
```

**定义 1.2 (三色不变式)**
三色不变式是并发 GC 正确性的核心保证：

```
∀b ∈ Black, ∀w ∈ White: ¬(b → w)
```

即：黑色对象不能直接引用白色对象。

**定理 1.1 (三色算法正确性)**
当灰色集合为空时，白色对象即为垃圾。

*证明*：

1. 初始时，所有对象都是白色
2. 根对象被标记为灰色
3. 处理灰色对象：标记为黑色，子对象标记为灰色
4. 重复直到灰色集合为空
5. 此时，所有从根可达对象都是黑色
6. 根据不变式，黑色对象不引用白色对象
7. 因此白色对象不可达，是垃圾

### 1.2 基本算法流程

**算法 1.1 (串行三色标记)**

```
初始化:
  ∀o ∈ Heap: color[o] = white
  GreySet = ∅

标记根:
  ∀r ∈ Roots:
    color[r] = grey
    GreySet.add(r)

标记阶段:
  while GreySet ≠ ∅:
    o = GreySet.remove()
    color[o] = black
    for child in o.children:
      if color[child] == white:
        color[child] = grey
        GreySet.add(child)

清除阶段:
  ∀o ∈ Heap:
    if color[o] == white:
      free(o)
```

---

## 2. 并发三色标记

### 2.1 并发问题

当 GC 与应用程序（mutator）并发执行时，mutator 可能修改对象引用关系，破坏三色不变式：

```
问题场景:
  1. GC: 标记 A 为黑色
  2. Mutator: A.child = B (B 是白色)
  3. GC: 完成标记
  4. B 被错误回收（尽管被黑色 A 引用）
```

### 2.2 写屏障

**定义 2.1 (Dijkstra 插入写屏障)**

```
当修改指针 slot = ptr 时:
  if GC 进行中 AND slot 是黑色 AND ptr 是白色:
    shade(ptr)  // 标记 ptr 为灰色
  *slot = ptr
```

**定义 2.2 (Yuasa 删除写屏障)**

```
当修改指针 slot = ptr 时:
  if GC 进行中 AND slot 是白色:
    shade(slot)  // 标记 slot 为灰色（保留旧引用）
  *slot = ptr
```

**定义 2.3 (混合写屏障)**
Go 使用 Dijkstra + Yuasa 混合写屏障：

```
shade(ptr)        // Dijkstra: 保护新引用
shade(*slot)      // Yuasa: 保护旧引用
*slot = ptr
```

### 2.3 Go 写屏障实现

```go
// runtime/mbarrier.go

// gcWriteBarrier 是写屏障入口
func gcWriteBarrier(ptr, slot uintptr) {
    // 标记新值
    if gcphase == _GCmark {
        greyobject(ptr)
    }
    // 标记旧值（删除屏障）
    if gcphase == _GCmark {
        greyobject(*(*uintptr)(unsafe.Pointer(slot)))
    }
    // 执行写操作
    *(*uintptr)(unsafe.Pointer(slot)) = ptr
}

// greyobject 将对象标记为灰色
func greyobject(obj uintptr) {
    if obj == 0 {
        return
    }
    // 设置标记位
    if !markBits.isMarked(obj) {
        markBits.setMarked(obj)
        // 加入工作队列
        work.putFast(obj)
    }
}
```

---

## 3. Go GC 实现细节

### 3.1 GC 状态机

```
┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐
│  _GCoff │────►│_GCmark  │────►│_GCmarktm│────►│  _GCoff │
│ (清扫)  │     │(标记)   │     │(标记结束)│    │ (循环)  │
└─────────┘     └─────────┘     └─────────┘     └─────────┘
     ▲                                              │
     └──────────────────────────────────────────────┘
```

**状态定义**

```go
const (
    _GCoff = iota      // GC 关闭，清扫阶段
    _GCmark            // 标记阶段
    _GCmarktermination // 标记终止
)
```

### 3.2 GC 周期详解

**阶段 1: Sweep Termination (STW)**

```
1. 停止所有 mutator
2. 等待当前清扫完成
3. 清除准备标记阶段
4. 启动世界（ Start the World）
```

**阶段 2: Mark (并发)**

```
1. 启用写屏障
2. 标记根对象（goroutine 栈、全局变量）
3. 并发标记对象图
4. 工作窃取平衡负载
```

**阶段 3: Mark Termination (STW)**

```
1. 停止 mutator
2. 完成剩余标记工作
3. 计算下次 GC 触发点
4. 禁用写屏障
5. 启动世界
```

**阶段 4: Sweep (并发)**

```
1. 逐 span 清扫
2. 回收白色对象内存
3. 延迟清扫（在分配时）
```

### 3.3 标记工作队列

**工作窃取算法**

```go
type gcWork struct {
    wbuf1, wbuf2 *workbuf  // 双缓冲
    bytesMarked  uint64
    heapScanWork int64
}

// 添加工作
func (w *gcWork) put(obj uintptr) {
    // 优先本地缓冲区
    if w.wbuf1.putFast(obj) {
        return
    }
    // 慢路径
    w.putSlow(obj)
}

// 获取工作
func (w *gcWork) get() uintptr {
    // 1. 本地缓冲区
    if obj := w.wbuf1.getFast(); obj != 0 {
        return obj
    }
    // 2. 交换缓冲区
    w.wbuf1, w.wbuf2 = w.wbuf2, w.wbuf1
    // 3. 窃取其他 P
    if obj := w.steal(); obj != 0 {
        return obj
    }
    // 4. 全局队列
    return w.getFromGlobal()
}
```

### 3.4 根对象扫描

**根集合组成**

```
1. goroutine 栈（每个 goroutine 的栈帧）
2. 全局变量
3. 寄存器中的指针
4. 运行时数据结构
```

**栈扫描实现**

```go
func scanstack(gp *g, gcw *gcWork) {
    // 保守扫描栈帧
    // 标记所有可能的指针
    var n uintptr
    for n = 0; n < gp.stack.hi-gp.stack.lo; n += ptrSize {
        p := *(*uintptr)(unsafe.Pointer(gp.stack.lo + n))
        if isHeapPointer(p) {
            gcw.put(p)
        }
    }
}
```

---

## 4. 清扫阶段

### 4.1 延迟清扫

Go 采用延迟清扫策略，在分配时惰性地清扫：

```go
// runtime/malloc.go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    // 1. 检查是否需要清扫当前 span
    if gcphase == _GCoff && span.sweepgen != sweep.gen {
        span.sweep(true)
    }

    // 2. 分配对象
    // ...
}
```

### 4.2 位图标记

**gcBits 结构**

```go
// 每个 word (8 bytes) 对应 2 bits
type gcBits struct {
    bits []byte
}

// 位含义
const (
    gcBits_white = 0  // 00: 白色
    gcBits_grey  = 1  // 01: 灰色
    gcBits_black = 2  // 10: 黑色
    gcBits_reserved = 3  // 11: 保留
)
```

**内存布局**

```
Heap:
┌─────────────────────────────────────────┐
│  Object 0 │ Object 1 │ ... │ Object N  │
└─────────────────────────────────────────┘
     │           │              │
     ▼           ▼              ▼
┌─────────────────────────────────────────┐
│  00 (白)  │  10 (黑) │ ... │  01 (灰)  │  gcBits
└─────────────────────────────────────────┘
```

---

## 5. GC 触发与调优

### 5.1 触发条件

**GC 触发公式**

```
目标堆大小 = 存活对象大小 + 存活对象大小 × GOGC/100

默认 GOGC=100:
  触发堆大小 = 2 × 存活对象大小
```

**触发方式**

```go
// 1. 自动触发（堆增长）
// 由 gcTriggerHeap 监控

// 2. 手动触发
runtime.GC()

// 3. 时间触发
gcTriggerTime // 2 分钟未 GC
```

### 5.2 GC 调优参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| GOGC | 100 | GC 触发百分比 |
| GOMEMLIMIT | 无 | 内存软限制 |
| debug.SetGCPercent | 100 | 运行时设置 |
| debug.SetMemoryLimit | MaxInt64 | 内存限制 |

### 5.3 调优策略

```go
// 减少 GC 频率（增加延迟）
debug.SetGCPercent(200)  // 堆增长 3 倍才触发

// 内存限制
debug.SetMemoryLimit(10 << 30)  // 10GB

// 读取统计
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("GC #%d, Pause: %dμs\n", m.NumGC, m.PauseNs[(m.NumGC+255)%256]/1000)
```

---

## 6. 性能分析

### 6.1 GC 开销模型

```
GC CPU 目标 = GOGC / (GOGC + 100)

GOGC=100:  50% CPU
GOGC=50:   33% CPU
GOGC=200:  67% CPU
```

### 6.2 停顿时间分解

| 阶段 | 目标 | 实际 (Go 1.20+) | 优化策略 |
|------|------|-----------------|----------|
| STW 1 | < 10μs | ~1-5μs | 并行标记准备 |
| STW 2 | < 100μs | ~10-50μs | 快速终止 |
| 总停顿 | < 500μs | ~100-500μs | 分布式终止 |

### 6.3 内存开销

| 组件 | 开销 | 说明 |
|------|------|------|
| GC 位图 | 1/64 堆大小 | 每对象 2 bits |
| Work buffer | 可变 | 标记工作队列 |
| Stack barrier | 栈大小 | 扫描 goroutine 栈 |

---

## 7. 多元表征

### 7.1 对象图演变

```
初始状态:
  Root → A → B → C
         ↓
         D
  所有节点白色

标记根后:
  Root(黑) → A(灰) → B(白) → C(白)
             ↓
             D(白)

标记 A 后:
  Root(黑) → A(黑) → B(灰) → C(白)
             ↓
             D(灰)

完成标记:
  Root(黑) → A(黑) → B(黑) → C(黑)
             ↓
             D(黑)
  白色集合为空

新增对象（写屏障保护）:
  A(黑) ──新引用──► X(灰)

  写屏障: A 变灰或 X 变灰
```

### 7.2 GC 决策树

```
遇到内存问题?
│
├── 高内存使用?
│   ├── 检查对象数量: runtime.ReadMemStats
│   ├── 使用 pprof heap 分析
│   ├── 检查内存泄漏
│   └── 考虑对象池 sync.Pool
│
├── GC 停顿高?
│   ├── 检查 GC 频率: GODEBUG=gctrace=1
│   ├── 减少堆分配
│   ├── 增加 GOGC (牺牲内存换 CPU)
│   └── 检查 finalizers
│
├── CPU 使用高?
│   ├── 使用 CPU profile
│   ├── 减少不必要的 allocation
│   ├── 使用 arena 分配 (Go 1.20+)
│   └── 调整 GOGC
│
└── 内存泄漏?
    ├── 使用 pprof 比较 heap profile
    ├── 检查 goroutine 泄漏
    ├── 检查未关闭的 channel
    └── 检查全局缓存
```

### 7.3 GC 算法对比

| 算法 | 停顿 | 吞吐量 | 内存碎片 | 实现复杂度 | 适用场景 |
|------|------|--------|----------|------------|----------|
| 标记-清除 | 长 | 中 | 有 | 低 | 简单应用 |
| 引用计数 | 无 | 高 | 无 | 中 | 实时系统 |
| 复制 | 中 | 高 | 无 | 中 | 对象存活率低 |
| 分代 | 短 | 高 | 有 | 高 | 对象生命周期差异大 |
| Go 并发三色 | 极短 | 中 | 有 | 高 | 低延迟要求 |

---

## 8. 代码示例

### 8.1 GC 统计监控

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func printGCStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("=== Memory Stats ===\n")
    fmt.Printf("Alloc: %d MB\n", m.Alloc/1024/1024)
    fmt.Printf("TotalAlloc: %d MB\n", m.TotalAlloc/1024/1024)
    fmt.Printf("Sys: %d MB\n", m.Sys/1024/1024)
    fmt.Printf("NumGC: %d\n", m.NumGC)
    fmt.Printf("PauseNs (last): %d μs\n", m.PauseNs[(m.NumGC+255)%256]/1000)
    fmt.Printf("GCCPUFraction: %.4f\n", m.GCCPUFraction)
    fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("HeapObjects: %d\n", m.HeapObjects)
}

func main() {
    // 分配一些内存
    data := make([][]byte, 100)
    for i := range data {
        data[i] = make([]byte, 1024*1024) // 1MB
    }

    printGCStats()

    // 释放引用
    data = nil
    runtime.GC()

    time.Sleep(time.Second)
    printGCStats()
}
```

### 8.2 对象池模式

```go
package main

import (
    "sync"
    "sync/atomic"
)

var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

var allocCount int64

func getBuffer() []byte {
    atomic.AddInt64(&allocCount, 1)
    return bufferPool.Get().([]byte)
}

func putBuffer(buf []byte) {
    if cap(buf) >= 1024 {
        bufferPool.Put(buf[:1024])
    }
}

// 使用示例
func process(data []byte) []byte {
    buf := getBuffer()
    defer putBuffer(buf)

    // 处理数据...
    copy(buf, data)
    return buf[:len(data)]
}
```

### 8.3 减少分配

```go
package main

import (
    "strings"
)

// Bad: 每次调用都分配
func concatBad(items []string) string {
    var result string
    for _, item := range items {
        result += item  // 每次分配新字符串
    }
    return result
}

// Good: 预分配容量
func concatGood(items []string) string {
    var n int
    for _, item := range items {
        n += len(item)
    }

    var buf strings.Builder
    buf.Grow(n)  // 预分配
    for _, item := range items {
        buf.WriteString(item)
    }
    return buf.String()
}

// Better: 复用 buffer
func concatBetter(items []string, buf *strings.Builder) string {
    buf.Reset()
    for _, item := range items {
        buf.WriteString(item)
    }
    return buf.String()
}
```

---

## 9. 关系网络

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go GC Context                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  GC 算法家族                                                     │
│  ├── Mark-Sweep (McCarthy, 1960)                                │
│  ├── Reference Counting (Collins, 1960)                         │
│  ├── Copying (Fenichel, 1969)                                   │
│  ├── Generational (Lieberman, 1983)                             │
│  ├── Tri-color (Dijkstra, 1978)                                 │
│  └── Concurrent GC (Boehm, 1993)                                │
│                                                                  │
│  Go GC 演进                                                      │
│  ├── Go 1.0: 停止世界标记-清除                                  │
│  ├── Go 1.3: 并行清扫                                           │
│  ├── Go 1.5: 并发三色标记 (1ms STW)                             │
│  ├── Go 1.8: 亚毫秒 STW                                         │
│  ├── Go 1.19: Soft memory limit                                 │
│  └── Go 1.20+: Pacer 改进                                       │
│                                                                  │
│  调优工具                                                        │
│  ├── GODEBUG=gctrace=1                                          │
│  ├── runtime.ReadMemStats                                       │
│  ├── debug.SetGCPercent                                         │
│  ├── debug.SetMemoryLimit                                       │
│  └── pprof heap                                                 │
│                                                                  │
│  相关概念                                                        │
│  ├── Write Barrier                                              │
│  ├── Work Stealing                                              │
│  ├── Lazy Sweeping                                              │
│  └── Escape Analysis                                            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 10. 参考文献

### 经典论文

1. **Dijkstra, E. W. et al.** (1978). On-the-fly Garbage Collection: An Exercise in Cooperation.
2. **McCarthy, J.** (1960). Recursive Functions of Symbolic Expressions.
3. **Lieberman, H. & Hewitt, C.** (1983). A Real-Time Garbage Collector Based on the Lifetimes of Objects.

### Go 相关

1. **Hudson, R.** (2015). Go 1.5 Concurrent Garbage Collector.
2. **Go Authors.** runtime/mgc.go, runtime/mbarrier.go.

### 书籍

1. **Jones, R. & Lins, R.** (1996). Garbage Collection: Algorithms for Automatic Dynamic Memory Management.

---

## 11. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go GC Toolkit                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  核心原则                                                        │
│  ═══════════════════════════════════════════════════════════    │
│  1. 低延迟优先 (< 100μs STW)                                    │
│  2. 并发执行 (GC 与 mutator 并行)                               │
│  3. 自适应 (根据负载调整)                                        │
│                                                                  │
│  三色标记记忆法:                                                 │
│  • 白色 = 垃圾（需要清除）                                       │
│  • 灰色 = 待处理（工作队列）                                     │
│  • 黑色 = 存活（已确认）                                         │
│                                                                  │
│  写屏障规则:                                                     │
│  黑色对象不能引用白色对象!                                       │
│  当黑→白引用出现时，写屏障将白色对象变灰                         │
│                                                                  │
│  调优检查清单:                                                   │
│  □ 使用 runtime.ReadMemStats 监控                               │
│  □ 用 GODEBUG=gctrace=1 查看 GC 行为                            │
│  □ 适当调整 GOGC（默认 100）                                    │
│  □ 使用 sync.Pool 减少分配                                       │
│  □ 关注 pprof heap profile                                       │
│                                                                  │
│  常见陷阱:                                                       │
│  ❌ 大量小对象分配导致频繁 GC                                   │
│  ❌ 全局 map/cache 持有对象不释放                               │
│  ❌ goroutine 泄漏导致栈不回收                                  │
│  ❌ 使用 finalizer 影响 GC 性能                                 │
│                                                                  │
│  性能优化:                                                       │
│  • 预分配 slice 和 map                                          │
│  • 复用对象 (sync.Pool)                                         │
│  • 避免不必要的指针                                             │
│  • 使用值类型而非指针（小对象）                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16KB)
**完成日期**: 2026-04-02
