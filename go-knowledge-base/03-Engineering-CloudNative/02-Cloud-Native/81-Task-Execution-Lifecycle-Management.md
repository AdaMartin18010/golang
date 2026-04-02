# 任务执行生命周期管理 (Task Execution Lifecycle Management)

> **分类**: 工程与云原生
> **标签**: #task-lifecycle #state-management #execution-flow
> **参考**: AWS Step Functions, Temporal Workflow Engine

---

## 任务生命周期状态机（生产级实现）

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Task Execution State Machine                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌──────────┐                                                               │
│   │  PENDING │◄────────────────────────────────────────────────────┐        │
│   └────┬─────┘                                                     │        │
│        │ trigger                                                   │        │
│        ▼                                                           │        │
│   ┌──────────┐    ┌──────────┐    ┌──────────┐                    │        │
│   │ ENQUEUED │───►│ SCHEDULED│───►│ RUNNING  │                    │        │
│   └──────────┘    └────┬─────┘    └────┬─────┘                    │        │
│                        │               │                          │        │
│                        │    ┌──────────┼──────────┐               │        │
│                        │    │          │          │               │        │
│                        ▼    ▼          ▼          ▼               │        │
│   ┌──────────┐    ┌──────────┐   ┌──────────┐  ┌──────────┐      │        │
│   │  PAUSED  │    │ CANCELLED│   │ COMPLETED│  │  FAILED  │      │        │
│   └────┬─────┘    └──────────┘   └──────────┘  └────┬─────┘      │        │
│        │                                            │             │        │
│        │ resume                                     │ retry       │        │
│        └────────────► (back to SCHEDULED) ◄─────────┘             │        │
│                                                                    │        │
│   ┌──────────┐    ┌──────────┐                                    │        │
│   │ TIMED_OUT│───►│  RETRYING│────────────────────────────────────┘        │
│   └──────────┘    └──────────┘                                            │
│                                                                              │
│   State Transition Rules:                                                   │
│   - PENDING → ENQUEUED: Task created and added to queue                     │
│   - ENQUEUED → SCHEDULED: Worker picked up task                             │
│   - SCHEDULED → RUNNING: Execution started                                  │
│   - RUNNING → COMPLETED: Successful execution                               │
│   - RUNNING → FAILED: Error occurred (retryable or non-retryable)          │
│   - RUNNING → TIMED_OUT: Execution exceeded timeout                         │
│   - FAILED → RETRYING: Retry policy triggered                               │
│   - RETRYING → SCHEDULED: Delay before retry                                │
│   - RUNNING → CANCELLED: User or system cancellation                        │
│   - SCHEDULED → PAUSED: Dependency not met or manual pause                  │
│   - PAUSED → SCHEDULED: Dependency resolved or manual resume                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整生命周期管理实现

```go
package lifecycle

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// TaskState 任务状态类型
type TaskState string

const (
    TaskStatePending    TaskState = "PENDING"
    TaskStateEnqueued   TaskState = "ENQUEUED"
    TaskStateScheduled  TaskState = "SCHEDULED"
    TaskStateRunning    TaskState = "RUNNING"
    TaskStateCompleted  TaskState = "COMPLETED"
    TaskStateFailed     TaskState = "FAILED"
    TaskStateCancelled  TaskState = "CANCELLED"
    TaskStatePaused     TaskState = "PAUSED"
    TaskStateTimedOut   TaskState = "TIMED_OUT"
    TaskStateRetrying   TaskState = "RETRYING"
)

// IsTerminal 是否为终止状态
func (s TaskState) IsTerminal() bool {
    return s == TaskStateCompleted || s == TaskStateFailed ||
           s == TaskStateCancelled || s == TaskStateTimedOut
}

// CanTransitionTo 检查状态转换是否合法
func (s TaskState) CanTransitionTo(target TaskState) bool {
    validTransitions := map[TaskState][]TaskState{
        TaskStatePending:   {TaskStateEnqueued, TaskStateCancelled},
        TaskStateEnqueued:  {TaskStateScheduled, TaskStateCancelled},
        TaskStateScheduled: {TaskStateRunning, TaskStatePaused, TaskStateCancelled},
        TaskStateRunning:   {TaskStateCompleted, TaskStateFailed, TaskStateTimedOut, TaskStateCancelled},
        TaskStateFailed:    {TaskStateRetrying, TaskStateCancelled},
        TaskStateTimedOut:  {TaskStateRetrying, TaskStateFailed, TaskStateCancelled},
        TaskStateRetrying:  {TaskStateScheduled, TaskStateCancelled},
        TaskStatePaused:    {TaskStateScheduled, TaskStateCancelled},
    }

    valid, ok := validTransitions[s]
    if !ok {
        return false
    }

    for _, v := range valid {
        if v == target {
            return true
        }
    }
    return false
}

// TaskExecution 任务执行实例
type TaskExecution struct {
    ID           string                 `json:"id"`
    TaskID       string                 `json:"task_id"`
    State        TaskState              `json:"state"`
    Input        map[string]interface{} `json:"input"`
    Output       map[string]interface{} `json:"output"`
    Error        *ExecutionError        `json:"error,omitempty"`

    // 时间戳
    CreatedAt    time.Time              `json:"created_at"`
    EnqueuedAt   *time.Time             `json:"enqueued_at,omitempty"`
    ScheduledAt  *time.Time             `json:"scheduled_at,omitempty"`
    StartedAt    *time.Time             `json:"started_at,omitempty"`
    CompletedAt  *time.Time             `json:"completed_at,omitempty"`

    // 执行信息
    WorkerID     string                 `json:"worker_id,omitempty"`
    Attempt      int                    `json:"attempt"`
    MaxAttempts  int                    `json:"max_attempts"`

    // 超时配置
    Timeout      time.Duration          `json:"timeout"`
    Deadline     *time.Time             `json:"deadline,omitempty"`

    // 上下文传播
    TraceID      string                 `json:"trace_id"`
    SpanID       string                 `json:"span_id"`
    ParentSpanID string                 `json:"parent_span_id,omitempty"`

    // 元数据
    Metadata     map[string]string      `json:"metadata,omitempty"`

    mu           sync.RWMutex
}

// ExecutionError 执行错误
type ExecutionError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    StackTrace string `json:"stack_trace,omitempty"`
    IsRetryable bool  `json:"is_retryable"`
}

func (e *ExecutionError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// LifecycleManager 生命周期管理器
type LifecycleManager struct {
    store        ExecutionStore
    hooks        map[TaskState][]LifecycleHook
    mu           sync.RWMutex

    // 状态机处理器
    handlers     map[TaskState]StateHandler
}

// ExecutionStore 执行存储接口
type ExecutionStore interface {
    Create(ctx context.Context, exec *TaskExecution) error
    Update(ctx context.Context, exec *TaskExecution) error
    Get(ctx context.Context, id string) (*TaskExecution, error)
    GetByTaskID(ctx context.Context, taskID string) ([]*TaskExecution, error)
    ListActive(ctx context.Context) ([]*TaskExecution, error)
}

// LifecycleHook 生命周期钩子
type LifecycleHook interface {
    OnStateTransition(ctx context.Context, exec *TaskExecution, from, to TaskState) error
}

// StateHandler 状态处理器
type StateHandler interface {
    Enter(ctx context.Context, exec *TaskExecution) error
    Execute(ctx context.Context, exec *TaskExecution) (*Transition, error)
    Exit(ctx context.Context, exec *TaskExecution) error
}

// Transition 状态转换
type Transition struct {
    To      TaskState
    Reason  string
    Payload map[string]interface{}
}

// NewLifecycleManager 创建生命周期管理器
func NewLifecycleManager(store ExecutionStore) *LifecycleManager {
    lm := &LifecycleManager{
        store:    store,
        hooks:    make(map[TaskState][]LifecycleHook),
        handlers: make(map[TaskState]StateHandler),
    }

    // 注册默认处理器
    lm.registerDefaultHandlers()

    return lm
}

// RegisterHook 注册生命周期钩子
func (lm *LifecycleManager) RegisterHook(state TaskState, hook LifecycleHook) {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    lm.hooks[state] = append(lm.hooks[state], hook)
}

// RegisterHandler 注册状态处理器
func (lm *LifecycleManager) RegisterHandler(state TaskState, handler StateHandler) {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    lm.handlers[state] = handler
}

// CreateExecution 创建执行实例
func (lm *LifecycleManager) CreateExecution(ctx context.Context, taskID string, input map[string]interface{}, opts ExecutionOptions) (*TaskExecution, error) {
    exec := &TaskExecution{
        ID:          generateExecutionID(),
        TaskID:      taskID,
        State:       TaskStatePending,
        Input:       input,
        Attempt:     1,
        MaxAttempts: opts.MaxAttempts,
        Timeout:     opts.Timeout,
        CreatedAt:   time.Now(),
        TraceID:     opts.TraceID,
        SpanID:      generateSpanID(),
        Metadata:    opts.Metadata,
    }

    if opts.TraceID == "" {
        exec.TraceID = generateTraceID()
    }

    if err := lm.store.Create(ctx, exec); err != nil {
        return nil, fmt.Errorf("failed to create execution: %w", err)
    }

    return exec, nil
}

// TransitionState 执行状态转换
func (lm *LifecycleManager) TransitionState(ctx context.Context, exec *TaskExecution, to TaskState, reason string) error {
    exec.mu.Lock()
    defer exec.mu.Unlock()

    from := exec.State

    // 验证状态转换
    if !from.CanTransitionTo(to) {
        return fmt.Errorf("invalid state transition from %s to %s", from, to)
    }

    // 执行钩子
    if err := lm.executeHooks(ctx, exec, from, to); err != nil {
        return fmt.Errorf("lifecycle hook failed: %w", err)
    }

    // 更新状态
    exec.State = to

    // 更新时间戳
    now := time.Now()
    switch to {
    case TaskStateEnqueued:
        exec.EnqueuedAt = &now
    case TaskStateScheduled:
        exec.ScheduledAt = &now
    case TaskStateRunning:
        exec.StartedAt = &now
        if exec.Timeout > 0 {
            deadline := now.Add(exec.Timeout)
            exec.Deadline = &deadline
        }
    case TaskStateCompleted, TaskStateFailed, TaskStateCancelled, TaskStateTimedOut:
        exec.CompletedAt = &now
    }

    // 持久化
    if err := lm.store.Update(ctx, exec); err != nil {
        return fmt.Errorf("failed to update execution: %w", err)
    }

    return nil
}

// Execute 执行状态机
func (lm *LifecycleManager) Execute(ctx context.Context, exec *TaskExecution) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // 检查是否终止
        if exec.State.IsTerminal() {
            return nil
        }

        // 获取处理器
        handler := lm.getHandler(exec.State)
        if handler == nil {
            return fmt.Errorf("no handler for state %s", exec.State)
        }

        // 进入状态
        if err := handler.Enter(ctx, exec); err != nil {
            return fmt.Errorf("enter state failed: %w", err)
        }

        // 执行状态逻辑
        transition, err := handler.Execute(ctx, exec)

        // 退出状态
        if exitErr := handler.Exit(ctx, exec); exitErr != nil {
            // 记录退出错误但不中断流程
        }

        if err != nil {
            // 执行失败，转换到失败状态
            exec.Error = &ExecutionError{
                Code:        "EXECUTION_ERROR",
                Message:     err.Error(),
                IsRetryable: isRetryableError(err),
            }

            if exec.Attempt < exec.MaxAttempts && exec.Error.IsRetryable {
                if err := lm.TransitionState(ctx, exec, TaskStateRetrying, "execution failed, will retry"); err != nil {
                    return err
                }
            } else {
                if err := lm.TransitionState(ctx, exec, TaskStateFailed, "execution failed"); err != nil {
                    return err
                }
            }
            continue
        }

        if transition == nil {
            // 没有转换，保持当前状态
            continue
        }

        // 执行状态转换
        if err := lm.TransitionState(ctx, exec, transition.To, transition.Reason); err != nil {
            return err
        }
    }
}

// CancelExecution 取消执行
func (lm *LifecycleManager) CancelExecution(ctx context.Context, execID string, reason string) error {
    exec, err := lm.store.Get(ctx, execID)
    if err != nil {
        return err
    }

    if exec.State.IsTerminal() {
        return fmt.Errorf("execution already in terminal state: %s", exec.State)
    }

    return lm.TransitionState(ctx, exec, TaskStateCancelled, reason)
}

// PauseExecution 暂停执行
func (lm *LifecycleManager) PauseExecution(ctx context.Context, execID string, reason string) error {
    exec, err := lm.store.Get(ctx, execID)
    if err != nil {
        return err
    }

    if exec.State != TaskStateScheduled {
        return fmt.Errorf("can only pause scheduled execution, current state: %s", exec.State)
    }

    return lm.TransitionState(ctx, exec, TaskStatePaused, reason)
}

// ResumeExecution 恢复执行
func (lm *LifecycleManager) ResumeExecution(ctx context.Context, execID string) error {
    exec, err := lm.store.Get(ctx, execID)
    if err != nil {
        return err
    }

    if exec.State != TaskStatePaused {
        return fmt.Errorf("can only resume paused execution, current state: %s", exec.State)
    }

    return lm.TransitionState(ctx, exec, TaskStateScheduled, "manual resume")
}

// GetExecutionHistory 获取执行历史
func (lm *LifecycleManager) GetExecutionHistory(ctx context.Context, taskID string) ([]*TaskExecution, error) {
    return lm.store.GetByTaskID(ctx, taskID)
}

// executeHooks 执行生命周期钩子
func (lm *LifecycleManager) executeHooks(ctx context.Context, exec *TaskExecution, from, to TaskState) error {
    lm.mu.RLock()
    hooks := lm.hooks[to]
    lm.mu.RUnlock()

    for _, hook := range hooks {
        if err := hook.OnStateTransition(ctx, exec, from, to); err != nil {
            return err
        }
    }

    return nil
}

// getHandler 获取状态处理器
func (lm *LifecycleManager) getHandler(state TaskState) StateHandler {
    lm.mu.RLock()
    defer lm.mu.RUnlock()
    return lm.handlers[state]
}

// registerDefaultHandlers 注册默认处理器
func (lm *LifecycleManager) registerDefaultHandlers() {
    // 待实现具体处理器
}

// ExecutionOptions 执行选项
type ExecutionOptions struct {
    MaxAttempts int
    Timeout     time.Duration
    TraceID     string
    Metadata    map[string]string
}

// Helper functions
func generateExecutionID() string {
    return fmt.Sprintf("exec-%d", time.Now().UnixNano())
}

func generateTraceID() string {
    return fmt.Sprintf("trace-%d", time.Now().UnixNano())
}

func generateSpanID() string {
    return fmt.Sprintf("span-%d", time.Now().UnixNano())
}

func isRetryableError(err error) bool {
    // 根据错误类型判断是否可重试
    // 网络错误、超时等通常可重试
    // 业务逻辑错误通常不可重试
    return true
}
```

---

## 钩子实现示例

```go
package lifecycle

import (
    "context"
    "fmt"

    "github.com/prometheus/client_golang/prometheus"
)

// MetricsHook 指标收集钩子
type MetricsHook struct {
    stateTransitions *prometheus.CounterVec
    executionDuration *prometheus.HistogramVec
}

func NewMetricsHook() *MetricsHook {
    return &MetricsHook{
        stateTransitions: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "task_state_transitions_total",
                Help: "Total number of task state transitions",
            },
            []string{"from_state", "to_state"},
        ),
        executionDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "task_execution_duration_seconds",
                Help: "Task execution duration",
            },
            []string{"final_state"},
        ),
    }
}

func (h *MetricsHook) OnStateTransition(ctx context.Context, exec *TaskExecution, from, to TaskState) error {
    h.stateTransitions.WithLabelValues(string(from), string(to)).Inc()

    // 记录执行时间
    if to.IsTerminal() && exec.StartedAt != nil {
        duration := time.Since(*exec.StartedAt)
        h.executionDuration.WithLabelValues(string(to)).Observe(duration.Seconds())
    }

    return nil
}

// NotificationHook 通知钩子
type NotificationHook struct {
    notifier Notifier
}

type Notifier interface {
    Notify(ctx context.Context, event string, exec *TaskExecution) error
}

func (h *NotificationHook) OnStateTransition(ctx context.Context, exec *TaskExecution, from, to TaskState) error {
    // 只在关键状态转换时通知
    switch to {
    case TaskStateFailed, TaskStateTimedOut, TaskStateCancelled:
        return h.notifier.Notify(ctx, string(to), exec)
    }
    return nil
}

// AuditLogHook 审计日志钩子
type AuditLogHook struct {
    logger Logger
}

type Logger interface {
    Info(msg string, fields ...interface{})
}

func (h *AuditLogHook) OnStateTransition(ctx context.Context, exec *TaskExecution, from, to TaskState) error {
    h.logger.Info("task state transition",
        "execution_id", exec.ID,
        "task_id", exec.TaskID,
        "from", from,
        "to", to,
        "timestamp", time.Now(),
    )
    return nil
}
```

---

## 状态处理器实现

```go
package lifecycle

import (
    "context"
    "time"
)

// RunningStateHandler 运行状态处理器
type RunningStateHandler struct {
    taskExecutor TaskExecutor
    cancelFunc   context.CancelFunc
}

type TaskExecutor interface {
    Execute(ctx context.Context, exec *TaskExecution) (map[string]interface{}, error)
}

func (h *RunningStateHandler) Enter(ctx context.Context, exec *TaskExecution) error {
    // 创建带超时的上下文
    execCtx, cancel := context.WithTimeout(ctx, exec.Timeout)
    h.cancelFunc = cancel
    _ = execCtx
    return nil
}

func (h *RunningStateHandler) Execute(ctx context.Context, exec *TaskExecution) (*Transition, error) {
    // 创建执行上下文
    execCtx, cancel := context.WithTimeout(ctx, exec.Timeout)
    defer cancel()

    // 监听取消信号
    done := make(chan struct{})
    var result map[string]interface{}
    var execErr error

    go func() {
        defer close(done)
        result, execErr = h.taskExecutor.Execute(execCtx, exec)
    }()

    select {
    case <-done:
        if execErr != nil {
            return nil, execErr
        }
        exec.Output = result
        return &Transition{
            To:     TaskStateCompleted,
            Reason: "execution completed successfully",
        }, nil

    case <-execCtx.Done():
        if execCtx.Err() == context.DeadlineExceeded {
            return &Transition{
                To:     TaskStateTimedOut,
                Reason: "execution timed out",
            }, nil
        }
        return &Transition{
            To:     TaskStateCancelled,
            Reason: "execution cancelled",
        }, nil
    }
}

func (h *RunningStateHandler) Exit(ctx context.Context, exec *TaskExecution) error {
    if h.cancelFunc != nil {
        h.cancelFunc()
    }
    return nil
}

// RetryingStateHandler 重试状态处理器
type RetryingStateHandler struct {
    backoffCalculator BackoffCalculator
}

type BackoffCalculator interface {
    Calculate(attempt int) time.Duration
}

func (h *RetryingStateHandler) Enter(ctx context.Context, exec *TaskExecution) error {
    return nil
}

func (h *RetryingStateHandler) Execute(ctx context.Context, exec *TaskExecution) (*Transition, error) {
    // 计算退避时间
    delay := h.backoffCalculator.Calculate(exec.Attempt)

    select {
    case <-time.After(delay):
        exec.Attempt++
        return &Transition{
            To:     TaskStateScheduled,
            Reason: fmt.Sprintf("retry attempt %d/%d", exec.Attempt, exec.MaxAttempts),
        }, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func (h *RetryingStateHandler) Exit(ctx context.Context, exec *TaskExecution) error {
    return nil
}
```

---

## 生产级使用示例

```go
package main

import (
    "context"
    "fmt"
    "time"

    "lifecycle"
)

func main() {
    // 创建存储
    store := NewPostgresExecutionStore("postgres://...")

    // 创建生命周期管理器
    lm := lifecycle.NewLifecycleManager(store)

    // 注册钩子
    lm.RegisterHook(lifecycle.TaskStateCompleted, lifecycle.NewMetricsHook())
    lm.RegisterHook(lifecycle.TaskStateFailed, lifecycle.NewNotificationHook(&EmailNotifier{}))

    // 注册处理器
    lm.RegisterHandler(lifecycle.TaskStateRunning, &lifecycle.RunningStateHandler{
        taskExecutor: &MyTaskExecutor{},
    })

    // 创建执行
    ctx := context.Background()
    exec, err := lm.CreateExecution(ctx, "task-123", map[string]interface{}{
        "data": "value",
    }, lifecycle.ExecutionOptions{
        MaxAttempts: 3,
        Timeout:     30 * time.Second,
    })
    if err != nil {
        panic(err)
    }

    // 转换到 ENQUEUED
    lm.TransitionState(ctx, exec, lifecycle.TaskStateEnqueued, "task created")

    // 转换到 SCHEDULED
    lm.TransitionState(ctx, exec, lifecycle.TaskStateScheduled, "worker picked up")

    // 执行状态机
    lm.Execute(ctx, exec)

    fmt.Printf("Execution completed with state: %s\n", exec.State)
}

type MyTaskExecutor struct{}

func (e *MyTaskExecutor) Execute(ctx context.Context, exec *lifecycle.TaskExecution) (map[string]interface{}, error) {
    // 实际任务逻辑
    time.Sleep(1 * time.Second)
    return map[string]interface{}{"result": "success"}, nil
}
```
