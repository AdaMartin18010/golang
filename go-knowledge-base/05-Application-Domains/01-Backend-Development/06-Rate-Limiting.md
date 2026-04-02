# 限流算法 (Rate Limiting)

> **分类**: 成熟应用领域

---

## 令牌桶算法

```go
type TokenBucket struct {
    rate       float64    // 每秒产生令牌数
    capacity   int        // 桶容量
    tokens     float64    // 当前令牌数
    lastUpdate time.Time  // 上次更新时间
    mu         sync.Mutex
}

func NewTokenBucket(rate float64, capacity int) *TokenBucket {
    return &TokenBucket{
        rate:       rate,
        capacity:   capacity,
        tokens:     float64(capacity),
        lastUpdate: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(tb.lastUpdate).Seconds()
    tb.lastUpdate = now
    
    // 添加新令牌
    tb.tokens += elapsed * tb.rate
    if tb.tokens > float64(tb.capacity) {
        tb.tokens = float64(tb.capacity)
    }
    
    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }
    return false
}
```

---

## 漏桶算法

```go
type LeakyBucket struct {
    rate       time.Duration  // 漏水间隔
    capacity   int            // 桶容量
    water      int            // 当前水量
    lastLeak   time.Time
    mu         sync.Mutex
}

func (lb *LeakyBucket) Allow() bool {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(lb.lastLeak)
    
    // 漏水
    leaked := int(elapsed / lb.rate)
    lb.water -= leaked
    if lb.water < 0 {
        lb.water = 0
    }
    lb.lastLeak = now
    
    if lb.water < lb.capacity {
        lb.water++
        return true
    }
    return false
}
```

---

## Gin 中间件限流

```go
func RateLimitMiddleware(tb *TokenBucket) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !tb.Allow() {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "rate limit exceeded",
            })
            return
        }
        c.Next()
    }
}

// 使用
r.Use(RateLimitMiddleware(NewTokenBucket(10, 100)))
```

---

## Redis 分布式限流

```go
func RedisRateLimiter(ctx context.Context, client *redis.Client, key string, limit int, window time.Duration) (bool, error) {
    pipe := client.Pipeline()
    now := time.Now().Unix()
    
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(window.Seconds())))
    pipe.ZCard(ctx, key)
    pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
    pipe.Expire(ctx, key, window)
    
    cmders, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    
    count := cmders[1].(*redis.IntCmd).Val()
    return count < int64(limit), nil
}
```
