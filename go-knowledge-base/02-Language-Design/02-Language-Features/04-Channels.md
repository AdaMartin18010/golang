# Channels

> **分类**: 语言设计

---

## 定义

Channel 是 goroutine 间的类型安全通信机制。

```go
ch := make(chan int)      // 无缓冲
ch := make(chan int, 10)  // 缓冲 10
```

---

## 操作

### 发送与接收

```go
ch <- v      // 发送
v := <-ch    // 接收
```

### 关闭

```go
close(ch)
```

---

## 缓冲 vs 无缓冲

### 无缓冲 (同步)

```go
ch := make(chan int)

// 发送阻塞直到有接收者
ch <- 42

// 接收阻塞直到有发送者
v := <-ch
```

**特性**: 同步通信，保证 happens-before

### 有缓冲 (异步)

```go
ch := make(chan int, 2)

ch <- 1  // 不阻塞
ch <- 2  // 不阻塞
ch <- 3  // 阻塞（缓冲满）
```

**特性**: 异步通信，解耦生产者消费者

---

## Select

```go
select {
case v1 := <-ch1:
    // 从 ch1 接收
    fmt.Println("ch1:", v1)
case v2 := <-ch2:
    // 从 ch2 接收
    fmt.Println("ch2:", v2)
case ch3 <- 100:
    // 发送到 ch3
    fmt.Println("sent to ch3")
default:
    // 默认分支
}
```

**特性**: 非确定性选择可用 case

---

## 方向

```go
// 双向
func process(ch chan int) { }

// 只接收
func receiveOnly(ch <-chan int) { }

// 只发送
func sendOnly(ch chan<- int) { }
```

---

## 模式

### 1. 工作池

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}
```

### 2. 扇出/扇入

```go
// 扇出: 多个 goroutine 读取同一 channel
for i := 0; i < 3; i++ {
    go worker(in)
}

// 扇入: 多个 channel 合并到一个
for i := 0; i < 3; i++ {
    go func() { out <- <-in }()
}
```

---

## 最佳实践

### 1. 由发送者关闭

```go
// 好
go func() {
    defer close(ch)
    for _, v := range values {
        ch <- v
    }
}()

// 不好: 由接收者关闭
```

### 2. 检查关闭

```go
v, ok := <-ch
if !ok {
    // channel 已关闭
}
```

### 3. for-range

```go
for v := range ch {
    // 自动处理关闭
}
```
