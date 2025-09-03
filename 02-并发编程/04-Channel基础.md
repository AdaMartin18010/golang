# 2.1 Channel基础

<!-- TOC START -->
- [2.1 Channel基础](#channel基础)
  - [2.1.1 理论基础](#理论基础)
  - [2.1.2 典型用法](#典型用法)
    - [2.1.2.1 1. 无缓冲channel](#1-无缓冲channel)
    - [2.1.2.2 2. 缓冲channel](#2-缓冲channel)
    - [2.1.2.3 3. 单向channel](#3-单向channel)
    - [2.1.2.4 4. select多路复用](#4-select多路复用)
  - [2.1.3 工程分析与最佳实践](#工程分析与最佳实践)
  - [2.1.4 常见陷阱](#常见陷阱)
  - [2.1.5 单元测试建议](#单元测试建议)
  - [2.1.6 参考文献](#参考文献)
<!-- TOC END -->














## 2.1.1 理论基础

Go的Channel是CSP（Communicating Sequential Processes）模型的核心实现，支持Goroutine间安全通信。

- **CSP理论**：进程间通过消息传递（channel）而非共享内存通信。
- **Channel类型系统**：Go的channel是类型安全的，支持单向/双向、缓冲/无缓冲。

形式化描述：

- $\text{chan}\langle T \rangle$ 表示元素类型为T的channel。
- $P \xrightarrow{c!v} Q$ 表示P通过channel c发送v给Q。

---

## 2.1.2 典型用法

### 2.1.2.1 1. 无缓冲channel

```go
ch := make(chan int)
go func() { ch <- 42 }()
val := <-ch
fmt.Println(val)
```

### 2.1.2.2 2. 缓冲channel

```go
ch := make(chan string, 2)
ch <- "hello"
ch <- "world"
fmt.Println(<-ch)
fmt.Println(<-ch)
```

### 2.1.2.3 3. 单向channel

```go
func send(ch chan<- int) { ch <- 1 }
func recv(ch <-chan int) { fmt.Println(<-ch) }
```

### 2.1.2.4 4. select多路复用

```go
select {
case v := <-ch1:
    fmt.Println("ch1:", v)
case v := <-ch2:
    fmt.Println("ch2:", v)
default:
    fmt.Println("no data")
}
```

---

## 2.1.3 工程分析与最佳实践

- 推荐优先使用无缓冲channel，保证同步性。
- select语句可实现超时、取消等高级控制。
- 单向channel提升类型安全，利于接口设计。
- 注意channel关闭（close）时机，避免panic。
- 避免死锁：所有发送/接收都必须有对应方。
- 使用range遍历channel时，需配合close。

---

## 2.1.4 常见陷阱

- 向已关闭的channel发送数据会panic。
- 从已关闭的channel接收数据会返回零值。
- 忽略channel缓冲区满/空的情况易导致死锁。

---

## 2.1.5 单元测试建议

- 测试并发场景下的channel通信正确性。
- 覆盖关闭、超时、异常等边界情况。

---

## 2.1.6 参考文献

- Go官方文档：<https://golang.org/doc/>
- Go Blog: <https://blog.golang.org/pipelines>
- 《Go语言高级编程》
