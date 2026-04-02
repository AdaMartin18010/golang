# 任务队列完整实现 (Task Queue Implementation)

> **分类**: 工程与云原生
> **标签**: #task-queue #priority-queue #delayed-queue #distributed
> **参考**: Redis Streams, RabbitMQ, SQS

---

## 任务队列架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Distributed Task Queue System                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Producer Layer                                  │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Submit Task │  │   Delay     │  │  Schedule   │  │   Batch     │ │   │
│  │  │ (Enqueue)   │  │   Task      │  │   Task      │  │   Submit    │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                       Queue Router                                     │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │  Priority   │  │   Round     │  │   Hash      │  │   Load      │ │   │
│  │  │   Queue     │  │   Robin     │  │   Routing   │  │   Balance   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                     Queue Implementations                              │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Redis     │  │   Kafka     │  │   RabbitMQ  │  │  In-Memory  │ │   │
│  │  │  Streams    │  │  Partition  │  │   Queue     │  │   Channel   │ │   │
│  │  │  (ZSET)     │  │             │  │             │  │             │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                      Consumer Layer                                    │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Worker    │  │  Parallel   │  │   Retry     │  │   Dead      │ │   │
│  │  │   Pool      │  │   Consume   │  │   Handler   │  │   Letter    │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心接口定义

```go
package taskqueue

import (
    "context"
    "time"
)

// Task 任务定义
type Task struct {
    ID          string
    Type        string
    Payload     []byte
    Priority    int       // 优先级，数字越小优先级越高
    Delay       time.Duration // 延迟执行时间
    MaxRetries  int       // 最大重试次数
    RetryCount  int       // 当前重试次数

    // 元数据
    CreatedAt   time.Time
    ScheduledAt *time.Time // 计划执行时间

    // 追踪
    TraceID     string
    SpanID      string
}

// TaskResult 任务执行结果
type TaskResult struct {
    TaskID  string
    Success bool
    Error   error
    Retry   bool // 是否需要重试
}

// TaskHandler 任务处理器
type TaskHandler func(ctx context.Context, task *Task) error

// Queue 队列接口
type Queue interface {
    // 入队
    Enqueue(ctx context.Context, task *Task) error

    // 出队（阻塞）
    Dequeue(ctx context.Context) (*Task, error)

    // 确认完成
    Ack(ctx context.Context, taskID string) error

    // 拒绝/重试
    Nack(ctx context.Context, taskID string, requeue bool) error

    // 队列深度
    Size(ctx context.Context) (int64, error)

    // 关闭
    Close() error
}

// PriorityQueue 优先级队列接口
type PriorityQueue interface {
    Queue

    // 入队带优先级
    EnqueueWithPriority(ctx context.Context, task *Task, priority int) error

    // 获取优先级范围内的任务
    DequeueByPriority(ctx context.Context, maxPriority int) (*Task, error)
}

// DelayedQueue 延迟队列接口
type DelayedQueue interface {
    Queue

    // 延迟入队
    EnqueueDelayed(ctx context.Context, task *Task, delay time.Duration) error

    // 处理到期的延迟任务
    ProcessDelayed(ctx context.Context) error
}

// Consumer 消费者接口
type Consumer interface {
    // 开始消费
    Start(ctx context.Context) error

    // 停止消费
    Stop() error

    // 注册处理器
    RegisterHandler(taskType string, handler TaskHandler)
}

// Producer 生产者接口
type Producer interface {
    // 发送任务
    Send(ctx context.Context, task *Task) error

    // 延迟发送
    SendDelayed(ctx context.Context, task *Task, delay time.Duration) error

    // 批量发送
    SendBatch(ctx context.Context, tasks []*Task) error
}
```

---

## Redis 实现（Streams + ZSET）

```go
package taskqueue

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisQueue Redis 队列实现
type RedisQueue struct {
    client     *redis.Client
    stream     string // 主队列 (Streams)
    group      string // 消费者组
    consumer   string // 消费者名称

    delayedKey string // 延迟队列 (ZSET)
    dlqKey     string // 死信队列

    maxLen     int64
}

func NewRedisQueue(client *redis.Client, stream, group, consumer string) *RedisQueue {
    return &RedisQueue{
        client:     client,
        stream:     stream,
        group:      group,
        consumer:   consumer,
        delayedKey: stream + ":delayed",
        dlqKey:     stream + ":dlq",
        maxLen:     10000,
    }
}

// CreateGroup 创建消费者组
func (q *RedisQueue) CreateGroup(ctx context.Context) error {
    return q.client.XGroupCreateMkStream(ctx, q.stream, q.group, "0").Err()
}

// Enqueue 入队
func (q *RedisQueue) Enqueue(ctx context.Context, task *Task) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }

    args := &redis.XAddArgs{
        Stream: q.stream,
        MaxLen: q.maxLen,
        Approx: true,
        Values: map[string]interface{}{
            "task": string(data),
        },
    }

    _, err = q.client.XAdd(ctx, args).Result()
    return err
}

// EnqueueWithPriority 优先级入队
func (q *RedisQueue) EnqueueWithPriority(ctx context.Context, task *Task, priority int) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }

    // 使用优先级分数：高分在前
    score := float64(^int(priority)) // 反转优先级

    return q.client.ZAdd(ctx, q.stream+":priority", redis.Z{
        Score:  score,
        Member: string(data),
    }).Err()
}

// EnqueueDelayed 延迟入队
func (q *RedisQueue) EnqueueDelayed(ctx context.Context, task *Task, delay time.Duration) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }

    executeAt := time.Now().Add(delay).UnixMilli()

    return q.client.ZAdd(ctx, q.delayedKey, redis.Z{
        Score:  float64(executeAt),
        Member: string(data),
    }).Err()
}

// Dequeue 出队
func (q *RedisQueue) Dequeue(ctx context.Context) (*Task, error) {
    // 1. 先处理延迟队列
    if err := q.ProcessDelayed(ctx); err != nil {
        // 非致命错误，继续
    }

    // 2. 从主队列读取
    streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
        Group:    q.group,
        Consumer: q.consumer,
        Streams:  []string{q.stream, ">"},
        Count:    1,
        Block:    5 * time.Second,
    }).Result()

    if err != nil {
        if err == redis.Nil {
            return nil, nil // 无消息
        }
        return nil, err
    }

    if len(streams) == 0 || len(streams[0].Messages) == 0 {
        return nil, nil
    }

    msg := streams[0].Messages[0]
    return q.parseTask(msg.Values["task"].(string))
}

// DequeueByPriority 按优先级出队
func (q *RedisQueue) DequeueByPriority(ctx context.Context, maxPriority int) (*Task, error) {
    // 获取最高优先级的任务
    result, err := q.client.ZPopMax(ctx, q.stream+":priority", 1).Result()
    if err != nil {
        return nil, err
    }

    if len(result) == 0 {
        return nil, nil
    }

    return q.parseTask(result[0].Member.(string))
}

// parseTask 解析任务
func (q *RedisQueue) parseTask(data string) (*Task, error) {
    var task Task
    if err := json.Unmarshal([]byte(data), &task); err != nil {
        return nil, err
    }
    return &task, nil
}

// Ack 确认
func (q *RedisQueue) Ack(ctx context.Context, taskID string) error {
    // 从 pending 列表确认
    return q.client.XAck(ctx, q.stream, q.group, taskID).Err()
}

// Nack 拒绝
func (q *RedisQueue) Nack(ctx context.Context, taskID string, requeue bool) error {
    if !requeue {
        // 直接确认（丢弃）
        return q.Ack(ctx, taskID)
    }

    // 重新入队：XCLAIM 修改空闲时间，让其被其他消费者处理
    // 或 XADD 重新添加
    _, err := q.client.XClaim(ctx, &redis.XClaimArgs{
        Stream:   q.stream,
        Group:    q.group,
        Consumer: q.consumer,
        MinIdle:  0,
        Messages: []string{taskID},
    }).Result()

    return err
}

// ProcessDelayed 处理延迟任务
func (q *RedisQueue) ProcessDelayed(ctx context.Context) error {
    now := time.Now().UnixMilli()

    // 获取到期的任务
    tasks, err := q.client.ZRangeByScoreWithScores(ctx, q.delayedKey, &redis.ZRangeBy{
        Min: "0",
        Max: fmt.Sprintf("%d", now),
    }).Result()

    if err != nil {
        return err
    }

    for _, t := range tasks {
        task, err := q.parseTask(t.Member.(string))
        if err != nil {
            continue
        }

        // 移动到主队列
        if err := q.Enqueue(ctx, task); err != nil {
            continue
        }

        // 从延迟队列移除
        q.client.ZRem(ctx, q.delayedKey, t.Member)
    }

    return nil
}

// ClaimStuckTasks 认领卡住的任务（其他消费者崩溃）
func (q *RedisQueue) ClaimStuckTasks(ctx context.Context, minIdle time.Duration, count int64) ([]*Task, error) {
    // XPENDING 获取挂起任务
    pending, err := q.client.XPendingExt(ctx, &redis.XPendingExtArgs{
        Stream: q.stream,
        Group:  q.group,
        Start:  "-",
        End:    "+",
        Count:  count,
    }).Result()

    if err != nil {
        return nil, err
    }

    var tasks []*Task
    for _, p := range pending {
        if p.Idle < minIdle {
            continue
        }

        // XCLAIM 认领
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
            task, err := q.parseTask(msg.Values["task"].(string))
            if err != nil {
                continue
            }
            tasks = append(tasks, task)
        }
    }

    return tasks, nil
}

// MoveToDLQ 移动到死信队列
func (q *RedisQueue) MoveToDLQ(ctx context.Context, task *Task, reason string) error {
    taskData, _ := json.Marshal(task)

    dlqEntry := map[string]interface{}{
        "task":   string(taskData),
        "reason": reason,
        "time":   time.Now().Unix(),
    }

    entryData, _ := json.Marshal(dlqEntry)

    return q.client.LPush(ctx, q.dlqKey, entryData).Err()
}

// Size 队列深度
func (q *RedisQueue) Size(ctx context.Context) (int64, error) {
    info, err := q.client.XInfoStream(ctx, q.stream).Result()
    if err != nil {
        return 0, err
    }
    return info.Length, nil
}

func (q *RedisQueue) Close() error {
    return q.client.Close()
}
```

---

## 消费者实现

```go
package taskqueue

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
    WorkerCount    int           // 工作线程数
    MaxRetries     int           // 最大重试次数
    RetryDelay     time.Duration // 重试延迟
    BatchSize      int           // 批处理大小
    PollInterval   time.Duration // 轮询间隔
    HandleTimeout  time.Duration // 处理超时
}

// DefaultConsumerConfig 默认配置
var DefaultConsumerConfig = ConsumerConfig{
    WorkerCount:   10,
    MaxRetries:    3,
    RetryDelay:    5 * time.Second,
    BatchSize:     1,
    PollInterval:  1 * time.Second,
    HandleTimeout: 30 * time.Second,
}

// TaskConsumer 任务消费者
type TaskConsumer struct {
    queue    Queue
    config   ConsumerConfig
    handlers map[string]TaskHandler

    wg       sync.WaitGroup
    ctx      context.Context
    cancel   context.CancelFunc

    mu       sync.RWMutex
    running  bool
}

func NewTaskConsumer(queue Queue, config ConsumerConfig) *TaskConsumer {
    ctx, cancel := context.WithCancel(context.Background())

    return &TaskConsumer{
        queue:    queue,
        config:   config,
        handlers: make(map[string]TaskHandler),
        ctx:      ctx,
        cancel:   cancel,
    }
}

// RegisterHandler 注册处理器
func (c *TaskConsumer) RegisterHandler(taskType string, handler TaskHandler) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.handlers[taskType] = handler
}

// Start 启动消费
func (c *TaskConsumer) Start(ctx context.Context) error {
    c.mu.Lock()
    if c.running {
        c.mu.Unlock()
        return fmt.Errorf("already running")
    }
    c.running = true
    c.mu.Unlock()

    // 启动工作线程
    for i := 0; i < c.config.WorkerCount; i++ {
        c.wg.Add(1)
        go c.worker(i)
    }

    return nil
}

// Stop 停止消费
func (c *TaskConsumer) Stop() error {
    c.mu.Lock()
    if !c.running {
        c.mu.Unlock()
        return nil
    }
    c.running = false
    c.mu.Unlock()

    c.cancel()
    c.wg.Wait()

    return nil
}

func (c *TaskConsumer) worker(id int) {
    defer c.wg.Done()

    for {
        select {
        case <-c.ctx.Done():
            return
        default:
        }

        // 获取任务
        task, err := c.queue.Dequeue(c.ctx)
        if err != nil {
            time.Sleep(c.config.PollInterval)
            continue
        }

        if task == nil {
            time.Sleep(c.config.PollInterval)
            continue
        }

        // 处理任务
        c.processTask(c.ctx, task)
    }
}

func (c *TaskConsumer) processTask(ctx context.Context, task *Task) {
    // 获取处理器
    c.mu.RLock()
    handler, ok := c.handlers[task.Type]
    c.mu.RUnlock()

    if !ok {
        // 无处理器，确认丢弃
        c.queue.Ack(ctx, task.ID)
        return
    }

    // 创建带超时的上下文
    taskCtx, cancel := context.WithTimeout(ctx, c.config.HandleTimeout)
    defer cancel()

    // 执行任务
    err := handler(taskCtx, task)

    if err == nil {
        // 成功，确认
        c.queue.Ack(ctx, task.ID)
        return
    }

    // 失败处理
    if task.RetryCount < c.config.MaxRetries {
        // 重试
        task.RetryCount++
        time.Sleep(c.config.RetryDelay)
        c.queue.Nack(ctx, task.ID, true)
    } else {
        // 超过重试次数，移动到 DLQ
        c.queue.Ack(ctx, task.ID) // 从原队列移除
        if dlq, ok := c.queue.(*RedisQueue); ok {
            dlq.MoveToDLQ(ctx, task, err.Error())
        }
    }
}
```

---

## 生产者实现

```go
package taskqueue

import (
    "context"
    "sync"

    "github.com/google/uuid"
)

// TaskProducer 任务生产者
type TaskProducer struct {
    queue  Queue
    config ProducerConfig
    mu     sync.Mutex
}

type ProducerConfig struct {
    DefaultMaxRetries int
    DefaultPriority   int
}

func NewTaskProducer(queue Queue, config ProducerConfig) *TaskProducer {
    return &TaskProducer{
        queue:  queue,
        config: config,
    }
}

// Send 发送任务
func (p *TaskProducer) Send(ctx context.Context, taskType string, payload []byte) error {
    task := &Task{
        ID:         uuid.New().String(),
        Type:       taskType,
        Payload:    payload,
        Priority:   p.config.DefaultPriority,
        MaxRetries: p.config.DefaultMaxRetries,
    }

    return p.queue.Enqueue(ctx, task)
}

// SendDelayed 延迟发送
func (p *TaskProducer) SendDelayed(ctx context.Context, taskType string, payload []byte, delay time.Duration) error {
    task := &Task{
        ID:         uuid.New().String(),
        Type:       taskType,
        Payload:    payload,
        Priority:   p.config.DefaultPriority,
        MaxRetries: p.config.DefaultMaxRetries,
    }

    if dq, ok := p.queue.(DelayedQueue); ok {
        return dq.EnqueueDelayed(ctx, task, delay)
    }

    // 不支持延迟队列，直接入队
    return p.queue.Enqueue(ctx, task)
}

// SendWithPriority 带优先级发送
func (p *TaskProducer) SendWithPriority(ctx context.Context, taskType string, payload []byte, priority int) error {
    task := &Task{
        ID:         uuid.New().String(),
        Type:       taskType,
        Payload:    payload,
        Priority:   priority,
        MaxRetries: p.config.DefaultMaxRetries,
    }

    if pq, ok := p.queue.(PriorityQueue); ok {
        return pq.EnqueueWithPriority(ctx, task, priority)
    }

    return p.queue.Enqueue(ctx, task)
}

// SendBatch 批量发送
func (p *TaskProducer) SendBatch(ctx context.Context, tasks []*Task) error {
    // 简单实现：顺序发送
    // 生产环境应使用管道/事务批量写入
    for _, task := range tasks {
        if err := p.queue.Enqueue(ctx, task); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 内存队列实现（测试用）

```go
package taskqueue

import (
    "container/heap"
    "context"
    "sync"
    "time"
)

// MemoryQueue 内存队列实现
type MemoryQueue struct {
    tasks    chan *Task
    delayed  *DelayedTaskHeap
    mu       sync.RWMutex
    closed   bool
}

func NewMemoryQueue(size int) *MemoryQueue {
    mq := &MemoryQueue{
        tasks:   make(chan *Task, size),
        delayed: &DelayedTaskHeap{},
    }
    heap.Init(mq.delayed)
    return mq
}

func (q *MemoryQueue) Enqueue(ctx context.Context, task *Task) error {
    q.mu.RLock()
    if q.closed {
        q.mu.RUnlock()
        return fmt.Errorf("queue closed")
    }
    q.mu.RUnlock()

    select {
    case q.tasks <- task:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (q *MemoryQueue) EnqueueDelayed(ctx context.Context, task *Task, delay time.Duration) error {
    q.mu.Lock()
    defer q.mu.Unlock()

    executeAt := time.Now().Add(delay)
    task.ScheduledAt = &executeAt

    heap.Push(q.delayed, &DelayedTask{
        Task:      task,
        ExecuteAt: executeAt,
    })

    return nil
}

func (q *MemoryQueue) Dequeue(ctx context.Context) (*Task, error) {
    // 检查延迟队列
    q.processDelayed()

    select {
    case task := <-q.tasks:
        return task, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func (q *MemoryQueue) processDelayed() {
    q.mu.Lock()
    defer q.mu.Unlock()

    now := time.Now()
    for q.delayed.Len() > 0 {
        dt := (*q.delayed)[0]
        if dt.ExecuteAt.After(now) {
            break
        }

        heap.Pop(q.delayed)
        select {
        case q.tasks <- dt.Task:
        default:
            // 队列满，放回
            heap.Push(q.delayed, dt)
            return
        }
    }
}

func (q *MemoryQueue) Ack(ctx context.Context, taskID string) error {
    return nil // 内存队列无需确认
}

func (q *MemoryQueue) Nack(ctx context.Context, taskID string, requeue bool) error {
    return nil
}

func (q *MemoryQueue) Size(ctx context.Context) (int64, error) {
    return int64(len(q.tasks)), nil
}

func (q *MemoryQueue) Close() error {
    q.mu.Lock()
    defer q.mu.Unlock()

    if !q.closed {
        q.closed = true
        close(q.tasks)
    }
    return nil
}

// DelayedTask 延迟任务
type DelayedTask struct {
    *Task
    ExecuteAt time.Time
}

// DelayedTaskHeap 延迟任务堆
type DelayedTaskHeap []*DelayedTask

func (h DelayedTaskHeap) Len() int { return len(h) }
func (h DelayedTaskHeap) Less(i, j int) bool { return h[i].ExecuteAt.Before(h[j].ExecuteAt) }
func (h DelayedTaskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *DelayedTaskHeap) Push(x interface{}) {
    *h = append(*h, x.(*DelayedTask))
}

func (h *DelayedTaskHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}
```
