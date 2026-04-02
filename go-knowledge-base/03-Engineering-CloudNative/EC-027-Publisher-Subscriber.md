# EC-027: Publisher-Subscriber Pattern

## Problem Formalization

### The Decoupling Challenge

In distributed systems, direct service-to-service communication creates tight coupling, making systems fragile and difficult to evolve. The Pub/Sub pattern decouples senders from receivers, enabling flexible, scalable architectures.

#### Problem Statement

Given:

- Services S = {s₁, s₂, ..., sₙ} that need to communicate
- Events E = {e₁, e₂, ..., eₘ} representing state changes

Replace direct communication with:

```
Direct: sᵢ → sⱼ (point-to-point, tightly coupled)

Pub/Sub: sᵢ → Topic → {sⱼ, sₖ, ...} (loosely coupled)

Benefits:
    - Publishers don't know about subscribers
    - Subscribers don't know about publishers
    - Easy to add new consumers
    - Natural fan-out for broadcasts
```

### Communication Pattern Comparison

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Communication Patterns Comparison                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Request-Response (Tight Coupling):                                     │
│  ┌──────────┐     Request      ┌──────────┐                            │
│  │ Service A│─────────────────►│ Service B│                            │
│  │          │◄─────────────────│          │                            │
│  └──────────┘     Response     └──────────┘                            │
│                                                                         │
│  Problems:                                                              │
│  • Service A must know Service B's location                             │
│  • Service B must be available when Service A calls                     │
│  • Adding Service C requires changes to Service A                       │
│  • Cascading failures                                                   │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Pub/Sub (Loose Coupling):                                              │
│  ┌──────────┐                                                          │
│  │ Service A│──Publish──►┌─────────────┐                               │
│  │(Publisher)│            │   Topic     │                               │
│  └──────────┘            └──────┬──────┘                               │
│                                 │                                       │
│                    ┌────────────┼────────────┐                         │
│                    ▼            ▼            ▼                         │
│               ┌─────────┐  ┌─────────┐  ┌─────────┐                    │
│               │Service B│  │Service C│  │Service D│                    │
│               │(Sub)    │  │(Sub)    │  │(Sub)    │                    │
│               └─────────┘  └─────────┘  └─────────┘                    │
│                                                                         │
│  Benefits:                                                              │
│  • Service A publishes without knowing subscribers                      │
│  • Services B, C, D receive without knowing publisher                   │
│  • Easy to add Service E - just subscribe                               │
│  • Message broker handles availability                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Pub/Sub Topology Patterns

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Pub/Sub Topology Patterns                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Pattern 1: Single Topic                                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  Publishers ──►  orders.topic  ──► Subscribers                  │   │
│  │                                                                  │   │
│  │  Simple but limited filtering options                            │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 2: Hierarchical Topics                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  orders.created                                                 │   │
│  │       ├── orders.created.success                                │   │
│  │       └── orders.created.failed                                 │   │
│  │                                                                  │   │
│  │  orders.updated                                                 │   │
│  │       ├── orders.updated.payment                                │   │
│  │       └── orders.updated.shipping                               │   │
│  │                                                                  │   │
│  │  Subscriptions can use wildcards: orders.*.success              │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 3: Topic with Partitions                                       │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │  ┌─────────────────────────────────────────────────────────────┐│   │
│  │  │  orders.topic (3 partitions)                                ││   │
│  │  │  ┌─────────┬─────────┬─────────┐                           ││   │
│  │  │  │ Part 0  │ Part 1  │ Part 2  │                           ││   │
│  │  │  │ [1,4,7] │ [2,5,8] │ [3,6,9] │                           ││   │
│  │  │  └────┬────┴────┬────┴────┬────┘                           ││   │
│  │  │       │         │         │                                ││   │
│  │  │       ▼         ▼         ▼                                ││   │
│  │  │  Consumer 1  Consumer 2  Consumer 3                        ││   │
│  │  └─────────────────────────────────────────────────────────────┘│   │
│  │                                                                  │   │
│  │  Enables parallel processing while preserving order per key      │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Message Delivery Guarantees

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Delivery Guarantee Levels                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  At-Most-Once (Fire and Forget):                                        │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐                          │
│  │ Publisher│───►│  Broker  │───►│Subscriber│                          │
│  └──────────┘    └──────────┘    └──────────┘                          │
│       │               │               │                                │
│       └───────────────┴───────────────┘                                │
│       No acknowledgment                                                 │
│                                                                         │
│  Pros: Lowest latency                                                   │
│  Cons: Messages can be lost                                             │
│  Use: Telemetry, metrics                                                │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  At-Least-Once (Acknowledged Delivery):                                 │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐                          │
│  │ Publisher│───►│  Broker  │───►│Subscriber│                          │
│  └──────────┘    └──────────┘    └────┬─────┘                          │
│       │               │◄───────────────┘                                │
│       │          ACK received                                           │
│       │◄──────────────┘                                                 │
│  Delivery confirmed                                                     │
│                                                                         │
│  Pros: No message loss                                                  │
│  Cons: Possible duplicates                                              │
│  Use: Most business events                                              │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Exactly-Once (Transactional):                                          │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐                          │
│  │ Publisher│───►│  Broker  │───►│Subscriber│                          │
│  └────┬─────┘    └────┬─────┘    └────┬─────┘                          │
│       │               │◄───────────────┘                                │
│       │          Transaction ID                                         │
│       │◄──────────────┘                                                 │
│  Idempotent processing                                                  │
│                                                                         │
│  Pros: No loss, no duplicates                                           │
│  Cons: Higher latency, complexity                                       │
│  Use: Financial transactions                                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Publisher Implementation

```go
// pkg/pubsub/publisher.go
package pubsub

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/IBM/sarama"
    "go.uber.org/zap"
)

// Publisher publishes messages to topics
type Publisher struct {
    producer sarama.AsyncProducer
    config   *PublisherConfig
    logger   *zap.Logger

    // Metrics
    published    uint64
    publishErrors uint64

    // Shutdown
    closed bool
}

type PublisherConfig struct {
    Brokers         []string
    MaxRetries      int
    RetryBackoff    time.Duration
    RequiredAcks    sarama.RequiredAcks // NoResponse, WaitForLocal, WaitForAll
    Compression     sarama.CompressionCodec
    MaxMessageBytes int
}

// Message represents a pub/sub message
type Message struct {
    Topic       string
    Key         []byte
    Value       []byte
    Headers     map[string]string
    Timestamp   time.Time
    Partition   int32 // Optional: for ordering guarantees
}

// Event is a strongly-typed event for publishing
type Event struct {
    Type      string
    Aggregate string
    ID        string
    Timestamp time.Time
    Payload   interface{}
}

func NewPublisher(cfg *PublisherConfig, logger *zap.Logger) (*Publisher, error) {
    saramaConfig := sarama.NewConfig()
    saramaConfig.Producer.RequiredAcks = cfg.RequiredAcks
    saramaConfig.Producer.Compression = cfg.Compression
    saramaConfig.Producer.MaxMessageBytes = cfg.MaxMessageBytes
    saramaConfig.Producer.Retry.Max = cfg.MaxRetries
    saramaConfig.Producer.Retry.Backoff = cfg.RetryBackoff
    saramaConfig.Producer.Return.Successes = true
    saramaConfig.Producer.Return.Errors = true

    producer, err := sarama.NewAsyncProducer(cfg.Brokers, saramaConfig)
    if err != nil {
        return nil, fmt.Errorf("creating producer: %w", err)
    }

    pub := &Publisher{
        producer: producer,
        config:   cfg,
        logger:   logger,
    }

    // Start error handling goroutine
    go pub.handleErrors()
    go pub.handleSuccesses()

    return pub, nil
}

func (p *Publisher) Publish(ctx context.Context, msg *Message) error {
    if p.closed {
        return fmt.Errorf("publisher is closed")
    }

    headers := make([]sarama.RecordHeader, 0, len(msg.Headers))
    for k, v := range msg.Headers {
        headers = append(headers, sarama.RecordHeader{
            Key:   []byte(k),
            Value: []byte(v),
        })
    }

    kafkaMsg := &sarama.ProducerMessage{
        Topic:     msg.Topic,
        Key:       sarama.ByteEncoder(msg.Key),
        Value:     sarama.ByteEncoder(msg.Value),
        Headers:   headers,
        Timestamp: msg.Timestamp,
        Partition: msg.Partition,
    }

    select {
    case p.producer.Input() <- kafkaMsg:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (p *Publisher) PublishEvent(ctx context.Context, topic string, event *Event) error {
    payload, err := json.Marshal(event.Payload)
    if err != nil {
        return fmt.Errorf("marshaling payload: %w", err)
    }

    headers := map[string]string{
        "event_type":   event.Type,
        "aggregate":    event.Aggregate,
        "event_id":     event.ID,
        "timestamp":    event.Timestamp.Format(time.RFC3339),
    }

    // Use aggregate ID as key for ordering
    key := []byte(event.Aggregate + ":" + event.ID)

    msg := &Message{
        Topic:     topic,
        Key:       key,
        Value:     payload,
        Headers:   headers,
        Timestamp: event.Timestamp,
    }

    return p.Publish(ctx, msg)
}

func (p *Publisher) PublishBatch(ctx context.Context, msgs []*Message) error {
    for _, msg := range msgs {
        if err := p.Publish(ctx, msg); err != nil {
            return err
        }
    }
    return nil
}

func (p *Publisher) handleErrors() {
    for err := range p.producer.Errors() {
        p.publishErrors++
        p.logger.Error("publish error",
            zap.Error(err.Err),
            zap.String("topic", err.Msg.Topic),
            zap.Int32("partition", err.Msg.Partition))
    }
}

func (p *Publisher) handleSuccesses() {
    for range p.producer.Successes() {
        p.published++
    }
}

func (p *Publisher) Close() error {
    p.closed = true
    return p.producer.Close()
}
```

### Subscriber Implementation

```go
// pkg/pubsub/subscriber.go
package pubsub

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/IBM/sarama"
    "go.uber.org/zap"
)

// Subscriber subscribes to topics and processes messages
type Subscriber struct {
    config     *SubscriberConfig
    consumer   sarama.ConsumerGroup
    handler    Handler
    logger     *zap.Logger

    // Lifecycle
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
    ready      chan bool

    // Metrics
    processed  uint64
    failed     uint64
}

type SubscriberConfig struct {
    Brokers           []string
    GroupID           string
    Topics            []string
    InitialOffset     int64 // OffsetNewest or OffsetOldest

    // Processing
    Workers           int
    MaxProcessingTime time.Duration
    RetryLimit        int

    // Session
    SessionTimeout    time.Duration
    HeartbeatInterval time.Duration
}

// Handler processes messages
type Handler interface {
    Handle(ctx context.Context, msg *Message) error
}

type HandlerFunc func(ctx context.Context, msg *Message) error

func (f HandlerFunc) Handle(ctx context.Context, msg *Message) error {
    return f(ctx, msg)
}

func NewSubscriber(cfg *SubscriberConfig, handler Handler, logger *zap.Logger) (*Subscriber, error) {
    ctx, cancel := context.WithCancel(context.Background())

    saramaConfig := sarama.NewConfig()
    saramaConfig.Version = sarama.V2_6_0_0
    saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    saramaConfig.Consumer.Group.Session.Timeout = cfg.SessionTimeout
    saramaConfig.Consumer.Group.Heartbeat.Interval = cfg.HeartbeatInterval
    saramaConfig.Consumer.Offsets.Initial = cfg.InitialOffset
    saramaConfig.Consumer.MaxProcessingTime = cfg.MaxProcessingTime

    consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConfig)
    if err != nil {
        cancel()
        return nil, fmt.Errorf("creating consumer group: %w", err)
    }

    return &Subscriber{
        config:   cfg,
        consumer: consumer,
        handler:  handler,
        logger:   logger,
        ctx:      ctx,
        cancel:   cancel,
        ready:    make(chan bool),
    }, nil
}

func (s *Subscriber) Start() error {
    s.wg.Add(1)
    go s.consumeLoop()

    <-s.ready
    s.logger.Info("subscriber started",
        zap.String("group", s.config.GroupID),
        zap.Strings("topics", s.config.Topics))

    return nil
}

func (s *Subscriber) Stop() error {
    s.cancel()
    s.wg.Wait()
    return s.consumer.Close()
}

func (s *Subscriber) consumeLoop() {
    defer s.wg.Done()

    for {
        select {
        case <-s.ctx.Done():
            return
        default:
        }

        handler := &consumerGroupHandler{
            subscriber: s,
            ready:      s.ready,
        }

        err := s.consumer.Consume(s.ctx, s.config.Topics, handler)
        if err != nil {
            s.logger.Error("consume error", zap.Error(err))
        }

        if s.ctx.Err() != nil {
            return
        }

        s.ready = make(chan bool)
    }
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
    subscriber *Subscriber
    ready      chan bool
}

func (h *consumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
    close(h.ready)
    return nil
}

func (h *consumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
    return nil
}

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    // Create worker pool
    sem := make(chan struct{}, h.subscriber.config.Workers)

    for msg := range claim.Messages() {
        select {
        case sem <- struct{}{}:
        case <-sess.Context().Done():
            return nil
        }

        go func(m *sarama.ConsumerMessage) {
            defer func() { <-sem }()

            h.processMessage(sess, m)
        }(msg)
    }

    return nil
}

func (h *consumerGroupHandler) processMessage(sess sarama.ConsumerGroupSession, kafkaMsg *sarama.ConsumerMessage) {
    headers := make(map[string]string)
    for _, h := range kafkaMsg.Headers {
        headers[string(h.Key)] = string(h.Value)
    }

    msg := &Message{
        Topic:     kafkaMsg.Topic,
        Partition: kafkaMsg.Partition,
        Offset:    kafkaMsg.Offset,
        Key:       kafkaMsg.Key,
        Value:     kafkaMsg.Value,
        Timestamp: kafkaMsg.Timestamp,
        Headers:   headers,
    }

    ctx, cancel := context.WithTimeout(h.subscriber.ctx, h.subscriber.config.MaxProcessingTime)
    defer cancel()

    var err error
    for attempt := 0; attempt <= h.subscriber.config.RetryLimit; attempt++ {
        err = h.subscriber.handler.Handle(ctx, msg)
        if err == nil {
            break
        }

        if attempt < h.subscriber.config.RetryLimit {
            time.Sleep(time.Duration(attempt+1) * time.Second)
        }
    }

    if err != nil {
        h.subscriber.failed++
        h.subscriber.logger.Error("message handling failed",
            zap.Error(err),
            zap.String("topic", msg.Topic),
            zap.Int32("partition", msg.Partition),
            zap.Int64("offset", msg.Offset))

        // Send to dead letter
        h.sendToDeadLetter(msg, err)
    } else {
        h.subscriber.processed++
    }

    sess.MarkMessage(kafkaMsg, "")
}

func (h *consumerGroupHandler) sendToDeadLetter(msg *Message, err error) {
    // Implementation: send to DLQ topic or storage
}
```

### Event Bus with Type-Safe Handlers

```go
// pkg/pubsub/eventbus.go
package pubsub

import (
    "context"
    "fmt"
    "reflect"
    "sync"
)

// EventBus provides type-safe event publishing and subscribing
type EventBus struct {
    publisher  *Publisher
    subscriber *Subscriber

    handlers   map[string][]EventHandler
    mu         sync.RWMutex

    typeMap    map[string]reflect.Type
}

// EventHandler handles typed events
type EventHandler interface {
    HandleEvent(ctx context.Context, event interface{}) error
    EventType() string
}

type typedHandler struct {
    handlerFunc func(ctx context.Context, event interface{}) error
    eventType   reflect.Type
}

func (h *typedHandler) HandleEvent(ctx context.Context, event interface{}) error {
    return h.handlerFunc(ctx, event)
}

func (h *typedHandler) EventType() string {
    return h.eventType.String()
}

func NewEventBus(publisher *Publisher, subscriber *Subscriber) *EventBus {
    return &EventBus{
        publisher:  publisher,
        subscriber: subscriber,
        handlers:   make(map[string][]EventHandler),
        typeMap:    make(map[string]reflect.Type),
    }
}

// Subscribe registers a handler for an event type
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()

    eb.handlers[eventType] = append(eb.handlers[eventType], handler)
    eb.typeMap[eventType] = reflect.TypeOf(handler).Elem()
}

// SubscribeFunc registers a function handler for an event type
func SubscribeFunc[T any](eb *EventBus, handler func(ctx context.Context, event T) error) {
    var zero T
    eventType := reflect.TypeOf(zero).String()

    wrapped := func(ctx context.Context, evt interface{}) error {
        typed, ok := evt.(T)
        if !ok {
            return fmt.Errorf("event type mismatch")
        }
        return handler(ctx, typed)
    }

    eb.Subscribe(eventType, &typedHandler{
        handlerFunc: wrapped,
        eventType:   reflect.TypeOf(zero),
    })
}

// Publish publishes a typed event
func Publish[T any](ctx context.Context, eb *EventBus, topic string, event T) error {
    payload, err := json.Marshal(event)
    if err != nil {
        return err
    }

    var zero T
    eventType := reflect.TypeOf(zero).String()

    msg := &Message{
        Topic: topic,
        Value: payload,
        Headers: map[string]string{
            "event_type": eventType,
        },
    }

    return eb.publisher.Publish(ctx, msg)
}

// Start begins processing events
func (eb *EventBus) Start() error {
    handler := HandlerFunc(func(ctx context.Context, msg *Message) error {
        eventType := msg.Headers["event_type"]
        if eventType == "" {
            return fmt.Errorf("no event type in message")
        }

        eb.mu.RLock()
        handlers := eb.handlers[eventType]
        eb.mu.RUnlock()

        if len(handlers) == 0 {
            return fmt.Errorf("no handlers for event type: %s", eventType)
        }

        // Deserialize
        eventTypeRef := eb.typeMap[eventType]
        event := reflect.New(eventTypeRef).Interface()
        if err := json.Unmarshal(msg.Value, event); err != nil {
            return err
        }

        // Call handlers
        for _, h := range handlers {
            if err := h.HandleEvent(ctx, event); err != nil {
                return err
            }
        }

        return nil
    })

    eb.subscriber.handler = handler
    return eb.subscriber.Start()
}
```

## Trade-off Analysis

### Pub/Sub vs Direct Communication

| Aspect | Pub/Sub | Direct (HTTP/gRPC) | Notes |
|--------|---------|-------------------|-------|
| **Coupling** | Loose | Tight | Pub/Sub decouples in time and space |
| **Latency** | Higher (ms) | Lower (μs-ms) | Pub/Sub adds broker hop |
| **Reliability** | Higher | Lower | Broker persists messages |
| **Ordering** | Per partition | Per connection | Both can guarantee ordering |
| **Observability** | Centralized | Distributed | Broker provides unified view |
| **Cost** | Higher (infrastructure) | Lower | Broker adds cost |

### Message Size and Frequency Guidelines

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Pub/Sub Usage Guidelines                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Message Size:                                                          │
│  • 0-1KB:     Ideal for pub/sub                                         │
│  • 1KB-100KB: Acceptable                                                │
│  • 100KB-1MB: Consider compression or claim-check pattern               │
│  • >1MB:      Use claim-check pattern (store payload, pass reference)   │
│                                                                         │
│  Message Frequency:                                                     │
│  • < 1K/sec:  Any broker works well                                     │
│  • 1K-100K/sec: Kafka, Pulsar recommended                               │
│  • > 100K/sec: Partition carefully, consider batching                   │
│                                                                         │
│  When NOT to use Pub/Sub:                                               │
│  • Synchronous request-response                                         │
│  • Real-time gaming (<10ms latency required)                            │
│  • Simple point-to-point with no future fan-out                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Pub/Sub Testing

```go
// test/pubsub/pubsub_test.go
package pubsub

import (
    "context"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPublisherSubscriber(t *testing.T) {
    // Setup
    pubCfg := &PublisherConfig{
        Brokers:      []string{"localhost:9092"},
        RequiredAcks: sarama.WaitForLocal,
    }

    subCfg := &SubscriberConfig{
        Brokers:       []string{"localhost:9092"},
        GroupID:       "test-group",
        Topics:        []string{"test-topic"},
        Workers:       3,
        InitialOffset: sarama.OffsetNewest,
    }

    publisher, err := NewPublisher(pubCfg, zap.NewNop())
    require.NoError(t, err)
    defer publisher.Close()

    received := make(chan *Message, 100)
    handler := HandlerFunc(func(ctx context.Context, msg *Message) error {
        received <- msg
        return nil
    })

    subscriber, err := NewSubscriber(subCfg, handler, zap.NewNop())
    require.NoError(t, err)

    err = subscriber.Start()
    require.NoError(t, err)
    defer subscriber.Stop()

    // Publish message
    ctx := context.Background()
    msg := &Message{
        Topic: "test-topic",
        Key:   []byte("key-1"),
        Value: []byte(`{"data":"test"}`),
    }

    err = publisher.Publish(ctx, msg)
    require.NoError(t, err)

    // Verify receipt
    select {
    case receivedMsg := <-received:
        assert.Equal(t, "test-topic", receivedMsg.Topic)
        assert.Equal(t, "key-1", string(receivedMsg.Key))
    case <-time.After(10 * time.Second):
        t.Fatal("timeout waiting for message")
    }
}

func TestEventBus(t *testing.T) {
    type OrderCreated struct {
        OrderID string
        Amount  float64
    }

    received := make(chan OrderCreated, 10)

    bus := NewEventBus(nil, nil) // Mock for test

    SubscribeFunc(bus, func(ctx context.Context, evt OrderCreated) error {
        received <- evt
        return nil
    })

    // Simulate publish
    event := OrderCreated{OrderID: "123", Amount: 99.99}

    handlers := bus.handlers["pubsub.OrderCreated"]
    require.Len(t, handlers, 1)

    err := handlers[0].HandleEvent(context.Background(), event)
    require.NoError(t, err)

    select {
    case receivedEvent := <-received:
        assert.Equal(t, "123", receivedEvent.OrderID)
        assert.Equal(t, 99.99, receivedEvent.Amount)
    case <-time.After(time.Second):
        t.Fatal("timeout")
    }
}

func TestAtLeastOnceDelivery(t *testing.T) {
    // Test that messages are redelivered on failure

    attemptCount := 0
    var mu sync.Mutex

    handler := HandlerFunc(func(ctx context.Context, msg *Message) error {
        mu.Lock()
        attemptCount++
        mu.Unlock()
        return assert.AnError // Always fail
    })

    cfg := &SubscriberConfig{
        Brokers:    []string{"localhost:9092"},
        GroupID:    "retry-test",
        Topics:     []string{"retry-topic"},
        RetryLimit: 2,
    }

    subscriber, _ := NewSubscriber(cfg, handler, zap.NewNop())
    subscriber.Start()
    defer subscriber.Stop()

    // Publish and wait
    time.Sleep(5 * time.Second)

    mu.Lock()
    assert.GreaterOrEqual(t, attemptCount, 2) // At least retry limit attempts
    mu.Unlock()
}
```

## Summary

The Publisher-Subscriber Pattern provides:

1. **Decoupling**: Publishers and subscribers are independent
2. **Scalability**: Easy to add consumers
3. **Flexibility**: Dynamic subscription changes
4. **Reliability**: Persistent messaging
5. **Observability**: Central message flow visibility

Key considerations:

- Choose appropriate delivery guarantees
- Design topics for proper partitioning
- Handle message ordering requirements
- Monitor consumer lag
- Plan for dead letter queues
