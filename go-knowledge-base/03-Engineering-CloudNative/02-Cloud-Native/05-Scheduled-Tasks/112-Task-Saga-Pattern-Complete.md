# Saga 模式完整实现 (Saga Pattern Complete Implementation)

> **分类**: 工程与云原生
> **标签**: #saga #distributed-transactions #compensation #orchestration
> **参考**: Saga Pattern (Hector Garcia-Molina), Temporal Saga, Axon Saga

---

## Saga 核心原理

```
传统分布式事务 (2PC)              Saga 模式
       │                            │
       ▼                            ▼
┌──────────────┐              ┌──────────────┐
│  Coordinator │              │   Saga       │
│  (全局锁)     │              │  (无锁+补偿)  │
└──────┬───────┘              └──────┬───────┘
       │                            │
   Prepare?                      Step 1 ✓
       │                            │
   Commit?                       Step 2 ✓
       │                            │
   Rollback?                     Step 3 ✗
       │                            │
       │                      Compensation 3
       │                      Compensation 2
       │                      Compensation 1
```

---

## 完整 Saga 编排器实现

```go
package saga

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "sync"
 "time"

 "github.com/google/uuid"
 "go.uber.org/zap"
)

// SagaState Saga状态
type SagaState int

const (
 SagaStatePending     SagaState = iota // 待执行
 SagaStateRunning                      // 执行中
 SagaStateCompleted                    // 完成
 SagaStateCompensating                 // 补偿中
 SagaStateCompensated                  // 已补偿
 SagaStateFailed                       // 失败（无法补偿）
)

// SagaDefinition Saga定义
type SagaDefinition struct {
 Name        string
 Description string
 Steps       []StepDefinition

 // 超时配置
 Timeout            time.Duration
 StepTimeout        time.Duration
 CompensationTimeout time.Duration

 // 并行度
 MaxParallelism int

 // 失败策略
 FailurePolicy FailurePolicy
}

// StepDefinition 步骤定义
type StepDefinition struct {
 Name            string
 Description     string

 // 执行动作
 Action      Action
 Compensate  Action

 // 依赖
 DependsOn   []string // 依赖的步骤名称

 // 重试
 MaxRetries  int
 RetryDelay  time.Duration

 // 超时
 Timeout     time.Duration

 // 并行组
 ParallelGroup string
}

// Action 动作函数
type Action func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)

// FailurePolicy 失败策略
type FailurePolicy int

const (
 FailFast FailurePolicy = iota        // 立即失败
 ContinueOnFailure                    // 继续执行
 RetryThenFail                        // 重试后失败
 ManualIntervention                   // 人工介入
)

// SagaInstance Saga实例
type SagaInstance struct {
 ID          string
 Definition  *SagaDefinition
 State       SagaState

 // 执行上下文
 Context     map[string]interface{}
 StepResults map[string]StepResult

 // 执行状态
 CurrentStep int
 StartedAt   time.Time
 CompletedAt *time.Time

 // 步骤执行状态
 stepStates  map[string]StepState
 stepMu      sync.RWMutex

 // 取消
 cancelFunc  context.CancelFunc
}

// StepState 步骤状态
type StepState struct {
 Name        string
 Status      StepStatus
 StartedAt   *time.Time
 CompletedAt *time.Time
 Result      map[string]interface{}
 Error       string
 Attempts    int
}

type StepStatus int

const (
 StepStatusPending StepStatus = iota
 StepStatusRunning
 StepStatusCompleted
 StepStatusFailed
 StepStatusCompensating
 StepStatusCompensated
)

// StepResult 步骤结果
type StepResult struct {
 Success bool
 Output  map[string]interface{}
 Error   error
}

// Orchestrator Saga编排器
type Orchestrator struct {
 definitions map[string]*SagaDefinition
 instances   map[string]*SagaInstance

 store   SagaStore
 logger  *zap.Logger

 // 执行
 workerPool chan struct{}
 mu         sync.RWMutex
}

// NewOrchestrator 创建编排器
func NewOrchestrator(store SagaStore, logger *zap.Logger, maxConcurrency int) *Orchestrator {
 return &Orchestrator{
  definitions: make(map[string]*SagaDefinition),
  instances:   make(map[string]*SagaInstance),
  store:       store,
  logger:      logger,
  workerPool:  make(chan struct{}, maxConcurrency),
 }
}

// Register 注册Saga定义
func (o *Orchestrator) Register(def *SagaDefinition) error {
 if err := o.validateDefinition(def); err != nil {
  return err
 }

 o.mu.Lock()
 defer o.mu.Unlock()

 o.definitions[def.Name] = def
 return nil
}

// validateDefinition 验证定义
func (o *Orchestrator) validateDefinition(def *SagaDefinition) error {
 if def.Name == "" {
  return errors.New("saga name is required")
 }

 if len(def.Steps) == 0 {
  return errors.New("saga must have at least one step")
 }

 // 检查步骤名称唯一性
 stepNames := make(map[string]bool)
 for _, step := range def.Steps {
  if step.Name == "" {
   return errors.New("step name is required")
  }
  if stepNames[step.Name] {
   return fmt.Errorf("duplicate step name: %s", step.Name)
  }
  stepNames[step.Name] = true

  // 检查依赖存在性
  for _, dep := range step.DependsOn {
   if !stepNames[dep] {
    return fmt.Errorf("step %s depends on non-existent step %s",
     step.Name, dep)
   }
  }
 }

 return nil
}

// Start 启动Saga
func (o *Orchestrator) Start(ctx context.Context, sagaName string,
 input map[string]interface{}) (*SagaInstance, error) {

 o.mu.RLock()
 def, exists := o.definitions[sagaName]
 o.mu.RUnlock()

 if !exists {
  return nil, fmt.Errorf("saga not found: %s", sagaName)
 }

 instance := &SagaInstance{
  ID:          uuid.New().String(),
  Definition:  def,
  State:       SagaStatePending,
  Context:     input,
  StepResults: make(map[string]StepResult),
  stepStates:  make(map[string]StepState),
  StartedAt:   time.Now(),
 }

 // 初始化步骤状态
 for _, step := range def.Steps {
  instance.stepStates[step.Name] = StepState{
   Name:   step.Name,
   Status: StepStatusPending,
  }
 }

 // 持久化
 if err := o.store.Save(ctx, instance); err != nil {
  return nil, err
 }

 o.mu.Lock()
 o.instances[instance.ID] = instance
 o.mu.Unlock()

 // 执行
 go o.execute(instance)

 return instance, nil
}

// execute 执行Saga
func (o *Orchestrator) execute(instance *SagaInstance) {
 instance.State = SagaStateRunning

 // 创建带超时的上下文
 ctx, cancel := context.WithTimeout(context.Background(),
  instance.Definition.Timeout)
 instance.cancelFunc = cancel
 defer cancel()

 // 获取可并行执行的步骤
 steps := o.getExecutableSteps(instance)

 // 执行步骤
 var wg sync.WaitGroup
 errChan := make(chan error, len(steps))

 for _, step := range steps {
  wg.Add(1)
  go func(s StepDefinition) {
   defer wg.Done()

   if err := o.executeStep(ctx, instance, s); err != nil {
    errChan <- err
   }
  }(step)
 }

 wg.Wait()
 close(errChan)

 // 检查是否有错误
 var execErr error
 for err := range errChan {
  if execErr == nil {
   execErr = err
  }
 }

 if execErr != nil {
  // 开始补偿
  o.compensate(instance)
  return
 }

 // 检查是否全部完成
 if o.isComplete(instance) {
  instance.State = SagaStateCompleted
  now := time.Now()
  instance.CompletedAt = &now
  o.store.Save(ctx, instance)
  o.logger.Info("saga completed",
   zap.String("saga_id", instance.ID),
   zap.String("definition", instance.Definition.Name))
 } else {
  // 继续执行依赖步骤
  go o.execute(instance)
 }
}

// executeStep 执行单个步骤
func (o *Orchestrator) executeStep(ctx context.Context, instance *SagaInstance,
 step StepDefinition) error {

 // 更新状态
 o.updateStepState(instance, step.Name, StepStatusRunning, nil)

 // 准备输入（合并上下文和依赖结果）
 input := o.prepareInput(instance, step)

 // 执行动作（带重试）
 var result map[string]interface{}
 var err error

 for attempt := 0; attempt <= step.MaxRetries; attempt++ {
  stepCtx, cancel := context.WithTimeout(ctx, step.Timeout)
  result, err = step.Action(stepCtx, input)
  cancel()

  if err == nil {
   break
  }

  if attempt < step.MaxRetries {
   time.Sleep(step.RetryDelay * time.Duration(attempt+1))
  }
 }

 if err != nil {
  o.updateStepState(instance, step.Name, StepStatusFailed, err)
  o.logger.Error("step failed",
   zap.String("saga_id", instance.ID),
   zap.String("step", step.Name),
   zap.Error(err))
  return err
 }

 // 保存结果
 o.updateStepState(instance, step.Name, StepStatusCompleted, result)
 instance.StepResults[step.Name] = StepResult{
  Success: true,
  Output:  result,
 }

 o.logger.Info("step completed",
  zap.String("saga_id", instance.ID),
  zap.String("step", step.Name))

 return nil
}

// compensate 执行补偿
func (o *Orchestrator) compensate(instance *SagaInstance) {
 instance.State = SagaStateCompensating

 ctx, cancel := context.WithTimeout(context.Background(),
  instance.Definition.CompensationTimeout)
 defer cancel()

 o.logger.Info("starting compensation",
  zap.String("saga_id", instance.ID))

 // 按相反顺序补偿
 for i := len(instance.Definition.Steps) - 1; i >= 0; i-- {
  step := instance.Definition.Steps[i]
  stepState := instance.stepStates[step.Name]

  // 只补偿已完成的步骤
  if stepState.Status != StepStatusCompleted {
   continue
  }

  // 更新状态
  o.updateStepState(instance, step.Name, StepStatusCompensating, nil)

  // 执行补偿
  input := o.prepareInput(instance, step)
  if _, err := step.Compensate(ctx, input); err != nil {
   o.logger.Error("compensation failed",
    zap.String("saga_id", instance.ID),
    zap.String("step", step.Name),
    zap.Error(err))

   instance.State = SagaStateFailed
   o.store.Save(ctx, instance)
   return
  }

  o.updateStepState(instance, step.Name, StepStatusCompensated, nil)
 }

 instance.State = SagaStateCompensated
 o.store.Save(ctx, instance)

 o.logger.Info("compensation completed",
  zap.String("saga_id", instance.ID))
}

// getExecutableSteps 获取可执行的步骤
func (o *Orchestrator) getExecutableSteps(instance *SagaInstance) []StepDefinition {
 var executable []StepDefinition

 for _, step := range instance.Definition.Steps {
  state := instance.stepStates[step.Name]

  // 跳过非pending状态
  if state.Status != StepStatusPending {
   continue
  }

  // 检查依赖
  depsSatisfied := true
  for _, dep := range step.DependsOn {
   depState := instance.stepStates[dep]
   if depState.Status != StepStatusCompleted {
    depsSatisfied = false
    break
   }
  }

  if depsSatisfied {
   executable = append(executable, step)
  }
 }

 return executable
}

// prepareInput 准备步骤输入
func (o *Orchestrator) prepareInput(instance *SagaInstance,
 step StepDefinition) map[string]interface{} {

 input := make(map[string]interface{})

 // 复制上下文
 for k, v := range instance.Context {
  input[k] = v
 }

 // 添加依赖结果
 for _, dep := range step.DependsOn {
  if result, ok := instance.StepResults[dep]; ok && result.Success {
   input[dep+"_result"] = result.Output
  }
 }

 return input
}

// updateStepState 更新步骤状态
func (o *Orchestrator) updateStepState(instance *SagaInstance,
 stepName string, status StepStatus, result interface{}) {

 instance.stepMu.Lock()
 defer instance.stepMu.Unlock()

 state := instance.stepStates[stepName]
 state.Status = status

 now := time.Now()
 switch status {
 case StepStatusRunning:
  state.StartedAt = &now
  state.Attempts++
 case StepStatusCompleted, StepStatusFailed, StepStatusCompensated:
  state.CompletedAt = &now
 }

 if result != nil {
  switch r := result.(type) {
  case map[string]interface{}:
   state.Result = r
  case error:
   state.Error = r.Error()
  }
 }

 instance.stepStates[stepName] = state
}

// isComplete 检查是否完成
func (o *Orchestrator) isComplete(instance *SagaInstance) bool {
 for _, step := range instance.Definition.Steps {
  state := instance.stepStates[step.Name]
  if state.Status != StepStatusCompleted {
   return false
  }
 }
 return true
}

// SagaStore 存储接口
type SagaStore interface {
 Save(ctx context.Context, instance *SagaInstance) error
 Get(ctx context.Context, id string) (*SagaInstance, error)
 List(ctx context.Context, state SagaState) ([]*SagaInstance, error)
}
```

---

## 实际使用示例

```go
// 订单处理 Saga
func CreateOrderSaga() *saga.SagaDefinition {
 return &saga.SagaDefinition{
  Name:        "CreateOrder",
  Description: "Create order with inventory and payment",
  Steps: []saga.StepDefinition{
   {
    Name:        "reserve_inventory",
    Description: "Reserve inventory",
    Action: func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
     productID := input["product_id"].(string)
     quantity := input["quantity"].(int)

     // 调用库存服务
     reservationID, err := inventoryClient.Reserve(ctx, productID, quantity)
     return map[string]interface{}{"reservation_id": reservationID}, err
    },
    Compensate: func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
     reservationID := input["reserve_inventory_result"].(map[string]interface{})["reservation_id"].(string)
     return nil, inventoryClient.Release(ctx, reservationID)
    },
    MaxRetries: 3,
    Timeout:    5 * time.Second,
   },
   {
    Name:        "process_payment",
    Description: "Process payment",
    DependsOn:   []string{"reserve_inventory"},
    Action: func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
     amount := input["amount"].(float64)

     // 调用支付服务
     transactionID, err := paymentClient.Charge(ctx, amount)
     return map[string]interface{}{"transaction_id": transactionID}, err
    },
    Compensate: func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
     transactionID := input["process_payment_result"].(map[string]interface{})["transaction_id"].(string)
     return nil, paymentClient.Refund(ctx, transactionID)
    },
    MaxRetries: 3,
    Timeout:    10 * time.Second,
   },
   {
    Name:        "create_shipment",
    Description: "Create shipment",
    DependsOn:   []string{"process_payment"},
    Action: func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
     address := input["address"].(string)

     shipmentID, err := shipmentClient.Create(ctx, address)
     return map[string]interface{}{"shipment_id": shipmentID}, err
    },
    Compensate: func(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
     shipmentID := input["create_shipment_result"].(map[string]interface{})["shipment_id"].(string)
     return nil, shipmentClient.Cancel(ctx, shipmentID)
    },
    Timeout: 5 * time.Second,
   },
  },
  Timeout:             30 * time.Second,
  StepTimeout:         10 * time.Second,
  CompensationTimeout: 60 * time.Second,
  FailurePolicy:       saga.RetryThenFail,
 }
}

// 使用
func PlaceOrder(ctx context.Context, orchestrator *saga.Orchestrator) error {
 instance, err := orchestrator.Start(ctx, "CreateOrder", map[string]interface{}{
  "product_id": "PROD-123",
  "quantity":   2,
  "amount":     199.99,
  "address":    "123 Main St",
 })

 if err != nil {
  return err
 }

 // 等待完成或失败
 return orchestrator.Wait(ctx, instance.ID)
}
```
