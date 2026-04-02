# 简洁性原则 (Simplicity)

> **分类**: 语言设计
> **难度**: 入门

---

## 核心思想

**"Less is More"** - 用更少的方式做更多的事。

Go 的设计哲学：

- 一种方式做一件事
- 显式优于隐式
- 简单优于复杂

---

## 设计体现

### 1. 错误处理

```go
// Go: 显式错误处理
if err != nil {
    return err
}

// 对比 Java: 异常
// 对比 Rust: ? 运算符
```

### 2. 继承

```go
// Go: 无继承，只有组合
type Reader struct { }
type Writer struct { }

type ReadWriter struct {
    Reader
    Writer
}
```

### 3. 泛型限制

```go
// Go 1.18: 简化泛型
func Max[T Ordered](a, b T) T

// 对比 C++: 模板元编程
// 对比 Rust: 复杂 trait 系统
```

---

## 代价

### 优点

- 易于学习
- 易于阅读
- 易于维护

### 缺点

- 代码重复（无泛型时）
- 错误处理冗长
- 缺少某些抽象

---

## 引用

> "The complexity of the language must not exceed the complexity of the problem." - Rob Pike
