# 任务上下文取消模式 (Task Context Cancellation Patterns)

> **分类**: 工程与云原生  
> **标签**: #context #cancellation #graceful-shutdown #patterns

---

## 协作式取消

```go
// 协作式取消模式
// 任务主动检查取消信号并清理资源

type CancellableTask struct {
    id       string
    cancel   context.CancelFunc
    done     chan struct{}
    cleanup  []func()
}

func (ct *CancellableTask) Run(ctx context.Context) error {
    // 添加清理函数
    defer ct.runCleanup()
    
    // 主要处理循环
    for {
        select {
        case <-ctx.Done():
            // 收到取消信号
            return ct.handleCancellation(ctx)
            
        case work := <-ct.workQueue:
            // 检查取消状态
            if err := ct.checkContext(ctx); err != nil {
                // 将未处理的工作重新入队
                ct.requeue(work)
                return err
            }
            
            if err := ct.processWork(ctx, work); err != nil {
                return err
            }
        }
    }
}

func (ct *CancellableTask) handleCancellation(ctx context.Context) error {
    // 记录取消原因
    cause := context.Cause(ctx)
    
    switch {
    case errors.Is(cause, context.DeadlineExceeded):
        log.Printf("Task %s cancelled due to timeout", ct.id)
        return &TaskCancelledError{Reason: "timeout", Cause: cause}
        
    case errors.Is(cause, context.Canceled):
        log.Printf("Task %s cancelled by user", ct.id)
        return &TaskCancelledError{Reason: "user_request", Cause: cause}
        
    default:
        log.Printf("Task %s cancelled: %v", ct.id, cause)
        return &TaskCancelledError{Reason: "unknown", Cause: cause}
    }
}

func (ct *CancellableTask) checkContext(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        return nil
    }
}

func (ct *CancellableTask) AddCleanup(fn func()) {
    ct.mu.Lock()
    defer ct.mu.Unlock()
    ct.cleanup = append(ct.cleanup, fn)
}

func (ct *CancellableTask) runCleanup() {
    ct.mu.Lock()
    cleanup := ct.cleanup
    ct.cleanup = nil
    ct.mu.Unlock()
    
    // 逆序执行清理
    for i := len(cleanup) - 1; i >= 0; i-- {
        cleanup[i]()
    }
}
```

---

## 级联取消

```go
// 父子任务级联取消
type HierarchicalCanceller struct {
    mu        sync.RWMutex
    children  map[string]*HierarchicalCanceller
    parent    *HierarchicalCanceller
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewHierarchicalCanceller(parent *HierarchicalCanceller) *HierarchicalCanceller {
    hc := &HierarchicalCanceller{
        children: make(map[string]*HierarchicalCanceller),
    }
    
    if parent != nil {
        hc.parent = parent
        hc.ctx, hc.cancel = context.WithCancel(parent.ctx)
        parent.addChild(hc)
    } else {
        hc.ctx, hc.cancel = context.WithCancel(context.Background())
    }
    
    return hc
}

func (hc *HierarchicalCanceller) addChild(child *HierarchicalCanceller) {
    hc.mu.Lock()
    defer hc.mu.Unlock()
    hc.children[child.id] = child
}

func (hc *HierarchicalCanceller) Cancel(cause error) {
    // 取消自己
    hc.cancel()
    
    // 级联取消所有子任务
    hc.mu.RLock()
    children := make([]*HierarchicalCanceller, 0, len(hc.children))
    for _, child := range hc.children {
        children = append(children, child)
    }
    hc.mu.RUnlock()
    
    var wg sync.WaitGroup
    for _, child := range children {
        wg.Add(1)
        go func(c *HierarchicalCanceller) {
            defer wg.Done()
            c.Cancel(fmt.Errorf("parent cancelled: %w", cause))
        }(child)
    }
    
    wg.Wait()
}

func (hc *HierarchicalCanceller) RemoveChild(id string) {
    hc.mu.Lock()
    defer hc.mu.Unlock()
    delete(hc.children, id)
}
```

---

## 取消传播策略

```go
// 取消传播策略
type CancellationPolicy int

const (
    Immediate CancellationPolicy = iota  // 立即取消
    Graceful                              // 优雅取消，等待当前工作完成
    Drain                                 // 排空模式，完成队列中所有工作
)

type CancellationManager struct {
    policy CancellationPolicy
    gracePeriod time.Duration
}

func (cm *CancellationManager) CancelWithPolicy(ctx context.Context, executor *TaskExecutor, policy CancellationPolicy) error {
    switch policy {
    case Immediate:
        return cm.cancelImmediate(ctx, executor)
        
    case Graceful:
        return cm.cancelGraceful(ctx, executor)
        
    case Drain:
        return cm.cancelDrain(ctx, executor)
        
    default:
        return cm.cancelImmediate(ctx, executor)
    }
}

func (cm *CancellationManager) cancelGraceful(ctx context.Context, executor *TaskExecutor) error {
    // 停止接受新任务
    executor.StopAccepting()
    
    // 等待进行中的任务完成
    ctx, cancel := context.WithTimeout(ctx, cm.gracePeriod)
    defer cancel()
    
    done := make(chan struct{})
    go func() {
        executor.WaitForRunningTasks()
        close(done)
    }()
    
    select {
    case <-done:
        log.Println("All tasks completed gracefully")
        return nil
        
    case <-ctx.Done():
        log.Println("Grace period expired, forcing cancellation")
        return cm.cancelImmediate(ctx, executor)
    }
}

func (cm *CancellationManager) cancelDrain(ctx context.Context, executor *TaskExecutor) error {
    // 停止接受新任务
    executor.StopAccepting()
    
    // 处理队列中所有任务
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
            
        default:
            task, err := executor.Dequeue(ctx)
            if err == ErrQueueEmpty {
                // 队列已空，取消进行中的任务
                return cm.cancelImmediate(ctx, executor)
            }
            
            // 处理任务
            if err := executor.ExecuteTask(ctx, task); err != nil {
                log.Printf("Task execution failed during drain: %v", err)
            }
        }
    }
}
```

---

## 取消超时控制

```go
// 取消操作的超时控制
type CancellationWithTimeout struct {
    timeout time.Duration
}

func (cwt *CancellationWithTimeout) CancelTask(ctx context.Context, taskID string) error {
    ctx, cancel := context.WithTimeout(ctx, cwt.timeout)
    defer cancel()
    
    result := make(chan error, 1)
    
    go func() {
        result <- cwt.doCancel(taskID)
    }()
    
    select {
    case err := <-result:
        return err
        
    case <-ctx.Done():
        return fmt.Errorf("cancel operation timed out: %w", ctx.Err())
    }
}

// 批量取消
func (cwt *CancellationWithTimeout) CancelBatch(ctx context.Context, taskIDs []string) CancelBatchResult {
    ctx, cancel := context.WithTimeout(ctx, cwt.timeout)
    defer cancel()
    
    result := CancelBatchResult{
        Succeeded: make([]string, 0),
        Failed:    make(map[string]error),
    }
    
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for _, id := range taskIDs {
        wg.Add(1)
        go func(taskID string) {
            defer wg.Done()
            
            err := cwt.doCancel(taskID)
            
            mu.Lock()
            defer mu.Unlock()
            
            if err != nil {
                result.Failed[taskID] = err
            } else {
                result.Succeeded = append(result.Succeeded, taskID)
            }
        }(id)
    }
    
    // 等待或超时
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        return result
    case <-ctx.Done():
        result.Incomplete = true
        return result
    }
}
```
