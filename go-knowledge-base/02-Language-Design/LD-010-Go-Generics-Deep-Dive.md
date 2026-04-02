# LD-010: Go 泛型深度解析 (Go Generics Deep Dive)

> **维度**: Language Design
> **级别**: S (17+ KB)
> **标签**: #go-generics #type-parameters #constraints #type-inference
> **权威来源**: [Go Generics Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md), [Type Parameters](https://go.dev/tour/generics/1)
> **Go 版本**: 1.18+

---

## 核心概念

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Go Generics Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  类型参数 (Type Parameters)                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  func Map[K comparable, V any](keys []K, f func(K) V) []V           │    │
│  │         └───────┘  └─────┘                                          │    │
│  │         类型参数     约束                                             │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  类型约束 (Constraints)                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  type Number interface {                                            │    │
│  │      ~int | ~int8 | ~int16 | ~int32 | ~int64 |                     │    │
│  │      ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |                 │    │
│  │      ~float32 | ~float64                                           │    │
│  │  }                                                                 │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  类型推导 (Type Inference)                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  result := Map([]int{1,2,3}, func(x int) int { return x*2 })        │    │
│  │  // 编译器自动推导出 K=int, V=int                                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 基础语法

### 函数泛型

```go
package main

import "fmt"

// 非泛型版本 (Go 1.17 及之前)
func IntMin(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func Float64Min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// 泛型版本 (Go 1.18+)
func Min[T comparable](a, b T) T {
    if a < b {  // 错误: comparable 不支持 <
        return a
    }
    return b
}

// 正确的约束
func Min[T interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64 }](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 或使用预定义约束
import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

func main() {
    fmt.Println(Min(1, 2))        // int
    fmt.Println(Min(1.5, 2.5))    // float64
    fmt.Println(Min("a", "b"))    // string (字典序)
}
```

### 结构体泛型

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
intStack := Stack[int]{}
intStack.Push(1)
intStack.Push(2)

strStack := Stack[string]{}
strStack.Push("hello")
```

### 接口泛型

```go
// 泛型接口
type Processor[T, R any] interface {
    Process(input T) (R, error)
}

// 实现
type StringToInt struct{}

func (s StringToInt) Process(input string) (int, error) {
    return strconv.Atoi(input)
}
```

---

## 高级模式

### 类型约束定义

```go
// 基本约束
import "constraints"

type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
    Signed | Unsigned
}

type Float interface {
    ~float32 | ~float64
}

type Ordered interface {
    Integer | Float | ~string
}

// 复杂约束
type ComparableHasher interface {
    comparable
    Hash() uint64
}

// 方法约束
type Stringer interface {
    String() string
}

func ToString[T Stringer](items []T) []string {
    result := make([]string, len(items))
    for i, item := range items {
        result[i] = item.String()
    }
    return result
}
```

### 类型推导

```go
// 显式指定类型参数
result := Min[int](1, 2)

// 类型推导 (更简洁)
result := Min(1, 2)  // 编译器推导出 T=int

// 多类型参数推导
func Convert[A, B any](input A, converter func(A) B) B {
    return converter(input)
}

// 部分推导
result := Convert(42, func(x int) string {
    return strconv.Itoa(x)
})  // 推导 A=int, B=string
```

---

## 常见数据结构实现

### 泛型 Map/Filter/Reduce

```go
package slice

// Map 映射
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

// Filter 过滤
func Filter[T any](s []T, f func(T) bool) []T {
    var result []T
    for _, v := range s {
        if f(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce 归约
func Reduce[T, U any](s []T, initial U, f func(U, T) U) U {
    result := initial
    for _, v := range s {
        result = f(result, v)
    }
    return result
}

// 使用示例
func main() {
    numbers := []int{1, 2, 3, 4, 5}

    // Map: 平方
    squares := Map(numbers, func(n int) int {
        return n * n
    })
    fmt.Println(squares) // [1 4 9 16 25]

    // Filter: 偶数
    evens := Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // [2 4]

    // Reduce: 求和
    sum := Reduce(numbers, 0, func(acc, n int) int {
        return acc + n
    })
    fmt.Println(sum) // 15
}
```

### 泛型缓存

```go
package cache

import "sync"

// Cache 泛型缓存
type Cache[K comparable, V any] struct {
    mu    sync.RWMutex
    items map[K]V
}

func New[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        items: make(map[K]V),
    }
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

func (c *Cache[K, V]) Set(key K, val V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = val
}

// 使用
func main() {
    // 字符串缓存
    strCache := cache.New[int, string]()
    strCache.Set(1, "hello")

    // 结构体缓存
    type User struct {
        Name string
        Age  int
    }
    userCache := cache.New[string, User]()
    userCache.Set("user-1", User{Name: "Alice", Age: 30})
}
```

---

## 性能考量

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Generics Performance                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Go 1.18+ 实现: GC Shape Stenciling + Dictionaries                          │
│                                                                              │
│  1. 编译时: 为每个底层类型 (GC Shape) 生成一个实现                             │
│  2. 运行时: 字典传递类型信息                                                  │
│                                                                              │
│  性能影响:                                                                   │
│  - 与手动实现的非泛型代码性能相当                                              │
│  - 无装箱开销 (与 interface{} 不同)                                          │
│  - 编译时间略有增加                                                           │
│                                                                              │
│  代码膨胀:                                                                   │
│  - 为每个不同的类型参数生成代码                                                │
│  - int, int64 会生成不同代码                                                  │
│  - 但 GC Shape 相同的类型共享代码 (如 *int, *string 共享指针代码)              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 最佳实践

### ✅ 推荐

```go
// 1. 使用预定义约束
import "golang.org/x/exp/constraints"

func Sum[T constraints.Ordered](items []T) T { ... }

// 2. 类型参数命名简洁
func Transform[T, R any](input T, fn func(T) R) R { ... }
// 而不是 func Transform[InputType, ReturnType any]...

// 3. 复杂逻辑使用类型断言
func Equal[T any](a, b T) bool {
    if comparable, ok := any(a).(interface{ Equal(T) bool }); ok {
        return comparable.Equal(b)
    }
    return any(a) == any(b)
}
```

### ❌ 避免

```go
// 1. 过度泛型 (YAGNI)
func process[T any, U any, V any](...)

// 2. 不必要的约束
func Print[T fmt.Stringer](x T)  // 如果只需要 String()，直接用接口

// 3. 类型参数过多
// 通常 1-2 个类型参数足够
```

---

## 参考文献

1. [Go Generics Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
2. [Type Parameters](https://go.dev/tour/generics/1)
3. [Generics in Go](https://go.dev/blog/intro-generics)
