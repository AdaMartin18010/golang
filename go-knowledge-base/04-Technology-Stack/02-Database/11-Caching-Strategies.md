# 缓存策略 (Caching Strategies)

> **分类**: 开源技术堆栈  
> **标签**: #cache #redis #performance

---

## 缓存模式

### Cache-Aside (旁路缓存)

```go
type CacheAside struct {
    cache cache.Cache
    db    *sql.DB
}

func (ca *CacheAside) Get(ctx context.Context, key string) (*User, error) {
    // 1. 先查缓存
    if val, err := ca.cache.Get(ctx, key); err == nil {
        return val.(*User), nil
    }
    
    // 2. 缓存未命中，查数据库
    user, err := ca.getFromDB(ctx, key)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    ca.cache.Set(ctx, key, user, 10*time.Minute)
    
    return user, nil
}

func (ca *CacheAside) Set(ctx context.Context, user *User) error {
    // 1. 更新数据库
    if err := ca.saveToDB(ctx, user); err != nil {
        return err
    }
    
    // 2. 删除缓存（或更新）
    ca.cache.Delete(ctx, user.ID)
    
    return nil
}
```

### Read-Through (穿透缓存)

```go
func (c *ReadThroughCache) Get(ctx context.Context, key string) (*User, error) {
    return c.cache.GetOrSet(ctx, key, func() (interface{}, error) {
        // 缓存未命中时自动加载
        return c.db.GetUser(ctx, key)
    }, 10*time.Minute)
}
```

### Write-Through (直写缓存)

```go
func (c *WriteThroughCache) Set(ctx context.Context, user *User) error {
    // 同步更新缓存和数据库
    if err := c.cache.Set(ctx, user.ID, user, 10*time.Minute); err != nil {
        return err
    }
    
    return c.db.SaveUser(ctx, user)
}
```

### Write-Behind (异步写)

```go
func (c *WriteBehindCache) Set(ctx context.Context, user *User) error {
    // 先写缓存
    if err := c.cache.Set(ctx, user.ID, user, 10*time.Minute); err != nil {
        return err
    }
    
    // 异步写数据库
    c.writeQueue <- user
    
    return nil
}

func (c *WriteBehindCache) flush() {
    for user := range c.writeQueue {
        // 批量写入
        c.batch <- user
        
        if len(c.batch) >= c.batchSize {
            c.flushBatch()
        }
    }
}
```

---

## 缓存问题

### 缓存穿透

```go
func (ca *CacheAside) GetWithBloomFilter(ctx context.Context, key string) (*User, error) {
    // 布隆过滤器检查
    if !ca.bloomFilter.MayContain(key) {
        return nil, ErrNotFound  // 直接返回，不查缓存和数据库
    }
    
    return ca.Get(ctx, key)
}

// 或缓存空值
func (ca *CacheAside) GetWithNullCache(ctx context.Context, key string) (*User, error) {
    val, err := ca.cache.Get(ctx, key)
    if err == nil {
        if val == nil {
            return nil, ErrNotFound  // 缓存的空值
        }
        return val.(*User), nil
    }
    
    user, err := ca.getFromDB(ctx, key)
    if err == ErrNotFound {
        // 缓存空值，防止穿透
        ca.cache.Set(ctx, key, nil, 5*time.Minute)
        return nil, ErrNotFound
    }
    
    ca.cache.Set(ctx, key, user, 10*time.Minute)
    return user, nil
}
```

### 缓存击穿

```go
func (ca *CacheAside) GetWithMutex(ctx context.Context, key string) (*User, error) {
    // 先查缓存
    if val, err := ca.cache.Get(ctx, key); err == nil {
        return val.(*User), nil
    }
    
    // 获取分布式锁
    lockKey := "lock:" + key
    if !ca.cache.Lock(ctx, lockKey, 5*time.Second) {
        // 没拿到锁，等待后重试
        time.Sleep(100 * time.Millisecond)
        return ca.GetWithMutex(ctx, key)
    }
    defer ca.cache.Unlock(ctx, lockKey)
    
    // 双重检查
    if val, err := ca.cache.Get(ctx, key); err == nil {
        return val.(*User), nil
    }
    
    // 查数据库并缓存
    user, err := ca.getFromDB(ctx, key)
    if err != nil {
        return nil, err
    }
    
    ca.cache.Set(ctx, key, user, 10*time.Minute)
    return user, nil
}
```

### 缓存雪崩

```go
func (ca *CacheAside) SetWithJitter(ctx context.Context, key string, user *User) error {
    // 随机过期时间，避免同时失效
    jitter := time.Duration(rand.Intn(300)) * time.Second
    ttl := 10*time.Minute + jitter
    
    return ca.cache.Set(ctx, key, user, ttl)
}
```

---

## 缓存一致性

```go
// 最终一致性：延迟双删
func (ca *CacheAside) UpdateWithDelayDelete(ctx context.Context, user *User) error {
    // 1. 删除缓存
    ca.cache.Delete(ctx, user.ID)
    
    // 2. 更新数据库
    if err := ca.saveToDB(ctx, user); err != nil {
        return err
    }
    
    // 3. 延迟删除缓存（处理并发读）
    go func() {
        time.Sleep(500 * time.Millisecond)
        ca.cache.Delete(ctx, user.ID)
    }()
    
    return nil
}
```
