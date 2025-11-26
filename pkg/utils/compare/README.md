# 比较工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [比较工具](#比较工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

比较工具提供了各种比较功能，包括基本类型比较、值比较、范围检查、零值检查、nil检查、切片比较、映射比较等。

---

## 2. 功能特性

### 2.1 基本比较

- `Equal`: 检查两个值是否相等
- `NotEqual`: 检查两个值是否不相等
- `CompareInt`: 比较两个整数
- `CompareInt64`: 比较两个int64
- `CompareFloat64`: 比较两个float64
- `CompareString`: 比较两个字符串
- `CompareTime`: 比较两个时间

### 2.2 大小比较

- `Less`: 检查a是否小于b
- `Greater`: 检查a是否大于b
- `LessOrEqual`: 检查a是否小于等于b
- `GreaterOrEqual`: 检查a是否大于等于b

### 2.3 最值

- `Min`: 返回两个值中的较小值
- `Max`: 返回两个值中的较大值
- `MinInt`: 返回两个整数中的较小值
- `MaxInt`: 返回两个整数中的较大值
- `MinInt64`: 返回两个int64中的较小值
- `MaxInt64`: 返回两个int64中的较大值
- `MinFloat64`: 返回两个float64中的较小值
- `MaxFloat64`: 返回两个float64中的较大值
- `MinString`: 返回两个字符串中的较小值（字典序）
- `MaxString`: 返回两个字符串中的较大值（字典序）
- `MinTime`: 返回两个时间中的较早时间
- `MaxTime`: 返回两个时间中的较晚时间

### 2.4 范围检查

- `InRange`: 检查值是否在范围内
- `InRangeInt`: 检查整数是否在范围内
- `InRangeInt64`: 检查int64是否在范围内
- `InRangeFloat64`: 检查float64是否在范围内

### 2.5 范围限制

- `Clamp`: 将值限制在[min, max]范围内
- `ClampInt`: 将整数限制在[min, max]范围内
- `ClampInt64`: 将int64限制在[min, max]范围内
- `ClampFloat64`: 将float64限制在[min, max]范围内

### 2.6 零值检查

- `IsZero`: 检查值是否为零值
- `IsNil`: 检查值是否为nil
- `IsEmpty`: 检查值是否为空（nil、零值或空集合）

### 2.7 函数式比较

- `Compare`: 使用比较函数比较两个值
- `LessThan`: 使用小于函数比较两个值
- `EqualTo`: 使用相等函数比较两个值
- `CompareBy`: 根据键函数比较两个值
- `LessBy`: 根据键函数检查a是否小于b
- `EqualBy`: 根据键函数检查a是否等于b

### 2.8 集合比较

- `CompareSlice`: 比较两个切片
- `CompareSliceFunc`: 使用比较函数比较两个切片
- `CompareMap`: 比较两个映射
- `CompareMapFunc`: 使用比较函数比较两个映射

---

## 3. 使用示例

### 3.1 基本比较

```go
import "github.com/yourusername/golang/pkg/utils/compare"

// 相等检查
equal := compare.Equal(1, 1)  // true
notEqual := compare.NotEqual(1, 2)  // true

// 整数比较
result := compare.CompareInt(1, 2)  // -1 (小于)
result = compare.CompareInt(2, 1)   // 1 (大于)
result = compare.CompareInt(1, 1)   // 0 (相等)

// 浮点数比较
result = compare.CompareFloat64(1.0, 2.0)  // -1

// 字符串比较
result = compare.CompareString("a", "b")  // -1

// 时间比较
t1 := time.Now()
t2 := t1.Add(time.Hour)
result = compare.CompareTime(t1, t2)  // -1
```

### 3.2 大小比较

```go
// 大小比较
less := compare.Less(1, 2)  // true
greater := compare.Greater(2, 1)  // true
lessOrEqual := compare.LessOrEqual(1, 2)  // true
greaterOrEqual := compare.GreaterOrEqual(2, 1)  // true
```

### 3.3 最值

```go
// 最值
min := compare.MinInt(1, 2)  // 1
max := compare.MaxInt(1, 2)  // 2

minFloat := compare.MinFloat64(1.0, 2.0)  // 1.0
maxFloat := compare.MaxFloat64(1.0, 2.0)  // 2.0

minStr := compare.MinString("a", "b")  // "a"
maxStr := compare.MaxString("a", "b")  // "b"

t1 := time.Now()
t2 := t1.Add(time.Hour)
minTime := compare.MinTime(t1, t2)  // t1
maxTime := compare.MaxTime(t1, t2)  // t2
```

### 3.4 范围检查

```go
// 范围检查
inRange := compare.InRangeInt(5, 1, 10)  // true
inRange = compare.InRangeInt(15, 1, 10)  // false

inRangeFloat := compare.InRangeFloat64(5.5, 1.0, 10.0)  // true
```

### 3.5 范围限制

```go
// 范围限制
clamped := compare.ClampInt(15, 1, 10)  // 10
clamped = compare.ClampInt(-5, 1, 10)   // 1
clamped = compare.ClampInt(5, 1, 10)    // 5

clampedFloat := compare.ClampFloat64(15.5, 1.0, 10.0)  // 10.0
```

### 3.6 零值检查

```go
// 零值检查
isZero := compare.IsZero(0)  // true
isZero = compare.IsZero("")  // true
isZero = compare.IsZero(1)   // false

// nil检查
var s []int
isNil := compare.IsNil(s)  // true
isNil = compare.IsNil(1)   // false

// 空值检查
isEmpty := compare.IsEmpty(0)  // true
isEmpty = compare.IsEmpty("")  // true
isEmpty = compare.IsEmpty(nil) // true
```

### 3.7 函数式比较

```go
type Person struct {
    Name string
    Age  int
}

// 使用比较函数
people := []Person{
    {"Alice", 30},
    {"Bob", 25},
}

result := compare.Compare(people[0], people[1], func(a, b Person) int {
    return compare.CompareInt(a.Age, b.Age)
})

// 根据键函数比较
less := compare.LessBy(people[0], people[1], func(p Person) int {
    return p.Age
}, func(a, b int) bool {
    return a < b
})

equal := compare.EqualBy(people[0], people[1], func(p Person) string {
    return p.Name
})
```

### 3.8 集合比较

```go
// 切片比较
a := []int{1, 2, 3}
b := []int{1, 2, 4}
result := compare.CompareSlice(a, b)  // -1

// 使用比较函数比较切片
result = compare.CompareSliceFunc(a, b, compare.CompareInt)

// 映射比较
m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"a": 1, "b": 2}
equal := compare.CompareMap(m1, m2)  // true

// 使用比较函数比较映射
equal = compare.CompareMapFunc(m1, m2, func(a, b int) bool {
    return a == b
})
```

### 3.9 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/compare"
)

func main() {
    // 基本比较
    fmt.Printf("Equal: %v\n", compare.Equal(1, 1))
    fmt.Printf("Compare: %d\n", compare.CompareInt(1, 2))

    // 最值
    fmt.Printf("Min: %d\n", compare.MinInt(1, 2))
    fmt.Printf("Max: %d\n", compare.MaxInt(1, 2))

    // 范围检查
    fmt.Printf("InRange: %v\n", compare.InRangeInt(5, 1, 10))

    // 范围限制
    fmt.Printf("Clamp: %d\n", compare.ClampInt(15, 1, 10))

    // 零值检查
    fmt.Printf("IsZero: %v\n", compare.IsZero(0))
    fmt.Printf("IsNil: %v\n", compare.IsNil(nil))
    fmt.Printf("IsEmpty: %v\n", compare.IsEmpty(""))

    // 切片比较
    a := []int{1, 2, 3}
    b := []int{1, 2, 4}
    fmt.Printf("CompareSlice: %d\n", compare.CompareSlice(a, b))
}
```

---

**更新日期**: 2025-11-11
