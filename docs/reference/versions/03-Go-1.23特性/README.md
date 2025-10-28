# Go 1.23特性

Go 1.23版本特性完整指南，涵盖迭代器、语言改进和标准库更新。

---

## 🎯 核心特性

### 1. 迭代器 (Iterators) ⭐⭐⭐⭐⭐

**新增iter包**:
```go
import "iter"

// Seq: 单值迭代器
func Count(start, end int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := start; i < end; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

// 使用
for v := range Count(1, 10) {
    fmt.Println(v)
}

// Seq2: 键值对迭代器
func Enumerate[T any](slice []T) iter.Seq2[int, T] {
    return func(yield func(int, T) bool) {
        for i, v := range slice {
            if !yield(i, v) {
                return
            }
        }
    }
}

// 使用
for i, v := range Enumerate([]string{"a", "b", "c"}) {
    fmt.Printf("%d: %s\n", i, v)
}
```

### 2. range over func ⭐⭐⭐⭐⭐

```go
// 自定义迭代器
func Fibonacci(n int) func(func(int) bool) {
    return func(yield func(int) bool) {
        a, b := 0, 1
        for i := 0; i < n; i++ {
            if !yield(a) {
                return
            }
            a, b = b, a+b
        }
    }
}

// 使用
for v := range Fibonacci(10) {
    fmt.Println(v)  // 0, 1, 1, 2, 3, 5, 8, 13, 21, 34
}
```

### 3. slog改进 ⭐⭐⭐⭐

**更好的日志处理**:
```go
import "log/slog"

// 结构化日志
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("User action",
    "user_id", 123,
    "action", "login",
    "ip", "192.168.1.1",
)
```

### 4. 语言微调

- `struct`标签改进
- 类型推断增强
- 编译器优化

---

## 📚 详细文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

## 🔗 相关资源

- [Go 1.23发布说明](https://go.dev/doc/go1.23)
- [迭代器提案](https://go.dev/wiki/RangefuncExperiment)
- [版本对比](../00-版本对比与选择指南.md)

---

**发布时间**: 2024年8月  
**最后更新**: 2025-10-28
