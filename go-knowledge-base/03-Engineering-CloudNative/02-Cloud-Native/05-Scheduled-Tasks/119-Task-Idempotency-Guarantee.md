# 幂等性保证机制 (Idempotency Guarantee Mechanism)

> **分类**: 工程与云原生
> **标签**: #idempotency #exactly-once #deduplication
> **参考**: Stripe Idempotency Keys, AWS Lambda, Kafka Idempotent Producer

---

## 幂等性核心问题

```
客户端                    服务端
  │                        │
  ├───── 请求 A ─────────► │
  │    (网络中断)           │ 执行 A
  │                        │
  ├───── 请求 A (重试) ──►  │ ?
  │                        │

问题：服务端如何判断这是重试，而不是新请求？
答案：Idempotency Key
```

---

## 幂等键实现

```go
package idempotency

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "time"

 "github.com/google/uuid"
)

// Key 幂等键
type Key struct {
 ID        string    // 唯一标识
 Resource  string    // 资源类型
 Operation string    // 操作类型
 CreatedAt time.Time // 创建时间
 ExpiresAt time.Time // 过期时间
}

// GenerateKey 生成幂等键
func GenerateKey(resource, operation string, payload interface{}) (*Key, error) {
 // 基于资源+操作+载荷生成确定性键
 data, err := json.Marshal(payload)
 if err != nil {
  return nil, err
 }

 h := sha256.New()
 h.Write([]byte(resource))
 h.Write([]byte(operation))
 h.Write(data)
 hash := hex.EncodeToString(h.Sum(nil))

 return &Key{
  ID:        hash[:32],
  Resource:  resource,
  Operation: operation,
  CreatedAt: time.Now(),
  ExpiresAt: time.Now().Add(24 * time.Hour),
 }, nil
}

// Store 幂等存储接口
type Store interface {
 // 获取或创建
 GetOrCreate(ctx context.Context, key *Key) (*Lock, error)

 // 更新结果
 UpdateResult(ctx context.Context, keyID string, result *Result) error

 // 获取结果
 GetResult(ctx context.Context, keyID string) (*Result, error)

 // 清理过期
 Cleanup(ctx context.Context, before time.Time) error
}

// Lock 分布式锁
type Lock struct {
 KeyID     string
 Status    LockStatus
 Acquired  bool
 Release   func()
}

type LockStatus int

const (
 LockStatusNew       LockStatus = iota // 新请求
 LockStatusInProgress                  // 处理中
 LockStatusCompleted                   // 已完成
)

// Result 执行结果
type Result struct {
 Status     ResultStatus
 Data       interface{}
 Error      string
 CreatedAt  time.Time
}

type ResultStatus int

const (
 ResultStatusSuccess ResultStatus = iota
 ResultStatusError
 ResultStatusPending
)

// Guard 幂等守卫
type Guard struct {
 store Store
}

// Execute 执行幂等操作
func (g *Guard) Execute(ctx context.Context, key *Key,
 fn func() (interface{}, error)) (interface{}, error) {

 // 1. 尝试获取锁
 lock, err := g.store.GetOrCreate(ctx, key)
 if err != nil {
  return nil, err
 }
 defer lock.Release()

 switch lock.Status {
 case LockStatusCompleted:
  // 已处理过，直接返回缓存结果
  result, err := g.store.GetResult(ctx, key.ID)
  if err != nil {
   return nil, err
  }

  if result.Status == ResultStatusError {
   return nil, fmt.Errorf(result.Error)
  }
  return result.Data, nil

 case LockStatusInProgress:
  // 正在处理，等待或返回
  return g.waitForResult(ctx, key.ID)

 case LockStatusNew:
  // 新请求，执行
  return g.doExecute(ctx, key.ID, fn)
 }

 return nil, fmt.Errorf("unknown lock status")
}

// doExecute 实际执行
func (g *Guard) doExecute(ctx context.Context, keyID string,
 fn func() (interface{}, error)) (interface{}, error) {

 data, err := fn()

 // 保存结果
 result := &Result{
  CreatedAt: time.Now(),
 }

 if err != nil {
  result.Status = ResultStatusError
  result.Error = err.Error()
 } else {
  result.Status = ResultStatusSuccess
  result.Data = data
 }

 g.store.UpdateResult(ctx, keyID, result)

 return data, err
}

// waitForResult 等待结果
func (g *Guard) waitForResult(ctx context.Context, keyID string) (interface{}, error) {
 ticker := time.NewTicker(100 * time.Millisecond)
 defer ticker.Stop()

 for {
  select {
  case <-ctx.Done():
   return nil, ctx.Err()
  case <-ticker.C:
   result, err := g.store.GetResult(context.Background(), keyID)
   if err != nil {
    continue
   }

   switch result.Status {
   case ResultStatusSuccess:
    return result.Data, nil
   case ResultStatusError:
    return nil, fmt.Errorf(result.Error)
   }
  }
 }
}
```

---

## Redis 存储实现

```go
// RedisStore Redis幂等存储
type RedisStore struct {
 client *redis.Client
 ttl    time.Duration
}

// GetOrCreate 获取或创建锁
func (r *RedisStore) GetOrCreate(ctx context.Context, key *Key) (*Lock, error) {
 lockKey := fmt.Sprintf("idempotency:%s:lock", key.ID)
 dataKey := fmt.Sprintf("idempotency:%s:data", key.ID)

 // 尝试获取锁
 acquired, err := r.client.SetNX(ctx, lockKey, "in_progress", r.ttl).Result()
 if err != nil {
  return nil, err
 }

 if acquired {
  // 新请求
  return &Lock{
   KeyID:    key.ID,
   Status:   LockStatusNew,
   Acquired: true,
   Release: func() {
    // 不立即释放，由 UpdateResult 处理
   },
  }, nil
 }

 // 检查是否已有结果
 exists, err := r.client.Exists(ctx, dataKey).Result()
 if err != nil {
  return nil, err
 }

 if exists > 0 {
  return &Lock{
   KeyID:    key.ID,
   Status:   LockStatusCompleted,
   Acquired: false,
   Release:  func() {},
  }, nil
 }

 // 正在处理中
 return &Lock{
  KeyID:    key.ID,
  Status:   LockStatusInProgress,
  Acquired: false,
  Release:  func() {},
 }, nil
}

// UpdateResult 更新结果
func (r *RedisStore) UpdateResult(ctx context.Context, keyID string, result *Result) error {
 dataKey := fmt.Sprintf("idempotency:%s:data", keyID)

 data, _ := json.Marshal(result)

 pipe := r.client.Pipeline()
 pipe.Set(ctx, dataKey, data, r.ttl)
 pipe.Del(ctx, fmt.Sprintf("idempotency:%s:lock", keyID))

 _, err := pipe.Exec(ctx)
 return err
}
```

---

## Exactly-Once 语义

```go
// ExactlyOnceProcessor 精确一次处理器
type ExactlyOnceProcessor struct {
 idempotency *Guard
 producer    KafkaProducer
 consumer    KafkaConsumer
}

// Process 精确一次处理
func (p *ExactlyOnceProcessor) Process(msg *kafka.Message) error {
 // 生成幂等键（基于消息 offset + partition）
 key := &Key{
  ID:        fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset),
  Resource:  msg.Topic,
  Operation: "process",
  CreatedAt: time.Now(),
 }

 _, err := p.idempotency.Execute(context.Background(), key, func() (interface{}, error) {
  // 业务逻辑
  result, err := p.doBusinessLogic(msg)
  if err != nil {
   return nil, err
  }

  // 发送结果（事务性）
  if err := p.producer.Send(result); err != nil {
   return nil, err
  }

  return result, nil
 })

 return err
}
```
