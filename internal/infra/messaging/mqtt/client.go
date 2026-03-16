// Package mqtt provides MQTT client implementation for IoT and messaging applications.
//
// MQTT（Message Queuing Telemetry Transport）是一个轻量级的消息传输协议，
// 专为低带宽、高延迟或不稳定的网络环境设计。
//
// 设计原则：
// 1. 轻量级：协议简单，开销小，适合资源受限的设备
// 2. 可靠性：支持 QoS（Quality of Service）保证消息传递
// 3. 发布订阅：基于主题（Topic）的发布订阅模式
// 4. 自动重连：支持自动重连和连接保持
//
// 核心功能：
// - 发布消息：向指定主题发布消息
// - 订阅主题：订阅主题并接收消息
// - 连接管理：自动重连、心跳保持等
//
// 使用场景：
// - IoT 设备通信：传感器数据采集、设备控制
// - 实时消息推送：通知、告警等
// - 移动应用：消息推送、状态同步
// - 边缘计算：边缘设备与云端通信
//
// 示例：
//
//	// 创建客户端
//	client, err := mqtt.NewClient("tcp://localhost:1883", "client-id", "user", "pass")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// 订阅主题
//	handler := func(ctx context.Context, topic string, payload []byte) error {
//	    log.Printf("Received: topic=%s, payload=%s", topic, string(payload))
//	    return nil
//	}
//	err = client.Subscribe(ctx, "sensors/temperature", 1, handler)
//
//	// 发布消息
//	err = client.Publish(ctx, "sensors/temperature", 1, false, "25.5")
package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client 是 MQTT 客户端的封装，用于与 MQTT Broker 通信。
//
// 功能说明：
// - 连接到 MQTT Broker
// - 发布消息到指定主题
// - 订阅主题并接收消息
// - 自动重连和连接保持
//
// 设计说明：
// - 基于 Eclipse Paho MQTT Go 客户端
// - 支持自动重连和连接保持
// - 支持 QoS 0、1、2 三种消息质量等级
//
// 使用示例：
//
//	client, err := mqtt.NewClient("tcp://localhost:1883", "client-id", "", "")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
type Client struct {
	client mqtt.Client
}

// MessageHandler 是消息处理函数的类型定义。
//
// 功能说明：
// - 处理从 MQTT Broker 接收到的消息
// - 接收主题和消息负载（原始字节）
// - 返回处理错误（如果处理失败）
//
// 参数：
// - ctx: 上下文，用于控制处理超时和取消
// - topic: 消息主题
// - payload: 消息负载（原始字节）
//
// 返回：
// - error: 如果处理失败，返回错误信息
//
// 使用示例：
//
//	handler := func(ctx context.Context, topic string, payload []byte) error {
//	    log.Printf("Received message on topic %s: %s", topic, string(payload))
//	    return nil
//	}
type MessageHandler func(ctx context.Context, topic string, payload []byte) error

// NewClient 创建并连接到 MQTT Broker 的客户端。
//
// 功能说明：
// - 配置客户端选项（Broker 地址、认证信息等）
// - 连接到 MQTT Broker
// - 配置自动重连和连接保持
//
// 参数：
// - broker: MQTT Broker 地址
//   格式：tcp://host:port 或 ssl://host:port
//   示例：tcp://localhost:1883、ssl://mqtt.example.com:8883
// - clientID: 客户端 ID（必须唯一）
//   用于标识客户端，相同 ID 的客户端会互相踢下线
// - username: 用户名（可选，如果为空则不使用认证）
// - password: 密码（可选，如果为空则不使用密码）
//
// 返回：
// - *Client: 配置好的客户端实例
// - error: 如果连接失败，返回错误信息
//
// 配置说明：
// - AutoReconnect: 自动重连（默认启用）
// - ConnectRetry: 连接重试（默认启用）
// - ConnectRetryInterval: 重试间隔（默认 5 秒）
// - KeepAlive: 心跳间隔（默认 60 秒）
// - PingTimeout: Ping 超时（默认 10 秒）
//
// 使用示例：
//
//	// 连接到本地 MQTT Broker（无认证）
//	client, err := mqtt.NewClient("tcp://localhost:1883", "my-client", "", "")
//
//	// 连接到远程 MQTT Broker（带认证）
//	client, err := mqtt.NewClient(
//	    "ssl://mqtt.example.com:8883",
//	    "my-client",
//	    "username",
//	    "password",
//	)
//
// 注意事项：
// - 确保 MQTT Broker 已启动并可访问
// - 客户端 ID 应具有唯一性，避免冲突
// - 生产环境建议使用 TLS/SSL 连接
// - 应在应用程序生命周期中复用客户端实例
func NewClient(broker, clientID, username, password string) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)                    // Broker 地址
	opts.SetClientID(clientID)                // 客户端 ID
	opts.SetUsername(username)                // 用户名
	opts.SetPassword(password)                // 密码
	opts.SetAutoReconnect(true)               // 自动重连
	opts.SetConnectRetry(true)                // 连接重试
	opts.SetConnectRetryInterval(5 * time.Second) // 重试间隔
	opts.SetKeepAlive(60 * time.Second)       // 心跳间隔
	opts.SetPingTimeout(10 * time.Second)     // Ping 超时
	// 其他可选配置：
	// opts.SetCleanSession(true)              // 清理会话
	// opts.SetOrderMatters(false)             // 消息顺序
	// opts.SetResumeSubs(true)                // 恢复订阅

	client := mqtt.NewClient(opts)

	// 连接到 Broker
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect: %w", token.Error())
	}

	return &Client{client: client}, nil
}

// Publish 发布消息到指定的 MQTT 主题。
//
// 功能说明：
// - 将消息发布到指定的主题
// - 支持多种负载类型（字节数组、字符串、JSON 对象等）
// - 支持 QoS 和保留消息（Retained Message）
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消（当前未使用，保留用于未来扩展）
// - topic: MQTT 主题名称
//   支持通配符：+（单级）、#（多级）
//   示例：sensors/temperature、sensors/+/temperature
// - qos: 消息质量等级
//   - 0: 最多一次（At most once），不保证消息到达
//   - 1: 至少一次（At least once），保证消息至少到达一次
//   - 2: 恰好一次（Exactly once），保证消息恰好到达一次
// - retained: 是否保留消息
//   - true: 保留消息，新订阅者会收到最后一条保留消息
//   - false: 不保留消息
// - payload: 消息负载
//   支持类型：
//   - []byte: 原始字节数组
//   - string: 字符串（自动转换为字节数组）
//   - 其他类型: 自动序列化为 JSON
//
// 返回：
// - error: 如果发布失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//
//	// 发布字符串消息
//	err := client.Publish(ctx, "sensors/temperature", 1, false, "25.5")
//
//	// 发布 JSON 对象
//	data := map[string]interface{}{
//	    "sensor_id": "sensor-001",
//	    "value":     25.5,
//	    "unit":      "celsius",
//	}
//	err = client.Publish(ctx, "sensors/temperature", 1, false, data)
//
//	// 发布保留消息
//	err = client.Publish(ctx, "status/online", 1, true, "true")
//
// 注意事项：
// - QoS 0 最快但不保证消息到达
// - QoS 1 保证消息至少到达一次，但可能重复
// - QoS 2 保证消息恰好到达一次，但性能较低
// - 保留消息会占用 Broker 存储空间
func (c *Client) Publish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	var data []byte
	var err error

	// 根据负载类型进行转换
	switch v := payload.(type) {
	case []byte:
		// 已经是字节数组，直接使用
		data = v
	case string:
		// 字符串，转换为字节数组
		data = []byte(v)
	default:
		// 其他类型，序列化为 JSON
		data, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	// 发布消息
	token := c.client.Publish(topic, qos, retained, data)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish: %w", token.Error())
	}

	return nil
}

// Subscribe 订阅 MQTT 主题并设置消息处理函数。
//
// 功能说明：
// - 订阅指定的主题
// - 当收到消息时，调用处理函数处理消息
// - 支持通配符订阅（+、#）
//
// 参数：
// - ctx: 上下文，用于控制订阅超时和取消（当前未使用，保留用于未来扩展）
// - topic: MQTT 主题名称或通配符
//   支持通配符：
//   - +: 单级通配符，匹配一个主题级别
//     示例：sensors/+/temperature 匹配 sensors/room1/temperature
//   - #: 多级通配符，匹配多个主题级别
//     示例：sensors/# 匹配 sensors/room1/temperature、sensors/room2/humidity 等
// - qos: 订阅的 QoS 等级（0、1、2）
//   订阅的 QoS 不能高于发布消息的 QoS
// - handler: 消息处理函数
//   当收到消息时，会调用此函数处理消息
//
// 返回：
// - error: 如果订阅失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//
//	// 订阅单个主题
//	handler := func(ctx context.Context, topic string, payload []byte) error {
//	    log.Printf("Received: %s = %s", topic, string(payload))
//	    return nil
//	}
//	err := client.Subscribe(ctx, "sensors/temperature", 1, handler)
//
//	// 订阅通配符主题
//	err = client.Subscribe(ctx, "sensors/+", 1, handler)
//	err = client.Subscribe(ctx, "sensors/#", 1, handler)
//
// 注意事项：
// - 订阅的 QoS 不能高于发布消息的 QoS
// - 处理函数中的错误会被记录，但不会中断订阅
// - 可以多次订阅不同的主题
// - 使用 Unsubscribe 取消订阅
func (c *Client) Subscribe(ctx context.Context, topic string, qos byte, handler MessageHandler) error {
	// 订阅主题并设置消息回调
	token := c.client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		// 调用处理函数处理消息
		if err := handler(ctx, msg.Topic(), msg.Payload()); err != nil {
			// 记录错误，但不中断订阅
			// 生产环境应使用日志库记录错误
			fmt.Printf("Error handling message: %v\n", err)
		}
	})

	// 等待订阅完成
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe: %w", token.Error())
	}

	return nil
}

// Unsubscribe 取消订阅一个或多个 MQTT 主题。
//
// 功能说明：
// - 取消之前订阅的主题
// - 支持批量取消订阅
// - 取消订阅后不再接收该主题的消息
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消（当前未使用，保留用于未来扩展）
// - topics: 要取消订阅的主题列表（可变参数）
//
// 返回：
// - error: 如果取消订阅失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//
//	// 取消订阅单个主题
//	err := client.Unsubscribe(ctx, "sensors/temperature")
//
//	// 批量取消订阅
//	err = client.Unsubscribe(ctx, "sensors/temperature", "sensors/humidity")
//
// 注意事项：
// - 取消订阅后，该主题的消息处理函数不会再被调用
// - 可以取消订阅不存在的主题（不会报错）
func (c *Client) Unsubscribe(ctx context.Context, topics ...string) error {
	token := c.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe: %w", token.Error())
	}
	return nil
}

// Close 关闭 MQTT 客户端连接。
//
// 功能说明：
// - 断开与 MQTT Broker 的连接
// - 发送断开连接消息
// - 释放客户端资源
//
// 参数说明：
// - Disconnect 方法的参数是断开连接的等待时间（毫秒）
//   250 毫秒是默认值，等待未完成的消息发送完成
//
// 使用示例：
//
//	defer client.Close()
//
// 注意事项：
// - 应在应用程序退出前调用
// - 关闭后不应再使用该客户端
// - 关闭会发送断开连接消息给 Broker
func (c *Client) Close() {
	c.client.Disconnect(250) // 等待 250 毫秒后断开连接
}
