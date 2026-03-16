# 计数器工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [计数器工具](#计数器工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

计数器工具提供了多种计数器实现，包括简单计数器、最大计数器、最小计数器、速率计数器、滑动窗口计数器、多键计数器等，帮助开发者进行各种计数统计。

---

## 2. 功能特性

### 2.1 简单计数器

- `SimpleCounter`: 简单计数器实现
- `NewSimpleCounter`: 创建简单计数器
- `Increment`: 增加1
- `Decrement`: 减少1
- `Add`: 增加指定值
- `Get`: 获取当前值
- `Set`: 设置值
- `Reset`: 重置计数器

### 2.2 最大计数器

- `MaxCounter`: 最大计数器实现（只增不减）
- `NewMaxCounter`: 创建最大计数器
- `Increment`: 增加1
- `Add`: 增加指定值
- `Get`: 获取当前值
- `Reset`: 重置计数器

### 2.3 最小计数器

- `MinCounter`: 最小计数器实现（只减不增）
- `NewMinCounter`: 创建最小计数器
- `Decrement`: 减少1
- `Subtract`: 减少指定值
- `Get`: 获取当前值
- `Reset`: 重置计数器

### 2.4 速率计数器

- `RateCounter`: 速率计数器实现
- `NewRateCounter`: 创建速率计数器
- `Increment`: 增加1
- `Add`: 增加指定值
- `Get`: 获取当前速率（每秒）
- `Reset`: 重置计数器

### 2.5 滑动窗口计数器

- `SlidingWindowCounter`: 滑动窗口计数器实现
- `NewSlidingWindowCounter`: 创建滑动窗口计数器
- `Increment`: 增加1
- `Add`: 增加指定值
- `Get`: 获取窗口内的总数
- `Reset`: 重置计数器

### 2.6 多键计数器

- `MultiCounter`: 多键计数器实现
- `NewMultiCounter`: 创建多键计数器
- `Increment`: 增加指定键的计数
- `Decrement`: 减少指定键的计数
- `Add`: 增加指定键的计数
- `Get`: 获取指定键的计数
- `GetAll`: 获取所有计数
- `Reset`: 重置指定键的计数
- `ResetAll`: 重置所有计数
- `Keys`: 获取所有键

---

## 3. 使用示例

### 3.1 简单计数器

```go
import "github.com/yourusername/golang/pkg/utils/counter"

// 创建简单计数器
c := counter.NewSimpleCounter()

// 增加
c.Increment()
c.Increment()

// 减少
c.Decrement()

// 增加指定值
c.Add(5)

// 获取值
value := c.Get()

// 设置值
c.Set(10)

// 重置
c.Reset()
```

### 3.2 最大计数器

```go
// 创建最大计数器
c := counter.NewMaxCounter()

// 只能增加
c.Increment()
c.Add(5)

// 尝试减少（无效）
c.Add(-1)  // 不会减少
```

### 3.3 最小计数器

```go
// 创建最小计数器（初始值10）
c := counter.NewMinCounter(10)

// 只能减少
c.Decrement()
c.Subtract(5)

// 重置
c.Reset(10)
```

### 3.4 速率计数器

```go
// 创建速率计数器（窗口1分钟，间隔1秒）
rc := counter.NewRateCounter(1*time.Minute, 1*time.Second)

// 增加计数
rc.Increment()
rc.Add(5)

// 获取速率（每秒）
rate := rc.Get()
fmt.Printf("Rate: %.2f/s\n", rate)
```

### 3.5 滑动窗口计数器

```go
// 创建滑动窗口计数器（窗口1分钟，间隔1秒）
swc := counter.NewSlidingWindowCounter(1*time.Minute, 1*time.Second)

// 增加计数
swc.Increment()
swc.Add(5)

// 获取窗口内的总数
total := swc.Get()
fmt.Printf("Total: %d\n", total)
```

### 3.6 多键计数器

```go
// 创建多键计数器
mc := counter.NewMultiCounter()

// 增加不同键的计数
mc.Increment("key1")
mc.Increment("key1")
mc.Increment("key2")

// 获取指定键的计数
count1 := mc.Get("key1")  // 2
count2 := mc.Get("key2")  // 1

// 获取所有计数
all := mc.GetAll()
fmt.Printf("All counts: %v\n", all)

// 获取所有键
keys := mc.Keys()

// 重置指定键
mc.Reset("key1")

// 重置所有
mc.ResetAll()
```

### 3.7 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/golang/pkg/utils/counter"
)

func main() {
    // 简单计数器
    c := counter.NewSimpleCounter()
    c.Increment()
    c.Increment()
    fmt.Printf("Count: %d\n", c.Get())
    
    // 速率计数器
    rc := counter.NewRateCounter(1*time.Minute, 1*time.Second)
    for i := 0; i < 10; i++ {
        rc.Increment()
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Printf("Rate: %.2f/s\n", rc.Get())
    
    // 多键计数器
    mc := counter.NewMultiCounter()
    mc.Increment("user1")
    mc.Increment("user1")
    mc.Increment("user2")
    fmt.Printf("User1: %d, User2: %d\n", mc.Get("user1"), mc.Get("user2"))
}
```

---

**更新日期**: 2025-11-11

