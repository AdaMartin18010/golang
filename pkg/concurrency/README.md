# pkg/concurrency - 并发模式库

> **版本**: v2.0
> **Go版本**: 1.25.3+

本包提供了一系列实用的Go并发模式实现，帮助开发者更好地处理并发场景。

---

## 📦 包含的模式

### 1. Context传播模式 (context.go)

Context用于在goroutine之间传递取消信号、超时和请求范围的值。

**功能**:

- ✅ `WithTimeout` - 带超时的任务执行
- ✅ `WithCancel` - 可取消的任务执行
- ✅ `WithValue` - Context值传播
- ✅ `ContextAwarePipeline` - 支持Context的管道
- ✅ `ContextAwareWorkerPool` - 支持Context的Worker Pool

**示例**:

```go
import "github.com/yourusername/golang/pkg/concurrency/patterns"

// 带超时的任务
err := patterns.WithTimeout(context.Background(), 5*time.Second, func(ctx context.Context) error {
    // 你的任务逻辑
    return nil
})
```

---

### 2. Semaphore信号量 (semaphore.go)

用于限制并发访问的信号量实现。

**功能**:

- ✅ `Semaphore` - 基础信号量
- ✅ `WeightedSemaphore` - 加权信号量
- ✅ `ParallelExecuteWithLimit` - 限制并发数的并行执行

**示例**:

```go
// 创建信号量，最多5个并发
sem := patterns.NewSemaphore(5)

sem.Acquire()
defer sem.Release()

// 执行任务
doWork()
```

---

### 3. Rate Limiter限流器 (ratelimiter.go)

多种限流策略实现。

**功能**:

- ✅ `RateLimiter` - 令牌桶限流器
- ✅ `LeakyBucket` - 漏桶限流器
- ✅ `SlidingWindowLimiter` - 滑动窗口限流器

**示例**:

```go
// 创建限流器: 每秒100个请求，桶容量200
limiter := patterns.NewRateLimiter(100, 200)

if limiter.Allow() {
    // 处理请求
    handleRequest()
} else {
    // 拒绝请求
    rejectRequest()
}
```

---

### 4. Timeout超时控制 (timeout.go)

超时控制和断路器模式。

**功能**:

- ✅ `WithTimeoutFunc` - 带超时的函数执行
- ✅ `TimeoutRetry` - 带超时的重试机制
- ✅ `CircuitBreaker` - 断路器模式
- ✅ `BatchWithTimeout` - 批量任务超时控制

**示例**:

```go
// 5秒超时
result, err := patterns.WithTimeoutFunc(5*time.Second, func() (interface{}, error) {
    return heavyComputation()
})

// 断路器
cb := patterns.NewCircuitBreaker(3, 10*time.Second)
err := cb.Execute(func() error {
    return callExternalService()
})
```

---

### 5. Pipeline管道模式

数据流式处理。

**示例**:

```go
// 生成数据
nums := generateNumbers(1, 2, 3, 4, 5)

// 转换数据
squared := squareNumbers(nums)

// 消费数据
for n := range squared {
    fmt.Println(n)
}
```

---

### 6. Worker Pool工作池

固定数量的worker处理任务。

**示例**:

```go
jobs := make(chan int, 100)
results := make(chan int, 100)

// 启动5个workers
for w := 1; w <= 5; w++ {
    go worker(jobs, results)
}

// 发送任务
for j := 1; j <= 10; j++ {
    jobs <- j
}
close(jobs)

// 收集结果
for r := 1; r <= 10; r++ {
    result := <-results
    fmt.Println(result)
}
```

---

## 🚀 快速开始

### 安装

```bash
go get github.com/yourusername/golang/pkg/concurrency
```

### 使用

```go
import (
    "github.com/yourusername/golang/pkg/concurrency/patterns"
    "context"
    "time"
)

func main() {
    // 使用限流器
    limiter := patterns.NewRateLimiter(10, 20)

    if limiter.Allow() {
        // 处理请求
    }

    // 使用信号量
    sem := patterns.NewSemaphore(5)
    sem.Acquire()
    defer sem.Release()

    // 使用超时控制
    result, err := patterns.WithTimeoutFunc(5*time.Second, func() (interface{}, error) {
        return doSomething()
    })
}
```

---

## 📊 性能特点

### Semaphore

- **Acquire**: O(1)
- **Release**: O(1)
- **内存**: 每个信号量 ~200 bytes

### RateLimiter

- **Allow**: O(1)
- **内存**: ~300 bytes
- **适用**: 高QPS场景 (10K+ req/s)

### CircuitBreaker

- **Execute**: O(1)
- **内存**: ~150 bytes
- **适用**: 外部服务调用保护

---

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行基准测试
go test -bench=. -benchmem ./...

# 查看覆盖率
go test -cover ./...
```

---

## 📝 最佳实践

### 1. Context传播

- 总是传递Context
- 及时检查Context取消
- 使用defer cleanup

### 2. 信号量使用

- 使用defer确保释放
- 避免死锁
- 合理设置并发数

### 3. 限流策略

- 根据场景选择限流器
- 设置合理的速率
- 监控限流效果

### 4. 超时控制

- 设置合理的超时时间
- 处理超时错误
- 使用断路器保护外部调用

---

## 🔗 相关资源

- [Go并发模式](https://go.dev/blog/pipelines)
- [Context包](https://pkg.go.dev/context)
- [项目文档](../../docs/)

---

## 📞 问题反馈

遇到问题？欢迎提Issue或PR！

---

**维护者**: Project Team
**最后更新**: 2025-10-22
**License**: MIT
