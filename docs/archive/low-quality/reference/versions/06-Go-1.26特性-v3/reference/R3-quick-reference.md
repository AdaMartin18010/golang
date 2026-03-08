# R3: 快速参考

> **层级**: 参考层 (Reference)
> **地位**: 速查手册

---

## new表达式

```go
// 语法
new(expression)

// 示例
ptr := new(42)              // *int
ptr := new("hello")         // *string
ptr := new(time.Second)     // *time.Duration

// 用于可选字段
type Config struct {
    Timeout *time.Duration
}
cfg := Config{Timeout: new(30 * time.Second)}
```

## 递归泛型

```go
// 约束定义
type Node[T Node[T]] interface {
    Children() []T
}

// 使用
func Walk[T Node[T]](node T, fn func(T)) {
    fn(node)
    for _, child := range node.Children() {
        Walk(child, fn)
    }
}
```

## 定理速查

| 定理 | 内容 | 应用 |
|------|------|------|
| Th1.1 | new(v) ≡ &v' | 可选字段模式 |
| Th1.2 | 递归约束终止 | 树遍历算法 |
| Th2.1 | GC pause < 1ms (99%) | 低延迟优化 |

---

**入口**: [README](../README.md)
