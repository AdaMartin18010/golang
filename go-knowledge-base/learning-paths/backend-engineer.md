# Backend Engineer Learning Path

> **Version**: 1.0.0
> **Last Updated**: 2026-04-02
> **Duration**: 16 weeks (full-time) / 24 weeks (part-time)
> **Prerequisites**: Programming fundamentals, basic Go syntax
> **Outcome**: Production-ready backend engineer with microservices expertise

---

## 🎯 Path Overview

### Target Competencies

Upon completion, you will be able to:

- Design and implement scalable RESTful and gRPC APIs
- Build resilient services with circuit breakers, retries, and rate limiting
- Implement secure authentication and authorization
- Design efficient database schemas and queries
- Deploy and monitor services in production
- Debug performance issues and optimize bottlenecks

### Prerequisites Graph

```
Programming Fundamentals
    ↓
Basic Go Syntax (Tour of Go)
    ↓
┌─────────────────────────────────────────────────────────────┐
│              BACKEND ENGINEER LEARNING PATH                  │
│                                                              │
│  Phase 1: Go Mastery (Weeks 1-4)                            │
│    ├── Go Memory Model → Context → Concurrency              │
│    └── Outcome: Idiomatic Go code                           │
│                                                              │
│  Phase 2: API Development (Weeks 5-8)                       │
│    ├── REST → gRPC → Middleware → Validation                │
│    └── Outcome: Production APIs                             │
│                                                              │
│  Phase 3: Data Layer (Weeks 9-12)                           │
│    ├── PostgreSQL → Redis → ORM → Caching                   │
│    └── Outcome: Efficient data access                       │
│                                                              │
│  Phase 4: Production Systems (Weeks 13-16)                  │
│    ├── Resilience → Observability → Testing → Deployment    │
│    └── Outcome: Production-ready services                   │
└─────────────────────────────────────────────────────────────┘
    ↓
Advanced Paths (Optional)
    ├── Cloud-Native Engineer
    ├── Distributed Systems Engineer
    └── Go Specialist
```

---

## 📚 Phase 1: Go Mastery (Weeks 1-4)

### Week 1: Deep Dive into Go Memory and Concurrency

**Goal**: Understand Go's memory model and basic concurrency primitives

#### Day 1-2: Memory Model

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [LD-001] Go Memory Model Formal | 4h | Happens-before relationships, safe concurrency |
| [LD-001] Go Memory Model Happens-Before | 3h | Synchronization guarantees |

**Study Notes**:

- Understand the happens-before relation: `a < b` means `a` is visible to `b`
- Channel communication establishes happens-before
- `sync.Mutex` and `sync.WaitGroup` synchronization
- Safe publication of shared data

#### Day 3-4: Goroutines and Channels

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 02-Language-Design/02-Language-Features/03-Goroutines.md | 3h | Goroutine lifecycle, scheduling |
| 02-Language-Design/02-Language-Features/04-Channels.md | 4h | Buffered vs unbuffered, patterns |
| 01-Formal-Theory/03-Concurrency-Models/01-CSP-Theory.md | 2h | CSP foundations |

**Study Notes**:

- Goroutines are lightweight (2KB stack)
- Channels are for communication AND synchronization
- "Share memory by communicating, don't communicate by sharing memory"
- Common patterns: fan-out, fan-in, pipeline

#### Day 5-7: Context Package Deep Dive

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-005] Context Management | 3h | Context usage patterns |
| 04-Technology-Stack/01-Core-Library/04-Context-Package.md | 2h | Standard library context |
| [EC-011] Context Cancellation Patterns | 3h | Graceful shutdown |

**Study Notes**:

- Context is the first parameter, always
- Propagate cancellation through the call chain
- Use `context.WithTimeout` for external calls
- Never store context in structs

**Week 1 Capstone**:

```go
// Build a concurrent URL fetcher with:
// - Configurable worker pool
// - Context cancellation
// - Timeout handling
// - Rate limiting
```

### Week 2: Advanced Concurrency and Sync Primitives

**Goal**: Master advanced concurrency patterns and synchronization

#### Day 1-3: sync Package Internals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/01-Core-Library/05-Sync-Package.md | 4h | Mutex, RWMutex, Once, Pool |
| [LD-030] Go sync Package Internals | 4h | Implementation details |
| 03-Engineering-CloudNative/03-Performance/06-Lock-Free-Programming.md | 2h | Lock-free patterns |

**Study Notes**:

- `sync.Pool` for reducing GC pressure
- `sync.Map` vs regular map with RWMutex
- `sync.Once` for lazy initialization
- Avoiding deadlocks: consistent lock ordering

#### Day 4-5: Select Statement and Patterns

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 02-Language-Design/02-Language-Features/12-Select-Statement.md | 3h | Non-blocking operations |
| [EC-013] Concurrent Patterns | 4h | Worker pools, pipelines |

**Study Notes**:

- `select` for multiplexing channels
- Default case for non-blocking receive
- `done` channel pattern for cancellation
- Pipeline pattern for data transformation

#### Day 6-7: Error Handling and Panic Recovery

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [LD-008] Go Error Handling Patterns | 3h | Idiomatic error handling |
| 02-Language-Design/02-Language-Features/11-Defer-Panic-Recover.md | 2h | Exception handling |
| 03-Engineering-CloudNative/01-Methodology/06-Error-Handling-Patterns.md | 2h | Application errors |

**Study Notes**:

- Errors are values - check them explicitly
- Custom error types with context
- Panic only for unrecoverable errors
- Recover in goroutine boundaries

**Week 2 Capstone**:

```go
// Build a worker pool with:
// - Dynamic scaling
// - Error propagation
// - Graceful shutdown
// - Metrics collection
```

### Week 3: Go Runtime and Performance Basics

**Goal**: Understand Go runtime fundamentals and basic profiling

#### Day 1-3: Runtime and Scheduler

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [LD-004] Go Runtime GMP Deep Dive | 4h | GMP model explained |
| [FT-002] GMP Scheduler Deep Dive | 3h | Formal understanding |

**Study Notes**:

- G: Goroutine, M: OS Thread, P: Processor
- Work-stealing scheduler
- `GOMAXPROCS` and parallelism
- Syscalls and scheduler interaction

#### Day 4-5: Garbage Collection

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [LD-003] Go Garbage Collector Formal | 3h | GC theory |
| [LD-003] Tri-Color Mark-Sweep | 2h | Algorithm details |
| 03-Engineering-CloudNative/03-Performance/05-Memory-Leak-Detection.md | 2h | Debugging leaks |

**Study Notes**:

- Tri-color concurrent mark-sweep
- GC pacer and heap growth
- Reducing allocations: object pooling
- `runtime.ReadMemStats` for monitoring

#### Day 6-7: Profiling Introduction

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/03-Performance/01-Profiling.md | 4h | pprof usage |
| 03-Engineering-CloudNative/03-Performance/04-Race-Detection.md | 2h | Data races |

**Study Notes**:

- CPU profiling: `go tool pprof cpu.prof`
- Memory profiling: heap and allocs
- Race detector: `-race` flag
- Profile-guided optimization

**Week 3 Capstone**:

```go
// Optimize a slow service:
// - Profile to find bottlenecks
// - Reduce allocations
// - Fix race conditions
// - Measure improvement
```

### Week 4: Testing and Project Structure

**Goal**: Write comprehensive tests and organize projects properly

#### Day 1-3: Testing Fundamentals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/01-Core-Library/08-Testing-Package.md | 3h | Unit testing |
| [LD-009] Go Testing Patterns | 4h | Table-driven tests, subtests |
| 03-Engineering-CloudNative/01-Methodology/03-Testing-Strategies.md | 3h | Test strategy |

**Study Notes**:

- Table-driven tests for coverage
- `testing/quick` for property testing
- `testify/assert` vs standard library
- Test fixtures and helpers

#### Day 4-5: Advanced Testing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-037] Task Testing Strategies | 3h | Integration tests |
| [EC-095] Testing Strategies | 2h | E2E testing |
| 04-Technology-Stack/04-Development-Tools/10-Go-Fuzzing.md | 2h | Fuzz testing |

**Study Notes**:

- Test containers for integration tests
- HTTP test recorder
- Mocking with `gomock` or `mockery`
- Fuzzing for security

#### Day 6-7: Project Structure

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/01-Methodology/05-Project-Structure.md | 3h | Standard layout |
| 04-Technology-Stack/04-Development-Tools/01-Go-Modules.md | 2h | Module management |
| 03-Engineering-CloudNative/01-Methodology/01-Clean-Code.md | 2h | Code quality |

**Study Notes**:

- Standard Go Project Layout
- Internal vs pkg vs cmd
- Semantic versioning with modules
- Code review checklist

**Week 4 Capstone**:

```go
// Build a well-structured module:
// - Clean project layout
// - >80% test coverage
// - Integration tests
// - Makefile for tasks
```

---

## 📚 Phase 2: API Development (Weeks 5-8)

### Week 5: RESTful API Design

**Goal**: Design and implement production-quality REST APIs

#### Day 1-3: API Design Principles

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-004] API Design | 4h | REST principles |
| 05-Application-Domains/01-Backend-Development/01-RESTful-API.md | 3h | HTTP semantics |
| 05-Application-Domains/01-Backend-Development/11-API-Versioning.md | 2h | Versioning strategies |

**Study Notes**:

- Resource-oriented URLs
- Proper HTTP methods and status codes
- Content negotiation
- API versioning: URL vs header vs media type

#### Day 4-5: Web Frameworks

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/03-Network/01-Gin-Framework.md | 4h | Gin framework |
| 04-Technology-Stack/03-Network/03-Echo-Framework.md | 2h | Echo framework |

**Study Notes**:

- Gin: middleware chain, binding, validation
- Router groups for API versions
- Custom middleware
- Performance considerations

#### Day 6-7: Request Handling

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/13-Request-Validation.md | 3h | Input validation |
| 05-Application-Domains/01-Backend-Development/14-Content-Negotiation.md | 2h | Content types |
| [EC-043] Task API Design | 2h | Practical design |

**Study Notes**:

- Struct tags for validation
- Custom validators
- Request/response DTOs
- Content negotiation middleware

**Week 5 Capstone**:

```go
// Build a task management API:
// - RESTful endpoints
// - Request validation
// - Error responses (RFC 7807)
// - API documentation
```

### Week 6: Middleware and Security

**Goal**: Implement secure middleware chains

#### Day 1-3: Middleware Patterns

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/03-Middleware-Patterns.md | 4h | Middleware design |
| [EC-005] Context Management | 3h | Context propagation |

**Study Notes**:

- Middleware as onion layers
- Context for request-scoped values
- Middleware ordering matters
- Recover middleware for panics

#### Day 4-5: Authentication

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/02-Authentication.md | 4h | Auth patterns |
| 03-Engineering-CloudNative/04-Security/03-Cryptography.md | 3h | Crypto basics |

**Study Notes**:

- JWT: structure, signing, validation
- OAuth 2.0 flows
- Password hashing (bcrypt/argon2)
- Session management

#### Day 6-7: Security Best Practices

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/04-Security/05-OWASP-Top-10.md | 3h | Common vulnerabilities |
| 03-Engineering-CloudNative/04-Security/08-Security-Headers.md | 2h | HTTP security |
| 05-Application-Domains/01-Backend-Development/10-Webhook-Security.md | 2h | Webhook security |

**Study Notes**:

- Input sanitization
- CSRF protection
- Security headers (HSTS, CSP)
- Rate limiting for auth endpoints

**Week 6 Capstone**:

```go
// Add to task API:
// - JWT authentication middleware
// - Role-based access control
// - Rate limiting
// - Security headers
```

### Week 7: gRPC and Protocol Buffers

**Goal**: Build high-performance RPC services

#### Day 1-3: Protocol Buffers

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/03-Network/11-Protocol-Buffers.md | 4h | Protobuf basics |
| 04-Technology-Stack/03-Network/02-gRPC.md | 4h | gRPC framework |

**Study Notes**:

- Protobuf syntax and types
- Message evolution rules
- Code generation
- gRPC vs REST trade-offs

#### Day 4-5: gRPC Implementation

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-014] gRPC Internals | 3h | Implementation |
| 04-Technology-Stack/03-Network/13-API-Documentation.md | 2h | Documentation |

**Study Notes**:

- Unary vs streaming
- Interceptors (middleware)
- Deadline and cancellation
- Load balancing

#### Day 6-7: API Gateway

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/04-API-Gateway.md | 3h | Gateway patterns |
| [AD-006] API Gateway Design | 3h | Design principles |

**Study Notes**:

- Gateway responsibilities
- Request routing
- Aggregation pattern
- gRPC-Web translation

**Week 7 Capstone**:

```go
// Build a gRPC service:
// - Proto definitions
// - Server implementation
// - Client with interceptors
// - Gateway for REST conversion
```

### Week 8: Real-Time and GraphQL

**Goal**: Support real-time and flexible query APIs

#### Day 1-3: WebSocket

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/03-Network/04-WebSocket.md | 4h | WebSocket protocol |
| 05-Application-Domains/01-Backend-Development/08-Real-Time-Communication.md | 3h | Implementation |

**Study Notes**:

- WebSocket upgrade handshake
- Connection management
- Broadcast patterns
- Fallback to SSE/long-polling

#### Day 4-5: GraphQL

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/05-GraphQL.md | 4h | GraphQL basics |

**Study Notes**:

- Schema definition
- Resolvers
- N+1 problem and DataLoader
- Subscriptions

#### Day 6-7: Integration

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-049] Integration Patterns | 3h | System integration |

**Week 8 Capstone**:

```go
// Enhance task API:
// - WebSocket for real-time updates
// - GraphQL endpoint
// - Unified API gateway
```

---

## 📚 Phase 3: Data Layer (Weeks 9-12)

### Week 9: PostgreSQL Deep Dive

**Goal**: Master PostgreSQL for production use

#### Day 1-3: PostgreSQL Internals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-001] PostgreSQL Transaction Internals | 4h | MVCC, isolation |
| [TS-001] PostgreSQL Transaction Formal | 3h | Formal semantics |

**Study Notes**:

- MVCC implementation
- Isolation levels and anomalies
- Vacuum and bloat
- WAL and durability

#### Day 4-5: Query Optimization

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/02-Database/13-Database-Pooling.md | 3h | Connection pooling |
| EC-005 | Database Patterns | 3h | Access patterns |

**Study Notes**:

- `EXPLAIN ANALYZE` interpretation
- Index types and selection
- Query planning
- Connection pool sizing

#### Day 6-7: Advanced PostgreSQL

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-065] Transaction Isolation & MVCC | 3h | Deep dive |
| 04-Technology-Stack/02-Database/09-Database-Migration.md | 2h | Migrations |

**Study Notes**:

- Advisory locks
- Partial indexes
- Partitioning strategies
- Migration best practices

**Week 9 Capstone**:

```sql
-- Design schema for:
-- - E-commerce order system
-- - Proper indexing
-- - Partitioning by date
-- - Optimized queries
```

### Week 10: Redis and Caching

**Goal**: Implement effective caching strategies

#### Day 1-3: Redis Internals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-002] Redis Multithreaded IO | 3h | IO model |
| [TS-006] Redis Data Structures | 3h | Data types |
| [TS-007] Redis Data Structures Deep Dive | 2h | Internals |

**Study Notes**:

- Redis IO threads
- Data structure implementations
- Memory encoding
- Persistence options

#### Day 4-5: Caching Strategies

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/02-Database/11-Caching-Strategies.md | 4h | Cache patterns |
| 04-Technology-Stack/02-Database/04-Redis.md | 3h | Redis usage |

**Study Notes**:

- Cache-aside vs write-through
- Cache invalidation strategies
- TTL and eviction policies
- Cache stampede prevention

#### Day 6-7: Advanced Redis

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/02-Database/12-Database-Replication.md | 3h | Replication |
| [TS-003] Redis Internals Formal | 2h | Formal model |

**Study Notes**:

- Redis Sentinel
- Redis Cluster
- Distributed locks (Redlock)
- Stream data type

**Week 10 Capstone**:

```go
// Implement caching layer:
// - Multi-level cache (L1/L2)
// - Cache warming
// - Circuit breaker for Redis
// - Metrics and eviction stats
```

### Week 11: ORM and Database Access

**Goal**: Efficient database access patterns

#### Day 1-3: GORM

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/02-Database/02-ORM-GORM.md | 4h | GORM basics |
| 04-Technology-Stack/02-Database/03-SQLC.md | 3h | SQLC type-safe |

**Study Notes**:

- GORM hooks and callbacks
- Preloading associations
- Raw SQL when needed
- SQLC for compile-time safety

#### Day 4-5: Database Patterns

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-005] Database Patterns | 4h | Repository pattern |
| [AD-001] DDD Strategic Patterns | 3h | Aggregate roots |

**Study Notes**:

- Repository pattern
- Unit of Work
- Specification pattern
- CQRS introduction

#### Day 6-7: Sharding and Scaling

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 04-Technology-Stack/02-Database/10-Database-Sharding.md | 4h | Sharding strategies |
| [EC-082] Distributed Task Sharding | 2h | Implementation |

**Study Notes**:

- Horizontal vs vertical sharding
- Consistent hashing for sharding
- Cross-shard queries
- Shard rebalancing

**Week 11 Capstone**:

```go
// Build data access layer:
// - Repository pattern
// - Unit of Work
// - Multi-tenant data isolation
// - Query optimization
```

### Week 12: Distributed Transactions

**Goal**: Handle transactions across services

#### Day 1-3: Transaction Theory

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/07-Distributed-Transactions.md | 4h | Distributed TX |
| [EC-008] Saga Pattern | 4h | Saga implementation |

**Study Notes**:

- 2PC limitations
- Saga pattern: choreography vs orchestration
- Compensating transactions
- Outbox pattern

#### Day 4-5: Saga Implementation

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-090] Saga Implementation | 3h | Code patterns |
| [EC-112] Saga Pattern Complete | 3h | Full example |
| ../examples/saga/ | 2h | Working example |

**Study Notes**:

- Saga orchestrator
- Compensation logic
- Idempotency keys
- State machine for sagas

#### Day 6-7: Idempotency

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/01-Backend-Development/09-Idempotency.md | 3h | Idempotency patterns |
| [EC-013] Idempotency Pattern | 3h | Implementation |
| [EC-119] Idempotency Guarantee | 2h | Production ready |

**Study Notes**:

- Idempotency key generation
- Storage strategies
- TTL and cleanup
- Duplicate detection

**Week 12 Capstone**:

```go
// Implement order processing:
// - Saga for checkout flow
// - Idempotent payment processing
// - Outbox pattern
// - Event publishing
```

---

## 📚 Phase 4: Production Systems (Weeks 13-16)

### Week 13: Resilience Patterns

**Goal**: Build fault-tolerant services

#### Day 1-3: Circuit Breaker

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-007] Circuit Breaker | 4h | Pattern basics |
| [EC-008] Circuit Breaker Advanced | 3h | Production |
| [EC-117] Circuit Breaker Advanced | 3h | Deep dive |

**Study Notes**:

- States: Closed, Open, Half-Open
- Failure threshold configuration
- Success threshold for recovery
- Bulkhead pattern combination

#### Day 4-5: Retry and Timeout

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-009] Retry Pattern | 3h | Retry strategies |
| [EC-075] Retry & Backoff | 3h | Exponential backoff |
| [EC-083] Timeout Control | 2h | Timeout patterns |

**Study Notes**:

- Exponential backoff with jitter
- Context deadline propagation
- Idempotency for retries
- Max retry limits

#### Day 6-7: Rate Limiting

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-012] Rate Limiting | 4h | Algorithms |
| [EC-030] Rate Limiting | 3h | Implementation |
| [EC-078] Rate Limiting & Throttling | 3h | Advanced |

**Study Notes**:

- Token bucket vs leaky bucket
- Distributed rate limiting
- Sliding window counter
- Per-user vs global limits

**Week 13 Capstone**:

```go
// Build resilient HTTP client:
// - Circuit breaker
// - Retry with backoff
// - Timeout handling
// - Rate limiting
```

### Week 14: Observability

**Goal**: Implement comprehensive observability

#### Day 1-3: Distributed Tracing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-006] Distributed Tracing | 4h | Tracing basics |
| [EC-056] Distributed Tracing Deep Dive | 4h | Implementation |
| [EC-070] W3C Trace Context | 2h | Standards |

**Study Notes**:

- Trace, span, span context
- Context propagation
- Sampling strategies
- W3C Trace Context headers

#### Day 4-5: Metrics and Logging

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-013] Prometheus | 4h | Metrics |
| [EC-074] Context-Aware Logging | 3h | Structured logs |
| EC-022 | Context-Aware Logging | 2h | Correlation |

**Study Notes**:

- RED metrics (Rate, Errors, Duration)
- USE metrics (Utilization, Saturation, Errors)
- Structured JSON logging
- Log correlation with trace IDs

#### Day 6-7: Production Observability

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-044] Observability Production | 4h | Full stack |
| [EC-080] Metrics Integration | 3h | Dashboards |

**Study Notes**:

- OpenTelemetry integration
- Grafana dashboards
- Alerting rules
- SLO/SLI definitions

**Week 14 Capstone**:

```go
// Add observability to task API:
// - OpenTelemetry tracing
// - Prometheus metrics
// - Structured logging
// - Grafana dashboards
```

### Week 15: Health Checks and Graceful Shutdown

**Goal**: Implement production-ready lifecycle management

#### Day 1-3: Health Checks

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-014] Health Checks | 4h | Health patterns |
| [EC-086] Health Check Patterns | 3h | Production |

**Study Notes**:

- Liveness vs readiness probes
- Deep health checks
- Dependency health
- Kubernetes probe configuration

#### Day 4-5: Graceful Shutdown

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-007] Graceful Shutdown Complete | 4h | Shutdown patterns |
| [EC-079] Graceful Shutdown | 3h | Implementation |
| [EC-120] Graceful Shutdown Complete | 3h | Production |

**Study Notes**:

- Signal handling
- Draining in-flight requests
- Closing connections
- Cleanup tasks

#### Day 6-7: Backpressure

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-118] Backpressure & Flow Control | 4h | Flow control |

**Study Notes**:

- Load shedding
- Request queuing
- Adaptive concurrency
- Client-side backoff

**Week 15 Capstone**:

```go
// Implement lifecycle management:
// - Kubernetes health probes
// - Graceful shutdown
// - Backpressure handling
// - Zero-downtime deployment
```

### Week 16: Deployment and DevOps

**Goal**: Deploy and operate services in production

#### Day 1-3: Containers and Kubernetes

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-003] Container Design | 4h | Container patterns |
| [TS-005] Kubernetes Operators | 4h | K8s basics |

**Study Notes**:

- Multi-stage Docker builds
- Security contexts
- Resource limits
- Pod disruption budgets

#### Day 4-5: CI/CD

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/03-DevOps-Tools/04-CI-CD.md | 4h | Pipeline design |
| 05-Application-Domains/02-Cloud-Infrastructure/08-GitOps.md | 3h | GitOps |

**Study Notes**:

- GitHub Actions workflows
- Build caching
- Automated testing
- GitOps with ArgoCD

#### Day 6-7: SRE Practices

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-121] Google SRE | 4h | Reliability |
| 05-Application-Domains/03-DevOps-Tools/12-Platform-Engineering.md | 3h | Platform |

**Study Notes**:

- Error budgets
- Incident response
- Postmortems
- Toil reduction

**Week 16 Capstone**:

```yaml
# Deploy task API:
# - Kubernetes manifests
# - Helm chart
# - CI/CD pipeline
# - Monitoring stack
```

---

## 🎓 Capstone Project: Complete Backend System

### Project: E-Commerce Order Management System

**Requirements**:

1. **API Layer**
   - REST API for order management
   - gRPC for internal services
   - GraphQL for flexible queries
   - WebSocket for real-time order updates

2. **Authentication & Security**
   - JWT-based authentication
   - Role-based access control (customer, admin)
   - Rate limiting
   - Security headers

3. **Data Layer**
   - PostgreSQL for transactional data
   - Redis for caching and sessions
   - Database migrations
   - Connection pooling

4. **Distributed Transactions**
   - Saga pattern for order processing
   - Idempotency for payment
   - Outbox pattern for events

5. **Resilience**
   - Circuit breaker for external calls
   - Retry with exponential backoff
   - Timeout handling
   - Graceful degradation

6. **Observability**
   - Distributed tracing with OpenTelemetry
   - Prometheus metrics
   - Structured logging
   - Health checks

7. **Deployment**
   - Docker containers
   - Kubernetes deployment
   - CI/CD pipeline
   - GitOps workflow

**Evaluation Criteria**:

- [ ] Code quality and organization
- [ ] Test coverage >80%
- [ ] API documentation
- [ ] Performance benchmarks
- [ ] Security audit pass
- [ ] Deployment runbook

---

## 📖 Supplementary Resources

### Books

- "Building Microservices" by Sam Newman
- "Designing Data-Intensive Applications" by Martin Kleppmann
- "The Go Programming Language" by Donovan & Kernighan
- "Cloud Native Go" by Matthew Titmus

### Online Courses

- Coursera: Cloud Computing Specialization
- Udemy: Go: The Complete Developer's Guide
- Pluralsight: Building Scalable APIs with Go

### Communities

- Go Forum: forum.golangbridge.org
- Gopher Slack: invite.gopherjs.io
- Reddit: r/golang

---

## ✅ Progress Tracker

| Phase | Week | Topic | Complete |
|-------|------|-------|----------|
| 1 | 1 | Memory & Concurrency | [ ] |
| 1 | 2 | Advanced Concurrency | [ ] |
| 1 | 3 | Runtime & Performance | [ ] |
| 1 | 4 | Testing & Structure | [ ] |
| 2 | 5 | RESTful API Design | [ ] |
| 2 | 6 | Middleware & Security | [ ] |
| 2 | 7 | gRPC & Protocol Buffers | [ ] |
| 2 | 8 | Real-Time & GraphQL | [ ] |
| 3 | 9 | PostgreSQL | [ ] |
| 3 | 10 | Redis & Caching | [ ] |
| 3 | 11 | ORM & Patterns | [ ] |
| 3 | 12 | Distributed Transactions | [ ] |
| 4 | 13 | Resilience Patterns | [ ] |
| 4 | 14 | Observability | [ ] |
| 4 | 15 | Health & Shutdown | [ ] |
| 4 | 16 | Deployment & SRE | [ ] |

---

*This learning path provides a comprehensive foundation for backend engineering with Go. Complete all phases and the capstone project to achieve production-ready competency.*
