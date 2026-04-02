# Glossary

> **Version**: 1.0 S-Level  
> **Created**: 2026-04-02  
> **Status**: Active  
> **Scope**: Go Knowledge Base Terminology

---

## Table of Contents

1. [How to Use This Glossary](#how-to-use-this-glossary)
2. [A](#a)
3. [B](#b)
4. [C](#c)
5. [D](#d)
6. [E](#e)
7. [F](#f)
8. [G](#g)
9. [H](#h)
10. [I](#i)
11. [J](#j)
12. [K](#k)
13. [L](#l)
14. [M](#m)
15. [N](#n)
16. [O](#o)
17. [P](#p)
18. [Q](#q)
19. [R](#r)
20. [S](#s)
21. [T](#t)
22. [U](#u)
23. [V](#v)
24. [W](#w)
25. [X](#x)
26. [Y](#y)
27. [Z](#z)

---

## How to Use This Glossary

### Definition Format

Each entry follows this structure:

```
Term
├── Category: [Technical/Process/Concept]
├── Related: [Related terms]
├── See: [Link to detailed documentation]
└── Definition: [Precise definition]
```

### Cross-Reference Symbols

| Symbol | Meaning |
|--------|---------|
| → | Synonym or preferred term |
| See | Related concept |
| Also | Additional related terms |
| Cf. | Compare with (contrast) |

---

## A

### ACID
**Category**: Database  
**Related**: BASE, Transaction, Consistency  
**See**: [FT-009-Distributed-Transactions-Formal.md](./01-Formal-Theory/FT-009-Distributed-Transactions-Formal.md)

Atomicity, Consistency, Isolation, Durability — properties guaranteeing reliable database transactions.
- **Atomicity**: All operations succeed or all fail
- **Consistency**: Database remains in valid state
- **Isolation**: Concurrent transactions don't interfere
- **Durability**: Committed transactions persist

### Actor Model
**Category**: Concurrency  
**Related**: CSP, Goroutine, Message Passing  
**See**: [01-Formal-Theory/03-Concurrency-Models/](./01-Formal-Theory/03-Concurrency-Models/)

Concurrency model where actors are units of computation that communicate exclusively through asynchronous message passing.

### Axiomatic Semantics
**Category**: Formal Methods  
**Related**: Operational Semantics, Denotational Semantics  
**See**: [01-Formal-Theory/01-Semantics/03-Axiomatic-Semantics.md](./01-Formal-Theory/01-Semantics/03-Axiomatic-Semantics.md)

Approach to defining program semantics through logical assertions about program states.

---

## B

### BASE
**Category**: Database  
**Related**: ACID, Eventual Consistency, NoSQL  
**See**: [FT-004-Distributed-Systems-Fundamentals-CAP-BASE-ACID.md](./01-Formal-Theory/FT-004-Distributed-Systems-Fundamentals-CAP-BASE-ACID.md)

Basically Available, Soft state, Eventually consistent — properties of distributed databases prioritizing availability.

### Bloom Filter
**Category**: Data Structure  
**Related**: Probabilistic Data Structures, Hashing  
**See**: [FT-008-Probabilistic-Data-Structures.md](./01-Formal-Theory/FT-008-Probabilistic-Data-Structures.md)

Space-efficient probabilistic data structure for set membership testing with possible false positives.

### Byzantine Fault
**Category**: Distributed Systems  
**Related**: Byzantine Fault Tolerance, Consensus  
**See**: [FT-007-Byzantine-Fault-Tolerance.md](./01-Formal-Theory/FT-007-Byzantine-Fault-Tolerance.md)

Arbitrary node failure where faulty nodes may behave maliciously or inconsistently.

### Byzantine Fault Tolerance (BFT)
**Category**: Distributed Systems  
**Related**: Consensus, PBFT, Fault Tolerance  
**See**: [FT-008-Byzantine-Consensus-Formal.md](./01-Formal-Theory/FT-008-Byzantine-Consensus-Formal.md)

System's ability to withstand Byzantine faults, requiring ≥3f+1 nodes to tolerate f faults.

---

## C

### CAP Theorem
**Category**: Distributed Systems  
**Related**: Consistency, Availability, Partition Tolerance  
**See**: [FT-001-Distributed-Systems-Foundation-Formal.md](./01-Formal-Theory/FT-001-Distributed-Systems-Foundation-Formal.md)

Theorem stating distributed systems can guarantee at most two of: Consistency, Availability, Partition tolerance.

### Channel (Go)
**Category**: Go Language  
**Related**: Goroutine, CSP, Select  
**See**: [02-Language-Design/02-Language-Features/04-Channels.md](./02-Language-Design/02-Language-Features/04-Channels.md)

Typed conduit for communication between goroutines; implements CSP communication model.

```go
ch := make(chan int)    // Unbuffered
ch := make(chan int, 5) // Buffered
```

### Circuit Breaker
**Category**: Resilience Pattern  
**Related**: Retry, Timeout, Bulkhead  
**See**: [EC-001-Circuit-Breaker-Pattern.md](./03-Engineering-CloudNative/EC-001-Circuit-Breaker-Pattern.md)

Pattern preventing cascade failures by stopping requests to failing services.

States: **CLOSED** (normal) → **OPEN** (failing) → **HALF-OPEN** (testing)

### Consistency Model
**Category**: Distributed Systems  
**Related**: Linearizability, Eventual Consistency, CAP  
**See**: [FT-010-Linearizability-Formal.md](./01-Formal-Theory/FT-010-Linearizability-Formal.md)

Guarantees provided by distributed systems regarding the ordering and visibility of operations.

**Hierarchy**: Linearizability → Sequential → Causal → Eventual

### Context (Go)
**Category**: Go Language  
**Related**: Cancellation, Deadline, Request-scoped Values  
**See**: [04-Technology-Stack/01-Core-Library/04-Context-Package.md](./04-Technology-Stack/01-Core-Library/04-Context-Package.md)

Package for carrying deadlines, cancellation signals, and request-scoped values across API boundaries.

### CRDT
**Category**: Distributed Systems  
**Related**: Eventual Consistency, Gossip Protocols  
**See**: [FT-018-CRDT-Formal.md](./01-Formal-Theory/FT-018-CRDT-Formal.md)

Conflict-free Replicated Data Type — data structure that can be replicated and modified concurrently without coordination.

### CSP
**Category**: Concurrency Model  
**Related**: Actor Model, Channels, Goroutines  
**See**: [01-Formal-Theory/03-Concurrency-Models/01-CSP-Theory.md](./01-Formal-Theory/03-Concurrency-Models/01-CSP-Theory.md)

Communicating Sequential Processes — formal language for describing patterns of interaction in concurrent systems; Go's concurrency model is based on CSP.

---

## D

### Deadlock
**Category**: Concurrency  
**Related**: Livelock, Starvation, Race Condition  
**See**: [03-Engineering-CloudNative/01-Methodology/](./03-Engineering-CloudNative/01-Methodology/)

State where two or more processes block forever, waiting for each other.

**Coffman Conditions** (all necessary):
1. Mutual exclusion
2. Hold and wait
3. No preemption
4. Circular wait

### Defer (Go)
**Category**: Go Language  
**Related**: Panic, Recover, Cleanup  
**See**: [02-Language-Design/02-Language-Features/11-Defer-Panic-Recover.md](./02-Language-Design/02-Language-Features/11-Defer-Panic-Recover.md)

Statement scheduling function call to execute after surrounding function returns.

```go
func example() {
    f := openFile()
    defer f.Close()  // Executes when example() returns
    // ... use f
}
```

### DRF-SC
**Category**: Memory Model  
**Related**: Happens-Before, Memory Model, Data Race  
**See**: [01-Formal-Theory/04-Memory-Models/02-DRF-SC.md](./01-Formal-Theory/04-Memory-Models/02-DRF-SC.md)

Data-Race-Free Sequential Consistency — property guaranteeing sequentially consistent execution for programs without data races.

---

## E

### Embedding (Go)
**Category**: Go Language  
**Related**: Composition, Inheritance, Struct  
**See**: [02-Language-Design/02-Language-Features/13-Struct-Embedding.md](./02-Language-Design/02-Language-Features/13-Struct-Embedding.md)

Go's approach to code reuse where embedded types' methods are promoted to the embedding struct.

### Ephemeral Storage
**Category**: Cloud-Native  
**Related**: Persistent Storage, Stateful, Stateless  
**See**: [03-Engineering-CloudNative/02-Cloud-Native/](./03-Engineering-CloudNative/02-Cloud-Native/)

Temporary storage that exists only for container lifetime; destroyed when container restarts.

### Event Sourcing
**Category**: Architecture Pattern  
**Related**: CQRS, Event-Driven, Audit Trail  
**See**: [EC-015-Event-Sourcing-Pattern.md](./03-Engineering-CloudNative/EC-015-Event-Sourcing-Pattern.md)

Pattern storing application state as a sequence of events rather than current state.

### Eventual Consistency
**Category**: Distributed Systems  
**Related**: Strong Consistency, BASE, CAP  
**See**: [FT-013-Eventual-Consistency-Formal.md](./01-Formal-Theory/FT-013-Eventual-Consistency-Formal.md)

Consistency model guaranteeing that if no new updates are made, eventually all accesses return the last updated value.

---

## F

### F-Bounded Polymorphism
**Category**: Type Theory  
**Related**: Generics, Type Constraints  
**See**: [01-Formal-Theory/02-Type-Theory/03-Generics-Theory/01-F-Bounded-Polymorphism.md](./01-Formal-Theory/02-Type-Theory/03-Generics-Theory/01-F-Bounded-Polymorphism.md)

Recursive type constraint where a type parameter appears in its own bound.

### Featherweight Go
**Category**: Formal Methods  
**Related**: Operational Semantics, Type System  
**See**: [01-Formal-Theory/01-Semantics/04-Featherweight-Go.md](./01-Formal-Theory/01-Semantics/04-Featherweight-Go.md)

Minimal core calculus of Go used for formal reasoning about the type system.

### FLP Impossibility
**Category**: Distributed Systems  
**Related**: Consensus, Asynchronous Systems  
**See**: [FT-015-FLP-Impossibility-Formal.md](./01-Formal-Theory/FT-015-FLP-Impossibility-Formal.md)

Fischer, Lynch, Paterson result proving impossibility of deterministic consensus in asynchronous systems with even one process crash.

---

## G

### Generics (Go)
**Category**: Go Language  
**Related**: Type Parameters, Constraints, Type Sets  
**See**: [02-Language-Design/02-Language-Features/06-Generics.md](./02-Language-Design/02-Language-Features/06-Generics.md)

Feature (Go 1.18+) enabling types and functions to operate with any type satisfying constraints.

```go
func Max[T comparable](a, b T) T {
    if a > b { return a }
    return b
}
```

### Goroutine
**Category**: Go Language  
**Related**: Channel, Scheduler, Concurrency  
**See**: [02-Language-Design/02-Language-Features/03-Goroutines.md](./02-Language-Design/02-Language-Features/03-Goroutines.md)

Lightweight thread managed by Go runtime; starts with 2KB stack growing/shrinking as needed.

```go
go function()  // Start goroutine
```

### Gossip Protocol
**Category**: Distributed Systems  
**Related**: Epidemic Broadcast, Membership, CRDT  
**See**: [FT-011-Gossip-Protocols.md](./01-Formal-Theory/FT-011-Gossip-Protocols.md)

Probabilistic communication protocol where nodes randomly exchange state, propagating information epidemically.

### Graceful Shutdown
**Category**: Operations  
**Related**: Signal Handling, Cleanup, Zero-Downtime  
**See**: [EC-009-Graceful-Shutdown.md](./03-Engineering-CloudNative/EC-009-Graceful-Shutdown.md)

Process of shutting down application allowing in-flight requests to complete before exiting.

---

## H

### Happens-Before
**Category**: Memory Model  
**Related**: Synchronization, Visibility, DRF-SC  
**See**: [01-Formal-Theory/04-Memory-Models/01-Happens-Before.md](./01-Formal-Theory/04-Memory-Models/01-Happens-Before.md)

Partial order relation defining when memory writes are visible to reads; fundamental to Go's memory model.

### Horizontal Scaling
**Category**: Architecture  
**Related**: Vertical Scaling, Load Balancing, Sharding  
**See**: [05-Application-Domains/AD-009-Capacity-Planning.md](./05-Application-Domains/AD-009-Capacity-Planning.md)

Scaling by adding more machines rather than increasing resources on existing machines.

---

## I

### Idempotency
**Category**: API Design  
**Related**: Retry, Safety, HTTP Methods  
**See**: [EC-013-Idempotency-Patterns.md](./03-Engineering-CloudNative/EC-013-Idempotency-Patterns.md)

Property where operation produces the same result whether executed once or multiple times.

### Interface (Go)
**Category**: Go Language  
**Related**: Structural Typing, Duck Typing, Composition  
**See**: [02-Language-Design/02-Language-Features/02-Interfaces.md](./02-Language-Design/02-Language-Features/02-Interfaces.md)

Type defining method signatures; implemented implicitly by types with matching methods (structural typing).

### Invariant
**Category**: Formal Methods  
**Related**: Precondition, Postcondition, Assertion  
**See**: [01-Formal-Theory/03-Program-Verification/](./01-Formal-Theory/03-Program-Verification/)

Condition that remains true before and after an operation.

---

## J

### JWT (JSON Web Token)
**Category**: Security  
**Related**: Authentication, Authorization, OAuth  
**See**: [03-Engineering-CloudNative/04-Security/](./03-Engineering-CloudNative/04-Security/)

Compact, URL-safe token format for claims transfer; structure: `header.payload.signature`

---

## K

### Kubernetes
**Category**: Cloud-Native  
**Related**: Container Orchestration, Pods, Services  
**See**: [04-Technology-Stack/03-Network/09-Service-Mesh.md](./04-Technology-Stack/03-Network/09-Service-Mesh.md)

Container orchestration platform automating deployment, scaling, and management of containerized applications.

---

## L

### Liveness
**Category**: Distributed Systems  
**Related**: Safety, Consensus, Fairness  
**See**: [01-Formal-Theory/03-Concurrency-Models/](./01-Formal-Theory/03-Concurrency-Models/)

Property ensuring something good eventually happens (e.g., consensus is reached).

### Linearizability
**Category**: Consistency  
**Related**: Sequential Consistency, Strong Consistency  
**See**: [FT-010-Linearizability-Formal.md](./01-Formal-Theory/FT-010-Linearizability-Formal.md)

Strongest consistency model where operations appear to occur instantaneously at some point between invocation and completion.

### Liveness Probe
**Category**: Cloud-Native  
**Related**: Readiness, Health Check, Self-Healing  
**See**: [EC-014-Health-Checks.md](./03-Engineering-CloudNative/EC-014-Health-Checks.md)

Kubernetes check determining if container is running (not deadlocked); failure triggers restart.

### Load Balancing
**Category**: Architecture  
**Related**: Horizontal Scaling, Health Checks, Proxies  
**See**: [EC-006-Load-Balancing-Algorithms.md](./03-Engineering-CloudNative/EC-006-Load-Balancing-Algorithms.md)

Distribution of traffic across multiple servers for availability and scalability.

**Algorithms**: Round-robin, Least connections, IP hash, Weighted

---

## M

### Memory Model
**Category**: Language  
**Related**: Happens-Before, Visibility, Ordering  
**See**: [LD-001-Go-Memory-Model-Formal.md](./02-Language-Design/LD-001-Go-Memory-Model-Formal.md)

Specification defining when memory writes are visible to reads in concurrent programs.

### Microservices
**Category**: Architecture  
**Related**: Monolith, Service Mesh, Bounded Context  
**See**: [EC-001-Microservices.md](./03-Engineering-CloudNative/EC-001-Microservices.md)

Architectural style structuring application as loosely coupled, independently deployable services.

### Middleware
**Category**: Web Development  
**Related**: Chain, Handler, Decorator  
**See**: [05-Application-Domains/01-Backend-Development/03-Middleware-Patterns.md](./05-Application-Domains/01-Backend-Development/03-Middleware-Patterns.md)

Software layer processing requests/responses; chainable functions wrapping handlers.

### Monolith
**Category**: Architecture  
**Related**: Microservices, Modular Monolith  
**See**: [EC-001-Microservices.md](./03-Engineering-CloudNative/EC-001-Microservices.md)

Architecture where all components are in a single deployable unit.

### Mutex
**Category**: Concurrency  
**Related**: Lock, RWMutex, sync  
**See**: [04-Technology-Stack/01-Core-Library/05-Sync-Package.md](./04-Technology-Stack/01-Core-Library/05-Sync-Package.md)

Mutual exclusion lock; allows one goroutine at a time to access critical section.

---

## N

### NATS
**Category**: Messaging  
**Related**: Pub/Sub, Message Queue, Microservices  
**See**: [04-Technology-Stack/03-Network/06-NATS.md](./04-Technology-Stack/03-Network/06-NATS.md)

Lightweight, high-performance messaging system for distributed systems.

### Non-Blocking
**Category**: I/O  
**Related**: Async, Event Loop, goroutine  
**See**: [02-Language-Design/02-Language-Features/03-Goroutines.md](./02-Language-Design/02-Language-Features/03-Goroutines.md)

Operation that returns immediately without waiting for completion; caller polls or uses callback.

---

## O

### Observability
**Category**: Operations  
**Related**: Monitoring, Tracing, Logging  
**See**: [EC-044-Observability-Production.md](./03-Engineering-CloudNative/EC-044-Observability-Production.md)

Ability to understand system state from outputs (metrics, logs, traces); "monitoring you can ask questions of."

**Pillars**: Metrics, Logs, Traces, Profiles

### OpenTelemetry
**Category**: Observability  
**Related**: Distributed Tracing, Metrics, OTLP  
**See**: [EC-060-OpenTelemetry-Distributed-Tracing-Production.md](./03-Engineering-CloudNative/EC-060-OpenTelemetry-Distributed-Tracing-Production.md)

Open-source observability framework providing APIs, libraries, and agents for telemetry collection.

### Operational Semantics
**Category**: Formal Methods  
**Related**: Denotational Semantics, Axiomatic Semantics  
**See**: [01-Formal-Theory/01-Semantics/01-Operational-Semantics.md](./01-Formal-Theory/01-Semantics/01-Operational-Semantics.md)

Approach defining program meaning through execution rules specifying how program state evolves.

### Outbox Pattern
**Category**: Architecture  
**Related**: Saga, Event Sourcing, Transaction  
**See**: [EC-013-Outbox-Pattern.md](./03-Engineering-CloudNative/EC-013-Outbox-Pattern.md)

Pattern ensuring reliable message publishing by storing messages in database table within same transaction as business update.

---

## P

### Paxos
**Category**: Consensus  
**Related**: Raft, Multi-Paxos, FLP  
**See**: [FT-003-Paxos-Consensus-Formal.md](./01-Formal-Theory/FT-003-Paxos-Consensus-Formal.md)

Family of consensus protocols for unreliable distributed networks; proven safe but complex.

### Pointer Receiver (Go)
**Category**: Go Language  
**Related**: Value Receiver, Method Set, Mutability  
**See**: [02-Language-Design/LD-005-Go-126-Pointer-Receiver-Constraints.md](./02-Language-Design/LD-005-Go-126-Pointer-Receiver-Constraints.md)

Method receiver type that can modify the receiver and shares the value across calls.

```go
func (p *Type) Method()  // Pointer receiver
```

### Polymorphism
**Category**: Type Theory  
**Related**: Generics, Interfaces, Subtyping  
**See**: [01-Formal-Theory/02-Type-Theory/04-Subtyping.md](./01-Formal-Theory/02-Type-Theory/04-Subtyping.md)

Ability for interface or function to work with multiple types; Go achieves via interfaces and generics.

---

## Q

### Quorum
**Category**: Distributed Systems  
**Related**: Consensus, Majority, RAFT  
**See**: [FT-017-Quorum-Consensus-Formal.md](./01-Formal-Theory/FT-017-Quorum-Consensus-Formal.md)

Minimum number of votes required for an operation; typically ⌊N/2⌋+1 for N nodes.

### Quality Levels
**Category**: Documentation  
**Related**: S-Level, A-Level, Standards  
**See**: [QUALITY-STANDARDS.md](./QUALITY-STANDARDS.md)

Classification system for documentation quality: S (Supreme), A (Advanced), B (Basic), C (Concise).

---

## R

### Race Condition
**Category**: Concurrency  
**Related**: Data Race, Deadlock, Synchronization  
**See**: [03-Engineering-CloudNative/03-Performance/04-Race-Detection.md](./03-Engineering-CloudNative/03-Performance/04-Race-Detection.md)

Bug where behavior depends on relative timing of events; often leads to data races.

### Raft
**Category**: Consensus  
**Related**: Paxos, Leader Election, Log Replication  
**See**: [FT-002-Raft-Consensus-Formal.md](./01-Formal-Theory/FT-002-Raft-Consensus-Formal.md)

Consensus algorithm designed for understandability; uses leader election and log replication.

### Rate Limiting
**Category**: Resilience  
**Related**: Throttling, Circuit Breaker, Token Bucket  
**See**: [EC-012-Rate-Limiting-Formal.md](./03-Engineering-CloudNative/EC-012-Rate-Limiting-Formal.md)

Controlling request rate to prevent overload; algorithms include token bucket, leaky bucket, sliding window.

### Retry Pattern
**Category**: Resilience  
**Related**: Backoff, Circuit Breaker, Idempotency  
**See**: [EC-002-Retry-Pattern.md](./03-Engineering-CloudNative/EC-002-Retry-Pattern.md)

Pattern for transient failure handling with exponential backoff and jitter.

---

## S

### Saga Pattern
**Category**: Distributed Transactions  
**Related**: Compensation, Orchestration, Choreography  
**See**: [EC-008-Saga-Pattern-Formal.md](./03-Engineering-CloudNative/EC-008-Saga-Pattern-Formal.md)

Pattern managing distributed transactions through sequence of local transactions with compensating actions.

### Safety
**Category**: Distributed Systems  
**Related**: Liveness, Consensus, Properties  
**See**: [01-Formal-Theory/03-Concurrency-Models/](./01-Formal-Theory/03-Concurrency-Models/)

Property ensuring something bad never happens (e.g., two different values decided).

### Select (Go)
**Category**: Go Language  
**Related**: Channel, Goroutine, Non-blocking  
**See**: [02-Language-Design/02-Language-Features/12-Select-Statement.md](./02-Language-Design/02-Language-Features/12-Select-Statement.md)

Control structure waiting on multiple channel operations; analogous to switch for channels.

```go
select {
case v1 := <-ch1:
    // Use v1
case v2 := <-ch2:
    // Use v2
case <-timeout:
    // Timeout
}
```

### Service Mesh
**Category**: Cloud-Native  
**Related**: Sidecar, Istio, mTLS  
**See**: [04-Technology-Stack/03-Network/09-Service-Mesh.md](./04-Technology-Stack/03-Network/09-Service-Mesh.md)

Infrastructure layer handling service-to-service communication; provides traffic management, security, observability.

### Sidecar Pattern
**Category**: Architecture  
**Related**: Service Mesh, Container, Proxy  
**See**: [EC-014-Sidecar-Pattern-Formal.md](./03-Engineering-CloudNative/EC-014-Sidecar-Pattern-Formal.md)

Pattern deploying helper container alongside application container in same pod.

### Slice (Go)
**Category**: Go Language  
**Related**: Array, Map, Reference Type  
**See**: [02-Language-Design/02-Language-Features/17-Slice-Internals.md](./02-Language-Design/02-Language-Features/17-Slice-Internals.md)

Reference to contiguous array segment; dynamic-size view into array.

```go
slice := []int{1, 2, 3}
slice = append(slice, 4)
```

### Structural Typing
**Category**: Type Theory  
**Related**: Nominal Typing, Interfaces, Duck Typing  
**See**: [01-Formal-Theory/02-Type-Theory/01-Structural-Typing.md](./01-Formal-Theory/02-Type-Theory/01-Structural-Typing.md)

Type system where type compatibility is determined by structure, not declaration; Go interfaces use structural typing.

---

## T

### Temporal
**Category**: Workflow  
**Related**: Saga, Stateful, Durability  
**See**: [EC-100-Temporal-Workflow-Engine.md](./03-Engineering-CloudNative/EC-100-Temporal-Workflow-Engine.md)

Durable execution platform for reliable workflow orchestration.

### Timeout
**Category**: Resilience  
**Related**: Cancellation, Context, Retry  
**See**: [EC-003-Timeout-Pattern.md](./03-Engineering-CloudNative/EC-003-Timeout-Pattern.md)

Pattern preventing indefinite waiting by specifying maximum wait time.

### Tracing
**Category**: Observability  
**Related**: Span, Trace, Distributed  
**See**: [EC-056-Task-Distributed-Tracing-Deep-Dive.md](./03-Engineering-CloudNative/EC-056-Task-Distributed-Tracing-Deep-Dive.md)

Recording request path through distributed system; tracks latency and dependencies.

### Two-Phase Commit (2PC)
**Category**: Distributed Transactions  
**Related**: Saga, Consensus, Atomic Commit  
**See**: [FT-021-Two-Phase-Commit-Formal.md](./01-Formal-Theory/FT-021-Two-Phase-Commit-Formal.md)

Protocol for atomic transaction commit across distributed systems using prepare and commit phases.

### Type Assertion (Go)
**Category**: Go Language  
**Related**: Interface, Type Switch, Reflection  
**See**: [02-Language-Design/02-Language-Features/20-Type-Assertions.md](./02-Language-Design/02-Language-Features/20-Type-Assertions.md)

Accessing concrete value from interface variable.

```go
v := i.(Type)      // Panics if wrong
v, ok := i.(Type)  // Safe check
```

---

## U

### Unary RPC
**Category**: gRPC  
**Related**: Streaming, Client, Server  
**See**: [04-Technology-Stack/03-Network/02-gRPC.md](./04-Technology-Stack/03-Network/02-gRPC.md)

Simplest RPC pattern: single request, single response.

---

## V

### Vector Clock
**Category**: Distributed Systems  
**Related**: Logical Time, Causality, Version Vector  
**See**: [FT-005-Vector-Clocks-Formal.md](./01-Formal-Theory/FT-005-Vector-Clocks-Formal.md)

Algorithm for generating partial ordering of events and detecting causality violations.

### Vertical Scaling
**Category**: Architecture  
**Related**: Horizontal Scaling, Resources  
**See**: [05-Application-Domains/AD-009-Capacity-Planning.md](./05-Application-Domains/AD-009-Capacity-Planning.md)

Scaling by adding resources (CPU, memory) to existing machine.

### Visibility (Memory)
**Category**: Concurrency  
**Related**: Happens-Before, Memory Model, Flush  
**See**: [01-Formal-Theory/04-Memory-Models/01-Happens-Before.md](./01-Formal-Theory/04-Memory-Models/01-Happens-Before.md)

Guarantee that memory write becomes visible to other threads/processors.

---

## W

### Worker Pool
**Category**: Concurrency Pattern  
**Related**: Goroutine, Channel, Fan-Out  
**See**: [EC-013-Concurrent-Patterns.md](./03-Engineering-CloudNative/EC-013-Concurrent-Patterns.md)

Pattern using fixed set of workers processing jobs from shared queue; controls concurrency level.

### W3C Trace Context
**Category**: Observability  
**Related**: OpenTelemetry, Distributed Tracing, Headers  
**See**: [EC-070-OpenTelemetry-W3C-Trace-Context.md](./03-Engineering-CloudNative/EC-070-OpenTelemetry-W3C-Trace-Context.md)

Standard for propagating trace context across service boundaries using HTTP headers.

---

## X

### X-Request-ID
**Category**: Observability  
**Related**: Correlation ID, Tracing, Middleware  
**See**: [03-Engineering-CloudNative/02-Cloud-Native/](./03-Engineering-CloudNative/02-Cloud-Native/)

HTTP header carrying unique request identifier for request tracing across services.

---

## Y

### YAML
**Category**: Configuration  
**Related**: JSON, TOML, Configuration  
**See**: Various configuration documents

Human-readable data serialization standard commonly used for configuration files.

---

## Z

### Zero-Downtime Deployment
**Category**: DevOps  
**Related**: Blue-Green, Canary, Rolling  
**See**: [03-Engineering-CloudNative/02-Cloud-Native/](./03-Engineering-CloudNative/02-Cloud-Native/)

Deployment strategy ensuring continuous service availability during updates.

---

## Document History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | 2026-04-02 | Initial comprehensive glossary | Knowledge Base Team |

---

*For acronyms not found here, check [REFERENCES.md](./REFERENCES.md) or [INDEX.md](./INDEX.md).*
