# EC-010: Graceful Degradation Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (18+ KB)
> **Tags**: #graceful-degradation #fallback #circuit-breaker #feature-flags #resilience
> **Authoritative Sources**:
>
> - [Graceful Degradation](https://docs.microsoft.com/en-us/azure/architecture/patterns/gatekeeper) - Microsoft Azure
> - [Release It!](https://pragprog.com/titles/mnee2/release-it-second-edition/) - Michael Nygard (2018)
> - [Google SRE](https://sre.google/sre-book/handling-overload/) - Google
> - [Feature Toggles](https://martinfowler.com/articles/feature-toggles.html) - Martin Fowler
> - [Resilience Patterns](https://aws.amazon.com/blogs/architecture/resilience-in-distributed-systems/) - AWS

---

## 1. Pattern Overview

### 1.1 Problem Statement

When services experience failures or overload, complete failure is often unacceptable. Users expect core functionality to continue working even when non-essential features are unavailable.

**Degradation Scenarios:**

- Recommendation service down → Show popular items instead
- Payment processor slow → Queue for async processing
- Analytics service failing → Disable analytics, keep main features
- CDN unavailable → Serve content directly

### 1.2 Solution Overview

Graceful Degradation ensures:

- Core functionality remains available
- Non-critical features fail silently
- Users receive reduced but functional service
- System can recover automatically when dependencies heal

---

## 2. Design Pattern Formalization

### 2.1 Degradation Model

**Definition 2.1 (Feature Criticality)**
Each feature $f$ has a criticality level $c(f)$:

$$
c(f) \in \{\text{CRITICAL}, \text{IMPORTANT}, \text{NICE_TO_HAVE}\}
$$

**Definition 2.2 (Degradation Policy)**
A degradation policy $P$ defines fallback behavior:

$$
P: F \times S \to \{\text{full}, \text{reduced}, \text{disabled}\}
$$

Where:

- $F$: Set of features
- $S$: System state (normal, degraded, critical)

**Definition 2.3 (Graceful Degradation)**
Given a feature request $req_f$ and system state $s$:

$$
\text{GD}(req_f, s) = \begin{cases}
execute(req_f) & \text{if } availability(f, s) = 1 \\
fallback(f, s) & \text{if } 0 < availability(f, s) < 1 \\
error(f, s) & \text{if } availability(f, s) = 0 \land c(f) = \text{CRITICAL} \\
silent(f) & \text{if } availability(f, s) = 0 \land c(f) \neq \text{CRITICAL}
\end{cases}
$$

### 2.2 Degradation Levels

| Level | Description | User Impact | Implementation |
|-------|-------------|-------------|----------------|
| **Full Service** | All features available | None | Normal operation |
| **Reduced Service** | Non-critical features disabled | Minimal | Feature flags, fallbacks |
| **Core Service** | Only critical features | Moderate | Circuit breakers, caching |
| **Emergency** | Minimal survival mode | Significant | Static responses, queues |

---

## 3. Visual Representations

### 3.1 Degradation Decision Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Graceful Degradation Decision Flow                        │
└─────────────────────────────────────────────────────────────────────────────┘

Request Received
       │
       ▼
┌─────────────────┐
│ Check System    │
│ State           │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     Normal      ┌─────────────────┐
│ System State?   │────────────────►│ Full Service    │
└────────┬────────┘                 │ Execute Normal  │
         │                          │ Operation       │
         │ Degraded                 └─────────────────┘
         ▼
┌─────────────────┐     Critical    ┌─────────────────┐
│ Feature         │────────────────►│ Core Service    │
│ Criticality?    │                 │ Only Critical   │
└────────┬────────┘                 │ Features        │
         │                          └─────────────────┘
         │ Important
         ▼
┌─────────────────┐     Available   ┌─────────────────┐
│ Dependency      │────────────────►│ Reduced Service │
│ Available?      │                 │ Use Fallback    │
└────────┬────────┘                 └─────────────────┘
         │
         │ Unavailable
         ▼
┌─────────────────┐     Critical    ┌─────────────────┐
│ Use Fallback?   │────────────────►│ Try Fallback    │
└────────┬────────┘                 │ or Queue        │
         │                          └─────────────────┘
         │ No Fallback
         ▼
┌─────────────────┐     Critical    ┌─────────────────┐
│ Return Error    │                 │ Show Error to   │
│ or Silent?      │────────────────►│ User            │
└────────┬────────┘                 └─────────────────┘
         │
         │ Non-Critical
         ▼
┌─────────────────┐
│ Silent Failure  │
│ (Feature Hidden)│
└─────────────────┘
```

### 3.2 Feature Degradation Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Feature Degradation Matrix                                │
└─────────────────────────────────────────────────────────────────────────────┘

E-Commerce Application:

┌─────────────────────┬───────────┬──────────────┬──────────────┬────────────┐
│ Feature             │ Normal    │ Reduced      │ Core Only    │ Emergency  │
├─────────────────────┼───────────┼──────────────┼──────────────┼────────────┤
│ Browse Products     │ Full      │ Full         │ Full         │ Static     │
│ Search              │ Full      │ Cached       │ Basic        │ Disabled   │
│ Product Details     │ Full      │ Cached       │ Cached       │ Static     │
│ Add to Cart         │ Full      │ Full         │ Full         │ Queue      │
│ Checkout            │ Full      │ Async        │ Async        │ Disabled   │
│ Payment             │ Full      │ Stripe only  │ Stripe only  │ Disabled   │
│ Recommendations     │ ML-based  │ Popular      │ Disabled     │ Disabled   │
│ Reviews             │ Full      │ Cached       │ Disabled     │ Disabled   │
│ User Profile        │ Full      │ Read-only    │ Read-only    │ Disabled   │
│ Order History       │ Full      │ Cached       │ Cached       │ Disabled   │
│ Live Chat           │ Full      │ Bot only     │ Disabled     │ Disabled   │
│ Analytics           │ Full      │ Sampled      │ Disabled     │ Disabled   │
└─────────────────────┴───────────┴──────────────┴──────────────┴────────────┘

Fallback Strategies by Feature:

Recommendations:
  Normal:  Personalized ML recommendations
  Fallback: Trending items (pre-computed)
  Fallback2: Popular category items
  Emergency: Static featured items

Search:
  Normal:  Elasticsearch with ranking
  Fallback: Database search (slower)
  Fallback2: Cached recent results
  Emergency: Category browse only

Payment:
  Normal:  Multiple processors (Stripe, PayPal, etc.)
  Fallback: Primary processor only
  Fallback2: Queue for async processing
  Emergency: "Try again later" message
```

### 3.3 Degradation Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Graceful Degradation Architecture                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Request                                  │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Degradation Router                                   │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  • System Load Monitor                                                │  │
│  │  • Dependency Health Checker                                          │  │
│  │  • Feature Flag Evaluator                                             │  │
│  │  • Circuit Breaker State                                              │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
              ┌─────────────────────┼─────────────────────┐
              │                     │                     │
              ▼                     ▼                     ▼
    ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
    │  Normal Path    │   │ Degraded Path   │   │  Fallback Path  │
    │                 │   │                 │   │                 │
    │ ┌─────────────┐ │   │ ┌─────────────┐ │   │ ┌─────────────┐ │
    │ │  Primary    │ │   │ │  Reduced    │ │   │ │  Cached/    │ │
    │ │  Service    │ │   │ │  Service    │ │   │ │  Static     │ │
    │ └─────────────┘ │   │ └─────────────┘ │   │ └─────────────┘ │
    └────────┬────────┘   └────────┬────────┘   └────────┬────────┘
             │                     │                     │
             ▼                     ▼                     ▼
    ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
    │  Full Response  │   │  Partial Resp   │   │  Minimal Resp   │
    └─────────────────┘   └─────────────────┘   └─────────────────┘

Fallback Cache Layer:
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                  │
│  │   Hot Cache  │    │  Warm Cache  │    │  Cold Cache  │                  │
│  │  (In-Memory) │    │   (Redis)    │    │  (Static)    │                  │
│  │              │    │              │    │              │                  │
│  │ • User prefs │    │ • Top items  │    │ • Default    │                  │
│  │ • Sessions   │    │ • Search idx │    │ • Home page  │                  │
│  │ • Cart       │    │ • Analytics  │    │ • Error pages│                  │
│  └──────────────┘    └──────────────┘    └──────────────┘                  │
│                                                                              │
│  Stale Data Tolerance:                                                       │
│  • Hot:  0-5 minutes (high freshness)                                        │
│  • Warm: 5-60 minutes (acceptable staleness)                                 │
│  • Cold: Infinite (static fallback)                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

Queue for Async Processing:
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  Request ──► ┌─────────┐ ──► ┌─────────┐ ──► ┌─────────┐ ──► Success       │
│              │ Accept  │     │  Queue  │     │ Process │                    │
│              │         │     │         │     │ Later   │                    │
│              └─────────┘     └─────────┘     └─────────┘                    │
│                    │                               │                        │
│                    ▼                               ▼                        │
│              Return 202                      Notify User                     │
│              (Accepted)                      (Webhook/Email)                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production-Ready Implementation

### 4.1 Degradation Manager

```go
package degradation

import (
 "context"
 "errors"
 "sync"
 "time"

 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/metric"
)

// Criticality level
type Criticality int

const (
 CriticalityCritical Criticality = iota
 CriticalityImportant
 CriticalityNiceToHave
)

func (c Criticality) String() string {
 switch c {
 case CriticalityCritical:
  return "critical"
 case CriticalityImportant:
  return "important"
 case CriticalityNiceToHave:
  return "nice_to_have"
 default:
  return "unknown"
 }
}

// Feature defines a degradable feature
type Feature struct {
 Name        string
 Criticality Criticality
 Execute     func(ctx context.Context) (interface{}, error)
 Fallback    func(ctx context.Context) (interface{}, error)
 SilentFail  bool
}

// Manager handles graceful degradation
type Manager struct {
 features map[string]*Feature
 state    SystemState
 mu       sync.RWMutex

 // Metrics
 degradationCounter metric.Int64Counter
 fallbackCounter    metric.Int64Counter
}

// SystemState represents system health state
type SystemState int

const (
 StateNormal SystemState = iota
 StateDegraded
 StateCritical
 StateEmergency
)

// NewManager creates a degradation manager
func NewManager(meter metric.Meter) *Manager {
 m := &Manager{
  features: make(map[string]*Feature),
  state:    StateNormal,
 }

 if meter != nil {
  var err error
  m.degradationCounter, err = meter.Int64Counter(
   "degradation_events_total",
   metric.WithDescription("Total degradation events"),
  )
  if err != nil {
   // Log error
  }

  m.fallbackCounter, err = meter.Int64Counter(
   "fallback_executions_total",
   metric.WithDescription("Total fallback executions"),
  )
  if err != nil {
   // Log error
  }
 }

 return m
}

// Register registers a feature
func (m *Manager) Register(feature *Feature) {
 m.mu.Lock()
 defer m.mu.Unlock()
 m.features[feature.Name] = feature
}

// SetState updates system state
func (m *Manager) SetState(state SystemState) {
 m.mu.Lock()
 defer m.mu.Unlock()
 m.state = state
}

// Execute executes a feature with degradation
func (m *Manager) Execute(ctx context.Context, featureName string) (interface{}, error) {
 m.mu.RLock()
 feature, ok := m.features[featureName]
 state := m.state
 m.mu.RUnlock()

 if !ok {
  return nil, errors.New("feature not found")
 }

 // Check if feature should be executed based on state
 if !m.shouldExecute(feature, state) {
  if feature.SilentFail {
   return nil, nil
  }
  return nil, errors.New("feature disabled due to system state")
 }

 // Try primary execution
 result, err := feature.Execute(ctx)
 if err == nil {
  return result, nil
 }

 // Record degradation
 if m.degradationCounter != nil {
  m.degradationCounter.Add(ctx, 1, metric.WithAttributes(
   attribute.String("feature", featureName),
   attribute.String("reason", "primary_failed"),
  ))
 }

 // Try fallback
 if feature.Fallback != nil {
  if m.fallbackCounter != nil {
   m.fallbackCounter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("feature", featureName),
   ))
  }
  return feature.Fallback(ctx)
 }

 // Handle based on criticality
 if feature.Criticality == CriticalityNiceToHave {
  if feature.SilentFail {
   return nil, nil
  }
 }

 return nil, err
}

func (m *Manager) shouldExecute(feature *Feature, state SystemState) bool {
 switch state {
 case StateEmergency:
  return feature.Criticality == CriticalityCritical
 case StateCritical:
  return feature.Criticality != CriticalityNiceToHave
 case StateDegraded:
  return true
 default:
  return true
 }
}
```

### 4.2 Fallback Cache

```go
package degradation

import (
 "context"
 "encoding/json"
 "time"

 "github.com/redis/go-redis/v9"
)

// FallbackCache provides cached fallbacks
type FallbackCache struct {
 redis      *redis.Client
 localCache map[string]cacheEntry
 ttl        time.Duration
}

type cacheEntry struct {
 value     interface{}
 expireAt  time.Time
}

// NewFallbackCache creates a fallback cache
func NewFallbackCache(redis *redis.Client, ttl time.Duration) *FallbackCache {
 return &FallbackCache{
  redis:      redis,
  localCache: make(map[string]cacheEntry),
  ttl:        ttl,
 }
}

// Get retrieves cached value
func (fc *FallbackCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
 // Check local cache first
 if entry, ok := fc.localCache[key]; ok {
  if time.Now().Before(entry.expireAt) {
   return true, fc.copyValue(entry.value, dest)
  }
  delete(fc.localCache, key)
 }

 // Check Redis
 if fc.redis != nil {
  data, err := fc.redis.Get(ctx, key).Bytes()
  if err == nil {
   if err := json.Unmarshal(data, dest); err == nil {
    // Update local cache
    fc.localCache[key] = cacheEntry{
     value:    dest,
     expireAt: time.Now().Add(fc.ttl),
    }
    return true, nil
   }
  }
 }

 return false, nil
}

// Set stores value in cache
func (fc *FallbackCache) Set(ctx context.Context, key string, value interface{}) error {
 // Update local cache
 fc.localCache[key] = cacheEntry{
  value:    value,
  expireAt: time.Now().Add(fc.ttl),
 }

 // Store in Redis
 if fc.redis != nil {
  data, err := json.Marshal(value)
  if err != nil {
   return err
  }
  return fc.redis.Set(ctx, key, data, fc.ttl).Err()
 }

 return nil
}

func (fc *FallbackCache) copyValue(src, dest interface{}) error {
 data, err := json.Marshal(src)
 if err != nil {
  return err
 }
 return json.Unmarshal(data, dest)
}
```

### 4.3 Static Fallback Generator

```go
package degradation

import (
 "encoding/json"
 "time"
)

// StaticFallback provides static fallback data
type StaticFallback struct {
 data map[string]interface{}
}

// NewStaticFallback creates static fallback provider
func NewStaticFallback() *StaticFallback {
 return &StaticFallback{
  data: make(map[string]interface{}),
 }
}

// Register registers static fallback data
func (sf *StaticFallback) Register(key string, data interface{}) {
 sf.data[key] = data
}

// Get retrieves static fallback
func (sf *StaticFallback) Get(key string, dest interface{}) (bool, error) {
 data, ok := sf.data[key]
 if !ok {
  return false, nil
 }

 // Copy via JSON
 jsonData, err := json.Marshal(data)
 if err != nil {
  return false, err
 }

 return true, json.Unmarshal(jsonData, dest)
}

// Common Static Fallbacks
func DefaultStaticFallbacks() *StaticFallback {
 sf := NewStaticFallback()

 // Home page data
 sf.Register("home_page", map[string]interface{}{
  "featured_products": []map[string]interface{}{
   {"id": 1, "name": "Popular Item 1", "price": 29.99},
   {"id": 2, "name": "Popular Item 2", "price": 49.99},
  },
  "categories": []string{"Electronics", "Books", "Home"},
 })

 // Search results
 sf.Register("search_empty", map[string]interface{}{
  "query":   "",
  "results": []interface{}{},
  "message": "Search temporarily unavailable. Please browse categories instead.",
 })

 // User profile (minimal)
 sf.Register("user_profile_minimal", map[string]interface{}{
  "id":       0,
  "name":     "Guest",
  "email":    "",
  "is_guest": true,
 })

 return sf
}
```

---

## 5. Failure Scenarios and Mitigation

| Scenario | Symptom | Cause | Mitigation |
|----------|---------|-------|------------|
| **Hidden Failures** | Silent data loss | Silent fail on critical feature | Criticality classification, alerts |
| **Stale Data** | Users see old data | Over-reliance on cache | TTL limits, freshness indicators |
| **Cascading Degradation** | All features disabled | Dependency chain failure | Circuit breakers, bulkheads |
| **Recovery Failure** | System stays degraded | No auto-recovery | Health checks, state monitoring |

---

## 6. Observability Integration

```go
// DegradationMetrics for monitoring
type DegradationMetrics struct {
 stateGauge         metric.Int64Gauge
 degradationCounter metric.Int64Counter
 fallbackLatency    metric.Float64Histogram
}
```

---

## 7. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Degradation Security Checklist                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Data Consistency:                                                           │
│  □ Ensure fallback data doesn't expose sensitive information                 │
│  □ Validate cached data integrity                                            │
│  □ Don't serve stale authentication data                                     │
│                                                                              │
│  Access Control:                                                             │
│  □ Maintain access controls in degraded mode                                 │
│  □ Don't bypass authentication in fallback paths                             │
│  □ Validate permissions in static fallbacks                                  │
│                                                                              │
│  Denial of Service:                                                          │
│  □ Prevent degradation from being triggered maliciously                      │
│  □ Rate limit degradation transitions                                        │
│  □ Monitor for abuse patterns                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Best Practices

### 8.1 Feature Classification

| Feature | Criticality | Fallback | Max Staleness |
|---------|-------------|----------|---------------|
| Login | Critical | None | N/A |
| Browse | Critical | Static | 1 hour |
| Search | Important | Cached | 15 min |
| Recommendations | Nice-to-have | Popular | 1 hour |
| Reviews | Nice-to-have | Disabled | N/A |
| Analytics | Nice-to-have | Sampled | 24 hours |

---

## 9. References

1. **Nygard, M. T.** *Release It!* Pragmatic Bookshelf.
2. **Fowler, M.** [Feature Toggles](https://martinfowler.com/articles/feature-toggles.html).
3. **Google**. [Handling Overload](https://sre.google/sre-book/handling-overload/).

---

**Quality Rating**: S (18KB+, Complete Formalization + Production Code + Visualizations)
