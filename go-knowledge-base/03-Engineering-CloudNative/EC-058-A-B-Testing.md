# EC-058: A/B Testing Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #ab-testing #experimentation #data-driven #conversion-optimization #hypothesis-testing
> **Authoritative Sources**:
>
> - [Trustworthy Online Controlled Experiments](https://experimentguide.com/) - Kohavi et al. (2020)
> - [Lean Analytics](https://www.oreilly.com/library/view/lean-analytics/9781449334915/) - Croll & Yoskovitz (2013)
> - [Online Controlled Experiments at Scale](https://ai.stanford.edu/~ronnyk/2015%20Controlled%20Experiments%20on%20the%20Web%20Survey.pdf) - Kohavi et al. (2009)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (A/B Testing Domain)**
Let $\mathcal{E}$ be an experiment with:

- Control variant $A$ (current version)
- Treatment variant $B$ (new version)
- Traffic split $\pi_A + \pi_B = 1$
- Metric $M$ measuring experiment success

**Statistical Requirements:**

| Requirement | Formal Definition | Impact |
|-------------|-------------------|--------|
| **Randomization** | $P(user \in A) = \pi_A \land P(user \in B) = \pi_B$ | Eliminates selection bias |
| **Sample Size** | $n > n_{min}(\alpha, \beta, \delta)$ | Statistical power |
| **Independence** | $\forall u_i, u_j: variant(u_i) \perp variant(u_j)$ | Valid variance estimation |
| **SRM Check** | $\chi^2(actual, expected) < threshold$ | Detects assignment issues |

### 1.2 Problem Statement

**Problem 1.1 (Causal Inference)**
Given experiment $\mathcal{E}$, estimate the Average Treatment Effect (ATE):

$$ATE = E[M|B] - E[M|A]$$

With confidence interval $CI_{95\%}$ that doesn't include 0.

**Key Challenges:**

1. **Sample Ratio Mismatch**: Uneven traffic split
2. **Novelty Effect**: Initial bias in user behavior
3. **Network Effects**: Users affecting each other
4. **Multiple Testing**: Increased false positive rate
5. **Segment Analysis**: Spurious findings in subgroups

---

## 2. Solution Architecture

### 2.1 Experiment Lifecycle

| Phase | Duration | Activity |
|-------|----------|----------|
| **Design** | 1-2 days | Define hypothesis, metrics, sample size |
| **Setup** | 1 day | Configure variants, targeting |
| **Run** | 1-4 weeks | Collect data, monitor health |
| **Analysis** | 2-3 days | Calculate results, validate |
| **Decision** | 1 day | Ship, iterate, or abandon |

### 2.2 Statistical Framework

- **Significance Level**: $\alpha = 0.05$ (5% false positive rate)
- **Power**: $1 - \beta = 0.8$ (80% chance to detect true effect)
- **Minimum Detectable Effect**: $\delta$ (practical significance)
- **Two-tailed t-test**: For comparing means

---

## 3. Visual Representations

### 3.1 A/B Testing Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    A/B TESTING ARCHITECTURE                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  EXPERIMENT DEFINITION                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Experiment Configuration                         │   │
│  │                                                                      │   │
│  │  Name: "New Checkout Flow"                                           │   │
│  │  Hypothesis: "Simplified checkout will increase conversion by 5%"    │   │
│  │                                                                      │   │
│  │  Variants:                                                           │   │
│  │  • A (Control): Current checkout - 50% traffic                       │   │
│  │  • B (Treatment): New checkout - 50% traffic                         │   │
│  │                                                                      │   │
│  │  Primary Metric: Purchase completion rate                            │   │
│  │  Secondary Metrics: Avg order value, Time to checkout                │   │
│  │  Guardrail Metrics: Error rate, Support contacts                     │   │
│  │                                                                      │   │
│  │  Targeting: All users, Desktop only                                  │   │
│  │  Duration: 2 weeks (minimum sample: 10,000 per variant)              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  TRAFFIC ASSIGNMENT                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Assignment Service                              │   │
│  │                                                                      │   │
│  │  User-123 ──► hash(user_id + experiment_id) % 100 = 42 ──► Variant B │   │
│  │                                                                      │   │
│  │  Assignment Logic:                                                   │   │
│  │  • Consistent hashing (same user always same variant)                │   │
│  │  • Sticky assignment (survives sessions)                             │   │
│  │  • Stratified sampling (preserve user segments)                      │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│         ┌──────────────────────────┴──────────────────────────┐            │
│         │                          │                          │            │
│         ▼                          ▼                          ▼            │
│  ┌───────────────┐          ┌───────────────┐          ┌───────────────┐   │
│  │   VARIANT A   │          │   VARIANT B   │          │   ANALYTICS   │   │
│  │   (Control)   │          │  (Treatment)  │          │               │   │
│  │               │          │               │          │  Event Store  │   │
│  │  ┌─────────┐  │          │  ┌─────────┐  │          │               │   │
│  │  │ User    │  │          │  │ User    │  │─────────►│  • Exposures  │   │
│  │  │ Session │  │          │  │ Session │  │          │  • Conversions│   │
│  │  └─────────┘  │          │  └─────────┘  │          │  • Metrics    │   │
│  │       │       │          │       │       │          │               │   │
│  │       ▼       │          │       ▼       │          └───────┬───────┘   │
│  │  ┌─────────┐  │          │  ┌─────────┐  │                  │           │
│  │  │ Events  │──┼──────────┼─►│ Events  │──┘                  ▼           │
│  │  └─────────┘  │          │  └─────────┘  │          ┌───────────────┐   │
│  │               │          │               │          │  Stats Engine │   │
│  └───────────────┘          └───────────────┘          │               │   │
│                                                        │ • Calculate   │   │
│                                                        │   p-values    │   │
│                                                        │ • Confidence  │   │
│                                                        │   intervals   │   │
│                                                        │ • Sequential  │   │
│                                                        │   testing     │   │
│                                                        └───────┬───────┘   │
│                                                                │           │
│  RESULTS                                                       ▼           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Experiment Results                               │   │
│  │                                                                      │   │
│  │  Variant A: 12.5% conversion (n=10,523)                              │   │
│  │  Variant B: 13.2% conversion (n=10,489)                              │   │
│  │                                                                      │   │
│  │  Relative Lift: +5.6%                                                │   │
│  │  P-value: 0.032 (statistically significant)                          │   │
│  │  95% CI: [+0.4%, +10.8%]                                             │   │
│  │                                                                      │   │
│  │  Decision: ✅ SHIP Variant B                                         │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Statistical Analysis Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    STATISTICAL ANALYSIS FLOW                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  RAW DATA                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Variant A: [0, 1, 0, 1, 1, 0, 1, 0, 0, 1, ...] (n=10,523)          │   │
│  │  Variant B: [0, 1, 1, 1, 0, 1, 1, 1, 0, 1, ...] (n=10,489)          │   │
│  │  (0 = no conversion, 1 = conversion)                                 │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  DESCRIPTIVE STATISTICS                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  Variant A:                                                          │   │
│  │    • Mean: 0.125 (12.5%)                                             │   │
│  │    • Std: 0.331                                                      │   │
│  │    • N: 10,523                                                       │   │
│  │                                                                      │   │
│  │  Variant B:                                                          │   │
│  │    • Mean: 0.132 (13.2%)                                             │   │
│  │    • Std: 0.338                                                      │   │
│  │    • N: 10,489                                                       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  HYPOTHESIS TEST                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  H0: μB - μA = 0  (no difference)                                    │   │
│  │  H1: μB - μA ≠ 0  (there is a difference)                            │   │
│  │                                                                      │   │
│  │  Two-proportion z-test:                                              │   │
│  │                                                                      │   │
│  │  z = (pB - pA) / sqrt(p(1-p)(1/nA + 1/nB))                           │   │
│  │                                                                      │   │
│  │  z = (0.132 - 0.125) / sqrt(0.128(0.872)(1/10523 + 1/10489))         │   │
│  │  z = 2.14                                                            │   │
│  │                                                                      │   │
│  │  p-value = 0.032                                                     │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  DECISION                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  p-value (0.032) < α (0.05)?  YES                                    │   │
│  │                                                                      │   │
│  │  → Reject H0                                                         │   │
│  │  → Result is statistically significant                               │   │
│  │                                                                      │   │
│  │  Practical significance: +5.6% lift                                  │   │
│  │  Business impact: +$500K/year (projected)                            │   │
│  │                                                                      │   │
│  │  Decision: Ship Variant B                                            │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Experiment Assignment

```go
package experiment

import (
 "crypto/sha256"
 "encoding/binary"
 "fmt"
 "sync"
 "time"
)

// Experiment represents an A/B test
type Experiment struct {
 ID          string
 Name        string
 Status      Status
 Variants    []Variant
 TrafficSplit map[string]float64 // variant_id -> percentage
 StartTime   time.Time
 EndTime     *time.Time
}

type Status int

const (
 StatusDraft Status = iota
 StatusRunning
 StatusPaused
 StatusCompleted
)

// Variant represents an experiment variant
type Variant struct {
 ID          string
 Name        string
 Description string
 Weight      int // Relative weight (not percentage)
}

// Assignment represents a user's assignment to a variant
type Assignment struct {
 UserID        string
 ExperimentID  string
 VariantID     string
 AssignedAt    time.Time
}

// Assigner handles experiment assignments
type Assigner struct {
 experiments map[string]*Experiment
 assignments map[string]Assignment // cache: user_id:experiment_id -> Assignment
 mu          sync.RWMutex
}

// NewAssigner creates a new experiment assigner
func NewAssigner() *Assigner {
 return &Assigner{
  experiments: make(map[string]*Experiment),
  assignments: make(map[string]Assignment),
 }
}

// RegisterExperiment registers an experiment
func (a *Assigner) RegisterExperiment(exp *Experiment) {
 a.mu.Lock()
 defer a.mu.Unlock()
 a.experiments[exp.ID] = exp
}

// Assign assigns a user to a variant
func (a *Assigner) Assign(userID, experimentID string) (string, error) {
 // Check cache first
 cacheKey := fmt.Sprintf("%s:%s", userID, experimentID)
 a.mu.RLock()
 if assignment, ok := a.assignments[cacheKey]; ok {
  a.mu.RUnlock()
  return assignment.VariantID, nil
 }
 a.mu.RUnlock()

 // Get experiment
 a.mu.RLock()
 exp, ok := a.experiments[experimentID]
 a.mu.RUnlock()

 if !ok {
  return "", fmt.Errorf("experiment not found: %s", experimentID)
 }

 if exp.Status != StatusRunning {
  return "", fmt.Errorf("experiment not running: %s", experimentID)
 }

 // Deterministic assignment using hash
 variantID := a.deterministicAssign(userID, exp)

 // Cache assignment
 assignment := Assignment{
  UserID:       userID,
  ExperimentID: experimentID,
  VariantID:    variantID,
  AssignedAt:   time.Now(),
 }

 a.mu.Lock()
 a.assignments[cacheKey] = assignment
 a.mu.Unlock()

 return variantID, nil
}

func (a *Assigner) deterministicAssign(userID string, exp *Experiment) string {
 // Create deterministic hash
 h := sha256.New()
 h.Write([]byte(fmt.Sprintf("%s:%s", exp.ID, userID)))
 hash := h.Sum(nil)

 // Convert to int in range [0, 9999]
 val := binary.BigEndian.Uint32(hash[:4]) % 10000

 // Find variant based on traffic split
 cumulative := 0
 for _, variant := range exp.Variants {
  percentage := int(exp.TrafficSplit[variant.ID] * 100)
  cumulative += percentage
  if int(val) < cumulative {
   return variant.ID
  }
 }

 // Fallback to control
 return exp.Variants[0].ID
}

// GetVariant returns the variant details
func (a *Assigner) GetVariant(experimentID, variantID string) (*Variant, error) {
 a.mu.RLock()
 defer a.mu.RUnlock()

 exp, ok := a.experiments[experimentID]
 if !ok {
  return nil, fmt.Errorf("experiment not found")
 }

 for _, v := range exp.Variants {
  if v.ID == variantID {
   return &v, nil
  }
 }

 return nil, fmt.Errorf("variant not found")
}
```

### 4.2 Statistical Analysis

```go
package experiment

import (
 "math"
)

// Result represents experiment results
type Result struct {
 VariantID     string
 SampleSize    int
 Conversions   int
 ConversionRate float64
}

// Analysis represents statistical analysis
type Analysis struct {
 Control     Result
 Treatment   Result
 Lift        float64
 PValue      float64
 CI Lower    float64
 CI Upper    float64
 Significant bool
}

// Analyze performs statistical analysis
func Analyze(control, treatment Result) *Analysis {
 // Calculate conversion rates
 p1 := control.ConversionRate
 p2 := treatment.ConversionRate
 n1 := float64(control.SampleSize)
 n2 := float64(treatment.SampleSize)

 // Pooled proportion
 pPool := (float64(control.Conversions) + float64(treatment.Conversions)) / (n1 + n2)

 // Standard error
 se := math.Sqrt(pPool * (1 - pPool) * (1/n1 + 1/n2))

 // Z-score
 z := (p2 - p1) / se

 // P-value (two-tailed)
 pValue := 2 * (1 - normalCDF(math.Abs(z)))

 // 95% Confidence interval
 seDiff := math.Sqrt(p1*(1-p1)/n1 + p2*(1-p2)/n2)
 margin := 1.96 * seDiff
 ciLower := (p2 - p1) - margin
 ciUpper := (p2 - p1) + margin

 // Relative lift
 lift := (p2 - p1) / p1 * 100

 return &Analysis{
  Control:     control,
  Treatment:   treatment,
  Lift:        lift,
  PValue:      pValue,
  CI Lower:    ciLower * 100,
  CI Upper:    ciUpper * 100,
  Significant: pValue < 0.05,
 }
}

// normalCDF calculates the cumulative distribution function
func normalCDF(x float64) float64 {
 return 0.5 * (1 + math.Erf(x/math.Sqrt(2)))
}
```

---

## 5. Failure Scenarios and Mitigations

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **SRM** | Biased results | Chi-square test | Auto-pause experiment |
| **Novelty Effect** | False positive | Time-based analysis | Run longer, segment by time |
| **Network Effects** | Contaminated results | Cross-user metrics | Cluster-based analysis |
| **Peeking** | False positive | Multiple comparison correction | Sequential testing |

---

## 6. Semantic Trade-off Analysis

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    EXPERIMENT DESIGN TRADE-OFFS                              │
├─────────────────────┬─────────────────┬─────────────────────────────────────┤
│     Dimension       │   Short Run     │            Long Run                 │
├─────────────────────┼─────────────────┼─────────────────────────────────────┤
│ Statistical Power   │ Lower           │ Higher                              │
│ Novelty Effect      │ Higher          │ Lower                               │
│ Business Risk       │ Lower           │ Higher (if negative)                │
│ Time to Decision    │ Faster          │ Slower                              │
│ Seasonality Bias    │ Higher          │ Lower                               │
└─────────────────────┴─────────────────┴─────────────────────────────────────┘
```

---

## 7. References

1. Kohavi, R., Tang, D., & Xu, Y. (2020). *Trustworthy Online Controlled Experiments*. Cambridge University Press.
2. Croll, A., & Yoskovitz, B. (2013). *Lean Analytics*. O'Reilly Media.
3. Kohavi, R., et al. (2009). Online Controlled Experiments at Scale. *KDD*.
