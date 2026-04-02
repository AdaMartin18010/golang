# 任务优先级队列 (Task Priority Queue)

> **分类**: 工程与云原生
> **标签**: #priority-queue #heap #scheduling
> **参考**: Linux CFS Scheduler, Priority Queue Algorithms

---

## 优先级队列架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Priority Queue Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Binary Heap (Binary Tree)                         │   │
│  │                                                                      │   │
│  │                              1 (highest)                             │   │
│  │                             / \                                     │   │
│  │                            /   \                                    │   │
│  │                           3     5                                   │   │
│  │                          / \   /                                    │   │
│  │                         7   9 11                                    │   │
│  │                                                                      │   │
│  │   Array representation: [1, 3, 5, 7, 9, 11]                         │   │
│  │   Parent(i) = (i-1)/2, Left(i) = 2i+1, Right(i) = 2i+2              │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Multi-Level Priority Queue                        │   │
│  │                                                                      │   │
│  │   Level 0 (Critical):    ┌─────┐ ┌─────┐ ┌─────┐                    │   │
│  │                          │  P0 │ │  P1 │ │  P2 │  (Immediate)       │   │
│  │                          └─────┘ └─────┘ └─────┘                    │   │
│  │                                                                      │   │
│  │   Level 1 (High):        ┌─────┐ ┌─────┐                            │   │
│  │                          │  P3 │ │  P4 │        (Process next)      │   │
│  │                          └─────┘ └─────┘                            │   │
│  │                                                                      │   │
│  │   Level 2 (Normal):      ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐           │   │
│  │                          │  P5 │ │  P6 │ │  P7 │ │  P8 │            │   │
│  │                          └─────┘ └─────┘ └─────┘ └─────┘           │   │
│  │                                                                      │   │
│  │   Level 3 (Low):         ┌─────┐ ┌─────┐ ...                        │   │
│  │                          │ P9  │ │ P10 │                             │   │
│  │                          └─────┘ └─────┘                            │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Priority Inversion Prevention                     │   │
│  │                                                                      │   │
│  │   Priority Boost: Low-priority task holding lock gets boosted       │   │
│  │   when high-priority task waits                                       │   │
│  │                                                                      │   │
│  │   ┌──────────┐          Lock          ┌──────────┐                 │   │
│  │   │  High    │ ─────────────────────► │   Low    │                 │   │
│  │   │ Priority │         Wait           │ Priority │ ──► Boost!     │   │
│  │   └──────────┘                        └──────────┘                 │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整优先级队列实现

```go
package priority

import (
    "container/heap"
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// Priority 优先级类型
type Priority int

const (
    PriorityCritical Priority = iota // 关键
    PriorityHigh                     // 高
    PriorityNormal                   // 正常
    PriorityLow                      // 低
    PriorityBackground               // 后台
)

func (p Priority) String() string {
    switch p {
    case PriorityCritical:
        return "CRITICAL"
    case PriorityHigh:
        return "HIGH"
    case PriorityNormal:
        return "NORMAL"
    case PriorityLow:
        return "LOW"
    case PriorityBackground:
        return "BACKGROUND"
    default:
        return "UNKNOWN"
    }
}

// PriorityTask 优先级任务
type PriorityTask struct {
    ID          string
    Priority    Priority
    Payload     interface{}
    CreatedAt   time.Time
    Deadline    *time.Time

    // 内部字段
    sequence    int64 // 用于相同优先级的 FIFO
    index       int   // heap 索引
}

// PriorityQueue 优先级队列
type PriorityQueue struct {
    tasks    []*PriorityTask
    sequence int64
    mu       sync.RWMutex
}

// NewPriorityQueue 创建优先级队列
func NewPriorityQueue() *PriorityQueue {
    pq := &PriorityQueue{
        tasks: make([]*PriorityTask, 0),
    }
    heap.Init(pq)
    return pq
}

func (pq *PriorityQueue) Len() int {
    pq.mu.RLock()
    defer pq.mu.RUnlock()
    return len(pq.tasks)
}

func (pq *PriorityQueue) Less(i, j int) bool {
    // 优先级高的在前
    if pq.tasks[i].Priority != pq.tasks[j].Priority {
        return pq.tasks[i].Priority < pq.tasks[j].Priority
    }
    // 相同优先级按时间顺序
    return pq.tasks[i].sequence < pq.tasks[j].sequence
}

func (pq *PriorityQueue) Swap(i, j int) {
    pq.tasks[i], pq.tasks[j] = pq.tasks[j], pq.tasks[i]
    pq.tasks[i].index = i
    pq.tasks[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(pq.tasks)
    task := x.(*PriorityTask)
    task.index = n
    task.sequence = atomic.AddInt64(&pq.sequence, 1)
    pq.tasks = append(pq.tasks, task)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := pq.tasks
    n := len(old)
    task := old[n-1]
    old[n-1] = nil
    task.index = -1
    pq.tasks = old[0 : n-1]
    return task
}

// Enqueue 入队
func (pq *PriorityQueue) Enqueue(task *PriorityTask) {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    heap.Push(pq, task)
}

// Dequeue 出队
func (pq *PriorityQueue) Dequeue() (*PriorityTask, bool) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    if len(pq.tasks) == 0 {
        return nil, false
    }

    task := heap.Pop(pq).(*PriorityTask)
    return task, true
}

// Peek 查看队首
func (pq *PriorityQueue) Peek() (*PriorityTask, bool) {
    pq.mu.RLock()
    defer pq.mu.RUnlock()

    if len(pq.tasks) == 0 {
        return nil, false
    }

    return pq.tasks[0], true
}

// UpdatePriority 更新优先级
func (pq *PriorityQueue) UpdatePriority(taskID string, newPriority Priority) bool {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    for _, task := range pq.tasks {
        if task.ID == taskID {
            task.Priority = newPriority
            heap.Fix(pq, task.index)
            return true
        }
    }

    return false
}

// Remove 移除任务
func (pq *PriorityQueue) Remove(taskID string) bool {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    for i, task := range pq.tasks {
        if task.ID == taskID {
            heap.Remove(pq, i)
            return true
        }
    }

    return false
}

// PriorityScheduler 优先级调度器
type PriorityScheduler struct {
    queues      map[Priority]*PriorityQueue
    weights     map[Priority]int // 权重

    // 执行控制
    workers     int
    executor    Executor

    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

// Executor 执行器接口
type Executor interface {
    Execute(ctx context.Context, task *PriorityTask) error
}

// NewPriorityScheduler 创建优先级调度器
func NewPriorityScheduler(workers int, executor Executor) *PriorityScheduler {
    ctx, cancel := context.WithCancel(context.Background())

    ps := &PriorityScheduler{
        queues:   make(map[Priority]*PriorityQueue),
        weights:  make(map[Priority]int),
        workers:  workers,
        executor: executor,
        ctx:      ctx,
        cancel:   cancel,
    }

    // 初始化各级队列
    for p := PriorityCritical; p <= PriorityBackground; p++ {
        ps.queues[p] = NewPriorityQueue()
    }

    // 默认权重
    ps.weights[PriorityCritical] = 50
    ps.weights[PriorityHigh] = 30
    ps.weights[PriorityNormal] = 15
    ps.weights[PriorityLow] = 4
    ps.weights[PriorityBackground] = 1

    return ps
}

// Submit 提交任务
func (ps *PriorityScheduler) Submit(task *PriorityTask) {
    if task.CreatedAt.IsZero() {
        task.CreatedAt = time.Now()
    }

    ps.queues[task.Priority].Enqueue(task)
}

// Start 启动调度器
func (ps *PriorityScheduler) Start() {
    for i := 0; i < ps.workers; i++ {
        ps.wg.Add(1)
        go ps.worker()
    }
}

// Stop 停止调度器
func (ps *PriorityScheduler) Stop() {
    ps.cancel()
    ps.wg.Wait()
}

func (ps *PriorityScheduler) worker() {
    defer ps.wg.Done()

    for {
        select {
        case <-ps.ctx.Done():
            return
        default:
        }

        task := ps.selectTask()
        if task == nil {
            time.Sleep(10 * time.Millisecond)
            continue
        }

        // 检查截止时间
        if task.Deadline != nil && time.Now().After(*task.Deadline) {
            // 任务超时
            continue
        }

        // 执行任务
        ctx, cancel := context.WithTimeout(ps.ctx, 5*time.Minute)
        _ = ps.executor.Execute(ctx, task)
        cancel()
    }
}

// selectTask 选择任务（加权轮询）
func (ps *PriorityScheduler) selectTask() *PriorityTask {
    for p := PriorityCritical; p <= PriorityBackground; p++ {
        if task, ok := ps.queues[p].Dequeue(); ok {
            return task
        }
    }
    return nil
}

// MultiLevelPriorityQueue 多级优先级队列
type MultiLevelPriorityQueue struct {
    levels      []PriorityQueue
    currentIdx  int
    weights     []int
    counters    []int

    mu          sync.Mutex
}

// NewMultiLevelPriorityQueue 创建多级队列
func NewMultiLevelPriorityQueue(levels int, weights []int) *MultiLevelPriorityQueue {
    mlpq := &MultiLevelPriorityQueue{
        levels:   make([]PriorityQueue, levels),
        weights:  weights,
        counters: make([]int, levels),
    }

    for i := range mlpq.levels {
        mlpq.levels[i] = *NewPriorityQueue()
    }

    return mlpq
}

// Enqueue 入队
func (mlpq *MultiLevelPriorityQueue) Enqueue(level int, task *PriorityTask) {
    if level < 0 || level >= len(mlpq.levels) {
        level = len(mlpq.levels) - 1
    }
    mlpq.levels[level].Enqueue(task)
}

// Dequeue 出队（加权轮询）
func (mlpq *MultiLevelPriorityQueue) Dequeue() (*PriorityTask, bool) {
    mlpq.mu.Lock()
    defer mlpq.mu.Unlock()

    for i := 0; i < len(mlpq.levels); i++ {
        level := (mlpq.currentIdx + i) % len(mlpq.levels)

        mlpq.counters[level]++
        if mlpq.counters[level] >= mlpq.weights[level] {
            mlpq.counters[level] = 0
            mlpq.currentIdx = (level + 1) % len(mlpq.levels)
        }

        if task, ok := mlpq.levels[level].Dequeue(); ok {
            return task, true
n        }
    }

    return nil, false
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

    "priority"
)

func main() {
    // 创建执行器
    executor := &priority.Executor{
        Execute: func(ctx context.Context, task *priority.PriorityTask) error {
            fmt.Printf("[%s] Executing task %s\n", task.Priority, task.ID)
            time.Sleep(100 * time.Millisecond)
            return nil
        },
    }

    // 创建调度器
    scheduler := priority.NewPriorityScheduler(5, executor)
    scheduler.Start()
    defer scheduler.Stop()

    // 提交不同优先级的任务
    priorities := []priority.Priority{
        priority.PriorityLow,
        priority.PriorityCritical,
        priority.PriorityNormal,
        priority.PriorityHigh,
        priority.PriorityCritical,
        priority.PriorityBackground,
    }

    for i, p := range priorities {
        scheduler.Submit(&priority.PriorityTask{
            ID:       fmt.Sprintf("task-%d", i),
            Priority: p,
            Payload:  fmt.Sprintf("data-%d", i),
        })
    }

    // 等待执行
    time.Sleep(2 * time.Second)
}
```
