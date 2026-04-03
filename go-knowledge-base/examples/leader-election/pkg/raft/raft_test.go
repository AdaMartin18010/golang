package raft

import (
	"testing"
	"time"
)

func TestNewRaft(t *testing.T) {
	config := &Config{
		ElectionTimeoutMin: 100 * time.Millisecond,
		ElectionTimeoutMax: 200 * time.Millisecond,
		HeartbeatInterval:  50 * time.Millisecond,
	}

	peers := []string{"node-1", "node-2", "node-3"}
	r := New("node-1", peers, config)

	if r == nil {
		t.Fatal("Expected Raft instance, got nil")
	}

	if r.State() != Follower {
		t.Errorf("Expected initial state to be Follower, got %s", r.State())
	}

	term, isLeader := r.GetState()
	if term != 0 {
		t.Errorf("Expected initial term to be 0, got %d", term)
	}
	if isLeader {
		t.Error("Expected not to be leader initially")
	}
}

func TestHandleRequestVote(t *testing.T) {
	config := DefaultConfig()
	peers := []string{"node-1", "node-2", "node-3"}
	r := New("node-1", peers, config)

	// Test voting for candidate with higher term
	args := &RequestVoteArgs{
		Term:         1,
		CandidateId:  "node-2",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	reply := r.HandleRequestVote(args)

	if !reply.VoteGranted {
		t.Error("Expected vote to be granted")
	}

	if reply.Term != 1 {
		t.Errorf("Expected term to be 1, got %d", reply.Term)
	}

	// Test not voting for candidate with lower term
	args2 := &RequestVoteArgs{
		Term:         0,
		CandidateId:  "node-3",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	reply2 := r.HandleRequestVote(args2)

	if reply2.VoteGranted {
		t.Error("Expected vote to not be granted for lower term")
	}
}

func TestHandleAppendEntries(t *testing.T) {
	config := DefaultConfig()
	peers := []string{"node-1", "node-2", "node-3"}
	r := New("node-1", peers, config)

	// Test heartbeat from leader
	args := &AppendEntriesArgs{
		Term:         1,
		LeaderId:     "node-2",
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries:      []LogEntry{},
		LeaderCommit: 0,
	}

	reply := r.HandleAppendEntries(args)

	if !reply.Success {
		t.Error("Expected AppendEntries to succeed")
	}

	// Check that we converted to follower
	if r.State() != Follower {
		t.Errorf("Expected state to be Follower after AppendEntries, got %s", r.State())
	}
}

func TestLogOperations(t *testing.T) {
	config := DefaultConfig()
	peers := []string{"node-1", "node-2"}
	r := New("node-1", peers, config)

	// Test initial log state
	if r.lastLogIndex() != 0 {
		t.Errorf("Expected initial log index to be 0, got %d", r.lastLogIndex())
	}

	if r.lastLogTerm() != 0 {
		t.Errorf("Expected initial log term to be 0, got %d", r.lastLogTerm())
	}
}

func BenchmarkHandleRequestVote(b *testing.B) {
	config := DefaultConfig()
	peers := []string{"node-1", "node-2", "node-3"}
	r := New("node-1", peers, config)

	args := &RequestVoteArgs{
		Term:         1,
		CandidateId:  "node-2",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.HandleRequestVote(args)
	}
}

func BenchmarkHandleAppendEntries(b *testing.B) {
	config := DefaultConfig()
	peers := []string{"node-1", "node-2", "node-3"}
	r := New("node-1", peers, config)

	args := &AppendEntriesArgs{
		Term:         1,
		LeaderId:     "node-2",
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries:      []LogEntry{},
		LeaderCommit: 0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.HandleAppendEntries(args)
	}
}
