# TS-003: Kafka Architecture - Internals & Go Implementation

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #kafka #streaming #distributed #internals #go
> **权威来源**:
>
> - [Apache Kafka Documentation](https://kafka.apache.org/documentation/) - Apache Software Foundation
> - [Kafka: The Definitive Guide](https://www.oreilly.com/library/view/kafka-the-definitive/) - O'Reilly Media
> - [KIP-500](https://cwiki.apache.org/confluence/display/KAFKA/KIP-500) - Kafka Raft Metadata Mode
> - [Confluent Kafka Internals](https://www.confluent.io/blog/) - Confluent Engineering

---

## 1. Kafka Internal Architecture

### 1.1 High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Apache Kafka System Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         Kafka Cluster (KRaft Mode)                     │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                   │  │
│  │  │  Controller │  │  Controller │  │  Controller │                   │  │
│  │  │  (Leader)   │  │  (Follower) │  │  (Follower) │  Metadata Quorum  │  │
│  │  │  Node 1     │  │  Node 2     │  │  Node 3     │                   │  │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                   │  │
│  │         │                │                │                          │  │
│  │         └────────────────┼────────────────┘                          │  │
│  │                          │ Raft Consensus (KRaft)                    │  │
│  └──────────────────────────┼───────────────────────────────────────────┘  │
│                             │                                              │
│                             ▼                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │                      Kafka Brokers                                    │ │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │ │
│  │  │   Broker 1  │◄──►│   Broker 2  │◄──►│   Broker 3  │              │ │
│  │  │  ┌───────┐  │    │  ┌───────┐  │    │  ┌───────┐  │              │ │
│  │  │  │TopicA │  │    │  │TopicA │  │    │  │TopicA │  │ Replication  │ │
│  │  │  │ -P0   │  │    │  │ -P1   │  │    │  │ -P2   │  │              │ │
│  │  │  │ -P1(R)│  │    │  │ -P2(R)│  │    │  │ -P0(R)│  │              │ │
│  │  │  └───────┘  │    │  └───────┘  │    │  └───────┘  │              │ │
│  │  │  ┌───────┐  │    │  ┌───────┐  │    │  ┌───────┐  │              │ │
│  │  │  │TopicB │  │    │  │TopicB │  │    │  │TopicB │  │              │ │
│  │  │  │ -P0   │  │    │  │ -P0(R)│  │    │  │ -P1   │  │              │ │
│  │  │  │ -P1(R)│  │    │  │ -P1   │  │    │  │ -P0(R)│  │              │ │
│  │  │  └───────┘  │    │  └───────┘  │    │  └───────┘  │              │ │
│  │  └─────────────┘    └─────────────┘    └─────────────┘              │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │                      Clients                                          │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │ │
│  │  │  Producers   │  │  Consumers   │  │  Connect/    │                │ │
│  │  │              │  │  (Groups)    │  │  Streams API │                │ │
│  │  └──────────────┘  └──────────────┘  └──────────────┘                │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 KRaft (Kafka Raft) Metadata Architecture

**KRaft Mode (Kafka 3.0+)**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    KRaft Metadata Management                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  传统模式 (ZooKeeper)                    KRaft 模式 (Kafka 3.0+)           │
│  ┌─────────┐    ┌─────────┐             ┌─────────────────────────┐        │
│  │  Kafka  │◄──►│   ZK    │             │  Kafka Metadata Quorum  │        │
│  │ Cluster │    │ Ensemble│             │  (Raft Consensus)       │        │
│  └─────────┘    └─────────┘             └─────────────────────────┘        │
│       │                                       │                            │
│       │ Metadata RPCs                         │ Log Replication              │
│       ▼                                       ▼                            │
│  ┌─────────┐    ┌─────────┐             ┌─────────┐  ┌─────────┐          │
│  │ /brokers│    │/topics  │             │ Node 1  │  │ Node 2  │          │
│  │ /ids    │    │/partitions│           │ (Leader)│  │(Follower)│         │
│  │ /topics │    │/config   │            └────┬────┘  └────┬────┘          │
│  └─────────┘    └─────────┘                  │            │                │
│                                              └────────────┘                │
│                                                   │                          │
│                                              ┌────┴────┐                     │
│                                              │ Node 3  │                     │
│                                              │(Follower)│                    │
│                                              └─────────┘                     │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                       Raft Log Structure                              │  │
│  ├──────────────────────────────────────────────────────────────────────┤  │
│  │                                                                       │  │
│  │  Offset │ Term │ Metadata Record                                       │  │
│  │  ───────┼──────┼────────────────────────────────────────────────      │  │
│  │     0   │  1   │ Cluster UUID                                          │  │
│  │     1   │  1   │ Broker Registration (node 1)                          │  │
│  │     2   │  1   │ Broker Registration (node 2)                          │  │
│  │     3   │  1   │ Topic Creation: "orders"                              │  │
│  │     4   │  1   │ Partition Assignment: orders-0 → [1, 2, 3]            │  │
│  │     5   │  2   │ Leader Election (term incremented)                    │  │
│  │    ...  │ ...  │ ...                                                   │  │
│  │                                                                       │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Quorum Controller Benefits**:

- **简化运维**: 无需维护 ZooKeeper 集群
- **更高性能**: 元数据操作吞吐量提升 10x+
- **更快恢复**: Controller failover < 3s (vs 30s+ with ZK)
- **更大规模**: 支持百万级分区

### 1.3 Partition & Replication Internals

**Partition Storage Structure**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Partition Storage Layout                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Topic: user-events, Partition: 0, Replica on Broker 1                      │
│                                                                              │
│  /var/lib/kafka-logs/user-events-0/                                          │
│  ├── 00000000000000000000.log      ← Segment 0 (base offset = 0)            │
│  ├── 00000000000000000000.index    ← Offset to position index               │
│  ├── 00000000000000000000.timeindex ← Timestamp to offset index             │
│  ├── 00000000000000356892.log      ← Segment 1 (base offset = 356892)      │
│  ├── 00000000000000356892.index                                              │
│  ├── 00000000000000356892.timeindex                                          │
│  ├── 00000000000000789123.log                                                │
│  ├── ...                                                                      │
│  └── leader-epoch-checkpoint     ← Leader epoch tracking                    │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Log Segment Structure                            │   │
│  ├──────────────────────────────────────────────────────────────────────┤   │
│  │                                                                       │   │
│  │  RecordBatch 1                        RecordBatch 2                   │   │
│  │  ┌─────────────────┐                  ┌─────────────────┐            │   │
│  │  │ BaseOffset: 0   │                  │ BaseOffset: 100 │            │   │
│  │  │ RecordCount: 100│                  │ RecordCount: 100│            │   │
│  │  ├─────────────────┤                  ├─────────────────┤            │   │
│  │  │ Record 0        │                  │ Record 100      │            │   │
│  │  │ Record 1        │                  │ Record 101      │            │   │
│  │  │ ...             │                  │ ...             │            │   │
│  │  │ Record 99       │                  │ Record 199      │            │   │
│  │  └─────────────────┘                  └─────────────────┘            │   │
│  │                                                                       │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Index Structure                                  │   │
│  ├──────────────────────────────────────────────────────────────────────┤   │
│  │                                                                       │   │
│  │  Relative Offset (4B) │ Position (4B) │                              │   │
│  │  ─────────────────────┼───────────────┼─                             │   │
│  │           0           │       0       │  → RecordBatch 1 starts at 0 │   │
│  │         100           │    2048       │  → RecordBatch 2 starts at 2048│  │
│  │         200           │    4096       │                              │   │
│  │         ...           │     ...       │                              │   │
│  │                                                                       │   │
│  │  Sparse Index: 每 4KB 写入一个索引项 (log.index.interval.bytes)       │   │
│  │                                                                       │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**ISR (In-Sync Replicas) Management**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ISR Management & Replication Flow                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Replication Factor = 3                                                      │
│                                                                              │
│  ┌───────────────┐      ┌───────────────┐      ┌───────────────┐           │
│  │    Leader     │      │  Follower 1   │      │  Follower 2   │           │
│  │   (HW = 100)  │      │   (HW = 98)   │      │   (HW = 95)   │           │
│  │   (LEO = 105) │      │   (LEO = 100) │      │   (LEO = 98)  │           │
│  └───────┬───────┘      └───────┬───────┘      └───────┬───────┘           │
│          │                      │                      │                    │
│          │  1. Fetch Request    │                      │                    │
│          │◄─────────────────────┤                      │                    │
│          │  (offset = 98)       │                      │                    │
│          │                      │                      │                    │
│          │  2. Return Records   │                      │                    │
│          │─────────────────────►│                      │                    │
│          │  (98-105)            │                      │                    │
│          │                      │                      │                    │
│          │                      │  3. Append to Log    │                    │
│          │                      │  LEO → 106           │                    │
│          │                      │                      │                    │
│          │  4. Fetch Request    │                      │                    │
│          │◄────────────────────────────────────────────┤                    │
│          │  (offset = 95)                              │                    │
│          │                      │                      │                    │
│          │  5. Return Records   │                      │                    │
│          │────────────────────────────────────────────►│                    │
│          │  (95-105)                                   │                    │
│          │                      │                      │                    │
│          │  6. Update ISR (if replica caught up)       │                    │
│          │  ISR = {Leader, Follower 1}                 │                    │
│          │  Follower 2 removed from ISR (replica.lag.time.max.ms exceeded)  │
│          │                      │                      │                    │
│  ┌───────┴───────┐      ┌───────┴───────┐      ┌───────┴───────┐           │
│  │  HW = 100     │      │  HW = 100     │      │  HW = 98      │           │
│  │  ISR updated  │      │  Joined ISR   │      │  Out of ISR   │           │
│  └───────────────┘      └───────────────┘      └───────────────┘           │
│                                                                              │
│  Key Metrics:                                                                │
│  - HW (High Watermark): 已提交的最大 offset，消费者只能读到 HW               │
│  - LEO (Log End Offset): 下一条待写入消息的 offset                          │
│  - replica.lag.time.max.ms: 副本最大允许滞后时间 (默认 30s)                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.4 Producer Internals

**Producer Message Flow**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Producer Message Flow                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         Producer Architecture                          │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Application → Serializer → Partitioner → RecordAccumulator → Sender   │  │
│  │                                                                        │  │
│  │  ┌───────────┐   ┌──────────┐   ┌──────────┐   ┌───────────────────┐  │  │
│  │  │   User    │──►│ Serializer│──►│ Partition│──►│ RecordAccumulator │  │  │
│  │  │   Code    │   │ (Avro/    │   │  er      │   │  (Buffer Pool)    │  │  │
│  │  │           │   │  Protobuf)│   │          │   │                   │  │  │
│  │  └───────────┘   └──────────┘   └──────────┘   └─────────┬─────────┘  │  │
│  │                                                           │            │  │
│  │                              ┌────────────────────────────┘            │  │
│  │                              │                                         │  │
│  │                              ▼                                         │  │
│  │                    ┌───────────────────┐                               │  │
│  │                    │   Record Batches  │                               │  │
│  │                    │   per Partition   │                               │  │
│  │                    ├───────────────────┤                               │  │
│  │                    │ TopicA-Partition0 │ ──┐                           │  │
│  │                    │ TopicA-Partition1 │ ──┤                           │  │
│  │                    │ TopicB-Partition0 │ ──┼──► Sender Thread          │  │
│  │                    │ TopicB-Partition1 │ ──┤                           │  │
│  │                    └───────────────────┘   │                           │  │
│  │                                            ▼                           │  │
│  │                              ┌───────────────────────┐                 │  │
│  │                              │      Network Client   │                 │  │
│  │                              │  (Selector, NIO-based)│                 │  │
│  │                              └───────────┬───────────┘                 │  │
│  │                                          │                             │  │
│  │                              ┌───────────┼───────────┐                 │  │
│  │                              ▼           ▼           ▼                 │  │
│  │                         ┌────────┐  ┌────────┐  ┌────────┐            │  │
│  │                         │Broker 1│  │Broker 2│  │Broker 3│            │  │
│  │                         └────────┘  └────────┘  └────────┘            │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                      Partitioner Strategies                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. DefaultPartitioner (Sticky with Batching)                         │  │
│  │                                                                        │  │
│  │     ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐            │  │
│  │     │Batch P0 │   │Batch P1 │   │Batch P0 │   │Batch P2 │            │  │
│  │     │  1-10   │   │  1-8    │   │  11-20  │   │  1-12   │            │  │
│  │     └─────────┘   └─────────┘   └─────────┘   └─────────┘            │  │
│  │          │             │             │             │                  │  │
│  │          └─────────────┴─────────────┴─────────────┘                  │  │
│  │                          │                                            │  │
│  │                    Sticky for batching efficiency                     │  │
│  │                                                                        │  │
│  │  2. RoundRobinPartitioner                                             │  │
│  │     P0 → P1 → P2 → P0 → P1 → P2 ...                                  │  │
│  │                                                                        │  │
│  │  3. UniformStickyPartitioner                                          │  │
│  │     Stick to partition until batch is full, then round-robin          │  │
│  │                                                                        │  │
│  │  4. Custom Partitioner (based on key)                                 │  │
│  │     partition = murmur2(key) % numPartitions                          │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Producer Acknowledgment Modes**:

| acks Setting | Behavior | Durability | Throughput |
|--------------|----------|------------|------------|
| `acks=0` | Fire-and-forget, no acknowledgment | Lowest | Highest |
| `acks=1` | Leader acknowledgment only | Medium | High |
| `acks=all` | All ISRs acknowledgment | Highest | Lower |

### 1.5 Consumer Group Internals

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Consumer Group Rebalance Protocol                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Consumer Group: analytics-group                                             │
│  Topic: user-events (6 partitions)                                           │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Consumer Group State Machine                        │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────┐    Join    ┌──────────┐   Sync   ┌──────────┐           │  │
│  │  │ Stable  │ ─────────►│ Preparing│─────────►│ Awaiting │           │  │
│  │  │         │           │ Rebalance│          │ Sync     │           │  │
│  │  └────┬────┘           └──────────┘          └────┬─────┘           │  │
│  │       ▲                                           │                  │  │
│  │       │              Rebalance Complete           │                  │  │
│  │       └───────────────────────────────────────────┘                  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Rebalance Process (Eager Protocol)                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Initial State:                                                        │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐                            │  │
│  │  │Consumer 1│  │Consumer 2│  │Consumer 3│                            │  │
│  │  │ P0, P1   │  │ P2, P3   │  │ P4, P5   │                            │  │
│  │  └──────────┘  └──────────┘  └──────────┘                            │  │
│  │                                                                        │  │
│  │  Step 1: New consumer (C4) joins                                       │  │
│  │                                                                        │  │
│  │  C4 ──JoinGroup──► Group Coordinator (Broker)                          │  │
│  │                                                                        │  │
│  │  Step 2: Coordinator triggers rebalance                                │  │
│  │                                                                        │  │
│  │  Coordinator ──Rebalance──► C1, C2, C3 (revoke partitions)             │  │
│  │                                                                        │  │
│  │  Step 3: All consumers re-join                                         │  │
│  │                                                                        │  │
│  │  C1, C2, C3, C4 ──JoinGroup──► Coordinator                             │  │
│  │  (Coordinator selects leader, usually first to join)                   │  │
│  │                                                                        │  │
│  │  Step 4: Leader performs assignment                                    │  │
│  │                                                                        │  │
│  │  Leader (C1) ──Assign:                                                 │  │
│  │    C1: P0, P1                                                          │  │
│  │    C2: P2, P3                                                          │  │
│  │    C3: P4                                                              │  │
│  │    C4: P5                                                              │  │
│  │                                                                        │  │
│  │  Step 5: SyncGroup to all members                                      │  │
│  │                                                                        │  │
│  │  Final State:                                                          │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐               │  │
│  │  │Consumer 1│  │Consumer 2│  │Consumer 3│  │Consumer 4│               │  │
│  │  │ P0, P1   │  │ P2, P3   │  │   P4     │  │   P5     │               │  │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘               │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Cooperative Rebalance (Incremental)                 │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Unlike eager rebalance, cooperative protocol:                         │  │
│  │  1. Doesn't revoke all partitions at once                              │  │
│  │  2. Only revokes partitions being reassigned                           │  │
│  │  3. Consumers can continue processing unaffected partitions            │  │
│  │  4. Reduces stop-the-world during rebalance                            │  │
│  │                                                                        │  │
│  │  partition.assignment.strategy = org.apache.kafka.clients.consumer.    │  │
│  │                                   CooperativeStickyAssignor            │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Configuration Best Practices

### 2.1 Broker Configuration

```properties
# ==================== Server Basics ====================
# KRaft 模式节点 ID
node.id=1
process.roles=broker,controller
controller.quorum.voters=1@localhost:9093,2@localhost:9093,3@localhost:9093

# 监听器配置
listeners=PLAINTEXT://:9092,CONTROLLER://:9093
inter.broker.listener.name=PLAINTEXT
advertised.listeners=PLAINTEXT://localhost:9092
listener.security.protocol.map=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT

# ==================== Log Configuration ====================
# 日志目录 (RAID 10 推荐)
log.dirs=/var/lib/kafka-logs

# 日志段大小 (默认 1GB)
log.segment.bytes=1073741824

# 日志保留时间 (默认 7 天)
log.retention.hours=168

# 日志保留大小 (可选)
# log.retention.bytes=107374182400

# 日志清理策略 (delete 或 compact)
log.cleanup.policy=delete

# 段滚动时间 (默认 7 天)
log.roll.hours=168

# ==================== Replication ====================
# 默认副本因子
default.replication.factor=3

# 最小 ISR 大小
default.min.insync.replicas=2

# 副本拉取线程数
num.replica.fetchers=4

# 副本最大滞后时间
replica.lag.time.max.ms=30000

# 非 ISR 副本允许作为 leader (数据丢失风险!)
unclean.leader.election.enable=false

# ==================== Network & I/O ====================
# 网络线程数 (num.network.threads)
num.network.threads=8

# I/O 线程数 (num.io.threads)
num.io.threads=16

# Socket 发送/接收缓冲区
socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400
socket.request.max.bytes=104857600

# ==================== Performance ====================
# 消息压缩类型 (none, gzip, snappy, lz4, zstd)
compression.type=producer

# 批量大小
log.flush.interval.messages=10000
log.flush.interval.ms=1000

# 页缓存刷新
log.flush.scheduler.interval.ms=3000
```

### 2.2 Producer Configuration

```properties
# ==================== Bootstrap & Connection ====================
bootstrap.servers=kafka1:9092,kafka2:9092,kafka3:9092
client.id=order-service-producer

# ==================== Reliability ====================
# 确认级别 (0, 1, all)
acks=all

# 重试次数 (无限重试直到 delivery.timeout.ms)
retries=2147483647

# 幂等生产者 (exactly-once 必需)
enable.idempotence=true

# 最大 inflight 请求 (幂等生产者最大为 5)
max.in.flight.requests.per.connection=5

# 事务 ID (exactly-once 语义)
transactional.id=order-service-tx-001

# ==================== Batching ====================
# 批量大小
batch.size=16384

# 批处理延迟 (等待更多记录加入批次)
linger.ms=5

# 缓冲区大小
buffer.memory=33554432

# 压缩类型
compression.type=lz4

# ==================== Timeout ====================
# 请求超时
request.timeout.ms=30000

# 元数据过期
metadata.max.age.ms=300000

# 连接最大空闲时间
connections.max.idle.ms=540000
```

### 2.3 Consumer Configuration

```properties
# ==================== Bootstrap & Group ====================
bootstrap.servers=kafka1:9092,kafka2:9092,kafka3:9092
group.id=order-processing-group
client.id=order-consumer-1

# ==================== Offset Management ====================
# 自动提交 (生产环境建议关闭)
enable.auto.commit=false

# 自动提交间隔
# auto.commit.interval.ms=5000

# 偏移重置策略 (earliest, latest, none)
auto.offset.reset=earliest

# ==================== Fetch Configuration ====================
# 最小拉取字节 (减少空轮询)
fetch.min.bytes=1

# 最大拉取字节
fetch.max.bytes=52428800

# 分区最大字节
max.partition.fetch.bytes=1048576

# 拉取超时
fetch.max.wait.ms=500

# 最大轮询记录 (单次 poll 返回的最大记录数)
max.poll.records=500

# 最大轮询间隔 (超过则认为 consumer 死亡)
max.poll.interval.ms=300000

# ==================== Heartbeat ====================
# 心跳间隔
heartbeat.interval.ms=3000

# 会话超时
session.timeout.ms=10000

# ==================== Partition Assignment ====================
# 分区分配策略
partition.assignment.strategy=org.apache.kafka.clients.consumer.CooperativeStickyAssignor

# ==================== Exactly-Once ====================
# 隔离级别 (read_committed, read_uncommitted)
isolation.level=read_committed
```

---

## 3. Go Implementation with franz-go

### 3.1 Producer Implementation

```go
package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/twmb/franz-go/pkg/kgo"
    "github.com/twmb/franz-go/pkg/sasl/scram"
)

// Producer Kafka 生产者封装
type Producer struct {
    client *kgo.Client
    topic  string
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
    Brokers      []string
    Topic        string
    ClientID     string
    EnableIdempotency bool
    EnableCompression bool
    SASLUser     string
    SASLPassword string
}

// NewProducer 创建生产者
func NewProducer(cfg *ProducerConfig) (*Producer, error) {
    opts := []kgo.Opt{
        kgo.SeedBrokers(cfg.Brokers...),
        kgo.ClientID(cfg.ClientID),
        kgo.DefaultProduceTopic(cfg.Topic),
        kgo.RequiredAcks(kgo.AllISRAcks()),
        kgo.BatchMaxBytes(1000000),
        kgo.BatchTimeout(5 * time.Millisecond),
    }

    if cfg.EnableIdempotency {
        opts = append(opts, kgo.Idempotent())
    }

    if cfg.EnableCompression {
        opts = append(opts, kgo.ProducerBatchCompression(kgo.Lz4Compression()))
    }

    if cfg.SASLUser != "" {
        auth := scram.Auth{
            User: cfg.SASLUser,
            Pass: cfg.SASLPassword,
        }.AsSha256Mechanism()
        opts = append(opts, kgo.SASL(auth))
    }

    client, err := kgo.NewClient(opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to create kafka producer: %w", err)
    }

    return &Producer{
        client: client,
        topic:  cfg.Topic,
    }, nil
}

// Produce 同步发送消息
func (p *Producer) Produce(ctx context.Context, key, value []byte) error {
    record := &kgo.Record{
        Topic: p.topic,
        Key:   key,
        Value: value,
    }

    var wg sync.WaitGroup
    wg.Add(1)

    var produceErr error
    p.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
        defer wg.Done()
        if err != nil {
            produceErr = err
        }
    })

    wg.Wait()
    return produceErr
}

// ProduceAsync 异步发送消息
func (p *Producer) ProduceAsync(ctx context.Context, key, value []byte, callback func(error)) {
    record := &kgo.Record{
        Topic: p.topic,
        Key:   key,
        Value: value,
    }

    p.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
        if callback != nil {
            callback(err)
        }
    })
}

// ProduceJSON 发送 JSON 消息
func (p *Producer) ProduceJSON(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return p.Produce(ctx, []byte(key), data)
}

// ProduceBatch 批量发送
func (p *Producer) ProduceBatch(ctx context.Context, records []*kgo.Record) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(records))

    for _, record := range records {
        wg.Add(1)
        rec := record // 避免闭包问题

        p.client.Produce(ctx, rec, func(_ *kgo.Record, err error) {
            defer wg.Done()
            if err != nil {
                select {
                case errChan <- err:
                default:
                }
            }
        })
    }

    wg.Wait()
    close(errChan)

    // 返回第一个错误
    for err := range errChan {
        return err
    }
    return nil
}

// Flush 刷新所有待发送消息
func (p *Producer) Flush(ctx context.Context) error {
    p.client.Flush(ctx)
    return nil
}

// Close 关闭生产者
func (p *Producer) Close() {
    p.client.Close()
}

// TransactionalProducer 事务生产者
type TransactionalProducer struct {
    client *kgo.Client
}

// NewTransactionalProducer 创建事务生产者
func NewTransactionalProducer(brokers []string, transactionalID string) (*TransactionalProducer, error) {
    client, err := kgo.NewClient(
        kgo.SeedBrokers(brokers...),
        kgo.TransactionalID(transactionalID),
        kgo.Idempotent(),
        kgo.RequiredAcks(kgo.AllISRAcks()),
    )
    if err != nil {
        return nil, err
    }

    return &TransactionalProducer{client: client}, nil
}

// BeginTransaction 开始事务
func (tp *TransactionalProducer) BeginTransaction() error {
    return tp.client.BeginTransaction()
}

// CommitTransaction 提交事务
func (tp *TransactionalProducer) CommitTransaction(ctx context.Context) error {
    return tp.client.CommitTransaction(ctx)
}

// AbortTransaction 中止事务
func (tp *TransactionalProducer) AbortTransaction(ctx context.Context) error {
    return tp.client.AbortTransaction(ctx)
}

// SendOffsetsToTransaction 发送偏移量到事务
func (tp *TransactionalProducer) SendOffsetsToTransaction(
    ctx context.Context,
    offsets kgo.GroupTransactSessions,
    group string,
) error {
    return tp.client.SendOffsetsToTransaction(ctx, offsets, group)
}
```

### 3.2 Consumer Implementation

```go
package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/twmb/franz-go/pkg/kgo"
)

// MessageHandler 消息处理函数类型
type MessageHandler func(ctx context.Context, record *kgo.Record) error

// Consumer Kafka 消费者封装
type Consumer struct {
    client  *kgo.Client
    handler MessageHandler
    topics  []string
    wg      sync.WaitGroup
    stopCh  chan struct{}
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
    Brokers          []string
    Topics           []string
    GroupID          string
    ClientID         string
    AutoCommit       bool
    MaxPollRecords   int
    SessionTimeout   time.Duration
    HeartbeatInterval time.Duration
}

// NewConsumer 创建消费者
func NewConsumer(cfg *ConsumerConfig, handler MessageHandler) (*Consumer, error) {
    opts := []kgo.Opt{
        kgo.SeedBrokers(cfg.Brokers...),
        kgo.ConsumerGroup(cfg.GroupID),
        kgo.ConsumeTopics(cfg.Topics...),
        kgo.ClientID(cfg.ClientID),
        kgo.FetchMinBytes(1),
        kgo.FetchMaxBytes(50 << 20),
        kgo.FetchMaxWait(500 * time.Millisecond),
        kgo.SessionTimeout(cfg.SessionTimeout),
        kgo.HeartbeatInterval(cfg.HeartbeatInterval),
        kgo.RebalanceGroupProtocol(kgo.CooperativeStickyBalancer()),
    }

    if !cfg.AutoCommit {
        // 手动提交模式
        opts = append(opts, kgo.DisableAutoCommit())
    }

    if cfg.MaxPollRecords > 0 {
        opts = append(opts, kgo.MaxPollRecords(int32(cfg.MaxPollRecords)))
    }

    client, err := kgo.NewClient(opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
    }

    return &Consumer{
        client:  client,
        handler: handler,
        topics:  cfg.Topics,
        stopCh:  make(chan struct{}),
    }, nil
}

// Start 开始消费
func (c *Consumer) Start(ctx context.Context) error {
    defer c.wg.Wait()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-c.stopCh:
            return nil
        default:
        }

        // 拉取记录
        fetchRecords := c.client.PollRecords(ctx, 500)
        if fetchRecords.IsClientClosed() {
            return fmt.Errorf("kafka client closed")
        }

        if err := fetchRecords.Err0(); err != nil {
            // 处理错误，如断线等
            time.Sleep(time.Second)
            continue
        }

        // 处理记录
        iter := fetchRecords.RecordIter()
        for !iter.Done() {
            record := iter.Next()

            if err := c.handler(ctx, record); err != nil {
                // 处理错误：可以记录、发送到死信队列等
                fmt.Printf("handler error: %v\n", err)
                continue
            }
        }

        // 手动提交偏移量
        if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
            fmt.Printf("commit error: %v\n", err)
        }
    }
}

// Stop 停止消费
func (c *Consumer) Stop() {
    close(c.stopCh)
    c.client.Close()
}

// ConsumerGroupRebalanceCallback 再平衡回调
type ConsumerGroupRebalanceCallback interface {
    OnAssigned(ctx context.Context, client *kgo.Client, assigned map[string][]int32)
    OnRevoked(ctx context.Context, client *kgo.Client, revoked map[string][]int32)
    OnLost(ctx context.Context, client *kgo.Client, lost map[string][]int32)
}

// MultiWorkerConsumer 多工作线程消费者
type MultiWorkerConsumer struct {
    client      *kgo.Client
    handler     MessageHandler
    workerCount int
    wg          sync.WaitGroup
}

// NewMultiWorkerConsumer 创建多工作线程消费者
func NewMultiWorkerConsumer(cfg *ConsumerConfig, handler MessageHandler, workerCount int) (*MultiWorkerConsumer, error) {
    opts := []kgo.Opt{
        kgo.SeedBrokers(cfg.Brokers...),
        kgo.ConsumerGroup(cfg.GroupID),
        kgo.ConsumeTopics(cfg.Topics...),
        kgo.ClientID(cfg.ClientID),
        kgo.DisableAutoCommit(),
        kgo.RebalanceGroupProtocol(kgo.CooperativeStickyBalancer()),
    }

    client, err := kgo.NewClient(opts...)
    if err != nil {
        return nil, err
    }

    return &MultiWorkerConsumer{
        client:      client,
        handler:     handler,
        workerCount: workerCount,
    }, nil
}

// Start 启动多工作线程消费
func (c *MultiWorkerConsumer) Start(ctx context.Context) {
    recordCh := make(chan *kgo.Record, c.workerCount*2)

    // 启动工作线程
    for i := 0; i < c.workerCount; i++ {
        c.wg.Add(1)
        go c.worker(ctx, recordCh)
    }

    // 拉取循环
    go func() {
        defer close(recordCh)

        for {
            select {
            case <-ctx.Done():
                return
            default:
            }

            fetchRecords := c.client.PollRecords(ctx, 100)
            if fetchRecords.IsClientClosed() {
                return
            }

            iter := fetchRecords.RecordIter()
            for !iter.Done() {
                record := iter.Next()
                select {
                case recordCh <- record:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()

    c.wg.Wait()
}

func (c *MultiWorkerConsumer) worker(ctx context.Context, recordCh <-chan *kgo.Record) {
    defer c.wg.Done()

    for {
        select {
        case <-ctx.Done():
            return
        case record, ok := <-recordCh:
            if !ok {
                return
            }

            if err := c.handler(ctx, record); err != nil {
                fmt.Printf("worker handler error: %v\n", err)
            }

            // 标记为可提交
            c.client.MarkCommitRecords(record)
        }
    }
}

// CommitOffsets 提交偏移量
func (c *MultiWorkerConsumer) CommitOffsets(ctx context.Context) error {
    return c.client.CommitMarkedOffsets(ctx)
}

// Stop 停止
func (c *MultiWorkerConsumer) Stop() {
    c.client.Close()
}

// JSONHandler 创建 JSON 消息处理器
func JSONHandler(v interface{}, fn func(context.Context, string, interface{}) error) MessageHandler {
    return func(ctx context.Context, record *kgo.Record) error {
        val := v // 创建新实例
        if err := json.Unmarshal(record.Value, &val); err != nil {
            return fmt.Errorf("unmarshal error: %w", err)
        }
        return fn(ctx, string(record.Key), val)
    }
}
```

### 3.3 Admin Client Operations

```go
package kafka

import (
    "context"
    "fmt"
    "time"

    "github.com/twmb/franz-go/pkg/kadm"
    "github.com/twmb/franz-go/pkg/kgo"
)

// AdminClient Kafka 管理客户端
type AdminClient struct {
    client *kadm.Client
}

// NewAdminClient 创建管理客户端
func NewAdminClient(brokers []string) (*AdminClient, error) {
    client, err := kgo.NewClient(
        kgo.SeedBrokers(brokers...),
    )
    if err != nil {
        return nil, err
    }

    return &AdminClient{
        client: kadm.NewClient(client),
    }, nil
}

// CreateTopic 创建主题
func (a *AdminClient) CreateTopic(
    ctx context.Context,
    topic string,
    partitions int32,
    replicationFactor int16,
    configs map[string]string,
) error {
    req := kadm.CreateTopicConfig{
        Topic:             topic,
        NumPartitions:     partitions,
        ReplicationFactor: replicationFactor,
    }

    for k, v := range configs {
        req.ConfigEntries = append(req.ConfigEntries, kadm.CreateTopicConfigEntry{
            Name:  k,
            Value: v,
        })
    }

    resp, err := a.client.CreateTopic(ctx, partitions, replicationFactor, configs, topic)
    if err != nil {
        return err
    }

    if resp.Err != nil {
        return fmt.Errorf("failed to create topic %s: %w", topic, resp.Err)
    }

    return nil
}

// DeleteTopic 删除主题
func (a *AdminClient) DeleteTopic(ctx context.Context, topics ...string) error {
    resps, err := a.client.DeleteTopics(ctx, topics...)
    if err != nil {
        return err
    }

    for _, resp := range resps {
        if resp.Err != nil {
            return fmt.Errorf("failed to delete topic %s: %w", resp.Topic, resp.Err)
        }
    }

    return nil
}

// ListTopics 列出主题
func (a *AdminClient) ListTopics(ctx context.Context) ([]string, error) {
    topics, err := a.client.ListTopics(ctx)
    if err != nil {
        return nil, err
    }

    var names []string
    topics.Each(func(topic kadm.ListedTopic) {
        if !topic.IsInternal {
            names = append(names, topic.Topic)
        }
    })

    return names, nil
}

// DescribeTopic 描述主题详情
func (a *AdminClient) DescribeTopic(ctx context.Context, topic string) (*TopicDescription, error) {
    details, err := a.client.DescribeTopicConfigs(ctx, topic)
    if err != nil {
        return nil, err
    }

    // 获取分区信息
    offsets, err := a.client.ListEndOffsets(ctx, topic)
    if err != nil {
        return nil, err
    }

    desc := &TopicDescription{
        Name:       topic,
        Partitions: make(map[int32]PartitionInfo),
    }

    for _, detail := range details {
        for _, entry := range detail.Configs {
            desc.Configs[entry.Name] = entry.Value
        }
    }

    offsets.Each(func(o kadm.ListedOffset) {
        desc.Partitions[o.Partition] = PartitionInfo{
            Leader:     o.Leader,
            HighWatermark: o.Offset,
        }
    })

    return desc, nil
}

// TopicDescription 主题描述
type TopicDescription struct {
    Name       string
    Configs    map[string]string
    Partitions map[int32]PartitionInfo
}

// PartitionInfo 分区信息
type PartitionInfo struct {
    Leader        int32
    HighWatermark int64
}

// AlterPartitionCount 修改分区数
func (a *AdminClient) AlterPartitionCount(ctx context.Context, topic string, newCount int32) error {
    resp, err := a.client.CreatePartitions(ctx, newCount, topic)
    if err != nil {
        return err
    }

    if resp.Err != nil {
        return fmt.Errorf("failed to alter partitions: %w", resp.Err)
    }

    return nil
}

// ConsumerGroupOffsets 获取消费者组偏移量
func (a *AdminClient) ConsumerGroupOffsets(ctx context.Context, group string) (map[string]map[int32]int64, error) {
    offsets, err := a.client.FetchOffsets(ctx, group)
    if err != nil {
        return nil, err
    }

    result := make(map[string]map[int32]int64)
    offsets.Each(func(o kadm.OffsetResponse) {
        if result[o.Topic] == nil {
            result[o.Topic] = make(map[int32]int64)
        }
        result[o.Topic][o.Partition] = o.At
    })

    return result, nil
}

// ResetConsumerGroupOffsets 重置消费者组偏移量
func (a *AdminClient) ResetConsumerGroupOffsets(
    ctx context.Context,
    group string,
    topic string,
    partitionOffsets map[int32]int64,
) error {
    offsets := make map[string]map[int32]kadm.Offset
    offsets[topic] = make(map[int32]kadm.Offset)

    for partition, offset := range partitionOffsets {
        offsets[topic][partition] = kadm.NewOffset().At(offset)
    }

    return a.client.CommitAllOffsets(ctx, group, offsets)
}

// Close 关闭客户端
func (a *AdminClient) Close() {
    a.client.Close()
}
```

---

## 4. Performance Tuning

### 4.1 Producer Tuning

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Producer Performance Tuning                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     Throughput Optimization                           │   │
│  ├──────────────────────────────────────────────────────────────────────┤   │
│  │                                                                       │   │
│  │  1. Batching (批处理)                                                  │   │
│  │     ┌────────────────────────────────────────────────────────────┐   │   │
│  │     │  linger.ms: 5-100ms (等待更多记录加入批次)                   │   │   │
│  │     │  batch.size: 32KB-256KB (批次大小)                           │   │   │
│  │     │  buffer.memory: 64MB+ (缓冲区大小)                           │   │   │
│  │     └────────────────────────────────────────────────────────────┘   │   │
│  │                                                                       │   │
│  │  2. Compression (压缩)                                                 │   │
│  │     ┌────────────────────────────────────────────────────────────┐   │   │
│  │     │  none:  无压缩，CPU 最低，网络最高                           │   │   │
│  │     │  gzip:  高压缩比，CPU 开销大                                 │   │   │
│  │     │  snappy: 平衡选择 (推荐)                                     │   │   │
│  │     │  lz4:   低延迟，高吞吐 (推荐)                                │   │   │
│  │     │  zstd:  最高压缩比，Kafka 2.1+ (推荐大数据量)                │   │   │
│  │     └────────────────────────────────────────────────────────────┘   │   │
│  │                                                                       │   │
│  │  3. Idempotent Producer (幂等生产者)                                   │   │
│  │     - enable.idempotence=true                                         │   │
│  │     - 自动处理重试时的重复消息                                        │   │
│  │     - max.in.flight.requests.per.connection 可以提高到 5              │   │
│  │                                                                       │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     Latency Optimization                              │   │
│  ├──────────────────────────────────────────────────────────────────────┤   │
│  │                                                                       │   │
│  │  低延迟配置 (实时场景):                                                │   │
│  │  ┌────────────────────────────────────────────────────────────┐       │   │
│  │  │  linger.ms = 0       (立即发送)                             │       │   │
│  │  │  batch.size = 1      (最小批次)                             │       │   │
│  │  │  compression.type = none (无压缩)                           │       │   │
│  │  │  acks = 1            (仅 leader 确认)                       │       │   │
│  │  └────────────────────────────────────────────────────────────┘       │   │
│  │                                                                       │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Consumer Tuning

```go
// ConsumerPerformanceConfig 高性能消费者配置
func ConsumerPerformanceConfig() []kgo.Opt {
    return []kgo.Opt{
        // 提高拉取效率
        kgo.FetchMinBytes(1024),           // 至少 1KB 数据才返回
        kgo.FetchMaxBytes(50 << 20),       // 最大 50MB
        kgo.FetchMaxWait(500 * time.Millisecond),

        // 增加单次 poll 记录数
        kgo.MaxPollRecords(1000),

        // 消费者组再平衡优化
        kgo.RebalanceGroupProtocol(kgo.CooperativeStickyBalancer()),

        // 会话配置
        kgo.SessionTimeout(30 * time.Second),
        kgo.HeartbeatInterval(3 * time.Second),

        // 读取隔离级别 (exactly-once)
        kgo.IsolationLevel(kgo.ReadCommitted()),
    }
}

// ParallelConsumer 并行消费模式
type ParallelConsumer struct {
    client      *kgo.Client
    partitions  map[int32]chan *kgo.Record
    workers     map[int32]context.CancelFunc
}

// NewParallelConsumer 创建分区级并行消费者
func NewParallelConsumer(brokers []string, topics []string, groupID string) (*ParallelConsumer, error) {
    pc := &ParallelConsumer{
        partitions: make(map[int32]chan *kgo.Record),
        workers:    make(map[int32]context.CancelFunc),
    }

    opts := []kgo.Opt{
        kgo.SeedBrokers(brokers...),
        kgo.ConsumerGroup(groupID),
        kgo.ConsumeTopics(topics...),
        kgo.OnPartitionsAssigned(pc.onAssigned),
        kgo.OnPartitionsRevoked(pc.onRevoked),
        kgo.OnPartitionsLost(pc.onLost),
    }

    client, err := kgo.NewClient(opts...)
    if err != nil {
        return nil, err
    }

    pc.client = client
    return pc, nil
}

func (pc *ParallelConsumer) onAssigned(ctx context.Context, client *kgo.Client, assigned map[string][]int32) {
    for topic, partitions := range assigned {
        for _, p := range partitions {
            ch := make(chan *kgo.Record, 1000)
            pc.partitions[p] = ch

            workerCtx, cancel := context.WithCancel(ctx)
            pc.workers[p] = cancel

            go pc.partitionWorker(workerCtx, topic, p, ch)
        }
    }
}

func (pc *ParallelConsumer) onRevoked(ctx context.Context, client *kgo.Client, revoked map[string][]int32) {
    for _, partitions := range revoked {
        for _, p := range partitions {
            if cancel, ok := pc.workers[p]; ok {
                cancel()
                delete(pc.workers, p)
            }
            if ch, ok := pc.partitions[p]; ok {
                close(ch)
                delete(pc.partitions, p)
            }
        }
    }
}

func (pc *ParallelConsumer) onLost(ctx context.Context, client *kgo.Client, lost map[string][]int32) {
    pc.onRevoked(ctx, client, lost)
}

func (pc *ParallelConsumer) partitionWorker(ctx context.Context, topic string, partition int32, ch <-chan *kgo.Record) {
    for {
        select {
        case <-ctx.Done():
            return
        case record, ok := <-ch:
            if !ok {
                return
            }
            // 处理记录
            _ = record
        }
    }
}
```

---

## 5. Visual Representations

### 5.1 Kafka Log Storage Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Kafka Log Storage Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                      Topic Partition Layout                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Topic: "orders"  Partition: 0                                         │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Segment 0 (Offset 0 - 999)                                      │  │  │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐       ┌─────────┐        │  │  │
│  │  │  │Batch 0  │  │Batch 1  │  │Batch 2  │  ...  │Batch n  │        │  │  │
│  │  │  │(0-99)   │  │(100-199)│  │(200-299)│       │(900-999)│        │  │  │
│  │  │  └─────────┘  └─────────┘  └─────────┘       └─────────┘        │  │  │
│  │  │                                                                    │  │  │
│  │  │  Offset Index:                                                     │  │  │
│  │  │  ┌─────────┬──────────┐                                            │  │  │
│  │  │  │Rel. Off │ Position │                                            │  │  │
│  │  │  ├─────────┼──────────┤                                            │  │  │
│  │  │  │    0    │    0     │                                            │  │  │
│  │  │  │  100    │   2048   │                                            │  │  │
│  │  │  │  200    │   4096   │                                            │  │  │
│  │  │  │  ...    │   ...    │                                            │  │  │
│  │  │  └─────────┴──────────┘                                            │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Segment 1 (Offset 1000 - 1999)                                  │  │  │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐       ┌─────────┐        │  │  │
│  │  │  │Batch 0  │  │Batch 1  │  │Batch 2  │  ...  │Batch n  │        │  │  │
│  │  │  │(1000-..)│  │         │  │         │       │         │        │  │  │
│  │  │  └─────────┘  └─────────┘  └─────────┘       └─────────┘        │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Active Segment (Offset 2000+)                                   │  │  │
│  │  │  (写入进行中，达到 1GB 或 7 天后滚动)                             │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Record Batch Structure (V2)                         │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Record Batch Header                                             │  │  │
│  │  │  ┌────────────────┬──────────────────────────────────────────┐ │  │  │
│  │  │  │ BaseOffset     │ int64    (8 bytes)                       │ │  │  │
│  │  │  │ BatchLength    │ int32    (4 bytes)                       │ │  │  │
│  │  │  │ PartitionLeaderEpoch │ int32 (4 bytes)                    │ │  │  │
│  │  │  │ Magic          │ int8     (1 byte, = 2 for V2)            │ │  │  │
│  │  │  │ CRC            │ int32    (4 bytes)                       │ │  │  │
│  │  │  │ Attributes     │ int16    (compression, timestamp type)   │ │  │  │
│  │  │  │ LastOffsetDelta│ int32    (relative to BaseOffset)        │ │  │  │
│  │  │  │ BaseTimestamp  │ int64    (first record timestamp)        │ │  │  │
│  │  │  │ MaxTimestamp   │ int64    (last record timestamp)         │ │  │  │
│  │  │  │ ProducerId     │ int64    (for idempotent producer)       │ │  │  │
│  │  │  │ ProducerEpoch  │ int16    (for transactions)              │ │  │  │
│  │  │  │ BaseSequence   │ int32    (for duplicate detection)       │ │  │  │
│  │  │  │ RecordsCount   │ int32    (number of records)             │ │  │  │
│  │  │  └────────────────┴──────────────────────────────────────────┘ │  │  │
│  │  │                                                                    │  │  │
│  │  │  Records...                                                        │  │  │
│  │  │                                                                    │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Kafka Exactly-Once Semantics Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Exactly-Once Semantics (EOS) Flow                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Idempotent Producer Flow                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Producer                          Broker (Partition)                  │  │
│  │                                                                        │  │
│  │  PID: 1000, Epoch: 0                                                   │  │
│  │     │                                                                  │  │
│  │     │  1. Produce: PID=1000, Epoch=0, Seq=0, Records=[A,B]             │  │
│  │     ├────────────────────────────────────────────────────────────►    │  │
│  │     │                                                                  │  │
│  │     │  2. Ack (retry: network timeout)                                 │  │
│  │     │◄─────────────────────────────────────────────────────────────    │  │
│  │     │     [LOST]                                                       │  │
│  │     │                                                                  │  │
│  │     │  3. Retry: PID=1000, Epoch=0, Seq=0, Records=[A,B]               │  │
│  │     ├────────────────────────────────────────────────────────────►    │  │
│  │     │                                                                  │  │
│  │     │  Broker Check:                                                   │  │
│  │     │  - PID=1000, Epoch=0 valid? YES                                  │  │
│  │     │  - Seq=0 already committed? YES                                  │  │
│  │     │  - Result: DEDUPLICATE, return last ack                          │  │
│  │     │                                                                  │  │
│  │     │  4. Duplicate Ack                                                │  │
│  │     │◄─────────────────────────────────────────────────────────────    │  │
│  │     │                                                                  │  │
│  │     │  5. Continue: PID=1000, Epoch=0, Seq=2, Records=[C,D]            │  │
│  │     ├────────────────────────────────────────────────────────────►    │  │
│  │     │                                                                  │  │
│  │     │  6. Ack: BaseOffset=2                                             │  │
│  │     │◄─────────────────────────────────────────────────────────────    │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Transactional Producer Flow                         │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Step 1: Initialize Transaction                                        │  │
│  │  ┌──────────────┐                                                      │  │
│  │  │ Producer     │                                                      │  │
│  │  │ TID: "app-1" │──InitTransactions()──► Transaction Coordinator       │  │
│  │  └──────────────┘                      (特殊分区 __transaction_state)   │  │
│  │                                                                        │  │
│  │  Step 2: Begin Transaction                                             │  │
│  │  ┌──────────────┐                                                      │  │
│  │  │ Producer     │                                                      │  │
│  │  │ PID: 1000    │──BeginTransaction()──► Start TX in coordinator      │  │
│  │  └──────────────┘                                                      │  │
│  │                                                                        │  │
│  │  Step 3: Send Messages to Partitions                                   │  │
│  │  ┌──────────────┐                                                      │  │
│  │  │ Partition 0  │◄────Produce(PID=1000, TID="app-1")────┐             │  │
│  │  └──────────────┘                                     │              │  │
│  │  ┌──────────────┐                                     │              │  │
│  │  │ Partition 1  │◄────Produce(PID=1000, TID="app-1")────┤             │  │
│  │  └──────────────┘                                     │              │  │
│  │  ┌──────────────┐                                     │              │  │
│  │  │ Partition 2  │◄────Produce(PID=1000, TID="app-1")────┘             │  │
│  │  └──────────────┘           Producer                   │              │  │
│  │                             (Transaction Marker)       │              │  │
│  │                                                                        │  │
│  │  Step 4: Commit Transaction                                            │  │
│  │  ┌──────────────┐                                                      │  │
│  │  │ Producer     │──CommitTransaction()──► Coordinator                  │  │
│  │  └──────────────┘                          │                           │  │
│  │                                            │  Write COMMIT markers     │  │
│  │                                            ▼                           │  │
│  │                            ┌───────────────────────────┐               │  │
│  │                            │  Partition 0: [A,B,COMMIT]│               │  │
│  │                            │  Partition 1: [D,E,COMMIT]│               │  │
│  │                            │  Partition 2: [F,G,COMMIT]│               │  │
│  │                            └───────────────────────────┘               │  │
│  │                                                                        │  │
│  │  Step 5: Consumer reads (isolation.level=read_committed)               │  │
│  │  ┌──────────────┐                                                      │  │
│  │  │ Consumer     │                                                      │  │
│  │  │              │──Poll()──► Sees: [A,B,D,E,F,G] (committed only)     │  │
│  │  │              │           Does NOT see: [X,Y] (aborted TX)          │  │
│  │  └──────────────┘                                                      │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │              Consumer-Producer Transaction (Consume-Transform-Produce) │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     │  │
│  │  │ Input    │────►│ Consumer │────►│ Transform│────►│ Producer │────►│  │
│  │  │ Topic    │     │ (Group)  │     │ Logic    │     │          │     │  │
│  │  └──────────┘     └──────────┘     └──────────┘     └──────────┘     │  │
│  │                          │                               │           │  │
│  │                          │  SendOffsetsToTransaction()   │           │  │
│  │                          └───────────────────────────────┘           │  │
│  │                                                                        │  │
│  │  Atomic Operation:                                                     │  │
│  │  - Input offsets committed                                             │  │
│  │  - Output records committed                                            │  │
│  │  - BOTH succeed or BOTH fail                                           │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Kafka Cluster Failover Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Kafka Cluster Failover & Replication                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Leader Failure Scenario                             │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Initial State:                                                        │  │
│  │  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐            │  │
│  │  │  Broker 1   │      │  Broker 2   │      │  Broker 3   │            │  │
│  │  │  (Leader)   │◄────►│  (Follower) │◄────►│  (Follower) │            │  │
│  │  │  HW: 100    │      │  HW: 100    │      │  HW: 98     │            │  │
│  │  │  LEO: 105   │      │  LEO: 105   │      │  LEO: 102   │            │  │
│  │  │  ISR: {1,2} │      │  In ISR     │      │  Out of ISR │            │  │
│  │  └─────────────┘      └─────────────┘      └─────────────┘            │  │
│  │         │                    │                    │                   │  │
│  │                                                                        │  │
│  │  Failure Event: Broker 1 Crashes                                       │  │
│  │  ═════════════════════════════════                                     │  │
│  │                                                                        │  │
│  │  Step 1: Controller detects failure (via zookeeper/session timeout)    │  │
│  │                                                                        │  │
│  │  Step 2: Controller selects new leader from ISR                        │  │
│  │          ┌─────────────────────────────────────────────┐               │  │
│  │          │  Leader Selection Criteria:                 │               │  │
│  │          │  1. Must be in ISR                          │               │  │
│  │          │  2. Highest LEO preferred                   │               │  │
│  │          │  → Broker 2 selected as new leader          │               │  │
│  │          └─────────────────────────────────────────────┘               │  │
│  │                                                                        │  │
│  │  Step 3: New Leader Election                                           │  │
│  │  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐            │  │
│  │  │  Broker 1   │      │  Broker 2   │      │  Broker 3   │            │  │
│  │  │  OFFLINE    │      │  (Leader)   │◄────►│  (Follower) │            │  │
│  │  │             │      │  HW: 100    │      │  HW: 98     │            │  │
│  │  │             │      │  LEO: 105   │      │  LEO: 102   │            │  │
│  │  │             │      │  ISR: {2,3} │      │  Truncating │            │  │
│  │  └─────────────┘      └─────────────┘      └─────────────┘            │  │
│  │                                                                        │  │
│  │  Step 4: ISR Recovery                                                  │  │
│  │  - Broker 3 truncates to HW: 98 (数据一致性)                            │  │
│  │  - Broker 3 fetches from offset 98 to catch up                         │  │
│  │  - When caught up, re-joins ISR                                        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │              Unclean Leader Election (Availability vs Consistency)     │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Scenario: All replicas in ISR are offline                             │  │
│  │                                                                        │  │
│  │  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐            │  │
│  │  │  Broker 1   │      │  Broker 2   │      │  Broker 3   │            │  │
│  │  │  (Leader)   │      │  OFFLINE    │      │  OFFLINE    │            │  │
│  │  │  OFFLINE    │      │  (in ISR)   │      │  (in ISR)   │            │  │
│  │  │  HW: 100    │      │             │      │             │            │  │
│  │  └─────────────┘      └─────────────┘      └─────────────┘            │  │
│  │                                                                        │  │
│  │  Option 1: unclean.leader.election.enable = FALSE (Default, Safe)     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│  │  │  Behavior: Partition remains unavailable                         │   │  │
│  │  │  Trade-off: Availability sacrificed for data consistency         │   │  │
│  │  │  Use case: Financial transactions, critical data                 │   │  │
│  │  └─────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                        │  │
│  │  Option 2: unclean.leader.election.enable = TRUE (Availability)       │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │  │
│  │  │  Behavior: Broker 3 (not in ISR) can be elected leader           │   │  │
│  │  │  Risk: Data loss (messages 99-102 on Broker 3 may be missing)    │   │  │
│  │  │  Use case: Log aggregation, metrics (some loss acceptable)       │   │  │
│  │  └─────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. References

1. **Kreps, J., Narkhede, N., & Rao, J.** (2011). Kafka: A Distributed Messaging System for Log Processing. *NetDB*.
2. **Kleppmann, M.** (2017). Designing Data-Intensive Applications. O'Reilly Media.
3. **Apache Kafka Documentation** (2024). kafka.apache.org/documentation
4. **Confluent** (2024). Kafka Internals. confluent.io/blog
5. **Wang, G., et al.** (2021). Building a High-Availability Distributed Streaming System. *VLDB*.

---

*Document Version: 1.0 | Last Updated: 2024*

---

## 10. Performance Benchmarking

### 10.1 Technology Stack Benchmarks

```go
package techstack_test

import (
	"context"
	"testing"
	"time"
)

// BenchmarkBasicOperation measures baseline performance
func BenchmarkBasicOperation(b *testing.B) {
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		// Simulate operation
	}
}

// BenchmarkConcurrentLoad tests concurrent operations
func BenchmarkConcurrentLoad(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate concurrent operation
			time.Sleep(1 * time.Microsecond)
		}
	})
}
```

### 10.2 Performance Characteristics

| Operation | Latency | Throughput | Resource Usage |
|-----------|---------|------------|----------------|
| **Simple** | 1ms | 1K RPS | Low |
| **Complex** | 10ms | 100 RPS | Medium |
| **Batch** | 100ms | 10K records | High |

### 10.3 Production Metrics

| Metric | Target | Alert | Critical |
|--------|--------|-------|----------|
| Latency p99 | < 100ms | > 200ms | > 500ms |
| Error Rate | < 0.1% | > 0.5% | > 1% |
| Throughput | > 1K | < 500 | < 100 |
| CPU Usage | < 70% | > 80% | > 95% |

### 10.4 Optimization Checklist

- [ ] Connection pooling configured
- [ ] Read replicas for read-heavy workloads
- [ ] Caching layer implemented
- [ ] Batch operations for bulk inserts
- [ ] Proper indexing strategy
- [ ] Query optimization completed
- [ ] Resource limits configured
