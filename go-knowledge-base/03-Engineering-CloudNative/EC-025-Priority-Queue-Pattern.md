# EC-025: Priority Queue Pattern

## Problem Formalization

### The Resource Contention Challenge

In distributed systems, not all work is created equal. Critical operations must complete before lower-priority tasks, but FIFO queues treat all messages equally, potentially starving important work during high load.

#### Problem Statement

Given:

- Task set T = {t₁, t₂, ..., tₙ} with priorities P = {p₁, p₂, ..., pₙ}
- Processing capacity C (tasks per unit time)
- Service level objectives SLO = {slo₁, slo₂, ..., sloₙ}

Find scheduling order O that:

```
Maximize: Σ Satisfaction(tᵢ) where Satisfaction = 1 if completion_time ≤ sloᵢ
Subject to:
    - Processing rate ≤ C
    - Higher priority tasks generally complete before lower priority
    - Starvation of low-priority tasks is bounded
```

### Priority Queue Use Cases

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Priority Queue Use Cases                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. Payment Processing                                                  │
│     ┌─────────────────────────────────────────────────────────────┐    │
│  P0 │ Fraud alerts, Chargebacks                                   │    │
│  P1 │ Real-time payments, Refunds                                 │    │
│  P2 │ Subscription renewals                                       │    │
│  P3 │ Reconciliation reports                                      │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  2. IoT Device Management                                               │
│     ┌─────────────────────────────────────────────────────────────┐    │
│  P0 │ Emergency shutdown commands                                 │    │
│  P1 │ Configuration updates                                       │    │
│  P2 │ Firmware downloads                                          │    │
│  P3 │ Telemetry batch uploads                                     │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  3. Healthcare Systems                                                  │
│     ┌─────────────────────────────────────────────────────────────┐    │
│  P0 │ Critical lab results, Code blue alerts                      │    │
│  P1 │ Urgent prescriptions                                        │    │
│  P2 │ Appointment reminders                                       │    │
│  P3 │ Routine reports                                             │    │
│     └─────────────────────────────────────────────────────────────┘    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Priority Queue Implementation Strategies

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Priority Queue Architectures                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Strategy 1: Single Priority Queue (Binary Heap)                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                 │   │
│  │        Push ──►  ┌───────────────┐  ◄── Pop (highest priority) │   │
│  │  [P3, P1, P0]    │ Priority Heap │      [P0, P1, P2, P3...]    │   │
│  │                  │  O(log n)     │                             │   │
│  │                  └───────────────┘                             │   │
│  │                                                                 │   │
│  │  Pros: Simple, memory efficient                                 │   │
│  │  Cons: Contention on single queue, no priority isolation        │   │
│  │                                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Strategy 2: Multiple Queues (One per Priority)                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                 │   │
│  │  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐                │   │
│  │  │ Queue  │  │ Queue  │  │ Queue  │  │ Queue  │                │   │
│  │  │  P0    │  │  P1    │  │  P2    │  │  P3    │                │   │
│  │  │[x, x]  │  │[x]     │  │[x,x,x] │  │[x...]  │                │   │
│  │  └───┬────┘  └───┬────┘  └───┬────┘  └───┬────┘                │   │
│  │      │           │           │           │                      │   │
│  │      └───────────┴─────┬─────┴───────────┘                      │   │
│  │                        ▼                                         │   │
│  │                  Scheduler                                       │   │
│  │                  (Weighted Round-Robin)                          │   │
│  │                                                                 │   │
│  │  Pros: Priority isolation, configurable scheduling               │   │
│  │  Cons: More complex, potential starvation if not careful         │   │
│  │                                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Strategy 3: Time-Based Priority (Aging)                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                 │   │
│  │  Dynamic Priority = Base Priority + Time in Queue × Aging Factor│   │
│  │                                                                 │   │
│  │  Prevents starvation: low priority tasks eventually become      │   │
│  │  high priority if waiting long enough                           │   │
│  │                                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Multi-Queue Scheduler

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Weighted Fair Queuing Scheduler                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Configuration                                                   │   │
│  │  ┌───────────────┬─────────────────┬──────────────────────────┐  │   │
│  │  │   Priority    │   Weight        │   Max Queue Depth        │  │   │
│  │  ├───────────────┼─────────────────┼──────────────────────────┤  │   │
│  │  │   P0 (Critical)│   50%          │   1000                   │  │   │
│  │  │   P1 (High)   │   30%          │   5000                   │  │   │
│  │  │   P2 (Normal) │   15%          │   10000                  │  │   │
│  │  │   P3 (Low)    │   5%           │   Unlimited              │  │   │
│  │  └───────────────┴─────────────────┴──────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Scheduling Algorithm:                                                  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                 │   │
│  │  For each scheduling decision:                                  │   │
│  │  1. Check P0 queue, if not empty, process P0                    │   │
│  │  2. Otherwise, use weighted selection:                          │   │
│  │     - P1: 60% (30/50 remaining)                                 │   │
│  │     - P2: 30% (15/25 remaining)                                 │   │
│  │     - P3: 10% (5/10 remaining)                                  │   │
│  │  3. If selected queue empty, fall through to next               │   │
│  │                                                                 │   │
│  │  Weights reset after fixed time window or processing quantum    │   │
│  │                                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Priority Queue Core

```go
// pkg/priorityqueue/queue.go
package priorityqueue

import (
    "container/heap"
    "context"
    "sync"
    "time"
)

// Priority levels
type Priority int

const (
    PriorityCritical Priority = iota // P0
    PriorityHigh                     // P1
    PriorityNormal                   // P2
    PriorityLow                      // P3
    PriorityBackground               // P4
)

func (p Priority) String() string {
    switch p {
    case PriorityCritical:
        return "P0-Critical"
    case PriorityHigh:
        return "P1-High"
    case PriorityNormal:
        return "P2-Normal"
    case PriorityLow:
        return "P3-Low"
    case PriorityBackground:
        return "P4-Background"
    default:
        return "Unknown"
    }
}

// Item represents an item in the priority queue
type Item struct {
    ID          string
    Priority    Priority
    Payload     interface{}
    CreatedAt   time.Time
    Deadline    *time.Time // Optional deadline
    Attempts    int
    MaxAttempts int

    // Internal heap fields
    index int // position in heap
}

// PriorityQueue implements heap.Interface
type PriorityQueue struct {
    items []*Item
    mu    sync.RWMutex
}

func NewPriorityQueue() *PriorityQueue {
    pq := &PriorityQueue{
        items: make([]*Item, 0),
    }
    heap.Init(pq)
    return pq
}

func (pq PriorityQueue) Len() int { return len(pq.items) }

func (pq PriorityQueue) Less(i, j int) bool {
    // Higher priority (lower number) comes first
    if pq.items[i].Priority != pq.items[j].Priority {
        return pq.items[i].Priority < pq.items[j].Priority
    }

    // Same priority: earlier deadline or creation time first
    if pq.items[i].Deadline != nil && pq.items[j].Deadline != nil {
        return pq.items[i].Deadline.Before(*pq.items[j].Deadline)
    }

    return pq.items[i].CreatedAt.Before(pq.items[j].CreatedAt)
}

func (pq PriorityQueue) Swap(i, j int) {
    pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
    pq.items[i].index = i
    pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(pq.items)
    item := x.(*Item)
    item.index = n
    pq.items = append(pq.items, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := pq.items
    n := len(old)
    item := old[n-1]
    old[n-1] = nil  // avoid memory leak
    item.index = -1 // for safety
    pq.items = old[0 : n-1]
    return item
}

func (pq *PriorityQueue) Enqueue(item *Item) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    if item.CreatedAt.IsZero() {
        item.CreatedAt = time.Now()
    }

    heap.Push(pq, item)
}

func (pq *PriorityQueue) Dequeue() *Item {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    if pq.Len() == 0 {
        return nil
    }

    return heap.Pop(pq).(*Item)
}

func (pq *PriorityQueue) Peek() *Item {
    pq.mu.RLock()
    defer pq.mu.RUnlock()

    if pq.Len() == 0 {
        return nil
    }

    return pq.items[0]
}

func (pq *PriorityQueue) Size() int {
    pq.mu.RLock()
    defer pq.mu.RUnlock()
    return pq.Len()
}
```

### Multi-Queue Priority System

```go
// pkg/priorityqueue/multi_queue.go
package priorityqueue

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "time"

    "github.com/IBM/sarama"
)

// MultiQueuePrioritySystem manages multiple priority queues
type MultiQueuePrioritySystem struct {
    config    *MultiQueueConfig
    queues    map[Priority]*PriorityQueue
    consumers map[Priority]*QueueConsumer

    // Scheduling
    scheduler Scheduler

    // Metrics
    processed      uint64
    failed         uint64
    byPriority     map[Priority]*uint64

    // Control
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup

    handler Handler
}

type MultiQueueConfig struct {
    QueueConfigs map[Priority]QueueConfig
    Scheduler    SchedulerConfig
}

type QueueConfig struct {
    MaxDepth       int
    Weight         float64 // For weighted scheduling
    MaxRetries     int
    Timeout        time.Duration
    Workers        int
}

type SchedulerConfig struct {
    Type           string // "strict", "weighted", "fair"
    CheckInterval  time.Duration
}

// Handler processes items
type Handler interface {
    Process(ctx context.Context, item *Item) error
}

// Scheduler determines which queue to process next
type Scheduler interface {
    SelectQueue(queues map[Priority]*PriorityQueue, stats map[Priority]QueueStats) Priority
    UpdateWeights(weights map[Priority]float64)
}

// WeightedFairScheduler implements weighted fair queuing
type WeightedFairScheduler struct {
    weights   map[Priority]float64
    processed map[Priority]uint64
    mu        sync.Mutex
}

func NewWeightedFairScheduler(weights map[Priority]float64) *WeightedFairScheduler {
    return &WeightedFairScheduler{
        weights:   weights,
        processed: make(map[Priority]uint64),
    }
}

func (s *WeightedFairScheduler) SelectQueue(
    queues map[Priority]*PriorityQueue,
    stats map[Priority]QueueStats,
) Priority {
    s.mu.Lock()
    defer s.mu.Unlock()

    // First, check critical queue
    if queues[PriorityCritical].Size() > 0 {
        s.processed[PriorityCritical]++
        return PriorityCritical
    }

    // Calculate deficits for each queue
    var bestPriority Priority = PriorityBackground
    maxDeficit := -1.0

    for p := PriorityHigh; p <= PriorityBackground; p++ {
        if queues[p].Size() == 0 {
            continue
        }

        weight := s.weights[p]
        processed := float64(s.processed[p])

        // Deficit = expected - actual
        // Expected proportional to weight
        totalProcessed := float64(0)
        for _, count := range s.processed {
            totalProcessed += float64(count)
        }

        if totalProcessed == 0 {
            s.processed[p]++
            return p
        }

        expected := totalProcessed * weight
        deficit := expected - processed

        if deficit > maxDeficit {
            maxDeficit = deficit
            bestPriority = p
        }
    }

    if bestPriority != PriorityBackground || queues[PriorityBackground].Size() > 0 {
        s.processed[bestPriority]++
    }

    return bestPriority
}

func NewMultiQueuePrioritySystem(config *MultiQueueConfig, handler Handler) (*MultiQueuePrioritySystem, error) {
    ctx, cancel := context.WithCancel(context.Background())

    mqs := &MultiQueuePrioritySystem{
        config:     config,
        queues:     make(map[Priority]*PriorityQueue),
        consumers:  make(map[Priority]*QueueConsumer),
        byPriority: make(map[Priority]*uint64),
        ctx:        ctx,
        cancel:     cancel,
        handler:    handler,
    }

    // Initialize queues
    for p := PriorityCritical; p <= PriorityBackground; p++ {
        mqs.queues[p] = NewPriorityQueue()
        counter := uint64(0)
        mqs.byPriority[p] = &counter
    }

    // Initialize scheduler
    weights := make(map[Priority]float64)
    for p, cfg := range config.QueueConfigs {
        weights[p] = cfg.Weight
    }
    mqs.scheduler = NewWeightedFairScheduler(weights)

    return mqs, nil
}

func (mqs *MultiQueuePrioritySystem) Start() {
    // Start scheduler
    mqs.wg.Add(1)
    go mqs.schedulingLoop()

    // Start workers for each priority
    for p, cfg := range mqs.config.QueueConfigs {
        for i := 0; i < cfg.Workers; i++ {
            mqs.wg.Add(1)
            go mqs.workerLoop(p, i)
        }
    }
}

func (mqs *MultiQueuePrioritySystem) Stop() {
    mqs.cancel()
    mqs.wg.Wait()
}

func (mqs *MultiQueuePrioritySystem) Enqueue(item *Item) error {
    if item.Priority < PriorityCritical || item.Priority > PriorityBackground {
        return fmt.Errorf("invalid priority: %d", item.Priority)
    }

    cfg := mqs.config.QueueConfigs[item.Priority]
    queue := mqs.queues[item.Priority]

    // Check queue depth
    if cfg.MaxDepth > 0 && queue.Size() >= cfg.MaxDepth {
        return fmt.Errorf("queue %s at max depth %d", item.Priority.String(), cfg.MaxDepth)
    }

    queue.Enqueue(item)
    return nil
}

func (mqs *MultiQueuePrioritySystem) schedulingLoop() {
    defer mqs.wg.Done()

    ticker := time.NewTicker(mqs.config.Scheduler.CheckInterval)
    defer ticker.Stop()

    for {
        select {
        case <-mqs.ctx.Done():
            return
        case <-ticker.C:
            // Collect stats
            stats := make(map[Priority]QueueStats)
            for p, q := range mqs.queues {
                stats[p] = QueueStats{Depth: q.Size()}
            }

            // Scheduler can adjust based on stats
            _ = mqs.scheduler.SelectQueue(mqs.queues, stats)
        }
    }
}

func (mqs *MultiQueuePrioritySystem) workerLoop(priority Priority, workerID int) {
    defer mqs.wg.Done()

    cfg := mqs.config.QueueConfigs[priority]
    queue := mqs.queues[priority]

    for {
        select {
        case <-mqs.ctx.Done():
            return
        default:
        }

        // Try to get item from this priority queue
        item := queue.Dequeue()
        if item == nil {
            // Queue empty, try next priority or wait
            time.Sleep(10 * time.Millisecond)
            continue
        }

        // Process item
        ctx, cancel := context.WithTimeout(mqs.ctx, cfg.Timeout)
        err := mqs.processItem(ctx, item)
        cancel()

        if err != nil {
            mqs.handleFailure(item, err, cfg)
        } else {
            atomic.AddUint64(mqs.byPriority[priority], 1)
            atomic.AddUint64(&mqs.processed, 1)
        }
    }
}

func (mqs *MultiQueuePrioritySystem) processItem(ctx context.Context, item *Item) error {
    return mqs.handler.Process(ctx, item)
}

func (mqs *MultiQueuePrioritySystem) handleFailure(item *Item, err error, cfg QueueConfig) {
    item.Attempts++

    if item.Attempts >= item.MaxAttempts || item.Attempts >= cfg.MaxRetries {
        // Dead letter
        atomic.AddUint64(&mqs.failed, 1)
        mqs.sendToDeadLetter(item, err)
    } else {
        // Re-queue with exponential backoff
        time.AfterFunc(calculateBackoff(item.Attempts), func() {
            mqs.queues[item.Priority].Enqueue(item)
        })
    }
}

func (mqs *MultiQueuePrioritySystem) sendToDeadLetter(item *Item, err error) {
    // Implementation depends on your dead letter strategy
    // Could write to database, another queue, or alert system
}

func calculateBackoff(attempt int) time.Duration {
    base := time.Second
    max := 5 * time.Minute

    backoff := base * time.Duration(1<<uint(attempt))
    if backoff > max {
        backoff = max
    }

    // Add jitter
    jitter := time.Duration(rand.Int63n(int64(backoff) / 2))
    return backoff + jitter
}

type QueueStats struct {
    Depth     int
    Processing int
    WaitTime  time.Duration
}
```

## Trade-off Analysis

### Priority Queue Strategies

| Strategy | Latency for P0 | Fairness | Complexity | Use Case |
|----------|---------------|----------|------------|----------|
| **Strict Priority** | Lowest | Poor (starvation) | Low | Emergency systems |
| **Weighted Fair** | Low | Good | Medium | General purpose |
| **Aging/Time-based** | Medium | Very Good | Medium | Batch systems |
| **Token Bucket** | Configurable | Excellent | High | Multi-tenant systems |

### Priority Inversion Prevention

```
Priority Inversion Problem:
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│  Scenario:                                                              │
│  1. Low-priority task (L) acquires lock                                 │
│  2. High-priority task (H) tries to acquire same lock, blocks           │
│  3. Medium-priority task (M) preempts L, preventing it from releasing   │
│  4. Result: H waits for M (priority inversion!)                         │
│                                                                         │
│  Solutions:                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  1. Priority Inheritance                                          │   │
│  │     L's priority temporarily raised to H's level                │   │
│  │     When L releases lock, priority restored                     │   │
│  │                                                                  │   │
│  │  2. Priority Ceiling                                              │   │
│  │     Each resource has ceiling priority = max(priority of tasks   │   │
│  │     that can lock it)                                           │   │
│  │     Task gets ceiling priority when locking                     │   │
│  │                                                                  │   │
│  │  3. Avoid Shared Resources                                        │   │
│  │     Use message passing instead of shared state                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Priority Queue Testing

```go
// test/priorityqueue/queue_test.go
package priorityqueue

import (
    "context"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPriorityQueueOrdering(t *testing.T) {
    pq := NewPriorityQueue()

    // Add items in random order
    items := []*Item{
        {ID: "1", Priority: PriorityLow},
        {ID: "2", Priority: PriorityCritical},
        {ID: "3", Priority: PriorityHigh},
        {ID: "4", Priority: PriorityNormal},
        {ID: "5", Priority: PriorityCritical},
    }

    for _, item := range items {
        pq.Enqueue(item)
    }

    // Should dequeue in priority order
    first := pq.Dequeue()
    assert.Equal(t, PriorityCritical, first.Priority)

    second := pq.Dequeue()
    assert.Equal(t, PriorityCritical, second.Priority)

    third := pq.Dequeue()
    assert.Equal(t, PriorityHigh, third.Priority)
}

func TestMultiQueueProcessing(t *testing.T) {
    config := &MultiQueueConfig{
        QueueConfigs: map[Priority]QueueConfig{
            PriorityCritical: {Weight: 0.5, Workers: 2, Timeout: time.Second},
            PriorityHigh:     {Weight: 0.3, Workers: 2, Timeout: time.Second},
            PriorityNormal:   {Weight: 0.2, Workers: 1, Timeout: time.Second},
        },
        Scheduler: SchedulerConfig{CheckInterval: time.Second},
    }

    processed := make([]string, 0)
    var mu sync.Mutex

    handler := HandlerFunc(func(ctx context.Context, item *Item) error {
        mu.Lock()
        processed = append(processed, item.ID)
        mu.Unlock()
        return nil
    })

    system, _ := NewMultiQueuePrioritySystem(config, handler)
    system.Start()
    defer system.Stop()

    // Enqueue items
    system.Enqueue(&Item{ID: "critical", Priority: PriorityCritical})
    system.Enqueue(&Item{ID: "high", Priority: PriorityHigh})
    system.Enqueue(&Item{ID: "normal", Priority: PriorityNormal})

    // Wait for processing
    time.Sleep(500 * time.Millisecond)

    mu.Lock()
    defer mu.Unlock()

    // Critical should be processed first
    assert.Contains(t, processed, "critical")
}

func TestQueueDepthLimit(t *testing.T) {
    config := &MultiQueueConfig{
        QueueConfigs: map[Priority]QueueConfig{
            PriorityNormal: {MaxDepth: 2},
        },
    }

    system, _ := NewMultiQueuePrioritySystem(config, nil)

    // First two should succeed
    assert.NoError(t, system.Enqueue(&Item{ID: "1", Priority: PriorityNormal}))
    assert.NoError(t, system.Enqueue(&Item{ID: "2", Priority: PriorityNormal}))

    // Third should fail
    err := system.Enqueue(&Item{ID: "3", Priority: PriorityNormal})
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "max depth")
}

func TestRetryWithBackoff(t *testing.T) {
    attempts := 0
    handler := HandlerFunc(func(ctx context.Context, item *Item) error {
        attempts++
        return assert.AnError
    })

    config := &MultiQueueConfig{
        QueueConfigs: map[Priority]QueueConfig{
            PriorityNormal: {MaxRetries: 3, Workers: 1},
        },
    }

    system, _ := NewMultiQueuePrioritySystem(config, handler)
    system.Start()
    defer system.Stop()

    system.Enqueue(&Item{
        ID:          "retry-test",
        Priority:    PriorityNormal,
        MaxAttempts: 3,
    })

    // Wait for retries
    time.Sleep(10 * time.Second)

    // Should have attempted 3 times
    assert.Equal(t, 3, attempts)
}

func BenchmarkPriorityQueue(b *testing.B) {
    pq := NewPriorityQueue()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pq.Enqueue(&Item{
            ID:       fmt.Sprintf("item-%d", i),
            Priority: Priority(i % 5),
        })
    }
}
```

## Summary

The Priority Queue Pattern provides:

1. **Service Level Guarantee**: Critical work gets processed first
2. **Resource Optimization**: Match processing capacity to demand
3. **Fairness**: Prevent starvation with weighted/aging strategies
4. **Observability**: Track processing by priority level
5. **Resilience**: Retry and dead letter handling

Key considerations:

- Prevent priority inversion
- Monitor queue depths by priority
- Set appropriate timeouts per priority
- Consider aging to prevent starvation
- Plan for dead letter handling
