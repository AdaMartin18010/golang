# Redis 客户端

> **分类**: 开源技术堆栈

---

## go-redis

```go
import "github.com/redis/go-redis/v9"

rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // 无密码
    DB:       0,  // 默认 DB
})
```

---

## 基本操作

```go
ctx := context.Background()

// String
rdb.Set(ctx, "key", "value", 0)
val, err := rdb.Get(ctx, "key").Result()

// Hash
rdb.HSet(ctx, "user:1", "name", "Alice")
rdb.HGet(ctx, "user:1", "name")

// List
rdb.LPush(ctx, "queue", "task1")
rdb.RPop(ctx, "queue")

// Set
rdb.SAdd(ctx, "tags", "go", "redis")
rdb.SMembers(ctx, "tags")
```

---

## 连接池

```go
rdb := redis.NewClient(&redis.Options{
    PoolSize:     10,
    MinIdleConns: 5,
    MaxRetries:   3,
})
```

---

## 分布式锁

```go
lock := rdb.Lock(ctx, "my-lock", 10*time.Second)
if err := lock.Err(); err != nil {
    return err
}
defer lock.Release(ctx)
```
