# AD-XXX: [Application Domain Topic] - Quick Contribution Template

> **Dimension**: Application Domains (AD)
> **Level**: S/A/B - Target >[TODO: 15KB/10KB/5KB]
> **Status**: [TODO: Draft/Review/Complete]
> **Tags**: #[TODO: architecture] #[TODO: system-design] #[TODO: domain]
> **Author**: [TODO: Your Name]
> **Created**: [TODO: YYYY-MM-DD]
> **Estimated Reading Time**: [TODO: XX minutes]

---

## Table of Contents

1. [AD-XXX: [Application Domain Topic] - Quick Contribution Template](#executive-summary)
2. [Introduction](#introduction)
3. [Requirements Analysis](#requirements-analysis)
4. [Architecture Design](#architecture-design)
5. [Component Details](#component-details)
6. [Implementation Guide](#implementation-guide)
7. [Visual Representations](#visual-representations)
8. [Code Examples](#code-examples)
9. [Operational Considerations](#operational-considerations)
10. [Cross-References](#cross-references)
11. [References](#references)

---

## Executive Summary

[TODO: 2-3 paragraph overview for architects and tech leads]

**System at a Glance**:

- **Type**: [TODO: E-commerce/IoT/FinTech/etc.]
- **Scale**: [TODO: Expected traffic/data volume]
- **Complexity**: [TODO: Low/Medium/High]
- **Key Technologies**: [TODO: Go, PostgreSQL, Redis, etc.]

---

## Introduction

### Problem Domain

[TODO: What business/technical problem does this system solve?]

### Scope

**In Scope**:

- [TODO: Feature 1]
- [TODO: Feature 2]
- [TODO: Feature 3]

**Out of Scope**:

- [TODO: Out of scope item 1]
- [TODO: Out of scope item 2]

### Target Audience

- [TODO: Software Architects]
- [TODO: Backend Engineers]
- [TODO: DevOps Engineers]

### Prerequisites

- [TODO: [Microservices Patterns](../03-Engineering-CloudNative/EC-001-Microservices.md)]
- [TODO: [Domain-Driven Design](../03-Engineering-CloudNative/EC-XXX-DDD.md)]

---

## Requirements Analysis

### Functional Requirements

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| FR-1 | [TODO: Requirement] | Must | [TODO] |
| FR-2 | [TODO: Requirement] | Must | [TODO] |
| FR-3 | [TODO: Requirement] | Should | [TODO] |

### Non-Functional Requirements

| ID | Requirement | Target | Measurement |
|----|-------------|--------|-------------|
| NFR-1 | Availability | 99.99% | Uptime SLA |
| NFR-2 | Latency (p99) | < 100ms | API response time |
| NFR-3 | Throughput | 10K RPS | Peak load |
| NFR-4 | Data Durability | 99.9999% | Backup success rate |

### Constraints

- **Budget**: [TODO]
- **Timeline**: [TODO]
- **Compliance**: [TODO: GDPR/SOC2/etc.]
- **Technology Stack**: [TODO: Must use Go]

---

## Architecture Design

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    SYSTEM ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐  │
│   │                     Client Layer                        │  │
│   │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐  │  │
│   │  │  Web App │  │ Mobile   │  │  CLI     │  │  API   │  │  │
│   │  └──────────┘  └──────────┘  └──────────┘  └────────┘  │  │
│   └────────────────────────┬────────────────────────────────┘  │
│                            │                                    │
│                            ▼                                    │
│   ┌─────────────────────────────────────────────────────────┐  │
│   │                    API Gateway Layer                    │  │
│   │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐  │  │
│   │  │  Rate    │  │   Auth   │  │  Routing │  │  SSL   │  │  │
│   │  │  Limit   │  │          │  │          │  │Term.   │  │  │
│   │  └──────────┘  └──────────┘  └──────────┘  └────────┘  │  │
│   └────────────────────────┬────────────────────────────────┘  │
│                            │                                    │
│                            ▼                                    │
│   ┌─────────────────────────────────────────────────────────┐  │
│   │                  Service Layer (Go)                     │  │
│   │  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐        │  │
│   │  │Service1│  │Service2│  │Service3│  │Service4│        │  │
│   │  │ (Go)   │  │ (Go)   │  │ (Go)   │  │ (Go)   │        │  │
│   │  └────┬───┘  └────┬───┘  └────┬───┘  └────┬───┘        │  │
│   │       └───────────┴───────────┴───────────┘            │  │
│   │                    Message Bus                         │  │
│   └────────────────────────┬────────────────────────────────┘  │
│                            │                                    │
│           ┌────────────────┼────────────────┐                  │
│           ▼                ▼                ▼                  │
│   ┌────────────┐  ┌────────────┐  ┌────────────┐              │
│   │PostgreSQL  │  │   Redis    │  │ Elasticsearch│              │
│   │ (Primary)  │  │  (Cache)   │  │  (Search)   │              │
│   └────────────┘  └────────────┘  └────────────┘              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Architecture Patterns

**Patterns Used**:

- [TODO: Microservices Architecture]
- [TODO: CQRS (Command Query Responsibility Segregation)]
- [TODO: Event Sourcing]
- [TODO: Saga Pattern for distributed transactions]

**Rationale**: [TODO: Why these patterns were chosen]

### Data Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     DATA ARCHITECTURE                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐                                              │
│   │   Command    │───┐                                          │
│   │   Side       │   │  Write Operations                        │
│   └──────────────┘   │                                          │
│                      ▼                                          │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│   │  Event       │  │  Aggregate   │  │  Event       │          │
│   │  Store       │◀─│   Store      │◀─│  Bus         │          │
│   └──────────────┘  └──────────────┘  └──────────────┘          │
│          │               │                                      │
│          │               │                                      │
│          ▼               ▼                                      │
│   ┌──────────────┐  ┌──────────────┐                            │
│   │  Read Model  │  │  Snapshot    │                            │
│   │  Projector   │  │  Store       │                            │
│   └──────┬───────┘  └──────────────┘                            │
│          │                                                      │
│          ▼                                                      │
│   ┌──────────────┐                                              │
│   │   Query      │───┐  Read Operations                         │
│   │   Side       │   │                                          │
│   └──────────────┘   │                                          │
│                      ▼                                          │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│   │  Read DB     │  │    Cache     │  │   Search     │          │
│   │ (PostgreSQL) │  │   (Redis)    │  │ (Elastic)    │          │
│   └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Component Details

### Component 1: [Service Name]

**Responsibility**: [TODO: What this service does]

**API Specification**:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/resource` | GET | [TODO] |
| `/api/v1/resource` | POST | [TODO] |
| `/api/v1/resource/:id` | PUT | [TODO] |
| `/api/v1/resource/:id` | DELETE | [TODO] |

**Data Model**:

```go
// file: models.go
// description: Domain models for [Service]
package domain

import (
    "time"

    "github.com/google/uuid"
)

// [Entity] represents [description].
type Entity struct {
    ID        uuid.UUID `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Status    Status    `json:"status" db:"status"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
    Version   int       `json:"version" db:"version"`
}

// Status represents entity state.
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
)

// Validate validates the entity.
func (e *Entity) Validate() error {
    if e.Name == "" {
        return fmt.Errorf("name is required")
    }
    return nil
}
```

**Dependencies**:

- [TODO: Database]
- [TODO: Cache]
- [TODO: Other services]

### Component 2: [Service Name]

[TODO: Additional components]

---

## Implementation Guide

### Project Structure

```
my-service/
├── cmd/
│   ├── api/
│   │   └── main.go              # API server entry point
│   ├── worker/
│   │   └── main.go              # Background worker entry point
│   └── migrate/
│       └── main.go              # Database migration tool
├── internal/
│   ├── application/
│   │   ├── commands/            # Command handlers (CQRS)
│   │   ├── queries/             # Query handlers (CQRS)
│   │   └── services/            # Application services
│   ├── domain/
│   │   ├── entity/              # Domain entities
│   │   ├── events/              # Domain events
│   │   ├── repository/          # Repository interfaces
│   │   └── service/             # Domain services
│   ├── infrastructure/
│   │   ├── persistence/         # Repository implementations
│   │   │   ├── postgres/
│   │   │   └── redis/
│   │   ├── messaging/           # Event bus implementation
│   │   ├── http/                # HTTP handlers
│   │   └── config/              # Configuration
│   └── interfaces/
│       └── http/                # HTTP controllers
├── pkg/
│   ├── errors/                  # Shared errors
│   ├── logger/                  # Shared logger
│   └── middleware/              # Shared middleware
├── api/
│   └── openapi.yaml             # OpenAPI specification
├── deployments/
│   ├── docker-compose.yml
│   └── kubernetes/
├── migrations/
│   └── *.sql
├── go.mod
├── go.sum
└── README.md
```

### Core Implementation

```go
// file: internal/application/commands/create_entity.go
// description: Command handler for creating entities
package commands

import (
    "context"
    "fmt"

    "github.com/google/uuid"
)

// CreateEntityCommand represents the create command.
type CreateEntityCommand struct {
    Name string
    // [TODO: Other fields]
}

// CreateEntityHandler handles CreateEntityCommand.
type CreateEntityHandler struct {
    repo      repository.EntityRepository
    eventBus  messaging.EventBus
    logger    *zap.Logger
}

// NewCreateEntityHandler creates a new handler.
func NewCreateEntityHandler(
    repo repository.EntityRepository,
    eventBus messaging.EventBus,
    logger *zap.Logger,
) *CreateEntityHandler {
    return &CreateEntityHandler{
        repo:     repo,
        eventBus: eventBus,
        logger:   logger,
    }
}

// Handle processes the command.
func (h *CreateEntityHandler) Handle(ctx context.Context, cmd CreateEntityCommand) (*domain.Entity, error) {
    // Create entity
    entity := &domain.Entity{
        ID:        uuid.New(),
        Name:      cmd.Name,
        Status:    domain.StatusPending,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Version:   1,
    }

    // Validate
    if err := entity.Validate(); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Persist
    if err := h.repo.Save(ctx, entity); err != nil {
        return nil, fmt.Errorf("save failed: %w", err)
    }

    // Publish event
    event := events.EntityCreated{
        EntityID:  entity.ID,
        Name:      entity.Name,
        CreatedAt: entity.CreatedAt,
    }

    if err := h.eventBus.Publish(ctx, event); err != nil {
        h.logger.Error("failed to publish event", zap.Error(err))
        // Don't fail the command if event publishing fails
    }

    h.logger.Info("entity created",
        zap.String("entity_id", entity.ID.String()),
        zap.String("name", entity.Name),
    )

    return entity, nil
}
```

### Infrastructure Layer

```go
// file: internal/infrastructure/persistence/postgres/entity_repository.go
// description: PostgreSQL implementation of EntityRepository
package postgres

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

type EntityRepository struct {
    db *pgxpool.Pool
}

func NewEntityRepository(db *pgxpool.Pool) *EntityRepository {
    return &EntityRepository{db: db}
}

func (r *EntityRepository) Save(ctx context.Context, entity *domain.Entity) error {
    query := `
        INSERT INTO entities (id, name, status, created_at, updated_at, version)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err := r.db.Exec(ctx, query,
        entity.ID,
        entity.Name,
        entity.Status,
        entity.CreatedAt,
        entity.UpdatedAt,
        entity.Version,
    )

    if err != nil {
        return fmt.Errorf("insert entity: %w", err)
    }

    return nil
}

func (r *EntityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Entity, error) {
    query := `
        SELECT id, name, status, created_at, updated_at, version
        FROM entities
        WHERE id = $1
    `

    entity := &domain.Entity{}
    err := r.db.QueryRow(ctx, query, id).Scan(
        &entity.ID,
        &entity.Name,
        &entity.Status,
        &entity.CreatedAt,
        &entity.UpdatedAt,
        &entity.Version,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, repository.ErrNotFound
        }
        return nil, fmt.Errorf("get entity: %w", err)
    }

    return entity, nil
}
```

---

## Visual Representations

### System Context Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    SYSTEM CONTEXT                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│                          ┌──────────┐                          │
│                          │  Users   │                          │
│                          └────┬─────┘                          │
│                               │                                  │
│                               │ Uses                              │
│                               ▼                                  │
│   ┌──────────┐         ┌──────────┐         ┌──────────┐        │
│   │ External │◀───────▶│ [System] │◀───────▶│ External │        │
│   │ System 1 │         │  Name    │         │ System 2 │        │
│   └──────────┘         └────┬─────┘         └──────────┘        │
│                             │                                    │
│                             │ Uses                                │
│                             ▼                                    │
│                      ┌──────────┐                               │
│                      │ Database │                               │
│                      └──────────┘                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Sequence Diagram

```
                    ┌──────────┐  ┌──────────┐  ┌──────────┐
                    │  Client  │  │   API    │  │ Service  │
                    └────┬─────┘  └────┬─────┘  └────┬─────┘
                         │              │              │
    1. Request           │─────────────▶│              │
                         │              │              │
    2. Validate          │              │──────▶       │
                         │              │              │
    3. Process           │              │       ──────▶│
                         │              │              │
    4. Persist           │              │       ──────▶│
                         │              │              │
    5. Event             │              │       ──────▶│
                         │              │              │
    6. Response          │              │◀──────       │
                         │◀─────────────│              │
                         │              │              │
```

### Deployment Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    DEPLOYMENT ARCHITECTURE                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   Kubernetes Cluster                                            │
│   ┌─────────────────────────────────────────────────────────┐  │
│   │  Namespace: production                                   │  │
│   │                                                          │  │
│   │  ┌────────────┐  ┌────────────┐  ┌────────────┐        │  │
│   │  │ API Pod 1  │  │ API Pod 2  │  │ API Pod 3  │        │  │
│   │  │ (Go 1.21)  │  │ (Go 1.21)  │  │ (Go 1.21)  │        │  │
│   │  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘        │  │
│   │        │               │               │               │  │
│   │        └───────────────┼───────────────┘               │  │
│   │                        │                                │  │
│   │                  ┌─────┴─────┐                          │  │
│   │                  │  Service  │                          │  │
│   │                  │  (LB)     │                          │  │
│   │                  └─────┬─────┘                          │  │
│   │                        │                                │  │
│   │  ┌─────────────────────┼────────────────────────────┐   │  │
│   │  │                     ▼                             │   │  │
│   │  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  │   │  │
│   │  │  │ PostgreSQL │  │   Redis    │  │  Kafka     │  │   │  │
│   │  │  │  Primary   │  │  Cluster   │  │  Cluster   │  │   │  │
│   │  │  └────────────┘  └────────────┘  └────────────┘  │   │  │
│   │  │                                                  │   │  │
│   │  └─────────────────── StatefulSet ──────────────────┘   │  │
│   │                                                          │  │
│   └─────────────────────────────────────────────────────────┘  │
│                                                                  │
│   Ingress: api.example.com                                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Code Examples

### API Handler

```go
// file: internal/interfaces/http/entity_handler.go
// description: HTTP handlers for entity endpoints
package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type EntityHandler struct {
    createHandler *commands.CreateEntityHandler
    getHandler    *queries.GetEntityHandler
    logger        *zap.Logger
}

func NewEntityHandler(
    createHandler *commands.CreateEntityHandler,
    getHandler *queries.GetEntityHandler,
    logger *zap.Logger,
) *EntityHandler {
    return &EntityHandler{
        createHandler: createHandler,
        getHandler:    getHandler,
        logger:        logger,
    }
}

func (h *EntityHandler) Register(r *gin.RouterGroup) {
    r.POST("/entities", h.Create)
    r.GET("/entities/:id", h.Get)
}

func (h *EntityHandler) Create(c *gin.Context) {
    var req CreateEntityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }

    cmd := commands.CreateEntityCommand{
        Name: req.Name,
    }

    entity, err := h.createHandler.Handle(c.Request.Context(), cmd)
    if err != nil {
        h.logger.Error("create entity failed", zap.Error(err))
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
        return
    }

    c.JSON(http.StatusCreated, EntityResponse{
        ID:        entity.ID.String(),
        Name:      entity.Name,
        Status:    entity.Status.String(),
        CreatedAt: entity.CreatedAt,
    })
}

func (h *EntityHandler) Get(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
        return
    }

    query := queries.GetEntityQuery{ID: id}
    entity, err := h.getHandler.Handle(c.Request.Context(), query)

    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            c.JSON(http.StatusNotFound, ErrorResponse{Error: "entity not found"})
            return
        }
        h.logger.Error("get entity failed", zap.Error(err))
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
        return
    }

    c.JSON(http.StatusOK, EntityResponse{
        ID:        entity.ID.String(),
        Name:      entity.Name,
        Status:    entity.Status.String(),
        CreatedAt: entity.CreatedAt,
    })
}
```

### Docker Configuration

```dockerfile
# file: Dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary
COPY --from=builder /app/api .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["./api"]
```

### Kubernetes Deployment

```yaml
# file: deployments/kubernetes/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
  namespace: production
  labels:
    app: my-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-service
  template:
    metadata:
      labels:
        app: my-service
    spec:
      containers:
        - name: api
          image: my-service:v1.0.0
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: db-credentials
                  key: url
            - name: REDIS_URL
              valueFrom:
                secretKeyRef:
                  name: redis-credentials
                  key: url
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

---

## Operational Considerations

### Deployment Strategy

**Strategy**: Rolling Update with Blue-Green Deployment

**Rollback Plan**:

1. Monitor error rates and latency
2. Trigger automatic rollback if error rate > 1%
3. Manual rollback available via `kubectl rollout undo`

### Monitoring

**Key Metrics**:

- Request rate, latency, errors (Golden Signals)
- Business metrics: [TODO]
- Infrastructure metrics: CPU, memory, disk

**Alerting Rules**:

- High error rate (> 1% for 5 minutes)
- High latency (p99 > 200ms for 5 minutes)
- Resource saturation (CPU > 80% for 10 minutes)

### Security

- Authentication: JWT tokens
- Authorization: RBAC
- Encryption: TLS 1.3 in transit, AES-256 at rest
- Secrets: Kubernetes Secrets + External Secrets Operator

---

## Cross-References

### Prerequisites

- [TODO: [Microservices](../03-Engineering-CloudNative/EC-001-Microservices.md)]
- [TODO: [DDD](../03-Engineering-CloudNative/EC-XXX-DDD.md)]
- [TODO: [CQRS](../03-Engineering-CloudNative/EC-014-CQRS-Pattern.md)]

### Related Application Documents

- [TODO: [AD-XXX: Related](../05-Application-Domains/AD-XXX-Related.md)]

### Other Dimensions

- **Formal Theory**: [TODO: [FT-XXX](../01-Formal-Theory/FT-XXX-Name.md)]
- **Language Design**: [TODO: [LD-XXX](../02-Language-Design/LD-XXX-Name.md)]
- **Engineering**: [TODO: [EC-XXX](../03-Engineering-CloudNative/EC-XXX-Name.md)]
- **Technology**: [TODO: [TS-XXX](../04-Technology-Stack/TS-XXX-Name.md)]

---

## References

### Books

[1] [TODO: Book Title] - [TODO: Author]

### Articles

[2] [TODO: Article Title](https://) - [TODO: Source]

### Case Studies

[3] [TODO: Company] - [TODO: How they implemented similar]

---

## Document History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | [TODO: YYYY-MM-DD] | Initial architecture document | [TODO: Name] |

---

*Template: AD-XXX - Application Domain Document (S/A-Level)*
*For contribution guidelines, see [CONTRIBUTING.md](../CONTRIBUTING.md)*
