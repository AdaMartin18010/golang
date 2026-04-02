# Cadence/Temporal 工作流引擎深度解析 (Cadence/Temporal Workflow Engine Deep Dive)

> **分类**: 工程与云原生
> **标签**: #cadence #temporal #workflow-engine #saga-pattern
> **参考**: Uber Cadence, Temporal.io, AWS SWF

---

## 架构核心概念

Cadence/Temporal 是一种用于构建可容错、可扩展的长时间运行工作流的编程框架。它将工作流编排逻辑与业务活动分离。

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Temporal/Cadence Cluster                       │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐  ┌───────────┐ │
│  │   Frontend    │  │    History    │  │   Matching    │  │  Worker   │ │
│  │   (Gateway)   │  │   (Event Sourcing)│  │   (Task Queue)  │  │ (System)  │ │
│  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘  └─────┬─────┘ │
│          │                  │                  │                │       │
│          └──────────────────┴──────────────────┘                │       │
│                             │                                   │       │
│  ┌──────────────────────────┴───────────────────────────────────┘       │
│  │                         Persistence Layer                             │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                    │
│  │  │  Cassandra  │  │  PostgreSQL │  │    MySQL    │                    │
│  │  │  (Events)   │  │  (Metadata) │  │  (Visibility)│                   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                    │
│  └──────────────────────────────────────────────────────────────────────┘
└─────────────────────────────────────────────────────────────────────────┘
                                    │
              ┌─────────────────────┼─────────────────────┐
              ▼                     ▼                     ▼
        ┌───────────┐        ┌───────────┐        ┌───────────┐
        │  Worker 1 │        │  Worker 2 │        │  Worker N │
        │           │        │           │        │           │
        │ • Workflow│        │ • Workflow│        │ • Workflow│
        │ • Activity│        │ • Activity│        │ • Activity│
        │ • Local   │        │ • Local   │        │ • Local   │
        └───────────┘        └───────────┘        └───────────┘
```

---

## 工作流实现原理

```go
// Workflow 函数实现约束
// 1. 必须接受 workflow.Context 作为第一个参数
// 2. 必须返回 error
// 3. 使用 framework 提供的 API 进行异步操作

package workflows

import (
    "go.temporal.io/sdk/workflow"
    "go.temporal.io/sdk/temporal"
    "time"
)

// OrderWorkflow 订单处理工作流示例
func OrderWorkflow(ctx workflow.Context, orderID string, items []OrderItem) (string, error) {
    // 1. 设置 Activity 选项
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Minute,
        // 重试策略
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval:    time.Minute,
            MaximumAttempts:    3,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 2. 记录工作流开始
    workflow.GetLogger(ctx).Info("Order workflow started", "orderID", orderID)

    // 3. 执行库存检查 Activity
    var inventoryResult InventoryResult
    err := workflow.ExecuteActivity(ctx, activities.CheckInventory, items).Get(ctx, &inventoryResult)
    if err != nil {
        return "", err
    }

    if !inventoryResult.Available {
        return "", fmt.Errorf("inventory not available")
    }

    // 4. 执行支付 Activity
    var paymentResult PaymentResult
    err = workflow.ExecuteActivity(ctx, activities.ProcessPayment, orderID, calculateTotal(items)).Get(ctx, &paymentResult)
    if err != nil {
        // 支付失败，回滚库存
        _ = workflow.ExecuteActivity(ctx, activities.ReleaseInventory, items).Get(ctx, nil)
        return "", err
    }

    // 5. 执行发货 Activity
    var shippingResult ShippingResult
    err = workflow.ExecuteActivity(ctx, activities.CreateShipment, orderID, items, paymentResult.TransactionID).Get(ctx, &shippingResult)
    if err != nil {
        // 发货失败，触发补偿流程
        return "", handleShippingFailure(ctx, orderID, paymentResult.TransactionID)
    }

    // 6. 发送通知
    _ = workflow.ExecuteActivity(ctx, activities.SendNotification, orderID, "order_shipped").Get(ctx, nil)

    return shippingResult.TrackingID, nil
}

// handleShippingFailure 处理发货失败的补偿逻辑
func handleShippingFailure(ctx workflow.Context, orderID, transactionID string) error {
    // Saga 补偿模式
    // 1. 退款
    err := workflow.ExecuteActivity(ctx, activities.RefundPayment, transactionID).Get(ctx, nil)
    if err != nil {
        return fmt.Errorf("refund failed: %w", err)
    }

    // 2. 释放库存
    err = workflow.ExecuteActivity(ctx, activities.ReleaseInventory, orderID).Get(ctx, nil)
    if err != nil {
        return fmt.Errorf("release inventory failed: %w", err)
    }

    return fmt.Errorf("shipping failed, compensation completed")
}
```

---

## Activity 实现模式

```go
package activities

import (
    "context"
    "go.temporal.io/sdk/activity"
)

// Activity 结构体封装依赖
type OrderActivities struct {
    inventoryClient InventoryClient
    paymentClient   PaymentClient
    shippingClient  ShippingClient
    notificationSvc NotificationService
}

// CheckInventory 库存检查 Activity
func (a *OrderActivities) CheckInventory(ctx context.Context, items []OrderItem) (InventoryResult, error) {
    // 获取 Activity 信息
    info := activity.GetInfo(ctx)
    activity.GetLogger(ctx).Info("Checking inventory",
        "activityID", info.ActivityID,
        "attempt", info.Attempt,
    )

    // 心跳机制（长时间运行的 Activity）
    activity.RecordHeartbeat(ctx, "started")

    result, err := a.inventoryClient.CheckAvailability(items)
    if err != nil {
        return InventoryResult{}, err
    }

    activity.RecordHeartbeat(ctx, "completed")
    return result, nil
}

// ProcessPayment 支付处理 Activity
func (a *OrderActivities) ProcessPayment(ctx context.Context, orderID string, amount float64) (PaymentResult, error) {
    // 实现幂等性
    // 使用 orderID 作为幂等键
    result, err := a.paymentClient.Charge(PaymentRequest{
        OrderID:       orderID,
        Amount:        amount,
        IdempotencyKey: orderID,
    })

    return result, err
}

// 本地 Activity（在同进程执行，不经过 Task Queue）
func (a *OrderActivities) CalculateDiscount(ctx context.Context, items []OrderItem, customerTier string) (float64, error) {
    // 本地计算，不涉及外部服务
    var total float64
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }

    discount := calculateTierDiscount(customerTier)
    return total * discount, nil
}
```

---

## 事件驱动工作流模式

```go
// Signal 处理模式
// 工作流可以接收外部 Signal 并响应

func OrderWorkflowWithSignals(ctx workflow.Context, orderID string) error {
    // 状态机定义
    var state OrderState
    state.Status = "pending"

    // 创建 Signal Channel
    paymentSignalCh := workflow.GetSignalChannel(ctx, "payment-completed")
    cancelSignalCh := workflow.GetSignalChannel(ctx, "cancel-order")

    // Selector 用于等待多个事件
    selector := workflow.NewSelector(ctx)

    // 处理支付完成 Signal
    selector.AddReceive(paymentSignalCh, func(c workflow.ReceiveChannel, more bool) {
        var paymentInfo PaymentInfo
        c.Receive(ctx, &paymentInfo)

        state.Status = "paid"
        state.PaymentID = paymentInfo.TransactionID

        // 触发发货
        workflow.ExecuteActivity(ctx, activities.CreateShipment, orderID).Get(ctx, nil)
    })

    // 处理取消 Signal
    selector.AddReceive(cancelSignalCh, func(c workflow.ReceiveChannel, more bool) {
        var cancelInfo CancelInfo
        c.Receive(ctx, &cancelInfo)

        if state.Status == "paid" {
            // 已支付，需要退款
            workflow.ExecuteActivity(ctx, activities.RefundPayment, state.PaymentID).Get(ctx, nil)
        }

        state.Status = "cancelled"
    })

    // 设置超时
    timer := workflow.NewTimer(ctx, 24*time.Hour)
    selector.AddFuture(timer, func(f workflow.Future) {
        state.Status = "expired"
    })

    // 等待任一事件
    selector.Select(ctx)

    // 持久化最终状态
    _ = workflow.ExecuteActivity(ctx, activities.UpdateOrderState, orderID, state).Get(ctx, nil)

    return nil
}

// Query Handler 提供工作流状态查询
func init() {
    workflow.RegisterQueryHandler("get-order-state", func(state *OrderState) (*OrderState, error) {
        return state, nil
    })
}

// 外部发送 Signal
func SendPaymentSignal(c client.Client, workflowID string, paymentInfo PaymentInfo) error {
    return c.SignalWorkflow(context.Background(), workflowID, "", "payment-completed", paymentInfo)
}
```

---

## 长时间运行工作流模式

```go
// ContinueAsNew 模式处理无限长时间运行
func DataPipelineWorkflow(ctx workflow.Context, config PipelineConfig, lastCheckpoint time.Time) error {
    // 获取新数据批次
    var batches []DataBatch
    err := workflow.ExecuteActivity(ctx, activities.FetchDataBatches, lastCheckpoint).Get(ctx, &batches)
    if err != nil {
        return err
    }

    // 处理批次
    for _, batch := range batches {
        err := workflow.ExecuteActivity(ctx, activities.ProcessBatch, batch).Get(ctx, nil)
        if err != nil {
            return err
        }
        lastCheckpoint = batch.Timestamp
    }

    // 检查是否需要 ContinueAsNew（控制历史大小）
    if workflow.GetInfo(ctx).GetCurrentHistoryLength() > 10000 {
        // 启动新的工作流实例，携带状态
        return workflow.NewContinueAsNewError(ctx, DataPipelineWorkflow, config, lastCheckpoint)
    }

    // 等待下一个处理周期
    _ = workflow.Sleep(ctx, 5*time.Minute)

    // 递归调用自身
    return workflow.NewContinueAsNewError(ctx, DataPipelineWorkflow, config, lastCheckpoint)
}

// Child Workflow 模式
func ParentWorkflow(ctx workflow.Context, parentArgs ParentArgs) error {
    // 启动多个子工作流
    futures := make([]workflow.ChildWorkflowFuture, 0, len(parentArgs.SubTasks))

    for _, subTask := range parentArgs.SubTasks {
        cwo := workflow.ChildWorkflowOptions{
            WorkflowID: fmt.Sprintf("%s-child-%s", workflow.GetInfo(ctx).WorkflowExecution.ID, subTask.ID),
            // 父工作流取消时，子工作流继续运行
            ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
        }

        childCtx := workflow.WithChildOptions(ctx, cwo)
        future := workflow.ExecuteChildWorkflow(childCtx, ChildWorkflow, subTask)
        futures = append(futures, future)
    }

    // 等待所有子工作流完成
    for _, future := range futures {
        var result ChildResult
        if err := future.Get(ctx, &result); err != nil {
            // 处理子工作流错误
            workflow.GetLogger(ctx).Error("Child workflow failed", "error", err)
        }
    }

    return nil
}
```

---

## Worker 实现与配置

```go
package worker

import (
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    "go.temporal.io/sdk/workflow"
)

func StartWorker() {
    // 创建客户端
    c, err := client.NewClient(client.Options{
        HostPort: "localhost:7233",
        Namespace: "default",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 创建 Worker
    w := worker.New(c, "order-task-queue", worker.Options{
        // 并发控制
        MaxConcurrentActivityExecutionSize:     100,
        MaxConcurrentWorkflowTaskExecutionSize: 50,
        MaxConcurrentLocalActivityExecutionSize: 100,

        // Worker 身份
        Identity: "worker-1",

        // 死信队列配置
        DeadlockDetectionTimeout: 10 * time.Second,

        // 缓存配置
        WorkflowCacheSize: 600,

        // 速率限制
        WorkerActivitiesPerSecond: 100,
    })

    // 注册工作流
    w.RegisterWorkflow(OrderWorkflow)
    w.RegisterWorkflow(OrderWorkflowWithSignals)

    // 注册 Activities
    activities := &OrderActivities{
        inventoryClient: NewInventoryClient(),
        paymentClient:   NewPaymentClient(),
        shippingClient:  NewShippingClient(),
    }
    w.RegisterActivity(activities)

    // 启动 Worker
    err = w.Run(worker.InterruptCh())
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## Saga 模式完整实现

```go
// Saga 补偿事务模式

type Saga struct {
    compensations []func() error
}

func (s *Saga) Add(compensation func() error) {
    s.compensations = append(s.compensations, compensation)
}

func (s *Saga) Compensate() error {
    // 逆序执行补偿
    for i := len(s.compensations) - 1; i >= 0; i-- {
        if err := s.compensations[i](); err != nil {
            return err
        }
    }
    return nil
}

func SagaOrderWorkflow(ctx workflow.Context, order Order) (string, error) {
    saga := &Saga{}

    defer func() {
        if err := recover(); err != nil {
            // 发生 panic，执行补偿
            saga.Compensate()
        }
    }()

    // 步骤1: 预留库存
    var reservationID string
    err := workflow.ExecuteActivity(ctx, activities.ReserveInventory, order.Items).Get(ctx, &reservationID)
    if err != nil {
        return "", err
    }
    saga.Add(func() error {
        return workflow.ExecuteActivity(ctx, activities.CancelReservation, reservationID).Get(ctx, nil)
    })

    // 步骤2: 处理支付
    var payment PaymentResult
    err = workflow.ExecuteActivity(ctx, activities.ChargePayment, order.Total).Get(ctx, &payment)
    if err != nil {
        saga.Compensate()
        return "", err
    }
    saga.Add(func() error {
        return workflow.ExecuteActivity(ctx, activities.RefundPayment, payment.TransactionID).Get(ctx, nil)
    })

    // 步骤3: 确认库存扣减
    err = workflow.ExecuteActivity(ctx, activities.ConfirmReservation, reservationID).Get(ctx, nil)
    if err != nil {
        saga.Compensate()
        return "", err
    }

    // 步骤4: 创建发货
    var shipment ShipmentResult
    err = workflow.ExecuteActivity(ctx, activities.CreateShipment, order, payment.TransactionID).Get(ctx, &shipment)
    if err != nil {
        saga.Compensate()
        return "", err
    }

    return shipment.TrackingID, nil
}
```

---

## 高级模式：动态工作流

```go
// DSL 驱动的工作流
type WorkflowDSL struct {
    Steps []WorkflowStep `json:"steps"`
}

type WorkflowStep struct {
    Type       string                 `json:"type"`       // "activity", "condition", "parallel"
    Name       string                 `json:"name"`
    Activity   string                 `json:"activity,omitempty"`
    Args       map[string]interface{} `json:"args,omitempty"`
    Condition  string                 `json:"condition,omitempty"` // 条件表达式
    Branches   []WorkflowStep         `json:"branches,omitempty"`
    ParallelOf []WorkflowStep         `json:"parallel_of,omitempty"`
}

func DynamicWorkflow(ctx workflow.Context, dsl WorkflowDSL, input map[string]interface{}) (map[string]interface{}, error) {
    state := make(map[string]interface{})
    for k, v := range input {
        state[k] = v
    }

    for _, step := range dsl.Steps {
        switch step.Type {
        case "activity":
            result, err := executeDynamicActivity(ctx, step, state)
            if err != nil {
                return nil, err
            }
            state[step.Name] = result

        case "condition":
            condition := evaluateCondition(step.Condition, state)
            if condition {
                result, err := DynamicWorkflow(ctx, WorkflowDSL{Steps: step.Branches}, state)
                if err != nil {
                    return nil, err
                }
                for k, v := range result {
                    state[k] = v
                }
            }

        case "parallel":
            futures := make([]workflow.Future, 0, len(step.ParallelOf))
            for _, branch := range step.ParallelOf {
                future := workflow.ExecuteChildWorkflow(ctx, DynamicWorkflow, WorkflowDSL{Steps: []WorkflowStep{branch}}, state)
                futures = append(futures, future)
            }

            for _, future := range futures {
                var result map[string]interface{}
                if err := future.Get(ctx, &result); err != nil {
                    return nil, err
                }
                for k, v := range result {
                    state[k] = v
                }
            }
        }
    }

    return state, nil
}

func executeDynamicActivity(ctx workflow.Context, step WorkflowStep, state map[string]interface{}) (interface{}, error) {
    // 使用反射或注册表调用对应的 Activity
    // 实际实现需要 Activity 注册表
    return nil, nil
}
```
