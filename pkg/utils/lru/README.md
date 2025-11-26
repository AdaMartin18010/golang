# LRU缓存工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [LRU缓存工具](#lru缓存工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

LRU缓存工具提供了LRU（Least Recently Used）缓存实现，用于缓存最近使用的数据，自动淘汰最久未使用的数据。

---

## 2. 功能特性

### 2.1 LRU缓存

- `LRUCache`: LRU缓存实现
- `NewLRUCache`: 创建LRU缓存
- `Get`: 获取值（会更新访问顺序）
- `Put`: 设置值
- `Delete`: 删除键
- `Contains`: 检查键是否存在
- `Size`: 获取缓存大小
- `Capacity`: 获取缓存容量
- `Clear`: 清空缓存
- `Keys`: 获取所有键
- `Values`: 获取所有值
- `Peek`: 查看值（不更新访问顺序）
- `GetOldest`: 获取最旧的键值对
- `GetNewest`: 获取最新的键值对
- `Resize`: 调整容量

---

## 3. 使用示例

### 3.1 基本使用

```go
import "github.com/yourusername/golang/pkg/utils/lru"

// 创建LRU缓存（容量为3）
cache := lru.NewLRUCache[string, int](3)

// 设置值
cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3)

// 获取值
val, ok := cache.Get("a")
if ok {
    fmt.Printf("Value: %d\n", val)
}

// 添加新值，会自动淘汰最旧的
cache.Put("d", 4)  // "b"会被淘汰
```

### 3.2 访问顺序

```go
cache := lru.NewLRUCache[string, int](3)

cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3)

// 访问"a"会将其移到最前面
cache.Get("a")

// 添加新值，"b"会被淘汰（因为"a"和"c"最近被访问）
cache.Put("d", 4)
```

### 3.3 删除和清空

```go
// 删除键
deleted := cache.Delete("a")

// 检查键是否存在
exists := cache.Contains("a")

// 清空缓存
cache.Clear()
```

### 3.4 查看和遍历

```go
// 查看值（不更新访问顺序）
val, ok := cache.Peek("a")

// 获取最旧的键值对
oldestKey, oldestVal, ok := cache.GetOldest()

// 获取最新的键值对
newestKey, newestVal, ok := cache.GetNewest()

// 获取所有键
keys := cache.Keys()

// 获取所有值
values := cache.Values()
```

### 3.5 调整容量

```go
// 调整容量
cache.Resize(5)

// 如果当前大小超过新容量，会自动删除多余的条目
```

### 3.6 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/lru"
)

func main() {
    // 创建LRU缓存
    cache := lru.NewLRUCache[string, string](3)
    
    // 添加数据
    cache.Put("user1", "Alice")
    cache.Put("user2", "Bob")
    cache.Put("user3", "Charlie")
    
    // 访问数据
    val, ok := cache.Get("user1")
    if ok {
        fmt.Printf("User1: %s\n", val)
    }
    
    // 添加新数据，会淘汰最旧的
    cache.Put("user4", "David")
    
    // 检查数据
    if cache.Contains("user2") {
        fmt.Println("user2 exists")
    } else {
        fmt.Println("user2 was evicted")
    }
    
    // 获取缓存大小
    fmt.Printf("Cache size: %d\n", cache.Size())
}
```

---

**更新日期**: 2025-11-11

