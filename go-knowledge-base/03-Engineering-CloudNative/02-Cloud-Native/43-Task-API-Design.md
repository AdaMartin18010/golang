# 任务 API 设计 (Task API Design)

> **分类**: 工程与云原生  
> **标签**: #api-design #rest #graphql

---

## REST API

```go
// API 路由定义
func RegisterTaskAPI(r *gin.Engine, service TaskService) {
    handler := &TaskHandler{service: service}
    
    v1 := r.Group("/api/v1")
    {
        tasks := v1.Group("/tasks")
        {
            tasks.GET("", handler.ListTasks)
            tasks.POST("", handler.CreateTask)
            tasks.GET("/:id", handler.GetTask)
            tasks.PATCH("/:id", handler.UpdateTask)
            tasks.DELETE("/:id", handler.CancelTask)
            
            // 任务操作
            tasks.POST("/:id/retry", handler.RetryTask)
            tasks.POST("/:id/pause", handler.PauseTask)
            tasks.POST("/:id/resume", handler.ResumeTask)
            
            // 任务日志
            tasks.GET("/:id/logs", handler.GetTaskLogs)
            tasks.GET("/:id/events", handler.GetTaskEvents)
        }
        
        // 批处理
        v1.POST("/batch/submit", handler.BatchSubmit)
        v1.POST("/batch/cancel", handler.BatchCancel)
    }
}

// 请求/响应结构
type CreateTaskRequest struct {
    Name        string            `json:"name" binding:"required"`
    Type        string            `json:"type" binding:"required"`
    Payload     json.RawMessage   `json:"payload"`
    Priority    int               `json:"priority"`
    ScheduleAt  *time.Time        `json:"schedule_at,omitempty"`
    Timeout     *time.Duration    `json:"timeout,omitempty"`
    Retries     *int              `json:"retries,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}

type TaskResponse struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        string            `json:"type"`
    Status      string            `json:"status"`
    Priority    int               `json:"priority"`
    Payload     json.RawMessage   `json:"payload,omitempty"`
    Result      json.RawMessage   `json:"result,omitempty"`
    Error       *ErrorInfo        `json:"error,omitempty"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
    StartedAt   *time.Time        `json:"started_at,omitempty"`
    CompletedAt *time.Time        `json:"completed_at,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}

type ListTasksResponse struct {
    Data       []TaskResponse `json:"data"`
    Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
    Total  int64 `json:"total"`
    Limit  int   `json:"limit"`
    Offset int   `json:"offset"`
}
```

---

## 分页与过滤

```go
func (th *TaskHandler) ListTasks(c *gin.Context) {
    var query ListTasksQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }
    
    // 构建过滤器
    filter := TaskFilter{
        Status:    query.Status,
        Type:      query.Type,
        Priority:  query.Priority,
        CreatedAfter: query.CreatedAfter,
        CreatedBefore: query.CreatedBefore,
        Metadata:  parseMetadataFilter(query.Metadata),
    }
    
    // 排序
    sort := SortOptions{
        Field: query.SortBy,
        Order: query.SortOrder,
    }
    
    // 分页
    pagination := PaginationOptions{
        Limit:  query.Limit,
        Offset: query.Offset,
        Cursor: query.Cursor,
    }
    
    result, err := th.service.List(c.Request.Context(), filter, sort, pagination)
    if err != nil {
        c.JSON(500, ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(200, ListTasksResponse{
        Data:       convertToResponses(result.Tasks),
        Pagination: result.Pagination,
    })
}

type ListTasksQuery struct {
    Status        string    `form:"status"`
    Type          string    `form:"type"`
    Priority      int       `form:"priority"`
    CreatedAfter  time.Time `form:"created_after"`
    CreatedBefore time.Time `form:"created_before"`
    Metadata      string    `form:"metadata"`
    SortBy        string    `form:"sort_by"`
    SortOrder     string    `form:"sort_order"`
    Limit         int       `form:"limit,default=20"`
    Offset        int       `form:"offset,default=0"`
    Cursor        string    `form:"cursor"`
}
```

---

## GraphQL API

```go
// Schema
const TaskSchema = `
type Task {
    id: ID!
    name: String!
    type: String!
    status: TaskStatus!
    priority: Int!
    payload: JSON
    result: JSON
    error: TaskError
    createdAt: Time!
    updatedAt: Time!
    scheduledAt: Time
    startedAt: Time
    completedAt: Time
    duration: Duration
    metadata: Map
    logs: [TaskLog!]
    events: [TaskEvent!]
}

enum TaskStatus {
    PENDING
    SCHEDULED
    RUNNING
    COMPLETED
    FAILED
    CANCELLED
    RETRYING
}

type Query {
    task(id: ID!): Task
    tasks(
        status: TaskStatus
        type: String
        priority: Int
        limit: Int
        offset: Int
    ): TaskConnection!
    taskStats: TaskStats!
}

type Mutation {
    createTask(input: CreateTaskInput!): Task!
    cancelTask(id: ID!): Task!
    retryTask(id: ID!): Task!
    pauseTask(id: ID!): Task!
    resumeTask(id: ID!): Task!
    batchSubmit(tasks: [CreateTaskInput!]!): [Task!]!
}

type Subscription {
    taskUpdated(id: ID): Task!
    taskStatsUpdated: TaskStats!
}
`

// Resolver
type TaskResolver struct {
    service TaskService
}

func (tr *TaskResolver) Task(ctx context.Context, args struct{ ID string }) (*Task, error) {
    return tr.service.Get(ctx, args.ID)
}

func (tr *TaskResolver) Tasks(ctx context.Context, args struct {
    Status *string
    Type   *string
    Limit  *int32
    Offset *int32
}) (*TaskConnection, error) {
    filter := TaskFilter{
        Status: dereference(args.Status),
        Type:   dereference(args.Type),
    }
    
    limit := 20
    if args.Limit != nil {
        limit = int(*args.Limit)
    }
    
    result, _ := tr.service.List(ctx, filter, SortOptions{}, PaginationOptions{Limit: limit})
    
    return &TaskConnection{
        Edges: convertToEdges(result.Tasks),
        PageInfo: PageInfo{
            TotalCount: int32(result.Pagination.Total),
        },
    }, nil
}
```

---

## API 限流

```go
type RateLimitedHandler struct {
    handler   http.Handler
    limiter   *rate.Limiter
    keyFunc   func(*http.Request) string
}

func (rlh *RateLimitedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    key := rlh.keyFunc(r)
    
    if !rlh.limiter.Allow(key) {
        w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rlh.limiter.Limit()))
        w.Header().Set("X-RateLimit-Remaining", "0")
        
        retryAfter := rlh.limiter.RetryAfter(key)
        w.Header().Set("Retry-After", fmt.Sprintf("%d", int(retryAfter.Seconds())))
        
        http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
        return
    }
    
    rlh.handler.ServeHTTP(w, r)
}

// 不同端点不同限流
func RateLimitByEndpoint() gin.HandlerFunc {
    limits := map[string]rate.Limit{
        "POST /api/v1/tasks":      rate.Limit(10, time.Minute),
        "GET /api/v1/tasks":       rate.Limit(100, time.Minute),
        "POST /api/v1/batch":      rate.Limit(5, time.Minute),
        "GET /api/v1/tasks/:id":   rate.Limit(200, time.Minute),
    }
    
    return func(c *gin.Context) {
        key := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
        limit, ok := limits[key]
        if !ok {
            c.Next()
            return
        }
        
        if !limit.Allow(c.ClientIP()) {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "rate limit exceeded",
            })
            return
        }
        
        c.Next()
    }
}
```
