package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

// MessageHandler 是消息处理函数的类型定义。
//
// 功能说明：
// - 处理从 Kafka 接收到的消息
// - 接收消息键和值（原始字节）
// - 返回处理错误（如果处理失败）
//
// 参数：
// - ctx: 上下文，用于控制处理超时和取消
// - key: 消息键（字符串格式）
// - value: 消息值（原始字节，需要自行反序列化）
//
// 返回：
// - error: 如果处理失败，返回错误信息
//   返回错误会导致消息处理失败，可能需要重试
//
// 使用示例：
//
//	handler := func(ctx context.Context, key string, value []byte) error {
//	    var data map[string]interface{}
//	    if err := json.Unmarshal(value, &data); err != nil {
//	        return err
//	    }
//	    // 处理消息
//	    return processMessage(data)
//	}
type MessageHandler func(ctx context.Context, key string, value []byte) error

// Consumer 是 Kafka 消费者的封装，用于从 Kafka 主题消费消息。
//
// 功能说明：
// - 使用消费者组（Consumer Group）消费消息
// - 支持多个消费者实例负载均衡
// - 自动管理消息偏移量（Offset）
// - 支持重平衡（Rebalance）处理
//
// 设计说明：
// - 消费者组：多个消费者实例可以组成一个消费者组
// - 负载均衡：消费者组内的消费者会分配不同的分区
// - 偏移量管理：自动提交消息偏移量，确保消息不重复消费
//
// 使用示例：
//
//	handler := func(ctx context.Context, key string, value []byte) error {
//	    log.Printf("Received message: key=%s, value=%s", key, string(value))
//	    return nil
//	}
//
//	consumer, err := kafka.NewConsumer(
//	    []string{"localhost:9092"},
//	    "my-consumer-group",
//	    handler,
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer consumer.Close()
//
//	ctx := context.Background()
//	err = consumer.Consume(ctx, []string{"my-topic"})
type Consumer struct {
	consumer sarama.ConsumerGroup
	handler  MessageHandler
}

// NewConsumer 创建并初始化 Kafka 消费者。
//
// 功能说明：
// - 连接到 Kafka 集群
// - 加入指定的消费者组
// - 配置消费者选项（重平衡策略、初始偏移量等）
//
// 参数：
// - brokers: Kafka Broker 地址列表
//   格式：[]string{"host1:9092", "host2:9092"}
// - groupID: 消费者组 ID
//   相同 groupID 的消费者会组成一个消费者组
//   用于负载均衡和消息分发
// - handler: 消息处理函数
//   当收到消息时，会调用此函数处理消息
//
// 返回：
// - *Consumer: 配置好的消费者实例
// - error: 如果创建失败，返回错误信息
//
// 配置说明：
// - Rebalance.Strategy: 重平衡策略
//   - RoundRobin: 轮询分配（默认）
//   - Range: 范围分配
//   - Sticky: 粘性分配（减少重平衡）
// - Offsets.Initial: 初始偏移量
//   - OffsetNewest: 从最新消息开始消费（默认）
//   - OffsetOldest: 从最早消息开始消费
//
// 使用示例：
//
//	// 创建消费者
//	handler := func(ctx context.Context, key string, value []byte) error {
//	    // 处理消息
//	    return nil
//	}
//	consumer, err := kafka.NewConsumer(
//	    []string{"localhost:9092"},
//	    "my-consumer-group",
//	    handler,
//	)
//
// 注意事项：
// - 确保 Kafka 集群已启动并可访问
// - 消费者组 ID 应具有业务意义，便于管理
// - 多个消费者实例使用相同的 groupID 可以实现负载均衡
// - 应在应用程序生命周期中复用消费者实例
func NewConsumer(brokers []string, groupID string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin() // 轮询分配策略
	config.Consumer.Offsets.Initial = sarama.OffsetNewest                            // 从最新消息开始
	// 其他可选配置：
	// config.Consumer.Offsets.AutoCommit.Enable = true  // 自动提交偏移量
	// config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	// config.Consumer.MaxProcessingTime = 30 * time.Second

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		handler:  handler,
	}, nil
}

// Consume 开始消费指定主题的消息。
//
// 功能说明：
// - 订阅指定的主题列表
// - 从 Kafka 接收消息并调用处理函数
// - 方法会阻塞，直到上下文取消或发生错误
//
// 参数：
// - ctx: 上下文，用于控制消费过程
//   取消上下文会停止消费
// - topics: 要订阅的主题列表
//
// 返回：
// - error: 如果消费过程中发生错误，返回错误信息
//
// 工作流程：
// 1. 订阅指定的主题
// 2. 加入消费者组并分配分区
// 3. 从分配的分区接收消息
// 4. 调用消息处理函数处理消息
// 5. 标记消息已处理（提交偏移量）
//
// 使用示例：
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	// 在 goroutine 中消费消息
//	go func() {
//	    if err := consumer.Consume(ctx, []string{"topic1", "topic2"}); err != nil {
//	        log.Printf("Consumer error: %v", err)
//	    }
//	}()
//
//	// 在需要时停止消费
//	cancel()
//
// 注意事项：
// - 方法会阻塞，应在单独的 goroutine 中运行
// - 使用上下文控制消费的开始和停止
// - 如果处理函数返回错误，消息处理会失败
// - 消费者组会自动管理分区分配和重平衡
func (c *Consumer) Consume(ctx context.Context, topics []string) error {
	handler := &consumerGroupHandler{handler: c.handler}

	// 持续消费消息，直到上下文取消
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 开始消费，这个方法会阻塞直到发生重平衡或错误
			err := c.consumer.Consume(ctx, topics, handler)
			if err != nil {
				return fmt.Errorf("consumer error: %w", err)
			}
		}
	}
}

// Close 关闭 Kafka 消费者。
//
// 功能说明：
// - 离开消费者组
// - 关闭与 Kafka 集群的连接
// - 释放消费者资源
//
// 返回：
// - error: 如果关闭过程中出现错误，返回错误信息
//
// 使用示例：
//
//	defer consumer.Close()
//
// 注意事项：
// - 应在应用程序退出前调用
// - 关闭后不应再使用该消费者
// - 关闭会触发重平衡，其他消费者会接管分区
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

// consumerGroupHandler 是消费者组处理器的实现。
//
// 功能说明：
// - 实现 sarama.ConsumerGroupHandler 接口
// - 处理消费者组的生命周期事件
// - 处理接收到的消息
//
// 生命周期：
// 1. Setup: 消费者加入消费者组时调用
// 2. ConsumeClaim: 处理分配的分区中的消息
// 3. Cleanup: 消费者离开消费者组时调用
type consumerGroupHandler struct {
	handler MessageHandler
}

// Setup 在消费者加入消费者组时调用。
//
// 功能说明：
// - 消费者成功加入消费者组后调用
// - 可以在此处进行初始化操作
//
// 参数：
// - session: 消费者组会话
//
// 返回：
// - error: 如果设置失败，返回错误信息
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 在消费者离开消费者组时调用。
//
// 功能说明：
// - 消费者离开消费者组前调用
// - 可以在此处进行清理操作
//
// 参数：
// - session: 消费者组会话
//
// 返回：
// - error: 如果清理失败，返回错误信息
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 处理分配的分区中的消息。
//
// 功能说明：
// - 从分配的分区接收消息
// - 调用消息处理函数处理消息
// - 标记消息已处理（提交偏移量）
//
// 参数：
// - session: 消费者组会话，用于标记消息和提交偏移量
// - claim: 分区声明，包含分配的分区中的消息
//
// 返回：
// - error: 如果处理失败，返回错误信息
//   返回错误会导致消费者离开消费者组并触发重平衡
//
// 工作流程：
// 1. 从分区声明中接收消息
// 2. 调用消息处理函数处理消息
// 3. 如果处理成功，标记消息已处理
// 4. 如果处理失败，返回错误（可能导致重试）
//
// 注意事项：
// - 消息处理应该是幂等的（可以安全地重复处理）
// - 处理函数返回错误会导致消息处理失败
// - 标记消息后，偏移量会在下次提交时更新
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			// 会话已取消，停止处理
			return nil
		case message := <-claim.Messages():
			// 接收消息
			if message == nil {
				// 分区已关闭，退出
				continue
			}

			// 处理消息
			key := string(message.Key)
			if err := h.handler(session.Context(), key, message.Value); err != nil {
				// 处理失败，返回错误
				// 这会导致消费者离开消费者组并触发重平衡
				return err
			}

			// 标记消息已处理
			// 偏移量会在下次自动提交时更新
			session.MarkMessage(message, "")
		}
	}
}
