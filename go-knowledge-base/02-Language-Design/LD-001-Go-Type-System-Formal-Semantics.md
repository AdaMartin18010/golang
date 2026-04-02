# LD-001: Go 类型系统形式语义 (Go Type System Formal Semantics)

> **维度**: Language Design
> **级别**: S (25+ KB)
> **标签**: #type-system #generics #interface #structural-typing
> **权威来源**: [Go Spec](https://go.dev/ref/spec), [Featherweight Go](https://arxiv.org/abs/2005.11710), [Go Generics Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)

---

## 类型系统概述

Go 采用结构化类型系统（Structural Typing）而非名义类型系统（Nominal Typing）。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go Type System Taxonomy                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Basic Types                                                                  │
│  ───────────                                                                  │
│  bool | string | numeric (int, float, complex) | unsafe.Pointer            │
│                                                                              │
│  Composite Types                                                              │
│  ──────────────                                                               │
│  Array | Slice | Struct | Pointer | Function | Interface | Map | Channel     │
│                                                                              │
│  Generic Types (Go 1.18+)                                                     │
│  ─────────────────────────                                                    │
│  Type Parameters | Type Sets | Constraints                                    │
│                                                                              │
│  Type Identity: Two types are identical if their underlying type trees        │
│  are structurally equivalent.                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 接口实现：结构化类型

### 形式化定义

$$
\begin{aligned}
&\text{Type } T \text{ implements interface } I \text{ iff:} \\
&\forall m \in I.Methods(): \exists m' \in T.Methods() \\
&\quad \text{such that } m.Name = m'.Name \land m.Signature = m'.Signature \\
\\
&\text{Note: No explicit declaration required (implicit satisfaction)}
\end{aligned}
```

### 代码示例

```go
// 接口定义
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 类型定义 - 没有显式声明 implements
type MyReader struct{}

func (r MyReader) Read(p []byte) (n int, err error) {
    // implementation
    return 0, nil
}

// 编译时检查
var _ Reader = MyReader{}  // OK: MyReader implements Reader
```

### 接口内部表示

```go
// src/runtime/iface.go

// 非空接口
type iface struct {
    tab  *itab          // 类型和方法表
    data unsafe.Pointer // 实际数据指针
}

// itab 结构
type itab struct {
    inter *interfacetype  // 接口类型
    _type *_type          // 动态类型
    hash  uint32          // 类型哈希（用于类型 switch）
    _     [4]byte
    fun   [1]uintptr      // 方法表（变长）
}

// 空接口（interface{}）
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

---

## 泛型类型参数（Go 1.18+）

### 类型约束

```go
// 类型集（Type Set）表示法
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64 |
    ~complex64 | ~complex128
}

// 近似元素 ~T 表示所有底层类型为 T 的类型
type MyInt int  // ~int 包含 MyInt

// 约束中可以包含方法
type Ordered interface {
    Integer | Float | ~string
}
```

### 泛型函数实现

```go
// 源码级泛型
func Max[T Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 编译器生成（单态化/Stenciling）
// 实际生成多个实例：
// func Max_int(a, b int) int
// func Max_float64(a, b float64) float64
```

### 类型字典（GC Shape Stenciling）

```go
// 为了减少代码膨胀，Go 使用 GC Shape 分组
// 相同内存布局的类型共享同一个实例

// 相同 shape：指针大小
func f[T any](x T) // 一个实例处理所有指针类型

// 不同 shape：不同大小
func g[T any](x []T) // 每个元素大小不同，需要不同实例
```

---

## 类型推断

```go
// 显式实例化
var m = Max[int](1, 2)

// 类型推断
var n = Max(1, 2)        // T 推断为 int
var f = Max(1.0, 2.0)    // T 推断为 float64

// 约束类型推断
func Map[T, R any](s []T, f func(T) R) []R

// 调用时推断
ints := []int{1, 2, 3}
strings := Map(ints, strconv.Itoa)  // T=int, R=string
```

---

## 协变与逆变

Go 的类型系统：**不支持**协变/逆变（数组、切片、channel 都是不变的）

```go
// 错误：Go 切片是不变的
var animals []Animal = []Cat{}  // Compile error

// 正确：使用泛型
func FeedAll[T Animal](animals []T)  // OK
```

---

## 类型安全证明

Featherweight Go (FG) 证明了 Go 类型系统的安全性：

$$
\text{Progress} + \text{Preservation} = \text{Type Safety}
$$

**Progress**: 良类型的程序不会 stuck（要么继续执行，要么终止）

**Preservation**: 如果 $e : T$ 且 $e \rightarrow e'$，则 $e' : T$

---

## 参考文献

1. [The Go Programming Language Specification](https://go.dev/ref/spec) - 官方规范
2. [Featherweight Go](https://arxiv.org/abs/2005.11710) - Griesemer et al.
3. [Type Parameters Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
4. [Go Generics Implementation](https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md)
