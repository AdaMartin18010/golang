# 延迟任务调度 (Delayed Task Scheduling)

> **分类**: 工程与云原生
> **标签**: #delayed-tasks #scheduling #timing-wheel
> **参考**: Kafka Delayed Queue, Timing Wheel Algorithm

---

## 延迟任务架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Delayed Task Scheduling Architecture                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Approach 1: Priority Queue (ZSET)                 │   │
│  │                                                                      │   │
│  │   Delay Queue (Redis ZSET)                                           │   │
│  │   ┌─────────────────────────────────────────────────────────────┐    │   │
│  │   │  Score (execute_at) │ Task ID │ Payload                     │    │   │
│  │   ├─────────────────────┼─────────┼─────────────────────────────┤    │   │
│  │   │  1640995200000      │ task-1  │ {"type":"email","to":"a"}    │   │   │
│  │   │  1640995260000      │ task-2  │ {"type":"sms","to":"b"}      │   │   │
│  │   │  1640995320000      │ task-3  │ {"type":"push","to":"c"}     │   │   │
│  │   └─────────────────────────────────────────────────────────────┘    │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Approach 2: Timing Wheel                          │   │
│  │                                                                      │   │
│  │                     1 tick = 1ms                                     │   │
│  │                     Wheel size = 512 slots                           │   │
│  │                     1 round = 512ms                                  │   │
│  │                                                                      │   │
│  │   Current: 247                                                       │   │
│  │         │                                                            │   │
│  │         ▼                                                            │   │
│  │   ┌─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┐                  │   │
│  │   │ 244 │ 245 │ 246 │ 247 │ 248 │ 249 │ 250 │ 251 │ ...              │   │
│  │   │     │     │     │ [●] │     │     │     │     │                  │   │
│  │   └─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┘                  │   │
│  │                        ▲                                             │   │
│  │                        │                                             │   │
│  │                   Task expires at 247                                │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Approach 3: Hierarchical Timing Wheel             │   │
│  │                                                                      │   │
│  │   Wheel 1: 20ms tick, 20 slots  (covers 400ms)                       │   │
│  │   Wheel 2: 400ms tick, 20 slots (covers 8s)                          │   │
│  │   Wheel 3: 8s tick,   20 slots  (covers 160s)                        │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整延迟任务实现

```go
package delayed

import (
    "container/heap"
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// DelayedTask 延迟任务
type DelayedTask struct {
    ID        string
    Payload   interface{}
    ExecuteAt time.Time

    // 内部使用
    index     int
    cancelled int32
}

// IsCancelled 是否已取消
func (dt *DelayedTask) IsCancelled() bool {
    return atomic.LoadInt32(&dt.cancelled) == 1
}

// Cancel 取消任务
func (dt *DelayedTask) Cancel() {
    atomic.StoreInt32(&dt.cancelled, 1)
}

// DelayedQueue 延迟队列接口
type DelayedQueue interface {
    Offer(task *DelayedTask) error
    Poll() (*DelayedTask, bool)
    Peek() (*DelayedTask, bool)
    Size() int
}

// PriorityDelayedQueue 基于优先队列的延迟队列
type PriorityDelayedQueue struct {
    tasks []*DelayedTask
    mu    sync.RWMutex
}

// NewPriorityDelayedQueue 创建优先队列
func NewPriorityDelayedQueue() *PriorityDelayedQueue {
    pq := &PriorityDelayedQueue{
        tasks: make([]*DelayedTask, 0),
    }
    heap.Init(pq)
    return pq
}

func (pq *PriorityDelayedQueue) Len() int {
    pq.mu.RLock()
    defer pq.mu.RUnlock()
    return len(pq.tasks)
}

func (pq *PriorityDelayedQueue) Less(i, j int) bool {
    return pq.tasks[i].ExecuteAt.Before(pq.tasks[j].ExecuteAt)
}

func (pq *PriorityDelayedQueue) Swap(i, j int) {
    pq.tasks[i], pq.tasks[j] = pq.tasks[j], pq.tasks[i]
    pq.tasks[i].index = i
    pq.tasks[j].index = j
}

func (pq *PriorityDelayedQueue) Push(x interface{}) {
    n := len(pq.tasks)
    task := x.(*DelayedTask)
    task.index = n
    pq.tasks = append(pq.tasks, task)
}

func (pq *PriorityDelayedQueue) Pop() interface{} {
    old := pq.tasks
    n := len(old)
    task := old[n-1]
    old[n-1] = nil
    task.index = -1
    pq.tasks = old[0 : n-1]
    return task
}

func (pq *PriorityDelayedQueue) Offer(task *DelayedTask) error {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    heap.Push(pq, task)
    return nil
}

func (pq *PriorityDelayedQueue) Poll() (*DelayedTask, bool) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    if len(pq.tasks) == 0 {
        return nil, false
    }

    task := heap.Pop(pq).(*DelayedTask)
    return task, true
}

func (pq *PriorityDelayedQueue) Peek() (*DelayedTask, bool) {
    pq.mu.RLock()
    defer pq.mu.RUnlock()

    if len(pq.tasks) == 0 {
        return nil, false
    }

    return pq.tasks[0], true
}

func (pq *PriorityDelayedQueue) Size() int {
    return pq.Len()
}

// Scheduler 延迟任务调度器
type Scheduler struct {
    queue       DelayedQueue
    executor    Executor

    // 控制
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup

    // 状态
    running     int32
}

// Executor 执行器接口
type Executor interface {
    Execute(ctx context.Context, task *DelayedTask) error
}

// NewScheduler 创建调度器
func NewScheduler(queue DelayedQueue, executor Executor) *Scheduler {
    ctx, cancel := context.WithCancel(context.Background())

    return &Scheduler{
        queue:    queue,
        executor: executor,
        ctx:      ctx,
        cancel:   cancel,
    }
}

// Start 启动调度器
func (s *Scheduler) Start() {
    if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
        return
    }

    s.wg.Add(1)
    go s.scheduleLoop()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
    if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
        return
    }

    s.cancel()
    s.wg.Wait()
}

// Schedule 调度任务
func (s *Scheduler) Schedule(task *DelayedTask) error {
    return s.queue.Offer(task)
}

// ScheduleFunc 调度函数
func (s *Scheduler) ScheduleFunc(id string, payload interface{}, delay time.Duration, fn func(ctx context.Context, payload interface{}) error) error {
    task := &DelayedTask{
        ID:        id,
        Payload:   payload,
        ExecuteAt: time.Now().Add(delay),
    }

    // 包装执行器
    executor := &funcExecutor{fn: fn}
    s.executor = executor

    return s.Schedule(task)
}

func (s *Scheduler) scheduleLoop() {
    defer s.wg.Done()

    for atomic.LoadInt32(&s.running) == 1 {
        task, ok := s.queue.Peek()
        if !ok {
            time.Sleep(100 * time.Millisecond)
            continue
        }

        // 计算等待时间
        waitTime := time.Until(task.ExecuteAt)
        if waitTime > 0 {
            select {
            case <-time.After(waitTime):
            case <-s.ctx.Done():
                return
            }
        }

        // 取出任务
        task, ok = s.queue.Poll()
        if !ok {
            continue
        }

        // 检查是否已取消
        if task.IsCancelled() {
            continue
        }

        // 执行任务
        s.wg.Add(1)
        go func(t *DelayedTask) {
            defer s.wg.Done()

            ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
            defer cancel()

            _ = s.executor.Execute(ctx, t)
        }(task)
    }
}

type funcExecutor struct {
    fn func(ctx context.Context, payload interface{}) error
}

func (fe *funcExecutor) Execute(ctx context.Context, task *DelayedTask) error {
    return fe.fn(ctx, task.Payload)
}

// TimingWheel 时间轮
type TimingWheel struct {
    tickMs     int64         // 每tick的毫秒数
    wheelSize  int           // 轮大小
    interval   int64         // 一轮的毫秒数

    currentTime int64        // 当前时间（ms）
    buckets    []*TimingWheelBucket

    queue      DelayedQueue  // 上级时间轮

    ctx        context.Context
    cancel     context.CancelFunc

    overflowWheel *TimingWheel // 溢出时间轮
    mu         sync.Mutex
}

// TimingWheelBucket 时间轮槽
type TimingWheelBucket struct {
    expiration int64
    tasks      map[string]*DelayedTask
    mu         sync.Mutex
}

// NewTimingWheel 创建时间轮
func NewTimingWheel(tickMs int64, wheelSize int) *TimingWheel {
    ctx, cancel := context.WithCancel(context.Background())

    tw := &TimingWheel{
        tickMs:      tickMs,
        wheelSize:   wheelSize,
        interval:    tickMs * int64(wheelSize),
        currentTime: time.Now().UnixMilli(),
        buckets:     make([]*TimingWheelBucket, wheelSize),
        ctx:         ctx,
        cancel:      cancel,
    }

    for i := 0; i < wheelSize; i++ {
        tw.buckets[i] = &TimingWheelBucket{
            tasks: make(map[string]*DelayedTask),
        }
    }

    return tw
}

// Add 添加任务
func (tw *TimingWheel) Add(task *DelayedTask) error {
    tw.mu.Lock()
    defer tw.mu.Unlock()

    if task.IsCancelled() {
        return nil
    }

    expiration := task.ExecuteAt.UnixMilli()

    // 计算延迟时间
    delay := expiration - tw.currentTime

    if delay < tw.tickMs {
        // 立即执行
        return nil
    } else if delay < tw.interval {
        // 在当前轮
        virtualID := expiration / tw.tickMs
        bucket := tw.buckets[virtualID%int64(tw.wheelSize)]
        bucket.mu.Lock()
        bucket.tasks[task.ID] = task
        bucket.mu.Unlock()
    } else {
        // 在溢出轮
        if tw.overflowWheel == nil {
            tw.addOverflowWheel()
        }
        return tw.overflowWheel.Add(task)
    }

    return nil
}

func (tw *TimingWheel) addOverflowWheel() {
    tw.overflowWheel = NewTimingWheel(tw.interval, tw.wheelSize)
}

// AdvanceClock 推进时钟
func (tw *TimingWheel) AdvanceClock(timeMs int64) {
    tw.mu.Lock()
    defer tw.mu.Unlock()

    if timeMs >= tw.currentTime+tw.tickMs {
        tw.currentTime = timeMs - (timeMs % tw.tickMs)

        if tw.overflowWheel != nil {
            tw.overflowWheel.AdvanceClock(timeMs)
        }
    }
}

// GetExpiredTasks 获取过期任务
func (tw *TimingWheel) GetExpiredTasks() []*DelayedTask {
    tw.mu.Lock()
    defer tw.mu.Unlock()

    var expired []*DelayedTask
    now := time.Now().UnixMilli()

    for _, bucket := range tw.buckets {
        bucket.mu.Lock()
        for id, task := range bucket.tasks {
            if task.ExecuteAt.UnixMilli() <= now {
                expired = append(expired, task)
                delete(bucket.tasks, id)
            }
        }
        bucket.mu.Unlock()
    }

    return expired
}

// RedisDelayedQueue Redis延迟队列
type RedisDelayedQueue struct {
    client RedisClient
    key    string
}

type RedisClient interface {
    ZAdd(ctx context.Context, key string, members ...interface{}) error
    ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error)
    ZRem(ctx context.Context, key string, members ...interface{}) error
}

// Offer 添加任务
func (rdq *RedisDelayedQueue) Offer(task *DelayedTask) error {
    ctx := context.Background()
    score := float64(task.ExecuteAt.UnixMilli())

    data, _ := json.Marshal(task)
    return rdq.client.ZAdd(ctx, rdq.key, score, string(data))
}

// PollDueTasks 拉取到期任务
func (rdq *RedisDelayedQueue) PollDueTasks(ctx context.Context, batchSize int) ([]*DelayedTask, error) {
    max := fmt.Sprintf("%d", time.Now().UnixMilli())

    results, err := rdq.client.ZRangeByScore(ctx, rdq.key, "0", max)
    if err != nil {
        return nil, err
    }

    var tasks []*DelayedTask
    for _, data := range results {
        var task DelayedTask
        if err := json.Unmarshal([]byte(data), &task); err != nil {
            continue
        }
        tasks = append(tasks, &task)

        // 从队列移除
        rdq.client.ZRem(ctx, rdq.key, data)
    }

    return tasks, nil
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "time"

    "delayed"
)

func main() {
    // 优先队列调度器
    queue := delayed.NewPriorityDelayedQueue()

    executor := &delayed.Executor{
        Execute: func(ctx context.Context, task *delayed.DelayedTask) error {
            fmt.Printf("Executing task %s at %v\n", task.ID, time.Now())
            return nil
        },
    }

    scheduler := delayed.NewScheduler(queue, executor)
    scheduler.Start()
    defer scheduler.Stop()

    // 调度延迟任务
    now := time.Now()

    scheduler.Schedule(&delayed.DelayedTask{
        ID:        "task-1",
        Payload:   "email-1",
        ExecuteAt: now.Add(5 * time.Second),
    })

    scheduler.Schedule(&delayed.DelayedTask{
        ID:        "task-2",
        Payload:   "email-2",
        ExecuteAt: now.Add(10 * time.Second),
    })

    scheduler.Schedule(&delayed.DelayedTask{
        ID:        "task-3",
        Payload:   "email-3",
        ExecuteAt: now.Add(3 * time.Second),
    })

    fmt.Println("Scheduled 3 tasks")

    // 等待执行
    time.Sleep(15 * time.Second)
}
```
