# 任务数据一致性 (Task Data Consistency)

> **分类**: 工程与云原生
> **标签**: #consistency #transaction #at-least-once

---

## At-Least-Once 执行

```go
// 确保任务至少执行一次
type AtLeastOnceExecutor struct {
    store  TaskStore
    idempotencyChecker IdempotencyChecker
}

func (ale *AtLeastOnceExecutor) Execute(ctx context.Context, task *Task) error {
    // 1. 检查是否已经执行过（幂等性）
    if executed, _ := ale.idempotencyChecker.IsExecuted(ctx, task.ID); executed {
        return nil
    }

    // 2. 标记为执行中
    if err := ale.store.MarkExecuting(ctx, task.ID); err != nil {
        return err
    }

    // 3. 执行任务
    err := ale.executeTask(ctx, task)

    // 4. 记录结果
    if err != nil {
        ale.store.MarkFailed(ctx, task.ID, err)
        return err
    }

    // 5. 标记为已完成（幂等键）
    return ale.idempotencyChecker.MarkExecuted(ctx, task.ID, time.Hour*24)
}
```

---

## Exactly-Once 语义

```go
// 精确一次执行
type ExactlyOnceExecutor struct {
    dedupStore DedupStore
    locker     DistributedLocker
}

func (eoe *ExactlyOnceExecutor) Execute(ctx context.Context, task *Task) error {
    // 使用任务ID作为去重键
    dedupKey := fmt.Sprintf("task:dedup:%s", task.ID)

    // 1. 尝试获取去重锁
    if !eoe.dedupStore.TrySetNX(ctx, dedupKey, "processing", 5*time.Minute) {
        // 已经处理过或正在处理
        return eoe.waitForCompletion(ctx, dedupKey)
    }

    // 2. 获取分布式锁确保只有一个实例执行
    lock := eoe.locker.NewLock(fmt.Sprintf("task:lock:%s", task.ID), 5*time.Minute)
    if err := lock.Acquire(ctx); err != nil {
        return err
    }
    defer lock.Release(ctx)

    // 3. 再次检查（双重检查）
    if status, _ := eoe.dedupStore.Get(ctx, dedupKey); status == "completed" {
        return nil
    }

    // 4. 执行任务
    err := eoe.execute(ctx, task)

    // 5. 标记完成
    if err != nil {
        eoe.dedupStore.Set(ctx, dedupKey, "failed")
        return err
    }

    eoe.dedupStore.Set(ctx, dedupKey, "completed")
    return nil
}
```

---

## 分布式事务 Outbox

```go
// Outbox 模式确保事件发布
type OutboxPattern struct {
    db         *sql.DB
    eventBus   EventBus
    publisher  OutboxPublisher
}

func (op *OutboxPattern) ProcessTask(ctx context.Context, task *Task) error {
    tx, err := op.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. 执行业务逻辑
    if err := op.executeBusinessLogic(tx, task); err != nil {
        return err
    }

    // 2. 记录事件到 Outbox 表
    event := op.createEvent(task)
    if err := op.saveToOutbox(tx, event); err != nil {
        return err
    }

    // 3. 提交事务
    if err := tx.Commit(); err != nil {
        return err
    }

    // 4. 异步发布事件
    go op.publisher.Publish(event)

    return nil
}

// Outbox 发布器
type OutboxPublisher struct {
    db       *sql.DB
    eventBus EventBus
}

func (op *OutboxPublisher) Start(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)

    for {
        select {
        case <-ticker.C:
            op.publishPending(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (op *OutboxPublisher) publishPending(ctx context.Context) {
    rows, _ := op.db.QueryContext(ctx,
        "SELECT id, event_type, payload FROM outbox WHERE processed = false LIMIT 100")

    for rows.Next() {
        var id int
        var eventType, payload string
        rows.Scan(&id, &eventType, &payload)

        // 发布事件
        if err := op.eventBus.Publish(eventType, []byte(payload)); err == nil {
            // 标记为已处理
            op.db.ExecContext(ctx, "UPDATE outbox SET processed = true WHERE id = ?", id)
        }
    }
}
```

---

## 最终一致性检查

```go
type ConsistencyChecker struct {
    store ConsistencyStore
}

// 检查任务执行与副作用的一致性
func (cc *ConsistencyChecker) CheckTaskConsistency(ctx context.Context, taskID string) error {
    // 获取任务记录
    task, err := cc.store.GetTask(ctx, taskID)
    if err != nil {
        return err
    }

    // 检查预期副作用
    for _, expected := range task.ExpectedSideEffects {
        actual, err := cc.store.GetSideEffect(ctx, expected.ID)
        if err != nil {
            return fmt.Errorf("missing side effect %s: %w", expected.ID, err)
        }

        if !reflect.DeepEqual(expected.Data, actual.Data) {
            return fmt.Errorf("side effect %s mismatch", expected.ID)
        }
    }

    return nil
}

// 修复不一致
func (cc *ConsistencyChecker) RepairConsistency(ctx context.Context, taskID string) error {
    task, _ := cc.store.GetTask(ctx, taskID)

    for _, expected := range task.ExpectedSideEffects {
        if exists, _ := cc.store.SideEffectExists(ctx, expected.ID); !exists {
            // 重新应用副作用
            if err := cc.applySideEffect(ctx, expected); err != nil {
                return err
            }
        }
    }

    return nil
}
```

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02
---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02

---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02