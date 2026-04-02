# 任务事件溯源 (Task Event Sourcing)

> **分类**: 工程与云原生
> **标签**: #event-sourcing #cqrs #audit

---

## 事件定义

```go
// 任务领域事件
type TaskEvent interface {
    EventType() string
    EventID() string
    AggregateID() string  // TaskID
    OccurredAt() time.Time
}

type TaskCreatedEvent struct {
    ID          string
    TaskID      string
    Name        string
    Type        string
    Payload     []byte
    ScheduledAt time.Time
    Timestamp   time.Time
}

func (e TaskCreatedEvent) EventType() string   { return "TaskCreated" }
func (e TaskCreatedEvent) EventID() string     { return e.ID }
func (e TaskCreatedEvent) AggregateID() string { return e.TaskID }
func (e TaskCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

type TaskStartedEvent struct {
    ID        string
    TaskID    string
    WorkerID  string
    Timestamp time.Time
}

type TaskCompletedEvent struct {
    ID        string
    TaskID    string
    Result    []byte
    Duration  time.Duration
    Timestamp time.Time
}

type TaskFailedEvent struct {
    ID        string
    TaskID    string
    Error     string
    Retryable bool
    Timestamp time.Time
}
```

---

## 事件存储

```go
type EventStore interface {
    Append(ctx context.Context, events []TaskEvent) error
    GetEvents(ctx context.Context, taskID string, fromVersion int) ([]TaskEvent, error)
    GetAllEvents(ctx context.Context, afterPosition int64) ([]TaskEvent, error)
}

// PostgreSQL 实现
type PostgresEventStore struct {
    db *sql.DB
}

func (pes *PostgresEventStore) Append(ctx context.Context, events []TaskEvent) error {
    tx, err := pes.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO task_events (event_id, task_id, event_type, payload, version, position)
        VALUES ($1, $2, $3, $4, $5, nextval('event_position_seq'))
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, event := range events {
        payload, _ := json.Marshal(event)
        _, err := stmt.ExecContext(ctx,
            event.EventID(),
            event.AggregateID(),
            event.EventType(),
            payload,
            0, // version
        )
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

---

## 事件溯源聚合

```go
type TaskAggregate struct {
    TaskID   string
    Version  int
    Events   []TaskEvent
    State    TaskState
}

type TaskState struct {
    Status   string
    WorkerID string
    Result   []byte
    Error    string
    Attempts int
}

func (ta *TaskAggregate) Apply(event TaskEvent) {
    switch e := event.(type) {
    case TaskCreatedEvent:
        ta.State.Status = "pending"

    case TaskStartedEvent:
        ta.State.Status = "running"
        ta.State.WorkerID = e.WorkerID

    case TaskCompletedEvent:
        ta.State.Status = "completed"
        ta.State.Result = e.Result

    case TaskFailedEvent:
        ta.State.Status = "failed"
        ta.State.Error = e.Error
        ta.State.Attempts++
    }

    ta.Version++
    ta.Events = append(ta.Events, event)
}

func (ta *TaskAggregate) Replay(events []TaskEvent) {
    for _, event := range events {
        ta.Apply(event)
    }
}
```

---

## 事件投影 (CQRS)

```go
// 读模型投影
type TaskProjection struct {
    store TaskReadStore
}

func (tp *TaskProjection) HandleEvent(ctx context.Context, event TaskEvent) error {
    switch e := event.(type) {
    case TaskCreatedEvent:
        return tp.store.Create(ctx, TaskView{
            ID:        e.TaskID,
            Name:      e.Name,
            Type:      e.Type,
            Status:    "pending",
            CreatedAt: e.Timestamp,
        })

    case TaskStartedEvent:
        return tp.store.Update(ctx, e.TaskID, map[string]interface{}{
            "status":    "running",
            "worker_id": e.WorkerID,
            "started_at": e.Timestamp,
        })

    case TaskCompletedEvent:
        return tp.store.Update(ctx, e.TaskID, map[string]interface{}{
            "status":     "completed",
            "result":     e.Result,
            "duration":   e.Duration,
            "completed_at": e.Timestamp,
        })
    }

    return nil
}

// 物化视图
type TaskView struct {
    ID          string
    Name        string
    Type        string
    Status      string
    WorkerID    string
    Result      []byte
    Duration    time.Duration
    CreatedAt   time.Time
    StartedAt   *time.Time
    CompletedAt *time.Time
}
```

---

## 审计日志

```go
type TaskAuditLog struct {
    store EventStore
}

func (tal *TaskAuditLog) GetAuditTrail(ctx context.Context, taskID string) ([]AuditEntry, error) {
    events, err := tal.store.GetEvents(ctx, taskID, 0)
    if err != nil {
        return nil, err
    }

    var entries []AuditEntry
    for _, event := range events {
        entries = append(entries, AuditEntry{
            Time:      event.OccurredAt(),
            Action:    event.EventType(),
            Actor:     getActor(event),
            Details:   event,
        })
    }

    return entries, nil
}

func (tal *TaskAuditLog) WhoChanged(ctx context.Context, taskID string, field string) ([]ChangeRecord, error) {
    events, _ := tal.store.GetEvents(ctx, taskID, 0)

    var changes []ChangeRecord
    var prevState TaskState

    for _, event := range events {
        // 重建状态
        newState := prevState
        applyEvent(&newState, event)

        // 检测字段变化
        if diff := compareStates(prevState, newState); diff != nil {
            for _, change := range diff.Changes {
                if change.Field == field {
                    changes = append(changes, ChangeRecord{
                        Time:     event.OccurredAt(),
                        OldValue: change.Old,
                        NewValue: change.New,
                        Actor:    getActor(event),
                    })
                }
            }
        }

        prevState = newState
    }

    return changes, nil
}
```
