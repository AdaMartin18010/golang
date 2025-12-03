# Go并发编程

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go并发编程](#go并发编程)
  - [📋 目录](#-目录)
  - [📚 文档列表](#-文档列表)
  - [🎯 核心概念](#-核心概念)
    - [Goroutine](#goroutine)
    - [Channel](#channel)
    - [Context](#context)
  - [📖 系统文档](#-系统文档)

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
    fmt.Println("Hello from Goroutine!")
}()
```

### Channel

```go
ch := make(Channel int, 10)
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

### 📚 核心系统文档

这些文档提供结构化的知识体系：

- **[概念定义体系](./00-概念定义体系.md)** - 并发概念的形式化定义
- **[对比矩阵](./00-对比矩阵.md)** - 并发概念的系统对比
- **[知识图谱](./00-知识图谱.md)** - 并发知识的结构化总览

### 📖 学习建议

#### 初学者路径（3天）

```text
Day 1: 基础理论
  → 并发基础概念 (理解并发vs并行、CSP模型)
  → Goroutine基础 (创建goroutine)

Day 2: 通信机制
  → Channel深入 (无缓冲/缓冲channel)
  → 实践: 使用channel通信

Day 3: 控制与模式
  → Context应用 (取消、超时)
  → 并发模式 (Pipeline、Worker Pool)
```

#### 进阶学习（1周）

```text
Week 1: 完整体系
  → 深入学习所有5个主题
  → 理解底层机制（调度、栈管理）
  → 掌握常见并发模式
  → 实战项目实践
```

### 🎯 快速查找

**想理解概念**:
- 查看 [概念定义体系](./00-概念定义体系.md)
- 查看 [对比矩阵](./00-对比矩阵.md)

**规划学习路径**:
- 查看 [知识图谱](./00-知识图谱.md)

**学习具体主题**:
- 查看文档列表中的对应主题

### 🔗 相关资源

**基础主题**:
- [Go语言基础](../) - Go基础知识
- [函数式编程](../functions/) - 函数和闭包

**进阶主题**:
- [GMP调度器](../../advanced/runtime/scheduler.md) - 深入调度机制
- [性能优化](../../advanced/performance/) - 并发性能优化
- [测试](../testing/) - 并发代码测试

**实战资源**:
- [示例代码](../../../examples/concurrency/) - 并发示例
- [最佳实践](../../practices/concurrency-best-practices.md) - 最佳实践
- [常见陷阱](../../practices/concurrency-pitfalls.md) - 避免错误

---

**上次更新**: 2025-12-03
**维护者**: Go Framework Team
**反馈**: 欢迎通过Issue提供反馈
