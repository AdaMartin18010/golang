# 分布式共识 Raft 实现 (Distributed Consensus Raft Implementation)

> **分类**: 工程与云原生
> **标签**: #raft #consensus #distributed-systems #etcd
> **参考**: Raft Paper, etcd Raft, Consul

---

## 目录

- [分布式共识 Raft 实现 (Distributed Consensus Raft Implementation)](#分布式共识-raft-实现-distributed-consensus-raft-implementation)
  - [目录](#目录)
  - [Raft 架构](#raft-架构)
  - [完整 Raft 实现](#完整-raft-实现)

## Raft 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Raft Consensus Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Raft Node States                                  │   │
│  │                                                                      │   │
│  │     ┌──────────┐                                                    │   │
│  │     │ Follower │◄───────────────────────────────────┐              │   │
│  │     └────┬─────┘                                    │              │   │
│  │          │ election timeout                         │              │   │
│  │          ▼                                          │              │   │
│  │     ┌──────────┐          ┌──────────┐              │              │   │
│  │     │ Candidate│─────────►│  Leader  │──────────────┘              │   │
│  │     └──────────┘  majority └────┬─────┘                             │   │
│  │                                 │                                   │   │
│  │          ┌──────────────────────┼──────────────────────┐            │   │
│  │          │                      │                      │            │   │
│  │          ▼                      ▼                      ▼            │   │
│  │     ┌──────────┐          ┌──────────┐          ┌──────────┐       │   │
│  │     │ Follower │          │ Follower │          │ Follower │       │   │
│  │     └──────────┘          └──────────┘          └──────────┘       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Log Replication:                                                            │
│  Leader ──► AppendEntries RPC ──► Followers                                  │
│         ◄── Response (success/fail)                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整 Raft 实现

```go
package raft

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"
)

// NodeState Raft 节点状态
type NodeState int

const (
    StateFollower NodeState = iota
    StateCandidate
    StateLeader
)

func (s NodeState) String() string {
    switch s {
    case StateFollower:
        return "Follower"
    case StateCandidate:
        return "Candidate"
    case StateLeader:
        return "Leader"
    default:
        return "Unknown"
    }
}

// Config Raft 配置
type Config struct {
    NodeID uint64

    // 选举超时 (随机范围)
    ElectionTimeoutMin time.Duration
    ElectionTimeoutMax time.Duration

    // 心跳间隔
    HeartbeatInterval time.Duration

    // 集群节点
    Peers []uint64
}

// Raft 节点
type Raft struct {
    config Config

    // 持久化状态
    currentTerm uint64
    votedFor    uint64
    log         []LogEntry

    // 易失状态
    state        NodeState
    leaderID     uint64

    // 提交状态
    commitIndex uint64
    lastApplied uint64

    // 领导者状态 (仅 leader)
    nextIndex  map[uint64]uint64
    matchIndex map[uint64]uint64

    // 通信
    rpcCh     chan RPC
    applyCh   chan ApplyMsg

    // 控制
    shutdownCh chan struct{}
    wg         sync.WaitGroup

    // 选举计时器
    electionTimer  *time.Timer
    heartbeatTimer *time.Timer

    mu sync.RWMutex
}

// LogEntry 日志条目
type LogEntry struct {
    Term    uint64
    Index   uint64
    Command interface{}
}

// RPC 请求/响应接口
type RPC interface {
    GetTerm() uint64
}

// RequestVoteRequest 请求投票请求
type RequestVoteRequest struct {
    Term         uint64
    CandidateID  uint64
    LastLogIndex uint64
    LastLogTerm  uint64
}

func (r *RequestVoteRequest) GetTerm() uint64 { return r.Term }

// RequestVoteResponse 请求投票响应
type RequestVoteResponse struct {
    Term        uint64
    VoteGranted bool
}

// AppendEntriesRequest 追加条目请求
type AppendEntriesRequest struct {
    Term         uint64
    LeaderID     uint64
    PrevLogIndex uint64
    PrevLogTerm  uint64
    Entries      []LogEntry
    LeaderCommit uint64
}

func (a *AppendEntriesRequest) GetTerm() uint64 { return a.Term }

// AppendEntriesResponse 追加条目响应
type AppendEntriesResponse struct {
    Term    uint64
    Success bool
    // 快速回退优化
    ConflictIndex uint64
    ConflictTerm  uint64
}

// ApplyMsg 应用消息
type ApplyMsg struct {
    CommandValid bool
    Command      interface{}
    CommandIndex uint64
}

// NewRaft 创建 Raft 节点
func NewRaft(config Config, applyCh chan ApplyMsg) *Raft {
    r := &Raft{
        config:     config,
        log:        make([]LogEntry, 0),
        state:      StateFollower,
        rpcCh:      make(chan RPC, 100),
        applyCh:    applyCh,
        shutdownCh: make(chan struct{}),
        nextIndex:  make(map[uint64]uint64),
        matchIndex: make(map[uint64]uint64),
    }

    // 启动主循环
    r.wg.Add(1)
    go r.run()

    return r
}

// run 主循环
func (r *Raft) run() {
    defer r.wg.Done()

    // 启动选举计时器
    r.resetElectionTimer()

    for {
        select {
        case <-r.shutdownCh:
            return

        case rpc := <-r.rpcCh:
            r.handleRPC(rpc)

        default:
            switch r.state {
            case StateFollower:
                r.runFollower()
            case StateCandidate:
                r.runCandidate()
            case StateLeader:
                r.runLeader()
            }
        }
    }
}

func (r *Raft) runFollower() {
    select {
    case <-r.electionTimer.C:
        // 选举超时，转换为候选者
        r.becomeCandidate()
    }
}

func (r *Raft) runCandidate() {
    // 开始选举
    r.startElection()

    // 等待选举结果或超时
    electionTimeout := r.randomElectionTimeout()

    select {
    case <-time.After(electionTimeout):
        // 选举超时，重新开始
        return
    }
}

func (r *Raft) runLeader() {
    // 发送心跳
    r.broadcastHeartbeat()

    time.Sleep(r.config.HeartbeatInterval)
}

// becomeCandidate 成为候选者
func (r *Raft) becomeCandidate() {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.state = StateCandidate
    r.currentTerm++
    r.votedFor = r.config.NodeID

    fmt.Printf("Node %d became Candidate at term %d\n", r.config.NodeID, r.currentTerm)
}

// becomeLeader 成为领导者
func (r *Raft) becomeLeader() {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.state = StateLeader
    r.leaderID = r.config.NodeID

    // 初始化领导者状态
    lastLogIndex := r.getLastLogIndex()
    for _, peer := range r.config.Peers {
        r.nextIndex[peer] = lastLogIndex + 1
        r.matchIndex[peer] = 0
    }

    fmt.Printf("Node %d became Leader at term %d\n", r.config.NodeID, r.currentTerm)

    // 立即发送心跳
    go r.broadcastHeartbeat()
}

// startElection 开始选举
func (r *Raft) startElection() {
    r.mu.RLock()
    term := r.currentTerm
    lastLogIndex := r.getLastLogIndex()
    lastLogTerm := r.getLastLogTerm()
    r.mu.RUnlock()

    // 并行向所有节点请求投票
    votes := 1 // 投给自己
    var voteMu sync.Mutex

    for _, peer := range r.config.Peers {
        if peer == r.config.NodeID {
            continue
        }

        go func(peer uint64) {
            req := &RequestVoteRequest{
                Term:         term,
                CandidateID:  r.config.NodeID,
                LastLogIndex: lastLogIndex,
                LastLogTerm:  lastLogTerm,
            }

            resp := r.sendRequestVote(peer, req)

            if resp.VoteGranted {
                voteMu.Lock()
                votes++
                if votes > len(r.config.Peers)/2 {
                    r.becomeLeader()
                }
                voteMu.Unlock()
            } else if resp.Term > term {
                r.stepDown(resp.Term)
            }
        }(peer)
    }
}

// broadcastHeartbeat 广播心跳
func (r *Raft) broadcastHeartbeat() {
    r.mu.RLock()
    term := r.currentTerm
    leaderCommit := r.commitIndex
    r.mu.RUnlock()

    for _, peer := range r.config.Peers {
        if peer == r.config.NodeID {
            continue
        }

        go func(peer uint64) {
            r.mu.RLock()
            nextIdx := r.nextIndex[peer]
            prevLogIndex := nextIdx - 1
            prevLogTerm := r.getLogTerm(prevLogIndex)

            // 获取需要发送的条目
            entries := r.getEntriesFrom(nextIdx)
            r.mu.RUnlock()

            req := &AppendEntriesRequest{
                Term:         term,
                LeaderID:     r.config.NodeID,
                PrevLogIndex: prevLogIndex,
                PrevLogTerm:  prevLogTerm,
                Entries:      entries,
                LeaderCommit: leaderCommit,
            }

            resp := r.sendAppendEntries(peer, req)

            r.handleAppendEntriesResponse(peer, resp, nextIdx)
        }(peer)
    }
}

// handleRPC 处理 RPC
func (r *Raft) handleRPC(rpc RPC) {
    switch req := rpc.(type) {
    case *RequestVoteRequest:
        r.handleRequestVote(req)
    case *AppendEntriesRequest:
        r.handleAppendEntries(req)
    }
}

// handleRequestVote 处理投票请求
func (r *Raft) handleRequestVote(req *RequestVoteRequest) *RequestVoteResponse {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 如果请求任期更小，拒绝
    if req.Term < r.currentTerm {
        return &RequestVoteResponse{Term: r.currentTerm, VoteGranted: false}
    }

    // 如果请求任期更大，更新任期
    if req.Term > r.currentTerm {
        r.stepDown(req.Term)
    }

    // 检查是否已投票
    if r.votedFor != 0 && r.votedFor != req.CandidateID {
        return &RequestVoteResponse{Term: r.currentTerm, VoteGranted: false}
    }

    // 检查日志是否最新
    if !r.isLogUpToDate(req.LastLogIndex, req.LastLogTerm) {
        return &RequestVoteResponse{Term: r.currentTerm, VoteGranted: false}
    }

    // 投票
    r.votedFor = req.CandidateID
    r.resetElectionTimer()

    return &RequestVoteResponse{Term: r.currentTerm, VoteGranted: true}
}

// handleAppendEntries 处理追加条目
func (r *Raft) handleAppendEntries(req *AppendEntriesRequest) *AppendEntriesResponse {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 如果请求任期更小，拒绝
    if req.Term < r.currentTerm {
        return &AppendEntriesResponse{Term: r.currentTerm, Success: false}
    }

    // 如果请求任期更大，更新任期
    if req.Term > r.currentTerm {
        r.stepDown(req.Term)
    }

    // 重置选举计时器（收到有效心跳）
    r.resetElectionTimer()
    r.leaderID = req.LeaderID

    // 检查日志一致性
    if req.PrevLogIndex > 0 {
        if req.PrevLogIndex > r.getLastLogIndex() {
            return &AppendEntriesResponse{
                Term:          r.currentTerm,
                Success:       false,
                ConflictIndex: r.getLastLogIndex() + 1,
            }
        }

        if r.getLogTerm(req.PrevLogIndex) != req.PrevLogTerm {
            // 找到冲突任期的第一个索引
            conflictTerm := r.getLogTerm(req.PrevLogIndex)
            conflictIndex := req.PrevLogIndex

            for conflictIndex > 1 && r.getLogTerm(conflictIndex-1) == conflictTerm {
                conflictIndex--
            }

            return &AppendEntriesResponse{
                Term:          r.currentTerm,
                Success:       false,
                ConflictIndex: conflictIndex,
                ConflictTerm:  conflictTerm,
            }
        }
    }

    // 追加条目
    r.appendEntries(req.PrevLogIndex, req.Entries)

    // 更新提交索引
    if req.LeaderCommit > r.commitIndex {
        r.commitIndex = min(req.LeaderCommit, r.getLastLogIndex())
        r.applyCommitted()
    }

    return &AppendEntriesResponse{Term: r.currentTerm, Success: true}
}

// Propose 提议命令
func (r *Raft) Propose(command interface{}) (uint64, uint64, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.state != StateLeader {
        return 0, 0, fmt.Errorf("not leader")
    }

    entry := LogEntry{
        Term:    r.currentTerm,
        Index:   r.getLastLogIndex() + 1,
        Command: command,
    }

    r.log = append(r.log, entry)

    // 触发立即复制
    go r.broadcastHeartbeat()

    return entry.Term, entry.Index, nil
}

// Helper methods

func (r *Raft) stepDown(term uint64) {
    r.currentTerm = term
    r.votedFor = 0
    r.state = StateFollower
    r.leaderID = 0
}

func (r *Raft) getLastLogIndex() uint64 {
    if len(r.log) == 0 {
        return 0
    }
    return r.log[len(r.log)-1].Index
}

func (r *Raft) getLastLogTerm() uint64 {
    if len(r.log) == 0 {
        return 0
    }
    return r.log[len(r.log)-1].Term
}

func (r *Raft) getLogTerm(index uint64) uint64 {
    if index == 0 || index > uint64(len(r.log)) {
        return 0
    }
    return r.log[index-1].Term
}

func (r *Raft) getEntriesFrom(index uint64) []LogEntry {
    if index > uint64(len(r.log)) {
        return nil
    }
    return r.log[index-1:]
}

func (r *Raft) appendEntries(prevIndex uint64, entries []LogEntry) {
    // 删除冲突条目并追加新条目
    for i, entry := range entries {
        idx := prevIndex + uint64(i) + 1
        if idx <= uint64(len(r.log)) {
            if r.log[idx-1].Term != entry.Term {
                // 删除冲突及之后的所有条目
                r.log = r.log[:idx-1]
                r.log = append(r.log, entries[i:]...)
                break
            }
        } else {
            r.log = append(r.log, entries[i:]...)
            break
        }
    }
}

func (r *Raft) applyCommitted() {
    for r.lastApplied < r.commitIndex {
        r.lastApplied++
        entry := r.log[r.lastApplied-1]

        r.applyCh <- ApplyMsg{
            CommandValid: true,
            Command:      entry.Command,
            CommandIndex: entry.Index,
        }
    }
}

func (r *Raft) isLogUpToDate(lastLogIndex, lastLogTerm uint64) bool {
    myLastIndex := r.getLastLogIndex()
    myLastTerm := r.getLastLogTerm()

    if lastLogTerm != myLastTerm {
        return lastLogTerm > myLastTerm
    }
    return lastLogIndex >= myLastIndex
}

func (r *Raft) resetElectionTimer() {
    if r.electionTimer != nil {
        r.electionTimer.Stop()
    }
    r.electionTimer = time.NewTimer(r.randomElectionTimeout())
}

func (r *Raft) randomElectionTimeout() time.Duration {
    min := r.config.ElectionTimeoutMin
    max := r.config.ElectionTimeoutMax
    return min + time.Duration(rand.Int63n(int64(max-min)))
}

func (r *Raft) sendRequestVote(peer uint64, req *RequestVoteRequest) *RequestVoteResponse {
    // 实际实现需要网络通信
    return &RequestVoteResponse{Term: req.Term, VoteGranted: false}
}

func (r *Raft) sendAppendEntries(peer uint64, req *AppendEntriesRequest) *AppendEntriesResponse {
    // 实际实现需要网络通信
    return &AppendEntriesResponse{Term: req.Term, Success: true}
}

func (r *Raft) handleAppendEntriesResponse(peer uint64, resp *AppendEntriesResponse, nextIdx uint64) {
    // 处理追加条目响应
}

func min(a, b uint64) uint64 {
    if a < b {
        return a
    }
    return b
}
```
