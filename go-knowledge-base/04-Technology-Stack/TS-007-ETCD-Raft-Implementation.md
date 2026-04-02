# TS-007: etcd Raft Implementation - Distributed Consensus Internals

> **维度**: Technology Stack
> **级别**: S (16+ KB)
> **标签**: #etcd #raft #consensus #distributed-systems #go
> **权威来源**:
>
> - [etcd Raft Paper](https://raft.github.io/raft.pdf) - Diego Ongaro & John Ousterhout
> - [etcd Documentation](https://etcd.io/docs/) - CNCF
> - [Raft Consensus Algorithm](https://raft.github.io/) - raft.github.io

---

## 1. Raft Consensus Algorithm

### 1.1 Raft State Machine

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Raft State Machine                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Server States                                       │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │         ┌─────────────┐                                               │  │
│  │         │   Follower  │◄────────────────────────┐                     │  │
│  │         │             │                         │                     │  │
│  │         │ • Passive   │                         │                     │  │
│  │         │ • Responds  │                         │                     │  │
│  │         │   to RPCs   │                         │                     │  │
│  │         └──────┬──────┘                         │                     │  │
│  │                │                                │                     │  │
│  │                │ Election timeout               │                     │  │
│  │                │ without leader                 │                     │  │
│  │                │                                │                     │  │
│  │                ▼                                │                     │  │
│  │         ┌─────────────┐    Discover higher    │                     │  │
│  │    ┌───►│  Candidate  │────term or new leader─┘                     │  │
│  │    │    │             │                                               │  │
│  │    │    │ • Votes for │                                               │  │
│  │    │    │   itself    │                                               │  │
│  │    │    │ • Sends     │                                               │  │
│  │    │    │   RequestVote                                               │  │
│  │    │    └──────┬──────┘                                               │  │
│  │    │           │                                                       │  │
│  │    │           │ Majority votes received                               │  │
│  │    │           │                                                       │  │
│  │    │           ▼                                                       │  │
│  │    │    ┌─────────────┐                                               │  │
│  │    └────┤    Leader   │                                               │  │
│  │         │             │                                               │  │
│  │         │ • Handles   │                                               │  │
│  │         │   all client│                                               │  │
│  │         │   requests  │                                               │  │
│  │         │ • Sends     │                                               │  │
│  │         │   heartbeats│                                               │  │
│  │         └─────────────┘                                               │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    State Variables                                     │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Persistent state (on all servers):                                    │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │ currentTerm │ int    │ Latest term server has seen            │  │  │
│  │  │ votedFor    │ int    │ CandidateId that received vote         │  │  │
│  │  │ log[]       │ Log    │ Log entries; each: (term, command)     │  │  │
│  │  │             │        │   index starts at 1                    │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Volatile state (on all servers):                                      │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │ commitIndex │ int    │ Highest log entry known to be committed  │  │  │
│  │  │ lastApplied │ int    │ Highest log entry applied to state mach  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Volatile state (on leaders only):                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │ nextIndex[] │ int[]  │ For each server, index of next log entry │  │  │
│  │  │             │        │ to send to that server                   │  │  │
│  │  │ matchIndex[]│ int[]  │ For each server, index of highest log    │  │  │
│  │  │             │        │ entry known to be replicated             │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Log Replication Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Raft Log Replication                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    AppendEntries RPC                                   │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Request from Leader:                                                  │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │ term         │ Leader's term                                     │  │  │
│  │  │ leaderId     │ So follower can redirect clients                  │  │  │
│  │  │ prevLogIndex │ Index of log entry immediately preceding new ones │  │  │
│  │  │ prevLogTerm  │ Term of prevLogIndex entry                        │  │  │
│  │  │ entries[]    │ Log entries to store (empty for heartbeat)        │  │  │
│  │  │ leaderCommit │ Leader's commitIndex                              │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Response from Follower:                                               │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │ term    │ CurrentTerm, for leader to update itself              │  │  │
│  │  │ success │ True if follower contained entry matching             │  │  │
│  │  │         │ prevLogIndex and prevLogTerm                          │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Replication Flow Example                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Leader (S1)          Follower (S2)          Follower (S3)            │  │
│  │  Log: [1,1][2,1][3,2]  Log: [1,1][2,1]       Log: [1,1][2,1]         │  │
│  │        │                   │                      │                   │  │
│  │        │ AppendEntries     │                      │                   │  │
│  │        │ prev=2,1          │                      │                   │  │
│  │        │ entry=[3,2]       │                      │                   │  │
│  │        ├──────────────────►│                      │                   │  │
│  │        │                   │ Append to log        │                   │  │
│  │        │ success=true      │                      │                   │  │
│  │        │◄──────────────────┤                      │                   │  │
│  │        │                   │                      │                   │  │
│  │        │ AppendEntries     │                      │                   │  │
│  │        ├───────────────────┼─────────────────────►│                   │  │
│  │        │                   │                      │ Append to log     │  │
│  │        │ success=true      │                      │                   │  │
│  │        │◄──────────────────┼──────────────────────┤                   │  │
│  │        │                   │                      │                   │  │
│  │        │ Majority achieved │                      │                   │  │
│  │        │ (2 of 3 replicas) │                      │                   │  │
│  │        ▼                   ▼                      ▼                   │  │
│  │  Update commitIndex to 3                                             │  │
│  │        │                   │                      │                   │  │
│  │        │ Next heartbeat    │                      │                   │  │
│  │        │ leaderCommit=3    │                      │                   │  │
│  │        ├──────────────────►│                      │                   │  │
│  │        │                   │ Apply to state mach  │                   │  │
│  │        ├───────────────────┼─────────────────────►│                   │  │
│  │        │                   │                      │ Apply to state    │  │
│  │        │                   │                      │ machine           │  │
│  │                                                                        │  │
│  │  Safety Rule: Only commit entry from current term once stored         │  │
│  │  on majority. This ensures all future leaders will have it.           │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Log Consistency Maintenance                         │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Scenario: Leader has new entries, follower is behind                  │  │
│  │                                                                        │  │
│  │  Leader S1: [1,1][2,1][3,2][4,2][5,3]                                 │  │
│  │                    ▲                                                  │  │
│  │                    └─ nextIndex[S2] = 5 (optimized)                    │  │
│  │                                                                        │  │
│  │  Follower S2: [1,1][2,1][3,3][4,3]                                    │  │
│  │                    ▲                                                  │  │
│  │                    └─ Conflict at index 3 (term 2 vs 3)                │  │
│  │                                                                        │  │
│  │  Step 1: Leader sends AppendEntries(prev=4,2)                         │  │
│  │  Step 2: S2 rejects (log doesn't match at prevLogIndex)               │  │
│  │  Step 3: Leader decrements nextIndex[S2] = 4                          │  │
│  │  Step 4: Leader sends AppendEntries(prev=3,2)                         │  │
│  │  Step 5: S2 rejects (term mismatch at 3)                              │  │
│  │  Step 6: Leader decrements nextIndex[S2] = 3                          │  │
│  │  Step 7: Leader sends AppendEntries(prev=2,1)                         │  │
│  │  Step 8: S2 accepts!                                                  │  │
│  │  Step 9: Leader sends entries [3,2][4,2][5,3]                         │  │
│  │  Step 10: S2 overwrites [3,3][4,3] and appends new entries            │  │
│  │                                                                        │  │
│  │  Result: Logs converge to leader's log                                │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. etcd Raft Implementation

```go
package raft

import (
    "context"
    "fmt"
    "sync"
    "time"

    "go.etcd.io/etcd/raft/v3"
    "go.etcd.io/etcd/raft/v3/raftpb"
    "go.etcd.io/etcd/server/v3/etcdserver/api/snap"
    "go.etcd.io/etcd/server/v3/wal"
    "go.etcd.io/etcd/server/v3/wal/walpb"
)

// RaftNode Raft 节点封装
type RaftNode struct {
    id          uint64
    peers       []uint64
    raftNode    raft.Node
    transport   *Transport
    storage     *Storage

    // 通道
    proposeC    chan []byte
    confChangeC chan raftpb.ConfChange
    commitC     chan *Commit
    errorC      chan error

    // 状态
    mu          sync.RWMutex
    snapState   snap.Snapshot
}

// Commit 已提交的日志条目
type Commit struct {
    Data       []byte
    Index      uint64
}

// Config Raft 配置
type Config struct {
    ID              uint64
    Peers           []uint64
    ElectionTick    int
    HeartbeatTick   int
    MaxSizePerMsg   uint64
    MaxInflightMsgs int
    DataDir         string
}

// NewRaftNode 创建 Raft 节点
func NewRaftNode(cfg *Config) (*RaftNode, error) {
    // 创建存储
    storage, err := NewStorage(cfg.DataDir)
    if err != nil {
        return nil, err
    }

    // Raft 配置
    rc := &raft.Config{
        ID:              cfg.ID,
        ElectionTick:    cfg.ElectionTick,
        HeartbeatTick:   cfg.HeartbeatTick,
        Storage:         storage.raftStorage,
        MaxSizePerMsg:   cfg.MaxSizePerMsg,
        MaxInflightMsgs: cfg.MaxInflightMsgs,
    }

    // 创建节点
    var r raft.Node
    if storage.IsEmpty() {
        // 新集群
        r = raft.StartNode(rc, cfg.Peers)
    } else {
        // 重启
        r = raft.RestartNode(rc)
    }

    rn := &RaftNode{
        id:          cfg.ID,
        peers:       cfg.Peers,
        raftNode:    r,
        storage:     storage,
        proposeC:    make(chan []byte),
        confChangeC: make(chan raftpb.ConfChange),
        commitC:     make(chan *Commit, 1000),
        errorC:      make(chan error),
    }

    return rn, nil
}

// Run 运行 Raft 节点
func (rn *RaftNode) Run(ctx context.Context) {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            rn.raftNode.Stop()
            return

        case <-ticker.C:
            rn.raftNode.Tick()

        case rd := <-rn.raftNode.Ready():
            // 保存到 WAL
            if err := rn.storage.Save(rd.HardState, rd.Entries); err != nil {
                rn.errorC <- err
                return
            }

            // 发送消息给其他节点
            rn.transport.Send(rd.Messages)

            // 应用已提交的条目
            if ok := rn.publishEntries(rd.CommittedEntries); !ok {
                rn.errorC <- fmt.Errorf("failed to publish entries")
                return
            }

            // 应用快照
            if !raft.IsEmptySnap(rd.Snapshot) {
                if err := rn.storage.ApplySnapshot(rd.Snapshot); err != nil {
                    rn.errorC <- err
                    return
                }
            }

            // 通知 Ready 处理完成
            rn.raftNode.Advance()

        case prop := <-rn.proposeC:
            // 提议数据
            if err := rn.raftNode.Propose(ctx, prop); err != nil {
                rn.errorC <- err
            }

        case cc := <-rn.confChangeC:
            // 配置变更
            rn.raftNode.ProposeConfChange(ctx, cc)
        }
    }
}

// publishEntries 发布已提交的条目
func (rn *RaftNode) publishEntries(ents []raftpb.Entry) bool {
    for i := range ents {
        switch ents[i].Type {
        case raftpb.EntryNormal:
            if len(ents[i].Data) == 0 {
                // 忽略空条目
                continue
            }
            select {
            case rn.commitC <- &Commit{Data: ents[i].Data, Index: ents[i].Index}:
            case <-time.After(10 * time.Second):
                return false
            }

        case raftpb.EntryConfChange:
            var cc raftpb.ConfChange
            if err := cc.Unmarshal(ents[i].Data); err != nil {
                continue
            }
            rn.raftNode.ApplyConfChange(cc)
        }
    }
    return true
}

// Propose 提议数据
func (rn *RaftNode) Propose(data []byte) {
    rn.proposeC <- data
}

// AddPeer 添加节点
func (rn *RaftNode) AddPeer(nodeID uint64) {
    cc := raftpb.ConfChange{
        Type:    raftpb.ConfChangeAddNode,
        NodeID:  nodeID,
    }
    rn.confChangeC <- cc
}

// RemovePeer 移除节点
func (rn *RaftNode) RemovePeer(nodeID uint64) {
    cc := raftpb.ConfChange{
        Type:    raftpb.ConfChangeRemoveNode,
        NodeID:  nodeID,
    }
    rn.confChangeC <- cc
}

// Storage 存储封装
type Storage struct {
    wal         *wal.WAL
    snap        *snap.Snapshotter
    raftStorage *raft.MemoryStorage
}

// NewStorage 创建存储
func NewStorage(dataDir string) (*Storage, error) {
    // 实现 WAL 和快照存储
    return &Storage{}, nil
}

// Save 保存硬状态和条目
func (s *Storage) Save(st raftpb.HardState, ents []raftpb.Entry) error {
    // 保存到 WAL
    return nil
}

// ApplySnapshot 应用快照
func (s *Storage) ApplySnapshot(snap raftpb.Snapshot) error {
    return nil
}

// IsEmpty 检查是否为空
func (s *Storage) IsEmpty() bool {
    return false
}

// Transport 网络传输层
type Transport struct {
    peers   map[uint64]*Peer
}

// Peer 对等节点
type Peer struct {
    ID      uint64
    Address string
    Client  *RaftClient
}

// RaftClient Raft RPC 客户端
type RaftClient struct {
    // gRPC 连接
}

// Send 发送消息
func (t *Transport) Send(msgs []raftpb.Message) {
    for _, msg := range msgs {
        if peer, ok := t.peers[msg.To]; ok {
            peer.Client.Send(msg)
        }
    }
}

// Send 发送单条消息
func (c *RaftClient) Send(msg raftpb.Message) error {
    // 发送 AppendEntries/RequestVote/Heartbeat
    return nil
}
```

---

## 3. Configuration Best Practices

```yaml
# etcd 配置
name: 'etcd-node-1'
data-dir: '/var/lib/etcd'
wal-dir: '/var/lib/etcd/wal'

# 集群配置
initial-cluster: 'etcd-node-1=http://192.168.1.1:2380,etcd-node-2=http://192.168.1.2:2380,etcd-node-3=http://192.168.1.3:2380'
initial-cluster-token: 'etcd-cluster-1'
initial-cluster-state: 'new'

# 网络配置
listen-peer-urls: 'http://0.0.0.0:2380'
listen-client-urls: 'http://0.0.0.0:2379'
advertise-client-urls: 'http://192.168.1.1:2379'

# 性能调优
heartbeat-interval: 100    # 毫秒
election-timeout: 1000     # 毫秒
snapshot-count: 100000
max-snapshots: 5
max-wals: 5

# 配额
quota-backend-bytes: 8589934592  # 8GB
```

---

## 4. Visual Representations

### Raft Election

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Raft Leader Election                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Node 1 (Term 1, Follower)    Node 2 (Term 1, Follower)    Node 3          │
│  ┌─────────────────────┐      ┌─────────────────────┐      ┌─────────────┐│
│  │ Election timeout    │      │ Election timeout    │      │             ││
│  │ (random 150-300ms)  │      │ (random 150-300ms)  │      │             ││
│  │ Expires first!      │      │                     │      │             ││
│  └──────────┬──────────┘      └─────────────────────┘      └─────────────┘│
│             │                                                               │
│             ▼                                                               │
│  ┌─────────────────────┐                                                    │
│  │ Become Candidate    │                                                    │
│  │ currentTerm = 2     │                                                    │
│  │ votedFor = 1        │                                                    │
│  └──────────┬──────────┘                                                    │
│             │                                                               │
│             │ RequestVote RPC                                               │
│             │ term=2, candidateId=1, lastLogIndex=10, lastLogTerm=1         │
│             ├────────────────────────►┌─────────────────────┐               │
│             ├────────────────────────►│ Node 2: Votes YES   │               │
│             │                         │ (higher term)       │               │
│             │◄────────────────────────└─────────────────────┘               │
│             │                                                               │
│             ├────────────────────────►┌─────────────────────┐               │
│             ├────────────────────────►│ Node 3: Votes YES   │               │
│             │                         │ (higher term)       │               │
│             │◄────────────────────────└─────────────────────┘               │
│             │                                                               │
│             │ Majority achieved (2/3)                                       │
│             ▼                                                               │
│  ┌─────────────────────┐                                                    │
│  │ Become LEADER       │                                                    │
│  │ Send heartbeats     │                                                    │
│  │ to all followers    │                                                    │
│  └─────────────────────┘                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Membership Change

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Raft Membership Change (Joint Consensus)                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Phase 1: C_old ∪ C_new (Joint Consensus)                                   │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  Old Config: {S1, S2, S3}                                             │  │
│  │  New Config: {S1, S2, S3, S4, S5}  (add 2 nodes)                      │  │
│  │  Joint Config: Both configs must agree                                │  │
│  │                                                                        │  │
│  │  Entry committed when:                                                │  │
│  │  - Majority of C_old acknowledges AND                                 │  │
│  │  - Majority of C_new acknowledges                                     │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Phase 2: C_new                                                           │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  After joint config committed, leader replicates C_new entry          │  │
│  │  New config: {S1, S2, S3, S4, S5}                                     │  │
│  │                                                                        │  │
│  │  S4, S5 can now participate in consensus                              │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Safety: Never removes majority in single step                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. References

1. **Ongaro, D., & Ousterhout, J.** (2014). In Search of an Understandable Consensus Algorithm. *USENIX ATC*.
2. **etcd Documentation** (2024). etcd.io/docs
3. **etcd Raft Implementation** (2024). github.com/etcd-io/raft

---

*Document Version: 1.0 | Last Updated: 2024*
