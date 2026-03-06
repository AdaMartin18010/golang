# Go 1.23 并发编程模式全面指南

> 本文档深入剖析 **Go 1.23** 语言并发编程的核心机制与设计模式，涵盖从基础概念到高级模式的完整知识体系。
>
> **Go 1.23 更新**：
>
> - Timer/Ticker 改进（无缓冲channel、立即GC）
> - 新增 `sync.Map.Clear` 方法
> - 新增 `sync/atomic` 的 `And` / `Or` 操作
> - 迭代器模式在并发场景的应用

---

## 目录

- [Go 1.23 并发编程模式全面指南](#go-123-并发编程模式全面指南)
  - [目录](#目录)
  - [1. Goroutine基础](#1-goroutine基础)
    - [1.1 Goroutine创建与生命周期](#11-goroutine创建与生命周期)
      - [概念定义](#概念定义)
      - [工作原理](#工作原理)
      - [形式论证](#形式论证)
      - [完整示例](#完整示例)
      - [反例说明](#反例说明)
      - [性能分析](#性能分析)
      - [最佳实践](#最佳实践)
    - [1.2 Goroutine与线程区别](#12-goroutine与线程区别)
      - [概念定义](#概念定义-1)
      - [工作原理](#工作原理-1)
      - [完整示例](#完整示例-1)
      - [性能分析](#性能分析-1)
      - [最佳实践](#最佳实践-1)
    - [1.3 Goroutine泄漏与检测](#13-goroutine泄漏与检测)
      - [概念定义](#概念定义-2)
      - [工作原理](#工作原理-2)
      - [完整示例](#完整示例-2)
      - [检测方法](#检测方法)
      - [最佳实践](#最佳实践-2)
    - [1.4 Goroutine数量控制](#14-goroutine数量控制)
      - [概念定义](#概念定义-3)
      - [工作原理](#工作原理-3)
      - [完整示例](#完整示例-3)
      - [性能分析](#性能分析-2)
      - [最佳实践](#最佳实践-3)
  - [2. Channel模式](#2-channel模式)
    - [2.1 无缓冲Channel](#21-无缓冲channel)
      - [概念定义](#概念定义-4)
      - [工作原理](#工作原理-4)
      - [形式论证](#形式论证-1)
      - [完整示例](#完整示例-4)
      - [反例说明](#反例说明-1)
      - [性能分析](#性能分析-3)
      - [最佳实践](#最佳实践-4)
    - [2.2 有缓冲Channel](#22-有缓冲channel)
      - [概念定义](#概念定义-5)
      - [工作原理](#工作原理-5)
      - [形式论证](#形式论证-2)
      - [完整示例](#完整示例-5)
      - [反例说明](#反例说明-2)
      - [性能分析](#性能分析-4)
      - [最佳实践](#最佳实践-5)
    - [2.3 单向Channel](#23-单向channel)
      - [概念定义](#概念定义-6)
      - [工作原理](#工作原理-6)
      - [形式论证](#形式论证-3)
      - [完整示例](#完整示例-6)
      - [反例说明](#反例说明-3)
      - [最佳实践](#最佳实践-6)
    - [2.4 Channel关闭与检测](#24-channel关闭与检测)
      - [概念定义](#概念定义-7)
      - [工作原理](#工作原理-7)
      - [形式论证](#形式论证-4)
      - [完整示例](#完整示例-7)
      - [反例说明](#反例说明-4)
      - [最佳实践](#最佳实践-7)
    - [2.5 Channel Select多路复用](#25-channel-select多路复用)
      - [概念定义](#概念定义-8)
      - [工作原理](#工作原理-8)
      - [形式论证](#形式论证-5)
      - [完整示例](#完整示例-8)
      - [反例说明](#反例说明-5)
      - [最佳实践](#最佳实践-8)
    - [2.6 超时模式](#26-超时模式)
      - [概念定义](#概念定义-9)
      - [工作原理](#工作原理-9)
      - [完整示例](#完整示例-9)
      - [性能分析](#性能分析-5)
      - [最佳实践](#最佳实践-9)
    - [2.7 心跳模式](#27-心跳模式)
      - [概念定义](#概念定义-10)
      - [工作原理](#工作原理-10)
      - [完整示例](#完整示例-10)
      - [最佳实践](#最佳实践-10)
    - [2.8 流水线模式](#28-流水线模式)
      - [概念定义](#概念定义-11)
      - [工作原理](#工作原理-11)
      - [完整示例](#完整示例-11)
      - [性能分析](#性能分析-6)
      - [最佳实践](#最佳实践-11)
  - [3. 同步原语（sync包）](#3-同步原语sync包)
    - [3.1 Mutex与RWMutex](#31-mutex与rwmutex)
      - [概念定义](#概念定义-12)
      - [工作原理](#工作原理-12)
      - [形式论证](#形式论证-6)
      - [完整示例](#完整示例-12)
      - [反例说明](#反例说明-6)
      - [性能分析](#性能分析-7)
      - [最佳实践](#最佳实践-12)
    - [3.2 WaitGroup](#32-waitgroup)
      - [概念定义](#概念定义-13)
      - [工作原理](#工作原理-13)
      - [完整示例](#完整示例-13)
      - [反例说明](#反例说明-7)
      - [最佳实践](#最佳实践-13)
    - [3.3 Once](#33-once)
      - [概念定义](#概念定义-14)
      - [工作原理](#工作原理-14)
      - [完整示例](#完整示例-14)
      - [反例说明](#反例说明-8)
      - [最佳实践](#最佳实践-14)
    - [3.4 Cond](#34-cond)
      - [概念定义](#概念定义-15)
      - [工作原理](#工作原理-15)
      - [完整示例](#完整示例-15)
      - [反例说明](#反例说明-9)
      - [最佳实践](#最佳实践-15)
    - [3.5 Pool](#35-pool)
      - [概念定义](#概念定义-16)
      - [工作原理](#工作原理-16)
      - [完整示例](#完整示例-16)
      - [反例说明](#反例说明-10)
      - [最佳实践](#最佳实践-16)
    - [3.6 Map](#36-map)
      - [概念定义](#概念定义-17)
      - [工作原理](#工作原理-17)
      - [完整示例](#完整示例-17)
      - [反例说明](#反例说明-11)
      - [最佳实践](#最佳实践-17)
    - [3.7 Atomic操作](#37-atomic操作)
      - [概念定义](#概念定义-18)
      - [工作原理](#工作原理-18)
      - [完整示例](#完整示例-18)
      - [反例说明](#反例说明-12)
      - [性能分析](#性能分析-8)
      - [最佳实践](#最佳实践-18)
  - [4. Context模式](#4-context模式)
    - [4.1 取消传播](#41-取消传播)
      - [概念定义](#概念定义-19)
      - [工作原理](#工作原理-19)
      - [完整示例](#完整示例-19)
      - [反例说明](#反例说明-13)
      - [最佳实践](#最佳实践-19)
    - [4.2 超时控制](#42-超时控制)
      - [概念定义](#概念定义-20)
      - [完整示例](#完整示例-20)
      - [最佳实践](#最佳实践-20)
    - [4.3 值传递](#43-值传递)
      - [概念定义](#概念定义-21)
      - [完整示例](#完整示例-21)
      - [反例说明](#反例说明-14)
      - [最佳实践](#最佳实践-21)
    - [4.4 Context链](#44-context链)
      - [概念定义](#概念定义-22)
      - [完整示例](#完整示例-22)
      - [最佳实践](#最佳实践-22)
  - [5. 常见并发模式](#5-常见并发模式)
    - [5.1 Worker Pool（工作池）](#51-worker-pool工作池)
      - [概念定义](#概念定义-23)
      - [完整示例](#完整示例-23)
      - [最佳实践](#最佳实践-23)
    - [5.2 Fan-Out/Fan-In](#52-fan-outfan-in)
      - [概念定义](#概念定义-24)
      - [完整示例](#完整示例-24)
      - [最佳实践](#最佳实践-24)
    - [5.3 Pipeline（管道）](#53-pipeline管道)
      - [概念定义](#概念定义-25)
      - [完整示例](#完整示例-25)
      - [最佳实践](#最佳实践-25)
    - [5.4 Tee（分流）](#54-tee分流)
      - [概念定义](#概念定义-26)
      - [完整示例](#完整示例-26)
      - [最佳实践](#最佳实践-26)
    - [5.5 Bridge（桥接）](#55-bridge桥接)
      - [概念定义](#概念定义-27)
      - [完整示例](#完整示例-27)
      - [最佳实践](#最佳实践-27)
    - [5.6 Or-Done（或完成）](#56-or-done或完成)
      - [概念定义](#概念定义-28)
      - [完整示例](#完整示例-28)
      - [最佳实践](#最佳实践-28)
    - [5.7 Or-Channel（或通道）](#57-or-channel或通道)
      - [概念定义](#概念定义-29)
      - [完整示例](#完整示例-29)
      - [最佳实践](#最佳实践-29)
    - [5.8 Quit信号模式](#58-quit信号模式)
      - [概念定义](#概念定义-30)
      - [完整示例](#完整示例-30)
      - [最佳实践](#最佳实践-30)
    - [5.9 速率限制（Rate Limiting）](#59-速率限制rate-limiting)
      - [概念定义](#概念定义-31)
      - [完整示例](#完整示例-31)
      - [最佳实践](#最佳实践-31)
    - [5.10 防抖与节流](#510-防抖与节流)
      - [概念定义](#概念定义-32)
      - [完整示例](#完整示例-32)
      - [最佳实践](#最佳实践-32)
  - [6. 并行计算模式](#6-并行计算模式)
    - [6.1 并行Map](#61-并行map)
      - [概念定义](#概念定义-33)
      - [完整示例](#完整示例-33)
      - [最佳实践](#最佳实践-33)
    - [6.2 并行Reduce](#62-并行reduce)
      - [概念定义](#概念定义-34)
      - [完整示例](#完整示例-34)
      - [最佳实践](#最佳实践-34)
    - [6.3 并行For](#63-并行for)
      - [概念定义](#概念定义-35)
      - [完整示例](#完整示例-35)
      - [最佳实践](#最佳实践-35)
    - [6.4 分治并行](#64-分治并行)
      - [概念定义](#概念定义-36)
      - [完整示例](#完整示例-36)
      - [最佳实践](#最佳实践-36)
  - [7. 异步编程模式](#7-异步编程模式)
    - [7.1 Future/Promise模式](#71-futurepromise模式)
      - [概念定义](#概念定义-37)
      - [完整示例](#完整示例-37)
      - [最佳实践](#最佳实践-37)
    - [7.2 Callback模式](#72-callback模式)
      - [概念定义](#概念定义-38)
      - [完整示例](#完整示例-38)
      - [最佳实践](#最佳实践-38)
    - [7.3 Async/Await模拟](#73-asyncawait模拟)
      - [概念定义](#概念定义-39)
      - [完整示例](#完整示例-39)
      - [最佳实践](#最佳实践-39)
    - [7.4 事件驱动模式](#74-事件驱动模式)
      - [概念定义](#概念定义-40)
      - [完整示例](#完整示例-40)
      - [最佳实践](#最佳实践-40)
  - [8. 并发安全](#8-并发安全)
    - [8.1 数据竞争检测](#81-数据竞争检测)
      - [概念定义](#概念定义-41)
      - [工作原理](#工作原理-20)
      - [完整示例](#完整示例-41)
      - [反例说明](#反例说明-15)
      - [最佳实践](#最佳实践-41)
    - [8.2 Happens-Before关系](#82-happens-before关系)
      - [概念定义](#概念定义-42)
      - [工作原理](#工作原理-21)
      - [完整示例](#完整示例-42)
      - [最佳实践](#最佳实践-42)
    - [8.3 内存同步](#83-内存同步)
      - [概念定义](#概念定义-43)
      - [完整示例](#完整示例-43)
      - [最佳实践](#最佳实践-43)
    - [8.4 死锁检测与避免](#84-死锁检测与避免)
      - [概念定义](#概念定义-44)
      - [完整示例](#完整示例-44)
      - [最佳实践](#最佳实践-44)
    - [8.5 活锁与饥饿](#85-活锁与饥饿)
      - [概念定义](#概念定义-45)
      - [完整示例](#完整示例-45)
      - [最佳实践](#最佳实践-45)
  - [总结](#总结)
    - [核心概念](#核心概念)
    - [设计原则](#设计原则)
    - [最佳实践](#最佳实践-46)
  - [参考资源](#参考资源)

---

## 1. Goroutine基础

### 1.1 Goroutine创建与生命周期

#### 概念定义

**Goroutine** 是Go语言中的轻量级线程，由Go运行时（runtime）管理而非操作系统内核管理。
它是Go并发模型的核心抽象，允许函数或方法与程序的其余部分并发执行。

#### 工作原理

Goroutine的实现基于**M:N调度模型**：

- **M** (Machine): 操作系统线程
- **P** (Processor): 逻辑处理器，维护本地可运行Goroutine队列
- **G** (Goroutine): 待执行的任务

```text
┌─────────────────────────────────────────────────────────┐
│                    Go Runtime Scheduler                 │
├─────────────────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐    │
│  │    P    │  │    P    │  │    P    │  │    P    │    │
│  │  ┌───┐  │  │  ┌───┐  │  │  ┌───┐  │  │  ┌───┐  │    │
│  │  │ G │  │  │  │ G │  │  │  │ G │  │  │  │ G │  │    │
│  │  │ G │  │  │  │ G │  │  │  │ G │  │  │  │ G │  │    │
│  │  └───┘  │  │  └───┘  │  │  └───┘  │  │  └───┘  │    │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘    │
│       └─────────────┴─────────────┴─────────────┘       │
│                         │                               │
│                    ┌────┴────┐                          │
│                    │ Global  │                          │
│                    │  Queue  │                          │
│                    └─────────┘                          │
│  ┌─────────────────────────────────────────────────┐   │
│  │  M (OS Thread)  M (OS Thread)  M (OS Thread)    │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

**调度机制**：

1. **Work Stealing**: 当P的本地队列为空时，从其他P窃取Goroutine
2. **Handoff**: 当Goroutine阻塞时，M与P分离，P寻找新的M继续执行
3. **Sysmon**: 系统监控线程，处理长时间运行的Goroutine抢占

#### 形式论证

**定理1.1**: Goroutine的创建是O(1)操作。

**证明**:

- Goroutine初始栈大小为2KB（远小于线程的MB级栈）
- 创建仅需分配栈空间和初始化少量寄存器
- 不涉及系统调用，完全在用户态完成
- 因此时间复杂度为常数时间

**定理1.2**: Go调度器保证公平性。

**证明**:

- 每个P维护本地队列，减少锁竞争
- 全局队列用于平衡负载
- Sysmon定期检查，防止Goroutine饿死

#### 完整示例

```go
package main

import (
 "fmt"
 "runtime"
 "sync"
 "time"
)

// 基本Goroutine创建
func basicGoroutine() {
 // 使用go关键字创建Goroutine
 go func() {
  fmt.Println("Hello from goroutine!")
 }()

 // 主Goroutine需要等待，否则程序可能立即退出
 time.Sleep(100 * time.Millisecond)
}

// 带参数的Goroutine
func parameterizedGoroutine() {
 messages := []string{"Hello", "World", "from", "Go"}

 var wg sync.WaitGroup
 for _, msg := range messages {
  wg.Add(1)
  // 必须传递参数，避免闭包陷阱
  go func(m string) {
   defer wg.Done()
   fmt.Println(m)
  }(msg)
 }
 wg.Wait()
}

// Goroutine生命周期管理
func lifecycleDemo() {
 // 使用channel进行同步
 done := make(chan struct{})

 go func() {
  fmt.Println("Goroutine started")
  time.Sleep(500 * time.Millisecond)
  fmt.Println("Goroutine finished")
  close(done) // 通知完成
 }()

 <-done // 等待Goroutine完成
 fmt.Println("Main goroutine continues")
}

// 获取Goroutine信息
func goroutineInfo() {
 fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
 fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
 fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
}

func main() {
 fmt.Println("=== Basic Goroutine ===")
 basicGoroutine()

 fmt.Println("\n=== Parameterized Goroutine ===")
 parameterizedGoroutine()

 fmt.Println("\n=== Lifecycle Demo ===")
 lifecycleDemo()

 fmt.Println("\n=== Goroutine Info ===")
 goroutineInfo()
}
```

#### 反例说明

```go
// ❌ 错误：闭包陷阱
func closureTrapWrong() {
 for i := 0; i < 5; i++ {
  go func() {
   // 所有Goroutine共享同一个i变量
   fmt.Println(i) // 可能输出相同的值或超出范围
  }()
 }
 time.Sleep(time.Second)
}

// ❌ 错误：变量地址陷阱
func addressTrapWrong() {
 var wg sync.WaitGroup
 for i := 0; i < 5; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   fmt.Println(&i) // 所有Goroutine打印相同的地址
  }()
 }
 wg.Wait()
}

// ✅ 正确：传递参数
func closureTrapCorrect() {
 for i := 0; i < 5; i++ {
  go func(n int) {
   fmt.Println(n) // 每个Goroutine有自己的n副本
  }(i)
 }
 time.Sleep(time.Second)
}
```

**问题分析**:

1. **闭包陷阱**: 循环变量在Goroutine启动时才被捕获，此时循环可能已结束
2. **地址陷阱**: 所有Goroutine共享同一个变量地址

#### 性能分析

| 指标 | Goroutine | OS Thread |
|------|-----------|-----------|
| 初始栈大小 | 2KB | 1-8MB |
| 创建时间 | ~2μs | ~1ms |
| 上下文切换 | ~200ns | ~1-2μs |
| 最大数量 | 数百万 | 数千 |
| 调度方式 | 用户态 | 内核态 |

#### 最佳实践

1. **始终传递循环变量作为参数**
2. **使用sync.WaitGroup等待Goroutine完成**
3. **避免在Goroutine中panic，使用recover**
4. **合理设置GOMAXPROCS**

---

### 1.2 Goroutine与线程区别

#### 概念定义

| 特性 | Goroutine | OS Thread |
|------|-----------|-----------|
| 调度 | Go运行时调度（用户态） | 操作系统调度（内核态） |
| 栈管理 | 动态增长/收缩（2KB起） | 固定大小（通常MB级） |
| 创建成本 | 极低（~2KB内存） | 较高（~1MB内存+内核资源） |
| 切换成本 | ~200ns | ~1-2μs |
| 通信方式 | Channel（CSP模型） | 共享内存+锁 |
| 标识 | 无唯一标识（设计哲学） | 有线程ID |

#### 工作原理

**Goroutine调度器演进**:

```
Go 1.0: 单全局队列 + 单M锁
        ↓
Go 1.1: 多M多P，但G在全局队列
        ↓
Go 1.2: 每个P有本地队列，Work Stealing
        ↓
Go 1.5: GOMAXPROCS默认等于CPU数
        ↓
Go 1.14: 基于信号的抢占式调度
```

**抢占式调度**:

```
Before Go 1.14: 协作式调度（函数调用点检查）
After Go 1.14:  信号驱动的抢占式调度

┌────────────────────────────────────────┐
│           Goroutine Execution          │
├────────────────────────────────────────┤
│  func longRunning() {                  │
│      for {                             │
│          // 无函数调用 = 无法抢占       │
│          // Go 1.14+: 可被信号抢占      │
│      }                                 │
│  }                                     │
└────────────────────────────────────────┘
```

#### 完整示例

```go
package main

import (
 "fmt"
 "runtime"
 "sync"
 "time"
)

// 对比Goroutine和线程的创建成本
func creationBenchmark() {
 const n = 10000

 // Goroutine创建
 start := time.Now()
 var wg sync.WaitGroup
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   _ = 1 + 1
  }()
 }
 wg.Wait()
 fmt.Printf("%d goroutines: %v\n", n, time.Since(start))

 // 对比：线程创建（使用runtime.LockOSThread模拟）
 // 实际线程创建需要CGO或syscall，这里仅作概念说明
}

// 展示Goroutine的轻量级特性
func lightweightDemo() {
 fmt.Println("Creating 100,000 goroutines...")

 var wg sync.WaitGroup
 sem := make(chan struct{}, 10000) // 限制并发数

 for i := 0; i < 100000; i++ {
  wg.Add(1)
  sem <- struct{}{}
  go func(n int) {
   defer wg.Done()
   defer func() { <-sem }()
   time.Sleep(10 * time.Millisecond)
  }(i)
 }

 wg.Wait()
 fmt.Println("All goroutines completed")
}

// 展示调度行为
func schedulerDemo() {
 runtime.GOMAXPROCS(4) // 使用4个逻辑处理器

 var wg sync.WaitGroup
 for i := 0; i < 8; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   for j := 0; j < 3; j++ {
    fmt.Printf("Goroutine %d: iteration %d on %d\n",
     id, j, runtime.GOMAXPROCS(0))
    time.Sleep(10 * time.Millisecond)
   }
  }(i)
 }
 wg.Wait()
}

func main() {
 fmt.Println("=== Creation Benchmark ===")
 creationBenchmark()

 fmt.Println("\n=== Lightweight Demo ===")
 lightweightDemo()

 fmt.Println("\n=== Scheduler Demo ===")
 schedulerDemo()
}
```

#### 性能分析

**内存使用对比**:

```
100,000 Goroutines: ~200MB (2KB each)
100,000 OS Threads: ~100GB (1MB each) - 实际上不可能创建
```

**上下文切换对比**:

```
Goroutine切换: ~200ns (用户态)
OS Thread切换: ~1-2μs (涉及内核态)
```

#### 最佳实践

1. **不要滥用runtime.LockOSThread()**
   - 会绑定Goroutine到特定OS线程
   - 破坏调度器优化

2. **合理设置GOMAXPROCS**
   - 默认值（CPU核心数）通常最优
   - I/O密集型可适当增加

3. **避免在Goroutine中做大量计算而不yield**
   - 虽然Go 1.14+支持抢占，但仍应设计良好

---

### 1.3 Goroutine泄漏与检测

#### 概念定义

**Goroutine泄漏**指Goroutine由于某种原因无法退出，持续占用内存和调度资源。与内存泄漏类似，但泄漏的是执行上下文而非堆内存。

#### 工作原理

**常见泄漏场景**:

```
1. Channel发送/接收阻塞
   ┌─────────┐      ┌─────────┐
   │   G1    │ ──→  │ Channel │  (无人接收)
   │ (send)  │      │ (full)  │
   └─────────┘      └─────────┘

2. 无限等待锁
   ┌─────────┐      ┌─────────┐
   │   G1    │ ──→  │  Mutex  │  (被G2持有不释放)
   │ (Lock)  │      │ (locked)│
   └─────────┘      └─────────┘

3. 无限循环
   ┌─────────┐
   │   G1    │ ──→ for { ... }  (无退出条件)
   │ (loop)  │
   └─────────┘

4. select阻塞
   ┌─────────┐
   │   G1    │ ──→ select{}  (无case可执行)
   │ (select)│
   └─────────┘
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "net/http"
 _ "net/http/pprof"
 "runtime"
 "sync"
 "time"
)

// 泄漏场景1: Channel发送阻塞
func leakChannelSend() {
 ch := make(chan int)
 go func() {
  ch <- 42 // 无人接收，永久阻塞
 }()
}

// 泄漏场景2: Channel接收阻塞
func leakChannelRecv() {
 ch := make(chan int)
 go func() {
  <-ch // 无人发送，永久阻塞
 }()
}

// 泄漏场景3: 锁未释放
func leakMutex() {
 var mu sync.Mutex
 mu.Lock()
 go func() {
  mu.Lock() // 等待锁，但主Goroutine已退出
  defer mu.Unlock()
 }()
}

// 泄漏场景4: 无限循环
func leakInfiniteLoop() {
 go func() {
  for {
   // 无退出条件
   time.Sleep(100 * time.Millisecond)
  }
 }()
}

// ✅ 正确做法1: 使用带缓冲的channel
func correctBuffered() {
 ch := make(chan int, 1) // 缓冲为1
 go func() {
  ch <- 42 // 不会阻塞
 }()
}

// ✅ 正确做法2: 使用select+超时
func correctTimeout() {
 ch := make(chan int)
 go func() {
  select {
  case ch <- 42:
  case <-time.After(time.Second):
   // 超时退出
  }
 }()
}

// ✅ 正确做法3: 使用context取消
func correctContext() {
 ctx, cancel := context.WithCancel(context.Background())
 defer cancel()

 go func(ctx context.Context) {
  for {
   select {
   case <-ctx.Done():
    return // 收到取消信号，退出
   default:
    // 执行任务
    time.Sleep(100 * time.Millisecond)
   }
  }
 }(ctx)
}

// 泄漏检测函数
func detectLeak() {
 // 获取当前Goroutine数量
 before := runtime.NumGoroutine()
 fmt.Printf("Goroutines before: %d\n", before)

 // 执行可能泄漏的代码
 leakChannelSend()
 time.Sleep(100 * time.Millisecond)

 after := runtime.NumGoroutine()
 fmt.Printf("Goroutines after: %d\n", after)

 if after > before {
  fmt.Printf("Possible leak: %d goroutines\n", after-before)
 }
}

// 使用pprof检测泄漏
func startPprof() {
 go func() {
  http.ListenAndServe("localhost:6060", nil)
 }()
 fmt.Println("pprof started at http://localhost:6060/debug/pprof/goroutine")
}

func main() {
 // 启动pprof服务用于检测
 startPprof()

 fmt.Println("=== Leak Detection ===")
 detectLeak()

 fmt.Println("\n=== Testing Correct Patterns ===")
 correctBuffered()
 correctTimeout()
 correctContext()

 time.Sleep(200 * time.Millisecond)
 fmt.Printf("Final goroutines: %d\n", runtime.NumGoroutine())

 // 保持程序运行以便查看pprof
 time.Sleep(30 * time.Second)
}
```

#### 检测方法

```bash
# 1. 使用runtime.NumGoroutine()
go test -v -run TestLeak

# 2. 使用pprof
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 3. 使用goleak测试库
go.uber.org/goleak
```

#### 最佳实践

1. **使用context管理Goroutine生命周期**
2. **设置合理的超时机制**
3. **使用带缓冲的channel避免阻塞**
4. **在测试中集成泄漏检测**

---

### 1.4 Goroutine数量控制

#### 概念定义

**Goroutine数量控制**是通过限制并发执行的Goroutine数量，防止资源耗尽和系统过载的技术。

#### 工作原理

**控制策略**:

```
1. 固定大小Worker Pool
   ┌─────────────────────────────────────┐
   │           Worker Pool (N)            │
   │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐   │
   │  │  W  │ │  W  │ │  W  │ │  W  │   │
   │  │  o  │ │  o  │ │  o  │ │  o  │   │
   │  │  r  │ │  r  │ │  r  │ │  r  │   │
   │  │  k  │ │  k  │ │  k  │ │  k  │   │
   │  │  e  │ │  e  │ │  e  │ │  e  │   │
   │  │  r  │ │  r  │ │  r  │ │  r  │   │
   │  └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘   │
   │     └───────┴───────┴───────┘       │
   │              │                       │
   │         ┌────┴────┐                  │
   │         │  Task   │                  │
   │         │  Queue  │                  │
   │         └─────────┘                  │
   └─────────────────────────────────────┘

2. 信号量(Semaphore)控制
   ┌─────────────────────────────────────┐
   │        Semaphore (Capacity N)        │
   │                                      │
   │   ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐    │
   │   │ 1 │ │ 1 │ │ 1 │ │ 0 │ │ 0 │    │
   │   └───┘ └───┘ └───┘ └───┘ └───┘    │
   │    ↑    ↑    ↑                      │
   │   可用槽位: 3                        │
   └─────────────────────────────────────┘

3. 自适应控制
   ┌─────────────────────────────────────┐
   │     Adaptive Controller              │
   │                                      │
   │   ┌─────────┐    ┌─────────┐        │
   │   │  Load   │───→│  Scale  │        │
   │   │ Monitor │    │  Logic  │        │
   │   └─────────┘    └────┬────┘        │
   │                        │             │
   │                   ┌────┴────┐        │
   │                   │ Workers │        │
   │                   └─────────┘        │
   └─────────────────────────────────────┘
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "runtime"
 "sync"
 "sync/atomic"
 "time"
)

// 1. 使用Channel实现信号量
type Semaphore struct {
 sem chan struct{}
}

func NewSemaphore(n int) *Semaphore {
 return &Semaphore{sem: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire() {
 s.sem <- struct{}{}
}

func (s *Semaphore) Release() {
 <-s.sem
}

func semaphoreExample() {
 sem := NewSemaphore(10) // 最多10个并发

 var wg sync.WaitGroup
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func(n int) {
   defer wg.Done()
   sem.Acquire()
   defer sem.Release()

   // 执行任务
   fmt.Printf("Task %d executing, active: %d\n",
    n, len(sem.sem))
   time.Sleep(100 * time.Millisecond)
  }(i)
 }
 wg.Wait()
}

// 2. 固定大小Worker Pool
type WorkerPool struct {
 workers   int
 taskQueue chan func()
 wg        sync.WaitGroup
}

func NewWorkerPool(workers, queueSize int) *WorkerPool {
 return &WorkerPool{
  workers:   workers,
  taskQueue: make(chan func(), queueSize),
 }
}

func (wp *WorkerPool) Start() {
 for i := 0; i < wp.workers; i++ {
  wp.wg.Add(1)
  go func(id int) {
   defer wp.wg.Done()
   for task := range wp.taskQueue {
    fmt.Printf("Worker %d processing task\n", id)
    task()
   }
  }(i)
 }
}

func (wp *WorkerPool) Submit(task func()) bool {
 select {
 case wp.taskQueue <- task:
  return true
 default:
  return false // 队列满
 }
}

func (wp *WorkerPool) Stop() {
 close(wp.taskQueue)
 wp.wg.Wait()
}

func workerPoolExample() {
 pool := NewWorkerPool(5, 100)
 pool.Start()

 for i := 0; i < 50; i++ {
  n := i
  pool.Submit(func() {
   fmt.Printf("Task %d completed\n", n)
   time.Sleep(50 * time.Millisecond)
  })
 }

 pool.Stop()
}

// 3. 带超时的Worker Pool
type TimeoutWorkerPool struct {
 workers   int
 taskQueue chan func()
 ctx       context.Context
 cancel    context.CancelFunc
 wg        sync.WaitGroup
}

func NewTimeoutWorkerPool(workers int, timeout time.Duration) *TimeoutWorkerPool {
 ctx, cancel := context.WithTimeout(context.Background(), timeout)
 return &TimeoutWorkerPool{
  workers:   workers,
  taskQueue: make(chan func()),
  ctx:       ctx,
  cancel:    cancel,
 }
}

func (wp *TimeoutWorkerPool) Start() {
 for i := 0; i < wp.workers; i++ {
  wp.wg.Add(1)
  go func(id int) {
   defer wp.wg.Done()
   for {
    select {
    case task, ok := <-wp.taskQueue:
     if !ok {
      return
     }
     task()
    case <-wp.ctx.Done():
     fmt.Printf("Worker %d timeout\n", id)
     return
    }
   }
  }(i)
 }
}

func (wp *TimeoutWorkerPool) Submit(task func()) {
 select {
 case wp.taskQueue <- task:
 case <-wp.ctx.Done():
  fmt.Println("Submit timeout")
 }
}

func (wp *TimeoutWorkerPool) Stop() {
 wp.cancel()
 close(wp.taskQueue)
 wp.wg.Wait()
}

// 4. 自适应Goroutine池
type AdaptivePool struct {
 minWorkers    int32
 maxWorkers    int32
 currentWorker int32
 taskQueue     chan func()
 wg            sync.WaitGroup
 mu            sync.Mutex
}

func NewAdaptivePool(min, max int) *AdaptivePool {
 return &AdaptivePool{
  minWorkers: int32(min),
  maxWorkers: int32(max),
  taskQueue:  make(chan func(), max*2),
 }
}

func (ap *AdaptivePool) Start() {
 // 启动最小数量的worker
 for i := 0; i < int(ap.minWorkers); i++ {
  ap.addWorker()
 }

 // 监控并动态调整
 go ap.monitor()
}

func (ap *AdaptivePool) addWorker() {
 ap.mu.Lock()
 defer ap.mu.Unlock()

 if atomic.LoadInt32(&ap.currentWorker) >= ap.maxWorkers {
  return
 }

 atomic.AddInt32(&ap.currentWorker, 1)
 ap.wg.Add(1)

 go func(id int32) {
  defer ap.wg.Done()
  defer atomic.AddInt32(&ap.currentWorker, -1)

  for task := range ap.taskQueue {
   task()
  }
 }(ap.currentWorker)
}

func (ap *AdaptivePool) monitor() {
 ticker := time.NewTicker(time.Second)
 defer ticker.Stop()

 for range ticker.C {
  queueLen := len(ap.taskQueue)
  current := atomic.LoadInt32(&ap.currentWorker)

  // 如果队列积压且未达到最大值，增加worker
  if queueLen > int(current)*2 && current < ap.maxWorkers {
   ap.addWorker()
   fmt.Printf("Scaled up: %d workers, queue: %d\n", current+1, queueLen)
  }
 }
}

func (ap *AdaptivePool) Submit(task func()) {
 ap.taskQueue <- task
}

func (ap *AdaptivePool) Stop() {
 close(ap.taskQueue)
 ap.wg.Wait()
}

func main() {
 fmt.Println("=== Semaphore Example ===")
 semaphoreExample()

 fmt.Println("\n=== Worker Pool Example ===")
 workerPoolExample()

 fmt.Printf("\nFinal goroutines: %d\n", runtime.NumGoroutine())
}
```

#### 性能分析

| 控制方式 | 优点 | 缺点 | 适用场景 |
|---------|------|------|---------|
| 信号量 | 简单、轻量 | 无任务队列 | 简单限流 |
| 固定Pool | 资源可控 | 不够灵活 | CPU密集型 |
| 自适应Pool | 动态调整 | 实现复杂 | 负载波动大 |

#### 最佳实践

1. **根据任务类型选择控制策略**
   - CPU密集型: workers ≈ CPU核心数
   - I/O密集型: workers可更大

2. **设置合理的队列大小**
   - 太小: 任务被拒绝
   - 太大: 内存占用高

3. **监控Pool状态**
   - 队列长度
   - 活跃worker数
   - 任务处理延迟

---

## 2. Channel模式

### 2.1 无缓冲Channel

#### 概念定义

**无缓冲Channel**（Unbuffered Channel）是容量为0的Channel，发送操作必须等待接收方准备好，实现**同步通信**。

#### 工作原理

```
无缓冲Channel的握手过程:

发送方G1                    Channel                    接收方G2
   │                          │                          │
   │──── ch <- v ────────────→│                          │
   │                          │                          │
   │◄──── 阻塞 ───────────────│                          │
   │                          │                          │
   │                          │◄─────────── <-ch ────────│
   │                          │                          │
   │◄──── 解除阻塞 ───────────│                          │
   │                          │                          │
   │                          │──── 数据传递 ────────────→│
   │                          │                          │
   │                          │                          │◄── 继续执行
   │◄── 继续执行              │                          │
```

**状态转换**:

```
Empty (无发送者, 无接收者)
    │
    │ 发送者到达
    ↓
Send (有发送者等待)
    │
    │ 接收者到达
    ↓
Empty ──→ 数据传递 ──→ 双方继续
```

#### 形式论证

**定理2.1**: 无缓冲Channel保证发送和接收的**happens-before**关系。

**证明**:

- 发送操作 `ch <- v` 在接收操作 `<-ch` 完成之前阻塞
- 接收操作 `<-ch` 在发送操作 `ch <- v` 完成之前阻塞
- 因此，发送和接收之间存在同步点
- 根据Go内存模型，这建立了happens-before关系

**定理2.2**: 无缓冲Channel实现**CSP（Communicating Sequential Processes）**模型。

**证明**:

- CSP要求通信是同步的
- 无缓冲Channel的发送必须等待接收
- 这确保了通信双方在同一时刻交互
- 符合CSP的同步通信语义

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 基本无缓冲Channel
func basicUnbuffered() {
 ch := make(chan int) // 无缓冲

 go func() {
  fmt.Println("Sender: sending 42")
  ch <- 42 // 阻塞，直到有接收者
  fmt.Println("Sender: sent successfully")
 }()

 time.Sleep(100 * time.Millisecond) // 确保sender先执行
 fmt.Println("Receiver: about to receive")
 v := <-ch // 接收，解除sender阻塞
 fmt.Printf("Receiver: received %d\n", v)
}

// 同步屏障模式
func synchronizationBarrier() {
 ch := make(chan struct{})

 var wg sync.WaitGroup
 for i := 0; i < 3; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   fmt.Printf("Worker %d: doing work...\n", id)
   time.Sleep(time.Duration(id*100) * time.Millisecond)
   fmt.Printf("Worker %d: waiting at barrier\n", id)
   ch <- struct{}{} // 到达屏障
   <-ch              // 等待其他worker
   fmt.Printf("Worker %d: passed barrier\n", id)
  }(i)
 }

 // 等待所有worker到达
 for i := 0; i < 3; i++ {
  <-ch
 }
 // 释放所有worker
 for i := 0; i < 3; i++ {
  ch <- struct{}{}
 }

 wg.Wait()
}

// 请求-响应模式
func requestResponse() {
 type Request struct {
  Data     int
  Response chan int
 }

 requestCh := make(chan Request)

 // 服务端
 go func() {
  for req := range requestCh {
   result := req.Data * 2
   req.Response <- result
  }
 }()

 // 客户端
 var wg sync.WaitGroup
 for i := 1; i <= 5; i++ {
  wg.Add(1)
  go func(n int) {
   defer wg.Done()

   responseCh := make(chan int)
   requestCh <- Request{Data: n, Response: responseCh}
   result := <-responseCh

   fmt.Printf("Request %d * 2 = %d\n", n, result)
  }(i)
 }

 wg.Wait()
}

// 信号通知模式
func signalNotification() {
 done := make(chan struct{})

 go func() {
  fmt.Println("Worker: starting long task...")
  time.Sleep(500 * time.Millisecond)
  fmt.Println("Worker: task completed")
  close(done) // 通知完成
 }()

 fmt.Println("Main: waiting for worker...")
 <-done // 等待信号
 fmt.Println("Main: worker done, continuing...")
}

func main() {
 fmt.Println("=== Basic Unbuffered Channel ===")
 basicUnbuffered()

 fmt.Println("\n=== Synchronization Barrier ===")
 synchronizationBarrier()

 fmt.Println("\n=== Request-Response Pattern ===")
 requestResponse()

 fmt.Println("\n=== Signal Notification ===")
 signalNotification()
}
```

#### 反例说明

```go
// ❌ 错误：死锁 - 无接收者的发送
func deadlockNoReceiver() {
 ch := make(chan int)
 ch <- 42 // 永久阻塞，无接收者
 fmt.Println("This will never print")
}

// ❌ 错误：死锁 - 无发送者的接收
func deadlockNoSender() {
 ch := make(chan int)
 <-ch // 永久阻塞，无发送者
 fmt.Println("This will never print")
}

// ❌ 错误：单向发送导致死锁
func deadlockSingleDirection() {
 ch := make(chan int)

 // 发送
 go func() {
  for i := 0; i < 3; i++ {
   ch <- i
  }
  // 忘记关闭channel
 }()

 // 接收不足
 fmt.Println(<-ch)
 fmt.Println(<-ch)
 // 第三个值无人接收，sender阻塞
 // 程序死锁
}

// ✅ 正确：确保有对应的接收者
func correctPaired() {
 ch := make(chan int)

 go func() {
  ch <- 42 // 发送
 }()

 fmt.Println(<-ch) // 接收
}

// ✅ 正确：使用buffered channel或关闭channel
func correctBufferedOrClose() {
 ch := make(chan int, 3) // 有缓冲

 go func() {
  for i := 0; i < 3; i++ {
   ch <- i
  }
  close(ch) // 完成后关闭
 }()

 for v := range ch {
  fmt.Println(v)
 }
}
```

#### 性能分析

| 特性 | 无缓冲Channel | 有缓冲Channel |
|------|--------------|--------------|
| 同步性 | 同步（阻塞） | 异步（非阻塞至满） |
| 延迟 | 高（需等待） | 低（缓冲命中） |
| 内存 | 低（无缓冲） | 高（缓冲大小） |
| 适用场景 | 同步、信号传递 | 解耦、批量处理 |

#### 最佳实践

1. **用于同步场景**：确保两个Goroutine在同一时刻交互
2. **用于信号传递**：close(channel)作为广播信号
3. **避免在单Goroutine中既发送又接收**
4. **确保有对应的接收者或发送者**

---

### 2.2 有缓冲Channel

#### 概念定义

**有缓冲Channel**（Buffered Channel）是具有固定容量的Channel，允许发送者在缓冲区未满时非阻塞发送，接收者在缓冲区非空时非阻塞接收。

#### 工作原理

```
有缓冲Channel的状态:

Empty (len=0, cap=N)
    │
    │ 发送 (缓冲区有空位)
    ↓
Partial (0 < len < N)
    │
    │ 继续发送至满
    ↓
Full (len=N, cap=N)
    │
    │ 发送阻塞，等待接收
    ↓
    │ 接收释放空间
    ↓
Partial

缓冲区实现:
┌─────────────────────────────────────┐
│         Circular Buffer              │
│                                      │
│   ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐    │
│   │ A │ │ B │ │ C │ │   │ │   │    │
│   └───┘ └───┘ └───┘ └───┘ └───┘    │
│    ↑         ↑                       │
│   send      recv                    │
│   index    index                     │
│                                      │
│   len=3, cap=5                       │
└─────────────────────────────────────┘
```

#### 形式论证

**定理2.3**: 有缓冲Channel的happens-before关系仅在缓冲区满/空时建立。

**证明**:

- 当 `len < cap` 时，发送不阻塞，无同步
- 当 `len > 0` 时，接收不阻塞，无同步
- 仅当 `len == cap`（发送阻塞）或 `len == 0`（接收阻塞）时建立同步
- 因此，有缓冲Channel的同步是条件性的

**定理2.4**: 有缓冲Channel实现**异步消息传递**。

**证明**:

- 发送者无需等待接收者就绪
- 消息暂存于缓冲区
- 接收者按FIFO顺序获取
- 符合异步消息队列语义

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 基本有缓冲Channel
func basicBuffered() {
 ch := make(chan int, 3) // 容量为3

 // 非阻塞发送（缓冲未满）
 ch <- 1
 ch <- 2
 ch <- 3
 fmt.Println("Sent 3 values without blocking")

 // 阻塞发送（缓冲已满）
 go func() {
  fmt.Println("Trying to send 4...")
  ch <- 4 // 阻塞，等待接收
  fmt.Println("Sent 4")
 }()

 time.Sleep(100 * time.Millisecond)

 // 接收释放空间
 fmt.Printf("Received: %d\n", <-ch)
 time.Sleep(50 * time.Millisecond) // 让sender完成
}

// 批量处理模式
func batchProcessing() {
 const batchSize = 5
 ch := make(chan int, batchSize)

 // 生产者
 go func() {
  for i := 1; i <= 20; i++ {
   ch <- i
   if i%batchSize == 0 {
    fmt.Printf("Batch of %d sent\n", batchSize)
   }
  }
  close(ch)
 }()

 // 消费者 - 批量处理
 var batch []int
 for v := range ch {
  batch = append(batch, v)
  if len(batch) == batchSize {
   fmt.Printf("Processing batch: %v\n", batch)
   batch = batch[:0] // 清空
  }
 }

 // 处理剩余
 if len(batch) > 0 {
  fmt.Printf("Processing remaining: %v\n", batch)
 }
}

// 工作队列模式
func workQueue() {
 const (
  numWorkers = 3
  queueSize  = 10
 )

 type Task struct {
  ID   int
  Data string
 }

 taskQueue := make(chan Task, queueSize)
 var wg sync.WaitGroup

 // 启动workers
 for i := 0; i < numWorkers; i++ {
  wg.Add(1)
  go func(workerID int) {
   defer wg.Done()
   for task := range taskQueue {
    fmt.Printf("Worker %d processing task %d: %s\n",
     workerID, task.ID, task.Data)
    time.Sleep(100 * time.Millisecond) // 模拟工作
   }
  }(i)
 }

 // 提交任务（非阻塞，直到队列满）
 for i := 0; i < 15; i++ {
  select {
  case taskQueue <- Task{ID: i, Data: fmt.Sprintf("data-%d", i)}:
   fmt.Printf("Task %d queued\n", i)
  default:
   fmt.Printf("Task %d dropped (queue full)\n", i)
  }
 }

 close(taskQueue)
 wg.Wait()
}

// 速率限制模式
func rateLimiting() {
 const (
  requestsPerSecond = 5
  burstSize         = 3
 )

 // 令牌桶
 tokens := make(chan struct{}, burstSize)

 // 填充令牌
 go func() {
  ticker := time.NewTicker(time.Second / requestsPerSecond)
  defer ticker.Stop()
  for range ticker.C {
   select {
   case tokens <- struct{}{}:
   default: // 桶满，丢弃令牌
   }
  }
 }()

 // 初始填充
 for i := 0; i < burstSize; i++ {
  tokens <- struct{}{}
 }

 // 处理请求
 for i := 0; i < 15; i++ {
  <-tokens // 获取令牌
  fmt.Printf("Request %d processed at %v\n", i, time.Now().Format("15:04:05.000"))
 }
}

func main() {
 fmt.Println("=== Basic Buffered Channel ===")
 basicBuffered()

 fmt.Println("\n=== Batch Processing ===")
 batchProcessing()

 fmt.Println("\n=== Work Queue ===")
 workQueue()

 fmt.Println("\n=== Rate Limiting ===")
 rateLimiting()
}
```

#### 反例说明

```go
// ❌ 错误：缓冲区过大导致内存浪费
func wasteMemory() {
 ch := make(chan int, 1000000) // 过大的缓冲
 // 实际只使用少量
 ch <- 1
 ch <- 2
 // 浪费大量内存
}

// ❌ 错误：缓冲区过小失去异步优势
func tooSmallBuffer() {
 ch := make(chan int, 1) // 几乎等同于无缓冲
 // 频繁阻塞
}

// ❌ 错误：忘记关闭导致goroutine泄漏
func forgetClose() {
 ch := make(chan int, 10)

 go func() {
  for range ch { // 永久等待
   // 处理
  }
 }()

 ch <- 1
 ch <- 2
 // 忘记close(ch) - goroutine泄漏
}

// ✅ 正确：合理设置缓冲区大小
func correctBufferSize() {
 // 根据实际场景设置
 ch := make(chan int, 100) // 合理的缓冲

 go func() {
  for v := range ch {
   _ = v
  }
 }()

 // 使用
 for i := 0; i < 10; i++ {
  ch <- i
 }
 close(ch)
}
```

#### 性能分析

**缓冲区大小选择**:

| 场景 | 推荐大小 | 原因 |
|------|---------|------|
| 同步协调 | 0 | 需要强同步 |
| 批量处理 | 批量大小的倍数 | 匹配处理单元 |
| 速率限制 | burst大小 | 允许突发流量 |
| 事件流 | 100-1000 | 平滑流量波动 |

#### 最佳实践

1. **根据场景选择缓冲区大小**
   - 同步场景: 0
   - 批量处理: 批量大小的倍数
   - 解耦: 根据生产消费速度差

2. **使用select处理满缓冲**

   ```go
   select {
   case ch <- v:
   default:
       // 缓冲满的处理
   }
   ```

3. **及时关闭channel**
   - 通知接收者数据结束
   - 避免goroutine泄漏

---

### 2.3 单向Channel

#### 概念定义

**单向Channel**是对Channel方向的限制，分为**只发送**（send-only）和**只接收**（receive-only）两种类型，用于编译期约束Channel使用方式。

#### 工作原理

```
类型定义:
- chan T      : 双向Channel
- chan<- T    : 只发送Channel
- <-chan T    : 只接收Channel

函数签名约束:
┌─────────────────────────────────────────┐
│ func producer(ch chan<- int)            │
│     │                                   │
│     └── 只能发送，编译器保证            │
│                                         │
│ func consumer(ch <-chan int)            │
│     │                                   │
│     └── 只能接收，编译器保证            │
└─────────────────────────────────────────┘

类型转换:
双向 ──→ 单向: 允许（隐式转换）
单向 ──→ 双向: 不允许（编译错误）
```

#### 形式论证

**定理2.5**: 单向Channel提供**编译期类型安全**。

**证明**:

- `chan<- T` 类型禁止接收操作
- `<-chan T` 类型禁止发送操作
- 违反约束在编译时报错
- 因此错误在运行前被发现

**定理2.6**: 单向Channel实现**接口隔离原则**。

**证明**:

- 生产者只需发送能力
- 消费者只需接收能力
- 单向类型精确表达需求
- 符合最小权限原则

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 生产者 - 只发送
func producer(ch chan<- int, id int, wg *sync.WaitGroup) {
 defer wg.Done()
 for i := 0; i < 3; i++ {
  value := id*10 + i
  fmt.Printf("Producer %d sending: %d\n", id, value)
  ch <- value
  time.Sleep(50 * time.Millisecond)
 }
}

// 消费者 - 只接收
func consumer(ch <-chan int, id int, wg *sync.WaitGroup) {
 defer wg.Done()
 for value := range ch {
  fmt.Printf("Consumer %d received: %d\n", id, value)
  time.Sleep(100 * time.Millisecond)
 }
}

// 转换器 - 接收一个channel，发送另一个
func transformer(in <-chan int, out chan<- int, wg *sync.WaitGroup) {
 defer wg.Done()
 defer close(out)

 for v := range in {
  out <- v * 2 // 转换数据
 }
}

// 多生产者单消费者
func multiProducerSingleConsumer() {
 ch := make(chan int, 10)
 var wg sync.WaitGroup

 // 启动多个生产者
 for i := 0; i < 3; i++ {
  wg.Add(1)
  go producer(ch, i, &wg)
 }

 // 等待生产者完成，然后关闭channel
 go func() {
  wg.Wait()
  close(ch)
 }()

 // 消费者
 for v := range ch {
  fmt.Printf("Received: %d\n", v)
 }
}

// Pipeline模式
func pipelinePattern() {
 // Stage 1: 生成数据
 generator := func() <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for i := 1; i <= 5; i++ {
    out <- i
   }
  }()
  return out
 }

 // Stage 2: 平方
 square := func(in <-chan int) <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for v := range in {
    out <- v * v
   }
  }()
  return out
 }

 // Stage 3: 加倍
 double := func(in <-chan int) <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for v := range in {
    out <- v * 2
   }
  }()
  return out
 }

 // 连接pipeline
 ch := double(square(generator()))

 for v := range ch {
  fmt.Printf("Result: %d\n", v)
 }
}

// 类型安全演示
func typeSafetyDemo() {
 ch := make(chan int)

 // 可以隐式转换为单向
 var sendOnly chan<- int = ch
 var recvOnly <-chan int = ch

 // 使用
 go func() {
  sendOnly <- 42 // 只能发送
 }()

 fmt.Println(<-recvOnly) // 只能接收

 // 以下会编译错误:
 // <-sendOnly  // 错误: 不能从send-only channel接收
 // recvOnly <- 1  // 错误: 不能向receive-only channel发送
}

func main() {
 fmt.Println("=== Multi Producer Single Consumer ===")
 multiProducerSingleConsumer()

 fmt.Println("\n=== Pipeline Pattern ===")
 pipelinePattern()

 fmt.Println("\n=== Type Safety Demo ===")
 typeSafetyDemo()
}
```

#### 反例说明

```go
// ❌ 错误：尝试从send-only channel接收
func wrongReceive() {
 ch := make(chan int)
 sendCh := chan<- int(ch)

 // <-sendCh  // 编译错误: invalid operation: cannot receive from send-only channel
}

// ❌ 错误：尝试向receive-only channel发送
func wrongSend() {
 ch := make(chan int)
 recvCh := <-chan int(ch)

 // recvCh <- 1  // 编译错误: invalid operation: cannot send to receive-only channel
}

// ❌ 错误：单向转双向
func wrongConversion() {
 ch := make(chan int)
 sendCh := chan<- int(ch)

 // bidirectional := chan int(sendCh)  // 编译错误: cannot convert
}

// ✅ 正确：使用单向channel作为函数参数
func correctUsage() {
 process := func(send chan<- int, recv <-chan int) {
  for v := range recv {
   send <- v * 2
  }
 }

 ch1 := make(chan int)
 ch2 := make(chan int)

 go process(ch2, ch1)

 ch1 <- 1
 ch1 <- 2
 close(ch1)

 fmt.Println(<-ch2)
 fmt.Println(<-ch2)
}
```

#### 最佳实践

1. **函数参数使用单向Channel**
   - 明确表达意图
   - 编译期检查

2. **Pipeline每个stage返回receive-only**

   ```go
   func stage(in <-chan T) <-chan T
   ```

3. **不要在API中暴露双向Channel**
   - 限制调用者操作
   - 防止误用

---

### 2.4 Channel关闭与检测

#### 概念定义

**Channel关闭**是发送方通知接收方不再有数据发送的机制。关闭后的Channel可以无限次接收零值，用于广播通知。

#### 工作原理

```
Channel关闭语义:

发送方:                     接收方:
close(ch) ──────────────→  <-ch 返回 (value, false)
   │                          │
   │ 关闭后禁止发送           │ 关闭后可无限接收
   │ ch <- v 会panic          │ 返回零值和false

关闭状态检测:
┌─────────────────────────────────────────┐
│ v, ok := <-ch                          │
│                                         │
│ ok == true  : channel未关闭，v是实际值  │
│ ok == false : channel已关闭，v是零值    │
└─────────────────────────────────────────┘

广播通知模式:
┌─────────────────────────────────────────┐
│ done := make(chan struct{})             │
│                                         │
│ close(done)  // 广播给所有接收者        │
│                                         │
│ G1 ◄──┐                                 │
│ G2 ◄──┼── 所有goroutine同时收到通知     │
│ G3 ◄──┘                                 │
└─────────────────────────────────────────┘
```

#### 形式论证

**定理2.7**: 关闭Channel实现**广播语义**。

**证明**:

- 关闭后，所有等待的接收者都被唤醒
- 后续接收立即返回（零值, false）
- 不保证接收顺序，但保证所有接收者都被通知
- 因此实现了一对多广播

**定理2.8**: Channel关闭是**幂等**的。

**证明**:

- 多次close同一channel会导致panic
- 因此需要保证只关闭一次
- 通常由发送方负责关闭
- 或使用sync.Once确保

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 基本关闭操作
func basicClose() {
 ch := make(chan int, 3)

 ch <- 1
 ch <- 2
 ch <- 3
 close(ch)

 // 接收直到关闭
 for {
  v, ok := <-ch
  if !ok {
   fmt.Println("Channel closed")
   break
  }
  fmt.Printf("Received: %d\n", v)
 }
}

// 使用range简化
func rangeClose() {
 ch := make(chan int, 3)

 ch <- 1
 ch <- 2
 ch <- 3
 close(ch)

 // range自动检测关闭
 for v := range ch {
  fmt.Printf("Received: %d\n", v)
 }
 fmt.Println("Channel closed (range exited)")
}

// 广播通知模式
func broadcastNotification() {
 notify := make(chan struct{})
 var wg sync.WaitGroup

 // 启动多个worker
 for i := 0; i < 5; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()

   <-notify // 等待通知
   fmt.Printf("Worker %d notified at %v\n", id, time.Now().UnixNano())
  }(i)
 }

 time.Sleep(100 * time.Millisecond)
 fmt.Println("Broadcasting...")
 close(notify) // 广播通知

 wg.Wait()
}

// 安全关闭模式 - sync.Once
func safeCloseOnce() {
 ch := make(chan int)
 var once sync.Once

 closeFn := func() {
  once.Do(func() {
   close(ch)
   fmt.Println("Channel closed safely")
  })
 }

 // 多个goroutine尝试关闭
 var wg sync.WaitGroup
 for i := 0; i < 3; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   closeFn()
  }()
 }
 wg.Wait()
}

// 检测channel是否关闭（不推荐）
func isClosed(ch <-chan int) bool {
 select {
 case <-ch:
  return true
 default:
  return false
 }
}

// 更好的模式 - 使用额外channel
func closeDetection() {
 dataCh := make(chan int)
 doneCh := make(chan struct{})

 go func() {
  defer close(doneCh)
  for v := range dataCh {
   fmt.Printf("Processing: %d\n", v)
  }
 }()

 // 发送数据
 for i := 0; i < 5; i++ {
  dataCh <- i
 }
 close(dataCh)

 // 等待处理完成
 <-doneCh
 fmt.Println("All data processed")
}

// 优雅关闭模式
func gracefulShutdown() {
 jobs := make(chan int, 10)
 done := make(chan struct{})
 var wg sync.WaitGroup

 // Worker
 wg.Add(1)
 go func() {
  defer wg.Done()
  for {
   select {
   case job, ok := <-jobs:
    if !ok {
     fmt.Println("Worker: jobs channel closed, exiting")
     return
    }
    fmt.Printf("Worker: processing job %d\n", job)
    time.Sleep(100 * time.Millisecond)
   case <-done:
    fmt.Println("Worker: shutdown signal received")
    // 处理剩余jobs
    for job := range jobs {
     fmt.Printf("Worker: finishing job %d\n", job)
    }
    return
   }
  }
 }()

 // 提交jobs
 for i := 0; i < 10; i++ {
  jobs <- i
 }

 // 方式1: 正常完成
 close(jobs)

 // 方式2: 优雅关闭（取消上面，启用下面）
 // close(done)

 wg.Wait()
}

func main() {
 fmt.Println("=== Basic Close ===")
 basicClose()

 fmt.Println("\n=== Range Close ===")
 rangeClose()

 fmt.Println("\n=== Broadcast Notification ===")
 broadcastNotification()

 fmt.Println("\n=== Safe Close Once ===")
 safeCloseOnce()

 fmt.Println("\n=== Close Detection ===")
 closeDetection()

 fmt.Println("\n=== Graceful Shutdown ===")
 gracefulShutdown()
}
```

#### 反例说明

```go
// ❌ 错误：向已关闭的channel发送
func sendToClosed() {
 ch := make(chan int)
 close(ch)
 ch <- 1 // panic: send on closed channel
}

// ❌ 错误：重复关闭channel
func doubleClose() {
 ch := make(chan int)
 close(ch)
 close(ch) // panic: close of closed channel
}

// ❌ 错误：多个goroutine关闭（竞态）
func raceClose() {
 ch := make(chan int)

 go func() { close(ch) }()
 go func() { close(ch) }() // 可能panic
}

// ❌ 错误：接收方关闭channel
func receiverClose() {
 ch := make(chan int)

 go func() {
  ch <- 1
 }()

 <-ch
 close(ch) // 语义错误：接收方不应关闭
}

// ✅ 正确：发送方关闭
func senderClose() {
 ch := make(chan int)

 go func() {
  defer close(ch) // 发送方关闭
  for i := 0; i < 5; i++ {
   ch <- i
  }
 }()

 for v := range ch {
  fmt.Println(v)
 }
}

// ✅ 正确：使用sync.Once
func onceClose() {
 ch := make(chan int)
 var once sync.Once

 closeSafe := func() {
  once.Do(func() { close(ch) })
 }

 // 可以安全地多次调用
 closeSafe()
 closeSafe()
 closeSafe()
}
```

#### 最佳实践

1. **发送方负责关闭Channel**
2. **使用`range`接收直到关闭**
3. **使用`sync.Once`确保只关闭一次**
4. **不要依赖关闭状态检测（isClosed）**
5. **使用`select`处理关闭和取消**

---

### 2.5 Channel Select多路复用

#### 概念定义

**Select**是Go语言提供的多路复用原语，允许同时等待多个Channel操作，哪个就绪就执行哪个，实现非确定性选择和超时控制。

#### 工作原理

```
Select执行流程:

1. 评估所有case表达式
2. 如果有多个case就绪，随机选择一个执行
3. 如果没有case就绪，阻塞等待
4. 如果有default，立即执行（非阻塞）

伪代码实现:
func select(cases []Case) {
    // 1. 锁定所有涉及的channel
    lockAllChannels(cases)

    // 2. 检查就绪的case
    ready := findReadyCases(cases)

    if len(ready) > 0 {
        // 3. 随机选择一个
        chosen := randomSelect(ready)
        unlockAllChannels(cases)
        execute(chosen)
    } else if hasDefault(cases) {
        // 4. 执行default
        unlockAllChannels(cases)
        executeDefault()
    } else {
        // 5. 阻塞等待
        parkGoroutine(cases)
        unlockAllChannels(cases)
    }
}

随机选择保证公平性:
case ch1 <- v1:
case ch2 <- v2:
case v3 := <-ch3:

如果ch1和ch2都可写，随机选择其中一个
避免某个channel被饿死
```

#### 形式论证

**定理2.9**: Select实现**公平的多路复用**。

**证明**:

- 多个就绪case时，使用均匀随机选择
- 每个就绪case被选中的概率相等
- 避免任何channel被优先或饿死
- 因此实现公平性

**定理2.10**: Select的default case实现**非阻塞操作**。

**证明**:

- 当所有case都阻塞时，如果有default则立即执行
- 不等待任何channel就绪
- 因此是非阻塞的
- 可用于实现try-send/try-receive

#### 完整示例

```go
package main

import (
 "fmt"
 "math/rand"
 "time"
)

// 基本select
func basicSelect() {
 ch1 := make(chan string)
 ch2 := make(chan string)

 go func() {
  time.Sleep(100 * time.Millisecond)
  ch1 <- "from ch1"
 }()

 go func() {
  time.Sleep(200 * time.Millisecond)
  ch2 <- "from ch2"
 }()

 // 等待任一channel
 select {
 case msg1 := <-ch1:
  fmt.Println(msg1)
 case msg2 := <-ch2:
  fmt.Println(msg2)
 }
}

// 非阻塞select (default)
func nonBlockingSelect() {
 ch := make(chan int)

 // 非阻塞发送
 select {
 case ch <- 1:
  fmt.Println("Sent successfully")
 default:
  fmt.Println("Channel full, dropped")
 }

 // 非阻塞接收
 select {
 case v := <-ch:
  fmt.Printf("Received: %d\n", v)
 default:
  fmt.Println("No data available")
 }
}

// 随机选择演示
func randomSelection() {
 ch := make(chan int, 10)

 // 填充channel
 for i := 0; i < 10; i++ {
  ch <- i
 }

 // 多次select，观察随机性
 counts := make(map[int]int)
 for i := 0; i < 1000; i++ {
  select {
  case <-ch:
   counts[0]++
  case <-ch:
   counts[1]++
  case <-ch:
   counts[2]++
  }
  // 重新填充
  for len(ch) < 10 {
   ch <- 1
  }
 }

 fmt.Println("Selection counts:", counts)
}

// 多路复用 - 合并channels
func multiplexChannels() {
 ch1 := make(chan int)
 ch2 := make(chan int)
 out := make(chan int)

 // 合并goroutine
 go func() {
  defer close(out)
  for {
   select {
   case v, ok := <-ch1:
    if !ok {
     ch1 = nil // 禁用此case
     continue
    }
    out <- v
   case v, ok := <-ch2:
    if !ok {
     ch2 = nil // 禁用此case
     continue
    }
    out <- v
   }
   if ch1 == nil && ch2 == nil {
    return
   }
  }
 }()

 // 发送数据
 go func() {
  for i := 0; i < 3; i++ {
   ch1 <- i
  }
  close(ch1)
 }()

 go func() {
  for i := 10; i < 13; i++ {
   ch2 <- i
  }
  close(ch2)
 }()

 // 接收合并后的数据
 for v := range out {
  fmt.Printf("Received: %d\n", v)
 }
}

// 优先级select
func prioritySelect() {
 highPriority := make(chan string)
 lowPriority := make(chan string)

 go func() {
  for {
   // 优先处理高优先级
   select {
   case msg := <-highPriority:
    fmt.Printf("High priority: %s\n", msg)
   default:
    // 高优先级无数据，处理低优先级
    select {
    case msg := <-highPriority:
     fmt.Printf("High priority: %s\n", msg)
    case msg := <-lowPriority:
     fmt.Printf("Low priority: %s\n", msg)
    }
   }
  }
 }()

 // 发送消息
 lowPriority <- "low 1"
 highPriority <- "high 1"
 lowPriority <- "low 2"

 time.Sleep(100 * time.Millisecond)
}

// 动态case (反射实现)
func dynamicSelect(channels []chan int) {
 // 注意：实际应使用reflect.Select，但这里展示概念
 // 简单版本：轮询
 for {
  for i, ch := range channels {
   select {
   case v := <-ch:
    fmt.Printf("From channel %d: %d\n", i, v)
   default:
   }
  }
 }
}

func main() {
 fmt.Println("=== Basic Select ===")
 basicSelect()

 fmt.Println("\n=== Non-blocking Select ===")
 nonBlockingSelect()

 fmt.Println("\n=== Random Selection ===")
 randomSelection()

 fmt.Println("\n=== Multiplex Channels ===")
 multiplexChannels()

 fmt.Println("\n=== Priority Select ===")
 prioritySelect()
}
```

#### 反例说明

```go
// ❌ 错误：空的select
func emptySelect() {
 select {} // 永久阻塞
}

// ❌ 错误：所有case阻塞且无default
func allBlocking() {
 ch1 := make(chan int)
 ch2 := make(chan int)

 select {
 case <-ch1: // 阻塞
 case <-ch2: // 阻塞
 }
 // 永久阻塞
}

// ❌ 错误：在select中发送nil channel
func nilChannelSelect() {
 var ch chan int // nil

 select {
 case ch <- 1: // 永远不会执行
  fmt.Println("Sent")
 default:
  fmt.Println("Default") // 总是执行
 }
}

// ❌ 错误：在循环中无退出的select
func infiniteSelectLoop() {
 ch := make(chan int)

 for {
  select {
  case v := <-ch:
   fmt.Println(v)
   // 没有退出条件
  }
 }
}

// ✅ 正确：使用done channel退出
func correctSelectLoop() {
 ch := make(chan int)
 done := make(chan struct{})

 go func() {
  time.Sleep(500 * time.Millisecond)
  close(done)
 }()

 for {
  select {
  case v := <-ch:
   fmt.Println(v)
  case <-done:
   fmt.Println("Exiting")
   return
  }
 }
}
```

#### 最佳实践

1. **使用default实现非阻塞操作**
2. **使用nil channel禁用case**
3. **使用done channel实现优雅退出**
4. **避免空的select（除非有意阻塞）**
5. **注意select的随机性，不要依赖顺序**

---

### 2.6 超时模式

#### 概念定义

**超时模式**是通过select和time.After实现的操作超时控制，防止Goroutine无限期等待。

#### 工作原理

```
超时模式实现:

select {
case result := <-operation:
    // 操作完成
    process(result)
case <-time.After(timeout):
    // 超时处理
    handleTimeout()
}

time.After实现:
time.After(d) ──→ 创建timer ──→ 在d后发送时间到channel

Timer内部:
┌─────────────────────────────────────┐
│           time.Timer                │
│                                      │
│   ┌─────────┐      ┌─────────┐      │
│   │  Timer  │─────→│ Channel │      │
│   │ (heap)  │      │ (send)  │      │
│   └─────────┘      └─────────┘      │
│        ↑                             │
│   Go运行时统一管理                    │
└─────────────────────────────────────┘
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "time"
)

// 基本超时
func basicTimeout() {
 ch := make(chan string)

 go func() {
  time.Sleep(2 * time.Second)
  ch <- "result"
 }()

 select {
 case result := <-ch:
  fmt.Println("Result:", result)
 case <-time.After(1 * time.Second):
  fmt.Println("Timeout!")
 }
}

// 可复用timer（避免频繁创建）
func reusableTimer() {
 timer := time.NewTimer(1 * time.Second)
 defer timer.Stop()

 ch := make(chan string)

 go func() {
  time.Sleep(500 * time.Millisecond)
  ch <- "quick result"
 }()

 select {
 case result := <-ch:
  // 停止timer，释放资源
  if !timer.Stop() {
   <-timer.C
  }
  fmt.Println("Result:", result)
 case <-timer.C:
  fmt.Println("Timeout!")
 }
}

// 多级超时
func multiLevelTimeout() {
 ch := make(chan string)

 go func() {
  time.Sleep(1500 * time.Millisecond)
  ch <- "result"
 }()

 // 快速超时
 select {
 case result := <-ch:
  fmt.Println("Fast result:", result)
  return
 case <-time.After(500 * time.Millisecond):
  fmt.Println("Fast timeout, trying medium...")
 }

 // 中等超时
 select {
 case result := <-ch:
  fmt.Println("Medium result:", result)
  return
 case <-time.After(1 * time.Second):
  fmt.Println("Medium timeout, trying slow...")
 }

 // 慢速超时
 select {
 case result := <-ch:
  fmt.Println("Slow result:", result)
 case <-time.After(2 * time.Second):
  fmt.Println("Slow timeout!")
 }
}

// 使用context实现超时
func contextTimeout() {
 ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
 defer cancel()

 ch := make(chan string, 1)

 go func() {
  // 模拟长时间操作
  select {
  case <-time.After(2 * time.Second):
   ch <- "result"
  case <-ctx.Done():
   fmt.Println("Operation cancelled")
   return
  }
 }()

 select {
 case result := <-ch:
  fmt.Println("Result:", result)
 case <-ctx.Done():
  fmt.Println("Timeout:", ctx.Err())
 }
}

// 带超时的函数调用
type Result struct {
 Value string
 Error error
}

func callWithTimeout(fn func() (string, error), timeout time.Duration) (string, error) {
 resultCh := make(chan Result, 1)

 go func() {
  value, err := fn()
  resultCh <- Result{Value: value, Error: err}
 }()

 select {
 case result := <-resultCh:
  return result.Value, result.Error
 case <-time.After(timeout):
  return "", fmt.Errorf("operation timed out after %v", timeout)
 }
}

func timeoutFunctionDemo() {
 slowFunc := func() (string, error) {
  time.Sleep(2 * time.Second)
  return "success", nil
 }

 result, err := callWithTimeout(slowFunc, 1*time.Second)
 if err != nil {
  fmt.Println("Error:", err)
 } else {
  fmt.Println("Result:", result)
 }
}

// 超时后清理资源
func timeoutWithCleanup() {
 type Resource struct {
  ID int
 }

 ch := make(chan *Resource)

 go func() {
  // 分配资源
  res := &Resource{ID: 1}
  time.Sleep(2 * time.Second)
  ch <- res
 }()

 var res *Resource
 select {
 case res = <-ch:
  fmt.Println("Got resource:", res.ID)
 case <-time.After(1 * time.Second):
  fmt.Println("Timeout, resource will leak if not handled")
  // 需要额外的goroutine来清理
  go func() {
   r := <-ch
   fmt.Println("Cleaned up resource:", r.ID)
  }()
 }
}

func main() {
 fmt.Println("=== Basic Timeout ===")
 basicTimeout()

 fmt.Println("\n=== Reusable Timer ===")
 reusableTimer()

 fmt.Println("\n=== Multi-level Timeout ===")
 multiLevelTimeout()

 fmt.Println("\n=== Context Timeout ===")
 contextTimeout()

 fmt.Println("\n=== Timeout Function ===")
 timeoutFunctionDemo()

 fmt.Println("\n=== Timeout with Cleanup ===")
 timeoutWithCleanup()
}
```

#### 性能分析

| 方式 | 优点 | 缺点 |
|------|------|------|
| time.After | 简单 | 每次创建timer |
| time.NewTimer | 可复用 | 需手动管理 |
| context.WithTimeout | 可取消、可传递 | 稍复杂 |

#### 最佳实践

1. **优先使用context.WithTimeout**
2. **高频调用使用可复用timer**
3. **超时后考虑资源清理**
4. **设置合理的超时时间**

---

### 2.7 心跳模式

#### 概念定义

**心跳模式**（Heartbeat）是定期发送信号表明组件仍在运行的机制，用于健康检查、超时检测和连接保持。

#### 工作原理

```
心跳模式结构:

┌─────────────┐         ┌─────────────┐
│   Worker    │         │   Monitor   │
│             │         │             │
│  ┌───────┐  │  beat   │  ┌───────┐  │
│  │ Task  │  │ ──────→ │  │ Check │  │
│  │ Loop  │  │         │  │ Timer │  │
│  └───┬───┘  │         │  └───┬───┘  │
│      │      │         │      │      │
│  ┌───┴───┐  │         │  ┌───┴───┐  │
│  │ Heart │  │  beat   │  │ Alert │  │
│  │ beat  │  │ ──────→ │  │ if    │  │
│  │ Timer │  │         │  │ timeout│  │
│  └───────┘  │         │  └───────┘  │
└─────────────┘         └─────────────┘

心跳间隔 < 超时阈值（通常2-3倍关系）
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "time"
)

// 基本心跳
func basicHeartbeat() {
 heartbeat := make(chan struct{})
 done := make(chan struct{})

 // Worker
 go func() {
  defer close(done)
  ticker := time.NewTicker(500 * time.Millisecond)
  defer ticker.Stop()

  for i := 0; i < 5; i++ {
   select {
   case <-ticker.C:
    heartbeat <- struct{}{}
    fmt.Println("Heartbeat sent")
   }
  }
 }()

 // Monitor
 go func() {
  for {
   select {
   case <-heartbeat:
    fmt.Println("Heartbeat received at", time.Now().Format("15:04:05"))
   case <-done:
    fmt.Println("Worker done")
    return
   }
  }
 }()

 <-done
}

// 带超时检测的心跳
func heartbeatWithTimeout() {
 heartbeat := make(chan struct{})

 // Worker
 go func() {
  ticker := time.NewTicker(500 * time.Millisecond)
  defer ticker.Stop()

  for i := 0; i < 10; i++ {
   if i == 5 {
    // 模拟卡顿
    fmt.Println("Worker: simulating stall...")
    time.Sleep(2 * time.Second)
   }

   select {
   case <-ticker.C:
    select {
    case heartbeat <- struct{}{}:
    default: // 非阻塞发送
    }
   }
  }
  close(heartbeat)
 }()

 // Monitor with timeout
 timeout := 1 * time.Second
 for {
  select {
  case _, ok := <-heartbeat:
   if !ok {
    fmt.Println("Worker finished")
    return
   }
   fmt.Println("Heartbeat OK at", time.Now().Format("15:04:05"))
  case <-time.After(timeout):
   fmt.Println("ALERT: Heartbeat timeout!")
   return
  }
 }
}

// 工作单元+心跳
func workUnitWithHeartbeat() {
 type Work struct {
  ID   int
  Data string
 }

 type Result struct {
  WorkID int
  Output string
 }

 workCh := make(chan Work)
 resultCh := make(chan Result)
 heartbeat := make(chan int) // 发送workID

 // Worker
 go func() {
  defer close(resultCh)

  for work := range workCh {
   // 模拟长时间工作，定期发送心跳
   steps := 5
   for i := 0; i < steps; i++ {
    time.Sleep(200 * time.Millisecond)

    // 发送进度心跳
    select {
    case heartbeat <- work.ID:
    default:
    }

    fmt.Printf("Work %d: progress %d/%d\n", work.ID, i+1, steps)
   }

   resultCh <- Result{
    WorkID: work.ID,
    Output: fmt.Sprintf("processed-%s", work.Data),
   }
  }
 }()

 // 提交工作
 go func() {
  for i := 1; i <= 3; i++ {
   workCh <- Work{ID: i, Data: fmt.Sprintf("data-%d", i)}
  }
  close(workCh)
 }()

 // 监控心跳和结果
 for {
  select {
  case workID := <-heartbeat:
   fmt.Printf("Heartbeat: work %d is alive\n", workID)
  case result, ok := <-resultCh:
   if !ok {
    fmt.Println("All work completed")
    return
   }
   fmt.Printf("Result: work %d = %s\n", result.WorkID, result.Output)
  case <-time.After(2 * time.Second):
   fmt.Println("No heartbeat - possible deadlock!")
   return
  }
 }
}

// 使用context的心跳
func contextHeartbeat() {
 ctx, cancel := context.WithCancel(context.Background())
 defer cancel()

 heartbeat := make(chan struct{})

 // Heartbeat generator
 go func() {
  ticker := time.NewTicker(500 * time.Millisecond)
  defer ticker.Stop()

  for {
   select {
   case <-ticker.C:
    select {
    case heartbeat <- struct{}{}:
    case <-ctx.Done():
     return
    }
   case <-ctx.Done():
    return
   }
  }
 }()

 // Monitor
 go func() {
  timeout := 2 * time.Second
  for {
   select {
   case <-heartbeat:
    fmt.Println("Heartbeat OK")
   case <-time.After(timeout):
    fmt.Println("Heartbeat timeout, cancelling...")
    cancel()
    return
   case <-ctx.Done():
    return
   }
  }
 }()

 // 模拟工作
 select {
 case <-ctx.Done():
  fmt.Println("Work cancelled:", ctx.Err())
 case <-time.After(5 * time.Second):
  fmt.Println("Work completed")
 }
}

func main() {
 fmt.Println("=== Basic Heartbeat ===")
 basicHeartbeat()

 fmt.Println("\n=== Heartbeat with Timeout ===")
 heartbeatWithTimeout()

 fmt.Println("\n=== Work Unit with Heartbeat ===")
 workUnitWithHeartbeat()

 fmt.Println("\n=== Context Heartbeat ===")
 contextHeartbeat()
}
```

#### 最佳实践

1. **心跳间隔 < 超时阈值**（通常1:2或1:3）
2. **使用非阻塞发送避免worker阻塞**
3. **心跳channel应带缓冲**
4. **结合context实现取消**

---

### 2.8 流水线模式

#### 概念定义

**流水线模式**（Pipeline）是将数据处理分解为多个阶段（stage），每个阶段通过Channel连接，实现数据流的顺序处理。

#### 工作原理

```
流水线结构:

Input ──→ [Stage 1] ──→ [Stage 2] ──→ [Stage 3] ──→ Output
              │             │             │
              ↓             ↓             ↓
           Channel       Channel       Channel

每个Stage:
┌─────────────────────────────────────────┐
│ func stage(in <-chan T) <-chan U {      │
│     out := make(chan U)                 │
│     go func() {                         │
│         defer close(out)                │
│         for v := range in {             │
│             out <- process(v)           │
│         }                               │
│     }()                                 │
│     return out                          │
│ }                                       │
└─────────────────────────────────────────┘

组合:
result := stage3(stage2(stage1(input)))
```

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
)

// Stage 1: 生成数据
func generator(nums ...int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  for _, n := range nums {
   out <- n
  }
 }()
 return out
}

// Stage 2: 平方
func square(in <-chan int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  for n := range in {
   out <- n * n
  }
 }()
 return out
}

// Stage 3: 过滤（只保留偶数）
func filterEven(in <-chan int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  for n := range in {
   if n%2 == 0 {
    out <- n
   }
  }
 }()
 return out
}

// Stage 4: 求和
func sum(in <-chan int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  total := 0
  for n := range in {
   total += n
  }
  out <- total
 }()
 return out
}

// 合并多个channels
func merge(channels ...<-chan int) <-chan int {
 out := make(chan int)
 var wg sync.WaitGroup

 output := func(c <-chan int) {
  defer wg.Done()
  for n := range c {
   out <- n
  }
 }

 wg.Add(len(channels))
 for _, c := range channels {
  go output(c)
 }

 go func() {
  wg.Wait()
  close(out)
 }()

 return out
}

// 带缓冲的pipeline（提高吞吐量）
func bufferedPipeline() {
 generator := func() <-chan int {
  out := make(chan int, 10)
  go func() {
   defer close(out)
   for i := 1; i <= 20; i++ {
    out <- i
   }
  }()
  return out
 }

 process := func(in <-chan int) <-chan int {
  out := make(chan int, 10)
  go func() {
   defer close(out)
   for n := range in {
    out <- n * 2
   }
  }()
  return out
 }

 ch := process(generator())
 for n := range ch {
  fmt.Printf("%d ", n)
 }
 fmt.Println()
}

// 并行Pipeline（Fan-out/Fan-in）
func parallelPipeline() {
 // Generator
 generator := func() <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for i := 1; i <= 20; i++ {
    out <- i
   }
  }()
  return out
 }

 // Worker
 worker := func(in <-chan int, id int) <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for n := range in {
    fmt.Printf("Worker %d processing %d\n", id, n)
    out <- n * n
   }
  }()
  return out
 }

 // Fan-out
 in := generator()
 numWorkers := 3
 var channels []<-chan int
 for i := 0; i < numWorkers; i++ {
  channels = append(channels, worker(in, i))
 }

 // Fan-in
 out := merge(channels...)

 for n := range out {
  fmt.Printf("Result: %d\n", n)
 }
}

// 带取消的Pipeline
func cancellablePipeline() {
 done := make(chan struct{})

 // 模拟取消
 go func() {
  // 稍后取消
  // close(done)
 }()

 generator := func() <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for i := 1; i <= 100; i++ {
    select {
    case out <- i:
    case <-done:
     return
    }
   }
  }()
  return out
 }

 process := func(in <-chan int) <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for n := range in {
    select {
    case out <- n * 2:
    case <-done:
     return
    }
   }
  }()
  return out
 }

 for n := range process(generator()) {
  fmt.Printf("%d ", n)
  if n > 20 {
   // close(done) // 取消
   break
  }
 }
}

func main() {
 fmt.Println("=== Basic Pipeline ===")
 // generator -> square -> filter -> sum
 result := <-sum(filterEven(square(generator(1, 2, 3, 4, 5))))
 fmt.Printf("Result: %d\n", result)

 fmt.Println("\n=== Buffered Pipeline ===")
 bufferedPipeline()

 fmt.Println("\n=== Parallel Pipeline ===")
 parallelPipeline()
}
```

#### 性能分析

| 模式 | 吞吐量 | 延迟 | 适用场景 |
|------|--------|------|---------|
| 串行Pipeline | 低 | 低 | 简单处理 |
| 缓冲Pipeline | 中 | 中 | 平滑流量 |
| 并行Pipeline | 高 | 中 | CPU密集型 |

#### 最佳实践

1. **每个stage独立关闭输出channel**
2. **使用单向channel类型**
3. **适当使用缓冲提高吞吐量**
4. **支持取消机制**
5. **使用Fan-out/Fan-in提高并行度**

---

## 3. 同步原语（sync包）

### 3.1 Mutex与RWMutex

#### 概念定义

**Mutex**（互斥锁）提供对共享资源的独占访问，**RWMutex**（读写锁）允许多个读操作并发执行，但写操作独占。

#### 工作原理

```
Mutex状态机:

Unlocked ──Lock()──→ Locked ──Unlock()──→ Unlocked
                          ↑                │
                          └────────────────┘

RWMutex状态机:

┌─────────────────────────────────────────────────────┐
│                  RWMutex States                      │
├─────────────────────────────────────────────────────┤
│  Unlocked: 无读者，无写者                            │
│       │                                             │
│       ├── RLock() ──→ Reading (n readers)          │
│       │                    │                        │
│       │                    ├── RUnlock() ──→ 计数-1 │
│       │                    │                        │
│       │                    └── 计数=0 ──→ Unlocked  │
│       │                                             │
│       └── Lock() ──→ Writing (1 writer)            │
│                           │                        │
│                           └── Unlock() ──→ Unlocked│
└─────────────────────────────────────────────────────┘

写者优先 vs 读者优先:
- Go的RWMutex: 写者优先（防止写饥饿）
- 新读者在有等待写者时阻塞
```

#### 形式论证

**定理3.1**: Mutex保证**互斥性**。

**证明**:

- Lock()原子地将状态从Unlocked变为Locked
- 如果已为Locked，调用者阻塞
- Unlock()原子地将状态变回Unlocked
- 因此任意时刻最多一个Goroutine持有锁

**定理3.2**: RWMutex保证**读写互斥，读读兼容**。

**证明**:

- 读者计数器记录当前读者数
- 写者需要读者计数为0且无其他写者
- 多个读者可同时持有RLock
- 写者独占，与所有读者互斥

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 基本Mutex使用
type Counter struct {
 mu    sync.Mutex
 value int
}

func (c *Counter) Increment() {
 c.mu.Lock()
 defer c.mu.Unlock()
 c.value++
}

func (c *Counter) Value() int {
 c.mu.Lock()
 defer c.mu.Unlock()
 return c.value
}

// RWMutex使用
type Cache struct {
 mu    sync.RWMutex
 data  map[string]string
}

func NewCache() *Cache {
 return &Cache{data: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
 c.mu.RLock()
 defer c.mu.RUnlock()
 val, ok := c.data[key]
 return val, ok
}

func (c *Cache) Set(key, value string) {
 c.mu.Lock()
 defer c.mu.Unlock()
 c.data[key] = value
}

func (c *Cache) Delete(key string) {
 c.mu.Lock()
 defer c.mu.Unlock()
 delete(c.data, key)
}

// 嵌套锁（危险）
type NestedStruct struct {
 mu     sync.Mutex
 value  int
 other  *NestedStruct
}

func (n *NestedStruct) Sum() int {
 n.mu.Lock()
 defer n.mu.Unlock()

 sum := n.value
 if n.other != nil {
  // 可能导致死锁！
  sum += n.other.Sum()
 }
 return sum
}

// 锁排序避免死锁
func (n *NestedStruct) SumSafe() int {
 // 获取两个锁的顺序
 first, second := n, n.other
 if second != nil && first > second {
  first, second = second, first
 }

 first.mu.Lock()
 defer first.mu.Unlock()

 if second != nil {
  second.mu.Lock()
  defer second.mu.Unlock()
 }

 sum := n.value
 if n.other != nil {
  sum += n.other.value
 }
 return sum
}

// 锁粒度优化
type FineGrainedCache struct {
 mu     sync.RWMutex
 shards [16]shard
}

type shard struct {
 mu   sync.RWMutex
 data map[string]string
}

func (c *FineGrainedCache) getShard(key string) *shard {
 // 根据key选择shard
 hash := 0
 for i := 0; i < len(key); i++ {
  hash += int(key[i])
 }
 return &c.shards[hash%16]
}

func (c *FineGrainedCache) Get(key string) (string, bool) {
 s := c.getShard(key)
 s.mu.RLock()
 defer s.mu.RUnlock()
 val, ok := s.data[key]
 return val, ok
}

// 性能对比
func benchmarkMutex() {
 const n = 100000

 // Mutex
 counter := &Counter{}
 start := time.Now()
 var wg sync.WaitGroup
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter.Increment()
  }()
 }
 wg.Wait()
 fmt.Printf("Mutex: %v, value=%d\n", time.Since(start), counter.Value())

 // RWMutex - 多读
 cache := NewCache()
 cache.Set("key", "value")

 start = time.Now()
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   cache.Get("key")
  }()
 }
 wg.Wait()
 fmt.Printf("RWMutex (read-only): %v\n", time.Since(start))
}

func main() {
 fmt.Println("=== Basic Mutex ===")
 counter := &Counter{}
 var wg sync.WaitGroup
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter.Increment()
  }()
 }
 wg.Wait()
 fmt.Printf("Counter: %d\n", counter.Value())

 fmt.Println("\n=== RWMutex Cache ===")
 cache := NewCache()
 cache.Set("name", "Go")
 if val, ok := cache.Get("name"); ok {
  fmt.Printf("Cache hit: %s\n", val)
 }

 fmt.Println("\n=== Benchmark ===")
 benchmarkMutex()
}
```

#### 反例说明

```go
// ❌ 错误：重复解锁
func doubleUnlock() {
 var mu sync.Mutex
 mu.Lock()
 mu.Unlock()
 mu.Unlock() // panic: sync: unlock of unlocked mutex
}

// ❌ 错误：未锁定就解锁
func unlockWithoutLock() {
 var mu sync.Mutex
 mu.Unlock() // panic: sync: unlock of unlocked mutex
}

// ❌ 错误：拷贝Mutex
func copyMutex() {
 type Data struct {
  mu sync.Mutex
  v  int
 }

 d1 := Data{v: 1}
 d2 := d1 // 拷贝Mutex - 危险！

 d1.mu.Lock()
 // d2.mu也"被锁定"了（实际上状态混乱）
}

// ❌ 错误：死锁 - 循环依赖
func deadlock() {
 var mu1, mu2 sync.Mutex

 go func() {
  mu1.Lock()
  defer mu1.Unlock()
  time.Sleep(10 * time.Millisecond)
  mu2.Lock() // 等待mu2
  defer mu2.Unlock()
 }()

 mu2.Lock()
 defer mu2.Unlock()
 time.Sleep(10 * time.Millisecond)
 mu1.Lock() // 等待mu1 - 死锁！
 defer mu1.Unlock()
}

// ✅ 正确：使用defer确保解锁
func correctDefer() {
 var mu sync.Mutex
 mu.Lock()
 defer mu.Unlock()
 // 执行业务逻辑
}

// ✅ 正确：锁排序避免死锁
func correctOrdering() {
 var mu1, mu2 sync.Mutex

 // 总是先锁mu1，再锁mu2
 mu1.Lock()
 defer mu1.Unlock()
 mu2.Lock()
 defer mu2.Unlock()
}
```

#### 性能分析

| 锁类型 | 读并发 | 写并发 | 适用场景 |
|--------|--------|--------|---------|
| Mutex | 无 | 无 | 读写均衡 |
| RWMutex | 高 | 无 | 读多写少 |
| 无锁 | 最高 | 最高 | 简单操作 |

#### 最佳实践

1. **使用defer确保Unlock**
2. **锁保护范围尽量小**
3. **不要拷贝Mutex**
4. **避免锁嵌套**
5. **必要时使用锁排序**

---

### 3.2 WaitGroup

#### 概念定义

**WaitGroup**用于等待一组Goroutine完成，通过计数器跟踪未完成的Goroutine数量。

#### 工作原理

```
WaitGroup状态:

┌─────────────────────────────────────────┐
│            WaitGroup Counter            │
├─────────────────────────────────────────┤
│                                         │
│  Add(n) ──→ counter += n                │
│                                         │
│  Done() ──→ counter -= 1                │
│      (等价于 Add(-1))                   │
│                                         │
│  Wait() ──→ 阻塞直到 counter == 0       │
│                                         │
│  内部使用信号量实现:                     │
│  - sema: 等待者的信号量                  │
│  - state: 高32位=counter, 低32位=waiter │
│                                         │
└─────────────────────────────────────────┘
```

#### 完整示例

```go
package main

import (
 "fmt"
 "net/http"
 "sync"
 "time"
)

// 基本WaitGroup
func basicWaitGroup() {
 var wg sync.WaitGroup

 for i := 1; i <= 3; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   fmt.Printf("Worker %d starting\n", id)
   time.Sleep(time.Duration(id) * 100 * time.Millisecond)
   fmt.Printf("Worker %d done\n", id)
  }(i)
 }

 wg.Wait()
 fmt.Println("All workers completed")
}

// 并行下载
func parallelDownload(urls []string) {
 var wg sync.WaitGroup
 results := make(map[string]string)
 var mu sync.Mutex

 for _, url := range urls {
  wg.Add(1)
  go func(u string) {
   defer wg.Done()

   resp, err := http.Get(u)
   if err != nil {
    fmt.Printf("Error downloading %s: %v\n", u, err)
    return
   }
   defer resp.Body.Close()

   mu.Lock()
   results[u] = resp.Status
   mu.Unlock()

   fmt.Printf("Downloaded %s: %s\n", u, resp.Status)
  }(url)
 }

 wg.Wait()
 fmt.Printf("Total downloaded: %d\n", len(results))
}

// 错误处理 + WaitGroup
type Result struct {
 Value string
 Error error
}

func parallelWithError() {
 tasks := []func() (string, error){
  func() (string, error) { return "task1", nil },
  func() (string, error) { return "", fmt.Errorf("task2 failed") },
  func() (string, error) { return "task3", nil },
 }

 var wg sync.WaitGroup
 results := make(chan Result, len(tasks))

 for i, task := range tasks {
  wg.Add(1)
  go func(id int, t func() (string, error)) {
   defer wg.Done()
   val, err := t()
   results <- Result{Value: val, Error: err}
  }(i, task)
 }

 // 等待所有任务完成，然后关闭channel
 go func() {
  wg.Wait()
  close(results)
 }()

 // 收集结果
 for r := range results {
  if r.Error != nil {
   fmt.Printf("Error: %v\n", r.Error)
  } else {
   fmt.Printf("Success: %s\n", r.Value)
  }
 }
}

// 嵌套WaitGroup
func nestedWaitGroup() {
 var outerWg sync.WaitGroup

 for i := 0; i < 2; i++ {
  outerWg.Add(1)
  go func(outerID int) {
   defer outerWg.Done()

   var innerWg sync.WaitGroup
   for j := 0; j < 3; j++ {
    innerWg.Add(1)
    go func(innerID int) {
     defer innerWg.Done()
     fmt.Printf("Outer %d, Inner %d\n", outerID, innerID)
    }(j)
   }
   innerWg.Wait()
   fmt.Printf("Outer %d: all inner done\n", outerID)
  }(i)
 }

 outerWg.Wait()
 fmt.Println("All done")
}

// 超时WaitGroup
func timeoutWaitGroup() {
 var wg sync.WaitGroup
 done := make(chan struct{})

 // 启动长时间任务
 wg.Add(1)
 go func() {
  defer wg.Done()
  time.Sleep(2 * time.Second)
  fmt.Println("Task completed")
 }()

 // 等待完成或超时
 go func() {
  wg.Wait()
  close(done)
 }()

 select {
 case <-done:
  fmt.Println("All tasks completed")
 case <-time.After(1 * time.Second):
  fmt.Println("Timeout!")
 }
}

func main() {
 fmt.Println("=== Basic WaitGroup ===")
 basicWaitGroup()

 fmt.Println("\n=== Parallel with Error ===")
 parallelWithError()

 fmt.Println("\n=== Nested WaitGroup ===")
 nestedWaitGroup()

 fmt.Println("\n=== Timeout WaitGroup ===")
 timeoutWaitGroup()
}
```

#### 反例说明

```go
// ❌ 错误：负计数
func negativeCounter() {
 var wg sync.WaitGroup
 wg.Add(1)
 wg.Done()
 wg.Done() // panic: sync: negative WaitGroup counter
}

// ❌ 错误：Add在Wait之后
func addAfterWait() {
 var wg sync.WaitGroup

 go func() {
  wg.Wait()
  fmt.Println("Wait returned")
 }()

 wg.Add(1) // 竞态条件！
 go func() {
  defer wg.Done()
  fmt.Println("Done")
 }()
}

// ❌ 错误：WaitGroup拷贝
func copyWaitGroup() {
 var wg sync.WaitGroup

 worker := func(wg sync.WaitGroup) { // 拷贝！
  defer wg.Done()
  fmt.Println("Working")
 }

 wg.Add(1)
 go worker(wg)
 wg.Wait() // 永远等待（等待的是原始wg）
}

// ❌ 错误：忘记Done
func forgetDone() {
 var wg sync.WaitGroup
 wg.Add(1)
 go func() {
  // 忘记调用wg.Done()
  fmt.Println("Working")
 }()
 wg.Wait() // 永远等待
}

// ✅ 正确：使用指针传递
func correctPointer() {
 var wg sync.WaitGroup

 worker := func(wg *sync.WaitGroup) {
  defer wg.Done()
  fmt.Println("Working")
 }

 wg.Add(1)
 go worker(&wg)
 wg.Wait()
}

// ✅ 正确：Add在启动goroutine之前
func correctAdd() {
 var wg sync.WaitGroup

 for i := 0; i < 3; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   fmt.Println("Working")
  }()
 }
 wg.Wait()
}
```

#### 最佳实践

1. **Add在启动Goroutine之前**
2. **使用指针传递WaitGroup**
3. **使用defer确保Done**
4. **不要拷贝WaitGroup**
5. **结合channel处理错误**

---

### 3.3 Once

#### 概念定义

**Once**确保某段代码在程序运行期间只执行一次，常用于单例模式、延迟初始化等场景。

#### 工作原理

```
Once实现原理:

┌─────────────────────────────────────────┐
│              sync.Once                   │
├─────────────────────────────────────────┤
│  type Once struct {                     │
│      done uint32  // 原子标志            │
│      m    Mutex   // 互斥锁              │
│  }                                      │
│                                         │
│  func (o *Once) Do(f func()) {          │
│      // 快速路径：检查done               │
│      if atomic.LoadUint32(&o.done) == 0 │{
│          // 慢速路径：获取锁执行          │
│          o.m.Lock()                     │
│          defer o.m.Unlock()             │
│          if o.done == 0 {               │
│              defer atomic.StoreUint32(  │
│                  &o.done, 1)            │
│              f()                        │
│          }                              │
│      }                                  │
│  }                                      │
└─────────────────────────────────────────┘

双检锁（Double-Checked Locking）:
1. 第一次检查（无锁）：快速路径
2. 第二次检查（有锁）：确保只有一个执行
```

#### 完整示例

```go
package main

import (
 "fmt"
 "net"
 "sync"
 "time"
)

// 基本Once
func basicOnce() {
 var once sync.Once

 for i := 0; i < 5; i++ {
  once.Do(func() {
   fmt.Println("Executed only once")
  })
  fmt.Printf("Iteration %d\n", i)
 }
}

// 单例模式
var (
 instance *Singleton
 once     sync.Once
)

type Singleton struct {
 data string
}

func GetInstance() *Singleton {
 once.Do(func() {
  instance = &Singleton{data: "initialized"}
  fmt.Println("Singleton initialized")
 })
 return instance
}

// 延迟初始化
type LazyResource struct {
 value  int
 once   sync.Once
 initFn func()
}

func (r *LazyResource) Init() {
 r.once.Do(r.initFn)
}

func (r *LazyResource) Value() int {
 r.Init()
 return r.value
}

// 连接池初始化
type ConnectionPool struct {
 connections []net.Conn
 once        sync.Once
}

func (p *ConnectionPool) Initialize() {
 p.once.Do(func() {
  fmt.Println("Initializing connection pool...")
  p.connections = make([]net.Conn, 0, 10)
  // 实际初始化连接...
  time.Sleep(100 * time.Millisecond)
  fmt.Println("Connection pool initialized")
 })
}

// 配置加载
type Config struct {
 DatabaseURL string
 APIKey      string
}

var (
 config     *Config
 configOnce sync.Once
)

func LoadConfig() *Config {
 configOnce.Do(func() {
  fmt.Println("Loading config...")
  config = &Config{
   DatabaseURL: "postgres://localhost/db",
   APIKey:      "secret-key",
  }
  time.Sleep(50 * time.Millisecond)
  fmt.Println("Config loaded")
 })
 return config
}

// 并发测试
func concurrentOnce() {
 var once sync.Once
 var wg sync.WaitGroup

 counter := 0
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   once.Do(func() {
    counter++
    fmt.Printf("Executed by goroutine %d\n", id)
   })
  }(i)
 }
 wg.Wait()
 fmt.Printf("Counter: %d (should be 1)\n", counter)
}

// Once与错误处理
type SafeInitializer struct {
 once  sync.Once
 err   error
 value interface{}
}

func (s *SafeInitializer) Do(initFunc func() (interface{}, error)) (interface{}, error) {
 s.once.Do(func() {
  s.value, s.err = initFunc()
 })
 return s.value, s.err
}

func main() {
 fmt.Println("=== Basic Once ===")
 basicOnce()

 fmt.Println("\n=== Singleton Pattern ===")
 for i := 0; i < 3; i++ {
  _ = GetInstance()
 }

 fmt.Println("\n=== Lazy Resource ===")
 resource := &LazyResource{
  initFn: func() {
   fmt.Println("Initializing resource...")
   resource.value = 42
  },
 }
 fmt.Printf("Value: %d\n", resource.Value())
 fmt.Printf("Value again: %d\n", resource.Value())

 fmt.Println("\n=== Connection Pool ===")
 pool := &ConnectionPool{}
 var wg sync.WaitGroup
 for i := 0; i < 5; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   pool.Initialize()
  }()
 }
 wg.Wait()

 fmt.Println("\n=== Concurrent Once ===")
 concurrentOnce()
}
```

#### 反例说明

```go
// ❌ 错误：Once拷贝
func copyOnce() {
 once := sync.Once{}

 doSomething := func(o sync.Once) { // 拷贝
  o.Do(func() {
   fmt.Println("Executed")
  })
 }

 doSomething(once)
 doSomething(once) // 再次执行（不同的Once）
}

// ❌ 错误：Once中的panic
func panicInOnce() {
 var once sync.Once

 defer func() {
  if r := recover(); r != nil {
   fmt.Println("Recovered:", r)
  }
 }()

 once.Do(func() {
  panic("oops")
 })

 // 再次调用 - 不会执行（Once认为已执行过）
 once.Do(func() {
  fmt.Println("This won't execute")
 })
}

// ❌ 错误：递归Once
func recursiveOnce() {
 var once sync.Once

 var recursiveFunc func()
 recursiveFunc = func() {
  once.Do(func() {
   fmt.Println("Before recursion")
   recursiveFunc() // 递归调用 - 死锁！
   fmt.Println("After recursion")
  })
 }

 recursiveFunc()
}

// ✅ 正确：使用指针
func correctOncePointer() {
 var once sync.Once

 doSomething := func(o *sync.Once) {
  o.Do(func() {
   fmt.Println("Executed")
  })
 }

 doSomething(&once)
 doSomething(&once) // 不会再次执行
}
```

#### 最佳实践

1. **使用指针传递Once**
2. **Once中的函数不要panic**
3. **不要在Once中递归调用Once**
4. **Once不能重置，需要重置使用其他方案**

---

### 3.4 Cond

#### 概念定义

**Cond**（条件变量）用于在特定条件满足时唤醒等待的Goroutine，实现复杂的同步场景。

#### 工作原理

```
Cond结构:

┌─────────────────────────────────────────┐
│              sync.Cond                   │
├─────────────────────────────────────────┤
│  type Cond struct {                     │
│      L Locker      // 关联的锁           │
│      notify notifyList // 等待列表       │
│  }                                      │
│                                         │
│  Wait():                                │
│    1. 释放锁                            │
│    2. 加入等待列表                       │
│    3. 等待信号                           │
│    4. 重新获取锁                         │
│                                         │
│  Signal(): 唤醒一个等待者                │
│  Broadcast(): 唤醒所有等待者             │
└─────────────────────────────────────────┘

使用模式:
cond.L.Lock()
for !condition() {  // 必须用for，不能用if
    cond.Wait()
}
// 执行业务
cond.L.Unlock()
```

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 基本Cond - 生产者消费者
type Queue struct {
 items []int
 cond  *sync.Cond
}

func NewQueue() *Queue {
 return &Queue{
  items: make([]int, 0),
  cond:  sync.NewCond(&sync.Mutex{}),
 }
}

func (q *Queue) Enqueue(item int) {
 q.cond.L.Lock()
 defer q.cond.L.Unlock()

 q.items = append(q.items, item)
 fmt.Printf("Enqueued: %d, queue size: %d\n", item, len(q.items))
 q.cond.Signal() // 通知等待的消费者
}

func (q *Queue) Dequeue() int {
 q.cond.L.Lock()
 defer q.cond.L.Unlock()

 for len(q.items) == 0 {
  fmt.Println("Queue empty, waiting...")
  q.cond.Wait() // 等待信号
 }

 item := q.items[0]
 q.items = q.items[1:]
 fmt.Printf("Dequeued: %d, queue size: %d\n", item, len(q.items))
 return item
}

// 广播模式 - 资源就绪通知
type ResourcePool struct {
 resources []string
 ready     bool
 cond      *sync.Cond
}

func NewResourcePool() *ResourcePool {
 return &ResourcePool{
  resources: make([]string, 0),
  cond:      sync.NewCond(&sync.Mutex{}),
 }
}

func (p *ResourcePool) AddResource(r string) {
 p.cond.L.Lock()
 defer p.cond.L.Unlock()

 p.resources = append(p.resources, r)
 if !p.ready && len(p.resources) >= 3 {
  p.ready = true
  fmt.Println("Resources ready, broadcasting...")
  p.cond.Broadcast() // 通知所有等待者
 }
}

func (p *ResourcePool) WaitForReady() {
 p.cond.L.Lock()
 defer p.cond.L.Unlock()

 for !p.ready {
  fmt.Println("Waiting for resources to be ready...")
  p.cond.Wait()
 }
 fmt.Println("Resources are ready!")
}

// 多条件等待
type MultiCondExample struct {
 value     int
 target    int
 cond      *sync.Cond
}

func NewMultiCondExample() *MultiCondExample {
 return &MultiCondExample{
  cond: sync.NewCond(&sync.Mutex{}),
 }
}

func (m *MultiCondExample) SetTarget(t int) {
 m.cond.L.Lock()
 defer m.cond.L.Unlock()

 m.target = t
 m.cond.Broadcast() // 目标改变，通知所有等待者
}

func (m *MultiCondExample) Increment() {
 m.cond.L.Lock()
 defer m.cond.L.Unlock()

 m.value++
 if m.value >= m.target {
  m.cond.Broadcast()
 }
}

func (m *MultiCondExample) WaitForTarget() {
 m.cond.L.Lock()
 defer m.cond.L.Unlock()

 for m.value < m.target {
  m.cond.Wait()
 }
 fmt.Printf("Target reached: %d >= %d\n", m.value, m.target)
}

// 超时Cond（使用channel实现）
type TimeoutCond struct {
 cond *sync.Cond
}

func NewTimeoutCond() *TimeoutCond {
 return &TimeoutCond{
  cond: sync.NewCond(&sync.Mutex{}),
 }
}

func (t *TimeoutCond) WaitWithTimeout(timeout time.Duration) bool {
 timeoutCh := make(chan struct{})
 go func() {
  time.Sleep(timeout)
  close(timeoutCh)
 }()

 done := make(chan struct{})
 go func() {
  t.cond.L.Lock()
  t.cond.Wait()
  t.cond.L.Unlock()
  close(done)
 }()

 select {
 case <-done:
  return true // 正常唤醒
 case <-timeoutCh:
  return false // 超时
 }
}

func main() {
 fmt.Println("=== Basic Cond (Queue) ===")
 queue := NewQueue()

 // 消费者
 go func() {
  for i := 0; i < 5; i++ {
   queue.Dequeue()
   time.Sleep(100 * time.Millisecond)
  }
 }()

 // 生产者
 for i := 0; i < 5; i++ {
  queue.Enqueue(i)
  time.Sleep(150 * time.Millisecond)
 }

 time.Sleep(500 * time.Millisecond)

 fmt.Println("\n=== Broadcast Pattern ===")
 pool := NewResourcePool()

 // 多个等待者
 for i := 0; i < 3; i++ {
  go func(id int) {
   fmt.Printf("Worker %d: ", id)
   pool.WaitForReady()
  }(i)
 }

 time.Sleep(100 * time.Millisecond)
 pool.AddResource("r1")
 pool.AddResource("r2")
 pool.AddResource("r3")

 time.Sleep(500 * time.Millisecond)
}
```

#### 反例说明

```go
// ❌ 错误：用if而不是for
func wrongIfCond() {
 cond := sync.NewCond(&sync.Mutex{})
 ready := false

 go func() {
  cond.L.Lock()
  if !ready {  // 错误：应该用for
   cond.Wait()
  }
  // 可能被虚假唤醒，ready仍为false
  fmt.Println("Ready")
  cond.L.Unlock()
 }()
}

// ❌ 错误：忘记Lock/Unlock
func forgetLock() {
 cond := sync.NewCond(&sync.Mutex{})

 cond.Wait() // panic: sync: condition variable unlocked wait
}

// ❌ 错误：Signal在Lock之外
func signalOutsideLock() {
 cond := sync.NewCond(&sync.Mutex{})

 cond.Signal() // 不安全，应该在Lock保护下
}

// ✅ 正确：使用for循环
func correctForCond() {
 cond := sync.NewCond(&sync.Mutex{})
 ready := false

 go func() {
  cond.L.Lock()
  for !ready {  // 正确：用for处理虚假唤醒
   cond.Wait()
  }
  fmt.Println("Ready")
  cond.L.Unlock()
 }()
}
```

#### 最佳实践

1. **总是用for循环检查条件**
2. **在Lock保护下调用Wait/Signal/Broadcast**
3. **优先使用channel，Cond用于复杂场景**
4. **注意虚假唤醒（spurious wakeup）**

---

### 3.5 Pool

#### 概念定义

**Pool**是一组可复用的临时对象集合，用于减少GC压力和内存分配开销。

#### 工作原理

```
Pool结构:

┌─────────────────────────────────────────┐
│              sync.Pool                   │
├─────────────────────────────────────────┤
│  type Pool struct {                     │
│      New func() interface{}             │
│      local     []poolLocal  // P本地    │
│      victim    []poolLocal  // 上一周期  │
│  }                                        │
│                                         │
│  Get(): 从本地队列获取，或从victim获取   │
│  Put(): 放回本地队列                     │
│                                         │
│  GC时：local → victim, victim → 清空     │
└─────────────────────────────────────────┘

对象生命周期:
Get ──→ 使用 ──→ Put ──→ 可能被复用
                    ↓
              GC后可能释放
```

#### 完整示例

```go
package main

import (
 "bytes"
 "fmt"
 "sync"
 "time"
)

// 基本Pool
var bufferPool = sync.Pool{
 New: func() interface{} {
  fmt.Println("Creating new buffer")
  return new(bytes.Buffer)
 },
}

func useBuffer() {
 // 获取buffer
 buf := bufferPool.Get().(*bytes.Buffer)
 defer bufferPool.Put(buf)

 // 重置并复用
 buf.Reset()
 buf.WriteString("Hello, Pool!")
 fmt.Println(buf.String())
}

// 自定义对象池
type Connection struct {
 ID   int
 Addr string
}

var connectionPool = sync.Pool{
 New: func() interface{} {
  fmt.Println("Creating new connection")
  return &Connection{ID: 0, Addr: ""}
 },
}

func getConnection() *Connection {
 conn := connectionPool.Get().(*Connection)
 conn.ID++
 return conn
}

func releaseConnection(conn *Connection) {
 conn.Addr = "" // 清理状态
 connectionPool.Put(conn)
}

// 带初始化的Pool
type Worker struct {
 ID     int
 active bool
}

var workerPool = sync.Pool{
 New: func() interface{} {
  return &Worker{ID: 0, active: false}
 },
}

func (w *Worker) Initialize(id int) {
 w.ID = id
 w.active = true
}

func (w *Worker) Reset() {
 w.active = false
}

// 性能对比
func benchmarkPool() {
 const n = 1000000

 // 不使用Pool
 start := time.Now()
 for i := 0; i < n; i++ {
  _ = new(bytes.Buffer)
 }
 fmt.Printf("Without Pool: %v\n", time.Since(start))

 // 使用Pool
 start = time.Now()
 for i := 0; i < n; i++ {
  buf := bufferPool.Get().(*bytes.Buffer)
  bufferPool.Put(buf)
 }
 fmt.Printf("With Pool: %v\n", time.Since(start))
}

// 并发使用Pool
func concurrentPool() {
 var wg sync.WaitGroup

 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()

   buf := bufferPool.Get().(*bytes.Buffer)
   defer bufferPool.Put(buf)

   buf.Reset()
   fmt.Fprintf(buf, "Worker %d", id)
  }(i)
 }

 wg.Wait()
}

// Pool的局限性演示
func poolLimitation() {
 // Pool中的对象可能在任何时候被GC
 pool := sync.Pool{
  New: func() interface{} {
   fmt.Println("Creating object")
   return &struct{}{}
  },
 }

 obj := pool.Get()
 pool.Put(obj)

 // 强制GC
 fmt.Println("Before GC:", pool.Get())
 pool.Put(obj)

 runtime.GC()
 time.Sleep(100 * time.Millisecond)

 // 可能创建新对象
 fmt.Println("After GC:", pool.Get())
}

import "runtime"

func main() {
 fmt.Println("=== Basic Pool ===")
 for i := 0; i < 5; i++ {
  useBuffer()
 }

 fmt.Println("\n=== Connection Pool ===")
 conn := getConnection()
 conn.Addr = "localhost:8080"
 fmt.Printf("Connection: ID=%d, Addr=%s\n", conn.ID, conn.Addr)
 releaseConnection(conn)

 conn2 := getConnection()
 fmt.Printf("Reused Connection: ID=%d, Addr=%s\n", conn2.ID, conn2.Addr)
 releaseConnection(conn2)

 fmt.Println("\n=== Benchmark ===")
 benchmarkPool()

 fmt.Println("\n=== Pool Limitation ===")
 poolLimitation()
}
```

#### 反例说明

```go
// ❌ 错误：放回Pool后继续使用
func useAfterPut() {
 pool := sync.Pool{
  New: func() interface{} {
   return &bytes.Buffer{}
  },
 }

 buf := pool.Get().(*bytes.Buffer)
 pool.Put(buf)

 buf.WriteString("after put") // 危险：对象可能已被其他goroutine获取
}

// ❌ 错误：Pool中存储有状态对象
func statefulPool() {
 type Stateful struct {
  Data []int
 }

 pool := sync.Pool{
  New: func() interface{} {
   return &Stateful{Data: make([]int, 0, 100)}
  },
 }

 s := pool.Get().(*Stateful)
 s.Data = append(s.Data, 1, 2, 3)
 pool.Put(s)

 s2 := pool.Get().(*Stateful)
 // s2.Data可能包含之前的数据！
}

// ❌ 错误：假设Pool一定返回对象
func assumeObject() {
 pool := sync.Pool{}
 // 没有New函数

 obj := pool.Get()
 if obj == nil {
  fmt.Println("Got nil")
 }
}

// ✅ 正确：重置对象状态
func correctReset() {
 pool := sync.Pool{
  New: func() interface{} {
   return &bytes.Buffer{}
  },
 }

 buf := pool.Get().(*bytes.Buffer)
 buf.Reset() // 重置状态
 buf.WriteString("data")
 // 使用...
 pool.Put(buf)
}
```

#### 最佳实践

1. **在Put前重置对象状态**
2. **不要假设Get一定返回非nil**
3. **Pool适合临时对象，不适合持久对象**
4. **不要存储需要确定释放的资源**
5. **Pool中的对象可能在GC时被释放**

---

### 3.6 Map

#### 概念定义

**sync.Map**是并发安全的map实现，针对读多写少场景优化，避免使用Mutex的锁竞争。

#### 工作原理

```
sync.Map结构:

┌─────────────────────────────────────────┐
│              sync.Map                    │
├─────────────────────────────────────────┤
│  type Map struct {                      │
│      mu Mutex                           │
│      read atomic.Value  // read-only    │
│      dirty map[interface{}]*entry       │
│      misses int                         │
│  }                                      │
│                                         │
│  读路径：                                │
│    1. 先查read（无锁）                   │
│    2. 未命中且dirty有新数据，加锁查dirty │
│    3. misses++，达到阈值时提升dirty      │
│                                         │
│  写路径：                                │
│    1. 如果read中有，CAS更新              │
│    2. 否则加锁，写入dirty                │
└─────────────────────────────────────────┘
```

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 基本sync.Map
func basicSyncMap() {
 var m sync.Map

 // 存储
 m.Store("key1", "value1")
 m.Store("key2", 123)
 m.Store("key3", 3.14)

 // 读取
 if val, ok := m.Load("key1"); ok {
  fmt.Printf("key1: %v\n", val)
 }

 // 不存在则存储
 if val, loaded := m.LoadOrStore("key4", "new"); !loaded {
  fmt.Printf("Stored new value: %v\n", val)
 }

 // 删除
 m.Delete("key2")

 // 遍历
 fmt.Println("All entries:")
 m.Range(func(key, value interface{}) bool {
  fmt.Printf("  %v: %v\n", key, value)
  return true // 继续遍历
 })
}

// 对比sync.Map和map+Mutex
func benchmarkMap() {
 const n = 100000

 // sync.Map
 syncMap := &sync.Map{}
 start := time.Now()
 var wg sync.WaitGroup

 // 写入
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func(i int) {
   defer wg.Done()
   syncMap.Store(i, i)
  }(i)
 }
 wg.Wait()

 // 读取
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func(i int) {
   defer wg.Done()
   syncMap.Load(i)
  }(i)
 }
 wg.Wait()

 fmt.Printf("sync.Map: %v\n", time.Since(start))

 // map + RWMutex
 type SafeMap struct {
  mu sync.RWMutex
  m  map[int]int
 }

 safeMap := &SafeMap{m: make(map[int]int)}
 start = time.Now()

 // 写入
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func(i int) {
   defer wg.Done()
   safeMap.mu.Lock()
   safeMap.m[i] = i
   safeMap.mu.Unlock()
  }(i)
 }
 wg.Wait()

 // 读取
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func(i int) {
   defer wg.Done()
   safeMap.mu.RLock()
   _ = safeMap.m[i]
   safeMap.mu.RUnlock()
  }(i)
 }
 wg.Wait()

 fmt.Printf("map+RWMutex: %v\n", time.Since(start))
}

// 缓存实现
type Cache struct {
 data sync.Map
}

func (c *Cache) Get(key string) (interface{}, bool) {
 return c.data.Load(key)
}

func (c *Cache) Set(key string, value interface{}) {
 c.data.Store(key, value)
}

func (c *Cache) Delete(key string) {
 c.data.Delete(key)
}

// 带过期时间的缓存
type ExpiringCache struct {
 data sync.Map
}

type expiringValue struct {
 value      interface{}
 expiration time.Time
}

func (c *ExpiringCache) Set(key string, value interface{}, ttl time.Duration) {
 c.data.Store(key, expiringValue{
  value:      value,
  expiration: time.Now().Add(ttl),
 })
}

func (c *ExpiringCache) Get(key string) (interface{}, bool) {
 if val, ok := c.data.Load(key); ok {
  ev := val.(expiringValue)
  if time.Now().Before(ev.expiration) {
   return ev.value, true
  }
  // 过期，删除
  c.data.Delete(key)
 }
 return nil, false
}

// 计数器
type Counter struct {
 counts sync.Map
}

func (c *Counter) Increment(key string) int {
 actual, _ := c.counts.LoadOrStore(key, 0)
 for {
  current := actual.(int)
  next := current + 1
  if c.counts.CompareAndSwap(key, current, next) {
   return next
  }
  actual, _ = c.counts.Load(key)
 }
}

func (c *Counter) Get(key string) int {
 if val, ok := c.counts.Load(key); ok {
  return val.(int)
 }
 return 0
}

func main() {
 fmt.Println("=== Basic sync.Map ===")
 basicSyncMap()

 fmt.Println("\n=== Cache Implementation ===")
 cache := &Cache{}
 cache.Set("user:1", "Alice")
 cache.Set("user:2", "Bob")

 if val, ok := cache.Get("user:1"); ok {
  fmt.Printf("User 1: %v\n", val)
 }

 fmt.Println("\n=== Expiring Cache ===")
 ec := &ExpiringCache{}
 ec.Set("key1", "value1", 100*time.Millisecond)

 if val, ok := ec.Get("key1"); ok {
  fmt.Printf("Before expire: %v\n", val)
 }

 time.Sleep(150 * time.Millisecond)

 if _, ok := ec.Get("key1"); !ok {
  fmt.Println("After expire: not found")
 }

 fmt.Println("\n=== Counter ===")
 counter := &Counter{}
 var wg sync.WaitGroup
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter.Increment("requests")
  }()
 }
 wg.Wait()
 fmt.Printf("Total requests: %d\n", counter.Get("requests"))
}
```

#### 反例说明

```go
// ❌ 错误：使用复杂类型作为key
func complexKey() {
 var m sync.Map

 type Key struct {
  Data []int // 包含slice，不能作为map key
 }

 // m.Store(Key{Data: []int{1}}, "value") // 可以存储，但不能正确比较
}

// ❌ 错误：在Range中修改Map
func modifyInRange() {
 var m sync.Map
 m.Store("key1", "value1")
 m.Store("key2", "value2")

 m.Range(func(k, v interface{}) bool {
  m.Delete(k) // 危险：可能导致未定义行为
  return true
 })
}

// ❌ 错误：假设Load的顺序
func assumeOrder() {
 var m sync.Map
 m.Store("a", 1)
 m.Store("b", 2)
 m.Store("c", 3)

 // Range不保证顺序
 m.Range(func(k, v interface{}) bool {
  fmt.Println(k, v)
  return true
 })
}

// ✅ 正确：使用简单类型作为key
func correctKey() {
 var m sync.Map
 m.Store("string_key", "value")
 m.Store(123, "int_key")
 m.Store(struct{ A int }{A: 1}, "struct_key")
}
```

#### 最佳实践

1. **读多写少场景使用sync.Map**
2. **使用简单类型作为key**
3. **不要在Range中修改Map**
4. **Range不保证顺序**
5. **普通map+Mutex在写多场景可能更好**

---

### 3.7 Atomic操作

#### 概念定义

**Atomic操作**是不可分割的原子操作，用于无锁并发编程，提供最高性能的基本同步原语。

#### 工作原理

```
Atomic操作类型:

┌─────────────────────────────────────────┐
│           sync/atomic                    │
├─────────────────────────────────────────┤
│  基本操作:                               │
│  - AddInt32/64: 原子加法                 │
│  - CompareAndSwapInt32/64: CAS           │
│  - LoadInt32/64: 原子读取                │
│  - StoreInt32/64: 原子写入               │
│  - SwapInt32/64: 原子交换                │
│                                         │
│  指针操作:                               │
│  - LoadPointer, StorePointer             │
│  - CompareAndSwapPointer                 │
│                                         │
│  Value类型:                              │
│  - Load, Store, CompareAndSwap           │
└─────────────────────────────────────────┘

CAS循环模式:
for {
    old := atomic.LoadInt32(&value)
    new := old + 1
    if atomic.CompareAndSwapInt32(&value, old, new) {
        break
    }
}
```

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "sync/atomic"
 "time"
)

// 基本Atomic操作
func basicAtomic() {
 var counter int64 = 0

 // 原子加法
 atomic.AddInt64(&counter, 1)
 fmt.Printf("After Add: %d\n", counter)

 // 原子读取
 val := atomic.LoadInt64(&counter)
 fmt.Printf("Load: %d\n", val)

 // 原子写入
 atomic.StoreInt64(&counter, 100)
 fmt.Printf("After Store: %d\n", counter)

 // 原子交换
 old := atomic.SwapInt64(&counter, 200)
 fmt.Printf("Swap: old=%d, new=%d\n", old, counter)

 // CAS
 if atomic.CompareAndSwapInt64(&counter, 200, 300) {
  fmt.Printf("CAS success: %d\n", counter)
 }
}

// 无锁计数器
type AtomicCounter struct {
 value int64
}

func (c *AtomicCounter) Increment() int64 {
 return atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Decrement() int64 {
 return atomic.AddInt64(&c.value, -1)
}

func (c *AtomicCounter) Value() int64 {
 return atomic.LoadInt64(&c.value)
}

// 无锁栈（Treiber Stack）
type Node struct {
 value int
 next  *Node
}

type LockFreeStack struct {
 top unsafe.Pointer
}

func (s *LockFreeStack) Push(value int) {
 newNode := &Node{value: value}
 for {
  oldTop := atomic.LoadPointer(&s.top)
  newNode.next = (*Node)(oldTop)
  if atomic.CompareAndSwapPointer(&s.top, oldTop, unsafe.Pointer(newNode)) {
   return
  }
 }
}

func (s *LockFreeStack) Pop() (int, bool) {
 for {
  oldTop := atomic.LoadPointer(&s.top)
  if oldTop == nil {
   return 0, false
  }
  node := (*Node)(oldTop)
  newTop := unsafe.Pointer(node.next)
  if atomic.CompareAndSwapPointer(&s.top, oldTop, newTop) {
   return node.value, true
  }
 }
}

import "unsafe"

// 原子Value
type Config struct {
 Version int
 Data    string
}

var configValue atomic.Value

func loadConfig() *Config {
 return configValue.Load().(*Config)
}

func storeConfig(cfg *Config) {
 configValue.Store(cfg)
}

// 性能对比
func benchmarkAtomic() {
 const n = 10000000

 // Atomic
 var atomicCounter int64
 start := time.Now()
 var wg sync.WaitGroup
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   atomic.AddInt64(&atomicCounter, 1)
  }()
 }
 wg.Wait()
 fmt.Printf("Atomic: %v, value=%d\n", time.Since(start), atomicCounter)

 // Mutex
 var mu sync.Mutex
 var mutexCounter int64
 start = time.Now()
 for i := 0; i < n; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   mu.Lock()
   mutexCounter++
   mu.Unlock()
  }()
 }
 wg.Wait()
 fmt.Printf("Mutex: %v, value=%d\n", time.Since(start), mutexCounter)
}

// 自旋锁
 type SpinLock struct {
 flag int32
}

func (s *SpinLock) Lock() {
 for !atomic.CompareAndSwapInt32(&s.flag, 0, 1) {
  // 自旋等待
  runtime.Gosched()
 }
}

func (s *SpinLock) Unlock() {
 atomic.StoreInt32(&s.flag, 0)
}

import "runtime"

// 序列号生成器
type SequenceGenerator struct {
 sequence int64
}

func (s *SequenceGenerator) Next() int64 {
 return atomic.AddInt64(&s.sequence, 1)
}

// 引用计数
type RefCounted struct {
 refCount int32
 data     interface{}
}

func (r *RefCounted) AddRef() int32 {
 return atomic.AddInt32(&r.refCount, 1)
}

func (r *RefCounted) Release() int32 {
 newCount := atomic.AddInt32(&r.refCount, -1)
 if newCount == 0 {
  // 释放资源
  fmt.Println("Resource released")
 }
 return newCount
}

func main() {
 fmt.Println("=== Basic Atomic ===")
 basicAtomic()

 fmt.Println("\n=== Atomic Counter ===")
 counter := &AtomicCounter{}
 var wg sync.WaitGroup
 for i := 0; i < 1000; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter.Increment()
  }()
 }
 wg.Wait()
 fmt.Printf("Counter: %d\n", counter.Value())

 fmt.Println("\n=== Atomic Value ===")
 storeConfig(&Config{Version: 1, Data: "initial"})
 cfg := loadConfig()
 fmt.Printf("Config: %+v\n", cfg)

 storeConfig(&Config{Version: 2, Data: "updated"})
 cfg = loadConfig()
 fmt.Printf("Updated Config: %+v\n", cfg)

 fmt.Println("\n=== Sequence Generator ===")
 seq := &SequenceGenerator{}
 for i := 0; i < 5; i++ {
  fmt.Printf("Sequence: %d\n", seq.Next())
 }

 fmt.Println("\n=== Ref Counted ===")
 ref := &RefCounted{refCount: 1, data: "resource"}
 ref.AddRef()
 ref.AddRef()
 ref.Release()
 ref.Release()
 ref.Release()
}
```

#### 反例说明

```go
// ❌ 错误：非原子操作
func nonAtomic() {
 var counter int64

 // 非原子：读取+写入是两个操作
 counter = counter + 1 // 竞态条件！
}

// ❌ 错误：ABA问题
func abaProblem() {
 var ptr unsafe.Pointer

 // 假设ptr指向A
 old := atomic.LoadPointer(&ptr)

 // 其他goroutine将ptr改为B，再改回A

 // CAS成功，但ptr已被修改过
 atomic.CompareAndSwapPointer(&ptr, old, newPtr)
}

// ❌ 错误：64位对齐问题（32位系统）
func alignmentIssue() {
 struct {
  _     int32 // 填充
  value int64 // 可能不对齐
 }
 // 在32位系统上，int64可能不对齐，导致panic
}

// ✅ 正确：使用CAS循环
func correctCASLoop() {
 var value int64 = 0

 for {
  old := atomic.LoadInt64(&value)
  new := old + 1
  if atomic.CompareAndSwapInt64(&value, old, new) {
   break
  }
 }
}
```

#### 性能分析

| 操作 | Atomic | Mutex |
|------|--------|-------|
| 计数器 | ~10ns | ~100ns |
| 读 | ~5ns | ~50ns |
| 写 | ~10ns | ~100ns |

#### 最佳实践

1. **简单计数器使用Atomic**
2. **复杂操作使用Mutex**
3. **注意64位对齐（32位系统）**
4. **CAS循环要有退出条件**
5. **警惕ABA问题**

---

## 4. Context模式

### 4.1 取消传播

#### 概念定义

**Context**是Go语言中用于传递截止时间、取消信号和请求范围值的API，实现跨Goroutine的取消传播。

#### 工作原理

```
Context结构:

┌─────────────────────────────────────────┐
│              Context Tree                │
├─────────────────────────────────────────┤
│                                         │
│              Background()               │
│                   │                     │
│         ┌────────┼────────┐             │
│         ↓        ↓        ↓             │
│      ctx1     ctx2     ctx3            │
│      (TO)     (CA)     (VAL)           │
│         │                              │
│    ┌────┴────┐                         │
│    ↓         ↓                         │
│  ctx4     ctx5                         │
│                                         │
│  TO: WithTimeout                        │
│  CA: WithCancel                         │
│  VAL: WithValue                         │
│                                         │
│  取消传播: 父取消 → 所有子取消           │
└─────────────────────────────────────────┘

取消机制:
parent.Done() ──close──→ child.Done() ──close──→ grandchild.Done()
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "time"
)

// 基本取消
func basicCancel() {
 ctx, cancel := context.WithCancel(context.Background())

 go func() {
  for {
   select {
   case <-ctx.Done():
    fmt.Println("Worker cancelled:", ctx.Err())
    return
   default:
    fmt.Println("Working...")
    time.Sleep(200 * time.Millisecond)
   }
  }
 }()

 time.Sleep(500 * time.Millisecond)
 fmt.Println("Cancelling...")
 cancel()

 time.Sleep(200 * time.Millisecond)
}

// 取消传播
func cancelPropagation() {
 rootCtx, rootCancel := context.WithCancel(context.Background())

 // 第一层
 ctx1, _ := context.WithCancel(rootCtx)
 ctx2, _ := context.WithCancel(rootCtx)

 // 第二层
 ctx1a, _ := context.WithCancel(ctx1)
 ctx1b, _ := context.WithCancel(ctx1)

 // 启动多个goroutine
 for i, ctx := range []context.Context{ctx1, ctx2, ctx1a, ctx1b} {
  go func(id int, c context.Context) {
   <-c.Done()
   fmt.Printf("Context %d cancelled: %v\n", id, c.Err())
  }(i, ctx)
 }

 time.Sleep(200 * time.Millisecond)
 fmt.Println("Cancelling root...")
 rootCancel()

 time.Sleep(200 * time.Millisecond)
}

// 取消函数传递
func cancelFunctionPassing() {
 worker := func(ctx context.Context, id int) {
  for {
   select {
   case <-ctx.Done():
    fmt.Printf("Worker %d stopped\n", id)
    return
   default:
    fmt.Printf("Worker %d working\n", id)
    time.Sleep(300 * time.Millisecond)
   }
  }
 }

 ctx, cancel := context.WithCancel(context.Background())

 for i := 0; i < 3; i++ {
  go worker(ctx, i)
 }

 time.Sleep(500 * time.Millisecond)
 cancel() // 取消所有worker

 time.Sleep(200 * time.Millisecond)
}

// 选择性取消
func selectiveCancel() {
 ctx1, cancel1 := context.WithCancel(context.Background())
 ctx2, cancel2 := context.WithCancel(context.Background())

 go func() {
  <-ctx1.Done()
  fmt.Println("Context 1 cancelled")
 }()

 go func() {
  <-ctx2.Done()
  fmt.Println("Context 2 cancelled")
 }()

 time.Sleep(200 * time.Millisecond)
 cancel1() // 只取消ctx1

 time.Sleep(200 * time.Millisecond)
 cancel2() // 取消ctx2
}

// 取消后清理
func cleanupOnCancel() {
 ctx, cancel := context.WithCancel(context.Background())

 go func() {
  defer func() {
   fmt.Println("Cleanup: releasing resources")
  }()

  for {
   select {
   case <-ctx.Done():
    fmt.Println("Cancelled, cleaning up...")
    return
   default:
    fmt.Println("Processing...")
    time.Sleep(200 * time.Millisecond)
   }
  }
 }()

 time.Sleep(500 * time.Millisecond)
 cancel()

 time.Sleep(200 * time.Millisecond)
}

// HTTP请求取消示例
func httpRequestCancel() {
 ctx, cancel := context.WithCancel(context.Background())

 // 模拟HTTP请求
 go func() {
  req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
  client := &http.Client{Timeout: 5 * time.Second}

  fmt.Println("Starting request...")
  resp, err := client.Do(req)
  if err != nil {
   fmt.Printf("Request failed: %v\n", err)
   return
  }
  defer resp.Body.Close()
  fmt.Println("Request completed")
 }()

 time.Sleep(100 * time.Millisecond)
 fmt.Println("Cancelling request...")
 cancel()

 time.Sleep(200 * time.Millisecond)
}

import "net/http"

func main() {
 fmt.Println("=== Basic Cancel ===")
 basicCancel()

 fmt.Println("\n=== Cancel Propagation ===")
 cancelPropagation()

 fmt.Println("\n=== Selective Cancel ===")
 selectiveCancel()

 fmt.Println("\n=== Cleanup on Cancel ===")
 cleanupOnCancel()
}
```

#### 反例说明

```go
// ❌ 错误：存储cancel函数
func storeCancel() {
 var cancelFunc context.CancelFunc

 func() {
  _, cancel := context.WithCancel(context.Background())
  cancelFunc = cancel // 可能泄漏
 }()

 // cancelFunc可能在不应该的时候被调用
}

// ❌ 错误：在子goroutine中取消父context
func childCancelParent() {
 parent, _ := context.WithCancel(context.Background())
 child, childCancel := context.WithCancel(parent)

 go func() {
  childCancel() // 正确：取消子context
  // parentCancel() // 错误：不应该取消父context
 }()

 _ = child
}

// ❌ 错误：多次取消
func multipleCancel() {
 _, cancel := context.WithCancel(context.Background())

 cancel()
 cancel() // 安全但无意义
 cancel() // 同上
}

// ✅ 正确：及时调用cancel
func correctCancel() {
 ctx, cancel := context.WithCancel(context.Background())
 defer cancel() // 确保取消

 // 使用ctx
 _ = ctx
}
```

#### 最佳实践

1. **将context作为第一个参数传递**
2. **及时调用cancel函数（defer cancel()）**
3. **不要存储cancel函数**
4. **只取消自己的子context**
5. **在API边界检查ctx.Done()**

---

### 4.2 超时控制

#### 概念定义

**超时控制**通过`WithTimeout`或`WithDeadline`为操作设置时间限制，超时后自动取消context。

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "time"
)

// 基本超时
func basicTimeout() {
 ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
 defer cancel()

 select {
 case <-time.After(2 * time.Second):
  fmt.Println("Operation completed")
 case <-ctx.Done():
  fmt.Println("Timeout:", ctx.Err())
 }
}

// 截止时间
func deadlineExample() {
 deadline := time.Now().Add(1 * time.Second)
 ctx, cancel := context.WithDeadline(context.Background(), deadline)
 defer cancel()

 fmt.Println("Deadline:", ctx.Err()) // nil

 <-ctx.Done()
 fmt.Println("After deadline:", ctx.Err()) // context deadline exceeded
}

// 嵌套超时
func nestedTimeout() {
 outerCtx, outerCancel := context.WithTimeout(context.Background(), 3*time.Second)
 defer outerCancel()

 innerCtx, innerCancel := context.WithTimeout(outerCtx, 1*time.Second)
 defer innerCancel()

 select {
 case <-innerCtx.Done():
  fmt.Println("Inner timeout:", innerCtx.Err())
 case <-outerCtx.Done():
  fmt.Println("Outer timeout:", outerCtx.Err())
 }
}

// 带超时的函数调用
func callWithTimeout(fn func() error, timeout time.Duration) error {
 ctx, cancel := context.WithTimeout(context.Background(), timeout)
 defer cancel()

 done := make(chan error, 1)
 go func() {
  done <- fn()
 }()

 select {
 case err := <-done:
  return err
 case <-ctx.Done():
  return ctx.Err()
 }
}

// 超时后重试
func retryWithTimeout() {
 operation := func() error {
  time.Sleep(500 * time.Millisecond)
  return fmt.Errorf("operation failed")
 }

 maxRetries := 3
 for i := 0; i < maxRetries; i++ {
  err := callWithTimeout(operation, 200*time.Millisecond)
  if err == nil {
   fmt.Println("Success")
   return
  }
  fmt.Printf("Attempt %d failed: %v\n", i+1, err)
 }
 fmt.Println("All retries failed")
}

// 动态超时调整
func dynamicTimeout() {
 baseTimeout := 100 * time.Millisecond

 for i := 0; i < 5; i++ {
  // 逐渐增加超时
  timeout := baseTimeout * time.Duration(i+1)

  ctx, cancel := context.WithTimeout(context.Background(), timeout)

  start := time.Now()
  select {
  case <-time.After(300 * time.Millisecond):
   fmt.Printf("Attempt %d succeeded with timeout %v\n", i+1, timeout)
   cancel()
   return
  case <-ctx.Done():
   fmt.Printf("Attempt %d timeout after %v\n", i+1, time.Since(start))
   cancel()
  }
 }
}

// 超时统计
func timeoutStats() {
 timeouts := 0
 successes := 0

 for i := 0; i < 10; i++ {
  ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

  select {
  case <-time.After(time.Duration(i*10) * time.Millisecond):
   successes++
  case <-ctx.Done():
   timeouts++
  }

  cancel()
 }

 fmt.Printf("Successes: %d, Timeouts: %d\n", successes, timeouts)
}

func main() {
 fmt.Println("=== Basic Timeout ===")
 basicTimeout()

 fmt.Println("\n=== Deadline Example ===")
 deadlineExample()

 fmt.Println("\n=== Nested Timeout ===")
 nestedTimeout()

 fmt.Println("\n=== Retry with Timeout ===")
 retryWithTimeout()

 fmt.Println("\n=== Dynamic Timeout ===")
 dynamicTimeout()

 fmt.Println("\n=== Timeout Stats ===")
 timeoutStats()
}
```

#### 最佳实践

1. **总是使用defer cancel()**
2. **超时时间应该是可配置的**
3. **嵌套超时取最小值**
4. **超时后记录日志**
5. **考虑重试策略**

---

### 4.3 值传递

#### 概念定义

**值传递**通过`WithValue`在context树中存储键值对，用于传递请求范围的元数据（如trace ID、用户信息等）。

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
)

// 定义key类型（避免冲突）
type contextKey string

const (
 RequestIDKey contextKey = "request_id"
 UserIDKey    contextKey = "user_id"
 TraceIDKey   contextKey = "trace_id"
)

// 基本值传递
func basicValue() {
 ctx := context.Background()

 // 添加值
 ctx = context.WithValue(ctx, RequestIDKey, "req-123")
 ctx = context.WithValue(ctx, UserIDKey, "user-456")

 // 获取值
 if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
  fmt.Println("Request ID:", reqID)
 }

 if userID, ok := ctx.Value(UserIDKey).(string); ok {
  fmt.Println("User ID:", userID)
 }
}

// 值传播
func valuePropagation() {
 root := context.Background()

 // 父context添加值
 parent := context.WithValue(root, "key1", "value1")

 // 子context继承值
 child := context.WithValue(parent, "key2", "value2")

 // 孙子context
 grandchild := context.WithValue(child, "key3", "value3")

 // 可以访问所有祖先的值
 fmt.Println("key1:", grandchild.Value("key1"))
 fmt.Println("key2:", grandchild.Value("key2"))
 fmt.Println("key3:", grandchild.Value("key3"))
}

// 中间件模式
func middlewarePattern() {
 type Handler func(ctx context.Context) error

 // 日志中间件
 loggingMiddleware := func(next Handler) Handler {
  return func(ctx context.Context) error {
   if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
    fmt.Printf("[%s] Start processing\n", reqID)
   }
   err := next(ctx)
   if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
    fmt.Printf("[%s] End processing\n", reqID)
   }
   return err
  }
 }

 // 认证中间件
 authMiddleware := func(next Handler) Handler {
  return func(ctx context.Context) error {
   if userID, ok := ctx.Value(UserIDKey).(string); ok {
    fmt.Printf("[Auth] User %s authenticated\n", userID)
   }
   return next(ctx)
  }
 }

 // 业务处理
 handler := func(ctx context.Context) error {
  fmt.Println("Business logic executed")
  return nil
 }

 // 组合中间件
 wrapped := loggingMiddleware(authMiddleware(handler))

 // 创建context并执行
 ctx := context.WithValue(context.Background(), RequestIDKey, "req-789")
 ctx = context.WithValue(ctx, UserIDKey, "user-101")

 wrapped(ctx)
}

// 值覆盖
func valueOverride() {
 ctx := context.Background()

 ctx = context.WithValue(ctx, "key", "original")
 fmt.Println("Original:", ctx.Value("key"))

 // 子context覆盖值
 child := context.WithValue(ctx, "key", "overridden")
 fmt.Println("Overridden:", child.Value("key"))

 // 父context不变
 fmt.Println("Parent:", ctx.Value("key"))
}

// 辅助函数
type ContextUtils struct{}

func (c ContextUtils) GetRequestID(ctx context.Context) string {
 if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
  return reqID
 }
 return ""
}

func (c ContextUtils) GetUserID(ctx context.Context) string {
 if userID, ok := ctx.Value(UserIDKey).(string); ok {
  return userID
 }
 return ""
}

func (c ContextUtils) WithRequestID(ctx context.Context, reqID string) context.Context {
 return context.WithValue(ctx, RequestIDKey, reqID)
}

func (c ContextUtils) WithUserID(ctx context.Context, userID string) context.Context {
 return context.WithValue(ctx, UserIDKey, userID)
}

// 使用辅助函数
func useHelpers() {
 utils := ContextUtils{}

 ctx := context.Background()
 ctx = utils.WithRequestID(ctx, "req-abc")
 ctx = utils.WithUserID(ctx, "user-xyz")

 fmt.Println("Request ID:", utils.GetRequestID(ctx))
 fmt.Println("User ID:", utils.GetUserID(ctx))
}

func main() {
 fmt.Println("=== Basic Value ===")
 basicValue()

 fmt.Println("\n=== Value Propagation ===")
 valuePropagation()

 fmt.Println("\n=== Middleware Pattern ===")
 middlewarePattern()

 fmt.Println("\n=== Value Override ===")
 valueOverride()

 fmt.Println("\n=== Use Helpers ===")
 useHelpers()
}
```

#### 反例说明

```go
// ❌ 错误：使用string作为key
func stringKey() {
 ctx := context.WithValue(context.Background(), "key", "value")
 // 可能被其他包覆盖
}

// ❌ 错误：存储敏感信息
func sensitiveInfo() {
 ctx := context.WithValue(context.Background(), "password", "secret")
 // 敏感信息可能泄露到日志
}

// ❌ 错误：存储大对象
func largeObject() {
 largeData := make([]byte, 1000000)
 ctx := context.WithValue(context.Background(), "data", largeData)
 // 大对象会随context传递，影响性能
 _ = ctx
}

// ❌ 错误：频繁修改context
func frequentModify() {
 ctx := context.Background()
 for i := 0; i < 1000; i++ {
  ctx = context.WithValue(ctx, i, i) // 创建大量context
 }
}

// ✅ 正确：使用自定义类型作为key
func correctKeyType() {
 type myKey struct{}
 ctx := context.WithValue(context.Background(), myKey{}, "value")
 _ = ctx
}
```

#### 最佳实践

1. **使用自定义类型作为key**
2. **只存储元数据，不存储业务数据**
3. **不要存储敏感信息**
4. **提供辅助函数访问值**
5. **文档化context中存储的值**

---

### 4.4 Context链

#### 概念定义

**Context链**是通过嵌套创建context形成的树形结构，实现取消、超时、值的层级传播。

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "time"
)

// 构建context链
func buildContextChain() context.Context {
 // 根context
 root := context.Background()
 fmt.Println("1. Root context created")

 // 添加值
 ctx := context.WithValue(root, "app", "myapp")
 fmt.Println("2. Added app value")

 // 添加超时
 ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
 defer cancel()
 fmt.Println("3. Added 5s timeout")

 // 添加更多值
 ctx = context.WithValue(ctx, "request_id", "req-123")
 fmt.Println("4. Added request_id")

 // 添加子超时
 ctx, subCancel := context.WithTimeout(ctx, 2*time.Second)
 defer subCancel()
 fmt.Println("5. Added 2s sub-timeout")

 return ctx
}

// 遍历context链
func traverseContextChain(ctx context.Context) {
 fmt.Println("\nContext chain:")
 for ctx != nil {
  fmt.Printf("  - Type: %T\n", ctx)

  // 检查取消状态
  select {
  case <-ctx.Done():
   fmt.Printf("    Done: %v\n", ctx.Err())
  default:
   fmt.Println("    Not done")
  }

  // 检查截止时间
  if deadline, ok := ctx.Deadline(); ok {
   fmt.Printf("    Deadline: %v\n", deadline)
  }

  ctx = ctx.Value(struct{}{}) // 获取父context（hack方式）
 }
}

// 实际应用：HTTP请求处理链
func httpRequestChain() {
 // 请求级别context
 reqCtx, reqCancel := context.WithTimeout(context.Background(), 30*time.Second)
 defer reqCancel()
 reqCtx = context.WithValue(reqCtx, "request_id", "req-abc-123")

 // 数据库查询级别（更短的超时）
 dbCtx, dbCancel := context.WithTimeout(reqCtx, 5*time.Second)
 defer dbCancel()

 // 缓存查询级别
 cacheCtx, cacheCancel := context.WithTimeout(reqCtx, 100*time.Millisecond)
 defer cacheCancel()

 // 执行操作
 go func() {
  select {
  case <-dbCtx.Done():
   fmt.Println("DB query cancelled:", dbCtx.Err())
  case <-time.After(3 * time.Second):
   fmt.Println("DB query completed")
  }
 }()

 go func() {
  select {
  case <-cacheCtx.Done():
   fmt.Println("Cache query cancelled:", cacheCtx.Err())
  case <-time.After(50 * time.Millisecond):
   fmt.Println("Cache query completed")
  }
 }()

 time.Sleep(4 * time.Second)
}

// 取消传播演示
func cancelPropagationDemo() {
 root, rootCancel := context.WithCancel(context.Background())

 // 创建多层context
 level1, _ := context.WithCancel(root)
 level2a, _ := context.WithCancel(level1)
 level2b, _ := context.WithCancel(level1)
 level3, _ := context.WithCancel(level2a)

 contexts := map[string]context.Context{
  "root":    root,
  "level1":  level1,
  "level2a": level2a,
  "level2b": level2b,
  "level3":  level3,
 }

 // 监控所有context
 for name, ctx := range contexts {
  go func(n string, c context.Context) {
   <-c.Done()
   fmt.Printf("%s cancelled\n", n)
  }(name, ctx)
 }

 time.Sleep(100 * time.Millisecond)
 fmt.Println("\nCancelling root...")
 rootCancel()

 time.Sleep(100 * time.Millisecond)
}

// Context链最佳实践
type Service struct {
 db     *Database
 cache  *Cache
 logger *Logger
}

type Database struct{}
type Cache struct{}
type Logger struct{}

func (s *Service) HandleRequest(ctx context.Context, reqID string) error {
 // 添加请求ID
 ctx = context.WithValue(ctx, "request_id", reqID)

 // 尝试从缓存获取
 cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
 defer cancel()

 if data, err := s.cache.Get(cacheCtx, reqID); err == nil {
  s.logger.Log(ctx, "Cache hit")
  return s.processData(ctx, data)
 }

 // 从数据库获取
 dbCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
 defer cancel()

 data, err := s.db.Query(dbCtx, reqID)
 if err != nil {
  return err
 }

 // 更新缓存（后台）
 go s.cache.Set(context.WithoutCancel(ctx), reqID, data)

 return s.processData(ctx, data)
}

func (s *Service) processData(ctx context.Context, data interface{}) error {
 // 处理数据
 return nil
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
 return nil, fmt.Errorf("not found")
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) {
}

func (db *Database) Query(ctx context.Context, query string) (interface{}, error) {
 return nil, nil
}

func (l *Logger) Log(ctx context.Context, msg string) {
 if reqID, ok := ctx.Value("request_id").(string); ok {
  fmt.Printf("[%s] %s\n", reqID, msg)
 }
}

func main() {
 fmt.Println("=== Build Context Chain ===")
 ctx := buildContextChain()
 _ = ctx

 fmt.Println("\n=== HTTP Request Chain ===")
 httpRequestChain()

 fmt.Println("\n=== Cancel Propagation ===")
 cancelPropagationDemo()

 fmt.Println("\n=== Service Example ===")
 service := &Service{}
 reqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()
 service.HandleRequest(reqCtx, "req-test-001")
}
```

#### 最佳实践

1. **Context作为函数第一个参数**
2. **不要存储context，只传递**
3. **及时取消子context**
4. **使用context.WithoutCancel进行后台任务**
5. **保持context链的清晰**

---

## 5. 常见并发模式

### 5.1 Worker Pool（工作池）

#### 概念定义

**Worker Pool**模式维护一组固定数量的工作Goroutine，从任务队列中获取并执行任务，限制并发数量。

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// 基本Worker Pool
type WorkerPool struct {
 workers   int
 taskQueue chan func()
 wg        sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
 return &WorkerPool{
  workers:   workers,
  taskQueue: make(chan func()),
 }
}

func (p *WorkerPool) Start() {
 for i := 0; i < p.workers; i++ {
  p.wg.Add(1)
  go p.worker(i)
 }
}

func (p *WorkerPool) worker(id int) {
 defer p.wg.Done()
 for task := range p.taskQueue {
  fmt.Printf("Worker %d processing task\n", id)
  task()
 }
}

func (p *WorkerPool) Submit(task func()) {
 p.taskQueue <- task
}

func (p *WorkerPool) Stop() {
 close(p.taskQueue)
 p.wg.Wait()
}

// 带结果的Worker Pool
type Task struct {
 ID   int
 Data int
}

type Result struct {
 TaskID int
 Output int
 Error  error
}

type ResultWorkerPool struct {
 workers   int
 taskQueue chan Task
 results   chan Result
 wg        sync.WaitGroup
}

func NewResultWorkerPool(workers int) *ResultWorkerPool {
 return &ResultWorkerPool{
  workers:   workers,
  taskQueue: make(chan Task),
  results:   make(chan Result, 100),
 }
}

func (p *ResultWorkerPool) Start() {
 for i := 0; i < p.workers; i++ {
  p.wg.Add(1)
  go p.worker(i)
 }

 go func() {
  p.wg.Wait()
  close(p.results)
 }()
}

func (p *ResultWorkerPool) worker(id int) {
 defer p.wg.Done()
 for task := range p.taskQueue {
  output := task.Data * task.Data
  p.results <- Result{
   TaskID: task.ID,
   Output: output,
  }
 }
}

func (p *ResultWorkerPool) Submit(task Task) {
 p.taskQueue <- task
}

func (p *ResultWorkerPool) Results() <-chan Result {
 return p.results
}

func (p *ResultWorkerPool) Stop() {
 close(p.taskQueue)
}

// 带Context的Worker Pool
type ContextWorkerPool struct {
 workers   int
 taskQueue chan func(context.Context)
 ctx       context.Context
 cancel    context.CancelFunc
 wg        sync.WaitGroup
}

func NewContextWorkerPool(workers int) *ContextWorkerPool {
 ctx, cancel := context.WithCancel(context.Background())
 return &ContextWorkerPool{
  workers:   workers,
  taskQueue: make(chan func(context.Context)),
  ctx:       ctx,
  cancel:    cancel,
 }
}

func (p *ContextWorkerPool) Start() {
 for i := 0; i < p.workers; i++ {
  p.wg.Add(1)
  go p.worker(i)
 }
}

func (p *ContextWorkerPool) worker(id int) {
 defer p.wg.Done()
 for {
  select {
  case task, ok := <-p.taskQueue:
   if !ok {
    return
   }
   task(p.ctx)
  case <-p.ctx.Done():
   return
  }
 }
}

func (p *ContextWorkerPool) Submit(task func(context.Context)) {
 select {
 case p.taskQueue <- task:
 case <-p.ctx.Done():
 }
}

func (p *ContextWorkerPool) Stop() {
 p.cancel()
 close(p.taskQueue)
 p.wg.Wait()
}

// 动态Worker Pool
type DynamicWorkerPool struct {
 minWorkers    int
 maxWorkers    int
 currentWorker int32
 taskQueue     chan func()
 wg            sync.WaitGroup
 mu            sync.Mutex
}

func NewDynamicWorkerPool(min, max int) *DynamicWorkerPool {
 return &DynamicWorkerPool{
  minWorkers: min,
  maxWorkers: max,
  taskQueue:  make(chan func()),
 }
}

func (p *DynamicWorkerPool) Start() {
 for i := 0; i < p.minWorkers; i++ {
  p.addWorker()
 }
 go p.monitor()
}

func (p *DynamicWorkerPool) addWorker() {
 p.mu.Lock()
 defer p.mu.Unlock()

 if p.currentWorker >= int32(p.maxWorkers) {
  return
 }

 p.currentWorker++
 p.wg.Add(1)

 go func(id int32) {
  defer p.wg.Done()
  defer func() {
   p.mu.Lock()
   p.currentWorker--
   p.mu.Unlock()
  }()

  for task := range p.taskQueue {
   task()
  }
 }(p.currentWorker)
}

func (p *DynamicWorkerPool) monitor() {
 ticker := time.NewTicker(time.Second)
 defer ticker.Stop()

 for range ticker.C {
  queueLen := len(p.taskQueue)
  p.mu.Lock()
  current := p.currentWorker
  p.mu.Unlock()

  if queueLen > int(current)*2 && current < int32(p.maxWorkers) {
   p.addWorker()
   fmt.Printf("Scaled up: %d workers\n", current+1)
  }
 }
}

func (p *DynamicWorkerPool) Submit(task func()) {
 p.taskQueue <- task
}

func (p *DynamicWorkerPool) Stop() {
 close(p.taskQueue)
 p.wg.Wait()
}

func main() {
 fmt.Println("=== Basic Worker Pool ===")
 pool := NewWorkerPool(3)
 pool.Start()

 for i := 0; i < 10; i++ {
  n := i
  pool.Submit(func() {
   fmt.Printf("Task %d executed\n", n)
   time.Sleep(100 * time.Millisecond)
  })
 }

 pool.Stop()

 fmt.Println("\n=== Result Worker Pool ===")
 resultPool := NewResultWorkerPool(3)
 resultPool.Start()

 for i := 0; i < 10; i++ {
  resultPool.Submit(Task{ID: i, Data: i})
 }
 resultPool.Stop()

 for result := range resultPool.Results() {
  fmt.Printf("Task %d result: %d\n", result.TaskID, result.Output)
 }
}
```

#### 最佳实践

1. **根据CPU核心数设置worker数量**
2. **使用有缓冲channel作为任务队列**
3. **支持context取消**
4. **考虑动态扩缩容**
5. **监控队列长度和worker利用率**

---

### 5.2 Fan-Out/Fan-In

#### 概念定义

**Fan-Out**将任务分发给多个worker并行处理，**Fan-In**将多个结果合并到一个channel。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// Fan-Out: 分发任务到多个worker
func fanOut() {
 input := make(chan int, 10)

 // 填充输入
 go func() {
  for i := 1; i <= 10; i++ {
   input <- i
  }
  close(input)
 }()

 // 创建多个worker
 numWorkers := 3
 var outputs []<-chan int

 for i := 0; i < numWorkers; i++ {
  outputs = append(outputs, worker(input, i))
 }

 // Fan-In: 合并结果
 merged := merge(outputs...)

 for result := range merged {
  fmt.Printf("Result: %d\n", result)
 }
}

func worker(input <-chan int, id int) <-chan int {
 output := make(chan int)

 go func() {
  defer close(output)
  for n := range input {
   // 模拟处理
   fmt.Printf("Worker %d processing %d\n", id, n)
   time.Sleep(100 * time.Millisecond)
   output <- n * n
  }
 }()

 return output
}

// Fan-In: 合并多个channels
func merge(channels ...<-chan int) <-chan int {
 var wg sync.WaitGroup
 merged := make(chan int)

 output := func(c <-chan int) {
  defer wg.Done()
  for n := range c {
   merged <- n
  }
 }

 wg.Add(len(channels))
 for _, c := range channels {
  go output(c)
 }

 go func() {
  wg.Wait()
  close(merged)
 }()

 return merged
}

// 带错误处理的Fan-Out/Fan-In
type TaskResult struct {
 Value int
 Error error
}

func fanOutFanInWithError() {
 input := make(chan int, 10)

 go func() {
  for i := 1; i <= 10; i++ {
   input <- i
  }
  close(input)
 }()

 numWorkers := 3
 var outputs []<-chan TaskResult

 for i := 0; i < numWorkers; i++ {
  outputs = append(outputs, workerWithError(input, i))
 }

 merged := mergeWithError(outputs...)

 for result := range merged {
  if result.Error != nil {
   fmt.Printf("Error: %v\n", result.Error)
  } else {
   fmt.Printf("Result: %d\n", result.Value)
  }
 }
}

func workerWithError(input <-chan int, id int) <-chan TaskResult {
 output := make(chan TaskResult)

 go func() {
  defer close(output)
  for n := range input {
   if n%7 == 0 {
    output <- TaskResult{Error: fmt.Errorf("worker %d: error processing %d", id, n)}
   } else {
    output <- TaskResult{Value: n * n}
   }
  }
 }()

 return output
}

func mergeWithError(channels ...<-chan TaskResult) <-chan TaskResult {
 var wg sync.WaitGroup
 merged := make(chan TaskResult)

 output := func(c <-chan TaskResult) {
  defer wg.Done()
  for r := range c {
   merged <- r
  }
 }

 wg.Add(len(channels))
 for _, c := range channels {
  go output(c)
 }

 go func() {
  wg.Wait()
  close(merged)
 }()

 return merged
}

// 有序Fan-In
func orderedFanIn() {
 input := make(chan int, 10)

 go func() {
  for i := 1; i <= 10; i++ {
   input <- i
  }
  close(input)
 }()

 numWorkers := 3
 var outputs []<-chan int

 for i := 0; i < numWorkers; i++ {
  outputs = append(outputs, worker(input, i))
 }

 // 有序合并
 ordered := orderedMerge(outputs...)

 for result := range ordered {
  fmt.Printf("Ordered result: %d\n", result)
 }
}

func orderedMerge(channels ...<-chan int) <-chan int {
 // 使用选择器保持顺序
 return nil // 简化实现
}

// 限流Fan-Out
func throttledFanOut() {
 input := make(chan int, 100)

 go func() {
  for i := 1; i <= 100; i++ {
   input <- i
  }
  close(input)
 }()

 // 使用信号量限流
 sem := make(chan struct{}, 5) // 最多5个并发

 var wg sync.WaitGroup
 for n := range input {
  wg.Add(1)
  sem <- struct{}{}

  go func(num int) {
   defer wg.Done()
   defer func() { <-sem }()

   fmt.Printf("Processing %d\n", num)
   time.Sleep(100 * time.Millisecond)
  }(n)
 }

 wg.Wait()
}

func main() {
 fmt.Println("=== Fan-Out/Fan-In ===")
 fanOut()

 fmt.Println("\n=== Fan-Out/Fan-In with Error ===")
 fanOutFanInWithError()

 fmt.Println("\n=== Throttled Fan-Out ===")
 throttledFanOut()
}
```

#### 最佳实践

1. **worker数量根据任务特性调整**
2. **使用WaitGroup确保所有结果收集**
3. **考虑有序输出需求**
4. **实现错误传播机制**
5. **必要时使用限流**

---

### 5.3 Pipeline（管道）

#### 概念定义

**Pipeline**将数据处理分解为多个阶段，每个阶段通过Channel连接，数据顺序流过各阶段。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
)

// Pipeline阶段定义
// Stage 1: 生成数据
func generator(nums ...int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  for _, n := range nums {
   out <- n
  }
 }()
 return out
}

// Stage 2: 平方
func square(in <-chan int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  for n := range in {
   out <- n * n
  }
 }()
 return out
}

// Stage 3: 过滤奇数
func filterOdd(in <-chan int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  for n := range in {
   if n%2 == 0 {
    out <- n
   }
  }
 }()
 return out
}

// Stage 4: 求和
func sum(in <-chan int) <-chan int {
 out := make(chan int)
 go func() {
  defer close(out)
  total := 0
  for n := range in {
   total += n
  }
  out <- total
 }()
 return out
}

// 带缓冲的Pipeline
func bufferedPipeline() {
 generator := func() <-chan int {
  out := make(chan int, 10)
  go func() {
   defer close(out)
   for i := 1; i <= 20; i++ {
    out <- i
   }
  }()
  return out
 }

 process := func(in <-chan int) <-chan int {
  out := make(chan int, 10)
  go func() {
   defer close(out)
   for n := range in {
    out <- n * 2
   }
  }()
  return out
 }

 ch := process(generator())
 for n := range ch {
  fmt.Printf("%d ", n)
 }
 fmt.Println()
}

// 并行Pipeline
func parallelPipeline() {
 generator := func() <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for i := 1; i <= 20; i++ {
    out <- i
   }
  }()
  return out
 }

 // 多个并行worker
 parallelProcess := func(in <-chan int, workers int) <-chan int {
  var outputs []<-chan int

  for i := 0; i < workers; i++ {
   out := make(chan int)
   outputs = append(outputs, out)

   go func(id int, output chan<- int) {
    defer close(output)
    for n := range in {
     fmt.Printf("Worker %d processing %d\n", id, n)
     output <- n * n
    }
   }(i, out)
  }

  return merge(outputs...)
 }

 ch := parallelProcess(generator(), 3)
 for n := range ch {
  fmt.Printf("Result: %d\n", n)
 }
}

// 带取消的Pipeline
func cancellablePipeline() {
 done := make(chan struct{})

 // 稍后取消
 go func() {
  // close(done)
 }()

 generator := func() <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for i := 1; i <= 100; i++ {
    select {
    case out <- i:
    case <-done:
     return
    }
   }
  }()
  return out
 }

 process := func(in <-chan int) <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for n := range in {
    select {
    case out <- n * 2:
    case <-done:
     return
    }
   }
  }()
  return out
 }

 for n := range process(generator()) {
  fmt.Printf("%d ", n)
  if n > 20 {
   // close(done)
   break
  }
 }
}

// 多输入Pipeline
type Message struct {
 ID      int
 Content string
}

func multiInputPipeline() {
 // 输入1: 用户消息
 userMessages := make(chan Message)
 go func() {
  defer close(userMessages)
  for i := 1; i <= 5; i++ {
   userMessages <- Message{ID: i, Content: fmt.Sprintf("User message %d", i)}
  }
 }()

 // 输入2: 系统消息
 systemMessages := make(chan Message)
 go func() {
  defer close(systemMessages)
  for i := 1; i <= 3; i++ {
   systemMessages <- Message{ID: i + 100, Content: fmt.Sprintf("System message %d", i)}
  }
 }()

 // 合并输入
 merged := mergeMessages(userMessages, systemMessages)

 // 处理
 processed := processMessages(merged)

 for msg := range processed {
  fmt.Printf("Processed: %+v\n", msg)
 }
}

func mergeMessages(channels ...<-chan Message) <-chan Message {
 var wg sync.WaitGroup
 merged := make(chan Message)

 output := func(c <-chan Message) {
  defer wg.Done()
  for msg := range c {
   merged <- msg
  }
 }

 wg.Add(len(channels))
 for _, c := range channels {
  go output(c)
 }

 go func() {
  wg.Wait()
  close(merged)
 }()

 return merged
}

func processMessages(in <-chan Message) <-chan Message {
 out := make(chan Message)
 go func() {
  defer close(out)
  for msg := range in {
   msg.Content = "[PROCESSED] " + msg.Content
   out <- msg
  }
 }()
 return out
}

func main() {
 fmt.Println("=== Basic Pipeline ===")
 result := <-sum(filterOdd(square(generator(1, 2, 3, 4, 5))))
 fmt.Printf("Result: %d\n", result)

 fmt.Println("\n=== Buffered Pipeline ===")
 bufferedPipeline()

 fmt.Println("\n=== Multi Input Pipeline ===")
 multiInputPipeline()
}
```

#### 最佳实践

1. **每个stage独立关闭输出channel**
2. **使用单向channel类型**
3. **支持取消机制**
4. **适当使用缓冲提高吞吐量**
5. **考虑并行stage提高性能**

---

### 5.4 Tee（分流）

#### 概念定义

**Tee**模式将一个输入channel复制到多个输出channel，实现数据的多路分发。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
)

// 基本Tee
tee := func(in <-chan int, outs ...chan<- int) {
 go func() {
  defer func() {
   for _, out := range outs {
    close(out)
   }
  }()

  for v := range in {
   for _, out := range outs {
    out <- v
   }
  }
 }()
}

// 带条件的Tee
teeWithCondition := func(in <-chan int, evenOut, oddOut chan<- int) {
 go func() {
  defer close(evenOut)
  defer close(oddOut)

  for v := range in {
   if v%2 == 0 {
    evenOut <- v
   } else {
    oddOut <- v
   }
  }
 }()
}

// 广播Tee
type BroadcastTee struct {
 listeners []chan int
 mu        sync.RWMutex
}

func NewBroadcastTee() *BroadcastTee {
 return &BroadcastTee{
  listeners: make([]chan int, 0),
 }
}

func (b *BroadcastTee) Subscribe() <-chan int {
 ch := make(chan int, 10)
 b.mu.Lock()
 b.listeners = append(b.listeners, ch)
 b.mu.Unlock()
 return ch
}

func (b *BroadcastTee) Unsubscribe(ch <-chan int) {
 b.mu.Lock()
 defer b.mu.Unlock()

 for i, listener := range b.listeners {
  if listener == ch {
   b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
   close(listener)
   break
  }
 }
}

func (b *BroadcastTee) Send(v int) {
 b.mu.RLock()
 listeners := make([]chan int, len(b.listeners))
 copy(listeners, b.listeners)
 b.mu.RUnlock()

 for _, ch := range listeners {
  select {
  case ch <- v:
  default:
   // 缓冲区满，丢弃
  }
 }
}

func (b *BroadcastTee) Close() {
 b.mu.Lock()
 defer b.mu.Unlock()

 for _, ch := range b.listeners {
  close(ch)
 }
 b.listeners = nil
}

// 使用示例
func teeExample() {
 input := make(chan int)
 out1, out2 := make(chan int), make(chan int)

 // 启动tee
 go func() {
  defer close(out1)
  defer close(out2)

  for v := range input {
   out1 <- v
   out2 <- v
  }
 }()

 // 消费者1
 go func() {
  for v := range out1 {
   fmt.Printf("Consumer 1: %d\n", v)
  }
 }()

 // 消费者2
 go func() {
  for v := range out2 {
   fmt.Printf("Consumer 2: %d\n", v)
  }
 }()

 // 发送数据
 for i := 1; i <= 5; i++ {
  input <- i
 }
 close(input)
}

// 日志分流示例
func logTeeExample() {
 type LogEntry struct {
  Level   string
  Message string
 }

 logs := make(chan LogEntry)

 // 分流到不同级别
 debugLogs := make(chan LogEntry)
 infoLogs := make(chan LogEntry)
 errorLogs := make(chan LogEntry)

 go func() {
  defer close(debugLogs)
  defer close(infoLogs)
  defer close(errorLogs)

  for log := range logs {
   switch log.Level {
   case "DEBUG":
    debugLogs <- log
   case "INFO":
    infoLogs <- log
   case "ERROR":
    errorLogs <- log
   }
  }
 }()

 // 处理器
 go func() {
  for log := range debugLogs {
   fmt.Printf("[DEBUG] %s\n", log.Message)
  }
 }()

 go func() {
  for log := range infoLogs {
   fmt.Printf("[INFO] %s\n", log.Message)
  }
 }()

 go func() {
  for log := range errorLogs {
   fmt.Printf("[ERROR] %s\n", log.Message)
  }
 }()

 // 发送日志
 logs <- LogEntry{"DEBUG", "debug message"}
 logs <- LogEntry{"INFO", "info message"}
 logs <- LogEntry{"ERROR", "error message"}

 close(logs)
}

func main() {
 fmt.Println("=== Tee Example ===")
 teeExample()

 fmt.Println("\n=== Log Tee Example ===")
 logTeeExample()
}
```

#### 最佳实践

1. **确保所有输出channel都被关闭**
2. **考虑使用缓冲避免阻塞**
3. **处理慢消费者问题**
4. **支持动态订阅/取消订阅**
5. **考虑背压机制**

---

### 5.5 Bridge（桥接）

#### 概念定义

**Bridge**模式将channel的channel（`<-chan <-chan T`）扁平化为单个channel，用于动态创建子pipeline。

#### 完整示例

```go
package main

import (
 "fmt"
)

// Bridge模式: 将chan chan T 转换为 chan T
func bridge[T any](chanStream <-chan <-chan T) <-chan T {
 valStream := make(chan T)

 go func() {
  defer close(valStream)

  for {
   var stream <-chan T
   select {
   case maybeStream, ok := <-chanStream:
    if !ok {
     return
    }
    stream = maybeStream
   }

   for val := range stream {
    select {
    case valStream <- val:
    }
   }
  }
 }()

 return valStream
}

// 使用示例
genVals := func() <-chan <-chan int {
 chanStream := make(chan (<-chan int))

 go func() {
  defer close(chanStream)

  for i := 0; i < 3; i++ {
   stream := make(chan int, 3)
   stream <- i * 1
   stream <- i * 2
   stream <- i * 3
   close(stream)
   chanStream <- stream
  }
 }()

 return chanStream
}

// 动态任务生成
func dynamicTaskGeneration() {
 type Task struct {
  ID   int
  Data string
 }

 // 生成任务流
 taskStreams := make(chan (<-chan Task))

 go func() {
  defer close(taskStreams)

  // 模拟动态生成任务批次
  for batch := 0; batch < 3; batch++ {
   tasks := make(chan Task, 5)

   for i := 0; i < 5; i++ {
    tasks <- Task{
     ID:   batch*10 + i,
     Data: fmt.Sprintf("batch-%d-task-%d", batch, i),
    }
   }

   close(tasks)
   taskStreams <- tasks
  }
 }()

 // 使用bridge处理所有任务
 for task := range bridge(taskStreams) {
  fmt.Printf("Processing task: %+v\n", task)
 }
}

// 嵌套Pipeline
func nestedPipeline() {
 // 阶段1: 生成批次
 batches := func() <-chan <-chan int {
  out := make(chan (<-chan int))

  go func() {
   defer close(out)

   for i := 0; i < 3; i++ {
    batch := make(chan int, 5)
    for j := 0; j < 5; j++ {
     batch <- i*10 + j
    }
    close(batch)
    out <- batch
   }
  }()

  return out
 }()

 // 阶段2: 使用bridge处理
 processed := make(chan int)
 go func() {
  defer close(processed)
  for val := range bridge(batches) {
   processed <- val * 2
  }
 }()

 // 收集结果
 for result := range processed {
  fmt.Printf("Result: %d\n", result)
 }
}

func main() {
 fmt.Println("=== Basic Bridge ===")
 for val := range bridge(genVals()) {
  fmt.Printf("Value: %d\n", val)
 }

 fmt.Println("\n=== Dynamic Task Generation ===")
 dynamicTaskGeneration()

 fmt.Println("\n=== Nested Pipeline ===")
 nestedPipeline()
}
```

#### 最佳实践

1. **确保所有子channel都被消费**
2. **处理nil channel情况**
3. **支持取消机制**
4. **考虑背压处理**

---

### 5.6 Or-Done（或完成）

#### 概念定义

**Or-Done**模式将输入channel与done channel结合，当任一channel关闭时返回。

#### 完整示例

```go
package main

import (
 "fmt"
 "time"
)

// Or-Done: 当input或done任一关闭时返回
func orDone[T any](done <-chan struct{}, c <-chan T) <-chan T {
 valStream := make(chan T)

 go func() {
  defer close(valStream)

  for {
   select {
   case <-done:
    return
   case v, ok := <-c:
    if !ok {
     return
    }
    select {
    case valStream <- v:
    case <-done:
     return
    }
   }
  }
 }()

 return valStream
}

// 使用示例
func orDoneExample() {
 done := make(chan struct{})
 input := make(chan int)

 // 发送数据
 go func() {
  for i := 0; i < 10; i++ {
   input <- i
   time.Sleep(100 * time.Millisecond)
  }
  close(input)
 }()

 // 稍后取消
 go func() {
  time.Sleep(350 * time.Millisecond)
  close(done)
 }()

 // 消费
 for val := range orDone(done, input) {
  fmt.Printf("Received: %d\n", val)
 }

 fmt.Println("Done")
}

// Pipeline中的Or-Done
func pipelineOrDone() {
 generator := func(done <-chan struct{}) <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for i := 0; i < 100; i++ {
    select {
    case out <- i:
    case <-done:
     return
    }
   }
  }()
  return out
 }

 process := func(done <-chan struct{}, in <-chan int) <-chan int {
  out := make(chan int)
  go func() {
   defer close(out)
   for val := range orDone(done, in) {
    select {
    case out <- val * 2:
    case <-done:
     return
    }
   }
  }()
  return out
 }

 done := make(chan struct{})

 // 稍后取消
 go func() {
  time.Sleep(200 * time.Millisecond)
  close(done)
 }()

 for val := range process(done, generator(done)) {
  fmt.Printf("%d ", val)
 }
 fmt.Println()
}

func main() {
 fmt.Println("=== Or-Done Example ===")
 orDoneExample()

 fmt.Println("\n=== Pipeline Or-Done ===")
 pipelineOrDone()
}
```

#### 最佳实践

1. **在每个pipeline stage使用orDone**
2. **优先检查done channel**
3. **确保及时关闭输出channel**

---

### 5.7 Or-Channel（或通道）

#### 概念定义

**Or-Channel**模式将多个channel合并，当任一channel有数据时返回。

#### 完整示例

```go
package main

import (
 "fmt"
 "time"
)

// Or-Channel: 返回任一channel的数据
func or[T any](channels ...<-chan T) <-chan T {
 switch len(channels) {
 case 0:
  return nil
 case 1:
  return channels[0]
 }

 orDone := make(chan T)
 go func() {
  defer close(orDone)

  switch len(channels) {
  case 2:
   select {
   case <-channels[0]:
   case <-channels[1]:
   }
  default:
   select {
   case <-channels[0]:
   case <-channels[1]:
   case <-channels[2]:
   case <-or(append(channels[3:], orDone)...):
   }
  }
 }()

 return orDone
}

// 使用示例
func orChannelExample() {
 sig := func(after time.Duration) <-chan struct{} {
  c := make(chan struct{})
  go func() {
   defer close(c)
   time.Sleep(after)
  }()
  return c
 }

 start := time.Now()
 <-or(sig(2*time.Hour), sig(5*time.Minute), sig(1*time.Second), sig(1*time.Hour), sig(1*time.Minute))
 fmt.Printf("Done after %v\n", time.Since(start))
}

// 超时等待多个事件
func waitForAnyEvent() {
 eventA := make(chan string)
 eventB := make(chan string)
 eventC := make(chan string)

 go func() { time.Sleep(100 * time.Millisecond); eventA <- "Event A" }()
 go func() { time.Sleep(200 * time.Millisecond); eventB <- "Event B" }()
 go func() { time.Sleep(300 * time.Millisecond); eventC <- "Event C" }()

 select {
 case msg := <-eventA:
  fmt.Println("Received:", msg)
 case msg := <-eventB:
  fmt.Println("Received:", msg)
 case msg := <-eventC:
  fmt.Println("Received:", msg)
 case <-time.After(150 * time.Millisecond):
  fmt.Println("Timeout")
 }
}

// 服务器关闭信号
func serverShutdownSignal() {
 shutdown := make(chan struct{})
 sigint := make(chan struct{})
 sigterm := make(chan struct{})

 go func() {
  time.Sleep(500 * time.Millisecond)
  close(sigint)
 }()

 go func() {
  time.Sleep(1 * time.Second)
  close(sigterm)
 }()

 go func() {
  time.Sleep(2 * time.Second)
  close(shutdown)
 }()

 select {
 case <-sigint:
  fmt.Println("Received SIGINT")
 case <-sigterm:
  fmt.Println("Received SIGTERM")
 case <-shutdown:
  fmt.Println("Shutdown requested")
 }
}

func main() {
 fmt.Println("=== Or-Channel Example ===")
 orChannelExample()

 fmt.Println("\n=== Wait for Any Event ===")
 waitForAnyEvent()

 fmt.Println("\n=== Server Shutdown Signal ===")
 serverShutdownSignal()
}
```

#### 最佳实践

1. **用于等待多个信号**
2. **实现超时机制**
3. **处理服务器关闭**
4. **注意递归深度限制**

---

### 5.8 Quit信号模式

#### 概念定义

**Quit信号模式**使用专门的channel来通知Goroutine优雅退出。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 基本Quit信号
func basicQuitSignal() {
 quit := make(chan struct{})

 go func() {
  for {
   select {
   case <-quit:
    fmt.Println("Worker quitting...")
    return
   default:
    fmt.Println("Working...")
    time.Sleep(200 * time.Millisecond)
   }
  }
 }()

 time.Sleep(500 * time.Millisecond)
 close(quit)
 time.Sleep(100 * time.Millisecond)
}

// 带确认的Quit
func acknowledgedQuit() {
 quit := make(chan struct{})
 ack := make(chan struct{})

 go func() {
  for {
   select {
   case <-quit:
    fmt.Println("Worker cleaning up...")
    time.Sleep(200 * time.Millisecond)
    close(ack)
    return
   default:
    fmt.Println("Working...")
    time.Sleep(100 * time.Millisecond)
   }
  }
 }()

 time.Sleep(300 * time.Millisecond)
 close(quit)

 <-ack
 fmt.Println("Worker acknowledged quit")
}

// 多个Worker的Quit
func multiWorkerQuit() {
 quit := make(chan struct{})
 var wg sync.WaitGroup

 for i := 0; i < 3; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()

   for {
    select {
    case <-quit:
     fmt.Printf("Worker %d quitting\n", id)
     return
    default:
     fmt.Printf("Worker %d working\n", id)
     time.Sleep(100 * time.Millisecond)
    }
   }
  }(i)
 }

 time.Sleep(300 * time.Millisecond)
 close(quit)
 wg.Wait()
 fmt.Println("All workers quit")
}

// 优雅关闭服务器
func gracefulServerShutdown() {
 quit := make(chan struct{})

 // 模拟服务器
 server := func() {
  connections := make(chan int, 10)

  // 接受连接
  go func() {
   for i := 0; i < 10; i++ {
    connections <- i
    time.Sleep(50 * time.Millisecond)
   }
   close(connections)
  }()

  // 处理连接
  var wg sync.WaitGroup
  for conn := range connections {
   wg.Add(1)
   go func(c int) {
    defer wg.Done()

    select {
    case <-quit:
     fmt.Printf("Connection %d: server shutting down\n", c)
     return
    default:
     fmt.Printf("Connection %d: processing\n", c)
     time.Sleep(200 * time.Millisecond)
     fmt.Printf("Connection %d: completed\n", c)
    }
   }(conn)
  }

  wg.Wait()
 }

 go server()

 time.Sleep(300 * time.Millisecond)
 fmt.Println("Shutting down server...")
 close(quit)

 time.Sleep(500 * time.Millisecond)
}

// 带超时的Quit
func quitWithTimeout() {
 quit := make(chan struct{})
 done := make(chan struct{})

 go func() {
  defer close(done)

  for {
   select {
   case <-quit:
    fmt.Println("Received quit, cleaning up...")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Cleanup done")
    return
   default:
    fmt.Println("Working...")
    time.Sleep(100 * time.Millisecond)
   }
  }
 }()

 close(quit)

 select {
 case <-done:
  fmt.Println("Worker quit gracefully")
 case <-time.After(200 * time.Millisecond):
  fmt.Println("Timeout waiting for quit")
 }
}

func main() {
 fmt.Println("=== Basic Quit Signal ===")
 basicQuitSignal()

 fmt.Println("\n=== Acknowledged Quit ===")
 acknowledgedQuit()

 fmt.Println("\n=== Multi Worker Quit ===")
 multiWorkerQuit()

 fmt.Println("\n=== Graceful Server Shutdown ===")
 gracefulServerShutdown()

 fmt.Println("\n=== Quit with Timeout ===")
 quitWithTimeout()
}
```

#### 最佳实践

1. **使用struct{}类型节省空间**
2. **优先使用close而非发送值**
3. **支持优雅关闭和清理**
4. **考虑超时机制**
5. **使用WaitGroup等待所有worker退出**

---

### 5.9 速率限制（Rate Limiting）

#### 概念定义

**速率限制**控制操作的执行频率，防止系统过载。

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "time"
)

// 令牌桶速率限制
type TokenBucket struct {
 tokens   chan struct{}
 interval time.Duration
}

func NewTokenBucket(rate int, interval time.Duration) *TokenBucket {
 tb := &TokenBucket{
  tokens:   make(chan struct{}, rate),
  interval: interval,
 }

 // 填充令牌
 for i := 0; i < rate; i++ {
  tb.tokens <- struct{}{}
 }

 // 持续补充
 go func() {
  ticker := time.NewTicker(interval / time.Duration(rate))
  defer ticker.Stop()

  for range ticker.C {
   select {
   case tb.tokens <- struct{}{}:
   default:
    // 桶满，丢弃
   }
  }
 }()

 return tb
}

func (tb *TokenBucket) Wait() {
 <-tb.tokens
}

// 漏桶速率限制
type LeakyBucket struct {
 requests chan struct{}
 rate     time.Duration
}

func NewLeakyBucket(rate time.Duration) *LeakyBucket {
 lb := &LeakyBucket{
  requests: make(chan struct{}, 100),
  rate:     rate,
 }

 go func() {
  ticker := time.NewTicker(rate)
  defer ticker.Stop()

  for range ticker.C {
   select {
   case <-lb.requests:
    // 处理一个请求
   default:
    // 无请求
   }
  }
 }()

 return lb
}

func (lb *LeakyBucket) Allow() bool {
 select {
 case lb.requests <- struct{}{}:
  return true
 default:
  return false
 }
}

// 固定窗口计数器
type FixedWindow struct {
 limit    int
 count    int
 window   time.Duration
 resetAt  time.Time
}

func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
 return &FixedWindow{
  limit:   limit,
  window:  window,
  resetAt: time.Now().Add(window),
 }
}

func (fw *FixedWindow) Allow() bool {
 now := time.Now()

 if now.After(fw.resetAt) {
  fw.count = 0
  fw.resetAt = now.Add(fw.window)
 }

 if fw.count < fw.limit {
  fw.count++
  return true
 }

 return false
}

// 滑动窗口日志
type SlidingWindow struct {
 limit   int
 window  time.Duration
 timestamps []time.Time
}

func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
 return &SlidingWindow{
  limit:  limit,
  window: window,
 }
}

func (sw *SlidingWindow) Allow() bool {
 now := time.Now()
 cutoff := now.Add(-sw.window)

 // 移除过期时间戳
 var valid []time.Time
 for _, ts := range sw.timestamps {
  if ts.After(cutoff) {
   valid = append(valid, ts)
  }
 }
 sw.timestamps = valid

 if len(sw.timestamps) < sw.limit {
  sw.timestamps = append(sw.timestamps, now)
  return true
 }

 return false
}

// 使用示例
func rateLimitExample() {
 // 令牌桶: 每秒10个请求
 tokenBucket := NewTokenBucket(10, time.Second)

 fmt.Println("Token Bucket:")
 for i := 0; i < 15; i++ {
  start := time.Now()
  tokenBucket.Wait()
  fmt.Printf("Request %d at %v\n", i+1, time.Since(start))
 }
}

// API速率限制中间件
func apiRateLimit() {
 limiter := NewTokenBucket(5, time.Second)

 apiHandler := func(id int) {
  limiter.Wait()
  fmt.Printf("API request %d processed at %v\n", id, time.Now().Format("15:04:05.000"))
 }

 // 模拟并发请求
 for i := 0; i < 10; i++ {
  go apiHandler(i)
 }

 time.Sleep(3 * time.Second)
}

// 自适应速率限制
type AdaptiveRateLimiter struct {
 limiter    *TokenBucket
 successes  int
 failures   int
 mu         sync.Mutex
}

import "sync"

func (a *AdaptiveRateLimiter) RecordSuccess() {
 a.mu.Lock()
 a.successes++
 a.mu.Unlock()
}

func (a *AdaptiveRateLimiter) RecordFailure() {
 a.mu.Lock()
 a.failures++
 a.mu.Unlock()
}

func main() {
 fmt.Println("=== Rate Limit Example ===")
 rateLimitExample()

 fmt.Println("\n=== API Rate Limit ===")
 apiRateLimit()
}
```

#### 最佳实践

1. **根据场景选择合适的算法**
2. **令牌桶适合突发流量**
3. **漏桶适合平滑流量**
4. **滑动窗口更精确但开销大**
5. **考虑分布式速率限制**

---

### 5.10 防抖与节流

#### 概念定义

**防抖**（Debounce）延迟执行直到停止触发一段时间后执行，**节流**（Throttle）限制执行频率。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 防抖
type Debouncer struct {
 delay   time.Duration
 timer   *time.Timer
 mu      sync.Mutex
}

func NewDebouncer(delay time.Duration) *Debouncer {
 return &Debouncer{delay: delay}
}

func (d *Debouncer) Do(f func()) {
 d.mu.Lock()
 defer d.mu.Unlock()

 if d.timer != nil {
  d.timer.Stop()
 }

 d.timer = time.AfterFunc(d.delay, f)
}

// 节流
type Throttler struct {
 delay    time.Duration
 lastCall time.Time
 mu       sync.Mutex
}

func NewThrottler(delay time.Duration) *Throttler {
 return &Throttler{
  delay:    delay,
  lastCall: time.Now().Add(-delay),
 }
}

func (t *Throttler) Do(f func()) {
 t.mu.Lock()
 defer t.mu.Unlock()

 now := time.Now()
 if now.Sub(t.lastCall) >= t.delay {
  t.lastCall = now
  f()
 }
}

// 带 leading/trailing 的防抖
type AdvancedDebouncer struct {
 delay    time.Duration
 leading  bool
 trailing bool
 timer    *time.Timer
 mu       sync.Mutex
}

func NewAdvancedDebouncer(delay time.Duration, leading, trailing bool) *AdvancedDebouncer {
 return &AdvancedDebouncer{
  delay:    delay,
  leading:  leading,
  trailing: trailing,
 }
}

func (d *AdvancedDebouncer) Do(f func()) {
 d.mu.Lock()
 defer d.mu.Unlock()

 isFirstCall := d.timer == nil

 if d.timer != nil {
  d.timer.Stop()
 }

 if d.leading && isFirstCall {
  f()
 }

 if d.trailing {
  d.timer = time.AfterFunc(d.delay, f)
 }
}

// 使用示例
func debounceExample() {
 debouncer := NewDebouncer(300 * time.Millisecond)

 counter := 0
 for i := 0; i < 10; i++ {
  debouncer.Do(func() {
   counter++
   fmt.Printf("Debounced call #%d at %v\n", counter, time.Now().Format("15:04:05.000"))
  })
  time.Sleep(100 * time.Millisecond)
 }

 time.Sleep(500 * time.Millisecond)
 fmt.Printf("Total calls: %d (expected 1)\n", counter)
}

func throttleExample() {
 throttler := NewThrottler(200 * time.Millisecond)

 counter := 0
 for i := 0; i < 10; i++ {
  throttler.Do(func() {
   counter++
   fmt.Printf("Throttled call #%d at %v\n", counter, time.Now().Format("15:04:05.000"))
  })
  time.Sleep(50 * time.Millisecond)
 }

 time.Sleep(300 * time.Millisecond)
 fmt.Printf("Total calls: %d (expected ~3)\n", counter)
}

// 搜索防抖
func searchDebounce() {
 type SearchFunc func(query string)

 debouncedSearch := func(search SearchFunc, delay time.Duration) SearchFunc {
  var timer *time.Timer
  var mu sync.Mutex

  return func(query string) {
   mu.Lock()
   defer mu.Unlock()

   if timer != nil {
    timer.Stop()
   }

   timer = time.AfterFunc(delay, func() {
    search(query)
   })
  }
 }

 search := func(query string) {
  fmt.Printf("Searching for: %s\n", query)
 }

 debounced := debouncedSearch(search, 300*time.Millisecond)

 // 模拟快速输入
 debounced("a")
 time.Sleep(100 * time.Millisecond)
 debounced("ab")
 time.Sleep(100 * time.Millisecond)
 debounced("abc")

 time.Sleep(500 * time.Millisecond)
}

// 滚动节流
func scrollThrottle() {
 throttler := NewThrottler(100 * time.Millisecond)

 scrollHandler := func() {
  fmt.Printf("Scroll event at %v\n", time.Now().Format("15:04:05.000"))
 }

 // 模拟频繁滚动
 for i := 0; i < 20; i++ {
  throttler.Do(scrollHandler)
  time.Sleep(10 * time.Millisecond)
 }
}

func main() {
 fmt.Println("=== Debounce Example ===")
 debounceExample()

 fmt.Println("\n=== Throttle Example ===")
 throttleExample()

 fmt.Println("\n=== Search Debounce ===")
 searchDebounce()

 fmt.Println("\n=== Scroll Throttle ===")
 scrollThrottle()
}
```

#### 最佳实践

1. **防抖适合输入、搜索等场景**
2. **节流适合滚动、resize等高频事件**
3. **考虑leading/trailing选项**
4. **注意内存泄漏（及时停止timer）**
5. **使用mutex保证线程安全**

---

## 6. 并行计算模式

### 6.1 并行Map

#### 概念定义

**并行Map**将map操作并行化，对集合中的每个元素应用函数，利用多核提高处理速度。

#### 完整示例

```go
package main

import (
 "fmt"
 "math"
 "runtime"
 "sync"
 "time"
)

// 基本并行Map
func parallelMap[T any, R any](input []T, fn func(T) R) []R {
 numCPU := runtime.NumCPU()
 chunkSize := (len(input) + numCPU - 1) / numCPU

 results := make([]R, len(input))
 var wg sync.WaitGroup

 for i := 0; i < numCPU; i++ {
  start := i * chunkSize
  end := start + chunkSize
  if end > len(input) {
   end = len(input)
  }

  if start >= len(input) {
   break
  }

  wg.Add(1)
  go func(s, e int) {
   defer wg.Done()
   for j := s; j < e; j++ {
    results[j] = fn(input[j])
   }
  }(start, end)
 }

 wg.Wait()
 return results
}

// 带顺序保证的并行Map
func parallelMapOrdered[T any, R any](input []T, fn func(T) R) []R {
 results := make([]R, len(input))
 var wg sync.WaitGroup

 for i, v := range input {
  wg.Add(1)
  go func(index int, value T) {
   defer wg.Done()
   results[index] = fn(value)
  }(i, v)
 }

 wg.Wait()
 return results
}

// 带错误处理的并行Map
func parallelMapWithError[T any, R any](input []T, fn func(T) (R, error)) ([]R, error) {
 type result struct {
  index int
  value R
  err   error
 }

 results := make([]R, len(input))
 resultCh := make(chan result, len(input))

 var wg sync.WaitGroup
 for i, v := range input {
  wg.Add(1)
  go func(index int, value T) {
   defer wg.Done()
   val, err := fn(value)
   resultCh <- result{index, val, err}
  }(i, v)
 }

 go func() {
  wg.Wait()
  close(resultCh)
 }()

 for r := range resultCh {
  if r.err != nil {
   return nil, r.err
  }
  results[r.index] = r.value
 }

 return results, nil
}

// 带限制的并行Map
func parallelMapLimited[T any, R any](input []T, fn func(T) R, maxConcurrency int) []R {
 results := make([]R, len(input))
 sem := make(chan struct{}, maxConcurrency)
 var wg sync.WaitGroup

 for i, v := range input {
  wg.Add(1)
  sem <- struct{}{}
  go func(index int, value T) {
   defer wg.Done()
   defer func() { <-sem }()
   results[index] = fn(value)
  }(i, v)
 }

 wg.Wait()
 return results
}

// 使用示例
func mapExample() {
 numbers := make([]int, 1000000)
 for i := range numbers {
  numbers[i] = i + 1
 }

 // 串行
 start := time.Now()
 serialResults := make([]float64, len(numbers))
 for i, n := range numbers {
  serialResults[i] = math.Sqrt(float64(n))
 }
 serialTime := time.Since(start)

 // 并行
 start = time.Now()
 parallelResults := parallelMap(numbers, func(n int) float64 {
  return math.Sqrt(float64(n))
 })
 parallelTime := time.Since(start)

 fmt.Printf("Serial: %v\n", serialTime)
 fmt.Printf("Parallel: %v\n", parallelTime)
 fmt.Printf("Speedup: %.2fx\n", float64(serialTime)/float64(parallelTime))

 // 验证结果
 for i := range serialResults {
  if serialResults[i] != parallelResults[i] {
   fmt.Printf("Mismatch at index %d\n", i)
   break
  }
 }
}

// 图像处理示例
func imageProcessingExample() {
 type Pixel struct{ R, G, B uint8 }

 // 模拟图像
 width, height := 1920, 1080
 image := make([][]Pixel, height)
 for i := range image {
  image[i] = make([]Pixel, width)
 }

 // 并行灰度转换
 grayscale := parallelMap(image, func(row []Pixel) []uint8 {
  result := make([]uint8, len(row))
  for i, p := range row {
   gray := uint8(0.299*float64(p.R) + 0.587*float64(p.G) + 0.114*float64(p.B))
   result[i] = gray
  }
  return result
 })

 fmt.Printf("Processed %d rows in parallel\n", len(grayscale))
}

func main() {
 fmt.Println("=== Map Example ===")
 mapExample()

 fmt.Println("\n=== Image Processing Example ===")
 imageProcessingExample()
}
```

#### 最佳实践

1. **任务粒度要足够大**
2. **使用分块减少goroutine数量**
3. **考虑CPU核心数**
4. **验证并行结果正确性**
5. **注意内存带宽限制**

---

### 6.2 并行Reduce

#### 概念定义

**并行Reduce**将reduce操作并行化，通过分治策略并行聚合数据。

#### 完整示例

```go
package main

import (
 "fmt"
 "runtime"
 "sync"
 "time"
)

// 基本并行Reduce
func parallelReduce[T any](input []T, identity T, fn func(T, T) T) T {
 if len(input) == 0 {
  return identity
 }

 if len(input) == 1 {
  return input[0]
 }

 numCPU := runtime.NumCPU()
 if len(input) < numCPU*2 {
  // 数据量小，串行处理
  result := input[0]
  for i := 1; i < len(input); i++ {
   result = fn(result, input[i])
  }
  return result
 }

 chunkSize := (len(input) + numCPU - 1) / numCPU
 partialResults := make([]T, numCPU)
 var wg sync.WaitGroup

 for i := 0; i < numCPU; i++ {
  start := i * chunkSize
  end := start + chunkSize
  if end > len(input) {
   end = len(input)
  }

  if start >= len(input) {
   partialResults[i] = identity
   continue
  }

  wg.Add(1)
  go func(idx, s, e int) {
   defer wg.Done()
   result := input[s]
   for j := s + 1; j < e; j++ {
    result = fn(result, input[j])
   }
   partialResults[idx] = result
  }(i, start, end)
 }

 wg.Wait()

 // 合并部分结果
 finalResult := partialResults[0]
 for i := 1; i < numCPU; i++ {
  finalResult = fn(finalResult, partialResults[i])
 }

 return finalResult
}

// 并行Sum
func parallelSum(input []int) int {
 return parallelReduce(input, 0, func(a, b int) int {
  return a + b
 })
}

// 并行Max
func parallelMax(input []int) int {
 if len(input) == 0 {
  return 0
 }
 return parallelReduce(input, input[0], func(a, b int) int {
  if a > b {
   return a
  }
  return b
 })
}

// 并行字符串连接
func parallelStringConcat(input []string) string {
 return parallelReduce(input, "", func(a, b string) string {
  return a + b
 })
}

// 使用channel的并行Reduce
func parallelReduceChannel[T any](input []T, identity T, fn func(T, T) T) T {
 if len(input) == 0 {
  return identity
 }

 ch := make(chan T, len(input))

 // 发送所有元素
 for _, v := range input {
  ch <- v
 }
 close(ch)

 // 两两合并
 for len(ch) > 1 {
  a := <-ch
  b := <-ch
  go func(x, y T) {
   ch <- fn(x, y)
  }(a, b)
 }

 return <-ch
}

// MapReduce组合
func mapReduce[T any, M any, R any](
 input []T,
 mapFn func(T) M,
 reduceFn func(R, M) R,
 identity R,
) R {
 // 并行Map
 mapped := parallelMap(input, mapFn)

 // 并行Reduce
 return parallelReduce(mapped, identity, reduceFn)
}

// 使用示例
func reduceExample() {
 numbers := make([]int, 10000000)
 for i := range numbers {
  numbers[i] = i + 1
 }

 // 串行Sum
 start := time.Now()
 serialSum := 0
 for _, n := range numbers {
  serialSum += n
 }
 serialTime := time.Since(start)

 // 并行Sum
 start = time.Now()
 parallelSumResult := parallelSum(numbers)
 parallelTime := time.Since(start)

 fmt.Printf("Serial sum: %d, time: %v\n", serialSum, serialTime)
 fmt.Printf("Parallel sum: %d, time: %v\n", parallelSumResult, parallelTime)
 fmt.Printf("Speedup: %.2fx\n", float64(serialTime)/float64(parallelTime))

 // Max
 max := parallelMax(numbers)
 fmt.Printf("Max: %d\n", max)
}

// Word Count MapReduce
func wordCountMapReduce() {
 documents := []string{
  "hello world",
  "hello go",
  "go concurrency",
  "hello concurrency",
 }

 // Map: 统计每个文档的词频
 mapFn := func(doc string) map[string]int {
  freq := make(map[string]int)
  words := strings.Fields(doc)
  for _, word := range words {
   freq[word]++
  }
  return freq
 }

 // Reduce: 合并词频
 reduceFn := func(acc, curr map[string]int) map[string]int {
  for word, count := range curr {
   acc[word] += count
  }
  return acc
 }

 result := mapReduce(documents, mapFn, reduceFn, make(map[string]int))

 fmt.Println("Word counts:")
 for word, count := range result {
  fmt.Printf("  %s: %d\n", word, count)
 }
}

import "strings"

func main() {
 fmt.Println("=== Reduce Example ===")
 reduceExample()

 fmt.Println("\n=== Word Count MapReduce ===")
 wordCountMapReduce()
}
```

#### 最佳实践

1. **确保reduce函数满足结合律**
2. **小数据量使用串行**
3. **注意浮点数精度问题**
4. **合理设置分块大小**

---

### 6.3 并行For

#### 概念定义

**并行For**将循环迭代并行执行，适用于无依赖的迭代。

#### 完整示例

```go
package main

import (
 "fmt"
 "runtime"
 "sync"
 "sync/atomic"
 "time"
)

// 基本并行For
func parallelFor(start, end int, fn func(int)) {
 numCPU := runtime.NumCPU()
 chunkSize := (end - start + numCPU - 1) / numCPU

 var wg sync.WaitGroup

 for i := 0; i < numCPU; i++ {
  chunkStart := start + i*chunkSize
  chunkEnd := chunkStart + chunkSize
  if chunkEnd > end {
   chunkEnd = end
  }

  if chunkStart >= end {
   break
  }

  wg.Add(1)
  go func(s, e int) {
   defer wg.Done()
   for j := s; j < e; j++ {
    fn(j)
   }
  }(chunkStart, chunkEnd)
 }

 wg.Wait()
}

// 带动态调度的并行For
func parallelForDynamic(start, end int, fn func(int)) {
 var counter int64 = int64(start)

 numWorkers := runtime.NumCPU()
 var wg sync.WaitGroup

 for i := 0; i < numWorkers; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   for {
    idx := int(atomic.AddInt64(&counter, 1) - 1)
    if idx >= end {
     return
    }
    fn(idx)
   }
  }()
 }

 wg.Wait()
}

// 带原子累加的并行For
func parallelForSum(start, end int, fn func(int) int) int {
 var sum int64

 parallelFor(start, end, func(i int) {
  atomic.AddInt64(&sum, int64(fn(i)))
 })

 return int(sum)
}

// 矩阵乘法并行
func parallelMatrixMultiply(a, b [][]int) [][]int {
 rows := len(a)
 cols := len(b[0])
 inner := len(b)

 result := make([][]int, rows)
 for i := range result {
  result[i] = make([]int, cols)
 }

 parallelFor(0, rows, func(i int) {
  for j := 0; j < cols; j++ {
   sum := 0
   for k := 0; k < inner; k++ {
    sum += a[i][k] * b[k][j]
   }
   result[i][j] = sum
  }
 })

 return result
}

// 图像处理并行
func parallelImageProcess(width, height int, process func(x, y int)) {
 parallelFor(0, height, func(y int) {
  for x := 0; x < width; x++ {
   process(x, y)
  }
 })
}

// 使用示例
func forExample() {
 n := 10000000
 data := make([]int, n)

 // 串行填充
 start := time.Now()
 for i := 0; i < n; i++ {
  data[i] = i * i
 }
 serialTime := time.Since(start)

 // 并行填充
 data = make([]int, n)
 start = time.Now()
 parallelFor(0, n, func(i int) {
  data[i] = i * i
 })
 parallelTime := time.Since(start)

 fmt.Printf("Serial: %v\n", serialTime)
 fmt.Printf("Parallel: %v\n", parallelTime)
 fmt.Printf("Speedup: %.2fx\n", float64(serialTime)/float64(parallelTime))
}

// 蒙特卡洛模拟
func monteCarloPi(samples int) float64 {
 var inside int64

 parallelFor(0, samples, func(i int) {
  x := float64(i%10000) / 10000.0
  y := float64(i/10000%10000) / 10000.0

  // 使用更好的随机数
  // 这里简化处理
  if x*x+y*y <= 1 {
   atomic.AddInt64(&inside, 1)
  }
 })

 return 4.0 * float64(inside) / float64(samples)
}

func main() {
 fmt.Println("=== For Example ===")
 forExample()

 fmt.Println("\n=== Monte Carlo Pi ===")
 pi := monteCarloPi(1000000)
 fmt.Printf("Estimated Pi: %f\n", pi)
}
```

#### 最佳实践

1. **确保迭代间无依赖**
2. **使用动态调度处理不均匀任务**
3. **注意false sharing**
4. **使用原子操作或channel通信**

---

### 6.4 分治并行

#### 概念定义

**分治并行**将问题递归分解为子问题，并行解决后合并结果。

#### 完整示例

```go
package main

import (
 "fmt"
 "runtime"
 "sync"
)

// 并行快速排序
func parallelQuickSort(arr []int) {
 if len(arr) <= 1 {
  return
 }

 if len(arr) < 1000 {
  // 小数组串行排序
  quickSort(arr)
  return
 }

 pivot := partition(arr)

 var wg sync.WaitGroup
 wg.Add(2)

 go func() {
  defer wg.Done()
  parallelQuickSort(arr[:pivot])
 }()

 go func() {
  defer wg.Done()
  parallelQuickSort(arr[pivot+1:])
 }()

 wg.Wait()
}

func quickSort(arr []int) {
 if len(arr) <= 1 {
  return
 }

 pivot := partition(arr)
 quickSort(arr[:pivot])
 quickSort(arr[pivot+1:])
}

func partition(arr []int) int {
 pivot := arr[len(arr)-1]
 i := 0

 for j := 0; j < len(arr)-1; j++ {
  if arr[j] < pivot {
   arr[i], arr[j] = arr[j], arr[i]
   i++
  }
 }

 arr[i], arr[len(arr)-1] = arr[len(arr)-1], arr[i]
 return i
}

// 并行归并排序
func parallelMergeSort(arr []int) []int {
 if len(arr) <= 1 {
  return arr
 }

 if len(arr) < 1000 {
  return mergeSort(arr)
 }

 mid := len(arr) / 2

 var left, right []int
 var wg sync.WaitGroup
 wg.Add(2)

 go func() {
  defer wg.Done()
  left = parallelMergeSort(arr[:mid])
 }()

 go func() {
  defer wg.Done()
  right = parallelMergeSort(arr[mid:])
 }()

 wg.Wait()

 return merge(left, right)
}

func mergeSort(arr []int) []int {
 if len(arr) <= 1 {
  return arr
 }

 mid := len(arr) / 2
 left := mergeSort(arr[:mid])
 right := mergeSort(arr[mid:])

 return merge(left, right)
}

func merge(left, right []int) []int {
 result := make([]int, 0, len(left)+len(right))

 i, j := 0, 0
 for i < len(left) && j < len(right) {
  if left[i] <= right[j] {
   result = append(result, left[i])
   i++
  } else {
   result = append(result, right[j])
   j++
  }
 }

 result = append(result, left[i:]...)
 result = append(result, right[j:]...)

 return result
}

// 并行树遍历
type TreeNode struct {
 Value int
 Left  *TreeNode
 Right *TreeNode
}

func parallelTreeSum(root *TreeNode) int {
 if root == nil {
  return 0
 }

 if root.Left == nil && root.Right == nil {
  return root.Value
 }

 var leftSum, rightSum int
 var wg sync.WaitGroup
 wg.Add(2)

 go func() {
  defer wg.Done()
  leftSum = parallelTreeSum(root.Left)
 }()

 go func() {
  defer wg.Done()
  rightSum = parallelTreeSum(root.Right)
 }()

 wg.Wait()

 return root.Value + leftSum + rightSum
}

// 并行矩阵分块乘法
func parallelBlockMatrixMultiply(a, b [][]int, blockSize int) [][]int {
 n := len(a)
 result := make([][]int, n)
 for i := range result {
  result[i] = make([]int, n)
 }

 var wg sync.WaitGroup

 for i := 0; i < n; i += blockSize {
  for j := 0; j < n; j += blockSize {
   for k := 0; k < n; k += blockSize {
    wg.Add(1)
    go func(i0, j0, k0 int) {
     defer wg.Done()

     for i := i0; i < min(i0+blockSize, n); i++ {
      for j := j0; j < min(j0+blockSize, n); j++ {
       for k := k0; k < min(k0+blockSize, n); k++ {
        result[i][j] += a[i][k] * b[k][j]
       }
      }
     }
    }(i, j, k)
   }
  }
 }

 wg.Wait()
 return result
}

func min(a, b int) int {
 if a < b {
  return a
 }
 return b
}

// 使用示例
func quickSortExample() {
 arr := make([]int, 1000000)
 for i := range arr {
  arr[i] = len(arr) - i
 }

 parallelQuickSort(arr)

 // 验证
 for i := 1; i < len(arr); i++ {
  if arr[i] < arr[i-1] {
   fmt.Println("Sort failed!")
   return
  }
 }
 fmt.Println("Sort successful!")
}

// 斐波那契并行计算
func parallelFibonacci(n int) int {
 if n <= 1 {
  return n
 }

 if n < 20 {
  return fibonacci(n)
 }

 var left, right int
 var wg sync.WaitGroup
 wg.Add(2)

 go func() {
  defer wg.Done()
  left = parallelFibonacci(n - 1)
 }()

 go func() {
  defer wg.Done()
  right = parallelFibonacci(n - 2)
 }()

 wg.Wait()

 return left + right
}

func fibonacci(n int) int {
 if n <= 1 {
  return n
 }
 return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
 fmt.Println("=== Quick Sort Example ===")
 quickSortExample()

 fmt.Println("\n=== Parallel Fibonacci ===")
 result := parallelFibonacci(30)
 fmt.Printf("Fibonacci(30) = %d\n", result)
}
```

#### 最佳实践

1. **设置递归终止阈值**
2. **避免过度并行化**
3. **注意任务粒度平衡**
4. **考虑内存局部性**
5. **使用WaitGroup等待子任务**

---

## 7. 异步编程模式

### 7.1 Future/Promise模式

#### 概念定义

**Future/Promise**模式表示异步计算的结果，Future用于获取结果，Promise用于设置结果。

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// Future接口
type Future[T any] interface {
 Get() (T, error)
 GetWithTimeout(timeout time.Duration) (T, error)
 IsDone() bool
}

// Promise接口
type Promise[T any] interface {
 Future[T]
 Complete(value T, err error)
}

// 实现
type futureImpl[T any] struct {
 done    chan struct{}
 value   T
 err     error
 mu      sync.RWMutex
}

func NewFuture[T any]() Promise[T] {
 return &futureImpl[T]{
  done: make(chan struct{}),
 }
}

func (f *futureImpl[T]) Complete(value T, err error) {
 f.mu.Lock()
 defer f.mu.Unlock()

 select {
 case <-f.done:
  panic("future already completed")
 default:
  f.value = value
  f.err = err
  close(f.done)
 }
}

func (f *futureImpl[T]) Get() (T, error) {
 <-f.done
 return f.value, f.err
}

func (f *futureImpl[T]) GetWithTimeout(timeout time.Duration) (T, error) {
 select {
 case <-f.done:
  return f.value, f.err
 case <-time.After(timeout):
  var zero T
  return zero, fmt.Errorf("timeout")
 }
}

func (f *futureImpl[T]) IsDone() bool {
 select {
 case <-f.done:
  return true
 default:
  return false
 }
}

// 异步执行
func Async[T any](fn func() (T, error)) Future[T] {
 future := NewFuture[T]()

 go func() {
  value, err := fn()
  future.Complete(value, err)
 }()

 return future
}

// 带Context的异步执行
func AsyncWithContext[T any](ctx context.Context, fn func() (T, error)) Future[T] {
 future := NewFuture[T]()

 go func() {
  done := make(chan struct {
   value T
   err   error
  }, 1)

  go func() {
   value, err := fn()
   done <- struct {
    value T
    err   error
   }{value, err}
  }()

  select {
  case result := <-done:
   future.Complete(result.value, result.err)
  case <-ctx.Done():
   var zero T
   future.Complete(zero, ctx.Err())
  }
 }()

 return future
}

// Future组合
func Then[T any, R any](future Future[T], fn func(T) (R, error)) Future[R] {
 result := NewFuture[R]()

 go func() {
  value, err := future.Get()
  if err != nil {
   var zero R
   result.Complete(zero, err)
   return
  }

  newValue, newErr := fn(value)
  result.Complete(newValue, newErr)
 }()

 return result
}

// 等待多个Future
func All[T any](futures ...Future[T]) Future[[]T] {
 result := NewFuture[[]T]()

 go func() {
  values := make([]T, len(futures))
  var firstErr error

  for i, f := range futures {
   value, err := f.Get()
   if err != nil && firstErr == nil {
    firstErr = err
   }
   values[i] = value
  }

  result.Complete(values, firstErr)
 }()

 return result
}

// 任一Future完成
func Any[T any](futures ...Future[T]) Future[T] {
 result := NewFuture[T]()

 go func() {
  done := make(chan struct {
   value T
   err   error
  }, len(futures))

  for _, f := range futures {
   go func(future Future[T]) {
    value, err := future.Get()
    done <- struct {
     value T
     err   error
    }{value, err}
   }(f)
  }

  r := <-done
  result.Complete(r.value, r.err)
 }()

 return result
}

// 使用示例
func futureExample() {
 // 基本使用
 future := Async(func() (int, error) {
  time.Sleep(100 * time.Millisecond)
  return 42, nil
 })

 fmt.Println("Waiting for result...")
 result, err := future.Get()
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 } else {
  fmt.Printf("Result: %d\n", result)
 }
}

// 链式调用
func chainExample() {
 future := Async(func() (int, error) {
  return 10, nil
 })

 chained := Then(future, func(n int) (string, error) {
  return fmt.Sprintf("Result: %d", n*2), nil
 })

 result, err := chained.Get()
 fmt.Printf("Chained result: %s, err: %v\n", result, err)
}

// 等待多个
func allExample() {
 futures := []Future[int]{
  Async(func() (int, error) {
   time.Sleep(100 * time.Millisecond)
   return 1, nil
  }),
  Async(func() (int, error) {
   time.Sleep(200 * time.Millisecond)
   return 2, nil
  }),
  Async(func() (int, error) {
   time.Sleep(150 * time.Millisecond)
   return 3, nil
  }),
 }

 allFuture := All(futures...)
 results, err := allFuture.Get()
 fmt.Printf("All results: %v, err: %v\n", results, err)
}

// 并行下载示例
func parallelDownloadExample() {
 urls := []string{
  "https://example.com/1",
  "https://example.com/2",
  "https://example.com/3",
 }

 var futures []Future[string]
 for _, url := range urls {
  u := url
  future := Async(func() (string, error) {
   // 模拟下载
   time.Sleep(time.Duration(100+len(u)*10) * time.Millisecond)
   return "Content of " + u, nil
  })
  futures = append(futures, future)
 }

 results, err := All(futures...).Get()
 if err != nil {
  fmt.Printf("Error: %v\n", err)
  return
 }

 for i, content := range results {
  fmt.Printf("Downloaded %s: %s\n", urls[i], content)
 }
}

func main() {
 fmt.Println("=== Future Example ===")
 futureExample()

 fmt.Println("\n=== Chain Example ===")
 chainExample()

 fmt.Println("\n=== All Example ===")
 allExample()

 fmt.Println("\n=== Parallel Download Example ===")
 parallelDownloadExample()
}
```

#### 最佳实践

1. **Future只能完成一次**
2. **提供超时获取方法**
3. **支持链式组合**
4. **正确处理错误传播**
5. **使用泛型提高类型安全**

---

### 7.2 Callback模式

#### 概念定义

**Callback模式**在异步操作完成时调用指定的回调函数处理结果。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 回调类型
type SuccessCallback[T any] func(T)
type ErrorCallback func(error)
type CompleteCallback func()

// 异步任务
type AsyncTask[T any] struct {
 onSuccess SuccessCallback[T]
 onError   ErrorCallback
 onComplete CompleteCallback
}

func NewAsyncTask[T any]() *AsyncTask[T] {
 return &AsyncTask[T]{}
}

func (t *AsyncTask[T]) OnSuccess(cb SuccessCallback[T]) *AsyncTask[T] {
 t.onSuccess = cb
 return t
}

func (t *AsyncTask[T]) OnError(cb ErrorCallback) *AsyncTask[T] {
 t.onError = cb
 return t
}

func (t *AsyncTask[T]) OnComplete(cb CompleteCallback) *AsyncTask[T] {
 t.onComplete = cb
 return t
}

func (t *AsyncTask[T]) Execute(fn func() (T, error)) {
 go func() {
  result, err := fn()

  if err != nil {
   if t.onError != nil {
    t.onError(err)
   }
  } else {
   if t.onSuccess != nil {
    t.onSuccess(result)
   }
  }

  if t.onComplete != nil {
   t.onComplete()
  }
 }()
}

// 事件驱动回调
type EventEmitter struct {
 listeners map[string][]func(...interface{})
 mu        sync.RWMutex
}

func NewEventEmitter() *EventEmitter {
 return &EventEmitter{
  listeners: make(map[string][]func(...interface{})),
 }
}

func (e *EventEmitter) On(event string, listener func(...interface{})) {
 e.mu.Lock()
 defer e.mu.Unlock()
 e.listeners[event] = append(e.listeners[event], listener)
}

func (e *EventEmitter) Off(event string, listener func(...interface{})) {
 e.mu.Lock()
 defer e.mu.Unlock()

 listeners := e.listeners[event]
 for i, l := range listeners {
  if &l == &listener {
   e.listeners[event] = append(listeners[:i], listeners[i+1:]...)
   break
  }
 }
}

func (e *EventEmitter) Emit(event string, args ...interface{}) {
 e.mu.RLock()
 listeners := e.listeners[event]
 e.mu.RUnlock()

 for _, listener := range listeners {
  go listener(args...)
 }
}

// Promise风格的回调链
type CallbackPromise[T any] struct {
 thenCbs []func(T) (T, error)
 catchCb func(error) error
 finallyCb func()
}

func NewCallbackPromise[T any]() *CallbackPromise[T] {
 return &CallbackPromise[T]{
  thenCbs: make([]func(T) (T, error), 0),
 }
}

func (p *CallbackPromise[T]) Then(fn func(T) (T, error)) *CallbackPromise[T] {
 p.thenCbs = append(p.thenCbs, fn)
 return p
}

func (p *CallbackPromise[T]) Catch(fn func(error) error) *CallbackPromise[T] {
 p.catchCb = fn
 return p
}

func (p *CallbackPromise[T]) Finally(fn func()) *CallbackPromise[T] {
 p.finallyCb = fn
 return p
}

func (p *CallbackPromise[T]) Resolve(value T, err error) {
 if err != nil {
  if p.catchCb != nil {
   err = p.catchCb(err)
  }
  if p.finallyCb != nil {
   p.finallyCb()
  }
  return
 }

 for _, fn := range p.thenCbs {
  value, err = fn(value)
  if err != nil {
   if p.catchCb != nil {
    p.catchCb(err)
   }
   break
  }
 }

 if p.finallyCb != nil {
  p.finallyCb()
 }
}

// 使用示例
func callbackExample() {
 NewAsyncTask[int]().
  OnSuccess(func(result int) {
   fmt.Printf("Success: %d\n", result)
  }).
  OnError(func(err error) {
   fmt.Printf("Error: %v\n", err)
  }).
  OnComplete(func() {
   fmt.Println("Completed")
  }).
  Execute(func() (int, error) {
   time.Sleep(100 * time.Millisecond)
   return 42, nil
  })

 time.Sleep(200 * time.Millisecond)
}

// 事件驱动示例
func eventDrivenExample() {
 emitter := NewEventEmitter()

 emitter.On("user:login", func(args ...interface{}) {
  username := args[0].(string)
  fmt.Printf("User logged in: %s\n", username)
 })

 emitter.On("user:logout", func(args ...interface{}) {
  username := args[0].(string)
  fmt.Printf("User logged out: %s\n", username)
 })

 emitter.Emit("user:login", "alice")
 emitter.Emit("user:logout", "alice")

 time.Sleep(100 * time.Millisecond)
}

// 回调地狱解决方案
func callbackHellSolution() {
 // 链式回调
 promise := NewCallbackPromise[int]()

 promise.
  Then(func(n int) (int, error) {
   fmt.Printf("Step 1: %d\n", n)
   return n * 2, nil
  }).
  Then(func(n int) (int, error) {
   fmt.Printf("Step 2: %d\n", n)
   return n + 10, nil
  }).
  Then(func(n int) (int, error) {
   fmt.Printf("Step 3: %d\n", n)
   return n / 2, nil
  }).
  Catch(func(err error) error {
   fmt.Printf("Error: %v\n", err)
   return err
  }).
  Finally(func() {
   fmt.Println("Done")
  })

 // 模拟异步操作
 go func() {
  time.Sleep(100 * time.Millisecond)
  promise.Resolve(5, nil)
 }()

 time.Sleep(200 * time.Millisecond)
}

func main() {
 fmt.Println("=== Callback Example ===")
 callbackExample()

 fmt.Println("\n=== Event Driven Example ===")
 eventDrivenExample()

 fmt.Println("\n=== Callback Hell Solution ===")
 callbackHellSolution()
}
```

#### 最佳实践

1. **避免回调地狱，使用链式API**
2. **始终处理错误回调**
3. **使用事件驱动解耦**
4. **注意回调执行顺序**
5. **考虑使用channel替代回调**

---

### 7.3 Async/Await模拟

#### 概念定义

**Async/Await**模式使用同步风格的代码编写异步逻辑，Go可以通过channel和select模拟。

#### 完整示例

```go
package main

import (
 "fmt"
 "time"
)

// Async函数返回channel
type AsyncResult[T any] struct {
 Value T
 Error error
}

func AsyncFn[T any](fn func() (T, error)) <-chan AsyncResult[T] {
 ch := make(chan AsyncResult[T], 1)

 go func() {
  value, err := fn()
  ch <- AsyncResult[T]{Value: value, Error: err}
  close(ch)
 }()

 return ch
}

// Await等待结果
func Await[T any](ch <-chan AsyncResult[T]) (T, error) {
 result := <-ch
 return result.Value, result.Error
}

// AwaitWithTimeout
func AwaitWithTimeout[T any](ch <-chan AsyncResult[T], timeout time.Duration) (T, error) {
 select {
 case result := <-ch:
  return result.Value, result.Error
 case <-time.After(timeout):
  var zero T
  return zero, fmt.Errorf("timeout")
 }
}

// 使用示例
func asyncAwaitExample() {
 // 模拟异步操作
 fetchUser := func(id int) <-chan AsyncResult[string] {
  return AsyncFn(func() (string, error) {
   time.Sleep(100 * time.Millisecond)
   return fmt.Sprintf("User-%d", id), nil
  })
 }

 fetchOrders := func(user string) <-chan AsyncResult[[]string] {
  return AsyncFn(func() ([]string, error) {
   time.Sleep(100 * time.Millisecond)
   return []string{"order1", "order2"}, nil
  })
 }

 // 顺序执行（类似await）
 userCh := fetchUser(1)
 user, err := Await(userCh)
 if err != nil {
  fmt.Printf("Error: %v\n", err)
  return
 }
 fmt.Printf("User: %s\n", user)

 ordersCh := fetchOrders(user)
 orders, err := Await(ordersCh)
 if err != nil {
  fmt.Printf("Error: %v\n", err)
  return
 }
 fmt.Printf("Orders: %v\n", orders)
}

// 并行Await
func parallelAsyncAwait() {
 task1 := AsyncFn(func() (int, error) {
  time.Sleep(100 * time.Millisecond)
  return 1, nil
 })

 task2 := AsyncFn(func() (int, error) {
  time.Sleep(200 * time.Millisecond)
  return 2, nil
 })

 task3 := AsyncFn(func() (int, error) {
  time.Sleep(150 * time.Millisecond)
  return 3, nil
 })

 // 使用select等待任一完成
 for i := 0; i < 3; i++ {
  select {
  case result := <-task1:
   fmt.Printf("Task 1: %v, err: %v\n", result.Value, result.Error)
   task1 = nil
  case result := <-task2:
   fmt.Printf("Task 2: %v, err: %v\n", result.Value, result.Error)
   task2 = nil
  case result := <-task3:
   fmt.Printf("Task 3: %v, err: %v\n", result.Value, result.Error)
   task3 = nil
  }
 }
}

// Promise风格的Async/Await
type Promise[T any] struct {
 ch <-chan AsyncResult[T]
}

func NewPromise[T any](fn func() (T, error)) *Promise[T] {
 return &Promise[T]{ch: AsyncFn(fn)}
}

func (p *Promise[T]) Await() (T, error) {
 return Await(p.ch)
}

func (p *Promise[T]) Then(fn func(T) (T, error)) *Promise[T] {
 newCh := make(chan AsyncResult[T], 1)

 go func() {
  result, err := Await(p.ch)
  if err != nil {
   newCh <- AsyncResult[T]{Error: err}
  } else {
   newValue, newErr := fn(result)
   newCh <- AsyncResult[T]{Value: newValue, Error: newErr}
  }
  close(newCh)
 }()

 return &Promise[T]{ch: newCh}
}

// 使用Promise
func promiseExample() {
 result, err := NewPromise(func() (int, error) {
  return 10, nil
 }).Then(func(n int) (int, error) {
  return n * 2, nil
 }).Then(func(n int) (int, error) {
  return n + 5, nil
 }).Await()

 fmt.Printf("Result: %d, err: %v\n", result, err)
}

// 协程池 + Async/Await
type TaskPool struct {
 tasks   chan func()
 workers int
}

func NewTaskPool(workers int) *TaskPool {
 pool := &TaskPool{
  tasks:   make(chan func()),
  workers: workers,
 }

 for i := 0; i < workers; i++ {
  go pool.worker()
 }

 return pool
}

func (p *TaskPool) worker() {
 for task := range p.tasks {
  task()
 }
}

func (p *TaskPool) Submit(fn func()) {
 p.tasks <- fn
}

func (p *TaskPool) Async[T any](fn func() (T, error)) <-chan AsyncResult[T] {
 ch := make(chan AsyncResult[T], 1)

 p.Submit(func() {
  value, err := fn()
  ch <- AsyncResult[T]{Value: value, Error: err}
  close(ch)
 })

 return ch
}

// 使用TaskPool
func taskPoolExample() {
 pool := NewTaskPool(4)

 var results []<-chan AsyncResult[int]
 for i := 0; i < 10; i++ {
  n := i
  ch := pool.Async(func() (int, error) {
   time.Sleep(50 * time.Millisecond)
   return n * n, nil
  })
  results = append(results, ch)
 }

 for i, ch := range results {
  result, _ := Await(ch)
  fmt.Printf("Task %d result: %d\n", i, result)
 }
}

func main() {
 fmt.Println("=== Async/Await Example ===")
 asyncAwaitExample()

 fmt.Println("\n=== Parallel Async/Await ===")
 parallelAsyncAwait()

 fmt.Println("\n=== Promise Example ===")
 promiseExample()

 fmt.Println("\n=== Task Pool Example ===")
 taskPoolExample()
}
```

#### 最佳实践

1. **使用channel模拟async/await**
2. **提供超时版本**
3. **支持链式调用**
4. **考虑使用协程池**
5. **注意goroutine泄漏**

---

### 7.4 事件驱动模式

#### 概念定义

**事件驱动模式**通过事件的生产和消费实现松耦合的异步通信。

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// 事件定义
type Event struct {
 Type      string
 Payload   interface{}
 Timestamp time.Time
}

// 事件处理器
type EventHandler func(Event)

// 事件总线
type EventBus struct {
 handlers map[string][]EventHandler
 ch       chan Event
 ctx      context.Context
 cancel   context.CancelFunc
 wg       sync.WaitGroup
}

func NewEventBus(bufferSize int) *EventBus {
 ctx, cancel := context.WithCancel(context.Background())

 bus := &EventBus{
  handlers: make(map[string][]EventHandler),
  ch:       make(chan Event, bufferSize),
  ctx:      ctx,
  cancel:   cancel,
 }

 bus.start()
 return bus
}

func (b *EventBus) start() {
 b.wg.Add(1)
 go func() {
  defer b.wg.Done()

  for {
   select {
   case event := <-b.ch:
    dispatch(b.handlers, event)
   case <-b.ctx.Done():
    // 处理剩余事件
    for {
     select {
     case event := <-b.ch:
      dispatch(b.handlers, event)
     default:
      return
     }
    }
   }
  }
 }()
}

func dispatch(handlers map[string][]EventHandler, event Event) {
 hs, ok := handlers[event.Type]
 if !ok {
  return
 }

 for _, h := range hs {
  go h(event) // 并发处理
 }
}

func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
 b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(event Event) {
 select {
 case b.ch <- event:
 case <-b.ctx.Done():
 }
}

func (b *EventBus) Stop() {
 b.cancel()
 b.wg.Wait()
}

// 带优先级的事件
type PriorityEvent struct {
 Event
 Priority int
}

type PriorityEventBus struct {
 highPriority   chan PriorityEvent
 mediumPriority chan PriorityEvent
 lowPriority    chan PriorityEvent
 handlers       map[string][]EventHandler
}

func NewPriorityEventBus() *PriorityEventBus {
 return &PriorityEventBus{
  highPriority:   make(chan PriorityEvent, 100),
  mediumPriority: make(chan PriorityEvent, 100),
  lowPriority:    make(chan PriorityEvent, 100),
  handlers:       make(map[string][]EventHandler),
 }
}

func (b *PriorityEventBus) Publish(event PriorityEvent) {
 switch event.Priority {
 case 1:
  b.highPriority <- event
 case 2:
  b.mediumPriority <- event
 default:
  b.lowPriority <- event
 }
}

func (b *PriorityEventBus) Start() {
 go func() {
  for {
   select {
   case event := <-b.highPriority:
    b.dispatch(event.Event)
   default:
    select {
    case event := <-b.highPriority:
     b.dispatch(event.Event)
    case event := <-b.mediumPriority:
     b.dispatch(event.Event)
    default:
     select {
     case event := <-b.highPriority:
      b.dispatch(event.Event)
     case event := <-b.mediumPriority:
      b.dispatch(event.Event)
     case event := <-b.lowPriority:
      b.dispatch(event.Event)
     }
    }
   }
  }
 }()
}

func (b *PriorityEventBus) dispatch(event Event) {
 if handlers, ok := b.handlers[event.Type]; ok {
  for _, h := range handlers {
   go h(event)
  }
 }
}

// 使用示例
func eventBusExample() {
 bus := NewEventBus(100)

 // 订阅事件
 bus.Subscribe("user:created", func(e Event) {
  fmt.Printf("[Handler 1] User created: %v\n", e.Payload)
 })

 bus.Subscribe("user:created", func(e Event) {
  fmt.Printf("[Handler 2] User created: %v\n", e.Payload)
 })

 bus.Subscribe("order:placed", func(e Event) {
  fmt.Printf("Order placed: %v\n", e.Payload)
 })

 // 发布事件
 bus.Publish(Event{
  Type:      "user:created",
  Payload:   map[string]string{"id": "1", "name": "Alice"},
  Timestamp: time.Now(),
 })

 bus.Publish(Event{
  Type:      "order:placed",
  Payload:   map[string]interface{}{"id": "100", "amount": 99.99},
  Timestamp: time.Now(),
 })

 time.Sleep(100 * time.Millisecond)
 bus.Stop()
}

// CQRS风格的事件溯源
type Command interface {
 Type() string
}

type Event2 interface {
 Type() string
 AggregateID() string
}

type CommandHandler func(Command) ([]Event2, error)

type EventStore struct {
 events map[string][]Event2
 mu     sync.RWMutex
}

func NewEventStore() *EventStore {
 return &EventStore{
  events: make(map[string][]Event2),
 }
}

func (s *EventStore) Append(events []Event2) {
 s.mu.Lock()
 defer s.mu.Unlock()

 for _, e := range events {
  id := e.AggregateID()
  s.events[id] = append(s.events[id], e)
 }
}

func (s *EventStore) Get(aggregateID string) []Event2 {
 s.mu.RLock()
 defer s.mu.RUnlock()
 return s.events[aggregateID]
}

// Saga模式
type Saga struct {
 steps []SagaStep
}

type SagaStep struct {
 Action   func() error
 Compensate func()
}

func (s *Saga) Execute() error {
 completed := make([]int, 0)

 for i, step := range s.steps {
  if err := step.Action(); err != nil {
   // 补偿已完成的步骤
   for j := len(completed) - 1; j >= 0; j-- {
    s.steps[completed[j]].Compensate()
   }
   return err
  }
  completed = append(completed, i)
 }

 return nil
}

// 使用示例
func sagaExample() {
 saga := &Saga{
  steps: []SagaStep{
   {
    Action: func() error {
     fmt.Println("Step 1: Reserve inventory")
     return nil
    },
    Compensate: func() {
     fmt.Println("Compensate Step 1: Release inventory")
    },
   },
   {
    Action: func() error {
     fmt.Println("Step 2: Process payment")
     return nil
    },
    Compensate: func() {
     fmt.Println("Compensate Step 2: Refund payment")
    },
   },
   {
    Action: func() error {
     fmt.Println("Step 3: Ship order")
     return nil
    },
    Compensate: func() {
     fmt.Println("Compensate Step 3: Cancel shipment")
    },
   },
  },
 }

 if err := saga.Execute(); err != nil {
  fmt.Printf("Saga failed: %v\n", err)
 } else {
  fmt.Println("Saga completed successfully")
 }
}

func main() {
 fmt.Println("=== Event Bus Example ===")
 eventBusExample()

 fmt.Println("\n=== Saga Example ===")
 sagaExample()
}
```

#### 最佳实践

1. **使用channel缓冲事件**
2. **支持优雅关闭**
3. **考虑事件优先级**
4. **实现Saga模式处理长事务**
5. **使用CQRS分离读写**

---

## 8. 并发安全

### 8.1 数据竞争检测

#### 概念定义

**数据竞争**（Data Race）指两个或多个Goroutine同时访问同一内存位置，且至少一个是写操作，没有同步机制。

#### 工作原理

```
数据竞争示例:

Goroutine 1          Goroutine 2
    │                      │
    │ 读取 x               │ 写入 x
    │  (R)                 │  (W)
    ↓                      ↓
┌─────────┐            ┌─────────┐
│  x = ?  │            │  x = 5  │
└─────────┘            └─────────┘

结果不确定！

Go Race Detector:
- 基于Happens-Before关系检测
- 使用Shadow Memory跟踪访问
- 运行时开销约5-10x
```

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "sync/atomic"
)

// 数据竞争示例
func dataRaceExample() {
 var counter int
 var wg sync.WaitGroup

 for i := 0; i < 1000; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter++ // 数据竞争！
  }()
 }

 wg.Wait()
 fmt.Printf("Counter: %d (expected 1000)\n", counter)
}

// 修复：使用Mutex
func fixedWithMutex() {
 var counter int
 var mu sync.Mutex
 var wg sync.WaitGroup

 for i := 0; i < 1000; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   mu.Lock()
   counter++
   mu.Unlock()
  }()
 }

 wg.Wait()
 fmt.Printf("Counter (Mutex): %d\n", counter)
}

// 修复：使用Atomic
func fixedWithAtomic() {
 var counter int64
 var wg sync.WaitGroup

 for i := 0; i < 1000; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   atomic.AddInt64(&counter, 1)
  }()
 }

 wg.Wait()
 fmt.Printf("Counter (Atomic): %d\n", counter)
}

// 修复：使用Channel
func fixedWithChannel() {
 counterCh := make(chan int, 1)
 counterCh <- 0

 var wg sync.WaitGroup
 for i := 0; i < 1000; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter := <-counterCh
   counter++
   counterCh <- counter
  }()
 }

 wg.Wait()
 counter := <-counterCh
 fmt.Printf("Counter (Channel): %d\n", counter)
}

// 复杂数据竞争
type Counter struct {
 value int
}

func (c *Counter) Increment() {
 c.value++
}

func (c *Counter) Value() int {
 return c.value
}

func complexDataRace() {
 counter := &Counter{}
 var wg sync.WaitGroup

 // 写操作
 for i := 0; i < 500; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter.Increment()
  }()
 }

 // 读操作
 for i := 0; i < 500; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   _ = counter.Value()
  }()
 }

 wg.Wait()
 fmt.Printf("Counter value: %d\n", counter.Value())
}

// 修复后的Counter
type SafeCounter struct {
 mu    sync.RWMutex
 value int
}

func (c *SafeCounter) Increment() {
 c.mu.Lock()
 defer c.mu.Unlock()
 c.value++
}

func (c *SafeCounter) Value() int {
 c.mu.RLock()
 defer c.mu.RUnlock()
 return c.value
}

// 检测工具使用
func raceDetectionTools() {
 // 1. 使用 -race 标志编译
 // go run -race main.go

 // 2. 使用 go test -race
 // go test -race ./...

 // 3. 使用 go build -race
 // go build -race -o myapp

 fmt.Println("Use 'go run -race' to detect data races")
}

// 常见的数据竞争场景
func commonRaceScenarios() {
 // 场景1: 循环变量捕获
 fmt.Println("Scenario 1: Loop variable capture")
 for i := 0; i < 10; i++ {
  go func() {
   // 错误：i被多个goroutine共享
   // fmt.Println(i)
  }()
 }

 // 正确
 for i := 0; i < 10; i++ {
  i := i // 创建副本
  go func() {
   fmt.Println(i)
  }()
 }

 // 场景2: 切片/Map并发访问
 fmt.Println("\nScenario 2: Concurrent map access")
 m := make(map[string]int)

 // 错误
 // go func() { m["key"] = 1 }()
 // go func() { _ = m["key"] }()

 // 正确：使用sync.Map
 var safeMap sync.Map
 go func() { safeMap.Store("key", 1) }()
 go func() { safeMap.Load("key") }()

 // 场景3: 闭包修改外部变量
 fmt.Println("\nScenario 3: Closure modifying outer variable")
 var sum int

 // 错误
 // for i := 0; i < 10; i++ {
 //     go func() { sum += i }()
 // }

 // 正确
 var mu sync.Mutex
 for i := 0; i < 10; i++ {
  i := i
  go func() {
   mu.Lock()
   sum += i
   mu.Unlock()
  }()
 }
}

func main() {
 fmt.Println("=== Data Race Example ===")
 dataRaceExample()

 fmt.Println("\n=== Fixed with Mutex ===")
 fixedWithMutex()

 fmt.Println("\n=== Fixed with Atomic ===")
 fixedWithAtomic()

 fmt.Println("\n=== Fixed with Channel ===")
 fixedWithChannel()

 fmt.Println("\n=== Common Race Scenarios ===")
 commonRaceScenarios()
}
```

#### 反例说明

```go
// ❌ 错误：无保护的共享变量
func unprotectedSharedVar() {
 var x int
 go func() { x = 1 }()
 go func() { x = 2 }()
}

// ❌ 错误：读写不同步
func readWriteUnsync() {
 var x int
 go func() { x = 1 }()
 _ = x // 读操作与写操作无同步
}

// ❌ 错误：多个goroutine写map
func concurrentMapWrite() {
 m := make(map[int]int)
 for i := 0; i < 10; i++ {
  go func(n int) {
   m[n] = n // 数据竞争！
  }(i)
 }
}

// ❌ 错误：WaitGroup拷贝
func waitGroupCopy() {
 var wg sync.WaitGroup

 worker := func(wg sync.WaitGroup) { // 拷贝
  defer wg.Done()
 }

 wg.Add(1)
 go worker(wg)
 wg.Wait()
}

// ✅ 正确：使用指针
func correctWaitGroup() {
 var wg sync.WaitGroup

 worker := func(wg *sync.WaitGroup) {
  defer wg.Done()
 }

 wg.Add(1)
 go worker(&wg)
 wg.Wait()
}
```

#### 最佳实践

1. **始终使用-race标志测试**
2. **优先使用channel通信**
3. **使用sync包保护共享状态**
4. **避免闭包捕获循环变量**
5. **使用sync.Map替代普通map**

---

### 8.2 Happens-Before关系

#### 概念定义

**Happens-Before**定义了内存操作的偏序关系，确保一个操作的结果对后续操作可见。

#### 工作原理

```
Happens-Before规则:

1. 程序内顺序:
   同一goroutine内，语句按程序顺序执行

   a = 1
   b = 2  // a=1 happens-before b=2

2. Channel通信:
   - 发送 happens-before 接收完成
   - 关闭 happens-before 接收返回零值

   ch <- v    // 发送
   ...
   <-ch       // 接收，保证看到v

3. Mutex:
   - Unlock happens-before 后续的Lock

   mu.Unlock()  // 之前的写入可见
   ...
   mu.Lock()    // 保证看到之前的写入

4. WaitGroup:
   - Wait happens-before Wait返回

   wg.Done()
   ...
   wg.Wait()    // 保证看到Done前的写入

5. Once:
   - Do中的函数 happens-before 任何Do返回

   once.Do(f)
   ...
   // f() happens-before 这里
```

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
)

// Channel Happens-Before
func channelHappensBefore() {
 var msg string
 done := make(chan struct{})

 go func() {
  msg = "hello" // 写入
  close(done)   // 发送信号
 }()

 <-done          // 接收信号
 fmt.Println(msg) // 保证看到"hello"
}

// Mutex Happens-Before
func mutexHappensBefore() {
 var counter int
 var mu sync.Mutex
 var wg sync.WaitGroup

 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   mu.Lock()
   counter++ // 写入
   mu.Unlock() // happens-before后续的Lock
  }()
 }

 wg.Wait()
 fmt.Println(counter) // 保证看到所有写入
}

// WaitGroup Happens-Before
func waitGroupHappensBefore() {
 var result string
 var wg sync.WaitGroup

 wg.Add(1)
 go func() {
  defer wg.Done()
  result = "computed" // 写入
 }()

 wg.Wait() // happens-before这里
 fmt.Println(result) // 保证看到"computed"
}

// Once Happens-Before
func onceHappensBefore() {
 var config string
 var once sync.Once

 once.Do(func() {
  config = "initialized" // 写入
 })

 // config的写入 happens-before这里
 fmt.Println(config) // 保证看到"initialized"
}

// 没有Happens-Before的情况
func noHappensBefore() {
 var x, y int

 go func() {
  x = 1
  if y == 0 {
   // x=1 不一定对另一个goroutine可见
  }
 }()

 go func() {
  y = 1
  if x == 0 {
   // y=1 不一定对另一个goroutine可见
  }
 }()

 // 可能两个都执行！
}

// 建立Happens-Before
func establishHappensBefore() {
 var x, y int
 done := make(chan struct{})

 go func() {
  x = 1
  done <- struct{}{} // 发送
 }()

 go func() {
  <-done // 接收，建立happens-before
  if x == 1 {
   // 保证看到x=1
   y = 1
  }
 }()
}

// 复杂的Happens-Before链
func happensBeforeChain() {
 var a, b, c int
 ch1 := make(chan struct{})
 ch2 := make(chan struct{})

 go func() {
  a = 1
  close(ch1) // a=1 happens-before ch1关闭
 }()

 go func() {
  <-ch1       // ch1关闭 happens-before这里
  b = 1
  close(ch2) // b=1 happens-before ch2关闭
 }()

 go func() {
  <-ch2       // ch2关闭 happens-before这里
  c = 1
 }()

 // a=1 happens-before b=1 happens-before c=1
}

// 内存重排序示例
func memoryReordering() {
 var a, b int
 var wg sync.WaitGroup
 wg.Add(2)

 go func() {
  defer wg.Done()
  a = 1
  if b == 0 {
   // 可能被重排序，先检查b再写入a
  }
 }()

 go func() {
  defer wg.Done()
  b = 1
  if a == 0 {
   // 可能被重排序，先检查a再写入b
  }
 }()

 wg.Wait()
 // 可能两个条件都满足！
}

func main() {
 fmt.Println("=== Channel Happens-Before ===")
 channelHappensBefore()

 fmt.Println("\n=== Mutex Happens-Before ===")
 mutexHappensBefore()

 fmt.Println("\n=== WaitGroup Happens-Before ===")
 waitGroupHappensBefore()

 fmt.Println("\n=== Once Happens-Before ===")
 onceHappensBefore()
}
```

#### 最佳实践

1. **理解Happens-Before规则**
2. **使用sync包建立同步**
3. **不要依赖内存可见性的假设**
4. **使用channel明确同步点**
5. **使用-race检测器验证**

---

### 8.3 内存同步

#### 概念定义

**内存同步**确保跨Goroutine的内存操作按预期顺序执行，避免编译器和CPU重排序导致的问题。

#### 完整示例

```go
package main

import (
 "fmt"
 "runtime"
 "sync"
 "sync/atomic"
)

// 内存屏障示例
func memoryBarrier() {
 var x, y int
 var wg sync.WaitGroup

 wg.Add(2)

 go func() {
  defer wg.Done()
  x = 1
  // 内存屏障（通过sync操作隐式插入）
  runtime.Gosched()
  _ = y
 }()

 go func() {
  defer wg.Done()
  y = 1
  // 内存屏障
  runtime.Gosched()
  _ = x
 }()

 wg.Wait()
}

// 原子操作提供顺序一致性
func atomicOrdering() {
 var counter int64
 var done int32

 go func() {
  for i := 0; i < 1000; i++ {
   atomic.AddInt64(&counter, 1)
  }
  atomic.StoreInt32(&done, 1)
 }()

 for atomic.LoadInt32(&done) == 0 {
  runtime.Gosched()
 }

 // 保证看到所有counter的更新
 fmt.Println(atomic.LoadInt64(&counter))
}

// 发布-订阅模式
func publishSubscribe() {
 var config *Config
 var ready int32

 // Publisher
 go func() {
  config = &Config{Version: 1, Data: "initial"}
  atomic.StoreInt32(&ready, 1) // 发布
 }()

 // Subscriber
 for atomic.LoadInt32(&ready) == 0 {
  runtime.Gosched()
 }

 // 保证看到config的完整初始化
 fmt.Printf("Config: %+v\n", config)
}

type Config struct {
 Version int
 Data    string
}

// 双检锁（需要正确实现）
type Singleton struct {
 value int
}

var instance *Singleton
var once sync.Once

func getSingleton() *Singleton {
 once.Do(func() {
  instance = &Singleton{value: 42}
 })
 return instance
}

// 错误的双检锁（不要这样做）
var badInstance *Singleton
var initialized uint32

func badGetSingleton() *Singleton {
 if atomic.LoadUint32(&initialized) == 1 {
  return badInstance
 }

 // 错误：instance可能未完全初始化
 badInstance = &Singleton{value: 42}
 atomic.StoreUint32(&initialized, 1)

 return badInstance
}

// 正确的双检锁
var correctInstance atomic.Value

func correctGetSingleton() *Singleton {
 if inst := correctInstance.Load(); inst != nil {
  return inst.(*Singleton)
 }

 inst := &Singleton{value: 42}
 correctInstance.Store(inst)
 return inst
}

// 缓存一致性
func cacheCoherence() {
 const numCPU = 4
 var counters [numCPU]int64

 var wg sync.WaitGroup
 for i := 0; i < numCPU; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   for j := 0; j < 1000000; j++ {
    counters[id]++ // 每个CPU修改自己的缓存行
   }
  }(i)
 }

 wg.Wait()

 var total int64
 for _, c := range counters {
  total += c
 }
 fmt.Printf("Total: %d\n", total)
}

// False Sharing示例
func falseSharing() {
 const numCPU = 4

 // 错误：变量可能在同一缓存行
 var counters [numCPU]int64

 // 正确：使用填充避免false sharing
 type PaddedCounter struct {
  value int64
  _     [56]byte // 填充到64字节（缓存行大小）
 }

 var paddedCounters [numCPU]PaddedCounter

 var wg sync.WaitGroup
 for i := 0; i < numCPU; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   for j := 0; j < 1000000; j++ {
    paddedCounters[id].value++
   }
  }(i)
 }

 wg.Wait()
}

func main() {
 fmt.Println("=== Atomic Ordering ===")
 atomicOrdering()

 fmt.Println("\n=== Publish Subscribe ===")
 publishSubscribe()

 fmt.Println("\n=== Singleton ===")
 s := getSingleton()
 fmt.Printf("Singleton: %+v\n", s)

 fmt.Println("\n=== Cache Coherence ===")
 cacheCoherence()
}
```

#### 最佳实践

1. **使用sync包而非手动同步**
2. **理解缓存行和false sharing**
3. **使用原子操作保证顺序**
4. **避免复杂的同步模式**
5. **使用atomic.Value进行无锁发布**

---

### 8.4 死锁检测与避免

#### 概念定义

**死锁**指两个或多个Goroutine互相等待对方释放资源，导致永久阻塞。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 死锁示例1: 互相等待
func deadlockExample1() {
 ch1 := make(chan int)
 ch2 := make(chan int)

 go func() {
  <-ch1 // 等待ch1
  ch2 <- 1
 }()

 <-ch2 // 等待ch2
 ch1 <- 1
}

// 死锁示例2: 循环等待锁
func deadlockExample2() {
 var mu1, mu2 sync.Mutex

 go func() {
  mu1.Lock()
  defer mu1.Unlock()

  time.Sleep(10 * time.Millisecond)

  mu2.Lock() // 等待mu2
  defer mu2.Unlock()
 }()

 mu2.Lock()
 defer mu2.Unlock()

 time.Sleep(10 * time.Millisecond)

 mu1.Lock() // 等待mu1 - 死锁！
 defer mu1.Unlock()
}

// 死锁示例3: channel容量不足
func deadlockExample3() {
 ch := make(chan int) // 无缓冲

 ch <- 1 // 阻塞，无接收者
 <-ch
}

// 死锁示例4: WaitGroup错误使用
func deadlockExample4() {
 var wg sync.WaitGroup

 wg.Add(1)
 go func() {
  // 忘记调用wg.Done()
 }()

 wg.Wait() // 永远等待
}

// 避免死锁：锁排序
func avoidDeadlockOrdering() {
 var mu1, mu2 sync.Mutex

 // 总是先锁mu1，再锁mu2
 lockBoth := func() {
  mu1.Lock()
  mu2.Lock()
  defer mu2.Unlock()
  defer mu1.Unlock()
  // 执行业务
 }

 go lockBoth()
 lockBoth()
}

// 避免死锁：超时机制
func avoidDeadlockTimeout() {
 ch := make(chan int)

 select {
 case ch <- 1:
  fmt.Println("Sent")
 case <-time.After(100 * time.Millisecond):
  fmt.Println("Timeout, avoiding deadlock")
 }
}

// 避免死锁：TryLock模式
func avoidDeadlockTryLock() {
 var mu sync.Mutex

 // 尝试获取锁
 tryLock := func() bool {
  select {
  case <-func() chan struct{} {
   ch := make(chan struct{}, 1)
   go func() {
    mu.Lock()
    ch <- struct{}{}
   }()
   return ch
  }():
   return true
  case <-time.After(100 * time.Millisecond):
   return false
  }
 }

 if tryLock() {
  defer mu.Unlock()
  fmt.Println("Got lock")
 } else {
  fmt.Println("Could not get lock, continuing...")
 }
}

// 死锁检测工具
func deadlockDetection() {
 // 1. 使用 go-deadlock 库
 // import "github.com/sasha-s/go-deadlock"
 // deadlock.Opts.DeadlockTimeout = time.Second
 // deadlock.Opts.OnPotentialDeadlock = func() {}

 // 2. 使用 go test -timeout
 // go test -timeout 30s ./...

 // 3. 使用 pprof 查看goroutine状态
 // go tool pprof http://localhost:6060/debug/pprof/goroutine

 fmt.Println("Use deadlock detection tools")
}

// 资源分配图检测
func resourceAllocationGraph() {
 // 理论方法：构建资源分配图，检测环
 // 实际中：使用超时和日志

 type Resource struct {
  ID    int
  Owner int // goroutine ID
 }

 resources := make(map[int]*Resource)
 var mu sync.Mutex

 acquire := func(gid, rid int) bool {
  mu.Lock()
  defer mu.Unlock()

  if r, ok := resources[rid]; ok && r.Owner != gid {
   // 资源被占用，可能死锁
   fmt.Printf("Goroutine %d waiting for resource %d (owned by %d)\n",
    gid, rid, r.Owner)
   return false
  }

  resources[rid] = &Resource{ID: rid, Owner: gid}
  return true
 }

 release := func(gid, rid int) {
  mu.Lock()
  defer mu.Unlock()
  delete(resources, rid)
 }

 _ = acquire
 _ = release
}

// 哲学家就餐问题（死锁经典案例）
func diningPhilosophers() {
 const n = 5
 var forks [n]sync.Mutex

 philosopher := func(id int) {
  left := id
  right := (id + 1) % n

  // 解决方案1：限制同时就餐人数
  // 解决方案2：奇偶哲学家不同顺序拿叉子
  // 解决方案3：同时拿两把叉子

  // 这里使用方案2
  if id%2 == 0 {
   forks[left].Lock()
   forks[right].Lock()
  } else {
   forks[right].Lock()
   forks[left].Lock()
  }

  fmt.Printf("Philosopher %d is eating\n", id)
  time.Sleep(10 * time.Millisecond)

  forks[left].Unlock()
  forks[right].Unlock()
 }

 for i := 0; i < n; i++ {
  go philosopher(i)
 }

 time.Sleep(100 * time.Millisecond)
}

func main() {
 fmt.Println("=== Avoid Deadlock Ordering ===")
 avoidDeadlockOrdering()

 fmt.Println("\n=== Avoid Deadlock Timeout ===")
 avoidDeadlockTimeout()

 fmt.Println("\n=== Dining Philosophers ===")
 diningPhilosophers()
}
```

#### 最佳实践

1. **统一锁获取顺序**
2. **使用超时机制**
3. **避免锁嵌套**
4. **使用channel替代锁**
5. **使用死锁检测工具**

---

### 8.5 活锁与饥饿

#### 概念定义

**活锁**指Goroutine不断改变状态但无法前进，**饥饿**指某些Goroutine长期无法获得资源。

#### 完整示例

```go
package main

import (
 "fmt"
 "sync"
 "sync/atomic"
 "time"
)

// 活锁示例
func livelockExample() {
 type Person struct {
  name      string
  direction string
 }

 alice := &Person{name: "Alice", direction: "left"}
 bob := &Person{name: "Bob", direction: "right"}

 // 两人在走廊相遇，互相让路
 var wg sync.WaitGroup
 wg.Add(2)

 move := func(p *Person, other *Person) {
  defer wg.Done()

  for i := 0; i < 10; i++ {
   // 如果对方挡路，就让开
   if p.direction == other.direction {
    fmt.Printf("%s steps aside\n", p.name)
    // 改变方向（但可能又挡住对方）
    if p.direction == "left" {
     p.direction = "right"
    } else {
     p.direction = "left"
    }
   } else {
    fmt.Printf("%s passes\n", p.name)
    return
   }
   time.Sleep(10 * time.Millisecond)
  }
 }

 go move(alice, bob)
 go move(bob, alice)

 wg.Wait()
}

// 解决活锁：随机等待
func solveLivelock() {
 var wg sync.WaitGroup
 wg.Add(2)

 move := func(name string) {
  defer wg.Done()

  for {
   // 随机等待，打破对称性
   time.Sleep(time.Duration(10+int(time.Now().UnixNano())%50) * time.Millisecond)

   fmt.Printf("%s moves\n", name)
   return
  }
 }

 go move("Alice")
 go move("Bob")

 wg.Wait()
}

// 饥饿示例
func starvationExample() {
 var mu sync.Mutex
 var wg sync.WaitGroup

 // 贪婪的reader
 for i := 0; i < 10; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   for j := 0; j < 100; j++ {
    mu.Lock()
    fmt.Printf("Reader %d: %d\n", id, j)
    mu.Unlock()
   }
  }(i)
 }

 // writer很难获得锁
 wg.Add(1)
 go func() {
  defer wg.Done()
  mu.Lock()
  fmt.Println("Writer acquired lock")
  mu.Unlock()
 }()

 wg.Wait()
}

// 使用RWMutex解决reader饥饿writer
func rwMutexSolution() {
 var mu sync.RWMutex
 var wg sync.WaitGroup

 // 使用读锁的reader
 for i := 0; i < 10; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   for j := 0; j < 100; j++ {
    mu.RLock()
    _ = j // 读取操作
    mu.RUnlock()
   }
  }(i)
 }

 // writer可以获得写锁
 wg.Add(1)
 go func() {
  defer wg.Done()
  mu.Lock()
  fmt.Println("Writer acquired lock")
  mu.Unlock()
 }()

 wg.Wait()
}

// 优先级反转
func priorityInversion() {
 // 低优先级任务持有锁，高优先级任务等待
 // 中等优先级任务抢占CPU，导致高优先级饿死

 var mu sync.Mutex

 // 低优先级任务
 go func() {
  mu.Lock()
  fmt.Println("Low priority: acquired lock")
  time.Sleep(500 * time.Millisecond) // 长时间持有
  mu.Unlock()
 }()

 time.Sleep(50 * time.Millisecond)

 // 高优先级任务
 go func() {
  fmt.Println("High priority: waiting for lock")
  mu.Lock()
  fmt.Println("High priority: acquired lock")
  mu.Unlock()
 }()

 // 中优先级任务（抢占CPU）
 go func() {
  for i := 0; i < 100000000; i++ {
   // 消耗CPU
  }
 }()

 time.Sleep(1 * time.Second)
}

// 公平锁实现
type FairMutex struct {
 ch chan struct{}
}

func NewFairMutex() *FairMutex {
 m := &FairMutex{ch: make(chan struct{}, 1)}
 m.ch <- struct{}{}
 return m
}

func (m *FairMutex) Lock() {
 <-m.ch
}

func (m *FairMutex) Unlock() {
 m.ch <- struct{}{}
}

// 使用公平锁
func fairLockExample() {
 mu := NewFairMutex()
 var wg sync.WaitGroup

 for i := 0; i < 5; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()
   mu.Lock()
   fmt.Printf("Goroutine %d acquired lock\n", id)
   time.Sleep(10 * time.Millisecond)
   mu.Unlock()
  }(i)
 }

 wg.Wait()
}

// 工作窃取避免饥饿
func workStealing() {
 type Queue struct {
  tasks []func()
  mu    sync.Mutex
 }

 queues := make([]*Queue, 4)
 for i := range queues {
  queues[i] = &Queue{}
 }

 // 分发任务
 for i := 0; i < 100; i++ {
  queueID := i % 4
  n := i
  queues[queueID].mu.Lock()
  queues[queueID].tasks = append(queues[queueID].tasks, func() {
   fmt.Printf("Task %d executed\n", n)
  })
  queues[queueID].mu.Unlock()
 }

 // worker处理自己的队列，空闲时窃取其他队列
 var wg sync.WaitGroup
 for i := 0; i < 4; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()

   myQueue := queues[id]

   for {
    myQueue.mu.Lock()
    if len(myQueue.tasks) > 0 {
     task := myQueue.tasks[0]
     myQueue.tasks = myQueue.tasks[1:]
     myQueue.mu.Unlock()
     task()
     continue
    }
    myQueue.mu.Unlock()

    // 尝试窃取
    stolen := false
    for j := 0; j < 4; j++ {
     if j == id {
      continue
     }

     otherQueue := queues[j]
     otherQueue.mu.Lock()
     if len(otherQueue.tasks) > 0 {
      task := otherQueue.tasks[len(otherQueue.tasks)-1]
      otherQueue.tasks = otherQueue.tasks[:len(otherQueue.tasks)-1]
      otherQueue.mu.Unlock()
      task()
      stolen = true
      break
     }
     otherQueue.mu.Unlock()
    }

    if !stolen {
     return
    }
   }
  }(i)
 }

 wg.Wait()
}

// 监控和检测饥饿
func starvationDetection() {
 var mu sync.Mutex
 waitTimes := make(map[int]time.Duration)

 monitor := func(id int, start time.Time) {
  mu.Lock()
  waitTimes[id] = time.Since(start)
  mu.Unlock()
 }

 var wg sync.WaitGroup
 for i := 0; i < 10; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()

   start := time.Now()
   mu.Lock()
   monitor(id, start)

   // 长时间持有
   time.Sleep(100 * time.Millisecond)
   mu.Unlock()
  }(i)
 }

 wg.Wait()

 // 分析等待时间
 for id, waitTime := range waitTimes {
  fmt.Printf("Goroutine %d waited %v\n", id, waitTime)
 }
}

func main() {
 fmt.Println("=== Solve Livelock ===")
 solveLivelock()

 fmt.Println("\n=== Fair Lock Example ===")
 fairLockExample()

 fmt.Println("\n=== Work Stealing ===")
 workStealing()

 fmt.Println("\n=== Starvation Detection ===")
 starvationDetection()
}
```

#### 最佳实践

1. **使用随机延迟打破对称性**
2. **使用公平锁避免饥饿**
3. **限制高优先级任务数量**
4. **使用工作窃取平衡负载**
5. **监控任务等待时间**

---

## 总结

本文档全面梳理了Go语言并发编程的核心概念和模式：

### 核心概念

1. **Goroutine**: 轻量级线程，由Go运行时管理
2. **Channel**: CSP模型的实现，用于Goroutine间通信
3. **Select**: 多路复用，实现非确定性选择
4. **Sync包**: 提供底层同步原语
5. **Context**: 传递取消信号和元数据

### 设计原则

1. **不要通过共享内存通信，通过通信共享内存**
2. **Channel用于协调，Mutex用于保护**
3. **优先使用channel，必要时使用sync包**
4. **使用context管理生命周期**
5. **避免过早优化，先保证正确性**

### 最佳实践

1. **使用-race标志检测数据竞争**
2. **理解Happens-Before关系**
3. **避免死锁、活锁和饥饿**
4. **合理设置GOMAXPROCS**
5. **使用pprof分析性能**

---

## 参考资源

- [Go Memory Model](https://golang.org/ref/mem)
- [Effective Go](https://golang.org/doc/effective_go.html#concurrency)
- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide)
- [Advanced Go Concurrency Patterns](https://talks.golang.org/2013/advconc.slide)

---

*文档完成日期: 2024年*
