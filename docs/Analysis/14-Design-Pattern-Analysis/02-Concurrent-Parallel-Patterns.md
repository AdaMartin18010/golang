# 14.1 并发并行模式分析

<!-- TOC START -->
- [14.1 并发并行模式分析](#并发并行模式分析)
  - [14.1.1 概述](#概述)
  - [14.1.2 1. 活动对象模式 (Active Object)](#1-活动对象模式-active-object)
    - [14.1.2.1 定义](#定义)
    - [14.1.2.2 形式化定义](#形式化定义)
    - [14.1.2.3 Golang实现](#golang实现)
    - [14.1.2.4 性能分析](#性能分析)
  - [14.1.3 2. 管程模式 (Monitor)](#2-管程模式-monitor)
    - [14.1.3.1 定义](#定义)
    - [14.1.3.2 Golang实现](#golang实现)
  - [14.1.4 3. 线程池模式 (Thread Pool)](#3-线程池模式-thread-pool)
    - [14.1.4.1 定义](#定义)
    - [14.1.4.2 Golang实现](#golang实现)
  - [14.1.5 4. 生产者-消费者模式 (Producer-Consumer)](#4-生产者-消费者模式-producer-consumer)
    - [14.1.5.1 定义](#定义)
    - [14.1.5.2 Golang实现](#golang实现)
  - [14.1.6 5. 读写锁模式 (Readers-Writer Lock)](#5-读写锁模式-readers-writer-lock)
    - [14.1.6.1 定义](#定义)
    - [14.1.6.2 Golang实现](#golang实现)
  - [14.1.7 6. Future/Promise模式](#6-futurepromise模式)
    - [14.1.7.1 定义](#定义)
    - [14.1.7.2 Golang实现](#golang实现)
  - [14.1.8 7. Actor模型](#7-actor模型)
    - [14.1.8.1 定义](#定义)
    - [14.1.8.2 Golang实现](#golang实现)
  - [14.1.9 8. 性能分析](#8-性能分析)
    - [14.1.9.1 模式性能对比](#模式性能对比)
    - [14.1.9.2 内存使用分析](#内存使用分析)
    - [14.1.9.3 并发性能分析](#并发性能分析)
  - [14.1.10 9. 最佳实践](#9-最佳实践)
    - [14.1.10.1 模式选择原则](#模式选择原则)
    - [14.1.10.2 实现建议](#实现建议)
    - [14.1.10.3 常见陷阱](#常见陷阱)
  - [14.1.11 10. 应用场景](#10-应用场景)
    - [14.1.11.1 活动对象模式](#活动对象模式)
    - [14.1.11.2 管程模式](#管程模式)
    - [14.1.11.3 线程池模式](#线程池模式)
    - [14.1.11.4 生产者-消费者模式](#生产者-消费者模式)
    - [14.1.11.5 读写锁模式](#读写锁模式)
    - [14.1.11.6 Future/Promise模式](#futurepromise模式)
    - [14.1.11.7 Actor模式](#actor模式)
  - [14.1.12 11. 总结](#11-总结)
    - [14.1.12.1 关键优势](#关键优势)
    - [14.1.12.2 成功要素](#成功要素)
<!-- TOC END -->

## 14.1.1 概述

并发并行模式是处理多线程、异步编程和分布式计算的重要设计模式。本文档基于Golang技术栈，深入分析各种并发并行模式的设计、实现和性能特征。

## 14.1.2 1. 活动对象模式 (Active Object)

### 14.1.2.1 定义

将方法调用与执行分离，使方法调用异步执行，避免阻塞调用者。

### 14.1.2.2 形式化定义

$$\text{ActiveObject} = (M, Q, E, R)$$

其中：

- $M$ 是方法集合
- $Q$ 是请求队列
- $E$ 是执行器
- $R$ 是结果处理器

### 14.1.2.3 Golang实现

```go
package activeobject

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Request 请求结构
type Request struct {
    ID       string
    Method   string
    Params   map[string]interface{}
    Response chan interface{}
    Error    chan error
}

// ActiveObject 活动对象
type ActiveObject struct {
    requestQueue chan *Request
    worker       *Worker
    ctx          context.Context
    cancel       context.CancelFunc
    wg           sync.WaitGroup
}

// Worker 工作器
type Worker struct {
    name string
}

func (w *Worker) Process(data map[string]interface{}) (interface{}, error) {
    // 模拟处理时间
    time.Sleep(100 * time.Millisecond)
    
    result := fmt.Sprintf("Worker %s processed: %v", w.name, data)
    return result, nil
}

// NewActiveObject 创建活动对象
func NewActiveObject(workerName string) *ActiveObject {
    ctx, cancel := context.WithCancel(context.Background())
    
    ao := &ActiveObject{
        requestQueue: make(chan *Request, 100),
        worker:       &Worker{name: workerName},
        ctx:          ctx,
        cancel:       cancel,
    }
    
    ao.wg.Add(1)
    go ao.scheduler()
    
    return ao
}

// scheduler 调度器
func (ao *ActiveObject) scheduler() {
    defer ao.wg.Done()
    
    for {
        select {
        case request := <-ao.requestQueue:
            ao.wg.Add(1)
            go ao.executeRequest(request)
        case <-ao.ctx.Done():
            return
        }
    }
}

// executeRequest 执行请求
func (ao *ActiveObject) executeRequest(request *Request) {
    defer ao.wg.Done()
    
    result, err := ao.worker.Process(request.Params)
    
    if err != nil {
        request.Error <- err
    } else {
        request.Response <- result
    }
}

// ProcessAsync 异步处理
func (ao *ActiveObject) ProcessAsync(params map[string]interface{}) *Request {
    request := &Request{
        ID:       generateID(),
        Method:   "Process",
        Params:   params,
        Response: make(chan interface{}, 1),
        Error:    make(chan error, 1),
    }
    
    ao.requestQueue <- request
    return request
}

// ProcessSync 同步处理
func (ao *ActiveObject) ProcessSync(params map[string]interface{}) (interface{}, error) {
    request := ao.ProcessAsync(params)
    
    select {
    case result := <-request.Response:
        return result, nil
    case err := <-request.Error:
        return nil, err
    case <-time.After(5 * time.Second):
        return nil, fmt.Errorf("timeout")
    }
}

// Shutdown 关闭活动对象
func (ao *ActiveObject) Shutdown() {
    ao.cancel()
    ao.wg.Wait()
    close(ao.requestQueue)
}

func generateID() string {
    return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

```

### 14.1.2.4 性能分析

- **响应时间**: $O(1)$ - 立即返回
- **处理时间**: $O(n)$ - n为请求数量
- **内存使用**: $O(n)$ - 队列大小

## 14.1.3 2. 管程模式 (Monitor)

### 14.1.3.1 定义

提供线程安全的资源访问机制，确保同一时间只有一个线程能访问共享资源。

### 14.1.3.2 Golang实现

```go
package monitor

import (
    "fmt"
    "sync"
    "time"
)

// Monitor 管程
type Monitor struct {
    mu      sync.Mutex
    cond    *sync.Cond
    data    map[string]interface{}
    version int
}

// NewMonitor 创建管程
func NewMonitor() *Monitor {
    m := &Monitor{
        data: make(map[string]interface{}),
    }
    m.cond = sync.NewCond(&m.mu)
    return m
}

// Set 设置数据
func (m *Monitor) Set(key string, value interface{}) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.data[key] = value
    m.version++
    m.cond.Broadcast() // 通知所有等待的goroutine
}

// Get 获取数据
func (m *Monitor) Get(key string) (interface{}, bool) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    value, exists := m.data[key]
    return value, exists
}

// WaitForCondition 等待条件
func (m *Monitor) WaitForCondition(condition func() bool) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    for !condition() {
        m.cond.Wait()
    }
}

// WaitForValue 等待特定值
func (m *Monitor) WaitForValue(key string, expectedValue interface{}) {
    m.WaitForCondition(func() bool {
        value, exists := m.data[key]
        return exists && value == expectedValue
    })
}

// GetVersion 获取版本号
func (m *Monitor) GetVersion() int {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.version
}

// Snapshot 获取快照
func (m *Monitor) Snapshot() map[string]interface{} {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    snapshot := make(map[string]interface{})
    for k, v := range m.data {
        snapshot[k] = v
    }
    return snapshot
}

```

## 14.1.4 3. 线程池模式 (Thread Pool)

### 14.1.4.1 定义

预先创建一组线程，用于执行任务，避免频繁创建和销毁线程的开销。

### 14.1.4.2 Golang实现

```go
package threadpool

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Task 任务接口
type Task interface {
    Execute() (interface{}, error)
    GetID() string
}

// SimpleTask 简单任务
type SimpleTask struct {
    ID       string
    Function func() (interface{}, error)
}

func (t *SimpleTask) Execute() (interface{}, error) {
    return t.Function()
}

func (t *SimpleTask) GetID() string {
    return t.ID
}

// Worker 工作器
type Worker struct {
    ID       int
    taskChan <-chan Task
    resultChan chan<- TaskResult
    ctx       context.Context
    wg        *sync.WaitGroup
}

// TaskResult 任务结果
type TaskResult struct {
    TaskID string
    Result interface{}
    Error  error
    WorkerID int
}

func (w *Worker) Start() {
    defer w.wg.Done()
    
    for {
        select {
        case task := <-w.taskChan:
            if task == nil {
                return
            }
            
            result, err := task.Execute()
            w.resultChan <- TaskResult{
                TaskID:   task.GetID(),
                Result:   result,
                Error:    err,
                WorkerID: w.ID,
            }
        case <-w.ctx.Done():
            return
        }
    }
}

// ThreadPool 线程池
type ThreadPool struct {
    workers    []*Worker
    taskChan   chan Task
    resultChan chan TaskResult
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
    mu         sync.RWMutex
    stats      *PoolStats
}

// PoolStats 池统计
type PoolStats struct {
    TotalTasks    int64
    CompletedTasks int64
    FailedTasks   int64
    ActiveWorkers int
    mu            sync.RWMutex
}

func (s *PoolStats) IncrementTotal() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.TotalTasks++
}

func (s *PoolStats) IncrementCompleted() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.CompletedTasks++
}

func (s *PoolStats) IncrementFailed() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.FailedTasks++
}

func (s *PoolStats) SetActiveWorkers(count int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.ActiveWorkers = count
}

func (s *PoolStats) GetStats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return map[string]interface{}{
        "total_tasks":     s.TotalTasks,
        "completed_tasks": s.CompletedTasks,
        "failed_tasks":    s.FailedTasks,
        "active_workers":  s.ActiveWorkers,
        "success_rate":    float64(s.CompletedTasks) / float64(s.TotalTasks),
    }
}

// NewThreadPool 创建线程池
func NewThreadPool(workerCount int) *ThreadPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    pool := &ThreadPool{
        taskChan:   make(chan Task, workerCount*2),
        resultChan: make(chan TaskResult, workerCount*2),
        ctx:        ctx,
        cancel:     cancel,
        stats:      &PoolStats{},
    }
    
    // 创建工作器
    for i := 0; i < workerCount; i++ {
        worker := &Worker{
            ID:        i,
            taskChan:  pool.taskChan,
            resultChan: pool.resultChan,
            ctx:       ctx,
            wg:        &pool.wg,
        }
        pool.workers = append(pool.workers, worker)
    }
    
    // 启动工作器
    pool.wg.Add(workerCount)
    for _, worker := range pool.workers {
        go worker.Start()
    }
    
    // 启动结果处理器
    go pool.resultProcessor()
    
    pool.stats.SetActiveWorkers(workerCount)
    
    return pool
}

// Submit 提交任务
func (p *ThreadPool) Submit(task Task) error {
    select {
    case p.taskChan <- task:
        p.stats.IncrementTotal()
        return nil
    case <-p.ctx.Done():
        return fmt.Errorf("thread pool is closed")
    default:
        return fmt.Errorf("task queue is full")
    }
}

// SubmitFunc 提交函数任务
func (p *ThreadPool) SubmitFunc(id string, fn func() (interface{}, error)) error {
    task := &SimpleTask{
        ID:       id,
        Function: fn,
    }
    return p.Submit(task)
}

// resultProcessor 结果处理器
func (p *ThreadPool) resultProcessor() {
    for result := range p.resultChan {
        if result.Error != nil {
            p.stats.IncrementFailed()
        } else {
            p.stats.IncrementCompleted()
        }
    }
}

// Shutdown 关闭线程池
func (p *ThreadPool) Shutdown() {
    p.cancel()
    close(p.taskChan)
    p.wg.Wait()
    close(p.resultChan)
}

// GetStats 获取统计信息
func (p *ThreadPool) GetStats() map[string]interface{} {
    return p.stats.GetStats()
}

```

## 14.1.5 4. 生产者-消费者模式 (Producer-Consumer)

### 14.1.5.1 定义

通过队列解耦生产者和消费者，实现异步数据处理。

### 14.1.5.2 Golang实现

```go
package producerconsumer

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Item 数据项
type Item struct {
    ID      string
    Data    interface{}
    Created time.Time
}

// Producer 生产者
type Producer struct {
    ID       string
    itemChan chan<- *Item
    ctx      context.Context
    wg       *sync.WaitGroup
    rate     time.Duration
}

func (p *Producer) Start() {
    defer p.wg.Done()
    
    ticker := time.NewTicker(p.rate)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            item := &Item{
                ID:      fmt.Sprintf("%s_%d", p.ID, time.Now().UnixNano()),
                Data:    fmt.Sprintf("Data from producer %s", p.ID),
                Created: time.Now(),
            }
            
            select {
            case p.itemChan <- item:
                fmt.Printf("Producer %s produced item: %s\n", p.ID, item.ID)
            case <-p.ctx.Done():
                return
            }
        case <-p.ctx.Done():
            return
        }
    }
}

// Consumer 消费者
type Consumer struct {
    ID       string
    itemChan <-chan *Item
    ctx      context.Context
    wg       *sync.WaitGroup
    processor func(*Item) error
}

func (c *Consumer) Start() {
    defer c.wg.Done()
    
    for {
        select {
        case item := <-c.itemChan:
            if item == nil {
                return
            }
            
            if c.processor != nil {
                if err := c.processor(item); err != nil {
                    fmt.Printf("Consumer %s failed to process item %s: %v\n", 
                        c.ID, item.ID, err)
                } else {
                    fmt.Printf("Consumer %s processed item: %s\n", c.ID, item.ID)
                }
            } else {
                // 默认处理
                time.Sleep(50 * time.Millisecond) // 模拟处理时间
                fmt.Printf("Consumer %s processed item: %s\n", c.ID, item.ID)
            }
        case <-c.ctx.Done():
            return
        }
    }
}

// Buffer 缓冲区
type Buffer struct {
    items    chan *Item
    capacity int
    mu       sync.RWMutex
    stats    *BufferStats
}

// BufferStats 缓冲区统计
type BufferStats struct {
    Produced int64
    Consumed int64
    Dropped  int64
    mu       sync.RWMutex
}

func (s *BufferStats) IncrementProduced() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Produced++
}

func (s *BufferStats) IncrementConsumed() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Consumed++
}

func (s *BufferStats) IncrementDropped() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Dropped++
}

func (s *BufferStats) GetStats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return map[string]interface{}{
        "produced": s.Produced,
        "consumed": s.Consumed,
        "dropped":  s.Dropped,
        "utilization": float64(s.Consumed) / float64(s.Produced),
    }
}

// NewBuffer 创建缓冲区
func NewBuffer(capacity int) *Buffer {
    return &Buffer{
        items:    make(chan *Item, capacity),
        capacity: capacity,
        stats:    &BufferStats{},
    }
}

// Put 放入数据
func (b *Buffer) Put(item *Item) bool {
    select {
    case b.items <- item:
        b.stats.IncrementProduced()
        return true
    default:
        b.stats.IncrementDropped()
        return false
    }
}

// Get 获取数据
func (b *Buffer) Get() *Item {
    select {
    case item := <-b.items:
        b.stats.IncrementConsumed()
        return item
    default:
        return nil
    }
}

// ProducerConsumerSystem 生产者消费者系统
type ProducerConsumerSystem struct {
    buffer    *Buffer
    producers []*Producer
    consumers []*Consumer
    ctx       context.Context
    cancel    context.CancelFunc
    wg        sync.WaitGroup
}

// NewProducerConsumerSystem 创建生产者消费者系统
func NewProducerConsumerSystem(bufferSize, producerCount, consumerCount int) *ProducerConsumerSystem {
    ctx, cancel := context.WithCancel(context.Background())
    
    system := &ProducerConsumerSystem{
        buffer: NewBuffer(bufferSize),
        ctx:    ctx,
        cancel: cancel,
    }
    
    // 创建生产者
    for i := 0; i < producerCount; i++ {
        producer := &Producer{
            ID:       fmt.Sprintf("P%d", i),
            itemChan: system.buffer.items,
            ctx:      ctx,
            wg:       &system.wg,
            rate:     time.Duration(100+i*50) * time.Millisecond,
        }
        system.producers = append(system.producers, producer)
    }
    
    // 创建消费者
    for i := 0; i < consumerCount; i++ {
        consumer := &Consumer{
            ID:       fmt.Sprintf("C%d", i),
            itemChan: system.buffer.items,
            ctx:      ctx,
            wg:       &system.wg,
        }
        system.consumers = append(system.consumers, consumer)
    }
    
    return system
}

// Start 启动系统
func (s *ProducerConsumerSystem) Start() {
    // 启动生产者
    s.wg.Add(len(s.producers))
    for _, producer := range s.producers {
        go producer.Start()
    }
    
    // 启动消费者
    s.wg.Add(len(s.consumers))
    for _, consumer := range s.consumers {
        go consumer.Start()
    }
}

// Stop 停止系统
func (s *ProducerConsumerSystem) Stop() {
    s.cancel()
    s.wg.Wait()
    close(s.buffer.items)
}

// GetStats 获取统计信息
func (s *ProducerConsumerSystem) GetStats() map[string]interface{} {
    return s.buffer.stats.GetStats()
}

```

## 14.1.6 5. 读写锁模式 (Readers-Writer Lock)

### 14.1.6.1 定义

允许多个读者同时访问资源，但只允许一个写者访问。

### 14.1.6.2 Golang实现

```go
package rwlock

import (
    "fmt"
    "sync"
    "time"
)

// RWLock 读写锁
type RWLock struct {
    mu       sync.RWMutex
    data     map[string]interface{}
    version  int
    stats    *RWLockStats
}

// RWLockStats 读写锁统计
type RWLockStats struct {
    ReadCount   int64
    WriteCount  int64
    ReadWait    time.Duration
    WriteWait   time.Duration
    mu          sync.RWMutex
}

func (s *RWLockStats) IncrementRead() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.ReadCount++
}

func (s *RWLockStats) IncrementWrite() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.WriteCount++
}

func (s *RWLockStats) AddReadWait(duration time.Duration) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.ReadWait += duration
}

func (s *RWLockStats) AddWriteWait(duration time.Duration) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.WriteWait += duration
}

func (s *RWLockStats) GetStats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    avgReadWait := time.Duration(0)
    avgWriteWait := time.Duration(0)
    
    if s.ReadCount > 0 {
        avgReadWait = s.ReadWait / time.Duration(s.ReadCount)
    }
    if s.WriteCount > 0 {
        avgWriteWait = s.WriteWait / time.Duration(s.WriteCount)
    }
    
    return map[string]interface{}{
        "read_count":    s.ReadCount,
        "write_count":   s.WriteCount,
        "avg_read_wait": avgReadWait,
        "avg_write_wait": avgWriteWait,
        "read_write_ratio": float64(s.ReadCount) / float64(s.WriteCount),
    }
}

// NewRWLock 创建读写锁
func NewRWLock() *RWLock {
    return &RWLock{
        data:  make(map[string]interface{}),
        stats: &RWLockStats{},
    }
}

// Read 读取操作
func (r *RWLock) Read(key string) (interface{}, bool) {
    start := time.Now()
    
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    r.stats.AddReadWait(time.Since(start))
    r.stats.IncrementRead()
    
    value, exists := r.data[key]
    return value, exists
}

// Write 写入操作
func (r *RWLock) Write(key string, value interface{}) {
    start := time.Now()
    
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.stats.AddWriteWait(time.Since(start))
    r.stats.IncrementWrite()
    
    r.data[key] = value
    r.version++
}

// Delete 删除操作
func (r *RWLock) Delete(key string) bool {
    start := time.Now()
    
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.stats.AddWriteWait(time.Since(start))
    r.stats.IncrementWrite()
    
    if _, exists := r.data[key]; exists {
        delete(r.data, key)
        r.version++
        return true
    }
    return false
}

// GetVersion 获取版本号
func (r *RWLock) GetVersion() int {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.version
}

// Snapshot 获取快照
func (r *RWLock) Snapshot() map[string]interface{} {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    snapshot := make(map[string]interface{})
    for k, v := range r.data {
        snapshot[k] = v
    }
    return snapshot
}

// GetStats 获取统计信息
func (r *RWLock) GetStats() map[string]interface{} {
    return r.stats.GetStats()
}

```

## 14.1.7 6. Future/Promise模式

### 14.1.7.1 定义

表示异步计算的结果，可以在计算完成后获取结果。

### 14.1.7.2 Golang实现

```go
package future

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Future 未来对象
type Future struct {
    result    interface{}
    error     error
    done      chan struct{}
    mu        sync.RWMutex
    completed bool
}

// NewFuture 创建未来对象
func NewFuture() *Future {
    return &Future{
        done: make(chan struct{}),
    }
}

// SetResult 设置结果
func (f *Future) SetResult(result interface{}) {
    f.mu.Lock()
    defer f.mu.Unlock()
    
    if !f.completed {
        f.result = result
        f.completed = true
        close(f.done)
    }
}

// SetError 设置错误
func (f *Future) SetError(err error) {
    f.mu.Lock()
    defer f.mu.Unlock()
    
    if !f.completed {
        f.error = err
        f.completed = true
        close(f.done)
    }
}

// Get 获取结果
func (f *Future) Get() (interface{}, error) {
    <-f.done
    return f.result, f.error
}

// GetWithTimeout 带超时的获取结果
func (f *Future) GetWithTimeout(timeout time.Duration) (interface{}, error) {
    select {
    case <-f.done:
        return f.result, f.error
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout")
    }
}

// GetWithContext 带上下文的获取结果
func (f *Future) GetWithContext(ctx context.Context) (interface{}, error) {
    select {
    case <-f.done:
        return f.result, f.error
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// IsDone 是否完成
func (f *Future) IsDone() bool {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return f.completed
}

// Then 链式调用
func (f *Future) Then(fn func(interface{}) (interface{}, error)) *Future {
    newFuture := NewFuture()
    
    go func() {
        result, err := f.Get()
        if err != nil {
            newFuture.SetError(err)
            return
        }
        
        newResult, newErr := fn(result)
        if newErr != nil {
            newFuture.SetError(newErr)
        } else {
            newFuture.SetResult(newResult)
        }
    }()
    
    return newFuture
}

// Promise 承诺对象
type Promise struct {
    future *Future
}

// NewPromise 创建承诺对象
func NewPromise() *Promise {
    return &Promise{
        future: NewFuture(),
    }
}

// Resolve 解决承诺
func (p *Promise) Resolve(result interface{}) {
    p.future.SetResult(result)
}

// Reject 拒绝承诺
func (p *Promise) Reject(err error) {
    p.future.SetError(err)
}

// Future 获取未来对象
func (p *Promise) Future() *Future {
    return p.future
}

// FutureExecutor 未来执行器
type FutureExecutor struct {
    workers chan struct{}
}

// NewFutureExecutor 创建未来执行器
func NewFutureExecutor(maxWorkers int) *FutureExecutor {
    return &FutureExecutor{
        workers: make(chan struct{}, maxWorkers),
    }
}

// Submit 提交任务
func (p *FutureExecutor) Submit(task func() (interface{}, error)) *Future {
    future := NewFuture()
    
    go func() {
        p.workers <- struct{}{} // 获取工作槽
        defer func() {
            <-p.workers // 释放工作槽
        }()
        
        result, err := task()
        if err != nil {
            future.SetError(err)
        } else {
            future.SetResult(result)
        }
    }()
    
    return future
}

// All 等待所有未来完成
func (p *FutureExecutor) All(futures []*Future) *Future {
    resultFuture := NewFuture()
    
    go func() {
        results := make([]interface{}, len(futures))
        for i, future := range futures {
            result, err := future.Get()
            if err != nil {
                resultFuture.SetError(err)
                return
            }
            results[i] = result
        }
        resultFuture.SetResult(results)
    }()
    
    return resultFuture
}

// Any 等待任一未来完成
func (p *FutureExecutor) Any(futures []*Future) *Future {
    resultFuture := NewFuture()
    
    go func() {
        for _, future := range futures {
            result, err := future.Get()
            if err == nil {
                resultFuture.SetResult(result)
                return
            }
        }
        resultFuture.SetError(fmt.Errorf("all futures failed"))
    }()
    
    return resultFuture
}

```

## 14.1.8 7. Actor模型

### 14.1.8.1 定义

基于消息传递的并发模型，每个Actor都是独立的计算单元。

### 14.1.8.2 Golang实现

```go
package actor

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Message 消息接口
type Message interface {
    GetType() string
    GetData() interface{}
}

// SimpleMessage 简单消息
type SimpleMessage struct {
    Type string
    Data interface{}
}

func (m *SimpleMessage) GetType() string {
    return m.Type
}

func (m *SimpleMessage) GetData() interface{} {
    return m.Data
}

// Actor 演员接口
type Actor interface {
    Receive(message Message)
    GetID() string
}

// BaseActor 基础演员
type BaseActor struct {
    ID       string
    mailbox  chan Message
    ctx      context.Context
    cancel   context.CancelFunc
    wg       sync.WaitGroup
    behavior func(Message)
}

// NewBaseActor 创建基础演员
func NewBaseActor(id string, behavior func(Message)) *BaseActor {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &BaseActor{
        ID:       id,
        mailbox:  make(chan Message, 100),
        ctx:      ctx,
        cancel:   cancel,
        behavior: behavior,
    }
}

// Start 启动演员
func (a *BaseActor) Start() {
    a.wg.Add(1)
    go a.run()
}

// run 运行演员
func (a *BaseActor) run() {
    defer a.wg.Done()
    
    for {
        select {
        case message := <-a.mailbox:
            if message == nil {
                return
            }
            a.behavior(message)
        case <-a.ctx.Done():
            return
        }
    }
}

// Send 发送消息
func (a *BaseActor) Send(message Message) {
    select {
    case a.mailbox <- message:
    case <-a.ctx.Done():
    }
}

// Stop 停止演员
func (a *BaseActor) Stop() {
    a.cancel()
    close(a.mailbox)
    a.wg.Wait()
}

// ActorSystem 演员系统
type ActorSystem struct {
    actors map[string]Actor
    mu     sync.RWMutex
    ctx    context.Context
    cancel context.CancelFunc
}

// NewActorSystem 创建演员系统
func NewActorSystem() *ActorSystem {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &ActorSystem{
        actors: make(map[string]Actor),
        ctx:    ctx,
        cancel: cancel,
    }
}

// RegisterActor 注册演员
func (s *ActorSystem) RegisterActor(actor Actor) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.actors[actor.GetID()] = actor
    
    if baseActor, ok := actor.(*BaseActor); ok {
        baseActor.Start()
    }
}

// SendMessage 发送消息
func (s *ActorSystem) SendMessage(actorID string, message Message) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    actor, exists := s.actors[actorID]
    if !exists {
        return fmt.Errorf("actor not found: %s", actorID)
    }
    
    actor.Send(message)
    return nil
}

// GetActor 获取演员
func (s *ActorSystem) GetActor(actorID string) (Actor, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    actor, exists := s.actors[actorID]
    return actor, exists
}

// Shutdown 关闭系统
func (s *ActorSystem) Shutdown() {
    s.cancel()
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for _, actor := range s.actors {
        if baseActor, ok := actor.(*BaseActor); ok {
            baseActor.Stop()
        }
    }
}

// ExampleActor 示例演员
type ExampleActor struct {
    *BaseActor
    data map[string]interface{}
    mu   sync.RWMutex
}

// NewExampleActor 创建示例演员
func NewExampleActor(id string) *ExampleActor {
    actor := &ExampleActor{
        data: make(map[string]interface{}),
    }
    
    actor.BaseActor = NewBaseActor(id, actor.handleMessage)
    return actor
}

// handleMessage 处理消息
func (a *ExampleActor) handleMessage(message Message) {
    switch message.GetType() {
    case "SET":
        if data, ok := message.GetData().(map[string]interface{}); ok {
            a.mu.Lock()
            for k, v := range data {
                a.data[k] = v
            }
            a.mu.Unlock()
            fmt.Printf("Actor %s: SET data %v\n", a.ID, data)
        }
    case "GET":
        if key, ok := message.GetData().(string); ok {
            a.mu.RLock()
            value, exists := a.data[key]
            a.mu.RUnlock()
            
            if exists {
                fmt.Printf("Actor %s: GET %s = %v\n", a.ID, key, value)
            } else {
                fmt.Printf("Actor %s: GET %s not found\n", a.ID, key)
            }
        }
    case "DELETE":
        if key, ok := message.GetData().(string); ok {
            a.mu.Lock()
            delete(a.data, key)
            a.mu.Unlock()
            fmt.Printf("Actor %s: DELETE %s\n", a.ID, key)
        }
    default:
        fmt.Printf("Actor %s: Unknown message type %s\n", a.ID, message.GetType())
    }
}

```

## 14.1.9 8. 性能分析

### 14.1.9.1 模式性能对比

| 模式 | 创建开销 | 内存使用 | 并发度 | 延迟 | 吞吐量 |
|------|---------|---------|--------|------|--------|
| 活动对象 | 低 | 中 | 高 | 低 | 高 |
| 管程 | 低 | 低 | 中 | 中 | 中 |
| 线程池 | 中 | 中 | 高 | 低 | 高 |
| 生产者-消费者 | 低 | 中 | 高 | 低 | 高 |
| 读写锁 | 低 | 低 | 高 | 中 | 高 |
| Future/Promise | 低 | 低 | 高 | 低 | 高 |
| Actor | 中 | 中 | 高 | 中 | 高 |

### 14.1.9.2 内存使用分析

**活动对象模式**:
$$M_{active} = \sum_{i=1}^{n} \text{sizeof}(request_i) + \text{sizeof}(worker)$$

**线程池模式**:
$$M_{pool} = \sum_{i=1}^{w} \text{sizeof}(worker_i) + \text{sizeof}(taskQueue)$$

**Actor模式**:
$$M_{actor} = \sum_{i=1}^{a} \text{sizeof}(actor_i) + \text{sizeof}(mailbox_i)$$

### 14.1.9.3 并发性能分析

**吞吐量计算**:
$$\text{Throughput} = \frac{\text{ConcurrentWorkers} \times \text{ProcessingRate}}{\text{Latency}}$$

**延迟分析**:
$$\text{Latency} = \text{QueueTime} + \text{ProcessingTime} + \text{NetworkTime}$$

## 14.1.10 9. 最佳实践

### 14.1.10.1 模式选择原则

1. **简单优先**: 优先使用简单的并发模式
2. **性能考虑**: 根据性能需求选择合适的模式
3. **可维护性**: 考虑代码的可维护性
4. **扩展性**: 考虑系统的扩展性

### 14.1.10.2 实现建议

1. **错误处理**: 完善的错误处理机制
2. **资源管理**: 及时释放资源
3. **监控告警**: 实时监控系统状态
4. **测试策略**: 全面的并发测试

### 14.1.10.3 常见陷阱

1. **死锁**: 避免循环等待
2. **竞态条件**: 使用适当的同步机制
3. **内存泄漏**: 注意资源管理
4. **性能瓶颈**: 避免热点资源

## 14.1.11 10. 应用场景

### 14.1.11.1 活动对象模式

- 异步API调用
- 事件处理
- 消息队列
- 任务调度

### 14.1.11.2 管程模式

- 共享资源管理
- 数据库连接池
- 缓存管理
- 配置管理

### 14.1.11.3 线程池模式

- Web服务器
- 数据处理
- 图像处理
- 批量任务

### 14.1.11.4 生产者-消费者模式

- 日志处理
- 数据流处理
- 事件流处理
- 消息队列

### 14.1.11.5 读写锁模式

- 配置管理
- 缓存系统
- 数据库访问
- 文件系统

### 14.1.11.6 Future/Promise模式

- 异步API
- 并行计算
- 网络请求
- 数据处理

### 14.1.11.7 Actor模式

- 分布式系统
- 游戏服务器
- 实时系统
- 事件驱动架构

## 14.1.12 11. 总结

并发并行模式为高并发系统提供了重要的设计指导。通过合理应用这些模式，可以构建出高性能、高可用的并发系统。

### 14.1.12.1 关键优势

- **高性能**: 充分利用多核资源
- **高可用**: 提高系统可靠性
- **可扩展**: 支持水平扩展
- **可维护**: 清晰的代码结构

### 14.1.12.2 成功要素

1. **合理选择**: 根据需求选择合适的模式
2. **性能优化**: 持续的性能优化
3. **监控告警**: 完善的监控体系
4. **测试验证**: 全面的测试覆盖

通过合理应用并发并行模式，可以构建出高质量的并发系统，为业务发展提供强有力的技术支撑。
