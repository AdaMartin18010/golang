# 并发优化策略 - 无锁算法与Worker池

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [并发优化策略 - 无锁算法与Worker池](#并发优化策略---无锁算法与worker池)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
    - [1.1 并发优化目标](#11-并发优化目标)
  - [2. 无锁算法](#2-无锁算法)
    - [2.1 原子操作](#21-原子操作)
    - [2.2 无锁队列](#22-无锁队列)
    - [2.3 无锁栈](#23-无锁栈)
  - [3. Worker池设计](#3-worker池设计)
    - [3.1 基础Worker池](#31-基础worker池)
    - [3.2 动态Worker池](#32-动态worker池)
  - [4. 并发模式](#4-并发模式)
    - [4.1 Fan-Out/Fan-In](#41-fan-outfan-in)
    - [4.2 Pipeline模式](#42-pipeline模式)
  - [5. 最佳实践](#5-最佳实践)
    - [5.1 并发优化清单](#51-并发优化清单)
    - [5.2 性能对比](#52-性能对比)

---

## 1. 概述

### 1.1 并发优化目标

```text
并发优化四大目标:

┌─────────────────────────────────────┐
│         并发优化目标                │
├─────────────────────────────────────┤
│                                     │
│  1. 提升吞吐量                      │
│     └─ 充分利用多核CPU              │
│                                     │
│  2. 降低延迟                        │
│     └─ 减少锁竞争和等待             │
│                                     │
│  3. 避免竞态条件                    │
│     └─ 保证数据一致性               │
│                                     │
│  4. 资源高效利用                    │
│     └─ Goroutine池化和复用          │
│                                     │
└─────────────────────────────────────┘
```

---

## 2. 无锁算法

### 2.1 原子操作

```go
// pkg/lockfree/atomic.go

package lockfree

import (
    "sync/atomic"
    "unsafe"
)

// AtomicCounter 原子计数器
type AtomicCounter struct {
    value int64
}

// Inc 增加
func (c *AtomicCounter) Inc() int64 {
    return atomic.AddInt64(&c.value, 1)
}

// Dec 减少
func (c *AtomicCounter) Dec() int64 {
    return atomic.AddInt64(&c.value, -1)
}

// Get 获取值
func (c *AtomicCounter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}

// Set 设置值
func (c *AtomicCounter) Set(val int64) {
    atomic.StoreInt64(&c.value, val)
}

// CompareAndSwap CAS操作
func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
    return atomic.CompareAndSwapInt64(&c.value, old, new)
}
```

---

### 2.2 无锁队列

```go
// pkg/lockfree/queue.go

package lockfree

import (
    "sync/atomic"
    "unsafe"
)

// LockFreeQueue 无锁队列
type LockFreeQueue struct {
    head unsafe.Pointer
    tail unsafe.Pointer
}

type node struct {
    value interface{}
    next  unsafe.Pointer
}

// NewLockFreeQueue 创建无锁队列
func NewLockFreeQueue() *LockFreeQueue {
    n := unsafe.Pointer(&node{})
    return &LockFreeQueue{
        head: n,
        tail: n,
    }
}

// Enqueue 入队
func (q *LockFreeQueue) Enqueue(value interface{}) {
    n := &node{value: value}

    for {
        tail := load(&q.tail)
        next := load(&tail.next)

        if tail == load(&q.tail) {
            if next == nil {
                if cas(&tail.next, next, n) {
                    cas(&q.tail, tail, n)
                    return
                }
            } else {
                cas(&q.tail, tail, next)
            }
        }
    }
}

// Dequeue 出队
func (q *LockFreeQueue) Dequeue() (interface{}, bool) {
    for {
        head := load(&q.head)
        tail := load(&q.tail)
        next := load(&head.next)

        if head == load(&q.head) {
            if head == tail {
                if next == nil {
                    return nil, false
                }
                cas(&q.tail, tail, next)
            } else {
                value := next.value
                if cas(&q.head, head, next) {
                    return value, true
                }
            }
        }
    }
}

// 辅助函数
func load(p *unsafe.Pointer) *node {
    return (*node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *node) bool {
    return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
```

---

### 2.3 无锁栈

```go
// pkg/lockfree/stack.go

package lockfree

import (
    "sync/atomic"
    "unsafe"
)

// LockFreeStack 无锁栈
type LockFreeStack struct {
    head unsafe.Pointer
}

type stackNode struct {
    value interface{}
    next  unsafe.Pointer
}

// NewLockFreeStack 创建无锁栈
func NewLockFreeStack() *LockFreeStack {
    return &LockFreeStack{}
}

// Push 压栈
func (s *LockFreeStack) Push(value interface{}) {
    n := &stackNode{value: value}

    for {
        old := atomic.LoadPointer(&s.head)
        n.next = old
        if atomic.CompareAndSwapPointer(&s.head, old, unsafe.Pointer(n)) {
            return
        }
    }
}

// Pop 弹栈
func (s *LockFreeStack) Pop() (interface{}, bool) {
    for {
        old := atomic.LoadPointer(&s.head)
        if old == nil {
            return nil, false
        }

        node := (*stackNode)(old)
        next := atomic.LoadPointer(&node.next)

        if atomic.CompareAndSwapPointer(&s.head, old, next) {
            return node.value, true
        }
    }
}
```

---

## 3. Worker池设计

### 3.1 基础Worker池

```go
// pkg/worker/pool.go

package worker

import (
    "context"
    "fmt"
    "sync"
)

// Task 任务接口
type Task interface {
    Execute(ctx Context.Context) error
}

// WorkerPool Worker池
type WorkerPool struct {
    workers    int
    taskQueue  Channel Task
    wg         sync.WaitGroup
    ctx        Context.Context
    cancel     Context.CancelFunc
    stats      *PoolStats
}

// PoolStats 池统计
type PoolStats struct {
    mu             sync.RWMutex
    tasksSubmitted int64
    tasksCompleted int64
    tasksFailed    int64
}

// NewWorkerPool 创建Worker池
func NewWorkerPool(workers, queueSize int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())

    return &WorkerPool{
        workers:   workers,
        taskQueue: make(Channel Task, queueSize),
        ctx:       ctx,
        cancel:    cancel,
        stats:     &PoolStats{},
    }
}

// Start 启动Worker池
func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

// worker Worker Goroutine
func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()

    for {
        select {
        case task, ok := <-p.taskQueue:
            if !ok {
                return
            }

            if err := task.Execute(p.ctx); err != nil {
                p.stats.recordFailed()
            } else {
                p.stats.recordCompleted()
            }

        case <-p.ctx.Done():
            return
        }
    }
}

// Submit 提交任务
func (p *WorkerPool) Submit(task Task) error {
    select {
    case p.taskQueue <- task:
        p.stats.recordSubmitted()
        return nil
    case <-p.ctx.Done():
        return fmt.Errorf("pool is closed")
    }
}

// Stop 停止Worker池
func (p *WorkerPool) Stop() {
    p.cancel()
    close(p.taskQueue)
    p.wg.Wait()
}

// Stats 获取统计信息
func (p *WorkerPool) Stats() (submitted, completed, failed int64) {
    p.stats.mu.RLock()
    defer p.stats.mu.RUnlock()

    return p.stats.tasksSubmitted, p.stats.tasksCompleted, p.stats.tasksFailed
}

func (s *PoolStats) recordSubmitted() {
    s.mu.Lock()
    s.tasksSubmitted++
    s.mu.Unlock()
}

func (s *PoolStats) recordCompleted() {
    s.mu.Lock()
    s.tasksCompleted++
    s.mu.Unlock()
}

func (s *PoolStats) recordFailed() {
    s.mu.Lock()
    s.tasksFailed++
    s.mu.Unlock()
}
```

---

### 3.2 动态Worker池

```go
// pkg/worker/dynamic_pool.go

package worker

import (
    "context"
    "sync"
    "sync/atomic"
    "time"
)

// DynamicPool 动态Worker池
type DynamicPool struct {
    minWorkers    int
    maxWorkers    int
    currentWorkers int32
    taskQueue     Channel Task
    wg            sync.WaitGroup
    ctx           Context.Context
    cancel        Context.CancelFunc
    scaleInterval time.Duration
}

// NewDynamicPool 创建动态Worker池
func NewDynamicPool(min, max, queueSize int) *DynamicPool {
    ctx, cancel := context.WithCancel(context.Background())

    return &DynamicPool{
        minWorkers:    min,
        maxWorkers:    max,
        currentWorkers: 0,
        taskQueue:     make(Channel Task, queueSize),
        ctx:           ctx,
        cancel:        cancel,
        scaleInterval: 5 * time.Second,
    }
}

// Start 启动动态Worker池
func (p *DynamicPool) Start() {
    // 启动最小数量的worker
    for i := 0; i < p.minWorkers; i++ {
        p.addWorker()
    }

    // 启动自动缩放
    go p.autoScale()
}

// addWorker 添加worker
func (p *DynamicPool) addWorker() {
    current := atomic.LoadInt32(&p.currentWorkers)
    if int(current) >= p.maxWorkers {
        return
    }

    atomic.AddInt32(&p.currentWorkers, 1)
    p.wg.Add(1)

    go func() {
        defer p.wg.Done()
        defer atomic.AddInt32(&p.currentWorkers, -1)

        idleCount := 0
        maxIdle := 5

        for {
            select {
            case task, ok := <-p.taskQueue:
                if !ok {
                    return
                }

                idleCount = 0
                task.Execute(p.ctx)

            case <-time.After(time.Second):
                idleCount++

                // 如果空闲太久且超过最小worker数，退出
                current := atomic.LoadInt32(&p.currentWorkers)
                if idleCount >= maxIdle && int(current) > p.minWorkers {
                    return
                }

            case <-p.ctx.Done():
                return
            }
        }
    }()
}

// autoScale 自动缩放
func (p *DynamicPool) autoScale() {
    ticker := time.NewTicker(p.scaleInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            queueLen := len(p.taskQueue)
            currentWorkers := int(atomic.LoadInt32(&p.currentWorkers))

            // 队列积压，增加worker
            if queueLen > currentWorkers && currentWorkers < p.maxWorkers {
                needed := (queueLen - currentWorkers) / 2
                if needed < 1 {
                    needed = 1
                }

                for i := 0; i < needed && currentWorkers+i < p.maxWorkers; i++ {
                    p.addWorker()
                }
            }

        case <-p.ctx.Done():
            return
        }
    }
}

// Submit 提交任务
func (p *DynamicPool) Submit(task Task) error {
    select {
    case p.taskQueue <- task:
        return nil
    case <-p.ctx.Done():
        return fmt.Errorf("pool is closed")
    }
}

// Stop 停止动态Worker池
func (p *DynamicPool) Stop() {
    p.cancel()
    close(p.taskQueue)
    p.wg.Wait()
}

// WorkerCount 获取当前worker数量
func (p *DynamicPool) WorkerCount() int {
    return int(atomic.LoadInt32(&p.currentWorkers))
}
```

---

## 4. 并发模式

### 4.1 Fan-Out/Fan-In

```go
// pkg/patterns/fan.go

package patterns

import (
    "context"
    "sync"
)

// FanOut 扇出模式
func FanOut(ctx Context.Context, input <-Channel interface{}, workers int, process func(interface{}) interface{}) []<-Channel interface{} {
    outputs := make([]<-Channel interface{}, workers)

    for i := 0; i < workers; i++ {
        outputs[i] = worker(ctx, input, process)
    }

    return outputs
}

func worker(ctx Context.Context, input <-Channel interface{}, process func(interface{}) interface{}) <-Channel interface{} {
    output := make(Channel interface{})

    go func() {
        defer close(output)

        for {
            select {
            case data, ok := <-input:
                if !ok {
                    return
                }

                result := process(data)

                select {
                case output <- result:
                case <-ctx.Done():
                    return
                }

            case <-ctx.Done():
                return
            }
        }
    }()

    return output
}

// FanIn 扇入模式
func FanIn(ctx Context.Context, inputs ...<-Channel interface{}) <-Channel interface{} {
    output := make(Channel interface{})
    var wg sync.WaitGroup

    multiplex := func(input <-Channel interface{}) {
        defer wg.Done()

        for {
            select {
            case data, ok := <-input:
                if !ok {
                    return
                }

                select {
                case output <- data:
                case <-ctx.Done():
                    return
                }

            case <-ctx.Done():
                return
            }
        }
    }

    wg.Add(len(inputs))
    for _, input := range inputs {
        go multiplex(input)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}
```

---

### 4.2 Pipeline模式

```go
// pkg/patterns/pipeline.go

package patterns

import "context"

// Stage 管道阶段
type Stage func(Context.Context, <-Channel interface{}) <-Channel interface{}

// Pipeline 创建管道
func Pipeline(ctx Context.Context, input <-Channel interface{}, stages ...Stage) <-Channel interface{} {
    output := input

    for _, stage := range stages {
        output = stage(ctx, output)
    }

    return output
}

// 示例阶段：过滤
func FilterStage(predicate func(interface{}) bool) Stage {
    return func(ctx Context.Context, input <-Channel interface{}) <-Channel interface{} {
        output := make(Channel interface{})

        go func() {
            defer close(output)

            for {
                select {
                case data, ok := <-input:
                    if !ok {
                        return
                    }

                    if predicate(data) {
                        select {
                        case output <- data:
                        case <-ctx.Done():
                            return
                        }
                    }

                case <-ctx.Done():
                    return
                }
            }
        }()

        return output
    }
}

// 示例阶段：转换
func MapStage(transform func(interface{}) interface{}) Stage {
    return func(ctx Context.Context, input <-Channel interface{}) <-Channel interface{} {
        output := make(Channel interface{})

        go func() {
            defer close(output)

            for {
                select {
                case data, ok := <-input:
                    if !ok {
                        return
                    }

                    result := transform(data)

                    select {
                    case output <- result:
                    case <-ctx.Done():
                        return
                    }

                case <-ctx.Done():
                    return
                }
            }
        }()

        return output
    }
}
```

---

## 5. 最佳实践

### 5.1 并发优化清单

```text
✅ 并发优化检查清单:

□ 1. 避免过度并发
   - 根据CPU核心数设置worker数量
   - 使用runtime.GOMAXPROCS

□ 2. 减少锁竞争
   - 缩小临界区
   - 使用读写锁
   - 考虑无锁算法

□ 3. Channel使用
   - 适当的缓冲区大小
   - 及时关闭channel
   - 避免阻塞

□ 4. Goroutine池化
   - 使用Worker池
   - 避免无限创建goroutine

□ 5. 上下文管理
   - 使用context控制生命周期
   - 传递取消信号

□ 6. 避免数据竞争
   - 使用-race检测
   - 合理使用sync包

□ 7. 性能监控
   - 监控goroutine数量
   - 监控锁竞争

□ 8. 资源清理
   - 使用defer
   - 正确处理panic
```

---

### 5.2 性能对比
