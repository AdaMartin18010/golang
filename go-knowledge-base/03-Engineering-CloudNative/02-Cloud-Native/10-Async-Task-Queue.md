# 异步任务队列 (Async Task Queue)

> **分类**: 工程与云原生  
> **标签**: #async #queue #task #background

---

## 基于 Channel 的任务队列

```go
type TaskQueue struct {
    tasks   chan Task
    workers int
    wg      sync.WaitGroup
    ctx     context.Context
    cancel  context.CancelFunc
}

type Task struct {
    ID       string
    Execute  func(ctx context.Context) error
    Callback func(result interface{}, err error)
}

func NewTaskQueue(workers, buffer int) *TaskQueue {
    ctx, cancel := context.WithCancel(context.Background())
    return &TaskQueue{
        tasks:   make(chan Task, buffer),
        workers: workers,
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (q *TaskQueue) Start() {
    for i := 0; i < q.workers; i++ {
        q.wg.Add(1)
        go q.worker(i)
    }
}

func (q *TaskQueue) worker(id int) {
    defer q.wg.Done()
    
    for task := range q.tasks {
        select {
        case <-q.ctx.Done():
            return
        default:
            q.executeTask(task)
        }
    }
}

func (q *TaskQueue) executeTask(task Task) {
    ctx, cancel := context.WithTimeout(q.ctx, 5*time.Minute)
    defer cancel()
    
    err := task.Execute(ctx)
    
    if task.Callback != nil {
        task.Callback(nil, err)
    }
}

func (q *TaskQueue) Enqueue(task Task) bool {
    select {
    case q.tasks <- task:
        return true
    default:
        return false
    }
}

func (q *TaskQueue) Stop() {
    q.cancel()
    close(q.tasks)
    q.wg.Wait()
}
```

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
    cond  *sync.Cond
}

func NewPriorityQueue() *PriorityQueue {
    pq := &PriorityQueue{}
    pq.cond = sync.NewCond(&pq.mu)
    return pq
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
    
    pq.cond.Signal()
}

func (pq *PriorityQueue) Pop() (PriorityTask, bool) {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    
    for len(pq.items) == 0 {
        pq.cond.Wait()
    }
    
    if len(pq.items) == 0 {
        return PriorityTask{}, false
    }
    
    task := pq.items[0]
    pq.items = pq.items[1:]
    return task, true
}
```

---

## 持久化任务队列 (基于 SQLite)

```go
type PersistentQueue struct {
    db *sql.DB
}

func (q *PersistentQueue) Init() error {
    _, err := q.db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id TEXT PRIMARY KEY,
            type TEXT NOT NULL,
            payload BLOB,
            status TEXT DEFAULT 'pending',
            retry_count INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            scheduled_at DATETIME,
            processed_at DATETIME
        )
    `)
    return err
}

func (q *PersistentQueue) Enqueue(task Task, scheduleAt time.Time) error {
    _, err := q.db.Exec(
        "INSERT INTO tasks (id, type, payload, scheduled_at) VALUES (?, ?, ?, ?)",
        task.ID, task.Type, task.Payload, scheduleAt,
    )
    return err
}

func (q *PersistentQueue) Dequeue() (*Task, error) {
    tx, err := q.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    
    // 获取并锁定任务
    var task Task
    err = tx.QueryRow(`
        SELECT id, type, payload FROM tasks 
        WHERE status = 'pending' AND scheduled_at <= ?
        ORDER BY scheduled_at ASC
        LIMIT 1
    `, time.Now()).Scan(&task.ID, &task.Type, &task.Payload)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    // 标记为处理中
    _, err = tx.Exec("UPDATE tasks SET status = 'processing' WHERE id = ?", task.ID)
    if err != nil {
        return nil, err
    }
    
    return &task, tx.Commit()
}

func (q *PersistentQueue) Complete(taskID string, success bool) error {
    status := "completed"
    if !success {
        status = "failed"
    }
    
    _, err := q.db.Exec(
        "UPDATE tasks SET status = ?, processed_at = ? WHERE id = ?",
        status, time.Now(), taskID,
    )
    return err
}
```

---

## 任务编排 (DAG)

```go
type TaskNode struct {
    ID       string
    Execute  func(ctx context.Context) error
    DependsOn []string
}

type DAGExecutor struct {
    tasks    map[string]*TaskNode
    completed map[string]bool
    mu        sync.Mutex
}

func (e *DAGExecutor) Execute(ctx context.Context) error {
    for len(e.completed) < len(e.tasks) {
        // 找到可执行的任务（依赖已完成）
        ready := e.findReadyTasks()
        if len(ready) == 0 {
            if len(e.completed) < len(e.tasks) {
                return fmt.Errorf("circular dependency detected")
            }
            break
        }
        
        // 并行执行就绪任务
        var wg sync.WaitGroup
        errChan := make(chan error, len(ready))
        
        for _, task := range ready {
            wg.Add(1)
            go func(t *TaskNode) {
                defer wg.Done()
                if err := t.Execute(ctx); err != nil {
                    errChan <- fmt.Errorf("task %s failed: %w", t.ID, err)
                }
            }(task)
        }
        
        wg.Wait()
        close(errChan)
        
        // 收集错误
        for err := range errChan {
            return err
        }
        
        // 标记完成
        for _, task := range ready {
            e.completed[task.ID] = true
        }
    }
    
    return nil
}

func (e *DAGExecutor) findReadyTasks() []*TaskNode {
    var ready []*TaskNode
    
    for id, task := range e.tasks {
        if e.completed[id] {
            continue
        }
        
        // 检查依赖是否完成
        allDepsCompleted := true
        for _, dep := range task.DependsOn {
            if !e.completed[dep] {
                allDepsCompleted = false
                break
            }
        }
        
        if allDepsCompleted {
            ready = append(ready, task)
        }
    }
    
    return ready
}
```
