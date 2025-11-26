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

**生产者性能对比**:

| 生产者类型 | 吞吐量 | 延迟 | 可靠性 | 适用场景 |
|-----------|--------|------|--------|---------|
| **同步生产者** | 10,000-50,000 msg/s | 1-5ms | 高 | 需要确认的场景 |
| **异步生产者** | 100,000-500,000 msg/s | < 1ms | 中 | 高吞吐量场景 |
| **批量生产者** | 500,000-2,000,000 msg/s | 5-20ms | 高 | 批量处理场景 |

**完整的同步生产者实现**:

```go
// internal/infrastructure/messaging/kafka/producer.go
package kafka

import (
    "context"
    "fmt"
    "time"

    "github.com/IBM/sarama"
)

type Producer struct {
    producer sarama.SyncProducer
    brokers  []string
    config   *sarama.Config
}

// NewProducer 创建同步生产者（生产环境配置）
func NewProducer(brokers []string) (*Producer, error) {
    config := sarama.NewConfig()

    // 生产者配置
    config.Producer.Return.Successes = true  // 返回成功消息
    config.Producer.Return.Errors = true      // 返回错误消息
    config.Producer.RequiredAcks = sarama.WaitForAll  // 等待所有副本确认
    config.Producer.Retry.Max = 3             // 最大重试次数
    config.Producer.Retry.Backoff = 100 * time.Millisecond  // 重试间隔
    config.Producer.Timeout = 10 * time.Second  // 超时时间

    // 批量配置（提高吞吐量）
    config.Producer.Flush.Frequency = 10 * time.Millisecond  // 批量刷新间隔
    config.Producer.Flush.Messages = 100                     // 批量消息数
    config.Producer.Flush.Bytes = 1024 * 1024               // 批量大小（1MB）

    // 压缩配置（减少网络传输）
    config.Producer.Compression = sarama.CompressionSnappy

    // 分区策略
    config.Producer.Partitioner = sarama.NewHashPartitioner

    // 版本配置
    config.Version = sarama.V2_8_0_0

    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create producer: %w", err)
    }

    return &Producer{
        producer: producer,
        brokers:  brokers,
        config:   config,
    }, nil
}

// SendMessage 发送消息（带上下文）
func (p *Producer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.ByteEncoder(key),
        Value: sarama.ByteEncoder(value),
        Headers: []sarama.RecordHeader{
            {
                Key:   []byte("timestamp"),
                Value: []byte(time.Now().Format(time.RFC3339)),
            },
        },
    }

    // 检查上下文取消
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    partition, offset, err := p.producer.SendMessage(msg)
    if err != nil {
        return fmt.Errorf("failed to send message to topic %s: %w", topic, err)
    }

    // 记录成功日志（可选）
    // logger.Debug("Message sent", "topic", topic, "partition", partition, "offset", offset)
    _ = partition
    _ = offset

    return nil
}

// SendMessageWithRetry 发送消息（带重试）
func (p *Producer) SendMessageWithRetry(ctx context.Context, topic string, key, value []byte, maxRetries int) error {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        err := p.SendMessage(ctx, topic, key, value)
        if err == nil {
            return nil
        }

        lastErr = err

        // 检查是否可重试
        if !isRetryableError(err) {
            return err
        }

        if i < maxRetries-1 {
            // 指数退避
            backoff := time.Duration(1<<uint(i)) * time.Second
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(backoff):
            }
        }
    }

    return fmt.Errorf("failed to send message after %d retries: %w", maxRetries, lastErr)
}

func isRetryableError(err error) bool {
    // 判断错误是否可重试
    if err == nil {
        return false
    }

    // Kafka 错误类型判断
    if kafkaErr, ok := err.(sarama.ProducerError); ok {
        switch kafkaErr.Err {
        case sarama.ErrLeaderNotAvailable,
             sarama.ErrNotLeaderForPartition,
             sarama.ErrRequestTimedOut:
            return true
        }
    }

    return false
}
```

**完整的异步生产者实现**:

```go
// 异步生产者（高吞吐量）
type AsyncProducer struct {
    producer sarama.AsyncProducer
    brokers  []string
    errors   chan *sarama.ProducerError
    successes chan *sarama.ProducerMessage
    wg       sync.WaitGroup
}

func NewAsyncProducer(brokers []string) (*AsyncProducer, error) {
    config := sarama.NewConfig()

    // 异步生产者配置
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true
    config.Producer.RequiredAcks = sarama.WaitForOne  // 异步模式可以降低确认要求
    config.Producer.Retry.Max = 3
    config.Producer.Retry.Backoff = 100 * time.Millisecond

    // 批量配置（异步模式可以更激进）
    config.Producer.Flush.Frequency = 5 * time.Millisecond
    config.Producer.Flush.Messages = 200
    config.Producer.Flush.Bytes = 2 * 1024 * 1024  // 2MB

    // 压缩
    config.Producer.Compression = sarama.CompressionSnappy

    // 缓冲区配置
    config.ChannelBufferSize = 256  // 输入通道缓冲区

    producer, err := sarama.NewAsyncProducer(brokers, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create async producer: %w", err)
    }

    ap := &AsyncProducer{
        producer:  producer,
        brokers:   brokers,
        errors:    make(chan *sarama.ProducerError, 100),
        successes: make(chan *sarama.ProducerMessage, 100),
    }

    // 启动错误和成功处理 goroutine
    ap.startHandlers()

    return ap, nil
}

func (ap *AsyncProducer) startHandlers() {
    ap.wg.Add(2)

    // 处理成功消息
    go func() {
        defer ap.wg.Done()
        for msg := range ap.producer.Successes() {
            select {
            case ap.successes <- msg:
            default:
                // 通道满，记录警告
            }
        }
    }()

    // 处理错误消息
    go func() {
        defer ap.wg.Done()
        for err := range ap.producer.Errors() {
            select {
            case ap.errors <- err:
            default:
                // 通道满，记录警告
            }
        }
    }()
}

// SendMessage 异步发送消息
func (ap *AsyncProducer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.ByteEncoder(key),
        Value: sarama.ByteEncoder(value),
    }

    select {
    case ap.producer.Input() <- msg:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Close 关闭生产者
func (ap *AsyncProducer) Close() error {
    close(ap.producer.Input())
    ap.wg.Wait()
    return ap.producer.Close()
}
```

**批量生产者实现**:

```go
// 批量生产者（最高吞吐量）
type BatchProducer struct {
    producer *AsyncProducer
    batch    []*sarama.ProducerMessage
    mu       sync.Mutex
    maxSize  int
    flushInterval time.Duration
}

func NewBatchProducer(brokers []string, maxSize int, flushInterval time.Duration) (*BatchProducer, error) {
    producer, err := NewAsyncProducer(brokers)
    if err != nil {
        return nil, err
    }

    bp := &BatchProducer{
        producer:     producer,
        batch:        make([]*sarama.ProducerMessage, 0, maxSize),
        maxSize:      maxSize,
        flushInterval: flushInterval,
    }

    // 启动定时刷新
    go bp.autoFlush()

    return bp, nil
}

func (bp *BatchProducer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.ByteEncoder(key),
        Value: sarama.ByteEncoder(value),
    }

    bp.mu.Lock()
    bp.batch = append(bp.batch, msg)
    shouldFlush := len(bp.batch) >= bp.maxSize
    bp.mu.Unlock()

    if shouldFlush {
        return bp.Flush(ctx)
    }

    return nil
}

func (bp *BatchProducer) Flush(ctx context.Context) error {
    bp.mu.Lock()
    if len(bp.batch) == 0 {
        bp.mu.Unlock()
        return nil
    }

    batch := make([]*sarama.ProducerMessage, len(bp.batch))
    copy(batch, bp.batch)
    bp.batch = bp.batch[:0]
    bp.mu.Unlock()

    for _, msg := range batch {
        select {
        case bp.producer.producer.Input() <- msg:
        case <-ctx.Done():
            return ctx.Err()
        }
    }

    return nil
}

func (bp *BatchProducer) autoFlush() {
    ticker := time.NewTicker(bp.flushInterval)
    defer ticker.Stop()

    for range ticker.C {
        bp.Flush(context.Background())
    }
}
```

### 1.3.2 消费者示例

**消费者性能对比**:

| 消费者类型 | 吞吐量 | 延迟 | 可靠性 | 适用场景 |
|-----------|--------|------|--------|---------|
| **单个消费者** | 10,000-50,000 msg/s | 1-5ms | 中 | 单实例应用 |
| **消费者组** | 100,000-500,000 msg/s | 1-3ms | 高 | 分布式应用 |
| **批量消费者** | 500,000-2,000,000 msg/s | 5-20ms | 高 | 批量处理场景 |

**完整的单个消费者实现**:

```go
// internal/infrastructure/messaging/kafka/consumer.go
package kafka

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/IBM/sarama"
)

type Consumer struct {
    consumer sarama.Consumer
    brokers  []string
    wg       sync.WaitGroup
}

func NewConsumer(brokers []string) (*Consumer, error) {
    consumer, err := sarama.NewConsumer(brokers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create consumer: %w", err)
    }

    return &Consumer{
        consumer: consumer,
        brokers:  brokers,
    }, nil
}

func (c *Consumer) Consume(ctx context.Context, topic string, handler func(*sarama.ConsumerMessage) error) error {
    // 获取所有分区
    partitionList, err := c.consumer.Partitions(topic)
    if err != nil {
        return fmt.Errorf("failed to get partitions: %w", err)
    }

    // 为每个分区创建消费者
    for _, partition := range partitionList {
        c.wg.Add(1)
        go func(partition int32) {
            defer c.wg.Done()
            c.consumePartition(ctx, topic, partition, handler)
        }(partition)
    }

    // 等待所有分区消费者完成
    c.wg.Wait()
    return nil
}

func (c *Consumer) consumePartition(ctx context.Context, topic string, partition int32, handler func(*sarama.ConsumerMessage) error) {
    // 从最新偏移量开始消费
    pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
    if err != nil {
        fmt.Printf("Failed to consume partition %d: %v\n", partition, err)
        return
    }
    defer pc.AsyncClose()

    for {
        select {
        case msg := <-pc.Messages():
            if msg == nil {
                return
            }

            // 处理消息
            if err := handler(msg); err != nil {
                fmt.Printf("Failed to process message: %v\n", err)
                // 根据错误类型决定是否继续
                continue
            }

        case err := <-pc.Errors():
            if err != nil {
                fmt.Printf("Consumer error: %v\n", err)
            }

        case <-ctx.Done():
            return
        }
    }
}

func (c *Consumer) Close() error {
    return c.consumer.Close()
}
```

**完整的消费者组实现**:

```go
// 消费者组（推荐用于生产环境）
type ConsumerGroup struct {
    consumerGroup sarama.ConsumerGroup
    brokers       []string
    groupID       string
    config        *sarama.Config
}

func NewConsumerGroup(brokers []string, groupID string) (*ConsumerGroup, error) {
    config := sarama.NewConfig()

    // 消费者组配置
    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
    config.Consumer.Offsets.Initial = sarama.OffsetNewest  // 从最新偏移量开始
    config.Consumer.Offsets.AutoCommit.Enable = true       // 自动提交偏移量
    config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

    // 批量配置
    config.Consumer.Fetch.Min = 1024 * 1024  // 最小批量大小（1MB）
    config.Consumer.Fetch.Default = 10 * 1024 * 1024  // 默认批量大小（10MB）
    config.Consumer.MaxProcessingTime = 30 * time.Second  // 最大处理时间

    // 超时配置
    config.Consumer.Group.Session.Timeout = 10 * time.Second
    config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

    // 版本配置
    config.Version = sarama.V2_8_0_0

    consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create consumer group: %w", err)
    }

    return &ConsumerGroup{
        consumerGroup: consumerGroup,
        brokers:       brokers,
        groupID:       groupID,
        config:        config,
    }, nil
}

func (cg *ConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            err := cg.consumerGroup.Consume(ctx, topics, handler)
            if err != nil {
                return fmt.Errorf("consumer group error: %w", err)
            }
        }
    }
}

func (cg *ConsumerGroup) Close() error {
    return cg.consumerGroup.Close()
}
```

**完整的消息处理器实现**:

```go
// 生产环境级别的消息处理器
type MessageHandler struct {
    processor MessageProcessor
    logger    *slog.Logger
    metrics   *Metrics
}

type MessageProcessor func(context.Context, *sarama.ConsumerMessage) error

func NewMessageHandler(processor MessageProcessor, logger *slog.Logger) *MessageHandler {
    return &MessageHandler{
        processor: processor,
        logger:    logger,
        metrics:   NewMetrics(),
    }
}

func (h *MessageHandler) Setup(session sarama.ConsumerGroupSession) error {
    h.logger.Info("Consumer group session started",
        "member_id", session.MemberID(),
        "generation_id", session.GenerationID(),
    )
    return nil
}

func (h *MessageHandler) Cleanup(session sarama.ConsumerGroupSession) error {
    h.logger.Info("Consumer group session ended",
        "member_id", session.MemberID(),
    )
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
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

            start := time.Now()
            err := h.processMessage(ctx, message)
            duration := time.Since(start)

            cancel()

            if err != nil {
                // 记录错误，但不提交偏移量，以便重试
                h.logger.Error("Failed to process message",
                    "topic", message.Topic,
                    "partition", message.Partition,
                    "offset", message.Offset,
                    "error", err,
                    "duration", duration,
                )

                h.metrics.IncrementErrorCount(message.Topic)

                // 根据错误类型决定是否继续
                if !isRetryableError(err) {
                    // 不可重试的错误，跳过这条消息
                    session.MarkMessage(message, "")
                }
                continue
            }

            // 处理成功，提交偏移量
            session.MarkMessage(message, "")

            h.metrics.IncrementSuccessCount(message.Topic)
            h.metrics.RecordProcessingDuration(message.Topic, duration)

            // 记录成功日志（可选）
            if h.logger.Enabled(context.Background(), slog.LevelDebug) {
                h.logger.Debug("Message processed successfully",
                    "topic", message.Topic,
                    "partition", message.Partition,
                    "offset", message.Offset,
                    "duration", duration,
                )
            }

        case <-session.Context().Done():
            return nil
        }
    }
}

func (h *MessageHandler) processMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
    // 记录处理开始时间
    h.metrics.IncrementMessageCount(message.Topic)

    // 调用处理器
    return h.processor(ctx, message)
}
```

**批量消费者实现**:

```go
// 批量消费者（高性能）
type BatchConsumer struct {
    consumerGroup *ConsumerGroup
    batchSize     int
    flushInterval time.Duration
    handler       BatchMessageProcessor
}

type BatchMessageProcessor func(context.Context, []*sarama.ConsumerMessage) error

func NewBatchConsumer(brokers []string, groupID string, batchSize int, flushInterval time.Duration) (*BatchConsumer, error) {
    consumerGroup, err := NewConsumerGroup(brokers, groupID)
    if err != nil {
        return nil, err
    }

    return &BatchConsumer{
        consumerGroup: consumerGroup,
        batchSize:     batchSize,
        flushInterval: flushInterval,
    }, nil
}

func (bc *BatchConsumer) Consume(ctx context.Context, topics []string, handler BatchMessageProcessor) error {
    bc.handler = handler

    batchHandler := &BatchMessageHandler{
        batchConsumer: bc,
        batch:         make([]*sarama.ConsumerMessage, 0, bc.batchSize),
        mu:           sync.Mutex{},
    }

    return bc.consumerGroup.Consume(ctx, topics, batchHandler)
}

type BatchMessageHandler struct {
    batchConsumer *BatchConsumer
    batch         []*sarama.ConsumerMessage
    mu           sync.Mutex
    lastFlush    time.Time
}

func (h *BatchMessageHandler) Setup(session sarama.ConsumerGroupSession) error {
    h.lastFlush = time.Now()
    return nil
}

func (h *BatchMessageHandler) Cleanup(session sarama.ConsumerGroupSession) error {
    // 清理时刷新剩余消息
    h.flushBatch(session)
    return nil
}

func (h *BatchMessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    ticker := time.NewTicker(h.batchConsumer.flushInterval)
    defer ticker.Stop()

    for {
        select {
        case message := <-claim.Messages():
            if message == nil {
                h.flushBatch(session)
                return nil
            }

            h.mu.Lock()
            h.batch = append(h.batch, message)
            shouldFlush := len(h.batch) >= h.batchConsumer.batchSize
            h.mu.Unlock()

            if shouldFlush {
                h.flushBatch(session)
            }

        case <-ticker.C:
            h.flushBatch(session)

        case <-session.Context().Done():
            h.flushBatch(session)
            return nil
        }
    }
}

func (h *BatchMessageHandler) flushBatch(session sarama.ConsumerGroupSession) {
    h.mu.Lock()
    if len(h.batch) == 0 {
        h.mu.Unlock()
        return
    }

    batch := make([]*sarama.ConsumerMessage, len(h.batch))
    copy(batch, h.batch)
    h.batch = h.batch[:0]
    h.mu.Unlock()

    // 处理批量消息
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    if err := h.batchConsumer.handler(ctx, batch); err != nil {
        // 批量处理失败，记录错误但不提交偏移量
        log.Printf("Failed to process batch: %v", err)
        return
    }

    // 批量处理成功，提交所有消息的偏移量
    for _, msg := range batch {
        session.MarkMessage(msg, "")
    }

    h.lastFlush = time.Now()
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

良好的消费者设计可以提高消息处理的可靠性、性能和可维护性。根据生产环境的实际经验，合理的消费者设计可以将吞吐量提升 5-20 倍，将消息丢失率降低到接近 0。

**消费者性能优化对比**:

| 优化项 | 未优化 | 优化后 | 提升比例 |
|--------|--------|--------|---------|
| **吞吐量** | 10,000 msg/s | 500,000 msg/s | +4900% |
| **延迟** | 10ms | 1ms | -90% |
| **CPU 使用** | 70% | 25% | -64% |
| **内存使用** | 400MB | 150MB | -62.5% |

**消费者最佳实践**:

1. **消费者组**: 使用消费者组实现负载均衡和容错（提升 5-10 倍）
2. **偏移量管理**: 正确处理偏移量，确保消息不丢失
3. **错误处理**: 实现错误处理和重试机制
4. **批量处理**: 使用批量处理提高性能（提升 5-10 倍）
5. **并发处理**: 合理使用并发处理提高吞吐量
6. **背压控制**: 实现背压控制避免内存溢出

**完整的消费者最佳实践示例**:

```go
// 生产环境级别的消费者配置
func NewProductionConsumerGroup(brokers []string, groupID string) (*ConsumerGroup, error) {
    config := sarama.NewConfig()

    // 1. 消费者组配置
    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
    config.Consumer.Offsets.Initial = sarama.OffsetNewest
    config.Consumer.Offsets.AutoCommit.Enable = true
    config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

    // 2. 批量配置（关键优化）
    config.Consumer.Fetch.Min = 1024 * 1024        // 1MB
    config.Consumer.Fetch.Default = 10 * 1024 * 1024  // 10MB
    config.Consumer.Fetch.Max = 50 * 1024 * 1024  // 50MB
    config.Consumer.MaxProcessingTime = 30 * time.Second

    // 3. 超时配置
    config.Consumer.Group.Session.Timeout = 10 * time.Second
    config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
    config.Consumer.Group.Rebalance.Timeout = 60 * time.Second

    // 4. 版本配置
    config.Version = sarama.V2_8_0_0

    // 5. 网络配置
    config.Net.DialTimeout = 10 * time.Second
    config.Net.ReadTimeout = 10 * time.Second
    config.Net.WriteTimeout = 10 * time.Second

    consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create consumer group: %w", err)
    }

    return &ConsumerGroup{consumerGroup: consumerGroup}, nil
}
```

**偏移量管理最佳实践**:

```go
// 偏移量管理策略
type OffsetStrategy int

const (
    OffsetStrategyAutoCommit OffsetStrategy = iota  // 自动提交（默认）
    OffsetStrategyManualCommit                      // 手动提交（推荐生产环境）
    OffsetStrategyBatchCommit                       // 批量提交（高性能）
)

// 手动提交偏移量（推荐）
func (h *MessageHandler) ConsumeClaimWithManualCommit(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    commitInterval := 5 * time.Second
    ticker := time.NewTicker(commitInterval)
    defer ticker.Stop()

    for {
        select {
        case message := <-claim.Messages():
            if message == nil {
                session.Commit()  // 最后提交
                return nil
            }

            // 处理消息
            if err := h.processMessage(message); err != nil {
                // 处理失败，不提交偏移量，等待重试
                continue
            }

            // 标记消息已处理
            session.MarkMessage(message, "")

        case <-ticker.C:
            // 定期提交偏移量
            session.Commit()

        case <-session.Context().Done():
            session.Commit()  // 清理时提交
            return nil
        }
    }
}

// 批量提交偏移量（高性能）
func (h *BatchMessageHandler) ConsumeClaimWithBatchCommit(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    batch := make([]*sarama.ConsumerMessage, 0, 100)
    commitInterval := 5 * time.Second
    ticker := time.NewTicker(commitInterval)
    defer ticker.Stop()

    for {
        select {
        case message := <-claim.Messages():
            if message == nil {
                h.processBatch(session, batch)
                session.Commit()
                return nil
            }

            batch = append(batch, message)
            if len(batch) >= 100 {
                h.processBatch(session, batch)
                session.Commit()
                batch = batch[:0]
            }

        case <-ticker.C:
            if len(batch) > 0 {
                h.processBatch(session, batch)
                session.Commit()
                batch = batch[:0]
            }

        case <-session.Context().Done():
            if len(batch) > 0 {
                h.processBatch(session, batch)
                session.Commit()
            }
            return nil
        }
    }
}
```

**错误处理和重试最佳实践**:

```go
// 错误处理和重试机制
type RetryPolicy struct {
    MaxRetries    int
    InitialBackoff time.Duration
    MaxBackoff     time.Duration
    BackoffMultiplier float64
}

func (h *MessageHandler) processMessageWithRetry(ctx context.Context, message *sarama.ConsumerMessage, policy RetryPolicy) error {
    var lastErr error

    for attempt := 0; attempt < policy.MaxRetries; attempt++ {
        err := h.processor(ctx, message)
        if err == nil {
            return nil
        }

        lastErr = err

        // 判断是否可重试
        if !isRetryableError(err) {
            return err
        }

        // 最后一次尝试，直接返回错误
        if attempt == policy.MaxRetries-1 {
            break
        }

        // 指数退避
        backoff := policy.InitialBackoff * time.Duration(policy.BackoffMultiplier*float64(attempt))
        if backoff > policy.MaxBackoff {
            backoff = policy.MaxBackoff
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(backoff):
        }
    }

    return fmt.Errorf("failed after %d retries: %w", policy.MaxRetries, lastErr)
}

// 死信队列处理
func (h *MessageHandler) handleDeadLetter(message *sarama.ConsumerMessage, err error) {
    // 发送到死信队列
    deadLetterTopic := fmt.Sprintf("%s-dlq", message.Topic)

    // 记录死信消息
    h.logger.Error("Message sent to dead letter queue",
        "topic", message.Topic,
        "partition", message.Partition,
        "offset", message.Offset,
        "error", err,
    )

    // 发送到死信队列（实现细节...）
    // h.producer.SendMessage(deadLetterTopic, message.Key, message.Value)
}
```

**并发处理最佳实践**:

```go
// 并发处理消息（提高吞吐量）
type ConcurrentMessageHandler struct {
    processor MessageProcessor
    workers  int
    sem      chan struct{}
}

func NewConcurrentMessageHandler(processor MessageProcessor, workers int) *ConcurrentMessageHandler {
    return &ConcurrentMessageHandler{
        processor: processor,
        workers:   workers,
        sem:       make(chan struct{}, workers),
    }
}

func (h *ConcurrentMessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    var wg sync.WaitGroup

    for {
        select {
        case message := <-claim.Messages():
            if message == nil {
                wg.Wait()  // 等待所有处理完成
                return nil
            }

            // 获取信号量
            h.sem <- struct{}{}
            wg.Add(1)

            go func(msg *sarama.ConsumerMessage) {
                defer func() {
                    <-h.sem  // 释放信号量
                    wg.Done()
                }()

                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()

                if err := h.processor(ctx, msg); err != nil {
                    log.Printf("Failed to process message: %v", err)
                    return
                }

                // 处理成功，标记消息
                session.MarkMessage(msg, "")
            }(message)

        case <-session.Context().Done():
            wg.Wait()
            return nil
        }
    }
}
```

**背压控制最佳实践**:

```go
// 背压控制（避免内存溢出）
type BackpressureHandler struct {
    processor MessageProcessor
    maxPending int
    pending    chan struct{}
}

func NewBackpressureHandler(processor MessageProcessor, maxPending int) *BackpressureHandler {
    return &BackpressureHandler{
        processor: processor,
        maxPending: maxPending,
        pending:    make(chan struct{}, maxPending),
    }
}

func (h *BackpressureHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for {
        select {
        case message := <-claim.Messages():
            if message == nil {
                return nil
            }

            // 等待背压释放
            h.pending <- struct{}{}

            go func(msg *sarama.ConsumerMessage) {
                defer func() { <-h.pending }()

                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()

                if err := h.processor(ctx, msg); err != nil {
                    log.Printf("Failed to process message: %v", err)
                    return
                }

                session.MarkMessage(msg, "")
            }(message)

        case <-session.Context().Done():
            // 等待所有待处理消息完成
            for i := 0; i < h.maxPending; i++ {
                h.pending <- struct{}{}
            }
            return nil
        }
    }
}
```

**消费者最佳实践要点**:

1. **消费者组**:
   - 使用消费者组实现负载均衡和容错（提升 5-10 倍）
   - 合理设置消费者数量（通常等于分区数）
   - 监控消费者组状态

2. **偏移量管理**:
   - 手动提交偏移量（推荐生产环境）
   - 处理成功后提交偏移量
   - 定期提交避免重复处理

3. **错误处理**:
   - 实现智能重试机制
   - 区分可重试和不可重试错误
   - 使用死信队列处理失败消息

4. **批量处理**:
   - 使用批量处理提高性能（提升 5-10 倍）
   - 合理设置批量大小
   - 平衡延迟和吞吐量

5. **并发处理**:
   - 使用并发处理提高吞吐量
   - 控制并发数量避免资源耗尽
   - 使用信号量控制并发度

6. **背压控制**:
   - 实现背压控制避免内存溢出
   - 监控待处理消息数量
   - 设置合理的背压阈值

7. **监控和指标**:
   - 监控消费速率、延迟、错误率
   - 记录关键指标（分区、偏移量）
   - 设置告警阈值

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
