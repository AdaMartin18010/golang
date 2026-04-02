# Temporal 工作流引擎架构与实现

> **分类**: 工程与云原生
> **标签**: #temporal #workflow #cadence #state-machine
> **参考**: Temporal SDK, Cadence Architecture Papers

---

## Temporal 架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Temporal System Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Frontend Service                              │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  gRPC API   │  │  Namespace  │  │   Rate      │  │   Auth      │ │   │
│  │  │             │  │   Router    │  │   Limit     │  │             │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                     Matching Service (Task Queue)                      │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │   │
│  │  │  Workflow Task  │  │  Activity Task  │  │  Worker Poll    │      │   │
│  │  │     Queue       │  │     Queue       │  │    Dispatcher   │      │   │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────┘      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    History Service (Event Sourcing)                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  Event Log  │  │  Command    │  │  Workflow   │  │  State      │ │   │
│  │  │  (Append)   │  │  Processing │  │  State      │  │  Rebuild    │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Worker Service (Go SDK)                             │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  Workflow   │  │  Activity   │  │  Interceptor│  │   Logger    │ │   │
│  │  │  Executor   │  │  Executor   │  │  (Metrics)  │  │  (Context)  │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Persistence Layer                                   │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  Cassandra  │  │  PostgreSQL │  │  MySQL      │  │  Elasticsearch│ │   │
│  │  │  (Events)   │  │  (Metadata) │  │  (Tasks)    │  │   (Visibility)│ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 事件溯源模型

```go
package temporal

import (
    "time"

    commonpb "go.temporal.io/api/common/v1"
    historypb "go.temporal.io/api/history/v1"
)

// Event 事件定义
type Event struct {
    EventID   int64
    Timestamp time.Time
    EventType EventType
    Attributes interface{}
}

type EventType int

const (
    EventTypeWorkflowExecutionStarted EventType = iota
    EventTypeWorkflowExecutionCompleted
    EventTypeWorkflowExecutionFailed
    EventTypeWorkflowExecutionTimedOut
    EventTypeActivityTaskScheduled
    EventTypeActivityTaskStarted
    EventTypeActivityTaskCompleted
    EventTypeActivityTaskFailed
    EventTypeTimerStarted
    EventTypeTimerFired
    EventTypeWorkflowExecutionSignaled
    EventTypeWorkflowExecutionCanceled
)

// Command 命令定义（工作流生成的待执行操作）
type Command struct {
    CommandType CommandType
    Attributes  interface{}
}

type CommandType int

const (
    CommandTypeScheduleActivityTask CommandType = iota
    CommandTypeStartTimer
    CommandTypeCompleteWorkflowExecution
    CommandTypeFailWorkflowExecution
    CommandTypeRequestCancelActivityTask
)

// WorkflowState 工作流状态
type WorkflowState struct {
    WorkflowID string
    RunID      string
    Status     WorkflowStatus

    // 执行状态
    PendingActivities map[string]*PendingActivity
    PendingTimers     map[string]*PendingTimer

    // 事件历史
    Events []*historypb.HistoryEvent

    // 当前 WFT (Workflow Task) 状态
    NextEventID int64
}

type WorkflowStatus int

const (
    WorkflowStatusRunning WorkflowStatus = iota
    WorkflowStatusCompleted
    WorkflowStatusFailed
    WorkflowStatusCanceled
    WorkflowStatusTimedOut
)

type PendingActivity struct {
    ActivityID   string
    ActivityType string
    Input        *commonpb.Payloads
    ScheduledTime time.Time
    StartedTime   *time.Time
    Attempt       int32
}

type PendingTimer struct {
    TimerID   string
    FireTime  time.Time
    StartedID int64
}

// RebuildWorkflowState 从历史事件重建工作流状态
func RebuildWorkflowState(events []*historypb.HistoryEvent) (*WorkflowState, error) {
    state := &WorkflowState{
        Status:            WorkflowStatusRunning,
        PendingActivities: make(map[string]*PendingActivity),
        PendingTimers:     make(map[string]*PendingTimer),
        Events:            events,
    }

    for _, event := range events {
        switch event.EventType {
        case historypb.EVENT_TYPE_WORKFLOW_EXECUTION_STARTED:
            attrs := event.GetWorkflowExecutionStartedEventAttributes()
            state.WorkflowID = event.WorkflowExecution.GetWorkflowId()
            state.RunID = event.WorkflowExecution.GetRunId()

        case historypb.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
            attrs := event.GetActivityTaskScheduledEventAttributes()
            activityID := attrs.GetActivityId()
            state.PendingActivities[activityID] = &PendingActivity{
                ActivityID:    activityID,
                ActivityType:  attrs.GetActivityType().GetName(),
                Input:         attrs.GetInput(),
                ScheduledTime: event.EventTime.AsTime(),
                Attempt:       1,
            }

        case historypb.EVENT_TYPE_ACTIVITY_TASK_STARTED:
            attrs := event.GetActivityTaskStartedEventAttributes()
            activityID := getActivityIDFromScheduledEventID(state.Events, attrs.GetScheduledEventId())
            if pending, ok := state.PendingActivities[activityID]; ok {
                now := event.EventTime.AsTime()
                pending.StartedTime = &now
            }

        case historypb.EVENT_TYPE_ACTIVITY_TASK_COMPLETED:
            attrs := event.GetActivityTaskCompletedEventAttributes()
            activityID := getActivityIDFromScheduledEventID(state.Events, attrs.GetScheduledEventId())
            delete(state.PendingActivities, activityID)

        case historypb.EVENT_TYPE_ACTIVITY_TASK_FAILED:
            attrs := event.GetActivityTaskFailedEventAttributes()
            activityID := getActivityIDFromScheduledEventID(state.Events, attrs.GetScheduledEventId())
            if pending, ok := state.PendingActivities[activityID]; ok {
                pending.Attempt++
                pending.StartedTime = nil
            }

        case historypb.EVENT_TYPE_TIMER_STARTED:
            attrs := event.GetTimerStartedEventAttributes()
            timerID := attrs.GetTimerId()
            state.PendingTimers[timerID] = &PendingTimer{
                TimerID:   timerID,
                FireTime:  event.EventTime.AsTime().Add(attrs.GetStartToFireTimeout().AsDuration()),
                StartedID: event.EventId,
            }

        case historypb.EVENT_TYPE_TIMER_FIRED:
            attrs := event.GetTimerFiredEventAttributes()
            delete(state.PendingTimers, attrs.GetTimerId())

        case historypb.EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED:
            state.Status = WorkflowStatusCompleted

        case historypb.EVENT_TYPE_WORKFLOW_EXECUTION_FAILED:
            state.Status = WorkflowStatusFailed

        case historypb.EVENT_TYPE_WORKFLOW_EXECUTION_TIMED_OUT:
            state.Status = WorkflowStatusTimedOut
        }

        state.NextEventID = event.EventId + 1
    }

    return state, nil
}
```

---

## Go SDK 工作流实现

```go
package temporal

import (
    "context"
    "fmt"
    "time"

    "go.temporal.io/sdk/workflow"
)

// OrderWorkflow 订单处理工作流示例
func OrderWorkflow(ctx workflow.Context, orderID string) error {
    // 1. 设置工作流选项
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Minute,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval:    time.Minute,
            MaximumAttempts:    3,
            NonRetryableErrorTypes: []string{"InvalidOrderError"},
        },
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 2. 验证订单
    var validationResult ValidationResult
    err := workflow.ExecuteActivity(ctx, ValidateOrderActivity, orderID).Get(ctx, &validationResult)
    if err != nil {
        return fmt.Errorf("order validation failed: %w", err)
    }

    // 3. 扣减库存（ Saga 补偿模式）
    compensations := []func() error{}

    var inventoryResult InventoryResult
    err = workflow.ExecuteActivity(ctx, DeductInventoryActivity, orderID).Get(ctx, &inventoryResult)
    if err != nil {
        return fmt.Errorf("inventory deduction failed: %w", err)
    }
    compensations = append(compensations, func() error {
        return workflow.ExecuteActivity(ctx, RestoreInventoryActivity, orderID).Get(ctx, nil)
    })

    // 4. 处理支付
    var paymentResult PaymentResult
    err = workflow.ExecuteActivity(ctx, ProcessPaymentActivity, orderID).Get(ctx, &paymentResult)
    if err != nil {
        // 执行补偿
        for i := len(compensations) - 1; i >= 0; i-- {
            _ = compensations[i]()
        }
        return fmt.Errorf("payment failed: %w", err)
    }
    compensations = append(compensations, func() error {
        return workflow.ExecuteActivity(ctx, RefundPaymentActivity, orderID).Get(ctx, nil)
    })

    // 5. 发货
    var shippingResult ShippingResult
    err = workflow.ExecuteActivity(ctx, ShipOrderActivity, orderID).Get(ctx, &shippingResult)
    if err != nil {
        // 执行补偿
        for i := len(compensations) - 1; i >= 0; i-- {
            _ = compensations[i]()
        }
        return fmt.Errorf("shipping failed: %w", err)
    }

    // 6. 发送通知
    _ = workflow.ExecuteActivity(ctx, SendNotificationActivity, orderID).Get(ctx, nil)

    return nil
}

// Activity 实现
func ValidateOrderActivity(ctx context.Context, orderID string) (*ValidationResult, error) {
    // 实际业务逻辑
    return &ValidationResult{Valid: true}, nil
}

func DeductInventoryActivity(ctx context.Context, orderID string) (*InventoryResult, error) {
    return &InventoryResult{Deducted: true}, nil
}

func ProcessPaymentActivity(ctx context.Context, orderID string) (*PaymentResult, error) {
    return &PaymentResult{Success: true}, nil
}

func ShipOrderActivity(ctx context.Context, orderID string) (*ShippingResult, error) {
    return &ShippingResult{Shipped: true}, nil
}

// Saga 模式实现
type Saga struct {
    compensations []func() error
}

func NewSaga() *Saga {
    return &Saga{}
}

func (s *Saga) Add(compensation func() error) {
    s.compensations = append(s.compensations, compensation)
}

func (s *Saga) Compensate() error {
    var errs []error
    // 逆序执行补偿
    for i := len(s.compensations) - 1; i >= 0; i-- {
        if err := s.compensations[i](); err != nil {
            errs = append(errs, err)
        }
    }
    if len(errs) > 0 {
        return fmt.Errorf("compensation failed: %v", errs)
    }
    return nil
}

// 使用 Saga 的工作流
func OrderWorkflowWithSaga(ctx workflow.Context, orderID string) error {
    saga := NewSaga()
    defer func() {
        if err := recover(); err != nil {
            saga.Compensate()
            panic(err)
        }
    }()

    // 每个操作后添加补偿
    _ = workflow.ExecuteActivity(ctx, DeductInventoryActivity, orderID).Get(ctx, nil)
    saga.Add(func() error {
        return workflow.ExecuteActivity(ctx, RestoreInventoryActivity, orderID).Get(ctx, nil)
    })

    _ = workflow.ExecuteActivity(ctx, ProcessPaymentActivity, orderID).Get(ctx, nil)
    saga.Add(func() error {
        return workflow.ExecuteActivity(ctx, RefundPaymentActivity, orderID).Get(ctx, nil)
    })

    return nil
}
```

---

## 定时任务与调度

```go
package temporal

import (
    "time"

    "go.temporal.io/sdk/workflow"
)

// CronWorkflow Cron 定时工作流
func CronWorkflow(ctx workflow.Context, schedule string) error {
    // 设置定时器选项
    options := workflow.ChildWorkflowOptions{
        WorkflowExecutionTimeout: 5 * time.Minute,
    }
    childCtx := workflow.WithChildOptions(ctx, options)

    // 启动定时子工作流
    cronWorkflow := func(ctx workflow.Context) error {
        // 实际执行的任务
        return workflow.ExecuteActivity(ctx, ScheduledTaskActivity).Get(ctx, nil)
    }

    // 使用 CronSchedule 选项
    cwo := workflow.ChildWorkflowOptions{
        CronSchedule: schedule, // "0 9 * * *" 每天9点
    }
    childCtx = workflow.WithChildOptions(childCtx, cwo)

    future := workflow.ExecuteChildWorkflow(childCtx, cronWorkflow)

    var result interface{}
    if err := future.Get(childCtx, &result); err != nil {
        return err
    }

    return nil
}

// 定时触发器模式
func ScheduledWorkflow(ctx workflow.Context) error {
    // 每小时的第0分钟执行
    timer := workflow.NewTimer(ctx, getNextHour())

    for {
        selector := workflow.NewSelector(ctx)
        selector.AddFuture(timer, func(f workflow.Future) {
            // 定时器触发
            _ = workflow.ExecuteActivity(ctx, HourlyTaskActivity).Get(ctx, nil)

            // 设置下一个定时器
            timer = workflow.NewTimer(ctx, getNextHour())
        })

        // 等待信号或定时器
        signalChan := workflow.GetSignalChannel(ctx, "update-schedule")
        selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
            var newSchedule ScheduleUpdate
            c.Receive(ctx, &newSchedule)
            // 更新调度逻辑
        })

        selector.Select(ctx)

        // 检查是否应该继续
        if shouldStop(ctx) {
            break
        }
    }

    return nil
}

func getNextHour() time.Duration {
    now := time.Now()
    next := now.Truncate(time.Hour).Add(time.Hour)
    return next.Sub(now)
}

// 延迟任务模式
func DelayedTaskWorkflow(ctx workflow.Context, taskID string, delay time.Duration) error {
    // 等待指定时间
    _ = workflow.Sleep(ctx, delay)

    // 执行任务
    return workflow.ExecuteActivity(ctx, ExecuteDelayedTask, taskID).Get(ctx, nil)
}
```

---

## 外部信号与查询

```go
package temporal

import (
    "go.temporal.io/sdk/workflow"
)

// OrderState 订单状态
type OrderState struct {
    Status      string
    Items       []OrderItem
    TotalAmount float64
    Shipping    ShippingInfo
}

type OrderItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

type ShippingInfo struct {
    Address string
    Status  string
}

// OrderProcessingWorkflow 带信号处理的工作流
func OrderProcessingWorkflow(ctx workflow.Context, orderID string) (*OrderState, error) {
    state := &OrderState{
        Status: "PENDING",
        Items:  []OrderItem{},
    }

    // 设置查询处理
    err := workflow.SetQueryHandler(ctx, "get-state", func() (*OrderState, error) {
        return state, nil
    })
    if err != nil {
        return nil, err
    }

    // 等待订单确认信号
    confirmChan := workflow.GetSignalChannel(ctx, "confirm-order")
    var confirmData ConfirmOrderSignal
    confirmChan.Receive(ctx, &confirmData)

    state.Status = "CONFIRMED"
    state.Items = confirmData.Items
    state.TotalAmount = calculateTotal(state.Items)

    // 等待支付信号
    paymentChan := workflow.GetSignalChannel(ctx, "payment-received")
    var paymentData PaymentSignal
    paymentChan.Receive(ctx, &paymentData)

    state.Status = "PAID"

    // 等待发货信号
    shippingChan := workflow.GetSignalChannel(ctx, "shipped")
    var shippingData ShippingSignal
    shippingChan.Receive(ctx, &shippingData)

    state.Shipping = shippingData.Info
    state.Status = "SHIPPED"

    // 设置最终查询处理器
    _ = workflow.SetQueryHandler(ctx, "get-state", func() (*OrderState, error) {
        return state, nil
    })

    return state, nil
}

// 信号数据结构
type ConfirmOrderSignal struct {
    Items []OrderItem
}

type PaymentSignal struct {
    TransactionID string
    Amount        float64
}

type ShippingSignal struct {
    Info ShippingInfo
}

// 发送信号示例
func SendSignalExample() {
    // c.SignalWorkflow(ctx, workflowID, runID, "confirm-order", signalData)
}

// 查询状态示例
func QueryStateExample() {
    // resp, err := c.QueryWorkflow(ctx, workflowID, runID, "get-state")
}
```

---

## Worker 实现

```go
package temporal

import (
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
)

// WorkerConfig Worker 配置
type WorkerConfig struct {
    // 任务队列名称
    TaskQueue string

    // 并发度
    MaxConcurrentActivityExecutionSize     int
    MaxConcurrentWorkflowTaskExecutionSize int

    // Worker 数量
    MaxConcurrentActivityTaskPollers     int
    MaxConcurrentWorkflowTaskPollers     int

    // 会话选项
    EnableSessionWorker bool
}

// CreateWorker 创建 Temporal Worker
func CreateWorker(c client.Client, cfg WorkerConfig) worker.Worker {
    w := worker.New(c, cfg.TaskQueue, worker.Options{
        MaxConcurrentActivityExecutionSize:     cfg.MaxConcurrentActivityExecutionSize,
        MaxConcurrentWorkflowTaskExecutionSize: cfg.MaxConcurrentWorkflowTaskExecutionSize,
        MaxConcurrentActivityTaskPollers:       cfg.MaxConcurrentActivityTaskPollers,
        MaxConcurrentWorkflowTaskPollers:       cfg.MaxConcurrentWorkflowTaskPollers,

        // 拦截器
        Interceptors: []worker.Interceptor{
            &MetricsInterceptor{},
            &LoggingInterceptor{},
        },
    })

    // 注册工作流
    w.RegisterWorkflow(OrderWorkflow)
    w.RegisterWorkflow(CronWorkflow)
    w.RegisterWorkflow(OrderProcessingWorkflow)

    // 注册 Activities
    w.RegisterActivity(ValidateOrderActivity)
    w.RegisterActivity(DeductInventoryActivity)
    w.RegisterActivity(ProcessPaymentActivity)
    w.RegisterActivity(ShipOrderActivity)

    return w
}

// MetricsInterceptor 指标拦截器
type MetricsInterceptor struct{}

func (m *MetricsInterceptor) InterceptActivity(
    ctx context.Context,
    next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
    return &activityMetricsInterceptor{next: next}
}

type activityMetricsInterceptor struct {
    next interceptor.ActivityInboundInterceptor
}

func (a *activityMetricsInterceptor) ExecuteActivity(
    ctx context.Context,
    in *interceptor.ExecuteActivityInput,
) (interface{}, error) {
    start := time.Now()

    result, err := a.next.ExecuteActivity(ctx, in)

    duration := time.Since(start)
    activityType := in.ActivityType

    // 记录指标
    if err != nil {
        activityFailureCounter.WithLabelValues(activityType).Inc()
    } else {
        activitySuccessCounter.WithLabelValues(activityType).Inc()
    }
    activityDurationHistogram.WithLabelValues(activityType).Observe(duration.Seconds())

    return result, err
}
```
