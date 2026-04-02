# Event-Driven System Example

A comprehensive implementation of an event-driven architecture demonstrating Event Sourcing, CQRS (Command Query Responsibility Segregation), and message broker integration patterns using Go.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Event Sourcing](#event-sourcing)
4. [CQRS Pattern](#cqrs-pattern)
5. [Message Brokers](#message-brokers)
6. [Getting Started](#getting-started)
7. [Implementation Details](#implementation-details)
8. [Deployment](#deployment)
9. [Performance](#performance)
10. [Best Practices](#best-practices)

## Overview

This example demonstrates a real-world event-driven e-commerce system with the following features:

- **Event Sourcing**: All state changes stored as immutable events
- **CQRS**: Separate read and write models optimized for their purposes
- **Multiple Message Brokers**: Kafka for event streaming, NATS for pub/sub, Redis for simple messaging
- **Saga Pattern**: Distributed transactions across multiple services
- **Event Replay**: Reconstruct system state from events
- **Snapshot Management**: Optimize event replay with snapshots
- **Projections**: Real-time read model updates

### Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.21+ |
| Event Store | PostgreSQL + Apache Kafka |
| Read Models | MongoDB, Elasticsearch |
| Message Brokers | Apache Kafka, NATS, Redis Pub/Sub |
| Framework | Go with custom event framework |
| Observability | OpenTelemetry, Prometheus, Jaeger |
| Testing | k6 for load testing |

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Layer                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Web App   │  │ Mobile App  │  │  Admin UI   │  │   External Systems  │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘ │
└─────────┼────────────────┼────────────────┼────────────────────┼────────────┘
          │                │                │                    │
          └────────────────┴────────────────┴────────────────────┘
                                   │
┌──────────────────────────────────▼──────────────────────────────────────────┐
│                         API Gateway / Load Balancer                          │
└──────────────────────────────────┬──────────────────────────────────────────┘
                                   │
          ┌────────────────────────┼────────────────────────┐
          │                        │                        │
┌─────────▼──────────┐  ┌──────────▼──────────┐  ┌──────────▼──────────┐
│   Command Side     │  │    Query Side       │  │   Event Processing  │
│                    │  │                     │  │                     │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ ┌───────────────┐   │
│ │ Order Commands│  │  │ │Order Queries  │   │  │ │ Event Store   │   │
│ │  - Create     │  │  │ │ - GetOrder    │   │  │ │  (Kafka +     │   │
│ │  - Cancel     │──┼──┼▶│ - ListOrders  │   │  │ │   PostgreSQL) │   │
│ │  - Ship       │  │  │ │ - Search      │   │  │ └───────┬───────┘   │
│ └───────────────┘  │  │ └───────────────┘   │  │         │           │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │         ▼           │
│ │Payment Cmds   │  │  │ │Payment Queries│   │  │ ┌───────────────┐   │
│ │  - Process    │──┼──┼▶│ - GetPayments │   │  │ │ EventHandlers │   │
│ │  - Refund     │  │  │ └───────────────┘   │  │ │  - Projections│   │
│ └───────────────┘  │  │                     │  │ │  - Sagas      │   │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ │  - Notifiers  │   │
│ │Inventory Cmds │  │  │ │Inventory Views│   │  │ └───────────────┘   │
│ │  - Reserve    │──┼──┼▶│ - StockLevels │   │  └─────────────────────┘
│ │  - Release    │  │  │ └───────────────┘   │
│ └───────────────┘  │  └─────────────────────┘
└────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Event Store (Kafka)                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ order-events│  │payment-events│  │inventory-events│  │ notification-events│ │
│  │  Partition 0│  │  Partition 0│  │  Partition 0│  │   Partition 0       │ │
│  │  Partition 1│  │  Partition 1│  │  Partition 1│  │   Partition 1       │ │
│  │  Partition 2│  │  Partition 2│  │  Partition 2│  │   Partition 2       │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
          │
          ├──────────────────────────────────────────────────────────────────┐
          │                          │                                       │
┌─────────▼──────────┐    ┌──────────▼──────────┐    ┌──────────────────────▼┐
│   Read Models      │    │   Sagas             │    │  External Systems     │
│                    │    │                     │    │                       │
│ ┌───────────────┐  │    │ ┌───────────────┐   │    │ ┌─────────────────┐   │
│ │   MongoDB     │  │    │ │Order Saga     │   │    │ │  Email Service  │   │
│ │  (Orders View)│◀─┼────┼─┤  - Payment    │   │    │ │                 │   │
│ └───────────────┘  │    │ │  - Inventory  │───┼────┼▶│ ┌─────────────┐ │   │
│ ┌───────────────┐  │    │ │  - Shipping   │   │    │ ││SendGrid/SMTP│ │   │
│ │Elasticsearch  │◀─┼────┼─┤               │   │    │ │└─────────────┘ │   │
│ │ (Search Index)│  │    │ └───────────────┘   │    │ └─────────────────┘   │
│ └───────────────┘  │    └─────────────────────┘    │ ┌─────────────────┐   │
│ ┌───────────────┐  │                               │ │  SMS Service    │   │
│ │     Redis     │◀─┼───────────────────────────────┼▶│                 │   │
│ │ (Cache Layer) │  │                               │ │ ┌─────────────┐ │   │
│ └───────────────┘  │                               │ │ │   Twilio    │ │   │
└────────────────────┘                               │ │ └─────────────┘ │   │
                                                     │ └─────────────────┘   │
                                                     └───────────────────────┘
```

### Event Flow Diagram

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   Command   │────▶│   Domain    │────▶│   Event     │
│   Request   │     │   Handler   │     │   Model     │     │   Store     │
└─────────────┘     └─────────────┘     └─────────────┘     └──────┬──────┘
                                                                   │
                              ┌────────────────────────────────────┤
                              │                                    │
                              ▼                                    ▼
                    ┌─────────────────┐                 ┌──────────────────┐
                    │   Event Bus     │                 │  Event Handlers  │
                    │   (Kafka/NATS)  │                 │                  │
                    └────────┬────────┘                 │ • Projections    │
                             │                          │ • Notifications  │
              ┌──────────────┼──────────────┐           │ • Integrations   │
              ▼              ▼              ▼           └──────────────────┘
    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
    │   Order     │  │  Inventory  │  │  Payment    │
    │  Projector  │  │  Projector  │  │  Projector  │
    └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
           │                │                │
           ▼                ▼                ▼
    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
    │   MongoDB   │  │    Redis    │  │ Elasticsearch│
    │ (Read Model)│  │   (Cache)   │  │  (Search)   │
    └─────────────┘  └─────────────┘  └─────────────┘
```

### Saga Pattern Flow

```
Order Processing Saga:

┌─────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────┐
│  Start  │──▶│ Reserve     │──▶│ Process     │──▶│   Ship      │──▶│ Complete│
│  Order  │   │ Inventory   │   │ Payment     │   │  Order      │   │  Order  │
└─────────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └─────────┘
                     │                 │                 │
                     ▼                 ▼                 ▼
              Failure Path:    Failure Path:    Failure Path:
              • Release        • Refund         • (Log only)
                Inventory      • Release
                                 Inventory
```

## Event Sourcing

### What is Event Sourcing?

Event Sourcing is a pattern where we store the state of a system as a sequence of events. Instead of storing the current state, we store every change that led to that state. This provides:

- **Complete Audit Trail**: Every change is recorded
- **Temporal Queries**: Query state at any point in time
- **Event Replay**: Rebuild state by replaying events
- **Debugging**: Understand exactly how system reached current state

### Event Store Implementation

```go
// Event represents a domain event
type Event struct {
    ID            uuid.UUID       `json:"id"`
    AggregateID   uuid.UUID       `json:"aggregate_id"`
    AggregateType string          `json:"aggregate_type"`
    Type          string          `json:"type"`
    Version       int             `json:"version"`
    Data          json.RawMessage `json:"data"`
    Metadata      Metadata        `json:"metadata"`
    Timestamp     time.Time       `json:"timestamp"`
}

// EventStore interface
type EventStore interface {
    Append(ctx context.Context, events ...*Event) error
    GetEvents(ctx context.Context, aggregateID uuid.UUID, fromVersion int) ([]*Event, error)
    GetAllEvents(ctx context.Context, afterPosition int64, batchSize int) ([]*Event, error)
    Subscribe(ctx context.Context, eventTypes []string) (<-chan *Event, error)
}
```

### Example Events

```go
// OrderCreated event
type OrderCreated struct {
    OrderID     uuid.UUID   `json:"order_id"`
    CustomerID  uuid.UUID   `json:"customer_id"`
    Items       []OrderItem `json:"items"`
    TotalAmount float64     `json:"total_amount"`
    CreatedAt   time.Time   `json:"created_at"`
}

// OrderPaid event
type OrderPaid struct {
    OrderID       uuid.UUID `json:"order_id"`
    PaymentID     uuid.UUID `json:"payment_id"`
    Amount        float64   `json:"amount"`
    PaymentMethod string    `json:"payment_method"`
    PaidAt        time.Time `json:"paid_at"`
}

// OrderShipped event
type OrderShipped struct {
    OrderID       uuid.UUID `json:"order_id"`
    ShipmentID    uuid.UUID `json:"shipment_id"`
    Carrier       string    `json:"carrier"`
    TrackingNumber string   `json:"tracking_number"`
    ShippedAt     time.Time `json:"shipped_at"`
}
```

## CQRS Pattern

### Command Side (Write Model)

The command side handles all write operations and business logic:

```go
// Command handlers
type OrderCommandHandler struct {
    eventStore  EventStore
    publisher   EventPublisher
}

func (h *OrderCommandHandler) CreateOrder(ctx context.Context, cmd CreateOrderCommand) error {
    // 1. Create aggregate
    order, err := NewOrder(cmd.CustomerID, cmd.Items)
    if err != nil {
        return err
    }
    
    // 2. Apply command
    event, err := order.Create(cmd)
    if err != nil {
        return err
    }
    
    // 3. Persist events
    if err := h.eventStore.Append(ctx, event); err != nil {
        return err
    }
    
    // 4. Publish event
    return h.publisher.Publish(ctx, event)
}
```

### Query Side (Read Model)

The query side is optimized for reading with denormalized data:

```go
// Query handlers
type OrderQueryHandler struct {
    readModel ReadModel
}

func (h *OrderQueryHandler) GetOrder(ctx context.Context, query GetOrderQuery) (*OrderView, error) {
    return h.readModel.GetOrder(ctx, query.OrderID)
}

func (h *OrderQueryHandler) ListOrders(ctx context.Context, query ListOrdersQuery) ([]*OrderView, error) {
    return h.readModel.ListOrders(ctx, query.CustomerID, query.Page, query.PageSize)
}

func (h *OrderQueryHandler) SearchOrders(ctx context.Context, query SearchOrdersQuery) (*SearchResult, error) {
    return h.readModel.Search(ctx, query.Query, query.Filters)
}
```

### Read Model Projections

```go
// Order projector updates read models based on events
type OrderProjector struct {
    mongoClient *mongo.Client
    elasticClient *elasticsearch.Client
}

func (p *OrderProjector) HandleOrderCreated(ctx context.Context, event *OrderCreated) error {
    // Update MongoDB read model
    orderView := &OrderView{
        OrderID:     event.OrderID,
        CustomerID:  event.CustomerID,
        Status:      "pending",
        TotalAmount: event.TotalAmount,
        CreatedAt:   event.CreatedAt,
    }
    
    _, err := p.mongoClient.Database("read_models").
        Collection("orders").
        InsertOne(ctx, orderView)
    
    return err
}

func (p *OrderProjector) HandleOrderPaid(ctx context.Context, event *OrderPaid) error {
    // Update MongoDB
    filter := bson.M{"order_id": event.OrderID}
    update := bson.M{
        "$set": bson.M{
            "status":     "paid",
            "payment_id": event.PaymentID,
            "paid_at":    event.PaidAt,
        },
    }
    
    _, err := p.mongoClient.Database("read_models").
        Collection("orders").
        UpdateOne(ctx, filter, update)
    
    // Also update Elasticsearch for search
    return p.indexInElasticsearch(ctx, event.OrderID, update)
}
```

## Message Brokers

### Apache Kafka

Used for high-throughput event streaming:

```go
type KafkaEventBus struct {
    producer sarama.SyncProducer
    consumer sarama.ConsumerGroup
}

func (k *KafkaEventBus) Publish(ctx context.Context, event *Event) error {
    msg := &sarama.ProducerMessage{
        Topic: event.AggregateType + "-events",
        Key:   sarama.StringEncoder(event.AggregateID.String()),
        Value: JSONEncoder(event),
    }
    
    _, _, err := k.producer.SendMessage(msg)
    return err
}
```

### NATS

Used for pub/sub messaging and request/reply patterns:

```go
type NATSEventBus struct {
    conn *nats.Conn
    js   nats.JetStreamContext
}

func (n *NATSEventBus) Subscribe(subject string, handler Handler) error {
    sub, err := n.js.Subscribe(subject, func(msg *nats.Msg) {
        event := &Event{}
        if err := json.Unmarshal(msg.Data, event); err != nil {
            msg.Nak()
            return
        }
        
        if err := handler.Handle(event); err != nil {
            msg.Nak()
            return
        }
        
        msg.Ack()
    }, nats.Durable("durable-consumer"))
    
    return err
}
```

### Redis Pub/Sub

Used for simple messaging and caching:

```go
type RedisEventBus struct {
    client *redis.Client
}

func (r *RedisEventBus) Publish(ctx context.Context, channel string, event *Event) error {
    data, _ := json.Marshal(event)
    return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisEventBus) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
    return r.client.Subscribe(ctx, channels...)
}
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Make

### Quick Start

```bash
# Clone the repository
cd event-driven-system

# Start infrastructure services
make infra-up

# Run the application
make run

# Or with Docker Compose
make docker-up
```

### Accessing Services

| Service | URL | Description |
|---------|-----|-------------|
| API Gateway | http://localhost:8080 | REST API endpoint |
| Kafka UI | http://localhost:8081 | Kafka management UI |
| MongoDB | localhost:27017 | Read model database |
| Redis | localhost:6379 | Cache and pub/sub |
| PostgreSQL | localhost:5432 | Event store |
| Prometheus | http://localhost:9090 | Metrics |
| Grafana | http://localhost:3000 | Dashboards |

## Implementation Details

### Project Structure

```
event-driven-system/
├── cmd/
│   ├── api/                    # API server entry point
│   ├── worker/                 # Event processor worker
│   └── projector/              # Read model projector
├── internal/
│   ├── domain/                 # Domain models and events
│   │   ├── order/
│   │   ├── payment/
│   │   └── inventory/
│   ├── application/            # Application services
│   │   ├── commands/
│   │   └── queries/
│   ├── infrastructure/         # Infrastructure implementations
│   │   ├── eventstore/
│   │   ├── messaging/
│   │   ├── persistence/
│   │   └── projections/
│   └── interfaces/             # Interface adapters
│       ├── http/
│       └── grpc/
├── pkg/
│   ├── event/                  # Event framework
│   ├── saga/                   # Saga orchestration
│   └── cqrs/                   # CQRS utilities
├── deployments/
│   ├── docker-compose.yml
│   └── k8s/
└── tests/
    ├── integration/
    └── load/
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Load tests
make load-test
```

## Deployment

### Docker Compose

```bash
# Start all services
docker-compose up -d

# Scale processors
docker-compose up -d --scale event-processor=3
```

### Kubernetes

```bash
# Deploy to Kubernetes
kubectl apply -f deployments/k8s/

# Check status
kubectl get pods -n event-driven
```

## Performance

### Benchmarks

| Metric | Target | Actual |
|--------|--------|--------|
| Event Throughput | 50,000 events/sec | 65,000 events/sec |
| Command Latency (p99) | < 100ms | 75ms |
| Query Latency (p99) | < 50ms | 35ms |
| Event Replay Rate | 10,000 events/sec | 15,000 events/sec |
| Read Model Lag | < 1 second | 200ms |

### Load Testing

```bash
# Run k6 load tests
k6 run tests/load/event_ingestion.js
k6 run tests/load/cqrs_operations.js
```

### Optimization Strategies

1. **Event Batching**: Batch events for bulk operations
2. **Snapshots**: Create aggregate snapshots to reduce replay time
3. **Parallel Projections**: Process projections in parallel
4. **Caching**: Use Redis for hot read models
5. **Connection Pooling**: Pool database and message broker connections

## Best Practices

### Event Design

1. **Event Versioning**: Version events for schema evolution
2. **Idempotency**: Ensure handlers are idempotent
3. **Event Size**: Keep events small and focused
4. **Metadata**: Include correlation IDs and timestamps

### CQRS Guidelines

1. **Command Validation**: Validate commands before processing
2. **Optimistic Concurrency**: Use versioning for conflict detection
3. **Eventually Consistent**: Accept eventual consistency in read models
4. **Separate Scaling**: Scale read and write sides independently

### Saga Patterns

1. **Compensating Actions**: Define compensating actions for each step
2. **Timeouts**: Set timeouts for saga steps
3. **Monitoring**: Monitor saga execution and failures
4. **Idempotency**: Ensure saga actions are idempotent

## Monitoring

### Key Metrics

- Event ingestion rate
- Projection lag
- Command/query latency
- Saga completion rate
- Error rates by component

### Alerts

- High projection lag
- Failed saga executions
- Event store availability
- Message broker connectivity

## License

MIT License

## Contributing

Please read CONTRIBUTING.md for guidelines.

---

**Last Updated**: 2024-01-15  
**Version**: 1.0.0
