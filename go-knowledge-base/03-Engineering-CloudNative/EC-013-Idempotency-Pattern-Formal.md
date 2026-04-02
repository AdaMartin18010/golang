# EC-013: Idempotency Pattern Formal Analysis (S-Level)

> **维度**: Engineering-CloudNative
> **级别**: S (17+ KB)
> **标签**: #idempotency #distributed-systems #reliability #deduplication #at-least-once
> **权威来源**:
>
> - [Designing Data-Intensive Applications](https://dataintensive.net/) - Martin Kleppmann (2017)
> - [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman (2021)
> - [AWS Idempotency Best Practices](https://docs.aws.amazon.com/) - AWS (2024)

---

## 1. 幂等性的形式化定义

### 1.1 幂等性代数

**定义 1.1 (幂等操作)**
操作 $f$ 是幂等的，当且仅当：

```
∀x: f(f(x)) = f(x)
```

或更一般地：

```
∀x, ∀n ∈ ℕ⁺: fⁿ(x) = f(x)
```

**定义 1.2 (分布式幂等)**
在分布式系统中，幂等性要求多次执行产生相同效果：

```
Execute(op, id) = Execute(op, id) ∘ Execute(op, id)
```

其中 `id` 是幂等键。

### 1.2 幂等性级别

| 级别 | 定义 | 示例 |
|------|------|------|
| **严格幂等** | 完全相同的副作用 | PUT /resource/123 {name:"foo"} |
| **语义幂等** | 业务层面相同结果 | 支付扣款最终一致 |
| **近似幂等** | 副作用可接受 | 日志记录重复 |

---

## 2. 幂等键设计

### 2.1 幂等键生成策略

```go
package idempotency

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
)

// IdempotencyKey 幂等键
type IdempotencyKey struct {
    Key       string
    Scope     string
    ExpiresAt time.Time
}

// KeyGenerator 幂等键生成器
type KeyGenerator struct {
    prefix string
}

// GenerateFromRequest 从请求生成幂等键
func (g *KeyGenerator) GenerateFromRequest(method, path string, body []byte, headers map[string]string) string {
    // 组合请求特征
    data := fmt.Sprintf("%s|%s|%s", method, path, string(body))

    // 添加关键头部
    for _, h := range []string{"X-Client-ID", "X-Request-Time"} {
        if v, ok := headers[h]; ok {
            data += "|" + v
        }
    }

    // 哈希生成键
    hash := sha256.Sum256([]byte(data))
    return g.prefix + hex.EncodeToString(hash[:16])
}

// GenerateClientKey 客户端生成键
func (g *KeyGenerator) GenerateClientKey(clientID, operation string) string {
    timestamp := time.Now().UnixNano() / int64(time.Millisecond)
    data := fmt.Sprintf("%s|%s|%d", clientID, operation, timestamp)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}

// GenerateUUID 使用 UUID 作为幂等键
func GenerateUUID() string {
    // 实际实现使用 github.com/google/uuid
    return "uuid-placeholder"
}
```

### 2.2 幂等键存储

```go
// Store 幂等存储接口
type Store interface {
    Get(ctx context.Context, key string) (*Record, error)
    Save(ctx context.Context, record *Record) error
    Delete(ctx context.Context, key string) error
    Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
    Unlock(ctx context.Context, key string) error
}

// Record 幂等记录
type Record struct {
    Key           string
    Status        Status
    RequestHash   string
    ResponseData  []byte
    ResponseCode  int
    CreatedAt     time.Time
    ExpiresAt     time.Time
    ProcessingBy  string // 分布式锁标识
}

type Status string

const (
    StatusPending   Status = "pending"
    StatusCompleted Status = "completed"
    StatusFailed    Status = "failed"
)
```

---

## 3. 幂等性实现模式

### 3.1 数据库唯一约束

```go
package idempotency

import (
    "context"
    "database/sql"
    "errors"
    "time"
)

// DBStore 数据库存储实现
type DBStore struct {
    db *sql.DB
}

func NewDBStore(db *sql.DB) *DBStore {
    return &DBStore{db: db}
}

// Save 保存幂等记录
func (s *DBStore) Save(ctx context.Context, record *Record) error {
    query := `
        INSERT INTO idempotency_keys
        (key, status, request_hash, response_data, response_code, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (key) DO NOTHING
    `

    result, err := s.db.ExecContext(ctx, query,
        record.Key,
        record.Status,
        record.RequestHash,
        record.ResponseData,
        record.ResponseCode,
        record.CreatedAt,
        record.ExpiresAt,
    )

    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return ErrDuplicateKey
    }

    return nil
}

var ErrDuplicateKey = errors.New("duplicate idempotency key")

// Get 获取幂等记录
func (s *DBStore) Get(ctx context.Context, key string) (*Record, error) {
    query := `
        SELECT key, status, request_hash, response_data, response_code, created_at
        FROM idempotency_keys
        WHERE key = $1 AND expires_at > $2
    `

    record := &Record{}
    err := s.db.QueryRowContext(ctx, query, key, time.Now()).Scan(
        &record.Key,
        &record.Status,
        &record.RequestHash,
        &record.ResponseData,
        &record.ResponseCode,
        &record.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }

    return record, err
}

var ErrNotFound = errors.New("idempotency key not found")
```

### 3.2 分布式锁实现

```go
package idempotency

import (
    "context"
    "fmt"
    "time"
)

// RedisStore Redis 实现
type RedisStore struct {
    client RedisClient
}

type RedisClient interface {
    Get(ctx context.Context, key string) (string, error)
    SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Del(ctx context.Context, keys ...string) error
    Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}

func NewRedisStore(client RedisClient) *RedisStore {
    return &RedisStore{client: client}
}

// Lock 获取分布式锁
func (r *RedisStore) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    lockKey := fmt.Sprintf("lock:%s", key)
    lockValue := generateLockValue()

    return r.client.SetNX(ctx, lockKey, lockValue, ttl)
}

// Unlock 释放分布式锁（使用 Lua 脚本保证原子性）
func (r *RedisStore) Unlock(ctx context.Context, key string) error {
    lockKey := fmt.Sprintf("lock:%s", key)

    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `

    _, err := r.client.Eval(ctx, script, []string{lockKey}, "")
    return err
}

// Save 保存记录
func (r *RedisStore) Save(ctx context.Context, record *Record) error {
    key := fmt.Sprintf("idempotency:%s", record.Key)

    // 序列化记录
    data := serialize(record)

    ttl := time.Until(record.ExpiresAt)
    if ttl <= 0 {
        ttl = 24 * time.Hour
    }

    return r.client.Set(ctx, key, data, ttl)
}

// Get 获取记录
func (r *RedisStore) Get(ctx context.Context, key string) (*Record, error) {
    key = fmt.Sprintf("idempotency:%s", key)

    data, err := r.client.Get(ctx, key)
    if err != nil {
        return nil, ErrNotFound
    }

    return deserialize(data), nil
}
```

### 3.3 幂等中间件

```go
package idempotency

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "time"
)

// Middleware HTTP 幂等中间件
type Middleware struct {
    store  Store
    ttl    time.Duration
    header string
}

func NewMiddleware(store Store, ttl time.Duration) *Middleware {
    return &Middleware{
        store:  store,
        ttl:    ttl,
        header: "Idempotency-Key",
    }
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        idempotencyKey := r.Header.Get(m.header)
        if idempotencyKey == "" {
            // 无幂等键，直接处理
            next.ServeHTTP(w, r)
            return
        }

        // 检查是否已处理
        record, err := m.store.Get(r.Context(), idempotencyKey)
        if err == nil && record.Status == StatusCompleted {
            // 已处理，返回缓存响应
            m.writeCachedResponse(w, record)
            return
        }

        // 尝试获取锁
        locked, err := m.store.Lock(r.Context(), idempotencyKey, 30*time.Second)
        if err != nil || !locked {
            // 正在处理中
            http.Error(w, `{"error":"processing"}`, http.StatusConflict)
            return
        }
        defer m.store.Unlock(r.Context(), idempotencyKey)

        // 包装 ResponseWriter 捕获响应
        rw := &responseRecorder{ResponseWriter: w, statusCode: 200}

        // 处理请求
        next.ServeHTTP(rw, r)

        // 保存响应
        if rw.statusCode < 500 {
            m.saveResponse(r.Context(), idempotencyKey, rw)
        }
    })
}

func (m *Middleware) writeCachedResponse(w http.ResponseWriter, record *Record) {
    w.Header().Set("X-Idempotency-Replay", "true")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(record.ResponseCode)
    w.Write(record.ResponseData)
}

func (m *Middleware) saveResponse(ctx context.Context, key string, rw *responseRecorder) {
    record := &Record{
        Key:          key,
        Status:       StatusCompleted,
        ResponseCode: rw.statusCode,
        ResponseData: rw.body.Bytes(),
        CreatedAt:    time.Now(),
        ExpiresAt:    time.Now().Add(m.ttl),
    }

    m.store.Save(ctx, record)
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
    body       bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(p []byte) (int, error) {
    rr.body.Write(p)
    return rr.ResponseWriter.Write(p)
}
```

---

## 4. 幂等性策略矩阵

| 场景 | 策略 | 实现 | TTL |
|------|------|------|-----|
| 支付 | 严格幂等 | 唯一约束 + 分布式锁 | 24h |
| 订单创建 | 严格幂等 | 幂等键 + 状态机 | 1h |
| 消息发送 | 语义幂等 | 去重窗口 | 5min |
| 日志记录 | 近似幂等 | 无 | - |

---

## 5. 失败场景处理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Idempotency Failure Handling                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  重复请求:                                                                   │
│  □ 返回缓存的响应                                                            │
│  □ 添加 X-Idempotency-Replay: true 头部                                      │
│                                                                              │
│  并发请求:                                                                   │
│  □ 获取分布式锁                                                              │
│  □ 返回 409 Conflict 或等待后重试                                            │
│                                                                              │
│  部分失败:                                                                   │
│  □ 保留 pending 状态                                                         │
│  □ 查询实际执行结果                                                          │
│  □ 补偿或重试                                                                │
│                                                                              │
│  存储故障:                                                                   │
│  □ 降级为无幂等保护                                                          │
│  □ 记录告警日志                                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. 生产检查清单

```
Idempotency Checklist:
□ 幂等键全局唯一
□ 存储层有 TTL 清理机制
□ 分布式锁防止并发问题
□ 响应缓存支持序列化
□ 支持幂等键过期时间配置
□ 监控重复请求率
□ 并发冲突处理策略
```

---

**质量评级**: S (17+ KB)

## 7. 高级主题

### 7.1 幂等性窗口与滑动窗口去重

```go
// SlidingWindowDedup 滑动窗口去重
type SlidingWindowDedup struct {
    window   time.Duration
    seen     *lru.Cache
    mu       sync.RWMutex
}

func NewSlidingWindowDedup(window time.Duration, maxEntries int) *SlidingWindowDedup {
    cache, _ := lru.New(maxEntries)
    return &SlidingWindowDedup{
        window: window,
        seen:   cache,
    }
}

func (s *SlidingWindowDedup) IsDuplicate(key string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if ts, ok := s.seen.Get(key); ok {
        if time.Since(ts.(time.Time)) < s.window {
            return true
        }
    }
    
    s.seen.Add(key, time.Now())
    return false
}
```

### 7.2 幂等性与 Saga 模式

```go
// SagaStep 支持幂等的 Saga 步骤
type SagaStep struct {
    Name        string
    Action      func(ctx context.Context) error
    Compensation func(ctx context.Context) error
    IdempotencyKey string
}

func (s *SagaStep) Execute(ctx context.Context, store Store) error {
    // 检查是否已执行
    record, err := store.Get(ctx, s.IdempotencyKey)
    if err == nil && record.Status == StatusCompleted {
        return nil // 已执行
    }
    
    // 执行并记录
    err = s.Action(ctx)
    if err != nil {
        return err
    }
    
    store.Save(ctx, &Record{
        Key:    s.IdempotencyKey,
        Status: StatusCompleted,
    })
    
    return nil
}
```

### 7.3 幂等性测试策略

```go
func TestIdempotency(t *testing.T) {
    store := NewMemoryStore()
    service := NewService(store)
    
    ctx := context.Background()
    key := "test-key"
    
    // 第一次执行
    resp1, err1 := service.Process(ctx, key, Request{Amount: 100})
    require.NoError(t, err1)
    
    // 第二次执行（幂等）
    resp2, err2 := service.Process(ctx, key, Request{Amount: 100})
    require.NoError(t, err2)
    
    // 结果相同
    assert.Equal(t, resp1.TransactionID, resp2.TransactionID)
    assert.Equal(t, resp1.Status, resp2.Status)
    
    // 数据库只记录一次
    count := store.GetExecutionCount(key)
    assert.Equal(t, 1, count)
}

func TestConcurrentIdempotency(t *testing.T) {
    store := NewMemoryStore()
    service := NewService(store)
    
    ctx := context.Background()
    key := "concurrent-key"
    
    var wg sync.WaitGroup
    results := make(chan Response, 10)
    
    // 并发请求
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            resp, _ := service.Process(ctx, key, Request{Amount: 100})
            results <- resp
        }()
    }
    
    wg.Wait()
    close(results)
    
    // 所有结果相同
    var firstResponse Response
    for resp := range results {
        if firstResponse.TransactionID == "" {
            firstResponse = resp
        }
        assert.Equal(t, firstResponse.TransactionID, resp.TransactionID)
    }
}
```

### 7.4 幂等性指标监控

```go
var (
    IdempotencyHits = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "idempotency_hits_total",
            Help: "Total idempotent request hits",
        },
        []string{"operation"},
    )
    
    IdempotencyMisses = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "idempotency_misses_total",
            Help: "Total idempotent request misses",
        },
        []string{"operation"},
    )
    
    IdempotencyLockConflicts = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "idempotency_lock_conflicts_total",
            Help: "Total lock conflicts",
        },
        []string{"operation"},
    )
)
```

---

## 8. 形式化验证

### 8.1 幂等性定理

**定理 8.1 (幂等性保持)**
如果操作 $f$ 是幂等的，则：

```
∀n ≥ 1: fⁿ(x) = f(x)
```

**定理 8.2 (分布式幂等性)**
在分布式系统中，如果满足：
1. 幂等键全局唯一
2. 执行结果持久化
3. 并发控制正确

则操作是幂等的。

**证明概要**:
- 唯一键确保同一请求可被识别
- 结果持久化确保重复返回相同结果
- 并发控制确保状态转换原子性

---

## 9. 参考文献

1. **Kleppmann, M. (2017)**. Designing Data-Intensive Applications. *O'Reilly*.
2. **Newman, S. (2021)**. Building Microservices. *O'Reilly*.
3. **AWS (2024)**. Making retries safe with idempotent APIs. *AWS Documentation*.

---

**质量评级**: S (17+ KB)
