# TS-039: TiDB 8.5 Distributed Transactions - S-Level Technical Reference

**Version:** TiDB 8.5
**Status:** S-Level (Expert/Architectural)
**Last Updated:** 2026-04-03
**Classification:** Distributed Databases / Transactions / Spanner-like Systems

---

## 1. Executive Summary

TiDB 8.5 introduces significant enhancements to its distributed transaction engine, including the optimized Percolator-based 2PC protocol, improved Timestamp Oracle (TSO) allocation, and the new partitioned Raft engine. This document provides comprehensive technical analysis of TiDB's transaction architecture, distributed consensus mechanisms, and performance optimization strategies.

---

## 2. TiDB Architecture Overview

### 2.1 System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       TiDB 8.5 Cluster Architecture                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         TiDB Layer (SQL)                               │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │  │
│  │  │   TiDB       │  │   TiDB       │  │   TiDB       │                │  │
│  │  │   Server 1   │  │   Server 2   │  │   Server N   │  Stateless     │  │
│  │  │              │  │              │  │              │  Compute       │  │
│  │  │ • SQL Parse  │  │ • SQL Parse  │  │ • SQL Parse  │  Nodes         │  │
│  │  │ • Optimizer  │  │ • Optimizer  │  │ • Optimizer  │                │  │
│  │  │ • Execution  │  │ • Execution  │  │ • Execution  │                │  │
│  │  │ • Txn Coord  │  │ • Txn Coord  │  │ • Txn Coord  │                │  │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                │  │
│  │         │                 │                 │                         │  │
│  │         └─────────────────┴─────────────────┘                         │  │
│  │                           │                                           │  │
│  │                           ▼                                           │  │
│  │                    Load Balancer (LVS/HAProxy/F5)                     │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                                    ▼                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        PD (Placement Driver)                           │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │  │
│  │  │     PD       │  │     PD       │  │     PD       │                │  │
│  │  │   Server 1   │  │   Server 2   │  │   Server 3   │                │  │
│  │  │   (Leader)   │  │  (Follower)  │  │  (Follower)  │                │  │
│  │  │              │  │              │  │              │                │  │
│  │  │ • TSO Alloc  │  │ • TSO Sync   │  │ • TSO Sync   │                │  │
│  │  │ • Metadata   │  │ • Metadata   │  │ • Metadata   │                │  │
│  │  │ • Scheduling │  │ • HA         │  │ • HA         │                │  │
│  │  └──────┬───────┘  └──────────────┘  └──────────────┘                │  │
│  │         │                                                             │  │
│  │         │  TSO: 46-bit physical + 18-bit logical = 64-bit timestamp   │  │
│  │         │                                                             │  │
│  └─────────┼─────────────────────────────────────────────────────────────┘  │
│            │                                                                 │
│            ▼                                                                 │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    TiKV Layer (Distributed Storage)                    │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    TiKV Node 1                                   │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │  │  │
│  │  │  │   Region 1   │  │   Region 4   │  │   Region 7   │          │  │  │
│  │  │  │   (Leader)   │  │   (Leader)   │  │  (Follower)  │          │  │  │
│  │  │  │              │  │              │  │              │          │  │  │
│  │  │  │  ┌────────┐  │  │  ┌────────┐  │  │  ┌────────┐  │          │  │  │
│  │  │  │  │  Raft  │  │  │  │  Raft  │  │  │  │  Raft  │  │          │  │  │
│  │  │  │  │  State │  │  │  │  State │  │  │  │  State │  │          │  │  │
│  │  │  │  │Machine │  │  │  │Machine │  │  │  │Machine │  │          │  │  │
│  │  │  │  └────────┘  │  │  └────────┘  │  │  └────────┘  │          │  │  │
│  │  │  │              │  │              │  │              │          │  │  │
│  │  │  │  ┌────────┐  │  │  ┌────────┐  │  │  ┌────────┐  │          │  │  │
│  │  │  │  │RocksDB │  │  │  │RocksDB │  │  │  │RocksDB │  │          │  │  │
│  │  │  │  │  LSM   │  │  │  │  LSM   │  │  │  │  LSM   │  │          │  │  │
│  │  │  │  └────────┘  │  │  └────────┘  │  │  └────────┘  │          │  │  │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘          │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    TiKV Node 2                                   │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │  │  │
│  │  │  │   Region 1   │  │   Region 2   │  │   Region 5   │          │  │  │
│  │  │  │  (Follower)  │  │   (Leader)   │  │   (Leader)   │          │  │  │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘          │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    TiKV Node 3                                   │  │  │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │  │  │
│  │  │  │   Region 1   │  │   Region 2   │  │   Region 3   │          │  │  │
│  │  │  │  (Follower)  │  │  (Follower)  │  │   (Leader)   │          │  │  │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘          │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Region: 96MB default, Raft consensus per region                      │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Region and Raft Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TiKV Region and Raft Implementation                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Region Key Range Partitioning:                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Key Space:                                                            │   │
│  │  ┌─────────┬─────────┬─────────┬─────────┬─────────┬─────────┐       │   │
│  │  │    |    │    |    │    |    │    |    │    |    │    |    │       │   │
│  │  │ Region 1│ Region 2│ Region 3│ Region 4│ Region 5│ Region N│       │   │
│  │  │         │         │         │         │         │         │       │   │
│  │  │ [ - , a)│ [a, b)  │ [b, c)  │ [c, d)  │ [d, e)  │ ...     │       │   │
│  │  │         │         │         │         │         │         │       │   │
│  │  │ 96MB    │ 96MB    │ 96MB    │ 96MB    │ 96MB    │ 96MB    │       │   │
│  │  └─────────┴─────────┴─────────┴─────────┴─────────┴─────────┘       │   │
│  │                                                                        │   │
│  │  Split happens when region size exceeds threshold                      │   │
│  │  Merge happens when adjacent regions are small                         │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Raft Consensus within Region:                                               │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │                   Region Leader (TiKV Node 1)                          │   │
│  │                         │                                              │   │
│  │        ┌────────────────┼────────────────┐                             │   │
│  │        │                │                │                             │   │
│  │        ▼                ▼                ▼                             │   │
│  │   ┌─────────┐      ┌─────────┐      ┌─────────┐                       │   │
│  │   │ Follower│◀────▶│ Follower│      │ Learner │  (optional)           │   │
│  │   │Node 2   │      │Node 3   │      │Node 4   │  (replication only)   │   │
│  │   └─────────┘      └─────────┘      └─────────┘                       │   │
│  │                                                                        │   │
│  │   Replication Flow:                                                    │   │
│  │   1. Client write ──▶ Leader                                           │   │
│  │   2. Leader appends to local Raft log                                  │   │
│  │   3. Leader sends AppendEntries RPC to followers                       │   │
│  │   4. Followers append to local log, send Ack                           │   │
│  │   5. Leader commits when majority (2/3) acks                           │   │
│  │   6. Leader applies to RocksDB, notifies client                        │   │
│  │   7. Leader sends ApplyMsg to followers                                │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Distributed Transaction Protocol

### 3.1 Percolator-based 2PC

TiDB uses an optimized Percolator transaction model with the following phases:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TiDB Distributed Transaction (Percolator)                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Transaction Phases:                                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Prewrite Phase:                                                       │   │
│  │  ┌─────────┐     ┌─────────┐     ┌─────────┐                          │   │
│  │  │  TiDB   │────▶│  TiKV   │     │  TiKV   │                          │   │
│  │  │         │     │Region 1 │     │Region 2 │                          │   │
│  │  └────┬────┘     └────┬────┘     └────┬────┘                          │   │
│  │       │               │               │                                │   │
│  │       │  1. Prewrite  │               │                                │   │
│  │       │──────────────▶│               │                                │   │
│  │       │  (Lock CF)    │               │                                │   │
│  │       │               │               │                                │   │
│  │       │     OK        │               │                                │   │
│  │       │◀──────────────│               │                                │   │
│  │       │               │               │                                │   │
│  │       │  2. Prewrite  │               │                                │   │
│  │       │──────────────────────────────▶│                                │   │
│  │       │               │               │                                │   │
│  │       │               │     OK        │                                │   │
│  │       │◀──────────────────────────────│                                │   │
│  │       │               │               │                                │   │
│  │       │  [All regions locked]          │                                │   │
│  │       ▼               │               │                                │   │
│  │                                                                        │   │
│  │  Commit Phase (Primary Key):                                           │   │
│  │  ┌─────────┐     ┌─────────┐                                           │   │
│  │  │  TiDB   │────▶│  TiKV   │                                           │   │
│  │  │         │     │Region 1 │  (Primary)                                │   │
│  │  └────┬────┘     └────┬────┘                                           │   │
│  │       │               │                                                │   │
│  │       │  3. Commit    │                                                │   │
│  │       │──────────────▶│                                                │   │
│  │       │  (Write CF +  │                                                │   │
│  │       │   Remove Lock)│                                                │   │
│  │       │               │                                                │   │
│  │       │     OK        │                                                │   │
│  │       │◀──────────────│                                                │   │
│  │       │               │                                                │   │
│  │       │  [Transaction committed]                                       │   │
│  │       ▼               │                                                │   │
│  │                                                                        │   │
│  │  Async Cleanup (Secondary Keys):                                       │   │
│  │  ┌─────────┐     ┌─────────┐                                           │   │
│  │  │  TiDB   │────▶│  TiKV   │                                           │   │
│  │  │         │     │Region 2 │  (Secondary)                              │   │
│  │  └────┬────┘     └────┬────┘                                           │   │
│  │       │               │                                                │   │
│  │       │  4. Cleanup   │                                                │   │
│  │       │  (Async)      │                                                │   │
│  │       │──────────────────────────────▶                                 │   │
│  │       │               │                                                │   │
│  └───────┼───────────────┼────────────────────────────────────────────────┘   │
│          │               │                                                    │
│          ▼               ▼                                                    │
│  Storage Format (RocksDB Column Families):                                    │
│  ┌───────────────────────────────────────────────────────────────────────┐    │
│  │  CF: Default (Data)    CF: Lock         CF: Write                     │    │
│  │  ┌───────────────┐     ┌───────────────┐  ┌───────────────┐          │    │
│  │  │ Key | Value   │     │ Key | LockInfo│  │Key | CommitTS |          │    │
│  │  │               │     │               │  │    StartTS    |          │    │
│  │  │ row_key_1     │     │ row_key_1     │  │               │          │    │
│  │  │ ───▶ data     │     │ ───▶ {primary,│  │ row_key_1     │          │    │
│  │  │               │     │      lock_ts, │  │ ───▶ {commit  │          │    │
│  │  │ row_key_2     │     │      ttl}     │  │        _ts,    │          │    │
│  │  │ ───▶ data     │     │               │  │        start_ts│          │    │
│  │  │               │     │ row_key_2     │  │ }             │          │    │
│  │  │               │     │ ───▶ {primary,│  │               │          │    │
│  │  │               │     │      ...}     │  │               │          │    │
│  │  └───────────────┘     └───────────────┘  └───────────────┘          │    │
│  └───────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Transaction Algorithm

```
ALGORITHM PercolatorPrewrite(keys, primary, start_ts):
    INPUT:  keys - List of keys to lock
            primary - Primary key for the transaction
            start_ts - Transaction start timestamp
    OUTPUT: success - Boolean

    1. FOR each key in keys:
       a. // Check for write conflicts
          latest_write ← CF_Write.GetLatest(key)
          IF latest_write.commit_ts > start_ts:
              RETURN WriteConflictError(key, latest_write)

       b. // Check for locks
          existing_lock ← CF_Lock.Get(key)
          IF existing_lock != null:
              IF existing_lock.is_expired():
                  // Try to cleanup stale lock
                  cleanup_result ← cleanupLock(key, existing_lock)
                  IF cleanup_result.failed:
                      RETURN LockConflictError(key, existing_lock)
              ELSE:
                  RETURN LockConflictError(key, existing_lock)

       c. // Write lock and data
          lock_info ← {
              primary: primary,
              start_ts: start_ts,
              ttl: current_time + 3s,
              type: IF key == primary THEN PrimaryLock ELSE SecondaryLock
          }

          // Atomic batch write
          batch ← new WriteBatch()
          batch.Put(CF_Lock, key, lock_info)
          batch.Put(CF_Default, key@{start_ts}, value)
          RocksDB.Write(batch)

    2. RETURN Success

ALGORITHM PercolatorCommit(primary, commit_ts):
    INPUT:  primary - Primary key
            commit_ts - Commit timestamp from PD
    OUTPUT: success - Boolean

    1. // Check primary lock
       lock ← CF_Lock.Get(primary)
       IF lock == null:
          // Already committed or rolled back
          write ← CF_Write.Get(primary)
          IF write != null AND write.start_ts == start_ts:
              RETURN AlreadyCommitted
          ELSE:
              RETURN Rollbacked

       IF lock.start_ts != start_ts:
          RETURN LockNotFoundError

    2. // Atomically commit primary
       batch ← new WriteBatch()
       batch.Delete(CF_Lock, primary)
       batch.Put(CF_Write, primary, {
           start_ts: start_ts,
           commit_ts: commit_ts,
           type: Put
       })
       batch.Put(CF_Default, primary@{commit_ts},
                 CF_Default.Get(primary@{start_ts}))
       RocksDB.Write(batch)

    3. // Asynchronously commit secondary keys
       FOR each secondary_key in secondary_keys:
          async_commit(secondary_key, start_ts, commit_ts)

    4. RETURN Success

ALGORITHM ReadWithSnapshot(key, start_ts):
    INPUT:  key - Key to read
            start_ts - Snapshot timestamp
    OUTPUT: value or null

    1. // Check for locks
       lock ← CF_Lock.Get(key)
       IF lock != null AND lock.start_ts <= start_ts:
          // There might be an uncommitted transaction
          IF lock.primary == key:
              // This is the primary, check if committed
              write ← CF_Write.Get(key)
              IF write != null AND write.start_ts == lock.start_ts:
                  // Committed, cleanup lock and return value
                  cleanup_async(key, lock)
                  RETURN CF_Default.Get(key@{write.commit_ts})
              ELSE:
                  // Not committed yet, check TTL
                  IF lock.is_expired():
                      rollback_primary(lock.primary)
                      RETURN read_before(key, lock.start_ts)
                  ELSE:
                      RETURN WaitOrError(LockNotResolved)
          ELSE:
              // Secondary lock, check primary
              primary_write ← CF_Write.Get(lock.primary)
              IF primary_write != null AND
                 primary_write.start_ts == lock.start_ts:
                  // Primary committed, commit secondary
                  commit_secondary(key, lock.start_ts,
                                  primary_write.commit_ts)
                  RETURN CF_Default.Get(key@{primary_write.commit_ts})
              ELSE:
                  // Primary not committed, treat as rollback
                  rollback_secondary(key, lock.start_ts)
                  RETURN read_before(key, lock.start_ts)

    2. // No lock, find most recent write before start_ts
       iterator ← CF_Write.NewIterator()
       iterator.SeekForPrev(key, start_ts)

       WHILE iterator.Valid() AND iterator.Key().key == key:
           write ← iterator.Value()
           IF write.type == Put:
              RETURN CF_Default.Get(key@{write.commit_ts})
           ELSE IF write.type == Delete:
              RETURN null
           iterator.Prev()

    3. RETURN null  // Key not found
```

---

## 4. Timestamp Oracle (TSO)

### 4.1 TSO Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    TiDB TSO (Timestamp Oracle) Architecture                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  TSO Timestamp Format (64 bits):                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  Bits 63-18 (46 bits)    │  Bits 17-0 (18 bits)                      │   │
│  │  Physical time (ms)      │  Logical counter                          │   │
│  │  (68 years since 2020)   │  (262,144 values per ms)                  │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  TSO Allocation Flow:                                                        │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  PD Leader                                    TiDB/TiKV Nodes          │   │
│  │     │                                              │                   │   │
│  │     │  1. Batch pre-allocate timestamps          │                   │   │
│  │     │     (e.g., 100,000 timestamps)             │                   │   │
│  │     │◀─────────────────────────────────────────────                   │   │
│  │     │                                              │                   │   │
│  │     │  2. Store in memory                          │                   │   │
│  │     │     ┌────────────────────────┐             │                   │   │
│  │     │     │  TSO Cache:            │             │                   │   │
│  │     │     │  physical: 1699123456  │             │                   │   │
│  │     │     │  logical_start: 0      │             │                   │   │
│  │     │     │  logical_end: 100000   │             │                   │   │
│  │     │     └────────────────────────┘             │                   │   │
│  │     │                                              │                   │   │
│  │     │  3. Individual timestamp requests          │                   │   │
│  │     │◀─────────────────────────────────────────────                   │   │
│  │     │     (fast, in-memory allocation)           │                   │   │
│  │     │     atomic.Add(&logical_count, 1)         │                   │   │
│  │     │                                              │                   │   │
│  │     │  4. Timestamp response                       │                   │   │
│  │     │─────────────────────────────────────────────▶                   │   │
│  │     │     timestamp = (physical << 18) | logical │                   │   │
│  │     │                                              │                   │   │
│  │                                                                        │   │
│  │  Synchronization between PD nodes (etcd):                                │   │
│  │  • Leader election uses etcd raft                                        │   │
│  │  • TSO allocation only from leader                                       │   │
│  │  • Followers forward requests to leader                                  │   │
│  │  • On failover: max(ts) + 1 to ensure monotonicity                       │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Performance Benchmarks

### 5.1 Transaction Performance

| Workload | Latency (avg) | Latency (p99) | Throughput |
|----------|---------------|---------------|------------|
| Point SELECT | 0.3ms | 1.2ms | 150K ops/s |
| Simple UPDATE | 2.1ms | 8ms | 45K txns/s |
| Distributed TX (2 regions) | 5.2ms | 18ms | 18K txns/s |
| Distributed TX (5 regions) | 12ms | 45ms | 8K txns/s |
| Large TX (1000 keys) | 150ms | 350ms | 120 txns/s |

### 5.2 TSO Performance

| Metric | Value |
|--------|-------|
| TSO allocation latency | ~0.1ms |
| TSO throughput | 1M+ timestamps/s |
| Batch size | 100,000 (configurable) |
| Clock drift tolerance | 100ms |

---

## 6. References

1. **TiDB Architecture Documentation**
   - URL: <https://docs.pingcap.com/tidb/stable/tidb-architecture>

2. **Percolator Paper**
   - URL: <https://research.google/pubs/pub36726/>

3. **TiKV Documentation**
   - URL: <https://tikv.org/docs/>

4. **Raft Consensus Paper**
   - URL: <https://raft.github.io/raft.pdf>

---

*Document generated for S-Level technical reference.*
