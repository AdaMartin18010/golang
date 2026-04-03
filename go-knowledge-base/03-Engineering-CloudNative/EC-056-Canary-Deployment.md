# EC-056: Canary Deployment Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #canary-deployment # progressive-delivery #deployment-strategy #risk-mitigation #kubernetes
> **Authoritative Sources**:
>
> - [Continuous Delivery](https://continuousdelivery.com/) - Humble & Farley (2010)
> - [Progressive Delivery](https://redmonk.com/jgovernor/2018/08/06/progressive-delivery/) - Governor (2018)
> - [Flagger Documentation](https://docs.flagger.app/) - FluxCD (2024)
> - [Spinnaker](https://spinnaker.io/docs/guides/user/canary/) - Netflix (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Canary Deployment Domain)**
Let $V_{current}$ be the current version serving 100% traffic. A canary deployment introduces $V_{new}$ with traffic split $\alpha(t)$ where:

- $\alpha(0) = 0$ (no traffic to new version initially)
- $\alpha(T) = 1$ (full rollout if successful)
- $\alpha(t)$ is monotonically increasing during promotion

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Blast Radius** | $|\{affected\_users\}| < \theta_{max}$ | Must limit impact of failures |
| **Rollback Time** | $T_{rollback} < T_{SLA}$ | Must recover quickly |
| **Metric Confidence** | $confidence(metrics) > 0.95$ | Decisions need statistical validity |
| **Version Compatibility** | $compatible(V_{current}, V_{new})$ | Database/schema migrations |

### 1.2 Problem Statement

**Problem 1.1 (Safe Deployment)**
Find traffic split function $\alpha(t)$ such that:

$$\forall t: risk(V_{new}, \alpha(t)) < risk_{acceptable} \land detect(failure) \Rightarrow \alpha(t+\Delta) = 0$$

**Key Challenges:**

1. **Traffic Splitting**: Routing subset of users to new version
2. **Health Monitoring**: Detecting issues at low traffic volumes
3. **Automatic Promotion**: Increasing traffic when healthy
4. **Instant Rollback**: Reverting on failure detection
5. **Metric Validation**: Determining success criteria

---

## 2. Solution Architecture

### 2.1 Canary Stage Definition

| Stage | Traffic % | Duration | Criteria | Action on Failure |
|-------|-----------|----------|----------|-------------------|
| **Baseline** | 0% | 5 min | Metrics collection | N/A |
| **Canary 1** | 5% | 10 min | Error rate < 1% | Rollback |
| **Canary 2** | 25% | 15 min | P99 latency < baseline + 20% | Rollback |
| **Canary 3** | 50% | 15 min | All golden signals healthy | Rollback |
| **Full Rollout** | 100% | - | Sustained health | N/A |

---

## 3. Visual Representations

### 3.1 Canary Deployment Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CANARY DEPLOYMENT ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         TRAFFIC SOURCE                               │   │
│  │                           (Users/API)                                │   │
│  └──────────────────────────────────┬──────────────────────────────────┘   │
│                                     │                                       │
│                                     ▼                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        LOAD BALANCER / INGRESS                       │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                    Traffic Split Controller                  │    │   │
│  │  │                                                               │    │   │
│  │  │   Current: 75% ───────────────────────────────────────┐      │    │   │
│  │  │   Canary:  25% ──────────────┐                        │      │    │   │
│  │  │                               │                        │      │    │   │
│  │  │                               ▼                        ▼      │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                                                                      │   │
│  └──────────────────────────────────┬──────────────────────────────────┘   │
│                                     │                                       │
│         ┌───────────────────────────┴───────────────────────────┐          │
│         │                                                       │          │
│         ▼                                                       ▼          │
│  ┌───────────────────┐                                  ┌──────────────────┐│
│  │  STABLE VERSION   │                                  │  CANARY VERSION  ││
│  │  (v1.2.3)         │                                  │  (v1.3.0)        ││
│  │                   │                                  │                  ││
│  │  Replicas: 6      │                                  │  Replicas: 2     ││
│  │  Traffic: 75%     │                                  │  Traffic: 25%    ││
│  │                   │                                  │                  ││
│  │  ┌─────────────┐  │                                  │  ┌─────────────┐ ││
│  │  │ Pod 1      │  │                                  │  │ Pod 1      │ ││
│  │  │ Pod 2      │  │                                  │  │ Pod 2      │ ││
│  │  │ ...        │  │                                  │  └─────────────┘ ││
│  │  │ Pod 6      │  │                                  │                  ││
│  │  └─────────────┘  │                                  │                  ││
│  └─────────┬─────────┘                                  └──────────┬───────┘│
│            │                                                       │        │
│            └───────────────────┬───────────────────────────────────┘        │
│                                │                                            │
│                                ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      OBSERVABILITY STACK                             │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │  Metrics    │  │    Logs     │  │   Traces    │                  │   │
│  │  │  (Prometheus)│  │   (Fluentd) │  │   (Jaeger)  │                  │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                  │   │
│  │         │                │                │                         │   │
│  │         └────────────────┼────────────────┘                         │   │
│  │                          ▼                                         │   │
│  │                 ┌─────────────────┐                                 │   │
│  │                 │  Canary Analysis │                                │   │
│  │                 │                  │                                │   │
│  │                 │ • Error rate     │                                │   │
│  │                 │ • Latency P99    │                                │   │
│  │                 │ • Throughput     │                                │   │
│  │                 │ • Custom metrics │                                │   │
│  │                 └────────┬────────┘                                 │   │
│  │                          │                                         │   │
│  │         ┌────────────────┼────────────────┐                        │   │
│  │         │                │                │                        │   │
│  │         ▼                ▼                ▼                        │   │
│  │  ┌───────────┐    ┌───────────┐    ┌───────────┐                   │   │
│  │  │ Promote   │    │  Hold     │    │ Rollback  │                   │   │
│  │  │ (+25%)    │    │ (current) │    │ (to v1.2.3)                  │   │
│  │  └───────────┘    └───────────┘    └───────────┘                   │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Canary Promotion Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CANARY PROMOTION FLOW                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  TIME ─────────────────────────────────────────────────────────────────►   │
│                                                                             │
│  TRAFFIC:                                                                  │
│                                                                             │
│  100% │                                           ┌──────────────────────┐ │
│       │                                           │  FULL ROLLOUT        │ │
│   75% │                    ┌──────────────────────┤  (v1.3.0 @ 100%)     │ │
│       │                    │  CANARY 3            │                      │ │
│   50% │     ┌──────────────┤  (v1.3.0 @ 50%)      │                      │ │
│       │     │  CANARY 2    │                      │                      │ │
│   25% │     │  (v1.3.0 @ 25%)                    │                      │ │
│       │     │              │                      │                      │ │
│    5% │─────┘              │                      │                      │ │
│       │  CANARY 1          │                      │                      │ │
│    0% ├──────────────────────────────────────────────────────────────────│ │
│       │  BASELINE          │                      │                      │ │
│       └──────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  METRICS:                                                                  │
│                                                                             │
│  Error Rate:                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Stable:  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │ │
│  │  Canary:        ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │ │
│  │  Threshold: ───────────────────────────────────────────────────────  │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  Latency P99:                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Stable:  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │ │
│  │  Canary:        ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │ │
│  │  Baseline: ────────────────────────────────────────────────────────  │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  DECISION POINTS:                                                          │
│                                                                             │
│  Baseline ──► Canary 1 ──► Canary 2 ──► Canary 3 ──► Full Rollout        │
│     │            │            │            │            │                 │
│     │            │            │            │            │                 │
│     ▼            ▼            ▼            ▼            ▼                 │
│   Collect    Error rate   Latency      All metrics    Complete           │
│   metrics    < 1%?        < baseline   healthy        deployment         │
│   (5 min)    ✓ PASS      + 20%?      ✓ PASS                              │
│              (10 min)    ✓ PASS      (15 min)                            │
│                          (15 min)                                        │
│                                                                             │
│  ROLLBACK PATH:                                                            │
│                                                                             │
│  If any check fails at any stage:                                         │
│  Immediately route 100% traffic back to stable version                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Canary Analysis Engine

```go
package canary

import (
 "context"
 "fmt"
 "math"
 "time"

 "go.uber.org/zap"
)

// AnalysisResult represents the result of canary analysis
type AnalysisResult int

const (
 AnalysisResultPending AnalysisResult = iota
 AnalysisResultSuccess
 AnalysisResultFail
 AnalysisResultWarning
)

// MetricThreshold defines a threshold for a metric
type MetricThreshold struct {
 Name      string
 Query     string
 Min       *float64
 Max       *float64
 Baseline  bool // Compare against baseline
}

// AnalysisConfig holds configuration for canary analysis
type AnalysisConfig struct {
 Interval          time.Duration
 Thresholds        []MetricThreshold
 SuccessRate       float64 // Required success rate (0-1)
 RequiredSuccessCount int
 MaxFailures       int
}

// Analyzer performs canary analysis
type Analyzer struct {
 config    AnalysisConfig
 metrics   MetricsProvider
 logger    *zap.Logger
}

// MetricsProvider provides metric queries
type MetricsProvider interface {
 Query(ctx context.Context, query string, start, end time.Time) (float64, error)
}

// NewAnalyzer creates a new canary analyzer
func NewAnalyzer(config AnalysisConfig, metrics MetricsProvider, logger *zap.Logger) *Analyzer {
 return &Analyzer{
  config:  config,
  metrics: metrics,
  logger:  logger,
 }
}

// Analyze performs canary analysis
func (a *Analyzer) Analyze(ctx context.Context, canaryVersion, baselineVersion string) AnalysisResult {
 now := time.Now()
 start := now.Add(-a.config.Interval)

 successCount := 0
 failureCount := 0

 for _, threshold := range a.config.Thresholds {
  // Query canary metric
  canaryQuery := fmt.Sprintf(threshold.Query, canaryVersion)
  canaryValue, err := a.metrics.Query(ctx, canaryQuery, start, now)
  if err != nil {
   a.logger.Error("Failed to query canary metric",
    zap.String("metric", threshold.Name),
    zap.Error(err))
   failureCount++
   continue
  }

  var compareValue float64
  if threshold.Baseline {
   // Query baseline metric for comparison
   baselineQuery := fmt.Sprintf(threshold.Query, baselineVersion)
   compareValue, err = a.metrics.Query(ctx, baselineQuery, start, now)
   if err != nil {
    a.logger.Error("Failed to query baseline metric",
     zap.String("metric", threshold.Name),
     zap.Error(err))
    failureCount++
    continue
   }
  }

  // Evaluate threshold
  passed := a.evaluateThreshold(threshold, canaryValue, compareValue)
  if passed {
   successCount++
  } else {
   failureCount++
   a.logger.Warn("Canary threshold failed",
    zap.String("metric", threshold.Name),
    zap.Float64("canary", canaryValue),
    zap.Float64("baseline", compareValue))
  }
 }

 // Determine overall result
 if failureCount >= a.config.MaxFailures {
  return AnalysisResultFail
 }

 if successCount >= a.config.RequiredSuccessCount {
  return AnalysisResultSuccess
 }

 if failureCount > 0 {
  return AnalysisResultWarning
 }

 return AnalysisResultPending
}

func (a *Analyzer) evaluateThreshold(threshold MetricThreshold, canaryValue, baselineValue float64) bool {
 value := canaryValue

 // If comparing to baseline, calculate relative value
 if threshold.Baseline && baselineValue > 0 {
  value = (canaryValue - baselineValue) / baselineValue * 100
 }

 // Check min threshold
 if threshold.Min != nil && value < *threshold.Min {
  return false
 }

 // Check max threshold
 if threshold.Max != nil && value > *threshold.Max {
  return false
 }

 return true
}

// Standard thresholds for common metrics
func DefaultThresholds() []MetricThreshold {
 maxErrorRate := 1.0  // 1%
 maxLatencyDelta := 20.0  // 20% increase

 return []MetricThreshold{
  {
   Name:     "error_rate",
   Query:    "sum(rate(http_requests_total{version=\"%s\",status=~\"5..\"}[5m])) / sum(rate(http_requests_total{version=\"%s\"}[5m])) * 100",
   Max:      &maxErrorRate,
   Baseline: false,
  },
  {
   Name:     "latency_p99",
   Query:    "histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{version=\"%s\"}[5m])) by (le))",
   Max:      &maxLatencyDelta,
   Baseline: true,
  },
 }
}
```

---

## 5. Failure Scenarios and Mitigations

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Canary Crash** | 5-25% errors | Error rate spike | Automatic rollback |
| **Performance Regression** | Slow responses | Latency increase | Threshold-based rollback |
| **Metric False Positive** | Unnecessary rollback | Inconsistent signals | Multiple metric correlation |
| **Stuck Promotion** | Incomplete rollout | Timeout | Manual intervention |

---

## 6. Semantic Trade-off Analysis

| Aspect | Canary | Blue-Green | Rolling Update |
|--------|--------|------------|----------------|
| **Risk Level** | Low | Low | Medium |
| **Resource Cost** | Medium | High | Low |
| **Rollback Speed** | Fast | Instant | Medium |
| **User Impact** | Minimal | None | Some |
| **Complexity** | High | Medium | Low |

---

## 7. References

1. Humble, J., & Farley, D. (2010). *Continuous Delivery*. Addison-Wesley.
2. Governor, J. (2018). Progressive Delivery. RedMonk.
3. FluxCD. (2024). *Flagger Documentation*. flagger.app.
4. Netflix. (2024). *Spinnaker Canary Documentation*. spinnaker.io.
