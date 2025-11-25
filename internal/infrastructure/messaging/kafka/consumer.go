package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

// Consumer Kafka 消费者
type Consumer struct {
	consumer sarama.ConsumerGroup
}

// Config 配置
type Config struct {
	Brokers []string
	GroupID string
}

// Handler 消息处理函数
type Handler func(ctx context.Context, topic string, key string, value []byte) error

// NewConsumer 创建消费者
func NewConsumer(cfg Config) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
	}, nil
}

// Consume 消费消息
func (c *Consumer) Consume(ctx context.Context, topics []string, handler Handler) error {
	consumerHandler := &consumerGroupHandler{
		handler: handler,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := c.consumer.Consume(ctx, topics, consumerHandler); err != nil {
				return fmt.Errorf("failed to consume: %w", err)
			}
		}
	}
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	if c.consumer != nil {
		return c.consumer.Close()
	}
	return nil
}

// consumerGroupHandler 消费者组处理器
type consumerGroupHandler struct {
	handler Handler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

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
			value := message.Value

			if err := h.handler(session.Context(), message.Topic, key, value); err != nil {
				return err
			}

			session.MarkMessage(message, "")
		}
	}
}
