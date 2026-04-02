# FT-024: Consensus Variations - Formal Analysis

> **Dimension**: Formal Theory
> **Level**: S (>15KB)
> **Tags**: #consensus #multi-paxos #epaxos #flexible-paxos #distributed-systems
> **Authoritative Sources**:
>
> - Lamport, L. (2001). "Paxos Made Simple". ACM SIGACT News
> - Moraru, I., et al. (2013). "EPaxos: There Is More Consensus in Egalitarian Parliaments". SOSP
> - Howard, H., et al. (2016). "Flexible Paxos: Quorum Intersection Revisited". OPODIS
> - Van Renesse, R., et al. (2015). "Paxos Made Moderately Complex". ACM CSUR

---

## 1. Theoretical Foundations

### 1.1 System Model for Consensus Variants

**Definition 1.1 (Consensus Protocol Family)**: A consensus protocol family $\mathcal{P}$ is a set of protocols $\{P_1, P_2, ..., P_n\}$ that solve the consensus problem under varying assumptions and optimizations.

$$
\mathcal{P} = \{P | P = \langle \Pi, Q, M, \Delta, \Sigma \rangle\}
$$

where:

- $\Pi$: Set of participating processes
- $Q$: Quorum system
- $M$: Message space
- $\Delta$: Decision function
- $\Sigma$: Safety specification

**Definition 1.2 (Protocol Equivalence)**: Two protocols $P_1$ and $P_2$ are equivalent ($P_1 \equiv P_2$) if they satisfy the same safety and liveness properties:

$$
P_1 \equiv P_2 \iff (\Sigma_{P_1} = \Sigma_{P_2}) \land (\Lambda_{P_1} = \Lambda_{P_2})
$$

### 1.2 Quorum System Variations

**Definition 1.3 (Classic Paxos Quorum)**: The classic Paxos quorum system $\mathcal{Q}_{classic}$ requires any two quorums to intersect:

$$
\forall Q_1, Q_2 \in \mathcal{Q}_{classic}: Q_1 \cap Q_2 \neq \emptyset
$$

**Theorem 1.1 (Quorum Intersection Necessity)**: For consensus safety, any quorum system must satisfy the intersection property.

*Proof*:

- Assume $\exists Q_1, Q_2 \in \mathcal{Q}: Q_1 \cap Q_2 = \emptyset$
- $Q_1$ could decide value $v_1$ while $Q_2$ decides $v_2 \neq v_1$
- This violates agreement
- Therefore, intersection is necessary $\square$

**Definition 1.4 (Flexible Quorum System)**: Howard et al. (2016) showed that phase-1 and phase-2 quorums need not be the same:

$$
\mathcal{Q}_{flexible} = \langle \mathcal{Q}_1, \mathcal{Q}_2 \rangle \text{ where } \forall Q_1 \in \mathcal{Q}_1, \forall Q_2 \in \mathcal{Q}_2: Q_1 \cap Q_2 \neq \emptyset
$$

**Theorem 1.2 (Flexible Quorum Safety)**: Flexible Paxos maintains safety if $|Q_1| + |Q_2| > n$ for any $Q_1 \in \mathcal{Q}_1$, $Q_2 \in \mathcal{Q}_2$.

*Proof*:

- By pigeonhole principle: if $|Q_1| + |Q_2| > n$ and $Q_1, Q_2 \subseteq \Pi$ with $|\Pi| = n$
- Then $|Q_1 \cap Q_2| = |Q_1| + |Q_2| - |Q_1 \cup Q_2| \geq |Q_1| + |Q_2| - n > 0$
- Therefore $Q_1 \cap Q_2 \neq \emptyset$ $\square$

### 1.3 Latency and Throughput Models

**Definition 1.5 (Consensus Latency)**: The latency $L$ of a consensus instance is:

$$
L = \begin{cases}
2\delta & \text{for single-decree (Paxos)} \\
2\delta + k\epsilon & \text{for multi-decree with } k \text{ pipelined} \\
\delta + \max_{i}(\delta_i) & \text{for leaderless (EPaxos)}
\end{cases}
$$

where $\delta$ is network latency and $\epsilon$ is processing overhead.

---

## 2. Multi-Paxos Formalization

### 2.1 State Space Extension

**Definition 2.1 (Log-Based State Machine)**: Multi-Paxos extends single-decree Paxos to a sequence of consensus instances:

$$
\mathcal{L} = \langle c_1, c_2, ..., c_n \rangle \text{ where each } c_i \in \mathcal{C}
$$

**Definition 2.2 (Leader State)**: The leader maintains additional state:

| Variable | Type | Description |
|----------|------|-------------|
| $nextIndex$ | $\Pi \rightarrow \mathbb{N}$ | Next slot to send to each follower |
| $matchIndex$ | $\Pi \rightarrow \mathbb{N}$ | Highest slot known replicated |
| $pending$ | $\mathcal{P}(\mathbb{N})$ | Slots with proposed but uncommitted commands |

### 2.2 Multi-Paxos Protocol

**Phase 1: Leader Election** (executed once per leadership term):

$$
\text{Prepare}(n) \rightarrow \forall p \in \Pi: \text{Prepare}(n) \xrightarrow{} p
$$

**Phase 2: Command Replication** (executed for each command):

$$
\forall c \in \text{Commands}: \text{Accept}(n, c) \rightarrow \forall p \in Q: \text{Accept}(n, c) \xrightarrow{} p
$$

**Theorem 2.1 (Multi-Paxos Throughput)**: Multi-Paxos achieves throughput $\Theta(N/\delta)$ where $N$ is the number of replicas and $\delta$ is network latency.

*Proof*:

- After leader election (1 RTT), leader can pipeline commands
- Each command requires 1 RTT (leader to followers)
- Leader can send up to $N-1$ commands in parallel
- Throughput = $\frac{N-1}{\delta} = \Theta(N/\delta)$ $\square$

---

## 3. EPaxos (Egalitarian Paxos) Formalization

### 3.1 Dependency Tracking

**Definition 3.1 (Conflict Relation)**: Commands $c_1$ and $c_2$ conflict if they access overlapping keys with at least one write:

$$
\text{Conflict}(c_1, c_2) \iff (\text{Keys}(c_1) \cap \text{Keys}(c_2) \neq \emptyset) \land (\text{Write}(c_1) \lor \text{Write}(c_2))
$$

**Definition 3.2 (Dependency Graph)**: Each command $c$ carries a dependency set:

$$
\text{Deps}(c) = \{c' | \text{Conflict}(c, c') \land \text{Seen}(c')\}
$$

### 3.2 Fast Path vs Slow Path

**Definition 3.3 (Fast Quorum)**: EPaxos uses fast quorums of size $F + \lfloor(F+1)/2\rfloor$ where $F$ is the fault tolerance:

$$
Q_{fast} = F + \left\lfloor\frac{F+1}{2}\right\rfloor + 1
$$

For $n = 2F + 1$ replicas: $Q_{fast} = \lceil 3n/4 \rceil$

**Theorem 3.1 (Fast Path Condition)**: A command commits on the fast path if all fast quorum members report identical dependencies.

*Proof*:

- Fast quorum size: $\lceil 3n/4 \rceil$
- Any two fast quorums intersect in at least $\lceil n/2 \rceil$ nodes
- For identical dependencies, intersection nodes must agree
- This guarantees safety even with conflicting commands $\square$

**Theorem 3.2 (EPaxos Latency)**: EPaxos achieves commit latency of:

- Fast path: $\delta$ (1 RTT to closest replica)
- Slow path: $2\delta$ (2 RTTs, same as Paxos)

### 3.3 Execution Algorithm

```
Algorithm: EPaxos Command Execution

State:
  cmds: map<InstanceID, Command>
  deps: map<InstanceID, Set<InstanceID>>
  executed: Set<InstanceID>
  committed: Set<InstanceID>

Procedure ExecuteCommand(cmd):
  // Add to dependency graph
  cmds[cmd.id] = cmd
  deps[cmd.id] = FindConflicts(cmd)

  // Try to execute
  ExecuteReady()

Procedure ExecuteReady():
  // Find commands with all dependencies executed
  ready = {c ∈ committed | deps[c] ⊆ executed}

  while ready ≠ ∅:
    // Execute in reverse topological order
    c = SelectNoIncomingEdges(ready)
    ApplyToStateMachine(cmds[c])
    executed = executed ∪ {c}
    ready = ready \ {c}
```

---

## 4. Flexible Paxos Formalization

### 4.1 Asymmetric Quorums

**Definition 4.1 (Phase-1 Quorum)**: For leader election (phase 1):

$$
\mathcal{Q}_1 = \{Q \subseteq \Pi : |Q| \geq n - f\}
$$

**Definition 4.2 (Phase-2 Quorum)**: For value acceptance (phase 2):

$$
\mathcal{Q}_2 = \{Q \subseteq \Pi : |Q| \geq f + 1\}
$$

**Theorem 4.1 (Flexible Paxos Safety)**: With $\mathcal{Q}_1$ and $\mathcal{Q}_2$ as defined, any $Q_1 \in \mathcal{Q}_1$ and $Q_2 \in \mathcal{Q}_2$ intersect.

*Proof*:

- $|Q_1| \geq n - f$
- $|Q_2| \geq f + 1$
- $|Q_1| + |Q_2| \geq n - f + f + 1 = n + 1 > n$
- By pigeonhole principle: $Q_1 \cap Q_2 \neq \emptyset$ $\square$

### 4.2 Optimization Strategies

**Strategy 1: Large Phase-1, Small Phase-2**:

- Phase-1: Contact all $n$ replicas (slower but safer)
- Phase-2: Contact only $f+1$ replicas (faster commit)

**Strategy 2: Geographic Optimization**:

- Phase-1: Quorum from nearby replicas (low latency)
- Phase-2: Quorum includes phase-1 replicas

---

## 5. TLA+ Specifications

### 5.1 Multi-Paxos TLA+

```tla
----------------------------- MODULE MultiPaxos -----------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Replicas,
          MaxSlot,
          Values,
          MaxBallot

VARIABLES ballot,
          accepted,
          chosen,
          leader,
          log

Ballot == 0..MaxBallot
Slot == 1..MaxSlot

Init ==
  /\ ballot = [r \in Replicas |-> 0]
  /\ accepted = [r \in Replicas |-> [s \in Slot |-> {}]]
  /\ chosen = [s \in Slot |-> {}]
  /\ leader = [r \in Replicas |-> FALSE]
  /\ log = [r \in Replicas |-> [s \in Slot |-> None]]

Phase1a(r, b) ==
  /\ b > ballot[r]
  /\ ballot' = [ballot EXCEPT ![r] = b]
  /\ leader' = [leader EXCEPT ![r] = FALSE]
  /\ UNCHANGED <<accepted, chosen, log>>

Phase1b(r, s, b) ==
  /\ b >= ballot[r]
  /\ ballot' = [ballot EXCEPT ![r] = b]
  /\ UNCHANGED <<accepted, chosen, leader, log>>

Phase2a(r, s, b, v) ==
  /\ leader[r]
  /\ ballot[r] = b
  /\ \A Q \in Quorum: \E q \in Q: <b, v> \in accepted[q][s]
  /\ chosen' = [chosen EXCEPT ![s] = chosen[s] \union {v}]
  /\ log' = [log EXCEPT ![r][s] = v]
  /\ UNCHANGED <<ballot, accepted, leader>>

BecomeLeader(r, b) ==
  /\ ballot[r] = b
  /\ \E Q \in Quorum: \A q \in Q: ballot[q] <= b
  /\ leader' = [leader EXCEPT ![r] = TRUE]
  /\ UNCHANGED <<ballot, accepted, chosen, log>>

Safety ==
  /\ \A s \in Slot: Cardinality(chosen[s]) <= 1
  /\ \A s1, s2 \in Slot, v \in Values:
       v \in chosen[s1] /\ v \in chosen[s2] => s1 = s2

=============================================================================
```

### 5.2 EPaxos TLA+

```tla
----------------------------- MODULE EPaxos -----------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Replicas,
          Commands,
          MaxInstance

VARIABLES inst,
          deps,
          status,
          seq,
          executed

Instance == [replica: Replicas, index: 1..MaxInstance]

CommandStatus == {"preaccepted", "accepted", "committed"}

Init ==
  /\ inst = [r \in Replicas |-> 0]
  /\ deps = [i \in Instance |-> {}]
  /\ status = [i \in Instance |-> "none"]
  /\ seq = [i \in Instance |-> 0]
  /\ executed = {}

PreAccept(r, cmd, initialDeps) ==
  /\ inst[r] < MaxInstance
  /\ inst' = [inst EXCEPT ![r] = @ + 1]
  /\ LET i == [replica |-> r, index |-> inst'[r]]
     IN /\ deps' = [deps EXCEPT ![i] = initialDeps]
        /\ status' = [status EXCEPT ![i] = "preaccepted"]
  /\ UNCHANGED <<seq, executed>>

Accept(r, i, cmd, finalDeps, finalSeq) ==
  /\ status[i] = "preaccepted"
  /\ deps' = [deps EXCEPT ![i] = finalDeps]
  /\ seq' = [seq EXCEPT ![i] = finalSeq]
  /\ status' = [status EXCEPT ![i] = "accepted"]
  /\ UNCHANGED <<inst, executed>>

Commit(i) ==
  /\ status[i] \in {"preaccepted", "accepted"}
  /\ status' = [status EXCEPT ![i] = "committed"]
  /\ UNCHANGED <<inst, deps, seq, executed>>

Execute ==
  /\ LET ready == {i \in Instance : status[i] = "committed" /\ deps[i] \subseteq executed}
     IN /\ ready /= {}
        /\ \E i \in ready:
           /\ executed' = executed \union {i}
  /\ UNCHANGED <<inst, deps, status, seq>>

DependencyCorrectness ==
  \A i \in Instance:
    status[i] = "committed" =>
      \A j \in deps[i]: status[j] = "committed" /\ j \in executed

=============================================================================
```

---

## 6. Algorithm Pseudocode

### 6.1 Multi-Paxos Optimized Leader

```
Algorithm: Multi-Paxos with Batching and Pipelining

State:
  ballot: integer                         // Current ballot number
  isLeader: boolean                       // Leadership status
  log: array[slot] of Command             // Local log
  nextIndex: map<Replica, slot>           // Next slot to send each replica
  matchIndex: map<Replica, slot>          // Highest matching slot each replica
  pendingCommands: queue<Command>         // Commands waiting to be proposed

Constants:
  BATCH_SIZE = 100
  PIPELINE_DEPTH = 1000

On Startup():
  RunLeaderElection()

On ClientRequest(cmd):
  if not isLeader:
    ForwardToLeader(cmd)
    return

  pendingCommands.enqueue(cmd)

  if pendingCommands.size() >= BATCH_SIZE:
    ProposeBatch()

ProposeBatch():
  batch = pendingCommands.dequeue(BATCH_SIZE)
  startSlot = GetNextFreeSlot()

  for i, cmd in enumerate(batch):
    slot = startSlot + i
    log[slot] = cmd
    SendAccept(ballot, slot, cmd)

  StartTimer(slot)  // For detecting slow replicas

On AcceptAck(replica, slot, ballotNum):
  if ballotNum != ballot:
    return  // Stale response

  // Update match index
  matchIndex[replica] = max(matchIndex[replica], slot)

  // Check if slot can be committed
  if CountMatching(slot) > n/2:
    Commit(slot)

  // Advance nextIndex for pipelining
  if slot == nextIndex[replica] - 1:
    nextIndex[replica] = slot + 1
    SendMoreCommands(replica)

CountMatching(slot):
  count = 0
  for replica in Replicas:
    if matchIndex[replica] >= slot:
      count++
  return count

SendMoreCommands(replica):
  // Pipeline more commands to this replica
  while CanPipeline(replica):
    slot = nextIndex[replica]
    if log[slot] != null:
      SendAccept(ballot, slot, log[slot], replica)
      nextIndex[replica]++
    else:
      break

CanPipeline(replica):
  return (nextIndex[replica] - matchIndex[replica]) < PIPELINE_DEPTH
```

### 6.2 EPaxos Replica

```
Algorithm: EPaxos Replica Protocol

State:
  cmdId: integer                          // Local command counter
  cmds: map<InstanceID, Command>          // Command storage
  deps: map<InstanceID, Set<InstanceID>>  // Dependencies
  status: map<InstanceID, Status>         // Command status
  seq: map<InstanceID, integer>           // Sequence numbers
  executed: Set<InstanceID>               // Executed commands

Constants:
  N = total replicas
  F = floor((N-1)/2)                      // Fault tolerance
  FAST_QUORUM = F + floor((F+1)/2) + 1
  SLOW_QUORUM = F + 1

On ReceiveClientCommand(cmd):
  instance = [replica: self.id, index: ++cmdId]

  // Phase 1: Pre-accept
  initialDeps = FindConflicts(cmd)
  initialSeq = MaxSeq(initialDeps) + 1

  SendPreAccept(instance, cmd, initialDeps, initialSeq)

  // Wait for responses
  responses = WaitForPreAcceptAcks(FAST_QUORUM)

  if AllSameDeps(responses):
    // Fast path
    status[instance] = COMMITTED
    ReplyToClient(instance)
    BroadcastCommit(instance, cmd, initialDeps, initialSeq)
  else:
    // Slow path
    finalDeps = UnionAllDeps(responses)
    finalSeq = MaxSeq(finalDeps) + 1
    SendAccept(instance, cmd, finalDeps, finalSeq)
    WaitForAcceptAcks(SLOW_QUORUM)
    status[instance] = COMMITTED
    ReplyToClient(instance)
    BroadcastCommit(instance, cmd, finalDeps, finalSeq)

FindConflicts(cmd):
  conflicts = {}
  for instance, storedCmd in cmds:
    if status[instance] != NONE and KeysOverlap(cmd, storedCmd):
      if IsWrite(cmd) or IsWrite(storedCmd):
        conflicts.add(instance)
  return conflicts

On ReceivePreAccept(instance, cmd, theirDeps, theirSeq):
  myDeps = FindConflicts(cmd)

  // Union of dependencies
  combinedDeps = myDeps ∪ theirDeps
  seq[instance] = max(theirSeq, MaxSeq(combinedDeps) + 1)
  deps[instance] = combinedDeps
  cmds[instance] = cmd
  status[instance] = PREACCEPTED

  SendPreAcceptAck(instance, combinedDeps, seq[instance])

On ReceiveCommit(instance, cmd, finalDeps, finalSeq):
  cmds[instance] = cmd
  deps[instance] = finalDeps
  seq[instance] = finalSeq
  status[instance] = COMMITTED

  // Try to execute
  TryExecute()

TryExecute():
  // Build dependency graph of committed commands
  graph = BuildDependencyGraph()

  // Find strongly connected components (SCCs)
  sccs = TarjanSCC(graph)

  // Execute in topological order
  for scc in TopologicalSort(sccs):
    if scc.size() == 1:
      ExecuteSequential(scc[0])
    else:
      ExecuteSCC(scc)  // Within SCC, use deterministic ordering

def ExecuteSCC(scc):
  // Sort by sequence number, then by replica ID
  sorted = sorted(scc, key=lambda i: (seq[i], i.replica))
  for instance in sorted:
    if instance not in executed:
      ApplyToStateMachine(cmds[instance])
      executed.add(instance)
```

### 6.3 Flexible Paxos

```
Algorithm: Flexible Paxos with Asymmetric Quorums

State:
  ballot: integer
  phase1Quorum: Set<Replica>
  phase2Quorum: Set<Replica>
  preparedBallot: map<Replica, Ballot>
  acceptedValue: map<Replica, Value>

Constants:
  N = total replicas
  F = floor((N-1)/2)
  PHASE1_SIZE = N       // Large phase-1 quorum
  PHASE2_SIZE = F + 1   // Small phase-2 quorum

PreparePhase(b):
  ballot = b
  phase1Quorum = SelectQuorum(PHASE1_SIZE)

  for replica in phase1Quorum:
    SendPrepare(replica, b)

  responses = WaitForPromises(phase1Quorum)

  // Find highest accepted value
  maxBallot = 0
  valueToPropose = clientValue

  for response in responses:
    if response.acceptedBallot > maxBallot:
      maxBallot = response.acceptedBallot
      valueToPropose = response.acceptedValue

  return valueToPropose

AcceptPhase(b, v):
  phase2Quorum = SelectQuorum(PHASE2_SIZE)

  // phase2Quorum must intersect with phase1Quorum
  phase2Quorum = phase2Quorum ∪ (phase1Quorum ∩ AnyReplica)

  for replica in phase2Quorum:
    SendAccept(replica, b, v)

  responses = WaitForAccepted(phase2Quorum)
  return |responses| >= PHASE2_SIZE

SelectQuorum(size):
  // Strategy: Choose closest replicas for low latency
  return GetClosestReplicas(self.location, size)
```

---

## 7. Go Implementation

```go
// Package consensus provides implementations of consensus protocol variations
package consensus

import (
 "context"
 "fmt"
 "math/rand"
 "sort"
 "sync"
 "sync/atomic"
 "time"
)

// ============================================
// Multi-Paxos Implementation
// ============================================

// Command represents a client command
type Command struct {
 ID      uint64
 Key     string
 Value   []byte
 IsWrite bool
}

// Slot represents a log entry
type Slot struct {
 Index   uint64
 Ballot  uint64
 Command Command
 Committed bool
}

// MultiPaxos implements the Multi-Paxos protocol
type MultiPaxos struct {
 id       string
 replicas []string

 mu           sync.RWMutex
 ballot       uint64
 isLeader     bool
 log          map[uint64]*Slot
 nextIndex    map[string]uint64
 matchIndex   map[string]uint64
 commitIndex  uint64

 pendingCmds  chan Command
 proposeCh    chan Command

 batchSize    int
 pipelineSize int
}

// NewMultiPaxos creates a new Multi-Paxos instance
func NewMultiPaxos(id string, replicas []string) *MultiPaxos {
 return &MultiPaxos{
  id:           id,
  replicas:     replicas,
  log:          make(map[uint64]*Slot),
  nextIndex:    make(map[string]uint64),
  matchIndex:   make(map[string]uint64),
  pendingCmds:  make(chan Command, 1000),
  proposeCh:    make(chan Command, 100),
  batchSize:    100,
  pipelineSize: 1000,
 }
}

// Propose submits a command for consensus
func (mp *MultiPaxos) Propose(ctx context.Context, cmd Command) error {
 select {
 case mp.pendingCmds <- cmd:
  return nil
 case <-ctx.Done():
  return ctx.Err()
 }
}

// runLeader runs the leader logic
func (mp *MultiPaxos) runLeader() {
 ticker := time.NewTicker(10 * time.Millisecond)
 defer ticker.Stop()

 batch := make([]Command, 0, mp.batchSize)

 for {
  select {
  case cmd := <-mp.pendingCmds:
   batch = append(batch, cmd)
   if len(batch) >= mp.batchSize {
    mp.proposeBatch(batch)
    batch = batch[:0]
   }

  case <-ticker.C:
   if len(batch) > 0 {
    mp.proposeBatch(batch)
    batch = batch[:0]
   }
  }
 }
}

func (mp *MultiPaxos) proposeBatch(batch []Command) {
 mp.mu.Lock()
 defer mp.mu.Unlock()

 if !mp.isLeader {
  return
 }

 startSlot := mp.getNextFreeSlot()
 ballot := mp.ballot

 for i, cmd := range batch {
  slotIdx := startSlot + uint64(i)
  mp.log[slotIdx] = &Slot{
   Index:   slotIdx,
   Ballot:  ballot,
   Command: cmd,
  }

  // Send Accept messages to all replicas
  for _, replica := range mp.replicas {
   if replica != mp.id {
    go mp.sendAccept(replica, slotIdx, ballot, cmd)
   }
  }
 }
}

func (mp *MultiPaxos) sendAccept(replica string, slotIdx, ballot uint64, cmd Command) {
 // Simulated network call
 // In real implementation, use RPC
}

func (mp *MultiPaxos) getNextFreeSlot() uint64 {
 var maxSlot uint64
 for idx := range mp.log {
  if idx > maxSlot {
   maxSlot = idx
  }
 }
 return maxSlot + 1
}

// HandleAcceptAck processes an accept acknowledgment
func (mp *MultiPaxos) HandleAcceptAck(replica string, slotIdx uint64, ballot uint64) {
 mp.mu.Lock()
 defer mp.mu.Unlock()

 if ballot != mp.ballot {
  return // Stale response
 }

 mp.matchIndex[replica] = max(mp.matchIndex[replica], slotIdx)

 // Check if slot can be committed
 if mp.countMatching(slotIdx) > len(mp.replicas)/2 {
  mp.commit(slotIdx)
 }
}

func (mp *MultiPaxos) countMatching(slotIdx uint64) int {
 count := 1 // Count self
 for _, idx := range mp.matchIndex {
  if idx >= slotIdx {
   count++
  }
 }
 return count
}

func (mp *MultiPaxos) commit(slotIdx uint64) {
 if slot, ok := mp.log[slotIdx]; ok && !slot.Committed {
  slot.Committed = true
  // Apply to state machine
 }
}

// ============================================
// EPaxos Implementation
// ============================================

// InstanceID uniquely identifies a command instance
type InstanceID struct {
 Replica string
 Index   uint64
}

// EPaxosCommand extends Command with EPaxos metadata
type EPaxosCommand struct {
 Command
 Instance InstanceID
 Deps     map[InstanceID]struct{}
 Seq      uint64
 Status   CommandStatus
}

// CommandStatus represents the status of a command
type CommandStatus int

const (
 StatusNone CommandStatus = iota
 StatusPreAccepted
 StatusAccepted
 StatusCommitted
)

// EPaxos implements the Egalitarian Paxos protocol
type EPaxos struct {
 id       string
 replicas []string

 mu        sync.RWMutex
 cmdID     uint64
 cmds      map[InstanceID]*EPaxosCommand
 executed  map[InstanceID]bool

 fastQuorumSize int
 slowQuorumSize int
}

// NewEPaxos creates a new EPaxos instance
func NewEPaxos(id string, replicas []string) *EPaxos {
 n := len(replicas) + 1 // Include self
 f := (n - 1) / 2

 return &EPaxos{
  id:             id,
  replicas:       replicas,
  cmds:           make(map[InstanceID]*EPaxosCommand),
  executed:       make(map[InstanceID]bool),
  fastQuorumSize: f + (f+1)/2 + 1,
  slowQuorumSize: f + 1,
 }
}

// Propose proposes a command using EPaxos
func (ep *EPaxos) Propose(ctx context.Context, cmd Command) (*EPaxosCommand, error) {
 ep.mu.Lock()
 ep.cmdID++
 instance := InstanceID{Replica: ep.id, Index: ep.cmdID}

 initialDeps := ep.findConflicts(cmd)
 initialSeq := ep.maxSeq(initialDeps) + 1

 epCmd := &EPaxosCommand{
  Command:  cmd,
  Instance: instance,
  Deps:     initialDeps,
  Seq:      initialSeq,
  Status:   StatusNone,
 }
 ep.cmds[instance] = epCmd
 ep.mu.Unlock()

 // Phase 1: Pre-accept
 responses := ep.sendPreAccept(instance, cmd, initialDeps, initialSeq)

 if ep.allSameDeps(responses) {
  // Fast path
  ep.mu.Lock()
  epCmd.Status = StatusCommitted
  ep.mu.Unlock()
  ep.broadcastCommit(epCmd)
 } else {
  // Slow path
  finalDeps := ep.unionAllDeps(responses)
  finalSeq := ep.maxSeq(finalDeps) + 1

  ep.sendAccept(instance, cmd, finalDeps, finalSeq)

  ep.mu.Lock()
  epCmd.Deps = finalDeps
  epCmd.Seq = finalSeq
  epCmd.Status = StatusCommitted
  ep.mu.Unlock()
  ep.broadcastCommit(epCmd)
 }

 return epCmd, nil
}

func (ep *EPaxos) findConflicts(cmd Command) map[InstanceID]struct{} {
 conflicts := make(map[InstanceID]struct{})

 for instance, stored := range ep.cmds {
  if stored.Status == StatusNone {
   continue
  }
  if stored.Key == cmd.Key && (cmd.IsWrite || stored.IsWrite) {
   conflicts[instance] = struct{}{}
  }
 }

 return conflicts
}

func (ep *EPaxos) maxSeq(deps map[InstanceID]struct{}) uint64 {
 var maxSeq uint64
 for instance := range deps {
  if cmd, ok := ep.cmds[instance]; ok && cmd.Seq > maxSeq {
   maxSeq = cmd.Seq
  }
 }
 return maxSeq
}

func (ep *EPaxos) allSameDeps(responses []PreAcceptResponse) bool {
 if len(responses) == 0 {
  return false
 }
 first := responses[0].Deps
 for _, resp := range responses[1:] {
  if !depsEqual(first, resp.Deps) {
   return false
  }
 }
 return true
}

func depsEqual(a, b map[InstanceID]struct{}) bool {
 if len(a) != len(b) {
  return false
 }
 for k := range a {
  if _, ok := b[k]; !ok {
   return false
  }
 }
 return true
}

func (ep *EPaxos) unionAllDeps(responses []PreAcceptResponse) map[InstanceID]struct{} {
 union := make(map[InstanceID]struct{})
 for _, resp := range responses {
  for dep := range resp.Deps {
   union[dep] = struct{}{}
  }
 }
 return union
}

// PreAcceptResponse represents a response to PreAccept
type PreAcceptResponse struct {
 Replica string
 Deps    map[InstanceID]struct{}
 Seq     uint64
}

func (ep *EPaxos) sendPreAccept(instance InstanceID, cmd Command, deps map[InstanceID]struct{}, seq uint64) []PreAcceptResponse {
 // Simulated - in real implementation, send RPCs
 return nil
}

func (ep *EPaxos) sendAccept(instance InstanceID, cmd Command, deps map[InstanceID]struct{}, seq uint64) {
 // Simulated - in real implementation, send RPCs
}

func (ep *EPaxos) broadcastCommit(cmd *EPaxosCommand) {
 // Simulated - in real implementation, broadcast to all
}

// TryExecute attempts to execute committed commands
func (ep *EPaxos) TryExecute() {
 ep.mu.Lock()
 defer ep.mu.Unlock()

 // Build dependency graph
 graph := make(map[InstanceID]map[InstanceID]struct{})
 for instance, cmd := range ep.cmds {
  if cmd.Status == StatusCommitted {
   graph[instance] = cmd.Deps
  }
 }

 // Topological sort with SCC handling
 visited := make(map[InstanceID]bool)
 var order []InstanceID

 var visit func(InstanceID)
 visit = func(node InstanceID) {
  if visited[node] {
   return
  }
  visited[node] = true

  for dep := range graph[node] {
   visit(dep)
  }

  if !ep.executed[node] {
   order = append(order, node)
  }
 }

 for instance := range graph {
  visit(instance)
 }

 // Execute in order
 for _, instance := range order {
  if cmd, ok := ep.cmds[instance]; ok {
   ep.applyToStateMachine(cmd)
   ep.executed[instance] = true
  }
 }
}

func (ep *EPaxos) applyToStateMachine(cmd *EPaxosCommand) {
 // Apply command to state machine
 fmt.Printf("Executing command %d on key %s\n", cmd.ID, cmd.Key)
}

// ============================================
// Flexible Paxos Implementation
// ============================================

// FlexiblePaxos implements Paxos with asymmetric quorums
type FlexiblePaxos struct {
 id       string
 replicas []string

 mu              sync.RWMutex
 ballot          uint64
 phase1Quorum    []string
 phase2Quorum    []string
 promises        map[string]*Promise
 accepted        map[string]*AcceptedValue

 phase1Size      int
 phase2Size      int
}

// Promise represents a phase 1b response
type Promise struct {
 Replica       string
 Ballot        uint64
 AcceptedBallot uint64
 AcceptedValue  []byte
}

// AcceptedValue represents an accepted value
type AcceptedValue struct {
 Ballot uint64
 Value  []byte
}

// NewFlexiblePaxos creates a new Flexible Paxos instance
func NewFlexiblePaxos(id string, replicas []string) *FlexiblePaxos {
 n := len(replicas) + 1
 f := (n - 1) / 2

 return &FlexiblePaxos{
  id:           id,
  replicas:     replicas,
  promises:     make(map[string]*Promise),
  accepted:     make(map[string]*AcceptedValue),
  phase1Size:   n,       // Large phase-1
  phase2Size:   f + 1,   // Small phase-2
 }
}

// Propose runs the Flexible Paxos protocol
func (fp *FlexiblePaxos) Propose(ctx context.Context, value []byte) ([]byte, error) {
 // Phase 1: Prepare
 ballot := fp.incrementBallot()
 fp.selectPhase1Quorum()

 promises := fp.sendPrepare(ballot)

 // Find highest accepted value
 var maxBallot uint64
 valueToPropose := value

 for _, promise := range promises {
  if promise.AcceptedBallot > maxBallot {
   maxBallot = promise.AcceptedBallot
   valueToPropose = promise.AcceptedValue
  }
 }

 // Phase 2: Accept
 fp.selectPhase2Quorum()
 accepted := fp.sendAccept(ballot, valueToPropose)

 if accepted {
  return valueToPropose, nil
 }

 return nil, fmt.Errorf("failed to achieve consensus")
}

func (fp *FlexiblePaxos) incrementBallot() uint64 {
 fp.mu.Lock()
 defer fp.mu.Unlock()
 fp.ballot++
 return fp.ballot
}

func (fp *FlexiblePaxos) selectPhase1Quorum() {
 fp.mu.Lock()
 defer fp.mu.Unlock()
 fp.phase1Quorum = fp.selectClosest(fp.phase1Size)
}

func (fp *FlexiblePaxos) selectPhase2Quorum() {
 fp.mu.Lock()
 defer fp.mu.Unlock()

 // Phase-2 quorum must intersect with phase-1
 fp.phase2Quorum = fp.selectClosest(fp.phase2Size)

 // Ensure intersection
 hasIntersection := false
 for _, p1 := range fp.phase1Quorum {
  for _, p2 := range fp.phase2Quorum {
   if p1 == p2 {
    hasIntersection = true
    break
   }
  }
 }

 if !hasIntersection {
  // Add a phase-1 member to phase-2
  fp.phase2Quorum = append(fp.phase2Quorum, fp.phase1Quorum[0])
 }
}

func (fp *FlexiblePaxos) selectClosest(k int) []string {
 // In real implementation, select based on network distance
 all := append([]string{fp.id}, fp.replicas...)
 if k >= len(all) {
  return all
 }
 return all[:k]
}

func (fp *FlexiblePaxos) sendPrepare(ballot uint64) []*Promise {
 // Simulated - in real implementation, send RPCs
 return nil
}

func (fp *FlexiblePaxos) sendAccept(ballot uint64, value []byte) bool {
 // Simulated - in real implementation, send RPCs
 return true
}

func max(a, b uint64) uint64 {
 if a > b {
  return a
 }
 return b
}

// Replica represents a consensus participant
type Replica struct {
 ID       string
 Location string
 Latency  time.Duration
}

// ConsensusFactory creates appropriate consensus instances
type ConsensusFactory struct {
 config Config
}

// Config holds consensus configuration
type Config struct {
 Protocol       string // "multipaxos", "epaxos", "flexible"
 Replicas       []string
 BatchSize      int
 PipelineDepth  int
}

// Create creates a consensus instance based on config
func (f *ConsensusFactory) Create(id string) (interface{}, error) {
 switch f.config.Protocol {
 case "multipaxos":
  return NewMultiPaxos(id, f.config.Replicas), nil
 case "epaxos":
  return NewEPaxos(id, f.config.Replicas), nil
 case "flexible":
  return NewFlexiblePaxos(id, f.config.Replicas), nil
 default:
  return nil, fmt.Errorf("unknown protocol: %s", f.config.Protocol)
 }
}
```

---

## 8. Comparison with Related Approaches

### 8.1 Protocol Comparison Matrix

| Property | Classic Paxos | Multi-Paxos | EPaxos | Flexible Paxos |
|----------|---------------|-------------|--------|----------------|
| **Leader** | Single leader | Stable leader | Leaderless | Configurable |
| **Latency (uncontested)** | 2 RTT | 2 RTT (1 RTT with leader) | 1 RTT (fast path) | 1-2 RTT |
| **Latency (contested)** | 2 RTT | 2 RTT | 2 RTT (slow path) | 1-2 RTT |
| **Throughput** | Low | High | Very High | High |
| **Message Complexity** | O(n²) | O(n) per op | O(n) per op | O(n) - O(2n) |
| **Geo-replication** | Poor | Moderate | Excellent | Configurable |
| **Conflict Handling** | Sequential | Sequential | Parallel | Sequential |
| **Implementation Complexity** | High | Medium | High | Medium |

### 8.2 Quorum Comparison

| Protocol | Phase-1 Quorum | Phase-2 Quorum | Intersection |
|----------|----------------|----------------|--------------|
| Classic Paxos | ⌈n/2⌉ + 1 | ⌈n/2⌉ + 1 | Required |
| Multi-Paxos | ⌈n/2⌉ + 1 | ⌈n/2⌉ + 1 | Required |
| EPaxos (fast) | ⌈3n/4⌉ | N/A | Implicit |
| EPaxos (slow) | ⌈n/2⌉ + 1 | ⌈n/2⌉ + 1 | Required |
| Flexible Paxos | n - f | f + 1 | |Q₁| + |Q₂| > n |

### 8.3 Use Case Recommendations

| Scenario | Recommended Protocol | Rationale |
|----------|---------------------|-----------|
| Single datacenter | Multi-Paxos | Simple, high throughput |
| Geo-distributed | EPaxos | Minimizes WAN RTT |
| Read-heavy | Flexible Paxos | Small phase-2 quorums |
| Write-heavy | EPaxos | Parallel commit paths |
| Strong consistency | Multi-Paxos | Well-understood |
| Latency-sensitive | EPaxos (fast path) | 1 RTT commit |

---

## 9. Visual Representations

### 9.1 Consensus Protocol Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONSENSUS PROTOCOL TAXONOMY                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Consensus Protocols                                                          │
│  │                                                                            │
│  ├── Single-Decree                                                            │
│  │   └── Classic Paxos                                                        │
│  │       ├── Basic Paxos (single value)                                       │
│  │       └── Multi-Paxos (log replication)                                    │
│  │           ├── Standard Multi-Paxos                                        │
│  │           └── Optimized Variants                                           │
│  │               ├── Batched Multi-Paxos                                      │
│  │               ├── Pipelined Multi-Paxos                                    │
│  │               └── Flexible Paxos ←── Asymmetric quorums                    │
│  │                                                                            │
│  └── Multi-Decree                                                             │
│      ├── Leader-Based                                                         │
│      │   └── Multi-Paxos derivatives                                          │
│      ├── Leaderless                                                           │
│      │   ├── EPaxos ←── Conflict-based parallelization                        │
│      │   ├── Mencius ←── Rotating leaders                                     │
│      │   └── S-Paxos ←── Separated leader election                            │
│      └── Speculative                                                          │
│          ├── Speculative Paxos                                                │
│          └── NOPaxos                                                          │
│                                                                              │
│  Key Innovation Dimensions:                                                   │
│  • Quorum Structure: Symmetric → Asymmetric                                   │
│  • Leadership: Single → Rotating → Leaderless                                 │
│  • Parallelism: Sequential → Conflict-based → Full Parallel                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 EPaxos Fast vs Slow Path

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    EPAXOS EXECUTION PATHS                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  FAST PATH (1 RTT) ────────────────────────────────────────────────          │
│                                                                              │
│  Client        Replica A (Command Leader)      Fast Quorum (⌈3n/4⌉)          │
│    │                      │                            │                     │
│    │─── Command C ───────>│                            │                     │
│    │                      │                            │                     │
│    │                      │──── PreAccept(C, Deps) ───>│                     │
│    │                      │                            │                     │
│    │                      │<───── PreAcceptOK ────────│                     │
│    │                      │   (All same dependencies)  │                     │
│    │                      │                            │                     │
│    │<── Committed ───────│                            │                     │
│    │                      │                            │                     │
│    [Total: 1 RTT to closest replica]                                         │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────           │
│                                                                              │
│  SLOW PATH (2 RTT) ────────────────────────────────────────────────          │
│                                                                              │
│  Client        Replica A                      Fast Quorum    Slow Quorum    │
│    │                      │                            │            │        │
│    │─── Command C ───────>│                            │            │        │
│    │                      │                            │            │        │
│    │                      │──── PreAccept(C) ─────────>│            │        │
│    │                      │                            │            │        │
│    │                      │<── PreAccept (diff deps) ──│            │        │
│    │                      │   (Conflicts detected)     │            │        │
│    │                      │                            │            │        │
│    │                      │──── Accept(C, UnionDeps) ───────────────>│        │
│    │                      │                            │            │        │
│    │                      │<──────── Accepted ──────────────────────│        │
│    │                      │                                        │        │
│    │<── Committed ───────│                                        │        │
│    │                      │                                        │        │
│    [Total: 2 RTT, same as Paxos]                                            │
│                                                                              │
│  Slow Path Trigger Conditions:                                               │
│  • Conflicting commands in flight                                            │
│  • Fast quorum members report different dependencies                          │
│  • Network partitions affecting fast quorum                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.3 Flexible Quorum Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FLEXIBLE PAXOS QUORUM STRUCTURE                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Classic Paxos (Symmetric):                                                  │
│                                                                              │
│  Phase 1: ─────[R1]────[R2]────[R3]────[R4]────[R5]────                     │
│                └──────────┬──────────┘                                       │
│                       Quorum (3)                                             │
│                                                                              │
│  Phase 2: ─────[R1]────[R2]────[R3]────[R4]────[R5]────                     │
│                     └──────────┬──────────┘                                  │
│                            Quorum (3)                                        │
│                                                                              │
│  Intersection: R2, R3, R4 all in both quorums (guaranteed)                   │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────           │
│                                                                              │
│  Flexible Paxos (Asymmetric, n=5, f=2):                                      │
│                                                                              │
│  Phase 1 (Large): ─────[R1]────[R2]────[R3]────[R4]────[R5]────             │
│                      └──────────────────┬──────────────────┘                 │
│                                   |Q₁| = 5 = n                               │
│                                                                              │
│  Phase 2 (Small): ─────[R1]────[R2]────[R3]────[R4]────[R5]────             │
│                                    └────┬────┘                               │
│                                     |Q₂| = 3 = f+1                           │
│                                                                              │
│  Intersection Proof: |Q₁| + |Q₂| = 5 + 3 = 8 > n = 5                         │
│                      ∴ Q₁ ∩ Q₂ ≠ ∅ (at least 3 nodes)                        │
│                                                                              │
│  Benefits:                                                                   │
│  • Phase 1: Contact all for maximum safety (one-time cost)                   │
│  • Phase 2: Contact minimum for fast commits (repeated)                      │
│  • Recovery: Smaller phase-2 quorum = faster failover                        │
│                                                                              │
│  Geographic Optimization Example:                                            │
│                                                                              │
│   [US-West]    [US-East]    [EU]    [Asia]    [AU]                          │
│      R1          R2         R3       R4       R5                             │
│       │          │          │        │        │                              │
│       └──────────┴──────────┘        │        │                              │
│              Phase 1 (local + backup)                                        │
│                                                                              │
│       │          │                    │                                     │
│       └──────────┘                    │                                      │
│              Phase 2 (local only) ────┘                                      │
│                                                                              │
│  [R1, R2 contacted for Phase 1, only R1 for Phase 2 → lower latency]        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 10. Academic References

### 10.1 Foundational Papers

1. **Lamport, L. (1998)**. "The Part-Time Parliament". *ACM Transactions on Computer Systems*, 16(2), 133-169.
   - Original Paxos protocol

2. **Lamport, L. (2001)**. "Paxos Made Simple". *ACM SIGACT News*, 32(4), 18-25.
   - Simplified explanation of Paxos

3. **Chandra, T. D., Griesemer, R., & Redstone, J. (2007)**. "Paxos Made Live: An Engineering Perspective". *PODC*.
   - Practical implementation experiences

### 10.2 Multi-Paxos and Optimizations

1. **Van Renesse, R., & Altinbuken, D. (2015)**. "Paxos Made Moderately Complex". *ACM Computing Surveys*, 47(3), 1-36.
   - Comprehensive survey of Paxos variants

2. **Kirsch, J., & Amir, Y. (2008)**. "Paxos for System Builders". *Johns Hopkins University Tech Report*.
   - Practical implementation guide

### 10.3 EPaxos and Leaderless Consensus

1. **Moraru, I., Andersen, D. G., & Kaminsky, M. (2013)**. "There Is More Consensus in Egalitarian Parliaments". *SOSP*.
   - EPaxos: Leaderless consensus with fast path

2. **Moraru, I. (2014)**. "EPaxos: A New Approach to Efficient, Scalable, and Robust Consensus". *PhD Thesis, CMU*.
   - Extended EPaxos analysis

3. **Poke, M., & Schiavoni, V. (2016)**. "Dare: High-Performance State Machine Replication on RDMA". *DSN*.
   - High-performance leaderless replication

### 10.4 Flexible Paxos

1. **Howard, H., Malkhi, D., & Schwarzmann, A. A. (2016)**. "Flexible Paxos: Quorum Intersection Revisited". *OPODIS*.
   - Asymmetric quorum systems

2. **Howard, H., & Malkhi, D. (2019)**. "Paxos vs Raft: Have We Reached Consensus on Distributed Consensus?". *PaPoC*.
    - Comparative analysis of consensus protocols

### 10.5 Recent Advances (2020-2025)

1. **Pires, R., et al. (2025)**. "Can 1000 Nodes Agree? Scalability of Consensus Protocols". *NSDI*.
    - Large-scale consensus analysis

2. **Charapko, A., et al. (2024)**. "SEDA: A Scalable, Efficient, and Dynamic Architecture for Consensus". *EuroSys*.
    - Modern consensus architecture

---

## 11. Summary and Key Insights

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONSENSUS VARIATIONS: KEY INSIGHTS                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Quorum Flexibility                                                       │
│     • Symmetric quorums: Simple but potentially inefficient                   │
│     • Asymmetric quorums: Optimize for different phases                       │
│     • Intersection property: The fundamental safety requirement               │
│                                                                              │
│  2. Leadership Models                                                         │
│     • Stable leader: High throughput, single point of contention              │
│     • Rotating leader: Better load distribution                               │
│     • Leaderless: Maximum parallelism, conflict resolution needed             │
│                                                                              │
│  3. Latency-Throughput Tradeoffs                                              │
│     • Multi-Paxos: 2 RTT, O(N) throughput                                     │
│     • EPaxos fast: 1 RTT, O(N) throughput, conflict-dependent                 │
│     • EPaxos slow: 2 RTT, falls back to Paxos-like behavior                   │
│                                                                              │
│  4. Geographic Considerations                                                 │
│     • WAN latency dominates: Optimize for closest quorum                      │
│     • Flexible Paxos: Separate phase-1/phase-2 quorum selection               │
│     • EPaxos: Natural affinity for geo-distribution                           │
│                                                                              │
│  5. Implementation Guidance                                                   │
│     • Start with Multi-Paxos for simplicity                                   │
│     • Add batching and pipelining for throughput                              │
│     • Consider EPaxos for multi-region deployments                            │
│     • Use Flexible Paxos for specialized latency requirements                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
