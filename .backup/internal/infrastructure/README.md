# Infrastructure Layer (基础设施层)

Clean Architecture 的基础设施层，包含技术栈的实现。

## 📋 概述

基础设施层负责实现所有技术细节和外部依赖，包括数据库访问、消息队列、缓存、可观测性等。这一层是框架的核心，提供了完整的技术栈实现。

## ⚠️ 重要说明

**本框架提供技术栈的实现**，不包含具体业务的数据模型。用户需要根据自己的业务需求定义 Ent Schema 和实现具体的仓储。

## 🎯 设计原则

1. **技术实现隔离**：所有技术实现细节都在这一层
2. **接口实现**：实现领域层和应用层定义的接口
3. **可替换性**：可以轻松替换技术实现（如从 PostgreSQL 切换到 MySQL）
4. **无业务逻辑**：不包含任何业务逻辑，只负责技术实现

## 结构

```text
infrastructure/
├── database/      # 数据库实现
│   ├── postgres/  # PostgreSQL 连接管理
│   └── ent/       # Ent ORM 客户端和工具
├── messaging/     # 消息队列
│   ├── kafka/     # Kafka 生产者/消费者
│   └── mqtt/      # MQTT 客户端
├── cache/         # 缓存（待完善）
└── observability/ # 可观测性
    ├── otlp/      # OpenTelemetry 集成
    └── ebpf/      # eBPF 收集器
```

## 规则

- ✅ 实现技术栈的具体功能
- ✅ 包含技术实现细节
- ✅ 可以导入外部库
- ❌ 不包含具体业务的数据模型和仓储实现

## 数据库实现

### PostgreSQL

- **连接管理** (`database/postgres/connection.go`) - PostgreSQL 连接池管理
- **配置示例** - 提供连接配置示例

### SQLite3

- **连接管理** (`database/sqlite3/connection.go`) - SQLite3 连接池管理
- **配置支持** - 支持 WAL 模式、共享缓存等配置
- **使用示例** - 提供完整的使用示例

### Ent ORM

- **客户端** (`database/ent/`) - Ent 生成的客户端代码
- **工具函数** - Ent 辅助函数和工具
- **迁移工具** - 数据库迁移脚本

**注意**: Ent Schema 定义应该由用户在自己的项目中定义，框架不提供具体的 Schema 定义。示例请参考 `examples/ent-schema/`。

## 消息队列

### Kafka

- **生产者** (`messaging/kafka/producer.go`) - Kafka 消息生产者
- **消费者** (`messaging/kafka/consumer.go`) - Kafka 消息消费者

### MQTT

- **客户端** (`messaging/mqtt/client.go`) - MQTT 客户端封装

## 缓存

### Redis

- **客户端封装** (`cache/redis/client.go`) - Redis 客户端封装
- **连接管理** - 连接池管理和配置
- **常用操作** - Set、Get、Del、Exists 等常用操作封装
- **与中间件集成** - 支持限流中间件的分布式限流

## 可观测性

### OpenTelemetry

- **Logger** (`observability/otlp/logger.go`) - OpenTelemetry 日志集成
- **Metrics** (`observability/otlp/metrics.go`) - OpenTelemetry 指标集成
- **Tracer** (`observability/otlp/tracer.go`) - OpenTelemetry 追踪集成

### eBPF

- **收集器** (`observability/ebpf/collector.go`) - eBPF 数据收集器

## 用户如何使用

### 1. 定义 Ent Schema

用户在自己的项目中定义 Ent Schema：

```go
// 用户项目中的 Ent Schema
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("id").Unique().Immutable(),
        field.String("email").Unique().NotEmpty(),
        // ...
    }
}
```

### 2. 实现仓储

用户在自己的项目中实现仓储：

```go
// 用户项目中的仓储实现
package infrastructure

type UserRepository struct {
    client *ent.Client
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    // 使用 Ent 客户端实现
}
```

### 3. 使用消息队列

```go
// 使用 Kafka 生产者
producer := kafka.NewProducer(kafka.Config{...})
producer.Publish(ctx, "topic", message)

// 使用 MQTT 客户端
client := mqtt.NewClient(mqtt.Config{...})
client.Publish(ctx, "topic", message)
```

## 🔧 使用指南

### 1. 数据库连接

#### 1.1 PostgreSQL 连接示例

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/database/ent"
    "github.com/yourusername/golang/internal/infrastructure/database/postgres"
)

// 创建连接
conn, err := postgres.NewConnection(cfg.Database)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// 创建 Ent 客户端
client, err := ent.NewClientFromConfig(
    ctx,
    cfg.Database.Host,
    strconv.Itoa(cfg.Database.Port),
    cfg.Database.User,
    cfg.Database.Password,
    cfg.Database.Database,
    cfg.Database.SSLMode,
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 运行迁移
if err := client.Migrate(ctx); err != nil {
    log.Fatal(err)
}
```

#### 1.2 SQLite3 连接示例

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/database/sqlite3"
)

// 创建连接
conn, err := sqlite3.NewConnection(cfg.Database)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

### 2. 缓存使用

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/cache/redis"
)

// 创建 Redis 客户端
config := redis.DefaultConfig()
config.Addr = cfg.Redis.Addr
config.Password = cfg.Redis.Password
client, err := redis.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用缓存
ctx := context.Background()
err = client.Set(ctx, "key", "value", time.Hour)
value, err := client.Get(ctx, "key")
```

### 3. 消息队列使用

#### 3.1 Kafka 使用示例

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/messaging/kafka"
)

// 创建生产者
producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
if err != nil {
    log.Fatal(err)
}
defer producer.Close()

// 发送消息
ctx := context.Background()
err = producer.SendMessage(ctx, "topic", "key", messageData)

// 创建消费者
handler := func(ctx context.Context, key string, value []byte) error {
    // 处理消息
    return nil
}
consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, "group-id", handler)
if err != nil {
    log.Fatal(err)
}
defer consumer.Close()

// 消费消息
err = consumer.Consume(ctx, []string{"topic"})
```

#### 3.2 MQTT 使用示例

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/messaging/mqtt"
)

// 创建客户端
client, err := mqtt.NewClient(
    cfg.MQTT.Broker,
    cfg.MQTT.ClientID,
    cfg.MQTT.Username,
    cfg.MQTT.Password,
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 发布消息
ctx := context.Background()
err = client.Publish(ctx, "topic", 1, false, "message")

// 订阅主题
handler := func(ctx context.Context, topic string, payload []byte) error {
    // 处理消息
    return nil
}
err = client.Subscribe(ctx, "topic", 1, handler)
```

### 4. 可观测性使用

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/observability/otlp"
)

// 创建日志记录器
logger := otlp.NewLogger()
slog.SetDefault(logger.Logger)

// 创建追踪提供者
ctx := context.Background()
shutdownTracer, err := otlp.NewTracerProvider(
    ctx,
    cfg.OTLP.Endpoint,
    cfg.OTLP.Insecure,
)
if err != nil {
    log.Fatal(err)
}
defer shutdownTracer(ctx)

// 创建指标提供者
metricsProvider, err := otlp.NewMetricsProvider(
    ctx,
    cfg.OTLP.Endpoint,
    cfg.OTLP.Insecure,
)
if err != nil {
    log.Fatal(err)
}
defer metricsProvider.Shutdown(ctx)
```

### 5. 工作流使用

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/workflow/temporal"
)

// 创建客户端
client, err := temporal.NewClient(cfg.Temporal.Address)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 启动工作流
workflowOptions := client.StartWorkflowOptions{
    ID:        "workflow-id",
    TaskQueue: cfg.Temporal.TaskQueue,
}
workflowRun, err := client.ExecuteWorkflow(ctx, workflowOptions, MyWorkflow, input)

// 创建 Worker
worker := temporal.NewWorkerFromClient(client, cfg.Temporal.TaskQueue)
worker.RegisterWorkflow(MyWorkflow)
worker.RegisterActivity(MyActivity)
err = worker.Run()
```

## 📚 相关资源

- [Ent Schema 定义示例](../../examples/ent-schema/) - Ent Schema 定义示例
- [仓储实现示例](../../examples/repository/) - 仓储实现示例
- [配置管理](../config/) - 配置管理说明
- [应用层](../application/) - 应用层说明
- [接口层](../interfaces/) - 接口层说明
