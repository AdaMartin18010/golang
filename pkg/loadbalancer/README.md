# 负载均衡器

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26.2

---

## 📋 目录

- [负载均衡器](#负载均衡器)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 支持的算法](#2-支持的算法)
    - [2.1 轮询算法 (Round Robin)](#21-轮询算法-round-robin)
    - [2.2 随机算法 (Random)](#22-随机算法-random)
    - [2.3 加权轮询 (Weighted Round Robin)](#23-加权轮询-weighted-round-robin)
    - [2.4 最少连接 (Least Connections)](#24-最少连接-least-connections)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 与服务注册中心集成](#32-与服务注册中心集成)
    - [3.3 最少连接算法](#33-最少连接算法)
    - [3.4 加权轮询](#34-加权轮询)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

负载均衡器提供了多种负载均衡算法：

- ✅ **轮询算法**: 依次选择服务实例
- ✅ **随机算法**: 随机选择服务实例
- ✅ **加权轮询**: 根据权重选择服务实例
- ✅ **最少连接**: 选择连接数最少的服务实例

---

## 2. 支持的算法

### 2.1 轮询算法 (Round Robin)

依次选择服务实例，适合所有服务实例性能相近的场景。

### 2.2 随机算法 (Random)

随机选择服务实例，适合服务实例性能相近且请求分布均匀的场景。

### 2.3 加权轮询 (Weighted Round Robin)

根据服务实例的权重选择，适合服务实例性能不同的场景。

### 2.4 最少连接 (Least Connections)

选择当前连接数最少的服务实例，适合长连接场景。

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "context"
    "github.com/yourusername/golang/pkg/loadbalancer"
    "github.com/yourusername/golang/pkg/registry"
)

// 创建负载均衡器
lb := loadbalancer.NewRoundRobin()

// 选择服务
services := []*registry.Service{
    {ID: "service-1", Name: "user-service"},
    {ID: "service-2", Name: "user-service"},
}

selected, err := lb.Select(context.Background(), services)
if err != nil {
    // 处理错误
}
```

### 3.2 与服务注册中心集成

```go
// 创建服务选择器
reg := registry.NewInMemoryRegistry()
lb := loadbalancer.NewRoundRobin()
selector := loadbalancer.NewServiceSelector(reg, lb, "user-service")

// 选择服务
service, err := selector.Select(context.Background())
if err != nil {
    // 处理错误
}

// 使用服务
url := fmt.Sprintf("http://%s:%d", service.Address, service.Port)
```

### 3.3 最少连接算法

```go
lb := loadbalancer.NewLeastConnections()

// 选择服务
service, err := lb.Select(context.Background(), services)
if err != nil {
    // 处理错误
}

// 使用服务后释放连接
defer lb.Release(service.ID)
```

### 3.4 加权轮询

```go
lb := loadbalancer.NewWeightedRoundRobin()

services := []*registry.Service{
    {
        ID: "service-1",
        Name: "user-service",
        Metadata: map[string]string{"weight": "3"},
    },
    {
        ID: "service-2",
        Name: "user-service",
        Metadata: map[string]string{"weight": "1"},
    },
}

selected, err := lb.Select(context.Background(), services)
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **选择合适的算法**: 根据场景选择合适的负载均衡算法
2. **健康检查**: 结合健康检查过滤不健康的服务
3. **连接管理**: 使用最少连接算法时记得释放连接
4. **权重配置**: 根据服务实例性能配置合理的权重

### 4.2 DON'Ts ❌

1. **不要忽略错误**: 选择服务可能失败
2. **不要忘记释放连接**: 使用最少连接算法时必须释放连接
3. **不要使用过时的服务列表**: 定期更新服务列表

---

## 5. 相关资源

- [服务注册中心](../registry/README.md)
- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
