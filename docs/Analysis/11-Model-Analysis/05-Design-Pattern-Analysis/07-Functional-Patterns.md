# 11.5.1 函数式设计模式分析

<!-- TOC START -->
- [11.5.1 函数式设计模式分析](#函数式设计模式分析)
  - [11.5.1.1 目录](#目录)
  - [11.5.1.2 1. 概述](#1-概述)
    - [11.5.1.2.1 函数式编程范式](#函数式编程范式)
    - [11.5.1.2.2 核心概念](#核心概念)
  - [11.5.1.3 2. 形式化定义](#2-形式化定义)
    - [11.5.1.3.1 函数式系统模型](#函数式系统模型)
    - [11.5.1.3.2 纯函数定义](#纯函数定义)
  - [11.5.1.4 3. 高阶函数模式](#3-高阶函数模式)
    - [11.5.1.4.1 Map模式](#map模式)
    - [11.5.1.4.2 Filter模式](#filter模式)
    - [11.5.1.4.3 Reduce模式](#reduce模式)
  - [11.5.1.5 4. 函数组合模式](#4-函数组合模式)
    - [11.5.1.5.1 函数组合](#函数组合)
    - [11.5.1.5.2 部分应用](#部分应用)
  - [11.5.1.6 5. 不可变数据结构模式](#5-不可变数据结构模式)
    - [11.5.1.6.1 不可变切片](#不可变切片)
    - [11.5.1.6.2 不可变映射](#不可变映射)
  - [11.5.1.7 6. 惰性求值模式](#6-惰性求值模式)
    - [11.5.1.7.1 惰性列表](#惰性列表)
    - [11.5.1.7.2 惰性求值器](#惰性求值器)
  - [11.5.1.8 7. 函子与单子模式](#7-函子与单子模式)
    - [11.5.1.8.1 函子模式](#函子模式)
    - [11.5.1.8.2 单子模式](#单子模式)
    - [11.5.1.8.3 applicative函子](#applicative函子)
  - [11.5.1.9 8. 模式匹配](#8-模式匹配)
    - [11.5.1.9.1 模式匹配器](#模式匹配器)
  - [11.5.1.10 9. 最佳实践](#9-最佳实践)
    - [11.5.1.10.1 设计原则](#设计原则)
    - [11.5.1.10.2 实现建议](#实现建议)
  - [11.5.1.11 10. 案例分析](#10-案例分析)
    - [11.5.1.11.1 数据处理管道](#数据处理管道)
    - [11.5.1.11.2 配置构建器](#配置构建器)
<!-- TOC END -->

## 11.5.1.1 目录

1. [概述](#1-概述)
2. [形式化定义](#2-形式化定义)
3. [高阶函数模式](#3-高阶函数模式)
4. [函数组合模式](#4-函数组合模式)
5. [不可变数据结构模式](#5-不可变数据结构模式)
6. [惰性求值模式](#6-惰性求值模式)
7. [函子与单子模式](#7-函子与单子模式)
8. [模式匹配](#8-模式匹配)
9. [最佳实践](#9-最佳实践)
10. [案例分析](#10-案例分析)

## 11.5.1.2 1. 概述

### 11.5.1.2.1 函数式编程范式

函数式编程是一种编程范式，强调使用纯函数、不可变数据和函数组合来构建程序。在Golang中，虽然语言本身不是纯函数式的，但我们可以应用函数式编程的思想和模式。

### 11.5.1.2.2 核心概念

- **纯函数**: 相同输入总是产生相同输出，无副作用
- **不可变性**: 数据一旦创建就不能修改
- **高阶函数**: 接受函数作为参数或返回函数的函数
- **函数组合**: 将多个函数组合成新函数
- **惰性求值**: 延迟计算直到需要结果

## 11.5.1.3 2. 形式化定义

### 11.5.1.3.1 函数式系统模型

**定义 2.1** (函数式系统): 一个函数式系统是一个五元组 $FS = (T, F, C, E, M)$，其中：

- $T = \{t_1, t_2, ..., t_n\}$ 是类型集合
- $F = \{f_1, f_2, ..., f_m\}$ 是函数集合
- $C: F \times F \rightarrow F$ 是函数组合操作
- $E: T \times F \rightarrow T$ 是函数求值操作
- $M: T \times T \rightarrow T$ 是类型映射操作

### 11.5.1.3.2 纯函数定义

**定义 2.2** (纯函数): 纯函数是一个三元组 $PF = (D, R, f)$，其中：

- $D$ 是定义域
- $R$ 是值域
- $f: D \rightarrow R$ 是映射函数，满足 $\forall x, y \in D: x = y \Rightarrow f(x) = f(y)$

## 11.5.1.4 3. 高阶函数模式

### 11.5.1.4.1 Map模式

**定义 3.1** (Map模式): Map模式是一个三元组 $Map = (T, U, f)$，其中：

- $T$ 是输入类型
- $U$ 是输出类型
- $f: T \rightarrow U$ 是转换函数

```go
// Map模式实现
package functional

import (
    "reflect"
)

// Map 高阶函数：对切片中的每个元素应用函数
func Map[T any, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, item := range slice {
        result[i] = fn(item)
    }
    return result
}

// MapWithIndex 带索引的Map函数
func MapWithIndex[T any, U any](slice []T, fn func(int, T) U) []U {
    result := make([]U, len(slice))
    for i, item := range slice {
        result[i] = fn(i, item)
    }
    return result
}

// MapToMap 将map中的值进行转换
func MapToMap[K comparable, V any, U any](m map[K]V, fn func(V) U) map[K]U {
    result := make(map[K]U)
    for k, v := range m {
        result[k] = fn(v)
    }
    return result
}

// MapKeys 将map中的键进行转换
func MapKeys[K comparable, V any, U comparable](m map[K]V, fn func(K) U) map[U]V {
    result := make(map[U]V)
    for k, v := range m {
        result[fn(k)] = v
    }
    return result
}

// 使用示例
func MapExample() {
    // 数字切片转换
    numbers := []int{1, 2, 3, 4, 5}
    doubled := Map(numbers, func(n int) int {
        return n * 2
    })
    fmt.Println("Doubled:", doubled) // [2, 4, 6, 8, 10]
    
    // 字符串转换
    strings := []string{"hello", "world", "golang"}
    lengths := Map(strings, func(s string) int {
        return len(s)
    })
    fmt.Println("Lengths:", lengths) // [5, 5, 6]
    
    // 带索引的转换
    indexed := MapWithIndex(strings, func(i int, s string) string {
        return fmt.Sprintf("%d: %s", i, s)
    })
    fmt.Println("Indexed:", indexed) // ["0: hello", "1: world", "2: golang"]
}

```

### 11.5.1.4.2 Filter模式

**定义 3.2** (Filter模式): Filter模式是一个三元组 $Filter = (T, P, f)$，其中：

- $T$ 是数据类型
- $P$ 是谓词类型
- $f: T \rightarrow bool$ 是过滤函数

```go
// Filter模式实现
package functional

// Filter 高阶函数：过滤切片中的元素
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, item := range slice {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}

// FilterWithIndex 带索引的Filter函数
func FilterWithIndex[T any](slice []T, predicate func(int, T) bool) []T {
    var result []T
    for i, item := range slice {
        if predicate(i, item) {
            result = append(result, item)
        }
    }
    return result
}

// FilterMap 同时进行过滤和转换
func FilterMap[T any, U any](slice []T, fn func(T) (U, bool)) []U {
    var result []U
    for _, item := range slice {
        if value, ok := fn(item); ok {
            result = append(result, value)
        }
    }
    return result
}

// 使用示例
func FilterExample() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // 过滤偶数
    evens := Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("Evens:", evens) // [2, 4, 6, 8, 10]
    
    // 过滤大于5的数
    large := Filter(numbers, func(n int) bool {
        return n > 5
    })
    fmt.Println("Large:", large) // [6, 7, 8, 9, 10]
    
    // 带索引的过滤
    indexed := FilterWithIndex(numbers, func(i int, n int) bool {
        return i%2 == 0 // 保留偶数索引
    })
    fmt.Println("Indexed:", indexed) // [1, 3, 5, 7, 9]
}

```

### 11.5.1.4.3 Reduce模式

**定义 3.3** (Reduce模式): Reduce模式是一个四元组 $Reduce = (T, U, f, init)$，其中：

- $T$ 是输入类型
- $U$ 是累积类型
- $f: U \times T \rightarrow U$ 是归约函数
- $init$ 是初始值

```go
// Reduce模式实现
package functional

// Reduce 高阶函数：将切片归约为单个值
func Reduce[T any, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, item := range slice {
        result = fn(result, item)
    }
    return result
}

// ReduceWithIndex 带索引的Reduce函数
func ReduceWithIndex[T any, U any](slice []T, initial U, fn func(U, int, T) U) U {
    result := initial
    for i, item := range slice {
        result = fn(result, i, item)
    }
    return result
}

// FoldLeft 左折叠
func FoldLeft[T any, U any](slice []T, initial U, fn func(U, T) U) U {
    return Reduce(slice, initial, fn)
}

// FoldRight 右折叠
func FoldRight[T any, U any](slice []T, initial U, fn func(T, U) U) U {
    result := initial
    for i := len(slice) - 1; i >= 0; i-- {
        result = fn(slice[i], result)
    }
    return result
}

// 使用示例
func ReduceExample() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // 求和
    sum := Reduce(numbers, 0, func(acc, n int) int {
        return acc + n
    })
    fmt.Println("Sum:", sum) // 15
    
    // 求积
    product := Reduce(numbers, 1, func(acc, n int) int {
        return acc * n
    })
    fmt.Println("Product:", product) // 120
    
    // 字符串连接
    strings := []string{"hello", " ", "world", "!"}
    concatenated := Reduce(strings, "", func(acc, s string) string {
        return acc + s
    })
    fmt.Println("Concatenated:", concatenated) // "hello world!"
    
    // 带索引的归约
    indexed := ReduceWithIndex(numbers, "", func(acc string, i int, n int) string {
        if acc == "" {
            return fmt.Sprintf("%d:%d", i, n)
        }
        return fmt.Sprintf("%s, %d:%d", acc, i, n)
    })
    fmt.Println("Indexed:", indexed) // "0:1, 1:2, 2:3, 3:4, 4:5"
}

```

## 11.5.1.5 4. 函数组合模式

### 11.5.1.5.1 函数组合

**定义 4.1** (函数组合): 函数组合是一个三元组 $Compose = (f, g, h)$，其中：

- $f: B \rightarrow C$ 是第一个函数
- $g: A \rightarrow B$ 是第二个函数
- $h: A \rightarrow C$ 是组合函数，满足 $h(x) = f(g(x))$

```go
// 函数组合模式实现
package functional

// Compose 函数组合：f ∘ g
func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C {
    return func(a A) C {
        return f(g(a))
    }
}

// ComposeMany 组合多个函数
func ComposeMany[A any](functions ...func(A) A) func(A) A {
    if len(functions) == 0 {
        return func(a A) A { return a }
    }
    
    return func(a A) A {
        result := a
        for _, fn := range functions {
            result = fn(result)
        }
        return result
    }
}

// Pipe 管道操作：从左到右执行函数
func Pipe[A any](a A, functions ...func(A) A) A {
    result := a
    for _, fn := range functions {
        result = fn(result)
    }
    return result
}

// Curry 柯里化：将多参数函数转换为单参数函数链
func Curry[A, B, C any](fn func(A, B) C) func(A) func(B) C {
    return func(a A) func(B) C {
        return func(b B) C {
            return fn(a, b)
        }
    }
}

// Uncurry 反柯里化：将单参数函数链转换为多参数函数
func Uncurry[A, B, C any](fn func(A) func(B) C) func(A, B) C {
    return func(a A, b B) C {
        return fn(a)(b)
    }
}

// 使用示例
func ComposeExample() {
    // 基本函数组合
    addOne := func(x int) int { return x + 1 }
    multiplyByTwo := func(x int) int { return x * 2 }
    square := func(x int) int { return x * x }
    
    // 组合函数
    composed := Compose(square, Compose(multiplyByTwo, addOne))
    result := composed(3) // ((3 + 1) * 2)^2 = 64
    fmt.Println("Composed result:", result)
    
    // 管道操作
    piped := Pipe(3, addOne, multiplyByTwo, square)
    fmt.Println("Piped result:", piped) // 64
    
    // 柯里化
    add := func(a, b int) int { return a + b }
    curriedAdd := Curry(add)
    addFive := curriedAdd(5)
    result = addFive(3) // 8
    fmt.Println("Curried result:", result)
    
    // 反柯里化
    uncurriedAdd := Uncurry(curriedAdd)
    result = uncurriedAdd(5, 3) // 8
    fmt.Println("Uncurried result:", result)
}

```

### 11.5.1.5.2 部分应用

**定义 4.2** (部分应用): 部分应用是一个三元组 $Partial = (f, args, g)$，其中：

- $f$ 是原函数
- $args$ 是部分参数
- $g$ 是部分应用的函数

```go
// 部分应用模式实现
package functional

// Partial 部分应用：固定部分参数
func Partial[A, B, C any](fn func(A, B) C, a A) func(B) C {
    return func(b B) C {
        return fn(a, b)
    }
}

// Partial2 两个参数的部分应用
func Partial2[A, B, C, D any](fn func(A, B, C) D, a A, b B) func(C) D {
    return func(c C) D {
        return fn(a, b, c)
    }
}

// Partial3 三个参数的部分应用
func Partial3[A, B, C, D, E any](fn func(A, B, C, D) E, a A, b B, c C) func(D) E {
    return func(d D) E {
        return fn(a, b, c, d)
    }
}

// 使用示例
func PartialExample() {
    // 基本部分应用
    add := func(a, b int) int { return a + b }
    addFive := Partial(add, 5)
    result := addFive(3) // 8
    fmt.Println("Partial result:", result)
    
    // 多参数部分应用
    format := func(prefix, name, suffix string) string {
        return prefix + name + suffix
    }
    greet := Partial2(format, "Hello, ", "!")
    message := greet("World") // "Hello, World!"
    fmt.Println("Greet:", message)
}

```

## 11.5.1.6 5. 不可变数据结构模式

### 11.5.1.6.1 不可变切片

**定义 5.1** (不可变切片): 不可变切片是一个三元组 $IS = (T, data, operations)$，其中：

- $T$ 是元素类型
- $data$ 是底层数据
- $operations$ 是操作集合，每个操作返回新的切片

```go
// 不可变数据结构模式实现
package functional

// ImmutableSlice 不可变切片
type ImmutableSlice[T any] struct {
    data []T
}

// NewImmutableSlice 创建新的不可变切片
func NewImmutableSlice[T any](data []T) *ImmutableSlice[T] {
    // 创建副本以避免外部修改
    copy := make([]T, len(data))
    copy(copy, data)
    return &ImmutableSlice[T]{data: copy}
}

// Get 获取元素（不修改原切片）
func (is *ImmutableSlice[T]) Get(index int) (T, bool) {
    if index < 0 || index >= len(is.data) {
        var zero T
        return zero, false
    }
    return is.data[index], true
}

// Set 设置元素（返回新切片）
func (is *ImmutableSlice[T]) Set(index int, value T) *ImmutableSlice[T] {
    if index < 0 || index >= len(is.data) {
        return is
    }
    
    newData := make([]T, len(is.data))
    copy(newData, is.data)
    newData[index] = value
    
    return &ImmutableSlice[T]{data: newData}
}

// Append 追加元素（返回新切片）
func (is *ImmutableSlice[T]) Append(values ...T) *ImmutableSlice[T] {
    newData := make([]T, len(is.data)+len(values))
    copy(newData, is.data)
    copy(newData[len(is.data):], values)
    
    return &ImmutableSlice[T]{data: newData}
}

// Slice 切片操作（返回新切片）
func (is *ImmutableSlice[T]) Slice(start, end int) *ImmutableSlice[T] {
    if start < 0 {
        start = 0
    }
    if end > len(is.data) {
        end = len(is.data)
    }
    if start >= end {
        return &ImmutableSlice[T]{data: []T{}}
    }
    
    newData := make([]T, end-start)
    copy(newData, is.data[start:end])
    
    return &ImmutableSlice[T]{data: newData}
}

// ToSlice 转换为普通切片
func (is *ImmutableSlice[T]) ToSlice() []T {
    result := make([]T, len(is.data))
    copy(result, is.data)
    return result
}

// Len 获取长度
func (is *ImmutableSlice[T]) Len() int {
    return len(is.data)
}

// 使用示例
func ImmutableSliceExample() {
    // 创建不可变切片
    original := []int{1, 2, 3, 4, 5}
    immutable := NewImmutableSlice(original)
    
    // 获取元素
    if value, ok := immutable.Get(2); ok {
        fmt.Println("Value at index 2:", value) // 3
    }
    
    // 设置元素（返回新切片）
    newSlice := immutable.Set(2, 10)
    fmt.Println("Original:", immutable.ToSlice()) // [1, 2, 3, 4, 5]
    fmt.Println("Modified:", newSlice.ToSlice())  // [1, 2, 10, 4, 5]
    
    // 追加元素
    appended := immutable.Append(6, 7, 8)
    fmt.Println("Appended:", appended.ToSlice()) // [1, 2, 3, 4, 5, 6, 7, 8]
    
    // 切片操作
    sliced := immutable.Slice(1, 4)
    fmt.Println("Sliced:", sliced.ToSlice()) // [2, 3, 4]
}

```

### 11.5.1.6.2 不可变映射

```go
// ImmutableMap 不可变映射
type ImmutableMap[K comparable, V any] struct {
    data map[K]V
}

// NewImmutableMap 创建新的不可变映射
func NewImmutableMap[K comparable, V any](data map[K]V) *ImmutableMap[K, V] {
    // 创建副本
    copy := make(map[K]V, len(data))
    for k, v := range data {
        copy[k] = v
    }
    return &ImmutableMap[K, V]{data: copy}
}

// Get 获取值
func (im *ImmutableMap[K, V]) Get(key K) (V, bool) {
    value, exists := im.data[key]
    return value, exists
}

// Set 设置值（返回新映射）
func (im *ImmutableMap[K, V]) Set(key K, value V) *ImmutableMap[K, V] {
    newData := make(map[K]V, len(im.data)+1)
    for k, v := range im.data {
        newData[k] = v
    }
    newData[key] = value
    
    return &ImmutableMap[K, V]{data: newData}
}

// Delete 删除键（返回新映射）
func (im *ImmutableMap[K, V]) Delete(key K) *ImmutableMap[K, V] {
    newData := make(map[K]V, len(im.data))
    for k, v := range im.data {
        if k != key {
            newData[k] = v
        }
    }
    
    return &ImmutableMap[K, V]{data: newData}
}

// Keys 获取所有键
func (im *ImmutableMap[K, V]) Keys() []K {
    keys := make([]K, 0, len(im.data))
    for k := range im.data {
        keys = append(keys, k)
    }
    return keys
}

// Values 获取所有值
func (im *ImmutableMap[K, V]) Values() []V {
    values := make([]V, 0, len(im.data))
    for _, v := range im.data {
        values = append(values, v)
    }
    return values
}

// ToMap 转换为普通映射
func (im *ImmutableMap[K, V]) ToMap() map[K]V {
    result := make(map[K]V, len(im.data))
    for k, v := range im.data {
        result[k] = v
    }
    return result
}

// Len 获取长度
func (im *ImmutableMap[K, V]) Len() int {
    return len(im.data)
}

```

## 11.5.1.7 6. 惰性求值模式

### 11.5.1.7.1 惰性列表

**定义 6.1** (惰性列表): 惰性列表是一个四元组 $LL = (T, head, tail, force)$，其中：

- $T$ 是元素类型
- $head$ 是头部元素
- $tail$ 是尾部列表（惰性）
- $force$ 是强制求值函数

```go
// 惰性求值模式实现
package functional

// LazyList 惰性列表
type LazyList[T any] struct {
    head T
    tail *LazyList[T]
    computed bool
    computeFn func() *LazyList[T]
}

// NewLazyList 创建新的惰性列表
func NewLazyList[T any](head T, computeFn func() *LazyList[T]) *LazyList[T] {
    return &LazyList[T]{
        head:      head,
        computeFn: computeFn,
        computed:  false,
    }
}

// EmptyLazyList 创建空惰性列表
func EmptyLazyList[T any]() *LazyList[T] {
    return nil
}

// IsEmpty 检查是否为空
func (ll *LazyList[T]) IsEmpty() bool {
    return ll == nil
}

// Head 获取头部元素
func (ll *LazyList[T]) Head() (T, bool) {
    if ll == nil {
        var zero T
        return zero, false
    }
    return ll.head, true
}

// Tail 获取尾部（惰性求值）
func (ll *LazyList[T]) Tail() *LazyList[T] {
    if ll == nil {
        return nil
    }
    
    if !ll.computed {
        ll.tail = ll.computeFn()
        ll.computed = true
    }
    
    return ll.tail
}

// Take 获取前n个元素
func (ll *LazyList[T]) Take(n int) *LazyList[T] {
    if n <= 0 || ll == nil {
        return nil
    }
    
    return NewLazyList(ll.head, func() *LazyList[T] {
        return ll.Tail().Take(n - 1)
    })
}

// Drop 跳过前n个元素
func (ll *LazyList[T]) Drop(n int) *LazyList[T] {
    if n <= 0 {
        return ll
    }
    if ll == nil {
        return nil
    }
    
    return ll.Tail().Drop(n - 1)
}

// ToSlice 转换为切片（强制求值）
func (ll *LazyList[T]) ToSlice() []T {
    if ll == nil {
        return []T{}
    }
    
    var result []T
    current := ll
    for current != nil {
        result = append(result, current.head)
        current = current.Tail()
    }
    
    return result
}

// 使用示例
func LazyListExample() {
    // 创建无限序列
    naturals := func() *LazyList[int] {
        var generate func(int) *LazyList[int]
        generate = func(n int) *LazyList[int] {
            return NewLazyList(n, func() *LazyList[int] {
                return generate(n + 1)
            })
        }
        return generate(1)
    }()
    
    // 获取前10个自然数
    firstTen := naturals.Take(10)
    fmt.Println("First 10 naturals:", firstTen.ToSlice()) // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    
    // 跳过前5个，再取5个
    skipped := naturals.Drop(5).Take(5)
    fmt.Println("Skip 5, take 5:", skipped.ToSlice()) // [6, 7, 8, 9, 10]
}

```

### 11.5.1.7.2 惰性求值器

```go
// Lazy 惰性求值器
type Lazy[T any] struct {
    value    T
    computed bool
    computeFn func() T
}

// NewLazy 创建新的惰性求值器
func NewLazy[T any](computeFn func() T) *Lazy[T] {
    return &Lazy[T]{
        computeFn: computeFn,
        computed:  false,
    }
}

// Force 强制求值
func (l *Lazy[T]) Force() T {
    if !l.computed {
        l.value = l.computeFn()
        l.computed = true
    }
    return l.value
}

// IsComputed 检查是否已求值
func (l *Lazy[T]) IsComputed() bool {
    return l.computed
}

// 使用示例
func LazyExample() {
    // 创建惰性计算
    expensive := NewLazy(func() int {
        fmt.Println("Computing expensive value...")
        time.Sleep(1 * time.Second)
        return 42
    })
    
    fmt.Println("Lazy created, not computed yet")
    
    // 第一次访问时计算
    value := expensive.Force()
    fmt.Println("Value:", value) // 42
    
    // 第二次访问使用缓存
    value2 := expensive.Force()
    fmt.Println("Value again:", value2) // 42 (no computation)
}

```

## 11.5.1.8 7. 函子与单子模式

### 11.5.1.8.1 函子模式

**定义 7.1** (函子): 函子是一个三元组 $Functor = (F, map, laws)$，其中：

- $F$ 是类型构造器
- $map: (A \rightarrow B) \rightarrow F[A] \rightarrow F[B]$ 是映射函数
- $laws$ 是函子定律

```go
// 函子模式实现
package functional

// Functor 函子接口
type Functor[T any] interface {
    Map[U any](fn func(T) U) Functor[U]
}

// Option 可选类型（函子）
type Option[T any] struct {
    value T
    hasValue bool
}

// Some 创建有值的Option
func Some[T any](value T) Option[T] {
    return Option[T]{
        value:    value,
        hasValue: true,
    }
}

// None 创建空Option
func None[T any]() Option[T] {
    return Option[T]{
        hasValue: false,
    }
}

// IsSome 检查是否有值
func (o Option[T]) IsSome() bool {
    return o.hasValue
}

// IsNone 检查是否为空
func (o Option[T]) IsNone() bool {
    return !o.hasValue
}

// Unwrap 获取值（不安全）
func (o Option[T]) Unwrap() T {
    if !o.hasValue {
        panic("Option is None")
    }
    return o.value
}

// UnwrapOr 获取值或默认值
func (o Option[T]) UnwrapOr(defaultValue T) T {
    if o.hasValue {
        return o.value
    }
    return defaultValue
}

// Map 函子映射
func (o Option[T]) Map[U any](fn func(T) U) Option[U] {
    if o.hasValue {
        return Some(fn(o.value))
    }
    return None[U]()
}

// AndThen 链式操作
func (o Option[T]) AndThen[U any](fn func(T) Option[U]) Option[U] {
    if o.hasValue {
        return fn(o.value)
    }
    return None[U]()
}

// 使用示例
func OptionExample() {
    // 基本操作
    someValue := Some(42)
    noneValue := None[int]()
    
    fmt.Println("Some is some:", someValue.IsSome()) // true
    fmt.Println("None is none:", noneValue.IsNone()) // true
    
    // 映射操作
    doubled := someValue.Map(func(x int) int {
        return x * 2
    })
    fmt.Println("Doubled:", doubled.Unwrap()) // 84
    
    // 链式操作
    result := someValue.
        Map(func(x int) int { return x + 10 }).
        AndThen(func(x int) Option[string] {
            if x > 50 {
                return Some(fmt.Sprintf("Large: %d", x))
            }
            return None[string]()
        })
    
    if result.IsSome() {
        fmt.Println("Result:", result.Unwrap()) // "Large: 52"
    }
}

```

### 11.5.1.8.2 单子模式

**定义 7.2** (单子): 单子是一个四元组 $Monad = (F, map, bind, unit)$，其中：

- $F$ 是类型构造器
- $map: (A \rightarrow B) \rightarrow F[A] \rightarrow F[B]$ 是映射函数
- $bind: F[A] \times (A \rightarrow F[B]) \rightarrow F[B]$ 是绑定函数
- $unit: A \rightarrow F[A]$ 是单位函数

```go
// 单子模式实现
package functional

// Monad 单子接口
type Monad[T any] interface {
    Functor[T]
    Bind[U any](fn func(T) Monad[U]) Monad[U]
}

// Result 结果类型（单子）
type Result[T any] struct {
    value T
    err   error
}

// Ok 创建成功结果
func Ok[T any](value T) Result[T] {
    return Result[T]{
        value: value,
        err:   nil,
    }
}

// Err 创建错误结果
func Err[T any](err error) Result[T] {
    return Result[T]{
        err: err,
    }
}

// IsOk 检查是否成功
func (r Result[T]) IsOk() bool {
    return r.err == nil
}

// IsErr 检查是否出错
func (r Result[T]) IsErr() bool {
    return r.err != nil
}

// Unwrap 获取值（不安全）
func (r Result[T]) Unwrap() T {
    if r.err != nil {
        panic(r.err)
    }
    return r.value
}

// UnwrapOr 获取值或默认值
func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.IsOk() {
        return r.value
    }
    return defaultValue
}

// Error 获取错误
func (r Result[T]) Error() error {
    return r.err
}

// Map 函子映射
func (r Result[T]) Map[U any](fn func(T) U) Result[U] {
    if r.IsOk() {
        return Ok(fn(r.value))
    }
    return Err[U](r.err)
}

// Bind 单子绑定
func (r Result[T]) Bind[U any](fn func(T) Result[U]) Result[U] {
    if r.IsOk() {
        return fn(r.value)
    }
    return Err[U](r.err)
}

// 使用示例
func ResultExample() {
    // 基本操作
    divide := func(a, b int) Result[int] {
        if b == 0 {
            return Err[int](errors.New("division by zero"))
        }
        return Ok(a / b)
    }
    
    // 成功案例
    result1 := divide(10, 2)
    if result1.IsOk() {
        fmt.Println("Result1:", result1.Unwrap()) // 5
    }
    
    // 失败案例
    result2 := divide(10, 0)
    if result2.IsErr() {
        fmt.Println("Error:", result2.Error()) // "division by zero"
    }
    
    // 链式操作
    result3 := divide(20, 4).
        Bind(func(n int) Result[int] {
            return divide(100, n)
        }).
        Map(func(n int) string {
            return fmt.Sprintf("Final result: %d", n)
        })
    
    if result3.IsOk() {
        fmt.Println(result3.Unwrap()) // "Final result: 20"
    }
}

```

### 11.5.1.8.3 applicative函子

**定义 7.3** (applicative函子): applicative函子是一个四元组 $Applicative = (F, map, pure, ap)$，其中：

- $F$ 是类型构造器
- $map: (A \rightarrow B) \rightarrow F[A] \rightarrow F[B]$ 是映射函数
- $pure: A \rightarrow F[A]$ 是纯函数
- $ap: F[A \rightarrow B] \times F[A] \rightarrow F[B]$ 是应用函数

```go
// Applicative函子实现
package functional

// Applicative applicative函子接口
type Applicative[T any] interface {
    Functor[T]
    Apply[U any](fn Applicative[func(T) U]) Applicative[U]
}

// ZipWith 合并两个Applicative
func ZipWith[A, B, C any](a Applicative[A], b Applicative[B], fn func(A, B) C) Applicative[C] {
    return a.Map(func(x A) func(B) C {
        return func(y B) C {
            return fn(x, y)
        }
    }).(Applicative[func(B) C]).Apply(b.(Applicative[B]))
}

// Validation 验证类型
type Validation[T any] struct {
    value  T
    errors []error
}

// Success 创建成功验证
func Success[T any](value T) Validation[T] {
    return Validation[T]{
        value:  value,
        errors: []error{},
    }
}

// Failure 创建失败验证
func Failure[T any](errs ...error) Validation[T] {
    return Validation[T]{
        errors: errs,
    }
}

// IsValid 检查是否有效
func (v Validation[T]) IsValid() bool {
    return len(v.errors) == 0
}

// Errors 获取所有错误
func (v Validation[T]) Errors() []error {
    return v.errors
}

// Map 函子映射
func (v Validation[T]) Map[U any](fn func(T) U) Validation[U] {
    if v.IsValid() {
        return Success(fn(v.value))
    }
    return Failure[U](v.errors...)
}

// Apply applicative应用
func (v Validation[T]) Apply[U any](fn Validation[func(T) U]) Validation[U] {
    if v.IsValid() && fn.IsValid() {
        return Success(fn.value(v.value))
    }
    
    errors := append([]error{}, v.errors...)
    errors = append(errors, fn.errors...)
    return Failure[U](errors...)
}

// 使用示例
func ValidationExample() {
    // 验证函数
    validateName := func(name string) Validation[string] {
        if len(name) == 0 {
            return Failure[string](errors.New("name cannot be empty"))
        }
        return Success(name)
    }
    
    validateAge := func(age int) Validation[int] {
        if age < 0 || age > 120 {
            return Failure[int](errors.New("age must be between 0 and 120"))
        }
        return Success(age)
    }
    
    // 组合验证
    createPerson := func(name string, age int) string {
        return fmt.Sprintf("%s (%d)", name, age)
    }
    
    // 有效数据
    valid := validateName("John").Map(func(name string) func(int) string {
        return func(age int) string {
            return createPerson(name, age)
        }
    }).(Validation[func(int) string]).Apply(validateAge(30))
    
    if valid.IsValid() {
        fmt.Println("Valid person:", valid.value) // "John (30)"
    }
    
    // 无效数据
    invalid := validateName("").Map(func(name string) func(int) string {
        return func(age int) string {
            return createPerson(name, age)
        }
    }).(Validation[func(int) string]).Apply(validateAge(150))
    
    if !invalid.IsValid() {
        for _, err := range invalid.Errors() {
            fmt.Println("Error:", err)
        }
    }
}

```

## 11.5.1.9 8. 模式匹配

### 11.5.1.9.1 模式匹配器

```go
// 模式匹配模式实现
package functional

// Pattern 模式接口
type Pattern[T, R any] interface {
    Matches(value T) bool
    Apply(value T) R
}

// PatternMatch 模式匹配器
type PatternMatch[T, R any] struct {
    patterns []Pattern[T, R]
    defaultFn func(T) R
}

// NewPatternMatch 创建新的模式匹配器
func NewPatternMatch[T, R any]() *PatternMatch[T, R] {
    return &PatternMatch[T, R]{
        patterns: make([]Pattern[T, R], 0),
    }
}

// AddPattern 添加模式
func (pm *PatternMatch[T, R]) AddPattern(pattern Pattern[T, R]) *PatternMatch[T, R] {
    pm.patterns = append(pm.patterns, pattern)
    return pm
}

// SetDefault 设置默认处理
func (pm *PatternMatch[T, R]) SetDefault(defaultFn func(T) R) *PatternMatch[T, R] {
    pm.defaultFn = defaultFn
    return pm
}

// Match 执行模式匹配
func (pm *PatternMatch[T, R]) Match(value T) R {
    for _, pattern := range pm.patterns {
        if pattern.Matches(value) {
            return pattern.Apply(value)
        }
    }
    
    if pm.defaultFn != nil {
        return pm.defaultFn(value)
    }
    
    panic("no matching pattern and no default handler")
}

// 具体模式实现
type EqualsPattern[T comparable, R any] struct {
    expected T
    handler  func(T) R
}

func NewEqualsPattern[T comparable, R any](expected T, handler func(T) R) *EqualsPattern[T, R] {
    return &EqualsPattern[T, R]{
        expected: expected,
        handler:  handler,
    }
}

func (ep *EqualsPattern[T, R]) Matches(value T) bool {
    return value == ep.expected
}

func (ep *EqualsPattern[T, R]) Apply(value T) R {
    return ep.handler(value)
}

type RangePattern[T any, R any] struct {
    predicate func(T) bool
 handler   func(T) R
}

func NewRangePattern[T any, R any](predicate func(T) bool, handler func(T) R) *RangePattern[T, R] {
    return &RangePattern[T, R]{
        predicate: predicate,
        handler:   handler,
    }
}

func (rp *RangePattern[T, R]) Matches(value T) bool {
    return rp.predicate(value)
}

func (rp *RangePattern[T, R]) Apply(value T) R {
    return rp.handler(value)
}

// 使用示例
func PatternMatchExample() {
    // 创建模式匹配器
    matcher := NewPatternMatch[int, string]().
        AddPattern(NewEqualsPattern(1, func(x int) string {
            return "one"
        })).
        AddPattern(NewEqualsPattern(2, func(x int) string {
            return "two"
        })).
        AddPattern(NewRangePattern(func(x int) bool {
            return x > 10
        }, func(x int) string {
            return fmt.Sprintf("large: %d", x)
        })).
        SetDefault(func(x int) string {
            return fmt.Sprintf("other: %d", x)
        })
    
    // 测试模式匹配
    fmt.Println(matcher.Match(1))   // "one"
    fmt.Println(matcher.Match(2))   // "two"
    fmt.Println(matcher.Match(15))  // "large: 15"
    fmt.Println(matcher.Match(5))   // "other: 5"
}

```

## 11.5.1.10 9. 最佳实践

### 11.5.1.10.1 设计原则

1. **纯函数优先**: 尽可能使用纯函数，避免副作用
2. **不可变性**: 使用不可变数据结构
3. **函数组合**: 通过组合小函数构建复杂功能
4. **高阶函数**: 使用高阶函数提高代码复用性
5. **惰性求值**: 在适当场景使用惰性求值提高性能

### 11.5.1.10.2 实现建议

1. **使用泛型**: 利用Golang泛型实现类型安全的函数式模式
2. **错误处理**: 使用Result类型进行函数式错误处理
3. **性能考虑**: 在性能关键路径避免过度使用函数式模式
4. **可读性**: 保持代码可读性，不要过度抽象
5. **测试**: 为函数式代码编写充分的测试

## 11.5.1.11 10. 案例分析

### 11.5.1.11.1 数据处理管道

```go
// 数据处理管道示例
package functional

// DataProcessor 数据处理管道
type DataProcessor struct {
    pipeline []func([]int) []int
}

// NewDataProcessor 创建新的数据处理器
func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        pipeline: make([]func([]int) []int, 0),
    }
}

// AddStep 添加处理步骤
func (dp *DataProcessor) AddStep(step func([]int) []int) *DataProcessor {
    dp.pipeline = append(dp.pipeline, step)
    return dp
}

// Process 执行处理管道
func (dp *DataProcessor) Process(data []int) []int {
    result := data
    for _, step := range dp.pipeline {
        result = step(result)
    }
    return result
}

// 使用示例
func DataProcessingExample() {
    // 创建处理管道
    processor := NewDataProcessor().
        AddStep(func(data []int) []int {
            // 过滤偶数
            return Filter(data, func(x int) bool {
                return x%2 == 0
            })
        }).
        AddStep(func(data []int) []int {
            // 平方
            return Map(data, func(x int) int {
                return x * x
            })
        }).
        AddStep(func(data []int) []int {
            // 过滤大于10的数
            return Filter(data, func(x int) bool {
                return x > 10
            })
        })
    
    // 处理数据
    input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    result := processor.Process(input)
    fmt.Println("Input:", input)
    fmt.Println("Output:", result) // [16, 36, 64, 100]
}

```

### 11.5.1.11.2 配置构建器

```go
// 配置构建器示例
package functional

// Config 配置结构
type Config struct {
    Host     string
    Port     int
    Username string
    Password string
    SSL      bool
}

// ConfigBuilder 配置构建器
type ConfigBuilder struct {
    config Config
}

// NewConfigBuilder 创建新的配置构建器
func NewConfigBuilder() *ConfigBuilder {
    return &ConfigBuilder{
        config: Config{},
    }
}

// WithHost 设置主机
func (cb *ConfigBuilder) WithHost(host string) *ConfigBuilder {
    cb.config.Host = host
    return cb
}

// WithPort 设置端口
func (cb *ConfigBuilder) WithPort(port int) *ConfigBuilder {
    cb.config.Port = port
    return cb
}

// WithCredentials 设置凭据
func (cb *ConfigBuilder) WithCredentials(username, password string) *ConfigBuilder {
    cb.config.Username = username
    cb.config.Password = password
    return cb
}

// WithSSL 设置SSL
func (cb *ConfigBuilder) WithSSL(ssl bool) *ConfigBuilder {
    cb.config.SSL = ssl
    return cb
}

// Build 构建配置
func (cb *ConfigBuilder) Build() Config {
    return cb.config
}

// 使用示例
func ConfigBuilderExample() {
    // 构建配置
    config := NewConfigBuilder().
        WithHost("localhost").
        WithPort(8080).
        WithCredentials("admin", "password").
        WithSSL(true).
        Build()
    
    fmt.Printf("Config: %+v\n", config)
}

```

---

**总结**: 本文档提供了函数式设计模式的完整分析，包括形式化定义、Golang实现和最佳实践。这些模式为构建函数式风格的Golang程序提供了重要的理论基础和实践指导，支持各种业务场景的需求。
