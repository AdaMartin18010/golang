// Package mqtt provides tests for the factory pattern used in NewClient.
package mqtt

import (
	"errors"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient_Success_WithFactory 测试使用factory成功创建客户端
func TestNewClient_Success_WithFactory(t *testing.T) {
	// 保存原始factory
	originalFactory := defaultFactory
	defer func() {
		defaultFactory = originalFactory
	}()

	// 设置mock factory
	mock := newMockMqttClient()
	defaultFactory = func(options *mqtt.ClientOptions) mqtt.Client {
		return mock
	}

	client, err := NewClient("tcp://localhost:1883", "test-client", "", "")
	require.NoError(t, err)
	assert.NotNil(t, client)
}

// TestNewClient_ConnectError_WithFactory 测试连接失败
func TestNewClient_ConnectError_WithFactory(t *testing.T) {
	// 保存原始factory
	originalFactory := defaultFactory
	defer func() {
		defaultFactory = originalFactory
	}()

	// 设置mock factory返回一个会连接失败的客户端
	mock := newMockMqttClientWithConnect(errors.New("connection refused"))
	defaultFactory = func(options *mqtt.ClientOptions) mqtt.Client {
		return mock
	}

	client, err := NewClient("tcp://localhost:1883", "test-client", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
	assert.Nil(t, client)
}

// TestNewClient_WithCredentials 测试带认证信息创建客户端
func TestNewClient_WithCredentials(t *testing.T) {
	// 保存原始factory
	originalFactory := defaultFactory
	defer func() {
		defaultFactory = originalFactory
	}()

	// 设置mock factory并验证选项
	mock := newMockMqttClient()
	defaultFactory = func(options *mqtt.ClientOptions) mqtt.Client {
		// 验证选项被正确设置
		assert.NotNil(t, options)
		return mock
	}

	client, err := NewClient("tcp://localhost:1883", "test-client", "username", "password")
	require.NoError(t, err)
	assert.NotNil(t, client)
}

// TestSetMqttClientFactory 测试设置factory函数
func TestSetMqttClientFactory(t *testing.T) {
	// 保存原始factory
	originalFactory := defaultFactory
	defer func() {
		defaultFactory = originalFactory
	}()

	// 创建自定义factory
	customCalled := false
	customFactory := func(options *mqtt.ClientOptions) mqtt.Client {
		customCalled = true
		return newMockMqttClient()
	}

	// 设置自定义factory
	setMqttClientFactory(customFactory)

	// 调用NewClient
	_, _ = NewClient("tcp://localhost:1883", "test-client", "", "")

	// 验证自定义factory被调用
	assert.True(t, customCalled)
}

// TestResetMqttClientFactory 测试重置factory函数
func TestResetMqttClientFactory(t *testing.T) {
	// 注意：不能比较函数指针，所以这里只验证调用不会panic
	// 设置自定义factory
	setMqttClientFactory(func(options *mqtt.ClientOptions) mqtt.Client {
		return newMockMqttClient()
	})

	// 重置factory - 不应panic
	resetMqttClientFactory()

	// 验证factory被重置后NewClient可以正常工作（使用默认factory）
	// 由于我们无法比较函数，这里只确保没有panic
	assert.NotNil(t, defaultFactory)
}

// TestNewClient_WithSSL 测试SSL连接配置
func TestNewClient_WithSSL(t *testing.T) {
	// 保存原始factory
	originalFactory := defaultFactory
	defer func() {
		defaultFactory = originalFactory
	}()

	mock := newMockMqttClient()
	defaultFactory = func(options *mqtt.ClientOptions) mqtt.Client {
		return mock
	}

	client, err := NewClient("ssl://mqtt.example.com:8883", "ssl-client", "user", "pass")
	require.NoError(t, err)
	assert.NotNil(t, client)
}

// TestNewClient_DifferentBrokers 测试不同的broker配置
func TestNewClient_DifferentBrokers(t *testing.T) {
	// 保存原始factory
	originalFactory := defaultFactory
	defer func() {
		defaultFactory = originalFactory
	}()

	mock := newMockMqttClient()
	defaultFactory = func(options *mqtt.ClientOptions) mqtt.Client {
		return mock
	}

	tests := []struct {
		name     string
		broker   string
		clientID string
	}{
		{"localhost", "tcp://localhost:1883", "client1"},
		{"remote", "tcp://remote.broker:1883", "client2"},
		{"websocket", "ws://localhost:8080/mqtt", "client3"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(tc.broker, tc.clientID, "", "")
			require.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}
