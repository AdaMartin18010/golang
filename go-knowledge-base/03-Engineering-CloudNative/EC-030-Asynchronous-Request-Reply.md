# EC-030: Asynchronous Request-Reply Pattern

## Problem Formalization

### The Synchronous Timeout Problem

Long-running operations cannot be handled through synchronous HTTP requests without causing timeouts, resource exhaustion, and poor user experience. The Asynchronous Request-Reply pattern decouples request submission from result retrieval.

#### Problem Statement

Given:

- Operation O with execution time T_exec
- Client timeout T_timeout where T_exec > T_timeout
- System resources R (connections, threads)

Find a communication pattern that:

```
Minimizes: Client waiting time
Maximizes: System throughput
Subject to:
    - Client receives result eventually
    - System resources bounded
    - Operation can be monitored and cancelled
    - Client can retrieve result at their convenience
```

### Pattern Comparison

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Synchronous vs Asynchronous                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Synchronous (Problematic for Long Operations):                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Client ──Request──►┌─────────────┐                              │   │
│  │                     │   Server    │                              │   │
│  │                     │  (working   │ 5 min ──► Timeout!           │   │
│  │  Client ◄──Response─│   5 min)    │ Connection closed            │   │
│  │                     └─────────────┘                              │   │
│  │                                                                  │   │
│  │  Problems:                                                       │   │
│  │  • Client blocked for 5 minutes                                  │   │
│  │  • Server thread held for 5 minutes                              │   │
│  │  • Timeout forces retry (duplicate work)                         │   │
│  │  • No progress visibility                                        │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Asynchronous (Scalable):                                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Step 1: Submit Request                                          │   │
│  │  Client ──Request──►┌─────────┐                                  │   │
│  │  Client ◄─202/URI──│  API    │                                  │   │
│  │                     └────┬────┘                                  │   │
│  │                          │ Returns immediately                    │   │
│  │                          ▼                                       │   │
│  │                     ┌─────────────┐                              │   │
│  │                     │   Queue     │                              │   │
│  │                     └──────┬──────┘                              │   │
│  │                            │                                     │   │
│  │  Step 2: Poll for Status                                         │   │
│  │  Client ──GET /status──►┌─────────────┐                         │   │
│  │  Client ◄─202 Pending───│   Worker    │                         │   │
│  │  ...                    │  (async)    │                         │   │
│  │  Client ──GET /status──►│  processing │                         │   │
│  │  Client ◄─200 Complete──│             │                         │   │
│  │                         └─────────────┘                         │   │
│  │                                                                  │   │
│  │  Step 3: Retrieve Result                                         │   │
│  │  Client ──GET /result──►┌─────────────┐                         │   │
│  │  Client ◄─200 + Data────│  Storage    │                         │   │
│  │                         └─────────────┘                         │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Asynchronous Request-Reply Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Asynchronous Request-Reply Flow                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  1. Request Submission                                           │   │
│  │                                                                  │   │
│  │  POST /api/v1/jobs                                               │   │
│  │  {                                                               │   │
│  │    "operation": "generate_report",                               │   │
│  │    "parameters": {...},                                          │   │
│  │    "callback_url": "https://client.com/webhook"                  │   │
│  │  }                                                               │   │
│  │                                                                  │   │
│  │  Response: 202 Accepted                                          │   │
│  │  {                                                               │   │
│  │    "job_id": "550e8400-e29b-41d4-a716-446655440000",            │   │
│  │    "status": "pending",                                          │   │
│  │    "status_url": "/api/v1/jobs/550e8400-e29b-41d4-a716-446655440000", │   │
│  │    "estimated_completion": "2024-01-15T10:30:00Z",               │   │
│  │    "retry_after": 30                                             │   │
│  │  }                                                               │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  2. Request Processing                                           │   │
│  │                                                                  │   │
│  │  ┌──────────┐    ┌──────────┐    ┌──────────┐                  │   │
│  │  │   API    │───►│  Queue   │───►│  Worker  │                  │   │
│  │  │  Server  │    │          │    │  Pool    │                  │   │
│  │  └──────────┘    └──────────┘    └────┬─────┘                  │   │
│  │                                        │                        │   │
│  │                                        ▼                        │   │
│  │                              ┌──────────────────┐               │   │
│  │                              │ Status: running  │               │   │
│  │                              │ Progress: 45%    │               │   │
│  │                              │ (Redis/Database) │               │   │
│  │                              └──────────────────┘               │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  3. Status Polling / Callback                                    │   │
│  │                                                                  │   │
│  │  Option A: Polling                                               │   │
│  │  GET /api/v1/jobs/550e8400-e29b-41d4-a716-446655440000           │   │
│  │                                                                  │   │
│  │  Response (in progress):                                         │   │
│  │  {                                                               │   │
│  │    "job_id": "550e8400-...",                                     │   │
│  │    "status": "running",                                          │   │
│  │    "progress": 45,                                               │   │
│  │    "started_at": "2024-01-15T10:00:00Z",                         │   │
│  │    "retry_after": 30                                             │   │
│  │  }                                                               │   │
│  │                                                                  │   │
│  │  Response (complete):                                            │   │
│  │  {                                                               │   │
│  │    "job_id": "550e8400-...",                                     │   │
│  │    "status": "completed",                                        │   │
│  │    "completed_at": "2024-01-15T10:25:00Z",                       │   │
│  │    "result_url": "/api/v1/jobs/550e8400-.../result",            │   │
│  │    "expires_at": "2024-01-16T10:25:00Z"                          │   │
│  │  }                                                               │   │
│  │                                                                  │   │
│  │  Option B: Webhook Callback                                      │   │
│  │  POST https://client.com/webhook                                 │   │
│  │  {                                                               │   │
│  │    "job_id": "550e8400-...",                                     │   │
│  │    "status": "completed",                                        │   │
│  │    "result_url": "..."                                           │   │
│  │  }                                                               │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              │                                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  4. Result Retrieval                                             │   │
│  │                                                                  │   │
│  │  GET /api/v1/jobs/550e8400-e29b-41d4-a716-446655440000/result    │   │
│  │                                                                  │   │
│  │  Response: 200 OK                                                │   │
│  │  {                                                               │   │
│  │    "report_url": "https://storage.../report.pdf",                │   │
│  │    "metadata": {...},                                            │   │
│  │    "expires_at": "2024-01-16T10:25:00Z"                          │   │
│  │  }                                                               │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Job Manager

```go
// pkg/asyncjob/manager.go
package asyncjob

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/google/uuid"
    "github.com/go-redis/redis/v8"
)

// JobStatus represents the current state of a job
type JobStatus string

const (
    JobStatusPending    JobStatus = "pending"
    JobStatusRunning    JobStatus = "running"
    JobStatusCompleted  JobStatus = "completed"
    JobStatusFailed     JobStatus = "failed"
    JobStatusCancelled  JobStatus = "cancelled"
)

// Job represents an asynchronous job
type Job struct {
    ID              string                 `json:"id"`
    Type            string                 `json:"type"`
    Status          JobStatus              `json:"status"`
    Parameters      map[string]interface{} `json:"parameters"`
    Result          interface{}            `json:"result,omitempty"`
    Error           string                 `json:"error,omitempty"`
    Progress        int                    `json:"progress"`
    CreatedAt       time.Time              `json:"created_at"`
    StartedAt       *time.Time             `json:"started_at,omitempty"`
    CompletedAt     *time.Time             `json:"completed_at,omitempty"`
    ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
    CallbackURL     string                 `json:"callback_url,omitempty"`
    CorrelationID   string                 `json:"correlation_id,omitempty"`
}

// JobStore persists job state
type JobStore interface {
    Create(ctx context.Context, job *Job) error
    Get(ctx context.Context, jobID string) (*Job, error)
    Update(ctx context.Context, job *Job) error
    Delete(ctx context.Context, jobID string) error
    List(ctx context.Context, filter JobFilter) ([]*Job, error)
}

// JobExecutor processes jobs
type JobExecutor interface {
    Execute(ctx context.Context, job *Job) (interface{}, error)
    GetType() string
}

// Manager orchestrates asynchronous jobs
type Manager struct {
    store     JobStore
    queue     Queue
    executors map[string]JobExecutor
    callback  CallbackClient

    ctx       context.Context
    cancel    context.CancelFunc
    wg        sync.WaitGroup
}

type Queue interface {
    Enqueue(ctx context.Context, job *Job) error
    Dequeue(ctx context.Context) (*Job, error)
    Ack(ctx context.Context, job *Job) error
    Nack(ctx context.Context, job *Job, requeue bool) error
}

type CallbackClient interface {
    Send(ctx context.Context, url string, payload interface{}) error
}

func NewManager(store JobStore, queue Queue) *Manager {
    ctx, cancel := context.WithCancel(context.Background())

    return &Manager{
        store:     store,
        queue:     queue,
        executors: make(map[string]JobExecutor),
        ctx:       ctx,
        cancel:    cancel,
    }
}

func (m *Manager) RegisterExecutor(executor JobExecutor) {
    m.executors[executor.GetType()] = executor
}

func (m *Manager) SetCallbackClient(client CallbackClient) {
    m.callback = client
}

// Submit creates and queues a new job
func (m *Manager) Submit(ctx context.Context, jobType string, parameters map[string]interface{}, callbackURL string) (*Job, error) {
    executor, ok := m.executors[jobType]
    if !ok {
        return nil, fmt.Errorf("unknown job type: %s", jobType)
    }

    job := &Job{
        ID:          uuid.New().String(),
        Type:        jobType,
        Status:      JobStatusPending,
        Parameters:  parameters,
        CreatedAt:   time.Now(),
        CallbackURL: callbackURL,
    }

    // Persist job
    if err := m.store.Create(ctx, job); err != nil {
        return nil, fmt.Errorf("creating job: %w", err)
    }

    // Queue for execution
    if err := m.queue.Enqueue(ctx, job); err != nil {
        return nil, fmt.Errorf("queuing job: %w", err)
    }

    return job, nil
}

// Get retrieves job status
func (m *Manager) Get(ctx context.Context, jobID string) (*Job, error) {
    return m.store.Get(ctx, jobID)
}

// Cancel attempts to cancel a running job
func (m *Manager) Cancel(ctx context.Context, jobID string) error {
    job, err := m.store.Get(ctx, jobID)
    if err != nil {
        return err
    }

    if job.Status != JobStatusPending && job.Status != JobStatusRunning {
        return fmt.Errorf("cannot cancel job in status: %s", job.Status)
    }

    job.Status = JobStatusCancelled
    job.CompletedAt = timePtr(time.Now())

    return m.store.Update(ctx, job)
}

// Start begins processing jobs
func (m *Manager) Start(workers int) {
    for i := 0; i < workers; i++ {
        m.wg.Add(1)
        go m.worker(i)
    }
}

func (m *Manager) Stop() {
    m.cancel()
    m.wg.Wait()
}

func (m *Manager) worker(id int) {
    defer m.wg.Done()

    for {
        select {
        case <-m.ctx.Done():
            return
        default:
        }

        // Dequeue job
        job, err := m.queue.Dequeue(m.ctx)
        if err != nil {
            continue
        }

        // Process job
        m.processJob(job)
    }
}

func (m *Manager) processJob(job *Job) {
    executor, ok := m.executors[job.Type]
    if !ok {
        m.failJob(job, fmt.Sprintf("no executor for type: %s", job.Type))
        return
    }

    // Update status to running
    job.Status = JobStatusRunning
    now := time.Now()
    job.StartedAt = &now

    if err := m.store.Update(m.ctx, job); err != nil {
        m.queue.Nack(m.ctx, job, true)
        return
    }

    // Create cancellable context
    ctx, cancel := context.WithCancel(m.ctx)
    defer cancel()

    // Execute with progress updates
    result, err := executor.Execute(ctx, job)

    if err != nil {
        m.failJob(job, err.Error())
        return
    }

    // Complete successfully
    m.completeJob(job, result)
}

func (m *Manager) completeJob(job *Job, result interface{}) {
    job.Status = JobStatusCompleted
    job.Result = result
    job.Progress = 100
    now := time.Now()
    job.CompletedAt = &now

    // Set expiration (24 hours)
    expires := now.Add(24 * time.Hour)
    job.ExpiresAt = &expires

    m.store.Update(m.ctx, job)
    m.queue.Ack(m.ctx, job)

    // Send callback if configured
    if job.CallbackURL != "" && m.callback != nil {
        m.callback.Send(m.ctx, job.CallbackURL, map[string]interface{}{
            "job_id":     job.ID,
            "status":     job.Status,
            "result":     job.Result,
            "completed_at": job.CompletedAt,
        })
    }
}

func (m *Manager) failJob(job *Job, errorMsg string) {
    job.Status = JobStatusFailed
    job.Error = errorMsg
    now := time.Now()
    job.CompletedAt = &now

    m.store.Update(m.ctx, job)
    m.queue.Ack(m.ctx, job)

    // Send callback if configured
    if job.CallbackURL != "" && m.callback != nil {
        m.callback.Send(m.ctx, job.CallbackURL, map[string]interface{}{
            "job_id": job.ID,
            "status": job.Status,
            "error":  job.Error,
        })
    }
}

func timePtr(t time.Time) *time.Time {
    return &t
}
```

### HTTP Handler

```go
// internal/api/job_handler.go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/company/project/pkg/asyncjob"
)

type JobHandler struct {
    manager *asyncjob.Manager
}

func NewJobHandler(manager *asyncjob.Manager) *JobHandler {
    return &JobHandler{manager: manager}
}

func (h *JobHandler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Post("/", h.CreateJob)
    r.Get("/{jobID}", h.GetJob)
    r.Delete("/{jobID}", h.CancelJob)
    r.Get("/{jobID}/result", h.GetResult)
    return r
}

// CreateJob handles job submission
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Type        string                 `json:"type"`
        Parameters  map[string]interface{} `json:"parameters"`
        CallbackURL string                 `json:"callback_url,omitempty"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    job, err := h.manager.Submit(r.Context(), req.Type, req.Parameters, req.CallbackURL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Return 202 Accepted
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Location", "/api/v1/jobs/"+job.ID)
    w.WriteHeader(http.StatusAccepted)

    response := struct {
        JobID      string `json:"job_id"`
        Status     string `json:"status"`
        StatusURL  string `json:"status_url"`
        RetryAfter int    `json:"retry_after"`
    }{
        JobID:      job.ID,
        Status:     string(job.Status),
        StatusURL:  "/api/v1/jobs/" + job.ID,
        RetryAfter: 5,
    }

    json.NewEncoder(w).Encode(response)
}

// GetJob handles status polling
func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
    jobID := chi.URLParam(r, "jobID")

    job, err := h.manager.Get(r.Context(), jobID)
    if err != nil {
        http.Error(w, "Job not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    // Set appropriate status code
    switch job.Status {
    case asyncjob.JobStatusPending, asyncjob.JobStatusRunning:
        w.WriteHeader(http.StatusAccepted)
        w.Header().Set("Retry-After", "5")
    case asyncjob.JobStatusCompleted, asyncjob.JobStatusFailed, asyncjob.JobStatusCancelled:
        w.WriteHeader(http.StatusOK)
    }

    json.NewEncoder(w).Encode(job)
}

// CancelJob handles job cancellation
func (h *JobHandler) CancelJob(w http.ResponseWriter, r *http.Request) {
    jobID := chi.URLParam(r, "jobID")

    if err := h.manager.Cancel(r.Context(), jobID); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// GetResult handles result retrieval
func (h *JobHandler) GetResult(w http.ResponseWriter, r *http.Request) {
    jobID := chi.URLParam(r, "jobID")

    job, err := h.manager.Get(r.Context(), jobID)
    if err != nil {
        http.Error(w, "Job not found", http.StatusNotFound)
        return
    }

    if job.Status != asyncjob.JobStatusCompleted {
        http.Error(w, "Job not complete", http.StatusConflict)
        return
    }

    // Check if result expired
    if job.ExpiresAt != nil && time.Now().After(*job.ExpiresAt) {
        http.Error(w, "Result expired", http.StatusGone)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    result := struct {
        JobID       string      `json:"job_id"`
        Result      interface{} `json:"result"`
        ExpiresAt   *time.Time  `json:"expires_at"`
    }{
        JobID:     job.ID,
        Result:    job.Result,
        ExpiresAt: job.ExpiresAt,
    }

    json.NewEncoder(w).Encode(result)
}
```

## Trade-off Analysis

### Polling vs Webhooks

| Aspect | Polling | Webhooks | Notes |
|--------|---------|----------|-------|
| **Latency** | Polling interval dependent | Near real-time | Webhooks faster |
| **Reliability** | Client controls retry | Must handle failures | Polling more reliable |
| **Firewall** | No issues | Client must accept inbound | Polling simpler |
| **Server Load** | Scales with clients | Constant per completion | Polling can overload |
| **Implementation** | Simple on server | Complex (retry, idempotency) | Polling easier |
| **Best For** | Few clients, simple needs | Many clients, low latency | Hybrid often works |

### Status Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Job Status State Machine                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐     ┌──────────┐     ┌───────────┐     ┌──────────┐       │
│  │ Pending │────►│ Running  │────►│Completed  │     │ Expired  │       │
│  └────┬────┘     └────┬─────┘     └───────────┘     └──────────┘       │
│       │               │                                               │
│       │               └──────┐                                        │
│       │                      ▼                                        │
│       │               ┌──────────┐     ┌──────────┐                   │
│       │               │  Failed  │────►│ Retrying │ (auto-retry)      │
│       │               └──────────┘     └──────────┘                   │
│       │                                                               │
│       └──────────────►┌──────────┐                                    │
│         (cancel)      │Cancelled │                                    │
│                       └──────────┘                                    │
│                                                                         │
│  Transitions:                                                           │
│  • Pending → Running: Worker picks up job                              │
│  • Pending → Cancelled: Client cancels before start                    │
│  • Running → Completed: Success                                         │
│  • Running → Failed: Error, may retry                                  │
│  • Running → Cancelled: Client cancels                                 │
│  • Completed → Expired: TTL reached                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Asynchronous Job Testing

```go
// test/asyncjob/manager_test.go
package asyncjob

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestJobSubmission(t *testing.T) {
    store := new(MockJobStore)
    queue := new(MockQueue)

    manager := NewManager(store, queue)

    executor := new(MockJobExecutor)
    executor.On("GetType").Return("test-job")
    manager.RegisterExecutor(executor)

    // Expectations
    store.On("Create", mock.Anything, mock.AnythingOfType("*asyncjob.Job")).Return(nil)
    queue.On("Enqueue", mock.Anything, mock.AnythingOfType("*asyncjob.Job")).Return(nil)

    // Submit job
    job, err := manager.Submit(context.Background(), "test-job", map[string]interface{}{
        "param1": "value1",
    }, "")

    assert.NoError(t, err)
    assert.NotEmpty(t, job.ID)
    assert.Equal(t, JobStatusPending, job.Status)

    store.AssertExpectations(t)
    queue.AssertExpectations(t)
}

func TestJobExecution(t *testing.T) {
    store := new(MockJobStore)
    queue := new(MockQueue)

    manager := NewManager(store, queue)

    // Setup executor
    executor := new(MockJobExecutor)
    executor.On("GetType").Return("test-job")
    executor.On("Execute", mock.Anything, mock.Anything).Return(
        map[string]string{"result": "success"}, nil)
    manager.RegisterExecutor(executor)

    // Setup job
    job := &Job{
        ID:     "test-id",
        Type:   "test-job",
        Status: JobStatusPending,
    }

    store.On("Update", mock.Anything, mock.Anything).Return(nil)
    queue.On("Ack", mock.Anything, mock.Anything).Return(nil)

    // Process job
    manager.Start(1)

    // Simulate job completion
    // (Would need actual queue implementation for full test)

    manager.Stop()
}

func TestJobCancellation(t *testing.T) {
    store := new(MockJobStore)
    queue := new(MockQueue)

    manager := NewManager(store, queue)

    job := &Job{
        ID:     "test-id",
        Status: JobStatusPending,
    }

    store.On("Get", mock.Anything, "test-id").Return(job, nil)
    store.On("Update", mock.Anything, mock.MatchedBy(func(j *Job) bool {
        return j.Status == JobStatusCancelled
    })).Return(nil)

    err := manager.Cancel(context.Background(), "test-id")

    assert.NoError(t, err)
    store.AssertExpectations(t)
}

func TestJobExpiration(t *testing.T) {
    store := new(MockJobStore)
    queue := new(MockQueue)

    manager := NewManager(store, queue)

    // Completed but expired job
    completedAt := time.Now().Add(-48 * time.Hour)
    expiresAt := time.Now().Add(-24 * time.Hour)

    job := &Job{
        ID:          "expired-id",
        Status:      JobStatusCompleted,
        CompletedAt: &completedAt,
        ExpiresAt:   &expiresAt,
        Result:      map[string]string{"data": "old"},
    }

    store.On("Get", mock.Anything, "expired-id").Return(job, nil)

    // Simulate result retrieval - should return expired
    // (Test would be in HTTP handler)
}
```

## Summary

The Asynchronous Request-Reply Pattern provides:

1. **Scalability**: Handle long operations without blocking
2. **Resilience**: Survive client and server restarts
3. **Observability**: Track operation progress
4. **Flexibility**: Support polling and callbacks
5. **Resource Efficiency**: Bounded thread/connection usage

Key considerations:

- Set appropriate TTLs for results
- Implement idempotency for retries
- Handle webhook delivery failures
- Monitor queue depth and processing lag
- Provide cancellation capability
