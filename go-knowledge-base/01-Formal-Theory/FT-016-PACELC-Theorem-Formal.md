# FT-016: PACELC Theorem - Formal Specification

## Overview

PACELC is an extension of the CAP theorem that unifies the trade-off between consistency and latency. It states that in a distributed system:

- **If** there is a **P**artition, one must choose between **A**vailability and **C**onsistency
- **E**lse (when the system is running normally), one must choose between **L**atency and **C**onsistency

This theorem provides a more complete framework for understanding consistency trade-offs in distributed systems.

## Theoretical Foundations

### 1.1 System Model

**Distributed Data Store**:

```
System D = ⟨N, D, R, C, L, P⟩

where:
- N = {n₁, n₂, ..., nₙ}: set of nodes
- D: data items with replication factor r
- R: set of replicas for each data item
- C: consistency protocol
- L: latency function L: Operations → Time
- P: partition detection mechanism
```

**Network Model**:

```
Network state at time t:
- Connected: ∀nᵢ, nⱼ ∈ N: reachable(nᵢ, nⱼ, t)
- Partitioned: ∃N₁, N₂ ⊂ N: N₁ ∪ N₂ = N ∧ N₁ ∩ N₂ = ∅
               ∧ ∀nᵢ ∈ N₁, nⱼ ∈ N₂: ¬reachable(nᵢ, nⱼ, t)
```

### 1.2 PACELC Formal Statement

```
Theorem (PACELC):
For any replicated distributed data store D:

P ∧ A ⇒ ¬C  (If partitioned and available, not consistent)
P ∧ C ⇒ ¬A  (If partitioned and consistent, not available)
¬P ∧ L ⇒ ¬C (If not partitioned and low latency, not consistent)
¬P ∧ C ⇒ ¬L (If not partitioned and consistent, not low latency)

Where:
- P: Network partition exists
- A: System remains available
- C: System maintains consistency
- L: System maintains low latency
```

### 1.3 Proof of PACELC

**Part 1: Partition Scenario (PAC)**

```
Lemma 1 (Partition Consistency Trade-off):
Given a partition dividing N into N₁ and N₂, a write to N₁ cannot
be observed by reads in N₂ without cross-partition communication.

Proof:
1. Assume partition separates N₁ and N₂
2. Consider write w to data item d in N₁
3. For a read r of d in N₂ to observe w:
   - Synchronization must occur between N₁ and N₂
4. But partition prevents communication
5. Therefore: either w is not observed (¬C) or r waits (¬A)
∎

Theorem 1 (P ⇒ (A XOR C)):
During a partition, availability and consistency are mutually exclusive.

Proof:
Case 1: Choose Availability
  - Both N₁ and N₂ accept writes
  - Divergent histories possible
  - ⇒ ¬C (no global consistency)

Case 2: Choose Consistency
  - One partition must reject writes to prevent divergence
  - ⇒ ¬A (not fully available)
∎
```

**Part 2: Normal Operation (ELC)**

```
Lemma 2 (Latency vs. Consistency):
Achieving strong consistency requires coordination among replicas.

Proof:
1. Consider strong consistency: all reads see latest write
2. Read must contact sufficient replicas to find latest version
3. This requires:
   - Multiple network round-trips, OR
   - Synchronous replication to all replicas
4. Both approaches increase latency
5. Therefore: C ⇒ coordination overhead ⇒ ¬L
∎

Theorem 2 (¬P ⇒ (L XOR C)):
During normal operation, low latency and strong consistency are in tension.

Proof:
Case 1: Choose Low Latency
  - Read from single replica
  - No coordination overhead
  - May read stale data
  - ⇒ ¬C (weak consistency)

Case 2: Choose Consistency
  - Contact multiple replicas or use consensus
  - Coordination adds latency
  - ⇒ ¬L (higher latency)
∎
```

### 1.4 Formal Characterization of Systems

```
System Classification:

┌──────────────┬─────────────────────────────────────────────────┐
│ System       │ PACELC Choice                                    │
├──────────────┼─────────────────────────────────────────────────┤
│ Dynamo/Voldemort │ PA/EL (Available, Low Latency, Eventual)   │
│ Cassandra    │ PA/EL or PA/EC (configurable)                  │
│ MongoDB      │ PA/EC or PC/EC (configurable)                  │
│ HBase        │ PC/EC (Consistent, not low latency)            │
│ BigTable     │ PC/EC (Consistent, not low latency)            │
│ Spanner      │ PC/EC (Strong consistency, synchronized clocks)│
│ CockroachDB  │ PC/EC (Serializable default)                   │
│ PNUTS        │ PA/EL (Per-record timeline consistency)        │
└──────────────┴─────────────────────────────────────────────────┘
```

## TLA+ Specification

```tla
--------------------------- MODULE PACELC ---------------------------
EXTENDS Integers, Sequences, FiniteSets, Reals

CONSTANTS Nodes,          \* Set of nodes
          DataItems,      \* Set of data items
          Clients,        \* Set of clients
          MaxLatency      \* Maximum acceptable latency

VARIABLES nodeState,      \* State of each node
          replicas,       \* Replica placement
          operations,     \* Operation log
          networkState,   \* Current network state
          metrics         \* Performance metrics

\* Type definitions
ReplicaState == [data: [DataItems → Values ∪ {⊥}],
                 version: [DataItems → Nat],
                 timestamp: [DataItems → Nat]]

NetworkStatus == {"connected", "partitioned"}

Partition == [left: SUBSET Nodes, right: SUBSET Nodes]

\* Initial state
Init ==
  ∧ nodeState = [n ∈ Nodes ↦ [data ↦ [d ∈ DataItems ↦ ⊥],
                               version ↦ [d ∈ DataItems ↦ 0],
                               timestamp ↦ [d ∈ DataItems ↦ 0]]]
  ∧ replicas = [d ∈ DataItems ↦ CHOOSE s ∈ SUBSET Nodes : Cardinality(s) ≥ 3]
  ∧ operations = ⟨⟩
  ∧ networkState = [status ↦ "connected", partition ↦ ⊥]
  ∧ metrics = [latency ↦ 0, availability ↦ 1.0, consistency ↦ 1.0]

\* Network partition occurs
NetworkPartition ==
  ∧ networkState.status = "connected"
  ∧ ∃N₁, N₂ ∈ SUBSET Nodes:
      ∧ N₁ ∪ N₂ = Nodes
      ∧ N₁ ∩ N₂ = {}
      ∧ Cardinality(N₁) > 0
      ∧ Cardinality(N₂) > 0
      ∧ networkState' = [status ↦ "partitioned",
                         partition ↦ [left ↦ N₁, right ↦ N₂]]
  ∧ UNCHANGED ⟨nodeState, replicas, operations, metrics⟩

\* Network heals
NetworkHeal ==
  ∧ networkState.status = "partitioned"
  ∧ networkState' = [status ↦ "connected", partition ↦ ⊥]
  ∧ UNCHANGED ⟨nodeState, replicas, operations, metrics⟩

\* Read operation with consistency level
Read(client, dataItem, consistencyLevel) ==
  LET replicaSet == replicas[dataItem]
  IN CASE consistencyLevel = "ONE" →
          \* Low latency, weak consistency
          ∃n ∈ replicaSet:
            operations' = Append(operations,
              [type ↦ "read", client ↦ client, data ↦ dataItem,
               result ↦ nodeState[n].data[dataItem],
               latency ↦ 1,
               consistency ↦ "eventual"])
       [] consistencyLevel = "QUORUM" →
          \* Higher latency, stronger consistency
          ∃responses ∈ QuorumRead(replicaSet, dataItem):
            LET latest == MaxVersion(responses)
            IN operations' = Append(operations,
                 [type ↦ "read", client ↦ client, data ↦ dataItem,
                  result ↦ latest.value,
                  latency ↦ 3,
                  consistency ↦ "quorum"])
       [] consistencyLevel = "ALL" →
          \* Highest latency, strongest consistency
          ∃responses ∈ AllRead(replicaSet, dataItem):
            LET latest == MaxVersion(responses)
            IN operations' = Append(operations,
                 [type ↦ "read", client ↦ client, data ↦ dataItem,
                  result ↦ latest.value,
                  latency ↦ 5,
                  consistency ↦ "strong"])

\* Write operation with handling of partitions
Write(client, dataItem, value) ==
  LET replicaSet == replicas[dataItem]
  IN IF networkState.status = "connected"
     THEN \* Normal operation: write to all replicas
          nodeState' = [n ∈ Nodes ↦
            IF n ∈ replicaSet
            THEN [nodeState[n] EXCEPT
                   !.data[dataItem] = value,
                   !.version[dataItem] = @ + 1,
                   !.timestamp[dataItem] = @ + 1]
            ELSE nodeState[n]]
          operations' = Append(operations,
            [type ↦ "write", client ↦ client, data ↦ dataItem,
             value ↦ value, status ↦ "committed"])
     ELSE \* Partitioned: must choose availability or consistency
          IF ChooseAvailability  \* PA system
          THEN \* Write to available partition
               nodeState' = [n ∈ Nodes ↦
                 IF n ∈ AvailablePartition(networkState.partition)
                 THEN [nodeState[n] EXCEPT
                        !.data[dataItem] = value,
                        !.version[dataItem] = @ + 1]
                 ELSE nodeState[n]]
               operations' = Append(operations,
                 [type ↦ "write", client ↦ client, data ↦ dataItem,
                  value ↦ value, status ↦ "deferred"])
          ELSE \* PC system - reject write
               operations' = Append(operations,
                 [type ↦ "write", client ↦ client, data ↦ dataItem,
                  value ↦ value, status ↦ "rejected"])
               UNCHANGED nodeState

\* PACELC Theorem invariants

\* If partitioned and available, not consistent
PAC_Invariant ==
  networkState.status = "partitioned"
    ⇒ (AvailabilityMetric() > 0.5) ⇒ (ConsistencyMetric() < 1.0)

\* If not partitioned and low latency, not consistent
ELC_Invariant ==
  networkState.status = "connected"
    ⇒ (LatencyMetric() < MaxLatency) ⇒ (ConsistencyMetric() < 1.0)

\* Metrics calculations
AvailabilityMetric() ==
  LET total == Len(operations)
      successful == Cardinality({i ∈ 1..total :
                     operations[i].status ∈ {"committed", "deferred"}})
  IN IF total = 0 THEN 1.0 ELSE successful / total

ConsistencyMetric() ==
  LET reads == {i ∈ 1..Len(operations) : operations[i].type = "read"}
      consistentReads == Cardinality({i ∈ reads :
        operations[i].consistency ≠ "eventual"})
  IN IF Cardinality(reads) = 0 THEN 1.0
     ELSE consistentReads / Cardinality(reads)

LatencyMetric() ==
  LET totalLatency == Sum({operations[i].latency :
                     i ∈ 1..Len(operations)})
  IN IF Len(operations) = 0 THEN 0
     ELSE totalLatency / Len(operations)

=============================================================================
```

## Algorithm Pseudocode

### Adaptive Consistency Algorithm

```
Algorithm: Adaptive Consistency Based on Network Conditions

Constants:
  LATENCY_THRESHOLD: maximum acceptable latency
  CONSISTENCY_LEVELS: [ONE, QUORUM, ALL]
  PARTITION_TIMEOUT: time to detect partition

Variables:
  currentConsistency: current consistency level
  networkStatus: CONNECTED or PARTITIONED
  latencyHistory: circular buffer of recent latencies
  partitionDetector: failure detection module

Procedure Initialize():
  currentConsistency ← QUORUM
  networkStatus ← CONNECTED
  latencyHistory ← empty circular buffer of size 100
  StartPartitionDetector()

On NetworkStatusChange(newStatus):
  networkStatus ← newStatus

  if newStatus = PARTITIONED:
    // PAC scenario: choose between A and C
    if SystemPolicy = "AVAILABILITY_FIRST":
      currentConsistency ← ONE
      EnableHintedHandoff()
    else:
      currentConsistency ← ALL
      EnableStrictQuorum()
  else:
    // ELC scenario: resume normal operation
    currentConsistency ← QUORUM
    ResolveDeferredWrites()

On ReadRequest(key):
  // Check latency constraint
  if Average(latencyHistory) > LATENCY_THRESHOLD:
    // Degrade consistency for latency
    if currentConsistency = ALL:
      currentConsistency ← QUORUM
    else if currentConsistency = QUORUM:
      currentConsistency ← ONE

  // Perform read
  replicas ← GetReplicas(key)

  switch currentConsistency:
    case ONE:
      return ReadFromOne(replicas)
    case QUORUM:
      return ReadFromQuorum(replicas)
    case ALL:
      return ReadFromAll(replicas)

On WriteRequest(key, value):
  replicas ← GetReplicas(key)

  if networkStatus = PARTITIONED:
    if SystemPolicy = "AVAILABILITY_FIRST":
      // PA choice: write to available replicas
      available ← Filter(replicas, IsReachable)
      if Length(available) > 0:
        WriteToReplicas(available, key, value)
        QueueForReplay(key, value, unavailableReplicas)
        return SUCCESS
      else:
        return FAILURE
    else:
      // PC choice: require strict quorum
      if Length(Filter(replicas, IsReachable)) > Length(replicas) / 2:
        WriteToAllReachable(replicas, key, value)
        return SUCCESS
      else:
        return FAILURE
  else:
    // Normal operation: write to quorum
    return WriteToQuorum(replicas, key, value)

Procedure ResolveDeferredWrites():
  // After partition heals, reconcile divergent histories
  deferredWrites ← GetDeferredWrites()

  for each write in deferredWrites:
    if HasConflict(write):
      resolution ← ResolveConflict(write)
      ApplyResolution(resolution)
    else:
      PropagateWrite(write)

Procedure HasConflict(write):
  // Check if write conflicts with existing data
  currentValue ← ReadFromAll(GetReplicas(write.key))
  return write.timestamp < currentValue.timestamp

Procedure ResolveConflict(conflictingWrites):
  // Apply conflict resolution strategy
  switch ConflictResolutionPolicy:
    case LAST_WRITE_WINS:
      return MaxBy(conflictingWrites, w → w.timestamp)
    case VECTOR_CLOCKS:
      return MergeVectorClocks(conflictingWrites)
    case APPLICATION_MERGE:
      return InvokeMergeHook(conflictingWrites)
```

### Latency-Aware Routing

```
Algorithm: Latency-Aware Request Routing

Data Structures:
  NodeMetrics {
    nodeId: string
    averageLatency: float
    failureRate: float
    loadFactor: float
    lastUpdated: timestamp
  }

  RoutingTable {
    metrics: Map<Node, NodeMetrics>
    consistencyLatency: Map<ConsistencyLevel, float>
  }

Procedure SelectReplicas(key, consistencyLevel, maxLatency):
  candidates ← GetReplicas(key)

  // Filter by health
  healthy ← Filter(candidates, n →
    routingTable.metrics[n].failureRate < 0.01)

  // Estimate latency for each consistency level
  estimates ← Map(consistencyLevels, level → {
    level: level,
    estimatedLatency: EstimateLatency(healthy, level),
    consistency: GetConsistencyStrength(level)
  }))

  // Find best match: highest consistency under latency budget
  validOptions ← Filter(estimates, e →
    e.estimatedLatency ≤ maxLatency)

  if validOptions is empty:
    // Relax latency constraint
    return SelectByConsistency(candidates, consistencyLevel)

  // Select option with best consistency
  return MaxBy(validOptions, o → o.consistency)

Function EstimateLatency(nodes, consistencyLevel):
  baseLatency ← Min({routingTable.metrics[n].averageLatency : n ∈ nodes})

  switch consistencyLevel:
    case ONE:
      return baseLatency
    case QUORUM:
      // Need responses from majority
      sorted ← SortByLatency(nodes)
      quorumSize ← Length(nodes) / 2 + 1
      return sorted[quorumSize - 1].averageLatency
    case ALL:
      // Wait for slowest
      return Max({routingTable.metrics[n].averageLatency : n ∈ nodes})

Procedure UpdateMetrics(node, latency, success):
  metrics ← routingTable.metrics[node]

  // Exponential moving average for latency
  α ← 0.3
  metrics.averageLatency ← α * latency + (1 - α) * metrics.averageLatency

  // Update failure rate
  if success:
    metrics.failureRate ← 0.99 * metrics.failureRate
  else:
    metrics.failureRate ← 0.99 * metrics.failureRate + 0.01

  metrics.lastUpdated ← Now()
```

## Go Implementation

```go
// Package pacelc implements PACELC theorem-based adaptive consistency
package pacelc

import (
 "context"
 "fmt"
 "math"
 "sync"
 "time"
)

// ConsistencyLevel represents the level of consistency
type ConsistencyLevel int

const (
 // ONE reads from single replica (lowest latency)
 ONE ConsistencyLevel = iota
 // QUORUM reads from majority (balanced)
 QUORUM
 // ALL reads from all replicas (strongest)
 ALL
)

func (c ConsistencyLevel) String() string {
 switch c {
 case ONE:
  return "ONE"
 case QUORUM:
  return "QUORUM"
 case ALL:
  return "ALL"
 default:
  return "UNKNOWN"
 }
}

// SystemPolicy determines PACELC trade-offs
type SystemPolicy int

const (
 // AvailabilityFirst chooses PA during partitions
 AvailabilityFirst SystemPolicy = iota
 // ConsistencyFirst chooses PC during partitions
 ConsistencyFirst
)

// NetworkStatus represents the current network state
type NetworkStatus int

const (
 // Connected means normal operation
 Connected NetworkStatus = iota
 // Partitioned means network partition detected
 Partitioned
)

// NodeMetrics tracks performance metrics for a node
type NodeMetrics struct {
 NodeID         string
 AverageLatency time.Duration
 FailureRate    float64
 LoadFactor     float64
 LastUpdated    time.Time
 mu             sync.RWMutex
}

// UpdateLatency updates the average latency using EMA
func (m *NodeMetrics) UpdateLatency(latency time.Duration) {
 m.mu.Lock()
 defer m.mu.Unlock()

 alpha := 0.3
 current := float64(m.AverageLatency)
 new := float64(latency)
 m.AverageLatency = time.Duration(alpha*new + (1-alpha)*current)
 m.LastUpdated = time.Now()
}

// RecordFailure records a failure
func (m *NodeMetrics) RecordFailure() {
 m.mu.Lock()
 defer m.mu.Unlock()
 m.FailureRate = 0.99*m.FailureRate + 0.01
}

// RecordSuccess records a success
func (m *NodeMetrics) RecordSuccess() {
 m.mu.Lock()
 defer m.mu.Unlock()
 m.FailureRate = 0.99 * m.FailureRate
}

// RoutingTable maintains node metrics and routing decisions
type RoutingTable struct {
 mu              sync.RWMutex
 metrics         map[string]*NodeMetrics
 latencyHistory  *CircularBuffer
 maxLatency      time.Duration
 defaultConsistency ConsistencyLevel
}

// CircularBuffer implements a fixed-size circular buffer
type CircularBuffer struct {
 data []time.Duration
 head int
 size int
 mu   sync.Mutex
}

// NewCircularBuffer creates a new circular buffer
func NewCircularBuffer(size int) *CircularBuffer {
 return &CircularBuffer{
  data: make([]time.Duration, size),
  size: size,
 }
}

// Add adds a value to the buffer
func (cb *CircularBuffer) Add(value time.Duration) {
 cb.mu.Lock()
 defer cb.mu.Unlock()
 cb.data[cb.head] = value
 cb.head = (cb.head + 1) % cb.size
}

// Average returns the average of all values
func (cb *CircularBuffer) Average() time.Duration {
 cb.mu.Lock()
 defer cb.mu.Unlock()

 if cb.head == 0 {
  return 0
 }

 var sum time.Duration
 count := 0
 for i := 0; i < cb.size && i < cb.head; i++ {
  if cb.data[i] > 0 {
   sum += cb.data[i]
   count++
  }
 }

 if count == 0 {
  return 0
 }
 return sum / time.Duration(count)
}

// PACELCStore implements an adaptive consistency data store
type PACELCStore struct {
 mu               sync.RWMutex
 nodes            map[string]*Node
 replicaSets      map[string][]string // key -> node IDs
 networkStatus    NetworkStatus
 policy           SystemPolicy
 routingTable     *RoutingTable
 partitionDetector *PartitionDetector
 conflictResolver ConflictResolver
}

// Node represents a storage node
type Node struct {
 ID      string
 Address string
 Data    map[string]*DataItem
 mu      sync.RWMutex
}

// DataItem represents a stored data item
type DataItem struct {
 Key       string
 Value     []byte
 Version   uint64
 Timestamp time.Time
 VectorClock map[string]uint64
}

// PartitionDetector detects network partitions
type PartitionDetector struct {
 heartbeatInterval time.Duration
 heartbeatTimeout  time.Duration
 suspectedNodes    map[string]time.Time
 mu                sync.Mutex
}

// ConflictResolver handles divergent histories
type ConflictResolver interface {
 Resolve(writes []*WriteOperation) *DataItem
}

// WriteOperation represents a write request
type WriteOperation struct {
 Key         string
 Value       []byte
 Timestamp   time.Time
 NodeID      string
 VectorClock map[string]uint64
}

// NewPACELCStore creates a new PACELC store
func NewPACELCStore(policy SystemPolicy, maxLatency time.Duration) *PACELCStore {
 return &PACELCStore{
  nodes:         make(map[string]*Node),
  replicaSets:   make(map[string][]string),
  networkStatus: Connected,
  policy:        policy,
  routingTable: &RoutingTable{
   metrics:            make(map[string]*NodeMetrics),
   latencyHistory:     NewCircularBuffer(100),
   maxLatency:         maxLatency,
   defaultConsistency: QUORUM,
  },
  partitionDetector: &PartitionDetector{
   heartbeatInterval: 1 * time.Second,
   heartbeatTimeout:  5 * time.Second,
   suspectedNodes:    make(map[string]time.Time),
  },
 }
}

// Get performs a read with adaptive consistency
func (s *PACELCStore) Get(ctx context.Context, key string) (*DataItem, error) {
 start := time.Now()
 defer func() {
  s.routingTable.latencyHistory.Add(time.Since(start))
 }()

 // Determine consistency level based on current conditions
 consistency := s.selectConsistencyLevel()

 // Get replica set
 replicas := s.getReplicas(key)
 if len(replicas) == 0 {
  return nil, fmt.Errorf("no replicas for key %s", key)
 }

 // Select nodes based on consistency level
 selectedNodes := s.selectNodesForRead(replicas, consistency)

 // Perform read
 results := make(chan *DataItem, len(selectedNodes))
 errors := make(chan error, len(selectedNodes))

 for _, nodeID := range selectedNodes {
  go func(id string) {
   node, ok := s.nodes[id]
   if !ok {
    errors <- fmt.Errorf("node %s not found", id)
    return
   }

   item := s.readFromNode(node, key)
   if item != nil {
    results <- item
   }
  }(nodeID)
 }

 // Collect results based on consistency level
 var items []*DataItem
 timeout := time.NewTimer(s.routingTable.maxLatency)
 defer timeout.Stop()

 required := s.requiredResponses(len(selectedNodes), consistency)

 for len(items) < required {
  select {
  case item := <-results:
   items = append(items, item)
  case <-timeout.C:
   return nil, fmt.Errorf("read timeout")
  case <-ctx.Done():
   return nil, ctx.Err()
  }
 }

 // Return item with highest version
 return s.selectLatest(items), nil
}

// Put performs a write with PACELC handling
func (s *PACELCStore) Put(ctx context.Context, key string, value []byte) error {
 s.mu.RLock()
 status := s.networkStatus
 policy := s.policy
 s.mu.RUnlock()

 replicas := s.getReplicas(key)

 if status == Partitioned {
  // PAC scenario: choose between A and C
  if policy == AvailabilityFirst {
   return s.writeAvailable(ctx, replicas, key, value)
  }
  return s.writeConsistent(ctx, replicas, key, value)
 }

 // ELC scenario: normal quorum write
 return s.writeQuorum(ctx, replicas, key, value)
}

// selectConsistencyLevel chooses consistency based on latency
func (s *PACELCStore) selectConsistencyLevel() ConsistencyLevel {
 avgLatency := s.routingTable.latencyHistory.Average()
 maxLatency := s.routingTable.maxLatency

 if avgLatency > maxLatency {
  // Degrade consistency for latency
  current := s.routingTable.defaultConsistency
  if current > ONE {
   return current - 1
  }
 }

 return s.routingTable.defaultConsistency
}

// writeAvailable implements PA choice during partition
func (s *PACELCStore) writeAvailable(ctx context.Context, replicas []string, key string, value []byte) error {
 available := s.getAvailableNodes(replicas)

 if len(available) == 0 {
  return fmt.Errorf("no available nodes during partition")
 }

 // Write to available nodes
 var wg sync.WaitGroup
 errors := make(chan error, len(available))

 for _, nodeID := range available {
  wg.Add(1)
  go func(id string) {
   defer wg.Done()
   if err := s.writeToNode(id, key, value); err != nil {
    errors <- err
   }
  }(nodeID)
 }

 wg.Wait()
 close(errors)

 // Queue for later replication to unavailable nodes
 unavailable := s.getUnavailableNodes(replicas)
 s.queueDeferredWrite(key, value, unavailable)

 return nil
}

// writeConsistent implements PC choice during partition
func (s *PACELCStore) writeConsistent(ctx context.Context, replicas []string, key string, value []byte) error {
 available := len(s.getAvailableNodes(replicas))
 required := len(replicas)/2 + 1

 if available < required {
  return fmt.Errorf("insufficient nodes for consistent write during partition")
 }

 return s.writeQuorum(ctx, replicas, key, value)
}

// writeQuorum performs a quorum write
func (s *PACELCStore) writeQuorum(ctx context.Context, replicas []string, key string, value []byte) error {
 quorumSize := len(replicas)/2 + 1
 acks := 0

 var wg sync.WaitGroup
 ackChan := make(chan bool, len(replicas))

 for _, nodeID := range replicas {
  wg.Add(1)
  go func(id string) {
   defer wg.Done()
   if err := s.writeToNode(id, key, value); err == nil {
    ackChan <- true
   }
  }(nodeID)
 }

 go func() {
  wg.Wait()
  close(ackChan)
 }()

 for range ackChan {
  acks++
  if acks >= quorumSize {
   return nil
  }
 }

 return fmt.Errorf("failed to achieve quorum (%d/%d)", acks, quorumSize)
}

// Helper methods

func (s *PACELCStore) getReplicas(key string) []string {
 s.mu.RLock()
 defer s.mu.RUnlock()
 return s.replicaSets[key]
}

func (s *PACELCStore) selectNodesForRead(replicas []string, level ConsistencyLevel) []string {
 // Sort by latency and select based on consistency level
 s.routingTable.mu.RLock()
 defer s.routingTable.mu.RUnlock()

 sorted := make([]string, len(replicas))
 copy(sorted, replicas)

 // Simple selection - in production, sort by actual latency
 switch level {
 case ONE:
  if len(sorted) > 0 {
   return []string{sorted[0]}
  }
 case QUORUM:
  quorumSize := len(sorted)/2 + 1
  if len(sorted) >= quorumSize {
   return sorted[:quorumSize]
  }
 case ALL:
  return sorted
 }

 return sorted
}

func (s *PACELCStore) requiredResponses(total int, level ConsistencyLevel) int {
 switch level {
 case ONE:
  return 1
 case QUORUM:
  return total/2 + 1
 case ALL:
  return total
 }
 return 1
}

func (s *PACELCStore) readFromNode(node *Node, key string) *DataItem {
 node.mu.RLock()
 defer node.mu.RUnlock()
 return node.Data[key]
}

func (s *PACELCStore) writeToNode(nodeID string, key string, value []byte) error {
 s.mu.RLock()
 node, ok := s.nodes[nodeID]
 s.mu.RUnlock()

 if !ok {
  return fmt.Errorf("node not found")
 }

 node.mu.Lock()
 defer node.mu.Unlock()

 if node.Data == nil {
  node.Data = make(map[string]*DataItem)
 }

 existing := node.Data[key]
 version := uint64(1)
 if existing != nil {
  version = existing.Version + 1
 }

 node.Data[key] = &DataItem{
  Key:       key,
  Value:     value,
  Version:   version,
  Timestamp: time.Now(),
 }

 return nil
}

func (s *PACELCStore) selectLatest(items []*DataItem) *DataItem {
 if len(items) == 0 {
  return nil
 }

 latest := items[0]
 for _, item := range items {
  if item.Version > latest.Version {
   latest = item
  }
 }

 return latest
}

func (s *PACELCStore) getAvailableNodes(replicas []string) []string {
 // In production: check partition detector
 return replicas
}

func (s *PACELCStore) getUnavailableNodes(replicas []string) []string {
 // In production: check partition detector
 return []string{}
}

func (s *PACELCStore) queueDeferredWrite(key string, value []byte, nodes []string) {
 // In production: add to deferred write queue
}
```

## System Comparison Matrix

| System | P Scenario | E Scenario | Replica Factor | Conflict Resolution |
|--------|-----------|-----------|----------------|---------------------|
| **Dynamo** | PA (hinted handoff) | EL (async replication) | Configurable | Vector clocks + application merge |
| **Cassandra** | Configurable | Configurable | Configurable | Last-write-wins or custom |
| **Voldemort** | PA | EL | Fixed | Vector clocks |
| **Riak** | PA | Configurable | Configurable | Vector clocks + CRDTs |
| **MongoDB** | PA/PC | EC | Replica sets | Primary-based |
| **CockroachDB** | PC | EC | Default 3 | Serializable transactions |
| **Spanner** | PC | EC | Configurable | TrueTime + 2PC |
| **Cosmos DB** | Configurable | Configurable | Configurable | Multiple models |

## Visual Representations

### Figure 1: PACELC Decision Tree

```
┌─────────────────────────────────────────────────────────────────┐
│                  PACELC DECISION TREE                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                    ┌─────────────┐                              │
│                    │   START     │                              │
│                    └──────┬──────┘                              │
│                           │                                     │
│                    Is there a                                   │
│                    partition?                                   │
│                           │                                     │
│              ┌────────────┴────────────┐                        │
│              ▼                         ▼                        │
│            YES                        NO                        │
│              │                         │                        │
│     ┌────────┴────────┐       ┌────────┴────────┐              │
│     ▼                 ▼       ▼                 ▼              │
│  ┌──────┐        ┌──────┐  ┌──────┐        ┌──────┐            │
│  │ PA   │        │ PC   │  │ EL   │        │ EC   │            │
│  │      │        │      │  │      │        │      │            │
│  │Prefer│        │Prefer│  │Prefer│        │Prefer│            │
│  │Avail-│        │Consis│  │Low   │        │Consis│            │
│  │ability│       │tency │  │Latency│       │tency │            │
│  └──────┘        └──────┘  └──────┘        └──────┘            │
│     │               │         │               │                │
│     ▼               ▼         ▼               ▼                │
│  ┌──────┐      ┌──────┐   ┌──────┐      ┌──────┐              │
│  │Dynamo│      │HBase │   │Dynamo│      │Spanner│             │
│  │Cassan│      │Spanner│  │Volden│      │Cockroach│           │
│  │dra   │      │      │   │mort  │      │DB     │             │
│  └──────┘      └──────┘   └──────┘      └──────┘              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 2: Latency vs Consistency Trade-off

```
┌─────────────────────────────────────────────────────────────────┐
│              LATENCY VS CONSISTENCY TRADE-OFF                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Latency                                                        │
│    ▲                                                            │
│    │  ┌────┐                                                    │
│  H │  │ALL │                                                    │
│  I │  │    │     ┌────┐                                         │
│  G │  │    │     │Q.U.│                                         │
│  H │  │    │     │O.R.│                                         │
│    │  │    │     │U.M │     ┌────┐                              │
│    │  │    │     │    │     │ONE │                              │
│  L │  │    │     │    │     │    │                              │
│  O │  │    │     │    │     │    │                              │
│  W │  │    │     │    │     │    │                              │
│    └──┴────┴─────┴────┴─────┴────┴─────► Consistency            │
│       WEAK              MEDIUM          STRONG                  │
│                                                                 │
│  System Placement:                                              │
│  • Dynamo/Cassandra ONE: Low latency, eventual consistency      │
│  • Cassandra QUORUM: Balanced latency and consistency           │
│  • HBase/Spanner ALL: High latency, strong consistency          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 3: Network Partition Handling

```
┌─────────────────────────────────────────────────────────────────┐
│              NETWORK PARTITION HANDLING                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Normal Operation:                                              │
│  ┌─────────┐         ┌─────────┐         ┌─────────┐           │
│  │ Node A  │◄───────►│ Node B  │◄───────►│ Node C  │           │
│  │ (master)│         │(replica)│         │(replica)│           │
│  └─────────┘         └─────────┘         └─────────┘           │
│       │                   │                   │                 │
│       └───────────────────┴───────────────────┘                 │
│              All writes replicate to quorum                     │
│                                                                 │
│  Partition Occurs:                                              │
│  ┌─────────┐    X    ┌─────────┐         ┌─────────┐           │
│  │ Node A  │◄───╳───►│ Node B  │◄───────►│ Node C  │           │
│  │Partition│    X    │Partition│         │   OK    │           │
│  │    1    │    X    │    2    │         │         │           │
│  └─────────┘         └─────────┘         └─────────┘           │
│       │                   │                   │                 │
│                                                                 │
│  PA Choice (Dynamo):                                            │
│  ┌─────────┐              ┌─────────┐    ┌─────────┐           │
│  │ Node A  │              │ Node B  │◄──►│ Node C  │           │
│  │write(x,1│              │write(x,2│    │(replica)│           │
│  │)        │              │)        │    │         │           │
│  └─────────┘              └─────────┘    └─────────┘           │
│       │                        │                                │
│       ▼                        ▼                                │
│  [deferred]               [accepted]                            │
│  write queued                                                  │
│  for healing                                                    │
│                                                                 │
│  PC Choice (Spanner):                                           │
│  ┌─────────┐              ┌─────────┐    ┌─────────┐           │
│  │ Node A  │              │ Node B  │◄──►│ Node C  │           │
│  │REJECTED │              │write(x,2│    │(accept) │           │
│  │(no qrm) │              │)        │    │         │           │
│  └─────────┘              └─────────┘    └─────────┘           │
│       X                        │                                │
│  writes fail              writes succeed                        │
│  during partition         (majority side)                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## References

1. Abadi, D. J. (2012). "Consistency tradeoffs in modern distributed database system design." IEEE Computer.
2. Brewer, E. (2012). "CAP twelve years later: How the 'rules' have changed." IEEE Computer.
3. Vogels, W. (2009). "Eventually consistent." ACM Queue.
4. Bailis, P., & Ghodsi, A. (2013). "Eventual consistency today: Limitations, extensions, and beyond." ACM Queue.
5. DeCandia, G., et al. (2007). "Dynamo: Amazon's highly available key-value store." ACM SOSP.
