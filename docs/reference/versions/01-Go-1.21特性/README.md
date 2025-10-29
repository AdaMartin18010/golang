# Go 1.21特性

Go 1.21版本特性完整指南，涵盖新增功能、性能优化和标准库更新。

---

## 🎯 核心特性

### 1. 内置函数增强 ⭐⭐⭐⭐⭐

**min/max/clear函数**:

```go
// min/max: 返回最小/最大值
x := min(1, 2, 3)  // 1
y := max(1, 2, 3)  // 3

// clear: 清空map/slice
m := map[string]int{"a": 1, "b": 2}
clear(m)  // m现在为空

s := []int{1, 2, 3}
clear(s)  // s所有元素置为零值
```

### 2. PGO (Profile-Guided Optimization) ⭐⭐⭐⭐⭐

**性能优化**:

```bash
# 1. 收集性能数据
go test -cpuprofile=cpu.pprof

# 2. 使用PGO构建
go build -pgo=cpu.pprof
```

**性能提升**: 2-14%

### 3. GC改进 ⭐⭐⭐⭐

- 尾延迟降低40%
- 内存开销降低
- 更平滑的GC表现

### 4. 标准库更新

**新增包**:

- `log/slog`: 结构化日志
- `cmp`: 比较函数

**slog示例**:

```go
import "log/slog"

slog.Info("User login", "user_id", 123, "ip", "192.168.1.1")
slog.Error("Database error", "error", err)
```

---

## 📚 详细文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

## 🔗 相关资源

- [Go 1.21发布说明](https://go.dev/doc/go1.21)
- [版本对比](../00-版本对比与选择指南.md)

---

**发布时间**: 2023年8月  
**最后更新**: 2025-10-29
