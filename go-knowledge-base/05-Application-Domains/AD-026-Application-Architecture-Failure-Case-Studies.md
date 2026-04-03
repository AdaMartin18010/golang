# AD-026: Application Architecture Failure Case Studies

> **维度**: Application Domains | **级别**: S (15+ KB)
> **标签**: #architecture #microservices #failures #production-incidents #postmortem
> **权威来源**: Industry Postmortems, Architecture Reviews, Real-world Incidents

---

## Overview

This document contains 10 detailed application architecture failure case studies, covering microservices decomposition issues, API design failures, data consistency problems, and scalability bottlenecks. Each case study includes incident description, root cause analysis, timeline, resolution steps, lessons learned, and prevention recommendations.

---

## Case Study 1: Death Star Architecture - Tight Service Coupling

### 1.1 Incident Description

**System**: E-commerce platform with 200+ microservices
**Impact**: Single service failure cascaded to complete platform outage
**Duration**: 4 hours
**Date**: March 2024

A seemingly simple inventory service update triggered a cascade failure across the entire platform. The "Death Star" architecture - where all services depend on each other - caused a complete outage affecting 50M users.

### 1.2 Root Cause Analysis

```
Architecture Anti-Pattern:
┌─────────────────────────────────────────────────────────────┐
│                        API Gateway                           │
│                           │                                 │
│        ┌──────────────────┼──────────────────┐              │
│        ▼                  ▼                  ▼              │
│   ┌─────────┐       ┌─────────┐       ┌─────────┐          │
│   │  Order  │◄─────►│Payment  │◄─────►│Inventory│          │
│   │ Service │       │ Service │       │ Service │          │
│   └────┬────┘       └────┬────┘       └────┬────┘          │
│        │                  │                  │               │
│        ▼                  ▼                  ▼               │
│   ┌─────────┐       ┌─────────┐       ┌─────────┐          │
│   │  User   │◄─────►│ Notification│◄──►│ Shipping│          │
│   │ Service │       │ Service │       │ Service │          │
│   └────┬────┘       └────┬────┘       └────┬────┘          │
│        │                  │                  │               │
│        └──────────────────┼──────────────────┘              │
│                           ▼                                 │
│                    ┌─────────────┐                          │
│                    │  Analytics  │                          │
│                    │   Service   │                          │
│                    └─────────────┘                          │
└─────────────────────────────────────────────────────────────┘

Failure Cascade:
1. Inventory Service deployed with DB migration
2. Migration locked inventory table for 30s
3. Order Service timeouts waiting for inventory check
4. Order Service thread pool exhausted
5. API Gateway retries amplified load
6. Payment Service dependent on Order status
7. Shipping Service waiting for Payment confirmation
8. Complete platform deadlock
```

### 1.3 Timeline of Events

| Time (UTC) | Event | Error Rate |
|------------|-------|------------|
| 14:00:00 | Inventory Service deployment starts | 0% |
| 14:05:00 | Database migration begins | 0.1% |
| 14:10:00 | Migration still running, timeouts begin | 5% |
| 14:15:00 | Order Service thread pool exhausted | 30% |
| 14:20:00 | API Gateway retry storm | 60% |
| 14:25:00 | Payment Service degradation | 80% |
| 14:30:00 | Complete platform outage | 100% |
| 15:00:00 | Emergency rollback initiated | 90% |
| 16:00:00 | Services recovering | 40% |
| 18:00:00 | Full recovery | 0% |

### 1.4 Resolution Steps

```go
// Circuit breaker and bulkhead pattern
type ResilientInventoryClient struct {
    circuitBreaker *gobreaker.CircuitBreaker
    bulkhead       *Bulkhead
    timeout        time.Duration
}

func (c *ResilientInventoryClient) CheckAvailability(ctx context.Context, items []Item) (*AvailabilityResult, error) {
    return c.bulkhead.Execute(ctx, func() (interface{}, error) {
        return c.circuitBreaker.Execute(func() (interface{}, error) {
            ctx, cancel := context.WithTimeout(ctx, c.timeout)
            defer cancel()
            return c.client.CheckAvailability(ctx, items)
        })
    })
}

// Async inventory check with fallback
func (c *ResilientInventoryClient) CheckAvailabilityAsync(ctx context.Context, items []Item) (*AvailabilityResult, error) {
    resultChan := make(chan *AvailabilityResult, 1)

    go func() {
        result, err := c.CheckAvailability(ctx, items)
        if err != nil {
            // Return cached availability
            resultChan <- c.getCachedAvailability(items)
            return
        }
        resultChan <- result
    }()

    select {
    case result := <-resultChan:
        return result, nil
    case <-time.After(100 * time.Millisecond):
        return c.getCachedAvailability(items), ErrTimeoutUsingCache
    }
}
```

### 1.5 Lessons Learned

1. **Death Star architectures** fail catastrophically
2. **Synchronous chains** multiply failure probability
3. **Shared databases** create coupling through data
4. **Fallback strategies** must be designed upfront

### 1.6 Prevention Recommendations

```go
// Decoupled architecture with event sourcing
type OrderService struct {
    eventBus      EventBus
    orderRepo     OrderRepository
}

func (s *OrderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) error {
    order, err := domain.NewOrder(cmd)
    if err != nil {
        return err
    }

    if err := s.orderRepo.Save(ctx, order); err != nil {
        return err
    }

    return s.eventBus.Publish(ctx, &OrderCreatedEvent{
        OrderID: order.ID,
        Items:   order.Items,
    })
}
```

---

## Case Study 2: Database Per Service - Distributed Transaction Hell

### 2.1 Incident Description

**System**: Financial trading platform with microservices
**Impact**: Data inconsistency, $2M trading losses
**Duration**: 6 hours (detection + recovery)
**Date**: February 2024

The "database per service" pattern led to complex distributed transactions. A failure in the Saga orchestrator left transactions in inconsistent states, causing duplicate trades and incorrect account balances.

### 2.2 Root Cause Analysis

```
Distributed Transaction Failure:
- Trade created in TradeDB
- Account debited in AcctDB
- Risk check incomplete
- Compensation failed
- $2M unbalanced transactions
```

### 2.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 09:30:00 | Trading day begins |
| 10:15:00 | Network latency spike begins |
| 10:20:00 | First Saga timeout |
| 10:25:00 | Compensation failures begin |
| 14:00:00 | End-of-day reconciliation shows $2M discrepancy |
| 18:00:00 | Manual reconciliation completed |

### 2.4 Resolution Steps

```go
// Durable Saga with persistent state
type DurableSaga struct {
    id            string
    state         SagaState
    steps         []SagaStep
    currentStep   int
    log           SagaLog
    eventStore    EventStore
}

func (s *DurableSaga) Execute(ctx context.Context) error {
    s.state = SagaRunning
    s.persistState(ctx)

    for i, step := range s.steps[s.currentStep:] {
        s.currentStep = s.currentStep + i
        s.log.LogStepStarted(ctx, s.id, step.Name)

        err := s.executeStep(ctx, step)
        if err != nil {
            s.log.LogStepFailed(ctx, s.id, step.Name, err)
            s.persistState(ctx)
            return s.compensate(ctx)
        }

        s.log.LogStepCompleted(ctx, s.id, step.Name)
        s.persistState(ctx)
    }

    s.state = SagaCompleted
    s.persistState(ctx)
    return nil
}
```

### 2.5 Lessons Learned

1. **Distributed transactions are hard** - avoid if possible
2. **Saga compensation** must be idempotent
3. **Persistent saga state** enables recovery
4. **Event sourcing** provides better audit trail

### 2.6 Prevention Recommendations

```go
// Outbox pattern for eventual consistency
type OutboxPattern struct {
    db        *sql.DB
    eventBus  EventBus
    publisher *MessagePublisher
}

func (o *OutboxPattern) ProcessTransaction(ctx context.Context, txFunc func(*sql.Tx) error) error {
    tx, err := o.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    if err := txFunc(tx); err != nil {
        return err
    }

    events := o.collectPendingEvents()
    for _, event := range events {
        if err := o.saveToOutbox(tx, event); err != nil {
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    go o.publishPendingEvents()
    return nil
}
```

---

## Case Study 3: API Versioning Breakdown

### 3.1 Incident Description

**System**: Mobile banking API with 3M active users
**Impact**: Mobile app crashes for 500K users, $5M in failed transactions
**Duration**: 8 hours
**Date**: January 2024

A breaking API change without proper versioning caused mobile app crashes. The old app versions couldn't parse the new response format, leading to widespread service disruption.

### 3.2 Root Cause Analysis

```
Versioning Failure:
- No version negotiation in API contract
- Mobile app didn't send Accept-Version header
- Backend didn't validate version requirements
- Response structure completely changed
- No deprecation period for v1
- Old app versions (30% of users) couldn't parse response
```

### 3.3 Timeline of Events

| Time (UTC) | Event | Affected Users |
|------------|-------|----------------|
| 06:00:00 | API v2 deployment | 0 |
| 06:30:00 | Mobile app crashes begin | 10K |
| 08:00:00 | Root cause identified | 300K |
| 10:00:00 | Rollback failed (DB schema changed) | 500K |
| 11:00:00 | Compatibility layer deployed | 500K |
| 14:00:00 | All users recovered | 0 |

### 3.4 Resolution Steps

```go
// Version negotiation middleware
type VersionMiddleware struct {
    versions map[string]APIVersion
    defaultVersion string
}

func (m *VersionMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        version := r.Header.Get("Accept-Version")
        if version == "" {
            version = m.defaultVersion
        }

        apiVersion, ok := m.versions[version]
        if !ok {
            http.Error(w, "Unsupported API version", http.StatusNotAcceptable)
            return
        }

        if apiVersion.Deprecated {
            w.Header().Set("Deprecation", "true")
        }

        ctx := WithAPIVersion(r.Context(), version)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 3.5 Lessons Learned

1. **Never break backward compatibility** without deprecation
2. **Version negotiation** must be explicit
3. **Mobile apps** can't be updated instantly
4. **Compatibility layers** are essential for transitions

### 3.6 Prevention Recommendations

```yaml
# API versioning policy
apiVersioning:
  strategy: header-based
  header: Accept-Version
  default: "1.0"
  supported:
    - version: "1.0"
      status: deprecated
      sunset_date: "2024-06-01"
    - version: "2.0"
      status: current
  compatibility:
    enforce_backward_compat: true
```

---

## Case Study 4: Cache Stampede on Black Friday

### 4.1 Incident Description

**System**: E-commerce platform with Redis caching layer
**Impact**: Database overload, checkout failures, $10M revenue loss
**Duration**: 2 hours (peak shopping time)
**Date**: November 2023

A cache stampede occurred when the product catalog cache expired during Black Friday traffic. Thousands of concurrent requests hit the database simultaneously, causing connection pool exhaustion and checkout failures.

### 4.2 Root Cause Analysis

```
Cache Stampede Pattern:
- Cache expires for popular product
- 10,000 concurrent requests for same product
- All hit database simultaneously
- DB Connection Pool: 100/100 used
- Remaining requests: Queue full, timeout
```

### 4.3 Timeline of Events

| Time (UTC) | Event | DB Connections |
|------------|-------|----------------|
| 00:20:00 | Cache expires on top 100 products | 100/100 |
| 00:21:00 | Cache stampede begins | 100/100 |
| 00:25:00 | Checkout failures spike | 100/100 |
| 01:00:00 | Cache warming completed | 50/100 |
| 02:00:00 | Full recovery | 30/100 |

### 4.4 Resolution Steps

```go
// Singleflight pattern for cache stampede prevention
type CacheClient struct {
    redis     *redis.Client
    singleflight.Group
}

func (c *CacheClient) GetProduct(ctx context.Context, productID string) (*Product, error) {
    cached, err := c.redis.Get(ctx, productKey(productID)).Result()
    if err == nil {
        var product Product
        if err := json.Unmarshal([]byte(cached), &product); err == nil {
            return &product, nil
        }
    }

    v, err, _ := c.Do(productID, func() (interface{}, error) {
        return c.loadFromDB(ctx, productID)
    })

    if err != nil {
        return nil, err
    }

    return v.(*Product), nil
}
```

### 4.5 Lessons Learned

1. **Cache stampedes** can kill databases instantly
2. **Singleflight pattern** prevents duplicate work
3. **Jittered TTLs** spread expiration times
4. **Probabilistic early expiration** smooths load

### 4.6 Prevention Recommendations

```yaml
# Cache configuration
cache:
  redis:
    ttl:
      default: 5m
      jitter: 30s
    stampede_protection:
      singleflight: true
      early_expiration:
        enabled: true
        probability: 0.2
```

---

## Case Study 5: Message Queue Poison Pill

### 5.1 Incident Description

**System**: Order processing system with RabbitMQ
**Impact**: 50,000 orders stuck in queue, processing halted
**Duration**: 6 hours
**Date**: December 2023

A poison message with malformed JSON caused the consumer to crash and restart continuously. The message was requeued each time, blocking the entire queue and preventing any order processing.

### 5.2 Root Cause Analysis

```
Poison Message Pattern:
1. Consumer receives malformed message
2. JSON parsing fails
3. Consumer crashes (panic)
4. Container restarts
5. Message requeued (not acknowledged)
6. Infinite crash loop
```

### 5.3 Timeline of Events

| Time (UTC) | Event | Queue Depth |
|------------|-------|-------------|
| 10:15:00 | Poison message published | 150 |
| 10:15:01 | First consumer crash | 150 |
| 10:15:30 | Container restart loop begins | 150 |
| 14:00:00 | Queue depth critical | 30000 |
| 16:00:00 | Poison message purged | 0 |

### 5.4 Resolution Steps

```go
// Poison message handling with DLQ
type ResilientConsumer struct {
    channel       *amqp.Channel
    maxRetries    int
    dlqName       string
}

func (c *ResilientConsumer) Consume() {
    msgs, _ := c.channel.Consume("order-queue", "", false, false, false, false, nil)

    for msg := range msgs {
        retryCount := c.getRetryCount(msg)

        if retryCount >= c.maxRetries {
            c.sendToDLQ(msg)
            msg.Nack(false, false)
            continue
        }

        if err := c.processMessage(msg); err != nil {
            c.incrementRetryCount(&msg)
            msg.Nack(false, true)
            continue
        }

        msg.Ack(false)
    }
}

func (c *ResilientConsumer) processMessage(msg amqp.Delivery) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic recovered: %v", r)
        }
    }()

    var order Order
    if err := json.Unmarshal(msg.Body, &order); err != nil {
        return fmt.Errorf("unmarshal error: %w", err)
    }

    return c.processor.Process(order)
}
```

### 5.5 Lessons Learned

1. **Poison messages** can block entire queues
2. **Dead letter queues** are essential
3. **Panic recovery** prevents crash loops
4. **Retry limits** prevent infinite loops

### 5.6 Prevention Recommendations

```yaml
# RabbitMQ queue configuration
queues:
  order-processing:
    durable: true
    arguments:
      x-dead-letter-exchange: "order-dlx"
      x-dead-letter-routing-key: "order-dlq"
      x-max-retries: 3
```

---

## Case Study 6: Connection Pool Exhaustion

### 6.1 Incident Description

**System**: Payment processing service with PostgreSQL
**Impact**: Payment failures, $2M in stuck transactions
**Duration**: 1 hour 30 minutes
**Date**: October 2023

A goroutine leak in a payment webhook handler caused database connections to accumulate. When the connection pool was exhausted, all new payment requests failed, leaving transactions in an ambiguous state.

### 6.2 Root Cause Analysis

```
Connection Leak Pattern:
- Goroutine starts with connection
- Connection not released on error
- Pool size: 100 connections
- Webhook rate: 200/minute
- After 30 minutes: Pool exhausted
```

### 6.3 Timeline of Events

| Time (UTC) | Event | Pool Available |
|------------|-------|----------------|
| 14:00:00 | Deployment with new webhook handler | 100/100 |
| 14:30:00 | Pool exhausted | 0/100 |
| 14:45:00 | Payment failure rate: 80% | 0/100 |
| 15:15:00 | Hotfix deployed | 100/100 |

### 6.4 Resolution Steps

```go
// Proper connection handling with defer
type PaymentHandler struct {
    db *sql.DB
}

func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
    defer cancel()

    result, err := h.processPayment(ctx, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(result)
}

func (h *PaymentHandler) processPayment(ctx context.Context, r *http.Request) (*PaymentResult, error) {
    tx, err := h.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return &PaymentResult{Status: "success"}, nil
}
```

### 6.5 Lessons Learned

1. **Always defer resource cleanup**
2. **Don't pass connections to goroutines**
3. **Monitor connection pool metrics**
4. **Use context for cancellation**

### 6.6 Prevention Recommendations

```go
// Connection pool monitoring
type PoolMonitor struct {
    db     *sql.DB
    alerts AlertManager
}

func (m *PoolMonitor) Start() {
    ticker := time.NewTicker(10 * time.Second)
    go func() {
        for range ticker.C {
            stats := m.db.Stats()
            usage := float64(stats.InUse) / float64(stats.MaxOpenConnections)
            if usage > 0.8 {
                m.alerts.Send("DB pool at %.0f%% capacity", usage*100)
            }
        }
    }()
}
```

---

## Case Study 7: Rate Limiting Bypass

### 7.1 Incident Description

**System**: API gateway with per-IP rate limiting
**Impact**: Credential stuffing attack succeeded, 10K accounts compromised
**Duration**: 4 hours (attack duration)
**Date**: September 2023

The rate limiting implementation only tracked by IP address. Attackers used a distributed botnet (10K IPs) to bypass rate limits and perform credential stuffing attacks, compromising user accounts.

### 7.2 Root Cause Analysis

```
Rate Limiting Flaw:
- Implementation: 10 req/sec per IP
- Attack: 10,000 IPs x 10 req/sec = 100K req/sec
- Result: 10,000 accounts compromised
```

### 7.3 Timeline of Events

| Time (UTC) | Event | Failed Logins |
|------------|-------|---------------|
| 02:00:00 | Attack begins | 100/min |
| 04:00:00 | Customer reports suspicious activity | 20000/min |
| 06:30:00 | Multi-dimensional rate limiting deployed | 0/min |

### 7.4 Resolution Steps

```go
// Multi-dimensional rate limiting
type RateLimiter struct {
    redis       *redis.Client
    strategies  []LimitStrategy
}

type LimitStrategy struct {
    Name      string
    KeyFunc   func(*http.Request) string
    Rate      rate.Limit
    Burst     int
    Window    time.Duration
}

var defaultStrategies = []LimitStrategy{
    {
        Name: "ip",
        KeyFunc: func(r *http.Request) string {
            return fmt.Sprintf("ratelimit:ip:%s", getClientIP(r))
        },
        Rate:   10, Burst: 20, Window: time.Minute,
    },
    {
        Name: "global",
        KeyFunc: func(r *http.Request) string { return "ratelimit:global" },
        Rate:   10000, Burst: 20000, Window: time.Second,
    },
}
```

### 7.5 Lessons Learned

1. **IP-based rate limiting** is insufficient
2. **Multi-dimensional strategies** catch distributed attacks
3. **Endpoint-specific limits** protect sensitive operations
4. **User-based tracking** prevents account-level abuse

### 7.6 Prevention Recommendations

```go
// Behavioral detection
type BehavioralLimiter struct {
    redis *redis.Client
}

func (bl *BehavioralLimiter) DetectAttack(req *http.Request) bool {
    ip := getClientIP(req)
    endpoint := req.URL.Path
    key := fmt.Sprintf("behavior:%s:ips", endpoint)
    bl.redis.SAdd(context.Background(), key, ip)
    uniqueIPs := bl.redis.SCard(context.Background(), key).Val()
    return uniqueIPs > 1000
}
```

---

## Case Study 8: Async Processing Deadlock

### 8.1 Incident Description

**System**: Order fulfillment system with async workers
**Impact**: Orders stuck in "processing" state for 24 hours
**Duration**: 24 hours
**Date**: August 2023

A circular dependency between two async workers caused a deadlock. Worker A waited for Worker B, and Worker B waited for Worker A, leaving 50,000 orders stuck in an unprocessable state.

### 8.2 Root Cause Analysis

```
Circular Dependency:
Worker A: Process payment -> Wait for inventory
Worker B: Reserve inventory -> Verify payment
Result: Deadlock!
```

### 8.3 Timeline of Events

| Time (UTC) | Event | Stuck Orders |
|------------|-------|--------------|
| 10:00:00 | Deployment with new workers | 0 |
| Day 2 10:00 | Deadlock identified | 50000 |
| Day 2 11:00 | Fix deployed | 0 |

### 8.4 Resolution Steps

```go
// Acyclic workflow with event-driven architecture
type OrderWorkflow struct {
    eventBus EventBus
}

func (w *OrderWorkflow) Start(ctx context.Context, order Order) error {
    return w.eventBus.Publish(ctx, &PaymentRequestedEvent{
        OrderID: order.ID,
        Amount:  order.Total,
    })
}

// Payment service publishes event - no waiting
func (s *PaymentService) HandlePaymentRequested(ctx context.Context, event *PaymentRequestedEvent) error {
    result, err := s.processPayment(ctx, event)
    if err != nil {
        return s.eventBus.Publish(ctx, &PaymentFailedEvent{OrderID: event.OrderID})
    }
    return s.eventBus.Publish(ctx, &PaymentCompletedEvent{OrderID: event.OrderID})
}
```

### 8.5 Lessons Learned

1. **Never have circular dependencies** in async workflows
2. **Event-driven architecture** breaks direct dependencies
3. **Timeouts** must be set on all waits
4. **Deadlock detection** should monitor stuck workflows

### 8.6 Prevention Recommendations

```go
// Workflow timeout monitoring
type WorkflowTimeout struct {
    store WorkflowStore
}

func (t *WorkflowTimeout) CheckTimeouts(ctx context.Context) {
    stuckWorkflows := t.store.GetStuckWorkflows(ctx, 30*time.Minute)
    for _, workflow := range stuckWorkflows {
        alert("Workflow %s stuck, initiating compensation", workflow.ID)
    }
}
```

---

## Case Study 9: Memory Leak in Session Store

### 9.1 Incident Description

**System**: Web application with in-memory session store
**Impact**: OOM kills, session loss, user logouts
**Duration**: 3 hours
**Date**: July 2023

A memory leak in the session store implementation caused heap usage to grow unbounded. Sessions were never cleaned up, leading to OOM kills and complete session loss for all users.

### 9.2 Root Cause Analysis

```
Memory Leak Pattern:
- Session created on login
- Session accessed on each request
- Session never deleted
- No TTL on session data
- Heap grew by 100MB/hour
```

### 9.3 Timeline of Events

| Time (UTC) | Event | Memory Usage |
|------------|-------|--------------|
| 09:00:00 | Heap at 80% | 8GB/10GB |
| 10:00:00 | First OOM kill | - |
| 11:00:00 | Session store identified as cause | 9GB/10GB |
| 12:00:00 | Fix deployed with TTL | 2GB/10GB |

### 9.4 Resolution Steps

```go
// Session store with TTL and cleanup
type SessionStore struct {
    cache  *ristretto.Cache
    db     *sql.DB
}

func NewSessionStore() *SessionStore {
    cache, _ := ristretto.NewCache(&ristretto.Config{
        NumCounters: 1e7,
        MaxCost:     1 << 28, // 256MB
        BufferItems: 64,
    })

    store := &SessionStore{cache: cache}
    go store.periodicCleanup()
    return store
}

func (s *SessionStore) Get(sessionID string) (*Session, error) {
    if val, found := s.cache.Get(sessionID); found {
        return val.(*Session), nil
    }
    return s.loadFromDB(sessionID)
}

func (s *SessionStore) Set(sessionID string, session *Session, ttl time.Duration) {
    s.cache.SetWithTTL(sessionID, session, 1, ttl)
}

func (s *SessionStore) Delete(sessionID string) {
    s.cache.Del(sessionID)
    s.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
}

func (s *SessionStore) periodicCleanup() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        s.db.Exec("DELETE FROM sessions WHERE expires_at < NOW()")
    }
}
```

### 9.5 Lessons Learned

1. **Always set TTLs** on cached data
2. **Monitor heap growth** continuously
3. **Periodic cleanup** is essential
4. **Offload to persistent store** for long-term storage

### 9.6 Prevention Recommendations

```go
// Memory monitoring
type MemoryMonitor struct {
    limit uint64
}

func (m *MemoryMonitor) Start() {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            var stats runtime.MemStats
            runtime.ReadMemStats(&stats)

            if stats.HeapAlloc > m.limit {
                alert("Heap usage %.0f%% of limit", float64(stats.HeapAlloc)/float64(m.limit)*100)
            }
        }
    }()
}
```

---

## Case Study 10: Third-Party API Cascade Failure

### 10.1 Incident Description

**System**: Payment platform dependent on external fraud detection API
**Impact**: All payments blocked, $5M revenue loss
**Duration**: 4 hours
**Date**: June 2023

The external fraud detection API went down. Without proper fallback handling, all payment requests waited for the API to timeout, causing complete payment processing failure.

### 10.2 Root Cause Analysis

```
Cascade Failure:
1. Fraud API latency spiked to 30s
2. Payment service timeout: 30s
3. All payment threads blocked
4. Thread pool exhausted
5. New payment requests queued
6. Queue full, requests rejected
7. Complete payment outage
```

### 10.3 Timeline of Events

| Time (UTC) | Event | Payment Success Rate |
|------------|-------|---------------------|
| 14:00:00 | Fraud API latency spike | 80% |
| 14:15:00 | Thread pool exhaustion | 20% |
| 14:30:00 | Complete outage | 0% |
| 15:00:00 | Circuit breaker opened | 0% |
| 15:30:00 | Fallback mode activated | 60% |
| 18:00:00 | Full recovery | 100% |

### 10.4 Resolution Steps

```go
// Circuit breaker with fallback
type FraudChecker struct {
    client          *http.Client
    circuitBreaker  *gobreaker.CircuitBreaker
    fallbackChecker FallbackChecker
    timeout         time.Duration
}

func (fc *FraudChecker) CheckTransaction(ctx context.Context, tx Transaction) (*FraudResult, error) {
    result, err := fc.circuitBreaker.Execute(func() (interface{}, error) {
        ctx, cancel := context.WithTimeout(ctx, fc.timeout)
        defer cancel()
        return fc.callExternalAPI(ctx, tx)
    })

    if err != nil {
        if err == gobreaker.ErrOpenState {
            // Circuit open - use fallback
            return fc.fallbackChecker.Check(ctx, tx)
        }
        // Other error - allow transaction with elevated risk
        return &FraudResult{Risk: RiskElevated, Source: "fallback"}, nil
    }

    return result.(*FraudResult), nil
}

// Local fallback rules
type FallbackChecker struct {
    rules []FraudRule
}

func (fc *FallbackChecker) Check(ctx context.Context, tx Transaction) (*FraudResult, error) {
    risk := RiskLow

    for _, rule := range fc.rules {
        if rule.Matches(tx) {
            risk = max(risk, rule.Risk)
        }
    }

    return &FraudResult{Risk: risk, Source: "local_fallback"}, nil
}
```

### 10.5 Lessons Learned

1. **Never depend on single external API**
2. **Circuit breakers** must have fallbacks
3. **Graceful degradation** preserves core functionality
4. **Async fraud checks** can be done post-payment

### 10.6 Prevention Recommendations

```go
// Multi-provider with fallback
type MultiProviderFraudChecker struct {
    providers []FraudProvider
    timeout   time.Duration
}

func (mpc *MultiProviderFraudChecker) Check(ctx context.Context, tx Transaction) (*FraudResult, error) {
    ctx, cancel := context.WithTimeout(ctx, mpc.timeout)
    defer cancel()

    for _, provider := range mpc.providers {
        result, err := provider.Check(ctx, tx)
        if err == nil {
            return result, nil
        }
        // Log failure, try next provider
    }

    // All providers failed - use default
    return &FraudResult{Risk: RiskMedium, Source: "default"}, nil
}
```

---

## Summary and Best Practices

### Common Failure Patterns

| Pattern | Frequency | Impact | Detectability |
|---------|-----------|--------|---------------|
| Tight Coupling | High | Critical | Medium |
| Distributed Transaction Issues | Medium | Critical | Low |
| API Breaking Changes | Medium | High | Medium |
| Cache Stampede | Medium | High | High |
| Poison Messages | Medium | High | Medium |
| Connection Leaks | High | High | Medium |
| Rate Limit Bypass | Low | Critical | Low |
| Async Deadlocks | Low | Critical | Low |
| Memory Leaks | High | Medium | Low |
| Third-Party Failures | High | Critical | Medium |

### Prevention Checklist

- [ ] Use circuit breakers and bulkheads for external calls
- [ ] Implement idempotent saga compensation
- [ ] Maintain backward API compatibility
- [ ] Apply singleflight for cache stampede prevention
- [ ] Configure dead letter queues
- [ ] Monitor connection pool usage
- [ ] Implement multi-dimensional rate limiting
- [ ] Avoid circular dependencies in workflows
- [ ] Set TTLs on all cached/session data
- [ ] Design graceful degradation for third-party dependencies

### References

1. "Building Microservices" - Sam Newman
2. "Release It!" - Michael Nygard
3. "Designing Data-Intensive Applications" - Martin Kleppmann
4. "Cloud Native Patterns" - Cornelia Davis
5. "The Site Reliability Workbook" - Google

---

*Document Size: 15+ KB | Level: S | Last Updated: 2026-04-03*
