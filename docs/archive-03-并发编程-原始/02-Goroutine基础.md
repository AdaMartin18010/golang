# 12.1.2 Goroutine基础

## 概述

Goroutine是Go语言并发编程的核心，它是一种轻量级的线程，由Go运行时管理。每个goroutine只需要几KB的栈空间，可以轻松创建成千上万个goroutine。

## Goroutine创建

### 基本语法

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello %s, iteration %d\n", name, i+1)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // 启动goroutine
    go sayHello("Alice")
    go sayHello("Bob")
    
    // 主goroutine等待
    time.Sleep(1 * time.Second)
    fmt.Println("Main goroutine finished")
}
```

### 匿名函数启动

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 使用匿名函数启动goroutine
    go func(name string) {
        fmt.Printf("Hello %s from anonymous function!\n", name)
    }("Charlie")
    
    time.Sleep(100 * time.Millisecond)
}
```

## Goroutine生命周期

### 1. 创建阶段

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    // 获取当前goroutine数量
    fmt.Printf("Goroutines before: %d\n", runtime.NumGoroutine())
    
    // 创建goroutine
    go func() {
        fmt.Println("New goroutine started")
    }()
    
    // 获取创建后的goroutine数量
    fmt.Printf("Goroutines after: %d\n", runtime.NumGoroutine())
}
```

### 2. 执行阶段

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        time.Sleep(100 * time.Millisecond)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // 启动3个worker goroutines
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // 发送任务
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)
    
    // 收集结果
    for a := 1; a <= 5; a++ {
        <-results
    }
}
```

### 3. 结束阶段

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func longRunningTask(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Task %s cancelled\n", name)
            return
        default:
            fmt.Printf("Task %s is running...\n", name)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // 创建可取消的context
    ctx, cancel := context.WithCancel(context.Background())
    
    // 启动长时间运行的任务
    go longRunningTask(ctx, "Task1")
    go longRunningTask(ctx, "Task2")
    
    // 3秒后取消所有任务
    time.Sleep(3 * time.Second)
    cancel()
    
    // 等待goroutine结束
    time.Sleep(100 * time.Millisecond)
    fmt.Println("All tasks finished")
}
```

## Goroutine调度

### 1. 协作式调度

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // 设置使用单核
    runtime.GOMAXPROCS(1)
    
    go func() {
        for i := 0; i < 5; i++ {
            fmt.Printf("Goroutine 1: %d\n", i)
            runtime.Gosched() // 主动让出CPU
        }
    }()
    
    go func() {
        for i := 0; i < 5; i++ {
            fmt.Printf("Goroutine 2: %d\n", i)
            runtime.Gosched() // 主动让出CPU
        }
    }()
    
    time.Sleep(1 * time.Second)
}
```

### 2. 抢占式调度

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func cpuIntensiveTask() {
    for i := 0; i < 1000000; i++ {
        // CPU密集型任务
        _ = i * i
    }
    fmt.Println("CPU intensive task completed")
}

func main() {
    // 启动多个CPU密集型goroutine
    for i := 0; i < 4; i++ {
        go func(id int) {
            fmt.Printf("Starting CPU task %d\n", id)
            cpuIntensiveTask()
            fmt.Printf("Finished CPU task %d\n", id)
        }(i)
    }
    
    time.Sleep(5 * time.Second)
}
```

## 常见模式

### 1. Fan-out/Fan-in模式

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func producer(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

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

func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    // 启动输出goroutine
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            out <- n
        }
    }
    
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }
    
    // 启动goroutine在完成后关闭out
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func main() {
    // 设置输入
    in := producer(2, 3, 4, 5)
    
    // 分发工作到两个goroutine
    c1 := square(in)
    c2 := square(in)
    
    // 合并结果
    for n := range merge(c1, c2) {
        fmt.Println(n)
    }
}
```

### 2. Pipeline模式

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // 创建pipeline
    c0 := make(chan string)
    c1 := make(chan string)
    c2 := make(chan string)
    
    // 第一阶段：生成数据
    go func() {
        defer close(c0)
        c0 <- "hello world"
        c0 <- "golang programming"
        c0 <- "concurrent processing"
    }()
    
    // 第二阶段：转换为大写
    go func() {
        defer close(c1)
        for s := range c0 {
            c1 <- strings.ToUpper(s)
        }
    }()
    
    // 第三阶段：添加前缀
    go func() {
        defer close(c2)
        for s := range c1 {
            c2 <- "PREFIX: " + s
        }
    }()
    
    // 输出结果
    for s := range c2 {
        fmt.Println(s)
    }
}
```

## 最佳实践

1. **避免goroutine泄漏**: 确保所有goroutine都能正常结束
2. **合理使用缓冲**: 根据场景选择合适的channel缓冲大小
3. **错误处理**: 在goroutine中正确处理错误
4. **资源管理**: 及时释放不再使用的资源

---

**下一节**: [03-Go调度器与GPM模型](03-Go调度器与GPM模型.md)
