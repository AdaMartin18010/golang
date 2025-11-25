# 工作流架构设计

## 概述

本项目使用 Temporal 作为工作流编排引擎，实现可靠的分布式业务流程。

## 架构图

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

## 组件说明

### 1. Temporal Server

- **职责**: 工作流状态管理、调度、持久化
- **部署**: Docker 容器或 Kubernetes
- **端口**: 7233 (gRPC), 8088 (UI)

### 2. Worker

- **职责**: 执行工作流和活动
- **位置**: `cmd/temporal-worker/`
- **配置**: Task Queue, Namespace

### 3. Client

- **职责**: 启动工作流、查询状态、发送信号
- **位置**: `internal/infrastructure/workflow/temporal/`

### 4. Workflow

- **职责**: 定义业务流程
- **位置**: `internal/application/workflow/`
- **特点**: 确定性执行、可恢复、可查询

### 5. Activity

- **职责**: 执行具体业务逻辑
- **位置**: `internal/application/workflow/`
- **特点**: 可重试、可超时、可取消

## 工作流模式

### 1. 顺序执行

```go
result1 := workflow.ExecuteActivity(ctx, Activity1, input1)
result2 := workflow.ExecuteActivity(ctx, Activity2, result1)
```

### 2. 并行执行

```go
future1 := workflow.ExecuteActivity(ctx, Activity1, input1)
future2 := workflow.ExecuteActivity(ctx, Activity2, input2)
workflow.Await(future1, future2)
```

### 3. 条件执行

```go
if condition {
    workflow.ExecuteActivity(ctx, Activity1, input)
} else {
    workflow.ExecuteActivity(ctx, Activity2, input)
}
```

### 4. 循环执行

```go
for i := 0; i < 10; i++ {
    workflow.ExecuteActivity(ctx, Activity, i)
}
```

## 错误处理

### 重试策略

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

### 错误传播

工作流中的错误会自动传播，可以通过返回值或错误处理。

## 信号和查询

### 信号（Signal）

用于从外部向工作流发送异步消息：

```go
workflow.GetSignalChannel(ctx, "signal-name")
```

### 查询（Query）

用于查询工作流状态：

```go
workflow.SetQueryHandler(ctx, "query-name", handler)
```

## 持久化

Temporal 自动持久化工作流状态，包括：

- 工作流历史
- 活动结果
- 定时器
- 信号和查询

## 可扩展性

### 水平扩展

- 多个 Worker 可以处理同一个 Task Queue
- Temporal Server 支持集群部署

### 性能优化

- 使用合适的 Task Queue 分区
- 优化活动执行时间
- 合理设置超时和重试策略

## 监控和调试

### Temporal UI

- 查看工作流执行历史
- 调试工作流问题
- 监控工作流性能

### 日志

- 工作流和活动的日志
- 集成 OpenTelemetry 追踪

## 最佳实践

1. **保持工作流确定性** - 不要使用随机数、时间等非确定性操作
2. **合理设置超时** - 为活动设置合适的超时时间
3. **使用版本控制** - 工作流变更时使用版本控制
4. **错误处理** - 正确处理和传播错误
5. **幂等性** - 确保活动是幂等的
