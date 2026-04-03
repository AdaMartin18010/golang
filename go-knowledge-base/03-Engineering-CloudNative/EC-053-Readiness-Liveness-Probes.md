# EC-053: Readiness and Liveness Probes Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #kubernetes #health-probes #container-lifecycle #orchestration #resilience
> **Authoritative Sources**:
>
> - [Kubernetes Probes](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes) - Kubernetes (2024)
> - [Container Lifecycle Hooks](https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/) - Kubernetes (2024)
> - [Graceful Shutdown](https://cloud.google.com/blog/topics/containers-kubernetes/kubernetes-best-practices-terminating-with-grace) - Google Cloud (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Container Lifecycle Domain)**
Let $C$ be a container with lifecycle states $\Sigma = \{Pending, Running, Succeeded, Failed, Unknown\}$. An orchestrator $O$ manages $C$ through transitions triggered by probes $P = \{liveness, readiness, startup\}$.

**Probe Types:**

| Probe | Purpose | Action on Failure | Timing |
|-------|---------|-------------------|--------|
| **Startup** | Detect initialization complete | N/A (blocks others) | Until success |
| **Liveness** | Detect dead/stuck state | Restart container | Continuous |
| **Readiness** | Determine traffic eligibility | Remove from service | Continuous |

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Startup Time Variance** | $\sigma_{startup} > \mu_{startup}$ | Fixed delays are suboptimal |
| **False Positive Cost** | $cost(restart) \gg cost(check)$ | Incorrect restarts hurt availability |
| **Resource Competition** | $resources(check) \cap resources(app) \neq \emptyset$ | Health checks affect performance |
| **Probe Isolation** | $failure(check) \nRightarrow failure(app)$ | Check errors shouldn't cascade |

### 1.2 Problem Statement

**Problem 1.1 (Probe Design Problem)**
Given container $C$ with initialization time $T_{init} \sim \mathcal{N}(\mu, \sigma^2)$, design probe configuration $\pi = \{type, interval, timeout, threshold\}$ such that:

$$\forall t: action(O, \pi, C(t)) = optimal\_action(C(t))$$

**Key Challenges:**

1. **Startup Detection**: Determining when initialization is complete vs. failed
2. **Deadlock Detection**: Distinguishing slow processing from stuck processes
3. **False Trigger Avoidance**: Preventing premature restarts
4. **Resource Efficiency**: Minimizing probe overhead
5. **Graceful Degradation**: Handling partial failures

---

## 2. Solution Architecture

### 2.1 Probe State Machine

```
States:
- PENDING: Container created, probes not started
- STARTING: Startup probe running
- RUNNING_HEALTHY: All probes passing
- RUNNING_NOT_READY: Liveness passing, readiness failing
- RUNNING_UNHEALTHY: Liveness failing
- TERMINATING: Container being stopped

Transitions:
PENDING ──start container──► STARTING ──startup success──► RUNNING_HEALTHY
                                                      │
                                 readiness fail ◄─────┼────► readiness pass
                                     │                │
                                     ▼                │
                              RUNNING_NOT_READY ──────┘
                                     │
                                     │ liveness fail
                                     ▼
                              RUNNING_UNHEALTHY
                                     │
                                     ▼
                                 restart
```

---

## 3. Visual Representations

### 3.1 Probe Interaction Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PROBE INTERACTION LIFECYCLE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  KUBERNETES                    CONTAINER                         EXTERNAL   │
│  CONTROL PLANE                 RUNTIME                           SERVICES   │
│                                                                             │
│       │                           │                                  │      │
│       │ Create Pod                │                                  │      │
│       ├──────────────────────────►│                                  │      │
│       │                           │                                  │      │
│       │                           │ Pull Image                       │      │
│       │                           │◄────────────────────────────────►│      │
│       │                           │                                  │      │
│       │                           │ Start Container                  │      │
│       │                           │◄────────────────────────────────►│      │
│       │                           │                                  │      │
│       │                           │     ┌─────────────────────┐      │      │
│       │ Start Startup Probe       │     │  Application Init   │      │      │
│       ├──────────────────────────►│────►│  • Load config      │      │      │
│       │                           │     │  • Connect DB       │      │      │
│       │                           │     │  • Warm cache       │      │      │
│       │                           │     └─────────────────────┘      │      │
│       │                           │                                  │      │
│       │ GET /health/startup       │                                  │      │
│       ├──────────────────────────►│                                  │      │
│       │                           │                                  │      │
│       │ HTTP 200 ◄────────────────┤                                  │      │
│       │                           │                                  │      │
│       │ Stop Startup              │                                  │      │
│       │ Start Liveness/Readiness  │                                  │      │
│       ├──────────────────────────►│                                  │      │
│       │                           │                                  │      │
│       ═══════════════════════════════════════════════════════════    │      │
│       │                           │    STEADY STATE                  │      │
│       ═══════════════════════════════════════════════════════════    │      │
│       │                           │                                  │      │
│       │ GET /health/live          │                                  │      │
│       ├──────────────────────────►│                                  │      │
│       │ HTTP 200 ◄────────────────┤                                  │      │
│       │                           │                                  │      │
│       │ GET /health/ready         │                                  │      │
│       ├──────────────────────────►│                                  │      │
│       │ HTTP 200 ◄────────────────┤                                  │      │
│       │                           │                                  │      │
│       │ Update Pod Status         │                                  │      │
│       │ Ready: True               │                                  │      │
│       │◄──────────────────────────┤                                  │      │
│       │                           │                                  │      │
│       ▼                           ▼                                  ▼      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Failure Detection and Recovery

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FAILURE DETECTION AND RECOVERY                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SCENARIO 1: Application Deadlock                                           │
│  ═══════════════════════════════                                            │
│                                                                             │
│  Time ─────────────────────────────────────────────────────────────────►   │
│                                                                             │
│  App State:   [Running]    [Stuck/Deadlock]                                 │
│                 │               │                                           │
│  Liveness:    [OK]  [OK]  [OK]  [FAIL] [FAIL] [FAIL]                       │
│                 10s   10s   10s   10s   10s   10s                           │
│                                              │                              │
│                                              ▼                              │
│                                         [RESTART]                           │
│                                              │                              │
│  App State:                                  ▼                              │
│                                         [Running]                           │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  SCENARIO 2: Database Unavailable                                           │
│  ════════════════════════════════                                           │
│                                                                             │
│  Time ─────────────────────────────────────────────────────────────────►   │
│                                                                             │
│  App State:   [Healthy]    [DB Down]                                        │
│                 │               │                                           │
│  Liveness:    [OK]  [OK]  [OK]  [OK]  [OK]  [OK]  ← Shallow check passes   │
│                 10s   10s   10s   10s   10s   10s                           │
│                                                                             │
│  Readiness:   [OK]  [OK]  [OK]  [FAIL] [FAIL] [FAIL]                        │
│                 5s    5s    5s    5s    5s    5s                            │
│                                              │                              │
│                                              ▼                              │
│                                    [REMOVE FROM LB]                         │
│                                              │                              │
│                                              │ DB recovers                 │
│                                              ▼                              │
│                                         [ADD TO LB]                         │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  SCENARIO 3: Slow Startup                                                   │
│  ════════════════════════════                                               │
│                                                                             │
│  Time ─────────────────────────────────────────────────────────────────►   │
│                                                                             │
│  App State:   [Initializing]         [Ready]                                │
│                    │                      │                                 │
│  Startup:     [FAIL] [FAIL] [FAIL] ... [OK]                                 │
│               10s    10s    10s        10s                                  │
│                    │                        │                               │
│                    │ (retrying)             │ (success)                     │
│                    │                        │                               │
│  Without startup probe:                                                     │
│  Liveness would restart container repeatedly during initialization!         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Probe Configuration Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PROBE CONFIGURATION MATRIX                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  APPLICATION TYPE              │ STARTUP │ LIVENESS │ READINESS │ NOTES   │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  Fast-starting (< 10s)         │  Skip   │  10s/3   │  5s/3     │ Simple  │
│  (Go static binary)            │         │          │           │ apps    │
│                                                                             │
│  Medium-starting (10s-2min)    │  10s/30 │  10s/3   │  5s/3     │ Most    │
│  (Java, Node.js)               │         │          │           │ apps    │
│                                                                             │
│  Slow-starting (> 2min)        │  10s/60 │  10s/3   │  5s/3     │ Legacy  │
│  (Large JVM heap)              │         │          │           │ systems │
│                                                                             │
│  Memory-intensive              │  10s/60 │  15s/5   │  5s/3     │ Prevent │
│  (Data processing)             │         │          │           │ thrash  │
│                                                                             │
│  External dependency-heavy     │  10s/60 │  10s/3   │  5s/3     │ Deep    │
│  (Microservices)               │         │          │           │ checks  │
│                                                                             │
│  LEGEND: period/failureThreshold                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Probe Server Implementation

```go
package probes

import (
 "context"
 "fmt"
 "net/http"
 "sync"
 "time"

 "go.uber.org/zap"
)

// ProbeType represents the type of probe
type ProbeType string

const (
 ProbeStartup  ProbeType = "startup"
 ProbeLiveness ProbeType = "liveness"
 ProbeReadiness ProbeType = "readiness"
)

// State represents the application state
type State int

const (
 StateInitializing State = iota
 StateRunning
 StateShuttingDown
)

// Server manages health probes
type Server struct {
 port      int
 logger    *zap.Logger

 state     State
 stateMu   sync.RWMutex

 // Probe functions
 startupCheck   func() error
 livenessCheck  func() error
 readinessCheck func() error

 server *http.Server
}

// Config holds probe server configuration
type Config struct {
 Port           int
 Logger         *zap.Logger
 StartupCheck   func() error
 LivenessCheck  func() error
 ReadinessCheck func() error
}

// New creates a new probe server
func New(cfg Config) *Server {
 s := &Server{
  port:           cfg.Port,
  logger:         cfg.Logger,
  state:          StateInitializing,
  startupCheck:   cfg.StartupCheck,
  livenessCheck:  cfg.LivenessCheck,
  readinessCheck: cfg.ReadinessCheck,
 }

 mux := http.NewServeMux()
 mux.HandleFunc("/health/startup", s.handleStartup)
 mux.HandleFunc("/health/live", s.handleLiveness)
 mux.HandleFunc("/health/ready", s.handleReadiness)

 s.server = &http.Server{
  Addr:    fmt.Sprintf(":%d", cfg.Port),
  Handler: mux,
 }

 return s
}

// Start starts the probe server
func (s *Server) Start() error {
 s.logger.Info("Starting probe server", zap.Int("port", s.port))
 go func() {
  if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
   s.logger.Error("Probe server failed", zap.Error(err))
  }
 }()
 return nil
}

// Stop gracefully stops the probe server
func (s *Server) Stop(ctx context.Context) error {
 s.SetState(StateShuttingDown)
 return s.server.Shutdown(ctx)
}

// SetState updates the application state
func (s *Server) SetState(state State) {
 s.stateMu.Lock()
 defer s.stateMu.Unlock()
 s.state = state
 s.logger.Info("Application state changed", zap.String("state", stateString(state)))
}

// GetState returns the current application state
func (s *Server) GetState() State {
 s.stateMu.RLock()
 defer s.stateMu.RUnlock()
 return s.state
}

func stateString(s State) string {
 switch s {
 case StateInitializing:
  return "initializing"
 case StateRunning:
  return "running"
 case StateShuttingDown:
  return "shutting_down"
 default:
  return "unknown"
 }
}

func (s *Server) handleStartup(w http.ResponseWriter, r *http.Request) {
 // Startup probe: verify initialization complete
 state := s.GetState()

 if state == StateShuttingDown {
  http.Error(w, "Application is shutting down", http.StatusServiceUnavailable)
  return
 }

 if state == StateInitializing {
  // Still initializing - check if we can transition to running
  if s.startupCheck != nil {
   if err := s.startupCheck(); err != nil {
    s.logger.Debug("Startup check failed", zap.Error(err))
    http.Error(w, fmt.Sprintf("Initializing: %v", err), http.StatusServiceUnavailable)
    return
   }
  }
  // Transition to running state
  s.SetState(StateRunning)
 }

 w.WriteHeader(http.StatusOK)
 w.Write([]byte("OK"))
}

func (s *Server) handleLiveness(w http.ResponseWriter, r *http.Request) {
 // Liveness probe: simple check that process is not deadlocked
 state := s.GetState()

 if state == StateShuttingDown {
  // During graceful shutdown, we're still "alive" until complete
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("OK - shutting down"))
  return
 }

 if state == StateInitializing {
  // Not yet ready, but alive
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("OK - initializing"))
  return
 }

 // Perform liveness check if configured
 if s.livenessCheck != nil {
  if err := s.livenessCheck(); err != nil {
   s.logger.Error("Liveness check failed", zap.Error(err))
   http.Error(w, err.Error(), http.StatusInternalServerError)
   return
  }
 }

 w.WriteHeader(http.StatusOK)
 w.Write([]byte("OK"))
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
 // Readiness probe: check if ready to receive traffic
 state := s.GetState()

 if state == StateShuttingDown {
  http.Error(w, "Application is shutting down", http.StatusServiceUnavailable)
  return
 }

 if state == StateInitializing {
  http.Error(w, "Application is initializing", http.StatusServiceUnavailable)
  return
 }

 // Perform readiness check if configured
 if s.readinessCheck != nil {
  if err := s.readinessCheck(); err != nil {
   s.logger.Warn("Readiness check failed", zap.Error(err))
   http.Error(w, err.Error(), http.StatusServiceUnavailable)
   return
  }
 }

 w.WriteHeader(http.StatusOK)
 w.Write([]byte("OK"))
}
```

### 4.2 Kubernetes Deployment Configuration

```go
package probes

import (
 corev1 "k8s.io/api/core/v1"
 "k8s.io/apimachinery/pkg/util/intstr"
)

// ProbeConfig holds Kubernetes probe configuration
type ProbeConfig struct {
 // Startup probe
 StartupEnabled       bool
 StartupInitialDelay  int32
 StartupPeriod        int32
 StartupTimeout       int32
 StartupFailureThreshold int32

 // Liveness probe
 LivenessInitialDelay int32
 LivenessPeriod       int32
 LivenessTimeout      int32
 LivenessFailureThreshold int32

 // Readiness probe
 ReadinessInitialDelay int32
 ReadinessPeriod       int32
 ReadinessTimeout      int32
 ReadinessFailureThreshold int32
 ReadinessSuccessThreshold int32
}

// DefaultProbeConfig returns default configuration for a typical Go application
func DefaultProbeConfig() ProbeConfig {
 return ProbeConfig{
  StartupEnabled:        true,
  StartupInitialDelay:   0,
  StartupPeriod:         10,
  StartupTimeout:        5,
  StartupFailureThreshold: 30, // 5 minutes for slow startup

  LivenessInitialDelay:  0,
  LivenessPeriod:        10,
  LivenessTimeout:       5,
  LivenessFailureThreshold: 3,

  ReadinessInitialDelay:  0,
  ReadinessPeriod:        5,
  ReadinessTimeout:       3,
  ReadinessFailureThreshold: 3,
  ReadinessSuccessThreshold: 1,
 }
}

// ConfigureProbes adds probe configuration to a container
func ConfigureProbes(container *corev1.Container, cfg ProbeConfig) {
 port := intstr.FromInt(8080)

 // Startup probe
 if cfg.StartupEnabled {
  container.StartupProbe = &corev1.Probe{
   ProbeHandler: corev1.ProbeHandler{
    HTTPGet: &corev1.HTTPGetAction{
     Path: "/health/startup",
     Port: port,
    },
   },
   InitialDelaySeconds: cfg.StartupInitialDelay,
   PeriodSeconds:       cfg.StartupPeriod,
   TimeoutSeconds:      cfg.StartupTimeout,
   FailureThreshold:    cfg.StartupFailureThreshold,
  }
 }

 // Liveness probe
 container.LivenessProbe = &corev1.Probe{
  ProbeHandler: corev1.ProbeHandler{
   HTTPGet: &corev1.HTTPGetAction{
    Path: "/health/live",
    Port: port,
   },
  },
  InitialDelaySeconds: cfg.LivenessInitialDelay,
  PeriodSeconds:       cfg.LivenessPeriod,
  TimeoutSeconds:      cfg.LivenessTimeout,
  FailureThreshold:    cfg.LivenessFailureThreshold,
 }

 // Readiness probe
 container.ReadinessProbe = &corev1.Probe{
  ProbeHandler: corev1.ProbeHandler{
   HTTPGet: &corev1.HTTPGetAction{
    Path: "/health/ready",
    Port: port,
   },
  },
  InitialDelaySeconds: cfg.ReadinessInitialDelay,
  PeriodSeconds:       cfg.ReadinessPeriod,
  TimeoutSeconds:      cfg.ReadinessTimeout,
  FailureThreshold:    cfg.ReadinessFailureThreshold,
  SuccessThreshold:    cfg.ReadinessSuccessThreshold,
 }
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Probe Failure Taxonomy

| Scenario | Symptoms | Root Cause | Mitigation |
|----------|----------|------------|------------|
| **Probe Timeout** | Container restarting | Check too slow | Increase timeout, simplify check |
| **Premature Restart** | Startup loop | Low threshold | Increase failureThreshold |
| **Traffic to Unready** | 500 errors | Missing readiness check | Add readiness probe |
| **Zombie Process** | Not detected | Liveness too simple | Add application-specific check |
| **Cascade Failure** | All pods unhealthy | Deep dependency check | Use shallow checks |

---

## 6. Semantic Trade-off Analysis

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PROBE BALANCE DECISION MATRIX                             │
├────────────────────────┬─────────────────┬─────────────────┬─────────────────┤
│        Goal            │   Aggressive    │   Balanced      │   Conservative  │
├────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│ Startup Detection      │ Fast (low T)    │ Moderate        │ Slow (high T)   │
│ Deadlock Detection     │ Fast restart    │ Balanced        │ Slow restart    │
│ False Positive Rate    │ High            │ Medium          │ Low             │
│ Availability Impact    │ May hurt        │ Balanced        │ Protects        │
│ Resource Usage         │ High            │ Moderate        │ Lower           │
└────────────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

---

## 7. References

1. Kubernetes Documentation. (2024). *Container Probes*. kubernetes.io.
2. Google Cloud. (2024). *Kubernetes Best Practices: Terminating with Grace*.
3. Richardson, C. (2020). *Health Check API*. microservices.io.
