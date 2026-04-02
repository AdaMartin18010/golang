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
