# 任务可观测性 (Task Observability)

> **分类**: 工程与云原生
> **标签**: #observability #tracing #metrics #logging

---

## 分布式追踪

```go
import "go.opentelemetry.io/otel"

func (e *TaskExecutor) executeWithTracing(ctx context.Context, task *Task) error {
    tracer := otel.Tracer("task-executor")

    ctx, span := tracer.Start(ctx, fmt.Sprintf("execute-task-%s", task.Type),
        trace.WithAttributes(
            attribute.String("task.id", task.ID),
            attribute.String("task.name", task.Name),
            attribute.String("task.type", task.Type),
            attribute.Int("task.priority", task.Priority),
        ),
    )
    defer span.End()

    // 记录开始
    span.AddEvent("task started")

    // 执行
    err := e.execute(ctx, task)

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    } else {
        span.AddEvent("task completed")
        span.SetStatus(codes.Ok, "success")
    }

    return err
}

// 子任务追踪
func (e *TaskExecutor) executeSubTask(ctx context.Context, parentTaskID string, subTask *SubTask) error {
    tracer := otel.Tracer("task-executor")

    // 创建子span
    ctx, span := tracer.Start(ctx, fmt.Sprintf("subtask-%s", subTask.Name),
        trace.WithAttributes(
            attribute.String("subtask.name", subTask.Name),
            attribute.String("parent.task.id", parentTaskID),
        ),
    )
    defer span.End()

    return e.runSubTask(ctx, subTask)
}
```

---

## 任务指标

```go
var (
    taskExecutions = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "task_executions_total",
            Help: "Total task executions",
        },
        []string{"type", "status"},
    )

    taskDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "task_duration_seconds",
            Help:    "Task execution duration",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
        },
        []string{"type"},
    )

    taskQueueWait = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "task_queue_wait_seconds",
            Help:    "Time spent waiting in queue",
            Buckets: prometheus.ExponentialBuckets(0.01, 2, 10),
        },
    )

    activeTasksGauge = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "active_tasks",
            Help: "Number of active tasks",
        },
        []string{"type"},
    )
)

func recordTaskMetrics(task *Task, status string, duration time.Duration) {
    taskExecutions.WithLabelValues(task.Type, status).Inc()
    taskDuration.WithLabelValues(task.Type).Observe(duration.Seconds())
}
```

---

## 任务链路日志

```go
type TaskLogger struct {
    logger *zap.Logger
}

func (tl *TaskLogger) LogTaskLifecycle(ctx context.Context, task *Task) {
    fields := []zap.Field{
        zap.String("task_id", task.ID),
        zap.String("task_type", task.Type),
        zap.String("trace_id", TraceIDFromContext(ctx)),
        zap.String("span_id", SpanIDFromContext(ctx)),
    }

    tl.logger.Info("task lifecycle event",
        append(fields,
            zap.String("event", "created"),
            zap.Time("timestamp", time.Now()),
        )...,
    )
}

func (tl *TaskLogger) LogTaskExecution(ctx context.Context, task *Task, stage string, fields ...zap.Field) {
    baseFields := []zap.Field{
        zap.String("task_id", task.ID),
        zap.String("stage", stage),
        zap.String("trace_id", TraceIDFromContext(ctx)),
    }

    tl.logger.Info(fmt.Sprintf("task %s", stage),
        append(baseFields, fields...)...,
    )
}

// 使用
func processTask(ctx context.Context, task *Task) {
    logger := TaskLogger{logger: zap.L()}

    logger.LogTaskExecution(ctx, task, "started")

    if err := doWork(ctx); err != nil {
        logger.LogTaskExecution(ctx, task, "failed", zap.Error(err))
        return
    }

    logger.LogTaskExecution(ctx, task, "completed",
        zap.Duration("duration", time.Since(start)))
}
```

---

## 任务可视化

```go
type TaskVisualizer struct {
    store TaskStore
}

func (tv *TaskVisualizer) GetTaskFlow(taskID string) (*TaskFlow, error) {
    task, _ := tv.store.Get(taskID)

    flow := &TaskFlow{
        Root: task.ID,
        Nodes: make(map[string]*FlowNode),
    }

    // 构建流程图
    tv.buildFlow(task, flow, 0)

    return flow, nil
}

func (tv *TaskVisualizer) buildFlow(task *Task, flow *TaskFlow, depth int) {
    node := &FlowNode{
        ID:        task.ID,
        Name:      task.Name,
        Status:    task.Status,
        StartTime: task.StartTime,
        EndTime:   task.EndTime,
        Depth:     depth,
    }

    flow.Nodes[task.ID] = node

    // 递归子任务
    for _, subTaskID := range task.SubTaskIDs {
        if subTask, err := tv.store.Get(subTaskID); err == nil {
            tv.buildFlow(subTask, flow, depth+1)
            flow.Edges = append(flow.Edges, FlowEdge{
                From: task.ID,
                To:   subTaskID,
            })
        }
    }
}

// 生成 Gantt 图数据
func (tv *TaskVisualizer) GetGanttData(timeRange TimeRange) ([]GanttItem, error) {
    tasks, _ := tv.store.ListInRange(timeRange)

    var items []GanttItem
    for _, task := range tasks {
        if task.StartTime != nil && task.EndTime != nil {
            items = append(items, GanttItem{
                ID:       task.ID,
                Name:     task.Name,
                Start:    *task.StartTime,
                End:      *task.EndTime,
                Progress: task.Progress,
                Status:   task.Status,
            })
        }
    }

    return items, nil
}
```
