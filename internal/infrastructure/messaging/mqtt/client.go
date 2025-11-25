package mqtt

import (
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client MQTT 客户端
type Client struct {
	client mqtt.Client
}

// Config 配置
type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

// MessageHandler 消息处理函数
type MessageHandler func(ctx context.Context, topic string, payload []byte)

// NewClient 创建 MQTT 客户端
func NewClient(cfg Config) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		// 默认处理函数
	})
	opts.SetPingTimeout(1 * time.Second)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return &Client{
		client: client,
	}, nil
}

// Subscribe 订阅主题
func (c *Client) Subscribe(topic string, qos byte, handler MessageHandler) error {
	token := c.client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		ctx := context.Background()
		handler(ctx, msg.Topic(), msg.Payload())
	})

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe: %w", token.Error())
	}

	return nil
}

// Publish 发布消息
func (c *Client) Publish(topic string, qos byte, retained bool, payload []byte) error {
	token := c.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish: %w", token.Error())
	}
	return nil
}

// Disconnect 断开连接
func (c *Client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	}
}
