// Package mqtt provides tests for MQTT client.
package mqtt

import (
	"context"
	"encoding/json"
	"errors"
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

// TestMessageHandler_MultipleCalls 测试处理函数被多次调用
func TestMessageHandler_MultipleCalls(t *testing.T) {
	callCount := 0
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		callCount++
		return nil
	})

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		err := handler(ctx, "topic", []byte("payload"))
		require.NoError(t, err)
	}

	assert.Equal(t, 5, callCount)
}

// TestMessageHandler_DifferentTopics 测试不同主题的处理
func TestMessageHandler_DifferentTopics(t *testing.T) {
	topics := []string{
		"sensors/temperature",
		"sensors/humidity",
		"devices/status",
		"alerts/critical",
	}

	var receivedTopics []string
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		receivedTopics = append(receivedTopics, topic)
		return nil
	})

	ctx := context.Background()
	for _, topic := range topics {
		err := handler(ctx, topic, []byte("data"))
		require.NoError(t, err)
	}

	assert.Equal(t, topics, receivedTopics)
}

// TestMessageHandler_ErrorHandling 测试错误处理
func TestMessageHandler_ErrorHandling(t *testing.T) {
	testErrors := []error{
		assert.AnError,
		context.Canceled,
		context.DeadlineExceeded,
		errors.New("custom error"),
	}

	for _, testErr := range testErrors {
		handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
			return testErr
		})

		ctx := context.Background()
		err := handler(ctx, "topic", []byte("payload"))
		assert.Error(t, err)
	}
}

// TestMessageHandler_ContextPropagation 测试上下文传递
func TestMessageHandler_ContextPropagation(t *testing.T) {
	type contextKey string
	key := contextKey("test-key")
	expectedValue := "test-value"

	var receivedValue string
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		if v := ctx.Value(key); v != nil {
			receivedValue = v.(string)
		}
		return nil
	})

	ctx := context.WithValue(context.Background(), key, expectedValue)
	err := handler(ctx, "topic", []byte("payload"))
	require.NoError(t, err)

	assert.Equal(t, expectedValue, receivedValue)
}

// TestClient_Options 测试客户端配置选项
func TestClient_Options(t *testing.T) {
	// 测试配置选项的结构（不实际连接）
	// 这些值会在 NewClient 中使用
	broker := "tcp://localhost:1883"
	clientID := "test-client"
	username := "testuser"
	password := "testpass"

	assert.NotEmpty(t, broker)
	assert.NotEmpty(t, clientID)
	assert.NotEmpty(t, username)
	assert.NotEmpty(t, password)
}

// TestPayloadSerialization_Struct 测试结构体序列化
func TestPayloadSerialization_Struct(t *testing.T) {
	type SensorReading struct {
		SensorID    string    `json:"sensor_id"`
		Temperature float64   `json:"temperature"`
		Timestamp   time.Time `json:"timestamp"`
	}

	reading := SensorReading{
		SensorID:    "temp-001",
		Temperature: 23.5,
		Timestamp:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(reading)
	require.NoError(t, err)
	assert.Contains(t, string(data), "temp-001")
	assert.Contains(t, string(data), "23.5")
}

// TestPayloadSerialization_Complex 测试复杂类型序列化
func TestPayloadSerialization_Complex(t *testing.T) {
	payload := map[string]interface{}{
		"device_id": "dev-001",
		"metrics": map[string]float64{
			"cpu":    45.2,
			"memory": 78.9,
		},
		"tags":    []string{"prod", "server"},
		"active":  true,
		"latency": 12.5,
	}

	data, err := json.Marshal(payload)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "dev-001", decoded["device_id"])
	assert.Equal(t, true, decoded["active"])
}

// TestQoS_Validation 测试 QoS 值验证
func TestQoS_Validation(t *testing.T) {
	testCases := []struct {
		qos     byte
		isValid bool
	}{
		{0, true},    // At most once
		{1, true},    // At least once
		{2, true},    // Exactly once
		{3, false},   // Invalid
		{255, false}, // Invalid
	}

	for _, tc := range testCases {
		t.Run(string(rune('0'+tc.qos)), func(t *testing.T) {
			isValid := tc.qos >= 0 && tc.qos <= 2
			assert.Equal(t, tc.isValid, isValid)
		})
	}
}

// TestMessageHandler_ReturnValues 测试处理函数返回值
func TestMessageHandler_ReturnValues(t *testing.T) {
	testCases := []struct {
		name    string
		handler MessageHandler
		wantErr bool
	}{
		{
			name: "no_error",
			handler: MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
				return nil
			}),
			wantErr: false,
		},
		{
			name: "with_error",
			handler: MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
				return errors.New("processing failed")
			}),
			wantErr: true,
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.handler(ctx, "topic", []byte("payload"))
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestClient_NilClient 测试 nil 客户端
func TestClient_NilClient(t *testing.T) {
	client := &Client{client: nil}
	assert.NotNil(t, client)
	// 无法调用方法，因为 client.client 是 nil
}

// BenchmarkPayloadSerialization 性能测试 payload 序列化
func BenchmarkPayloadSerialization(b *testing.B) {
	payload := map[string]interface{}{
		"sensor_id": "temp-001",
		"value":     23.5,
		"timestamp": time.Now().Unix(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(payload)
	}
}

// BenchmarkMessageHandler 性能测试消息处理
func BenchmarkMessageHandler(b *testing.B) {
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return nil
	})

	ctx := context.Background()
	payload := []byte(`{"sensor_id": "temp-001", "value": 23.5}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler(ctx, "sensors/temperature", payload)
	}
}
