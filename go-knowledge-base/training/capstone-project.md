# Capstone Project: Distributed Task Scheduler

## Project Overview

**Duration:** 1 week (40 hours)
**Prerequisites:** Completion of Week 1-4 training
**Goal:** Demonstrate mastery of Go fundamentals, concurrency, cloud-native patterns, and system design

---

## Project Description

Build a production-ready distributed task scheduler that can:

- Schedule millions of tasks with cron expressions
- Execute tasks reliably across multiple worker nodes
- Handle failures gracefully with automatic retries
- Provide comprehensive observability
- Scale horizontally

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         Clients                             │
│                    (CLI, REST API, SDK)                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    API Gateway                              │
│              (Rate Limiting, Auth, Routing)                 │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
┌───────▼──────┐ ┌─────▼──────┐ ┌────▼─────┐
│  Scheduler   │ │  Worker    │ │  Admin   │
│   Service    │ │   Nodes    │ │   API    │
└───────┬──────┘ └─────┬──────┘ └────┬─────┘
        │              │             │
        └──────────────┼─────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
┌───────▼──────┐ ┌─────▼──────┐ ┌────▼─────┐
│ PostgreSQL   │ │   Redis    │ │  Kafka   │
│ (Task Store) │ │  (Queue)   │ │ (Events) │
└──────────────┘ └────────────┘ └──────────┘
```

---

## Requirements

### Functional Requirements

1. **Task Management**
   - Create tasks with payload and schedule
   - Support cron expressions for recurring tasks
   - Update and delete scheduled tasks
   - Query task status and history

2. **Task Execution**
   - Execute tasks at scheduled time
   - Support multiple task types (HTTP, gRPC, custom)
   - Execute tasks concurrently
   - Distribute across worker nodes

3. **Reliability**
   - Automatic retries with exponential backoff
   - Dead letter queue for failed tasks
   - At-least-once delivery guarantee
   - Graceful handling of worker failures

4. **Observability**
   - Task execution metrics
   - Structured logging
   - Distributed tracing
   - Health checks

### Non-Functional Requirements

1. **Performance**
   - Schedule 10,000 tasks/second
   - Execute 50,000 tasks/second
   - P99 latency < 100ms for API

2. **Scalability**
   - Horizontal scaling of scheduler
   - Horizontal scaling of workers
   - Handle 1M scheduled tasks

3. **Availability**
   - 99.9% uptime target
   - No single point of failure
   - Graceful degradation

4. **Security**
   - API authentication
   - Task payload encryption at rest
   - Secure inter-service communication

---

## Core Components

### 1. Task Domain Model

```go
package domain

import (
    "context"
    "time"
)

// Task represents a schedulable unit of work
type Task struct {
    ID          string
    Name        string
    Type        TaskType
    Status      TaskStatus
    Schedule    Schedule
    Payload     Payload
    RetryPolicy RetryPolicy
    Metadata    Metadata
    CreatedAt   time.Time
    UpdatedAt   time.Time
    NextRunAt   *time.Time
    LastRunAt   *time.Time
}

type TaskType string

const (
    TaskTypeHTTP   TaskType = "http"
    TaskTypeGRPC   TaskType = "grpc"
    TaskTypeCustom TaskType = "custom"
)

type TaskStatus string

const (
    TaskStatusPending    TaskStatus = "pending"
    TaskStatusScheduled  TaskStatus = "scheduled"
    TaskStatusRunning    TaskStatus = "running"
    TaskStatusCompleted  TaskStatus = "completed"
    TaskStatusFailed     TaskStatus = "failed"
    TaskStatusCancelled  TaskStatus = "cancelled"
)

type Schedule struct {
    Type      ScheduleType
    CronExpr  string
    Delay     *time.Duration
    RunAt     *time.Time
    Timezone  string
}

type ScheduleType string

const (
    ScheduleTypeOnce   ScheduleType = "once"
    ScheduleTypeCron   ScheduleType = "cron"
    ScheduleTypeDelay  ScheduleType = "delay"
)

type Payload struct {
    ContentType string
    Data        []byte
}

type RetryPolicy struct {
    MaxRetries  int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}

type Metadata struct {
    Tags        []string
    Priority    int
    Timeout     time.Duration
    UserID      string
    CorrelationID string
}

// ExecutionResult captures task execution outcome
type ExecutionResult struct {
    TaskID      string
    ExecutionID string
    Status      ExecutionStatus
    StartedAt   time.Time
    CompletedAt time.Time
    Output      []byte
    Error       *ExecutionError
}

type ExecutionStatus string

const (
    ExecutionStatusSuccess ExecutionStatus = "success"
    ExecutionStatusFailure ExecutionStatus = "failure"
    ExecutionStatusTimeout ExecutionStatus = "timeout"
)

type ExecutionError struct {
    Code    string
    Message string
    Details map[string]interface{}
}

// TaskRepository defines task persistence operations
type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (*Task, error)
    List(ctx context.Context, filter TaskFilter) ([]*Task, error)
    GetDueTasks(ctx context.Context, before time.Time, limit int) ([]*Task, error)
    UpdateStatus(ctx context.Context, id string, status TaskStatus) error
}

type TaskFilter struct {
    Status   *TaskStatus
    Type     *TaskType
    UserID   string
    Tags     []string
    FromTime *time.Time
    ToTime   *time.Time
    Limit    int
    Offset   int
}
```

### 2. Scheduler Service

```go
package scheduler

import (
    "context"
    "time"

    "github.com/robfig/cron/v3"
)

// Service manages task scheduling
type Service struct {
    repo        domain.TaskRepository
    queue       Queue
    cronParser  cron.Parser
    mu          sync.RWMutex
    schedules   map[string]cron.Schedule
    stopCh      chan struct{}
}

func NewService(repo domain.TaskRepository, queue Queue) *Service {
    return &Service{
        repo:       repo,
        queue:      queue,
        cronParser: cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
        schedules:  make(map[string]cron.Schedule),
        stopCh:     make(chan struct{}),
    }
}

// Start begins the scheduling loop
func (s *Service) Start(ctx context.Context) error {
    // Load existing scheduled tasks
    tasks, err := s.repo.List(ctx, domain.TaskFilter{
        Status: &domain.TaskStatusScheduled,
    })
    if err != nil {
        return err
    }

    for _, task := range tasks {
        if err := s.scheduleTask(task); err != nil {
            log.Printf("Failed to schedule task %s: %v", task.ID, err)
        }
    }

    // Start scheduling loop
    go s.schedulingLoop(ctx)

    return nil
}

func (s *Service) schedulingLoop(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-s.stopCh:
            return
        case now := <-ticker.C:
            s.checkAndDispatch(ctx, now)
        }
    }
}

func (s *Service) checkAndDispatch(ctx context.Context, now time.Time) {
    tasks, err := s.repo.GetDueTasks(ctx, now, 100)
    if err != nil {
        log.Printf("Failed to get due tasks: %v", err)
        return
    }

    for _, task := range tasks {
        if err := s.dispatchTask(ctx, task); err != nil {
            log.Printf("Failed to dispatch task %s: %v", task.ID, err)
            continue
        }
    }
}

func (s *Service) dispatchTask(ctx context.Context, task *domain.Task) error {
    // Create execution record
    execution := &domain.Execution{
        ID:        uuid.New().String(),
        TaskID:    task.ID,
        Status:    domain.ExecutionStatusPending,
        CreatedAt: time.Now(),
    }

    // Push to queue
    if err := s.queue.Push(ctx, execution); err != nil {
        return err
    }

    // Update task status
    if task.Schedule.Type == domain.ScheduleTypeOnce {
        task.Status = domain.TaskStatusRunning
    } else {
        // Calculate next run
        nextRun := s.calculateNextRun(task)
        task.NextRunAt = &nextRun
    }

    return s.repo.Update(ctx, task)
}

func (s *Service) calculateNextRun(task *domain.Task) time.Time {
    s.mu.RLock()
    schedule, ok := s.schedules[task.ID]
    s.mu.RUnlock()

    if !ok {
        // Parse schedule
        sched, err := s.cronParser.Parse(task.Schedule.CronExpr)
        if err != nil {
            return time.Time{}
        }
        schedule = sched
    }

    return schedule.Next(time.Now())
}

// CreateTask creates and schedules a new task
func (s *Service) CreateTask(ctx context.Context, req CreateTaskRequest) (*domain.Task, error) {
    task := &domain.Task{
        ID:        uuid.New().String(),
        Name:      req.Name,
        Type:      req.Type,
        Status:    domain.TaskStatusScheduled,
        Schedule:  req.Schedule,
        Payload:   req.Payload,
        RetryPolicy: req.RetryPolicy,
        Metadata:  req.Metadata,
        CreatedAt: time.Now(),
    }

    // Validate and calculate initial run time
    if err := s.validateAndPrepare(task); err != nil {
        return nil, err
    }

    // Persist
    if err := s.repo.Create(ctx, task); err != nil {
        return nil, err
    }

    // Schedule
    if err := s.scheduleTask(task); err != nil {
        return nil, err
    }

    return task, nil
}

func (s *Service) validateAndPrepare(task *domain.Task) error {
    switch task.Schedule.Type {
    case domain.ScheduleTypeOnce:
        if task.Schedule.RunAt == nil {
            now := time.Now()
            task.Schedule.RunAt = &now
        }
        task.NextRunAt = task.Schedule.RunAt

    case domain.ScheduleTypeCron:
        if task.Schedule.CronExpr == "" {
            return fmt.Errorf("cron expression required")
        }
        sched, err := s.cronParser.Parse(task.Schedule.CronExpr)
        if err != nil {
            return fmt.Errorf("invalid cron expression: %w", err)
        }
        nextRun := sched.Next(time.Now())
        task.NextRunAt = &nextRun

        s.mu.Lock()
        s.schedules[task.ID] = sched
        s.mu.Unlock()

    case domain.ScheduleTypeDelay:
        if task.Schedule.Delay == nil {
            return fmt.Errorf("delay duration required")
        }
        nextRun := time.Now().Add(*task.Schedule.Delay)
        task.NextRunAt = &nextRun
    }

    // Set defaults
    if task.RetryPolicy.MaxRetries == 0 {
        task.RetryPolicy.MaxRetries = 3
    }
    if task.RetryPolicy.InitialDelay == 0 {
        task.RetryPolicy.InitialDelay = 1 * time.Second
    }
    if task.RetryPolicy.Multiplier == 0 {
        task.RetryPolicy.Multiplier = 2.0
    }

    return nil
}
```

### 3. Worker Implementation

```go
package worker

import (
    "context"
    "fmt"
    "time"
)

// Worker executes tasks from the queue
type Worker struct {
    id           string
    queue        Queue
    executor     Executor
    resultStore  ResultStore
    stopCh       chan struct{}
    maxConcurrent int
}

func NewWorker(id string, queue Queue, executor Executor, resultStore ResultStore) *Worker {
    return &Worker{
        id:            id,
        queue:         queue,
        executor:      executor,
        resultStore:   resultStore,
        stopCh:        make(chan struct{}),
        maxConcurrent: 10,
    }
}

// Start begins processing tasks
func (w *Worker) Start(ctx context.Context) {
    semaphore := make(chan struct{}, w.maxConcurrent)

    for {
        select {
        case <-ctx.Done():
            return
        case <-w.stopCh:
            return
        default:
        }

        // Poll for task
        execution, err := w.queue.Pop(ctx, 5*time.Second)
        if err != nil {
            if err != context.DeadlineExceeded {
                log.Printf("Failed to pop from queue: %v", err)
            }
            continue
        }

        if execution == nil {
            continue
        }

        // Acquire semaphore
        semaphore <- struct{}{}

        go func(exec *domain.Execution) {
            defer func() { <-semaphore }()
            w.executeTask(ctx, exec)
        }(execution)
    }
}

func (w *Worker) executeTask(ctx context.Context, execution *domain.Execution) {
    execution.StartedAt = time.Now()
    execution.Status = domain.ExecutionStatusRunning

    // Get task details
    task, err := w.getTask(ctx, execution.TaskID)
    if err != nil {
        w.handleError(ctx, execution, err)
        return
    }

    // Create execution context with timeout
    execCtx, cancel := context.WithTimeout(ctx, task.Metadata.Timeout)
    defer cancel()

    // Execute
    result, err := w.executor.Execute(execCtx, task)
    execution.CompletedAt = time.Now()

    if err != nil {
        execution.Error = &domain.ExecutionError{
            Code:    "EXECUTION_ERROR",
            Message: err.Error(),
        }
        execution.Status = domain.ExecutionStatusFailure

        // Handle retry
        w.handleRetry(ctx, task, execution)
    } else {
        execution.Status = domain.ExecutionStatusSuccess
        execution.Output = result
    }

    // Store result
    if err := w.resultStore.Save(ctx, execution); err != nil {
        log.Printf("Failed to save execution result: %v", err)
    }
}

func (w *Worker) handleRetry(ctx context.Context, task *domain.Task, execution *domain.Execution) {
    if execution.Attempt >= task.RetryPolicy.MaxRetries {
        // Move to dead letter queue
        w.queue.PushToDLQ(ctx, execution)
        return
    }

    // Calculate retry delay
    delay := calculateBackoff(execution.Attempt, task.RetryPolicy)
    execution.ScheduledAt = time.Now().Add(delay)
    execution.Attempt++

    // Re-queue
    if err := w.queue.Push(ctx, execution); err != nil {
        log.Printf("Failed to re-queue task: %v", err)
    }
}

func calculateBackoff(attempt int, policy domain.RetryPolicy) time.Duration {
    delay := float64(policy.InitialDelay) * math.Pow(policy.Multiplier, float64(attempt))
    if delay > float64(policy.MaxDelay) {
        delay = float64(policy.MaxDelay)
    }

    // Add jitter
    jitter := delay * 0.1 * (rand.Float64()*2 - 1)
    delay += jitter

    return time.Duration(delay)
}

// Executor interface for task execution
type Executor interface {
    Execute(ctx context.Context, task *domain.Task) ([]byte, error)
}

// HTTPExecutor executes HTTP tasks
type HTTPExecutor struct {
    client *http.Client
}

func (e *HTTPExecutor) Execute(ctx context.Context, task *domain.Task) ([]byte, error) {
    var payload struct {
        URL     string            `json:"url"`
        Method  string            `json:"method"`
        Headers map[string]string `json:"headers"`
        Body    string            `json:"body"`
    }

    if err := json.Unmarshal(task.Payload.Data, &payload); err != nil {
        return nil, fmt.Errorf("invalid payload: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, payload.Method, payload.URL, strings.NewReader(payload.Body))
    if err != nil {
        return nil, err
    }

    for k, v := range payload.Headers {
        req.Header.Set(k, v)
    }

    resp, err := e.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }

    return io.ReadAll(resp.Body)
}
```

### 4. API Layer

```go
package api

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

// Server provides HTTP API for task management
type Server struct {
    echo      *echo.Echo
    scheduler *scheduler.Service
}

func NewServer(scheduler *scheduler.Service) *Server {
    e := echo.New()

    s := &Server{
        echo:      e,
        scheduler: scheduler,
    }

    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    // Middleware
    s.echo.Use(middleware.Recover())
    s.echo.Use(middleware.Logger())
    s.echo.Use(middleware.RequestID())

    // Health
    s.echo.GET("/health", s.handleHealth)

    // API v1
    v1 := s.echo.Group("/api/v1")

    // Task routes
    tasks := v1.Group("/tasks")
    tasks.POST("", s.createTask)
    tasks.GET("", s.listTasks)
    tasks.GET("/:id", s.getTask)
    tasks.PUT("/:id", s.updateTask)
    tasks.DELETE("/:id", s.deleteTask)
    tasks.POST("/:id/cancel", s.cancelTask)

    // Execution routes
    tasks.GET("/:id/executions", s.listExecutions)
    tasks.GET("/:id/executions/:executionId", s.getExecution)
}

type CreateTaskRequest struct {
    Name        string             `json:"name" validate:"required"`
    Type        string             `json:"type" validate:"required,oneof=http grpc custom"`
    Schedule    ScheduleRequest    `json:"schedule" validate:"required"`
    Payload     json.RawMessage    `json:"payload"`
    RetryPolicy *RetryPolicyRequest `json:"retry_policy,omitempty"`
    Metadata    *MetadataRequest    `json:"metadata,omitempty"`
}

type ScheduleRequest struct {
    Type     string     `json:"type" validate:"required,oneof=once cron delay"`
    CronExpr string     `json:"cron_expr,omitempty"`
    RunAt    *time.Time `json:"run_at,omitempty"`
    Delay    string     `json:"delay,omitempty"`
}

func (s *Server) createTask(c echo.Context) error {
    var req CreateTaskRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
    }

    if err := c.Validate(req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
    }

    // Convert request to domain
    task, err := s.scheduler.CreateTask(c.Request().Context(), convertRequest(req))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
    }

    return c.JSON(http.StatusCreated, TaskResponse{Task: convertTask(task)})
}

func (s *Server) handleHealth(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{
        "status": "healthy",
    })
}

func (s *Server) Start(address string) error {
    return s.echo.Start(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.echo.Shutdown(ctx)
}
```

---

## Project Deliverables

### Code Requirements

1. **Source Code**
   - Clean, idiomatic Go code
   - Comprehensive error handling
   - Context propagation throughout
   - No race conditions (verified with -race)

2. **Tests**
   - Unit tests (>80% coverage)
   - Integration tests
   - Benchmarks for critical paths
   - Race condition tests

3. **Documentation**
   - README with setup instructions
   - API documentation (OpenAPI)
   - Architecture decision records
   - Deployment guide

4. **Configuration**
   - Environment-based config
   - Docker Compose for local development
   - Kubernetes manifests
   - Helm chart

### Evaluation Criteria

| Category | Weight | Criteria |
|----------|--------|----------|
| Functionality | 30% | All requirements implemented |
| Code Quality | 25% | Idiomatic Go, clean architecture |
| Concurrency | 15% | Proper goroutine management |
| Testing | 15% | Coverage and test quality |
| Observability | 10% | Metrics, logs, traces |
| Documentation | 5% | Clear and comprehensive |

### Submission Checklist

- [ ] All functional requirements implemented
- [ ] Unit test coverage > 80%
- [ ] Integration tests passing
- [ ] Race detector passes
- [ ] Docker build working
- [ ] Kubernetes deployment tested
- [ ] README complete
- [ ] API documentation provided
- [ ] Demo video recorded (optional)

---

## Bonus Challenges

1. **Web UI**: Build a React/Vue frontend
2. **Multi-tenancy**: Isolate tasks by tenant
3. **Workflows**: Support task dependencies
4. **Auto-scaling**: Kubernetes HPA integration
5. **Chaos Engineering**: Failure injection tests

---

## Timeline

| Day | Focus | Deliverable |
|-----|-------|-------------|
| 1 | Domain model, repository | Core data layer |
| 2 | Scheduler service | Task scheduling |
| 3 | Worker, executor | Task execution |
| 4 | API, integration | Full system |
| 5 | Testing, polish | Production ready |

---

*Successful completion of this capstone demonstrates readiness for production Go development.*
