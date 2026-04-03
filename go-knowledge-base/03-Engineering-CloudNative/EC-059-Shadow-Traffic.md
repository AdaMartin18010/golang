# EC-059: Shadow Traffic Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #shadow-traffic #dark-launch #production-testing #mirroring #risk-reduction
> **Authoritative Sources**:
>
> - [Production Testing](https://landing.google.com/sre/book/chapters/testing-reliability.html) - Google SRE (2017)
| **Request Mirroring** | Mirror 100% traffic | Real load testing | Infrastructure cost |

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Shadow Traffic Domain)**
Let $R$ be a request to service $S_{production}$. Shadow traffic creates a replica $R'$ sent to $S_{shadow}$ where:

- $R'$ has identical payload to $R$
- Response from $S_{shadow}$ is discarded (fire-and-forget)
- $S_{production}$ is unaffected by $S_{shadow}$ performance

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Isolation** | $response(S_{shadow}) \nRightarrow response(client)$ | Must not affect users |
| **Fidelity** | $R' = R$ in all relevant aspects | Test must be realistic |
| **Side Effect Prevention** | $\forall side\_effect(e): execute(R') \nRightarrow e$ | No duplicate effects |
| **Resource Overhead** | $resources(S_{shadow}) \leq \theta_{budget}$ | Cost control |

### 1.2 Problem Statement

**Problem 1.1 (Safe Production Testing)**
Given new service version $S_{new}$, test under production load without risk:

$$\forall R: shadow(R, S_{new}) \Rightarrow observe(R') \land \neg affect(R')$$

**Key Challenges:**

1. **Side Effect Isolation**: Preventing duplicate writes
2. **State Synchronization**: Keeping shadow database in sync
3. **Response Comparison**: Validating correctness without affecting latency
4. **Resource Cost**: Running duplicate infrastructure
5. **Sensitive Data**: Handling PII in shadow environment

---

## 2. Solution Architecture

### 2.1 Shadow Traffic Modes

| Mode | Traffic % | Response | Use Case |
|------|-----------|----------|----------|
| **Sampled** | 1-10% | Ignored | Early validation |
| **Full** | 100% | Ignored | Load testing |
| **Comparison** | 100% | Compared | Regression testing |

---

## 3. Visual Representations

### 3.1 Shadow Traffic Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SHADOW TRAFFIC ARCHITECTURE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CLIENT                                                                     │
│  ┌─────────────┐                                                            │
│  │   Request   │                                                            │
│  └──────┬──────┘                                                            │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         TRAFFIC MIRROR                             │   │
│  │                                                                      │   │
│  │  ┌─────────────┐      ┌─────────────────────────────────────────┐   │   │
│  │  │   Clone     │      │         Request Duplication             │   │   │
│  │  │   Request   │─────►│                                         │   │   │
│  │  │             │      │  • Copy headers                         │   │   │
│  │  │  ID: abc123 │      │  • Copy body                            │   │   │
│  │  │             │      │  • Add shadow marker: X-Shadow: true    │   │   │
│  │  │  Sampling:  │      │  • Remove auth tokens (or use test)     │   │   │
│  │  │  10%        │      │                                         │   │   │
│  │  └─────────────┘      └─────────────────────────────────────────┘   │   │
│  │                                                                      │   │
│  └──────────────────────────────────┬──────────────────────────────────┘   │
│                                     │                                       │
│         ┌───────────────────────────┴───────────────────────────┐          │
│         │                           │                           │          │
│         ▼                           ▼                           ▼          │
│  ┌───────────────┐          ┌───────────────┐          ┌───────────────┐   │
│  │   PRODUCTION  │          │   SHADOW      │          │  COMPARISON   │   │
│  │               │          │   SERVICE     │          │  (Optional)   │   │
│  │  ┌─────────┐  │          │  ┌─────────┐  │          │               │   │
│  │  │ Process │  │          │  │ Process │  │─────────►│  • Compare    │   │
│  │  │ Request │  │          │  │ Request │  │          │    responses  │   │
│  │  └────┬────┘  │          │  └────┬────┘  │          │  • Log diffs  │   │
│  │       │       │          │       │       │          │  • Metrics    │   │
│  │       ▼       │          │       ▼       │          │               │   │
│  │  ┌─────────┐  │          │  ┌─────────┐  │          └───────────────┘   │
│  │  │  DB     │  │          │  │  DB     │  │                              │
│  │  │ Write   │  │          │  │ Write   │  │ (Mocked/No-op)               │
│  │  └─────────┘  │          │  └─────────┘  │                              │
│  │       │       │          │       │       │                              │
│  │       ▼       │          │       X       │ (Discarded)                  │
│  │  ┌─────────┐  │          │               │                              │
│  │  │Response │  │          │               │                              │
│  │  │ Return  │──┼──────────┼───────────────┼───────► Client               │
│  │  └─────────┘  │          │               │                              │
│  └───────────────┘          └───────────────┘                              │
│                                                                             │
│  KEY DIFFERENCES:                                                           │
│  ─────────────────                                                          │
│  Production:                                                                │
│  • Real database writes                                                     │
│  • Real external API calls                                                  │
│  • Response returned to client                                              │
│                                                                             │
│  Shadow:                                                                    │
│  • DB writes mocked/redirected to test DB                                   │
│  • External calls stubbed or sent to sandbox                                │
│  • Response discarded (fire-and-forget)                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Shadow Traffic with Side Effect Isolation

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SIDE EFFECT ISOLATION IN SHADOW                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  REQUEST FLOW:                                                              │
│                                                                             │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────────────────────┐   │
│  │   Request   │────►│   Shadow    │────►│   Service Handler           │   │
│  │             │     │   Marker    │     │                             │   │
│  └─────────────┘     └─────────────┘     └─────────────┬───────────────┘   │
│                                                        │                   │
│                                                        ▼                   │
│                                              ┌─────────────────┐           │
│                                              │  Decision Point │           │
│                                              │  (Is Shadow?)   │           │
│                                              └────────┬────────┘           │
│                                                       │                    │
│                                    ┌──────────────────┼──────────────────┐ │
│                                    │ NO               │ YES              │ │
│                                    ▼                  ▼                  │ │
│                           ┌───────────────┐  ┌────────────────┐         │ │
│                           │ Real Database │  │  Mock/No-op    │         │ │
│                           │    Write      │  │  Write         │         │ │
│                           └───────┬───────┘  └───────┬────────┘         │ │
│                                   │                  │                  │ │
│                                   ▼                  ▼                  │ │
│                           ┌───────────────┐  ┌────────────────┐         │ │
│                           │ Real External │  │  Sandbox/Stub  │         │ │
│                           │   API Call    │  │  API Call      │         │ │
│                           └───────┬───────┘  └───────┬────────┘         │ │
│                                   │                  │                  │ │
│                                   ▼                  ▼                  │ │
│                           ┌───────────────┐  ┌────────────────┐         │ │
│                           │  Return       │  │  Log & Discard │         │ │
│                           │  Response     │  │  Response      │         │ │
│                           └───────────────┘  └────────────────┘         │ │
│                                                                         │ │
└─────────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  IMPLEMENTATION STRATEGIES:                                                 │
│  ──────────────────────────                                                 │
│                                                                             │
│  1. Middleware Approach:                                                    │
│     • Intercept at HTTP/gRPC layer                                          │
│     • Check X-Shadow header                                                 │
│     • Route to mock implementations                                         │
│                                                                             │
│  2. Repository Pattern:                                                     │
│     • Abstract all DB operations                                            │
│     • Shadow mode returns success without write                             │
│     • Use interface injection for test doubles                              │
│                                                                             │
│  3. Feature Flags:                                                          │
│     • flag.isEnabled("shadow-mode")                                         │
│     • Conditional logic in business layer                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Shadow Traffic Middleware

```go
package shadow

import (
 "bytes"
 "context"
 "io"
 "net/http"
 "time"

 "go.uber.org/zap"
)

const ShadowHeader = "X-Shadow-Traffic"

// ShadowHandler handles shadow traffic
type ShadowHandler struct {
 targetURL    string
 sampler      Sampler
 transport    http.RoundTripper
 logger       *zap.Logger
}

// Sampler decides whether to shadow a request
type Sampler interface {
 ShouldShadow(r *http.Request) bool
}

// PercentageSampler samples based on percentage
type PercentageSampler struct {
 Percentage float64
}

func (s *PercentageSampler) ShouldShadow(r *http.Request) bool {
 // Hash request ID and check against percentage
 return hashPercentage(r.Header.Get("X-Request-ID")) < s.Percentage
}

// NewShadowHandler creates a new shadow handler
func NewShadowHandler(targetURL string, sampler Sampler, logger *zap.Logger) *ShadowHandler {
 return &ShadowHandler{
  targetURL: targetURL,
  sampler:   sampler,
  transport: &http.Transport{
   MaxIdleConns:        100,
   MaxIdleConnsPerHost: 10,
   IdleConnTimeout:     90 * time.Second,
  },
  logger: logger,
 }
}

// Middleware returns HTTP middleware that handles shadow traffic
func (sh *ShadowHandler) Middleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  // Process request normally
  next.ServeHTTP(w, r)

  // Decide whether to shadow
  if !sh.sampler.ShouldShadow(r) {
   return
  }

  // Fire shadow request asynchronously
  go sh.sendShadowRequest(r)
 })
}

func (sh *ShadowHandler) sendShadowRequest(original *http.Request) {
 ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
 defer cancel()

 // Clone request
 body, err := io.ReadAll(original.Body)
 if err != nil {
  sh.logger.Error("Failed to read request body for shadow", zap.Error(err))
  return
 }
 original.Body = io.NopCloser(bytes.NewReader(body))

 // Create shadow request
 shadowReq, err := http.NewRequestWithContext(ctx, original.Method,
  sh.targetURL+original.URL.Path, bytes.NewReader(body))
 if err != nil {
  sh.logger.Error("Failed to create shadow request", zap.Error(err))
  return
 }

 // Copy headers
 for key, values := range original.Header {
  for _, value := range values {
   shadowReq.Header.Add(key, value)
  }
 }

 // Mark as shadow traffic
 shadowReq.Header.Set(ShadowHeader, "true")

 // Send shadow request
 start := time.Now()
 resp, err := sh.transport.RoundTrip(shadowReq)
 duration := time.Since(start)

 if err != nil {
  sh.logger.Error("Shadow request failed",
   zap.Error(err),
   zap.Duration("duration", duration))
  return
 }
 defer resp.Body.Close()

 // Log shadow response (don't consume body to avoid memory pressure)
 sh.logger.Debug("Shadow request completed",
  zap.Int("status", resp.StatusCode),
  zap.Duration("duration", duration),
  zap.String("path", original.URL.Path))
}

func hashPercentage(id string) float64 {
 // Simple hash - use proper hash in production
 var h int64
 for _, c := range id {
  h = h*31 + int64(c)
 }
 return float64(h%10000) / 100.0
}
```

### 4.2 Side Effect Isolation

```go
package shadow

import (
 "context"
)

// Repository is the data access interface
type Repository interface {
 Create(ctx context.Context, entity interface{}) error
 Update(ctx context.Context, entity interface{}) error
 Delete(ctx context.Context, id string) error
 Get(ctx context.Context, id string) (interface{}, error)
}

// ShadowRepository wraps a repository with shadow mode support
type ShadowRepository struct {
 delegate Repository
 isShadow func(ctx context.Context) bool
}

// NewShadowRepository creates a shadow-aware repository
func NewShadowRepository(delegate Repository, isShadow func(ctx context.Context) bool) *ShadowRepository {
 return &ShadowRepository{
  delegate: delegate,
  isShadow: isShadow,
 }
}

func (r *ShadowRepository) Create(ctx context.Context, entity interface{}) error {
 if r.isShadow(ctx) {
  // In shadow mode: validate but don't persist
  // Could also write to separate test database
  return nil
 }
 return r.delegate.Create(ctx, entity)
}

func (r *ShadowRepository) Update(ctx context.Context, entity interface{}) error {
 if r.isShadow(ctx) {
  return nil
 }
 return r.delegate.Update(ctx, entity)
}

func (r *ShadowRepository) Delete(ctx context.Context, id string) error {
 if r.isShadow(ctx) {
  return nil
 }
 return r.delegate.Delete(ctx, id)
}

func (r *ShadowRepository) Get(ctx context.Context, id string) (interface{}, error) {
 // Reads are allowed in shadow mode
 return r.delegate.Get(ctx, id)
}

// IsShadowContext checks if context is for shadow traffic
func IsShadowContext(ctx context.Context) bool {
 return ctx.Value(ShadowHeader) != nil
}
```

---

## 5. Failure Scenarios and Mitigations

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Shadow Overload** | Production slowdown | Latency increase | Rate limiting |
| **Data Leak** | PII in shadow logs | Audit scan | Data masking |
| **Side Effect Bug** | Duplicate writes | Anomaly detection | Strict mocking |
| **Resource Exhaustion** | Cost overrun | Budget alert | Auto-scaling limits |

---

## 6. Semantic Trade-off Analysis

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SHADOW TRAFFIC COMPARISON                                 │
├─────────────────────┬─────────────────┬─────────────────────────────────────┤
│     Dimension       │  Shadow Traffic │         Canary Deployment           │
├─────────────────────┼─────────────────┼─────────────────────────────────────┤
│ User Impact         │ Zero            │ Minimal (small % affected)          │
│ Risk Level          │ Very Low        │ Low                                 │
│ Validation Depth    │ Deep (full load)│ Limited (sampled)                   │
│ Cost                │ 2x (duplicate)  │ 1.25x-1.5x                          │
│ Setup Complexity    │ High (mocking)  │ Medium                              │
│ Feedback Quality    │ No user metrics │ Real user feedback                  │
└─────────────────────┴─────────────────┴─────────────────────────────────────┘
```

---

## 7. References

1. Google SRE Team. (2017). *Testing for Reliability*. Google SRE Book.
2. Fowler, M. (2014). *Dark Launching*. martinfowler.com.
