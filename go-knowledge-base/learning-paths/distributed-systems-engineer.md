# Distributed Systems Engineer Learning Path

> **Version**: 1.0.0
> **Last Updated**: 2026-04-02
> **Duration**: 24 weeks (full-time) / 36 weeks (part-time)
> **Prerequisites**: Strong CS fundamentals, distributed systems basics, Go proficiency
> **Outcome**: Expert in designing and implementing fault-tolerant, scalable distributed systems

---

## 🎯 Path Overview

### Target Competencies

Upon completion, you will be able to:

- Design consensus protocols and distributed coordination systems
- Implement replication, sharding, and partitioning strategies
- Build fault-tolerant systems handling network partitions and failures
- Design eventually consistent and strongly consistent systems
- Implement distributed transactions and saga patterns
- Analyze systems using formal methods and verification
- Contribute to distributed databases, message queues, or coordination services

### Prerequisites Graph

```
Strong CS Fundamentals
    ├── Algorithms & Data Structures
    ├── Operating Systems
    ├── Computer Networks
    └── Database Systems
            ↓
Go Proficiency + Backend Experience
    ├── Production Go experience
    ├── Concurrency patterns
    └── API design
            ↓
┌─────────────────────────────────────────────────────────────────────┐
│           DISTRIBUTED SYSTEMS ENGINEER LEARNING PATH                 │
│                                                                      │
│  Phase 1: Theoretical Foundations (Weeks 1-6)                       │
│    ├── Distributed Systems Theory → CAP → Consistency Models        │
│    └── Outcome: Theoretical foundation                              │
│                                                                      │
│  Phase 2: Consensus & Coordination (Weeks 7-12)                     │
│    ├── Raft → Paxos → Byzantine Fault Tolerance                     │
│    └── Outcome: Consensus algorithm mastery                         │
│                                                                      │
│  Phase 3: Data & Storage (Weeks 13-18)                              │
│    ├── Replication → Sharding → CRDTs → Event Sourcing              │
│    └── Outcome: Distributed data expertise                          │
│                                                                      │
│  Phase 4: Advanced Topics (Weeks 19-24)                             │
│    ├── Formal Verification → System Design → Research               │
│    └── Outcome: Research-level expertise                            │
└─────────────────────────────────────────────────────────────────────┘
    ↓
Career Paths
    ├── Distributed Database Engineer
    ├── Platform Engineer (Infrastructure)
    ├── Research Engineer
    └── Staff/Principal Engineer
```

---

## 📚 Phase 1: Theoretical Foundations (Weeks 1-6)

### Week 1: Distributed Systems Fundamentals

**Goal**: Establish theoretical foundation

#### Day 1-3: Core Concepts

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-001] Distributed Systems Foundation | 6h | CAP, BASE, ACID |
| [FT-004] Distributed Systems Fundamentals | 4h | CAP/BASE/ACID deep dive |

**Study Notes**:

- **CAP Theorem**: Choose 2 of {Consistency, Availability, Partition tolerance}
- **PACELC**: If partitioned, choose {Availability, Consistency}; else choose {Latency, Consistency}
- **BASE**: Basically Available, Soft state, Eventually consistent
- **ACID** vs **BASE** trade-offs

**Key Insights**:

- Partitions are inevitable, so CP or AP choice is critical
- Most systems are actually PA/EL or PC/EC
- Network latency dominates distributed system design

#### Day 4-5: Consistency Models

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-001] Consistency Models | 4h | Linearizability, sequential, causal |

**Study Notes**:

- **Linearizability**: Strongest consistency, every operation appears instantaneous
- **Sequential**: Operations appear in some sequential order
- **Causal**: If op A happens before B, all see A before B
- **Eventual**: If no new updates, all converge to same value

#### Day 6-7: Failure Models

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-008] Network Partition & Brain Split | 4h | Partition handling |

**Study Notes**:

- Crash-stop vs crash-recovery
- Byzantine failures
- Network partitions and split-brain
- Failure detection

**Week 1 Capstone**:

```
Analyze these systems:
1. Cassandra - What consistency model? Why?
2. ZooKeeper - CP or AP? Evidence?
3. DynamoDB - How does it handle partitions?
4. Design a photo storage system - What trade-offs?
```

### Week 2: Time and Ordering

**Goal**: Understand logical time and event ordering

#### Day 1-3: Logical Time

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-005] Vector Clocks | 5h | Vector clock algorithm |
| [FT-006] Vector Clocks & Logical Time | 4h | Causal ordering |
| [FT-010] Time Clocks & Ordering | 4h | Ordering theory |

**Study Notes**:

- **Lamport Timestamps**: Partial ordering, compact
- **Vector Clocks**: Full causal history, detects concurrency
- **Version Vectors**: Track replica divergences
- **HLC (Hybrid Logical Clocks)**: Physical + logical time

**Algorithms**:

```
Vector Clock Update:
  On local event: VC[self]++
  On receive(m, VC_m): VC[i] = max(VC[i], VC_m[i]) for all i

Concurrent Events:
  VC_a || VC_b iff neither VC_a < VC_b nor VC_b < VC_a
```

#### Day 4-5: Physical Time

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-010] Time Clocks & Ordering | 3h | NTP, clock skew |

**Study Notes**:

- NTP and time synchronization
- Clock skew and drift
- TrueTime (Spanner)
- Timestamping transactions

#### Day 6-7: Causality Tracking

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-005] Vector Clocks | 3h | Practical usage |

**Week 2 Capstone**:

```go
// Implement vector clocks:
// - Merge operation
// - Concurrency detection
// - Pruning strategies
// - Visualize causality
```

### Week 3: Distributed Algorithms

**Goal**: Learn fundamental distributed algorithms

#### Day 1-3: Gossip Protocols

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-011] Gossip Protocols | 5h | Epidemic broadcast |

**Study Notes**:

- **Push gossip**: Infect others proactively
- **Pull gossip**: Request updates from peers
- **Push-pull gossip**: Combined approach
- Applications: membership, failure detection, aggregation

**Analysis**:

- Time to infect all: O(log N) rounds
- Bandwidth per node: O(log N)
- Robust to failures

#### Day 4-5: Membership Protocols

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-011] Gossip Protocols | 3h | SWIM protocol |

**Study Notes**:

- **SWIM**: Scalable Weakly-consistent Infection-style Process Group Membership
- Failure detection via probing
- Suspicion mechanism
- Lifeguard extensions

#### Day 6-7: Quorum Systems

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-008] Quorum Consensus Theory | 5h | Quorum design |
| [FT-009] Quorum Consensus | 4h | Implementation |

**Study Notes**:

- **Read quorum (R)** + **Write quorum (W)**
- Constraint: R + W > N (total replicas)
- Grid quorums for efficiency
- Probabilistic quorums

**Week 3 Capstone**:

```go
// Implement gossip protocol:
// - Membership list dissemination
// - Failure detection
// - Message broadcast
// - Measure convergence time
```

### Week 4: Sharding and Partitioning

**Goal**: Master data partitioning strategies

#### Day 1-3: Consistent Hashing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-004] Consistent Hashing | 5h | Algorithm |
| [FT-005] Consistent Hashing | 3h | Virtual nodes |

**Study Notes**:

- **Ring hash**: Map keys and nodes to ring
- **Virtual nodes**: Better distribution
- **Replication**: Next N nodes hold replica
- **Rebalancing**: Only k/n keys move when node added

**Algorithm**:

```
Hash space: [0, 2^160) for SHA-1
Map node to multiple points (virtual nodes)
Map key to closest node clockwise
```

#### Day 4-5: Partitioning Strategies

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-082] Distributed Task Sharding | 4h | Sharding patterns |
| 04-Technology-Stack/02-Database/10-Database-Sharding.md | 3h | DB sharding |

**Study Notes**:

- **Hash partitioning**: Good distribution, range queries hard
- **Range partitioning**: Efficient range queries, hot spots
- **List partitioning**: Explicit mapping
- **Composite**: Combine strategies

#### Day 6-7: Rebalancing

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-004] Consistent Hashing | 3h | Rebalancing |

**Study Notes**:

- Online rebalancing
- Consistent hashing minimal movement
- Rendezvous hashing (HRW)
- Jump consistent hash

**Week 4 Capstone**:

```go
// Implement consistent hashing:
// - Ring with virtual nodes
// - Add/remove nodes
// - Replication
// - Benchmark rebalancing cost
```

### Week 5: Formal Semantics

**Goal**: Understand formal specification methods

#### Day 1-3: Operational Semantics

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 01-Formal-Theory/01-Semantics/01-Operational-Semantics.md | 5h | SOS rules |
| 01-Formal-Theory/01-Semantics/02-Denotational-Semantics.md | 4h | Denotational |

**Study Notes**:

- **SOS (Structural Operational Semantics)**: Small-step rules
- **Big-step semantics**: Direct evaluation
- **Transition systems**: States and transitions
- **Labeled transitions**: Actions and observations

#### Day 4-5: Process Calculi

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 01-Formal-Theory/03-Concurrency-Models/01-CSP-Theory.md | 5h | CSP |

**Study Notes**:

- **CSP (Communicating Sequential Processes)**: Hoare's formalism
- Events and processes
- Parallel composition
- Trace semantics

#### Day 6-7: Applying Formal Methods

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 01-Formal-Theory/01-Semantics/04-Featherweight-Go.md | 4h | FG calculus |

**Study Notes**:

- Featherweight Go: Core Go formalization
- Structural typing
- Type safety proofs

### Week 6: Distributed Systems Case Studies

**Goal**: Study real-world systems

#### Day 1-2: Dynamo

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research papers + [FT-004] | 4h | Dynamo design |

**Study Notes**:

- Consistent hashing with virtual nodes
- Quorum reads/writes
- Vector clocks for versioning
- Gossip for membership

#### Day 3-4: Spanner

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research papers + [FT-010] | 4h | Spanner design |

**Study Notes**:

- TrueTime API
- External consistency (linearizability)
- 2PC with Paxos groups
- Global transactions

#### Day 5-6: Ceph

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research papers | 4h | CRUSH algorithm |

**Study Notes**:

- CRUSH: Controlled Replication Under Scalable Hashing
- Pseudorandom placement
- Device failure handling
- Data distribution

#### Day 7: Review

| Document | Time | Key Takeaways |
|----------|------|---------------|
| All week 1-6 docs | 6h | Synthesis |

---

## 📚 Phase 2: Consensus & Coordination (Weeks 7-12)

### Week 7: Consensus Fundamentals

**Goal**: Understand consensus problem deeply

#### Day 1-3: Consensus Problem

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-003] Distributed Consensus | 5h | Problem definition |
| [FT-015] Consensus Lower Bounds | 5h | Impossibility results |

**Study Notes**:

- **FLP impossibility**: No deterministic consensus in async systems
- **Safety**: All agree on same value
- **Liveness**: Eventually decide
- **Byzantine vs Crash-stop**

**Lower Bounds**:

- Minimum rounds for consensus
- Message complexity
- Lower bounds under synchrony assumptions

#### Day 4-5: Two-Phase Commit

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-014] Two-Phase Commit Formalization | 5h | 2PC protocol |

**Study Notes**:

- Coordinator and participants
- Voting phase
- Commit/abort phase
- Blocking problem

#### Day 6-7: Three-Phase Commit

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-014] Two-Phase Commit Formalization | 3h | 3PC overview |

**Study Notes**:

- Non-blocking during coordinator failure
- CanCommit, PreCommit, DoCommit
- Trade-offs vs 2PC

**Week 7 Capstone**:

```
Compare consensus protocols:
- 2PC: When appropriate? Limitations?
- 3PC: Does it solve blocking?
- When to use each?
```

### Week 8: Raft Consensus

**Goal**: Master Raft algorithm

#### Day 1-3: Raft Basics

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-002] Raft Consensus | 6h | Raft paper |
| [FT-002] Raft Consensus Formal | 4h | Formal specification |

**Study Notes**:

- **Leader election**: Randomized timeouts
- **Log replication**: Append entries
- **Safety**: Leader completeness
- **Membership changes**: Joint consensus

**Key Properties**:

```
Election Safety: At most one leader per term
Leader Append-Only: Leaders never overwrite entries
Log Matching: If logs have same index/term, identical up to there
Leader Completeness: Committed entries in all future leaders
State Machine Safety: Same index = same command
```

#### Day 4-5: Raft Implementation

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-108] Raft Implementation | 6h | Building Raft |

**Study Notes**:

- State machine design
- RPC handling
- Persistent state
- Snapshotting

#### Day 6-7: Raft Optimizations

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-057] ETCD Distributed Scheduler | 4h | Production Raft |

**Study Notes**:

- Pre-vote mechanism
- Leader leasing
- Batching and pipelining
- Checkpoints

**Week 8 Capstone**:

```go
// Implement Raft:
// - Leader election
// - Log replication
// - Commit logic
// - Pass Raft test suite
```

### Week 9: Paxos and Multi-Paxos

**Goal**: Understand Paxos family

#### Day 1-3: Paxos Basics

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-003] Paxos Consensus | 6h | Classic Paxos |
| [FT-003] Paxos Consensus Formal | 4h | Formalization |

**Study Notes**:

- **Prepare/Promise**: Establish round
- **Accept**: Propose value
- **Quorum intersection**: Ensures safety
- **Liveness**: Requires fair leader election

**Algorithm**:

```
Proposer:
  Phase 1: Send prepare(n) to acceptors
           If majority promise, proceed
  Phase 2: Send accept(n, v) to acceptors
           If majority accept, decided

Acceptor:
  Promise not to accept lower n
  Accept if n matches promise
```

#### Day 4-5: Multi-Paxos

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-003] Paxos Consensus | 4h | Multi-Paxos |

**Study Notes**:

- Stable leader optimization
- Eliminate prepare for each slot
- Log compaction
- Leader failover

#### Day 6-7: Raft vs Paxos

| Document | Time | Key Takeaways |
|----------|------|---------------|
| ../COMPARISON-Raft-vs-Paxos.md | 4h | Comparison |

**Study Notes**:

- Understandability vs optimality
- When to choose each
- Implementation complexity
- Performance characteristics

**Week 9 Capstone**:

```
Deep dive comparison:
- Safety proofs for both
- Liveness guarantees
- Implementation complexity
- Performance benchmarks
```

### Week 10: Byzantine Fault Tolerance

**Goal**: Understand BFT consensus

#### Day 1-3: Byzantine Generals Problem

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-006] Byzantine Fault Tolerance | 6h | BFT theory |
| [FT-013] BFT | 5h | Algorithms |

**Study Notes**:

- **Byzantine failures**: Arbitrary/malicious behavior
- **3f+1** nodes for f faults (vs 2f+1 for crash-stop)
- **PBFT**: Practical BFT
- **HotStuff**: Linear communication

#### Day 4-5: PBFT Deep Dive

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-013] BFT | 4h | PBFT details |

**Study Notes**:

- Three-phase protocol
- View changes
- Checkpointing
- Performance optimizations

#### Day 6-7: Modern BFT

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research papers (HotStuff, Tendermint) | 4h | Modern BFT |

**Study Notes**:

- Chained BFT (HotStuff)
- Proposer rotation
- Responsiveness
- Blockchain applications

**Week 10 Capstone**:

```
Analyze BFT requirements:
- When is BFT needed?
- Cost of BFT (3f+1)
- Applications: blockchain, aerospace
- Comparison with crash-stop consensus
```

### Week 11: Distributed Locking and Coordination

**Goal**: Implement coordination primitives

#### Day 1-3: Distributed Locks

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-091] Distributed Lock Implementation | 5h | Lock algorithms |

**Study Notes**:

- **Redlock**: Redis-based locking
- **ZooKeeper locks**: Ephemeral sequential nodes
- **Fencing tokens**: Prevent delayed lock holders
- **Lease-based**: Timeouts for safety

#### Day 4-5: Leader Election

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-057] ETCD Distributed Scheduler | 4h | Leader election |
| [EC-071] ETCD Coordination | 3h | Coordination |

**Study Notes**:

- Bully algorithm
- Ring algorithm
- ZooKeeper recipe
- Kubernetes leader election

#### Day 6-7: Barriers and Queues

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-071] ETCD Coordination | 3h | Barriers |

**Study Notes**:

- Double barriers
- Distributed queues
- Priority queues
- ZooKeeper recipes

**Week 11 Capstone**:

```go
// Implement distributed lock:
// - ETCD-based
// - Fencing tokens
// - Lease renewal
// - Test with chaos
```

### Week 12: State Machine Replication

**Goal**: Understand SMR deeply

#### Day 1-3: SMR Theory

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-009] State Machine Replication | 5h | SMR foundations |

**Study Notes**:

- **Deterministic state machines**: Same input → same output
- **Log replication**: Consensus on command sequence
- **Execution**: Apply in order
- **Snapshots**: Truncate logs

#### Day 4-5: SMR Implementation

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-108] Raft Implementation | 4h | SMR with Raft |

**Study Notes**:

- Command log
- State machine interface
- Snapshotting strategies
- Membership changes

#### Day 6-7: Advanced SMR

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research papers | 4h | Multi-core SMR, parallel SMR |

**Study Notes**:

- Parallel execution
- Speculative execution
- Early commit
- Partitioned state machines

---

## 📚 Phase 3: Data & Storage (Weeks 13-18)

### Week 13: Replication Strategies

**Goal**: Master data replication

#### Day 1-3: Replication Models

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-001] PostgreSQL Transaction Internals | 4h | PostgreSQL replication |
| 04-Technology-Stack/02-Database/12-Database-Replication.md | 4h | Replication |

**Study Notes**:

- **Primary-backup**: Single primary, async/sync replicas
- **Multi-primary**: Conflicts, convergence
- **Chain replication**: Throughput vs latency
- **Quorum replication**: Flexible consistency

#### Day 4-5: Consensus-Based Replication

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-002] Raft Consensus | 4h | Raft replication |

**Study Notes**:

- Strong consistency via consensus
- Log shipping
- State transfer
- Dynamic membership

#### Day 6-7: Conflict Resolution

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-113] CRDT Conflict Resolution | 4h | Conflict handling |

**Week 13 Capstone**:

```
Design replication for:
- Banking (strong consistency)
- Social media likes (eventual)
- Shopping cart (session)
- Analytics (batch)
```

### Week 14: CRDTs

**Goal**: Understand Conflict-Free Replicated Data Types

#### Day 1-3: CRDT Theory

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-012] CRDTs | 6h | CRDT theory |

**Study Notes**:

- **State-based CRDTs**: Merge function
- **Operation-based CRDTs**: Causal broadcast
- **Strong eventual consistency**: Convergence guaranteed
- **Bounded growth**: Garbage collection

**Examples**:

```
G-Counter (Grow-only counter):
  merge(a, b) = max(a[i], b[i]) for all i

PN-Counter (Increment/decrement):
  Two G-counters: P for increments, N for decrements

G-Set (Grow-only set):
  merge(a, b) = union(a, b)

OR-Set (Observed-remove set):
  Add unique tags, remove observed tags
```

#### Day 4-5: CRDT Implementations

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-113] CRDT Conflict Resolution | 5h | Implementation |

**Study Notes**:

- Delta state CRDTs
- CRDTs in Redis
- Riak DT
- Automerge

#### Day 6-7: CRDT Applications

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [FT-012] CRDTs | 3h | Applications |

**Study Notes**:

- Collaborative editing
- Shopping carts
- Presence indicators
- Counters and gauges

**Week 14 Capstone**:

```go
// Implement CRDTs:
// - G-Counter
// - PN-Counter
// - G-Set
// - OR-Set
// - Test convergence
```

### Week 15: Event Sourcing

**Goal**: Master event-driven persistence

#### Day 1-3: Event Sourcing Fundamentals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-015] Event Sourcing | 5h | Patterns |
| [EC-016] CQRS Pattern | 4h | CQRS |

**Study Notes**:

- **Event store**: Immutable event log
- **Aggregates**: Domain entities
- **Projections**: Read models
- **Snapshots**: Performance optimization

#### Day 4-5: Event Store Design

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-034] Task Event Sourcing | 4h | Implementation |
| [EC-092] Event Sourcing Persistence | 4h | Storage |
| [EC-111] Event Sourcing Implementation | 4h | Complete |

**Study Notes**:

- Event store schema
- Optimistic concurrency
- Event versioning
- Schema evolution

#### Day 6-7: Event-Driven Architecture

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [AD-004] Event-Driven Architecture | 4h | EDA patterns |

**Week 15 Capstone**:

```go
// Implement event sourcing:
// - Event store
// - Aggregate roots
// - Projections
// - Snapshotting
```

### Week 16: Distributed Transactions

**Goal**: Handle transactions across services

#### Day 1-3: Saga Pattern

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-008] Saga Pattern | 5h | Saga theory |
| [EC-090] Saga Implementation | 4h | Implementation |
| [EC-112] Saga Pattern Complete | 4h | Complete guide |

**Study Notes**:

- **Choreography**: Events trigger actions
- **Orchestration**: Central coordinator
- **Compensating transactions**: Undo operations
- **Saga log**: For recovery

#### Day 4-5: Practical Implementation

| Document | Time | Key Takeaways |
|----------|------|---------------|
| ../examples/saga/ | 6h | Working example |

**Study Notes**:

- Saga orchestrator
- Compensation logic
- Timeout handling
- Idempotency

#### Day 6-7: Outbox Pattern

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-111] Event Sourcing | 3h | Outbox |

**Study Notes**:

- Atomic database + message
- Relay process
- De-duplication
- Ordering guarantees

### Week 17: Kafka Internals

**Goal**: Deep dive into distributed messaging

#### Day 1-3: Kafka Architecture

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-003] Kafka KRaft Internals | 5h | KRaft mode |
| [TS-011] Kafka Internals | 5h | Architecture |
| [TS-011] Kafka Internals Formal | 4h | Formal model |

**Study Notes**:

- **KRaft**: Kafka without ZooKeeper
- **Partition leadership**: Controller election
- **Replication**: ISR (In-Sync Replicas)
- **Log compaction**: Key-based retention

#### Day 4-5: Kafka Replication

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-003] Kafka Internals Replication | 4h | Replication |

**Study Notes**:

- acks=0,1,all trade-offs
- Min ISR configuration
- Unclean leader election
- Partition reassignment

#### Day 6-7: Exactly-Once Semantics

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-011] Kafka Internals | 4h | EOS |

**Study Notes**:

- Idempotent producer
- Transactions
- Consumer isolation levels
- End-to-end exactly-once

### Week 18: Database Internals

**Goal**: Understand distributed databases

#### Day 1-3: PostgreSQL Internals

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-001] PostgreSQL Transaction Internals | 5h | MVCC |
| [TS-001] PostgreSQL Transaction Formal | 4h | Formal |
| [EC-065] Transaction Isolation & MVCC | 4h | Deep dive |

**Study Notes**:

- MVCC implementation
- Transaction ID wraparound
- Vacuum process
- HOT updates

#### Day 4-5: Distributed SQL

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research (CockroachDB, TiDB, Yugabyte) | 6h | Distributed SQL |

**Study Notes**:

- Spanner-like architecture
- Distributed transactions
- Partitioning and rebalancing
- Consistent backups

#### Day 6-7: Redis Cluster

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [TS-002] Redis Internals | 4h | Redis Cluster |
| [TS-007] Redis Data Structures | 3h | Internals |

---

## 📚 Phase 4: Advanced Topics (Weeks 19-24)

### Week 19: Formal Verification

**Goal**: Learn formal verification methods

#### Day 1-3: Model Checking

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 01-Formal-Theory/03-Program-Verification/03-Model-Checking.md | 5h | Model checking |
| [EC-101] Formal Verification | 4h | Application |

**Study Notes**:

- **TLA+**: Temporal Logic of Actions
- **State space exploration**: Reachability
- **Invariants**: Properties to verify
- **Liveness**: Eventually properties

#### Day 4-5: TLA+ for Distributed Systems

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research + tutorials | 6h | TLA+ specs |

**Study Notes**:

- Specifying consensus
- PlusCal algorithm language
- Model checker (TLC)
- Proof system (TLAPS)

#### Day 6-7: Verification in Practice

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-101] Formal Verification | 4h | Practical |

**Week 19 Capstone**:

```tla
// Write TLA+ spec for:
// - Simple consensus
// - Verify safety
// - Check liveness
```

### Week 20: System Design

**Goal**: Design large-scale systems

#### Day 1-3: System Design Methodology

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [AD-010] System Design Interview | 5h | Design approach |
| [AD-010] System Design Interview Formal | 4h | Formal methods |

**Study Notes**:

- Requirements gathering
- Back-of-envelope calculations
- High-level design
- Deep dive components

#### Day 4-5: Case Studies

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-103] Real-World Case Studies | 5h | Case studies |

**Study Notes**:

- Design trade-offs
- Scalability patterns
- Failure scenarios
- Evolution paths

#### Day 6-7: Capacity Planning

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [AD-009] Capacity Planning | 4h | Planning |
| [AD-009] Capacity Planning Formal | 3h | Formal |

### Week 21: Performance Optimization

**Goal**: Optimize distributed systems

#### Day 1-3: Performance Theory

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [AD-008] Performance Optimization | 5h | Optimization |
| [AD-008] Performance Optimization Formal | 4h | Formal |

**Study Notes**:

- Queueing theory
- Little's Law
- Bottleneck analysis
- Latency tail tolerance

#### Day 4-5: Benchmarking

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-102] Performance Benchmarking | 5h | Methodology |
| 03-Engineering-CloudNative/03-Performance/03-Benchmarking.md | 3h | Go benchmarking |

**Study Notes**:

- Jepsen testing
- Chaos testing
- Load testing
- Profiling distributed systems

#### Day 6-7: Compiler Optimizations

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [EC-106] Compiler Optimizations | 4h | Go compiler |

### Week 22: Security Patterns

**Goal**: Secure distributed systems

#### Day 1-3: Security Foundations

| Document | Time | Key Takeaways |
|----------|------|---------------|
| [AD-007] Security Patterns | 5h | Patterns |
| [AD-007] Security Patterns Formal | 4h | Formal |
| [EC-045] Task Security Hardening | 4h | Hardening |

**Study Notes**:

- Authentication in distributed systems
- Authorization (RBAC, ABAC)
- Encryption in transit and at rest
- Secret distribution

#### Day 4-5: Cryptography

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/04-Security/03-Cryptography.md | 5h | Crypto |

**Study Notes**:

- Symmetric vs asymmetric
- Certificate management
- mTLS
- Key rotation

#### Day 6-7: Zero Trust

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 03-Engineering-CloudNative/04-Security/06-Zero-Trust.md | 4h | Zero trust |

### Week 23: Emerging Research

**Goal**: Stay current with research

#### Day 1-3: Recent Consensus Research

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research papers | 6h | Flexible Paxos, EPaxos |

**Study Notes**:

- Flexible Paxos
- EPaxos (leaderless)
- Caeser
- New directions

#### Day 4-5: Distributed ML

| Document | Time | Key Takeaways |
|----------|------|---------------|
| Research papers | 6h | Parameter servers, federated |

#### Day 6-7: Edge and IoT

| Document | Time | Key Takeaways |
|----------|------|---------------|
| 05-Application-Domains/02-Cloud-Infrastructure/09-Edge-Computing.md | 4h | Edge |

### Week 24: Research Project

**Goal**: Contribute to distributed systems knowledge

#### Options

1. **Implement a novel consensus variant**
2. **Formal verification of a distributed algorithm**
3. **Performance analysis of existing systems**
4. **Design a new distributed data structure**
5. **Contribute to open source** (etcd, Consul, CockroachDB)

**Deliverables**:

- Working implementation
- Benchmark results
- Design document
- Presentation

---

## 🎓 Capstone Project: Distributed Key-Value Store

### Requirements

**Architecture**:

```
┌──────────────────────────────────────────────────────────────┐
│                         Client Layer                          │
│                  (SDK with consistency levels)                │
└───────────────────────────┬──────────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────┐
│                      API Gateway Layer                        │
│              (Request routing, authentication)                │
└───────────────────────────┬──────────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────┐
│                    Partition Layer                            │
│         (Consistent hashing, request routing)                 │
└───────────────────────────┬──────────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────┐
│                    Replication Layer                          │
│         (Raft consensus for each partition)                   │
└───────────────────────────┬──────────────────────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────┐
│                    Storage Layer                              │
│            (LSM-tree or B+tree storage engine)                │
└──────────────────────────────────────────────────────────────┘
```

**Features**:

1. **Data Model**
   - Key-value with versioning
   - TTL support
   - Multi-key transactions

2. **Consistency Levels**
   - Strong consistency (Raft)
   - Eventual consistency (async replication)
   - Read-your-writes
   - Monotonic reads

3. **Partitioning**
   - Consistent hashing
   - Virtual nodes
   - Online rebalancing
   - Range queries

4. **Replication**
   - Raft consensus per partition
   - Configurable replication factor
   - Automatic failover
   - Leader leasing

5. **Transactions**
   - Multi-key transactions
   - Optimistic concurrency control
   - Snapshot isolation
   - Deadlock detection

6. **Advanced Features**
   - Secondary indexes
   - Change data capture
   - Backup and restore
   - Metrics and tracing

7. **Testing**
   - Jepsen-style testing
   - Chaos engineering
   - Linearizability checker
   - Performance benchmarks

---

## ✅ Progress Tracker

| Phase | Week | Topic | Complete |
|-------|------|-------|----------|
| 1 | 1 | Distributed Systems Fundamentals | [ ] |
| 1 | 2 | Time and Ordering | [ ] |
| 1 | 3 | Distributed Algorithms | [ ] |
| 1 | 4 | Sharding and Partitioning | [ ] |
| 1 | 5 | Formal Semantics | [ ] |
| 1 | 6 | Case Studies | [ ] |
| 2 | 7 | Consensus Fundamentals | [ ] |
| 2 | 8 | Raft | [ ] |
| 2 | 9 | Paxos | [ ] |
| 2 | 10 | BFT | [ ] |
| 2 | 11 | Coordination | [ ] |
| 2 | 12 | SMR | [ ] |
| 3 | 13 | Replication | [ ] |
| 3 | 14 | CRDTs | [ ] |
| 3 | 15 | Event Sourcing | [ ] |
| 3 | 16 | Distributed Transactions | [ ] |
| 3 | 17 | Kafka | [ ] |
| 3 | 18 | Database Internals | [ ] |
| 4 | 19 | Formal Verification | [ ] |
| 4 | 20 | System Design | [ ] |
| 4 | 21 | Performance | [ ] |
| 4 | 22 | Security | [ ] |
| 4 | 23 | Research | [ ] |
| 4 | 24 | Project | [ ] |

---

*This learning path prepares you for research and advanced engineering roles in distributed systems. The capstone project demonstrates mastery of consensus, replication, and distributed data management.*
