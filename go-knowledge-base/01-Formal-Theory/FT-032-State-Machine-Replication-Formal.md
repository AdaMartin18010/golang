# FT-032: State Machine Replication - Formal Theory

> **Dimension**: Formal Theory
> **Level**: S (>15KB)
> **Tags**: #smr #state-machine-replication #paxos #raft #determinism
> **Authoritative Sources**:
>
> - Schneider, F. B. (1990). "Implementing Fault-Tolerant Services Using the State Machine Approach". ACM CSUR
> - Lamport, L. (1978). "Time, Clocks, and the Ordering of Events". CACM
> - Ongaro, D., & Ousterhout, J. (2014). "In Search of an Understandable Consensus Algorithm". USENIX ATC

---

## 1. Theoretical Foundations

### 1.1 State Machine Model

**Definition 1.1 (State Machine)**: A state machine $\mathcal{M}$ is a tuple $(S, S_0, \Sigma, \delta, \Lambda, \lambda)$ where:

- $S$: Set of states
- $S_0 \in S$: Initial state
- $\Sigma$: Input alphabet (commands)
- $\delta: S \times \Sigma \rightarrow S$: Transition function
- $\Lambda$: Output alphabet
- $\lambda: S \times \Sigma \rightarrow \Lambda$: Output function

**Definition 1.2 (Deterministic State Machine)**: A state machine is deterministic if:

$$
\forall s \in S, \sigma \in \Sigma: |\delta(s, \sigma)| = 1
$$

### 1.2 State Machine Replication

**Definition 1.3 (SMR)**: A protocol that maintains multiple copies of a deterministic state machine, ensuring all copies start in the same state and process the same commands in the same order.

**Requirements**:

1. **Agreement**: All replicas apply commands in the same order
2. **Validity**: Applied commands are valid client requests

**Theorem 1.1 (SMR Safety)**: If all replicas start in the same state, and apply the same sequence of deterministic commands, they remain in identical states.

*Proof by induction*:

- Base: All start in $S_0$ ✓
- Inductive step: If all in state $s_i$ before command $c$, all transition to $\delta(s_i, c)$ ✓
- Therefore, all states identical after any command sequence ∎

---

## 2. Log-Based SMR

### 2.1 Log Structure

```
Log Entry Format:
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│  Index          │  Term           │  Command        │  Timestamp      │
│  (uint64)       │  (uint64)       │  (bytes)        │  (int64)        │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘

Log Consistency Properties:
1. If entry i has term t, then entries 0..i-1 also exist
2. If two logs have entry i with same term, then entries 0..i match
3. Committed entries are durable
```

### 2.2 Command Execution

```
Algorithm ExecuteLog:

State:
  log: array of Entry
  state: ApplicationState
  lastApplied: int ← 0

On Commit(entry):
  log[entry.index] = entry

  // Apply committed entries in order
  while lastApplied < commitIndex:
    lastApplied++
    entry = log[lastApplied]
    state = Apply(state, entry.command)
    ReplyToClient(entry.client, result)
```

---

## 3. Replicated State Machine Correctness

### 3.1 Safety Properties

**Property 1 (State Machine Safety)**: If any replica applies command $c$ at position $i$, no replica applies a different command $c' \neq c$ at position $i$.

**Property 2 (Leader Completeness)**: If a command is committed in term $t$, it will be present in all leaders' logs for terms $> t$.

**Property 3 (State Machine Progress)**: All replicas eventually apply all committed commands.

### 3.2 Liveness Properties

**Property 4 (Availability)**: If a correct client submits a command, and a quorum of replicas are correct, the command is eventually committed.

---

## 4. Optimizations

### 4.1 Log Compaction

```
Algorithm Snapshotting:

Periodically (when log exceeds threshold):
  1. Pause command processing
  2. Serialize current state
  3. Write snapshot to disk
  4. Truncate log before snapshot
  5. Resume processing

State Transfer (catching up slow replica):
  Leader sends:
    - Latest snapshot
    - Log entries after snapshot
```

### 4.2 Read Optimization

```
Read Strategies:

1. Leader Read (strong consistency):
   - Contact leader
   - Leader commits "read command"
   - Returns result
   Latency: 1 RTT

2. Lease Read (optimized):
   - Leader grants time lease
   - Reads served locally during lease
   - Risk: stale reads if leader fails
   Latency: local

3. Quorum Read (follower reads):
   - Contact f+1 replicas
   - Return latest version
   - No log entry needed
   Latency: 1 RTT (parallel)
```

---

## 5. TLA+ Specification

```tla
----------------------------- MODULE StateMachineReplication -----------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Replicas,
          Commands,
          MaxLogSize

VARIABLES log,           \* Log at each replica
          state,         \* Application state
          commitIndex,   \* Index of highest committed entry
          appliedIndex   \* Index of highest applied entry

Init ==
  /\ log = [r \in Replicas |-> <<>>]
  /\ state = [r \in Replicas |-> "initial"]
  /\ commitIndex = [r \in Replicas |-> 0]
  /\ appliedIndex = [r \in Replicas |-> 0]

AppendEntry(leader, cmd) ==
  /\ Len(log[leader]) < MaxLogSize
  /\ log' = [log EXCEPT ![leader] = Append(@, cmd)]
  /\ UNCHANGED <<state, commitIndex, appliedIndex>>

ReplicateEntry(follower, leader, idx) ==
  /\ idx <= Len(log[leader])
  /\ idx > Len(log[follower])
  /\ log' = [log EXCEPT ![follower] = Append(@, log[leader][idx])]
  /\ UNCHANGED <<state, commitIndex, appliedIndex>>

CommitEntry(leader, idx) ==
  /\ idx <= Len(log[leader])
  /\ commitIndex' = [commitIndex EXCEPT ![leader] =
                         IF idx > @ THEN idx ELSE @]
  /\ UNCHANGED <<log, state, appliedIndex>>

ApplyEntries(r) ==
  /\ appliedIndex[r] < commitIndex[r]
  /\ LET nextIdx == appliedIndex[r] + 1
         cmd == log[r][nextIdx]
     IN /\ state' = [state EXCEPT ![r] = ApplyCommand(@, cmd)]
        /\ appliedIndex' = [appliedIndex EXCEPT ![r] = nextIdx]
  /\ UNCHANGED <<log, commitIndex>>

Safety ==
  \* All replicas apply same commands in order
  \A r1, r2 \in Replicas:
    \A i \in 1..appliedIndex[r1]:
      i <= appliedIndex[r2] => log[r1][i] = log[r2][i]

Liveness ==
  \* Committed entries eventually applied
  \A r \in Replicas:
    commitIndex[r] > appliedIndex[r] ~> appliedIndex[r] = commitIndex[r]

=============================================================================
```

---

## 6. Go Implementation

```go
// Package smr provides state machine replication
package smr

import (
 "context"
 "fmt"
 "sync"
 "sync/atomic"
)

// StateMachine defines the application interface
type StateMachine interface {
 Apply(cmd Command) (Result, error)
 Snapshot() ([]byte, error)
 Restore(snapshot []byte) error
}

// Command represents a state machine command
type Command struct {
 ID      uint64
 Op      string
 Args    []byte
 ClientID string
}

// Result represents command result
type Result struct {
 Value []byte
 Error error
}

// LogEntry represents a log entry
type LogEntry struct {
 Index     uint64
 Term      uint64
 Command   Command
 Committed bool
}

// SMR implements state machine replication
type SMR struct {
 mu           sync.RWMutex
 id           string
 stateMachine StateMachine

 log          []*LogEntry
 commitIndex  uint64
 lastApplied  uint64
 currentTerm  uint64

 applyCh      chan *LogEntry
 resultCh     map[uint64]chan Result

 stopCh       chan struct{}
 wg           sync.WaitGroup
}

// NewSMR creates a new SMR instance
func NewSMR(id string, sm StateMachine) *SMR {
 s := &SMR{
  id:           id,
  stateMachine: sm,
  log:          make([]*LogEntry, 0),
  applyCh:      make(chan *LogEntry, 100),
  resultCh:     make(map[uint64]chan Result),
  stopCh:       make(chan struct{}),
 }

 s.wg.Add(1)
 go s.applier()

 return s
}

// Submit submits a command for replication
func (s *SMR) Submit(cmd Command) (Result, error) {
 // Append to log
 s.mu.Lock()
 idx := uint64(len(s.log)) + 1
 entry := &LogEntry{
  Index:   idx,
  Term:    s.currentTerm,
  Command: cmd,
 }
 s.log = append(s.log, entry)

 // Create result channel
 ch := make(chan Result, 1)
 s.resultCh[idx] = ch
 s.mu.Unlock()

 // Wait for result
 select {
 case result := <-ch:
  return result, nil
 case <-s.stopCh:
  return Result{}, fmt.Errorf("SMR stopped")
 }
}

// Commit marks entries as committed up to index
func (s *SMR) Commit(index uint64) {
 s.mu.Lock()
 defer s.mu.Unlock()

 if index <= s.commitIndex {
  return
 }

 for i := s.commitIndex; i < index && i < uint64(len(s.log)); i++ {
  s.log[i].Committed = true
  select {
  case s.applyCh <- s.log[i]:
  default:
  }
 }

 s.commitIndex = index
}

func (s *SMR) applier() {
 defer s.wg.Done()

 for {
  select {
  case <-s.stopCh:
   return
  case entry := <-s.applyCh:
   s.apply(entry)
  }
 }
}

func (s *SMR) apply(entry *LogEntry) {
 result, err := s.stateMachine.Apply(entry.Command)

 s.mu.Lock()
 defer s.mu.Unlock()

 s.lastApplied = entry.Index

 if ch, ok := s.resultCh[entry.Index]; ok {
  ch <- Result{Value: result.Value, Error: err}
  delete(s.resultCh, entry.Index)
 }
}

// GetLogEntry returns a log entry at index
func (s *SMR) GetLogEntry(index uint64) (*LogEntry, bool) {
 s.mu.RLock()
 defer s.mu.RUnlock()

 if index == 0 || index > uint64(len(s.log)) {
  return nil, false
 }

 return s.log[index-1], true
}

// LogSize returns the current log size
func (s *SMR) LogSize() int {
 s.mu.RLock()
 defer s.mu.RUnlock()
 return len(s.log)
}

// Stop stops the SMR
func (s *SMR) Stop() {
 close(s.stopCh)
 s.wg.Wait()
}

// Snapshot creates a snapshot
func (s *SMR) Snapshot() ([]byte, error) {
 s.mu.Lock()
 defer s.mu.Unlock()

 return s.stateMachine.Snapshot()
}

// Truncate removes entries before index
func (s *SMR) Truncate(before uint64) {
 s.mu.Lock()
 defer s.mu.Unlock()

 if before > uint64(len(s.log)) {
  return
 }

 s.log = s.log[before-1:]
}
```

---

## 7. Visual Representations

### 7.1 SMR Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    STATE MACHINE REPLICATION ARCHITECTURE                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Client                                                                      │
│     │                                                                         │
│     │ Submit(Command)                                                         │
│     ▼                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      Consensus Module                               │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │    │
│  │  │ Replica │  │ Replica │  │ Replica │  │ Replica │  │ Replica │   │    │
│  │  │   R1    │  │   R2    │  │   R3    │  │   R4    │  │   R5    │   │    │
│  │  │(Leader) │  │         │  │         │  │         │  │         │   │    │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘   │    │
│  │       │            │            │            │            │        │    │
│  │       └────────────┴────────────┴────────────┴────────────┘        │    │
│  │                          Consensus Log                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│           │                                                                 │
│           │ Apply (in order)                                                │
│           ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Application State Machine                        │    │
│  │                                                                     │    │
│  │  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐               │    │
│  │  │   State     │──>│  Command 1  │──>│  Command 2  │──> ...        │    │
│  │  │    S0       │   │   (Log[1])  │   │   (Log[2])  │               │    │
│  │  └─────────────┘   └─────────────┘   └─────────────┘               │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│           │                                                                 │
│           │ Response                                                        │
│           ▼                                                                 │
│        Result                                                               │
│                                                                              │
│  Key Properties:                                                             │
│  • All replicas start in same state S0                                      │
│  • All replicas apply same commands in same order                           │
│  • Therefore all replicas reach same state                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Log Replication Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LOG REPLICATION FLOW                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Leader (R1)        Follower (R2)        Follower (R3)                      │
│     │                   │                   │                               │
│     │  AppendEntries    │                   │                               │
│     │  Index: 5         │                   │                               │
│     │  Term: 3          │                   │                               │
│     │  Entries: [Cmd]   │                   │                               │
│     │──────────────────>│                   │                               │
│     │                   │                   │                               │
│     │  AppendEntries    │                   │                               │
│     │──────────────────────────────────────>│                               │
│     │                   │                   │                               │
│     │                   │  Success          │                               │
│     │<──────────────────│  MatchIndex: 5    │                               │
│     │                   │                   │                               │
│     │                   │  Success          │                               │
│     │<──────────────────────────────────────│                               │
│     │                   │                   │                               │
│     │ [Quorum achieved] │                   │                               │
│     │                   │                   │                               │
│     │  CommitIndex = 5  │                   │                               │
│     │  (Apply to SM)    │                   │                               │
│     │                   │                   │                               │
│     │  Heartbeat        │                   │                               │
│     │  CommitIndex: 5   │                   │                               │
│     │──────────────────>│                   │                               │
│     │                   │  CommitIndex = 5  │                               │
│     │                   │  (Apply to SM)    │                               │
│     │                   │                   │                               │
│     │                   │                   │  CommitIndex = 5              │
│     │                   │                   │  (Apply to SM)                │
│     │                   │                   │                               │
│                                                                              │
│  Consistency: All apply command 5 only after committed by leader           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Snapshot and Recovery

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SNAPSHOT AND RECOVERY                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Log with Snapshot:                                                          │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐    │
│  │  Snapshotted      |        Log Entries                             │    │
│  │  (Compacted)      |                                                │    │
│  │                   │   1001   1002   1003   1004   1005             │    │
│  │  ┌──────────┐    │  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐           │    │
│  │  │ State at │    │  │Cmd1│ │Cmd2│ │Cmd3│ │Cmd4│ │Cmd5│           │    │
│  │  │ Index 1000     │  └────┘ └────┘ └────┘ └────┘ └────┘           │    │
│  │  │ (Binary)│    │                                                │    │
│  │  └──────────┘    │  CommitIndex = 1005                            │    │
│  │                   │  LastApplied = 1005                            │    │
│  └────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  Recovery Scenario:                                                          │
│                                                                              │
│  Replica R2 crashes, restarts:                                               │
│                                                                              │
│  1. Load snapshot (index 1000)                                               │
│     State = Snapshot.State                                                   │
│                                                                              │
│  2. Request log from leader                                                  │
│     Send: "Need entries from 1001"                                          │
│                                                                              │
│  3. Leader sends missing entries                                             │
│     Entries: [1001, 1002, 1003, 1004, 1005]                                 │
│                                                                              │
│  4. Apply entries                                                            │
│     For each entry in [1001..1005]:                                         │
│       State = Apply(State, entry.Command)                                    │
│                                                                              │
│  5. Caught up!                                                               │
│     State matches other replicas                                             │
│                                                                              │
│  Benefits of Snapshotting:                                                   │
│  • Bounded log size                                                          │
│  • Fast recovery                                                             │
│  • Reduced memory usage                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    STATE MACHINE REPLICATION SUMMARY                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Core Principle:                                                             │
│  Deterministic(State, Commands) → Same Result                                │
│                                                                              │
│  Key Requirements:                                                           │
│  1. Deterministic state machine                                             │
│  2. All replicas start in same state                                        │
│  3. Same commands in same order                                             │
│                                                                              │
│  Optimizations:                                                              │
│  • Batching: Amortize consensus overhead                                    │
│  • Pipelining: Overlap consensus instances                                  │
│  • Snapshots: Bound log size                                                │
│  • Read leases: Scale reads                                                 │
│                                                                              │
│  Use Cases:                                                                  │
│  • Key-value stores (etcd, ZooKeeper)                                       │
│  • Coordination services                                                    │
│  • Configuration management                                                 │
│  • Distributed locking                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
