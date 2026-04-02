# Go 泛型深度解析

> **创建**: 2026-04-02
> **状态**: 持续更新
> **关联**: [[F-有界多态性]], [[Featherweight Generic Go]]

---

## 概述

Go 1.18 引入泛型，Go 1.26 进一步增强 (F-有界多态性)。

---

## 类型参数基础

### 声明语法

```go
// 函数泛型
func Max[T Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 类型泛型
type Stack[T any] struct {
    items []T
}
```

### 类型约束

```go
// 接口约束
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64 |
    ~string
}
```

---

## 实现机制

### 字典传递 (Dictionary Passing)

编译时生成类型字典：

```go
// 伪代码
func Max(dict *TypeDict, a, b interface{}) interface{}
```

### 单态化 (Monomorphization)

为常用类型生成专用代码：

```go
// 编译器生成
func Max_int(a, b int) int
func Max_float64(a, b float64) float64
```

### 混合策略

Go 编译器使用混合策略：

- 常用类型：单态化
- 其他类型：字典传递

---

## Go 1.26 增强

### F-有界多态性

```go
type Adder[A Adder[A]] interface {
    Add(A) A
}
```

**应用场景**:

- 数学抽象
- 构建器模式
- 流畅接口

### 递归类型约束

```go
type Node[T Node[T]] interface {
    Parent() T
    Children() []T
}
```

---

## 性能特征

| 特性 | 开销 | 说明 |
|------|------|------|
| 单态化调用 | 0 | 直接调用 |
| 字典传递 | ~1-2ns | 间接调用 |
| 类型推断 | 编译时 | 无运行时开销 |

---

## 最佳实践

### Do ✅

```go
// 使用类型约束
func Sum[T ~int | ~float64](nums []T) T

// 利用类型推断
result := Sum([]int{1, 2, 3})  // 无需显式指定 T
```

### Don't ❌

```go
// 过度泛化
func Process[T any](x T)  // 不如使用 interface{}

// 复杂嵌套约束
// 保持约束简单可理解
```

---

## 与接口对比

| 特性 | 泛型 | 接口 |
|------|------|------|
| 类型检查 | 编译时 | 运行时 |
| 性能 | 更好 | 有间接开销 |
| 灵活性 | 类型安全 | 更灵活 |
| 代码大小 | 可能膨胀 | 固定 |

---

## 深入阅读

- [[Featherweight Generic Go]]
- [[Dictionary Passing Translation]]
- [[Type Set Semantics]]
