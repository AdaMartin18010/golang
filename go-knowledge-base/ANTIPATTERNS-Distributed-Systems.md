# 反模式：分布式系统 (Antipatterns: Distributed Systems)

> **分类**: 工程反模式
> **相关文档**: EC-042, EC-112, FT-004

---

## 1. 超时设置不当

### 反模式

```go
// 固定超时
timeout := 5 * time.Second  // 对所有请求都一样
```

### 问题

- 慢请求拖垮系统
- 级联故障

### 正解

```go
// 自适应超时 + 熔断
client := &http.Client{
    Timeout: calculateAdaptiveTimeout(),
}

// 配合熔断器
if cb.Allow() {
    resp, err := client.Do(req)
    cb.Record(err)
}
```

---

## 2. 忽略网络分区

### 反模式

```go
// 假设网络总是可用
if err := callServiceB(); err != nil {
    return err  // 直接返回错误
}
```

### 问题

- 分区时系统不可用
- 数据不一致

### 正解

```go
// 降级策略
if err := callServiceB(); err != nil {
    // 使用缓存
    return getFromCache(key)

    // 或返回默认值
    return defaultResponse()

    // 或进入只读模式
    return errServiceUnavailable
}
```

---

## 3. 分布式事务 (2PC)

### 反模式

```go
// 使用 2PC 处理所有跨服务事务
// 问题：性能差、单点故障、阻塞
```

### 正解

```go
// 使用 Saga 模式
saga := NewSaga()
saga.AddStep(reserveInventory)
saga.AddStep(processPayment)
saga.AddStep(shipOrder)
// 每个步骤独立，失败时补偿
```

---

## 4. 无界重试

### 反模式

```go
for {
    err := callService()  // 无限重试
    if err == nil {
        break
    }
}
```

### 问题

- 雪崩效应
- 资源耗尽

### 正解

```go
// 指数退避 + 最大重试
retry.WithMaxRetries(3,
    retry.WithBackoff(retry.ExponentialBackoff(
        100*time.Millisecond,
        2.0,  // multiplier
        1*time.Second,  // max
    )),
    func() error {
        return callService()
    },
)
```

---

## 5. 共享数据库

### 反模式

```
Service A ──┐
            ├──► Shared DB
Service B ──┘
```

### 问题

- 紧耦合
- 难以扩展
- schema 冲突

### 正解

```
Service A ──► DB A
Service B ──► DB B
            │
            └── 通过 API/Event 通信
```

---

## 6. 同步调用链

### 反模式

```
A ──sync──► B ──sync──► C ──sync──► D
总延迟 = A+B+C+D
```

### 正解

```
// 异步 + 缓存 + 批处理
A ──async──► B
         │
         └──► Cache

Batch Process ──► C
```

---

## 反模式清单

| 反模式 | 影响 | 解决方案 |
|--------|------|---------|
| 硬编码超时 | 级联故障 | 自适应超时 |
| 忽略分区 | 不可用 | 降级策略 |
| 2PC | 性能差 | Saga |
| 无限重试 | 雪崩 | 指数退避 |
| 共享DB | 紧耦合 | 服务自有数据 |
| 同步链 | 高延迟 | 异步解耦 |
| 无熔断 | 级联失败 | Circuit Breaker |
| 无幂等 | 重复处理 | 幂等键 |
| 大事务 | 锁竞争 | 小事务 |
| 无监控 | 黑盒 | 可观测性 |

---

## 参考文献

1. [Fallacies of Distributed Computing](https://en.wikipedia.org/wiki/Fallacies_of_distributed_computing)
2. [Release It!](https://pragprog.com/titles/mnee2/release-it-second-edition/) - Michael Nygard
3. [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman
