# EC-026: Competing Consumers Pattern

## Problem Formalization

### The Load Distribution Challenge

When processing messages from a queue, a single consumer may become a bottleneck. The Competing Consumers pattern enables multiple parallel consumers to process messages from the same queue, improving throughput and reliability.

#### Problem Statement

Given:

- Message queue Q with arrival rate λ (messages/second)
- Single consumer processing rate μ (messages/second)
- Target latency L_max

When λ approaches or exceeds μ, we need:

```
Find optimal consumer count N such that:
    Effective processing rate = N × μ × efficiency(N)
    Effective processing rate > λ
    Latency < L_max
    Cost is minimized

Where efficiency(N) accounts for:
    - Coordination overhead
    - Resource contention
    - Queue synchronization costs
```

### Queue Throughput Model

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Queue Throughput Analysis                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Single Consumer:                                                       │
│  ┌──────────┐    λ    ┌────────┐    μ    ┌──────────┐                 │
│  │  Queue   │────────►│Consumer│────────►│ Output   │                 │
│  │          │         │   1    │         │          │                 │
│  └──────────┘         └────────┘         └──────────┘                 │
│                                                                         │
│  Max throughput: μ                                                      │
│  If λ > μ: queue grows unbounded, latency → ∞                          │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Competing Consumers (N=3):                                             │
│  ┌──────────┐         ┌────────┐                                        │
│  │          │────────►│Consumer│───┐                                    │
│  │          │    λ    │   1    │   │                                    │
│  │          │         └────────┘   │    ┌──────────┐                   │
│  │  Queue   │────────►┌────────┐   ├───►│ Aggregated│                  │
│  │          │         │Consumer│   │    │  Output   │                  │
│  │          │         │   2    │───┘    └──────────┘                   │
│  │          │         └────────┘                                       │
│  │          │────────►┌────────┐                                       │
│  │          │         │Consumer│                                       │
│  └──────────┘         │   3    │                                       │
│                       └────────┘                                       │
│                                                                         │
│  Max throughput: 3μ (ideal)                                             │
│  Actual: 3μ × efficiency factor                                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Consumer Coordination

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Competing Consumers Architecture                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Message Source (Kafka/RabbitMQ/SQS)                                    │
│       │                                                                 │
│       │ Messages with partition keys                                    │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Partition Assignment                                            │   │
│  │  ┌───────────────┬─────────────────────────────────────────────┐ │   │
│  │  │  Partition 0  │ Consumer Group Coordination                 │ │   │
│  │  │  Partition 1  │                                             │ │   │
│  │  │  Partition 2  │  • Each partition consumed by one consumer  │ │   │
│  │  │  Partition 3  │  • Rebalance on consumer join/leave         │ │   │
│  │  │  ...          │  • Automatic failover                       │ │   │
│  │  └───────────────┴─────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Consumer Pool                                                   │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │   │
│  │  │Consumer 1│  │Consumer 2│  │Consumer 3│  │Consumer N│        │   │
│  │  │ (Pod 1)  │  │ (Pod 2)  │  │ (Pod 3)  │  │ (Pod N)  │        │   │
│  │  │          │  │          │  │          │  │          │        │   │
│  │  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │        │   │
│  │  │ │Worker│ │  │ │Worker│ │  │ │Worker│ │  │ │Worker│ │        │   │
│  │  │ │Pool  │ │  │ │Pool  │ │  │ │Pool  │ │  │ │Pool  │ │        │   │
│  │  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │        │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │   │
│  │       │              │              │              │            │   │
│  │       └──────────────┴──────────────┴──────────────┘            │   │
│  │                          │                                      │   │
│  │                   Auto-scaling (HPA)                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Work Distribution Strategies

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Work Distribution Strategies                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. Round-Robin (Simple, stateless)                                     │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Message 1 ──► Consumer 1                                       │   │
│  │  Message 2 ──► Consumer 2                                       │   │
│  │  Message 3 ──► Consumer 3                                       │   │
│  │  Message 4 ──► Consumer 1                                       │   │
│  │  ...                                                            │   │
│  │                                                                  │   │
│  │  Pros: Even distribution, simple                                │   │
│  │  Cons: No ordering guarantees, varying message sizes            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  2. Hash-based (Ordering guarantees)                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  hash(customer_id) % N ──► determines consumer                  │   │
│  │                                                                  │   │
│  │  Customer A orders ──► Consumer 1 (always)                      │   │
│  │  Customer B orders ──► Consumer 2 (always)                      │   │
│  │                                                                  │   │
│  │  Pros: Same key → same consumer (ordering preserved)            │   │
│  │  Cons: Hot keys create imbalance                                │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  3. Priority-based                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  High priority messages ──► Dedicated fast consumers            │   │
│  │  Normal messages ──► General pool                               │   │
│  │                                                                  │   │
│  │  Pros: SLA guarantees for critical work                         │   │
│  │  Cons: Complex resource management                              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  4. Load-based (Dynamic)                                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Monitor consumer load:                                         │   │
│  │  - CPU usage                                                    │   │
│  │  - Queue depth per consumer                                     │   │
│  │  - Processing latency                                           │   │
│  │                                                                  │   │
│  │  Route to least loaded consumer                                 │   │
│  │                                                                  │   │
│  │  Pros: Optimal resource utilization                             │   │
│  │  Cons: Requires coordination, complexity                        │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Consumer Group Implementation

```go
// pkg/consumers/consumer_group.go
package consumers

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/IBM/sarama"
    "go.uber.org/zap"
)

// ConsumerGroup manages a group of competing consumers
type ConsumerGroup struct {
    config     *Config
    group      sarama.ConsumerGroup
    handler    MessageHandler
    logger     *zap.Logger

    // Lifecycle
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup

    // Metrics
    messagesProcessed uint64
    messagesFailed    uint64

    // Ready signal
    ready        chan bool
}

type Config struct {
    Brokers          []string
    GroupID          string
    Topics           []string
    InitialOffset    int64 // sarama.OffsetNewest or sarama.OffsetOldest

    // Processing config
    WorkersPerConsumer int
    MaxProcessingTime  time.Duration
    RetryLimit         int

    // Kafka config
    SessionTimeout     time.Duration
    HeartbeatInterval  time.Duration
    RebalanceStrategy  string // range, roundrobin, sticky
}

// MessageHandler processes individual messages
type MessageHandler interface {
    Process(ctx context.Context, msg *Message) error
}

type Message struct {
    Topic     string
    Partition int32
    Offset    int64
    Key       []byte
    Value     []byte
    Timestamp time.Time
    Headers   map[string]string
}

func NewConsumerGroup(cfg *Config, handler MessageHandler, logger *zap.Logger) (*ConsumerGroup, error) {
    ctx, cancel := context.WithCancel(context.Background())

    // Configure Sarama
    saramaConfig := sarama.NewConfig()
    saramaConfig.Version = sarama.V2_6_0_0
    saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    saramaConfig.Consumer.Group.Session.Timeout = cfg.SessionTimeout
    saramaConfig.Consumer.Group.Heartbeat.Interval = cfg.HeartbeatInterval
    saramaConfig.Consumer.Offsets.Initial = cfg.InitialOffset
    saramaConfig.Consumer.MaxProcessingTime = cfg.MaxProcessingTime

    group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConfig)
    if err != nil {
        cancel()
        return nil, fmt.Errorf("creating consumer group: %w", err)
    }

    cg := &ConsumerGroup{
        config:  cfg,
        group:   group,
        handler: handler,
        logger:  logger,
        ctx:     ctx,
        cancel:  cancel,
        ready:   make(chan bool),
    }

    return cg, nil
}

func (cg *ConsumerGroup) Start() error {
    cg.wg.Add(1)
    go cg.consumeLoop()

    <-cg.ready // Wait until consumer is ready
    cg.logger.Info("consumer group started",
        zap.String("group_id", cg.config.GroupID),
        zap.Strings("topics", cg.config.Topics))

    return nil
}

func (cg *ConsumerGroup) Stop() error {
    cg.cancel()
    cg.wg.Wait()
    return cg.group.Close()
}

func (cg *ConsumerGroup) consumeLoop() {
    defer cg.wg.Done()

    for {
        select {
        case <-cg.ctx.Done():
            return
        default:
        }

        // Setup consumer group handler
        handler := &consumerGroupHandler{
            parent: cg,
            ready:  cg.ready,
        }

        // Consume
        err := cg.group.Consume(cg.ctx, cg.config.Topics, handler)
        if err != nil {
            cg.logger.Error("consume error", zap.Error(err))
        }

        // Check if context was cancelled
        if cg.ctx.Err() != nil {
            return
        }

        cg.ready = make(chan bool)
    }
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
    parent *ConsumerGroup
    ready  chan bool
}

func (h *consumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
    close(h.ready)
    h.parent.logger.Info("consumer group session setup",
        zap.String("member_id", sess.MemberID()))
    return nil
}

func (h *consumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
    h.parent.logger.Info("consumer group session cleanup",
        zap.String("member_id", sess.MemberID()))
    return nil
}

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    // Create worker pool for parallel processing within this consumer
    workerPool := make(chan struct{}, h.parent.config.WorkersPerConsumer)

    for {
        select {
        case msg := <-claim.Messages():
            if msg == nil {
                return nil
            }

            // Acquire worker slot
            workerPool <- struct{}{}

            go func(msg *sarama.ConsumerMessage) {
                defer func() { <-workerPool }()

                h.processMessage(sess, msg)
            }(msg)

        case <-sess.Context().Done():
            return nil
        }
    }
}

func (h *consumerGroupHandler) processMessage(sess sarama.ConsumerGroupSession, kafkaMsg *sarama.ConsumerMessage) {
    msg := &Message{
        Topic:     kafkaMsg.Topic,
        Partition: kafkaMsg.Partition,
        Offset:    kafkaMsg.Offset,
        Key:       kafkaMsg.Key,
        Value:     kafkaMsg.Value,
        Timestamp: kafkaMsg.Timestamp,
        Headers:   make(map[string]string),
    }

    for _, header := range kafkaMsg.Headers {
        msg.Headers[string(header.Key)] = string(header.Value)
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(h.parent.ctx, h.parent.config.MaxProcessingTime)
    defer cancel()

    // Process with retry
    var err error
    for attempt := 0; attempt <= h.parent.config.RetryLimit; attempt++ {
        err = h.parent.handler.Process(ctx, msg)
        if err == nil {
            break
        }

        if attempt < h.parent.config.RetryLimit {
            backoff := calculateBackoff(attempt)
            time.Sleep(backoff)
        }
    }

    if err != nil {
        h.parent.logger.Error("message processing failed",
            zap.Error(err),
            zap.String("topic", msg.Topic),
            zap.Int32("partition", msg.Partition),
            zap.Int64("offset", msg.Offset))

        // Send to DLQ or handle error
        h.handleError(msg, err)

        // Still mark as consumed to prevent blocking
        sess.MarkMessage(kafkaMsg, "")
    } else {
        sess.MarkMessage(kafkaMsg, "")
    }
}

func (h *consumerGroupHandler) handleError(msg *Message, err error) {
    // Implementation: send to dead letter queue, alert, etc.
}
```

### Worker Pool Implementation

```go
// pkg/consumers/worker_pool.go
package consumers

import (
    "context"
    "sync"
)

// WorkerPool manages a pool of workers for message processing
type WorkerPool struct {
    size      int
    jobs      chan Job
    results   chan Result
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
    processor Processor
}

type Job struct {
    Message *Message
    Context context.Context
}

type Result struct {
    Message *Message
    Error   error
}

type Processor func(ctx context.Context, msg *Message) error

func NewWorkerPool(size int, processor Processor) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())

    return &WorkerPool{
        size:      size,
        jobs:      make(chan Job, size*2),
        results:   make(chan Result, size*2),
        ctx:       ctx,
        cancel:    cancel,
        processor: processor,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.size; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) Stop() {
    wp.cancel()
    close(wp.jobs)
    wp.wg.Wait()
    close(wp.results)
}

func (wp *WorkerPool) Submit(msg *Message) <-chan Result {
    resultCh := make(chan Result, 1)

    select {
    case wp.jobs <- Job{Message: msg, Context: wp.ctx}:
        // Job queued
    case <-wp.ctx.Done():
        resultCh <- Result{Message: msg, Error: wp.ctx.Err()}
        close(resultCh)
        return resultCh
    }

    // Get result asynchronously
    go func() {
        select {
        case result := <-wp.results:
            resultCh <- result
        case <-wp.ctx.Done():
            resultCh <- Result{Message: msg, Error: wp.ctx.Err()}
        }
        close(resultCh)
    }()

    return resultCh
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()

    for job := range wp.jobs {
        err := wp.processor(job.Context, job.Message)

        select {
        case wp.results <- Result{Message: job.Message, Error: err}:
        case <-wp.ctx.Done():
            return
        }
    }
}
```

### Dynamic Scaling

```go
// pkg/consumers/scaler.go
package consumers

import (
    "context"
    "fmt"
    "math"
    "time"

    "go.uber.org/zap"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

// AutoScaler automatically adjusts consumer count based on load
type AutoScaler struct {
    config     ScalerConfig
    k8sClient  kubernetes.Interface
    logger     *zap.Logger

    // Metrics
    metrics    *QueueMetrics

    // Control
    ctx        context.Context
    cancel     context.CancelFunc
}

type ScalerConfig struct {
    MinReplicas       int
    MaxReplicas       int
    TargetLag         int64 // Target consumer lag
    ScaleUpThreshold  float64
    ScaleDownThreshold float64
    CooldownPeriod    time.Duration
    EvaluationPeriod  time.Duration

    // Kubernetes
    DeploymentName    string
    Namespace         string
}

type QueueMetrics struct {
    ConsumerLag      int64
    ProcessingRate   float64
    MessageRate      float64
    ConsumerUtilization float64
}

func NewAutoScaler(cfg ScalerConfig, logger *zap.Logger) (*AutoScaler, error) {
    // Create Kubernetes client
    k8sConfig, err := rest.InClusterConfig()
    if err != nil {
        return nil, fmt.Errorf("getting k8s config: %w", err)
    }

    k8sClient, err := kubernetes.NewForConfig(k8sConfig)
    if err != nil {
        return nil, fmt.Errorf("creating k8s client: %w", err)
    }

    ctx, cancel := context.WithCancel(context.Background())

    return &AutoScaler{
        config:    cfg,
        k8sClient: k8sClient,
        logger:    logger,
        ctx:       ctx,
        cancel:    cancel,
    }, nil
}

func (as *AutoScaler) Start() {
    go as.scalingLoop()
}

func (as *AutoScaler) Stop() {
    as.cancel()
}

func (as *AutoScaler) scalingLoop() {
    ticker := time.NewTicker(as.config.EvaluationPeriod)
    defer ticker.Stop()

    lastScaleTime := time.Time{}

    for {
        select {
        case <-as.ctx.Done():
            return
        case <-ticker.C:
            // Check cooldown
            if time.Since(lastScaleTime) < as.config.CooldownPeriod {
                continue
            }

            // Get current metrics
            metrics := as.collectMetrics()

            // Calculate desired replicas
            desired := as.calculateDesiredReplicas(metrics)

            // Get current replicas
            current, err := as.getCurrentReplicas()
            if err != nil {
                as.logger.Error("failed to get current replicas", zap.Error(err))
                continue
            }

            // Scale if needed
            if desired != current {
                as.logger.Info("scaling consumers",
                    zap.Int("current", current),
                    zap.Int("desired", desired),
                    zap.Int64("lag", metrics.ConsumerLag))

                if err := as.scale(desired); err != nil {
                    as.logger.Error("scaling failed", zap.Error(err))
                    continue
                }

                lastScaleTime = time.Now()
            }
        }
    }
}

func (as *AutoScaler) collectMetrics() *QueueMetrics {
    // In production, these would come from Kafka admin client or metrics system
    // For example, using Kafka's AdminClient to describe consumer groups

    return &QueueMetrics{
        ConsumerLag:      getConsumerLag(),      // Sum of lag across all partitions
        ProcessingRate:   getProcessingRate(),   // Messages/second processed
        MessageRate:      getMessageRate(),      // Messages/second produced
        ConsumerUtilization: getConsumerUtilization(),
    }
}

func (as *AutoScaler) calculateDesiredReplicas(metrics *QueueMetrics) int {
    if metrics.ProcessingRate == 0 {
        return as.config.MinReplicas
    }

    // Calculate based on lag
    // replicas = current_lag / (target_processing_time * processing_rate)

    // Simple formula: desired = ceil(lag / target_lag)
    desiredFloat := float64(metrics.ConsumerLag) / float64(as.config.TargetLag)
    desired := int(math.Ceil(desiredFloat))

    // Apply bounds
    if desired < as.config.MinReplicas {
        desired = as.config.MinReplicas
    }
    if desired > as.config.MaxReplicas {
        desired = as.config.MaxReplicas
    }

    // Consider message rate vs processing rate
    if metrics.MessageRate > metrics.ProcessingRate {
        // Need more consumers
        bufferFactor := metrics.MessageRate / metrics.ProcessingRate
        desired = int(math.Ceil(float64(desired) * bufferFactor))
    }

    // Apply bounds again after adjustment
    if desired > as.config.MaxReplicas {
        desired = as.config.MaxReplicas
    }

    return desired
}

func (as *AutoScaler) getCurrentReplicas() (int, error) {
    deployment, err := as.k8sClient.AppsV1().Deployments(as.config.Namespace).
        Get(as.ctx, as.config.DeploymentName, metav1.GetOptions{})
    if err != nil {
        return 0, err
    }

    return int(deployment.Status.Replicas), nil
}

func (as *AutoScaler) scale(replicas int) error {
    deployment, err := as.k8sClient.AppsV1().Deployments(as.config.Namespace).
        Get(as.ctx, as.config.DeploymentName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    deployment.Spec.Replicas = int32Ptr(int32(replicas))

    _, err = as.k8sClient.AppsV1().Deployments(as.config.Namespace).
        Update(as.ctx, deployment, metav1.UpdateOptions{})
    return err
}

func int32Ptr(i int32) *int32 {
    return &i
}
```

## Trade-off Analysis

### Consumer Count Optimization

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Optimal Consumer Count                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Throughput vs Consumer Count:                                          │
│                                                                         │
│  Throughput ↑                                                           │
│             │        ┌────────────────────────┐                        │
│             │       /  Saturation Point      │                        │
│             │      /   (Network/IO limited)  │                        │
│             │     /                          │                        │
│             │    /  Linear scaling           │                        │
│             │   /                            │                        │
│             │  /                             │                        │
│             │ /                              │                        │
│             │/                               │                        │
│             └────────────────────────────────────────►                  │
│                1    5    10   15   20   25   30                        │
│                        Consumer Count                                   │
│                                                                         │
│  Key Metrics:                                                           │
│  • Latency: Initially flat, then increases due to contention           │
│  • Throughput: Linear until saturation, then plateaus                  │
│  • Cost: Linear with consumer count                                     │
│  • Coordination overhead: O(N²) for rebalances                          │
│                                                                         │
│  Rule of thumb: Consumers = Partitions (up to saturation)               │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Consumer Group Testing

```go
// test/consumers/consumer_test.go
package consumers

import (
    "context"
    "sync"
    "sync/atomic"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestConsumerGroupProcessing(t *testing.T) {
    processed := uint32(0)
    var mu sync.Mutex
    received := make(map[string]bool)

    handler := MessageHandlerFunc(func(ctx context.Context, msg *Message) error {
        atomic.AddUint32(&processed, 1)
        mu.Lock()
        received[string(msg.Key)] = true
        mu.Unlock()
        return nil
    })

    cfg := &Config{
        Brokers:            []string{"localhost:9092"},
        GroupID:            "test-group",
        Topics:             []string{"test-topic"},
        WorkersPerConsumer: 3,
    }

    cg, err := NewConsumerGroup(cfg, handler, zap.NewNop())
    require.NoError(t, err)

    err = cg.Start()
    require.NoError(t, err)
    defer cg.Stop()

    // Produce test messages
    produceTestMessages(t, "test-topic", 100)

    // Wait for processing
    time.Sleep(5 * time.Second)

    // Verify
    count := atomic.LoadUint32(&processed)
    assert.GreaterOrEqual(t, count, uint32(90)) // At least 90% processed
}

func TestConsumerRebalance(t *testing.T) {
    // Test that when consumers join/leave, partitions are redistributed

    cfg := &Config{
        Brokers: []string{"localhost:9092"},
        GroupID: "rebalance-test",
        Topics:  []string{"rebalance-topic"},
    }

    // Start first consumer
    cg1, _ := NewConsumerGroup(cfg, dummyHandler, zap.NewNop())
    cg1.Start()
    defer cg1.Stop()

    time.Sleep(2 * time.Second)

    // Start second consumer
    cg2, _ := NewConsumerGroup(cfg, dummyHandler, zap.NewNop())
    cg2.Start()
    defer cg2.Stop()

    // Both should receive partitions
    // Verify by checking logs or metrics
}

func TestWorkerPool(t *testing.T) {
    processed := 0
    var mu sync.Mutex

    processor := func(ctx context.Context, msg *Message) error {
        mu.Lock()
        processed++
        mu.Unlock()
        return nil
    }

    pool := NewWorkerPool(3, processor)
    pool.Start()
    defer pool.Stop()

    // Submit jobs
    for i := 0; i < 10; i++ {
        msg := &Message{Key: []byte(fmt.Sprintf("key-%d", i))}
        resultCh := pool.Submit(msg)

        select {
        case result := <-resultCh:
            assert.NoError(t, result.Error)
        case <-time.After(time.Second):
            t.Fatal("timeout waiting for result")
        }
    }

    mu.Lock()
    assert.Equal(t, 10, processed)
    mu.Unlock()
}

func BenchmarkConsumerThroughput(b *testing.B) {
    handler := MessageHandlerFunc(func(ctx context.Context, msg *Message) error {
        return nil
    })

    pool := NewWorkerPool(10, handler.Process)
    pool.Start()
    defer pool.Stop()

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            msg := &Message{Value: []byte("test")}
            <-pool.Submit(msg)
        }
    })
}
```

## Summary

The Competing Consumers Pattern provides:

1. **Horizontal Scaling**: Add consumers to increase throughput
2. **High Availability**: Automatic failover when consumers fail
3. **Load Distribution**: Even distribution across consumers
4. **Ordering Guarantees**: Per-partition ordering preserved
5. **Dynamic Scaling**: Adjust to changing load patterns

Key considerations:

- Number of partitions limits parallelism
- Rebalance costs during scaling
- Message ordering requirements
- Consumer group coordination overhead
- Dead letter handling for failed messages
