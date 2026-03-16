# 并发工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [并发工具](#并发工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

并发工具提供了丰富的并发控制原语，包括goroutine池、worker池、信号量、互斥锁、屏障等，简化并发编程任务。

---

## 2. 功能特性

### 2.1 Goroutine池

- `Pool`: goroutine池，用于管理固定数量的goroutine
- `NewPool`: 创建新的goroutine池
- `Start`: 启动goroutine池
- `Submit`: 提交任务
- `SubmitWithContext`: 使用context提交任务
- `Stop`: 停止goroutine池
- `Wait`: 等待所有任务完成

### 2.2 Worker池

- `WorkerPool`: worker池，用于处理任务并返回结果
- `NewWorkerPool`: 创建新的worker池
- `Start`: 启动worker池
- `Submit`: 提交任务
- `GetResult`: 获取结果
- `Stop`: 停止worker池
- `Wait`: 等待所有任务完成

### 2.3 信号量

- `Semaphore`: 信号量，用于控制并发数量
- `NewSemaphore`: 创建新的信号量
- `Acquire`: 获取信号量
- `Release`: 释放信号量
- `TryAcquire`: 尝试获取信号量（非阻塞）
- `AcquireWithContext`: 使用context获取信号量
- `Size`: 获取当前信号量大小
- `Capacity`: 获取信号量容量

### 2.4 互斥锁

- `Mutex`: 互斥锁（带超时）
- `NewMutex`: 创建新的互斥锁
- `Lock`: 加锁
- `Unlock`: 解锁
- `TryLock`: 尝试加锁（非阻塞）
- `LockWithContext`: 使用context加锁

### 2.5 Once

- `Once`: 只执行一次（带错误处理）
- `Do`: 执行函数（只执行一次）
- `Reset`: 重置Once

### 2.6 屏障

- `Barrier`: 屏障，用于同步多个goroutine
- `NewBarrier`: 创建新的屏障
- `Wait`: 等待屏障
- `Reset`: 重置屏障

### 2.7 等待组

- `WaitGroup`: 等待组（带超时）
- `NewWaitGroup`: 创建新的等待组
- `Add`: 添加计数
- `Done`: 完成计数
- `Wait`: 等待完成（支持超时）

---

## 3. 使用示例

### 3.1 Goroutine池

```go
import "github.com/yourusername/golang/pkg/utils/concurrency"

// 创建goroutine池
pool := concurrency.NewPool(3, 10)
pool.Start()
defer pool.Stop()

// 提交任务
for i := 0; i < 10; i++ {
    j := i
    pool.Submit(func() {
        // 处理任务
        fmt.Println(j)
    })
}

// 使用context提交任务
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
err := pool.SubmitWithContext(ctx, func() {
    // 处理任务
})
```

### 3.2 Worker池

```go
// 创建worker池
processor := func(job interface{}) interface{} {
    return job.(int) * 2
}
pool := concurrency.NewWorkerPool(3, 10, processor)
pool.Start()
defer pool.Stop()

// 提交任务
for i := 0; i < 10; i++ {
    pool.Submit(i)
}

// 获取结果
for i := 0; i < 10; i++ {
    result, err := pool.GetResult()
    if err != nil {
        // 处理错误
    }
    fmt.Println(result)
}
```

### 3.3 信号量

```go
// 创建信号量
sem := concurrency.NewSemaphore(3)

// 获取信号量
sem.Acquire()
defer sem.Release()

// 尝试获取信号量（非阻塞）
if sem.TryAcquire() {
    defer sem.Release()
    // 处理任务
}

// 使用context获取信号量
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
err := sem.AcquireWithContext(ctx)
if err != nil {
    // 处理错误
}
defer sem.Release()
```

### 3.4 互斥锁

```go
// 创建互斥锁
mutex := concurrency.NewMutex()

// 加锁
mutex.Lock()
defer mutex.Unlock()

// 尝试加锁（非阻塞）
if mutex.TryLock() {
    defer mutex.Unlock()
    // 处理任务
}

// 使用context加锁
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
err := mutex.LockWithContext(ctx)
if err != nil {
    // 处理错误
}
defer mutex.Unlock()
```

### 3.5 Once

```go
// 创建Once
once := &concurrency.Once{}

// 执行函数（只执行一次）
result, err := once.Do(func() (interface{}, error) {
    // 初始化操作
    return "result", nil
})

// 重置Once
once.Reset()
```

### 3.6 屏障

```go
// 创建屏障
barrier := concurrency.NewBarrier(3)

// 在多个goroutine中等待屏障
for i := 0; i < 3; i++ {
    go func() {
        // 准备工作
        barrier.Wait()
        // 所有goroutine都到达屏障后继续执行
    }()
}
```

### 3.7 等待组

```go
// 创建等待组（带超时）
wg := concurrency.NewWaitGroup(time.Second)

// 添加计数
wg.Add(3)

// 启动goroutine
for i := 0; i < 3; i++ {
    go func() {
        defer wg.Done()
        // 处理任务
    }()
}

// 等待完成（支持超时）
err := wg.Wait()
if err != nil {
    // 处理超时错误
}
```

---

**更新日期**: 2025-11-11
