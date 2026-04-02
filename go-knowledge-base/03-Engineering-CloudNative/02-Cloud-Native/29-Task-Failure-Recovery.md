# 任务故障恢复 (Task Failure Recovery)

> **分类**: 工程与云原生
> **标签**: #failure-recovery #disaster-recovery #resilience

---

## 故障检测

```go
type FailureDetector struct {
    healthChecks []HealthCheck
    observers    []FailureObserver
}

type HealthCheck struct {
    Name      string
    Check     func(ctx context.Context) error
    Interval  time.Duration
    Timeout   time.Duration
    Threshold int  // 连续失败次数阈值
}

func (fd *FailureDetector) Start(ctx context.Context) {
    for _, check := range fd.healthChecks {
        go fd.runCheck(ctx, check)
    }
}

func (fd *FailureDetector) runCheck(ctx context.Context, check HealthCheck) {
    ticker := time.NewTicker(check.Interval)
    defer ticker.Stop()

    failures := 0

    for {
        select {
        case <-ticker.C:
            checkCtx, cancel := context.WithTimeout(ctx, check.Timeout)
            err := check.Check(checkCtx)
            cancel()

            if err != nil {
                failures++
                if failures >= check.Threshold {
                    fd.notifyFailure(check.Name, err)
                }
            } else {
                if failures >= check.Threshold {
                    fd.notifyRecovery(check.Name)
                }
                failures = 0
            }
        case <-ctx.Done():
            return
        }
    }
}
```

---

## 自动恢复策略

```go
type RecoveryStrategy interface {
    CanRecover(task *Task, err error) bool
    Recover(ctx context.Context, task *Task) error
}

// 重启恢复
type RestartRecovery struct{}

func (rr *RestartRecovery) CanRecover(task *Task, err error) bool {
    // 只有特定错误可以恢复
    var retryable RetryableError
    return errors.As(err, &retryable)
}

func (rr *RestartRecovery) Recover(ctx context.Context, task *Task) error {
    // 清理状态
    if err := cleanupTaskState(task); err != nil {
        return err
    }

    // 重新提交任务
    return scheduler.Reschedule(task)
}

// 跳过恢复
type SkipRecovery struct{}

func (sr *SkipRecovery) CanRecover(task *Task, err error) bool {
    // 非关键任务可以跳过
    return task.Priority == PriorityLow
}

func (sr *SkipRecovery) Recover(ctx context.Context, task *Task) error {
    task.Status = TaskStatusSkipped
    return store.Update(ctx, task)
}

// 降级恢复
type DegradedRecovery struct {
    fallbackHandler TaskHandler
}

func (dr *DegradedRecovery) CanRecover(task *Task, err error) bool {
    return dr.fallbackHandler != nil
}

func (dr *DegradedRecovery) Recover(ctx context.Context, task *Task) error {
    // 使用降级逻辑
    return dr.fallbackHandler.Handle(ctx, task)
}

// 恢复管理器
type RecoveryManager struct {
    strategies []RecoveryStrategy
}

func (rm *RecoveryManager) AttemptRecovery(ctx context.Context, task *Task, err error) error {
    for _, strategy := range rm.strategies {
        if strategy.CanRecover(task, err) {
            if recoveryErr := strategy.Recover(ctx, task); recoveryErr == nil {
                return nil  // 恢复成功
            }
        }
    }

    return fmt.Errorf("all recovery strategies failed for task %s", task.ID)
}
```

---

## 检查点恢复

```go
type Checkpointer struct {
    store CheckpointStore
    interval time.Duration
}

type Checkpoint struct {
    TaskID    string
    State     interface{}
    Progress  float64
    Timestamp time.Time
}

func (cp *Checkpointer) SaveCheckpoint(ctx context.Context, taskID string, state interface{}, progress float64) error {
    checkpoint := Checkpoint{
        TaskID:    taskID,
        State:     state,
        Progress:  progress,
        Timestamp: time.Now(),
    }

    return cp.store.Save(ctx, checkpoint)
}

func (cp *Checkpointer) RestoreFromCheckpoint(ctx context.Context, taskID string) (*Checkpoint, error) {
    return cp.store.GetLatest(ctx, taskID)
}

// 使用检查点的任务
func processWithCheckpoint(ctx context.Context, task *Task, checkpointer *Checkpointer) error {
    // 尝试从检查点恢复
    if checkpoint, err := checkpointer.RestoreFromCheckpoint(ctx, task.ID); err == nil {
        // 从检查点恢复状态
        task.State = checkpoint.State
        task.Progress = checkpoint.Progress
    }

    // 处理剩余工作
    for i := int(task.Progress); i < len(task.Items); i++ {
        // 处理每个项目
        if err := processItem(task.Items[i]); err != nil {
            // 保存检查点
            checkpointer.SaveCheckpoint(ctx, task.ID, task.State, float64(i))
            return err
        }

        // 定期保存检查点
        if i%100 == 0 {
            checkpointer.SaveCheckpoint(ctx, task.ID, task.State, float64(i))
        }
    }

    return nil
}
```

---

## 灾难恢复

```go
type DisasterRecovery struct {
    backupManager *BackupManager
    clusterManager *ClusterManager
}

func (dr *DisasterRecovery) ExecuteDRPlan(ctx context.Context) error {
    // 1. 评估影响
    impact := dr.assessImpact()

    // 2. 通知相关人员
    dr.notifyStakeholders(impact)

    // 3. 切换到备用集群
    if err := dr.clusterManager.Failover(ctx); err != nil {
        return fmt.Errorf("failover failed: %w", err)
    }

    // 4. 恢复数据
    if err := dr.backupManager.RestoreLatest(ctx); err != nil {
        return fmt.Errorf("restore failed: %w", err)
    }

    // 5. 恢复任务状态
    if err := dr.recoverTaskStates(ctx); err != nil {
        return fmt.Errorf("task recovery failed: %w", err)
    }

    // 6. 验证恢复
    if err := dr.verifyRecovery(ctx); err != nil {
        return fmt.Errorf("verification failed: %w", err)
    }

    return nil
}

func (dr *DisasterRecovery) recoverTaskStates(ctx context.Context) error {
    // 获取所有进行中任务
    tasks, _ := dr.store.ListRunningTasks(ctx)

    for _, task := range tasks {
        // 检查是否需要重新执行
        if time.Since(task.LastHeartbeat) > 5*time.Minute {
            // 任务可能已死，重置为待执行
            task.Status = TaskStatusPending
            dr.store.Update(ctx, task)
        }
    }

    return nil
}
```
