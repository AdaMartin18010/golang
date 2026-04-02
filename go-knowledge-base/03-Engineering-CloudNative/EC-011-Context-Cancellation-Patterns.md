# Context 取消模式 (Context Cancellation Patterns)

> **分类**: 工程与云原生
> **标签**: #context #cancellation #graceful-shutdown

---

## 取消传播链

```go
func ProcessWithCancellation(parentCtx context.Context) error {
    // 创建可取消的上下文
    ctx, cancel := context.WithCancel(parentCtx)
    defer cancel()

    // 启动多个子任务
    errChan := make(chan error, 3)

    go func() {
        errChan <- processStep1(ctx)
    }()

    go func() {
        errChan <- processStep2(ctx)
    }()

    go func() {
        errChan <- processStep3(ctx)
    }()

    // 等待任一任务完成或出错
    for i := 0; i < 3; i++ {
        if err := <-errChan; err != nil {
            cancel()  // 取消其他任务
            return err
        }
    }

    return nil
}
```

---

## 优雅取消 HTTP 请求

```go
func HTTPRequestWithCancellation(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        // 检查是否是取消错误
        if ctx.Err() == context.Canceled {
            return nil, fmt.Errorf("request cancelled: %w", err)
        }
        if ctx.Err() == context.DeadlineExceeded {
            return nil, fmt.Errorf("request timeout: %w", err)
        }
        return nil, err
    }

    return resp, nil
}
```

---

## 数据库查询取消

```go
func QueryWithCancellation(ctx context.Context, db *sql.DB, query string) (*sql.Rows, error) {
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        select {
        case <-ctx.Done():
            // 上下文被取消
            return nil, fmt.Errorf("query cancelled: %w", ctx.Err())
        default:
            return nil, err
        }
    }

    return rows, nil
}

// 扫描时检查取消
func ScanWithCancellation(ctx context.Context, rows *sql.Rows, dest interface{}) error {
    if err := ctx.Err(); err != nil {
        return err
    }

    return rows.Scan(dest)
}
```

---

## 可取消的工作池

```go
type CancellableWorkerPool struct {
    ctx     context.Context
    cancel  context.CancelFunc
    jobs    chan Job
    workers int
    wg      sync.WaitGroup
}

func NewCancellableWorkerPool(workers int) *CancellableWorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &CancellableWorkerPool{
        ctx:     ctx,
        cancel:  cancel,
        jobs:    make(chan Job),
        workers: workers,
    }
}

func (p *CancellableWorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *CancellableWorkerPool) worker(id int) {
    defer p.wg.Done()

    for {
        select {
        case job, ok := <-p.jobs:
            if !ok {
                return
            }

            // 每个任务使用派生上下文
            jobCtx, cancel := context.WithTimeout(p.ctx, 30*time.Second)
            p.executeJob(jobCtx, job)
            cancel()

        case <-p.ctx.Done():
            // 取消信号，处理剩余任务
            for {
                select {
                case job := <-p.jobs:
                    job.OnCancelled()
                default:
                    return
                }
            }
        }
    }
}

func (p *CancellableWorkerPool) Cancel() {
    p.cancel()
}

func (p *CancellableWorkerPool) Stop() {
    close(p.jobs)
    p.wg.Wait()
}
```

---

## 级联超时控制

```go
func CascadeTimeout(parentCtx context.Context, stages []Stage) error {
    remaining := 10 * time.Second  // 总超时

    for i, stage := range stages {
        start := time.Now()

        // 每个阶段分配剩余时间的一部分
        stageTimeout := remaining / time.Duration(len(stages)-i)

        ctx, cancel := context.WithTimeout(parentCtx, stageTimeout)

        err := stage.Execute(ctx)
        cancel()

        if err != nil {
            return fmt.Errorf("stage %d failed: %w", i, err)
        }

        // 扣除已用时间
        elapsed := time.Since(start)
        remaining -= elapsed

        if remaining <= 0 {
            return fmt.Errorf("timeout exceeded")
        }
    }

    return nil
}
```

---

## 取消原因传递

```go
type CancellableTask struct {
    ctx    context.Context
    cancel context.CancelCauseFunc
}

func (t *CancellableTask) Run() error {
    ctx, cancel := context.WithCancelCause(t.ctx)
    t.cancel = cancel

    go func() {
        if err := doWork(ctx); err != nil {
            cancel(err)  // 传递取消原因
        }
    }()

    <-ctx.Done()

    if err := context.Cause(ctx); err != nil {
        return fmt.Errorf("task failed: %w", err)
    }

    return nil
}

func (t *CancellableTask) Stop(reason error) {
    if t.cancel != nil {
        t.cancel(reason)
    }
}
```
