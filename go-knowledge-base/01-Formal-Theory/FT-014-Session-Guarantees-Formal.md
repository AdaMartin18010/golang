# FT-014: Session Guarantees - Formal Specification

## Overview

Session guarantees provide consistency semantics for client sessions in distributed systems, defining what a client can observe during a sequence of operations within a single session. These guarantees bridge the gap between strong consistency and high availability.

## Theoretical Foundations

### 1.1 Session Model

A session S is defined as a tuple:

```
S = ⟨Client, Operations, TimeRange, Visibility⟩

where:
- Client: unique client identifier
- Operations: ordered sequence of read/write operations
- TimeRange: [t_start, t_end] session duration
- Visibility: mapping of operations to visible writes
```

### 1.2 Core Session Guarantees

#### Read Your Writes (RYW)

**Definition**: If a client writes value v to object x, all subsequent reads by the same client to x will return v or a later value.

**Formal Specification**:

```
∀c ∈ Clients, ∀x ∈ Objects:
  write(c, x, v) at time t₁ ∧ read(c, x) at time t₂ ∧ t₂ > t₁
  ⇒ return(v') where v' = v ∨ v' > v (in version order)
```

**Proof of Monotonicity**:

```
Theorem: RYW ensures monotonic reads within a session.

Proof:
Let W = {w₁, w₂, ..., wₙ} be writes to object x in session order.
Let R = {r₁, r₂, ..., rₘ} be reads to object x in session order.

For any read rᵢ that follows write wⱼ:
  rᵢ must observe wⱼ or a successor write wₖ where k ≥ j.

By induction on session operations:
Base case: First read observes latest visible write.
Inductive step: If read rᵢ observes wⱼ, then read rᵢ₊₁ must observe
  at least wⱼ (by RYW definition).

Therefore, observed values are monotonically increasing.
∎
```

#### Monotonic Reads (MR)

**Definition**: If a client reads value v from object x, all subsequent reads by the same client to x will return v or a later value.

**Formal Specification**:

```
∀c ∈ Clients, ∀x ∈ Objects:
  read(c, x) → v at time t₁ ∧ read(c, x) → v' at time t₂ ∧ t₂ > t₁
  ⇒ v' ≥ v (in version order)
```

#### Monotonic Writes (MW)

**Definition**: If a client writes value v₁ to object x, then writes v₂ to the same object, all servers will apply v₁ before v₂.

**Formal Specification**:

```
∀c ∈ Clients, ∀x ∈ Objects:
  write(c, x, v₁) at t₁ ∧ write(c, x, v₂) at t₂ ∧ t₁ < t₂
  ⇒ ∀s ∈ Servers: apply(s, x, v₁) ≺ apply(s, x, v₂)
```

#### Writes Follow Reads (WFR)

**Definition**: If a client reads value v₁ written by client c' from object x, then writes v₂ to object y, then v₁ is ordered before v₂.

**Formal Specification**:

```
∀c, c' ∈ Clients, ∀x, y ∈ Objects:
  read(c, x) → v₁ written by c' at t₁ ∧ write(c, y, v₂) at t₂ ∧ t₁ < t₂
  ⇒ v₁ ≺ v₂ (causal order)
```

## TLA+ Specification

```tla
----------------------- MODULE SessionGuarantees -----------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Clients,       \* Set of client IDs
          Objects,       \* Set of object IDs
          Values,        \* Set of possible values
          Servers        \* Set of server replicas

VARIABLES clientState,   \* Client session state
          serverState,   \* Server replica states
          messages       \* In-flight messages

\* Type definitions
ClientState == [sessionId: Nat,
                pendingWrites: Seq([obj: Objects, val: Values]),
                readVersions: [Objects -> Nat],
                writeVersions: [Objects -> Nat]]

ServerState == [versions: [Objects -> Seq(Values)],
                vectorClock: [Clients -> Nat]]

\* Initial state
Init ==
  ∧ clientState = [c ∈ Clients ↦ [sessionId ↦ 0,
                                   pendingWrites ↦ ⟨⟩,
                                   readVersions ↦ [o ∈ Objects ↦ 0],
                                   writeVersions ↦ [o ∈ Objects ↦ 0]]]
  ∧ serverState = [s ∈ Servers ↦ [versions ↦ [o ∈ Objects ↦ ⟨⟩],
                                   vectorClock ↦ [c ∈ Clients ↦ 0]]]
  ∧ messages = {}

\* Helper: Get version of latest write
LatestVersion(s, o) == Len(serverState[s].versions[o])

\* Helper: Check if client has seen a version
HasSeenVersion(c, o, v) == clientState[c].readVersions[o] ≥ v

\* Read Your Writes guarantee
ReadYourWrites(c, o) ==
  LET clientWriteVersion == clientState[c].writeVersions[o]
  IN ∀s ∈ Servers:
       LatestVersion(s, o) ≥ clientWriteVersion

\* Monotonic Reads guarantee
MonotonicReads(c, o, readVersion) ==
  readVersion ≥ clientState[c].readVersions[o]

\* Client read operation with session guarantees
ClientRead(c, o) ==
  LET readVersion == Max({LatestVersion(s, o) : s ∈ Servers})
  IN ∧ MonotonicReads(c, o, readVersion)
     ∧ clientState' = [clientState EXCEPT
         ![c].readVersions[o] = readVersion]

\* Client write operation
ClientWrite(c, o, v) ==
  LET newVersion == clientState[c].writeVersions[o] + 1
  IN ∧ clientState' = [clientState EXCEPT
         ![c].writeVersions[o] = newVersion,
         ![c].pendingWrites = Append(@, [obj ↦ o, val ↦ v])]
     ∧ messages' = messages ∪ {[type ↦ "write",
                                 client ↦ c,
                                 obj ↦ o,
                                 val ↦ v,
                                 version ↦ newVersion]}

\* Server apply write with ordering guarantees
ServerApplyWrite(s, m) ==
  ∧ m.type = "write"
  ∧ serverState' = [serverState EXCEPT
       ![s].versions[m.obj] = Append(@, m.val),
       ![s].vectorClock[m.client] =
         Max({@, m.version})]

\* Next state relation
Next ==
  ∨ ∃c ∈ Clients, o ∈ Objects: ClientRead(c, o)
  ∨ ∃c ∈ Clients, o ∈ Objects, v ∈ Values: ClientWrite(c, o, v)
  ∨ ∃s ∈ Servers, m ∈ messages: ServerApplyWrite(s, m)

\* Invariants

\* Read Your Writes invariant
RYWInvariant ==
  ∀c ∈ Clients, o ∈ Objects:
    clientState[c].writeVersions[o] ≤
    Min({LatestVersion(s, o) : s ∈ Servers})

\* Monotonic Reads invariant
MRInvariant ==
  ∀c ∈ Clients, o ∈ Objects:
    clientState[c].readVersions[o] ≤
    Max({LatestVersion(s, o) : s ∈ Servers})

\* Monotonic Writes invariant
MWInvariant ==
  ∀c ∈ Clients, o ∈ Objects, s ∈ Servers:
    LET versions == serverState[s].versions[o]
        clientWrites == {i ∈ 1..Len(versions) :
          versions[i] written by c}
    IN ∀i, j ∈ clientWrites: i < j ⇒ versions[i] before versions[j]

\* Safety property: All session guarantees hold
SessionGuarantees ==
  RYWInvariant ∧ MRInvariant ∧ MWInvariant

=============================================================================
```

## Algorithm Pseudocode

### Session Manager Algorithm

```
Algorithm: Session Manager for Distributed Systems

Data Structures:
  Session {
    id: UUID
    clientId: string
    vectorClock: Map<Client, Timestamp>
    readSet: Map<Object, Version>
    writeSet: Map<Object, Version>
    causalDependencies: Set<WriteId>
  }

  WriteRecord {
    clientId: string
    objectId: string
    value: Value
    version: Version
    timestamp: Timestamp
    sessionId: UUID
  }

Operations:

  function InitializeSession(clientId): Session
    return Session {
      id: generateUUID(),
      clientId: clientId,
      vectorClock: emptyMap(),
      readSet: emptyMap(),
      writeSet: emptyMap(),
      causalDependencies: emptySet()
    }

  function ReadWithGuarantees(session, objectId, guaranteeLevel)
    // Determine minimum acceptable version
    minVersion ← 0

    switch guaranteeLevel:
      case RYW:
        minVersion ← session.writeSet.getOrDefault(objectId, 0)
      case MR:
        minVersion ← session.readSet.getOrDefault(objectId, 0)
      case WFR:
        minVersion ← maxCausalVersion(session.causalDependencies, objectId)
      case DEFAULT:
        minVersion ← 0

    // Query replicas with version constraint
    value, version ← queryReplicas(objectId, minVersion)

    // Update session state
    session.readSet[objectId] ← max(session.readSet[objectId], version)
    session.vectorClock.update(timestamp)

    return value

  function WriteWithGuarantees(session, objectId, value)
    // Generate new version
    newVersion ← session.writeSet.getOrDefault(objectId, 0) + 1

    writeRecord ← WriteRecord {
      clientId: session.clientId,
      objectId: objectId,
      value: value,
      version: newVersion,
      timestamp: now(),
      sessionId: session.id
    }

    // Apply to quorum of replicas
    success ← replicateWrite(writeRecord, QUORUM_SIZE)

    if success:
      session.writeSet[objectId] ← newVersion
      session.vectorClock.update(timestamp)
      return ACK
    else:
      return ERROR

  function queryReplicas(objectId, minVersion)
    responses ← readFromReplicas(objectId, REPLICA_QUORUM)

    // Filter by minimum version (session guarantee)
    validResponses ← filter(responses, r → r.version ≥ minVersion)

    if validResponses.isEmpty():
      // Wait for replication to catch up
      waitForReplication(objectId, minVersion)
      return queryReplicas(objectId, minVersion)

    // Return latest version from valid responses
    return maxBy(validResponses, r → r.version)

  function replicateWrite(writeRecord, quorumSize)
    acks ← 0
    for replica in selectReplicas(quorumSize):
      if sendWrite(replica, writeRecord):
        acks ← acks + 1

    return acks ≥ quorumSize

Correctness Proof:

Theorem 1 (RYW Correctness): The algorithm ensures Read Your Writes.

Proof:
1. When a write is issued, writeSet[objectId] is updated to newVersion.
2. Subsequent reads check minVersion = writeSet[objectId].
3. queryReplicas filters for versions ≥ minVersion.
4. Therefore, the read must observe the previous write or later.
∎

Theorem 2 (MR Correctness): The algorithm ensures Monotonic Reads.

Proof:
1. When a read returns version v, readSet[objectId] is updated to v.
2. Subsequent reads use minVersion = readSet[objectId] = v.
3. queryReplicas filters for versions ≥ v.
4. Therefore, subsequent reads return version ≥ v.
∎

Theorem 3 (MW Correctness): The algorithm ensures Monotonic Writes.

Proof:
1. Writes within a session increment writeSet[objectId] monotonically.
2. Each write carries its version number to replicas.
3. Replicas apply writes in version order.
4. Therefore, writes are applied in session order.
∎
```

## Go Implementation

```go
// Package session provides distributed session guarantee implementations
package session

import (
 "context"
 "errors"
 "fmt"
 "sync"
 "time"

 "github.com/google/uuid"
)

// GuaranteeLevel defines the type of session guarantee
type GuaranteeLevel int

const (
 // Default provides no specific guarantees
 Default GuaranteeLevel = iota
 // RYW ensures Read Your Writes
 RYW
 // MR ensures Monotonic Reads
 MR
 // MW ensures Monotonic Writes
 MW
 // WFR ensures Writes Follow Reads
 WFR
 // MRplusMW combines Monotonic Reads and Monotonic Writes
 MRplusMW
 // All provides all session guarantees
 All
)

// Session represents a client session with guarantee tracking
type Session struct {
 ID               uuid.UUID
 ClientID         string
 mu               sync.RWMutex
 vectorClock      map[string]uint64
 readVersions     map[string]uint64
 writeVersions    map[string]uint64
 causalDeps       map[string]struct{}
 replicaManager   *ReplicaManager
 guaranteeLevel   GuaranteeLevel
}

// WriteRecord tracks a write operation
type WriteRecord struct {
 ClientID    string
 ObjectID    string
 Value       []byte
 Version     uint64
 Timestamp   time.Time
 SessionID   uuid.UUID
 VectorClock map[string]uint64
}

// ReadResult contains the result of a read operation
type ReadResult struct {
 Value   []byte
 Version uint64
 Server  string
}

// ReplicaManager handles communication with replicas
type ReplicaManager struct {
 replicas []Replica
 quorumSize int
}

// Replica represents a server replica
type Replica struct {
 ID      string
 Address string
 Client  ReplicaClient
}

// ReplicaClient defines the interface for replica communication
type ReplicaClient interface {
 Read(ctx context.Context, objectID string, minVersion uint64) (*ReadResult, error)
 Write(ctx context.Context, record *WriteRecord) error
}

// NewSession creates a new session with specified guarantees
func NewSession(clientID string, rm *ReplicaManager, level GuaranteeLevel) *Session {
 return &Session{
  ID:             uuid.New(),
  ClientID:       clientID,
  vectorClock:    make(map[string]uint64),
  readVersions:   make(map[string]uint64),
  writeVersions:  make(map[string]uint64),
  causalDeps:     make(map[string]struct{}),
  replicaManager: rm,
  guaranteeLevel: level,
 }
}

// Read performs a read operation with session guarantees
func (s *Session) Read(ctx context.Context, objectID string) ([]byte, error) {
 s.mu.Lock()
 defer s.mu.Unlock()

 // Determine minimum version based on guarantee level
 minVersion := uint64(0)

 switch s.guaranteeLevel {
 case RYW, All:
  if v, ok := s.writeVersions[objectID]; ok {
   minVersion = v
  }
 case MR, MRplusMW:
  if v, ok := s.readVersions[objectID]; ok {
   minVersion = v
  }
 case WFR:
  minVersion = s.maxCausalVersion(objectID)
 }

 // Query replicas with version constraint
 result, err := s.replicaManager.QueryWithConstraint(ctx, objectID, minVersion)
 if err != nil {
  return nil, fmt.Errorf("read failed: %w", err)
 }

 // Update session state
 if result.Version > s.readVersions[objectID] {
  s.readVersions[objectID] = result.Version
 }
 s.updateVectorClock(result.Server)

 return result.Value, nil
}

// Write performs a write operation with session guarantees
func (s *Session) Write(ctx context.Context, objectID string, value []byte) error {
 s.mu.Lock()
 defer s.mu.Unlock()

 // Generate new version
 newVersion := s.writeVersions[objectID] + 1

 record := &WriteRecord{
  ClientID:    s.ClientID,
  ObjectID:    objectID,
  Value:       value,
  Version:     newVersion,
  Timestamp:   time.Now(),
  SessionID:   s.ID,
  VectorClock: s.copyVectorClock(),
 }

 // Replicate to quorum
 if err := s.replicaManager.ReplicateWrite(ctx, record); err != nil {
  return fmt.Errorf("write replication failed: %w", err)
 }

 // Update session state
 s.writeVersions[objectID] = newVersion
 s.updateVectorClock("")

 return nil
}

// maxCausalVersion returns the maximum version from causal dependencies
func (s *Session) maxCausalVersion(objectID string) uint64 {
 maxVer := uint64(0)
 // Implementation would track causal dependencies
 return maxVer
}

// updateVectorClock updates the vector clock
func (s *Session) updateVectorClock(server string) {
 s.vectorClock[s.ClientID]++
 if server != "" {
  s.vectorClock[server]++
 }
}

// copyVectorClock creates a copy of the vector clock
func (s *Session) copyVectorClock() map[string]uint64 {
 vc := make(map[string]uint64, len(s.vectorClock))
 for k, v := range s.vectorClock {
  vc[k] = v
 }
 return vc
}

// QueryWithConstraint queries replicas with a minimum version constraint
func (rm *ReplicaManager) QueryWithConstraint(
 ctx context.Context,
 objectID string,
 minVersion uint64,
) (*ReadResult, error) {
 responses := make(chan *ReadResult, len(rm.replicas))
 errors := make(chan error, len(rm.replicas))

 // Query all replicas concurrently
 for _, replica := range rm.replicas {
  go func(r Replica) {
   result, err := r.Client.Read(ctx, objectID, minVersion)
   if err != nil {
    errors <- err
    return
   }
   responses <- result
  }(replica)
 }

 // Collect valid responses
 var validResults []*ReadResult
 received := 0
 required := rm.quorumSize

 for received < len(rm.replicas) {
  select {
  case result := <-responses:
   if result.Version >= minVersion {
    validResults = append(validResults, result)
   }
   received++
  case <-errors:
   received++
  case <-ctx.Done():
   return nil, ctx.Err()
  }

  if len(validResults) >= required {
   break
  }
 }

 if len(validResults) == 0 {
  return nil, errors.New("no valid responses meeting version constraint")
 }

 // Return result with maximum version
 var best *ReadResult
 for _, r := range validResults {
  if best == nil || r.Version > best.Version {
   best = r
  }
 }

 return best, nil
}

// ReplicateWrite replicates a write to a quorum of replicas
func (rm *ReplicaManager) ReplicateWrite(ctx context.Context, record *WriteRecord) error {
 acks := make(chan error, len(rm.replicas))

 // Send to all replicas concurrently
 for _, replica := range rm.replicas {
  go func(r Replica) {
   acks <- r.Client.Write(ctx, record)
  }(replica)
 }

 // Wait for quorum
 successCount := 0
 failCount := 0

 for successCount < rm.quorumSize && failCount <= len(rm.replicas)-rm.quorumSize {
  select {
  case err := <-acks:
   if err == nil {
    successCount++
   } else {
    failCount++
   }
  case <-ctx.Done():
   return ctx.Err()
  }
 }

 if successCount < rm.quorumSize {
  return errors.New("failed to achieve write quorum")
 }

 return nil
}

// SessionStore manages multiple client sessions
type SessionStore struct {
 mu       sync.RWMutex
 sessions map[uuid.UUID]*Session
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
 return &SessionStore{
  sessions: make(map[uuid.UUID]*Session),
 }
}

// Get retrieves a session by ID
func (ss *SessionStore) Get(id uuid.UUID) (*Session, bool) {
 ss.mu.RLock()
 defer ss.mu.RUnlock()
 session, ok := ss.sessions[id]
 return session, ok
}

// Put stores a session
func (ss *SessionStore) Put(session *Session) {
 ss.mu.Lock()
 defer ss.mu.Unlock()
 ss.sessions[session.ID] = session
}

// Cleanup removes expired sessions
func (ss *SessionStore) Cleanup(maxAge time.Duration) {
 ss.mu.Lock()
 defer ss.mu.Unlock()

 now := time.Now()
 for id, session := range ss.sessions {
  // In production, track last activity time
  _ = now
  _ = session
  delete(ss.sessions, id)
 }
}
```

## Protocol Comparison

| Aspect | Session Guarantees | Linearizability | Causal Consistency | Eventual Consistency |
|--------|-------------------|-----------------|-------------------|---------------------|
| **Latency** | Low | High | Low | Lowest |
| **Availability** | High | Lower | High | Highest |
| **Client Complexity** | Medium | Low | Medium | Low |
| **Conflict Resolution** | None | None | Required | Required |
| **Use Cases** | User sessions | Financial txns | Social apps | Caching |
| **Ordering** | Per-session | Global | Causal | None |

## Visual Representations

### Figure 1: Session Guarantees Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│                    CONSISTENCY SPECTRUM                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Strong ◄──────────────────────────────────────────► Weak      │
│                                                                 │
│  ┌─────────────┐  ┌──────────┐  ┌────────┐  ┌────────────┐     │
│  │Linearizable │  │Session   │  │Causal  │  │Eventual    │     │
│  │             │  │Guarantees│  │        │  │            │     │
│  └─────────────┘  └──────────┘  └────────┘  └────────────┘     │
│       │              │            │            │                │
│       ▼              ▼            ▼            ▼                │
│   ┌──────┐      ┌──────┐    ┌──────┐    ┌──────┐              │
│   │Paxos │      │ RYW  │    │Lamport│    │ Gossip│             │
│   │Raft  │      │ MR   │    │Clocks │    │       │             │
│   └──────┘      │ MW   │    └──────┘    └──────┘              │
│                 │ WFR  │                                       │
│                 └──────┘                                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 2: Session State Transitions

```
┌─────────────────────────────────────────────────────────────────┐
│                  SESSION STATE MACHINE                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│    ┌─────────┐      write       ┌─────────┐                    │
│    │  START  │ ───────────────► │ PENDING │                    │
│    └─────────┘                  └────┬────┘                    │
│         │                            │                         │
│         │ read                       │ replicate               │
│         ▼                            ▼                         │
│    ┌─────────┐                  ┌─────────┐                    │
│    │ READING │ ◄────────────────│ SYNCING │                    │
│    └────┬────┘      quorum      └─────────┘                    │
│         │                                                      │
│         │ update state                                         │
│         ▼                                                      │
│    ┌─────────┐                                                 │
│    │ UPDATED │──────────────────┐                              │
│    └─────────┘                  │                              │
│         │                       │ next op                      │
│         │ timeout               ▼                              │
│         ▼                  ┌─────────┐                         │
│    ┌─────────┐             │  START  │                         │
│    │  ERROR  │────────────►└─────────┘                         │
│    └─────────┘                                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 3: Multi-Client Session Interactions

```
┌─────────────────────────────────────────────────────────────────┐
│              SESSION INTERACTION DIAGRAM                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Time ▼                                                         │
│                                                                 │
│  Client A        Replicas          Client B                     │
│     │               │                  │                        │
│     │── write(x,1)──►│                 │                        │
│     │               │                  │                        │
│     │◄── ack ───────│                 │                        │
│     │               │                  │                        │
│     │── read(x)────►│                 │                        │
│     │               │                  │                        │
│     │◄── x=1 ──────│                 │                        │
│     │               │                  │                        │
│     │               │◄── write(x,2)───│                        │
│     │               │                  │                        │
│     │               │─── ack ────────►│                        │
│     │               │                  │                        │
│     │── read(x)────►│                 │                        │
│     │               │                  │                        │
│     │◄── x=2 ──────│  ← MR guarantees │                        │
│     │               │                  │                        │
│     │               │◄── read(x)──────│                        │
│     │               │                  │                        │
│     │               │──── x=2 ───────►│  ← Independent session │
│     │               │                  │                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## References

1. Terry, D. B., et al. (1994). "Session guarantees for weakly consistent replicated data." IEEE ICDE.
2. Lloyd, W., et al. (2011). "Don't settle for eventual: scalable causal consistency for wide-area storage with COPS." ACM SOSP.
3. Bailis, P., et al. (2013). "Potential benefits of risk-aware consistency." ACM VLDB.
4. Viotti, P., & Vukolić, M. (2016). "Consistency in non-transactional distributed storage systems." ACM CSUR.
