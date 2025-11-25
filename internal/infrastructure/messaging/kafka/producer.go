package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

// Producer Kafka 生产者
type Producer struct {
	producer sarama.SyncProducer
}

// Config 配置
type Config struct {
	Brokers []string
}

// NewProducer 创建生产者
func NewProducer(cfg Config) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
	}, nil
}

// SendMessage 发送消息
func (p *Producer) SendMessage(ctx context.Context, topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// Close 关闭生产者
func (p *Producer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}
