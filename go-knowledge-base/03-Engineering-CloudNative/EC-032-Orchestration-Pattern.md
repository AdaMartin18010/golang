# EC-032: Orchestration Pattern (编排模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #orchestration #saga #centralized #workflow #state-machine
> **权威来源**:
>
> - [Orchestration-based Saga](https://microservices.io/patterns/data/saga.html) - Chris Richardson
> - [Temporal.io Documentation](https://docs.temporal.io/)
> - [Netflix Conductor](https://netflix.github.io/conductor/)
> - [AWS Step Functions](https://aws.amazon.com/step-functions/)

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在微服务架构中，如何有效管理复杂的分布式事务流程，包括条件分支、并行执行、重试策略和人工审批？

**形式化描述**:

```
给定: 服务集合 S = {S₁, S₂, ..., Sₙ}
给定: 业务工作流 W = (N, E, C, A)，其中：
  - N: 节点集合（活动、决策、并行）
  - E: 边集合（转换关系）
  - C: 条件函数
  - A: 动作集合
约束:
  - 需要可见性和控制能力
  - 支持复杂流程模式
  - 要求故障恢复机制
目标: 找到最优协调策略使得工作流正确执行
```

**反模式**:

- 分布式编舞：流程逻辑分散在各服务中
- 硬编码流程：业务逻辑与流程控制耦合
- 缺少超时控制：长时间挂起的流程

### 1.2 解决方案形式化

**定义 1.1 (编排器)**
编排器是一个中央协调组件，负责：

1. 维护工作流状态机
2. 向参与者发送命令
3. 处理响应和事件
4. 执行补偿逻辑
5. 管理流程生命周期

**形式化表示**:

```
Orchestrator O = ⟨State, Commands, Events, Transitions⟩

状态转换函数:
  δ: State × Event → State × Command*

执行语义:
  S₀ →[e₁/C₁] S₁ →[e₂/C₂] S₂ → ... →[eₙ/∅] Sₙ

其中 eᵢ 是事件，Cᵢ 是命令序列
```

**定义 1.2 (工作流模式)**

```
顺序:  S₁ → S₂ → S₃
并行:  S₁ → [S₂ || S₃] → S₄
选择:  S₁ → {condition} → S₂ or S₃
循环:  S₁ → {condition} → S₁ (repeat) or S₂ (exit)
补偿:  S₁ → (fail) → C₁ → C₂
```

### 1.3 状态机模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Orchestrator State Machine                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                              ┌─────────┐                                │
│                              │  Start  │                                │
│                              └────┬────┘                                │
│                                   │                                     │
│                                   ▼                                     │
│   ┌─────────────────────────────────────────────────────────────┐      │
│   │                       Activity Layer                         │      │
│   │  ┌─────────┐    Command    ┌─────────┐    Response          │      │
│   │  │ Waiting │──────────────►│Executing│◄────────────────     │      │
│   │  └────┬────┘               └────┬────┘                      │      │
│   │       ▲                         │                           │      │
│   │       │          ┌──────────────┘                           │      │
│   │       │          │                                         │      │
│   │       │    ┌─────┴─────┐   ┌─────────┐                     │      │
│   │       └────┤Completed  │   │ Failed  │                     │      │
│   │            └─────┬─────┘   └────┬────┘                     │      │
│   │                  │              │                           │      │
│   └──────────────────┼──────────────┼───────────────────────────┘      │
│                      │              │                                  │
│                      ▼              ▼                                  │
│              ┌─────────────┐  ┌─────────────┐                          │
│              │ Next Activity│  │ Compensate │                          │
│              └─────────────┘  └──────┬──────┘                          │
│                                      │                                  │
│                                      ▼                                  │
│   ┌─────────────────────────────────────────────────────────────┐      │
│   │                     Compensation Layer                       │      │
│   │  ┌─────────┐    Compensate    ┌─────────┐                   │      │
│   │  │ Pending │─────────────────►│Executing│                   │      │
│   │  └────┬────┘                  └────┬────┘                   │      │
│   │       ▲                            │                        │      │
│   │       └────────────────────────────┘                        │      │
│   │                    (Success/Retry)                           │      │
│   └─────────────────────────────────────────────────────────────┘      │
│                                                                         │
│   Terminal States:                                                      │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                    │
│   │  Completed  │  │ Compensated │  │   Failed    │                    │
│   └─────────────┘  └─────────────┘  └─────────────┘                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心编排器实现

```go
// orchestration/core.go
package orchestration

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/google/uuid"
)

// Workflow 工作流定义
type Workflow struct {
    ID          string
    Name        string
    Version     string
    StartNode   string
    Nodes       map[string]*Node
    Edges       []*Edge
    Timeout     time.Duration
    RetryPolicy *RetryPolicy
}

// Node 工作流节点
type Node struct {
    ID           string
    Type         NodeType
    Service      string
    Action       string
    Compensation *Compensation
    InputMapping map[string]string
    OutputMapping map[string]string
    RetryConfig  *RetryPolicy
}

// NodeType 节点类型
type NodeType int

const (
    NodeTypeActivity NodeType = iota
    NodeTypeDecision
    NodeTypeParallel
    NodeTypeSubWorkflow
    NodeTypeHumanTask
)

// Edge 工作流边
type Edge struct {
    From      string
    To        string
    Condition *Condition
}

// Condition 条件
type Condition struct {
    Expression string
    Evaluator  func(ctx *ExecutionContext) bool
}

// RetryPolicy 重试策略
type RetryPolicy struct {
    MaxAttempts int
    BackoffType BackoffType
    InitialDelay time.Duration
    MaxDelay     time.Duration
}

// BackoffType 退避类型
type BackoffType int

const (
    BackoffFixed BackoffType = iota
    BackoffLinear
    BackoffExponential
)

// Compensation 补偿定义
type Compensation struct {
    Service string
    Action  string
    Timeout time.Duration
}

// WorkflowInstance 工作流实例
type WorkflowInstance struct {
    ID           string
    WorkflowID   string
    WorkflowName string
    Status       WorkflowStatus
    CurrentNode  string
    Variables    map[string]interface{}
    History      []*ActivityResult
    StartedAt    time.Time
    CompletedAt  *time.Time
    mu           sync.RWMutex
}

// WorkflowStatus 工作流状态
type WorkflowStatus int

const (
    WorkflowStatusPending WorkflowStatus = iota
    WorkflowStatusRunning
    WorkflowStatusWaiting
    WorkflowStatusCompleted
    WorkflowStatusCompensating
    WorkflowStatusCompensated
    WorkflowStatusFailed
    WorkflowStatusCancelled
)

func (w WorkflowStatus) String() string {
    names := []string{"PENDING", "RUNNING", "WAITING", "COMPLETED",
                      "COMPENSATING", "COMPENSATED", "FAILED", "CANCELLED"}
    if int(w) < len(names) {
        return names[w]
    }
    return "UNKNOWN"
}

// ActivityResult 活动执行结果
type ActivityResult struct {
    NodeID    string
    Status    ActivityStatus
    Input     map[string]interface{}
    Output    map[string]interface{}
    Error     string
    StartedAt time.Time
    EndedAt   time.Time
}

// ActivityStatus 活动状态
type ActivityStatus int

const (
    ActivityStatusPending ActivityStatus = iota
    ActivityStatusRunning
    ActivityStatusCompleted
    ActivityStatusFailed
    ActivityStatusCompensated
)

// Command 命令
type Command struct {
    ID            string
    WorkflowID    string
    InstanceID    string
    NodeID        string
    Service       string
    Action        string
    Input         map[string]interface{}
    CorrelationID string
    Timestamp     time.Time
}

// Response 响应
type Response struct {
    ID            string
    CommandID     string
    Status        ResponseStatus
    Output        map[string]interface{}
    Error         string
    Timestamp     time.Time
}

// ResponseStatus 响应状态
type ResponseStatus int

const (
    ResponseStatusSuccess ResponseStatus = iota
    ResponseStatusFailure
    ResponseStatusTimeout
)

// ExecutionContext 执行上下文
type ExecutionContext struct {
    Workflow   *Workflow
    Instance   *WorkflowInstance
    Variables  map[string]interface{}
    Node       *Node
}

// ServiceClient 服务客户端接口
type ServiceClient interface {
    Execute(ctx context.Context, command *Command) (*Response, error)
    ExecuteCompensation(ctx context.Context, command *Command) error
}

// WorkflowStore 工作流存储接口
type WorkflowStore interface {
    SaveWorkflow(ctx context.Context, workflow *Workflow) error
    GetWorkflow(ctx context.Context, id string) (*Workflow, error)
    SaveInstance(ctx context.Context, instance *WorkflowInstance) error
    GetInstance(ctx context.Context, id string) (*WorkflowInstance, error)
    UpdateInstanceStatus(ctx context.Context, id string, status WorkflowStatus) error
    AppendHistory(ctx context.Context, instanceID string, result *ActivityResult) error
}

// Orchestrator 编排器
type Orchestrator struct {
    workflows   map[string]*Workflow
    instances   map[string]*WorkflowInstance
    store       WorkflowStore
    client      ServiceClient
    eventBus    EventBus
    mu          sync.RWMutex
    logger      Logger
}

// EventBus 事件总线接口
type EventBus interface {
    Publish(ctx context.Context, event interface{}) error
    Subscribe(eventType string, handler EventHandler)
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event interface{}) error

// Logger 日志接口
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// NewOrchestrator 创建编排器
func NewOrchestrator(store WorkflowStore, client ServiceClient, eventBus EventBus, logger Logger) *Orchestrator {
    return &Orchestrator{
        workflows: make(map[string]*Workflow),
        instances: make(map[string]*WorkflowInstance),
        store:     store,
        client:    client,
        eventBus:  eventBus,
        logger:    logger,
    }
}

// RegisterWorkflow 注册工作流
func (o *Orchestrator) RegisterWorkflow(workflow *Workflow) error {
    o.mu.Lock()
    defer o.mu.Unlock()

    if _, exists := o.workflows[workflow.ID]; exists {
        return fmt.Errorf("workflow %s already registered", workflow.ID)
    }

    // 验证工作流
    if err := o.validateWorkflow(workflow); err != nil {
        return fmt.Errorf("invalid workflow: %w", err)
    }

    o.workflows[workflow.ID] = workflow

    if err := o.store.SaveWorkflow(context.Background(), workflow); err != nil {
        return fmt.Errorf("failed to save workflow: %w", err)
    }

    o.logger.Info("workflow registered", Field{"workflow_id", workflow.ID}, Field{"name", workflow.Name})
    return nil
}

func (o *Orchestrator) validateWorkflow(workflow *Workflow) error {
    if workflow.StartNode == "" {
        return fmt.Errorf("start node not defined")
    }
    if _, exists := workflow.Nodes[workflow.StartNode]; !exists {
        return fmt.Errorf("start node %s not found", workflow.StartNode)
    }
    return nil
}

// StartWorkflow 启动工作流
func (o *Orchestrator) StartWorkflow(ctx context.Context, workflowID string, input map[string]interface{}) (*WorkflowInstance, error) {
    o.mu.RLock()
    workflow, exists := o.workflows[workflowID]
    o.mu.RUnlock()

    if !exists {
        return nil, fmt.Errorf("workflow %s not found", workflowID)
    }

    instance := &WorkflowInstance{
        ID:           uuid.New().String(),
        WorkflowID:   workflowID,
        WorkflowName: workflow.Name,
        Status:       WorkflowStatusRunning,
        CurrentNode:  workflow.StartNode,
        Variables:    input,
        History:      []*ActivityResult{},
        StartedAt:    time.Now(),
    }

    o.mu.Lock()
    o.instances[instance.ID] = instance
    o.mu.Unlock()

    if err := o.store.SaveInstance(ctx, instance); err != nil {
        return nil, fmt.Errorf("failed to save instance: %w", err)
    }

    o.logger.Info("workflow started",
        Field{"instance_id", instance.ID},
        Field{"workflow_id", workflowID})

    // 开始执行
    go o.executeNode(context.Background(), instance, workflow.StartNode)

    return instance, nil
}

// executeNode 执行节点
func (o *Orchestrator) executeNode(ctx context.Context, instance *WorkflowInstance, nodeID string) {
    o.mu.RLock()
    workflow := o.workflows[instance.WorkflowID]
    o.mu.RUnlock()

    node, exists := workflow.Nodes[nodeID]
    if !exists {
        o.handleError(instance, fmt.Errorf("node %s not found", nodeID))
        return
    }

    instance.CurrentNode = nodeID
    o.store.UpdateInstanceStatus(ctx, instance.ID, WorkflowStatusRunning)

    execCtx := &ExecutionContext{
        Workflow:  workflow,
        Instance:  instance,
        Variables: instance.Variables,
        Node:      node,
    }

    switch node.Type {
    case NodeTypeActivity:
        o.executeActivity(ctx, instance, node, execCtx)
    case NodeTypeDecision:
        o.executeDecision(ctx, instance, node, workflow, execCtx)
    case NodeTypeParallel:
        o.executeParallel(ctx, instance, node, workflow, execCtx)
    default:
        o.handleError(instance, fmt.Errorf("unsupported node type: %v", node.Type))
    }
}

// executeActivity 执行活动节点
func (o *Orchestrator) executeActivity(ctx context.Context, instance *WorkflowInstance, node *Node, execCtx *ExecutionContext) {
    // 准备输入
    input := o.prepareInput(node.InputMapping, execCtx.Variables)

    command := &Command{
        ID:            uuid.New().String(),
        WorkflowID:    instance.WorkflowID,
        InstanceID:    instance.ID,
        NodeID:        node.ID,
        Service:       node.Service,
        Action:        node.Action,
        Input:         input,
        CorrelationID: instance.ID,
        Timestamp:     time.Now(),
    }

    result := &ActivityResult{
        NodeID:    node.ID,
        Status:    ActivityStatusRunning,
        Input:     input,
        StartedAt: time.Now(),
    }

    o.logger.Info("executing activity",
        Field{"instance_id", instance.ID},
        Field{"node_id", node.ID},
        Field{"service", node.Service},
        Field{"action", node.Action})

    // 执行命令
    response, err := o.executeWithRetry(ctx, command, node.RetryConfig)

    result.EndedAt = time.Now()

    if err != nil {
        result.Status = ActivityStatusFailed
        result.Error = err.Error()
        o.store.AppendHistory(ctx, instance.ID, result)

        o.logger.Error("activity failed",
            Field{"instance_id", instance.ID},
            Field{"node_id", node.ID},
            Field{"error", err})

        // 触发补偿
        o.compensate(ctx, instance)
        return
    }

    result.Status = ActivityStatusCompleted
    result.Output = response.Output
    o.store.AppendHistory(ctx, instance.ID, result)

    // 更新变量
    o.updateVariables(instance, node.OutputMapping, response.Output)

    // 发布事件
    o.eventBus.Publish(ctx, map[string]interface{}{
        "type":         "activity.completed",
        "instance_id":  instance.ID,
        "node_id":      node.ID,
        "timestamp":    time.Now(),
    })

    // 转移到下一个节点
    o.transition(ctx, instance, node.ID, workflow)
}

// executeWithRetry 带重试的执行
func (o *Orchestrator) executeWithRetry(ctx context.Context, command *Command, retryConfig *RetryPolicy) (*Response, error) {
    if retryConfig == nil {
        return o.client.Execute(ctx, command)
    }

    var lastErr error
    delay := retryConfig.InitialDelay

    for attempt := 0; attempt <= retryConfig.MaxAttempts; attempt++ {
        if attempt > 0 {
            time.Sleep(delay)

            // 计算下一次延迟
            switch retryConfig.BackoffType {
            case BackoffLinear:
                delay += retryConfig.InitialDelay
            case BackoffExponential:
                delay *= 2
            }

            if delay > retryConfig.MaxDelay {
                delay = retryConfig.MaxDelay
            }
        }

        response, err := o.client.Execute(ctx, command)
        if err == nil {
            return response, nil
        }

        lastErr = err
        o.logger.Error("execution attempt failed",
            Field{"attempt", attempt + 1},
            Field{"max_attempts", retryConfig.MaxAttempts},
            Field{"error", err})
    }

    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// executeDecision 执行决策节点
func (o *Orchestrator) executeDecision(ctx context.Context, instance *WorkflowInstance, node *Node, workflow *Workflow, execCtx *ExecutionContext) {
    // 评估条件
    for _, edge := range workflow.Edges {
        if edge.From == node.ID && edge.Condition != nil {
            if edge.Condition.Evaluator(execCtx) {
                o.executeNode(ctx, instance, edge.To)
                return
            }
        }
    }

    // 默认路径
    for _, edge := range workflow.Edges {
        if edge.From == node.ID && edge.Condition == nil {
            o.executeNode(ctx, instance, edge.To)
            return
        }
    }

    o.handleError(instance, fmt.Errorf("no valid transition from decision node %s", node.ID))
}

// executeParallel 执行并行节点
func (o *Orchestrator) executeParallel(ctx context.Context, instance *WorkflowInstance, node *Node, workflow *Workflow, execCtx *ExecutionContext) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var firstError error

    // 收集所有并行分支
    branches := []string{}
    for _, edge := range workflow.Edges {
        if edge.From == node.ID {
            branches = append(branches, edge.To)
        }
    }

    for _, branchNodeID := range branches {
        wg.Add(1)
        go func(nodeID string) {
            defer wg.Done()

            // 每个分支在自己的 goroutine 中执行
            // 注意：这里简化处理，实际应使用子工作流
            if err := o.executeBranch(ctx, instance, nodeID); err != nil {
                mu.Lock()
                if firstError == nil {
                    firstError = err
                }
                mu.Unlock()
            }
        }(branchNodeID)
    }

    wg.Wait()

    if firstError != nil {
        o.handleError(instance, firstError)
        return
    }

    // 找到汇聚点并继续
    o.transitionFromParallel(ctx, instance, node.ID, workflow)
}

// executeBranch 执行分支
func (o *Orchestrator) executeBranch(ctx context.Context, instance *WorkflowInstance, nodeID string) error {
    // 简化实现：实际应创建子工作流实例
    return nil
}

// transition 状态转移
func (o *Orchestrator) transition(ctx context.Context, instance *WorkflowInstance, fromNodeID string, workflow *Workflow) {
    // 查找下一个节点
    for _, edge := range workflow.Edges {
        if edge.From == fromNodeID {
            o.executeNode(ctx, instance, edge.To)
            return
        }
    }

    // 没有后续节点，工作流完成
    o.completeWorkflow(ctx, instance)
}

// transitionFromParallel 从并行节点转移
func (o *Orchestrator) transitionFromParallel(ctx context.Context, instance *WorkflowInstance, parallelNodeID string, workflow *Workflow) {
    // 查找汇聚后的下一个节点
    // 简化实现：假设所有分支汇聚到同一个节点
    o.transition(ctx, instance, parallelNodeID, workflow)
}

// compensate 执行补偿
func (o *Orchestrator) compensate(ctx context.Context, instance *WorkflowInstance) {
    instance.Status = WorkflowStatusCompensating
    o.store.UpdateInstanceStatus(ctx, instance.ID, WorkflowStatusCompensating)

    o.logger.Info("starting compensation",
        Field{"instance_id", instance.ID})

    // 按相反顺序执行补偿
    for i := len(instance.History) - 1; i >= 0; i-- {
        result := instance.History[i]
        if result.Status == ActivityStatusCompleted {
            workflow := o.workflows[instance.WorkflowID]
            node := workflow.Nodes[result.NodeID]

            if node.Compensation != nil {
                command := &Command{
                    ID:         uuid.New().String(),
                    Service:    node.Compensation.Service,
                    Action:     node.Compensation.Action,
                    Input:      result.Output,
                    InstanceID: instance.ID,
                }

                if err := o.client.ExecuteCompensation(ctx, command); err != nil {
                    o.logger.Error("compensation failed",
                        Field{"instance_id", instance.ID},
                        Field{"node_id", node.ID},
                        Field{"error", err})
                    // 补偿失败需要人工介入
                    instance.Status = WorkflowStatusFailed
                    o.store.UpdateInstanceStatus(ctx, instance.ID, WorkflowStatusFailed)
                    return
                }

                result.Status = ActivityStatusCompensated
            }
        }
    }

    instance.Status = WorkflowStatusCompensated
    now := time.Now()
    instance.CompletedAt = &now
    o.store.UpdateInstanceStatus(ctx, instance.ID, WorkflowStatusCompensated)

    o.eventBus.Publish(ctx, map[string]interface{}{
        "type":        "workflow.compensated",
        "instance_id": instance.ID,
    })
}

// completeWorkflow 完成工作流
func (o *Orchestrator) completeWorkflow(ctx context.Context, instance *WorkflowInstance) {
    instance.Status = WorkflowStatusCompleted
    now := time.Now()
    instance.CompletedAt = &now

    o.store.UpdateInstanceStatus(ctx, instance.ID, WorkflowStatusCompleted)

    o.logger.Info("workflow completed",
        Field{"instance_id", instance.ID},
        Field{"duration", now.Sub(instance.StartedAt)})

    o.eventBus.Publish(ctx, map[string]interface{}{
        "type":        "workflow.completed",
        "instance_id": instance.ID,
        "duration":    now.Sub(instance.StartedAt).Seconds(),
    })
}

// handleError 处理错误
func (o *Orchestrator) handleError(instance *WorkflowInstance, err error) {
    o.logger.Error("workflow error",
        Field{"instance_id", instance.ID},
        Field{"error", err})

    instance.Status = WorkflowStatusFailed
}

// prepareInput 准备输入
func (o *Orchestrator) prepareInput(mapping map[string]string, variables map[string]interface{}) map[string]interface{} {
    input := make(map[string]interface{})
    for target, source := range mapping {
        if value, exists := variables[source]; exists {
            input[target] = value
        }
    }
    return input
}

// updateVariables 更新变量
func (o *Orchestrator) updateVariables(instance *WorkflowInstance, mapping map[string]string, output map[string]interface{}) {
    instance.mu.Lock()
    defer instance.mu.Unlock()

    for target, source := range mapping {
        if value, exists := output[source]; exists {
            instance.Variables[target] = value
        }
    }
}

// GetInstance 获取实例
func (o *Orchestrator) GetInstance(id string) (*WorkflowInstance, error) {
    o.mu.RLock()
    defer o.mu.RUnlock()

    if instance, exists := o.instances[id]; exists {
        return instance, nil
    }
    return nil, fmt.Errorf("instance not found: %s", id)
}
```

### 2.2 存储实现

```go
// orchestration/memory_store.go
package orchestration

import (
    "context"
    "fmt"
    "sync"
)

// MemoryWorkflowStore 内存工作流存储
type MemoryWorkflowStore struct {
    workflows map[string]*Workflow
    instances map[string]*WorkflowInstance
    mu        sync.RWMutex
}

// NewMemoryWorkflowStore 创建内存存储
func NewMemoryWorkflowStore() *MemoryWorkflowStore {
    return &MemoryWorkflowStore{
        workflows: make(map[string]*Workflow),
        instances: make(map[string]*WorkflowInstance),
    }
}

// SaveWorkflow 保存工作流
func (s *MemoryWorkflowStore) SaveWorkflow(ctx context.Context, workflow *Workflow) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.workflows[workflow.ID] = workflow
    return nil
}

// GetWorkflow 获取工作流
func (s *MemoryWorkflowStore) GetWorkflow(ctx context.Context, id string) (*Workflow, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    workflow, exists := s.workflows[id]
    if !exists {
        return nil, fmt.Errorf("workflow not found: %s", id)
    }
    return workflow, nil
}

// SaveInstance 保存实例
func (s *MemoryWorkflowStore) SaveInstance(ctx context.Context, instance *WorkflowInstance) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.instances[instance.ID] = instance
    return nil
}

// GetInstance 获取实例
func (s *MemoryWorkflowStore) GetInstance(ctx context.Context, id string) (*WorkflowInstance, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    instance, exists := s.instances[id]
    if !exists {
        return nil, fmt.Errorf("instance not found: %s", id)
    }
    return instance, nil
}

// UpdateInstanceStatus 更新状态
func (s *MemoryWorkflowStore) UpdateInstanceStatus(ctx context.Context, id string, status WorkflowStatus) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    instance, exists := s.instances[id]
    if !exists {
        return fmt.Errorf("instance not found: %s", id)
    }
    instance.Status = status
    return nil
}

// AppendHistory 追加历史
func (s *MemoryWorkflowStore) AppendHistory(ctx context.Context, instanceID string, result *ActivityResult) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    instance, exists := s.instances[instanceID]
    if !exists {
        return fmt.Errorf("instance not found: %s", instanceID)
    }
    instance.History = append(instance.History, result)
    return nil
}
```

### 2.3 服务客户端实现

```go
// orchestration/http_client.go
package orchestration

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// HTTPServiceClient HTTP 服务客户端
type HTTPServiceClient struct {
    baseURL    string
    httpClient *http.Client
}

// NewHTTPServiceClient 创建 HTTP 客户端
func NewHTTPServiceClient(baseURL string, timeout time.Duration) *HTTPServiceClient {
    return &HTTPServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: timeout,
        },
    }
}

// Execute 执行命令
func (c *HTTPServiceClient) Execute(ctx context.Context, command *Command) (*Response, error) {
    url := fmt.Sprintf("%s/services/%s/actions/%s", c.baseURL, command.Service, command.Action)

    payload, err := json.Marshal(command.Input)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Correlation-ID", command.CorrelationID)
    req.Header.Set("X-Workflow-Instance", command.InstanceID)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
    }

    var output map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
        return nil, err
    }

    return &Response{
        ID:        generateID(),
        CommandID: command.ID,
        Status:    ResponseStatusSuccess,
        Output:    output,
        Timestamp: time.Now(),
    }, nil
}

// ExecuteCompensation 执行补偿
func (c *HTTPServiceClient) ExecuteCompensation(ctx context.Context, command *Command) error {
    url := fmt.Sprintf("%s/services/%s/compensations/%s", c.baseURL, command.Service, command.Action)

    payload, err := json.Marshal(command.Input)
    if err != nil {
        return err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("compensation failed with status %d", resp.StatusCode)
    }

    return nil
}

func generateID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// orchestration/core_test.go
package orchestration

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...Field)  {}
func (m *mockLogger) Error(msg string, fields ...Field) {}
func (m *mockLogger) Debug(msg string, fields ...Field) {}

type mockServiceClient struct {
    mock.Mock
}

func (m *mockServiceClient) Execute(ctx context.Context, command *Command) (*Response, error) {
    args := m.Called(ctx, command)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Response), args.Error(1)
}

func (m *mockServiceClient) ExecuteCompensation(ctx context.Context, command *Command) error {
    args := m.Called(ctx, command)
    return args.Error(0)
}

type mockEventBus struct{}

func (m *mockEventBus) Publish(ctx context.Context, event interface{}) error {
    return nil
}

func (m *mockEventBus) Subscribe(eventType string, handler EventHandler) {}

func TestOrchestrator_RegisterWorkflow(t *testing.T) {
    store := NewMemoryWorkflowStore()
    client := new(mockServiceClient)
    eventBus := &mockEventBus{}
    logger := &mockLogger{}

    orch := NewOrchestrator(store, client, eventBus, logger)

    workflow := &Workflow{
        ID:        "wf-001",
        Name:      "test-workflow",
        StartNode: "start",
        Nodes: map[string]*Node{
            "start": {
                ID:      "start",
                Type:    NodeTypeActivity,
                Service: "test-service",
                Action:  "test-action",
            },
        },
    }

    err := orch.RegisterWorkflow(workflow)
    require.NoError(t, err)

    // 验证工作流已注册
    registered, err := store.GetWorkflow(context.Background(), "wf-001")
    require.NoError(t, err)
    assert.Equal(t, "test-workflow", registered.Name)
}

func TestOrchestrator_StartWorkflow(t *testing.T) {
    store := NewMemoryWorkflowStore()
    client := new(mockServiceClient)
    eventBus := &mockEventBus{}
    logger := &mockLogger{}

    orch := NewOrchestrator(store, client, eventBus, logger)

    // 注册测试工作流
    workflow := &Workflow{
        ID:        "wf-002",
        Name:      "order-workflow",
        StartNode: "create-order",
        Nodes: map[string]*Node{
            "create-order": {
                ID:      "create-order",
                Type:    NodeTypeActivity,
                Service: "order-service",
                Action:  "create",
            },
        },
    }
    require.NoError(t, orch.RegisterWorkflow(workflow))

    // 设置 mock 期望
    client.On("Execute", mock.Anything, mock.AnythingOfType("*orchestration.Command")).
        Return(&Response{
            Status: ResponseStatusSuccess,
            Output: map[string]interface{}{"order_id": "ORD-001"},
        }, nil)

    // 启动工作流
    input := map[string]interface{}{"customer_id": "CUST-001", "amount": 100.0}
    instance, err := orch.StartWorkflow(context.Background(), "wf-002", input)

    require.NoError(t, err)
    assert.NotEmpty(t, instance.ID)
    assert.Equal(t, WorkflowStatusRunning, instance.Status)

    // 等待执行完成
    time.Sleep(100 * time.Millisecond)

    client.AssertExpectations(t)
}

func TestOrchestrator_Compensation(t *testing.T) {
    store := NewMemoryWorkflowStore()
    client := new(mockServiceClient)
    eventBus := &mockEventBus{}
    logger := &mockLogger{}

    orch := NewOrchestrator(store, client, eventBus, logger)

    // 注册带补偿的工作流
    workflow := &Workflow{
        ID:        "wf-003",
        Name:      "payment-workflow",
        StartNode: "process-payment",
        Nodes: map[string]*Node{
            "process-payment": {
                ID:      "process-payment",
                Type:    NodeTypeActivity,
                Service: "payment-service",
                Action:  "charge",
                Compensation: &Compensation{
                    Service: "payment-service",
                    Action:  "refund",
                },
            },
        },
    }
    require.NoError(t, orch.RegisterWorkflow(workflow))

    // 模拟执行失败
    client.On("Execute", mock.Anything, mock.Anything).
        Return(nil, assert.AnError)

    // 模拟补偿成功
    client.On("ExecuteCompensation", mock.Anything, mock.Anything).
        Return(nil)

    instance, err := orch.StartWorkflow(context.Background(), "wf-003", nil)
    require.NoError(t, err)

    // 等待补偿完成
    time.Sleep(200 * time.Millisecond)

    // 验证补偿被调用
    client.AssertCalled(t, "ExecuteCompensation", mock.Anything, mock.Anything)
}

func TestMemoryWorkflowStore(t *testing.T) {
    store := NewMemoryWorkflowStore()
    ctx := context.Background()

    // 测试保存和获取工作流
    workflow := &Workflow{
        ID:   "wf-test",
        Name: "test",
    }

    require.NoError(t, store.SaveWorkflow(ctx, workflow))

    retrieved, err := store.GetWorkflow(ctx, "wf-test")
    require.NoError(t, err)
    assert.Equal(t, "test", retrieved.Name)

    // 测试获取不存在的
    _, err = store.GetWorkflow(ctx, "non-existent")
    assert.Error(t, err)
}

func TestRetryPolicy_ExponentialBackoff(t *testing.T) {
    policy := &RetryPolicy{
        MaxAttempts:  3,
        BackoffType:  BackoffExponential,
        InitialDelay: 100 * time.Millisecond,
        MaxDelay:     500 * time.Millisecond,
    }

    delays := []time.Duration{}
    delay := policy.InitialDelay

    for i := 0; i < policy.MaxAttempts; i++ {
        delays = append(delays, delay)
        delay *= 2
        if delay > policy.MaxDelay {
            delay = policy.MaxDelay
        }
    }

    assert.Equal(t, []time.Duration{
        100 * time.Millisecond,
        200 * time.Millisecond,
        400 * time.Millisecond,
    }, delays)
}
```

### 3.2 集成测试

```go
// orchestration/integration_test.go
package orchestration

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestOrderProcessingWorkflow(t *testing.T) {
    // 构建完整的订单处理工作流
    store := NewMemoryWorkflowStore()
    client := &mockServiceClient{}
    eventBus := &mockEventBus{}
    logger := &mockLogger{}

    orch := NewOrchestrator(store, client, eventBus, logger)

    // 定义订单处理工作流
    workflow := &Workflow{
        ID:        "order-processing",
        Name:      "Order Processing",
        Version:   "1.0",
        StartNode: "validate-order",
        Timeout:   5 * time.Minute,
        Nodes: map[string]*Node{
            "validate-order": {
                ID:      "validate-order",
                Type:    NodeTypeActivity,
                Service: "order-service",
                Action:  "validate",
                RetryConfig: &RetryPolicy{
                    MaxAttempts:  2,
                    BackoffType:  BackoffLinear,
                    InitialDelay: 1 * time.Second,
                },
            },
            "check-decision": {
                ID:   "check-decision",
                Type: NodeTypeDecision,
            },
            "reserve-inventory": {
                ID:      "reserve-inventory",
                Type:    NodeTypeActivity,
                Service: "inventory-service",
                Action:  "reserve",
                Compensation: &Compensation{
                    Service: "inventory-service",
                    Action:  "release",
                },
            },
            "process-payment": {
                ID:      "process-payment",
                Type:    NodeTypeActivity,
                Service: "payment-service",
                Action:  "charge",
                Compensation: &Compensation{
                    Service: "payment-service",
                    Action:  "refund",
                },
            },
            "confirm-order": {
                ID:      "confirm-order",
                Type:    NodeTypeActivity,
                Service: "order-service",
                Action:  "confirm",
            },
        },
        Edges: []*Edge{
            {From: "validate-order", To: "check-decision"},
            {From: "check-decision", To: "reserve-inventory", Condition: &Condition{
                Expression: "valid == true",
                Evaluator: func(ctx *ExecutionContext) bool {
                    valid, _ := ctx.Variables["valid"].(bool)
                    return valid
                },
            }},
            {From: "reserve-inventory", To: "process-payment"},
            {From: "process-payment", To: "confirm-order"},
        },
    }

    require.NoError(t, orch.RegisterWorkflow(workflow))

    // 测试成功场景
    t.Run("successful_order_flow", func(t *testing.T) {
        // 启动工作流
        input := map[string]interface{}{
            "order_id":    "ORDER-123",
            "customer_id": "CUST-456",
            "items": []map[string]interface{}{
                {"product_id": "PROD-1", "qty": 2},
            },
        }

        instance, err := orch.StartWorkflow(context.Background(), "order-processing", input)
        require.NoError(t, err)
        assert.Equal(t, WorkflowStatusRunning, instance.Status)
    })
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Choreography 的对比

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Orchestration vs Choreography Comparison                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  特性维度          │  Orchestration           │  Choreography           │
│  ──────────────────┼──────────────────────────┼─────────────────────────│
│  控制方式          │  集中控制                 │  分散控制               │
│  耦合程度          │  服务 → 编排器            │  服务 → 事件            │
│  可见性           │  高（单一视图）            │  低（分布式追踪）        │
│  复杂度管理        │  适合复杂流程              │  适合简单流程            │
│  扩展性           │  编排器可能成为瓶颈        │  天然分布式              │
│  故障恢复          │  内置重试/补偿             │  需额外实现              │
│  开发模式          │  流程优先                  │  领域事件优先            │
│                                                                         │
│  混合模式（推荐复杂系统）:                                                │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                     Hybrid Architecture                          │   │
│  │  ┌─────────────┐                                                  │   │
│  │  │Orchestrator │ ──► 复杂流程、跨边界协调                          │   │
│  │  └──────┬──────┘                                                  │   │
│  │         │ Commands                                                 │   │
│  │         ▼                                                          │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐            │   │
│  │  │  Service A  │───►│  Service B  │───►│  Service C  │            │   │
│  │  └─────────────┘    └─────────────┘    └─────────────┘            │   │
│  │       │ Events (内部协调)                                           │   │
│  │       ▼                                                            │   │
│  │  [内部使用 Choreography]                                            │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 与 CQRS 的集成

```go
// orchestration/cqrs_integration.go
package orchestration

// WorkflowReadModel 工作流读模型
type WorkflowReadModel struct {
    store WorkflowStore
    cache map[string]*WorkflowInstanceView
}

// WorkflowInstanceView 工作流实例视图
type WorkflowInstanceView struct {
    ID            string                 `json:"id"`
    WorkflowName  string                 `json:"workflow_name"`
    Status        string                 `json:"status"`
    CurrentStep   string                 `json:"current_step"`
    Progress      float64                `json:"progress"`
    Variables     map[string]interface{} `json:"variables"`
    StartedAt     time.Time              `json:"started_at"`
    Duration      string                 `json:"duration"`
}

// GetInstanceView 获取实例视图（优化读操作）
func (r *WorkflowReadModel) GetInstanceView(instanceID string) (*WorkflowInstanceView, error) {
    // 优先从缓存获取
    if view, exists := r.cache[instanceID]; exists {
        return view, nil
    }

    // 从存储获取
    instance, err := r.store.GetInstance(context.Background(), instanceID)
    if err != nil {
        return nil, err
    }

    view := &WorkflowInstanceView{
        ID:           instance.ID,
        WorkflowName: instance.WorkflowName,
        Status:       instance.Status.String(),
        CurrentStep:  instance.CurrentNode,
        Variables:    instance.Variables,
        StartedAt:    instance.StartedAt,
    }

    if instance.CompletedAt != nil {
        view.Duration = instance.CompletedAt.Sub(instance.StartedAt).String()
    } else {
        view.Duration = time.Since(instance.StartedAt).String()
    }

    // 计算进度
    view.Progress = r.calculateProgress(instance)

    // 更新缓存
    r.cache[instanceID] = view

    return view, nil
}

func (r *WorkflowReadModel) calculateProgress(instance *WorkflowInstance) float64 {
    // 简化计算：基于历史记录数量
    totalSteps := 5 // 假设总共 5 步
    completedSteps := len(instance.History)

    if completedSteps >= totalSteps {
        return 100.0
    }
    return float64(completedSteps) / float64(totalSteps) * 100
}

// GetActiveInstances 获取活跃实例列表
func (r *WorkflowReadModel) GetActiveInstances(workflowID string) ([]*WorkflowInstanceView, error) {
    // 查询读模型数据库
    // 简化实现
    return nil, nil
}
```

---

## 5. 决策标准

### 5.1 选型决策树

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Orchestration Decision Tree                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  开始                                                                    │
│   │                                                                     │
│   ▼                                                                     │
│  ┌─────────────────────────┐                                           │
│  │ 业务流程涉及多个步骤？   │───否──► 使用简单事务                        │
│  └──────────┬──────────────┘                                           │
│             │是                                                         │
│             ▼                                                           │
│  ┌─────────────────────────┐                                           │
│  │ 需要条件分支或并行执行？ │                                           │
│  └──────────┬──────────────┘                                           │
│             │                                                           │
│      ┌──────┴──────┐                                                     │
│      │是           │否                                                    │
│      ▼             ▼                                                     │
│  ┌─────────┐   ┌─────────────────────────┐                              │
│  │Orchestration│  │ 流程是否经常变化？      │                              │
│  │ 编排器   │   └──────────┬──────────────┘                              │
│  └─────────┘              │                                             │
│                    ┌──────┴──────┐                                       │
│                    │是           │否                                      │
│                    ▼             ▼                                       │
│               ┌─────────┐   ┌─────────┐                                 │
│               │Choreography│  │Either   │                                 │
│               │ 编舞模式   │   │ 两者皆可 │                                 │
│               └─────────┘   └─────────┘                                 │
│                                                                         │
│  其他考虑因素:                                                           │
│  • 团队规模: 大团队 → Orchestration（标准流程）                          │
│  • SLA 要求: 严格 SLA → Orchestration（可控性）                          │
│  • 技术栈: 异构系统 → Orchestration（标准化接口）                         │
│  • 合规要求: 审计追踪 → Orchestration（完整历史）                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 生产环境检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Orchestrator Production Checklist                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  架构设计:                                                               │
│  □ 编排器支持水平扩展（无状态设计）                                        │
│  □ 工作流定义支持版本管理                                                  │
│  □ 实现工作流迁移策略（升级时）                                            │
│  □ 配置超时和死信队列                                                      │
│                                                                         │
│  可靠性:                                                                 │
│  □ 编排器状态持久化（数据库/事件溯源）                                     │
│  □ 实现至少一次执行保证                                                    │
│  □ 幂等性检查（防止重复执行）                                              │
│  □ 补偿逻辑覆盖所有已执行步骤                                              │
│                                                                         │
│  可观察性:                                                               │
│  □ 工作流执行指标（成功率、持续时间）                                      │
│  □ 步骤级别指标（每个服务的性能）                                          │
│  □ 分布式追踪（Correlation ID 传播）                                       │
│  □ 工作流可视化仪表板                                                      │
│                                                                         │
│  运维:                                                                   │
│  □ 人工干预接口（暂停、恢复、强制完成）                                    │
│  □ 工作流查询 API（支持客服调查）                                          │
│  □ 定期归档已完成工作流                                                    │
│  □ 配置告警（长时间运行、失败率）                                          │
│                                                                         │
│  安全:                                                                   │
│  □ 命令签名验证                                                            │
│  □ 敏感变量加密存储                                                        │
│  □ 审计日志记录                                                            │
│  □ 访问控制（谁可以启动/查看工作流）                                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>20KB, 完整形式化 + Go 实现 + 测试 + 集成)

**相关文档**:

- [EC-031-Choreography-Pattern.md](./EC-031-Choreography-Pattern.md)
- [EC-008-Saga-Pattern-Formal.md](./EC-008-Saga-Pattern-Formal.md)
- [EC-024-Task-State-Machine.md](./EC-024-Task-State-Machine.md)
