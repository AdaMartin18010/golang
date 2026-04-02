# 任务调试与诊断 (Task Debugging & Diagnostics)

> **分类**: 工程与云原生  
> **标签**: #debugging #diagnostics #troubleshooting

---

## 任务调试接口

```go
type TaskDebugger struct {
    store     TaskStore
    executor  *TaskExecutor
}

// 获取任务详细信息
func (td *TaskDebugger) GetTaskDetails(ctx context.Context, taskID string) (*TaskDetails, error) {
    task, err := td.store.Get(ctx, taskID)
    if err != nil {
        return nil, err
    }
    
    details := &TaskDetails{
        Task:        task,
        StackTrace:  td.getStackTrace(taskID),
        Variables:   td.getVariables(taskID),
        Logs:        td.getRecentLogs(taskID, 100),
        Events:      td.getEventHistory(taskID),
        Performance: td.getPerformanceMetrics(taskID),
    }
    
    return details, nil
}

// 单步执行
func (td *TaskDebugger) StepExecute(ctx context.Context, taskID string) error {
    task, _ := td.store.Get(ctx, taskID)
    
    // 设置断点模式
    task.DebugMode = true
    task.Breakpoints = []string{"next"}
    
    // 执行一步
    return td.executor.Step(ctx, task)
}

// 设置断点
func (td *TaskDebugger) SetBreakpoint(ctx context.Context, taskID string, step string) error {
    return td.store.AddBreakpoint(ctx, taskID, step)
}

// 修改变量
func (td *TaskDebugger) ModifyVariable(ctx context.Context, taskID string, name string, value interface{}) error {
    return td.executor.SetVariable(ctx, taskID, name, value)
}
```

---

## 诊断工具

```go
type TaskDiagnostics struct {
    analyzer *TaskAnalyzer
}

// 分析任务失败原因
func (td *TaskDiagnostics) DiagnoseFailure(ctx context.Context, taskID string) (*Diagnosis, error) {
    task, _ := td.analyzer.store.Get(ctx, taskID)
    
    diagnosis := &Diagnosis{
        TaskID: taskID,
        Issues: []Issue{},
    }
    
    // 检查超时
    if task.Status == TaskStatusTimeout {
        diagnosis.Issues = append(diagnosis.Issues, Issue{
            Type:        "timeout",
            Severity:    "critical",
            Description: fmt.Sprintf("Task exceeded timeout of %v", task.Timeout),
            Suggestion:  "Consider increasing timeout or optimizing task",
        })
    }
    
    // 检查内存
    if task.ResourceUsage.Memory > task.ResourceRequest.Memory*1.5 {
        diagnosis.Issues = append(diagnosis.Issues, Issue{
            Type:        "memory_leak",
            Severity:    "warning",
            Description: "Task used significantly more memory than requested",
            Suggestion:  "Review memory allocation in task",
        })
    }
    
    // 检查重试
    if task.RetryCount > 5 {
        diagnosis.Issues = append(diagnosis.Issues, Issue{
            Type:        "excessive_retries",
            Severity:    "warning",
            Description: fmt.Sprintf("Task retried %d times", task.RetryCount),
            Suggestion:  "Check for flaky dependencies",
        })
    }
    
    return diagnosis, nil
}

// 性能分析
func (td *TaskDiagnostics) ProfileTask(ctx context.Context, taskID string) (*Profile, error) {
    task, _ := td.analyzer.store.Get(ctx, taskID)
    
    profile := &Profile{
        TaskID: taskID,
    }
    
    // CPU 分析
    profile.CPUProfile = CPUProfile{
        TotalTime: task.Duration,
        Hotspots:  td.analyzer.findCPUHotspots(task),
    }
    
    // 内存分析
    profile.MemoryProfile = MemoryProfile{
        PeakUsage: task.ResourceUsage.Memory,
        Allocations: td.analyzer.getAllocationStats(task),
    }
    
    // 阻塞分析
    profile.BlockProfile = BlockProfile{
        BlockedTime: td.analyzer.getBlockedTime(task),
        BlockPoints: td.analyzer.getBlockPoints(task),
    }
    
    return profile, nil
}
```

---

## 任务回放

```go
type TaskReplayer struct {
    eventStore EventStore
    executor   *TaskExecutor
}

func (tr *TaskReplayer) Replay(ctx context.Context, taskID string, fromEvent int) error {
    // 获取任务事件历史
    events, err := tr.eventStore.GetEvents(ctx, taskID, fromEvent)
    if err != nil {
        return err
    }
    
    // 重建任务状态
    task := tr.reconstructTask(events)
    
    // 在隔离环境重放
    replayCtx := context.WithValue(ctx, "replay_mode", true)
    
    return tr.executor.Execute(replayCtx, task)
}

func (tr *TaskReplayer) CompareRuns(ctx context.Context, taskID string, run1, run2 int) (*Comparison, error) {
    // 获取两次运行的结果
    result1, _ := tr.store.GetRunResult(ctx, taskID, run1)
    result2, _ := tr.store.GetRunResult(ctx, taskID, run2)
    
    return &Comparison{
        Differences: tr.compareResults(result1, result2),
        Similarity:  tr.calculateSimilarity(result1, result2),
    }, nil
}
```

---

## 实时监控诊断

```go
type LiveDiagnostics struct {
    subscribers map[string][]chan DiagnosticEvent
    mu          sync.RWMutex
}

func (ld *LiveDiagnostics) Subscribe(taskID string) chan DiagnosticEvent {
    ch := make(chan DiagnosticEvent, 100)
    
    ld.mu.Lock()
    ld.subscribers[taskID] = append(ld.subscribers[taskID], ch)
    ld.mu.Unlock()
    
    return ch
}

func (ld *LiveDiagnostics) Publish(taskID string, event DiagnosticEvent) {
    ld.mu.RLock()
    subs := ld.subscribers[taskID]
    ld.mu.RUnlock()
    
    for _, ch := range subs {
        select {
        case ch <- event:
        default:
            // 通道满，丢弃
        }
    }
}

// WebSocket 推送
func (ld *LiveDiagnostics) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    taskID := r.URL.Query().Get("task_id")
    
    conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
    defer conn.Close()
    
    events := ld.Subscribe(taskID)
    defer ld.Unsubscribe(taskID, events)
    
    for event := range events {
        conn.WriteJSON(event)
    }
}
```
