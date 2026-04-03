# LD-006: Go 内存分配器内部原理 (Go Memory Allocator Internals)

> **维度**: Language Design
> **级别**: S (40+ KB)
> **标签**: #memory-allocator #tcmalloc #heap #stack #gc #performance
> **权威来源**:
>
> - [Go Memory Allocator](https://github.com/golang/go/tree/master/src/runtime/malloc.go) - Go Authors
> - [TCMalloc](https://goog-perftools.sourceforge.net/doc/tcmalloc.html) - Google
> - [A Fast Storage Allocator](https://dl.acm.org/doi/10.1145/363267.363275) - Knuth

---

## 1. 内存分配基础

### 1.1 内存层次

```
┌─────────────────────────────────────────┐
│            Virtual Memory               │
├─────────────────────────────────────────┤
│  Stack │  Heap  │ Data/BSS │ Text/Code  │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│              mheap                      │
│  ┌────────┐ ┌────────┐ ┌────────┐      │
│  │  span  │ │  span  │ │  span  │ ...  │
│  │(mspan) │ │(mspan) │ │(mspan) │      │
│  └────────┘ └────────┘ └────────┘      │
└─────────────────────────────────────────┘
```

### 1.2 分配策略

| 大小 | 分配器 | 说明 |
|------|--------|------|
| ≤ 16B | Tiny allocator | 合并小对象 |
| ≤ 32KB | Size-class allocator | 分级分配 |
| > 32KB | Large allocator | 直接分配 span |

---

## 2. 内存分配器结构

### 2.1 三级缓存

```go
// mcache: per-P cache (无锁)
type mcache struct {
    tiny             uintptr
    tinyoffset       uintptr
    local_tinyallocs uintptr
    alloc [numSpanClasses]*mspan
}

// mcentral: per-size-class central cache
type mcentral struct {
    spanclass spanClass
    partial [2]mSpanList // 有空闲的 span
    full    [2]mSpanList // 已满的 span
}

// mheap: global heap
type mheap struct {
    lock      mutex
    free      mTreap // 空闲 span 树
    sweepgen  uint32
    arenas    [1 << arenaL1Bits]*[1 << arenaL2Bits]*heapArena
}
```

### 2.2 Span 结构

```go
type mspan struct {
    next      *mspan     // 链表下一个
    prev      *mspan     // 链表上一个
    startAddr uintptr    // 起始地址
    npages    uintptr    // 页数 (8KB per page)

    // 分配相关
    freeindex uintptr    // 下一个空闲对象索引
    nelems    uintptr    // 对象总数
    allocBits *gcBits    // 分配位图
    gcmarkBits *gcBits   // GC 标记位图

    allocCache uint64    // 快速分配缓存

    spanclass   spanClass
    state       mSpanState
}
```

### 2.3 Size Classes

| Class | Size | Objects/Span | Waste |
|-------|------|--------------|-------|
| 0 | 8B | 1024 | 0% |
| 1 | 16B | 512 | 0% |
| 2 | 32B | 256 | 0% |
| ... | ... | ... | ... |
| 66 | 32KB | 1 | 0% |

---

## 3. 分配流程

### 3.1 小对象分配

```go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    // 1. 选择 size class
    spc := makeSpanClass(sizeclass, noscan)

    // 2. 从 mcache 分配
    span := c.alloc[spc]
    v := nextFreeFast(span)
    if v == 0 {
        v = c.nextFree(spc)
    }

    // 3. 返回对象
    x = unsafe.Pointer(v)
    if needzero {
        memclrNoHeapPointers(x, size)
    }
    return x
}
```

### 3.2 大对象分配

```go
// 大于 32KB
func largeAlloc(size uintptr, needzero bool, noscan bool) *mspan {
    // 计算需要的页数
    npages := size >> pageShift
    if size&pageMask != 0 {
        npages++
    }

    // 从 heap 分配 span
    s := mheap_.alloc(npages, spanClass(0))

    if needzero {
        memclrNoHeapPointers(unsafe.Pointer(s.base()), s.npages<<pageShift)
    }
    return s
}
```

### 3.3 Tiny 分配器

```go
// ≤ 16B 的对象合并分配
type mcache struct {
    tiny       uintptr  // 当前 tiny 块指针
    tinyoffset uintptr  // 当前偏移
}

func (c *mcache) tinyAlloc(size uintptr) unsafe.Pointer {
    // 对齐
    roundedSize := roundUpSize(size)

    // 检查是否有空间
    if c.tinyoffset+roundedSize <= maxTinySize {
        x := c.tiny + c.tinyoffset
        c.tinyoffset += roundedSize
        return unsafe.Pointer(x)
    }

    // 分配新的 tiny 块
    c.tiny = nextFreeFast(span)
    c.tinyoffset = roundedSize
    return unsafe.Pointer(c.tiny)
}
```

---

## 4. 运行时行为分析

### 4.1 内存分配完整流程

```
┌─────────────────────────────────────────────────────────────────┐
│                    Memory Allocation Flow                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  new/make/size ≤ 16B                                            │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────────────────────────┐                        │
│  │ Tiny Allocator                      │                        │
│  │ ├── 检查当前 tiny 块空间             │                        │
│  │ ├── 有空间: 直接偏移分配             │                        │
│  │ └── 无空间: 从 mcache 获取新块       │                        │
│  └─────────────────────────────────────┘                        │
│                                                                  │
│  16B < size ≤ 32KB                                               │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────────────────────────┐                        │
│  │ Small Object Allocator              │                        │
│  │                                     │                        │
│  │  1. 确定 size class                 │                        │
│  │     sizeclass = size_to_class(size) │                        │
│  │                                     │                        │
│  │  2. 从 mcache 分配                  │                        │
│  │     span = mcache.alloc[sizeclass]  │                        │
│  │     obj = span.nextFreeIndex()      │                        │
│  │     [无锁，O(1)]                    │                        │
│  │                                     │                        │
│  │  3. mcache 不足                     │                        │
│  │     └── refill from mcentral        │                        │
│  │         ├── 从 partial 链表取 span  │                        │
│  │         └── 或从 mheap 分配新 span  │                        │
│  │                                     │                        │
│  │  4. mcentral 不足                   │                        │
│  │     └── grow from mheap             │                        │
│  │         ├── 从 mheap.free 取 span   │                        │
│  │         └── 或从 OS mmap 新内存     │                        │
│  └─────────────────────────────────────┘                        │
│                                                                  │
│  size > 32KB                                                     │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────────────────────────┐                        │
│  │ Large Object Allocator              │                        │
│  │                                     │                        │
│  │  1. 计算页数: npages = size/pagesize│                        │
│  │                                     │                        │
│  │  2. 从 mheap 直接分配               │                        │
│  │     span = mheap.alloc(npages)      │                        │
│  │     [需要全局锁]                    │                        │
│  │                                     │                        │
│  │  3. mheap 不足                      │                        │
│  │     └── 从 OS mmap 新 arena         │                        │
│  └─────────────────────────────────────┘                        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Span 状态机

```
┌─────────┐    alloc     ┌─────────┐    sweep    ┌─────────┐
│  free   │─────────────►│  inuse  │────────────►│  dead   │
│         │              │         │              │         │
└─────────┘              └─────────┘              └─────────┘
     ▲                                               │
     │              return to heap                   │
     └───────────────────────────────────────────────┘

状态说明:
- free: 空闲，可分配
- inuse: 正在使用，包含已分配对象
- dead: 所有对象可回收，等待清扫
```

### 4.3 内存释放流程

```go
// 延迟释放
func free(v unsafe.Pointer) {
    // 不立即释放，等待 GC 清扫
    // GC 会标记对象为可回收
}

// GC 清扫阶段
func (s *mspan) sweep(preserve bool) bool {
    // 释放未标记的对象
    for i := uintptr(0); i < s.nelems; i++ {
        if !s.gcmarkBits.isMarked(i) {
            // 对象可回收
            s.freeindex = i
            s.allocBits.clear(i)
        }
    }
    // 重置标记位
    s.gcmarkBits, s.allocBits = s.allocBits, s.gcmarkBits
    s.allocBits.clearAll()
}
```

---

## 5. 内存与性能特性

### 5.1 分配延迟

| 操作 | 时间 | 说明 |
|------|------|------|
| Tiny alloc | ~5ns | 无锁，直接偏移 |
| Small alloc | ~25ns | mcache，无锁 |
| Large alloc | ~100ns | 需要锁，可能 mmap |
| Free | ~0ns | 延迟到 GC |

### 5.2 内存碎片

```
内部碎片: 对象大小 < size class 大小
外部碎片: 不连续的空闲内存

Go 策略:
- 67 个 size class 减少内部碎片
- Span-based 分配减少外部碎片
- 延迟归还减少频繁 mmap/munmap
```

**Size Class 内部碎片分析**

| 请求大小 | Size Class | 内部碎片 | 碎片率 |
|----------|------------|----------|--------|
| 1B | 8B | 7B | 87.5% |
| 9B | 16B | 7B | 43.7% |
| 17B | 32B | 15B | 46.9% |
| 100B | 112B | 12B | 10.7% |
| 1000B | 1024B | 24B | 2.3% |

### 5.3 并发性能

```
无锁分配路径:
1. Tiny allocator: 完全无锁，原子操作
2. mcache: per-P 缓存，无竞争
3. mcache refill: 偶尔竞争（从 mcentral）

有锁分配路径:
1. mcentral: 每个 size class 一个锁
2. mheap: 全局锁（大对象分配）
3. OS mmap: 系统调用
```

---

## 6. 多元表征

### 6.1 分配器层次图

```
┌─────────────────────────────────────────┐
│           User Code                     │
│              │                          │
│              ▼ new/make                 │
├─────────────────────────────────────────┤
│         mallocgc                        │
│    ┌────────┬────────┬────────┐        │
│    ▼        ▼        ▼        ▼        │
│  Tiny    Small    Large   Stack        │
│  ≤16B    ≤32KB   >32KB                │
│    │       │       │                  │
├────┼───────┼───────┼──────────────────┤
│    │       │       │                  │
│    ▼       ▼       ▼                  │
│  mcache  mcache  mheap                │
│ (per-P)  (per-P)  (global)            │
│    │       │       │                  │
│    └──┬────┘       │                  │
│       ▼            ▼                  │
│    mcentral      mheap                │
│  (per-size)    (global)               │
│       │            │                  │
│       └────────────┘                  │
│              │                        │
│              ▼                        │
│           OS (mmap)                   │
└─────────────────────────────────────────┘
```

### 6.2 分配路径决策树

```
分配大小?
│
├── ≤ 16B?
│   └── Tiny allocator
│       ├── 当前 tiny 块有空间? → 直接偏移
│       └── 无空间 → 分配新 tiny 块
│
├── ≤ 32KB?
│   └── Size-class allocator
│       ├── mcache 有空闲? → 直接分配
│       ├── mcentral 有空闲 span? → 获取 span
│       └── mheap 有空闲? → 分配新 span
│
└── > 32KB?
    └── Large allocator
        └── 直接从 mheap 分配
```

### 6.3 内存分配器架构全景图

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go Memory Allocator Architecture             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                      User Code                          │    │
│  │                     (new, make)                         │    │
│  └─────────────────────────┬───────────────────────────────┘    │
│                            │                                     │
│                            ▼                                     │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              runtime.mallocgc(size, type, needzero)     │    │
│  └─────────────────────────┬───────────────────────────────┘    │
│                            │                                     │
│         ┌──────────────────┼──────────────────┐                  │
│         │                  │                  │                  │
│         ▼                  ▼                  ▼                  │
│  ┌───────────┐      ┌───────────┐      ┌───────────┐            │
│  │ Tiny      │      │ Small     │      │ Large     │            │
│  │ ≤16B      │      │ ≤32KB     │      │ >32KB     │            │
│  └─────┬─────┘      └─────┬─────┘      └─────┬─────┘            │
│        │                  │                  │                   │
│        ▼                  ▼                  ▼                   │
│  ┌───────────┐      ┌───────────┐      ┌───────────┐            │
│  │ mcache    │      │ mcache    │      │ mheap     │            │
│  │ (per-P)   │      │ (per-P)   │      │ (global)  │            │
│  │ 无锁      │      │ 无锁      │      │ 有锁      │            │
│  └─────┬─────┘      └─────┬─────┘      └─────┬─────┘            │
│        │                  │                  │                   │
│        │         ┌────────┴────────┐         │                   │
│        │         │                 │         │                   │
│        │         ▼                 ▼         │                   │
│        │  ┌───────────┐     ┌───────────┐    │                   │
│        │  │ mcentral  │     │ mcentral  │    │                   │
│        └──►│ (empty)   │     │ (full)    │◄───┘                   │
│           └─────┬─────┘     └─────┬─────┘                        │
│                 │                 │                               │
│                 └────────┬────────┘                               │
│                          │                                        │
│                          ▼                                        │
│                   ┌───────────┐                                   │
│                   │ mheap     │                                   │
│                   │ (global)  │                                   │
│                   └─────┬─────┘                                   │
│                         │                                         │
│              ┌──────────┼──────────┐                              │
│              ▼          ▼          ▼                              │
│        ┌─────────┐ ┌─────────┐ ┌─────────┐                       │
│        │ arenas  │ │ treap   │ │ scavenger│                       │
│        │ (bitmap)│ │ (free)  │ │ (回收)  │                       │
│        └────┬────┘ └─────────┘ └─────────┘                       │
│             │                                                     │
│             ▼                                                     │
│        ┌─────────┐                                                │
│        │ OS mmap │                                                │
│        └─────────┘                                                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 7. 完整代码示例

### 7.1 内存统计

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    var m1, m2 runtime.MemStats

    // 分配前
    runtime.GC()
    runtime.ReadMemStats(&m1)

    // 分配内存
    data := make([][]byte, 1000)
    for i := range data {
        data[i] = make([]byte, 1024)
    }

    // 分配后
    runtime.ReadMemStats(&m2)

    fmt.Printf("Alloc increase: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
    fmt.Printf("TotalAlloc: %d KB\n", m2.TotalAlloc/1024)
    fmt.Printf("Mallocs: %d\n", m2.Mallocs-m1.Mallocs)
    fmt.Printf("HeapSys: %d KB\n", m2.HeapSys/1024)
    fmt.Printf("HeapIdle: %d KB\n", m2.HeapIdle/1024)
}
```

### 7.2 对象池

```go
package main

import (
    "sync"
    "testing"
)

// 减少 GC 压力的对象池
type Buffer struct {
    data [1024]byte
}

var pool = sync.Pool{
    New: func() interface{} {
        return &Buffer{}
    },
}

func getBuffer() *Buffer {
    return pool.Get().(*Buffer)
}

func putBuffer(b *Buffer) {
    pool.Put(b)
}

func process(data []byte) {
    buf := getBuffer()
    defer putBuffer(buf)

    copy(buf.data[:], data)
    // 处理...
}

// 基准测试
func BenchmarkWithPool(b *testing.B) {
    data := make([]byte, 100)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        process(data)
    }
}

func BenchmarkWithoutPool(b *testing.B) {
    data := make([]byte, 100)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        buf := &Buffer{}
        copy(buf.data[:], data)
    }
}
```

### 7.3 减少分配

```go
package main

import (
    "strings"
    "testing"
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

// 基准测试
func BenchmarkConcatBad(b *testing.B) {
    items := []string{"hello", " ", "world", "!"}
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = concatBad(items)
    }
}

func BenchmarkConcatGood(b *testing.B) {
    items := []string{"hello", " ", "world", "!"}
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = concatGood(items)
    }
}

func BenchmarkConcatBetter(b *testing.B) {
    items := []string{"hello", " ", "world", "!"}
    var buf strings.Builder
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = concatBetter(items, &buf)
    }
}
```

### 7.4 逃逸分析示例

```go
package main

// 栈分配 - 不逃逸
func stackAlloc() int {
    x := 42
    return x
}

// 堆分配 - 返回指针，逃逸
func heapAlloc() *int {
    x := 42
    return &x  // 逃逸到堆
}

// 切片逃逸 - 容量不确定
func sliceEscape(n int) []int {
    return make([]int, n)  // n 不是常量，逃逸
}

// 不逃逸 - 容量确定
func sliceNoEscape() []int {
    return make([]int, 100)  // 容量常量，在栈上
}

// 接口逃逸 - 装箱
func interfaceEscape() interface{} {
    x := 42
    return x  // 装箱，逃逸
}

// 闭包逃逸
func closureEscape() func() int {
    x := 42
    return func() int {
        return x  // x 逃逸
    }
}

// go build -gcflags="-m" 查看逃逸分析结果
```

---

## 8. 最佳实践与反模式

### 8.1 ✅ 最佳实践

```go
// 1. 预分配切片容量
func collect(n int) []int {
    result := make([]int, 0, n)  // 预分配
    for i := 0; i < n; i++ {
        result = append(result, i)
    }
    return result
}

// 2. 使用 sync.Pool 复用临时对象
var bufPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func process(data []byte) {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf)
    // 使用 buf...
}

// 3. 避免在热路径装箱
// Bad
func process(iface interface{}) {
    n := iface.(int)
    _ = n
}

// Good
func processInt(n int) {
    _ = n
}

// 4. 使用值类型减少 GC 扫描
// Bad
type Node struct {
    Value interface{}  // 指针，需要扫描
}

// Good
type IntNode struct {
    Value int  // 非指针，无需扫描
}
```

### 8.2 ❌ 反模式

```go
// 1. 在循环中分配
func badLoop() {
    for i := 0; i < 1000000; i++ {
        buf := make([]byte, 1024)  // 每次迭代都分配
        _ = buf
    }
}

// 2. 不必要的指针
func badPointer(items []*Item) {  // 指针切片，更多 GC 压力
    for _, item := range items {
        _ = item.Value
    }
}

// 3. 接口过度使用
type Stringer interface {
    String() string
}

func badPrint(s Stringer) {  // 接口装箱
    fmt.Println(s.String())
}

// 4. 闭包捕获大对象
func badClosure() {
    bigData := make([]byte, 1000000)
    _ = bigData

    go func() {
        // 只使用一小部分
        _ = bigData[0]  // 但整个 bigData 都逃逸了
    }()
}
```

---

## 9. 关系网络

```
Go Memory Allocator
├── TCMalloc-inspired Design
│   ├── Thread-local cache (mcache)
│   ├── Size classes (67 classes)
│   ├── Central free lists (mcentral)
│   └── Global heap (mheap)
├── Garbage Collection
│   ├── Mark-sweep algorithm
│   ├── Write barrier
│   ├── Lazy sweep on allocate
│   └── Span-based management
├── OS Interface
│   ├── mmap/munmap
│   ├── madvise (MADV_DONTNEED)
│   └── sysAlloc/sysFree
└── Performance Optimizations
    ├── Lock-free fast path
    ├── Cache-friendly layout
    ├── Lazy scavenging
    └── Transparent huge pages
```

---

## 10. 参考文献

1. **Ghemawat, S. & Menage, P.** TCMalloc: Thread-Caching Malloc.
2. **Knuth, D.** The Art of Computer Programming, Vol 1: Fundamental Algorithms.
3. **Go Authors.** Go Runtime Source Code (src/runtime/malloc.go).
4. **Aken, J.** Understanding Go's Memory Allocator.

---

**质量评级**: S (40KB)
**完成日期**: 2026-04-02
