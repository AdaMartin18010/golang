# Redis 客户端封装

Redis 客户端封装，提供连接管理和常用操作。

## 功能特性

- ✅ 连接池管理
- ✅ 连接参数配置
- ✅ 连接健康检查
- ✅ 常用操作封装
- ✅ 自动重连支持

## 使用示例

### 基本使用

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/cache/redis"
)

// 使用默认配置
client, err := redis.NewClient(redis.DefaultConfig())
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### 自定义配置

```go
config := redis.Config{
    Addr:         "localhost:6379",
    Password:     "password",
    DB:           0,
    PoolSize:     10,
    MinIdleConns: 5,
}

client, err := redis.NewClient(config)
```

### 基本操作

```go
ctx := context.Background()

// 设置值
err := client.Set(ctx, "key", "value", time.Hour)

// 获取值
val, err := client.Get(ctx, "key")

// 删除键
err := client.Del(ctx, "key1", "key2")

// 检查键是否存在
count, err := client.Exists(ctx, "key")

// 设置过期时间
err := client.Expire(ctx, "key", time.Hour)
```

### 与限流中间件集成

```go
import (
    "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
    "github.com/yourusername/golang/internal/infrastructure/cache/redis"
)

// 创建 Redis 客户端
redisClient, _ := redis.NewClient(redis.DefaultConfig())

// 使用 Redis 分布式限流
r.Use(middleware.RateLimitMiddleware(middleware.RateLimitConfig{
    RequestsPerSecond: 100,
    Window:            time.Second,
    Algorithm:         middleware.AlgorithmSlidingWindow,
    RedisClient:       middleware.NewRedisAdapter(redisClient.GetClient()),
    RedisKeyPrefix:    "ratelimit:api",
}))
```

## 相关资源

- [go-redis 文档](https://github.com/redis/go-redis)
- [Redis 官方文档](https://redis.io/docs/)
