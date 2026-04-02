# 真实世界案例研究 (Real-World Case Studies)

> **分类**: 工程与云原生
> **标签**: #case-study #production #lessons-learned
> **参考**: Uber Cadence, Netflix Conductor, Airbnb Chronos

---

## 案例1: Uber Cadence 工作流引擎

### 系统规模

- **日处理工作流**: 1亿+
- **并发执行**: 100万+
- **延迟要求**: P99 < 100ms

### 架构决策

```
Cadence Architecture:

Frontend Service (gRPC API)
    │
    ▼
Matching Service (Task Queue)
    │
    ▼
History Service (Event Sourcing)
    │
    ▼
Persistence (Cassandra + MySQL)
```

### 关键教训

**1. 事件溯源的权衡**

```go
// Cadence 历史记录存储优化
// 问题: 每个工作流执行产生大量事件，存储成本高
// 解决方案: 压缩历史记录

type HistoryCompressor struct {
    threshold int // 超过此阈值压缩
}

func (hc *HistoryCompressor) Compress(events []HistoryEvent) ([]byte, error) {
    // 使用 Protocol Buffers + Snappy 压缩
    // 压缩率: ~80%
    pbEvents := ToProto(events)
    data, err := proto.Marshal(pbEvents)
    if err != nil {
        return nil, err
    }

    return snappy.Encode(nil, data), nil
}
```

**2. 长运行工作流的处理**

```go
// 问题: 工作流运行数周/数月，进程重启怎么办？
// 解决方案: 状态持久化 + 惰性加载

type LongRunningWorkflow struct {
    state WorkflowState

    // 定期持久化检查点
    checkpointInterval time.Duration
}

func (lrw *LongRunningWorkflow) Run(ctx context.Context) error {
    ticker := time.NewTicker(lrw.checkpointInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            // 保存检查点
            return lrw.checkpoint()

        case <-ticker.C:
            if err := lrw.checkpoint(); err != nil {
                return err
            }

        default:
            // 处理下一个活动
            if err := lrw.processNextActivity(); err != nil {
                return err
            }
        }
    }
}
```

---

## 案例2: Netflix Conductor 编排引擎

### 挑战

- **峰值负载**: 300万任务/分钟
- **动态工作流**: 支持运行时修改
- **多租户**: 100+ 团队

### 解决方案

**1. 可扩展的任务队列**

```go
// Netflix 的 Dyno-queues (Redis 基础)
type DynoQueue struct {
    redis *redis.Client
    shard int // 分区键
}

func (dq *DynoQueue) Push(task *Task) error {
    // 一致性哈希分区
    shard := dq.getShard(task.WorkflowID)
    key := fmt.Sprintf("conductor.queue.%d", shard)

    // 使用 Redis Sorted Set 实现优先级队列
    return dq.redis.ZAdd(ctx, key, redis.Z{
        Score:  float64(task.Priority),
        Member: task.ID,
    }).Err()
}
```

**2. 动态工作流版本控制**

```go
// 支持工作流运行时更新
type VersionedWorkflow struct {
    Definition WorkflowDef
    Version    int

    // 迁移策略
    MigrationStrategy string // "strict" | "lenient"
}

func (vw *VersionedWorkflow) ValidateMigration(oldDef, newDef WorkflowDef) error {
    // 严格模式: 不允许删除步骤
    if vw.MigrationStrategy == "strict" {
        for _, step := range oldDef.Steps {
            if !newDef.HasStep(step.Name) {
                return fmt.Errorf("cannot remove step %s in strict mode", step.Name)
            }
        }
    }

    return nil
}
```

---

## 案例3: Airbnb Chronos 分布式调度

### 规模

- **定时任务**: 50万+
- **执行节点**: 200+
- **可用性**: 99.99%

### 架构要点

**1. 两阶段调度**

```go
// 解决: 任务调度与执行的解耦

type ChronosScheduler struct {
    // 阶段1: 计划调度
    planner *Planner

    // 阶段2: 执行调度
    executor *Executor
}

func (cs *ChronosScheduler) Schedule() {
    // 1. 计划器生成执行计划
    plan := cs.planner.GeneratePlan(time.Now())

    // 2. 执行器根据资源可用性调度
    for _, job := range plan.Jobs {
        // 检查资源
        if cs.executor.HasCapacity(job.ResourceRequirements) {
            cs.executor.Submit(job)
        } else {
            // 延迟到下一个周期
            cs.planner.Delay(job, time.Minute)
        }
    }
}
```

**2. 故障转移机制**

```go
// 领导者选举 + 热备
type ChronosFailover struct {
    leader     string
    followers  []string

    leaderLock *etcd.Lock
}

func (cf *ChronosFailover) Run(ctx context.Context) {
    // 尝试成为领导者
    if err := cf.electLeader(); err != nil {
        // 成为跟随者
        cf.runAsFollower(ctx)
        return
    }

    // 领导者模式
    cf.runAsLeader(ctx)
}

func (cf *ChronosFailover) runAsFollower(ctx context.Context) {
    // 监听领导者健康
    watch := cf.leaderLock.Watch()

    select {
    case <-watch.LeaderLost:
        // 领导者丢失，重新选举
        cf.Run(ctx)
    case <-ctx.Done():
        return
    }
}
```

---

## 通用设计模式总结

### 1. 背压控制

```go
// 防止系统过载
type BackpressureController struct {
    maxPendingTasks int
    currentTasks    int32
}

func (bc *BackpressureController) Admit() bool {
    current := atomic.LoadInt32(&bc.currentTasks)
    if int(current) >= bc.maxPendingTasks {
        return false
    }

    atomic.AddInt32(&bc.currentTasks, 1)
    return true
}
```

### 2. 断路器模式

```go
type CircuitBreaker struct {
    failureThreshold int
    successThreshold int
    timeout          time.Duration

    state    State
    failures int
    successes int
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == StateOpen {
        return errors.New("circuit breaker open")
    }

    err := fn()

    if err != nil {
        cb.recordFailure()
    } else {
        cb.recordSuccess()
    }

    return err
}
```

### 3. 批量处理优化

```go
// 减少数据库往返
type BatchProcessor struct {
    batchSize int
    buffer    []Task
    flushInterval time.Duration
}

func (bp *BatchProcessor) Add(task Task) {
    bp.buffer = append(bp.buffer, task)

    if len(bp.buffer) >= bp.batchSize {
        bp.Flush()
    }
}

func (bp *BatchProcessor) Flush() {
    if len(bp.buffer) == 0 {
        return
    }

    // 批量写入
    db.InsertBatch(bp.buffer)
    bp.buffer = bp.buffer[:0]
}
```

---

## 反模式教训

### 1. 数据库轮询

```go
// ❌ 错误: 频繁轮询数据库
for {
    tasks := db.Query("SELECT * FROM tasks WHERE status = 'pending'")
    for _, task := range tasks {
        process(task)
    }
    time.Sleep(100 * time.Millisecond) // 高负载
}

// ✅ 正确: 使用事件驱动 + 数据库变更流
for event := range db.ChangeStream("tasks") {
    if event.Status == "pending" {
        process(event.Task)
    }
}
```

### 2. 内存泄漏

```go
// ❌ 错误: 无限制的任务缓存
type TaskCache struct {
    tasks map[string]*Task // 永不清理
}

// ✅ 正确: LRU 缓存 + TTL
type TaskCache struct {
    cache *lru.Cache
    ttl   time.Duration
}
```

### 3. 级联故障

```go
// ❌ 错误: 同步调用依赖服务
func process(task *Task) error {
    result := syncCallDependency(task) // 阻塞
    return saveResult(result)
}

// ✅ 正确: 异步 + 超时 + 降级
func process(task *Task) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := asyncCallDependency(ctx, task)
    if err != nil {
        return fallback(task)
    }

    return saveResult(result)
}
```
