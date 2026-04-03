package raft

import (
	"sync"
	"time"
)

// State represents the Raft node state
type State int

const (
	Follower State = iota
	Candidate
	Leader
)

func (s State) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

// LogEntry represents a single log entry
type LogEntry struct {
	Index   int
	Term    int
	Command interface{}
}

// ApplyMsg is sent to application when log entry is committed
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

// Config holds Raft configuration
type Config struct {
	ElectionTimeoutMin  time.Duration
	ElectionTimeoutMax  time.Duration
	HeartbeatInterval   time.Duration
	MaxLogEntriesPerMsg int
	SnapshotThreshold   int
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		ElectionTimeoutMin:  150 * time.Millisecond,
		ElectionTimeoutMax:  300 * time.Millisecond,
		HeartbeatInterval:   50 * time.Millisecond,
		MaxLogEntriesPerMsg: 100,
		SnapshotThreshold:   10000,
	}
}

// Raft represents a Raft consensus node
type Raft struct {
	id      string
	peers   []string
	config  *Config

	state    State
	stateMu  sync.RWMutex

	currentTerm int
	votedFor    string
	log         []LogEntry
	persistMu   sync.Mutex

	commitIndex int
	lastApplied int

	nextIndex  map[string]int
	matchIndex map[string]int

	applyCh    chan ApplyMsg
	shutdownCh chan struct{}

	electionTimer  *time.Timer
	heartbeatTimer *time.Timer
}

// New creates a new Raft node
func New(id string, peers []string, config *Config) *Raft {
	if config == nil {
		config = DefaultConfig()
	}

	r := &Raft{
		id:         id,
		peers:      peers,
		config:     config,
		state:      Follower,
		log:        make([]LogEntry, 0),
		nextIndex:  make(map[string]int),
		matchIndex: make(map[string]int),
		applyCh:    make(chan ApplyMsg, 100),
		shutdownCh: make(chan struct{}),
	}

	// Add dummy entry at index 0
	r.log = append(r.log, LogEntry{Term: 0})

	r.resetElectionTimer()
	go r.run()

	return r
}

// run is the main event loop
func (r *Raft) run() {
	for {
		select {
		case <-r.shutdownCh:
			return

		case <-r.electionTimer.C:
			r.handleElectionTimeout()

		case <-r.heartbeatTimer.C:
			if r.State() == Leader {
				r.broadcastHeartbeat()
				r.resetHeartbeatTimer()
			}
		}
	}
}

// State returns the current state
func (r *Raft) State() State {
	r.stateMu.RLock()
	defer r.stateMu.RUnlock()
	return r.state
}

// IsLeader returns true if this node is the leader
func (r *Raft) IsLeader() bool {
	return r.State() == Leader
}

// GetState returns current term and whether this node is leader
func (r *Raft) GetState() (int, bool) {
	r.persistMu.Lock()
	defer r.persistMu.Unlock()
	return r.currentTerm, r.state == Leader
}

// lastLogIndex returns the index of the last log entry
func (r *Raft) lastLogIndex() int {
	return len(r.log) - 1
}

// lastLogTerm returns the term of the last log entry
func (r *Raft) lastLogTerm() int {
	if len(r.log) == 0 {
		return 0
	}
	return r.log[len(r.log)-1].Term
}

// resetElectionTimer resets the election timer with random timeout
func (r *Raft) resetElectionTimer() {
	if r.electionTimer != nil {
		r.electionTimer.Stop()
	}

	timeout := r.config.ElectionTimeoutMin +
		time.Duration(randInt(int(r.config.ElectionTimeoutMax-r.config.ElectionTimeoutMin)))

	r.electionTimer = time.NewTimer(timeout)
}

// resetHeartbeatTimer resets the heartbeat timer
func (r *Raft) resetHeartbeatTimer() {
	if r.heartbeatTimer != nil {
		r.heartbeatTimer.Stop()
	}
	r.heartbeatTimer = time.NewTimer(r.config.HeartbeatInterval)
}

// handleElectionTimeout handles election timeout
func (r *Raft) handleElectionTimeout() {
	r.persistMu.Lock()
	defer r.persistMu.Unlock()

	// Become candidate
	r.state = Candidate
	r.currentTerm++
	r.votedFor = r.id

	// Reset election timer
	r.resetElectionTimer()

	// In a real implementation, send RequestVote RPCs here
}

// broadcastHeartbeat sends heartbeats to all peers
func (r *Raft) broadcastHeartbeat() {
	// In a real implementation, send AppendEntries RPCs here
}

// randInt returns a random int between 0 and n
func randInt(n int) int {
	// Simple implementation - in production use crypto/rand
	return int(time.Now().UnixNano() % int64(n))
}

// RequestVoteArgs is the argument for RequestVote RPC
type RequestVoteArgs struct {
	Term         int
	CandidateId  string
	LastLogIndex int
	LastLogTerm  int
}

// RequestVoteReply is the reply for RequestVote RPC
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// AppendEntriesArgs is the argument for AppendEntries RPC
type AppendEntriesArgs struct {
	Term         int
	LeaderId     string
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

// AppendEntriesReply is the reply for AppendEntries RPC
type AppendEntriesReply struct {
	Term          int
	Success       bool
	ConflictIndex int
	ConflictTerm  int
}

// HandleRequestVote handles vote requests
func (r *Raft) HandleRequestVote(args *RequestVoteArgs) *RequestVoteReply {
	r.persistMu.Lock()
	defer r.persistMu.Unlock()

	reply := &RequestVoteReply{
		Term:        r.currentTerm,
		VoteGranted: false,
	}

	// Reply false if term < currentTerm
	if args.Term < r.currentTerm {
		return reply
	}

	// If term > currentTerm, update and convert to follower
	if args.Term > r.currentTerm {
		r.currentTerm = args.Term
		r.state = Follower
		r.votedFor = ""
	}

	// Check if we can vote for this candidate
	if r.votedFor == "" || r.votedFor == args.CandidateId {
		// Check if candidate's log is at least as up-to-date
		if r.isLogUpToDate(args.LastLogIndex, args.LastLogTerm) {
			r.votedFor = args.CandidateId
			reply.VoteGranted = true
			r.resetElectionTimer()
		}
	}

	return reply
}

// isLogUpToDate checks if candidate's log is at least as up-to-date
func (r *Raft) isLogUpToDate(lastLogIndex, lastLogTerm int) bool {
	myLastTerm := r.lastLogTerm()
	if lastLogTerm != myLastTerm {
		return lastLogTerm > myLastTerm
	}
	return lastLogIndex >= r.lastLogIndex()
}

// HandleAppendEntries handles append entries requests
func (r *Raft) HandleAppendEntries(args *AppendEntriesArgs) *AppendEntriesReply {
	r.persistMu.Lock()
	defer r.persistMu.Unlock()

	reply := &AppendEntriesReply{
		Term: r.currentTerm,
	}

	// Reply false if term < currentTerm
	if args.Term < r.currentTerm {
		return reply
	}

	// Convert to follower if term >= currentTerm
	if args.Term > r.currentTerm {
		r.currentTerm = args.Term
		r.votedFor = ""
	}

	r.state = Follower
	r.resetElectionTimer()

	// Check log consistency
	if args.PrevLogIndex >= len(r.log) {
		reply.ConflictIndex = len(r.log)
		return reply
	}

	if r.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.ConflictTerm = r.log[args.PrevLogIndex].Term
		// Find first index with this term
		for i := args.PrevLogIndex; i >= 0; i-- {
			if r.log[i].Term != reply.ConflictTerm {
				reply.ConflictIndex = i + 1
				break
			}
		}
		return reply
	}

	// Append new entries
	if len(args.Entries) > 0 {
		// Remove conflicting entries and append new ones
		for i, entry := range args.Entries {
			idx := args.PrevLogIndex + 1 + i
			if idx < len(r.log) {
				if r.log[idx].Term != entry.Term {
					r.log = r.log[:idx]
				}
			}
			if idx >= len(r.log) {
				r.log = append(r.log, entry)
			}
		}
	}

	// Update commit index
	if args.LeaderCommit > r.commitIndex {
		r.commitIndex = min(args.LeaderCommit, r.lastLogIndex())
	}

	reply.Success = true
	return reply
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
