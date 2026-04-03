# EC-060: Chaos Engineering Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #chaos-engineering #resilience-testing #failure-injection #gremlin #litmus
> **Authoritative Sources**:
>
> - [Chaos Engineering Book](https://www.oreilly.com/library/view/chaos-engineering/9781491983850/) - Basiri et al. (2017)
> - [Principles of Chaos](https://principlesofchaos.org/) - Chaos Engineering Community (2024)
> - [AWS Fault Injection Simulator](https://aws.amazon.com/fis/) - AWS (2024)
> - [Chaos Mesh](https://chaos-mesh.org/) - CNCF (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Chaos Engineering Domain)**
Let $\mathcal{S}$ be a distributed system with steady-state behavior $SS(\mathcal{S})$. Chaos engineering introduces faults $F = \{f_1, f_2, ..., f_n\}$ to verify:

$$\forall f \in F: fault(f) \land SS(\mathcal{S}) = SS(\mathcal{S}') \lor graceful\_degradation(\mathcal{S})$$

**Steady-State Hypothesis:**

- System has measurable normal behavior
- Faults should not change steady-state (or degrade gracefully)

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Blast Radius** | $|\{affected\_components\}| < \theta_{max}$ | Limit experiment scope |
| **Abort Condition** | $\exists m \in metrics: m > threshold \Rightarrow abort$ | Automatic safety |
| **Time Bound** | $T_{experiment} < T_{max}$ | Limited window |
| **Reversibility** | $\forall f: \exists f^{-1}: apply(f^{-1}) \Rightarrow recover$ | Can undo faults |

### 1.2 Problem Statement

**Problem 1.1 (Resilience Verification)**
Given system $\mathcal{S}$ and fault $f$, verify resilience property $P$:

$$verify(\mathcal{S}, f, P) = apply(f) \land observe(\mathcal{S}) \models P \land recover(f)$$

**Key Challenges:**

1. **Hypothesis Formulation**: Defining meaningful steady-state
2. **Fault Selection**: Choosing realistic failure modes
3. **Safety**: Preventing actual outages
4. **Measurement**: Detecting subtle degradation
5. **Cultural Adoption**: Building organizational confidence

---

## 2. Solution Architecture

### 2.1 Chaos Experiment Types

| Level | Fault Type | Example | Tool |
|-------|------------|---------|------|
| **Infrastructure** | Node failure | Terminate VM/instance | AWS FIS, Gremlin |
| **Network** | Latency, partition | Delay packets, drop connections | Chaos Mesh, Toxiproxy |
| **Application** | Resource exhaustion | CPU/memory stress | Stress-ng, Chaos Monkey |
| **Dependency** | Service failure | Return 503 from downstream | MockServer |
| **Data** | Database corruption | Inject latency in queries | pgbench |

### 2.2 Experiment Lifecycle

```
Steady State Definition → Fault Hypothesis → Experiment Design → Run → Analyze → Improve
```

---

## 3. Visual Representations

### 3.1 Chaos Engineering Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CHAOS ENGINEERING ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CONTROL PLANE                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Chaos Engineering Platform                       │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │   │
│  │  │   Experiment│  │   Steady    │  │   Safety    │  │  Analysis  │  │   │
│  │  │   Designer  │  │   State     │  │   Guardrails│  │  Engine    │  │   │
│  │  │             │  │   Monitor   │  │             │  │            │  │   │
│  │  │ • Define    │  │ • Baseline  │  │ • Abort     │  │ • Compare  │  │   │
│  │  │   faults    │  │   metrics   │  │   thresholds│  │   results  │  │   │
│  │  │ • Set scope │  │ • Track     │  │ • Auto-stop │  │ • Generate │  │   │
│  │  │ • Schedule  │  │   deviations│  │   on alert  │  │   reports  │  │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └─────┬──────┘  │   │
│  │         │                │                │               │         │   │
│  │         └────────────────┴────────────────┴───────────────┘         │   │
│  │                          │                                         │   │
│  └──────────────────────────┼─────────────────────────────────────────┘   │
│                             │                                             │
│  EXECUTION PLANE            │ Orchestration                               │
│                             ▼                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Chaos Agents / Daemons                          │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │   Node      │  │   Network   │  │   Container │                  │   │
│  │  │   Agent     │  │   Agent     │  │   Agent     │                  │   │
│  │  │             │  │             │  │             │                  │   │
│  │  │ • CPU burn  │  │ • Latency   │  │ • Kill pod  │                  │   │
│  │  │ • Memory    │  │ • Packet    │  │ • Stress    │                  │   │
│  │  │   pressure  │  │   loss      │  │   resource  │                  │   │
│  │  │ • Disk fill │  │ • Partition │  │ • IO delay  │                  │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                  │   │
│  │         │                │                │                         │   │
│  └─────────┼────────────────┼────────────────┼─────────────────────────┘   │
│            │                │                │                              │
│  TARGET SYSTEM              │                │                              │
│            │                │                │                              │
│  ┌─────────┴────────────────┴────────────────┴─────────────────────────┐   │
│  │                                                                     │   │
│  │  Service Mesh / Kubernetes Cluster                                  │   │
│  │                                                                     │   │
│  │  ┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐    │   │
│  │  │ Pod 1   │◄────►│ Pod 2   │◄────►│ Pod 3   │◄────►│ Pod 4   │    │   │
│  │  │         │      │         │      │         │      │         │    │   │
│  │  │ [Agent] │      │ [Agent] │      │ [Agent] │      │         │    │   │
│  │  └─────────┘      └─────────┘      └─────────┘      └─────────┘    │   │
│  │       ▲                                  ▲                          │   │
│  │       │                                  │                          │   │
│  │       └──────┬───────────────────────────┘                          │   │
│  │              │                                                       │   │
│  │         [Network Latency Injected]                                   │   │
│  │         [Pod 3 Terminated]                                           │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  OBSERVABILITY                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Metrics ──► Alert? ──► ABORT ──► Rollback Fault                    │   │
│  │                                                                     │   │
│  │  Dashboard:                                                         │   │
│  │  • Error rate: 0.5% (baseline: 0.1%) ⚠️                             │   │
│  │  • Latency P99: 450ms (baseline: 200ms) ⚠️                          │   │
│  │  • Availability: 99.9% ✅                                           │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Game Day Exercise Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GAME DAY EXERCISE FLOW                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  PREPARATION (Day -7)                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  1. Define Steady State                                        │  │   │
│  │  │     • Error rate < 0.1%                                        │  │   │
│  │  │     • P99 latency < 500ms                                      │  │   │
│  │  │     • Throughput > 1000 req/s                                  │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                     │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  2. Form Hypothesis                                            │  │   │
│  │  │     "If database primary fails, read replicas will handle      │  │   │
│  │  │      traffic with < 5% error rate"                             │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                     │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  3. Set Safety Limits                                          │  │   │
│  │  │     • Abort if error rate > 10%                                │  │   │
│  │  │     • Max blast radius: 25% of traffic                         │  │   │
│  │  │     • Auto-rollback after 10 minutes                           │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  EXECUTION (Game Day)                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │  T+0:00  ──► Inject: Terminate database primary                     │   │
│  │            System: Automatic failover to replica                    │   │
│  │                                                                     │   │
│  │  T+0:30  ──► Monitor: Error rate spike 0.1% → 2%                    │   │
│  │            Status: ⚠️ Within acceptable limits                      │   │
│  │                                                                     │   │
│  │  T+1:00  ──► Monitor: Error rate stabilizes at 0.5%                 │   │
│  │            Status: ✅ Hypothesis confirmed                          │   │
│  │                                                                     │   │
│  │  T+5:00  ──► Recovery: Restore primary database                     │   │
│  │            System: Replication catch-up in progress                 │   │
│  │                                                                     │   │
│  │  T+10:00 ──► Complete: Full capacity restored                       │   │
│  │            Status: ✅ Experiment successful                         │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  POST-MORTEM (Day +1)                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │  Findings:                                                          │   │
│  │  • Failover worked as expected (90s to complete)                    │   │
│  │  • Error rate during failover was 2% (higher than expected)         │   │
│  │  • Connection pool exhaustion caused extended recovery              │   │
│  │                                                                     │   │
│  │  Action Items:                                                      │   │
│  │  1. Increase connection pool size [OWNER: Platform]                 │   │
│  │  2. Add retry with exponential backoff [OWNER: Backend]             │   │
│  │  3. Improve failover speed [OWNER: DBA]                             │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Chaos Experiment Framework

```go
package chaos

import (
 "context"
 "fmt"
 "sync"
 "time"

 "go.uber.org/zap"
)

// Experiment represents a chaos experiment
type Experiment struct {
 ID          string
 Name        string
 Description string
 Hypothesis  string
 Faults      []Fault
 AbortConditions []AbortCondition
 Duration    time.Duration
 Cooldown    time.Duration
}

// Fault represents a fault to inject
type Fault interface {
 Name() string
 Inject(ctx context.Context) error
 Recover(ctx context.Context) error
}

// AbortCondition defines when to stop the experiment
type AbortCondition struct {
 Name      string
 Query     string
 Threshold float64
 Operator  string // >, <, ==
}

// SteadyState represents baseline metrics
type SteadyState struct {
 Metrics map[string]float64
}

// Runner executes chaos experiments
type Runner struct {
 metrics  MetricsClient
 logger   *zap.Logger
 mu       sync.Mutex
 active   map[string]context.CancelFunc
}

// MetricsClient queries system metrics
type MetricsClient interface {
 Query(ctx context.Context, query string) (float64, error)
}

// NewRunner creates a new chaos runner
func NewRunner(metrics MetricsClient, logger *zap.Logger) *Runner {
 return &Runner{
  metrics: metrics,
  logger:  logger,
  active:  make(map[string]context.CancelFunc),
 }
}

// Run executes a chaos experiment
func (r *Runner) Run(ctx context.Context, exp *Experiment) (*Result, error) {
 r.mu.Lock()
 if _, exists := r.active[exp.ID]; exists {
  r.mu.Unlock()
  return nil, fmt.Errorf("experiment %s already running", exp.ID)
 }

 ctx, cancel := context.WithCancel(ctx)
 r.active[exp.ID] = cancel
 r.mu.Unlock()

 defer func() {
  r.mu.Lock()
  delete(r.active, exp.ID)
  r.mu.Unlock()
 }()

 r.logger.Info("Starting chaos experiment",
  zap.String("id", exp.ID),
  zap.String("name", exp.Name))

 // Capture steady state
 steadyState, err := r.captureSteadyState(ctx, exp)
 if err != nil {
  return nil, fmt.Errorf("failed to capture steady state: %w", err)
 }

 // Run experiment with monitoring
 result := &Result{
  ExperimentID: exp.ID,
  StartTime:    time.Now(),
  SteadyState:  steadyState,
 }

 // Start abort monitor
 abortCh := make(chan string)
 go r.monitorAbortConditions(ctx, exp, abortCh)

 // Inject faults
 for _, fault := range exp.Faults {
  select {
  case <-ctx.Done():
   return result, ctx.Err()
  case reason := <-abortCh:
   result.Aborted = true
   result.AbortReason = reason
   r.recoverAll(exp)
   return result, nil
  default:
  }

  r.logger.Info("Injecting fault",
   zap.String("experiment", exp.ID),
   zap.String("fault", fault.Name()))

  if err := fault.Inject(ctx); err != nil {
   r.logger.Error("Failed to inject fault",
    zap.String("fault", fault.Name()),
    zap.Error(err))
   result.FaultErrors = append(result.FaultErrors, err)
  }

  result.InjectedFaults = append(result.InjectedFaults, fault.Name())
 }

 // Wait for duration or abort
 select {
 case <-time.After(exp.Duration):
  // Experiment completed
 case <-ctx.Done():
  result.Aborted = true
 case reason := <-abortCh:
  result.Aborted = true
  result.AbortReason = reason
 }

 // Recover all faults
 r.recoverAll(exp)

 result.EndTime = time.Now()
 result.Duration = result.EndTime.Sub(result.StartTime)

 r.logger.Info("Chaos experiment completed",
  zap.String("id", exp.ID),
  zap.Duration("duration", result.Duration),
  zap.Bool("aborted", result.Aborted))

 return result, nil
}

func (r *Runner) captureSteadyState(ctx context.Context, exp *Experiment) (*SteadyState, error) {
 // Capture baseline metrics before injection
 state := &SteadyState{
  Metrics: make(map[string]float64),
 }

 for _, condition := range exp.AbortConditions {
  value, err := r.metrics.Query(ctx, condition.Query)
  if err != nil {
   return nil, err
  }
  state.Metrics[condition.Name] = value
 }

 return state, nil
}

func (r *Runner) monitorAbortConditions(ctx context.Context, exp *Experiment, abortCh chan<- string) {
 ticker := time.NewTicker(5 * time.Second)
 defer ticker.Stop()

 for {
  select {
  case <-ctx.Done():
   return
  case <-ticker.C:
   for _, condition := range exp.AbortConditions {
    value, err := r.metrics.Query(ctx, condition.Query)
    if err != nil {
     r.logger.Error("Failed to query metric", zap.Error(err))
     continue
    }

    shouldAbort := false
    switch condition.Operator {
    case ">":
     shouldAbort = value > condition.Threshold
    case "<":
     shouldAbort = value < condition.Threshold
    }

    if shouldAbort {
     r.logger.Warn("Abort condition triggered",
      zap.String("condition", condition.Name),
      zap.Float64("value", value),
      zap.Float64("threshold", condition.Threshold))

     select {
     case abortCh <- condition.Name:
     case <-ctx.Done():
     }
     return
    }
   }
  }
 }
}

func (r *Runner) recoverAll(exp *Experiment) {
 r.logger.Info("Recovering all faults", zap.String("experiment", exp.ID))

 for _, fault := range exp.Faults {
  ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  if err := fault.Recover(ctx); err != nil {
   r.logger.Error("Failed to recover fault",
    zap.String("fault", fault.Name()),
    zap.Error(err))
  }
  cancel()
 }
}

// Result represents experiment results
type Result struct {
 ExperimentID   string
 StartTime      time.Time
 EndTime        time.Time
 Duration       time.Duration
 SteadyState    *SteadyState
 InjectedFaults []string
 FaultErrors    []error
 Aborted        bool
 AbortReason    string
}

// Stop stops a running experiment
func (r *Runner) Stop(experimentID string) {
 r.mu.Lock()
 cancel, exists := r.active[experimentID]
 r.mu.Unlock()

 if exists {
  cancel()
 }
}
```

### 4.2 Example Fault Implementations

```go
package chaos

import (
 "context"
 "fmt"
 "os"
 "syscall"
 "time"
)

// HTTPDelayFault adds latency to HTTP responses
type HTTPDelayFault struct {
 Duration time.Duration
}

func (f *HTTPDelayFault) Name() string {
 return fmt.Sprintf("http_delay_%s", f.Duration)
}

func (f *HTTPDelayFault) Inject(ctx context.Context) error {
 // Implementation would register with HTTP middleware
 return nil
}

func (f *HTTPDelayFault) Recover(ctx context.Context) error {
 // Remove delay from HTTP middleware
 return nil
}

// PodKillFault simulates Kubernetes pod failure
type PodKillFault struct {
 Namespace string
 Label     string
 Count     int
}

func (f *PodKillFault) Name() string {
 return fmt.Sprintf("pod_kill_%s_%s", f.Namespace, f.Label)
}

func (f *PodKillFault) Inject(ctx context.Context) error {
 // Use Kubernetes API to delete pods
 return nil
}

func (f *PodKillFault) Recover(ctx context.Context) error {
 // Pods auto-recover via replica set
 return nil
}

// CPUStressFault consumes CPU resources
type CPUStressFault struct {
 Cores    int
 Duration time.Duration
}

func (f *CPUStressFault) Name() string {
 return fmt.Sprintf("cpu_stress_%d_cores", f.Cores)
}

func (f *CPUStressFault) Inject(ctx context.Context) error {
 // Start CPU-intensive goroutines
 for i := 0; i < f.Cores; i++ {
  go func() {
   for {
    select {
    case <-ctx.Done():
     return
    default:
     // Busy loop
    }
   }
  }()
 }
 return nil
}

func (f *CPUStressFault) Recover(ctx context.Context) error {
 // Goroutines stop when context is cancelled
 return nil
}

// DiskFillFault fills disk space
type DiskFillFault struct {
 Path     string
 SizeGB   int
}

func (f *DiskFillFault) Name() string {
 return fmt.Sprintf("disk_fill_%s_%dgb", f.Path, f.SizeGB)
}

func (f *DiskFillFault) Inject(ctx context.Context) error {
 // Create large file to fill disk
 file, err := os.CreateTemp(f.Path, "chaos_fill_*.dat")
 if err != nil {
  return err
 }
 defer file.Close()

 size := int64(f.SizeGB) * 1024 * 1024 * 1024
 if err := file.Truncate(size); err != nil {
  return err
 }

 return nil
}

func (f *DiskFillFault) Recover(ctx context.Context) error {
 // Remove the fill file
 // Implementation would track created file
 return nil
}
```

---

## 5. Failure Scenarios and Mitigations

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Abort Missed** | Outage | Multiple alert channels | Automatic timeout |
| **Recovery Failure** | Permanent degradation | Post-experiment health check | Manual intervention |
| **Scope Creep** | Too many services affected | Blast radius validation | Strict scoping rules |
| **Data Corruption** | Permanent data loss | Pre-experiment backups | Read-only tests first |

---

## 6. Semantic Trade-off Analysis

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CHAOS ENGINEERING APPROACHES                              │
├─────────────────────┬─────────────────┬─────────────────────────────────────┤
│     Dimension       │  Game Day       │         Automated Continuous        │
├─────────────────────┼─────────────────┼─────────────────────────────────────┤
│ Frequency           │ Monthly/Quarter │ Daily/Hourly                        │
│ Scope               │ Large, complex  │ Small, targeted                     │
│ Team Involvement    │ High            │ Low                                 │
│ Risk Level          │ Higher          │ Lower                               │
│ Learning Depth      │ Deep            │ Incremental                         │
│ Cultural Impact     │ High            │ Gradual                             │
└─────────────────────┴─────────────────┴─────────────────────────────────────┘
```

---

## 7. References

1. Basiri, A., et al. (2017). *Chaos Engineering*. O'Reilly Media.
2. Principles of Chaos. (2024). *Chaos Engineering Principles*. principlesofchaos.org.
3. AWS. (2024). *Fault Injection Simulator*. aws.amazon.com/fis.
4. CNCF. (2024). *Chaos Mesh*. chaos-mesh.org.
