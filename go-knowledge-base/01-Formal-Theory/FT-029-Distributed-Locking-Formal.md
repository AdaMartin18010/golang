# FT-029: Distributed Locking - Formal Theory and Analysis

> **Dimension**: Formal Theory
> **Level**: S (>15KB)
> **Tags**: #distributed-locking #mutex #deadlock #redlock #consensus
> **Authoritative Sources**:
>
> - Lamport, L. (1978). "Time, Clocks, and the Ordering of Events in a Distributed System". CACM
> - Rosenkrantz, D. J., et al. (1978). "System Level Concurrency Control for Distributed Database Systems". ACM TODS
> - Redlock Algorithm. Redis Documentation (2016)
> - Chandra, T. D., et al. (2007). "Paxos Made Live". ACM SIGOPS

---

## 1. Theoretical Foundations

### 1.1 Problem Definition

**Definition 1.1 (Distributed Mutual Exclusion)**: Given a distributed system with $n$ processes, ensure that at most one process can execute in its critical section at any time.

**Formal Specification**:

$$
\text{Safety}: \quad \forall t: \sum_{i=1}^{n} \text{inCS}_i(t) \leq 1
$$

$$
\text{Liveness}: \quad \forall i: \Diamond\square(\text{request}_i \rightarrow \Diamond\text{enter}_i)
$$

$$
\text{Fairness}: \quad \forall i, j: \text{request}_i \prec \text{request}_j \rightarrow \text{enter}_i \prec \text{enter}_j
$$

### 1.2 System Model

**Message Passing Model**:

$$
\mathcal{M} = \langle \Pi, \mathcal{C}, \Delta \rangle
$$

where:

- $\Pi$: Set of processes
- $\mathcal{C}$: Communication channels with properties:
  - Reliable: Messages may be lost but not corrupted
  - FIFO: Messages delivered in order
- $\Delta$: Message delay (bounded or unbounded)

**Properties**:

| Property | Description |
|----------|-------------|
| Fairness | Requests granted in order |
| Deadlock-free | No infinite waiting |
| Starvation-free | Every request eventually granted |
| Fault tolerance | Tolerates $f$ faults |

### 1.3 Impossibility Results

**Theorem 1.1 (Consensus Impossibility for Locking)**: In an asynchronous system with even one faulty process, no deterministic algorithm can provide both safety and liveness for distributed locking.

*Proof*:

- Distributed locking reduces to consensus (agreeing on who holds lock)
- By FLP impossibility, consensus is impossible in async systems with even one fault ∎

**Circumvention**: Use timeouts, failure detectors, or randomized algorithms.

---

## 2. Token-Based Algorithms

### 2.1 Suzuki-Kasami Algorithm

**Concept**: A unique token circulates among processes. Only the token holder can enter the critical section.

**State Variables**:

| Variable | Type | Description |
|----------|------|-------------|
| $RN$ | $\mathbb{N}^n$ | Request numbers from all processes |
| $LN$ | $\mathbb{N}^n$ | Last request numbers executed |
| $Q$ | Queue | Pending requests |
| $token$ | Boolean | Token possession |

**Algorithm 1: Suzuki-Kasami**:

```
Process p_i:

On RequestCS:
  if not hasToken:
    RN[i] ← RN[i] + 1
    Broadcast(REQUEST, i, RN[i])
    Wait until token received
  EnterCS()

On ReceiveREQUEST(j, n):
  RN[j] ← max(RN[j], n)
  if hasToken and not inCS:
    if RN[j] = LN[j] + 1:
      SendToken(j)

On ExitCS:
  LN[i] ← RN[i]
  for each j where RN[j] = LN[j] + 1:
    if j not in Q:
      Q.enqueue(j)
  if Q not empty:
    next ← Q.dequeue()
    SendToken(next)
  hasToken ← false
```

**Theorem 2.1 (Suzuki-Kasami Safety)**: At most one process holds the token at any time.

*Proof*: Token is unique and only transferred explicitly ∎

**Theorem 2.2 (Suzuki-Kasami Message Complexity)**: $O(n)$ messages per CS entry in worst case, $O(1)$ amortized.

---

## 3. Permission-Based Algorithms

### 3.1 Ricart-Agrawala Algorithm

**Concept**: Process enters CS after receiving permission (REPLY) from all other processes.

**Algorithm 2: Ricart-Agrawala**:

```
Process p_i:

State:
  ourSeqNum: integer ← 0
  highestSeqNum: integer ← 0
  outstandingReply: integer ← 0
  replyDeferred: array[1..n] of boolean ← {false}
  requestingCS: boolean ← false

On RequestCS:
  requestingCS ← true
  ourSeqNum ← highestSeqNum + 1
  outstandingReply ← n - 1
  for all j ≠ i:
    Send(REQUEST, ourSeqNum, i) to p_j
  Wait until outstandingReply = 0
  EnterCS()

On ReceiveREQUEST(seqNum, j):
  highestSeqNum ← max(highestSeqNum, seqNum)
  defer ← false

  if requestingCS:
    if (seqNum, j) > (ourSeqNum, i):  // Compare tuples
      defer ← true

  if defer:
    replyDeferred[j] ← true
  else:
    Send(REPLY) to p_j

On ExitCS:
  requestingCS ← false
  for all j ≠ i:
    if replyDeferred[j]:
      replyDeferred[j] ← false
      Send(REPLY) to p_j
```

**Theorem 3.1 (Ricart-Agrawala Safety)**: At most one process in CS at any time.

*Proof*:

- Process enters CS only after $n-1$ REPLY messages
- Process defers REPLY if its request has higher priority (lower seqNum or tie-breaker)
- If $p_i$ and $p_j$ both in CS, each must have sent REPLY to other
- But each would defer based on priority, contradiction ∎

**Theorem 3.2 (Ricart-Agrawala Message Complexity)**: $2(n-1)$ messages per CS entry.

---

## 4. Quorum-Based Algorithms

### 4.1 Maekawa's Algorithm

**Concept**: Each process requests permission from a quorum of processes, not all.

**Quorum Construction**:

$$
\forall i, j: V_i \cap V_j \neq \emptyset
$$

where $V_i$ is the voting set for process $p_i$.

**Optimal Quorum Size**: $|V_i| = K \approx \sqrt{n}$, giving $O(\sqrt{n})$ messages.

**Algorithm 3: Maekawa's Algorithm**:

```
Process p_i:

State:
  voted: boolean ← false
  lockedProcess: Process ← nil
  requestQueue: Queue ← empty

On RequestCS:
  for all j in V_i:
    Send(REQUEST, i) to p_j
  Wait until |V_i| replies received
  EnterCS()

On ReceiveREQUEST(j):
  if not voted:
    voted ← true
    lockedProcess ← j
    Send(REPLY) to p_j
  else:
    requestQueue.enqueue(j)

On ReleaseCS:
  for all j in V_i:
    Send(RELEASE) to p_j

On ReceiveRELEASE(j):
  if requestQueue not empty:
    next ← requestQueue.dequeue()
    lockedProcess ← next
    Send(REPLY) to next
  else:
    voted ← false
    lockedProcess ← nil
```

**Theorem 4.1 (Maekawa Safety)**: At most one process in CS at any time.

*Proof*:

- $p_i$ enters CS only after receiving OK from all in $V_i$
- If $p_j$ also in CS, must have OK from all in $V_j$
- $V_i \cap V_j \neq \emptyset$, so some process voted for both
- But voted flag ensures at most one vote per process, contradiction ∎

---

## 5. Redlock Algorithm

### 5.1 Redis Distributed Lock

**Algorithm 4: Redlock**:

```
Algorithm Redlock(resource, ttl):

  // Attempt to acquire lock on N Redis instances
  value ← GenerateUniqueToken()
  quorum ← N / 2 + 1

  for each Redis instance:
    reply ← SET resource value NX PX ttl
    if reply = OK:
      locksAcquired ← locksAcquired + 1

  if locksAcquired >= quorum:
    // Check validity time
    validityTime ← ttl - ClockDrift - (CurrentTime - StartTime)
    if validityTime > 0:
      return LOCK_ACQUIRED(value, validityTime)

  // Failed, release all locks
  for each Redis instance:
    EVAL "if redis.call('get', KEYS[1]) == ARGV[1] then ..."

  return LOCK_FAILED

Algorithm Unlock(resource, value):
  for each Redis instance:
    EVAL "if redis.call('get', KEYS[1]) == ARGV[1] then
            return redis.call('del', KEYS[1])
          else return 0 end"
```

**Safety Properties**:

1. Quorum ensures at most one client holds lock
2. Token ensures only lock holder can unlock
3. TTL provides liveness (lock expires)

---

## 6. TLA+ Specifications

### 6.1 Distributed Mutex TLA+

```tla
----------------------------- MODULE DistributedMutex -----------------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Processes,
          MaxRequests

VARIABLES state,        \* State: idle, waiting, critical
          requestNum,   \* Request sequence numbers
          replyCount,   \* Outstanding replies
          inCS          \* Set of processes in critical section

Init ==
  /\ state = [p \in Processes |-> "idle"]
  /\ requestNum = [p \in Processes |-> 0]
  /\ replyCount = [p \in Processes |-> 0]
  /\ inCS = {}

Request(p) ==
  /\ state[p] = "idle"
  /\ state' = [state EXCEPT ![p] = "waiting"]
  /\ requestNum' = [requestNum EXCEPT ![p] = @ + 1]
  /\ replyCount' = [replyCount EXCEPT ![p] = Cardinality(Processes) - 1]
  /\ UNCHANGED inCS

ReceiveReply(p) ==
  /\ state[p] = "waiting"
  /\ replyCount[p] > 0
  /\ replyCount' = [replyCount EXCEPT ![p] = @ - 1]
  /\ UNCHANGED <<state, requestNum, inCS>>

EnterCS(p) ==
  /\ state[p] = "waiting"
  /\ replyCount[p] = 0
  /\ state' = [state EXCEPT ![p] = "critical"]
  /\ inCS' = inCS \union {p}
  /\ UNCHANGED <<requestNum, replyCount>>

ExitCS(p) ==
  /\ state[p] = "critical"
  /\ state' = [state EXCEPT ![p] = "idle"]
  /\ inCS' = inCS \\ {p}
  /\ UNCHANGED <<requestNum, replyCount>>

Next ==
  \E p \in Processes:
    Request(p) \/ ReceiveReply(p) \/ EnterCS(p) \/ ExitCS(p)

Safety ==
  \* Mutual exclusion
  Cardinality(inCS) <= 1

Liveness ==
  \* No starvation (with fairness assumption)
  \A p \in Processes:
    state[p] = "waiting" ~> state[p] = "critical"

=============================================================================
```

---

## 7. Go Implementation

```go
// Package distlock provides distributed locking implementations
package distlock

import (
 "context"
 "crypto/rand"
 "encoding/hex"
 "fmt"
 "sort"
 "sync"
 "sync/atomic"
 "time"
)

// ============================================
// Common Interfaces
// ============================================

// Lock represents a distributed lock
type Lock interface {
 Acquire(ctx context.Context) error
 Release() error
 IsHeld() bool
}

// LockProvider provides lock instances
type LockProvider interface {
 GetLock(resource string) Lock
}

// Quorum represents a quorum of nodes
type Quorum interface {
 Size() int
 Members() []string
}

// ============================================
// Token-Based Lock (Suzuki-Kasami)
// ============================================

type TokenLock struct {
 mu        sync.Mutex
 nodeID    string
 nodes     []string
 hasToken  bool
 tokenCh   chan struct{}

 RN        map[string]uint64 // Request numbers
 LN        map[string]uint64 // Last executed
 queue     []string          // Pending requests

 transport Transport
}

type Transport interface {
 Broadcast(msg Message) error
 Send(to string, msg Message) error
}

type Message struct {
 Type      string
 From      string
 RequestID uint64
}

// NewTokenLock creates a new token-based lock
func NewTokenLock(nodeID string, nodes []string, transport Transport) *TokenLock {
 // Determine initial token holder (lowest ID)
 hasToken := true
 for _, n := range nodes {
  if n < nodeID {
   hasToken = false
   break
  }
 }

 tl := &TokenLock{
  nodeID:    nodeID,
  nodes:     nodes,
  hasToken:  hasToken,
  tokenCh:   make(chan struct{}, 1),
  RN:        make(map[string]uint64),
  LN:        make(map[string]uint64),
  transport: transport,
 }

 if hasToken {
  tl.tokenCh <- struct{}{}
 }

 for _, n := range nodes {
  tl.RN[n] = 0
  tl.LN[n] = 0
 }

 return tl
}

// Acquire acquires the lock
func (l *TokenLock) Acquire(ctx context.Context) error {
 l.mu.Lock()
 l.RN[l.nodeID]++
 myReq := l.RN[l.nodeID]
 l.mu.Unlock()

 if l.hasToken {
  return nil
 }

 // Broadcast request
 msg := Message{
  Type:      "REQUEST",
  From:      l.nodeID,
  RequestID: myReq,
 }
 l.transport.Broadcast(msg)

 // Wait for token
 select {
 case <-l.tokenCh:
  l.mu.Lock()
  l.hasToken = true
  l.mu.Unlock()
  return nil
 case <-ctx.Done():
  return ctx.Err()
 }
}

// Release releases the lock
func (l *TokenLock) Release() error {
 l.mu.Lock()
 defer l.mu.Unlock()

 if !l.hasToken {
  return fmt.Errorf("don't have token")
 }

 l.LN[l.nodeID] = l.RN[l.nodeID]

 // Check pending requests
 for _, node := range l.nodes {
  if l.RN[node] == l.LN[node]+1 {
   // Check if already in queue
   found := false
   for _, q := range l.queue {
    if q == node {
     found = true
     break
    }
   }
   if !found {
    l.queue = append(l.queue, node)
   }
  }
 }

 // Send token to next in queue
 if len(l.queue) > 0 {
  next := l.queue[0]
  l.queue = l.queue[1:]
  l.hasToken = false

  msg := Message{Type: "TOKEN", From: l.nodeID}
  l.transport.Send(next, msg)
 }

 return nil
}

// HandleMessage processes incoming messages
func (l *TokenLock) HandleMessage(msg Message) {
 l.mu.Lock()
 defer l.mu.Unlock()

 switch msg.Type {
 case "REQUEST":
  if msg.RequestID > l.RN[msg.From] {
   l.RN[msg.From] = msg.RequestID
  }

  // Send token if we have it and not in CS
  if l.hasToken && len(l.queue) == 0 {
   if l.RN[msg.From] == l.LN[msg.From]+1 {
    l.hasToken = false
    go l.transport.Send(msg.From, Message{Type: "TOKEN"})
   }
  }

 case "TOKEN":
  select {
  case l.tokenCh <- struct{}{}:
  default:
  }
 }
}

// ============================================
// Quorum Lock (Maekawa)
// ============================================

type QuorumLock struct {
 nodeID       string
 votingSet    []string
 transport    Transport

 mu           sync.Mutex
 voted        bool
 lockedBy     string
 requestQueue []string

 waiting     bool
 votesNeeded int
 voteCh      chan struct{}
}

// NewQuorumLock creates a new quorum-based lock
func NewQuorumLock(nodeID string, allNodes []string, transport Transport) *QuorumLock {
 // Compute voting set (simplified: use sqrt(n) nodes)
 k := int(float64(len(allNodes)) * 0.6) // sqrt approximation
 if k < 2 {
  k = 2
 }

 // Deterministic assignment
 votingSet := make([]string, 0, k)
 for _, n := range allNodes {
  if len(votingSet) < k {
   votingSet = append(votingSet, n)
  }
 }

 return &QuorumLock{
  nodeID:    nodeID,
  votingSet: votingSet,
  transport: transport,
  voteCh:    make(chan struct{}, 1),
 }
}

// Acquire acquires the lock
func (l *QuorumLock) Acquire(ctx context.Context) error {
 l.mu.Lock()
 l.waiting = true
 l.votesNeeded = len(l.votingSet)
 l.mu.Unlock()

 // Send requests to voting set
 for _, node := range l.votingSet {
  msg := Message{
   Type: "REQUEST",
   From: l.nodeID,
  }
  l.transport.Send(node, msg)
 }

 // Wait for quorum
 select {
 case <-l.voteCh:
  l.mu.Lock()
  l.waiting = false
  l.mu.Unlock()
  return nil
 case <-ctx.Done():
  return ctx.Err()
 }
}

// Release releases the lock
func (l *QuorumLock) Release() error {
 for _, node := range l.votingSet {
  msg := Message{
   Type: "RELEASE",
   From: l.nodeID,
  }
  l.transport.Send(node, msg)
 }
 return nil
}

// HandleMessage processes incoming messages
func (l *QuorumLock) HandleMessage(msg Message) {
 l.mu.Lock()
 defer l.mu.Unlock()

 switch msg.Type {
 case "REQUEST":
  if !l.voted {
   l.voted = true
   l.lockedBy = msg.From
   go l.transport.Send(msg.From, Message{Type: "REPLY"})
  } else {
   l.requestQueue = append(l.requestQueue, msg.From)
  }

 case "REPLY":
  if l.waiting {
   l.votesNeeded--
   if l.votesNeeded <= 0 {
    select {
    case l.voteCh <- struct{}{}:
    default:
    }
   }
  }

 case "RELEASE":
  if len(l.requestQueue) > 0 {
   next := l.requestQueue[0]
   l.requestQueue = l.requestQueue[1:]
   l.lockedBy = next
   go l.transport.Send(next, Message{Type: "REPLY"})
  } else {
   l.voted = false
   l.lockedBy = ""
  }
 }
}

// ============================================
// Redlock Implementation
// ============================================

type RedisInstance interface {
 SetNX(key, value string, ttl time.Duration) (bool, error)
 Eval(script string, keys []string, args []string) (interface{}, error)
}

type Redlock struct {
 instances []RedisInstance
 quorum    int
 ttl       time.Duration
 drift     time.Duration

 mu        sync.Mutex
 resource  string
 value     string
 hasLock   bool
 validity  time.Duration
}

// NewRedlock creates a new Redlock instance
func NewRedlock(instances []RedisInstance, ttl time.Duration) *Redlock {
 return &Redlock{
  instances: instances,
  quorum:    len(instances)/2 + 1,
  ttl:       ttl,
  drift:     time.Millisecond * 10,
 }
}

// Lock acquires a distributed lock
func (r *Redlock) Lock(resource string) (*LockContext, error) {
 value := generateToken()

 startTime := time.Now()
 n := 0

 for _, inst := range r.instances {
  ok, err := inst.SetNX(resource, value, r.ttl)
  if err == nil && ok {
   n++
  }
 }

 elapsed := time.Since(startTime)
 validity := r.ttl - elapsed - r.drift

 if n >= r.quorum && validity > 0 {
  r.mu.Lock()
  r.resource = resource
  r.value = value
  r.hasLock = true
  r.validity = validity
  r.mu.Unlock()

  return &LockContext{
   Validity: validity,
   Unlock: func() {
    r.Unlock()
   },
  }, nil
 }

 // Failed, unlock all
 r.unlockAll(resource, value)
 return nil, fmt.Errorf("failed to acquire lock")
}

// Unlock releases the lock
func (r *Redlock) Unlock() error {
 r.mu.Lock()
 resource := r.resource
 value := r.value
 r.hasLock = false
 r.mu.Unlock()

 if resource == "" {
  return nil
 }

 r.unlockAll(resource, value)
 return nil
}

func (r *Redlock) unlockAll(resource, value string) {
 script := `
  if redis.call("get", KEYS[1]) == ARGV[1] then
   return redis.call("del", KEYS[1])
  else
   return 0
  end
 `

 for _, inst := range r.instances {
  inst.Eval(script, []string{resource}, []string{value})
 }
}

type LockContext struct {
 Validity time.Duration
 Unlock   func()
}

func generateToken() string {
 b := make([]byte, 16)
 rand.Read(b)
 return hex.EncodeToString(b)
}
```

---

## 8. Visual Representations

### 8.1 Distributed Locking Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DISTRIBUTED LOCKING ALGORITHMS TAXONOMY                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Distributed Mutual Exclusion                                                │
│  │                                                                           │
│  ├── Token-Based                                                             │
│  │   ├── Suzuki-Kasami (broadcast requests)                                  │
│  │   ├── Raymond (tree-based)                                                │
│  │   └── Token Ring (logical ring)                                           │
│  │                                                                           │
│  ├── Permission-Based                                                        │
│  │   ├── Ricart-Agrawala (all-to-all)                                        │
│  │   ├── Lamport (logical clocks)                                            │
│  │   └── Multiple (k-out-of-n permissions)                                   │
│  │                                                                           │
│  ├── Quorum-Based                                                            │
│  │   ├── Maekawa (√n quorum)                                                 │
│  │   ├── Voting (weighted quorums)                                           │
│  │   └── Grid Protocol (2D quorums)                                          │
│  │                                                                           │
│  └── Consensus-Based                                                         │
│      ├── Paxos Lock (state machine replication)                              │
│      ├── Raft Lock (leader-based)                                            │
│      └── Redlock (quorum on external store)                                  │
│                                                                              │
│  Comparison:                                                                 │
│  ┌────────────────┬───────────┬─────────────┬────────────┬────────────┐   │
│  │ Algorithm      │ Messages  │ Delay       │ Fault Tol. │ Fairness   │   │
│  ├────────────────┼───────────┼─────────────┼────────────┼────────────┤   │
│  │ Suzuki-Kasami  │ 0 or n    │ 0 or T      │ No         │ FIFO       │   │
│  │ Ricart-Agrawala│ 2(n-1)    │ T           │ No         │ FIFO       │   │
│  │ Maekawa        │ 2√n       │ 2T          │ Limited    │ FIFO       │   │
│  │ Redlock        │ N         │ T           │ Yes (f<N/2)│ No         │   │
│  └────────────────┴───────────┴─────────────┴────────────┴────────────┘   │
│                                                                              │
│  T = message delay, N = number of lock instances                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Deadlock Detection and Prevention

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DEADLOCK HANDLING IN DISTRIBUTED SYSTEMS                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  DEADLOCK PREVENTION                                                         │
│  ─────────────────────                                                       │
│                                                                              │
│  1. Resource Ordering                                                        │
│     Assign total order to resources: R1 < R2 < R3 < ... < Rn                │
│     Processes must acquire locks in ascending order                          │
│                                                                              │
│     Process A: Lock(R1) → Lock(R2) → Lock(R3)                               │
│     Process B: Lock(R1) → Lock(R2) → Lock(R3)  [Same order - no deadlock]   │
│                                                                              │
│  2. Timeout-Based (Redlock approach)                                         │
│     • All locks have TTL                                                     │
│     • If process fails, lock automatically releases                          │
│     • Requires fencing tokens for safety                                     │
│                                                                              │
│  Time: ────[Lock]────────[Renew]────────[Renew]───────[Expire]               │
│        │               │               │                                    │
│        └── TTL ────────┴── TTL ───────┴── TTL ────>                       │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────           │
│                                                                              │
│  DEADLOCK DETECTION                                                          │
│  ──────────────────                                                          │
│                                                                              │
│  Wait-For Graph (distributed):                                               │
│                                                                              │
│  R1 ──waits──> R2                                                            │
│  ▲             │                                                             │
│  │             waits                                                        │
│  │             ▼                                                             │
│  └──────────── R3                                                            │
│                                                                              │
│  Cycle detected: Deadlock!                                                   │
│  Resolution: Abort youngest transaction                                      │
│                                                                              │
│  Chandy-Misra-Haas Algorithm:                                                │
│  • Probe messages traverse wait-for graph                                    │
│  • If probe returns to initiator: deadlock detected                          │
│  • Message complexity: O(e) where e = edges in wait-for graph               │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────           │
│                                                                              │
│  DEADLOCK AVOIDANCE                                                          │
│  ──────────────────                                                          │
│                                                                              │
│  Banker's Algorithm (distributed variant):                                   │
│  • Each process declares max resources needed                                │
│  • System checks if granting request leaves safe state                       │
│  • If unsafe, request is delayed                                             │
│                                                                              │
│  Safe State: There exists an ordering of processes such that                 │
│              each can complete with available resources                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.3 Redlock Execution Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    REDLOCK EXECUTION FLOW                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Client wants to acquire lock on "resource:X"                               │
│                                                                              │
│  Step 1: Generate unique token                                               │
│  ┌─────────────┐                                                             │
│  │ Token: abc  │                                                             │
│  └──────┬──────┘                                                             │
│         │                                                                    │
│  Step 2: Attempt to acquire lock on N Redis instances                        │
│         │                                                                    │
│         ├──────────SET resource:X abc NX PX 30000─────────> Redis-1         │
│         ├──────────SET resource:X abc NX PX 30000─────────> Redis-2         │
│         ├──────────SET resource:X abc NX PX 30000─────────> Redis-3         │
│         ├──────────SET resource:X abc NX PX 30000─────────> Redis-4         │
│         └──────────SET resource:X abc NX PX 30000─────────> Redis-5         │
│         │                                                                    │
│  Step 3: Collect responses (parallel)                                        │
│         │                                                                    │
│         │<──────────────OK────────────────────────────────Redis-1           │
│         │<──────────────OK────────────────────────────────Redis-2           │
│         │<──────────────OK────────────────────────────────Redis-3           │
│         │<──────────────(timeout)─────────────────────────Redis-4           │
│         │<──────────────OK────────────────────────────────Redis-5           │
│         │                                                                    │
│  Step 4: Count successful acquisitions                                       │
│         │                                                                    │
│         │  Acquired: 4 locks  (Redis 1,2,3,5)                                │
│         │  Quorum needed: N/2 + 1 = 3                                        │
│         │  Quorum achieved: YES                                              │
│         │                                                                    │
│  Step 5: Compute validity time                                               │
│         │                                                                    │
│         │  Start time: T0                                                    │
│         │  End time: T1                                                      │
│         │  Elapsed: T1 - T0 = 50ms                                           │
│         │  TTL: 30000ms                                                      │
│         │  Drift factor: 10ms                                                │
│         │                                                                    │
│         │  Validity = TTL - Elapsed - Drift                                  │
│         │          = 30000 - 50 - 10 = 29940ms                               │
│         │                                                                    │
│  Step 6: Check validity                                                      │
│         │                                                                    │
│         │  Validity (29940ms) > 0: YES                                       │
│         │                                                                    │
│         ▼                                                                    │
│  ┌─────────────────────┐                                                     │
│  │ LOCK ACQUIRED       │                                                     │
│  │ Token: abc          │                                                     │
│  │ Validity: ~29.9s    │                                                     │
│  └─────────────────────┘                                                     │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────           │
│                                                                              │
│  UNLOCK SEQUENCE:                                                            │
│                                                                              │
│  Client unlocks using same token:                                            │
│                                                                              │
│  EVAL "if redis.call('get', KEYS[1]) == ARGV[1] then                         │
│           return redis.call('del', KEYS[1])                                  │
│        else return 0 end"                                                    │
│  KEYS: [resource:X]                                                          │
│  ARGV: [abc]  ← Must match or unlock fails (fencing)                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 9. Academic References

1. **Lamport, L. (1978)**. "Time, Clocks, and the Ordering of Events in a Distributed System". *CACM*.

2. **Ricart, G., & Agrawala, A. K. (1981)**. "An Optimal Algorithm for Mutual Exclusion in Computer Networks". *CACM*.

3. **Maekawa, M. (1985)**. "A √N Algorithm for Mutual Exclusion in Decentralized Systems". *ACM TOCS*.

4. **Suzuki, I., & Kasami, T. (1985)**. "A Distributed Mutual Exclusion Algorithm". *ACM TOCS*.

5. **Antirez (2016)**. "Redis Distlock - A distributed lock manager". *Redis Documentation*.

---

## 10. Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DISTRIBUTED LOCKING SUMMARY                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Algorithm Selection Guide:                                                  │
│                                                                              │
│  ┌─────────────────────┬────────────────────────────────────────────────┐   │
│  │ Scenario            │ Recommended Algorithm                          │   │
│  ├─────────────────────┼────────────────────────────────────────────────┤   │
│  │ Small cluster (<10) │ Ricart-Agrawala                                │   │
│  │                     │ Simple, fair, 2(n-1) messages                 │   │
│  ├─────────────────────┼────────────────────────────────────────────────┤   │
│  │ Large cluster       │ Maekawa                                        │   │
│  │                     │ O(√n) messages                                 │   │
│  ├─────────────────────┼────────────────────────────────────────────────┤   │
│  │ High availability   │ Redlock                                        │   │
│  │ (with Redis)        │ Fault-tolerant, requires clock sync           │   │
│  ├─────────────────────┼────────────────────────────────────────────────┤   │
│  │ Strong consistency  │ Paxos/Raft based                               │   │
│  │                     │ Consensus-level safety                         │   │
│  └─────────────────────┴────────────────────────────────────────────────┘   │
│                                                                              │
│  Critical Considerations:                                                    │
│  1. Clock synchronization for TTL-based locks                               │
│  2. Fencing tokens to prevent delayed unlocks                               │
│  3. Deadlock detection or timeout-based release                             │
│  4. Network partitions and split-brain scenarios                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
