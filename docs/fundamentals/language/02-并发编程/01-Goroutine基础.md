# Goroutine基础

> **简介**: Goroutine基础完整指南，包括创建、调度、GMP模型和最佳实践

> **版本**: Go 1.25.3  
> **难度**: ⭐⭐⭐  
> **标签**: #并发 #Goroutine #GMP #调度器

---

## 📚 目录

1. [Goroutine简介](#goroutine简介)
2. [创建与使用](#创建与使用)
3. [GMP调度模型](#gmp调度模型)
4. [最佳实践](#最佳实践)
5. [常见陷阱](#常见陷阱)

---

## 1. Goroutine简介

### 什么是Goroutine

**Goroutine** 是Go语言的轻量级线程：
- 由Go运行时管理，而非操作系统
- 启动成本低（约2KB栈空间）
- 可以创建成千上万个
- 自动调度到多个CPU核心

### Goroutine vs 线程

| 特性 | Goroutine | 线程 |
|------|-----------|------|
| **栈空间** | 2KB起（动态增长） | 1-2MB固定 |
| **创建成本** | 微秒级 | 毫秒级 |
| **切换成本** | 纳秒级 | 微秒级 |
| **调度器** | Go运行时 | 操作系统 |
| **数量** | 百万级 | 千级 |

---

## 2. 创建与使用

### 基本语法

```go
func main() {
    // 启动一个Goroutine
    go sayHello()
    
    // 启动匿名函数
    go func() {
        fmt.Println("Hello from anonymous goroutine")
    }()
    
    // 等待goroutine执行
    time.Sleep(1 * time.Second)
}

func sayHello() {
    fmt.Println("Hello from goroutine")
}
```

---

### 传递参数

```go
func main() {
    name := "World"
    
    // ❌ 错误：捕获循环变量
    for i := 0; i < 5; i++ {
        go func() {
            fmt.Println(i)  // 可能打印5个5
        }()
    }
    
    // ✅ 正确：传递参数
    for i := 0; i < 5; i++ {
        go func(n int) {
            fmt.Println(n)  // 正确打印0-4
        }(i)
    }
    
    time.Sleep(1 * time.Second)
}
```

---

### 等待Goroutine完成

#### 方法1: WaitGroup

```go
func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)  // 增加计数
        go func(n int) {
            defer wg.Done()  // 完成时减少计数
            process(n)
        }(i)
    }
    
    wg.Wait()  // 等待所有goroutine完成
    fmt.Println("All done")
}
```

#### 方法2: Channel

```go
func main() {
    done := make(chan bool)
    
    go func() {
        process()
        done <- true  // 发送完成信号
    }()
    
    <-done  // 等待完成
    fmt.Println("Done")
}
```

---

## 3. GMP调度模型

### GMP组件

- **G (Goroutine)**: Goroutine实例
- **M (Machine)**: 操作系统线程
- **P (Processor)**: 调度器的上下文（逻辑处理器）

```
┌─────────────────────────────────────┐
│         Global Queue                │
│     (全局Goroutine队列)              │
└─────────────────────────────────────┘
          │
          ↓
┌───────────┐   ┌───────────┐   ┌───────────┐
│   P0      │   │   P1      │   │   P2      │
│  ┌─────┐  │   │  ┌─────┐  │   │  ┌─────┐  │
│  │Local│  │   │  │Local│  │   │  │Local│  │
│  │Queue│  │   │  │Queue│  │   │  │Queue│  │
│  └─────┘  │   │  └─────┘  │   │  └─────┘  │
└───────────┘   └───────────┘   └───────────┘
      │               │               │
      ↓               ↓               ↓
┌───────────┐   ┌───────────┐   ┌───────────┐
│   M0      │   │   M1      │   │   M2      │
│ (Thread)  │   │ (Thread)  │   │ (Thread)  │
└───────────┘   └───────────┘   └───────────┘
```

---

### 调度策略

1. **工作窃取 (Work Stealing)**:
   - P的本地队列为空时，从其他P窃取一半的goroutine
   
2. **系统调用**:
   - Goroutine进行系统调用时，M会与P分离
   - P会寻找其他空闲M或创建新M

3. **抢占式调度**:
   - Go 1.14+: 基于信号的抢占
   - 防止单个goroutine长时间占用CPU

---

### GOMAXPROCS

```go
import "runtime"

func main() {
    // 获取当前P的数量
    fmt.Println(runtime.GOMAXPROCS(0))
    
    // 设置P的数量
    runtime.GOMAXPROCS(4)
    
    // Go 1.25+: 容器环境自动适配
    // 无需手动设置GOMAXPROCS
}
```

---

## 4. 最佳实践

### 1. 控制Goroutine数量

```go
// Worker Pool模式
func workerPool(jobs <-chan int, results chan<- int) {
    const numWorkers = 10
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    
    wg.Wait()
    close(results)
}
```

---

### 2. 使用Context管理生命周期

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker stopped")
            return
        default:
            doWork()
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go worker(ctx)
    
    time.Sleep(5 * time.Second)
    cancel()  // 停止worker
}
```

---

### 3. 避免Goroutine泄漏

```go
// ❌ 泄漏：Goroutine永远阻塞
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // 永远等待
        fmt.Println(val)
    }()
    // ch从未被写入
}

// ✅ 正确：使用带缓冲的channel或context
func noLeak(ctx context.Context) {
    ch := make(chan int, 1)  // 带缓冲
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return
        }
    }()
}
```

---

## 5. 常见陷阱

### 陷阱1: 循环变量捕获

```go
// ❌ 错误
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // 可能全部打印5
    }()
}

// ✅ 正确
for i := 0; i < 5; i++ {
    i := i  // Go 1.22+不需要这行
    go func() {
        fmt.Println(i)
    }()
}

// ✅ 正确：传递参数
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}
```

---

### 陷阱2: Goroutine数量爆炸

```go
// ❌ 错误：无限制创建goroutine
for _, item := range millionsOfItems {
    go process(item)  // 可能创建数百万个goroutine
}

// ✅ 正确：使用Worker Pool
jobs := make(chan Item, 100)
results := make(chan Result, 100)

// 启动固定数量的worker
for i := 0; i < 100; i++ {
    go worker(jobs, results)
}

// 发送任务
for _, item := range millionsOfItems {
    jobs <- item
}
close(jobs)
```

---

### 陷阱3: 未等待Goroutine完成

```go
// ❌ 错误：主函数可能在goroutine执行前退出
func main() {
    go doWork()
    // 程序立即退出
}

// ✅ 正确：使用WaitGroup等待
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        doWork()
    }()
    wg.Wait()
}
```

---

## 🔗 相关资源

- [Channel详解](./02-Channel详解.md)
- [Context应用](./03-Context应用.md)
- [并发模式](./05-并发模式.md)
- [性能优化](../../../advanced/performance/03-并发优化.md)

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3

