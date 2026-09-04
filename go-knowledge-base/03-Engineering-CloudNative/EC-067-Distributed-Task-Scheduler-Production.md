# 分布式任务调度器生产实践 (Distributed Task Scheduler Production)

> **分类**: 工程与云原生
> **标签**: #distributed-scheduler #production #scalability #reliability
> **参考**: Uber Cadence, Temporal, Kubernetes Scheduler

> **维度**: Engineering & CloudNative
> **级别**: S (29 KB)
> **Go 版本**: 1.27+
---

## 生产级架构设计

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Production Distributed Scheduler                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        API Gateway Layer                            │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │    │
│  │  │  Rate Limit │  │   Auth      │  │  Validate   │  │   Route     │ │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐  │
│  │                      Scheduler Cluster (HA)                           │  │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                │  │
│  │  │ Scheduler 1 │◄──►│ Scheduler 2 │◄──►│ Scheduler N │                │  │
│  │  │  (Leader)   │    │  (Follower) │    │  (Follower) │               │   │
│  │  └──────┬──────┘    └─────────────┘    └─────────────┘               │   │
│  │         │                                                            │   │
│  │         │  Leader Election (etcd/ZooKeeper)                          │   │
│  │         ▼                                                            │   │
│  │  ┌─────────────────────────────────────────────────────────────┐     │   │
│  │  │                    Task State Machine                       │     │   │
│  │  │  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐      │     │   │
│  │  │  │ Pending │──►│Scheduled│──►│ Running │──►│Completed│      │     │   │
│  │  │  └─────────┘   └─────────┘   └────┬────┘   └─────────┘      │     │   │
│  │  │         │           │           │           │               │     │   │
│  │  │         ▼           ▼           ▼           ▼               │     │   │
│  │  │     ┌──────┐    ┌──────┐    ┌──────┐    ┌──────┐            │     │   │
│  │  │     │Cancel│    │Retry │    │Timeout│   │Fail  │            │     │   │
│  │  │     └──────┘    └──────┘    └──────┘    └──────┘            │     │   │
│  │  └─────────────────────────────────────────────────────────────┘     │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐  │
│  │                      Worker Pool Layer                                │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐              │   │
│  │  │ Worker 1 │  │ Worker 2 │  │ Worker 3 │  │ Worker N │              │   │
│  │  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │              │   │
│  │  │ │Task Q│ │  │ │Task Q│ │  │ │Task Q│ │  │ │Task Q│ │              │   │
│  │  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │              │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Storage Layer                                     │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Task DB   │  │   Queue     │  │    Log    │  │   Metrics   │ │   │
│  │  │ (PostgreSQL)│  │   (Redis)   │  │  (Kafka)  │  │ (Prometheus)│ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 任务状态机实现

```go
package scheduler

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// TaskState 任务状态
type TaskState string

const (
    TaskStatePending    TaskState = "pending"
    TaskStateScheduled  TaskState = "scheduled"
    TaskStateRunning    TaskState = "running"
    TaskStateCompleted  TaskState = "completed"
    TaskStateFailed     TaskState = "failed"
    TaskStateCancelled  TaskState = "cancelled"
    TaskStateRetrying   TaskState = "retrying"
    TaskStateTimeout    TaskState = "timeout"
)

// Task 任务定义
type Task struct {
    ID          string
    Type        string
    Payload     []byte
    State       TaskState
    Priority    int

    // 调度信息
    ScheduledAt *time.Time
    StartedAt   *time.Time
    CompletedAt *time.Time

    // 重试信息
    RetryCount  int
    MaxRetries  int
    RetryDelay  time.Duration

    // 超时控制
    Timeout     time.Duration

    // 工作节点
    WorkerID    string

    // 错误信息
    LastError   string

    // 上下文传播
    TraceID     string
    SpanID      string
    Baggage     map[string]string

    mu          sync.RWMutex
}

// StateMachine 任务状态机
type StateMachine struct {
    transitions map[TaskState]map[TaskState]bool
    handlers    map[TaskState]StateHandler
}

type StateHandler func(ctx context.Context, task *Task) error

func NewStateMachine() *StateMachine {
    sm := &StateMachine{
        transitions: make(map[TaskState]map[TaskState]bool),
        handlers:    make(map[TaskState]StateHandler),
    }

    // 定义合法状态转换
    sm.addTransition(TaskStatePending, TaskStateScheduled)
    sm.addTransition(TaskStatePending, TaskStateCancelled)

    sm.addTransition(TaskStateScheduled, TaskStateRunning)
    sm.addTransition(TaskStateScheduled, TaskStateCancelled)
    sm.addTransition(TaskStateScheduled, TaskStateTimeout)

    sm.addTransition(TaskStateRunning, TaskStateCompleted)
    sm.addTransition(TaskStateRunning, TaskStateFailed)
    sm.addTransition(TaskStateRunning, TaskStateTimeout)
    sm.addTransition(TaskStateRunning, TaskStateCancelled)

    sm.addTransition(TaskStateFailed, TaskStateRetrying)
    sm.addTransition(TaskStateTimeout, TaskStateRetrying)

    sm.addTransition(TaskStateRetrying, TaskStateScheduled)
    sm.addTransition(TaskStateRetrying, TaskStateFailed) // 超过重试次数

    return sm
}

func (sm *StateMachine) addTransition(from, to TaskState) {
    if sm.transitions[from] == nil {
        sm.transitions[from] = make(map[TaskState]bool)
    }
    sm.transitions[from][to] = true
}

func (sm *StateMachine) CanTransition(from, to TaskState) bool {
    if trans, ok := sm.transitions[from]; ok {
        return trans[to]
    }
    return false
}

func (sm *StateMachine) Transition(ctx context.Context, task *Task, to TaskState) error {
    task.mu.Lock()
    defer task.mu.Unlock()

    from := task.State

    if !sm.CanTransition(from, to) {
        return fmt.Errorf("invalid transition from %s to %s", from, to)
    }

    // 执行状态处理器
    if handler, ok := sm.handlers[to]; ok {
        if err := handler(ctx, task); err != nil {
            return fmt.Errorf("state handler failed: %w", err)
        }
    }

    task.State = to

    switch to {
    case TaskStateScheduled:
        now := time.Now()
        task.ScheduledAt = &now
    case TaskStateRunning:
        now := time.Now()
        task.StartedAt = &now
    case TaskStateCompleted, TaskStateFailed, TaskStateCancelled, TaskStateTimeout:
        now := time.Now()
        task.CompletedAt = &now
    }

    return nil
}

func (sm *StateMachine) RegisterHandler(state TaskState, handler StateHandler) {
    sm.handlers[state] = handler
}
```

---

## 分布式锁实现（基于 Redis）

```go
package scheduler

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisLock 基于 Redis 的分布式锁
type RedisLock struct {
    client *redis.Client
    key    string
    value  string // 唯一标识，防止误释放
    ttl    time.Duration
}

// NewRedisLock 创建分布式锁
func NewRedisLock(client *redis.Client, key string, ttl time.Duration) *RedisLock {
    // 生成唯一标识
    b := make([]byte, 16)
    rand.Read(b)
    value := hex.EncodeToString(b)

    return &RedisLock{
        client: client,
        key:    key,
        value:  value,
        ttl:    ttl,
    }
}

// Lock 获取锁（阻塞）
func (l *RedisLock) Lock(ctx context.Context) error {
    for {
        ok, err := l.TryLock(ctx)
        if err != nil {
            return err
        }
        if ok {
            return nil
        }

        // 等待重试
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(100 * time.Millisecond):
        }
    }
}

// TryLock 尝试获取锁（非阻塞）
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
    // SET key value NX EX ttl
    ok, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
    if err != nil {
        return false, err
    }

    if ok {
        // 启动续约 goroutine
        go l.renew(ctx)
    }

    return ok, nil
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
    // Lua 脚本：只有 value 匹配时才删除
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `

    result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
    if err != nil {
        return err
    }

    if result.(int64) == 0 {
        return fmt.Errorf("lock not held or expired")
    }

    return nil
}

// renew 自动续约
func (l *RedisLock) renew(ctx context.Context) {
    ticker := time.NewTicker(l.ttl / 3)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // 续约
            script := `
                if redis.call("get", KEYS[1]) == ARGV[1] then
                    return redis.call("pexpire", KEYS[1], ARGV[2])
                else
                    return 0
                end
            `
            ttlMs := l.ttl.Milliseconds()
            result, err := l.client.Eval(ctx, script, []string{l.key}, l.value, ttlMs).Result()
            if err != nil || result.(int64) == 0 {
                return
            }
        }
    }
}

// RedLock Redis RedLock 算法（多节点）
type RedLock struct {
    clients []*redis.Client
    quorum  int
}

func NewRedLock(clients []*redis.Client) *RedLock {
    return &RedLock{
        clients: clients,
        quorum:  len(clients)/2 + 1,
    }
}

func (rl *RedLock) Lock(ctx context.Context, key string, ttl time.Duration) (*RedisLock, error) {
    value := generateUniqueID()

    successes := 0
    startTime := time.Now()

    for _, client := range rl.clients {
        lock := &RedisLock{
            client: client,
            key:    key,
            value:  value,
            ttl:    ttl,
        }

        ok, _ := lock.TryLock(ctx)
        if ok {
            successes++
        }
    }

    elapsed := time.Since(startTime)
    validity := ttl - elapsed - 2*time.Millisecond // 时钟漂移补偿

    if successes >= rl.quorum && validity > 0 {
        // 获取成功
        return &RedisLock{
            clients: rl.clients,
            key:     key,
            value:   value,
            ttl:     ttl,
        }, nil
    }

    // 获取失败，释放所有锁
    for _, client := range rl.clients {
        lock := &RedisLock{client: client, key: key, value: value}
        lock.Unlock(ctx)
    }

    return nil, fmt.Errorf("failed to acquire redlock")
}

func generateUniqueID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}
```

---

## 任务队列实现（基于 Redis Streams）

```go
package scheduler

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisStreamQueue 基于 Redis Streams 的任务队列
type RedisStreamQueue struct {
    client    *redis.Client
    stream    string
    group     string
    consumer  string
    maxLen    int64
}

func NewRedisStreamQueue(client *redis.Client, stream, group, consumer string) *RedisStreamQueue {
    return &RedisStreamQueue{
        client:   client,
        stream:   stream,
        group:    group,
        consumer: consumer,
        maxLen:   10000,
    }
}

// CreateGroup 创建消费者组
func (q *RedisStreamQueue) CreateGroup(ctx context.Context) error {
    // MKSTREAM 选项：如果 stream 不存在则创建
    err := q.client.XGroupCreateMkStream(ctx, q.stream, q.group, "$").Err()
    if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
        return err
    }
    return nil
}

// Enqueue 添加任务到队列
func (q *RedisStreamQueue) Enqueue(ctx context.Context, task *Task) (string, error) {
    data, err := json.Marshal(task)
    if err != nil {
        return "", err
    }

    // XADD stream MAXLEN ~ maxlen * task data
    id, err := q.client.XAdd(ctx, &redis.XAddArgs{
        Stream: q.stream,
        MaxLen: q.maxLen,
        Approx: true,
        Values: map[string]interface{}{
            "task": string(data),
        },
    }).Result()

    return id, err
}

// Dequeue 从队列获取任务（阻塞）
func (q *RedisStreamQueue) Dequeue(ctx context.Context, timeout time.Duration) (*Task, string, error) {
    // XREADGROUP GROUP group consumer BLOCK timeout STREAMS stream >
    streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
        Group:    q.group,
        Consumer: q.consumer,
        Streams:  []string{q.stream, ">"},
        Block:    timeout,
        Count:    1,
    }).Result()

    if err != nil {
        if err == redis.Nil {
            return nil, "", nil
        }
        return nil, "", err
    }

    if len(streams) == 0 || len(streams[0].Messages) == 0 {
        return nil, "", nil
    }

    msg := streams[0].Messages[0]
    taskData := msg.Values["task"].(string)

    var task Task
    if err := json.Unmarshal([]byte(taskData), &task); err != nil {
        return nil, "", err
    }

    return &task, msg.ID, nil
}

// Ack 确认任务完成
func (q *RedisStreamQueue) Ack(ctx context.Context, msgID string) error {
    return q.client.XAck(ctx, q.stream, q.group, msgID).Err()
}

// Claim 认领挂起的任务（其他消费者崩溃）
func (q *RedisStreamQueue) Claim(ctx context.Context, minIdle time.Duration) ([]*Task, error) {
    // XPENDING stream group - 获取挂起任务
    pending, err := q.client.XPendingExt(ctx, &redis.XPendingExtArgs{
        Stream: q.stream,
        Group:  q.group,
        Start:  "-",
        End:    "+",
        Count:  10,
    }).Result()

    if err != nil {
        return nil, err
    }

    var tasks []*Task
    for _, p := range pending {
        if p.Idle < minIdle {
            continue
        }

        // XCLAIM stream group consumer min-idle-time ID
        claimed, err := q.client.XClaim(ctx, &redis.XClaimArgs{
            Stream:   q.stream,
            Group:    q.group,
            Consumer: q.consumer,
            MinIdle:  minIdle,
            Messages: []string{p.ID},
        }).Result()

        if err != nil {
            continue
        }

        for _, msg := range claimed {
            taskData := msg.Values["task"].(string)
            var task Task
            if err := json.Unmarshal([]byte(taskData), &task); err != nil {
                continue
            }
            tasks = append(tasks, &task)
        }
    }

    return tasks, nil
}

// ScheduleDelayed 延迟任务（使用 ZSET）
func (q *RedisStreamQueue) ScheduleDelayed(ctx context.Context, task *Task, delay time.Duration) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }

    executeAt := time.Now().Add(delay).UnixMilli()

    // ZADD delayed_queue score task_data
    return q.client.ZAdd(ctx, q.stream+":delayed", redis.Z{
        Score:  float64(executeAt),
        Member: string(data),
    }).Err()
}

// ProcessDelayed 处理到期的延迟任务
func (q *RedisStreamQueue) ProcessDelayed(ctx context.Context) error {
    now := time.Now().UnixMilli()

    // ZRANGEBYSCORE delayed_queue 0 now
    tasks, err := q.client.ZRangeByScoreWithScores(ctx, q.stream+":delayed", &redis.ZRangeBy{
        Min: "0",
        Max: fmt.Sprintf("%d", now),
    }).Result()

    if err != nil {
        return err
    }

    for _, t := range tasks {
        var task Task
        if err := json.Unmarshal([]byte(t.Member.(string)), &task); err != nil {
            continue
        }

        // 添加到主队列
        if _, err := q.Enqueue(ctx, &task); err != nil {
            continue
        }

        // 从延迟队列移除
        q.client.ZRem(ctx, q.stream+":delayed", t.Member)
    }

    return nil
}
```

---

## 任务调度器实现

```go
package scheduler

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Scheduler 任务调度器
type Scheduler struct {
    stateMachine *StateMachine
    queue        *RedisStreamQueue
    lock         *RedisLock
    store        TaskStore

    workers      int
    wg           sync.WaitGroup
    ctx          context.Context
    cancel       context.CancelFunc

    // 处理器注册表
    handlers     map[string]TaskHandler
    handlersMu   sync.RWMutex
}

type TaskHandler func(ctx context.Context, task *Task) error
type TaskStore interface {
    Save(ctx context.Context, task *Task) error
    Get(ctx context.Context, id string) (*Task, error)
    Update(ctx context.Context, task *Task) error
}

func NewScheduler(queue *RedisStreamQueue, store TaskStore, workers int) *Scheduler {
    ctx, cancel := context.WithCancel(context.Background())

    return &Scheduler{
        stateMachine: NewStateMachine(),
        queue:        queue,
        store:        store,
        workers:      workers,
        ctx:          ctx,
        cancel:       cancel,
        handlers:     make(map[string]TaskHandler),
    }
}

// RegisterHandler 注册任务处理器
func (s *Scheduler) RegisterHandler(taskType string, handler TaskHandler) {
    s.handlersMu.Lock()
    defer s.handlersMu.Unlock()
    s.handlers[taskType] = handler
}

// Start 启动调度器
func (s *Scheduler) Start() error {
    // 创建消费者组
    if err := s.queue.CreateGroup(s.ctx); err != nil {
        return err
    }

    // 启动工作线程
    for i := 0; i < s.workers; i++ {
        s.wg.Add(1)
        go s.worker(i)
    }

    // 启动延迟任务处理器
    go s.delayedProcessor()

    return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
    s.cancel()
    s.wg.Wait()
}

func (s *Scheduler) worker(id int) {
    defer s.wg.Done()

    for {
        select {
        case <-s.ctx.Done():
            return
        default:
        }

        // 获取任务
        task, msgID, err := s.queue.Dequeue(s.ctx, 5*time.Second)
        if err != nil {
            continue
        }
        if task == nil {
            continue
        }

        // 处理任务
        if err := s.processTask(s.ctx, task); err != nil {
            // 处理失败，可能重试
            s.handleFailure(s.ctx, task, err)
        }

        // 确认完成
        s.queue.Ack(s.ctx, msgID)
    }
}

func (s *Scheduler) processTask(ctx context.Context, task *Task) error {
    // 获取处理器
    s.handlersMu.RLock()
    handler, ok := s.handlers[task.Type]
    s.handlersMu.RUnlock()

    if !ok {
        return fmt.Errorf("no handler for task type: %s", task.Type)
    }

    // 状态转换：Scheduled -> Running
    if err := s.stateMachine.Transition(ctx, task, TaskStateRunning); err != nil {
        return err
    }

    // 保存状态
    if err := s.store.Update(ctx, task); err != nil {
        return err
    }

    // 创建带超时的上下文
    taskCtx, cancel := context.WithTimeout(ctx, task.Timeout)
    defer cancel()

    // 执行处理器
    err := handler(taskCtx, task)

    if err != nil {
        task.LastError = err.Error()
        s.stateMachine.Transition(ctx, task, TaskStateFailed)
    } else {
        s.stateMachine.Transition(ctx, task, TaskStateCompleted)
    }

    // 保存最终状态
    s.store.Update(ctx, task)

    return err
}

func (s *Scheduler) handleFailure(ctx context.Context, task *Task, err error) {
    if task.RetryCount < task.MaxRetries {
        // 可以重试
        task.RetryCount++
        s.stateMachine.Transition(ctx, task, TaskStateRetrying)

        // 延迟重试
        s.queue.ScheduleDelayed(ctx, task, task.RetryDelay)
    } else {
        // 超过重试次数
        s.stateMachine.Transition(ctx, task, TaskStateFailed)
    }

    s.store.Update(ctx, task)
}

func (s *Scheduler) delayedProcessor() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-s.ctx.Done():
            return
        case <-ticker.C:
            s.queue.ProcessDelayed(s.ctx)
        }
    }
}

// SubmitTask 提交任务
func (s *Scheduler) SubmitTask(ctx context.Context, task *Task) error {
    // 初始化任务
    task.State = TaskStatePending
    task.RetryCount = 0

    if task.MaxRetries == 0 {
        task.MaxRetries = 3
    }
    if task.Timeout == 0 {
        task.Timeout = 5 * time.Minute
    }

    // 保存到存储
    if err := s.store.Save(ctx, task); err != nil {
        return err
    }

    // 状态转换
    s.stateMachine.Transition(ctx, task, TaskStateScheduled)

    // 添加到队列
    if task.ScheduledAt != nil && task.ScheduledAt.After(time.Now()) {
        // 延迟任务
        delay := time.Until(*task.ScheduledAt)
        return s.queue.ScheduleDelayed(ctx, task, delay)
    }

    _, err := s.queue.Enqueue(ctx, task)
    return err
}
```

---

## 生产级监控

```go
package scheduler

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // 任务计数
    tasksSubmitted = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "scheduler_tasks_submitted_total",
        Help: "Total number of tasks submitted",
    }, []string{"type"})

    tasksCompleted = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "scheduler_tasks_completed_total",
        Help: "Total number of tasks completed",
    }, []string{"type"})

    tasksFailed = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "scheduler_tasks_failed_total",
        Help: "Total number of tasks failed",
    }, []string{"type"})

    // 任务执行时间
    taskDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "scheduler_task_duration_seconds",
        Help:    "Task execution duration",
        Buckets: prometheus.DefBuckets,
    }, []string{"type"})

    // 队列深度
    queueDepth = promauto.NewGaugeVec(prometheus.GaugeOpts{
        Name: "scheduler_queue_depth",
        Help: "Current queue depth",
    }, []string{"queue"})

    // 工作线程状态
    workerActive = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "scheduler_workers_active",
        Help: "Number of active workers",
    })
)

// MetricsCollector 监控收集器
type MetricsCollector struct {
    scheduler *Scheduler
}

func (m *MetricsCollector) Collect() {
    // 定期收集队列深度等指标
    // ...
}
```

---

## 最佳实践

```go
// 1. 任务幂等性
type IdempotentHandler struct {
    processed map[string]bool
    mu        sync.RWMutex
}

func (h *IdempotentHandler) Handle(ctx context.Context, task *Task) error {
    h.mu.RLock()
    if h.processed[task.ID] {
        h.mu.RUnlock()
        return nil // 已处理，直接返回
    }
    h.mu.RUnlock()

    // 处理任务
    if err := h.doProcess(ctx, task); err != nil {
        return err
    }

    h.mu.Lock()
    h.processed[task.ID] = true
    h.mu.Unlock()

    return nil
}

// 2. 任务去重（基于业务键）
type DeduplicationMiddleware struct {
    store DeduplicationStore
    next  TaskHandler
}

func (m *DeduplicationMiddleware) Handle(ctx context.Context, task *Task) error {
    dedupKey := task.Type + ":" + task.BusinessKey

    if exists, _ := m.store.Exists(ctx, dedupKey); exists {
        return nil // 重复任务，忽略
    }

    // 设置去重标记（带 TTL）
    m.store.Set(ctx, dedupKey, "", task.Timeout)

    return m.next(ctx, task)
}

// 3. 优雅关闭
func GracefulShutdown(scheduler *Scheduler, timeout time.Duration) {
    // 停止接受新任务
    scheduler.StopAccepting()

    // 等待正在处理的任务完成
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    scheduler.Wait(ctx)
}
```
