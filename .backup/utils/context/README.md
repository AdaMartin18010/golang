# 上下文工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [上下文工具](#上下文工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

上下文工具提供了context包的便捷封装，简化context的创建、管理和使用任务。

---

## 2. 功能特性

### 2.1 Context创建

- `WithTimeout`: 创建带超时的context
- `WithDeadline`: 创建带截止时间的context
- `WithCancel`: 创建可取消的context
- `WithValue`: 创建带值的context
- `Background`: 返回非nil的空context
- `TODO`: 返回非nil的空context
- `WithTimeoutSeconds`: 创建带超时的context（秒为单位）
- `WithTimeoutMinutes`: 创建带超时的context（分钟为单位）
- `WithTimeoutHours`: 创建带超时的context（小时为单位）

### 2.2 Context检查

- `IsDone`: 检查context是否已取消
- `IsCancelled`: 检查context是否已取消（别名）
- `GetError`: 获取context的错误
- `GetDeadline`: 获取context的截止时间

### 2.3 Context值操作

- `GetValue`: 获取context中的值
- `GetStringValue`: 获取context中的字符串值
- `GetIntValue`: 获取context中的整数值
- `GetInt64Value`: 获取context中的64位整数值
- `GetBoolValue`: 获取context中的布尔值
- `MustGetStringValue`: 获取字符串值，如果不存在则panic
- `MustGetIntValue`: 获取整数值，如果不存在则panic
- `MustGetBoolValue`: 获取布尔值，如果不存在则panic
- `WithValues`: 批量设置context值
- `WithStringValue`: 使用字符串键设置值
- `GetStringKeyValue`: 使用字符串键获取值

### 2.4 Context等待

- `Wait`: 等待context完成
- `WaitWithTimeout`: 等待context完成或超时
- `Sleep`: 睡眠，但可以被context取消

### 2.5 Context执行

- `DoWithTimeout`: 在超时时间内执行函数
- `DoWithDeadline`: 在截止时间前执行函数
- `DoWithCancel`: 执行函数，支持取消
- `RetryWithContext`: 使用context重试函数

### 2.6 Context合并

- `Merge`: 合并多个context（任一取消则取消）

### 2.7 常用键

- `KeyTraceID`: TraceID键
- `KeySpanID`: SpanID键
- `KeyUserID`: UserID键
- `KeyRequestID`: RequestID键
- `KeyIP`: IP键
- `KeyUserAgent`: UserAgent键

### 2.8 快捷函数

- `WithTraceID`: 设置TraceID
- `GetTraceID`: 获取TraceID
- `WithSpanID`: 设置SpanID
- `GetSpanID`: 获取SpanID
- `WithUserID`: 设置UserID
- `GetUserID`: 获取UserID
- `WithRequestID`: 设置RequestID
- `GetRequestID`: 获取RequestID
- `WithIP`: 设置IP
- `GetIP`: 获取IP
- `WithUserAgent`: 设置UserAgent
- `GetUserAgent`: 获取UserAgent

### 2.9 Context构建器

- `Chain`: 链式设置多个值
- `ContextBuilder`: 上下文构建器

---

## 3. 使用示例

### 3.1 Context创建

```go
import "github.com/yourusername/golang/pkg/utils/context"

// 创建带超时的context
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

// 创建带超时的context（秒为单位）
ctx, cancel := context.WithTimeoutSeconds(context.Background(), 5)

// 创建带截止时间的context
deadline := time.Now().Add(time.Hour)
ctx, cancel := context.WithDeadline(context.Background(), deadline)

// 创建可取消的context
ctx, cancel := context.WithCancel(context.Background())
```

### 3.2 Context检查

```go
// 检查context是否已取消
if context.IsDone(ctx) {
    // context已取消
}

// 获取context的错误
err := context.GetError(ctx)

// 获取context的截止时间
deadline, ok := context.GetDeadline(ctx)
```

### 3.3 Context值操作

```go
// 设置值
ctx := context.WithValue(context.Background(), "key", "value")

// 获取值
value := context.GetValue(ctx, "key")

// 获取字符串值
str, ok := context.GetStringValue(ctx, "key")

// 获取整数值
num, ok := context.GetIntValue(ctx, "key")

// 使用字符串键设置值
ctx = context.WithStringValue(ctx, "key", "value")
value := context.GetStringKeyValue(ctx, "key")
```

### 3.4 Context等待

```go
// 等待context完成
context.Wait(ctx)

// 等待context完成或超时
completed := context.WaitWithTimeout(ctx, time.Second)

// 睡眠，但可以被context取消
err := context.Sleep(ctx, time.Second)
```

### 3.5 Context执行

```go
// 在超时时间内执行函数
err := context.DoWithTimeout(ctx, time.Second, func(ctx context.Context) error {
    // 执行操作
    return nil
})

// 在截止时间前执行函数
err := context.DoWithDeadline(ctx, deadline, func(ctx context.Context) error {
    // 执行操作
    return nil
})

// 使用context重试函数
err := context.RetryWithContext(ctx, 3, func(ctx context.Context) error {
    // 执行操作
    return nil
})
```

### 3.6 常用键操作

```go
// 设置TraceID
ctx := context.WithTraceID(context.Background(), "trace123")

// 获取TraceID
traceID, ok := context.GetTraceID(ctx)

// 设置UserID
ctx = context.WithUserID(ctx, "user123")

// 获取UserID
userID, ok := context.GetUserID(ctx)

// 设置RequestID
ctx = context.WithRequestID(ctx, "req123")

// 获取RequestID
requestID, ok := context.GetRequestID(ctx)
```

### 3.7 Context构建器

```go
// 链式设置多个值
ctx := context.Chain(context.Background()).
    WithTraceID("trace123").
    WithUserID("user123").
    WithRequestID("req123").
    WithStringValue("custom_key", "custom_value").
    Build()
```

### 3.8 Context合并

```go
// 合并多个context（任一取消则取消）
ctx1, cancel1 := context.WithCancel(context.Background())
ctx2, cancel2 := context.WithCancel(context.Background())
merged, cancel := context.Merge(ctx1, ctx2)
defer cancel()
```

### 3.9 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/golang/pkg/utils/context"
)

func main() {
    // 创建带超时的context
    ctx, cancel := context.WithTimeoutSeconds(context.Background(), 5)
    defer cancel()

    // 设置常用值
    ctx = context.WithTraceID(ctx, "trace123")
    ctx = context.WithUserID(ctx, "user123")

    // 执行操作
    err := context.DoWithTimeout(ctx, 3*time.Second, func(ctx context.Context) error {
        // 获取值
        traceID, _ := context.GetTraceID(ctx)
        userID, _ := context.GetUserID(ctx)

        fmt.Printf("TraceID: %s, UserID: %s\n", traceID, userID)

        // 执行操作
        time.Sleep(2 * time.Second)
        return nil
    })

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

---

**更新日期**: 2025-11-11
