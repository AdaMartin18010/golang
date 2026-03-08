# A3: 性能优化

> **层级**: 应用层 (Application)
> **地位**: 基于 Th2.1 的性能优化实践
> **依赖**: A2

---

## 优化 1: GC调优

### 基于 Th2.1 的配置

```
由 Th2.1: P(pause < 1ms) ≥ 0.99

优化目标:
  - 保持低延迟保证
  - 最小化GC开销
  - 适应应用负载
```

### 配置

```go
package main

import (
    "runtime"
    "runtime/debug"
)

func init() {
    // 低延迟服务配置
    // 更频繁的GC，但每次工作量更少
    debug.SetGCPercent(50)

    // 内存限制（容器环境）
    debug.SetMemoryLimit(6 << 30)  // 6GB
}

// 或环境变量配置
// GOGC=50
// GOMEMLIMIT=6GiB
```

---

## 优化 2: 减少堆分配

```go
// ❌ 每次调用都分配
func process(data []byte) []byte {
    result := make([]byte, len(data))  // 堆分配
    // 处理...
    return result
}

// ✅ 使用对象池
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func processOptimized(data []byte) []byte {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)

    // 使用buf处理
    n := copy(buf, data)
    // 处理buf[:n]...

    result := make([]byte, n)
    copy(result, buf[:n])
    return result
}
```

---

## 优化 3: 栈分配优化

```go
// ✅ 小数组栈分配
func fastProcess() {
    var buf [1024]byte  // 大概率栈分配
    use(buf[:])
}

// ✅ 避免不必要的指针
func goodSum(data []int) int {
    sum := 0
    for _, v := range data {
        sum += v
    }
    return sum
}

// ❌ 不必要的指针
func badSum(data []int) *int {
    sum := 0
    for _, v := range data {
        sum += v
    }
    return &sum  // 强制堆分配
}
```

---

**下一章**: [A4-迁移指南](A4-migration.md)
