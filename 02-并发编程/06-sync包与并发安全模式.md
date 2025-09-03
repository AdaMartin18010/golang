# 2.1 sync包与并发安全模式

<!-- TOC START -->
- [2.1 sync包与并发安全模式](#sync包与并发安全模式)
  - [2.1.1 1. 理论基础](#1-理论基础)
  - [2.1.2 2. 典型用法](#2-典型用法)
    - [2.1.2.1 互斥锁Mutex](#互斥锁mutex)
    - [2.1.2.2 读写锁RWMutex](#读写锁rwmutex)
    - [2.1.2.3 WaitGroup](#waitgroup)
    - [2.1.2.4 Once](#once)
  - [2.1.3 3. 工程分析与最佳实践](#3-工程分析与最佳实践)
  - [2.1.4 4. 常见陷阱](#4-常见陷阱)
  - [2.1.5 5. 单元测试建议](#5-单元测试建议)
  - [2.1.6 6. 参考文献](#6-参考文献)
<!-- TOC END -->














## 2.1.1 1. 理论基础

Go的sync包提供了多种并发原语，保障多Goroutine环境下的数据一致性和同步。

- **互斥锁（Mutex）**：保证同一时刻只有一个Goroutine访问临界区。
- **读写锁（RWMutex）**：读操作可并发，写操作独占。
- **等待组（WaitGroup）**：用于等待一组Goroutine完成。
- **Once**：确保某段代码只执行一次。
- **Cond**：条件变量，支持复杂同步。

---

## 2.1.2 2. 典型用法

### 2.1.2.1 互斥锁Mutex

```go
var mu sync.Mutex
mu.Lock()
// 临界区
mu.Unlock()
```

### 2.1.2.2 读写锁RWMutex

```go
var rw sync.RWMutex
rw.RLock()
// 只读区
rw.RUnlock()
rw.Lock()
// 写区
rw.Unlock()
```

### 2.1.2.3 WaitGroup

```go
var wg sync.WaitGroup
wg.Add(2)
go func() {
    defer wg.Done()
    // 任务1
}()
go func() {
    defer wg.Done()
    // 任务2
}()
wg.Wait()
```

### 2.1.2.4 Once

```go
var once sync.Once
once.Do(func() {
    // 只执行一次
})
```

---

## 2.1.3 3. 工程分析与最佳实践

- 推荐优先使用channel实现同步，sync适合低层并发控制。
- Mutex/RWMutex适合保护共享资源，避免数据竞争。
- WaitGroup适合任务编排，避免忙等。
- Once适合单例、懒加载等场景。
- Cond适合复杂同步，需谨慎使用。
- 尽量缩小锁的粒度，减少锁竞争。

---

## 2.1.4 4. 常见陷阱

- 忘记Unlock会导致死锁。
- 多次Unlock会panic。
- WaitGroup的Add/Done不匹配会导致永久阻塞。
- RWMutex写锁不可重入。

---

## 2.1.5 5. 单元测试建议

- 测试并发场景下的数据一致性与死锁边界。
- 使用-race检测数据竞争。

---

## 2.1.6 6. 参考文献

- Go官方文档：<https://golang.org/pkg/sync/>
- Go Blog: <https://blog.golang.org/share-memory-by-communicating>
- 《Go语言高级编程》
