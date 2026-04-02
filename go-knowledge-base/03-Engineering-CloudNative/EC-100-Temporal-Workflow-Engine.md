# EC-100: Temporal 工作流引擎深度分析 (Temporal Workflow Engine Deep Dive)

> **维度**: Engineering CloudNative
> **级别**: S (25+ KB)
> **标签**: #temporal #workflow-engine #durable-execution #stateful
> **相关**: EC-099, EC-112, FT-018

---

## 整合说明

本文档合并了：

- `58-Cadence-Temporal-Workflow-Engine.md` (19 KB)
- `69-Temporal-Workflow-Engine.md` (22 KB)
- `115-Task-Temporal-Workflow-Deep-Dive.md` (14 KB)

---

## 核心架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Temporal Architecture                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Client                    Server                     Workers            │
│  ──────                    ──────                     ───────            │
│                                                                          │
│  ┌─────────────┐          ┌──────────────┐          ┌─────────────┐     │
│  │ Temporal SDK│◄────────►│ Frontend     │◄────────►│ Worker      │     │
│  │ (Go/Java/   │  gRPC    │ Service      │  Poll    │ Process     │     │
│  │  TypeScript)│          │              │          │             │     │
│  └─────────────┘          └──────┬───────┘          └─────────────┘     │
│                                  │                                       │
│                                  ▼                                       │
│                          ┌──────────────┐                               │
│                          │ Matching     │                               │
│                          │ Service      │  任务路由                      │
│                          └──────┬───────┘                               │
│                                  │                                       │
│                    ┌─────────────┼─────────────┐                        │
│                    ▼             ▼             ▼                        │
│             ┌──────────┐ ┌──────────┐ ┌──────────┐                     │
│             │ History  │ │  Shard   │ │ Visibility│                    │
│             │ Service  │ │ Manager  │ │ Store     │                    │
│             └────┬─────┘ └────┬─────┘ └────┬─────┘                    │
│                  │            │            │                           │
│                  ▼            ▼            ▼                           │
│             ┌─────────────────────────────────┐                        │
│             │        Persistence              │                        │
│             │  (Cassandra/MySQL/PostgreSQL)   │                        │
│             └─────────────────────────────────┘                        │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Cadence vs Temporal

| 特性 | Cadence (Uber) | Temporal (独立) |
|------|---------------|-----------------|
| 维护 | Uber | Temporal Technologies |
| 协议 | Thrift | gRPC |
| 云托管 | 无 | Temporal Cloud |
| 特性更新 | 慢 | 快 |
| 社区 | 较小 | 较大 |

---

## 工作流定义

```go
func TaskWorkflow(ctx workflow.Context, task TaskInput) (TaskResult, error) {
 options := workflow.ActivityOptions{
  StartToCloseTimeout: 30 * time.Second,
  RetryPolicy: &temporal.RetryPolicy{
   InitialInterval:    time.Second,
   BackoffCoefficient: 2.0,
   MaximumInterval:    time.Minute,
   MaximumAttempts:    3,
  },
 }
 ctx = workflow.WithActivityOptions(ctx, options)

 // 1. 预处理
 var preprocessResult PreprocessResult
 if err := workflow.ExecuteActivity(ctx, PreprocessActivity, task).Get(ctx, &preprocessResult); err != nil {
  return TaskResult{}, err
 }

 // 2. 并行执行子任务
 selector := workflow.NewSelector(ctx)
 results := make([]SubTaskResult, len(task.SubTasks))

 for i, subTask := range task.SubTasks {
  i, subTask := i, subTask
  f := workflow.ExecuteActivity(ctx, ExecuteSubTaskActivity, subTask)
  selector.AddFuture(f, func(f workflow.Future) {
   f.Get(ctx, &results[i])
  })
 }

 for i := 0; i < len(task.SubTasks); i++ {
  selector.Select(ctx)
 }

 // 3. 聚合结果
 var aggregateResult AggregateResult
 if err := workflow.ExecuteActivity(ctx, AggregateActivity, results).Get(ctx, &aggregateResult); err != nil {
  return TaskResult{}, err
 }

 return TaskResult{Output: aggregateResult.Output}, nil
}
```

---

## 状态持久化

```go
// WorkflowState 工作流状态
type WorkflowState struct {
 RunID        string
 WorkflowType string
 Status       WorkflowStatus
 History      []*HistoryEvent
 MutableState *MutableState
}

type MutableState struct {
 NextEventID           int64
 LastProcessedEvent    int64
 PendingActivities     map[string]*PendingActivityInfo
 PendingTimers         map[string]*PendingTimerInfo
 PendingChildWorkflows map[string]*ChildWorkflowInfo
}
```

---

## 关键设计

| 特性 | 实现 | 优势 |
|------|------|------|
| 持久化执行 | History + Mutable State | 容错、可审计 |
| 确定性重放 | 捕获非确定性操作 | 状态恢复一致性 |
| 异步完成 | 命令模式 | 高吞吐、低延迟 |
| 版本控制 | 工作流类型版本 | 安全升级 |
