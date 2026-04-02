# LD-003: Go 垃圾回收器三色标记 (Go GC: Tri-Color Mark & Sweep)

> **维度**: Language Design
> **级别**: S (25+ KB)
> **标签**: #go-gc #garbage-collection #tri-color-marking #concurrent-gc
> **权威来源**: [Go GC Guide](https://tip.golang.org/doc/gc-guide), [Go 1.5 GC](https://talks.golang.org/2015/gogc.slide), [Go Runtime Source](https://github.com/golang/go/tree/master/src/runtime)

---

## GC 演进历史

```
Go 1.0-1.4 (2012-2015)      Go 1.5 (2015)               Go 1.8+ (2017+)
      │                            │                            │
      ▼                            ▼                            ▼
┌─────────────┐            ┌───────────────┐            ┌─────────────────┐
│  串行 GC    │───────────►│  并发 GC      │───────────►│  亚毫秒 GC      │
│  Stop-The-  │            │  Tri-Color    │            │  Pacing 算法    │
│  World      │            │  Concurrent   │            │  优化           │
│  ~100ms+    │            │  ~10-50ms     │            │  ~100μs-1ms     │
└─────────────┘            └───────────────┘            └─────────────────┘
```

---

## 三色标记算法

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Tri-Color Marking Algorithm                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Color Meaning:                                                              │
│  ──────────────                                                              │
│  • White: 未访问，可能是垃圾                                                 │
│  • Gray:  已访问，但引用的对象未完全扫描                                       │
│  • Black: 已访问，引用的对象已全部扫描                                         │
│                                                                              │
│  Algorithm:                                                                  │
│  1. 所有对象初始为 White                                                      │
│  2. 从 GC Roots 开始，将直接引用的对象标记为 Gray                               │
│  3. 从 Gray 集合取出一个对象：                                                │
│     a. 将其引用的所有 White 对象标记为 Gray                                     │
│     b. 将该对象标记为 Black                                                    │
│  4. 重复步骤 3 直到 Gray 集合为空                                              │
│  5. 所有 White 对象即为垃圾，可以回收                                          │
│                                                                              │
│  GC Roots:                                                                   │
│  ──────────                                                                  │
│  • 全局变量                                                                  │
│  • 每个 Goroutine 的栈上的变量                                                 │
│  • 寄存器中的指针                                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 并发标记的问题：三色不变性破坏

```
问题场景（写屏障前）：

A (Black) ──► B (White)          用户程序并发执行：
                              1. A 删除对 B 的引用
                              2. A 添加对 C 的引用
C (White) ◄──                 3. B 被标记为垃圾（但 A 还引用 C）

结果：C 被错误地回收（A 还引用它）

解决：写屏障 (Write Barrier)
```

### 写屏障 (Write Barrier)

```go
// 伪代码：写屏障实现

func writePointer(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    //  shade(ptr) 将 ptr 标记为 Gray
    shade(ptr)

    // 原子的写操作
    *slot = ptr
}

// 或 Dijkstra 写屏障（Go 1.8+ 使用）
func writeBarrierDijkstra(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    // 如果当前在 GC 标记阶段
    if gcphase == _GCmark {
        // 将要写入的指针 shade
        shade(ptr)
    }
    *slot = ptr
}
```

---

## Go GC 实现

### GC 阶段

```go
// src/runtime/mgc.go

// GC 周期
const (
    _GCoff = iota  // GC 未运行
    _GCmark        // 标记阶段
    _GCmarktermination // 标记终止
    _GCoff         // 清扫阶段（实际是 _GCoff 的一部分）
)

// GC 触发条件
// 1. 堆大小达到 GOGC 百分比（默认 100%）
// 2. 手动调用 runtime.GC()
// 3. 系统内存压力

// GC 目标：标记阶段 CPU 使用率不超过 25%
// GOGC=100 表示堆大小翻倍时触发 GC
```

### 标记实现

```go
// src/runtime/mgcmark.go

// gcDrain 扫描灰色对象
gcDrain(gcw *gcWork, flags gcDrainFlags) {
    for {
        // 1. 从本地 work buffer 获取对象
        obj := gcw.get()
        if obj == 0 {
            // 2. 从全局队列获取
            obj = gcw.getFromGlobal()
            if obj == 0 {
                break
            }
        }

        // 3. 扫描对象
        scanobject(obj, gcw)
    }
}

// scanobject 扫描对象，将引用的白色对象标记为灰色
func scanobject(obj uintptr, gcw *gcWork) {
    // 获取对象类型信息
    s := spanOf(obj)
    n := s.elemsize

    // 遍历对象的指针字段
    for i := uintptr(0); i < n; i += sys.PtrSize {
        p := *(*uintptr)(unsafe.Pointer(obj + i))

        // 如果是有效指针，标记为灰色
        if p != 0 && inHeap(p) {
            if obj, span, objIndex := findObject(p); obj != 0 {
                greyobject(obj, span, objIndex, gcw)
            }
        }
    }
}
```

### GC Pacing 算法

```
目标：在堆大小翻倍前完成标记

用户分配内存速率: A bytes/ms
标记速率: M bytes/ms (由 dedicates workers 控制)

约束: M >= A (标记速度 >= 分配速度)

GO GC Pacing:
1. 计算目标堆大小: heap_goal = live_data * (1 + GOGC/100)
2. 根据当前 live data 和 heap_goal 计算所需的标记速度
3. 启动适当数量的 dedicated/fractional workers

触发时机:
trigger_ratio = 0.6  // 在达到 heap_goal 的 60% 时开始 GC
```

---

## GC 调优

### 配置参数

```go
// GOGC 环境变量
// 默认 100，表示堆大小翻倍时触发 GC
// GOGC=off 关闭 GC（仅在内存无限时使用）
// GOGC=50 更频繁的 GC，更低的内存占用
// GOGC=200 较少的 GC，更高的内存占用

// GOMEMLIMIT (Go 1.19+)
// 软内存限制，GC 会尽量保持内存在此限制以下
// 超过限制会触发更激进的 GC

// 代码中调整
import "runtime"

// 设置 GOGC 等效值
debug.SetGCPercent(100)

// 设置内存限制 (Go 1.19+)
debug.SetMemoryLimit(10 << 30)  // 10GB

// 强制 GC
runtime.GC()

// 释放内存给 OS (Go 1.13+)
debug.FreeOSMemory()
```

### 性能分析

```go
import (
    "runtime"
    "runtime/trace"
    "os"
)

func main() {
    // 开启 GC trace
    trace.Start(os.Stderr)
    defer trace.Stop()

    // 查看 GC 统计
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    println("Alloc:", m.Alloc)
    println("TotalAlloc:", m.TotalAlloc)
    println("Sys:", m.Sys)
    println("NumGC:", m.NumGC)
    println("PauseNs (last GC pause):", m.PauseNs[(m.NumGC+255)%256])
    println("GCCPUFraction:", m.GCCPUFraction)
}
```

### GC 优化技巧

| 问题 | 症状 | 解决方案 |
|------|------|---------|
| GC 频繁 | 小对象多 | 对象池 sync.Pool |
| GC 耗时 | 大堆 | 减少对象引用，扁平化数据结构 |
| 内存高 | 存活对象多 | 及时释放，避免全局变量 |
| 分配快 | mark 跟不上 | 增加 GOGC，或优化分配模式 |

---

## 与其他 GC 对比

| 特性 | Go GC | Java G1 | .NET GC | V8 | LuaJIT |
|------|-------|---------|---------|-----|--------|
| 算法 | Tri-Color | Region-based | Generational | Generational | Reference Counting |
| STW | < 100μs | ~1-10ms | ~1-10ms | ~1-10ms | None |
| 并发 | Full | Partial | Partial | Partial | N/A |
| 压缩 | No | Yes | Yes | Yes | No |
| 分代 | No | Yes | Yes | Yes | N/A |

---

## 参考文献

1. [Go GC Guide](https://tip.golang.org/doc/gc-guide) - 官方 GC 指南
2. [Go 1.5 GC](https://talks.golang.org/2015/gogc.slide) - Austin Clements
3. [Tracing Garbage Collection](https://www.memorymanagement.org/mmref/recycle.html#tracing-garbage-collection) - Memory Management Reference
4. [Garbage Collection Handbook](http://gchandbook.org/) - Jones & Lins
