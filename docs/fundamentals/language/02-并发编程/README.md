# Go并发编程

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go并发编程](#go并发编程)
  - [📋 目录](#-目录)
  - [📚 文档列表](#-文档列表)
  - [🚀 快速示例](#-快速示例)
    - [Goroutine](#goroutine)
    - [Channel](#channel)
    - [Context](#context)
    - [Worker Pool](#worker-pool)
  - [📖 系统文档](#-系统文档)

---

## 📚 文档列表

1. **[Goroutine基础](./01-Goroutine基础.md)** ⭐⭐⭐⭐⭐
   - 创建与启动
   - GMP调度模型
   - 性能优化

2. **[Channel详解](./02-Channel详解.md)** ⭐⭐⭐⭐⭐
   - 无缓冲/有缓冲Channel
   - 关闭Channel
   - select多路复用

3. **[Context应用](./03-Context应用.md)** ⭐⭐⭐⭐⭐
   - 超时控制
   - 取消传播
   - 值传递

4. **[同步原语](./04-同步原语.md)** ⭐⭐⭐⭐⭐
   - Mutex, RWMutex
   - WaitGroup, Once
   - Cond, atomic

5. **[并发模式](./05-并发模式.md)** ⭐⭐⭐⭐⭐
   - Worker Pool
   - Pipeline
   - Fan-out/Fan-in
   - Context取消

6. **[内存模型](./06-内存模型.md)** ⭐⭐⭐⭐⭐
   - Happens-Before
   - 数据竞争检测

---

## 🚀 快速示例

### Goroutine

```go
go func() {
    fmt.Println("Hello from Goroutine!")
}()
```

### Channel

```go
ch := make(Channel int, 10)
go func() { ch <- 42 }()
value := <-ch
```

### Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case <-ctx.Done():
    fmt.Println("Timeout!")
case result := <-ch:
    fmt.Println(result)
}
```

### Worker Pool

```go
jobs := make(Channel int, 100)
results := make(Channel int, 100)

for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}

for j := 1; j <= 9; j++ {
    jobs <- j
}
close(jobs)
```

---

## 📖 系统文档
