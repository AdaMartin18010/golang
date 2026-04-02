# LD-010: Go 泛型的形式化理论 (Go Generics: Formal Theory)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #generics #type-parameters #constraints #contracts #go118
> **权威来源**:
>
> - [Go Generics Proposal](https://go.googlesource.com/proposal/+/HEAD/design/43651-type-parameters.md) - Ian Lance Taylor
> - [Type Parameters](https://go.dev/tour/generics/1) - Go Authors
> - [Parameterized Types](https://go.dev/doc/tutorial/generics) - Go Tutorial

---

## 1. 泛型基础

### 1.1 类型参数

**定义 1.1 (类型参数)**
类型参数是类型的占位符，在实例化时替换为具体类型。

```go
// 泛型函数
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 泛型类型
type Stack[T any] struct {
    items []T
}
```

### 1.2 约束

**定义 1.2 (约束)**
约束定义了类型参数必须满足的条件。

```go
// any 约束 - 允许任何类型
func Print[T any](v T) {
    fmt.Println(v)
}

// Ordered 约束 - 可比较排序
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 自定义约束
type Number interface {
    ~int | ~int64 | ~float64
}
```

---

## 2. 形式化定义

### 2.1 类型参数语法

```
TypeParameters = "[" TypeParamDecl {"," TypeParamDecl} "]"
TypeParamDecl  = Identifier TypeConstraint
TypeConstraint = InterfaceType
```

### 2.2 类型推导

**定义 2.1 (类型推导)**
编译器可从函数参数推导出类型参数。

```go
// 显式指定
min := Min[int](1, 2)

// 类型推导
min := Min(1, 2)  // T 推导为 int
```

**定理 2.1 (推导完备性)**
若类型参数可从参数推导，则显式指定是可选的。

---

## 3. 约束详解

### 3.1 预定义约束

| 约束 | 要求 | 适用操作 |
|------|------|----------|
| any | 无 | 任何操作 |
| comparable | ==, != | 相等比较 |
| constraints.Ordered | <, <=, >, >= | 排序 |
| constraints.Signed | 有符号整数 | 算术 |
| constraints.Unsigned | 无符号整数 | 算术 |
| constraints.Integer | 整数 | 算术 |
| constraints.Float | 浮点 | 算术 |

### 3.2 类型集

```go
// 类型集约束
type Number interface {
    ~int | ~int64 | ~float64 | ~float32
}

// ~ 表示底层类型
// type MyInt int 满足 ~int
```

### 3.3 方法约束

```go
type Stringer interface {
    String() string
}

func ToString[T Stringer](v T) string {
    return v.String()
}
```

---

## 4. 泛型实现

### 4.1 GCShape  stencil

Go 使用 GCShape 和 stencil 实现泛型：

```
相同 GCShape 的类型共享代码:
- 指针类型: 共享一个实现
- 整数类型: 分别编译
```

### 4.2 类型字典

运行时通过类型字典传递类型信息：

```go
// 编译器生成
type dict struct {
    typeInfo *typeInfo
    methods  []func()
}
```

---

## 5. 泛型模式

### 5.1 泛型数据结构

```go
// 链表
type List[T any] struct {
    head *Node[T]
    tail *Node[T]
    len  int
}

type Node[T any] struct {
    value T
    next  *Node[T]
}

func (l *List[T]) Push(v T) {
    n := &Node[T]{value: v}
    // ...
}
```

### 5.2 泛型算法

```go
// 过滤
func Filter[T any](s []T, fn func(T) bool) []T {
    var result []T
    for _, v := range s {
        if fn(v) {
            result = append(result, v)
        }
    }
    return result
}

// 映射
func Map[T, R any](s []T, fn func(T) R) []R {
    result := make([]R, len(s))
    for i, v := range s {
        result[i] = fn(v)
    }
    return result
}

// 归约
func Reduce[T, R any](s []T, initial R, fn func(R, T) R) R {
    result := initial
    for _, v := range s {
        result = fn(result, v)
    }
    return result
}
```

---

## 6. 性能考虑

### 6.1 编译时开销

- 泛型增加编译时间
- 每个不同的类型参数组合生成代码

### 6.2 运行时开销

- 与手写代码性能相当
- 无装箱开销
- 内联优化有效

---

## 7. 代码示例

### 7.1 泛型 Map

```go
type Map[K comparable, V any] struct {
    data map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
    return &Map[K, V]{data: make(map[K]V)}
}

func (m *Map[K, V]) Get(k K) (V, bool) {
    v, ok := m.data[k]
    return v, ok
}

func (m *Map[K, V]) Set(k K, v V) {
    m.data[k] = v
}

// 使用
m := NewMap[string, int]()
m.Set("foo", 42)
v, ok := m.Get("foo")
```

---

## 8. 关系网络

```
Go Generics
├── Type Parameters
│   ├── Functions
│   ├── Types
│   └── Methods
├── Constraints
│   ├── Any
│   ├── Comparable
│   ├── Ordered
│   └── Custom
└── Implementation
    ├── GCShape
    ├── Stencil
    └── Type Dictionary
```

---

**质量评级**: S (15KB)
**完成日期**: 2026-04-02
