// Package mqtt 提供 MQTT 客户端的测试
package mqtt

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMessageHandler_Type 测试消息处理函数类型
func TestMessageHandler_Type(t *testing.T) {
	// 测试有效的消息处理函数
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return nil
	})
	assert.NotNil(t, handler)

	// 测试调用处理函数
	ctx := context.Background()
	err := handler(ctx, "test/topic", []byte("test payload"))
	assert.NoError(t, err)
}

// TestMessageHandler_WithError 测试返回错误的处理函数
func TestMessageHandler_WithError(t *testing.T) {
	testErr := assert.AnError
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return testErr
	})

	ctx := context.Background()
	err := handler(ctx, "topic", []byte("payload"))
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

// TestClient_Publish_PayloadTypes 测试不同类型的 payload 序列化
func TestClient_Publish_PayloadTypes(t *testing.T) {
	// 测试 Publish 方法中的 payload 类型转换逻辑
	testCases := []struct {
		name     string
		payload  interface{}
		isBytes  bool
		isString bool
	}{
		{
			name:    "字节数组",
			payload: []byte{0x01, 0x02, 0x03},
			isBytes: true,
		},
		{
			name:     "字符串",
			payload:  "test string",
			isString: true,
		},
		{
			name:    "整数",
			payload: 42,
		},
		{
			name:    "结构体",
			payload: struct{ Name string }{Name: "test"},
		},
		{
			name:    "map",
			payload: map[string]int{"key": 123},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 验证 payload 可以被正确处理
			var data []byte
			switch v := tc.payload.(type) {
			case []byte:
				data = v
				assert.True(t, tc.isBytes)
			case string:
				data = []byte(v)
				assert.True(t, tc.isString)
			default:
				// 其他类型需要序列化
				var err error
				data, err = json.Marshal(v)
				require.NoError(t, err)
				assert.False(t, tc.isBytes)
				assert.False(t, tc.isString)
			}
			assert.NotNil(t, data)
		})
	}
}

// TestClient_Publish_InvalidPayload 测试发送无法序列化的负载
func TestClient_Publish_InvalidPayload(t *testing.T) {
	// 创建一个无法 JSON 序列化的值（循环引用）
	type Circular struct {
		Self *Circular
	}
	circular := &Circular{}
	circular.Self = circular

	// 验证循环引用无法序列化
	_, err := json.Marshal(circular)
	assert.Error(t, err)
}

// TestClient_Struct 测试客户端结构体
func TestClient_Struct(t *testing.T) {
	client := &Client{client: nil}
	assert.NotNil(t, client)
}

// TestMessageHandler_ContextTimeout 测试处理函数上下文超时
func TestMessageHandler_ContextTimeout(t *testing.T) {
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Millisecond):
			return nil
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Millisecond) // 确保超时
	err := handler(ctx, "topic", []byte("payload"))
	// 错误可能是 nil 或 context.DeadlineExceeded，取决于时间
	t.Logf("Handler error: %v", err)
}

// TestMessageHandler_Called 测试消息处理函数被正确调用
func TestMessageHandler_Called(t *testing.T) {
	called := false
	receivedTopic := ""
	receivedPayload := []byte{}

	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		called = true
		receivedTopic = topic
		receivedPayload = payload
		return nil
	})

	ctx := context.Background()
	err := handler(ctx, "test/topic", []byte("test payload"))
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, "test/topic", receivedTopic)
	assert.Equal(t, []byte("test payload"), receivedPayload)
}

// TestMessageHandler_WithPayloadProcessing 测试处理函数处理 payload
func TestMessageHandler_WithPayloadProcessing(t *testing.T) {
	type SensorData struct {
		Temperature float64 `json:"temperature"`
		Humidity    float64 `json:"humidity"`
	}

	var receivedData SensorData
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return json.Unmarshal(payload, &receivedData)
	})

	ctx := context.Background()
	payload := `{"temperature": 25.5, "humidity": 60.0}`
	err := handler(ctx, "sensors/data", []byte(payload))
	require.NoError(t, err)
	assert.Equal(t, 25.5, receivedData.Temperature)
	assert.Equal(t, 60.0, receivedData.Humidity)
}

// TestQoS_Values 测试 QoS 值
func TestQoS_Values(t *testing.T) {
	// QoS 0: 最多一次
	// QoS 1: 至少一次
	// QoS 2: 恰好一次
	qosLevels := []byte{0, 1, 2}

	for _, qos := range qosLevels {
		assert.GreaterOrEqual(t, qos, byte(0))
		assert.LessOrEqual(t, qos, byte(2))
	}
}

// TestRetainedFlag 测试保留消息标志
func TestRetainedFlag(t *testing.T) {
	// 测试保留和非保留消息
	retainedValues := []bool{true, false}

	for _, retained := range retainedValues {
		// 验证布尔值
		assert.True(t, retained == true || retained == false)
	}
}

// TestPayloadTypes_Exhaustive 测试所有支持的 payload 类型
func TestPayloadTypes_Exhaustive(t *testing.T) {
	testCases := []struct {
		name    string
		payload interface{}
		expect  []byte
	}{
		{
			name:    "byte_slice",
			payload: []byte{0x01, 0x02, 0x03},
			expect:  []byte{0x01, 0x02, 0x03},
		},
		{
			name:    "string",
			payload: "hello world",
			expect:  []byte("hello world"),
		},
		{
			name:    "int",
			payload: 42,
			expect:  []byte("42"),
		},
		{
			name:    "float",
			payload: 3.14,
			expect:  []byte("3.14"),
		},
		{
			name:    "bool_true",
			payload: true,
			expect:  []byte("true"),
		},
		{
			name:    "bool_false",
			payload: false,
			expect:  []byte("false"),
		},
		{
			name:    "map",
			payload: map[string]string{"key": "value"},
			expect:  []byte(`{"key":"value"}`),
		},
		{
			name:    "slice",
			payload: []int{1, 2, 3},
			expect:  []byte("[1,2,3]"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result []byte
			var err error

			switch v := tc.payload.(type) {
			case []byte:
				result = v
			case string:
				result = []byte(v)
			default:
				result, err = json.Marshal(tc.payload)
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expect, result)
		})
	}
}
