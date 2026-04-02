# LD-006: Go 内存分配器内部机制 (Go Memory Allocator Internals)

> **维度**: Language Design
> **级别**: S (20+ KB)
> **标签**: #go-memory #allocator #mcache #mcentral #mheap #tcmalloc
> **权威来源**: [Go Memory Allocator](https://github.com/golang/go/blob/master/src/runtime/malloc.go), [TCMalloc](http://goog-perftools.sourceforge.net/doc/tcmalloc.html)

---

## 架构层次

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Go Memory Allocator Hierarchy                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Goroutine (G)                                                              │
│       │                                                                       │
│       ▼                                                                       │
│  ┌─────────────────┐     无锁，per-P                                          │
│  │    mcache       │     本地缓存，<32KB 对象                                  │
│  │  (per-P cache)  │                                                         │
│  │                 │                                                         │
│  │  tiny allocator │     <16B 小对象特殊处理                                   │
│  │  size classes   │     67 种规格 (8B ~ 32KB)                                │
│  └────────┬────────┘                                                         │
│           │ 缓存未命中                                                          │
│           ▼                                                                   │
│  ┌─────────────────┐     部分加锁，全局                                          │
│  │    mcentral     │     中心缓存，按 size class 组织                           │
│  │  (per-size      │                                                         │
│  │   central)      │                                                         │
│  │                 │                                                         │
│  │  nonempty list  │     有空闲 span 的列表                                    │
│  │  empty list     │     需要 GC 扫描的列表                                    │
│  └────────┬────────┘                                                         │
│           │ 无可用 span                                                         │
│           ▼                                                                   │
│  ┌─────────────────┐     全局加锁，系统调用                                      │
│  │     mheap       │     堆管理，>32KB 大对象                                    │
│  │   (global)      │                                                         │
│  │                 │                                                         │
│  │  arenas         │     64MB 连续内存块                                        │
│  │  spans          │     内存页管理                                             │
│  └─────────────────┘                                                         │
│           │                                                                   │
│           ▼                                                                   │
│      sysAlloc                                                                 │
│     (mmap/brk)                                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心数据结构

```go
// src/runtime/mcache.go

type mcache struct {
    // Tiny allocator
    tiny       uintptr
    tinyoffset uintptr
    tinyAllocs uintptr

    // mspan caches (per size class)
    alloc [numSpanClasses]*mspan  // 136 种 (67 * 2, noscan/scanned)

    // Local stack allocator
    stackcache [stackCacheSize]stackfreelist

    // Flushed to central
    local_largefree  uintptr
    local_nlargefree uintptr
    local_smallfree  [numSpanClasses]uintptr
}

// 67 size classes (8B ~ 32KB)
// class  bytes/obj  bytes/span  objects  waste
//     1          8        8192     1024   87.50%
//     2         16        8192      512   43.75%
//     3         24        8192      341   29.24%
//    ...        ...         ...      ...      ...
//    67      32768       32768        1    0.00%
```

---

## 分配流程

```go
// mallocgc 是 Go 的内存分配入口
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    // 1. 检查 size
    if size == 0 {
        return unsafe.Pointer(&zerobase)
    }

    // 2. 获取当前 P 的 mcache
    mp := acquirem()
    c := mp.mcache

    // 3. 小对象 (<32KB): 从 mcache 分配
    if size <= maxSmallSize {
        if noscan && size < maxTinySize {
            // Tiny allocator (更小的对象)
            off := c.tinyoffset
            if off+size <= maxTinySize && c.tiny != 0 {
                // 使用 tiny allocator
                x = unsafe.Pointer(c.tiny + off)
                c.tinyoffset = off + size
                c.tinyAllocs++
                return x
            }
            // 重新申请 tiny block
        }

        // 普通小对象
        sizeclass := sizeToClass(size)
        span := c.alloc[sizeclass]
        x = nextFreeFast(span)
        if x == 0 {
            x = c.nextFree(sizeclass)  // 从 mcentral 补充
        }
    } else {
        // 4. 大对象 (>32KB): 直接从 mheap 分配
        span := largeAlloc(size, needzero)
        x = unsafe.Pointer(span.base())
    }

    // 5. 清零（如果需要）
    if needzero {
        memclrNoHeapPointers(x, size)
    }

    return x
}
```

---

## Tiny Allocator

```go
// Tiny allocator 用于 <16B 的对象
// 合并多个小对象到一个 16B block，减少碎片

func (c *mcache) nextTiny(size uintptr) unsafe.Pointer {
    // 尝试在当前 tiny block 分配
    off := c.tinyoffset
    if off+size <= maxTinySize && c.tiny != 0 {
        x := unsafe.Pointer(c.tiny + off)
        c.tinyoffset = off + roundUpSize(size)
        c.tinyAllocs++
        return x
    }

    // 需要新的 tiny block
    span := c.alloc[tinySpanClass]
    v := nextFreeFast(span)
    if v == 0 {
        v = c.nextFree(tinySpanClass)
    }

    c.tiny = uintptr(v)
    c.tinyoffset = size
    return v
}

// 优点：
// 1. 减少 8B/16B 对象的分配开销
// 2. 更好的内存局部性
// 3. 减少 mspan 数量
```

---

## Size Class 计算

```go
// size class 到实际大小的映射
// 不是线性的，而是根据浪费率优化的

func initSizeClasses() {
    // 算法：
    // 1. 从 8B 开始
    // 2. 每次增加，确保浪费率 < 12.5%
    // 3. 对齐到 8 字节（64位系统）

    var sizeclass int
    var size uintptr = 8

    for size <= maxSmallSize {
        // 计算页数
        npages := (size + pageSize - 1) / pageSize

        // 每页可容纳的对象数
        nobj := pageSize * npages / size

        // 浪费率
        waste := float64(pageSize*npages-size*nobj) / float64(pageSize*npages)

        if waste <= maxWaste || sizeclass == 0 {
            class_to_size[sizeclass] = uint32(size)
            class_to_allocnpages[sizeclass] = uint8(npages)
            sizeclass++
        }

        // 增加 size，对齐到 8 字节
        if size < 512 {
            size += 8
        } else if size < 4096 {
            size += 64
        } else {
            size += 512
        }
    }
}
```

---

## 性能优化

| 优化 | 效果 |
|------|------|
| mcache | 无锁分配，90%+ 分配在这里完成 |
| Tiny allocator | 减少小对象开销 50%+ |
| Size classes | 减少碎片，提高缓存命中率 |
| Span reuse | 减少系统调用 |

---

## 参考文献

1. [Go Memory Allocator](https://github.com/golang/go/blob/master/src/runtime/malloc.go)
2. [TCMalloc Design](http://goog-perftools.sourceforge.net/doc/tcmalloc.html)
3. [A Quick Guide to Go's Assembler](https://golang.org/doc/asm) - 了解底层实现
