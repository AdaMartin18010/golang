# 工作流架构设计

> **简介**: 本文档介绍本项目使用 Temporal 作为工作流编排引擎的架构设计，包括组件说明、工作流模式和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [工作流架构设计](#工作流架构设计)
  - [📋 目录](#-目录)
  - [1. 📚 概述](#1--概述)
  - [2. 🏗️ 架构图](#2-️-架构图)
  - [3. 🔧 组件说明](#3--组件说明)
  - [4. 🔄 工作流模式](#4--工作流模式)
  - [5. ⚠️ 错误处理](#5-️-错误处理)
  - [6. 📡 信号和查询](#6--信号和查询)
  - [7. 💾 持久化](#7--持久化)
  - [8. 📈 可扩展性](#8--可扩展性)
  - [9. 🔍 监控和调试](#9--监控和调试)
  - [10. 🎯 最佳实践](#10--最佳实践)

---

## 1. 📚 概述

本项目使用 **Temporal** 作为工作流编排引擎，实现可靠的分布式业务流程。Temporal 提供了工作流状态管理、自动重试、持久化等特性，确保业务流程的可靠性。

---

## 2. 🏗️ 架构图

```text
┌─────────────┐
│   Client    │
│  (HTTP/gRPC)│
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Handler   │
│  (Temporal) │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Temporal   │
│   Server    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Worker    │
│ (Activities)│
└─────────────┘
```

---

## 3. 🔧 组件说明

### 3.1 Temporal Server

- **职责**: 工作流状态管理、调度、持久化
- **部署**: Docker 容器或 Kubernetes
- **端口**: 7233 (gRPC), 8088 (UI)

### 3.2 Worker

- **职责**: 执行工作流和活动
- **位置**: `cmd/temporal-worker/`
- **配置**: Task Queue, Namespace

### 3.3 Client

- **职责**: 启动工作流、查询状态、发送信号
- **位置**: `internal/infrastructure/workflow/temporal/`

### 3.4 Workflow

- **职责**: 定义业务流程
- **位置**: `internal/application/workflow/`
- **特点**: 确定性执行、可恢复、可查询

### 3.5 Activity

- **职责**: 执行具体业务逻辑
- **位置**: `internal/application/workflow/`
- **特点**: 可重试、可超时、可取消

---

## 4. 🔄 工作流模式

### 4.1 顺序执行

```go
result1 := workflow.ExecuteActivity(ctx, Activity1, input1)
result2 := workflow.ExecuteActivity(ctx, Activity2, result1)
```

### 4.2 并行执行

```go
future1 := workflow.ExecuteActivity(ctx, Activity1, input1)
future2 := workflow.ExecuteActivity(ctx, Activity2, input2)
workflow.Await(future1, future2)
```

### 4.3 条件执行

```go
if condition {
    workflow.ExecuteActivity(ctx, Activity1, input)
} else {
    workflow.ExecuteActivity(ctx, Activity2, input)
}
```

### 4.4 循环执行

```go
for i := 0; i < 10; i++ {
    workflow.ExecuteActivity(ctx, Activity, i)
}
```

---

## 5. ⚠️ 错误处理

### 5.1 重试策略

```go
ao := workflow.ActivityOptions{
    RetryPolicy: &workflow.RetryPolicy{
        InitialInterval:    time.Second,
        BackoffCoefficient: 2.0,
        MaximumInterval:    time.Minute,
        MaximumAttempts:    3,
    },
}
```

### 5.2 错误传播

工作流中的错误会自动传播，可以通过返回值或错误处理。

---

## 6. 📡 信号和查询

### 6.1 信号（Signal）

用于从外部向工作流发送异步消息：

```go
workflow.GetSignalChannel(ctx, "signal-name")
```

### 6.2 查询（Query）

用于查询工作流状态：

```go
workflow.SetQueryHandler(ctx, "query-name", handler)
```

---

## 7. 💾 持久化

Temporal 自动持久化工作流状态，包括：

- 工作流历史
- 活动结果
- 定时器
- 信号和查询

---

## 8. 📈 可扩展性

### 8.1 水平扩展

- 多个 Worker 可以处理同一个 Task Queue
- Temporal Server 支持集群部署

### 8.2 性能优化

- 使用合适的 Task Queue 分区
- 优化活动执行时间
- 合理设置超时和重试策略

---

## 9. 🔍 监控和调试

### 9.1 Temporal UI

- 查看工作流执行历史
- 调试工作流问题
- 监控工作流性能

### 9.2 日志

- 工作流和活动的日志
- 集成 OpenTelemetry 追踪

---

## 10. 🎯 最佳实践

1. **保持工作流确定性** - 不要使用随机数、时间等非确定性操作
2. **合理设置超时** - 为活动设置合适的超时时间
3. **使用版本控制** - 工作流变更时使用版本控制
4. **错误处理** - 正确处理和传播错误
5. **幂等性** - 确保活动是幂等的

---

## 📚 扩展阅读

- [工作流使用指南](../guides/workflow.md) - 工作流使用指南
- [Temporal 官方文档](https://docs.temporal.io/) - Temporal 官方文档
- [Clean Architecture](./clean-architecture.md) - 架构设计详解

---

> 📚 **简介**
> 本文深入探讨 Temporal 工作流在本项目中的应用，系统讲解工作流架构、组件、模式和最佳实践。通过本文，您将全面掌握工作流的设计和使用方法。
