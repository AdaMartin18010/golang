# 定时器工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [定时器工具](#定时器工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

定时器工具提供了多种定时器实现，包括简单定时器、一次性定时器、防抖定时器、节流定时器、间隔定时器等，帮助开发者处理各种定时任务场景。

---

## 2. 功能特性

### 2.1 简单定时器

- `SimpleTimer`: 简单定时器实现
- `NewSimpleTimer`: 创建简单定时器
- `Start`: 启动定时器
- `Stop`: 停止定时器
- `Reset`: 重置定时器
- `IsRunning`: 检查是否运行中

### 2.2 一次性定时器

- `OneShotTimer`: 一次性定时器实现
- `NewOneShotTimer`: 创建一次性定时器
- `Start`: 启动定时器
- `Stop`: 停止定时器
- `Reset`: 重置定时器
- `IsRunning`: 检查是否运行中

### 2.3 防抖定时器

- `DebounceTimer`: 防抖定时器实现
- `NewDebounceTimer`: 创建防抖定时器
- `Trigger`: 触发防抖
- `Cancel`: 取消防抖

### 2.4 节流定时器

- `ThrottleTimer`: 节流定时器实现
- `NewThrottleTimer`: 创建节流定时器
- `Trigger`: 触发节流
- `Reset`: 重置节流定时器

### 2.5 间隔定时器

- `IntervalTimer`: 间隔定时器实现
- `NewIntervalTimer`: 创建间隔定时器
- `Start`: 启动定时器
- `Stop`: 停止定时器
- `Reset`: 重置定时器
- `IsRunning`: 检查是否运行中
- `ExecutionCount`: 获取执行次数
- `ResetCount`: 重置执行次数

### 2.6 快捷函数

- `After`: 延迟执行
- `Every`: 定期执行
- `Schedule`: 调度执行

---

## 3. 使用示例

### 3.1 简单定时器

```go
import "github.com/yourusername/golang/pkg/utils/timer"

// 创建简单定时器
t := timer.NewSimpleTimer(1*time.Second, func() {
    fmt.Println("Timer ticked")
})

// 启动定时器
t.Start()

// 停止定时器
time.Sleep(5 * time.Second)
t.Stop()
```

### 3.2 一次性定时器

```go
// 创建一次性定时器
t := timer.NewOneShotTimer(2*time.Second, func() {
    fmt.Println("Timer executed once")
})

// 启动定时器
t.Start()
```

### 3.3 防抖定时器

```go
// 创建防抖定时器
debounce := timer.NewDebounceTimer(500*time.Millisecond, func() {
    fmt.Println("Debounced action")
})

// 快速触发多次，只会在最后一次触发后500ms执行
for i := 0; i < 10; i++ {
    debounce.Trigger()
    time.Sleep(50 * time.Millisecond)
}

// 取消防抖
debounce.Cancel()
```

### 3.4 节流定时器

```go
// 创建节流定时器
throttle := timer.NewThrottleTimer(1*time.Second, func() {
    fmt.Println("Throttled action")
})

// 快速触发多次，但每秒最多执行一次
for i := 0; i < 10; i++ {
    throttle.Trigger()
    time.Sleep(100 * time.Millisecond)
}
```

### 3.5 间隔定时器

```go
// 创建间隔定时器
interval := timer.NewIntervalTimer(1*time.Second, func() {
    fmt.Println("Interval ticked")
})

// 启动定时器
interval.Start()

// 获取执行次数
count := interval.ExecutionCount()
fmt.Printf("Executed %d times\n", count)

// 停止定时器
time.Sleep(5 * time.Second)
interval.Stop()
```

### 3.6 快捷函数

```go
// 延迟执行
timer.After(2*time.Second, func() {
    fmt.Println("Executed after 2 seconds")
})

// 定期执行
ticker := timer.Every(1*time.Second, func() {
    fmt.Println("Every second")
})
defer ticker.Stop()

// 调度执行（延迟后定期执行）
ticker = timer.Schedule(2*time.Second, 1*time.Second, func() {
    fmt.Println("Scheduled execution")
})
defer ticker.Stop()
```

### 3.7 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/golang/pkg/utils/timer"
)

func main() {
    // 简单定时器
    t := timer.NewSimpleTimer(1*time.Second, func() {
        fmt.Println("Tick")
    })
    t.Start()
    time.Sleep(3 * time.Second)
    t.Stop()
    
    // 防抖定时器
    debounce := timer.NewDebounceTimer(500*time.Millisecond, func() {
        fmt.Println("Debounced")
    })
    for i := 0; i < 5; i++ {
        debounce.Trigger()
        time.Sleep(100 * time.Millisecond)
    }
    time.Sleep(600 * time.Millisecond)
}
```

---

**更新日期**: 2025-11-11

