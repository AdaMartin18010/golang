# 重试工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [重试工具](#重试工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 重试策略](#2-重试策略)
    - [2.1 指数退避策略](#21-指数退避策略)
    - [2.2 线性退避策略](#22-线性退避策略)
    - [2.3 固定延迟策略](#23-固定延迟策略)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 带回调的重试](#32-带回调的重试)
    - [3.3 自定义策略](#33-自定义策略)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 选择合适策略](#41-选择合适策略)
    - [4.2 设置合理的重试次数](#42-设置合理的重试次数)
    - [4.3 使用Context控制超时](#43-使用context控制超时)

---

## 1. 概述

重试工具提供了灵活的重试机制，支持多种重试策略：

- ✅ **指数退避策略**: 延迟时间指数增长
- ✅ **线性退避策略**: 延迟时间线性增长
- ✅ **固定延迟策略**: 固定延迟时间
- ✅ **Context支持**: 支持取消和超时
- ✅ **回调支持**: 重试时执行回调函数

---

## 2. 重试策略

### 2.1 指数退避策略

延迟时间按指数增长：`delay = initialDelay * multiplier^(attempt-1)`

```go
strategy := &retry.ExponentialBackoff{
    InitialDelay: 100 * time.Millisecond,
    MaxDelay:     30 * time.Second,
    Multiplier:   2.0,
    MaxAttempts:  5,
}
```

### 2.2 线性退避策略

延迟时间按线性增长：`delay = initialDelay + increment * (attempt-1)`

```go
strategy := &retry.LinearBackoff{
    InitialDelay: 100 * time.Millisecond,
    MaxDelay:     5 * time.Second,
    Increment:    100 * time.Millisecond,
    MaxAttempts:  5,
}
```

### 2.3 固定延迟策略

延迟时间固定不变

```go
strategy := &retry.FixedBackoff{
    Delay:       1 * time.Second,
    MaxAttempts: 3,
}
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "context"
    "github.com/yourusername/golang/pkg/utils/retry"
)

ctx := context.Background()
strategy := retry.DefaultExponentialBackoff()

err := retry.Retry(ctx, strategy, func(ctx context.Context) error {
    // 执行可能失败的操作
    return someOperation()
})
```

### 3.2 带回调的重试

```go
err := retry.RetryWithCallback(ctx, strategy, func(ctx context.Context) error {
    return someOperation()
}, func(attempt int, err error) {
    log.Printf("Retry attempt %d, error: %v", attempt, err)
})
```

### 3.3 自定义策略

```go
strategy := &retry.ExponentialBackoff{
    InitialDelay: 200 * time.Millisecond,
    MaxDelay:     10 * time.Second,
    Multiplier:   1.5,
    MaxAttempts:  10,
}

err := retry.Retry(ctx, strategy, func(ctx context.Context) error {
    return apiCall()
})
```

---

## 4. 最佳实践

### 4.1 选择合适策略

- **指数退避**: 适用于网络请求、API调用
- **线性退避**: 适用于资源竞争场景
- **固定延迟**: 适用于简单重试场景

### 4.2 设置合理的重试次数

- 网络请求: 3-5次
- 数据库操作: 2-3次
- 文件操作: 1-2次

### 4.3 使用Context控制超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := retry.Retry(ctx, strategy, func(ctx context.Context) error {
    return operation(ctx)
})
```

---

**更新日期**: 2025-11-11
