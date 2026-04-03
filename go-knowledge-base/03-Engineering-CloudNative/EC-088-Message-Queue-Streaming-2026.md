# EC-088-Message-Queue-Streaming-2026

> **Dimension**: 03-Engineering-CloudNative
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: 2026 (Kafka, RabbitMQ, Pulsar, NATS, Redis Streams)
> **Size**: >20KB

---

## 1. 消息队列概览

### 1.1 分类

| 类型 | 代表产品 | 特点 |
|------|---------|------|
| 日志型 | Kafka, Pulsar | 高吞吐，持久化，流处理 |
| 队列型 | RabbitMQ, ActiveMQ | 路由灵活，协议丰富 |
| 内存型 | NATS, Redis Streams | 低延迟，轻量 |
| 云原生 | Kinesis, Event Hubs | 托管服务，集成生态 |

### 1.2 选型决策树

```
┌─────────────────────────────────────────┐
│         Message Queue Selection         │
├─────────────────────────────────────────┤
│                                         │
│  高吞吐(>1M/s)? ──Yes──► Kafka/Pulsar  │
│       │                                 │
│       No                                │
│       │                                 │
│  复杂路由? ──Yes──► RabbitMQ           │
│       │                                 │
│       No                                │
│       │                                 │
│  低延迟(<1ms)? ──Yes──► NATS/Redis     │
│       │                                 │
│       No                                │
│       │                                 │
│  云托管优先? ──Yes──► Kinesis/EventHub │
│       │                                 │
│       No                                │
│       ▼                                 │
│  推荐: RabbitMQ (通用场景)              │
│                                         │
└─────────────────────────────────────────┘
```

---

## 2. Apache Kafka 2026

### 2.1 Kafka 4.0 特性

**KRaft模式 GA**:

- 完全移除ZooKeeper依赖
- 简化部署
- 更好的可扩展性

**分层存储 (Tiered Storage)**:

```
┌─────────────────────────────────────────┐
│         Kafka Tiered Storage            │
├─────────────────────────────────────────┤
│                                         │
│  Hot Tier (Local)    Cold Tier (S3)     │
│  ┌─────────────┐     ┌─────────────┐    │
│  │ Recent Data │────►│ Historical  │    │
│  │ (<7 days)   │     │ Data        │    │
│  │ Low Latency │     │ (7+ days)   │    │
│  └─────────────┘     │ Archive     │    │
│                      └─────────────┘    │
│                                         │
│  优势:                                  │
│  - 无限存储扩展                         │
│  - 降低存储成本 80%+                    │
│  - 保留期从几天延长到几年               │
│                                         │
└─────────────────────────────────────────┘
```

### 2.2 生产者配置优化

```java
Properties props = new Properties();
props.put("bootstrap.servers", "kafka:9092");
props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");

// 高吞吐配置
props.put("batch.size", 65536);              // 64KB批次
props.put("linger.ms", 10);                   // 等待10ms批量
props.put("compression.type", "zstd");        // ZSTD压缩
props.put("acks", "1");                       // Leader确认即可
props.put("retries", 3);                      // 失败重试
props.put("max.in.flight.requests.per.connection", 5);

// 幂等性 (exactly-once)
props.put("enable.idempotence", "true");

Producer<String, String> producer = new KafkaProducer<>(props);
```

### 2.3 消费者组

```java
Properties props = new Properties();
props.put("bootstrap.servers", "kafka:9092");
props.put("group.id", "order-processing-group");
props.put("key.deserializer", "org.apache.kafka.common.serialization.StringDeserializer");
props.put("value.deserializer", "org.apache.kafka.common.serialization.StringDeserializer");

// 自动提交偏移量
props.put("enable.auto.commit", "false");  // 手动提交更安全
props.put("auto.offset.reset", "earliest");
props.put("max.poll.records", 500);
props.put("max.poll.interval.ms", 300000);

KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props);
consumer.subscribe(Arrays.asList("orders"));

// 手动提交模式
try {
    while (true) {
        ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));

        for (ConsumerRecord<String, String> record : records) {
            processRecord(record);
        }

        // 同步提交
        consumer.commitSync();
    }
} finally {
    consumer.close();
}
```

### 2.4 Exactly-Once语义

```java
// 事务性生产者
props.put("transactional.id", "order-producer-1");
KafkaProducer<String, String> producer = new KafkaProducer<>(props);
producer.initTransactions();

try {
    producer.beginTransaction();

    // 发送消息
    producer.send(new ProducerRecord<>("orders", orderId, orderData));

    // 发送偏移量到消费者组
    producer.sendOffsetsToTransaction(
        consumer.position(consumer.assignment()),
        consumer.groupMetadata()
    );

    producer.commitTransaction();
} catch (Exception e) {
    producer.abortTransaction();
    throw e;
}
```

---

## 3. Apache Pulsar

### 3.1 架构

```
┌─────────────────────────────────────────┐
│          Pulsar Architecture            │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────┐    ┌─────────────────┐    │
│  │ Producers│───►│  Broker Layer   │    │
│  └─────────┘    │  ┌─────────────┐ │    │
│                 │  │ Broker 1    │ │    │
│  ┌─────────┐    │  │ Broker 2    │ │    │
│  │Consumers│◄───│  │ Broker 3    │ │    │
│  └─────────┘    │  └─────────────┘ │    │
│                 │     ▲            │    │
│                 └─────┼────────────┘    │
│                       │                 │
│                 ┌─────┴────────────┐    │
│                 │  BookKeeper      │    │
│                 │  (Ledger Storage)│    │
│                 └──────────────────┘    │
│                                         │
│  多租户、多区域复制、分层存储            │
└─────────────────────────────────────────┘
```

### 3.2 统一消息模型

| 模式 | API | 适用场景 |
|------|-----|---------|
| 队列 | Shared Subscription | 工作队列，负载均衡 |
| 流 | Failover Subscription | 事件溯源，顺序处理 |
| 发布订阅 | Exclusive/Failover | 广播，扇出 |

### 3.3 Geo-Replication

```yaml
# 多区域复制配置
clusters:
  - name: us-east
    serviceUrl: http://pulsar-us-east:8080
  - name: us-west
    serviceUrl: http://pulsar-us-west:8080
  - name: eu-west
    serviceUrl: http://pulsar-eu-west:8080

# 复制策略
replication:
  - from: us-east
    to: [us-west, eu-west]
  - from: us-west
    to: [us-east]
```

---

## 4. RabbitMQ 4.0

### 4.1 新特性

**原生MQTT 5.0支持**:

- 不再需要插件
- 更好的物联网集成

**流队列增强**:

```
┌─────────────────────────────────────────┐
│        RabbitMQ Streams                 │
├─────────────────────────────────────────┤
│                                         │
│  append-only log结构                     │
│  非破坏性消费                            │
│  支持重放                                │
│  高吞吐(百万/秒)                         │
│                                         │
│  适用场景:                               │
│  - 事件溯源                             │
│  - 流处理                               │
│  - 大数据摄取                           │
│                                         │
└─────────────────────────────────────────┘
```

### 4.2 交换机类型

```python
import pika

connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
channel = connection.channel()

# 1. Direct Exchange - 精确匹配
channel.exchange_declare(exchange='direct_logs', exchange_type='direct')
channel.basic_publish(
    exchange='direct_logs',
    routing_key='error',  # 精确匹配
    body=message
)

# 2. Topic Exchange - 模式匹配
channel.exchange_declare(exchange='topic_logs', exchange_type='topic')
channel.basic_publish(
    exchange='topic_logs',
    routing_key='kern.critical',  # facility.severity
    body=message
)

# 绑定: kern.* 匹配 kern.info, kern.critical
# 绑定: kern.# 匹配 kern.info.test 等

# 3. Fanout Exchange - 广播
channel.exchange_declare(exchange='notifications', exchange_type='fanout')
channel.basic_publish(
    exchange='notifications',
    routing_key='',  # 忽略routing key
    body=message
)

# 4. Headers Exchange - 头属性匹配
channel.exchange_declare(exchange='headers_exchange', exchange_type='headers')
channel.basic_publish(
    exchange='headers_exchange',
    routing_key='',
    properties=pika.BasicProperties(
        headers={'format': 'pdf', 'type': 'report'}
    ),
    body=message
)
```

### 4.3 死信队列模式

```python
# 声明主队列和死信队列
args = {
    'x-message-ttl': 30000,  # 30秒TTL
    'x-dead-letter-exchange': 'dlx',  # 死信交换机
    'x-dead-letter-routing-key': 'failed',
    'x-max-retries': 3  # 最大重试次数
}

channel.queue_declare(queue='main_queue', arguments=args)
channel.queue_declare(queue='dlq')  # 死信队列
channel.exchange_declare(exchange='dlx', exchange_type='direct')
channel.queue_bind(queue='dlq', exchange='dlx', routing_key='failed')

# 消费消息
channel.basic_consume(queue='main_queue', on_message_callback=process_message)

def process_message(ch, method, properties, body):
    try:
        # 处理消息
        result = do_work(body)
        ch.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as e:
        # 拒绝消息，进入死信队列
        ch.basic_reject(
            delivery_tag=method.delivery_tag,
            requeue=False  # False = 进入DLQ
        )
```

---

## 5. NATS

### 5.1 架构

```
┌─────────────────────────────────────────┐
│           NATS Architecture             │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────┐     ┌─────────────┐       │
│  │Publisher│────►│ NATS Server │       │
│  └─────────┘     │  (Core NATS)│       │
│                  └──────┬──────┘       │
│                         │              │
│           ┌─────────────┼─────────────┐│
│           │             │             ││
│      ┌────┴────┐   ┌────┴────┐   ┌────┴────┐
│      │Subscriber1│  │Subscriber2│  │Subscriber3│
│      └─────────┘   └─────────┘   └─────────┘
│                                         │
│  JetStream: 持久化流                    │
│  Key-Value: 分布式KV存储                │
│  Object Store: 对象存储                 │
│                                         │
└─────────────────────────────────────────┘
```

### 5.2 JetStream

```go
package main

import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
)

func main() {
    nc, _ := nats.Connect(nats.DefaultURL)
    js, _ := jetstream.New(nc)

    ctx := context.Background()

    // 创建流
    stream, _ := js.CreateStream(ctx, jetstream.StreamConfig{
        Name:     "ORDERS",
        Subjects: []string{"orders.>"},
        Retention: jetstream.WorkQueuePolicy,
        MaxMsgs:  1000000,
        Storage:  jetstream.FileStorage,
        Replicas: 3,
    })

    // 发布消息
    ack, _ := js.Publish(ctx, "orders.created", []byte(`{"id": 1}`))
    fmt.Printf("Published: %s\n", ack.Stream)

    // 创建消费者
    cons, _ := stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
        Name:          "order-processor",
        Durable:       "order-processor",
        DeliverPolicy: jetstream.DeliverAllPolicy,
        AckPolicy:     jetstream.AckExplicitPolicy,
        MaxDeliver:    3,
    })

    // 消费消息
    cons.Consume(func(msg jetstream.Msg) {
        fmt.Printf("Received: %s\n", msg.Data())
        msg.Ack()
    })
}
```

### 5.3 请求-响应模式

```go
// 服务端
nc.Subscribe("help.request", func(m *nats.Msg) {
    nc.Publish(m.Reply, []byte("I can help!"))
})

// 客户端
msg, _ := nc.Request("help.request", []byte("help"), 2*time.Second)
fmt.Printf("Response: %s\n", msg.Data)
```

---

## 6. Redis Streams

### 6.1 基础操作

```bash
# 添加消息到流
XADD mystream * sensor-id 1234 temperature 19.8

# 读取消息
XREAD COUNT 2 STREAMS mystream 0

# 消费者组
XGROUP CREATE mystream mygroup $ MKSTREAM

# 消费者读取
XREADGROUP GROUP mygroup consumer1 COUNT 1 STREAMS mystream >

# 确认消息
XACK mystream mygroup 1526569495631-0

# 查看待处理消息
XPENDING mystream mygroup

# 范围查询
XRANGE mystream - +
XREVRANGE mystream + - COUNT 10
```

### 6.2 Go示例

```go
package main

import (
    "github.com/redis/go-redis/v9"
)

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    ctx := context.Background()

    // 添加消息
    id, err := rdb.XAdd(ctx, &redis.XAddArgs{
        Stream: "mystream",
        Values: map[string]interface{}{
            "sensor-id":   "1234",
            "temperature": 19.8,
        },
    }).Result()

    // 创建消费者组
    rdb.XGroupCreate(ctx, "mystream", "mygroup", "0")

    // 消费消息
    msgs, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
        Group:    "mygroup",
        Consumer: "consumer1",
        Streams:  []string{"mystream", ">"},
        Count:    10,
        Block:    0,
    }).Result()

    for _, msg := range msgs[0].Messages {
        // 处理消息
        fmt.Printf("ID: %s, Values: %v\n", msg.ID, msg.Values)

        // 确认
        rdb.XAck(ctx, "mystream", "mygroup", msg.ID)
    }
}
```

---

## 7. 消息队列选型对比

### 7.1 性能对比

| 特性 | Kafka | Pulsar | RabbitMQ | NATS | Redis Streams |
|------|-------|--------|----------|------|---------------|
| 吞吐量 | 极高(1M+/s) | 极高(1M+/s) | 高(100K/s) | 极高 | 高(500K/s) |
| 延迟 | 10-100ms | 5-50ms | 1-10ms | 亚毫秒 | 亚毫秒 |
| 持久化 | 是 | 是 | 是 | 可选 | 是 |
| 复制 | 是 | 是 | 是 | 是 | 是 |
| 多租户 | 否 | 是 | 是 | 是 | 否 |
| 地理复制 | MirrorMaker | 原生 | Shovel | 网关 | Redis Cluster |
| 流处理 | Kafka Streams | Pulsar Functions | 有限 | 有限 | 有限 |

### 7.2 场景推荐

| 场景 | 推荐 |
|------|------|
| 大数据管道 | Kafka |
| 多租户SaaS | Pulsar |
| 企业集成 | RabbitMQ |
| 微服务通信 | NATS |
| 实时分析 | Redis Streams |
| 物联网 | NATS/MQTT |
| 金融交易 | Pulsar/Kafka |

---

## 8. 模式与反模式

### 8.1 消息幂等性

```go
// 幂等消费者
func processMessage(msg Message) error {
    // 检查是否已处理
    if cache.IsProcessed(msg.ID) {
        return nil  // 已处理，跳过
    }

    // 业务处理
    if err := doWork(msg); err != nil {
        return err
    }

    // 标记已处理
    cache.MarkProcessed(msg.ID, time.Hour*24)
    return nil
}
```

### 8.2 死信处理

```go
// 带重试的处理器
func processWithRetry(msg Message, maxRetries int) error {
    retries := msg.Headers.GetInt("x-retries")

    if err := doWork(msg); err != nil {
        if retries < maxRetries {
            // 重新排队，延迟指数退避
            delay := time.Duration(math.Pow(2, float64(retries))) * time.Second
            msg.Headers.Set("x-retries", retries+1)
            scheduler.Schedule(msg, delay)
        } else {
            // 送入死信队列
            dlq.Publish(msg)
        }
        return err
    }
    return nil
}
```

---

## 9. 参考文献

1. Apache Kafka Documentation
2. Apache Pulsar Architecture
3. RabbitMQ Best Practices
4. NATS Documentation
5. Redis Streams Guide

---

*Last Updated: 2026-04-03*
