# FT-017: Quorum Consensus - Formal Specification

## Overview

Quorum consensus is a fundamental mechanism in distributed systems for achieving agreement among replicas without requiring all nodes to participate. A quorum is a minimum number of votes required to perform an operation, ensuring that any two operations on the same data item overlap in at least one replica.

## Theoretical Foundations

### 1.1 Quorum System Definition

**Quorum System**:

```
A quorum system Q over a universe U = {n₁, n₂, ..., nₙ} is a collection
of subsets of U (quorums) such that:

∀Q₁, Q₂ ∈ Q: Q₁ ∩ Q₂ ≠ ∅

This is the intersection property, ensuring any two quorums share at least one node.
```

**Strict Quorum System**:

```
A strict quorum system additionally satisfies:

∀Q₁, Q₂ ∈ Q, ∀B ⊂ U where |B| ≤ f:
  (Q₁ \\ B) ∩ (Q₂ \\ B) ≠ ∅

Where f is the maximum number of faulty nodes the system can tolerate.
```

### 1.2 Quorum Types

#### Majority Quorum

```
Q_majority = {Q ⊆ U : |Q| > |U|/2}

Properties:
- |Q₁ ∩ Q₂| ≥ 1 for any Q₁, Q₂ ∈ Q_majority
- Fault tolerance: f < |U|/2
- Quorum size: ⌊n/2⌋ + 1
```

**Proof of Intersection**:

```
Given Q₁, Q₂ ∈ Q_majority:
|Q₁| > n/2 and |Q₂| > n/2

Assume Q₁ ∩ Q₂ = ∅:
|Q₁ ∪ Q₂| = |Q₁| + |Q₂| > n/2 + n/2 = n

But Q₁ ∪ Q₂ ⊆ U, so |Q₁ ∪ Q₂| ≤ n
Contradiction! Therefore Q₁ ∩ Q₂ ≠ ∅.
∎
```

#### Read/Write Quorums (Weighted Voting)

```
Given:
- Read quorum requirement: R
- Write quorum requirement: W
- Total votes: V = Σ votes(nᵢ)

Constraints:
1. W > V/2 (write quorum is majority)
2. R + W > V (read and write quorums intersect)
3. W > f (write quorum survives f failures)

Read/Write Quorum System:
Q_read = {Q ⊆ U : Σ_{n∈Q} votes(n) ≥ R}
Q_write = {Q ⊆ U : Σ_{n∈Q} votes(n) ≥ W}
```

#### Grid Quorum

```
For U arranged in a √n × √n grid:

Q_grid = {rows} ∪ {columns}

Properties:
- Any row intersects any column
- Quorum size: 2√n - 1 (one row + one column - intersection)
- Better load balancing than majority
- Lower fault tolerance
```

#### Tree Quorum

```
For U arranged in a tree structure:

Q_tree = {root} ∪ {paths from root to leaves}

Properties:
- Root is in all quorums
- If root fails: use subtrees
- Adaptive quorum size based on failures
```

### 1.3 Quorum Consensus Protocol

**Core Operations**:

```
Write(value v):
  1. Select a write quorum Q_w ∈ Q_write
  2. Send v to all nodes in Q_w
  3. Wait for acknowledgments from quorum
  4. Return success

Read():
  1. Select a read quorum Q_r ∈ Q_read
  2. Query all nodes in Q_r for latest version
  3. Return value with maximum version
```

**Correctness Proof**:

```
Theorem: The quorum consensus protocol ensures linearizable reads and writes.

Proof:
1. Write-Write Ordering:
   Let w₁, w₂ be two writes with w₁ completing before w₂ starts.
   w₁ completes after receiving acks from Q_w₁.
   w₂ starts by selecting Q_w₂.
   Since Q_w₁ ∩ Q_w₂ ≠ ∅, at least one node has seen w₁ before w₂.
   Therefore, version ordering is preserved.

2. Read-Write Ordering:
   Let w be a completed write, and r be a subsequent read.
   w writes to Q_w, r reads from Q_r.
   Since Q_w ∩ Q_r ≠ ∅, r will observe w or a later write.

3. Read-Read Ordering:
   Two reads may return different versions if concurrent writes occur.
   However, each read returns the latest version from its quorum.
   Subsequent reads from intersecting quorums will observe monotonic versions.

Therefore, the protocol provides linearizability.
∎
```

## TLA+ Specification

```tla
------------------------ MODULE QuorumConsensus ------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Nodes,           \* Set of all nodes
          Values,          \* Set of possible values
          R,               \* Read quorum size
          W,               \* Write quorum size
          MaxVersion       \* Maximum version for model checking

ASSUME R + W > Cardinality(Nodes)
ASSUME W > Cardinality(Nodes) \div 2

VARIABLES nodeState,       \* State of each node
          operationLog,    \* Log of operations
          quorumHistory    \* History of quorum selections

\* Node state includes stored value and version
NodeState == [value: Values ∪ {⊥},
              version: 0..MaxVersion,
              lastWriter: Nodes ∪ {⊥}]

\* Quorum is a subset of Nodes meeting size requirements
IsReadQuorum(Q) == Q ⊆ Nodes ∧ Cardinality(Q) ≥ R
IsWriteQuorum(Q) == Q ⊆ Nodes ∧ Cardinality(Q) ≥ W

\* Initial state
Init ==
  ∧ nodeState = [n ∈ Nodes ↦ [value ↦ ⊥,
                               version ↦ 0,
                               lastWriter ↦ ⊥]]
  ∧ operationLog = ⟨⟩
  ∧ quorumHistory = ⟨⟩

\* Write operation using write quorum
Write(writer, value) ==
  ∃Q ∈ SUBSET Nodes:
    ∧ IsWriteQuorum(Q)
    ∧ nodeState' = [n ∈ Nodes ↦
         IF n ∈ Q
         THEN [value ↦ value,
               version ↦ nodeState[n].version + 1,
               lastWriter ↦ writer]
         ELSE nodeState[n]]
    ∧ operationLog' = Append(operationLog,
         [type ↦ "write",
          writer ↦ writer,
          value ↦ value,
          quorum ↦ Q])
    ∧ quorumHistory' = Append(quorumHistory, Q)

\* Read operation using read quorum
Read(reader) ==
  ∃Q ∈ SUBSET Nodes:
    ∧ IsReadQuorum(Q)
    ∧ LET responses == {nodeState[n] : n ∈ Q}
          maxVersion == MAX {r.version : r ∈ responses}
          result == CHOOSE r ∈ responses : r.version = maxVersion
      IN operationLog' = Append(operationLog,
           [type ↦ "read",
            reader ↦ reader,
            result ↦ result.value,
            version ↦ result.version,
            quorum ↦ Q])
    ∧ UNCHANGED nodeState
    ∧ quorumHistory' = Append(quorumHistory, Q)

\* Next state relation
Next ==
  ∨ ∃n ∈ Nodes, v ∈ Values: Write(n, v)
  ∨ ∃n ∈ Nodes: Read(n)

\* Invariants

\* Quorum intersection property
QuorumIntersection ==
  ∀i, j ∈ 1..Len(quorumHistory):
    quorumHistory[i] ∩ quorumHistory[j] ≠ {}

\* Any write quorum intersects with any read quorum
WriteReadIntersection ==
  ∀Qw ∈ {Q ∈ SUBSET Nodes : IsWriteQuorum(Q)},
   Qr ∈ {Q ∈ SUBSET Nodes : IsReadQuorum(Q)}:
    Qw ∩ Qr ≠ {}

\* Version monotonicity: writes increase versions
VersionMonotonicity ==
  ∀i, j ∈ 1..Len(operationLog):
    i < j ∧ operationLog[i].type = "write" ∧ operationLog[j].type = "write"
    ⇒ operationLog[i].quorum ∩ operationLog[j].quorum ≠ {}

\* Safety: reads see previously written values
ReadConsistency ==
  ∀readIdx ∈ 1..Len(operationLog):
    operationLog[readIdx].type = "read"
    ⇒ ∃writeIdx ∈ 1..readIdx:
        ∧ operationLog[writeIdx].type = "write"
        ∧ operationLog[readIdx].quorum ∩ operationLog[writeIdx].quorum ≠ {}

=============================================================================
```

## Algorithm Pseudocode

### Dynamic Quorum Consensus

```
Algorithm: Dynamic Quorum Consensus with Failure Detection

Constants:
  N: total number of nodes
  F: maximum number of tolerated failures
  DEFAULT_R: default read quorum size
  DEFAULT_W: default write quorum size

Variables:
  activeNodes: set of currently active nodes
  failedNodes: set of detected failed nodes
  nodeVotes: map of node to vote weight
  currentR: current read quorum threshold
  currentW: current write quorum threshold

Initialization:
  activeNodes ← all nodes
  failedNodes ← ∅
  nodeVotes ← {n → 1 : n ∈ activeNodes}
  currentR ← DEFAULT_R
  currentW ← DEFAULT_W

On FailureDetected(node):
  activeNodes ← activeNodes \\ {node}
  failedNodes ← failedNodes ∪ {node}

  // Recalculate quorum sizes
  RecalculateQuorums()

  // Trigger reconfiguration if necessary
  if |activeNodes| < currentW:
    EnterDegradedMode()

RecalculateQuorums():
  totalVotes ← Σ nodeVotes[n] for n ∈ activeNodes

  // Maintain constraints:
  // 1. W > totalVotes/2
  // 2. R + W > totalVotes
  // 3. W > F (if possible)

  currentW ← ⌊totalVotes/2⌋ + 1
  currentR ← totalVotes - currentW + 1

  // Ensure minimum availability
  if currentR > |activeNodes|:
    currentR ← |activeNodes|
  if currentW > |activeNodes|:
    currentW ← |activeNodes|

Write(value):
  // Select optimal write quorum
  Qw ← SelectQuorum(activeNodes, currentW, WRITE_OPTIMIZED)

  // Attempt write
  acks ← 0
  for node in Qw:
    if SendWrite(node, value):
      acks ← acks + nodeVotes[node]
    else:
      MarkSuspected(node)

  if acks ≥ currentW:
    return SUCCESS
  else:
    // Retry with different quorum
    remainingNodes ← activeNodes \\ Qw
    if |remainingNodes| ≥ currentW - acks:
      Qw2 ← SelectQuorum(remainingNodes, currentW - acks, ANY)
      for node in Qw2:
        if SendWrite(node, value):
          acks ← acks + nodeVotes[node]

      if acks ≥ currentW:
        return SUCCESS

    return FAILURE

Read():
  // Select optimal read quorum
  Qr ← SelectQuorum(activeNodes, currentR, READ_OPTIMIZED)

  // Collect responses
  responses ← ∅
  for node in Qr:
    response ← SendRead(node)
    if response ≠ timeout:
      responses ← responses ∪ {response}
    else:
      MarkSuspected(node)

  if |responses| < currentR:
    // Retry with different nodes
    return RetryRead()

  // Find value with maximum version
  maxVersion ← 0
  result ← ⊥
  for response in responses:
    if response.version > maxVersion:
      maxVersion ← response.version
      result ← response.value

  return result

SelectQuorum(candidates, required, optimization):
  switch optimization:
    case WRITE_OPTIMIZED:
      // Select lowest latency nodes for faster writes
      return LowestLatency(candidates, required)

    case READ_OPTIMIZED:
      // Select nodes with highest read throughput
      return HighestThroughput(candidates, required)

    case ANY:
      // Select any available nodes
      return RandomSelection(candidates, required)

EnterDegradedMode():
  // System cannot maintain strict quorums
  // Options:
  // 1. Accept reads only
  // 2. Enter split-brain mode with reconciliation
  // 3. Pause until nodes recover

  if DegradedModePolicy = "READ_ONLY":
    DisableWrites()
  else if DegradedModePolicy = "EVENTUAL":
    SwitchToEventualConsistency()
  else if DegradedModePolicy = "STALL":
    WaitForQuorumRecovery()

Correctness Proof:

Theorem: Dynamic Quorum Consensus maintains safety (intersection
property) despite up to F node failures.

Proof:
1. Initialization: R + W > N, W > N/2 ensures intersection.

2. Node Failure:
   When a node fails, RecalculateQuorums() maintains:
   - new_W > new_N/2 (where new_N = N - |failedNodes|)
   - new_R + new_W > new_N

   Therefore, any write quorum and read quorum still intersect.

3. Safety Property:
   Let Qw be a write quorum and Qr be a read quorum after recalculation.
   |Qw| ≥ new_W and |Qr| ≥ new_R
   |Qw| + |Qr| ≥ new_W + new_R > new_N = |activeNodes|

   By pigeonhole principle, Qw ∩ Qr ≠ ∅.

Therefore, the protocol remains safe despite failures.
∎
```

### Hierarchical Quorum

```
Algorithm: Hierarchical Quorum for Large-Scale Systems

Structure:
  - Nodes organized in hierarchy: regions → racks → nodes
  - Quorums must include nodes from multiple levels

Quorum Requirements:
  Level 1 (Region): R₁ nodes from R regions
  Level 2 (Rack):   R₂ nodes from each selected region
  Level 3 (Node):   R₃ nodes from each selected rack

Total Quorum Size: R₁ × R₂ × R₃

SelectHierarchicalQuorum():
  // Level 1: Select regions
  selectedRegions ← SelectRegions(R₁)

  quorum ← ∅
  for region in selectedRegions:
    // Level 2: Select racks within region
    selectedRacks ← SelectRacks(region, R₂)

    for rack in selectedRacks:
      // Level 3: Select nodes within rack
      selectedNodes ← SelectNodes(rack, R₃)
      quorum ← quorum ∪ selectedNodes

  return quorum

Why Hierarchical Quorums:
1. Locality: Prefer local communication
2. Failure Isolation: Rack/region failures don't break all quorums
3. Scalability: Smaller quorums in large deployments
4. Cost: Reduce cross-region traffic

Quorum Intersection Proof:

Theorem: Two hierarchical quorums Q₁ and Q₂ intersect.

Proof:
1. At level 1: Q₁ and Q₂ select R₁ regions each
   Total regions = R
   2×R₁ > R (by construction)
   Therefore, Q₁ and Q₂ share at least 2×R₁ - R regions

2. For each shared region:
   Q₁ and Q₂ select R₂ racks each
   Total racks in region = RacksPerRegion
   2×R₂ > RacksPerRegion
   Therefore, they share at least 2×R₂ - RacksPerRegion racks

3. For each shared rack:
   Q₁ and Q₂ select R₃ nodes each
   Total nodes in rack = NodesPerRack
   2×R₃ > NodesPerRack
   Therefore, they share at least 2×R₃ - NodesPerRack nodes

Thus, Q₁ ∩ Q₂ ≠ ∅.
∎
```

## Go Implementation

```go
// Package quorum implements quorum consensus protocols
package quorum

import (
 "context"
 "fmt"
 "math/rand"
 "sort"
 "sync"
 "time"
)

// Node represents a quorum member
type Node struct {
 ID       string
 Address  string
 Weight   int
 Latency  time.Duration
 IsActive bool
 mu       sync.RWMutex
}

// GetLatency returns the node's latency
func (n *Node) GetLatency() time.Duration {
 n.mu.RLock()
 defer n.mu.RUnlock()
 return n.Latency
}

// IsHealthy returns whether the node is active
func (n *Node) IsHealthy() bool {
 n.mu.RLock()
 defer n.mu.RUnlock()
 return n.IsActive
}

// QuorumType represents the type of quorum
type QuorumType int

const (
 ReadQuorum QuorumType = iota
 WriteQuorum
)

// QuorumManager manages quorum selection and validation
type QuorumManager struct {
 mu           sync.RWMutex
 nodes        map[string]*Node
 readQuorum   int
 writeQuorum  int
 totalWeight  int
 failureDetector *FailureDetector
}

// FailureDetector monitors node health
type FailureDetector struct {
 suspectedNodes map[string]time.Time
 mu             sync.Mutex
}

// NewQuorumManager creates a new quorum manager
func NewQuorumManager(nodes []*Node, readQ, writeQ int) (*QuorumManager, error) {
 if readQ+writeQ <= len(nodes) {
  return nil, fmt.Errorf("read quorum (%d) + write quorum (%d) must exceed node count (%d)",
   readQ, writeQ, len(nodes))
 }

 if writeQ <= len(nodes)/2 {
  return nil, fmt.Errorf("write quorum must be majority (> %d)", len(nodes)/2)
 }

 qm := &QuorumManager{
  nodes:       make(map[string]*Node),
  readQuorum:  readQ,
  writeQuorum: writeQ,
  failureDetector: &FailureDetector{
   suspectedNodes: make(map[string]time.Time),
  },
 }

 for _, node := range nodes {
  qm.nodes[node.ID] = node
  qm.totalWeight += node.Weight
 }

 return qm, nil
}

// SelectQuorum selects an optimal quorum of the specified type
func (qm *QuorumManager) SelectQuorum(qtype QuorumType, excluded map[string]bool) ([]*Node, error) {
 qm.mu.RLock()
 defer qm.mu.RUnlock()

 // Get available nodes
 var available []*Node
 for _, node := range qm.nodes {
  if node.IsHealthy() && !excluded[node.ID] {
   available = append(available, node)
  }
 }

 required := qm.readQuorum
 if qtype == WriteQuorum {
  required = qm.writeQuorum
 }

 if len(available) < required {
  return nil, fmt.Errorf("insufficient healthy nodes: %d available, %d required",
   len(available), required)
 }

 // Sort by latency for optimization
 sort.Slice(available, func(i, j int) bool {
  return available[i].GetLatency() < available[j].GetLatency()
 })

 // Select best nodes
 quorum := make([]*Node, required)
 copy(quorum, available[:required])

 return quorum, nil
}

// ValidateQuorum checks if a set of nodes forms a valid quorum
func (qm *QuorumManager) ValidateQuorum(nodes []*Node, qtype QuorumType) bool {
 qm.mu.RLock()
 required := qm.readQuorum
 if qtype == WriteQuorum {
  required = qm.writeQuorum
 }
 qm.mu.RUnlock()

 if len(nodes) < required {
  return false
 }

 // Check all nodes are healthy
 for _, node := range nodes {
  if !node.IsHealthy() {
   return false
  }
 }

 return true
}

// QuorumIntersection checks if two quorums intersect
func QuorumIntersection(q1, q2 []*Node) bool {
 set := make(map[string]bool)
 for _, n := range q1 {
  set[n.ID] = true
 }
 for _, n := range q2 {
  if set[n.ID] {
   return true
  }
 }
 return false
}

// ConsensusStore implements quorum-based consensus storage
type ConsensusStore struct {
 mu            sync.RWMutex
 qm            *QuorumManager
 data          map[string]*VersionedValue
 pendingWrites map[string]*PendingWrite
}

// VersionedValue stores a value with version metadata
type VersionedValue struct {
 Value     []byte
 Version   uint64
 Timestamp time.Time
 Writer    string
}

// PendingWrite tracks a write in progress
type PendingWrite struct {
 Key       string
 Value     []byte
 Version   uint64
 Quorum    []*Node
 Acks      int
 Timestamp time.Time
}

// NewConsensusStore creates a new consensus store
func NewConsensusStore(qm *QuorumManager) *ConsensusStore {
 return &ConsensusStore{
  qm:            qm,
  data:          make(map[string]*VersionedValue),
  pendingWrites: make(map[string]*PendingWrite),
 }
}

// Write performs a quorum write
func (cs *ConsensusStore) Write(ctx context.Context, key string, value []byte) error {
 // Select write quorum
 quorum, err := cs.qm.SelectQuorum(WriteQuorum, nil)
 if err != nil {
  return fmt.Errorf("failed to select quorum: %w", err)
 }

 // Generate new version
 cs.mu.Lock()
 existing := cs.data[key]
 newVersion := uint64(1)
 if existing != nil {
  newVersion = existing.Version + 1
 }

 pending := &PendingWrite{
  Key:       key,
  Value:     value,
  Version:   newVersion,
  Quorum:    quorum,
  Timestamp: time.Now(),
 }
 cs.pendingWrites[key] = pending
 cs.mu.Unlock()

 // Send writes to quorum
 acks := make(chan bool, len(quorum))
 for _, node := range quorum {
  go func(n *Node) {
   err := cs.sendWrite(ctx, n, key, value, newVersion)
   acks <- (err == nil)
  }(node)
 }

 // Collect acknowledgments
 ackCount := 0
 required := cs.qm.writeQuorum
 timeout := time.NewTimer(5 * time.Second)
 defer timeout.Stop()

 for ackCount < required {
  select {
  case ack := <-acks:
   if ack {
    ackCount++
   }
  case <-timeout.C:
   cs.mu.Lock()
   delete(cs.pendingWrites, key)
   cs.mu.Unlock()
   return fmt.Errorf("write timeout: %d/%d acks received", ackCount, required)
  case <-ctx.Done():
   cs.mu.Lock()
   delete(cs.pendingWrites, key)
   cs.mu.Unlock()
   return ctx.Err()
  }
 }

 // Write successful
 cs.mu.Lock()
 cs.data[key] = &VersionedValue{
  Value:     value,
  Version:   newVersion,
  Timestamp: time.Now(),
 }
 delete(cs.pendingWrites, key)
 cs.mu.Unlock()

 return nil
}

// Read performs a quorum read
func (cs *ConsensusStore) Read(ctx context.Context, key string) ([]byte, error) {
 // Select read quorum
 quorum, err := cs.qm.SelectQuorum(ReadQuorum, nil)
 if err != nil {
  return nil, fmt.Errorf("failed to select quorum: %w", err)
 }

 // Query quorum nodes
 responses := make(chan *VersionedValue, len(quorum))
 for _, node := range quorum {
  go func(n *Node) {
   val := cs.sendRead(ctx, n, key)
   responses <- val
  }(node)
 }

 // Collect responses
 var values []*VersionedValue
 timeout := time.NewTimer(5 * time.Second)
 defer timeout.Stop()

 required := cs.qm.readQuorum
 for len(values) < required {
  select {
  case val := <-responses:
   if val != nil {
    values = append(values, val)
   }
  case <-timeout.C:
   return nil, fmt.Errorf("read timeout: %d/%d responses received", len(values), required)
  case <-ctx.Done():
   return nil, ctx.Err()
  }
 }

 // Return value with highest version
 var latest *VersionedValue
 for _, val := range values {
  if latest == nil || val.Version > latest.Version {
   latest = val
  }
 }

 if latest == nil {
  return nil, fmt.Errorf("key not found")
 }

 return latest.Value, nil
}

// sendWrite sends a write to a node
func (cs *ConsensusStore) sendWrite(ctx context.Context, node *Node, key string, value []byte, version uint64) error {
 // In production: actual network call
 // Simulated with delay based on node latency
 select {
 case <-time.After(node.GetLatency()):
  return nil
 case <-ctx.Done():
  return ctx.Err()
 }
}

// sendRead sends a read to a node
func (cs *ConsensusStore) sendRead(ctx context.Context, node *Node, key string) *VersionedValue {
 // In production: actual network call
 select {
 case <-time.After(node.GetLatency()):
  cs.mu.RLock()
  defer cs.mu.RUnlock()
  return cs.data[node.ID+"-"+key] // Simulated per-node data
 case <-ctx.Done():
  return nil
 }
}

// HierarchicalQuorumManager implements hierarchical quorums
type HierarchicalQuorumManager struct {
 regions      map[string]*Region
 regionsPerQ  int
 racksPerRegion int
 nodesPerRack   int
}

// Region represents a geographic region
type Region struct {
 ID    string
 Racks map[string]*Rack
}

// Rack represents a server rack
type Rack struct {
 ID    string
 Nodes []*Node
}

// SelectHierarchicalQuorum selects a quorum across the hierarchy
func (hqm *HierarchicalQuorumManager) SelectHierarchicalQuorum() ([]*Node, error) {
 // Select regions
 var regionList []*Region
 for _, r := range hqm.regions {
  regionList = append(regionList, r)
 }

 if len(regionList) < hqm.regionsPerQ {
  return nil, fmt.Errorf("insufficient regions: %d available, %d required",
   len(regionList), hqm.regionsPerQ)
 }

 // Shuffle and select
 rand.Shuffle(len(regionList), func(i, j int) {
  regionList[i], regionList[j] = regionList[j], regionList[i]
 })

 selectedRegions := regionList[:hqm.regionsPerQ]

 // Collect nodes from selected regions
 var quorum []*Node
 for _, region := range selectedRegions {
  rackList := make([]*Rack, 0, len(region.Racks))
  for _, r := range region.Racks {
   rackList = append(rackList, r)
  }

  // Select racks within region
  rand.Shuffle(len(rackList), func(i, j int) {
   rackList[i], rackList[j] = rackList[j], rackList[i]
  })

  selectedRacks := rackList
  if len(rackList) > hqm.racksPerRegion {
   selectedRacks = rackList[:hqm.racksPerRegion]
  }

  // Collect nodes from selected racks
  for _, rack := range selectedRacks {
   selectedNodes := rack.Nodes
   if len(rack.Nodes) > hqm.nodesPerRack {
    selectedNodes = rack.Nodes[:hqm.nodesPerRack]
   }
   quorum = append(quorum, selectedNodes...)
  }
 }

 return quorum, nil
}

// VerifyIntersection verifies that two hierarchical quorums intersect
func (hqm *HierarchicalQuorumManager) VerifyIntersection(q1, q2 []*Node) bool {
 // Quorums intersect if they share at least one node
 set := make(map[string]bool)
 for _, n := range q1 {
  set[n.ID] = true
 }
 for _, n := range q2 {
  if set[n.ID] {
   return true
  }
 }
 return false
}
```

## Comparison with Related Protocols

| Aspect | Quorum Consensus | Two-Phase Commit | Raft | Paxos |
|--------|------------------|------------------|------|-------|
| **Coordination** | Partial (quorum only) | Full (all nodes) | Leader-based | Leader-based |
| **Fault Tolerance** | Configurable | Low (coordinator failure) | Majority | Majority |
| **Latency** | Single RTT (quorum) | 2+ RTTs | 1-2 RTTs | 2+ RTTs |
| **Scalability** | High | Low | Medium | Medium |
| **Use Case** | Key-value stores | Transactions | Log replication | Consensus |
| **Read Scaling** | Excellent | Poor | Good | Good |
| **Complexity** | Medium | Medium | Medium | High |

## Visual Representations

### Figure 1: Quorum System Taxonomy

```
┌─────────────────────────────────────────────────────────────────┐
│                  QUORUM SYSTEM TAXONOMY                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    QUORUM SYSTEMS                        │   │
│  └────────────────────┬────────────────────────────────────┘   │
│                       │                                         │
│         ┌─────────────┼─────────────┬─────────────┐            │
│         ▼             ▼             ▼             ▼            │
│    ┌─────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐       │
│    │Majority │  │Weighted  │  │Hierarch. │  │Grid/Tree │       │
│    │Quorums  │  │Voting    │  │Quorums   │  │Quorums   │       │
│    └────┬────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘       │
│         │            │             │             │              │
│    ┌────┴────┐  ┌────┴─────┐  ┌────┴─────┐  ┌────┴─────┐       │
│    │R = W =  │  │R + W > V │  │Multi-level│  │Row + Col │       │
│    │⌊N/2⌋+1  │  │W > V/2   │  │selection │  │or Tree   │       │
│    └─────────┘  └──────────┘  └──────────┘  └──────────┘       │
│                                                                 │
│  Examples:                                                      │
│    • Majority: Dynamo, Cassandra (default)                      │
│    • Weighted: Gifford's weighted voting                        │
│    • Hierarchical: WAN-distributed systems                      │
│    • Grid: High-availability clusters                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 2: Quorum Intersection Property

```
┌─────────────────────────────────────────────────────────────────┐
│                 QUORUM INTERSECTION PROPERTY                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Universe U = {A, B, C, D, E, F, G}                            │
│                                                                 │
│  Write Quorum Qw (size 4):                                      │
│  ┌─────────────────────────────────────┐                       │
│  │  ● A   ● B   ● C   ● D   ○ E   ...  │                       │
│  │   [==== Qw ====]                     │                       │
│  └─────────────────────────────────────┘                       │
│                                                                 │
│  Read Quorum Qr (size 3):                                       │
│  ┌─────────────────────────────────────┐                       │
│  │  ○ A   ● B   ● C   ○ D   ● E   ...  │                       │
│  │        [== Qr ==]                    │                       │
│  └─────────────────────────────────────┘                       │
│                                                                 │
│  Intersection: Qw ∩ Qr = {B, C} ≠ ∅                            │
│                                                                 │
│  Proof:                                                         │
│  |Qw| + |Qr| = 4 + 3 = 7 > |U| = 7? No, need R + W > N         │
│                                                                 │
│  With R + W > N:                                                │
│  Let Qw = {A, B, C, D}, Qr = {C, D, E}                          │
│  |Qw| + |Qr| = 4 + 3 = 7, N = 6                                 │
│  7 > 6, so intersection guaranteed                              │
│  Qw ∩ Qr = {C, D} ✓                                             │
│                                                                 │
│  ┌─────────────────────────────────────┐                       │
│  │  ● A   ● B   ● C   ● D   ● E   ○ F  │                       │
│  │   [==== Qw ====]                     │                       │
│  │            [== Qr ==]                │                       │
│  │            [intersect]               │                       │
│  └─────────────────────────────────────┘                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 3: Write Operation Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                  WRITE OPERATION FLOW                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Client                                                         │
│    │                                                            │
│    │ 1. Write(key, value)                                       │
│    ▼                                                            │
│  ┌─────────────┐                                                │
│  │   Client    │                                                │
│  │   Library   │                                                │
│  └──────┬──────┘                                                │
│         │ 2. Select Write Quorum (W = 3)                        │
│         ▼                                                       │
│  ┌─────────────────────────────────────────┐                    │
│  │         Write Quorum Nodes              │                    │
│  │                                         │                    │
│  │  ┌─────┐   ┌─────┐   ┌─────┐          │                    │
│  │  │ N1  │   │ N2  │   │ N3  │          │                    │
│  │  │     │   │     │   │     │          │                    │
│  │  └─────┘   └─────┘   └─────┘          │                    │
│  │     │         │         │              │                    │
│  └─────┼─────────┼─────────┼──────────────┘                    │
│         │         │         │                                   │
│         ▼         ▼         ▼                                   │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                           │
│  │3a.Write │ │3b.Write │ │3c.Write │                           │
│  │  key=v  │ │  key=v  │ │  key=v  │                           │
│  │ ver=n+1 │ │ ver=n+1 │ │ ver=n+1 │                           │
│  └────┬────┘ └────┬────┘ └────┬────┘                           │
│       │           │           │                                 │
│       ▼           ▼           ▼                                 │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                           │
│  │ 4a.ACK  │ │ 4b.ACK  │ │ 4c.TIME │                           │
│  │         │ │         │ │   OUT   │                           │
│  └────┬────┘ └────┬────┘ └─────────┘                           │
│       │           │                                             │
│       └───────────┘                                             │
│              │                                                  │
│              ▼                                                  │
│  ┌─────────────────┐                                            │
│  │ 5. W=3, ACK=2   │                                            │
│  │ 2 ≥ W/2+1? YES  │                                            │
│  │ Write Succeeds  │                                            │
│  └────────┬────────┘                                            │
│           │                                                     │
│           ▼                                                     │
│  ┌─────────────────┐                                            │
│  │ 6. Return OK    │                                            │
│  └─────────────────┘                                            │
│                                                                 │
│  Note: Write succeeds with W-1 acks if W > N/2                  │
│        (hinted handoff handles the slow/failed node)            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## References

1. Gifford, D. K. (1979). "Weighted voting for replicated data." ACM SOSP.
2. Thomas, R. H. (1979). "A majority consensus approach to concurrency control." ACM TODS.
3. Naor, M., & Wool, A. (1998). "The load, capacity, and availability of quorum systems." SIAM Journal on Computing.
4. Malkhi, D., & Reiter, M. (1998). "Byzantine quorum systems." Distributed Computing.
5. DeCandia, G., et al. (2007). "Dynamo: Amazon's highly available key-value store." ACM SOSP.
