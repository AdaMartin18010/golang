# FT-034: Distributed System Failure Case Studies

> **维度**: Formal Theory | **级别**: S (15+ KB)
> **标签**: #distributed-systems #failure-cases #production-incidents #postmortem
> **权威来源**: Industry Postmortems, Academic Papers, Real-world Incidents

---

## Overview

This document contains 10 detailed production failure case studies from distributed systems, covering consensus failures, network partitions, clock skew issues, and more. Each case study includes incident description, root cause analysis, timeline, resolution steps, lessons learned, and prevention recommendations.

---

## Case Study 1: Split-Brain in Redis Cluster

### 1.1 Incident Description

**System**: E-commerce platform with Redis Cluster (6 master nodes, 6 replicas)
**Impact**: Data inconsistency, duplicate orders, inventory corruption
**Duration**: 47 minutes
**Date**: March 2024

During a routine network maintenance, a partial network partition occurred between two data centers. The Redis Cluster experienced a split-brain scenario where both sides of the partition elected new primary nodes, resulting in concurrent write operations on different masters for the same data shards.

### 1.2 Root Cause Analysis

```
Root Cause Chain:
1. Network maintenance triggered asymmetric network partition
2. Redis Cluster node timeout (cluster-node-timeout: 15s) too aggressive
3. Minimum master nodes check (cluster-require-full-coverage) disabled
4. Clients on both sides continued writing to their respective masters
5. No fencing mechanism to prevent dual-primary scenario
```

**Technical Details**:

- The partition isolated 3 masters in DC-A and 3 masters in DC-B
- Each side had quorum (50%+1) and formed independent clusters
- 1,247 keys were modified on both sides during the partition
- Key conflicts occurred in order IDs, inventory counters, and session data

### 1.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 02:14:32 | Network maintenance begins |
| 02:15:47 | First node timeout detected |
| 02:16:03 | DC-A side initiates failover |
| 02:16:08 | DC-B side initiates failover |
| 02:16:15 | Both sides complete failover with new primaries |
| 02:16:20 | Split-brain confirmed - dual primaries active |
| 02:20:00 | Customer complaints about duplicate orders |
| 02:45:00 | Data inconsistency detected in inventory system |
| 02:50:00 | Incident declared, traffic diverted to DC-A only |
| 03:01:15 | Network partition resolved |
| 03:30:00 | Manual conflict resolution completed |

### 1.4 Resolution Steps

```go
// Emergency fencing script to prevent split-brain
func emergencyFencing(dc string, masters []string) error {
    ctx := context.Background()

    for _, master := range masters {
        // Force demote to replica
        client := redis.NewClient(&redis.Options{Addr: master})

        // Set cluster-slave-no-one to 0
        err := client.ClusterReplicate(ctx, "no").Err()
        if err != nil {
            log.Printf("Failed to demote %s: %v", master, err)
        }

        // Enable read-only mode
        err = client.ConfigSet(ctx, "slave-read-only", "yes").Err()
        if err != nil {
            log.Printf("Failed to set read-only on %s: %v", master, err)
        }
    }

    return nil
}

// Data conflict detection
func detectConflicts(dcAData, dcBData map[string]string) []Conflict {
    var conflicts []Conflict

    for key, valueA := range dcAData {
        if valueB, exists := dcBData[key]; exists && valueA != valueB {
            conflicts = append(conflicts, Conflict{
                Key:      key,
                ValueA:   valueA,
                ValueB:   valueB,
                Strategy: determineResolutionStrategy(key),
            })
        }
    }

    return conflicts
}
```

### 1.5 Lessons Learned

1. **Quorum-based decisions are insufficient** for preventing split-brain in partitioned networks
2. **Fencing mechanisms** must be in place before failover completes
3. **Data conflict resolution** strategies should be predefined per data type
4. **Asymmetric partitions** are harder to detect than symmetric ones

### 1.6 Prevention Recommendations

```yaml
# Redis Cluster Configuration
cluster-require-full-coverage: yes
cluster-node-timeout: 30000  # Increased from 15s
cluster-replica-validity-factor: 10
cluster-migration-barrier: 1

# Sentinel Configuration for monitoring
sentinel:
  down-after-milliseconds: 30000
  parallel-syncs: 1
  failover-timeout: 180000

# External fencing with ZooKeeper
fencing:
  enabled: true
  provider: zookeeper
  lock-path: /redis/fencing
  session-timeout: 15000
```

---

## Case Study 2: Raft Consensus Liveness Failure

### 2.1 Incident Description

**System**: Distributed configuration store using Raft consensus (5-node cluster)
**Impact**: Configuration updates blocked, service registration failures
**Duration**: 2 hours 15 minutes
**Date**: January 2024

A Raft-based configuration management system experienced a liveness failure where the cluster could not elect a leader despite having a majority of nodes available. This resulted in a complete inability to update configurations or register new services.

### 2.2 Root Cause Analysis

```
Failure Scenario:
1. Node A (Leader) experienced GC pause exceeding election timeout
2. Node B and C started election with term T+1
3. Node A recovered and rejected votes (old term T)
4. Node D had clock skew +200ms, rejected vote requests as "stale"
5. Node B and C couldn't achieve quorum (only 2 votes)
6. Continuous election cycles exhausted network bandwidth
7. Pre-vote mechanism not implemented, causing disruptive elections
```

**Critical Factors**:

- Election timeout: 100-200ms (too short for production)
- No pre-vote phase in Raft implementation
- Clock skew detection disabled
- Network partition detection based on heartbeats only

### 2.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 14:23:10 | Leader Node A GC pause begins |
| 14:23:15 | Followers timeout, election starts |
| 14:23:16 | Node A GC completes, rejects vote requests |
| 14:23:20 | First election fails (split vote) |
| 14:23:25 | Second election cycle begins |
| 14:25:00 | Election storm detected (15+ elections in 2 minutes) |
| 14:30:00 | Manual intervention attempted |
| 15:00:00 | Clock skew on Node D detected |
| 15:20:00 | Node D restarted with NTP sync |
| 15:35:00 | Leader elected successfully |
| 16:38:10 | Service fully recovered |

### 2.4 Resolution Steps

```go
// Raft Pre-Vote Implementation
func (r *Raft) preVote() bool {
    // Ask peers if they would grant vote without incrementing term
    votes := 1 // Self vote

    for _, peer := range r.peers {
        resp, err := peer.RequestPreVote(&PreVoteRequest{
            Term:         r.currentTerm,
            CandidateId:  r.id,
            LastLogIndex: r.lastLogIndex,
            LastLogTerm:  r.lastLogTerm,
        })

        if err == nil && resp.VoteGranted {
            votes++
        }
    }

    // Only start real election if pre-vote succeeds
    return votes > len(r.peers)/2
}

// Clock skew detection
func detectClockSkew() error {
    ntpTime, err := queryNTP("pool.ntp.org")
    if err != nil {
        return err
    }

    localTime := time.Now()
    drift := localTime.Sub(ntpTime)

    if math.Abs(drift.Seconds()) > 0.1 { // 100ms threshold
        return fmt.Errorf("clock skew detected: %v", drift)
    }

    return nil
}
```

### 2.5 Lessons Learned

1. **Pre-vote is essential** in production Raft implementations
2. **Clock synchronization** is not optional for consensus systems
3. **Election timeouts** should be jittered and configurable per environment
4. **Leader stickiness** can reduce unnecessary elections

### 2.6 Prevention Recommendations

```go
// Production Raft Configuration
type RaftConfig struct {
    ElectionTimeout    time.Duration // 300-500ms
    HeartbeatInterval  time.Duration // 50ms
    MaxLogEntries      int           // 1000
    SnapshotInterval   time.Duration // 1 hour

    // Liveness improvements
    PreVoteEnabled     bool          // true
    CheckQuorum        bool          // true
    LeaderStickiness   time.Duration // 5s

    // Clock sync
    NTPCheckInterval   time.Duration // 1 minute
    MaxClockSkew       time.Duration // 50ms
}

// CheckQuorum mechanism
func (r *Raft) checkQuorum() {
    if !r.isLeader {
        return
    }

    acks := 0
    for _, peer := range r.peers {
        if time.Since(peer.lastAck) < r.electionTimeout {
            acks++
        }
    }

    if acks <= len(r.peers)/2 {
        r.stepDown()
    }
}
```

---

## Case Study 3: Cassandra Hinted Handoff Overload

### 3.1 Incident Description

**System**: Cassandra cluster (12 nodes, RF=3) storing time-series metrics
**Impact**: Cascading node failures, cluster unavailability
**Duration**: 3 hours 45 minutes
**Date**: February 2024

A network blip caused brief node disconnections. Cassandra's hinted handoff mechanism queued millions of write hints. When nodes reconnected, the aggressive hint replay overwhelmed the cluster, causing cascading failures and complete service unavailability.

### 3.2 Root Cause Analysis

```
Cascade Failure:
1. 30-second network blip disconnected 4 nodes
2. Coordinators stored 50M+ write hints in memory
3. Nodes reconnected, hint replay began at max throttle
4. Hint replay consumed all disk I/O capacity
5. Compaction couldn't keep up, SSTable count exploded
6. Read latency spiked, clients started retry storms
7. GC pressure caused nodes to drop out of gossip ring
8. Remaining nodes couldn't handle increased load
9. Complete cluster meltdown
```

### 3.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 09:15:00 | Network blip begins |
| 09:15:30 | Network recovers, 4 nodes marked DOWN |
| 09:15:35 | Hinted handoff activated |
| 09:45:00 | Hint queue reaches 50M entries |
| 10:00:00 | Nodes marked UP, hint replay starts |
| 10:05:00 | Disk I/O saturation detected |
| 10:15:00 | First node OOM due to hint buffer |
| 10:30:00 | Read latency > 30s |
| 11:00:00 | Retry storm from clients |
| 12:00:00 | Cascading node failures begin |
| 13:00:00 | Complete cluster unavailability |

### 3.4 Resolution Steps

```yaml
# Emergency cassandra.yaml modifications
hinted_handoff_enabled: false  # Disable temporarily
hinted_handoff_throttle_in_kb: 1024  # Reduced from 2048
max_hint_window_in_ms: 3600000  # 1 hour limit
hints_directory: /var/lib/cassandra/hints

# Clear hint files
# nodetool disablehandoff
# rm -rf /var/lib/cassandra/hints/*
# nodetool enablehandoff

# Rate limiting for replay
dynamic_snitch: false
read_request_timeout_in_ms: 5000
write_request_timeout_in_ms: 2000
```

### 3.5 Lessons Learned

1. **Hinted handoff is not free** - it can destabilize the cluster
2. **Throttle limits** must be environment-specific
3. **Hint expiration** should be shorter than default (3 hours)
4. **Circuit breakers** needed for replay mechanism

### 3.6 Prevention Recommendations

```go
// Adaptive hint throttling
type HintThrottler struct {
    baseRate      int // KB/s
    currentRate   int
    diskUtil      float64
    compactionLag int
}

func (ht *HintThrottler) adjustRate() {
    // Reduce rate if disk is busy
    if ht.diskUtil > 0.8 {
        ht.currentRate = ht.baseRate / 2
    } else if ht.diskUtil < 0.5 && ht.compactionLag < 10 {
        ht.currentRate = min(ht.currentRate*2, ht.baseRate*2)
    }
}

// Hint age-based filtering
func shouldReplayHint(hint *Hint) bool {
    age := time.Since(hint.timestamp)

    // Discard old hints for time-series data
    if hint.table == "metrics" && age > 15*time.Minute {
        return false
    }

    return age < maxHintAge
}
```

---

## Case Study 4: ZooKeeper Session Expiration Cascade

### 4.1 Incident Description

**System**: Microservices platform using ZooKeeper for service discovery (500+ services)
**Impact**: Mass service deregistration, thundering herd on recovery
**Duration**: 1 hour 20 minutes
**Date**: December 2023

A JVM GC pause on the ZooKeeper leader caused session expirations for hundreds of clients. The simultaneous reconnection attempt created a thundering herd, overwhelming the remaining ZK servers and causing additional session expirations.

### 4.2 Root Cause Analysis

```
Failure Cascade:
1. ZK Leader JVM heap exhaustion (Xmx: 4GB insufficient)
2. GC pause > session timeout (30s)
3. 350+ client sessions expired simultaneously
4. Clients received SESSION_EXPIRED event
5. All clients attempted immediate reconnection
6. ZK servers overwhelmed, more sessions dropped
7. Services deregistered from load balancers
8. Traffic routed to non-existent instances
```

### 4.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 16:45:00 | ZK Leader GC begins |
| 16:45:35 | GC pause exceeds session timeout |
| 16:45:36 | Mass session expiration (350 clients) |
| 16:45:37 | Clients begin reconnection |
| 16:46:00 | Thundering herd detected |
| 16:47:00 | Follower servers overwhelmed |
| 16:50:00 | Service discovery returns empty lists |
| 16:55:00 | 40% of services marked unhealthy |
| 17:00:00 | Manual intervention: rate limiting applied |
| 18:05:00 | All services re-registered |

### 4.4 Resolution Steps

```go
// Exponential backoff for ZK reconnection
type ZKReconnectBackoff struct {
    baseDelay    time.Duration
    maxDelay     time.Duration
    jitterFactor float64
    attempt      int
}

func (rb *ZKReconnectBackoff) nextDelay() time.Duration {
    // Exponential backoff with full jitter
    delay := rb.baseDelay * time.Duration(math.Pow(2, float64(rb.attempt)))
    if delay > rb.maxDelay {
        delay = rb.maxDelay
    }

    jitter := time.Duration(rand.Float64() * float64(delay) * rb.jitterFactor)
    rb.attempt++

    return delay + jitter
}

// Client-side rate limiting
func (c *ZKClient) reconnect() {
    // Check cluster health before reconnecting
    if !c.isClusterHealthy() {
        time.Sleep(c.backoff.nextDelay())
    }

    // Stagger reconnections
    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

    c.doConnect()
}
```

### 4.5 Lessons Learned

1. **ZK session timeouts** need to account for GC pauses
2. **JVM heap sizing** for ZK should include safety margins
3. **Reconnection storms** can be worse than the original failure
4. **Client backoff strategies** are critical

### 4.6 Prevention Recommendations

```yaml
# zoo.cfg optimizations
tickTime: 2000
initLimit: 10
syncLimit: 5
maxClientCnxns: 300

# Session timeout configuration
minSessionTimeout: 40000  # 2x tickTime * initLimit
maxSessionTimeout: 80000

# JVM settings
jvm:
  heap_size: 8g
  gc: G1GC
  g1_heap_region_size: 16m
  max_gc_pause_millis: 200
```

---

## Case Study 5: Kafka Consumer Group Rebalance Storm

### 5.1 Incident Description

**System**: Kafka cluster processing 2M messages/sec, 200 consumers
**Impact**: Message processing stopped, consumer lag grew to 50M messages
**Duration**: 45 minutes
**Date**: November 2023

A deployment triggered consumer group rebalances. The "stop-the-world" rebalance protocol caused cascading rebalances as consumers joined/left, resulting in a rebalance storm where no processing occurred for 45 minutes.

### 5.2 Root Cause Analysis

```
Rebalance Storm:
1. Deployment updated 50 consumers simultaneously
2. Consumers left group, triggering rebalance
3. Rebalance took 30s (200 consumers, 1000 partitions)
4. During rebalance, heartbeat timeouts occurred
5. More consumers considered dead, new rebalance triggered
6. This cycle repeated continuously
7. No consumer was stable enough to process messages
8. Consumer lag grew by 50M messages
```

### 5.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 11:00:00 | Deployment begins |
| 11:00:05 | First batch of consumers restart |
| 11:00:10 | First rebalance triggered |
| 11:00:40 | Rebalance completes, processing resumes |
| 11:01:00 | Second batch restart, rebalance #2 |
| 11:01:30 | Heartbeat timeouts during rebalance |
| 11:02:00 | Consumers removed from group |
| 11:02:05 | Rebalance storm begins |
| 11:15:00 | Consumer lag: 20M messages |
| 11:30:00 | Consumer lag: 50M messages |
| 11:45:00 | Rolling restart with staggered deployment |

### 5.4 Resolution Steps

```go
// Cooperative rebalancing with IncrementalRebalanceProtocol
config := kafka.ConfigMap{
    "bootstrap.servers":  "kafka:9092",
    "group.id":          "consumer-group",
    "partition.assignment.strategy": "cooperative-sticky",
    "session.timeout.ms": 10000,
    "heartbeat.interval.ms": 3000,
    "max.poll.interval.ms": 300000,
}

// Static membership to reduce rebalances
config["group.instance.id"] = getPodName() // Unique per pod
```

### 5.5 Lessons Learned

1. **Eager rebalancing** doesn't scale to large consumer groups
2. **Static membership** reduces unnecessary rebalances
3. **Staggered deployments** are essential for Kafka consumers
4. **Incremental cooperative rebalancing** should be the default

### 5.6 Prevention Recommendations

```yaml
# Consumer configuration
consumer:
  partition.assignment.strategy: org.apache.kafka.clients.consumer.CooperativeStickyAssignor
  session.timeout.ms: 10000
  heartbeat.interval.ms: 3000
  max.poll.interval.ms: 300000

  # Static membership
  group.instance.id: ${POD_NAME}

# Deployment strategy
deployment:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 10%
      maxUnavailable: 10%
```

---

## Case Study 6: etcd MVCC Database Bloat

### 6.1 Incident Description

**System**: Kubernetes cluster using etcd (3-node cluster, 8GB disk)
**Impact**: API server failures, cluster control plane down
**Duration**: 2 hours
**Date**: October 2023

etcd's MVCC store accumulated uncompacted revisions due to a bug in the compaction cron job. The database grew from 2GB to 8GB, causing disk exhaustion and cluster-wide control plane failures.

### 6.2 Root Cause Analysis

```
Failure Chain:
1. Automated compaction job failed silently (permission issue)
2. Revisions accumulated for 30 days
3. etcd db size reached 8GB limit (quota-backend-bytes)
4. Alarm raised: NOSPACE
5. etcd switched to read-only mode
6. Kubernetes API server writes failed
7. Pods couldn't be scheduled, services not updated
```

### 6.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 30 days ago | Compaction job last successful |
| 08:00:00 | etcd db size: 7.5GB |
| 08:15:00 | NOSPACE alarm triggered |
| 08:16:00 | etcd enters read-only mode |
| 08:20:00 | Pod scheduling failures begin |
| 08:30:00 | Kubernetes API errors spike |
| 09:00:00 | Manual compaction initiated |
| 10:00:00 | Compaction completed, db size: 500MB |
| 10:15:00 | Alarm disarmed, writes restored |

### 6.4 Resolution Steps

```bash
# Check etcd database size
ETCDCTL_API=3 etcdctl --endpoints=$ENDPOINTS endpoint status -w table

# Manual compaction
ETCDCTL_API=3 etcdctl --endpoints=$ENDPOINTS compaction $(ETCDCTL_API=3 etcdctl --endpoints=$ENDPOINTS get "" --prefix --keys-only | wc -l)

# Defragment to reclaim space
ETCDCTL_API=3 etcdctl --endpoints=$ENDPOINTS defrag

# Disarm alarm
ETCDCTL_API=3 etcdctl --endpoints=$ENDPOINTS alarm disarm
```

### 6.5 Lessons Learned

1. **Compaction is critical** for etcd maintenance
2. **Silent failures** in maintenance jobs are dangerous
3. **Monitoring db size** with alerts is essential
4. **Quota limits** should trigger alerts before alarms

### 6.6 Prevention Recommendations

```yaml
# etcd configuration
etcd:
  quota-backend-bytes: 8589934592  # 8GB
  auto-compaction-mode: periodic
  auto-compaction-retention: "1h"

# Monitoring alerts
alerts:
  - name: EtcdDatabaseSizeHigh
    expr: etcd_mvcc_db_total_size_in_bytes / etcd_server_quota_backend_bytes > 0.8
    for: 5m
    severity: warning

  - name: EtcdDatabaseSizeCritical
    expr: etcd_mvcc_db_total_size_in_bytes / etcd_server_quota_backend_bytes > 0.9
    for: 1m
    severity: critical
```

---

## Case Study 7: MongoDB Replica Set Stale Secondary

### 7.1 Incident Description

**System**: MongoDB replica set (1 primary, 2 secondaries) for user profiles
**Impact**: Read preference failures, query routing to stale data
**Duration**: 6 hours (undetected)
**Date**: September 2023

A secondary node fell behind replication due to disk issues. Applications using `readPreference: secondary` served 6-hour stale data, causing user profile updates to appear lost and cache inconsistencies.

### 7.2 Root Cause Analysis

```
Root Causes:
1. Secondary's disk performance degraded (SSD wear)
2. Replication lag grew gradually
3. Oplog window (24h) prevented recovery
4. No replication lag monitoring on read path
5. Driver continued routing reads to stale secondary
6. Stale data cached by application layer
```

### 7.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 00:00:00 | Secondary disk performance degrades |
| 02:00:00 | Replication lag: 2 hours |
| 08:00:00 | User complaints about profile changes not saving |
| 10:00:00 | Replication lag: 6 hours |
| 14:00:00 | Issue detected in monitoring dashboard |
| 14:30:00 | Secondary removed from replica set |
| 15:00:00 | New secondary provisioned |
| 16:00:00 | Initial sync completed |

### 7.4 Resolution Steps

```go
// Driver configuration with max staleness
clientOpts := options.Client().ApplyURI("mongodb://mongodb:27017").
    SetReadPreference(readpref.Secondary(readpref.WithMaxStaleness(90 * time.Second)))

// Custom server selector
func staleServerSelector(lagThreshold time.Duration) description.ServerSelectorFunc {
    return func(srv description.Topology, cs description.ConnectionState, candidates []description.Server) ([]description.Server, error) {
        var valid []description.Server
        for _, s := range candidates {
            if s.ReplicationStatus != description.ReplicaSetSecondary {
                valid = append(valid, s)
                continue
            }

            // Check replication lag
            if s.ReplicationLag <= lagThreshold {
                valid = append(valid, s)
            }
        }
        return valid, nil
    }
}
```

### 7.5 Lessons Learned

1. **Default read preferences** don't account for replication lag
2. **Secondary reads need staleness limits**
3. **Replication lag monitoring** must be actionable
4. **Application caching** amplifies stale data issues

### 7.6 Prevention Recommendations

```yaml
# MongoDB configuration
replication:
  replSetName: rs0
  secondaryIndexPrefetch: all

# Monitoring
metrics:
  replication_lag_seconds:
    warning: 10
    critical: 60

# Application driver settings
mongodb:
  uri: mongodb://mongodb:27017/?readPreference=secondaryMaxStalenessMS=90000
  maxStalenessSeconds: 90
```

---

## Case Study 8: Consul KV Store Inconsistency

### 8.1 Incident Description

**System**: Service mesh using Consul for configuration (5-server cluster)
**Impact**: Routing rules inconsistent, traffic routed to wrong services
**Duration**: 1 hour 10 minutes
**Date**: August 2023

Network partition created divergent Consul KV stores. When the partition healed, conflicting configuration values caused the service mesh to route traffic incorrectly, resulting in cross-service data exposure.

### 8.2 Root Cause Analysis

```
Failure Scenario:
1. Network partition: 3 nodes (DC-A) | 2 nodes (DC-B)
2. DC-A maintained quorum, continued accepting writes
3. DC-B entered read-only mode but cached old values
4. Configuration updated in DC-A (routing rules changed)
5. Network partition healed
6. Log replay caused temporary inconsistencies
7. Services in DC-B received stale routing rules
8. Traffic routed to incorrect upstream services
```

### 8.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 13:20:00 | Network partition occurs |
| 13:25:00 | Configuration update in DC-A |
| 13:30:00 | Routing rule change deployed |
| 13:35:00 | Network partition heals |
| 13:40:00 | First routing anomalies detected |
| 13:45:00 | Cross-service traffic identified |
| 14:00:00 | Emergency config rollback |
| 14:30:00 | All nodes synchronized |

### 8.4 Resolution Steps

```go
// Consul write with consistency check
func safeKVPut(client *api.Client, key string, value []byte, index uint64) error {
    kv := client.KV()

    // Use Check-And-Set to prevent overwriting
    p := &api.KVPair{
        Key:         key,
        Value:       value,
        ModifyIndex: index,
    }

    ok, _, err := kv.CAS(p, nil)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("CAS failed: key modified")
    }

    return nil
}

// Read with consistent mode
func consistentRead(client *api.Client, key string) (*api.KVPair, error) {
    kv := client.KV()

    // Use consistent=true for critical reads
    pair, _, err := kv.Get(key, &api.QueryOptions{
        Consistency: api.ConsistencyModeConsistent,
    })

    return pair, err
}
```

### 8.5 Lessons Learned

1. **Eventual consistency** can cause real issues during partitions
2. **CAS operations** prevent silent overwrites
3. **Consistent reads** are necessary for critical configuration
4. **Post-partition reconciliation** should be automated

### 8.6 Prevention Recommendations

```hcl
# Consul configuration
consul {
  performance {
    raft_multiplier = 1  # Faster consensus
  }

  acl {
    enabled = true
    default_policy = "deny"
  }
}

# Application configuration retrieval
config {
  consistency = "consistent"  # Not "default"
  max_stale = "0s"
  require_consistent = true
}
```

---

## Case Study 9: RabbitMQ Mirrored Queue Partition Handling

### 9.1 Incident Description

**System**: RabbitMQ cluster with mirrored queues (3 nodes) for order processing
**Impact**: Message loss, duplicate order processing
**Duration**: 2 hours 30 minutes
**Date**: July 2023

Network partition caused RabbitMQ cluster to split. The `pause_minority` mode paused nodes in the minority partition, but when the partition healed, unsynchronized messages caused duplicates and losses during queue synchronization.

### 9.2 Root Cause Analysis

```
Failure Sequence:
1. Network partition isolated 1 node from 2-node majority
2. Minority node paused (pause_minority mode)
3. Publishers continued to majority cluster
4. Consumers on minority node couldn't process
5. Some publishers failed over to minority node (DNS cache)
6. Messages queued on minority node (unmirrored)
7. Partition healed, minority node resumed
8. Queue synchronization: conflicting message IDs
9. Some messages lost, others duplicated
```

### 9.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 10:00:00 | Network partition detected |
| 10:00:05 | Minority node paused |
| 10:00:10 | Publishers fail over to majority |
| 10:05:00 | Some publishers stuck to minority (DNS TTL) |
| 10:15:00 | 500 messages queued on minority node |
| 10:30:00 | Network partition heals |
| 10:30:05 | Minority node resumes |
| 10:30:10 | Queue synchronization begins |
| 10:35:00 | Message inconsistencies detected |
| 12:30:00 | Manual reconciliation completed |

### 9.4 Resolution Steps

```erlang
%% RabbitMQ configuration for partition handling
[
  {rabbit, [
    {cluster_partition_handling, pause_minority},
    {queue_master_locator, <<"min-masters">>},
    {mirroring_metrics, true}
  ]},
  {rabbitmq_shovel, [
    {shovels, [
      {emergency_draining, [
        {sources, [
          {broker, "amqp://minority-node"}
        ]},
        {destinations, [
          {broker, "amqp://majority-node"}
        ]},
        {queue, <<">>},
        {prefetch_count, 1000},
        {publish_properties, [
          {delivery_mode, 2}
        ]}
      ]}
    ]}
  ]}
].
```

### 9.5 Lessons Learned

1. **Pause minority mode** is safer but not perfect
2. **Publisher failover** must be handled carefully
3. **Message deduplication** should be implemented
4. **Queue synchronization** can cause message loss

### 9.6 Prevention Recommendations

```go
// Idempotent consumer with deduplication
type IdempotentConsumer struct {
    store     *redis.Client
    ttl       time.Duration
    handler   func(Message) error
}

func (c *IdempotentConsumer) Consume(msg Message) error {
    key := fmt.Sprintf("msg:%s", msg.MessageId)

    // Check if already processed
    exists, err := c.store.SetNX(ctx, key, "1", c.ttl).Result()
    if err != nil {
        return err
    }

    if !exists {
        log.Printf("Duplicate message detected: %s", msg.MessageId)
        return nil // Acknowledge without processing
    }

    return c.handler(msg)
}
```

---

## Case Study 10: DynamoDB Global Table Conflict Resolution

### 10.1 Incident Description

**System**: Global application using DynamoDB Global Tables (3 regions)
**Impact**: Order status conflicts, inventory over-selling
**Duration**: 4 hours
**Date**: June 2023

Simultaneous writes to the same order record in different regions caused last-write-wins conflicts. The default conflict resolution resulted in incorrect order statuses, leading to inventory over-selling and fulfillment errors.

### 10.2 Root Cause Analysis

```
Conflict Scenario:
1. Order created in us-east-1
2. Order update initiated simultaneously in us-west-2 and eu-west-1
3. Both regions updated order status (PENDING → PROCESSING)
4. Different processing centers assigned in each region
5. DynamoDB replicated both updates
6. Last-write-wins based on timestamp
7. One update lost (but inventory already reserved)
8. Double inventory reservation for same order
```

### 10.3 Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 15:00:00 | Order created: ORDER-12345 |
| 15:05:00 | Concurrent updates in us-west-2 and eu-west-1 |
| 15:05:01 | Both updates applied locally |
| 15:05:05 | Replication begins |
| 15:05:10 | Conflict detected, resolved by timestamp |
| 15:06:00 | Inventory reserved in both regions |
| 15:10:00 | Order appears in both fulfillment queues |
| 16:00:00 | Over-selling detected in inventory system |
| 19:00:00 | Manual reconciliation completed |

### 10.4 Resolution Steps

```go
// CRDT-based conflict resolution
type OrderCRDT struct {
    OrderID       string
    Status        LWWRegister[string]
    ProcessingCenter LWWRegister[string]
    Inventory     PNCounter
    Version       VectorClock
}

func (o *OrderCRDT) Merge(other *OrderCRDT) *OrderCRDT {
    return &OrderCRDT{
        OrderID:          o.OrderID,
        Status:           o.Status.Merge(other.Status),
        ProcessingCenter: o.ProcessingCenter.Merge(other.ProcessingCenter),
        Inventory:        o.Inventory.Merge(other.Inventory),
        Version:          o.Version.Merge(other.Version),
    }
}

// Conditional write with version check
func updateOrder(client *dynamodb.Client, order *Order) error {
    _, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
        TableName: aws.String("orders"),
        Key: map[string]types.AttributeValue{
            "order_id": &types.AttributeValueS{Value: order.OrderID},
        },
        ConditionExpression: aws.String("version = :expected_version"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":expected_version": &types.AttributeValueN{Value: fmt.Sprintf("%d", order.Version)},
            ":new_status":      &types.AttributeValueS{Value: order.Status},
        },
        UpdateExpression: aws.String("SET status = :new_status, version = version + :inc"),
    })

    return err
}
```

### 10.5 Lessons Learned

1. **Last-write-wins** is rarely the right conflict resolution
2. **Vector clocks** or **CRDTs** should be used for concurrent writes
3. **Regional affinity** can reduce conflicts
4. **Conflict detection** should be explicit, not silent

### 10.6 Prevention Recommendations

```go
// Optimistic locking with version vectors
type VersionedOrder struct {
    OrderID string
    Status  string
    Version VectorClock
}

type VectorClock map[string]uint64

func (vc VectorClock) Increment(replica string) {
    vc[replica]++
}

func (vc VectorClock) Merge(other VectorClock) VectorClock {
    merged := make(VectorClock)
    for k, v := range vc {
        merged[k] = v
    }
    for k, v := range other {
        if v > merged[k] {
            merged[k] = v
        }
    }
    return merged
}

func (vc VectorClock) HappensBefore(other VectorClock) bool {
    allLessOrEqual := true
    atLeastOneLess := false

    for k, v := range vc {
        if v > other[k] {
            allLessOrEqual = false
            break
        }
        if v < other[k] {
            atLeastOneLess = true
        }
    }

    return allLessOrEqual && atLeastOneLess
}
```

---

## Summary and Best Practices

### Common Failure Patterns

| Pattern | Frequency | Impact | Detectability |
|---------|-----------|--------|---------------|
| Split-Brain | High | Critical | Medium |
| Consensus Liveness | Medium | High | Low |
| Hinted Handoff Overload | Medium | High | Medium |
| Session Expiration | High | Medium | High |
| Rebalance Storm | Medium | High | Medium |
| MVCC Bloat | Low | Critical | Low |
| Stale Reads | High | Medium | Low |
| KV Inconsistency | Low | Critical | Low |
| Queue Partition | Medium | High | Medium |
| Multi-Master Conflict | Low | Critical | Low |

### Prevention Checklist

- [ ] Implement proper fencing mechanisms
- [ ] Use pre-vote in Raft implementations
- [ ] Configure appropriate timeouts with jitter
- [ ] Implement circuit breakers for hint replay
- [ ] Use exponential backoff for reconnections
- [ ] Enable cooperative rebalancing
- [ ] Monitor replication lag and db size
- [ ] Use max staleness for secondary reads
- [ ] Implement idempotent consumers
- [ ] Use CRDTs or vector clocks for multi-master

### References

1. "Brewer's Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web Services" - Gilbert & Lynch
2. "The Chubby Lock Service for Loosely-Coupled Distributed Systems" - Burrows
3. "In Search of an Understandable Consensus Algorithm" - Ongaro & Ousterhout
4. "Dynamo: Amazon's Highly Available Key-value Store" - DeCandia et al.
5. "Kafka: A Distributed Messaging System for Log Processing" - Kreps et al.

---

*Document Size: 15+ KB | Level: S | Last Updated: 2026-04-03*
