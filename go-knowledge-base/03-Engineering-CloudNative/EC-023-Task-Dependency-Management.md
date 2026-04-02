# 任务依赖管理 (Task Dependency Management)

> **分类**: 工程与云原生
> **标签**: #task-dependency #dag #workflow

---

## DAG 任务依赖

```go
// 任务节点
type TaskNode struct {
    ID           string
    Name         string
    Execute      func(ctx context.Context) error
    Dependencies []string  // 依赖的任务ID
    Dependents   []string  // 被依赖的任务ID

    // 执行状态
    Status       TaskStatus
    Result       interface{}
    Error        error
    StartTime    *time.Time
    EndTime      *time.Time
}

// DAG 执行器
type DAGExecutor struct {
    nodes      map[string]*TaskNode
    mu         sync.RWMutex
    parallelism int
}

func (de *DAGExecutor) Execute(ctx context.Context) error {
    // 1. 拓扑排序检测循环依赖
    sorted, err := de.topologicalSort()
    if err != nil {
        return err
    }

    // 2. 构建依赖计数
    inDegree := make(map[string]int)
    for _, node := range de.nodes {
        inDegree[node.ID] = len(node.Dependencies)
    }

    // 3. 找到入度为0的节点（无依赖）
    var ready []*TaskNode
    for _, node := range sorted {
        if inDegree[node.ID] == 0 {
            ready = append(ready, node)
        }
    }

    // 4. 并行执行
    var wg sync.WaitGroup
    errChan := make(chan error, len(de.nodes))
    completed := make(chan string, len(de.nodes))

    // 启动调度器
    go func() {
        for completedID := range completed {
            de.mu.Lock()
            // 减少依赖该任务的节点入度
            for _, node := range de.nodes {
                for _, dep := range node.Dependencies {
                    if dep == completedID {
                        inDegree[node.ID]--
                        if inDegree[node.ID] == 0 {
                            ready = append(ready, node)
                        }
                    }
                }
            }
            de.mu.Unlock()
        }
    }()

    // 执行任务
    for len(ready) > 0 {
        node := ready[0]
        ready = ready[1:]

        wg.Add(1)
        go func(n *TaskNode) {
            defer wg.Done()

            // 执行前检查依赖是否都成功
            if !de.checkDependenciesSuccess(n) {
                n.Status = TaskStatusSkipped
                completed <- n.ID
                return
            }

            // 执行任务
            start := time.Now()
            n.StartTime = &start
            n.Status = TaskStatusRunning

            err := n.Execute(ctx)

            end := time.Now()
            n.EndTime = &end

            if err != nil {
                n.Status = TaskStatusFailed
                n.Error = err
                errChan <- fmt.Errorf("task %s failed: %w", n.ID, err)
            } else {
                n.Status = TaskStatusSuccess
            }

            completed <- n.ID
        }(node)
    }

    wg.Wait()
    close(completed)
    close(errChan)

    // 收集错误
    var errs []error
    for err := range errChan {
        errs = append(errs, err)
    }

    if len(errs) > 0 {
        return fmt.Errorf("DAG execution failed: %v", errs)
    }

    return nil
}

func (de *DAGExecutor) checkDependenciesSuccess(node *TaskNode) bool {
    for _, depID := range node.Dependencies {
        if dep, ok := de.nodes[depID]; ok {
            if dep.Status != TaskStatusSuccess {
                return false
            }
        }
    }
    return true
}
```

---

## 条件依赖

```go
// 条件执行
type ConditionalTask struct {
    TaskNode
    Condition func(ctx context.Context) bool
    IfTrue    string  // 条件为真时执行
    IfFalse   string  // 条件为假时执行
}

func (ct *ConditionalTask) Execute(ctx context.Context) error {
    if ct.Condition(ctx) {
        if ct.IfTrue != "" {
            return de.executeTaskByID(ctx, ct.IfTrue)
        }
    } else {
        if ct.IfFalse != "" {
            return de.executeTaskByID(ctx, ct.IfFalse)
        }
    }
    return nil
}

// 使用示例
validationTask := &ConditionalTask{
    TaskNode: TaskNode{
        ID:   "validate",
        Name: "Validate Data",
        Execute: func(ctx context.Context) error {
            // 验证数据
            return nil
        },
    },
    Condition: func(ctx context.Context) bool {
        data := ctx.Value("data").(*Data)
        return data.Amount > 1000
    },
    IfTrue:  "high-value-process",
    IfFalse: "standard-process",
}
```

---

## 动态依赖

```go
// 运行时动态确定依赖
type DynamicDependencyTask struct {
    TaskNode
    DependencyResolver func(ctx context.Context) []string
}

func (ddt *DynamicDependencyTask) ResolveDependencies(ctx context.Context) []string {
    return ddt.DependencyResolver(ctx)
}

// 使用示例
dataSplitTask := &DynamicDependencyTask{
    TaskNode: TaskNode{
        ID:   "data-split",
        Name: "Split Data",
    },
    DependencyResolver: func(ctx context.Context) []string {
        // 根据数据量动态决定分片数
        dataSize := getDataSize(ctx)
        shards := calculateShards(dataSize)

        var deps []string
        for i := 0; i < shards; i++ {
            shardID := fmt.Sprintf("process-shard-%d", i)
            deps = append(deps, shardID)

            // 动态注册子任务
            registerShardTask(shardID, i)
        }

        return deps
    },
}
```

---

## 依赖可视化

```go
// 生成 Mermaid 图
type DAGVisualizer struct {
    nodes map[string]*TaskNode
}

func (dv *DAGVisualizer) ToMermaid() string {
    var sb strings.Builder
    sb.WriteString("graph TD\n")

    for _, node := range dv.nodes {
        // 节点定义
        status := node.Status
        color := "gray"
        switch status {
        case TaskStatusSuccess:
            color = "green"
        case TaskStatusFailed:
            color = "red"
        case TaskStatusRunning:
            color = "blue"
        }

        sb.WriteString(fmt.Sprintf("    %s[%s]:::status-%s\n",
            node.ID, node.Name, color))

        // 依赖关系
        for _, dep := range node.Dependencies {
            sb.WriteString(fmt.Sprintf("    %s --> %s\n", dep, node.ID))
        }
    }

    // 样式定义
    sb.WriteString("    classDef status-green fill:#90EE90\n")
    sb.WriteString("    classDef status-red fill:#FFB6C1\n")
    sb.WriteString("    classDef status-blue fill:#87CEEB\n")
    sb.WriteString("    classDef status-gray fill:#D3D3D3\n")

    return sb.String()
}
```
