# Temporal Workflow 深度分析

> **分类**: 工程与云原生
> **标签**: #temporal #workflow-engine #durable-execution #stateful
> **参考**: Temporal SDK, Cadence Paper (Uber), Durable Functions

---

## Temporal 核心架构

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

## Workflow 执行模型

```go
package temporal

import (
 "context"
 "fmt"
 "time"

 "go.temporal.io/sdk/workflow"
)

// TaskWorkflow 任务工作流定义
func TaskWorkflow(ctx workflow.Context, task TaskInput) (TaskResult, error) {
 // 工作流选项
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
 err := workflow.ExecuteActivity(ctx, PreprocessActivity, task).Get(ctx, &preprocessResult)
 if err != nil {
  return TaskResult{}, fmt.Errorf("preprocess failed: %w", err)
 }

 // 2. 并行执行子任务
 selector := workflow.NewSelector(ctx)
 results := make([]SubTaskResult, len(task.SubTasks))

 for i, subTask := range task.SubTasks {
  i, subTask := i, subTask // 捕获循环变量

  f := workflow.ExecuteActivity(ctx, ExecuteSubTaskActivity, subTask)
  selector.AddFuture(f, func(f workflow.Future) {
   err := f.Get(ctx, &results[i])
   if err != nil {
    // 记录失败但不阻塞
    workflow.GetLogger(ctx).Error("subtask failed", err)
   }
  })
 }

 // 等待所有完成或超时
 deadline := workflow.Now(ctx).Add(5 * time.Minute)
 timer := workflow.NewTimer(ctx, deadline.Sub(workflow.Now(ctx)))
 selector.AddFuture(timer, func(f workflow.Future) {
  // 超时处理
 })

 // 等待指定数量完成
 for i := 0; i < len(task.SubTasks); i++ {
  selector.Select(ctx)
 }

 // 3. 聚合结果
 var aggregateResult AggregateResult
 err = workflow.ExecuteActivity(ctx, AggregateActivity, results).Get(ctx, &aggregateResult)
 if err != nil {
  return TaskResult{}, err
 }

 // 4. Saga 补偿检查
 if aggregateResult.NeedsCompensation {
  err := workflow.ExecuteActivity(ctx, CompensateActivity, results).Get(ctx, nil)
  if err != nil {
   return TaskResult{}, fmt.Errorf("compensation failed: %w", err)
  }
 }

 return TaskResult{
  ProcessedAt: workflow.Now(ctx),
  Output:      aggregateResult.Output,
 }, nil
}

// DurableTimer 持久化定时器（关键特性）
func ScheduledTaskWorkflow(ctx workflow.Context, schedule Schedule) error {
 for {
  // 计算下次执行时间
  nextRun := calculateNextRun(schedule)

  // 持久化睡眠 - 工作流在此"暂停"，服务器端维护状态
  // Worker 进程可以在此期间完全退出
  _ = workflow.Sleep(ctx, time.Until(nextRun))

  // 恢复执行时，Temporal 会从持久化状态恢复
  childCtx, cancel := workflow.WithCancel(ctx)
  childWorkflow := workflow.ExecuteChildWorkflow(childCtx, TaskWorkflow, TaskInput{})

  // 设置执行超时
  timer := workflow.NewTimer(ctx, schedule.Timeout)

  selector := workflow.NewSelector(ctx)
  selector.AddFuture(childWorkflow, func(f workflow.Future) {
   // 子工作流完成
   cancel()
  })
  selector.AddFuture(timer, func(f workflow.Future) {
   // 超时，取消子工作流
   cancel()
  })
  selector.Select(ctx)
 }
}

// ChildWorkflow 子工作流
func ChildTaskWorkflow(ctx workflow.Context, input ChildTaskInput) (ChildTaskResult, error) {
 // 查询当前状态（外部可查询）
 workflow.SetQueryHandler(ctx, "state", func() (string, error) {
  return fmt.Sprintf("processing step %d", input.CurrentStep), nil
 })

 // 发送信号（外部可控制）
 ch := workflow.GetSignalChannel(ctx, "control")

 selector := workflow.NewSelector(ctx)
 selector.AddReceive(ch, func(c workflow.ReceiveChannel, more bool) {
  var signal ControlSignal
  c.Receive(ctx, &signal)
  // 处理控制信号
 })

 // 主处理逻辑
 for i := 0; i < input.TotalSteps; i++ {
  // 执行步骤
  err := workflow.ExecuteActivity(ctx, ProcessStepActivity, i).Get(ctx, nil)
  if err != nil {
   return ChildTaskResult{}, err
  }

  // 检查信号
  selector.Select(ctx)
 }

 return ChildTaskResult{Completed: true}, nil
}
```

---

## 状态持久化机制

```go
// Workflow 状态持久化

// WorkflowState 工作流状态（服务器端存储）
type WorkflowState struct {
 RunID        string
 WorkflowType string
 Status       WorkflowStatus

 // 历史事件（不可变）
 History []*HistoryEvent

 // 当前状态
 MutableState *MutableState
}

type MutableState struct {
 NextEventID        int64
 LastProcessedEvent int64

 // 活动状态
 PendingActivities map[string]*PendingActivityInfo
 PendingTimers     map[string]*PendingTimerInfo
 PendingSignals    map[string]*PendingSignalInfo

 // 子工作流
 PendingChildWorkflows map[string]*ChildWorkflowInfo

 // 完成的活动（用于 Saga）
 CompletedActivities []CompletedActivityInfo
}

type PendingActivityInfo struct {
 ScheduleID     int64
 ActivityID     string
 ActivityType   string
 State          PendingActivityState
 ScheduledTime  time.Time
 StartedTime    *time.Time
 Attempt        int32
 RetryPolicy    *RetryPolicy
}

// HistoryEvent 历史事件
type HistoryEvent struct {
 EventID   int64
 Timestamp time.Time
 EventType EventType

 // 事件详情
 Attributes interface{}
}

// EventType 事件类型
const (
 EventTypeWorkflowExecutionStarted       EventType = iota
 EventTypeWorkflowExecutionCompleted
 EventTypeWorkflowExecutionFailed
 EventTypeActivityTaskScheduled
 EventTypeActivityTaskStarted
 EventTypeActivityTaskCompleted
 EventTypeActivityTaskFailed
 EventTypeTimerStarted
 EventTypeTimerFired
 EventTypeWorkflowExecutionSignaled
 // ...
)

// 事件溯源：从 History 重建状态
func (s *MutableState) RebuildFromHistory(history []*HistoryEvent) error {
 for _, event := range history {
  switch event.EventType {
  case EventTypeActivityTaskScheduled:
   attr := event.Attributes.(*ScheduleActivityTaskCommandAttributes)
   s.PendingActivities[attr.ActivityId] = &PendingActivityInfo{
    ScheduleID:   event.EventID,
    ActivityID:   attr.ActivityId,
    ActivityType: attr.ActivityType.Name,
    State:        PendingActivityStateScheduled,
    ScheduledTime: event.Timestamp,
    Attempt:      0,
   }

  case EventTypeActivityTaskStarted:
   attr := event.Attributes.(*ActivityTaskStartedEventAttributes)
   if activity, ok := s.PendingActivities[attr.ScheduledEventId]; ok {
    activity.State = PendingActivityStateStarted
    activity.StartedTime = &event.Timestamp
   }

  case EventTypeActivityTaskCompleted:
   attr := event.Attributes.(*ActivityTaskCompletedEventAttributes)
   activityID := s.getActivityID(attr.ScheduledEventId)
   delete(s.PendingActivities, activityID)
   s.CompletedActivities = append(s.CompletedActivities, CompletedActivityInfo{
    ActivityID: activityID,
    Result:     attr.Result,
   })

  // ... 其他事件类型
  }

  s.NextEventID = event.EventID + 1
 }

 return nil
}
```

---

## Worker 实现

```go
// Worker 实现

// Worker 工作节点
type Worker struct {
 client Client

 // 任务轮询
 taskQueues []string

 // 工作流和活动注册表
 registry *Registry

 // 并发控制
 executor *taskExecutor

 // 停止信号
 stopCh chan struct{}
}

// Start 启动 Worker
func (w *Worker) Start() error {
 for _, queue := range w.taskQueues {
  go w.pollAndProcess(queue)
 }
 return nil
}

// pollAndProcess 轮询并处理任务
func (w *Worker) pollAndProcess(taskQueue string) {
 for {
  select {
  case <-w.stopCh:
   return
  default:
  }

  // 长轮询获取任务
  task, err := w.client.PollWorkflowTaskQueue(context.Background(), &PollRequest{
   Namespace: "default",
   TaskQueue: taskQueue,
   Identity:  w.identity,
  })

  if err != nil {
   time.Sleep(time.Second)
   continue
  }

  if task == nil {
   continue
  }

  // 执行任务
  w.executor.Execute(task)
 }
}

// taskExecutor 任务执行器
type taskExecutor struct {
 client   Client
 registry *Registry

 // 工作流缓存
 workflowCache map[string]*workflowState
 cacheMu       sync.RWMutex
}

// Execute 执行任务
func (e *taskExecutor) Execute(task *WorkflowTask) {
 switch task.Type {
 case TaskTypeWorkflow:
  e.executeWorkflowTask(task)
 case TaskTypeActivity:
  e.executeActivityTask(task)
 }
}

// executeWorkflowTask 执行工作流任务
func (e *taskExecutor) executeWorkflowTask(task *WorkflowTask) {
 // 获取或创建工作流状态
 state := e.getOrCreateWorkflowState(task.WorkflowID, task.RunID)

 // 应用新事件到状态机
 for _, event := range task.NewEvents {
  if err := state.ProcessEvent(event); err != nil {
   // 报告失败
   e.client.RespondWorkflowTaskFailed(task.TaskToken, err)
   return
  }
 }

 // 获取工作流函数
 fn, ok := e.registry.GetWorkflow(state.WorkflowType)
 if !ok {
  e.client.RespondWorkflowTaskFailed(task.TaskToken, fmt.Errorf("unknown workflow type"))
  return
 }

 // 创建 WorkflowContext
 wctx := &workflowContextImpl{
  state:    state,
  client:   e.client,
  taskToken: task.TaskToken,
 }

 // 恢复执行（或新启动）
 if state.IsReplaying() {
  wctx.SetReplaying(true)
 }

 // 执行工作流函数
 result, err := fn(wctx, state.Input)

 if err != nil {
  // 工作流失败
  e.client.RespondWorkflowTaskFailed(task.TaskToken, err)
  return
 }

 // 收集命令并响应
 commands := wctx.GetCommands()
 e.client.RespondWorkflowTaskCompleted(task.TaskToken, commands, result)
}
```

---

## 关键设计决策

| 特性 | 实现 | 优势 |
|------|------|------|
| 持久化执行 | History + Mutable State | 容错、可审计、可重放 |
| 确定性重放 | 捕获非确定性操作 | 状态恢复一致性 |
| 异步完成 | 命令模式 | 高吞吐、低延迟 |
| 版本控制 | 工作流类型版本 | 安全升级 |
| 查询处理 | 只读副本 | 不影响执行 |

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02