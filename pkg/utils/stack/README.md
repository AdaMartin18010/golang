# 栈工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [栈工具](#栈工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

栈工具提供了多种栈实现，包括简单栈、最大栈、最小栈等，帮助开发者处理各种栈场景。

---

## 2. 功能特性

### 2.1 简单栈

- `SimpleStack`: 简单栈实现
- `NewSimpleStack`: 创建简单栈
- `Push`: 入栈
- `Pop`: 出栈
- `Peek`: 查看栈顶元素
- `Size`: 获取栈大小
- `IsEmpty`: 检查栈是否为空
- `Clear`: 清空栈
- `ToSlice`: 转换为切片

### 2.2 最大栈

- `MaxStack`: 最大栈实现（支持O(1)获取最大值）
- `NewMaxStack`: 创建最大栈
- `Push`: 入栈
- `Pop`: 出栈
- `Peek`: 查看栈顶元素
- `Max`: 获取最大值
- `Size`: 获取栈大小
- `IsEmpty`: 检查栈是否为空
- `Clear`: 清空栈

### 2.3 最小栈

- `MinStack`: 最小栈实现（支持O(1)获取最小值）
- `NewMinStack`: 创建最小栈
- `Push`: 入栈
- `Pop`: 出栈
- `Peek`: 查看栈顶元素
- `Min`: 获取最小值
- `Size`: 获取栈大小
- `IsEmpty`: 检查栈是否为空
- `Clear`: 清空栈

---

## 3. 使用示例

### 3.1 简单栈

```go
import "github.com/yourusername/golang/pkg/utils/stack"

// 创建栈
s := stack.NewSimpleStack[string]()

// 入栈
s.Push("first")
s.Push("second")
s.Push("third")

// 查看栈顶
item, ok := s.Peek()
if ok {
    fmt.Printf("Top: %s\n", item)  // "third"
}

// 出栈
item, ok = s.Pop()
if ok {
    fmt.Printf("Popped: %s\n", item)  // "third"
}

// 获取大小
size := s.Size()
fmt.Printf("Size: %d\n", size)

// 转换为切片
items := s.ToSlice()
fmt.Printf("Items: %v\n", items)

// 清空栈
s.Clear()
```

### 3.2 最大栈

```go
// 创建最大栈
ms := stack.NewMaxStack[int](func(a, b int) bool {
    return a > b
})

// 入栈
ms.Push(3)
ms.Push(1)
ms.Push(5)
ms.Push(2)

// 获取最大值
max, ok := ms.Max()  // 5

// 出栈
item, ok := ms.Pop()  // 2
max, ok = ms.Max()    // 5

item, ok = ms.Pop()   // 5
max, ok = ms.Max()    // 3
```

### 3.3 最小栈

```go
// 创建最小栈
ms := stack.NewMinStack[int](func(a, b int) bool {
    return a < b
})

// 入栈
ms.Push(3)
ms.Push(1)
ms.Push(5)
ms.Push(2)

// 获取最小值
min, ok := ms.Min()  // 1

// 出栈
item, ok := ms.Pop()  // 2
min, ok = ms.Min()    // 1

item, ok = ms.Pop()   // 5
min, ok = ms.Min()    // 1
```

### 3.4 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/stack"
)

func main() {
    // 简单栈
    s := stack.NewSimpleStack[int]()
    s.Push(1)
    s.Push(2)
    s.Push(3)
    
    for !s.IsEmpty() {
        item, _ := s.Pop()
        fmt.Printf("Popped: %d\n", item)
    }
    
    // 最大栈
    ms := stack.NewMaxStack[int](func(a, b int) bool {
        return a > b
    })
    ms.Push(3)
    ms.Push(1)
    ms.Push(5)
    
    max, _ := ms.Max()
    fmt.Printf("Max: %d\n", max)
}
```

---

**更新日期**: 2025-11-11

