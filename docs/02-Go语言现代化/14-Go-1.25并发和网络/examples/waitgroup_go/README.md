# WaitGroup.Go() 示例

> **Go 版本**: 1.25+  
> **目的**: 演示 `sync.WaitGroup` 的新 `Go()` 方法

---

## 快速开始

### 编译和运行

```bash
# 编译
go build -o waitgroup_demo basic_example.go

# 运行
./waitgroup_demo
```

---

## 示例说明

### 示例 1: 基本使用

演示 `WaitGroup.Go()` 的基本用法:

```go
var wg sync.WaitGroup

wg.Go(func() {
    // 任务 1
})

wg.Go(func() {
    // 任务 2
})

wg.Wait()
```

---

### 示例 2: 并行处理切片

演示如何并行处理切片元素:

```go
items := []string{"a", "b", "c"}

for _, item := range items {
    wg.Go(func() {
        process(item)  // Go 1.22+ 自动捕获 item 副本
    })
}

wg.Wait()
```

---

### 示例 3: 限制并发数

使用信号量限制并发 goroutine 数量:

```go
maxConcurrency := 3
sem := make(chan struct{}, maxConcurrency)

for _, item := range items {
    sem <- struct{}{}  // 获取信号量
    
    wg.Go(func() {
        defer func() { <-sem }()  // 释放信号量
        process(item)
    })
}

wg.Wait()
```

---

### 示例 4: 收集结果

通过 channel 收集 goroutine 的结果:

```go
results := make(chan int, len(items))

for _, item := range items {
    wg.Go(func() {
        result := compute(item)
        results <- result
    })
}

wg.Wait()
close(results)

// 收集结果
for result := range results {
    fmt.Println(result)
}
```

---

### 示例 5: 错误处理

使用互斥锁收集错误:

```go
var mu sync.Mutex
var errors []error

for _, item := range items {
    wg.Go(func() {
        if err := process(item); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        }
    })
}

wg.Wait()
```

---

### 示例 6: 传统方式对比

对比传统 WaitGroup 和 `WaitGroup.Go()`:

```go
// 传统方式: 4 行
wg.Add(1)
go func() {
    defer wg.Done()
    work()
}()

// WaitGroup.Go(): 1 行
wg.Go(func() {
    work()
})
```

---

## 运行输出

```text
Go 1.25 WaitGroup.Go() 示例

=== 示例 1: 基本使用 ===
Task 1: Starting
Task 2: Starting
Task 3: Starting
Task 1: Done
Task 3: Done
Task 2: Done
All tasks completed!

=== 示例 2: 并行处理切片 ===
Processed: apple
Processed: banana
Processed: cherry
Processed: date
Processed: elderberry
All items processed!

=== 示例 3: 限制并发数 ===
Processing item 1
Processing item 2
Processing item 3
Completed item 1
Processing item 4
Completed item 2
Processing item 5
...

🎉 所有示例运行完成!
```

---

## 最佳实践

### 1. 优先使用 WaitGroup.Go()

```go
// ✅ 推荐
wg.Go(func() {
    work()
})

// ❌ 不推荐 (除非有特殊需求)
wg.Add(1)
go func() {
    defer wg.Done()
    work()
}()
```

---

### 2. 结合 Go 1.22+ 循环变量特性

```go
// Go 1.22+ 自动捕获副本,无需传参
for _, item := range items {
    wg.Go(func() {
        process(item)  // ✅ 安全
    })
}
```

---

### 3. 限制并发数避免资源耗尽

```go
// 使用信号量
sem := make(chan struct{}, maxConcurrency)

for _, item := range items {
    sem <- struct{}{}
    wg.Go(func() {
        defer func() { <-sem }()
        process(item)
    })
}
```

---

### 4. 错误处理使用 errgroup

对于需要错误处理的场景,使用 `errgroup.Group`:

```go
import "golang.org/x/sync/errgroup"

g := new(errgroup.Group)

for _, item := range items {
    g.Go(func() error {
        return process(item)  // 可以返回错误
    })
}

if err := g.Wait(); err != nil {
    // 处理错误
}
```

---

## 常见问题

### Q: WaitGroup.Go() 是否线程安全?

**A**: ✅ 是的,完全线程安全。

---

### Q: 性能有影响吗?

**A**: 影响极小 (多一层闭包,~1% 开销),可以忽略。

---

### Q: 支持返回值吗?

**A**: ❌ 不支持。需要通过 channel 或 `errgroup.Group`。

---

## 相关资源

- 📘 [WaitGroup.Go() 技术文档](../01-WaitGroup-Go方法.md)
- 📘 [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup)
- 📘 [errgroup.Group](https://pkg.go.dev/golang.org/x/sync/errgroup)

---

**创建日期**: 2025年10月18日  
**作者**: AI Assistant

