# 并发编程基础

**章节定位**: Go语言基础 > 并发编程深入  
**难度级别**: 中级  
**预计学习时间**: 8-10小时

---

## 📖 章节概述

本章节深入讲解Go语言并发编程的核心概念和实战技巧。Go的并发模型基于CSP（Communicating Sequential Processes），通过Goroutine和Channel实现轻量级并发，是Go语言最具特色和威力的特性之一。

### 🎯 学习目标

完成本章学习后，你将能够：

- ✅ 深入理解Go并发模型和CSP原理
- ✅ 掌握Goroutine的创建、调度和生命周期管理
- ✅ 熟练使用各种类型的Channel进行通信
- ✅ 使用Context实现超时控制和取消传播
- ✅ 实现常见的并发模式和最佳实践
- ✅ 避免常见的并发陷阱和问题

---

## 📚 章节内容

### [01-并发基础概念](./01-并发基础概念.md)
**难度**: ⭐⭐  
**预计阅读**: 15分钟

- 并发 vs 并行的区别
- CSP模型详解
- Go并发模型的优势
- 并发原语概览

**关键概念**: CSP、Goroutine、Channel、并发安全

---

### [02-Goroutine深入](./02-Goroutine深入.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 20分钟

- Goroutine的创建和启动
- G-P-M调度模型
- Goroutine生命周期管理
- 栈管理和内存开销
- 避免Goroutine泄漏

**关键概念**: 协程、调度器、栈增长、泄漏检测

---

### [03-Channel深入](./03-Channel深入.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 25分钟

- Channel的类型和特性
- 缓冲Channel vs 无缓冲Channel
- Channel的关闭和检测
- select多路复用
- Channel的内部实现

**关键概念**: 通道、阻塞、非阻塞、select、range

---

### [04-Context应用](./04-Context应用.md)
**难度**: ⭐⭐⭐  
**预计阅读**: 20分钟

- Context的设计理念
- 四种Context类型
- 超时控制实战
- 取消信号传播
- 值传递最佳实践
- Context在HTTP中的应用

**关键概念**: 上下文、取消、超时、值传递

---

### [05-并发模式](./05-并发模式.md)
**难度**: ⭐⭐⭐⭐  
**预计阅读**: 30分钟

- Pipeline模式
- Fan-out/Fan-in模式
- Worker Pool模式
- Or-Channel模式
- Tee-Channel模式
- Bridge-Channel模式

**关键概念**: 并发模式、通道组合、任务编排

---

## 🎓 学习路径

### 初学者路线
```
01并发基础概念 → 02Goroutine深入 → 03Channel深入
```

### 进阶路线
```
04Context应用 → 05并发模式 → 实战项目
```

### 推荐学习顺序

1. **第一步**: 理解并发基础概念（1小时）
   - 阅读 01-并发基础概念.md
   - 运行基础示例代码

2. **第二步**: 掌握Goroutine（2小时）
   - 阅读 02-Goroutine深入.md
   - 实践创建和管理Goroutine
   - 学习避免泄漏

3. **第三步**: 精通Channel（2-3小时）
   - 阅读 03-Channel深入.md
   - 实践各种Channel用法
   - 掌握select多路复用

4. **第四步**: 应用Context（1-2小时）
   - 阅读 04-Context应用.md
   - 实现超时控制
   - 练习取消传播

5. **第五步**: 并发模式实战（2-3小时）
   - 阅读 05-并发模式.md
   - 实现常见并发模式
   - 完成实战项目

---

## 🔗 相关章节

### 前置知识
- [语法基础](../language/01-语法基础/) - Go基础语法
- [函数](../language/01-语法基础/05-函数.md) - 函数定义和调用
- [错误处理](../language/01-语法基础/11-错误处理.md) - 错误处理机制

### 后续进阶
- [Go调度器](../language/02-并发编程/04-Go调度器.md) - 深入调度器原理
- [sync包](../language/02-并发编程/06-sync包.md) - 同步原语
- [并发优化](../../advanced/performance/03-并发优化.md) - 并发性能优化
- [并发型模式](../../advanced/architecture/04-并发型模式.md) - 高级并发模式

---

## 💻 实战项目

完成学习后，建议实践以下项目：

1. **并发爬虫** - 使用Worker Pool模式
2. **实时数据处理** - Pipeline模式
3. **微服务超时控制** - Context应用
4. **高并发Web服务器** - 综合应用

---

## ⚠️ 常见问题

### Q1: Goroutine和线程有什么区别？
- Goroutine更轻量（2KB初始栈 vs 1-2MB线程栈）
- 由Go运行时调度，而非操作系统
- M:N调度模型，多个Goroutine复用少量线程

### Q2: 什么时候使用缓冲Channel？
- 解耦生产者和消费者速度
- 实现异步通信
- 减少阻塞，提高吞吐量

### Q3: 如何避免Goroutine泄漏？
- 确保所有Goroutine都有退出条件
- 使用Context实现取消机制
- 使用defer和recover处理panic
- 定期监控Goroutine数量

### Q4: select应该怎么用？
- 多路复用多个Channel
- 实现超时控制（time.After）
- 实现非阻塞操作（default分支）
- 实现取消操作（Context.Done）

---

## 📚 推荐资源

### 官方文档
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Blog - Concurrency Patterns](https://go.dev/blog/pipelines)

### 推荐书籍
- 《Concurrency in Go》 by Katherine Cox-Buday
- 《Go并发编程实战》

### 在线资源
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)

---

## 🎯 下一步

完成本章学习后，你可以：

1. **深入调度器**: 学习 [Go调度器](../language/02-并发编程/04-Go调度器.md)
2. **同步原语**: 学习 [sync包](../language/02-并发编程/06-sync包.md)
3. **性能优化**: 学习 [并发优化](../../advanced/performance/03-并发优化.md)
4. **设计模式**: 学习 [并发型模式](../../advanced/architecture/04-并发型模式.md)

---

**维护者**: Documentation Team  
**最后更新**: 2025-10-27  
**版本**: v1.0

