# TS-011: Kafka 内部机制深度解析 (Apache Kafka Internals)

> **维度**: Technology Stack
> **级别**: S (17+ KB)
> **标签**: #kafka #streaming #log-structure #distributed-messaging
> **权威来源**: [Kafka Documentation](https://kafka.apache.org/documentation/), [Kafka: The Definitive Guide](https://www.oreilly.com/library/view/kafka-the-definitive/9781491936153/)
> **版本**: Kafka 4.0+

---

## Kafka 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Kafka Architecture                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      Kafka Cluster                                  │    │
│  │                                                                      │    │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │    │
│  │  │  Broker-1   │    │  Broker-2   │    │  Broker-3   │              │    │
│  │  │             │    │             │    │             │              │    │
│  │  │ ┌─────────┐ │    │ ┌─────────┐ │    │ ┌─────────┐ │              │    │
│  │  │ │Topic-A  │ │    │ │Topic-A  │ │    │ │Topic-A  │ │  Replica     │    │
│  │  │ │P0 (L)   │ │    │ │P0 (F)   │ │    │ │P0 (F)   │ │  Set         │    │
│  │  │ │P1 (F)   │ │    │ │P1 (L)   │ │    │ │P1 (F)   │ │  ISR={0,1,2} │    │
│  │  │ │P2 (F)   │ │    │ │P2 (F)   │ │    │ │P2 (L)   │ │              │    │
│  │  │ └─────────┘ │    │ └─────────┘ │    │ └─────────┘ │              │    │
│  │  └─────────────┘    └─────────────┘    └─────────────┘              │    │
│  │                                                                      │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │  ZooKeeper / KRaft (Metadata)                               │    │    │
│  │  │  - Broker 注册                                               │    │    │
│  │  │  - Topic/Partition 元数据                                    │    │    │
│  │  │  - Controller 选举 (Broker-1)                                │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  关键概念:                                                                   │
│  - Topic: 逻辑消息流                                                        │
│  - Partition: 物理分片，有序日志                                              │
│  - Replica: 分区副本，Leader/Follower                                         │
│  - ISR (In-Sync Replicas): 同步副本集合                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 日志存储

### 日志结构存储

```
Topic: order-events
Partition: 0

文件结构:
order-events-0/
├── 00000000000000000000.log      # 消息数据 (0 偏移开始)
├── 00000000000000000000.index    # 稀疏索引
├── 00000000000000000000.timeindex # 时间戳索引
├── 00000000000356897219.log      # 下一个日志段 (偏移 356897219)
├── 00000000000356897219.index
└── 00000000000356897219.timeindex

日志段 (Log Segment):
- 默认 1GB 或 7 天滚动
- 每个段独立索引

索引机制:
- 稀疏索引: 每 4KB 数据记录一个索引项
- 时间戳索引: 支持按时间查找
- 二分查找定位消息
```

### 消息格式 (Record Batch)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Record Batch Format (V2)                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Record Batch (Magic = 2)                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ BaseOffset: int64     (第一条消息偏移)                               │    │
│  │ BatchLength: int32    (批次长度)                                     │    │
│  │ PartitionLeaderEpoch: int32                                          │    │
│  │ Magic: int8 (2)                                                      │    │
│  │ CRC: int32                                                            │    │
│  │ Attributes: int16    (压缩类型、事务、控制批次)                        │    │
│  │ LastOffsetDelta: int32                                               │    │
│  │ FirstTimestamp: int64                                                │    │
│  │ MaxTimestamp: int64                                                  │    │
│  │ ProducerId: int64    (事务 ID)                                        │    │
│  │ ProducerEpoch: int16                                                 │    │
│  │ BaseSequence: int32  (序列号，用于幂等)                                │    │
│  │ RecordsCount: int32                                                  │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │ Records: []Record                                                   │    │
│  │                                                                     │    │
│  │ ┌─────────────────────────────────────────────────────────────────┐ │    │
│  │ │ Record:                                                         │ │    │
│  │ │ - Length: varint                                                │ │    │
│  │ │ - Attributes: int8                                              │ │    │
│  │ │ - TimestampDelta: varint                                        │ │    │
│  │ │ - OffsetDelta: varint                                           │ │    │
│  │ │ - KeyLen: varint, Key: []byte                                   │ │    │
│  │ │ - ValueLen: varint, Value: []byte                               │ │    │
│  │ │ - HeadersCount: varint                                          │ │    │
│  │ │ - Headers: []Header                                             │ │    │
│  │ └─────────────────────────────────────────────────────────────────┘ │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  优化特性:                                                                   │
│  - 批次压缩 (GZIP, Snappy, LZ4, ZSTD, Zstd)                                  │
│  - 零拷贝传输                                                               │
│  - 页缓存友好                                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 生产者机制

### 发送流程

```go
package kafka

import (
    "context"
    "github.com/IBM/sarama"  // Kafka 4.0+ 兼容
)

// ProducerConfig 生产者配置
type ProducerConfig struct {
    // 可靠性
    RequiredAcks sarama.RequiredAcks // WaitForLocal, WaitForAll, NoResponse
    Retries      int                // 重试次数

    // 批量
    BatchSize    int               // 批量大小
    BatchTimeout time.Duration     // 批量等待时间

    // 幂等性 (Kafka 4.0+ 默认开启)
    Idempotent bool

    // 事务
    Transactional bool
}

// 配置示例
func createProducer() (sarama.SyncProducer, error) {
    config := sarama.NewConfig()

    // 高可靠配置
    config.Producer.RequiredAcks = sarama.WaitForAll  // 所有 ISR 确认
    config.Producer.Retry.Max = 3
    config.Producer.Return.Successes = true

    // 幂等性 (确保消息不重复)
    config.Producer.Idempotent = true
    config.Net.MaxOpenRequests = 1

    // 压缩
    config.Producer.Compression = sarama.CompressionZSTD

    return sarama.NewSyncProducer([]string{"localhost:9092"}, config)
}

// 发送消息
func sendMessage(producer sarama.SyncProducer, topic string, key, value []byte) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.ByteEncoder(key),
        Value: sarama.ByteEncoder(value),
        Headers: []sarama.RecordHeader{
            {Key: []byte("version"), Value: []byte("1.0")},
        },
    }

    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        return err
    }

    log.Printf("Message sent to partition %d at offset %d", partition, offset)
    return nil
}
```

---

## 消费者机制

### 消费者组

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Consumer Group Rebalance                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Topic: order-events (6 partitions)                                         │
│                                                                              │
│  初始状态:                                                                   │
│  ┌───────────────┐      ┌───────────────┐      ┌───────────────┐           │
│  │ Consumer-1    │      │ Consumer-2    │      │ Consumer-3    │           │
│  │ P0, P1        │      │ P2, P3        │      │ P4, P5        │           │
│  └───────────────┘      └───────────────┘      └───────────────┘           │
│                                                                              │
│  Consumer-2 离开后 (Rebalance):                                              │
│  ┌───────────────┐                          ┌───────────────┐              │
│  │ Consumer-1    │                          │ Consumer-3    │              │
│  │ P0, P1, P2    │                          │ P3, P4, P5    │              │
│  └───────────────┘                          └───────────────┘              │
│                                                                              │
│  分区分配策略:                                                                │
│  - Range: 按范围分配 (默认)                                                   │
│  - RoundRobin: 轮询                                                          │
│  - Sticky: 最小化重新分配                                                     │
│  - CooperativeSticky: Kafka 4.0+ 推荐 (增量再平衡)                             │
│                                                                              │
│  Rebalance 协议:                                                              │
│  1. 消费者加入组 (JoinGroup)                                                  │
│  2. 选举 Leader (协调器)                                                      │
│  3. Leader 计算分区分配 (SyncGroup)                                           │
│  4. 所有消费者接收分配                                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 偏移管理

```go
// 自动提交 (默认)
config.Consumer.Offsets.AutoCommit.Enable = true
config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

// 手动提交 (精确控制)
config.Consumer.Offsets.AutoCommit.Enable = false

func consumeWithManualCommit(consumer sarama.ConsumerGroup) {
    handler := ConsumerGroupHandler{
        commitFunc: func(session sarama.ConsumerGroupSession) {
            session.Commit() // 手动提交偏移
        },
    }
    consumer.Consume(context.Background(), []string{"topic"}, handler)
}

// 存储偏移到 Kafka (__consumer_offsets topic)
// 或外部存储 (Redis, DB) 用于 Exactly-Once 语义
```

---

## 性能优化

| 优化项 | 配置 | 说明 |
|--------|------|------|
| 批量发送 | `linger.ms=5`, `batch.size=16384` | 提高吞吐量 |
| 压缩 | `compression.type=zstd` | 减少网络/磁盘 |
| 零拷贝 | 自动启用 | 传输优化 |
| 页缓存 | `log.flush.interval.messages=10000` | 刷盘策略 |
| 分区数 | 根据消费者数 | 并行度 |

---

## 参考文献

1. [Kafka Documentation](https://kafka.apache.org/documentation/)
2. [Kafka: The Definitive Guide](https://www.oreilly.com/library/view/kafka-the-definitive/9781491936153/)
3. [KIP-500: Replace ZooKeeper with a Self-Managed Metadata Quorum](https://cwiki.apache.org/confluence/display/KAFKA/KIP-500)
