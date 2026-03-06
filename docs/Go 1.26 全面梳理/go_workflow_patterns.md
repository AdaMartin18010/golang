# Go 1.23 工作流设计模式全面指南

> 本文档系统梳理 **Go 1.23** 语言实现的43种工作流设计模式中的23种可判断模式，提供完整的实现方案。
>
> **Go 1.23 更新**：
>
> - 利用迭代器模式简化工作流节点遍历
> - 使用 `iter.Seq` / `iter.Seq2` 实现工作流状态迭代
> - 结合 `slices` / `maps` 迭代器优化工作流数据管理

---

## 目录

- [Go 1.23 工作流设计模式全面指南](#go-123-工作流设计模式全面指南)
  - [目录](#目录)
  - [1. 基础控制流模式](#1-基础控制流模式)
    - [1.1 顺序（Sequence）](#11-顺序sequence)
      - [概念定义](#概念定义)
      - [BPMN描述](#bpmn描述)
      - [流程图](#流程图)
      - [Go实现](#go实现)
      - [完整示例](#完整示例)
      - [反例说明](#反例说明)
      - [适用场景](#适用场景)
      - [与其他模式关系](#与其他模式关系)
    - [1.2 并行分支（Parallel Split）](#12-并行分支parallel-split)
      - [概念定义](#概念定义-1)
      - [BPMN描述](#bpmn描述-1)
      - [流程图](#流程图-1)
      - [Go实现](#go实现-1)
      - [完整示例](#完整示例-1)
      - [反例说明](#反例说明-1)
      - [适用场景](#适用场景-1)
      - [与其他模式关系](#与其他模式关系-1)
    - [1.3 同步（Synchronization）](#13-同步synchronization)
      - [概念定义](#概念定义-2)
      - [BPMN描述](#bpmn描述-2)
      - [流程图](#流程图-2)
      - [Go实现](#go实现-2)
      - [完整示例](#完整示例-2)
      - [反例说明](#反例说明-2)
      - [适用场景](#适用场景-2)
      - [与其他模式关系](#与其他模式关系-2)
    - [1.4 排他选择（Exclusive Choice）](#14-排他选择exclusive-choice)
      - [概念定义](#概念定义-3)
      - [BPMN描述](#bpmn描述-3)
      - [流程图](#流程图-3)
      - [Go实现](#go实现-3)
      - [完整示例](#完整示例-3)
      - [反例说明](#反例说明-3)
      - [适用场景](#适用场景-3)
      - [与其他模式关系](#与其他模式关系-3)
    - [1.5 简单合并（Simple Merge）](#15-简单合并simple-merge)
      - [概念定义](#概念定义-4)
      - [BPMN描述](#bpmn描述-4)
      - [流程图](#流程图-4)
      - [Go实现](#go实现-4)
      - [完整示例](#完整示例-4)
      - [反例说明](#反例说明-4)
      - [适用场景](#适用场景-4)
      - [与其他模式关系](#与其他模式关系-4)
  - [2. 高级分支与同步模式](#2-高级分支与同步模式)
    - [2.1 多选（Multi-Choice）](#21-多选multi-choice)
      - [概念定义](#概念定义-5)
      - [BPMN描述](#bpmn描述-5)
      - [流程图](#流程图-5)
      - [Go实现](#go实现-5)
      - [完整示例](#完整示例-5)
      - [反例说明](#反例说明-5)
      - [适用场景](#适用场景-5)
      - [与其他模式关系](#与其他模式关系-5)
    - [2.2 同步合并（Synchronizing Merge）](#22-同步合并synchronizing-merge)
      - [概念定义](#概念定义-6)
      - [BPMN描述](#bpmn描述-6)
      - [流程图](#流程图-6)
      - [Go实现](#go实现-6)
      - [完整示例](#完整示例-6)
      - [反例说明](#反例说明-6)
      - [适用场景](#适用场景-6)
      - [与其他模式关系](#与其他模式关系-6)
    - [2.3 多合并（Multi-Merge）](#23-多合并multi-merge)
      - [概念定义](#概念定义-7)
      - [BPMN描述](#bpmn描述-7)
      - [流程图](#流程图-7)
      - [Go实现](#go实现-7)
      - [完整示例](#完整示例-7)
      - [反例说明](#反例说明-7)
      - [适用场景](#适用场景-7)
      - [与其他模式关系](#与其他模式关系-7)
    - [2.4 鉴别器（Discriminator）](#24-鉴别器discriminator)
      - [概念定义](#概念定义-8)
      - [BPMN描述](#bpmn描述-8)
      - [流程图](#流程图-8)
      - [Go实现](#go实现-8)
      - [完整示例](#完整示例-8)
      - [反例说明](#反例说明-8)
      - [适用场景](#适用场景-8)
      - [与其他模式关系](#与其他模式关系-8)
    - [2.5 N选M（N-out-of-M Join）](#25-n选mn-out-of-m-join)
      - [概念定义](#概念定义-9)
      - [BPMN描述](#bpmn描述-9)
      - [流程图](#流程图-9)
      - [Go实现](#go实现-9)
      - [完整示例](#完整示例-9)
      - [反例说明](#反例说明-9)
      - [适用场景](#适用场景-9)
      - [与其他模式关系](#与其他模式关系-9)
  - [3. 结构化模式](#3-结构化模式)
    - [3.1 任意循环（Arbitrary Cycles）](#31-任意循环arbitrary-cycles)
      - [概念定义](#概念定义-10)
      - [BPMN描述](#bpmn描述-10)
      - [流程图](#流程图-10)
      - [Go实现](#go实现-10)
      - [完整示例](#完整示例-10)
      - [反例说明](#反例说明-10)
      - [适用场景](#适用场景-10)
      - [与其他模式关系](#与其他模式关系-10)
    - [3.2 隐式终止（Implicit Termination）](#32-隐式终止implicit-termination)
      - [概念定义](#概念定义-11)
      - [BPMN描述](#bpmn描述-11)
      - [流程图](#流程图-11)
      - [Go实现](#go实现-11)
      - [完整示例](#完整示例-11)
      - [反例说明](#反例说明-11)
      - [适用场景](#适用场景-11)
      - [与其他模式关系](#与其他模式关系-11)
  - [4. 多实例模式](#4-多实例模式)
    - [4.1 多实例无同步（Multiple Instances without Synchronization）](#41-多实例无同步multiple-instances-without-synchronization)
      - [概念定义](#概念定义-12)
      - [BPMN描述](#bpmn描述-12)
      - [流程图](#流程图-12)
      - [Go实现](#go实现-12)
      - [完整示例](#完整示例-12)
      - [反例说明](#反例说明-12)
      - [适用场景](#适用场景-12)
      - [与其他模式关系](#与其他模式关系-12)
    - [4.2 多实例需同步（Multiple Instances with a Priori Design Time Knowledge）](#42-多实例需同步multiple-instances-with-a-priori-design-time-knowledge)
      - [概念定义](#概念定义-13)
      - [BPMN描述](#bpmn描述-13)
      - [流程图](#流程图-13)
      - [Go实现](#go实现-13)
      - [完整示例](#完整示例-13)
      - [反例说明](#反例说明-13)
      - [适用场景](#适用场景-13)
      - [与其他模式关系](#与其他模式关系-13)
    - [4.3 多实例运行时确定（Multiple Instances with a Priori Runtime Knowledge）](#43-多实例运行时确定multiple-instances-with-a-priori-runtime-knowledge)
      - [概念定义](#概念定义-14)
      - [BPMN描述](#bpmn描述-14)
      - [流程图](#流程图-14)
      - [Go实现](#go实现-14)
      - [完整示例](#完整示例-14)
      - [反例说明](#反例说明-14)
      - [适用场景](#适用场景-14)
      - [与其他模式关系](#与其他模式关系-14)
  - [5. 状态模式](#5-状态模式)
    - [5.1 延迟选择（Deferred Choice）](#51-延迟选择deferred-choice)
      - [概念定义](#概念定义-15)
      - [BPMN描述](#bpmn描述-15)
      - [流程图](#流程图-15)
      - [Go实现](#go实现-15)
      - [完整示例](#完整示例-15)
      - [反例说明](#反例说明-15)
      - [适用场景](#适用场景-15)
      - [与其他模式关系](#与其他模式关系-15)
    - [5.2 交错并行路由（Interleaved Parallel Routing）](#52-交错并行路由interleaved-parallel-routing)
      - [概念定义](#概念定义-16)
      - [BPMN描述](#bpmn描述-16)
      - [流程图](#流程图-16)
      - [Go实现](#go实现-16)
      - [完整示例](#完整示例-16)
      - [反例说明](#反例说明-16)
      - [适用场景](#适用场景-16)
      - [与其他模式关系](#与其他模式关系-16)
    - [5.3 里程碑（Milestone）](#53-里程碑milestone)
      - [概念定义](#概念定义-17)
      - [BPMN描述](#bpmn描述-17)
      - [流程图](#流程图-17)
      - [Go实现](#go实现-17)
      - [完整示例](#完整示例-17)
      - [反例说明](#反例说明-17)
      - [适用场景](#适用场景-17)
      - [与其他模式关系](#与其他模式关系-17)
  - [6. 取消与补偿模式](#6-取消与补偿模式)
    - [6.1 取消任务（Cancel Activity）](#61-取消任务cancel-activity)
      - [概念定义](#概念定义-18)
      - [BPMN描述](#bpmn描述-18)
      - [流程图](#流程图-18)
      - [Go实现](#go实现-18)
      - [完整示例](#完整示例-18)
      - [反例说明](#反例说明-18)
      - [适用场景](#适用场景-18)
      - [与其他模式关系](#与其他模式关系-18)
    - [6.2 取消案例（Cancel Case）](#62-取消案例cancel-case)
      - [概念定义](#概念定义-19)
      - [BPMN描述](#bpmn描述-19)
      - [流程图](#流程图-19)
      - [Go实现](#go实现-19)
      - [完整示例](#完整示例-19)
      - [反例说明](#反例说明-19)
      - [适用场景](#适用场景-19)
      - [与其他模式关系](#与其他模式关系-19)
  - [7. 其他可判断模式](#7-其他可判断模式)
    - [7.1 递归（Recursion）](#71-递归recursion)
      - [概念定义](#概念定义-20)
      - [BPMN描述](#bpmn描述-20)
      - [流程图](#流程图-20)
      - [Go实现](#go实现-20)
      - [完整示例](#完整示例-20)
      - [反例说明](#反例说明-20)
      - [适用场景](#适用场景-20)
      - [与其他模式关系](#与其他模式关系-20)
    - [7.2 临时触发器（Transient Trigger）](#72-临时触发器transient-trigger)
      - [概念定义](#概念定义-21)
      - [BPMN描述](#bpmn描述-21)
      - [流程图](#流程图-21)
      - [Go实现](#go实现-21)
      - [完整示例](#完整示例-21)
      - [反例说明](#反例说明-21)
      - [适用场景](#适用场景-21)
      - [与其他模式关系](#与其他模式关系-21)
    - [7.3 持久触发器（Persistent Trigger）](#73-持久触发器persistent-trigger)
      - [概念定义](#概念定义-22)
      - [BPMN描述](#bpmn描述-22)
      - [流程图](#流程图-22)
      - [Go实现](#go实现-22)
      - [完整示例](#完整示例-22)
      - [反例说明](#反例说明-22)
      - [适用场景](#适用场景-22)
      - [与其他模式关系](#与其他模式关系-22)
  - [总结](#总结)
    - [模式分类统计](#模式分类统计)
    - [实现要点](#实现要点)
    - [选择建议](#选择建议)

---

## 1. 基础控制流模式

### 1.1 顺序（Sequence）

#### 概念定义

**顺序模式**是工作流中最基本的控制流模式，指任务按照预定义的先后顺序依次执行。前一个任务完成后，后一个任务才能开始。这是所有工作流的基础构建块。

**形式化定义**：给定任务集合 T = {t₁, t₂, ..., tₙ}，顺序模式定义了一个全序关系 ≺，使得 t₁ ≺ t₂ ≺ ... ≺ tₙ，其中 tᵢ ≺ tⱼ 表示任务 tᵢ 必须在任务 tⱼ 之前完成。

#### BPMN描述

在BPMN中，顺序模式通过**顺序流（Sequence Flow）**表示：

```
[开始事件] → [任务A] → [任务B] → [任务C] → [结束事件]
```

- 使用实线箭头连接各个活动
- 箭头方向表示执行顺序
- 无分支、无合并的直线流程

#### 流程图

```
┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐
│  开始   │────→│  任务A  │────→│  任务B  │────→│  任务C  │
└─────────┘     └─────────┘     └─────────┘     └─────────┘
                                                    │
                                                    ▼
                                               ┌─────────┐
                                               │  结束   │
                                               └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// Task 表示一个工作流任务
type Task struct {
 Name string
 Fn   func(ctx context.Context) error
}

// Sequence 顺序执行器
type Sequence struct {
 tasks []Task
}

// NewSequence 创建顺序执行器
func NewSequence() *Sequence {
 return &Sequence{
  tasks: make([]Task, 0),
 }
}

// Add 添加任务
func (s *Sequence) Add(name string, fn func(ctx context.Context) error) {
 s.tasks = append(s.tasks, Task{Name: name, Fn: fn})
}

// Execute 顺序执行所有任务
func (s *Sequence) Execute(ctx context.Context) error {
 for i, task := range s.tasks {
  log.Printf("[顺序执行] 开始任务 %d/%d: %s", i+1, len(s.tasks), task.Name)

  if err := task.Fn(ctx); err != nil {
   return fmt.Errorf("任务 '%s' 执行失败: %w", task.Name, err)
  }

  log.Printf("[顺序执行] 完成任务 %d/%d: %s", i+1, len(s.tasks), task.Name)
 }
 return nil
}

// 使用示例
func main() {
 seq := NewSequence()

 // 添加顺序任务
 seq.Add("数据验证", func(ctx context.Context) error {
  fmt.Println("→ 执行数据验证...")
  time.Sleep(100 * time.Millisecond)
  return nil
 })

 seq.Add("数据转换", func(ctx context.Context) error {
  fmt.Println("→ 执行数据转换...")
  time.Sleep(100 * time.Millisecond)
  return nil
 })

 seq.Add("数据存储", func(ctx context.Context) error {
  fmt.Println("→ 执行数据存储...")
  time.Sleep(100 * time.Millisecond)
  return nil
 })

 ctx := context.Background()
 if err := seq.Execute(ctx); err != nil {
  log.Fatal(err)
 }
 fmt.Println("✓ 所有任务顺序执行完成")
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// OrderProcessor 订单处理工作流
type OrderProcessor struct {
 sequence *Sequence
}

// NewOrderProcessor 创建订单处理器
func NewOrderProcessor() *OrderProcessor {
 op := &OrderProcessor{
  sequence: NewSequence(),
 }

 // 定义订单处理流程
 op.sequence.Add("验证订单", op.validateOrder)
 op.sequence.Add("检查库存", op.checkInventory)
 op.sequence.Add("计算价格", op.calculatePrice)
 op.sequence.Add("处理支付", op.processPayment)
 op.sequence.Add("生成发货单", op.generateShipping)

 return op
}

func (op *OrderProcessor) validateOrder(ctx context.Context) error {
 fmt.Println("[1/5] 验证订单信息...")
 time.Sleep(200 * time.Millisecond)
 return nil
}

func (op *OrderProcessor) checkInventory(ctx context.Context) error {
 fmt.Println("[2/5] 检查库存可用性...")
 time.Sleep(300 * time.Millisecond)
 return nil
}

func (op *OrderProcessor) calculatePrice(ctx context.Context) error {
 fmt.Println("[3/5] 计算订单价格...")
 time.Sleep(150 * time.Millisecond)
 return nil
}

func (op *OrderProcessor) processPayment(ctx context.Context) error {
 fmt.Println("[4/5] 处理支付...")
 time.Sleep(400 * time.Millisecond)
 return nil
}

func (op *OrderProcessor) generateShipping(ctx context.Context) error {
 fmt.Println("[5/5] 生成发货单...")
 time.Sleep(100 * time.Millisecond)
 return nil
}

func (op *OrderProcessor) Process(ctx context.Context) error {
 return op.sequence.Execute(ctx)
}

func main() {
 processor := NewOrderProcessor()
 ctx := context.Background()

 start := time.Now()
 if err := processor.Process(ctx); err != nil {
  log.Fatal(err)
 }
 fmt.Printf("\n订单处理完成，耗时: %v\n", time.Since(start))
}
```

#### 反例说明

**错误实现1：并发执行顺序任务**

```go
// ❌ 错误：使用goroutine并发执行本应顺序的任务
func WrongConcurrentExecution() {
 var wg sync.WaitGroup
 wg.Add(3)

 go func() { defer wg.Done(); validateOrder() }()
 go func() { defer wg.Done(); processPayment() }()  // 错误：支付可能在验证前执行
 go func() { defer wg.Done(); generateShipping() }() // 错误：发货可能在支付前执行

 wg.Wait()
}
```

**问题**：破坏了任务间的依赖关系，可能导致数据不一致或业务逻辑错误。

**错误实现2：缺少错误处理**

```go
// ❌ 错误：忽略错误继续执行
func WrongIgnoreErrors() {
 validateOrder()      // 可能失败
 processPayment()     // 即使验证失败也会执行
 generateShipping()   // 即使支付失败也会执行
}
```

**问题**：错误被忽略，导致无效操作继续执行，可能产生脏数据。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 数据处理流程 | ETL流程：提取→转换→加载 |
| 订单处理 | 验证→库存检查→支付→发货 |
| 审批流程 | 提交→初审→复审→终审 |
| 构建流水线 | 编译→测试→打包→部署 |
| 事务处理 | 开始→操作→提交/回滚 |

#### 与其他模式关系

- **基础模式**：所有其他模式的基础构建块
- **与并行分支**：顺序是并行的对立面，两者常组合使用
- **与子流程**：子流程内部通常使用顺序模式组织任务

---

### 1.2 并行分支（Parallel Split）

#### 概念定义

**并行分支模式**允许工作流在执行过程中的某一点将控制流分成多个并行的执行路径，这些路径可以同时执行。这是实现并行处理、提高效率的关键模式。

**形式化定义**：给定任务 t，并行分支模式创建一组并行任务集合 P = {p₁, p₂, ..., pₙ}，使得 ∀pᵢ, pⱼ ∈ P，pᵢ 和 pⱼ 可以并发执行，且执行顺序无关。

#### BPMN描述

在BPMN中，并行分支通过**并行网关（Parallel Gateway）**实现：

```
                    ┌→ [任务A] →┐
[前置任务] → [◇] →┼→ [任务B] →┼→ [同步点]
                    └→ [任务C] →┘
```

- 并行网关（◇）表示分支点
- 所有出分支同时激活
- 需要后续同步网关汇聚

#### 流程图

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  开始   │────→│  任务A  │     │  任务B  │
└─────────┘     └────┬────┘     └────┬────┘
                     │               │
                ┌────┴───────────────┴────┐
                │      并行网关(分支)      │
                └────┬───────────────┬────┘
                     │               │
                     ▼               ▼
               ┌─────────┐     ┌─────────┐
               │ 分支任务1│     │ 分支任务2│
               └────┬────┘     └────┬────┘
                    │               │
                    └───────┬───────┘
                            ▼
                     ┌─────────────┐
                     │  同步网关   │
                     └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// ParallelBranch 并行分支执行器
type ParallelBranch struct {
 branches []func(ctx context.Context) error
}

// NewParallelBranch 创建并行分支执行器
func NewParallelBranch() *ParallelBranch {
 return &ParallelBranch{
  branches: make([]func(ctx context.Context) error, 0),
 }
}

// Add 添加并行分支
func (p *ParallelBranch) Add(fn func(ctx context.Context) error) {
 p.branches = append(p.branches, fn)
}

// Execute 并行执行所有分支
func (p *ParallelBranch) Execute(ctx context.Context) error {
 if len(p.branches) == 0 {
  return nil
 }

 var wg sync.WaitGroup
 errChan := make(chan error, len(p.branches))

 log.Printf("[并行分支] 启动 %d 个并行任务", len(p.branches))

 for i, branch := range p.branches {
  wg.Add(1)
  go func(index int, fn func(ctx context.Context) error) {
   defer wg.Done()

   log.Printf("[并行分支] 任务 %d 开始执行", index+1)
   start := time.Now()

   if err := fn(ctx); err != nil {
    errChan <- fmt.Errorf("分支 %d 执行失败: %w", index+1, err)
    return
   }

   log.Printf("[并行分支] 任务 %d 完成，耗时: %v", index+1, time.Since(start))
  }(i, branch)
 }

 wg.Wait()
 close(errChan)

 // 收集所有错误
 var errs []error
 for err := range errChan {
  errs = append(errs, err)
 }

 if len(errs) > 0 {
  return fmt.Errorf("并行执行出现 %d 个错误: %v", len(errs), errs)
 }

 return nil
}

// ExecuteWithCancel 并行执行，任一失败则取消其他
func (p *ParallelBranch) ExecuteWithCancel(ctx context.Context) error {
 if len(p.branches) == 0 {
  return nil
 }

 ctx, cancel := context.WithCancel(ctx)
 defer cancel()

 var wg sync.WaitGroup
 errChan := make(chan error, 1) // 缓冲1，避免goroutine泄漏

 for i, branch := range p.branches {
  wg.Add(1)
  go func(index int, fn func(ctx context.Context) error) {
   defer wg.Done()

   if err := fn(ctx); err != nil {
    select {
    case errChan <- fmt.Errorf("分支 %d 失败: %w", index+1, err):
     cancel() // 通知其他分支取消
    default:
    }
   }
  }(i, branch)
 }

 // 等待所有分支完成或出错
 go func() {
  wg.Wait()
  close(errChan)
 }()

 if err := <-errChan; err != nil {
  return err
 }
 return nil
}

// 使用示例
func main() {
 parallel := NewParallelBranch()

 // 添加并行任务
 parallel.Add(func(ctx context.Context) error {
  fmt.Println("→ 并行任务1: 发送邮件通知")
  time.Sleep(200 * time.Millisecond)
  return nil
 })

 parallel.Add(func(ctx context.Context) error {
  fmt.Println("→ 并行任务2: 更新用户统计")
  time.Sleep(300 * time.Millisecond)
  return nil
 })

 parallel.Add(func(ctx context.Context) error {
  fmt.Println("→ 并行任务3: 记录操作日志")
  time.Sleep(100 * time.Millisecond)
  return nil
 })

 ctx := context.Background()
 start := time.Now()

 if err := parallel.Execute(ctx); err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n✓ 所有并行任务完成，总耗时: %v\n", time.Since(start))
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// DataPipeline 数据处理并行流水线
type DataPipeline struct {
 data []int
}

// NewDataPipeline 创建数据流水线
func NewDataPipeline(data []int) *DataPipeline {
 return &DataPipeline{data: data}
}

// ProcessParallel 并行处理数据
func (dp *DataPipeline) ProcessParallel(ctx context.Context) ([]int, []int, []int, error) {
 var wg sync.WaitGroup

 // 三个并行处理分支
 var filtered []int
 var transformed []int
 var analyzed []int

 var err1, err2, err3 error

 wg.Add(3)

 // 分支1: 数据过滤
 go func() {
  defer wg.Done()
  filtered, err1 = dp.filterData(ctx)
 }()

 // 分支2: 数据转换
 go func() {
  defer wg.Done()
  transformed, err2 = dp.transformData(ctx)
 }()

 // 分支3: 数据分析
 go func() {
  defer wg.Done()
  analyzed, err3 = dp.analyzeData(ctx)
 }()

 wg.Wait()

 // 检查错误
 if err1 != nil {
  return nil, nil, nil, fmt.Errorf("过滤失败: %w", err1)
 }
 if err2 != nil {
  return nil, nil, nil, fmt.Errorf("转换失败: %w", err2)
 }
 if err3 != nil {
  return nil, nil, nil, fmt.Errorf("分析失败: %w", err3)
 }

 return filtered, transformed, analyzed, nil
}

func (dp *DataPipeline) filterData(ctx context.Context) ([]int, error) {
 fmt.Println("[分支1] 开始数据过滤...")
 time.Sleep(200 * time.Millisecond)

 var result []int
 for _, v := range dp.data {
  if v > 10 {
   result = append(result, v)
  }
 }
 fmt.Printf("[分支1] 过滤完成: %d 条数据\n", len(result))
 return result, nil
}

func (dp *DataPipeline) transformData(ctx context.Context) ([]int, error) {
 fmt.Println("[分支2] 开始数据转换...")
 time.Sleep(300 * time.Millisecond)

 result := make([]int, len(dp.data))
 for i, v := range dp.data {
  result[i] = v * 2
 }
 fmt.Printf("[分支2] 转换完成: %d 条数据\n", len(result))
 return result, nil
}

func (dp *DataPipeline) analyzeData(ctx context.Context) ([]int, error) {
 fmt.Println("[分支3] 开始数据分析...")
 time.Sleep(150 * time.Millisecond)

 sum := 0
 for _, v := range dp.data {
  sum += v
 }
 avg := sum / len(dp.data)
 fmt.Printf("[分支3] 分析完成: 平均值=%d\n", avg)
 return []int{sum, avg}, nil
}

func main() {
 data := []int{5, 15, 8, 20, 3, 25, 12}
 pipeline := NewDataPipeline(data)

 ctx := context.Background()
 start := time.Now()

 filtered, transformed, analyzed, err := pipeline.ProcessParallel(ctx)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n处理结果:\n")
 fmt.Printf("- 过滤结果: %v\n", filtered)
 fmt.Printf("- 转换结果: %v\n", transformed)
 fmt.Printf("- 分析结果: %v\n", analyzed)
 fmt.Printf("\n总耗时: %v\n", time.Since(start))
}
```

#### 反例说明

**错误实现1：顺序执行并行任务**

```go
// ❌ 错误：将可并行的任务顺序执行
func WrongSequential() {
 sendEmail()        // 耗时 200ms
 updateStats()      // 耗时 300ms
 logOperation()     // 耗时 100ms
}
// 总耗时: 600ms
```

**问题**：任务之间没有依赖关系，顺序执行浪费时间。正确做法应并行执行，总耗时约300ms。

**错误实现2：goroutine泄漏**

```go
// ❌ 错误：未正确等待goroutine完成
func WrongGoroutineLeak() {
 go task1()
 go task2()
 // 函数返回，goroutine可能还在运行
}
```

**问题**：函数返回时goroutine可能仍在运行，导致不可预测的行为或资源泄漏。

**错误实现3：竞态条件**

```go
// ❌ 错误：多个goroutine同时修改共享变量
var counter int

func WrongRaceCondition() {
 var wg sync.WaitGroup
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter++  // 非原子操作，存在竞态
  }()
 }
 wg.Wait()
}
```

**问题**：存在竞态条件，counter的最终值不确定。应使用sync.Mutex或atomic操作。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 批量处理 | 同时处理多个独立的数据项 |
| 多通知渠道 | 同时发送邮件、短信、推送 |
| 数据并行计算 | 大数据集的分片并行处理 |
| 多服务调用 | 同时调用多个无依赖的服务 |
| 资源预加载 | 同时加载多个资源文件 |

#### 与其他模式关系

- **与同步模式**：并行分支通常需要同步模式来汇聚结果
- **与多实例模式**：多实例是并行分支的特例
- **与排他选择**：两者都是分支模式，但一个是并行一个是互斥

---

### 1.3 同步（Synchronization）

#### 概念定义

**同步模式**用于将多个并行执行路径汇聚成一个单一的执行路径。只有当所有入分支都完成后，后续任务才能开始执行。这是并行分支的配套模式。

**形式化定义**：给定并行任务集合 P = {p₁, p₂, ..., pₙ}，同步模式定义一个汇聚点 j，使得 j 只有在 ∀pᵢ ∈ P 都完成时才能被执行。

#### BPMN描述

在BPMN中，同步通过**并行网关（Parallel Gateway）**的汇聚功能实现：

```
[任务A] →┐
[任务B] →┼→ [◇] → [后续任务]
[任务C] →┘
```

- 并行网关等待所有入分支完成
- 只有全部到达后才触发后续任务
- 与分支网关配对使用

#### 流程图

```
┌─────────┐     ┌─────────┐
│  任务A  │────→│         │
└─────────┘     │         │
                │  同步   │────→┌─────────┐
┌─────────┐     │  网关   │     │  后续   │
│  任务B  │────→│         │     │  任务   │
└─────────┘     │         │     └─────────┘
                │         │
┌─────────┐     │         │
│  任务C  │────→│         │
└─────────┘     └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// Synchronizer 同步器
type Synchronizer struct {
 results map[string]interface{}
 mu      sync.RWMutex
}

// NewSynchronizer 创建同步器
func NewSynchronizer() *Synchronizer {
 return &Synchronizer{
  results: make(map[string]interface{}),
 }
}

// SetResult 设置任务结果
func (s *Synchronizer) SetResult(key string, value interface{}) {
 s.mu.Lock()
 defer s.mu.Unlock()
 s.results[key] = value
}

// GetResult 获取任务结果
func (s *Synchronizer) GetResult(key string) (interface{}, bool) {
 s.mu.RLock()
 defer s.mu.RUnlock()
 val, ok := s.results[key]
 return val, ok
}

// GetAllResults 获取所有结果
func (s *Synchronizer) GetAllResults() map[string]interface{} {
 s.mu.RLock()
 defer s.mu.RUnlock()

 // 返回副本
 result := make(map[string]interface{})
 for k, v := range s.results {
  result[k] = v
 }
 return result
}

// ParallelWithSync 并行执行并同步结果
func ParallelWithSync(
 ctx context.Context,
 tasks map[string]func(ctx context.Context) (interface{}, error),
) (map[string]interface{}, error) {
 var wg sync.WaitGroup
 syncer := NewSynchronizer()
 errChan := make(chan error, len(tasks))

 log.Printf("[同步] 启动 %d 个并行任务，等待全部完成", len(tasks))

 for name, task := range tasks {
  wg.Add(1)
  go func(taskName string, fn func(ctx context.Context) (interface{}, error)) {
   defer wg.Done()

   start := time.Now()
   result, err := fn(ctx)

   if err != nil {
    errChan <- fmt.Errorf("任务 '%s' 失败: %w", taskName, err)
    return
   }

   syncer.SetResult(taskName, result)
   log.Printf("[同步] 任务 '%s' 完成，耗时: %v", taskName, time.Since(start))
  }(name, task)
 }

 // 等待所有任务完成
 wg.Wait()
 close(errChan)

 // 检查错误
 for err := range errChan {
  return nil, err
 }

 log.Println("[同步] 所有任务已完成，继续执行")
 return syncer.GetAllResults(), nil
}

// 使用示例
func main() {
 tasks := map[string]func(ctx context.Context) (interface{}, error){
  "fetchUser": func(ctx context.Context) (interface{}, error) {
   fmt.Println("→ 获取用户信息...")
   time.Sleep(200 * time.Millisecond)
   return map[string]string{"name": "张三", "id": "001"}, nil
  },
  "fetchOrders": func(ctx context.Context) (interface{}, error) {
   fmt.Println("→ 获取订单列表...")
   time.Sleep(300 * time.Millisecond)
   return []string{"order1", "order2", "order3"}, nil
  },
  "fetchStats": func(ctx context.Context) (interface{}, error) {
   fmt.Println("→ 获取统计数据...")
   time.Sleep(150 * time.Millisecond)
   return map[string]int{"total": 100, "active": 80}, nil
  },
 }

 ctx := context.Background()
 start := time.Now()

 results, err := ParallelWithSync(ctx, tasks)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n同步结果:\n")
 for name, result := range results {
  fmt.Printf("- %s: %v\n", name, result)
 }
 fmt.Printf("\n总耗时: %v\n", time.Since(start))
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// ReportGenerator 报告生成器
type ReportGenerator struct {
 userID string
}

// NewReportGenerator 创建报告生成器
func NewReportGenerator(userID string) *ReportGenerator {
 return &ReportGenerator{userID: userID}
}

// Generate 生成完整报告
func (rg *ReportGenerator) Generate(ctx context.Context) (*Report, error) {
 // 定义并行数据收集任务
 dataTasks := map[string]func(ctx context.Context) (interface{}, error){
  "profile":    rg.fetchUserProfile,
  "activities": rg.fetchUserActivities,
  "payments":   rg.fetchPaymentHistory,
  "preferences": rg.fetchPreferences,
 }

 // 并行收集数据并同步
 results, err := ParallelWithSync(ctx, dataTasks)
 if err != nil {
  return nil, fmt.Errorf("数据收集失败: %w", err)
 }

 // 所有数据收集完成后，生成报告
 report := &Report{
  UserID:      rg.userID,
  Profile:     results["profile"],
  Activities:  results["activities"],
  Payments:    results["payments"],
  Preferences: results["preferences"],
  GeneratedAt: time.Now(),
 }

 return report, nil
}

func (rg *ReportGenerator) fetchUserProfile(ctx context.Context) (interface{}, error) {
 time.Sleep(100 * time.Millisecond)
 return map[string]string{"name": "张三", "level": "VIP"}, nil
}

func (rg *ReportGenerator) fetchUserActivities(ctx context.Context) (interface{}, error) {
 time.Sleep(200 * time.Millisecond)
 return []string{"登录", "浏览商品", "下单"}, nil
}

func (rg *ReportGenerator) fetchPaymentHistory(ctx context.Context) (interface{}, error) {
 time.Sleep(150 * time.Millisecond)
 return []map[string]interface{}{
  {"amount": 100, "date": "2024-01-01"},
  {"amount": 200, "date": "2024-01-15"},
 }, nil
}

func (rg *ReportGenerator) fetchPreferences(ctx context.Context) (interface{}, error) {
 time.Sleep(80 * time.Millisecond)
 return map[string][]string{
  "categories": {"电子产品", "图书"},
 }, nil
}

// Report 报告结构
type Report struct {
 UserID      string
 Profile     interface{}
 Activities  interface{}
 Payments    interface{}
 Preferences interface{}
 GeneratedAt time.Time
}

func main() {
 generator := NewReportGenerator("user_001")
 ctx := context.Background()

 start := time.Now()
 report, err := generator.Generate(ctx)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("报告生成完成!\n")
 fmt.Printf("用户: %s\n", report.UserID)
 fmt.Printf("生成时间: %v\n", report.GeneratedAt)
 fmt.Printf("总耗时: %v\n", time.Since(start))
}
```

#### 反例说明

**错误实现1：过早继续**

```go
// ❌ 错误：不等待所有任务完成就继续
func WrongEarlyContinue() {
 var wg sync.WaitGroup

 wg.Add(1)
 go func() {
  defer wg.Done()
  fetchData1()
 }()

 wg.Add(1)
 go func() {
  defer wg.Done()
  fetchData2()
 }()

 // 错误：只等待一个任务
 // 应该使用 wg.Wait()

 processResult()  // 可能在数据2还未准备好时执行
}
```

**问题**：后续处理可能在部分数据未准备好时执行，导致数据不完整或错误。

**错误实现2：死锁**

```go
// ❌ 错误：错误的WaitGroup使用导致死锁
func WrongDeadlock() {
 var wg sync.WaitGroup

 for i := 0; i < 3; i++ {
  wg.Wait()  // 错误：在Add之前调用Wait
  wg.Add(1)
  go func() {
   defer wg.Done()
   task()
  }()
 }
}
```

**问题**：Wait在Add之前调用会导致死锁。正确的顺序是先Add再启动goroutine，最后Wait。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 数据聚合 | 从多个源收集数据后统一处理 |
| 并行计算 | 分片计算后合并结果 |
| 多服务调用 | 调用多个服务后整合响应 |
| 资源加载 | 加载多个资源后开始渲染 |
| 批量处理 | 批量操作完成后提交事务 |

#### 与其他模式关系

- **与并行分支**：同步是并行分支的配套模式，两者常成对出现
- **与多选模式**：同步合并是多选的配套汇聚模式
- **与N选M模式**：N选M是同步的泛化形式

---

### 1.4 排他选择（Exclusive Choice）

#### 概念定义

**排他选择模式**允许工作流根据条件选择一条且仅一条执行路径。这是实现条件分支、决策逻辑的核心模式。

**形式化定义**：给定条件集合 C = {c₁, c₂, ..., cₙ} 和对应的任务集合 T = {t₁, t₂, ..., tₙ}，排他选择模式选择满足 cᵢ 的唯一任务 tᵢ 执行，且 ∀i≠j，cᵢ ∧ cⱼ = false。

#### BPMN描述

在BPMN中，排他选择通过**排他网关（Exclusive Gateway，XOR）**实现：

```
                    ┌→ [条件A] → [任务A] →┐
[前置任务] → [X] →┼→ [条件B] → [任务B] →┼→ [合并点]
                    └→ [条件C] → [任务C] →┘
```

- 排他网关（X）表示决策点
- 每个出分支有条件表达式
- 只有一条路径会被激活

#### 流程图

```
┌─────────┐     ┌─────────┐
│  开始   │────→│  输入   │
└─────────┘     └────┬────┘
                     │
                     ▼
               ┌─────────────┐
               │  排他网关   │
               │  (XOR)      │
               └──────┬──────┘
                      │
         ┌────────────┼────────────┐
         │            │            │
    [条件1]      [条件2]      [条件3]
         │            │            │
         ▼            ▼            ▼
   ┌─────────┐  ┌─────────┐  ┌─────────┐
   │  分支A  │  │  分支B  │  │  分支C  │
   └────┬────┘  └────┬────┘  └────┬────┘
        │            │            │
        └────────────┼────────────┘
                     │
                     ▼
               ┌─────────────┐
               │   合并点    │
               └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
)

// Condition 条件函数类型
type Condition func(ctx context.Context) bool

// ExclusiveChoice 排他选择器
type ExclusiveChoice struct {
 branches []struct {
  name      string
  condition Condition
  action    func(ctx context.Context) error
 }
 defaultAction func(ctx context.Context) error
}

// NewExclusiveChoice 创建排他选择器
func NewExclusiveChoice() *ExclusiveChoice {
 return &ExclusiveChoice{
  branches: make([]struct {
   name      string
   condition Condition
   action    func(ctx context.Context) error
  }, 0),
 }
}

// When 添加条件分支
func (e *ExclusiveChoice) When(name string, condition Condition, action func(ctx context.Context) error) *ExclusiveChoice {
 e.branches = append(e.branches, struct {
  name      string
  condition Condition
  action    func(ctx context.Context) error
 }{name, condition, action})
 return e
}

// Otherwise 设置默认分支
func (e *ExclusiveChoice) Otherwise(action func(ctx context.Context) error) *ExclusiveChoice {
 e.defaultAction = action
 return e
}

// Execute 执行排他选择
func (e *ExclusiveChoice) Execute(ctx context.Context) error {
 for _, branch := range e.branches {
  if branch.condition(ctx) {
   log.Printf("[排他选择] 条件满足: %s", branch.name)
   return branch.action(ctx)
  }
 }

 if e.defaultAction != nil {
  log.Println("[排他选择] 执行默认分支")
  return e.defaultAction(ctx)
 }

 return fmt.Errorf("没有条件满足且未设置默认分支")
}

// 使用示例
func main() {
 // 模拟订单金额
 orderAmount := 150.0

 choice := NewExclusiveChoice().
  When("小额订单", func(ctx context.Context) bool {
   return orderAmount < 100
  }, func(ctx context.Context) error {
   fmt.Println("→ 处理小额订单: 快速通道")
   return nil
  }).
  When("中额订单", func(ctx context.Context) bool {
   return orderAmount >= 100 && orderAmount < 500
  }, func(ctx context.Context) error {
   fmt.Println("→ 处理中额订单: 标准流程")
   return nil
  }).
  When("大额订单", func(ctx context.Context) bool {
   return orderAmount >= 500
  }, func(ctx context.Context) error {
   fmt.Println("→ 处理大额订单: 需要审批")
   return nil
  }).
  Otherwise(func(ctx context.Context) error {
   fmt.Println("→ 处理默认情况")
   return nil
  })

 ctx := context.Background()
 if err := choice.Execute(ctx); err != nil {
  log.Fatal(err)
 }
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
)

// PaymentRouter 支付路由
type PaymentRouter struct {
 choice *ExclusiveChoice
}

// NewPaymentRouter 创建支付路由器
func NewPaymentRouter() *PaymentRouter {
 return &PaymentRouter{
  choice: NewExclusiveChoice(),
 }
}

// Route 路由支付请求
func (pr *PaymentRouter) Route(ctx context.Context, method string, amount float64) error {
 router := NewExclusiveChoice().
  When("支付宝", func(ctx context.Context) bool {
   return method == "alipay"
  }, func(ctx context.Context) error {
   return pr.processAlipay(ctx, amount)
  }).
  When("微信支付", func(ctx context.Context) bool {
   return method == "wechat"
  }, func(ctx context.Context) error {
   return pr.processWechatPay(ctx, amount)
  }).
  When("银行卡", func(ctx context.Context) bool {
   return method == "card"
  }, func(ctx context.Context) error {
   return pr.processCardPay(ctx, amount)
  }).
  Otherwise(func(ctx context.Context) error {
   return fmt.Errorf("不支持的支付方式: %s", method)
  })

 return router.Execute(ctx)
}

func (pr *PaymentRouter) processAlipay(ctx context.Context, amount float64) error {
 fmt.Printf("→ 使用支付宝处理支付: ¥%.2f\n", amount)
 // 调用支付宝API
 return nil
}

func (pr *PaymentRouter) processWechatPay(ctx context.Context, amount float64) error {
 fmt.Printf("→ 使用微信支付处理支付: ¥%.2f\n", amount)
 // 调用微信支付API
 return nil
}

func (pr *PaymentRouter) processCardPay(ctx context.Context, amount float64) error {
 fmt.Printf("→ 使用银行卡处理支付: ¥%.2f\n", amount)
 // 调用银行API
 return nil
}

// 状态机风格实现
type StateMachine struct {
 currentState string
 transitions  map[string]map[string]string // state -> event -> newState
 actions      map[string]func(ctx context.Context) error
}

func NewStateMachine(initialState string) *StateMachine {
 return &StateMachine{
  currentState: initialState,
  transitions:  make(map[string]map[string]string),
  actions:      make(map[string]func(ctx context.Context) error),
 }
}

func (sm *StateMachine) AddTransition(from, event, to string) {
 if sm.transitions[from] == nil {
  sm.transitions[from] = make(map[string]string)
 }
 sm.transitions[from][event] = to
}

func (sm *StateMachine) OnState(state string, action func(ctx context.Context) error) {
 sm.actions[state] = action
}

func (sm *StateMachine) Trigger(ctx context.Context, event string) error {
 transitions := sm.transitions[sm.currentState]
 if transitions == nil {
  return fmt.Errorf("状态 '%s' 没有定义转换", sm.currentState)
 }

 newState, ok := transitions[event]
 if !ok {
  return fmt.Errorf("事件 '%s' 在状态 '%s' 下无效", event, sm.currentState)
 }

 log.Printf("[状态机] %s --%s--> %s", sm.currentState, event, newState)
 sm.currentState = newState

 if action, ok := sm.actions[sm.currentState]; ok {
  return action(ctx)
 }
 return nil
}

func main() {
 // 支付路由示例
 router := NewPaymentRouter()
 ctx := context.Background()

 fmt.Println("=== 支付路由示例 ===")
 if err := router.Route(ctx, "alipay", 199.99); err != nil {
  log.Fatal(err)
 }

 // 状态机示例
 fmt.Println("\n=== 状态机示例 ===")
 sm := NewStateMachine("待支付")
 sm.AddTransition("待支付", "支付", "已支付")
 sm.AddTransition("已支付", "发货", "已发货")
 sm.AddTransition("已发货", "签收", "已完成")

 sm.OnState("已支付", func(ctx context.Context) error {
  fmt.Println("→ 执行: 发送支付成功通知")
  return nil
 })

 sm.OnState("已发货", func(ctx context.Context) error {
  fmt.Println("→ 执行: 发送物流信息")
  return nil
 })

 sm.Trigger(ctx, "支付")
 sm.Trigger(ctx, "发货")
}
```

#### 反例说明

**错误实现1：多个条件同时满足**

```go
// ❌ 错误：条件不互斥
func WrongNonExclusive() {
 score := 85

 if score >= 60 {
  fmt.Println("及格")
 }
 if score >= 80 {  // 错误：85分同时满足两个条件
  fmt.Println("良好")
 }
}
```

**问题**：条件不互斥，导致多个分支被执行。应使用if-else链或确保条件互斥。

**错误实现2：条件覆盖不全**

```go
// ❌ 错误：缺少默认处理
func WrongNoDefault() {
 status := "unknown"

 switch status {
 case "pending":
  processPending()
 case "completed":
  processCompleted()
 }
 // 错误：unknown状态未被处理
}
```

**问题**：未处理所有可能的情况，导致某些输入被静默忽略。应添加default分支。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 支付路由 | 根据支付方式选择处理通道 |
| 审批流程 | 根据金额选择审批级别 |
| 错误处理 | 根据错误类型选择处理策略 |
| 功能开关 | 根据配置启用不同功能 |
| 用户分级 | 根据用户类型提供不同服务 |

#### 与其他模式关系

- **与简单合并**：排他选择通常与简单合并配对使用
- **与多选**：多选是排他选择的泛化，允许多个条件同时满足
- **与延迟选择**：延迟选择推迟决策点到运行时

---

### 1.5 简单合并（Simple Merge）

#### 概念定义

**简单合并模式**用于将多个互斥的执行路径汇聚成一个单一的执行路径。与同步不同，简单合并不需要等待所有入分支，任意一个入分支到达即可触发后续任务。

**形式化定义**：给定互斥任务集合 T = {t₁, t₂, ..., tₙ}，简单合并模式定义汇聚点 m，使得 m 在任意 tᵢ ∈ T 完成时即可执行，且 ∀i≠j，tᵢ 和 tⱼ 不会同时执行。

#### BPMN描述

在BPMN中，简单合并通过**排他网关（Exclusive Gateway）**的汇聚功能实现：

```
[任务A] →┐
[任务B] →┼→ [X] → [后续任务]
[任务C] →┘
```

- 排他网关作为汇聚点
- 任意入分支到达即触发后续
- 与排他选择网关配对使用

#### 流程图

```
┌─────────┐     ┌─────────┐
│  条件1  │────→│  任务A  │────┐
└─────────┘     └─────────┘    │
                               │
┌─────────┐     ┌─────────┐    │    ┌─────────┐
│  条件2  │────→│  任务B  │────┼───→│  后续   │
└─────────┘     └─────────┘    │    │  任务   │
                               │    └─────────┘
┌─────────┐     ┌─────────┐    │
│  条件3  │────→│  任务C  │────┘
└─────────┘     └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
)

// SimpleMerge 简单合并器
type SimpleMerge struct {
 branches []struct {
  name      string
  condition func(ctx context.Context) bool
  action    func(ctx context.Context) (interface{}, error)
 }
 mergeHandler func(ctx context.Context, result interface{}, branchName string) error
}

// NewSimpleMerge 创建简单合并器
func NewSimpleMerge() *SimpleMerge {
 return &SimpleMerge{
  branches: make([]struct {
   name      string
   condition func(ctx context.Context) bool
   action    func(ctx context.Context) (interface{}, error)
  }, 0),
 }
}

// AddBranch 添加分支
func (s *SimpleMerge) AddBranch(name string, condition func(ctx context.Context) bool, action func(ctx context.Context) (interface{}, error)) *SimpleMerge {
 s.branches = append(s.branches, struct {
  name      string
  condition func(ctx context.Context) bool
  action    func(ctx context.Context) (interface{}, error)
 }{name, condition, action})
 return s
}

// OnMerge 设置合并处理器
func (s *SimpleMerge) OnMerge(handler func(ctx context.Context, result interface{}, branchName string) error) *SimpleMerge {
 s.mergeHandler = handler
 return s
}

// Execute 执行并合并
func (s *SimpleMerge) Execute(ctx context.Context) error {
 // 找到第一个满足条件的分支
 for _, branch := range s.branches {
  if branch.condition(ctx) {
   log.Printf("[简单合并] 执行分支: %s", branch.name)

   result, err := branch.action(ctx)
   if err != nil {
    return fmt.Errorf("分支 '%s' 执行失败: %w", branch.name, err)
   }

   if s.mergeHandler != nil {
    return s.mergeHandler(ctx, result, branch.name)
   }
   return nil
  }
 }

 return fmt.Errorf("没有分支的条件被满足")
}

// 使用示例
func main() {
 // 根据用户类型选择不同的认证方式，然后统一处理
 userType := "vip"

 merge := NewSimpleMerge().
  AddBranch("VIP用户",
   func(ctx context.Context) bool { return userType == "vip" },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ VIP认证: 快速通道")
    return map[string]string{"level": "VIP", "discount": "20%"}, nil
   }).
  AddBranch("普通用户",
   func(ctx context.Context) bool { return userType == "normal" },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 普通认证: 标准流程")
    return map[string]string{"level": "Normal", "discount": "5%"}, nil
   }).
  AddBranch("游客",
   func(ctx context.Context) bool { return userType == "guest" },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 游客认证: 限制访问")
    return map[string]string{"level": "Guest", "discount": "0%"}, nil
   }).
  OnMerge(func(ctx context.Context, result interface{}, branchName string) error {
   fmt.Printf("→ 合并处理: 来自 '%s' 分支的结果 %v\n", branchName, result)
   return nil
  })

 ctx := context.Background()
 if err := merge.Execute(ctx); err != nil {
  log.Fatal(err)
 }
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
)

// OrderProcessor 订单处理器
type OrderProcessor struct{}

// Process 处理订单
func (op *OrderProcessor) Process(ctx context.Context, orderType string, amount float64) error {
 processor := NewSimpleMerge().
  AddBranch("线上订单",
   func(ctx context.Context) bool { return orderType == "online" },
   func(ctx context.Context) (interface{}, error) {
    return op.processOnlineOrder(ctx, amount)
   }).
  AddBranch("线下订单",
   func(ctx context.Context) bool { return orderType == "offline" },
   func(ctx context.Context) (interface{}, error) {
    return op.processOfflineOrder(ctx, amount)
   }).
  AddBranch("电话订单",
   func(ctx context.Context) bool { return orderType == "phone" },
   func(ctx context.Context) (interface{}, error) {
    return op.processPhoneOrder(ctx, amount)
   }).
  OnMerge(func(ctx context.Context, result interface{}, branchName string) error {
   return op.finalizeOrder(ctx, result, branchName)
  })

 return processor.Execute(ctx)
}

func (op *OrderProcessor) processOnlineOrder(ctx context.Context, amount float64) (interface{}, error) {
 fmt.Println("→ 处理线上订单...")
 return map[string]interface{}{
  "type":     "online",
  "amount":   amount,
  "channel":  "web",
  "discount": amount * 0.1,
 }, nil
}

func (op *OrderProcessor) processOfflineOrder(ctx context.Context, amount float64) (interface{}, error) {
 fmt.Println("→ 处理线下订单...")
 return map[string]interface{}{
  "type":     "offline",
  "amount":   amount,
  "channel":  "store",
  "discount": amount * 0.05,
 }, nil
}

func (op *OrderProcessor) processPhoneOrder(ctx context.Context, amount float64) (interface{}, error) {
 fmt.Println("→ 处理电话订单...")
 return map[string]interface{}{
  "type":     "phone",
  "amount":   amount,
  "channel":  "call",
  "discount": amount * 0.08,
 }, nil
}

func (op *OrderProcessor) finalizeOrder(ctx context.Context, result interface{}, branchName string) error {
 order := result.(map[string]interface{})
 finalAmount := order["amount"].(float64) - order["discount"].(float64)

 fmt.Printf("→ 订单合并处理:\n")
 fmt.Printf("  - 来源: %s\n", branchName)
 fmt.Printf("  - 类型: %s\n", order["type"])
 fmt.Printf("  - 原始金额: ¥%.2f\n", order["amount"])
 fmt.Printf("  - 优惠: ¥%.2f\n", order["discount"])
 fmt.Printf("  - 实付: ¥%.2f\n", finalAmount)
 return nil
}

func main() {
 processor := &OrderProcessor{}
 ctx := context.Background()

 fmt.Println("=== 线上订单 ===")
 processor.Process(ctx, "online", 1000)

 fmt.Println("\n=== 线下订单 ===")
 processor.Process(ctx, "offline", 500)
}
```

#### 反例说明

**错误实现1：使用同步代替简单合并**

```go
// ❌ 错误：使用WaitGroup等待所有分支
func WrongUseSync() {
 var wg sync.WaitGroup
 var result interface{}

 wg.Add(3)
 go func() { defer wg.Done(); if condition1() { result = task1() } }()
 go func() { defer wg.Done(); if condition2() { result = task2() } }()
 go func() { defer wg.Done(); if condition3() { result = task3() } }()

 wg.Wait()  // 错误：等待所有，即使只有一个会执行
 process(result)
}
```

**问题**：排他选择的分支是互斥的，使用同步会等待不必要的分支，浪费时间。

**错误实现2：重复执行**

```go
// ❌ 错误：多个条件同时满足时执行多个分支
func WrongMultipleExecution() {
 if condition1() {
  task1()
 }
 if condition2() {  // 错误：可能同时满足
  task2()
 }
 process()  // 可能被执行多次
}
```

**问题**：条件不互斥时会导致多个分支执行，后续处理可能被重复执行。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 订单处理 | 不同来源订单的不同处理后统一结算 |
| 认证流程 | 不同认证方式后的统一授权 |
| 支付方式 | 不同支付方式后的统一确认 |
| 渠道接入 | 不同渠道请求的统一响应 |
| 异常处理 | 不同异常类型的统一日志记录 |

#### 与其他模式关系

- **与排他选择**：简单合并是排他选择的配套汇聚模式
- **与同步合并**：同步合并是多选的配套汇聚模式，功能更强大
- **与多合并**：多合并允许多个分支到达，功能更复杂

---

## 2. 高级分支与同步模式

### 2.1 多选（Multi-Choice）

#### 概念定义

**多选模式**允许工作流根据多个条件选择一条或多条执行路径。与排他选择不同，多选允许多个条件同时满足，从而激活多个分支。

**形式化定义**：给定条件集合 C = {c₁, c₂, ..., cₙ} 和对应的任务集合 T = {t₁, t₂, ..., tₙ}，多选模式选择所有满足 cᵢ 的任务 tᵢ 执行，允许多个条件同时满足。

#### BPMN描述

在BPMN中，多选通过**包容网关（Inclusive Gateway，OR）**实现：

```
                    ┌→ [条件A] → [任务A] →┐
[前置任务] → [O] →┼→ [条件B] → [任务B] →┼→ [同步合并]
                    └→ [条件C] → [任务C] →┘
```

- 包容网关（O）表示多选决策点
- 每个出分支有独立的条件
- 多个条件可同时满足
- 需要同步合并汇聚

#### 流程图

```
┌─────────┐     ┌─────────────┐
│  开始   │────→│  包容网关   │
└─────────┘     │   (OR)      │
                └──────┬──────┘
                       │
         ┌─────────────┼─────────────┐
         │             │             │
    [条件A]       [条件B]       [条件C]
    (可满)        (可满)        (可满)
         │             │             │
         ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │  任务A  │   │  任务B  │   │  任务C  │
   └────┬────┘   └────┬────┘   └────┬────┘
        │             │             │
        └─────────────┼─────────────┘
                      │
                      ▼
               ┌─────────────┐
               │  同步合并   │
               └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
)

// MultiChoice 多选器
type MultiChoice struct {
 branches []struct {
  name      string
  condition func(ctx context.Context) bool
  action    func(ctx context.Context) (interface{}, error)
 }
}

// NewMultiChoice 创建多选器
func NewMultiChoice() *MultiChoice {
 return &MultiChoice{
  branches: make([]struct {
   name      string
   condition func(ctx context.Context) bool
   action    func(ctx context.Context) (interface{}, error)
  }, 0),
 }
}

// When 添加条件分支
func (m *MultiChoice) When(name string, condition func(ctx context.Context) bool, action func(ctx context.Context) (interface{}, error)) *MultiChoice {
 m.branches = append(m.branches, struct {
  name      string
  condition func(ctx context.Context) bool
  action    func(ctx context.Context) (interface{}, error)
 }{name, condition, action})
 return m
}

// Execute 执行多选
func (m *MultiChoice) Execute(ctx context.Context) (map[string]interface{}, error) {
 var wg sync.WaitGroup
 results := make(map[string]interface{})
 var mu sync.Mutex
 errChan := make(chan error, len(m.branches))

 activeCount := 0

 for _, branch := range m.branches {
  if branch.condition(ctx) {
   activeCount++
   wg.Add(1)
   go func(name string, action func(ctx context.Context) (interface{}, error)) {
    defer wg.Done()

    log.Printf("[多选] 执行分支: %s", name)
    result, err := action(ctx)

    if err != nil {
     errChan <- fmt.Errorf("分支 '%s' 失败: %w", name, err)
     return
    }

    mu.Lock()
    results[name] = result
    mu.Unlock()
   }(branch.name, branch.action)
  }
 }

 if activeCount == 0 {
  return nil, fmt.Errorf("没有条件被满足")
 }

 log.Printf("[多选] %d 个分支被激活", activeCount)

 wg.Wait()
 close(errChan)

 for err := range errChan {
  return nil, err
 }

 return results, nil
}

// 使用示例
func main() {
 // 用户通知场景：根据用户偏好选择多个通知渠道
 userPrefs := struct {
  Email    bool
  SMS      bool
  Push     bool
  Wechat   bool
 }{Email: true, SMS: true, Push: false, Wechat: true}

 notifier := NewMultiChoice().
  When("邮件通知",
   func(ctx context.Context) bool { return userPrefs.Email },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 发送邮件通知")
    return "email_sent", nil
   }).
  When("短信通知",
   func(ctx context.Context) bool { return userPrefs.SMS },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 发送短信通知")
    return "sms_sent", nil
   }).
  When("推送通知",
   func(ctx context.Context) bool { return userPrefs.Push },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 发送推送通知")
    return "push_sent", nil
   }).
  When("微信通知",
   func(ctx context.Context) bool { return userPrefs.Wechat },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 发送微信通知")
    return "wechat_sent", nil
   })

 ctx := context.Background()
 results, err := notifier.Execute(ctx)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n通知结果: %v\n", results)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
)

// DataProcessor 数据处理器
type DataProcessor struct {
 data []int
}

// NewDataProcessor 创建数据处理器
func NewDataProcessor(data []int) *DataProcessor {
 return &DataProcessor{data: data}
}

// Process 根据条件选择多个处理方式
func (dp *DataProcessor) Process(ctx context.Context, config ProcessingConfig) (*ProcessingResult, error) {
 processor := NewMultiChoice().
  When("过滤",
   func(ctx context.Context) bool { return config.EnableFilter },
   func(ctx context.Context) (interface{}, error) {
    return dp.filter(ctx)
   }).
  When("排序",
   func(ctx context.Context) bool { return config.EnableSort },
   func(ctx context.Context) (interface{}, error) {
    return dp.sort(ctx)
   }).
  When("去重",
   func(ctx context.Context) bool { return config.EnableDeduplicate },
   func(ctx context.Context) (interface{}, error) {
    return dp.deduplicate(ctx)
   }).
  When("聚合",
   func(ctx context.Context) bool { return config.EnableAggregate },
   func(ctx context.Context) (interface{}, error) {
    return dp.aggregate(ctx)
   })

 results, err := processor.Execute(ctx)
 if err != nil {
  return nil, err
 }

 return &ProcessingResult{
  Operations: results,
  InputSize:  len(dp.data),
 }, nil
}

func (dp *DataProcessor) filter(ctx context.Context) (interface{}, error) {
 var result []int
 for _, v := range dp.data {
  if v > 10 {
   result = append(result, v)
  }
 }
 fmt.Printf("→ 过滤: %d -> %d 条\n", len(dp.data), len(result))
 return result, nil
}

func (dp *DataProcessor) sort(ctx context.Context) (interface{}, error) {
 result := make([]int, len(dp.data))
 copy(result, dp.data)
 // 简单冒泡排序
 for i := 0; i < len(result); i++ {
  for j := i + 1; j < len(result); j++ {
   if result[i] > result[j] {
    result[i], result[j] = result[j], result[i]
   }
  }
 }
 fmt.Println("→ 排序完成")
 return result, nil
}

func (dp *DataProcessor) deduplicate(ctx context.Context) (interface{}, error) {
 seen := make(map[int]bool)
 var result []int
 for _, v := range dp.data {
  if !seen[v] {
   seen[v] = true
   result = append(result, v)
  }
 }
 fmt.Printf("→ 去重: %d -> %d 条\n", len(dp.data), len(result))
 return result, nil
}

func (dp *DataProcessor) aggregate(ctx context.Context) (interface{}, error) {
 sum := 0
 for _, v := range dp.data {
  sum += v
 }
 avg := sum / len(dp.data)
 fmt.Printf("→ 聚合: sum=%d, avg=%d\n", sum, avg)
 return map[string]int{"sum": sum, "avg": avg}, nil
}

// ProcessingConfig 处理配置
type ProcessingConfig struct {
 EnableFilter        bool
 EnableSort          bool
 EnableDeduplicate   bool
 EnableAggregate     bool
}

// ProcessingResult 处理结果
type ProcessingResult struct {
 Operations map[string]interface{}
 InputSize  int
}

func main() {
 data := []int{5, 15, 8, 20, 5, 15, 3, 25}
 processor := NewDataProcessor(data)
 ctx := context.Background()

 config := ProcessingConfig{
  EnableFilter:      true,
  EnableSort:        true,
  EnableDeduplicate: false,
  EnableAggregate:   true,
 }

 result, err := processor.Process(ctx, config)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n处理完成，输入: %d 条\n", result.InputSize)
 fmt.Println("各操作结果:")
 for op, res := range result.Operations {
  fmt.Printf("  - %s: %v\n", op, res)
 }
}
```

#### 反例说明

**错误实现1：使用排他选择实现多选**

```go
// ❌ 错误：使用if-else链，只能执行一个分支
func WrongUseExclusive() {
 if conditionA() {
  taskA()
 } else if conditionB() {  // 错误：A和B可能都需要执行
  taskB()
 } else if conditionC() {
  taskC()
 }
}
```

**问题**：使用if-else链强制互斥，即使多个条件都满足也只会执行一个分支。

**错误实现2：缺少同步**

```go
// ❌ 错误：多选后不进行同步
func WrongNoSynchronization() {
 if conditionA() {
  go taskA()  // 异步执行
 }
 if conditionB() {
  go taskB()  // 异步执行
 }
 processResult()  // 错误：可能在任务完成前就执行
}
```

**问题**：多选后的分支需要同步汇聚，否则后续处理可能在分支完成前执行。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 多渠道通知 | 根据用户偏好同时发送多种通知 |
| 数据处理 | 根据配置同时执行多个处理步骤 |
| 功能组合 | 根据权限同时启用多个功能 |
| 数据验证 | 同时执行多个独立的验证规则 |
| 日志记录 | 同时记录到多个目标 |

#### 与其他模式关系

- **与排他选择**：多选是排他选择的泛化，允许多个分支同时激活
- **与同步合并**：多选必须与同步合并配对使用
- **与并行分支**：多选是带条件的并行分支

---

### 2.2 同步合并（Synchronizing Merge）

#### 概念定义

**同步合并模式**是多选模式的配套汇聚模式。它等待所有被激活的分支完成后才触发后续任务，能够动态适应实际被激活的分支数量。

**形式化定义**：给定多选激活的分支集合 A ⊆ {1, 2, ..., n}，同步合并模式定义汇聚点 s，使得 s 只有在所有 i ∈ A 的分支都完成时才能执行。

#### BPMN描述

在BPMN中，同步合并通过**包容网关（Inclusive Gateway，OR）**的汇聚功能实现：

```
[任务A] →┐
[任务B] →┼→ [O] → [后续任务]
[任务C] →┘
```

- 包容网关作为汇聚点
- 动态等待实际激活的分支
- 不需要知道设计时有多少分支会被激活

#### 流程图

```
┌─────────┐     ┌─────────┐
│  任务A  │────→│         │
└─────────┘     │         │
                │  同步   │────→┌─────────┐
┌─────────┐     │  合并   │     │  后续   │
│  任务B  │────→│  (OR)   │     │  任务   │
└─────────┘     │         │     └─────────┘
                │         │
┌─────────┐     │         │
│  任务C  │────→│         │
└─────────┘     └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "sync/atomic"
)

// SynchronizingMerge 同步合并器
type SynchronizingMerge struct {
 activeCount int32
 wg          sync.WaitGroup
 results     map[string]interface{}
 mu          sync.Mutex
 errChan     chan error
}

// NewSynchronizingMerge 创建同步合并器
func NewSynchronizingMerge() *SynchronizingMerge {
 return &SynchronizingMerge{
  results: make(map[string]interface{}),
  errChan: make(chan error, 100),
 }
}

// Activate 激活一个分支
func (s *SynchronizingMerge) Activate(name string, fn func(ctx context.Context) (interface{}, error)) {
 atomic.AddInt32(&s.activeCount, 1)
 s.wg.Add(1)

 go func() {
  defer s.wg.Done()

  result, err := fn(context.Background())
  if err != nil {
   s.errChan <- fmt.Errorf("分支 '%s' 失败: %w", name, err)
   return
  }

  s.mu.Lock()
  s.results[name] = result
  s.mu.Unlock()

  log.Printf("[同步合并] 分支 '%s' 完成", name)
 }()
}

// Wait 等待所有激活的分支完成
func (s *SynchronizingMerge) Wait() (map[string]interface{}, error) {
 active := atomic.LoadInt32(&s.activeCount)
 log.Printf("[同步合并] 等待 %d 个分支完成", active)

 s.wg.Wait()
 close(s.errChan)

 for err := range s.errChan {
  return nil, err
 }

 log.Println("[同步合并] 所有分支已完成")
 return s.results, nil
}

// 使用示例
func main() {
 merge := NewSynchronizingMerge()

 // 根据条件动态激活分支
 conditions := map[string]bool{
  "branchA": true,
  "branchB": true,
  "branchC": false,
  "branchD": true,
 }

 for name, active := range conditions {
  if active {
   merge.Activate(name, func(ctx context.Context) (interface{}, error) {
    fmt.Printf("→ 执行 %s\n", name)
    return fmt.Sprintf("result_%s", name), nil
   })
  }
 }

 results, err := merge.Wait()
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n合并结果: %v\n", results)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// WorkflowEngine 工作流引擎
type WorkflowEngine struct {
 merge *SynchronizingMerge
}

// ExecuteDynamicWorkflow 执行动态工作流
func (we *WorkflowEngine) ExecuteDynamicWorkflow(ctx context.Context, tasks map[string]bool) (map[string]interface{}, error) {
 merge := NewSynchronizingMerge()

 // 根据运行时条件动态激活任务
 for taskName, shouldExecute := range tasks {
  if shouldExecute {
   task := we.createTask(taskName)
   merge.Activate(taskName, task)
  }
 }

 return merge.Wait()
}

func (we *WorkflowEngine) createTask(name string) func(ctx context.Context) (interface{}, error) {
 tasks := map[string]func(ctx context.Context) (interface{}, error){
  "validate": func(ctx context.Context) (interface{}, error) {
   time.Sleep(100 * time.Millisecond)
   return "validation_passed", nil
  },
  "enrich": func(ctx context.Context) (interface{}, error) {
   time.Sleep(200 * time.Millisecond)
   return "data_enriched", nil
  },
  "transform": func(ctx context.Context) (interface{}, error) {
   time.Sleep(150 * time.Millisecond)
   return "data_transformed", nil
  },
  "audit": func(ctx context.Context) (interface{}, error) {
   time.Sleep(50 * time.Millisecond)
   return "audit_logged", nil
  },
 }

 if task, ok := tasks[name]; ok {
  return task
 }
 return func(ctx context.Context) (interface{}, error) {
  return fmt.Sprintf("task_%s_done", name), nil
 }
}

// MultiChoiceWithSync 多选+同步合并组合
type MultiChoiceWithSync struct {
 branches []struct {
  name      string
  condition func(ctx context.Context) bool
  action    func(ctx context.Context) (interface{}, error)
 }
}

func NewMultiChoiceWithSync() *MultiChoiceWithSync {
 return &MultiChoiceWithSync{
  branches: make([]struct {
   name      string
   condition func(ctx context.Context) bool
   action    func(ctx context.Context) (interface{}, error)
  }, 0),
 }
}

func (m *MultiChoiceWithSync) When(name string, condition func(ctx context.Context) bool, action func(ctx context.Context) (interface{}, error)) *MultiChoiceWithSync {
 m.branches = append(m.branches, struct {
  name      string
  condition func(ctx context.Context) bool
  action    func(ctx context.Context) (interface{}, error)
 }{name, condition, action})
 return m
}

func (m *MultiChoiceWithSync) Execute(ctx context.Context) (map[string]interface{}, error) {
 merge := NewSynchronizingMerge()

 for _, branch := range m.branches {
  if branch.condition(ctx) {
   merge.Activate(branch.name, branch.action)
  }
 }

 return merge.Wait()
}

func main() {
 // 场景：根据配置动态选择数据处理步骤
 config := struct {
  Validate    bool
  Enrich      bool
  Transform   bool
  Audit       bool
 }{Validate: true, Enrich: true, Transform: false, Audit: true}

 workflow := NewMultiChoiceWithSync().
  When("validate",
   func(ctx context.Context) bool { return config.Validate },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 执行数据验证")
    time.Sleep(100 * time.Millisecond)
    return "validated", nil
   }).
  When("enrich",
   func(ctx context.Context) bool { return config.Enrich },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 执行数据增强")
    time.Sleep(200 * time.Millisecond)
    return "enriched", nil
   }).
  When("transform",
   func(ctx context.Context) bool { return config.Transform },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 执行数据转换")
    time.Sleep(150 * time.Millisecond)
    return "transformed", nil
   }).
  When("audit",
   func(ctx context.Context) bool { return config.Audit },
   func(ctx context.Context) (interface{}, error) {
    fmt.Println("→ 执行审计日志")
    time.Sleep(50 * time.Millisecond)
    return "audited", nil
   })

 ctx := context.Background()
 start := time.Now()

 results, err := workflow.Execute(ctx)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n工作流完成，耗时: %v\n", time.Since(start))
 fmt.Println("执行结果:")
 for name, result := range results {
  fmt.Printf("  - %s: %v\n", name, result)
 }
}
```

#### 反例说明

**错误实现1：固定等待所有分支**

```go
// ❌ 错误：使用固定数量的WaitGroup
func WrongFixedWait() {
 var wg sync.WaitGroup
 wg.Add(4)  // 错误：假设所有4个分支都会执行

 if conditionA() {
  go func() { defer wg.Done(); taskA() }()
 }
 if conditionB() {
  go func() { defer wg.Done(); taskB() }()
 }
 // 如果只有2个条件满足，Wait会永远阻塞

 wg.Wait()
}
```

**问题**：WaitGroup计数与实际执行的分支数不匹配，可能导致死锁或panic。

**错误实现2：使用普通同步代替同步合并**

```go
// ❌ 错误：使用并行网关等待所有分支
func WrongUseParallelSync() {
 // 假设有4个分支
 var wg sync.WaitGroup
 wg.Add(4)

 go func() { defer wg.Done(); if conditionA() { taskA() } }()
 go func() { defer wg.Done(); if conditionB() { taskB() } }()
 go func() { defer wg.Done(); if conditionC() { taskC() } }()
 go func() { defer wg.Done(); if conditionD() { taskD() } }()

 wg.Wait()  // 错误：等待所有4个，即使有些条件不满足
}
```

**问题**：普通同步等待所有分支，而同步合并只等待被激活的分支。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 动态工作流 | 运行时决定执行哪些任务 |
| 条件处理 | 根据条件执行不同组合的任务 |
| 插件系统 | 动态加载和执行插件 |
| 数据处理 | 根据配置选择处理步骤 |
| 审批流程 | 根据条件选择审批人 |

#### 与其他模式关系

- **与多选**：同步合并是多选的配套汇聚模式
- **与同步**：同步合并是同步的泛化，支持动态分支数
- **与鉴别器**：鉴别器只等待部分分支，同步合并等待所有激活分支

---

### 2.3 多合并（Multi-Merge）

#### 概念定义

**多合并模式**允许每个入分支到达时都触发后续任务执行。与同步合并不同，多合并不等待所有分支，而是每个分支独立触发后续。

**形式化定义**：给定分支集合 B = {b₁, b₂, ..., bₙ}，多合并模式定义触发点 t，使得每个 bᵢ ∈ B 到达时都会独立触发后续任务执行。

#### BPMN描述

在BPMN中，多合并没有专门的网关符号，通常通过多个顺序流直接连接到同一个活动实现：

```
[任务A] →┐
[任务B] →┼→ [后续任务]
[任务C] →┘
```

- 多个入分支直接连接到同一活动
- 每个入分支到达都触发活动执行
- 活动可能被多次执行

#### 流程图

```
┌─────────┐     ┌─────────┐
│  任务A  │────→│         │
└─────────┘     │         │────→┌─────────┐
                │  多合并 │     │  后续   │
┌─────────┐     │  (触发) │     │  任务   │
│  任务B  │────→│         │────→│ (可多次)│
└─────────┘     │         │     └─────────┘
                │         │────→
┌─────────┐     │         │
│  任务C  │────→│         │
└─────────┘     └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
)

// MultiMerge 多合并器
type MultiMerge struct {
 handler func(ctx context.Context, source string, data interface{}) error
}

// NewMultiMerge 创建多合并器
func NewMultiMerge(handler func(ctx context.Context, source string, data interface{}) error) *MultiMerge {
 return &MultiMerge{handler: handler}
}

// Branch 创建分支处理器
func (m *MultiMerge) Branch(name string) func(data interface{}) {
 return func(data interface{}) {
  ctx := context.Background()
  log.Printf("[多合并] 分支 '%s' 到达", name)

  if err := m.handler(ctx, name, data); err != nil {
   log.Printf("[多合并] 处理分支 '%s' 失败: %v", name, err)
  }
 }
}

// 使用示例
func main() {
 executionCount := 0
 var mu sync.Mutex

 merge := NewMultiMerge(func(ctx context.Context, source string, data interface{}) error {
  mu.Lock()
  executionCount++
  count := executionCount
  mu.Unlock()

  fmt.Printf("→ 第%d次执行: 来自 '%s' 的数据: %v\n", count, source, data)
  return nil
 })

 // 模拟多个分支到达
 var wg sync.WaitGroup

 wg.Add(1)
 go func() {
  defer wg.Done()
  merge.Branch("A")("data_from_A")
 }()

 wg.Add(1)
 go func() {
  defer wg.Done()
  merge.Branch("B")("data_from_B")
 }()

 wg.Add(1)
 go func() {
  defer wg.Done()
  merge.Branch("C")("data_from_C")
 }()

 wg.Wait()

 fmt.Printf("\n后续任务共执行了 %d 次\n", executionCount)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// EventProcessor 事件处理器
type EventProcessor struct {
 merge *MultiMerge
}

// NewEventProcessor 创建事件处理器
func NewEventProcessor() *EventProcessor {
 return &EventProcessor{}
}

// ProcessEvents 处理多个来源的事件
func (ep *EventProcessor) ProcessEvents(ctx context.Context, events []Event) error {
 var processedCount int32
 var mu sync.Mutex
 results := make([]string, 0)

 merge := NewMultiMerge(func(ctx context.Context, source string, data interface{}) error {
  event := data.(Event)

  // 处理事件
  result := ep.handleEvent(event)

  mu.Lock()
  processedCount++
  results = append(results, fmt.Sprintf("%s: %s", source, result))
  mu.Unlock()

  return nil
 })

 var wg sync.WaitGroup

 // 模拟多个事件源
 for i, event := range events {
  wg.Add(1)
  go func(index int, evt Event) {
   defer wg.Done()
   time.Sleep(time.Duration(index*50) * time.Millisecond)
   merge.Branch(fmt.Sprintf("source_%d", index))(evt)
  }(i, event)
 }

 wg.Wait()

 fmt.Printf("\n处理完成，共处理 %d 个事件\n", len(results))
 for _, r := range results {
  fmt.Printf("  - %s\n", r)
 }

 return nil
}

func (ep *EventProcessor) handleEvent(event Event) string {
 time.Sleep(50 * time.Millisecond)
 return fmt.Sprintf("processed_%s", event.Type)
}

// Event 事件结构
type Event struct {
 Type string
 Data map[string]interface{}
}

// NotificationService 通知服务
type NotificationService struct {
 notifications []Notification
}

// NewNotificationService 创建通知服务
func NewNotificationService() *NotificationService {
 return &NotificationService{
  notifications: make([]Notification, 0),
 }
}

// SendNotifications 发送多渠道通知
func (ns *NotificationService) SendNotifications(ctx context.Context, message string, channels []string) error {
 var mu sync.Mutex

 merge := NewMultiMerge(func(ctx context.Context, channel string, data interface{}) error {
  msg := data.(string)

  // 发送通知
  notif := Notification{
   Channel:   channel,
   Message:   msg,
   Timestamp: time.Now(),
  }

  mu.Lock()
  ns.notifications = append(ns.notifications, notif)
  mu.Unlock()

  fmt.Printf("→ [%s] 发送: %s\n", channel, msg)
  return nil
 })

 var wg sync.WaitGroup

 for _, channel := range channels {
  wg.Add(1)
  go func(ch string) {
   defer wg.Done()
   merge.Branch(ch)(message)
  }(channel)
 }

 wg.Wait()

 fmt.Printf("\n共发送 %d 条通知\n", len(ns.notifications))
 return nil
}

// Notification 通知结构
type Notification struct {
 Channel   string
 Message   string
 Timestamp time.Time
}

func main() {
 // 事件处理示例
 fmt.Println("=== 事件处理示例 ===")
 events := []Event{
  {Type: "click", Data: map[string]interface{}{"x": 100, "y": 200}},
  {Type: "scroll", Data: map[string]interface{}{"delta": 50}},
  {Type: "submit", Data: map[string]interface{}{"form": "contact"}},
 }

 processor := NewEventProcessor()
 if err := processor.ProcessEvents(context.Background(), events); err != nil {
  log.Fatal(err)
 }

 // 通知服务示例
 fmt.Println("\n=== 通知服务示例 ===")
 notifService := NewNotificationService()
 channels := []string{"email", "sms", "push", "wechat"}

 if err := notifService.SendNotifications(context.Background(), "您的订单已发货", channels); err != nil {
  log.Fatal(err)
 }
}
```

#### 反例说明

**错误实现1：使用同步合并代替多合并**

```go
// ❌ 错误：等待所有分支完成后再执行
func WrongUseSyncMerge() {
 var wg sync.WaitGroup
 results := make(chan interface{}, 3)

 wg.Add(3)
 go func() { defer wg.Done(); results <- taskA() }()
 go func() { defer wg.Done(); results <- taskB() }()
 go func() { defer wg.Done(); results <- taskC() }()

 wg.Wait()
 close(results)

 for r := range results {  // 错误：只执行一次
  process(r)
 }
}
```

**问题**：多合并要求每个分支到达都触发后续，而同步合并只触发一次。

**错误实现2：竞态条件**

```go
// ❌ 错误：多个goroutine同时修改共享状态
var counter int

func WrongRaceCondition() {
 merge := NewMultiMerge(func(ctx context.Context, source string, data interface{}) error {
  counter++  // 非原子操作
  process(data)
  return nil
 })

 // 多个分支同时触发
 go merge.Branch("A")(data1)
 go merge.Branch("B")(data2)
}
```

**问题**：存在竞态条件，需要使用互斥锁保护共享状态。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 事件处理 | 多个事件源的事件独立处理 |
| 通知发送 | 每个通知渠道独立发送 |
| 日志记录 | 每个操作独立记录日志 |
| 数据同步 | 多个数据源独立同步 |
| 监控告警 | 每个监控项独立告警 |

#### 与其他模式关系

- **与同步合并**：多合并不等待，每个分支独立触发；同步合并等待所有分支
- **与简单合并**：简单合并是排他选择的汇聚，多合并是多选的汇聚变体
- **与鉴别器**：鉴别器只触发一次后续，多合并触发多次

---

### 2.4 鉴别器（Discriminator）

#### 概念定义

**鉴别器模式**等待多个并行分支中的第一个完成，然后立即触发后续任务，同时取消或忽略其他分支的结果。这是一种"先到先得"的汇聚模式。

**形式化定义**：给定并行任务集合 P = {p₁, p₂, ..., pₙ}，鉴别器模式定义汇聚点 d，使得 d 在第一个 pᵢ 完成时触发，其他任务的结果被忽略。

#### BPMN描述

在BPMN中，鉴别器通过**复杂网关（Complex Gateway）**或特定的鉴别器语义实现：

```
[任务A] →┐
[任务B] →┼→ [鉴别器] → [后续任务]
[任务C] →┘
```

- 等待第一个完成的任务
- 触发后续任务
- 其他任务结果被忽略

#### 流程图

```
┌─────────┐     ┌─────────┐
│  任务A  │────→│         │
└─────────┘     │         │
                │  鉴别器 │────→┌─────────┐
┌─────────┐     │ (第一个 │     │  后续   │
│  任务B  │────→│  完成)  │     │  任务   │
└─────────┘     │         │     └─────────┘
                │         │
┌─────────┐     │         │
│  任务C  │────→│         │
└─────────┘     └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "sync/atomic"
 "time"
)

// Discriminator 鉴别器
type Discriminator struct {
 completed int32
 result    interface{}
 mu        sync.Mutex
}

// NewDiscriminator 创建鉴别器
func NewDiscriminator() *Discriminator {
 return &Discriminator{}
}

// Race 竞争执行，返回第一个成功的结果
func (d *Discriminator) Race(ctx context.Context, tasks map[string]func(ctx context.Context) (interface{}, error)) (string, interface{}, error) {
 ctx, cancel := context.WithCancel(ctx)
 defer cancel()

 resultChan := make(chan struct {
  name   string
  result interface{}
  err    error
 }, 1)

 var wg sync.WaitGroup

 for name, task := range tasks {
  wg.Add(1)
  go func(n string, t func(ctx context.Context) (interface{}, error)) {
   defer wg.Done()

   start := time.Now()
   result, err := t(ctx)

   if err != nil {
    log.Printf("[鉴别器] 任务 '%s' 失败: %v", n, err)
    return
   }

   // 尝试成为第一个完成的
   if atomic.CompareAndSwapInt32(&d.completed, 0, 1) {
    log.Printf("[鉴别器] 任务 '%s' 第一个完成，耗时: %v", n, time.Since(start))
    cancel() // 取消其他任务

    select {
    case resultChan <- struct {
     name   string
     result interface{}
     err    error
    }{n, result, nil}:
    default:
    }
   } else {
    log.Printf("[鉴别器] 任务 '%s' 完成但被忽略", n)
   }
  }(name, task)
 }

 // 等待所有任务完成或取消
 go func() {
  wg.Wait()
  close(resultChan)
 }()

 select {
 case r := <-resultChan:
  if r.err != nil {
   return "", nil, r.err
  }
  return r.name, r.result, nil
 case <-ctx.Done():
  return "", nil, ctx.Err()
 }
}

// 使用示例
func main() {
 discriminator := NewDiscriminator()

 // 模拟多个服务提供商，选择最快响应的
 providers := map[string]func(ctx context.Context) (interface{}, error){
  "providerA": func(ctx context.Context) (interface{}, error) {
   time.Sleep(200 * time.Millisecond)
   return "result_from_A", nil
  },
  "providerB": func(ctx context.Context) (interface{}, error) {
   time.Sleep(100 * time.Millisecond)
   return "result_from_B", nil
  },
  "providerC": func(ctx context.Context) (interface{}, error) {
   time.Sleep(300 * time.Millisecond)
   return "result_from_C", nil
  },
 }

 ctx := context.Background()
 start := time.Now()

 winner, result, err := discriminator.Race(ctx, providers)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n获胜者: %s\n", winner)
 fmt.Printf("结果: %v\n", result)
 fmt.Printf("总耗时: %v\n", time.Since(start))
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// FastestProvider 最快提供者选择器
type FastestProvider struct {
 discriminator *Discriminator
}

// NewFastestProvider 创建最快提供者选择器
func NewFastestProvider() *FastestProvider {
 return &FastestProvider{
  discriminator: NewDiscriminator(),
 }
}

// Query 查询最快响应的提供者
func (fp *FastestProvider) Query(ctx context.Context, query string, providers []Provider) (*QueryResult, error) {
 tasks := make(map[string]func(ctx context.Context) (interface{}, error))

 for _, provider := range providers {
  p := provider // 捕获循环变量
  tasks[p.Name] = func(ctx context.Context) (interface{}, error) {
   return p.Query(ctx, query)
  }
 }

 winner, result, err := fp.discriminator.Race(ctx, tasks)
 if err != nil {
  return nil, err
 }

 return &QueryResult{
  Provider: winner,
  Data:     result,
 }, nil
}

// Provider 提供者接口
type Provider struct {
 Name    string
 Latency time.Duration
}

// Query 查询
func (p Provider) Query(ctx context.Context, query string) (interface{}, error) {
 select {
 case <-time.After(p.Latency):
  return fmt.Sprintf("%s_response_for_%s", p.Name, query), nil
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

// QueryResult 查询结果
type QueryResult struct {
 Provider string
 Data     interface{}
}

// CacheSelector 缓存选择器
type CacheSelector struct {
 discriminator *Discriminator
}

// NewCacheSelector 创建缓存选择器
func NewCacheSelector() *CacheSelector {
 return &CacheSelector{
  discriminator: NewDiscriminator(),
 }
}

// Get 从多个缓存中获取数据
func (cs *CacheSelector) Get(ctx context.Context, key string, caches []Cache) (interface{}, error) {
 tasks := make(map[string]func(ctx context.Context) (interface{}, error))

 for _, cache := range caches {
  c := cache
  tasks[c.Name] = func(ctx context.Context) (interface{}, error) {
   return c.Get(ctx, key)
  }
 }

 _, result, err := cs.discriminator.Race(ctx, tasks)
 return result, err
}

// Cache 缓存接口
type Cache struct {
 Name    string
 Latency time.Duration
 Data    map[string]interface{}
}

// Get 获取数据
func (c Cache) Get(ctx context.Context, key string) (interface{}, error) {
 select {
 case <-time.After(c.Latency):
  if val, ok := c.Data[key]; ok {
   return val, nil
  }
  return nil, fmt.Errorf("key not found")
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

func main() {
 // 最快提供者示例
 fmt.Println("=== 最快提供者选择 ===")
 fp := NewFastestProvider()

 providers := []Provider{
  {Name: "AWS", Latency: 150 * time.Millisecond},
  {Name: "Azure", Latency: 100 * time.Millisecond},
  {Name: "GCP", Latency: 200 * time.Millisecond},
 }

 ctx := context.Background()
 start := time.Now()

 result, err := fp.Query(ctx, "user_data", providers)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("最快响应: %s\n", result.Provider)
 fmt.Printf("结果: %v\n", result.Data)
 fmt.Printf("耗时: %v\n", time.Since(start))

 // 缓存选择示例
 fmt.Println("\n=== 缓存选择 ===")
 cs := NewCacheSelector()

 caches := []Cache{
  {
   Name:    "Redis",
   Latency: 50 * time.Millisecond,
   Data:    map[string]interface{}{"key1": "redis_value"},
  },
  {
   Name:    "Memcached",
   Latency: 30 * time.Millisecond,
   Data:    map[string]interface{}{"key1": "memcached_value"},
  },
  {
   Name:    "Local",
   Latency: 10 * time.Millisecond,
   Data:    map[string]interface{}{"key1": "local_value"},
  },
 }

 start = time.Now()
 val, err := cs.Get(ctx, "key1", caches)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("缓存值: %v\n", val)
 fmt.Printf("耗时: %v\n", time.Since(start))
}
```

#### 反例说明

**错误实现1：使用同步代替鉴别器**

```go
// ❌ 错误：等待所有任务完成
func WrongUseSync() {
 var wg sync.WaitGroup
 results := make([]interface{}, 0)
 var mu sync.Mutex

 for _, task := range tasks {
  wg.Add(1)
  go func(t Task) {
   defer wg.Done()
   result := t.Execute()
   mu.Lock()
   results = append(results, result)
   mu.Unlock()
  }(task)
 }

 wg.Wait()
 useResult(results[0])  // 错误：浪费了等待其他任务的时间
}
```

**问题**：等待所有任务完成浪费时间，鉴别器只需要第一个结果。

**错误实现2：没有取消其他任务**

```go
// ❌ 错误：获得第一个结果后不取消其他任务
func WrongNoCancel() {
 resultChan := make(chan interface{}, len(tasks))

 for _, task := range tasks {
  go func(t Task) {
   resultChan <- t.Execute()  // 即使第一个完成，其他还会继续执行
  }(task)
 }

 result := <-resultChan
 useResult(result)
 // 其他任务仍在浪费资源
}
```

**问题**：没有取消机制，其他任务会继续执行浪费资源。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 多服务查询 | 查询多个服务，使用最快响应 |
| 缓存读取 | 从多个缓存层读取，使用最快命中 |
| 服务发现 | 从多个注册中心发现服务 |
| DNS解析 | 查询多个DNS服务器 |
| 负载均衡 | 选择最快响应的后端 |

#### 与其他模式关系

- **与N选M**：鉴别器是N选M的特例（1-out-of-N）
- **与同步**：鉴别器只等待一个，同步等待所有
- **与并行分支**：鉴别器是并行分支的汇聚变体

---

### 2.5 N选M（N-out-of-M Join）

#### 概念定义

**N选M模式**是同步模式的泛化形式，它等待M个并行分支中的N个完成后才触发后续任务。这是实现部分汇聚的灵活模式。

**形式化定义**：给定并行任务集合 P = {p₁, p₂, ..., pₘ}，N选M模式定义汇聚点 j，使得 j 在任意 N 个任务完成时触发（N ≤ M）。

#### BPMN描述

在BPMN中，N选M通过**复杂网关（Complex Gateway）**实现：

```
[任务A] →┐
[任务B] →┼→ [复杂网关 N/M] → [后续任务]
[任务C] →┘
```

- 复杂网关配置N/M参数
- 等待N个分支完成
- 可配置是否取消剩余分支

#### 流程图

```
┌─────────┐     ┌─────────┐
│  任务A  │────→│         │
└─────────┘     │         │
                │  N选M   │
┌─────────┐     │  网关   │────→┌─────────┐
│  任务B  │────→│ (2/3)   │     │  后续   │
└─────────┘     │         │     │  任务   │
                │         │     └─────────┘
┌─────────┐     │         │
│  任务C  │────→│         │
└─────────┘     └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "sync/atomic"
 "time"
)

// NOutOfM N选M汇聚器
type NOutOfM struct {
 n         int
 m         int
 completed int32
 results   map[string]interface{}
 mu        sync.Mutex
 done      chan struct{}
 once      sync.Once
}

// NewNOutOfM 创建N选M汇聚器
func NewNOutOfM(n, m int) *NOutOfM {
 if n > m {
  n = m
 }
 return &NOutOfM{
  n:       n,
  m:       m,
  results: make(map[string]interface{}),
  done:    make(chan struct{}),
 }
}

// Submit 提交任务结果
func (n *NOutOfM) Submit(name string, result interface{}) bool {
 n.mu.Lock()
 n.results[name] = result
 n.mu.Unlock()

 completed := atomic.AddInt32(&n.completed, 1)
 log.Printf("[N选M] 任务 '%s' 完成，进度: %d/%d", name, completed, n.n)

 if int(completed) >= n.n {
  n.once.Do(func() {
   close(n.done)
  })
  return true
 }
 return false
}

// Wait 等待N个任务完成
func (n *NOutOfM) Wait() map[string]interface{} {
 <-n.done

 n.mu.Lock()
 defer n.mu.Unlock()

 // 返回副本
 result := make(map[string]interface{})
 for k, v := range n.results {
  result[k] = v
 }
 return result
}

// Execute 执行N选M任务
func ExecuteNOutOfM(ctx context.Context, n int, tasks map[string]func(ctx context.Context) (interface{}, error)) (map[string]interface{}, error) {
 m := len(tasks)
 join := NewNOutOfM(n, m)

 ctx, cancel := context.WithCancel(ctx)
 defer cancel()

 var wg sync.WaitGroup

 for name, task := range tasks {
  wg.Add(1)
  go func(n string, t func(ctx context.Context) (interface{}, error)) {
   defer wg.Done()

   result, err := t(ctx)
   if err != nil {
    log.Printf("[N选M] 任务 '%s' 失败: %v", n, err)
    return
   }

   if join.Submit(n, result) {
    log.Printf("[N选M] 已达到 %d 个任务，取消其他任务", n)
    cancel()
   }
  }(name, task)
 }

 results := join.Wait()
 wg.Wait()

 return results, nil
}

// 使用示例
func main() {
 // 模拟数据备份到多个存储，只需要2个成功即可
 backends := map[string]func(ctx context.Context) (interface{}, error){
  "s3": func(ctx context.Context) (interface{}, error) {
   time.Sleep(100 * time.Millisecond)
   return "s3_backup_ok", nil
  },
  "gcs": func(ctx context.Context) (interface{}, error) {
   time.Sleep(150 * time.Millisecond)
   return "gcs_backup_ok", nil
  },
  "azure": func(ctx context.Context) (interface{}, error) {
   time.Sleep(200 * time.Millisecond)
   return "azure_backup_ok", nil
  },
  "minio": func(ctx context.Context) (interface{}, error) {
   time.Sleep(80 * time.Millisecond)
   return "minio_backup_ok", nil
  },
 }

 ctx := context.Background()
 start := time.Now()

 // 只需要2个成功
 results, err := ExecuteNOutOfM(ctx, 2, backends)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n备份结果 (2/4):\n")
 for name, result := range results {
  fmt.Printf("  - %s: %v\n", name, result)
 }
 fmt.Printf("耗时: %v\n", time.Since(start))
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// DistributedQuery 分布式查询
type DistributedQuery struct{}

// QueryWithQuorum 使用法定人数的查询
func (dq *DistributedQuery) QueryWithQuorum(ctx context.Context, query string, nodes []Node, quorum int) (*QueryResponse, error) {
 tasks := make(map[string]func(ctx context.Context) (interface{}, error))

 for _, node := range nodes {
  n := node
  tasks[n.ID] = func(ctx context.Context) (interface{}, error) {
   return n.Query(ctx, query)
  }
 }

 results, err := ExecuteNOutOfM(ctx, quorum, tasks)
 if err != nil {
  return nil, err
 }

 // 合并结果
 return &QueryResponse{
  Query:      query,
  Results:    results,
  NodeCount:  len(results),
  TotalNodes: len(nodes),
 }, nil
}

// Node 节点
type Node struct {
 ID      string
 Latency time.Duration
 Data    map[string]string
}

// Query 查询
func (n Node) Query(ctx context.Context, query string) (interface{}, error) {
 select {
 case <-time.After(n.Latency):
  if val, ok := n.Data[query]; ok {
   return val, nil
  }
  return nil, fmt.Errorf("not found")
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

// QueryResponse 查询响应
type QueryResponse struct {
 Query      string
 Results    map[string]interface{}
 NodeCount  int
 TotalNodes int
}

// RedundantService 冗余服务
type RedundantService struct {
 instances []ServiceInstance
}

// CallWithRedundancy 带冗余的服务调用
func (rs *RedundantService) CallWithRedundancy(ctx context.Context, request interface{}, minSuccess int) (*ServiceResponse, error) {
 tasks := make(map[string]func(ctx context.Context) (interface{}, error))

 for i, inst := range rs.instances {
  instance := inst
  tasks[fmt.Sprintf("instance_%d", i)] = func(ctx context.Context) (interface{}, error) {
   return instance.Call(ctx, request)
  }
 }

 results, err := ExecuteNOutOfM(ctx, minSuccess, tasks)
 if err != nil {
  return nil, err
 }

 return &ServiceResponse{
  Data:         results,
  SuccessCount: len(results),
 }, nil
}

// ServiceInstance 服务实例
type ServiceInstance struct {
 Name    string
 Latency time.Duration
}

// Call 调用
func (si ServiceInstance) Call(ctx context.Context, request interface{}) (interface{}, error) {
 select {
 case <-time.After(si.Latency):
  return fmt.Sprintf("response_from_%s", si.Name), nil
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

// ServiceResponse 服务响应
type ServiceResponse struct {
 Data         map[string]interface{}
 SuccessCount int
}

func main() {
 // 分布式查询示例
 fmt.Println("=== 分布式查询 (法定人数: 2/4) ===")
 dq := &DistributedQuery{}

 nodes := []Node{
  {ID: "node1", Latency: 50 * time.Millisecond, Data: map[string]string{"user:1": "Alice"}},
  {ID: "node2", Latency: 80 * time.Millisecond, Data: map[string]string{"user:1": "Alice"}},
  {ID: "node3", Latency: 120 * time.Millisecond, Data: map[string]string{"user:1": "Alice"}},
  {ID: "node4", Latency: 30 * time.Millisecond, Data: map[string]string{"user:1": "Alice"}},
 }

 ctx := context.Background()
 start := time.Now()

 response, err := dq.QueryWithQuorum(ctx, "user:1", nodes, 2)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("查询: %s\n", response.Query)
 fmt.Printf("响应节点数: %d/%d\n", response.NodeCount, response.TotalNodes)
 fmt.Printf("结果: %v\n", response.Results)
 fmt.Printf("耗时: %v\n", time.Since(start))

 // 冗余服务示例
 fmt.Println("\n=== 冗余服务调用 (最少成功: 2/3) ===")
 rs := &RedundantService{
  instances: []ServiceInstance{
   {Name: "instance1", Latency: 100 * time.Millisecond},
   {Name: "instance2", Latency: 150 * time.Millisecond},
   {Name: "instance3", Latency: 80 * time.Millisecond},
  },
 }

 start = time.Now()
 svcResponse, err := rs.CallWithRedundancy(ctx, "request_data", 2)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("成功调用数: %d\n", svcResponse.SuccessCount)
 fmt.Printf("结果: %v\n", svcResponse.Data)
 fmt.Printf("耗时: %v\n", time.Since(start))
}
```

#### 反例说明

**错误实现1：N大于M**

```go
// ❌ 错误：N大于可用任务数
func WrongNLargerThanM() {
 tasks := []Task{task1, task2}  // 只有2个任务
 join := NewNOutOfM(3, 2)  // 错误：需要3个成功，但只有2个任务
 // 永远不会触发
}
```

**问题**：N大于M会导致永远无法触发，应该确保N ≤ M。

**错误实现2：没有取消机制**

```go
// ❌ 错误：达到N后继续执行其他任务
func WrongNoCancel() {
 var completed int32

 for _, task := range tasks {
  go func(t Task) {
   result := t.Execute()
   if atomic.AddInt32(&completed, 1) <= int32(n) {
    results = append(results, result)
   }
   // 错误：即使达到N，其他任务仍继续执行
  }(task)
 }
}
```

**问题**：没有取消机制，其他任务会继续浪费资源。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 分布式共识 | 等待法定人数（quorum）的响应 |
| 冗余备份 | 数据备份到多个位置，只需要部分成功 |
| 多服务调用 | 调用多个服务，只需要部分响应 |
| 数据同步 | 从多个节点同步，只需要部分确认 |
| 故障转移 | 多个备用服务，只需要一个可用 |

#### 与其他模式关系

- **与同步**：同步是N选M的特例（M-out-of-M）
- **与鉴别器**：鉴别器是N选M的特例（1-out-of-M）
- **与并行分支**：N选M是并行分支的灵活汇聚模式

---

## 3. 结构化模式

### 3.1 任意循环（Arbitrary Cycles）

#### 概念定义

**任意循环模式**允许工作流中的任务通过任意路径回到之前的任务点，形成循环结构。这是实现迭代处理、重试机制的核心模式。

**形式化定义**：给定任务集合 T，任意循环模式允许存在边 (tⱼ, tᵢ) 其中 j > i，使得执行序列可以回到之前的任务。

#### BPMN描述

在BPMN中，任意循环通过**顺序流回环**实现：

```
         ┌─────────────────────────┐
         │                         │
         ▼                         │
[任务A] → [任务B] → [任务C] → [条件] ─┘
                      │
                      ▼
                   [任务D]
```

- 使用顺序流创建回环
- 条件决定是继续还是循环
- 需要防止无限循环

#### 流程图

```
        ┌───────────────────┐
        │                   │
        ▼                   │
┌─────────┐    ┌─────────┐  │
│  开始   │───→│  任务A  │  │
└─────────┘    └────┬────┘  │
                    │       │
                    ▼       │
               ┌─────────┐  │
               │  任务B  │  │
               └────┬────┘  │
                    │       │
                    ▼       │
               ┌─────────┐  │
               │  条件   │──┘
               └────┬────┘
                    │
           ┌────────┴────────┐
           │                 │
      [满足]              [不满足]
           │                 │
           ▼                 ▼
      ┌─────────┐       ┌─────────┐
      │  任务C  │       │  结束   │
      └─────────┘       └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// Loop 循环控制器
type Loop struct {
 maxIterations int
}

// NewLoop 创建循环控制器
func NewLoop(maxIterations int) *Loop {
 return &Loop{maxIterations: maxIterations}
}

// While while循环
func (l *Loop) While(ctx context.Context, condition func() bool, action func() error) error {
 iteration := 0

 for condition() {
  if iteration >= l.maxIterations {
   return fmt.Errorf("超过最大迭代次数: %d", l.maxIterations)
  }

  iteration++
  log.Printf("[循环] 第 %d 次迭代", iteration)

  if err := action(); err != nil {
   return fmt.Errorf("第 %d 次迭代失败: %w", iteration, err)
  }

  // 检查上下文取消
  select {
  case <-ctx.Done():
   return ctx.Err()
  default:
  }
 }

 log.Printf("[循环] 完成，共 %d 次迭代", iteration)
 return nil
}

// DoWhile do-while循环
func (l *Loop) DoWhile(ctx context.Context, action func() error, condition func() bool) error {
 iteration := 0

 for {
  if iteration >= l.maxIterations {
   return fmt.Errorf("超过最大迭代次数: %d", l.maxIterations)
  }

  iteration++
  log.Printf("[循环] 第 %d 次迭代", iteration)

  if err := action(); err != nil {
   return fmt.Errorf("第 %d 次迭代失败: %w", iteration, err)
  }

  if !condition() {
   break
  }

  select {
  case <-ctx.Done():
   return ctx.Err()
  default:
  }
 }

 log.Printf("[循环] 完成，共 %d 次迭代", iteration)
 return nil
}

// ForEach 遍历循环
func (l *Loop) ForEach(ctx context.Context, items []interface{}, action func(item interface{}, index int) error) error {
 for i, item := range items {
  if i >= l.maxIterations {
   return fmt.Errorf("超过最大迭代次数: %d", l.maxIterations)
  }

  log.Printf("[循环] 处理第 %d/%d 项", i+1, len(items))

  if err := action(item, i); err != nil {
   return fmt.Errorf("处理第 %d 项失败: %w", i, err)
  }

  select {
  case <-ctx.Done():
   return ctx.Err()
  default:
  }
 }

 return nil
}

// 使用示例
func main() {
 loop := NewLoop(10)
 ctx := context.Background()

 // While循环示例：重试机制
 fmt.Println("=== While循环：重试机制 ===")
 attempt := 0
 err := loop.While(ctx, func() bool {
  attempt++
  return attempt <= 3  // 最多重试3次
 }, func() error {
  fmt.Printf("→ 尝试第 %d 次\n", attempt)
  time.Sleep(100 * time.Millisecond)
  if attempt < 3 {
   return fmt.Errorf("暂时失败")
  }
  return nil
 })

 if err != nil {
  log.Printf("重试失败: %v", err)
 } else {
  fmt.Println("✓ 重试成功")
 }

 // DoWhile循环示例
 fmt.Println("\n=== DoWhile循环 ===")
 count := 0
 err = loop.DoWhile(ctx, func() error {
  count++
  fmt.Printf("→ 执行操作，count=%d\n", count)
  return nil
 }, func() bool {
  return count < 5
 })

 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("✓ 最终count=%d\n", count)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "math/rand"
 "time"
)

// RetryPolicy 重试策略
type RetryPolicy struct {
 MaxAttempts     int
 InitialInterval time.Duration
 MaxInterval     time.Duration
 Multiplier      float64
}

// DefaultRetryPolicy 默认重试策略
func DefaultRetryPolicy() *RetryPolicy {
 return &RetryPolicy{
  MaxAttempts:     5,
  InitialInterval: 100 * time.Millisecond,
  MaxInterval:     5 * time.Second,
  Multiplier:      2.0,
 }
}

// RetryWithBackoff 带退避的重试
func RetryWithBackoff(ctx context.Context, policy *RetryPolicy, operation func() error) error {
 interval := policy.InitialInterval

 for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
  log.Printf("[重试] 第 %d/%d 次尝试", attempt, policy.MaxAttempts)

  err := operation()
  if err == nil {
   log.Printf("[重试] 第 %d 次尝试成功", attempt)
   return nil
  }

  if attempt == policy.MaxAttempts {
   return fmt.Errorf("所有 %d 次尝试都失败: %w", policy.MaxAttempts, err)
  }

  // 计算下次重试间隔
  sleepDuration := interval + time.Duration(rand.Int63n(int64(interval)))
  log.Printf("[重试] 等待 %v 后重试...", sleepDuration)

  select {
  case <-time.After(sleepDuration):
   // 继续重试
  case <-ctx.Done():
   return ctx.Err()
  }

  // 增加间隔
  interval = time.Duration(float64(interval) * policy.Multiplier)
  if interval > policy.MaxInterval {
   interval = policy.MaxInterval
  }
 }

 return nil
}

// WorkflowLoop 工作流循环
type WorkflowLoop struct {
 loop *Loop
}

// NewWorkflowLoop 创建工作流循环
func NewWorkflowLoop() *WorkflowLoop {
 return &WorkflowLoop{
  loop: NewLoop(100),
 }
}

// ProcessWithRetry 带重试的处理
func (wl *WorkflowLoop) ProcessWithRetry(ctx context.Context, processor func() error) error {
 return wl.loop.While(ctx, func() bool {
  // 这里可以根据实际条件决定是否继续
  return true
 }, func() error {
  err := processor()
  if err != nil {
   // 重试
   return nil
  }
  // 成功，退出循环
  return fmt.Errorf("success")
 })
}

// PollUntilComplete 轮询直到完成
func PollUntilComplete(ctx context.Context, check func() (bool, error), interval time.Duration) error {
 ticker := time.NewTicker(interval)
 defer ticker.Stop()

 for {
  select {
  case <-ctx.Done():
   return ctx.Err()
  case <-ticker.C:
   completed, err := check()
   if err != nil {
    return err
   }
   if completed {
    return nil
   }
   log.Println("[轮询] 任务未完成，继续等待...")
  }
 }
}

// Iterator 迭代器模式
type Iterator struct {
 items []interface{}
 index int
}

// NewIterator 创建迭代器
func NewIterator(items []interface{}) *Iterator {
 return &Iterator{items: items, index: 0}
}

// HasNext 是否有下一个
func (it *Iterator) HasNext() bool {
 return it.index < len(it.items)
}

// Next 获取下一个
func (it *Iterator) Next() interface{} {
 if !it.HasNext() {
  return nil
 }
 item := it.items[it.index]
 it.index++
 return item
}

// ForEach 遍历
func (it *Iterator) ForEach(fn func(item interface{}) error) error {
 for it.HasNext() {
  if err := fn(it.Next()); err != nil {
   return err
  }
 }
 return nil
}

func main() {
 ctx := context.Background()

 // 重试示例
 fmt.Println("=== 带退避的重试 ===")
 attempt := 0
 err := RetryWithBackoff(ctx, DefaultRetryPolicy(), func() error {
  attempt++
  fmt.Printf("→ 执行操作 (attempt %d)\n", attempt)
  if attempt < 3 {
   return fmt.Errorf("暂时失败")
  }
  return nil
 })

 if err != nil {
  log.Fatal(err)
 }
 fmt.Println("✓ 操作成功")

 // 轮询示例
 fmt.Println("\n=== 轮询示例 ===")
 pollCount := 0
 ctx2, cancel := context.WithTimeout(ctx, 2*time.Second)
 defer cancel()

 err = PollUntilComplete(ctx2, func() (bool, error) {
  pollCount++
  fmt.Printf("→ 轮询检查 #%d\n", pollCount)
  return pollCount >= 3, nil
 }, 500*time.Millisecond)

 if err != nil {
  log.Printf("轮询结束: %v", err)
 } else {
  fmt.Println("✓ 任务完成")
 }

 // 迭代器示例
 fmt.Println("\n=== 迭代器示例 ===")
 items := []interface{}{"item1", "item2", "item3", "item4"}
 iterator := NewIterator(items)

 for iterator.HasNext() {
  item := iterator.Next()
  fmt.Printf("→ 处理: %v\n", item)
 }
}
```

#### 反例说明

**错误实现1：无限循环**

```go
// ❌ 错误：没有退出条件的循环
func WrongInfiniteLoop() {
 for {
  process()  // 错误：永远不会退出
 }
}
```

**问题**：没有退出条件的循环会导致程序挂起。应该设置最大迭代次数或退出条件。

**错误实现2：没有上下文检查**

```go
// ❌ 错误：循环中不检查上下文取消
func WrongNoContextCheck(ctx context.Context) {
 for condition() {
  process()  // 错误：即使上下文取消也会继续
 }
}
```

**问题**：不检查上下文取消会导致即使请求被取消，循环仍继续执行。

**错误实现3：竞态条件**

```go
// ❌ 错误：多个goroutine同时修改循环变量
var counter int

func WrongRaceCondition() {
 for i := 0; i < 10; i++ {
  go func() {
   counter++  // 竞态条件
  }()
 }
}
```

**问题**：存在竞态条件，需要使用同步机制。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 重试机制 | 操作失败后的自动重试 |
| 轮询检查 | 定期检查任务状态 |
| 批处理 | 循环处理批量数据 |
| 迭代处理 | 遍历集合并处理每个元素 |
| 条件等待 | 等待某个条件满足 |

#### 与其他模式关系

- **与多实例模式**：循环可以创建多实例
- **与递归模式**：循环和递归可以互相转换
- **与状态模式**：循环可以实现状态转换

---

### 3.2 隐式终止（Implicit Termination）

#### 概念定义

**隐式终止模式**指工作流在没有显式结束事件的情况下自然终止。当没有更多的任务可以执行时，工作流自动结束。这是最简单的工作流终止方式。

**形式化定义**：给定工作流 W，隐式终止发生在状态 S 下没有可执行的出转移时。

#### BPMN描述

在BPMN中，隐式终止通过**没有出顺序流的活动**实现：

```
[任务A] → [任务B] → [任务C]
                           (无出流，隐式终止)
```

- 活动没有出顺序流
- 工作流自然结束
- 不需要显式的结束事件

#### 流程图

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  开始   │────→│  任务A  │────→│  任务B  │
└─────────┘     └─────────┘     └────┬────┘
                                     │
                                     ▼
                                ┌─────────┐
                                │  任务C  │
                                └─────────┘
                                (隐式终止)
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
)

// ImplicitTermination 隐式终止执行器
type ImplicitTermination struct {
 tasks []func(ctx context.Context) error
}

// NewImplicitTermination 创建隐式终止执行器
func NewImplicitTermination() *ImplicitTermination {
 return &ImplicitTermination{
  tasks: make([]func(ctx context.Context) error, 0),
 }
}

// Add 添加任务
func (it *ImplicitTermination) Add(task func(ctx context.Context) error) {
 it.tasks = append(it.tasks, task)
}

// Execute 执行任务，自然终止
func (it *ImplicitTermination) Execute(ctx context.Context) error {
 for i, task := range it.tasks {
  select {
  case <-ctx.Done():
   return ctx.Err()
  default:
  }

  log.Printf("[隐式终止] 执行任务 %d/%d", i+1, len(it.tasks))

  if err := task(ctx); err != nil {
   return fmt.Errorf("任务 %d 失败: %w", i+1, err)
  }
 }

 // 所有任务完成，隐式终止
 log.Println("[隐式终止] 所有任务完成，工作流自然终止")
 return nil
}

// Workflow 工作流
type Workflow struct {
 steps []Step
}

// Step 工作流步骤
type Step struct {
 Name string
 Fn   func(ctx context.Context) error
}

// NewWorkflow 创建工作流
func NewWorkflow() *Workflow {
 return &Workflow{
  steps: make([]Step, 0),
 }
}

// AddStep 添加步骤
func (w *Workflow) AddStep(name string, fn func(ctx context.Context) error) {
 w.steps = append(w.steps, Step{Name: name, Fn: fn})
}

// Run 运行工作流
func (w *Workflow) Run(ctx context.Context) error {
 for _, step := range w.steps {
  select {
  case <-ctx.Done():
   return ctx.Err()
  default:
  }

  log.Printf("[工作流] 执行步骤: %s", step.Name)

  if err := step.Fn(ctx); err != nil {
   return fmt.Errorf("步骤 '%s' 失败: %w", step.Name, err)
  }
 }

 log.Println("[工作流] 执行完成")
 return nil
}

// 使用示例
func main() {
 workflow := NewWorkflow()

 // 添加工作流步骤
 workflow.AddStep("初始化", func(ctx context.Context) error {
  fmt.Println("→ 初始化系统...")
  return nil
 })

 workflow.AddStep("加载配置", func(ctx context.Context) error {
  fmt.Println("→ 加载配置...")
  return nil
 })

 workflow.AddStep("处理数据", func(ctx context.Context) error {
  fmt.Println("→ 处理数据...")
  return nil
 })

 workflow.AddStep("清理资源", func(ctx context.Context) error {
  fmt.Println("→ 清理资源...")
  return nil
 })

 // 隐式终止 - 没有显式的结束步骤
 ctx := context.Background()
 if err := workflow.Run(ctx); err != nil {
  log.Fatal(err)
 }

 fmt.Println("\n✓ 工作流自然终止")
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "os"
 "os/signal"
 "syscall"
)

// GracefulShutdown 优雅关闭
type GracefulShutdown struct {
 cleanupFuncs []func() error
}

// NewGracefulShutdown 创建优雅关闭管理器
func NewGracefulShutdown() *GracefulShutdown {
 return &GracefulShutdown{
  cleanupFuncs: make([]func() error, 0),
 }
}

// OnShutdown 注册关闭回调
func (gs *GracefulShutdown) OnShutdown(fn func() error) {
 gs.cleanupFuncs = append(gs.cleanupFuncs, fn)
}

// Wait 等待关闭信号
func (gs *GracefulShutdown) Wait() {
 sigChan := make(chan os.Signal, 1)
 signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

 log.Println("[优雅关闭] 等待关闭信号...")
 <-sigChan

 log.Println("[优雅关闭] 收到关闭信号，执行清理...")

 for i, fn := range gs.cleanupFuncs {
  if err := fn(); err != nil {
   log.Printf("[优雅关闭] 清理函数 %d 失败: %v", i+1, err)
  }
 }

 log.Println("[优雅关闭] 清理完成，程序终止")
}

// Pipeline 管道处理
type Pipeline struct {
 stages []Stage
}

// Stage 管道阶段
type Stage struct {
 Name string
 Fn   func(ctx context.Context, input interface{}) (interface{}, error)
}

// NewPipeline 创建管道
func NewPipeline() *Pipeline {
 return &Pipeline{
  stages: make([]Stage, 0),
 }
}

// AddStage 添加阶段
func (p *Pipeline) AddStage(name string, fn func(ctx context.Context, input interface{}) (interface{}, error)) {
 p.stages = append(p.stages, Stage{Name: name, Fn: fn})
}

// Process 处理数据
func (p *Pipeline) Process(ctx context.Context, input interface{}) (interface{}, error) {
 result := input

 for _, stage := range p.stages {
  select {
  case <-ctx.Done():
   return nil, ctx.Err()
  default:
  }

  log.Printf("[管道] 执行阶段: %s", stage.Name)

  var err error
  result, err = stage.Fn(ctx, result)
  if err != nil {
   return nil, fmt.Errorf("阶段 '%s' 失败: %w", stage.Name, err)
  }
 }

 return result, nil
}

// TaskChain 任务链
type TaskChain struct {
 tasks []Task
}

// Task 任务
type Task struct {
 Name string
 Fn   func() error
}

// NewTaskChain 创建任务链
func NewTaskChain() *TaskChain {
 return &TaskChain{
  tasks: make([]Task, 0),
 }
}

// Add 添加任务
func (tc *TaskChain) Add(name string, fn func() error) {
 tc.tasks = append(tc.tasks, Task{Name: name, Fn: fn})
}

// Execute 执行任务链
func (tc *TaskChain) Execute() error {
 for _, task := range tc.tasks {
  log.Printf("[任务链] 执行任务: %s", task.Name)

  if err := task.Fn(); err != nil {
   return fmt.Errorf("任务 '%s' 失败: %w", task.Name, err)
  }
 }

 return nil
}

// AutoClose 自动关闭资源
type AutoClose struct {
 resources []func()
}

// NewAutoClose 创建自动关闭管理器
func NewAutoClose() *AutoClose {
 return &AutoClose{
  resources: make([]func(), 0),
 }
}

// Register 注册关闭函数
func (ac *AutoClose) Register(fn func()) {
 ac.resources = append(ac.resources, fn)
}

// Close 关闭所有资源
func (ac *AutoClose) Close() {
 // 逆序关闭
 for i := len(ac.resources) - 1; i >= 0; i-- {
  ac.resources[i]()
 }
}

func main() {
 // 管道示例
 fmt.Println("=== 管道处理示例 ===")
 pipeline := NewPipeline()

 pipeline.AddStage("解析", func(ctx context.Context, input interface{}) (interface{}, error) {
  data := input.(string)
  fmt.Printf("→ 解析: %s\n", data)
  return map[string]string{"parsed": data}, nil
 })

 pipeline.AddStage("验证", func(ctx context.Context, input interface{}) (interface{}, error) {
  data := input.(map[string]string)
  fmt.Printf("→ 验证: %v\n", data)
  return data, nil
 })

 pipeline.AddStage("处理", func(ctx context.Context, input interface{}) (interface{}, error) {
  data := input.(map[string]string)
  fmt.Printf("→ 处理: %v\n", data)
  return "processed_result", nil
 })

 ctx := context.Background()
 result, err := pipeline.Process(ctx, "input_data")
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("✓ 管道处理完成: %v\n", result)

 // 任务链示例
 fmt.Println("\n=== 任务链示例 ===")
 taskChain := NewTaskChain()

 step := 0
 taskChain.Add("步骤1", func() error {
  step++
  fmt.Printf("→ 执行步骤 %d\n", step)
  return nil
 })
 taskChain.Add("步骤2", func() error {
  step++
  fmt.Printf("→ 执行步骤 %d\n", step)
  return nil
 })
 taskChain.Add("步骤3", func() error {
  step++
  fmt.Printf("→ 执行步骤 %d\n", step)
  return nil
 })

 if err := taskChain.Execute(); err != nil {
  log.Fatal(err)
 }

 fmt.Printf("✓ 任务链完成，共 %d 步\n", step)
}
```

#### 反例说明

**错误实现1：资源泄漏**

```go
// ❌ 错误：不关闭资源
func WrongResourceLeak() {
 file, _ := os.Open("data.txt")
 process(file)
 // 错误：没有关闭文件
}
```

**问题**：资源没有正确关闭，导致资源泄漏。应该使用defer或显式关闭。

**错误实现2：没有错误处理**

```go
// ❌ 错误：忽略错误
func WrongIgnoreErrors() {
 task1()  // 可能失败
 task2()  // 即使task1失败也会执行
 task3()  // 即使task2失败也会执行
}
```

**问题**：错误被忽略，可能导致无效操作继续执行。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 简单工作流 | 线性执行的工作流 |
| 批处理任务 | 批量数据处理 |
| 管道处理 | 数据流经多个处理阶段 |
| 资源管理 | 确保资源正确释放 |
| 优雅关闭 | 程序退出前的清理工作 |

#### 与其他模式关系

- **与显式终止**：隐式终止是显式终止的简化形式
- **与取消模式**：隐式终止不考虑取消，取消模式处理主动终止
- **与子流程**：子流程通常使用隐式终止

---

## 4. 多实例模式

### 4.1 多实例无同步（Multiple Instances without Synchronization）

#### 概念定义

**多实例无同步模式**允许在运行时创建多个活动实例，这些实例独立执行，不需要同步。这是实现批量并行处理的基础模式。

**形式化定义**：给定活动 A 和实例数 n，多实例无同步模式创建 n 个独立的 A 实例，每个实例独立执行，不等待其他实例。

#### BPMN描述

在BPMN中，多实例无同步通过**多实例活动（Multi-Instance Activity）**的并行设置实现：

```
[任务] * (并行，无同步)
```

- 多实例标记（三条竖线）
- 并行执行设置
- 无汇聚点

#### 流程图

```
┌─────────┐     ┌─────────────┐
│  数据   │────→│  多实例活动 │
│  集合   │     │  (无同步)   │
└─────────┘     └──────┬──────┘
                       │
         ┌─────────────┼─────────────┐
         │             │             │
         ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │ 实例1   │   │ 实例2   │   │ 实例3   │
   │ (独立)  │   │ (独立)  │   │ (独立)  │
   └─────────┘   └─────────┘   └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
)

// MultiInstanceNoSync 多实例无同步执行器
type MultiInstanceNoSync struct{}

// Execute 执行多实例，无同步
func (m *MultiInstanceNoSync) Execute(ctx context.Context, items []interface{}, handler func(ctx context.Context, item interface{}) error) error {
 var wg sync.WaitGroup
 errChan := make(chan error, len(items))

 log.Printf("[多实例无同步] 启动 %d 个实例", len(items))

 for i, item := range items {
  wg.Add(1)
  go func(index int, it interface{}) {
   defer wg.Done()

   log.Printf("[多实例无同步] 实例 %d 开始", index+1)

   if err := handler(ctx, it); err != nil {
    errChan <- fmt.Errorf("实例 %d 失败: %w", index+1, err)
    return
   }

   log.Printf("[多实例无同步] 实例 %d 完成", index+1)
  }(i, item)
 }

 // 启动goroutine等待所有实例完成
 go func() {
  wg.Wait()
  close(errChan)
 }()

 // 收集错误（非阻塞）
 var errs []error
 for err := range errChan {
  errs = append(errs, err)
 }

 if len(errs) > 0 {
  return fmt.Errorf("部分实例失败: %v", errs)
 }

 return nil
}

// 使用示例
func main() {
 mi := &MultiInstanceNoSync{}
 ctx := context.Background()

 // 批量处理订单
 orders := []interface{}{
  map[string]string{"id": "001", "amount": "100"},
  map[string]string{"id": "002", "amount": "200"},
  map[string]string{"id": "003", "amount": "150"},
  map[string]string{"id": "004", "amount": "300"},
 }

 err := mi.Execute(ctx, orders, func(ctx context.Context, item interface{}) error {
  order := item.(map[string]string)
  fmt.Printf("→ 处理订单 %s, 金额: %s\n", order["id"], order["amount"])
  return nil
 })

 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("\n✓ 所有实例处理完成")
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync/atomic"
 "time"
)

// BatchProcessor 批处理器
type BatchProcessor struct {
 mi *MultiInstanceNoSync
}

// NewBatchProcessor 创建批处理器
func NewBatchProcessor() *BatchProcessor {
 return &BatchProcessor{
  mi: &MultiInstanceNoSync{},
 }
}

// ProcessBatch 批量处理
func (bp *BatchProcessor) ProcessBatch(ctx context.Context, items []interface{}) (*BatchResult, error) {
 var successCount int32
 var failCount int32
 results := make([]ProcessResult, 0)
 var mu sync.Mutex

 err := bp.mi.Execute(ctx, items, func(ctx context.Context, item interface{}) error {
  result := bp.processItem(item)

  if result.Success {
   atomic.AddInt32(&successCount, 1)
  } else {
   atomic.AddInt32(&failCount, 1)
  }

  mu.Lock()
  results = append(results, result)
  mu.Unlock()

  return result.Error
 })

 return &BatchResult{
  Total:        len(items),
  SuccessCount: int(successCount),
  FailCount:    int(failCount),
  Results:      results,
 }, err
}

func (bp *BatchProcessor) processItem(item interface{}) ProcessResult {
 data := item.(map[string]interface{})
 id := data["id"].(string)

 time.Sleep(50 * time.Millisecond) // 模拟处理

 return ProcessResult{
  ID:      id,
  Success: true,
  Message: fmt.Sprintf("Processed %s", id),
 }
}

// ProcessResult 处理结果
type ProcessResult struct {
 ID      string
 Success bool
 Message string
 Error   error
}

// BatchResult 批量结果
type BatchResult struct {
 Total        int
 SuccessCount int
 FailCount    int
 Results      []ProcessResult
}

// ParallelDownloader 并行下载器
type ParallelDownloader struct {
 mi *MultiInstanceNoSync
}

// NewParallelDownloader 创建并行下载器
func NewParallelDownloader() *ParallelDownloader {
 return &ParallelDownloader{
  mi: &MultiInstanceNoSync{},
 }
}

// DownloadAll 并行下载所有文件
func (pd *ParallelDownloader) DownloadAll(ctx context.Context, urls []string) ([]DownloadResult, error) {
 items := make([]interface{}, len(urls))
 for i, url := range urls {
  items[i] = url
 }

 results := make([]DownloadResult, len(urls))

 err := pd.mi.Execute(ctx, items, func(ctx context.Context, item interface{}) error {
  url := item.(string)
  result := pd.downloadFile(ctx, url)

  // 找到对应的位置存储结果
  for i, u := range urls {
   if u == url {
    results[i] = result
    break
   }
  }

  return result.Error
 })

 return results, err
}

func (pd *ParallelDownloader) downloadFile(ctx context.Context, url string) DownloadResult {
 time.Sleep(100 * time.Millisecond) // 模拟下载

 return DownloadResult{
  URL:      url,
  Success:  true,
  FileSize: 1024,
 }
}

// DownloadResult 下载结果
type DownloadResult struct {
 URL      string
 Success  bool
 FileSize int
 Error    error
}

func main() {
 ctx := context.Background()

 // 批处理示例
 fmt.Println("=== 批处理示例 ===")
 processor := NewBatchProcessor()

 items := []interface{}{
  map[string]interface{}{"id": "item1", "value": 100},
  map[string]interface{}{"id": "item2", "value": 200},
  map[string]interface{}{"id": "item3", "value": 300},
  map[string]interface{}{"id": "item4", "value": 400},
  map[string]interface{}{"id": "item5", "value": 500},
 }

 result, err := processor.ProcessBatch(ctx, items)
 if err != nil {
  log.Printf("批处理部分失败: %v", err)
 }

 fmt.Printf("处理完成: 总数=%d, 成功=%d, 失败=%d\n",
  result.Total, result.SuccessCount, result.FailCount)

 // 并行下载示例
 fmt.Println("\n=== 并行下载示例 ===")
 downloader := NewParallelDownloader()
 urls := []string{
  "https://example.com/file1.pdf",
  "https://example.com/file2.pdf",
  "https://example.com/file3.pdf",
 }

 downloadResults, err := downloader.DownloadAll(ctx, urls)
 if err != nil {
  log.Printf("部分下载失败: %v", err)
 }

 for _, r := range downloadResults {
  fmt.Printf("- %s: %d bytes\n", r.URL, r.FileSize)
 }
}
```

#### 反例说明

**错误实现1：同步等待**

```go
// ❌ 错误：使用同步等待所有实例
func WrongSyncWait() {
 var wg sync.WaitGroup
 for _, item := range items {
  wg.Add(1)
  go func(i Item) {
   defer wg.Done()
   process(i)
  }(item)
 }
 wg.Wait()  // 错误：无同步模式不应该等待
}
```

**问题**：无同步模式不应该等待所有实例完成，应该立即返回。

**错误实现2：资源耗尽**

```go
// ❌ 错误：创建过多goroutine
func WrongResourceExhaustion() {
 for _, item := range hugeList {  // 假设有100万个项目
  go process(item)  // 错误：创建100万个goroutine
 }
}
```

**问题**：创建过多goroutine会导致资源耗尽。应该使用工作池限制并发数。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 批量通知 | 批量发送通知，不需要等待 |
| 日志记录 | 并行记录日志到多个目标 |
| 数据导入 | 批量导入数据到数据库 |
| 文件处理 | 并行处理多个文件 |
| 事件发布 | 并行发布事件到多个订阅者 |

#### 与其他模式关系

- **与并行分支**：多实例无同步是并行分支的泛化
- **与多实例需同步**：无同步不需要等待，需同步需要等待
- **与N选M**：N选M可以基于多实例实现

---

### 4.2 多实例需同步（Multiple Instances with a Priori Design Time Knowledge）

#### 概念定义

**多实例需同步模式**在设计时就知道实例数量，所有实例完成后才触发后续任务。这是实现批量并行处理并等待结果的模式。

**形式化定义**：给定活动 A 和固定的实例数 n（设计时已知），多实例需同步模式创建 n 个 A 实例，并等待所有实例完成后才继续。

#### BPMN描述

在BPMN中，多实例需同步通过**多实例活动**的顺序设置实现：

```
[任务] * (顺序，需同步)
```

- 多实例标记（三条竖线）
- 顺序执行设置
- 所有实例完成后汇聚

#### 流程图

```
┌─────────┐     ┌─────────────┐
│  数据   │────→│  多实例活动 │
│  集合   │     │  (需同步)   │
└─────────┘     └──────┬──────┘
                       │
         ┌─────────────┼─────────────┐
         │             │             │
         ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │ 实例1   │   │ 实例2   │   │ 实例3   │
   └────┬────┘   └────┬────┘   └────┬────┘
        │             │             │
        └─────────────┼─────────────┘
                      │
                      ▼
               ┌─────────────┐
               │   同步点    │
               └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
)

// MultiInstanceWithSync 多实例需同步执行器
type MultiInstanceWithSync struct{}

// Execute 执行多实例，需同步
func (m *MultiInstanceWithSync) Execute(ctx context.Context, items []interface{}, handler func(ctx context.Context, item interface{}) (interface{}, error)) ([]interface{}, error) {
 var wg sync.WaitGroup
 results := make([]interface{}, len(items))
 errChan := make(chan error, len(items))

 log.Printf("[多实例需同步] 启动 %d 个实例", len(items))

 for i, item := range items {
  wg.Add(1)
  go func(index int, it interface{}) {
   defer wg.Done()

   log.Printf("[多实例需同步] 实例 %d 开始", index+1)

   result, err := handler(ctx, it)
   if err != nil {
    errChan <- fmt.Errorf("实例 %d 失败: %w", index+1, err)
    return
   }

   results[index] = result
   log.Printf("[多实例需同步] 实例 %d 完成", index+1)
  }(i, item)
 }

 wg.Wait()
 close(errChan)

 for err := range errChan {
  return nil, err
 }

 log.Println("[多实例需同步] 所有实例已完成")
 return results, nil
}

// 使用示例
func main() {
 mi := &MultiInstanceWithSync{}
 ctx := context.Background()

 // 批量处理并收集结果
 items := []interface{}{1, 2, 3, 4, 5}

 results, err := mi.Execute(ctx, items, func(ctx context.Context, item interface{}) (interface{}, error) {
  num := item.(int)
  result := num * num
  fmt.Printf("→ 处理 %d -> %d\n", num, result)
  return result, nil
 })

 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n✓ 所有结果: %v\n", results)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// ParallelMapper 并行映射器
type ParallelMapper struct {
 mi *MultiInstanceWithSync
}

// NewParallelMapper 创建并行映射器
func NewParallelMapper() *ParallelMapper {
 return &ParallelMapper{
  mi: &MultiInstanceWithSync{},
 }
}

// Map 并行映射
func (pm *ParallelMapper) Map(ctx context.Context, items []interface{}, mapper func(interface{}) interface{}) ([]interface{}, error) {
 return pm.mi.Execute(ctx, items, func(ctx context.Context, item interface{}) (interface{}, error) {
  return mapper(item), nil
 })
}

// ParallelFilter 并行过滤器
type ParallelFilter struct {
 mi *MultiInstanceWithSync
}

// NewParallelFilter 创建并行过滤器
func NewParallelFilter() *ParallelFilter {
 return &ParallelFilter{
  mi: &MultiInstanceWithSync{},
 }
}

// Filter 并行过滤
func (pf *ParallelFilter) Filter(ctx context.Context, items []interface{}, predicate func(interface{}) bool) ([]interface{}, error) {
 results, err := pf.mi.Execute(ctx, items, func(ctx context.Context, item interface{}) (interface{}, error) {
  return predicate(item), nil
 })

 if err != nil {
  return nil, err
 }

 // 收集过滤后的结果
 var filtered []interface{}
 for i, result := range results {
  if result.(bool) {
   filtered = append(filtered, items[i])
  }
 }

 return filtered, nil
}

// DataEnricher 数据增强器
type DataEnricher struct {
 mi *MultiInstanceWithSync
}

// NewDataEnricher 创建数据增强器
func NewDataEnricher() *DataEnricher {
 return &DataEnricher{
  mi: &MultiInstanceWithSync{},
 }
}

// Enrich 并行增强数据
func (de *DataEnricher) Enrich(ctx context.Context, records []Record) ([]EnrichedRecord, error) {
 items := make([]interface{}, len(records))
 for i, r := range records {
  items[i] = r
 }

 results, err := de.mi.Execute(ctx, items, func(ctx context.Context, item interface{}) (interface{}, error) {
  record := item.(Record)
  return de.enrichRecord(ctx, record)
 })

 if err != nil {
  return nil, err
 }

 enriched := make([]EnrichedRecord, len(results))
 for i, r := range results {
  enriched[i] = r.(EnrichedRecord)
 }

 return enriched, nil
}

func (de *DataEnricher) enrichRecord(ctx context.Context, record Record) (EnrichedRecord, error) {
 // 模拟数据增强
 time.Sleep(50 * time.Millisecond)

 return EnrichedRecord{
  ID:        record.ID,
  Name:      record.Name,
  ExtraData: fmt.Sprintf("enriched_%s", record.ID),
 }, nil
}

// Record 记录
type Record struct {
 ID   string
 Name string
}

// EnrichedRecord 增强后的记录
type EnrichedRecord struct {
 ID        string
 Name      string
 ExtraData string
}

func main() {
 ctx := context.Background()

 // 并行映射示例
 fmt.Println("=== 并行映射示例 ===")
 mapper := NewParallelMapper()

 items := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
 results, err := mapper.Map(ctx, items, func(item interface{}) interface{} {
  num := item.(int)
  return num * num
 })

 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("映射结果: %v\n", results)

 // 并行过滤示例
 fmt.Println("\n=== 并行过滤示例 ===")
 filter := NewParallelFilter()

 numbers := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
 evenNumbers, err := filter.Filter(ctx, numbers, func(item interface{}) bool {
  num := item.(int)
  return num%2 == 0
 })

 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("偶数: %v\n", evenNumbers)

 // 数据增强示例
 fmt.Println("\n=== 数据增强示例 ===")
 enricher := NewDataEnricher()

 records := []Record{
  {ID: "001", Name: "Alice"},
  {ID: "002", Name: "Bob"},
  {ID: "003", Name: "Charlie"},
 }

 enriched, err := enricher.Enrich(ctx, records)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("增强结果:")
 for _, r := range enriched {
  fmt.Printf("  - %s: %s (%s)\n", r.ID, r.Name, r.ExtraData)
 }
}
```

#### 反例说明

**错误实现1：不等待所有实例**

```go
// ❌ 错误：不等待所有实例完成
func WrongNoWait() {
 for _, item := range items {
  go func(i Item) {
   process(i)
  }(item)
 }
 // 错误：没有等待，函数立即返回
}
```

**问题**：函数立即返回，实例可能还在执行中，导致结果不完整。

**错误实现2：结果顺序混乱**

```go
// ❌ 错误：结果顺序与输入不一致
func WrongOrder() {
 results := make([]Result, 0)

 for _, item := range items {
  go func(i Item) {
   result := process(i)
   results = append(results, result)  // 错误：顺序不确定
  }(item)
 }
}
```

**问题**：结果顺序与输入顺序不一致。应该使用索引存储结果。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 批量处理 | 批量处理并收集所有结果 |
| 数据转换 | 并行转换数据并收集结果 |
| 并行计算 | 分片计算后合并结果 |
| 数据增强 | 并行增强数据记录 |
| 批量验证 | 批量验证并收集结果 |

#### 与其他模式关系

- **与多实例无同步**：需同步需要等待所有实例，无同步不需要
- **与同步模式**：需同步使用同步模式汇聚
- **与并行分支**：需同步是并行分支的批量形式

---

### 4.3 多实例运行时确定（Multiple Instances with a Priori Runtime Knowledge）

#### 概念定义

**多实例运行时确定模式**在运行时才能确定实例数量，所有实例完成后才触发后续任务。这是实现动态批量处理的模式。

**形式化定义**：给定活动 A 和实例数 n（运行时确定），多实例运行时确定模式在运行时创建 n 个 A 实例，并等待所有实例完成后才继续。

#### BPMN描述

在BPMN中，多实例运行时确定通过**多实例活动**的动态设置实现：

```
[任务] * (运行时确定数量)
```

- 多实例标记（三条竖线）
- 运行时确定实例数
- 所有实例完成后汇聚

#### 流程图

```
┌─────────┐     ┌─────────────┐     ┌─────────────┐
│  查询   │────→│  确定实例数 │────→│  多实例活动 │
│  数据   │     │  (运行时)   │     │  (需同步)   │
└─────────┘     └─────────────┘     └──────┬──────┘
                                           │
                              ┌────────────┼────────────┐
                              │            │            │
                              ▼            ▼            ▼
                        ┌─────────┐  ┌─────────┐  ┌─────────┐
                        │ 实例1   │  │ 实例2   │  │ 实例N   │
                        └────┬────┘  └────┬────┘  └────┬────┘
                             │            │            │
                             └────────────┼────────────┘
                                          │
                                          ▼
                                   ┌─────────────┐
                                   │   同步点    │
                                   └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
)

// MultiInstanceRuntime 多实例运行时确定执行器
type MultiInstanceRuntime struct {
 getItems func(ctx context.Context) ([]interface{}, error)
}

// NewMultiInstanceRuntime 创建多实例运行时确定执行器
func NewMultiInstanceRuntime(getItems func(ctx context.Context) ([]interface{}, error)) *MultiInstanceRuntime {
 return &MultiInstanceRuntime{getItems: getItems}
}

// Execute 执行多实例，运行时确定数量
func (m *MultiInstanceRuntime) Execute(ctx context.Context, handler func(ctx context.Context, item interface{}) (interface{}, error)) ([]interface{}, error) {
 // 运行时获取实例列表
 items, err := m.getItems(ctx)
 if err != nil {
  return nil, fmt.Errorf("获取实例列表失败: %w", err)
 }

 log.Printf("[多实例运行时] 运行时确定 %d 个实例", len(items))

 var wg sync.WaitGroup
 results := make([]interface{}, len(items))
 errChan := make(chan error, len(items))

 for i, item := range items {
  wg.Add(1)
  go func(index int, it interface{}) {
   defer wg.Done()

   result, err := handler(ctx, it)
   if err != nil {
    errChan <- fmt.Errorf("实例 %d 失败: %w", index+1, err)
    return
   }

   results[index] = result
  }(i, item)
 }

 wg.Wait()
 close(errChan)

 for err := range errChan {
  return nil, err
 }

 return results, nil
}

// 使用示例
func main() {
 // 运行时确定实例数量
 itemProvider := func(ctx context.Context) ([]interface{}, error) {
  // 模拟从数据库查询需要处理的任务
  return []interface{}{
   "task_from_db_1",
   "task_from_db_2",
   "task_from_db_3",
  }, nil
 }

 mi := NewMultiInstanceRuntime(itemProvider)
 ctx := context.Background()

 results, err := mi.Execute(ctx, func(ctx context.Context, item interface{}) (interface{}, error) {
  task := item.(string)
  fmt.Printf("→ 处理任务: %s\n", task)
  return fmt.Sprintf("processed_%s", task), nil
 })

 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n✓ 所有结果: %v\n", results)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// DynamicProcessor 动态处理器
type DynamicProcessor struct {
 itemProvider func(ctx context.Context) ([]interface{}, error)
}

// NewDynamicProcessor 创建动态处理器
func NewDynamicProcessor(provider func(ctx context.Context) ([]interface{}, error)) *DynamicProcessor {
 return &DynamicProcessor{itemProvider: provider}
}

// Process 动态处理
func (dp *DynamicProcessor) Process(ctx context.Context, processor func(ctx context.Context, item interface{}) (interface{}, error)) ([]interface{}, error) {
 mi := NewMultiInstanceRuntime(dp.itemProvider)
 return mi.Execute(ctx, processor)
}

// QueueConsumer 队列消费者
type QueueConsumer struct {
 queueName string
}

// NewQueueConsumer 创建队列消费者
func NewQueueConsumer(queueName string) *QueueConsumer {
 return &QueueConsumer{queueName: queueName}
}

// Consume 消费队列消息
func (qc *QueueConsumer) Consume(ctx context.Context, batchSize int, processor func(msg Message) error) ([]ProcessingResult, error) {
 // 运行时获取消息
 itemProvider := func(ctx context.Context) ([]interface{}, error) {
  messages := qc.fetchMessages(batchSize)
  items := make([]interface{}, len(messages))
  for i, m := range messages {
   items[i] = m
  }
  return items, nil
 }

 mi := NewMultiInstanceRuntime(itemProvider)

 results, err := mi.Execute(ctx, func(ctx context.Context, item interface{}) (interface{}, error) {
  msg := item.(Message)
  start := time.Now()

  if err := processor(msg); err != nil {
   return ProcessingResult{
    MessageID: msg.ID,
    Success:   false,
    Error:     err,
   }, nil
  }

  return ProcessingResult{
   MessageID: msg.ID,
   Success:   true,
   Duration:  time.Since(start),
  }, nil
 })

 if err != nil {
  return nil, err
 }

 processingResults := make([]ProcessingResult, len(results))
 for i, r := range results {
  processingResults[i] = r.(ProcessingResult)
 }

 return processingResults, nil
}

func (qc *QueueConsumer) fetchMessages(limit int) []Message {
 // 模拟从队列获取消息
 messages := []Message{
  {ID: "msg1", Content: "content1"},
  {ID: "msg2", Content: "content2"},
  {ID: "msg3", Content: "content3"},
 }

 if limit < len(messages) {
  return messages[:limit]
 }
 return messages
}

// Message 消息
type Message struct {
 ID      string
 Content string
}

// ProcessingResult 处理结果
type ProcessingResult struct {
 MessageID string
 Success   bool
 Duration  time.Duration
 Error     error
}

// PaginatedFetcher 分页获取器
type PaginatedFetcher struct {
 fetchPage func(page, pageSize int) ([]interface{}, error)
}

// NewPaginatedFetcher 创建分页获取器
func NewPaginatedFetcher(fetcher func(page, pageSize int) ([]interface{}, error)) *PaginatedFetcher {
 return &PaginatedFetcher{fetchPage: fetcher}
}

// FetchAll 获取所有数据
func (pf *PaginatedFetcher) FetchAll(ctx context.Context, pageSize int, processor func(ctx context.Context, item interface{}) (interface{}, error)) ([]interface{}, error) {
 var allResults []interface{}
 page := 1

 for {
  // 运行时获取当前页的数据
  itemProvider := func(ctx context.Context) ([]interface{}, error) {
   return pf.fetchPage(page, pageSize)
  }

  mi := NewMultiInstanceRuntime(itemProvider)
  results, err := mi.Execute(ctx, processor)

  if err != nil {
   return nil, err
  }

  if len(results) == 0 {
   break
  }

  allResults = append(allResults, results...)
  page++

  // 假设每页最多pageSize条，如果少于pageSize说明是最后一页
  if len(results) < pageSize {
   break
  }
 }

 return allResults, nil
}

func main() {
 ctx := context.Background()

 // 队列消费示例
 fmt.Println("=== 队列消费示例 ===")
 consumer := NewQueueConsumer("order_queue")

 results, err := consumer.Consume(ctx, 10, func(msg Message) error {
  fmt.Printf("→ 处理消息: %s\n", msg.ID)
  time.Sleep(50 * time.Millisecond)
  return nil
 })

 if err != nil {
  log.Fatal(err)
 }

 successCount := 0
 for _, r := range results {
  if r.Success {
   successCount++
  }
 }

 fmt.Printf("处理完成: 成功=%d, 总数=%d\n", successCount, len(results))

 // 分页获取示例
 fmt.Println("\n=== 分页获取示例 ===")

 fetchCount := 0
 fetcher := NewPaginatedFetcher(func(page, pageSize int) ([]interface{}, error) {
  fetchCount++
  // 模拟分页数据
  if page > 2 {
   return nil, nil
  }

  items := []interface{}{}
  for i := 0; i < pageSize && i < 5; i++ {
   items = append(items, fmt.Sprintf("item_page%d_%d", page, i))
  }
  return items, nil
 })

 allResults, err := fetcher.FetchAll(ctx, 3, func(ctx context.Context, item interface{}) (interface{}, error) {
  fmt.Printf("→ 处理: %s\n", item)
  return fmt.Sprintf("processed_%s", item), nil
 })

 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n总共获取 %d 条数据\n", len(allResults))
}
```

#### 反例说明

**错误实现1：设计时硬编码数量**

```go
// ❌ 错误：设计时硬编码数量
func WrongHardcodedCount() {
 // 错误：假设总是有3个实例
 for i := 0; i < 3; i++ {
  go process(items[i])
 }
}
```

**问题**：硬编码数量无法适应动态变化的数据量。

**错误实现2：运行时panic**

```go
// ❌ 错误：运行时可能panic
func WrongRuntimePanic() {
 items := getItemsFromDB()  // 可能返回nil
 for _, item := range items {  // 错误：nil range会panic
  process(item)
 }
}
```

**问题**：没有处理边界情况，可能导致panic。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 队列消费 | 运行时确定队列中的消息数量 |
| 分页处理 | 运行时确定每页的数据量 |
| 动态任务 | 运行时确定需要处理的任务数 |
| 批量导入 | 运行时确定导入文件中的记录数 |
| 流式处理 | 运行时确定缓冲区中的数据量 |

#### 与其他模式关系

- **与多实例需同步**：运行时确定是需同步的动态版本
- **与任意循环**：可以使用循环实现运行时确定
- **与多实例无同步**：运行时确定可以与无同步结合

---

## 5. 状态模式

### 5.1 延迟选择（Deferred Choice）

#### 概念定义

**延迟选择模式**将决策点推迟到运行时，基于外部事件或数据来决定执行路径。与排他选择不同，延迟选择在设计时不确定哪个分支会被执行。

**形式化定义**：给定分支集合 B = {b₁, b₂, ..., bₙ}，延迟选择模式在运行时基于事件 e 选择分支 bᵢ，其中 e 触发 bᵢ 的执行。

#### BPMN描述

在BPMN中，延迟选择通过**基于事件网关（Event-Based Gateway）**实现：

```
[前置任务] → [⚡] →┬→ [事件A] → [任务A]
                    ├→ [事件B] → [任务B]
                    └→ [事件C] → [任务C]
```

- 基于事件网关（⚡）
- 等待第一个发生的事件
- 触发对应的分支

#### 流程图

```
┌─────────┐     ┌─────────────┐
│  开始   │────→│  延迟选择   │
└─────────┘     │  网关       │
                └──────┬──────┘
                       │
         ┌─────────────┼─────────────┐
         │             │             │
         ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │ 事件A   │   │ 事件B   │   │ 事件C   │
   └────┬────┘   └────┬────┘   └────┬────┘
        │             │             │
        ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │  任务A  │   │  任务B  │   │  任务C  │
   └─────────┘   └─────────┘   └─────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
)

// DeferredChoice 延迟选择器
type DeferredChoice struct {
 branches map[string]struct {
  event   chan interface{}
  handler func(ctx context.Context, data interface{}) error
 }
 mu sync.RWMutex
}

// NewDeferredChoice 创建延迟选择器
func NewDeferredChoice() *DeferredChoice {
 return &DeferredChoice{
  branches: make(map[string]struct {
   event   chan interface{}
   handler func(ctx context.Context, data interface{}) error
  }),
 }
}

// On 注册事件处理器
func (d *DeferredChoice) On(eventName string, handler func(ctx context.Context, data interface{}) error) {
 d.mu.Lock()
 defer d.mu.Unlock()

 d.branches[eventName] = struct {
  event   chan interface{}
  handler func(ctx context.Context, data interface{}) error
 }{
  event:   make(chan interface{}, 1),
  handler: handler,
 }
}

// Trigger 触发事件
func (d *DeferredChoice) Trigger(eventName string, data interface{}) error {
 d.mu.RLock()
 branch, ok := d.branches[eventName]
 d.mu.RUnlock()

 if !ok {
  return fmt.Errorf("未知事件: %s", eventName)
 }

 select {
 case branch.event <- data:
  return nil
 default:
  return fmt.Errorf("事件 '%s' 已经被触发", eventName)
 }
}

// Wait 等待第一个事件并执行
func (d *DeferredChoice) Wait(ctx context.Context) (string, error) {
 d.mu.RLock()
 cases := make([]<-chan interface{}, 0, len(d.branches))
 names := make([]string, 0, len(d.branches))

 for name, branch := range d.branches {
  cases = append(cases, branch.event)
  names = append(names, name)
 }
 d.mu.RUnlock()

 // 使用select等待第一个事件
 selectCases := make([]reflect.SelectCase, len(cases))
 for i, ch := range cases {
  selectCases[i] = reflect.SelectCase{
   Dir:  reflect.SelectRecv,
   Chan: reflect.ValueOf(ch),
  }
 }

 // 添加上下文取消
 ctxCase := reflect.SelectCase{
  Dir:  reflect.SelectRecv,
  Chan: reflect.ValueOf(ctx.Done()),
 }
 selectCases = append(selectCases, ctxCase)

 chosen, value, _ := reflect.Select(selectCases)

 if chosen == len(selectCases)-1 {
  return "", ctx.Err()
 }

 eventName := names[chosen]

 d.mu.RLock()
 branch := d.branches[eventName]
 d.mu.RUnlock()

 log.Printf("[延迟选择] 事件 '%s' 被触发", eventName)

 if err := branch.handler(ctx, value.Interface()); err != nil {
  return eventName, err
 }

 return eventName, nil
}

// 简化实现（不使用reflect）
type SimpleDeferredChoice struct {
 branches map[string]chan Event
 handlers map[string]func(ctx context.Context, data interface{}) error
 mu       sync.RWMutex
}

type Event struct {
 Name string
 Data interface{}
}

// NewSimpleDeferredChoice 创建简化延迟选择器
func NewSimpleDeferredChoice() *SimpleDeferredChoice {
 return &SimpleDeferredChoice{
  branches: make(map[string]chan Event),
  handlers: make(map[string]func(ctx context.Context, data interface{}) error),
 }
}

// On 注册事件
func (d *SimpleDeferredChoice) On(eventName string, handler func(ctx context.Context, data interface{}) error) {
 d.mu.Lock()
 defer d.mu.Unlock()

 d.branches[eventName] = make(chan Event, 1)
 d.handlers[eventName] = handler
}

// Trigger 触发事件
func (d *SimpleDeferredChoice) Trigger(eventName string, data interface{}) {
 d.mu.RLock()
 ch, ok := d.branches[eventName]
 d.mu.RUnlock()

 if !ok {
  return
 }

 select {
 case ch <- Event{Name: eventName, Data: data}:
 default:
 }
}

// Wait 等待并处理第一个事件
func (d *SimpleDeferredChoice) Wait(ctx context.Context) (string, error) {
 d.mu.RLock()
 chans := make([]chan Event, 0, len(d.branches))
 names := make([]string, 0, len(d.branches))

 for name, ch := range d.branches {
  chans = append(chans, ch)
  names = append(names, name)
 }
 d.mu.RUnlock()

 // 创建合并的channel
 merged := make(chan Event)
 var wg sync.WaitGroup

 for _, ch := range chans {
  wg.Add(1)
  go func(c chan Event) {
   defer wg.Done()
   select {
   case evt := <-c:
    merged <- evt
   case <-ctx.Done():
   }
  }(ch)
 }

 go func() {
  wg.Wait()
  close(merged)
 }()

 select {
 case evt := <-merged:
  d.mu.RLock()
  handler := d.handlers[evt.Name]
  d.mu.RUnlock()

  if handler != nil {
   return evt.Name, handler(ctx, evt.Data)
  }
  return evt.Name, nil

 case <-ctx.Done():
  return "", ctx.Err()
 }
}

// 使用示例
func main() {
 choice := NewSimpleDeferredChoice()

 // 注册事件处理器
 choice.On("payment_success", func(ctx context.Context, data interface{}) error {
  fmt.Printf("→ 处理支付成功: %v\n", data)
  return nil
 })

 choice.On("payment_failed", func(ctx context.Context, data interface{}) error {
  fmt.Printf("→ 处理支付失败: %v\n", data)
  return nil
 })

 choice.On("payment_timeout", func(ctx context.Context, data interface{}) error {
  fmt.Printf("→ 处理支付超时: %v\n", data)
  return nil
 })

 // 模拟异步事件触发
 go func() {
  time.Sleep(100 * time.Millisecond)
  choice.Trigger("payment_success", map[string]string{"order_id": "123"})
 }()

 ctx := context.Background()
 chosen, err := choice.Wait(ctx)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("\n✓ 选择了分支: %s\n", chosen)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// OrderWorkflow 订单工作流
type OrderWorkflow struct {
 choice *SimpleDeferredChoice
}

// NewOrderWorkflow 创建订单工作流
func NewOrderWorkflow() *OrderWorkflow {
 return &OrderWorkflow{
  choice: NewSimpleDeferredChoice(),
 }
}

// ProcessOrder 处理订单
func (ow *OrderWorkflow) ProcessOrder(ctx context.Context, orderID string) (string, error) {
 // 注册可能的事件
 ow.choice.On("paid", func(ctx context.Context, data interface{}) error {
  fmt.Printf("→ 订单 %s 已支付\n", orderID)
  return nil
 })

 ow.choice.On("cancelled", func(ctx context.Context, data interface{}) error {
  fmt.Printf("→ 订单 %s 已取消\n", orderID)
  return nil
 })

 ow.choice.On("expired", func(ctx context.Context, data interface{}) error {
  fmt.Printf("→ 订单 %s 已过期\n", orderID)
  return nil
 })

 fmt.Printf("等待订单 %s 的结果...\n", orderID)

 return ow.choice.Wait(ctx)
}

// ExternalEventSource 外部事件源
type ExternalEventSource struct {
 listeners []func(event string, data interface{})
 mu        sync.Mutex
}

// NewExternalEventSource 创建外部事件源
func NewExternalEventSource() *ExternalEventSource {
 return &ExternalEventSource{
  listeners: make([]func(event string, data interface{}), 0),
 }
}

// Subscribe 订阅事件
func (es *ExternalEventSource) Subscribe(listener func(event string, data interface{})) {
 es.mu.Lock()
 defer es.mu.Unlock()
 es.listeners = append(es.listeners, listener)
}

// Emit 触发事件
func (es *ExternalEventSource) Emit(event string, data interface{}) {
 es.mu.Lock()
 listeners := make([]func(event string, data interface{}), len(es.listeners))
 copy(listeners, es.listeners)
 es.mu.Unlock()

 for _, listener := range listeners {
  go listener(event, data)
 }
}

// TimeoutManager 超时管理器
type TimeoutManager struct {
 timeouts map[string]*time.Timer
 mu       sync.Mutex
}

// NewTimeoutManager 创建超时管理器
func NewTimeoutManager() *TimeoutManager {
 return &TimeoutManager{
  timeouts: make(map[string]*time.Timer),
 }
}

// SetTimeout 设置超时
func (tm *TimeoutManager) SetTimeout(id string, duration time.Duration, callback func()) {
 tm.mu.Lock()
 defer tm.mu.Unlock()

 // 取消已有的超时
 if timer, ok := tm.timeouts[id]; ok {
  timer.Stop()
 }

 tm.timeouts[id] = time.AfterFunc(duration, callback)
}

// CancelTimeout 取消超时
func (tm *TimeoutManager) CancelTimeout(id string) {
 tm.mu.Lock()
 defer tm.mu.Unlock()

 if timer, ok := tm.timeouts[id]; ok {
  timer.Stop()
  delete(tm.timeouts, id)
 }
}

// AsyncTaskManager 异步任务管理器
type AsyncTaskManager struct {
 tasks map[string]chan TaskResult
 mu    sync.Mutex
}

type TaskResult struct {
 TaskID string
 Result interface{}
 Error  error
}

// NewAsyncTaskManager 创建异步任务管理器
func NewAsyncTaskManager() *AsyncTaskManager {
 return &AsyncTaskManager{
  tasks: make(map[string]chan TaskResult),
 }
}

// SubmitTask 提交任务
func (atm *AsyncTaskManager) SubmitTask(taskID string, task func() (interface{}, error)) {
 atm.mu.Lock()
 resultChan := make(chan TaskResult, 1)
 atm.tasks[taskID] = resultChan
 atm.mu.Unlock()

 go func() {
  result, err := task()
  resultChan <- TaskResult{
   TaskID: taskID,
   Result: result,
   Error:  err,
  }
 }()
}

// WaitForAny 等待任意任务完成
func (atm *AsyncTaskManager) WaitForAny(ctx context.Context) (*TaskResult, error) {
 atm.mu.Lock()
 chans := make([]chan TaskResult, 0, len(atm.tasks))
 for _, ch := range atm.tasks {
  chans = append(chans, ch)
 }
 atm.mu.Unlock()

 if len(chans) == 0 {
  return nil, fmt.Errorf("没有任务")
 }

 // 合并所有channel
 merged := make(chan TaskResult)
 var wg sync.WaitGroup

 for _, ch := range chans {
  wg.Add(1)
  go func(c chan TaskResult) {
   defer wg.Done()
   select {
   case result := <-c:
    merged <- result
   case <-ctx.Done():
   }
  }(ch)
 }

 go func() {
  wg.Wait()
  close(merged)
 }()

 select {
 case result := <-merged:
  return &result, nil
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

func main() {
 // 订单工作流示例
 fmt.Println("=== 订单工作流示例 ===")
 workflow := NewOrderWorkflow()

 // 模拟异步事件
 go func() {
  time.Sleep(500 * time.Millisecond)
  workflow.choice.Trigger("paid", map[string]string{"amount": "100"})
 }()

 ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
 defer cancel()

 result, err := workflow.ProcessOrder(ctx, "ORDER_001")
 if err != nil {
  log.Printf("工作流错误: %v", err)
 } else {
  fmt.Printf("✓ 订单处理结果: %s\n", result)
 }

 // 异步任务管理器示例
 fmt.Println("\n=== 异步任务管理器示例 ===")
 atm := NewAsyncTaskManager()

 atm.SubmitTask("task1", func() (interface{}, error) {
  time.Sleep(200 * time.Millisecond)
  return "result1", nil
 })

 atm.SubmitTask("task2", func() (interface{}, error) {
  time.Sleep(100 * time.Millisecond)
  return "result2", nil
 })

 atm.SubmitTask("task3", func() (interface{}, error) {
  time.Sleep(300 * time.Millisecond)
  return "result3", nil
 })

 ctx2 := context.Background()
 firstResult, err := atm.WaitForAny(ctx2)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("✓ 第一个完成的任务: %s, 结果: %v\n", firstResult.TaskID, firstResult.Result)
}
```

#### 反例说明

**错误实现1：设计时硬编码选择**

```go
// ❌ 错误：设计时硬编码选择
func WrongHardcodedChoice() {
 if config.PaymentMethod == "alipay" {
  processAlipay()
 } else if config.PaymentMethod == "wechat" {
  processWechat()
 }
 // 错误：无法处理运行时的新支付方式
}
```

**问题**：设计时硬编码无法适应运行时变化。

**错误实现2：轮询代替事件驱动**

```go
// ❌ 错误：使用轮询代替事件驱动
func WrongPolling() {
 for {
  if checkPaymentSuccess() {
   processSuccess()
   break
  }
  if checkPaymentFailed() {
   processFailed()
   break
  }
  time.Sleep(100 * time.Millisecond)  // 浪费资源
 }
}
```

**问题**：轮询浪费资源，应该使用事件驱动。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 异步处理 | 等待异步操作完成 |
| 事件驱动 | 基于事件选择执行路径 |
| 超时处理 | 等待操作完成或超时 |
| 竞争条件 | 多个操作竞争，选择第一个完成的 |
| 用户交互 | 等待用户输入 |

#### 与其他模式关系

- **与排他选择**：延迟选择推迟决策到运行时
- **与鉴别器**：两者都选择第一个完成的
- **与事件驱动**：延迟选择基于事件驱动

---

### 5.2 交错并行路由（Interleaved Parallel Routing）

#### 概念定义

**交错并行路由模式**允许多个任务交错执行，但每个任务必须完整执行后才能开始下一个任务。这是实现协作式多任务的模式。

**形式化定义**：给定任务集合 T = {t₁, t₂, ..., tₙ}，交错并行路由模式允许任务以任意顺序执行，但每个任务必须原子性完成。

#### BPMN描述

在BPMN中，交错并行路由通过**复杂网关**或特定的交错语义实现：

```
[任务A] ─┐
[任务B] ─┼→ [交错路由] → [后续任务]
[任务C] ─┘
```

- 任务可以交错执行
- 每个任务原子性完成
- 所有任务完成后继续

#### 流程图

```
┌─────────┐     ┌─────────────┐
│  开始   │────→│  交错并行   │
└─────────┘     │   路由      │
                └──────┬──────┘
                       │
         ┌─────────────┼─────────────┐
         │             │             │
         ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │  任务A  │   │  任务B  │   │  任务C  │
   │ (原子)  │   │ (原子)  │   │ (原子)  │
   └────┬────┘   └────┬────┘   └────┬────┘
        │             │             │
        └─────────────┼─────────────┘
                      │
                      ▼
               ┌─────────────┐
               │   后续任务  │
               └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
)

// InterleavedRouting 交错并行路由器
type InterleavedRouting struct {
 tasks []InterleavedTask
}

// InterleavedTask 交错任务
type InterleavedTask struct {
 Name   string
 Action func(ctx context.Context) error
}

// NewInterleavedRouting 创建交错并行路由器
func NewInterleavedRouting() *InterleavedRouting {
 return &InterleavedRouting{
  tasks: make([]InterleavedTask, 0),
 }
}

// Add 添加任务
func (ir *InterleavedRouting) Add(name string, action func(ctx context.Context) error) {
 ir.tasks = append(ir.tasks, InterleavedTask{Name: name, Action: action})
}

// Execute 交错执行所有任务
func (ir *InterleavedRouting) Execute(ctx context.Context) error {
 if len(ir.tasks) == 0 {
  return nil
 }

 // 使用channel实现交错执行
 taskChan := make(chan InterleavedTask, len(ir.tasks))
 resultChan := make(chan struct {
  name string
  err  error
 }, len(ir.tasks))

 // 将任务放入channel
 for _, task := range ir.tasks {
  taskChan <- task
 }
 close(taskChan)

 // 启动固定数量的worker
 workerCount := 3
 var wg sync.WaitGroup

 for i := 0; i < workerCount; i++ {
  wg.Add(1)
  go func(workerID int) {
   defer wg.Done()

   for task := range taskChan {
    // 每个任务原子性执行
    log.Printf("[交错路由] Worker %d 执行任务: %s", workerID, task.Name)

    err := task.Action(ctx)
    resultChan <- struct {
     name string
     err  error
    }{task.Name, err}
   }
  }(i + 1)
 }

 // 等待所有worker完成
 go func() {
  wg.Wait()
  close(resultChan)
 }()

 // 收集结果
 var errs []error
 for result := range resultChan {
  if result.err != nil {
   errs = append(errs, fmt.Errorf("任务 '%s' 失败: %w", result.name, result.err))
  }
 }

 if len(errs) > 0 {
  return fmt.Errorf("部分任务失败: %v", errs)
 }

 return nil
}

// 使用示例
func main() {
 routing := NewInterleavedRouting()

 // 添加交错任务
 routing.Add("任务A", func(ctx context.Context) error {
  fmt.Println("→ 执行任务A")
  return nil
 })

 routing.Add("任务B", func(ctx context.Context) error {
  fmt.Println("→ 执行任务B")
  return nil
 })

 routing.Add("任务C", func(ctx context.Context) error {
  fmt.Println("→ 执行任务C")
  return nil
 })

 routing.Add("任务D", func(ctx context.Context) error {
  fmt.Println("→ 执行任务D")
  return nil
 })

 ctx := context.Background()
 if err := routing.Execute(ctx); err != nil {
  log.Fatal(err)
 }

 fmt.Println("\n✓ 所有交错任务完成")
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "sync/atomic"
 "time"
)

// CooperativeScheduler 协作式调度器
type CooperativeScheduler struct {
 tasks   []CooperativeTask
 quantum time.Duration
}

// CooperativeTask 协作式任务
type CooperativeTask struct {
 ID       string
 Priority int
 Work     func(ctx context.Context) error
}

// NewCooperativeScheduler 创建协作式调度器
func NewCooperativeScheduler(quantum time.Duration) *CooperativeScheduler {
 return &CooperativeScheduler{
  tasks:   make([]CooperativeTask, 0),
  quantum: quantum,
 }
}

// AddTask 添加任务
func (cs *CooperativeScheduler) AddTask(id string, priority int, work func(ctx context.Context) error) {
 cs.tasks = append(cs.tasks, CooperativeTask{
  ID:       id,
  Priority: priority,
  Work:     work,
 })
}

// Run 运行调度器
func (cs *CooperativeScheduler) Run(ctx context.Context) error {
 if len(cs.tasks) == 0 {
  return nil
 }

 // 按优先级排序
 sortedTasks := make([]CooperativeTask, len(cs.tasks))
 copy(sortedTasks, cs.tasks)

 // 简单的优先级排序
 for i := 0; i < len(sortedTasks); i++ {
  for j := i + 1; j < len(sortedTasks); j++ {
   if sortedTasks[i].Priority < sortedTasks[j].Priority {
    sortedTasks[i], sortedTasks[j] = sortedTasks[j], sortedTasks[i]
   }
  }
 }

 // 交错执行
 var wg sync.WaitGroup
 errChan := make(chan error, len(sortedTasks))

 for _, task := range sortedTasks {
  wg.Add(1)
  go func(t CooperativeTask) {
   defer wg.Done()

   // 创建带超时的上下文
   taskCtx, cancel := context.WithTimeout(ctx, cs.quantum)
   defer cancel()

   log.Printf("[协作调度] 开始任务: %s (优先级: %d)", t.ID, t.Priority)

   if err := t.Work(taskCtx); err != nil {
    errChan <- fmt.Errorf("任务 '%s' 失败: %w", t.ID, err)
    return
   }

   log.Printf("[协作调度] 完成任务: %s", t.ID)
  }(task)
 }

 wg.Wait()
 close(errChan)

 for err := range errChan {
  return err
 }

 return nil
}

// ResourcePool 资源池
type ResourcePool struct {
 resources chan Resource
 mu        sync.Mutex
}

// Resource 资源
type Resource struct {
 ID   string
 Name string
}

// NewResourcePool 创建资源池
func NewResourcePool(size int) *ResourcePool {
 pool := &ResourcePool{
  resources: make(chan Resource, size),
 }

 for i := 0; i < size; i++ {
  pool.resources <- Resource{
   ID:   fmt.Sprintf("resource_%d", i),
   Name: fmt.Sprintf("Resource %d", i),
  }
 }

 return pool
}

// Acquire 获取资源
func (rp *ResourcePool) Acquire(ctx context.Context) (Resource, error) {
 select {
 case resource := <-rp.resources:
  return resource, nil
 case <-ctx.Done():
  return Resource{}, ctx.Err()
 }
}

// Release 释放资源
func (rp *ResourcePool) Release(resource Resource) {
 select {
 case rp.resources <- resource:
 default:
  // 资源池已满
 }
}

// ConcurrentLimiter 并发限制器
type ConcurrentLimiter struct {
 semaphore chan struct{}
}

// NewConcurrentLimiter 创建并发限制器
func NewConcurrentLimiter(limit int) *ConcurrentLimiter {
 return &ConcurrentLimiter{
  semaphore: make(chan struct{}, limit),
 }
}

// Acquire 获取许可
func (cl *ConcurrentLimiter) Acquire(ctx context.Context) error {
 select {
 case cl.semaphore <- struct{}{}:
  return nil
 case <-ctx.Done():
  return ctx.Err()
 }
}

// Release 释放许可
func (cl *ConcurrentLimiter) Release() {
 select {
 case <-cl.semaphore:
 default:
 }
}

// Execute 在限制下执行
func (cl *ConcurrentLimiter) Execute(ctx context.Context, fn func() error) error {
 if err := cl.Acquire(ctx); err != nil {
  return err
 }
 defer cl.Release()

 return fn()
}

// WorkStealingPool 工作窃取池
type WorkStealingPool struct {
 workers int
 tasks   chan func()
 wg      sync.WaitGroup
 active  int32
}

// NewWorkStealingPool 创建工作窃取池
func NewWorkStealingPool(workers int) *WorkStealingPool {
 return &WorkStealingPool{
  workers: workers,
  tasks:   make(chan func(), 100),
 }
}

// Start 启动工作池
func (wsp *WorkStealingPool) Start() {
 for i := 0; i < wsp.workers; i++ {
  wsp.wg.Add(1)
  go func(workerID int) {
   defer wsp.wg.Done()

   for task := range wsp.tasks {
    atomic.AddInt32(&wsp.active, 1)
    log.Printf("[工作池] Worker %d 执行任务", workerID)
    task()
    atomic.AddInt32(&wsp.active, -1)
   }
  }(i + 1)
 }
}

// Submit 提交任务
func (wsp *WorkStealingPool) Submit(task func()) {
 wsp.tasks <- task
}

// Stop 停止工作池
func (wsp *WorkStealingPool) Stop() {
 close(wsp.tasks)
 wsp.wg.Wait()
}

func main() {
 ctx := context.Background()

 // 协作式调度器示例
 fmt.Println("=== 协作式调度器示例 ===")
 scheduler := NewCooperativeScheduler(1 * time.Second)

 scheduler.AddTask("high_priority", 10, func(ctx context.Context) error {
  fmt.Println("→ 执行高优先级任务")
  time.Sleep(100 * time.Millisecond)
  return nil
 })

 scheduler.AddTask("medium_priority", 5, func(ctx context.Context) error {
  fmt.Println("→ 执行中优先级任务")
  time.Sleep(100 * time.Millisecond)
  return nil
 })

 scheduler.AddTask("low_priority", 1, func(ctx context.Context) error {
  fmt.Println("→ 执行低优先级任务")
  time.Sleep(100 * time.Millisecond)
  return nil
 })

 if err := scheduler.Run(ctx); err != nil {
  log.Fatal(err)
 }

 // 并发限制器示例
 fmt.Println("\n=== 并发限制器示例 ===")
 limiter := NewConcurrentLimiter(2)

 var wg sync.WaitGroup
 for i := 0; i < 5; i++ {
  wg.Add(1)
  go func(id int) {
   defer wg.Done()

   err := limiter.Execute(ctx, func() error {
    fmt.Printf("→ 任务 %d 执行中\n", id)
    time.Sleep(200 * time.Millisecond)
    return nil
   })

   if err != nil {
    log.Printf("任务 %d 错误: %v", id, err)
   }
  }(i)
 }

 wg.Wait()
 fmt.Println("✓ 所有任务完成")
}
```

#### 反例说明

**错误实现1：任务被中断**

```go
// ❌ 错误：任务可能被中断
func WrongInterruptedTask() {
 for _, task := range tasks {
  go func(t Task) {
   // 执行部分工作
   partialWork()
   // 错误：任务可能被其他任务中断
  }(task)
 }
}
```

**问题**：任务应该原子性完成，不应该被中断。

**错误实现2：资源竞争**

```go
// ❌ 错误：资源竞争
var counter int

func WrongResourceRace() {
 for _, task := range tasks {
  go func(t Task) {
   counter++  // 竞态条件
  }(task)
 }
}
```

**问题**：多个任务同时访问共享资源，需要使用同步机制。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 协作式多任务 | 多个任务协作执行 |
| 资源限制 | 限制并发数以保护资源 |
| 优先级调度 | 按优先级交错执行任务 |
| 工作窃取 | 动态负载均衡 |
| 批处理 | 批量任务交错执行 |

#### 与其他模式关系

- **与并行分支**：交错路由是并行分支的协作版本
- **与多实例模式**：可以使用交错路由实现多实例
- **与同步模式**：交错路由内部需要同步

---

### 5.3 里程碑（Milestone）

#### 概念定义

**里程碑模式**表示工作流中的关键检查点，用于验证某些条件是否满足，只有满足条件才能继续执行。这是实现流程控制的验证模式。

**形式化定义**：给定条件 c 和后续任务 t，里程碑模式定义检查点 m，使得 t 只有在 c 满足时才能执行。

#### BPMN描述

在BPMN中，里程碑通过**中间事件**或**条件顺序流**实现：

```
[前置任务] → [里程碑/条件检查] → [后续任务]
```

- 中间事件作为里程碑
- 条件不满足时流程暂停
- 条件满足后继续执行

#### 流程图

```
┌─────────┐     ┌─────────────┐     ┌─────────┐
│  任务A  │────→│   里程碑    │────→│  任务B  │
└─────────┘     │  (条件检查) │     └─────────┘
                └──────┬──────┘
                       │
                  [条件不满足]
                       │
                       ▼
                ┌─────────────┐
                │   等待/暂停 │
                └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"
)

// Milestone 里程碑
type Milestone struct {
 name      string
 condition func(ctx context.Context) (bool, error)
 maxWait   time.Duration
}

// NewMilestone 创建里程碑
func NewMilestone(name string, condition func(ctx context.Context) (bool, error), maxWait time.Duration) *Milestone {
 return &Milestone{
  name:      name,
  condition: condition,
  maxWait:   maxWait,
 }
}

// Wait 等待条件满足
func (m *Milestone) Wait(ctx context.Context) error {
 log.Printf("[里程碑] 到达里程碑: %s", m.name)

 start := time.Now()
 ticker := time.NewTicker(100 * time.Millisecond)
 defer ticker.Stop()

 for {
  select {
  case <-ctx.Done():
   return ctx.Err()

  case <-ticker.C:
   satisfied, err := m.condition(ctx)
   if err != nil {
    return fmt.Errorf("条件检查失败: %w", err)
   }

   if satisfied {
    log.Printf("[里程碑] 条件满足，通过里程碑: %s (耗时: %v)", m.name, time.Since(start))
    return nil
   }

   if time.Since(start) > m.maxWait {
    return fmt.Errorf("等待里程碑超时: %s", m.name)
   }
  }
 }
}

// WorkflowWithMilestones 带里程碑的工作流
type WorkflowWithMilestones struct {
 milestones []Milestone
 steps      []WorkflowStep
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
 Name      string
 Action    func(ctx context.Context) error
 Milestone *Milestone
}

// NewWorkflowWithMilestones 创建工作流
func NewWorkflowWithMilestones() *WorkflowWithMilestones {
 return &WorkflowWithMilestones{
  milestones: make([]Milestone, 0),
  steps:      make([]WorkflowStep, 0),
 }
}

// AddStep 添加步骤
func (w *WorkflowWithMilestones) AddStep(name string, action func(ctx context.Context) error, milestone *Milestone) {
 w.steps = append(w.steps, WorkflowStep{
  Name:      name,
  Action:    action,
  Milestone: milestone,
 })
}

// Execute 执行工作流
func (w *WorkflowWithMilestones) Execute(ctx context.Context) error {
 for i, step := range w.steps {
  log.Printf("[工作流] 执行步骤 %d/%d: %s", i+1, len(w.steps), step.Name)

  // 检查里程碑
  if step.Milestone != nil {
   if err := step.Milestone.Wait(ctx); err != nil {
    return fmt.Errorf("步骤 '%s' 的里程碑检查失败: %w", step.Name, err)
   }
  }

  // 执行步骤
  if err := step.Action(ctx); err != nil {
   return fmt.Errorf("步骤 '%s' 执行失败: %w", step.Name, err)
  }
 }

 return nil
}

// 使用示例
func main() {
 workflow := NewWorkflowWithMilestones()

 // 模拟状态
 paymentConfirmed := false

 // 添加带里程碑的步骤
 workflow.AddStep("创建订单", func(ctx context.Context) error {
  fmt.Println("→ 创建订单")
  return nil
 }, nil)

 workflow.AddStep("等待支付", func(ctx context.Context) error {
  fmt.Println("→ 等待支付完成")
  return nil
 }, NewMilestone("支付确认", func(ctx context.Context) (bool, error) {
  return paymentConfirmed, nil
 }, 5*time.Second))

 workflow.AddStep("处理订单", func(ctx context.Context) error {
  fmt.Println("→ 处理订单")
  return nil
 }, nil)

 // 异步确认支付
 go func() {
  time.Sleep(1 * time.Second)
  paymentConfirmed = true
  fmt.Println("→ 支付已确认")
 }()

 ctx := context.Background()
 if err := workflow.Execute(ctx); err != nil {
  log.Fatal(err)
 }

 fmt.Println("\n✓ 工作流完成")
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "sync/atomic"
 "time"
)

// ApprovalWorkflow 审批工作流
type ApprovalWorkflow struct {
 approvals map[string]bool
 mu        sync.RWMutex
}

// NewApprovalWorkflow 创建审批工作流
func NewApprovalWorkflow() *ApprovalWorkflow {
 return &ApprovalWorkflow{
  approvals: make(map[string]bool),
 }
}

// RequestApproval 请求审批
func (aw *ApprovalWorkflow) RequestApproval(ctx context.Context, approvers []string, task func() error) error {
 // 创建里程碑等待所有审批
 milestone := NewMilestone("审批完成", func(ctx context.Context) (bool, error) {
  aw.mu.RLock()
  defer aw.mu.RUnlock()

  for _, approver := range approvers {
   if !aw.approvals[approver] {
    return false, nil
   }
  }
  return true, nil
 }, 10*time.Second)

 fmt.Printf("等待 %d 位审批者审批...\n", len(approvers))

 if err := milestone.Wait(ctx); err != nil {
  return err
 }

 return task()
}

// Approve 审批通过
func (aw *ApprovalWorkflow) Approve(approver string) {
 aw.mu.Lock()
 defer aw.mu.Unlock()
 aw.approvals[approver] = true
 fmt.Printf("→ %s 已审批\n", approver)
}

// CheckpointManager 检查点管理器
type CheckpointManager struct {
 checkpoints map[string]*Checkpoint
 mu          sync.RWMutex
}

// Checkpoint 检查点
type Checkpoint struct {
 Name      string
 Condition func() bool
 Completed bool
}

// NewCheckpointManager 创建检查点管理器
func NewCheckpointManager() *CheckpointManager {
 return &CheckpointManager{
  checkpoints: make(map[string]*Checkpoint),
 }
}

// Register 注册检查点
func (cm *CheckpointManager) Register(name string, condition func() bool) {
 cm.mu.Lock()
 defer cm.mu.Unlock()

 cm.checkpoints[name] = &Checkpoint{
  Name:      name,
  Condition: condition,
  Completed: false,
 }
}

// Check 检查检查点
func (cm *CheckpointManager) Check(name string) bool {
 cm.mu.RLock()
 defer cm.mu.RUnlock()

 checkpoint, ok := cm.checkpoints[name]
 if !ok {
  return false
 }

 if checkpoint.Completed {
  return true
 }

 if checkpoint.Condition() {
  checkpoint.Completed = true
  return true
 }

 return false
}

// WaitForCheckpoint 等待检查点
func (cm *CheckpointManager) WaitForCheckpoint(ctx context.Context, name string, timeout time.Duration) error {
 start := time.Now()

 for {
  if cm.Check(name) {
   return nil
  }

  if time.Since(start) > timeout {
   return fmt.Errorf("等待检查点 '%s' 超时", name)
  }

  select {
  case <-ctx.Done():
   return ctx.Err()
  case <-time.After(100 * time.Millisecond):
   // 继续检查
  }
 }
}

// ProgressTracker 进度追踪器
type ProgressTracker struct {
 total     int32
 completed int32
}

// NewProgressTracker 创建进度追踪器
func NewProgressTracker(total int) *ProgressTracker {
 return &ProgressTracker{
  total: int32(total),
 }
}

// Complete 完成一项
func (pt *ProgressTracker) Complete() {
 atomic.AddInt32(&pt.completed, 1)
}

// GetProgress 获取进度
func (pt *ProgressTracker) GetProgress() (completed, total int, percentage float64) {
 c := atomic.LoadInt32(&pt.completed)
 t := atomic.LoadInt32(&pt.total)
 return int(c), int(t), float64(c) / float64(t) * 100
}

// IsComplete 是否完成
func (pt *ProgressTracker) IsComplete() bool {
 return atomic.LoadInt32(&pt.completed) >= atomic.LoadInt32(&pt.total)
}

// WaitForCompletion 等待完成
func (pt *ProgressTracker) WaitForCompletion(ctx context.Context, timeout time.Duration) error {
 start := time.Now()

 for !pt.IsComplete() {
  if time.Since(start) > timeout {
   return fmt.Errorf("等待完成超时")
  }

  select {
  case <-ctx.Done():
   return ctx.Err()
  case <-time.After(100 * time.Millisecond):
   completed, total, pct := pt.GetProgress()
   fmt.Printf("进度: %d/%d (%.1f%%)\n", completed, total, pct)
  }
 }

 return nil
}

func main() {
 ctx := context.Background()

 // 审批工作流示例
 fmt.Println("=== 审批工作流示例 ===")
 approval := NewApprovalWorkflow()

 approvers := []string{"manager", "director", "cto"}

 // 模拟异步审批
 go func() {
  time.Sleep(500 * time.Millisecond)
  approval.Approve("manager")

  time.Sleep(500 * time.Millisecond)
  approval.Approve("director")

  time.Sleep(500 * time.Millisecond)
  approval.Approve("cto")
 }()

 err := approval.RequestApproval(ctx, approvers, func() error {
  fmt.Println("→ 所有审批通过，执行任务")
  return nil
 })

 if err != nil {
  log.Fatal(err)
 }

 // 进度追踪示例
 fmt.Println("\n=== 进度追踪示例 ===")
 tracker := NewProgressTracker(5)

 // 模拟异步任务
 go func() {
  for i := 0; i < 5; i++ {
   time.Sleep(300 * time.Millisecond)
   tracker.Complete()
  }
 }()

 if err := tracker.WaitForCompletion(ctx, 5*time.Second); err != nil {
  log.Fatal(err)
 }

 fmt.Println("✓ 所有任务完成")
}
```

#### 反例说明

**错误实现1：没有超时机制**

```go
// ❌ 错误：没有超时，可能永久阻塞
func WrongNoTimeout() {
 for {
  if checkCondition() {
   break
  }
  time.Sleep(100 * time.Millisecond)
 }
}
```

**问题**：没有超时机制可能导致永久阻塞。

**错误实现2：忙等待**

```go
// ❌ 错误：忙等待浪费CPU
func WrongBusyWait() {
 for !checkCondition() {
  // 错误：没有sleep，忙等待
 }
}
```

**问题**：忙等待浪费CPU资源，应该使用定时检查或事件驱动。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 审批流程 | 等待审批通过 |
| 依赖检查 | 等待依赖满足 |
| 进度追踪 | 追踪任务进度 |
| 条件等待 | 等待条件满足 |
| 流程控制 | 控制流程执行 |

#### 与其他模式关系

- **与延迟选择**：里程碑可以基于事件触发
- **与同步模式**：里程碑等待条件满足
- **与任意循环**：里程碑可以在循环中使用

---

## 6. 取消与补偿模式

### 6.1 取消任务（Cancel Activity）

#### 概念定义

**取消任务模式**允许在工作流执行过程中取消一个或多个正在执行的活动。这是实现异常处理、超时控制的核心模式。

**形式化定义**：给定活动集合 A 和取消条件 c，取消任务模式定义取消操作，使得当 c 满足时，活动 a ∈ A 被终止。

#### BPMN描述

在BPMN中，取消任务通过**取消事件**或**边界事件**实现：

```
                    ┌→ [取消事件] → [取消处理]
[前置任务] → [任务] │
                    └→ [正常完成]
```

- 边界取消事件附加到任务
- 取消事件触发时终止任务
- 执行取消处理逻辑

#### 流程图

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  开始   │────→│  任务   │────→│  完成   │
└─────────┘     └────┬────┘     └─────────┘
                     │
                [取消边界]
                     │
                     ▼
               ┌─────────────┐
               │  取消处理   │
               └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// CancellableTask 可取消任务
type CancellableTask struct {
 name   string
 cancel chan struct{}
 done   chan error
}

// NewCancellableTask 创建可取消任务
func NewCancellableTask(name string) *CancellableTask {
 return &CancellableTask{
  name:   name,
  cancel: make(chan struct{}),
  done:   make(chan error, 1),
 }
}

// Execute 执行任务
func (ct *CancellableTask) Execute(ctx context.Context, action func(ctx context.Context) error) {
 go func() {
  // 创建可取消的上下文
  taskCtx, cancel := context.WithCancel(ctx)
  defer cancel()

  // 监听取消信号
  go func() {
   select {
   case <-ct.cancel:
    cancel()
   case <-taskCtx.Done():
   }
  }()

  ct.done <- action(taskCtx)
 }()
}

// Cancel 取消任务
func (ct *CancellableTask) Cancel() {
 close(ct.cancel)
}

// Wait 等待任务完成
func (ct *CancellableTask) Wait() error {
 return <-ct.done
}

// TaskManager 任务管理器
type TaskManager struct {
 tasks map[string]*CancellableTask
 mu    sync.RWMutex
}

// NewTaskManager 创建任务管理器
func NewTaskManager() *TaskManager {
 return &TaskManager{
  tasks: make(map[string]*CancellableTask),
 }
}

// Register 注册任务
func (tm *TaskManager) Register(name string, action func(ctx context.Context) error) {
 tm.mu.Lock()
 defer tm.mu.Unlock()

 task := NewCancellableTask(name)
 tm.tasks[name] = task

 ctx := context.Background()
 task.Execute(ctx, action)
}

// Cancel 取消任务
func (tm *TaskManager) Cancel(name string) error {
 tm.mu.RLock()
 task, ok := tm.tasks[name]
 tm.mu.RUnlock()

 if !ok {
  return fmt.Errorf("任务 '%s' 不存在", name)
 }

 log.Printf("[取消任务] 取消任务: %s", name)
 task.Cancel()
 return nil
}

// CancelAll 取消所有任务
func (tm *TaskManager) CancelAll() {
 tm.mu.RLock()
 tasks := make([]*CancellableTask, 0, len(tm.tasks))
 for _, task := range tm.tasks {
  tasks = append(tasks, task)
 }
 tm.mu.RUnlock()

 for _, task := range tasks {
  task.Cancel()
 }
}

// 使用示例
func main() {
 tm := NewTaskManager()

 // 注册任务
 tm.Register("long_task", func(ctx context.Context) error {
  for i := 0; i < 10; i++ {
   select {
   case <-ctx.Done():
    fmt.Println("→ 长任务被取消")
    return ctx.Err()
   default:
    fmt.Printf("→ 长任务执行中: %d/10\n", i+1)
    time.Sleep(200 * time.Millisecond)
   }
  }
  return nil
 })

 // 3秒后取消任务
 go func() {
  time.Sleep(600 * time.Millisecond)
  tm.Cancel("long_task")
 }()

 // 等待任务完成
 tm.mu.RLock()
 task := tm.tasks["long_task"]
 tm.mu.RUnlock()

 err := task.Wait()
 if err != nil {
  log.Printf("任务执行结果: %v", err)
 }
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// TimeoutTask 带超时的任务
type TimeoutTask struct {
 name    string
 timeout time.Duration
 action  func(ctx context.Context) error
}

// NewTimeoutTask 创建带超时的任务
func NewTimeoutTask(name string, timeout time.Duration, action func(ctx context.Context) error) *TimeoutTask {
 return &TimeoutTask{
  name:    name,
  timeout: timeout,
  action:  action,
 }
}

// Run 运行任务
func (tt *TimeoutTask) Run(ctx context.Context) error {
 ctx, cancel := context.WithTimeout(ctx, tt.timeout)
 defer cancel()

 done := make(chan error, 1)

 go func() {
  done <- tt.action(ctx)
 }()

 select {
 case err := <-done:
  return err
 case <-ctx.Done():
  return fmt.Errorf("任务 '%s' 超时", tt.name)
 }
}

// CancellationToken 取消令牌
type CancellationToken struct {
 ctx    context.Context
 cancel context.CancelFunc
 mu     sync.Mutex
}

// NewCancellationToken 创建取消令牌
func NewCancellationToken() *CancellationToken {
 ctx, cancel := context.WithCancel(context.Background())
 return &CancellationToken{
  ctx:    ctx,
  cancel: cancel,
 }
}

// Context 获取上下文
func (ct *CancellationToken) Context() context.Context {
 return ct.ctx
}

// Cancel 取消
func (ct *CancellationToken) Cancel() {
 ct.mu.Lock()
 defer ct.mu.Unlock()

 if ct.cancel != nil {
  ct.cancel()
  ct.cancel = nil
 }
}

// IsCancelled 是否已取消
func (ct *CancellationToken) IsCancelled() bool {
 select {
 case <-ct.ctx.Done():
  return true
 default:
  return false
 }
}

// CooperativeTask 协作式任务
type CooperativeTask struct {
 name      string
 cancelled int32
}

// NewCooperativeTask 创建协作式任务
func NewCooperativeTask(name string) *CooperativeTask {
 return &CooperativeTask{name: name}
}

// Execute 执行任务
func (ct *CooperativeTask) Execute(ctx context.Context, steps []func() error) error {
 for i, step := range steps {
  select {
  case <-ctx.Done():
   fmt.Printf("→ 任务 '%s' 在第 %d 步被取消\n", ct.name, i+1)
   return ctx.Err()
  default:
  }

  if err := step(); err != nil {
   return err
  }
 }

 return nil
}

// WorkflowCancellation 工作流取消
type WorkflowCancellation struct {
 tasks []CancellableFunc
 mu    sync.Mutex
}

// CancellableFunc 可取消函数
type CancellableFunc struct {
 Name   string
 Cancel context.CancelFunc
}

// NewWorkflowCancellation 创建工作流取消管理器
func NewWorkflowCancellation() *WorkflowCancellation {
 return &WorkflowCancellation{
  tasks: make([]CancellableFunc, 0),
 }
}

// AddTask 添加任务
func (wc *WorkflowCancellation) AddTask(name string, action func(ctx context.Context) error) {
 ctx, cancel := context.WithCancel(context.Background())

 wc.mu.Lock()
 wc.tasks = append(wc.tasks, CancellableFunc{Name: name, Cancel: cancel})
 wc.mu.Unlock()

 go func() {
  if err := action(ctx); err != nil {
   log.Printf("任务 '%s' 错误: %v", name, err)
  }
 }()
}

// CancelTask 取消指定任务
func (wc *WorkflowCancellation) CancelTask(name string) {
 wc.mu.Lock()
 defer wc.mu.Unlock()

 for _, task := range wc.tasks {
  if task.Name == name {
   task.Cancel()
   log.Printf("取消任务: %s", name)
   return
  }
 }
}

// CancelAll 取消所有任务
func (wc *WorkflowCancellation) CancelAll() {
 wc.mu.Lock()
 defer wc.mu.Unlock()

 for _, task := range wc.tasks {
  task.Cancel()
 }

 log.Println("取消所有任务")
}

func main() {
 ctx := context.Background()

 // 超时任务示例
 fmt.Println("=== 超时任务示例 ===")
 timeoutTask := NewTimeoutTask("data_fetch", 500*time.Millisecond, func(ctx context.Context) error {
  for i := 0; i < 10; i++ {
   select {
   case <-ctx.Done():
    return ctx.Err()
   default:
    fmt.Printf("→ 获取数据中: %d/10\n", i+1)
    time.Sleep(100 * time.Millisecond)
   }
  }
  return nil
 })

 err := timeoutTask.Run(ctx)
 if err != nil {
  log.Printf("任务结果: %v", err)
 }

 // 取消令牌示例
 fmt.Println("\n=== 取消令牌示例 ===")
 token := NewCancellationToken()

 go func() {
  time.Sleep(300 * time.Millisecond)
  token.Cancel()
 }()

 select {
 case <-token.Context().Done():
  fmt.Println("→ 令牌已取消")
 case <-time.After(1 * time.Second):
  fmt.Println("→ 等待超时")
 }

 // 工作流取消示例
 fmt.Println("\n=== 工作流取消示例 ===")
 wf := NewWorkflowCancellation()

 wf.AddTask("task1", func(ctx context.Context) error {
  for i := 0; i < 5; i++ {
   select {
   case <-ctx.Done():
    fmt.Println("→ task1 被取消")
    return ctx.Err()
   default:
    fmt.Println("→ task1 执行中")
    time.Sleep(200 * time.Millisecond)
   }
  }
  return nil
 })

 wf.AddTask("task2", func(ctx context.Context) error {
  for i := 0; i < 5; i++ {
   select {
   case <-ctx.Done():
    fmt.Println("→ task2 被取消")
    return ctx.Err()
   default:
    fmt.Println("→ task2 执行中")
    time.Sleep(200 * time.Millisecond)
   }
  }
  return nil
 })

 // 500ms后取消所有任务
 time.Sleep(500 * time.Millisecond)
 wf.CancelAll()

 time.Sleep(200 * time.Millisecond)
 fmt.Println("✓ 工作流取消完成")
}
```

#### 反例说明

**错误实现1：没有检查取消信号**

```go
// ❌ 错误：没有检查取消信号
func WrongNoCancelCheck(ctx context.Context) {
 for i := 0; i < 100; i++ {
  process()  // 错误：即使ctx取消也会继续
 }
}
```

**问题**：没有检查取消信号，任务无法响应取消。

**错误实现2：强制终止goroutine**

```go
// ❌ 错误：Go没有强制终止goroutine的方法
func WrongForceKill() {
 go longRunningTask()
 // 错误：无法强制终止goroutine
}
```

**问题**：Go没有强制终止goroutine的方法，必须使用协作式取消。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 超时控制 | 任务执行超时后取消 |
| 用户取消 | 用户主动取消操作 |
| 错误处理 | 出错时取消相关任务 |
| 资源释放 | 取消后释放资源 |
| 流程终止 | 终止整个工作流 |

#### 与其他模式关系

- **与取消案例**：取消任务是取消案例的基础
- **与补偿模式**：取消后可能需要补偿
- **与超时模式**：超时是取消的常见触发条件

---

### 6.2 取消案例（Cancel Case）

#### 概念定义

**取消案例模式**允许取消整个工作流实例，包括所有正在执行的活动和子流程。这是实现工作流终止的核心模式。

**形式化定义**：给定工作流实例 W，取消案例模式定义取消操作，使得 W 中的所有活动都被终止。

#### BPMN描述

在BPMN中，取消案例通过**取消结束事件**实现：

```
[任务] → [取消事件] → (终止整个流程)
```

- 取消结束事件终止整个流程
- 所有活动被强制终止
- 可选执行清理逻辑

#### 流程图

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  开始   │────→│  任务A  │────→│  任务B  │
└─────────┘     └────┬────┘     └────┬────┘
                     │               │
                     │          [错误]
                     │               │
                     │               ▼
                     │         ┌─────────────┐
                     │         │  取消案例   │
                     │         │  (终止全部) │
                     │         └─────────────┘
                     │
                     └────────→ [正常结束]
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// WorkflowInstance 工作流实例
type WorkflowInstance struct {
 id       string
 ctx      context.Context
 cancel   context.CancelFunc
 tasks    map[string]*Task
 mu       sync.RWMutex
 status   WorkflowStatus
}

// WorkflowStatus 工作流状态
type WorkflowStatus int

const (
 StatusRunning WorkflowStatus = iota
 StatusCompleted
 StatusCancelled
 StatusFailed
)

// Task 任务
type Task struct {
 Name   string
 Action func(ctx context.Context) error
}

// NewWorkflowInstance 创建工作流实例
func NewWorkflowInstance(id string) *WorkflowInstance {
 ctx, cancel := context.WithCancel(context.Background())
 return &WorkflowInstance{
  id:     id,
  ctx:    ctx,
  cancel: cancel,
  tasks:  make(map[string]*Task),
  status: StatusRunning,
 }
}

// AddTask 添加任务
func (wi *WorkflowInstance) AddTask(name string, action func(ctx context.Context) error) {
 wi.mu.Lock()
 defer wi.mu.Unlock()

 wi.tasks[name] = &Task{Name: name, Action: action}
}

// Start 启动工作流
func (wi *WorkflowInstance) Start() error {
 var wg sync.WaitGroup
 errChan := make(chan error, len(wi.tasks))

 wi.mu.RLock()
 tasks := make([]*Task, 0, len(wi.tasks))
 for _, task := range wi.tasks {
  tasks = append(tasks, task)
 }
 wi.mu.RUnlock()

 for _, task := range tasks {
  wg.Add(1)
  go func(t *Task) {
   defer wg.Done()

   log.Printf("[工作流 %s] 启动任务: %s", wi.id, t.Name)

   if err := t.Action(wi.ctx); err != nil {
    errChan <- fmt.Errorf("任务 '%s' 失败: %w", t.Name, err)
   }
  }(task)
 }

 // 等待所有任务完成或被取消
 go func() {
  wg.Wait()
  close(errChan)
 }()

 // 检查错误
 for err := range errChan {
  if err != nil {
   wi.status = StatusFailed
   return err
  }
 }

 if wi.status != StatusCancelled {
  wi.status = StatusCompleted
 }

 return nil
}

// Cancel 取消工作流
func (wi *WorkflowInstance) Cancel() {
 wi.mu.Lock()
 defer wi.mu.Unlock()

 if wi.status == StatusRunning {
  log.Printf("[工作流 %s] 取消工作流", wi.id)
  wi.status = StatusCancelled
  wi.cancel()
 }
}

// GetStatus 获取状态
func (wi *WorkflowInstance) GetStatus() WorkflowStatus {
 wi.mu.RLock()
 defer wi.mu.RUnlock()
 return wi.status
}

// WorkflowManager 工作流管理器
type WorkflowManager struct {
 workflows map[string]*WorkflowInstance
 mu        sync.RWMutex
}

// NewWorkflowManager 创建工作流管理器
func NewWorkflowManager() *WorkflowManager {
 return &WorkflowManager{
  workflows: make(map[string]*WorkflowInstance),
 }
}

// CreateWorkflow 创建工作流
func (wm *WorkflowManager) CreateWorkflow(id string) *WorkflowInstance {
 wm.mu.Lock()
 defer wm.mu.Unlock()

 workflow := NewWorkflowInstance(id)
 wm.workflows[id] = workflow
 return workflow
}

// CancelWorkflow 取消工作流
func (wm *WorkflowManager) CancelWorkflow(id string) error {
 wm.mu.RLock()
 workflow, ok := wm.workflows[id]
 wm.mu.RUnlock()

 if !ok {
  return fmt.Errorf("工作流 '%s' 不存在", id)
 }

 workflow.Cancel()
 return nil
}

// CancelAll 取消所有工作流
func (wm *WorkflowManager) CancelAll() {
 wm.mu.RLock()
 workflows := make([]*WorkflowInstance, 0, len(wm.workflows))
 for _, wf := range wm.workflows {
  workflows = append(workflows, wf)
 }
 wm.mu.RUnlock()

 for _, wf := range workflows {
  wf.Cancel()
 }
}

// 使用示例
func main() {
 manager := NewWorkflowManager()

 // 创建工作流
 workflow := manager.CreateWorkflow("order_001")

 // 添加任务
 workflow.AddTask("validate", func(ctx context.Context) error {
  for i := 0; i < 5; i++ {
   select {
   case <-ctx.Done():
    fmt.Println("→ 验证任务被取消")
    return ctx.Err()
   default:
    fmt.Println("→ 验证中...")
    time.Sleep(200 * time.Millisecond)
   }
  }
  return nil
 })

 workflow.AddTask("process", func(ctx context.Context) error {
  for i := 0; i < 5; i++ {
   select {
   case <-ctx.Done():
    fmt.Println("→ 处理任务被取消")
    return ctx.Err()
   default:
    fmt.Println("→ 处理中...")
    time.Sleep(200 * time.Millisecond)
   }
  }
  return nil
 })

 // 启动工作流
 go func() {
  if err := workflow.Start(); err != nil {
   log.Printf("工作流错误: %v", err)
  }
 }()

 // 600ms后取消工作流
 time.Sleep(600 * time.Millisecond)
 manager.CancelWorkflow("order_001")

 time.Sleep(200 * time.Millisecond)
 fmt.Printf("\n✓ 工作流状态: %v\n", workflow.GetStatus())
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// Saga  Saga模式实现
type Saga struct {
 steps       []SagaStep
 compensations []Compensation
 mu          sync.Mutex
}

// SagaStep Saga步骤
type SagaStep struct {
 Name      string
 Action    func(ctx context.Context) error
 Compensate func(ctx context.Context) error
}

// Compensation 补偿
type Compensation struct {
 StepName string
 Action   func(ctx context.Context) error
}

// NewSaga 创建Saga
func NewSaga() *Saga {
 return &Saga{
  steps:         make([]SagaStep, 0),
  compensations: make([]Compensation, 0),
 }
}

// AddStep 添加步骤
func (s *Saga) AddStep(name string, action func(ctx context.Context) error, compensate func(ctx context.Context) error) {
 s.steps = append(s.steps, SagaStep{
  Name:       name,
  Action:     action,
  Compensate: compensate,
 })
}

// Execute 执行Saga
func (s *Saga) Execute(ctx context.Context) error {
 completedSteps := make([]string, 0)

 for _, step := range s.steps {
  log.Printf("[Saga] 执行步骤: %s", step.Name)

  if err := step.Action(ctx); err != nil {
   log.Printf("[Saga] 步骤 '%s' 失败: %v", step.Name, err)

   // 执行补偿
   s.compensate(ctx, completedSteps)
   return err
  }

  completedSteps = append(completedSteps, step.Name)
 }

 return nil
}

// compensate 执行补偿
func (s *Saga) compensate(ctx context.Context, completedSteps []string) {
 log.Println("[Saga] 开始补偿...")

 // 逆序执行补偿
 for i := len(completedSteps) - 1; i >= 0; i-- {
  stepName := completedSteps[i]

  for _, step := range s.steps {
   if step.Name == stepName && step.Compensate != nil {
    log.Printf("[Saga] 补偿步骤: %s", stepName)
    if err := step.Compensate(ctx); err != nil {
     log.Printf("[Saga] 补偿 '%s' 失败: %v", stepName, err)
    }
   }
  }
 }
}

// TransactionManager 事务管理器
type TransactionManager struct {
 transactions map[string]*Transaction
 mu           sync.RWMutex
}

// Transaction 事务
type Transaction struct {
 ID      string
 Ctx     context.Context
 Cancel  context.CancelFunc
 Status  TransactionStatus
}

// TransactionStatus 事务状态
type TransactionStatus int

const (
 TxPending TransactionStatus = iota
 TxCommitted
 TxRolledBack
)

// NewTransactionManager 创建事务管理器
func NewTransactionManager() *TransactionManager {
 return &TransactionManager{
  transactions: make(map[string]*Transaction),
 }
}

// Begin 开始事务
func (tm *TransactionManager) Begin(id string) *Transaction {
 ctx, cancel := context.WithCancel(context.Background())

 tx := &Transaction{
  ID:     id,
  Ctx:    ctx,
  Cancel: cancel,
  Status: TxPending,
 }

 tm.mu.Lock()
 tm.transactions[id] = tx
 tm.mu.Unlock()

 return tx
}

// Rollback 回滚事务
func (tm *TransactionManager) Rollback(id string) error {
 tm.mu.RLock()
 tx, ok := tm.transactions[id]
 tm.mu.RUnlock()

 if !ok {
  return fmt.Errorf("事务 '%s' 不存在", id)
 }

 log.Printf("[事务] 回滚事务: %s", id)
 tx.Status = TxRolledBack
 tx.Cancel()

 return nil
}

// CleanupManager 清理管理器
type CleanupManager struct {
 cleanups []CleanupFunc
 mu       sync.Mutex
}

// CleanupFunc 清理函数
type CleanupFunc func() error

// NewCleanupManager 创建清理管理器
func NewCleanupManager() *CleanupManager {
 return &CleanupManager{
  cleanups: make([]CleanupFunc, 0),
 }
}

// Register 注册清理函数
func (cm *CleanupManager) Register(fn CleanupFunc) {
 cm.mu.Lock()
 defer cm.mu.Unlock()
 cm.cleanups = append(cm.cleanups, fn)
}

// Cleanup 执行清理
func (cm *CleanupManager) Cleanup() error {
 cm.mu.Lock()
 cleanups := make([]CleanupFunc, len(cm.cleanups))
 copy(cleanups, cm.cleanups)
 cm.mu.Unlock()

 var errs []error

 // 逆序执行清理
 for i := len(cleanups) - 1; i >= 0; i-- {
  if err := cleanups[i](); err != nil {
   errs = append(errs, err)
  }
 }

 if len(errs) > 0 {
  return fmt.Errorf("清理错误: %v", errs)
 }

 return nil
}

func main() {
 ctx := context.Background()

 // Saga模式示例
 fmt.Println("=== Saga模式示例 ===")
 saga := NewSaga()

 // 添加步骤和补偿
 saga.AddStep(
  "扣减库存",
  func(ctx context.Context) error {
   fmt.Println("→ 扣减库存")
   return nil
  },
  func(ctx context.Context) error {
   fmt.Println("→ 补偿：恢复库存")
   return nil
  },
 )

 saga.AddStep(
  "扣减余额",
  func(ctx context.Context) error {
   fmt.Println("→ 扣减余额")
   return nil
  },
  func(ctx context.Context) error {
   fmt.Println("→ 补偿：恢复余额")
   return nil
  },
 )

 saga.AddStep(
  "创建订单",
  func(ctx context.Context) error {
   fmt.Println("→ 创建订单")
   return fmt.Errorf("创建订单失败")
  },
  func(ctx context.Context) error {
   fmt.Println("→ 补偿：删除订单")
   return nil
  },
 )

 err := saga.Execute(ctx)
 if err != nil {
  log.Printf("Saga执行失败: %v", err)
 }

 // 清理管理器示例
 fmt.Println("\n=== 清理管理器示例 ===")
 cleanup := NewCleanupManager()

 cleanup.Register(func() error {
  fmt.Println("→ 清理资源1")
  return nil
 })

 cleanup.Register(func() error {
  fmt.Println("→ 清理资源2")
  return nil
 })

 cleanup.Register(func() error {
  fmt.Println("→ 清理资源3")
  return nil
 })

 if err := cleanup.Cleanup(); err != nil {
  log.Fatal(err)
 }

 fmt.Println("✓ 清理完成")
}
```

#### 反例说明

**错误实现1：没有清理资源**

```go
// ❌ 错误：取消后没有清理资源
func WrongNoCleanup() {
 file, _ := os.Open("data.txt")

 ctx, cancel := context.WithCancel(context.Background())
 go processFile(ctx, file)

 time.Sleep(1 * time.Second)
 cancel()
 // 错误：文件没有关闭
}
```

**问题**：取消后没有清理资源，导致资源泄漏。

**错误实现2：没有补偿机制**

```go
// ❌ 错误：取消后没有补偿
func WrongNoCompensation() {
 chargeUser()
 reserveInventory()
 createOrder()  // 如果这里失败
 // 错误：没有补偿之前的操作
}
```

**问题**：取消后没有补偿，导致数据不一致。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 分布式事务 | 分布式事务的补偿处理 |
| 工作流终止 | 终止整个工作流实例 |
| 错误恢复 | 出错后的恢复处理 |
| 资源清理 | 取消后的资源释放 |
| Saga模式 | 长事务的补偿处理 |

#### 与其他模式关系

- **与取消任务**：取消案例是取消任务的扩展
- **与补偿模式**：取消后可能需要补偿
- **与Saga模式**：Saga使用补偿实现长事务

---

## 7. 其他可判断模式

### 7.1 递归（Recursion）

#### 概念定义

**递归模式**允许活动调用自身，形成递归执行结构。这是实现分治算法、树形遍历的核心模式。

**形式化定义**：给定活动 A，递归模式允许 A 在执行过程中调用 A 自身，直到满足终止条件。

#### BPMN描述

在BPMN中，递归通过**调用活动**实现，调用活动可以调用包含自身的流程：

```
[开始] → [任务A] → [条件] →┬→ [调用自身] →┐
                            │              │
                            └→ [结束] ←────┘
```

- 调用活动调用子流程
- 子流程可以包含调用活动自身
- 需要终止条件防止无限递归

#### 流程图

```
        ┌──────────────────────────┐
        │                          │
        ▼                          │
┌─────────────┐    ┌─────────┐     │
│   递归函数  │───→│  条件   │─────┤
└─────────────┘    └────┬────┘     │
                        │          │
               ┌────────┴────────┐ │
               │                 │ │
          [终止条件]        [递归调用]│
               │                 │ │
               ▼                 │ │
          ┌─────────┐            │ │
          │  返回   │            │ │
          └─────────┘            │ │
                                 │ │
                                 ▼ │
                          ┌─────────────┐
                          │  处理+递归  │
                          └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
)

// RecursiveProcessor 递归处理器
type RecursiveProcessor struct {
 maxDepth int
}

// NewRecursiveProcessor 创建递归处理器
func NewRecursiveProcessor(maxDepth int) *RecursiveProcessor {
 return &RecursiveProcessor{maxDepth: maxDepth}
}

// Factorial 阶乘计算
func (rp *RecursiveProcessor) Factorial(n int) (int, error) {
 if n < 0 {
  return 0, fmt.Errorf("负数没有阶乘")
 }

 // 终止条件
 if n == 0 || n == 1 {
  return 1, nil
 }

 // 递归调用
 result, err := rp.Factorial(n - 1)
 if err != nil {
  return 0, err
 }

 return n * result, nil
}

// Fibonacci 斐波那契数列
func (rp *RecursiveProcessor) Fibonacci(n int) (int, error) {
 if n < 0 {
  return 0, fmt.Errorf("负数没有斐波那契数")
 }

 // 终止条件
 if n == 0 {
  return 0, nil
 }
 if n == 1 {
  return 1, nil
 }

 // 递归调用
 n1, err := rp.Fibonacci(n - 1)
 if err != nil {
  return 0, err
 }

 n2, err := rp.Fibonacci(n - 2)
 if err != nil {
  return 0, err
 }

 return n1 + n2, nil
}

// TreeTraversal 树遍历
type TreeNode struct {
 Value    int
 Children []*TreeNode
}

// DepthFirstSearch 深度优先搜索
func (rp *RecursiveProcessor) DepthFirstSearch(node *TreeNode, visit func(int)) {
 if node == nil {
  return
 }

 // 访问当前节点
 visit(node.Value)

 // 递归遍历子节点
 for _, child := range node.Children {
  rp.DepthFirstSearch(child, visit)
 }
}

// BinarySearch 二分查找
func (rp *RecursiveProcessor) BinarySearch(arr []int, target int) int {
 return rp.binarySearchHelper(arr, target, 0, len(arr)-1)
}

func (rp *RecursiveProcessor) binarySearchHelper(arr []int, target, left, right int) int {
 // 终止条件
 if left > right {
  return -1
 }

 mid := left + (right-left)/2

 if arr[mid] == target {
  return mid
 }

 if arr[mid] > target {
  return rp.binarySearchHelper(arr, target, left, mid-1)
 }

 return rp.binarySearchHelper(arr, target, mid+1, right)
}

// 使用示例
func main() {
 processor := NewRecursiveProcessor(100)

 // 阶乘示例
 fmt.Println("=== 阶乘计算 ===")
 for i := 0; i <= 10; i++ {
  result, err := processor.Factorial(i)
  if err != nil {
   log.Fatal(err)
  }
  fmt.Printf("%d! = %d\n", i, result)
 }

 // 斐波那契示例
 fmt.Println("\n=== 斐波那契数列 ===")
 for i := 0; i <= 10; i++ {
  result, err := processor.Fibonacci(i)
  if err != nil {
   log.Fatal(err)
  }
  fmt.Printf("F(%d) = %d\n", i, result)
 }

 // 树遍历示例
 fmt.Println("\n=== 树遍历 ===")
 root := &TreeNode{
  Value: 1,
  Children: []*TreeNode{
   {
    Value: 2,
    Children: []*TreeNode{
     {Value: 4},
     {Value: 5},
    },
   },
   {
    Value: 3,
    Children: []*TreeNode{
     {Value: 6},
     {Value: 7},
    },
   },
  },
 }

 fmt.Print("DFS: ")
 processor.DepthFirstSearch(root, func(v int) {
  fmt.Printf("%d ", v)
 })
 fmt.Println()

 // 二分查找示例
 fmt.Println("\n=== 二分查找 ===")
 arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
 target := 7
 index := processor.BinarySearch(arr, target)
 fmt.Printf("在 %v 中查找 %d: 索引 = %d\n", arr, target, index)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "path/filepath"
 "strings"
 "sync"
)

// FileProcessor 文件处理器
type FileProcessor struct {
 maxDepth int
}

// NewFileProcessor 创建文件处理器
func NewFileProcessor(maxDepth int) *FileProcessor {
 return &FileProcessor{maxDepth: maxDepth}
}

// ProcessDirectory 递归处理目录
func (fp *FileProcessor) ProcessDirectory(ctx context.Context, dir string, depth int, processor func(path string) error) error {
 // 检查深度限制
 if depth > fp.maxDepth {
  return fmt.Errorf("超过最大深度: %d", fp.maxDepth)
 }

 // 检查上下文取消
 select {
 case <-ctx.Done():
  return ctx.Err()
 default:
 }

 // 模拟目录处理
 fmt.Printf("%s处理目录: %s (深度: %d)\n", strings.Repeat("  ", depth), dir, depth)

 // 模拟子目录
 subDirs := []string{"subdir1", "subdir2"}

 for _, subDir := range subDirs {
  subPath := filepath.Join(dir, subDir)

  // 递归处理子目录
  if err := fp.ProcessDirectory(ctx, subPath, depth+1, processor); err != nil {
   return err
  }
 }

 return nil
}

// MergeSort 归并排序
func MergeSort(arr []int) []int {
 // 终止条件
 if len(arr) <= 1 {
  return arr
 }

 // 分治
 mid := len(arr) / 2
 left := MergeSort(arr[:mid])
 right := MergeSort(arr[mid:])

 // 合并
 return merge(left, right)
}

func merge(left, right []int) []int {
 result := make([]int, 0, len(left)+len(right))

 i, j := 0, 0
 for i < len(left) && j < len(right) {
  if left[i] <= right[j] {
   result = append(result, left[i])
   i++
  } else {
   result = append(result, right[j])
   j++
  }
 }

 result = append(result, left[i:]...)
 result = append(result, right[j:]...)

 return result
}

// QuickSort 快速排序
func QuickSort(arr []int) []int {
 if len(arr) <= 1 {
  return arr
 }

 pivot := arr[len(arr)/2]
 var left, middle, right []int

 for _, v := range arr {
  if v < pivot {
   left = append(left, v)
  } else if v == pivot {
   middle = append(middle, v)
  } else {
   right = append(right, v)
  }
 }

 left = QuickSort(left)
 right = QuickSort(right)

 return append(append(left, middle...), right...)
}

// Memoization 记忆化
type Memoization struct {
 cache map[int]int
 mu    sync.RWMutex
}

// NewMemoization 创建记忆化
func NewMemoization() *Memoization {
 return &Memoization{
  cache: make(map[int]int),
 }
}

// Fibonacci 带记忆化的斐波那契
func (m *Memoization) Fibonacci(n int) int {
 // 检查缓存
 m.mu.RLock()
 if val, ok := m.cache[n]; ok {
  m.mu.RUnlock()
  return val
 }
 m.mu.RUnlock()

 // 终止条件
 if n <= 1 {
  return n
 }

 // 递归计算
 result := m.Fibonacci(n-1) + m.Fibonacci(n-2)

 // 缓存结果
 m.mu.Lock()
 m.cache[n] = result
 m.mu.Unlock()

 return result
}

// TailRecursion 尾递归
type TailRecursion struct{}

// Factorial 尾递归阶乘
func (tr *TailRecursion) Factorial(n, acc int) int {
 if n == 0 {
  return acc
 }
 return tr.Factorial(n-1, n*acc)
}

func main() {
 // 归并排序示例
 fmt.Println("=== 归并排序 ===")
 arr := []int{64, 34, 25, 12, 22, 11, 90}
 fmt.Printf("原始: %v\n", arr)
 sorted := MergeSort(arr)
 fmt.Printf("排序后: %v\n", sorted)

 // 快速排序示例
 fmt.Println("\n=== 快速排序 ===")
 arr2 := []int{64, 34, 25, 12, 22, 11, 90}
 fmt.Printf("原始: %v\n", arr2)
 sorted2 := QuickSort(arr2)
 fmt.Printf("排序后: %v\n", sorted2)

 // 记忆化示例
 fmt.Println("\n=== 记忆化斐波那契 ===")
 memo := NewMemoization()
 for i := 0; i <= 20; i++ {
  result := memo.Fibonacci(i)
  fmt.Printf("F(%d) = %d\n", i, result)
 }

 // 尾递归示例
 fmt.Println("\n=== 尾递归阶乘 ===")
 tr := &TailRecursion{}
 for i := 0; i <= 10; i++ {
  result := tr.Factorial(i, 1)
  fmt.Printf("%d! = %d\n", i, result)
 }

 // 目录处理示例
 fmt.Println("\n=== 目录递归处理 ===")
 processor := NewFileProcessor(3)
 ctx := context.Background()

 err := processor.ProcessDirectory(ctx, "/root", 0, func(path string) error {
  fmt.Printf("处理文件: %s\n", path)
  return nil
 })

 if err != nil {
  log.Printf("目录处理错误: %v", err)
 }
}
```

#### 反例说明

**错误实现1：没有终止条件**

```go
// ❌ 错误：没有终止条件
func WrongNoBaseCase(n int) int {
 return WrongNoBaseCase(n - 1)  // 无限递归
}
```

**问题**：没有终止条件导致无限递归，最终栈溢出。

**错误实现2：重复计算**

```go
// ❌ 错误：没有记忆化，重复计算
func WrongNoMemoization(n int) int {
 if n <= 1 {
  return n
 }
 return WrongNoMemoization(n-1) + WrongNoMemoization(n-2)  // 大量重复计算
}
```

**问题**：没有记忆化导致大量重复计算，时间复杂度指数级增长。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 分治算法 | 归并排序、快速排序 |
| 树形遍历 | 文件系统遍历、DOM遍历 |
| 数学计算 | 阶乘、斐波那契数列 |
| 搜索算法 | 二分查找、深度优先搜索 |
| 图算法 | 图的遍历、最短路径 |

#### 与其他模式关系

- **与任意循环**：递归和循环可以互相转换
- **与分治模式**：递归是实现分治的基础
- **与动态规划**：递归是动态规划的基础

---

### 7.2 临时触发器（Transient Trigger）

#### 概念定义

**临时触发器模式**允许活动等待一个临时事件触发，事件只被消费一次，如果没有活动等待，事件会被丢弃。这是实现即时通知的模式。

**形式化定义**：给定活动 A 和事件 e，临时触发器模式使得 A 等待 e，e 只触发一个等待的 A，如果没有 A 等待，e 被丢弃。

#### BPMN描述

在BPMN中，临时触发器通过**消息中间事件**或**信号中间事件**实现：

```
[任务] → [中间事件] → [后续任务]
```

- 中间事件等待触发
- 事件只触发一次
- 无等待者时事件丢弃

#### 流程图

```
┌─────────┐     ┌─────────────┐     ┌─────────┐
│  任务A  │────→│ 临时触发器  │────→│  任务B  │
└─────────┘     │ (等待事件)  │     └─────────┘
                └──────┬──────┘
                       │
                  [事件到达]
                       │
                       ▼
                ┌─────────────┐
                │  触发执行   │
                └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// TransientTrigger 临时触发器
type TransientTrigger struct {
 mu        sync.Mutex
 waiting   chan struct{}
 handler   func(data interface{})
}

// NewTransientTrigger 创建临时触发器
func NewTransientTrigger(handler func(data interface{})) *TransientTrigger {
 return &TransientTrigger{
  waiting: make(chan struct{}, 1),
  handler: handler,
 }
}

// Wait 等待触发
func (tt *TransientTrigger) Wait(ctx context.Context, timeout time.Duration) (interface{}, error) {
 timer := time.NewTimer(timeout)
 defer timer.Stop()

 select {
 case <-tt.waiting:
  return nil, nil
 case <-timer.C:
  return nil, fmt.Errorf("等待超时")
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

// Trigger 触发事件
func (tt *TransientTrigger) Trigger(data interface{}) bool {
 tt.mu.Lock()
 defer tt.mu.Unlock()

 select {
 case tt.waiting <- struct{}{}:
  if tt.handler != nil {
   go tt.handler(data)
  }
  return true
 default:
  // 没有等待者，事件丢弃
  log.Println("[临时触发器] 无等待者，事件丢弃")
  return false
 }
}

// OneShotEvent 一次性事件
type OneShotEvent struct {
 ch     chan interface{}
 once   sync.Once
 closed bool
}

// NewOneShotEvent 创建一次性事件
func NewOneShotEvent() *OneShotEvent {
 return &OneShotEvent{
  ch: make(chan interface{}, 1),
 }
}

// Send 发送事件
func (ose *OneShotEvent) Send(data interface{}) bool {
 select {
 case ose.ch <- data:
  return true
 default:
  return false
 }
}

// Receive 接收事件
func (ose *OneShotEvent) Receive(ctx context.Context) (interface{}, error) {
 select {
 case data := <-ose.ch:
  return data, nil
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

// Close 关闭事件
func (ose *OneShotEvent) Close() {
 ose.once.Do(func() {
  close(ose.ch)
  ose.closed = true
 })
}

// 使用示例
func main() {
 trigger := NewTransientTrigger(func(data interface{}) {
  fmt.Printf("→ 触发处理: %v\n", data)
 })

 // 先触发，后等待（事件被丢弃）
 fmt.Println("=== 先触发，后等待 ===")
 trigger.Trigger("event1")

 ctx := context.Background()
 _, err := trigger.Wait(ctx, 500*time.Millisecond)
 if err != nil {
  log.Printf("等待结果: %v", err)
 }

 // 一次性事件示例
 fmt.Println("\n=== 一次性事件 ===")
 event := NewOneShotEvent()

 // 发送事件
 sent := event.Send("data")
 fmt.Printf("发送结果: %v\n", sent)

 // 再次发送失败
 sent = event.Send("data2")
 fmt.Printf("再次发送结果: %v\n", sent)

 // 接收事件
 data, err := event.Receive(ctx)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("接收数据: %v\n", data)
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// NotificationService 通知服务
type NotificationService struct {
 triggers map[string]*TransientTrigger
 mu       sync.RWMutex
}

// NewNotificationService 创建通知服务
func NewNotificationService() *NotificationService {
 return &NotificationService{
  triggers: make(map[string]*TransientTrigger),
 }
}

// Register 注册触发器
func (ns *NotificationService) Register(eventType string, handler func(data interface{})) {
 ns.mu.Lock()
 defer ns.mu.Unlock()

 ns.triggers[eventType] = NewTransientTrigger(handler)
}

// Notify 发送通知
func (ns *NotificationService) Notify(eventType string, data interface{}) bool {
 ns.mu.RLock()
 trigger, ok := ns.triggers[eventType]
 ns.mu.RUnlock()

 if !ok {
  return false
 }

 return trigger.Trigger(data)
}

// WaitFor 等待通知
func (ns *NotificationService) WaitFor(eventType string, ctx context.Context, timeout time.Duration) error {
 ns.mu.RLock()
 trigger, ok := ns.triggers[eventType]
 ns.mu.RUnlock()

 if !ok {
  return fmt.Errorf("未知事件类型: %s", eventType)
 }

 _, err := trigger.Wait(ctx, timeout)
 return err
}

// EventBus 事件总线
type EventBus struct {
 subscribers map[string][]chan interface{}
 mu          sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
 return &EventBus{
  subscribers: make(map[string][]chan interface{}),
 }
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType string) chan interface{} {
 eb.mu.Lock()
 defer eb.mu.Unlock()

 ch := make(chan interface{}, 1)
 eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)

 return ch
}

// Publish 发布事件
func (eb *EventBus) Publish(eventType string, data interface{}) int {
 eb.mu.Lock()
 subscribers := eb.subscribers[eventType]
 // 清空订阅者列表（临时触发器只触发一次）
 delete(eb.subscribers, eventType)
 eb.mu.Unlock()

 count := 0
 for _, ch := range subscribers {
  select {
  case ch <- data:
   count++
  default:
  }
 }

 return count
}

// CallbackRegistry 回调注册表
type CallbackRegistry struct {
 callbacks map[string]func(data interface{})
 mu        sync.RWMutex
}

// NewCallbackRegistry 创建回调注册表
func NewCallbackRegistry() *CallbackRegistry {
 return &CallbackRegistry{
  callbacks: make(map[string]func(data interface{})),
 }
}

// Register 注册回调
func (cr *CallbackRegistry) Register(id string, callback func(data interface{})) {
 cr.mu.Lock()
 defer cr.mu.Unlock()

 cr.callbacks[id] = callback
}

// Trigger 触发回调
func (cr *CallbackRegistry) Trigger(id string, data interface{}) bool {
 cr.mu.Lock()
 callback, ok := cr.callbacks[id]
 // 删除回调（一次性）
 delete(cr.callbacks, id)
 cr.mu.Unlock()

 if !ok {
  return false
 }

 go callback(data)
 return true
}

func main() {
 // 通知服务示例
 fmt.Println("=== 通知服务示例 ===")
 ns := NewNotificationService()

 ns.Register("order_complete", func(data interface{}) {
  fmt.Printf("→ 订单完成通知: %v\n", data)
 })

 // 先通知，后注册（通知被丢弃）
 ns.Notify("order_complete", "order_001")

 // 重新注册
 ns.Register("order_complete", func(data interface{}) {
  fmt.Printf("→ 订单完成通知: %v\n", data)
 })

 // 再次通知
 ns.Notify("order_complete", "order_002")

 // 事件总线示例
 fmt.Println("\n=== 事件总线示例 ===")
 eb := NewEventBus()

 // 订阅事件
 ch1 := eb.Subscribe("user_login")
 ch2 := eb.Subscribe("user_login")

 // 发布事件
 count := eb.Publish("user_login", map[string]string{"user": "alice"})
 fmt.Printf("触发 %d 个订阅者\n", count)

 // 接收事件
 select {
 case data := <-ch1:
  fmt.Printf("订阅者1收到: %v\n", data)
 case <-time.After(100 * time.Millisecond):
 }

 select {
 case data := <-ch2:
  fmt.Printf("订阅者2收到: %v\n", data)
 case <-time.After(100 * time.Millisecond):
 }

 // 再次发布（没有订阅者）
 count = eb.Publish("user_login", map[string]string{"user": "bob"})
 fmt.Printf("再次发布，触发 %d 个订阅者\n", count)

 // 回调注册表示例
 fmt.Println("\n=== 回调注册表示例 ===")
 cr := NewCallbackRegistry()

 cr.Register("callback_1", func(data interface{}) {
  fmt.Printf("→ 回调1执行: %v\n", data)
 })

 cr.Trigger("callback_1", "data1")

 // 再次触发（回调已被删除）
 triggered := cr.Trigger("callback_1", "data2")
 fmt.Printf("再次触发结果: %v\n", triggered)
}
```

#### 反例说明

**错误实现1：事件被缓存**

```go
// ❌ 错误：事件被缓存而不是丢弃
func WrongCachedEvent() {
 events := make(chan interface{}, 100)  // 错误：有缓冲

 // 事件会被缓存，不是临时触发器
}
```

**问题**：有缓冲channel会缓存事件，不符合临时触发器的语义。

**错误实现2：事件触发多个等待者**

```go
// ❌ 错误：事件触发多个等待者
func WrongMultipleTrigger() {
 for _, waiter := range waiters {
  waiter.Notify(event)  // 错误：触发了所有等待者
 }
}
```

**问题**：临时触发器只应该触发一个等待者。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 即时通知 | 即时通知等待的组件 |
| 一次性事件 | 只处理一次的事件 |
| 信号通知 | 信号量通知 |
| 回调触发 | 一次性回调触发 |
| 竞速条件 | 多个等待者竞争事件 |

#### 与其他模式关系

- **与持久触发器**：临时触发器事件丢弃，持久触发器事件保留
- **与延迟选择**：临时触发器可以实现延迟选择
- **与事件驱动**：临时触发器是事件驱动的一种形式

---

### 7.3 持久触发器（Persistent Trigger）

#### 概念定义

**持久触发器模式**允许活动等待一个持久事件触发，事件会被保留直到被消费，即使没有活动等待。这是实现可靠通知的模式。

**形式化定义**：给定活动 A 和事件 e，持久触发器模式使得 A 等待 e，e 被保留直到被 A 消费。

#### BPMN描述

在BPMN中，持久触发器通过**消息开始事件**或**条件事件**实现：

```
[消息开始事件] → [任务]
```

- 事件被持久化存储
- 活动启动时消费事件
- 事件不会被丢弃

#### 流程图

```
┌─────────────┐     ┌─────────────┐     ┌─────────┐
│  事件产生   │────→│  持久存储   │────→│  任务   │
└─────────────┘     │  (队列)     │     └─────────┘
                    └──────┬──────┘
                           │
                      [活动启动]
                           │
                           ▼
                    ┌─────────────┐
                    │  消费事件   │
                    └─────────────┘
```

#### Go实现

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// PersistentTrigger 持久触发器
type PersistentTrigger struct {
 mu     sync.Mutex
 events []interface{}
 waiter chan interface{}
}

// NewPersistentTrigger 创建持久触发器
func NewPersistentTrigger() *PersistentTrigger {
 return &PersistentTrigger{
  events: make([]interface{}, 0),
  waiter: make(chan interface{}, 1),
 }
}

// Trigger 触发事件（持久化）
func (pt *PersistentTrigger) Trigger(data interface{}) {
 pt.mu.Lock()
 defer pt.mu.Unlock()

 // 如果有等待者，直接通知
 select {
 case pt.waiter <- data:
  return
 default:
 }

 // 否则持久化存储
 pt.events = append(pt.events, data)
 log.Printf("[持久触发器] 事件持久化存储，当前队列: %d", len(pt.events))
}

// Wait 等待事件
func (pt *PersistentTrigger) Wait(ctx context.Context) (interface{}, error) {
 pt.mu.Lock()

 // 检查是否有已持久化的事件
 if len(pt.events) > 0 {
  data := pt.events[0]
  pt.events = pt.events[1:]
  pt.mu.Unlock()
  return data, nil
 }

 pt.mu.Unlock()

 // 等待新事件
 select {
 case data := <-pt.waiter:
  return data, nil
 case <-ctx.Done():
  return nil, ctx.Err()
 }
}

// Queue 队列
type Queue struct {
 items []interface{}
 mu    sync.Mutex
 notEmpty chan struct{}
}

// NewQueue 创建队列
func NewQueue() *Queue {
 return &Queue{
  items:    make([]interface{}, 0),
  notEmpty: make(chan struct{}, 1),
 }
}

// Enqueue 入队
func (q *Queue) Enqueue(item interface{}) {
 q.mu.Lock()
 defer q.mu.Unlock()

 q.items = append(q.items, item)

 // 通知等待者
 select {
 case q.notEmpty <- struct{}{}:
 default:
 }
}

// Dequeue 出队
func (q *Queue) Dequeue(ctx context.Context) (interface{}, error) {
 for {
  q.mu.Lock()
  if len(q.items) > 0 {
   item := q.items[0]
   q.items = q.items[1:]
   q.mu.Unlock()
   return item, nil
  }
  q.mu.Unlock()

  // 等待队列非空
  select {
  case <-q.notEmpty:
   continue
  case <-ctx.Done():
   return nil, ctx.Err()
  }
 }
}

// Size 获取队列大小
func (q *Queue) Size() int {
 q.mu.Lock()
 defer q.mu.Unlock()
 return len(q.items)
}

// 使用示例
func main() {
 trigger := NewPersistentTrigger()

 // 先触发事件
 fmt.Println("=== 先触发，后等待 ===")
 trigger.Trigger("event1")
 trigger.Trigger("event2")

 // 后等待，应该能获取到持久化的事件
 ctx := context.Background()

 data1, err := trigger.Wait(ctx)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("获取到事件: %v\n", data1)

 data2, err := trigger.Wait(ctx)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("获取到事件: %v\n", data2)

 // 队列示例
 fmt.Println("\n=== 队列示例 ===")
 queue := NewQueue()

 // 先入队
 queue.Enqueue("item1")
 queue.Enqueue("item2")
 queue.Enqueue("item3")

 fmt.Printf("队列大小: %d\n", queue.Size())

 // 后出队
 item, err := queue.Dequeue(ctx)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("出队: %v\n", item)

 fmt.Printf("队列大小: %d\n", queue.Size())
}
```

#### 完整示例

```go
package main

import (
 "context"
 "fmt"
 "log"
 "sync"
 "time"
)

// MessageQueue 消息队列
type MessageQueue struct {
 name     string
 messages []Message
 mu       sync.Mutex
 waiters  []chan Message
}

// Message 消息
type Message struct {
 ID        string
 Content   string
 Timestamp time.Time
}

// NewMessageQueue 创建消息队列
func NewMessageQueue(name string) *MessageQueue {
 return &MessageQueue{
  name:     name,
  messages: make([]Message, 0),
  waiters:  make([]chan Message, 0),
 }
}

// Publish 发布消息
func (mq *MessageQueue) Publish(content string) string {
 msg := Message{
  ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
  Content:   content,
  Timestamp: time.Now(),
 }

 mq.mu.Lock()
 defer mq.mu.Unlock()

 // 如果有等待者，直接发送
 if len(mq.waiters) > 0 {
  waiter := mq.waiters[0]
  mq.waiters = mq.waiters[1:]

  select {
  case waiter <- msg:
   return msg.ID
  default:
  }
 }

 // 否则持久化存储
 mq.messages = append(mq.messages, msg)
 log.Printf("[消息队列 %s] 消息持久化: %s, 队列长度: %d", mq.name, msg.ID, len(mq.messages))

 return msg.ID
}

// Consume 消费消息
func (mq *MessageQueue) Consume(ctx context.Context) (Message, error) {
 mq.mu.Lock()

 // 检查是否有已持久化的消息
 if len(mq.messages) > 0 {
  msg := mq.messages[0]
  mq.messages = mq.messages[1:]
  mq.mu.Unlock()
  return msg, nil
 }

 // 创建等待channel
 waiter := make(chan Message, 1)
 mq.waiters = append(mq.waiters, waiter)
 mq.mu.Unlock()

 // 等待消息
 select {
 case msg := <-waiter:
  return msg, nil
 case <-ctx.Done():
  // 取消等待
  mq.mu.Lock()
  for i, w := range mq.waiters {
   if w == waiter {
    mq.waiters = append(mq.waiters[:i], mq.waiters[i+1:]...)
    break
   }
  }
  mq.mu.Unlock()
  return Message{}, ctx.Err()
 }
}

// EventStore 事件存储
type EventStore struct {
 events map[string][]Event
 mu     sync.RWMutex
}

// Event 事件
type Event struct {
 Type      string
 Data      interface{}
 Timestamp time.Time
}

// NewEventStore 创建事件存储
func NewEventStore() *EventStore {
 return &EventStore{
  events: make(map[string][]Event),
 }
}

// Append 追加事件
func (es *EventStore) Append(stream string, eventType string, data interface{}) {
 es.mu.Lock()
 defer es.mu.Unlock()

 event := Event{
  Type:      eventType,
  Data:      data,
  Timestamp: time.Now(),
 }

 es.events[stream] = append(es.events[stream], event)
}

// GetEvents 获取事件
func (es *EventStore) GetEvents(stream string) []Event {
 es.mu.RLock()
 defer es.mu.RUnlock()

 events := make([]Event, len(es.events[stream]))
 copy(events, es.events[stream])
 return events
}

// Replay 重放事件
func (es *EventStore) Replay(stream string, handler func(event Event)) {
 events := es.GetEvents(stream)

 for _, event := range events {
  handler(event)
 }
}

// ScheduledTask 定时任务
type ScheduledTask struct {
 tasks   []TaskSchedule
 mu      sync.Mutex
 running bool
}

// TaskSchedule 任务调度
type TaskSchedule struct {
 Name     string
 ExecuteAt time.Time
 Action   func()
}

// NewScheduledTask 创建定时任务
func NewScheduledTask() *ScheduledTask {
 return &ScheduledTask{
  tasks: make([]TaskSchedule, 0),
 }
}

// Schedule 调度任务
func (st *ScheduledTask) Schedule(name string, executeAt time.Time, action func()) {
 st.mu.Lock()
 defer st.mu.Unlock()

 st.tasks = append(st.tasks, TaskSchedule{
  Name:      name,
  ExecuteAt: executeAt,
  Action:    action,
 })

 log.Printf("[定时任务] 调度任务: %s, 执行时间: %v", name, executeAt)
}

// Start 启动调度器
func (st *ScheduledTask) Start() {
 st.mu.Lock()
 if st.running {
  st.mu.Unlock()
  return
 }
 st.running = true
 st.mu.Unlock()

 go func() {
  for {
   st.mu.Lock()
   tasks := make([]TaskSchedule, len(st.tasks))
   copy(tasks, st.tasks)
   st.mu.Unlock()

   now := time.Now()
   for _, task := range tasks {
    if now.After(task.ExecuteAt) {
     log.Printf("[定时任务] 执行任务: %s", task.Name)
     go task.Action()

     // 从列表中移除
     st.mu.Lock()
     for i, t := range st.tasks {
      if t.Name == task.Name {
       st.tasks = append(st.tasks[:i], st.tasks[i+1:]...)
       break
      }
     }
     st.mu.Unlock()
    }
   }

   time.Sleep(100 * time.Millisecond)
  }
 }()
}

func main() {
 ctx := context.Background()

 // 消息队列示例
 fmt.Println("=== 消息队列示例 ===")
 queue := NewMessageQueue("order_queue")

 // 先发布消息
 queue.Publish("订单创建: order_001")
 queue.Publish("订单创建: order_002")
 queue.Publish("订单创建: order_003")

 // 后消费消息
 for i := 0; i < 3; i++ {
  msg, err := queue.Consume(ctx)
  if err != nil {
   log.Fatal(err)
  }
  fmt.Printf("消费消息: %s - %s\n", msg.ID, msg.Content)
 }

 // 事件存储示例
 fmt.Println("\n=== 事件存储示例 ===")
 eventStore := NewEventStore()

 // 追加事件
 eventStore.Append("user_001", "UserCreated", map[string]string{"name": "Alice"})
 eventStore.Append("user_001", "UserUpdated", map[string]string{"email": "alice@example.com"})
 eventStore.Append("user_001", "UserLoggedIn", map[string]string{"ip": "192.168.1.1"})

 // 重放事件
 fmt.Println("重放用户事件:")
 eventStore.Replay("user_001", func(event Event) {
  fmt.Printf("  - %s: %v\n", event.Type, event.Data)
 })

 // 定时任务示例
 fmt.Println("\n=== 定时任务示例 ===")
 scheduler := NewScheduledTask()
 scheduler.Start()

 scheduler.Schedule("task1", time.Now().Add(200*time.Millisecond), func() {
  fmt.Println("→ 定时任务1执行")
 })

 scheduler.Schedule("task2", time.Now().Add(400*time.Millisecond), func() {
  fmt.Println("→ 定时任务2执行")
 })

 time.Sleep(600 * time.Millisecond)
 fmt.Println("✓ 定时任务完成")
}
```

#### 反例说明

**错误实现1：事件丢失**

```go
// ❌ 错误：事件没有持久化
func WrongNoPersistence() {
 ch := make(chan interface{})  // 无缓冲

 // 发送事件
 go func() {
  ch <- "event"  // 如果没有接收者，会阻塞
 }()
}
```

**问题**：没有持久化，事件可能丢失或阻塞。

**错误实现2：没有顺序保证**

```go
// ❌ 错误：使用map存储，没有顺序保证
func WrongNoOrder() {
 events := make(map[string]interface{})  // 错误：map无序
}
```

**问题**：map无序，不能保证事件的处理顺序。

#### 适用场景

| 场景 | 说明 |
|------|------|
| 消息队列 | 可靠的消息传递 |
| 事件溯源 | 事件存储和重放 |
| 任务调度 | 定时任务执行 |
| 工作流触发 | 持久化的工作流触发 |
| 通知系统 | 可靠的通知传递 |

#### 与其他模式关系

- **与临时触发器**：持久触发器保留事件，临时触发器丢弃事件
- **与延迟选择**：持久触发器可以实现延迟选择
- **与里程碑**：持久触发器可以触发里程碑

---

## 总结

本文档全面梳理了Go语言实现的23种工作流设计模式，涵盖了：

### 模式分类统计

| 类别 | 模式数量 | 主要特点 |
|------|----------|----------|
| 基础控制流模式 | 5 | 顺序、并行、选择、合并 |
| 高级分支与同步模式 | 5 | 多选、同步合并、鉴别器、N选M |
| 结构化模式 | 2 | 循环、终止 |
| 多实例模式 | 3 | 无同步、需同步、运行时确定 |
| 状态模式 | 3 | 延迟选择、交错路由、里程碑 |
| 取消与补偿模式 | 2 | 取消任务、取消案例 |
| 其他可判断模式 | 3 | 递归、临时触发器、持久触发器 |

### 实现要点

1. **上下文管理**：所有模式都应支持context.Context，实现取消和超时控制
2. **错误处理**：完善的错误处理和传播机制
3. **并发安全**：使用sync包确保并发安全
4. **资源管理**：使用defer确保资源释放
5. **可测试性**：设计可测试的接口和实现

### 选择建议

- **简单流程**：使用顺序、排他选择、简单合并
- **并行处理**：使用并行分支、同步、多实例模式
- **动态流程**：使用延迟选择、多选、同步合并
- **可靠性**：使用持久触发器、Saga模式、补偿机制
- **性能优化**：使用鉴别器、N选M、交错并行路由

---

*文档生成时间：2024年*

*作者：AI Assistant*
