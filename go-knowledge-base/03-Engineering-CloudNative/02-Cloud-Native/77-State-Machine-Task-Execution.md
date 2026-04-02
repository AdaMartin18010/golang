# 状态机任务执行 (State Machine Task Execution)

> **分类**: 工程与云原生
> **标签**: #state-machine #workflow #execution-engine
> **参考**: AWS Step Functions, Temporal State Machines

---

## 状态机架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    State Machine Execution Engine                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    State Machine Definition                          │   │
│  │                                                                      │   │
│  │  States:                                                             │   │
│  │    - StartAt: "ValidateOrder"                                        │   │
│  │    - States:                                                         │   │
│  │      - ValidateOrder:                                                │   │
│  │          Type: Task                                                  │   │
│  │          Next: CheckInventory                                        │   │
│  │      - CheckInventory:                                               │   │
│  │          Type: Task                                                  │   │
│  │          Next: ProcessPayment                                        │   │
│  │      - ProcessPayment:                                               │   │
│  │          Type: Task                                                  │   │
│  │          Catch: [{Error: ["PaymentFailed"], Next: "HandleFailure"}] │   │
│  │          Next: ShipOrder                                             │   │
│  │      - ShipOrder:                                                    │   │
│  │          Type: Task                                                  │   │
│  │          End: true                                                   │   │
│  │      - HandleFailure:                                                │   │
│  │          Type: Task                                                  │   │
│  │          End: true                                                   │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    State Machine Execution                             │   │
│  │                                                                      │   │
│  │  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐       │   │
│  │  │ Validate │───►│ Inventory│───►│ Payment  │───►│  Ship    │       │   │
│  │  │  Order   │    │  Check   │    │ Process  │    │  Order   │       │   │
│  │  └──────────┘    └──────────┘    └────┬─────┘    └──────────┘       │   │
│  │                                      │                              │   │
│  │                                      ▼                              │   │
│  │                                 ┌──────────┐                         │   │
│  │                                 │  Failure │                         │   │
│  │                                 │ Handler  │                         │   │
│  │                                 └──────────┘                         │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    State Types                                       │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │   Task   │  │  Choice  │  │   Wait   │  │  Parallel│            │   │
│  │  │(Execute) │  │(Branch)  │  │(Delay)   │  │(Fork)    │            │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐                          │   │
│  │  │   Map    │  │   Pass   │  │  Succeed │                          │   │
│  │  │(Iterate) │  │(No-op)   │  │  / Fail  │                          │   │
│  │  └──────────┘  └──────────┘  └──────────┘                          │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心实现

```go
package statemachine

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

// StateType 状态类型
type StateType string

const (
    StateTypeTask     StateType = "Task"
    StateTypeChoice   StateType = "Choice"
    StateTypeWait     StateType = "Wait"
    StateTypeParallel StateType = "Parallel"
    StateTypeMap      StateType = "Map"
    StateTypePass     StateType = "Pass"
    StateTypeSucceed  StateType = "Succeed"
    StateTypeFail     StateType = "Fail"
)

// State 状态定义
type State struct {
    Name string    `json:"Name"`
    Type StateType `json:"Type"`

    // Task 状态
    Resource string                 `json:"Resource,omitempty"` // 资源 ARN/函数名
    InputPath  string               `json:"InputPath,omitempty"`
    ResultPath string               `json:"ResultPath,omitempty"`
    OutputPath string               `json:"OutputPath,omitempty"`
    Parameters map[string]interface{} `json:"Parameters,omitempty"`
    ResultSelector map[string]interface{} `json:"ResultSelector,omitempty"`
    Retry      []RetryPolicy        `json:"Retry,omitempty"`
    Catch      []CatchPolicy        `json:"Catch,omitempty"`
    Timeout    *int                 `json:"Timeout,omitempty"` // 秒

    // 转换
    Next       string `json:"Next,omitempty"`
    End        bool   `json:"End,omitempty"`

    // Choice 状态
    Choices    []ChoiceRule `json:"Choices,omitempty"`
    Default    string       `json:"Default,omitempty"`

    // Wait 状态
    Seconds    *int       `json:"Seconds,omitempty"`
    Timestamp  *time.Time `json:"Timestamp,omitempty"`
    TimestampPath string  `json:"TimestampPath,omitempty"`

    // Parallel 状态
    Branches   []StateMachine `json:"Branches,omitempty"`

    // Map 状态
    Iterator   *StateMachine `json:"Iterator,omitempty"`
    ItemsPath  string        `json:"ItemsPath,omitempty"`
    MaxConcurrency int       `json:"MaxConcurrency,omitempty"`

    // Fail 状态
    Error      string `json:"Error,omitempty"`
    Cause      string `json:"Cause,omitempty"`

    // Succeed 状态 (无额外字段)
}

// StateMachine 状态机定义
type StateMachine struct {
    Comment    string            `json:"Comment,omitempty"`
    StartAt    string            `json:"StartAt"`
    States     map[string]*State `json:"States"`
    Timeout    *int              `json:"Timeout,omitempty"`
    Version    string            `json:"Version,omitempty"`
}

// RetryPolicy 重试策略
type RetryPolicy struct {
    ErrorEquals    []string `json:"ErrorEquals"`
    Interval       int      `json:"IntervalSeconds,omitempty"` // 默认 1
    MaxAttempts    int      `json:"MaxAttempts,omitempty"`     // 默认 3
    BackoffRate    float64  `json:"BackoffRate,omitempty"`     // 默认 2.0
}

// CatchPolicy 错误处理策略
type CatchPolicy struct {
    ErrorEquals []string `json:"ErrorEquals"`
    Next        string   `json:"Next"`
    ResultPath  string   `json:"ResultPath,omitempty"`
}

// ChoiceRule 选择规则
type ChoiceRule struct {
    Variable    string      `json:"Variable"` // $.variable 路径
    Next        string      `json:"Next"`

    // 比较操作
    StringEquals  string    `json:"StringEquals,omitempty"`
    NumericEquals float64   `json:"NumericEquals,omitempty"`
    BooleanEquals bool      `json:"BooleanEquals,omitempty"`
    TimestampEquals time.Time `json:"TimestampEquals,omitempty"`

    // 其他比较...
    And []ChoiceRule `json:"And,omitempty"`
    Or  []ChoiceRule `json:"Or,omitempty"`
    Not *ChoiceRule  `json:"Not,omitempty"`
}

// Execution 执行实例
type Execution struct {
    ID          string
    StateMachine *StateMachine
    Input       interface{}
    Output      interface{}

    // 执行状态
    CurrentState string
    StateHistory []StateHistory
    Status       ExecutionStatus
    Error        string
    Cause        string

    StartTime    time.Time
    EndTime      *time.Time
}

// ExecutionStatus 执行状态
type ExecutionStatus string

const (
    ExecutionStatusRunning   ExecutionStatus = "RUNNING"
    ExecutionStatusSucceeded ExecutionStatus = "SUCCEEDED"
    ExecutionStatusFailed    ExecutionStatus = "FAILED"
    ExecutionStatusTimedOut  ExecutionStatus = "TIMED_OUT"
    ExecutionStatusAborted   ExecutionStatus = "ABORTED"
)

// StateHistory 状态历史
type StateHistory struct {
    EnteredTime time.Time
    Name        string
    Type        StateType
    Input       interface{}
    Output      interface{}
    Error       string
}
```

---

## 执行引擎

```go
package statemachine

import (
    "context"
    "encoding/json"
    "fmt"
    "reflect"
    "time"
)

// Executor 状态机执行器
type Executor struct {
    taskHandlers map[string]TaskHandler
}

// TaskHandler 任务处理器
type TaskHandler func(ctx context.Context, input interface{}) (interface{}, error)

// NewExecutor 创建执行器
func NewExecutor() *Executor {
    return &Executor{
        taskHandlers: make(map[string]TaskHandler),
    }
}

// RegisterTaskHandler 注册任务处理器
func (e *Executor) RegisterTaskHandler(resource string, handler TaskHandler) {
    e.taskHandlers[resource] = handler
}

// Execute 执行状态机
func (e *Executor) Execute(ctx context.Context, sm *StateMachine, input interface{}) (*Execution, error) {
    exec := &Execution{
        ID:           generateExecutionID(),
        StateMachine: sm,
        Input:        input,
        CurrentState: sm.StartAt,
        Status:       ExecutionStatusRunning,
        StartTime:    time.Now(),
    }

    currentInput := input

    for exec.Status == ExecutionStatusRunning {
        state, ok := sm.States[exec.CurrentState]
        if !ok {
            exec.Status = ExecutionStatusFailed
            exec.Error = "StateNotFound"
            exec.Cause = fmt.Sprintf("State %s not found", exec.CurrentState)
            break
        }

        // 记录历史
        history := StateHistory{
            EnteredTime: time.Now(),
            Name:        state.Name,
            Type:        state.Type,
            Input:       currentInput,
        }

        // 执行状态
        output, nextState, err := e.executeState(ctx, state, currentInput)

        history.Output = output
        if err != nil {
            history.Error = err.Error()
        }
        exec.StateHistory = append(exec.StateHistory, history)

        if err != nil {
            // 检查是否有 Catch
            handled := false
            for _, catch := range state.Catch {
                if e.errorMatches(err, catch.ErrorEquals) {
                    exec.CurrentState = catch.Next
                    currentInput = map[string]interface{}{
                        "Error": err.Error(),
                    }
                    handled = true
                    break
                }
            }

            if !handled {
                exec.Status = ExecutionStatusFailed
                exec.Error = "States.TaskFailed"
                exec.Cause = err.Error()
                break
            }
        } else {
            currentInput = output

            if state.End {
                exec.Status = ExecutionStatusSucceeded
                exec.Output = output
            } else if state.Next != "" {
                exec.CurrentState = state.Next
            } else if state.Type == StateTypeChoice {
                // Choice 状态在 executeState 中设置 Next
            } else if state.Type == StateTypeFail {
                exec.Status = ExecutionStatusFailed
                exec.Error = state.Error
                exec.Cause = state.Cause
            } else if state.Type == StateTypeSucceed {
                exec.Status = ExecutionStatusSucceeded
                exec.Output = output
            }
        }

        // 检查超时
        if sm.Timeout != nil {
            elapsed := time.Since(exec.StartTime)
            if elapsed > time.Duration(*sm.Timeout)*time.Second {
                exec.Status = ExecutionStatusTimedOut
                break
            }
        }
    }

    now := time.Now()
    exec.EndTime = &now

    return exec, nil
}

func (e *Executor) executeState(ctx context.Context, state *State, input interface{}) (interface{}, string, error) {
    switch state.Type {
    case StateTypeTask:
        return e.executeTask(ctx, state, input)

    case StateTypeChoice:
        return e.executeChoice(state, input)

    case StateTypeWait:
        return e.executeWait(ctx, state, input)

    case StateTypeParallel:
        return e.executeParallel(ctx, state, input)

    case StateTypeMap:
        return e.executeMap(ctx, state, input)

    case StateTypePass:
        return e.executePass(state, input)

    case StateTypeSucceed:
        return input, "", nil

    case StateTypeFail:
        return nil, "", fmt.Errorf(state.Cause)

    default:
        return nil, "", fmt.Errorf("unknown state type: %s", state.Type)
    }
}

func (e *Executor) executeTask(ctx context.Context, state *State, input interface{}) (interface{}, string, error) {
    handler, ok := e.taskHandlers[state.Resource]
    if !ok {
        return nil, "", fmt.Errorf("no handler for resource: %s", state.Resource)
    }

    // 应用 InputPath 过滤
    taskInput := e.applyPath(input, state.InputPath)

    // 应用 Parameters 转换
    if state.Parameters != nil {
        taskInput = e.applyParameters(taskInput, state.Parameters)
    }

    // 执行重试逻辑
    var result interface{}
    var err error

    retryPolicy := e.getRetryPolicy(state.Retry)
    for attempt := 0; attempt <= retryPolicy.MaxAttempts; attempt++ {
        if attempt > 0 {
            // 退避等待
            delay := time.Duration(retryPolicy.Interval) * time.Second
            delay = time.Duration(float64(delay) * pow(retryPolicy.BackoffRate, float64(attempt-1)))
            time.Sleep(delay)
        }

        result, err = handler(ctx, taskInput)
        if err == nil {
            break
        }

        // 检查是否应该重试
        if attempt < retryPolicy.MaxAttempts && e.shouldRetry(err, retryPolicy.ErrorEquals) {
            continue
        }
        break
    }

    if err != nil {
        return nil, "", err
    }

    // 应用 ResultSelector
    if state.ResultSelector != nil {
        result = e.applyParameters(result, state.ResultSelector)
    }

    // 应用 ResultPath
    output := e.applyResultPath(input, result, state.ResultPath)

    // 应用 OutputPath
    output = e.applyPath(output, state.OutputPath)

    return output, state.Next, nil
}

func (e *Executor) executeChoice(state *State, input interface{}) (interface{}, string, error) {
    for _, rule := range state.Choices {
        if e.evaluateChoiceRule(rule, input) {
            return input, rule.Next, nil
        }
    }

    // 没有匹配的规则，使用 Default
    if state.Default != "" {
        return input, state.Default, nil
    }

    return nil, "", fmt.Errorf("no matching choice rule and no default")
}

func (e *Executor) executeWait(ctx context.Context, state *State, input interface{}) (interface{}, string, error) {
    var duration time.Duration

    if state.Seconds != nil {
        duration = time.Duration(*state.Seconds) * time.Second
    } else if state.Timestamp != nil {
        duration = state.Timestamp.Sub(time.Now())
        if duration < 0 {
            duration = 0
        }
    }

    select {
    case <-ctx.Done():
        return nil, "", ctx.Err()
    case <-time.After(duration):
    }

    return input, state.Next, nil
}

func (e *Executor) executeParallel(ctx context.Context, state *State, input interface{}) (interface{}, string, error) {
    // 并行执行多个分支
    results := make([]interface{}, len(state.Branches))
    errors := make([]error, len(state.Branches))

    // 使用 WaitGroup 并行执行
    // 实际实现中需要考虑并发限制和错误处理

    for i, err := range errors {
        if err != nil {
            return nil, "", fmt.Errorf("branch %d failed: %w", i, err)
        }
    }

    return results, state.Next, nil
}

func (e *Executor) executeMap(ctx context.Context, state *State, input interface{}) (interface{}, string, error) {
    // 迭代处理集合
    items := e.getItems(input, state.ItemsPath)

    results := make([]interface{}, 0, len(items))

    for _, item := range items {
        // 为每个元素创建执行
        // 实际实现中需要考虑 MaxConcurrency
        _ = item
    }

    return results, state.Next, nil
}

func (e *Executor) executePass(state *State, input interface{}) (interface{}, string, error) {
    result := input

    // Pass 状态可以直接设置 Result
    // result = state.Result

    // 应用 ResultPath
    output := e.applyResultPath(input, result, state.ResultPath)

    return output, state.Next, nil
}

// 辅助方法

func (e *Executor) applyPath(input interface{}, path string) interface{} {
    if path == "" || path == "$" {
        return input
    }
    // 实现 JSONPath 解析
    return input
}

func (e *Executor) applyParameters(input interface{}, params map[string]interface{}) interface{} {
    // 实现参数替换
    return params
}

func (e *Executor) applyResultPath(original, result interface{}, path string) interface{} {
    if path == "" {
        return result
    }
    // 将结果写入指定路径
    return result
}

func (e *Executor) evaluateChoiceRule(rule ChoiceRule, input interface{}) bool {
    // 获取变量值
    value := e.getValueByPath(input, rule.Variable)

    // 比较操作
    if rule.StringEquals != "" {
        return value == rule.StringEquals
    }
    if rule.NumericEquals != 0 {
        if v, ok := value.(float64); ok {
            return v == rule.NumericEquals
        }
    }
    // ... 其他比较

    return false
}

func (e *Executor) getValueByPath(input interface{}, path string) interface{} {
    // 实现 JSONPath 解析
    return input
}

func (e *Executor) getItems(input interface{}, path string) []interface{} {
    // 获取集合元素
    return []interface{}{input}
}

func (e *Executor) getRetryPolicy(policies []RetryPolicy) RetryPolicy {
    if len(policies) == 0 {
        return RetryPolicy{
            ErrorEquals: []string{"States.ALL"},
            Interval:    1,
            MaxAttempts: 3,
            BackoffRate: 2.0,
        }
    }
    return policies[0]
}

func (e *Executor) shouldRetry(err error, errorEquals []string) bool {
    for _, ee := range errorEquals {
        if ee == "States.ALL" {
            return true
        }
        // 检查错误类型匹配
    }
    return false
}

func (e *Executor) errorMatches(err error, errorEquals []string) bool {
    for _, ee := range errorEquals {
        if ee == "States.ALL" {
            return true
        }
        if ee == "States.TaskFailed" {
            return true
        }
        // ...
    }
    return false
}

func generateExecutionID() string {
    return fmt.Sprintf("exec-%d", time.Now().UnixNano())
}

func pow(x, y float64) float64 {
    result := 1.0
    for i := 0; i < int(y); i++ {
        result *= x
    }
    return result
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"

    "statemachine"
)

func main() {
    // 定义状态机
    smJSON := `{
        "Comment": "Order Processing",
        "StartAt": "ValidateOrder",
        "States": {
            "ValidateOrder": {
                "Type": "Task",
                "Resource": "validate-order",
                "Next": "CheckInventory"
            },
            "CheckInventory": {
                "Type": "Task",
                "Resource": "check-inventory",
                "Next": "ProcessPayment"
            },
            "ProcessPayment": {
                "Type": "Task",
                "Resource": "process-payment",
                "Retry": [{
                    "ErrorEquals": ["PaymentServiceUnavailable"],
                    "IntervalSeconds": 2,
                    "MaxAttempts": 3,
                    "BackoffRate": 2
                }],
                "Catch": [{
                    "ErrorEquals": ["PaymentFailed"],
                    "Next": "HandlePaymentFailure"
                }],
                "Next": "ShipOrder"
            },
            "ShipOrder": {
                "Type": "Task",
                "Resource": "ship-order",
                "End": true
            },
            "HandlePaymentFailure": {
                "Type": "Task",
                "Resource": "handle-failure",
                "End": true
            }
        }
    }`

    var sm statemachine.StateMachine
    if err := json.Unmarshal([]byte(smJSON), &sm); err != nil {
        panic(err)
    }

    // 创建执行器
    executor := statemachine.NewExecutor()

    // 注册任务处理器
    executor.RegisterTaskHandler("validate-order", func(ctx context.Context, input interface{}) (interface{}, error) {
        fmt.Println("Validating order...")
        return map[string]interface{}{"valid": true}, nil
    })

    executor.RegisterTaskHandler("check-inventory", func(ctx context.Context, input interface{}) (interface{}, error) {
        fmt.Println("Checking inventory...")
        return map[string]interface{}{"inStock": true}, nil
    })

    executor.RegisterTaskHandler("process-payment", func(ctx context.Context, input interface{}) (interface{}, error) {
        fmt.Println("Processing payment...")
        return map[string]interface{}{"transactionId": "txn-123"}, nil
    })

    executor.RegisterTaskHandler("ship-order", func(ctx context.Context, input interface{}) (interface{}, error) {
        fmt.Println("Shipping order...")
        return map[string]interface{}{"trackingId": "track-456"}, nil
    })

    // 执行
    input := map[string]interface{}{
        "orderId": "order-789",
        "items": []string{"item1", "item2"},
    }

    execution, err := executor.Execute(context.Background(), &sm, input)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Execution status: %s\n", execution.Status)
    fmt.Printf("Output: %v\n", execution.Output)
}
```
