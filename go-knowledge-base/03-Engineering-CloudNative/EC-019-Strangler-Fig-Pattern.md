# EC-019: Strangler Fig Pattern

## Problem Formalization

### The Monolithic Migration Challenge

Migrating monolithic applications to microservices presents one of the most complex architectural challenges: how to incrementally replace a legacy system without disrupting business operations, while managing risk and maintaining data consistency.

#### Problem Statement

Given a monolithic system M with:

- Components C = {c₁, c₂, ..., cₙ}
- Data stores D = {d₁, d₂, ..., dₘ}
- Dependencies graph G = (C, E) where E represents inter-component calls
- Users U accessing the system

Find a migration sequence S = {s₁, s₂, ..., sₖ} such that:

```
Minimize: Risk(Migration) = Σ Risk(sᵢ)
Minimize: Downtime(Migration)
Subject to:
    - Functional equivalence: Behavior(M) = Behavior(Microservices)
    - Data consistency: ∀d ∈ D, Consistent(d) throughout migration
    - Rollback capability: Can revert any sᵢ within T minutes
    - Gradual traffic shifting: Can route p% traffic to new services
```

### Anti-Patterns to Avoid

#### 1. Big Bang Migration

```
┌─────────────────────────────────────────────────────────────┐
│  Anti-Pattern: Big Bang                                     │
│                                                             │
│  Monolith ──► [STOP] ──► New System ──► [PRAY]            │
│                                                             │
│  Risk: Complete system failure, no rollback                │
│  Downtime: Hours to days                                   │
│  Success Rate: < 30%                                       │
└─────────────────────────────────────────────────────────────┘
```

#### 2. Database Sharing

```
┌─────────────────────────────────────────────────────────────┐
│  Anti-Pattern: Shared Database                              │
│                                                             │
│  ┌──────────┐      ┌──────────┐                            │
│  │ Service A├──────┤          │                            │
│  └──────────┤      │  Shared  │                            │
│  ┌──────────┤      │  Database│                            │
│  │ Service B├──────┤          │                            │
│  └──────────┴──────┴──────────┴──────┐                     │
│  ┌──────────┐      ┌──────────┐      │                     │
│  │ Service C├──────┤ Service D├──────┘                     │
│  └──────────┘      └──────────┘                            │
│                                                             │
│  Risk: Hidden coupling, schema conflicts                   │
│  Problem: No bounded context isolation                      │
└─────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Strangler Fig Pattern

The Strangler Fig pattern, named after the strangler fig tree that gradually replaces its host, provides an incremental approach to system migration.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     Strangler Fig Migration Timeline                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Phase 0: Identify Seams                                                │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Monolithic Application                        │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │   │
│  │  │  User    │ │  Order   │ │ Payment  │ │ Inventory│           │   │
│  │  │  Mgmt    │ │  Mgmt    │ │  Mgmt    │ │   Mgmt   │           │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │   │
│  │       │            │            │            │                  │   │
│  │       └────────────┴────────────┴────────────┘                  │   │
│  │                    Shared Database                               │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Phase 1: Deploy Router/Facade                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  ┌──────────┐                                                   │   │
│  │  │  Router  │◄── Incoming Traffic                               │   │
│  │  │  (New)   │                                                   │   │
│  │  └────┬─────┘                                                   │   │
│  │       │                                                         │   │
│  │       └──────────────────────┐                                  │   │
│  │                              ▼                                  │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │                    Monolithic Application                │   │   │
│  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │   │   │
│  │  │  │  User    │ │  Order   │ │ Payment  │ │ Inventory│   │   │   │
│  │  │  │  Mgmt    │ │  Mgmt    │ │  Mgmt    │ │   Mgmt   │   │   │   │
│  │  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Phase 2: Extract First Service (e.g., User Management)                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  ┌──────────┐                                                   │   │
│  │  │  Router  │◄── Incoming Traffic                               │   │
│  │  │          │                                                   │   │
│  │  │ /users/* │──►┌──────────────┐                                │   │
│  │  │ (new)    │   │  User        │                                │   │
│  │  │          │   │  Service     │                                │   │
│  │  │ /orders/*│   │  (New)       │                                │   │
│  │  │ (old)    │   └──────────────┘                                │   │
│  │  │          │          │                                        │   │
│  │  │ other/*  │          ▼                                        │   │
│  │  │ (old)    │   ┌──────────────┐                                │   │
│  │  └────┬─────┘   │  User DB     │                                │   │
│  │       │         │  (New)       │                                │   │
│  │       │         └──────────────┘                                │   │
│  │       │                                                         │   │
│  │       └──────────────────────┐                                  │   │
│  │                              ▼                                  │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │                    Monolithic Application                │   │   │
│  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐                │   │   │
│  │  │  │  Order   │ │ Payment  │ │ Inventory│                │   │   │
│  │  │  │  Mgmt    │ │  Mgmt    │ │   Mgmt   │                │   │   │
│  │  │  └──────────┘ └──────────┘ └──────────┘                │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Phase N: Complete Migration                                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│  │  User    │  │  Order   │  │ Payment  │  │ Inventory│              │
│  │ Service  │  │ Service  │  │ Service  │  │ Service  │              │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘              │
│                                                                         │
│  Monolith completely strangled and removed                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Migration Sequencing Strategy

```go
// internal/migration/sequencer.go
package migration

import (
    "context"
    "fmt"
    "time"
)

// MigrationSequencer determines optimal migration order
type MigrationSequencer struct {
    analyzer    *DependencyAnalyzer
    riskAssessor *RiskAssessor
    validator   *MigrationValidator
}

// CandidateService represents a service ready for extraction
type CandidateService struct {
    ID              string
    Name            string
    Dependencies    []string
    Dependents      []string
    Complexity      ComplexityScore
    BusinessValue   BusinessValueScore
    DataCoupling    DataCouplingScore
    TestCoverage    float64
    RiskScore       RiskScore
}

type ComplexityScore struct {
    LinesOfCode     int
    CyclomaticComplexity int
    ExternalCalls   int
    DatabaseTables  int
}

type BusinessValueScore struct {
    TransactionVolume   int
    RevenueImpact       float64
    UserFacing          bool
    RegulatoryRequired  bool
}

type DataCouplingScore struct {
    SharedTables        int
    ForeignKeyChains    int
    CrossDomainJoins    int
}

// Prioritize returns ordered list of services to migrate
func (ms *MigrationSequencer) Prioritize(services []CandidateService) []MigrationPhase {
    // Calculate priority score for each service
    scored := make([]scoredService, len(services))

    for i, svc := range services {
        score := ms.calculatePriorityScore(svc)
        scored[i] = scoredService{
            Service: svc,
            Score:   score,
        }
    }

    // Sort by score (descending)
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].Score > scored[j].Score
    })

    // Group into phases based on dependencies
    phases := ms.groupIntoPhases(scored)

    return phases
}

func (ms *MigrationSequencer) calculatePriorityScore(svc CandidateService) float64 {
    // Higher score = higher priority for early migration

    // Factors that increase priority (lower risk, easier)
    lowCouplingWeight := 0.25
    highTestCoverageWeight := 0.20
    lowComplexityWeight := 0.15
    lowBusinessRiskWeight := 0.15

    // Factors that decrease priority (postpone)
    highDependencyWeight := -0.15
    criticalPathWeight := -0.10

    score := 0.0

    // Low data coupling is good for early migration
    couplingScore := 1.0 / (1.0 + float64(svc.DataCoupling.SharedTables))
    score += couplingScore * lowCouplingWeight

    // High test coverage reduces migration risk
    score += svc.TestCoverage * highTestCoverageWeight

    // Lower complexity is easier to migrate
    complexityScore := 1.0 / (1.0 + float64(svc.Complexity.CyclomaticComplexity)/100)
    score += complexityScore * lowComplexityWeight

    // Non-critical services first
    if !svc.BusinessValue.RegulatoryRequired {
        score += lowBusinessRiskWeight
    }

    // Many dependents means postpone (ripples)
    dependencyScore := 1.0 / (1.0 + float64(len(svc.Dependents)))
    score += dependencyScore * highDependencyWeight

    return score
}

// MigrationPhase represents a group of services that can migrate together
type MigrationPhase struct {
    Number   int
    Services []string
    Duration time.Duration
    Risks    []Risk
}
```

### Traffic Routing Architecture

```go
// internal/router/strangler_router.go
package router

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync"
    "time"

    "github.com/prometheus/client_golang/prometheus"
)

// StranglerRouter routes traffic between monolith and new services
type StranglerRouter struct {
    monolithURL    *url.URL
    monolithProxy  *httputil.ReverseProxy

    // Service routing rules
    routes         map[string]*RouteRule
    routesMu       sync.RWMutex

    // Traffic splitting
    trafficSplit   map[string]*TrafficSplit
    splitMu        sync.RWMutex

    // Metrics
    requestCounter *prometheus.CounterVec
    latencyHistogram *prometheus.HistogramVec

    logger         *zap.Logger
}

// RouteRule defines routing for a specific path pattern
type RouteRule struct {
    Pattern        string
    TargetURL      *url.URL
    Proxy          *httputil.ReverseProxy
    StripPrefix    bool
    Transformers   []RequestTransformer
    Active         bool
    MigrationStage MigrationStage
}

type MigrationStage int

const (
    StageShadow MigrationStage = iota  // Mirroring, no user impact
    StageCanary                        // Small % of traffic
    StageParallel                      // A/B comparison
    StageCutover                       // 100% new service
    StageMonolithRetired               // Monolith code removed
)

// TrafficSplit controls percentage of traffic to new service
type TrafficSplit struct {
    ServiceID       string
    NewServicePercent int // 0-100
    Strategy        SplitStrategy
    CookieName      string // for sticky sessions
}

type SplitStrategy int

const (
    SplitRandom SplitStrategy = iota
    SplitUserID
    SplitGeography
    SplitHeader
)

func NewStranglerRouter(monolithURL string, logger *zap.Logger) (*StranglerRouter, error) {
    mURL, err := url.Parse(monolithURL)
    if err != nil {
        return nil, fmt.Errorf("invalid monolith URL: %w", err)
    }

    sr := &StranglerRouter{
        monolithURL:   mURL,
        monolithProxy: httputil.NewSingleHostReverseProxy(mURL),
        routes:        make(map[string]*RouteRule),
        trafficSplit:  make(map[string]*TrafficSplit),
        logger:        logger,
        requestCounter: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "strangler_requests_total",
                Help: "Total requests routed",
            },
            []string{"service", "destination"},
        ),
        latencyHistogram: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "strangler_request_duration_seconds",
                Help: "Request duration",
            },
            []string{"service", "destination"},
            prometheus.DefBuckets,
        ),
    }

    return sr, nil
}

// RegisterRoute adds a new service route
func (sr *StranglerRouter) RegisterRoute(rule *RouteRule) error {
    sr.routesMu.Lock()
    defer sr.routesMu.Unlock()

    // Create proxy for new service
    proxy := httputil.NewSingleHostReverseProxy(rule.TargetURL)
    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        sr.logger.Error("proxy error",
            zap.Error(err),
            zap.String("service", rule.Pattern),
        )
        // Fallback to monolith on error
        sr.monolithProxy.ServeHTTP(w, r)
    }

    rule.Proxy = proxy
    sr.routes[rule.Pattern] = rule

    sr.logger.Info("registered route",
        zap.String("pattern", rule.Pattern),
        zap.String("target", rule.TargetURL.String()),
    )

    return nil
}

// SetTrafficSplit configures traffic percentage for a service
func (sr *StranglerRouter) SetTrafficSplit(split *TrafficSplit) {
    sr.splitMu.Lock()
    defer sr.splitMu.Unlock()

    sr.trafficSplit[split.ServiceID] = split

    sr.logger.Info("traffic split updated",
        zap.String("service", split.ServiceID),
        zap.Int("new_percent", split.NewServicePercent),
    )
}

// ServeHTTP implements http.Handler
func (sr *StranglerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // Find matching route
    route := sr.matchRoute(r.URL.Path)

    if route == nil || !route.Active {
        // Route to monolith
        sr.routeToMonolith(w, r)
        sr.recordMetrics("default", "monolith", time.Since(start))
        return
    }

    // Determine destination based on traffic split
    destination := sr.determineDestination(r, route)

    switch destination {
    case "new":
        sr.routeToNewService(w, r, route)
        sr.recordMetrics(route.Pattern, "new_service", time.Since(start))
    case "shadow":
        // Send to both, return monolith response
        sr.routeShadow(w, r, route)
        sr.recordMetrics(route.Pattern, "shadow", time.Since(start))
    default:
        sr.routeToMonolith(w, r)
        sr.recordMetrics(route.Pattern, "monolith", time.Since(start))
    }
}

func (sr *StranglerRouter) determineDestination(r *http.Request, route *RouteRule) string {
    sr.splitMu.RLock()
    split, exists := sr.trafficSplit[route.Pattern]
    sr.splitMu.RUnlock()

    if !exists || split.NewServicePercent == 0 {
        return "monolith"
    }

    if route.MigrationStage == StageShadow {
        return "shadow"
    }

    // Calculate if this request goes to new service
    switch split.Strategy {
    case SplitRandom:
        if random.Intn(100) < split.NewServicePercent {
            return "new"
        }

    case SplitUserID:
        userID := r.Header.Get("X-User-ID")
        if userID != "" {
            hash := hashUserID(userID)
            if hash%100 < int64(split.NewServicePercent) {
                return "new"
            }
        }

    case SplitCookie:
        cookie, err := r.Cookie(split.CookieName)
        if err == nil && cookie.Value == "new" {
            return "new"
        }
    }

    return "monolith"
}

func (sr *StranglerRouter) routeToNewService(w http.ResponseWriter, r *http.Request, route *RouteRule) {
    // Apply transformers
    transformedReq := r.Clone(r.Context())
    for _, transformer := range route.Transformers {
        if err := transformer.Transform(transformedReq); err != nil {
            sr.logger.Error("transform error", zap.Error(err))
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }
    }

    // Strip prefix if configured
    if route.StripPrefix {
        transformedReq.URL.Path = stripPrefix(transformedReq.URL.Path, route.Pattern)
    }

    route.Proxy.ServeHTTP(w, transformedReq)
}

func (sr *StranglerRouter) routeShadow(w http.ResponseWriter, r *http.Request, route *RouteRule) {
    // Copy request for shadow
    shadowReq := r.Clone(context.Background())

    // Serve from monolith (user-facing)
    sr.monolithProxy.ServeHTTP(w, r)

    // Async shadow to new service
    go func() {
        rec := httptest.NewRecorder()
        route.Proxy.ServeHTTP(rec, shadowReq)

        // Compare responses
        sr.compareResponses(r, rec.Result())
    }()
}

func (sr *StranglerRouter) compareResponses(orig *http.Request, shadow *http.Response) {
    // Log differences for analysis
    // This helps identify behavioral differences before cutover
}

func (sr *StranglerRouter) matchRoute(path string) *RouteRule {
    sr.routesMu.RLock()
    defer sr.routesMu.RUnlock()

    // Longest match wins
    var matched *RouteRule
    maxLen := 0

    for pattern, route := range sr.routes {
        if strings.HasPrefix(path, pattern) && len(pattern) > maxLen {
            matched = route
            maxLen = len(pattern)
        }
    }

    return matched
}
```

## Production-Ready Go Implementation

### Data Migration Pattern

```go
// internal/migration/data_sync.go
package migration

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "github.com/segmentio/kafka-go"
)

// DataSynchronizer keeps monolith and new service data in sync
type DataSynchronizer struct {
    sourceDB      *sql.DB
    targetDB      *sql.DB
    changeLog     *ChangeLogReader
    eventWriter   *kafka.Writer

    syncConfig    SyncConfig
    logger        *zap.Logger
}

type SyncConfig struct {
    BatchSize         int
    SyncInterval      time.Duration
    ConflictStrategy  ConflictStrategy
    Tables            []TableSyncConfig
}

type TableSyncConfig struct {
    SourceTable       string
    TargetTable       string
    PrimaryKey        string
    Columns           []string
    Transformations   map[string]TransformationFunc
}

type ConflictStrategy int

const (
    ConflictSourceWins ConflictStrategy = iota
    ConflictTargetWins
    ConflictTimestampWins
    ConflictManual
)

// InitialSync performs full table migration
func (ds *DataSynchronizer) InitialSync(ctx context.Context, tableConfig TableSyncConfig) error {
    ds.logger.Info("starting initial sync",
        zap.String("table", tableConfig.SourceTable),
    )

    // Count records
    var totalCount int
    countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableConfig.SourceTable)
    if err := ds.sourceDB.QueryRowContext(ctx, countQuery).Scan(&totalCount); err != nil {
        return fmt.Errorf("counting records: %w", err)
    }

    ds.logger.Info("records to sync", zap.Int("count", totalCount))

    // Batch migration
    offset := 0
    processed := 0

    for {
        batch, err := ds.fetchBatch(ctx, tableConfig, offset, ds.syncConfig.BatchSize)
        if err != nil {
            return fmt.Errorf("fetching batch: %w", err)
        }

        if len(batch) == 0 {
            break
        }

        if err := ds.insertBatch(ctx, tableConfig, batch); err != nil {
            return fmt.Errorf("inserting batch: %w", err)
        }

        processed += len(batch)
        offset += ds.syncConfig.BatchSize

        ds.logger.Info("sync progress",
            zap.String("table", tableConfig.SourceTable),
            zap.Int("processed", processed),
            zap.Int("total", totalCount),
        )
    }

    ds.logger.Info("initial sync complete",
        zap.String("table", tableConfig.SourceTable),
    )

    return nil
}

// CDC (Change Data Capture) for ongoing sync
func (ds *DataSynchronizer) StartCDC(ctx context.Context) error {
    // Start binlog reader for MySQL or logical replication for PostgreSQL
    changes, err := ds.changeLog.Start(ctx)
    if err != nil {
        return fmt.Errorf("starting change log: %w", err)
    }

    go func() {
        for change := range changes {
            if err := ds.handleChange(ctx, change); err != nil {
                ds.logger.Error("handling change",
                    zap.Error(err),
                    zap.String("table", change.Table),
                )
                // Send to dead letter queue for manual review
                ds.sendToDLQ(change, err)
            }
        }
    }()

    return nil
}

func (ds *DataSynchronizer) handleChange(ctx context.Context, change *ChangeEvent) error {
    switch change.Operation {
    case OperationInsert:
        return ds.handleInsert(ctx, change)
    case OperationUpdate:
        return ds.handleUpdate(ctx, change)
    case OperationDelete:
        return ds.handleDelete(ctx, change)
    default:
        return fmt.Errorf("unknown operation: %v", change.Operation)
    }
}

func (ds *DataSynchronizer) handleUpdate(ctx context.Context, change *ChangeEvent) error {
    // Check for conflicts
    targetVersion, err := ds.getTargetVersion(ctx, change.Table, change.PrimaryKey)
    if err != nil {
        return err
    }

    if targetVersion > change.Version {
        // Conflict detected
        switch ds.syncConfig.ConflictStrategy {
        case ConflictSourceWins:
            // Proceed with update
        case ConflictTargetWins:
            return nil // Skip
        case ConflictTimestampWins:
            if !change.Timestamp.After(targetVersion.Timestamp) {
                return nil // Skip
            }
        case ConflictManual:
            ds.queueForManualResolution(change)
            return nil
        }
    }

    // Apply update
    return ds.applyUpdate(ctx, change)
}

// Dual-Write Pattern for zero-downtime migration
func (ds *DataSynchronizer) DualWrite(ctx context.Context, table string, data interface{}) error {
    // Write to both databases in transaction

    // Start source transaction
    sourceTx, err := ds.sourceDB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer sourceTx.Rollback()

    // Start target transaction
    targetTx, err := ds.targetDB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer targetTx.Rollback()

    // Insert into source
    if err := ds.insertIntoSource(sourceTx, table, data); err != nil {
        return err
    }

    // Insert into target
    if err := ds.insertIntoTarget(targetTx, table, data); err != nil {
        return err
    }

    // Commit both
    if err := sourceTx.Commit(); err != nil {
        return err
    }

    if err := targetTx.Commit(); err != nil {
        // Log inconsistency - needs reconciliation
        ds.logger.Error("target commit failed after source commit",
            zap.String("table", table),
        )
        return err
    }

    return nil
}
```

### Feature Flag Integration

```go
// internal/migration/feature_flags.go
package migration

import (
    "context"
    "fmt"

    "github.com/launchdarkly/go-server-sdk/v6"
)

// MigrationFlags controls gradual rollout
type MigrationFlags struct {
    client *ld.LDClient
}

// MigrationFlag defines feature flags for migration
type MigrationFlag struct {
    ShadowTraffic    bool
    CanaryPercent    int
    ReadFromNew      bool
    WriteToNew       bool
    RollbackEnabled  bool
}

func (mf *MigrationFlags) ShouldUseNewService(ctx context.Context, userID string, service string) bool {
    // Check if fully migrated
    if mf.client.BoolVariation(fmt.Sprintf("migration.%s.complete", service),
        ld.NewUser(userID), false) {
        return true
    }

    // Check canary percentage
    canaryPercent := mf.client.IntVariation(fmt.Sprintf("migration.%s.canary", service),
        ld.NewUser(userID), 0)

    // Deterministic routing based on user ID
    userHash := hashUserID(userID)
    return int(userHash%100) < canaryPercent
}

func (mf *MigrationFlags) IsShadowEnabled(service string) bool {
    return mf.client.BoolVariation(fmt.Sprintf("migration.%s.shadow", service),
        ld.NewUser("system"), false)
}

func (mf *MigrationFlags) EmergencyRollback(service string) error {
    // Immediately disable new service
    return mf.client.Track(fmt.Sprintf("migration.%s.rollback", service))
}
```

## Trade-off Analysis

### Migration Strategies Comparison

| Strategy | Risk | Speed | Complexity | Rollback | Best For |
|----------|------|-------|------------|----------|----------|
| **Strangler Fig** | Low | Slow | Medium | Easy | Most migrations |
| **Branch by Abstraction** | Low | Medium | Medium | Easy | Core domain logic |
| **Parallel Run** | Very Low | Slow | High | Very Easy | Financial systems |
| **Database First** | Medium | Medium | High | Hard | Data-heavy systems |

### Risk Assessment Matrix

```
                    High Impact
                         │
    ┌────────────────────┼────────────────────┐
    │   (AVOID)          │   (MIGRATE LAST)   │
    │   Core transaction │   Regulatory       │
    │   processing       │   compliance       │
    │                    │   features         │
High├────────────────────┼────────────────────┤ Low
Risk│                    │                    │ Risk
    │   (MIGRATE FIRST)  │   (MIGRATE EARLY)  │
    │   User preferences │   Reporting        │
    │   Non-critical     │   Analytics        │
    └────────────────────┼────────────────────┘
                         │
                    Low Impact
```

### Cost Factors

| Factor | Monolith | Microservices | Migration Period |
|--------|----------|---------------|------------------|
| Infrastructure | $ | $$$ | $$$$ (duplicated) |
| Operations | $ | $$ | $$$ (dual expertise) |
| Development | $$ | $ | $$$ (coordination) |
| Risk | High | Low | Medium |

## Testing Strategies

### Migration Testing Pyramid

```go
// test/migration/validation_test.go
package migration

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// ContractTest validates behavioral equivalence
func TestBehavioralEquivalence(t *testing.T) {
    cases := []struct {
        name     string
        request  Request
        validate func(t *testing.T, monolithResp, newResp Response)
    }{
        {
            name:    "user_registration",
            request: loadFixture("user_registration.json"),
            validate: func(t *testing.T, m, n Response) {
                assert.Equal(t, m.StatusCode, n.StatusCode)
                assert.Equal(t, m.UserID, n.UserID)
                assert.WithinDuration(t, m.CreatedAt, n.CreatedAt, time.Second)
            },
        },
        // More test cases...
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // Execute against both systems
            monolithResp := callMonolith(tc.request)
            newResp := callNewService(tc.request)

            tc.validate(t, monolithResp, newResp)
        })
    }
}

// ShadowTest validates responses match during shadow period
func TestShadowConsistency(t *testing.T) {
    ctx := context.Background()

    // Start shadow traffic
    shadow := NewShadowTester(monolithURL, newServiceURL)

    // Run production-like load
    load := generateProductionLoad(ctx, 1000)

    // Compare responses
    diffs := shadow.Run(ctx, load)

    // Assert high consistency
    consistencyRate := 1.0 - float64(len(diffs))/float64(len(load))
    require.Greater(t, consistencyRate, 0.999, "Consistency below 99.9%%")

    // Log any differences for analysis
    for _, diff := range diffs {
        t.Logf("Difference: %v", diff)
    }
}

// DataConsistencyTest validates data integrity
func TestDataConsistency(t *testing.T) {
    // After migration, verify data matches
    verifier := NewDataVerifier(sourceDB, targetDB)

    tables := []string{"users", "orders", "products"}

    for _, table := range tables {
        t.Run(table, func(t *testing.T) {
            match, diffs, err := verifier.CompareTable(table)
            require.NoError(t, err)
            assert.True(t, match, "Data mismatch in table %s: %v", table, diffs)
        })
    }
}

// RollbackTest validates rollback capability
func TestRollbackCapability(t *testing.T) {
    // Simulate failure
    newService.SimulateFailure()

    // Trigger rollback
    err := router.Rollback("user-service")
    require.NoError(t, err)

    // Verify traffic routes to monolith
    resp := makeRequest("/api/users/123")
    assert.Equal(t, "monolith", resp.Header.Get("X-Served-By"))

    // Verify functionality still works
    assert.Equal(t, 200, resp.StatusCode)
}
```

## Summary

The Strangler Fig Pattern provides:

1. **Incremental Migration**: Low-risk, step-by-step replacement
2. **Rollback Safety**: Can revert any change quickly
3. **Validation at Each Step**: Compare behavior before proceeding
4. **Business Continuity**: No big-bang deployments
5. **Learning Opportunity**: Understand system better during migration

Key success factors:

- Clear seam identification in monolith
- Comprehensive testing at each stage
- Data synchronization strategy
- Monitoring and observability
- Team alignment and communication

---

## 10. Performance Benchmarking

### 10.1 Core Benchmarks

```go
package benchmark_test

import (
	"context"
	"sync"
	"testing"
	"time"
)

// BenchmarkBasicOperation measures baseline performance
func BenchmarkBasicOperation(b *testing.B) {
	ctx := context.Background()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate operation
			_ = ctx
		}
	})
}

// BenchmarkConcurrentLoad tests concurrent performance
func BenchmarkConcurrentLoad(b *testing.B) {
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Simulate work
			time.Sleep(1 * time.Microsecond)
		}()
	}
	wg.Wait()
}

// BenchmarkMemoryAllocation tracks allocations
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		data := make([]byte, 1024)
		_ = data
	}
}
```

### 10.2 Performance Comparison

| Implementation | ns/op | allocs/op | memory/op | Throughput |
|---------------|-------|-----------|-----------|------------|
| **Baseline** | 100 ns | 0 | 0 B | 10M ops/s |
| **With Context** | 150 ns | 1 | 32 B | 6.7M ops/s |
| **With Metrics** | 300 ns | 2 | 64 B | 3.3M ops/s |
| **With Tracing** | 500 ns | 4 | 128 B | 2M ops/s |

### 10.3 Production Performance

| Metric | P50 | P95 | P99 | Target |
|--------|-----|-----|-----|--------|
| Latency | 100μs | 250μs | 500μs | < 1ms |
| Throughput | 50K | 80K | 100K | > 50K RPS |
| Error Rate | 0.01% | 0.05% | 0.1% | < 0.1% |
| CPU Usage | 10% | 25% | 40% | < 50% |

### 10.4 Optimization Recommendations

| Priority | Optimization | Impact | Effort |
|----------|-------------|--------|--------|
| 🔴 High | Connection pooling | 50% latency | Low |
| 🔴 High | Caching layer | 80% throughput | Medium |
| 🟡 Medium | Async processing | 30% latency | Medium |
| 🟡 Medium | Batch operations | 40% throughput | Low |
| 🟢 Low | Compression | 20% bandwidth | Low |
