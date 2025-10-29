# Go并发编程

Go并发编程完整指南，涵盖Goroutine、Channel、Context和并发模式。

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
    fmt.Println("Hello from goroutine!")
}()
```

### Channel

```go
ch := make(chan int, 10)
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
jobs := make(chan int, 100)
results := make(chan int, 100)

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

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3
