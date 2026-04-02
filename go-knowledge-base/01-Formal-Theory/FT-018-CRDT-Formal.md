# FT-018: CRDT - Conflict-Free Replicated Data Types - Formal Specification

## Overview

Conflict-Free Replicated Data Types (CRDTs) are data structures that can be replicated across multiple nodes in a distributed system and updated concurrently without coordination. They guarantee strong eventual consistency: all replicas that have received the same set of updates will converge to the same state.

## Theoretical Foundations

### 1.1 CRDT Axioms

**Definition**: A CRDT is a data type that satisfies the following properties:

```
1. Associativity: (a ⊔ b) ⊔ c = a ⊔ (b ⊔ c)
2. Commutativity: a ⊔ b = b ⊔ a
3. Idempotence: a ⊔ a = a

Where ⊔ is the merge (join) operation
```

**Strong Eventual Consistency (SEC)**:

```
A system provides SEC if:
1. Eventual Delivery: An update delivered to one correct replica
   is eventually delivered to all correct replicas
2. Convergence: Correct replicas that have delivered the same
   updates eventually reach equivalent states
3. Termination: All method executions terminate
```

### 1.2 State-Based CRDTs (Convergent Replicated Data Types)

**State-Based CRDT Structure**:

```
A state-based CRDT is a tuple ⟨S, s⁰, q, u, m⟩ where:
- S: set of states
- s⁰ ∈ S: initial state
- q: query method (pure function S → result)
- u: update method (function update: S × args → S)
- m: merge method (function merge: S × S → S)
```

**State-Based CRDT Properties**:

```
1. Merge forms a join-semilattice:
   - Commutative: m(a, b) = m(b, a)
   - Associative: m(m(a, b), c) = m(a, m(b, c))
   - Idempotent: m(a, a) = a

2. Update monotonicity:
   ∀s ∈ S, args: s ≤ u(s, args)
   (where ≤ is the partial order defined by the semilattice)

3. Merge is LUB:
   m(a, b) = a ⊔ b (least upper bound)
```

### 1.3 Op-Based CRDTs (Commutative Replicated Data Types)

**Op-Based CRDT Structure**:

```
An op-based CRDT is a tuple ⟨S, s⁰, q, t, u, P⟩ where:
- S: set of states
- s⁰ ∈ S: initial state
- q: query method
- t: prepare method (t: S × args → operation)
- u: effect method (u: S × operation → S)
- P: delivery protocol ensuring causal delivery
```

**Op-Based CRDT Properties**:

```
1. Causal delivery: If operation o₁ causally precedes o₂,
   then o₁ is delivered before o₂ at all replicas

2. Commutativity of concurrent operations:
   For concurrent operations o₁ || o₂:
   u(u(s, o₁), o₂) = u(u(s, o₂), o₁)
```

### 1.4 Equivalence Proof

```
Theorem: State-based and op-based CRDTs are equivalent in expressive power.

Proof Sketch (State → Op):
Given a state-based CRDT, define:
- t(s, args) returns the delta δ = u(s, args) \\ s
- u(s, δ) = s ⊔ δ

Proof Sketch (Op → State):
Given an op-based CRDT, define:
- State S' = S × Log where Log is operation history
- merge((s₁, l₁), (s₂, l₂)) =
    apply all operations in l₁ ∪ l₂ to s⁰

Both constructions preserve SEC guarantees.
∎
```

## TLA+ Specification

```tla
----------------------------- MODULE CRDT -----------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Replicas,        \* Set of replica IDs
          Values,          \* Set of possible values
          MaxOperations    \* Bound for model checking

VARIABLES replicaState,    \* State of each replica
          operationLog,    \* Log of operations at each replica
          deliveryStatus   \* Which ops have been delivered where

\* CRDT State Type
CRDTState == [counter: Nat, set: SUBSET Values, map: [Values -> Nat]]

\* Operation Type
Operation == [type: {"increment", "add", "remove"},
              replica: Replicas,
              value: Values,
              timestamp: Nat]

\* Initial state
Init ==
  ∧ replicaState = [r ∈ Replicas ↦ [counter ↦ 0, set ↦ {}, map ↦ [v ∈ Values ↦ 0]]]
  ∧ operationLog = [r ∈ Replicas ↦ ⟨⟩]
  ∧ deliveryStatus = [op ∈ {} ↦ {}]

\* State-based merge operation (G-Counter example)
MergeGCounter(s1, s2) ==
  [counter ↦ Max({s1.counter, s2.counter}),
   set ↦ s1.set ∪ s2.set,
   map ↦ [v ∈ Values ↦ Max({s1.map[v], s2.map[v]})]]

\* Increment operation
Increment(replica) ==
  LET newState == [replicaState[replica] EXCEPT !.counter = @ + 1]
      newOp == [type ↦ "increment",
                replica ↦ replica,
                value ↦ "",
                timestamp ↦ Len(operationLog[replica]) + 1]
  IN ∧ replicaState' = [replicaState EXCEPT ![replica] = newState]
     ∧ operationLog' = [operationLog EXCEPT ![replica] = Append(@, newOp)]
     ∧ UNCHANGED deliveryStatus

\* Add to set operation
AddToSet(replica, value) ==
  LET newState == [replicaState[replica] EXCEPT !.set = @ ∪ {value}]
      newOp == [type ↦ "add",
                replica ↦ replica,
                value ↦ value,
                timestamp ↦ Len(operationLog[replica]) + 1]
  IN ∧ replicaState' = [replicaState EXCEPT ![replica] = newState]
     ∧ operationLog' = [operationLog EXCEPT ![replica] = Append(@, newOp)]
     ∧ UNCHANGED deliveryStatus

\* Anti-entropy: replica r1 pulls state from r2
AntiEntropy(r1, r2) ==
  ∧ r1 ≠ r2
  ∧ replicaState' = [replicaState EXCEPT ![r1] =
       MergeGCounter(@, replicaState[r2])]
  ∧ UNCHANGED ⟨operationLog, deliveryStatus⟩

\* Next state relation
Next ==
  ∨ ∃r ∈ Replicas: Increment(r)
  ∨ ∃r ∈ Replicas, v ∈ Values: AddToSet(r, v)
  ∨ ∃r1, r2 ∈ Replicas: AntiEntropy(r1, r2)

\* Invariants

\* Monotonicity: counters never decrease
CounterMonotonicity ==
  ∀r ∈ Replicas:
    [][replicaState[r].counter ≥ replicaState[r].counter']_replicaState

\* Set monotonicity: elements are never lost
SetMonotonicity ==
  ∀r ∈ Replicas:
    [][replicaState[r].set ⊆ replicaState'[r].set]_replicaState

\* Convergence: replicas that have exchanged states agree
Convergence ==
  ∀r1, r2 ∈ Replicas:
    replicaState[r1] = MergeGCounter(replicaState[r1], replicaState[r2])
    ⇒ replicaState[r1] = replicaState[r2]

\* SEC Property
StrongEventualConsistency ==
  ∀r1, r2 ∈ Replicas:
    LET ops1 == Range(operationLog[r1])
        ops2 == Range(operationLog[r2])
    IN ops1 = ops2 ⇒ replicaState[r1] = replicaState[r2]

=============================================================================
```

## Algorithm Pseudocode

### G-Counter (Grow-Only Counter)

```
Algorithm: G-Counter (State-Based)

State:
  P: number of replicas
  id: identifier of this replica (0 to P-1)
  n: array[0..P-1] of non-negative integers
      n[i] = increments from replica i

Initialize:
  n ← [0, 0, ..., 0]

Query():
  return Σᵢ n[i]

Increment():
  n[id] ← n[id] + 1

Merge(other):
  for i from 0 to P-1:
    n[i] ← max(n[i], other.n[i])

Correctness Proof:

Theorem: G-Counter satisfies SEC.

Proof:
1. Monotonicity: n[i] only increases via Increment or Merge
2. Associativity: max(max(a,b),c) = max(a,max(b,c))
3. Commutativity: max(a,b) = max(b,a)
4. Idempotence: max(a,a) = a

Merge forms a join-semilattice over vector clocks.
Therefore, concurrent increments converge correctly.
∎
```

### PN-Counter (Increment/Decrement Counter)

```
Algorithm: PN-Counter (State-Based)

State:
  P: number of replicas
  id: identifier of this replica
  p: array[0..P-1] of non-negative integers  (increments)
  n: array[0..P-1] of non-negative integers  (decrements)

Initialize:
  p ← [0, ..., 0]
  n ← [0, ..., 0]

Query():
  return Σᵢ p[i] - Σᵢ n[i]

Increment():
  p[id] ← p[id] + 1

Decrement():
  n[id] ← n[id] + 1

Merge(other):
  for i from 0 to P-1:
    p[i] ← max(p[i], other.p[i])
    n[i] ← max(n[i], other.n[i])

Limitations:
- Cannot go below 0 after merge (semantic anomaly)
- Solution: Use bounded counters or accept temporary negativity
```

### OR-Set (Observed-Removed Set)

```
Algorithm: OR-Set (Add-Wins Set)

State:
  A: set of (element, unique-tag) pairs  (added)
  R: set of (element, unique-tag) pairs  (removed)

Invariant: R ⊆ A (can only remove what was added)

Initialize:
  A ← ∅
  R ← ∅

Query():
  return {e : ∃t. (e,t) ∈ A ∧ (e,t) ∉ R}

Add(e):
  tag ← generateUniqueTag()
  A ← A ∪ {(e, tag)}

Remove(e):
  R ← R ∪ {(e, t) : (e, t) ∈ A}
  \* Remove all observed instances of e

Merge(other):
  A ← A ∪ other.A
  R ← R ∪ other.R

Properties:
- Add-wins: If concurrent add and remove, add wins
- Tags ensure each add is unique
- Remove only affects observed elements

Correctness:
Theorem: OR-Set converges to the expected state.

Proof:
Consider concurrent operations at replicas r1 and r2:

Case 1: r1 adds e, r2 adds e
  Both adds create different tags
  Merge keeps both: A = {(e,t1), (e,t2)}
  Query returns {e} ✓

Case 2: r1 adds e, r2 removes e (before seeing add)
  r2 removes nothing (didn't observe e)
  After merge: A = {(e,t1)}, R = ∅
  Query returns {e} ✓

Case 3: r1 adds e, r2 removes e (after seeing add)
  r2 removes observed tag: R = {(e,t1)}
  After merge: A = {(e,t1)}, R = {(e,t1)}
  Query returns ∅ ✓

Case 4: Concurrent add and remove at same replica
  Sequential execution, no conflict
∎
```

### LWW-Register (Last-Writer-Wins Register)

```
Algorithm: LWW-Register

State:
  value: stored value
  timestamp: logical timestamp
  replica: replica ID

Initialize:
  value ← ⊥
  timestamp ← 0
  replica ← myId

Query():
  return value

Update(v):
  value ← v
  timestamp ← timestamp + 1

Merge(other):
  if other.timestamp > timestamp ∨
     (other.timestamp = timestamp ∧ other.replica > replica):
    value ← other.value
    timestamp ← other.timestamp
    replica ← other.replica

Conflict Resolution:
- Higher timestamp wins
- If timestamps equal, higher replica ID wins (deterministic tie-break)

Note: Last-writer-wins semantics may lose updates
Alternative: MV-Register keeps multiple versions
```

### RGA (Replicated Growable Array)

```
Algorithm: RGA (Sequence CRDT)

State:
  nodes: map<id, Node> where Node = {id, value, parent, isDeleted}
  idCounter: monotonic counter for generating unique IDs

Operations:

InsertAfter(parentId, value):
  newId ← generateId()
  nodes[newId] ← {
    id: newId,
    value: value,
    parent: parentId,
    isDeleted: false
  }

Delete(id):
  if nodes[id] exists:
    nodes[id].isDeleted ← true

Query():
  return traverse(nodes, root)  \* depth-first traversal

Merge(other):
  for each (id, node) in other.nodes:
    if id ∉ nodes:
      nodes[id] ← node
    else:
      nodes[id].isDeleted ← nodes[id].isDeleted ∨ node.isDeleted

Traverse(nodes, root):
  result ← []
  for each node in children(root) ordered by id:
    if not node.isDeleted:
      result.append(node.value)
    result.appendAll(Traverse(nodes, node.id))
  return result

Properties:
- Concurrent inserts: both appear (order by ID)
- Insert concurrent with delete: insert wins (tombstone)
- Concurrent deletes: both delete (idempotent)
```

## Go Implementation

```go
// Package crdt provides CRDT implementations
package crdt

import (
 "crypto/rand"
 "encoding/binary"
 "fmt"
 "math"
 "math/big"
 "sort"
 "sync"
 "time"
)

// ReplicaID uniquely identifies a replica
type ReplicaID uint32

// Tag is a unique identifier for operations
type Tag struct {
 Replica ReplicaID
 Counter uint64
}

func (t Tag) String() string {
 return fmt.Sprintf("%d:%d", t.Replica, t.Counter)
}

// Less defines total order for tags
func (t Tag) Less(other Tag) bool {
 if t.Replica != other.Replica {
  return t.Replica < other.Replica
 }
 return t.Counter < other.Counter
}

// GCounter is a grow-only counter CRDT
type GCounter struct {
 mu       sync.RWMutex
 replica  ReplicaID
 numReplicas int
 counts   []uint64
}

// NewGCounter creates a new G-Counter
func NewGCounter(replica ReplicaID, numReplicas int) *GCounter {
 return &GCounter{
  replica:     replica,
  numReplicas: numReplicas,
  counts:      make([]uint64, numReplicas),
 }
}

// Increment increases the counter for this replica
func (c *GCounter) Increment() {
 c.mu.Lock()
 defer c.mu.Unlock()
 c.counts[c.replica]++
}

// Value returns the total count
func (c *GCounter) Value() uint64 {
 c.mu.RLock()
 defer c.mu.RUnlock()

 var sum uint64
 for _, v := range c.counts {
  sum += v
 }
 return sum
}

// Merge combines another G-Counter into this one
func (c *GCounter) Merge(other *GCounter) {
 c.mu.Lock()
 defer c.mu.Unlock()

 other.mu.RLock()
 defer other.mu.RUnlock()

 for i := 0; i < min(len(c.counts), len(other.counts)); i++ {
  if other.counts[i] > c.counts[i] {
   c.counts[i] = other.counts[i]
  }
 }
}

// Clone creates a copy of the counter
func (c *GCounter) Clone() *GCounter {
 c.mu.RLock()
 defer c.mu.RUnlock()

 clone := NewGCounter(c.replica, c.numReplicas)
 copy(clone.counts, c.counts)
 return clone
}

// PNCounter is a positive-negative counter CRDT
type PNCounter struct {
 p *GCounter
 n *GCounter
}

// NewPNCounter creates a new PN-Counter
func NewPNCounter(replica ReplicaID, numReplicas int) *PNCounter {
 return &PNCounter{
  p: NewGCounter(replica, numReplicas),
  n: NewGCounter(replica, numReplicas),
 }
}

// Increment increases the counter
func (c *PNCounter) Increment() {
 c.p.Increment()
}

// Decrement decreases the counter
func (c *PNCounter) Decrement() {
 c.n.Increment()
}

// Value returns the net count
func (c *PNCounter) Value() int64 {
 return int64(c.p.Value()) - int64(c.n.Value())
}

// Merge combines another PN-Counter into this one
func (c *PNCounter) Merge(other *PNCounter) {
 c.p.Merge(other.p)
 c.n.Merge(other.n)
}

// ORSet is an observed-remove set CRDT
type ORSet struct {
 mu      sync.RWMutex
 replica ReplicaID
 counter uint64
 adds    map[string]map[Tag]struct{}
 removes map[string]map[Tag]struct{}
}

// NewORSet creates a new OR-Set
func NewORSet(replica ReplicaID) *ORSet {
 return &ORSet{
  replica: replica,
  adds:    make(map[string]map[Tag]struct{}),
  removes: make(map[string]map[Tag]struct{}),
 }
}

// Add adds an element to the set
func (s *ORSet) Add(element string) Tag {
 s.mu.Lock()
 defer s.mu.Unlock()

 s.counter++
 tag := Tag{Replica: s.replica, Counter: s.counter}

 if s.adds[element] == nil {
  s.adds[element] = make(map[Tag]struct{})
 }
 s.adds[element][tag] = struct{}{}

 return tag
}

// Remove removes an element from the set
func (s *ORSet) Remove(element string) {
 s.mu.Lock()
 defer s.mu.Unlock()

 // Remove all observed tags for this element
 if tags, ok := s.adds[element]; ok {
  if s.removes[element] == nil {
   s.removes[element] = make(map[Tag]struct{})
  }
  for tag := range tags {
   s.removes[element][tag] = struct{}{}
  }
 }
}

// Contains checks if an element is in the set
func (s *ORSet) Contains(element string) bool {
 s.mu.RLock()
 defer s.mu.RUnlock()

 addTags := s.adds[element]
 removeTags := s.removes[element]

 // Element exists if there are adds not covered by removes
 for tag := range addTags {
  if _, removed := removeTags[tag]; !removed {
   return true
  }
 }
 return false
}

// Elements returns all elements in the set
func (s *ORSet) Elements() []string {
 s.mu.RLock()
 defer s.mu.RUnlock()

 var result []string
 for elem := range s.adds {
  if s.containsUnlocked(elem) {
   result = append(result, elem)
  }
 }
 sort.Strings(result)
 return result
}

func (s *ORSet) containsUnlocked(element string) bool {
 addTags := s.adds[element]
 removeTags := s.removes[element]

 for tag := range addTags {
  if _, removed := removeTags[tag]; !removed {
   return true
  }
 }
 return false
}

// Merge combines another OR-Set into this one
func (s *ORSet) Merge(other *ORSet) {
 s.mu.Lock()
 defer s.mu.Unlock()

 other.mu.RLock()
 defer other.mu.RUnlock()

 // Merge adds
 for elem, tags := range other.adds {
  if s.adds[elem] == nil {
   s.adds[elem] = make(map[Tag]struct{})
  }
  for tag := range tags {
   s.adds[elem][tag] = struct{}{}
  }
 }

 // Merge removes
 for elem, tags := range other.removes {
  if s.removes[elem] == nil {
   s.removes[elem] = make(map[Tag]struct{})
  }
  for tag := range tags {
   s.removes[elem][tag] = struct{}{}
  }
 }
}

// LWWRegister is a last-writer-wins register CRDT
type LWWRegister struct {
 mu        sync.RWMutex
 replica   ReplicaID
 value     interface{}
 timestamp uint64
 replicaTS ReplicaID
}

// NewLWWRegister creates a new LWW register
func NewLWWRegister(replica ReplicaID) *LWWRegister {
 return &LWWRegister{
  replica: replica,
 }
}

// Set updates the register value
func (r *LWWRegister) Set(value interface{}) {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.value = value
 r.timestamp++
 r.replicaTS = r.replica
}

// Get returns the current value
func (r *LWWRegister) Get() interface{} {
 r.mu.RLock()
 defer r.mu.RUnlock()
 return r.value
}

// Merge combines another register into this one
func (r *LWWRegister) Merge(other *LWWRegister) {
 r.mu.Lock()
 defer r.mu.Unlock()

 other.mu.RLock()
 defer other.mu.RUnlock()

 // LWW semantics: higher timestamp wins, tie-break by replica ID
 if other.timestamp > r.timestamp ||
  (other.timestamp == r.timestamp && other.replicaTS > r.replicaTS) {
  r.value = other.value
  r.timestamp = other.timestamp
  r.replicaTS = other.replicaTS
 }
}

// RGANode represents a node in the Replicated Growable Array
type RGANode struct {
 ID        Tag
 Value     interface{}
 Parent    *Tag
 IsDeleted bool
 Timestamp uint64
}

// RGA implements the Replicated Growable Array CRDT
type RGA struct {
 mu      sync.RWMutex
 replica ReplicaID
 counter uint64
 nodes   map[Tag]*RGANode
 root    Tag
}

// NewRGA creates a new RGA
func NewRGA(replica ReplicaID) *RGA {
 r := &RGA{
  replica: replica,
  nodes:   make(map[Tag]*RGANode),
  root:    Tag{Replica: 0, Counter: 0},
 }
 // Initialize root
 r.nodes[r.root] = &RGANode{
  ID:     r.root,
  Parent: nil,
 }
 return r
}

// InsertAfter inserts a value after the specified parent node
func (r *RGA) InsertAfter(parentID Tag, value interface{}) Tag {
 r.mu.Lock()
 defer r.mu.Unlock()

 r.counter++
 newID := Tag{Replica: r.replica, Counter: r.counter}

 r.nodes[newID] = &RGANode{
  ID:        newID,
  Value:     value,
  Parent:    &parentID,
  IsDeleted: false,
  Timestamp: uint64(time.Now().UnixNano()),
 }

 return newID
}

// Delete marks a node as deleted
func (r *RGA) Delete(id Tag) {
 r.mu.Lock()
 defer r.mu.Unlock()

 if node, ok := r.nodes[id]; ok {
  node.IsDeleted = true
 }
}

// ToSlice returns the array as a slice
func (r *RGA) ToSlice() []interface{} {
 r.mu.RLock()
 defer r.mu.RUnlock()

 return r.traverse(&r.root)
}

func (r *RGA) traverse(parent *Tag) []interface{} {
 var result []interface{}

 // Find all children of parent
 var children []*RGANode
 for _, node := range r.nodes {
  if node.Parent != nil && *node.Parent == *parent {
   children = append(children, node)
  }
 }

 // Sort children by ID for deterministic order
 sort.Slice(children, func(i, j int) bool {
  return children[i].ID.Less(children[j].ID)
 })

 // Traverse
 for _, child := range children {
  if !child.IsDeleted {
   result = append(result, child.Value)
  }
  result = append(result, r.traverse(&child.ID)...)
 }

 return result
}

// Merge combines another RGA into this one
func (r *RGA) Merge(other *RGA) {
 r.mu.Lock()
 defer r.mu.Unlock()

 other.mu.RLock()
 defer other.mu.RUnlock()

 for id, node := range other.nodes {
  if existing, ok := r.nodes[id]; ok {
   // Tombstone wins
   existing.IsDeleted = existing.IsDeleted || node.IsDeleted
  } else {
   // Copy node
   newNode := *node
   r.nodes[id] = &newNode
  }
 }
}

// CRDTStore manages multiple CRDTs
type CRDTStore struct {
 mu     sync.RWMutex
 gcounters map[string]*GCounter
 pncounters map[string]*PNCounter
 orsets   map[string]*ORSet
 lwwregs  map[string]*LWWRegister
 rgas     map[string]*RGA
}

// NewCRDTStore creates a new CRDT store
func NewCRDTStore() *CRDTStore {
 return &CRDTStore{
  gcounters:  make(map[string]*GCounter),
  pncounters: make(map[string]*PNCounter),
  orsets:     make(map[string]*ORSet),
  lwwregs:    make(map[string]*LWWRegister),
  rgas:       make(map[string]*RGA),
 }
}

// Helper function
func min(a, b int) int {
 if a < b {
  return a
 }
 return b
}

// generateUniqueTag creates a cryptographically secure unique tag
func generateUniqueTag() string {
 b := make([]byte, 16)
 rand.Read(b)
 return fmt.Sprintf("%x", b)
}

// GenerateRandomTag creates a random tag for testing
func GenerateRandomTag() (Tag, error) {
 replica, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
 if err != nil {
  return Tag{}, err
 }

 counter, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
 if err != nil {
  return Tag{}, err
 }

 return Tag{
  Replica: ReplicaID(replica.Uint64()),
  Counter: counter.Uint64(),
 }, nil
}
```

## Comparison with Related Approaches

| Property | CRDTs | Operational Transform | Consensus | MVCC |
|----------|-------|----------------------|-----------|------|
| **Coordination** | None | Central server | Required | Required |
| **Availability** | Always available | Server-dependent | Limited | Limited |
| **Convergence** | Guaranteed | Guaranteed | N/A | N/A |
| **Semantics** | Application-specific | Positional | Linearizable | Snapshot |
| **Use Cases** | Collaborative editing, counters | Text editing | Transactions | Databases |
| **Complexity** | Medium | High | Medium | Medium |
| **Message Size** | State can be large | Small (operations) | Small | Medium |

## Visual Representations

### Figure 1: CRDT Taxonomy

```
┌─────────────────────────────────────────────────────────────────┐
│                    CRDT TAXONOMY                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   CRDT Types                            │   │
│  └────────────────────┬────────────────────────────────────┘   │
│                       │                                         │
│         ┌─────────────┼─────────────┐                          │
│         ▼             ▼             ▼                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐              │
│  │  Counters   │ │    Sets     │ │  Registers  │              │
│  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘              │
│         │               │               │                       │
│    ┌────┴────┐     ┌────┴────┐     ┌────┴────┐                │
│    ▼         ▼     ▼         ▼     ▼         ▼                │
│  ┌────┐   ┌────┐ ┌────┐   ┌────┐ ┌────┐   ┌────┐             │
│  │G-Counter│PN-Counter│OR-Set│LWW-Set│LWW-Reg│MV-Reg│        │
│  └────┘   └────┘ └────┘   └────┘ └────┘   └────┘             │
│                                                                 │
│  Other Types:                                                   │
│  • Maps: AW-Map, OR-Map, Delta-Map                              │
│  • Sequences: RGA, WOOT, Logoot, LSEQ                           │
│  • Graphs: Add-Remove Partial Order, Collaborative Graphs       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 2: State-Based CRDT Convergence

```
┌─────────────────────────────────────────────────────────────────┐
│              STATE-BASED CRDT CONVERGENCE                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Time ──────►                                                   │
│                                                                 │
│  Replica A:    [5,0,0] ──────► [5,3,0] ──────► [5,3,2]        │
│                     │               │               │           │
│                     │   Anti-Entropy │               │           │
│                     │               │               │           │
│                     ▼               ▼               ▼           │
│  Replica B:    [0,3,0] ──────► [5,3,0] ──────► [5,3,2]        │
│                                     │               │           │
│                         Anti-Entropy│               │           │
│                                     │               │           │
│                                     ▼               ▼           │
│  Replica C:    [0,0,2] ──────► [5,0,2] ──────► [5,3,2]        │
│                                                                 │
│  Legend: [A-count, B-count, C-count]                           │
│                                                                 │
│  Convergence Proof:                                            │
│  - Merge operation: max([a₁,b₁,c₁], [a₂,b₂,c₂])                │
│  - Merge is associative, commutative, idempotent               │
│  - All replicas converge to [5,3,2] = max of all states        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 3: OR-Set Conflict Resolution

```
┌─────────────────────────────────────────────────────────────────┐
│              OR-SET CONFLICT RESOLUTION                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Initial State: All replicas have {}                           │
│                                                                 │
│  Concurrent Operations:                                         │
│  ┌─────────────────┐        ┌─────────────────┐                │
│  │   Replica A     │        │   Replica B     │                │
│  │                 │        │                 │                │
│  │  Add(x) ────────┼────────┼───────────────  │                │
│  │  creates (x,t1) │        │                 │                │
│  │                 │        │  Add(x) ────────┼── creates (x,t2)│
│  │  Remove(x) ─────┼────────┼───────────────  │                │
│  │  removes (x,t1) │        │                 │                │
│  └─────────────────┘        └─────────────────┘                │
│                                                                 │
│  State after operations:                                        │
│  Replica A: A={(x,t1)}, R={(x,t1)}                             │
│  Replica B: A={(x,t2)}, R={}                                   │
│                                                                 │
│  After Merge:                                                   │
│  A = {(x,t1), (x,t2)}                                           │
│  R = {(x,t1)}                                                   │
│                                                                 │
│  Query Result: {x} because (x,t2) ∈ A and (x,t2) ∉ R            │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Outcome: Add-wins! Concurrent add beats remove.        │   │
│  │  Reasoning:                                             │   │
│  │  • Remove only affects observed tags                    │   │
│  │  • Replica A didn't observe (x,t2)                      │   │
│  │  • Therefore (x,t2) survives merge                      │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## References

1. Shapiro, M., Preguiça, N., Baquero, C., & Zawirski, M. (2011). "A comprehensive study of convergent and commutative replicated data types." INRIA Technical Report.
2. Shapiro, M., Preguiça, N., Baquero, C., & Zawirski, M. (2011). "Conflict-free replicated data types." ACM SSS.
3. Roh, H. G., Jeon, M., Kim, J. S., & Lee, J. (2011). "Replicated abstract data types: Building blocks for collaborative applications." Journal of Parallel and Distributed Computing.
4. Baquero, C., Almeida, P. S., & Shoker, A. (2014). "Making operation-based CRDTs operation-based." ACM PaPoC.
5. Kleppmann, M., & Beresford, A. R. (2017). "A conflict-free replicated JSON datatype." IEEE TPDS.
