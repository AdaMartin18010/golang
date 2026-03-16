# 限流工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [限流工具](#限流工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

限流工具提供了多种限流算法实现，包括令牌桶、漏桶、滑动窗口、固定窗口等，帮助开发者控制请求速率。

---

## 2. 功能特性

### 2.1 令牌桶限流器

- `TokenBucket`: 令牌桶限流器实现
- `NewTokenBucket`: 创建令牌桶限流器
- `Allow`: 检查是否允许通过
- `Wait`: 等待直到允许通过
- `Reset`: 重置限流器

### 2.2 漏桶限流器

- `LeakyBucket`: 漏桶限流器实现
- `NewLeakyBucket`: 创建漏桶限流器
- `Allow`: 检查是否允许通过
- `Wait`: 等待直到允许通过
- `Reset`: 重置限流器

### 2.3 滑动窗口限流器

- `SlidingWindow`: 滑动窗口限流器实现
- `NewSlidingWindow`: 创建滑动窗口限流器
- `Allow`: 检查是否允许通过
- `Wait`: 等待直到允许通过
- `Reset`: 重置限流器

### 2.4 固定窗口限流器

- `FixedWindow`: 固定窗口限流器实现
- `NewFixedWindow`: 创建固定窗口限流器
- `Allow`: 检查是否允许通过
- `Wait`: 等待直到允许通过
- `Reset`: 重置限流器

---

## 3. 使用示例

### 3.1 令牌桶限流器

```go
import "github.com/yourusername/golang/pkg/utils/ratelimit"

// 创建令牌桶限流器（容量10，每秒补充5个令牌）
tb := ratelimit.NewTokenBucket(10, 5)

// 检查是否允许通过
if tb.Allow() {
    // 处理请求
    processRequest()
} else {
    // 被限流
    handleRateLimit()
}

// 等待直到允许通过
tb.Wait()
processRequest()

// 重置限流器
tb.Reset()
```

### 3.2 漏桶限流器

```go
// 创建漏桶限流器（容量10，每秒漏出5个）
lb := ratelimit.NewLeakyBucket(10, 5)

// 检查是否允许通过
if lb.Allow() {
    processRequest()
}

// 等待直到允许通过
lb.Wait()
processRequest()
```

### 3.3 滑动窗口限流器

```go
// 创建滑动窗口限流器（窗口大小1秒，最大请求数5）
sw := ratelimit.NewSlidingWindow(1*time.Second, 5)

// 检查是否允许通过
if sw.Allow() {
    processRequest()
}

// 等待直到允许通过
sw.Wait()
processRequest()
```

### 3.4 固定窗口限流器

```go
// 创建固定窗口限流器（窗口大小1秒，最大请求数5）
fw := ratelimit.NewFixedWindow(1*time.Second, 5)

// 检查是否允许通过
if fw.Allow() {
    processRequest()
}

// 等待直到允许通过
fw.Wait()
processRequest()
```

### 3.5 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/golang/pkg/utils/ratelimit"
)

func main() {
    // 令牌桶限流器
    tb := ratelimit.NewTokenBucket(10, 5)
    
    for i := 0; i < 15; i++ {
        if tb.Allow() {
            fmt.Printf("Request %d: Allowed\n", i+1)
        } else {
            fmt.Printf("Request %d: Rate limited\n", i+1)
            time.Sleep(200 * time.Millisecond)
        }
    }
    
    // 滑动窗口限流器
    sw := ratelimit.NewSlidingWindow(1*time.Second, 5)
    
    for i := 0; i < 10; i++ {
        if sw.Allow() {
            fmt.Printf("Request %d: Allowed\n", i+1)
        } else {
            fmt.Printf("Request %d: Rate limited\n", i+1)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

---

**更新日期**: 2025-11-11

