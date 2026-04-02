# 幂等性设计 (Idempotency)

> **分类**: 成熟应用领域  
> **标签**: #idempotency #distributed-systems

---

## 幂等键模式

```go
type IdempotentHandler struct {
    store IdempotencyStore
}

type IdempotencyStore interface {
    Get(key string) (*IdempotencyRecord, error)
    Save(record *IdempotencyRecord) error
}

type IdempotencyRecord struct {
    Key        string
    Status     string  // processing, completed
    Response   []byte
    ExpiresAt  time.Time
}

func (h *IdempotentHandler) Handle(ctx context.Context, req Request) (Response, error) {
    // 1. 检查幂等键
    record, err := h.store.Get(req.IdempotencyKey)
    if err == nil && record != nil {
        if record.Status == "completed" {
            // 返回缓存结果
            var resp Response
            json.Unmarshal(record.Response, &resp)
            return resp, nil
        }
        if record.Status == "processing" {
            return Response{}, ErrProcessing
        }
    }
    
    // 2. 标记为处理中
    h.store.Save(&IdempotencyRecord{
        Key:       req.IdempotencyKey,
        Status:    "processing",
        ExpiresAt: time.Now().Add(24 * time.Hour),
    })
    
    // 3. 执行业务逻辑
    resp, err := h.process(ctx, req)
    
    // 4. 保存结果
    respBytes, _ := json.Marshal(resp)
    h.store.Save(&IdempotencyRecord{
        Key:       req.IdempotencyKey,
        Status:    "completed",
        Response:  respBytes,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    })
    
    return resp, err
}
```

---

## Redis 实现

```go
type RedisIdempotencyStore struct {
    client *redis.Client
}

func (s *RedisIdempotencyStore) Get(key string) (*IdempotencyRecord, error) {
    data, err := s.client.Get(ctx, "idempotency:"+key).Result()
    if err == redis.Nil {
        return nil, ErrNotFound
    }
    
    var record IdempotencyRecord
    json.Unmarshal([]byte(data), &record)
    return &record, nil
}

func (s *RedisIdempotencyStore) Save(record *IdempotencyRecord) error {
    data, _ := json.Marshal(record)
    return s.client.Set(ctx, "idempotency:"+record.Key, data, 24*time.Hour).Err()
}
```

---

## 乐观锁幂等

```go
func UpdateBalance(ctx context.Context, userID string, amount int64) error {
    for {
        // 读取当前版本
        user, err := db.GetUser(ctx, userID)
        if err != nil {
            return err
        }
        
        // 尝试更新
        result, err := db.ExecContext(ctx,
            "UPDATE users SET balance = balance + ?, version = version + 1 WHERE id = ? AND version = ?",
            amount, userID, user.Version,
        )
        
        affected, _ := result.RowsAffected()
        if affected == 1 {
            return nil  // 成功
        }
        
        // 版本冲突，重试
        time.Sleep(time.Millisecond * 10)
    }
}
```

---

## HTTP 幂等

```go
func IdempotentMiddleware(store IdempotencyStore) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 幂等键从 Header 获取
        key := c.GetHeader("Idempotency-Key")
        if key == "" {
            c.Next()
            return
        }
        
        // 检查是否处理过
        record, _ := store.Get(key)
        if record != nil && record.Status == "completed" {
            c.Data(http.StatusOK, "application/json", record.Response)
            c.Abort()
            return
        }
        
        // 包装 ResponseWriter 捕获响应
        w := &responseRecorder{ResponseWriter: c.Writer}
        c.Writer = w
        
        c.Next()
        
        // 保存响应
        if c.Writer.Status() == http.StatusOK {
            store.Save(&IdempotencyRecord{
                Key:      key,
                Status:   "completed",
                Response: w.body.Bytes(),
            })
        }
    }
}
```

---

## 幂等键生成策略

| 场景 | 策略 | 示例 |
|------|------|------|
| 客户端生成 | UUID | `Idempotency-Key: <uuid>` |
| 服务端生成 | 哈希 | SHA256(user + action + params) |
| 自然键 | 业务键 | order_id + action |
