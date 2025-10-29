# Go并发编程

Go并发编程完整指南，涵盖Goroutine、Channel、Context和并发模式。

---

## 📚 文档列表

1. **[并发基础概念](./01-并发基础概念.md)** ⭐⭐⭐⭐⭐
2. **[Goroutine深入](./02-Goroutine深入.md)** ⭐⭐⭐⭐⭐
3. **[Channel深入](./03-Channel深入.md)** ⭐⭐⭐⭐⭐
4. **[Context应用](./04-Context应用.md)** ⭐⭐⭐⭐⭐
5. **[并发模式](./05-并发模式.md)** ⭐⭐⭐⭐⭐

---

## 🎯 核心概念

### Goroutine
```go
go func() {
    fmt.Println("Hello from goroutine!")
}()
```

### Channel
```go
ch := make(chan int, 10)
ch <- 42
value := <-ch
```

### Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
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
