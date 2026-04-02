---
id: 2026040201
title: F-有界多态性 (F-Bounded Polymorphism)
date: 2026-04-02
tags: [generics, type-theory, go126]
references: [go-blog-2026-02, fgg-paper]
status: active
---

## 定义

F-有界多态性允许类型参数引用自身，实现自指类型约束。

```go
type Adder[A Adder[A]] interface {
    Add(A) A
}
```

## 理论基础

源于类型理论中的 F-有界量化 (F-bounded quantification)，类似 Java 和 Rust 中的模式。

## 应用场景

### 1. 数学运算抽象

```go
type Number[N Number[N]] interface {
    Add(N) N
    Mul(N) N
}
```

### 2. 构建器模式

```go
type Builder[B Builder[B]] interface {
    WithName(string) B
    Build() *Product
}
```

### 3. 比较器

```go
type Ordered[T Ordered[T]] interface {
    Compare(T) int
}
```

## 形式化语义

```
Γ ⊢ A <: Adder[A]
───────────────────────
Γ ⊢ algo[A](x, y) : A
```

## 实现机制

Go 1.26 编译器通过延迟类型检查支持自引用约束。

## 优势

- 表达能力接近 Java/Rust
- 简化数学抽象
- 支持流畅接口

## 限制

- 循环检查更复杂
- 错误信息可能难以理解

## 关联

- [[Go Generics Type System]]
- [[Structural Subtyping]]
- [[Featherweight Generic Go]]

## 待研究

- [ ] 与 Java F-bounds 的详细对比
- [ ] 编译器实现细节
- [ ] 性能影响分析
