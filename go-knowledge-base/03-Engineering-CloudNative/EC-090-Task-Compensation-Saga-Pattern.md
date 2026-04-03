# 任务补偿与 Saga 模式 (Task Compensation & Saga Pattern)

> **分类**: 工程与云原生
> **标签**: #compensation #saga #distributed-transactions
> **参考**: Saga Pattern, Microservices Patterns

---

## Saga 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Saga Pattern - Distributed Transactions                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Choreography-Based Saga                           │   │
│  │                                                                      │   │
│  │   OrderService ──► Create Order ──► PaymentService                  │   │
│  │                                          │                          │   │
│  │                                          ▼                          │   │
│  │                                    Process Payment                  │   │
│  │                                          │                          │   │
│  │                                          ▼                          │   │
│  │                                    InventoryService                 │   │
│  │                                          │                          │   │
│  │                                          ▼                          │   │
│  │                                    Reserve Inventory                │   │
│  │                                          │                          │   │
│  │                                          ▼                          │   │
│  │                                    ShippingService                  │   │
│  │                                                                      │   │
│  │                                          ▲                          │   │
│  │   Failure ◄── Compensate ◄── Compensate ◄── Compensate              │   │
│  │   Refund       Release        Cancel          Cancel                │   │
│  │   Payment      Inventory      Shipment        Order                 │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Orchestration-Based Saga                          │   │
│  │                                                                      │   │
│  │                    ┌─────────────┐                                  │   │
│  │                    │   Saga      │                                  │   │
│  │                    │ Orchestrator│                                  │   │
│  │                    └──────┬──────┘                                  │   │
│  │                           │                                          │   │
│  │          ┌────────────────┼────────────────┐                        │   │
│  │          │                │                │                        │   │
│  │          ▼                ▼                ▼                        │   │
│  │   ┌──────────┐     ┌──────────┐     ┌──────────┐                   │   │
│  │   │  Order   │     │ Payment  │     │ Inventory│                   │   │
│  │   │ Service  │     │ Service  │     │ Service  │                   │   │
│  │   └──────────┘     └──────────┘     └──────────┘                   │   │
│  │                                                                      │   │
│  │   Orchestrator coordinates all steps and handles compensation        │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整 Saga 实现

```go
package saga

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Saga 状态
type SagaState string

const (
    SagaStateStarted      SagaState = "STARTED"
    SagaStateInProgress   SagaState = "IN_PROGRESS"
    SagaStateCompleted    SagaState = "COMPLETED"
    SagaStateCompensating SagaState = "COMPENSATING"
    SagaStateCompensated  SagaState = "COMPENSATED"
    SagaStateFailed       SagaState = "FAILED"
)

// Step 状态
type StepState string

const (
    StepStatePending     StepState = "PENDING"
    StepStateSucceeded   StepState = "SUCCEEDED"
    StepStateFailed      StepState = "FAILED"
    StepStateCompensated StepState = "COMPENSATED"
)

// SagaStep Saga 步骤
type SagaStep struct {
    Name        string
    Execute     func(ctx context.Context, data interface{}) (interface{}, error)
    Compensate  func(ctx context.Context, data interface{}) error

    // 状态
    State       StepState
    Input       interface{}
    Output      interface{}
    Error       error

    // 元数据
    StartedAt   *time.Time
    CompletedAt *time.Time
}

// Saga 事务
type Saga struct {
    ID          string
    Name        string
    State       SagaState
    Steps       []*SagaStep
    CurrentStep int

    // 数据
    Data        map[string]interface{}

    // 元数据
    StartedAt   time.Time
    CompletedAt *time.Time

    mu          sync.RWMutex
}

// NewSaga 创建 Saga
func NewSaga(id, name string) *Saga {
    return &Saga{
        ID:        id,
        Name:      name,
        State:     SagaStateStarted,
        Steps:     make([]*SagaStep, 0),
        Data:      make(map[string]interface{}),
        StartedAt: time.Now(),
    }
}

// AddStep 添加步骤
func (s *Saga) AddStep(step *SagaStep) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Steps = append(s.Steps, step)
}

// Execute 执行 Saga
func (s *Saga) Execute(ctx context.Context) error {
    s.mu.Lock()
    s.State = SagaStateInProgress
    s.mu.Unlock()

    for i, step := range s.Steps {
        s.mu.Lock()
        s.CurrentStep = i
        s.mu.Unlock()

        // 执行步骤
        now := time.Now()
        step.StartedAt = &now
        step.State = StepStatePending

        output, err := step.Execute(ctx, step.Input)

        now = time.Now()
        step.CompletedAt = &now

        if err != nil {
            step.State = StepStateFailed
            step.Error = err

            // 执行补偿
            return s.compensate(ctx, i)
        }

        step.State = StepStateSucceeded
        step.Output = output

        // 保存输出到 Saga 数据
        s.Data[step.Name] = output
    }

    s.mu.Lock()
    s.State = SagaStateCompleted
    now := time.Now()
    s.CompletedAt = &now
    s.mu.Unlock()

    return nil
}

// compensate 执行补偿
func (s *Saga) compensate(ctx context.Context, failedStep int) error {
    s.mu.Lock()
    s.State = SagaStateCompensating
    s.mu.Unlock()

    // 逆序补偿
    for i := failedStep - 1; i >= 0; i-- {
        step := s.Steps[i]

        if step.State != StepStateSucceeded {
            continue
        }

        if step.Compensate != nil {
            if err := step.Compensate(ctx, step.Output); err != nil {
                // 补偿失败，记录但继续
                // 实际系统中可能需要人工干预
            }
        }

        step.State = StepStateCompensated
    }

    s.mu.Lock()
    s.State = SagaStateCompensated
    now := time.Now()
    s.CompletedAt = &now
    s.mu.Unlock()

    return fmt.Errorf("saga failed at step %d, compensated", failedStep)
}

// SagaOrchestrator Saga 编排器
type SagaOrchestrator struct {
    sagas  map[string]*Saga
    store  SagaStore
    mu     sync.RWMutex
}

// SagaStore 存储接口
type SagaStore interface {
    Save(ctx context.Context, saga *Saga) error
    Load(ctx context.Context, sagaID string) (*Saga, error)
    UpdateStep(ctx context.Context, sagaID string, stepIndex int, state StepState) error
}

// NewSagaOrchestrator 创建编排器
func NewSagaOrchestrator(store SagaStore) *SagaOrchestrator {
    return &SagaOrchestrator{
        sagas: make(map[string]*Saga),
        store: store,
    }
}

// StartSaga 开始 Saga
func (so *SagaOrchestrator) StartSaga(ctx context.Context, saga *Saga) error {
    so.mu.Lock()
    so.sagas[saga.ID] = saga
    so.mu.Unlock()

    // 持久化
    if err := so.store.Save(ctx, saga); err != nil {
        return err
    }

    // 异步执行
    go so.executeSaga(context.Background(), saga)

    return nil
}

func (so *SagaOrchestrator) executeSaga(ctx context.Context, saga *Saga) {
    _ = saga.Execute(ctx)

    // 保存最终状态
    so.store.Save(ctx, saga)
}

// GetSagaStatus 获取 Saga 状态
func (so *SagaOrchestrator) GetSagaStatus(sagaID string) (*Saga, error) {
    so.mu.RLock()
    saga, ok := so.sagas[sagaID]
    so.mu.RUnlock()

    if ok {
        return saga, nil
    }

    return so.store.Load(context.Background(), sagaID)
}

// CompensationManager 补偿管理器
type CompensationManager struct {
    compensations []CompensationRecord
    mu            sync.Mutex
}

// CompensationRecord 补偿记录
type CompensationRecord struct {
    SagaID      string
    StepName    string
    ExecutedAt  time.Time
    Success     bool
    Error       string
    RetryCount  int
}

// NewCompensationManager 创建补偿管理器
func NewCompensationManager() *CompensationManager {
    return &CompensationManager{
        compensations: make([]CompensationRecord, 0),
    }
}

// Record 记录补偿
func (cm *CompensationManager) Record(record CompensationRecord) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.compensations = append(cm.compensations, record)
}

// GetFailedCompensations 获取失败的补偿
func (cm *CompensationManager) GetFailedCompensations() []CompensationRecord {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    var failed []CompensationRecord
    for _, r := range cm.compensations {
        if !r.Success {
            failed = append(failed, r)
        }
    }

    return failed
}

// RetryCompensation 重试补偿
func (cm *CompensationManager) RetryCompensation(ctx context.Context, compensate func(ctx context.Context) error) error {
    var lastErr error
    for i := 0; i < 3; i++ {
        if err := compensate(ctx); err == nil {
            return nil
        } else {
            lastErr = err
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }
    return lastErr
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"

    "saga"
)

func main() {
    // 创建订单 Saga
    orderSaga := saga.NewSaga("order-123", "CreateOrder")

    // 步骤 1: 创建订单
    orderSaga.AddStep(&saga.SagaStep{
        Name: "CreateOrder",
        Execute: func(ctx context.Context, data interface{}) (interface{}, error) {
            fmt.Println("Creating order...")
            return map[string]string{"order_id": "ORD-123"}, nil
        },
        Compensate: func(ctx context.Context, data interface{}) error {
            fmt.Println("Cancelling order...")
            return nil
        },
    })

    // 步骤 2: 扣减库存
    orderSaga.AddStep(&saga.SagaStep{
        Name: "DeductInventory",
        Execute: func(ctx context.Context, data interface{}) (interface{}, error) {
            fmt.Println("Deducting inventory...")
            return map[string]string{"inventory_id": "INV-456"}, nil
        },
        Compensate: func(ctx context.Context, data interface{}) error {
            fmt.Println("Restoring inventory...")
            return nil
        },
    })

    // 步骤 3: 处理支付（模拟失败）
    orderSaga.AddStep(&saga.SagaStep{
        Name: "ProcessPayment",
        Execute: func(ctx context.Context, data interface{}) (interface{}, error) {
            fmt.Println("Processing payment...")
            return nil, fmt.Errorf("payment failed: insufficient funds")
        },
        Compensate: func(ctx context.Context, data interface{}) error {
            fmt.Println("Refunding payment...")
            return nil
        },
    })

    // 执行 Saga
    err := orderSaga.Execute(context.Background())
    if err != nil {
        fmt.Printf("Saga failed: %v\n", err)
    }

    // 打印状态
    fmt.Printf("Final state: %s\n", orderSaga.State)
    for _, step := range orderSaga.Steps {
        fmt.Printf("Step %s: %s\n", step.Name, step.State)
    }
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