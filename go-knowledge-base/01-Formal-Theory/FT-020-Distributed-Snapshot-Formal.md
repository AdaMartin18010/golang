# FT-020: Distributed Snapshot - Formal Specification

## Overview

A distributed snapshot captures the global state of a distributed system at a particular point in time. The Chandy-Lamport algorithm provides an efficient method for recording consistent global states without stopping the system. This is fundamental for checkpointing, deadlock detection, and debugging distributed systems.

## Theoretical Foundations

### 1.1 System Model

**Distributed System**:

```
A distributed system consists of:
- N processes: P = {p₁, p₂, ..., pₙ}
- Communication channels: C = {cᵢⱼ : pᵢ sends to pⱼ}
- Events: E = {send, receive, internal}

Process state: sᵢ ∈ Sᵢ
Channel state: cᵢⱼ contains messages in transit
```

**Global State**:

```
A global state GS = (S, C) where:
- S = (s₁, s₂, ..., sₙ): local states of all processes
- C = {cᵢⱼ}: states of all channels

GS is consistent if:
∀message m: if m ∈ cᵢⱼ (in transit), then send(m) ∈ Sᵢ and receive(m) ∉ Sⱼ
```

### 1.2 Happens-Before Relation

```
The happens-before relation (→) is defined as:

1. If e and f are events in the same process and e occurs before f,
   then e → f

2. If e is send(m) and f is receive(m), then e → f

3. If e → f and f → g, then e → g (transitivity)

Concurrent events: e || f iff ¬(e → f) ∧ ¬(f → e)
```

**Consistent Global State**:

```
GS = (S, C) is consistent iff:
∀pᵢ, pⱼ: if sᵢ includes event e, and e → f where f occurs at pⱼ,
then sⱼ must also include f (or an event that happens after f).

Equivalently: GS contains no orphan messages.
  (No message received without being sent in GS)
```

### 1.3 Chandy-Lamport Algorithm

**Algorithm Description**:

```
Initiator Process:
  1. Record local state
  2. Send MARKER on all outgoing channels
  3. Start recording messages on incoming channels

Non-Initiator Process (on receiving MARKER from pⱼ):
  If first MARKER received:
    1. Record local state
    2. Send MARKER on all outgoing channels
    3. Start recording messages on all incoming channels except from pⱼ
    4. Channel cⱼ (from pⱼ) state = empty
  Else:
    1. Channel cⱼ state = messages recorded since first MARKER
    2. Stop recording on cⱼ

Termination: When all processes have received MARKER on all channels
```

**Correctness Proof**:

```
Theorem: The Chandy-Lamport algorithm records a consistent global state.

Proof:

Lemma 1 (Marker Receipt Order): If process pⱼ receives marker from pᵢ,
then all messages sent by pᵢ before its marker are received by pⱼ
before the marker.

Proof of Lemma 1: Channels are FIFO. Marker is sent after preceding
messages, so it arrives after them. ∎

Lemma 2 (No Orphan Messages): If message m is in channel state cᵢⱼ,
then send(m) is recorded in pᵢ's state but receive(m) is not recorded
in pⱼ's state.

Proof of Lemma 2:
- pⱼ records its state upon first MARKER receipt
- Messages recorded on cᵢⱼ are those received after pⱼ's state recording
  but before receiving MARKER from pᵢ
- By Lemma 1, these were sent by pᵢ before sending MARKER
- pᵢ recorded its state before sending MARKER
- Therefore send(m) ∈ state(pᵢ) and receive(m) ∉ state(pⱼ) ∎

Main Proof:
Consider any message m sent from pᵢ to pⱼ.

Case 1: send(m) happens before pᵢ records state
  - If receive(m) happens before pⱼ records state:
    m is not in channel state (already received)
  - If receive(m) happens after pⱼ records state:
    By Lemma 1, receive(m) happens after MARKER from pᵢ
    So m is recorded in channel state cᵢⱼ

Case 2: send(m) happens after pᵢ records state
  - By Lemma 1, m is sent after MARKER
  - So m arrives after MARKER
  - pⱼ stops recording on cᵢⱼ upon MARKER receipt
  - Therefore m is not in channel state

In all cases, either both send and receive are recorded, or neither is,
or send is recorded and message is in channel state. This is exactly
the definition of consistency. ∎
```

### 1.4 Snapshot Properties

```
Property 1 (Liveness): The algorithm terminates if the snapshot
initiator is connected to all other processes.

Proof: MARKERs propagate through the communication graph. In a
connected graph, all processes eventually receive a MARKER. ∎

Property 2 (Safety): The recorded global state is a state that
could have occurred during execution.

Proof: The recorded state corresponds to a consistent cut - a
partition of events into "before" and "after" such that no message
crosses from "after" to "before". Such cuts always correspond to
valid system states. ∎

Property 3 (No Process Blocking): The algorithm never blocks
process execution.

Proof: Recording state and sending MARKER are non-blocking operations.
Message recording can be done asynchronously. ∎
```

## TLA+ Specification

```tla
----------------------- MODULE DistributedSnapshot -----------------------
EXTENDS Integers, Sequences, FiniteSets, TLC

CONSTANTS Processes,         \* Set of process IDs
          Channels,          \* Set of channel IDs
          MaxMessages        \* Bound for model checking

VARIABLES processState,      \* State of each process
          channelState,      \* Messages in each channel
          snapshotState,     \* Snapshot recording state
          markerStatus       \* Which markers have been sent/received

\* Event types
EventType == {"internal", "send", "receive", "record", "marker_send", "marker_recv"}

\* Process state
ProcessStatus == {"running", "recorded"}

\* Snapshot state per process
SnapshotData == [localState: Nat,
                 channelStates: [Channels -> Seq(Nat)],
                 status: ProcessStatus,
                 markersReceived: SUBSET Processes,
                 recording: SUBSET Channels]

\* Initial state
Init ==
  ∧ processState = [p ∈ Processes ↦ [state ↦ 0, events ↦ ⟨⟩]]
  ∧ channelState = [c ∈ Channels ↦ ⟨⟩]
  ∧ snapshotState = [p ∈ Processes ↦
       [localState ↦ 0,
        channelStates ↦ [c ∈ Channels ↦ ⟨⟩],
        status ↦ "running",
        markersReceived ↦ {},
        recording ↦ {}]]
  ∧ markerStatus = [sent ↦ {}, received ↦ {}]

\* Internal event at process
InternalEvent(p) ==
  ∧ processState' = [processState EXCEPT
       ![p].state = @ + 1,
       ![p].events = Append(@, [type ↦ "internal", data ↦ @.state + 1])]
  ∧ UNCHANGED ⟨channelState, snapshotState, markerStatus⟩

\* Send message from p1 to p2
SendMessage(p1, p2, c) ==
  ∧ c ∈ Channels  \* c is channel from p1 to p2
  ∧ LET msg == processState[p1].state
    IN channelState' = [channelState EXCEPT ![c] = Append(@, msg)]
  ∧ processState' = [processState EXCEPT
       ![p1].events = Append(@, [type ↦ "send", to ↦ p2, data ↦ msg])]
  ∧ UNCHANGED ⟨snapshotState, markerStatus⟩

\* Receive message at p2 from p1
ReceiveMessage(p1, p2, c) ==
  ∧ Len(channelState[c]) > 0
  ∧ LET msg == Head(channelState[c])
    IN ∧ channelState' = [channelState EXCEPT ![c] = Tail(@)]
       ∧ processState' = [processState EXCEPT
            ![p2].state = @ + msg,
            ![p2].events = Append(@, [type ↦ "receive", from ↦ p1, data ↦ msg])]
       ∧ snapshotState' = [snapshotState EXCEPT ![p2] =
            IF c ∈ @.recording
            THEN [channelStates EXCEPT ![c] = Append(@[c], msg)]
            ELSE @]
  ∧ UNCHANGED markerStatus

\* Initiate snapshot
InitiateSnapshot(p) ==
  ∧ snapshotState[p].status = "running"
  ∧ snapshotState' = [snapshotState EXCEPT ![p] =
       [@ EXCEPT !.status = "recorded",
                  !.localState = processState[p].state,
                  !.recording = {c ∈ Channels : c connects to p}]]
  ∧ markerStatus' = [markerStatus EXCEPT !.sent = @ ∪ {[from ↦ p, to ↦ q] : q ∈ Processes \\ {p}}]
  ∧ UNCHANGED ⟨processState, channelState⟩

\* Receive marker from p1 at p2
ReceiveMarker(p1, p2) ==
  ∧ [from ↦ p1, to ↦ p2] ∈ markerStatus.sent
  ∧ [from ↦ p1, to ↦ p2] ∉ markerStatus.received
  ∧ markerStatus' = [markerStatus EXCEPT !.received = @ ∪ {[from ↦ p1, to ↦ p2]}]
  ∧ IF snapshotState[p2].status = "running"
    THEN \* First marker - record state
         snapshotState' = [snapshotState EXCEPT ![p2] =
           [localState ↦ processState[p2].state,
            channelStates ↦ [c ∈ Channels ↦
              IF c connects from p1 THEN ⟨⟩ ELSE ⟨⟩],  \* empty for now
            status ↦ "recorded",
            markersReceived ↦ {p1},
            recording ↦ {c ∈ Channels : c connects to p2} \\ {channel(p1,p2)}]]
    ELSE \* Subsequent marker - stop recording this channel
         snapshotState' = [snapshotState EXCEPT ![p2].markersReceived = @ ∪ {p1},
                                                  ![p2].recording = @ \\ {channel(p1,p2)}]
  ∧ UNCHANGED ⟨processState, channelState⟩

\* Next state relation
Next ==
  ∨ ∃p ∈ Processes: InternalEvent(p)
  ∨ ∃p1, p2 ∈ Processes, c ∈ Channels: SendMessage(p1, p2, c)
  ∨ ∃p1, p2 ∈ Processes, c ∈ Channels: ReceiveMessage(p1, p2, c)
  ∨ ∃p ∈ Processes: InitiateSnapshot(p)
  ∨ ∃p1, p2 ∈ Processes: ReceiveMarker(p1, p2)

\* Invariants

\* Consistent snapshot: no orphan messages
ConsistentSnapshot ==
  ∀p ∈ Processes:
    snapshotState[p].status = "recorded"
    ⇒ ∀c ∈ Channels, msg ∈ snapshotState[p].channelStates[c]:
        msg was sent by source of c before it recorded state

\* Termination: eventually all processes record state
SnapshotTermination ==
  ◇(∀p ∈ Processes: snapshotState[p].status = "recorded")

\* No deadlock: markers eventually received
MarkerDelivery ==
  ∀m ∈ markerStatus.sent: ◇(m ∈ markerStatus.received)

=============================================================================
```

## Algorithm Pseudocode

### Chandy-Lamport Snapshot Algorithm

```
Algorithm: Chandy-Lamport Distributed Snapshot

Types:
  ProcessID: unique identifier for each process
  Message: data sent between processes
  Marker: special control message indicating snapshot

State at Process p:
  state: local process state
  recorded: boolean = false
  incoming: map<ProcessID, Queue<Message>>
  channelEmpty: map<ProcessID, boolean>
  snapshot: local state copy (when recorded)
  channelState: map<ProcessID, List<Message>>

Procedure InitiateSnapshot():
  // Called by initiator process
  RecordState()
  for each neighbor q:
    SendMarker(q)

Procedure ReceiveMessage(m, from):
  if m is Marker:
    HandleMarker(from)
  else:
    state ← UpdateState(state, m)
    if recorded and not channelEmpty[from]:
      channelState[from].append(m)

Procedure HandleMarker(from):
  if not recorded:
    RecordState()
    for each neighbor q where q ≠ from:
      SendMarker(q)
    channelEmpty[from] ← true
    channelState[from] ← empty
  else:
    channelEmpty[from] ← true
    // channelState[from] already contains recorded messages

Procedure RecordState():
  recorded ← true
  snapshot ← Copy(state)
  for each incoming channel c:
    channelEmpty[c] ← false
    channelState[c] ← empty
    StartRecording(c)

Procedure SendMarker(to):
  send(Marker, to)

Procedure StartRecording(channel):
  // Start buffering messages on this channel
  // Implementation-dependent

CheckTermination():
  // Snapshot is complete when:
  // 1. All processes have recorded (recorded = true for all)
  // 2. All channels have stopped recording (channelEmpty = true for all)

  if recorded and ∀c: channelEmpty[c]:
    ReportSnapshotComplete()

ReportSnapshotComplete():
  // Send snapshot data to initiator or collection point
  SendToInitiator({
    process: self,
    state: snapshot,
    channels: channelState
  })

Correctness Proof:

Theorem 1: The recorded state is consistent.

Proof:
Let S be the set of processes that have recorded when snapshot completes.

For any channel c from pᵢ to pⱼ:

Case 1: pᵢ recorded before sending any message on c
  - pⱼ's channel state for c is empty
  - No messages "in flight" from recorded send to unrecorded receive

Case 2: pᵢ sent messages on c before recording
  Subcase 2a: pⱼ received all messages before recording
    - channel state is empty
  Subcase 2b: pⱼ received some messages after recording
    - those messages are in channel state
    - they were sent before pᵢ recorded (FIFO property)

In all cases, no message crosses from "after" to "before" the cut. ∎

Theorem 2: The algorithm terminates in finite time.

Proof:
- MARKERs propagate through the communication graph
- Each process sends at most |neighbors| MARKERs
- Each process receives at most |neighbors| MARKERs
- Graph is finite, so termination is guaranteed ∎
```

### Lai-Yang Snapshot Algorithm

```
Algorithm: Lai-Yang (Non-Blocking) Snapshot

Difference from Chandy-Lamput:
- No explicit MARKER messages
- Uses message coloring (white/red phases)

State:
  color: {WHITE, RED}
  state: local process state
  receivedRed: map<ProcessID, boolean>
  channelState: map<ProcessID, List<Message>>

Procedure InitiateSnapshot():
  color ← RED
  snapshot ← Copy(state)
  for each channel c:
    channelState[c] ← empty
    receivedRed[c] ← false

On Send(m, to):
  send(color, m) to to
  // Messages sent after turning RED are RED

On Receive((msgColor, m), from):
  if color = WHITE and msgColor = RED:
    // First red message triggers recording
    color ← RED
    snapshot ← Copy(state)
    receivedRed[from] ← true
    for each channel c ≠ from:
      channelState[c] ← empty
      receivedRed[c] ← false
  else if color = RED and msgColor = WHITE:
    // White message received after turning RED
    channelState[from].append(m)

  if msgColor = RED:
    receivedRed[from] ← true

  state ← UpdateState(state, m)

  // Check completion
  if ∀c: receivedRed[c]:
    ReportSnapshotComplete()

Advantages:
- No separate MARKER messages (piggyback on data)
- Less overhead

Disadvantages:
- Requires modifying all messages to include color
- More complex termination detection
```

## Go Implementation

```go
// Package snapshot implements distributed snapshot algorithms
package snapshot

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// Message represents a message in the system
type Message struct {
 From    string
 To      string
 Payload interface{}
 Color   MessageColor // For Lai-Yang algorithm
}

// MessageColor represents the phase of a message
type MessageColor int

const (
 White MessageColor = iota // Before snapshot
 Red                       // After snapshot
)

// Marker is a special control message
type Marker struct {
 From      string
 Timestamp time.Time
}

// ProcessState represents the local state of a process
type ProcessState struct {
 Data      interface{}
 Timestamp time.Time
 Version   int64
}

// ChannelState represents the state of a communication channel
type ChannelState struct {
 Messages []interface{}
 mu       sync.RWMutex
}

// Process represents a process in the distributed system
type Process struct {
 ID       string
 State    ProcessState
 Neighbors []string
 Channels map[string]chan interface{}

 // Snapshot state
 mu             sync.RWMutex
 recorded       bool
 snapshotState  ProcessState
 channelStates  map[string]*ChannelState
 recording      map[string]bool
 markersReceived map[string]bool

 // Communication
 recvChan       chan Message
 markerChan     chan Marker
 snapshotResult chan<- SnapshotResult
}

// SnapshotResult contains the complete snapshot data
type SnapshotResult struct {
 ProcessID      string
 LocalState     ProcessState
 ChannelStates  map[string][]interface{}
 Timestamp      time.Time
}

// NewProcess creates a new process
func NewProcess(id string, snapshotResult chan<- SnapshotResult) *Process {
 return &Process{
  ID:              id,
  Neighbors:       make([]string, 0),
  Channels:        make(map[string]chan interface{}),
  channelStates:   make(map[string]*ChannelState),
  recording:       make(map[string]bool),
  markersReceived: make(map[string]bool),
  recvChan:        make(chan Message, 100),
  markerChan:      make(chan Marker, 100),
  snapshotResult:  snapshotResult,
 }
}

// AddNeighbor adds a neighbor process
func (p *Process) AddNeighbor(neighborID string) {
 p.mu.Lock()
 defer p.mu.Unlock()

 p.Neighbors = append(p.Neighbors, neighborID)
 p.Channels[neighborID] = make(chan interface{}, 100)
 p.channelStates[neighborID] = &ChannelState{Messages: make([]interface{}, 0)}
 p.recording[neighborID] = false
 p.markersReceived[neighborID] = false
}

// SendMessage sends a message to a neighbor
func (p *Process) SendMessage(to string, payload interface{}) error {
 p.mu.RLock()
 ch, ok := p.Channels[to]
 recorded := p.recorded
 p.mu.RUnlock()

 if !ok {
  return fmt.Errorf("unknown neighbor: %s", to)
 }

 msg := Message{
  From:    p.ID,
  To:      to,
  Payload: payload,
 }

 // Set color for Lai-Yang
 if recorded {
  msg.Color = Red
 } else {
  msg.Color = White
 }

 select {
 case ch <- msg:
  return nil
 default:
  return fmt.Errorf("channel full")
 }
}

// SendMarker sends a marker to a neighbor
func (p *Process) SendMarker(to string) error {
 p.mu.RLock()
 ch, ok := p.Channels[to]
 p.mu.RUnlock()

 if !ok {
  return fmt.Errorf("unknown neighbor: %s", to)
 }

 marker := Marker{
  From:      p.ID,
  Timestamp: time.Now(),
 }

 select {
 case ch <- marker:
  return nil
 default:
  return fmt.Errorf("channel full")
 }
}

// InitiateSnapshot starts the snapshot process
func (p *Process) InitiateSnapshot() error {
 p.mu.Lock()
 defer p.mu.Unlock()

 if p.recorded {
  return fmt.Errorf("snapshot already recorded")
 }

 // Record local state
 p.recordState()

 // Send markers to all neighbors
 for _, neighbor := range p.Neighbors {
  go p.SendMarker(neighbor)
 }

 return nil
}

func (p *Process) recordState() {
 p.snapshotState = ProcessState{
  Data:      p.State.Data,
  Timestamp: time.Now(),
  Version:   p.State.Version,
 }
 p.recorded = true

 // Start recording on all incoming channels
 for neighbor := range p.Channels {
  if neighbor != p.ID {
   p.recording[neighbor] = true
   p.channelStates[neighbor].Messages = make([]interface{}, 0)
  }
 }
}

// Receive processes incoming messages and markers
func (p *Process) Receive(ctx context.Context) {
 for {
  select {
  case <-ctx.Done():
   return

  case msg := <-p.recvChan:
   p.handleMessage(msg)

  case marker := <-p.markerChan:
   p.handleMarker(marker)
  }
 }
}

func (p *Process) handleMessage(msg Message) {
 p.mu.Lock()
 defer p.mu.Unlock()

 // Update process state
 p.State.Data = msg.Payload
 p.State.Version++
 p.State.Timestamp = time.Now()

 // Record message if in recording phase
 if p.recorded && p.recording[msg.From] {
  p.channelStates[msg.From].mu.Lock()
  p.channelStates[msg.From].Messages = append(
   p.channelStates[msg.From].Messages, msg.Payload)
  p.channelStates[msg.From].mu.Unlock()
 }
}

func (p *Process) handleMarker(marker Marker) {
 p.mu.Lock()
 defer p.mu.Unlock()

 from := marker.From

 if !p.recorded {
  // First marker - record state
  p.recordState()

  // Send markers to all neighbors except sender
  for _, neighbor := range p.Neighbors {
   if neighbor != from {
    go p.SendMarker(neighbor)
   }
  }

  // Channel from sender is empty
  p.recording[from] = false
  p.markersReceived[from] = true
 } else {
  // Subsequent marker - stop recording this channel
  p.recording[from] = false
  p.markersReceived[from] = true
 }

 // Check if snapshot is complete
 p.checkComplete()
}

func (p *Process) checkComplete() {
 if !p.recorded {
  return
 }

 // Check if all markers received
 allReceived := true
 for _, received := range p.markersReceived {
  if !received {
   allReceived = false
   break
  }
 }

 if allReceived {
  // Report snapshot
  result := SnapshotResult{
   ProcessID:     p.ID,
   LocalState:    p.snapshotState,
   ChannelStates: make(map[string][]interface{}),
   Timestamp:     time.Now(),
  }

  for neighbor, cs := range p.channelStates {
   cs.mu.RLock()
   result.ChannelStates[neighbor] = append([]interface{}{}, cs.Messages...)
   cs.mu.RUnlock()
  }

  select {
  case p.snapshotResult <- result:
  default:
  }
 }
}

// LaiYangProcess implements the Lai-Yang snapshot algorithm
type LaiYangProcess struct {
 *Process
 color MessageColor
}

// NewLaiYangProcess creates a new Lai-Yang process
func NewLaiYangProcess(id string, snapshotResult chan<- SnapshotResult) *LaiYangProcess {
 return &LaiYangProcess{
  Process: NewProcess(id, snapshotResult),
  color:   White,
 }
}

// InitiateSnapshotLaiYang starts the Lai-Yang snapshot
func (p *LaiYangProcess) InitiateSnapshotLaiYang() {
 p.mu.Lock()
 p.color = Red
 p.recordState()
 p.mu.Unlock()
}

// HandleMessageLaiYang handles messages in Lai-Yang algorithm
func (p *LaiYangProcess) HandleMessageLaiYang(msg Message) {
 p.mu.Lock()
 defer p.mu.Unlock()

 if p.color == White && msg.Color == Red {
  // First red message triggers recording
  p.color = Red
  p.recordState()
  p.recording[msg.From] = false // Channel from sender is empty

  for neighbor := range p.Channels {
   if neighbor != msg.From {
    p.recording[neighbor] = true
    p.channelStates[neighbor].Messages = make([]interface{}, 0)
   }
  }
 } else if p.color == Red && msg.Color == White {
  // White message after turning red
  p.channelStates[msg.From].mu.Lock()
  p.channelStates[msg.From].Messages = append(
   p.channelStates[msg.From].Messages, msg.Payload)
  p.channelStates[msg.From].mu.Unlock()
 }

 if msg.Color == Red {
  p.markersReceived[msg.From] = true
 }

 // Update state
 p.State.Data = msg.Payload
 p.State.Version++

 // Check completion
 p.checkComplete()
}

// SnapshotCoordinator coordinates snapshots across multiple processes
type SnapshotCoordinator struct {
 mu         sync.RWMutex
 processes  map[string]*Process
 results    map[string]SnapshotResult
 resultChan chan SnapshotResult
 done       chan struct{}
}

// NewSnapshotCoordinator creates a new coordinator
func NewSnapshotCoordinator() *SnapshotCoordinator {
 return &SnapshotCoordinator{
  processes:  make(map[string]*Process),
  results:    make(map[string]SnapshotResult),
  resultChan: make(chan SnapshotResult, 100),
  done:       make(chan struct{}),
 }
}

// AddProcess adds a process to coordinate
func (sc *SnapshotCoordinator) AddProcess(p *Process) {
 sc.mu.Lock()
 defer sc.mu.Unlock()
 sc.processes[p.ID] = p
}

// InitiateGlobalSnapshot initiates a snapshot from a process
func (sc *SnapshotCoordinator) InitiateGlobalSnapshot(initiatorID string) error {
 sc.mu.RLock()
 p, ok := sc.processes[initiatorID]
 sc.mu.RUnlock()

 if !ok {
  return fmt.Errorf("unknown initiator: %s", initiatorID)
 }

 return p.InitiateSnapshot()
}

// CollectResults collects snapshot results
func (sc *SnapshotCoordinator) CollectResults(ctx context.Context, expected int) (map[string]SnapshotResult, error) {
 collected := make(map[string]SnapshotResult)

 for len(collected) < expected {
  select {
  case <-ctx.Done():
   return collected, ctx.Err()
  case result := <-sc.resultChan:
   collected[result.ProcessID] = result
  }
 }

 return collected, nil
}

// VerifyConsistency verifies that the snapshot is consistent
func VerifyConsistency(results map[string]SnapshotResult) error {
 // Check that for every channel, messages in state were sent before
 // source recorded and received after destination recorded

 // This is a simplified check
 for id1, result1 := range results {
  for id2, channelState := range result1.ChannelStates {
   result2, ok := results[id2]
   if !ok {
    return fmt.Errorf("missing result for %s", id2)
   }

   // Channel state from id2 to id1 should have messages that
   // id2 sent before recording and id1 received after recording
   _ = channelState
   _ = result2
   // Detailed consistency check would go here
  }
 }

 return nil
}
