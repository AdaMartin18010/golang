# 锁工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [锁工具](#锁工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

锁工具提供了多种锁实现，包括互斥锁、读写锁、键控互斥锁、键控读写锁、自旋锁等，帮助开发者进行并发控制。

---

## 2. 功能特性

### 2.1 互斥锁

- `Mutex`: 互斥锁实现
- `NewMutex`: 创建互斥锁
- `Lock`: 加锁
- `Unlock`: 解锁
- `TryLock`: 尝试加锁
- `TryLockWithTimeout`: 带超时的尝试加锁

### 2.2 读写锁

- `RWMutex`: 读写锁实现
- `NewRWMutex`: 创建读写锁
- `Lock`: 写锁
- `Unlock`: 写解锁
- `RLock`: 读锁
- `RUnlock`: 读解锁
- `TryLock`: 尝试写锁
- `TryRLock`: 尝试读锁
- `TryLockWithTimeout`: 带超时的尝试写锁
- `TryRLockWithTimeout`: 带超时的尝试读锁

### 2.3 键控互斥锁

- `KeyedMutex`: 键控互斥锁实现
- `NewKeyedMutex`: 创建键控互斥锁
- `Lock`: 加锁
- `Unlock`: 解锁
- `TryLock`: 尝试加锁

### 2.4 键控读写锁

- `KeyedRWMutex`: 键控读写锁实现
- `NewKeyedRWMutex`: 创建键控读写锁
- `Lock`: 写锁
- `Unlock`: 写解锁
- `RLock`: 读锁
- `RUnlock`: 读解锁
- `TryLock`: 尝试写锁
- `TryRLock`: 尝试读锁

### 2.5 自旋锁

- `SpinLock`: 自旋锁实现
- `NewSpinLock`: 创建自旋锁
- `Lock`: 加锁
- `Unlock`: 解锁
- `TryLock`: 尝试加锁
- `TryLockWithTimeout`: 带超时的尝试加锁

---

## 3. 使用示例

### 3.1 互斥锁

```go
import "github.com/yourusername/golang/pkg/utils/lock"

// 创建互斥锁
m := lock.NewMutex()

// 加锁
m.Lock()
// 临界区代码
m.Unlock()

// 尝试加锁
if m.TryLock() {
    // 临界区代码
    m.Unlock()
}

// 带超时的尝试加锁
if m.TryLockWithTimeout(1 * time.Second) {
    // 临界区代码
    m.Unlock()
}
```

### 3.2 读写锁

```go
// 创建读写锁
rw := lock.NewRWMutex()

// 读锁
rw.RLock()
// 读取操作
rw.RUnlock()

// 写锁
rw.Lock()
// 写入操作
rw.Unlock()

// 尝试读锁
if rw.TryRLock() {
    // 读取操作
    rw.RUnlock()
}
```

### 3.3 键控互斥锁

```go
// 创建键控互斥锁
km := lock.NewKeyedMutex()

// 不同键可以同时加锁
km.Lock("key1")
km.Lock("key2")

// 临界区代码
km.Unlock("key1")
km.Unlock("key2")

// 相同键不能同时加锁
km.Lock("key1")
if !km.TryLock("key1") {
    fmt.Println("Cannot lock same key twice")
}
km.Unlock("key1")
```

### 3.4 键控读写锁

```go
// 创建键控读写锁
km := lock.NewKeyedRWMutex()

// 读锁
km.RLock("key1")
// 读取操作
km.RUnlock("key1")

// 写锁
km.Lock("key1")
// 写入操作
km.Unlock("key1")
```

### 3.5 自旋锁

```go
// 创建自旋锁
sl := lock.NewSpinLock()

// 加锁
sl.Lock()
// 临界区代码
sl.Unlock()

// 尝试加锁
if sl.TryLock() {
    // 临界区代码
    sl.Unlock()
}
```

### 3.6 完整示例

```go
package main

import (
    "fmt"
    "sync"
    "github.com/yourusername/golang/pkg/utils/lock"
)

func main() {
    // 互斥锁
    m := lock.NewMutex()
    var counter int
    
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            m.Lock()
            counter++
            m.Unlock()
        }()
    }
    
    wg.Wait()
    fmt.Printf("Counter: %d\n", counter)
    
    // 键控互斥锁
    km := lock.NewKeyedMutex()
    km.Lock("user1")
    // 处理user1的数据
    km.Unlock("user1")
}
```

---

**更新日期**: 2025-11-11

