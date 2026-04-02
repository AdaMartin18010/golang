# 任务队列实现模式 (Task Queue Implementation Patterns)

> **分类**: 工程与云原生
> **标签**: #task-queue #patterns #implementation #redis
> **参考**: Redis Streams, Kafka, RabbitMQ, Amazon SQS

---

## 队列架构模式

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        Task Queue Architectures                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. 简单队列 (Simple Queue)                                              │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐                              │
│  │Producer │───▶│  Queue  │───▶│Consumer │                             │
│  └─────────┘    └─────────┘    └─────────┘                              │
│                                                                         │
│  2. 工作队列 (Work Queue)                                                │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐                              │
│  │Producer │───▶│  Queue  │───▶│Worker 1 │                             │
│  └─────────┘    └─────────┘    ├─────────┤                              │
│                                  │Worker 2 │                             │
│                                  ├─────────┤                             │
│                                  │Worker N │                             │
│                                  └─────────┘                             │
│                                                                          │
│  3. 发布/订阅 (Pub/Sub)                                                  │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐                              │
│  │Producer │───▶│Exchange │───▶│Queue 1  │───▶ Consumer 1              │
│  └─────────┘    │(Fanout) │    └─────────┘                              │
│                 │         │    ┌─────────┐                              │
│                 │         ├───▶│Queue 2  │───▶ Consumer 2              │
│                 │         │    └─────────┘                              │
│                 └─────────┘    ┌─────────┐                              │
│                                │Queue N  │───▶ Consumer N              │
│                                └─────────┘                              │
│                                                                         │
│  4. 优先级队列 (Priority Queue)                                          │
│  ┌─────────┐    ┌─────────────────┐    ┌─────────┐                      │
│  │Producer │───▶│High Priority    │───▶│Consumer │                     │
│  │         │    ├─────────────────┤    └─────────┘                      │
│  │         │───▶│Medium Priority  │                                    │
│  │         │    ├─────────────────┤                                     │
│  │         │───▶│Low Priority     │                                    │
│  └─────────┘    └─────────────────┘                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Redis 队列实现

```go
// Redis 队列实现
package queue

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisQueue Redis 队列实现
type RedisQueue struct {
    client *redis.Client
    name   string
}

type Task struct {
    ID        string          `json:"id"`
    Type      string          `json:"type"`
    Payload   json.RawMessage `json:"payload"`
    Priority  int             `json:"priority"`
    CreatedAt time.Time       `json:"created_at"`
    DelayedTo *time.Time      `json:"delayed_to,omitempty"`
}

// Push 添加任务到队列
func (q *RedisQueue) Push(ctx context.Context, task *Task) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }

    return q.client.LPush(ctx, q.name, data).Err()
}

// PushPriority 添加优先级任务 (使用有序集合)
func (q *RedisQueue) PushPriority(ctx context.Context, task *Task) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }

    // 使用负优先级，使高优先级排在前面
    score := -float64(task.Priority)

    return q.client.ZAdd(ctx, q.name+":priority", redis.Z{
        Score:  score,
        Member: data,
    }).Err()
}

// Pop 阻塞式获取任务
func (q *RedisQueue) Pop(ctx context.Context, timeout time.Duration) (*Task, error) {
    result, err := q.client.BRPop(ctx, timeout, q.name).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var task Task
    if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
        return nil, err
    }

    return &task, nil
}

// PopPriority 获取高优先级任务
func (q *RedisQueue) PopPriority(ctx context.Context) (*Task, error) {
    // 获取优先级最高的任务
    result, err := q.client.ZPopMax(ctx, q.name+":priority").Result()
    if err != nil {
        return nil, err
    }
    if len(result) == 0 {
        return nil, nil
    }

    var task Task
    if err := json.Unmarshal([]byte(result[0].Member.(string)), &task); err != nil {
        return nil, err
    }

    return &task, nil
}

// DelayedPush 延迟任务
func (q *RedisQueue) DelayedPush(ctx context.Context, task *Task, delay time.Duration) error {
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }

    executeAt := time.Now().Add(delay)

    return q.client.ZAdd(ctx, q.name+":delayed", redis.Z{
        Score:  float64(executeAt.Unix()),
        Member: data,
    }).Err()
}

// ProcessDelayed 处理延迟任务
func (q *RedisQueue) ProcessDelayed(ctx context.Context) error {
    now := float64(time.Now().Unix())

    // 获取已到期的延迟任务
    tasks, err := q.client.ZRangeByScoreWithScores(ctx, q.name+":delayed", &redis.ZRangeBy{
        Min: "0",
        Max: fmt.Sprintf("%f", now),
    }).Result()
    if err != nil {
        return err
    }

    for _, task := range tasks {
        // 添加到主队列
        if err := q.client.LPush(ctx, q.name, task.Member).Err(); err != nil {
            return err
        }

        // 从延迟队列移除
        if err := q.client.ZRem(ctx, q.name+":delayed", task.Member).Err(); err != nil {
            return err
        }
    }

    return nil
}

// Acknowledge 确认任务完成
func (q *RedisQueue) Acknowledge(ctx context.Context, taskID string) error {
    return q.client.HDel(ctx, q.name+":processing", taskID).Err()
}

// RequeueFailed 重新入队失败任务
func (q *RedisQueue) RequeueFailed(ctx context.Context, task *Task) error {
    // 增加重试计数
    task.RetryCount++

    // 计算退避延迟
    backoff := calculateBackoff(task.RetryCount)

    if backoff > 0 {
        return q.DelayedPush(ctx, task, backoff)
    }

    return q.Push(ctx, task)
}

func calculateBackoff(retryCount int) time.Duration {
    // 指数退避
    baseDelay := time.Second
    maxDelay := time.Hour

    delay := baseDelay * time.Duration(1<<uint(retryCount))
    if delay > maxDelay {
        delay = maxDelay
    }

    // 添加 jitter
    jitter := time.Duration(rand.Int63n(int64(delay) / 2))
    return delay + jitter
}
```

---

## 可靠队列模式

```go
// ReliableQueue 可靠队列 (至少一次交付)
type ReliableQueue struct {
    redis      *redis.Client
    queueKey   string
    processingKey string
    deadLetterKey string
    maxRetries    int
}

// Consume 消费任务 (带确认机制)
func (q *ReliableQueue) Consume(ctx context.Context, handler Handler) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // 使用 BRPOPLPUSH 实现可靠队列
        result, err := q.redis.BRPopLPush(ctx, q.queueKey, q.processingKey, 5*time.Second).Result()
        if err == redis.Nil {
            continue
        }
        if err != nil {
            return err
        }

        var task Task
        if err := json.Unmarshal([]byte(result), &task); err != nil {
            continue
        }

        // 处理任务
        err = handler(ctx, &task)

        if err == nil {
            // 成功：从处理队列移除
            q.redis.LRem(ctx, q.processingKey, 0, result)
        } else {
            // 失败：处理重试或死信
            q.handleFailure(ctx, &task, result, err)
        }
    }
}

func (q *ReliableQueue) handleFailure(ctx context.Context, task *Task, rawTask string, err error) {
    task.RetryCount++

    if task.RetryCount >= q.maxRetries {
        // 移至死信队列
        q.redis.LRem(ctx, q.processingKey, 0, rawTask)
        q.redis.LPush(ctx, q.deadLetterKey, rawTask)

        // 记录失败事件
        log.Printf("Task %s moved to dead letter queue after %d retries: %v",
            task.ID, task.RetryCount, err)
    } else {
        // 重新入队
        q.redis.LRem(ctx, q.processingKey, 0, rawTask)

        backoff := calculateBackoff(task.RetryCount)
        delayedTask, _ := json.Marshal(task)

        q.redis.ZAdd(ctx, q.queueKey+":delayed", redis.Z{
            Score:  float64(time.Now().Add(backoff).Unix()),
            Member: delayedTask,
        })
    }
}

// RecoverStalledTasks 恢复卡住的任务
func (q *ReliableQueue) RecoverStalledTasks(ctx context.Context, timeout time.Duration) error {
    // 获取处理队列中的所有任务
    tasks, err := q.redis.LRange(ctx, q.processingKey, 0, -1).Result()
    if err != nil {
        return err
    }

    for _, rawTask := range tasks {
        var task Task
        if err := json.Unmarshal([]byte(rawTask), &task); err != nil {
            continue
        }

        // 检查任务是否超时
        if time.Since(task.StartedAt) > timeout {
            // 重新入队
            q.redis.LRem(ctx, q.processingKey, 0, rawTask)
            q.redis.LPush(ctx, q.queueKey, rawTask)

            log.Printf("Recovered stalled task: %s", task.ID)
        }
    }

    return nil
}
```

---

## 批量处理模式

```go
// BatchQueue 批量队列
type BatchQueue struct {
    redis      *redis.Client
    queueKey   string
    batchSize  int
    flushInterval time.Duration
}

// ConsumeBatch 批量消费
func (q *BatchQueue) ConsumeBatch(ctx context.Context, handler BatchHandler) error {
    ticker := time.NewTicker(q.flushInterval)
    defer ticker.Stop()

    var batch []Task

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case <-ticker.C:
            if len(batch) > 0 {
                if err := handler(ctx, batch); err != nil {
                    // 逐个重试
                    q.retryIndividually(ctx, batch)
                }
                batch = batch[:0]
            }

        default:
            // 非阻塞获取
            result, err := q.redis.RPop(ctx, q.queueKey).Result()
            if err == redis.Nil {
                time.Sleep(10 * time.Millisecond)
                continue
            }
            if err != nil {
                return err
            }

            var task Task
            if err := json.Unmarshal([]byte(result), &task); err != nil {
                continue
            }

            batch = append(batch, task)

            // 达到批次大小，立即处理
            if len(batch) >= q.batchSize {
                if err := handler(ctx, batch); err != nil {
                    q.retryIndividually(ctx, batch)
                }
                batch = batch[:0]
            }
        }
    }
}

func (q *BatchQueue) retryIndividually(ctx context.Context, batch []Task) {
    for _, task := range batch {
        data, _ := json.Marshal(task)
        q.redis.LPush(ctx, q.queueKey, data)
    }
}
```

---

## 流式队列 (Redis Streams)

```go
// RedisStream Redis Streams 实现
type RedisStream struct {
    client   *redis.Client
    stream   string
    group    string
    consumer string
}

// Produce 生产消息
func (s *RedisStream) Produce(ctx context.Context, task *Task) error {
    data, _ := json.Marshal(task)

    values := map[string]interface{}{
        "task": string(data),
    }

    return s.client.XAdd(ctx, &redis.XAddArgs{
        Stream: s.stream,
        Values: values,
    }).Err()
}

// Consume 消费消息 (消费者组)
func (s *RedisStream) Consume(ctx context.Context, handler Handler) error {
    // 创建消费者组
    _ = s.client.XGroupCreateMkStream(ctx, s.stream, s.group, "$").Err()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // 读取消息
        streams, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    s.group,
            Consumer: s.consumer,
            Streams:  []string{s.stream, ">"},
            Count:    10,
            Block:    5 * time.Second,
        }).Result()

        if err == redis.Nil {
            continue
        }
        if err != nil {
            return err
        }

        for _, stream := range streams {
            for _, message := range stream.Messages {
                var task Task
                taskData := message.Values["task"].(string)
                if err := json.Unmarshal([]byte(taskData), &task); err != nil {
                    continue
                }

                // 处理任务
                err := handler(ctx, &task)

                if err == nil {
                    // 确认消息
                    s.client.XAck(ctx, s.stream, s.group, message.ID)
                } else {
                    // 处理失败，消息保留在 Pending 列表中
                    log.Printf("Task %s failed: %v", task.ID, err)
                }
            }
        }
    }
}

// ClaimPending 认领悬停消息 (其他消费者崩溃)
func (s *RedisStream) ClaimPending(ctx context.Context, minIdle time.Duration) error {
    // 获取 Pending 消息
    pending, err := s.client.XPendingExt(ctx, &redis.XPendingExtArgs{
        Stream: s.stream,
        Group:  s.group,
        Start:  "-",
        End:    "+",
        Count:  100,
    }).Result()

    if err != nil {
        return err
    }

    for _, p := range pending {
        if p.Idle >= minIdle {
            // 认领消息
            _, err := s.client.XClaim(ctx, &redis.XClaimArgs{
                Stream:   s.stream,
                Group:    s.group,
                Consumer: s.consumer,
                MinIdle:  minIdle,
                Messages: []string{p.ID},
            }).Result()

            if err != nil {
                log.Printf("Failed to claim message %s: %v", p.ID, err)
            }
        }
    }

    return nil
}
```

---

## 队列监控与管理

```go
// QueueMonitor 队列监控
type QueueMonitor struct {
    redis *redis.Client
}

// Metrics 队列指标
type Metrics struct {
    QueueDepth       int64
    ProcessingCount  int64
    DelayedCount     int64
    DeadLetterCount  int64
    Throughput       float64
}

func (m *QueueMonitor) GetMetrics(ctx context.Context, queueName string) (*Metrics, error) {
    pipe := m.redis.Pipeline()

    queueLen := pipe.LLen(ctx, queueName)
    processingLen := pipe.LLen(ctx, queueName+":processing")
    delayedLen := pipe.ZCard(ctx, queueName+":delayed")
    deadLetterLen := pipe.LLen(ctx, queueName+":dlq")

    _, err := pipe.Exec(ctx)
    if err != nil {
        return nil, err
    }

    return &Metrics{
        QueueDepth:      queueLen.Val(),
        ProcessingCount: processingLen.Val(),
        DelayedCount:    delayedLen.Val(),
        DeadLetterCount: deadLetterLen.Val(),
    }, nil
}

// PurgeQueue 清空队列
func (m *QueueMonitor) PurgeQueue(ctx context.Context, queueName string) error {
    pipe := m.redis.Pipeline()
    pipe.Del(ctx, queueName)
    pipe.Del(ctx, queueName+":processing")
    pipe.Del(ctx, queueName+":delayed")
    pipe.Del(ctx, queueName+":priority")

    _, err := pipe.Exec(ctx)
    return err
}

// MoveToDeadLetter 手动移入死信队列
func (m *QueueMonitor) MoveToDeadLetter(ctx context.Context, queueName string, taskID string) error {
    // 从处理队列移除并加入死信队列
    tasks, err := m.redis.LRange(ctx, queueName+":processing", 0, -1).Result()
    if err != nil {
        return err
    }

    for _, rawTask := range tasks {
        var task Task
        if err := json.Unmarshal([]byte(rawTask), &task); err != nil {
            continue
        }

        if task.ID == taskID {
            m.redis.LRem(ctx, queueName+":processing", 0, rawTask)
            m.redis.LPush(ctx, queueName+":dlq", rawTask)
            return nil
        }
    }

    return fmt.Errorf("task %s not found in processing queue", taskID)
}
```
