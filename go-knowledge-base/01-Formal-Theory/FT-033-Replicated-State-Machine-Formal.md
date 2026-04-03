# FT-033: Replicated State Machine - Formal Specification

> **Dimension**: Formal Theory
> **Level**: S (>15KB)
> **Tags**: #rsm #replication #consensus #formal-methods #safety
> **Authoritative Sources**:
>
> - Schneider, F. B. (1990). "Implementing Fault-Tolerant Services Using the State Machine Approach". ACM CSUR
> - Lamport, L. (1984). "Using Time Instead of Timeout for Fault-Tolerant Distributed Systems". ACM TOPLAS
> - Van Renesse, R., & Altinbuken, D. (2015). "Paxos Made Moderately Complex". ACM CSUR

---

## 1. Formal Model

### 1.1 State Machine Definition

A replicated state machine is formally defined as:

$$
\mathcal{R} = (\Pi, \Sigma, S, \delta, S_0, \mathcal{C})
$$

where:

- $\Pi = \{p_1, p_2, ..., p_n\}$: Set of replicas
- $\Sigma$: Input alphabet (commands)
- $S$: Set of states
- $\delta: S \times \Sigma \rightarrow S$: Deterministic transition function
- $S_0 \in S$: Initial state
- $\mathcal{C}$: Consensus protocol

### 1.2 Execution Semantics

**Definition 1.1 (Execution)**: An execution $E$ is a sequence of state transitions:

$$
E = s_0 \xrightarrow{c_1} s_1 \xrightarrow{c_2} s_2 \xrightarrow{c_3} ... \xrightarrow{c_k} s_k
$$

where $s_0 = S_0$ and $\forall i: s_i = \delta(s_{i-1}, c_i)$

**Definition 1.2 (Replica State)**: The state of replica $p$ at time $t$:

$$
\text{State}(p, t) = (\text{log}_p(t), \text{app}_p(t), \text{commit}_p(t))
$$

where:

- $\text{log}_p(t)$: Command log
- $\text{app}_p(t)$: Application state
- $\text{commit}_p(t)$: Commit index

---

## 2. Correctness Properties

### 2.1 Safety Properties

**Property 1 (Agreement)**: All correct replicas agree on the same sequence of committed commands.

$$
\forall p, q \in \text{Correct}, \forall i \leq \min(\text{commit}_p, \text{commit}_q): \text{log}_p[i] = \text{log}_q[i]
$$

**Property 2 (Validity)**: Committed commands are valid client requests.

$$
\forall p \in \Pi, \forall i \leq \text{commit}_p: \text{log}_p[i] \in \text{ValidCommands}
$$

**Property 3 (State Consistency)**: All correct replicas that have applied the same number of commands are in the same state.

$$
\forall p, q \in \text{Correct}: \text{applied}_p = \text{applied}_q \Rightarrow \text{app}_p = \text{app}_q
$$

### 2.2 Liveness Properties

**Property 4 (Progress)**: If a correct client submits a command, eventually some correct replica applies it.

$$
\text{Submit}(c) \Rightarrow \Diamond(\exists p \in \text{Correct}: c \in \text{log}_p)
$$

**Property 5 (Termination)**: All correct replicas eventually apply all committed commands.

$$
\forall p \in \text{Correct}: \Diamond(\text{applied}_p = \text{commit}_p)
$$

---

## 3. Protocol Specification

### 3.1 Client Protocol

```
Algorithm ClientProtocol:

Submit(cmd):
  // Phase 1: Submit to any replica
  replica ← SelectRandom(CorrectReplicas)
  Send(SUBMIT, cmd) to replica

  // Phase 2: Wait for response
  StartTimer(TIMEOUT)

  loop:
    On Receive(RESULT, result):
      return result

    On TimerExpired:
      replica ← SelectRandom(CorrectReplicas)
      Send(SUBMIT, cmd) to replica
      RestartTimer(TIMEOUT)
```

### 3.2 Replica Protocol

```
Algorithm ReplicaProtocol:

State:
  id: ReplicaID
  log: List<Command>
  state: ApplicationState
  commitIndex: integer ← 0
  lastApplied: integer ← 0
  currentView: integer ← 0

On ReceiveSubmit(cmd):
  if IsLeader():
    index ← log.Length() + 1
    log[index] ← cmd
    Broadcast(Propose, currentView, index, cmd)
  else:
    Forward(cmd, Leader())

On ReceivePropose(view, index, cmd):
  if view ≥ currentView and Valid(cmd):
    log[index] ← cmd
    Send(Accept, view, index, Hash(cmd)) to Leader

On ReceiveCommit(view, index):
  if view ≥ currentView:
    commitIndex ← max(commitIndex, index)
    ApplyCommitted()

Procedure ApplyCommitted():
  while lastApplied < commitIndex:
    lastApplied ← lastApplied + 1
    cmd ← log[lastApplied]
    result ← Execute(cmd)
    if cmd.client = self:
      ReplyToClient(result)
```

---

## 4. TLA+ Specification

```tla
----------------------------- MODULE ReplicatedStateMachine -----------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Replicas,
          Clients,
          Commands,
          MaxLogSize

VARIABLES replicaLog,      \* Log at each replica
          replicaState,    \* Application state
          commitIndex,     \* Commit index
          appliedIndex,    \* Applied index
          clientCommands,  \* Commands submitted by clients
          clientResults    \* Results received by clients

Replica == Replicas
Client == Clients
Command == Commands

Init ==
  /\ replicaLog = [r \in Replica |-> <<>>]
  /\ replicaState = [r \in Replica |-> [state |-> "initial", value |-> 0]]
  /\ commitIndex = [r \in Replica |-> 0]
  /\ appliedIndex = [r \in Replica |-> 0]
  /\ clientCommands = [c \in Client |-> {}]
  /\ clientResults = [c \in Client |-> {}]

\* Client submits command
SubmitCommand(c, cmd) ==
  /\ clientCommands' = [clientCommands EXCEPT ![c] = @ \union {cmd}]
  /\ UNCHANGED <<replicaLog, replicaState, commitIndex, appliedIndex, clientResults>>

\* Leader proposes command
ProposeCommand(leader, cmd, index) ==
  /\ Len(replicaLog[leader]) < MaxLogSize
  /\ index = Len(replicaLog[leader]) + 1
  /\ replicaLog' = [replicaLog EXCEPT ![leader] = Append(@, cmd)]
  /\ UNCHANGED <<replicaState, commitIndex, appliedIndex, clientCommands, clientResults>>

\* Replica accepts proposal
AcceptCommand(replica, leader, index) ==
  /\ index <= Len(replicaLog[leader])
  /\ replicaLog' = [replicaLog EXCEPT ![replica] =
                         IF Len(@) < index THEN Append(@, replicaLog[leader][index])
                         ELSE @]
  /\ UNCHANGED <<replicaState, commitIndex, appliedIndex, clientCommands, clientResults>>

\* Commit commands up to index
CommitCommands(leader, index) ==
  /\ index <= Len(replicaLog[leader])
  /\ commitIndex' = [commitIndex EXCEPT ![leader] =
                         IF index > @ THEN index ELSE @]
  /\ UNCHANGED <<replicaLog, replicaState, appliedIndex, clientCommands, clientResults>>

\* Apply committed commands
ApplyCommands(replica) ==
  /\ appliedIndex[replica] < commitIndex[replica]
  /\ LET nextIdx == appliedIndex[replica] + 1
         cmd == replicaLog[replica][nextIdx]
         currentState == replicaState[replica]
         newState == [value |-> currentState.value + cmd]
     IN /\ replicaState' = [replicaState EXCEPT ![replica] = newState]
        /\ appliedIndex' = [appliedIndex EXCEPT ![replica] = nextIdx]
  /\ UNCHANGED <<replicaLog, commitIndex, clientCommands, clientResults>>

Next ==
  /\ \E c \in Client, cmd \in Command: SubmitCommand(c, cmd)
  \/ \E leader \in Replica, cmd \in Command, index \in 1..MaxLogSize:
       ProposeCommand(leader, cmd, index)
  \/ \E replica, leader \in Replica, index \in 1..MaxLogSize:
       AcceptCommand(replica, leader, index)
  \/ \E leader \in Replica, index \in 1..MaxLogSize:
       CommitCommands(leader, index)
  \/ \E replica \in Replica: ApplyCommands(replica)

\* Safety: All replicas agree on committed prefix
Agreement ==
  \A r1, r2 \in Replica:
    LET minCommit == Min(commitIndex[r1], commitIndex[r2])
    IN \A i \in 1..minCommit: replicaLog[r1][i] = replicaLog[r2][i]

\* Safety: Committed commands are applied
CommitImpliesApplied ==
  \A r \in Replica: appliedIndex[r1] <= commitIndex[r1]

\* Safety: Applied state is deterministic
DeterministicState ==
  \A r1, r2 \in Replica:
    appliedIndex[r1] = appliedIndex[r2] => replicaState[r1] = replicaState[r2]

\* Liveness: Eventually committed
EventuallyCommitted ==
  \A c \in Client, cmd \in Command:
    cmd \in clientCommands[c] =>
      <>(\E r \in Replica: \E i \in 1..commitIndex[r]: replicaLog[r][i] = cmd)

=============================================================================
```

---

## 5. Go Implementation

```go
// Package rsm provides replicated state machine implementation
package rsm

import (
 "context"
 "fmt"
 "sync"
 "sync/atomic"
)

// Command represents a state machine command
type Command struct {
 ID       uint64
 ClientID string
 Op       Operation
 Args     []byte
}

// Operation is the operation type
type Operation int

const (
 OpRead Operation = iota
 OpWrite
 OpDelete
)

// State represents application state
type State interface {
 Apply(cmd Command) (Result, error)
 Clone() State
}

// Result is the command result
type Result struct {
 Value []byte
 Error error
}

// Entry is a log entry
type Entry struct {
 Index     uint64
 Term      uint64
 Command   Command
 Committed bool
}

// RSM implements replicated state machine
type RSM struct {
 mu          sync.RWMutex
 id          string
 replicas    []string

 log         []*Entry
 state       State
 commitIndex uint64
 lastApplied uint64
 currentTerm uint64

 // Channels
 proposeCh   chan Command
 commitCh    chan uint64
 applyCh     chan *Entry

 // Callbacks
 onCommit    func(uint64, Command)
 onApply     func(uint64, Command, Result)

 stopCh      chan struct{}
 wg          sync.WaitGroup
}

// NewRSM creates a new RSM
func NewRSM(id string, replicas []string, initialState State) *RSM {
 return &RSM{
  id:        id,
  replicas:  replicas,
  log:       make([]*Entry, 0),
  state:     initialState,
  proposeCh: make(chan Command, 100),
  commitCh:  make(chan uint64, 100),
  applyCh:   make(chan *Entry, 100),
  stopCh:    make(chan struct{}),
 }
}

// SetCallbacks sets event callbacks
func (r *RSM) SetCallbacks(onCommit func(uint64, Command), onApply func(uint64, Command, Result)) {
 r.onCommit = onCommit
 r.onApply = onApply
}

// Start starts the RSM
func (r *RSM) Start() error {
 r.wg.Add(2)
 go r.committer()
 go r.applier()
 return nil
}

// Stop stops the RSM
func (r *RSM) Stop() {
 close(r.stopCh)
 r.wg.Wait()
}

// Propose proposes a command
func (r *RSM) Propose(cmd Command) (uint64, error) {
 r.mu.Lock()
 defer r.mu.Unlock()

 index := uint64(len(r.log)) + 1
 entry := &Entry{
  Index:   index,
  Term:    r.currentTerm,
  Command: cmd,
 }
 r.log = append(r.log, entry)

 return index, nil
}

// Commit commits entries up to index
func (r *RSM) Commit(index uint64) {
 select {
 case r.commitCh <- index:
 case <-r.stopCh:
 }
}

func (r *RSM) committer() {
 defer r.wg.Done()

 for {
  select {
  case <-r.stopCh:
   return
  case index := <-r.commitCh:
   r.doCommit(index)
  }
 }
}

func (r *RSM) doCommit(index uint64) {
 r.mu.Lock()
 defer r.mu.Unlock()

 if index <= r.commitIndex || index > uint64(len(r.log)) {
  return
 }

 for i := r.commitIndex; i < index; i++ {
  r.log[i].Committed = true

  if r.onCommit != nil {
   go r.onCommit(i+1, r.log[i].Command)
  }

  select {
  case r.applyCh <- r.log[i]:
  default:
  }
 }

 r.commitIndex = index
}

func (r *RSM) applier() {
 defer r.wg.Done()

 for {
  select {
  case <-r.stopCh:
   return
  case entry := <-r.applyCh:
   r.doApply(entry)
  }
 }
}

func (r *RSM) doApply(entry *Entry) {
 r.mu.Lock()
 if entry.Index != r.lastApplied+1 {
  r.mu.Unlock()
  return
 }
 r.mu.Unlock()

 result, err := r.state.Apply(entry.Command)
 if err != nil {
  result.Error = err
 }

 r.mu.Lock()
 r.lastApplied = entry.Index
 r.mu.Unlock()

 if r.onApply != nil {
  go r.onApply(entry.Index, entry.Command, result)
 }
}

// GetLog returns log entry at index
func (r *RSM) GetLog(index uint64) (*Entry, bool) {
 r.mu.RLock()
 defer r.mu.RUnlock()

 if index == 0 || index > uint64(len(r.log)) {
  return nil, false
 }

 return r.log[index-1], true
}

// GetState returns current state
func (r *RSM) GetState() State {
 r.mu.RLock()
 defer r.mu.RUnlock()

 return r.state.Clone()
}

// GetCommitIndex returns commit index
func (r *RSM) GetCommitIndex() uint64 {
 return atomic.LoadUint64(&r.commitIndex)
}

// GetLastApplied returns last applied index
func (r *RSM) GetLastApplied() uint64 {
 return atomic.LoadUint64(&r.lastApplied)
}

// IsConsistent checks if state is consistent with another RSM
func (r *RSM) IsConsistent(other *RSM) bool {
 r.mu.RLock()
 other.mu.RLock()
 defer r.mu.RUnlock()
 defer other.mu.RUnlock()

 minCommit := r.commitIndex
 if other.commitIndex < minCommit {
  minCommit = other.commitIndex
 }

 for i := uint64(0); i < minCommit; i++ {
  if r.log[i].Command.ID != other.log[i].Command.ID {
   return false
  }
 }

 return true
}

// MemoryState implements State interface for in-memory state
type MemoryState struct {
 mu     sync.RWMutex
 data   map[string][]byte
}

// NewMemoryState creates new memory state
func NewMemoryState() *MemoryState {
 return &MemoryState{
  data: make(map[string][]byte),
 }
}

// Apply applies command to state
func (s *MemoryState) Apply(cmd Command) (Result, error) {
 s.mu.Lock()
 defer s.mu.Unlock()

 switch cmd.Op {
 case OpRead:
  val, ok := s.data[string(cmd.Args)]
  if !ok {
   return Result{}, fmt.Errorf("key not found")
  }
  return Result{Value: val}, nil

 case OpWrite:
  if len(cmd.Args) < 2 {
   return Result{}, fmt.Errorf("invalid write args")
  }
  key := string(cmd.Args[0])
  val := cmd.Args[1]
  s.data[key] = val
  return Result{Value: val}, nil

 case OpDelete:
  delete(s.data, string(cmd.Args))
  return Result{}, nil

 default:
  return Result{}, fmt.Errorf("unknown operation")
 }
}

// Clone clones the state
func (s *MemoryState) Clone() State {
 s.mu.RLock()
 defer s.mu.RUnlock()

 clone := NewMemoryState()
 for k, v := range s.data {
  clone.data[k] = v
 }
 return clone
}
```

---

## 6. Visual Representations

### 6.1 RSM State Transitions

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    RSM STATE TRANSITIONS                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐                                                             │
│  │   Initial   │                                                             │
│  │    State    │                                                             │
│  └──────┬──────┘                                                             │
│         │ Apply(Command 1)                                                   │
│         ▼                                                                    │
│  ┌─────────────┐                                                             │
│  │   State 1   │                                                             │
│  └──────┬──────┘                                                             │
│         │ Apply(Command 2)                                                   │
│         ▼                                                                    │
│  ┌─────────────┐                                                             │
│  │   State 2   │                                                             │
│  └──────┬──────┘                                                             │
│         │ Apply(Command 3)                                                   │
│         ▼                                                                    │
│  ┌─────────────┐                                                             │
│  │   State 3   │                                                             │
│  └─────────────┘                                                             │
│                                                                              │
│  Key Property:                                                               │
│  δ(δ(δ(S0, C1), C2), C3) = δ(δ(δ(S0, C1), C2), C3)                         │
│  All replicas reach State 3 given same command sequence                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 RSM Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    RSM ARCHITECTURE OVERVIEW                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Client Layer                                  │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                             │   │
│  │  │Client 1 │  │Client 2 │  │Client 3 │                             │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘                             │   │
│  └───────┼────────────┼────────────┼──────────────────────────────────┘   │
│          │            │            │                                        │
│          └────────────┴────────────┘                                        │
│                     │                                                       │
│                     ▼                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Consensus Layer                                 │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │   │
│  │  │ Replica │  │ Replica │  │ Replica │  │ Replica │  │ Replica │   │   │
│  │  │    1    │  │    2    │  │    3    │  │    4    │  │    5    │   │   │
│  │  │ (Leader)│  │         │  │         │  │         │  │         │   │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘   │   │
│  │       └────────────┴────────────┴────────────┴────────────┘        │   │
│  │                         Consensus Log                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                     │                                                       │
│                     ▼ Apply (in log order)                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Application Layer                                 │   │
│  │                                                                     │   │
│  │  ┌─────────────────────────────────────────────────────────────┐   │   │
│  │  │              Replicated State Machine                        │   │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐        │   │   │
│  │  │  │ State 1 │─>│ State 2 │─>│ State 3 │─>│ State 4 │ ...    │   │   │
│  │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘        │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Failure Handling in RSM

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FAILURE HANDLING IN RSM                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Normal Operation:                                                           │
│  ────────────────                                                            │
│                                                                              │
│  Leader (R1)    R2        R3        R4        R5                            │
│    │            │          │          │          │                          │
│    │─Propose──>│          │          │          │                          │
│    │─Propose──>│          │          │          │                          │
│    │           │          │          │          │                          │
│    │<─Accept──│          │          │          │                          │
│    │<─Accept──│          │          │          │                          │
│    │ ...      │          │          │          │                          │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────           │
│                                                                              │
│  Leader Failure:                                                             │
│  ───────────────                                                             │
│                                                                              │
│  Leader (R1)    R2        R3        R4        R5                            │
│    X            │          │          │          │  [R1 crashes]             │
│                 │          │          │          │                          │
│                 │─Timeout  │          │          │  [R2 times out]          │
│                 │          │          │          │                          │
│                 │─ViewChange────────────────────>│  [R2 starts view change] │
│                 │          │          │          │                          │
│                 │<────────ViewChange─────────────│  [Others respond]        │
│                 │          │          │          │                          │
│                 │─NewView───────────────────────>│  [R2 becomes leader]     │
│                 │          │          │          │                          │
│  New Leader (R2) R3       R4        R5          │  [Continue operation]    │
│    │            │          │          │          │                          │
│    │─Propose──>│          │          │          │                          │
│                                                                              │
│  Properties Maintained:                                                      │
│  ✓ Safety: Committed commands preserved                                       │
│  ✓ Liveness: New leader elected, progress continues                          │
│  ✓ Consistency: All replicas agree on committed prefix                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    REPLICATED STATE MACHINE SUMMARY                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Core Principle:                                                             │
│  If all replicas start in same state and apply same deterministic           │
│  commands in same order, they remain in identical states.                   │
│                                                                              │
│  Implementation Requirements:                                                │
│  1. Deterministic state machine                                              │
│  2. Consensus for command ordering                                          │
│  3. Log replication                                                          │
│  4. Failure detection and recovery                                           │
│                                                                              │
│  Guarantees:                                                                 │
│  • Safety: All replicas agree (despite failures)                            │
│  • Liveness: Commands eventually applied (if majority correct)              │
│  • Fault Tolerance: Tolerates f failures with 2f+1 replicas                 │
│                                                                              │
│  Applications:                                                               │
│  • Distributed databases                                                     │
│  • Coordination services                                                     │
│  • Configuration management                                                  │
│  • Leader election                                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
