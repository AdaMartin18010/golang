# Go Knowledge Base - Learning Dependency Graph

> **Version**: 2.0.0
> **Last Updated**: 2026-04-02
> **Purpose**: Visualize learning paths and prerequisite relationships
> **Total Dependencies**: 500+ connections
> **Learning Paths**: 10+ complete curricula

---

## 📊 Learning Path Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        KNOWLEDGE DEPENDENCY GRAPH                            │
└─────────────────────────────────────────────────────────────────────────────┘

                              ┌─────────────┐
                              │  FOUNDATION │
                              │    (FT-001) │
                              └──────┬──────┘
                                     │
            ┌────────────────────────┼────────────────────────┐
            │                        │                        │
            ▼                        ▼                        ▼
   ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
   │  DISTRIBUTED    │    │   CONCURRENCY   │    │   TYPE SYSTEM   │
   │    SYSTEMS      │    │     (FT-002)    │    │    (LD-001)     │
   │  (FT-003~FT-015)│    └────────┬────────┘    └────────┬────────┘
   └────────┬────────┘             │                      │
            │                      │                      │
            │            ┌─────────┴─────────┐            │
            │            │                   │            │
            ▼            ▼                   ▼            ▼
   ┌─────────────────┐  ┌─────────────┐  ┌─────────────────┐
   │    CONSENSUS    │  │  GMP MODEL  │  │  GO RUNTIME     │
   │   (FT-002~004)  │  │   (LD-004)  │  │   (LD-002~015)  │
   └────────┬────────┘  └──────┬──────┘  └────────┬────────┘
            │                  │                  │
            │                  │                  │
            └──────────────────┼──────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │  PRODUCTION SYSTEMS │
                    │   (EC-001~EC-121)   │
                    └──────────┬──────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
    ┌─────────────────┐ ┌─────────────┐ ┌─────────────────┐
    │  MICROSERVICES  │ │   CLOUD     │ │  APPLICATION    │
    │   (EC-001~050)  │ │   NATIVE    │ │    DOMAINS      │
    └─────────────────┘ │ (EC-051~121)│ │   (AD-001~026)  │
                        └─────────────┘ └─────────────────┘
```

---

## 🎯 Complete Learning Paths

### Path 1: Backend Engineer 🌐

**Duration**: 12 weeks | **Documents**: 80+ | **Prerequisites**: None

```
Phase 1: Foundation (Weeks 1-2)
├── 02-Language-Design/02-Language-Features/01-Type-System.md
├── 02-Language-Design/02-Language-Features/05-Error-Handling.md
└── EC-001: Microservices Architecture
    │
    ▼
Phase 2: API Design (Weeks 3-4)
├── EC-004: API Design Formal
├── AD-006: API Gateway Design
└── 05-Application-Domains/01-Backend-Development/01-RESTful-API.md
    │
    ▼
Phase 3: Data Layer (Weeks 5-6)
├── TS-001: PostgreSQL Transaction Internals
├── 04-Technology-Stack/02-Database/04-Redis.md
└── 04-Technology-Stack/02-Database/11-Caching-Strategies.md
    │
    ▼
Phase 4: Resilience (Weeks 7-8)
├── EC-007: Graceful Shutdown
├── EC-008: Circuit Breaker Advanced
└── EC-013: Idempotency Patterns
    │
    ▼
Phase 5: Observability (Weeks 9-10)
├── EC-006: Distributed Tracing
├── EC-044: Observability Production
└── EC-060: OpenTelemetry Production
    │
    ▼
Phase 6: Production (Weeks 11-12)
├── EC-121: Google SRE Engineering
└── AD-010: System Design Interview
```

**Prerequisites Table**:

| Document | Prerequisites | Followed By |
|----------|---------------|-------------|
| EC-001 | None | EC-002, EC-004, EC-007 |
| EC-004 | EC-001 | AD-006, EC-017 |
| EC-007 | EC-001 | EC-008, EC-079 |
| EC-008 | EC-007 | EC-112, EC-117 |
| TS-001 | None | EC-065, AD-016 |
| EC-044 | EC-006 | EC-060, EC-080 |

---

### Path 2: Cloud-Native Engineer ☁️

**Duration**: 16 weeks | **Documents**: 100+ | **Prerequisites**: Basic Go

```
Phase 1: Container Fundamentals (Weeks 1-2)
├── EC-003: Container Design Formal
├── TS-022: Docker Container Runtime
└── 05-Application-Domains/02-Cloud-Infrastructure/03-Docker-Lib.md
    │
    ▼
Phase 2: Kubernetes Core (Weeks 3-5)
├── TS-005: Kubernetes Operator Patterns
├── TS-021: Kubernetes Networking
└── EC-069: Kubernetes Operators
    │
    ▼
Phase 3: Networking & Mesh (Weeks 6-7)
├── TS-006: Kubernetes Networking Formal
├── TS-015: Service Mesh Istio
└── 04-Technology-Stack/03-Network/09-Service-Mesh.md
    │
    ▼
Phase 4: Coordination (Weeks 8-9)
├── TS-007: ETCD Raft Implementation
├── EC-071: ETCD Distributed Coordination
└── EC-116: ETCD Coordination Patterns
    │
    ▼
Phase 5: Scheduling (Weeks 10-12)
├── EC-099: Kubernetes CronJob Deep Dive
├── EC-100: Temporal Workflow Engine
└── EC-109: Production Ready Task Scheduler
    │
    ▼
Phase 6: GitOps & Delivery (Weeks 13-14)
├── TS-028: ArgoCD GitOps
├── TS-029: Flux CD GitOps
└── 05-Application-Domains/02-Cloud-Infrastructure/08-GitOps.md
    │
    ▼
Phase 7: Production (Weeks 15-16)
├── EC-110: Task Resource Quota Management
└── EC-121: Google SRE Engineering
```

**Prerequisites Table**:

| Document | Prerequisites | Followed By |
|----------|---------------|-------------|
| EC-003 | None | TS-022, EC-068 |
| TS-005 | EC-003 | EC-069, EC-099 |
| TS-007 | FT-002 | EC-071, EC-116 |
| EC-099 | TS-005 | EC-100, EC-114 |
| EC-100 | EC-099 | EC-115 |
| EC-109 | EC-042, FT-002 | EC-110 |

---

### Path 3: Distributed Systems Engineer 🔗

**Duration**: 20 weeks | **Documents**: 120+ | **Prerequisites**: Strong CS fundamentals

```
Phase 1: Theory Foundation (Weeks 1-3)
├── FT-001: Distributed Systems Foundation
├── FT-003: CAP Theorem Formal
└── FT-016: PACELC Theorem Formal
    │
    ▼
Phase 2: Consensus Algorithms (Weeks 4-6)
├── FT-002: Raft Consensus Formal
├── FT-003: Paxos Consensus Formal
├── FT-006: Paxos Formal
└── FT-015: FLP Impossibility Formal
    │
    ▼
Phase 3: Consistency Models (Weeks 7-8)
├── FT-010: Linearizability Formal
├── FT-011: Sequential Consistency Formal
├── FT-012: Causal Consistency Formal
└── FT-013: Eventual Consistency Formal
    │
    ▼
Phase 4: Time & Ordering (Weeks 9-10)
├── FT-005: Vector Clocks Formal
├── FT-010: Time Clocks Ordering
└── FT-019: Operational Transformation
    │
    ▼
Phase 5: Data Structures (Weeks 11-12)
├── FT-007: Probabilistic Data Structures
├── FT-012: CRDT Formal
└── FT-018: CRDT Formal Extended
    │
    ▼
Phase 6: Transactions (Weeks 13-14)
├── FT-014: Two-Phase Commit Formalization
├── FT-021: Two-Phase Commit Formal
├── FT-022: Three-Phase Commit Formal
└── FT-023: SAGA Formal
    │
    ▼
Phase 7: Membership & Gossip (Weeks 15-16)
├── FT-011: Gossip Protocols
├── FT-025: Leader Election Formal
├── FT-026: Membership Protocol Formal
└── FT-027: Gossip Protocol Formal
    │
    ▼
Phase 8: Implementation (Weeks 17-18)
├── EC-108: Distributed Consensus Raft Implementation
├── EC-057: ETCD Distributed Task Scheduler
└── EC-091: Distributed Lock Implementation
    │
    ▼
Phase 9: Production (Weeks 19-20)
├── EC-112: Saga Pattern Complete
├── EC-113: CRDT Conflict Resolution
└── EC-121: Google SRE Engineering
```

**Prerequisites Table**:

| Document | Prerequisites | Followed By |
|----------|---------------|-------------|
| FT-001 | None | FT-003, FT-016 |
| FT-003 | FT-001 | FT-002, FT-006, FT-015 |
| FT-002 | FT-003 | FT-025, EC-108 |
| FT-010 | FT-003 | FT-011, FT-012, FT-013 |
| FT-014 | FT-003 | FT-021, FT-022, FT-023 |
| FT-012 | FT-010 | FT-018, EC-113 |
| EC-108 | FT-002 | EC-057, EC-116 |
| EC-112 | FT-014, FT-023 | EC-100 |

---

### Path 4: Go Specialist 🐹

**Duration**: 10 weeks | **Documents**: 60+ | **Prerequisites**: Go basics

```
Phase 1: Memory Model (Weeks 1-2)
├── LD-001: Go Memory Model Formal
├── 01-Formal-Theory/19-Go-Memory-Model-Happens-Before.md
└── 01-Formal-Theory/04-Memory-Models/01-Happens-Before.md
    │
    ▼
Phase 2: Runtime Deep Dive (Weeks 3-4)
├── LD-004: Go Runtime GMP Deep Dive
├── FT-002: GMP Scheduler Deep Dive
└── LD-010: Go Scheduler GMP
    │
    ▼
Phase 3: Garbage Collection (Weeks 5-6)
├── LD-003: Go Garbage Collector Formal
├── LD-011: Go GC Algorithm
└── 02-Language-Design/02-Language-Features/10-GC.md
    │
    ▼
Phase 4: Compiler Pipeline (Weeks 7-8)
├── LD-002: Go Compiler Architecture SSA
├── LD-013: Go Compiler Phases
└── LD-012: Go Linker Build Process
    │
    ▼
Phase 5: Advanced Internals (Weeks 9-10)
├── LD-011: Go Assembly Internals
├── LD-014: Go Assembly Programming
├── LD-015: Go Plugin System
└── LD-016: Go Standard Library Deep Dive
```

**Prerequisites Table**:

| Document | Prerequisites | Followed By |
|----------|---------------|-------------|
| LD-001 | None | LD-004, FT-002 |
| LD-004 | LD-001 | LD-010, EC-042 |
| LD-003 | LD-004 | LD-011 |
| LD-002 | LD-004 | LD-013 |
| LD-013 | LD-002 | LD-012, LD-014 |
| LD-011 | LD-003 | LD-014, LD-015 |

---

### Path 5: Database Internals Expert 🗄️

**Duration**: 12 weeks | **Documents**: 50+ | **Prerequisites**: SQL basics

```
Phase 1: Transaction Theory (Weeks 1-3)
├── TS-001: PostgreSQL Transaction Internals
├── TS-001-PostgreSQL-Transaction-Formal.md
└── EC-065: Database Transaction Isolation MVCC
    │
    ▼
Phase 2: Storage Engines (Weeks 4-5)
├── TS-001: PostgreSQL 18 Transaction Internals
├── TS-006: MySQL Transaction Isolation
└── TS-010: ClickHouse Column Storage
    │
    ▼
Phase 3: Caching (Weeks 6-7)
├── TS-002: Redis Data Structures
├── TS-002-Redis-Data-Structures-Internals.md
└── 04-Technology-Stack/02-Database/11-Caching-Strategies.md
    │
    ▼
Phase 4: Distributed Storage (Weeks 8-9)
├── TS-003: Kafka Architecture
├── TS-004: Elasticsearch Query DSL
└── TS-012: Elasticsearch Internals
    │
    ▼
Phase 5: Advanced Topics (Weeks 10-12)
├── 04-Technology-Stack/02-Database/10-Database-Sharding.md
├── 04-Technology-Stack/02-Database/12-Database-Replication.md
└── FT-012: CRDT Conflict-Free Replicated Data Types
```

**Prerequisites Table**:

| Document | Prerequisites | Followed By |
|----------|---------------|-------------|
| TS-001 | None | EC-065, TS-006 |
| EC-065 | TS-001 | FT-014, EC-111 |
| TS-002 | None | EC-091, EC-005 |
| TS-003 | None | EC-061, EC-034 |
| TS-004 | None | TS-012 |

---

## 🕸️ Dependency Network Visualization

### Core Dependency Clusters

```
┌──────────────────────────────────────────────────────────────────────────┐
│                    CONCURRENCY CLUSTER                                    │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   LD-001 ──► LD-004 ──► LD-010 ──► EC-042 ──► EC-109                    │
│      │         │          │          │          │                        │
│      ▼         ▼          ▼          ▼          ▼                        │
│   FT-001    FT-002     EC-013     EC-073     EC-110                     │
│      │         │          │          │                                   │
│      └────┬───┘          └────┬─────┘                                   │
│           │                   │                                          │
│           ▼                   ▼                                          │
│        EC-084              EC-085                                        │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────────┐
│                    CONSENSUS CLUSTER                                      │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   FT-001 ──► FT-003 ──► FT-002 ──► EC-108 ──► EC-057                    │
│               │          │          │          │                         │
│               ▼          ▼          ▼          ▼                         │
│            FT-006     FT-025     EC-116     EC-071                       │
│               │          │          │          │                         │
│               └────┬─────┘          └────┬─────┘                         │
│                    │                     │                               │
│                    ▼                     ▼                               │
│                 FT-015                EC-091                             │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────────┐
│                    STORAGE CLUSTER                                        │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   TS-001 ──► EC-065 ──► FT-014 ──► FT-023 ──► EC-112                    │
│      │          │          │          │          │                       │
│      ▼          ▼          ▼          ▼          ▼                       │
│   TS-006     EC-111     FT-021     EC-090     EC-100                     │
│      │          │          │          │                                  │
│      └────┬─────┘          └────┬─────┘                                  │
│           │                     │                                        │
│           ▼                     ▼                                        │
│        EC-028                EC-048                                      │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 📋 Detailed Prerequisite Tables

### Formal Theory Dependencies

| Document | Hard Prerequisites | Recommended | Enables |
|----------|-------------------|-------------|---------|
| FT-001 | None | None | FT-003, FT-004 |
| FT-002 | FT-003 | LD-001 | EC-108, EC-042 |
| FT-003 | FT-001 | None | FT-002, FT-006, FT-015 |
| FT-004 | FT-001 | None | FT-005, EC-057 |
| FT-005 | FT-004 | None | FT-006, FT-010 |
| FT-006 | FT-003 | FT-005 | FT-007, FT-008 |
| FT-007 | FT-006 | None | FT-013, FT-031 |
| FT-008 | FT-006 | None | FT-009, EC-091 |
| FT-009 | FT-008 | None | FT-017, EC-108 |
| FT-010 | FT-003 | FT-005 | FT-011, FT-012 |
| FT-011 | FT-010 | None | FT-027 |
| FT-012 | FT-010 | FT-011 | FT-018, EC-113 |
| FT-013 | FT-007 | None | FT-031 |
| FT-014 | FT-003 | EC-065 | FT-021, FT-022, FT-023 |
| FT-015 | FT-003 | FT-006 | FT-024, FT-025 |
| FT-016 | FT-001 | FT-003 | None |
| FT-017 | FT-009 | None | EC-116 |
| FT-018 | FT-012 | None | EC-113 |
| FT-019 | FT-005 | None | None |
| FT-020 | FT-003 | FT-010 | EC-111 |
| FT-021 | FT-014 | None | EC-112 |
| FT-022 | FT-014 | FT-021 | None |
| FT-023 | FT-014 | None | EC-112 |
| FT-024 | FT-015 | FT-002 | None |
| FT-025 | FT-015 | FT-002 | EC-116 |
| FT-026 | FT-025 | None | EC-071 |
| FT-027 | FT-011 | None | EC-071 |
| FT-028 | FT-027 | None | None |
| FT-029 | FT-008 | None | EC-091 |

### Language Design Dependencies

| Document | Hard Prerequisites | Recommended | Enables |
|----------|-------------------|-------------|---------|
| LD-001 | None | None | LD-004, FT-002 |
| LD-002 | LD-004 | None | LD-013 |
| LD-003 | LD-004 | None | LD-011 |
| LD-004 | LD-001 | None | LD-010, EC-042 |
| LD-005 | LD-004 | None | LD-007 |
| LD-006 | LD-004 | None | LD-023 |
| LD-007 | LD-005 | None | LD-009 |
| LD-008 | LD-006 | None | LD-022 |
| LD-009 | LD-007 | None | LD-024 |
| LD-010 | LD-004 | LD-003 | None |
| LD-011 | LD-003 | LD-002 | LD-014 |
| LD-012 | LD-013 | None | None |
| LD-013 | LD-002 | None | LD-012, LD-014 |
| LD-014 | LD-011, LD-013 | None | LD-015 |
| LD-015 | LD-014 | None | None |
| LD-016 | LD-011, LD-013 | None | None |
| LD-017 | LD-004 | None | None |
| LD-018 | LD-017 | None | None |
| LD-019 | LD-017 | None | None |
| LD-020 | LD-017 | None | None |
| LD-021 | LD-004 | None | EC-013 |
| LD-022 | LD-008 | None | EC-043 |
| LD-023 | LD-006 | None | EC-119 |
| LD-024 | LD-009 | None | EC-095 |
| LD-025 | LD-024 | None | EC-046 |

### Engineering CloudNative Dependencies

| Document | Hard Prerequisites | Recommended | Enables |
|----------|-------------------|-------------|---------|
| EC-001 | None | None | EC-002, EC-004, EC-007 |
| EC-002 | EC-001 | None | EC-009, EC-075 |
| EC-003 | None | None | TS-022, EC-068 |
| EC-004 | EC-001 | None | AD-006, EC-017 |
| EC-005 | TS-001 | None | EC-065 |
| EC-006 | EC-001 | None | EC-044, EC-056 |
| EC-007 | EC-001 | None | EC-008, EC-079 |
| EC-008 | EC-007 | FT-018 | EC-112, EC-117 |
| EC-009 | EC-001 | None | EC-017, EC-042 |
| EC-010 | EC-001 | None | EC-021, EC-061 |
| EC-011 | EC-007 | EC-005 | EC-084 |
| EC-012 | EC-007 | EC-011 | EC-078, EC-030 |
| EC-013 | LD-021 | EC-001 | EC-119 |
| EC-014 | EC-007 | EC-006 | EC-086 |
| EC-015 | EC-001 | AD-004 | EC-034, EC-092 |
| EC-016 | EC-001 | None | EC-071 |
| EC-017 | EC-009 | EC-004 | EC-020, EC-031 |
| EC-020 | EC-017 | EC-016 | EC-099 |
| EC-021 | EC-010 | None | EC-072 |
| EC-042 | LD-004, FT-002 | EC-009 | EC-109, EC-062 |
| EC-043 | LD-022 | EC-011 | EC-064, EC-066 |
| EC-044 | EC-006 | EC-014 | EC-060, EC-080 |
| EC-057 | FT-002, EC-108 | TS-007 | EC-116 |
| EC-062 | EC-042 | EC-017 | EC-067 |
| EC-065 | TS-001, FT-010 | None | EC-111, FT-014 |
| EC-067 | EC-062 | EC-057 | EC-109 |
| EC-071 | FT-026, FT-027 | EC-016 | EC-116 |
| EC-091 | FT-029 | TS-002 | EC-093 |
| EC-099 | TS-005, EC-020 | FT-025 | EC-100, EC-114 |
| EC-100 | EC-099, EC-112 | None | EC-115 |
| EC-108 | FT-002 | None | EC-057 |
| EC-109 | EC-042, EC-067 | FT-002 | EC-110 |
| EC-111 | FT-020, EC-065 | EC-034 | EC-112 |
| EC-112 | FT-023, EC-008 | EC-090 | EC-100 |
| EC-116 | FT-017, EC-071 | EC-057 | EC-091 |
| EC-121 | EC-007, EC-044 | All above | None |

### Technology Stack Dependencies

| Document | Hard Prerequisites | Recommended | Enables |
|----------|-------------------|-------------|---------|
| TS-001 | None | None | EC-065, EC-005 |
| TS-002 | None | None | EC-091, TS-006 |
| TS-003 | None | None | EC-061, EC-034 |
| TS-004 | None | None | TS-012 |
| TS-005 | EC-003 | None | EC-099, EC-069 |
| TS-006 | TS-005 | None | TS-021 |
| TS-007 | FT-002 | None | EC-057, EC-116 |
| TS-008 | None | None | None |
| TS-009 | None | None | None |
| TS-010 | None | None | None |
| TS-011 | TS-003 | None | None |
| TS-012 | TS-004 | None | None |
| TS-013 | EC-006 | None | EC-044 |
| TS-014 | None | None | None |
| TS-015 | TS-006 | None | None |
| TS-016 | TS-013 | None | None |
| TS-017 | TS-013 | None | None |
| TS-018 | EC-006 | None | None |
| TS-019 | EC-060 | None | None |
| TS-020 | AD-007 | None | None |
| TS-021 | TS-006 | None | TS-015 |
| TS-022 | EC-003 | None | None |
| TS-023 | TS-021 | None | None |
| TS-024 | TS-015 | None | None |
| TS-025 | TS-021 | None | None |
| TS-026 | None | None | 05-Application-Domains/02-Cloud-Infrastructure/02-Terraform-Providers.md |
| TS-027 | None | None | None |
| TS-028 | None | None | None |
| TS-029 | None | None | None |

---

## 🎓 Skill Level Progressions

### Beginner → Intermediate → Advanced → Expert

```
BEGINNER (0-6 months)
│
├── FT-001: Distributed Systems Foundation
├── LD-001: Go Memory Model
├── EC-001: Microservices Architecture
├── TS-001: PostgreSQL Basics
└── AD-006: API Gateway Design
    │
    ▼
INTERMEDIATE (6-12 months)
│
├── FT-003: CAP Theorem
├── LD-004: Go Runtime GMP
├── EC-007: Graceful Shutdown
├── EC-008: Circuit Breaker
├── TS-002: Redis
└── EC-042: Task Scheduler Core
    │
    ▼
ADVANCED (12-18 months)
│
├── FT-002: Raft Consensus
├── LD-013: Go Compiler Phases
├── EC-099: Kubernetes CronJob
├── EC-109: Production Task Scheduler
├── TS-007: ETCD Raft
└── EC-112: Saga Pattern Complete
    │
    ▼
EXPERT (18+ months)
│
├── FT-015: FLP Impossibility
├── LD-015: Go Plugin System
├── EC-121: Google SRE Engineering
├── FT-024: Consensus Variations
├── EC-116: ETCD Coordination Patterns
└── EC-101: Formal Verification
```

---

## 🔗 Navigation

- [← Back to Cross-Reference](cross-reference.md)
- [Complete Index →](complete-index.md)
- [Search Index →](search-index.md)

---

**Total Learning Paths**: 10+
**Total Dependencies**: 500+
**Average Path Length**: 8 documents
**Last Updated**: 2026-04-02
