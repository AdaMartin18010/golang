# FT-030: Consensus Performance - Formal Analysis

> **Dimension**: Formal Theory
> **Level**: S (>15KB)
> **Tags**: #consensus #performance #latency #throughput #optimization
> **Authoritative Sources**:
>
> - Liskov, B., & Cowling, J. (2012). "Viewstamped Replication Revisited". MIT
> - Moraru, I., et al. (2013). "There Is More Consensus in Egalitarian Parliaments". SOSP
> - Poke, M., & Schiavoni, V. (2015). "Dare: High-Performance State Machine Replication". DSN

---

## 1. Theoretical Foundations

### 1.1 Performance Metrics

**Definition 1.1 (Consensus Latency)**: The time from proposal submission until the value is committed by a majority.

$$
L = T_{propose} + T_{agree} + T_{commit}
$$

**Definition 1.2 (Consensus Throughput)**: The number of consensus instances completed per unit time.

$$
\Theta = \frac{N_{committed}}{T_{elapsed}}
$$

**Definition 1.3 (Message Complexity)**: Total number of messages exchanged per consensus instance.

$$
M = \sum_{i=1}^{n} m_i
$$

### 1.2 Performance Models

**Network Model**:

| Model | Latency | Bandwidth | Description |
|-------|---------|-----------|-------------|
| LAN | $\delta \approx 0.1$ms | 1-10 Gbps | Datacenter |
| WAN | $\delta \approx 50-200$ms | 10-100 Mbps | Cross-region |
| Geo | $\delta \approx 100-300$ms | Variable | Global |

**Theorem 1.1 (Latency Lower Bound)**: Any consensus protocol requires at least $2\delta$ latency where $\delta$ is the one-way network delay.

*Proof*:

- Process must send to others: $\delta$
- Quorum must respond: $\delta$
- Total: $2\delta$ ∎

---

## 2. Multi-Paxos Performance

### 2.1 Batch Optimization

**Algorithm 1: Batched Multi-Paxos**:

```
Batch Size Optimization:

Variables:
  B: batch size
  D: command size
  N: number of nodes
  δ: network latency
  μ: processing overhead per command

Latency per batch:
  L_batch = 2δ + B × μ

Throughput:
  Θ = B / L_batch = B / (2δ + B × μ)

Optimal batch size depends on load:
  - Low load: B = 1 (minimize latency)
  - High load: B = B_max (maximize throughput)

Adaptive batching:
  if queue_length > threshold:
    B = min(queue_length, B_max)
  else:
    B = 1
```

**Theorem 2.1 (Batch Throughput Gain)**: Batching increases throughput by factor $O(B)$ when $B \times \mu \ll 2\delta$.

### 2.2 Pipeline Optimization

**Algorithm 2: Pipelined Multi-Paxos**:

```
Pipelining allows overlapping consensus instances:

Time: ──────────────────────────────────────────>

Instance 1: [Prepare]──[Accept]──[Commit]
Instance 2:      [Prepare]──[Accept]──[Commit]
Instance 3:           [Prepare]──[Accept]──[Commit]

Pipeline depth limited by:
  - Network bandwidth
  - Leader processing capacity
  - Follower apply rate

Optimal pipeline depth:
  P = floor(Bandwidth × RTT / Command_Size)
```

### 2.3 Leader Rotation

```
Leader rotation strategies:

1. Static Leader:
   Pros: No leader election overhead
   Cons: Hotspot, geographic bias

2. Random Rotation:
   Pros: Load balancing
   Cons: Frequent elections

3. Latency-Aware:
   Select leader closest to majority
   Leader = argmin_{p ∈ Π} Σ_{q ∈ Q} latency(p, q)
```

---

## 3. EPaxos Performance Analysis

### 3.1 Fast Path Analysis

**Theorem 3.1 (EPaxos Fast Path Probability)**: For workload with conflict rate $\alpha$, the fast path probability is approximately $(1-\alpha)^{F}$ where $F$ is the fast quorum size.

*Proof*:

- Fast path succeeds if all fast quorum members report same deps
- Conflict probability with any in-flight command: $\alpha$
- No conflict with any of $F$ members: $(1-\alpha)^F$ ∎

### 3.2 Dependency Graph Complexity

```
Dependency resolution cost:

Let C = number of commands
Let D = average dependencies per command

Execution time:
  - Build graph: O(C × D)
  - Topological sort: O(C + D)
  - SCC detection: O(C + D)

Total: O(C × D) where D depends on conflict rate

Conflict rate impact:
  α = 0.01 → D ≈ 0.01 × C
  α = 0.1  → D ≈ 0.1 × C
  α = 0.5  → D ≈ 0.5 × C (degrades to sequential)
```

---

## 4. Go Implementation

```go
// Package consensusperf provides performance-optimized consensus
package consensusperf

import (
 "context"
 "sync"
 "sync/atomic"
 "time"
)

// Batcher implements adaptive batching
type Batcher struct {
 minBatch    int
 maxBatch    int
 timeout     time.Duration

 queue       chan Command
 batchCh     chan []Command
 stopCh      chan struct{}
 wg          sync.WaitGroup
}

type Command struct {
 ID   uint64
 Data []byte
 Done chan<- Result
}

type Result struct {
 Index uint64
 Error error
}

// NewBatcher creates a new adaptive batcher
func NewBatcher(minBatch, maxBatch int, timeout time.Duration) *Batcher {
 return &Batcher{
  minBatch: minBatch,
  maxBatch: maxBatch,
  timeout:  timeout,
  queue:    make(chan Command, 10000),
  batchCh:  make(chan []Command),
  stopCh:   make(chan struct{}),
 }
}

// Start begins batching
func (b *Batcher) Start() {
 b.wg.Add(1)
 go b.run()
}

func (b *Batcher) run() {
 defer b.wg.Done()

 batch := make([]Command, 0, b.maxBatch)
 timer := time.NewTimer(b.timeout)
 defer timer.Stop()

 for {
  select {
  case <-b.stopCh:
   if len(batch) > 0 {
    b.batchCh <- batch
   }
   return

  case cmd := <-b.queue:
   batch = append(batch, cmd)

   if len(batch) >= b.maxBatch {
    b.batchCh <- batch
    batch = make([]Command, 0, b.maxBatch)
    timer.Reset(b.timeout)
   }

  case <-timer.C:
   if len(batch) >= b.minBatch {
    b.batchCh <- batch
    batch = make([]Command, 0, b.maxBatch)
   }
   timer.Reset(b.timeout)
  }
 }
}

// Submit submits a command for batching
func (b *Batcher) Submit(cmd Command) {
 b.queue <- cmd
}

// NextBatch returns the next batch
func (b *Batcher) NextBatch() <-chan []Command {
 return b.batchCh
}

// Stop stops the batcher
func (b *Batcher) Stop() {
 close(b.stopCh)
 b.wg.Wait()
}

// Pipeliner manages pipelined consensus
type Pipeliner struct {
 maxPipeline int
 inFlight    int64
 results     chan uint64
}

// NewPipeliner creates a new pipeliner
func NewPipeliner(maxPipeline int) *Pipeliner {
 return &Pipeliner{
  maxPipeline: maxPipeline,
  results:     make(chan uint64, maxPipeline),
 }
}

// CanPropose returns true if pipeline has room
func (p *Pipeliner) CanPropose() bool {
 return atomic.LoadInt64(&p.inFlight) < int64(p.maxPipeline)
}

// StartProposal records start of proposal
func (p *Pipeliner) StartProposal() {
 atomic.AddInt64(&p.inFlight, 1)
}

// CompleteProposal records completion
func (p *Pipeliner) CompleteProposal(index uint64) {
 atomic.AddInt64(&p.inFlight, -1)
 p.results <- index
}

// Metrics tracks consensus performance
type Metrics struct {
 proposals    uint64
 commits      uint64
 latencySum   uint64
 latencyCount uint64
 bytesIn      uint64
 bytesOut     uint64
}

func (m *Metrics) RecordProposal() {
 atomic.AddUint64(&m.proposals, 1)
}

func (m *Metrics) RecordCommit(latency time.Duration) {
 atomic.AddUint64(&m.commits, 1)
 atomic.AddUint64(&m.latencySum, uint64(latency.Microseconds()))
 atomic.AddUint64(&m.latencyCount, 1)
}

func (m *Metrics) RecordBytesIn(n int) {
 atomic.AddUint64(&m.bytesIn, uint64(n))
}

func (m *Metrics) RecordBytesOut(n int) {
 atomic.AddUint64(&m.bytesOut, uint64(n))
}

func (m *Metrics) Latency() time.Duration {
 count := atomic.LoadUint64(&m.latencyCount)
 if count == 0 {
  return 0
 }
 sum := atomic.LoadUint64(&m.latencySum)
 return time.Duration(sum/count) * time.Microsecond
}

func (m *Metrics) Throughput() float64 {
 return float64(atomic.LoadUint64(&m.commits))
}
```

---

## 5. Visual Representations

### 5.1 Latency vs Throughput Tradeoff

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONSENSUS LATENCY-THROUGHPUT TRADEOFF                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Throughput ▲                                                                │
│             │                                                                │
│     Θ_max   ├─────────────────────────────██████████                        │
│             │                         █████          │                        │
│             │                    █████               │                        │
│             │               █████                    │  Saturation Region    │
│             │          █████                         │                        │
│             │     █████                              │                        │
│             │ ████                                   │                        │
│             │██                                      │                        │
│             └─────────────────────────────────────────> Load                  │
│                                                      │                        │
│  Latency    ▲                                        │                        │
│             │                                        │                        │
│             │                                ████████████                    │
│             │                          ██████          │                      │
│             │                    █████                 │                      │
│             │              █████                       │  Queueing Delay      │
│             │        █████                           │                        │
│             │   █████                                │                        │
│             │████                                    │                        │
│             └─────────────────────────────────────────> Load                  │
│                                                      │                        │
│  Key Observations:                                                           │
│  1. Low load: Minimal latency, throughput proportional to load              │
│  2. Saturation: Throughput plateaus, latency increases linearly             │
│  3. Batching extends saturation point at cost of latency                    │
│  4. Pipelining increases throughput without latency penalty                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Multi-Paxos Optimization Stack

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MULTI-PAXOS OPTIMIZATION STACK                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Layer 4: Application Level                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ • Request coalescing                                               │   │
│  │ • Read-write splitting                                             │   │
│  │ • Consistency level selection                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Layer 3: Consensus Optimizations                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ • Batching (amortize overhead)                                     │   │
│  │ • Pipelining (overlap instances)                                   │   │
│  │ • Leader stickiness (reduce elections)                             │   │
│  │ • Fast path (skip prepare when safe)                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Layer 2: Protocol Level                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ • Selective replication (quorum sizes)                             │   │
│  │ • Learner nodes (scale reads)                                      │   │
│  │ • Log compaction (snapshotting)                                    │   │
│  │ • Chunked replication (large values)                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Layer 1: Network Level                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ • TCP tuning (Nagle, keepalive)                                    │   │
│  │ • Connection pooling                                               │   │
│  │ • Compression                                                      │   │
│  │ • Topology-aware placement                                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Performance Impact:                                                         │
│  ┌─────────────────┬─────────────────┬─────────────────┐                    │
│  │ Optimization    │ Latency Impact  │ Throughput Gain │                    │
│  ├─────────────────┼─────────────────┼─────────────────┤                    │
│  │ Batching        │ +10-50%         │ +5-10x          │                    │
│  │ Pipelining      │ Minimal         │ +3-5x           │                    │
│  │ Compression     │ +5-10%          │ +1-2x (WAN)     │                    │
│  │ Fast Path       │ -50% (happy path│ No change       │                    │
│  └─────────────────┴─────────────────┴─────────────────┘                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Geographic Consensus Performance

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GEOGRAPHIC CONSENSUS PERFORMANCE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Deployment Scenario: 5 replicas                                             │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┐              │
│  │ Location     │ us-west      │ us-east      │ eu-west      │              │
│  │              │ (leader)     │              │              │              │
│  ├──────────────┼──────────────┼──────────────┼──────────────┤              │
│  │ Latency      │ 0ms          │ 80ms         │ 150ms        │              │
│  └──────────────┴──────────────┴──────────────┴──────────────┘              │
│                                                                              │
│  Protocol Performance Comparison:                                            │
│                                                                              │
│  Single-Datacenter (all in us-west):                                         │
│  ┌──────────────────────────────────────────────────────────────┐           │
│  │  Phase 1 (Prepare):   1 RTT = 1ms                           │           │
│  │  Phase 2 (Accept):    1 RTT = 1ms                           │           │
│  │  Total:               2ms                                   │           │
│  └──────────────────────────────────────────────────────────────┘           │
│                                                                              │
│  Multi-Region (us-west leader):                                              │
│  ┌──────────────────────────────────────────────────────────────┐           │
│  │  Phase 1: Max(0, 80, 150) + response = 300ms                │           │
│  │  Phase 2: Max(0, 80, 150) + response = 300ms                │           │
│  │  Total:   600ms                                             │           │
│  │  Slowdown: 300x!                                            │           │
│  └──────────────────────────────────────────────────────────────┘           │
│                                                                              │
│  EPaxos (leaderless, closest replica):                                       │
│  ┌──────────────────────────────────────────────────────────────┐           │
│  │  Fast Path: 1 RTT to nearest 2 replicas = 2ms               │           │
│  │  Slow Path: 2 RTT = 4ms                                     │           │
│  │  Average: Depends on conflict rate                          │           │
│  └──────────────────────────────────────────────────────────────┘           │
│                                                                              │
│  Flexible Paxos (separate quorums):                                          │
│  ┌──────────────────────────────────────────────────────────────┐           │
│  │  Phase 1 (all nodes):  300ms one-time cost                  │           │
│  │  Phase 2 (local):      1ms                                  │           │
│  │  Amortized:            ~1ms after first command             │           │
│  └──────────────────────────────────────────────────────────────┘           │
│                                                                              │
│  Key Insight: WAN latency dominates, optimizations critical                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONSENSUS PERFORMANCE SUMMARY                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Theoretical Limits:                                                         │
│  • Minimum latency: 2 × RTT (fundamental network bound)                     │
│  • Maximum throughput: Limited by leader capacity                           │
│  • Optimal batch size: Balance latency vs throughput                        │
│                                                                              │
│  Optimization Priorities:                                                    │
│  1. Batching: Highest impact on throughput                                  │
│  2. Pipelining: Better resource utilization                                 │
│  3. Leader placement: Critical for WAN deployments                          │
│  4. Fast paths: Reduce latency for common cases                             │
│                                                                              │
│  Deployment Recommendations:                                                 │
│  • Single DC: Standard Multi-Paxos with batching                            │
│  • Multi-region: EPaxos or Flexible Paxos                                   │
│  • Read-heavy: Learner nodes for reads                                      │
│  • Write-heavy: Dedicated leader, pipelining                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
