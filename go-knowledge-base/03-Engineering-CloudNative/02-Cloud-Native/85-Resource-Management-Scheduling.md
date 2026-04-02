# 资源管理与调度 (Resource Management & Scheduling)

> **分类**: 工程与云原生
> **标签**: #resource-management #scheduling #pool #limiting
> **参考**: Kubernetes Scheduler, Linux Cgroups

---

## 资源调度架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Resource Management & Scheduling                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Resource Types                                    │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │   CPU    │  │  Memory  │  │  Network │  │   Disk   │            │   │
│  │  │ (cores)  │  │   (GB)   │  │ (MB/s)   │  │  (IOPS)  │            │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │   │
│  │                                                                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐                          │   │
│  │  │ Goroutine│  │  File    │  │ External │                          │   │
│  │  │   Pool   │  │ Descriptor│  │  API     │                          │   │
│  │  └──────────┘  └──────────┘  └──────────┘                          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│  ┌─────────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Scheduling Policies                               │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │   FIFO   │  │   LIFO   │  │  Priority│  │Weighted  │            │   │
│  │  │          │  │          │  │          │  │ Fair     │            │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │   │
│  │                                                                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐                          │   │
│  │  │RoundRobin│  │LeastConn │  │ Resource │                          │   │
│  │  │          │  │          │  │Based     │                          │   │
│  │  └──────────┘  └──────────┘  └──────────┘                          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 资源管理器实现

```go
package resource

import (
    "context"
    "fmt"
    "runtime"
    "sync"
    "sync/atomic"
    "time"
)

// ResourceType 资源类型
type ResourceType string

const (
    ResourceCPU        ResourceType = "cpu"
    ResourceMemory     ResourceType = "memory"
    ResourceNetwork    ResourceType = "network"
    ResourceDisk       ResourceType = "disk"
    ResourceGoroutine  ResourceType = "goroutine"
    ResourceConnection ResourceType = "connection"
)

// ResourceQuota 资源配额
type ResourceQuota struct {
    CPU        float64 // CPU核心数
    Memory     int64   // 内存字节数
    Network    int64   // 网络带宽 B/s
    Disk       int64   // 磁盘IOPS
    Goroutines int     // Goroutine数量
    Connections int    // 连接数
}

// ResourceUsage 资源使用
type ResourceUsage struct {
    CPU         float64
    Memory      int64
    Network     int64
    Disk        int64
    Goroutines  int
    Connections int
    Timestamp   time.Time
}

// ResourceManager 资源管理器
type ResourceManager struct {
    quota   ResourceQuota
    usage   ResourceUsage

    // 限制器
    limiters map[ResourceType]Limiter

    // 监控
    mu      sync.RWMutex
    history []ResourceUsage
}

// Limiter 资源限制器接口
type Limiter interface {
    Acquire(ctx context.Context, amount float64) error
    Release(amount float64)
    Available() float64
}

// NewResourceManager 创建资源管理器
func NewResourceManager(quota ResourceQuota) *ResourceManager {
    rm := &ResourceManager{
        quota:    quota,
        limiters: make(map[ResourceType]Limiter),
        history:  make([]ResourceUsage, 0, 100),
    }

    // 初始化限制器
    rm.limiters[ResourceCPU] = NewCPULimiter(quota.CPU)
    rm.limiters[ResourceMemory] = NewMemoryLimiter(quota.Memory)
    rm.limiters[ResourceGoroutine] = NewGoroutineLimiter(quota.Goroutines)

    // 启动监控
    go rm.monitor()

    return rm
}

// Acquire 获取资源
func (rm *ResourceManager) Acquire(ctx context.Context, resType ResourceType, amount float64) error {
    limiter, ok := rm.limiters[resType]
    if !ok {
        return fmt.Errorf("unknown resource type: %s", resType)
    }

    return limiter.Acquire(ctx, amount)
}

// Release 释放资源
func (rm *ResourceManager) Release(resType ResourceType, amount float64) {
    limiter, ok := rm.limiters[resType]
    if !ok {
        return
    }

    limiter.Release(amount)
}

// GetUsage 获取当前使用
func (rm *ResourceManager) GetUsage() ResourceUsage {
    rm.mu.RLock()
    defer rm.mu.RUnlock()
    return rm.usage
}

// GetQuota 获取配额
func (rm *ResourceManager) GetQuota() ResourceQuota {
    return rm.quota
}

// monitor 监控资源使用
func (rm *ResourceManager) monitor() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        usage := rm.collectUsage()

        rm.mu.Lock()
        rm.usage = usage
        rm.history = append(rm.history, usage)
        if len(rm.history) > 100 {
            rm.history = rm.history[1:]
        }
        rm.mu.Unlock()
    }
}

func (rm *ResourceManager) collectUsage() ResourceUsage {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return ResourceUsage{
        Memory:     int64(m.Alloc),
        Goroutines: runtime.NumGoroutine(),
        Timestamp:  time.Now(),
    }
}

// TokenBucketLimiter 令牌桶限制器
type TokenBucketLimiter struct {
    capacity   float64
    tokens     float64
    fillRate   float64
    lastFill   time.Time
    mu         sync.Mutex
}

// NewCPULimiter 创建CPU限制器
func NewCPULimiter(cores float64) Limiter {
    return &TokenBucketLimiter{
        capacity: cores,
        tokens:   cores,
        fillRate: cores, // 每秒填充核心数
        lastFill: time.Now(),
    }
}

// NewMemoryLimiter 创建内存限制器
func NewMemoryLimiter(bytes int64) Limiter {
    return &TokenBucketLimiter{
        capacity: float64(bytes),
        tokens:   float64(bytes),
        fillRate: 0, // 内存不自动填充
        lastFill: time.Now(),
    }
}

// NewGoroutineLimiter 创建Goroutine限制器
func NewGoroutineLimiter(max int) Limiter {
    return &TokenBucketLimiter{
        capacity: float64(max),
        tokens:   float64(max),
        fillRate: float64(max),
        lastFill: time.Now(),
    }
}

func (tbl *TokenBucketLimiter) Acquire(ctx context.Context, amount float64) error {
    tbl.mu.Lock()
    defer tbl.mu.Unlock()

    // 填充令牌
    tbl.refill()

    if tbl.tokens >= amount {
        tbl.tokens -= amount
        return nil
    }

    // 需要等待
    needed := amount - tbl.tokens
    waitTime := time.Duration(needed / tbl.fillRate * float64(time.Second))

    tbl.mu.Unlock()

    select {
    case <-time.After(waitTime):
        tbl.mu.Lock()
        tbl.refill()
        if tbl.tokens >= amount {
            tbl.tokens -= amount
            return nil
        }
        return fmt.Errorf("insufficient tokens")
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (tbl *TokenBucketLimiter) Release(amount float64) {
    tbl.mu.Lock()
    defer tbl.mu.Unlock()

    tbl.tokens += amount
    if tbl.tokens > tbl.capacity {
        tbl.tokens = tbl.capacity
    }
}

func (tbl *TokenBucketLimiter) Available() float64 {
    tbl.mu.Lock()
    defer tbl.mu.Unlock()
    tbl.refill()
    return tbl.tokens
}

func (tbl *TokenBucketLimiter) refill() {
    if tbl.fillRate <= 0 {
        return
    }

    now := time.Now()
    elapsed := now.Sub(tbl.lastFill).Seconds()

    tokensToAdd := elapsed * tbl.fillRate
    if tokensToAdd > 0 {
        tbl.tokens += tokensToAdd
        if tbl.tokens > tbl.capacity {
            tbl.tokens = tbl.capacity
        }
        tbl.lastFill = now
    }
}
```

---

## 资源感知调度器

```go
package resource

import (
    "container/heap"
    "context"
    "fmt"
    "sync"
)

// Task 调度任务
type Task struct {
    ID       string
    Priority int
    Resource ResourceRequirement
    Handler  func(ctx context.Context) error
}

// ResourceRequirement 资源需求
type ResourceRequirement struct {
    CPU      float64
    Memory   int64
    Duration time.Duration
}

// ResourceAwareScheduler 资源感知调度器
type ResourceAwareScheduler struct {
    rm         *ResourceManager

    // 任务队列
    readyQueue PriorityQueue
    pendingQueue []*Task

    // 执行控制
    maxConcurrent int
    running       int32

    // 生命周期
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
    mu     sync.Mutex
}

// PriorityQueue 优先队列
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Priority > pq[j].Priority }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
    *pq = append(*pq, x.(*Task))
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[0 : n-1]
    return item
}

// NewResourceAwareScheduler 创建调度器
func NewResourceAwareScheduler(rm *ResourceManager, maxConcurrent int) *ResourceAwareScheduler {
    ctx, cancel := context.WithCancel(context.Background())

    return &ResourceAwareScheduler{
        rm:            rm,
        maxConcurrent: maxConcurrent,
        ctx:           ctx,
        cancel:        cancel,
    }
}

// Submit 提交任务
func (ras *ResourceAwareScheduler) Submit(task *Task) error {
    ras.mu.Lock()
    defer ras.mu.Unlock()

    // 检查资源是否足够
    if ras.canSchedule(task) {
        heap.Push(&ras.readyQueue, task)
        ras.schedule()
    } else {
        ras.pendingQueue = append(ras.pendingQueue, task)
    }

    return nil
}

// canSchedule 检查是否可以调度
func (ras *ResourceAwareScheduler) canSchedule(task *Task) bool {
    // 检查当前资源使用情况
    usage := ras.rm.GetUsage()
    quota := ras.rm.GetQuota()

    // 检查并发限制
    if int(atomic.LoadInt32(&ras.running)) >= ras.maxConcurrent {
        return false
    }

    // 检查资源
    if usage.Memory+task.Resource.Memory > quota.Memory {
        return false
    }

    return true
}

// schedule 调度任务
func (ras *ResourceAwareScheduler) schedule() {
    for atomic.LoadInt32(&ras.running) < int32(ras.maxConcurrent) && ras.readyQueue.Len() > 0 {
        task := heap.Pop(&ras.readyQueue).(*Task)

        atomic.AddInt32(&ras.running, 1)
        ras.wg.Add(1)

        go ras.executeTask(task)
    }
}

func (ras *ResourceAwareScheduler) executeTask(task *Task) {
    defer ras.wg.Done()
    defer atomic.AddInt32(&ras.running, -1)

    // 获取资源
    ctx, cancel := context.WithTimeout(ras.ctx, task.Resource.Duration)
    defer cancel()

    if err := ras.rm.Acquire(ctx, ResourceCPU, task.Resource.CPU); err != nil {
        // 资源获取失败，重新入队
        ras.Submit(task)
        return
    }
    defer ras.rm.Release(ResourceCPU, task.Resource.CPU)

    if err := ras.rm.Acquire(ctx, ResourceMemory, float64(task.Resource.Memory)); err != nil {
        ras.rm.Release(ResourceCPU, task.Resource.CPU)
        ras.Submit(task)
        return
    }
    defer ras.rm.Release(ResourceMemory, float64(task.Resource.Memory))

    // 执行任务
    _ = task.Handler(ctx)

    // 检查等待队列
    ras.checkPendingQueue()
}

func (ras *ResourceAwareScheduler) checkPendingQueue() {
    ras.mu.Lock()
    defer ras.mu.Unlock()

    var stillPending []*Task
    for _, task := range ras.pendingQueue {
        if ras.canSchedule(task) {
            heap.Push(&ras.readyQueue, task)
        } else {
            stillPending = append(stillPending, task)
        }
    }

    ras.pendingQueue = stillPending
    ras.schedule()
}

// Stop 停止调度器
func (ras *ResourceAwareScheduler) Stop() {
    ras.cancel()
    ras.wg.Wait()
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

    "resource"
)

func main() {
    // 创建资源管理器
    quota := resource.ResourceQuota{
        CPU:        4.0,
        Memory:     1024 * 1024 * 1024, // 1GB
        Goroutines: 1000,
    }

    rm := resource.NewResourceManager(quota)

    // 创建调度器
    scheduler := resource.NewResourceAwareScheduler(rm, 10)

    // 提交任务
    for i := 0; i < 20; i++ {
        task := &resource.Task{
            ID:       fmt.Sprintf("task-%d", i),
            Priority: 20 - i,
            Resource: resource.ResourceRequirement{
                CPU:      0.5,
                Memory:   100 * 1024 * 1024,
                Duration: 5 * time.Second,
            },
            Handler: func(ctx context.Context) error {
                fmt.Println("Executing task...")
                time.Sleep(1 * time.Second)
                return nil
            },
        }

        scheduler.Submit(task)
    }

    // 等待
    time.Sleep(30 * time.Second)

    scheduler.Stop()
}
```
