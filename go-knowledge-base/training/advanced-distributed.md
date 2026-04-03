# Advanced Distributed Systems

## Module Overview

**Duration:** Self-paced (Ongoing after Week 4)
**Prerequisites:** Week 4 completion
**Learning Goal:** Master advanced distributed systems concepts and implementations

---

## Learning Objectives

By the end of this advanced track, you will be able to:

1. **Consensus Algorithms**
   - Understand Raft and Paxos protocols
   - Implement Raft consensus
   - Analyze consensus trade-offs

2. **Distributed Data Structures**
   - Implement CRDTs (Conflict-free Replicated Data Types)
   - Understand vector clocks and logical time
   - Design distributed caches

3. **Event Sourcing and CQRS**
   - Implement event-sourced systems
   - Design read model projections
   - Handle event versioning

4. **Byzantine Fault Tolerance**
   - Understand BFT consensus
   - Implement PBFT concepts
   - Analyze blockchain consensus

5. **Production Systems**
   - Design for operational excellence
   - Implement chaos engineering
   - Handle complex failure scenarios

---

## Reading Assignments

### Consensus Algorithms

1. **[Raft Consensus Formal](../01-Formal-Theory/FT-002-Raft-Consensus-Formal.md)**
   - Leader election
   - Log replication
   - Safety guarantees

2. **[Paxos Formal](../01-Formal-Theory/FT-006-Paxos-Formal.md)**
   - Single-decree Paxos
   - Multi-Paxos optimization
   - Compare with Raft

3. **[Distributed Consensus](../01-Formal-Theory/FT-003-Distributed-Consensus-Raft-Paxos.md)**
   - Consensus problem definition
   - FLP impossibility
   - Practical solutions

### Distributed Data

1. **[CRDT Formal](../01-Formal-Theory/FT-018-CRDT-Formal.md)**
   - State-based CRDTs
   - Operation-based CRDTs
   - Merge semantics

2. **[Vector Clocks](../01-Formal-Theory/FT-005-Vector-Clocks-Formal.md)**
   - Logical time
   - Happens-before detection
   - Version vectors

3. **[Gossip Protocols](../01-Formal-Theory/FT-011-Gossip-Protocols.md)**
   - Epidemic broadcast
   - Anti-entropy mechanisms
   - Membership protocols

### Advanced Patterns

1. **[Event Sourcing Formal](../03-Engineering-CloudNative/EC-015-Event-Sourcing-Formal.md)**
   - Event store design
   - Snapshot strategies
   - Projection patterns

2. **[Byzantine Fault Tolerance](../01-Formal-Theory/FT-007-Byzantine-Fault-Tolerance.md)**
   - Byzantine generals problem
   - PBFT algorithm
   - BFT vs CFT trade-offs

---

## Hands-on Projects

### Project 1: Raft Implementation

Implement a Raft consensus module:

```go
package raft

import (
    "context"
    "sync"
    "time"
)

// NodeState represents the state of a Raft node
type NodeState int

const (
    Follower NodeState = iota
    Candidate
    Leader
)

// Raft implements the Raft consensus algorithm
type Raft struct {
    id      string
    peers   []string
    state   NodeState

    // Persistent state
    currentTerm int64
    votedFor    string
    log         []LogEntry

    // Volatile state
    commitIndex int64
    lastApplied int64

    // Leader state
    nextIndex  map[string]int64
    matchIndex map[string]int64

    mu          sync.RWMutex
    electionTimer *time.Timer

    // Channels
    applyCh     chan ApplyMsg
    heartbeatCh chan struct{}
}

type LogEntry struct {
    Term    int64
    Index   int64
    Command interface{}
}

type ApplyMsg struct {
    CommandIndex int64
    Command      interface{}
}

// StartRaft initializes a new Raft node
func StartRaft(id string, peers []string, applyCh chan ApplyMsg) *Raft {
    r := &Raft{
        id:          id,
        peers:       peers,
        state:       Follower,
        nextIndex:   make(map[string]int64),
        matchIndex:  make(map[string]int64),
        applyCh:     applyCh,
        heartbeatCh: make(chan struct{}),
    }

    r.resetElectionTimer()
    go r.run()

    return r
}

func (r *Raft) run() {
    for {
        switch r.State() {
        case Follower:
            r.runFollower()
        case Candidate:
            r.runCandidate()
        case Leader:
            r.runLeader()
        }
    }
}

func (r *Raft) runFollower() {
    select {
    case <-r.electionTimer.C:
        r.becomeCandidate()
    case <-r.heartbeatCh:
        r.resetElectionTimer()
    }
}

func (r *Raft) runCandidate() {
    r.mu.Lock()
    r.currentTerm++
    r.votedFor = r.id
    term := r.currentTerm
    r.mu.Unlock()

    votes := 1
    voteCh := make(chan bool, len(r.peers))

    // Request votes from all peers
    for _, peer := range r.peers {
        if peer == r.id {
            continue
        }
        go func(peer string) {
            vote := r.requestVote(peer, term)
            voteCh <- vote
        }(peer)
    }

    // Count votes
    for i := 0; i < len(r.peers)-1; i++ {
        select {
        case vote := <-voteCh:
            if vote {
                votes++
                if votes > len(r.peers)/2 {
                    r.becomeLeader()
                    return
                }
            }
        case <-time.After(r.electionTimeout()):
            return // Election timeout, restart
        }
    }
}

func (r *Raft) runLeader() {
    // Send heartbeats
    r.broadcastHeartbeat()

    select {
    case <-time.After(50 * time.Millisecond):
        // Continue as leader
    }
}

func (r *Raft) becomeCandidate() {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.state = Candidate
}

func (r *Raft) becomeLeader() {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.state = Leader

    // Initialize leader state
    for _, peer := range r.peers {
        r.nextIndex[peer] = int64(len(r.log)) + 1
        r.matchIndex[peer] = 0
    }
}

// Propose submits a command to the Raft log
func (r *Raft) Propose(cmd interface{}) (int64, int64, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.state != Leader {
        return 0, 0, false
    }

    entry := LogEntry{
        Term:    r.currentTerm,
        Index:   int64(len(r.log)) + 1,
        Command: cmd,
    }

    r.log = append(r.log, entry)

    return entry.Index, entry.Term, true
}

// State returns current node state
func (r *Raft) State() NodeState {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.state
}
```

**Requirements:**

- Leader election
- Log replication
- Safety guarantees
- Membership changes
- Snapshot support

**Deliverable:** Working Raft implementation with tests

---

### Project 2: Distributed Key-Value Store

Build a distributed KV store using Raft:

```go
package kvstore

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"

    "github.com/hashicorp/raft"
)

// KVStore is a distributed key-value store
type KVStore struct {
    mu     sync.RWMutex
    data   map[string]string
    raft   *raft.Raft
}

type KVCommand struct {
    Op    string `json:"op"`
    Key   string `json:"key"`
    Value string `json:"value,omitempty"`
}

func NewKVStore() *KVStore {
    return &KVStore{
        data: make(map[string]string),
    }
}

// Apply implements raft.FSM
func (kv *KVStore) Apply(log *raft.Log) interface{} {
    var cmd KVCommand
    if err := json.Unmarshal(log.Data, &cmd); err != nil {
        return err
    }

    kv.mu.Lock()
    defer kv.mu.Unlock()

    switch cmd.Op {
    case "set":
        kv.data[cmd.Key] = cmd.Value
        return nil
    case "delete":
        delete(kv.data, cmd.Key)
        return nil
    default:
        return fmt.Errorf("unknown operation: %s", cmd.Op)
    }
}

// Get returns a value (local read, might be stale)
func (kv *KVStore) Get(key string) (string, bool) {
    kv.mu.RLock()
    defer kv.mu.RUnlock()
    val, ok := kv.data[key]
    return val, ok
}

// Set sets a value (goes through Raft)
func (kv *KVStore) Set(ctx context.Context, key, value string) error {
    cmd := KVCommand{
        Op:    "set",
        Key:   key,
        Value: value,
    }

    data, err := json.Marshal(cmd)
    if err != nil {
        return err
    }

    future := kv.raft.Apply(data, 10*time.Second)
    return future.Error()
}

// Snapshot implements raft.FSM
func (kv *KVStore) Snapshot() (raft.FSMSnapshot, error) {
    kv.mu.RLock()
    defer kv.mu.RUnlock()

    // Copy data
    data := make(map[string]string)
    for k, v := range kv.data {
        data[k] = v
    }

    return &KVSnapshot{data: data}, nil
}

// Restore implements raft.FSM
func (kv *KVStore) Restore(rc io.ReadCloser) error {
    var data map[string]string
    if err := json.NewDecoder(rc).Decode(&data); err != nil {
        return err
    }

    kv.mu.Lock()
    defer kv.mu.Unlock()

    kv.data = data
    return nil
}

type KVSnapshot struct {
    data map[string]string
}

func (s *KVSnapshot) Persist(sink raft.SnapshotSink) error {
    err := json.NewEncoder(sink).Encode(s.data)
    if err != nil {
        sink.Cancel()
        return err
    }
    return sink.Close()
}

func (s *KVSnapshot) Release() {}
```

**Deliverable:** Working distributed KV store with CLI

---

### Project 3: CRDT Counter

Implement a G-Counter CRDT:

```go
package crdt

import (
    "encoding/json"
    "sync"
)

// GCounter is a grow-only counter CRDT
type GCounter struct {
    mu     sync.RWMutex
    id     string
    values map[string]uint64
}

func NewGCounter(id string) *GCounter {
    return &GCounter{
        id:     id,
        values: make(map[string]uint64),
    }
}

// Increment increases the local counter
func (c *GCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.values[c.id]++
}

// Value returns the total count
func (c *GCounter) Value() uint64 {
    c.mu.RLock()
    defer c.mu.RUnlock()

    var total uint64
    for _, v := range c.values {
        total += v
    }
    return total
}

// Merge combines two GCounters (takes maximum for each replica)
func (c *GCounter) Merge(other *GCounter) {
    c.mu.Lock()
    defer c.mu.Unlock()

    other.mu.RLock()
    defer other.mu.RUnlock()

    for replica, value := range other.values {
        if current, ok := c.values[replica]; !ok || value > current {
            c.values[replica] = value
        }
    }
}

// PNCounter is a counter that supports increment and decrement
type PNCounter struct {
    increments *GCounter
    decrements *GCounter
}

func NewPNCounter(id string) *PNCounter {
    return &PNCounter{
        increments: NewGCounter(id),
        decrements: NewGCounter(id),
    }
}

func (c *PNCounter) Increment() {
    c.increments.Increment()
}

func (c *PNCounter) Decrement() {
    c.decrements.Increment()
}

func (c *PNCounter) Value() int64 {
    return int64(c.increments.Value()) - int64(c.decrements.Value())
}

func (c *PNCounter) Merge(other *PNCounter) {
    c.increments.Merge(other.increments)
    c.decrements.Merge(other.decrements)
}
```

**Extensions:**

- Implement LWW-Register CRDT
- Implement OR-Set CRDT
- Build distributed shopping cart

**Deliverable:** CRDT library with 3 types

---

### Project 4: Event Sourced System

Build an event-sourced aggregate:

```go
package eventsource

import (
    "context"
    "fmt"
    "time"
)

// Event represents a domain event
type Event interface {
    EventType() string
    AggregateID() string
    EventVersion() int
    OccurredAt() time.Time
}

// Aggregate is an event-sourced aggregate root
type Aggregate interface {
    AggregateID() string
    AggregateVersion() int
    ApplyEvent(event Event) error
    UncommittedEvents() []Event
    MarkCommitted()
}

// Order aggregate example
type Order struct {
    id       string
    version  int
    customer string
    items    []OrderItem
    status   OrderStatus

    uncommitted []Event
}

type OrderItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

type OrderStatus int

const (
    OrderPending OrderStatus = iota
    OrderConfirmed
    OrderShipped
    OrderDelivered
    OrderCancelled
)

// Events
type OrderCreated struct {
    OrderID   string
    Customer  string
    Timestamp time.Time
}

func (e OrderCreated) EventType() string     { return "OrderCreated" }
func (e OrderCreated) AggregateID() string   { return e.OrderID }
func (e OrderCreated) EventVersion() int     { return 1 }
func (e OrderCreated) OccurredAt() time.Time { return e.Timestamp }

type ItemAdded struct {
    OrderID   string
    ProductID string
    Quantity  int
    Price     float64
    Timestamp time.Time
}

func (e ItemAdded) EventType() string     { return "ItemAdded" }
func (e ItemAdded) AggregateID() string   { return e.OrderID }
func (e ItemAdded) EventVersion() int     { return 1 }
func (e ItemAdded) OccurredAt() time.Time { return e.Timestamp }

// Factory functions
func NewOrder(id, customer string) (*Order, error) {
    if id == "" || customer == "" {
        return nil, fmt.Errorf("id and customer are required")
    }

    order := &Order{
        id: id,
    }

    order.raiseEvent(OrderCreated{
        OrderID:   id,
        Customer:  customer,
        Timestamp: time.Now(),
    })

    return order, nil
}

func (o *Order) AddItem(productID string, quantity int, price float64) error {
    if o.status != OrderPending {
        return fmt.Errorf("cannot add items to order with status %v", o.status)
    }

    if quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }

    o.raiseEvent(ItemAdded{
        OrderID:   o.id,
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
        Timestamp: time.Now(),
    })

    return nil
}

func (o *Order) raiseEvent(event Event) {
    o.uncommitted = append(o.uncommitted, event)
    o.ApplyEvent(event)
}

// ApplyEvent applies an event to update state
func (o *Order) ApplyEvent(event Event) error {
    switch e := event.(type) {
    case OrderCreated:
        o.id = e.OrderID
        o.customer = e.Customer
        o.status = OrderPending

    case ItemAdded:
        o.items = append(o.items, OrderItem{
            ProductID: e.ProductID,
            Quantity:  e.Quantity,
            Price:     e.Price,
        })

    default:
        return fmt.Errorf("unknown event type: %s", event.EventType())
    }

    o.version++
    return nil
}

// Getters
func (o *Order) AggregateID() string       { return o.id }
func (o *Order) AggregateVersion() int     { return o.version }
func (o *Order) UncommittedEvents() []Event { return o.uncommitted }
func (o *Order) MarkCommitted()            { o.uncommitted = nil }

// Event Store
type EventStore interface {
    Append(ctx context.Context, aggregateID string, events []Event, expectedVersion int) error
    GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
}

// Repository
type OrderRepository struct {
    eventStore EventStore
}

func (r *OrderRepository) Load(ctx context.Context, id string) (*Order, error) {
    events, err := r.eventStore.GetEvents(ctx, id, 0)
    if err != nil {
        return nil, err
    }

    if len(events) == 0 {
        return nil, fmt.Errorf("order not found: %s", id)
    }

    order := &Order{id: id}
    for _, event := range events {
        if err := order.ApplyEvent(event); err != nil {
            return nil, err
        }
    }

    return order, nil
}

func (r *OrderRepository) Save(ctx context.Context, order *Order) error {
    events := order.UncommittedEvents()
    if len(events) == 0 {
        return nil
    }

    err := r.eventStore.Append(ctx, order.AggregateID(), events, order.AggregateVersion()-len(events))
    if err != nil {
        return err
    }

    order.MarkCommitted()
    return nil
}
```

**Deliverable:** Event-sourced order system with projections

---

## Assessment Criteria

### Advanced Projects (60%)

| Project | Weight | Requirements |
|---------|--------|--------------|
| Raft Implementation | 25% | Pass all Jepsen-style tests |
| Distributed KV Store | 20% | Linearizable reads/writes |
| CRDT Library | 10% | Convergence guarantees |
| Event Sourcing | 5% | Event versioning, snapshots |

### Code Review (25%)

- Algorithm correctness
- Test coverage (>90%)
- Documentation quality
- Performance considerations

### Presentation (15%)

- Explain consensus algorithm choice
- Demonstrate fault tolerance
- Discuss trade-offs

---

## Resources

### Papers

1. "In Search of an Understandable Consensus Algorithm" (Raft)
2. "Paxos Made Simple" (Lamport)
3. "Conflict-free Replicated Data Types" (Shapiro et al.)

### Books

- "Designing Data-Intensive Applications" - Martin Kleppmann
- "Database Internals" - Alex Petrov

### Courses

- MIT 6.824: Distributed Systems
- Coursera: Cloud Computing Concepts

---

*Complete all projects to achieve Specialist level certification*
