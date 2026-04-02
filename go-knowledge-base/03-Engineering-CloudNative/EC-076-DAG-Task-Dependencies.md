# DAG 任务依赖调度 (DAG Task Dependencies)

> **分类**: 工程与云原生
> **标签**: #dag #workflow #dependencies #graph
> **参考**: Airflow DAG, Temporal Workflow, Argo Workflows

---

## DAG 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DAG Task Scheduler Architecture                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  DAG Definition:                                                             │
│                                                                              │
│       ┌─────┐                                                               │
│       │  A  │ ◄── Start Task                                                │
│       └──┬──┘                                                               │
│          │                                                                  │
│     ┌────┴────┐                                                              │
│     ▼         ▼                                                              │
│  ┌─────┐   ┌─────┐                                                           │
│  │  B  │   │  C  │ ◄── Parallel Execution                                    │
│  └──┬──┘   └──┬──┘                                                           │
│     │         │                                                              │
│     └────┬────┘                                                              │
│          ▼                                                                  │
│       ┌─────┐                                                               │
│       │  D  │ ◄── Join (Wait for B & C)                                     │
│       └──┬──┘                                                               │
│          │                                                                  │
│          ▼                                                                  │
│       ┌─────┐                                                               │
│       │  E  │ ◄── End Task                                                  │
│       └─────┘                                                               │
│                                                                              │
│  Execution States:                                                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│  │ PENDING  │─►│ RUNNABLE │─►│ RUNNING  │─►│COMPLETED │  │  FAILED  │      │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘      │
│                      │                          │              │            │
│                      └──────────────────────────┴──────────────┘            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心实现

```go
package dag

import (
    "context"
    "fmt"
    "sync"
)

// TaskState 任务状态
type TaskState string

const (
    TaskStatePending    TaskState = "pending"
    TaskStateRunnable   TaskState = "runnable"
    TaskStateRunning    TaskState = "running"
    TaskStateCompleted  TaskState = "completed"
    TaskStateFailed     TaskState = "failed"
    TaskStateSkipped    TaskState = "skipped"
)

// Task DAG 任务
type Task struct {
    ID          string
    Name        string
    Handler     TaskHandler

    // 依赖关系
    Dependencies []string // 依赖的任务 ID
    Dependents   []string // 依赖于本任务的任务 ID

    // 状态
    State       TaskState
    Result      interface{}
    Error       error

    // 配置
    RetryCount  int
    MaxRetries  int
    SkipOnError bool // 失败后是否跳过而非失败
}

// TaskHandler 任务处理器
type TaskHandler func(ctx context.Context, inputs map[string]interface{}) (interface{}, error)

// DAG 有向无环图
type DAG struct {
    ID    string
    Name  string
    Tasks map[string]*Task

    // 缓存
    sortedTasks []*Task // 拓扑排序后的任务
}

// NewDAG 创建 DAG
func NewDAG(id, name string) *DAG {
    return &DAG{
        ID:    id,
        Name:  name,
        Tasks: make(map[string]*Task),
    }
}

// AddTask 添加任务
func (d *DAG) AddTask(task *Task) error {
    if _, exists := d.Tasks[task.ID]; exists {
        return fmt.Errorf("task %s already exists", task.ID)
    }

    d.Tasks[task.ID] = task
    d.sortedTasks = nil // 清除缓存

    return nil
}

// AddDependency 添加依赖
func (d *DAG) AddDependency(taskID, dependsOn string) error {
    task, ok := d.Tasks[taskID]
    if !ok {
        return fmt.Errorf("task %s not found", taskID)
    }

    dep, ok := d.Tasks[dependsOn]
    if !ok {
        return fmt.Errorf("dependency task %s not found", dependsOn)
    }

    // 检查循环依赖
    if err := d.checkCircularDependency(taskID, dependsOn); err != nil {
        return err
    }

    task.Dependencies = append(task.Dependencies, dependsOn)
    dep.Dependents = append(dep.Dependents, taskID)
    d.sortedTasks = nil

    return nil
}

// checkCircularDependency 检查循环依赖
func (d *DAG) checkCircularDependency(taskID, newDep string) error {
    // 使用 DFS 检查
    visited := make(map[string]bool)

    var dfs func(string) bool
    dfs = func(id string) bool {
        if id == taskID {
            return true // 发现循环
        }
        if visited[id] {
            return false
        }
        visited[id] = true

        task := d.Tasks[id]
        for _, dep := range task.Dependencies {
            if dfs(dep) {
                return true
            }
        }
        return false
    }

    if dfs(newDep) {
        return fmt.Errorf("circular dependency detected")
    }

    return nil
}

// TopologicalSort 拓扑排序
func (d *DAG) TopologicalSort() ([]*Task, error) {
    if d.sortedTasks != nil {
        return d.sortedTasks, nil
    }

    // Kahn 算法
    inDegree := make(map[string]int)
    for id := range d.Tasks {
        inDegree[id] = 0
    }

    for _, task := range d.Tasks {
        for _, dep := range task.Dependencies {
            inDegree[task.ID]++
        }
    }

    // 找到入度为 0 的节点
    var queue []string
    for id, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, id)
        }
    }

    var sorted []*Task

    for len(queue) > 0 {
        // 取出入度为 0 的节点
        id := queue[0]
        queue = queue[1:]

        sorted = append(sorted, d.Tasks[id])

        // 减少依赖此节点的其他节点的入度
        for _, dependentID := range d.Tasks[id].Dependents {
            inDegree[dependentID]--
            if inDegree[dependentID] == 0 {
                queue = append(queue, dependentID)
            }
        }
    }

    if len(sorted) != len(d.Tasks) {
        return nil, fmt.Errorf("cycle detected in DAG")
    }

    d.sortedTasks = sorted
    return sorted, nil
}

// GetRunnableTasks 获取可运行的任务
func (d *DAG) GetRunnableTasks() []*Task {
    var runnable []*Task

    for _, task := range d.Tasks {
        if task.State != TaskStatePending {
            continue
        }

        // 检查所有依赖是否已完成
        allCompleted := true
        for _, depID := range task.Dependencies {
            dep := d.Tasks[depID]
            if dep.State != TaskStateCompleted {
                allCompleted = false
                break
            }
        }

        if allCompleted {
            task.State = TaskStateRunnable
            runnable = append(runnable, task)
        }
    }

    return runnable
}

// IsCompleted 检查 DAG 是否完成
func (d *DAG) IsCompleted() bool {
    for _, task := range d.Tasks {
        if task.State != TaskStateCompleted && task.State != TaskStateSkipped && task.State != TaskStateFailed {
            return false
        }
    }
    return true
}

// HasFailed 检查是否有失败任务
func (d *DAG) HasFailed() bool {
    for _, task := range d.Tasks {
        if task.State == TaskStateFailed {
            return true
        }
    }
    return false
}
```

---

## DAG 执行器

```go
package dag

import (
    "context"
    "fmt"
    "sync"
)

// Executor DAG 执行器
type Executor struct {
    workerCount int

    // 状态
    mu       sync.RWMutex
    running  bool
}

// NewExecutor 创建执行器
func NewExecutor(workerCount int) *Executor {
    return &Executor{
        workerCount: workerCount,
    }
}

// Execute 执行 DAG
func (e *Executor) Execute(ctx context.Context, dag *DAG) error {
    e.mu.Lock()
    if e.running {
        e.mu.Unlock()
        return fmt.Errorf("executor already running")
    }
    e.running = true
    e.mu.Unlock()

    defer func() {
        e.mu.Lock()
        e.running = false
        e.mu.Unlock()
    }()

    // 拓扑排序
    _, err := dag.TopologicalSort()
    if err != nil {
        return err
    }

    // 创建工作池
    taskQueue := make(chan *Task, len(dag.Tasks))
    resultQueue := make(chan *Task, len(dag.Tasks))

    var wg sync.WaitGroup

    // 启动工作线程
    for i := 0; i < e.workerCount; i++ {
        wg.Add(1)
        go e.worker(ctx, &wg, taskQueue, resultQueue)
    }

    // 结果收集器
    go e.resultCollector(ctx, dag, resultQueue, taskQueue)

    // 启动初始任务（没有依赖的任务）
    runnable := dag.GetRunnableTasks()
    for _, task := range runnable {
        taskQueue <- task
    }

    // 等待完成
    for !dag.IsCompleted() && !dag.HasFailed() {
        select {
        case <-ctx.Done():
            close(taskQueue)
            wg.Wait()
            return ctx.Err()
        default:
        }
    }

    close(taskQueue)
    wg.Wait()

    if dag.HasFailed() {
        return fmt.Errorf("DAG execution failed")
    }

    return nil
}

func (e *Executor) worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan *Task, results chan<- *Task) {
    defer wg.Done()

    for task := range tasks {
        if ctx.Err() != nil {
            return
        }

        e.executeTask(ctx, task)
        results <- task
    }
}

func (e *Executor) executeTask(ctx context.Context, task *Task) {
    task.State = TaskStateRunning

    // 收集输入
    inputs := make(map[string]interface{})
    // 从依赖任务获取输出

    // 执行任务
    result, err := task.Handler(ctx, inputs)

    if err != nil {
        task.Error = err
        task.State = TaskStateFailed
    } else {
        task.Result = result
        task.State = TaskStateCompleted
    }
}

func (e *Executor) resultCollector(ctx context.Context, dag *DAG, results <-chan *Task, taskQueue chan<- *Task) {
    for {
        select {
        case <-ctx.Done():
            return
        case task, ok := <-results:
            if !ok {
                return
            }

            // 任务完成，触发依赖任务
            if task.State == TaskStateCompleted {
                for _, dependentID := range task.Dependents {
                    dependent := dag.Tasks[dependentID]
                    if e.canRun(dag, dependent) {
                        taskQueue <- dependent
                    }
                }
            }
        }
    }
}

func (e *Executor) canRun(dag *DAG, task *Task) bool {
    if task.State != TaskStatePending {
        return false
    }

    for _, depID := range task.Dependencies {
        dep := dag.Tasks[depID]
        if dep.State != TaskStateCompleted {
            return false
        }
    }

    return true
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"
    "time"

    "dag"
)

func main() {
    // 创建 DAG
    d := dag.NewDAG("data-pipeline", "Data Processing Pipeline")

    // 添加任务
    tasks := []*dag.Task{
        {
            ID:      "extract",
            Name:    "Extract Data",
            Handler: extractHandler,
        },
        {
            ID:      "transform-a",
            Name:    "Transform A",
            Handler: transformAHandler,
        },
        {
            ID:      "transform-b",
            Name:    "Transform B",
            Handler: transformBHandler,
        },
        {
            ID:      "validate",
            Name:    "Validate Data",
            Handler: validateHandler,
        },
        {
            ID:      "load",
            Name:    "Load to Database",
            Handler: loadHandler,
        },
    }

    for _, task := range tasks {
        d.AddTask(task)
    }

    // 设置依赖
    d.AddDependency("transform-a", "extract")
    d.AddDependency("transform-b", "extract")
    d.AddDependency("validate", "transform-a")
    d.AddDependency("validate", "transform-b")
    d.AddDependency("load", "validate")

    // 执行
    executor := dag.NewExecutor(4)
    if err := executor.Execute(context.Background(), d); err != nil {
        fmt.Printf("DAG execution failed: %v\n", err)
    }
}

func extractHandler(ctx context.Context, inputs map[string]interface{}) (interface{}, error) {
    fmt.Println("Extracting data...")
    time.Sleep(1 * time.Second)
    return "raw-data", nil
}

func transformAHandler(ctx context.Context, inputs map[string]interface{}) (interface{}, error) {
    fmt.Println("Transforming data A...")
    time.Sleep(1 * time.Second)
    return "transformed-a", nil
}

func transformBHandler(ctx context.Context, inputs map[string]interface{}) (interface{}, error) {
    fmt.Println("Transforming data B...")
    time.Sleep(1 * time.Second)
    return "transformed-b", nil
}

func validateHandler(ctx context.Context, inputs map[string]interface{}) (interface{}, error) {
    fmt.Println("Validating data...")
    time.Sleep(500 * time.Millisecond)
    return "validated-data", nil
}

func loadHandler(ctx context.Context, inputs map[string]interface{}) (interface{}, error) {
    fmt.Println("Loading to database...")
    time.Sleep(1 * time.Second)
    return "loaded", nil
}
```
