# 队列工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [队列工具](#队列工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

队列工具提供了多种队列实现，包括简单队列、优先队列、循环队列等，帮助开发者处理各种队列场景。

---

## 2. 功能特性

### 2.1 简单队列

- `SimpleQueue`: 简单队列实现
- `NewSimpleQueue`: 创建简单队列
- `Enqueue`: 入队
- `Dequeue`: 出队
- `Peek`: 查看队首元素
- `Size`: 获取队列大小
- `IsEmpty`: 检查队列是否为空
- `Clear`: 清空队列
- `ToSlice`: 转换为切片

### 2.2 优先队列

- `PriorityQueue`: 优先队列实现
- `NewPriorityQueue`: 创建优先队列
- `Enqueue`: 入队（带优先级）
- `Dequeue`: 出队
- `Peek`: 查看队首元素
- `Size`: 获取队列大小
- `IsEmpty`: 检查队列是否为空
- `Clear`: 清空队列

### 2.3 循环队列

- `CircularQueue`: 循环队列实现
- `NewCircularQueue`: 创建循环队列
- `Enqueue`: 入队
- `Dequeue`: 出队
- `Peek`: 查看队首元素
- `Size`: 获取队列大小
- `IsEmpty`: 检查队列是否为空
- `IsFull`: 检查队列是否已满
- `Clear`: 清空队列
- `Capacity`: 获取队列容量

---

## 3. 使用示例

### 3.1 简单队列

```go
import "github.com/yourusername/golang/pkg/utils/queue"

// 创建队列
q := queue.NewSimpleQueue[string]()

// 入队
q.Enqueue("first")
q.Enqueue("second")
q.Enqueue("third")

// 查看队首
item, ok := q.Peek()
if ok {
    fmt.Printf("Front: %s\n", item)
}

// 出队
item, ok = q.Dequeue()
if ok {
    fmt.Printf("Dequeued: %s\n", item)
}

// 获取大小
size := q.Size()
fmt.Printf("Size: %d\n", size)

// 转换为切片
items := q.ToSlice()
fmt.Printf("Items: %v\n", items)

// 清空队列
q.Clear()
```

### 3.2 优先队列

```go
// 创建优先队列
pq := queue.NewPriorityQueue[string]()

// 入队（带优先级）
pq.Enqueue("low priority", 1)
pq.Enqueue("high priority", 10)
pq.Enqueue("medium priority", 5)

// 出队（按优先级）
item, ok := pq.Dequeue()  // "high priority"
item, ok = pq.Dequeue()   // "medium priority"
item, ok = pq.Dequeue()   // "low priority"
```

### 3.3 循环队列

```go
// 创建循环队列（容量为3）
cq := queue.NewCircularQueue[int](3)

// 入队
cq.Enqueue(1)
cq.Enqueue(2)
cq.Enqueue(3)

// 检查是否已满
if cq.IsFull() {
    fmt.Println("Queue is full")
}

// 出队
item, ok := cq.Dequeue()  // 1

// 可以继续入队
cq.Enqueue(4)

// 获取容量
capacity := cq.Capacity()
fmt.Printf("Capacity: %d\n", capacity)
```

### 3.4 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/queue"
)

func main() {
    // 简单队列
    q := queue.NewSimpleQueue[int]()
    q.Enqueue(1)
    q.Enqueue(2)
    q.Enqueue(3)
    
    for !q.IsEmpty() {
        item, _ := q.Dequeue()
        fmt.Printf("Dequeued: %d\n", item)
    }
    
    // 优先队列
    pq := queue.NewPriorityQueue[string]()
    pq.Enqueue("task1", 1)
    pq.Enqueue("task2", 10)
    pq.Enqueue("task3", 5)
    
    for !pq.IsEmpty() {
        item, _ := pq.Dequeue()
        fmt.Printf("Processed: %s\n", item)
    }
}
```

---

**更新日期**: 2025-11-11

