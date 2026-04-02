# EC-029: Sequential Convoy Pattern

## Problem Formalization

### The Ordered Processing Challenge

In distributed systems, some operations must be processed in strict order for the same entity, while still allowing parallel processing across different entities. The Sequential Convoy pattern enables high-throughput ordered processing.

#### Problem Statement

Given:

- Message stream M with messages m₁, m₂, ..., mₙ
- Each message has an entity key k(m)
- Ordering constraint: For any entity e, messages with k(m) = e must be processed sequentially
- Throughput requirement: Process messages as fast as possible

Find an architecture that:

```
Maximizes: Throughput(M)
Subject to:
    - For all entities e: Process(mᵢ) before Process(mⱼ) if i < j and k(mᵢ) = k(mⱼ) = e
    - Horizontal scalability (add workers to increase throughput)
    - Fault tolerance (handle worker failures)
```

### Ordering Requirements

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Ordering Requirement Examples                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Example 1: Bank Account Transactions                                   │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Account A-123:                                                  │   │
│  │  ┌─────┐   ┌─────┐   ┌─────┐   ┌─────┐                         │   │
│  │  │$100 │──►│-$20 │──►│+$50 │──►│-$30 │  Must process in order!  │   │
│  │  │Dep  │   │W/d  │   │Dep  │   │W/d  │  Balance: 100→80→130→100 │   │
│  │  └─────┘   └─────┘   └─────┘   └─────┘                         │   │
│  │                                                                  │   │
│  │  Account B-456:                                                  │   │
│  │  ┌─────┐   ┌─────┐                                             │   │
│  │  │+$200│──►│-$100│  Can process in parallel with A-123         │   │
│  │  └─────┘   └─────┘                                             │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Example 2: Inventory Management                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Product SKU-XYZ:                                                │   │
│  │  ┌────────┐   ┌────────┐   ┌────────┐                          │   │
│  │  │Stock+10│──►│Stock-5 │──►│Stock+20│  Order matters!           │   │
│  │  └────────┘   └────────┘   └────────┘  Current: 10→15→35       │   │
│  │                                                                  │   │
│  │  Wrong order could result in negative stock or overselling       │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Sequential Convoy Pattern

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Sequential Convoy Architecture                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Partition by Entity Key                                         │   │
│  │                                                                  │   │
│  │  Input Stream:                                                   │   │
│  │  [A:1] [B:1] [A:2] [C:1] [A:3] [B:2] [C:2] ...                   │   │
│  │    │     │     │     │     │     │     │                        │   │
│  │    └─────┴─────┴─────┴─────┴─────┴─────┘                        │   │
│  │                          │                                       │   │
│  │                          ▼                                       │   │
│  │            Hash(entity_key) % N_partitions                       │   │
│  │                          │                                       │   │
│  │          ┌───────────────┼───────────────┐                       │   │
│  │          ▼               ▼               ▼                       │   │
│  │     ┌─────────┐     ┌─────────┐     ┌─────────┐                  │   │
│  │     │Part 0   │     │Part 1   │     │Part 2   │                  │   │
│  │     │[A:1,A:2│     │[B:1,B:2│     │[C:1,C:2│                  │   │
│  │     │ A:3...]│     │ ...]    │     │ ...]    │                  │   │
│  │     └────┬────┘     └────┬────┘     └────┬────┘                  │   │
│  │          │               │               │                       │   │
│  └──────────┼───────────────┼───────────────┼───────────────────────┘   │
│             │               │               │                           │
│             ▼               ▼               ▼                           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  One Consumer per Partition (Guarantees Ordering)                │   │
│  │                                                                  │   │
│  │  Consumer 0          Consumer 1          Consumer 2              │   │
│  │  ┌─────────┐         ┌─────────┐         ┌─────────┐             │   │
│  │  │Process A│         │Process B│         │Process C│             │   │
│  │  │messages │         │messages │         │messages │             │   │
│  │  │in order │         │in order │         │in order │             │   │
│  │  └─────────┘         └─────────┘         └─────────┘             │   │
│  │                                                                  │   │
│  │  Scaling: Add partitions + consumers                             │   │
│  │  (More partitions = more parallelism)                            │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Convoy Implementation Strategies

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Convoy Implementation Strategies                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Strategy 1: Partitioned Message Queue (Kafka)                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Kafka Topic with 3 partitions:                                  │   │
│  │  partition = hash(entity_id) % 3                                 │   │
│  │                                                                  │   │
│  │  Guarantees:                                                     │   │
│  │  • Same entity ID → same partition                               │   │
│  │  • Same partition → same consumer                                │   │
│  │  • Single consumer → ordered processing                          │   │
│  │                                                                  │   │
│  │  Limitation: Uneven distribution if skewed keys                  │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Strategy 2: Actor Model (per-entity workers)                           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐             │   │
│  │  │ Actor A │  │ Actor B │  │ Actor C │  │ Actor D │             │   │
│  │  │ (inbox) │  │ (inbox) │  │ (inbox) │  │ (inbox) │             │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘             │   │
│  │       └─────────────┴─────────────┴─────────────┘                │   │
│  │                     │                                            │   │
│  │                     ▼                                            │   │
│  │            Router (hash to actor)                                │   │
│  │                                                                  │   │
│  │  Actors created on-demand, garbage collected after inactivity   │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Strategy 3: Striped Lock Pool                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Worker Pool ──► Hash Ring ──► Lock Stripes                     │   │
│  │                                                                  │   │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐                    │   │
│  │  │Lock 0  │ │Lock 1  │ │Lock 2  │ │Lock 3  │                    │   │
│  │  │[A,C,E] │ │[B,D,F] │ │[G,H...]│ │[...]   │                    │   │
│  │  └────────┘ └────────┘ └────────┘ └────────┘                    │   │
│  │                                                                  │   │
│  │  Worker acquires lock for entity before processing              │   │
│  │  (More locks than workers to reduce contention)                 │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Kafka-Based Sequential Convoy

```go
// pkg/convoy/kafka_convoy.go
package convoy

import (
    "context"
    "fmt"
    "hash/fnv"
    "sync"
    "time"

    "github.com/IBM/sarama"
    "go.uber.org/zap"
)

// KafkaConvoy implements sequential convoy using Kafka partitions
type KafkaConvoy struct {
    config      *KafkaConfig
    consumer    sarama.ConsumerGroup
    producer    sarama.AsyncProducer
    handler     MessageHandler
    partitioner Partitioner
    logger      *zap.Logger

    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

type KafkaConfig struct {
    Brokers         []string
    Topic           string
    GroupID         string
    NumPartitions   int
    ReplicationFactor int

    // Processing
    MaxProcessingTime time.Duration
    RetryLimit        int
}

// Partitioner determines partition from entity key
type Partitioner interface {
    Partition(key string, numPartitions int) int32
}

// HashPartitioner uses FNV hash
type HashPartitioner struct{}

func (h *HashPartitioner) Partition(key string, numPartitions int) int32 {
    hasher := fnv.New32a()
    hasher.Write([]byte(key))
    return int32(hasher.Sum32() % uint32(numPartitions))
}

// MessageHandler processes messages in order
type MessageHandler interface {
    Handle(ctx context.Context, msg *OrderedMessage) error
    // Should return error to trigger retry/dead letter
}

type OrderedMessage struct {
    EntityKey   string
    Sequence    int64
    Payload     []byte
    Timestamp   time.Time
    Partition   int32
    Offset      int64
}

func NewKafkaConvoy(cfg *KafkaConfig, handler MessageHandler, logger *zap.Logger) (*KafkaConvoy, error) {
    ctx, cancel := context.WithCancel(context.Background())

    // Setup producer
    producerConfig := sarama.NewConfig()
    producerConfig.Producer.Partitioner = sarama.NewManualPartitioner
    producerConfig.Producer.RequiredAcks = sarama.WaitForAll
    producerConfig.Producer.Return.Successes = true

    producer, err := sarama.NewAsyncProducer(cfg.Brokers, producerConfig)
    if err != nil {
        cancel()
        return nil, fmt.Errorf("creating producer: %w", err)
    }

    // Setup consumer
    consumerConfig := sarama.NewConfig()
    consumerConfig.Version = sarama.V2_6_0_0
    consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
    consumerConfig.Consumer.MaxProcessingTime = cfg.MaxProcessingTime

    consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, consumerConfig)
    if err != nil {
        producer.Close()
        cancel()
        return nil, fmt.Errorf("creating consumer: %w", err)
    }

    return &KafkaConvoy{
        config:      cfg,
        consumer:    consumer,
        producer:    producer,
        handler:     handler,
        partitioner: &HashPartitioner{},
        logger:      logger,
        ctx:         ctx,
        cancel:      cancel,
    }, nil
}

// Send sends a message to the convoy, ensuring ordering by entity key
func (kc *KafkaConvoy) Send(ctx context.Context, entityKey string, payload []byte) error {
    partition := kc.partitioner.Partition(entityKey, kc.config.NumPartitions)

    msg := &sarama.ProducerMessage{
        Topic:     kc.config.Topic,
        Partition: partition,
        Key:       sarama.StringEncoder(entityKey),
        Value:     sarama.ByteEncoder(payload),
        Headers: []sarama.RecordHeader{
            {Key: []byte("entity_key"), Value: []byte(entityKey)},
        },
    }

    select {
    case kc.producer.Input() <- msg:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (kc *KafkaConvoy) Start() error {
    handler := &convoyConsumerHandler{
        convoy: kc,
        ready:  make(chan bool),
    }

    kc.wg.Add(1)
    go func() {
        defer kc.wg.Done()

        for {
            select {
            case <-kc.ctx.Done():
                return
            default:
            }

            err := kc.consumer.Consume(kc.ctx, []string{kc.config.Topic}, handler)
            if err != nil {
                kc.logger.Error("consume error", zap.Error(err))
            }
        }
    }()

    <-handler.ready
    kc.logger.Info("convoy started", zap.String("topic", kc.config.Topic))
    return nil
}

func (kc *KafkaConvoy) Stop() error {
    kc.cancel()
    kc.wg.Wait()

    if err := kc.consumer.Close(); err != nil {
        kc.logger.Error("consumer close error", zap.Error(err))
    }

    return kc.producer.Close()
}

// convoyConsumerHandler implements sarama.ConsumerGroupHandler
type convoyConsumerHandler struct {
    convoy *KafkaConvoy
    ready  chan bool
}

func (h *convoyConsumerHandler) Setup(sess sarama.ConsumerGroupSession) error {
    close(h.ready)
    return nil
}

func (h *convoyConsumerHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
    return nil
}

func (h *convoyConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    // Process messages sequentially - this is the key!
    // Kafka guarantees order within a partition
    // Consumer group assigns one consumer per partition

    for msg := range claim.Messages() {
        orderedMsg := &OrderedMessage{
            EntityKey: string(msg.Key),
            Payload:   msg.Value,
            Timestamp: msg.Timestamp,
            Partition: msg.Partition,
            Offset:    msg.Offset,
        }

        // Process with retry logic
        var err error
        for attempt := 0; attempt <= h.convoy.config.RetryLimit; attempt++ {
            ctx, cancel := context.WithTimeout(h.convoy.ctx, h.convoy.config.MaxProcessingTime)
            err = h.convoy.handler.Handle(ctx, orderedMsg)
            cancel()

            if err == nil {
                break
            }

            h.convoy.logger.Warn("message processing failed, retrying",
                zap.Error(err),
                zap.String("entity_key", orderedMsg.EntityKey),
                zap.Int("attempt", attempt))

            if attempt < h.convoy.config.RetryLimit {
                time.Sleep(time.Duration(attempt+1) * time.Second)
            }
        }

        if err != nil {
            h.convoy.logger.Error("message processing exhausted retries",
                zap.Error(err),
                zap.String("entity_key", orderedMsg.EntityKey))

            // Send to dead letter queue
            h.sendToDeadLetter(orderedMsg, err)
        }

        // Mark message as processed
        sess.MarkMessage(msg, "")
    }

    return nil
}

func (h *convoyConsumerHandler) sendToDeadLetter(msg *OrderedMessage, err error) {
    // Implementation: send to DLQ topic
}
```

### In-Memory Actor-Based Convoy

```go
// pkg/convoy/actor_convoy.go
package convoy

import (
    "context"
    "fmt"
    "hash/fnv"
    "sync"
    "time"
)

// ActorConvoy implements sequential convoy using actors per entity
type ActorConvoy struct {
    config      *ActorConfig
    actors      map[string]*EntityActor
    actorMu     sync.RWMutex
    workerPool  chan struct{}
    handler     MessageHandler

    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

type ActorConfig struct {
    MaxActors        int           // Maximum concurrent actors
    ActorIdleTimeout time.Duration // How long to keep idle actors
    WorkerPoolSize   int           // Concurrent workers
}

// EntityActor processes messages for a single entity
type EntityActor struct {
    entityKey   string
    inbox       chan *OrderedMessage
    lastActive  time.Time
    ctx         context.Context
    cancel      context.CancelFunc
    handler     MessageHandler
    wg          sync.WaitGroup
}

func NewActorConvoy(cfg *ActorConfig, handler MessageHandler) *ActorConvoy {
    ctx, cancel := context.WithCancel(context.Background())

    return &ActorConvoy{
        config:     cfg,
        actors:     make(map[string]*EntityActor),
        workerPool: make(chan struct{}, cfg.WorkerPoolSize),
        handler:    handler,
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (ac *ActorConvoy) Start() {
    // Initialize worker pool tokens
    for i := 0; i < ac.config.WorkerPoolSize; i++ {
        ac.workerPool <- struct{}{}
    }

    // Start actor cleanup goroutine
    ac.wg.Add(1)
    go ac.cleanupLoop()
}

func (ac *ActorConvoy) Stop() {
    ac.cancel()

    // Stop all actors
    ac.actorMu.Lock()
    for _, actor := range ac.actors {
        actor.stop()
    }
    ac.actorMu.Unlock()

    ac.wg.Wait()
}

func (ac *ActorConvoy) Send(entityKey string, msg *OrderedMessage) error {
    // Get or create actor for entity
    actor := ac.getOrCreateActor(entityKey)

    select {
    case actor.inbox <- msg:
        return nil
    case <-ac.ctx.Done():
        return ac.ctx.Err()
    }
}

func (ac *ActorConvoy) getOrCreateActor(entityKey string) *EntityActor {
    ac.actorMu.RLock()
    actor, exists := ac.actors[entityKey]
    ac.actorMu.RUnlock()

    if exists {
        actor.lastActive = time.Now()
        return actor
    }

    // Create new actor
    ac.actorMu.Lock()
    defer ac.actorMu.Unlock()

    // Double-check
    if actor, exists = ac.actors[entityKey]; exists {
        return actor
    }

    actor = newEntityActor(entityKey, ac.handler)
    ac.actors[entityKey] = actor
    actor.start()

    return actor
}

func (ac *ActorConvoy) cleanupLoop() {
    defer ac.wg.Done()

    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ac.ctx.Done():
            return
        case <-ticker.C:
            ac.cleanupIdleActors()
        }
    }
}

func (ac *ActorConvoy) cleanupIdleActors() {
    ac.actorMu.Lock()
    defer ac.actorMu.Unlock()

    for key, actor := range ac.actors {
        if time.Since(actor.lastActive) > ac.config.ActorIdleTimeout {
            actor.stop()
            delete(ac.actors, key)
        }
    }
}

func newEntityActor(entityKey string, handler MessageHandler) *EntityActor {
    ctx, cancel := context.WithCancel(context.Background())

    return &EntityActor{
        entityKey:  entityKey,
        inbox:      make(chan *OrderedMessage, 100),
        lastActive: time.Now(),
        ctx:        ctx,
        cancel:     cancel,
        handler:    handler,
    }
}

func (a *EntityActor) start() {
    a.wg.Add(1)
    go a.processLoop()
}

func (a *EntityActor) stop() {
    a.cancel()
    a.wg.Wait()
}

func (a *EntityActor) processLoop() {
    defer a.wg.Done()

    for {
        select {
        case <-a.ctx.Done():
            // Process remaining messages before stopping
            for {
                select {
                case msg := <-a.inbox:
                    a.processMessage(msg)
                default:
                    return
                }
            }

        case msg := <-a.inbox:
            a.processMessage(msg)
            a.lastActive = time.Now()
        }
    }
}

func (a *EntityActor) processMessage(msg *OrderedMessage) {
    ctx, cancel := context.WithTimeout(a.ctx, 30*time.Second)
    defer cancel()

    if err := a.handler.Handle(ctx, msg); err != nil {
        // Log error, potentially send to DLQ
    }
}
```

### Stripe-Based Lock Convoy

```go
// pkg/convoy/stripe_convoy.go
package convoy

import (
    "context"
    "hash/fnv"
    "sync"
    "sync/atomic"
)

// StripeConvoy uses lock striping for ordered processing
type StripeConvoy struct {
    config      *StripeConfig
    stripes     []sync.Mutex
    workerPool  chan struct{}
    handler     MessageHandler

    inFlight    int64
    ctx         context.Context
    cancel      context.CancelFunc
}

type StripeConfig struct {
    NumStripes     int // More stripes = less contention
    WorkerPoolSize int
}

func NewStripeConvoy(cfg *StripeConfig, handler MessageHandler) *StripeConvoy {
    ctx, cancel := context.WithCancel(context.Background())

    stripes := make([]sync.Mutex, cfg.NumStripes)

    return &StripeConvoy{
        config:     cfg,
        stripes:    stripes,
        workerPool: make(chan struct{}, cfg.WorkerPoolSize),
        handler:    handler,
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (sc *StripeConvoy) Start() {
    for i := 0; i < sc.config.WorkerPoolSize; i++ {
        sc.workerPool <- struct{}{}
    }
}

func (sc *StripeConvoy) Stop() {
    sc.cancel()

    // Wait for in-flight messages
    for atomic.LoadInt64(&sc.inFlight) > 0 {
        time.Sleep(10 * time.Millisecond)
    }
}

func (sc *StripeConvoy) Send(msg *OrderedMessage) error {
    stripeIdx := sc.getStripeIndex(msg.EntityKey)

    select {
    case <-sc.workerPool:
        // Acquired worker
    case <-sc.ctx.Done():
        return sc.ctx.Err()
    }

    atomic.AddInt64(&sc.inFlight, 1)

    go func() {
        defer func() {
            sc.workerPool <- struct{}{}
            atomic.AddInt64(&sc.inFlight, -1)
        }()

        // Acquire stripe lock for this entity
        sc.stripes[stripeIdx].Lock()
        defer sc.stripes[stripeIdx].Unlock()

        // Process with timeout
        ctx, cancel := context.WithTimeout(sc.ctx, 30*time.Second)
        defer cancel()

        _ = sc.handler.Handle(ctx, msg)
    }()

    return nil
}

func (sc *StripeConvoy) getStripeIndex(entityKey string) int {
    hasher := fnv.New32a()
    hasher.Write([]byte(entityKey))
    return int(hasher.Sum32() % uint32(sc.config.NumStripes))
}
```

## Trade-off Analysis

### Convoy Strategy Comparison

| Strategy | Ordering Guarantee | Throughput | Latency | Complexity | Best For |
|----------|-------------------|------------|---------|------------|----------|
| **Kafka Partition** | Strong | High | Low-Med | Low | Most use cases |
| **Actor Model** | Strong | Med | Low | Med | Variable entity count |
| **Lock Stripes** | Strong | Med-High | Low | Low | In-memory processing |
| **Database Queue** | Strong | Med | Med | Med | Small scale |

### Key Distribution Impact

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Key Distribution Effects                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Good Distribution (Uniform):                                           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Entities: A, B, C, D, E, F, G, H (even message rates)          │   │
│  │                                                                  │   │
│  │  Partitions:                                                     │   │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐                    │   │
│  │  │ P0: A,E│ │ P1: B,F│ │ P2: C,G│ │ P3: D,H│                    │   │
│  │  │  25%   │ │  25%   │ │  25%   │ │  25%   │                    │   │
│  │  └────────┘ └────────┘ └────────┘ └────────┘                    │   │
│  │                                                                  │   │
│  │  Result: Perfect parallelization                                │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Bad Distribution (Skewed):                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Entities: A (80%), B (5%), C (5%), D (5%), E (5%)              │   │
│  │                                                                  │   │
│  │  Partitions:                                                     │   │
│  │  ┌─────────────────────────────────────────────────────────┐    │   │
│  │  │ P0: A (80% load)                                        │    │   │
│  │  │ P1: B (5% load)    P2: C (5% load)    P3: D,E (10%)     │    │   │
│  │  └─────────────────────────────────────────────────────────┘    │   │
│  │                                                                  │   │
│  │  Result: P0 is bottleneck, others underutilized                 │   │
│  │                                                                  │   │
│  │  Solutions:                                                     │   │
│  │  • Use sub-entity keys (A-1, A-2, A-3, ...)                    │   │
│  │  • Split hot entities                                           │   │
│  │  • Use separate queue for hot entities                          │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Sequential Convoy Testing

```go
// test/convoy/convoy_test.go
package convoy

import (
    "context"
    "sync"
    "sync/atomic"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSequentialProcessing(t *testing.T) {
    // Track processing order per entity
    processed := make(map[string][]int)
    var mu sync.Mutex

    handler := MessageHandlerFunc(func(ctx context.Context, msg *OrderedMessage) error {
        mu.Lock()
        defer mu.Unlock()

        seq := extractSequence(msg.Payload)
        processed[msg.EntityKey] = append(processed[msg.EntityKey], seq)
        return nil
    })

    convoy := NewActorConvoy(&ActorConfig{
        MaxActors:        10,
        ActorIdleTimeout: time.Minute,
        WorkerPoolSize:   5,
    }, handler)

    convoy.Start()
    defer convoy.Stop()

    // Send messages out of order
    entities := []string{"A", "B", "C"}
    for _, entity := range entities {
        // Send 5, 3, 1, 4, 2 (out of order)
        for _, seq := range []int{5, 3, 1, 4, 2} {
            msg := &OrderedMessage{
                EntityKey: entity,
                Payload:   []byte(fmt.Sprintf("seq=%d", seq)),
            }
            convoy.Send(entity, msg)
        }
    }

    // Wait for processing
    time.Sleep(time.Second)

    // Verify order within each entity
    mu.Lock()
    defer mu.Unlock()

    for _, entity := range entities {
        sequences := processed[entity]
        assert.Equal(t, []int{5, 3, 1, 4, 2}, sequences,
            "Entity %s should be in send order", entity)
    }
}

func TestParallelEntities(t *testing.T) {
    var concurrentEntities int64
    var maxConcurrent int64
    var mu sync.Mutex

    handler := MessageHandlerFunc(func(ctx context.Context, msg *OrderedMessage) error {
        atomic.AddInt64(&concurrentEntities, 1)

        current := atomic.LoadInt64(&concurrentEntities)
        for {
            max := atomic.LoadInt64(&maxConcurrent)
            if current <= max || atomic.CompareAndSwapInt64(&maxConcurrent, max, current) {
                break
            }
        }

        time.Sleep(10 * time.Millisecond) // Simulate work

        atomic.AddInt64(&concurrentEntities, -1)
        return nil
    })

    convoy := NewActorConvoy(&ActorConfig{
        MaxActors:      100,
        WorkerPoolSize: 10,
    }, handler)

    convoy.Start()
    defer convoy.Stop()

    // Send to many different entities
    for i := 0; i < 50; i++ {
        entity := fmt.Sprintf("entity-%d", i)
        msg := &OrderedMessage{EntityKey: entity}
        convoy.Send(entity, msg)
    }

    time.Sleep(500 * time.Millisecond)

    // Should have processed multiple entities concurrently
    assert.Greater(t, atomic.LoadInt64(&maxConcurrent), int64(1))
}

func TestPartitioning(t *testing.T) {
    partitioner := &HashPartitioner{}
    numPartitions := 10

    // Test that same key always goes to same partition
    for i := 0; i < 100; i++ {
        p1 := partitioner.Partition("entity-A", numPartitions)
        assert.Equal(t, partitioner.Partition("entity-A", numPartitions), p1)
    }

    // Test distribution
    distribution := make(map[int32]int)
    for i := 0; i < 10000; i++ {
        key := fmt.Sprintf("entity-%d", i)
        p := partitioner.Partition(key, numPartitions)
        distribution[p]++
    }

    // All partitions should have some data
    assert.Equal(t, numPartitions, len(distribution))

    // Distribution should be roughly even (within 30%)
    expected := 10000 / numPartitions
    for _, count := range distribution {
        assert.InDelta(t, expected, count, float64(expected)*0.3)
    }
}
```

## Summary

The Sequential Convoy Pattern provides:

1. **Ordering Guarantee**: Messages for same entity processed sequentially
2. **Parallelism**: Different entities processed concurrently
3. **Scalability**: Add partitions/workers for more throughput
4. **Fault Tolerance**: Failed messages don't block entire queue
5. **Flexibility**: Multiple implementation strategies

Key considerations:

- Entity key distribution affects parallelism
- Hot entities can create bottlenecks
- Monitor per-partition lag
- Handle poison messages without blocking
- Consider dead letter queues for failed messages
