# Channel详解

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

---

## 📋 目录

- [1. Channel基础](#1.-channel基础)
  - [什么是Channel](#什么是channel)
  - [声明和创建](#声明和创建)
  - [基本操作](#基本操作)
- [2. 无缓冲Channel](#2.-无缓冲channel)
  - [特点](#特点)
  - [示例](#同步示例)
  - [同步示例](#同步示例)
- [3. 有缓冲Channel](#3.-有缓冲channel)
  - [特点](#特点)
  - [示例](#同步示例)
  - [查询状态](#查询状态)
- [4. 关闭Channel](#4.-关闭channel)
  - [close操作](#close操作)
  - [检查是否关闭](#检查是否关闭)
  - [关闭规则](#关闭规则)
- [5. select多路复用](#5.-select多路复用)
  - [基本语法](#基本语法)
  - [超时控制](#超时控制)
  - [非阻塞接收/发送](#非阻塞接收发送)
  - [多channel监听](#多channel监听)
- [6. Channel模式](#6.-channel模式)
  - [模式1: Generator](#模式1-generator)
  - [模式2: Pipeline](#模式2-pipeline)
  - [模式3: Fan-out/Fan-in](#模式3-fan-outfan-in)
  - [模式4: 退出通知](#模式4-退出通知)
- [7. 最佳实践](#7.-最佳实践)
  - [1. 发送者关闭Channel](#1.-发送者关闭channel)
  - [2. 使用有缓冲Channel避免阻塞](#2.-使用有缓冲channel避免阻塞)
  - [3. 使用select处理超时](#3.-使用select处理超时)
  - [4. nil Channel的使用](#4.-nil-channel的使用)
  - [5. 避免Channel泄漏](#5.-避免channel泄漏)
- [🔗 相关资源](#相关资源)

## 1. Channel基础

### 什么是Channel

**Channel** 是Go中的通信机制：

- 用于Goroutine之间传递数据
- 类型安全的管道
- 支持同步和异步通信
- "不要通过共享内存来通信，而要通过通信来共享内存"

### 声明和创建

```go
// 声明
var ch chan int          // nil channel
var ch chan string       // nil channel

// 创建无缓冲channel
ch := make(chan int)

// 创建有缓冲channel
ch := make(chan int, 10)

// 只读channel
var readOnly <-chan int = ch

// 只写channel
var writeOnly chan<- int = ch
```

---

### 基本操作

```go
ch := make(chan int)

// 发送（写入）
ch <- 42

// 接收（读取）
value := <-ch

// 接收并检查是否关闭
value, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}

// 关闭
close(ch)
```

---

## 2. 无缓冲Channel

### 特点

- **同步通信**：发送操作会阻塞，直到有接收者
- **握手机制**：发送者和接收者必须同时准备好
- **容量为0**：make(chan T)

### 示例

```go
func main() {
    ch := make(chan int)  // 无缓冲
    
    // ❌ 错误：死锁
    // ch <- 42  // 阻塞，没有接收者
    
    // ✅ 正确：在goroutine中发送
    go func() {
        ch <- 42
    }()
    
    value := <-ch  // 接收
    fmt.Println(value)  // 42
}
```

---

### 同步示例

```go
func worker(done chan bool) {
    fmt.Println("Working...")
    time.Sleep(1 * time.Second)
    fmt.Println("Done")
    done <- true  // 发送完成信号
}

func main() {
    done := make(chan bool)
    go worker(done)
    <-done  // 等待完成
    fmt.Println("All done")
}
```

---

## 3. 有缓冲Channel

### 特点

- **异步通信**：发送操作不会立即阻塞（缓冲区未满时）
- **有容量**：make(chan T, capacity)
- **FIFO队列**：先进先出

### 示例

```go
func main() {
    ch := make(chan int, 3)  // 缓冲大小为3
    
    // 可以连续发送3个值而不阻塞
    ch <- 1
    ch <- 2
    ch <- 3
    
    // ch <- 4  // 会阻塞，缓冲区满了
    
    fmt.Println(<-ch)  // 1
    fmt.Println(<-ch)  // 2
    fmt.Println(<-ch)  // 3
}
```

---

### 查询状态

```go
ch := make(chan int, 10)

// 当前长度
fmt.Println(len(ch))  // 0

// 容量
fmt.Println(cap(ch))  // 10

// 发送3个值
ch <- 1
ch <- 2
ch <- 3

fmt.Println(len(ch))  // 3
fmt.Println(cap(ch))  // 10
```

---

## 4. 关闭Channel

### close操作

```go
ch := make(chan int, 3)

// 发送数据
ch <- 1
ch <- 2
ch <- 3

// 关闭channel
close(ch)

// ❌ 错误：向已关闭的channel发送数据会panic
// ch <- 4  // panic: send on closed channel

// ✅ 正确：从已关闭的channel接收数据
fmt.Println(<-ch)  // 1
fmt.Println(<-ch)  // 2
fmt.Println(<-ch)  // 3
fmt.Println(<-ch)  // 0 (零值)
```

---

### 检查是否关闭

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)

// 方法1: 使用ok标识
for {
    value, ok := <-ch
    if !ok {
        break  // channel已关闭
    }
    fmt.Println(value)
}

// 方法2: 使用range（推荐）
ch2 := make(chan int, 2)
ch2 <- 1
ch2 <- 2
close(ch2)

for value := range ch2 {
    fmt.Println(value)  // 1, 2
}
// range会自动处理关闭的channel
```

---

### 关闭规则

```go
// ✅ 可以：发送者关闭channel
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)  // 发送者关闭
}

// ❌ 不要：接收者关闭channel
func consumer(ch <-chan int) {
    for v := range ch {
        fmt.Println(v)
    }
    // 不要在这里关闭
}

// ❌ 不要：多次关闭channel
close(ch)
// close(ch)  // panic: close of closed channel
```

---

## 5. select多路复用

### 基本语法

```go
select {
case value := <-ch1:
    // 从ch1接收
case ch2 <- value:
    // 向ch2发送
case <-ch3:
    // 从ch3接收（丢弃值）
default:
    // 所有case都阻塞时执行
}
```

---

### 超时控制

```go
func main() {
    ch := make(chan int)
    
    select {
    case value := <-ch:
        fmt.Println(value)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

---

### 非阻塞接收/发送

```go
// 非阻塞接收
select {
case value := <-ch:
    fmt.Println("Received:", value)
default:
    fmt.Println("No data available")
}

// 非阻塞发送
select {
case ch <- 42:
    fmt.Println("Sent")
default:
    fmt.Println("Channel full")
}
```

---

### 多channel监听

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "one"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "two"
    }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received:", msg2)
        }
    }
}
```

---

## 6. Channel模式

### 模式1: Generator

```go
func fibonacci(n int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        a, b := 0, 1
        for i := 0; i < n; i++ {
            ch <- a
            a, b = b, a+b
        }
    }()
    return ch
}

// 使用
for num := range fibonacci(10) {
    fmt.Println(num)
}
```

---

### 模式2: Pipeline

```go
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

// 使用
nums := generator(1, 2, 3, 4, 5)
squares := square(nums)
for sq := range squares {
    fmt.Println(sq)  // 1, 4, 9, 16, 25
}
```

---

### 模式3: Fan-out/Fan-in

```go
func fanOut(ch <-chan int, n int) []<-chan int {
    channels := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        channels[i] = worker(ch)
    }
    return channels
}

func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

---

### 模式4: 退出通知

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
    cancel()  // 通知退出
}
```

---

## 7. 最佳实践

### 1. 发送者关闭Channel

```go
// ✅ 正确
func producer(ch chan<- int) {
    defer close(ch)  // 发送者关闭
    for i := 0; i < 10; i++ {
        ch <- i
    }
}
```

---

### 2. 使用有缓冲Channel避免阻塞

```go
// ❌ 可能阻塞
ch := make(chan int)
ch <- 42  // 阻塞

// ✅ 使用缓冲
ch := make(chan int, 1)
ch <- 42  // 不阻塞
```

---

### 3. 使用select处理超时

```go
// ✅ 推荐
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case result := <-ch:
    handleResult(result)
case <-ctx.Done():
    handleTimeout()
}
```

---

### 4. nil Channel的使用

```go
var ch chan int  // nil channel

// 从nil channel接收：永远阻塞
// <-ch

// 向nil channel发送：永远阻塞
// ch <- 1

// 关闭nil channel：panic
// close(ch)

// 用途：在select中禁用某个case
ch1 := make(chan int)
ch2 := make(chan int)

select {
case v := <-ch1:
    // 处理ch1
    ch1 = nil  // 禁用这个case
case v := <-ch2:
    // 处理ch2
}
```

---

### 5. 避免Channel泄漏

```go
// ❌ 泄漏：goroutine永远阻塞
func leak() <-chan int {
    ch := make(chan int)
    go func() {
        ch <- 42  // 没有接收者，永远阻塞
    }()
    return ch
}

// ✅ 正确：使用缓冲或context
func noLeak(ctx context.Context) <-chan int {
    ch := make(chan int, 1)  // 有缓冲
    go func() {
        select {
        case ch <- 42:
        case <-ctx.Done():
            return
        }
    }()
    return ch
}
```

---

## 🔗 相关资源

- [Goroutine基础](./01-Goroutine基础.md)
- [Context应用](./03-Context应用.md)
- [并发模式](./05-并发模式.md)

---

**最后更新**: 2025-10-29  
**Go版本**: 1.25.3
