# Go 1.26.1 高级特性全面分析

## 目录

- [Go 1.26.1 高级特性全面分析](#go-1261-高级特性全面分析)
  - [目录](#目录)
  - [概述](#概述)
  - [泛型编程](#泛型编程)
    - [1. 概念定义](#1-概念定义)
    - [2. 属性特征](#2-属性特征)
    - [3. 关系依赖](#3-关系依赖)
    - [4. 详细示例代码](#4-详细示例代码)
      - [4.1 基础泛型函数](#41-基础泛型函数)
      - [4.2 泛型数据结构](#42-泛型数据结构)
      - [4.3 类型约束详解](#43-类型约束详解)
      - [4.4 泛型接口](#44-泛型接口)
    - [5. 反例说明](#5-反例说明)
  - [接口类型](#接口类型)
    - [1. 概念定义](#1-概念定义-1)
    - [2. 属性特征](#2-属性特征-1)
    - [3. 关系依赖](#3-关系依赖-1)
    - [4. 详细示例代码](#4-详细示例代码-1)
      - [4.1 基本接口](#41-基本接口)
      - [4.2 接口组合](#42-接口组合)
      - [4.3 泛型接口与约束](#43-泛型接口与约束)
    - [5. 反例说明](#5-反例说明-1)
  - [反射机制](#反射机制)
    - [1. 概念定义](#1-概念定义-2)
    - [2. 属性特征](#2-属性特征-2)
    - [3. 关系依赖](#3-关系依赖-2)
    - [4. 详细示例代码](#4-详细示例代码-2)
      - [4.1 基础反射操作](#41-基础反射操作)
      - [4.2 高级反射应用](#42-高级反射应用)
    - [5. 反例说明](#5-反例说明-2)
  - [元编程](#元编程)
    - [1. 概念定义](#1-概念定义-3)
    - [2. 属性特征](#2-属性特征-3)
    - [3. 关系依赖](#3-关系依赖-3)
    - [4. 详细示例代码](#4-详细示例代码-3)
      - [4.1 go generate 基础使用](#41-go-generate-基础使用)
      - [4.2 自定义代码生成工具](#42-自定义代码生成工具)
      - [4.3 使用 text/template 生成代码](#43-使用-texttemplate-生成代码)
      - [4.4 AST 操作示例](#44-ast-操作示例)
    - [5. 反例说明](#5-反例说明-3)
  - [Go 1.26.1 新特性深度分析](#go-1261-新特性深度分析)
    - [1. new() 表达式支持](#1-new-表达式支持)
      - [1.1 概念定义](#11-概念定义)
      - [1.2 属性特征](#12-属性特征)
      - [1.3 详细示例代码](#13-详细示例代码)
      - [1.4 反例说明](#14-反例说明)
    - [2. 自引用泛型约束](#2-自引用泛型约束)
      - [2.1 概念定义](#21-概念定义)
      - [2.2 属性特征](#22-属性特征)
      - [2.3 关系依赖](#23-关系依赖)
      - [2.4 详细示例代码](#24-详细示例代码)
      - [2.5 反例说明](#25-反例说明)
  - [CGO 和外部函数接口](#cgo-和外部函数接口)
    - [1. 概念定义](#1-概念定义-4)
    - [2. 属性特征](#2-属性特征-4)
    - [3. 关系依赖](#3-关系依赖-4)
    - [4. 详细示例代码](#4-详细示例代码-4)
      - [4.1 基础 CGO 使用](#41-基础-cgo-使用)
      - [4.2 内存管理和字符串处理](#42-内存管理和字符串处理)
      - [4.3 Go 函数导出给 C](#43-go-函数导出给-c)
      - [4.4 使用外部 C 库](#44-使用外部-c-库)
      - [4.5 复杂类型和结构体](#45-复杂类型和结构体)
    - [5. 反例说明](#5-反例说明-4)
  - [决策树图分析](#决策树图分析)
    - [1. 泛型使用决策树](#1-泛型使用决策树)
    - [2. 接口设计决策树](#2-接口设计决策树)
    - [3. 反射使用决策树](#3-反射使用决策树)
    - [4. CGO 使用决策树](#4-cgo-使用决策树)
  - [总结](#总结)
    - [特性对比表](#特性对比表)
    - [最佳实践建议](#最佳实践建议)

---

## 概述

Go 1.26.1 是 Go 语言的重要版本更新，带来了多项高级特性的增强和改进：

| 特性 | 版本引入 | 重要性 | 适用场景 |
|------|---------|--------|----------|
| 泛型编程 | Go 1.18 | 高 | 通用数据结构、算法复用 |
| 自引用泛型约束 | Go 1.26 | 极高 | F-bounded 多态、复杂接口 |
| new() 表达式 | Go 1.26 | 高 | 指针初始化、可选字段 |
| 接口类型 | Go 1.0+ | 高 | 多态、解耦设计 |
| 反射机制 | Go 1.0+ | 中 | 动态类型处理、序列化 |
| 元编程 | Go 1.4+ | 中 | 代码生成、自动化 |
| CGO | Go 1.0+ | 中 | C 语言互操作 |

---

## 泛型编程

### 1. 概念定义

泛型编程（Generics）允许编写类型参数化的代码，实现算法和数据结构的复用，而无需牺牲类型安全。

### 2. 属性特征

```
┌─────────────────────────────────────────────────────────────┐
│                    泛型编程核心要素                          │
├─────────────────────────────────────────────────────────────┤
│  类型参数(Type Parameters)  │  类型约束(Type Constraints)    │
│  ├─ 形式：T, K, V          │  ├─ 内置：any, comparable      │
│  ├─ 多参数：[K, V]         │  ├─ 接口约束：interface{}      │
│  └─ 约束限定：T Constraint │  └─ 类型集：~int \| ~float64   │
├─────────────────────────────────────────────────────────────┤
│  类型推断(Type Inference)    │  自引用约束(Self-Reference)    │
│  ├─ 函数参数推断            │  ├─ F-bounded 多态性          │
│  ├─ 约束类型推断            │  ├─ 递归类型定义              │
│  └─ 部分实例化              │  └─ Go 1.26 新特性            │
└─────────────────────────────────────────────────────────────┘
```

### 3. 关系依赖

```
                    泛型编程
                       │
        ┌──────────────┼──────────────┐
        │              │              │
    类型参数      类型约束       类型推断
        │              │              │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │         │    │         │    │         │
 单参数   多参数  内置约束  自定义  函数推断  约束推断
   │         │    │         │    │         │
   T       [K,V]  any    interface 泛化调用 类型推导
        comparable  ~int   简化语法
```

### 4. 详细示例代码

#### 4.1 基础泛型函数

```go
package main

import "fmt"

// Min 返回两个值中的较小值
// T 是类型参数，comparable 是约束
func Min[T comparable](a, b T) T {
    // 注意：comparable 只支持 == 和 !=
    // 不支持 <, > 比较
    if a == b {
        return a
    }
    // 无法直接比较大小，这里仅作示例
    return a
}

// MinOrdered 支持可排序类型
// 使用 constraints.Ordered 约束
import "golang.org/x/exp/constraints"

func MinOrdered[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

func main() {
    // 类型推断：编译器自动推断 T 为 int
    m1 := MinOrdered(10, 20)
    fmt.Println(m1) // 10

    // 显式指定类型参数
    m2 := MinOrdered[float64](3.14, 2.71)
    fmt.Println(m2) // 2.71
}
```

#### 4.2 泛型数据结构

```go
package main

import "fmt"

// Stack 泛型栈实现
type Stack[T any] struct {
    items []T
}

// Push 入栈
func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// Pop 出栈
func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// Peek 查看栈顶
func (s *Stack[T]) Peek() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

// Map 泛型 Map 函数
func Map[T, R any](input []T, fn func(T) R) []R {
    result := make([]R, len(input))
    for i, v := range input {
        result[i] = fn(v)
    }
    return result
}

// Filter 泛型 Filter 函数
func Filter[T any](input []T, predicate func(T) bool) []T {
    var result []T
    for _, v := range input {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce 泛型 Reduce 函数
func Reduce[T, R any](input []T, initial R, fn func(R, T) R) R {
    result := initial
    for _, v := range input {
        result = fn(result, v)
    }
    return result
}

func main() {
    // 整数栈
    intStack := Stack[int]{}
    intStack.Push(1)
    intStack.Push(2)
    intStack.Push(3)

    val, _ := intStack.Pop()
    fmt.Println(val) // 3

    // 字符串栈
    strStack := Stack[string]{}
    strStack.Push("hello")
    strStack.Push("world")

    s, _ := strStack.Pop()
    fmt.Println(s) // "world"

    // 使用 Map
    numbers := []int{1, 2, 3, 4, 5}
    doubled := Map(numbers, func(n int) int {
        return n * 2
    })
    fmt.Println(doubled) // [2 4 6 8 10]

    // 使用 Filter
    evens := Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // [2 4]

    // 使用 Reduce
    sum := Reduce(numbers, 0, func(acc, n int) int {
        return acc + n
    })
    fmt.Println(sum) // 15
}
```

#### 4.3 类型约束详解

```go
package main

import "fmt"

// 自定义类型约束
// Number 约束允许所有数值类型
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
        ~float32 | ~float64
}

// Signed 有符号整数约束
type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned 无符号整数约束
type Unsigned interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float 浮点数约束
type Float interface {
    ~float32 | ~float64
}

// Add 泛型加法函数
func Add[T Number](a, b T) T {
    return a + b
}

// Sum 计算切片总和
func Sum[T Number](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}

// Average 计算平均值
func Average[T Number](values []T) float64 {
    if len(values) == 0 {
        return 0
    }
    // 需要转换为 float64 进行除法
    var sum float64
    for _, v := range values {
        sum += float64(any(v).(float64))
    }
    return sum / float64(len(values))
}

// 更好的 Average 实现 - 使用类型参数约束
func AverageBetter[T ~int | ~int64 | ~float64](values []T) float64 {
    if len(values) == 0 {
        return 0
    }
    var sum float64
    for _, v := range values {
        sum += float64(v)
    }
    return sum / float64(len(values))
}

// 基于底层类型的约束
type MyInt int // MyInt 的底层类型是 int

func Process[T ~int](value T) {
    fmt.Printf("Processing: %v\n", value)
}

func main() {
    // 使用 Number 约束
    fmt.Println(Add(10, 20))       // 30
    fmt.Println(Add(3.14, 2.86))   // 6

    // 使用 Sum
    ints := []int{1, 2, 3, 4, 5}
    fmt.Println(Sum(ints)) // 15

    floats := []float64{1.5, 2.5, 3.5}
    fmt.Println(Sum(floats)) // 7.5

    // 底层类型匹配
    var mi MyInt = 100
    Process(mi) // 成功：MyInt 的底层类型是 int
}
```

#### 4.4 泛型接口

```go
package main

import "fmt"

// Comparable 可比较接口
type Comparable[T any] interface {
    Compare(other T) int
    Equal(other T) bool
}

// Container 泛型容器接口
type Container[T any] interface {
    Add(item T)
    Remove(item T) bool
    Contains(item T) bool
    Size() int
}

// Sorter 排序接口
type Sorter[T any] interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
    Get(i int) T
}

// 实现 Comparable 的类型
type Person struct {
    Name string
    Age  int
}

func (p Person) Compare(other Person) int {
    if p.Age < other.Age {
        return -1
    }
    if p.Age > other.Age {
        return 1
    }
    return 0
}

func (p Person) Equal(other Person) bool {
    return p.Name == other.Name && p.Age == other.Age
}

// Set 泛型集合实现
type Set[T comparable] struct {
    items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
    return &Set[T]{
        items: make(map[T]struct{}),
    }
}

func (s *Set[T]) Add(item T) {
    s.items[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) bool {
    if _, exists := s.items[item]; exists {
        delete(s.items, item)
        return true
    }
    return false
}

func (s *Set[T]) Contains(item T) bool {
    _, exists := s.items[item]
    return exists
}

func (s *Set[T]) Size() int {
    return len(s.items)
}

func (s *Set[T]) ToSlice() []T {
    result := make([]T, 0, len(s.items))
    for item := range s.items {
        result = append(result, item)
    }
    return result
}

func main() {
    // 整数集合
    intSet := NewSet[int]()
    intSet.Add(1)
    intSet.Add(2)
    intSet.Add(3)
    intSet.Add(2) // 重复，不会添加

    fmt.Println(intSet.Size())       // 3
    fmt.Println(intSet.Contains(2))  // true
    fmt.Println(intSet.Contains(10)) // false

    // 字符串集合
    strSet := NewSet[string]()
    strSet.Add("apple")
    strSet.Add("banana")
    strSet.Add("apple") // 重复

    fmt.Println(strSet.ToSlice()) // [apple banana] 或 [banana apple]

    // Person 比较
    p1 := Person{Name: "Alice", Age: 30}
    p2 := Person{Name: "Bob", Age: 25}

    fmt.Println(p1.Compare(p2)) // 1 (Alice 年龄更大)
    fmt.Println(p1.Equal(p1))   // true
}
```

### 5. 反例说明

```go
package main

// ❌ 错误示例 1: 类型参数不能用于常量声明
// const MaxSize[T any] = 100  // 编译错误

// ❌ 错误示例 2: 不能在方法上声明类型参数
// type Container struct{}
// func (c *Container) Push[T any](item T) {}  // 编译错误

// ✅ 正确做法：在类型上声明类型参数
type Container[T any] struct {
    items []T
}

func (c *Container[T]) Push(item T) { // 方法使用类型的类型参数
    c.items = append(c.items, item)
}

// ❌ 错误示例 3: 类型参数不能参与类型断言
// func process[T any](v T) {
//     if s, ok := v.(string); ok {  // 编译错误
//         fmt.Println(s)
//     }
// }

// ✅ 正确做法：使用 any 类型或类型约束
func processAny(v any) {
    if s, ok := v.(string); ok {
        println(s)
    }
}

// ❌ 错误示例 4: 不能对类型参数使用类型开关
// func switchType[T any](v T) {
//     switch v.(type) {  // 编译错误
//     case int:
//         println("int")
//     }
// }

// ❌ 错误示例 5: 泛型类型不能嵌入到接口中
// type MyInterface interface {
//     Container[int]  // 编译错误
// }

// ❌ 错误示例 6: 类型参数不能用于接收者类型
// func (T) Method() {}  // 编译错误

// ❌ 错误示例 7: 泛型函数不能作为值使用（除非实例化）
// var f = Min  // 编译错误：Min 是泛型函数

// ✅ 正确做法：实例化后使用
func Min[T comparable](a, b T) T {
    return a
}

var intMin = Min[int] // ✅ 正确：实例化后的函数

func main() {
    result := intMin(10, 20)
    println(result)
}
```

---

## 接口类型

### 1. 概念定义

Go 的接口定义了一组方法签名，任何实现了这些方法的类型都隐式实现了该接口。Go 1.18+ 引入了泛型接口，Go 1.20+ 引入了接口类型集和更灵活的约束。

### 2. 属性特征

```
┌─────────────────────────────────────────────────────────────┐
│                    接口类型分类                              │
├─────────────────────────────────────────────────────────────┤
│  基本接口 (Basic Interface)                                  │
│  ├─ 仅包含方法签名                                          │
│  ├─ 传统接口定义方式                                        │
│  └─ 例：type Reader interface { Read([]byte) (int, error) } │
├─────────────────────────────────────────────────────────────┤
│  泛型接口 (Generic Interface)                                │
│  ├─ 包含类型参数                                            │
│  ├─ 可用于泛型约束                                          │
│  └─ 例：type Container[T any] interface { Add(T) }          │
├─────────────────────────────────────────────────────────────┤
│  类型集接口 (Type Set Interface)                             │
│  ├─ 包含类型列表（Go 1.18+）                                │
│  ├─ 用于约束类型参数                                        │
│  └─ 例：type Number interface { ~int \| ~float64 }          │
├─────────────────────────────────────────────────────────────┤
│  接口组合 (Interface Composition)                            │
│  ├─ 嵌入其他接口                                            │
│  ├─ 组合多个接口能力                                        │
│  └─ 例：type ReadWriter interface { Reader; Writer }        │
└─────────────────────────────────────────────────────────────┘
```

### 3. 关系依赖

```
                        接口类型
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    基本接口          泛型接口          类型集接口
        │                  │                  │
   ┌────┴────┐        ┌────┴────┐        ┌────┴────┐
   │         │        │         │        │         │
 方法签名  接口组合  类型参数  约束使用  类型列表  ~底层类型
   │         │        │         │        │         │
 Reader   io.Read   [T any]  泛型函数  ~int    MyInt
 Writer   Closer    泛型类型  类型约束  ~string 底层int
```

### 4. 详细示例代码

#### 4.1 基本接口

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

// 定义基本接口
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Stringer interface {
    String() string
}

// Rectangle 实现 Shape 接口
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
    return fmt.Sprintf("Rectangle(%.2f x %.2f)", r.Width, r.Height)
}

// Circle 实现 Shape 接口
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

func (c Circle) String() string {
    return fmt.Sprintf("Circle(r=%.2f)", c.Radius)
}

// 多态函数
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// 空接口 - 可以存储任何类型
func PrintAnything(v interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}

// 类型断言
func ProcessValue(v interface{}) {
    // 类型断言
    if s, ok := v.(string); ok {
        fmt.Println("String:", s)
        return
    }

    // 类型开关
    switch val := v.(type) {
    case int:
        fmt.Println("Integer:", val)
    case float64:
        fmt.Println("Float:", val)
    case Shape:
        fmt.Println("Shape area:", val.Area())
    default:
        fmt.Println("Unknown type:", val)
    }
}

func main() {
    // 使用接口实现多态
    shapes := []Shape{
        Rectangle{Width: 5, Height: 3},
        Circle{Radius: 4},
    }

    for _, shape := range shapes {
        PrintShapeInfo(shape)
    }

    // 空接口
    PrintAnything(42)
    PrintAnything("hello")
    PrintAnything(Rectangle{5, 3})

    // 类型断言和类型开关
    ProcessValue(100)
    ProcessValue(3.14)
    ProcessValue("test")
    ProcessValue(Circle{Radius: 5})

    // 标准库接口使用
    var r io.Reader = strings.NewReader("Hello, World!")
    buf := make([]byte, 20)
    n, _ := r.Read(buf)
    fmt.Println(string(buf[:n]))
}
```

#### 4.2 接口组合

```go
package main

import (
    "fmt"
    "io"
)

// 基础接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// 接口组合
type ReadWriter interface {
    Reader
    Writer
}

type ReadCloser interface {
    Reader
    Closer
}

type WriteCloser interface {
    Writer
    Closer
}

// 组合三个接口
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// 自定义组合接口
type Named interface {
    Name() string
}

type Aged interface {
    Age() int
}

type Describable interface {
    Describe() string
}

// 组合多个接口
type PersonInterface interface {
    Named
    Aged
    Describable
}

// 实现类型
type Person struct {
    name string
    age  int
}

func (p Person) Name() string {
    return p.name
}

func (p Person) Age() int {
    return p.age
}

func (p Person) Describe() string {
    return fmt.Sprintf("%s is %d years old", p.name, p.age)
}

// 使用组合接口的函数
func PrintPersonInfo(p PersonInterface) {
    fmt.Printf("Name: %s\n", p.Name())
    fmt.Printf("Age: %d\n", p.Age())
    fmt.Printf("Description: %s\n", p.Describe())
}

// 带方法的组合接口
type AdvancedReader interface {
    Reader
    ReadString() (string, error)
    ReadLine() ([]byte, error)
}

// 与标准库接口兼容
func UseStandardLibrary(rw io.ReadWriter) {
    // 可以使用 io.ReadWriter 的所有方法
    buf := make([]byte, 100)
    rw.Read(buf)
    rw.Write([]byte("hello"))
}

func main() {
    person := Person{name: "Alice", age: 30}
    PrintPersonInfo(person)

    // 检查接口实现
    var _ PersonInterface = Person{}
    var _ Named = Person{}
    var _ Aged = Person{}
    var _ Describable = Person{}

    fmt.Println("All interfaces implemented!")
}
```

#### 4.3 泛型接口与约束

```go
package main

import (
    "fmt"
    "golang.org/x/exp/constraints"
)

// 泛型接口定义
type Container[T any] interface {
    Add(item T)
    Remove(item T) bool
    Get(index int) (T, bool)
    Size() int
}

// 可比较元素的容器
type ComparableContainer[T comparable] interface {
    Container[T]
    Contains(item T) bool
}

// 可排序容器
type OrderedContainer[T constraints.Ordered] interface {
    Container[T]
    Sort()
    Min() T
    Max() T
}

// 数值容器
type NumericContainer[T Number] interface {
    Container[T]
    Sum() T
    Average() float64
}

// 数值约束（类型集）
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
        ~float32 | ~float64
}

// 实现泛型接口的类型
type SliceContainer[T any] struct {
    items []T
}

func NewSliceContainer[T any]() *SliceContainer[T] {
    return &SliceContainer[T]{items: make([]T, 0)}
}

func (c *SliceContainer[T]) Add(item T) {
    c.items = append(c.items, item)
}

func (c *SliceContainer[T]) Remove(item T) bool {
    for i, v := range c.items {
        // 注意：这里需要 comparable 约束才能使用 ==
        // 暂时跳过实现
        _ = v
        _ = i
    }
    return false
}

func (c *SliceContainer[T]) Get(index int) (T, bool) {
    var zero T
    if index < 0 || index >= len(c.items) {
        return zero, false
    }
    return c.items[index], true
}

func (c *SliceContainer[T]) Size() int {
    return len(c.items)
}

// 使用泛型接口约束的函数
func ProcessContainer[T any](c Container[T], processor func(T) T) {
    for i := 0; i < c.Size(); i++ {
        if item, ok := c.Get(i); ok {
            processed := processor(item)
            c.Add(processed)
        }
    }
}

// 约束接口作为类型参数
type Adder[T any] interface {
    Add(a, b T) T
}

type IntAdder struct{}

func (IntAdder) Add(a, b int) int {
    return a + b
}

// 泛型函数使用接口约束
func SumWithAdder[T any](items []T, adder Adder[T]) T {
    var sum T
    for _, item := range items {
        sum = adder.Add(sum, item)
    }
    return sum
}

func main() {
    // 使用泛型容器
    container := NewSliceContainer[int]()
    container.Add(1)
    container.Add(2)
    container.Add(3)

    fmt.Println("Size:", container.Size())

    if val, ok := container.Get(1); ok {
        fmt.Println("Item at index 1:", val)
    }

    // 验证接口实现
    var _ Container[int] = container

    fmt.Println("Generic interface example completed!")
}
```

### 5. 反例说明

```go
package main

// ❌ 错误示例 1: 接口不能包含字段
// type BadInterface interface {
//     Name string  // 编译错误：接口只能包含方法
// }

// ✅ 正确做法：使用方法
 type GoodInterface interface {
     GetName() string
 }

// ❌ 错误示例 2: 接口方法不能有实现
// type BadInterface interface {
//     Method() { println("implementation") }  // 编译错误
// }

// ✅ 正确做法：接口只声明方法签名
 type GoodInterface2 interface {
     Method()
 }

// ❌ 错误示例 3: 不能循环嵌入接口
// type A interface {
//     B
// }
// type B interface {
//     A  // 编译错误：循环嵌入
// }

// ❌ 错误示例 4: 接口不能嵌入具体类型
// type MyStruct struct{}
// type BadInterface interface {
//     MyStruct  // 编译错误
// }

// ❌ 错误示例 5: 类型集接口不能用于变量声明
// type Number interface {
//     ~int | ~float64
// }
// var n Number  // 编译错误：Number 是约束接口

// ✅ 正确做法：类型集接口只用于泛型约束
 func UseNumber[T Number](v T) T {
     return v
 }

// ❌ 错误示例 6: 接口方法名冲突
// type Reader interface {
//     Read() byte
// }
// type Writer interface {
//     Read() string  // 方法名相同但返回类型不同
// }
// type BadCombo interface {
//     Reader
//     Writer  // 编译错误：方法 Read 冲突
// }

// ❌ 错误示例 7: 非接口类型不能嵌入接口
// type BadStruct struct {
//     io.Reader  // 这是合法的（嵌入接口类型）
//     string     // 这是合法的（嵌入具体类型）
// }
// 但非接口类型不能"实现"接口的嵌入

func main() {
    // 接口使用示例
    println("Interface examples")
}
```

---

## 反射机制

### 1. 概念定义

反射是程序在运行时检查和操作自身结构的能力。Go 通过 `reflect` 包提供反射功能，允许动态检查类型信息、访问和修改值。

### 2. 属性特征

```
┌─────────────────────────────────────────────────────────────┐
│                    反射核心组件                              │
├─────────────────────────────────────────────────────────────┤
│  reflect.Type  - 类型信息                                    │
│  ├─ 获取类型：reflect.TypeOf(v)                             │
│  ├─ 类型种类：Kind() - struct, slice, map, ptr 等           │
│  ├─ 结构体字段：NumField(), Field(i)                        │
│  └─ 方法信息：NumMethod(), Method(i)                        │
├─────────────────────────────────────────────────────────────┤
│  reflect.Value - 值操作                                      │
│  ├─ 获取值：reflect.ValueOf(v)                              │
│  ├─ 类型获取：Type()                                        │
│  ├─ 种类判断：Kind()                                        │
│  ├─ 值获取：Int(), String(), Bool() 等                      │
│  └─ 值设置：SetInt(), SetString() 等（需可设置）            │
├─────────────────────────────────────────────────────────────┤
│  可设置性 (Settability)                                      │
│  ├─ 值必须可寻址才能修改                                    │
│  ├─ 使用指针或 Elem() 获取可设置值                          │
│  └─ CanSet() 检查是否可设置                                 │
├─────────────────────────────────────────────────────────────┤
│  动态创建                                                    │
│  ├─ 创建值：reflect.New(), reflect.MakeSlice()              │
│  ├─ 创建映射：reflect.MakeMap(), reflect.MakeChan()         │
│  └─ 调用函数：reflect.Value.Call()                          │
└─────────────────────────────────────────────────────────────┘
```

### 3. 关系依赖

```
                    反射机制
                       │
        ┌──────────────┼──────────────┐
        │              │              │
    reflect.Type  reflect.Value   动态创建
        │              │              │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │         │    │         │    │         │
 TypeOf()  Kind() ValueOf()  Get() New()   MakeSlice()
 NumField() String() Set()   Set() MakeMap() MakeChan()
 Field()   Int()   CanSet()  Call() Zero()  Indirect()
```

### 4. 详细示例代码

#### 4.1 基础反射操作

```go
package main

import (
    "fmt"
    "reflect"
)

func basicReflection() {
    // 基本类型的反射
    x := 42
    t := reflect.TypeOf(x)
    v := reflect.ValueOf(x)

    fmt.Println("Type:", t)           // int
    fmt.Println("Value:", v)          // 42
    fmt.Println("Kind:", v.Kind())    // int

    // 字符串反射
    s := "hello"
    fmt.Println("String Type:", reflect.TypeOf(s))   // string
    fmt.Println("String Value:", reflect.ValueOf(s)) // hello

    // 指针反射
    p := &x
    fmt.Println("Pointer Type:", reflect.TypeOf(p))  // *int
    fmt.Println("Pointer Kind:", reflect.ValueOf(p).Kind()) // ptr

    // 获取指针指向的值
    elem := reflect.ValueOf(p).Elem()
    fmt.Println("Elem Value:", elem)      // 42
    fmt.Println("Elem Type:", elem.Type()) // int
}

func structReflection() {
    type Person struct {
        Name    string `json:"name" db:"user_name"`
        Age     int    `json:"age"`
        Email   string `json:"email,omitempty"`
        private string // 未导出字段
    }

    p := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

    t := reflect.TypeOf(p)
    v := reflect.ValueOf(p)

    fmt.Println("Struct Type:", t)
    fmt.Println("NumField:", t.NumField())

    // 遍历结构体字段
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)

        fmt.Printf("Field %d: Name=%s, Type=%s, Value=%v, Tag=%s\n",
            i, field.Name, field.Type, value, field.Tag)

        // 获取特定 tag
        jsonTag := field.Tag.Get("json")
        dbTag := field.Tag.Get("db")
        fmt.Printf("  json tag: %s, db tag: %s\n", jsonTag, dbTag)
    }

    // 通过名称获取字段
    if nameField, ok := t.FieldByName("Name"); ok {
        fmt.Println("Name field found:", nameField.Type)
    }
}

func modifyValue() {
    x := 42

    // ❌ 不能直接修改值
    // v := reflect.ValueOf(x)
    // v.SetInt(100) // panic: 值不可设置

    // ✅ 通过指针修改
    v := reflect.ValueOf(&x).Elem()
    fmt.Println("CanSet:", v.CanSet()) // true

    if v.CanSet() {
        v.SetInt(100)
    }
    fmt.Println("Modified x:", x) // 100

    // 修改结构体字段
    type Person struct {
        Name string
        Age  int
    }

    p := &Person{Name: "Bob", Age: 25}
    pv := reflect.ValueOf(p).Elem()

    nameField := pv.FieldByName("Name")
    if nameField.CanSet() {
        nameField.SetString("Alice")
    }

    ageField := pv.FieldByName("Age")
    if ageField.CanSet() {
        ageField.SetInt(30)
    }

    fmt.Println("Modified person:", *p) // {Alice 30}
}

func sliceAndMapReflection() {
    // 切片反射
    nums := []int{1, 2, 3, 4, 5}
    v := reflect.ValueOf(nums)

    fmt.Println("Slice Kind:", v.Kind())      // slice
    fmt.Println("Slice Len:", v.Len())        // 5
    fmt.Println("Slice Cap:", v.Cap())        // 5

    // 获取元素
    for i := 0; i < v.Len(); i++ {
        elem := v.Index(i)
        fmt.Printf("Element %d: %v\n", i, elem.Int())
    }

    // 修改切片元素
    if v.Index(0).CanSet() {
        v.Index(0).SetInt(100)
    }
    fmt.Println("Modified slice:", nums) // [100 2 3 4 5]

    // 映射反射
    m := map[string]int{"a": 1, "b": 2, "c": 3}
    mv := reflect.ValueOf(m)

    fmt.Println("Map Kind:", mv.Kind()) // map

    // 获取映射值
    key := reflect.ValueOf("b")
    val := mv.MapIndex(key)
    fmt.Println("Value for 'b':", val) // 2

    // 遍历映射
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        fmt.Printf("Key: %v, Value: %v\n", k, v)
    }

    // 修改映射
    mv.SetMapIndex(reflect.ValueOf("d"), reflect.ValueOf(4))
    fmt.Println("Modified map:", m) // map[a:1 b:2 c:3 d:4]

    // 删除映射元素
    mv.SetMapIndex(reflect.ValueOf("a"), reflect.Value{})
    fmt.Println("After delete:", m) // map[b:2 c:3 d:4]
}

func dynamicCreation() {
    // 动态创建值
    intType := reflect.TypeOf(0)
    newInt := reflect.New(intType) // 创建 *int
    newInt.Elem().SetInt(42)
    fmt.Println("New int:", newInt.Elem().Int()) // 42

    // 动态创建切片
    sliceType := reflect.TypeOf([]int{})
    newSlice := reflect.MakeSlice(sliceType, 3, 5)
    for i := 0; i < newSlice.Len(); i++ {
        newSlice.Index(i).SetInt(int64(i * 10))
    }
    fmt.Println("New slice:", newSlice.Interface()) // [0 10 20]

    // 动态创建映射
    mapType := reflect.TypeOf(map[string]int{})
    newMap := reflect.MakeMap(mapType)
    newMap.SetMapIndex(reflect.ValueOf("key1"), reflect.ValueOf(100))
    newMap.SetMapIndex(reflect.ValueOf("key2"), reflect.ValueOf(200))
    fmt.Println("New map:", newMap.Interface()) // map[key1:100 key2:200]

    // 动态创建结构体
    type Person struct {
        Name string
        Age  int
    }

    personType := reflect.TypeOf(Person{})
    newPerson := reflect.New(personType).Elem()
    newPerson.FieldByName("Name").SetString("Charlie")
    newPerson.FieldByName("Age").SetInt(35)
    fmt.Println("New person:", newPerson.Interface()) // {Charlie 35}
}

func functionReflection() {
    // 函数反射
    add := func(a, b int) int {
        return a + b
    }

    fn := reflect.ValueOf(add)
    fmt.Println("Function Kind:", fn.Kind()) // func
    fmt.Println("Function Type:", fn.Type()) // func(int, int) int

    // 调用函数
    args := []reflect.Value{
        reflect.ValueOf(10),
        reflect.ValueOf(20),
    }
    results := fn.Call(args)
    fmt.Println("Result:", results[0]) // 30

    // 变参函数
    sum := func(nums ...int) int {
        total := 0
        for _, n := range nums {
            total += n
        }
        return total
    }

    sumFn := reflect.ValueOf(sum)
    sumArgs := []reflect.Value{
        reflect.ValueOf([]int{1, 2, 3, 4, 5}),
    }
    sumResults := sumFn.Call(sumArgs)
    fmt.Println("Sum result:", sumResults[0]) // 15
}

func main() {
    fmt.Println("=== Basic Reflection ===")
    basicReflection()

    fmt.Println("\n=== Struct Reflection ===")
    structReflection()

    fmt.Println("\n=== Modify Value ===")
    modifyValue()

    fmt.Println("\n=== Slice and Map Reflection ===")
    sliceAndMapReflection()

    fmt.Println("\n=== Dynamic Creation ===")
    dynamicCreation()

    fmt.Println("\n=== Function Reflection ===")
    functionReflection()
}
```

#### 4.2 高级反射应用

```go
package main

import (
    "fmt"
    "reflect"
)

// 结构体验证器
type Validator struct {
    Errors map[string]string
}

func NewValidator() *Validator {
    return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Validate(s interface{}) bool {
    v.Errors = make(map[string]string)

    val := reflect.ValueOf(s)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        // 检查 validate tag
        tag := fieldType.Tag.Get("validate")
        if tag == "" {
            continue
        }

        fieldName := fieldType.Name

        switch tag {
        case "required":
            if isZero(field) {
                v.Errors[fieldName] = "field is required"
            }
        case "min=3":
            if field.Kind() == reflect.String && len(field.String()) < 3 {
                v.Errors[fieldName] = "must be at least 3 characters"
            }
        case "positive":
            if field.Kind() >= reflect.Int && field.Kind() <= reflect.Float64 {
                if field.Convert(reflect.TypeOf(float64(0))).Float() <= 0 {
                    v.Errors[fieldName] = "must be positive"
                }
            }
        }
    }

    return len(v.Errors) == 0
}

func isZero(v reflect.Value) bool {
    switch v.Kind() {
    case reflect.String:
        return v.String() == ""
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return v.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return v.Uint() == 0
    case reflect.Float32, reflect.Float64:
        return v.Float() == 0
    case reflect.Bool:
        return !v.Bool()
    case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan:
        return v.IsNil()
    default:
        return false
    }
}

// 深拷贝
defunc DeepCopy(dst, src interface{}) {
    dstVal := reflect.ValueOf(dst)
    srcVal := reflect.ValueOf(src)

    if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
        panic("dst must be a non-nil pointer")
    }

    dstVal = dstVal.Elem()
    if srcVal.Kind() == reflect.Ptr {
        srcVal = srcVal.Elem()
    }

    deepCopyValue(dstVal, srcVal)
}

func deepCopyValue(dst, src reflect.Value) {
    if !dst.CanSet() {
        return
    }

    switch src.Kind() {
    case reflect.Ptr:
        if src.IsNil() {
            return
        }
        dst.Set(reflect.New(src.Elem().Type()))
        deepCopyValue(dst.Elem(), src.Elem())

    case reflect.Interface:
        if src.IsNil() {
            return
        }
        dst.Set(reflect.New(src.Elem().Type()).Elem())
        deepCopyValue(dst.Elem(), src.Elem())

    case reflect.Struct:
        for i := 0; i < src.NumField(); i++ {
            if src.Type().Field(i).PkgPath != "" {
                continue // 跳过未导出字段
            }
            deepCopyValue(dst.Field(i), src.Field(i))
        }

    case reflect.Slice:
        if src.IsNil() {
            return
        }
        dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
        for i := 0; i < src.Len(); i++ {
            deepCopyValue(dst.Index(i), src.Index(i))
        }

    case reflect.Map:
        if src.IsNil() {
            return
        }
        dst.Set(reflect.MakeMap(src.Type()))
        for _, key := range src.MapKeys() {
            newVal := reflect.New(src.MapIndex(key).Type()).Elem()
            deepCopyValue(newVal, src.MapIndex(key))
            dst.SetMapIndex(key, newVal)
        }

    default:
        dst.Set(src)
    }
}

// 结构体转映射
func StructToMap(s interface{}) map[string]interface{} {
    result := make(map[string]interface{})

    val := reflect.ValueOf(s)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        // 跳过未导出字段
        if fieldType.PkgPath != "" {
            continue
        }

        // 使用 json tag 作为键，如果没有则使用字段名
        key := fieldType.Tag.Get("json")
        if key == "" || key == "-" {
            key = fieldType.Name
        }

        // 处理嵌套结构体
        if field.Kind() == reflect.Struct && fieldType.Anonymous {
            nested := StructToMap(field.Interface())
            for k, v := range nested {
                result[k] = v
            }
        } else {
            result[key] = field.Interface()
        }
    }

    return result
}

// 映射转结构体
func MapToStruct(m map[string]interface{}, s interface{}) error {
    val := reflect.ValueOf(s)
    if val.Kind() != reflect.Ptr || val.IsNil() {
        return fmt.Errorf("s must be a non-nil pointer")
    }

    val = val.Elem()
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)

        // 跳过未导出字段
        if !field.CanSet() {
            continue
        }

        // 获取字段对应的 map 键
        key := fieldType.Tag.Get("json")
        if key == "" || key == "-" {
            key = fieldType.Name
        }

        // 查找对应的值
        mapVal, exists := m[key]
        if !exists {
            continue
        }

        // 设置值
        if err := setField(field, reflect.ValueOf(mapVal)); err != nil {
            return fmt.Errorf("field %s: %v", fieldType.Name, err)
        }
    }

    return nil
}

func setField(field, value reflect.Value) error {
    if !field.CanSet() {
        return fmt.Errorf("cannot set field")
    }

    // 类型转换
    if field.Type() == value.Type() {
        field.Set(value)
        return nil
    }

    // 尝试类型转换
    if value.Type().ConvertibleTo(field.Type()) {
        field.Set(value.Convert(field.Type()))
        return nil
    }

    return fmt.Errorf("cannot convert %v to %v", value.Type(), field.Type())
}

// 使用示例
type Address struct {
    City    string `json:"city"`
    Country string `json:"country"`
}

type User struct {
    ID      int     `json:"id" validate:"positive"`
    Name    string  `json:"name" validate:"required,min=3"`
    Email   string  `json:"email" validate:"required"`
    Age     int     `json:"age"`
    Address Address `json:"address"`
}

func main() {
    // 验证器示例
    validator := NewValidator()

    user1 := User{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30}
    if validator.Validate(user1) {
        fmt.Println("User1 is valid")
    } else {
        fmt.Println("User1 errors:", validator.Errors)
    }

    user2 := User{ID: 0, Name: "Al", Email: ""}
    if validator.Validate(user2) {
        fmt.Println("User2 is valid")
    } else {
        fmt.Println("User2 errors:", validator.Errors)
    }

    // 深拷贝示例
    original := &User{
        ID:    1,
        Name:  "Bob",
        Email: "bob@example.com",
        Age:   25,
        Address: Address{
            City:    "New York",
            Country: "USA",
        },
    }

    var copied User
    DeepCopy(&copied, original)
    copied.Name = "Copied Bob"
    copied.Address.City = "Los Angeles"

    fmt.Println("Original:", original)
    fmt.Println("Copied:", copied)

    // 结构体转映射
    userMap := StructToMap(user1)
    fmt.Println("User map:", userMap)

    // 映射转结构体
    newUser := &User{}
    err := MapToStruct(userMap, newUser)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("New user:", newUser)
    }
}
```

### 5. 反例说明

```go
package main

import "reflect"

// ❌ 错误示例 1: 修改不可设置的值
// func badModify() {
//     x := 42
//     v := reflect.ValueOf(x)
//     v.SetInt(100) // panic: reflect.Value.SetInt using unaddressable value
// }

// ✅ 正确做法：使用指针
 func goodModify() {
     x := 42
     v := reflect.ValueOf(&x).Elem()
     v.SetInt(100) // ✅ 成功
 }

// ❌ 错误示例 2: 访问未导出字段
// type MyStruct struct {
//     privateField int
// }
// func badAccess() {
//     s := MyStruct{privateField: 42}
//     v := reflect.ValueOf(s)
//     field := v.FieldByName("privateField")
//     // 可以访问，但不能设置
//     // field.SetInt(100) // panic: reflect.Value.SetInt using value obtained using unexported field
// }

// ❌ 错误示例 3: 类型不匹配的操作
// func badType() {
//     x := "hello"
//     v := reflect.ValueOf(x)
//     v.SetInt(100) // panic: reflect: call of reflect.Value.SetInt on string Value
// }

// ❌ 错误示例 4: 对 nil 指针调用 Elem()
// func badNil() {
//     var p *int
//     v := reflect.ValueOf(p)
//     elem := v.Elem() // panic: reflect: call of reflect.Value.Elem on zero Value
// }

// ✅ 正确做法：检查 IsNil()
 func goodNilCheck() {
     var p *int
     v := reflect.ValueOf(p)
     if !v.IsNil() {
         elem := v.Elem()
         _ = elem
     }
 }

// ❌ 错误示例 5: 使用错误的 Kind 进行操作
// func badKind() {
//     x := []int{1, 2, 3}
//     v := reflect.ValueOf(x)
//     v.MapKeys() // panic: reflect: call of reflect.Value.MapKeys on slice Value
// }

// ❌ 错误示例 6: 索引越界
// func badIndex() {
//     x := []int{1, 2, 3}
//     v := reflect.ValueOf(x)
//     v.Index(10) // panic: reflect: slice index out of range
// }

// ❌ 错误示例 7: 调用非函数值
// func badCall() {
//     x := 42
//     v := reflect.ValueOf(x)
//     v.Call(nil) // panic: reflect: call of reflect.Value.Call on int Value
// }

// ❌ 错误示例 8: 参数数量不匹配
// func badArgCount() {
//     fn := func(a, b int) int { return a + b }
//     v := reflect.ValueOf(fn)
//     v.Call([]reflect.Value{reflect.ValueOf(1)}) // panic: reflect.Call with wrong number of arguments
// }

// ❌ 错误示例 9: 参数类型不匹配
// func badArgType() {
//     fn := func(a int) {}
//     v := reflect.ValueOf(fn)
//     v.Call([]reflect.Value{reflect.ValueOf("hello")}) // panic: reflect.Value.Call using string as type int
// }

func main() {
    // 反射使用示例
    println("Reflection examples")
}
```

---

## 元编程

### 1. 概念定义

元编程是指编写能够生成、操作或转换代码的代码。Go 通过 `go generate` 命令和代码生成工具提供元编程能力，允许在编译前自动生成代码。

### 2. 属性特征

```
┌─────────────────────────────────────────────────────────────┐
│                    元编程核心机制                            │
├─────────────────────────────────────────────────────────────┤
│  go generate 命令                                            │
│  ├─ 扫描 //go:generate 指令                                  │
│  ├─ 执行指定的命令或程序                                     │
│  ├─ 按文件名顺序执行                                         │
│  └─ 独立于 go build 运行                                     │
├─────────────────────────────────────────────────────────────┤
│  代码生成工具                                                │
│  ├─ stringer：为常量生成 String() 方法                      │
│  ├─ jsonenums：生成 JSON 序列化代码                         │
│  ├─ protoc-gen-go：Protocol Buffers 代码生成               │
│  └─ 自定义生成工具                                           │
├─────────────────────────────────────────────────────────────┤
│  文本模板 (text/template)                                    │
│  ├─ 定义代码模板                                             │
│  ├─ 数据驱动生成                                             │
│  └─ 支持控制结构                                             │
├─────────────────────────────────────────────────────────────┤
│  AST 操作 (go/ast, go/parser)                                │
│  ├─ 解析 Go 源码为 AST                                       │
│  ├─ 遍历和修改 AST                                           │
│  └─ 生成修改后的代码                                         │
└─────────────────────────────────────────────────────────────┘
```

### 3. 关系依赖

```
                    元编程
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   go generate    文本模板        AST操作
        │              │              │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │         │    │         │    │         │
 指令扫描  工具执行 模板定义  数据驱动  源码解析  AST遍历
 顺序执行  代码生成 控制结构  文件输出  语法分析  代码生成
```

### 4. 详细示例代码

#### 4.1 go generate 基础使用

```go
package main

//go:generate stringer -type=Status
//go:generate stringer -type=Priority -linecomment

// Status 状态枚举
type Status int

const (
    StatusPending Status = iota
    StatusProcessing
    StatusCompleted
    StatusFailed
)

// Priority 优先级枚举
type Priority int

const (
    PriorityLow Priority = iota
    PriorityMedium
    PriorityHigh
)

func (p Priority) String() string {
    switch p {
    case PriorityLow:
        return "低优先级"
    case PriorityMedium:
        return "中优先级"
    case PriorityHigh:
        return "高优先级"
    default:
        return "未知优先级"
    }
}

// 使用 go generate 生成的代码
// 运行: go generate
// 会生成 status_string.go 和 priority_string.go

func main() {
    // 使用生成的 String() 方法
    fmt.Println(StatusPending)    // StatusPending
    fmt.Println(StatusProcessing) // StatusProcessing
    fmt.Println(PriorityHigh)     // 高优先级
}
```

#### 4.2 自定义代码生成工具

```go
// 文件: cmd/gencode/main.go
// 自定义代码生成器

package main

import (
    "flag"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "strings"
    "text/template"
)

const structTemplate = `// Code generated by gencode; DO NOT EDIT.

package {{.Package}}

{{range .Structs}}
// New{{.Name}} creates a new {{.Name}} with default values
func New{{.Name}}() *{{.Name}} {
    return &{{.Name}}{}
}

// Validate validates the {{.Name}} fields
func (s *{{.Name}}) Validate() error {
{{range .RequiredFields}}
    if s.{{.}} == nil || *s.{{.}} == "" {
        return fmt.Errorf("{{.}} is required")
    }
{{end}}
    return nil
}

// SetDefaults sets default values for optional fields
func (s *{{.Name}}) SetDefaults() {
{{range .OptionalFields}}
    if s.{{.}} == nil {
        default{{.}} := {{DefaultValue .}}
        s.{{.}} = &default{{.}}
    }
{{end}}
}
{{end}}
`

type StructInfo struct {
    Name           string
    RequiredFields []string
    OptionalFields []string
}

type TemplateData struct {
    Package string
    Structs []StructInfo
}

func main() {
    var (
        input  = flag.String("input", "", "Input Go file")
        output = flag.String("output", "", "Output file")
    )
    flag.Parse()

    if *input == "" || *output == "" {
        fmt.Fprintf(os.Stderr, "Usage: gencode -input=<file> -output=<file>\n")
        os.Exit(1)
    }

    // 解析输入文件
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, *input, nil, parser.ParseComments)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
        os.Exit(1)
    }

    data := TemplateData{
        Package: node.Name.Name,
        Structs: extractStructs(node),
    }

    // 执行模板
    tmpl, err := template.New("code").Parse(structTemplate)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
        os.Exit(1)
    }

    file, err := os.Create(*output)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    if err := tmpl.Execute(file, data); err != nil {
        fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Generated: %s\n", *output)
}

func extractStructs(node *ast.File) []StructInfo {
    var structs []StructInfo

    ast.Inspect(node, func(n ast.Node) bool {
        genDecl, ok := n.(*ast.GenDecl)
        if !ok || genDecl.Tok != token.TYPE {
            return true
        }

        for _, spec := range genDecl.Specs {
            typeSpec, ok := spec.(*ast.TypeSpec)
            if !ok {
                continue
            }

            structType, ok := typeSpec.Type.(*ast.StructType)
            if !ok {
                continue
            }

            info := StructInfo{Name: typeSpec.Name.Name}

            for _, field := range structType.Fields.List {
                fieldName := field.Names[0].Name

                // 检查 tag
                if field.Tag != nil {
                    tag := strings.Trim(field.Tag.Value, "`")
                    if strings.Contains(tag, "required") {
                        info.RequiredFields = append(info.RequiredFields, fieldName)
                    } else if strings.Contains(tag, "optional") {
                        info.OptionalFields = append(info.OptionalFields, fieldName)
                    }
                }
            }

            structs = append(structs, info)
        }

        return true
    })

    return structs
}
```

#### 4.3 使用 text/template 生成代码

```go
package main

import (
    "os"
    "text/template"
)

// 定义代码生成模板
const crudTemplate = `// Code generated by go generate; DO NOT EDIT.

package {{.Package}}

import (
    "context"
    "database/sql"
    "fmt"
    "time"
)

// {{.ModelName}}Store 提供 {{.ModelName}} 的 CRUD 操作
type {{.ModelName}}Store struct {
    db *sql.DB
}

// New{{.ModelName}}Store 创建存储实例
func New{{.ModelName}}Store(db *sql.DB) *{{.ModelName}}Store {
    return &{{.ModelName}}Store{db: db}
}

// Create 创建新记录
func (s *{{.ModelName}}Store) Create(ctx context.Context, m *{{.ModelName}}) error {
    query := ` + "`" + `INSERT INTO {{.TableName}} (
        {{range .Fields}}{{if .IsDB}},{{.DBName}}{{end}}{{end}}
    ) VALUES (
        {{range .Fields}}{{if .IsDB}}?{{if not .IsLast}},{{end}}{{end}}{{end}}
    )` + "`" + `

    result, err := s.db.ExecContext(ctx, query,
        {{range .Fields}}{{if .IsDB}}m.{{.Name}},{{end}}{{end}}
    )
    if err != nil {
        return fmt.Errorf("failed to create {{.ModelName}}: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert id: %w", err)
    }

    m.ID = id
    m.CreatedAt = time.Now()
    m.UpdatedAt = time.Now()

    return nil
}

// GetByID 根据 ID 获取记录
func (s *{{.ModelName}}Store) GetByID(ctx context.Context, id int64) (*{{.ModelName}}, error) {
    query := ` + "`" + `SELECT
        {{range .Fields}}{{if .IsDB}}{{if not .IsFirst}},{{end}}{{.DBName}}{{end}}{{end}}
    FROM {{.TableName}} WHERE id = ?` + "`" + `

    m := &{{.ModelName}}{}
    err := s.db.QueryRowContext(ctx, query, id).Scan(
        {{range .Fields}}{{if .IsDB}}&m.{{.Name}},{{end}}{{end}}
    )

    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("{{.ModelName}} not found: %w", err)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get {{.ModelName}}: %w", err)
    }

    return m, nil
}

// Update 更新记录
func (s *{{.ModelName}}Store) Update(ctx context.Context, m *{{.ModelName}}) error {
    query := ` + "`" + `UPDATE {{.TableName}} SET
        {{range .Fields}}{{if and .IsDB (ne .DBName "id") (ne .DBName "created_at")}}{{if not .IsFirst}},{{end}}{{.DBName}} = ?{{end}}{{end}},
        updated_at = ?
    WHERE id = ?` + "`" + `

    m.UpdatedAt = time.Now()

    _, err := s.db.ExecContext(ctx, query,
        {{range .Fields}}{{if and .IsDB (ne .DBName "id") (ne .DBName "created_at")}}m.{{.Name}},{{end}}{{end}}
        m.UpdatedAt,
        m.ID,
    )

    if err != nil {
        return fmt.Errorf("failed to update {{.ModelName}}: %w", err)
    }

    return nil
}

// Delete 删除记录
func (s *{{.ModelName}}Store) Delete(ctx context.Context, id int64) error {
    query := ` + "`" + `DELETE FROM {{.TableName}} WHERE id = ?` + "`" + `

    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete {{.ModelName}}: %w", err)
    }

    return nil
}

// List 列出所有记录
func (s *{{.ModelName}}Store) List(ctx context.Context, limit, offset int) ([]*{{.ModelName}}, error) {
    query := ` + "`" + `SELECT
        {{range .Fields}}{{if .IsDB}}{{if not .IsFirst}},{{end}}{{.DBName}}{{end}}{{end}}
    FROM {{.TableName}} ORDER BY created_at DESC LIMIT ? OFFSET ?` + "`" + `

    rows, err := s.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list {{.ModelName}}: %w", err)
    }
    defer rows.Close()

    var results []*{{.ModelName}}
    for rows.Next() {
        m := &{{.ModelName}}{}
        err := rows.Scan(
            {{range .Fields}}{{if .IsDB}}&m.{{.Name}},{{end}}{{end}}
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan {{.ModelName}}: %w", err)
        }
        results = append(results, m)
    }

    return results, rows.Err()
}
`

// 模板数据结构
type Field struct {
    Name     string
    Type     string
    DBName   string
    IsDB     bool
    IsFirst  bool
    IsLast   bool
}

type CRUDData struct {
    Package   string
    ModelName string
    TableName string
    Fields    []Field
}

func generateCRUD() {
    data := CRUDData{
        Package:   "models",
        ModelName: "User",
        TableName: "users",
        Fields: []Field{
            {Name: "ID", Type: "int64", DBName: "id", IsDB: true, IsFirst: true},
            {Name: "Name", Type: "string", DBName: "name", IsDB: true},
            {Name: "Email", Type: "string", DBName: "email", IsDB: true},
            {Name: "Age", Type: "int", DBName: "age", IsDB: true},
            {Name: "CreatedAt", Type: "time.Time", DBName: "created_at", IsDB: true},
            {Name: "UpdatedAt", Type: "time.Time", DBName: "updated_at", IsDB: true, IsLast: true},
        },
    }

    tmpl, err := template.New("crud").Parse(crudTemplate)
    if err != nil {
        panic(err)
    }

    file, err := os.Create("user_store_gen.go")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    if err := tmpl.Execute(file, data); err != nil {
        panic(err)
    }
}

func main() {
    generateCRUD()
}
```

#### 4.4 AST 操作示例

```go
package main

import (
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
)

// 分析 Go 源码的 AST
func analyzeAST(filename string) {
    fset := token.NewFileSet()

    // 解析文件
    node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        return
    }

    fmt.Printf("Package: %s\n", node.Name.Name)
    fmt.Printf("Number of declarations: %d\n", len(node.Decls))

    // 遍历所有声明
    for _, decl := range node.Decls {
        switch d := decl.(type) {
        case *ast.FuncDecl:
            analyzeFunction(d)
        case *ast.GenDecl:
            analyzeGeneralDecl(d)
        }
    }

    // 使用 Inspect 遍历所有节点
    fmt.Println("\n=== All Types ===")
    ast.Inspect(node, func(n ast.Node) bool {
        if n == nil {
            return true
        }

        switch x := n.(type) {
        case *ast.TypeSpec:
            fmt.Printf("Type: %s\n", x.Name.Name)
        case *ast.StructType:
            fmt.Printf("  Struct with %d fields\n", len(x.Fields.List))
        case *ast.InterfaceType:
            fmt.Printf("  Interface with %d methods\n", len(x.Methods.List))
        }

        return true
    })
}

func analyzeFunction(fn *ast.FuncDecl) {
    fmt.Printf("\nFunction: %s\n", fn.Name.Name)

    if fn.Recv != nil {
        fmt.Println("  Method (has receiver)")
        for _, recv := range fn.Recv.List {
            fmt.Printf("  Receiver: %s\n", recv.Names[0].Name)
        }
    }

    if fn.Type.Params != nil {
        fmt.Printf("  Parameters: %d\n", len(fn.Type.Params.List))
    }

    if fn.Type.Results != nil {
        fmt.Printf("  Results: %d\n", len(fn.Type.Results.List))
    }
}

func analyzeGeneralDecl(decl *ast.GenDecl) {
    switch decl.Tok {
    case token.IMPORT:
        fmt.Printf("\nImports: %d\n", len(decl.Specs))
        for _, spec := range decl.Specs {
            if imp, ok := spec.(*ast.ImportSpec); ok {
                fmt.Printf("  %s\n", imp.Path.Value)
            }
        }

    case token.TYPE:
        fmt.Printf("\nType declarations: %d\n", len(decl.Specs))
        for _, spec := range decl.Specs {
            if ts, ok := spec.(*ast.TypeSpec); ok {
                fmt.Printf("  Type: %s\n", ts.Name.Name)
            }
        }

    case token.CONST:
        fmt.Printf("\nConst declarations: %d\n", len(decl.Specs))

    case token.VAR:
        fmt.Printf("\nVar declarations: %d\n", len(decl.Specs))
    }
}

// 修改 AST 并生成代码
func modifyAST() {
    src := `package example

type User struct {
    Name string
    Age  int
}

func (u *User) GetName() string {
    return u.Name
}
`

    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, "example.go", src, parser.ParseComments)
    if err != nil {
        panic(err)
    }

    // 遍历并修改 AST
    ast.Inspect(node, func(n ast.Node) bool {
        if fn, ok := n.(*ast.FuncDecl); ok {
            if fn.Name.Name == "GetName" {
                // 修改函数名
                fn.Name.Name = "GetUserName"
            }
        }
        return true
    })

    // 打印修改后的代码
    fmt.Println("\n=== Modified AST ===")
    ast.Print(fset, node)
}

func main() {
    // 创建测试文件
    testFile := "/tmp/test_ast.go"
    content := `package test

import "fmt"

type User struct {
    Name string
    Age  int
}

func NewUser(name string, age int) *User {
    return &User{Name: name, Age: age}
}

func (u *User) String() string {
    return fmt.Sprintf("User{Name: %s, Age: %d}", u.Name, u.Age)
}

const MaxAge = 150

var DefaultUser = NewUser("default", 0)
`

    os.WriteFile(testFile, []byte(content), 0644)

    fmt.Println("=== AST Analysis ===")
    analyzeAST(testFile)

    fmt.Println("\n=== AST Modification ===")
    modifyAST()

    os.Remove(testFile)
}
```

### 5. 反例说明

```go
package main

// ❌ 错误示例 1: go generate 指令格式错误
// //go:generate  // 缺少命令
// //go:generate unknown-tool  // 工具不存在

// ✅ 正确做法：完整的 generate 指令
//go:generate go run cmd/generator/main.go -input=$GOFILE -output=generated.go

// ❌ 错误示例 2: 生成代码中手动修改
// 生成的代码文件顶部应该有明确的注释：
// // Code generated by tool; DO NOT EDIT.

// ❌ 错误示例 3: 模板语法错误
// const badTemplate = `{{.Field`  // 缺少闭合括号

// ❌ 错误示例 4: 模板数据类型不匹配
// type Data struct { Name string }
// tmpl.Execute(w, 123)  // 数据类型不匹配

// ❌ 错误示例 5: AST 操作后未正确格式化输出
// 修改 AST 后应该使用 format.Node 输出

// ❌ 错误示例 6: 递归遍历 AST 时无限循环
// ast.Inspect 中返回 false 会停止遍历该分支

// ❌ 错误示例 7: 忽略解析错误
// node, _ := parser.ParseFile(...)  // 忽略错误

// ✅ 正确做法：正确处理错误
// node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
// if err != nil {
//     return err
// }

func main() {
    // 元编程示例
    println("Metaprogramming examples")
}
```

---

## Go 1.26.1 新特性深度分析

### 1. new() 表达式支持

#### 1.1 概念定义

Go 1.26 引入了 `new()` 内置函数接受表达式的能力，允许直接在 `new()` 中传递类型转换或函数调用结果，简化可选指针字段的初始化。

#### 1.2 属性特征

```
┌─────────────────────────────────────────────────────────────┐
│                  new() 表达式特性                            │
├─────────────────────────────────────────────────────────────┤
│  之前语法 (Go 1.25 及以前)                                   │
│  ├─ new(T) 只接受类型标识符                                 │
│  ├─ ptr := new(int)                                         │
│  └─ *ptr = 300  // 需要单独赋值                             │
├─────────────────────────────────────────────────────────────┤
│  新语法 (Go 1.26+)                                          │
│  ├─ new() 接受表达式                                        │
│  ├─ ptr := new(int64(300))                                  │
│  ├─ ptr := new(someFunc())                                  │
│  └─ 简化可选指针字段初始化                                  │
├─────────────────────────────────────────────────────────────┤
│  适用场景                                                    │
│  ├─ 可选配置结构体初始化                                    │
│  ├─ 函数返回值的指针获取                                    │
│  ├─ 类型转换后的指针创建                                    │
│  └─ 链式调用简化                                            │
└─────────────────────────────────────────────────────────────┘
```

#### 1.3 详细示例代码

```go
package main

import (
    "fmt"
    "time"
)

// ========== 基本用法 ==========

func basicNewExpression() {
    // Go 1.26 之前
    oldWay := new(int64)
    *oldWay = 300

    // Go 1.26 新方式
    newWay := new(int64(300))

    fmt.Println("Old way:", *oldWay) // 300
    fmt.Println("New way:", *newWay) // 300
}

// ========== 配置结构体初始化 ==========

// Config 应用配置
type Config struct {
    // 可选配置字段（指针类型）
    Timeout     *time.Duration
    MaxRetries  *int
    Concurrency *int
    EnableCache *bool
    LogLevel    *string
}

// 选项模式辅助函数
func Duration(d time.Duration) *time.Duration {
    return &d
}

func Int(i int) *int {
    return &i
}

func Bool(b bool) *bool {
    return &b
}

func String(s string) *string {
    return &s
}

// Go 1.26 之前的初始化方式
func oldConfigInit() *Config {
    timeout := 30 * time.Second
    maxRetries := 3
    enableCache := true
    logLevel := "info"

    return &Config{
        Timeout:     &timeout,
        MaxRetries:  &maxRetries,
        EnableCache: &enableCache,
        LogLevel:    &logLevel,
    }
}

// Go 1.26 新的初始化方式
func newConfigInit() *Config {
    return &Config{
        Timeout:     new(time.Duration(30 * time.Second)),
        MaxRetries:  new(int(3)),
        EnableCache: new(bool(true)),
        LogLevel:    new(string("info")),
    }
}

// 更简洁的辅助函数方式（Go 1.26 配合）
func NewConfig(options ...ConfigOption) *Config {
    c := &Config{}
    for _, opt := range options {
        opt(c)
    }
    return c
}

type ConfigOption func(*Config)

func WithTimeout(d time.Duration) ConfigOption {
    return func(c *Config) {
        c.Timeout = new(time.Duration(d))
    }
}

func WithMaxRetries(n int) ConfigOption {
    return func(c *Config) {
        c.MaxRetries = new(int(n))
    }
}

func WithEnableCache(b bool) ConfigOption {
    return func(c *Config) {
        c.EnableCache = new(bool(b))
    }
}

// ========== 函数返回值指针 ==========

func getDefaultTimeout() time.Duration {
    return 30 * time.Second
}

func getMaxConnections() int {
    return 100
}

// Go 1.26 可以直接获取函数返回值的指针
func functionReturnPointer() {
    // 之前需要先调用函数，再取地址
    timeout := getDefaultTimeout()
    timeoutPtr := &timeout

    // Go 1.26 新方式
    timeoutPtrNew := new(getDefaultTimeout())

    fmt.Println("Timeout ptr:", *timeoutPtr)
    fmt.Println("Timeout ptr new:", *timeoutPtrNew)

    // 更复杂的表达式
    maxConnPtr := new(getMaxConnections() * 2)
    fmt.Println("Max conn ptr:", *maxConnPtr) // 200
}

// ========== 类型转换场景 ==========

func typeConversionExamples() {
    // 数值类型转换
    intPtr := new(int64(42))
    fmt.Println("Int64 ptr:", *intPtr)

    // 浮点数转换
    floatPtr := new(float64(3.14159))
    fmt.Println("Float64 ptr:", *floatPtr)

    // 字符串转换（从其他类型）
    type MyString string
    strPtr := new(MyString("hello"))
    fmt.Println("MyString ptr:", *strPtr)

    // 复杂表达式
    computedPtr := new(int(10 + 20*3))
    fmt.Println("Computed ptr:", *computedPtr) // 70
}

// ========== 实际应用：API 请求配置 ==========

type APIRequest struct {
    URL         string
    Method      string
    Headers     map[string]string
    Body        []byte
    Timeout     *time.Duration
    RetryCount  *int
    RetryDelay  *time.Duration
}

// RequestOption 请求选项
type RequestOption func(*APIRequest)

func NewAPIRequest(url string, opts ...RequestOption) *APIRequest {
    req := &APIRequest{
        URL:     url,
        Method:  "GET",
        Headers: make(map[string]string),
    }

    for _, opt := range opts {
        opt(req)
    }

    return req
}

func WithMethod(method string) RequestOption {
    return func(r *APIRequest) {
        r.Method = method
    }
}

func WithTimeout(d time.Duration) RequestOption {
    return func(r *APIRequest) {
        r.Timeout = new(time.Duration(d))
    }
}

func WithRetry(count int, delay time.Duration) RequestOption {
    return func(r *APIRequest) {
        r.RetryCount = new(int(count))
        r.RetryDelay = new(time.Duration(delay))
    }
}

func WithHeader(key, value string) RequestOption {
    return func(r *APIRequest) {
        if r.Headers == nil {
            r.Headers = make(map[string]string)
        }
        r.Headers[key] = value
    }
}

// ========== 数据库连接配置 ==========

type DBConfig struct {
    Host        string
    Port        int
    Database    string
    Username    string
    Password    string

    // 可选配置
    MaxOpenConns    *int
    MaxIdleConns    *int
    ConnMaxLifetime *time.Duration
    ConnMaxIdleTime *time.Duration
    SSLMode         *string
}

func NewDBConfig(host string, port int, database, username, password string) *DBConfig {
    return &DBConfig{
        Host:     host,
        Port:     port,
        Database: database,
        Username: username,
        Password: password,
    }
}

func (c *DBConfig) WithPoolSettings(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) *DBConfig {
    c.MaxOpenConns = new(int(maxOpen))
    c.MaxIdleConns = new(int(maxIdle))
    c.ConnMaxLifetime = new(time.Duration(maxLifetime))
    c.ConnMaxIdleTime = new(time.Duration(maxIdleTime))
    return c
}

func (c *DBConfig) WithSSL(mode string) *DBConfig {
    c.SSLMode = new(string(mode))
    return c
}

// ========== 测试用例配置 ==========

type TestConfig struct {
    Name        string
    Parallel    *bool
    Timeout     *time.Duration
    Skip        *bool
    Verbose     *bool
}

func NewTest(name string) *TestConfig {
    return &TestConfig{Name: name}
}

func (t *TestConfig) Parallel() *TestConfig {
    t.Parallel = new(bool(true))
    return t
}

func (t *TestConfig) WithTimeout(d time.Duration) *TestConfig {
    t.Timeout = new(time.Duration(d))
    return t
}

func (t *TestConfig) Skip() *TestConfig {
    t.Skip = new(bool(true))
    return t
}

func (t *TestConfig) Verbose() *TestConfig {
    t.Verbose = new(bool(true))
    return t
}

func main() {
    fmt.Println("=== Basic new() Expression ===")
    basicNewExpression()

    fmt.Println("\n=== Config Initialization ===")
    oldCfg := oldConfigInit()
    newCfg := newConfigInit()
    fmt.Printf("Old config: %+v\n", oldCfg)
    fmt.Printf("New config: %+v\n", newCfg)

    fmt.Println("\n=== Function Return Pointer ===")
    functionReturnPointer()

    fmt.Println("\n=== Type Conversion ===")
    typeConversionExamples()

    fmt.Println("\n=== API Request ===")
    req := NewAPIRequest(
        "https://api.example.com/users",
        WithMethod("POST"),
        WithTimeout(30*time.Second),
        WithRetry(3, 1*time.Second),
        WithHeader("Content-Type", "application/json"),
        WithHeader("Authorization", "Bearer token"),
    )
    fmt.Printf("API Request: %+v\n", req)

    fmt.Println("\n=== DB Config ===")
    dbCfg := NewDBConfig("localhost", 5432, "mydb", "user", "pass").
        WithPoolSettings(100, 10, time.Hour, 10*time.Minute).
        WithSSL("require")
    fmt.Printf("DB Config: %+v\n", dbCfg)

    fmt.Println("\n=== Test Config ===")
    test := NewTest("Unit Test").
        Parallel().
        WithTimeout(5 * time.Minute).
        Verbose()
    fmt.Printf("Test Config: %+v\n", test)
}
```

#### 1.4 反例说明

```go
package main

// ❌ 错误示例 1: new() 不能接受变量
// func badNew() {
//     x := 100
//     ptr := new(x)  // 编译错误：x 不是类型或表达式
// }

// ❌ 错误示例 2: new() 不能接受赋值表达式
// ptr := new(a = 10)  // 编译错误

// ❌ 错误示例 3: new() 不能接受语句
// ptr := new(if true { 10 })  // 编译错误

// ✅ 正确做法：new() 接受类型转换或函数调用表达式
 func goodNew() {
     // 类型转换表达式
     ptr1 := new(int64(100))

     // 函数调用表达式
     ptr2 := new(someFunc())

     // 复合表达式
     ptr3 := new(int(10 + 20))

     _ = ptr1
     _ = ptr2
     _ = ptr3
 }

func someFunc() int {
    return 42
}

// ❌ 错误示例 4: 表达式类型必须与 new 的参数类型匹配
// ptr := new(int64(int32(100)))  // 这实际上是可以的，因为有类型转换

// ❌ 错误示例 5: 不能用于 nil 类型
// ptr := new(nil)  // 编译错误

func main() {
    // new() 表达式示例
    println("new() expression examples")
}
```

### 2. 自引用泛型约束

#### 2.1 概念定义

Go 1.26 引入了自引用泛型约束（Self-Referential Generic Constraints），允许泛型类型参数引用包含它的类型本身。这实现了 F-bounded 多态性，支持递归类型定义和更复杂的类型约束。

#### 2.2 属性特征

```
┌─────────────────────────────────────────────────────────────┐
│                自引用泛型约束特性                            │
├─────────────────────────────────────────────────────────────┤
│  F-bounded 多态性                                            │
│  ├─ 类型参数约束引用自身                                      │
│  ├─ 支持递归类型定义                                          │
│  ├─ 实现类型安全的递归操作                                    │
│  └─ 简化复杂数据结构和算法                                    │
├─────────────────────────────────────────────────────────────┤
│  语法形式                                                    │
│  ├─ type Adder[A Adder[A]] interface { Add(A) A }           │
│  ├─ type Node[T Node[T]] interface { Children() []T }        │
│  └─ type Comparable[T Comparable[T]] interface { ... }        │
├─────────────────────────────────────────────────────────────┤
│  应用场景                                                    │
│  ├─ 树形数据结构                                             │
│  ├─ 图算法                                                   │
│  ├─ 比较器接口                                               │
│  ├─ 数学运算抽象                                             │
│  └─ 递归算法泛化                                             │
└─────────────────────────────────────────────────────────────┘
```

#### 2.3 关系依赖

```
                    自引用约束
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   F-bounded      递归类型       类型安全
   多态性         定义           递归操作
        │              │              │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │         │    │         │    │         │
 Adder[A]  Node[T] 树结构   图结构 比较器   运算器
 Comparable 表达式  二叉树   有向图  Compare  Add
```

#### 2.4 详细示例代码

```go
package main

import (
    "fmt"
)

// ========== 基础自引用约束 ==========

// Adder 自引用加法接口
// A 必须是实现了 Adder[A] 的类型
type Adder[A Adder[A]] interface {
    Add(other A) A
    Zero() A
}

// IntAdder 整数加法实现
type IntAdder int

func (a IntAdder) Add(other IntAdder) IntAdder {
    return a + other
}

func (a IntAdder) Zero() IntAdder {
    return 0
}

// FloatAdder 浮点数加法实现
type FloatAdder float64

func (a FloatAdder) Add(other FloatAdder) FloatAdder {
    return a + other
}

func (a FloatAdder) Zero() FloatAdder {
    return 0
}

// 泛型求和函数
func SumAdder[A Adder[A]](items []A) A {
    var sum A
    for _, item := range items {
        sum = sum.Add(item)
    }
    return sum
}

// ========== 比较器接口 ==========

// Comparator 自引用比较器接口
type Comparator[C Comparator[C]] interface {
    Compare(other C) int  // 返回 -1, 0, 1
    Equal(other C) bool
}

// IntComparator 整数比较器
type IntComparator int

func (c IntComparator) Compare(other IntComparator) int {
    if c < other {
        return -1
    }
    if c > other {
        return 1
    }
    return 0
}

func (c IntComparator) Equal(other IntComparator) bool {
    return c == other
}

// StringComparator 字符串比较器
type StringComparator string

func (c StringComparator) Compare(other StringComparator) int {
    if c < other {
        return -1
    }
    if c > other {
        return 1
    }
    return 0
}

func (c StringComparator) Equal(other StringComparator) bool {
    return c == other
}

// 泛型查找最大值
func MaxComparator[C Comparator[C]](items []C) (C, bool) {
    var zero C
    if len(items) == 0 {
        return zero, false
    }

    max := items[0]
    for _, item := range items[1:] {
        if item.Compare(max) > 0 {
            max = item
        }
    }
    return max, true
}

// 泛型二分查找
func BinarySearch[C Comparator[C]](items []C, target C) int {
    left, right := 0, len(items)-1

    for left <= right {
        mid := left + (right-left)/2
        cmp := items[mid].Compare(target)

        if cmp == 0 {
            return mid
        } else if cmp < 0 {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }

    return -1
}

// ========== 树形数据结构 ==========

// TreeNode 自引用树节点接口
type TreeNode[T TreeNode[T]] interface {
    GetValue() int
    GetChildren() []T
    IsLeaf() bool
}

// BinaryTreeNode 二叉树节点
type BinaryTreeNode struct {
    Value       int
    Left, Right *BinaryTreeNode
}

func (n *BinaryTreeNode) GetValue() int {
    if n == nil {
        return 0
    }
    return n.Value
}

func (n *BinaryTreeNode) GetChildren() []*BinaryTreeNode {
    if n == nil {
        return nil
    }
    children := []*BinaryTreeNode{}
    if n.Left != nil {
        children = append(children, n.Left)
    }
    if n.Right != nil {
        children = append(children, n.Right)
    }
    return children
}

func (n *BinaryTreeNode) IsLeaf() bool {
    if n == nil {
        return false
    }
    return n.Left == nil && n.Right == nil
}

// NaryTreeNode N叉树节点
type NaryTreeNode struct {
    Value    int
    Children []*NaryTreeNode
}

func (n *NaryTreeNode) GetValue() int {
    if n == nil {
        return 0
    }
    return n.Value
}

func (n *NaryTreeNode) GetChildren() []*NaryTreeNode {
    if n == nil {
        return nil
    }
    return n.Children
}

func (n *NaryTreeNode) IsLeaf() bool {
    if n == nil {
        return false
    }
    return len(n.Children) == 0
}

// 泛型树遍历
func TraverseTree[T TreeNode[T]](root T, visit func(int)) {
    if root == nil {
        return
    }

    visit(root.GetValue())

    for _, child := range root.GetChildren() {
        TraverseTree(child, visit)
    }
}

// 泛型树深度计算
func TreeDepth[T TreeNode[T]](node T) int {
    if node == nil || node.IsLeaf() {
        return 0
    }

    maxDepth := 0
    for _, child := range node.GetChildren() {
        depth := TreeDepth(child)
        if depth > maxDepth {
            maxDepth = depth
        }
    }

    return maxDepth + 1
}

// ========== 表达式求值系统 ==========

// Expr 自引用表达式接口
type Expr[E Expr[E]] interface {
    Eval() float64
    String() string
    Children() []E
}

// Literal 字面量表达式
type Literal float64

func (l Literal) Eval() float64 {
    return float64(l)
}

func (l Literal) String() string {
    return fmt.Sprintf("%.2f", l)
}

func (l Literal) Children() []Expr[Literal] {
    return nil
}

// BinaryOp 二元运算表达式
type BinaryOp struct {
    Op          string
    Left, Right Expr[Expr[any]]
}

func (b *BinaryOp) Eval() float64 {
    left := b.Left.Eval()
    right := b.Right.Eval()

    switch b.Op {
    case "+":
        return left + right
    case "-":
        return left - right
    case "*":
        return left * right
    case "/":
        if right != 0 {
            return left / right
        }
        return 0
    default:
        return 0
    }
}

func (b *BinaryOp) String() string {
    return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Op, b.Right.String())
}

func (b *BinaryOp) Children() []Expr[Expr[any]] {
    return []Expr[Expr[any]]{b.Left, b.Right}
}

// ========== 图算法 ==========

// GraphNode 自引用图节点接口
type GraphNode[G GraphNode[G]] interface {
    GetID() string
    GetNeighbors() []G
    GetWeight(to G) float64
}

// SimpleGraphNode 简单图节点实现
type SimpleGraphNode struct {
    ID        string
    Neighbors map[string]*SimpleGraphNode
    Weights   map[string]float64
}

func (n *SimpleGraphNode) GetID() string {
    return n.ID
}

func (n *SimpleGraphNode) GetNeighbors() []*SimpleGraphNode {
    result := make([]*SimpleGraphNode, 0, len(n.Neighbors))
    for _, neighbor := range n.Neighbors {
        result = append(result, neighbor)
    }
    return result
}

func (n *SimpleGraphNode) GetWeight(to *SimpleGraphNode) float64 {
    if n.Weights == nil {
        return 0
    }
    return n.Weights[to.ID]
}

func (n *SimpleGraphNode) AddNeighbor(neighbor *SimpleGraphNode, weight float64) {
    if n.Neighbors == nil {
        n.Neighbors = make(map[string]*SimpleGraphNode)
        n.Weights = make(map[string]float64)
    }
    n.Neighbors[neighbor.ID] = neighbor
    n.Weights[neighbor.ID] = weight
}

// 泛型深度优先搜索
func DFS[G GraphNode[G]](start G, visit func(G)) {
    visited := make(map[string]bool)

    var dfs func(G)
    dfs = func(node G) {
        id := node.GetID()
        if visited[id] {
            return
        }

        visited[id] = true
        visit(node)

        for _, neighbor := range node.GetNeighbors() {
            dfs(neighbor)
        }
    }

    dfs(start)
}

// ========== 数学向量空间 ==========

// Vector 自引用向量接口
type Vector[V Vector[V]] interface {
    Add(other V) V
    Scale(scalar float64) V
    Dot(other V) float64
    Zero() V
}

// FloatVector 浮点数向量
type FloatVector []float64

func (v FloatVector) Add(other FloatVector) FloatVector {
    if len(v) != len(other) {
        panic("vector dimensions do not match")
    }

    result := make(FloatVector, len(v))
    for i := range v {
        result[i] = v[i] + other[i]
    }
    return result
}

func (v FloatVector) Scale(scalar float64) FloatVector {
    result := make(FloatVector, len(v))
    for i := range v {
        result[i] = v[i] * scalar
    }
    return result
}

func (v FloatVector) Dot(other FloatVector) float64 {
    if len(v) != len(other) {
        panic("vector dimensions do not match")
    }

    sum := 0.0
    for i := range v {
        sum += v[i] * other[i]
    }
    return sum
}

func (v FloatVector) Zero() FloatVector {
    return make(FloatVector, len(v))
}

// 泛型线性组合
func LinearCombination[V Vector[V]](vectors []V, scalars []float64) V {
    if len(vectors) == 0 {
        var zero V
        return zero
    }

    result := vectors[0].Zero()
    for i, v := range vectors {
        result = result.Add(v.Scale(scalars[i]))
    }
    return result
}

// ========== 使用示例 ==========

func main() {
    fmt.Println("=== Adder Interface ===")
    intItems := []IntAdder{1, 2, 3, 4, 5}
    intSum := SumAdder(intItems)
    fmt.Printf("Sum of %v = %d\n", intItems, intSum)

    floatItems := []FloatAdder{1.5, 2.5, 3.5}
    floatSum := SumAdder(floatItems)
    fmt.Printf("Sum of %v = %.2f\n", floatItems, floatSum)

    fmt.Println("\n=== Comparator Interface ===")
    intComps := []IntComparator{5, 2, 8, 1, 9, 3}
    if max, ok := MaxComparator(intComps); ok {
        fmt.Printf("Max of %v = %d\n", intComps, max)
    }

    strComps := []StringComparator{"apple", "banana", "cherry", "date"}
    if maxStr, ok := MaxComparator(strComps); ok {
        fmt.Printf("Max of %v = %s\n", strComps, maxStr)
    }

    sortedInts := []IntComparator{1, 3, 5, 7, 9, 11, 13}
    idx := BinarySearch(sortedInts, IntComparator(7))
    fmt.Printf("Binary search for 7: index = %d\n", idx)

    fmt.Println("\n=== Tree Traversal ===")
    // 构建二叉树
    //       1
    //      / \
    //     2   3
    //    / \
    //   4   5
    root := &BinaryTreeNode{
        Value: 1,
        Left: &BinaryTreeNode{
            Value: 2,
            Left:  &BinaryTreeNode{Value: 4},
            Right: &BinaryTreeNode{Value: 5},
        },
        Right: &BinaryTreeNode{Value: 3},
    }

    fmt.Print("Tree traversal: ")
    TraverseTree(root, func(v int) {
        fmt.Printf("%d ", v)
    })
    fmt.Println()

    depth := TreeDepth(root)
    fmt.Printf("Tree depth: %d\n", depth)

    fmt.Println("\n=== Expression Evaluation ===")
    // 构建表达式: (1 + 2) * 3
    expr := &BinaryOp{
        Op: "*",
        Left: &BinaryOp{
            Op:    "+",
            Left:  Literal(1),
            Right: Literal(2),
        },
        Right: Literal(3),
    }

    fmt.Printf("Expression: %s\n", expr.String())
    fmt.Printf("Result: %.2f\n", expr.Eval())

    fmt.Println("\n=== Graph DFS ===")
    // 构建图
    nodeA := &SimpleGraphNode{ID: "A"}
    nodeB := &SimpleGraphNode{ID: "B"}
    nodeC := &SimpleGraphNode{ID: "C"}
    nodeD := &SimpleGraphNode{ID: "D"}

    nodeA.AddNeighbor(nodeB, 1.0)
    nodeA.AddNeighbor(nodeC, 2.0)
    nodeB.AddNeighbor(nodeD, 3.0)
    nodeC.AddNeighbor(nodeD, 1.0)

    fmt.Print("DFS from A: ")
    DFS(nodeA, func(n *SimpleGraphNode) {
        fmt.Printf("%s ", n.GetID())
    })
    fmt.Println()

    fmt.Println("\n=== Vector Operations ===")
    v1 := FloatVector{1, 2, 3}
    v2 := FloatVector{4, 5, 6}

    fmt.Printf("v1 = %v\n", v1)
    fmt.Printf("v2 = %v\n", v2)
    fmt.Printf("v1 + v2 = %v\n", v1.Add(v2))
    fmt.Printf("v1 · v2 = %.2f\n", v1.Dot(v2))
    fmt.Printf("2 * v1 = %v\n", v1.Scale(2))

    // 线性组合
    vectors := []FloatVector{v1, v2}
    scalars := []float64{2.0, 3.0}
    combo := LinearCombination(vectors, scalars)
    fmt.Printf("2*v1 + 3*v2 = %v\n", combo)
}
```

#### 2.5 反例说明

```go
package main

// ❌ 错误示例 1: 自引用约束必须正确实现
// type BadAdder int
// func (a BadAdder) Add(other BadAdder) BadAdder { return a + other }
// 缺少 Zero() 方法，不满足 Adder[A] 约束

// ✅ 正确做法：实现所有约束方法
 type GoodAdder int
 func (a GoodAdder) Add(other GoodAdder) GoodAdder { return a + other }
 func (a GoodAdder) Zero() GoodAdder { return 0 }

// ❌ 错误示例 2: 自引用约束不能循环依赖
// type A[T B[T]] interface {}
// type B[T A[T]] interface {}
// 这种循环依赖会导致编译问题

// ❌ 错误示例 3: 类型参数不匹配
// type Node[T Node[T]] interface { GetChildren() []T }
// type MyNode struct {}
// func (n MyNode) GetChildren() []OtherNode { return nil }  // 类型不匹配

// ✅ 正确做法：确保方法签名完全匹配
 type MyNode struct {}
 func (n MyNode) GetChildren() []MyNode { return nil }

// ❌ 错误示例 4: 不能对非接口类型使用自引用
// type BadStruct[T BadStruct[T]] struct {}  // 编译错误

// ✅ 正确做法：自引用约束只用于接口类型
 type GoodInterface[T GoodInterface[T]] interface {
     Method() T
 }

// ❌ 错误示例 5: 泛型函数的类型推断限制
// func Process[T Node[T]](n T) {}
// 在某些复杂情况下，可能需要显式指定类型参数

// ❌ 错误示例 6: 自引用接口不能用于具体类型变量声明
// var node Node[MyNode]  // 这可能无法编译

// ✅ 正确做法：使用具体类型或泛型约束
 func ProcessNode[T Node[T]](n T) {
     // 使用泛型约束
 }

func main() {
    // 自引用约束示例
    println("Self-referential constraint examples")
}
```

---

## CGO 和外部函数接口

### 1. 概念定义

CGO（C Go）是 Go 语言提供的与 C 语言互操作的机制，允许 Go 程序调用 C 代码，也允许 C 代码调用 Go 函数。CGO 通过特殊的注释语法和运行时支持实现两种语言之间的无缝集成。

### 2. 属性特征

```
┌─────────────────────────────────────────────────────────────┐
│                    CGO 核心组件                              │
├─────────────────────────────────────────────────────────────┤
│  特殊注释语法                                                │
│  ├─ import "C" 启用 CGO                                     │
│  ├─ // #include <header.h>                                  │
│  ├─ // #cgo CFLAGS: -I/path                                │
│  └─ // #cgo LDFLAGS: -L/path -lname                        │
├─────────────────────────────────────────────────────────────┤
│  类型转换                                                    │
│  ├─ C.int ↔ C.int                                           │
│  ├─ C.char ↔ C.char                                         │
│  ├─ C.CString() / C.GoString()                              │
│  └─ unsafe.Pointer 用于复杂类型                             │
├─────────────────────────────────────────────────────────────┤
│  内存管理                                                    │
│  ├─ C 内存：C.malloc() / C.free()                          │
│  ├─ Go 内存：由 GC 管理                                     │
│  ├─ 字符串转换需要手动释放                                  │
│  └─ 指针传递注意事项                                        │
├─────────────────────────────────────────────────────────────┤
│  回调机制                                                    │
│  ├─ Go 函数导出给 C 调用                                    │
│  ├─ //export 注释                                           │
│  └─ 函数签名限制                                            │
├─────────────────────────────────────────────────────────────┤
│  并发安全                                                    │
│  ├─ C 代码可能阻塞 Go 调度器                                │
│  ├─ 长时间运行的 C 函数影响性能                             │
│  └─ 使用 Goroutine 隔离                                     │
└─────────────────────────────────────────────────────────────┘
```

### 3. 关系依赖

```
                    CGO
                       │
        ┌──────────────┼──────────────┐
        │              │              │
    Go代码调用C    C代码调用Go    类型转换
        │              │              │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │         │    │         │    │         │
 import "C" C函数调用 //export  导出  C类型   Go类型
 注释包含  直接调用  函数签名  回调注册  转换函数  内存管理
```

### 4. 详细示例代码

#### 4.1 基础 CGO 使用

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 简单 C 函数
int add(int a, int b) {
    return a + b;
}

// 字符串处理
void greet(const char* name) {
    printf("Hello, %s!\n", name);
}

// 返回字符串
const char* get_version() {
    return "C Library v1.0";
}

// 结构体定义
typedef struct {
    int x;
    int y;
} Point;

// 使用结构体
int point_distance(Point p1, Point p2) {
    int dx = p1.x - p2.x;
    int dy = p1.y - p2.y;
    return dx * dx + dy * dy;
}
*/
import "C"
import (
    "fmt"
    "unsafe"
)

func basicCGO() {
    // 调用 C 函数
    result := C.add(C.int(10), C.int(20))
    fmt.Printf("10 + 20 = %d\n", int(result))

    // 字符串传递
    name := C.CString("World")
    defer C.free(unsafe.Pointer(name)) // 重要：释放 C 字符串
    C.greet(name)

    // 接收 C 字符串
    version := C.GoString(C.get_version())
    fmt.Printf("Version: %s\n", version)
}

func structExample() {
    // 使用 C 结构体
    p1 := C.Point{x: 0, y: 0}
    p2 := C.Point{x: 3, y: 4}

    dist := C.point_distance(p1, p2)
    fmt.Printf("Distance squared: %d\n", int(dist))
}

func main() {
    fmt.Println("=== Basic CGO ===")
    basicCGO()

    fmt.Println("\n=== Struct Example ===")
    structExample()
}
```

#### 4.2 内存管理和字符串处理

```go
package main

/*
#include <stdlib.h>
#include <string.h>

// 分配内存
char* allocate_string(size_t size) {
    return (char*)malloc(size);
}

// 复制字符串
char* duplicate_string(const char* src) {
    return strdup(src);
}

// 释放内存
void free_memory(void* ptr) {
    free(ptr);
}

// 修改字符串
void to_uppercase(char* str) {
    for (int i = 0; str[i]; i++) {
        if (str[i] >= 'a' && str[i] <= 'z') {
            str[i] -= 32;
        }
    }
}
*/
import "C"
import (
    "fmt"
    "unsafe"
)

// 安全的字符串转换
func cStringToGoString(cstr *C.char) string {
    return C.GoString(cstr)
}

func goStringToCString(gostr string) *C.char {
    return C.CString(gostr)
}

// 内存分配示例
func memoryAllocationExample() {
    // 从 C 分配内存
    size := C.size_t(100)
    cstr := C.allocate_string(size)
    defer C.free_memory(unsafe.Pointer(cstr))

    // 复制 Go 字符串到 C 内存
    gostr := "Hello from Go"
    cstrCopy := C.CString(gostr)
    defer C.free(unsafe.Pointer(cstrCopy))

    C.strcpy(cstr, cstrCopy)

    // 读取结果
    result := C.GoString(cstr)
    fmt.Printf("C memory content: %s\n", result)
}

// 字符串修改示例
func stringModification() {
    // 创建可修改的 C 字符串
    gostr := "hello world"
    cstr := C.CString(gostr)
    defer C.free(unsafe.Pointer(cstr))

    fmt.Printf("Before: %s\n", gostr)

    // 在 C 中修改
    C.to_uppercase(cstr)

    // 读取修改后的结果
    result := C.GoString(cstr)
    fmt.Printf("After: %s\n", result)
}

// 批量字符串处理
func batchStringProcessing() {
    strings := []string{"apple", "banana", "cherry"}

    // 转换为 C 字符串数组
    cStrings := make([]*C.char, len(strings))
    for i, s := range strings {
        cStrings[i] = C.CString(s)
    }

    // 确保释放所有内存
    defer func() {
        for _, cs := range cStrings {
            C.free(unsafe.Pointer(cs))
        }
    }()

    // 处理每个字符串
    for i, cs := range cStrings {
        C.to_uppercase(cs)
        fmt.Printf("String %d: %s -> %s\n", i, strings[i], C.GoString(cs))
    }
}

// 安全的字节数组传递
func byteArrayExample() {
    data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

    // 转换为 C 数组
    cData := C.CBytes(data)
    defer C.free(cData)

    // 可以传递给 C 函数处理
    // C.process_bytes((*C.uchar)(cData), C.size_t(len(data)))

    // 转换回 Go 切片
    length := len(data)
    goData := C.GoBytes(cData, C.int(length))

    fmt.Printf("Original: %v\n", data)
    fmt.Printf("From C: %v\n", goData)
}

func main() {
    fmt.Println("=== Memory Allocation ===")
    memoryAllocationExample()

    fmt.Println("\n=== String Modification ===")
    stringModification()

    fmt.Println("\n=== Batch String Processing ===")
    batchStringProcessing()

    fmt.Println("\n=== Byte Array Example ===")
    byteArrayExample()
}
```

#### 4.3 Go 函数导出给 C

```go
package main

/*
#include <stdint.h>
#include <stdlib.h>

// 前向声明 Go 回调函数
typedef int (*callback_func)(int, int);
typedef void (*string_callback)(const char*);

// C 包装函数
int call_go_callback(callback_func fn, int a, int b) {
    return fn(a, b);
}

void call_string_callback(string_callback fn, const char* str) {
    fn(str);
}
*/
import "C"
import (
    "fmt"
    "sync"
    "unsafe"
)

//export Add
func Add(a, b C.int) C.int {
    return a + b
}

//export Multiply
func Multiply(a, b C.int) C.int {
    return a * b
}

//export ProcessString
func ProcessString(cstr *C.char) *C.char {
    gostr := C.GoString(cstr)
    result := fmt.Sprintf("Processed: %s", gostr)
    return C.CString(result)
}

//export FreeString
func FreeString(cstr *C.char) {
    C.free(unsafe.Pointer(cstr))
}

// 回调函数类型
type GoCallback func(int, int) int

var (
    callbacks   = make(map[uintptr]GoCallback)
    callbackMux sync.RWMutex
    callbackID  uintptr
)

//export GoCallbackWrapper
func GoCallbackWrapper(id uintptr, a, b C.int) C.int {
    callbackMux.RLock()
    cb, exists := callbacks[id]
    callbackMux.RUnlock()

    if !exists {
        return -1
    }

    return C.int(cb(int(a), int(b)))
}

func registerCallback(cb GoCallback) uintptr {
    callbackMux.Lock()
    defer callbackMux.Unlock()

    callbackID++
    callbacks[callbackID] = cb
    return callbackID
}

func unregisterCallback(id uintptr) {
    callbackMux.Lock()
    defer callbackMux.Unlock()
    delete(callbacks, id)
}

// 使用 C 调用 Go 回调
func useCallbackExample() {
    // 注册回调
    cb := func(a, b int) int {
        return a * a + b * b
    }

    id := registerCallback(cb)
    defer unregisterCallback(id)

    // 通过 C 调用
    result := C.GoCallbackWrapper(C.uintptr_t(id), 3, 4)
    fmt.Printf("Callback result: %d\n", int(result))
}

// 批量处理示例
//export BatchProcess
func BatchProcess(data *C.int, length C.int) *C.int {
    goData := make([]int, int(length))

    // 从 C 数组读取
    cArray := (*[1 << 30]C.int)(unsafe.Pointer(data))[:length:length]
    for i := 0; i < int(length); i++ {
        goData[i] = int(cArray[i])
    }

    // 处理数据
    for i := range goData {
        goData[i] = goData[i] * 2
    }

    // 分配 C 内存并返回
    result := C.malloc(C.size_t(length) * C.size_t(unsafe.Sizeof(C.int(0))))
    resultArray := (*[1 << 30]C.int)(result)[:length:length]

    for i, v := range goData {
        resultArray[i] = C.int(v)
    }

    return (*C.int)(result)
}

//export FreeBatchResult
func FreeBatchResult(data *C.int) {
    C.free(unsafe.Pointer(data))
}

func main() {
    fmt.Println("=== Go Functions Exported to C ===")

    // 直接调用导出的函数
    sum := Add(10, 20)
    fmt.Printf("Add(10, 20) = %d\n", int(sum))

    product := Multiply(5, 6)
    fmt.Printf("Multiply(5, 6) = %d\n", int(product))

    // 字符串处理
    input := C.CString("Hello")
    defer C.free(unsafe.Pointer(input))

    output := ProcessString(input)
    defer FreeString(output)

    fmt.Printf("ProcessString result: %s\n", C.GoString(output))

    // 回调示例
    fmt.Println("\n=== Callback Example ===")
    useCallbackExample()
}
```

#### 4.4 使用外部 C 库

```go
package main

/*
#cgo pkg-config: sqlite3
#cgo LDFLAGS: -lm

#include <sqlite3.h>
#include <math.h>
#include <stdlib.h>

// 错误处理辅助函数
const char* sqlite3_errmsg_wrapper(sqlite3* db) {
    return sqlite3_errmsg(db);
}

// 查询回调类型
typedef int (*query_callback)(void*, int, char**, char**);
*/
import "C"
import (
    "fmt"
    "unsafe"
)

// SQLite 包装器
type SQLiteDB struct {
    db *C.sqlite3
}

func OpenSQLiteDB(path string) (*SQLiteDB, error) {
    cpath := C.CString(path)
    defer C.free(unsafe.Pointer(cpath))

    var db *C.sqlite3
    result := C.sqlite3_open(cpath, &db)

    if result != C.SQLITE_OK {
        errMsg := C.GoString(C.sqlite3_errmsg(db))
        C.sqlite3_close(db)
        return nil, fmt.Errorf("failed to open database: %s", errMsg)
    }

    return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) Close() {
    C.sqlite3_close(s.db)
}

func (s *SQLiteDB) Execute(sql string) error {
    csql := C.CString(sql)
    defer C.free(unsafe.Pointer(csql))

    result := C.sqlite3_exec(s.db, csql, nil, nil, nil)
    if result != C.SQLITE_OK {
        errMsg := C.GoString(C.sqlite3_errmsg(s.db))
        return fmt.Errorf("SQL error: %s", errMsg)
    }

    return nil
}

func (s *SQLiteDB) Query(sql string) ([]map[string]string, error) {
    csql := C.CString(sql)
    defer C.free(unsafe.Pointer(csql))

    var stmt *C.sqlite3_stmt
    result := C.sqlite3_prepare_v2(s.db, csql, -1, &stmt, nil)

    if result != C.SQLITE_OK {
        errMsg := C.GoString(C.sqlite3_errmsg(s.db))
        return nil, fmt.Errorf("prepare error: %s", errMsg)
    }
    defer C.sqlite3_finalize(stmt)

    var results []map[string]string

    for {
        result = C.sqlite3_step(stmt)
        if result == C.SQLITE_DONE {
            break
        }
        if result != C.SQLITE_ROW {
            return nil, fmt.Errorf("step error: %d", result)
        }

        row := make(map[string]string)
        colCount := C.sqlite3_column_count(stmt)

        for i := 0; i < int(colCount); i++ {
            colName := C.GoString(C.sqlite3_column_name(stmt, C.int(i)))
            colValue := C.GoString((*C.char)(unsafe.Pointer(C.sqlite3_column_text(stmt, C.int(i)))))
            row[colName] = colValue
        }

        results = append(results, row)
    }

    return results, nil
}

// 数学函数包装
func Sqrt(x float64) float64 {
    return float64(C.sqrt(C.double(x)))
}

func Pow(x, y float64) float64 {
    return float64(C.pow(C.double(x), C.double(y)))
}

func Sin(x float64) float64 {
    return float64(C.sin(C.double(x)))
}

func Cos(x float64) float64 {
    return float64(C.cos(C.double(x)))
}

func main() {
    fmt.Println("=== Math Functions ===")
    fmt.Printf("sqrt(16) = %.2f\n", Sqrt(16))
    fmt.Printf("pow(2, 10) = %.2f\n", Pow(2, 10))
    fmt.Printf("sin(pi/2) = %.2f\n", Sin(3.14159/2))
    fmt.Printf("cos(0) = %.2f\n", Cos(0))

    fmt.Println("\n=== SQLite Example ===")

    // 打开内存数据库
    db, err := OpenSQLiteDB(":memory:")
    if err != nil {
        fmt.Printf("Error opening database: %v\n", err)
        return
    }
    defer db.Close()

    // 创建表
    err = db.Execute(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            age INTEGER
        )
    `)
    if err != nil {
        fmt.Printf("Error creating table: %v\n", err)
        return
    }

    // 插入数据
    err = db.Execute(`INSERT INTO users (name, age) VALUES ('Alice', 30)`)
    if err != nil {
        fmt.Printf("Error inserting: %v\n", err)
        return
    }

    err = db.Execute(`INSERT INTO users (name, age) VALUES ('Bob', 25)`)
    if err != nil {
        fmt.Printf("Error inserting: %v\n", err)
        return
    }

    // 查询数据
    results, err := db.Query(`SELECT * FROM users`)
    if err != nil {
        fmt.Printf("Error querying: %v\n", err)
        return
    }

    fmt.Println("Users:")
    for _, row := range results {
        fmt.Printf("  ID: %s, Name: %s, Age: %s\n", row["id"], row["name"], row["age"])
    }
}
```

#### 4.5 复杂类型和结构体

```go
package main

/*
#include <stdlib.h>
#include <string.h>

// 复杂结构体定义
typedef struct {
    char* name;
    int age;
    double salary;
    char** tags;
    int tag_count;
} Employee;

// 创建员工
typedef struct {
    Employee** items;
    int count;
    int capacity;
} EmployeeList;

// 辅助函数
Employee* create_employee(const char* name, int age, double salary) {
    Employee* emp = (Employee*)malloc(sizeof(Employee));
    emp->name = strdup(name);
    emp->age = age;
    emp->salary = salary;
    emp->tags = NULL;
    emp->tag_count = 0;
    return emp;
}

void free_employee(Employee* emp) {
    if (emp) {
        free(emp->name);
        for (int i = 0; i < emp->tag_count; i++) {
            free(emp->tags[i]);
        }
        free(emp->tags);
        free(emp);
    }
}

void add_tag(Employee* emp, const char* tag) {
    emp->tags = (char**)realloc(emp->tags, (emp->tag_count + 1) * sizeof(char*));
    emp->tags[emp->tag_count] = strdup(tag);
    emp->tag_count++;
}

EmployeeList* create_employee_list() {
    EmployeeList* list = (EmployeeList*)malloc(sizeof(EmployeeList));
    list->items = NULL;
    list->count = 0;
    list->capacity = 0;
    return list;
}

void add_to_list(EmployeeList* list, Employee* emp) {
    if (list->count >= list->capacity) {
        list->capacity = list->capacity == 0 ? 4 : list->capacity * 2;
        list->items = (Employee**)realloc(list->items, list->capacity * sizeof(Employee*));
    }
    list->items[list->count++] = emp;
}

void free_employee_list(EmployeeList* list) {
    if (list) {
        for (int i = 0; i < list->count; i++) {
            free_employee(list->items[i]);
        }
        free(list->items);
        free(list);
    }
}
*/
import "C"
import (
    "fmt"
    "unsafe"
)

// Go 版本的 Employee 结构体
type Employee struct {
    Name    string
    Age     int
    Salary  float64
    Tags    []string
}

// 从 C 结构体转换为 Go 结构体
func cEmployeeToGo(cEmp *C.Employee) *Employee {
    if cEmp == nil {
        return nil
    }

    emp := &Employee{
        Name:   C.GoString(cEmp.name),
        Age:    int(cEmp.age),
        Salary: float64(cEmp.salary),
    }

    // 转换标签
    if cEmp.tag_count > 0 {
        tags := (*[1 << 30]*C.char)(unsafe.Pointer(cEmp.tags))[:cEmp.tag_count:cEmp.tag_count]
        emp.Tags = make([]string, cEmp.tag_count)
        for i, tag := range tags {
            emp.Tags[i] = C.GoString(tag)
        }
    }

    return emp
}

// 从 Go 结构体创建 C 结构体
func goEmployeeToC(emp *Employee) *C.Employee {
    cName := C.CString(emp.Name)
    // 注意：cName 需要在 C 代码中释放

    cEmp := C.create_employee(cName, C.int(emp.Age), C.double(emp.Salary))
    C.free(unsafe.Pointer(cName))

    // 添加标签
    for _, tag := range emp.Tags {
        cTag := C.CString(tag)
        C.add_tag(cEmp, cTag)
        C.free(unsafe.Pointer(cTag))
    }

    return cEmp
}

// EmployeeList 包装器
type EmployeeList struct {
    cList *C.EmployeeList
}

func NewEmployeeList() *EmployeeList {
    return &EmployeeList{
        cList: C.create_employee_list(),
    }
}

func (el *EmployeeList) Add(emp *Employee) {
    cEmp := goEmployeeToC(emp)
    C.add_to_list(el.cList, cEmp)
}

func (el *EmployeeList) GetAll() []*Employee {
    count := int(el.cList.count)
    results := make([]*Employee, count)

    items := (*[1 << 30]*C.Employee)(unsafe.Pointer(el.cList.items))[:count:count]
    for i, cEmp := range items {
        results[i] = cEmployeeToGo(cEmp)
    }

    return results
}

func (el *EmployeeList) Free() {
    C.free_employee_list(el.cList)
}

func complexStructExample() {
    // 创建员工列表
    list := NewEmployeeList()
    defer list.Free()

    // 添加员工
    employees := []*Employee{
        {
            Name:   "Alice",
            Age:    30,
            Salary: 50000.0,
            Tags:   []string{"developer", "senior"},
        },
        {
            Name:   "Bob",
            Age:    25,
            Salary: 35000.0,
            Tags:   []string{"developer", "junior"},
        },
        {
            Name:   "Charlie",
            Age:    35,
            Salary: 70000.0,
            Tags:   []string{"manager"},
        },
    }

    for _, emp := range employees {
        list.Add(emp)
    }

    // 获取所有员工
    fmt.Println("Employee List:")
    for _, emp := range list.GetAll() {
        fmt.Printf("  Name: %s, Age: %d, Salary: %.2f, Tags: %v\n",
            emp.Name, emp.Age, emp.Salary, emp.Tags)
    }
}

// 批量处理示例
func batchProcessingExample() {
    // 创建大量员工数据
    list := NewEmployeeList()
    defer list.Free()

    for i := 0; i < 1000; i++ {
        emp := &Employee{
            Name:   fmt.Sprintf("Employee_%d", i),
            Age:    20 + i%40,
            Salary: 30000.0 + float64(i)*100,
            Tags:   []string{"batch", fmt.Sprintf("group_%d", i%10)},
        }
        list.Add(emp)
    }

    allEmployees := list.GetAll()

    // 计算统计数据
    var totalSalary float64
    ageGroups := make(map[int]int)

    for _, emp := range allEmployees {
        totalSalary += emp.Salary
        ageGroups[emp.Age/10*10]++
    }

    fmt.Printf("\nBatch Processing Results:\n")
    fmt.Printf("  Total employees: %d\n", len(allEmployees))
    fmt.Printf("  Average salary: %.2f\n", totalSalary/float64(len(allEmployees)))
    fmt.Println("  Age groups:")
    for age, count := range ageGroups {
        fmt.Printf("    %d-%d: %d\n", age, age+9, count)
    }
}

func main() {
    fmt.Println("=== Complex Struct Example ===")
    complexStructExample()

    fmt.Println("\n=== Batch Processing Example ===")
    batchProcessingExample()
}
```

### 5. 反例说明

```go
package main

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// ❌ 错误示例 1: 内存泄漏 - 忘记释放 C 字符串
// func badStringHandling() {
//     cstr := C.CString("hello")  // 分配了 C 内存
//     // 忘记 defer C.free(unsafe.Pointer(cstr))
//     // 内存泄漏！
// }

// ✅ 正确做法：使用 defer 确保释放
 func goodStringHandling() {
     cstr := C.CString("hello")
     defer C.free(unsafe.Pointer(cstr))
     // 安全使用 cstr
 }

// ❌ 错误示例 2: 在 C 中存储 Go 指针
// C 代码不能存储 Go 指针，因为 Go GC 会移动对象

// ❌ 错误示例 3: 导出函数使用不支持的类型
// //export BadFunction
// func BadFunction(s string) {  // 编译错误：不能使用 Go 字符串类型
// }

// ✅ 正确做法：使用 C 类型
// //export GoodFunction
// func GoodFunction(cstr *C.char) {
//     gostr := C.GoString(cstr)
//     // 处理 gostr
// }

// ❌ 错误示例 4: 并发访问 C 全局状态
// C 代码可能不是线程安全的，需要同步

// ❌ 错误示例 5: 长时间运行的 C 函数阻塞调度器
// func badLongRunning() {
//     C.long_running_function()  // 会阻塞 Go 调度器
// }

// ✅ 正确做法：在单独的 Goroutine 中运行
// func goodLongRunning() {
//     go func() {
//         C.long_running_function()
//     }()
// }

// ❌ 错误示例 6: 类型转换错误
// func badCast() {
//     var x int = 100
//     cPtr := (*C.int)(unsafe.Pointer(&x))  // 不安全
//     _ = cPtr
// }

// ✅ 正确做法：使用正确的 C 类型
// func goodCast() {
//     var x C.int = 100
//     cPtr := &x
//     _ = cPtr
// }

// ❌ 错误示例 7: 修改 C 结构体中的 Go 指针
// C 结构体中不应包含 Go 指针

// ❌ 错误示例 8: 在 C 回调中调用 Go 函数
// 需要确保 Go 函数在 C 回调期间保持有效

func main() {
    // CGO 示例
    println("CGO examples")
}
```

---

## 决策树图分析

### 1. 泛型使用决策树

```
                    是否需要泛型？
                         │
            ┌────────────┴────────────┐
            │                         │
           否                        是
            │                         │
        使用具体类型            需要类型参数？
                                     │
                    ┌────────────────┴────────────────┐
                    │                                 │
                   否                                是
                    │                                 │
            使用接口类型                    需要多个类型参数？
                                                     │
                                    ┌────────────────┴────────────────┐
                                    │                                 │
                                   否                                是
                                    │                                 │
                            单类型参数                      定义多个类型参数
                            [T Constraint]                  [K comparable, V any]
                                    │                                 │
                                    └────────────────┬────────────────┘
                                                     │
                                            需要类型约束？
                                                     │
                                    ┌────────────────┴────────────────┐
                                    │                                 │
                                   否                                是
                                    │                                 │
                            使用 any 约束                    自定义约束接口
                            [T any]                          type Number interface {
                                                                    ~int | ~float64
                                                                 }
                                                     │
                                                     ▼
                                            需要自引用约束？
                                                     │
                                    ┌────────────────┴────────────────┐
                                    │                                 │
                                   否                                是
                                    │                                 │
                            标准泛型约束                    使用自引用约束
                            [T Constraint]                 type Node[T Node[T]] interface {
                                                                    GetChildren() []T
                                                                 }
```

### 2. 接口设计决策树

```
                    设计接口？
                         │
            ┌────────────┴────────────┐
            │                         │
        具体类型                    接口类型
            │                         │
    定义结构体/类型          需要泛型接口？
                                     │
                    ┌────────────────┴────────────────┐
                    │                                 │
                   否                                是
                    │                                 │
            基本接口                        定义泛型接口
            type Reader interface {         type Container[T any] interface {
                Read([]byte) (int, error)       Add(T)
            }                                       Remove(T) bool
                                            }
                    │                                 │
                    └────────────────┬────────────────┘
                                     │
                            需要组合接口？
                                     │
                    ┌────────────────┴────────────────┐
                    │                                 │
                   否                                是
                    │                                 │
            单一接口                        组合接口
            type Writer interface {         type ReadWriter interface {
                Write([]byte) (int, error)      Reader
            }                                       Writer
                                            }
                                     │
                                     ▼
                            需要类型集约束？
                                     │
                    ┌────────────────┴────────────────┐
                    │                                 │
                   否                                是
                    │                                 │
            方法约束                        类型集约束
            type Stringer interface {       type Number interface {
                String() string                 ~int | ~float64
            }                                   ~int8 | ~int16
```

### 3. 反射使用决策树

```
                    需要反射？
                         │
            ┌────────────┴────────────┐
            │                         │
           否                        是
            │                         │
        使用直接调用            需要获取类型信息？
                                     │
                    ┌────────────────┴────────────────┐
                    │                                 │
                   否                                是
                    │                                 │
            值操作                          使用 reflect.TypeOf()
            reflect.ValueOf()               获取类型信息
                    │                         │
                    │                    需要结构体信息？
                    │                             │
                    │            ┌────────────────┴────────────────┐
                    │            │                                 │
                    │           否                                是
                    │            │                                 │
                    │    基本类型操作                    结构体反射
                    │    v.Int(), v.String()            t.NumField()
                    │                                    v.Field(i)
                    │                                         │
                    │                                         ▼
                    │                                需要修改值？
                    │                                         │
                    │                        ┌────────────────┴────────────────┐
                    │                        │                                 │
                    │                       否                                是
                    │                        │                                 │
                    │                只读访问                      检查可设置性
                    │                v.Interface()                 v.CanSet()
                    │                                                  │
                    │                                                  ▼
                    │                                         使用指针获取可设置值
                    │                                         reflect.ValueOf(&x).Elem()
                    │                                                  │
                    │                                         安全修改值
                    │                                         v.SetInt(100)
                    │
                    └────────────────┬────────────────┘
                                     │
                            需要动态创建？
                                     │
                    ┌────────────────┴────────────────┐
                    │                                 │
                   否                                是
                    │                                 │
            使用现有值                      动态创建值
            reflect.ValueOf(x)              reflect.New(type)
                                            reflect.MakeSlice(type, len, cap)
                                            reflect.MakeMap(type)
```

### 4. CGO 使用决策树

```
                    需要 CGO？
                         │
            ┌────────────┴────────────┐
            │                         │
           否                        是
            │                         │
        纯 Go 实现              调用 C 代码还是导出 Go？
                                     │
                    ┌────────────────┴────────────────┐
                    │                                 │
                调用 C                          导出 Go
                    │                                 │
            编写 C 注释                     使用 //export
            /*                              注释导出函数
            #include <header.h>
            */                                  │
            import "C"                          ▼
                    │                    检查函数签名
                    │                    只能使用 C 类型
                    │
                    ▼
            需要传递字符串？
                    │
    ┌───────────────┴───────────────┐
    │                               │
   否                              是
    │                               │
直接传递                      转换字符串
C.int, C.double               C.CString() / C.GoString()
                                    │
                            确保释放内存
                            defer C.free()
                                    │
                                    ▼
                            需要传递结构体？
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
                   否                              是
                    │                               │
            基本类型传递                  定义 C 结构体
            C.int, C.double               转换字段类型
                                            │
                                            ▼
                                    需要回调函数？
                                            │
                            ┌───────────────┴───────────────┐
                            │                               │
                           否                              是
                            │                               │
                    直接调用 C 函数                 定义回调类型
                    C.function()                    typedef int (*cb)(int)
                                                    //export GoCallback
                                                    注册回调函数
```

---

## 总结

### 特性对比表

| 特性 | 版本 | 复杂度 | 性能影响 | 主要用途 |
|------|------|--------|----------|----------|
| 泛型编程 | 1.18+ | 中 | 编译期展开 | 代码复用、类型安全 |
| 自引用约束 | 1.26+ | 高 | 无 | F-bounded 多态 |
| new() 表达式 | 1.26+ | 低 | 无 | 指针初始化简化 |
| 接口类型 | 1.0+ | 低 | 动态派发 | 多态、解耦 |
| 反射机制 | 1.0+ | 高 | 较大 | 动态类型处理 |
| 元编程 | 1.4+ | 中 | 编译期 | 代码生成 |
| CGO | 1.0+ | 高 | 较大 | C 语言互操作 |

### 最佳实践建议

1. **泛型编程**
   - 优先使用标准约束（`comparable`, `constraints.Ordered`）
   - 避免过度泛化，保持代码可读性
   - 利用自引用约束实现复杂递归类型

2. **接口设计**
   - 遵循"小接口"原则，接口应该小而专注
   - 使用接口组合构建复杂接口
   - 考虑使用类型集约束进行泛型编程

3. **反射使用**
   - 避免在热路径使用反射
   - 优先使用类型断言和类型开关
   - 注意内存安全和性能影响

4. **元编程**
   - 使用 `go generate` 自动化重复代码
   - 保持生成代码的可读性
   - 在生成代码中添加 `DO NOT EDIT` 注释

5. **CGO 使用**
   - 最小化 CGO 使用范围
   - 注意内存管理和释放
   - 避免在 C 代码中长时间阻塞

---

*文档生成时间：Go 1.26.1 高级特性分析*
*涵盖特性：泛型编程、接口类型、反射机制、元编程、new()表达式、自引用约束、CGO*
