# 1. 💬 MQTT 深度解析

> **简介**: 本文档详细阐述了 MQTT 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 💬 MQTT 深度解析](#1--mqtt-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 客户端连接](#131-客户端连接)
    - [1.3.2 发布消息](#132-发布消息)
    - [1.3.3 订阅主题](#133-订阅主题)
    - [1.3.4 QoS 级别使用](#134-qos-级别使用)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 主题设计最佳实践](#141-主题设计最佳实践)
    - [1.4.2 QoS 选择最佳实践](#142-qos-选择最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**MQTT 是什么？**

MQTT 是一个轻量级的消息传输协议，适用于 IoT 场景。

**核心特性**:

- ✅ **轻量级**: 协议简单，开销小
- ✅ **QoS**: 支持三种 QoS 级别
- ✅ **主题**: 基于主题的发布/订阅
- ✅ **低带宽**: 适合低带宽场景

---

## 1.2 选型论证

**为什么选择 MQTT？**

**论证矩阵**:

| 评估维度 | 权重 | MQTT | CoAP | AMQP | HTTP | 说明 |
|---------|------|------|------|------|------|------|
| **IoT 适配** | 35% | 10 | 9 | 5 | 3 | MQTT 最适合 IoT |
| **低带宽** | 25% | 10 | 9 | 6 | 4 | MQTT 协议开销小 |
| **QoS 支持** | 20% | 10 | 7 | 9 | 5 | MQTT QoS 完善 |
| **易用性** | 15% | 9 | 7 | 6 | 8 | MQTT 简单易用 |
| **生态支持** | 5% | 9 | 6 | 8 | 10 | MQTT 生态良好 |
| **加权总分** | - | **9.60** | 8.20 | 6.60 | 4.80 | MQTT 得分最高 |

**核心优势**:

1. **IoT 适配（权重 35%）**:
   - 专为 IoT 场景设计
   - 支持大量并发连接
   - 适合资源受限设备

2. **低带宽（权重 25%）**:
   - 协议开销小
   - 适合低带宽场景
   - 支持压缩

3. **QoS 支持（权重 20%）**:
   - 三种 QoS 级别
   - 保证消息可靠性
   - 适合不同场景需求

**为什么不选择其他协议？**

1. **CoAP**:
   - ✅ 专为 IoT 设计
   - ❌ 生态不如 MQTT 丰富
   - ❌ 使用不如 MQTT 广泛

2. **AMQP**:
   - ✅ 功能强大
   - ❌ 协议复杂，开销大
   - ❌ 不适合 IoT 场景

3. **HTTP**:
   - ✅ 标准协议，生态丰富
   - ❌ 协议开销大
   - ❌ 不适合 IoT 场景

---

## 1.3 实际应用

### 1.3.1 客户端连接

**连接示例**:

```go
// internal/infrastructure/messaging/mqtt/client.go
package mqtt

import (
    "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
    client mqtt.Client
}

func NewClient(broker string, clientID string) (*Client, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(broker)
    opts.SetClientID(clientID)
    opts.SetCleanSession(true)
    opts.SetAutoReconnect(true)

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return &Client{client: client}, nil
}
```

### 1.3.2 发布消息

**发布消息示例**:

```go
// 发布消息
func (c *Client) Publish(topic string, qos byte, retained bool, payload []byte) error {
    token := c.client.Publish(topic, qos, retained, payload)
    token.Wait()
    return token.Error()
}

// 使用示例
client.Publish("sensors/temperature", 1, false, []byte("25.5"))
```

### 1.3.3 订阅主题

**订阅主题示例**:

```go
// 订阅主题
func (c *Client) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
    token := c.client.Subscribe(topic, qos, handler)
    token.Wait()
    return token.Error()
}

// 使用示例
client.Subscribe("sensors/+", 1, func(client mqtt.Client, msg mqtt.Message) {
    logger.Info("Received message",
        "topic", msg.Topic(),
        "payload", string(msg.Payload()),
    )
})
```

### 1.3.4 QoS 级别使用

**QoS 级别说明**:

```go
// QoS 0: 最多一次，不保证消息到达
client.Publish("topic", 0, false, payload)

// QoS 1: 至少一次，保证消息至少到达一次
client.Publish("topic", 1, false, payload)

// QoS 2: 恰好一次，保证消息恰好到达一次
client.Publish("topic", 2, false, payload)
```

---

## 1.4 最佳实践

### 1.4.1 主题设计最佳实践

**为什么需要良好的主题设计？**

良好的主题设计可以提高消息路由的效率和可维护性。

**主题设计原则**:

1. **层次结构**: 使用层次结构组织主题
2. **命名规范**: 使用清晰的命名规范
3. **通配符使用**: 合理使用通配符
4. **主题长度**: 控制主题长度

**实际应用示例**:

```go
// 主题设计最佳实践
// 层次结构: {location}/{device_type}/{device_id}/{sensor_type}
// 示例: home/thermostat/living-room/temperature

// 订阅所有温度传感器
client.Subscribe("+/+/+/temperature", 1, handler)

// 订阅特定设备的所有传感器
client.Subscribe("home/thermostat/living-room/+", 1, handler)
```

**最佳实践要点**:

1. **层次结构**: 使用层次结构组织主题，便于管理和订阅
2. **命名规范**: 使用清晰的命名规范，便于理解
3. **通配符使用**: 合理使用通配符（+、#），提高订阅效率
4. **主题长度**: 控制主题长度，避免过长

### 1.4.2 QoS 选择最佳实践

**为什么需要合理选择 QoS？**

合理选择 QoS 可以平衡消息可靠性和性能。

**QoS 选择原则**:

1. **QoS 0**: 适用于不重要的数据，如传感器数据
2. **QoS 1**: 适用于重要数据，如控制命令
3. **QoS 2**: 适用于关键数据，如支付信息

**实际应用示例**:

```go
// QoS 选择最佳实践
// 传感器数据：QoS 0（不重要的数据）
client.Publish("sensors/temperature", 0, false, temperatureData)

// 控制命令：QoS 1（重要的数据）
client.Publish("devices/thermostat/command", 1, false, commandData)

// 关键数据：QoS 2（关键的数据）
client.Publish("payment/transaction", 2, false, paymentData)
```

**最佳实践要点**:

1. **QoS 0**: 适用于不重要的数据，提高性能
2. **QoS 1**: 适用于重要数据，平衡可靠性和性能
3. **QoS 2**: 适用于关键数据，保证可靠性

---

## 📚 扩展阅读

- [MQTT 官方文档](https://mqtt.org/)
- [Paho Go MQTT 客户端](https://github.com/eclipse/paho.mqtt.golang)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 MQTT 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
