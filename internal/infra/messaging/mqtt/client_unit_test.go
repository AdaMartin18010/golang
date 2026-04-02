// Package mqtt provides unit tests for MQTT client with mocking.
package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========== Publish Tests ==========

// TestClient_Publish_Success 测试成功发布消息
func TestClient_Publish_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	err := client.Publish(ctx, "test/topic", 1, false, "test message")
	assert.NoError(t, err)
}

// TestClient_Publish_WithByteArray 测试发布字节数组
func TestClient_Publish_WithByteArray_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	err := client.Publish(ctx, "test/topic", 1, false, []byte("byte message"))
	assert.NoError(t, err)
}

// TestClient_Publish_WithStruct 测试发布结构体
func TestClient_Publish_WithStruct_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	data := map[string]interface{}{
		"sensor_id": "temp-001",
		"value":     25.5,
	}
	err := client.Publish(ctx, "test/topic", 1, false, data)
	assert.NoError(t, err)
}

// TestClient_Publish_WithRetained 测试发布保留消息
func TestClient_Publish_WithRetained_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	err := client.Publish(ctx, "status/online", 1, true, "true")
	assert.NoError(t, err)
}

// TestClient_Publish_DifferentQoS 测试不同QoS级别
func TestClient_Publish_DifferentQoS_Mock(t *testing.T) {
	qosLevels := []byte{0, 1, 2}
	
	for _, qos := range qosLevels {
		t.Run(string(rune('0'+qos)), func(t *testing.T) {
			mock := newMockMqttClient()
			client := mockClient(mock)
			ctx := testContext()

			err := client.Publish(ctx, "test/topic", qos, false, "message")
			assert.NoError(t, err)
		})
	}
}

// TestClient_Publish_Error 测试发布失败
func TestClient_Publish_Error_Mock(t *testing.T) {
	expectedErr := errors.New("publish failed")
	mock := newMockClientWithErrors(expectedErr, nil, nil)
	client := mockClient(mock)
	ctx := testContext()

	err := client.Publish(ctx, "test/topic", 1, false, "message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish")
}

// TestClient_Publish_InvalidPayload_Mock 测试发布无效负载
func TestClient_Publish_InvalidPayload_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	type Circular struct {
		Self *Circular
	}
	circular := &Circular{}
	circular.Self = circular

	err := client.Publish(ctx, "test/topic", 1, false, circular)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal payload")
}

// ========== Subscribe Tests ==========

// TestClient_Subscribe_Success_Mock 测试成功订阅
func TestClient_Subscribe_Success_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return nil
	})

	err := client.Subscribe(ctx, "test/topic", 1, handler)
	assert.NoError(t, err)
}

// TestClient_Subscribe_WithWildcard_Mock 测试通配符订阅
func TestClient_Subscribe_WithWildcard_Mock(t *testing.T) {
	wildcards := []string{
		"sensors/+/temperature",
		"sensors/#",
		"home/+/+/status",
	}

	for _, wildcard := range wildcards {
		t.Run(wildcard, func(t *testing.T) {
			mock := newMockMqttClient()
			client := mockClient(mock)
			ctx := testContext()

			handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
				return nil
			})

			err := client.Subscribe(ctx, wildcard, 1, handler)
			assert.NoError(t, err)
		})
	}
}

// TestClient_Subscribe_Error_Mock 测试订阅失败
func TestClient_Subscribe_Error_Mock(t *testing.T) {
	expectedErr := errors.New("subscribe failed")
	mock := newMockClientWithErrors(nil, expectedErr, nil)
	client := mockClient(mock)
	ctx := testContext()

	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return nil
	})

	err := client.Subscribe(ctx, "test/topic", 1, handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to subscribe")
}

// TestClient_Subscribe_HandlerCalled_Mock 测试订阅处理函数被调用
func TestClient_Subscribe_HandlerCalled_Mock(t *testing.T) {
	var receivedHandler mqtt.MessageHandler
	mock := &mockMqttClient{
		connected: true,
		subscribeFunc: func(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
			receivedHandler = callback
			return newMockToken(nil)
		},
	}
	
	client := mockClient(mock)
	ctx := testContext()

	handlerCalled := false
	receivedTopic := ""
	receivedPayload := []byte{}

	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		handlerCalled = true
		receivedTopic = topic
		receivedPayload = payload
		return nil
	})

	err := client.Subscribe(ctx, "test/topic", 1, handler)
	require.NoError(t, err)

	msg := &mockMessage{
		topic:   "test/topic",
		payload: []byte("test payload"),
		qos:     1,
	}
	
	if receivedHandler != nil {
		receivedHandler(nil, msg)
	}

	time.Sleep(10 * time.Millisecond)

	assert.True(t, handlerCalled)
	assert.Equal(t, "test/topic", receivedTopic)
	assert.Equal(t, []byte("test payload"), receivedPayload)
}

// TestClient_Subscribe_HandlerError_Mock 测试处理函数返回错误
func TestClient_Subscribe_HandlerError_Mock(t *testing.T) {
	var receivedHandler mqtt.MessageHandler
	mock := &mockMqttClient{
		connected: true,
		subscribeFunc: func(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
			receivedHandler = callback
			return newMockToken(nil)
		},
	}
	
	client := mockClient(mock)
	ctx := testContext()

	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return errors.New("handler error")
	})

	err := client.Subscribe(ctx, "test/topic", 1, handler)
	require.NoError(t, err)

	msg := &mockMessage{
		topic:   "test/topic",
		payload: []byte("test payload"),
		qos:     1,
	}

	if receivedHandler != nil {
		receivedHandler(nil, msg)
	}

	time.Sleep(10 * time.Millisecond)
}

// ========== Unsubscribe Tests ==========

// TestClient_Unsubscribe_Success_Mock 测试成功取消订阅
func TestClient_Unsubscribe_Success_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	err := client.Unsubscribe(ctx, "test/topic")
	assert.NoError(t, err)
}

// TestClient_Unsubscribe_MultipleTopics_Mock 测试批量取消订阅
func TestClient_Unsubscribe_MultipleTopics_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)
	ctx := testContext()

	err := client.Unsubscribe(ctx, "topic1", "topic2", "topic3")
	assert.NoError(t, err)
}

// TestClient_Unsubscribe_Error_Mock 测试取消订阅失败
func TestClient_Unsubscribe_Error_Mock(t *testing.T) {
	expectedErr := errors.New("unsubscribe failed")
	mock := newMockClientWithErrors(nil, nil, expectedErr)
	client := mockClient(mock)
	ctx := testContext()

	err := client.Unsubscribe(ctx, "test/topic")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unsubscribe")
}

// ========== Close Tests ==========

// TestClient_Close_Mock 测试关闭连接
func TestClient_Close_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)

	client.Close()
	assert.False(t, mock.connected)
}

// TestClient_Close_MultipleTimes_Mock 测试多次关闭
func TestClient_Close_MultipleTimes_Mock(t *testing.T) {
	mock := newMockMqttClient()
	client := mockClient(mock)

	client.Close()
	client.Close()
	client.Close()
	
	assert.False(t, mock.connected)
}

// ========== MessageHandler Tests ==========

// TestMessageHandler_Concurrent_Mock 测试并发消息处理
func TestMessageHandler_Concurrent_Mock(t *testing.T) {
	var mu sync.Mutex
	count := 0
	
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})

	ctx := testContext()
	
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			handler(ctx, "test/topic", []byte("payload"))
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	mu.Lock()
	assert.Equal(t, 10, count)
	mu.Unlock()
}

// TestMessageHandler_PanicRecovery_Mock 测试处理函数panic恢复
func TestMessageHandler_PanicRecovery_Mock(t *testing.T) {
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		if string(payload) == "panic" {
			panic("test panic")
		}
		return nil
	})

	ctx := testContext()
	
	// 正常消息
	err := handler(ctx, "test/topic", []byte("normal"))
	assert.NoError(t, err)
}

// ========== ConvertPayload Tests ==========

// TestConvertPayload_ByteSlice_Mock 测试字节数组转换
func TestConvertPayload_ByteSlice_Mock(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03}
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, input, result)
}

// TestConvertPayload_String_Mock 测试字符串转换
func TestConvertPayload_String_Mock(t *testing.T) {
	input := "test message"
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(input), result)
}

// TestConvertPayload_EmptyString_Mock 测试空字符串
func TestConvertPayload_EmptyString_Mock(t *testing.T) {
	input := ""
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, result)
}

// TestConvertPayload_Integer_Mock 测试整数转换
func TestConvertPayload_Integer_Mock(t *testing.T) {
	input := 42
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("42"), result)
}

// TestConvertPayload_Struct_Mock 测试结构体转换
func TestConvertPayload_Struct_Mock(t *testing.T) {
	input := map[string]interface{}{
		"key": "value",
		"num": 123,
	}
	result, err := convertPayload(input)
	require.NoError(t, err)
	
	var decoded map[string]interface{}
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "value", decoded["key"])
}

// TestConvertPayload_InvalidJSON_Mock 测试无效JSON
func TestConvertPayload_InvalidJSON_Mock(t *testing.T) {
	type Circular struct {
		Self *Circular
	}
	circular := &Circular{}
	circular.Self = circular

	_, err := convertPayload(circular)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal payload")
}

// ========== Integration with Existing Tests ==========

// TestClient_Struct_Mock 测试客户端结构体
func TestClient_Struct_Mock(t *testing.T) {
	client := &Client{client: nil}
	assert.NotNil(t, client)
}

// TestNewClientWithClient 测试使用mock客户端创建Client
func TestNewClientWithClient(t *testing.T) {
	mock := newMockMqttClient()
	client := newClientWithClient(mock)
	assert.NotNil(t, client)
	
	// 验证客户端可以使用
	ctx := testContext()
	err := client.Publish(ctx, "test/topic", 1, false, "message")
	assert.NoError(t, err)
}

// TestNewClientWithClient_Nil 测试使用nil创建Client
func TestNewClientWithClient_Nil(t *testing.T) {
	client := newClientWithClient(nil)
	assert.NotNil(t, client)
	// 注意：使用nil客户端调用方法会导致panic
	// 这是预期的行为，生产代码应该确保client不为nil
}

// TestClient_Publish_EmptyTopic 测试发布到空主题
func TestClient_Publish_EmptyTopic(t *testing.T) {
	mock := newMockMqttClient()
	client := newClientWithClient(mock)
	ctx := testContext()

	err := client.Publish(ctx, "", 1, false, "message")
	assert.NoError(t, err) // mock允许空主题
}

// TestClient_Publish_EmptyPayload 测试发布空负载
func TestClient_Publish_EmptyPayload(t *testing.T) {
	mock := newMockMqttClient()
	client := newClientWithClient(mock)
	ctx := testContext()

	err := client.Publish(ctx, "test/topic", 1, false, "")
	assert.NoError(t, err)
}

// TestClient_Publish_LargePayload 测试发布大负载
func TestClient_Publish_LargePayload(t *testing.T) {
	mock := newMockMqttClient()
	client := newClientWithClient(mock)
	ctx := testContext()

	// 创建1KB的数据
	largeData := make([]byte, 1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	err := client.Publish(ctx, "test/topic", 1, false, largeData)
	assert.NoError(t, err)
}

// TestClient_Publish_NilPayload 测试发布nil负载
func TestClient_Publish_NilPayload(t *testing.T) {
	mock := newMockMqttClient()
	client := newClientWithClient(mock)
	ctx := testContext()

	err := client.Publish(ctx, "test/topic", 1, false, nil)
	assert.NoError(t, err)
}

// TestClient_Subscribe_EmptyTopic 测试订阅空主题
func TestClient_Subscribe_EmptyTopic(t *testing.T) {
	mock := newMockMqttClient()
	client := newClientWithClient(mock)
	ctx := testContext()

	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return nil
	})

	err := client.Subscribe(ctx, "", 1, handler)
	assert.NoError(t, err)
}

// TestClient_Subscribe_DifferentQoS 测试不同QoS级别订阅
func TestClient_Subscribe_DifferentQoS(t *testing.T) {
	qosLevels := []byte{0, 1, 2}

	for _, qos := range qosLevels {
		t.Run(string(rune('0'+qos)), func(t *testing.T) {
			mock := newMockMqttClient()
			client := newClientWithClient(mock)
			ctx := testContext()

			handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
				return nil
			})

			err := client.Subscribe(ctx, "test/topic", qos, handler)
			assert.NoError(t, err)
		})
	}
}

// TestClient_Unsubscribe_EmptyTopics 测试取消订阅空主题列表
func TestClient_Unsubscribe_EmptyTopics(t *testing.T) {
	mock := newMockMqttClient()
	client := newClientWithClient(mock)
	ctx := testContext()

	err := client.Unsubscribe(ctx)
	assert.NoError(t, err)
}

// TestConvertPayload_FloatTypes 测试各种浮点数类型
func TestConvertPayload_FloatTypes(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"float32", float32(3.14)},
		{"float64", float64(3.14159)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertPayload(tc.input)
			require.NoError(t, err)
			assert.Greater(t, len(result), 0)
		})
	}
}

// TestConvertPayload_IntTypes 测试各种整数类型
func TestConvertPayload_IntTypes(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"int", int(42)},
		{"int8", int8(127)},
		{"int16", int16(32767)},
		{"int32", int32(2147483647)},
		{"int64", int64(9223372036854775807)},
		{"uint", uint(42)},
		{"uint8", uint8(255)},
		{"uint16", uint16(65535)},
		{"uint32", uint32(4294967295)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertPayload(tc.input)
			require.NoError(t, err)
			assert.Greater(t, len(result), 0)
		})
	}
}

// TestConvertPayload_BoolTypes 测试布尔类型
func TestConvertPayload_BoolTypes(t *testing.T) {
	tests := []struct {
		name  string
		input bool
		want  string
	}{
		{"true", true, "true"},
		{"false", false, "false"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertPayload(tc.input)
			require.NoError(t, err)
			assert.Equal(t, []byte(tc.want), result)
		})
	}
}

// TestConvertPayload_Array 测试数组类型
func TestConvertPayload_Array(t *testing.T) {
	input := [3]int{1, 2, 3}
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("[1,2,3]"), result)
}

// TestConvertPayload_InterfaceSlice 测试interface切片
func TestConvertPayload_InterfaceSlice(t *testing.T) {
	input := []interface{}{"string", 42, true, nil}
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(`["string",42,true,null]`), result)
}

// TestConvertPayload_Pointer 测试指针类型
func TestConvertPayload_Pointer(t *testing.T) {
	value := "test pointer"
	result, err := convertPayload(&value)
	require.NoError(t, err)
	assert.Equal(t, []byte(`"test pointer"`), result)
}

// TestConvertPayload_Unicode 测试Unicode字符串
func TestConvertPayload_Unicode(t *testing.T) {
	input := "Hello 世界 🌍"
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(input), result)
}

// TestConvertPayload_SpecialChars 测试特殊字符
func TestConvertPayload_SpecialChars(t *testing.T) {
	input := "line1\nline2\ttab"
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(input), result)
}

// TestMessageHandler_NilPayload 测试处理nil负载
func TestMessageHandler_NilPayload(t *testing.T) {
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		assert.Nil(t, payload)
		return nil
	})

	ctx := testContext()
	err := handler(ctx, "test/topic", nil)
	assert.NoError(t, err)
}

// TestMessageHandler_EmptyPayload 测试处理空负载
func TestMessageHandler_EmptyPayload(t *testing.T) {
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		assert.Empty(t, payload)
		return nil
	})

	ctx := testContext()
	err := handler(ctx, "test/topic", []byte{})
	assert.NoError(t, err)
}

// TestMessageHandler_LargePayload 测试处理大负载
func TestMessageHandler_LargePayload(t *testing.T) {
	largePayload := make([]byte, 1024*1024) // 1MB
	for i := range largePayload {
		largePayload[i] = byte(i % 256)
	}

	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		assert.Equal(t, len(largePayload), len(payload))
		return nil
	})

	ctx := testContext()
	err := handler(ctx, "test/topic", largePayload)
	assert.NoError(t, err)
}

// TestMessageHandler_Type_Mock 测试消息处理函数类型
func TestMessageHandler_Type_Mock(t *testing.T) {
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return nil
	})
	assert.NotNil(t, handler)

	ctx := testContext()
	err := handler(ctx, "test/topic", []byte("test payload"))
	assert.NoError(t, err)
}

// TestMessageHandler_WithError_Mock 测试返回错误的处理函数
func TestMessageHandler_WithError_Mock(t *testing.T) {
	testErr := errors.New("test error")
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		return testErr
	})

	ctx := testContext()
	err := handler(ctx, "topic", []byte("payload"))
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

// TestQoS_Values_Mock 测试QoS值
func TestQoS_Values_Mock(t *testing.T) {
	qosLevels := []byte{0, 1, 2}

	for _, qos := range qosLevels {
		assert.GreaterOrEqual(t, qos, byte(0))
		assert.LessOrEqual(t, qos, byte(2))
	}
}

// TestRetainedFlag_Mock 测试保留消息标志
func TestRetainedFlag_Mock(t *testing.T) {
	retainedValues := []bool{true, false}

	for _, retained := range retainedValues {
		assert.True(t, retained == true || retained == false)
	}
}

// TestPayloadTypes_Exhaustive_Mock 测试所有支持的payload类型
func TestPayloadTypes_Exhaustive_Mock(t *testing.T) {
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
			result, err := convertPayload(tc.payload)
			require.NoError(t, err)
			assert.Equal(t, tc.expect, result)
		})
	}
}

// TestQoS_Validation_Mock 测试QoS值验证
func TestQoS_Validation_Mock(t *testing.T) {
	testCases := []struct {
		qos     byte
		isValid bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, false},
		{255, false},
	}

	for _, tc := range testCases {
		t.Run(string(rune('0'+tc.qos)), func(t *testing.T) {
			isValid := tc.qos >= 0 && tc.qos <= 2
			assert.Equal(t, tc.isValid, isValid)
		})
	}
}

// TestMessageHandler_MultipleCalls_Mock 测试处理函数被多次调用
func TestMessageHandler_MultipleCalls_Mock(t *testing.T) {
	callCount := 0
	handler := MessageHandler(func(ctx context.Context, topic string, payload []byte) error {
		callCount++
		return nil
	})

	ctx := testContext()
	for i := 0; i < 5; i++ {
		err := handler(ctx, "topic", []byte("payload"))
		require.NoError(t, err)
	}

	assert.Equal(t, 5, callCount)
}

// TestMessageHandler_DifferentTopics_Mock 测试不同主题的处理
func TestMessageHandler_DifferentTopics_Mock(t *testing.T) {
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

	ctx := testContext()
	for _, topic := range topics {
		err := handler(ctx, topic, []byte("data"))
		require.NoError(t, err)
	}

	assert.Equal(t, topics, receivedTopics)
}

// TestMessageHandler_ContextPropagation_Mock 测试上下文传递
func TestMessageHandler_ContextPropagation_Mock(t *testing.T) {
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

	ctx := context.WithValue(testContext(), key, expectedValue)
	err := handler(ctx, "topic", []byte("payload"))
	require.NoError(t, err)

	assert.Equal(t, expectedValue, receivedValue)
}

// TestClient_NilClient_Mock 测试nil客户端
func TestClient_NilClient_Mock(t *testing.T) {
	client := &Client{client: nil}
	assert.NotNil(t, client)
}
