# EC-052: Health Endpoint Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #health-checks #kubernetes #observability #load-balancing #graceful-degradation
> **Authoritative Sources**:
>
> - [Kubernetes Health Checks](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) - Kubernetes (2024)
> - [Microservices Health Check](https://microservices.io/patterns/observability/health-check-api.html) - Richardson (2020)
> - [AWS Well-Architected Framework](https://docs.aws.amazon.com/wellarchitected/latest/framework/) - AWS (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Service Health Domain)**
Let $S$ be a service with state space $\Sigma = \{\sigma_1, \sigma_2, ..., \sigma_n\}$ where:

- Each state $\sigma_i$ represents a configuration of component health
- Health is a predicate $H: \Sigma \to \{healthy, unhealthy, degraded\}$
- Observers $O = \{O_1, O_2, ...\}$ need to query health for routing decisions

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Availability Requirement** | $P(H(S) = healthy) \geq 99.99\%$ | Health checks must be accurate |
| **Response Time SLA** | $T_{health\_check} < 100ms$ | Cannot impact normal operations |
| **Resource Isolation** | $resources(check) \cap resources(S) = \emptyset$ | Check must not affect service |
| **Cascade Avoidance** | $\forall C \in dependencies(S): H(C) \nRightarrow H(S)$ | Deep checks must be bounded |

### 1.2 Problem Statement

**Problem 1.1 (Health Assessment Problem)**
Given service $S$ with dependencies $D(S) = \{d_1, d_2, ..., d_m\}$, determine health $H(S)$ such that:

$$H(S) = f(h_{local}(S), h_{cascade}(D(S)))$$

Where $h_{local}$ assesses internal state and $h_{cascade}$ aggregates dependency health with depth limit.

**Key Challenges:**

1. **Accuracy vs Performance**: Deep health checks provide accuracy but add latency
2. **False Positives**: Transient issues triggering unnecessary failovers
3. **False Negatives**: Delayed detection of actual failures
4. **Resource Exhaustion**: Aggressive health check frequency consuming resources
5. **Dependency Cascade**: Excessive dependency checking amplifying failures

### 1.3 Formal Requirements Specification

**Requirement 1.1 (Accuracy)**
$$P(H_{reported}(S) = H_{actual}(S)) > 0.99$$

**Requirement 1.2 (Latency)**
$$T_{response}(health\_endpoint) < T_{SLA}$$

**Requirement 1.3 (Isolation)**
$$failure(health\_check) \nRightarrow failure(S)$$

---

## 2. Solution Architecture

### 2.1 Health Check Taxonomy

| Check Type | Purpose | Kubernetes Mapping | Frequency |
|------------|---------|-------------------|-----------|
| **Liveness** | Detect deadlock/stuck state | livenessProbe | Every 10s |
| **Readiness** | Determine traffic eligibility | readinessProbe | Every 5s |
| **Startup** | Detect initialization complete | startupProbe | During startup |
| **Deep** | Dependency verification | Custom endpoint | On-demand |

### 2.2 Health State Machine

```
States:
- UNKNOWN: Initial state, no checks performed
- STARTING: Application initializing
- HEALTHY: All checks passing
- DEGRADED: Functioning with reduced capacity
- UNHEALTHY: Critical check failing

Transitions:
UNKNOWN ──start──► STARTING ──init_complete──► HEALTHY
                                         │
                                         ├──dependency_failure──► DEGRADED
                                         │
                                         └──critical_failure──► UNHEALTHY
```

---

## 3. Visual Representations

### 3.1 Health Check Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HEALTH CHECK ARCHITECTURE                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         LOAD BALANCER                                │   │
│  │                                                                      │   │
│  │  Health Check: GET /health/ready                                     │   │
│  │  Interval: 5s  Timeout: 3s  Threshold: 2 failures                     │   │
│  │                                                                      │   │
│  │  Routing Decision:                                                   │   │
│  │  • HTTP 200 → Route traffic                                          │   │
│  │  • HTTP 503 → Remove from pool                                       │   │
│  └──────────────────────────────────┬───────────────────────────────────┘   │
│                                     │                                       │
│         ┌───────────────────────────┼───────────────────────────┐          │
│         │                           │                           │          │
│         ▼                           ▼                           ▼          │
│  ┌───────────────┐          ┌───────────────┐          ┌───────────────┐   │
│  │   Instance 1  │          │   Instance 2  │          │   Instance 3  │   │
│  │   (Healthy)   │          │  (Degraded)   │          │  (Unhealthy)  │   │
│  │               │          │               │          │               │   │
│  │  /health/live │          │  /health/live │          │  /health/live │   │
│  │  → HTTP 200   │          │  → HTTP 200   │          │  → HTTP 500   │   │
│  │               │          │               │          │               │   │
│  │  /health/ready│          │  /health/ready│          │  /health/ready│   │
│  │  → HTTP 200   │          │  → HTTP 503   │          │  → HTTP 503   │   │
│  │               │          │  (DB slow)    │          │  (App crash)  │   │
│  └───────┬───────┘          └───────┬───────┘          └───────┬───────┘   │
│          │                          │                          │           │
│          │ Receive traffic          │ No new traffic           │ No traffic│
│          │                          │ Existing continues         │           │
│          ▼                          ▼                          ▼           │
│  ┌───────────────┐          ┌───────────────┐          ┌───────────────┐   │
│  │   Traffic     │          │   Traffic     │          │   Traffic     │   │
│  │   Active      │          │   Draining    │          │   Rejected    │   │
│  └───────────────┘          └───────────────┘          └───────────────┘   │
│                                                                             │
│  LEGEND:                                                                    │
│  ───────                                                                    │
│  ● /health/live  - Liveness (is the process running?)                       │
│  ○ /health/ready - Readiness (is it ready to serve traffic?)                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Kubernetes Probe Integration

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    KUBERNETES HEALTH PROBES                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  POD LIFECYCLE WITH PROBES:                                                 │
│                                                                             │
│  Time ─────────────────────────────────────────────────────────────────►   │
│                                                                             │
│  Scheduled                                                                  │
│     │                                                                       │
│     ▼                                                                       │
│  ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐          │
│  │  Pull    │────►│  Start   │────►│ Startup  │────►│  Ready   │          │
│  │  Image   │     │Container │     │  Probe   │     │  State   │          │
│  └──────────┘     └──────────┘     └────┬─────┘     └────┬─────┘          │
│                                         │                │                │
│                                         │ Success        │ Readiness      │
│                                         │                │ Probe          │
│                                    ┌────┴────┐      ┌────┴────┐           │
│                                    │ Liveness│      │ Receive │           │
│                                    │ Probe   │      │ Traffic │           │
│                                    │ (10s)   │      │         │           │
│                                    └────┬────┘      └────┬────┘           │
│                                         │                │                │
│                    ┌────────────────────┘                │                │
│                    │ Fail                                │ Fail           │
│                    ▼                                     ▼                │
│               ┌──────────┐                          ┌──────────┐          │
│               │ Restart  │                          │ Remove   │          │
│               │Container │                          │ from LB  │          │
│               └──────────┘                          └──────────┘          │
│                                                                             │
│  PROBE CONFIGURATION:                                                       │
│  ───────────────────                                                        │
│                                                                             │
│  startupProbe:                                                              │
│    httpGet:                                                                 │
│      path: /health/startup                                                  │
│      port: 8080                                                             │
│    failureThreshold: 30      # 30 * 10s = 5 min to start                    │
│    periodSeconds: 10                                                        │
│                                                                             │
│  livenessProbe:                                                             │
│    httpGet:                                                                 │
│      path: /health/live                                                     │
│      port: 8080                                                             │
│    initialDelaySeconds: 0                                                   │
│    periodSeconds: 10                                                        │
│    failureThreshold: 3       # Restart after 30s of failures                │
│                                                                             │
│  readinessProbe:                                                            │
│    httpGet:                                                                 │
│      path: /health/ready                                                    │
│      port: 8080                                                             │
│    initialDelaySeconds: 0                                                   │
│    periodSeconds: 5                                                         │
│    failureThreshold: 3       # Remove from LB after 15s                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Health Check Component Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HEALTH CHECK COMPONENT HIERARCHY                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  OVERALL HEALTH (/health)                                                   │
│  ═══════════════════════                                                    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Status: DEGRADED                                                    │   │
│  │  {                                                                   │   │
│  │    "status": "degraded",                                             │   │
│  │    "checks": {                                                       │   │
│  │      "disk":    {"status": "healthy"},                               │   │
│  │      "memory":  {"status": "healthy"},                               │   │
│  │      "database": {"status": "degraded", "response_time": "2.5s"},     │   │
│  │      "cache":   {"status": "healthy"}                                │   │
│  │    }                                                                 │   │
│  │  }                                                                   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│         ┌──────────────────────────┼──────────────────────────┐            │
│         │                          │                          │            │
│         ▼                          ▼                          ▼            │
│  ┌───────────────┐          ┌───────────────┐          ┌───────────────┐   │
│  │   SYSTEM      │          │   DATABASE    │          │    CACHE      │   │
│  │   CHECKS      │          │    CHECK      │          │    CHECK      │   │
│  │               │          │               │          │               │   │
│  │ • Disk space  │          │ • Connection  │          │ • Connection  │   │
│  │ • Memory      │          │ • Query perf  │          │ • Latency     │   │
│  │ • Goroutines  │          │ • Replication │          │ • Hit ratio   │   │
│  └───────────────┘          └───────┬───────┘          └───────────────┘   │
│                                     │                                       │
│                                     ▼                                       │
│                           ┌───────────────────┐                            │
│                           │  Deep Health      │                            │
│                           │  (Optional)       │                            │
│                           │                   │                            │
│                           │ • Dependency check│                            │
│                           │ • End-to-end test │                            │
│                           └───────────────────┘                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Core Health Check Framework

```go
package health

import (
 "context"
 "encoding/json"
 "fmt"
 "net/http"
 "sync"
 "time"

 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/metric"
 "go.uber.org/zap"
)

// Status represents the health status
type Status string

const (
 StatusHealthy   Status = "healthy"
 StatusDegraded  Status = "degraded"
 StatusUnhealthy Status = "unhealthy"
 StatusUnknown   Status = "unknown"
)

// CheckResult holds the result of a health check
type CheckResult struct {
 Name        string                 `json:"name"`
 Status      Status                 `json:"status"`
 Message     string                 `json:"message,omitempty"`
 Timestamp   time.Time              `json:"timestamp"`
 Duration    time.Duration          `json:"duration_ms"`
 Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// OverallHealth represents the aggregated health status
type OverallHealth struct {
 Status    string                 `json:"status"`
 Timestamp time.Time              `json:"timestamp"`
 Version   string                 `json:"version"`
 Service   string                 `json:"service"`
 Checks    map[string]CheckResult `json:"checks"`
}

// Checker defines the interface for health checks
type Checker interface {
 Name() string
 Check(ctx context.Context) CheckResult
}

// CheckFunc is a function that implements the Checker interface
type CheckFunc struct {
 name string
 fn   func(ctx context.Context) CheckResult
}

func (c *CheckFunc) Name() string {
 return c.name
}

func (c *CheckFunc) Check(ctx context.Context) CheckResult {
 return c.fn(ctx)
}

// Registry manages health checks
type Registry struct {
 checkers []Checker
 mu       sync.RWMutex

 logger *zap.Logger
 meter  metric.Meter

 // Metrics
 healthStatus  metric.Int64Gauge
 checkDuration metric.Float64Histogram
}

// NewRegistry creates a new health check registry
func NewRegistry(logger *zap.Logger, meter metric.Meter) *Registry {
 r := &Registry{
  checkers: make([]Checker, 0),
  logger:   logger,
  meter:    meter,
 }

 if meter != nil {
  var err error
  r.healthStatus, err = meter.Int64Gauge(
   "health_status",
   metric.WithDescription("Current health status (0=unknown, 1=healthy, 2=degraded, 3=unhealthy)"),
  )
  if err != nil {
   logger.Error("Failed to create health status gauge", zap.Error(err))
  }

  r.checkDuration, err = meter.Float64Histogram(
   "health_check_duration_seconds",
   metric.WithDescription("Duration of health checks"),
  )
  if err != nil {
   logger.Error("Failed to create check duration histogram", zap.Error(err))
  }
 }

 return r
}

// Register adds a health check to the registry
func (r *Registry) Register(checker Checker) {
 r.mu.Lock()
 defer r.mu.Unlock()
 r.checkers = append(r.checkers, checker)
}

// RegisterFunc registers a function as a health check
func (r *Registry) RegisterFunc(name string, fn func(ctx context.Context) CheckResult) {
 r.Register(&CheckFunc{name: name, fn: fn})
}

// Check runs all registered health checks
func (r *Registry) Check(ctx context.Context) OverallHealth {
 r.mu.RLock()
 checkers := make([]Checker, len(r.checkers))
 copy(checkers, r.checkers)
 r.mu.RUnlock()

 checks := make(map[string]CheckResult, len(checkers))
 overallStatus := StatusHealthy

 for _, checker := range checkers {
  start := time.Now()
  result := checker.Check(ctx)
  result.Duration = time.Since(start)
  result.Timestamp = time.Now().UTC()

  checks[checker.Name()] = result

  // Aggregate status
  if result.Status == StatusUnhealthy {
   overallStatus = StatusUnhealthy
  } else if result.Status == StatusDegraded && overallStatus == StatusHealthy {
   overallStatus = StatusDegraded
  }

  // Record metrics
  if r.checkDuration != nil {
   r.checkDuration.Record(ctx, result.Duration.Seconds(),
    metric.WithAttributes(attribute.String("check_name", checker.Name())))
  }
 }

 if r.healthStatus != nil {
  statusValue := map[Status]int64{
   StatusUnknown:   0,
   StatusHealthy:   1,
   StatusDegraded:  2,
   StatusUnhealthy: 3,
  }[overallStatus]
  r.healthStatus.Record(ctx, statusValue)
 }

 return OverallHealth{
  Status:    string(overallStatus),
  Timestamp: time.Now().UTC(),
  Version:   "1.0.0", // Should be injected from build info
  Service:   "app-service",
  Checks:    checks,
 }
}

// HTTPHandler returns an HTTP handler for health checks
func (r *Registry) HTTPHandler() http.HandlerFunc {
 return func(w http.ResponseWriter, req *http.Request) {
  ctx := req.Context()
  path := req.URL.Path

  var health OverallHealth
  var statusCode int

  switch path {
  case "/health", "/health/ready":
   health = r.Check(ctx)
   if health.Status == string(StatusHealthy) {
    statusCode = http.StatusOK
   } else if health.Status == string(StatusDegraded) {
    statusCode = http.StatusServiceUnavailable // 503
   } else {
    statusCode = http.StatusServiceUnavailable
   }

  case "/health/live":
   // Simple liveness check - just check if we can respond
   health = OverallHealth{
    Status:    string(StatusHealthy),
    Timestamp: time.Now().UTC(),
    Version:   "1.0.0",
    Service:   "app-service",
    Checks: map[string]CheckResult{
     "alive": {
      Name:      "alive",
      Status:    StatusHealthy,
      Timestamp: time.Now().UTC(),
     },
    },
   }
   statusCode = http.StatusOK

  case "/health/startup":
   // Startup check - verifies initialization complete
   health = r.Check(ctx)
   if health.Status != string(StatusUnhealthy) {
    statusCode = http.StatusOK
   } else {
    statusCode = http.StatusServiceUnavailable
   }

  default:
   http.NotFound(w, req)
   return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  json.NewEncoder(w).Encode(health)
 }
}
```

### 4.2 Built-in Health Checks

```go
package health

import (
 "context"
 "database/sql"
 "fmt"
 "runtime"
 "time"
)

// DatabaseCheck creates a database health checker
func DatabaseCheck(db *sql.DB, name string) Checker {
 return &CheckFunc{
  name: name,
  fn: func(ctx context.Context) CheckResult {
   start := time.Now()
   err := db.PingContext(ctx)
   duration := time.Since(start)

   if err != nil {
    return CheckResult{
     Name:    name,
     Status:  StatusUnhealthy,
     Message: fmt.Sprintf("Database ping failed: %v", err),
     Duration: duration,
    }
   }

   // Check connection pool stats
   stats := db.Stats()
   status := StatusHealthy
   message := "Database connection healthy"

   // Degraded if high connection usage
   if stats.InUse > int(float64(stats.MaxOpenConnections)*0.8) {
    status = StatusDegraded
    message = "High connection pool usage"
   }

   return CheckResult{
    Name:    name,
    Status:  status,
    Message: message,
    Duration: duration,
    Metadata: map[string]interface{}{
     "open_connections":    stats.OpenConnections,
     "in_use":              stats.InUse,
     "idle":                stats.Idle,
     "wait_count":          stats.WaitCount,
     "max_open_connections": stats.MaxOpenConnections,
    },
   }
  },
 }
}

// MemoryCheck creates a memory health checker
func MemoryCheck(threshold float64) Checker {
 return &CheckFunc{
  name: "memory",
  fn: func(ctx context.Context) CheckResult {
   var m runtime.MemStats
   runtime.ReadMemStats(&m)

   // Calculate memory usage percentage
   used := float64(m.Sys)
   // This is a simplified check - in production, compare to container limits

   status := StatusHealthy
   message := "Memory usage normal"

   // Check if GC pressure is high
   if m.NumGC > 100 && m.PauseNs[(m.NumGC+255)%256] > 1e6 {
    status = StatusDegraded
    message = "High GC pressure detected"
   }

   return CheckResult{
    Name:   "memory",
    Status: status,
    Message: message,
    Metadata: map[string]interface{}{
     "alloc_mb":       m.Alloc / 1024 / 1024,
     "sys_mb":         m.Sys / 1024 / 1024,
     "num_gc":         m.NumGC,
     "gc_pause_us":    m.PauseNs[(m.NumGC+255)%256] / 1000,
     "goroutines":     runtime.NumGoroutine(),
    },
   }
  },
 }
}

// DiskCheck creates a disk space health checker
func DiskCheck(path string, warningThreshold, criticalThreshold float64) Checker {
 return &CheckFunc{
  name: "disk",
  fn: func(ctx context.Context) CheckResult {
   // This is a simplified version - use syscall.Statfs in production
   return CheckResult{
    Name:   "disk",
    Status: StatusHealthy,
    Message: fmt.Sprintf("Disk check for %s", path),
    Metadata: map[string]interface{}{
     "path": path,
    },
   }
  },
 }
}

// ExternalDependencyCheck creates a check for external services
func ExternalDependencyCheck(name string, checkFunc func(ctx context.Context) error) Checker {
 return &CheckFunc{
  name: name,
  fn: func(ctx context.Context) CheckResult {
   // Use timeout to prevent hanging
   ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
   defer cancel()

   start := time.Now()
   err := checkFunc(ctx)
   duration := time.Since(start)

   if err != nil {
    return CheckResult{
     Name:     name,
     Status:   StatusDegraded,
     Message:  fmt.Sprintf("Dependency check failed: %v", err),
     Duration: duration,
    }
   }

   return CheckResult{
    Name:     name,
    Status:   StatusHealthy,
    Message:  "Dependency healthy",
    Duration: duration,
   }
  },
 }
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Failure Taxonomy

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **False Positive** | Unnecessary restart | Alert on restart rate | Threshold tuning |
| **False Negative** | Delayed failure detection | Outage duration | Multiple check types |
| **Check Storm** | Resource exhaustion | CPU/memory spikes | Check caching + Rate limiting |
| **Dependency Cascade** | All services unhealthy | Dependency depth | Shallow checks by default |
| **Startup Timeout** | Container never ready | Startup duration | Separate startup probe |

---

## 6. Semantic Trade-off Analysis

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HEALTH CHECK DEPTH TRADE-OFFS                             │
├─────────────────────┬─────────────────┬─────────────────┬───────────────────┤
│     Dimension       │    Shallow      │     Medium      │       Deep        │
├─────────────────────┼─────────────────┼─────────────────┼────────────────────┤
│ Accuracy            │ ⚠️  May miss    │ ✅ Good         │ ✅ Excellent      │
│ Response Time       │ ✅ < 10ms       │ ⚠️  ~50ms       │ ❌ > 100ms        │
│ Resource Usage      │ ✅ Minimal      │ ⚠️  Moderate    │ ❌ High           │
│ Dependency Impact   │ ✅ None         │ ⚠️  Some        │ ❌ High risk      │
│ False Positive Rate │ ⚠️  Higher      │ ✅ Lower        │ ✅ Lowest         │
├─────────────────────┼─────────────────┼─────────────────┼────────────────────┤
│ Use for Liveness    │ ✅ Yes          │ ❌ No           │ ❌ No             │
│ Use for Readiness   │ ⚠️  Partial     │ ✅ Yes          │ ⚠️  Carefully     │
│ Use for Deep Check  │ ❌ No           │ ⚠️  Maybe       │ ✅ Yes            │
└─────────────────────┴─────────────────┴─────────────────┴───────────────────┘
```

---

## 7. References

1. Kubernetes Documentation. (2024). *Configure Liveness, Readiness and Startup Probes*.
2. Richardson, C. (2020). *Health Check API Pattern*. microservices.io.
3. AWS. (2024). *Well-Architected Framework - Operational Excellence*.
4. Google SRE Team. (2017). *Monitoring Distributed Systems*. Google SRE Book.
