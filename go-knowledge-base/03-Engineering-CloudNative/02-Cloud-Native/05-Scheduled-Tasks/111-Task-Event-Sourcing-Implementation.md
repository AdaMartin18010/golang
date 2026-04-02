# 任务事件溯源实现 (Task Event Sourcing Implementation)

> **分类**: 工程与云原生
> **标签**: #event-sourcing #cqrs #event-store #audit
> **参考**: EventStoreDB, Axon Framework, Martin Fowler Event Sourcing

---

## 事件溯源核心概念

```
传统CRUD:                 事件溯源:
┌──────────┐             ┌──────────┐
│   Task   │             │  Events  │
│  (状态)   │             │(不可变序列)│
├──────────┤             ├──────────┤
│ status   │             │ Created  │ ─┐
│ retry    │  ← 问题:    │ Scheduled│  │
│ worker   │   丢失历史   │ Started  │  ├ 可重建
│ result   │             │ Retried  │  │  任意状态
│ ...      │             │ Completed│ ─┘
└──────────┘             └──────────┘
```

---

## 完整事件存储实现

```go
package eventsourcing

import (
 "context"
 "encoding/json"
 "fmt"
 "sync"
 "time"

 "github.com/google/uuid"
)

// Event 领域事件接口
type Event interface {
 EventID() string
 EventType() string
 AggregateID() string
 AggregateType() string
 EventVersion() int
 OccurredAt() time.Time
 Metadata() map[string]string
}

// BaseEvent 事件基类
type BaseEvent struct {
 ID            string
 Type          string
 AggregateID   string
 AggregateType string
 Version       int
 Timestamp     time.Time
 Meta          map[string]string
 Payload       json.RawMessage
}

func (e BaseEvent) EventID() string            { return e.ID }
func (e BaseEvent) EventType() string          { return e.Type }
func (e BaseEvent) AggregateID() string        { return e.AggregateID }
func (e BaseEvent) AggregateType() string      { return e.AggregateType }
func (e BaseEvent) EventVersion() int          { return e.Version }
func (e BaseEvent) OccurredAt() time.Time     { return e.Timestamp }
func (e BaseEvent) Metadata() map[string]string { return e.Meta }

// TaskCreatedEvent 任务创建事件
type TaskCreatedEvent struct {
 BaseEvent
 TaskType     string
 Payload      []byte
 Priority     int
 ScheduledAt  *time.Time
 MaxRetries   int
 Timeout      time.Duration
}

// TaskScheduledEvent 任务调度事件
type TaskScheduledEvent struct {
 BaseEvent
 WorkerID    string
 ScheduledAt time.Time
}

// TaskStartedEvent 任务开始事件
type TaskStartedEvent struct {
 BaseEvent
 WorkerID  string
 StartedAt time.Time
}

// TaskCompletedEvent 任务完成事件
type TaskCompletedEvent struct {
 BaseEvent
 WorkerID    string
 Result      []byte
 Duration    time.Duration
 CompletedAt time.Time
}

// TaskFailedEvent 任务失败事件
type TaskFailedEvent struct {
 BaseEvent
 WorkerID  string
 Error     string
 Retryable bool
 Attempt   int
}

// EventStore 事件存储接口
type EventStore interface {
 // 写入事件（原子性）
 Append(ctx context.Context, events []Event, expectedVersion int) error

 // 读取聚合事件
 GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)

 // 按类型读取
 GetEventsByType(ctx context.Context, eventType string, afterPosition int64, limit int) ([]Event, int64, error)

 // 获取当前版本
 GetVersion(ctx context.Context, aggregateID string) (int, error)

 // 快照支持
 SaveSnapshot(ctx context.Context, aggregateID string, version int, snapshot interface{}) error
 GetSnapshot(ctx context.Context, aggregateID string) (int, interface{}, error)
}

// Aggregate 聚合根接口
type Aggregate interface {
 ID() string
 Type() string
 Version() int
 Apply(event Event) error
 UncommittedEvents() []Event
 MarkCommitted()
}

// TaskAggregate 任务聚合根
type TaskAggregate struct {
 id      string
 version int

 // 状态
 status      TaskStatus
 workerID    string
 retryCount  int
 result      []byte

 // 未提交事件
 uncommitted []Event
 mu          sync.RWMutex
}

// NewTaskAggregate 创建任务聚合
func NewTaskAggregate(id string) *TaskAggregate {
 return &TaskAggregate{
  id:          id,
  status:      TaskStatusPending,
  uncommitted: make([]Event, 0),
 }
}

func (t *TaskAggregate) ID() string           { return t.id }
func (t *TaskAggregate) Type() string         { return "Task" }
func (t *TaskAggregate) Version() int         { return t.version }
func (t *TaskAggregate) UncommittedEvents() []Event {
 t.mu.RLock()
 defer t.mu.RUnlock()
 return append([]Event{}, t.uncommitted...)
}
func (t *TaskAggregate) MarkCommitted() {
 t.mu.Lock()
 defer t.mu.Unlock()
 t.uncommitted = t.uncommitted[:0]
}

// Create 创建任务
func (t *TaskAggregate) Create(taskType string, payload []byte, priority int) error {
 if t.status != TaskStatusNone {
  return fmt.Errorf("task already exists")
 }

 event := &TaskCreatedEvent{
  BaseEvent: BaseEvent{
   ID:            uuid.New().String(),
   Type:          "TaskCreated",
   AggregateID:   t.id,
   AggregateType: "Task",
   Version:       t.version + 1,
   Timestamp:     time.Now(),
  },
  TaskType: taskType,
  Payload:  payload,
  Priority: priority,
 }

 return t.apply(event)
}

// Schedule 调度任务
func (t *TaskAggregate) Schedule(workerID string) error {
 if t.status != TaskStatusPending {
  return fmt.Errorf("cannot schedule task in status %s", t.status)
 }

 event := &TaskScheduledEvent{
  BaseEvent: BaseEvent{
   ID:            uuid.New().String(),
   Type:          "TaskScheduled",
   AggregateID:   t.id,
   AggregateType: "Task",
   Version:       t.version + 1,
   Timestamp:     time.Now(),
  },
  WorkerID:    workerID,
  ScheduledAt: time.Now(),
 }

 return t.apply(event)
}

// Complete 完成任务
func (t *TaskAggregate) Complete(result []byte, duration time.Duration) error {
 if t.status != TaskStatusRunning {
  return fmt.Errorf("cannot complete task in status %s", t.status)
 }

 event := &TaskCompletedEvent{
  BaseEvent: BaseEvent{
   ID:            uuid.New().String(),
   Type:          "TaskCompleted",
   AggregateID:   t.id,
   AggregateType: "Task",
   Version:       t.version + 1,
   Timestamp:     time.Now(),
  },
  Result:      result,
  Duration:    duration,
  CompletedAt: time.Now(),
 }

 return t.apply(event)
}

// Fail 标记失败
func (t *TaskAggregate) Fail(err error, retryable bool) error {
 event := &TaskFailedEvent{
  BaseEvent: BaseEvent{
   ID:            uuid.New().String(),
   Type:          "TaskFailed",
   AggregateID:   t.id,
   AggregateType: "Task",
   Version:       t.version + 1,
   Timestamp:     time.Now(),
  },
  Error:     err.Error(),
  Retryable: retryable,
  Attempt:   t.retryCount + 1,
 }

 return t.apply(event)
}

// apply 应用事件（内部）
func (t *TaskAggregate) apply(event Event) error {
 t.mu.Lock()
 defer t.mu.Unlock()

 // 应用到状态
 if err := t.applyEvent(event); err != nil {
  return err
 }

 // 记录未提交事件
 t.uncommitted = append(t.uncommitted, event)
 t.version = event.EventVersion()

 return nil
}

// applyEvent 状态转换
func (t *TaskAggregate) applyEvent(event Event) error {
 switch e := event.(type) {
 case *TaskCreatedEvent:
  t.status = TaskStatusPending

 case *TaskScheduledEvent:
  t.status = TaskStatusScheduled
  t.workerID = e.WorkerID

 case *TaskStartedEvent:
  t.status = TaskStatusRunning

 case *TaskCompletedEvent:
  t.status = TaskStatusCompleted
  t.result = e.Result

 case *TaskFailedEvent:
  t.retryCount++
  if e.Retryable && t.retryCount < 3 {
   t.status = TaskStatusPending
  } else {
   t.status = TaskStatusFailed
  }

 default:
  return fmt.Errorf("unknown event type: %s", event.EventType())
 }

 return nil
}

// LoadFromHistory 从历史事件重建聚合
func (t *TaskAggregate) LoadFromHistory(events []Event) error {
 for _, event := range events {
  if err := t.applyEvent(event); err != nil {
   return err
  }
  t.version = event.EventVersion()
 }
 return nil
}

// TaskStatus 任务状态
type TaskStatus int

const (
 TaskStatusNone      TaskStatus = iota
 TaskStatusPending
 TaskStatusScheduled
 TaskStatusRunning
 TaskStatusCompleted
 TaskStatusFailed
)

// Repository 聚合仓库
type Repository struct {
 eventStore EventStore
 snapshots  map[string]*SnapshotEntry
 snapshotMu sync.RWMutex
}

type SnapshotEntry struct {
 Version  int
 Snapshot interface{}
 Time     time.Time
}

// Get 获取聚合
func (r *Repository) Get(ctx context.Context, aggregateID string) (*TaskAggregate, error) {
 aggregate := NewTaskAggregate(aggregateID)

 // 尝试从快照恢复
 if snapshotVersion, snapshot, err := r.eventStore.GetSnapshot(ctx, aggregateID); err == nil {
  if taskSnapshot, ok := snapshot.(*TaskSnapshot); ok {
   aggregate.restoreFromSnapshot(taskSnapshot)
   aggregate.version = snapshotVersion

   // 只加载快照后的事件
   events, err := r.eventStore.GetEvents(ctx, aggregateID, snapshotVersion+1)
   if err != nil {
    return nil, err
   }
   return aggregate, aggregate.LoadFromHistory(events)
  }
 }

 // 从事件历史重建
 events, err := r.eventStore.GetEvents(ctx, aggregateID, 0)
 if err != nil {
  return nil, err
 }

 if len(events) == 0 {
  return nil, fmt.Errorf("aggregate not found: %s", aggregateID)
 }

 return aggregate, aggregate.LoadFromHistory(events)
}

// Save 保存聚合
func (r *Repository) Save(ctx context.Context, aggregate *TaskAggregate) error {
 events := aggregate.UncommittedEvents()
 if len(events) == 0 {
  return nil
 }

 // 乐观并发控制
 err := r.eventStore.Append(ctx, events, aggregate.Version()-len(events))
 if err != nil {
  return fmt.Errorf("concurrency conflict: %w", err)
 }

 aggregate.MarkCommitted()

 // 检查是否需要创建快照
 if aggregate.Version()%100 == 0 {
  r.createSnapshot(ctx, aggregate)
 }

 return nil
}

// TaskSnapshot 任务快照
type TaskSnapshot struct {
 Status     TaskStatus
 WorkerID   string
 RetryCount int
 Result     []byte
}

func (t *TaskAggregate) restoreFromSnapshot(snapshot *TaskSnapshot) {
 t.status = snapshot.Status
 t.workerID = snapshot.WorkerID
 t.retryCount = snapshot.RetryCount
 t.result = snapshot.Result
}

func (r *Repository) createSnapshot(ctx context.Context, aggregate *TaskAggregate) {
 snapshot := &TaskSnapshot{
  Status:     aggregate.status,
  WorkerID:   aggregate.workerID,
  RetryCount: aggregate.retryCount,
  Result:     aggregate.result,
 }

 r.eventStore.SaveSnapshot(ctx, aggregate.ID(), aggregate.Version(), snapshot)
}

// PostgresEventStore PostgreSQL事件存储实现
type PostgresEventStore struct {
 db *sql.DB
}

// Append 原子性追加事件
func (s *PostgresEventStore) Append(ctx context.Context, events []Event, expectedVersion int) error {
 tx, err := s.db.BeginTx(ctx, nil)
 if err != nil {
  return err
 }
 defer tx.Rollback()

 // 获取当前版本
 var currentVersion int
 err = tx.QueryRowContext(ctx,
  "SELECT COALESCE(MAX(version), 0) FROM events WHERE aggregate_id = $1 FOR UPDATE",
  events[0].AggregateID()).Scan(&currentVersion)
 if err != nil {
  return err
 }

 // 乐观并发检查
 if currentVersion != expectedVersion {
  return fmt.Errorf("concurrency conflict: expected version %d, got %d",
   expectedVersion, currentVersion)
 }

 // 插入事件
 stmt, err := tx.PrepareContext(ctx, `
  INSERT INTO events (id, type, aggregate_id, aggregate_type, version,
       payload, metadata, occurred_at)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
 `)
 if err != nil {
  return err
 }
 defer stmt.Close()

 for _, event := range events {
  payload, _ := json.Marshal(event)
  meta, _ := json.Marshal(event.Metadata())

  _, err = stmt.ExecContext(ctx,
   event.EventID(),
   event.EventType(),
   event.AggregateID(),
   event.AggregateType(),
   event.EventVersion(),
   payload,
   meta,
   event.OccurredAt(),
  )
  if err != nil {
   return err
  }
 }

 return tx.Commit()
}
```

---

## 事件存储性能优化

```go
// EventBus 事件总线（CQRS读取模型更新）
type EventBus struct {
 subscribers map[string][]EventHandler
 mu          sync.RWMutex
}

type EventHandler func(ctx context.Context, event Event) error

// Projector 投影器（读取模型构建）
type Projector struct {
 readModel ReadModelStore
}

func (p *Projector) HandleTaskCreated(ctx context.Context, event *TaskCreatedEvent) error {
 return p.readModel.Insert(ctx, &TaskView{
  ID:        event.AggregateID(),
  Type:      event.TaskType,
  Status:    "pending",
  Priority:  event.Priority,
  CreatedAt: event.OccurredAt(),
 })
}

func (p *Projector) HandleTaskCompleted(ctx context.Context, event *TaskCompletedEvent) error {
 return p.readModel.Update(ctx, event.AggregateID(), map[string]interface{}{
  "status":       "completed",
  "completed_at": event.CompletedAt,
  "duration_ms":  event.Duration.Milliseconds(),
 })
}
```

---

## 形式化一致性保证

$$
\begin{aligned}
&\text{Event Store Invariants:} \\
&1. \text{Immutability: } \forall e \in \text{Events}: \neg \text{modified}(e) \\
&2. \text{Ordering: } \forall i,j: i < j \Rightarrow t_i \leq t_j \\
&3. \text{Version Monotonicity: } v_{n+1} = v_n + 1 \\
&4. \text{Idempotency: } \text{append}(E) = \text{append}(E') \Rightarrow E = E'
\end{aligned}
$$
