package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

// MessageHandler 消息处理函数
type MessageHandler func(ctx context.Context, key string, value []byte) error

// Consumer Kafka 消费者
type Consumer struct {
	consumer sarama.ConsumerGroup
	handler  MessageHandler
}

// NewConsumer 创建 Kafka 消费者
func NewConsumer(brokers []string, groupID string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		handler:  handler,
	}, nil
}

// Consume 消费消息
func (c *Consumer) Consume(ctx context.Context, topics []string) error {
	handler := &consumerGroupHandler{handler: c.handler}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := c.consumer.Consume(ctx, topics, handler)
			if err != nil {
				return fmt.Errorf("consumer error: %w", err)
			}
		}
	}
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

// consumerGroupHandler 消费者组处理器
type consumerGroupHandler struct {
	handler MessageHandler
}

// Setup 设置会话
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 清理会话
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费消息
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return nil
		case message := <-claim.Messages():
			if message == nil {
				continue
			}

			key := string(message.Key)
			if err := h.handler(session.Context(), key, message.Value); err != nil {
				return err
			}

			session.MarkMessage(message, "")
		}
	}
}
