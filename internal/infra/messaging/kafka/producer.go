// Package kafka provides Kafka producer and consumer implementations for message queuing.
//
// Kafka 是一个分布式流处理平台，用于构建实时数据管道和流式应用。
//
// 设计原则：
// 1. 可靠性：使用同步生产者确保消息发送成功
// 2. 持久化：消息持久化到磁盘，支持消息重放
// 3. 可扩展性：支持水平扩展，处理大规模消息流
// 4. 容错性：支持消息重试和错误处理
//
// 核心功能：
// - Producer: 发送消息到 Kafka 主题
// - Consumer: 从 Kafka 主题消费消息
// - 消息序列化：支持 JSON 格式的消息序列化
//
// 使用场景：
// - 事件驱动架构：服务间异步通信
// - 日志聚合：收集和聚合应用日志
// - 流式处理：实时数据处理和分析
// - 消息队列：解耦生产者和消费者
//
// 示例：
//
//	// 创建生产者
//	producer, err := kafka.NewProducer([]string{"localhost:9092"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer producer.Close()
//
//	// 发送消息
//	ctx := context.Background()
//	err = producer.SendMessage(ctx, "my-topic", "key", map[string]interface{}{
//	    "message": "Hello Kafka",
//	})
package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

// Producer 是 Kafka 生产者的封装，用于发送消息到 Kafka 主题。
//
// 功能说明：
// - 使用同步生产者确保消息发送成功
// - 支持消息序列化（JSON 格式）
// - 支持消息键（Key）用于分区路由
//
// 设计说明：
// - 同步发送：等待消息发送成功后才返回
// - 可靠性：使用 WaitForAll 确保所有副本都收到消息
// - 自动重试：网络错误时自动重试
//
// 使用示例：
//
//	producer, err := kafka.NewProducer([]string{"localhost:9092"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer producer.Close()
//
//	ctx := context.Background()
//	err = producer.SendMessage(ctx, "my-topic", "message-key", messageData)
type Producer struct {
	producer sarama.SyncProducer
}

// NewProducer 创建并初始化 Kafka 生产者。
//
// 功能说明：
// - 连接到 Kafka 集群
// - 配置生产者选项（确认机制、重试策略等）
// - 返回配置好的生产者实例
//
// 参数：
// - brokers: Kafka Broker 地址列表
//   格式：[]string{"host1:9092", "host2:9092"}
//   至少提供一个 Broker 地址
//
// 返回：
// - *Producer: 配置好的生产者实例
// - error: 如果创建失败，返回错误信息
//
// 配置说明：
// - Return.Successes: 返回成功发送的消息（用于确认）
// - RequiredAcks: 确认机制
//   - WaitForAll: 等待所有副本确认（最可靠）
//   - WaitForLocal: 等待本地副本确认
//   - NoResponse: 不等待确认（最快但不可靠）
// - Retry.Max: 最大重试次数
//
// 使用示例：
//
//	// 连接到单个 Broker
//	producer, err := kafka.NewProducer([]string{"localhost:9092"})
//
//	// 连接到集群
//	producer, err := kafka.NewProducer([]string{
//	    "kafka1:9092",
//	    "kafka2:9092",
//	    "kafka3:9092",
//	})
//
// 注意事项：
// - 确保 Kafka 集群已启动并可访问
// - 生产环境建议使用集群配置（多个 Broker）
// - 应在应用程序生命周期中复用生产者实例
// - 退出前应调用 Close() 关闭生产者
func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true  // 返回成功发送的消息
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = 5            // 最大重试 5 次
	// 其他可选配置：
	// config.Producer.Timeout = 10 * time.Second
	// config.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer: producer}, nil
}

// SendMessage 发送消息到指定的 Kafka 主题。
//
// 功能说明：
// - 将消息序列化为 JSON 格式
// - 发送到指定的主题
// - 使用消息键（Key）进行分区路由
// - 同步等待发送结果
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消（当前未使用，保留用于未来扩展）
// - topic: Kafka 主题名称
// - key: 消息键（用于分区路由，相同键的消息会发送到同一分区）
// - value: 消息值（可以是任意类型，会自动序列化为 JSON）
//
// 返回：
// - error: 如果发送失败，返回错误信息
//
// 分区路由：
// - 如果提供了 key，消息会根据 key 的哈希值路由到特定分区
// - 相同 key 的消息会发送到同一分区，保证顺序
// - 如果 key 为空，消息会轮询分配到各个分区
//
// 使用示例：
//
//	ctx := context.Background()
//
//	// 发送简单消息
//	err := producer.SendMessage(ctx, "my-topic", "key1", "Hello Kafka")
//
//	// 发送结构化消息
//	message := map[string]interface{}{
//	    "user_id": 123,
//	    "action":  "login",
//	    "time":    time.Now(),
//	}
//	err = producer.SendMessage(ctx, "user-events", "user-123", message)
//
//	// 发送自定义结构体
//	type OrderEvent struct {
//	    OrderID int    `json:"order_id"`
//	    Status  string `json:"status"`
//	}
//	event := OrderEvent{OrderID: 1, Status: "created"}
//	err = producer.SendMessage(ctx, "orders", fmt.Sprintf("order-%d", event.OrderID), event)
//
// 注意事项：
// - 消息值必须是可序列化为 JSON 的类型
// - 同步发送会阻塞直到消息发送成功或失败
// - 如果发送失败，会自动重试（最多 5 次）
// - 生产环境建议使用异步生产者以提高吞吐量
func (p *Producer) SendMessage(ctx context.Context, topic string, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
		// 其他可选字段：
		// Headers: []sarama.RecordHeader{...}, // 消息头
		// Timestamp: time.Now(),               // 消息时间戳
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// Close 关闭 Kafka 生产者。
//
// 功能说明：
// - 关闭与 Kafka 集群的连接
// - 刷新所有待发送的消息
// - 释放生产者资源
//
// 返回：
// - error: 如果关闭过程中出现错误，返回错误信息
//
// 使用示例：
//
//	defer producer.Close()
//
// 注意事项：
// - 应在应用程序退出前调用
// - 关闭后不应再使用该生产者
// - 关闭会等待所有待发送的消息完成
func (p *Producer) Close() error {
	return p.producer.Close()
}
