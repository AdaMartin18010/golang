// Package mqtt provides mock implementations for unit testing.
package mqtt

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// mockToken implements mqtt.Token interface for testing
type mockToken struct {
	err   error
	done  chan struct{}
}

func newMockToken(err error) *mockToken {
	t := &mockToken{
		err:  err,
		done: make(chan struct{}),
	}
	close(t.done)
	return t
}

func (m *mockToken) Wait() bool {
	<-m.done
	return true
}

func (m *mockToken) WaitTimeout(timeout time.Duration) bool {
	select {
	case <-m.done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (m *mockToken) Done() <-chan struct{} {
	return m.done
}

func (m *mockToken) Error() error {
	return m.err
}

// mockMqttClient implements mqtt.Client interface for testing
type mockMqttClient struct {
	connected    bool
	publishFunc  func(topic string, qos byte, retained bool, payload interface{}) mqtt.Token
	subscribeFunc func(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token
	unsubscribeFunc func(topics ...string) mqtt.Token
}

func newMockMqttClient() *mockMqttClient {
	return &mockMqttClient{
		connected: true,
		publishFunc: func(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
			return newMockToken(nil)
		},
		subscribeFunc: func(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
			return newMockToken(nil)
		},
		unsubscribeFunc: func(topics ...string) mqtt.Token {
			return newMockToken(nil)
		},
	}
}

func (m *mockMqttClient) IsConnected() bool {
	return m.connected
}

func (m *mockMqttClient) IsConnectionOpen() bool {
	return m.connected
}

func (m *mockMqttClient) Connect() mqtt.Token {
	return newMockToken(nil)
}

func (m *mockMqttClient) Disconnect(quiesce uint) {
	m.connected = false
}

func (m *mockMqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	if m.publishFunc != nil {
		return m.publishFunc(topic, qos, retained, payload)
	}
	return newMockToken(nil)
}

func (m *mockMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	if m.subscribeFunc != nil {
		return m.subscribeFunc(topic, qos, callback)
	}
	return newMockToken(nil)
}

func (m *mockMqttClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return newMockToken(nil)
}

func (m *mockMqttClient) Unsubscribe(topics ...string) mqtt.Token {
	if m.unsubscribeFunc != nil {
		return m.unsubscribeFunc(topics...)
	}
	return newMockToken(nil)
}

func (m *mockMqttClient) AddRoute(topic string, callback mqtt.MessageHandler) {}

func (m *mockMqttClient) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}

// mockMessage implements mqtt.Message interface for testing
type mockMessage struct {
	topic   string
	payload []byte
	qos     byte
	retained bool
	msgID   uint16
}

func (m *mockMessage) Duplicate() bool { return false }
func (m *mockMessage) Qos() byte       { return m.qos }
func (m *mockMessage) Retained() bool  { return m.retained }
func (m *mockMessage) Topic() string   { return m.topic }
func (m *mockMessage) MessageID() uint16 { return m.msgID }
func (m *mockMessage) Payload() []byte { return m.payload }
func (m *mockMessage) Ack()            {}

// mockClientWithErrors creates a mock client that returns errors for testing error handling
func newMockClientWithErrors(publishErr, subscribeErr, unsubscribeErr error) *mockMqttClient {
	return &mockMqttClient{
		connected: true,
		publishFunc: func(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
			return newMockToken(publishErr)
		},
		subscribeFunc: func(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
			return newMockToken(subscribeErr)
		},
		unsubscribeFunc: func(topics ...string) mqtt.Token {
			return newMockToken(unsubscribeErr)
		},
	}
}

// newClientWithMock creates a Client with a mock mqtt client for testing
func newClientWithMock(mock mqttClient) *Client {
	return newClientWithClient(mock)
}

// mockClient creates a Client with a mock mqtt client for testing
func mockClient(m mqtt.Client) *Client {
	return &Client{client: m}
}

// Helper context for testing
func testContext() context.Context {
	return context.Background()
}

// mockMqttClientWithConnect is a mock that simulates Connect behavior
type mockMqttClientWithConnect struct {
	mockMqttClient
	connectErr error
}

func newMockMqttClientWithConnect(connectErr error) *mockMqttClientWithConnect {
	return &mockMqttClientWithConnect{
		mockMqttClient: mockMqttClient{
			connected: true,
			publishFunc: func(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
				return newMockToken(nil)
			},
			subscribeFunc: func(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
				return newMockToken(nil)
			},
			unsubscribeFunc: func(topics ...string) mqtt.Token {
				return newMockToken(nil)
			},
		},
		connectErr: connectErr,
	}
}

func (m *mockMqttClientWithConnect) Connect() mqtt.Token {
	return newMockToken(m.connectErr)
}
