# 任务 Web UI (Task Web UI)

> **分类**: 工程与云原生
> **标签**: #web-ui #dashboard #visualization

---

## 管理界面后端

```go
type TaskDashboardHandler struct {
    taskService TaskService
    statsService StatsService
}

func (tdh *TaskDashboardHandler) RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {
        api.GET("/tasks", tdh.ListTasks)
        api.GET("/tasks/:id", tdh.GetTask)
        api.POST("/tasks", tdh.CreateTask)
        api.DELETE("/tasks/:id", tdh.CancelTask)

        api.GET("/stats", tdh.GetStats)
        api.GET("/stats/realtime", tdh.GetRealtimeStats)

        api.GET("/workers", tdh.ListWorkers)
        api.GET("/queues", tdh.ListQueues)

        api.GET("/logs/:taskId", tdh.GetTaskLogs)
    }

    // WebSocket 实时更新
    r.GET("/ws", tdh.WebSocketHandler)
}

func (tdh *TaskDashboardHandler) ListTasks(c *gin.Context) {
    options := ListOptions{
        Status: c.Query("status"),
        Type:   c.Query("type"),
        Limit:  parseInt(c.DefaultQuery("limit", "20")),
        Offset: parseInt(c.DefaultQuery("offset", "0")),
    }

    tasks, total, _ := tdh.taskService.List(c.Request.Context(), options)

    c.JSON(200, gin.H{
        "data":  tasks,
        "total": total,
    })
}

func (tdh *TaskDashboardHandler) GetStats(c *gin.Context) {
    stats, _ := tdh.statsService.GetDashboardStats(c.Request.Context())
    c.JSON(200, stats)
}

func (tdh *TaskDashboardHandler) WebSocketHandler(c *gin.Context) {
    conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
    if err != nil {
        return
    }
    defer conn.Close()

    // 订阅实时更新
    events := tdh.statsService.Subscribe()
    defer tdh.statsService.Unsubscribe(events)

    for event := range events {
        if err := conn.WriteJSON(event); err != nil {
            break
        }
    }
}
```

---

## 仪表盘数据

```go
type DashboardStats struct {
    Overview    TaskOverview    `json:"overview"`
    Performance PerformanceStats `json:"performance"`
    Workers     []WorkerStatus  `json:"workers"`
    Queues      []QueueStatus   `json:"queues"`
    Timeline    []TimelinePoint `json:"timeline"`
}

type TaskOverview struct {
    Total      int64 `json:"total"`
    Pending    int64 `json:"pending"`
    Running    int64 `json:"running"`
    Completed  int64 `json:"completed"`
    Failed     int64 `json:"failed"`
    Retrying   int64 `json:"retrying"`
}

type PerformanceStats struct {
    Throughput    float64 `json:"throughput"`     // tasks/sec
    AvgLatency    float64 `json:"avg_latency_ms"`
    P99Latency    float64 `json:"p99_latency_ms"`
    ErrorRate     float64 `json:"error_rate"`
    QueueDepth    int     `json:"queue_depth"`
}

func (ss *StatsService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
    return &DashboardStats{
        Overview:    ss.getOverview(ctx),
        Performance: ss.getPerformance(ctx),
        Workers:     ss.getWorkers(ctx),
        Queues:      ss.getQueues(ctx),
        Timeline:    ss.getTimeline(ctx, time.Hour),
    }, nil
}
```

---

## 前端组件接口

```go
// 为 React/Vue 提供的数据格式
type TaskTableRow struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Status      string    `json:"status"`
    Priority    int       `json:"priority"`
    Progress    float64   `json:"progress"`
    CreatedAt   time.Time `json:"created_at"`
    StartedAt   *time.Time `json:"started_at,omitempty"`
    Duration    *float64  `json:"duration_ms,omitempty"`
    WorkerID    string    `json:"worker_id,omitempty"`
    Actions     []Action  `json:"actions"`
}

type Action struct {
    Name    string `json:"name"`
    Icon    string `json:"icon"`
    Action  string `json:"action"`
    Enabled bool   `json:"enabled"`
}

// 图表数据
type ChartData struct {
    Labels []string    `json:"labels"`
    Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
    Label string    `json:"label"`
    Data  []float64 `json:"data"`
    Color string    `json:"color"`
}

func (tdh *TaskDashboardHandler) GetTimelineData(c *gin.Context) {
    rangeParam := c.DefaultQuery("range", "1h")
    duration := parseDuration(rangeParam)

    timeline, _ := tdh.statsService.GetTimeline(c.Request.Context(), duration)

    c.JSON(200, ChartData{
        Labels: extractLabels(timeline),
        Datasets: []Dataset{
            {
                Label: "Submitted",
                Data:  extractValues(timeline, "submitted"),
                Color: "#3498db",
            },
            {
                Label: "Completed",
                Data:  extractValues(timeline, "completed"),
                Color: "#2ecc71",
            },
            {
                Label: "Failed",
                Data:  extractValues(timeline, "failed"),
                Color: "#e74c3c",
            },
        },
    })
}
```
