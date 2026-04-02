# 任务性能调优 (Task Performance Tuning)

> **分类**: 工程与云原生
> **标签**: #performance #optimization #tuning

---

## 性能基准测试

```go
// 任务执行基准
type TaskBenchmark struct {
    executor *TaskExecutor
}

func (tb *TaskBenchmark) Run(b *testing.B, taskType string) {
    task := &Task{
        Type:    taskType,
        Payload: []byte(`{"test": "data"}`),
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ctx := context.Background()
            tb.executor.Execute(ctx, task)
        }
    })
}

// 调度延迟基准
func BenchmarkSchedulerLatency(b *testing.B) {
    scheduler := NewTaskScheduler()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        task := &Task{
            Type: "latency-test",
        }
        start := time.Now()
        scheduler.Schedule(context.Background(), task)
        latency := time.Since(start)

        // 记录延迟分布
        recordLatency(latency)
    }
}

// 吞吐量测试
func TestThroughput(t *testing.T) {
    scheduler := NewTaskScheduler(Config{
        Workers:    100,
        QueueSize:  10000,
    })

    ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    defer cancel()

    var completed int64
    var wg sync.WaitGroup

    // 启动生产者
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for ctx.Err() == nil {
                task := &Task{Type: "throughput-test"}
                scheduler.Schedule(ctx, task)
                atomic.AddInt64(&completed, 1)
            }
        }()
    }

    <-ctx.Done()
    wg.Wait()

    tps := float64(atomic.LoadInt64(&completed)) / 60
    t.Logf("Throughput: %.2f tasks/sec", tps)
}
```

---

## 性能分析

```go
type PerformanceProfiler struct {
    pprof    *pprof.Profile
    trace    *trace.Tracer
}

func (pp *PerformanceProfiler) ProfileCPU(ctx context.Context, duration time.Duration) ([]byte, error) {
    var buf bytes.Buffer

    if err := pprof.StartCPUProfile(&buf); err != nil {
        return nil, err
    }
    defer pprof.StopCPUProfile()

    select {
    case <-time.After(duration):
    case <-ctx.Done():
        return nil, ctx.Err()
    }

    return buf.Bytes(), nil
}

func (pp *PerformanceProfiler) ProfileMemory() ([]byte, error) {
    var buf bytes.Buffer

    runtime.GC() // 先 GC 获取准确数据
    if err := pprof.WriteHeapProfile(&buf); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

func (pp *PerformanceProfiler) TraceExecution(ctx context.Context, duration time.Duration) ([]byte, error) {
    var buf bytes.Buffer

    if err := trace.Start(&buf); err != nil {
        return nil, err
    }
    defer trace.Stop()

    // 执行一些任务
    <-time.After(duration)

    return buf.Bytes(), nil
}
```

---

## 瓶颈分析

```go
type BottleneckAnalyzer struct {
    metrics MetricsCollector
}

func (ba *BottleneckAnalyzer) Analyze() *BottleneckReport {
    report := &BottleneckReport{
        Timestamp: time.Now(),
    }

    // 1. 检查队列深度
    queueDepth := ba.metrics.GetQueueDepth()
    if queueDepth > 1000 {
        report.Bottlenecks = append(report.Bottlenecks, Bottleneck{
            Component: "task_queue",
            Severity:  "high",
            Metric:    queueDepth,
            Threshold: 1000,
            Suggestion: "Increase worker count or queue size",
        })
    }

    // 2. 检查 worker 利用率
    workerUtil := ba.metrics.GetWorkerUtilization()
    if workerUtil > 0.95 {
        report.Bottlenecks = append(report.Bottlenecks, Bottleneck{
            Component: "worker_pool",
            Severity:  "high",
            Metric:    workerUtil,
            Threshold: 0.95,
            Suggestion: "Add more workers or optimize task handlers",
        })
    }

    // 3. 检查数据库连接
    dbWaitTime := ba.metrics.GetDBWaitTime()
    if dbWaitTime > 100*time.Millisecond {
        report.Bottlenecks = append(report.Bottlenecks, Bottleneck{
            Component: "database",
            Severity:  "medium",
            Metric:    dbWaitTime,
            Threshold: 100 * time.Millisecond,
            Suggestion: "Increase connection pool size or optimize queries",
        })
    }

    // 4. 检查 GC 压力
    gcPause := ba.metrics.GetGCPauseTime()
    if gcPause > 10*time.Millisecond {
        report.Bottlenecks = append(report.Bottlenecks, Bottleneck{
            Component: "gc",
            Severity:  "low",
            Metric:    gcPause,
            Threshold: 10 * time.Millisecond,
            Suggestion: "Reduce allocations or tune GC",
        })
    }

    return report
}
```

---

## 性能优化技巧

```go
// 1. 对象池化
var taskPool = sync.Pool{
    New: func() interface{} {
        return &Task{}
    },
}

func getTask() *Task {
    return taskPool.Get().(*Task)
}

func putTask(t *Task) {
    t.Reset()
    taskPool.Put(t)
}

// 2. 批处理提交
func (te *TaskExecutor) BatchSubmit(tasks []*Task) error {
    // 批量写入，减少锁竞争
    te.queueMu.Lock()
    for _, task := range tasks {
        te.queue = append(te.queue, task)
    }
    te.queueMu.Unlock()

    // 批量通知
    te.cond.Broadcast()
    return nil
}

// 3. 无锁队列
func (te *TaskExecutor) useLockFreeQueue() {
    te.queue = NewRingBuffer(10000)
}

// 4. 内存预分配
func (t *Task) PreallocatePayload(size int) {
    if cap(t.Payload) < size {
        t.Payload = make([]byte, 0, size)
    }
}

// 5. 并行处理
func (tp *TaskProcessor) ProcessBatch(tasks []Task) []Result {
    results := make([]Result, len(tasks))

    parallel.ForEach(tasks, func(i int, task Task) {
        results[i] = tp.processSingle(task)
    })

    return results
}
```
