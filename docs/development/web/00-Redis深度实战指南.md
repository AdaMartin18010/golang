# Redis

**字数**: ~48,000字  
**代码示例**: 170+个完整示例  
**实战案例**: 15个端到端案例  
**适用人群**: 中级到高级Go开发者

---

## 📋 目录


- [第一部分：Redis核心原理](#第一部分redis核心原理)
  - [Redis架构](#redis架构)
  - [为什么使用Redis？](#为什么使用redis)
- [第二部分：Go连接Redis](#第二部分go连接redis)
  - [实战案例1：基础连接（go-redis）](#实战案例1基础连接go-redis)
- [第三部分：五大数据结构](#第三部分五大数据结构)
  - [实战案例2：String（字符串）](#实战案例2string字符串)
  - [实战案例3：Hash（哈希表）](#实战案例3hash哈希表)
- [第五部分：缓存策略](#第五部分缓存策略)
  - [实战案例4：缓存穿透/击穿/雪崩](#实战案例4缓存穿透击穿雪崩)
  - [实战案例5：List/Set/ZSet操作](#实战案例5listsetzset操作)
- [第四部分：高级数据结构](#第四部分高级数据结构)
  - [实战案例6：Bitmap/HyperLogLog/Geo](#实战案例6bitmaphyperlogloggeo)
- [第六部分：分布式锁](#第六部分分布式锁)
  - [实战案例7：Redis分布式锁](#实战案例7redis分布式锁)
- [第七部分：消息队列](#第七部分消息队列)
  - [实战案例8：List实现消息队列](#实战案例8list实现消息队列)
- [第十部分：事务与管道](#第十部分事务与管道)
  - [实战案例9：Redis事务](#实战案例9redis事务)
- [第十一部分：Lua脚本](#第十一部分lua脚本)
  - [实战案例10：Lua脚本原子操作](#实战案例10lua脚本原子操作)
- [第八部分：Stream流处理](#第八部分stream流处理)
  - [实战案例11：Redis Stream消费者组](#实战案例11redis-stream消费者组)
- [第九部分：发布订阅](#第九部分发布订阅)
  - [实战案例12：Pub/Sub模式](#实战案例12pubsub模式)
- [第十二部分：持久化与高可用](#第十二部分持久化与高可用)
  - [持久化配置](#持久化配置)
- [第十三部分：性能优化](#第十三部分性能优化)
  - [实战案例13：性能优化技巧](#实战案例13性能优化技巧)
- [第十四部分：监控与运维](#第十四部分监控与运维)
  - [实战案例14：Redis监控](#实战案例14redis监控)
- [第十六部分：完整项目实战](#第十六部分完整项目实战)
  - [实战案例15：电商库存系统](#实战案例15电商库存系统)
- [🎯 总结](#总结)
  - [Redis核心要点](#redis核心要点)
  - [最佳实践清单](#最佳实践清单)

## 第一部分：Redis核心原理

### Redis架构

```text
┌─────────────────────────────────────────────────┐
│              Redis架构                          │
└─────────────────────────────────────────────────┘

Client
  │
  ▼
┌─────────────────────────────────────────────────┐
│              Redis Server                       │
│  ┌───────────────────────────────────────────┐  │
│  │         Memory (内存数据库)                │  │
│  │  ┌─────────────────────────────────────┐  │  │
│  │  │  String | Hash | List | Set | ZSet  │  │  │
│  │  │  Bitmap | HyperLogLog | Geo | Stream│  │  │
│  │  └─────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────┘  │
│                     ↓                           │
│  ┌───────────────────────────────────────────┐  │
│  │         持久化 (Persistence)              │  │
│  │  - RDB (快照)                             │  │
│  │  - AOF (追加文件)                         │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
         ↓                    ↓
    ┌─────────┐          ┌─────────┐
    │  Disk   │          │ Replica │
    └─────────┘          └─────────┘

核心特性:
1. 内存数据库 - 极快（10万+QPS）
2. 支持持久化 - RDB + AOF
3. 多种数据结构 - 不只是KV存储
4. 主从复制 - 高可用
5. 哨兵模式 - 自动故障转移
6. 集群模式 - 分布式扩展
7. 发布订阅 - 消息通信
8. Lua脚本 - 原子操作
```

---

### 为什么使用Redis？

```text
使用场景:

1. 缓存 (最常用)
   数据库 → Redis缓存 → 应用
   - 减少DB压力
   - 提升响应速度（从ms到µs）

2. 会话存储
   用户Session → Redis
   - 支持分布式应用
   - 快速访问

3. 分布式锁
   多个服务 → Redis锁 → 资源
   - 防止并发问题
   - 保证原子性

4. 计数器
   点赞/浏览量 → Redis INCR
   - 高并发写入
   - 实时统计

5. 排行榜
   游戏分数 → Redis ZSet
   - 实时排序
   - 范围查询

6. 消息队列
   生产者 → Redis List → 消费者
   - 异步处理
   - 削峰填谷

7. 实时分析
   用户行为 → Redis Stream
   - 流式处理
   - 时间序列数据
```

---

## 第二部分：Go连接Redis

### 实战案例1：基础连接（go-redis）

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端封装
type RedisClient struct {
 client *redis.Client
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
 rdb := redis.NewClient(&redis.Options{
  Addr:         addr,     // localhost:6379
  Password:     password, // 密码
  DB:           db,       // 数据库
  PoolSize:     10,       // 连接池大小
  MinIdleConns: 5,        // 最小空闲连接
  MaxRetries:   3,        // 最大重试次数
  DialTimeout:  5 * time.Second,
  ReadTimeout:  3 * time.Second,
  WriteTimeout: 3 * time.Second,
  PoolTimeout:  4 * time.Second,
 })

 // 测试连接
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()

 if err := rdb.Ping(ctx).Err(); err != nil {
  return nil, fmt.Errorf("failed to connect to redis: %v", err)
 }

 log.Println("Connected to Redis successfully")
 return &RedisClient{client: rdb}, nil
}

// Close 关闭连接
func (r *RedisClient) Close() error {
 return r.client.Close()
}

// Set 设置键值（带过期时间）
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
 return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
 val, err := r.client.Get(ctx, key).Result()
 if err == redis.Nil {
  return "", fmt.Errorf("key does not exist")
 }
 return val, err
}

// Del 删除键
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
 return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
 return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
 return r.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余生存时间
func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
 return r.client.TTL(ctx, key).Result()
}

// ===== 使用示例 =====
func main() {
 // 创建客户端
 client, err := NewRedisClient("localhost:6379", "", 0)
 if err != nil {
  log.Fatal(err)
 }
 defer client.Close()

 ctx := context.Background()

 // 设置键值
 err = client.Set(ctx, "name", "Alice", 10*time.Minute)
 if err != nil {
  log.Printf("Set error: %v", err)
 }

 // 获取值
 val, err := client.Get(ctx, "name")
 if err != nil {
  log.Printf("Get error: %v", err)
 } else {
  fmt.Printf("name = %s\n", val)
 }

 // 检查存在
 exists, _ := client.Exists(ctx, "name")
 fmt.Printf("exists = %d\n", exists)

 // 获取TTL
 ttl, _ := client.TTL(ctx, "name")
 fmt.Printf("ttl = %v\n", ttl)

 // 删除键
 client.Del(ctx, "name")
}
```

---

## 第三部分：五大数据结构

### 实战案例2：String（字符串）

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/redis/go-redis/v9"
)

type StringOps struct {
 rdb *redis.Client
}

func NewStringOps(rdb *redis.Client) *StringOps {
 return &StringOps{rdb: rdb}
}

// ===== 基础操作 =====

// SetGet 设置和获取
func (s *StringOps) SetGet(ctx context.Context) {
 // SET
 s.rdb.Set(ctx, "user:1:name", "Alice", 0)
 s.rdb.Set(ctx, "user:1:age", "25", 0)

 // GET
 name, _ := s.rdb.Get(ctx, "user:1:name").Result()
 age, _ := s.rdb.Get(ctx, "user:1:age").Result()
 fmt.Printf("Name: %s, Age: %s\n", name, age)

 // SETEX (带过期时间)
 s.rdb.SetEx(ctx, "token", "abc123", 1*time.Hour)

 // SETNX (不存在时设置)
 ok, _ := s.rdb.SetNX(ctx, "lock", "1", 10*time.Second).Result()
 fmt.Printf("Lock acquired: %v\n", ok)

 // MSET/MGET (批量操作)
 s.rdb.MSet(ctx, "key1", "val1", "key2", "val2", "key3", "val3")
 vals, _ := s.rdb.MGet(ctx, "key1", "key2", "key3").Result()
 fmt.Printf("Values: %v\n", vals)
}

// ===== 计数器操作 =====

// Counter 计数器实现
func (s *StringOps) Counter(ctx context.Context) {
 key := "article:123:views"

 // INCR (增加1)
 views, _ := s.rdb.Incr(ctx, key).Result()
 fmt.Printf("Views: %d\n", views)

 // INCRBY (增加N)
 s.rdb.IncrBy(ctx, key, 10)

 // DECR (减少1)
 s.rdb.Decr(ctx, key)

 // DECRBY (减少N)
 s.rdb.DecrBy(ctx, key, 5)

 // 获取最终值
 finalViews, _ := s.rdb.Get(ctx, key).Int()
 fmt.Printf("Final views: %d\n", finalViews)
}

// ===== 分布式ID生成器 =====

// GenerateID 生成分布式ID
func (s *StringOps) GenerateID(ctx context.Context, key string) (int64, error) {
 // 使用INCR生成唯一ID
 return s.rdb.Incr(ctx, key).Result()
}

// ===== 实战案例：点赞系统 =====

type LikeService struct {
 rdb *redis.Client
}

func NewLikeService(rdb *redis.Client) *LikeService {
 return &LikeService{rdb: rdb}
}

// Like 点赞
func (l *LikeService) Like(ctx context.Context, postID string) (int64, error) {
 key := fmt.Sprintf("post:%s:likes", postID)
 return l.rdb.Incr(ctx, key).Result()
}

// Unlike 取消点赞
func (l *LikeService) Unlike(ctx context.Context, postID string) (int64, error) {
 key := fmt.Sprintf("post:%s:likes", postID)
 return l.rdb.Decr(ctx, key).Result()
}

// GetLikes 获取点赞数
func (l *LikeService) GetLikes(ctx context.Context, postID string) (int64, error) {
 key := fmt.Sprintf("post:%s:likes", postID)
 return l.rdb.Get(ctx, key).Int64()
}

// ===== 实战案例：限流器 =====

type RateLimiter struct {
 rdb *redis.Client
}

func NewRateLimiter(rdb *redis.Client) *RateLimiter {
 return &RateLimiter{rdb: rdb}
}

// Allow 检查是否允许访问（简单计数器限流）
func (r *RateLimiter) Allow(ctx context.Context, userID string, maxRequests int, window time.Duration) (bool, error) {
 key := fmt.Sprintf("rate_limit:%s", userID)

 // 获取当前计数
 count, err := r.rdb.Get(ctx, key).Int()
 if err != nil && err != redis.Nil {
  return false, err
 }

 // 如果是第一次访问，设置初始值和过期时间
 if err == redis.Nil {
  pipe := r.rdb.Pipeline()
  pipe.Set(ctx, key, 1, window)
  _, err = pipe.Exec(ctx)
  return true, err
 }

 // 检查是否超过限制
 if count >= maxRequests {
  return false, nil
 }

 // 增加计数
 r.rdb.Incr(ctx, key)
 return true, nil
}

func main() {
 rdb := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
 })
 defer rdb.Close()

 ctx := context.Background()

 // String操作
 stringOps := NewStringOps(rdb)
 stringOps.SetGet(ctx)
 stringOps.Counter(ctx)

 // 点赞系统
 likeService := NewLikeService(rdb)
 likes, _ := likeService.Like(ctx, "123")
 fmt.Printf("Post 123 likes: %d\n", likes)

 // 限流器
 limiter := NewRateLimiter(rdb)
 allowed, _ := limiter.Allow(ctx, "user1", 10, time.Minute)
 fmt.Printf("Request allowed: %v\n", allowed)
}
```

---

### 实战案例3：Hash（哈希表）

```go
package main

import (
 "context"
 "fmt"
 "log"

 "github.com/redis/go-redis/v9"
)

type HashOps struct {
 rdb *redis.Client
}

func NewHashOps(rdb *redis.Client) *HashOps {
 return &HashOps{rdb: rdb}
}

// ===== 基础操作 =====

// BasicOps Hash基础操作
func (h *HashOps) BasicOps(ctx context.Context) {
 key := "user:1000"

 // HSET (设置单个字段)
 h.rdb.HSet(ctx, key, "name", "Alice")
 h.rdb.HSet(ctx, key, "age", "25")
 h.rdb.HSet(ctx, key, "email", "alice@example.com")

 // HGET (获取单个字段)
 name, _ := h.rdb.HGet(ctx, key, "name").Result()
 fmt.Printf("Name: %s\n", name)

 // HMSET (批量设置)
 h.rdb.HMSet(ctx, key, map[string]interface{}{
  "city":  "Beijing",
  "score": "95",
 })

 // HMGET (批量获取)
 vals, _ := h.rdb.HMGet(ctx, key, "name", "city", "score").Result()
 fmt.Printf("Values: %v\n", vals)

 // HGETALL (获取所有字段)
 all, _ := h.rdb.HGetAll(ctx, key).Result()
 fmt.Printf("All fields: %v\n", all)

 // HEXISTS (检查字段是否存在)
 exists, _ := h.rdb.HExists(ctx, key, "email").Result()
 fmt.Printf("Email exists: %v\n", exists)

 // HDEL (删除字段)
 h.rdb.HDel(ctx, key, "email")

 // HLEN (获取字段数量)
 length, _ := h.rdb.HLen(ctx, key).Result()
 fmt.Printf("Field count: %d\n", length)

 // HKEYS (获取所有字段名)
 keys, _ := h.rdb.HKeys(ctx, key).Result()
 fmt.Printf("Keys: %v\n", keys)

 // HVALS (获取所有值)
 values, _ := h.rdb.HVals(ctx, key).Result()
 fmt.Printf("Values: %v\n", values)
}

// ===== 计数器操作 =====

// IncrementField Hash字段自增
func (h *HashOps) IncrementField(ctx context.Context) {
 key := "user:1000:stats"

 // HINCRBY (整数增加)
 h.rdb.HIncrBy(ctx, key, "visits", 1)
 h.rdb.HIncrBy(ctx, key, "posts", 1)

 // HINCRBYFLOAT (浮点数增加)
 h.rdb.HIncrByFloat(ctx, key, "rating", 0.5)

 stats, _ := h.rdb.HGetAll(ctx, key).Result()
 fmt.Printf("User stats: %v\n", stats)
}

// ===== 实战案例：用户信息存储 =====

type User struct {
 ID       string
 Username string
 Email    string
 Age      int
 City     string
}

type UserRepository struct {
 rdb *redis.Client
}

func NewUserRepository(rdb *redis.Client) *UserRepository {
 return &UserRepository{rdb: rdb}
}

// Save 保存用户信息
func (r *UserRepository) Save(ctx context.Context, user *User) error {
 key := fmt.Sprintf("user:%s", user.ID)
 return r.rdb.HMSet(ctx, key, map[string]interface{}{
  "username": user.Username,
  "email":    user.Email,
  "age":      user.Age,
  "city":     user.City,
 }).Err()
}

// Get 获取用户信息
func (r *UserRepository) Get(ctx context.Context, id string) (*User, error) {
 key := fmt.Sprintf("user:%s", id)
 data, err := r.rdb.HGetAll(ctx, key).Result()
 if err != nil {
  return nil, err
 }

 if len(data) == 0 {
  return nil, fmt.Errorf("user not found")
 }

 user := &User{
  ID:       id,
  Username: data["username"],
  Email:    data["email"],
  City:     data["city"],
 }

 // 转换age
 fmt.Sscanf(data["age"], "%d", &user.Age)

 return user, nil
}

// Update 更新字段
func (r *UserRepository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
 key := fmt.Sprintf("user:%s", id)
 return r.rdb.HMSet(ctx, key, fields).Err()
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id string) error {
 key := fmt.Sprintf("user:%s", id)
 return r.rdb.Del(ctx, key).Err()
}

// ===== 实战案例：购物车 =====

type ShoppingCart struct {
 rdb *redis.Client
}

func NewShoppingCart(rdb *redis.Client) *ShoppingCart {
 return &ShoppingCart{rdb: rdb}
}

// AddItem 添加商品到购物车
func (s *ShoppingCart) AddItem(ctx context.Context, userID, productID string, quantity int) error {
 key := fmt.Sprintf("cart:%s", userID)
 return s.rdb.HIncrBy(ctx, key, productID, int64(quantity)).Err()
}

// RemoveItem 从购物车移除商品
func (s *ShoppingCart) RemoveItem(ctx context.Context, userID, productID string) error {
 key := fmt.Sprintf("cart:%s", userID)
 return s.rdb.HDel(ctx, key, productID).Err()
}

// UpdateQuantity 更新商品数量
func (s *ShoppingCart) UpdateQuantity(ctx context.Context, userID, productID string, quantity int) error {
 key := fmt.Sprintf("cart:%s", userID)
 return s.rdb.HSet(ctx, key, productID, quantity).Err()
}

// GetCart 获取购物车
func (s *ShoppingCart) GetCart(ctx context.Context, userID string) (map[string]string, error) {
 key := fmt.Sprintf("cart:%s", userID)
 return s.rdb.HGetAll(ctx, key).Result()
}

// GetItemCount 获取购物车商品种类数
func (s *ShoppingCart) GetItemCount(ctx context.Context, userID string) (int64, error) {
 key := fmt.Sprintf("cart:%s", userID)
 return s.rdb.HLen(ctx, key).Result()
}

// Clear 清空购物车
func (s *ShoppingCart) Clear(ctx context.Context, userID string) error {
 key := fmt.Sprintf("cart:%s", userID)
 return s.rdb.Del(ctx, key).Err()
}

func main() {
 rdb := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
 })
 defer rdb.Close()

 ctx := context.Background()

 // Hash基础操作
 hashOps := NewHashOps(rdb)
 hashOps.BasicOps(ctx)
 hashOps.IncrementField(ctx)

 // 用户信息存储
 userRepo := NewUserRepository(rdb)
 user := &User{
  ID:       "1001",
  Username: "alice",
  Email:    "alice@example.com",
  Age:      25,
  City:     "Beijing",
 }
 userRepo.Save(ctx, user)

 // 获取用户
 loadedUser, _ := userRepo.Get(ctx, "1001")
 fmt.Printf("Loaded user: %+v\n", loadedUser)

 // 购物车
 cart := NewShoppingCart(rdb)
 cart.AddItem(ctx, "user1", "product123", 2)
 cart.AddItem(ctx, "user1", "product456", 1)
 items, _ := cart.GetCart(ctx, "user1")
 fmt.Printf("Cart items: %v\n", items)
}
```

---

## 第五部分：缓存策略

### 实战案例4：缓存穿透/击穿/雪崩

```go
package main

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "log"
 "math/rand"
 "sync"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== 缓存穿透解决方案：布隆过滤器 + 空值缓存 =====

type CacheService struct {
 rdb *redis.Client
 mu  sync.Mutex
}

func NewCacheService(rdb *redis.Client) *CacheService {
 return &CacheService{rdb: rdb}
}

// GetWithCachePenetration 防止缓存穿透
func (c *CacheService) GetWithCachePenetration(ctx context.Context, key string, loader func() (interface{}, error)) (string, error) {
 // 1. 查询缓存
 val, err := c.rdb.Get(ctx, key).Result()
 if err == nil {
  // 命中缓存
  if val == "NULL" {
   // 空值缓存
   return "", fmt.Errorf("not found")
  }
  return val, nil
 }

 // 2. 缓存未命中，查询数据库
 data, err := loader()
 if err != nil {
  return "", err
 }

 // 3. 数据不存在，缓存空值（防止穿透）
 if data == nil {
  c.rdb.Set(ctx, key, "NULL", 5*time.Minute)
  return "", fmt.Errorf("not found")
 }

 // 4. 数据存在，缓存数据
 jsonData, _ := json.Marshal(data)
 c.rdb.Set(ctx, key, jsonData, 30*time.Minute)

 return string(jsonData), nil
}

// ===== 缓存击穿解决方案：互斥锁 =====

// GetWithCacheBreakdown 防止缓存击穿（热点key）
func (c *CacheService) GetWithCacheBreakdown(ctx context.Context, key string, loader func() (interface{}, error), ttl time.Duration) (string, error) {
 // 1. 查询缓存
 val, err := c.rdb.Get(ctx, key).Result()
 if err == nil {
  return val, nil
 }

 // 2. 使用分布式锁
 lockKey := fmt.Sprintf("lock:%s", key)
 locked, err := c.rdb.SetNX(ctx, lockKey, "1", 10*time.Second).Result()
 if err != nil {
  return "", err
 }

 if !locked {
  // 获取锁失败，等待后重试
  time.Sleep(100 * time.Millisecond)
  return c.GetWithCacheBreakdown(ctx, key, loader, ttl)
 }

 // 获取锁成功，释放锁
 defer c.rdb.Del(ctx, lockKey)

 // 3. 双重检查（DCL）
 val, err = c.rdb.Get(ctx, key).Result()
 if err == nil {
  return val, nil
 }

 // 4. 加载数据
 data, err := loader()
 if err != nil {
  return "", err
 }

 // 5. 缓存数据
 jsonData, _ := json.Marshal(data)
 c.rdb.Set(ctx, key, jsonData, ttl)

 return string(jsonData), nil
}

// ===== 缓存雪崩解决方案：随机过期时间 =====

// SetWithRandomTTL 设置随机过期时间（防止雪崩）
func (c *CacheService) SetWithRandomTTL(ctx context.Context, key string, value interface{}, baseTTL time.Duration) error {
 // 在基础TTL上增加随机时间（0-300秒）
 randomTTL := baseTTL + time.Duration(rand.Intn(300))*time.Second
 
 jsonData, err := json.Marshal(value)
 if err != nil {
  return err
 }

 return c.rdb.Set(ctx, key, jsonData, randomTTL).Err()
}

// ===== 实战案例：多级缓存 =====

type MultiLevelCache struct {
 rdb        *redis.Client
 localCache *sync.Map // 本地缓存
}

func NewMultiLevelCache(rdb *redis.Client) *MultiLevelCache {
 return &MultiLevelCache{
  rdb:        rdb,
  localCache: &sync.Map{},
 }
}

type CacheItem struct {
 Value      interface{}
 ExpireTime time.Time
}

// Get 多级缓存获取
func (m *MultiLevelCache) Get(ctx context.Context, key string, loader func() (interface{}, error)) (interface{}, error) {
 // 1. 查询本地缓存
 if item, ok := m.localCache.Load(key); ok {
  cacheItem := item.(*CacheItem)
  if time.Now().Before(cacheItem.ExpireTime) {
   log.Printf("[Local Cache] Hit: %s", key)
   return cacheItem.Value, nil
  }
  // 本地缓存过期，删除
  m.localCache.Delete(key)
 }

 // 2. 查询Redis缓存
 val, err := m.rdb.Get(ctx, key).Result()
 if err == nil {
  log.Printf("[Redis Cache] Hit: %s", key)
  
  // 更新本地缓存
  m.localCache.Store(key, &CacheItem{
   Value:      val,
   ExpireTime: time.Now().Add(1 * time.Minute),
  })
  
  return val, nil
 }

 // 3. 缓存未命中，加载数据
 log.Printf("[Cache] Miss: %s", key)
 data, err := loader()
 if err != nil {
  return nil, err
 }

 // 4. 更新Redis缓存
 jsonData, _ := json.Marshal(data)
 m.rdb.Set(ctx, key, jsonData, 30*time.Minute)

 // 5. 更新本地缓存
 m.localCache.Store(key, &CacheItem{
  Value:      data,
  ExpireTime: time.Now().Add(1 * time.Minute),
 })

 return data, nil
}
```

---

### 实战案例5：List/Set/ZSet操作

```go
package main

import (
 "context"
 "fmt"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== List（列表）操作 =====

type ListOps struct {
 rdb *redis.Client
}

func NewListOps(rdb *redis.Client) *ListOps {
 return &ListOps{rdb: rdb}
}

// 消息队列示例
func (l *ListOps) MessageQueue(ctx context.Context) {
 queue := "task_queue"

 // 生产者：LPUSH（左侧插入）
 l.rdb.LPush(ctx, queue, "task1", "task2", "task3")

 // 消费者：BRPOP（右侧阻塞弹出）
 result, _ := l.rdb.BRPop(ctx, 5*time.Second, queue).Result()
 fmt.Printf("Consumed: %v\n", result)

 // 查看队列长度
 length, _ := l.rdb.LLen(ctx, queue).Result()
 fmt.Printf("Queue length: %d\n", length)
}

// ===== Set（集合）操作 =====

type SetOps struct {
 rdb *redis.Client
}

func NewSetOps(rdb *redis.Client) *SetOps {
 return &SetOps{rdb: rdb}
}

// 标签系统示例
func (s *SetOps) TagSystem(ctx context.Context) {
 // 文章标签
 s.rdb.SAdd(ctx, "article:1:tags", "go", "redis", "performance")
 s.rdb.SAdd(ctx, "article:2:tags", "go", "kubernetes", "cloud")

 // 获取文章1的所有标签
 tags, _ := s.rdb.SMembers(ctx, "article:1:tags").Result()
 fmt.Printf("Article 1 tags: %v\n", tags)

 // 查找共同标签（交集）
 common, _ := s.rdb.SInter(ctx, "article:1:tags", "article:2:tags").Result()
 fmt.Printf("Common tags: %v\n", common)

 // 检查标签是否存在
 exists, _ := s.rdb.SIsMember(ctx, "article:1:tags", "redis").Result()
 fmt.Printf("Has redis tag: %v\n", exists)
}

// ===== ZSet（有序集合）操作 =====

type ZSetOps struct {
 rdb *redis.Client
}

func NewZSetOps(rdb *redis.Client) *ZSetOps {
 return &ZSetOps{rdb: rdb}
}

// 排行榜示例
func (z *ZSetOps) Leaderboard(ctx context.Context) {
 leaderboard := "game:leaderboard"

 // 添加分数
 z.rdb.ZAdd(ctx, leaderboard, redis.Z{Score: 100, Member: "player1"})
 z.rdb.ZAdd(ctx, leaderboard, redis.Z{Score: 200, Member: "player2"})
 z.rdb.ZAdd(ctx, leaderboard, redis.Z{Score: 150, Member: "player3"})

 // 增加分数
 z.rdb.ZIncrBy(ctx, leaderboard, 50, "player1")

 // 获取排名（从高到低）
 topPlayers, _ := z.rdb.ZRevRangeWithScores(ctx, leaderboard, 0, 2).Result()
 fmt.Println("Top 3 players:")
 for i, player := range topPlayers {
  fmt.Printf("%d. %s: %.0f\n", i+1, player.Member, player.Score)
 }

 // 获取玩家排名
 rank, _ := z.rdb.ZRevRank(ctx, leaderboard, "player1").Result()
 fmt.Printf("player1 rank: %d\n", rank+1)

 // 获取分数
 score, _ := z.rdb.ZScore(ctx, leaderboard, "player1").Result()
 fmt.Printf("player1 score: %.0f\n", score)
}
```

---

## 第四部分：高级数据结构

### 实战案例6：Bitmap/HyperLogLog/Geo

```go
package main

import (
 "context"
 "fmt"

 "github.com/redis/go-redis/v9"
)

// ===== Bitmap（位图）操作 =====

type BitmapOps struct {
 rdb *redis.Client
}

func NewBitmapOps(rdb *redis.Client) *BitmapOps {
 return &BitmapOps{rdb: rdb}
}

// 用户签到系统
func (b *BitmapOps) UserCheckIn(ctx context.Context) {
 // 用户ID 123 在第1、3、5天签到
 b.rdb.SetBit(ctx, "checkin:123:2024-01", 1, 1)
 b.rdb.SetBit(ctx, "checkin:123:2024-01", 3, 1)
 b.rdb.SetBit(ctx, "checkin:123:2024-01", 5, 1)

 // 检查第3天是否签到
 checked, _ := b.rdb.GetBit(ctx, "checkin:123:2024-01", 3).Result()
 fmt.Printf("Day 3 checked in: %d\n", checked)

 // 统计签到天数
 count, _ := b.rdb.BitCount(ctx, "checkin:123:2024-01", nil).Result()
 fmt.Printf("Total check-in days: %d\n", count)
}

// ===== HyperLogLog（基数统计）=====

type HyperLogLogOps struct {
 rdb *redis.Client
}

func NewHyperLogLogOps(rdb *redis.Client) *HyperLogLogOps {
 return &HyperLogLogOps{rdb: rdb}
}

// UV统计（独立访客）
func (h *HyperLogLogOps) UniqueVisitors(ctx context.Context) {
 // 添加访客
 h.rdb.PFAdd(ctx, "page:home:uv:2024-01-01", "user1", "user2", "user3")
 h.rdb.PFAdd(ctx, "page:home:uv:2024-01-01", "user2", "user4") // user2重复

 // 获取UV数量
 uv, _ := h.rdb.PFCount(ctx, "page:home:uv:2024-01-01").Result()
 fmt.Printf("Unique visitors: %d\n", uv) // 输出4（去重）
}

// ===== Geo（地理位置）=====

type GeoOps struct {
 rdb *redis.Client
}

func NewGeoOps(rdb *redis.Client) *GeoOps {
 return &GeoOps{rdb: rdb}
}

// 附近的人
func (g *GeoOps) NearbyLocations(ctx context.Context) {
 locations := "drivers:locations"

 // 添加位置（经度、纬度、成员）
 g.rdb.GeoAdd(ctx, locations,
  &redis.GeoLocation{Longitude: 116.397128, Latitude: 39.916527, Name: "driver1"},
  &redis.GeoLocation{Longitude: 116.405285, Latitude: 39.904989, Name: "driver2"},
 )

 // 查找半径5km内的司机
 results, _ := g.rdb.GeoRadius(ctx, locations, 116.4, 39.9, &redis.GeoRadiusQuery{
  Radius:      5,
  Unit:        "km",
  WithCoord:   true,
  WithDist:    true,
  Count:       10,
  Sort:        "ASC",
 }).Result()

 fmt.Println("Nearby drivers:")
 for _, result := range results {
  fmt.Printf("- %s: %.2f km\n", result.Name, result.Dist)
 }
}
```

---

## 第六部分：分布式锁

### 实战案例7：Redis分布式锁

```go
package main

import (
 "context"
 "errors"
 "fmt"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== 简单分布式锁 =====

type DistributedLock struct {
 rdb   *redis.Client
 key   string
 value string
 ttl   time.Duration
}

func NewDistributedLock(rdb *redis.Client, key string, ttl time.Duration) *DistributedLock {
 return &DistributedLock{
  rdb:   rdb,
  key:   fmt.Sprintf("lock:%s", key),
  value: fmt.Sprintf("%d", time.Now().UnixNano()),
  ttl:   ttl,
 }
}

// Lock 获取锁
func (d *DistributedLock) Lock(ctx context.Context) (bool, error) {
 return d.rdb.SetNX(ctx, d.key, d.value, d.ttl).Result()
}

// Unlock 释放锁（使用Lua脚本确保原子性）
func (d *DistributedLock) Unlock(ctx context.Context) error {
 script := `
  if redis.call("get", KEYS[1]) == ARGV[1] then
   return redis.call("del", KEYS[1])
  else
   return 0
  end
 `
 result, err := d.rdb.Eval(ctx, script, []string{d.key}, d.value).Result()
 if err != nil {
  return err
 }
 if result == int64(0) {
  return errors.New("lock not held")
 }
 return nil
}

// ===== 使用示例 =====

func UseDistributedLock(ctx context.Context, rdb *redis.Client) {
 lock := NewDistributedLock(rdb, "order:123", 10*time.Second)

 // 尝试获取锁
 acquired, err := lock.Lock(ctx)
 if err != nil {
  fmt.Printf("Lock error: %v\n", err)
  return
 }

 if !acquired {
  fmt.Println("Failed to acquire lock")
  return
 }

 // 执行业务逻辑
 fmt.Println("Lock acquired, processing order...")
 time.Sleep(2 * time.Second)

 // 释放锁
 if err := lock.Unlock(ctx); err != nil {
  fmt.Printf("Unlock error: %v\n", err)
 } else {
  fmt.Println("Lock released")
 }
}

// ===== 可重入锁 =====

type ReentrantLock struct {
 rdb      *redis.Client
 key      string
 threadID string
 ttl      time.Duration
}

func NewReentrantLock(rdb *redis.Client, key, threadID string, ttl time.Duration) *ReentrantLock {
 return &ReentrantLock{
  rdb:      rdb,
  key:      fmt.Sprintf("reentrant_lock:%s", key),
  threadID: threadID,
  ttl:      ttl,
 }
}

// Lock 可重入锁
func (r *ReentrantLock) Lock(ctx context.Context) (bool, error) {
 script := `
  if redis.call("exists", KEYS[1]) == 0 then
   redis.call("hset", KEYS[1], ARGV[1], 1)
   redis.call("expire", KEYS[1], ARGV[2])
   return 1
  elseif redis.call("hexists", KEYS[1], ARGV[1]) == 1 then
   redis.call("hincrby", KEYS[1], ARGV[1], 1)
   redis.call("expire", KEYS[1], ARGV[2])
   return 1
  else
   return 0
  end
 `
 result, err := r.rdb.Eval(ctx, script, []string{r.key}, r.threadID, int(r.ttl.Seconds())).Result()
 if err != nil {
  return false, err
 }
 return result == int64(1), nil
}

// Unlock 释放可重入锁
func (r *ReentrantLock) Unlock(ctx context.Context) error {
 script := `
  if redis.call("hexists", KEYS[1], ARGV[1]) == 0 then
   return nil
  else
   local count = redis.call("hincrby", KEYS[1], ARGV[1], -1)
   if count > 0 then
    return 0
   else
    redis.call("del", KEYS[1])
    return 1
   end
  end
 `
 _, err := r.rdb.Eval(ctx, script, []string{r.key}, r.threadID).Result()
 return err
}
```

---

## 第七部分：消息队列

### 实战案例8：List实现消息队列

```go
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== 简单消息队列 =====

type Task struct {
 ID        string    `json:"id"`
 Type      string    `json:"type"`
 Payload   string    `json:"payload"`
 CreatedAt time.Time `json:"created_at"`
}

type TaskQueue struct {
 rdb       *redis.Client
 queueName string
}

func NewTaskQueue(rdb *redis.Client, queueName string) *TaskQueue {
 return &TaskQueue{
  rdb:       rdb,
  queueName: queueName,
 }
}

// Push 推送任务
func (q *TaskQueue) Push(ctx context.Context, task *Task) error {
 data, err := json.Marshal(task)
 if err != nil {
  return err
 }
 return q.rdb.LPush(ctx, q.queueName, data).Err()
}

// Pop 消费任务（阻塞式）
func (q *TaskQueue) Pop(ctx context.Context, timeout time.Duration) (*Task, error) {
 result, err := q.rdb.BRPop(ctx, timeout, q.queueName).Result()
 if err != nil {
  return nil, err
 }

 if len(result) < 2 {
  return nil, fmt.Errorf("invalid result")
 }

 var task Task
 if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
  return nil, err
 }

 return &task, nil
}

// GetLength 获取队列长度
func (q *TaskQueue) GetLength(ctx context.Context) (int64, error) {
 return q.rdb.LLen(ctx, q.queueName).Result()
}

// ===== 延迟队列 =====

type DelayQueue struct {
 rdb       *redis.Client
 queueName string
}

func NewDelayQueue(rdb *redis.Client, queueName string) *DelayQueue {
 return &DelayQueue{
  rdb:       rdb,
  queueName: queueName,
 }
}

// Push 推送延迟任务
func (d *DelayQueue) Push(ctx context.Context, task *Task, delay time.Duration) error {
 data, err := json.Marshal(task)
 if err != nil {
  return err
 }

 score := float64(time.Now().Add(delay).Unix())
 return d.rdb.ZAdd(ctx, d.queueName, redis.Z{
  Score:  score,
  Member: data,
 }).Err()
}

// PopReady 获取到期的任务
func (d *DelayQueue) PopReady(ctx context.Context) ([]*Task, error) {
 now := float64(time.Now().Unix())

 // 获取到期的任务
 results, err := d.rdb.ZRangeByScore(ctx, d.queueName, &redis.ZRangeBy{
  Min:   "-inf",
  Max:   fmt.Sprintf("%f", now),
  Count: 10,
 }).Result()

 if err != nil {
  return nil, err
 }

 var tasks []*Task
 for _, result := range results {
  var task Task
  if err := json.Unmarshal([]byte(result), &task); err != nil {
   continue
  }
  tasks = append(tasks, &task)

  // 删除已处理的任务
  d.rdb.ZRem(ctx, d.queueName, result)
 }

 return tasks, nil
}
```

---

## 第十部分：事务与管道

### 实战案例9：Redis事务

```go
package main

import (
 "context"
 "errors"
 "fmt"

 "github.com/redis/go-redis/v9"
)

// ===== Redis事务（MULTI/EXEC）=====

func TransferMoney(ctx context.Context, rdb *redis.Client, from, to string, amount int64) error {
 // 使用WATCH实现乐观锁
 txf := func(tx *redis.Tx) error {
  // 获取账户余额
  fromBalance, err := tx.Get(ctx, fmt.Sprintf("account:%s", from)).Int64()
  if err != nil && err != redis.Nil {
   return err
  }

  // 检查余额
  if fromBalance < amount {
   return errors.New("insufficient balance")
  }

  // 执行事务
  _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
   pipe.DecrBy(ctx, fmt.Sprintf("account:%s", from), amount)
   pipe.IncrBy(ctx, fmt.Sprintf("account:%s", to), amount)
   return nil
  })

  return err
 }

 // 重试机制
 for i := 0; i < 3; i++ {
  err := rdb.Watch(ctx, txf, fmt.Sprintf("account:%s", from))
  if err == nil {
   return nil
  }
  if err == redis.TxFailedErr {
   // 事务冲突，重试
   continue
  }
  return err
 }

 return errors.New("transaction failed after retries")
}

// ===== Pipeline（管道）=====

func UsePipeline(ctx context.Context, rdb *redis.Client) {
 pipe := rdb.Pipeline()

 // 批量操作
 incr := pipe.Incr(ctx, "counter")
 pipe.Set(ctx, "key1", "value1", 0)
 pipe.Set(ctx, "key2", "value2", 0)
 pipe.Get(ctx, "key1")

 // 执行管道
 cmds, err := pipe.Exec(ctx)
 if err != nil {
  fmt.Printf("Pipeline error: %v\n", err)
  return
 }

 // 获取结果
 fmt.Printf("Counter: %d\n", incr.Val())
 fmt.Printf("Commands executed: %d\n", len(cmds))
}
```

---

## 第十一部分：Lua脚本

### 实战案例10：Lua脚本原子操作

```go
package main

import (
 "context"
 "fmt"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== 限流器（令牌桶算法）=====

func RateLimitWithLua(ctx context.Context, rdb *redis.Client, key string, rate, capacity int64) (bool, error) {
 script := `
  local key = KEYS[1]
  local rate = tonumber(ARGV[1])
  local capacity = tonumber(ARGV[2])
  local now = tonumber(ARGV[3])
  local requested = 1

  local bucket = redis.call('hmget', key, 'tokens', 'last_time')
  local tokens = tonumber(bucket[1])
  local last_time = tonumber(bucket[2])

  if tokens == nil then
   tokens = capacity
   last_time = now
  end

  local delta = math.max(0, now - last_time)
  tokens = math.min(capacity, tokens + delta * rate)

  if tokens >= requested then
   tokens = tokens - requested
   redis.call('hmset', key, 'tokens', tokens, 'last_time', now)
   redis.call('expire', key, capacity / rate)
   return 1
  else
   return 0
  end
 `

 now := fmt.Sprintf("%d", time.Now().Unix())
 result, err := rdb.Eval(ctx, script,
  []string{key},
  rate, capacity, now,
 ).Result()

 if err != nil {
  return false, err
 }

 return result == int64(1), nil
}

// ===== 原子递增（带上限）=====

func IncrWithLimit(ctx context.Context, rdb *redis.Client, key string, limit int64) (int64, error) {
 script := `
  local key = KEYS[1]
  local limit = tonumber(ARGV[1])
  local current = tonumber(redis.call('get', key) or "0")

  if current < limit then
   return redis.call('incr', key)
  else
   return -1
  end
 `

 result, err := rdb.Eval(ctx, script, []string{key}, limit).Int64()
 return result, err
}
```

---

## 第八部分：Stream流处理

### 实战案例11：Redis Stream消费者组

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== Stream生产者 =====

type StreamProducer struct {
 rdb        *redis.Client
 streamName string
}

func NewStreamProducer(rdb *redis.Client, streamName string) *StreamProducer {
 return &StreamProducer{
  rdb:        rdb,
  streamName: streamName,
 }
}

// Publish 发布消息到Stream
func (p *StreamProducer) Publish(ctx context.Context, data map[string]interface{}) (string, error) {
 id, err := p.rdb.XAdd(ctx, &redis.XAddArgs{
  Stream: p.streamName,
  Values: data,
 }).Result()
 return id, err
}

// ===== Stream消费者组 =====

type StreamConsumer struct {
 rdb        *redis.Client
 streamName string
 groupName  string
 consumer   string
}

func NewStreamConsumer(rdb *redis.Client, streamName, groupName, consumerName string) *StreamConsumer {
 return &StreamConsumer{
  rdb:        rdb,
  streamName: streamName,
  groupName:  groupName,
  consumer:   consumerName,
 }
}

// CreateGroup 创建消费者组
func (c *StreamConsumer) CreateGroup(ctx context.Context) error {
 return c.rdb.XGroupCreate(ctx, c.streamName, c.groupName, "0").Err()
}

// Consume 消费消息
func (c *StreamConsumer) Consume(ctx context.Context) {
 for {
  // 读取消息
  streams, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
   Group:    c.groupName,
   Consumer: c.consumer,
   Streams:  []string{c.streamName, ">"},
   Count:    10,
   Block:    time.Second * 5,
  }).Result()

  if err != nil {
   log.Printf("Read error: %v", err)
   continue
  }

  for _, stream := range streams {
   for _, message := range stream.Messages {
    log.Printf("Consumer %s received: %s - %v", c.consumer, message.ID, message.Values)

    // 处理消息
    // ...

    // 确认消息
    c.rdb.XAck(ctx, c.streamName, c.groupName, message.ID)
   }
  }
 }
}
```

---

## 第九部分：发布订阅

### 实战案例12：Pub/Sub模式

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== 发布者 =====

type Publisher struct {
 rdb *redis.Client
}

func NewPublisher(rdb *redis.Client) *Publisher {
 return &Publisher{rdb: rdb}
}

// Publish 发布消息
func (p *Publisher) Publish(ctx context.Context, channel, message string) error {
 return p.rdb.Publish(ctx, channel, message).Err()
}

// ===== 订阅者 =====

type Subscriber struct {
 rdb *redis.Client
}

func NewSubscriber(rdb *redis.Client) *Subscriber {
 return &Subscriber{rdb: rdb}
}

// Subscribe 订阅频道
func (s *Subscriber) Subscribe(ctx context.Context, channels ...string) {
 pubsub := s.rdb.Subscribe(ctx, channels...)
 defer pubsub.Close()

 // 接收消息
 ch := pubsub.Channel()
 for msg := range ch {
  log.Printf("Received from %s: %s", msg.Channel, msg.Payload)
 }
}

// PSubscribe 模式订阅
func (s *Subscriber) PSubscribe(ctx context.Context, patterns ...string) {
 pubsub := s.rdb.PSubscribe(ctx, patterns...)
 defer pubsub.Close()

 ch := pubsub.Channel()
 for msg := range ch {
  log.Printf("Pattern %s matched %s: %s", msg.Pattern, msg.Channel, msg.Payload)
 }
}

// ===== 使用示例 =====

func PubSubExample() {
 rdb := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
 })

 ctx := context.Background()

 // 启动订阅者
 subscriber := NewSubscriber(rdb)
 go subscriber.Subscribe(ctx, "news", "sports")

 // 模式订阅
 go subscriber.PSubscribe(ctx, "event:*")

 time.Sleep(100 * time.Millisecond)

 // 发布消息
 publisher := NewPublisher(rdb)
 publisher.Publish(ctx, "news", "Breaking news!")
 publisher.Publish(ctx, "event:user", "User logged in")
}
```

---

## 第十二部分：持久化与高可用

### 持久化配置

```go
// ===== RDB持久化 =====

// 配置文件redis.conf:
// save 900 1        # 900秒内至少1个key变化
// save 300 10       # 300秒内至少10个key变化
// save 60 10000     # 60秒内至少10000个key变化

// 手动触发RDB快照
func SaveSnapshot(ctx context.Context, rdb *redis.Client) error {
 // SAVE（阻塞）
 return rdb.Save(ctx).Err()
}

func BackgroundSave(ctx context.Context, rdb *redis.Client) error {
 // BGSAVE（后台异步）
 return rdb.BgSave(ctx).Err()
}

// ===== AOF持久化 =====

// 配置文件redis.conf:
// appendonly yes
// appendfsync everysec  # always | everysec | no

// 重写AOF
func RewriteAOF(ctx context.Context, rdb *redis.Client) error {
 return rdb.BgRewriteAOF(ctx).Err()
}
```

---

## 第十三部分：性能优化

### 实战案例13：性能优化技巧

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== 1. 使用Pipeline批量操作 =====

func BatchOperations(ctx context.Context, rdb *redis.Client) {
 pipe := rdb.Pipeline()

 for i := 0; i < 1000; i++ {
  pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, 0)
 }

 // 一次性执行
 _, err := pipe.Exec(ctx)
 if err != nil {
  log.Printf("Pipeline error: %v", err)
 }
}

// ===== 2. 连接池配置 =====

func OptimizedClient() *redis.Client {
 return redis.NewClient(&redis.Options{
  Addr:         "localhost:6379",
  PoolSize:     100,              // 连接池大小
  MinIdleConns: 10,               // 最小空闲连接
  MaxRetries:   3,                // 重试次数
  PoolTimeout:  4 * time.Second,  // 池超时
  IdleTimeout:  5 * time.Minute,  // 空闲超时
 })
}

// ===== 3. 避免大key =====

func AvoidBigKeys(ctx context.Context, rdb *redis.Client) {
 // 不好：一个大Hash
 // for i := 0; i < 1000000; i++ {
 //     rdb.HSet(ctx, "big_hash", fmt.Sprintf("field%d", i), i)
 // }

 // 好：分片存储
 shardCount := 100
 for i := 0; i < 1000000; i++ {
  shardID := i % shardCount
  rdb.HSet(ctx, fmt.Sprintf("hash_shard:%d", shardID), fmt.Sprintf("field%d", i), i)
 }
}

// ===== 4. 使用Scan代替Keys =====

func ScanKeys(ctx context.Context, rdb *redis.Client, pattern string) []string {
 var keys []string
 var cursor uint64

 for {
  var err error
  var scanKeys []string

  // 使用SCAN代替KEYS
  scanKeys, cursor, err = rdb.Scan(ctx, cursor, pattern, 100).Result()
  if err != nil {
   break
  }

  keys = append(keys, scanKeys...)

  if cursor == 0 {
   break
  }
 }

 return keys
}
```

---

## 第十四部分：监控与运维

### 实战案例14：Redis监控

```go
package main

import (
 "context"
 "fmt"
 "log"

 "github.com/redis/go-redis/v9"
)

// ===== Redis INFO命令 =====

func MonitorRedis(ctx context.Context, rdb *redis.Client) {
 // 获取服务器信息
 info, err := rdb.Info(ctx).Result()
 if err != nil {
  log.Printf("Info error: %v", err)
  return
 }
 fmt.Println("Server Info:", info)

 // 获取统计信息
 stats, err := rdb.Info(ctx, "stats").Result()
 if err != nil {
  log.Printf("Stats error: %v", err)
  return
 }
 fmt.Println("Stats:", stats)

 // 获取内存信息
 memory, err := rdb.Info(ctx, "memory").Result()
 if err != nil {
  log.Printf("Memory error: %v", err)
  return
 }
 fmt.Println("Memory:", memory)
}

// ===== 慢查询监控 =====

func MonitorSlowLog(ctx context.Context, rdb *redis.Client) {
 // 获取慢查询日志
 slowLogs, err := rdb.SlowLogGet(ctx, 10).Result()
 if err != nil {
  log.Printf("SlowLog error: %v", err)
  return
 }

 for _, log := range slowLogs {
  fmt.Printf("ID: %d, Time: %v, Duration: %v, Cmd: %v\n",
   log.ID, log.Time, log.Duration, log.Args)
 }
}

// ===== 健康检查 =====

func HealthCheck(ctx context.Context, rdb *redis.Client) error {
 // PING命令
 pong, err := rdb.Ping(ctx).Result()
 if err != nil {
  return err
 }
 if pong != "PONG" {
  return fmt.Errorf("unexpected ping response: %s", pong)
 }

 // 检查内存使用
 info, err := rdb.Info(ctx, "memory").Result()
 if err != nil {
  return err
 }

 // 可以解析info字符串检查具体指标
 fmt.Println("Health check passed")
 return nil
}
```

---

## 第十六部分：完整项目实战

### 实战案例15：电商库存系统

```go
package main

import (
 "context"
 "errors"
 "fmt"
 "time"

 "github.com/redis/go-redis/v9"
)

// ===== 库存系统 =====

type InventoryService struct {
 rdb *redis.Client
}

func NewInventoryService(rdb *redis.Client) *InventoryService {
 return &InventoryService{rdb: rdb}
}

// SetStock 设置库存
func (s *InventoryService) SetStock(ctx context.Context, productID string, stock int64) error {
 key := fmt.Sprintf("inventory:%s", productID)
 return s.rdb.Set(ctx, key, stock, 0).Err()
}

// GetStock 获取库存
func (s *InventoryService) GetStock(ctx context.Context, productID string) (int64, error) {
 key := fmt.Sprintf("inventory:%s", productID)
 return s.rdb.Get(ctx, key).Int64()
}

// DeductStock 扣减库存（使用Lua脚本保证原子性）
func (s *InventoryService) DeductStock(ctx context.Context, productID string, quantity int64) error {
 script := `
  local key = KEYS[1]
  local quantity = tonumber(ARGV[1])
  local stock = tonumber(redis.call('get', key) or "0")

  if stock >= quantity then
   redis.call('decrby', key, quantity)
   return 1
  else
   return 0
  end
 `

 key := fmt.Sprintf("inventory:%s", productID)
 result, err := s.rdb.Eval(ctx, script, []string{key}, quantity).Int64()
 if err != nil {
  return err
 }

 if result == 0 {
  return errors.New("insufficient stock")
 }

 return nil
}

// PreemptStock 预占库存（秒杀场景）
func (s *InventoryService) PreemptStock(ctx context.Context, productID, orderID string, quantity int64, ttl time.Duration) error {
 // 1. 扣减总库存
 if err := s.DeductStock(ctx, productID, quantity); err != nil {
  return err
 }

 // 2. 记录预占（设置过期时间）
 preemptKey := fmt.Sprintf("preempt:%s:%s", productID, orderID)
 err := s.rdb.Set(ctx, preemptKey, quantity, ttl).Err()
 if err != nil {
  // 回滚库存
  s.AddStock(ctx, productID, quantity)
  return err
 }

 return nil
}

// ConfirmStock 确认预占（支付成功）
func (s *InventoryService) ConfirmStock(ctx context.Context, productID, orderID string) error {
 preemptKey := fmt.Sprintf("preempt:%s:%s", productID, orderID)
 
 // 删除预占记录
 deleted, err := s.rdb.Del(ctx, preemptKey).Result()
 if err != nil {
  return err
 }

 if deleted == 0 {
  return errors.New("preemption not found or expired")
 }

 return nil
}

// CancelStock 取消预占（支付失败或超时）
func (s *InventoryService) CancelStock(ctx context.Context, productID, orderID string) error {
 preemptKey := fmt.Sprintf("preempt:%s:%s", productID, orderID)

 // 获取预占数量
 quantity, err := s.rdb.Get(ctx, preemptKey).Int64()
 if err == redis.Nil {
  return errors.New("preemption not found")
 }
 if err != nil {
  return err
 }

 // 回滚库存
 if err := s.AddStock(ctx, productID, quantity); err != nil {
  return err
 }

 // 删除预占记录
 s.rdb.Del(ctx, preemptKey)

 return nil
}

// AddStock 增加库存
func (s *InventoryService) AddStock(ctx context.Context, productID string, quantity int64) error {
 key := fmt.Sprintf("inventory:%s", productID)
 return s.rdb.IncrBy(ctx, key, quantity).Err()
}

// ===== 使用示例 =====

func InventoryExample() {
 rdb := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
 })
 defer rdb.Close()

 ctx := context.Background()
 service := NewInventoryService(rdb)

 // 设置初始库存
 service.SetStock(ctx, "product123", 1000)

 // 秒杀场景：预占库存
 orderID := "order001"
 if err := service.PreemptStock(ctx, "product123", orderID, 1, 15*time.Minute); err != nil {
  log.Printf("Preempt failed: %v", err)
  return
 }

 // 支付成功：确认预占
 if err := service.ConfirmStock(ctx, "product123", orderID); err != nil {
  log.Printf("Confirm failed: %v", err)
  return
 }

 // 查询当前库存
 stock, _ := service.GetStock(ctx, "product123")
 fmt.Printf("Current stock: %d\n", stock)
}
```

---

## 🎯 总结

### Redis核心要点

1. **数据结构** - String/Hash/List/Set/ZSet + 高级结构
2. **缓存策略** - 缓存穿透/击穿/雪崩解决方案
3. **分布式锁** - SETNX + Lua脚本 + Redlock
4. **消息队列** - List/Stream + 发布订阅
5. **流处理** - Redis Stream（消费者组）
6. **事务** - MULTI/EXEC + WATCH乐观锁
7. **Lua脚本** - 原子操作 + 复杂逻辑
8. **持久化** - RDB快照 + AOF日志
9. **高可用** - 主从复制 + 哨兵 + 集群
10. **性能优化** - Pipeline + 连接池 + 慢查询

### 最佳实践清单

```text
✅ 合理设置过期时间
✅ 使用Pipeline批量操作
✅ 避免大key（单个key < 1MB）
✅ 使用Hash代替多个String
✅ 实施缓存预热
✅ 防止缓存穿透/击穿/雪崩
✅ 监控慢查询
✅ 配置持久化策略
✅ 使用连接池
✅ 实施多级缓存
```

---

**文档版本**: v17.0  

<div align="center">

Made with ❤️ for High-Performance System Developers

[⬆ 回到顶部](#回到顶部)

</div>

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025-10-29  
**文档状态**: 完成  
**适用版本**: Go 1.25.3+
