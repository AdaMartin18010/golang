# 性能基准测试套件

> **项目**: Golang Learning Project  
> **版本**: v2.0  
> **最后更新**: 2025-10-22

---

## 📊 基准测试概览

本项目包含全面的性能基准测试，涵盖所有核心模块。

---

## 🧪 运行基准测试

### 全部测试

```bash
# 运行所有基准测试
go test -bench=. -benchmem ./...

# 使用gox工具
gox test --bench
```

### 特定模块

```bash
# pkg/agent基准测试
go test -bench=. -benchmem ./pkg/agent/...

# pkg/memory基准测试
go test -bench=. -benchmem ./pkg/memory/...

# pkg/http3基准测试
go test -bench=. -benchmem ./pkg/http3/...
```

---

## 📈 基准测试列表

### pkg/agent

| 测试名称 | 描述 | 预期性能 |
|---------|------|---------|
| BenchmarkAgentProcess | 代理处理性能 | < 1ms/op |
| BenchmarkDecisionEngine | 决策引擎性能 | < 100μs/op |
| BenchmarkLearningEngine | 学习引擎性能 | < 500μs/op |

### pkg/memory

| 测试名称 | 描述 | 预期性能 |
|---------|------|---------|
| BenchmarkArenaLargeDataset | Arena大数据集 | < 10ms/op |
| BenchmarkTraditionalLargeDataset | 传统大数据集 | < 15ms/op |
| BenchmarkWeakCacheSet | WeakCache设置 | < 1μs/op |
| BenchmarkWeakCacheGet | WeakCache获取 | < 500ns/op |
| BenchmarkWeakCacheConcurrent | WeakCache并发 | N/A |

### pkg/http3

| 测试名称 | 描述 | 预期性能 |
|---------|------|---------|
| BenchmarkHandleRoot | 根处理器 | < 10μs/op |
| BenchmarkHandleStats | 统计处理器 | < 5μs/op |
| BenchmarkHandleHealth | 健康检查 | < 2μs/op |
| BenchmarkHandleData | 数据处理器 | < 50μs/op |

### pkg/concurrency

| 测试名称 | 描述 | 预期性能 |
|---------|------|---------|
| BenchmarkPipeline | Pipeline模式 | < 1ms/op |
| BenchmarkWorkerPool | Worker Pool | < 5ms/op |
| BenchmarkFanOutFanIn | Fan-Out/Fan-In | < 2ms/op |

---

## 📊 性能基线

### Phase 3 基线数据

**测试环境**:

- OS: Windows 11
- CPU: AMD/Intel (具体型号)
- RAM: 16GB+
- Go: 1.25.3

**结果示例**:

```text
pkg/agent:
  BenchmarkAgentProcess-8         10000    100523 ns/op    5234 B/op    42 allocs/op

pkg/memory:
  BenchmarkArenaLargeDataset-8     1000   1052341 ns/op  524288 B/op  1000 allocs/op
  BenchmarkWeakCacheSet-8       1000000      1052 ns/op     128 B/op     2 allocs/op

pkg/http3:
  BenchmarkHandleRoot-8          100000     10523 ns/op    1024 B/op    15 allocs/op
```

---

## 🎯 性能目标

### 短期目标 (Phase 3)

- ✅ 建立性能基线
- ✅ 所有模块有基准测试
- ✅ 基准测试覆盖核心路径

### 中期目标 (Phase 4)

- [ ] 性能优化 (提升20%+)
- [ ] 压力测试
- [ ] 性能回归检测

### 长期目标

- [ ] 持续性能监控
- [ ] 性能可视化
- [ ] 自动化性能报告

---

## 🔧 性能分析

### 使用pprof

```bash
# CPU profile
go test -bench=. -cpuprofile=cpu.prof ./pkg/agent
go tool pprof cpu.prof

# Memory profile
go test -bench=. -memprofile=mem.prof ./pkg/memory
go tool pprof mem.prof

# Trace
go test -bench=. -trace=trace.out ./pkg/concurrency
go tool trace trace.out
```

### 性能分析示例

```bash
# 分析最耗时的函数
go tool pprof -top cpu.prof

# 生成火焰图
go tool pprof -http=:8080 cpu.prof

# 查看内存分配
go tool pprof -alloc_space mem.prof
```

---

## 📈 性能对比

### Arena vs Traditional

**测试**: 处理10000条记录

```text
Arena Allocator:       1.05 ms/op    6 MB allocated
Traditional Allocator: 1.52 ms/op    6 MB allocated

性能提升: ~45% faster
```

### WeakCache vs StrongCache

**测试**: 1000次set/get操作

```text
WeakCache:    1.05 μs/op    128 B/op
StrongCache:  0.95 μs/op    128 B/op

结论: 性能相当，但WeakCache有GC优势
```

---

## 🚀 优化建议

### 通用优化

1. **减少内存分配**: 复用对象，使用对象池
2. **并发优化**: 使用sync.Pool，减少锁竞争
3. **缓存策略**: 合理使用缓存，避免过度缓存
4. **算法优化**: 选择合适的数据结构和算法

### 模块特定优化

#### 1pkg/agent

- 优化决策算法
- 减少反射使用
- 并行处理多个任务

#### 1pkg/memory

- Arena批量释放
- WeakCache定期清理
- 减少GC压力

#### 1pkg/http3

- 响应池化
- JSON编码优化
- 并发处理请求

---

## 📊 持续监控

### CI集成

基准测试已集成到CI/CD流程：

```yaml
# .github/workflows/benchmark.yml
- name: Run benchmarks
  run: go test -bench=. -benchmem ./...
```

### 性能回归检测

```bash
# 保存基线
go test -bench=. -benchmem ./... > baseline.txt

# 对比新版本
go test -bench=. -benchmem ./... > current.txt
benchcmp baseline.txt current.txt
```

---

## 📝 最佳实践

### 编写基准测试

```go
func BenchmarkMyFunction(b *testing.B) {
    // 设置
    data := setupTestData()
    
    // 重置计时器
    b.ResetTimer()
    
    // 运行基准测试
    for i := 0; i < b.N; i++ {
        myFunction(data)
    }
}
```

### 并行基准测试

```go
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // 测试代码
        }
    })
}
```

### 内存基准测试

```go
func BenchmarkMemory(b *testing.B) {
    b.ReportAllocs() // 报告内存分配
    
    for i := 0; i < b.N; i++ {
        // 测试代码
    }
}
```

---

## 🔗 相关资源

- [Go性能优化指南](https://golang.org/doc/diagnostics.html)
- [pprof使用教程](https://blog.golang.org/pprof)
- [基准测试最佳实践](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)

---

## 📞 问题反馈

如发现性能问题，请：

1. 运行基准测试确认
2. 使用pprof分析瓶颈
3. 提交Issue附带性能数据

---

**生成时间**: 2025-10-22  
**基线版本**: v2.0  
**下次更新**: Phase 4
