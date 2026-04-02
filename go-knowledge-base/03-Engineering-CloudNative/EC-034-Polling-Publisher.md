# EC-034: Polling Publisher Pattern (轮询发布者模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #polling-publisher #event-driven #reliability #outbox
> **权威来源**:
>
> - [Polling Publisher Pattern](https://microservices.io/patterns/data/polling-publisher.html) - Chris Richardson
> - [Enterprise Integration Patterns](https://www.enterpriseintegrationpatterns.com/) - Hohpe & Woolf
> - [Outbox Pattern Implementation](https://debezium.io/blog/2019/02/19/reliable-microservices-integration-with-the-outbox-pattern/) - Debezium

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在事务发件箱模式中，如何从数据库表中可靠地读取事件并发布到消息代理，同时确保至少一次语义、处理故障和维持合理的延迟？

**形式化描述**:

```
给定: 事件表 E = {e₁, e₂, ..., eₙ}，其中每个事件 eᵢ 具有状态 {unpublished, published}
给定: 消息代理 B
约束:
  - 原子性: 事件标记为 published 当且仅当成功发布到 B
  - 可用性: 系统在发布者故障时可恢复
  - 延迟: 发布延迟 < Δt
目标: 设计发布函数 P: E → B 满足上述约束
```

**挑战**:

- 高频率轮询导致数据库负载
- 多个发布者实例的竞争条件
- 发布失败后的重试策略
- 大量事件的批量处理

### 1.2 解决方案形式化

**定义 1.1 (轮询发布者)**
轮询发布者是一个后台进程，周期性地：

1. 从发件箱表查询未发布事件（有限数量）
2. 尝试发布每个事件到消息代理
3. 成功发布后标记为已发布（或删除）
4. 处理失败时根据策略重试

**形式化算法**:

```
算法: PollingPublisher
输入: 数据库连接 D，消息代理 B，配置 C
输出: 无

while not terminated:
    events ← D.Query("SELECT * FROM outbox WHERE published = false LIMIT C.batch_size")

    if events is empty:
        sleep(C.poll_interval)
        continue

    for each event in events:
        success ← B.Publish(event)

        if success:
            D.Execute("UPDATE outbox SET published = true WHERE id = event.id")
        else:
            RecordFailure(event, error)

    if len(events) < C.batch_size:
        sleep(C.poll_interval)
```

**定义 1.2 (竞争控制)**
在多实例部署时，需要防止重复发布：

- 行级锁定：`SELECT ... FOR UPDATE SKIP LOCKED`
- 分布式锁：使用数据库或 Redis
- 唯一消费者：基于分区的分配

### 1.3 状态机模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Polling Publisher State Machine                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌─────────┐     Poll      ┌─────────┐     Empty      ┌─────────┐     │
│   │  Idle   │ ─────────────►│ Fetching│ ──────────────►│ Waiting │     │
│   └────▲────┘               └────┬────┘               └────┬────┘     │
│        │                         │                         │          │
│        │                         ▼                         │          │
│        │                   ┌─────────┐                     │          │
│        │                   │Processing│◄───────────────────┘          │
│        │                   └────┬────┘                               │
│        │                         │                                   │
│        │          ┌─────────────┼─────────────┐                      │
│        │          │             │             │                      │
│        │          ▼             ▼             ▼                      │
│        │    ┌─────────┐   ┌─────────┐   ┌─────────┐                 │
│        │    │ Success │   │  Retry  │   │  Fail   │                 │
│        │    └────┬────┘   └────┬────┘   └────┬────┘                 │
│        │         │             │             │                       │
│        │         ▼             ▼             ▼                       │
│        │    ┌─────────┐   ┌─────────┐   ┌─────────┐                 │
│        │    │  Mark   │   │Backoff  │   │  DLQ    │                 │
│        │    │Published│   │  Delay  │   │  Route  │                 │
│        │    └────┬────┘   └────┬────┘   └─────────┘                 │
│        │         │             │                                    │
│        └─────────┴─────────────┘                                    │
│                                                                         │
│   Error Handling States:                                               │
│   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐                │
│   │ Max Retries │──►│  Dead Letter │──►│  Alert      │                │
│   │ Exceeded    │   │  Queue       │   │  Operator   │                │
│   └─────────────┘   └─────────────┘   └─────────────┘                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心轮询发布者实现

```go
// pollingpublisher/core.go
package pollingpublisher

import (
    "context"
    "database/sql"
    "fmt"
    "sync"
    "time"
)

// Event 事件结构
type Event struct {
    ID            string
    EventType     string
    AggregateType string
    AggregateID   string
    Payload       []byte
    Metadata      []byte
    CreatedAt     time.Time
    RetryCount    int
}

// Publisher 消息发布器接口
type Publisher interface {
    Publish(ctx context.Context, topic string, event *Event) error
    Close() error
}

// Store 事件存储接口
type Store interface {
    // FetchUnpublished 获取未发布事件（带锁定）
    FetchUnpublished(ctx context.Context, batchSize int) ([]*Event, error)

    // MarkPublished 标记为已发布
    MarkPublished(ctx context.Context, ids []string) error

    // MarkFailed 标记失败（增加重试计数）
    MarkFailed(ctx context.Context, id string, errorMsg string) error

    // MoveToDeadLetter 移动到死信队列
    MoveToDeadLetter(ctx context.Context, id string, reason string) error

    // ReleaseLock 释放锁（用于失败时）
    ReleaseLock(ctx context.Context, ids []string) error
}

// Config 发布者配置
type Config struct {
    // PollInterval 轮询间隔
    PollInterval time.Duration

    // BatchSize 每批处理的事件数
    BatchSize int

    // MaxRetries 最大重试次数
    MaxRetries int

    // RetryBackoff 重试退避策略
    RetryBackoff func(attempt int) time.Duration

    // WorkerCount 并发工作线程数
    WorkerCount int

    // LockTimeout 锁超时时间
    LockTimeout time.Duration
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
    return &Config{
        PollInterval: 5 * time.Second,
        BatchSize:    100,
        MaxRetries:   3,
        RetryBackoff: func(attempt int) time.Duration {
            return time.Duration(attempt) * time.Second
        },
        WorkerCount: 5,
        LockTimeout: 5 * time.Minute,
    }
}

// PollingPublisher 轮询发布者
type PollingPublisher struct {
    store     Store
    publisher Publisher
    config    *Config

    // 控制
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
    stopCh chan struct{}

    // 指标
    metrics *Metrics

    // 日志
    logger Logger
}

// Metrics 指标
type Metrics struct {
    EventsPublished   uint64
    EventsFailed      uint64
    EventsRetried     uint64
    EventsMovedToDLQ  uint64
    PublishLatency    []time.Duration
    mu                sync.RWMutex
}

// Logger 日志接口
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// NewPollingPublisher 创建轮询发布者
func NewPollingPublisher(store Store, publisher Publisher, logger Logger, config *Config) *PollingPublisher {
    if config == nil {
        config = DefaultConfig()
    }

    ctx, cancel := context.WithCancel(context.Background())

    return &PollingPublisher{
        store:     store,
        publisher: publisher,
        config:    config,
        ctx:       ctx,
        cancel:    cancel,
        stopCh:    make(chan struct{}),
        metrics:   &Metrics{},
        logger:    logger,
    }
}

// Start 启动发布者
func (p *PollingPublisher) Start() error {
    p.logger.Info("starting polling publisher",
        Field{"poll_interval", p.config.PollInterval},
        Field{"batch_size", p.config.BatchSize},
        Field{"worker_count", p.config.WorkerCount})

    // 启动工作线程
    for i := 0; i < p.config.WorkerCount; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }

    return nil
}

// Stop 停止发布者
func (p *PollingPublisher) Stop() error {
    p.logger.Info("stopping polling publisher")

    close(p.stopCh)
    p.cancel()

    done := make(chan struct{})
    go func() {
        p.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        p.logger.Info("polling publisher stopped gracefully")
    case <-time.After(30 * time.Second):
        p.logger.Warn("polling publisher stop timeout")
    }

    return p.publisher.Close()
}

// worker 工作线程
func (p *PollingPublisher) worker(id int) {
    defer p.wg.Done()

    p.logger.Info("worker started", Field{"worker_id", id})

    ticker := time.NewTicker(p.config.PollInterval)
    defer ticker.Stop()

    for {
        select {
        case <-p.ctx.Done():
            return
        case <-p.stopCh:
            return
        case <-ticker.C:
            if err := p.processBatch(); err != nil {
                p.logger.Error("failed to process batch",
                    Field{"worker_id", id},
                    Field{"error", err})
            }
        }
    }
}

// processBatch 处理一批事件
func (p *PollingPublisher) processBatch() error {
    // 获取未发布事件
    events, err := p.store.FetchUnpublished(p.ctx, p.config.BatchSize)
    if err != nil {
        return fmt.Errorf("failed to fetch unpublished events: %w", err)
    }

    if len(events) == 0 {
        return nil
    }

    p.logger.Debug("processing batch", Field{"count", len(events)})

    var published []string
    var failed []*Event

    for _, event := range events {
        if err := p.publishEvent(event); err != nil {
            p.logger.Error("failed to publish event",
                Field{"event_id", event.ID},
                Field{"error", err})
            failed = append(failed, event)
        } else {
            published = append(published, event.ID)
        }
    }

    // 标记已发布
    if len(published) > 0 {
        if err := p.store.MarkPublished(p.ctx, published); err != nil {
            p.logger.Error("failed to mark events as published",
                Field{"error", err})
            // 这些事件会被重新处理（幂等性保证）
        }

        p.updateMetricsPublished(len(published))
    }

    // 处理失败的事件
    for _, event := range failed {
        p.handleFailedEvent(event)
    }

    return nil
}

// publishEvent 发布单个事件
func (p *PollingPublisher) publishEvent(event *Event) error {
    start := time.Now()
    defer func() {
        p.recordLatency(time.Since(start))
    }()

    topic := fmt.Sprintf("%s.%s", event.AggregateType, event.EventType)

    ctx, cancel := context.WithTimeout(p.ctx, 30*time.Second)
    defer cancel()

    if err := p.publisher.Publish(ctx, topic, event); err != nil {
        return err
    }

    return nil
}

// handleFailedEvent 处理失败事件
func (p *PollingPublisher) handleFailedEvent(event *Event) {
    p.updateMetricsFailed()

    if event.RetryCount >= p.config.MaxRetries {
        // 移动到死信队列
        if err := p.store.MoveToDeadLetter(p.ctx, event.ID, "max_retries_exceeded"); err != nil {
            p.logger.Error("failed to move event to DLQ",
                Field{"event_id", event.ID},
                Field{"error", err})
        } else {
            p.updateMetricsDLQ()
            p.logger.Warn("event moved to DLQ",
                Field{"event_id", event.ID},
                Field{"retry_count", event.RetryCount})
        }
    } else {
        // 标记失败，增加重试计数
        if err := p.store.MarkFailed(p.ctx, event.ID, "publish_failed"); err != nil {
            p.logger.Error("failed to mark event as failed",
                Field{"event_id", event.ID},
                Field{"error", err})
        }
        p.updateMetricsRetried()
    }
}

// recordLatency 记录延迟
func (p *PollingPublisher) recordLatency(d time.Duration) {
    p.metrics.mu.Lock()
    defer p.metrics.mu.Unlock()
    p.metrics.PublishLatency = append(p.metrics.PublishLatency, d)

    // 只保留最近的 1000 个样本
    if len(p.metrics.PublishLatency) > 1000 {
        p.metrics.PublishLatency = p.metrics.PublishLatency[len(p.metrics.PublishLatency)-1000:]
    }
}

// updateMetricsPublished 更新发布计数
func (p *PollingPublisher) updateMetricsPublished(count int) {
    p.metrics.mu.Lock()
    defer p.metrics.mu.Unlock()
    p.metrics.EventsPublished += uint64(count)
}

// updateMetricsFailed 更新失败计数
func (p *PollingPublisher) updateMetricsFailed() {
    p.metrics.mu.Lock()
    defer p.metrics.mu.Unlock()
    p.metrics.EventsFailed++
}

// updateMetricsRetried 更新重试计数
func (p *PollingPublisher) updateMetricsRetried() {
    p.metrics.mu.Lock()
    defer p.metrics.mu.Unlock()
    p.metrics.EventsRetried++
}

// updateMetricsDLQ 更新 DLQ 计数
func (p *PollingPublisher) updateMetricsDLQ() {
    p.metrics.mu.Lock()
    defer p.metrics.mu.Unlock()
    p.metrics.EventsMovedToDLQ++
}

// GetMetrics 获取指标
func (p *PollingPublisher) GetMetrics() Metrics {
    p.metrics.mu.RLock()
    defer p.metrics.mu.RUnlock()
    return *p.metrics
}
```

### 2.2 数据库存储实现（带行级锁定）

```go
// pollingpublisher/sql_store.go
package pollingpublisher

import (
    "context"
    "database/sql"
    "fmt"
    "time"
)

// SQLStore SQL 存储实现
type SQLStore struct {
    db          *sql.DB
    outboxTable string
    dlqTable    string
    lockTimeout time.Duration
}

// NewSQLStore 创建 SQL 存储
func NewSQLStore(db *sql.DB, outboxTable, dlqTable string, lockTimeout time.Duration) *SQLStore {
    return &SQLStore{
        db:          db,
        outboxTable: outboxTable,
        dlqTable:    dlqTable,
        lockTimeout: lockTimeout,
    }
}

// FetchUnpublished 获取未发布事件（带 SKIP LOCKED）
func (s *SQLStore) FetchUnpublished(ctx context.Context, batchSize int) ([]*Event, error) {
    query := fmt.Sprintf(`
        SELECT id, event_type, aggregate_type, aggregate_id, payload, metadata, created_at, retry_count
        FROM %s
        WHERE published_at IS NULL
          AND (locked_at IS NULL OR locked_at < ?)
          AND (next_retry_at IS NULL OR next_retry_at <= ?)
        ORDER BY created_at ASC
        LIMIT ?
        FOR UPDATE SKIP LOCKED
    `, s.outboxTable)

    lockExpiry := time.Now().Add(-s.lockTimeout)
    now := time.Now()

    rows, err := s.db.QueryContext(ctx, query, lockExpiry, now, batchSize)
    if err != nil {
        return nil, fmt.Errorf("failed to query outbox: %w", err)
    }
    defer rows.Close()

    var events []*Event
    var ids []string

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
            &event.RetryCount,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan event: %w", err)
        }
        events = append(events, event)
        ids = append(ids, event.ID)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    // 设置锁定
    if len(ids) > 0 {
        if err := s.lockEvents(ctx, ids); err != nil {
            return nil, fmt.Errorf("failed to lock events: %w", err)
        }
    }

    return events, nil
}

// lockEvents 锁定事件
func (s *SQLStore) lockEvents(ctx context.Context, ids []string) error {
    placeholders := make([]string, len(ids))
    args := make([]interface{}, len(ids)+1)
    args[0] = time.Now()

    for i, id := range ids {
        placeholders[i] = "?"
        args[i+1] = id
    }

    query := fmt.Sprintf(`
        UPDATE %s
        SET locked_at = ?
        WHERE id IN (%s)
    `, s.outboxTable, joinPlaceholders(placeholders))

    _, err := s.db.ExecContext(ctx, query, args...)
    return err
}

// MarkPublished 标记为已发布
func (s *SQLStore) MarkPublished(ctx context.Context, ids []string) error {
    if len(ids) == 0 {
        return nil
    }

    placeholders := make([]string, len(ids))
    args := make([]interface{}, len(ids)+1)
    args[0] = time.Now()

    for i, id := range ids {
        placeholders[i] = "?"
        args[i+1] = id
    }

    query := fmt.Sprintf(`
        UPDATE %s
        SET published_at = ?, locked_at = NULL
        WHERE id IN (%s)
    `, s.outboxTable, joinPlaceholders(placeholders))

    _, err := s.db.ExecContext(ctx, query, args...)
    return err
}

// MarkFailed 标记失败
func (s *SQLStore) MarkFailed(ctx context.Context, id string, errorMsg string) error {
    query := fmt.Sprintf(`
        UPDATE %s
        SET retry_count = retry_count + 1,
            last_error = ?,
            locked_at = NULL,
            next_retry_at = ?
        WHERE id = ?
    `, s.outboxTable)

    // 简单的指数退避
    nextRetry := time.Now().Add(time.Duration(1) * time.Minute)

    _, err := s.db.ExecContext(ctx, query, errorMsg, nextRetry, id)
    return err
}

// MoveToDeadLetter 移动到死信队列
func (s *SQLStore) MoveToDeadLetter(ctx context.Context, id string, reason string) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 插入到死信队列
    insertQuery := fmt.Sprintf(`
        INSERT INTO %s (id, event_type, aggregate_type, aggregate_id, payload, metadata,
                        created_at, retry_count, last_error, moved_at, move_reason)
        SELECT id, event_type, aggregate_type, aggregate_id, payload, metadata,
               created_at, retry_count, last_error, ?, ?
        FROM %s
        WHERE id = ?
    `, s.dlqTable, s.outboxTable)

    _, err = tx.ExecContext(ctx, insertQuery, time.Now(), reason, id)
    if err != nil {
        return fmt.Errorf("failed to insert into DLQ: %w", err)
    }

    // 从发件箱删除
    deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, s.outboxTable)
    _, err = tx.ExecContext(ctx, deleteQuery, id)
    if err != nil {
        return fmt.Errorf("failed to delete from outbox: %w", err)
    }

    return tx.Commit()
}

// ReleaseLock 释放锁
func (s *SQLStore) ReleaseLock(ctx context.Context, ids []string) error {
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
        UPDATE %s
        SET locked_at = NULL
        WHERE id IN (%s)
    `, s.outboxTable, joinPlaceholders(placeholders))

    _, err := s.db.ExecContext(ctx, query, args...)
    return err
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

// CreateTablesSQL 创建表的 SQL
func (s *SQLStore) CreateTablesSQL() string {
    return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    id VARCHAR(36) PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    payload BLOB NOT NULL,
    metadata BLOB,
    created_at TIMESTAMP NOT NULL,
    published_at TIMESTAMP NULL,
    locked_at TIMESTAMP NULL,
    next_retry_at TIMESTAMP NULL,
    retry_count INT DEFAULT 0,
    last_error TEXT,
    INDEX idx_unpublished (published_at, next_retry_at, created_at),
    INDEX idx_locked (locked_at),
    INDEX idx_aggregate (aggregate_type, aggregate_id)
);

CREATE TABLE IF NOT EXISTS %s (
    id VARCHAR(36) PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    payload BLOB NOT NULL,
    metadata BLOB,
    created_at TIMESTAMP NOT NULL,
    retry_count INT DEFAULT 0,
    last_error TEXT,
    moved_at TIMESTAMP NOT NULL,
    move_reason VARCHAR(255),
    INDEX idx_moved_at (moved_at)
);`, s.outboxTable, s.dlqTable)
}
```

### 2.3 Kafka 发布器实现

```go
// pollingpublisher/kafka_publisher.go
package pollingpublisher

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

// KafkaConfig Kafka 配置
type KafkaConfig struct {
    Brokers      []string
    RequiredAcks sarama.RequiredAcks
    RetryMax     int
}

// DefaultKafkaConfig 默认 Kafka 配置
func DefaultKafkaConfig() *KafkaConfig {
    return &KafkaConfig{
        Brokers:      []string{"localhost:9092"},
        RequiredAcks: sarama.WaitForAll,
        RetryMax:     3,
    }
}

// NewKafkaPublisher 创建 Kafka 发布器
func NewKafkaPublisher(config *KafkaConfig) (*KafkaPublisher, error) {
    if config == nil {
        config = DefaultKafkaConfig()
    }

    saramaConfig := sarama.NewConfig()
    saramaConfig.Producer.RequiredAcks = config.RequiredAcks
    saramaConfig.Producer.Retry.Max = config.RetryMax
    saramaConfig.Producer.Return.Successes = true
    saramaConfig.Producer.Idempotent = true // 启用幂等生产者

    producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create kafka producer: %w", err)
    }

    return &KafkaPublisher{producer: producer}, nil
}

// Publish 发布事件到 Kafka
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

### 2.4 应用示例

```go
// examples/polling/main.go
package main

import (
    "context"
    "database/sql"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "go-knowledge-base/pollingpublisher"
)

type logger struct{}

func (l *logger) Info(msg string, fields ...pollingpublisher.Field) {
    log.Printf("[INFO] %s %+v", msg, fields)
}

func (l *logger) Error(msg string, fields ...pollingpublisher.Field) {
    log.Printf("[ERROR] %s %+v", msg, fields)
}

func (l *logger) Debug(msg string, fields ...pollingpublisher.Field) {
    log.Printf("[DEBUG] %s %+v", msg, fields)
}

func (l *logger) Warn(msg string, fields ...pollingpublisher.Field) {
    log.Printf("[WARN] %s %+v", msg, fields)
}

func main() {
    // 打开数据库
    db, err := sql.Open("sqlite3", "outbox.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 创建表
    store := pollingpublisher.NewSQLStore(db, "outbox", "dead_letter", 5*time.Minute)
    if _, err := db.Exec(store.CreateTablesSQL()); err != nil {
        log.Fatal(err)
    }

    // 创建 Kafka 发布器
    kafkaPub, err := pollingpublisher.NewKafkaPublisher(&pollingpublisher.KafkaConfig{
        Brokers: []string{"localhost:9092"},
    })
    if err != nil {
        log.Printf("Failed to create Kafka publisher (using mock): %v", err)
        kafkaPub = &mockPublisher{}
    }

    // 创建发布者
    logger := &logger{}
    config := &pollingpublisher.Config{
        PollInterval: 2 * time.Second,
        BatchSize:    50,
        MaxRetries:   3,
        WorkerCount:  3,
    }

    publisher := pollingpublisher.NewPollingPublisher(store, kafkaPub, logger, config)

    // 启动发布者
    if err := publisher.Start(); err != nil {
        log.Fatal(err)
    }

    // 等待中断信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    <-sigCh

    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    _ = ctx
    if err := publisher.Stop(); err != nil {
        log.Printf("Error stopping publisher: %v", err)
    }

    // 打印指标
    metrics := publisher.GetMetrics()
    log.Printf("Final metrics: Published=%d, Failed=%d, Retried=%d, DLQ=%d",
        metrics.EventsPublished,
        metrics.EventsFailed,
        metrics.EventsRetried,
        metrics.EventsMovedToDLQ)
}

// mockPublisher 模拟发布器（用于测试）
type mockPublisher struct{}

func (m *mockPublisher) Publish(ctx context.Context, topic string, event *pollingpublisher.Event) error {
    log.Printf("[MOCK] Publishing to %s: %s", topic, event.ID)
    return nil
}

func (m *mockPublisher) Close() error {
    return nil
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// pollingpublisher/core_test.go
package pollingpublisher

import (
    "context"
    "errors"
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
func (m *mockLogger) Warn(msg string, fields ...Field)  {}

type mockStore struct {
    mock.Mock
}

func (m *mockStore) FetchUnpublished(ctx context.Context, batchSize int) ([]*Event, error) {
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

func (m *mockStore) MarkFailed(ctx context.Context, id string, errorMsg string) error {
    args := m.Called(ctx, id, errorMsg)
    return args.Error(0)
}

func (m *mockStore) MoveToDeadLetter(ctx context.Context, id string, reason string) error {
    args := m.Called(ctx, id, reason)
    return args.Error(0)
}

func (m *mockStore) ReleaseLock(ctx context.Context, ids []string) error {
    args := m.Called(ctx, ids)
    return args.Error(0)
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

func TestPollingPublisher_processBatch(t *testing.T) {
    store := new(mockStore)
    publisher := new(mockPublisher)
    logger := &mockLogger{}

    config := &Config{
        PollInterval: 1 * time.Second,
        BatchSize:    10,
        MaxRetries:   3,
        WorkerCount:  1,
    }

    pub := NewPollingPublisher(store, publisher, logger, config)

    events := []*Event{
        {ID: "evt-001", EventType: "Test", AggregateType: "test"},
        {ID: "evt-002", EventType: "Test", AggregateType: "test"},
    }

    store.On("FetchUnpublished", mock.Anything, 10).Return(events, nil)
    publisher.On("Publish", mock.Anything, "test.Test", events[0]).Return(nil)
    publisher.On("Publish", mock.Anything, "test.Test", events[1]).Return(nil)
    store.On("MarkPublished", mock.Anything, []string{"evt-001", "evt-002"}).Return(nil)

    err := pub.processBatch()

    require.NoError(t, err)
    store.AssertExpectations(t)
    publisher.AssertExpectations(t)
}

func TestPollingPublisher_processBatch_WithFailures(t *testing.T) {
    store := new(mockStore)
    publisher := new(mockPublisher)
    logger := &mockLogger{}

    config := &Config{
        BatchSize:  10,
        MaxRetries: 3,
    }

    pub := NewPollingPublisher(store, publisher, logger, config)

    events := []*Event{
        {ID: "evt-001", EventType: "Test", AggregateType: "test", RetryCount: 0},
        {ID: "evt-002", EventType: "Test", AggregateType: "test", RetryCount: 3},
    }

    store.On("FetchUnpublished", mock.Anything, 10).Return(events, nil)
    publisher.On("Publish", mock.Anything, "test.Test", events[0]).Return(errors.New("publish failed"))
    publisher.On("Publish", mock.Anything, "test.Test", events[1]).Return(errors.New("publish failed"))
    store.On("MarkFailed", mock.Anything, "evt-001", "publish_failed").Return(nil)
    store.On("MoveToDeadLetter", mock.Anything, "evt-002", "max_retries_exceeded").Return(nil)

    err := pub.processBatch()

    require.NoError(t, err)
    store.AssertExpectations(t)
}

func TestConfig_DefaultConfig(t *testing.T) {
    config := DefaultConfig()

    assert.Equal(t, 5*time.Second, config.PollInterval)
    assert.Equal(t, 100, config.BatchSize)
    assert.Equal(t, 3, config.MaxRetries)
    assert.Equal(t, 5, config.WorkerCount)
}
```

### 3.2 性能测试

```go
// pollingpublisher/benchmark_test.go
package pollingpublisher

import (
    "context"
    "testing"
    "time"
)

func BenchmarkProcessBatch(b *testing.B) {
    // 创建模拟存储和发布器
    store := &benchmarkStore{events: make([]*Event, 100)}
    publisher := &benchmarkPublisher{}
    logger := &mockLogger{}

    config := &Config{
        BatchSize:  100,
        MaxRetries: 3,
    }

    pub := NewPollingPublisher(store, publisher, logger, config)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = pub.processBatch()
    }
}

type benchmarkStore struct {
    events []*Event
}

func (s *benchmarkStore) FetchUnpublished(ctx context.Context, batchSize int) ([]*Event, error) {
    return s.events, nil
}

func (s *benchmarkStore) MarkPublished(ctx context.Context, ids []string) error {
    return nil
}

func (s *benchmarkStore) MarkFailed(ctx context.Context, id string, errorMsg string) error {
    return nil
}

func (s *benchmarkStore) MoveToDeadLetter(ctx context.Context, id string, reason string) error {
    return nil
}

func (s *benchmarkStore) ReleaseLock(ctx context.Context, ids []string) error {
    return nil
}

type benchmarkPublisher struct{}

func (p *benchmarkPublisher) Publish(ctx context.Context, topic string, event *Event) error {
    return nil
}

func (p *benchmarkPublisher) Close() error {
    return nil
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Transactional Outbox 的关系

```
┌─────────────────────────────────────────────────────────────────────────┐
│          Transactional Outbox + Polling Publisher Flow                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────┐         ┌─────────────┐         ┌─────────────┐       │
│  │   Service   │────────►│  Database   │────────►│   Polling   │       │
│  │   (Write)   │  TXN    │   Outbox    │  Poll   │  Publisher  │       │
│  └─────────────┘         └─────────────┘         └──────┬──────┘       │
│        │                                                │              │
│        │                                                │ Publish      │
│        │                                                ▼              │
│  ┌─────┴─────┐                                   ┌─────────────┐       │
│  │  Business │                                   │   Kafka/    │       │
│  │   Tables  │                                   │  RabbitMQ   │       │
│  └───────────┘                                   └──────┬──────┘       │
│                                                         │               │
│                                            ┌────────────┼────────────┐ │
│                                            ▼            ▼            ▼ │
│                                       ┌─────────┐  ┌─────────┐  ┌─────────┐
│                                       │Consumer1│  │Consumer2│  │Consumer3│
│                                       └─────────┘  └─────────┘  └─────────┘
│                                                                         │
│  关键协作点:                                                               │
│  1. Service 在事务中同时写入业务表和 Outbox 表                               │
│  2. Polling Publisher 定期轮询 Outbox 表                                    │
│  3. 使用 FOR UPDATE SKIP LOCKED 避免多实例竞争                              │
│  4. 成功发布后标记为 published（或删除）                                     │
│  5. 失败时重试，超过阈值后移至 DLQ                                          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 与 CDC 的替代方案对比

| 特性 | Polling Publisher | CDC (Debezium) |
|------|-------------------|----------------|
| **延迟** | 秒级（可配置） | 毫秒级 |
| **数据库负载** | 中等（定期查询） | 低（读取 binlog）|
| **实现复杂度** | 低 | 高 |
| **部署复杂度** | 低 | 高（需 Kafka Connect）|
| **数据丢失风险** | 低 | 极低 |
| **适用场景** | 大多数场景 | 实时性要求高的场景 |

---

## 5. 决策标准

### 5.1 配置参数选择

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Polling Publisher Configuration Guide                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  PollInterval (轮询间隔)                                                │
│  ├── 低延迟需求: 100ms - 1s                                             │
│  ├── 一般场景: 1s - 5s  [默认]                                          │
│  └── 低频率: 5s - 30s                                                   │
│                                                                         │
│  BatchSize (批次大小)                                                   │
│  ├── 低吞吐量: 10-50                                                    │
│  ├── 中等吞吐量: 50-100  [默认]                                         │
│  ├── 高吞吐量: 100-500                                                  │
│  └── 极高吞吐量: 500+ (需测试数据库性能)                                 │
│                                                                         │
│  WorkerCount (工作线程)                                                 │
│  ├── 单实例: 1-3                                                        │
│  ├── 中等负载: 3-5  [默认]                                              │
│  └── 高负载: 5-10（配合多实例部署）                                      │
│                                                                         │
│  MaxRetries (最大重试)                                                  │
│  ├── 快速失败: 1-2                                                      │
│  ├── 平衡: 3  [默认]                                                    │
│  └── 高可靠性: 5+                                                       │
│                                                                         │
│  LockTimeout (锁超时)                                                   │
│  ├── 快速恢复: 30s-1m                                                   │
│  ├── 正常: 5m  [默认]                                                   │
│  └── 慢消费者: 10m+                                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 生产环境检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Polling Publisher Production Checklist               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  部署:                                                                   │
│  □ 多实例部署（水平扩展）                                                │
│  □ 使用行级锁定（SKIP LOCKED）防止竞争                                  │
│  □ 配置健康检查端点                                                      │
│  □ 实现优雅关闭（处理完当前批次）                                         │
│                                                                         │
│  监控:                                                                   │
│  □ 发布延迟（P95/P99）                                                  │
│  □ 积压事件数量                                                          │
│  □ 发布成功率                                                            │
│  □ 死信队列大小                                                          │
│  □ 数据库连接池使用率                                                    │
│                                                                         │
│  告警:                                                                   │
│  □ 发布延迟 > 阈值                                                       │
│  □ 积压事件持续增长                                                      │
│  □ 发布成功率 < 95%                                                     │
│  □ 死信队列有新事件                                                      │
│  □ 轮询间隔内未完成批次（处理能力不足）                                  │
│                                                                         │
│  运维:                                                                   │
│  □ 定期清理已发布事件（或归档）                                          │
│  □ 死信队列人工处理流程                                                  │
│  □ 发布者重启/迁移流程                                                   │
│  □ 性能基线和扩容策略                                                    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>16KB, 完整形式化 + Go 实现 + 测试 + 配置指南)

**相关文档**:

- [EC-033-Transactional-Outbox.md](./EC-033-Transactional-Outbox.md)
- [EC-031-Choreography-Pattern.md](./EC-031-Choreography-Pattern.md)
