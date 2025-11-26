# 服务注册中心

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [服务注册中心](#服务注册中心)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心功能](#2-核心功能)
    - [2.1 Service 结构](#21-service-结构)
    - [2.2 Registry 接口](#22-registry-接口)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本使用](#31-基本使用)
    - [3.2 服务监听](#32-服务监听)
    - [3.3 服务心跳](#33-服务心跳)
    - [3.4 清理过期服务](#34-清理过期服务)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 DO's ✅](#41-dos-)
    - [4.2 DON'Ts ❌](#42-donts-)
  - [5. 相关资源](#5-相关资源)

---

## 1. 概述

服务注册中心提供了服务注册与发现功能：

- ✅ **服务注册**: 注册服务实例
- ✅ **服务发现**: 发现服务实例
- ✅ **服务监听**: 监听服务变化
- ✅ **健康检查**: 服务健康检查
- ✅ **过期清理**: 自动清理过期服务

---

## 2. 核心功能

### 2.1 Service 结构

```go
type Service struct {
    ID       string            // 服务ID
    Name     string            // 服务名称
    Address  string            // 服务地址
    Port     int               // 服务端口
    Tags     []string          // 标签
    Metadata map[string]string // 元数据
    TTL      time.Duration     // 生存时间
    LastSeen time.Time         // 最后更新时间
}
```

### 2.2 Registry 接口

```go
type Registry interface {
    Register(ctx context.Context, service *Service) error
    Deregister(ctx context.Context, serviceID string) error
    GetService(ctx context.Context, serviceID string) (*Service, error)
    ListServices(ctx context.Context, name string) ([]*Service, error)
    Watch(ctx context.Context, name string) (<-chan []*Service, error)
    Health(ctx context.Context) error
}
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "context"
    "github.com/yourusername/golang/pkg/registry"
)

// 创建注册中心
reg := registry.NewInMemoryRegistry()

// 注册服务
service := &registry.Service{
    ID:      "user-service-1",
    Name:    "user-service",
    Address: "localhost",
    Port:    8080,
    Tags:    []string{"v1", "production"},
    Metadata: map[string]string{
        "version": "1.0.0",
    },
    TTL: 30 * time.Second,
}

err := reg.Register(context.Background(), service)
if err != nil {
    // 处理错误
}

// 发现服务
services, err := reg.ListServices(context.Background(), "user-service")
if err != nil {
    // 处理错误
}
```

### 3.2 服务监听

```go
// 监听服务变化
ch, err := reg.Watch(context.Background(), "user-service")
if err != nil {
    // 处理错误
}

go func() {
    for services := range ch {
        // 处理服务列表变化
        fmt.Printf("Services updated: %d instances\n", len(services))
    }
}()
```

### 3.3 服务心跳

```go
// 定期更新服务心跳
ticker := time.NewTicker(10 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        service.LastSeen = time.Now()
        reg.Register(context.Background(), service)
    case <-ctx.Done():
        return
    }
}
```

### 3.4 清理过期服务

```go
// 定期清理过期服务
ticker := time.NewTicker(1 * time.Minute)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        reg.CleanupExpiredServices(context.Background(), 60*time.Second)
    case <-ctx.Done():
        return
    }
}
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用TTL**: 为服务设置合理的TTL
2. **定期心跳**: 定期更新服务心跳
3. **优雅注销**: 应用关闭时注销服务
4. **监听变化**: 使用Watch监听服务变化
5. **健康检查**: 定期检查服务健康状态

### 4.2 DON'Ts ❌

1. **不要忘记注销**: 应用关闭时必须注销服务
2. **不要设置过长的TTL**: TTL过长会导致服务不可用检测延迟
3. **不要忽略错误**: 注册和注销操作可能失败
4. **不要阻塞监听**: Watch操作不应该阻塞主流程

---

## 5. 相关资源

- [负载均衡器](../loadbalancer/README.md)
- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
