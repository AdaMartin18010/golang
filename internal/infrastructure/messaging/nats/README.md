# NATS 客户端

NATS 客户端实现，支持发布/订阅、Request/Reply 等模式。

## 功能特性

- ✅ 连接管理
- ✅ 发布/订阅
- ✅ Request/Reply 模式
- ✅ 队列订阅（负载均衡）
- ✅ 自动重连
- ✅ 连接统计

## 使用示例

### 基本使用

```go
import (
    "log"
    "github.com/yourusername/golang/internal/infrastructure/messaging/nats"
)

// 创建客户端
client, err := nats.NewClient(nats.DefaultConfig())
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 发布消息
err = client.Publish("user.created", map[string]interface{}{
    "user_id": 123,
    "name":    "Alice",
})
if err != nil {
    log.Fatal(err)
}

// 订阅消息
sub, err := client.Subscribe("user.created", func(msg *nats.Msg) {
    var data map[string]interface{}
    json.Unmarshal(msg.Data, &data)
    log.Printf("Received: %+v", data)
})
if err != nil {
    log.Fatal(err)
}
defer sub.Unsubscribe()
```

### 自定义配置

```go
config := nats.Config{
    URL:           "nats://localhost:4222",
    MaxReconnects: 10,
    ReconnectWait: 2 * time.Second,
    Timeout:       5 * time.Second,
    Name:          "my-client",
    Username:      "user",
    Password:      "pass",
}

client, err := nats.NewClient(config)
```

### Request/Reply 模式

```go
// 服务端：订阅并回复
sub, err := client.Subscribe("user.get", func(msg *nats.Msg) {
    userID := string(msg.Data)
    // 处理请求
    response := fmt.Sprintf("user:%s", userID)
    msg.Respond([]byte(response))
})

// 客户端：发送请求
reply, err := client.Request("user.get", "123", 5*time.Second)
if err != nil {
    log.Fatal(err)
}
log.Printf("Response: %s", string(reply.Data))
```

### 队列订阅（负载均衡）

```go
// 多个订阅者使用相同的队列组名，消息会被负载均衡
sub, err := client.QueueSubscribe("tasks", "worker-group", func(msg *nats.Msg) {
    // 处理任务
    log.Printf("Processing task: %s", string(msg.Data))
})
```

## 配置说明

### Config 字段

- `URL`: NATS 服务器地址（默认: "nats://localhost:4222"）
- `MaxReconnects`: 最大重连次数（-1 表示无限重连）
- `ReconnectWait`: 重连等待时间
- `Timeout`: 连接超时时间
- `Name`: 客户端名称
- `Token`: 认证 Token（可选）
- `Username`: 用户名（可选）
- `Password`: 密码（可选）

## 相关资源

- [NATS 官方文档](https://docs.nats.io/)
- [NATS Go 客户端](https://github.com/nats-io/nats.go)
