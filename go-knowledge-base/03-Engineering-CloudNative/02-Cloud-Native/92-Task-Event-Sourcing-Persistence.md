# 任务事件溯源持久化 (Task Event Sourcing Persistence)

> **分类**: 工程与云原生
> **标签**: #event-sourcing #cqrs #persistence
> **参考**: Event Store, CQRS Pattern

---

## 事件溯源架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Event Sourcing Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Write Model (Command Side)                        │   │
│  │                                                                      │   │
│  │   Command ──► Validate ──► Generate Event ──► Append to Event Store  │   │
│  │                                                              │       │   │
│  │                                                              ▼       │   │
│  │                                                      ┌──────────┐   │    │
│  │                                                      │  Event   │   │    │
│  │                                                      │  Store   │   │    │
│  │                                                      │(Append-  │   │    │
│  │                                                      │ Only Log)│   │    │
│  │                                                      └──────────┘   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│                                    │ Projections                            │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Read Model (Query Side)                           │   │
│  │                                                                      │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐             │   │
│  │   │ Current  │  │  Task    │  │  Task    │  │  Audit   │             │   │
│  │   │  State   │  │ History  │  │ Metrics  │  │   Log    │             │   │
│  │   └──────────┘  └──────────┘  └──────────┘  └──────────┘             │   │
│  │                                                                      │   │
│  │   Projections built from event stream                                │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Event Store Schema                                │   │
│  │                                                                      │   │
│  │   ┌──────────┬──────────┬──────────┬──────────┬─────────────────┐    │   │
│  │   │Event ID  │Stream ID │Event Type│Version   │Payload          │    │   │
│  │   ├──────────┼──────────┼──────────┼──────────┼─────────────────┤    │   │
│  │   │UUID      │task-123  │Created   │1         │{...}            │    │   │
│  │   │UUID      │task-123  │Started   │2         │{worker: "A"}    │    │   │
│  │   │UUID      │task-123  │Completed │3         │{result: {...}}  │    │   │
│  │   └──────────┴──────────┴──────────┴──────────┴─────────────────┘    │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整事件溯源实现

```go
package eventsource

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

// Event 领域事件
type Event struct {
    ID          string          `json:"id"`
    StreamID    string          `json:"stream_id"`
    StreamType  string          `json:"stream_type"`
    EventType   string          `json:"event_type"`
    Version     int             `json:"version"`
    Payload     json.RawMessage `json:"payload"`
    Metadata    Metadata        `json:"metadata"`
    Timestamp   time.Time       `json:"timestamp"`
}

// Metadata 事件元数据
type Metadata struct {
    CorrelationID string            `json:"correlation_id"`
    CausationID   string            `json:"causation_id"`
    UserID        string            `json:"user_id"`
    Extra         map[string]string `json:"extra,omitempty"`
}

// EventStore 事件存储接口
type EventStore interface {
    Append(ctx context.Context, events []Event) error
    Read(ctx context.Context, streamID string, fromVersion int) ([]Event, error)
    ReadAll(ctx context.Context, position int64, count int) ([]Event, error)
    GetStreamVersion(ctx context.Context, streamID string) (int, error)
}

// TaskEvents 任务事件类型
const (
    TaskCreatedEvent     = "TaskCreated"
    TaskScheduledEvent   = "TaskScheduled"
    TaskStartedEvent     = "TaskStarted"
    TaskCompletedEvent   = "TaskCompleted"
    TaskFailedEvent      = "TaskFailed"
    TaskCancelledEvent   = "TaskCancelled"
    TaskRetriedEvent     = "TaskRetried"
)

// TaskCreated 任务创建事件
type TaskCreated struct {
    TaskID      string                 `json:"task_id"`
    Type        string                 `json:"type"`
    Payload     map[string]interface{} `json:"payload"`
    Priority    int                    `json:"priority"`
    ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

// TaskStarted 任务开始事件
type TaskStarted struct {
    TaskID   string    `json:"task_id"`
    WorkerID string    `json:"worker_id"`
    StartedAt time.Time `json:"started_at"`
}

// TaskCompleted 任务完成事件
type TaskCompleted struct {
    TaskID    string                 `json:"task_id"`
    Result    map[string]interface{} `json:"result"`
    Duration  int64                  `json:"duration_ms"`
}

// TaskFailed 任务失败事件
type TaskFailed struct {
    TaskID      string `json:"task_id"`
    Error       string `json:"error"`
    ErrorCode   string `json:"error_code"`
    Retryable   bool   `json:"retryable"`
    RetryCount  int    `json:"retry_count"`
}

// Aggregate 聚合根接口
type Aggregate interface {
    Apply(event Event) error
    GetUncommittedEvents() []Event
    ClearUncommittedEvents()
    GetVersion() int
}

// TaskAggregate 任务聚合根
type TaskAggregate struct {
    ID             string
    Version        int
    State          string
    Type           string
    Payload        map[string]interface{}
    WorkerID       string
    StartedAt      *time.Time
    CompletedAt    *time.Time
    RetryCount     int

    uncommittedEvents []Event
}

// NewTaskAggregate 创建任务聚合根
func NewTaskAggregate(id string) *TaskAggregate {
    return &TaskAggregate{
        ID:      id,
        State:   "pending",
        Payload: make(map[string]interface{}),
    }
}

// Create 创建任务
func (ta *TaskAggregate) Create(taskType string, payload map[string]interface{}, priority int) error {
    if ta.Version > 0 {
        return fmt.Errorf("task already exists")
    }

    event := Event{
        StreamID:   ta.ID,
        StreamType: "Task",
        EventType:  TaskCreatedEvent,
        Version:    ta.Version + 1,
        Timestamp:  time.Now(),
    }

    data, _ := json.Marshal(TaskCreated{
        TaskID:   ta.ID,
        Type:     taskType,
        Payload:  payload,
        Priority: priority,
    })
    event.Payload = data

    ta.uncommittedEvents = append(ta.uncommittedEvents, event)
    return ta.Apply(event)
}

// Start 开始任务
func (ta *TaskAggregate) Start(workerID string) error {
    if ta.State != "pending" && ta.State != "scheduled" {
        return fmt.Errorf("cannot start task in state: %s", ta.State)
    }

    event := Event{
        StreamID:   ta.ID,
        StreamType: "Task",
        EventType:  TaskStartedEvent,
        Version:    ta.Version + 1,
        Timestamp:  time.Now(),
    }

    data, _ := json.Marshal(TaskStarted{
        TaskID:    ta.ID,
        WorkerID:  workerID,
        StartedAt: time.Now(),
    })
    event.Payload = data

    ta.uncommittedEvents = append(ta.uncommittedEvents, event)
    return ta.Apply(event)
}

// Complete 完成任务
func (ta *TaskAggregate) Complete(result map[string]interface{}, duration int64) error {
    if ta.State != "running" {
        return fmt.Errorf("cannot complete task in state: %s", ta.State)
    }

    event := Event{
        StreamID:   ta.ID,
        StreamType: "Task",
        EventType:  TaskCompletedEvent,
        Version:    ta.Version + 1,
        Timestamp:  time.Now(),
    }

    data, _ := json.Marshal(TaskCompleted{
        TaskID:   ta.ID,
        Result:   result,
        Duration: duration,
    })
    event.Payload = data

    ta.uncommittedEvents = append(ta.uncommittedEvents, event)
    return ta.Apply(event)
}

// Fail 任务失败
func (ta *TaskAggregate) Fail(err error, errorCode string, retryable bool) error {
    event := Event{
        StreamID:   ta.ID,
        StreamType: "Task",
        EventType:  TaskFailedEvent,
        Version:    ta.Version + 1,
        Timestamp:  time.Now(),
    }

    data, _ := json.Marshal(TaskFailed{
        TaskID:     ta.ID,
        Error:      err.Error(),
        ErrorCode:  errorCode,
        Retryable:  retryable,
        RetryCount: ta.RetryCount,
    })
    event.Payload = data

    ta.uncommittedEvents = append(ta.uncommittedEvents, event)
    return ta.Apply(event)
}

// Apply 应用事件
func (ta *TaskAggregate) Apply(event Event) error {
    switch event.EventType {
    case TaskCreatedEvent:
        var data TaskCreated
        if err := json.Unmarshal(event.Payload, &data); err != nil {
            return err
        }
        ta.Type = data.Type
        ta.Payload = data.Payload
        ta.State = "created"

    case TaskStartedEvent:
        var data TaskStarted
        if err := json.Unmarshal(event.Payload, &data); err != nil {
            return err
        }
        ta.WorkerID = data.WorkerID
        ta.StartedAt = &data.StartedAt
        ta.State = "running"

    case TaskCompletedEvent:
        ta.State = "completed"
        now := time.Now()
        ta.CompletedAt = &now

    case TaskFailedEvent:
        var data TaskFailed
        if err := json.Unmarshal(event.Payload, &data); err != nil {
            return err
        }
        ta.RetryCount = data.RetryCount
        if data.Retryable {
            ta.State = "failed_retryable"
        } else {
            ta.State = "failed"
        }
    }

    ta.Version = event.Version
    return nil
}

// GetUncommittedEvents 获取未提交事件
func (ta *TaskAggregate) GetUncommittedEvents() []Event {
    return ta.uncommittedEvents
}

// ClearUncommittedEvents 清空未提交事件
func (ta *TaskAggregate) ClearUncommittedEvents() {
    ta.uncommittedEvents = nil
}

// GetVersion 获取版本
func (ta *TaskAggregate) GetVersion() int {
    return ta.Version
}

// Repository 仓储
type Repository struct {
    eventStore EventStore
}

// NewRepository 创建仓储
func NewRepository(es EventStore) *Repository {
    return &Repository{eventStore: es}
}

// Load 加载聚合根
func (r *Repository) Load(ctx context.Context, streamID string) (*TaskAggregate, error) {
    events, err := r.eventStore.Read(ctx, streamID, 0)
    if err != nil {
        return nil, err
    }

    if len(events) == 0 {
        return nil, fmt.Errorf("stream not found: %s", streamID)
    }

    aggregate := NewTaskAggregate(streamID)
    for _, event := range events {
        if err := aggregate.Apply(event); err != nil {
            return nil, err
        }
    }

    // 清除未提交事件（从存储加载的）
    aggregate.ClearUncommittedEvents()

    return aggregate, nil
}

// Save 保存聚合根
func (r *Repository) Save(ctx context.Context, aggregate *TaskAggregate) error {
    events := aggregate.GetUncommittedEvents()
    if len(events) == 0 {
        return nil
    }

    // 乐观并发控制
    currentVersion, err := r.eventStore.GetStreamVersion(ctx, aggregate.ID)
    if err != nil {
        return err
    }

    if currentVersion != aggregate.GetVersion()-len(events) {
        return fmt.Errorf("concurrency conflict: expected version %d, got %d",
            currentVersion, aggregate.GetVersion()-len(events))
    }

    // 追加事件
    if err := r.eventStore.Append(ctx, events); err != nil {
        return err
    }

    aggregate.ClearUncommittedEvents()
    return nil
}

// Projection 投影接口
type Projection interface {
    Handle(event Event) error
    GetState(streamID string) (interface{}, error)
}

// TaskProjection 任务投影
type TaskProjection struct {
    tasks map[string]*TaskReadModel
}

// TaskReadModel 任务读模型
type TaskReadModel struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    State       string                 `json:"state"`
    WorkerID    string                 `json:"worker_id,omitempty"`
    CreatedAt   time.Time              `json:"created_at"`
    StartedAt   *time.Time             `json:"started_at,omitempty"`
    CompletedAt *time.Time             `json:"completed_at,omitempty"`
    Payload     map[string]interface{} `json:"payload"`
    Result      map[string]interface{} `json:"result,omitempty"`
}

// NewTaskProjection 创建任务投影
func NewTaskProjection() *TaskProjection {
    return &TaskProjection{
        tasks: make(map[string]*TaskReadModel),
    }
}

// Handle 处理事件
func (tp *TaskProjection) Handle(event Event) error {
    switch event.EventType {
    case TaskCreatedEvent:
        var data TaskCreated
        if err := json.Unmarshal(event.Payload, &data); err != nil {
            return err
        }
        tp.tasks[event.StreamID] = &TaskReadModel{
            ID:        data.TaskID,
            Type:      data.Type,
            State:     "pending",
            CreatedAt: event.Timestamp,
            Payload:   data.Payload,
        }

    case TaskStartedEvent:
        var data TaskStarted
        if err := json.Unmarshal(event.Payload, &data); err != nil {
            return err
        }
        if task, ok := tp.tasks[event.StreamID]; ok {
            task.State = "running"
            task.WorkerID = data.WorkerID
            task.StartedAt = &data.StartedAt
        }

    case TaskCompletedEvent:
        var data TaskCompleted
        if err := json.Unmarshal(event.Payload, &data); err != nil {
            return err
        }
        if task, ok := tp.tasks[event.StreamID]; ok {
            task.State = "completed"
            task.Result = data.Result
            now := time.Now()
            task.CompletedAt = &now
        }
    }

    return nil
}

// GetState 获取状态
func (tp *TaskProjection) GetState(streamID string) (interface{}, error) {
    task, ok := tp.tasks[streamID]
    if !ok {
        return nil, fmt.Errorf("task not found: %s", streamID)
    }
    return task, nil
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"

    "eventsource"
)

func main() {
    // 创建事件存储
    eventStore := NewPostgresEventStore("postgres://...")

    // 创建仓储
    repo := eventsource.NewRepository(eventStore)

    // 创建任务
    task := eventsource.NewTaskAggregate("task-123")

    ctx := context.Background()

    // 执行命令
    task.Create("send-email", map[string]interface{}{
        "to":      "user@example.com",
        "subject": "Hello",
    }, 1)

    // 保存
    if err := repo.Save(ctx, task); err != nil {
        panic(err)
    }

    // 加载任务
    loadedTask, err := repo.Load(ctx, "task-123")
    if err != nil {
        panic(err)
    }

    // 开始任务
    loadedTask.Start("worker-1")

    // 保存
    repo.Save(ctx, loadedTask)

    // 完成任务
    loadedTask.Complete(map[string]string{"status": "sent"}, 1000)
    repo.Save(ctx, loadedTask)

    fmt.Printf("Task state: %s\n", loadedTask.State)
}
```
