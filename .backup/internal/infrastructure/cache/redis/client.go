// Package redis provides a Redis client wrapper for caching and data storage.
//
// Redis 是一个开源的内存数据结构存储系统，常用作缓存、消息队列和数据库。
//
// 设计原则：
// 1. 连接池管理：使用连接池提高性能和资源利用率
// 2. 超时控制：配置合理的超时时间，避免长时间阻塞
// 3. 错误处理：提供清晰的错误信息，便于问题排查
// 4. 资源管理：确保连接正确关闭，避免资源泄漏
//
// 核心功能：
// - 键值存储：Set、Get、Del 等基本操作
// - 过期管理：Expire、TTL 等过期时间管理
// - 连接管理：连接池、超时、重试等配置
//
// 使用场景：
// - 缓存：存储热点数据，减少数据库压力
// - 会话存储：存储用户会话信息
// - 分布式锁：实现分布式锁机制
// - 计数器：实现计数器和限流功能
// - 消息队列：使用 Redis 的发布订阅功能
//
// 示例：
//
//	// 创建客户端
//	config := redis.DefaultConfig()
//	config.Addr = "localhost:6379"
//	client, err := redis.NewClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// 设置和获取值
//	ctx := context.Background()
//	err = client.Set(ctx, "key", "value", time.Hour)
//	value, err := client.Get(ctx, "key")
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config 是 Redis 客户端的配置结构。
//
// 配置说明：
// - 所有超时时间使用 time.Duration 类型
// - 连接池大小应根据应用负载调整
// - 生产环境应配置密码和 TLS
//
// 字段说明：
// - Addr: Redis 服务器地址，格式为 "host:port"
// - Password: Redis 服务器密码，如果为空则不使用密码
// - DB: 数据库编号，默认为 0（Redis 支持多个逻辑数据库）
// - PoolSize: 连接池大小，控制最大并发连接数
// - MinIdleConns: 最小空闲连接数，保持一定数量的空闲连接
// - MaxRetries: 最大重试次数，网络错误时自动重试
// - DialTimeout: 建立连接的超时时间
// - ReadTimeout: 读取操作的超时时间
// - WriteTimeout: 写入操作的超时时间
// - PoolTimeout: 从连接池获取连接的超时时间
// - ConnMaxIdleTime: 连接最大空闲时间，超过此时间后关闭连接
// - ConnMaxLifetime: 连接最大生存时间，超过此时间后关闭连接
type Config struct {
	Addr            string        // 地址（例如: "localhost:6379"）
	Password        string        // 密码
	DB              int           // 数据库编号
	PoolSize        int           // 连接池大小
	MinIdleConns    int           // 最小空闲连接数
	MaxRetries      int           // 最大重试次数
	DialTimeout     time.Duration // 连接超时
	ReadTimeout     time.Duration // 读取超时
	WriteTimeout    time.Duration // 写入超时
	PoolTimeout     time.Duration // 连接池超时
	ConnMaxIdleTime time.Duration // 连接最大空闲时间
	ConnMaxLifetime time.Duration // 连接最大生存时间
}

// DefaultConfig 返回 Redis 客户端的默认配置。
//
// 默认配置说明：
// - 适用于开发和测试环境
// - 生产环境应根据实际情况调整
// - 连接池大小和超时时间应根据负载调整
//
// 返回：
// - Config: 包含默认值的配置结构
//
// 使用示例：
//
//	config := redis.DefaultConfig()
//	config.Addr = "redis.example.com:6379"
//	config.Password = "your-password"
//	client, err := redis.NewClient(config)
func DefaultConfig() Config {
	return Config{
		Addr:            "localhost:6379",
		Password:        "",
		DB:              0,
		PoolSize:        10,
		MinIdleConns:    5,
		MaxRetries:      3,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
	}
}

// Client 是 Redis 客户端的封装，提供简化的 API 接口。
//
// 功能说明：
// - 封装了 go-redis 客户端
// - 提供常用的 Redis 操作接口
// - 支持连接池管理和超时控制
//
// 设计说明：
// - 线程安全：可以在多个 goroutine 中并发使用
// - 连接复用：使用连接池复用连接，提高性能
// - 自动重试：网络错误时自动重试
//
// 使用示例：
//
//	client, err := redis.NewClient(redis.DefaultConfig())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	ctx := context.Background()
//	err = client.Set(ctx, "key", "value", time.Hour)
type Client struct {
	client *redis.Client
}

// NewClient 创建并初始化 Redis 客户端。
//
// 功能说明：
// - 根据配置创建 Redis 客户端
// - 建立连接池
// - 测试连接是否可用
//
// 参数：
//   - config: Redis 客户端配置
//     如果 Addr 为空，则使用默认配置
//
// 返回：
// - *Client: 配置好的客户端实例
// - error: 如果连接失败，返回错误信息
//
// 使用示例：
//
//	config := redis.Config{
//	    Addr:     "localhost:6379",
//	    Password: "password",
//	    DB:       0,
//	}
//	client, err := redis.NewClient(config)
//	if err != nil {
//	    log.Fatal("Failed to create Redis client:", err)
//	}
//	defer client.Close()
//
// 注意事项：
// - 创建客户端时会测试连接，确保 Redis 服务器可用
// - 应在应用程序生命周期中复用客户端实例
// - 退出前应调用 Close() 关闭客户端
func NewClient(config Config) (*Client, error) {
	if config.Addr == "" {
		config = DefaultConfig()
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:            config.Addr,
		Password:        config.Password,
		DB:              config.DB,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		MaxRetries:      config.MaxRetries,
		DialTimeout:     config.DialTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
		PoolTimeout:     config.PoolTimeout,
		ConnMaxIdleTime: config.ConnMaxIdleTime,
		ConnMaxLifetime: config.ConnMaxLifetime,
	})

	// 测试连接
	// 确保 Redis 服务器可用，避免后续操作失败
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{client: rdb}, nil
}

// Close 关闭 Redis 客户端连接。
//
// 功能说明：
// - 关闭所有连接池中的连接
// - 停止所有后台任务
// - 释放客户端资源
//
// 返回：
// - error: 如果关闭过程中出现错误，返回错误信息
//
// 使用示例：
//
//	defer client.Close()
//
// 注意事项：
// - 应在应用程序退出前调用
// - 关闭后不应再使用该客户端
// - 关闭是同步的，会等待所有操作完成
func (c *Client) Close() error {
	return c.client.Close()
}

// GetClient 返回底层的 go-redis 客户端实例。
//
// 功能说明：
// - 提供对底层客户端的访问
// - 用于需要使用 go-redis 高级功能的场景
//
// 返回：
// - *redis.Client: 底层的 go-redis 客户端
//
// 使用场景：
// - 需要使用 go-redis 的高级功能（如 Pipeline、Transaction）
// - 与第三方库集成
// - 需要直接访问底层客户端的方法
//
// 使用示例：
//
//	// 使用底层客户端执行 Pipeline
//	pipe := client.GetClient().Pipeline()
//	pipe.Set(ctx, "key1", "value1", 0)
//	pipe.Set(ctx, "key2", "value2", 0)
//	_, err := pipe.Exec(ctx)
//
// 注意事项：
// - 返回的客户端与封装客户端共享连接池
// - 不应直接关闭返回的客户端（应使用 Close() 方法）
func (c *Client) GetClient() *redis.Client {
	return c.client
}

// Ping 测试与 Redis 服务器的连接。
//
// 功能说明：
// - 发送 PING 命令到 Redis 服务器
// - 用于健康检查和连接验证
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
//
// 返回：
// - error: 如果连接失败，返回错误信息
//
// 使用场景：
// - 健康检查端点
// - 连接验证
// - 故障排查
//
// 使用示例：
//
//	ctx := context.Background()
//	if err := client.Ping(ctx); err != nil {
//	    log.Printf("Redis connection failed: %v", err)
//	}
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Set 设置键值对，并可选地设置过期时间。
//
// 功能说明：
// - 将键值对存储到 Redis
// - 如果键已存在，则覆盖旧值
// - 可以设置过期时间，过期后键自动删除
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
// - key: 键名
// - value: 值（可以是任意类型，会自动序列化）
// - expiration: 过期时间
//   - 0: 不过期
//   - > 0: 指定过期时间
//
// 返回：
// - error: 如果设置失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//	// 设置不过期的键
//	err := client.Set(ctx, "key", "value", 0)
//
//	// 设置带过期时间的键（1小时后过期）
//	err = client.Set(ctx, "key", "value", time.Hour)
//
//	// 设置复杂对象
//	user := User{ID: 1, Name: "Alice"}
//	data, _ := json.Marshal(user)
//	err = client.Set(ctx, "user:1", data, time.Hour)
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取指定键的值。
//
// 功能说明：
// - 从 Redis 获取键的值
// - 如果键不存在，返回 redis.Nil 错误
// - 返回字符串类型，需要自行反序列化
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
// - key: 键名
//
// 返回：
// - string: 键的值
// - error: 如果获取失败或键不存在，返回错误信息
//   - redis.Nil: 键不存在
//   - 其他错误: 网络错误、超时等
//
// 使用示例：
//
//	ctx := context.Background()
//	value, err := client.Get(ctx, "key")
//	if err == redis.Nil {
//	    log.Println("Key does not exist")
//	} else if err != nil {
//	    log.Printf("Error getting key: %v", err)
//	} else {
//	    log.Printf("Value: %s", value)
//	}
//
//	// 获取并反序列化 JSON
//	data, err := client.Get(ctx, "user:1")
//	if err == nil {
//	    var user User
//	    json.Unmarshal([]byte(data), &user)
//	}
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Del 删除一个或多个键。
//
// 功能说明：
// - 删除指定的键
// - 如果键不存在，则忽略
// - 支持批量删除多个键
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
// - keys: 要删除的键名（可变参数）
//
// 返回：
// - error: 如果删除失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//	// 删除单个键
//	err := client.Del(ctx, "key")
//
//	// 批量删除多个键
//	err = client.Del(ctx, "key1", "key2", "key3")
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查一个或多个键是否存在。
//
// 功能说明：
// - 检查指定的键是否存在
// - 支持检查多个键，返回存在的键的数量
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
// - keys: 要检查的键名（可变参数）
//
// 返回：
// - int64: 存在的键的数量
// - error: 如果检查失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//	// 检查单个键
//	count, err := client.Exists(ctx, "key")
//	if count > 0 {
//	    log.Println("Key exists")
//	}
//
//	// 检查多个键
//	count, err = client.Exists(ctx, "key1", "key2", "key3")
//	log.Printf("%d keys exist", count)
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire 为键设置过期时间。
//
// 功能说明：
// - 为已存在的键设置过期时间
// - 如果键不存在，返回 false
// - 过期后键会自动删除
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
// - key: 键名
// - expiration: 过期时间
//
// 返回：
// - error: 如果设置失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//	// 设置键在 1 小时后过期
//	err := client.Expire(ctx, "key", time.Hour)
//
//	// 延长键的过期时间
//	err = client.Expire(ctx, "key", 2*time.Hour)
//
// 注意事项：
// - 只能为已存在的键设置过期时间
// - 如果键已有过期时间，会覆盖旧的过期时间
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取键的剩余过期时间。
//
// 功能说明：
// - 返回键的剩余过期时间（Time To Live）
// - 如果键不存在，返回 -2
// - 如果键存在但没有设置过期时间，返回 -1
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
// - key: 键名
//
// 返回：
// - time.Duration: 剩余过期时间
//   - > 0: 剩余过期时间
//   - -1: 键存在但没有过期时间
//   - -2: 键不存在
//
// - error: 如果查询失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//	ttl, err := client.TTL(ctx, "key")
//	if ttl > 0 {
//	    log.Printf("Key will expire in %v", ttl)
//	} else if ttl == -1 {
//	    log.Println("Key exists but has no expiration")
//	} else {
//	    log.Println("Key does not exist")
//	}
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}
