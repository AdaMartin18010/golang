# TS-038: MongoDB 8 Replication Algorithm - S-Level Technical Reference

**Version:** MongoDB 8.0
**Status:** S-Level (Expert/Architectural)
**Last Updated:** 2026-04-03
**Classification:** Distributed Systems / Replication / Consensus

---

## 1. Executive Summary

MongoDB 8 introduces significant enhancements to its replication protocol, including an optimized Raft-based consensus implementation, improved majority acknowledgment mechanisms, and advanced transaction support across sharded clusters. This document provides comprehensive technical analysis of MongoDB 8's replication architecture, consensus algorithms, and operational characteristics.

---

## 2. Replication Architecture Overview

### 2.1 Replica Set Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB 8 Replica Set Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         Replica Set (rs0)                              │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                         Primary Node                             │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │  │  │
│  │  │  │   Client     │  │ Replication  │  │  Oplog       │           │  │  │
│  │  │  │   Operations │  │   Manager    │  │  Buffer      │           │  │  │
│  │  │  │              │  │              │  │              │           │  │  │
│  │  │  │  Writes ────▶│  │  Heartbeats  │  │  Entries     │           │  │  │
│  │  │  │  Reads  ◀────│  │  Elections   │  │  (capped)    │           │  │  │
│  │  │  │              │  │  Sync Source │  │              │           │  │  │
│  │  │  └──────────────┘  └──────┬───────┘  └──────┬───────┘           │  │  │
│  │  │                           │                  │                    │  │  │
│  │  │                           └──────────────────┘                    │  │  │
│  │  │                                           │                      │  │  │
│  │  │              ┌────────────────────────────┼──────────────────┐  │  │  │
│  │  │              │                            ▼                  │  │  │  │
│  │  │              │    ┌──────────────────────────────────────┐  │  │  │  │
│  │  │              │    │       Replication Protocol           │  │  │  │  │
│  │  │              │    │  ┌─────────┐  ┌─────────┐  ┌───────┐ │  │  │  │  │
│  │  │              │    │  │ Raft    │  │ Log     │  │ State │ │  │  │  │  │
│  │  │              │    │  │ Module  │  │ Manager │  │ Machine│ │  │  │  │  │
│  │  │              │    │  └────┬────┘  └────┬────┘  └───┬───┘ │  │  │  │  │
│  │  │              │    │       └────────────┴───────────┘     │  │  │  │  │
│  │  │              │    └──────────────────────────────────────┘  │  │  │  │
│  │  │              │                                              │  │  │  │
│  │  │              │    oplog replication to secondaries          │  │  │  │
│  │  │              └──────────────────────────────────────────────┘  │  │  │
│  │  │                                                                  │  │  │
│  │  └──────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                         │  │
│  │  ┌──────────────────────────┐  ┌──────────────────────────┐            │  │
│  │  │      Secondary 1         │  │      Secondary 2         │            │  │
│  │  │  (Priority: 2)           │  │  (Priority: 1)           │            │  │
│  │  │                          │  │                          │            │  │
│  │  │  ┌──────────────────┐   │  │  ┌──────────────────┐   │            │  │
│  │  │  │ Replication      │   │  │  │ Replication      │   │            │  │
│  │  │  │ - Apply oplog    │   │  │  │ - Apply oplog    │   │            │  │
│  │  │  │ - Heartbeat      │   │  │  │ - Heartbeat      │   │            │  │
│  │  │  │ - Vote in        │   │  │  │  - Vote in       │   │            │  │
│  │  │  │   elections      │   │  │  │   elections      │   │            │  │
│  │  │  └──────────────────┘   │  │  └──────────────────┘   │            │  │
│  │  │                          │  │                          │            │  │
│  │  │  Sync Source: PRIMARY   │  │  Sync Source: SECONDARY1 │            │  │
│  │  │  Optime: ts:100         │  │  Optime: ts:98           │            │  │
│  │  └──────────────────────────┘  └──────────────────────────┘            │  │
│  │                                                                         │  │
│  │  ┌──────────────────────────────────────────────────────────────────┐  │  │
│  │  │                      Arbiter Node                               │  │  │
│  │  │  (No data, participates in elections)                           │  │  │
│  │  └──────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                         │  │
│  │  Write Concern: majority                                                │  │
│  │  Read Concern: majority / linearizable / available                      │  │
│  │  Election Timeout: 10 seconds                                           │  │
│  │  Heartbeat Interval: 2 seconds                                          │  │
│  │                                                                         │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Oplog Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Oplog (Operations Log) Structure                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Oplog Entry Document:                                                       │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │ {                                                                      │   │
│  │   "ts": Timestamp(1699123456, 1),    // Operation timestamp            │   │
│  │   "t": 5,                            // Term number (Raft)             │   │
│  │   "h": NumberLong("-1234567890"),   // Hash of the oplog entry         │   │
│  │   "v": 2,                            // Version of oplog format        │   │
│  │   "op": "i",                         // Operation type:                │   │
│  │                                       // i=insert, u=update,           │   │
│  │                                       // d=delete, c=cmd, n=noop       │   │
│  │   "ns": "db.collection",             // Namespace                      │   │
│  │   "ui": UUID("..."),                 // Collection UUID                │   │
│  │   "o": {                             // Operation document             │   │
│  │     "_id": ObjectId("..."),                                         │   │
│  │     "field1": "value1",                                             │   │
│  │     "field2": "value2"                                              │   │
│  │   },                                                                   │   │
│  │   "o2": {                            // Update selector (for updates)  │   │
│  │     "_id": ObjectId("...")                                          │   │
│  │   },                                                                   │   │
│  │   "wall": ISODate("2023-11-04T12:30:56.123Z"),  // Wall clock time    │   │
│  │   "stmtId": 0,                       // Statement ID (transactions)    │   │
│  │   "prevOpTime": {                    // Previous operation timestamp    │   │
│  │     "ts": Timestamp(1699123456, 0),                                 │   │
│  │     "t": 5                                                          │   │
│  │   },                                                                   │   │
│  │   "lsid": {                          // Logical session ID             │   │
│  │     "id": UUID("..."),                                              │   │
│  │     "uid": BinData(0, "...")                                        │   │
│  │   },                                                                   │   │
│  │   "txnNumber": 1,                    // Transaction number             │   │
│  │   "authInfo": { ... }                // Authentication info (v8 new)   │   │
│  │ }                                                                      │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Oplog Capped Collection Properties:                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  • Stored in local.oplog.rs                                            │   │
│  │  • Capped collection with configurable size                            │   │
│  │    (default: 5% of free disk space, min 990MB)                         │   │
│  │  • Idempotent operations for replay safety                             │   │
│  │  • Natural order = insert order = chronological order                  │   │
│  │  • Entries replicated in order, cannot have gaps                        │   │
│  │  • Secondary truncate after: timestamp of oldest backup                 │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Raft Consensus Protocol Implementation

### 3.1 State Machine

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Raft State Machine                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                         States                                         │   │
│  │                                                                        │   │
│  │  ┌───────────────┐     ┌───────────────┐     ┌───────────────┐       │   │
│  │  │   FOLLOWER    │◀───▶│  CANDIDATE    │◀───▶│   PRIMARY     │       │   │
│  │  │               │     │               │     │  (LEADER)     │       │   │
│  │  │ • Replicates  │     │ • Request     │     │ • Processes   │       │   │
│  │  │   from primary│     │   votes       │     │   writes      │       │   │
│  │  │ • Votes in    │     │ • Increment   │     │ • Replicates  │       │   │
│  │  │   elections   │     │   term        │     │   to          │       │   │
│  │  │ • Handles     │     │ • Timeout if  │     │   secondaries │       │   │
│  │  │   heartbeats  │     │   no majority │     │ • Heartbeats  │       │   │
│  │  └───────┬───────┘     └───────┬───────┘     └───────┬───────┘       │   │
│  │          │                     │                     │                │   │
│  │          │                     │                     │                │   │
│  │          └─────────────────────┴─────────────────────┘                │   │
│  │                            │                                          │   │
│  │                            ▼                                          │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │   │
│  │  │                      Transitions                                 │  │   │
│  │  │                                                                  │  │   │
│  │  │  FOLLOWER ──[electionTimeout]──▶ CANDIDATE                      │  │   │
│  │  │                                                                  │  │   │
│  │  │  CANDIDATE ──[majority votes]──▶ PRIMARY                        │  │   │
│  │  │                                                                  │  │   │
│  │  │  CANDIDATE ──[discover higher term]──▶ FOLLOWER                 │  │   │
│  │  │                                                                  │  │   │
│  │  │  PRIMARY ──[stepDown command / higher term]──▶ FOLLOWER         │  │   │
│  │  │                                                                  │  │   │
│  │  │  ANY ──[higher term seen]──▶ FOLLOWER                           │  │   │
│  │  │                                                                  │  │   │
│  │  └─────────────────────────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  └────────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Election Algorithm

```
ALGORITHM RequestVote(candidateId, term, lastLogIndex, lastLogTerm):
    INPUT:  candidateId - ID of candidate requesting vote
            term - Candidate's term
            lastLogIndex - Index of candidate's last log entry
            lastLogTerm - Term of candidate's last log entry
    OUTPUT: voteGranted - Boolean
            currentTerm - For candidate to update itself

    1. // Reply false if term < currentTerm
       IF term < currentTerm:
          RETURN (voteGranted=false, currentTerm=currentTerm)

    2. // If term > currentTerm, convert to follower
       IF term > currentTerm:
          currentTerm ← term
          votedFor ← null
          state ← FOLLOWER
          resetElectionTimer()

    3. // Check if can vote for this candidate
       IF votedFor == null OR votedFor == candidateId:
          // Check if candidate's log is at least as up-to-date
          myLastLog ← getLastLogEntry()

          // Up-to-date check (Raft paper section 5.4.1)
          IF lastLogTerm > myLastLog.term:
             upToDate ← true
          ELSE IF lastLogTerm == myLastLog.term AND
                  lastLogIndex >= myLastLog.index:
             upToDate ← true
          ELSE:
             upToDate ← false

          IF upToDate:
             votedFor ← candidateId
             resetElectionTimer()
             RETURN (voteGranted=true, currentTerm=currentTerm)

    4. RETURN (voteGranted=false, currentTerm=currentTerm)

ALGORITHM StartElection():
    1. // Increment current term
       currentTerm ← currentTerm + 1

    2. // Transition to candidate state
       state ← CANDIDATE
       votedFor ← self.id
       votesReceived ← 1  // Vote for self

    3. // Reset election timer
       resetElectionTimer()

    4. // Send RequestVote RPCs to all other nodes
       FOR each server in configuration:
          IF server.id != self.id:
             lastLog ← getLastLogEntry()
             sendAsync(RequestVote, {
                term: currentTerm,
                candidateId: self.id,
                lastLogIndex: lastLog.index,
                lastLogTerm: lastLog.term
             }, server)

    5. // Wait for votes or timeout
       electionDeadline ← now() + electionTimeout
       WHILE now() < electionDeadline AND state == CANDIDATE:
          IF votesReceived > majoritySize:
             becomePrimary()
             RETURN
          sleep(10ms)

    6. // Election timeout - will retry on next timeout

FUNCTION becomePrimary():
    1. state ← PRIMARY
    2. FOR each server in configuration:
          nextIndex[server] ← lastLogIndex + 1
          matchIndex[server] ← 0

    3. // Start heartbeat goroutine
       startHeartbeatLoop()

    4. // Initialize write concern tracking
       writeConcernMajorityJournalDefault ← true
```

### 3.3 Log Replication

```
ALGORITHM AppendEntriesRPC(leaderTerm, leaderId, prevLogIndex,
                           prevLogTerm, entries[], leaderCommit):
    INPUT:  Standard Raft AppendEntries parameters
    OUTPUT: success - Boolean, term - currentTerm

    1. // Reply false if term < currentTerm
       IF leaderTerm < currentTerm:
          RETURN (success=false, term=currentTerm)

    2. // Convert to follower if term >= currentTerm
       IF leaderTerm >= currentTerm:
          currentTerm ← leaderTerm
          state ← FOLLOWER
          leaderId ← leaderId
          resetElectionTimer()

    3. // Reply false if log doesn't contain entry at prevLogIndex
       // with term matching prevLogTerm
       IF prevLogIndex > 0:
          IF prevLogIndex > lastLogIndex:
             RETURN (success=false, term=currentTerm)
          IF log[prevLogIndex].term != prevLogTerm:
             // Conflict detected - truncate log
             truncateLogFrom(prevLogIndex)
             RETURN (success=false, term=currentTerm)

    4. // Process entries
       FOR i ← 0 TO length(entries) - 1:
          entryIndex ← prevLogIndex + 1 + i

          IF entryIndex <= lastLogIndex AND
             log[entryIndex].term != entries[i].term:
             // Conflict - delete existing entry and all that follow
             truncateLogFrom(entryIndex)

          IF entryIndex > lastLogIndex:
             // Append new entry
             appendLogEntry(entries[i])

    5. // Update commitIndex
       IF leaderCommit > commitIndex:
          commitIndex ← min(leaderCommit, lastLogIndex)
          applyCommittedEntries()

    6. RETURN (success=true, term=currentTerm)

ALGORITHM ReplicateLogEntry(entry):
    // Called by PRIMARY to replicate new entry

    1. // Append to local log first
       appendLogEntry(entry)

    2. // Send to all secondaries
       acks ← 1  // Count self
       ackMutex ← new Mutex()
       ackCV ← new ConditionVariable()

       FOR each server in configuration WHERE server != self:
          sendAsync(AppendEntries, {
             term: currentTerm,
             leaderId: self.id,
             prevLogIndex: nextIndex[server] - 1,
             prevLogTerm: log[nextIndex[server] - 1].term,
             entries: log[nextIndex[server]:],
             leaderCommit: commitIndex
          }, server, callback=(success) => {
             IF success:
                ackMutex.lock()
                acks ← acks + 1
                matchIndex[server] ← lastLogIndex
                nextIndex[server] ← lastLogIndex + 1
                ackCV.notify()
                ackMutex.unlock()
             ELSE:
                // Decrement nextIndex and retry
                nextIndex[server] ← max(1, nextIndex[server] - 1)
          })

    3. // Wait for majority acknowledgment
       ackMutex.lock()
       WHILE acks <= majoritySize AND
             (now() - startTime) < replicationTimeout:
          ackCV.wait(replicationTimeout)
       ackMutex.unlock()

    4. IF acks > majoritySize:
          commitIndex ← lastLogIndex
          applyCommittedEntries()
          RETURN COMMITTED
       ELSE:
          RETURN TIMEOUT
```

---

## 4. Write Concern and Read Concern

### 4.1 Write Concern Levels

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MongoDB Write Concern Levels                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  writeConcern: { w: <value>, j: <boolean>, wtimeout: <number> }      │   │
│  │                                                                        │   │
│  │  w values:                                                             │   │
│  │  ┌──────────────┬───────────────────────────────────────────────────┐  │   │
│  │  │ w: 0         │ No acknowledgment (fire-and-forget)              │  │   │
│  │  │ w: 1         │ Acknowledged by primary only (default)           │  │   │
│  │  │ w: "majority"│ Acknowledged by majority of voting members       │  │   │
│  │  │ w: <number>  │ Acknowledged by specific number of members       │  │   │
│  │  │ w: <tag>     │ Acknowledged by members matching tag set         │  │   │
│  │  └──────────────┴───────────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  │  j: true   - Wait for journal commit (fsync to disk)                  │   │
│  │  wtimeout  - Timeout in milliseconds                                   │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Write Concern Implementation:                                               │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  PRIMARY          SECONDARY 1         SECONDARY 2                      │   │
│  │     │                  │                   │                           │   │
│  │     │  1. Write        │                   │                           │   │
│  │     │─────────────────▶│                   │                           │   │
│  │     │                  │                   │                           │   │
│  │     │  2. Apply        │                   │                           │   │
│  │     │◀─────────────────│                   │                           │   │
│  │     │                  │                   │                           │   │
│  │     │  3. Write        │                   │                           │   │
│  │     │─────────────────────────────────────▶│                           │   │
│  │     │                  │                   │                           │   │
│  │     │  4. Apply        │                   │                           │   │
│  │     │◀─────────────────────────────────────│                           │   │
│  │     │                  │                   │                           │   │
│  │     │  5. Majority ACK │                   │                           │   │
│  │     │◀─────────────────┼───────────────────┤                           │   │
│  │     │                  │                   │                           │   │
│  │     │  6. Response     │                   │                           │   │
│  │     │─────────────────▶│                   │                           │   │
│  │     │ (to client)      │                   │                           │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Read Concern Levels

| Level | Description | Use Case | Guarantees |
|-------|-------------|----------|------------|
| `local` | Return latest data on queried node | Default, max performance | No guarantees across nodes |
| `available` | Like local but can read from secondaries | Sharded clusters | Stale reads possible |
| `majority` | Return data acknowledged by majority | Strong consistency | No dirty reads, monotonic |
| `linearizable` | Read after write durability | Critical reads | Full linearizability |
| `snapshot` | Multi-document transaction isolation | Transactions | Repeatable reads |

---

## 5. Performance Benchmarks

### 5.1 Replication Latency

| Scenario | Avg Latency | P99 Latency | Throughput |
|----------|-------------|-------------|------------|
| Local network (w:1) | 0.5ms | 2ms | 50K ops/s |
| Local network (w:majority) | 2ms | 8ms | 15K ops/s |
| Cross-AZ (w:1) | 1.2ms | 5ms | 40K ops/s |
| Cross-AZ (w:majority) | 8ms | 25ms | 8K ops/s |
| Cross-region (w:1) | 15ms | 45ms | 25K ops/s |
| Cross-region (w:majority) | 120ms | 300ms | 2K ops/s |

### 5.2 Election Timing

| Scenario | Election Time | Data Unavailable |
|----------|---------------|------------------|
| Clean shutdown | 2-5s | 0s (planned) |
| Network partition | 10-12s | 10-12s |
| Primary crash | 10-15s | 10-15s |
| Priority takeover | 1-2s | 1-2s |
| Step down (rs.stepDown) | 1-2s | 1-2s |

---

## 6. References

1. **MongoDB Replication Documentation**
   - URL: <https://www.mongodb.com/docs/manual/replication/>

2. **Raft Consensus Paper**
   - URL: <https://raft.github.io/raft.pdf>

3. **MongoDB Architecture Guide**
   - URL: <https://www.mongodb.com/docs/manual/core/replica-set-internals/>

---

*Document generated for S-Level technical reference.*
