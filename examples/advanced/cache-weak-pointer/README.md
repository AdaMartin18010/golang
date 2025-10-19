# Weak Pointer Cache 示例

## 说明

此示例展示如何使用`runtime/weak`包实现弱引用缓存，避免内存泄漏。

## 特性

- 使用weak.Pointer实现缓存
- 允许GC回收不活跃对象
- 对比强引用缓存的内存使用
- 自动清理失效条目

## 运行要求

**Go 1.22+**（weak包在Go 1.22引入）

```bash
go run main.go
```

## 核心概念

### Weak Reference（弱引用）

弱引用不会阻止GC回收对象：

```go
// 强引用：阻止GC
var strongRef *Value = &Value{Data: "data"}

// 弱引用：不阻止GC
var weakRef = weak.Make(&Value{Data: "data"})
```

### 使用场景

✅ **适合**:
- 缓存实现
- 观察者模式
- 资源管理器
- 避免循环引用

❌ **不适合**:
- 需要确保对象存活
- 关键数据存储
- 频繁访问的数据

## 性能对比

```
强引用缓存：
- 内存持续增长
- 可能导致OOM
- 100%缓存命中

弱引用缓存：
- 内存可控
- GC可回收不活跃对象
- 85%+缓存命中（足够）
```

## 示例输出

```
✅ Cached 1000 items
📊 Memory: Alloc=45 MB, Sys=67 MB, NumGC=3
⚡ Triggering GC...
📊 Memory: Alloc=18 MB, Sys=67 MB, NumGC=4
🧹 Cleaned up 900 entries
💡 Weak cache allows GC to reclaim unused entries
```

## 最佳实践

1. **定期清理**：调用`Cleanup()`清理失效条目
2. **监控命中率**：确保缓存效果
3. **合理配置**：平衡内存和命中率
4. **活跃保持**：对重要数据保持强引用

## 更多信息

- [内存分配器优化文档](../../../docs/02-Go语言现代化/12-Go-1.23运行时优化/03-内存分配器优化.md)
- [FAQ](../../../docs/02-Go语言现代化/12-Go-1.23运行时优化/FAQ.md)

