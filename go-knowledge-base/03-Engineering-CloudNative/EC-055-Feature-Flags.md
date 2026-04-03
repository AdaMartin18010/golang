# EC-055: Feature Flags Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #feature-flags #feature-toggles #canary-releases #continuous-deployment #launchdarkly
> **Authoritative Sources**:
>
> - [Feature Toggles](https://martinfowler.com/articles/feature-toggles.html) - Fowler (2017)
> - [Continuous Delivery](https://continuousdelivery.com/) - Humble & Farley (2010)
> - [LaunchDarkly Documentation](https://docs.launchdarkly.com/) - LaunchDarkly (2024)
> - [Unleash Documentation](https://docs.getunleash.io/) - Unleash (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Feature Flag Domain)**
Let $\mathcal{F} = \{f_1, f_2, ..., f_n\}$ be a set of feature flags for application $A$ where each flag $f_i$ has:

- State $s(f_i) \in \{off, on, partial\}$
- Targeting rules $\mathcal{R}(f_i)$ defining activation conditions
- Variants $\mathcal{V}(f_i)$ for A/B testing

**Flag Types:**

| Type | Purpose | Lifecycle | Targeting |
|------|---------|-----------|-----------|
| **Release** | Dark launch | Short | User segments |
| **Experiment** | A/B testing | Medium | Percentage |
| **Operational** | Kill switch | Permanent | Emergency only |
| **Permission** | Entitlement | Long | User attributes |

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Evaluation Latency** | $T_{eval}(f) < 10ms$ | Cannot impact request performance |
| **Consistency** | $\forall req: eval(f, req) = deterministic$ | Same user must see same feature |
| **Propagation Speed** | $T_{flag\_change} < 1min$ | Quick enable/disable required |
| **Cardinality** | $|\mathcal{F}| \leq 1000$ | Too many flags create complexity |

### 1.2 Problem Statement

**Problem 1.1 (Feature Flag Evaluation)**
Given flag $f$, context $ctx$, and rules $\mathcal{R}(f)$, determine:

$$eval(f, ctx) = \begin{cases} true & \exists r \in \mathcal{R}(f): match(r, ctx) \\ false & otherwise \end{cases}$$

**Key Challenges:**

1. **Performance**: Evaluating flags without adding latency
2. **Consistency**: Same user always sees same feature variant
3. **Safety**: Ability to quickly disable problematic features
4. **Testing**: Validating behavior with flags on/off
5. **Cleanup**: Removing flags after feature release

---

## 2. Solution Architecture

### 2.1 Feature Flag Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FEATURE FLAG TAXONOMY                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    RELEASE FLAGS (Short-lived)                       │   │
│  │                                                                      │   │
│  │  Purpose: Decouple deployment from release                          │   │
│  │  Lifecycle: Days to weeks                                           │   │
│  │                                                                      │   │
│  │  Example:                                                            │   │
│  │  ┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐       │   │
│  │  │  Code   │────►│ Deploy  │────►│  Test   │────►│ Release │       │   │
│  │  │ Written │     │ to Prod │     │  (off)  │     │ (on)    │       │   │
│  │  └─────────┘     └─────────┘     └─────────┘     └─────────┘       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    EXPERIMENT FLAGS (Medium-lived)                   │   │
│  │                                                                      │   │
│  │  Purpose: A/B testing, data-driven decisions                        │   │
│  │  Lifecycle: Weeks to months                                         │   │
│  │                                                                      │   │
│  │  Example:                                                            │   │
│  │  ┌─────────┐     ┌─────────┐     ┌─────────┐                       │   │
│  │  │ Variant │     │ Variant │     │ Winner  │                       │   │
│  │  │   A     │     │   B     │     │ Deployed│                       │   │
│  │  │ (50%)   │     │ (50%)   │     │ (100%)  │                       │   │
│  │  └─────────┘     └─────────┘     └─────────┘                       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    OPERATIONAL FLAGS (Permanent)                     │   │
│  │                                                                      │   │
│  │  Purpose: Circuit breakers, kill switches                           │   │
│  │  Lifecycle: Permanent                                               │   │
│  │                                                                      │   │
│  │  Example:                                                            │   │
│  │  ┌─────────┐     ┌─────────┐                                        │   │
│  │  │ Feature │────►│ Problem │────► [Kill Switch] ───► Disabled       │   │
│  │  │ Enabled │     │ Detected│                                        │   │
│  │  └─────────┘     └─────────┘                                        │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    PERMISSION FLAGS (Long-lived)                     │   │
│  │                                                                      │   │
│  │  Purpose: Entitlement, tier-based features                          │   │
│  │  Lifecycle: Months to years                                         │   │
│  │                                                                      │   │
│  │  Example:                                                            │   │
│  │  ┌─────────┐     ┌─────────┐     ┌─────────┐                       │   │
│  │  │  Free   │     │  Pro    │     │Enterprise                          │   │
│  │  │ (basic) │     │(advanced)│    │ (full)  │                       │   │
│  │  └─────────┘     └─────────┘     └─────────┘                       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Visual Representations

### 3.1 Feature Flag System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FEATURE FLAG SYSTEM ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  MANAGEMENT LAYER                                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Feature Flag Service                            │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │   Admin UI  │  │   API       │  │   SDK Keys  │                  │   │
│  │  │             │  │             │  │             │                  │   │
│  │  │• Create flag│  │• CRUD flags │  │• Client SDK│                  │   │
│  │  │• Set rules  │  │• Evaluate   │  │• Server SDK│                  │   │
│  │  │• View stats │  │• Analytics  │  │• Mobile SDK│                  │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                  │   │
│  │         │                │                │                         │   │
│  │         └────────────────┼────────────────┘                         │   │
│  │                          ▼                                         │   │
│  │                 ┌─────────────────┐                                 │   │
│  │                 │  Flag Storage   │                                 │   │
│  │                 │  (Redis/etcd)   │                                 │   │
│  │                 └────────┬────────┘                                 │   │
│  │                          │                                         │   │
│  └──────────────────────────┼─────────────────────────────────────────┘   │
│                             │                                             │
│  STREAMING LAYER            │ WebSocket / SSE                            │
│                             │                                             │
│  ┌──────────────────────────┼─────────────────────────────────────────┐   │
│  │                          ▼                                         │   │
│  │                     Flag Updates                                   │   │
│  │  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐        │   │
│  │  │    SDK      │◄────►│    SDK      │◄────►│    SDK      │        │   │
│  │  │   Client    │      │   Client    │      │   Client    │        │   │
│  │  │             │      │             │      │             │        │   │
│  │  │• Local cache│      │• Local cache│      │• Local cache│        │   │
│  │  │• Real-time  │      │• Real-time  │      │• Real-time  │        │   │
│  │  │  updates    │      │  updates    │      │  updates    │        │   │
│  │  └──────┬──────┘      └──────┬──────┘      └──────┬──────┘        │   │
│  │         │                    │                    │                │   │
│  └─────────┼────────────────────┼────────────────────┼────────────────┘   │
│            │                    │                    │                    │
│  APPLICATION LAYER                                                    │    │
│            │                    │                    │                    │
│  ┌─────────┴────────────────────┴────────────────────┴────────────────┐   │
│  │                                                                     │   │
│  │  Application Code:                                                  │   │
│  │                                                                     │   │
│  │  if (flags.IsEnabled("new-checkout-flow", user)) {                  │   │
│  │      // New feature code                                            │   │
│  │      renderNewCheckout(user);                                       │   │
│  │  } else {                                                           │   │
│  │      // Old code path                                               │   │
│  │      renderOldCheckout(user);                                       │   │
│  │  }                                                                  │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ANALYTICS LAYER                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Analytics Pipeline                               │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │   │
│  │  │  Event      │  │  Aggregate  │  │  Dashboard  │                  │   │
│  │  │  Collection │──►│  Metrics    │──►│  & Reports  │                  │   │
│  │  │             │  │             │  │             │                  │   │
│  │  │• Flag evals │  │• Conversion │  │• A/B results│                  │   │
│  │  │• User segments│ │• Impact     │  │• Usage stats│                  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Flag Evaluation Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLAG EVALUATION FLOW                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Request: Is flag "new-feature" enabled for user-123?                       │
│                                                                             │
│  Step 1: Check local cache                                                  │
│  ┌─────────────┐                                                            │
│  │  SDK Cache  │                                                            │
│  │             │                                                            │
│  │ new-feature │──► User not in cache                                       │
│  │ default: off│     (stale-while-revalidate)                               │
│  └─────────────┘                                                            │
│         │                                                                   │
│         │ Cache miss / stale                                                │
│         ▼                                                                   │
│  Step 2: Evaluate targeting rules                                           │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Targeting Rules for "new-feature"                 │   │
│  │                                                                      │   │
│  │  Rule 1: user.id == "admin@example.com" ──► true (override)         │   │
│  │       No match                                                         │   │
│  │                                                                      │   │
│  │  Rule 2: user.group == "beta-testers" ──► true (50% rollout)        │   │
│  │       Checking...                                                      │   │
│  │       user-123.groups = ["premium", "beta-testers"] ✓ MATCH           │   │
│  │                                                                      │   │
│  │  Rule 3: percentage rollout 10% ──► false (not in range)            │   │
│  │       (Skipped - Rule 2 matched)                                       │   │
│  │                                                                      │   │
│  │  Result: MATCH on Rule 2                                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│         │                                                                   │
│         ▼                                                                   │
│  Step 3: Determine variant                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Variant Assignment                                │   │
│  │                                                                      │   │
│  │  Hash: sha256("new-feature:user-123") % 100 = 42                    │   │
│  │                                                                      │   │
│  │  Variants:                                                           │   │
│  │  • "control"  : 0-49  (50%)                                          │   │
│  │  • "treatment": 50-99 (50%)                                          │   │
│  │                                                                      │   │
│  │  42 falls in "control" range                                         │   │
│  │                                                                      │   │
│  │  Result: true (enabled), variant: "control"                          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│         │                                                                   │
│         ▼                                                                   │
│  Step 4: Update cache and return                                            │
│  ┌─────────────┐      ┌─────────────┐                                       │
│  │  SDK Cache  │      │  Analytics  │                                       │
│  │             │      │  (async)    │                                       │
│  │ user-123:   │      │             │                                       │
│  │  new-feature│      │ Flag eval:  │                                       │
│  │  = true     │      │  user-123   │                                       │
│  │  variant:   │      │  flag: new  │                                       │
│  │  control    │      │  result:    │                                       │
│  │  expires:   │      │  enabled    │                                       │
│  │  +5min      │      │  variant:   │                                       │
│  └─────────────┘      │  control    │                                       │
│                       └─────────────┘                                       │
│                                                                             │
│  Response: {enabled: true, variant: "control"}                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Feature Flag Client

```go
package featureflags

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "strconv"
 "sync"
 "time"

 "go.uber.org/zap"
)

// Flag represents a feature flag
type Flag struct {
 Key          string            `json:"key"`
 Enabled      bool              `json:"enabled"`
 DefaultValue bool              `json:"default_value"`
 Rules        []Rule            `json:"rules,omitempty"`
 Variants     map[string]int    `json:"variants,omitempty"` // variant -> percentage

 // Metadata
 Description string    `json:"description"`
 CreatedAt   time.Time `json:"created_at"`
 UpdatedAt   time.Time `json:"updated_at"`
}

// Rule defines targeting criteria
type Rule struct {
 ID       string      `json:"id"`
 Name     string      `json:"name"`
 Priority int         `json:"priority"`
 Conditions []Condition `json:"conditions"`
 Action   RuleAction  `json:"action"`
}

// Condition defines a single condition
type Condition struct {
 Attribute string      `json:"attribute"`
 Operator  string      `json:"operator"` // eq, neq, in, contains, gt, lt
 Value     interface{} `json:"value"`
}

// RuleAction defines the action when rule matches
type RuleAction struct {
 Type      string            `json:"type"` // enable, disable, variant
 Variant   string            `json:"variant,omitempty"`
 Percentage int              `json:"percentage,omitempty"` // for gradual rollout
}

// Context provides user/session context for evaluation
type Context struct {
 UserID    string                 `json:"user_id"`
 SessionID string                 `json:"session_id"`
 Groups    []string               `json:"groups,omitempty"`
 Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// EvalResult is the result of flag evaluation
type EvalResult struct {
 Enabled bool   `json:"enabled"`
 Variant string `json:"variant,omitempty"`
 RuleID  string `json:"rule_id,omitempty"`
}

// Client manages feature flags
type Client struct {
 config    Config
 store     Store
 cache     *cache
 logger    *zap.Logger

 mu        sync.RWMutex
 flags     map[string]*Flag

 stopCh    chan struct{}
}

// Config holds client configuration
type Config struct {
 SDKKey           string
 APIEndpoint      string
 RefreshInterval  time.Duration
 CacheTTL         time.Duration
 OfflineMode      bool
}

// Store defines the flag storage interface
type Store interface {
 Get(ctx context.Context, key string) (*Flag, error)
 GetAll(ctx context.Context) (map[string]*Flag, error)
 Subscribe(ctx context.Context) (<-chan FlagUpdate, error)
}

// FlagUpdate represents a flag change
type FlagUpdate struct {
 Type string // add, update, delete
 Flag *Flag
}

// cache provides in-memory caching
type cache struct {
 data map[string]cacheEntry
 mu   sync.RWMutex
 ttl  time.Duration
}

type cacheEntry struct {
 result    EvalResult
 timestamp time.Time
}

// New creates a new feature flag client
func New(cfg Config, store Store, logger *zap.Logger) (*Client, error) {
 client := &Client{
  config: cfg,
  store:  store,
  cache: &cache{
   data: make(map[string]cacheEntry),
   ttl:  cfg.CacheTTL,
  },
  logger: logger,
  flags:  make(map[string]*Flag),
  stopCh: make(chan struct{}),
 }

 // Initial load
 if err := client.refreshFlags(context.Background()); err != nil {
  logger.Warn("Failed to load initial flags", zap.Error(err))
 }

 // Start background refresh
 go client.backgroundRefresh()

 return client, nil
}

// IsEnabled checks if a flag is enabled for the given context
func (c *Client) IsEnabled(ctx context.Context, flagKey string, fctx Context) bool {
 result := c.Evaluate(ctx, flagKey, fctx)
 return result.Enabled
}

// Evaluate evaluates a flag for the given context
func (c *Client) Evaluate(ctx context.Context, flagKey string, fctx Context) EvalResult {
 // Check cache first
 cacheKey := c.cacheKey(flagKey, fctx)
 c.cache.mu.RLock()
 if entry, ok := c.cache.data[cacheKey]; ok {
  if time.Since(entry.timestamp) < c.cache.ttl {
   c.cache.mu.RUnlock()
   return entry.result
  }
 }
 c.cache.mu.RUnlock()

 // Get flag definition
 c.mu.RLock()
 flag, ok := c.flags[flagKey]
 c.mu.RUnlock()

 if !ok {
  // Flag not found, return default
  return EvalResult{Enabled: false}
 }

 // Evaluate flag
 result := c.evaluateFlag(flag, fctx)

 // Update cache
 c.cache.mu.Lock()
 c.cache.data[cacheKey] = cacheEntry{
  result:    result,
  timestamp: time.Now(),
 }
 c.cache.mu.Unlock()

 return result
}

func (c *Client) evaluateFlag(flag *Flag, fctx Context) EvalResult {
 // If flag is disabled globally, return false
 if !flag.Enabled {
  return EvalResult{Enabled: false}
 }

 // Check rules in priority order
 for _, rule := range flag.Rules {
  if c.matchesRule(rule, fctx) {
   return c.applyRule(rule, flag, fctx)
  }
 }

 // No rules matched, return default
 return EvalResult{Enabled: flag.DefaultValue}
}

func (c *Client) matchesRule(rule Rule, fctx Context) bool {
 for _, cond := range rule.Conditions {
  if !c.matchesCondition(cond, fctx) {
   return false
  }
 }
 return true
}

func (c *Client) matchesCondition(cond Condition, fctx Context) bool {
 value := c.getAttributeValue(cond.Attribute, fctx)

 switch cond.Operator {
 case "eq":
  return value == cond.Value
 case "neq":
  return value != cond.Value
 case "in":
  if arr, ok := cond.Value.([]interface{}); ok {
   for _, v := range arr {
    if value == v {
     return true
    }
   }
  }
  return false
 case "contains":
  if str, ok := value.(string); ok {
   if substr, ok := cond.Value.(string); ok {
    return contains(str, substr)
   }
  }
  return false
 default:
  return false
 }
}

func (c *Client) applyRule(rule Rule, flag *Flag, fctx Context) EvalResult {
 switch rule.Action.Type {
 case "enable":
  if rule.Action.Percentage > 0 {
   // Gradual rollout
   hash := c.hashPercentage(flag.Key, fctx.UserID)
   enabled := hash < rule.Action.Percentage
   return EvalResult{Enabled: enabled, RuleID: rule.ID}
  }
  return EvalResult{Enabled: true, RuleID: rule.ID}

 case "disable":
  return EvalResult{Enabled: false, RuleID: rule.ID}

 case "variant":
  variant := c.selectVariant(flag, fctx)
  return EvalResult{Enabled: true, Variant: variant, RuleID: rule.ID}

 default:
  return EvalResult{Enabled: flag.DefaultValue}
 }
}

func (c *Client) selectVariant(flag *Flag, fctx Context) string {
 if len(flag.Variants) == 0 {
  return ""
 }

 hash := c.hashPercentage(flag.Key, fctx.UserID)
 cumulative := 0

 for variant, percentage := range flag.Variants {
  cumulative += percentage
  if hash < cumulative {
   return variant
  }
 }

 return ""
}

func (c *Client) hashPercentage(flagKey, userID string) int {
 h := sha256.New()
 h.Write([]byte(flagKey + ":" + userID))
 hash := hex.EncodeToString(h.Sum(nil))

 // Take first 8 chars and convert to int
 val, _ := strconv.ParseInt(hash[:8], 16, 64)
 return int(val % 100)
}

func (c *Client) getAttributeValue(attr string, fctx Context) interface{} {
 switch attr {
 case "user.id":
  return fctx.UserID
 case "user.groups":
  return fctx.Groups
 default:
  return fctx.Attributes[attr]
 }
}

func (c *Client) cacheKey(flagKey string, fctx Context) string {
 return fmt.Sprintf("%s:%s", flagKey, fctx.UserID)
}

func (c *Client) refreshFlags(ctx context.Context) error {
 flags, err := c.store.GetAll(ctx)
 if err != nil {
  return err
 }

 c.mu.Lock()
 c.flags = flags
 c.mu.Unlock()

 return nil
}

func (c *Client) backgroundRefresh() {
 ticker := time.NewTicker(c.config.RefreshInterval)
 defer ticker.Stop()

 for {
  select {
  case <-ticker.C:
   if err := c.refreshFlags(context.Background()); err != nil {
    c.logger.Error("Failed to refresh flags", zap.Error(err))
   }
  case <-c.stopCh:
   return
  }
 }
}

// Close closes the client
func (c *Client) Close() error {
 close(c.stopCh)
 return nil
}

func contains(s, substr string) bool {
 return len(s) >= len(substr) && s[:len(substr)] == substr ||
     len(s) > len(substr) && contains(s[1:], substr)
}
```

---

## 5. Failure Scenarios and Mitigations

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Flag Service Down** | Cannot evaluate new flags | Connection error | Local cache + Default values |
| **Flag Leak** | Users see wrong features | Analytics anomaly | Gradual rollout + Monitoring |
| **Performance Degradation** | Request latency | Response time | Async evaluation + Caching |
| **Flag Sprawl** | Technical debt | Flag count | TTL policies + Cleanup automation |
| **Inconsistent Variant** | Wrong A/B results | Hash collision detection | Consistent hashing |

---

## 6. Semantic Trade-off Analysis

| Aspect | Client-Side | Server-Side | Hybrid |
|--------|-------------|-------------|--------|
| **Latency** | Zero network | Network call | Cached |
| **Security** | Exposed rules | Protected | Protected |
| **Consistency** | Harder | Easier | Balanced |
| **Flexibility** | Limited | Full | Full |

---

## 7. References

1. Fowler, M. (2017). *Feature Toggles*. martinfowler.com.
2. Humble, J., & Farley, D. (2010). *Continuous Delivery*. Addison-Wesley.
3. LaunchDarkly. (2024). *Feature Flag Best Practices*.
4. Unleash. (2024). *Feature Toggle Documentation*.
