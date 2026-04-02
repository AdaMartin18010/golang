# 任务补偿机制 (Task Compensation)

> **分类**: 工程与云原生  
> **标签**: #compensation #saga #distributed-transaction

---

## Saga 补偿模式

```go
// Saga 执行器
type Saga struct {
    steps []SagaStep
    ctx   context.Context
}

type SagaStep struct {
    Name       string
    Action     func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

func (s *Saga) Execute() error {
    completed := []int{}  // 记录已完成的步骤索引
    
    for i, step := range s.steps {
        if err := step.Action(s.ctx); err != nil {
            // 执行补偿
            return s.compensate(completed)
        }
        completed = append(completed, i)
    }
    
    return nil
}

func (s *Saga) compensate(completed []int) error {
    var errs []error
    
    // 逆序补偿
    for i := len(completed) - 1; i >= 0; i-- {
        stepIndex := completed[i]
        step := s.steps[stepIndex]
        
        if err := step.Compensate(s.ctx); err != nil {
            errs = append(errs, fmt.Errorf("compensate %s failed: %w", step.Name, err))
            // 记录补偿失败，需要人工介入
        }
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("compensation failed: %v", errs)
    }
    
    return nil
}

// 使用示例
func CreateOrderSaga(order Order) *Saga {
    return &Saga{
        steps: []SagaStep{
            {
                Name: "reserve_inventory",
                Action: func(ctx context.Context) error {
                    return inventoryService.Reserve(ctx, order.Items)
                },
                Compensate: func(ctx context.Context) error {
                    return inventoryService.Release(ctx, order.Items)
                },
            },
            {
                Name: "process_payment",
                Action: func(ctx context.Context) error {
                    return paymentService.Charge(ctx, order.Amount)
                },
                Compensate: func(ctx context.Context) error {
                    return paymentService.Refund(ctx, order.Amount)
                },
            },
            {
                Name: "create_shipment",
                Action: func(ctx context.Context) error {
                    return shippingService.Create(ctx, order)
                },
                Compensate: func(ctx context.Context) error {
                    return shippingService.Cancel(ctx, order)
                },
            },
        },
    }
}
```

---

## 补偿任务管理

```go
type CompensationManager struct {
    compensations []CompensationTask
    store         CompensationStore
}

type CompensationTask struct {
    ID          string
    OriginalTaskID string
    Action      func(ctx context.Context) error
    MaxRetries  int
    Status      CompensationStatus
}

func (cm *CompensationManager) Register(task CompensationTask) {
    cm.compensations = append(cm.compensations, task)
}

func (cm *CompensationManager) ExecuteAll(ctx context.Context) error {
    for _, comp := range cm.compensations {
        if err := cm.executeCompensation(ctx, comp); err != nil {
            return err
        }
    }
    return nil
}

func (cm *CompensationManager) executeCompensation(ctx context.Context, comp CompensationTask) error {
    var lastErr error
    
    for i := 0; i < comp.MaxRetries; i++ {
        if err := comp.Action(ctx); err != nil {
            lastErr = err
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }
        
        comp.Status = CompensationStatusSuccess
        return nil
    }
    
    comp.Status = CompensationStatusFailed
    
    // 记录需要人工处理的补偿失败
    cm.store.RecordFailedCompensation(ctx, comp, lastErr)
    
    return fmt.Errorf("compensation %s failed after %d retries: %w", 
        comp.ID, comp.MaxRetries, lastErr)
}
```

---

## 幂等补偿

```go
// 确保补偿操作幂等
type IdempotentCompensation struct {
    store IdempotencyStore
}

func (ic *IdempotentCompensation) Execute(ctx context.Context, key string, action func() error) error {
    // 检查是否已执行
    if status, _ := ic.store.Get(ctx, key); status == "completed" {
        return nil  // 已执行过
    }
    
    // 标记为进行中
    ic.store.Set(ctx, key, "in_progress")
    
    // 执行补偿
    if err := action(); err != nil {
        ic.store.Set(ctx, key, "failed")
        return err
    }
    
    // 标记为完成
    ic.store.Set(ctx, key, "completed")
    return nil
}

// 使用
compensation := &IdempotentCompensation{store: redisStore}

compensation.Execute(ctx, "refund:order-123", func() error {
    return paymentService.Refund(ctx, orderID, amount)
})
```

---

## 补偿监控

```go
type CompensationMonitor struct {
    metrics MetricsCollector
}

func (cm *CompensationMonitor) RecordCompensationStart(originalTaskID string) {
    cm.metrics.IncCounter("compensation_started")
    cm.metrics.GaugeInc("active_compensations")
}

func (cm *CompensationMonitor) RecordCompensationSuccess(originalTaskID string, duration time.Duration) {
    cm.metrics.IncCounter("compensation_succeeded")
    cm.metrics.GaugeDec("active_compensations")
    cm.metrics.RecordTiming("compensation_duration", duration)
}

func (cm *CompensationMonitor) RecordCompensationFailure(originalTaskID string, err error) {
    cm.metrics.IncCounter("compensation_failed")
    cm.metrics.GaugeDec("active_compensations")
    
    // 发送告警
    alertManager.Send(Alert{
        Severity: "critical",
        Message:  fmt.Sprintf("Compensation failed for task %s: %v", originalTaskID, err),
    })
}
```
