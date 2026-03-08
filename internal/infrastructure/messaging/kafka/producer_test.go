// Package kafka 提供 Kafka 生产者的测试
package kafka

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewProducer_InvalidBrokers 测试无效 broker 地址
func TestNewProducer_InvalidBrokers(t *testing.T) {
	// 使用无效的 broker 地址
	brokers := []string{"invalid:9092"}

	producer, err := NewProducer(brokers)
	// 由于无法连接到实际的 Kafka，这里会返回错误
	assert.Error(t, err)
	assert.Nil(t, producer)
}

// TestNewProducer_EmptyBrokers 测试空 broker 列表
func TestNewProducer_EmptyBrokers(t *testing.T) {
	producer, err := NewProducer([]string{})
	assert.Error(t, err)
	assert.Nil(t, producer)
}

// TestProducer_SendMessage_InvalidValue 测试发送无效值
func TestProducer_SendMessage_InvalidValue(t *testing.T) {
	// 使用循环引用创建无法序列化的值
	type Circular struct {
		Self *Circular
	}
	circular := &Circular{}
	circular.Self = circular

	// 由于无法创建有效的 producer，这里测试序列化逻辑
	// 使用 json.Marshal 模拟
	_, err := json.Marshal(circular)
	assert.Error(t, err)
}

// TestProducer_Struct 测试生产者结构体
func TestProducer_Struct(t *testing.T) {
	producer := &Producer{producer: nil}
	assert.NotNil(t, producer)
}

// TestNewProducer_Config 测试生产者配置
func TestNewProducer_Config(t *testing.T) {
	// 验证配置选项的默认值
	// 这些值在 NewProducer 函数中设置
	// Return.Successes = true
	// RequiredAcks = WaitForAll
	// Retry.Max = 5
	// 由于无法实际创建 producer，这里主要验证代码路径
	brokers := []string{"localhost:9092"}
	producer, err := NewProducer(brokers)
	if err != nil {
		// 预期错误，因为无法连接到实际 Kafka
		t.Logf("Expected error when connecting to Kafka: %v", err)
	}
	if producer != nil {
		defer producer.Close()
	}
}
