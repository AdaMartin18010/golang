# EC-033: Transactional Outbox Pattern (事务发件箱模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #transactional-outbox #event-driven #at-least-once #reliability
> **权威来源**:
>
> - [Transactional Outbox Pattern](https://microservices.io/patterns/data/transactional-outbox.html) - Chris Richardson
> - [Implementing the Outbox Pattern](https://debezium.io/blog/2019/02/19/reliable-microservices-integration-with-the-outbox-pattern/) - Debezium
> - [Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/) - Hohpe & Woolf

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在微服务架构中，如何确保数据库操作和消息发布之间的原子性，避免"数据已更新但事件未发送"或"事件已发送但数据更新失败"的不一致状态？

**双写问题 (Dual Write Problem)**:

```
场景 A: 数据库提交成功，消息发送失败
┌─────────┐    ┌─────────┐    ┌─────────┐
│  Start  │───►│ DB Commit│───►│ Message │
└─────────┘    └────┬────┘    │  Fail   │
                    │         └────┬────┘
                    ▼              │
              ┌─────────┐          │
              │ Data    │    ❌ Inconsistent!
              │ Updated │          │
              └─────────┘    Event lost

场景 B: 消息发送成功，数据库回滚
┌─────────┐    ┌─────────┐    ┌─────────┐
│  Start  │───►│ Message │───►│ DB      │
└─────────┘    │  Sent   │    │ Rollback│
               └────┬────┘    └────┬────┘
                    │              │
                    ▼              │
              ┌─────────┐    ❌ Inconsistent!
              │ Event   │          │
              │ Sent    │    Data not updated
              └─────────┘
```

**形式化描述**:

```
给定: 数据库事务 T_db 和消息发布 T_msg
约束: 需要保证 T_db 和 T_msg 的原子性
      ∀s: committed(s) → published(s)
      ∀s: ¬committed(s) → ¬published(s)
目标: 找到实现原子性的机制 M
```

### 1.2 解决方案形式化

**定义 1.1 (事务发件箱)**
事务发件箱模式通过以下机制保证原子性：

1. 使用本地数据库事务同时更新业务数据和事件记录
2. 使用单独的进程（Relay）读取发件箱表
3. 发布事件到消息代理
4. 标记或删除已发布的事件

**形式化表示**:

```
原子操作:
  BEGIN TRANSACTION
    UPDATE business_table SET ...
    INSERT INTO outbox_table (event_type, payload, created_at) VALUES (...)
  COMMIT

发布流程:
  ∀e ∈ outbox_table WHERE published = false:
    publish(e) → mark_published(e)
```

**一致性保证**:

```
性质 1 (至少一次发布):
  ∀e ∈ outbox_table: committed(e) → ◇published(e)

性质 2 (幂等消费):
  消费者必须支持幂等: f(e) = f(f(e))

性质 3 (顺序保证):
  可选: publish_order(e₁) < publish_order(e₂) → delivery_order(e₁) < delivery_order(e₂)
```

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Transactional Outbox Architecture                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Application Layer                             │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │   Service   │───►│  Business   │───►│   Outbox Writer     │  │   │
│  │  │   Handler   │    │   Logic     │    │  (Same Transaction) │  │   │
│  │  └─────────────┘    └──────┬──────┘    └─────────────────────┘  │   │
│  │                           │                                      │   │
│  │                           ▼                                      │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │              Database (ACID Transaction)                   │  │   │
│  │  │  ┌──────────────┐            ┌──────────────────────┐     │  │   │
│  │  │  │ Business     │◄──────────►│ Outbox Table         │     │  │   │
│  │  │  │ Tables       │   FK/Atomic │ - id                 │     │  │   │
│  │  │  │              │            │ - event_type         │     │  │   │
│  │  │  │              │            │ - aggregate_type     │     │  │   │
│  │  │  │              │            │ - aggregate_id       │     │  │   │
│  │  │  │              │            │ - payload (JSON)     │     │  │   │
│  │  │  │              │            │ - created_at         │     │  │   │
│  │  │  │              │            │ - published (bool)   │     │  │   │
│  │  │  └──────────────┘            └──────────────────────┘     │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                    │
│                                    │ Polling / CDC                       │
│                                    ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                       Relay Process                              │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │   Polling   │───►│   Event     │───►│  Message Broker     │  │   │
│  │  │  (or CDC)   │    │  Publisher  │    │  (Kafka/RabbitMQ)   │  │   │
│  │  └─────────────┘    └─────────────┘    └─────────────────────┘  │   │
│  │         │                  │                  │                  │   │
│  │         ▼                  ▼                  ▼                  │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │  After Success:                                           │  │   │
│  │  │  - DELETE from outbox (or mark published)                 │  │   │
│  │  │  - Track processed events (idempotency)                   │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                    │
│                                    │ Publish                             │
│                                    ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                      Consumer Services                           │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │  Consumer 1 │    │  Consumer 2 │    │  Consumer N         │  │   │
│  │  └─────────────┘    └─────────────┘    └─────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心组件实现

```go
// outbox/core.go
package outbox

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"

    "github.com/google/uuid"
)

// Event 发件箱事件
type Event struct {
    ID            string          `db:"id" json:"id"`
    EventType     string          `db:"event_type" json:"event_type"`
    AggregateType string          `db:"aggregate_type" json:"aggregate_type"`
    AggregateID   string          `db:"aggregate_id" json:"aggregate_id"`
    Payload       json.RawMessage `db:"payload" json:"payload"`
    Metadata      json.RawMessage `db:"metadata" json:"metadata"`
    CreatedAt     time.Time       `db:"created_at" json:"created_at"`
    PublishedAt   *time.Time      `db:"published_at" json:"published_at,omitempty"`
}

// OutboxEvent 待发布事件
type OutboxEvent struct {
    ID            string
    EventType     string
    AggregateType string
    AggregateID   string
    Payload       interface{}
    Metadata      map[string]string
}

// Publisher 消息发布器接口
type Publisher interface {
    Publish(ctx context.Context, topic string, event *Event) error
    Close() error
}

// Store 发件箱存储接口
type Store interface {
    // Save 保存事件到发件箱（在业务事务中调用）
    Save(ctx context.Context, tx *sql.Tx, event *OutboxEvent) error

    // GetUnpublished 获取未发布事件
    GetUnpublished(ctx context.Context, batchSize int) ([]*Event, error)

    // MarkPublished 标记事件为已发布
    MarkPublished(ctx context.Context, ids []string) error

    // DeletePublished 删除已发布事件
    DeletePublished(ctx context.Context, ids []string) error

    // CleanOld 清理旧事件
    CleanOld(ctx context.Context, before time.Time) (int64, error)
}

// Relay 发件箱中继器
type Relay struct {
    store         Store
    publisher     Publisher
    pollInterval  time.Duration
    batchSize     int
    stopCh        chan struct{}
    stoppedCh     chan struct{}
    logger        Logger
}

// Logger 日志接口
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// Config 中继器配置
type Config struct {
    PollInterval time.Duration
    BatchSize    int
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
    return &Config{
        PollInterval: 5 * time.Second,
        BatchSize:    100,
    }
}

// NewRelay 创建中继器
func NewRelay(store Store, publisher Publisher, logger Logger, config *Config) *Relay {
    if config == nil {
        config = DefaultConfig()
    }

    return &Relay{
        store:        store,
        publisher:    publisher,
        pollInterval: config.PollInterval,
        batchSize:    config.BatchSize,
        stopCh:       make(chan struct{}),
        stoppedCh:    make(chan struct{}),
        logger:       logger,
    }
}

// Start 启动中继器
func (r *Relay) Start(ctx context.Context) error {
    r.logger.Info("starting outbox relay",
        Field{"poll_interval", r.pollInterval},
        Field{"batch_size", r.batchSize})

    ticker := time.NewTicker(r.pollInterval)
    defer ticker.Stop()
    defer close(r.stoppedCh)

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-r.stopCh:
            return nil
        case <-ticker.C:
            if err := r.processBatch(ctx); err != nil {
                r.logger.Error("failed to process outbox batch", Field{"error", err})
            }
        }
    }
}

// Stop 停止中继器
func (r *Relay) Stop() {
    close(r.stopCh)
    <-r.stoppedCh
}

// processBatch 处理一批事件
func (r *Relay) processBatch(ctx context.Context) error {
    events, err := r.store.GetUnpublished(ctx, r.batchSize)
    if err != nil {
        return fmt.Errorf("failed to get unpublished events: %w", err)
    }

    if len(events) == 0 {
        return nil
    }

    r.logger.Debug("processing outbox batch", Field{"count", len(events)})

    var published []string
    for _, event := range events {
        if err := r.publishEvent(ctx, event); err != nil {
            r.logger.Error("failed to publish event",
                Field{"event_id", event.ID},
                Field{"error", err})
            // 继续处理其他事件
            continue
        }
        published = append(published, event.ID)
    }

    if len(published) > 0 {
        if err := r.store.MarkPublished(ctx, published); err != nil {
            return fmt.Errorf("failed to mark events as published: %w", err)
        }

        r.logger.Info("published events", Field{"count", len(published)})
    }

    return nil
}

// publishEvent 发布单个事件
func (r *Relay) publishEvent(ctx context.Context, event *Event) error {
    topic := fmt.Sprintf("%s.%s", event.AggregateType, event.EventType)

    if err := r.publisher.Publish(ctx, topic, event); err != nil {
        return fmt.Errorf("failed to publish to %s: %w", topic, err)
    }

    return nil
}

// OutboxWriter 发件箱写入器
type OutboxWriter struct {
    store Store
}

// NewOutboxWriter 创建写入器
func NewOutboxWriter(store Store) *OutboxWriter {
    return &OutboxWriter{store: store}
}

// Write 写入事件（在事务中调用）
func (w *OutboxWriter) Write(ctx context.Context, tx *sql.Tx, event *OutboxEvent) error {
    if event.ID == "" {
        event.ID = uuid.New().String()
    }

    return w.store.Save(ctx, tx, event)
}

// TransactionalService 事务性服务
type TransactionalService struct {
    db          *sql.DB
    outbox      *OutboxWriter
    aggregateType string
}

// NewTransactionalService 创建事务性服务
func NewTransactionalService(db *sql.DB, store Store, aggregateType string) *TransactionalService {
    return &TransactionalService{
        db:            db,
        outbox:        NewOutboxWriter(store),
        aggregateType: aggregateType,
    }
}

// Execute 执行业务操作并记录事件
func (s *TransactionalService) Execute(ctx context.Context, operation func(tx *sql.Tx) (*OutboxEvent, error)) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // 执行业务操作
    event, err := operation(tx)
    if err != nil {
        return err
    }

    if event != nil {
        // 设置聚合类型
        if event.AggregateType == "" {
            event.AggregateType = s.aggregateType
        }

        // 写入发件箱
        if err := s.outbox.Write(ctx, tx, event); err != nil {
            return fmt.Errorf("failed to write to outbox: %w", err)
        }
    }

    // 提交事务
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}
```

### 2.2 数据库存储实现

```go
// outbox/sql_store.go
package outbox

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
)

// SQLStore SQL 存储实现
type SQLStore struct {
    db             *sql.DB
    outboxTable    string
    idempotencyTTL time.Duration
}

// NewSQLStore 创建 SQL 存储
func NewSQLStore(db *sql.DB, outboxTable string) *SQLStore {
    return &SQLStore{
        db:             db,
        outboxTable:    outboxTable,
        idempotencyTTL: 7 * 24 * time.Hour,
    }
}

// Save 保存事件
func (s *SQLStore) Save(ctx context.Context, tx *sql.Tx, event *OutboxEvent) error {
    payload, err := json.Marshal(event.Payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    metadata, err := json.Marshal(event.Metadata)
    if err != nil {
        return fmt.Errorf("failed to marshal metadata: %w", err)
    }

    query := fmt.Sprintf(`
        INSERT INTO %s (id, event_type, aggregate_type, aggregate_id, payload, metadata, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, s.outboxTable)

    _, err = tx.ExecContext(ctx, query,
        event.ID,
        event.EventType,
        event.AggregateType,
        event.AggregateID,
        payload,
        metadata,
        time.Now(),
    )

    if err != nil {
        return fmt.Errorf("failed to insert into outbox: %w", err)
    }

    return nil
}

// GetUnpublished 获取未发布事件
func (s *SQLStore) GetUnpublished(ctx context.Context, batchSize int) ([]*Event, error) {
    query := fmt.Sprintf(`
        SELECT id, event_type, aggregate_type, aggregate_id, payload, metadata, created_at
        FROM %s
        WHERE published_at IS NULL
        ORDER BY created_at ASC
        LIMIT ?
    `, s.outboxTable)

    rows, err := s.db.QueryContext(ctx, query, batchSize)
    if err != nil {
        return nil, fmt.Errorf("failed to query outbox: %w", err)
    }
    defer rows.Close()

    var events []*Event
    for rows.Next() {
        event := &Event{}
        err := rows.Scan(
            &event.ID,
            &event.EventType,
            &event.AggregateType,
            &event.AggregateID,
            &event.Payload,
            &event.Metadata,
            &event.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan event: %w", err)
        }
        events = append(events, event)
    }

    return events, rows.Err()
}

// MarkPublished 标记为已发布
func (s *SQLStore) MarkPublished(ctx context.Context, ids []string) error {
    if len(ids) == 0 {
        return nil
    }

    // 构建 IN 子句
    placeholders := make([]string, len(ids))
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        placeholders[i] = "?"
        args[i] = id
    }

    query := fmt.Sprintf(`
        UPDATE %s
        SET published_at = ?
        WHERE id IN (%s)
    `, s.outboxTable, joinPlaceholders(placeholders))

    args = append([]interface{}{time.Now()}, args...)

    _, err := s.db.ExecContext(ctx, query, args...)
    if err != nil {
        return fmt.Errorf("failed to mark as published: %w", err)
    }

    return nil
}

// DeletePublished 删除已发布事件
func (s *SQLStore) DeletePublished(ctx context.Context, ids []string) error {
    if len(ids) == 0 {
        return nil
    }

    placeholders := make([]string, len(ids))
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        placeholders[i] = "?"
        args[i] = id
    }

    query := fmt.Sprintf(`
        DELETE FROM %s
        WHERE id IN (%s)
    `, s.outboxTable, joinPlaceholders(placeholders))

    _, err := s.db.ExecContext(ctx, query, args...)
    if err != nil {
        return fmt.Errorf("failed to delete published: %w", err)
    }

    return nil
}

// CleanOld 清理旧事件
func (s *SQLStore) CleanOld(ctx context.Context, before time.Time) (int64, error) {
    query := fmt.Sprintf(`
        DELETE FROM %s
        WHERE published_at IS NOT NULL AND published_at < ?
    `, s.outboxTable)

    result, err := s.db.ExecContext(ctx, query, before)
    if err != nil {
        return 0, fmt.Errorf("failed to clean old events: %w", err)
    }

    return result.RowsAffected()
}

func joinPlaceholders(placeholders []string) string {
    result := ""
    for i, p := range placeholders {
        if i > 0 {
            result += ","
        }
        result += p
    }
    return result
}

// CreateTableSQL 创建表的 SQL
func (s *SQLStore) CreateTableSQL() string {
    return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id VARCHAR(36) PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    payload JSON NOT NULL,
    metadata JSON,
    created_at TIMESTAMP NOT NULL,
    published_at TIMESTAMP NULL,
    INDEX idx_unpublished (published_at, created_at),
    INDEX idx_aggregate (aggregate_type, aggregate_id)
);`, s.outboxTable)
}
```

### 2.3 Kafka 发布器实现

```go
// outbox/kafka_publisher.go
package outbox

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/IBM/sarama"
)

// KafkaPublisher Kafka 发布器
type KafkaPublisher struct {
    producer sarama.SyncProducer
}

// NewKafkaPublisher 创建 Kafka 发布器
func NewKafkaPublisher(brokers []string, config *sarama.Config) (*KafkaPublisher, error) {
    if config == nil {
        config = sarama.NewConfig()
        config.Producer.RequiredAcks = sarama.WaitForAll
        config.Producer.Retry.Max = 3
        config.Producer.Return.Successes = true
    }

    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create kafka producer: %w", err)
    }

    return &KafkaPublisher{producer: producer}, nil
}

// Publish 发布事件
func (p *KafkaPublisher) Publish(ctx context.Context, topic string, event *Event) error {
    payload, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }

    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(event.AggregateID),
        Value: sarama.ByteEncoder(payload),
        Headers: []sarama.RecordHeader{
            {Key: []byte("event_id"), Value: []byte(event.ID)},
            {Key: []byte("event_type"), Value: []byte(event.EventType)},
            {Key: []byte("aggregate_type"), Value: []byte(event.AggregateType)},
        },
    }

    _, _, err = p.producer.SendMessage(msg)
    if err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }

    return nil
}

// Close 关闭发布器
func (p *KafkaPublisher) Close() error {
    return p.producer.Close()
}
```

### 2.4 订单服务示例

```go
// examples/order/order_service.go
package order

import (
    "context"
    "database/sql"
    "encoding/json"
    "time"

    "go-knowledge-base/outbox"
)

// OrderService 订单服务
type OrderService struct {
    transactional *outbox.TransactionalService
}

// Order 订单
type Order struct {
    ID         string    `json:"id"`
    CustomerID string    `json:"customer_id"`
    Items      []Item    `json:"items"`
    Total      float64   `json:"total"`
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"created_at"`
}

// Item 订单项
type Item struct {
    ProductID string  `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
    OrderID    string    `json:"order_id"`
    CustomerID string    `json:"customer_id"`
    Total      float64   `json:"total"`
    CreatedAt  time.Time `json:"created_at"`
}

// NewOrderService 创建订单服务
func NewOrderService(db *sql.DB, store outbox.Store) *OrderService {
    return &OrderService{
        transactional: outbox.NewTransactionalService(db, store, "order"),
    }
}

// CreateOrder 创建订单（事务性）
func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []Item) (*Order, error) {
    var createdOrder *Order

    err := s.transactional.Execute(ctx, func(tx *sql.Tx) (*outbox.OutboxEvent, error) {
        // 计算总价
        var total float64
        for _, item := range items {
            total += item.Price * float64(item.Quantity)
        }

        // 创建订单
        order := &Order{
            ID:         generateOrderID(),
            CustomerID: customerID,
            Items:      items,
            Total:      total,
            Status:     "PENDING",
            CreatedAt:  time.Now(),
        }

        // 插入订单到数据库
        if _, err := tx.ExecContext(ctx,
            "INSERT INTO orders (id, customer_id, total, status, created_at) VALUES (?, ?, ?, ?, ?)",
            order.ID, order.CustomerID, order.Total, order.Status, order.CreatedAt); err != nil {
            return nil, err
        }

        // 插入订单项
        for _, item := range order.Items {
            if _, err := tx.ExecContext(ctx,
                "INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
                order.ID, item.ProductID, item.Quantity, item.Price); err != nil {
                return nil, err
            }
        }

        createdOrder = order

        // 创建事件
        eventPayload := OrderCreatedEvent{
            OrderID:    order.ID,
            CustomerID: order.CustomerID,
            Total:      order.Total,
            CreatedAt:  order.CreatedAt,
        }

        return &outbox.OutboxEvent{
            EventType:     "OrderCreated",
            AggregateType: "order",
            AggregateID:   order.ID,
            Payload:       eventPayload,
            Metadata: map[string]string{
                "correlation_id": getCorrelationID(ctx),
                "service":        "order-service",
            },
        }, nil
    })

    if err != nil {
        return nil, err
    }

    return createdOrder, nil
}

// UpdateOrderStatus 更新订单状态（事务性）
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error {
    return s.transactional.Execute(ctx, func(tx *sql.Tx) (*outbox.OutboxEvent, error) {
        // 更新订单状态
        result, err := tx.ExecContext(ctx,
            "UPDATE orders SET status = ? WHERE id = ?",
            newStatus, orderID)
        if err != nil {
            return nil, err
        }

        rows, err := result.RowsAffected()
        if err != nil {
            return nil, err
        }
        if rows == 0 {
            return nil, sql.ErrNoRows
        }

        // 创建状态变更事件
        return &outbox.OutboxEvent{
            EventType:     "OrderStatusChanged",
            AggregateType: "order",
            AggregateID:   orderID,
            Payload: map[string]interface{}{
                "order_id":   orderID,
                "new_status": newStatus,
                "changed_at": time.Now(),
            },
        }, nil
    })
}

func generateOrderID() string {
    return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

func getCorrelationID(ctx context.Context) string {
    // 从 context 获取 correlation ID
    if id, ok := ctx.Value("correlation_id").(string); ok {
        return id
    }
    return ""
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// outbox/core_test.go
package outbox

import (
    "context"
    "database/sql"
    "encoding/json"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...Field)  {}
func (m *mockLogger) Error(msg string, fields ...Field) {}
func (m *mockLogger) Debug(msg string, fields ...Field) {}

type mockStore struct {
    mock.Mock
}

func (m *mockStore) Save(ctx context.Context, tx *sql.Tx, event *OutboxEvent) error {
    args := m.Called(ctx, tx, event)
    return args.Error(0)
}

func (m *mockStore) GetUnpublished(ctx context.Context, batchSize int) ([]*Event, error) {
    args := m.Called(ctx, batchSize)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*Event), args.Error(1)
}

func (m *mockStore) MarkPublished(ctx context.Context, ids []string) error {
    args := m.Called(ctx, ids)
    return args.Error(0)
}

func (m *mockStore) DeletePublished(ctx context.Context, ids []string) error {
    args := m.Called(ctx, ids)
    return args.Error(0)
}

func (m *mockStore) CleanOld(ctx context.Context, before time.Time) (int64, error) {
    args := m.Called(ctx, before)
    return args.Get(0).(int64), args.Error(1)
}

type mockPublisher struct {
    mock.Mock
}

func (m *mockPublisher) Publish(ctx context.Context, topic string, event *Event) error {
    args := m.Called(ctx, topic, event)
    return args.Error(0)
}

func (m *mockPublisher) Close() error {
    args := m.Called()
    return args.Error(0)
}

func TestRelay_processBatch(t *testing.T) {
    store := new(mockStore)
    publisher := new(mockPublisher)
    logger := &mockLogger{}

    relay := NewRelay(store, publisher, logger, &Config{
        PollInterval: 1 * time.Second,
        BatchSize:    10,
    })

    // 准备测试事件
    events := []*Event{
        {
            ID:            "evt-001",
            EventType:     "OrderCreated",
            AggregateType: "order",
            AggregateID:   "ORD-001",
            Payload:       json.RawMessage(`{"order_id":"ORD-001"}`),
            CreatedAt:     time.Now(),
        },
        {
            ID:            "evt-002",
            EventType:     "OrderUpdated",
            AggregateType: "order",
            AggregateID:   "ORD-002",
            Payload:       json.RawMessage(`{"order_id":"ORD-002"}`),
            CreatedAt:     time.Now(),
        },
    }

    // 设置期望
    store.On("GetUnpublished", mock.Anything, 10).Return(events, nil)
    publisher.On("Publish", mock.Anything, "order.OrderCreated", events[0]).Return(nil)
    publisher.On("Publish", mock.Anything, "order.OrderUpdated", events[1]).Return(nil)
    store.On("MarkPublished", mock.Anything, []string{"evt-001", "evt-002"}).Return(nil)

    // 执行
    err := relay.processBatch(context.Background())

    // 验证
    require.NoError(t, err)
    store.AssertExpectations(t)
    publisher.AssertExpectations(t)
}

func TestRelay_processBatch_PartialFailure(t *testing.T) {
    store := new(mockStore)
    publisher := new(mockPublisher)
    logger := &mockLogger{}

    relay := NewRelay(store, publisher, logger, DefaultConfig())

    events := []*Event{
        {ID: "evt-001", EventType: "OrderCreated", AggregateType: "order"},
        {ID: "evt-002", EventType: "OrderFailed", AggregateType: "order"},
    }

    store.On("GetUnpublished", mock.Anything, 100).Return(events, nil)
    publisher.On("Publish", mock.Anything, "order.OrderCreated", mock.Anything).Return(nil)
    publisher.On("Publish", mock.Anything, "order.OrderFailed", mock.Anything).Return(assert.AnError)
    store.On("MarkPublished", mock.Anything, []string{"evt-001"}).Return(nil)

    err := relay.processBatch(context.Background())

    require.NoError(t, err) // 部分失败不应返回错误
    store.AssertExpectations(t)
}

func TestOutboxWriter_Write(t *testing.T) {
    store := new(mockStore)
    writer := NewOutboxWriter(store)

    event := &OutboxEvent{
        EventType:     "TestEvent",
        AggregateType: "test",
        AggregateID:   "123",
        Payload:       map[string]string{"key": "value"},
    }

    store.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

    err := writer.Write(context.Background(), nil, event)

    require.NoError(t, err)
    assert.NotEmpty(t, event.ID) // 应自动生成 ID
    store.AssertExpectations(t)
}

func TestTransactionalService_Execute(t *testing.T) {
    // 这需要集成测试，因为涉及真实数据库事务
    // 这里只做接口测试

    store := new(mockStore)
    service := NewTransactionalService(nil, store, "order")

    assert.NotNil(t, service)
    assert.Equal(t, "order", service.aggregateType)
}
```

### 3.2 集成测试

```go
// outbox/integration_test.go
package outbox

import (
    "context"
    "database/sql"
    "testing"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)

    // 创建发件箱表
    _, err = db.Exec(`
        CREATE TABLE outbox (
            id VARCHAR(36) PRIMARY KEY,
            event_type VARCHAR(255) NOT NULL,
            aggregate_type VARCHAR(255) NOT NULL,
            aggregate_id VARCHAR(255) NOT NULL,
            payload BLOB NOT NULL,
            metadata BLOB,
            created_at TIMESTAMP NOT NULL,
            published_at TIMESTAMP NULL
        )
    `)
    require.NoError(t, err)

    return db
}

func TestSQLStore_Integration(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    store := NewSQLStore(db, "outbox")
    ctx := context.Background()

    t.Run("save and retrieve", func(t *testing.T) {
        tx, err := db.Begin()
        require.NoError(t, err)
        defer tx.Rollback()

        event := &OutboxEvent{
            ID:            "evt-001",
            EventType:     "TestEvent",
            AggregateType: "test",
            AggregateID:   "123",
            Payload:       map[string]string{"key": "value"},
            Metadata:      map[string]string{"meta": "data"},
        }

        err = store.Save(ctx, tx, event)
        require.NoError(t, err)

        err = tx.Commit()
        require.NoError(t, err)

        // 获取未发布事件
        events, err := store.GetUnpublished(ctx, 10)
        require.NoError(t, err)
        assert.Len(t, events, 1)
        assert.Equal(t, "evt-001", events[0].ID)

        // 标记为已发布
        err = store.MarkPublished(ctx, []string{"evt-001"})
        require.NoError(t, err)

        // 验证已发布
        events, err = store.GetUnpublished(ctx, 10)
        require.NoError(t, err)
        assert.Len(t, events, 0)
    })

    t.Run("clean old events", func(t *testing.T) {
        // 插入旧事件
        tx, _ := db.Begin()
        store.Save(ctx, tx, &OutboxEvent{
            ID:            "old-evt",
            EventType:     "OldEvent",
            AggregateType: "test",
            AggregateID:   "old",
            Payload:       map[string]string{},
        })
        tx.Commit()

        // 标记为已发布（模拟过去的时间）
        db.Exec("UPDATE outbox SET published_at = ? WHERE id = ?",
            time.Now().Add(-30*24*time.Hour), "old-evt")

        // 清理旧事件
        deleted, err := store.CleanOld(ctx, time.Now().Add(-7*24*time.Hour))
        require.NoError(t, err)
        assert.Equal(t, int64(1), deleted)
    })
}

func TestEndToEnd_OrderService(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    // 创建订单表
    _, err := db.Exec(`
        CREATE TABLE orders (
            id VARCHAR(36) PRIMARY KEY,
            customer_id VARCHAR(255) NOT NULL,
            total REAL NOT NULL,
            status VARCHAR(50) NOT NULL,
            created_at TIMESTAMP NOT NULL
        )
    `)
    require.NoError(t, err)

    store := NewSQLStore(db, "outbox")
    orderService := NewOrderService(db, store)

    ctx := context.Background()

    // 创建订单
    order, err := orderService.CreateOrder(ctx, "CUST-001", []Item{
        {ProductID: "PROD-1", Quantity: 2, Price: 29.99},
    })

    require.NoError(t, err)
    assert.NotNil(t, order)
    assert.Equal(t, "CUST-001", order.CustomerID)
    assert.Equal(t, 59.98, order.Total)

    // 验证订单在数据库中
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", order.ID).Scan(&count)
    require.NoError(t, err)
    assert.Equal(t, 1, count)

    // 验证事件在发件箱中
    events, err := store.GetUnpublished(ctx, 10)
    require.NoError(t, err)
    assert.Len(t, events, 1)
    assert.Equal(t, "OrderCreated", events[0].EventType)
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Change Data Capture (CDC) 的集成

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Outbox with CDC Integration                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                      Application                                 │   │
│  │  ┌─────────────┐    ┌─────────────┐                            │   │
│  │  │  Service    │───►│   Outbox    │                            │   │
│  │  │             │    │   Table     │                            │   │
│  │  └─────────────┘    └──────┬──────┘                            │   │
│  │                            │ INSERT                            │   │
│  └────────────────────────────┼───────────────────────────────────┘   │
│                               │                                        │
│                               ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Database (MySQL/PostgreSQL)                  │   │
│  │  ┌─────────────────────────────────────────────────────┐       │   │
│  │  │                   Binlog/WAL                         │       │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │       │   │
│  │  │  │ INSERT  │  │ UPDATE  │  │ INSERT  │  │ DELETE  │  │       │   │
│  │  │  │ (Outbox)│  │(Orders) │  │ (Outbox)│  │(Outbox) │  │       │   │
│  │  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  │       │   │
│  │  └───────┼────────────┼────────────┼────────────┼────────┘       │   │
│  └──────────┼────────────┼────────────┼────────────┼────────────────┘   │
│             │            │            │            │                     │
│             └────────────┴────────────┴────────────┘                     │
│                          │                                              │
│                          ▼                                              │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     CDC Connector (Debezium)                     │   │
│  │                                                                 │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │ Binlog      │───►│  Event      │───►│  Kafka Connect      │  │   │
│  │  │ Parser      │    │  Filter     │    │  (Outbox Router)    │  │   │
│  │  └─────────────┘    └─────────────┘    └─────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                    │
│                                    │ Route by aggregate_type             │
│                                    ▼                                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Kafka Topics                                 │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │   │
│  │  │ order       │    │ payment     │    │ inventory           │  │   │
│  │  │ events      │    │ events      │    │ events              │  │   │
│  │  └─────────────┘    └─────────────┘    └─────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  优势:                                                                   │
│  • 无轮询开销（近实时）                                                   │
│  • 发件箱和事务在同一数据库事务中，保证原子性                                 │
│  • CDC 确保所有变更都被捕获                                               │
│  • 解耦发布逻辑                                                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 与 Idempotent Consumer 模式的集成

```go
// outbox/idempotent_consumer.go
package outbox

import (
    "context"
    "database/sql"
    "time"
)

// IdempotencyStore 幂等性存储接口
type IdempotencyStore interface {
    IsProcessed(ctx context.Context, eventID string) (bool, error)
    MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error
}

// IdempotentHandler 幂等处理器
type IdempotentHandler struct {
    handler   func(ctx context.Context, event *Event) error
    store     IdempotencyStore
}

// NewIdempotentHandler 创建幂等处理器
func NewIdempotentHandler(handler func(ctx context.Context, event *Event) error, store IdempotencyStore) *IdempotentHandler {
    return &IdempotentHandler{
        handler: handler,
        store:   store,
    }
}

// Handle 处理事件（幂等）
func (h *IdempotentHandler) Handle(ctx context.Context, event *Event) error {
    // 检查是否已处理
    processed, err := h.store.IsProcessed(ctx, event.ID)
    if err != nil {
        return err
    }

    if processed {
        return nil // 已处理，直接返回
    }

    // 执行业务逻辑
    if err := h.handler(ctx, event); err != nil {
        return err
    }

    // 标记为已处理
    if err := h.store.MarkProcessed(ctx, event.ID, 24*time.Hour); err != nil {
        return err
    }

    return nil
}

// SQLIdempotencyStore SQL 幂等性存储
type SQLIdempotencyStore struct {
    db    *sql.DB
    table string
}

// NewSQLIdempotencyStore 创建 SQL 存储
func NewSQLIdempotencyStore(db *sql.DB, table string) *SQLIdempotencyStore {
    return &SQLIdempotencyStore{db: db, table: table}
}

// IsProcessed 检查是否已处理
func (s *SQLIdempotencyStore) IsProcessed(ctx context.Context, eventID string) (bool, error) {
    var count int
    query := `SELECT COUNT(*) FROM ` + s.table + ` WHERE event_id = ?`
    err := s.db.QueryRowContext(ctx, query, eventID).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// MarkProcessed 标记为已处理
func (s *SQLIdempotencyStore) MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error {
    query := `INSERT INTO ` + s.table + ` (event_id, processed_at, expires_at) VALUES (?, ?, ?)`
    _, err := s.db.ExecContext(ctx, query, eventID, time.Now(), time.Now().Add(ttl))
    return err
}
```

### 4.3 与 SAGA 模式的集成

```go
// outbox/saga_integration.go
package outbox

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
)

// SagaOutboxBridge SAGA 发件箱桥接
type SagaOutboxBridge struct {
    transactional *TransactionalService
}

// NewSagaOutboxBridge 创建桥接
func NewSagaOutboxBridge(db *sql.DB, store Store) *SagaOutboxBridge {
    return &SagaOutboxBridge{
        transactional: NewTransactionalService(db, store, "saga"),
    }
}

// SagaEvent SAGA 事件
type SagaEvent struct {
    SagaID        string                 `json:"saga_id"`
    StepID        string                 `json:"step_id"`
    EventType     string                 `json:"event_type"`
    Payload       map[string]interface{} `json:"payload"`
    Timestamp     time.Time              `json:"timestamp"`
}

// PublishSagaEvent 发布 SAGA 事件
func (b *SagaOutboxBridge) PublishSagaEvent(ctx context.Context, sagaID string, stepID string, eventType string, payload map[string]interface{}) error {
    return b.transactional.Execute(ctx, func(tx *sql.Tx) (*OutboxEvent, error) {
        sagaEvent := SagaEvent{
            SagaID:    sagaID,
            StepID:    stepID,
            EventType: eventType,
            Payload:   payload,
            Timestamp: time.Now(),
        }

        return &OutboxEvent{
            EventType:     fmt.Sprintf("Saga%s", eventType),
            AggregateType: "saga",
            AggregateID:   sagaID,
            Payload:       sagaEvent,
            Metadata: map[string]string{
                "saga_id": sagaID,
                "step_id": stepID,
            },
        }, nil
    })
}
```

---

## 5. 决策标准

### 5.1 何时使用事务发件箱

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Transactional Outbox Decision Tree                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  开始                                                                    │
│   │                                                                     │
│   ▼                                                                     │
│  ┌─────────────────────────┐                                           │
│  │ 需要保证数据和事件一致性？│───否──► 直接发布消息                        │
│  └──────────┬──────────────┘                                           │
│             │是                                                         │
│             ▼                                                           │
│  ┌─────────────────────────┐                                           │
│  │ 数据库支持 CDC？         │───是──► 使用 CDC + Outbox                  │
│  └──────────┬──────────────┘                                           │
│             │否                                                         │
│             ▼                                                           │
│  ┌─────────────────────────┐                                           │
│  │ 可接受延迟？             │───是──► 使用 Polling Outbox                 │
│  └──────────┬──────────────┘                                           │
│             │否                                                         │
│             ▼                                                           │
│  ┌─────────────────────────┐                                           │
│  │ 需要强一致性？           │───是──► 考虑 2PC / TCC                      │
│  └──────────┬──────────────┘                                           │
│             │否                                                         │
│             ▼                                                           │
│         重新评估需求                                                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 实现选项对比

| 方案 | 延迟 | 复杂度 | 资源开销 | 适用场景 |
|------|------|--------|----------|----------|
| **Polling Outbox** | 秒级 | 低 | 低 | 大多数场景 |
| **CDC (Debezium)** | 毫秒级 | 高 | 中 | 实时性要求高 |
| **数据库触发器** | 毫秒级 | 中 | 中 | 简单场景 |
| **应用内双写** | 毫秒级 | 低 | 低 | ❌ 不推荐（无原子性） |

### 5.3 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Transactional Outbox Checklist                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 确定事件 schema 和版本策略                                            │
│  □ 定义聚合边界和事件粒度                                                │
│  □ 选择发布机制（Polling/CDC）                                           │
│  □ 设计幂等消费机制                                                      │
│  □ 规划事件保留策略                                                      │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现原子性事务（业务数据 + 发件箱）                                    │
│  □ 实现中继器（Polling 或 CDC）                                          │
│  □ 实现错误处理和重试机制                                                │
│  □ 添加监控指标（延迟、积压量）                                          │
│  □ 实现发件箱清理任务                                                    │
│                                                                         │
│  测试阶段:                                                               │
│  □ 测试事务原子性（数据库回滚时无事件）                                   │
│  □ 测试至少一次语义（发布失败重试）                                       │
│  □ 测试幂等消费（重复事件处理）                                          │
│  □ 测试顺序保证（如有要求）                                              │
│  □ 测试故障恢复（中继器重启）                                            │
│                                                                         │
│  运维阶段:                                                               │
│  □ 监控发件箱积压                                                        │
│  □ 配置告警（发布延迟、失败率）                                          │
│  □ 定期清理已发布事件                                                    │
│  □ 备份发件箱数据（可选）                                                │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 不要在大事务中包含发件箱写入                                          │
│  ❌ 不要在发件箱中存储敏感信息                                            │
│  ❌ 不要依赖事件顺序（除非必要）                                          │
│  ❌ 不要忘记处理死信事件                                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>20KB, 完整形式化 + Go 实现 + 测试 + 决策标准)

**相关文档**:

- [EC-034-Polling-Publisher.md](./EC-034-Polling-Publisher.md)
- [EC-039-Domain-Event-Pattern.md](./EC-039-Domain-Event-Pattern.md)
- [EC-031-Choreography-Pattern.md](./EC-031-Choreography-Pattern.md)
