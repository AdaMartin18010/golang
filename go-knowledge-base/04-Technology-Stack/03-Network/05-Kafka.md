# TS-NET-005: Apache Kafka Architecture and Go Integration

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #kafka #streaming #messaging #distributed #event-driven
> **权威来源**:
>
> - [Apache Kafka Documentation](https://kafka.apache.org/documentation/) - Apache
> - [Kafka: The Definitive Guide](https://www.oreilly.com/library/view/kafka-the-definitive/9781491936153/) - O'Reilly
> - [Sarama (Go Client)](https://github.com/Shopify/sarama) - Shopify
> - [franz-go](https://github.com/twmb/franz-go) - Modern Go client

---

## 1. Kafka Architecture Overview

### 1.1 Distributed Log Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Apache Kafka Distributed Architecture                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Kafka Cluster                                 │   │
│  │  ┌─────────────────┬─────────────────┬─────────────────┐           │   │
│  │  │   Broker 1      │   Broker 2      │   Broker 3      │           │   │
│  │  │   (Leader)      │   (Follower)    │   (Follower)    │           │   │
│  │  │                 │                 │                 │           │   │
│  │  │  ┌───────────┐  │  ┌───────────┐  │  ┌───────────┐  │           │   │
│  │  │  │ Partition │  │  │ Partition │  │  │ Partition │  │           │   │
│  │  │  │ 0 (Leader)│  │  │ 0 (Replica)│  │  │ 0 (Replica)│  │           │   │
│  │  │  ├───────────┤  │  ├───────────┤  │  ├───────────┤  │           │   │
│  │  │  │ Partition │  │  │ Partition │  │  │ Partition │  │           │   │
│  │  │  │ 1 (Replica)│  │  │ 1 (Leader)│  │  │ 1 (Replica)│  │           │   │
│  │  │  ├───────────┤  │  ├───────────┤  │  ├───────────┤  │           │   │
│  │  │  │ Partition │  │  │ Partition │  │  │ Partition │  │           │   │
│  │  │  │ 2 (Replica)│  │  │ 2 (Replica)│  │  │ 2 (Leader)│  │           │   │
│  │  │  └───────────┘  │  └───────────┘  │  └───────────┘  │           │   │
│  │  └─────────────────┴─────────────────┴─────────────────┘           │   │
│  │                              ▲                                     │   │
│  └──────────────────────────────┼─────────────────────────────────────┘   │
│                                 │                                          │
│  ┌──────────────────────────────┼─────────────────────────────────────┐   │
│  │                              │         ZooKeeper / KRaft            │   │
│  │  ┌───────────────────┐       │  ┌───────────────────────────────┐   │   │
│  │  │   Producers       │───────┘  │  - Controller election        │   │   │
│  │  │                   │          │  - Cluster membership         │   │   │
│  │  │  ┌─────────────┐  │          │  - Topic configuration        │   │   │
│  │  │  │ Partitioner │  │          │  - ISR management             │   │   │
│  │  │  │ (hash/rr)   │  │          │  - ACL storage                │   │   │
│  │  │  └──────┬──────┘  │          └───────────────────────────────┘   │   │
│  │  │         │         │                                             │   │
│  │  └─────────┼─────────┘                                             │   │
│  │            │                                                        │   │
│  │  ┌─────────┴─────────┐                                             │   │
│  │  │   Consumers       │                                             │   │
│  │  │   (Consumer Group)│                                             │   │
│  │  │                   │                                             │   │
│  │  │  ┌─────────────┐  │                                             │   │
│  │  │  │ Group       │  │                                             │   │
│  │  │  │ Coordinator │  │                                             │   │
│  │  │  │ (on broker) │  │                                             │   │
│  │  │  └─────────────┘  │                                             │   │
│  │  └───────────────────┘                                             │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Partition Log Structure

```
Topic: orders
Partitions: 3
Replication Factor: 3

Partition 0 (Leader on Broker 1):
┌────────────────────────────────────────────────────────────────┐
│ Offset  │ Timestamp  │ Key         │ Value (Message)           │
├─────────┼────────────┼─────────────┼───────────────────────────┤
│ 0       │ T1         │ user-123    │ {"order": "A1", "amt": 100}│
│ 1       │ T2         │ user-456    │ {"order": "A2", "amt": 200}│
│ 2       │ T3         │ user-123    │ {"order": "A3", "amt": 150}│
│ ...     │ ...        │ ...         │ ...                       │
└─────────┴────────────┴─────────────┴───────────────────────────┘
         ▲
         │
    Consumer Offset: 1 (next read: offset 2)

Partition 1 (Leader on Broker 2):
┌────────────────────────────────────────────────────────────────┐
│ Offset  │ Timestamp  │ Key         │ Value                     │
├─────────┼────────────┼─────────────┼───────────────────────────┤
│ 0       │ T1         │ user-789    │ {"order": "B1", "amt": 300}│
│ 1       │ T4         │ user-123    │ {"order": "B2", "amt": 400}│
│ ...     │ ...        │ ...         │ ...                       │
└─────────┴────────────┴─────────────┴───────────────────────────┘

Physical Storage (Segment Files):
orders-0/
├── 00000000000000000000.log    (offset 0-999)
├── 00000000000000000000.index  (sparse index)
├── 00000000000000001000.log    (offset 1000-1999)
├── 00000000000000001000.index
└── ...

Log Retention:
- time-based: 7 days (default)
- size-based: 1GB per partition
- compaction: keep latest value per key
```

---

## 2. Producer Architecture

### 2.1 Producer Message Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Producer Message Flow                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Application                                                                  │
│     │                                                                         │
│     ▼                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Producer API (sarama/franz-go)                   │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │   │
│  │  │ Serializer  │  │ Partitioner │  │  Compressor │                 │   │
│  │  │ (JSON/Avro) │──►│ (hash/key)  │──►│ (gzip/snappy)│              │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                 │   │
│  │                                                                      │   │
│  │  ┌────────────────────────────────────────────────────────────────┐ │   │
│  │  │                      RecordAccumulator                        │ │   │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │ │   │
│  │  │  │ Batch 0      │  │ Batch 1      │  │ Batch 2      │        │ │   │
│  │  │  │ (Partition 0)│  │ (Partition 1)│  │ (Partition 2)│        │ │   │
│  │  │  │ - Message 1  │  │ - Message 3  │  │ - Message 2  │        │ │   │
│  │  │  │ - Message 4  │  │              │  │              │        │ │   │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘        │ │   │
│  │  └────────────────────────────────────────────────────────────────┘ │   │
│  │                            │                                        │   │
│  │  ┌─────────────────────────┴─────────────────────────────────┐    │   │
│  │  │                    Sender (background thread)               │    │   │
│  │  │  - batch.size: 16384 bytes (default)                      │    │   │
│  │  │  - linger.ms: 0 (send immediately) or batch               │    │   │
│  │  │  - max.in.flight.requests: 5 (pipelining)                 │    │   │
│  │  └───────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│                                    ▼                                         │
│                           ┌─────────────────┐                                │
│                           │  Broker (Leader)│                                │
│                           └─────────────────┘                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Go Producer Implementation

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/IBM/sarama"
)

type OrderEvent struct {
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
}

func createProducer() (sarama.SyncProducer, error) {
    config := sarama.NewConfig()

    // Producer configurations
    config.Producer.RequiredAcks = sarama.WaitForAll  // -1: Wait for all replicas
    config.Producer.Retry.Max = 5                      // Retry up to 5 times
    config.Producer.Retry.Backoff = 100 * time.Millisecond
    config.Producer.Return.Successes = true           // Return success channel
    config.Producer.Return.Errors = true              // Return error channel

    // Idempotence (exactly-once semantics)
    config.Producer.Idempotent = true
    config.Net.MaxOpenRequests = 1

    // Compression
    config.Producer.Compression = sarama.CompressionSnappy
    config.Producer.CompressionLevel = 9

    // Batching
    config.Producer.Flush.Bytes = 16384    // 16KB batches
    config.Producer.Flush.Messages = 100   // Or 100 messages
    config.Producer.Flush.Frequency = 100 * time.Millisecond
    config.Producer.MaxMessageBytes = 1000000  // 1MB max message

    brokers := []string{"localhost:9092", "localhost:9093"}
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    return producer, nil
}

func createAsyncProducer() (sarama.AsyncProducer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true

    brokers := []string{"localhost:9092"}
    return sarama.NewAsyncProducer(brokers, config)
}

// Synchronous produce
func produceMessage(producer sarama.SyncProducer, event OrderEvent) error {
    value, err := json.Marshal(event)
    if err != nil {
        return err
    }

    msg := &sarama.ProducerMessage{
        Topic: "orders",
        Key:   sarama.StringEncoder(event.UserID), // Same user -> same partition
        Value: sarama.ByteEncoder(value),
        Headers: []sarama.RecordHeader{
            {Key: []byte("source"), Value: []byte("order-service")},
            {Key: []byte("version"), Value: []byte("1.0")},
        },
    }

    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        return err
    }

    log.Printf("Message sent to partition %d at offset %d", partition, offset)
    return nil
}

// Asynchronous produce with goroutine
func produceAsync(asyncProducer sarama.AsyncProducer, events <-chan OrderEvent) {
    go func() {
        for err := range asyncProducer.Errors() {
            log.Printf("Failed to produce message: %v", err)
        }
    }()

    go func() {
        for success := range asyncProducer.Successes() {
            log.Printf("Successfully produced: partition=%d, offset=%d",
                success.Partition, success.Offset)
        }
    }()

    for event := range events {
        value, _ := json.Marshal(event)
        asyncProducer.Input() <- &sarama.ProducerMessage{
            Topic: "orders",
            Key:   sarama.StringEncoder(event.UserID),
            Value: sarama.ByteEncoder(value),
        }
    }
}

func main() {
    producer, err := createProducer()
    if err != nil {
        log.Fatal(err)
    }
    defer producer.Close()

    event := OrderEvent{
        OrderID:   "ORD-001",
        UserID:    "USER-123",
        Amount:    99.99,
        Status:    "CREATED",
        Timestamp: time.Now(),
    }

    if err := produceMessage(producer, event); err != nil {
        log.Printf("Failed to produce: %v", err)
    }
}
```

---

## 3. Consumer Architecture

### 3.1 Consumer Group Rebalancing

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Consumer Group Rebalancing                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Scenario 1: Initial Join                                                   │
│  ┌─────────────────────┐                                                    │
│  │  Topic: orders (4 partitions)                                            │
│  │  P0, P1, P2, P3                                                          │
│  └─────────────────────┘                                                    │
│           │                                                                  │
│           │   Join Group                                                     │
│           ▼                                                                  │
│  ┌─────────┐  ┌─────────┐                                                   │
│  │Consumer │  │Consumer │                                                   │
│  │    A    │  │    B    │                                                   │
│  └────┬────┘  └────┬────┘                                                   │
│       │            │                                                         │
│       │  ┌─────────┴─────────┐                                               │
│       └──►│ Group Coordinator │                                               │
│           │ (on broker)       │                                               │
│           └─────────┬─────────┘                                               │
│                     │                                                        │
│                     ▼                                                        │
│  ┌─────────────────────────────────────────┐                                 │
│  │ Assignment:                              │                                 │
│  │  Consumer A: P0, P1                      │                                 │
│  │  Consumer B: P2, P3                      │                                 │
│  └─────────────────────────────────────────┘                                 │
│                                                                              │
│  Scenario 2: Consumer C Joins                                               │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                                      │
│  │Consumer │  │Consumer │  │Consumer │                                      │
│  │    A    │  │    B    │  │    C    │                                      │
│  └─────────┘  └─────────┘  └─────────┘                                      │
│                                                                              │
│  Rebalance triggered → New assignment:                                      │
│  Consumer A: P0                                                            │
│  Consumer B: P1                                                            │
│  Consumer C: P2, P3                                                        │
│                                                                              │
│  Scenario 3: Consumer B Leaves                                              │
│  ┌─────────┐              ┌─────────┐                                       │
│  │Consumer │              │Consumer │                                       │
│  │    A    │              │    C    │                                       │
│  └─────────┘              └─────────┘                                       │
│                                                                              │
│  Rebalance triggered → New assignment:                                      │
│  Consumer A: P0, P1                                                        │
│  Consumer C: P2, P3                                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Go Consumer Implementation

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
    processor MessageProcessor
}

func (h ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
    log.Println("Consumer group session starting")
    return nil
}

func (h ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
    log.Println("Consumer group session ending")
    return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for {
        select {
        case message := <-claim.Messages():
            if message == nil {
                return nil
            }

            // Process message
            if err := h.processMessage(message); err != nil {
                log.Printf("Failed to process message: %v", err)
                // Don't commit - will be reprocessed
                continue
            }

            // Mark message as processed (commit offset)
            session.MarkMessage(message, "")

        case <-session.Context().Done():
            return nil
        }
    }
}

func (h ConsumerGroupHandler) processMessage(msg *sarama.ConsumerMessage) error {
    var event OrderEvent
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        return err
    }

    log.Printf("Processing: topic=%s, partition=%d, offset=%d, key=%s",
        msg.Topic, msg.Partition, msg.Offset, string(msg.Key))

    // Business logic here
    return h.processor.Process(event)
}

func createConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
    config := sarama.NewConfig()

    // Consumer configurations
    config.Version = sarama.V2_8_0_0
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    config.Consumer.Offsets.AutoCommit.Enable = false // Manual commit

    // Fetch configurations
    config.Consumer.Fetch.Min = 1           // 1 byte minimum
    config.Consumer.Fetch.Default = 1024 * 1024  // 1MB default
    config.Consumer.Fetch.Max = 10 * 1024 * 1024 // 10MB max
    config.Consumer.MaxWaitTime = 500 * time.Millisecond
    config.Consumer.MaxProcessingTime = 1 * time.Second

    // Retry configurations
    config.Consumer.Retry.Backoff = 2 * time.Second

    brokers := []string{"localhost:9092"}
    return sarama.NewConsumerGroup(brokers, groupID, config)
}

func main() {
    groupID := "order-processors"
    topics := []string{"orders"}

    consumerGroup, err := createConsumerGroup(groupID)
    if err != nil {
        log.Fatal(err)
    }
    defer consumerGroup.Close()

    handler := ConsumerGroupHandler{
        processor: NewOrderProcessor(),
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Graceful shutdown
    sigterm := make(chan os.Signal, 1)
    signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigterm
        log.Println("Shutdown signal received")
        cancel()
    }()

    // Consume
    for {
        err := consumerGroup.Consume(ctx, topics, handler)
        if err != nil {
            log.Printf("Error from consumer: %v", err)
        }

        if ctx.Err() != nil {
            break
        }
    }

    log.Println("Consumer stopped")
}
```

---

## 4. Configuration Best Practices

```
Producer Configuration:
□ acks=all (durability vs latency trade-off)
□ retries=MAX_INT (infinite retries)
□ enable.idempotence=true (exactly-once)
□ compression.type=snappy (good balance)
□ batch.size=16384 (16KB batches)
□ linger.ms=5 (small delay for batching)
□ max.in.flight.requests=5 (pipelining)

Consumer Configuration:
□ auto.offset.reset=earliest (for new groups)
□ enable.auto.commit=false (manual commit)
□ max.poll.records=500 (batch size)
□ session.timeout.ms=10000
□ heartbeat.interval.ms=3000
□ max.poll.interval.ms=300000
```

---

## 5. Checklist

```
Kafka Production Checklist:
□ Proper topic partitioning strategy
□ Configure replication factor >= 3
□ Monitor consumer lag
□ Implement dead letter queue
□ Configure retention policies
□ Enable monitoring (JMX metrics)
□ Set up alerting for under-replicated partitions
□ Test failover scenarios
□ Implement idempotent consumers
□ Configure appropriate batch sizes
□ Monitor disk usage
□ Plan capacity for growth
```
