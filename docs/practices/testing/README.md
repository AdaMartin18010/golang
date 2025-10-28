# Go测试实践

Go测试完整指南，涵盖单元测试、集成测试、性能测试和测试最佳实践。

---

## 📚 核心内容

1. **[单元测试](./01-单元测试.md)** ⭐⭐⭐⭐⭐
2. **[表格驱动测试](./02-表格驱动测试.md)** ⭐⭐⭐⭐⭐
3. **[集成测试](./03-集成测试.md)** ⭐⭐⭐⭐
4. **[性能测试](./04-性能测试.md)** ⭐⭐⭐⭐
5. **[测试覆盖率](./05-测试覆盖率.md)** ⭐⭐⭐⭐
6. **[Mock与Stub](./06-Mock与Stub.md)** ⭐⭐⭐⭐
7. **[测试最佳实践](./07-测试最佳实践.md)** ⭐⭐⭐⭐⭐

---

## 🚀 快速开始

### 表格驱动测试
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**最后更新**: 2025-10-28
