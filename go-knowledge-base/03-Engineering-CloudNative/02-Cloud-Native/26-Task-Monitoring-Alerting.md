# 任务监控与告警 (Task Monitoring & Alerting)

> **分类**: 工程与云原生  
> **标签**: #monitoring #alerting #observability

---

## 任务指标收集

```go
type TaskMetrics struct {
    registry prometheus.Registerer
}

var (
    taskTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "tasks_total",
            Help: "Total number of tasks",
        },
        []string{"type", "status"},
    )
    
    taskDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "task_duration_seconds",
            Help:    "Task execution duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"type"},
    )
    
    taskQueueSize = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "task_queue_size",
            Help: "Current task queue size",
        },
        []string{"queue"},
    )
    
    activeTasks = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_tasks",
            Help: "Number of currently running tasks",
        },
    )
    
    taskRetries = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "task_retries_total",
            Help: "Total number of task retries",
        },
        []string{"type"},
    )
)

func (tm *TaskMetrics) RecordTaskStart(taskType string) {
    activeTasks.Inc()
}

func (tm *TaskMetrics) RecordTaskComplete(taskType string, status string, duration time.Duration) {
    activeTasks.Dec()
    taskTotal.WithLabelValues(taskType, status).Inc()
    taskDuration.WithLabelValues(taskType).Observe(duration.Seconds())
}

func (tm *TaskMetrics) RecordRetry(taskType string) {
    taskRetries.WithLabelValues(taskType).Inc()
}

func (tm *TaskMetrics) SetQueueSize(queueName string, size int) {
    taskQueueSize.WithLabelValues(queueName).Set(float64(size))
}
```

---

## 健康检查

```go
type TaskHealthChecker struct {
    executor *TaskExecutor
    thresholds HealthThresholds
}

type HealthThresholds struct {
    MaxQueueSize      int
    MaxActiveTasks    int
    MaxErrorRate      float64
    MaxAvgDuration    time.Duration
}

func (thc *TaskHealthChecker) CheckHealth(ctx context.Context) HealthStatus {
    status := HealthStatus{Healthy: true}
    
    // 检查队列大小
    for queue, size := range thc.executor.GetQueueSizes() {
        if size > thc.thresholds.MaxQueueSize {
            status.Healthy = false
            status.Issues = append(status.Issues, 
                fmt.Sprintf("Queue %s size %d exceeds threshold %d", 
                    queue, size, thc.thresholds.MaxQueueSize))
        }
    }
    
    // 检查活跃任务数
    active := thc.executor.GetActiveTaskCount()
    if active > thc.thresholds.MaxActiveTasks {
        status.Healthy = false
        status.Issues = append(status.Issues,
            fmt.Sprintf("Active tasks %d exceeds threshold %d",
                active, thc.thresholds.MaxActiveTasks))
    }
    
    // 检查错误率
    errorRate := thc.calculateErrorRate()
    if errorRate > thc.thresholds.MaxErrorRate {
        status.Healthy = false
        status.Issues = append(status.Issues,
            fmt.Sprintf("Error rate %.2f%% exceeds threshold %.2f%%",
                errorRate*100, thc.thresholds.MaxErrorRate*100))
    }
    
    return status
}
```

---

## 告警规则

```go
type AlertRule struct {
    Name        string
    Condition   func(metrics Metrics) bool
    Severity    string
    Message     string
    Cooldown    time.Duration
    LastFired   time.Time
}

var alertRules = []AlertRule{
    {
        Name: "high_task_failure_rate",
        Condition: func(m Metrics) bool {
            return m.TaskFailureRate > 0.1  // 10%
        },
        Severity: "warning",
        Message:  "Task failure rate is above 10%",
        Cooldown: 5 * time.Minute,
    },
    {
        Name: "task_queue_backlog",
        Condition: func(m Metrics) bool {
            return m.QueueWaitTime > 5*time.Minute
        },
        Severity: "critical",
        Message:  "Tasks are waiting in queue for more than 5 minutes",
        Cooldown: 1 * time.Minute,
    },
    {
        Name: "long_running_tasks",
        Condition: func(m Metrics) bool {
            return m.MaxTaskDuration > 30*time.Minute
        },
        Severity: "warning",
        Message:  "Some tasks have been running for more than 30 minutes",
        Cooldown: 10 * time.Minute,
    },
}

func (ar *AlertRule) Evaluate(metrics Metrics) *Alert {
    if !ar.Condition(metrics) {
        return nil
    }
    
    // 检查冷却期
    if time.Since(ar.LastFired) < ar.Cooldown {
        return nil
    }
    
    ar.LastFired = time.Now()
    
    return &Alert{
        Rule:      ar.Name,
        Severity:  ar.Severity,
        Message:   ar.Message,
        Timestamp: time.Now(),
        Metrics:   metrics,
    }
}
```

---

## 实时仪表板

```go
// WebSocket 推送实时数据
type DashboardHub struct {
    clients map[*Client]bool
    broadcast chan TaskEvent
}

func (hub *DashboardHub) Run() {
    for {
        select {
        case event := <-hub.broadcast:
            for client := range hub.clients {
                select {
                case client.send <- event:
                default:
                    close(client.send)
                    delete(hub.clients, client)
                }
            }
        }
    }
}

func (hub *DashboardHub) BroadcastTaskEvent(event TaskEvent) {
    hub.broadcast <- event
}

// 事件类型
type TaskEvent struct {
    Type      string    `json:"type"`      // started, completed, failed
    TaskID    string    `json:"task_id"`
    TaskType  string    `json:"task_type"`
    Timestamp time.Time `json:"timestamp"`
    Duration  float64   `json:"duration,omitempty"`
    Error     string    `json:"error,omitempty"`
}
```

---

## 任务追踪

```go
type TaskTracer struct {
    spans map[string]*TaskSpan
}

type TaskSpan struct {
    TaskID      string
    ParentID    string
    Name        string
    StartTime   time.Time
    EndTime     *time.Time
    Tags        map[string]string
    Logs        []LogEntry
}

func (tt *TaskTracer) StartSpan(ctx context.Context, taskID, name string) context.Context {
    span := &TaskSpan{
        TaskID:    taskID,
        Name:      name,
        StartTime: time.Now(),
        Tags:      make(map[string]string),
    }
    
    // 提取父span
    if parent := SpanFromContext(ctx); parent != nil {
        span.ParentID = parent.TaskID
    }
    
    tt.spans[taskID] = span
    
    return WithSpan(ctx, span)
}

func (tt *TaskTracer) FinishSpan(ctx context.Context, taskID string) {
    if span, ok := tt.spans[taskID]; ok {
        now := time.Now()
        span.EndTime = &now
        
        // 记录到存储
        tt.store.SaveSpan(ctx, span)
    }
}

func (tt *TaskTracer) AddLog(ctx context.Context, taskID, message string) {
    if span, ok := tt.spans[taskID]; ok {
        span.Logs = append(span.Logs, LogEntry{
            Timestamp: time.Now(),
            Message:   message,
        })
    }
}
```
