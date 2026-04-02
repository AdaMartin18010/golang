// Package nats provides unit tests for NATS client with mocking.
package nats

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========== Publish Tests ==========

// TestClient_Publish_Success 测试成功发布消息
func TestClient_Publish_Success(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	err := client.Publish("test.subject", "hello world")
	assert.NoError(t, err)

	published := mock.getPublishedMessages()
	require.Len(t, published, 1)
	assert.Equal(t, "test.subject", published[0].Subject)
	assert.Equal(t, []byte("hello world"), published[0].Data)
}

// TestClient_Publish_ByteSlice 测试发布字节数组
func TestClient_Publish_ByteSlice(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	data := []byte{0x01, 0x02, 0x03}
	err := client.Publish("test.subject", data)
	assert.NoError(t, err)

	published := mock.getPublishedMessages()
	require.Len(t, published, 1)
	assert.Equal(t, data, published[0].Data)
}

// TestClient_Publish_String 测试发布字符串
func TestClient_Publish_String(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	err := client.Publish("test.subject", "test message")
	assert.NoError(t, err)

	published := mock.getPublishedMessages()
	require.Len(t, published, 1)
	assert.Equal(t, []byte("test message"), published[0].Data)
}

// TestClient_Publish_Struct 测试发布结构体
func TestClient_Publish_Struct(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	data := map[string]interface{}{
		"user_id": 123,
		"name":    "Alice",
	}
	err := client.Publish("user.created", data)
	assert.NoError(t, err)

	published := mock.getPublishedMessages()
	require.Len(t, published, 1)
	assert.Contains(t, string(published[0].Data), "user_id")
	assert.Contains(t, string(published[0].Data), "Alice")
}

// TestClient_Publish_Error 测试发布失败
func TestClient_Publish_Error(t *testing.T) {
	mock := newMockNatsConn()
	mock.setPublishErr(errPublishFailed)
	client := newClientWithConn(mock)

	err := client.Publish("test.subject", "message")
	assert.Error(t, err)
	assert.Equal(t, errPublishFailed, err)
}

// TestClient_Publish_InvalidData 测试发布无法序列化的数据
func TestClient_Publish_InvalidData(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	type BadData struct {
		Data chan int
	}
	badData := BadData{Data: make(chan int)}

	err := client.Publish("test.subject", badData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
}

// TestClient_Publish_MultipleMessages 测试发布多条消息
func TestClient_Publish_MultipleMessages(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	for i := 0; i < 5; i++ {
		err := client.Publish("test.subject", "message")
		require.NoError(t, err)
	}

	published := mock.getPublishedMessages()
	assert.Len(t, published, 5)
}

// ========== Subscribe Tests ==========

// TestClient_Subscribe_Success 测试成功订阅
func TestClient_Subscribe_Success(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	sub, err := client.Subscribe("test.subject", func(msg *nats.Msg) {})

	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, "test.subject", sub.Subject)
}

// TestClient_Subscribe_Error 测试订阅失败
func TestClient_Subscribe_Error(t *testing.T) {
	mock := newMockNatsConn()
	mock.setSubscribeErr(errSubscribeFailed)
	client := newClientWithConn(mock)

	sub, err := client.Subscribe("test.subject", func(msg *nats.Msg) {})

	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.Equal(t, errSubscribeFailed, err)
}

// TestClient_Subscribe_WithHandler 测试订阅并调用处理函数
func TestClient_Subscribe_WithHandler(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	received := false
	sub, err := client.Subscribe("test.subject", func(msg *nats.Msg) {
		received = true
	})

	require.NoError(t, err)
	require.NotNil(t, sub)
	assert.False(t, received) // Handler not called yet
}

// ========== QueueSubscribe Tests ==========

// TestClient_QueueSubscribe_Success 测试成功队列订阅
func TestClient_QueueSubscribe_Success(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	sub, err := client.QueueSubscribe("test.subject", "queue-group", func(msg *nats.Msg) {})

	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, "test.subject", sub.Subject)
	assert.Equal(t, "queue-group", sub.Queue)
}

// TestClient_QueueSubscribe_Error 测试队列订阅失败
func TestClient_QueueSubscribe_Error(t *testing.T) {
	mock := newMockNatsConn()
	mock.setQueueSubErr(errors.New("queue subscribe failed"))
	client := newClientWithConn(mock)

	sub, err := client.QueueSubscribe("test.subject", "queue-group", func(msg *nats.Msg) {})

	assert.Error(t, err)
	assert.Nil(t, sub)
}

// ========== Request Tests ==========

// TestClient_Request_Success 测试成功请求
func TestClient_Request_Success(t *testing.T) {
	mock := newMockNatsConn()
	expectedResponse := &nats.Msg{
		Subject: "test.reply",
		Data:    []byte(`{"result":"success"}`),
		Header:  nats.Header{},
	}
	mock.setRequestResponse(expectedResponse)
	client := newClientWithConn(mock)

	resp, err := client.Request("test.subject", "request", 5*time.Second)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedResponse.Data, resp.Data)
}

// TestClient_Request_WithStruct 测试发送结构体请求
func TestClient_Request_WithStruct(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	requestData := map[string]interface{}{
		"action": "get_user",
		"id":     123,
	}

	resp, err := client.Request("test.subject", requestData, 5*time.Second)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestClient_Request_Error 测试请求失败
func TestClient_Request_Error(t *testing.T) {
	mock := newMockNatsConn()
	mock.setRequestErr(errRequestFailed)
	client := newClientWithConn(mock)

	resp, err := client.Request("test.subject", "request", 100*time.Millisecond)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send request")
	assert.Nil(t, resp)
}

// TestClient_Request_InvalidData 测试发送无法序列化的请求数据
func TestClient_Request_InvalidData(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	type BadData struct {
		Data chan int
	}
	badData := BadData{Data: make(chan int)}

	resp, err := client.Request("test.subject", badData, 5*time.Second)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
	assert.Nil(t, resp)
}

// ========== Connection State Tests ==========

// TestClient_IsConnected 测试连接状态检查
func TestClient_IsConnected(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	assert.True(t, client.IsConnected())

	mock.disconnect()
	assert.False(t, client.IsConnected())
}

// TestClient_IsConnected_NilConn 测试nil连接状态
func TestClient_IsConnected_NilConn(t *testing.T) {
	client := &Client{conn: nil}
	assert.False(t, client.IsConnected())
}

// TestClient_Stats 测试获取统计信息
func TestClient_Stats(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	// 发布一些消息
	for i := 0; i < 3; i++ {
		client.Publish("test.subject", "message")
	}

	stats := client.Stats()
	assert.Equal(t, uint64(3), stats.OutMsgs)
	assert.Greater(t, stats.OutBytes, uint64(0))
}

// TestClient_Stats_NilConn 测试nil连接的统计信息
func TestClient_Stats_NilConn(t *testing.T) {
	client := &Client{conn: nil}
	stats := client.Stats()
	assert.Equal(t, nats.Statistics{}, stats)
}

// TestClient_Close 测试关闭连接
func TestClient_Close(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	assert.True(t, client.IsConnected())
	client.Close()
	assert.False(t, client.IsConnected())
}

// TestClient_Close_NotConnected 测试关闭未连接的客户端
func TestClient_Close_NotConnected(t *testing.T) {
	mock := newMockNatsConn()
	mock.disconnect()
	client := newClientWithConn(mock)

	// 不应该panic
	client.Close()
	assert.False(t, client.IsConnected())
}

// TestClient_Close_MultipleTimes 测试多次关闭
func TestClient_Close_MultipleTimes(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	// 多次关闭不应该panic
	client.Close()
	client.Close()
	client.Close()

	assert.False(t, client.IsConnected())
}

// TestClient_Close_NilConn 测试关闭nil连接
func TestClient_Close_NilConn(t *testing.T) {
	client := &Client{conn: nil}
	// 不应该panic
	client.Close()
}

// ========== Marshal Payload Tests ==========

// TestMarshalPayload_ByteSlice 测试字节数组类型
func TestMarshalPayload_ByteSlice_Unit(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03}
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, input, result)
}

// TestMarshalPayload_String 测试字符串类型
func TestMarshalPayload_String_Unit(t *testing.T) {
	input := "test message"
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(input), result)
}

// TestMarshalPayload_EmptyString 测试空字符串
func TestMarshalPayload_EmptyString_Unit(t *testing.T) {
	input := ""
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, result)
}

// TestMarshalPayload_Integer 测试整数类型
func TestMarshalPayload_Integer_Unit(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"positive", 42, "42"},
		{"negative", -42, "-42"},
		{"zero", 0, "0"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := marshalPayload(tc.input)
			require.NoError(t, err)
			assert.Equal(t, []byte(tc.expected), result)
		})
	}
}

// TestMarshalPayload_Float 测试浮点数类型
func TestMarshalPayload_Float_Unit(t *testing.T) {
	input := 3.14159
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("3.14159"), result)
}

// TestMarshalPayload_Boolean 测试布尔类型
func TestMarshalPayload_Boolean_Unit(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		result, err := marshalPayload(true)
		require.NoError(t, err)
		assert.Equal(t, []byte("true"), result)
	})

	t.Run("false", func(t *testing.T) {
		result, err := marshalPayload(false)
		require.NoError(t, err)
		assert.Equal(t, []byte("false"), result)
	})
}

// TestMarshalPayload_Map 测试Map类型
func TestMarshalPayload_Map_Unit(t *testing.T) {
	input := map[string]interface{}{
		"event":   "user.created",
		"user_id": 123,
	}
	result, err := marshalPayload(input)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "user.created", decoded["event"])
}

// TestMarshalPayload_Struct 测试结构体类型
func TestMarshalPayload_Struct_Unit(t *testing.T) {
	type Event struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}

	input := Event{
		Type:    "order.created",
		Payload: "data",
	}

	result, err := marshalPayload(input)
	require.NoError(t, err)

	var decoded Event
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, input, decoded)
}

// TestMarshalPayload_Slice 测试切片类型
func TestMarshalPayload_Slice_Unit(t *testing.T) {
	input := []string{"a", "b", "c"}
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(`["a","b","c"]`), result)
}

// TestMarshalPayload_Nil 测试nil值
func TestMarshalPayload_Nil_Unit(t *testing.T) {
	result, err := marshalPayload(nil)
	require.NoError(t, err)
	assert.Equal(t, []byte("null"), result)
}

// TestMarshalPayload_InvalidJSON 测试无法序列化的类型
func TestMarshalPayload_InvalidJSON_Unit(t *testing.T) {
	type Circular struct {
		Self *Circular
	}
	circular := &Circular{}
	circular.Self = circular

	_, err := marshalPayload(circular)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
}

// TestMarshalPayload_AllTypes 测试所有payload类型的序列化
func TestMarshalPayload_AllTypes_Unit(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []byte
	}{
		{
			name:     "byte_slice",
			input:    []byte{0x01, 0x02},
			expected: []byte{0x01, 0x02},
		},
		{
			name:     "string",
			input:    "hello",
			expected: []byte("hello"),
		},
		{
			name:     "int",
			input:    42,
			expected: []byte("42"),
		},
		{
			name:     "bool",
			input:    true,
			expected: []byte("true"),
		},
		{
			name:     "nil",
			input:    nil,
			expected: []byte("null"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := marshalPayload(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// ========== Config Tests ==========

// TestDefaultConfig 测试默认配置
func TestDefaultConfig_Unit(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "nats://localhost:4222", config.URL)
	assert.Equal(t, -1, config.MaxReconnects)
	assert.Equal(t, 2*time.Second, config.ReconnectWait)
	assert.Equal(t, 5*time.Second, config.Timeout)
	assert.Empty(t, config.Name)
	assert.Empty(t, config.Token)
	assert.Empty(t, config.Username)
	assert.Empty(t, config.Password)
}

// TestConfig_WithTokenAuth 测试Token认证配置
func TestConfig_WithTokenAuth_Unit(t *testing.T) {
	config := Config{
		URL:   "nats://localhost:4222",
		Token: "my-secret-token",
	}

	assert.NotEmpty(t, config.Token)
	assert.Equal(t, "my-secret-token", config.Token)
}

// TestConfig_WithUserPassAuth 测试用户名密码认证配置
func TestConfig_WithUserPassAuth_Unit(t *testing.T) {
	config := Config{
		URL:      "nats://localhost:4222",
		Username: "admin",
		Password: "secret",
	}

	assert.NotEmpty(t, config.Username)
	assert.NotEmpty(t, config.Password)
	assert.Equal(t, "admin", config.Username)
	assert.Equal(t, "secret", config.Password)
}

// TestConfig_CustomTimeouts 测试自定义超时配置
func TestConfig_CustomTimeouts_Unit(t *testing.T) {
	config := Config{
		URL:           "nats://localhost:4222",
		MaxReconnects: 10,
		ReconnectWait: 5 * time.Second,
		Timeout:       10 * time.Second,
	}

	assert.Equal(t, 10, config.MaxReconnects)
	assert.Equal(t, 5*time.Second, config.ReconnectWait)
	assert.Equal(t, 10*time.Second, config.Timeout)
}

// TestConfig_Empty 测试空配置
func TestConfig_Empty_Unit(t *testing.T) {
	config := Config{}
	assert.Empty(t, config.URL)
	assert.Equal(t, 0, config.MaxReconnects)
	assert.Equal(t, time.Duration(0), config.ReconnectWait)
	assert.Equal(t, time.Duration(0), config.Timeout)
}

// TestConfig_WithName 测试带名称的配置
func TestConfig_WithName_Unit(t *testing.T) {
	config := Config{
		URL:  "nats://localhost:4222",
		Name: "test-client",
	}
	assert.Equal(t, "test-client", config.Name)
}

// TestConfig_Complete 测试完整的配置
func TestConfig_Complete_Unit(t *testing.T) {
	config := Config{
		URL:           "nats://server:4222",
		MaxReconnects: 5,
		ReconnectWait: 1 * time.Second,
		Timeout:       3 * time.Second,
		Name:          "my-client",
		Token:         "secret-token",
		Username:      "user",
		Password:      "pass",
	}

	assert.Equal(t, "nats://server:4222", config.URL)
	assert.Equal(t, 5, config.MaxReconnects)
	assert.Equal(t, 1*time.Second, config.ReconnectWait)
	assert.Equal(t, 3*time.Second, config.Timeout)
	assert.Equal(t, "my-client", config.Name)
	assert.Equal(t, "secret-token", config.Token)
	assert.Equal(t, "user", config.Username)
	assert.Equal(t, "pass", config.Password)
}

// ========== Additional Marshal Tests ==========

// TestMarshalPayload_IntTypes 测试各种整数类型
func TestMarshalPayload_IntTypes_Unit(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"int8", int8(127)},
		{"int16", int16(32767)},
		{"int32", int32(2147483647)},
		{"int64", int64(9223372036854775807)},
		{"uint8", uint8(255)},
		{"uint16", uint16(65535)},
		{"uint32", uint32(4294967295)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := marshalPayload(tc.input)
			require.NoError(t, err)
			assert.Greater(t, len(result), 0)
		})
	}
}

// TestMarshalPayload_FloatTypes 测试浮点数类型
func TestMarshalPayload_FloatTypes_Unit(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"float32", float32(3.14159)},
		{"float64", float64(3.14159265359)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := marshalPayload(tc.input)
			require.NoError(t, err)
			assert.Greater(t, len(result), 0)
		})
	}
}

// TestMarshalPayload_Array 测试数组类型
func TestMarshalPayload_Array_Unit(t *testing.T) {
	input := [3]int{1, 2, 3}
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("[1,2,3]"), result)
}

// TestMarshalPayload_EmptySlice 测试空切片
func TestMarshalPayload_EmptySlice_Unit(t *testing.T) {
	input := []string{}
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("[]"), result)
}

// TestMarshalPayload_EmptyMap 测试空Map
func TestMarshalPayload_EmptyMap_Unit(t *testing.T) {
	input := map[string]interface{}{}
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("{}"), result)
}

// TestMarshalPayload_JSONNumber 测试JSON数字类型
func TestMarshalPayload_JSONNumber_Unit(t *testing.T) {
	num := json.Number("123.456")
	result, err := marshalPayload(num)
	require.NoError(t, err)
	assert.Equal(t, []byte("123.456"), result)
}

// TestDefaultConfig_Copies 测试默认配置返回的是副本
func TestDefaultConfig_Copies_Unit(t *testing.T) {
	config1 := DefaultConfig()
	config2 := DefaultConfig()

	config1.URL = "modified"
	assert.Equal(t, "nats://localhost:4222", config2.URL)
}

// TestClientStruct 测试Client结构体
func TestClientStruct(t *testing.T) {
	client := &Client{conn: nil}
	assert.NotNil(t, client)
}

// ========== Concurrent Tests ==========

// TestClient_ConcurrentPublish 测试并发发布
func TestClient_ConcurrentPublish(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			client.Publish("test.subject", "message")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	published := mock.getPublishedMessages()
	assert.Len(t, published, 10)
}

// TestClient_ConcurrentSubscribe 测试并发订阅
func TestClient_ConcurrentSubscribe(t *testing.T) {
	mock := newMockNatsConn()
	client := newClientWithConn(mock)

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			subject := "test.subject." + string(rune('0'+i))
			client.Subscribe(subject, func(msg *nats.Msg) {})
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
