// Package nats provides NATS client implementation for high-performance messaging.
//
// NATS（NATS Messaging System）是一个高性能、云原生的消息传递系统，
// 专为微服务、IoT 和云原生应用设计。
//
// 设计原则：
// 1. 高性能：微秒级延迟，高吞吐量
// 2. 轻量级：协议简单，资源占用小
// 3. 云原生：支持集群、流式处理
// 4. 可靠性：支持自动重连和连接保持
//
// 核心功能：
// - 发布消息：向指定主题发布消息
// - 订阅主题：订阅主题并接收消息
// - Request/Reply：请求-响应模式
// - 队列订阅：负载均衡订阅
//
// 使用场景：
// - 微服务间通信：服务发现、事件通知
// - 实时消息推送：低延迟消息传递
// - IoT 设备通信：设备状态同步
// - 云原生应用：服务网格、配置分发
//
// 示例：
//
//	// 创建客户端
//	client, err := nats.NewClient(nats.DefaultConfig())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// 发布消息
//	err = client.Publish("user.created", map[string]interface{}{
//	    "user_id": 123,
//	    "name":    "Alice",
//	})
//
//	// 订阅消息
//	sub, err := client.Subscribe("user.created", func(msg *nats.Msg) {
//	    var data map[string]interface{}
//	    json.Unmarshal(msg.Data, &data)
//	    log.Printf("Received: %+v", data)
//	})
//	defer sub.Unsubscribe()
package nats

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// Client NATS 客户端封装
type Client struct {
	conn *nats.Conn
}

// NewClient 创建 NATS 客户端
func NewClient(cfg Config) (*Client, error) {
	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.Timeout(cfg.Timeout),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			// 记录断开连接错误（可以集成日志系统）
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			// 记录重连成功（可以集成日志系统）
		}),
	}

	// 认证配置
	if cfg.Token != "" {
		opts = append(opts, nats.Token(cfg.Token))
	} else if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, nats.UserInfo(cfg.Username, cfg.Password))
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &Client{conn: conn}, nil
}

// Publish 发布消息到指定主题
func (c *Client) Publish(subject string, data interface{}) error {
	var payload []byte
	var err error

	switch v := data.(type) {
	case []byte:
		payload = v
	case string:
		payload = []byte(v)
	default:
		payload, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
	}

	return c.conn.Publish(subject, payload)
}

// Subscribe 订阅主题
func (c *Client) Subscribe(subject string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}

// QueueSubscribe 队列订阅（负载均衡）
func (c *Client) QueueSubscribe(subject, queue string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	return c.conn.QueueSubscribe(subject, queue, handler)
}

// Request 发送请求并等待响应（Request-Reply 模式）
func (c *Client) Request(subject string, data interface{}, timeout time.Duration) (*nats.Msg, error) {
	var payload []byte
	var err error

	switch v := data.(type) {
	case []byte:
		payload = v
	case string:
		payload = []byte(v)
	default:
		payload, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}
	}

	msg, err := c.conn.Request(subject, payload, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return msg, nil
}

// Close 关闭连接
func (c *Client) Close() {
	c.conn.Close()
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	return c.conn.IsConnected()
}

// Stats 获取连接统计信息
func (c *Client) Stats() nats.Stats {
	return c.conn.Stats()
}

