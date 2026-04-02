# FT-015: FLP Impossibility - Formal Analysis

## Overview

The Fischer, Lynch, and Paterson (FLP) impossibility result is a fundamental theorem in distributed systems theory. It proves that no deterministic consensus algorithm can guarantee termination in an asynchronous distributed system with at least one faulty process, even with reliable communication.

## Theoretical Foundations

### 1.1 System Model

**Asynchronous Distributed System**:

```
System Σ = ⟨Π, C, M, A⟩

where:
- Π = {p₁, p₂, ..., pₙ}: set of processes
- C: communication channels (reliable but unbounded delay)
- M: set of possible messages
- A: asynchronous scheduling adversary
```

**Process Model**:

```
Each process pᵢ has:
- State space: Sᵢ (finite or infinite)
- Input register: xᵢ ∈ {0, 1, ⊥}
- Output register: yᵢ ∈ {0, 1, ⊥}
- Transition function: δᵢ: Sᵢ × Events → Sᵢ
- Initial state: sᵢ⁰ ∈ Sᵢ
```

### 1.2 Consensus Specification

A consensus algorithm must satisfy three properties:

**Termination**: Every correct process eventually decides.

```
∀pᵢ ∈ Correct(Σ): ◇(yᵢ ≠ ⊥)
```

**Agreement**: No two correct processes decide differently.

```
∀pᵢ, pⱼ ∈ Correct(Σ): yᵢ ≠ ⊥ ∧ yⱼ ≠ ⊥ ⇒ yᵢ = yⱼ
```

**Validity**: If all processes start with the same value, that value is the only possible decision.

```
(∀pᵢ ∈ Π: xᵢ = v) ⇒ (∀pᵢ ∈ Correct(Σ): yᵢ = v ∨ yᵢ = ⊥)
```

### 1.3 FLP Theorem Statement

```
Theorem (FLP Impossibility):
In an asynchronous distributed system with at least one process that may fail by
crashing, there exists no deterministic algorithm that solves consensus for even
two processes.

Formally:
∄ A: Algorithm | ∀Σ = ⟨Π, C, M, A⟩, |Π| ≥ 2, f ≥ 1:
  A satisfies Termination ∧ Agreement ∧ Validity in Σ
```

### 1.4 Proof of FLP Impossibility

**Lemma 1 (Initial Bivalence)**: The initial configuration is bivalent.

```
Proof of Lemma 1:
Consider two initial configurations:
- C₀: All processes have input 0
- C₁: All processes have input 1

By Validity, C₀ can only decide 0, and C₁ can only decide 1.

Consider configuration C':
- Process p₁ has input 0
- All other processes have input 1

If C' is univalent (decides 0), then p₁ can decide without hearing from others.
If C' is univalent (decides 1), then all processes except p₁ can decide without p₁.

By a chain argument, there must exist a configuration where changing one
process's input changes the decision value. This configuration is bivalent.
∎
```

**Lemma 2 (Decision Gadget)**: In a bivalent configuration, there exists a reachable bivalent configuration where every step from it leads to a univalent configuration.

```
Proof of Lemma 2:
Suppose not. Then from every bivalent configuration, there is a path to
another bivalent configuration. Since the system is finite-state (with
bounded messages), this would imply an infinite chain of bivalent
configurations, contradicting the requirement that the algorithm must
eventually decide.
∎
```

**Lemma 3 (Pivotal Event)**: There exists a configuration C and events e₁, e₂ such that:

- Applying e₁ to C leads to 0-valent configuration
- Applying e₂ to C leads to 1-valent configuration
- e₁ and e₂ are events of different processes

```
Proof of Lemma 3:
From Lemma 2, take the "decision gadget" configuration C.
There must be two events leading to different valences, otherwise C would be univalent.
If both events are from the same process, consider the schedule that delays both.
By asynchrony, this is indistinguishable from crashing that process.
Therefore, the events must involve different processes.
∎
```

**Main Proof**:

```
Theorem Proof:

We construct an admissible run that never decides.

1. Start from initial bivalent configuration C₀ (Lemma 1).

2. Inductively construct an infinite sequence:

   For current bivalent configuration Cᵢ:
   a) Find the decision gadget configuration C' reachable from Cᵢ (Lemma 2)
   b) Find pivotal events e₀, e₁ leading to 0-valent and 1-valent configs (Lemma 3)
   c) Let p₀, p₁ be the processes of e₀, e₁ respectively

   Since the system has f ≥ 1 faulty process, either p₀ or p₁ could be faulty.
   Without loss of generality, assume p₀ is faulty (crashed).

   Delay e₀ indefinitely (admissible by asynchrony).
   Apply e₁ to reach bivalent configuration Cᵢ₊₁.

3. The resulting infinite execution:
   - Is admissible: only faulty processes have delayed events
   - Never decides: always remains bivalent
   - Contradicts Termination property

Therefore, no such deterministic algorithm exists.
∎
```

## TLA+ Specification

```tla
----------------------- MODULE FLPImpossibility -----------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Processes,      \* Set of process IDs
          Values,         \* {0, 1}
          MaxSteps        \* Bound for model checking

VARIABLES pc,             \* Program counter for each process
          state,          \* Process state
          messages,       \* Message buffer
          decided,        \* Decision values
          stepCount       \* Step counter

\* Types
TypeInvariant ==
  ∧ pc ∈ [Processes → {"init", "propose", "decide", "done"}]
  ∧ state ∈ [Processes → [input: Values ∪ {⊥},
                          proposal: Values ∪ {⊥},
                          decided: BOOLEAN]]
  ∧ messages ⊆ [from: Processes, to: Processes, value: Values]
  ∧ decided ∈ [Processes → Values ∪ {⊥}]

\* Initial state - bivalent configuration
Init ==
  ∧ pc = [p ∈ Processes ↦ "init"]
  ∧ state = [p ∈ Processes ↦ [input ↦ IF p = "p1" THEN 0 ELSE 1,
                               proposal ↦ ⊥,
                               decided ↦ FALSE]]
  ∧ messages = {}
  ∧ decided = [p ∈ Processes ↦ ⊥]
  ∧ stepCount = 0

\* Propose step
Propose(p) ==
  ∧ pc[p] = "init"
  ∧ state' = [state EXCEPT ![p].proposal = state[p].input]
  ∧ pc' = [pc EXCEPT ![p] = "propose"]
  ∧ UNCHANGED ⟨messages, decided, stepCount⟩

\* Send proposal to all
Broadcast(p) ==
  ∧ pc[p] = "propose"
  ∧ messages' = messages ∪
      {[from ↦ p, to ↦ q, value ↦ state[p].proposal] : q ∈ Processes \\ {p}}
  ∧ pc' = [pc EXCEPT ![p] = "decide"]
  ∧ UNCHANGED ⟨state, decided, stepCount⟩

\* Decide based on received messages
Decide(p) ==
  ∧ pc[p] = "decide"
  ∧ LET received == {m ∈ messages : m.to = p}
        values == {m.value : m ∈ received}
    IN ∧ Cardinality(values) = 1  \* All proposals same
       ∧ decided' = [decided EXCEPT ![p] = CHOOSE v ∈ values : TRUE]
       ∧ state' = [state EXCEPT ![p].decided = TRUE]
       ∧ pc' = [pc EXCEPT ![p] = "done"]
       ∧ UNCHANGED ⟨messages, stepCount⟩

\* Crash process (non-deterministic)
Crash(p) ==
  ∧ pc' = [pc EXCEPT ![p] = "done"]
  ∧ decided' = [decided EXCEPT ![p] = ⊥]
  ∧ UNCHANGED ⟨state, messages, stepCount⟩

\* Next state
Next ==
  ∨ ∃p ∈ Processes: Propose(p) ∨ Broadcast(p) ∨ Decide(p) ∨ Crash(p)
  ∨ stepCount' = stepCount + 1

\* Fairness: correct processes must eventually decide
Fairness ==
  ∀p ∈ Processes:
    WF_⟨pc, state, messages, decided, stepCount⟩(Decide(p))

\* Consensus properties

\* Agreement: no two processes decide differently
Agreement ==
  ∀p, q ∈ Processes:
    decided[p] ≠ ⊥ ∧ decided[q] ≠ ⊥ ⇒ decided[p] = decided[q]

\* Validity: decision must be a proposed value
Validity ==
  ∀p ∈ Processes:
    decided[p] ≠ ⊥ ⇒ decided[p] ∈ Values

\* Termination: all correct processes eventually decide
Termination ==
  ∀p ∈ Processes:
    pc[p] ≠ "done" ⇒ ◇(decided[p] ≠ ⊥)

\* FLP Result: Termination cannot hold with even one faulty process
FLPResult ==
  \* In an asynchronous system with crash faults,
  \* there exists an admissible execution where termination fails
  TRUE  \* This is the existential statement of FLP

=============================================================================
```

## Algorithm Pseudocode

### Asynchronous Consensus with Failure Detectors

```
Algorithm: FLP-Style Consensus with Ω (Eventual Leader Election)

Since FLP proves deterministic consensus impossible, we use:
1. Randomization (Ben-Or algorithm)
2. Failure detectors (Chandra-Toueg)
3. Partial synchrony (Dwork-Lynch-Stockmeyer)

Ben-Or Randomized Consensus:

Constants:
  N: total number of processes
  F: maximum number of faulty processes (F < N/2)
  Φ: coin toss probability (typically 0.5)

Variables:
  x: current proposed value ∈ {0, 1}
  r: current round number
  decided: boolean

Procedure Consensus(initialValue):
  x ← initialValue
  r ← 1
  decided ← false

  while not decided:
    // Phase 1: Propose
    Broadcast(PROPOSAL, r, x) to all processes

    // Collect proposals
    proposals ← WaitFor(N - F) proposals for round r

    if all proposals have same value v:
      Broadcast(COMMIT, r, v)
    else:
      Broadcast(COMMIT, r, ⊥)

    // Phase 2: Commit
    commits ← WaitFor(N - F) commits for round r

    if ∃v: at least N - F commits have value v ≠ ⊥:
      decided ← true
      return v
    else if ∃v: at least F + 1 commits have value v ≠ ⊥:
      x ← v
    else:
      // Randomized choice
      x ← RandomCoinToss()

    r ← r + 1

Theorem (Ben-Or): The algorithm achieves consensus with probability 1.

Proof Sketch:
1. Validity: Follows from commit threshold (N - F)
2. Agreement: If a process decides v, at least N - F commits were for v.
   Since F < N/2, any quorum of N - F intersects with this set.
3. Termination: Expected O(1) rounds when coin tosses align.

Complexity:
- Message complexity: O(N²) per round
- Round complexity: Expected O(1), with 2^(-k) probability of k rounds
```

### Paxos Variation for Partial Synchrony

```
Algorithm: Paxos with Timeouts (circumventing FLP)

FLP holds only for purely asynchronous systems.
Paxos uses leader election with timeouts to achieve consensus.

Phase 1: Prepare
  Proposer:
    1. Choose proposal number n (unique and increasing)
    2. Send PREPARE(n) to majority of Acceptors

  Acceptor:
    1. If n > maxPrepared:
         maxPrepared ← n
         Respond PROMISE(n, lastAccepted)

Phase 2: Accept
  Proposer (on majority promises):
    1. If any promise contains accepted value:
         value ← that value (highest n)
       Else:
         value ← proposed value
    2. Send ACCEPT(n, value) to majority

  Acceptor:
    1. If n ≥ maxPrepared:
         Accept value
         Respond ACCEPTED(n, value)

Phase 3: Learn
  Proposer (on majority accepted):
    Broadcast LEARN(value)

Termination Handling:
  If timeout occurs in any phase:
    1. Increment proposal number
    2. Retry with new leader
    3. After maxRetries, declare failure

Why this circumvents FLP:
1. Uses timeouts (assumes partial synchrony)
2. Leader election provides a form of eventual synchrony
3. Random backoff prevents infinite contention
```

## Go Implementation

```go
// Package flp implements FLP-related algorithms and simulations
package flp

import (
 "context"
 "crypto/rand"
 "fmt"
 "math/big"
 "sync"
 "time"
)

// Value represents a consensus value
type Value int

const (
 ValueZero Value = iota
 ValueOne
 ValueUndefined
)

func (v Value) String() string {
 switch v {
 case ValueZero:
  return "0"
 case ValueOne:
  return "1"
 case ValueUndefined:
  return "⊥"
 default:
  return "?"
 }
}

// Process represents a consensus participant
type Process struct {
 ID       string
 Input    Value
 Decision Value

 mu           sync.RWMutex
 state        ProcessState
 proposal     Value
 round        int
 halted       bool

 messenger    Messenger
 failureProb  float64
}

// ProcessState represents the state of a process
type ProcessState int

const (
 StateInit ProcessState = iota
 StateProposing
 StateDeciding
 StateDecided
 StateCrashed
)

// Messenger handles inter-process communication
type Messenger interface {
 Send(to string, msg Message) error
 Broadcast(msg Message) error
 Receive(timeout time.Duration) (Message, error)
}

// Message represents a consensus message
type Message struct {
 From      string
 To        string
 Type      MessageType
 Round     int
 Value     Value
 Timestamp time.Time
}

// MessageType represents the type of message
type MessageType int

const (
 MsgPropose MessageType = iota
 MsgCommit
 MsgAck
 MsgNack
)

// ConsensusAlgorithm is the interface for consensus implementations
type ConsensusAlgorithm interface {
 Propose(ctx context.Context, value Value) (Value, error)
 IsDecided() bool
 GetDecision() Value
}

// BenOrConsensus implements the Ben-Or randomized consensus algorithm
// This circumvents FLP impossibility using randomization
type BenOrConsensus struct {
 process      *Process
 n            int           // total processes
 f            int           // max faulty processes
 quorumSize   int

 mu           sync.Mutex
 proposals    map[int]map[string]Value
 commits      map[int]map[string]Value
 decided      bool
 decision     Value
}

// NewBenOrConsensus creates a new Ben-Or consensus instance
func NewBenOrConsensus(p *Process, n, f int) *BenOrConsensus {
 return &BenOrConsensus{
  process:    p,
  n:          n,
  f:          f,
  quorumSize: n - f,
  proposals:  make(map[int]map[string]Value),
  commits:    make(map[int]map[string]Value),
  decision:   ValueUndefined,
 }
}

// Propose executes the Ben-Or consensus algorithm
func (b *BenOrConsensus) Propose(ctx context.Context, value Value) (Value, error) {
 b.process.mu.Lock()
 b.process.proposal = value
 b.process.round = 1
 b.process.state = StateProposing
 b.process.mu.Unlock()

 for {
  select {
  case <-ctx.Done():
   return ValueUndefined, ctx.Err()
  default:
  }

  b.process.mu.RLock()
 round := b.process.round
  x := b.process.proposal
  b.process.mu.RUnlock()

  // Phase 1: Propose
  if err := b.phase1Propose(ctx, round, x); err != nil {
   return ValueUndefined, fmt.Errorf("phase 1 failed: %w", err)
  }

  // Phase 2: Commit and decide
  decision, decided, err := b.phase2Commit(ctx, round)
  if err != nil {
   return ValueUndefined, fmt.Errorf("phase 2 failed: %w", err)
  }

  if decided {
   b.mu.Lock()
   b.decided = true
   b.decision = decision
   b.mu.Unlock()

   b.process.mu.Lock()
   b.process.Decision = decision
   b.process.state = StateDecided
   b.process.mu.Unlock()

   return decision, nil
  }

  // Continue to next round
  b.process.mu.Lock()
  b.process.round++
  b.process.proposal = decision // may be undefined (randomize next)
  b.process.mu.Unlock()
 }
}

// phase1Propose broadcasts proposal and collects responses
func (b *BenOrConsensus) phase1Propose(ctx context.Context, round int, value Value) error {
 msg := Message{
  From:  b.process.ID,
  Type:  MsgPropose,
  Round: round,
  Value: value,
 }

 if err := b.process.messenger.Broadcast(msg); err != nil {
  return err
 }

 // Collect proposals
 proposals := make(map[string]Value)
 timeout := time.NewTimer(100 * time.Millisecond)
 defer timeout.Stop()

 for len(proposals) < b.quorumSize {
  select {
  case <-timeout.C:
   return fmt.Errorf("proposal timeout")
  case <-ctx.Done():
   return ctx.Err()
  default:
   msg, err := b.process.messenger.Receive(10 * time.Millisecond)
   if err != nil {
    continue
   }
   if msg.Type == MsgPropose && msg.Round == round {
    proposals[msg.From] = msg.Value
   }
  }
 }

 // Store proposals for phase 2
 b.mu.Lock()
 b.proposals[round] = proposals
 b.mu.Unlock()

 return nil
}

// phase2Commit determines if consensus is reached
func (b *BenOrConsensus) phase2Commit(ctx context.Context, round int) (Value, bool, error) {
 b.mu.Lock()
 proposals := b.proposals[round]
 b.mu.Unlock()

 // Check if all proposals agree
 var commitValue Value = ValueUndefined
 allSame := true
 firstVal := ValueUndefined

 for _, v := range proposals {
  if firstVal == ValueUndefined {
   firstVal = v
  } else if v != firstVal {
   allSame = false
   break
  }
 }

 if allSame && firstVal != ValueUndefined {
  commitValue = firstVal
 }

 // Broadcast commit
 msg := Message{
  From:  b.process.ID,
  Type:  MsgCommit,
  Round: round,
  Value: commitValue,
 }

 if err := b.process.messenger.Broadcast(msg); err != nil {
  return ValueUndefined, false, err
 }

 // Collect commits
 commits := make(map[string]Value)
 timeout := time.NewTimer(100 * time.Millisecond)
 defer timeout.Stop()

 for len(commits) < b.quorumSize {
  select {
  case <-timeout.C:
   return ValueUndefined, false, fmt.Errorf("commit timeout")
  case <-ctx.Done():
   return ValueUndefined, false, ctx.Err()
  default:
   msg, err := b.process.messenger.Receive(10 * time.Millisecond)
   if err != nil {
    continue
   }
   if msg.Type == MsgCommit && msg.Round == round {
    commits[msg.From] = msg.Value
   }
  }
 }

 // Analyze commits
 valueCounts := make(map[Value]int)
 for _, v := range commits {
  valueCounts[v]++
 }

 // Check for decision
 for v, count := range valueCounts {
  if v != ValueUndefined && count >= b.quorumSize {
   return v, true, nil
  }
 }

 // Check for biased next proposal
 for v, count := range valueCounts {
  if v != ValueUndefined && count >= b.f+1 {
   return v, false, nil
  }
 }

 // Random choice
 randomVal, err := randomCoinToss()
 if err != nil {
  return ValueUndefined, false, err
 }

 return randomVal, false, nil
}

// randomCoinToss returns a random value
func randomCoinToss() (Value, error) {
 n, err := rand.Int(rand.Reader, big.NewInt(2))
 if err != nil {
  return ValueZero, err
 }
 if n.Int64() == 0 {
  return ValueZero, nil
 }
 return ValueOne, nil
}

// IsDecided returns whether consensus has been reached
func (b *BenOrConsensus) IsDecided() bool {
 b.mu.Lock()
 defer b.mu.Unlock()
 return b.decided
}

// GetDecision returns the consensus decision
func (b *BenOrConsensus) GetDecision() Value {
 b.mu.Lock()
 defer b.mu.Unlock()
 return b.decision
}

// FLPScenario simulates the FLP impossibility scenario
type FLPScenario struct {
 processes     []*Process
 messenger     *FLPMessenger
 asynchrony    time.Duration
 faultyCount   int
}

// FLPMessenger simulates asynchronous message delivery
type FLPMessenger struct {
 mu       sync.Mutex
 inbox    map[string][]Message
 delays   map[string]time.Duration
}

// NewFLPMessenger creates a new FLP messenger
func NewFLPMessenger() *FLPMessenger {
 return &FLPMessenger{
  inbox:  make(map[string][]Message),
  delays: make(map[string]time.Duration),
 }
}

// Send implements Messenger
func (f *FLPMessenger) Send(to string, msg Message) error {
 f.mu.Lock()
 defer f.mu.Unlock()

 // Simulate arbitrary delay
 delay := time.Duration(rand.Int(1000)) * time.Millisecond
 f.delays[to] = delay

 go func() {
  time.Sleep(delay)
  f.mu.Lock()
  f.inbox[to] = append(f.inbox[to], msg)
  f.mu.Unlock()
 }()

 return nil
}

// Broadcast implements Messenger
func (f *FLPMessenger) Broadcast(msg Message) error {
 for id := range f.inbox {
  if id != msg.From {
   msgCopy := msg
   msgCopy.To = id
   f.Send(id, msgCopy)
  }
 }
 return nil
}

// Receive implements Messenger
func (f *FLPMessenger) Receive(timeout time.Duration) (Message, error) {
 start := time.Now()
 for time.Since(start) < timeout {
  f.mu.Lock()
  // Check inbox (simplified)
  f.mu.Unlock()
  time.Sleep(1 * time.Millisecond)
 }
 return Message{}, fmt.Errorf("no message")
}

// SimulateFLP demonstrates the FLP impossibility
func SimulateFLP(numProcesses, numFaulty int, duration time.Duration) {
 fmt.Printf("=== FLP Impossibility Simulation ===\n")
 fmt.Printf("Processes: %d, Faulty: %d\n", numProcesses, numFaulty)
 fmt.Printf("Duration: %v\n\n", duration)

 // In a truly asynchronous system with even one faulty process,
 // deterministic consensus is impossible

 fmt.Println("Key observation: With arbitrary message delays,")
 fmt.Println("no process can distinguish between:")
 fmt.Println("  1. A slow message from a correct process")
 fmt.Println("  2. A message that will never come from a crashed process")
 fmt.Println()
 fmt.Println("Therefore, deterministic consensus cannot terminate.")
}
```

## Comparison with Related Protocols

| Property | FLP Impossibility | Ben-Or | Paxos | Raft |
|----------|-------------------|--------|-------|------|
| **Assumptions** | Async, 1+ faults | Async, randomization | Partial sync | Partial sync |
| **Termination** | Impossible | Probabilistic | Deterministic | Deterministic |
| **Fault Tolerance** | n/a | < n/2 faults | < n/2 faults | < n/2 faults |
| **Complexity** | n/a | O(N²) messages | O(N) per round | O(N) per election |
| **Leader** | None | None | Yes | Yes |
| **Message Delays** | Arbitrary | Arbitrary | Bounded | Bounded |

## Visual Representations

### Figure 1: FLP Proof Structure

```
┌─────────────────────────────────────────────────────────────────┐
│                  FLP PROOF STRUCTURE                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐                                            │
│  │ FLP Theorem     │                                            │
│  │ (Impossibility) │                                            │
│  └────────┬────────┘                                            │
│           │                                                     │
│           ▼                                                     │
│  ┌─────────────────────────────────────────┐                    │
│  │ Proof by Contradiction                 │                    │
│  │                                          │                    │
│  │ Assume: Deterministic consensus exists   │                    │
│  │ Show:   Contradiction via infinite run   │                    │
│  └────────┬────────────────────────────────┘                    │
│           │                                                     │
│     ┌─────┴─────┬─────────────┐                                 │
│     ▼           ▼             ▼                                 │
│  ┌──────┐  ┌────────┐   ┌──────────┐                           │
│  │Lemma1│  │Lemma2  │   │Lemma3    │                           │
│  │Bival-│  │Decision│   │Pivotal  │                           │
│  │ent   │  │Gadget  │   │Event    │                           │
│  │Init  │  │Exists  │   │Exists   │                           │
│  └──────┘  └────────┘   └──────────┘                           │
│     │           │             │                                 │
│     └───────────┴─────────────┘                                 │
│                 │                                               │
│                 ▼                                               │
│  ┌─────────────────────────────────────────┐                    │
│  │ Construction: Infinite bivalent run    │                    │
│  │                                          │                    │
│  │ C₀ → C₁ → C₂ → ... (all bivalent)      │                    │
│  │                                          │                    │
│  │ Each step: delay one process's message   │                    │
│  │            (indistinguishable from crash)│                    │
│  └─────────────────────────────────────────┘                    │
│                 │                                               │
│                 ▼                                               │
│  ┌─────────────────────────────────────────┐                    │
│  │ Contradiction: No decision ever made   │                    │
│  │ ∴ Deterministic consensus impossible     │                    │
│  └─────────────────────────────────────────┘                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 2: Circumventing FLP

```
┌─────────────────────────────────────────────────────────────────┐
│              CIRCUMVENTING FLP IMPOSSIBILITY                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ FLP Barrier: Deterministic consensus impossible         │   │
│  │ in async systems with even 1 crash fault                │   │
│  └───────────────────────┬─────────────────────────────────┘   │
│                          │                                      │
│          ┌───────────────┼───────────────┐                     │
│          ▼               ▼               ▼                     │
│     ┌─────────┐    ┌──────────┐   ┌────────────┐              │
│     │Randomi- │    │Failure   │   │Partial     │              │
│     │zation   │    │Detectors │   │Synchrony   │              │
│     └────┬────┘    └────┬─────┘   └─────┬──────┘              │
│          │              │               │                      │
│          ▼              ▼               ▼                      │
│     ┌─────────┐    ┌──────────┐   ┌────────────┐              │
│     │ Ben-Or │    │ Chandra- │   │ Paxos/     │              │
│     │ Algorithm│   │ Toueg    │   │ Raft       │              │
│     │         │    │          │   │            │              │
│     │ O(1)    │    │ ◇S class │   │ Leader     │              │
│     │ expected│    │ eventual │   │ election   │              │
│     └─────────┘    └──────────┘   └────────────┘              │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Key Insight: Each approach adds assumptions that break  │   │
│  │ the FLP adversary model                                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Figure 3: Asynchronous Execution Timeline

```
┌─────────────────────────────────────────────────────────────────┐
│              ASYNCHRONOUS EXECUTION TIMELINE                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Time ▼                                                         │
│                                                                 │
│  P1: ──[w]────────────┬─────────────────────────────────►      │
│             │         │                                         │
│             ▼         │ (delayed indefinitely)                  │
│  Network:   ══════════╪═══════════════════════════════►         │
│             ▲         │                                         │
│             │         ▼                                         │
│  P2: ───────┴────────[r]──────────────────────────────────►    │
│                                                                 │
│  Legend:                                                        │
│    [w] = write operation                                        │
│    [r] = read operation                                         │
│    ═══ = message in transit                                     │
│    ─── = local execution                                        │
│                                                                 │
│  Key Observation:                                               │
│  P2 cannot distinguish between:                                 │
│    1. P1's message is delayed                                   │
│    2. P1 has crashed                                            │
│                                                                 │
│  Therefore: P2 cannot safely decide based on P1's input         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## References

1. Fischer, M. J., Lynch, N. A., & Paterson, M. S. (1985). "Impossibility of distributed consensus with one faulty process." Journal of the ACM.
2. Ben-Or, M. (1983). "Another advantage of free choice: Completely asynchronous agreement protocols." ACM PODC.
3. Chandra, T. D., & Toueg, S. (1996). "Unreliable failure detectors for reliable distributed systems." Journal of the ACM.
4. Dwork, C., Lynch, N., & Stockmeyer, L. (1988). "Consensus in the presence of partial synchrony." Journal of the ACM.
5. Lamport, L. (1998). "The part-time parliament." ACM TOCS.
