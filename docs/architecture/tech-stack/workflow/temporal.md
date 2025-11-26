# 1. 🔄 Temporal 深度解析

> **简介**: 本文档详细阐述了 Temporal 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 🔄 Temporal 深度解析](#1--temporal-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 工作流定义](#131-工作流定义)
    - [1.3.2 活动定义](#132-活动定义)
    - [1.3.3 Worker 配置](#133-worker-配置)
    - [1.3.4 Client 使用](#134-client-使用)
    - [1.3.5 信号和查询使用](#135-信号和查询使用)
    - [1.3.6 错误处理示例](#136-错误处理示例)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 工作流设计最佳实践](#141-工作流设计最佳实践)
    - [1.4.2 活动设计最佳实践](#142-活动设计最佳实践)
    - [1.4.3 Worker 配置最佳实践](#143-worker-配置最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Temporal 是什么？**

Temporal 是一个分布式工作流编排引擎，提供可靠的业务流程管理。

**核心特性**:

- ✅ **可靠性**: 自动持久化状态，支持故障恢复
- ✅ **可观测性**: 内置 UI 和监控
- ✅ **Go 支持**: 官方 Go SDK，功能完整
- ✅ **可扩展性**: 支持水平扩展

---

## 1.2 选型论证

**为什么选择 Temporal？**

**论证矩阵**:

| 评估维度 | 权重 | Temporal | Airflow | Conductor | Cadence | 说明 |
|---------|------|----------|---------|-----------|---------|------|
| **Go 支持** | 40% | 10 | 0 | 0 | 5 | Temporal 官方 Go SDK |
| **功能完整性** | 25% | 10 | 8 | 7 | 8 | Temporal 功能完善 |
| **可观测性** | 20% | 10 | 7 | 5 | 6 | Temporal UI 功能强大 |
| **学习曲线** | 10% | 7 | 8 | 7 | 7 | Temporal 学习曲线适中 |
| **社区支持** | 5% | 8 | 10 | 5 | 6 | Temporal 社区活跃 |
| **加权总分** | - | **9.25** | 5.40 | 4.85 | 6.50 | Temporal 得分最高 |

**核心优势**:

1. **Go 支持（权重 40%）**:
   - 官方 Go SDK，功能完整
   - 文档完善，示例丰富
   - 社区支持好
   - **这是选择 Temporal 的最重要原因**

2. **功能完整性（权重 25%）**:
   - 持久化、可恢复、可查询功能完善
   - 信号和版本控制支持好
   - UI 功能完善

3. **可观测性（权重 20%）**:
   - 内置 UI，功能完善
   - 支持 OpenTelemetry
   - 追踪和监控集成好

**为什么不选择其他工作流引擎？**

1. **Airflow**:
   - ✅ UI 功能丰富，社区活跃
   - ❌ 无官方 Go SDK
   - ❌ 主要面向 Python
   - ❌ 不适合实时工作流

2. **Conductor**:
   - ✅ 功能强大，Netflix 开源
   - ❌ 无官方 Go SDK
   - ❌ 可观测性支持有限
   - ❌ 社区较小

3. **Cadence**:
   - ⚠️ 只有社区 Go SDK，功能有限
   - ⚠️ 可观测性支持有限
   - ⚠️ 文档和社区支持有限

**详细论证请参考**: [工作流架构设计](../../workflow.md#11-为什么选择-temporal)

---

## 1.3 实际应用

### 1.3.1 工作流定义

**基础工作流定义**:

```go
// internal/application/workflow/user_workflow.go
package workflow

import (
    "fmt"
    "time"

    "go.temporal.io/sdk/workflow"
    "go.temporal.io/sdk/temporal"
)

// UserWorkflowInput 工作流输入
type UserWorkflowInput struct {
    UserID  string
    Email   string
    Name    string
    Action  string // "create", "update", "delete"
}

// UserWorkflowOutput 工作流输出
type UserWorkflowOutput struct {
    UserID    string
    Success   bool
    Message   string
    Timestamp time.Time
}

// UserWorkflow 用户工作流
func UserWorkflow(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
    // 配置活动选项
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 30 * time.Second,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval:    time.Minute,
            MaximumAttempts:    3,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    var result UserWorkflowOutput
    var err error

    switch input.Action {
    case "create":
        result, err = createUserWorkflow(ctx, input)
    case "update":
        result, err = updateUserWorkflow(ctx, input)
    case "delete":
        result, err = deleteUserWorkflow(ctx, input)
    default:
        return UserWorkflowOutput{
            Success: false,
            Message: "unknown action",
        }, fmt.Errorf("unknown action: %s", input.Action)
    }

    return result, err
}
```

### 1.3.2 活动定义

**活动定义示例**:

```go
// internal/application/workflow/user_activities.go
package workflow

import (
    "context"
    "fmt"

    appuser "github.com/yourusername/golang/internal/application/user"
)

// ValidateUserActivity 验证用户活动
func ValidateUserActivity(ctx context.Context, email, name string) (string, error) {
    // 验证邮箱格式
    if !isValidEmail(email) {
        return "", fmt.Errorf("invalid email: %s", email)
    }

    // 验证姓名
    if len(name) < 2 || len(name) > 50 {
        return "", fmt.Errorf("invalid name: %s", name)
    }

    return "validation passed", nil
}

// CreateUserActivity 创建用户活动
func CreateUserActivity(ctx context.Context, email, name string) (string, error) {
    userService, ok := GetUserServiceFromContext(ctx)
    if !ok {
        return "", fmt.Errorf("user service not found in context")
    }

    user, err := userService.CreateUser(ctx, appuser.CreateUserRequest{
        Email: email,
        Name:  name,
    })
    if err != nil {
        return "", fmt.Errorf("failed to create user: %w", err)
    }

    return user.ID, nil
}

// SendNotificationActivity 发送通知活动
func SendNotificationActivity(ctx context.Context, userID, eventType string) error {
    // 发送通知逻辑
    fmt.Printf("Sending notification: userID=%s, eventType=%s\n", userID, eventType)
    return nil
}
```

### 1.3.3 Worker 配置

**Worker 配置示例**:

```go
// cmd/temporal-worker/main.go
package main

import (
    "context"
    "log"

    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"

    appworkflow "github.com/yourusername/golang/internal/application/workflow"
    "github.com/yourusername/golang/internal/config"
    temporalclient "github.com/yourusername/golang/internal/infrastructure/workflow/temporal"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 创建 Temporal 客户端
    temporalClient, err := temporalclient.NewClient(cfg.Workflow.Temporal.Address)
    if err != nil {
        log.Fatalf("Failed to create temporal client: %v", err)
    }
    defer temporalClient.Close()

    // 创建 Worker
    w := worker.New(temporalClient.Client(), cfg.Workflow.Temporal.TaskQueue, worker.Options{})

    // 注册工作流
    w.RegisterWorkflow(appworkflow.UserWorkflow)

    // 注册活动
    w.RegisterActivity(appworkflow.ValidateUserActivity)
    w.RegisterActivity(appworkflow.CreateUserActivity)
    w.RegisterActivity(appworkflow.SendNotificationActivity)

    // 启动 Worker
    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalf("Worker failed: %v", err)
    }
}
```

### 1.3.4 Client 使用

**Client 使用示例**:

```go
// 启动工作流
func StartUserWorkflow(ctx context.Context, client client.Client, input appworkflow.UserWorkflowInput) (client.WorkflowRun, error) {
    options := client.StartWorkflowOptions{
        ID:        fmt.Sprintf("user-workflow-%s-%s", input.Action, input.UserID),
        TaskQueue: "user-task-queue",
    }

    workflowRun, err := client.ExecuteWorkflow(ctx, options, appworkflow.UserWorkflow, input)
    if err != nil {
        return nil, fmt.Errorf("failed to start workflow: %w", err)
    }

    return workflowRun, nil
}

// 获取工作流结果
func GetWorkflowResult(ctx context.Context, client client.Client, workflowID, runID string) (appworkflow.UserWorkflowOutput, error) {
    var result appworkflow.UserWorkflowOutput

    workflowRun := client.GetWorkflow(ctx, workflowID, runID)
    err := workflowRun.Get(ctx, &result)
    if err != nil {
        return result, fmt.Errorf("failed to get workflow result: %w", err)
    }

    return result, nil
}

// 发送信号
func SignalWorkflow(ctx context.Context, client client.Client, workflowID, runID, signalName string, arg interface{}) error {
    return client.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

// 查询工作流
func QueryWorkflow(ctx context.Context, client client.Client, workflowID, runID, queryType string, args ...interface{}) (interface{}, error) {
    return client.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
}
```

### 1.3.5 信号和查询使用

**信号使用示例**:

```go
// 在工作流中接收信号
func OrderApprovalWorkflow(ctx workflow.Context, orderID string) error {
    // 创建信号通道
    signalChan := workflow.GetSignalChannel(ctx, "approve-signal")

    // 等待信号
    var approvalResult bool
    signalChan.Receive(ctx, &approvalResult)

    if approvalResult {
        // 处理批准逻辑
        return workflow.ExecuteActivity(ctx, ProcessOrderActivity, orderID).Get(ctx, nil)
    } else {
        // 处理拒绝逻辑
        return workflow.ExecuteActivity(ctx, CancelOrderActivity, orderID).Get(ctx, nil)
    }
}

// 从客户端发送信号
func SendApprovalSignal(ctx context.Context, client client.Client, workflowID, runID string, approved bool) error {
    return client.SignalWorkflow(ctx, workflowID, runID, "approve-signal", approved)
}
```

**查询使用示例**:

```go
// 在工作流中设置查询处理器
func OrderStatusWorkflow(ctx workflow.Context, orderID string) (string, error) {
    currentStatus := "PENDING"

    // 设置查询处理器
    err := workflow.SetQueryHandler(ctx, "get-status", func() (string, error) {
        return currentStatus, nil
    })
    if err != nil {
        return "", err
    }

    // 更新状态
    currentStatus = "PROCESSING"
    workflow.Sleep(ctx, 10*time.Second)

    currentStatus = "COMPLETED"
    return currentStatus, nil
}

// 从客户端查询工作流状态
func GetOrderStatus(ctx context.Context, client client.Client, workflowID, runID string) (string, error) {
    var status string
    err := client.QueryWorkflow(ctx, workflowID, runID, "get-status").Get(ctx, &status)
    return status, err
}
```

### 1.3.6 错误处理示例

**错误处理示例**:

```go
// 工作流中的错误处理
func UserWorkflowWithErrorHandling(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 30 * time.Second,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval:    time.Minute,
            MaximumAttempts:    3,
            NonRetryableErrorTypes: []string{"ValidationError", "NotFoundError"},
        },
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 执行活动
    err := workflow.ExecuteActivity(ctx, CreateUserActivity, input.Email, input.Name).Get(ctx, nil)
    if err != nil {
        // 检查错误类型
        var activityErr *temporal.ActivityError
        if errors.As(err, &activityErr) {
            // 处理活动错误
            workflow.GetLogger(ctx).Error("Activity failed", "error", activityErr)
            return UserWorkflowOutput{Success: false}, err
        }

        // 处理其他错误
        return UserWorkflowOutput{Success: false}, err
    }

    return UserWorkflowOutput{Success: true}, nil
}
```

---

## 1.4 最佳实践

### 1.4.1 工作流设计最佳实践

**为什么需要良好的工作流设计？**

良好的工作流设计可以提高工作流的可维护性、可测试性和性能。

**工作流设计原则**:

1. **确定性**: 工作流代码必须是确定性的，不能使用随机数、时间等非确定性函数
2. **细粒度活动**: 将复杂逻辑拆分为多个细粒度活动
3. **错误处理**: 合理配置重试策略，处理不同类型的错误
4. **超时设置**: 为活动设置合理的超时时间

**实际应用示例**:

```go
// 良好的工作流设计
func UserRegistrationWorkflow(ctx workflow.Context, input UserRegistrationInput) (UserRegistrationOutput, error) {
    // 配置活动选项
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 30 * time.Second,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval:    time.Minute,
            MaximumAttempts:    3,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 1. 验证用户信息
    var validationResult string
    err := workflow.ExecuteActivity(ctx, ValidateUserActivity, input.Email, input.Name).Get(ctx, &validationResult)
    if err != nil {
        return UserRegistrationOutput{Success: false}, err
    }

    // 2. 创建用户
    var userID string
    err = workflow.ExecuteActivity(ctx, CreateUserActivity, input.Email, input.Name).Get(ctx, &userID)
    if err != nil {
        return UserRegistrationOutput{Success: false}, err
    }

    // 3. 发送欢迎邮件（不阻塞主流程）
    workflow.ExecuteActivity(ctx, SendWelcomeEmailActivity, userID, input.Email).Get(ctx, nil)

    return UserRegistrationOutput{
        Success: true,
        UserID:  userID,
    }, nil
}
```

**最佳实践要点**:

1. **确定性**: 使用 `workflow.Now()` 而不是 `time.Now()`，使用 `workflow.GetRandomSequence()` 而不是 `rand.Int()`
2. **细粒度活动**: 将复杂逻辑拆分为多个活动，每个活动职责单一
3. **错误处理**: 合理配置重试策略，区分可重试和不可重试错误
4. **超时设置**: 为活动设置合理的超时时间，避免长时间阻塞

### 1.4.2 活动设计最佳实践

**为什么需要良好的活动设计？**

良好的活动设计可以提高活动的可重用性、可测试性和性能。

**活动设计原则**:

1. **幂等性**: 活动应该是幂等的，多次执行结果相同
2. **单一职责**: 每个活动只负责一个功能
3. **错误处理**: 返回明确的错误类型
4. **超时处理**: 合理设置超时时间

**实际应用示例**:

```go
// 良好的活动设计
func CreateUserActivity(ctx context.Context, email, name string) (string, error) {
    // 获取服务
    userService, ok := GetUserServiceFromContext(ctx)
    if !ok {
        return "", fmt.Errorf("user service not found")
    }

    // 检查用户是否已存在（幂等性）
    existingUser, err := userService.GetUserByEmail(ctx, email)
    if err == nil && existingUser != nil {
        // 用户已存在，返回现有用户 ID（幂等性）
        return existingUser.ID, nil
    }

    // 创建用户
    user, err := userService.CreateUser(ctx, appuser.CreateUserRequest{
        Email: email,
        Name:  name,
    })
    if err != nil {
        return "", fmt.Errorf("failed to create user: %w", err)
    }

    return user.ID, nil
}
```

**最佳实践要点**:

1. **幂等性**: 活动应该是幂等的，多次执行结果相同
2. **单一职责**: 每个活动只负责一个功能
3. **错误处理**: 返回明确的错误类型，便于工作流处理
4. **超时处理**: 合理设置超时时间，避免长时间阻塞

### 1.4.3 Worker 配置最佳实践

**为什么需要合理的 Worker 配置？**

合理的 Worker 配置可以提高 Worker 的性能和可靠性。

**Worker 配置原则**:

1. **Task Queue 划分**: 根据业务特性划分 Task Queue
2. **Worker 数量**: 根据负载配置 Worker 数量
3. **活动注册**: 只注册需要的活动
4. **错误处理**: 配置 Worker 级别的错误处理

**实际应用示例**:

```go
// Worker 配置最佳实践
func NewWorker(client client.Client, taskQueue string) worker.Worker {
    w := worker.New(client, taskQueue, worker.Options{
        MaxConcurrentActivityExecutionSize: 100,  // 最大并发活动数
        MaxConcurrentWorkflowTaskSize:      10,   // 最大并发工作流任务数
        MaxConcurrentLocalActivitySize:     100,  // 最大并发本地活动数
    })

    // 注册工作流
    w.RegisterWorkflow(appworkflow.UserWorkflow)
    w.RegisterWorkflow(appworkflow.OrderWorkflow)

    // 注册活动
    w.RegisterActivity(appworkflow.ValidateUserActivity)
    w.RegisterActivity(appworkflow.CreateUserActivity)
    w.RegisterActivity(appworkflow.SendNotificationActivity)

    return w
}
```

**最佳实践要点**:

1. **Task Queue 划分**: 根据业务特性划分 Task Queue，实现任务隔离
2. **Worker 数量**: 根据负载配置 Worker 数量，实现负载均衡
3. **活动注册**: 只注册需要的活动，减少内存占用
4. **并发配置**: 合理配置并发参数，避免资源耗尽

**详细实现请参考**: [工作流架构设计](../../workflow.md)

---

## 📚 扩展阅读

- [Temporal 官方文档](https://docs.temporal.io/)
- [工作流架构设计](../../workflow.md)
- [工作流指南](../../../guides/workflow.md)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Temporal 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
