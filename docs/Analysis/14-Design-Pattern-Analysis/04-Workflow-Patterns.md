# 14.1 工作流模式分析

<!-- TOC START -->
- [14.1 工作流模式分析](#工作流模式分析)
  - [14.1.1 概述](#概述)
  - [14.1.2 1. 状态机模式 (State Machine)](#1-状态机模式-state-machine)
    - [14.1.2.1 定义](#定义)
    - [14.1.2.2 形式化定义](#形式化定义)
    - [14.1.2.3 Golang实现](#golang实现)
  - [14.1.3 2. 工作流引擎模式 (Workflow Engine)](#2-工作流引擎模式-workflow-engine)
    - [14.1.3.1 定义](#定义)
    - [14.1.3.2 Golang实现](#golang实现)
  - [14.1.4 3. 任务队列模式 (Task Queue)](#3-任务队列模式-task-queue)
    - [14.1.4.1 定义](#定义)
    - [14.1.4.2 Golang实现](#golang实现)
  - [14.1.5 4. 编排vs协同模式 (Orchestration vs Choreography)](#4-编排vs协同模式-orchestration-vs-choreography)
    - [14.1.5.1 编排模式 (Orchestration)](#编排模式-orchestration)
    - [14.1.5.2 协同模式 (Choreography)](#协同模式-choreography)
    - [14.1.5.3 Golang实现](#golang实现)
  - [14.1.6 5. 性能分析](#5-性能分析)
    - [14.1.6.1 模式性能对比](#模式性能对比)
    - [14.1.6.2 性能指标](#性能指标)
    - [14.1.6.3 容量规划](#容量规划)
  - [14.1.7 6. 最佳实践](#6-最佳实践)
    - [14.1.7.1 设计原则](#设计原则)
    - [14.1.7.2 实现建议](#实现建议)
    - [14.1.7.3 常见陷阱](#常见陷阱)
  - [14.1.8 7. 应用场景](#7-应用场景)
    - [14.1.8.1 状态机模式](#状态机模式)
    - [14.1.8.2 工作流引擎](#工作流引擎)
    - [14.1.8.3 任务队列](#任务队列)
    - [14.1.8.4 编排模式](#编排模式)
    - [14.1.8.5 协同模式](#协同模式)
  - [14.1.9 8. 总结](#8-总结)
    - [14.1.9.1 关键优势](#关键优势)
    - [14.1.9.2 成功要素](#成功要素)
<!-- TOC END -->

## 14.1.1 概述

工作流模式是处理业务流程自动化的核心设计模式。本文档基于Golang技术栈，深入分析各种工作流模式的设计、实现和性能特征。

## 14.1.2 1. 状态机模式 (State Machine)

### 14.1.2.1 定义

管理对象的状态转换，确保状态转换的合法性和一致性。

### 14.1.2.2 形式化定义

$$\text{StateMachine} = (S, E, T, I, F)$$

其中：

- $S$ 是状态集合
- $E$ 是事件集合
- $T$ 是转换函数
- $I$ 是初始状态
- $F$ 是最终状态集合

### 14.1.2.3 Golang实现

```go
package statemachine

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// State 状态接口
type State interface {
    Name() string
    Enter(ctx context.Context) error
    Exit(ctx context.Context) error
    Handle(ctx context.Context, event Event) (State, error)
}

// Event 事件接口
type Event interface {
    Name() string
    Data() interface{}
}

// SimpleEvent 简单事件
type SimpleEvent struct {
    eventName string
    eventData interface{}
}

func (e *SimpleEvent) Name() string {
    return e.eventName
}

func (e *SimpleEvent) Data() interface{} {
    return e.eventData
}

// Transition 状态转换
type Transition struct {
    FromState string
    Event     string
    ToState   string
    Condition func(ctx context.Context, event Event) bool
    Action    func(ctx context.Context, event Event) error
}

// StateMachine 状态机
type StateMachine struct {
    name        string
    currentState State
    states      map[string]State
    transitions []*Transition
    history     []*StateChange
    mu          sync.RWMutex
    ctx         context.Context
    cancel      context.CancelFunc
    listeners   []StateChangeListener
}

// StateChange 状态变化
type StateChange struct {
    FromState  string
    ToState    string
    Event      string
    Timestamp  time.Time
    Data       interface{}
}

// StateChangeListener 状态变化监听器
type StateChangeListener interface {
    OnStateChange(change *StateChange)
}

// NewStateMachine 创建状态机
func NewStateMachine(name string) *StateMachine {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &StateMachine{
        name:        name,
        states:      make(map[string]State),
        transitions: make([]*Transition, 0),
        history:     make([]*StateChange, 0),
        ctx:         ctx,
        cancel:      cancel,
        listeners:   make([]StateChangeListener, 0),
    }
}

// AddState 添加状态
func (sm *StateMachine) AddState(state State) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.states[state.Name()] = state
}

// AddTransition 添加转换
func (sm *StateMachine) AddTransition(transition *Transition) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.transitions = append(sm.transitions, transition)
}

// SetInitialState 设置初始状态
func (sm *StateMachine) SetInitialState(stateName string) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    state, exists := sm.states[stateName]
    if !exists {
        return fmt.Errorf("state not found: %s", stateName)
    }
    
    sm.currentState = state
    return state.Enter(sm.ctx)
}

// Trigger 触发事件
func (sm *StateMachine) Trigger(event Event) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    if sm.currentState == nil {
        return fmt.Errorf("no initial state set")
    }
    
    // 查找转换
    var transition *Transition
    for _, t := range sm.transitions {
        if t.FromState == sm.currentState.Name() && t.Event == event.Name() {
            if t.Condition == nil || t.Condition(sm.ctx, event) {
                transition = t
                break
            }
        }
    }
    
    if transition == nil {
        return fmt.Errorf("no valid transition for event %s from state %s", 
            event.Name(), sm.currentState.Name())
    }
    
    // 执行转换
    toState, exists := sm.states[transition.ToState]
    if !exists {
        return fmt.Errorf("target state not found: %s", transition.ToState)
    }
    
    // 执行动作
    if transition.Action != nil {
        if err := transition.Action(sm.ctx, event); err != nil {
            return fmt.Errorf("transition action failed: %v", err)
        }
    }
    
    // 退出当前状态
    if err := sm.currentState.Exit(sm.ctx); err != nil {
        return fmt.Errorf("exit state failed: %v", err)
    }
    
    // 记录状态变化
    change := &StateChange{
        FromState: sm.currentState.Name(),
        ToState:   toState.Name(),
        Event:     event.Name(),
        Timestamp: time.Now(),
        Data:      event.Data(),
    }
    sm.history = append(sm.history, change)
    
    // 通知监听器
    for _, listener := range sm.listeners {
        listener.OnStateChange(change)
    }
    
    // 进入新状态
    if err := toState.Enter(sm.ctx); err != nil {
        return fmt.Errorf("enter state failed: %v", err)
    }
    
    sm.currentState = toState
    return nil
}

// GetCurrentState 获取当前状态
func (sm *StateMachine) GetCurrentState() State {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.currentState
}

// GetHistory 获取历史记录
func (sm *StateMachine) GetHistory() []*StateChange {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    history := make([]*StateChange, len(sm.history))
    copy(history, sm.history)
    return history
}

// AddListener 添加监听器
func (sm *StateMachine) AddListener(listener StateChangeListener) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.listeners = append(sm.listeners, listener)
}

// Shutdown 关闭状态机
func (sm *StateMachine) Shutdown() {
    sm.cancel()
}

// 具体状态实现
type OrderState struct {
    name string
}

func (s *OrderState) Name() string {
    return s.name
}

func (s *OrderState) Enter(ctx context.Context) error {
    fmt.Printf("Entering state: %s\n", s.name)
    return nil
}

func (s *OrderState) Exit(ctx context.Context) error {
    fmt.Printf("Exiting state: %s\n", s.name)
    return nil
}

func (s *OrderState) Handle(ctx context.Context, event Event) (State, error) {
    // 默认处理，实际应该根据具体状态实现
    return nil, fmt.Errorf("not implemented")
}

```

## 14.1.3 2. 工作流引擎模式 (Workflow Engine)

### 14.1.3.1 定义

管理和执行复杂的工作流程，支持条件分支、并行执行、错误处理等。

### 14.1.3.2 Golang实现

```go
package workflowengine

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Task 任务接口
type Task interface {
    ID() string
    Name() string
    Execute(ctx context.Context) error
    GetDependencies() []string
    GetRetryCount() int
    GetTimeout() time.Duration
}

// SimpleTask 简单任务
type SimpleTask struct {
    taskID       string
    taskName     string
    dependencies []string
    retryCount   int
    timeout      time.Duration
    executor     func(ctx context.Context) error
}

func (t *SimpleTask) ID() string {
    return t.taskID
}

func (t *SimpleTask) Name() string {
    return t.taskName
}

func (t *SimpleTask) Execute(ctx context.Context) error {
    if t.executor != nil {
        return t.executor(ctx)
    }
    return nil
}

func (t *SimpleTask) GetDependencies() []string {
    return t.dependencies
}

func (t *SimpleTask) GetRetryCount() int {
    return t.retryCount
}

func (t *SimpleTask) GetTimeout() time.Duration {
    return t.timeout
}

// TaskStatus 任务状态
type TaskStatus int

const (
    TaskPending TaskStatus = iota
    TaskRunning
    TaskCompleted
    TaskFailed
    TaskSkipped
)

// TaskResult 任务结果
type TaskResult struct {
    TaskID   string
    Status   TaskStatus
    Error    error
    StartTime time.Time
    EndTime   time.Time
    Retries   int
}

// Workflow 工作流
type Workflow struct {
    ID          string
    Name        string
    Tasks       map[string]Task
    Results     map[string]*TaskResult
    Status      WorkflowStatus
    mu          sync.RWMutex
    ctx         context.Context
    cancel      context.CancelFunc
    listeners   []WorkflowListener
}

// WorkflowStatus 工作流状态
type WorkflowStatus int

const (
    WorkflowPending WorkflowStatus = iota
    WorkflowRunning
    WorkflowCompleted
    WorkflowFailed
    WorkflowCancelled
)

// WorkflowListener 工作流监听器
type WorkflowListener interface {
    OnTaskCompleted(taskID string, result *TaskResult)
    OnTaskFailed(taskID string, result *TaskResult)
    OnWorkflowCompleted(workflowID string)
    OnWorkflowFailed(workflowID string, error error)
}

// NewWorkflow 创建工作流
func NewWorkflow(id, name string) *Workflow {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &Workflow{
        ID:        id,
        Name:      name,
        Tasks:     make(map[string]Task),
        Results:   make(map[string]*TaskResult),
        Status:    WorkflowPending,
        ctx:       ctx,
        cancel:    cancel,
        listeners: make([]WorkflowListener, 0),
    }
}

// AddTask 添加任务
func (w *Workflow) AddTask(task Task) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.Tasks[task.ID()] = task
}

// Execute 执行工作流
func (w *Workflow) Execute() error {
    w.mu.Lock()
    if w.Status != WorkflowPending {
        w.mu.Unlock()
        return fmt.Errorf("workflow is not in pending status")
    }
    w.Status = WorkflowRunning
    w.mu.Unlock()
    
    // 构建依赖图
    graph := w.buildDependencyGraph()
    
    // 执行任务
    return w.executeTasks(graph)
}

// buildDependencyGraph 构建依赖图
func (w *Workflow) buildDependencyGraph() map[string][]string {
    graph := make(map[string][]string)
    
    for taskID, task := range w.Tasks {
        dependencies := task.GetDependencies()
        graph[taskID] = dependencies
    }
    
    return graph
}

// executeTasks 执行任务
func (w *Workflow) executeTasks(graph map[string][]string) error {
    completed := make(map[string]bool)
    running := make(map[string]bool)
    
    for {
        // 查找可执行的任务
        executable := w.findExecutableTasks(graph, completed, running)
        
        if len(executable) == 0 {
            if len(completed) == len(w.Tasks) {
                // 所有任务完成
                w.mu.Lock()
                w.Status = WorkflowCompleted
                w.mu.Unlock()
                
                for _, listener := range w.listeners {
                    listener.OnWorkflowCompleted(w.ID)
                }
                return nil
            } else if len(running) == 0 {
                // 死锁检测
                w.mu.Lock()
                w.Status = WorkflowFailed
                w.mu.Unlock()
                
                for _, listener := range w.listeners {
                    listener.OnWorkflowFailed(w.ID, fmt.Errorf("deadlock detected"))
                }
                return fmt.Errorf("deadlock detected")
            }
        }
        
        // 并行执行任务
        var wg sync.WaitGroup
        for _, taskID := range executable {
            wg.Add(1)
            go func(tid string) {
                defer wg.Done()
                w.executeTask(tid, graph, completed, running)
            }(taskID)
        }
        wg.Wait()
    }
}

// findExecutableTasks 查找可执行任务
func (w *Workflow) findExecutableTasks(graph map[string][]string, completed, running map[string]bool) []string {
    var executable []string
    
    for taskID, dependencies := range graph {
        if completed[taskID] || running[taskID] {
            continue
        }
        
        // 检查依赖是否完成
        allDepsCompleted := true
        for _, dep := range dependencies {
            if !completed[dep] {
                allDepsCompleted = false
                break
            }
        }
        
        if allDepsCompleted {
            executable = append(executable, taskID)
        }
    }
    
    return executable
}

// executeTask 执行单个任务
func (w *Workflow) executeTask(taskID string, graph map[string][]string, completed, running map[string]bool) {
    w.mu.Lock()
    running[taskID] = true
    w.mu.Unlock()
    
    task := w.Tasks[taskID]
    result := &TaskResult{
        TaskID:   taskID,
        Status:   TaskRunning,
        StartTime: time.Now(),
    }
    
    // 执行任务
    var err error
    for retry := 0; retry <= task.GetRetryCount(); retry++ {
        ctx, cancel := context.WithTimeout(w.ctx, task.GetTimeout())
        
        err = task.Execute(ctx)
        cancel()
        
        if err == nil {
            break
        }
        
        result.Retries = retry
        if retry < task.GetRetryCount() {
            time.Sleep(time.Duration(retry+1) * time.Second) // 指数退避
        }
    }
    
    result.EndTime = time.Now()
    
    if err != nil {
        result.Status = TaskFailed
        result.Error = err
        
        w.mu.Lock()
        w.Results[taskID] = result
        w.Status = WorkflowFailed
        w.mu.Unlock()
        
        for _, listener := range w.listeners {
            listener.OnTaskFailed(taskID, result)
        }
    } else {
        result.Status = TaskCompleted
        
        w.mu.Lock()
        w.Results[taskID] = result
        completed[taskID] = true
        w.mu.Unlock()
        
        for _, listener := range w.listeners {
            listener.OnTaskCompleted(taskID, result)
        }
    }
    
    w.mu.Lock()
    delete(running, taskID)
    w.mu.Unlock()
}

// GetStatus 获取状态
func (w *Workflow) GetStatus() WorkflowStatus {
    w.mu.RLock()
    defer w.mu.RUnlock()
    return w.Status
}

// GetResults 获取结果
func (w *Workflow) GetResults() map[string]*TaskResult {
    w.mu.RLock()
    defer w.mu.RUnlock()
    
    results := make(map[string]*TaskResult)
    for k, v := range w.Results {
        results[k] = v
    }
    return results
}

// AddListener 添加监听器
func (w *Workflow) AddListener(listener WorkflowListener) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.listeners = append(w.listeners, listener)
}

// Cancel 取消工作流
func (w *Workflow) Cancel() {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    if w.Status == WorkflowRunning {
        w.Status = WorkflowCancelled
        w.cancel()
    }
}

// WorkflowEngine 工作流引擎
type WorkflowEngine struct {
    workflows map[string]*Workflow
    mu        sync.RWMutex
    ctx       context.Context
    cancel    context.CancelFunc
}

// NewWorkflowEngine 创建工作流引擎
func NewWorkflowEngine() *WorkflowEngine {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &WorkflowEngine{
        workflows: make(map[string]*Workflow),
        ctx:       ctx,
        cancel:    cancel,
    }
}

// RegisterWorkflow 注册工作流
func (e *WorkflowEngine) RegisterWorkflow(workflow *Workflow) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.workflows[workflow.ID] = workflow
}

// ExecuteWorkflow 执行工作流
func (e *WorkflowEngine) ExecuteWorkflow(workflowID string) error {
    e.mu.RLock()
    workflow, exists := e.workflows[workflowID]
    e.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("workflow not found: %s", workflowID)
    }
    
    return workflow.Execute()
}

// GetWorkflow 获取工作流
func (e *WorkflowEngine) GetWorkflow(workflowID string) (*Workflow, bool) {
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    workflow, exists := e.workflows[workflowID]
    return workflow, exists
}

// Shutdown 关闭引擎
func (e *WorkflowEngine) Shutdown() {
    e.cancel()
}

```

## 14.1.4 3. 任务队列模式 (Task Queue)

### 14.1.4.1 定义

管理异步任务的执行，支持任务调度、优先级、重试等功能。

### 14.1.4.2 Golang实现

```go
package taskqueue

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Task 任务接口
type Task interface {
    ID() string
    Priority() int
    Execute(ctx context.Context) error
    GetRetryCount() int
    GetTimeout() time.Duration
    GetDelay() time.Duration
}

// SimpleTask 简单任务
type SimpleTask struct {
    taskID     string
    priority   int
    retryCount int
    timeout    time.Duration
    delay      time.Duration
    executor   func(ctx context.Context) error
}

func (t *SimpleTask) ID() string {
    return t.taskID
}

func (t *SimpleTask) Priority() int {
    return t.priority
}

func (t *SimpleTask) Execute(ctx context.Context) error {
    if t.executor != nil {
        return t.executor(ctx)
    }
    return nil
}

func (t *SimpleTask) GetRetryCount() int {
    return t.retryCount
}

func (t *SimpleTask) GetTimeout() time.Duration {
    return t.timeout
}

func (t *SimpleTask) GetDelay() time.Duration {
    return t.delay
}

// TaskStatus 任务状态
type TaskStatus int

const (
    TaskQueued TaskStatus = iota
    TaskRunning
    TaskCompleted
    TaskFailed
    TaskCancelled
)

// TaskResult 任务结果
type TaskResult struct {
    TaskID    string
    Status    TaskStatus
    Error     error
    StartTime time.Time
    EndTime   time.Time
    Retries   int
    Priority  int
}

// PriorityQueue 优先级队列
type PriorityQueue struct {
    tasks []*TaskItem
    mu    sync.RWMutex
}

// TaskItem 任务项
type TaskItem struct {
    task     Task
    priority int
    index    int
}

func NewPriorityQueue() *PriorityQueue {
    return &PriorityQueue{
        tasks: make([]*TaskItem, 0),
    }
}

func (pq *PriorityQueue) Push(task Task) {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    
    item := &TaskItem{
        task:     task,
        priority: task.Priority(),
    }
    
    pq.tasks = append(pq.tasks, item)
    pq.up(len(pq.tasks) - 1)
}

func (pq *PriorityQueue) Pop() Task {
    pq.mu.Lock()
    defer pq.mu.Unlock()
    
    if len(pq.tasks) == 0 {
        return nil
    }
    
    task := pq.tasks[0].task
    pq.tasks[0] = pq.tasks[len(pq.tasks)-1]
    pq.tasks = pq.tasks[:len(pq.tasks)-1]
    
    if len(pq.tasks) > 0 {
        pq.down(0)
    }
    
    return task
}

func (pq *PriorityQueue) Peek() Task {
    pq.mu.RLock()
    defer pq.mu.RUnlock()
    
    if len(pq.tasks) == 0 {
        return nil
    }
    
    return pq.tasks[0].task
}

func (pq *PriorityQueue) Len() int {
    pq.mu.RLock()
    defer pq.mu.RUnlock()
    return len(pq.tasks)
}

func (pq *PriorityQueue) up(index int) {
    for {
        parent := (index - 1) / 2
        if parent == index || pq.tasks[parent].priority <= pq.tasks[index].priority {
            break
        }
        pq.tasks[parent], pq.tasks[index] = pq.tasks[index], pq.tasks[parent]
        index = parent
    }
}

func (pq *PriorityQueue) down(index int) {
    for {
        child := 2*index + 1
        if child >= len(pq.tasks) {
            break
        }
        
        if child+1 < len(pq.tasks) && pq.tasks[child+1].priority < pq.tasks[child].priority {
            child++
        }
        
        if pq.tasks[index].priority <= pq.tasks[child].priority {
            break
        }
        
        pq.tasks[index], pq.tasks[child] = pq.tasks[child], pq.tasks[index]
        index = child
    }
}

// TaskQueue 任务队列
type TaskQueue struct {
    name           string
    queue          *PriorityQueue
    workers        int
    results        map[string]*TaskResult
    mu             sync.RWMutex
    ctx            context.Context
    cancel         context.CancelFunc
    wg             sync.WaitGroup
    taskChan       chan Task
    resultChan     chan *TaskResult
    listeners      []TaskListener
}

// TaskListener 任务监听器
type TaskListener interface {
    OnTaskQueued(taskID string)
    OnTaskStarted(taskID string)
    OnTaskCompleted(taskID string, result *TaskResult)
    OnTaskFailed(taskID string, result *TaskResult)
}

// NewTaskQueue 创建任务队列
func NewTaskQueue(name string, workers int) *TaskQueue {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &TaskQueue{
        name:       name,
        queue:      NewPriorityQueue(),
        workers:    workers,
        results:    make(map[string]*TaskResult),
        ctx:        ctx,
        cancel:     cancel,
        taskChan:   make(chan Task, workers*2),
        resultChan: make(chan *TaskResult, workers*2),
        listeners:  make([]TaskListener, 0),
    }
}

// Start 启动队列
func (tq *TaskQueue) Start() {
    // 启动工作器
    for i := 0; i < tq.workers; i++ {
        tq.wg.Add(1)
        go tq.worker(i)
    }
    
    // 启动调度器
    tq.wg.Add(1)
    go tq.scheduler()
    
    // 启动结果处理器
    tq.wg.Add(1)
    go tq.resultProcessor()
}

// worker 工作器
func (tq *TaskQueue) worker(id int) {
    defer tq.wg.Done()
    
    for {
        select {
        case task := <-tq.taskChan:
            if task == nil {
                return
            }
            
            tq.executeTask(task)
        case <-tq.ctx.Done():
            return
        }
    }
}

// scheduler 调度器
func (tq *TaskQueue) scheduler() {
    defer tq.wg.Done()
    
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // 检查队列中的任务
            for tq.queue.Len() > 0 {
                task := tq.queue.Peek()
                if task == nil {
                    break
                }
                
                // 检查延迟
                if task.GetDelay() > 0 {
                    // 这里应该实现延迟逻辑
                    break
                }
                
                // 发送任务到工作器
                select {
                case tq.taskChan <- task:
                    tq.queue.Pop()
                default:
                    // 工作器忙，等待
                    break
                }
            }
        case <-tq.ctx.Done():
            return
        }
    }
}

// executeTask 执行任务
func (tq *TaskQueue) executeTask(task Task) {
    result := &TaskResult{
        TaskID:    task.ID(),
        Status:    TaskRunning,
        StartTime: time.Now(),
        Priority:  task.Priority(),
    }
    
    // 通知任务开始
    for _, listener := range tq.listeners {
        listener.OnTaskStarted(task.ID())
    }
    
    // 执行任务
    var err error
    for retry := 0; retry <= task.GetRetryCount(); retry++ {
        ctx, cancel := context.WithTimeout(tq.ctx, task.GetTimeout())
        
        err = task.Execute(ctx)
        cancel()
        
        if err == nil {
            break
        }
        
        result.Retries = retry
        if retry < task.GetRetryCount() {
            time.Sleep(time.Duration(retry+1) * time.Second)
        }
    }
    
    result.EndTime = time.Now()
    
    if err != nil {
        result.Status = TaskFailed
        result.Error = err
        
        for _, listener := range tq.listeners {
            listener.OnTaskFailed(task.ID(), result)
        }
    } else {
        result.Status = TaskCompleted
        
        for _, listener := range tq.listeners {
            listener.OnTaskCompleted(task.ID(), result)
        }
    }
    
    // 发送结果
    select {
    case tq.resultChan <- result:
    default:
        // 结果通道已满，丢弃
    }
}

// resultProcessor 结果处理器
func (tq *TaskQueue) resultProcessor() {
    defer tq.wg.Done()
    
    for {
        select {
        case result := <-tq.resultChan:
            if result == nil {
                return
            }
            
            tq.mu.Lock()
            tq.results[result.TaskID] = result
            tq.mu.Unlock()
        case <-tq.ctx.Done():
            return
        }
    }
}

// Enqueue 入队任务
func (tq *TaskQueue) Enqueue(task Task) {
    tq.queue.Push(task)
    
    for _, listener := range tq.listeners {
        listener.OnTaskQueued(task.ID())
    }
}

// GetResult 获取结果
func (tq *TaskQueue) GetResult(taskID string) (*TaskResult, bool) {
    tq.mu.RLock()
    defer tq.mu.RUnlock()
    
    result, exists := tq.results[taskID]
    return result, exists
}

// GetStats 获取统计信息
func (tq *TaskQueue) GetStats() map[string]interface{} {
    tq.mu.RLock()
    defer tq.mu.RUnlock()
    
    total := len(tq.results)
    completed := 0
    failed := 0
    
    for _, result := range tq.results {
        switch result.Status {
        case TaskCompleted:
            completed++
        case TaskFailed:
            failed++
        }
    }
    
    return map[string]interface{}{
        "total":     total,
        "completed": completed,
        "failed":    failed,
        "queued":    tq.queue.Len(),
        "workers":   tq.workers,
    }
}

// AddListener 添加监听器
func (tq *TaskQueue) AddListener(listener TaskListener) {
    tq.mu.Lock()
    defer tq.mu.Unlock()
    tq.listeners = append(tq.listeners, listener)
}

// Shutdown 关闭队列
func (tq *TaskQueue) Shutdown() {
    tq.cancel()
    close(tq.taskChan)
    close(tq.resultChan)
    tq.wg.Wait()
}

```

## 14.1.5 4. 编排vs协同模式 (Orchestration vs Choreography)

### 14.1.5.1 编排模式 (Orchestration)

中央协调器控制整个工作流程的执行。

### 14.1.5.2 协同模式 (Choreography)

各个服务通过事件进行协作，没有中央协调器。

### 14.1.5.3 Golang实现

```go
package orchestration

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Orchestrator 编排器
type Orchestrator struct {
    name     string
    steps    []*OrchestrationStep
    results  map[string]interface{}
    mu       sync.RWMutex
    ctx      context.Context
    cancel   context.CancelFunc
}

// OrchestrationStep 编排步骤
type OrchestrationStep struct {
    ID       string
    Name     string
    Service  string
    Execute  func(ctx context.Context, data map[string]interface{}) (interface{}, error)
    Rollback func(ctx context.Context, data map[string]interface{}) error
}

// NewOrchestrator 创建编排器
func NewOrchestrator(name string) *Orchestrator {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &Orchestrator{
        name:    name,
        steps:   make([]*OrchestrationStep, 0),
        results: make(map[string]interface{}),
        ctx:     ctx,
        cancel:  cancel,
    }
}

// AddStep 添加步骤
func (o *Orchestrator) AddStep(step *OrchestrationStep) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.steps = append(o.steps, step)
}

// Execute 执行编排
func (o *Orchestrator) Execute() error {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    for i, step := range o.steps {
        // 执行步骤
        result, err := step.Execute(o.ctx, o.results)
        if err != nil {
            // 回滚前面的步骤
            return o.rollback(i - 1)
        }
        
        o.results[step.ID] = result
    }
    
    return nil
}

// rollback 回滚
func (o *Orchestrator) rollback(fromIndex int) error {
    for i := fromIndex; i >= 0; i-- {
        step := o.steps[i]
        if step.Rollback != nil {
            if err := step.Rollback(o.ctx, o.results); err != nil {
                return fmt.Errorf("rollback failed for step %s: %v", step.ID, err)
            }
        }
    }
    return nil
}

// GetResults 获取结果
func (o *Orchestrator) GetResults() map[string]interface{} {
    o.mu.RLock()
    defer o.mu.RUnlock()
    
    results := make(map[string]interface{})
    for k, v := range o.results {
        results[k] = v
    }
    return results
}

// EventBus 事件总线
type EventBus struct {
    subscribers map[string][]EventSubscriber
    mu          sync.RWMutex
    ctx         context.Context
    cancel      context.CancelFunc
}

// Event 事件
type Event struct {
    Type    string
    Data    interface{}
    Source  string
    Time    time.Time
}

// EventSubscriber 事件订阅者
type EventSubscriber interface {
    OnEvent(event *Event) error
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &EventBus{
        subscribers: make(map[string][]EventSubscriber),
        ctx:         ctx,
        cancel:      cancel,
    }
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType string, subscriber EventSubscriber) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    
    eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
}

// Publish 发布事件
func (eb *EventBus) Publish(event *Event) error {
    eb.mu.RLock()
    subscribers := eb.subscribers[event.Type]
    eb.mu.RUnlock()
    
    for _, subscriber := range subscribers {
        if err := subscriber.OnEvent(event); err != nil {
            return err
        }
    }
    
    return nil
}

// Service 服务
type Service struct {
    name    string
    eventBus *EventBus
}

// NewService 创建服务
func NewService(name string, eventBus *EventBus) *Service {
    return &Service{
        name:     name,
        eventBus: eventBus,
    }
}

// ProcessEvent 处理事件
func (s *Service) ProcessEvent(event *Event) error {
    // 处理事件逻辑
    fmt.Printf("Service %s processing event: %s\n", s.name, event.Type)
    
    // 发布新事件
    newEvent := &Event{
        Type:   "processed",
        Data:   fmt.Sprintf("Processed by %s", s.name),
        Source: s.name,
        Time:   time.Now(),
    }
    
    return s.eventBus.Publish(newEvent)
}

```

## 14.1.6 5. 性能分析

### 14.1.6.1 模式性能对比

| 模式 | 延迟 | 吞吐量 | 可扩展性 | 复杂度 | 一致性 |
|------|------|--------|----------|--------|--------|
| 状态机 | 低 | 高 | 中 | 低 | 强 |
| 工作流引擎 | 中 | 中 | 高 | 高 | 强 |
| 任务队列 | 低 | 高 | 高 | 中 | 最终 |
| 编排 | 高 | 中 | 中 | 高 | 强 |
| 协同 | 中 | 高 | 高 | 中 | 最终 |

### 14.1.6.2 性能指标

**状态机性能**:
$$\text{StateTransitionTime} = \text{ValidationTime} + \text{ActionTime}$$

**工作流性能**:
$$\text{WorkflowTime} = \sum_{i=1}^{n} \text{TaskTime}_i + \text{CoordinationTime}$$

**任务队列性能**:
$$\text{QueueLatency} = \text{EnqueueTime} + \text{WaitTime} + \text{ProcessingTime}$$

### 14.1.6.3 容量规划

**状态机容量**:
$$C_{statemachine} = \frac{\text{StateCount} \times \text{EventRate}}{\text{ProcessingCapacity}}$$

**工作流容量**:
$$C_{workflow} = \frac{\text{TaskCount} \times \text{Parallelism}}{\text{ResourceCapacity}}$$

## 14.1.7 6. 最佳实践

### 14.1.7.1 设计原则

1. **简单优先**: 优先使用简单的模式
2. **可扩展性**: 考虑系统的扩展性
3. **容错设计**: 处理异常情况
4. **监控告警**: 全面的监控体系

### 14.1.7.2 实现建议

1. **异步处理**: 使用异步模式提高性能
2. **错误处理**: 完善的错误处理机制
3. **重试策略**: 实现智能重试
4. **超时控制**: 设置合理的超时时间

### 14.1.7.3 常见陷阱

1. **状态不一致**: 处理状态同步问题
2. **死锁**: 避免循环依赖
3. **性能瓶颈**: 避免热点资源
4. **复杂性**: 避免过度设计

## 14.1.8 7. 应用场景

### 14.1.8.1 状态机模式

- 订单状态管理
- 用户状态跟踪
- 设备状态控制
- 游戏状态管理

### 14.1.8.2 工作流引擎

- 业务流程自动化
- 审批流程
- 数据处理管道
- CI/CD流程

### 14.1.8.3 任务队列

- 异步任务处理
- 批量数据处理
- 邮件发送
- 文件处理

### 14.1.8.4 编排模式

- 复杂业务流程
- 事务处理
- 服务编排
- 数据同步

### 14.1.8.5 协同模式

- 事件驱动架构
- 微服务通信
- 实时数据处理
- 消息传递系统

## 14.1.9 8. 总结

工作流模式为业务流程自动化提供了重要的设计指导。通过合理应用这些模式，可以构建出高效、可靠的工作流系统。

### 14.1.9.1 关键优势

- **自动化**: 减少人工干预
- **可靠性**: 提高系统稳定性
- **可扩展**: 支持业务增长
- **可维护**: 清晰的流程设计

### 14.1.9.2 成功要素

1. **合理选择**: 根据需求选择合适的模式
2. **性能优化**: 持续的性能优化
3. **监控告警**: 完善的监控体系
4. **测试验证**: 全面的测试覆盖

通过合理应用工作流模式，可以构建出高质量的工作流系统，为业务发展提供强有力的技术支撑。
