package nats

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startTestServer 启动测试用的 NATS 服务器
// 注意：此函数需要 NATS 服务器运行，如果没有运行则跳过测试
func getTestURL() string {
	// 默认测试 URL，需要 NATS 服务器运行在 localhost:4222
	return "nats://localhost:4222"
}

func TestNewClient(t *testing.T) {
	cfg := Config{
		URL:           getTestURL(),
		MaxReconnects: 3,
		ReconnectWait: 1 * time.Second,
		Timeout:       5 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skip("NATS server not available, skipping test. Start NATS server with: docker run -p 4222:4222 nats:latest")
		return
	}
	defer client.Close()

	assert.NotNil(t, client)
	assert.True(t, client.IsConnected())
}

func TestPublishSubscribe(t *testing.T) {
	cfg := Config{
		URL: getTestURL(),
	}

	client, err := NewClient(cfg)
	require.NoError(t, err)
	defer client.Close()

	// 订阅
	received := make(chan bool, 1)
	sub, err := client.Subscribe("test.subject", func(msg *nats.Msg) {
		assert.Equal(t, "test.payload", string(msg.Data))
		received <- true
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布
	err = client.Publish("test.subject", "test.payload")
	require.NoError(t, err)

	// 等待接收
	select {
	case <-received:
		// 成功
	case <-time.After(1 * time.Second):
		t.Fatal("message not received")
	}
}

func TestPublishSubscribeJSON(t *testing.T) {
	cfg := Config{
		URL: getTestURL(),
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skip("NATS server not available, skipping test. Start NATS server with: docker run -p 4222:4222 nats:latest")
		return
	}
	defer client.Close()

	// 订阅
	received := make(chan map[string]interface{}, 1)
	sub, err := client.Subscribe("test.json", func(msg *nats.Msg) {
		var data map[string]interface{}
		err := json.Unmarshal(msg.Data, &data)
		require.NoError(t, err)
		received <- data
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布 JSON 消息
	message := map[string]interface{}{
		"user_id": 123,
		"name":    "Alice",
	}
	err = client.Publish("test.json", message)
	require.NoError(t, err)

	// 等待接收
	select {
	case data := <-received:
		assert.Equal(t, float64(123), data["user_id"])
		assert.Equal(t, "Alice", data["name"])
	case <-time.After(1 * time.Second):
		t.Fatal("message not received")
	}
}

func TestQueueSubscribe(t *testing.T) {
	cfg := Config{
		URL: getTestURL(),
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skip("NATS server not available, skipping test. Start NATS server with: docker run -p 4222:4222 nats:latest")
		return
	}
	defer client.Close()

	// 队列订阅
	received := make(chan bool, 1)
	sub, err := client.QueueSubscribe("test.queue", "queue-group", func(msg *nats.Msg) {
		assert.Equal(t, "queue.message", string(msg.Data))
		received <- true
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布
	err = client.Publish("test.queue", "queue.message")
	require.NoError(t, err)

	// 等待接收
	select {
	case <-received:
		// 成功
	case <-time.After(1 * time.Second):
		t.Fatal("message not received")
	}
}

func TestRequestReply(t *testing.T) {
	cfg := Config{
		URL: getTestURL(),
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skip("NATS server not available, skipping test. Start NATS server with: docker run -p 4222:4222 nats:latest")
		return
	}
	defer client.Close()

	// 订阅并回复
	sub, err := client.Subscribe("test.request", func(msg *nats.Msg) {
		// 回复消息
		response := "response: " + string(msg.Data)
		msg.Respond([]byte(response))
	})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发送请求
	reply, err := client.Request("test.request", "request data", 5*time.Second)
	require.NoError(t, err)
	assert.Equal(t, "response: request data", string(reply.Data))
}

func TestIsConnected(t *testing.T) {
	cfg := Config{
		URL: getTestURL(),
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skip("NATS server not available, skipping test. Start NATS server with: docker run -p 4222:4222 nats:latest")
		return
	}
	defer client.Close()

	assert.True(t, client.IsConnected())

	client.Close()
	assert.False(t, client.IsConnected())
}

func TestStats(t *testing.T) {
	cfg := Config{
		URL: getTestURL(),
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Skip("NATS server not available, skipping test. Start NATS server with: docker run -p 4222:4222 nats:latest")
		return
	}
	defer client.Close()

	stats := client.Stats()
	assert.NotNil(t, stats)
}
