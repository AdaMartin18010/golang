# FT-021: Two-Phase Commit (2PC) - Formal Specification

## Overview

Two-Phase Commit (2PC) is a distributed algorithm that ensures all nodes in a distributed system agree to commit or abort a transaction. It is a blocking consensus protocol that guarantees atomicity across multiple nodes but is vulnerable to coordinator failure.

## Theoretical Foundations

### 1.1 System Model

**Distributed Transaction System**:

```
System T = ⟨C, P, M⟩

where:
- C: coordinator node
- P = {p₁, p₂, ..., pₙ}: set of participant nodes
- M: set of messages between coordinator and participants
```

**Transaction State**:

```
Transaction t has:
- Coordinator state: {INIT, PREPARING, COMMITTED, ABORTED}
- Participant state: {INIT, PREPARED, COMMITTED, ABORTED}
- Decision: COMMIT or ABORT
```

### 1.2 2PC Protocol Phases

**Phase 1: Prepare Phase (Voting)**

```
Coordinator:
  1. Write START-2PC to stable storage
  2. Send PREPARE message to all participants
  3. Wait for votes from all participants

Participant pᵢ (on receiving PREPARE):
  1. Check if can commit:
     - Locks acquired
     - Constraints satisfied
     - No local failures
  2. If can commit:
     - Write PREPARED to stable storage
     - Send YES vote to coordinator
     - Enter PREPARED state (blocking)
  3. If cannot commit:
     - Write ABORTED to stable storage
     - Send NO vote to coordinator
     - Abort transaction locally
```

**Phase 2: Commit/Abort Phase (Decision)**

```
Coordinator (after collecting all votes):
  Case 1: All votes are YES
    1. Write COMMIT to stable storage
    2. Send COMMIT message to all participants
    3. Enter COMMITTED state

  Case 2: Any vote is NO or timeout
    1. Write ABORT to stable storage
    2. Send ABORT message to all participants
    3. Enter ABORTED state

Participant (on receiving decision):
  On COMMIT:
    1. Write COMMITTED to stable storage
    2. Commit transaction locally
    3. Release locks
    4. Send ACK to coordinator

  On ABORT:
    1. Write ABORTED to stable storage
    2. Abort transaction locally
    3. Release locks
    4. Send ACK to coordinator

Coordinator (after receiving all ACKs):
  1. Write COMPLETE to stable storage
  2. Forget transaction
```

### 1.3 Correctness Proof

```
Theorem 1 (Atomicity): All participants commit if and only if all
voted YES, and all abort if any voted NO.

Proof:
(⇒) If all participants commit:
  - Coordinator received all YES votes
  - Coordinator wrote COMMIT to stable storage
  - Therefore all participants voted YES

(⇐) If all participants voted YES:
  - Coordinator receives all YES votes
  - Coordinator writes COMMIT to stable storage
  - Coordinator sends COMMIT to all participants
  - All participants commit

For abort: If any participant votes NO:
  - Coordinator receives at least one NO
  - Coordinator writes ABORT to stable storage
  - Coordinator sends ABORT to all participants
  - All participants abort

Therefore, atomicity is guaranteed. ∎

Theorem 2 (Consistency): All participants reach the same decision.

Proof:
Participants decide based on coordinator's decision message.
Coordinator sends the same decision (COMMIT or ABORT) to all participants.
Therefore, all participants reach the same decision. ∎

Theorem 3 (Durability): Once a participant commits, the transaction
stays committed even after failures.

Proof:
Participants write COMMITTED to stable storage before committing.
On recovery, participants check stable storage:
  - If COMMITTED, transaction is already committed
  - If PREPARED, contact coordinator for decision
  - If ABORTED, transaction is already aborted

Therefore, committed transactions remain committed. ∎
```

### 1.4 Failure Scenarios and Recovery

**Coordinator Failure**:

```
Case 1: Before writing COMMIT/ABORT
  - Participants timeout waiting for decision
  - Participants can unilaterally abort (if not PREPARED)
  - If PREPARED, participants must block until coordinator recovers

Case 2: After writing COMMIT/ABORT but before sending to all
  - Coordinator recovers and re-sends decision
  - Participants may receive duplicate COMMIT/ABORT messages
  - Idempotent handling required

Case 3: After sending to all but before receiving all ACKs
  - Coordinator recovers and re-sends decision
  - Participants respond with ACK again
```

**Participant Failure**:

```
Case 1: Before responding to PREPARE
  - Coordinator times out
  - Coordinator aborts transaction
  - Other participants abort

Case 2: After voting YES (PREPARED state)
  - Participant recovers and checks stable storage
  - If PREPARED, contact coordinator for decision
  - Must commit if coordinator decided COMMIT

Case 3: After committing/aborting
  - Transaction complete
  - Ignore duplicate COMMIT/ABORT messages
```

## TLA+ Specification

```tla
----------------------------- MODULE TwoPhaseCommit -----------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Participants,      \* Set of participant IDs
          Coordinator,       \* Coordinator ID
          Values             \* Values being committed

VARIABLES cState,            \* Coordinator state
          pState,            \* Participant states
          cLog,              \* Coordinator log
          pLog,              \* Participant logs
          messages,          \* In-flight messages
          decisions          \* Final decisions

\* States
CoordinatorState == {"init", "preparing", "committed", "aborted", "completed"}
ParticipantState == {"init", "prepared", "committed", "aborted"}

\* Decision values
Decision == {"yes", "no", "commit", "abort"}

\* Message types
Message == [type: {"prepare", "yes", "no", "commit", "abort", "ack"},
            from: {Coordinator} ∪ Participants,
            to: {Coordinator} ∪ Participants,
            tx: Nat]

\* Initial state
Init ==
  ∧ cState = [state ↦ "init", tx ↦ 0]
  ∧ pState = [p ∈ Participants ↦ [state ↦ "init", vote ↦ ⊥]]
  ∧ cLog = ⟨⟩
  ∧ pLog = [p ∈ Participants ↦ ⟨⟩]
  ∧ messages = {}
  ∧ decisions = [p ∈ Participants ↦ ⊥]

\* Phase 1: Coordinator sends PREPARE
CoordinatorPrepare ==
  ∧ cState.state = "init"
  ∧ cState' = [cState EXCEPT !.state = "preparing", !.tx = @ + 1]
  ∧ cLog' = Append(cLog, [type ↦ "start-2pc", tx ↦ cState.tx + 1])
  ∧ messages' = messages ∪
       {[type ↦ "prepare", from ↦ Coordinator, to ↦ p, tx ↦ cState.tx + 1] :
        p ∈ Participants}
  ∧ UNCHANGED ⟨pState, pLog, decisions⟩

\* Participant receives PREPARE and votes
ParticipantPrepare(p) ==
  ∧ pState[p].state = "init"
  ∧ ∃m ∈ messages:
      ∧ m.type = "prepare"
      ∧ m.to = p
      ∧ m.tx = cState.tx
      ∧ IF CanCommit(p)  \* Abstract predicate
        THEN ∧ pState' = [pState EXCEPT ![p] = [state ↦ "prepared", vote ↦ "yes"]]
             ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "prepared", tx ↦ m.tx])]
             ∧ messages' = messages ∪
                  {[type ↦ "yes", from ↦ p, to ↦ Coordinator, tx ↦ m.tx]}
        ELSE ∧ pState' = [pState EXCEPT ![p] = [state ↦ "aborted", vote ↦ "no"]]
             ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "aborted", tx ↦ m.tx])]
             ∧ messages' = messages ∪
                  {[type ↦ "no", from ↦ p, to ↦ Coordinator, tx ↦ m.tx]}
  ∧ UNCHANGED ⟨cState, cLog, decisions⟩

\* Coordinator makes decision
CoordinatorDecide ==
  ∧ cState.state = "preparing"
  ∧ LET votes == {m ∈ messages :
                   m.type ∈ {"yes", "no"} ∧ m.to = Coordinator ∧ m.tx = cState.tx}
        allVoted == Cardinality({m.from : m ∈ votes}) = Cardinality(Participants)
        allYes == ∀m ∈ votes : m.type = "yes"
    IN ∧ allVoted
       ∧ IF allYes
         THEN ∧ cState' = [cState EXCEPT !.state = "committed"]
              ∧ cLog' = Append(cLog, [type ↦ "commit", tx ↦ cState.tx])
              ∧ messages' = messages ∪
                   {[type ↦ "commit", from ↦ Coordinator, to ↦ p, tx ↦ cState.tx] :
                    p ∈ Participants}
         ELSE ∧ cState' = [cState EXCEPT !.state = "aborted"]
              ∧ cLog' = Append(cLog, [type ↦ "abort", tx ↦ cState.tx])
              ∧ messages' = messages ∪
                   {[type ↦ "abort", from ↦ Coordinator, to ↦ p, tx ↦ cState.tx] :
                    p ∈ Participants}
  ∧ UNCHANGED ⟨pState, pLog, decisions⟩

\* Participant receives decision
ParticipantDecide(p) ==
  ∧ pState[p].state ∈ {"prepared", "init"}
  ∧ ∃m ∈ messages:
      ∧ m.type ∈ {"commit", "abort"}
      ∧ m.to = p
      ∧ m.tx = cState.tx
      ∧ IF m.type = "commit"
        THEN ∧ pState' = [pState EXCEPT ![p].state = "committed"]
             ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "committed", tx ↦ m.tx])]
             ∧ decisions' = [decisions EXCEPT ![p] = "commit"]
        ELSE ∧ pState' = [pState EXCEPT ![p].state = "aborted"]
             ∧ pLog' = [pLog EXCEPT ![p] = Append(@, [type ↦ "aborted", tx ↦ m.tx])]
             ∧ decisions' = [decisions EXCEPT ![p] = "abort"]
      ∧ messages' = messages ∪
           {[type ↦ "ack", from ↦ p, to ↦ Coordinator, tx ↦ m.tx]}
  ∧ UNCHANGED ⟨cState, cLog⟩

\* Coordinator completes transaction
CoordinatorComplete ==
  ∧ cState.state ∈ {"committed", "aborted"}
  ∧ LET acks == {m ∈ messages :
                  m.type = "ack" ∧ m.to = Coordinator ∧ m.tx = cState.tx}
    IN Cardinality({m.from : m ∈ acks}) = Cardinality(Participants)
  ∧ cState' = [cState EXCEPT !.state = "completed"]
  ∧ cLog' = Append(cLog, [type ↦ "complete", tx ↦ cState.tx])
  ∧ UNCHANGED ⟨pState, pLog, messages, decisions⟩

\* Next state
Next ==
  ∨ CoordinatorPrepare
  ∨ ∃p ∈ Participants: ParticipantPrepare(p)
  ∨ CoordinatorDecide
  ∨ ∃p ∈ Participants: ParticipantDecide(p)
  ∨ CoordinatorComplete

\* Invariants

\* All participants reach same decision
Atomicity ==
  ∀p1, p2 ∈ Participants:
    decisions[p1] ≠ ⊥ ∧ decisions[p2] ≠ ⊥ ⇒ decisions[p1] = decisions[p2]

\* Coordinator decides commit only if all voted yes
CommitConsistency ==
  cState.state = "committed"
    ⇒ ∀p ∈ Participants: pState[p].vote = "yes"

\* No participant commits if coordinator aborted
AbortConsistency ==
  cState.state = "aborted"
    ⇒ ∀p ∈ Participants: decisions[p] ≠ "commit"

=============================================================================
```

## Algorithm Pseudocode

### Classic 2PC

```
Algorithm: Classic Two-Phase Commit

Types:
  TransactionID: unique identifier
  Vote: YES | NO
  Decision: COMMIT | ABORT

Coordinator:
  State:
    txId: TransactionID
    participants: Set<Participant>
    votes: Map<Participant, Vote>
    state: INIT | PREPARING | COMMITTED | ABORTED

  Procedure BeginTransaction(txId, participants):
    self.txId = txId
    self.participants = participants
    self.votes = empty map
    self.state = INIT

    WriteLog(START, txId)
    self.state = PREPARING

    for each p in participants:
      Send(PREPARE, txId) to p

    StartTimeout(TIMEOUT_PREPARE)

  On Receive(VOTE, txId, participant, vote):
    if state ≠ PREPARING:
      return  // Ignore late votes

    votes[participant] = vote

    if all participants have voted:
      CancelTimeout()

      if all votes are YES:
        WriteLog(COMMIT, txId)
        state = COMMITTED
        for each p in participants:
          Send(COMMIT, txId) to p
      else:
        WriteLog(ABORT, txId)
        state = ABORTED
        for each p in participants:
          Send(ABORT, txId) to p

      StartTimeout(TIMEOUT_COMPLETE)

  On Receive(ACK, txId, participant):
    Mark participant as acknowledged

    if all participants acknowledged:
      CancelTimeout()
      WriteLog(COMPLETE, txId)
      state = COMPLETED
      Cleanup(txId)

  On Timeout(timeoutType):
    if timeoutType = TIMEOUT_PREPARE:
      // Some participants didn't vote
      WriteLog(ABORT, txId)
      state = ABORTED
      for each p in participants who haven't voted:
        Send(ABORT, txId) to p

    else if timeoutType = TIMEOUT_COMPLETE:
      // Some participants didn't acknowledge
      // Re-send decision
      decision = (state = COMMITTED) ? COMMIT : ABORT
      for each p in participants who haven't acknowledged:
        Send(decision, txId) to p
      RestartTimeout(TIMEOUT_COMPLETE)

Participant:
  State:
    state: INIT | PREPARED | COMMITTED | ABORTED
    txId: TransactionID

  On Receive(PREPARE, txId):
    if not CanCommit(txId):
      WriteLog(ABORT, txId)
      state = ABORTED
      Send(VOTE, txId, self, NO) to coordinator
      return

    // Acquire locks, check constraints
    AcquireLocks(txId)

    if LockAcquisitionFails:
      ReleaseLocks(txId)
      WriteLog(ABORT, txId)
      state = ABORTED
      Send(VOTE, txId, self, NO) to coordinator
      return

    WriteLog(PREPARED, txId)
    state = PREPARED
    Send(VOTE, txId, self, YES) to coordinator

    StartTimeout(TIMEOUT_DECISION)

  On Receive(COMMIT, txId):
    CancelTimeout()
    WriteLog(COMMITTED, txId)
    state = COMMITTED
    Commit(txId)
    ReleaseLocks(txId)
    Send(ACK, txId, self) to coordinator

  On Receive(ABORT, txId):
    CancelTimeout()
    WriteLog(ABORTED, txId)
    state = ABORTED
    Abort(txId)
    ReleaseLocks(txId)
    Send(ACK, txId, self) to coordinator

  On Timeout(TIMEOUT_DECISION):
    // Coordinator may have failed
    // Contact coordinator or other participants
    ContactCoordinatorForDecision(txId)

Recovery Procedure (for both):
  logs = ReadStableLog()

  for each incomplete tx in logs:
    if I'm coordinator:
      if log contains COMMIT:
        ReSend(COMMIT, tx) to all participants
      else if log contains ABORT or START:
        ReSend(ABORT, tx) to all participants

    else (I'm participant):
      if log contains COMMITTED:
        // Already committed, nothing to do
      else if log contains ABORTED:
        // Already aborted, nothing to do
      else if log contains PREPARED:
        // Need to find out decision
        ContactCoordinatorForDecision(tx)
```

### Presumed Abort Optimization

```
Algorithm: 2PC with Presumed Abort Optimization

Optimization: If coordinator crashes before writing COMMIT,
assume ABORT on recovery (no need to contact participants).

Coordinator Changes:
  - Don't write ABORT to log
  - On recovery, if no COMMIT record exists for a transaction,
    assume it was aborted
  - Participants can clean up after timeout without coordinator ACK

Participant Changes:
  - If no decision received after timeout, assume ABORT
  - Can unilaterally abort if in INIT state

Benefits:
  - Fewer log writes for aborts
  - Faster recovery
  - Less blocking

Coordinator Recovery:
  logs = ReadStableLog()

  for each tx with START but no COMPLETE:
    if log contains COMMIT:
      state = COMMITTED
      ReSend(COMMIT, tx)
    else:
      // No COMMIT record - presumed abort
      state = ABORTED
      // No need to notify participants
      WriteLog(COMPLETE, tx)
```

## Go Implementation

```go
// Package twopc implements the Two-Phase Commit protocol
package twopc

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// TransactionState represents the state of a transaction
type TransactionState int

const (
 StateInit TransactionState = iota
 StatePreparing
 StatePrepared
 StateCommitted
 StateAborted
 StateCompleted
)

func (s TransactionState) String() string {
 switch s {
 case StateInit:
  return "INIT"
 case StatePreparing:
  return "PREPARING"
 case StatePrepared:
  return "PREPARED"
 case StateCommitted:
  return "COMMITTED"
 case StateAborted:
  return "ABORTED"
 case StateCompleted:
  return "COMPLETED"
 default:
  return "UNKNOWN"
 }
}

// Vote represents a participant's vote
type Vote int

const (
 VoteYes Vote = iota
 VoteNo
)

// Decision represents the coordinator's decision
type Decision int

const (
 DecisionCommit Decision = iota
 DecisionAbort
)

// LogEntry represents an entry in the stable log
type LogEntry struct {
 TxID      string
 Type      string // START, PREPARED, COMMIT, ABORT, COMPLETE
 Timestamp time.Time
 Data      interface{}
}

// StableLog represents the write-ahead log
type StableLog interface {
 Append(entry LogEntry) error
 Read(txID string) ([]LogEntry, error)
}

// Message types
type Message struct {
 Type   string // PREPARE, VOTE, COMMIT, ABORT, ACK
 TxID   string
 From   string
 To     string
 Vote   Vote
 Decision Decision
}

// Network represents the communication layer
type Network interface {
 Send(msg Message) error
 Receive(timeout time.Duration) (Message, bool)
}

// Participant defines the interface for transaction participants
type Participant interface {
 ID() string
 Prepare(txID string) error
 Commit(txID string) error
 Abort(txID string) error
 CanCommit(txID string) bool
}

// Coordinator manages the 2PC protocol
type Coordinator struct {
 mu           sync.RWMutex
 id           string
 network      Network
 log          StableLog
 transactions map[string]*CoordinatorTx
}

// CoordinatorTx tracks a transaction
type CoordinatorTx struct {
 TxID         string
 Participants []string
 Votes        map[string]Vote
 State        TransactionState
 mu           sync.Mutex
 timeoutChan  chan struct{}
}

// NewCoordinator creates a new coordinator
func NewCoordinator(id string, network Network, log StableLog) *Coordinator {
 return &Coordinator{
  id:           id,
  network:      network,
  log:          log,
  transactions: make(map[string]*CoordinatorTx),
 }
}

// BeginTransaction starts a new 2PC transaction
func (c *Coordinator) BeginTransaction(txID string, participants []string) error {
 c.mu.Lock()
 defer c.mu.Unlock()

 // Check if transaction already exists
 if _, exists := c.transactions[txID]; exists {
  return fmt.Errorf("transaction %s already exists", txID)
 }

 // Write start log
 if err := c.log.Append(LogEntry{
  TxID:      txID,
  Type:      "START-2PC",
  Timestamp: time.Now(),
 }); err != nil {
  return fmt.Errorf("failed to write start log: %w", err)
 }

 // Create transaction
 tx := &CoordinatorTx{
  TxID:         txID,
  Participants: participants,
  Votes:        make(map[string]Vote),
  State:        StatePreparing,
  timeoutChan:  make(chan struct{}),
 }
 c.transactions[txID] = tx

 // Send PREPARE to all participants
 for _, p := range participants {
  msg := Message{
   Type: "PREPARE",
   TxID: txID,
   From: c.id,
   To:   p,
  }
  if err := c.network.Send(msg); err != nil {
   return fmt.Errorf("failed to send PREPARE to %s: %w", p, err)
  }
 }

 // Start timeout
 go c.waitForVotes(txID)

 return nil
}

func (c *Coordinator) waitForVotes(txID string) {
 tx := c.transactions[txID]
 if tx == nil {
  return
 }

 timeout := time.NewTimer(30 * time.Second)
 defer timeout.Stop()

 select {
 case <-timeout.C:
  c.handleTimeout(txID)
 case <-tx.timeoutChan:
  // All votes received
 }
}

func (c *Coordinator) handleTimeout(txID string) {
 c.mu.Lock()
 tx, exists := c.transactions[txID]
 c.mu.Unlock()

 if !exists || tx.State != StatePreparing {
  return
 }

 // Timeout - abort transaction
 tx.mu.Lock()
 tx.State = StateAborted
 tx.mu.Unlock()

 // Write abort log
 c.log.Append(LogEntry{
  TxID:      txID,
  Type:      "ABORT",
  Timestamp: time.Now(),
 })

 // Notify all participants
 for _, p := range tx.Participants {
  if _, voted := tx.Votes[p]; !voted {
   c.network.Send(Message{
    Type: "ABORT",
    TxID: txID,
    From: c.id,
    To:   p,
   })
  }
 }
}

// HandleVote processes a vote from a participant
func (c *Coordinator) HandleVote(msg Message) error {
 c.mu.Lock()
 tx, exists := c.transactions[msg.TxID]
 c.mu.Unlock()

 if !exists {
  return fmt.Errorf("unknown transaction: %s", msg.TxID)
 }

 tx.mu.Lock()
 defer tx.mu.Unlock()

 if tx.State != StatePreparing {
  return nil // Ignore late votes
 }

 // Record vote
 tx.Votes[msg.From] = msg.Vote

 // Check if all participants have voted
 if len(tx.Votes) == len(tx.Participants) {
  // Cancel timeout
  close(tx.timeoutChan)

  // Count votes
  allYes := true
  for _, vote := range tx.Votes {
   if vote == VoteNo {
    allYes = false
    break
   }
  }

  // Make decision
  if allYes {
   return c.commitTransaction(tx)
  }
  return c.abortTransaction(tx)
 }

 return nil
}

func (c *Coordinator) commitTransaction(tx *CoordinatorTx) error {
 tx.State = StateCommitted

 // Write commit log
 if err := c.log.Append(LogEntry{
  TxID:      tx.TxID,
  Type:      "COMMIT",
  Timestamp: time.Now(),
 }); err != nil {
  return err
 }

 // Send COMMIT to all participants
 for _, p := range tx.Participants {
  c.network.Send(Message{
   Type:     "COMMIT",
   TxID:     tx.TxID,
   From:     c.id,
   To:       p,
   Decision: DecisionCommit,
  })
 }

 return nil
}

func (c *Coordinator) abortTransaction(tx *CoordinatorTx) error {
 tx.State = StateAborted

 // Write abort log
 if err := c.log.Append(LogEntry{
  TxID:      tx.TxID,
  Type:      "ABORT",
  Timestamp: time.Now(),
 }); err != nil {
  return err
 }

 // Send ABORT to all participants
 for _, p := range tx.Participants {
  c.network.Send(Message{
   Type:     "ABORT",
   TxID:     tx.TxID,
   From:     c.id,
   To:       p,
   Decision: DecisionAbort,
  })
 }

 return nil
}

// ParticipantNode represents a 2PC participant
type ParticipantNode struct {
 mu            sync.RWMutex
 id            string
 network       Network
 log           StableLog
 transactions  map[string]*ParticipantTx
 preparedFunc  func(string) bool
 commitFunc    func(string) error
 abortFunc     func(string) error
}

// ParticipantTx tracks a participant's transaction state
type ParticipantTx struct {
 TxID   string
 State  TransactionState
 mu     sync.Mutex
}

// NewParticipant creates a new participant
func NewParticipant(id string, network Network, log StableLog) *ParticipantNode {
 return &ParticipantNode{
  id:           id,
  network:      network,
  log:          log,
  transactions: make(map[string]*ParticipantTx),
 }
}

// SetCallbacks sets the transaction callbacks
func (p *ParticipantNode) SetCallbacks(
 canCommit func(string) bool,
 commit func(string) error,
 abort func(string) error,
) {
 p.preparedFunc = canCommit
 p.commitFunc = commit
 p.abortFunc = abort
}

// HandlePrepare handles a PREPARE message
func (p *ParticipantNode) HandlePrepare(msg Message) error {
 p.mu.Lock()
 defer p.mu.Unlock()

 // Check if already handling this transaction
 if tx, exists := p.transactions[msg.TxID]; exists {
  if tx.State == StatePrepared {
   // Already prepared, resend YES
   p.network.Send(Message{
    Type: "VOTE",
    TxID: msg.TxID,
    From: p.id,
    To:   msg.From,
    Vote: VoteYes,
   })
  }
  return nil
 }

 // Check if we can commit
 if p.preparedFunc == nil || !p.preparedFunc(msg.TxID) {
  // Cannot commit - vote NO
  p.network.Send(Message{
   Type: "VOTE",
   TxID: msg.TxID,
   From: p.id,
   To:   msg.From,
   Vote: VoteNo,
  })

  // Write abort log
  p.log.Append(LogEntry{
   TxID:      msg.TxID,
   Type:      "ABORTED",
   Timestamp: time.Now(),
  })

  return nil
 }

 // Can commit - enter PREPARED state
 tx := &ParticipantTx{
  TxID:  msg.TxID,
  State: StatePrepared,
 }
 p.transactions[msg.TxID] = tx

 // Write prepared log
 if err := p.log.Append(LogEntry{
  TxID:      msg.TxID,
  Type:      "PREPARED",
  Timestamp: time.Now(),
 }); err != nil {
  return err
 }

 // Send YES vote
 return p.network.Send(Message{
  Type: "VOTE",
  TxID: msg.TxID,
  From: p.id,
  To:   msg.From,
  Vote: VoteYes,
 })
}

// HandleDecision handles COMMIT or ABORT from coordinator
func (p *ParticipantNode) HandleDecision(msg Message) error {
 p.mu.Lock()
 tx, exists := p.transactions[msg.TxID]
 p.mu.Unlock()

 if !exists && msg.Decision == DecisionAbort {
  // Abort for unknown transaction - ok to ignore
  return nil
 }

 if !exists {
  return fmt.Errorf("no transaction %s to commit", msg.TxID)
 }

 tx.mu.Lock()
 defer tx.mu.Unlock()

 if msg.Decision == DecisionCommit {
  tx.State = StateCommitted

  // Write commit log
  p.log.Append(LogEntry{
   TxID:      msg.TxID,
   Type:      "COMMITTED",
   Timestamp: time.Now(),
  })

  // Commit locally
  if p.commitFunc != nil {
   p.commitFunc(msg.TxID)
  }
 } else {
  tx.State = StateAborted

  // Write abort log
  p.log.Append(LogEntry{
   TxID:      msg.TxID,
   Type:      "ABORTED",
   Timestamp: time.Now(),
  })

  // Abort locally
  if p.abortFunc != nil {
   p.abortFunc(msg.TxID)
  }
 }

 // Send ACK
 return p.network.Send(Message{
  Type: "ACK",
  TxID: msg.TxID,
  From: p.id,
  To:   msg.From,
 })
}

// Recover recovers from crash
func (p *ParticipantNode) Recover() error {
 // Read all log entries
 // For each PREPARED transaction without COMMITTED/ABORTED, contact coordinator
 return nil
}
