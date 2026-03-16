# 排序工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [排序工具](#排序工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

排序工具提供了各种排序功能，包括基本类型排序、自定义排序、稳定排序、反向排序、去重排序、TopN/BottomN、随机打乱、多字段排序等。

---

## 2. 功能特性

### 2.1 基本类型排序

- `Ints`: 对整数切片进行排序
- `IntsAreSorted`: 检查整数切片是否已排序
- `SearchInts`: 在已排序的整数切片中搜索
- `Float64s`: 对float64切片进行排序
- `Float64sAreSorted`: 检查float64切片是否已排序
- `SearchFloat64s`: 在已排序的float64切片中搜索
- `Strings`: 对字符串切片进行排序
- `StringsAreSorted`: 检查字符串切片是否已排序
- `SearchStrings`: 在已排序的字符串切片中搜索

### 2.2 反向排序

- `IntsReverse`: 对整数切片进行反向排序
- `Float64sReverse`: 对float64切片进行反向排序
- `StringsReverse`: 对字符串切片进行反向排序
- `Reverse`: 反转切片

### 2.3 自定义排序

- `SortBy`: 根据比较函数对切片进行排序
- `SortByFunc`: 根据比较函数对切片进行排序（使用元素比较）
- `SortStable`: 稳定排序
- `SortStableByFunc`: 稳定排序（使用元素比较）
- `SortByKey`: 根据键函数对切片进行排序
- `SortByKeyInt`: 根据整数键对切片进行排序
- `SortByKeyString`: 根据字符串键对切片进行排序
- `SortByKeyFloat64`: 根据float64键对切片进行排序

### 2.4 排序检查

- `IsSorted`: 检查切片是否已排序
- `IsSortedFunc`: 检查切片是否已排序（使用元素比较）

### 2.5 搜索

- `Search`: 在已排序的切片中搜索
- `SearchSlice`: 在已排序的切片中搜索元素

### 2.6 去重排序

- `Unique`: 去重并排序
- `UniqueInts`: 去重并排序整数切片
- `UniqueFloat64s`: 去重并排序float64切片
- `UniqueStrings`: 去重并排序字符串切片

### 2.7 TopN/BottomN

- `TopN`: 返回前N个最大元素
- `BottomN`: 返回前N个最小元素
- `TopNInts`: 返回前N个最大整数
- `BottomNInts`: 返回前N个最小整数
- `TopNFloat64s`: 返回前N个最大float64
- `BottomNFloat64s`: 返回前N个最小float64

### 2.8 随机打乱

- `Shuffle`: 随机打乱切片
- `ShuffleWithSeed`: 使用种子随机打乱切片

### 2.9 多字段排序

- `MultiSort`: 多字段排序

### 2.10 比较函数

- `CompareInt`: 比较两个整数
- `CompareFloat64`: 比较两个float64
- `CompareString`: 比较两个字符串

---

## 3. 使用示例

### 3.1 基本类型排序

```go
import "github.com/yourusername/golang/pkg/utils/sort"

// 整数排序
nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
sort.Ints(nums)  // [1, 1, 2, 3, 4, 5, 6, 9]

// 检查是否已排序
isSorted := sort.IntsAreSorted(nums)  // true

// 搜索
index := sort.SearchInts(nums, 5)  // 5

// float64排序
floats := []float64{3.1, 1.4, 4.1, 1.5}
sort.Float64s(floats)

// 字符串排序
strs := []string{"banana", "apple", "cherry"}
sort.Strings(strs)  // ["apple", "banana", "cherry"]
```

### 3.2 反向排序

```go
// 反向排序
nums := []int{1, 2, 3, 4, 5}
sort.IntsReverse(nums)  // [5, 4, 3, 2, 1]

// 反转切片
sort.Reverse(nums)  // [1, 2, 3, 4, 5]
```

### 3.3 自定义排序

```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
}

// 根据年龄排序
sort.SortByFunc(people, func(a, b Person) bool {
    return a.Age < b.Age
})

// 根据键排序
sort.SortByKeyInt(people, func(p Person) int {
    return p.Age
})
```

### 3.4 去重排序

```go
// 去重并排序
nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
unique := sort.UniqueInts(nums)  // [1, 2, 3, 4, 5, 6, 9]

// 自定义类型去重
type Item struct {
    ID   int
    Name string
}
items := []Item{
    {1, "A"},
    {2, "B"},
    {1, "A"},
}
uniqueItems := sort.Unique(items, func(a, b Item) bool {
    return a.ID < b.ID
})
```

### 3.5 TopN/BottomN

```go
// 获取前N个最大元素
nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
top3 := sort.TopNInts(nums, 3)  // [9, 6, 5]

// 获取前N个最小元素
bottom3 := sort.BottomNInts(nums, 3)  // [1, 1, 2]

// 自定义类型TopN
type Score struct {
    Name  string
    Value int
}
scores := []Score{
    {"Alice", 90},
    {"Bob", 85},
    {"Charlie", 95},
}
top2 := sort.TopN(scores, 2, func(a, b Score) bool {
    return a.Value > b.Value
})
```

### 3.6 随机打乱

```go
// 随机打乱
nums := []int{1, 2, 3, 4, 5}
sort.Shuffle(nums)

// 使用种子随机打乱
sort.ShuffleWithSeed(nums, 12345)
```

### 3.7 多字段排序

```go
type Person struct {
    Name string
    Age  int
    City string
}

people := []Person{
    {"Alice", 30, "Beijing"},
    {"Bob", 30, "Shanghai"},
    {"Charlie", 25, "Beijing"},
}

// 先按年龄排序，再按城市排序
sort.MultiSort(people,
    func(a, b Person) int {
        return sort.CompareInt(a.Age, b.Age)
    },
    func(a, b Person) int {
        return sort.CompareString(a.City, b.City)
    },
)
```

### 3.8 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/sort"
)

func main() {
    // 基本排序
    nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
    sort.Ints(nums)
    fmt.Printf("Sorted: %v\n", nums)

    // 反向排序
    sort.IntsReverse(nums)
    fmt.Printf("Reversed: %v\n", nums)

    // 去重排序
    unique := sort.UniqueInts([]int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3})
    fmt.Printf("Unique: %v\n", unique)

    // TopN
    top3 := sort.TopNInts([]int{3, 1, 4, 1, 5, 9, 2, 6}, 3)
    fmt.Printf("Top 3: %v\n", top3)

    // 自定义排序
    type Person struct {
        Name string
        Age  int
    }
    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
    }
    sort.SortByKeyInt(people, func(p Person) int {
        return p.Age
    })
    fmt.Printf("Sorted by age: %v\n", people)
}
```

---

**更新日期**: 2025-11-11
