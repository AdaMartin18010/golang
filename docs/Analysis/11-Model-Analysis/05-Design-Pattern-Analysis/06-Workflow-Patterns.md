# 11.5.1 工作流模式 (Workflow Patterns)

<!-- TOC START -->
- [11.5.1 工作流模式 (Workflow Patterns)](#工作流模式-workflow-patterns)
  - [11.5.1.1 目录](#目录)
  - [11.5.1.2 1. 概述](#1-概述)
    - [11.5.1.2.1 工作流模式的定义与意义](#工作流模式的定义与意义)
    - [11.5.1.2.2 工作流模式的分类体系](#工作流模式的分类体系)
    - [11.5.1.2.3 在Go语言中实现工作流模式的优势](#在go语言中实现工作流模式的优势)
  - [11.5.1.3 2. 工作流模式形式化定义](#2-工作流模式形式化定义)
    - [11.5.1.3.1 工作流系统的数学表示](#工作流系统的数学表示)
    - [11.5.1.3.2 工作流模式的形式化属性](#工作流模式的形式化属性)
    - [11.5.1.3.3 Petri网与工作流模式](#petri网与工作流模式)
  - [11.5.1.4 3. 基本控制流模式](#3-基本控制流模式)
    - [11.5.1.4.1 序列模式 (Sequence Pattern)](#序列模式-sequence-pattern)
      - [11.5.1.4.1.1 形式化定义](#形式化定义)
      - [11.5.1.4.1.2 Golang实现](#golang实现)
    - [11.5.1.4.2 并行拆分模式 (Parallel Split Pattern)](#并行拆分模式-parallel-split-pattern)
      - [11.5.1.4.2.1 形式化定义](#形式化定义)
      - [11.5.1.4.2.2 Golang实现](#golang实现)
    - [11.5.1.4.3 同步模式 (Synchronization Pattern)](#同步模式-synchronization-pattern)
      - [11.5.1.4.3.1 形式化定义](#形式化定义)
      - [11.5.1.4.3.2 Golang实现](#golang实现)
  - [11.5.1.5 4. 高级分支合并模式](#4-高级分支合并模式)
    - [11.5.1.5.1 多选模式](#多选模式)
    - [11.5.1.5.2 结构化同步合并模式](#结构化同步合并模式)
  - [11.5.1.6 5. 结构模式](#5-结构模式)
    - [11.5.1.6.1 任意循环模式](#任意循环模式)
  - [11.5.1.7 6. 多实例模式](#6-多实例模式)
    - [11.5.1.7.1 多实例不同步模式](#多实例不同步模式)
  - [11.5.1.8 7. 状态基础模式](#7-状态基础模式)
    - [11.5.1.8.1 延迟选择模式](#延迟选择模式)
    - [11.5.1.8.2 里程碑模式](#里程碑模式)
  - [11.5.1.9 8. 资源模式](#8-资源模式)
    - [11.5.1.9.1 直接分配模式](#直接分配模式)
    - [11.5.1.9.2 基于角色分配模式](#基于角色分配模式)
    - [11.5.1.9.3 授权模式](#授权模式)
  - [11.5.1.10 9. 数据流模式](#9-数据流模式)
    - [11.5.1.10.1 数据传递模式](#数据传递模式)
    - [11.5.1.10.2 数据转换模式](#数据转换模式)
    - [11.5.1.10.3 数据路由模式](#数据路由模式)
  - [11.5.1.11 10. Petri网建模分析](#10-petri网建模分析)
    - [11.5.1.11.1 形式化定义](#形式化定义)
  - [11.5.1.12 11. 工作流引擎实现](#11-工作流引擎实现)
    - [11.5.1.12.1 工作流引擎架构](#工作流引擎架构)
  - [11.5.1.13 12. 模式集成案例](#12-模式集成案例)
    - [11.5.1.13.1 订单处理工作流](#订单处理工作流)
    - [11.5.1.13.2 审批工作流](#审批工作流)
<!-- TOC END -->

## 11.5.1.1 目录

- [工作流模式 (Workflow Patterns)](#工作流模式-workflow-patterns)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 工作流模式的定义与意义](#11-工作流模式的定义与意义)
    - [1.2 工作流模式的分类体系](#12-工作流模式的分类体系)
    - [1.3 在Go语言中实现工作流模式的优势](#13-在go语言中实现工作流模式的优势)
  - [2. 工作流模式形式化定义](#2-工作流模式形式化定义)
    - [2.1 工作流系统的数学表示](#21-工作流系统的数学表示)
    - [2.2 工作流模式的形式化属性](#22-工作流模式的形式化属性)
    - [2.3 Petri网与工作流模式](#23-petri网与工作流模式)
  - [3. 基本控制流模式](#3-基本控制流模式)
    - [3.1 序列模式 (Sequence Pattern)](#31-序列模式-sequence-pattern)
      - [3.1.1 形式化定义](#311-形式化定义)
      - [3.1.2 Golang实现](#312-golang实现)
    - [3.2 并行拆分模式 (Parallel Split Pattern)](#32-并行拆分模式-parallel-split-pattern)
      - [3.2.1 形式化定义](#321-形式化定义)
      - [3.2.2 Golang实现](#322-golang实现)
    - [3.3 同步模式 (Synchronization Pattern)](#33-同步模式-synchronization-pattern)
      - [3.3.1 形式化定义](#331-形式化定义)
      - [3.3.2 Golang实现](#332-golang实现)
  - [4. 高级分支合并模式](#4-高级分支合并模式)
    - [4.1 多选模式](#41-多选模式)
    - [4.2 结构化同步合并模式](#42-结构化同步合并模式)
  - [5. 结构模式](#5-结构模式)
    - [5.1 任意循环模式](#51-任意循环模式)
  - [6. 多实例模式](#6-多实例模式)
    - [6.1 多实例不同步模式](#61-多实例不同步模式)
  - [7. 状态基础模式](#7-状态基础模式)
    - [7.1 延迟选择模式](#71-延迟选择模式)
    - [7.2 里程碑模式](#72-里程碑模式)
  - [8. 资源模式](#8-资源模式)
    - [8.1 直接分配模式](#81-直接分配模式)
    - [8.2 基于角色分配模式](#82-基于角色分配模式)
    - [8.3 授权模式](#83-授权模式)
  - [9. 数据流模式](#9-数据流模式)
    - [9.1 数据传递模式](#91-数据传递模式)
    - [9.2 数据转换模式](#92-数据转换模式)
    - [9.3 数据路由模式](#93-数据路由模式)
  - [10. Petri网建模分析](#10-petri网建模分析)
    - [10.1 形式化定义](#101-形式化定义)
  - [11. 工作流引擎实现](#11-工作流引擎实现)
    - [11.1 工作流引擎架构](#111-工作流引擎架构)
  - [12. 模式集成案例](#12-模式集成案例)
    - [12.1 订单处理工作流](#121-订单处理工作流)
    - [12.2 审批工作流](#122-审批工作流)

## 11.5.1.2 1. 概述

### 11.5.1.2.1 工作流模式的定义与意义

工作流模式是一组在工作流系统中经常出现的控制流、数据流和资源流的可重用解决方案。这些模式为设计和实现复杂工作流系统提供了形式化的方法和最佳实践。

### 11.5.1.2.2 工作流模式的分类体系

工作流模式主要分为以下几类：

1. **基本控制流模式**: 序列、并行拆分、同步、独占选择、简单合并等
2. **高级分支合并模式**: 多选、同步合并、多合并、结构化判别等
3. **结构模式**: 任意循环、隐式终止等
4. **多实例模式**: 无同步多实例、设计时多实例、运行时多实例等
5. **状态基础模式**: 延迟选择、交错并行路由、里程碑等
6. **取消模式**: 取消活动、取消案例等

### 11.5.1.2.3 在Go语言中实现工作流模式的优势

Go语言具有以下特性，使其特别适合实现工作流模式：

- **并发模型**: goroutine和channel为并行工作流提供了简洁的表达方式
- **类型系统**: 接口和结构体支持灵活定义工作流组件
- **错误处理**: 显式的错误处理机制有助于工作流中的异常管理
- **标准库**: context包为工作流的跟踪和取消提供了支持
- **编译速度**: 快速的编译周期方便工作流的迭代开发

## 11.5.1.3 2. 工作流模式形式化定义

### 11.5.1.3.1 工作流系统的数学表示

在工作流系统中，我们可以使用形式化的数学定义来描述工作流结构和行为。

**定义 2.1** (工作流): 工作流是一个七元组 $W = (A, C, D, F, R, S, T)$，其中：

- $A = \{a_1, a_2, ..., a_n\}$ 是活动(Activity)集合
- $C = \{c_1, c_2, ..., c_m\}$ 是条件(Condition)集合
- $D = \{d_1, d_2, ..., d_k\}$ 是数据对象(Data)集合
- $F \subseteq (A \times C) \cup (C \times A)$ 是控制流(Flow)关系
- $R = \{r_1, r_2, ..., r_j\}$ 是资源(Resource)集合
- $S: A \rightarrow 2^{R}$ 是资源分配函数
- $T: A \rightarrow 2^{D \times D}$ 是数据转换函数

### 11.5.1.3.2 工作流模式的形式化属性

工作流模式具有以下关键属性，可以通过形式化方法进行分析：

1. **可靠性(Reliability)**: 给定条件下，工作流能否正确完成
2. **可终止性(Termination)**: 工作流是否总能到达终止状态
3. **安全性(Safety)**: 工作流是否不会进入不良状态
4. **活性(Liveness)**: 工作流是否总能继续执行
5. **公平性(Fairness)**: 所有可能的分支是否都有执行的机会

### 11.5.1.3.3 Petri网与工作流模式

Petri网是一种强大的工作流形式化工具，可用于验证工作流的性质和行为。

**定义 2.2** (Petri网): Petri网是一个五元组 $PN = (P, T, F, W, M_0)$，其中：

- $P$ 是库所(Place)集合，表示状态或条件
- $T$ 是变迁(Transition)集合，表示活动或事件
- $F \subseteq (P \times T) \cup (T \times P)$ 是弧(Arc)集合
- $W: F \rightarrow \mathbb{N}^+$ 是权重函数
- $M_0: P \rightarrow \mathbb{N}$ 是初始标记(Initial marking)

工作流模式可以映射到相应的Petri网结构，从而利用Petri网理论进行形式化分析。

## 11.5.1.4 3. 基本控制流模式

### 11.5.1.4.1 序列模式 (Sequence Pattern)

#### 11.5.1.4.1.1 形式化定义

**定义 3.1** (序列模式): 给定两个活动 $a_1$ 和 $a_2$，序列模式定义了一个偏序关系 $a_1 < a_2$，表示 $a_1$ 必须完全完成后 $a_2$ 才能开始执行。在Petri网中，这可表示为：

1. 活动 $a_1$ 对应变迁 $t_1$
2. 活动 $a_2$ 对应变迁 $t_2$
3. 存在一个库所 $p$ 使得 $(t_1, p) \in F$ 且 $(p, t_2) \in F$

**定理 3.1**: 序列模式满足以下属性：

- 可靠性：如果 $a_1$ 和 $a_2$ 都能终止，则序列可靠
- 可终止性：如果 $a_1$ 和 $a_2$ 都是可终止的，则序列可终止
- 无死锁：如果 $a_1$ 能完成，序列模式不会引入新的死锁

#### 11.5.1.4.1.2 Golang实现

```go
package workflow

// Task 表示工作流中的任务
type Task[I any, O any] interface {
 Execute(input I) (O, error)
 Name() string
}

// SequentialTask 表示序列模式中的任务
type SequentialTask[I any, M any, O any] struct {
 name      string
 firstTask Task[I, M]
 lastTask  Task[M, O]
}

// 创建新的序列任务
func NewSequentialTask[I any, M any, O any](
 name string,
 firstTask Task[I, M],
 lastTask Task[M, O],
) *SequentialTask[I, M, O] {
 return &SequentialTask[I, M, O]{
  name:      name,
  firstTask: firstTask,
  lastTask:  lastTask,
 }
}

// Execute 执行序列任务
func (t *SequentialTask[I, M, O]) Execute(input I) (O, error) {
 var zero O
 
 // 执行第一个任务
 intermediate, err := t.firstTask.Execute(input)
 if err != nil {
  return zero, fmt.Errorf("error executing first task '%s': %w", t.firstTask.Name(), err)
 }
 
 // 执行第二个任务
 output, err := t.lastTask.Execute(intermediate)
 if err != nil {
  return zero, fmt.Errorf("error executing last task '%s': %w", t.lastTask.Name(), err)
 }
 
 return output, nil
}

// Name 返回序列任务名称
func (t *SequentialTask[I, M, O]) Name() string {
 return t.name
}

// 可以组合多个任务构建复杂序列
func BuildSequence[I any, O any](name string, tasks []Task[I, O]) Task[I, O] {
 if len(tasks) == 0 {
  panic("cannot build sequence with empty tasks")
 }
 
 if len(tasks) == 1 {
  return tasks[0]
 }
 
 // 通过递归组合构建序列
 // 注意：这里需要一个更复杂的实现来处理不同类型的任务序列
 // 此处为简化示例
 return tasks[0]
}

// 使用示例
func SequenceExample() {
 task1 := &ConcreteTask[int, string]{
  name: "IntToString",
  execute: func(input int) (string, error) {
   return strconv.Itoa(input), nil
  },
 }
 
 task2 := &ConcreteTask[string, int]{
  name: "StringLength",
  execute: func(input string) (int, error) {
   return len(input), nil
  },
 }
 
 sequence := NewSequentialTask[int, string, int](
  "IntToLengthSequence", 
  task1, 
  task2,
 )
 
 result, err := sequence.Execute(42)
 if err != nil {
  fmt.Printf("Error: %v\n", err)
  return
 }
 
 fmt.Printf("Result: %v\n", result) // 输出: Result: 2
}

```

### 11.5.1.4.2 并行拆分模式 (Parallel Split Pattern)

#### 11.5.1.4.2.1 形式化定义

**定义 3.2** (并行拆分模式): 给定一个活动 $a_0$ 和一组活动 $\{a_1, a_2, ..., a_n\}$，并行拆分模式允许在 $a_0$ 完成后同时启动所有 $a_i (1 \leq i \leq n)$。在Petri网中，这可表示为：

1. 活动 $a_0$ 对应变迁 $t_0$
2. 活动 $a_i$ 对应变迁 $t_i (1 \leq i \leq n)$
3. 存在库所 $p$ 使得 $(t_0, p) \in F$ 且对于所有 $i (1 \leq i \leq n)$，$(p, t_i) \in F$

**定理 3.2**: 并行拆分模式满足以下属性：

- 并发性：所有 $a_i$ 可以并发执行，无需相互等待
- 无数据竞争：如果 $a_i$ 操作的数据集合互不相交，则不存在数据竞争
- 有界性：如果所有 $a_i$ 都是有界的，则并行拆分是有界的

#### 11.5.1.4.2.2 Golang实现

```go
package workflow

import (
 "fmt"
 "sync"
 "context"
)

// ParallelTask 表示并行拆分模式中的任务组
type ParallelTask[I any, O any] struct {
 name  string
 tasks []Task[I, O]
}

// 创建新的并行任务组
func NewParallelTask[I any, O any](name string, tasks []Task[I, O]) *ParallelTask[I, O] {
 return &ParallelTask[I, O]{
  name:  name,
  tasks: tasks,
 }
}

// Execute 并行执行所有任务
func (t *ParallelTask[I, O]) Execute(input I) ([]O, error) {
 var wg sync.WaitGroup
 results := make([]O, len(t.tasks))
 errors := make([]error, len(t.tasks))
 
 // 为每个任务启动一个goroutine
 for i, task := range t.tasks {
  wg.Add(1)
  go func(idx int, tsk Task[I, O]) {
   defer wg.Done()
   result, err := tsk.Execute(input)
   results[idx] = result
   errors[idx] = err
  }(i, task)
 }
 
 // 等待所有任务完成
 wg.Wait()
 
 // 检查错误
 for i, err := range errors {
  if err != nil {
   return nil, fmt.Errorf("task '%s' failed: %w", t.tasks[i].Name(), err)
  }
 }
 
 return results, nil
}

// Name 返回并行任务组名称
func (t *ParallelTask[I, O]) Name() string {
 return t.name
}

// 支持取消的并行任务执行
func (t *ParallelTask[I, O]) ExecuteWithContext(ctx context.Context, input I) ([]O, error) {
 var wg sync.WaitGroup
 results := make([]O, len(t.tasks))
 errors := make([]error, len(t.tasks))
 
 // 创建用于传递取消信号的channel
 done := make(chan struct{})
 defer close(done)
 
 // 为每个任务启动一个goroutine
 for i, task := range t.tasks {
  wg.Add(1)
  go func(idx int, tsk Task[I, O]) {
   defer wg.Done()
   
   // 使用带有上下文的任务执行
   if ctxTask, ok := tsk.(ContextTask[I, O]); ok {
    result, err := ctxTask.ExecuteWithContext(ctx, input)
    results[idx] = result
    errors[idx] = err
   } else {
    // 回退到普通执行，但检查上下文取消
    select {
    case <-ctx.Done():
     errors[idx] = ctx.Err()
    default:
     result, err := tsk.Execute(input)
     results[idx] = result
     errors[idx] = err
    }
   }
  }(i, task)
 }
 
 // 等待所有任务完成或上下文取消
 go func() {
  wg.Wait()
  close(done)
 }()
 
 select {
 case <-done:
  // 所有任务已完成
 case <-ctx.Done():
  // 上下文已取消
  return nil, ctx.Err()
 }
 
 // 检查错误
 for i, err := range errors {
  if err != nil {
   return nil, fmt.Errorf("task '%s' failed: %w", t.tasks[i].Name(), err)
  }
 }
 
 return results, nil
}

// 使用示例
func ParallelExample() {
 task1 := &ConcreteTask[int, int]{
  name: "Double",
  execute: func(input int) (int, error) {
   return input * 2, nil
  },
 }
 
 task2 := &ConcreteTask[int, int]{
  name: "Triple",
  execute: func(input int) (int, error) {
   return input * 3, nil
  },
 }
 
 task3 := &ConcreteTask[int, int]{
  name: "Square",
  execute: func(input int) (int, error) {
   return input * input, nil
  },
 }
 
 parallel := NewParallelTask[int, int](
  "ParallelComputation", 
  []Task[int, int]{task1, task2, task3},
 )
 
 results, err := parallel.Execute(5)
 if err != nil {
  fmt.Printf("Error: %v\n", err)
  return
 }
 
 fmt.Printf("Results: %v\n", results) // 输出: Results: [10 15 25]
}

```

### 11.5.1.4.3 同步模式 (Synchronization Pattern)

#### 11.5.1.4.3.1 形式化定义

**定义 3.3** (同步模式): 给定一组活动 $\{a_1, a_2, ..., a_n\}$ 和一个活动 $a_{n+1}$，同步模式要求所有 $a_i (1 \leq i \leq n)$ 都完成后才能启动 $a_{n+1}$。在Petri网中，这可表示为：

1. 活动 $a_i$ 对应变迁 $t_i (1 \leq i \leq n)$
2. 活动 $a_{n+1}$ 对应变迁 $t_{n+1}$
3. 存在库所 $p$ 使得对于所有 $i (1 \leq i \leq n)$，$(t_i, p) \in F$ 且 $(p, t_{n+1}) \in F$
4. $W(t_i, p) = 1$ 对于所有 $i (1 \leq i \leq n)$
5. $W(p, t_{n+1}) = n$

**定理 3.3**: 同步模式满足以下属性：

- 同步保证：$a_{n+1}$ 只有在所有 $a_i$ 完成后才会启动
- 无死锁：如果所有 $a_i$ 最终都能完成，则不会发生死锁
- 合成性：同步模式可以与其他模式组合，形成复杂的工作流

#### 11.5.1.4.3.2 Golang实现

```go
package workflow

import (
 "fmt"
 "sync"
)

// SyncTask 表示同步模式中的任务
type SyncTask[I any, M any, O any] struct {
 name       string
 tasks      []Task[I, M]
 finalTask  Task[[]M, O]
}

// 创建新的同步任务
func NewSyncTask[I any, M any, O any](
 name string,
 tasks []Task[I, M],
 finalTask Task[[]M, O],
) *SyncTask[I, M, O] {
 return &SyncTask[I, M, O]{
  name:       name,
  tasks:      tasks,
  finalTask:  finalTask,
 }
}

// Execute 执行同步任务
func (t *SyncTask[I, M, O]) Execute(input I) (O, error) {
 var zero O
 var wg sync.WaitGroup
 results := make([]M, len(t.tasks))
 errors := make([]error, len(t.tasks))
 
 // 并行执行所有前置任务
 for i, task := range t.tasks {
  wg.Add(1)
  go func(idx int, tsk Task[I, M]) {
   defer wg.Done()
   result, err := tsk.Execute(input)
   results[idx] = result
   errors[idx] = err
  }(i, task)
 }
 
 // 等待所有前置任务完成
 wg.Wait()
 
 // 检查前置任务错误
 for i, err := range errors {
  if err != nil {
   return zero, fmt.Errorf("task '%s' failed: %w", t.tasks[i].Name(), err)
  }
 }
 
 // 所有前置任务成功后，执行最终任务
 output, err := t.finalTask.Execute(results)
 if err != nil {
  return zero, fmt.Errorf("final task '%s' failed: %w", t.finalTask.Name(), err)
 }
 
 return output, nil
}

// Name 返回同步任务名称
func (t *SyncTask[I, M, O]) Name() string {
 return t.name
}

// 使用示例
func SynchronizationExample() {
 // 前置任务：计算不同的统计数据
 task1 := &ConcreteTask[[]int, float64]{
  name: "CalculateAverage",
  execute: func(input []int) (float64, error) {
   sum := 0
   for _, v := range input {
    sum += v
   }
   return float64(sum) / float64(len(input)), nil
  },
 }
 
 task2 := &ConcreteTask[[]int, int]{
  name: "FindMaximum",
  execute: func(input []int) (int, error) {
   max := input[0]
   for _, v := range input[1:] {
    if v > max {
     max = v
    }
   }
   return max, nil
  },
 }
 
 task3 := &ConcreteTask[[]int, int]{
  name: "FindMinimum",
  execute: func(input []int) (int, error) {
   min := input[0]
   for _, v := range input[1:] {
    if v < min {
     min = v
    }
   }
   return min, nil
  },
 }
 
 // 最终任务：生成报告
 finalTask := &ConcreteTask[[]interface{}, string]{
  name: "GenerateReport",
  execute: func(input []interface{}) (string, error) {
   avg := input[0].(float64)
   max := input[1].(int)
   min := input[2].(int)
   return fmt.Sprintf("Report - Avg: %.2f, Max: %d, Min: %d", avg, max, min), nil
  },
 }
 
 // 创建同步任务
 syncTask := NewSyncTask[[]int, interface{}, string](
  "DataAnalysisSync",
  []Task[[]int, interface{}]{
   WrapTask[[]int, float64, interface{}](task1),
   WrapTask[[]int, int, interface{}](task2),
   WrapTask[[]int, int, interface{}](task3),
  },
  finalTask,
 )
 
 data := []int{5, 2, 9, 1, 7, 3}
 report, err := syncTask.Execute(data)
 if err != nil {
  fmt.Printf("Error: %v\n", err)
  return
 }
 
 fmt.Println(report) // 输出: Report - Avg: 4.50, Max: 9, Min: 1
}

// 辅助函数，用于类型转换
func WrapTask[I any, O any, R any](task Task[I, O]) Task[I, R] {
 return &wrappedTask[I, O, R]{task: task}
}

type wrappedTask[I any, O any, R any] struct {
 task Task[I, O]
}

func (t *wrappedTask[I, O, R]) Execute(input I) (R, error) {
 result, err := t.task.Execute(input)
 if err != nil {
  var zero R
  return zero, err
 }
 
 // 类型转换
 return interface{}(result).(R), nil
}

func (t *wrappedTask[I, O, R]) Name() string {
 return t.task.Name()
}

```

## 11.5.1.5 4. 高级分支合并模式

### 11.5.1.5.1 多选模式

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

### 11.5.1.5.2 结构化同步合并模式

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

## 11.5.1.6 5. 结构模式

### 11.5.1.6.1 任意循环模式

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

## 11.5.1.7 6. 多实例模式

### 11.5.1.7.1 多实例不同步模式

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

## 11.5.1.8 7. 状态基础模式

### 11.5.1.8.1 延迟选择模式

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

### 11.5.1.8.2 里程碑模式

**定义 7.2** (里程碑模式): 里程碑模式是一个三元组 $M = (P, E, \Phi)$，其中：

- $P$ 是前置条件集合
- $E$ 是里程碑事件集合
- $\Phi: P \rightarrow E$ 是里程碑触发函数

```go
// 里程碑模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// Milestone 里程碑接口
type Milestone interface {
    ID() string
    IsReached(ctx context.Context, data interface{}) bool
    OnReached(ctx context.Context, data interface{}) error
}

// MilestoneWorkflow 里程碑工作流
type MilestoneWorkflow struct {
    milestones []Milestone
    activities map[string]Activity
    mu         sync.RWMutex
}

// NewMilestoneWorkflow 创建新的里程碑工作流
func NewMilestoneWorkflow() *MilestoneWorkflow {
    return &MilestoneWorkflow{
        milestones: make([]Milestone, 0),
        activities: make(map[string]Activity),
    }
}

// AddMilestone 添加里程碑
func (mw *MilestoneWorkflow) AddMilestone(milestone Milestone, activity Activity) {
    mw.mu.Lock()
    defer mw.mu.Unlock()
    
    mw.milestones = append(mw.milestones, milestone)
    mw.activities[milestone.ID()] = activity
}

// Execute 执行里程碑工作流
func (mw *MilestoneWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    mw.mu.RLock()
    defer mw.mu.RUnlock()
    
    result := data
    
    for _, milestone := range mw.milestones {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            if milestone.IsReached(ctx, result) {
                if err := milestone.OnReached(ctx, result); err != nil {
                    return nil, err
                }
                
                if activity, exists := mw.activities[milestone.ID()]; exists {
                    var err error
                    result, err = activity.Execute(ctx, result)
                    if err != nil {
                        return nil, fmt.Errorf("activity at milestone %s failed: %w", milestone.ID(), err)
                    }
                }
            }
        }
    }
    
    return result, nil
}

```

## 11.5.1.9 8. 资源模式

资源模式关注工作流执行过程中的资源分配、角色授权和任务委派等问题。这些模式确保工作流系统能够有效地管理和利用可用资源。

### 11.5.1.9.1 直接分配模式

**定义 8.1** (直接分配模式): 直接分配模式是一个三元组 $DA = (A, R, \alpha)$，其中：

- $A$ 是活动集合
- $R$ 是资源集合
- $\alpha: A \rightarrow R$ 是分配函数

```go
// 直接分配模式实现
package workflow

import (
    "context"
    "fmt"
)

// Resource 资源接口
type Resource interface {
    ID() string
    IsAvailable(ctx context.Context) bool
    Acquire(ctx context.Context) error
    Release(ctx context.Context) error
}

// DirectAllocationWorkflow 直接分配工作流
type DirectAllocationWorkflow struct {
    activity Activity
    resource Resource
}

// NewDirectAllocationWorkflow 创建新的直接分配工作流
func NewDirectAllocationWorkflow(activity Activity, resource Resource) *DirectAllocationWorkflow {
    return &DirectAllocationWorkflow{
        activity: activity,
        resource: resource,
    }
}

// Execute 执行直接分配工作流
func (daw *DirectAllocationWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    if !daw.resource.IsAvailable(ctx) {
        return nil, fmt.Errorf("resource %s is not available", daw.resource.ID())
    }
    
    if err := daw.resource.Acquire(ctx); err != nil {
        return nil, fmt.Errorf("failed to acquire resource: %w", err)
    }
    
    defer daw.resource.Release(ctx)
    
    return daw.activity.Execute(ctx, data)
}

```

### 11.5.1.9.2 基于角色分配模式

**定义 8.2** (基于角色分配模式): 基于角色分配模式是一个四元组 $RA = (A, R, L, \beta)$，其中：

- $A$ 是活动集合
- $R$ 是资源集合
- $L$ 是角色集合
- $\beta: A \times L \rightarrow 2^R$ 是角色映射函数

```go
// 基于角色分配模式实现
package workflow

import (
    "context"
    "fmt"
    "sync"
)

// Role 角色接口
type Role interface {
    ID() string
    GetResources(ctx context.Context) []Resource
}

// RoleBasedWorkflow 基于角色的工作流
type RoleBasedWorkflow struct {
    activity Activity
    role     Role
    mu       sync.Mutex
}

// NewRoleBasedWorkflow 创建新的基于角色的工作流
func NewRoleBasedWorkflow(activity Activity, role Role) *RoleBasedWorkflow {
    return &RoleBasedWorkflow{
        activity: activity,
        role:     role,
    }
}

// Execute 执行基于角色的工作流
func (rbw *RoleBasedWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    rbw.mu.Lock()
    defer rbw.mu.Unlock()
    
    resources := rbw.role.GetResources(ctx)
    if len(resources) == 0 {
        return nil, fmt.Errorf("no resources available for role %s", rbw.role.ID())
    }
    
    // 选择第一个可用的资源
    var selectedResource Resource
    for _, resource := range resources {
        if resource.IsAvailable(ctx) {
            selectedResource = resource
            break
        }
    }
    
    if selectedResource == nil {
        return nil, fmt.Errorf("no available resources found for role %s", rbw.role.ID())
    }
    
    if err := selectedResource.Acquire(ctx); err != nil {
        return nil, fmt.Errorf("failed to acquire resource: %w", err)
    }
    
    defer selectedResource.Release(ctx)
    
    return rbw.activity.Execute(ctx, data)
}

```

### 11.5.1.9.3 授权模式

**定义 8.3** (授权模式): 授权模式是一个四元组 $AU = (A, R, P, \gamma)$，其中：

- $A$ 是活动集合
- $R$ 是资源集合
- $P$ 是权限集合
- $\gamma: A \times R \rightarrow 2^P$ 是授权函数

```go
// 授权模式实现
package workflow

import (
    "context"
    "fmt"
)

// Permission 权限接口
type Permission interface {
    ID() string
    Check(ctx context.Context, resource Resource) bool
}

// AuthorizationWorkflow 授权工作流
type AuthorizationWorkflow struct {
    activity    Activity
    resource    Resource
    permissions []Permission
}

// NewAuthorizationWorkflow 创建新的授权工作流
func NewAuthorizationWorkflow(activity Activity, resource Resource, permissions []Permission) *AuthorizationWorkflow {
    return &AuthorizationWorkflow{
        activity:    activity,
        resource:    resource,
        permissions: permissions,
    }
}

// Execute 执行授权工作流
func (aw *AuthorizationWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    // 检查权限
    for _, permission := range aw.permissions {
        if !permission.Check(ctx, aw.resource) {
            return nil, fmt.Errorf("permission %s not granted for resource %s", permission.ID(), aw.resource.ID())
        }
    }
    
    // 获取资源
    if err := aw.resource.Acquire(ctx); err != nil {
        return nil, fmt.Errorf("failed to acquire resource: %w", err)
    }
    
    defer aw.resource.Release(ctx)
    
    // 执行活动
    return aw.activity.Execute(ctx, data)
}

```

## 11.5.1.10 9. 数据流模式

数据流模式关注工作流执行过程中的数据转换、传递和处理等问题，确保数据在工作流各个环节的正确流动和处理。

### 11.5.1.10.1 数据传递模式

**定义 9.1** (数据传递模式): 数据传递模式是一个四元组 $DP = (A, D, \Phi, \Psi)$，其中：

- $A = \{a_1, a_2, ..., a_n\}$ 是活动集合
- $D = \{d_1, d_2, ..., d_m\}$ 是数据对象集合
- $\Phi: A \rightarrow 2^D$ 是输入映射函数
- $\Psi: A \rightarrow 2^D$ 是输出映射函数

```go
// 数据传递模式实现
package workflow

import (
    "context"
    "fmt"
)

// DataObject 数据对象接口
type DataObject interface {
    ID() string
    GetValue() interface{}
    SetValue(value interface{}) error
}

// DataPassingWorkflow 数据传递工作流
type DataPassingWorkflow struct {
    activities       []Activity
    inputMapping     map[string][]string
    outputMapping    map[string][]string
}

// NewDataPassingWorkflow 创建新的数据传递工作流
func NewDataPassingWorkflow() *DataPassingWorkflow {
    return &DataPassingWorkflow{
        activities:    make([]Activity, 0),
        inputMapping:  make(map[string][]string),
        outputMapping: make(map[string][]string),
    }
}

// AddActivity 添加活动和数据映射
func (dpw *DataPassingWorkflow) AddActivity(activity Activity, inputs []string, outputs []string) {
    dpw.activities = append(dpw.activities, activity)
    dpw.inputMapping[activity.ID()] = inputs
    dpw.outputMapping[activity.ID()] = outputs
}

// Execute 执行数据传递工作流
func (dpw *DataPassingWorkflow) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    for k, v := range data {
        result[k] = v
    }
    
    for _, activity := range dpw.activities {
        // 准备活动输入数据
        activityInput := make(map[string]interface{})
        for _, inputID := range dpw.inputMapping[activity.ID()] {
            if value, exists := result[inputID]; exists {
                activityInput[inputID] = value
            }
        }
        
        // 执行活动
        activityOutput, err := activity.Execute(ctx, activityInput)
        if err != nil {
            return nil, fmt.Errorf("activity %s execution failed: %w", activity.ID(), err)
        }
        
        // 处理活动输出数据
        outputMap, ok := activityOutput.(map[string]interface{})
        if !ok {
            return nil, fmt.Errorf("activity %s output is not a map", activity.ID())
        }
        
        // 更新数据存储
        for _, outputID := range dpw.outputMapping[activity.ID()] {
            if value, exists := outputMap[outputID]; exists {
                result[outputID] = value
            }
        }
    }
    
    return result, nil
}

```

### 11.5.1.10.2 数据转换模式

**定义 9.2** (数据转换模式): 数据转换模式是一个四元组 $DT = (D_{in}, D_{out}, T, \xi)$，其中：

- $D_{in}$ 是输入数据类型
- $D_{out}$ 是输出数据类型
- $T$ 是转换规则集合
- $\xi: D_{in} \times T \rightarrow D_{out}$ 是转换函数

```go
// 数据转换模式实现
package workflow

import (
    "context"
    "fmt"
)

// TransformationRule 转换规则接口
type TransformationRule interface {
    ID() string
    Apply(input interface{}) (interface{}, error)
}

// DataTransformationWorkflow 数据转换工作流
type DataTransformationWorkflow struct {
    rules []TransformationRule
}

// NewDataTransformationWorkflow 创建新的数据转换工作流
func NewDataTransformationWorkflow(rules []TransformationRule) *DataTransformationWorkflow {
    return &DataTransformationWorkflow{
        rules: rules,
    }
}

// Execute 执行数据转换工作流
func (dtw *DataTransformationWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    result := data
    
    for _, rule := range dtw.rules {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            var err error
            result, err = rule.Apply(result)
            if err != nil {
                return nil, fmt.Errorf("transformation rule %s failed: %w", rule.ID(), err)
            }
        }
    }
    
    return result, nil
}

// 基础转换规则
type BaseTransformationRule struct {
    id        string
    transform func(input interface{}) (interface{}, error)
}

// NewBaseTransformationRule 创建新的基础转换规则
func NewBaseTransformationRule(id string, transform func(input interface{}) (interface{}, error)) *BaseTransformationRule {
    return &BaseTransformationRule{
        id:        id,
        transform: transform,
    }
}

// ID 获取规则ID
func (btr *BaseTransformationRule) ID() string {
    return btr.id
}

// Apply 应用转换规则
func (btr *BaseTransformationRule) Apply(input interface{}) (interface{}, error) {
    return btr.transform(input)
}

```

### 11.5.1.10.3 数据路由模式

**定义 9.3** (数据路由模式): 数据路由模式是一个四元组 $DR = (D, R, C, \lambda)$，其中：

- $D$ 是数据对象
- $R = \{r_1, r_2, ..., r_n\}$ 是路由目标集合
- $C = \{c_1, c_2, ..., c_n\}$ 是路由条件集合
- $\lambda: D \times C \rightarrow R$ 是路由函数

```go
// 数据路由模式实现
package workflow

import (
    "context"
    "fmt"
)

// RoutingRule 路由规则接口
type RoutingRule interface {
    ID() string
    Matches(ctx context.Context, data interface{}) bool
    GetTargetActivity() Activity
}

// DataRoutingWorkflow 数据路由工作流
type DataRoutingWorkflow struct {
    rules           []RoutingRule
    defaultActivity Activity
}

// NewDataRoutingWorkflow 创建新的数据路由工作流
func NewDataRoutingWorkflow(defaultActivity Activity) *DataRoutingWorkflow {
    return &DataRoutingWorkflow{
        rules:           make([]RoutingRule, 0),
        defaultActivity: defaultActivity,
    }
}

// AddRule 添加路由规则
func (drw *DataRoutingWorkflow) AddRule(rule RoutingRule) {
    drw.rules = append(drw.rules, rule)
}

// Execute 执行数据路由工作流
func (drw *DataRoutingWorkflow) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    // 查找匹配的路由规则
    for _, rule := range drw.rules {
        if rule.Matches(ctx, data) {
            activity := rule.GetTargetActivity()
            return activity.Execute(ctx, data)
        }
    }
    
    // 使用默认活动
    if drw.defaultActivity != nil {
        return drw.defaultActivity.Execute(ctx, data)
    }
    
    return nil, fmt.Errorf("no matching routing rule found")
}

// 基础路由规则
type BaseRoutingRule struct {
    id       string
    predicate func(ctx context.Context, data interface{}) bool
    activity  Activity
}

// NewBaseRoutingRule 创建新的基础路由规则
func NewBaseRoutingRule(id string, predicate func(ctx context.Context, data interface{}) bool, activity Activity) *BaseRoutingRule {
    return &BaseRoutingRule{
        id:        id,
        predicate: predicate,
        activity:  activity,
    }
}

// ID 获取规则ID
func (brr *BaseRoutingRule) ID() string {
    return brr.id
}

// Matches 检查数据是否匹配规则
func (brr *BaseRoutingRule) Matches(ctx context.Context, data interface{}) bool {
    return brr.predicate(ctx, data)
}

// GetTargetActivity 获取目标活动
func (brr *BaseRoutingRule) GetTargetActivity() Activity {
    return brr.activity
}

```

## 11.5.1.11 10. Petri网建模分析

### 11.5.1.11.1 形式化定义

**定义 10.1** (Petri网): Petri网是一个五元组 $PN = (P, T, F, W, M_0)$，其中：

- $P$ 是库所(Place)集合，表示状态或条件
- $T$ 是变迁(Transition)集合，表示活动或事件
- $F \subseteq (P \times T) \cup (T \times P)$ 是弧(Arc)集合
- $W: F \rightarrow \mathbb{N}^+$ 是权重函数
- $M_0: P \rightarrow \mathbb{N}$ 是初始标记(Initial marking)

工作流模式可以映射到相应的Petri网结构，从而利用Petri网理论进行形式化分析。

## 11.5.1.12 11. 工作流引擎实现

### 11.5.1.12.1 工作流引擎架构

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

## 11.5.1.13 12. 模式集成案例

### 11.5.1.13.1 订单处理工作流

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

### 11.5.1.13.2 审批工作流

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
