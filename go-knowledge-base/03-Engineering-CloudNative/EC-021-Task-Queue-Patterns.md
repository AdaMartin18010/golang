# 任务队列模式 (Task Queue Patterns)

> **分类**: 工程与云原生
> **标签**: #task-queue #patterns #messaging

---

## 优先级队列

```go
type PriorityTask struct {
    Task
    Priority int  // 数字越小优先级越高
}

type PriorityQueue struct {
    items []PriorityTask
    mu    sync.Mutex
}

func (pq *PriorityQueue) Push(task PriorityTask) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    // 按优先级插入
    inserted := false
    for i, item := range pq.items {
        if task.Priority < item.Priority {
            pq.items = append(pq.items[:i], append([]PriorityTask{task}, pq.items[i:]...)...)
            inserted = true
            break
        }
    }

    if !inserted {
        pq.items = append(pq.items, task)
    }
}

func (pq *PriorityQueue) Pop() (PriorityTask, bool) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    if len(pq.items) == 0 {
        return PriorityTask{}, false
    }

    task := pq.items[0]
    pq.items = pq.items[1:]
    return task, true
}
```

---

## 延迟队列

```go
type DelayedQueue struct {
    items []DelayedItem
    mu    sync.Mutex
    cond  *sync.Cond
}

type DelayedItem struct {
    Task
    ExecuteAt time.Time
}

func (dq *DelayedQueue) Push(item DelayedItem) {
    dq.mu.Lock()
    defer dq.mu.Unlock()

    dq.items = append(dq.items, item)
    sort.Slice(dq.items, func(i, j int) bool {
        return dq.items[i].ExecuteAt.Before(dq.items[j].ExecuteAt)
    })

    dq.cond.Signal()
}

func (dq *DelayedQueue) Poll() (DelayedItem, bool) {
    dq.mu.Lock()
    defer dq.mu.Unlock()

    for {
        if len(dq.items) == 0 {
            dq.cond.Wait()
            continue
        }

        item := dq.items[0]
        if time.Now().Before(item.ExecuteAt) {
            // 等待直到执行时间
            dq.cond.Wait()
            continue
        }

        dq.items = dq.items[1:]
        return item, true
    }
}
```

---

## 死信队列

```go
type DeadLetterQueue struct {
    maxRetries int
    mainQueue  Queue
    dlq        Queue
}

func (dlq *DeadLetterQueue) Process(task *Task) {
    if task.RetryCount >= dlq.maxRetries {
        // 移到死信队列
        dlq.dlq.Push(task)
        return
    }

    // 正常处理
    dlq.mainQueue.Push(task)
}

// 死信处理
func (dlq *DeadLetterQueue) ProcessDeadLetter(task *Task) {
    // 记录日志
    log.Printf("Dead letter: %v, error: %v", task.ID, task.LastError)

    // 告警
    alertManager.Send(Alert{
        Severity: "warning",
        Message:  fmt.Sprintf("Task %s failed after %d retries", task.ID, task.RetryCount),
    })

    // 保存到持久化存储
    dlq.store.SaveDeadLetter(task)
}
```

---

## 背压模式

```go
type BackpressureQueue struct {
    queue      chan *Task
    limit      int
    strategy   BackpressureStrategy
}

type BackpressureStrategy int

const (
    Block BackpressureStrategy = iota  // 阻塞
    Drop                               // 丢弃新任务
    DropOldest                         // 丢弃最旧任务
    Reject                             // 返回错误
)

func (bp *BackpressureQueue) Push(task *Task) error {
    select {
    case bp.queue <- task:
        return nil
    default:
        // 队列满，应用背压策略
        switch bp.strategy {
        case Block:
            bp.queue <- task  // 阻塞等待
            return nil
        case Drop:
            return ErrDropped
        case DropOldest:
            // 丢弃最旧任务
            <-bp.queue
            bp.queue <- task
            return ErrDroppedOld
        case Reject:
            return ErrQueueFull
        default:
            return ErrQueueFull
        }
    }
}
```

---

## 消息路由

```go
type MessageRouter struct {
    rules []RoutingRule
}

type RoutingRule struct {
    Condition func(*Task) bool
    Target    string
}

func (mr *MessageRouter) Route(task *Task) (string, error) {
    for _, rule := range mr.rules {
        if rule.Condition(task) {
            return rule.Target, nil
        }
    }

    return "", ErrNoRoute
}

// 使用
router := &MessageRouter{
    rules: []RoutingRule{
        {
            Condition: func(t *Task) bool {
                return t.Priority == PriorityHigh
            },
            Target: "high-priority-queue",
        },
        {
            Condition: func(t *Task) bool {
                return t.Type == "report"
            },
            Target: "report-queue",
        },
        {
            Condition: func(t *Task) bool { return true },
            Target: "default-queue",
        },
    },
}
```
