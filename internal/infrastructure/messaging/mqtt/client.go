package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client MQTT 客户端
type Client struct {
	client mqtt.Client
}

// MessageHandler 消息处理函数
type MessageHandler func(ctx context.Context, topic string, payload []byte) error

// NewClient 创建 MQTT 客户端
func NewClient(broker, clientID, username, password string) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect: %w", token.Error())
	}

	return &Client{client: client}, nil
}

// Publish 发布消息
func (c *Client) Publish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	var data []byte
	var err error

	switch v := payload.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	token := c.client.Publish(topic, qos, retained, data)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish: %w", token.Error())
	}

	return nil
}

// Subscribe 订阅主题
func (c *Client) Subscribe(ctx context.Context, topic string, qos byte, handler MessageHandler) error {
	token := c.client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		if err := handler(ctx, msg.Topic(), msg.Payload()); err != nil {
			// 记录错误，但不中断订阅
			fmt.Printf("Error handling message: %v\n", err)
		}
	})

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe: %w", token.Error())
	}

	return nil
}

// Unsubscribe 取消订阅
func (c *Client) Unsubscribe(ctx context.Context, topics ...string) error {
	token := c.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe: %w", token.Error())
	}
	return nil
}

// Close 关闭客户端
func (c *Client) Close() {
	c.client.Disconnect(250)
}
