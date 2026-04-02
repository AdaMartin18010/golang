# 任务调度 (Job Scheduling)

> **分类**: 工程与云原生  
> **标签**: #scheduler #cron #distributed-job

---

## Cron 表达式调度

### robfig/cron

```go
import "github.com/robfig/cron/v3"

c := cron.New()

// 每分钟执行
c.AddFunc("*/1 * * * *", func() {
    fmt.Println("Every minute")
})

// 每小时执行
c.AddFunc("0 * * * *", func() {
    fmt.Println("Every hour")
})

// 每天凌晨2点
c.AddFunc("0 2 * * *", func() {
    fmt.Println("2 AM daily")
})

// 工作日每10分钟
c.AddFunc("*/10 * * * 1-5", func() {
    fmt.Println("Every 10 minutes on weekdays")
})

c.Start()
```

---

## 分布式任务调度

### 基于 Redis 的分布式锁

```go
type DistributedScheduler struct {
    redis      *redis.Client
    nodeID     string
    lockTTL    time.Duration
}

func (s *DistributedScheduler) Schedule(ctx context.Context, job Job) error {
    // 获取分布式锁
    lockKey := fmt.Sprintf("scheduler:lock:%s", job.Name)
    
    acquired, err := s.redis.SetNX(ctx, lockKey, s.nodeID, s.lockTTL).Result()
    if err != nil || !acquired {
        return fmt.Errorf("could not acquire lock: %w", err)
    }
    
    // 启动续期 goroutine
    stopRenew := make(chan struct{})
    go s.renewLock(ctx, lockKey, stopRenew)
    defer close(stopRenew)
    
    // 执行任务
    return s.executeJob(ctx, job)
}

func (s *DistributedScheduler) renewLock(ctx context.Context, key string, stop <-chan struct{}) {
    ticker := time.NewTicker(s.lockTTL / 3)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            s.redis.Expire(ctx, key, s.lockTTL)
        case <-stop:
            return
        case <-ctx.Done():
            return
        }
    }
}
```

---

## 任务队列与工作池

```go
type Job struct {
    ID      string
    Type    string
    Payload interface{}
    Context context.Context
}

type WorkerPool struct {
    workers   int
    jobQueue  chan Job
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewWorkerPool(workers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        workers:  workers,
        jobQueue: make(chan Job, workers*2),
        ctx:      ctx,
        cancel:   cancel,
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()
    
    for {
        select {
        case job := <-p.jobQueue:
            p.processJob(job)
        case <-p.ctx.Done():
            // 处理剩余任务
            for {
                select {
                case job := <-p.jobQueue:
                    p.processJob(job)
                default:
                    return
                }
            }
        }
    }
}

func (p *WorkerPool) processJob(job Job) {
    // 合并任务上下文和工作池上下文
    ctx, cancel := context.WithTimeout(job.Context, 30*time.Second)
    defer cancel()
    
    if err := executeWithContext(ctx, job); err != nil {
        log.Printf("Job %s failed: %v", job.ID, err)
    }
}

func (p *WorkerPool) Submit(job Job) error {
    select {
    case p.jobQueue <- job:
        return nil
    case <-time.After(5 * time.Second):
        return fmt.Errorf("job queue full")
    }
}

func (p *WorkerPool) Stop() {
    p.cancel()
    p.wg.Wait()
}
```

---

## 延迟任务

```go
type DelayedQueue struct {
    redis *redis.Client
}

func (q *DelayedQueue) Push(ctx context.Context, job Job, delay time.Duration) error {
    executeAt := time.Now().Add(delay)
    
    data, _ := json.Marshal(job)
    
    // 使用 Redis ZSet，score 为执行时间戳
    return q.redis.ZAdd(ctx, "delayed:jobs", &redis.Z{
        Score:  float64(executeAt.Unix()),
        Member: string(data),
    }).Err()
}

func (q *DelayedQueue) Poll(ctx context.Context) (*Job, error) {
    for {
        // 获取已到期的任务
        now := float64(time.Now().Unix())
        
        result, err := q.redis.ZRangeByScoreWithScores(ctx, "delayed:jobs", &redis.ZRangeBy{
            Min:   "0",
            Max:   fmt.Sprintf("%f", now),
            Count: 1,
        }).Result()
        
        if err != nil || len(result) == 0 {
            time.Sleep(100 * time.Millisecond)
            continue
        }
        
        // 移除并返回
        data := result[0].Member.(string)
        q.redis.ZRem(ctx, "delayed:jobs", data)
        
        var job Job
        json.Unmarshal([]byte(data), &job)
        return &job, nil
    }
}
```

---

## 任务重试与退避

```go
type RetryPolicy struct {
    MaxRetries  int
    InitialDelay time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
}

func (p *RetryPolicy) CalculateDelay(attempt int) time.Duration {
    if attempt >= p.MaxRetries {
        return 0
    }
    
    delay := float64(p.InitialDelay) * math.Pow(p.Multiplier, float64(attempt))
    if delay > float64(p.MaxDelay) {
        delay = float64(p.MaxDelay)
    }
    
    // 添加抖动
    jitter := rand.Float64() * 0.3 * delay
    delay += jitter
    
    return time.Duration(delay)
}

func ExecuteWithRetry(ctx context.Context, fn func() error, policy RetryPolicy) error {
    var err error
    
    for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
        err = fn()
        if err == nil {
            return nil
        }
        
        if attempt == policy.MaxRetries {
            break
        }
        
        delay := policy.CalculateDelay(attempt)
        
        select {
        case <-time.After(delay):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    
    return fmt.Errorf("max retries exceeded: %w", err)
}
```
