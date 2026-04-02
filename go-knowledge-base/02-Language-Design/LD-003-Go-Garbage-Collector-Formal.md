# LD-003: Go 垃圾回收器的形式化理论 (Go Garbage Collector: Formal Theory)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #garbage-collection #tricolor #concurrent-gc #memory-management
> **权威来源**:
>
> - [Go GC Guide](https://go.dev/doc/gc-guide) - Go Authors
> - [Concurrent Garbage Collection](https://dl.acm.org/doi/10.1145/359580.359587) - Dijkstra et al.
> - [Tri-color Marking](https://en.wikipedia.org/wiki/Tracing_garbage_collection) - Garbage Collection
> - [Go 1.5 GC](https://go.dev/s/go15gc) - Rick Hudson

---

## 1. 形式化基础

### 1.1 垃圾回收理论

**定义 1.1 (垃圾)**
垃圾是不再被任何可达对象引用的内存对象。

```
Garbage = { o ∈ Heap | ¬∃r ∈ Roots : r →* o }
```

**定义 1.2 (根集合)**
根集合是程序直接访问的对象集合：

```
Roots = Globals ∪ Stack ∪ Registers
```

**定义 1.3 (可达性)**
对象 o 从根 r 可达当存在引用链：

```
r →* o ≡ r → o1 → o2 → ... → o
```

**定理 1.1 (垃圾安全性)**
垃圾回收器不会回收可达对象。

```
∀o: Collected(o) ⇒ o ∈ Garbage
```

### 1.2 Go GC 设计目标

**公理 1.1 (低延迟优先)**
GC 暂停时间目标 < 100μs。

**公理 1.2 (并发执行)**
垃圾回收与 mutator（应用程序）并发执行。

**公理 1.3 (自适应)**
GC 根据堆增长率和 CPU 可用性自适应调整。

---

## 2. 三色标记-清除算法

### 2.1 三色不变式

**定义 2.1 (三色抽象)**

```
White: 待检查对象（可能是垃圾）
Grey:  已发现但未检查子引用的对象
Black: 已检查且子引用也是非白色的对象
```

**定义 2.2 (三色不变式)**

```
¬∃b ∈ Black, w ∈ White : b → w
```

黑色对象不直接引用白色对象。

**定理 2.1 (三色不变式保持)**
若维持三色不变式，则当灰色集合为空时，白色对象即为垃圾。

*证明*：

1. 所有从根可达对象要么是黑色（已检查），要么是灰色（待检查）
2. 灰色集合为空意味着所有可达对象都是黑色
3. 根据不变式，黑色对象不引用白色对象
4. 因此白色对象不可达，是垃圾

### 2.2 并发三色算法

**算法 2.1 (并发标记)**

```
1. 初始化: 所有对象白色，根对象灰色
2. while 灰色集合非空:
3.   从灰色取对象 g
4.   标记 g 为黑色
5.   for g 引用的每个对象 w:
6.     if w 是白色:
7.       标记 w 为灰色
8. 收集白色对象
```

**写屏障 (Write Barrier)**

```go
// 当 mutator 修改指针时触发
func writeBarrier(slot, ptr *Object) {
    if currentPhase == GCmark && slot is black && ptr is white {
        shade(ptr)  // 标记 ptr 为灰色
    }
    *slot = ptr
}
```

---

## 3. Go GC 实现

### 3.1 GC 阶段

```
┌─────────────────────────────────────────────────────────────┐
│                    Go GC Cycle                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  STW 1 (Sweep Termination)                                   │
│  ├── 停止所有 mutator                                        │
│  └── 准备标记阶段                                            │
│                                                              │
│  Mark Phase (并发)                                           │
│  ├── Mark Setup: 启用写屏障                                  │
│  ├── Mark: 遍历对象图                                        │
│  └── Mark Termination: 完成标记                              │
│                                                              │
│  STW 2 (Mark Termination)                                    │
│  ├── 停止 mutator                                            │
│  └── 计算回收统计                                            │
│                                                              │
│  Sweep Phase (并发)                                          │
│  ├── 逐 span 清扫                                            │
│  └── 回收白色对象内存                                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 内存分配器

**Span 结构**

```go
// runtime/mheap.go
type mspan struct {
    next      *mspan
    prev      *mspan
    startAddr uintptr
    npages    uintptr
    freeindex uintptr
    allocCache uint64
    allocBits  *gcBits
    gcmarkBits *gcBits
}
```

**大小分级**

| 类别 | 大小范围 | 用途 |
|------|----------|------|
| Tiny | ≤ 16B | 小对象合并 |
| Small | 16B - 32KB | 按大小分级分配 |
| Large | > 32KB | 直接分配 |

### 3.3 GC 触发条件

**触发条件**

```
触发堆大小 = GOGC/100 * 存活堆大小
```

- GOGC=100 (默认): 当堆增长到存活对象的两倍时触发
- GOGC=off: 禁用 GC
- GOGC=-1: 由 runtime.GC() 手动触发

---

## 4. 性能分析

### 4.1 GC 开销模型

```
GC CPU 目标 = GOGC / (GOGC + 100)
```

默认 GOGC=100 时，GC 使用 50% 可用 CPU。

### 4.2 停顿时间分解

| 阶段 | 目标 | 实际 (Go 1.20+) |
|------|------|-----------------|
| STW 1 | < 10μs | ~1-5μs |
| STW 2 | < 100μs | ~10-50μs |
| 总停顿 | < 500μs | ~100-500μs |

### 4.3 内存开销

| 组件 | 开销 | 说明 |
|------|------|------|
| GC 位图 | 1/64 堆大小 | 每对象 2 bits |
| Work buffer | 可变 | 标记工作队列 |
| Stack barrier | 栈大小 | 扫描 goroutine 栈 |

---

## 5. 多元表征

### 5.1 三色标记可视化

```
初始状态:
  Roots → A → B → C
          ↓
          D

所有节点白色

标记 Roots:
  Roots(黑) → A(灰) → B(白) → C(白)
              ↓
              D(白)

标记 A:
  Roots(黑) → A(黑) → B(灰) → C(白)
              ↓
              D(灰)

标记 B:
  Roots(黑) → A(黑) → B(黑) → C(灰)
              ↓
              D(灰)

完成 (灰色为空):
  Roots(黑) → A(黑) → B(黑) → C(黑)
              ↓
              D(黑)

白色集合为空，无垃圾
```

### 5.2 GC 决策树

```
遇到内存问题?
│
├── 高内存使用?
│   ├── 检查对象数量: runtime.ReadMemStats
│   ├── 使用 pprof heap 分析
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
│   └── 考虑 arena 分配 (Go 1.20+)
│
└── 内存泄漏?
    ├── 使用 pprof 比较 heap profile
    ├── 检查 goroutine 泄漏
    └── 检查未关闭的 channel
```

### 5.3 GC 算法对比

| 算法 | 停顿 | 吞吐量 | 内存碎片 | 实现复杂度 |
|------|------|--------|----------|------------|
| 标记-清除 | 长 | 中 | 有 | 低 |
| 引用计数 | 无 | 高 | 无 | 中 |
| 复制 | 中 | 高 | 无 | 中 |
| 分代 | 短 | 高 | 有 | 高 |
| Go 并发三色 | 极短 | 中 | 有 | 高 |

---

## 6. 代码示例

### 6.1 GC 调优

```go
package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
)

func main() {
    // 设置 GC 目标百分比
    debug.SetGCPercent(100)  // 默认

    // 手动触发 GC
    runtime.GC()

    // 读取 GC 统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("NumGC: %d\n", m.NumGC)
    fmt.Printf("PauseNs: %d μs\n", m.PauseNs[(m.NumGC+255)%256]/1000)
}
```

### 6.2 对象池模式

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process() {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)

    // 使用 buf...
}
```

### 6.3 减少分配

```go
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
```

---

## 7. 关系网络

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
│  └── Go 1.19: Soft memory limit                                 │
│                                                                  │
│  调优工具                                                        │
│  ├── GODEBUG=gctrace=1                                          │
│  ├── runtime.ReadMemStats                                       │
│  ├── debug.SetGCPercent                                         │
│  └── pprof heap                                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 8. 参考文献

1. **Dijkstra, E. W. et al.** On-the-fly Garbage Collection.
2. **Hudson, R.** Go 1.5 Concurrent Garbage Collector.
3. **Jones, R.** The Garbage Collection Handbook.

---

**质量评级**: S (16KB)
**完成日期**: 2026-04-02
