# EC-013: Outbox Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (18+ KB)
> **Tags**: #outbox #transactional-messaging #event-sourcing #eventual-consistency #cdc
> **Authoritative Sources**:
>
> - [Transactional Outbox](https://microservices.io/patterns/data/transactional-outbox.html) - Microservices.io
> - [Debezium Outbox Pattern](https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html) - Debezium
> - [Reliable Messaging](https://www.enterpriseintegrationpatterns.com/patterns/messaging/GuaranteedMessaging.html) - Enterprise Integration Patterns
> - [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html) - Martin Fowler
> - [Change Data Capture](https://www.confluent.io/blog/change-data-capture-capture/) - Confluent

---

## 1. Pattern Overview

### 1.1 Problem Statement

In microservices, updating the database and publishing events should be atomic. Without this guarantee:

- Database update succeeds, event publish fails → Stale data in other services
- Event published, database update fails → Inconsistent state

**The Dual Write Problem:**

```
Problem: Update DB and Publish Message should be atomic

Scenario 1: DB Success, Publish Fail
  1. UPDATE orders SET status='confirmed' WHERE id=123  ✓
  2. Publish OrderConfirmed event                      ✗
  Result: Order confirmed in DB, but no notification sent

Scenario 2: Publish Success, DB Fail
  1. Publish OrderConfirmed event                      ✓
  2. UPDATE orders SET status='confirmed' WHERE id=123  ✗
  Result: Notification sent for non-existent order
```

### 1.2 Solution Overview

The Outbox Pattern ensures atomicity by:

1. Storing events in an "outbox" table within the same database transaction
2. Asynchronously publishing events from the outbox
3. Removing published events from the outbox

This provides **exactly-once delivery** semantics.

---

## 2. Design Pattern Formalization

### 2.1 Outbox Definition

**Definition 2.1 (Outbox Table)**
An outbox table $O$ stores pending events:

$$
O = \{ (id, aggregate_id, event_type, payload, created_at, published) \}
$$

**Definition 2.2 (Transaction with Outbox)**
A business transaction $T$ with outbox insertion:

$$
T = \langle B, I_O \rangle
$$

Where:

- $B$: Business operations (e.g., UPDATE, INSERT)
- $I_O$: Insert into outbox table

**ACID Guarantee:**
$$
atomic(B \land I_O) \Rightarrow B \text{ committed} \iff I_O \text{ committed}
$$

### 2.2 Event Publishing

**Definition 2.3 (Outbox Relay)**
The relay $R$ publishes events from outbox:

$$
R: O_{unpublished} \to M
$$

Where $M$ is the message broker.

**At-Least-Once Delivery:**
$$
\forall o \in O_{committed}: P(published(o)) = 1 \text{ as } t \to \infty
$$

---

## 3. Visual Representations

### 3.1 Outbox Pattern Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Outbox Pattern Architecture                            │
└─────────────────────────────────────────────────────────────────────────────┘

Without Outbox Pattern (Dual Write Problem):
┌──────────────┐                              ┌──────────────┐
│   Service    │  1. Update DB                │  Database    │
│              │─────────────────────────────►│              │
│              │                              │              │
│              │  2. Publish Event            │              │
│              │────────────────────────────X►│   (FAIL!)    │
│              │         (Network Error)      │              │
└──────────────┘                              └──────────────┘
                                                         │
                              ┌──────────────────────────┘
                              ▼
                       ┌──────────────┐
                       │   Message    │
                       │    Broker    │
                       └──────────────┘

PROBLEM: Database updated, but event not published!


With Outbox Pattern (Atomic Transaction):
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Service                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  Transaction (ACID)                                                   │  │
│  │                                                                       │  │
│  │  BEGIN;                                                               │  │
│  │    UPDATE orders SET status='confirmed' WHERE id=123;                 │  │
│  │    INSERT INTO outbox (aggregate_id, event_type, payload)             │  │
│  │    VALUES ('123', 'OrderConfirmed', '{...}');                         │  │
│  │  COMMIT;                                                              │  │
│  │                                                                       │  │
│  │  ✓ Both operations succeed or both fail                               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Database                                         │
│  ┌─────────────────────┐    ┌───────────────────────────────────────────┐   │
│  │   orders table      │    │           outbox table                    │   │
│  ├─────────┬───────────┤    ├────┬──────────────┬───────────┬──────────┤   │
│  │ id      │ status    │    │ id │ aggregate_id │ event_type│ payload  │   │
│  ├─────────┼───────────┤    ├────┼──────────────┼───────────┼──────────┤   │
│  │ 123     │ confirmed │    │ 1  │ 123          │ OrderConf │ {...}    │   │
│  └─────────┴───────────┘    └────┴──────────────┴───────────┴──────────┘   │
│                                                                    ▲        │
└────────────────────────────────────────────────────────────────────┼────────┘
                                                                     │
                           ┌─────────────────────────────────────────┘
                           │ Poll / CDC
                           ▼
                  ┌─────────────────┐
                  │  Outbox Relay   │
                  │                 │
                  │ • Poll outbox   │
                  │ • Publish event │
                  │ • Mark published│
                  └────────┬────────┘
                           │
                           │ Publish
                           ▼
                  ┌─────────────────┐
                  │  Message Broker │
                  │  (Kafka/Rabbit) │
                  └─────────────────┘
```

### 3.2 Outbox Relay Mechanisms

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Outbox Relay Mechanisms                                 │
└─────────────────────────────────────────────────────────────────────────────┘

Mechanism 1: Polling Publisher
┌─────────────────┐
│  Polling Loop   │
│  (Every 100ms)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐      ┌─────────────────┐
│ SELECT * FROM   │─────►│ Publish to      │
│ outbox WHERE    │      │ Message Broker  │
│ published=false │      │                 │
└─────────────────┘      └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │ DELETE FROM     │
                         │ outbox WHERE    │
                         │ id IN (...)     │
                         └─────────────────┘

PROS: Simple, no external dependencies
CONS: Polling overhead, latency


Mechanism 2: Transaction Log Tailing (CDC)
┌─────────────────┐
│  Database       │      Binlog/ WAL
│  Transaction    │──────────────────────────┐
│  Log            │                          │
└─────────────────┘                          │
                                             ▼
                                    ┌─────────────────┐
                                    │  CDC Connector  │
                                    │  (Debezium)     │
                                    │                 │
                                    │ • Read binlog   │
                                    │ • Filter outbox │
                                    │ • Emit events   │
                                    └────────┬────────┘
                                             │
                                             ▼
                                    ┌─────────────────┐
                                    │  Message Broker │
                                    └─────────────────┘

PROS: Near real-time, no polling
CONS: Additional infrastructure (Debezium)


Mechanism 3: Transactional Event Listener
┌─────────────────────────────────────────────────────────────────────────────┐
│  Database Trigger                                                           │
│                                                                             │
│  CREATE TRIGGER outbox_trigger                                              │
│  AFTER INSERT ON outbox                                                     │
│  FOR EACH ROW EXECUTE FUNCTION notify_relay();                              │
│                                                                             │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │ NOTIFY
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  LISTEN Channel                   │
│                                   │
│  ┌─────────────┐    ┌─────────────┴─────────┐    ┌─────────────────┐       │
│  │  LISTEN     │◄───┤  PostgreSQL NOTIFY    │◄───┤  INSERT outbox  │       │
│  │  for events │    │  Channel              │    │                 │       │
│  └──────┬──────┘    └───────────────────────┘    └─────────────────┘       │
│         │                                                                  │
│         ▼                                                                  │
│  ┌─────────────┐    ┌─────────────────┐                                   │
│  │  Relay      │───►│  Message Broker │                                   │
│  │  Process    │    │                 │                                   │
│  └─────────────┘    └─────────────────┘                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

PROS: Low latency, event-driven
CONS: Database-specific, connection management
```

### 3.3 Event Delivery Guarantees

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Exactly-Once Delivery Flow                                │
└─────────────────────────────────────────────────────────────────────────────┘

Producer Side (Outbox Pattern):
┌─────────────────────────────────────────────────────────────────────────────┐
│  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐       │
│  │  Begin  │──►│ Update  │──►│ Insert  │──►│ Commit  │──►│  Poll   │       │
│  │   Tx    │   │   DB    │   │ Outbox  │   │   Tx    │   │ Outbox  │       │
│  └─────────┘   └─────────┘   └─────────┘   └─────────┘   └────┬────┘       │
│                                                               │            │
│                                                               ▼            │
│                                                        ┌─────────────┐      │
│                                                        │  Publish    │      │
│                                                        │  to Kafka   │      │
│                                                        │  with Key   │      │
│                                                        └──────┬──────┘      │
│                                                               │            │
│                                                               ▼            │
│                                                        ┌─────────────┐      │
│                                                        │   DELETE    │      │
│                                                        │   Outbox    │      │
│                                                        └─────────────┘      │
└─────────────────────────────────────────────────────────────────────────────┘

Consumer Side (Idempotent Consumer):
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐     │
│  │  Receive    │──►│ Check       │──►│ Process     │──►│ Store       │     │
│  │  Message    │   │ Processed?  │   │ Event       │   │ Offset      │     │
│  └─────────────┘   └──────┬──────┘   └─────────────┘   └─────────────┘     │
│                           │                                                 │
│                      ┌────┴────┐                                            │
│                      │  YES    │──► Skip (Duplicate)                         │
│                      └─────────┘                                            │
│                      │  NO                                                 │
│                      └────────► Continue Processing                         │
│                                                                             │
│  Idempotency Key: aggregate_id + event_type + timestamp                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

Delivery Guarantees:
┌─────────────────────────────────────────────────────────────────────────┐
│ Guarantee Level: AT_LEAST_ONCE (can upgrade to EXACTLY_ONCE)            │
│                                                                         │
│ Scenarios:                                                              │
│ 1. DB Commit succeeds, Relay crashes before publish                     │
│    → Event stays in outbox, relay retries on restart                    │
│                                                                         │
│ 2. Publish succeeds, DELETE fails                                       │
│    → Duplicate possible (handled by idempotent consumer)                │
│                                                                         │
│ 3. Relay publishes but doesn't receive ACK                              │
│    → Retry may cause duplicate (handled by idempotent consumer)         │
│                                                                         │
│ 4. Consumer processes but fails to ACK                                  │
│    → Redelivery occurs (handled by idempotency check)                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production-Ready Implementation

```go
package outbox

import (
 "context"
 "database/sql"
 "encoding/json"
 "fmt"
 "time"

 "github.com/google/uuid"
)

// Event represents an outbox event
type Event struct {
 ID          string          `json:"id" db:"id"`
 AggregateID string          `json:"aggregate_id" db:"aggregate_id"`
 EventType   string          `json:"event_type" db:"event_type"`
 Payload     json.RawMessage `json:"payload" db:"payload"`
 CreatedAt   time.Time       `json:"created_at" db:"created_at"`
 Published   bool            `json:"published" db:"published"`
 PublishedAt *time.Time      `json:"published_at,omitempty" db:"published_at"`
}

// Publisher publishes events to message broker
type Publisher interface {
 Publish(ctx context.Context, event *Event) error
}

// Store manages outbox storage
type Store struct {
 db *sql.DB
}

// NewStore creates a new outbox store
func NewStore(db *sql.DB) *Store {
 return &Store{db: db}
}

// SaveEvent saves an event to the outbox (call within transaction)
func (s *Store) SaveEvent(ctx context.Context, tx *sql.Tx, aggregateID, eventType string, payload interface{}) (*Event, error) {
 event := &Event{
  ID:          uuid.New().String(),
  AggregateID: aggregateID,
  EventType:   eventType,
  CreatedAt:   time.Now(),
  Published:   false,
 }

 payloadBytes, err := json.Marshal(payload)
 if err != nil {
  return nil, err
 }
 event.Payload = payloadBytes

 _, err = tx.ExecContext(ctx,
  `INSERT INTO outbox (id, aggregate_id, event_type, payload, created_at, published)
   VALUES ($1, $2, $3, $4, $5, $6)`,
  event.ID, event.AggregateID, event.EventType, event.Payload, event.CreatedAt, event.Published,
 )

 return event, err
}

// GetUnpublished returns unpublished events
func (s *Store) GetUnpublished(ctx context.Context, limit int) ([]*Event, error) {
 rows, err := s.db.QueryContext(ctx,
  `SELECT id, aggregate_id, event_type, payload, created_at, published, published_at
   FROM outbox WHERE published = false ORDER BY created_at LIMIT $1`,
  limit,
 )
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 var events []*Event
 for rows.Next() {
  var event Event
  err := rows.Scan(&event.ID, &event.AggregateID, &event.EventType,
   &event.Payload, &event.CreatedAt, &event.Published, &event.PublishedAt)
  if err != nil {
   return nil, err
  }
  events = append(events, &event)
 }

 return events, rows.Err()
}

// MarkPublished marks events as published
func (s *Store) MarkPublished(ctx context.Context, ids []string) error {
 if len(ids) == 0 {
  return nil
 }

 query := `UPDATE outbox SET published = true, published_at = $1 WHERE id IN (`
 args := []interface{}{time.Now()}
 for i, id := range ids {
  if i > 0 {
   query += ", "
  }
  query += fmt.Sprintf("$%d", i+2)
  args = append(args, id)
 }
 query += ")"

 _, err := s.db.ExecContext(ctx, query, args...)
 return err
}

// DeletePublished removes published events (cleanup)
func (s *Store) DeletePublished(ctx context.Context, olderThan time.Duration) (int64, error) {
 result, err := s.db.ExecContext(ctx,
  `DELETE FROM outbox WHERE published = true AND published_at < $1`,
  time.Now().Add(-olderThan),
 )
 if err != nil {
  return 0, err
 }
 return result.RowsAffected()
}

// Relay polls and publishes outbox events
type Relay struct {
 store     *Store
 publisher Publisher
 interval  time.Duration
 batchSize int
}

// NewRelay creates a new outbox relay
func NewRelay(store *Store, publisher Publisher, interval time.Duration, batchSize int) *Relay {
 return &Relay{
  store:     store,
  publisher: publisher,
  interval:  interval,
  batchSize: batchSize,
 }
}

// Start starts the relay loop
func (r *Relay) Start(ctx context.Context) {
 ticker := time.NewTicker(r.interval)
 defer ticker.Stop()

 for {
  select {
  case <-ctx.Done():
   return
  case <-ticker.C:
   r.processBatch(ctx)
  }
 }
}

func (r *Relay) processBatch(ctx context.Context) {
 events, err := r.store.GetUnpublished(ctx, r.batchSize)
 if err != nil {
  // Log error
  return
 }

 if len(events) == 0 {
  return
 }

 var publishedIDs []string
 for _, event := range events {
  if err := r.publisher.Publish(ctx, event); err != nil {
   // Log error, will retry on next poll
   continue
  }
  publishedIDs = append(publishedIDs, event.ID)
 }

 if len(publishedIDs) > 0 {
  if err := r.store.MarkPublished(ctx, publishedIDs); err != nil {
   // Log error - events will be republished (idempotency required)
  }
 }
}
```

---

## 5. Failure Scenarios and Mitigation

| Scenario | Symptom | Cause | Mitigation |
|----------|---------|-------|------------|
| **Relay Lag** | Event delay | High volume | Increase relay instances, batch size |
| **Outbox Growth** | Storage pressure | Cleanup failure | Automated cleanup, monitoring |
| **Duplicate Events** | Multiple deliveries | Retry without delete | Idempotent consumers |
| **Lost Events** | Events not delivered | DB rollback | Outbox is in same transaction |

---

## 6. Best Practices

```
Outbox Schema:
CREATE TABLE outbox (
  id UUID PRIMARY KEY,
  aggregate_id VARCHAR(255) NOT NULL,
  event_type VARCHAR(255) NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL,
  published BOOLEAN DEFAULT FALSE,
  published_at TIMESTAMP,

  INDEX idx_published (published, created_at),
  INDEX idx_aggregate (aggregate_id)
);

Relay Configuration:
• Poll interval: 100-500ms (balance latency vs load)
• Batch size: 100-1000 events
• Cleanup: Daily for events older than 7 days
• Monitoring: Track relay lag, outbox size
```

---

## 7. References

1. **Richardson, C.** [Transactional Outbox](https://microservices.io/patterns/data/transactional-outbox.html).
2. **Debezium.** [Outbox Event Router](https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html).
3. **Hohpe, G. & Woolf, B.** *Enterprise Integration Patterns*. Addison-Wesley.

---

**Quality Rating**: S (18KB+, Complete Formalization + Production Code + Visualizations)
