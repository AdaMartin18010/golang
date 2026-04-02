# LD-006: Go 内存分配器内部原理 (Go Memory Allocator Internals)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #memory-allocator #tcmalloc #heap #stack #gc
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

## 4. 内存释放

### 4.1 延迟释放

```go
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

### 4.2 Span 回收

```
当 span 所有对象都空闲时:
1. 从 mcentral 移除
2. 返回给 mheap
3. mheap 可能保留或归还给 OS
```

---

## 5. 性能分析

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

---

## 7. 代码示例

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
}
```

### 7.2 对象池

```go
package main

import (
    "sync"
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
```

---

## 8. 关系网络

```
Go Memory Allocator
├── TCMalloc-inspired Design
│   ├── Thread-local cache (mcache)
│   ├── Size classes
│   └── Central free lists
├── Garbage Collection
│   ├── Mark-sweep
│   ├── Write barrier
│   └── Sweep on allocate
└── OS Interface
    ├── mmap/munmap
    ├── madvise
    └── sysAlloc
```

---

**质量评级**: S (15KB)
**完成日期**: 2026-04-02
