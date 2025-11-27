# Infrastructure Layer (基础设施层)

Clean Architecture 的基础设施层，包含技术栈的实现。

## ⚠️ 重要说明

**本框架提供技术栈的实现**，不包含具体业务的数据模型。用户需要根据自己的业务需求定义 Ent Schema 和实现具体的仓储。

## 结构

```
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

## 相关资源

- [Ent Schema 定义示例](../../examples/ent-schema/) - Ent Schema 定义示例
- [仓储实现示例](../../examples/repository/) - 仓储实现示例
