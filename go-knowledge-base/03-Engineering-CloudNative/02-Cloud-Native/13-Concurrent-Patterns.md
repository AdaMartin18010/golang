# 并发模式 (Concurrent Patterns)

> **分类**: 工程与云原生
> **标签**: #concurrency #patterns #goroutine

---

## Fan-Out / Fan-In

```go
// Fan-Out: 多个 goroutine 处理任务
func FanOut(ctx context.Context, tasks []Task, workers int) []Result {
    taskCh := make(chan Task)
    resultCh := make(chan Result, len(tasks))

    var wg sync.WaitGroup

    // 启动 workers
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for task := range taskCh {
                select {
                case <-ctx.Done():
                    return
                default:
                    result := process(task)
                    resultCh <- result
                }
            }
        }(i)
    }

    // 分发任务
    go func() {
        for _, task := range tasks {
            taskCh <- task
        }
        close(taskCh)
    }()

    // 等待完成
    go func() {
        wg.Wait()
        close(resultCh)
    }()

    // 收集结果
    var results []Result
    for r := range resultCh {
        results = append(results, r)
    }

    return results
}

// Fan-In: 合并多个 channel
func FanIn(ctx context.Context, channels ...<-chan Result) <-chan Result {
    out := make(chan Result)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan Result) {
            defer wg.Done()
            for r := range c {
                select {
                case <-ctx.Done():
                    return
                case out <- r:
                }
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

---

## Pipeline

```go
func Pipeline(ctx context.Context, stages ...Stage) Stage {
    return func(in <-chan Data) <-chan Data {
        current := in
        for _, stage := range stages {
            current = stage(ctx, current)
        }
        return current
    }
}

// 使用
generator := func(ctx context.Context) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for i := 0; i < 100; i++ {
            select {
            case <-ctx.Done():
                return
            case out <- Data{Value: i}:
            }
        }
    }()
    return out
}

stage1 := func(ctx context.Context, in <-chan Data) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for d := range in {
            select {
            case <-ctx.Done():
                return
            case out <- Data{Value: d.Value * 2}:
            }
        }
    }()
    return out
}

stage2 := func(ctx context.Context, in <-chan Data) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for d := range in {
            select {
            case <-ctx.Done():
                return
            case out <- Data{Value: d.Value + 1}:
            }
        }
    }()
    return out
}

pipeline := Pipeline(context.Background(), stage1, stage2)
result := pipeline(generator(context.Background()))
```

---

## Worker Pool with Cancellation

```go
type Pool struct {
    workers int
    tasks   chan func(ctx context.Context)
    ctx     context.Context
    cancel  context.CancelFunc
    wg      sync.WaitGroup
}

func NewPool(workers int) *Pool {
    ctx, cancel := context.WithCancel(context.Background())
    return &Pool{
        workers: workers,
        tasks:   make(chan func(ctx context.Context)),
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (p *Pool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *Pool) worker(id int) {
    defer p.wg.Done()
    for {
        select {
        case task, ok := <-p.tasks:
            if !ok {
                return
            }
            task(p.ctx)
        case <-p.ctx.Done():
            return
        }
    }
}

func (p *Pool) Submit(task func(ctx context.Context)) bool {
    select {
    case p.tasks <- task:
        return true
    case <-p.ctx.Done():
        return false
    }
}

func (p *Pool) Stop() {
    p.cancel()
    close(p.tasks)
    p.wg.Wait()
}
```

---

## Semaphore 模式

```go
type Semaphore struct {
    ch chan struct{}
}

func NewSemaphore(n int) *Semaphore {
    return &Semaphore{ch: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire(ctx context.Context) error {
    select {
    case s.ch <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (s *Semaphore) Release() {
    select {
    case <-s.ch:
    default:
        panic("semaphore release without acquire")
    }
}

// 使用
func ProcessWithLimit(ctx context.Context, items []Item, limit int) {
    sem := NewSemaphore(limit)
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()

            if err := sem.Acquire(ctx); err != nil {
                return
            }
            defer sem.Release()

            process(i)
        }(item)
    }

    wg.Wait()
}
```

---

## Or-Done 模式

```go
func OrDone(ctx context.Context, c <-chan Data) <-chan Data {
    out := make(chan Data)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                return
            case v, ok := <-c:
                if !ok {
                    return
                }
                select {
                case out <- v:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}

// 使用
for v := range OrDone(ctx, dataCh) {
    // 处理数据，自动处理 ctx 取消
}
```
