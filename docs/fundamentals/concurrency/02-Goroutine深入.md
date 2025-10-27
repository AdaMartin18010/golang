# Goroutine深入

**难度**: 中级  
**预计阅读时间**: 20分钟  
**前置知识**: Go基础语法、并发基础概念

---

## 📖 概念介绍

Goroutine是Go语言实现并发的核心机制，它是一种轻量级的用户态线程。理解Goroutine的创建、调度和生命周期管理，是掌握Go并发编程的关键。

---

## 🎯 核心知识点

### 1. Goroutine的创建和启动

#### 基础创建

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // 普通函数调用（同步）
    sayHello("Alice")
    
    // 使用go关键字创建Goroutine（异步）
    go sayHello("Bob")
    go sayHello("Charlie")
    
    // 主Goroutine需要等待，否则程序会立即退出
    time.Sleep(1 * time.Second)
    fmt.Println("Main goroutine exiting")
}
```

#### 匿名函数Goroutine

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    
    // 方式1：匿名函数
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("匿名函数Goroutine")
    }()
    
    // 方式2：带参数的匿名函数
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) { // 通过参数传递，避免闭包陷阱
            defer wg.Done()
            fmt.Printf("Goroutine %d\n", id)
        }(i)
    }
    
    wg.Wait()
}
```

#### 闭包陷阱

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func closureTrap() {
    var wg sync.WaitGroup
    
    // ❌ 错误：所有Goroutine都会打印5
    fmt.Println("错误示例：")
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fmt.Printf("%d ", i) // i被所有goroutine共享
        }()
    }
    wg.Wait()
    fmt.Println()
    
    // ✅ 正确：通过参数传递
    fmt.Println("正确示例：")
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("%d ", id)
        }(i) // 将i作为参数传递
    }
    wg.Wait()
    fmt.Println()
}

func main() {
    closureTrap()
}
```

---

### 2. G-P-M调度模型

#### 模型组成

```
G (Goroutine)  - 用户态线程
P (Processor)  - 逻辑处理器（通常等于CPU核心数）
M (Machine)    - 系统线程（OS Thread）

关系：
- 多个G运行在P上
- P绑定到M上执行
- M由操作系统调度
```

#### 查看调度信息

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func schedulerInfo() {
    // 获取GOMAXPROCS（P的数量）
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
    
    // 获取CPU核心数
    fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
    
    // 获取当前Goroutine数量
    fmt.Printf("Initial Goroutines: %d\n", runtime.NumGoroutine())
    
    // 创建1000个Goroutine
    for i := 0; i < 1000; i++ {
        go func() {
            time.Sleep(1 * time.Second)
        }()
    }
    
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("After creating 1000 goroutines: %d\n", runtime.NumGoroutine())
    
    // 设置GOMAXPROCS
    runtime.GOMAXPROCS(2)
    fmt.Printf("Set GOMAXPROCS to: %d\n", runtime.GOMAXPROCS(0))
}

func main() {
    schedulerInfo()
    time.Sleep(2 * time.Second)
}
```

#### 调度策略

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

// Work Stealing演示
func workStealingDemo() {
    runtime.GOMAXPROCS(4) // 使用4个P
    
    var wg sync.WaitGroup
    
    // 创建不平衡的工作负载
    for i := 0; i < 8; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // 不同的工作量
            iterations := 1000000 * (id + 1)
            sum := 0
            for j := 0; j < iterations; j++ {
                sum += j
            }
            
            fmt.Printf("Goroutine %d finished with sum=%d\n", id, sum)
        }(i)
    }
    
    wg.Wait()
}

func main() {
    workStealingDemo()
}
```

---

### 3. Goroutine生命周期管理

#### 使用WaitGroup

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // 确保Done被调用
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Duration(id) * 100 * time.Millisecond)
    fmt.Printf("Worker %d done\n", id)
}

func waitGroupDemo() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    wg.Wait() // 等待所有worker完成
    fmt.Println("All workers completed")
}

func main() {
    waitGroupDemo()
}
```

#### 使用Channel同步

```go
package main

import (
    "fmt"
    "time"
)

func workerWithChannel(id int, done chan bool) {
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Duration(id) * 100 * time.Millisecond)
    fmt.Printf("Worker %d done\n", id)
    done <- true // 发送完成信号
}

func channelSyncDemo() {
    done := make(chan bool, 5) // 缓冲channel
    
    for i := 1; i <= 5; i++ {
        go workerWithChannel(i, done)
    }
    
    // 等待所有worker完成
    for i := 1; i <= 5; i++ {
        <-done
    }
    
    fmt.Println("All workers completed")
}

func main() {
    channelSyncDemo()
}
```

#### 使用Context控制

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func workerWithContext(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d cancelled\n", id)
            return
        default:
            fmt.Printf("Worker %d working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func contextDemo() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    for i := 1; i <= 3; i++ {
        go workerWithContext(ctx, i)
    }
    
    <-ctx.Done()
    fmt.Println("All workers should have stopped")
    time.Sleep(100 * time.Millisecond) // 等待打印完成
}

func main() {
    contextDemo()
}
```

---

### 4. 栈管理和内存开销

#### 栈大小演示

```go
package main

import (
    "fmt"
    "runtime"
)

func stackGrowth() {
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    
    // 创建10000个Goroutine
    done := make(chan bool)
    for i := 0; i < 10000; i++ {
        go func() {
            // 递归调用，触发栈增长
            var arr [1024]byte
            _ = arr
            <-done
        }()
    }
    
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    fmt.Printf("初始分配: %d KB\n", m1.Alloc/1024)
    fmt.Printf("创建10000个Goroutine后: %d KB\n", m2.Alloc/1024)
    fmt.Printf("平均每个Goroutine: %d KB\n", (m2.Alloc-m1.Alloc)/1024/10000)
    
    close(done)
}

func main() {
    stackGrowth()
}
```

#### 栈增长机制

```go
package main

import (
    "fmt"
    "runtime"
)

func deepRecursion(n int) int {
    if n <= 0 {
        return 0
    }
    // 每次递归都会占用栈空间
    var arr [100]int // 800字节
    arr[0] = n
    return arr[0] + deepRecursion(n-1)
}

func stackGrowthDemo() {
    // 打印初始栈大小
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    fmt.Printf("初始栈大小: %d KB\n", m1.StackInuse/1024)
    
    // 深度递归触发栈增长
    result := deepRecursion(1000)
    
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    fmt.Printf("递归后栈大小: %d KB\n", m2.StackInuse/1024)
    fmt.Printf("递归结果: %d\n", result)
}

func main() {
    stackGrowthDemo()
}
```

---

### 5. 避免Goroutine泄漏

#### 泄漏案例1：阻塞的Channel

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

// ❌ 错误：Goroutine泄漏
func leakyGoroutine() {
    ch := make(chan int)
    
    go func() {
        val := <-ch // 永远阻塞，Goroutine泄漏
        fmt.Println(val)
    }()
    
    // ch没有发送数据，Goroutine永远阻塞
}

// ✅ 正确：使用Context取消
func fixedGoroutine() {
    ch := make(chan int)
    done := make(chan bool)
    
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-done:
            fmt.Println("Goroutine cancelled")
            return
        }
    }()
    
    time.Sleep(1 * time.Second)
    close(done) // 取消Goroutine
}

func detectLeak() {
    fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())
    
    for i := 0; i < 10; i++ {
        leakyGoroutine()
    }
    
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("After leaky: %d goroutines (leaked 10)\n", runtime.NumGoroutine())
    
    for i := 0; i < 10; i++ {
        fixedGoroutine()
    }
    
    time.Sleep(2 * time.Second)
    fmt.Printf("After fixed: %d goroutines (no leak)\n", runtime.NumGoroutine())
}

func main() {
    detectLeak()
}
```

#### 泄漏案例2：无限循环

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// ❌ 错误：无退出条件
func leakyLoop() {
    go func() {
        for {
            // 无限循环，无退出条件
            time.Sleep(100 * time.Millisecond)
            fmt.Println("Working...")
        }
    }()
}

// ✅ 正确：使用Context控制退出
func fixedLoop() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Loop cancelled")
                return
            default:
                time.Sleep(100 * time.Millisecond)
                fmt.Println("Working...")
            }
        }
    }()
    
    time.Sleep(1 * time.Second)
    cancel() // 取消循环
}

func main() {
    fixedLoop()
    time.Sleep(200 * time.Millisecond)
}
```

#### 泄漏检测工具

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func monitorGoroutines() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for i := 0; i < 5; i++ {
        <-ticker.C
        fmt.Printf("[%d] Goroutines: %d\n", i, runtime.NumGoroutine())
    }
}

func main() {
    go monitorGoroutines()
    
    // 创建一些Goroutine
    for i := 0; i < 10; i++ {
        go func() {
            time.Sleep(3 * time.Second)
        }()
    }
    
    time.Sleep(6 * time.Second)
}
```

---

## ⚠️ 常见问题

### Q1: Goroutine的开销有多大？
- 初始栈大小：2KB
- 创建时间：约纳秒级
- 上下文切换：用户态，比线程快得多
- 可以轻松创建数万个Goroutine

### Q2: 如何限制Goroutine数量？
```go
// 使用缓冲Channel作为信号量
sem := make(chan struct{}, 100) // 最多100个并发

for i := 0; i < 1000; i++ {
    sem <- struct{}{} // 获取
    go func() {
        defer func() { <-sem }() // 释放
        // 工作...
    }()
}
```

### Q3: 如何知道Goroutine何时结束？
- 使用WaitGroup
- 使用Channel通信
- 使用Context传播取消信号

### Q4: Goroutine会被GC回收吗？
- Goroutine本身不会被GC
- 但Goroutine引用的对象会被GC
- 必须显式确保Goroutine退出

---

## 📚 相关资源

### 下一步学习
- [03-Channel深入](./03-Channel深入.md)
- [04-Context应用](./04-Context应用.md)
- [Go调度器](../language/02-并发编程/04-Go调度器.md)

### 推荐阅读
- [Go Scheduler Design](https://golang.org/s/go11sched)
- [Goroutine泄漏检测](https://github.com/uber-go/goleak)

---

**最后更新**: 2025-10-27  
**作者**: Documentation Team

