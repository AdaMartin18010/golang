# 工作池动态伸缩实现 (Worker Pool Dynamic Scaling)

> **分类**: 工程与云原生
> **标签**: #worker-pool #scaling #concurrency #resource-management
> **参考**: Go sync.Pool, Ants Goroutine Pool, Worker Pool Pattern

---

## 动态伸缩架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                   Dynamic Worker Pool Architecture                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Metrics Collector                                 │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │ Queue Depth │  │ Process Rate│  │   Latency   │  │   Errors    │ │   │
│  │  │  (tasks)    │  │  (tps)      │  │   (p99)     │  │   (rate)    │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Auto-Scaler Controller                            │   │
│  │  ┌─────────────────────────────────────────────────────────────┐   │   │
│  │  │                    Scaling Decision Logic                      │   │   │
│  │  │                                                              │   │   │
│  │  │  if queue_depth > scale_up_threshold for scale_up_duration:  │   │   │
│  │  │      scale_up(min_workers, max_workers, scale_step)          │   │   │
│  │  │                                                              │   │   │
│  │  │  if queue_depth < scale_down_threshold for scale_down_duration:│  │   │
│  │  │      scale_down(min_workers, scale_step)                     │   │   │
│  │  │                                                              │   │   │
│  │  │  if worker_idle_time > max_idle_duration:                    │   │   │
│  │  │      scale_down(1)                                           │   │   │
│  │  │                                                              │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Worker Pool                                       │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐        ┌──────────┐       │   │
│  │  │ Worker 1 │  │ Worker 2 │  │ Worker 3 │  ...   │ Worker N │       │   │
│  │  │  (Busy)  │  │  (Idle)  │  │  (Busy)  │        │  (Idle)  │       │   │
│  │  └──────────┘  └──────────┘  └──────────┘        └──────────┘       │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐   │   │
│  │  │                    Task Channel                                │   │   │
│  │  │           ┌─────┬─────┬─────┬─────┬─────┬─────┐              │   │   │
│  │  │  Task In  │ T1  │ T2  │ T3  │ ... │     │     │  Task Out    │   │   │
│  │  │           └─────┴─────┴─────┴─────┴─────┴─────┘              │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心实现

```go
package workerpool

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// Task 任务接口
type Task interface {
    Execute(ctx context.Context) error
}

// TaskFunc 任务函数类型
type TaskFunc func(ctx context.Context) error

func (f TaskFunc) Execute(ctx context.Context) error {
    return f(ctx)
}

// PoolStats 池统计
type PoolStats struct {
    WorkerCount     int32         // 当前工作线程数
    IdleWorkerCount int32         // 空闲工作线程数
    BusyWorkerCount int32         // 忙碌工作线程数
    PendingTasks    int32         // 待处理任务数
    ProcessedTasks  uint64        // 已处理任务数
    ErrorTasks      uint64        // 失败任务数
    AverageLatency  time.Duration // 平均处理延迟
}

// DynamicPool 动态工作池
type DynamicPool struct {
    // 配置
    minWorkers    int           // 最小工作线程数
    maxWorkers    int           // 最大工作线程数
    queueSize     int           // 任务队列大小
    scaleStep     int           // 伸缩步长
    scaleUpThreshold   int32    // 扩容阈值（队列深度）
    scaleDownThreshold int32    // 缩容阈值（队列深度）
    scaleUpDuration    time.Duration
    scaleDownDuration  time.Duration
    maxIdleDuration    time.Duration

    // 状态
    workers       []*Worker
    taskQueue     chan Task

    // 指标
    stats         PoolStats

    // 控制
    ctx           context.Context
    cancel        context.CancelFunc
    wg            sync.WaitGroup

    // 同步
    mu            sync.RWMutex
    running       int32

    // 自动伸缩
    scaler        *AutoScaler
}

// PoolConfig 池配置
type PoolConfig struct {
    MinWorkers         int
    MaxWorkers         int
    QueueSize          int
    ScaleStep          int
    ScaleUpThreshold   int32
    ScaleDownThreshold int32
    ScaleUpDuration    time.Duration
    ScaleDownDuration  time.Duration
    MaxIdleDuration    time.Duration
}

// DefaultConfig 默认配置
var DefaultConfig = PoolConfig{
    MinWorkers:         5,
    MaxWorkers:         100,
    QueueSize:          1000,
    ScaleStep:          5,
    ScaleUpThreshold:   10,
    ScaleDownThreshold: 2,
    ScaleUpDuration:    30 * time.Second,
    ScaleDownDuration:  60 * time.Second,
    MaxIdleDuration:    5 * time.Minute,
}

// NewDynamicPool 创建动态工作池
func NewDynamicPool(config PoolConfig) *DynamicPool {
    ctx, cancel := context.WithCancel(context.Background())

    pool := &DynamicPool{
        minWorkers:         config.MinWorkers,
        maxWorkers:         config.MaxWorkers,
        queueSize:          config.QueueSize,
        scaleStep:          config.ScaleStep,
        scaleUpThreshold:   config.ScaleUpThreshold,
        scaleDownThreshold: config.ScaleDownThreshold,
        scaleUpDuration:    config.ScaleUpDuration,
        scaleDownDuration:  config.ScaleDownDuration,
        maxIdleDuration:    config.MaxIdleDuration,
        taskQueue:          make(chan Task, config.QueueSize),
        ctx:                ctx,
        cancel:             cancel,
    }

    // 创建自动伸缩器
    pool.scaler = NewAutoScaler(pool)

    return pool
}

// Start 启动池
func (p *DynamicPool) Start() error {
    if !atomic.CompareAndSwapInt32(&p.running, 0, 1) {
        return fmt.Errorf("pool already running")
    }

    // 启动最小工作线程
    for i := 0; i < p.minWorkers; i++ {
        p.addWorker()
    }

    // 启动自动伸缩
    p.scaler.Start()

    return nil
}

// Stop 停止池
func (p *DynamicPool) Stop() {
    if !atomic.CompareAndSwapInt32(&p.running, 1, 0) {
        return
    }

    // 停止自动伸缩
    p.scaler.Stop()

    // 取消上下文
    p.cancel()

    // 等待所有工作线程结束
    p.wg.Wait()

    // 关闭任务队列
    close(p.taskQueue)
}

// Submit 提交任务
func (p *DynamicPool) Submit(ctx context.Context, task Task) error {
    if atomic.LoadInt32(&p.running) == 0 {
        return fmt.Errorf("pool not running")
    }

    select {
    case p.taskQueue <- task:
        atomic.AddInt32(&p.stats.PendingTasks, 1)
        return nil
    case <-ctx.Done():
        return ctx.Err()
    default:
        return fmt.Errorf("task queue full")
    }
}

// SubmitWithTimeout 带超时提交
func (p *DynamicPool) SubmitWithTimeout(task Task, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    return p.Submit(ctx, task)
}

// addWorker 添加工作线程
func (p *DynamicPool) addWorker() bool {
    p.mu.Lock()
    defer p.mu.Unlock()

    if len(p.workers) >= p.maxWorkers {
        return false
    }

    worker := NewWorker(p)
    p.workers = append(p.workers, worker)

    p.wg.Add(1)
    go worker.Run()

    atomic.AddInt32(&p.stats.WorkerCount, 1)
    atomic.AddInt32(&p.stats.IdleWorkerCount, 1)

    return true
}

// removeWorker 移除工作线程
func (p *DynamicPool) removeWorker() bool {
    p.mu.Lock()
    defer p.mu.Unlock()

    if len(p.workers) <= p.minWorkers {
        return false
    }

    // 找到空闲的工作线程
    for i, w := range p.workers {
        if w.IsIdle() {
            w.Stop()
            p.workers = append(p.workers[:i], p.workers[i+1:]...)

            atomic.AddInt32(&p.stats.WorkerCount, -1)
            return true
        }
    }

    return false
}

// GetStats 获取统计信息
func (p *DynamicPool) GetStats() PoolStats {
    return PoolStats{
        WorkerCount:     atomic.LoadInt32(&p.stats.WorkerCount),
        IdleWorkerCount: atomic.LoadInt32(&p.stats.IdleWorkerCount),
        BusyWorkerCount: atomic.LoadInt32(&p.stats.BusyWorkerCount),
        PendingTasks:    atomic.LoadInt32(&p.stats.PendingTasks),
        ProcessedTasks:  atomic.LoadUint64(&p.stats.ProcessedTasks),
        ErrorTasks:      atomic.LoadUint64(&p.stats.ErrorTasks),
    }
}

// Worker 工作线程
type Worker struct {
    pool     *DynamicPool
    id       int
    state    int32 // 0: idle, 1: busy
    lastWork time.Time
    stopCh   chan struct{}
}

const (
    WorkerIdle = iota
    WorkerBusy
)

func NewWorker(pool *DynamicPool) *Worker {
    return &Worker{
        pool:     pool,
        id:       int(atomic.AddInt32(&workerIDCounter, 1)),
        state:    WorkerIdle,
        lastWork: time.Now(),
        stopCh:   make(chan struct{}),
    }
}

var workerIDCounter int32

func (w *Worker) Run() {
    defer w.pool.wg.Done()

    for {
        select {
        case <-w.pool.ctx.Done():
            return
        case <-w.stopCh:
            return
        case task, ok := <-w.pool.taskQueue:
            if !ok {
                return
            }

            w.processTask(task)
        }
    }
}

func (w *Worker) processTask(task Task) {
    // 标记为忙碌
    atomic.StoreInt32(&w.state, WorkerBusy)
    atomic.AddInt32(&w.pool.stats.IdleWorkerCount, -1)
    atomic.AddInt32(&w.pool.stats.BusyWorkerCount, 1)
    atomic.AddInt32(&w.pool.stats.PendingTasks, -1)

    start := time.Now()

    // 执行任务
    err := task.Execute(w.pool.ctx)

    // 更新统计
    atomic.AddUint64(&w.pool.stats.ProcessedTasks, 1)
    if err != nil {
        atomic.AddUint64(&w.pool.stats.ErrorTasks, 1)
    }

    // 记录延迟
    latency := time.Since(start)
    // 可以在这里记录到指标系统

    // 标记为空闲
    atomic.StoreInt32(&w.state, WorkerIdle)
    atomic.AddInt32(&w.pool.stats.IdleWorkerCount, 1)
    atomic.AddInt32(&w.pool.stats.BusyWorkerCount, -1)
    w.lastWork = time.Now()
}

func (w *Worker) IsIdle() bool {
    return atomic.LoadInt32(&w.state) == WorkerIdle
}

func (w *Worker) Stop() {
    close(w.stopCh)
}

func (w *Worker) IdleDuration() time.Duration {
    if !w.IsIdle() {
        return 0
    }
    return time.Since(w.lastWork)
}
```

---

## 自动伸缩控制器

```go
package workerpool

import (
    "sync"
    "sync/atomic"
    "time"
)

// AutoScaler 自动伸缩器
type AutoScaler struct {
    pool           *DynamicPool

    // 指标历史
    queueDepthHistory []int32
    historyIndex      int
    historySize       int

    // 伸缩计时器
    scaleUpTimer    int32
    scaleDownTimer  int32

    // 控制
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup

    running int32
}

func NewAutoScaler(pool *DynamicPool) *AutoScaler {
    historySize := int(pool.scaleUpDuration / (5 * time.Second)) + 1

    return &AutoScaler{
        pool:              pool,
        queueDepthHistory: make([]int32, historySize),
        historySize:       historySize,
    }
}

// Start 启动伸缩器
func (s *AutoScaler) Start() {
    if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
        return
    }

    s.ctx, s.cancel = context.WithCancel(s.pool.ctx)

    s.wg.Add(1)
    go s.run()
}

// Stop 停止伸缩器
func (s *AutoScaler) Stop() {
    if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
        return
    }

    s.cancel()
    s.wg.Wait()
}

func (s *AutoScaler) run() {
    defer s.wg.Done()

    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-s.ctx.Done():
            return
        case <-ticker.C:
            s.evaluate()
        }
    }
}

func (s *AutoScaler) evaluate() {
    // 获取当前队列深度
    queueDepth := atomic.LoadInt32(&s.pool.stats.PendingTasks)

    // 记录历史
    s.queueDepthHistory[s.historyIndex] = queueDepth
    s.historyIndex = (s.historyIndex + 1) % s.historySize

    // 检查扩容条件
    if s.shouldScaleUp() {
        s.scaleUp()
    } else if s.shouldScaleDown() {
        s.scaleDown()
    }

    // 检查空闲超时
    s.checkIdleTimeout()
}

func (s *AutoScaler) shouldScaleUp() bool {
    // 检查是否连续超过扩容阈值
    consecutiveCount := 0
    requiredCount := int(s.pool.scaleUpDuration / (5 * time.Second))

    for i := 0; i < s.historySize; i++ {
        idx := (s.historyIndex - 1 - i + s.historySize) % s.historySize
        if s.queueDepthHistory[idx] >= s.pool.scaleUpThreshold {
            consecutiveCount++
        } else {
            break
        }
    }

    return consecutiveCount >= requiredCount
}

func (s *AutoScaler) shouldScaleDown() bool {
    // 检查是否连续低于缩容阈值
    consecutiveCount := 0
    requiredCount := int(s.pool.scaleDownDuration / (5 * time.Second))

    for i := 0; i < s.historySize; i++ {
        idx := (s.historyIndex - 1 - i + s.historySize) % s.historySize
        if s.queueDepthHistory[idx] <= s.pool.scaleDownThreshold {
            consecutiveCount++
        } else {
            break
        }
    }

    return consecutiveCount >= requiredCount
}

func (s *AutoScaler) scaleUp() {
    for i := 0; i < s.pool.scaleStep; i++ {
        if !s.pool.addWorker() {
            break // 达到最大限制
        }
    }
}

func (s *AutoScaler) scaleDown() {
    for i := 0; i < s.pool.scaleStep; i++ {
        if !s.pool.removeWorker() {
            break // 达到最小限制
        }
    }
}

func (s *AutoScaler) checkIdleTimeout() {
    s.pool.mu.RLock()
    workers := make([]*Worker, len(s.pool.workers))
    copy(workers, s.pool.workers)
    s.pool.mu.RUnlock()

    for _, w := range workers {
        if w.IsIdle() && w.IdleDuration() > s.pool.maxIdleDuration {
            // 移除空闲超时的工作线程
            if s.pool.removeWorker() {
                break // 一次只移除一个
            }
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
    "fmt"
    "time"

    "workerpool"
)

func main() {
    // 创建动态工作池
    config := workerpool.PoolConfig{
        MinWorkers:         5,
        MaxWorkers:         50,
        QueueSize:          1000,
        ScaleStep:          5,
        ScaleUpThreshold:   20,
        ScaleDownThreshold: 5,
        ScaleUpDuration:    30 * time.Second,
        ScaleDownDuration:  60 * time.Second,
        MaxIdleDuration:    5 * time.Minute,
    }

    pool := workerpool.NewDynamicPool(config)

    // 启动池
    if err := pool.Start(); err != nil {
        panic(err)
    }
    defer pool.Stop()

    // 提交任务
    for i := 0; i < 1000; i++ {
        taskID := i
        task := workerpool.TaskFunc(func(ctx context.Context) error {
            fmt.Printf("Processing task %d\n", taskID)
            time.Sleep(100 * time.Millisecond) // 模拟工作
            return nil
        })

        if err := pool.Submit(context.Background(), task); err != nil {
            fmt.Printf("Failed to submit task %d: %v\n", taskID, err)
        }
    }

    // 查看统计
    time.Sleep(1 * time.Second)
    stats := pool.GetStats()
    fmt.Printf("Workers: %d, Idle: %d, Busy: %d, Pending: %d, Processed: %d\n",
        stats.WorkerCount, stats.IdleWorkerCount, stats.BusyWorkerCount,
        stats.PendingTasks, stats.ProcessedTasks)

    // 等待完成
    time.Sleep(30 * time.Second)
}
```

---

## 高级功能

```go
package workerpool

import (
    "container/heap"
    "context"
)

// PriorityPool 优先级工作池
type PriorityPool struct {
    *DynamicPool
    priorityQueue *PriorityTaskQueue
}

// PriorityTask 带优先级的任务
type PriorityTask struct {
    Task     Task
    Priority int // 数字越小优先级越高
    index    int // 堆索引
}

// PriorityTaskQueue 优先级任务队列
type PriorityTaskQueue []*PriorityTask

func (pq PriorityTaskQueue) Len() int { return len(pq) }
func (pq PriorityTaskQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }
func (pq PriorityTaskQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityTaskQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*PriorityTask)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityTaskQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil
    item.index = -1
    *pq = old[0 : n-1]
    return item
}

// SubmitPriority 提交优先级任务
func (p *PriorityPool) SubmitPriority(ctx context.Context, task Task, priority int) error {
    pt := &PriorityTask{
        Task:     task,
        Priority: priority,
    }

    heap.Push(p.priorityQueue, pt)

    // 通知工作线程
    select {
    case p.taskQueue <- &priorityTaskWrapper{pt}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

type priorityTaskWrapper struct {
    *PriorityTask
}

func (w *priorityTaskWrapper) Execute(ctx context.Context) error {
    return w.Task.Execute(ctx)
}
```
