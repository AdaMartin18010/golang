# 布隆过滤器工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [布隆过滤器工具](#布隆过滤器工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

布隆过滤器工具提供了布隆过滤器实现，用于快速判断元素是否可能存在于集合中，适用于大规模数据去重、缓存穿透防护等场景。

---

## 2. 功能特性

### 2.1 布隆过滤器

- `BloomFilter`: 布隆过滤器实现
- `NewBloomFilter`: 创建布隆过滤器
- `Add`: 添加元素（字节数组）
- `AddString`: 添加字符串元素
- `Contains`: 检查元素是否存在（字节数组）
- `ContainsString`: 检查字符串元素是否存在
- `Clear`: 清空布隆过滤器
- `Size`: 获取位数组大小
- `Count`: 估算元素数量（近似值）
- `FalsePositiveRate`: 计算假阳性率

### 2.2 优化函数

- `OptimalSize`: 计算最优位数组大小
- `OptimalHashCount`: 计算最优哈希函数数量

---

## 3. 使用示例

### 3.1 基本使用

```go
import "github.com/yourusername/golang/pkg/utils/bloom"

// 创建布隆过滤器（位数组大小1000，3个哈希函数）
bf := bloom.NewBloomFilter(1000, 3)

// 添加元素
bf.AddString("user1")
bf.AddString("user2")
bf.AddString("user3")

// 检查元素是否存在
if bf.ContainsString("user1") {
    fmt.Println("user1 exists (or false positive)")
}

if !bf.ContainsString("user999") {
    fmt.Println("user999 definitely does not exist")
}
```

### 3.2 字节数组操作

```go
// 添加字节数组
bf.Add([]byte("data1"))
bf.Add([]byte("data2"))

// 检查字节数组
if bf.Contains([]byte("data1")) {
    fmt.Println("data1 exists")
}
```

### 3.3 优化配置

```go
// 计算最优位数组大小（1000个元素，1%假阳性率）
optimalSize := bloom.OptimalSize(1000, 0.01)
fmt.Printf("Optimal size: %d\n", optimalSize)

// 计算最优哈希函数数量
optimalHashes := bloom.OptimalHashCount(1000, optimalSize)
fmt.Printf("Optimal hash count: %d\n", optimalHashes)

// 使用优化配置创建布隆过滤器
bf := bloom.NewBloomFilter(optimalSize, optimalHashes)
```

### 3.4 估算和统计

```go
// 估算元素数量
count := bf.Count()
fmt.Printf("Estimated count: %d\n", count)

// 计算假阳性率
falsePositiveRate := bf.FalsePositiveRate(1000)
fmt.Printf("False positive rate: %.4f\n", falsePositiveRate)
```

### 3.5 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/bloom"
)

func main() {
    // 创建布隆过滤器
    bf := bloom.NewBloomFilter(10000, 3)
    
    // 添加用户ID
    users := []string{"user1", "user2", "user3", "user4", "user5"}
    for _, user := range users {
        bf.AddString(user)
    }
    
    // 检查用户是否存在
    testUsers := []string{"user1", "user2", "user999"}
    for _, user := range testUsers {
        if bf.ContainsString(user) {
            fmt.Printf("%s: exists (or false positive)\n", user)
        } else {
            fmt.Printf("%s: definitely does not exist\n", user)
        }
    }
    
    // 估算元素数量
    count := bf.Count()
    fmt.Printf("Estimated count: %d\n", count)
}
```

---

**更新日期**: 2025-11-11

