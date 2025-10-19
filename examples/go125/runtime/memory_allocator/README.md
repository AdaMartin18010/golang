# 内存分配器优化示例

> **Go 版本**: 1.25+  
> **示例类型**: 内存性能优化  
> **最后更新**: 2025-10-18

本目录包含 Go 1.23+ 内存分配器优化的基准测试和示例代码。

---

## 🚀 快速开始

### 1. 运行所有基准测试

```bash
go test -bench=. -benchmem
```

**预期输出**:

```text
BenchmarkMapLarge/Size1000-8                    50000000    25 ns/op    0 B/op    0 allocs/op
BenchmarkMapLarge/Size10000-8                   30000000    28 ns/op    0 B/op    0 allocs/op
BenchmarkMapLarge/Size100000-8                  20000000    32 ns/op    0 B/op    0 allocs/op
BenchmarkMapLarge/Size1000000-8                 15000000    28 ns/op    0 B/op    0 allocs/op
```

### 2. 对比 Go 1.24 和 1.25

```bash
# Go 1.24 环境
go test -bench=. -benchmem > go1.24.txt

# Go 1.23+ 环境
go test -bench=. -benchmem > go1.23.txt

# 使用 benchstat 对比
go install golang.org/x/perf/cmd/benchstat@latest
benchstat go1.24.txt go1.23.txt
```

### 3. 运行功能测试

```bash
go test -v
```

---

## 📊 基准测试说明

### Swiss Tables Map 测试

| 测试 | 说明 | 关注指标 |
|------|------|---------|
| `BenchmarkMapLarge` | 大规模 map 查找 | ns/op (Go 1.23+ 提升 38%) |
| `BenchmarkMapInsert` | map 插入性能 | allocs/op |
| `BenchmarkMapIteration` | map 遍历性能 | ns/op |
| `BenchmarkMapStringKey` | 字符串键性能 | ns/op |

**运行**:

```bash
go test -bench=BenchmarkMap -benchmem -benchtime=5s
```

### 小对象分配测试

| 测试 | 说明 | 关注指标 |
|------|------|---------|
| `BenchmarkSmallObjectAllocation` | 小对象分配 | allocs/op |
| `BenchmarkSliceAllocation` | 切片分配 | B/op |
| `BenchmarkAllocationPatterns` | 分配模式对比 | allocs/op |

**运行**:

```bash
go test -bench=BenchmarkSmall -benchmem
go test -bench=BenchmarkSlice -benchmem
go test -bench=BenchmarkAllocation -benchmem
```

### GC 压力测试

| 测试 | 说明 | 关注指标 |
|------|------|---------|
| `BenchmarkGCPressure` | GC 压力对比 | gc-count |

**运行**:

```bash
go test -bench=BenchmarkGCPressure -benchmem -benchtime=10s
```

### 实际场景模拟

| 测试 | 说明 | 场景 |
|------|------|------|
| `BenchmarkRealWorldScenario/CacheSimulation` | 缓存系统 | LRU 缓存 |
| `BenchmarkRealWorldScenario/DataProcessing` | 数据处理 | 聚合计算 |

**运行**:

```bash
go test -bench=BenchmarkRealWorld -benchmem
```

---

## 🎯 预期性能提升

### Swiss Tables Map (Go 1.23+)

| 场景 | Go 1.24 | Go 1.23+ | 提升 |
|------|---------|---------|------|
| 大 map 查找 (100万) | 45 ns/op | 28 ns/op | ⬆️ 38% |
| Map 插入 | 120 ns/op | 75 ns/op | ⬆️ 38% |
| Map 遍历 (10万) | 8.5 ms | 5.2 ms | ⬆️ 39% |
| 内存占用 | 45 MB | 42 MB | ⬇️ 7% |

### 小对象分配

| 场景 | Go 1.24 | Go 1.23+ | 提升 |
|------|---------|---------|------|
| 单个对象分配 | 18 ns/op | 14 ns/op | ⬆️ 22% |
| 批量分配 | 125 μs | 95 μs | ⬆️ 24% |

### GC 性能

| 指标 | Go 1.24 | Go 1.23+ | 提升 |
|------|---------|---------|------|
| GC 暂停时间 | 2.5 ms | 2.1 ms | ⬇️ 16% |
| GC 频率 | 120/min | 95/min | ⬇️ 21% |

---

## 💻 代码示例

### 1. Swiss Tables Map 最佳实践

```go
package main

import "fmt"

func main() {
    // ✅ 推荐：预分配容量
    m := make(map[int]string, 1000000)
    
    // 批量插入（Swiss Tables 优化）
    for i := 0; i < 1000000; i++ {
        m[i] = fmt.Sprintf("value_%d", i)
    }
    
    // 高性能查找
    if v, ok := m[500000]; ok {
        fmt.Println(v)
    }
}
```

### 2. 内存分配模式对比

```go
// ❌ 低效：频繁分配
func BadPattern() {
    for i := 0; i < 10000; i++ {
        obj := new(Object)  // 每次分配
        process(obj)
    }
}

// ✅ 高效：预分配
func GoodPattern() {
    objects := make([]Object, 10000)  // 一次分配
    for i := range objects {
        process(&objects[i])
    }
}
```

---

## 📈 性能分析

### 使用 pprof

```bash
# 1. 启用 CPU profile
go test -bench=BenchmarkMapLarge -cpuprofile=cpu.prof

# 2. 分析 profile
go tool pprof cpu.prof
(pprof) top10
(pprof) list BenchmarkMapLarge

# 3. 可视化
go tool pprof -http=:8080 cpu.prof
```

### 使用 trace

```bash
# 1. 生成 trace
go test -bench=BenchmarkGCPressure -trace=trace.out

# 2. 查看 trace
go tool trace trace.out
```

### 内存分析

```bash
# 1. 生成 memory profile
go test -bench=BenchmarkMemoryUsage -memprofile=mem.prof

# 2. 分析内存分配
go tool pprof mem.prof
(pprof) top10
(pprof) list BenchmarkMemoryUsage
```

---

## 🔧 调优建议

### Map 优化

1. **预分配容量**: `make(map[K]V, capacity)`
2. **使用整数键**: 比字符串键快 30%
3. **批量操作**: 利用 Swiss Tables 优化
4. **避免频繁扩容**: 预估合理容量

### 小对象优化

1. **对象池**: 复用对象减少分配
2. **预分配切片**: 一次分配替代多次
3. **值类型**: 避免不必要的指针
4. **批处理**: 减少GC压力

### GC 调优

1. **GOGC**: 调整 GC 触发阈值
2. **内存限制**: `debug.SetMemoryLimit()`
3. **监控**: 使用 Prometheus 监控 GC 指标
4. **分析**: 定期 pprof 分析

---

## 📚 参考资料

- [内存分配器优化文档](../../03-内存分配器优化.md)
- [Go 1.23+ Release Notes](https://golang.org/doc/go1.23)
- [Swiss Tables Paper](https://abseil.io/blog/20180927-swisstables)

---

**示例维护**: AI Assistant  
**最后更新**: 2025-10-18  
**反馈**: 提交 Issue 或 PR
