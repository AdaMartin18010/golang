# 异步任务模式 (Async Task Patterns)

> **分类**: 工程与云原生
> **标签**: #async #task #patterns #event-driven
> **参考**: CQRS, Event Sourcing, Saga Pattern

---

## 异步任务架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Async Task Processing Patterns                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Pattern 1: Fire and Forget                        │   │
│  │                                                                      │   │
│  │   Client ──► Submit Task ──► Queue ──► Worker                        │   │
│  │     │                            │                                   │   │
│  │     └──────► Immediate Ack ◄─────┘                                   │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Pattern 2: Request-Reply                          │   │
│  │                                                                      │   │
│  │   Client ──► Submit Task ──► Queue ──► Worker                        │   │
│  │     │                                              │                 │   │
│  │     │                                              ▼                 │   │
│  │     │                                         Result Queue           │   │
│  │     │                                              │                 │   │
│  │     └────────────── Wait ◄─────────────────────────┘                 │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Pattern 3: Callback/Promise                       │   │
│  │                                                                      │   │
│  │   Client ──► Submit Task ──► Queue ──► Worker                        │   │
│  │     │                                              │                 │   │
│  │     │                                              ▼                 │   │
│  │     │                                         Call Webhook           │   │
│  │     │                                              │                 │   │
│  │     └────────────── Callback ◄─────────────────────┘                 │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整异步任务实现

```go
package async

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// AsyncTask 异步任务
type AsyncTask struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Payload     json.RawMessage        `json:"payload"`
    CallbackURL string                 `json:"callback_url,omitempty"`
    ReplyTo     string                 `json:"reply_to,omitempty"`

    Status      TaskStatus             `json:"status"`
    Result      json.RawMessage        `json:"result,omitempty"`
    Error       string                 `json:"error,omitempty"`

    CreatedAt   time.Time              `json:"created_at"`
    StartedAt   *time.Time             `json:"started_at,omitempty"`
    CompletedAt *time.Time             `json:"completed_at,omitempty"`

    mu          sync.RWMutex
}

// TaskStatus 任务状态
type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "pending"
    TaskStatusProcessing TaskStatus = "processing"
    TaskStatusCompleted  TaskStatus = "completed"
    TaskStatusFailed     TaskStatus = "failed"
    TaskStatusCancelled  TaskStatus = "cancelled"
)

// TaskHandler 任务处理器
type TaskHandler func(ctx context.Context, payload json.RawMessage) (json.RawMessage, error)

// AsyncProcessor 异步处理器
type AsyncProcessor struct {
    queue       TaskQueue
    handlers    map[string]TaskHandler
    callbacks   map[string]CallbackHandler

    workerCount int
    wg          sync.WaitGroup
    ctx         context.Context
    cancel      context.CancelFunc

    mu          sync.RWMutex
}

// TaskQueue 任务队列接口
type TaskQueue interface {
    Enqueue(ctx context.Context, task *AsyncTask) error
    Dequeue(ctx context.Context) (*AsyncTask, error)
    Ack(ctx context.Context, taskID string) error
}

// CallbackHandler 回调处理器
type CallbackHandler func(ctx context.Context, task *AsyncTask) error

// NewAsyncProcessor 创建异步处理器
func NewAsyncProcessor(queue TaskQueue, workerCount int) *AsyncProcessor {
    ctx, cancel := context.WithCancel(context.Background())

    return &AsyncProcessor{
        queue:       queue,
        handlers:    make(map[string]TaskHandler),
        callbacks:   make(map[string]CallbackHandler),
        workerCount: workerCount,
        ctx:         ctx,
        cancel:      cancel,
    }
}

// RegisterHandler 注册任务处理器
func (ap *AsyncProcessor) RegisterHandler(taskType string, handler TaskHandler) {
    ap.mu.Lock()
    defer ap.mu.Unlock()
    ap.handlers[taskType] = handler
}

// RegisterCallback 注册回调处理器
func (ap *AsyncProcessor) RegisterCallback(callbackType string, handler CallbackHandler) {
    ap.mu.Lock()
    defer ap.mu.Unlock()
    ap.callbacks[callbackType] = handler
}

// Start 启动处理器
func (ap *AsyncProcessor) Start() {
    for i := 0; i < ap.workerCount; i++ {
        ap.wg.Add(1)
        go ap.worker(i)
    }
}

// Stop 停止处理器
func (ap *AsyncProcessor) Stop() {
    ap.cancel()
    ap.wg.Wait()
}

func (ap *AsyncProcessor) worker(id int) {
    defer ap.wg.Done()

    for {
        select {
        case <-ap.ctx.Done():
            return
        default:
        }

        task, err := ap.queue.Dequeue(ap.ctx)
        if err != nil {
            continue
        }
        if task == nil {
            time.Sleep(100 * time.Millisecond)
            continue
        }

        ap.processTask(ap.ctx, task)
    }
}

func (ap *AsyncProcessor) processTask(ctx context.Context, task *AsyncTask) {
    // 更新状态
    task.mu.Lock()
    task.Status = TaskStatusProcessing
    now := time.Now()
    task.StartedAt = &now
    task.mu.Unlock()

    // 获取处理器
    ap.mu.RLock()
    handler, ok := ap.handlers[task.Type]
    ap.mu.RUnlock()

    if !ok {
        task.mu.Lock()
        task.Status = TaskStatusFailed
        task.Error = fmt.Sprintf("no handler for task type: %s", task.Type)
        task.mu.Unlock()
        ap.handleCallback(ctx, task)
        return
    }

    // 执行任务
    result, err := handler(ctx, task.Payload)

    task.mu.Lock()
    now = time.Now()
    task.CompletedAt = &now

    if err != nil {
        task.Status = TaskStatusFailed
        task.Error = err.Error()
    } else {
        task.Status = TaskStatusCompleted
        task.Result = result
    }
    task.mu.Unlock()

    // 确认消息
    ap.queue.Ack(ctx, task.ID)

    // 处理回调
    ap.handleCallback(ctx, task)
}

func (ap *AsyncProcessor) handleCallback(ctx context.Context, task *AsyncTask) {
    if task.CallbackURL != "" {
        // HTTP 回调
        ap.executeHTTPCallback(ctx, task)
    }

    if task.ReplyTo != "" {
        // 队列回复
        ap.executeQueueCallback(ctx, task)
    }
}

func (ap *AsyncProcessor) executeHTTPCallback(ctx context.Context, task *AsyncTask) {
    // HTTP 回调实现
}

func (ap *AsyncProcessor) executeQueueCallback(ctx context.Context, task *AsyncTask) {
    // 队列回调实现
}

// Submit 提交任务 (Fire and Forget)
func (ap *AsyncProcessor) Submit(ctx context.Context, taskType string, payload interface{}) (string, error) {
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    task := &AsyncTask{
        ID:        generateID(),
        Type:      taskType,
        Payload:   payloadBytes,
        Status:    TaskStatusPending,
        CreatedAt: time.Now(),
    }

    if err := ap.queue.Enqueue(ctx, task); err != nil {
        return "", err
    }

    return task.ID, nil
}

// SubmitWithReply 提交任务并等待回复 (Request-Reply)
func (ap *AsyncProcessor) SubmitWithReply(ctx context.Context, taskType string, payload interface{}, timeout time.Duration) (*AsyncTask, error) {
    replyQueue := make(chan *AsyncTask, 1)

    callbackType := fmt.Sprintf("reply-%s", generateID())

    // 注册临时回调
    ap.RegisterCallback(callbackType, func(ctx context.Context, task *AsyncTask) error {
        replyQueue <- task
        return nil
    })
    defer delete(ap.callbacks, callbackType)

    // 提交任务
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    task := &AsyncTask{
        ID:        generateID(),
        Type:      taskType,
        Payload:   payloadBytes,
        ReplyTo:   callbackType,
        Status:    TaskStatusPending,
        CreatedAt: time.Now(),
    }

    if err := ap.queue.Enqueue(ctx, task); err != nil {
        return nil, err
    }

    // 等待回复
    select {
    case result := <-replyQueue:
        return result, nil
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout waiting for reply")
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// SubmitWithCallback 提交任务并设置回调 (Callback/Promise)
func (ap *AsyncProcessor) SubmitWithCallback(ctx context.Context, taskType string, payload interface{}, callbackURL string) (string, error) {
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    task := &AsyncTask{
        ID:          generateID(),
        Type:        taskType,
        Payload:     payloadBytes,
        CallbackURL: callbackURL,
        Status:      TaskStatusPending,
        CreatedAt:   time.Now(),
    }

    if err := ap.queue.Enqueue(ctx, task); err != nil {
        return "", err
    }

    return task.ID, nil
}

func generateID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

---

## Promise/Future 模式

```go
package async

import (
    "context"
    "encoding/json"
    "sync"
    "time"
)

// Promise 异步承诺
type Promise struct {
    ID      string
    task    *AsyncTask
    done    chan struct{}
    result  json.RawMessage
    err     error

    mu      sync.RWMutex
}

// Future 异步结果
type Future struct {
    promise *Promise
}

// NewPromise 创建 Promise
func NewPromise(id string) *Promise {
    return &Promise{
        ID:   id,
        done: make(chan struct{}),
    }
}

// Complete 完成 Promise
func (p *Promise) Complete(result json.RawMessage, err error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.result = result
    p.err = err
    close(p.done)
}

// GetFuture 获取 Future
func (p *Promise) GetFuture() *Future {
    return &Future{promise: p}
}

// Get 获取结果（阻塞）
func (f *Future) Get(ctx context.Context) (json.RawMessage, error) {
    select {
    case <-f.promise.done:
        return f.promise.result, f.promise.err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// GetWithTimeout 获取结果（带超时）
func (f *Future) GetWithTimeout(timeout time.Duration) (json.RawMessage, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    return f.Get(ctx)
}

// IsDone 是否完成
func (f *Future) IsDone() bool {
    select {
    case <-f.promise.done:
        return true
    default:
        return false
    }
}

// PromiseManager Promise 管理器
type PromiseManager struct {
    promises map[string]*Promise
    mu       sync.RWMutex
}

// NewPromiseManager 创建 Promise 管理器
func NewPromiseManager() *PromiseManager {
    return &PromiseManager{
        promises: make(map[string]*Promise),
    }
}

// Create 创建 Promise
func (pm *PromiseManager) Create(id string) *Promise {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    promise := NewPromise(id)
    pm.promises[id] = promise

    return promise
}

// Get 获取 Promise
func (pm *PromiseManager) Get(id string) (*Promise, bool) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    promise, ok := pm.promises[id]
    return promise, ok
}

// Complete 完成 Promise
func (pm *PromiseManager) Complete(id string, result json.RawMessage, err error) bool {
    pm.mu.Lock()
    promise, ok := pm.promises[id]
    if ok {
        delete(pm.promises, id)
    }
    pm.mu.Unlock()

    if !ok {
        return false
    }

    promise.Complete(result, err)
    return true
}

// Cleanup 清理过期的 Promise
func (pm *PromiseManager) Cleanup(maxAge time.Duration) {
    // 实现清理逻辑
}
```

---

## 事件驱动模式

```go
package async

import (
    "context"
    "sync"
)

// Event 事件
type Event struct {
    Type    string
    Payload interface{}
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) error

// EventBus 事件总线
type EventBus struct {
    subscribers map[string][]EventHandler
    mu          sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]EventHandler),
    }
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

// Publish 发布事件
func (eb *EventBus) Publish(ctx context.Context, event Event) {
    eb.mu.RLock()
    handlers := make([]EventHandler, len(eb.subscribers[event.Type]))
    copy(handlers, eb.subscribers[event.Type])
    eb.mu.RUnlock()

    for _, handler := range handlers {
        go handler(ctx, event)
    }
}

// Saga 编排器
type Saga struct {
    steps      []SagaStep
    compensations []CompensationStep

    mu         sync.Mutex
    completed  []int
}

// SagaStep Saga 步骤
type SagaStep struct {
    Name     string
    Execute  func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

// CompensationStep 补偿步骤
type CompensationStep struct {
    StepIndex int
    Execute   func(ctx context.Context) error
}

// NewSaga 创建 Saga
func NewSaga() *Saga {
    return &Saga{
        steps:      make([]SagaStep, 0),
        compensations: make([]CompensationStep, 0),
        completed:  make([]int, 0),
    }
}

// AddStep 添加步骤
func (s *Saga) AddStep(step SagaStep) {
    s.steps = append(s.steps, step)
}

// Execute 执行 Saga
func (s *Saga) Execute(ctx context.Context) error {
    for i, step := range s.steps {
        if err := step.Execute(ctx); err != nil {
            // 执行补偿
            s.compensate(ctx, i)
            return err
        }

        s.mu.Lock()
        s.completed = append(s.completed, i)
        s.mu.Unlock()
    }

    return nil
}

func (s *Saga) compensate(ctx context.Context, failedIndex int) {
    s.mu.Lock()
    completed := make([]int, len(s.completed))
    copy(completed, s.completed)
    s.mu.Unlock()

    // 逆序执行补偿
    for i := len(completed) - 1; i >= 0; i-- {
        stepIndex := completed[i]
        step := s.steps[stepIndex]
        if step.Compensate != nil {
            step.Compensate(ctx)
        }
    }
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
    "time"

    "async"
)

func main() {
    // Fire and Forget
    queue := NewMemoryQueue()
    processor := async.NewAsyncProcessor(queue, 10)

    processor.RegisterHandler("send-email", func(ctx context.Context, payload json.RawMessage) (json.RawMessage, error) {
        fmt.Println("Sending email...")
        time.Sleep(1 * time.Second)
        return json.Marshal(map[string]string{"status": "sent"})
    })

    processor.Start()
    defer processor.Stop()

    // Fire and forget
    taskID, _ := processor.Submit(context.Background(), "send-email", map[string]string{
        "to":      "user@example.com",
        "subject": "Hello",
    })
    fmt.Printf("Task submitted: %s\n", taskID)

    // Request-Reply
    processor.RegisterHandler("calculate", func(ctx context.Context, payload json.RawMessage) (json.RawMessage, error) {
        var req map[string]int
        json.Unmarshal(payload, &req)
        result := req["a"] + req["b"]
        return json.Marshal(map[string]int{"result": result})
    })

    result, _ := processor.SubmitWithReply(context.Background(), "calculate", map[string]int{
        "a": 10,
        "b": 20,
    }, 5*time.Second)

    fmt.Printf("Result: %s\n", result.Result)

    // Promise/Future
    pm := async.NewPromiseManager()
    promise := pm.Create("task-123")
    future := promise.GetFuture()

    // 异步完成
    go func() {
        time.Sleep(2 * time.Second)
        pm.Complete("task-123", json.RawMessage(`{"data": "ok"}`), nil)
    }()

    data, err := future.GetWithTimeout(5 * time.Second)
    fmt.Printf("Future result: %s, err: %v\n", data, err)

    // Saga
    saga := async.NewSaga()

    saga.AddStep(async.SagaStep{
        Name: "deduct-inventory",
        Execute: func(ctx context.Context) error {
            fmt.Println("Deducting inventory...")
            return nil
        },
        Compensate: func(ctx context.Context) error {
            fmt.Println("Restoring inventory...")
            return nil
        },
    })

    saga.AddStep(async.SagaStep{
        Name: "process-payment",
        Execute: func(ctx context.Context) error {
            fmt.Println("Processing payment...")
            return nil
        },
        Compensate: func(ctx context.Context) error {
            fmt.Println("Refunding payment...")
            return nil
        },
    })

    if err := saga.Execute(context.Background()); err != nil {
        fmt.Printf("Saga failed: %v\n", err)
    }
}
```
