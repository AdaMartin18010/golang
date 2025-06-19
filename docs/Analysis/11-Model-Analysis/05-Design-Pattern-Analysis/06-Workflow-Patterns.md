# 工作流设计模式分析

## 目录

1. [概述](#1-概述)
2. [形式化定义](#2-形式化定义)
3. [基本控制流模式](#3-基本控制流模式)
4. [高级分支合并模式](#4-高级分支合并模式)
5. [结构模式](#5-结构模式)
6. [多实例模式](#6-多实例模式)
7. [状态基础模式](#7-状态基础模式)
8. [取消和触发模式](#8-取消和触发模式)
9. [工作流引擎](#9-工作流引擎)
10. [最佳实践](#10-最佳实践)
11. [案例分析](#11-案例分析)

## 1. 概述

### 1.1 工作流模式定义

工作流模式是描述业务流程控制逻辑的可重用解决方案。这些模式提供了构建复杂业务流程的标准方法，确保流程的可预测性、可维护性和可扩展性。

### 1.2 核心概念

- **活动(Activity)**: 工作流中的基本执行单元
- **网关(Gateway)**: 控制流程分支和合并的决策点
- **事件(Event)**: 触发或响应工作流状态变化
- **流程(Process)**: 由活动、网关和事件组成的完整业务流程

## 2. 形式化定义

### 2.1 工作流系统模型

**定义 2.1** (工作流系统): 一个工作流系统是一个八元组 $WS = (A, G, E, S, T, F, I, O)$，其中：

- $A = \{a_1, a_2, ..., a_n\}$ 是活动集合
- $G = \{g_1, g_2, ..., g_m\}$ 是网关集合
- $E = \{e_1, e_2, ..., e_k\}$ 是事件集合
- $S = \{s_1, s_2, ..., s_p\}$ 是状态集合
- $T: A \times S \rightarrow S$ 是状态转换函数
- $F: G \times S \rightarrow 2^A$ 是流程控制函数
- $I: E \times S \rightarrow S$ 是事件处理函数
- $O: A \times S \rightarrow S$ 是输出处理函数

### 2.2 工作流模式分类

**定义 2.2** (工作流模式): 工作流模式是一个三元组 $WP = (P, C, R)$，其中：

- $P$ 是模式类型集合
- $C$ 是控制流集合
- $R$ 是规则集合

## 3. 基本控制流模式

### 3.1 序列模式

**定义 3.1** (序列模式): 序列模式是一个二元组 $Seq = (A, <)$，其中：

- $A$ 是活动集合
- $<$ 是活动间的顺序关系

```go
// 序列模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// Activity 活动接口
type Activity interface {
    ID() string
    Execute(ctx context.Context, data interface{}) (interface{}, error)
}

// SequentialWorkflow 序列工作流
type SequentialWorkflow struct {
    activities []Activity
    mu         sync.RWMutex
}

// NewSequentialWorkflow 创建新的序列工作流
func NewSequentialWorkflow() *SequentialWorkflow {
    return &SequentialWorkflow{
        activities: make([]Activity, 0),
    }
}

// AddActivity 添加活动
func (sw *SequentialWorkflow) AddActivity(activity Activity) {
    sw.mu.Lock()
    defer sw.mu.Unlock()
    sw.activities = append(sw.activities, activity)
}

// Execute 执行工作流
func (sw *SequentialWorkflow) Execute(ctx context.Context, initialData interface{}) (interface{}, error) {
    sw.mu.RLock()
    defer sw.mu.RUnlock()
    
    data := initialData
    
    for _, activity := range sw.activities {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            result, err := activity.Execute(ctx, data)
            if err != nil {
                return nil, fmt.Errorf("activity %s failed: %w", activity.ID(), err)
            }
            data = result
        }
    }
    
    return data, nil
}

// BaseActivity 基础活动
type BaseActivity struct {
    id       string
    executor func(ctx context.Context, data interface{}) (interface{}, error)
}

// NewBaseActivity 创建新的基础活动
func NewBaseActivity(id string, executor func(ctx context.Context, data interface{}) (interface{}, error)) *BaseActivity {
    return &BaseActivity{
        id:       id,
        executor: executor,
    }
}

// ID 获取活动ID
func (ba *BaseActivity) ID() string {
    return ba.id
}

// Execute 执行活动
func (ba *BaseActivity) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    return ba.executor(ctx, data)
}
```

### 3.2 并行拆分模式

**定义 3.2** (并行拆分模式): 并行拆分模式是一个三元组 $PS = (A, B, \parallel)$，其中：

- $A$ 是输入活动
- $B$ 是并行分支集合
- $\parallel$ 是并行关系

```go
// 并行拆分模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// ParallelWorkflow 并行工作流
type ParallelWorkflow struct {
    branches []Activity
    mu       sync.RWMutex
}

// NewParallelWorkflow 创建新的并行工作流
func NewParallelWorkflow() *ParallelWorkflow {
    return &ParallelWorkflow{
        branches: make([]Activity, 0),
    }
}

// AddBranch 添加分支
func (pw *ParallelWorkflow) AddBranch(activity Activity) {
    pw.mu.Lock()
    defer pw.mu.Unlock()
    pw.branches = append(pw.branches, activity)
}

// Execute 执行并行工作流
func (pw *ParallelWorkflow) Execute(ctx context.Context, data interface{}) ([]interface{}, error) {
    pw.mu.RLock()
    defer pw.mu.RUnlock()
    
    if len(pw.branches) == 0 {
        return nil, fmt.Errorf("no branches to execute")
    }
    
    results := make([]interface{}, len(pw.branches))
    errors := make([]error, len(pw.branches))
    
    var wg sync.WaitGroup
    
    for i, branch := range pw.branches {
        wg.Add(1)
        go func(index int, activity Activity) {
            defer wg.Done()
            
            select {
            case <-ctx.Done():
                errors[index] = ctx.Err()
            default:
                result, err := activity.Execute(ctx, data)
                if err != nil {
                    errors[index] = err
                } else {
                    results[index] = result
                }
            }
        }(i, branch)
    }
    
    wg.Wait()
    
    // 检查是否有错误
    for _, err := range errors {
        if err != nil {
            return nil, fmt.Errorf("parallel execution failed: %w", err)
        }
    }
    
    return results, nil
}
```

### 3.3 同步模式

**定义 3.3** (同步模式): 同步模式是一个三元组 $Sync = (B, A, \land)$，其中：

- $B$ 是并行分支集合
- $A$ 是后续活动
- $\land$ 是同步合并关系

```go
// 同步模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// SynchronizationWorkflow 同步工作流
type SynchronizationWorkflow struct {
    branches []Activity
    next     Activity
    mu       sync.RWMutex
}

// NewSynchronizationWorkflow 创建新的同步工作流
func NewSynchronizationWorkflow(next Activity) *SynchronizationWorkflow {
    return &SynchronizationWorkflow{
        branches: make([]Activity, 0),
        next:     next,
    }
}

// AddBranch 添加分支
func (sw *SynchronizationWorkflow) AddBranch(activity Activity) {
    sw.mu.Lock()
    defer sw.mu.Unlock()
    sw.branches = append(sw.branches, activity)
}

// Execute 执行同步工作流
func (sw *SynchronizationWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    sw.mu.RLock()
    defer sw.mu.RUnlock()
    
    // 并行执行所有分支
    results := make([]interface{}, len(sw.branches))
    errors := make([]error, len(sw.branches))
    
    var wg sync.WaitGroup
    
    for i, branch := range sw.branches {
        wg.Add(1)
        go func(index int, activity Activity) {
            defer wg.Done()
            
            select {
            case <-ctx.Done():
                errors[index] = ctx.Err()
            default:
                result, err := activity.Execute(ctx, data)
                if err != nil {
                    errors[index] = err
                } else {
                    results[index] = result
                }
            }
        }(i, branch)
    }
    
    wg.Wait()
    
    // 检查是否有错误
    for _, err := range errors {
        if err != nil {
            return nil, fmt.Errorf("synchronization failed: %w", err)
        }
    }
    
    // 执行后续活动
    if sw.next != nil {
        return sw.next.Execute(ctx, results)
    }
    
    return results, nil
}
```

## 4. 高级分支合并模式

### 4.1 多选模式

**定义 4.1** (多选模式): 多选模式是一个四元组 $MC = (A, B, C, \phi)$，其中：

- $A$ 是输入活动
- $B$ 是可选分支集合
- $C$ 是选择条件集合
- $\phi: C \rightarrow 2^B$ 是选择函数

```go
// 多选模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// Condition 条件接口
type Condition interface {
    Evaluate(ctx context.Context, data interface{}) bool
}

// MultiChoiceWorkflow 多选工作流
type MultiChoiceWorkflow struct {
    branches map[Condition]Activity
    mu       sync.RWMutex
}

// NewMultiChoiceWorkflow 创建新的多选工作流
func NewMultiChoiceWorkflow() *MultiChoiceWorkflow {
    return &MultiChoiceWorkflow{
        branches: make(map[Condition]Activity),
    }
}

// AddBranch 添加分支
func (mcw *MultiChoiceWorkflow) AddBranch(condition Condition, activity Activity) {
    mcw.mu.Lock()
    defer mcw.mu.Unlock()
    mcw.branches[condition] = activity
}

// Execute 执行多选工作流
func (mcw *MultiChoiceWorkflow) Execute(ctx context.Context, data interface{}) ([]interface{}, error) {
    mcw.mu.RLock()
    defer mcw.mu.RUnlock()
    
    var selectedActivities []Activity
    
    // 评估所有条件
    for condition, activity := range mcw.branches {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            if condition.Evaluate(ctx, data) {
                selectedActivities = append(selectedActivities, activity)
            }
        }
    }
    
    if len(selectedActivities) == 0 {
        return nil, fmt.Errorf("no branches selected")
    }
    
    // 并行执行选中的分支
    results := make([]interface{}, len(selectedActivities))
    errors := make([]error, len(selectedActivities))
    
    var wg sync.WaitGroup
    
    for i, activity := range selectedActivities {
        wg.Add(1)
        go func(index int, act Activity) {
            defer wg.Done()
            
            select {
            case <-ctx.Done():
                errors[index] = ctx.Err()
            default:
                result, err := act.Execute(ctx, data)
                if err != nil {
                    errors[index] = err
                } else {
                    results[index] = result
                }
            }
        }(i, activity)
    }
    
    wg.Wait()
    
    // 检查是否有错误
    for _, err := range errors {
        if err != nil {
            return nil, fmt.Errorf("multi-choice execution failed: %w", err)
        }
    }
    
    return results, nil
}

// BaseCondition 基础条件
type BaseCondition struct {
    evaluator func(ctx context.Context, data interface{}) bool
}

// NewBaseCondition 创建新的基础条件
func NewBaseCondition(evaluator func(ctx context.Context, data interface{}) bool) *BaseCondition {
    return &BaseCondition{
        evaluator: evaluator,
    }
}

// Evaluate 评估条件
func (bc *BaseCondition) Evaluate(ctx context.Context, data interface{}) bool {
    return bc.evaluator(ctx, data)
}
```

### 4.2 结构化同步合并模式

**定义 4.2** (结构化同步合并): 结构化同步合并是一个四元组 $SSM = (B, A, C, \psi)$，其中：

- $B$ 是分支集合
- $A$ 是合并活动
- $C$ 是分支完成条件集合
- $\psi: C \rightarrow bool$ 是合并条件函数

```go
// 结构化同步合并模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// MergeCondition 合并条件接口
type MergeCondition interface {
    CanMerge(ctx context.Context, completedBranches []string) bool
}

// StructuredSyncMergeWorkflow 结构化同步合并工作流
type StructuredSyncMergeWorkflow struct {
    branches       []Activity
    mergeActivity  Activity
    mergeCondition MergeCondition
    mu             sync.RWMutex
}

// NewStructuredSyncMergeWorkflow 创建新的结构化同步合并工作流
func NewStructuredSyncMergeWorkflow(mergeActivity Activity, mergeCondition MergeCondition) *StructuredSyncMergeWorkflow {
    return &StructuredSyncMergeWorkflow{
        branches:       make([]Activity, 0),
        mergeActivity:  mergeActivity,
        mergeCondition: mergeCondition,
    }
}

// AddBranch 添加分支
func (ssmw *StructuredSyncMergeWorkflow) AddBranch(activity Activity) {
    ssmw.mu.Lock()
    defer ssmw.mu.Unlock()
    ssmw.branches = append(ssmw.branches, activity)
}

// Execute 执行结构化同步合并工作流
func (ssmw *StructuredSyncMergeWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    ssmw.mu.RLock()
    defer ssmw.mu.RUnlock()
    
    results := make([]interface{}, len(ssmw.branches))
    completedBranches := make([]string, 0)
    
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for i, branch := range ssmw.branches {
        wg.Add(1)
        go func(index int, activity Activity) {
            defer wg.Done()
            
            select {
            case <-ctx.Done():
                return
            default:
                result, err := activity.Execute(ctx, data)
                if err == nil {
                    mu.Lock()
                    results[index] = result
                    completedBranches = append(completedBranches, activity.ID())
                    mu.Unlock()
                }
            }
        }(i, branch)
    }
    
    wg.Wait()
    
    // 检查合并条件
    if ssmw.mergeCondition.CanMerge(ctx, completedBranches) {
        if ssmw.mergeActivity != nil {
            return ssmw.mergeActivity.Execute(ctx, results)
        }
        return results, nil
    }
    
    return nil, fmt.Errorf("merge condition not satisfied")
}

// BaseMergeCondition 基础合并条件
type BaseMergeCondition struct {
    evaluator func(ctx context.Context, completedBranches []string) bool
}

// NewBaseMergeCondition 创建新的基础合并条件
func NewBaseMergeCondition(evaluator func(ctx context.Context, completedBranches []string) bool) *BaseMergeCondition {
    return &BaseMergeCondition{
        evaluator: evaluator,
    }
}

// CanMerge 检查是否可以合并
func (bmc *BaseMergeCondition) CanMerge(ctx context.Context, completedBranches []string) bool {
    return bmc.evaluator(ctx, completedBranches)
}
```

## 5. 结构模式

### 5.1 任意循环模式

**定义 5.1** (任意循环模式): 任意循环模式是一个三元组 $AC = (A, C, \rho)$，其中：

- $A$ 是循环活动
- $C$ 是循环条件
- $\rho: C \times A \rightarrow bool$ 是循环控制函数

```go
// 任意循环模式实现
package workflow

import (
    "context"
    "fmt"
)

// LoopCondition 循环条件接口
type LoopCondition interface {
    ShouldContinue(ctx context.Context, data interface{}) bool
}

// ArbitraryCycleWorkflow 任意循环工作流
type ArbitraryCycleWorkflow struct {
    activity Activity
    condition LoopCondition
    maxIterations int
}

// NewArbitraryCycleWorkflow 创建新的任意循环工作流
func NewArbitraryCycleWorkflow(activity Activity, condition LoopCondition, maxIterations int) *ArbitraryCycleWorkflow {
    return &ArbitraryCycleWorkflow{
        activity:       activity,
        condition:      condition,
        maxIterations:  maxIterations,
    }
}

// Execute 执行任意循环工作流
func (acw *ArbitraryCycleWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    currentData := data
    iteration := 0
    
    for iteration < acw.maxIterations {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            if !acw.condition.ShouldContinue(ctx, currentData) {
                break
            }
            
            result, err := acw.activity.Execute(ctx, currentData)
            if err != nil {
                return nil, fmt.Errorf("iteration %d failed: %w", iteration, err)
            }
            
            currentData = result
            iteration++
        }
    }
    
    if iteration >= acw.maxIterations {
        return nil, fmt.Errorf("maximum iterations reached")
    }
    
    return currentData, nil
}

// BaseLoopCondition 基础循环条件
type BaseLoopCondition struct {
    evaluator func(ctx context.Context, data interface{}) bool
}

// NewBaseLoopCondition 创建新的基础循环条件
func NewBaseLoopCondition(evaluator func(ctx context.Context, data interface{}) bool) *BaseLoopCondition {
    return &BaseLoopCondition{
        evaluator: evaluator,
    }
}

// ShouldContinue 检查是否应该继续循环
func (blc *BaseLoopCondition) ShouldContinue(ctx context.Context, data interface{}) bool {
    return blc.evaluator(ctx, data)
}
```

## 6. 多实例模式

### 6.1 多实例不同步模式

**定义 6.1** (多实例不同步模式): 多实例不同步模式是一个四元组 $MI = (A, N, I, \delta)$，其中：

- $A$ 是活动模板
- $N$ 是实例数量
- $I$ 是实例集合
- $\delta: A \times N \rightarrow I$ 是实例化函数

```go
// 多实例不同步模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// MultipleInstancesWorkflow 多实例工作流
type MultipleInstancesWorkflow struct {
    activityTemplate Activity
    instanceCount    int
    mu               sync.RWMutex
}

// NewMultipleInstancesWorkflow 创建新的多实例工作流
func NewMultipleInstancesWorkflow(activityTemplate Activity, instanceCount int) *MultipleInstancesWorkflow {
    return &MultipleInstancesWorkflow{
        activityTemplate: activityTemplate,
        instanceCount:    instanceCount,
    }
}

// Execute 执行多实例工作流
func (miw *MultipleInstancesWorkflow) Execute(ctx context.Context, data interface{}) ([]interface{}, error) {
    miw.mu.RLock()
    defer miw.mu.RUnlock()
    
    results := make([]interface{}, miw.instanceCount)
    errors := make([]error, miw.instanceCount)
    
    var wg sync.WaitGroup
    
    for i := 0; i < miw.instanceCount; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            
            select {
            case <-ctx.Done():
                errors[index] = ctx.Err()
            default:
                // 为每个实例创建独立的数据副本
                instanceData := map[string]interface{}{
                    "original": data,
                    "instance": index,
                }
                
                result, err := miw.activityTemplate.Execute(ctx, instanceData)
                if err != nil {
                    errors[index] = err
                } else {
                    results[index] = result
                }
            }
        }(i)
    }
    
    wg.Wait()
    
    // 检查是否有错误
    for _, err := range errors {
        if err != nil {
            return nil, fmt.Errorf("multiple instances execution failed: %w", err)
        }
    }
    
    return results, nil
}
```

## 7. 状态基础模式

### 7.1 延迟选择模式

**定义 7.1** (延迟选择模式): 延迟选择模式是一个四元组 $DC = (A, E, T, \tau)$，其中：

- $A$ 是活动集合
- $E$ 是事件集合
- $T$ 是时间约束集合
- $\tau: E \times T \rightarrow A$ 是选择函数

```go
// 延迟选择模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Event 事件接口
type Event interface {
    ID() string
    Trigger(ctx context.Context) bool
}

// DeferredChoiceWorkflow 延迟选择工作流
type DeferredChoiceWorkflow struct {
    activities map[Event]Activity
    timeout    time.Duration
    mu         sync.RWMutex
}

// NewDeferredChoiceWorkflow 创建新的延迟选择工作流
func NewDeferredChoiceWorkflow(timeout time.Duration) *DeferredChoiceWorkflow {
    return &DeferredChoiceWorkflow{
        activities: make(map[Event]Activity),
        timeout:    timeout,
    }
}

// AddChoice 添加选择
func (dcw *DeferredChoiceWorkflow) AddChoice(event Event, activity Activity) {
    dcw.mu.Lock()
    defer dcw.mu.Unlock()
    dcw.activities[event] = activity
}

// Execute 执行延迟选择工作流
func (dcw *DeferredChoiceWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    dcw.mu.RLock()
    defer dcw.mu.RUnlock()
    
    // 创建带超时的上下文
    timeoutCtx, cancel := context.WithTimeout(ctx, dcw.timeout)
    defer cancel()
    
    // 监听所有事件
    eventChan := make(chan Event, len(dcw.activities))
    
    var wg sync.WaitGroup
    
    for event := range dcw.activities {
        wg.Add(1)
        go func(e Event) {
            defer wg.Done()
            
            for {
                select {
                case <-timeoutCtx.Done():
                    return
                default:
                    if e.Trigger(timeoutCtx) {
                        select {
                        case eventChan <- e:
                        case <-timeoutCtx.Done():
                        }
                        return
                    }
                    time.Sleep(100 * time.Millisecond)
                }
            }
        }(event)
    }
    
    // 等待第一个触发的事件
    select {
    case triggeredEvent := <-eventChan:
        if activity, exists := dcw.activities[triggeredEvent]; exists {
            return activity.Execute(ctx, data)
        }
    case <-timeoutCtx.Done():
        return nil, fmt.Errorf("deferred choice timeout")
    }
    
    return nil, fmt.Errorf("no event triggered")
}

// BaseEvent 基础事件
type BaseEvent struct {
    id       string
    trigger  func(ctx context.Context) bool
}

// NewBaseEvent 创建新的基础事件
func NewBaseEvent(id string, trigger func(ctx context.Context) bool) *BaseEvent {
    return &BaseEvent{
        id:      id,
        trigger: trigger,
    }
}

// ID 获取事件ID
func (be *BaseEvent) ID() string {
    return be.id
}

// Trigger 触发事件
func (be *BaseEvent) Trigger(ctx context.Context) bool {
    return be.trigger(ctx)
}
```

## 8. 取消和触发模式

### 8.1 取消活动模式

**定义 8.1** (取消活动模式): 取消活动模式是一个三元组 $CA = (A, C, \gamma)$，其中：

- $A$ 是活动集合
- $C$ 是取消条件集合
- $\gamma: C \times A \rightarrow bool$ 是取消函数

```go
// 取消活动模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// CancelCondition 取消条件接口
type CancelCondition interface {
    ShouldCancel(ctx context.Context, data interface{}) bool
}

// CancellableActivity 可取消活动
type CancellableActivity struct {
    activity Activity
    condition CancelCondition
    mu        sync.RWMutex
}

// NewCancellableActivity 创建新的可取消活动
func NewCancellableActivity(activity Activity, condition CancelCondition) *CancellableActivity {
    return &CancellableActivity{
        activity:  activity,
        condition: condition,
    }
}

// Execute 执行可取消活动
func (ca *CancellableActivity) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    ca.mu.Lock()
    defer ca.mu.Unlock()
    
    // 检查是否应该取消
    if ca.condition.ShouldCancel(ctx, data) {
        return nil, fmt.Errorf("activity cancelled")
    }
    
    // 执行活动
    return ca.activity.Execute(ctx, data)
}

// ID 获取活动ID
func (ca *CancellableActivity) ID() string {
    return ca.activity.ID()
}

// BaseCancelCondition 基础取消条件
type BaseCancelCondition struct {
    evaluator func(ctx context.Context, data interface{}) bool
}

// NewBaseCancelCondition 创建新的基础取消条件
func NewBaseCancelCondition(evaluator func(ctx context.Context, data interface{}) bool) *BaseCancelCondition {
    return &BaseCancelCondition{
        evaluator: evaluator,
    }
}

// ShouldCancel 检查是否应该取消
func (bcc *BaseCancelCondition) ShouldCancel(ctx context.Context, data interface{}) bool {
    return bcc.evaluator(ctx, data)
}
```

## 9. 工作流引擎

### 9.1 工作流引擎架构

```go
// 工作流引擎实现
package workflow

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// WorkflowEngine 工作流引擎
type WorkflowEngine struct {
    workflows map[string]Workflow
    executor  *WorkflowExecutor
    storage   WorkflowStorage
    mu        sync.RWMutex
}

// Workflow 工作流接口
type Workflow interface {
    ID() string
    Execute(ctx context.Context, data interface{}) (interface{}, error)
}

// WorkflowExecutor 工作流执行器
type WorkflowExecutor struct {
    workers int
    queue   chan WorkflowTask
    mu      sync.RWMutex
}

// WorkflowTask 工作流任务
type WorkflowTask struct {
    WorkflowID string
    Data       interface{}
    Context    context.Context
    Result     chan WorkflowResult
}

// WorkflowResult 工作流结果
type WorkflowResult struct {
    Data  interface{}
    Error error
}

// WorkflowStorage 工作流存储接口
type WorkflowStorage interface {
    Save(ctx context.Context, workflowID string, data interface{}) error
    Load(ctx context.Context, workflowID string) (interface{}, error)
}

// NewWorkflowEngine 创建新的工作流引擎
func NewWorkflowEngine(storage WorkflowStorage, workers int) *WorkflowEngine {
    engine := &WorkflowEngine{
        workflows: make(map[string]Workflow),
        executor:  NewWorkflowExecutor(workers),
        storage:   storage,
    }
    
    return engine
}

// RegisterWorkflow 注册工作流
func (we *WorkflowEngine) RegisterWorkflow(workflow Workflow) {
    we.mu.Lock()
    defer we.mu.Unlock()
    we.workflows[workflow.ID()] = workflow
}

// ExecuteWorkflow 执行工作流
func (we *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID string, data interface{}) (interface{}, error) {
    we.mu.RLock()
    workflow, exists := we.workflows[workflowID]
    we.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("workflow %s not found", workflowID)
    }
    
    // 保存初始数据
    if err := we.storage.Save(ctx, workflowID, data); err != nil {
        return nil, fmt.Errorf("failed to save workflow data: %w", err)
    }
    
    // 执行工作流
    result, err := workflow.Execute(ctx, data)
    
    // 保存结果
    if err := we.storage.Save(ctx, workflowID+"_result", result); err != nil {
        return nil, fmt.Errorf("failed to save workflow result: %w", err)
    }
    
    return result, err
}

// NewWorkflowExecutor 创建新的工作流执行器
func NewWorkflowExecutor(workers int) *WorkflowExecutor {
    executor := &WorkflowExecutor{
        workers: workers,
        queue:   make(chan WorkflowTask, 100),
    }
    
    // 启动工作协程
    for i := 0; i < workers; i++ {
        go executor.worker()
    }
    
    return executor
}

// worker 工作协程
func (we *WorkflowExecutor) worker() {
    for task := range we.queue {
        // 执行工作流任务
        // 这里可以实现具体的执行逻辑
        select {
        case task.Result <- WorkflowResult{Data: nil, Error: nil}:
        case <-task.Context.Done():
        }
    }
}
```

## 10. 最佳实践

### 10.1 设计原则

1. **单一职责**: 每个活动只负责一个特定的任务
2. **可组合性**: 工作流应该能够组合成更复杂的流程
3. **可测试性**: 工作流应该易于单元测试和集成测试
4. **可观测性**: 提供完善的日志、监控和追踪
5. **容错性**: 设计时考虑各种故障场景

### 10.2 实现建议

1. **使用状态机**: 复杂工作流可以使用状态机模式
2. **实现重试机制**: 处理临时性故障
3. **使用事件驱动**: 提高系统的响应性和可扩展性
4. **实现版本控制**: 支持工作流版本管理
5. **提供可视化**: 工作流设计和监控的可视化界面

## 11. 案例分析

### 11.1 订单处理工作流

```go
// 订单处理工作流示例
package workflow

import (
    "context"
    "fmt"
)

// Order 订单结构
type Order struct {
    ID       string  `json:"id"`
    Amount   float64 `json:"amount"`
    Status   string  `json:"status"`
    Customer string  `json:"customer"`
}

// OrderProcessingWorkflow 订单处理工作流
type OrderProcessingWorkflow struct {
    sequential *SequentialWorkflow
}

// NewOrderProcessingWorkflow 创建新的订单处理工作流
func NewOrderProcessingWorkflow() *OrderProcessingWorkflow {
    sequential := NewSequentialWorkflow()
    
    // 添加订单验证活动
    sequential.AddActivity(NewBaseActivity("validate_order", func(ctx context.Context, data interface{}) (interface{}, error) {
        order := data.(Order)
        if order.Amount <= 0 {
            return nil, fmt.Errorf("invalid order amount")
        }
        order.Status = "validated"
        return order, nil
    }))
    
    // 添加库存检查活动
    sequential.AddActivity(NewBaseActivity("check_inventory", func(ctx context.Context, data interface{}) (interface{}, error) {
        order := data.(Order)
        // 模拟库存检查
        order.Status = "inventory_checked"
        return order, nil
    }))
    
    // 添加支付处理活动
    sequential.AddActivity(NewBaseActivity("process_payment", func(ctx context.Context, data interface{}) (interface{}, error) {
        order := data.(Order)
        // 模拟支付处理
        order.Status = "paid"
        return order, nil
    }))
    
    // 添加发货活动
    sequential.AddActivity(NewBaseActivity("ship_order", func(ctx context.Context, data interface{}) (interface{}, error) {
        order := data.(Order)
        // 模拟发货
        order.Status = "shipped"
        return order, nil
    }))
    
    return &OrderProcessingWorkflow{
        sequential: sequential,
    }
}

// Execute 执行订单处理工作流
func (opw *OrderProcessingWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    return opw.sequential.Execute(ctx, data)
}

// ID 获取工作流ID
func (opw *OrderProcessingWorkflow) ID() string {
    return "order_processing"
}
```

### 11.2 审批工作流

```go
// 审批工作流示例
package workflow

import (
    "context"
    "fmt"
)

// ApprovalRequest 审批请求
type ApprovalRequest struct {
    ID          string  `json:"id"`
    Amount      float64 `json:"amount"`
    Requester   string  `json:"requester"`
    Status      string  `json:"status"`
    Approvers   []string `json:"approvers"`
    ApprovedBy  []string `json:"approved_by"`
}

// ApprovalWorkflow 审批工作流
type ApprovalWorkflow struct {
    multiChoice *MultiChoiceWorkflow
}

// NewApprovalWorkflow 创建新的审批工作流
func NewApprovalWorkflow() *ApprovalWorkflow {
    multiChoice := NewMultiChoiceWorkflow()
    
    // 添加小额审批分支
    multiChoice.AddBranch(
        NewBaseCondition(func(ctx context.Context, data interface{}) bool {
            request := data.(ApprovalRequest)
            return request.Amount <= 1000
        }),
        NewBaseActivity("small_amount_approval", func(ctx context.Context, data interface{}) (interface{}, error) {
            request := data.(ApprovalRequest)
            request.Status = "approved"
            request.ApprovedBy = append(request.ApprovedBy, "manager")
            return request, nil
        }),
    )
    
    // 添加大额审批分支
    multiChoice.AddBranch(
        NewBaseCondition(func(ctx context.Context, data interface{}) bool {
            request := data.(ApprovalRequest)
            return request.Amount > 1000
        }),
        NewBaseActivity("large_amount_approval", func(ctx context.Context, data interface{}) (interface{}, error) {
            request := data.(ApprovalRequest)
            request.Status = "pending_director"
            return request, nil
        }),
    )
    
    return &ApprovalWorkflow{
        multiChoice: multiChoice,
    }
}

// Execute 执行审批工作流
func (aw *ApprovalWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    results, err := aw.multiChoice.Execute(ctx, data)
    if err != nil {
        return nil, err
    }
    
    // 返回第一个结果
    if len(results) > 0 {
        return results[0], nil
    }
    
    return nil, fmt.Errorf("no approval branch executed")
}

// ID 获取工作流ID
func (aw *ApprovalWorkflow) ID() string {
    return "approval_workflow"
}
```

---

**总结**: 本文档提供了工作流设计模式的完整分析，包括形式化定义、Golang实现和最佳实践。这些模式为构建复杂业务流程提供了重要的理论基础和实践指导，支持各种业务场景的需求。
