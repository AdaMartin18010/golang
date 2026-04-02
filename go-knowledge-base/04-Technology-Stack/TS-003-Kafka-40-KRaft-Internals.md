# TS-003: Kafka 4.0 KRaft 内部机制 (Kafka 4.0 KRaft Internals)

> **维度**: Technology Stack
> **级别**: S (20+ KB)
> **标签**: #kafka40 #kraft #raft #consensus #zookeeper-removal
> **权威来源**: [KIP-500](https://cwiki.apache.org/confluence/display/KAFKA/KIP-500), [Kafka 4.0 Release Notes](https://kafka.apache.org/documentation/#upgrade_4_0_0)

---

## KRaft 演进

```
Kafka 2.8 (2021)         Kafka 3.3 (2022)          Kafka 4.0 (2026) ⭐️
      │                          │                          │
      ▼                          ▼                          ▼
┌─────────────┐          ┌───────────────┐          ┌─────────────────┐
│ KRaft       │          │ KRaft         │          │ KRaft GA        │
│ Early Access│─────────►│ Production    │─────────►│ ZooKeeper      │
│             │          │ Ready         │          │ Removed         │
└─────────────┘          └───────────────┘          └─────────────────┘
      │                          │                          │
      • ZK 依赖                   • 支持两种模式              • 仅 KRaft
      • 双写                      • ZK 逐渐废弃               • 全新架构
                                   • 迁移工具                  • 更高性能
```

---

## KRaft 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Kafka 4.0 KRaft Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Legacy (ZK Mode)                    KRaft Mode (Kafka 4.0)                  │
│  ─────────────────                   ─────────────────────                   │
│                                                                              │
│  ┌─────────┐                         ┌─────────────┐                        │
│  │ZooKeeper│◄───────────────────────►│ Controller  │                        │
│  │Quorum   │  元数据管理               │ Quorum      │                        │
│  │(3-5节点)│                         │ (3+节点)    │                        │
│  └────┬────┘                         └──────┬──────┘                        │
│       │                                      │                               │
│       │  会话管理、配置                        │  元数据复制 (Raft)            │
│       │  ACL、ISR管理                         │  控制器选举                   │
│       │                                      │  配置管理                      │
│       │                                      │                               │
│  ┌────┴────┐                            ┌────┴────┐                        │
│  │ Brokers │                            │ Brokers │                        │
│  │         │                            │         │                        │
│  │ Pull ZK │                            │ Pull    │                        │
│  │  data   │                            │ metadata│                        │
│  └─────────┘                            │ from    │                        │
│                                         │ Quorum  │                        │
│                                         └─────────┘                        │
│                                                                              │
│  优势：                                                                       │
│  1. 单系统维护（不再需要 ZK）                                                 │
│  2. 更高性能（分区数扩展到 2M+）                                               │
│  3. 更快控制器选举（Raft vs ZK）                                              │
│  4. 更简单的部署                                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## KRaft 元数据复制

### 日志结构

```go
// KRaft 使用 Raft 共识复制元数据

// 元数据日志条目
{
    "offset": 100,
    "term": 5,
    "timestamp": 1712345678901,
    "type": "PARTITION_CHANGE",
    "data": {
        "topic": "orders",
        "partition": 0,
        "leader": 1,
        "isr": [1, 2, 3],
        "ar": [1, 2, 3, 4],
        "replicas": [1, 2, 3, 4]
    }
}

// 快照（定期压缩）
{
    "offset": 10000,
    "data": {
        "topics": [...],
        "partitions": [...],
        "configs": [...],
        "acls": [...]
    }
}
```

### 元数据状态机

```go
// 控制器应用日志到状态机
type MetadataStateMachine struct {
    topics      map[string]Topic
    partitions  map[PartitionID]Partition
    configs     map[ConfigKey]Config
    acls        map[ACLKey]ACL

    // 内存中维护的索引
    leaderIsr   map[PartitionID]LeaderAndISR
}

// 应用日志条目
func (sm *MetadataStateMachine) apply(record *Record) error {
    switch record.Type {
    case TOPIC_CHANGE:
        return sm.applyTopicChange(record)
    case PARTITION_CHANGE:
        return sm.applyPartitionChange(record)
    case CONFIG_CHANGE:
        return sm.applyConfigChange(record)
    case ACL_CHANGE:
        return sm.applyACLChange(record)
    default:
        return fmt.Errorf("unknown record type: %v", record.Type)
    }
}

// ISR 变更处理
func (sm *MetadataStateMachine) applyPartitionChange(record *Record) {
    partition := sm.partitions[record.PartitionID]
    partition.Leader = record.Leader
    partition.ISR = record.ISR
    partition.AR = record.AR  // Adding Replicas

    // 更新 leader 索引
    sm.leaderIsr[record.PartitionID] = LeaderAndISR{
        Leader: record.Leader,
        ISR:    record.ISR,
    }
}
```

---

## 控制器选举

### 与 ZK 模式对比

| 特性 | ZK 模式 | KRaft 模式 |
|------|---------|-----------|
| 选举算法 | ZAB (类似 Paxos) | Raft |
| 选举时间 | ~3-5s | ~100-500ms |
| 脑裂风险 | 有 (ZK 分区) | 无 (Quorum 保证) |
| 元数据延迟 | ~10-100ms | ~1-10ms |
| 最大分区数 | ~50K | ~2M+ |

### Raft 实现

```go
// KRaft Controller 使用 Raft

type Controller struct {
    raft *RaftNode

    // 当前状态
    state ControllerState

    // 元数据日志
    log *MetadataLog
}

// 成为 Leader 后
func (c *Controller) onElectedLeader() {
    log.Info("Became controller leader")

    // 1. 恢复元数据状态
    c.restoreMetadataState()

    // 2. 启动分区重平衡
    c.triggerRebalance()

    // 3. 开始 ISR 监控
    c.startISRMonitoring()
}

// 处理 Broker 注册
func (c *Controller) handleBrokerRegistration(req *RegisterBrokerRequest) {
    record := &Record{
        Type: BROKER_REGISTRATION,
        Data: BrokerRecord{
            ID:        req.BrokerID,
            Host:      req.Host,
            Port:      req.Port,
            Rack:      req.Rack,
            Timestamp: time.Now(),
        },
    }

    // 复制到 Quorum
    c.raft.Append(record)
}
```

---

## 元数据传播

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Metadata Propagation                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Controller Quorum                     Broker                               │
│  ─────────────────                     ──────                               │
│                                                                              │
│  ┌─────────────┐                       ┌─────────────┐                     │
│  │ Leader      │◄────FetchMetadata────│             │                     │
│  │             │       (pull)         │  Metadata   │                     │
│  │ Log:        │                       │  Cache      │                     │
│  │ offset 100  │────MetadataResponse──►│             │                     │
│  │ ...         │                       │  Applied:   │                     │
│  │ offset 200  │                       │  offset 95  │                     │
│  │ (latest)    │                       │             │                     │
│  └─────────────┘                       └─────────────┘                     │
│                                                                              │
│  增量更新：                                                                   │
│  1. Broker 记录上次获取的 offset                                             │
│  2. 下次只获取增量 (offset 96-200)                                           │
│  3. 应用到本地 Metadata Cache                                                │
│                                                                              │
│  全量更新（滞后太多）：                                                        │
│  1. 发送 snapshot 请求                                                       │
│  2. Leader 发送完整状态快照                                                   │
│  3. Broker 重建本地状态                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Kafka 4.0 新特性

### 1. 性能提升

```
• 分区数限制：50K → 2M+
• 控制器选举：3-5s → 100-500ms
• 元数据传播：~100ms → ~10ms
• 启动时间：分钟级 → 秒级
```

### 2. 移除 ZooKeeper

```yaml
# Kafka 4.0 配置 (server.properties)
# 不再需要 ZK 配置

# 旧配置（已移除）
# zookeeper.connect=localhost:2181

# 新配置
process.roles=broker,controller
node.id=1
controller.quorum.voters=1@localhost:9093,2@localhost:9094,3@localhost:9095

listeners=PLAINTEXT://localhost:9092,CONTROLLER://localhost:9093
```

### 3. 新的管理 API

```bash
# Kafka 4.0 使用 Admin API 替代 ZK 脚本

# 创建 topic
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic orders --partitions 3

# 查看 ISR
kafka-metadata-quorum.sh --bootstrap-server localhost:9092 --describe --status

# 动态配置
kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name orders --alter --add-config retention.ms=86400000
```

---

## 参考文献

1. [KIP-500: Replace ZooKeeper with a Self-Managed Metadata Quorum](https://cwiki.apache.org/confluence/display/KAFKA/KIP-500)
2. [Kafka 4.0 Release Notes](https://kafka.apache.org/documentation/#upgrade_4_0_0)
3. [Raft Consensus Algorithm](https://raft.github.io/)
4. [Running Kafka without ZooKeeper](https://kafka.apache.org/documentation/#kraft)
