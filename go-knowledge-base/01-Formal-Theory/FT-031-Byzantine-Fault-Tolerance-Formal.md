# FT-031: Byzantine Fault Tolerance - Formal Theory

> **Dimension**: Formal Theory
> **Level**: S (>15KB)
> **Tags**: #byzantine #bft #pbft #hotstuff #distributed-systems
> **Authoritative Sources**:
>
> - Lamport, L., Shostak, R., & Pease, M. (1982). "The Byzantine Generals Problem". ACM TOPLAS
> - Castro, M., & Liskov, B. (2002). "Practical Byzantine Fault Tolerance". OSDI
> - Yin, M., et al. (2019). "HotStuff: BFT Consensus in the Lens of Blockchain". PODC

---

## 1. Theoretical Foundations

### 1.1 Problem Definition

**Definition 1.1 (Byzantine Generals Problem)**: Given $n$ generals, where $f$ may be traitors, find an algorithm ensuring:

1. All loyal generals agree on the same plan
2. A small number of traitors cannot cause loyal generals to adopt a bad plan

**Definition 1.2 (Byzantine Fault)**: A fault where a process behaves arbitrarily, including sending conflicting information to different processes.

### 1.2 Impossibility Results

**Theorem 1.1 (Byzantine Impossibility)**: No solution exists with $n \leq 3f$ generals.

*Proof (for n=3, f=1)*:

```
Scenario:
    L1 ──m1──> T
    L2 <──m2── T

L1 sends "attack" to T
T (traitor) sends "retreat" to L2

From L1's view: T could be loyal, L2 is traitor
From L2's view: T could be loyal, L1 is traitor
No way to distinguish! Therefore no agreement possible. ∎
```

**Theorem 1.2 (Upper Bound)**: Byzantine consensus is possible if and only if $n \geq 3f + 1$.

### 1.3 Security Properties

| Property | Definition |
|----------|------------|
| **Safety** | No two correct replicas commit different values |
| **Liveness** | All correct replicas eventually commit |
| **Accountability** | Faulty replicas can be identified |

---

## 2. Practical Byzantine Fault Tolerance (PBFT)

### 2.1 Protocol Overview

PBFT operates in three phases:

1. **Pre-prepare**: Leader assigns sequence number
2. **Prepare**: Replicas validate and prepare
3. **Commit**: Replicas commit the request

**Quorum Size**: $2f + 1$ out of $3f + 1$

### 2.2 PBFT Algorithm

```
Algorithm PBFT at replica i:

State:
  view: integer                  // Current view number
  log: Map<int, Request>         // Request log
  prepared: Map<int, Digest>     // Prepared certificates
  committed: Map<int, Digest>    // Committed certificates

ClientRequest(request):
  Send <REQUEST, request, timestamp> to leader

// Leader (primary for view v)
On ReceiveRequest(request):
  n ← NextSequenceNumber()
  digest ← Hash(request)
  Broadcast <PRE-PREPARE, v, n, digest, request>

// All replicas
On ReceivePrePrepare(v, n, digest, request):
  if ValidPrePrepare(v, n, digest, request):
    Store request
    Broadcast <PREPARE, v, n, digest, i>

On ReceivePrepare(v, n, digest, j):
  Store prepare certificate from j
  if CountPrepares(v, n, digest) ≥ 2f:
    MarkPrepared(n, digest)
    Broadcast <COMMIT, v, n, digest, i>

On ReceiveCommit(v, n, digest, j):
  Store commit certificate from j
  if CountCommits(v, n, digest) ≥ 2f + 1:
    Execute(request)
    ReplyToClient(result)
```

### 2.3 View Change

```
ViewChange(newView):
  // Triggered by timeout or suspected leader failure

  // Collect P-certificate (proof of preparation)
  P ← {(n, digest) | prepared in current view}

  Broadcast <VIEW-CHANGE, v+1, P, i>

  // New primary (p = v+1 mod n)
  On ReceiveViewChange(view, P_set, j):
    if CountViewChanges(view) ≥ 2f + 1:
      NewView ← SelectMaxSequence(P_set)
      Broadcast <NEW-VIEW, view+1, NewView>
```

---

## 3. HotStuff Algorithm

### 3.1 Key Innovation

HotStuff introduces linear communication complexity using:

- Threshold signatures for certificate aggregation
- Chained consensus (pipelining)
- Responsive communication (leader only waits for $f+1$ responses)

### 3.2 HotStuff Protocol

```
Algorithm HotStuff at replica i:

State:
  view: integer
  lockedQC: QuorumCertificate  // Highest committed QC
  prepareQC: QuorumCertificate // Current prepare QC

ProposePhase(leader):
  // Leader proposes block extending highest locked QC
  block ← CreateBlock(parent=lockedQC.block)
  Broadcast <PROPOSE, block, prepareQC>

VotePhase(replica):
  On ReceivePropose(block, parentQC):
    if ExtendsLockedQC(block) and ValidBlock(block):
      vote ← SignVote(block)
      Send <VOTE, vote> to next leader

NewViewPhase:
  // Collect votes, form QC
  if votes ≥ 2f + 1:
    newQC ← AggregateVotes(votes)
    if view = currentView:
      prepareQC ← newQC
      Goto ProposePhase
```

**Theorem 3.1 (HotStuff Complexity)**: HotStuff achieves $O(n)$ authenticator complexity per view.

*Proof*:

- Leader sends to all $n$: $O(n)$
- Replicas send to leader: $O(n)$ total
- Threshold signatures aggregate to constant size ∎

---

## 4. BFT Analysis

### 4.1 Quorum Analysis

| Property | Requirement | Explanation |
|----------|-------------|-------------|
| Intersection | $|Q_1 \cap Q_2| \geq f + 1$ | Any two quorums share a correct replica |
| Commit | $|Q| \geq 2f + 1$ | With $n=3f+1$, ensures majority correct |
| Prepare | $|Q| \geq 2f + 1$ | Same as commit |

**Theorem 4.1 (Quorum Intersection)**: With $n = 3f + 1$ and quorum size $2f + 1$, any two quorums intersect in at least $f + 1$ replicas.

*Proof*:
$$|Q_1 \cap Q_2| = |Q_1| + |Q_2| - |Q_1 \cup Q_2| \geq (2f+1) + (2f+1) - (3f+1) = f + 1$$ ∎

---

## 5. Go Implementation

```go
// Package bft provides Byzantine fault tolerance implementations
package bft

import (
 "crypto/ecdsa"
 "crypto/sha256"
 "fmt"
 "sync"
)

// Replica represents a BFT node
type Replica struct {
 ID      int
 Address string
 PubKey  *ecdsa.PublicKey
}

// Request represents a client request
type Request struct {
 Op        []byte
 Timestamp int64
 ClientID  string
}

// Message types
type MessageType int

const (
 MsgRequest MessageType = iota
 MsgPrePrepare
 MsgPrepare
 MsgCommit
 MsgViewChange
 MsgNewView
)

// Message wraps protocol messages
type Message struct {
 Type       MessageType
 View       int
 Sequence   uint64
 Digest     []byte
 Request    *Request
 ReplicaID  int
 Signature  []byte
}

// PBFT implements Practical Byzantine Fault Tolerance
type PBFT struct {
 id       int
 replicas []*Replica
 f        int // max faults tolerated

 mu         sync.RWMutex
 view       int
 seq        uint64
 log        map[uint64]*Request
 prepares   map[uint64]map[int][]byte // seq -> replica -> digest
 commits    map[uint64]map[int][]byte

 prepared   map[uint64]bool
 committed  map[uint64]bool

 isLeader   bool
 privateKey *ecdsa.PrivateKey
}

// NewPBFT creates a new PBFT instance
func NewPBFT(id int, replicas []*Replica, privateKey *ecdsa.PrivateKey) *PBFT {
 n := len(replicas) + 1 // include self
 f := (n - 1) / 3

 return &PBFT{
  id:        id,
  replicas:  replicas,
  f:         f,
  view:      0,
  log:       make(map[uint64]*Request),
  prepares:  make(map[uint64]map[int][]byte),
  commits:   make(map[uint64]map[int][]byte),
  prepared:  make(map[uint64]bool),
  committed: make(map[uint64]bool),
  privateKey: privateKey,
 }
}

// IsLeader returns true if this replica is the leader
func (p *PBFT) IsLeader() bool {
 return p.view%len(p.replicas) == p.id
}

// HandleRequest processes client request (leader only)
func (p *PBFT) HandleRequest(req *Request) error {
 if !p.IsLeader() {
  return fmt.Errorf("not leader")
 }

 p.mu.Lock()
 p.seq++
 seq := p.seq
 p.log[seq] = req
 p.mu.Unlock()

 digest := hashRequest(req)

 // Broadcast pre-prepare
 msg := &Message{
  Type:     MsgPrePrepare,
  View:     p.view,
  Sequence: seq,
  Digest:   digest,
  Request:  req,
  ReplicaID: p.id,
 }

 p.broadcast(msg)
 return nil
}

// HandlePrePrepare processes pre-prepare message
func (p *PBFT) HandlePrePrepare(msg *Message) error {
 if msg.Type != MsgPrePrepare {
  return fmt.Errorf("invalid message type")
 }

 // Verify message
 if !p.verifyMessage(msg) {
  return fmt.Errorf("verification failed")
 }

 p.mu.Lock()
 p.log[msg.Sequence] = msg.Request
 p.mu.Unlock()

 // Send prepare
 prepare := &Message{
  Type:      MsgPrepare,
  View:      msg.View,
  Sequence:  msg.Sequence,
  Digest:    msg.Digest,
  ReplicaID: p.id,
 }

 p.broadcast(prepare)
 return nil
}

// HandlePrepare processes prepare message
func (p *PBFT) HandlePrepare(msg *Message) error {
 p.mu.Lock()
 defer p.mu.Unlock()

 if p.prepares[msg.Sequence] == nil {
  p.prepares[msg.Sequence] = make(map[int][]byte)
 }
 p.prepares[msg.Sequence][msg.ReplicaID] = msg.Digest

 // Check if prepared
 if len(p.prepares[msg.Sequence]) >= 2*p.f {
  p.prepared[msg.Sequence] = true

  // Send commit
  commit := &Message{
   Type:      MsgCommit,
   View:      msg.View,
   Sequence:  msg.Sequence,
   Digest:    msg.Digest,
   ReplicaID: p.id,
  }
  p.broadcast(commit)
 }

 return nil
}

// HandleCommit processes commit message
func (p *PBFT) HandleCommit(msg *Message) error {
 p.mu.Lock()
 defer p.mu.Unlock()

 if p.commits[msg.Sequence] == nil {
  p.commits[msg.Sequence] = make(map[int][]byte)
 }
 p.commits[msg.Sequence][msg.ReplicaID] = msg.Digest

 // Check if committed
 if !p.committed[msg.Sequence] && len(p.commits[msg.Sequence]) >= 2*p.f+1 {
  p.committed[msg.Sequence] = true
  req := p.log[msg.Sequence]
  p.execute(req)
 }

 return nil
}

func (p *PBFT) execute(req *Request) {
 // Execute the request
 fmt.Printf("Executing request from %s\n", req.ClientID)
}

func (p *PBFT) broadcast(msg *Message) {
 // Sign message
 msg.Signature = p.sign(msg)

 // Send to all replicas
 for _, r := range p.replicas {
  p.send(r, msg)
 }
}

func (p *PBFT) send(replica *Replica, msg *Message) {
 // Network send
}

func (p *PBFT) verifyMessage(msg *Message) bool {
 // Verify signature
 return true
}

func (p *PBFT) sign(msg *Message) []byte {
 // Sign with private key
 return nil
}

func hashRequest(req *Request) []byte {
 h := sha256.New()
 h.Write(req.Op)
 return h.Sum(nil)
}

// HotStuff implements the HotStuff BFT protocol
type HotStuff struct {
 id      int
 replicas []*Replica
 f       int

 mu        sync.RWMutex
 view      int
 lockedQC  *QuorumCertificate
 prepareQC *QuorumCertificate

 privateKey *ecdsa.PrivateKey
}

// QuorumCertificate represents aggregated votes
type QuorumCertificate struct {
 View     int
 BlockHash []byte
 Signatures map[int][]byte
}

// NewHotStuff creates a new HotStuff instance
func NewHotStuff(id int, replicas []*Replica, privateKey *ecdsa.PrivateKey) *HotStuff {
 n := len(replicas) + 1
 f := (n - 1) / 3

 return &HotStuff{
  id:         id,
  replicas:   replicas,
  f:          f,
  privateKey: privateKey,
 }
}

// Propose creates a new proposal
func (h *HotStuff) Propose(parentHash []byte, cmds [][]byte) *Block {
 block := &Block{
  View:       h.view,
  ParentHash: parentHash,
  Commands:   cmds,
 }
 block.Hash = hashBlock(block)
 return block
}

// Vote casts a vote on a block
func (h *HotStuff) Vote(block *Block) *Vote {
 return &Vote{
  View:      h.view,
  BlockHash: block.Hash,
  ReplicaID: h.id,
  Signature: h.signBlock(block),
 }
}

// CreateQC aggregates votes into quorum certificate
func (h *HotStuff) CreateQC(votes []*Vote) (*QuorumCertificate, error) {
 if len(votes) < 2*h.f+1 {
  return nil, fmt.Errorf("insufficient votes")
 }

 qc := &QuorumCertificate{
  View:       votes[0].View,
  BlockHash:  votes[0].BlockHash,
  Signatures: make(map[int][]byte),
 }

 for _, v := range votes {
  qc.Signatures[v.ReplicaID] = v.Signature
 }

 return qc, nil
}

type Block struct {
 View       int
 ParentHash []byte
 Commands   [][]byte
 Hash       []byte
}

type Vote struct {
 View      int
 BlockHash []byte
 ReplicaID int
 Signature []byte
}

func hashBlock(b *Block) []byte {
 h := sha256.New()
 // Serialize and hash
 return h.Sum(nil)
}

func (h *HotStuff) signBlock(b *Block) []byte {
 // Sign block hash
 return nil
}
```

---

## 6. Visual Representations

### 6.1 PBFT Message Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PBFT MESSAGE FLOW                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Client      Leader(R0)      R1        R2        R3                          │
│    │           │              │          │          │                        │
│    │─REQUEST──>│              │          │          │                        │
│    │           │              │          │          │                        │
│    │    ┌──────┴──────────────┴────┐     │          │                        │
│    │    │ PRE-PREPARE             │     │          │                        │
│    │    │ <v,n,d,request>         │     │          │                        │
│    │    └──────┬──────────────┬────┘     │          │                        │
│    │           │              │          │          │                        │
│    │    ┌──────┴──────────────┴────┐     │          │                        │
│    │    │ PREPARE                  │     │          │                        │
│    │    │ <v,n,d,i>                │     │          │                        │
│    │    │ (2f+1 prepares = cert)   │     │          │                        │
│    │    └──────┬──────────────┬────┘     │          │                        │
│    │           │              │          │          │                        │
│    │    ┌──────┴──────────────┴────┐     │          │                        │
│    │    │ COMMIT                   │     │          │                        │
│    │    │ <v,n,d,i>                │     │          │                        │
│    │    │ (2f+1 commits = exec)    │     │          │                        │
│    │    └──────┬──────────────┬────┘     │          │                        │
│    │           │              │          │          │                        │
│    │<─REPLY────│              │          │          │                        │
│    │           │              │          │          │                        │
│                                                                              │
│  Total messages per request: 1 + 3n (including client)                      │
│  With n=4, f=1: 1 + 12 = 13 messages                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 HotStuff Chained Consensus

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HOTSTUFF CHAINED CONSENSUS                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  View v: Leader Lv                                                           │
│                                                                              │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐                   │
│  │ Block   │───>│ Block   │───>│ Block   │───>│ Block   │                   │
│  │ B1      │    │ B2      │    │ B3      │    │ B4      │                   │
│  │         │    │ QC(B1)  │    │ QC(B2)  │    │ QC(B3)  │                   │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘                   │
│       │              │              │              │                         │
│       v              v              v              v                         │
│   PREPARE       PRE-COMMIT       COMMIT        DECIDE                      │
│       │              │              │              │                         │
│       │              │              │              │                         │
│  3-chain commit: B1 is committed when B4 is created                        │
│                                                                              │
│  Pipelining:                                                                 │
│  ────────────                                                                │
│                                                                              │
│  Time: ──────────────────────────────────────────────────────────>          │
│                                                                              │
│  View 1: [Prepare(B1)]──[PreCommit(B1)]──[Commit(B1)]                      │
│  View 2:                 [Prepare(B2)]───[PreCommit(B2)]──[Commit(B2)]      │
│  View 3:                                [Prepare(B3)]───[PreCommit(B3)]     │
│                                                                              │
│  Overlapping views enable higher throughput                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 BFT Comparison Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BFT PROTOCOL COMPARISON                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Property        │ PBFT      │ HotStuff  │ Tendermint │ SBFT      │        │
│  ─────────────────┼───────────┼───────────┼────────────┼───────────┤        │
│  Replica count   │ 3f+1      │ 3f+1      │ 3f+1       │ 3f+1      │        │
│  Quorum size     │ 2f+1      │ 2f+1      │ 2f+1       │ 2f+1      │        │
│  Messages/req    │ O(n²)     │ O(n)      │ O(n²)      │ O(n)      │        │
│  Latency (hops)  │ 3         │ 3         │ 3          │ 2         │        │
│  View change     │ Expensive │ Seamless  │ Expensive  │ Fast      │        │
│  Responsiveness  │ No        │ Yes       │ No         │ Yes       │        │
│  Optimistic path │ No        │ Yes       │ No         │ Yes       │        │
│  Implementation  │ Complex   │ Moderate  │ Complex    │ Complex   │        │
│  ─────────────────┴───────────┴───────────┴────────────┴───────────┘        │
│                                                                              │
│  Key Terms:                                                                  │
│  • Responsiveness: Leader can proceed as fast as network allows             │
│  • Optimistic path: Fast commit when leader is honest                        │
│  • View change: Replacing faulty leader                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BYZANTINE FAULT TOLERANCE SUMMARY                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Key Insight: Byzantine consensus requires n ≥ 3f+1 replicas                │
│                                                                              │
│  Protocol Selection:                                                         │
│  • Classic deployment: PBFT (proven, simple)                                │
│  • High throughput: HotStuff (linear communication)                         │
│  • Blockchain: Tendermint/HotStuff (chained)                                │
│                                                                              │
│  Trade-offs:                                                                 │
│  • More replicas = more fault tolerance, higher overhead                    │
│  • Threshold signatures reduce communication but add computation            │
│  • View changes are expensive - optimize for stability                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
