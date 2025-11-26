# 对象池工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [对象池工具](#对象池工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

对象池工具提供了多种对象池实现，包括简单对象池、有界对象池、缓冲区池、字符串构建器池等，帮助开发者减少对象分配，提高性能。

---

## 2. 功能特性

### 2.1 简单对象池

- `SimplePool`: 简单对象池实现
- `NewSimplePool`: 创建简单对象池
- `Get`: 获取对象
- `Put`: 归还对象
- `Clear`: 清空对象池
- `Size`: 获取对象池大小

### 2.2 有界对象池

- `BoundedPool`: 有界对象池实现
- `NewBoundedPool`: 创建有界对象池
- `Get`: 获取对象
- `Put`: 归还对象
- `Clear`: 清空对象池
- `Size`: 获取对象池大小
- `Capacity`: 获取对象池容量

### 2.3 缓冲区池

- `BufferPool`: 缓冲区池实现
- `NewBufferPool`: 创建缓冲区池
- `Get`: 获取缓冲区
- `Put`: 归还缓冲区

### 2.4 字符串构建器池

- `StringBuilderPool`: 字符串构建器池实现
- `NewStringBuilderPool`: 创建字符串构建器池
- `Get`: 获取字符串构建器
- `Put`: 归还字符串构建器

---

## 3. 使用示例

### 3.1 简单对象池

```go
import "github.com/yourusername/golang/pkg/utils/pool"

// 创建简单对象池
p := pool.NewSimplePool[[]byte](func() []byte {
    return make([]byte, 0, 1024)
})

// 获取对象
buf := p.Get()

// 使用对象
buf = append(buf, []byte("data")...)

// 归还对象
p.Put(buf)

// 清空对象池
p.Clear()
```

### 3.2 有界对象池

```go
// 创建有界对象池（最大容量10）
p := pool.NewBoundedPool[[]byte](10, func() []byte {
    return make([]byte, 0, 1024)
})

// 获取对象
buf := p.Get()

// 使用对象
buf = append(buf, []byte("data")...)

// 归还对象
p.Put(buf)

// 获取容量
capacity := p.Capacity()
```

### 3.3 缓冲区池

```go
// 创建缓冲区池
bp := pool.NewBufferPool()

// 获取缓冲区
buf := bp.Get()

// 使用缓冲区
buf = append(buf, []byte("data")...)

// 归还缓冲区
bp.Put(buf)
```

### 3.4 字符串构建器池

```go
// 创建字符串构建器池
sbp := pool.NewStringBuilderPool()

// 获取字符串构建器
sb := sbp.Get()
sb.WriteString("Hello")
sb.WriteString(" World")

// 获取字符串
result := sb.String()

// 归还字符串构建器
sbp.Put(sb)
```

### 3.5 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/pool"
)

func main() {
    // 简单对象池
    p := pool.NewSimplePool[[]byte](func() []byte {
        return make([]byte, 0, 1024)
    })
    
    buf := p.Get()
    buf = append(buf, []byte("test")...)
    p.Put(buf)
    
    // 缓冲区池
    bp := pool.NewBufferPool()
    buffer := bp.Get()
    buffer = append(buffer, []byte("data")...)
    bp.Put(buffer)
    
    // 字符串构建器池
    sbp := pool.NewStringBuilderPool()
    sb := sbp.Get()
    sb.WriteString("Hello")
    sb.WriteString(" World")
    fmt.Println(sb.String())
    sbp.Put(sb)
}
```

---

**更新日期**: 2025-11-11

