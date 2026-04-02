# 泛型 (Generics)

> **分类**: 语言设计
> **适用版本**: Go 1.18+

---

## 语法

```go
// 类型参数
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 类型约束
type Number interface {
    ~int | ~int64 | ~float64
}

func Sum[T Number](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}
```

---

## 类型集

```go
// 基本类型集
type Integer interface {
    int | int8 | int16 | int32 | int64
}

// 近似类型 (~)
type MyInt int  // 底层类型是 int

type AnyInt interface {
    ~int  // 包含 MyInt
}

// 复合约束
type SignedInteger interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

---

## 泛型类型

```go
// 泛型结构体
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// 使用
stack := Stack[int]{}
stack.Push(42)
```

---

## 类型推断

```go
// 显式指定
min := Min[int](1, 2)

// 推断
min := Min(1, 2)  // T 推断为 int
```

---

## 约束包

```go
import "golang.org/x/exp/constraints"

type Ordered interface {
    Integer | Float | ~string
}
```

---

## 最佳实践

### 1. 从具体开始

```go
// 先写具体实现
func MinInt(a, b int) int

// 再泛化
func Min[T constraints.Ordered](a, b T) T
```

### 2. 约束最小化

```go
// 好: 最小约束
func Print[T fmt.Stringer](v T)

// 不好: 过度约束
func Print[T any](v T)  // 如果只需要 String()
```

### 3. 避免过度泛型

```go
// 不好: 不必要的泛型
func Print[T any](v T) { fmt.Println(v) }

// 好: 直接用 interface{}
func Print(v interface{}) { fmt.Println(v) }
```
