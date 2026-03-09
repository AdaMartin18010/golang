# NATS 使用文档

> **版本**: v1.0
> **日期**: 2025-01-XX
> **位置**: `internal/infrastructure/messaging/nats/`

---

## 📋 目录

- [NATS 使用文档](#nats-使用文档)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
    - [特性](#特性)
  - [2. 快速开始](#2-快速开始)
    - [安装依赖](#安装依赖)
    - [基本使用](#基本使用)
  - [3. 核心功能](#3-核心功能)
    - [3.1 发布/订阅](#31-发布订阅)
    - [3.2 Request/Reply](#32-requestreply)
    - [3.3 队列订阅](#33-队列订阅)
  - [4. 配置说明](#4-配置说明)
    - [默认配置](#默认配置)
    - [自定义配置](#自定义配置)
  - [5. 使用示例](#5-使用示例)
  - [6. 最佳实践](#6-最佳实践)
  - [📚 相关资源](#-相关资源)

---

## 1. 概述

NATS 是一个高性能、云原生的消息传递系统，专为微服务、IoT 和云原生应用设计。

### 特性

- ✅ **高性能**: 微秒级延迟，高吞吐量
- ✅ **轻量级**: 协议简单，资源占用小
- ✅ **云原生**: 支持集群、流式处理
- ✅ **可靠性**: 支持自动重连和连接保持

---

## 2. 快速开始

### 安装依赖

```bash
go get github.com/nats-io/nats.go@v1.35.0
```

### 基本使用

```go
import "github.com/yourusername/golang/internal/infrastructure/messaging/nats"

// 创建客户端
client, err := nats.NewClient(nats.DefaultConfig())
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 发布消息
err = client.Publish("user.created", map[string]interface{}{
    "user_id": 123,
    "name":    "Alice",
})
```

---

## 3. 核心功能

### 3.1 发布/订阅

```go
// 订阅
sub, err := client.Subscribe("user.created", func(msg *nats.Msg) {
    log.Printf("Received: %s", string(msg.Data))
})
defer sub.Unsubscribe()

// 发布
err = client.Publish("user.created", "message data")
```

### 3.2 Request/Reply

```go
// 服务端
sub, err := client.Subscribe("user.get", func(msg *nats.Msg) {
    msg.Respond([]byte("response data"))
})

// 客户端
reply, err := client.Request("user.get", "request data", 5*time.Second)
```

### 3.3 队列订阅

```go
// 负载均衡订阅
sub, err := client.QueueSubscribe("tasks", "worker-group", func(msg *nats.Msg) {
    // 处理任务
})
```

---

## 4. 配置说明

### 默认配置

```go
config := nats.DefaultConfig()
// URL: "nats://localhost:4222"
// MaxReconnects: -1 (无限重连)
// ReconnectWait: 2秒
// Timeout: 5秒
```

### 自定义配置

```go
config := nats.Config{
    URL:           "nats://localhost:4222",
    MaxReconnects: 10,
    ReconnectWait: 2 * time.Second,
    Timeout:       5 * time.Second,
    Name:          "my-client",
    Username:      "user",
    Password:      "pass",
}
```

---

## 5. 使用示例

完整示例请参考：

- `examples/messaging/nats/publish_subscribe.go`
- `examples/messaging/nats/request_reply.go`

---

## 6. 最佳实践

1. **连接复用**: 在应用程序生命周期中复用客户端实例
2. **错误处理**: 始终检查错误并实现重试机制
3. **资源清理**: 使用 defer 确保连接和订阅被正确关闭
4. **消息序列化**: 使用 JSON 进行消息序列化

---

## 📚 相关资源

- [NATS 官方文档](https://docs.nats.io/)
- [NATS Go 客户端](https://github.com/nats-io/nats.go)
- [代码实现](../internal/infrastructure/messaging/nats/)

---

**最后更新**: 2025-01-XX
