# 1. 💬 Kafka (Sarama) 深度解析

> **简介**: 本文档详细阐述了 Kafka (Sarama) 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 💬 Kafka (Sarama) 深度解析](#1--kafka-sarama-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 生产者示例](#131-生产者示例)
    - [1.3.2 消费者示例](#132-消费者示例)
    - [1.3.3 分区策略](#133-分区策略)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 生产者最佳实践](#141-生产者最佳实践)
    - [1.4.2 消费者最佳实践](#142-消费者最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Kafka 是什么？**

Kafka 是一个分布式流处理平台，用于构建实时数据管道和流式应用。

**核心特性**:

- ✅ **高吞吐量**: 支持高吞吐量的消息处理
- ✅ **持久化**: 消息持久化存储
- ✅ **分区**: 支持分区和并行处理
- ✅ **顺序保证**: 保证分区内消息顺序

---

## 1.2 选型论证

**为什么选择 Kafka？**

**论证矩阵**:

| 评估维度 | 权重 | Kafka | RabbitMQ | NATS | Redis Streams | 说明 |
|---------|------|-------|----------|------|---------------|------|
| **吞吐量** | 30% | 10 | 6 | 9 | 8 | Kafka 吞吐量最高 |
| **持久化** | 25% | 10 | 9 | 6 | 7 | Kafka 持久化最可靠 |
| **分区支持** | 20% | 10 | 7 | 8 | 6 | Kafka 分区功能最完善 |
| **生态集成** | 15% | 10 | 8 | 7 | 6 | Kafka 生态最丰富 |
| **学习成本** | 10% | 6 | 8 | 8 | 9 | Kafka 学习曲线较陡 |
| **加权总分** | - | **9.20** | 7.50 | 7.80 | 7.30 | Kafka 得分最高 |

**核心优势**:

1. **吞吐量（权重 30%）**:
   - 高吞吐量，适合大数据流处理场景
   - 支持批量处理，提高效率
   - 分区并行处理，性能优秀

2. **持久化（权重 25%）**:
   - 消息持久化存储，保证可靠性
   - 支持消息重放，适合事件溯源
   - 数据保留策略灵活

3. **分区支持（权重 20%）**:
   - 完善的分区机制，支持顺序保证
   - 支持消费者组，实现负载均衡
   - 支持动态分区扩展

**为什么不选择其他消息队列？**

1. **RabbitMQ**:
   - ✅ 功能丰富，适合传统消息队列场景
   - ❌ 吞吐量不如 Kafka
   - ❌ 不适合大数据流处理
   - ❌ 分区支持不如 Kafka

2. **NATS**:
   - ✅ 性能优秀，延迟低
   - ❌ 持久化支持有限
   - ❌ 不适合大数据流处理
   - ❌ 生态不如 Kafka 丰富

3. **Redis Streams**:
   - ✅ 简单易用，与 Redis 集成好
   - ❌ 持久化不如 Kafka 可靠
   - ❌ 不适合大数据流处理
   - ❌ 功能不如 Kafka 完善

---

## 1.3 实际应用

### 1.3.1 生产者示例

**同步生产者**:

```go
// internal/infrastructure/messaging/kafka/producer.go
package kafka

import (
    "github.com/IBM/sarama"
)

type Producer struct {
    producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll

    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    return &Producer{producer: producer}, nil
}

func (p *Producer) SendMessage(topic string, key, value []byte) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.ByteEncoder(key),
        Value: sarama.ByteEncoder(value),
    }

    _, _, err := p.producer.SendMessage(msg)
    return err
}
```

**异步生产者**:

```go
// 异步生产者
func NewAsyncProducer(brokers []string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true

    producer, err := sarama.NewAsyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    // 处理成功和错误
    go func() {
        for {
            select {
            case success := <-producer.Successes():
                logger.Info("Message sent",
                    "topic", success.Topic,
                    "partition", success.Partition,
                    "offset", success.Offset,
                )
            case err := <-producer.Errors():
                logger.Error("Failed to send message", "error", err)
            }
        }
    }()

    return &Producer{producer: producer}, nil
}
```

### 1.3.2 消费者示例

**单个消费者**:

```go
// internal/infrastructure/messaging/kafka/consumer.go
package kafka

type Consumer struct {
    consumer sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        return nil, err
    }

    return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) Consume(topic string, handler func(*sarama.ConsumerMessage)) error {
    partitionList, err := c.consumer.Partitions(topic)
    if err != nil {
        return err
    }

    for partition := range partitionList {
        pc, err := c.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
        if err != nil {
            return err
        }

        go func(pc sarama.PartitionConsumer) {
            defer pc.AsyncClose()
            for msg := range pc.Messages() {
                handler(msg)
            }
        }(pc)
    }

    return nil
}
```

**消费者组**:

```go
// 消费者组
func NewConsumerGroup(brokers []string, groupID string) (*ConsumerGroup, error) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
    config.Consumer.Offsets.Initial = sarama.OffsetNewest

    consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, err
    }

    return &ConsumerGroup{consumerGroup: consumerGroup}, nil
}

func (cg *ConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
    for {
        err := cg.consumerGroup.Consume(ctx, topics, handler)
        if err != nil {
            return err
        }
    }
}
```

### 1.3.3 分区策略

**分区策略示例**:

```go
// 自定义分区策略
type CustomPartitioner struct{}

func (p *CustomPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
    // 根据 key 进行分区
    if message.Key != nil {
        hash := crc32.ChecksumIEEE(message.Key.(sarama.ByteEncoder))
        return int32(hash) % numPartitions, nil
    }

    // 轮询分区
    return rand.Int31n(numPartitions), nil
}

// 使用自定义分区策略
config.Producer.Partitioner = sarama.NewCustomPartitioner(&CustomPartitioner{})
```

---

## 1.4 最佳实践

### 1.4.1 生产者最佳实践

**为什么需要良好的生产者设计？**

良好的生产者设计可以提高消息发送的可靠性和性能。

**生产者最佳实践**:

1. **批量发送**: 使用批量发送提高吞吐量
2. **错误处理**: 正确处理发送错误，实现重试机制
3. **分区策略**: 选择合适的分区策略
4. **消息序列化**: 使用高效的序列化方式

**实际应用示例**:

```go
// 生产者最佳实践
func (p *Producer) SendMessageWithRetry(topic string, key, value []byte, maxRetries int) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.ByteEncoder(key),
        Value: sarama.ByteEncoder(value),
    }

    for i := 0; i < maxRetries; i++ {
        _, _, err := p.producer.SendMessage(msg)
        if err == nil {
            return nil
        }

        if i < maxRetries-1 {
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }

    return fmt.Errorf("failed to send message after %d retries", maxRetries)
}
```

**最佳实践要点**:

1. **批量发送**: 使用批量发送提高吞吐量
2. **错误处理**: 实现重试机制，处理发送失败
3. **分区策略**: 根据业务需求选择合适的分区策略
4. **消息序列化**: 使用高效的序列化方式（如 Protocol Buffers）

### 1.4.2 消费者最佳实践

**为什么需要良好的消费者设计？**

良好的消费者设计可以提高消息处理的可靠性和性能。

**消费者最佳实践**:

1. **消费者组**: 使用消费者组实现负载均衡
2. **偏移量管理**: 正确处理偏移量，避免消息丢失
3. **错误处理**: 实现错误处理和重试机制
4. **批量处理**: 使用批量处理提高性能

**实际应用示例**:

```go
// 消费者最佳实践
type MessageHandler struct{}

func (h *MessageHandler) Setup(sarama.ConsumerGroupSession) error {
    return nil
}

func (h *MessageHandler) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

func (h *MessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for {
        select {
        case message := <-claim.Messages():
            if message == nil {
                return nil
            }

            // 处理消息
            if err := h.processMessage(message); err != nil {
                // 记录错误，但不提交偏移量，以便重试
                logger.Error("Failed to process message",
                    "topic", message.Topic,
                    "partition", message.Partition,
                    "offset", message.Offset,
                    "error", err,
                )
                continue
            }

            // 提交偏移量
            session.MarkMessage(message, "")
        case <-session.Context().Done():
            return nil
        }
    }
}
```

**最佳实践要点**:

1. **消费者组**: 使用消费者组实现负载均衡和容错
2. **偏移量管理**: 正确处理偏移量，确保消息不丢失
3. **错误处理**: 实现错误处理和重试机制
4. **批量处理**: 使用批量处理提高性能

---

## 📚 扩展阅读

- [Kafka 官方文档](https://kafka.apache.org/)
- [Sarama 官方文档](https://github.com/IBM/sarama)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Kafka (Sarama) 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
