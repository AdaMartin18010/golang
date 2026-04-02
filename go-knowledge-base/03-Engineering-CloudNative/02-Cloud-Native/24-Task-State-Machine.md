# 任务状态机 (Task State Machine)

> **分类**: 工程与云原生  
> **标签**: #state-machine #task-lifecycle #workflow

---

## 状态定义

```go
type TaskStatus int

const (
    TaskStatusPending TaskStatus = iota
    TaskStatusScheduled
    TaskStatusRunning
    TaskStatusPaused
    TaskStatusSucceeded
    TaskStatusFailed
    TaskStatusCancelled
    TaskStatusRetrying
    TaskStatusTimeout
    TaskStatusSkipped
)

func (s TaskStatus) String() string {
    switch s {
    case TaskStatusPending:
        return "PENDING"
    case TaskStatusScheduled:
        return "SCHEDULED"
    case TaskStatusRunning:
        return "RUNNING"
    case TaskStatusPaused:
        return "PAUSED"
    case TaskStatusSucceeded:
        return "SUCCEEDED"
    case TaskStatusFailed:
        return "FAILED"
    case TaskStatusCancelled:
        return "CANCELLED"
    case TaskStatusRetrying:
        return "RETRYING"
    case TaskStatusTimeout:
        return "TIMEOUT"
    case TaskStatusSkipped:
        return "SKIPPED"
    default:
        return "UNKNOWN"
    }
}

// 状态转换规则
var validTransitions = map[TaskStatus][]TaskStatus{
    TaskStatusPending:    {TaskStatusScheduled, TaskStatusCancelled},
    TaskStatusScheduled:  {TaskStatusRunning, TaskStatusCancelled},
    TaskStatusRunning:    {TaskStatusSucceeded, TaskStatusFailed, TaskStatusPaused, TaskStatusTimeout, TaskStatusRetrying},
    TaskStatusPaused:     {TaskStatusRunning, TaskStatusCancelled},
    TaskStatusFailed:     {TaskStatusRetrying, TaskStatusCancelled},
    TaskStatusRetrying:   {TaskStatusRunning, TaskStatusFailed, TaskStatusCancelled},
    TaskStatusTimeout:    {TaskStatusRetrying, TaskStatusFailed, TaskStatusCancelled},
    TaskStatusSucceeded:  {}, // 终态
    TaskStatusCancelled:  {}, // 终态
    TaskStatusSkipped:    {}, // 终态
}

func CanTransition(from, to TaskStatus) bool {
    valid, ok := validTransitions[from]
    if !ok {
        return false
    }
    
    for _, v := range valid {
        if v == to {
            return true
        }
    }
    return false
}
```

---

## 状态机实现

```go
type TaskStateMachine struct {
    task       *Task
    transitions []StateTransition
    mu         sync.RWMutex
    observers  []StateObserver
}

type StateTransition struct {
    From      TaskStatus
    To        TaskStatus
    Timestamp time.Time
    Reason    string
    Operator  string
}

type StateObserver interface {
    OnStateChange(task *Task, from, to TaskStatus)
}

func (sm *TaskStateMachine) Transition(to TaskStatus, reason string) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    from := sm.task.Status
    
    // 验证转换
    if !CanTransition(from, to) {
        return fmt.Errorf("invalid transition from %s to %s", from, to)
    }
    
    // 执行转换
    transition := StateTransition{
        From:      from,
        To:        to,
        Timestamp: time.Now(),
        Reason:    reason,
    }
    
    sm.transitions = append(sm.transitions, transition)
    sm.task.Status = to
    
    // 通知观察者
    for _, obs := range sm.observers {
        obs.OnStateChange(sm.task, from, to)
    }
    
    return nil
}

// 持久化状态
func (sm *TaskStateMachine) Persist(ctx context.Context, store StateStore) error {
    return store.SaveTransition(ctx, sm.task.ID, StateTransitionRecord{
        TaskID:     sm.task.ID,
        From:       sm.task.Status,
        To:         sm.task.Status,
        Transitions: sm.transitions,
        UpdatedAt:  time.Now(),
    })
}
```

---

## 状态钩子

```go
type StateHooks struct {
    OnPending    func(*Task) error
    OnScheduled  func(*Task) error
    OnRunning    func(*Task) error
    OnSucceeded  func(*Task) error
    OnFailed     func(*Task) error
    OnCancelled  func(*Task) error
    OnRetrying   func(*Task) error
}

func (sm *TaskStateMachine) ExecuteWithHooks(to TaskStatus, hooks StateHooks) error {
    var hook func(*Task) error
    
    switch to {
    case TaskStatusRunning:
        hook = hooks.OnRunning
    case TaskStatusSucceeded:
        hook = hooks.OnSucceeded
    case TaskStatusFailed:
        hook = hooks.OnFailed
    }
    
    if hook != nil {
        if err := hook(sm.task); err != nil {
            return fmt.Errorf("hook failed: %w", err)
        }
    }
    
    return sm.Transition(to, "hook executed")
}
```

---

## 状态恢复

```go
func (sm *TaskStateMachine) Recover(ctx context.Context, store StateStore) error {
    record, err := store.GetLatestTransition(ctx, sm.task.ID)
    if err != nil {
        return err
    }
    
    // 恢复状态
    sm.task.Status = record.To
    sm.transitions = record.Transitions
    
    // 处理中断的任务
    switch sm.task.Status {
    case TaskStatusRunning:
        // 上次运行中断了，标记为失败并触发重试
        sm.Transition(TaskStatusFailed, "interrupted")
        if sm.task.RetryCount < sm.task.MaxRetries {
            sm.Transition(TaskStatusRetrying, "auto retry after recovery")
        }
    case TaskStatusRetrying:
        // 继续重试流程
        go sm.executeRetry()
    }
    
    return nil
}
```

---

## 状态查询

```go
func (sm *TaskStateMachine) GetHistory() []StateTransition {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    history := make([]StateTransition, len(sm.transitions))
    copy(history, sm.transitions)
    return history
}

func (sm *TaskStateMachine) TimeInState(status TaskStatus) time.Duration {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    var total time.Duration
    var enterTime time.Time
    
    for _, t := range sm.transitions {
        if t.To == status {
            enterTime = t.Timestamp
        }
        if t.From == status && !enterTime.IsZero() {
            total += t.Timestamp.Sub(enterTime)
            enterTime = time.Time{}
        }
    }
    
    // 如果当前就在该状态
    if sm.task.Status == status && !enterTime.IsZero() {
        total += time.Since(enterTime)
    }
    
    return total
}
```
