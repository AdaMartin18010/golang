# 健康检查框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [健康检查框架](#健康检查框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心功能](#2-核心功能)
    - [2.1 健康状态](#21-健康状态)
    - [2.2 检查类型](#22-检查类型)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 带超时的检查](#32-带超时的检查)
    - [3.3 定期检查（带缓存）](#33-定期检查带缓存)
    - [3.4 聚合检查](#34-聚合检查)
    - [3.5 在HTTP Handler中使用](#35-在http-handler中使用)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

健康检查框架提供了完整的服务健康检查功能：

- ✅ **多种检查类型**: 简单检查、超时检查、定期检查、聚合检查
- ✅ **健康状态管理**: 健康、不健康、降级状态
- ✅ **检查结果缓存**: 定期检查支持结果缓存
- ✅ **检查聚合**: 支持多个检查的聚合

---

## 2. 核心功能

### 2.1 健康状态

- `StatusHealthy` - 健康
- `StatusUnhealthy` - 不健康
- `StatusDegraded` - 降级

### 2.2 检查类型

- **SimpleCheck**: 简单健康检查
- **TimeoutCheck**: 带超时的健康检查
- **PeriodicCheck**: 定期健康检查（带缓存）
- **AggregateCheck**: 聚合健康检查

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "context"
    "github.com/yourusername/golang/pkg/health"
)

// 创建健康检查器
checker := health.NewHealthChecker()

// 注册简单检查
checker.Register(health.NewSimpleCheck("database", func(ctx context.Context) error {
    // 检查数据库连接
    return db.Ping(ctx)
}))

// 执行所有检查
results := checker.Check(context.Background())

// 获取整体状态
status := checker.OverallStatus(context.Background())
```

### 3.2 带超时的检查

```go
dbCheck := health.NewSimpleCheck("database", func(ctx context.Context) error {
    return db.Ping(ctx)
})

timeoutCheck := health.NewTimeoutCheck("database-timeout", 5*time.Second, dbCheck)
checker.Register(timeoutCheck)
```

### 3.3 定期检查（带缓存）

```go
dbCheck := health.NewSimpleCheck("database", func(ctx context.Context) error {
    return db.Ping(ctx)
})

// 每30秒检查一次，结果缓存30秒
periodicCheck := health.NewPeriodicCheck("database-periodic", 30*time.Second, dbCheck)
checker.Register(periodicCheck)
```

### 3.4 聚合检查

```go
dbCheck := health.NewSimpleCheck("database", func(ctx context.Context) error {
    return db.Ping(ctx)
})

cacheCheck := health.NewSimpleCheck("cache", func(ctx context.Context) error {
    return redis.Ping(ctx).Err()
})

// 聚合多个检查
aggregateCheck := health.NewAggregateCheck("storage", dbCheck, cacheCheck)
checker.Register(aggregateCheck)
```

### 3.5 在HTTP Handler中使用

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    status := checker.OverallStatus(r.Context())
    results := checker.Check(r.Context())

    response := map[string]interface{}{
        "status": status,
        "checks": results,
    }

    code := http.StatusOK
    if status == health.StatusUnhealthy {
        code = http.StatusServiceUnavailable
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(response)
}
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用超时**: 为所有外部依赖检查设置超时
2. **定期检查**: 对频繁检查使用定期检查以减少开销
3. **聚合相关检查**: 将相关的检查聚合在一起
4. **提供详细信息**: 在检查结果中包含有用的错误信息

### 4.2 DON'Ts ❌

1. **不要阻塞**: 健康检查不应该阻塞太久
2. **不要忽略错误**: 正确处理检查错误
3. **不要过度检查**: 避免过于频繁的健康检查

---

## 5. 相关资源

- [服务注册中心](../registry/README.md)
- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
