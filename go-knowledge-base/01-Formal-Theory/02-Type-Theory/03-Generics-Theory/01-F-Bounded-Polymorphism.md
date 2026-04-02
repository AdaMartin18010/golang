# F-有界多态性 (F-Bounded Polymorphism)

> **分类**: 形式理论
> **难度**: 专家
> **前置知识**: 泛型基础

---

## 概述

F-有界多态性允许类型参数引用自身，实现自指类型约束。Go 1.26 引入此特性。

---

## 语法

```go
type Adder[A Adder[A]] interface {
    Add(A) A
}
```

**类型参数 `A` 的约束是 `Adder[A]` 本身**。

---

## 形式化

```
type F[T <: F[T]] interface { ... }

递归约束: T 必须是 F[T] 的子类型
```

---

## 应用

### 数学运算

```go
type Number[N Number[N]] interface {
    Add(N) N
    Mul(N) N
}
```

### 构建器模式

```go
type Builder[B Builder[B]] interface {
    WithName(string) B
    Build() *Product
}
```

---

## 类型检查

```
检查: Int 是否满足 Number[Int]

1. Number[Int] = interface { Add(Int) Int; Mul(Int) Int }
2. Int 实现了所有方法
3. 结论: Int <: Number[Int] ✅
```

---

## 参考

- Go 1.26 Release Notes
- "F-bounded polymorphism" (Cardelli et al.)
