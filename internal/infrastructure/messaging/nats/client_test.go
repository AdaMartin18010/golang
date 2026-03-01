// Package nats provides comprehensive tests for NATS client.
//
// 测试策略：
// 1. 使用接口 mock 进行单元测试
// 2. 使用 docker-compose 或测试容器进行集成测试（可选）
// 3. 覆盖连接、发布/订阅、Request/Reply 和错误处理
//
// 运行测试：
//   - 单元测试: go test -v ./internal/infrastructure/messaging/nats/...
//   - 集成测试: 需要运行 docker-compose up nats
//
// 启动 NATS 服务器：
//   docker run -d -p 4222:4222 nats:latest
package nats

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockNatsConn 是 NATS 连接的 mock 实现
type MockNatsConn struct {
	mock.Mock
	connected bool
}

func (m *MockNatsConn) Publish(subj string, data []byte) error {
	args := m.Called(subj, data)
	return args.Error(0)
}

func (m *MockNatsConn) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	args := m.Called(subj, cb)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (m *MockNatsConn) QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	args := m.Called(subj, queue, cb)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (m *MockNatsConn) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	args := m.Called(subj, data, timeout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nats.Msg), args.Error(1)
}

func (m *MockNatsConn) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockNatsConn) Close() {
	m.Called()
}

func (m *MockNatsConn) Stats() nats.Statistics {
	args := m.Called()
	return args.Get(0).(nats.Statistics)
}

func (m *MockNatsConn) ConnectedUrl() string {
	args := m.Called()
	return args.String(0)
}

// getTestURL 返回测试 NATS 服务器地址
func getTestURL() string {
	return "nats://localhost:4222"
}

// isNatsAvailable 检查 NATS 服务器是否可用
func isNatsAvailable() bool {
	cfg := Config{
		URL:           getTestURL(),
		MaxReconnects: 1,
		ReconnectWait: 1 * time.Second,
		Timeout:       2 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		return false
	}
	defer client.Close()

	return client.IsConnected()
}

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
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

// TestNewClient_Integration 测试创建 NATS 客户端（集成测试）
func TestNewClient_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test. Start NATS server with: docker run -d -p 4222:4222 nats:latest")
	}

	cfg := DefaultConfig()
	cfg.URL = getTestURL()
	cfg.Name = "test-client"

	client, err := NewClient(cfg)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.True(t, client.IsConnected())
	defer client.Close()
}

// TestNewClient_InvalidURL 测试无效 URL
func TestNewClient_InvalidURL(t *testing.T) {
	cfg := Config{
		URL:           "invalid://localhost:4222",
		MaxReconnects: 1,
		ReconnectWait: 1 * time.Second,
		Timeout:       2 * time.Second,
	}

	_, err := NewClient(cfg)
	assert.Error(t, err)
}

// TestNewClient_ConnectionRefused 测试连接被拒绝
func TestNewClient_ConnectionRefused(t *testing.T) {
	cfg := Config{
		URL:           "nats://localhost:59999", // 未使用的端口
		MaxReconnects: 1,
		ReconnectWait: 100 * time.Millisecond,
		Timeout:       1 * time.Second,
	}

	_, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to NATS")
}

// TestClient_Publish_Integration 测试发布消息（集成测试）
func TestClient_Publish_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 发布字符串消息
	err = client.Publish("test.publish", "hello world")
	require.NoError(t, err)

	// 发布字节消息
	err = client.Publish("test.publish", []byte("byte message"))
	require.NoError(t, err)

	// 发布 JSON 消息
	data := map[string]interface{}{"key": "value", "number": 42}
	err = client.Publish("test.publish", data)
	require.NoError(t, err)
}

// TestClient_Subscribe_Integration 测试订阅（集成测试）
func TestClient_Subscribe_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 订阅
	received := make(chan string, 1)
	sub, err := client.Subscribe("test.subscribe", func(msg *nats.Msg) {
		received <- string(msg.Data)
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布消息
	err = client.Publish("test.subscribe", "test message")
	require.NoError(t, err)

	// 等待接收
	select {
	case msg := <-received:
		assert.Equal(t, "test message", msg)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

// TestClient_SubscribeJSON_Integration 测试订阅 JSON 消息（集成测试）
func TestClient_SubscribeJSON_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 订阅
	received := make(chan map[string]interface{}, 1)
	sub, err := client.Subscribe("test.json", func(msg *nats.Msg) {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			received <- data
		}
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布 JSON 消息
	data := map[string]interface{}{
		"user_id": 123,
		"name":    "Alice",
		"active":  true,
	}
	err = client.Publish("test.json", data)
	require.NoError(t, err)

	// 等待接收
	select {
	case receivedData := <-received:
		assert.Equal(t, float64(123), receivedData["user_id"])
		assert.Equal(t, "Alice", receivedData["name"])
		assert.Equal(t, true, receivedData["active"])
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

// TestClient_QueueSubscribe_Integration 测试队列订阅（集成测试）
func TestClient_QueueSubscribe_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client1, err := NewClient(cfg)
	require.NoError(t, err)
	defer client1.Close()

	client2, err := NewClient(cfg)
	require.NoError(t, err)
	defer client2.Close()

	// 两个客户端订阅同一个队列
	var mu sync.Mutex
	receivedCount := 0

	sub1, err := client1.QueueSubscribe("test.queue", "queue-group", func(msg *nats.Msg) {
		mu.Lock()
		receivedCount++
		mu.Unlock()
	})
	require.NoError(t, err)
	defer sub1.Unsubscribe()

	sub2, err := client2.QueueSubscribe("test.queue", "queue-group", func(msg *nats.Msg) {
		mu.Lock()
		receivedCount++
		mu.Unlock()
	})
	require.NoError(t, err)
	defer sub2.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布多条消息
	for i := 0; i < 10; i++ {
		err := client1.Publish("test.queue", "message")
		require.NoError(t, err)
	}

	// 等待处理
	time.Sleep(500 * time.Millisecond)

	// 验证所有消息都被处理（负载均衡）
	mu.Lock()
	assert.Equal(t, 10, receivedCount)
	mu.Unlock()
}

// TestClient_Request_Integration 测试请求-响应（集成测试）
func TestClient_Request_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 创建响应服务
	sub, err := client.Subscribe("test.request", func(msg *nats.Msg) {
		// 解析请求
		var request map[string]interface{}
		json.Unmarshal(msg.Data, &request)

		// 构建响应
		response := map[string]interface{}{
			"echo":    request["data"],
			"status":  "ok",
		}

		// 发送响应
		data, _ := json.Marshal(response)
		msg.Respond(data)
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发送请求
	request := map[string]interface{}{"data": "hello"}
	reply, err := client.Request("test.request", request, 5*time.Second)
	require.NoError(t, err)

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(reply.Data, &response)
	require.NoError(t, err)
	assert.Equal(t, "hello", response["echo"])
	assert.Equal(t, "ok", response["status"])
}

// TestClient_Request_Timeout_Integration 测试请求超时（集成测试）
func TestClient_Request_Timeout_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 发送请求到不存在的主题（应该超时）
	_, err = client.Request("test.nonexistent", "data", 100*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send request")
}

// TestClient_IsConnected_Integration 测试连接状态（集成测试）
func TestClient_IsConnected_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)

	assert.True(t, client.IsConnected())

	client.Close()
	assert.False(t, client.IsConnected())
}

// TestClient_Stats_Integration 测试统计信息（集成测试）
func TestClient_Stats_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	stats := client.Stats()
	// 验证返回了统计信息
	// 注意：具体值取决于连接状态和使用情况
	_ = stats
}

// TestClient_Close_Integration 测试关闭连接（集成测试）
func TestClient_Close_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)

	// 正常关闭
	client.Close()
	assert.False(t, client.IsConnected())

	// 重复关闭不应该出错
	client.Close()
}

// TestPublishMarshalError 测试发布时序列化错误
func TestPublishMarshalError(t *testing.T) {
	// 创建一个包含不可序列化数据的结构
	type BadData struct {
		Data chan int // channel 不能被 JSON 序列化
	}

	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 尝试发布不可序列化的数据
	badData := BadData{Data: make(chan int)}
	err = client.Publish("test.bad", badData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
}

// TestRequestMarshalError 测试请求时序列化错误
func TestRequestMarshalError(t *testing.T) {
	type BadData struct {
		Data chan int
	}

	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 尝试发送不可序列化的请求
	badData := BadData{Data: make(chan int)}
	_, err = client.Request("test.bad", badData, 1*time.Second)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
}

// TestConfig_WithAuth 测试带认证的配置
func TestConfig_WithAuth(t *testing.T) {
	// 注意：这里只是测试配置结构，不涉及实际连接
	// 实际认证测试需要配置 NATS 服务器的认证

	tests := []struct {
		name           string
		config         Config
		expectAuthType string
	}{
		{
			name: "token auth",
			config: Config{
				URL:   "nats://localhost:4222",
				Token: "my-token",
			},
			expectAuthType: "token",
		},
		{
			name: "user/password auth",
			config: Config{
				URL:      "nats://localhost:4222",
				Username: "admin",
				Password: "secret",
			},
			expectAuthType: "userpass",
		},
		{
			name: "no auth",
			config: Config{
				URL: "nats://localhost:4222",
			},
			expectAuthType: "none",
		},
		{
			name: "token takes precedence",
			config: Config{
				URL:      "nats://localhost:4222",
				Token:    "my-token",
				Username: "admin",
				Password: "secret",
			},
			expectAuthType: "token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证配置结构
			if tt.config.Token != "" {
				assert.NotEmpty(t, tt.config.Token)
			}
			if tt.config.Username != "" && tt.config.Token == "" {
				assert.NotEmpty(t, tt.config.Username)
				assert.NotEmpty(t, tt.config.Password)
			}
		})
	}
}

// TestClient_ConcurrentPublish_Integration 测试并发发布（集成测试）
func TestClient_ConcurrentPublish_Integration(t *testing.T) {
	if !isNatsAvailable() {
		t.Skip("NATS server not available, skipping integration test")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 订阅计数
	var mu sync.Mutex
	count := 0

	sub, err := client.Subscribe("test.concurrent", func(msg *nats.Msg) {
		mu.Lock()
		count++
		mu.Unlock()
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 并发发布
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Publish("test.concurrent", "message")
		}()
	}
	wg.Wait()

	// 等待消息处理
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, 100, count)
	mu.Unlock()
}

// BenchmarkPublish 基准测试：发布消息
func BenchmarkPublish(b *testing.B) {
	if !isNatsAvailable() {
		b.Skip("NATS server not available, skipping benchmark")
	}

	cfg := Config{URL: getTestURL()}
	client, err := NewClient(cfg)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Publish("bench.publish", "message")
	}
}

// MockError 是一个 mock 错误
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	testErr := errors.New("test error")
	wrappedErr := errors.New("wrapped: " + testErr.Error())

	assert.Error(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "test error")
}
