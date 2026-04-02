# 任务状态机实现 (Task State Machine Implementation)

> **分类**: 工程与云原生
> **标签**: #state-machine #workflow #task-lifecycle #event-driven
> **参考**: AWS Step Functions, Temporal, State Machine Cat

---

## 状态机模型

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Task State Machine                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────┐                                                                │
│  │  IDLE   │                                                                │
│  └────┬────┘                                                                │
│       │ create()                                                            │
│       ▼                                                                     │
│  ┌─────────┐                                                                │
│  │ PENDING │◄──────────────────────────────────────────┐                    │
│  └────┬────┘                                           │                    │
│       │ schedule()                                     │ retry()            │
│       ▼                                                │                    │
│  ┌─────────┐    timeout()    ┌─────────┐               │                    │
│  │SCHEDULED│────────────────▶│ TIMEOUT │──────────────┘                     │
│  └────┬────┘                  └─────────┘                                   │
│       │ start()                                                             │
│       ▼                                                                     │
│  ┌─────────┐    cancel()     ┌──────────┐                                   │
│  │ RUNNING │────────────────▶│CANCELLED │                                   │
│  └────┬────┘                  └──────────┘                                  │
│       │                                                                     │
│       ├────────────────┬────────────────┐                                   │
│       │                │                │                                   │
│       ▼                ▼                ▼                                   │
│  ┌─────────┐     ┌─────────┐     ┌─────────┐                                │
│  │COMPLETED│     │  FAILED │     │ PAUSED  │                                │
│  └─────────┘     └────┬────┘     └────┬────┘                                │
│                       │                │                                    │
│                       │ retry()        │ resume()                           │
│                       │                │                                    │
│                       └────────┬───────┘                                    │
│                                ▼                                            │
│                         ┌─────────────┐                                     │
│                         │   RETRYING  │                                     │
│                         └─────────────┘                                     │
│                                                                             │
│  State Transitions:                                                         │
│  • create:    IDLE → PENDING                                                │
│  • schedule:  PENDING → SCHEDULED                                           │
│  • start:     SCHEDULED → RUNNING                                           │
│  • complete:  RUNNING → COMPLETED                                           │
│  • fail:      RUNNING/RETRYING → FAILED                                     │
│  • cancel:    * → CANCELLED (except COMPLETED/FAILED)                       │
│  • pause:     RUNNING → PAUSED                                              │
│  • resume:    PAUSED → PENDING                                              │
│  • timeout:   SCHEDULED/RUNNING → TIMEOUT                                   │
│  • retry:     FAILED/TIMEOUT → RETRYING → PENDING                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 状态机核心实现

```go
package statemachine

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// State 状态类型
type State string

const (
    StateIdle       State = "idle"
    StatePending    State = "pending"
    StateScheduled  State = "scheduled"
    StateRunning    State = "running"
    StateCompleted  State = "completed"
    StateFailed     State = "failed"
    StateCancelled  State = "cancelled"
    StatePaused     State = "paused"
    StateTimeout    State = "timeout"
    StateRetrying   State = "retrying"
)

// Event 事件类型
type Event string

const (
    EventCreate   Event = "create"
    EventSchedule Event = "schedule"
    EventStart    Event = "start"
    EventComplete Event = "complete"
    EventFail     Event = "fail"
    EventCancel   Event = "cancel"
    EventPause    Event = "pause"
    EventResume   Event = "resume"
    EventTimeout  Event = "timeout"
    EventRetry    Event = "retry"
)

// Transition 状态转换
type Transition struct {
    From   State
    To     State
    Event  Event
    Guard  GuardFunc
    Action ActionFunc
}

type GuardFunc func(ctx context.Context, task *Task, data interface{}) bool
type ActionFunc func(ctx context.Context, task *Task, data interface{}) error

// StateMachine 状态机
type StateMachine struct {
    transitions map[State]map[Event]Transition
    handlers    map[Event][]HandlerFunc
    mu          sync.RWMutex
}

type HandlerFunc func(ctx context.Context, task *Task, from, to State, data interface{})

// NewStateMachine 创建状态机
func NewStateMachine() *StateMachine {
    sm := &StateMachine{
        transitions: make(map[State]map[Event]Transition),
        handlers:    make(map[Event][]HandlerFunc),
    }

    // 注册默认转换
    sm.registerDefaultTransitions()

    return sm
}

func (sm *StateMachine) registerDefaultTransitions() {
    // IDLE -> PENDING
    sm.RegisterTransition(StateIdle, StatePending, EventCreate, nil, nil)

    // PENDING -> SCHEDULED
    sm.RegisterTransition(StatePending, StateScheduled, EventSchedule, nil, nil)

    // PENDING -> CANCELLED
    sm.RegisterTransition(StatePending, StateCancelled, EventCancel, nil, nil)

    // SCHEDULED -> RUNNING
    sm.RegisterTransition(StateScheduled, StateRunning, EventStart, nil, nil)

    // SCHEDULED -> TIMEOUT
    sm.RegisterTransition(StateScheduled, StateTimeout, EventTimeout, nil, nil)

    // SCHEDULED -> CANCELLED
    sm.RegisterTransition(StateScheduled, StateCancelled, EventCancel, nil, nil)

    // RUNNING -> COMPLETED
    sm.RegisterTransition(StateRunning, StateCompleted, EventComplete, nil, nil)

    // RUNNING -> FAILED
    sm.RegisterTransition(StateRunning, StateFailed, EventFail, nil, nil)

    // RUNNING -> TIMEOUT
    sm.RegisterTransition(StateRunning, StateTimeout, EventTimeout, nil, nil)

    // RUNNING -> CANCELLED
    sm.RegisterTransition(StateRunning, StateCancelled, EventCancel, nil, nil)

    // RUNNING -> PAUSED
    sm.RegisterTransition(StateRunning, StatePaused, EventPause, nil, nil)

    // PAUSED -> PENDING
    sm.RegisterTransition(StatePaused, StatePending, EventResume, nil, nil)

    // FAILED -> RETRYING
    sm.RegisterTransition(StateFailed, StateRetrying, EventRetry,
        // Guard: 检查重试次数
        func(ctx context.Context, task *Task, data interface{}) bool {
            return task.RetryCount < task.MaxRetries
        },
        // Action: 增加重试计数
        func(ctx context.Context, task *Task, data interface{}) error {
            task.RetryCount++
            task.LastRetryAt = time.Now()
            return nil
        },
    )

    // RETRYING -> PENDING
    sm.RegisterTransition(StateRetrying, StatePending, EventSchedule, nil, nil)

    // TIMEOUT -> RETRYING
    sm.RegisterTransition(StateTimeout, StateRetrying, EventRetry,
        func(ctx context.Context, task *Task, data interface{}) bool {
            return task.RetryCount < task.MaxRetries
        }, nil,
    )
}

// RegisterTransition 注册状态转换
func (sm *StateMachine) RegisterTransition(from, to State, event Event,
    guard GuardFunc, action ActionFunc) {

    sm.mu.Lock()
    defer sm.mu.Unlock()

    if sm.transitions[from] == nil {
        sm.transitions[from] = make(map[Event]Transition)
    }

    sm.transitions[from][event] = Transition{
        From:   from,
        To:     to,
        Event:  event,
        Guard:  guard,
        Action: action,
    }
}

// CanTransition 检查是否可以转换
func (sm *StateMachine) CanTransition(task *Task, event Event) bool {
    sm.mu.RLock()
    defer sm.mu.RUnlock()

    transitions, ok := sm.transitions[task.State]
    if !ok {
        return false
    }

    transition, ok := transitions[event]
    if !ok {
        return false
    }

    if transition.Guard != nil {
        return transition.Guard(context.Background(), task, nil)
    }

    return true
}

// Transition 执行状态转换
func (sm *StateMachine) Transition(ctx context.Context, task *Task,
    event Event, data interface{}) error {

    sm.mu.Lock()
    defer sm.mu.Unlock()

    transitions, ok := sm.transitions[task.State]
    if !ok {
        return fmt.Errorf("no transitions defined for state %s", task.State)
    }

    transition, ok := transitions[event]
    if !ok {
        return fmt.Errorf("event %s not allowed in state %s", event, task.State)
    }

    // 检查 Guard
    if transition.Guard != nil && !transition.Guard(ctx, task, data) {
        return fmt.Errorf("guard condition not met for transition %s -> %s",
            task.State, transition.To)
    }

    fromState := task.State

    // 执行 Action
    if transition.Action != nil {
        if err := transition.Action(ctx, task, data); err != nil {
            return fmt.Errorf("action failed: %w", err)
        }
    }

    // 更新状态
    task.State = transition.To
    task.UpdatedAt = time.Now()

    // 记录历史
    task.History = append(task.History, StateHistory{
        From:      fromState,
        To:        transition.To,
        Event:     event,
        Timestamp: time.Now(),
        Data:      data,
    })

    // 触发事件处理器
    sm.triggerHandlers(ctx, task, fromState, transition.To, event, data)

    return nil
}

func (sm *StateMachine) triggerHandlers(ctx context.Context, task *Task,
    from, to State, event Event, data interface{}) {

    handlers := sm.handlers[event]
    for _, handler := range handlers {
        go handler(ctx, task, from, to, data)
    }
}

// On 注册事件处理器
func (sm *StateMachine) On(event Event, handler HandlerFunc) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sm.handlers[event] = append(sm.handlers[event], handler)
}
```

---

## 任务结构定义

```go
// Task 任务定义
type Task struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    State     State                  `json:"state"`
    Payload   map[string]interface{} `json:"payload"`

    // 重试配置
    RetryCount  int       `json:"retry_count"`
    MaxRetries  int       `json:"max_retries"`
    LastRetryAt time.Time `json:"last_retry_at,omitempty"`

    // 时间戳
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    ScheduledAt *time.Time     `json:"scheduled_at,omitempty"`
    StartedAt   *time.Time     `json:"started_at,omitempty"`
    CompletedAt *time.Time     `json:"completed_at,omitempty"`

    // 状态历史
    History []StateHistory `json:"history"`

    // 错误信息
    LastError string `json:"last_error,omitempty"`
}

// StateHistory 状态历史
type StateHistory struct {
    From      State       `json:"from"`
    To        State       `json:"to"`
    Event     Event       `json:"event"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data,omitempty"`
}

// TaskManager 任务管理器
type TaskManager struct {
    store StateMachineStore
    sm    *StateMachine
}

// CreateTask 创建任务
func (tm *TaskManager) CreateTask(ctx context.Context, taskType string,
    payload map[string]interface{}) (*Task, error) {

    task := &Task{
        ID:         generateID(),
        Type:       taskType,
        State:      StateIdle,
        Payload:    payload,
        MaxRetries: 3,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    // 触发创建事件
    if err := tm.sm.Transition(ctx, task, EventCreate, nil); err != nil {
        return nil, err
    }

    // 持久化
    if err := tm.store.Save(ctx, task); err != nil {
        return nil, err
    }

    return task, nil
}

// ExecuteTask 执行任务
func (tm *TaskManager) ExecuteTask(ctx context.Context, taskID string,
    executor TaskExecutor) error {

    task, err := tm.store.Get(ctx, taskID)
    if err != nil {
        return err
    }

    // 调度任务
    if err := tm.sm.Transition(ctx, task, EventSchedule, nil); err != nil {
        return err
    }

    now := time.Now()
    task.ScheduledAt = &now

    // 启动任务
    if err := tm.sm.Transition(ctx, task, EventStart, nil); err != nil {
        return err
    }

    task.StartedAt = &now

    // 保存状态
    tm.store.Save(ctx, task)

    // 异步执行
    go tm.runTask(context.Background(), task, executor)

    return nil
}

func (tm *TaskManager) runTask(ctx context.Context, task *Task,
    executor TaskExecutor) {

    // 执行任务
    result, err := executor.Execute(ctx, task)

    if err != nil {
        // 执行失败
        task.LastError = err.Error()

        // 检查是否需要重试
        if task.RetryCount < task.MaxRetries {
            tm.sm.Transition(ctx, task, EventRetry, result)
            tm.store.Save(ctx, task)

            // 延迟重试
            time.Sleep(calculateBackoff(task.RetryCount))
            tm.ExecuteTask(ctx, task.ID, executor)
        } else {
            tm.sm.Transition(ctx, task, EventFail, result)
            tm.store.Save(ctx, task)
        }
    } else {
        // 执行成功
        tm.sm.Transition(ctx, task, EventComplete, result)
        now := time.Now()
        task.CompletedAt = &now
        tm.store.Save(ctx, task)
    }
}

func calculateBackoff(retryCount int) time.Duration {
    // 指数退避
    baseDelay := time.Second
    maxDelay := time.Minute * 5

    delay := baseDelay * time.Duration(1<<uint(retryCount))
    if delay > maxDelay {
        delay = maxDelay
    }

    // 添加 jitter
    jitter := time.Duration(rand.Int63n(int64(delay) / 2))
    return delay + jitter
}
```

---

## 事件驱动工作流

```go
// EventDrivenWorkflow 事件驱动工作流
type EventDrivenWorkflow struct {
    sm       *StateMachine
    events   chan WorkflowEvent
    handlers map[string]EventHandler
}

type WorkflowEvent struct {
    Type      string
    TaskID    string
    Data      interface{}
    Timestamp time.Time
}

type EventHandler func(ctx context.Context, task *Task, event WorkflowEvent) error

// Start 启动工作流引擎
func (w *EventDrivenWorkflow) Start(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case event := <-w.events:
            w.handleEvent(ctx, event)
        }
    }
}

func (w *EventDrivenWorkflow) handleEvent(ctx context.Context, event WorkflowEvent) {
    // 获取任务
    task, err := w.getTask(ctx, event.TaskID)
    if err != nil {
        log.Printf("Failed to get task: %v", err)
        return
    }

    // 查找事件处理器
    handler, ok := w.handlers[event.Type]
    if !ok {
        log.Printf("No handler for event type: %s", event.Type)
        return
    }

    // 执行处理器
    if err := handler(ctx, task, event); err != nil {
        log.Printf("Event handler failed: %v", err)
        return
    }

    // 检查是否需要自动推进状态
    w.autoAdvance(ctx, task)
}

func (w *EventDrivenWorkflow) autoAdvance(ctx context.Context, task *Task) {
    // 根据当前状态和条件自动触发状态转换
    switch task.State {
    case StatePending:
        // 检查是否可以调度
        if w.canSchedule(task) {
            w.sm.Transition(ctx, task, EventSchedule, nil)
        }

    case StateScheduled:
        // 检查是否有可用 Worker
        if w.hasAvailableWorker() {
            w.sm.Transition(ctx, task, EventStart, nil)
        }

    case StateFailed:
        // 检查是否需要重试
        if task.RetryCount < task.MaxRetries {
            w.sm.Transition(ctx, task, EventRetry, nil)
        }
    }

    w.saveTask(ctx, task)
}

// RegisterEventHandler 注册事件处理器
func (w *EventDrivenWorkflow) RegisterEventHandler(eventType string,
    handler EventHandler) {
    w.handlers[eventType] = handler
}

// EmitEvent 发送事件
func (w *EventDrivenWorkflow) EmitEvent(event WorkflowEvent) {
    w.events <- event
}
```

---

## 持久化存储

```go
// StateMachineStore 状态机存储接口
type StateMachineStore interface {
    Save(ctx context.Context, task *Task) error
    Get(ctx context.Context, taskID string) (*Task, error)
    Update(ctx context.Context, taskID string, fn func(*Task) error) error
    ListByState(ctx context.Context, state State) ([]*Task, error)
    Watch(ctx context.Context, state State) (<-chan *Task, error)
}

// PostgresStore PostgreSQL 实现
type PostgresStore struct {
    db *sql.DB
}

func (s *PostgresStore) Save(ctx context.Context, task *Task) error {
    query := `
        INSERT INTO tasks (id, type, state, payload, retry_count, max_retries,
            created_at, updated_at, history)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (id) DO UPDATE SET
            state = EXCLUDED.state,
            payload = EXCLUDED.payload,
            retry_count = EXCLUDED.retry_count,
            updated_at = EXCLUDED.updated_at,
            history = EXCLUDED.history
    `

    historyJSON, _ := json.Marshal(task.History)

    _, err := s.db.ExecContext(ctx, query,
        task.ID, task.Type, task.State, task.Payload,
        task.RetryCount, task.MaxRetries,
        task.CreatedAt, task.UpdatedAt, historyJSON,
    )

    return err
}

func (s *PostgresStore) Update(ctx context.Context, taskID string,
    fn func(*Task) error) error {

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 获取并锁定任务
    task, err := s.getForUpdate(ctx, tx, taskID)
    if err != nil {
        return err
    }

    // 执行更新
    if err := fn(task); err != nil {
        return err
    }

    // 保存更新
    if err := s.saveInTx(ctx, tx, task); err != nil {
        return err
    }

    return tx.Commit()
}

// EventSourcingStore 事件溯源存储
type EventSourcingStore struct {
    eventStore EventStore
    projector  StateProjector
}

type TaskEvent struct {
    EventID     string          `json:"event_id"`
    TaskID      string          `json:"task_id"`
    EventType   string          `json:"event_type"`
    FromState   State           `json:"from_state"`
    ToState     State           `json:"to_state"`
    Data        interface{}     `json:"data"`
    Timestamp   time.Time       `json:"timestamp"`
    Version     int             `json:"version"`
}

func (s *EventSourcingStore) Save(ctx context.Context, task *Task) error {
    // 获取最新事件版本
    version := s.eventStore.GetLatestVersion(ctx, task.ID)

    // 创建事件
    lastHistory := task.History[len(task.History)-1]
    event := &TaskEvent{
        EventID:   generateID(),
        TaskID:    task.ID,
        EventType: string(lastHistory.Event),
        FromState: lastHistory.From,
        ToState:   lastHistory.To,
        Data:      lastHistory.Data,
        Timestamp: lastHistory.Timestamp,
        Version:   version + 1,
    }

    // 保存事件
    if err := s.eventStore.Append(ctx, event); err != nil {
        return err
    }

    // 更新投影
    return s.projector.Project(ctx, task)
}

func (s *EventSourcingStore) Get(ctx context.Context, taskID string) (*Task, error) {
    // 重放所有事件重建状态
    events, err := s.eventStore.GetEvents(ctx, taskID)
    if err != nil {
        return nil, err
    }

    task := &Task{ID: taskID}

    for _, event := range events {
        task.State = event.ToState
        task.History = append(task.History, StateHistory{
            From:      event.FromState,
            To:        event.ToState,
            Event:     Event(event.EventType),
            Timestamp: event.Timestamp,
            Data:      event.Data,
        })
    }

    return task, nil
}
```
