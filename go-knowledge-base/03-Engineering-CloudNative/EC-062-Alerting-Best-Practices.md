# Alerting Best Practices

> **分类**: 工程与云原生
> **标签**: #alerting #monitoring #sre #oncall #incident-response
> **参考**: Google SRE, Prometheus Alerting, PagerDuty Best Practices

---

## 1. Formal Definition

### 1.1 What is Alerting?

Alerting is the systematic process of notifying responsible parties when a system's observable state deviates from defined Service Level Objectives (SLOs) or acceptable operational parameters. Effective alerting is a cornerstone of Site Reliability Engineering (SRE) and enables rapid incident response.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Alerting System Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│   │   Metrics   │───→│   Rules     │───→│  Alert      │───→│ Notification│ │
│   │   Sources   │    │   Engine    │    │  Manager    │    │  Router     │ │
│   │             │    │             │    │             │    │             │ │
│   │ • Prometheus│    │ • PromQL    │    │ • Grouping  │    │ • PagerDuty │ │
│   │ • InfluxDB  │    │ • LogQL     │    │ • Inhibition│    │ • Slack     │ │
│   │ • CloudWatch│    │ • SignalFX  │    │ • Silencing │    │ • Email     │ │
│   │ • Datadog   │    │             │    │ • Routing   │    │ • Webhook   │ │
│   └─────────────┘    └─────────────┘    └──────┬──────┘    └──────┬──────┘ │
│                                                 │                  │       │
│                                                 ↓                  ↓       │
│                                        ┌─────────────────┐  ┌──────────┐  │
│                                        │  Alert Storage  │  │ On-Call  │  │
│                                        │  & History      │  │ Engineer │  │
│                                        └─────────────────┘  └──────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Alert Types Classification

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Alert Classification System                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  BY SEVERITY                    BY SOURCE                      BY TYPE     │
│  ───────────────────            ─────────────────              ────────────│
│                                                                             │
│  P0 - Critical      ◄────────── Infrastructure ───────────────► Threshold  │
│  ├── System down                  ├── CPU/Memory/Network       ├── Static  │
│  ├── Data loss                    ├── Disk/Storage             └── Dynamic │
│  └── Security breach           Application                     Anomaly      │
│                                ├── Error rate                  Predictive  │
│  P1 - High         ◄───────────┼── Latency/P99 ─────────────► Composite   │
│  ├── Degraded performance      ├── Throughput                               │
│  ├── Capacity critical         └── Availability              SLO-Based     │
│  └── Partial outage          Business                                     │
│                              ├── Revenue drop              Event-Based      │
│  P2 - Medium       ◄───────────┼── User signups ──────────────────────────│
│  ├── Warning conditions        └── Transaction volume                      │
│  └── Approaching limits                                                  │
│                              Security                                       │
│  P3 - Low          ◄───────────┼── Failed logins ─────────────────────────│
│  ├── Informational             ├── Unauthorized access                     │
│  └── Tracking                  └── Policy violations                       │
│                                                                             │
│  P4 - None       ◄───────────  Custom/ML  ───────────────────────────────│
│  └── Dashboards only           ├── Anomaly detection                       │
│                                └── Pattern matching                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 Alert Rule Engine

```go
package alerting

import (
    "context"
    "fmt"
    "math"
    "sync"
    "time"

    "github.com/prometheus/prometheus/promql"
)

// AlertState represents the current state of an alert
type AlertState string

const (
    AlertStateInactive   AlertState = "inactive"
    AlertStatePending    AlertState = "pending"
    AlertStateFiring     AlertState = "firing"
    AlertStateResolved   AlertState = "resolved"
    AlertStateSuppressed AlertState = "suppressed"
)

// Severity represents alert severity levels
type Severity string

const (
    SeverityCritical Severity = "critical"
    SeverityHigh     Severity = "high"
    SeverityMedium   Severity = "medium"
    SeverityLow      Severity = "low"
    SeverityInfo     Severity = "info"
)

// AlertRule defines a single alerting rule
type AlertRule struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Severity    Severity          `json:"severity"`
    Query       string            `json:"query"`           // PromQL-like query
    Duration    time.Duration     `json:"duration"`        // For: duration
    Labels      map[string]string `json:"labels"`
    Annotations map[string]string `json:"annotations"`

    // Advanced options
    EnableAutoResolve bool          `json:"enable_auto_resolve"`
    ResolveDuration   time.Duration `json:"resolve_duration"`
    RepeatInterval    time.Duration `json:"repeat_interval"`
    GroupBy           []string      `json:"group_by"`

    // State (not persisted)
    state        AlertState
    stateMu      sync.RWMutex
    lastEval     time.Time
    firedAt      *time.Time
    resolvedAt   *time.Time
    value        float64
}

// AlertInstance represents a specific firing instance
type AlertInstance struct {
    RuleID      string            `json:"rule_id"`
    RuleName    string            `json:"rule_name"`
    Severity    Severity          `json:"severity"`
    State       AlertState        `json:"state"`
    Labels      map[string]string `json:"labels"`
    Annotations map[string]string `json:"annotations"`
    Value       float64           `json:"value"`
    FiredAt     time.Time         `json:"fired_at"`
    ResolvedAt  *time.Time        `json:"resolved_at,omitempty"`
    LastEval    time.Time         `json:"last_eval"`
    Fingerprint string            `json:"fingerprint"`
}

// RuleEngine evaluates alert rules against metric data
type RuleEngine struct {
    rules      map[string]*AlertRule
    rulesMu    sync.RWMutex

    evalInterval time.Duration
    queryEngine  QueryEngine
    notifier     Notifier

    instances   map[string]*AlertInstance
    instancesMu sync.RWMutex

    stopCh chan struct{}
}

// QueryEngine defines the interface for querying metrics
type QueryEngine interface {
    Query(ctx context.Context, query string) ([]TimeSeries, error)
}

// TimeSeries represents a metric time series
type TimeSeries struct {
    Labels map[string]string
    Value  float64
    Timestamp time.Time
}

// Notifier defines the interface for sending notifications
type Notifier interface {
    Notify(ctx context.Context, alert *AlertInstance) error
}

// NewRuleEngine creates a new alert rule engine
func NewRuleEngine(evalInterval time.Duration, queryEngine QueryEngine, notifier Notifier) *RuleEngine {
    if evalInterval == 0 {
        evalInterval = 15 * time.Second
    }

    return &RuleEngine{
        rules:        make(map[string]*AlertRule),
        instances:    make(map[string]*AlertInstance),
        evalInterval: evalInterval,
        queryEngine:  queryEngine,
        notifier:     notifier,
        stopCh:       make(chan struct{}),
    }
}

// RegisterRule registers a new alert rule
func (e *RuleEngine) RegisterRule(rule *AlertRule) error {
    if rule.ID == "" {
        return fmt.Errorf("rule ID is required")
    }
    if rule.Query == "" {
        return fmt.Errorf("rule query is required")
    }
    if rule.Duration == 0 {
        rule.Duration = 1 * time.Minute
    }

    e.rulesMu.Lock()
    defer e.rulesMu.Unlock()

    e.rules[rule.ID] = rule
    return nil
}

// UnregisterRule removes an alert rule
func (e *RuleEngine) UnregisterRule(ruleID string) {
    e.rulesMu.Lock()
    defer e.rulesMu.Unlock()

    delete(e.rules, ruleID)
}

// Start begins the evaluation loop
func (e *RuleEngine) Start(ctx context.Context) {
    ticker := time.NewTicker(e.evalInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            e.evaluateAll(ctx)
        case <-ctx.Done():
            return
        case <-e.stopCh:
            return
        }
    }
}

// Stop stops the evaluation loop
func (e *RuleEngine) Stop() {
    close(e.stopCh)
}

// evaluateAll evaluates all registered rules
func (e *RuleEngine) evaluateAll(ctx context.Context) {
    e.rulesMu.RLock()
    rules := make([]*AlertRule, 0, len(e.rules))
    for _, rule := range e.rules {
        rules = append(rules, rule)
    }
    e.rulesMu.RUnlock()

    for _, rule := range rules {
        if err := e.evaluateRule(ctx, rule); err != nil {
            // Log error but continue with other rules
            fmt.Printf("Error evaluating rule %s: %v\n", rule.ID, err)
        }
    }
}

// evaluateRule evaluates a single rule
func (e *RuleEngine) evaluateRule(ctx context.Context, rule *AlertRule) error {
    series, err := e.queryEngine.Query(ctx, rule.Query)
    if err != nil {
        return fmt.Errorf("query failed: %w", err)
    }

    rule.lastEval = time.Now()

    for _, ts := range series {
        fingerprint := e.calculateFingerprint(rule.ID, ts.Labels)

        e.instancesMu.Lock()
        instance, exists := e.instances[fingerprint]

        // Determine if condition is met
        conditionMet := e.evaluateCondition(ts.Value, rule)

        if conditionMet {
            if !exists || instance.State == AlertStateResolved {
                // New firing instance
                now := time.Now()
                instance = &AlertInstance{
                    RuleID:      rule.ID,
                    RuleName:    rule.Name,
                    Severity:    rule.Severity,
                    State:       AlertStatePending,
                    Labels:      e.mergeLabels(rule.Labels, ts.Labels),
                    Annotations: rule.Annotations,
                    Value:       ts.Value,
                    FiredAt:     now,
                    LastEval:    now,
                    Fingerprint: fingerprint,
                }
                e.instances[fingerprint] = instance

                // Check if duration has passed to transition to firing
                if now.Sub(instance.FiredAt) >= rule.Duration {
                    instance.State = AlertStateFiring
                    if err := e.notifier.Notify(ctx, instance); err != nil {
                        fmt.Printf("Failed to notify: %v\n", err)
                    }
                }
            } else {
                // Update existing instance
                instance.Value = ts.Value
                instance.LastEval = time.Now()

                // Check if should transition from pending to firing
                if instance.State == AlertStatePending {
                    if time.Since(instance.FiredAt) >= rule.Duration {
                        instance.State = AlertStateFiring
                        if err := e.notifier.Notify(ctx, instance); err != nil {
                            fmt.Printf("Failed to notify: %v\n", err)
                        }
                    }
                }
            }
        } else {
            // Condition not met - handle resolution
            if exists && instance.State != AlertStateResolved {
                if rule.EnableAutoResolve {
                    now := time.Now()
                    instance.ResolvedAt = &now
                    instance.State = AlertStateResolved
                    if err := e.notifier.Notify(ctx, instance); err != nil {
                        fmt.Printf("Failed to notify resolution: %v\n", err)
                    }
                }
            }
        }
        e.instancesMu.Unlock()
    }

    return nil
}

// evaluateCondition evaluates if the condition is met
func (e *RuleEngine) evaluateCondition(value float64, rule *AlertRule) bool {
    // Simple threshold check - can be extended
    return !math.IsNaN(value) && value > 0
}

// calculateFingerprint creates a unique fingerprint for an alert instance
func (e *RuleEngine) calculateFingerprint(ruleID string, labels map[string]string) string {
    // Simple implementation - in production use proper hashing
    return fmt.Sprintf("%s:%v", ruleID, labels)
}

// mergeLabels merges rule labels with series labels
func (e *RuleEngine) mergeLabels(ruleLabels, seriesLabels map[string]string) map[string]string {
    merged := make(map[string]string)
    for k, v := range ruleLabels {
        merged[k] = v
    }
    for k, v := range seriesLabels {
        merged[k] = v
    }
    return merged
}

// GetActiveAlerts returns all active (non-resolved) alerts
func (e *RuleEngine) GetActiveAlerts() []*AlertInstance {
    e.instancesMu.RLock()
    defer e.instancesMu.RUnlock()

    alerts := make([]*AlertInstance, 0)
    for _, instance := range e.instances {
        if instance.State != AlertStateResolved {
            alerts = append(alerts, instance)
        }
    }
    return alerts
}
```

### 2.2 Alert Grouping and Routing

```go
package alerting

import (
    "context"
    "fmt"
    "sort"
    "strings"
    "sync"
    "time"
)

// Router routes alerts to appropriate notification channels
type Router struct {
    routes      []Route
    routesMu    sync.RWMutex

    groupCache  map[string]*AlertGroup
    cacheMu     sync.RWMutex
    groupWindow time.Duration
}

// Route defines an alert routing rule
type Route struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Matchers    []Matcher         `json:"matchers"`
    GroupBy     []string          `json:"group_by"`
    GroupWait   time.Duration     `json:"group_wait"`
    GroupInterval time.Duration   `json:"group_interval"`
    RepeatInterval time.Duration  `json:"repeat_interval"`
    Continue    bool              `json:"continue"`
    Receiver    string            `json:"receiver"`
    Routes      []Route           `json:"routes,omitempty"` // Nested routes
}

// Matcher defines a label matching condition
type Matcher struct {
    Name    string `json:"name"`
    Value   string `json:"value"`
    IsRegex bool   `json:"is_regex"`
}

// Matches checks if labels match the matcher
func (m Matcher) Matches(labels map[string]string) bool {
    value, exists := labels[m.Name]
    if !exists {
        return false
    }

    if m.IsRegex {
        // In production, use proper regex matching
        return strings.Contains(value, m.Value)
    }
    return value == m.Value
}

// AlertGroup represents a grouped set of alerts
type AlertGroup struct {
    ID          string           `json:"id"`
    Labels      map[string]string `json:"labels"`
    Alerts      []*AlertInstance  `json:"alerts"`
    CreatedAt   time.Time        `json:"created_at"`
    UpdatedAt   time.Time        `json:"updated_at"`
    Receiver    string           `json:"receiver"`
    SentAt      *time.Time       `json:"sent_at,omitempty"`
}

// NewRouter creates a new alert router
func NewRouter(groupWindow time.Duration) *Router {
    if groupWindow == 0 {
        groupWindow = 5 * time.Minute
    }

    return &Router{
        routes:      make([]Route, 0),
        groupCache:  make(map[string]*AlertGroup),
        groupWindow: groupWindow,
    }
}

// AddRoute adds a routing rule
func (r *Router) AddRoute(route Route) {
    r.routesMu.Lock()
    defer r.routesMu.Unlock()
    r.routes = append(r.routes, route)
}

// RouteAlert routes a single alert to appropriate receivers
func (r *Router) RouteAlert(alert *AlertInstance) []string {
    r.routesMu.RLock()
    defer r.routesMu.RUnlock()

    receivers := make([]string, 0)

    for _, route := range r.routes {
        if r.matchesRoute(alert, route) {
            receivers = append(receivers, route.Receiver)

            // Add to group
            r.addToGroup(alert, route)

            if !route.Continue {
                break
            }
        }
    }

    return receivers
}

// matchesRoute checks if an alert matches a route
func (r *Router) matchesRoute(alert *AlertInstance, route Route) bool {
    for _, matcher := range route.Matchers {
        if !matcher.Matches(alert.Labels) {
            return false
        }
    }
    return true
}

// addToGroup adds an alert to its appropriate group
func (r *Router) addToGroup(alert *AlertInstance, route Route) {
    // Calculate group key based on GroupBy labels
    groupKey := r.calculateGroupKey(alert, route.GroupBy)
    groupID := fmt.Sprintf("%s:%s", route.ID, groupKey)

    r.cacheMu.Lock()
    defer r.cacheMu.Unlock()

    group, exists := r.groupCache[groupID]
    if !exists {
        now := time.Now()
        group = &AlertGroup{
            ID:        groupID,
            Labels:    r.extractGroupLabels(alert, route.GroupBy),
            Alerts:    make([]*AlertInstance, 0),
            CreatedAt: now,
            UpdatedAt: now,
            Receiver:  route.Receiver,
        }
        r.groupCache[groupID] = group
    }

    // Check if alert already exists in group
    found := false
    for i, a := range group.Alerts {
        if a.Fingerprint == alert.Fingerprint {
            group.Alerts[i] = alert
            found = true
            break
        }
    }

    if !found {
        group.Alerts = append(group.Alerts, alert)
    }

    group.UpdatedAt = time.Now()
}

// calculateGroupKey creates a group key from labels
func (r *Router) calculateGroupKey(alert *AlertInstance, groupBy []string) string {
    if len(groupBy) == 0 {
        return "default"
    }

    parts := make([]string, 0, len(groupBy))
    for _, label := range groupBy {
        if value, exists := alert.Labels[label]; exists {
            parts = append(parts, fmt.Sprintf("%s=%s", label, value))
        }
    }

    sort.Strings(parts)
    return strings.Join(parts, ",")
}

// extractGroupLabels extracts the grouping labels
func (r *Router) extractGroupLabels(alert *AlertInstance, groupBy []string) map[string]string {
    labels := make(map[string]string)
    for _, label := range groupBy {
        if value, exists := alert.Labels[label]; exists {
            labels[label] = value
        }
    }
    return labels
}

// GetGroups returns all active alert groups
func (r *Router) GetGroups() []*AlertGroup {
    r.cacheMu.RLock()
    defer r.cacheMu.RUnlock()

    groups := make([]*AlertGroup, 0, len(r.groupCache))
    for _, group := range r.groupCache {
        groups = append(groups, group)
    }

    // Sort by creation time
    sort.Slice(groups, func(i, j int) bool {
        return groups[i].CreatedAt.Before(groups[j].CreatedAt)
    })

    return groups
}

// Cleanup removes old groups
func (r *Router) Cleanup(maxAge time.Duration) {
    r.cacheMu.Lock()
    defer r.cacheMu.Unlock()

    cutoff := time.Now().Add(-maxAge)
    for id, group := range r.groupCache {
        if group.UpdatedAt.Before(cutoff) {
            delete(r.groupCache, id)
        }
    }
}
```

### 2.3 Alert Inhibition

```go
package alerting

import (
    "context"
    "sync"
    "time"
)

// Inhibitor suppresses alerts based on other firing alerts
type Inhibitor struct {
    rules      []InhibitRule
    rulesMu    sync.RWMutex

    activeFiring map[string]*AlertInstance
    firingMu     sync.RWMutex

    inhibitedAlerts map[string]time.Time
    inhibitedMu     sync.RWMutex
}

// InhibitRule defines when to inhibit an alert
type InhibitRule struct {
    ID string `json:"id"`

    // Source match - alert that causes inhibition
    SourceMatch      map[string]string `json:"source_match"`
    SourceMatchRegex map[string]string `json:"source_match_regex,omitempty"`

    // Target match - alert to be inhibited
    TargetMatch      map[string]string `json:"target_match"`
    TargetMatchRegex map[string]string `json:"target_match_regex,omitempty"`

    // Equal labels must match between source and target
    Equal []string `json:"equal"`

    // Duration for inhibition
    Duration time.Duration `json:"duration"`
}

// NewInhibitor creates a new alert inhibitor
func NewInhibitor() *Inhibitor {
    return &Inhibitor{
        rules:           make([]InhibitRule, 0),
        activeFiring:    make(map[string]*AlertInstance),
        inhibitedAlerts: make(map[string]time.Time),
    }
}

// AddRule adds an inhibition rule
func (i *Inhibitor) AddRule(rule InhibitRule) {
    i.rulesMu.Lock()
    defer i.rulesMu.Unlock()
    i.rules = append(i.rules, rule)
}

// RecordFiring records a firing alert
func (i *Inhibitor) RecordFiring(alert *AlertInstance) {
    i.firingMu.Lock()
    defer i.firingMu.Unlock()
    i.activeFiring[alert.Fingerprint] = alert
}

// RecordResolved records a resolved alert
func (i *Inhibitor) RecordResolved(fingerprint string) {
    i.firingMu.Lock()
    defer i.firingMu.Unlock()
    delete(i.activeFiring, fingerprint)
}

// IsInhibited checks if an alert should be inhibited
func (i *Inhibitor) IsInhibited(alert *AlertInstance) (bool, string) {
    i.rulesMu.RLock()
    defer i.rulesMu.RUnlock()

    i.firingMu.RLock()
    defer i.firingMu.RUnlock()

    for _, rule := range i.rules {
        // Check each active firing alert as potential source
        for _, sourceAlert := range i.activeFiring {
            if sourceAlert.Fingerprint == alert.Fingerprint {
                continue // Can't inhibit itself
            }

            if i.matchesRule(rule, sourceAlert, alert) {
                // Record inhibition
                i.inhibitedMu.Lock()
                i.inhibitedAlerts[alert.Fingerprint] = time.Now()
                i.inhibitedMu.Unlock()

                return true, rule.ID
            }
        }
    }

    return false, ""
}

// matchesRule checks if the source/target pair matches an inhibit rule
func (i *Inhibitor) matchesRule(rule InhibitRule, source, target *AlertInstance) bool {
    // Check source matches
    if !i.matchesLabels(source.Labels, rule.SourceMatch, rule.SourceMatchRegex) {
        return false
    }

    // Check target matches
    if !i.matchesLabels(target.Labels, rule.TargetMatch, rule.TargetMatchRegex) {
        return false
    }

    // Check equal labels match between source and target
    for _, label := range rule.Equal {
        sourceVal, sourceExists := source.Labels[label]
        targetVal, targetExists := target.Labels[label]

        if !sourceExists || !targetExists {
            return false
        }
        if sourceVal != targetVal {
            return false
        }
    }

    return true
}

// matchesLabels checks if labels match the matchers
func (i *Inhibitor) matchesLabels(labels, match, matchRegex map[string]string) bool {
    // Check exact matches
    for key, value := range match {
        if labels[key] != value {
            return false
        }
    }

    // Check regex matches (simplified - in production use proper regex)
    for key, pattern := range matchRegex {
        if !strings.Contains(labels[key], pattern) {
            return false
        }
    }

    return true
}

// GetInhibitedAlerts returns currently inhibited alerts
func (i *Inhibitor) GetInhibitedAlerts() map[string]time.Time {
    i.inhibitedMu.RLock()
    defer i.inhibitedMu.RUnlock()

    // Return a copy
    result := make(map[string]time.Time)
    for k, v := range i.inhibitedAlerts {
        result[k] = v
    }
    return result
}

// Cleanup removes old inhibition records
func (i *Inhibitor) Cleanup(maxAge time.Duration) {
    i.inhibitedMu.Lock()
    defer i.inhibitedMu.Unlock()

    cutoff := time.Now().Add(-maxAge)
    for fp, t := range i.inhibitedAlerts {
        if t.Before(cutoff) {
            delete(i.inhibitedAlerts, fp)
        }
    }
}

// Common inhibition rules
var CommonInhibitRules = []InhibitRule{
    {
        ID:          "node-down-inhibits-node-alerts",
        SourceMatch: map[string]string{"alertname": "NodeDown"},
        TargetMatch: map[string]string{},
        Equal:       []string{"instance"},
        Duration:    30 * time.Minute,
    },
    {
        ID:          "cluster-down-inhibits-cluster-alerts",
        SourceMatch: map[string]string{"alertname": "ClusterDown"},
        TargetMatch: map[string]string{},
        Equal:       []string{"cluster"},
        Duration:    1 * time.Hour,
    },
    {
        ID:          "high-latency-inhibits-timeout-alerts",
        SourceMatch: map[string]string{"severity": "critical", "alertname": "HighLatency"},
        TargetMatch: map[string]string{"alertname": "RequestTimeout"},
        Equal:       []string{"service"},
        Duration:    15 * time.Minute,
    },
}
```

---

## 3. Production-Ready Configurations

### 3.1 Prometheus Alert Rules

```yaml
# prometheus-alerts.yml
groups:
  - name: infrastructure
    interval: 15s
    rules:
      # Node alerts
      - alert: NodeCPUHigh
        expr: |
          (
            100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
          ) > 80
        for: 5m
        labels:
          severity: warning
          team: infrastructure
        annotations:
          summary: "High CPU usage on {{ $labels.instance }}"
          description: "CPU usage is above 80% (current: {{ $value }}%)"
          runbook_url: "https://wiki/runbooks/node-cpu-high"

      - alert: NodeMemoryHigh
        expr: |
          (
            node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100
          ) < 10
        for: 5m
        labels:
          severity: critical
          team: infrastructure
        annotations:
          summary: "Low memory on {{ $labels.instance }}"
          description: "Available memory is below 10% (current: {{ $value }}%)"
          runbook_url: "https://wiki/runbooks/node-memory-high"

      - alert: NodeDiskFull
        expr: |
          (
            node_filesystem_avail_bytes / node_filesystem_size_bytes * 100
          ) < 5
        for: 1m
        labels:
          severity: critical
          team: infrastructure
        annotations:
          summary: "Disk almost full on {{ $labels.instance }}"
          description: "Disk {{ $labels.device }} is {{ $value }}% full"

  - name: application
    interval: 15s
    rules:
      # Error rate alerts
      - alert: HighErrorRate
        expr: |
          (
            sum(rate(http_requests_total{status=~"5.."}[5m])) by (service)
            /
            sum(rate(http_requests_total[5m])) by (service)
          ) > 0.05
        for: 5m
        labels:
          severity: critical
          team: sre
        annotations:
          summary: "High error rate for {{ $labels.service }}"
          description: "Error rate is {{ $value | humanizePercentage }} for last 5m"

      - alert: ElevatedErrorRate
        expr: |
          (
            sum(rate(http_requests_total{status=~"5.."}[5m])) by (service)
            /
            sum(rate(http_requests_total[5m])) by (service)
          ) > 0.01
        for: 10m
        labels:
          severity: warning
          team: sre
        annotations:
          summary: "Elevated error rate for {{ $labels.service }}"
          description: "Error rate is {{ $value | humanizePercentage }} for last 10m"

      # Latency alerts
      - alert: HighLatencyP99
        expr: |
          histogram_quantile(0.99,
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service)
          ) > 2
        for: 10m
        labels:
          severity: warning
          team: sre
        annotations:
          summary: "High P99 latency for {{ $labels.service }}"
          description: "P99 latency is {{ $value }}s"

      # SLO-based alerts
      - alert: SLOErrorBudgetBurn
        expr: |
          (
            sum(rate(http_requests_total{status=~"5.."}[1h])) by (service)
            /
            sum(rate(http_requests_total[1h])) by (service)
          ) > (14.4 * 0.001)
        for: 2m
        labels:
          severity: critical
          team: sre
          slo: "true"
        annotations:
          summary: "Fast error budget burn for {{ $labels.service }}"
          description: "2% error budget burned in last hour"

      - alert: SLOErrorBudgetBurnSlow
        expr: |
          (
            sum(rate(http_requests_total{status=~"5.."}[6h])) by (service)
            /
            sum(rate(http_requests_total[6h])) by (service)
          ) > (6 * 0.001)
        for: 15m
        labels:
          severity: warning
          team: sre
          slo: "true"
        annotations:
          summary: "Slow error budget burn for {{ $labels.service }}"
          description: "5% error budget burned in last 6 hours"

  - name: business
    interval: 60s
    rules:
      - alert: LowConversionRate
        expr: |
          (
            sum(rate(checkouts_total[1h]))
            /
            sum(rate(cart_adds_total[1h]))
          ) < 0.05
        for: 30m
        labels:
          severity: high
          team: product
        annotations:
          summary: "Low conversion rate detected"
          description: "Conversion rate dropped to {{ $value | humanizePercentage }}"
```

### 3.2 Alertmanager Configuration

```yaml
# alertmanager.yml
global:
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alerts@example.com'
  smtp_auth_username: 'alerts@example.com'
  smtp_auth_password: '<password>'
  slack_api_url: '<slack_webhook_url>'
  pagerduty_url: 'https://events.pagerduty.com/v2/enqueue'
  resolve_timeout: 5m

# Templates for notifications
templates:
  - '/etc/alertmanager/templates/*.tmpl'

# Inhibition rules
inhibit_rules:
  # Inhibit warning alerts if critical for same service
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['service', 'instance']

  # Inhibit node-level alerts if node is down
  - source_match:
      alertname: 'NodeDown'
    target_match_re:
      alertname: 'Node.+|System.+'
    equal: ['instance']

  # Inhibit all alerts if entire cluster is down
  - source_match:
      alertname: 'ClusterDown'
    target_match: {}
    equal: ['cluster']

# Routing tree
route:
  # Root receiver
  receiver: 'default'

  # How long to initially wait to send a notification
  group_wait: 30s

  # How long to wait before sending a notification about new alerts
  group_interval: 5m

  # How long to wait before sending a notification again
  repeat_interval: 4h

  # Default grouping
  group_by: ['alertname', 'severity', 'service']

  # Continue processing if matched (for multiple routes)
  continue: false

  # Sub-routes
  routes:
    # Critical alerts -> PagerDuty + Slack
    - match:
        severity: critical
      receiver: pagerduty-critical
      continue: true
      group_by: ['service']
      group_wait: 0s
      repeat_interval: 1h

    # Infrastructure team alerts
    - match_re:
        team: infrastructure|sre
      receiver: slack-infrastructure
      routes:
        - match:
            severity: critical
          receiver: pagerduty-infrastructure
          continue: true

    # Application team alerts
    - match_re:
        team: application|backend
      receiver: slack-application
      group_by: ['service', 'severity']

    # SLO-based alerts
    - match:
        slo: "true"
      receiver: slack-slo
      routes:
        - match:
            severity: critical
          receiver: pagerduty-slo
          continue: true

    # Business metrics -> Product team
    - match:
        team: product
      receiver: email-product
      group_by: ['alertname']
      group_wait: 5m
      repeat_interval: 24h

    # Informational alerts -> Logging only
    - match:
        severity: info
      receiver: null

# Receivers
receivers:
  - name: 'default'
    slack_configs:
      - channel: '#alerts-default'
        send_resolved: true
        title: '{{ template "slack.default.title" . }}'
        text: '{{ template "slack.default.text" . }}'

  - name: 'null'

  - name: 'pagerduty-critical'
    pagerduty_configs:
      - service_key: '<pagerduty-service-key>'
        severity: critical
        description: '{{ .CommonAnnotations.summary }}'
        details:
          firing: '{{ template "pagerduty.default.instances" .Alerts.Firing }}'
          resolved: '{{ template "pagerduty.default.instances" .Alerts.Resolved }}'
          runbook: '{{ .CommonAnnotations.runbook_url }}'

  - name: 'pagerduty-infrastructure'
    pagerduty_configs:
      - service_key: '<infrastructure-service-key>'
        severity: '{{ .CommonLabels.severity }}'
        description: '{{ .CommonAnnotations.summary }}'

  - name: 'pagerduty-slo'
    pagerduty_configs:
      - service_key: '<slo-service-key>'
        severity: warning
        description: 'SLO Alert: {{ .CommonAnnotations.summary }}'

  - name: 'slack-infrastructure'
    slack_configs:
      - channel: '#infrastructure-alerts'
        send_resolved: true
        color: '{{ if eq .CommonLabels.severity "critical" }}danger{{ else }}warning{{ end }}'
        title: 'Infrastructure Alert'
        text: |
          {{ range .Alerts }}
          *Alert:* {{ .Annotations.summary }}
          *Severity:* {{ .Labels.severity }}
          *Instance:* {{ .Labels.instance }}
          *Description:* {{ .Annotations.description }}
          {{ if .Annotations.runbook_url }}*Runbook:* {{ .Annotations.runbook_url }}{{ end }}
          {{ end }}

  - name: 'slack-application'
    slack_configs:
      - channel: '#app-alerts'
        send_resolved: true
        color: '{{ if eq .CommonLabels.severity "critical" }}danger{{ else if eq .CommonLabels.severity "warning" }}warning{{ else }}good{{ end }}'
        title: 'Application Alert - {{ .CommonLabels.service }}'
        text: |
          {{ range .Alerts }}
          *Service:* {{ .Labels.service }}
          *Alert:* {{ .Annotations.summary }}
          *Details:* {{ .Annotations.description }}
          {{ end }}

  - name: 'slack-slo'
    slack_configs:
      - channel: '#slo-alerts'
        send_resolved: true
        color: 'warning'
        title: 'SLO Alert - {{ .CommonLabels.service }}'
        text: |
          SLO Error Budget Burn Detected
          {{ range .Alerts }}
          Service: {{ .Labels.service }}
          {{ .Annotations.description }}
          {{ end }}

  - name: 'email-product'
    email_configs:
      - to: 'product-team@example.com'
        send_resolved: true
        headers:
          Subject: 'Business Alert: {{ .CommonAnnotations.summary }}'
        html: '{{ template "email.default.html" . }}'
```

---

## 4. Security Considerations

### 4.1 Alert Security Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Alerting Security Checklist                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ACCESS CONTROL                                                             │
│  [ ] Alertmanager UI behind authentication                                  │
│  [ ] RBAC for alert rule modification                                       │
│  [ ] Audit logging for rule changes                                         │
│  [ ] Separate credentials for different environments                        │
│  [ ] Regular rotation of API keys/tokens                                    │
│                                                                             │
│  DATA PROTECTION                                                            │
│  [ ] No PII in alert labels/annotations                                     │
│  [ ] Sanitize error messages in alerts                                      │
│  [ ] Encrypt alert notification channels                                    │
│  [ ] Secure storage of notification credentials                             │
│  [ ] TLS for all webhook endpoints                                          │
│                                                                             │
│  NOTIFICATION SECURITY                                                      │
│  [ ] Verify webhook signatures                                              │
│  [ ] Rate limiting on notification endpoints                                │
│  [ ] IP allowlisting for webhook sources                                    │
│  [ ] Validate notification payload schemas                                  │
│  [ ] DLP scanning for outbound notifications                                │
│                                                                             │
│  OPERATIONAL SECURITY                                                       │
│  [ ] Alert on authentication failures                                       │
│  [ ] Alert on privilege escalation                                          │
│  [ ] Alert on configuration changes                                         │
│  [ ] Alert on unusual query patterns                                        │
│  [ ] Separate alerting for security events                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Secure Webhook Implementation

```go
package alerting

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// SecureWebhookNotifier sends secure webhook notifications
type SecureWebhookNotifier struct {
    endpoint   string
    secretKey  string
    httpClient *http.Client
    timeout    time.Duration
}

// WebhookPayload represents the webhook payload
type WebhookPayload struct {
    Version   string             `json:"version"`
    Timestamp time.Time          `json:"timestamp"`
    Alerts    []*AlertInstance   `json:"alerts"`
    GroupKey  string             `json:"group_key"`
    Status    string             `json:"status"` // firing, resolved
    Signature string             `json:"signature,omitempty"`
}

// NewSecureWebhookNotifier creates a secure webhook notifier
func NewSecureWebhookNotifier(endpoint, secretKey string, timeout time.Duration) *SecureWebhookNotifier {
    if timeout == 0 {
        timeout = 30 * time.Second
    }

    return &SecureWebhookNotifier{
        endpoint:  endpoint,
        secretKey: secretKey,
        httpClient: &http.Client{
            Timeout: timeout,
        },
        timeout: timeout,
    }
}

// Notify sends a secure webhook notification
func (n *SecureWebhookNotifier) Notify(ctx context.Context, alert *AlertInstance) error {
    payload := WebhookPayload{
        Version:   "1.0",
        Timestamp: time.Now().UTC(),
        Alerts:    []*AlertInstance{alert},
        Status:    string(alert.State),
    }

    return n.sendWebhook(ctx, payload)
}

// NotifyGroup sends a grouped webhook notification
func (n *SecureWebhookNotifier) NotifyGroup(ctx context.Context, group *AlertGroup) error {
    status := "firing"
    for _, alert := range group.Alerts {
        if alert.State == AlertStateResolved {
            status = "resolved"
            break
        }
    }

    payload := WebhookPayload{
        Version:   "1.0",
        Timestamp: time.Now().UTC(),
        Alerts:    group.Alerts,
        GroupKey:  group.ID,
        Status:    status,
    }

    return n.sendWebhook(ctx, payload)
}

// sendWebhook sends the webhook with HMAC signature
func (n *SecureWebhookNotifier) sendWebhook(ctx context.Context, payload WebhookPayload) error {
    // Marshal payload
    body, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    // Generate HMAC signature
    signature := n.generateSignature(body)

    // Create request
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.endpoint, bytes.NewReader(body))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Webhook-Signature", "sha256="+signature)
    req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", payload.Timestamp.Unix()))
    req.Header.Set("X-Webhook-Version", payload.Version)

    // Send request
    resp, err := n.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("webhook returned status %d", resp.StatusCode)
    }

    return nil
}

// generateSignature generates HMAC-SHA256 signature
func (n *SecureWebhookNotifier) generateSignature(body []byte) string {
    mac := hmac.New(sha256.New, []byte(n.secretKey))
    mac.Write(body)
    return hex.EncodeToString(mac.Sum(nil))
}

// WebhookReceiver handles incoming webhook notifications securely
type WebhookReceiver struct {
    secretKey string
    handler   func(payload *WebhookPayload) error
}

// NewWebhookReceiver creates a new secure webhook receiver
func NewWebhookReceiver(secretKey string, handler func(payload *WebhookPayload) error) *WebhookReceiver {
    return &WebhookReceiver{
        secretKey: secretKey,
        handler:   handler,
    }
}

// HTTPHandler returns an HTTP handler for receiving webhooks
func (r *WebhookReceiver) HTTPHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        // Verify method
        if req.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Read body
        body, err := io.ReadAll(req.Body)
        if err != nil {
            http.Error(w, "Failed to read body", http.StatusBadRequest)
            return
        }
        defer req.Body.Close()

        // Verify signature
        signature := req.Header.Get("X-Webhook-Signature")
        if signature == "" {
            http.Error(w, "Missing signature", http.StatusUnauthorized)
            return
        }

        // Remove "sha256=" prefix if present
        signature = strings.TrimPrefix(signature, "sha256=")

        expectedSignature := r.generateSignature(body)
        if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
            http.Error(w, "Invalid signature", http.StatusUnauthorized)
            return
        }

        // Verify timestamp (prevent replay attacks)
        timestamp := req.Header.Get("X-Webhook-Timestamp")
        if timestamp != "" {
            ts, err := strconv.ParseInt(timestamp, 10, 64)
            if err == nil {
                // Reject if older than 5 minutes
                if time.Since(time.Unix(ts, 0)) > 5*time.Minute {
                    http.Error(w, "Request too old", http.StatusUnauthorized)
                    return
                }
            }
        }

        // Parse payload
        var payload WebhookPayload
        if err := json.Unmarshal(body, &payload); err != nil {
            http.Error(w, "Invalid payload", http.StatusBadRequest)
            return
        }

        // Handle payload
        if err := r.handler(&payload); err != nil {
            http.Error(w, "Handler error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    }
}

// generateSignature generates HMAC-SHA256 signature
func (r *WebhookReceiver) generateSignature(body []byte) string {
    mac := hmac.New(sha256.New, []byte(r.secretKey))
    mac.Write(body)
    return hex.EncodeToString(mac.Sum(nil))
}
```

---

## 5. Compliance Requirements

### 5.1 Audit and Compliance

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Alerting Compliance Requirements                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SOC 2 TYPE II                                                              │
│  ├─ CC6.1: Logical access security - Alert on unauthorized access attempts  │
│  ├─ CC6.2: Prior to access - Alert on new account creation                  │
│  ├─ CC6.3: Access removal - Alert on access modifications                   │
│  ├─ CC7.2: System monitoring - Comprehensive alerting coverage              │
│  └─ CC7.3: Incident detection - Automated incident creation                 │
│                                                                             │
│  ISO 27001                                                                  │
│  ├─ A.12.4: Logging and monitoring - Alert on security events               │
│  ├─ A.16.1: Incident management - Incident response alerts                  │
│  ├─ A.12.6: Technical vulnerability management - Vulnerability alerts       │
│  └─ A.9.4: System access control - Access violation alerts                  │
│                                                                             │
│  PCI DSS                                                                    │
│  ├─ Req 10.4: Synchronize clocks - Alert on clock skew                      │
│  ├─ Req 10.5: Secure audit trails - Tampering alerts                        │
│  ├─ Req 10.6: Review logs - Anomaly detection alerts                        │
│  ├─ Req 10.7: Retain audit history - Retention compliance alerts            │
│  └─ Req 11.4: IDS/IPS - Security incident alerts                            │
│                                                                             │
│  HIPAA                                                                      │
│  ├─ §164.312(b): Audit controls - PHI access alerts                         │
│  ├─ §164.312(c)(1): Integrity - Data integrity alerts                       │
│  ├─ §164.312(c)(2): Mechanism - Encryption alerts                           │
│  └─ §164.312(e)(1): Transmission security - Network security alerts         │
│                                                                             │
│  GDPR                                                                       │
│  ├─ Art. 33: Breach notification - 72-hour breach detection                 │
│  ├─ Art. 35: DPIA triggers - High-risk processing alerts                    │
│  └─ Art. 5: Principles - Data retention and accuracy alerts                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Decision Matrices

### 6.1 Alert Severity Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Alert Severity Decision Matrix                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Impact        │  Scope              │  Urgency     │  Severity   │  SLA    │
├────────────────┼─────────────────────┼──────────────┼─────────────┼─────────│
│  Complete      │  All users          │  Immediate   │  P0         │  5 min  │
│  outage        │  Multiple regions   │              │  Critical   │         │
│                │                     │              │             │         │
│  Major         │  Many users         │  < 15 min    │  P1         │  15 min │
│  degradation   │  Single region      │              │  High       │         │
│                │                     │              │             │         │
│  Minor         │  Some users         │  < 1 hour    │  P2         │  1 hour │
│  impact        │  Single service     │              │  Medium     │         │
│                │                     │              │             │         │
│  No immediate  │  Monitoring only    │  < 24 hours  │  P3         │  24 hrs │
│  impact        │  Single instance    │              │  Low        │         │
│                │                     │              │             │         │
│  Informational │  N/A                │  Next sprint │  P4         │  None   │
│                │                     │              │  Info       │         │
│                                                                             │
│  Decision Flow:                                                            │
│                                                                             │
│  Error Rate > 5% ──┐                                                        │
│  Revenue Impact    ├──► P0 (Critical) ──► Page immediately                │
│  Security Breach ──┘                                                        │
│                                                                             │
│  Error Rate 1-5% ──┐                                                        │
│  Latency > 2x      ├──► P1 (High) ──────► Page if sustained               │
│  Capacity > 90% ───┘                                                        │
│                                                                             │
│  Error Rate 0.1-1%─┐                                                        │
│  Latency > 1.5x    ├──► P2 (Medium) ────► Slack alert                     │
│  Capacity > 80% ───┘                                                        │
│                                                                             │
│  Error Rate < 0.1%─┐                                                        │
│  Latency > 1.2x    ├──► P3 (Low) ───────► Dashboard only                  │
│  Capacity > 70% ───┘                                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Notification Channel Selection Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                  Notification Channel Selection Matrix                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Severity  │  PagerDuty  │  SMS    │  Phone  │  Slack  │  Email  │  Ticket │
├────────────┼─────────────┼─────────┼─────────┼─────────┼─────────┼─────────│
│  P0        │     ✓       │    ✓    │    ✓    │    ✓    │    ✗    │    ✗    │
│  Critical  │   Page      │  Backup │  Escal  │  Follow │         │         │
├────────────┼─────────────┼─────────┼─────────┼─────────┼─────────┼─────────│
│  P1        │     ✓       │    ✗    │    ✗    │    ✓    │    ✗    │    ✗    │
│  High      │   Page      │         │         │  Follow │         │         │
├────────────┼─────────────┼─────────┼─────────┼─────────┼─────────┼─────────│
│  P2        │     ✗       │    ✗    │    ✗    │    ✓    │    ✗    │    ✓    │
│  Medium    │             │         │         │  Main   │         │  Track  │
├────────────┼─────────────┼─────────┼─────────┼─────────┼─────────┼─────────│
│  P3        │     ✗       │    ✗    │    ✗    │    ✗    │    ✓    │    ✓    │
│  Low       │             │         │         │         │  Main   │  Track  │
├────────────┼─────────────┼─────────┼─────────┼─────────┼─────────┼─────────│
│  P4        │     ✗       │    ✗    │    ✗    │    ✗    │    ✗    │    ✗    │
│  Info      │             │         │         │         │         │  Log    │
│                                                                             │
│  Notes:                                                                    │
│  • PagerDuty: Primary on-call paging system                                │
│  • SMS: Fallback for critical alerts only                                  │
│  • Phone: Auto-escalation after 10 min unack                               │
│  • Slack: Team coordination and updates                                    │
│  • Email: Non-urgent notifications and summaries                           │
│  • Ticket: JIRA/ServiceNow for tracking and post-mortems                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Alert Maintenance Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Alert Maintenance Decision Matrix                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Problem                │  Metric           │  Threshold │  Action          │
├─────────────────────────┼───────────────────┼────────────┼──────────────────│
│  Too many alerts        │  Alerts/week      │  > 50      │  Review & tune   │
│  (Alert fatigue)        │  per person       │            │  rules           │
│  ───────────────────────┼───────────────────┼────────────┼──────────────────│
│  False positive rate    │  FPs / Total      │  > 20%     │  Adjust          │
│  too high               │  alerts           │            │  thresholds      │
│  ───────────────────────┼───────────────────┼────────────┼──────────────────│
│  Alert latency          │  Time to detect   │  > 5 min   │  Reduce          │
│  too slow               │                   │            │  evaluation      │
│                         │                   │            │  interval        │
│  ───────────────────────┼───────────────────┼────────────┼──────────────────│
│  Missing critical       │  MTTD vs SLA      │  > SLA     │  Add coverage    │
│  incidents              │                   │            │  gaps            │
│  ───────────────────────┼───────────────────┼────────────┼──────────────────│
│  Alert storms           │  Alerts/min       │  > 10      │  Improve         │
│  during incidents       │  during incident  │            │  grouping        │
│  ───────────────────────┼───────────────────┼────────────┼──────────────────│
│  Stale alerts           │  Last modified    │  > 90 days │  Review and      │
│                         │                   │            │  update or       │
│                         │                   │            │  remove          │
│  ───────────────────────┼───────────────────┼────────────┼──────────────────│
│  Low actionability      │  Alerts without   │  > 30%     │  Add runbooks    │
│                         │  runbooks         │            │  or remove       │
│                                                                             │
│  Maintenance Schedule:                                                     │
│  • Daily: Review unacknowledged alerts                                      │
│  • Weekly: Alert quality review meeting                                     │
│  • Monthly: Tune noisy alerts                                               │
│  • Quarterly: Comprehensive alert audit                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Alerting Best Practices Summary                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DO                                                                         │
│  ✓ Alert on symptoms (user impact) rather than causes                      │
│  ✓ Include actionable information in alert descriptions                     │
│  ✓ Use SLO-based alerting for user-facing services                          │
│  ✓ Group related alerts to reduce noise                                     │
│  ✓ Provide runbook links in every alert                                     │
│  ✓ Test alert rules before deploying to production                          │
│  ✓ Set appropriate alert durations ("for" clause)                           │
│  ✓ Use severity labels consistently                                         │
│  ✓ Implement alert inhibition for known dependencies                        │
│  ✓ Monitor the monitoring system (meta-monitoring)                          │
│  ✓ Regular review of alert effectiveness                                    │
│  ✓ Include context (links to dashboards, logs)                              │
│                                                                             │
│  DON'T                                                                      │
│  ✗ Alert on every anomaly or spike                                          │
│  ✗ Send alerts to individuals instead of on-call rotation                   │
│  ✗ Create alerts without runbooks                                           │
│  ✗ Use same severity for all alerts                                         │
│  ✗ Alert on metrics that naturally fluctuate                                │
│  ✗ Create alerts that page during maintenance windows                       │
│  ✗ Forget to handle timezone differences                                    │
│  ✗ Ignore alert fatigue signals                                             │
│  ✗ Leave alerts firing without investigation                                │
│  ✗ Create alerts that can't be acted upon                                   │
│                                                                             │
│  KEY METRICS TO TRACK                                                       │
│  • MTTD (Mean Time To Detect)                                               │
│  • MTTR (Mean Time To Respond)                                              │
│  • False positive rate                                                      │
│  • Alert coverage (% of incidents detected)                                 │
│  • Alert fatigue index (alerts per on-call shift)                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Google SRE Book - Monitoring and Alerting
2. Prometheus Alerting Best Practices
3. PagerDuty Incident Response Guide
4. Site Reliability Workbook - Alerting on SLOs
5. My Philosophy on Alerting - Rob Ewaschuk
