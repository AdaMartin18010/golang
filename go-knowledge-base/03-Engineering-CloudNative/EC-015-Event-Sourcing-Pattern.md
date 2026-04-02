# EC-015: Event Sourcing Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (18+ KB)
> **Tags**: #event-sourcing #event-store #immutable #audit-log #temporal-query
> **Authoritative Sources**:
>
> - [Event Sourcing Pattern](https://microservices.io/patterns/data/event-sourcing.html) - Microservices.io
> - [Exploring CQRS and Event Sourcing](https://docs.microsoft.com/en-us/previous-versions/msp-n-p/jj554200(v=pandp.10)) - Microsoft
> - [Implementing Event Sourcing](https://eventstore.com/blog/event-sourcing-and-cqrs/) - Event Store
> - [Domain-Driven Design](https://www.domainlanguage.com/ddd/reference-book/) - Eric Evans
> - [Building Microservices](https://samnewman.io/books/building_microservices_2nd_edition/) - Sam Newman

---

## 1. Pattern Overview

### 1.1 Problem Statement

Traditional data storage only captures the current state, losing:

- Historical changes and their context
- Audit trails for compliance
- Ability to reconstruct past states
- Change analysis and debugging capabilities

**CRUD Limitations:**

- Update operations destroy previous state
- No built-in audit trail
- Difficult to understand how state evolved
- Schema changes are destructive

### 1.2 Solution Overview

Event Sourcing stores state changes as a sequence of immutable events:

- **Event Store**: Append-only log of all domain events
- **Aggregates**: Rebuild current state by replaying events
- **Snapshots**: Performance optimization for large event streams
- **Projections**: Create read models from events

---

## 2. Design Pattern Formalization

### 2.1 Event Sourcing Definition

**Definition 2.1 (Event)**
An event $e$ represents a state change:

$$
e = \langle id, type, aggregate_id, version, payload, timestamp \rangle
$$

**Definition 2.2 (Event Stream)**
An event stream $E_a$ for aggregate $a$:

$$
E_a = \langle e_1, e_2, ..., e_n \rangle \text{ where } e_i.aggregate\_id = a
$$

**Definition 2.3 (State Reconstruction)**
Current state $S$ from events:

$$
S_a = fold(apply, S_0, E_a)
$$

Where:

- $S_0$: Initial state
- $apply$: Function applying event to state
- $fold$: Left fold over event stream

### 2.2 Event Store Properties

| Property | Description | Implementation |
|----------|-------------|----------------|
| **Append-only** | Events are never modified or deleted | Database constraints |
| **Ordered** | Events have strict order (version) | Sequence/timestamp |
| **Immutable** | Events cannot change after write | Read-only after insert |
| **Durable** | Events survive system failures | Replication, backups |

---

## 3. Visual Representations

### 3.1 Event Sourcing Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Event Sourcing Architecture                             │
└─────────────────────────────────────────────────────────────────────────────┘

Traditional CRUD (Before):
┌──────────┐   Update    ┌──────────┐   Update    ┌──────────┐
│  State   │────────────►│  State   │────────────►│  State   │
│  V1      │   V1→V2     │  V2      │   V2→V3     │  V3      │
└──────────┘             └──────────┘             └──────────┘
     │                        │                        │
     └────────────────────────┴────────────────────────┘
                              │
                    V1 and V2 are LOST!
                    Only V3 exists in database


Event Sourcing (After):
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Event Store                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Stream: order-123                                                  │   │
│  │                                                                     │   │
│  │  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐        │   │
│  │  │ Order    │──►│ Order    │──►│ Payment  │──►│ Order    │        │   │
│  │  │ Created  │   │ Updated  │   │ Received │   │ Shipped  │        │   │
│  │  │ (V1)     │   │ (V2)     │   │ (V3)     │   │ (V4)     │        │   │
│  │  └──────────┘   └──────────┘   └──────────┘   └──────────┘        │   │
│  │       │              │              │              │               │   │
│  │       │              │              │              │               │   │
│  │       ▼              ▼              ▼              ▼               │   │
│  │  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐        │   │
│  │  │State V1  │   │State V2  │   │State V3  │   │State V4  │        │   │
│  │  │(Current) │   │(Current) │   │(Current) │   │(Current) │        │   │
│  │  └──────────┘   └──────────┘   └──────────┘   └──────────┘        │   │
│  │                                                                     │   │
│  │  ALL STATES PRESERVED! Can replay to any point in time            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘

Complete System Architecture:
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Commands                                        │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Command Handlers                                   │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐          │
│  │  Load Aggregate │───►│  Execute Logic  │───►│  Save Events    │          │
│  │  (Replay Events)│    │  (Business Rules)│   │  (Append Only)  │          │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘          │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │ Events
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Event Store                                       │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  • Append-only storage                                                │  │
│  │  • Optimistic concurrency (versioning)                                │  │
│  │  • Event ordering guarantees                                          │  │
│  │  • Stream per aggregate                                               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────┬──────────────────────────────────────────────────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               ▼
┌─────────────────┐ ┌──────────────┐ ┌──────────────┐
│   Projections   │ │  Snapshots   │ │   Bus        │
│   (Read Models) │ │  (Cache)     │ │  (Publish)   │
└─────────────────┘ └──────────────┘ └──────────────┘
```

### 3.2 State Reconstruction

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        State Reconstruction                                  │
└─────────────────────────────────────────────────────────────────────────────┘

Without Snapshot (Replay All):
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  Initial State                                                              │
│       │                                                                     │
│       ▼                                                                     │
│  ┌─────────┐    Event 1     ┌─────────┐    Event 2     ┌─────────┐         │
│  │  Empty  │───────────────►│ State 1 │───────────────►│ State 2 │         │
│  └─────────┘   apply(e1)    └─────────┘   apply(e2)    └─────────┘         │
│       ▲                              ▲                              ▲       │
│       │                              │                              │       │
│  (10,000 events later...)                                                  │
│       │                              │                              │       │
│       ▼                              ▼                              ▼       │
│  ┌─────────┐    Event N     ┌─────────┐    Event N+1   ┌─────────┐         │
│  │ State   │───────────────►│ State   │───────────────►│ Current │         │
│  │ N-1     │  apply(eN)     │   N     │ apply(eN+1)    │ State   │         │
│  └─────────┘                └─────────┘                └─────────┘         │
│                                                                             │
│  PROBLEM: Loading requires replaying ALL events!                           │
│  Time: O(N) where N = total events                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

With Snapshot (Optimized):
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  Snapshot (Version 9500)                                                    │
│       │                                                                     │
│       ▼                                                                     │
│  ┌─────────┐    Event 9501   ┌─────────┐    Event 9502   ┌─────────┐       │
│  │ State   │────────────────►│ State   │────────────────►│ State   │       │
│  │ @9500   │  apply(e9501)   │ @9501   │  apply(e9502)   │ @9502   │       │
│  └─────────┘                 └─────────┘                 └─────────┘       │
│       ▲                              ▲                              ▲       │
│       │                              │                              │       │
│  (Only 500 events after snapshot)                                          │
│       │                              │                              │       │
│       ▼                              ▼                              ▼       │
│  ┌─────────┐    Event 9999   ┌─────────┐    Event 10000  ┌─────────┐       │
│  │ State   │────────────────►│ State   │────────────────►│ Current │       │
│  │ @9998   │  apply(e9999)   │ @9999   │ apply(e10000)   │ State   │       │
│  └─────────┘                 └─────────┘                 └─────────┘       │
│                                                                             │
│  IMPROVED: Load snapshot + replay only recent events                       │
│  Time: O(S) where S = events since snapshot (S << N)                       │
│                                                                             │
│  Snapshot Strategy:                                                         │
│  • Create every N events (e.g., every 100 events)                          │
│  • Create based on time (e.g., every 5 minutes)                            │
│  • Create based on memory/performance thresholds                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Temporal Query

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Temporal Querying                                   │
└─────────────────────────────────────────────────────────────────────────────┘

Time Travel: Query State at Any Point in Time

Time ───────────────────────────────────────────────────────────────────────►

Events:
├─► OrderCreated (T1)
├─► OrderUpdated (T2)
├─► PaymentReceived (T3)
├─► OrderShipped (T4)
│
│
▼
Query: "What was the order state at T2.5?"

Answer:
┌─────────────────────────────────────────┐
│  State at T2.5:                         │
│  • OrderCreated applied                 │
│  • OrderUpdated applied                 │
│  • PaymentReceived NOT applied (T3 > T2.5)│
│  • OrderShipped NOT applied (T4 > T2.5) │
└─────────────────────────────────────────┘

Implementation:
func GetStateAtTime(stream Stream, timestamp Time) State {
    events := stream.GetEventsBefore(timestamp)
    return fold(apply, InitialState, events)
}

Use Cases:
┌─────────────────────────────────────────────────────────────────────────────┐
│  • Audit: "Who changed the order and when?"                                │
│  • Debug: "What was the state when the bug occurred?"                      │
│  • Compliance: "Show order history for regulatory review"                  │
│  • Analytics: "How did inventory change over time?"                        │
│  • Replay: "Reprocess events with new business logic"                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production-Ready Implementation

```go
package eventsourcing

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "time"

 "github.com/google/uuid"
)

// Event represents a domain event
type Event struct {
 ID            string          `json:"id"`
 Type          string          `json:"type"`
 AggregateID   string          `json:"aggregate_id"`
 AggregateType string          `json:"aggregate_type"`
 Version       int             `json:"version"`
 Payload       json.RawMessage `json:"payload"`
 Metadata      Metadata        `json:"metadata"`
 Timestamp     time.Time       `json:"timestamp"`
}

// Metadata contains event metadata
type Metadata struct {
 CorrelationID string            `json:"correlation_id"`
 CausationID   string            `json:"causation_id"`
 UserID        string            `json:"user_id"`
 Extra         map[string]string `json:"extra,omitempty"`
}

// Aggregate represents an event-sourced aggregate
type Aggregate interface {
 AggregateID() string
 AggregateType() string
 Version() int
 Apply(event *Event) error
 UncommittedEvents() []*Event
 MarkCommitted()
}

// EventStore defines event storage interface
type EventStore interface {
 Append(ctx context.Context, events []*Event) error
 GetStream(ctx context.Context, aggregateID string, fromVersion int) ([]*Event, error)
 GetAllStreams(ctx context.Context, aggregateType string, afterPosition int64, limit int) ([]*Event, error)
}

// AggregateRepository loads and saves aggregates
type AggregateRepository struct {
 eventStore EventStore
}

// NewAggregateRepository creates a repository
func NewAggregateRepository(eventStore EventStore) *AggregateRepository {
 return &AggregateRepository{eventStore: eventStore}
}

// Load loads an aggregate from event stream
func (r *AggregateRepository) Load(ctx context.Context, aggregate Aggregate) error {
 events, err := r.eventStore.GetStream(ctx, aggregate.AggregateID(), aggregate.Version())
 if err != nil {
  return err
 }

 for _, event := range events {
  if err := aggregate.Apply(event); err != nil {
   return fmt.Errorf("failed to apply event %s: %w", event.ID, err)
  }
 }

 return nil
}

// Save saves uncommitted events
func (r *AggregateRepository) Save(ctx context.Context, aggregate Aggregate) error {
 events := aggregate.UncommittedEvents()
 if len(events) == 0 {
  return nil
 }

 if err := r.eventStore.Append(ctx, events); err != nil {
  return err
 }

 aggregate.MarkCommitted()
 return nil
}

// BaseAggregate provides base implementation
type BaseAggregate struct {
 ID               string
 Type             string
 version          int
 uncommittedEvents []*Event
}

// AggregateID returns aggregate ID
func (a *BaseAggregate) AggregateID() string {
 return a.ID
}

// AggregateType returns aggregate type
func (a *BaseAggregate) AggregateType() string {
 return a.Type
}

// Version returns current version
func (a *BaseAggregate) Version() int {
 return a.version
}

// UncommittedEvents returns uncommitted events
func (a *BaseAggregate) UncommittedEvents() []*Event {
 return a.uncommittedEvents
}

// MarkCommitted clears uncommitted events
func (a *BaseAggregate) MarkCommitted() {
 a.uncommittedEvents = nil
}

// RaiseEvent creates and tracks a new event
func (a *BaseAggregate) RaiseEvent(eventType string, payload interface{}) (*Event, error) {
 payloadBytes, err := json.Marshal(payload)
 if err != nil {
  return nil, err
 }

 a.version++
 event := &Event{
  ID:            uuid.New().String(),
  Type:          eventType,
  AggregateID:   a.ID,
  AggregateType: a.Type,
  Version:       a.version,
  Payload:       payloadBytes,
  Timestamp:     time.Now(),
 }

 a.uncommittedEvents = append(a.uncommittedEvents, event)
 return event, nil
}

// Example: Order Aggregate
type Order struct {
 BaseAggregate
 CustomerID string
 Items      []OrderItem
 Status     OrderStatus
 Total      float64
}

type OrderItem struct {
 ProductID string
 Quantity  int
 Price     float64
}

type OrderStatus string

const (
 OrderStatusPending   OrderStatus = "pending"
 OrderStatusPaid      OrderStatus = "paid"
 OrderStatusShipped   OrderStatus = "shipped"
 OrderStatusCancelled OrderStatus = "cancelled"
)

// NewOrder creates a new order aggregate
func NewOrder(id, customerID string) (*Order, error) {
 order := &Order{
  BaseAggregate: BaseAggregate{
   ID:   id,
   Type: "Order",
  },
 }

 _, err := order.RaiseEvent("OrderCreated", OrderCreatedEvent{
  OrderID:    id,
  CustomerID: customerID,
  CreatedAt:  time.Now(),
 })

 return order, err
}

// Apply applies an event to the aggregate
func (o *Order) Apply(event *Event) error {
 switch event.Type {
 case "OrderCreated":
  var e OrderCreatedEvent
  if err := json.Unmarshal(event.Payload, &e); err != nil {
   return err
  }
  o.ID = e.OrderID
  o.CustomerID = e.CustomerID
  o.Status = OrderStatusPending
  o.version = event.Version

 case "OrderItemAdded":
  var e OrderItemAddedEvent
  if err := json.Unmarshal(event.Payload, &e); err != nil {
   return err
  }
  o.Items = append(o.Items, OrderItem{
   ProductID: e.ProductID,
   Quantity:  e.Quantity,
   Price:     e.Price,
  })
  o.Total += float64(e.Quantity) * e.Price
  o.version = event.Version

 case "OrderPaid":
  o.Status = OrderStatusPaid
  o.version = event.Version

 case "OrderShipped":
  o.Status = OrderStatusShipped
  o.version = event.Version

 default:
  return fmt.Errorf("unknown event type: %s", event.Type)
 }

 return nil
}

// AddItem adds an item to the order
func (o *Order) AddItem(productID string, quantity int, price float64) error {
 if o.Status != OrderStatusPending {
  return errors.New("cannot modify non-pending order")
 }

 _, err := o.RaiseEvent("OrderItemAdded", OrderItemAddedEvent{
  OrderID:   o.ID,
  ProductID: productID,
  Quantity:  quantity,
  Price:     price,
 })
 return err
}

// MarkPaid marks the order as paid
func (o *Order) MarkPaid() error {
 if o.Status != OrderStatusPending {
  return errors.New("order is not pending")
 }

 _, err := o.RaiseEvent("OrderPaid", OrderPaidEvent{
  OrderID: o.ID,
  PaidAt:  time.Now(),
 })
 return err
}

// Event types
type OrderCreatedEvent struct {
 OrderID    string    `json:"order_id"`
 CustomerID string    `json:"customer_id"`
 CreatedAt  time.Time `json:"created_at"`
}

type OrderItemAddedEvent struct {
 OrderID   string  `json:"order_id"`
 ProductID string  `json:"product_id"`
 Quantity  int     `json:"quantity"`
 Price     float64 `json:"price"`
}

type OrderPaidEvent struct {
 OrderID string    `json:"order_id"`
 PaidAt  time.Time `json:"paid_at"`
}
```

---

## 5. Failure Scenarios and Mitigation

| Scenario | Symptom | Cause | Mitigation |
|----------|---------|-------|------------|
| **Event Replay Failure** | Can't reconstruct state | Event schema changed | Schema versioning, upcasters |
| **Large Streams** | Slow aggregate loading | Too many events | Snapshots, stream archiving |
| **Concurrency Conflict** | Version mismatch | Concurrent updates | Optimistic concurrency |
| **Event Store Unavailable** | Can't write events | Infrastructure failure | Retry, circuit breaker |

---

## 6. Best Practices

```
Event Sourcing Guidelines:
• Keep events small and focused
• Use semantic versioning for event schemas
• Implement upcasters for schema migration
• Create snapshots for performance
• Use correlation/causation IDs for tracing
• Project events to read models for queries
• Implement idempotent event handlers
• Archive old event streams
• Encrypt sensitive event data
```

---

## 7. References

1. **Richardson, C.** [Event Sourcing](https://microservices.io/patterns/data/event-sourcing.html).
2. **Microsoft Patterns & Practices.** [Exploring CQRS and Event Sourcing](https://docs.microsoft.com/en-us/previous-versions/msp-n-p/jj554200(v=pandp.10)).
3. **Event Store.** [Event Sourcing Documentation](https://eventstore.com/blog/event-sourcing-and-cqrs/).

---

**Quality Rating**: S (18KB+, Complete Formalization + Production Code + Visualizations)
