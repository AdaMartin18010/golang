# TS-008: NATS Messaging Patterns - Architecture & Go Implementation

> **维度**: Technology Stack
> **级别**: S (16+ KB)
> **标签**: #nats #messaging #pubsub #jetstream #go
> **权威来源**:
>
> - [NATS Documentation](https://docs.nats.io/) - Synadia
> - [NATS Architecture](https://docs.nats.io/nats-concepts/architecture) - NATS.io
> - [JetStream Documentation](https://docs.nats.io/jetstream/jetstream) - NATS.io

---

## 1. NATS Core Architecture

### 1.1 Server Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    NATS Server Architecture                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Single NATS Server                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Client Connection Handling                                      │  │  │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐              │  │  │
│  │  │  │ Client  │ │ Client  │ │ Client  │ │ Client  │              │  │  │
│  │  │  │ Conn 1  │ │ Conn 2  │ │ Conn 3  │ │ Conn N  │              │  │  │
│  │  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘              │  │  │
│  │  │       └───────────┴───────────┴───────────┘                     │  │  │
│  │  │                   │                                             │  │  │
│  │  │       ┌───────────┴───────────┐                                 │  │  │
│  │  │       │   Read Loop (per conn)│  Parse protocol                │  │  │
│  │  │       └───────────┬───────────┘                                 │  │  │
│  │  │                   │                                             │  │  │
│  │  │       ┌───────────┴───────────┐                                 │  │  │
│  │  │       │   SUBS (Hash Map)     │  subject -> []subscribers     │  │  │
│  │  │       └───────────┬───────────┘                                 │  │  │
│  │  │                   │                                             │  │  │
│  │  │       ┌───────────┴───────────┐                                 │  │  │
│  │  │       │   Write Loop (per conn)│  Deliver messages             │  │  │
│  │  │       └───────────────────────┘                                 │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Key Characteristics:                                                  │  │
│  │  • Pure pub-sub: No persistence in core NATS                          │  │
│  │  • At-most-once delivery                                              │  │
│  │  • Fan-out to all matching subscribers                                │  │
│  │  • Subject-based addressing (dot notation)                            │  │
│  │  • Wildcards: * (single token), > (multi token)                       │  │
│  │                                                                        │  │
│  │  Example Subjects:                                                     │  │
│  │  • "orders.created"                                                    │  │
│  │  • "orders.processed.us"                                               │  │
│  │  • "orders.*.us"  (matches orders.created.us, orders.shipped.us)      │  │
│  │  • "orders.>"     (matches any orders.*.*)                            │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    NATS Cluster (Routes)                               │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │      Server A              Server B              Server C             │  │
│  │    ┌──────────┐          ┌──────────┐          ┌──────────┐          │  │
│  │    │ Client 1 │          │ Client 2 │          │ Client 3 │          │  │
│  │    │ Client 2 │◄────────►│ Client 3 │◄────────►│ Client 4 │          │  │
│  │    └──────────┘  Route   └──────────┘  Route   └──────────┘          │  │
│  │         ▲                  ▲    ▲                  ▲                 │  │
│  │         │                  │    │                  │                 │  │
│  │         └──────────────────┘    └──────────────────┘                 │  │
│  │                   Full mesh route connections                         │  │
│  │                                                                        │  │
│  │  Message Flow:                                                         │  │
│  │  1. Client 1 publishes to "orders.new" on Server A                    │  │
│  │  2. Server A forwards via routes to Server B and C                    │  │
│  │  3. All servers deliver to local matching subscribers                 │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Subject-Based Messaging

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    NATS Subject Hierarchy                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Subject Namespace                                   │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  orders.                                                              │  │
│  │    ├── created                                                        │  │
│  │    │     ├── us                                                       │  │
│  │    │     ├── eu                                                       │  │
│  │    │     └── asia                                                     │  │
│  │    ├── processed                                                      │  │
│  │    │     ├── us                                                       │  │
│  │    │     ├── eu                                                       │  │
│  │    │     └── asia                                                     │  │
│  │    ├── shipped                                                        │  │
│  │    └── cancelled                                                      │  │
│  │                                                                        │  │
│  │  payments.                                                            │  │
│  │    ├── authorized                                                     │  │
│  │    ├── captured                                                       │  │
│  │    └── refunded                                                       │  │
│  │                                                                        │  │
│  │  inventory.                                                           │  │
│  │    ├── updated                                                        │  │
│  │    └── lowstock                                                       │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Wildcard Matching                                   │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Publisher: PUB orders.created.us "{...}"                             │  │
│  │                                                                        │  │
│  │  Subscribers:                                                          │  │
│  │  ┌─────────────────────┬────────────────────────────────────────────┐  │  │
│  │  │ Subscription        │ Matches?                                   │  │  │
│  │  ├─────────────────────┼────────────────────────────────────────────┤  │  │
│  │  │ orders.created.us   │ ✓ Exact match                              │  │  │
│  │  │ orders.created.*    │ ✓ * matches any single token               │  │  │
│  │  │ orders.*.us         │ ✓ * matches "created"                      │  │  │
│  │  │ orders.>            │ ✓ > matches all remaining tokens           │  │  │
│  │  │ orders.created.eu   │ ✗ Different region                         │  │  │
│  │  │ payments.>          │ ✗ Different subject root                   │  │  │
│  │  └─────────────────────┴────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Multiple subscribers to same subject = fan-out                        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go Implementation

```go
package nats

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
)

// Client NATS 客户端封装
type Client struct {
    conn *nats.Conn
    js   jetstream.JetStream
}

// Config NATS 配置
type Config struct {
    URL           string
    MaxReconnects int
    ReconnectWait time.Duration
    Timeout       time.Duration
}

// NewClient 创建 NATS 客户端
func NewClient(cfg *Config) (*Client, error) {
    opts := []nats.Option{
        nats.MaxReconnects(cfg.MaxReconnects),
        nats.ReconnectWait(cfg.ReconnectWait),
        nats.Timeout(cfg.Timeout),
        nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
            fmt.Printf("Disconnected: %v\n", err)
        }),
        nats.ReconnectHandler(func(nc *nats.Conn) {
            fmt.Printf("Reconnected to %s\n", nc.ConnectedUrl())
        }),
    }

    conn, err := nats.Connect(cfg.URL, opts...)
    if err != nil {
        return nil, fmt.Errorf("connect failed: %w", err)
    }

    js, err := jetstream.New(conn)
    if err != nil {
        return nil, fmt.Errorf("jetstream init failed: %w", err)
    }

    return &Client{
        conn: conn,
        js:   js,
    }, nil
}

// Close 关闭连接
func (c *Client) Close() {
    c.conn.Close()
}

// Publish 发布消息
func (c *Client) Publish(subject string, data []byte) error {
    return c.conn.Publish(subject, data)
}

// PublishJSON 发布 JSON 消息
func (c *Client) PublishJSON(subject string, v interface{}) error {
    data, err := json.Marshal(v)
    if err != nil {
        return err
    }
    return c.conn.Publish(subject, data)
}

// Subscribe 订阅消息
func (c *Client) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
    return c.conn.Subscribe(subject, handler)
}

// QueueSubscribe 队列订阅 (load balancing)
func (c *Client) QueueSubscribe(subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
    return c.conn.QueueSubscribe(subject, queue, handler)
}

// Request 请求-响应模式
func (c *Client) Request(subject string, data []byte, timeout time.Duration) (*nats.Msg, error) {
    return c.conn.Request(subject, data, timeout)
}

// ==================== JetStream Implementation ====================

// StreamConfig 流配置
type StreamConfig struct {
    Name     string
    Subjects []string
    Retention jetstream.RetentionPolicy
    MaxMsgs  int64
    MaxBytes int64
    MaxAge   time.Duration
}

// CreateStream 创建 JetStream
func (c *Client) CreateStream(ctx context.Context, cfg *StreamConfig) (jetstream.Stream, error) {
    return c.js.CreateStream(ctx, jetstream.StreamConfig{
        Name:      cfg.Name,
        Subjects:  cfg.Subjects,
        Retention: cfg.Retention,
        MaxMsgs:   cfg.MaxMsgs,
        MaxBytes:  cfg.MaxBytes,
        MaxAge:    cfg.MaxAge,
    })
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
    Stream        string
    Name          string
    Durable       string
    DeliverPolicy jetstream.DeliverPolicy
    AckPolicy     jetstream.AckPolicy
    MaxDeliver    int
    FilterSubject string
}

// CreateConsumer 创建消费者
func (c *Client) CreateConsumer(ctx context.Context, cfg *ConsumerConfig) (jetstream.Consumer, error) {
    return c.js.CreateConsumer(ctx, cfg.Stream, jetstream.ConsumerConfig{
        Name:          cfg.Name,
        Durable:       cfg.Durable,
        DeliverPolicy: cfg.DeliverPolicy,
        AckPolicy:     cfg.AckPolicy,
        MaxDeliver:    cfg.MaxDeliver,
        FilterSubject: cfg.FilterSubject,
    })
}

// ConsumeMessages 消费消息
func (c *Client) ConsumeMessages(ctx context.Context, consumer jetstream.Consumer, handler func(jetstream.Msg)) error {
    cons, err := consumer.Consume(handler)
    if err != nil {
        return err
    }
    defer cons.Stop()

    <-ctx.Done()
    return ctx.Err()
}

// PublishToStream 发布到流
func (c *Client) PublishToStream(ctx context.Context, subject string, data []byte) (*jetstream.PubAck, error) {
    return c.js.Publish(ctx, subject, data)
}

// ==================== Patterns ====================

// RequestReplyPattern 请求-响应模式
type RequestReplyPattern struct {
    client *Client
}

// NewRequestReplyPattern 创建请求-响应模式
func NewRequestReplyPattern(client *Client) *RequestReplyPattern {
    return &RequestReplyPattern{client: client}
}

// Request 发送请求
func (r *RequestReplyPattern) Request(subject string, req interface{}, resp interface{}, timeout time.Duration) error {
    reqData, err := json.Marshal(req)
    if err != nil {
        return err
    }

    msg, err := r.client.Request(subject, reqData, timeout)
    if err != nil {
        return err
    }

    return json.Unmarshal(msg.Data, resp)
}

// ReplyHandler 注册响应处理器
func (r *RequestReplyPattern) ReplyHandler(subject string, handler func([]byte) ([]byte, error)) (*nats.Subscription, error) {
    return r.client.Subscribe(subject, func(msg *nats.Msg) {
        resp, err := handler(msg.Data)
        if err != nil {
            msg.Respond([]byte(`{"error":"` + err.Error() + `"}`))
            return
        }
        msg.Respond(resp)
    })
}
```

---

## 3. Configuration Best Practices

```hcl
# NATS 服务器配置 nats-server.conf

# 网络配置
port: 4222
http_port: 8222

# 集群配置
cluster {
  name: "prod-cluster"
  listen: 0.0.0.0:6222
  routes: [
    "nats://nats-1:6222",
    "nats://nats-2:6222",
    "nats://nats-3:6222"
  ]
}

# JetStream 配置
jetstream {
  store_dir: "/data/jetstream"
  max_memory_store: 1GB
  max_file_store: 10GB
}

# 认证配置
authorization {
  users: [
    {user: admin, password: $ADMIN_PASS}
    {user: app, password: $APP_PASS, permissions: {
      publish: ["orders.>", "payments.>"]
      subscribe: ["orders.>", "inventory.>"]
    }}
  ]
}

# 监控
monitoring: true
```

---

## 4. Visual Representations

### JetStream Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    NATS JetStream Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Stream (Message Log)                                │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Stream: ORDERS                                                        │  │
│  │  Subjects: orders.>                                                    │  │
│  │  Retention: Limits / Interest / WorkQueue                             │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Seq  │ Subject            │ Data         │ Timestamp          │  │  │
│  │  ├─────────────────────────────────────────────────────────────────┤  │  │
│  │  │  1    │ orders.created.us  │ {...}        │ 2024-01-01T10:00:00│  │  │
│  │  │  2    │ orders.created.eu  │ {...}        │ 2024-01-01T10:00:01│  │  │
│  │  │  3    │ orders.processed.us│ {...}        │ 2024-01-01T10:01:00│  │  │
│  │  │  ...  │ ...                │ ...          │ ...                │  │  │
│  │  │  N    │ orders.shipped.eu  │ {...}        │ 2024-01-01T11:00:00│  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Consumer Groups                                     │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Stream: ORDERS                                                        │  │
│  │                                                                        │  │
│  │     ┌─────────────────┐                                               │  │
│  │     │ Consumer A      │                                               │  │
│  │     │ (Durable)       │                                               │  │
│  │     │ Deliver: All    │                                               │  │
│  │     │ Ack: Explicit   │                                               │  │
│  │     └────────┬────────┘                                               │  │
│  │              │                                                         │  │
│  │     ┌────────┴────────┐                                               │  │
│  │     │ Consumer B      │                                               │  │
│  │     │ (Durable)       │                                               │  │
│  │     │ Deliver: New    │                                               │  │
│  │     │ Ack: Explicit   │                                               │  │
│  │     └────────┬────────┘                                               │  │
│  │              │                                                         │  │
│  │              ▼                                                         │  │
│  │     Applications receive messages independently                        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. References

1. **NATS Documentation** (2024). docs.nats.io
2. **Synadia Communications** (2024). NATS Architecture Whitepaper
3. **Colby Toland, et al.** (2020). Pratical NATS. Apress.

---

*Document Version: 1.0 | Last Updated: 2024*
