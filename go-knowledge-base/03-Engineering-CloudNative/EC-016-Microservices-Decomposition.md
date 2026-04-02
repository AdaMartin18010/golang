# EC-016: Microservices Decomposition Patterns

## Problem Formalization

### The Monolithic Decomposition Dilemma

Decomposing monolithic applications into microservices represents one of the most challenging architectural transformations in software engineering. Organizations face a fundamental paradox: they understand the destination (microservices architecture) but struggle with the journey (decomposition strategy).

#### Mathematical Problem Definition

Given a monolithic system M with components C = {c₁, c₂, ..., cₙ} and dependencies D = {(cᵢ, cⱼ, wᵢⱼ)}, find an optimal partitioning P = {S₁, S₂, ..., Sₘ} where:

```
Maximize: Cohesion(Sₖ) for all Sₖ ∈ P
Minimize: Coupling(Sᵢ, Sⱼ) for all Sᵢ, Sⱼ ∈ P, i ≠ j
Subject to:
  - Bounded Context Integrity
  - Data Consistency Requirements
  - Transactional Boundaries
  - Team Organizational Constraints
  - Deployment Independence
```

**Cohesion Metric:**

```go
// Cohesion measures internal relatedness of a service
type CohesionMetrics struct {
    FunctionalCohesion    float64 // Related business functions
    DataCohesion         float64 // Shared data access patterns
    TemporalCohesion     float64 // Synchronized execution needs
    CommunicationCohesion float64 // Communication frequency
}

func CalculateCohesion(service Service, metrics CohesionMetrics) float64 {
    weights := []float64{0.4, 0.3, 0.2, 0.1}
    values := []float64{
        metrics.FunctionalCohesion,
        metrics.DataCohesion,
        metrics.TemporalCohesion,
        metrics.CommunicationCohesion,
    }

    var cohesion float64
    for i, w := range weights {
        cohesion += w * values[i]
    }
    return cohesion
}
```

**Coupling Metric:**

```go
// Coupling measures inter-service dependencies
type CouplingMetrics struct {
    InterfaceCoupling    int     // Number of API endpoints consumed
    DataCoupling         int     // Shared data entities
    ControlCoupling      int     // Control flow dependencies
    TemporalCoupling     int     // Synchronization requirements
}

func CalculateCoupling(serviceA, serviceB Service) float64 {
    metrics := serviceA.GetCouplingMetrics(serviceB)

    // Weighted coupling score
    weights := map[string]float64{
        "interface": 0.35,
        "data":      0.35,
        "control":   0.20,
        "temporal":  0.10,
    }

    return float64(metrics.InterfaceCoupling)*weights["interface"] +
           float64(metrics.DataCoupling)*weights["data"] +
           float64(metrics.ControlCoupling)*weights["control"] +
           float64(metrics.TemporalCoupling)*weights["temporal"]
}
```

### Anti-Patterns in Decomposition

#### 1. Distributed Monolith

When services are decomposed technically but remain tightly coupled, creating the worst of both worlds.

**Symptoms:**

- Services must be deployed together
- Changes in one service require changes in others
- Distributed transactions across many services
- Cascading failures on minor updates

#### 2. Nanoservices Anti-Pattern

Over-decomposition creating services with minimal functionality, leading to operational complexity.

**Detection Formula:**

```
Nanoservice Risk Score = (Operations per Day) / (Business Value × Code Lines)
Risk Score > 0.1 indicates potential nanoservice
```

#### 3. Shared Database Fallacy

Multiple services accessing the same database schema, creating hidden coupling through data.

#### 4. Wrong Bounded Contexts

Decomposing by technical layers rather than business capabilities.

### Domain-Driven Design Integration

#### Strategic Patterns for Decomposition

```go
// Bounded Context represents a boundary within which domain model applies
type BoundedContext struct {
    ID          string
    Name        string
    Description string
    Subdomains  []Subdomain
    UbiquitousLanguage map[string]string
    ContextMap  []ContextMapping
    Teams       []Team
}

// ContextMapping defines relationships between bounded contexts
type ContextMapping struct {
    SourceContextID string
    TargetContextID string
    Relationship    ContextRelationship
    TranslationMap  map[string]string // Terminology translation
}

type ContextRelationship int

const (
    Partnership ContextRelationship = iota
    SharedKernel
    CustomerSupplier
    Conformist
    AntiCorruptionLayer
    OpenHostService
    PublishedLanguage
    SeparateWays
)
```

## Solution Architecture

### Decomposition Strategies

#### 1. Decompose by Business Capability

```
┌─────────────────────────────────────────────────────────────────┐
│                    E-Commerce Platform                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Catalog    │  │    Order     │  │   Payment    │         │
│  │   Service    │  │   Service    │  │   Service    │         │
│  │              │  │              │  │              │         │
│  │ • Products   │  │ • Cart       │  │ • Processing │         │
│  │ • Categories │  │ • Checkout   │  │ • Refunds    │         │
│  │ • Inventory  │  │ • History    │  │ • Fraud      │         │
│  │ • Search     │  │ • Returns    │  │ • Reporting  │         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘         │
│         │                 │                 │                  │
│         └─────────────────┼─────────────────┘                  │
│                           │                                    │
│                           ▼                                    │
│              ┌──────────────────────┐                         │
│              │   Event Bus (Kafka)  │                         │
│              └──────────────────────┘                         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### 2. Decompose by Subdomain

```
┌─────────────────────────────────────────────────────────────────┐
│                    Banking System                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Core Domain (Competitive Advantage)                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐        │   │
│  │  │  Lending   │  │  Trading   │  │  Risk Mgmt │        │   │
│  │  │  Engine    │  │  Platform  │  │  System    │        │   │
│  │  └────────────┘  └────────────┘  └────────────┘        │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Supporting Subdomains                                          │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐               │
│  │  Customer  │  │  Document  │  │ Notification│               │
│  │  Mgmt      │  │  Mgmt      │  │  Service    │               │
│  └────────────┘  └────────────┘  └────────────┘               │
│                                                                 │
│  Generic Subdomains                                             │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐               │
│  │    Auth    │  │   Audit    │  │   Search   │               │
│  │  Service   │  │   Log      │  │   Engine   │               │
│  └────────────┘  └────────────┘  └────────────┘               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### 3. Strangler Fig Pattern (Incremental Migration)

```
┌─────────────────────────────────────────────────────────────────┐
│                    Strangler Fig Migration                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Phase 1: Identify Seams                                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Monolithic Application                      │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │   │
│  │  │   User   │ │  Order   │ │ Payment  │ │ Inventory│  │   │
│  │  │   Mgmt   │ │  Mgmt    │ │  Mgmt    │ │   Mgmt   │  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Phase 2: Extract User Service                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  ┌──────────────┐         ┌─────────────────────────┐  │   │
│  │  │ User Service │◄───────►│   Monolith (Facade)     │  │   │
│  │  │  (New)       │         │  ┌───────────────────┐  │  │   │
│  │  └──────────────┘         │  │ Order │Pay │ Inv  │  │  │   │
│  │                           │  └───────────────────┘  │  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Phase 3: Extract Order Service                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  ┌──────────┐  ┌──────────┐    ┌───────────────────┐   │   │
│  │  │  User    │  │  Order   │◄──►│   Monolith        │   │   │
│  │  │ Service  │  │ Service  │    │  ┌─────────────┐  │   │   │
│  │  └──────────┘  └──────────┘    │  │ Pay │ Inv   │  │   │   │
│  │                                │  └─────────────┘  │   │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Phase N: Complete Migration                                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐         │
│  │  User    │ │  Order   │ │ Payment  │ │ Inventory│         │
│  │ Service  │ │ Service  │ │ Service  │ │ Service  │         │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Service Boundary Definition

```go
// ServiceBoundaryAnalyzer determines optimal service boundaries
type ServiceBoundaryAnalyzer struct {
    codebaseAnalyzer    *CodebaseAnalyzer
    dependencyGraph     *DependencyGraph
    domainAnalyzer      *DomainAnalyzer
    teamStructure       []Team
}

func (sba *ServiceBoundaryAnalyzer) AnalyzeBoundaries() ([]ServiceBoundary, error) {
    // Step 1: Analyze code dependencies
    dependencies := sba.codebaseAnalyzer.ExtractDependencies()

    // Step 2: Build weighted dependency graph
    graph := sba.buildWeightedGraph(dependencies)

    // Step 3: Apply clustering algorithm
    clusters := sba.clusterByModularity(graph)

    // Step 4: Validate against domain model
    boundaries := sba.validateAgainstDomain(clusters)

    // Step 5: Apply team topology constraints
    boundaries = sba.applyTeamConstraints(boundaries)

    return boundaries, nil
}

func (sba *ServiceBoundaryAnalyzer) clusterByModularity(graph *DependencyGraph) []Cluster {
    // Louvain algorithm for community detection
    communities := louvainClustering(graph)

    // Merge small communities
    merged := sba.mergeSmallCommunities(communities, minServiceSize)

    // Split oversized communities
    final := sba.splitLargeCommunities(merged, maxServiceSize)

    return final
}
```

### Data Decomposition Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                    Data Decomposition Pattern                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Before: Shared Database                                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Shared Database                       │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │   │
│  │  │ users   │ │ orders  │ │payments │ │inventory│       │   │
│  │  │ table   │ │ table   │ │ table   │ │ table   │       │   │
│  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘       │   │
│  │       │           │           │           │             │   │
│  │       └───────────┴───────────┴───────────┘             │   │
│  │                   (Cross-table joins)                    │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  After: Database-per-Service                                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐         │
│  │  User    │ │  Order   │ │ Payment  │ │ Inventory│         │
│  │ Service  │ │ Service  │ │ Service  │ │ Service  │         │
│  │ ┌──────┐ │ │ ┌──────┐ │ │ ┌──────┐ │ │ ┌──────┐ │         │
│  │ │ User │ │ │ │Order │ │ │ │ Pay  │ │ │ │ Inv  │ │         │
│  │ │  DB  │ │ │ │  DB  │ │ │ │  DB  │ │ │ │  DB  │ │         │
│  │ └──────┘ │ │ └──────┘ │ │ └──────┘ │ │ └──────┘ │         │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘         │
│       │            │            │            │                │
│       └────────────┴────────────┴────────────┘                │
│              (Event-driven synchronization)                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Service Template with Bounded Context

```go
// internal/boundedcontext/user/service.go
package user

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "go.uber.org/zap"
    "golang.org/x/sync/singleflight"
)

// Domain Events
type UserRegisteredEvent struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    OccurredAt time.Time `json:"occurred_at"`
}

type UserProfileUpdatedEvent struct {
    UserID    string            `json:"user_id"`
    Changes   map[string]interface{} `json:"changes"`
    OccurredAt time.Time         `json:"occurred_at"`
}

// Aggregate Root
type User struct {
    ID          string
    Email       string
    Profile     UserProfile
    Status      UserStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Version     int // Optimistic locking
    Events      []DomainEvent
}

type UserStatus int

const (
    UserStatusPending UserStatus = iota
    UserStatusActive
    UserStatusSuspended
    UserStatusDeleted
)

type UserProfile struct {
    FirstName string
    LastName  string
    Phone     string
    Address   Address
    Preferences UserPreferences
}

type Address struct {
    Street     string
    City       string
    PostalCode string
    Country    string
}

type UserPreferences struct {
    Language          string
    Currency          string
    NotificationPrefs NotificationPreferences
}

type NotificationPreferences struct {
    EmailEnabled bool
    SMSEnabled   bool
    PushEnabled  bool
}

type DomainEvent interface {
    EventName() string
    OccurredAt() time.Time
    AggregateID() string
}

func (e UserRegisteredEvent) EventName() string    { return "user.registered" }
func (e UserRegisteredEvent) OccurredAt() time.Time { return e.OccurredAt }
func (e UserRegisteredEvent) AggregateID() string   { return e.UserID }

// Repository Interface (Port)
type Repository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Save(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter UserFilter) ([]*User, error)
    Count(ctx context.Context, filter UserFilter) (int64, error)
}

type UserFilter struct {
    Status    *UserStatus
    Email     string
    CreatedAfter  *time.Time
    CreatedBefore *time.Time
    Limit     int
    Offset    int
}

// Event Publisher Interface
type EventPublisher interface {
    Publish(ctx context.Context, event DomainEvent) error
    PublishBatch(ctx context.Context, events []DomainEvent) error
}

// Service Implementation
type Service struct {
    repo       Repository
    publisher  EventPublisher
    validator  *UserValidator
    logger     *zap.Logger
    sf         singleflight.Group
    metrics    *ServiceMetrics
}

type ServiceConfig struct {
    Repository  Repository
    Publisher   EventPublisher
    Logger      *zap.Logger
    Metrics     *ServiceMetrics
}

func NewService(cfg ServiceConfig) *Service {
    return &Service{
        repo:      cfg.Repository,
        publisher: cfg.Publisher,
        validator: NewUserValidator(),
        logger:    cfg.Logger,
        metrics:   cfg.Metrics,
    }
}

// RegisterUser handles user registration with idempotency
type RegisterUserRequest struct {
    Email     string
    Password  string
    FirstName string
    LastName  string
    IdempotencyKey string
}

type RegisterUserResponse struct {
    UserID    string
    Email     string
    Status    UserStatus
    CreatedAt time.Time
}

func (s *Service) RegisterUser(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, error) {
    start := time.Now()
    defer func() {
        s.metrics.RecordOperationDuration("register_user", time.Since(start))
    }()

    // Idempotency check using singleflight
    result, err, _ := s.sf.Do(req.IdempotencyKey, func() (interface{}, error) {
        return s.registerUserInternal(ctx, req)
    })

    if err != nil {
        s.metrics.RecordOperationError("register_user", err)
        return nil, err
    }

    return result.(*RegisterUserResponse), nil
}

func (s *Service) registerUserInternal(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, error) {
    // Validate input
    if err := s.validator.ValidateRegistration(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Check email uniqueness
    existing, err := s.repo.FindByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, ErrUserNotFound) {
        return nil, fmt.Errorf("checking email existence: %w", err)
    }
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }

    // Create user aggregate
    user := &User{
        ID:        uuid.New().String(),
        Email:     req.Email,
        Status:    UserStatusPending,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Version:   1,
        Profile: UserProfile{
            FirstName: req.FirstName,
            LastName:  req.LastName,
        },
    }

    // Hash password (delegated to auth service via event)
    // In production, use bcrypt or Argon2

    // Persist user
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("saving user: %w", err)
    }

    // Publish domain event
    event := UserRegisteredEvent{
        UserID:     user.ID,
        Email:      user.Email,
        OccurredAt: time.Now().UTC(),
    }

    if err := s.publisher.Publish(ctx, event); err != nil {
        // Log error but don't fail - event will be replayed by outbox pattern
        s.logger.Error("failed to publish event",
            zap.String("event", event.EventName()),
            zap.Error(err),
        )
    }

    s.metrics.RecordUserCreated()

    return &RegisterUserResponse{
        UserID:    user.ID,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
    }, nil
}

// UpdateProfile handles profile updates with optimistic locking
type UpdateProfileRequest struct {
    UserID    string
    Profile   UserProfile
    ExpectedVersion int
}

func (s *Service) UpdateProfile(ctx context.Context, req UpdateProfileRequest) error {
    start := time.Now()
    defer func() {
        s.metrics.RecordOperationDuration("update_profile", time.Since(start))
    }()

    // Get current user
    user, err := s.repo.FindByID(ctx, req.UserID)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return ErrUserNotFound
        }
        return fmt.Errorf("finding user: %w", err)
    }

    // Optimistic locking check
    if user.Version != req.ExpectedVersion {
        return ErrConcurrentModification
    }

    // Calculate changes for event
    changes := calculateProfileChanges(user.Profile, req.Profile)
    if len(changes) == 0 {
        return nil // No changes, idempotent
    }

    // Update aggregate
    user.Profile = req.Profile
    user.UpdatedAt = time.Now().UTC()
    user.Version++
    user.Events = append(user.Events, UserProfileUpdatedEvent{
        UserID:     user.ID,
        Changes:    changes,
        OccurredAt: time.Now().UTC(),
    })

    // Persist with version check
    if err := s.repo.Update(ctx, user); err != nil {
        if errors.Is(err, ErrVersionConflict) {
            return ErrConcurrentModification
        }
        return fmt.Errorf("updating user: %w", err)
    }

    // Publish events
    for _, event := range user.Events {
        if err := s.publisher.Publish(ctx, event); err != nil {
            s.logger.Error("failed to publish event",
                zap.String("event", event.EventName()),
                zap.Error(err),
            )
        }
    }

    return nil
}

func calculateProfileChanges(old, new UserProfile) map[string]interface{} {
    changes := make(map[string]interface{})

    if old.FirstName != new.FirstName {
        changes["first_name"] = map[string]string{
            "from": old.FirstName,
            "to":   new.FirstName,
        }
    }
    if old.LastName != new.LastName {
        changes["last_name"] = map[string]string{
            "from": old.LastName,
            "to":   new.LastName,
        }
    }
    // ... more fields

    return changes
}

// Errors
var (
    ErrUserNotFound          = errors.New("user not found")
    ErrEmailAlreadyExists    = errors.New("email already exists")
    ErrConcurrentModification = errors.New("concurrent modification detected")
    ErrVersionConflict       = errors.New("version conflict")
)

// ServiceMetrics for observability
type ServiceMetrics struct {
    // Implementation using Prometheus or similar
}

func (m *ServiceMetrics) RecordOperationDuration(operation string, duration time.Duration) {}
func (m *ServiceMetrics) RecordOperationError(operation string, err error) {}
func (m *ServiceMetrics) RecordUserCreated() {}
```

### Database Per Service Implementation

```go
// internal/boundedcontext/user/repository/postgres.go
package repository

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"

    "github.com/lib/pq"
    "github.com/company/project/internal/boundedcontext/user"
)

type PostgresRepository struct {
    db     *sql.DB
    schema string
}

func NewPostgresRepository(db *sql.DB, schema string) *PostgresRepository {
    return &PostgresRepository{
        db:     db,
        schema: schema,
    }
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
    query := fmt.Sprintf(`
        SELECT id, email, status, profile, created_at, updated_at, version
        FROM %s.users
        WHERE id = $1 AND deleted_at IS NULL
    `, r.schema)

    var u user.User
    var profileJSON []byte

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &u.ID,
        &u.Email,
        &u.Status,
        &profileJSON,
        &u.CreatedAt,
        &u.UpdatedAt,
        &u.Version,
    )

    if err == sql.ErrNoRows {
        return nil, user.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query user: %w", err)
    }

    if err := json.Unmarshal(profileJSON, &u.Profile); err != nil {
        return nil, fmt.Errorf("unmarshal profile: %w", err)
    }

    return &u, nil
}

func (r *PostgresRepository) Save(ctx context.Context, u *user.User) error {
    profileJSON, err := json.Marshal(u.Profile)
    if err != nil {
        return fmt.Errorf("marshal profile: %w", err)
    }

    query := fmt.Sprintf(`
        INSERT INTO %s.users (id, email, status, profile, created_at, updated_at, version)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, r.schema)

    _, err = r.db.ExecContext(ctx, query,
        u.ID,
        u.Email,
        u.Status,
        profileJSON,
        u.CreatedAt,
        u.UpdatedAt,
        u.Version,
    )

    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            return user.ErrEmailAlreadyExists
        }
        return fmt.Errorf("insert user: %w", err)
    }

    return nil
}

func (r *PostgresRepository) Update(ctx context.Context, u *user.User) error {
    profileJSON, err := json.Marshal(u.Profile)
    if err != nil {
        return fmt.Errorf("marshal profile: %w", err)
    }

    query := fmt.Sprintf(`
        UPDATE %s.users
        SET email = $1, status = $2, profile = $3, updated_at = $4, version = $5
        WHERE id = $6 AND version = $7 AND deleted_at IS NULL
    `, r.schema)

    result, err := r.db.ExecContext(ctx, query,
        u.Email,
        u.Status,
        profileJSON,
        u.UpdatedAt,
        u.Version,
        u.ID,
        u.Version-1, // Expected previous version
    )

    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }

    if rows == 0 {
        return user.ErrVersionConflict
    }

    return nil
}

// Migration for isolated schema
const migration = `
CREATE SCHEMA IF NOT EXISTS user_service;

CREATE TABLE IF NOT EXISTS user_service.users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    profile JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_users_email ON user_service.users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON user_service.users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON user_service.users(created_at);

-- Outbox table for event publishing
CREATE TABLE IF NOT EXISTS user_service.outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP WITH TIME ZONE,
    publish_attempts INTEGER DEFAULT 0
);

CREATE INDEX idx_outbox_unpublished ON user_service.outbox(published_at) WHERE published_at IS NULL;
`
```

### Event-Driven Communication

```go
// internal/boundedcontext/user/events/publisher.go
package events

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/IBM/sarama"
    "github.com/company/project/internal/boundedcontext/user"
)

// KafkaEventPublisher implements event publishing via Kafka
type KafkaEventPublisher struct {
    producer sarama.AsyncProducer
    topic    string
}

func NewKafkaEventPublisher(brokers []string, topic string) (*KafkaEventPublisher, error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForLocal
    config.Producer.Compression = sarama.CompressionSnappy
    config.Producer.Flush.Frequency = 500 * time.Millisecond
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true

    producer, err := sarama.NewAsyncProducer(brokers, config)
    if err != nil {
        return nil, fmt.Errorf("creating kafka producer: %w", err)
    }

    go func() {
        for err := range producer.Errors() {
            // Log and handle errors, potentially retry
            fmt.Printf("Failed to produce message: %v\n", err)
        }
    }()

    return &KafkaEventPublisher{
        producer: producer,
        topic:    topic,
    }, nil
}

func (p *KafkaEventPublisher) Publish(ctx context.Context, event user.DomainEvent) error {
    payload, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshal event: %w", err)
    }

    msg := &sarama.ProducerMessage{
        Topic: p.topic,
        Key:   sarama.StringEncoder(event.AggregateID()),
        Value: sarama.ByteEncoder(payload),
        Headers: []sarama.RecordHeader{
            {
                Key:   []byte("event_type"),
                Value: []byte(event.EventName()),
            },
            {
                Key:   []byte("occurred_at"),
                Value: []byte(event.OccurredAt().Format(time.RFC3339)),
            },
        },
    }

    select {
    case p.producer.Input() <- msg:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (p *KafkaEventPublisher) Close() error {
    return p.producer.Close()
}

// Outbox pattern for guaranteed delivery
type OutboxPublisher struct {
    db        *sql.DB
    publisher user.EventPublisher
}

func (p *OutboxPublisher) Publish(ctx context.Context, event user.DomainEvent) error {
    // Save to outbox table in same transaction as business logic
    payload, _ := json.Marshal(event)

    query := `
        INSERT INTO user_service.outbox (aggregate_id, aggregate_type, event_type, payload, occurred_at)
        VALUES ($1, $2, $3, $4, $5)
    `

    _, err := p.db.ExecContext(ctx, query,
        event.AggregateID(),
        "user",
        event.EventName(),
        payload,
        event.OccurredAt(),
    )

    return err
}

// OutboxRelay polls outbox table and publishes events
func (p *OutboxPublisher) StartRelay(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            p.processOutbox(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (p *OutboxPublisher) processOutbox(ctx context.Context) {
    query := `
        SELECT id, aggregate_id, event_type, payload, occurred_at
        FROM user_service.outbox
        WHERE published_at IS NULL AND publish_attempts < 10
        ORDER BY occurred_at
        LIMIT 100
        FOR UPDATE SKIP LOCKED
    `

    rows, err := p.db.QueryContext(ctx, query)
    if err != nil {
        return
    }
    defer rows.Close()

    for rows.Next() {
        var event OutboxEvent
        if err := rows.Scan(&event.ID, &event.AggregateID, &event.EventType, &event.Payload, &event.OccurredAt); err != nil {
            continue
        }

        // Publish via actual publisher
        // Mark as published or increment attempts
    }
}
```

## Trade-off Analysis

### Decomposition Strategy Comparison

| Strategy | Pros | Cons | Best For |
|----------|------|------|----------|
| **By Business Capability** | Clear ownership, aligns with org | May have data duplication | Most greenfield projects |
| **By Subdomain** | DDD alignment, clear boundaries | Requires domain expertise | Complex domain systems |
| **By Transaction** | Data consistency preserved | Limited decomposition | Financial systems |
| **By Volatility** | Isolates change | Artificial boundaries | Legacy modernization |

### Cohesion vs. Coupling Matrix

```
                    High Coupling
                         │
    ┌────────────────────┼────────────────────┐
    │   (AVOID)          │   (REFACTOR)       │
    │   Big Ball of Mud  │   Distributed      │
    │                    │   Monolith         │
Low ├────────────────────┼────────────────────┤ High
Cohesion│                    │                    │ Cohesion
    │   (CONSOLIDATE)    │   (TARGET)         │
    │   Nanoservices     │   Well-Designed    │
    │                    │   Microservices    │
    └────────────────────┼────────────────────┘
                         │
                    Low Coupling
```

### Data Consistency Patterns

| Pattern | Consistency | Complexity | Latency | Use Case |
|---------|-------------|------------|---------|----------|
| **Saga (Orchestration)** | Eventual | High | Medium | Long-running processes |
| **Saga (Choreography)** | Eventual | Medium | Low | Simple workflows |
| **Two-Phase Commit** | Strong | Very High | High | Critical financial |
| **Transactional Outbox** | Eventual | Medium | Low | Most use cases |

### Organizational Impact

```
Team Topology Impact:
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Stream-Aligned Teams     │  Platform Teams                │
│  ┌───────────────────┐    │  ┌─────────────────────────┐   │
│  │ One team per      │    │  │ Shared services:        │   │
│  │ bounded context   │    │  │ - API Gateway           │   │
│  │                   │    │  │ - Message Bus           │   │
│  │ Benefits:         │    │  │ - Observability         │   │
│  │ - Clear ownership │    │  │ - Service Mesh          │   │
│  │ - Fast iteration  │    │  └─────────────────────────┘   │
│  │ - Domain focus    │    │                                │
│  └───────────────────┘    │  Enabling Teams                │
│                           │  ┌─────────────────────────┐   │
│  Complicated Subsystem    │  │ - Architecture guidance │   │
│  ┌───────────────────┐    │  │ - Best practices        │   │
│  │ Specialized team  │    │  │ - Tooling support       │   │
│  │ for ML/AI, Search │    │  └─────────────────────────┘   │
│  └───────────────────┘    │                                │
│                           └────────────────────────────────┘
```

## Testing Strategies

### Testing Pyramid for Microservices

```
┌─────────────────────────────────────────────────────────────┐
│                    End-to-End Tests                          │
│                    (5% - Critical Paths)                     │
│         ┌─────────────────────────────────────┐              │
│         │   Contract Tests (Pact)             │              │
│         │   Service-to-Service Verification   │              │
│         └─────────────────────────────────────┘              │
│                    Integration Tests                         │
│                    (15% - Repository, Message Bus)           │
│         ┌─────────────────────────────────────┐              │
│         │   Component Tests                   │              │
│         │   (Testcontainers, In-memory DB)    │              │
│         └─────────────────────────────────────┘              │
│                    Unit Tests                                │
│                    (80% - Business Logic)                    │
└─────────────────────────────────────────────────────────────┘
```

### Contract Testing with Pact

```go
// test/contract/user_service_test.go
package contract

import (
    "testing"

    "github.com/pact-foundation/pact-go/dsl"
)

func TestUserServiceProvider(t *testing.T) {
    pact := &dsl.Pact{
        Provider: "user-service",
        Consumer: "order-service",
    }
    defer pact.Teardown()

    pact.VerifyProvider(t, dsl.VerifyRequest{
        ProviderBaseURL: "http://localhost:8080",
        PactURLs:        []string{"./pacts/order-service-user-service.json"},
        StateHandlers: dsl.StateHandlers{
            "user exists": func() error {
                // Setup test data
                return setupTestUser("test-user-id")
            },
            "user does not exist": func() error {
                return nil
            },
        },
    })
}

func TestUserServiceConsumer(t *testing.T) {
    pact := &dsl.Pact{
        Consumer: "order-service",
        Provider: "user-service",
    }
    defer pact.Teardown()

    pact.AddInteraction().
        Given("user exists").
        UponReceiving("a request for user details").
        WithRequest(dsl.Request{
            Method: "GET",
            Path:   dsl.String("/users/123"),
        }).
        WillRespondWith(dsl.Response{
            Status: 200,
            Body: dsl.Like(map[string]interface{}{
                "id":     dsl.String("123"),
                "email":  dsl.String("user@example.com"),
                "status": dsl.String("active"),
            }),
        })

    err := pact.Verify(func() error {
        _, err := userClient.GetUser("123")
        return err
    })

    if err != nil {
        t.Fatalf("Error on Verify: %v", err)
    }
}
```

### Component Testing

```go
// test/component/user_service_component_test.go
package component

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
)

type UserServiceComponentTestSuite struct {
    suite.Suite
    postgresContainer *postgres.PostgresContainer
    service          *user.Service
    ctx              context.Context
}

func (s *UserServiceComponentTestSuite) SetupSuite() {
    s.ctx = context.Background()

    // Start PostgreSQL container
    container, err := postgres.RunContainer(s.ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Second),
        ),
    )
    s.Require().NoError(err)
    s.postgresContainer = container

    // Get connection string
    connStr, err := container.ConnectionString(s.ctx)
    s.Require().NoError(err)

    // Initialize database
    db, err := sql.Open("postgres", connStr)
    s.Require().NoError(err)

    // Run migrations
    _, err = db.Exec(migration)
    s.Require().NoError(err)

    // Setup service with test doubles
    repo := repository.NewPostgresRepository(db, "user_service")
    publisher := &mockEventPublisher{}

    s.service = user.NewService(user.ServiceConfig{
        Repository: repo,
        Publisher:  publisher,
        Logger:     zap.NewNop(),
    })
}

func (s *UserServiceComponentTestSuite) TearDownSuite() {
    if s.postgresContainer != nil {
        s.postgresContainer.Terminate(s.ctx)
    }
}

func (s *UserServiceComponentTestSuite) TestRegisterUser() {
    req := user.RegisterUserRequest{
        Email:          "test@example.com",
        FirstName:      "John",
        LastName:       "Doe",
        IdempotencyKey: "unique-key-123",
    }

    resp, err := s.service.RegisterUser(s.ctx, req)
    s.Require().NoError(err)
    s.NotEmpty(resp.UserID)
    s.Equal(req.Email, resp.Email)

    // Verify idempotency
    resp2, err := s.service.RegisterUser(s.ctx, req)
    s.Require().NoError(err)
    s.Equal(resp.UserID, resp2.UserID)
}

func (s *UserServiceComponentTestSuite) TestConcurrentProfileUpdate() {
    // Create user first
    createReq := user.RegisterUserRequest{
        Email:          "concurrent@example.com",
        FirstName:      "Jane",
        LastName:       "Doe",
        IdempotencyKey: "concurrent-key",
    }
    resp, err := s.service.RegisterUser(s.ctx, createReq)
    s.Require().NoError(err)

    // Attempt concurrent updates
    var wg sync.WaitGroup
    errors := make(chan error, 2)

    wg.Add(2)
    go func() {
        defer wg.Done()
        err := s.service.UpdateProfile(s.ctx, user.UpdateProfileRequest{
            UserID:          resp.UserID,
            ExpectedVersion: 1,
            Profile: user.UserProfile{
                FirstName: "Update1",
            },
        })
        errors <- err
    }()

    go func() {
        defer wg.Done()
        err := s.service.UpdateProfile(s.ctx, user.UpdateProfileRequest{
            UserID:          resp.UserID,
            ExpectedVersion: 1,
            Profile: user.UserProfile{
                FirstName: "Update2",
            },
        })
        errors <- err
    }()

    wg.Wait()
    close(errors)

    // One should succeed, one should fail with concurrent modification
    var successCount, conflictCount int
    for err := range errors {
        if err == nil {
            successCount++
        } else if errors.Is(err, user.ErrConcurrentModification) {
            conflictCount++
        }
    }

    s.Equal(1, successCount)
    s.Equal(1, conflictCount)
}

func TestUserServiceComponent(t *testing.T) {
    suite.Run(t, new(UserServiceComponentTestSuite))
}
```

### Integration Test for Event Publishing

```go
// test/integration/event_publishing_test.go
package integration

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "github.com/IBM/sarama"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/kafka"
)

func TestEventPublishingIntegration(t *testing.T) {
    ctx := context.Background()

    // Start Kafka container
    kafkaContainer, err := kafka.RunContainer(ctx,
        testcontainers.WithImage("confluentinc/cp-kafka:latest"),
    )
    require.NoError(t, err)
    defer kafkaContainer.Terminate(ctx)

    brokers, err := kafkaContainer.Brokers(ctx)
    require.NoError(t, err)

    // Create topic
    // ... setup code

    // Create publisher
    publisher, err := events.NewKafkaEventPublisher(brokers, "user-events")
    require.NoError(t, err)
    defer publisher.Close()

    // Create consumer for verification
    consumer := createTestConsumer(brokers, "user-events")

    // Publish event
    event := user.UserRegisteredEvent{
        UserID:     "test-id",
        Email:      "test@example.com",
        OccurredAt: time.Now().UTC(),
    }

    err = publisher.Publish(ctx, event)
    require.NoError(t, err)

    // Verify event received
    select {
    case msg := <-consumer.Messages():
        var received user.UserRegisteredEvent
        err := json.Unmarshal(msg.Value, &received)
        require.NoError(t, err)
        require.Equal(t, event.UserID, received.UserID)
    case <-time.After(10 * time.Second):
        t.Fatal("timeout waiting for message")
    }
}
```

## Summary

Microservices decomposition requires careful analysis of domain boundaries, technical constraints, and organizational structure. Success depends on:

1. **Strong Domain Modeling**: Use DDD to identify bounded contexts
2. **Incremental Migration**: Apply Strangler Fig pattern for gradual transition
3. **Data Isolation**: Implement database-per-service with event synchronization
4. **Observability**: Comprehensive metrics and distributed tracing
5. **Testing Strategy**: Contract tests, component tests, and chaos engineering

Key metrics to track:

- Service independence score (deployment frequency correlation)
- Cross-service call ratio (target <20% of total calls)
- Data consistency lag (p95 < 500ms)
- Team cognitive load (services per team < 10)
