# ANTIPATTERNS: 分布式系统反模式

> **维度**: Engineering-CloudNative
> **级别**: S (17+ KB)
> **标签**: #distributed-systems #antipatterns #reliability #scalability
> **权威来源**: [Release It!](https://pragprog.com/titles/mnee2/release-it-second-edition/) - Michael Nygard

---

## 反模式清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Distributed Systems Antipatterns                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  超时与重试                                                                 │
│  ├── ❌ 无超时 (No Timeout)                                                  │
│  ├── ❌ 快速重试 (Immediate Retry)                                          │
│  └── ❌ 无限重试 (Infinite Retry)                                           │
│                                                                              │
│  资源管理                                                                   │
│  ├── ❌ 连接泄漏 (Connection Leak)                                          │
│  ├── ❌ 无界队列 (Unbounded Queue)                                          │
│  └── ❌ 级联故障 (Cascading Failure)                                        │
│                                                                              │
│  数据一致性                                                                 │
│  ├── ❌ 分布式事务 (Distributed Transaction)                                 │
│  ├── ❌ 循环依赖 (Circular Dependency)                                      │
│  └── ❌ 共享数据库 (Shared Database)                                        │
│                                                                              │
│  服务设计                                                                   │
│  ├── ❌ 上帝服务 (God Service)                                              │
│  ├── ❌ 循环调用 (Circular Call)                                            │
│  └── ❌ 版本地狱 (Version Hell)                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. 无超时 (No Timeout)

### 问题

```go
// 反模式: 无超时
resp, err := http.Get("http://slow-service/api") // 可能永远等待!
if err != nil {
    return err
}
```

### 正解

```go
// 正确: 设置超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, "GET", "http://service/api", nil)
if err != nil {
    return err
}

resp, err := http.DefaultClient.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()
```

---

## 2. 快速重试 (Immediate Retry)

### 问题

```go
// 反模式: 立即重试
for i := 0; i < 3; i++ {
    resp, err := callService()  // 失败立即重试，导致服务雪崩
    if err == nil {
        return resp
    }
}
```

### 正解: 指数退避

```go
// 正确: 指数退避 + 抖动
func callWithRetry(fn func() error, maxRetries int) error {
    backoff := 100 * time.Millisecond
    maxBackoff := 10 * time.Second

    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }

        if i < maxRetries-1 {
            // 指数退避
            sleep := backoff * time.Duration(1<<i)
            if sleep > maxBackoff {
                sleep = maxBackoff
            }
            // 添加抖动避免共振
            jitter := time.Duration(rand.Int63n(int64(sleep) / 2))
            time.Sleep(sleep + jitter)
        }
    }
    return fmt.Errorf("max retries exceeded")
}
```

---

## 3. 级联故障 (Cascading Failure)

### 问题

```
用户请求 → API网关 → 订单服务 → 库存服务 (故障)
                            ↓
                     订单服务线程池耗尽
                            ↓
                     API网关线程池耗尽
                            ↓
                     整个系统不可用
```

### 正解: 熔断器

```go
// 正确: 熔断器保护
type CircuitBreaker struct {
    failures     int
    threshold    int
    timeout      time.Duration
    lastFailure  time.Time
    state        State // Closed, Open, HalfOpen
    mu           sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()

    // 检查熔断状态
    if cb.state == Open {
        if time.Since(cb.lastFailure) < cb.timeout {
            cb.mu.Unlock()
            return ErrCircuitOpen
        }
        cb.state = HalfOpen
    }

    cb.mu.Unlock()

    // 执行调用
    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        if cb.failures >= cb.threshold {
            cb.state = Open
        }
        return err
    }

    // 成功
    cb.failures = 0
    cb.state = Closed
    return nil
}
```

---

## 4. 分布式事务 (Distributed Transaction)

### 问题

```
反模式: 2PC (两阶段提交)
┌─────────┐               ┌─────────┐               ┌─────────┐
│ 协调者  │──Prepare─────►│ 服务A   │               │ 服务B   │
│         │◄────OK────────│         │               │         │
│         │──Prepare─────►│         │               │         │
│         │               │         │◄──────────────│         │
│         │               │         │───网络分区────►│         │
│         │ 阻塞!          │         │               │         │
└─────────┘               └─────────┘               └─────────┘
```

### 正解: Saga 模式

```go
// 正确: Saga 补偿事务
type Saga struct {
    steps []Step
}

type Step struct {
    Action    func() error
    Compensate func() error
}

func (s *Saga) Execute() error {
    completed := []int{}

    for i, step := range s.steps {
        if err := step.Action(); err != nil {
            // 补偿已完成的步骤
            for j := len(completed) - 1; j >= 0; j-- {
                s.steps[completed[j]].Compensate()
            }
            return err
        }
        completed = append(completed, i)
    }
    return nil
}

// 使用
saga := Saga{
    steps: []Step{
        {Action: createOrder, Compensate: cancelOrder},
        {Action: reserveInventory, Compensate: releaseInventory},
        {Action: processPayment, Compensate: refundPayment},
    },
}
```

---

## 5. 无界队列 (Unbounded Queue)

### 问题

```go
// 反模式: 无界队列
tasks := make(chan Task)  // 无缓冲，但生产者速度 > 消费者 → 内存溢出

go func() {
    for task := range tasks {
        process(task)  // 慢速处理
    }
}()

// 生产者快速推送
for {
    tasks <- newTask()  // 内存持续增长，最终 OOM
}
```

### 正解: 有界队列 + 背压

```go
// 正确: 有界队列
const maxQueueSize = 1000
tasks := make(chan Task, maxQueueSize)

// 生产者使用非阻塞发送
select {
case tasks <- task:
    // 成功入队
default:
    // 队列满，执行降级策略
    return ErrQueueFull
}

// 或: 超时入队
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

select {
case tasks <- task:
    // 成功
case <-ctx.Done():
    return ErrTimeout
}
```

---

## 6. 上帝服务 (God Service)

### 问题

```
反模式: 单服务处理所有业务
┌────────────────────────────────────────────────────────────┐
│                    OrderService (上帝服务)                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ 订单管理    │  │ 支付处理    │  │  库存管理            │ │
│  │ 优惠计算    │  │ 发票生成    │  │  物流跟踪            │ │
│  │ 用户积分    │  │ 退款处理    │  │  邮件通知            │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│                           50万行代码                        │
└────────────────────────────────────────────────────────────┘

问题:
- 部署困难 (任何改动都需全量发布)
- 团队冲突 (100+ 开发者)
- 扩展困难 (无法针对热点独立扩展)
- 故障隔离差 (一处故障影响全部)
```

### 正解: 领域拆分

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Order     │  │   Payment   │  │  Inventory  │  │Notification │
│   Service   │  │   Service   │  │   Service   │  │   Service   │
│  (订单领域)  │  │  (支付领域)  │  │  (库存领域)  │  │ (通知领域)  │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
      │                 │                │                │
      └─────────────────┴────────────────┴────────────────┘
                          │
                    Event Bus
```

---

## 反模式速查表

| 反模式 | 症状 | 正解 |
|--------|------|------|
| 无超时 | 请求卡住 | context.WithTimeout |
| 快速重试 | 服务雪崩 | 指数退避+抖动 |
| 无限重试 | 资源耗尽 | 最大重试次数+熔断 |
| 级联故障 | 全系统宕机 | 熔断器+舱壁模式 |
| 分布式事务 | 阻塞+不一致 | Saga模式 |
| 无界队列 | OOM | 有界队列+背压 |
| 上帝服务 | 部署困难 | 领域拆分 |
| 循环调用 | 死循环 | 重构+事件驱动 |

---

## 参考文献

1. [Release It!](https://pragprog.com/titles/mnee2/release-it-second-edition/) - Michael Nygard
2. [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman
3. [Designing Data-Intensive Applications](https://dataintensive.net/) - Martin Kleppmann
