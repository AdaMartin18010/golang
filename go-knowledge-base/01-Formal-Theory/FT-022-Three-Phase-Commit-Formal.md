# FT-022: Three-Phase Commit (3PC) - Formal Specification

## Overview

Three-Phase Commit (3PC) is a distributed consensus protocol that extends Two-Phase Commit (2PC) to handle coordinator failure without blocking. It adds an additional pre-commit phase, allowing participants to make progress even if the coordinator fails after making a decision.

## Theoretical Foundations

### 1.1 System Model

**Distributed System**:

```
System T = ⟨C, P, Δ, f⟩

where:
- C: coordinator node
- P = {p₁, p₂, ..., pₙ}: set of participant nodes (n ≥ 2)
- Δ: maximum message delay
- f: maximum number of failures (f < n/2)
```

**Failure Model**:

```
- Fail-stop failures: nodes halt and don't resume
- Network partitions: partial communication failures
- Timeout-based failure detection
```

### 1.2 3PC Protocol Phases

**Phase 1: CanCommit (Voting)**

```
Coordinator:
  1. Write START-3PC to stable storage
  2. Send CAN_COMMIT? to all participants
  3. Wait for responses (YES/NO)

Participant pᵢ (on receiving CAN_COMMIT?):
  1. Check local feasibility:
     - Can acquire necessary locks
     - Local resources available
     - No conflicts detected
  2. If feasible:
     - Reply YES
     - Enter UNCERTAIN state
  3. If not feasible:
     - Reply NO
     - Enter ABORTED state
```

**Phase 2: PreCommit**

```
Coordinator (if all YES received):
  1. Write PRE-COMMIT to stable storage
  2. Send PRE_COMMIT to all participants
  3. Enter COMMITTING state
  4. Wait for ACKs

Coordinator (if any NO or timeout):
  1. Write ABORT to stable storage
  2. Send ABORT to all participants
  3. Enter ABORTED state

Participant pᵢ (on receiving PRE_COMMIT):
  1. Write PRE-COMMIT to stable storage
  2. Send ACK to coordinator
  3. Enter COMMITTABLE state
  4. Start timeout timer

Participant pᵢ (on receiving ABORT):
  1. Write ABORT to stable storage
  2. Abort transaction locally
  3. Enter ABORTED state
```

**Phase 3: DoCommit**

```
Coordinator (after all ACKs):
  1. Write COMMIT to stable storage
  2. Send COMMIT to all participants
  3. Enter COMMITTED state

Participant pᵢ (on receiving COMMIT):
  1. Write COMMIT to stable storage
  2. Commit transaction
  3. Enter COMMITTED state
  4. Send ACK to coordinator

Coordinator (after all ACKs):
  1. Write COMPLETE to stable storage
  2. Clean up
```

### 1.3 State Machine

```
Coordinator States:
  INIT ──→ PREPARING ──→ COMMITTING ──→ COMMITTED
    │          │              │
    └──────────┴──────────────┴──→ ABORTED

Participant States:
  INIT ──→ UNCERTAIN ──→ COMMITTABLE ──→ COMMITTED
    │           │             │
    └───────────┴─────────────┴──→ ABORTED
```

### 1.4 Non-Blocking Property Proof

```
Theorem (Non-Blocking): 3PC is non-blocking, meaning that if
coordinator fails, participants can reach a decision without
waiting for coordinator recovery.

Proof by Cases:

Case 1: Coordinator fails before sending PRE_COMMIT
  - Participants are in UNCERTAIN state (voted YES but no PRE_COMMIT)
  - Participants timeout waiting for PRE_COMMIT
  - Participants can safely abort (coordinator couldn't have decided COMMIT)
  - No blocking

Case 2: Coordinator fails after sending PRE_COMMIT but before COMMIT
  - Some participants may be in COMMITTABLE state
  - If any participant receives PRE_COMMIT, all did (reliable broadcast)
  - Participants in COMMITTABLE can contact other participants
  - If any participant is COMMITTED, all commit
  - If all participants are COMMITTABLE, they can safely commit
  - No blocking

Case 3: Coordinator fails after sending COMMIT
  - Some participants may have committed
  - Participants that haven't committed can learn decision from others
  - No blocking

Therefore, 3PC is non-blocking. ∎

Theorem (Atomicity): All participants either commit or abort together.

Proof:
Coordinator only sends COMMIT if:
  1. All participants voted YES
  2. All participants acknowledged PRE_COMMIT

If coordinator decides COMMIT:
  - All participants receive COMMIT (or learn from others)
  - All commit

If coordinator decides ABORT:
  - ABORT can be sent at multiple points
  - All participants abort

Therefore, atomicity is maintained. ∎
```

### 1.5 Comparison with 2PC

| Property | 2PC | 3PC |
|----------|-----|-----|
| **Blocking** | Yes (coordinator failure) | No |
| **Message Rounds** | 2 | 3 |
| **Log Writes** | 2-3 per participant | 3 per participant |
| **Latency** | Lower | Higher |
| **Fault Tolerance** | Lower | Higher |
| **Complexity** | Simpler | More complex |

**Blocking Scenario in 2PC**:

```
Participant p votes YES and enters PREPARED state
Coordinator fails before sending decision
p cannot unilaterally decide (coordinator might have told others to commit)
p blocks until coordinator recovers
```

**Why 3PC is Non-Blocking**:

```
The PRE_COMMIT phase creates a window where:
- If any participant has PRE_COMMIT, all will eventually get it
- If coordinator fails after PRE_COMMIT, participants can:
  - Query other participants
  - If any has COMMIT, all commit
  - If all have COMMITTABLE, commit is safe
  - If some don't have PRE_COMMIT, abort is safe
```

## TLA+ Specification

```tla
----------------------------- MODULE ThreePhaseCommit -----------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Participants,      \* Set of participant IDs
          Coordinator,       \* Coordinator ID
          MaxFailures        \* Maximum failures to tolerate

VARIABLES cState,            \* Coordinator state
          pState,            \* Participant states
          cLog,              \* Coordinator log
          pLog,              \* Participant logs
          messages,          \* In-flight messages
          decisions,         \* Final decisions
          failed             \* Set of failed nodes

\* States
CoordinatorState == {"init", "preparing", "precommitted", "committing",
                     "committed", "aborted", "completed"}
ParticipantState == {"init", "uncertain", "committable", "committed", "aborted"}

\* Initial state
Init ==
  ∧ cState = [state ↦ "init", tx ↦ 0, acks ↦ 0]
  ∧ pState = [p ∈ Participants ↦ [state ↦ "init", vote ↦ ⊥]]
  ∧ cLog = ⟨⟩
  ∧ pLog = [p ∈ Participants ↦ ⟨⟩]
  ∧ messages = {}
  ∧ decisions = [p ∈ Participants ↦ ⊥]
  ∧ failed = {}

\* Phase 1: Coordinator sends CAN_COMMIT
CoordinatorCanCommit ==
  ∧ cState.state = "init"
  ∧ Coordinator ∉ failed
  ∧ cState' = [cState EXCEPT !.state = "preparing", !.tx = @ + 1]
  ∧ cLog' = Append(cLog, [type ↦ "start-3pc", tx ↦ cState.tx + 1])
  ∧ messages' = messages ∪
       {[type ↦ "can_commit", from ↦ Coordinator, to ↦ p, tx ↦ cState.tx + 1] :
        p ∈ Participants \\ failed}
  ∧ UNCHANGED ⟨pState, pLog, decisions, failed⟩

\* Participant receives CAN_COMMIT and votes
ParticipantVote(p) ==
  ∧ p ∉ failed
  ∧ pState[p].state = "init"
  ∧ ∃m ∈ messages:
      ∧ m.type = "can_commit"
      ∧ m.to = p
      ∧ IF CanCommit(p)  \* Abstract predicate
        THEN ∧ pState' = [pState EXCEPT ![p] = [state ↦ "uncertain", vote ↦ "yes"]]
             ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "uncertain", tx ↦ m.tx])]
             ∧ messages' = messages ∪
                  {[type ↦ "yes", from ↦ p, to ↦ Coordinator, tx ↦ m.tx]}
        ELSE ∧ pState' = [pState EXCEPT ![p] = [state ↦ "aborted", vote ↦ "no"]]
             ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "aborted", tx ↦ m.tx])]
             ∧ decisions' = [decisions EXCEPT ![p] = "abort"]
             ∧ messages' = messages ∪
                  {[type ↦ "no", from ↦ p, to ↦ Coordinator, tx ↦ m.tx]}
  ∧ UNCHANGED ⟨cState, cLog, failed⟩

\* Phase 2a: Coordinator sends PRE_COMMIT (all voted yes)
CoordinatorPreCommit ==
  ∧ cState.state = "preparing"
  ∧ Coordinator ∉ failed
  ∧ LET votes == {m ∈ messages :
                   m.type = "yes" ∧ m.to = Coordinator ∧ m.tx = cState.tx}
    IN Cardinality(votes) = Cardinality(Participants)
  ∧ cState' = [cState EXCEPT !.state = "precommitted"]
  ∧ cLog' = Append(cLog, [type ↦ "pre-commit", tx ↦ cState.tx])
  ∧ messages' = messages ∪
       {[type ↦ "pre_commit", from ↦ Coordinator, to ↦ p, tx ↦ cState.tx] :
        p ∈ Participants \\ failed}
  ∧ UNCHANGED ⟨pState, pLog, decisions, failed⟩

\* Phase 2b: Coordinator aborts (any voted no or timeout)
CoordinatorAbortPhase2 ==
  ∧ cState.state = "preparing"
  ∧ Coordinator ∉ failed
  ∧ (∃m ∈ messages : m.type = "no" ∧ m.to = Coordinator ∧ m.tx = cState.tx)
  ∧ cState' = [cState EXCEPT !.state = "aborted"]
  ∧ cLog' = Append(cLog, [type ↦ "abort", tx ↦ cState.tx])
  ∧ messages' = messages ∪
       {[type ↦ "abort", from ↦ Coordinator, to ↦ p, tx ↦ cState.tx] :
        p ∈ Participants \\ failed}
  ∧ UNCHANGED ⟨pState, pLog, decisions, failed⟩

\* Participant receives PRE_COMMIT
ParticipantPreCommit(p) ==
  ∧ p ∉ failed
  ∧ pState[p].state = "uncertain"
  ∧ ∃m ∈ messages:
      ∧ m.type = "pre_commit"
      ∧ m.to = p
      ∧ m.tx = cState.tx
      ∧ pState' = [pState EXCEPT ![p] = [state ↦ "committable", vote ↦ "ack"]]
      ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "committable", tx ↦ m.tx])]
      ∧ messages' = messages ∪
           {[type ↦ "ack", from ↦ p, to ↦ Coordinator, tx ↦ m.tx]}
  ∧ UNCHANGED ⟨cState, cLog, decisions, failed⟩

\* Phase 3: Coordinator sends COMMIT
CoordinatorDoCommit ==
  ∧ cState.state = "precommitted"
  ∧ Coordinator ∉ failed
  ∧ LET acks == {m ∈ messages :
                  m.type = "ack" ∧ m.to = Coordinator ∧ m.tx = cState.tx}
    IN Cardinality(acks) = Cardinality(Participants)
  ∧ cState' = [cState EXCEPT !.state = "committed"]
  ∧ cLog' = Append(cLog, [type ↦ "commit", tx ↦ cState.tx])
  ∧ messages' = messages ∪
       {[type ↦ "commit", from ↦ Coordinator, to ↦ p, tx ↦ cState.tx] :
        p ∈ Participants \\ failed}
  ∧ UNCHANGED ⟨pState, pLog, decisions, failed⟩

\* Participant receives COMMIT
ParticipantCommit(p) ==
  ∧ p ∉ failed
  ∧ pState[p].state = "committable"
  ∧ ∃m ∈ messages:
      ∧ m.type = "commit"
      ∧ m.to = p
      ∧ m.tx = cState.tx
      ∧ pState' = [pState EXCEPT ![p].state = "committed"]
      ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "committed", tx ↦ m.tx])]
      ∧ decisions' = [decisions EXCEPT ![p] = "commit"]
  ∧ UNCHANGED ⟨cState, cLog, messages, failed⟩

\* Participant receives ABORT
ParticipantAbort(p) ==
  ∧ p ∉ failed
  ∧ pState[p].state ∈ {"uncertain", "committable"}
  ∧ ∃m ∈ messages:
      ∧ m.type = "abort"
      ∧ m.to = p
      ∧ m.tx = cState.tx
      ∧ pState' = [pState EXCEPT ![p].state = "aborted"]
      ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "aborted", tx ↦ m.tx])]
      ∧ decisions' = [decisions EXCEPT ![p] = "abort"]
  ∧ UNCHANGED ⟨cState, cLog, messages, failed⟩

\* Failure transitions
CoordinatorFails ==
  ∧ Coordinator ∉ failed
  ∧ Cardinality(failed) < MaxFailures
  ∧ failed' = failed ∪ {Coordinator}
  ∧ UNCHANGED ⟨cState, pState, cLog, pLog, messages, decisions⟩

ParticipantFails(p) ==
  ∧ p ∉ failed
  ∧ Cardinality(failed) < MaxFailures
  ∧ pState[p].state ∉ {"committed", "aborted"}
  ∧ failed' = failed ∪ {p}
  ∧ UNCHANGED ⟨cState, pState, cLog, pLog, messages, decisions⟩

\* Non-blocking recovery (participants take initiative)
ParticipantTimeout(p) ==
  ∧ p ∉ failed
  ∧ pState[p].state = "uncertain"
  ∧ Coordinator ∈ failed
  ∧ pState' = [pState EXCEPT ![p].state = "aborted"]
  ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "timeout-abort", tx ↦ cState.tx])]
  ∧ decisions' = [decisions EXCEPT ![p] = "abort"]
  ∧ UNCHANGED ⟨cState, cLog, messages, failed⟩

ParticipantCooperativeDecision(p) ==
  ∧ p ∉ failed
  ∧ pState[p].state = "committable"
  ∧ Coordinator ∈ failed
  ∧ (∃q ∈ Participants \\ failed :
       pState[q].state = "committed" ∨ decisions[q] = "commit")
  ∧ pState' = [pState EXCEPT ![p].state = "committed"]
  ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "coop-commit", tx ↦ cState.tx])]
  ∧ decisions' = [decisions EXCEPT ![p] = "commit"]
  ∧ UNCHANGED ⟨cState, cLog, messages, failed⟩

\* Next state
Next ==
  ∨ CoordinatorCanCommit
  ∨ ∃p ∈ Participants: ParticipantVote(p)
  ∨ CoordinatorPreCommit
  ∨ CoordinatorAbortPhase2
  ∨ ∃p ∈ Participants: ParticipantPreCommit(p)
  ∨ CoordinatorDoCommit
  ∨ ∃p ∈ Participants: ParticipantCommit(p)
  ∨ ∃p ∈ Participants: ParticipantAbort(p)
  ∨ CoordinatorFails
  ∨ ∃p ∈ Participants: ParticipantFails(p)
  ∨ ∃p ∈ Participants: ParticipantTimeout(p)
  ∨ ∃p ∈ Participants: ParticipantCooperativeDecision(p)

\* Invariants

\* All participants reach same decision
Atomicity ==
  ∀p1, p2 ∈ Participants \\ failed:
    decisions[p1] ≠ ⊥ ∧ decisions[p2] ≠ ⊥ ⇒ decisions[p1] = decisions[p2]

\* Non-blocking: no participant blocks indefinitely
NonBlocking ==
  ∀p ∈ Participants:
    Coordinator ∈ failed ∧ p ∉ failed
    ⇒ ◇(decisions[p] ≠ ⊥)

=============================================================================
```

## Algorithm Pseudocode

### Classic 3PC

```
Algorithm: Classic Three-Phase Commit

Types:
  State: INIT | UNCERTAIN | COMMITTABLE | COMMITTED | ABORTED
  Decision: COMMIT | ABORT | PRE_COMMIT

Coordinator:
  Variables:
    state: State = INIT
    votes: Map<Participant, Vote>
    acks: Set<Participant>

  BeginTransaction(txId, participants):
    WriteLog(START, txId)
    state = PREPARING
    for each p in participants:
      Send(CAN_COMMIT, txId) to p
    StartTimer(TIMEOUT_PHASE1)

  On Receive(VOTE, txId, p, vote):
    if state ≠ PREPARING: return
    votes[p] = vote

    if all participants voted:
      CancelTimer()
      if all votes are YES:
        EnterPhase2Commit(txId)
      else:
        EnterPhase2Abort(txId)

  EnterPhase2Commit(txId):
    WriteLog(PRE_COMMIT, txId)
    state = PRE_COMMIT
    for each p in participants:
      Send(PRE_COMMIT, txId) to p
    StartTimer(TIMEOUT_PHASE2)

  EnterPhase2Abort(txId):
    WriteLog(ABORT, txId)
    state = ABORTED
    for each p in participants:
      Send(ABORT, txId) to p

  On Receive(ACK, txId, p):
    if state ≠ PRE_COMMIT: return
    acks.add(p)

    if all participants acked:
      CancelTimer()
      EnterPhase3(txId)

  EnterPhase3(txId):
    WriteLog(COMMIT, txId)
    state = COMMITTED
    for each p in participants:
      Send(COMMIT, txId) to p
    StartTimer(TIMEOUT_PHASE3)

  On Receive(FINAL_ACK, txId, p):
    if state ≠ COMMITTED: return
    // All done
    WriteLog(COMPLETE, txId)
    state = COMPLETED

  On Timer(timeout):
    if timeout == TIMEOUT_PHASE1:
      // Some didn't vote, abort
      EnterPhase2Abort(txId)
    else if timeout == TIMEOUT_PHASE2:
      // Some didn't ack, but we can still commit
      EnterPhase3(txId)
    else if timeout == TIMEOUT_PHASE3:
      // Re-send COMMIT
      for each p in participants:
        Send(COMMIT, txId) to p
      RestartTimer(TIMEOUT_PHASE3)

Participant:
  Variables:
    state: State = INIT

  On Receive(CAN_COMMIT, txId):
    if not CanExecute(txId):
      Send(VOTE, txId, self, NO) to coordinator
      state = ABORTED
      return

    PrepareResources(txId)
    Send(VOTE, txId, self, YES) to coordinator
    state = UNCERTAIN
    StartTimer(TIMEOUT_DECISION)

  On Receive(PRE_COMMIT, txId):
    CancelTimer()
    WriteLog(PRE_COMMIT, txId)
    Send(ACK, txId, self) to coordinator
    state = COMMITTABLE
    StartTimer(TIMEOUT_COMMIT)

  On Receive(COMMIT, txId):
    CancelTimer()
    WriteLog(COMMIT, txId)
    Commit(txId)
    state = COMMITTED
    Send(FINAL_ACK, txId, self) to coordinator

  On Receive(ABORT, txId):
    CancelTimer()
    WriteLog(ABORT, txId)
    Abort(txId)
    state = ABORTED
    Send(FINAL_ACK, txId, self) to coordinator

  On Timer(timeout):
    if timeout == TIMEOUT_DECISION:
      // Coordinator may have failed
      if state == UNCERTAIN:
        // Safe to abort - coordinator couldn't have committed
        Abort(txId)
        state = ABORTED
      else if state == COMMITTABLE:
        // Contact other participants
        CooperativeDecision()
    else if timeout == TIMEOUT_COMMIT:
      // Re-contact coordinator or participants
      CooperativeDecision()

  CooperativeDecision():
    // Query other participants
    for each p in participants where p ≠ self:
      status = QueryParticipant(p)
      if status == COMMITTED:
        Commit(txId)
        state = COMMITTED
        return
      if status == ABORTED:
        Abort(txId)
        state = ABORTED
        return

    // If all are COMMITTABLE, we can safely commit
    // If any are UNCERTAIN, coordinator didn't pre-commit, safe to abort
```

### Optimizations

```
Optimization 1: Early Abort
If a participant knows it cannot commit during Phase 1,
it can abort immediately without waiting for coordinator.
Other participants will eventually timeout and abort.

Optimization 2: Cooperative Recovery
When coordinator fails, participants use consensus protocol
to agree on outcome. Options:
  - Paxos/Raft among participants
  - Byzantine agreement for malicious failures
  - Simple majority vote (if all participants are COMMITTABLE)

Optimization 3: Reduced Logging
- Write logs asynchronously
- Batch multiple transactions
- Use group commit for disk writes
```

## Go Implementation

```go
// Package threepc implements the Three-Phase Commit protocol
package threepc

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// TransactionState represents the state in 3PC
type TransactionState int

const (
 StateInit TransactionState = iota
 StatePreparing
 StateUncertain
 StatePreCommitted
 StateCommittable
 StateCommitting
 StateCommitted
 StateAborting
 StateAborted
 StateCompleted
)

func (s TransactionState) String() string {
 switch s {
 case StateInit:
  return "INIT"
 case StatePreparing:
  return "PREPARING"
 case StateUncertain:
  return "UNCERTAIN"
 case StatePreCommitted:
  return "PRE-COMMITTED"
 case StateCommittable:
  return "COMMITTABLE"
 case StateCommitting:
  return "COMMITTING"
 case StateCommitted:
  return "COMMITTED"
 case StateAborting:
  return "ABORTING"
 case StateAborted:
  return "ABORTED"
 case StateCompleted:
  return "COMPLETED"
 default:
  return "UNKNOWN"
 }
}

// Phase represents the current phase of 3PC
type Phase int

const (
 PhaseCanCommit Phase = iota
 PhasePreCommit
 PhaseDoCommit
)

// Coordinator manages 3PC transactions
type Coordinator struct {
 mu           sync.RWMutex
 id           string
 participants []string
 network      Network
 log          StableLog
 timeouts     map[string]*time.Timer
 transactions map[string]*CoordinatorTx
}

// CoordinatorTx tracks a transaction at coordinator
type CoordinatorTx struct {
 TxID         string
 Phase        Phase
 State        TransactionState
 Votes        map[string]bool
 Acks         map[string]bool
 mu           sync.Mutex
}

// Network interface for communication
type Network interface {
 Send(msg Message) error
 Broadcast(msg Message, exclude []string) error
}

// StableLog interface for persistence
type StableLog interface {
 Append(entry LogEntry) error
 Read(txID string) ([]LogEntry, error)
}

// LogEntry represents a log record
type LogEntry struct {
 TxID      string
 Type      string
 Timestamp time.Time
}

// Message represents protocol messages
type Message struct {
 Type       string
 TxID       string
 From       string
 To         string
 Vote       bool
 Ack        bool
 Decision   string
}

// NewCoordinator creates a new 3PC coordinator
func NewCoordinator(id string, participants []string, network Network, log StableLog) *Coordinator {
 return &Coordinator{
  id:           id,
  participants: participants,
  network:      network,
  log:          log,
  timeouts:     make(map[string]*time.Timer),
  transactions: make(map[string]*CoordinatorTx),
 }
}

// BeginTransaction starts a new 3PC transaction
func (c *Coordinator) BeginTransaction(txID string) error {
 c.mu.Lock()
 defer c.mu.Unlock()

 // Write start log
 if err := c.log.Append(LogEntry{
  TxID:      txID,
  Type:      "START-3PC",
  Timestamp: time.Now(),
 }); err != nil {
  return err
 }

 // Create transaction
 tx := &CoordinatorTx{
  TxID:  txID,
  Phase: PhaseCanCommit,
  State: StatePreparing,
  Votes: make(map[string]bool),
  Acks:  make(map[string]bool),
 }
 c.transactions[txID] = tx

 // Send CAN_COMMIT to all participants
 for _, p := range c.participants {
  msg := Message{
   Type: "CAN_COMMIT",
   TxID: txID,
   From: c.id,
   To:   p,
  }
  if err := c.network.Send(msg); err != nil {
   return err
  }
 }

 // Start timeout
 c.timeouts[txID] = time.AfterFunc(30*time.Second, func() {
  c.handleTimeout(txID)
 })

 return nil
}

// HandleVote processes a vote from participant
func (c *Coordinator) HandleVote(msg Message) error {
 c.mu.Lock()
 tx, exists := c.transactions[msg.TxID]
 c.mu.Unlock()

 if !exists {
  return fmt.Errorf("unknown transaction: %s", msg.TxID)
 }

 tx.mu.Lock()
 defer tx.mu.Unlock()

 if tx.Phase != PhaseCanCommit {
  return nil // Late vote
 }

 tx.Votes[msg.From] = msg.Vote

 // Check if all voted
 if len(tx.Votes) == len(c.participants) {
  // Cancel timeout
  if timer, ok := c.timeouts[msg.TxID]; ok {
   timer.Stop()
  }

  // Count votes
  allYes := true
  for _, vote := range tx.Votes {
   if !vote {
    allYes = false
    break
   }
  }

  if allYes {
   return c.preCommit(tx)
  }
  return c.abort(tx)
 }

 return nil
}

// preCommit enters Phase 2
func (c *Coordinator) preCommit(tx *CoordinatorTx) error {
 tx.Phase = PhasePreCommit
 tx.State = StatePreCommitted

 // Write pre-commit log
 if err := c.log.Append(LogEntry{
  TxID:      tx.TxID,
  Type:      "PRE-COMMIT",
  Timestamp: time.Now(),
 }); err != nil {
  return err
 }

 // Send PRE_COMMIT to all
 for _, p := range c.participants {
  c.network.Send(Message{
   Type: "PRE_COMMIT",
   TxID: tx.TxID,
   From: c.id,
   To:   p,
  })
 }

 // Start new timeout
 c.timeouts[tx.TxID] = time.AfterFunc(30*time.Second, func() {
  c.handleTimeout(tx.TxID)
 })

 return nil
}

// doCommit enters Phase 3
func (c *Coordinator) doCommit(tx *CoordinatorTx) error {
 tx.Phase = PhaseDoCommit
 tx.State = StateCommitting

 // Write commit log
 if err := c.log.Append(LogEntry{
  TxID:      tx.TxID,
  Type:      "COMMIT",
  Timestamp: time.Now(),
 }); err != nil {
  return err
 }

 // Send COMMIT to all
 for _, p := range c.participants {
  c.network.Send(Message{
   Type:     "COMMIT",
   TxID:     tx.TxID,
   From:     c.id,
   To:       p,
   Decision: "COMMIT",
  })
 }

 return nil
}

// abort aborts the transaction
func (c *Coordinator) abort(tx *CoordinatorTx) error {
 tx.State = StateAborting

 // Write abort log
 if err := c.log.Append(LogEntry{
  TxID:      tx.TxID,
  Type:      "ABORT",
  Timestamp: time.Now(),
 }); err != nil {
  return err
 }

 // Send ABORT to all
 for _, p := range c.participants {
  c.network.Send(Message{
   Type:     "ABORT",
   TxID:     tx.TxID,
   From:     c.id,
   To:       p,
   Decision: "ABORT",
  })
 }

 return nil
}

func (c *Coordinator) handleTimeout(txID string) {
 c.mu.Lock()
 tx, exists := c.transactions[txID]
 c.mu.Unlock()

 if !exists {
  return
 }

 tx.mu.Lock()
 defer tx.mu.Unlock()

 switch tx.Phase {
 case PhaseCanCommit:
  // Some didn't vote, abort
  c.abort(tx)
 case PhasePreCommit:
  // Some didn't ack, but we can still commit
  c.doCommit(tx)
 case PhaseDoCommit:
  // Re-send commit
  for _, p := range c.participants {
   c.network.Send(Message{
    Type:     "COMMIT",
    TxID:     tx.TxID,
    From:     c.id,
    To:       p,
    Decision: "COMMIT",
   })
  }
 }
}

// Participant represents a 3PC participant
type Participant struct {
 mu            sync.RWMutex
 id            string
 coordinator   string
 network       Network
 log           StableLog
 timeouts      map[string]*time.Timer
 transactions  map[string]*ParticipantTx
 queryFunc     func(string) (string, error)
}

// ParticipantTx tracks participant's transaction
type ParticipantTx struct {
 TxID  string
 State TransactionState
}

// NewParticipant creates a new participant
func NewParticipant(id, coordinator string, network Network, log StableLog) *Participant {
 return &Participant{
  id:           id,
  coordinator:  coordinator,
  network:      network,
  log:          log,
  timeouts:     make(map[string]*time.Timer),
  transactions: make(map[string]*ParticipantTx),
 }
}

// HandleCanCommit handles CAN_COMMIT message
func (p *Participant) HandleCanCommit(msg Message, canExecute func(string) bool) error {
 p.mu.Lock()
 defer p.mu.Unlock()

 if canExecute != nil && !canExecute(msg.TxID) {
  // Cannot execute, vote NO
  p.network.Send(Message{
   Type: "VOTE",
   TxID: msg.TxID,
   From: p.id,
   To:   msg.From,
   Vote: false,
  })
  return nil
 }

 // Vote YES
 tx := &ParticipantTx{
  TxID:  msg.TxID,
  State: StateUncertain,
 }
 p.transactions[msg.TxID] = tx

 p.network.Send(Message{
  Type: "VOTE",
  TxID: msg.TxID,
  From: p.id,
  To:   msg.From,
  Vote: true,
 })

 // Start timeout
 p.timeouts[msg.TxID] = time.AfterFunc(30*time.Second, func() {
  p.handleTimeout(msg.TxID)
 })

 return nil
}

// HandlePreCommit handles PRE_COMMIT message
func (p *Participant) HandlePreCommit(msg Message) error {
 p.mu.Lock()
 defer p.mu.Unlock()

 tx, exists := p.transactions[msg.TxID]
 if !exists {
  return fmt.Errorf("unknown transaction")
 }

 // Cancel timeout
 if timer, ok := p.timeouts[msg.TxID]; ok {
  timer.Stop()
 }

 // Write log
 p.log.Append(LogEntry{
  TxID:      msg.TxID,
  Type:      "PRE-COMMIT",
  Timestamp: time.Now(),
 })

 tx.State = StateCommittable

 // Send ACK
 p.network.Send(Message{
  Type: "ACK",
  TxID: msg.TxID,
  From: p.id,
  To:   msg.From,
  Ack:  true,
 })

 // Start new timeout
 p.timeouts[msg.TxID] = time.AfterFunc(30*time.Second, func() {
  p.handleTimeout(msg.TxID)
 })

 return nil
}

// HandleCommit handles COMMIT message
func (p *Participant) HandleCommit(msg Message, commitFunc func(string) error) error {
 p.mu.Lock()
 defer p.mu.Unlock()

 tx, exists := p.transactions[msg.TxID]
 if !exists {
  return fmt.Errorf("unknown transaction")
 }

 // Cancel timeout
 if timer, ok := p.timeouts[msg.TxID]; ok {
  timer.Stop()
 }

 // Execute commit
 if commitFunc != nil {
  commitFunc(msg.TxID)
 }

 // Write log
 p.log.Append(LogEntry{
  TxID:      msg.TxID,
  Type:      "COMMITTED",
  Timestamp: time.Now(),
 })

 tx.State = StateCommitted

 return nil
}

func (p *Participant) handleTimeout(txID string) {
 p.mu.Lock()
 tx, exists := p.transactions[txID]
 p.mu.Unlock()

 if !exists {
  return
 }

 switch tx.State {
 case StateUncertain:
  // Safe to abort - coordinator couldn't have committed
  p.log.Append(LogEntry{
   TxID:      txID,
   Type:      "ABORTED",
   Timestamp: time.Now(),
  })
  tx.State = StateAborted
 case StateCommittable:
  // Need to query other participants
  p.cooperativeDecision(txID)
 }
}

func (p *Participant) cooperativeDecision(txID string) {
 // In a real implementation, query other participants
 // For now, commit if possible
 p.log.Append(LogEntry{
  TxID:      txID,
  Type:      "COMMITTED",
  Timestamp: time.Now(),
 })
}
