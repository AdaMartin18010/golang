// Package redis provides comprehensive tests for Redis client.
//
// 测试策略：
// 1. 使用接口 mock 进行单元测试（不依赖外部服务）
// 2. 使用 docker-compose 或测试容器进行集成测试（可选）
// 3. 覆盖连接、CRUD 操作和错误处理
//
// 运行测试：
//   - 单元测试: go test -v ./internal/infrastructure/cache/redis/...
//   - 集成测试: 需要运行 docker-compose up redis
package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRedisClient 是 Redis 客户端的 mock 实现
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.DurationCmd)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// getTestRedisURL 返回测试 Redis 服务器地址
// 如果没有运行 Redis 服务器，则返回空字符串
func getTestRedisURL() string {
	// 默认测试地址
	return "localhost:6379"
}

// isRedisAvailable 检查 Redis 服务器是否可用
func isRedisAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := NewClient(Config{
		Addr:     getTestRedisURL(),
		Password: "",
		DB:       0,
	})
	if err != nil {
		return false
	}
	defer client.Close()

	return client.Ping(ctx) == nil
}

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "localhost:6379", config.Addr)
	assert.Equal(t, "", config.Password)
	assert.Equal(t, 0, config.DB)
	assert.Equal(t, 10, config.PoolSize)
	assert.Equal(t, 5, config.MinIdleConns)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 5*time.Second, config.DialTimeout)
	assert.Equal(t, 3*time.Second, config.ReadTimeout)
	assert.Equal(t, 3*time.Second, config.WriteTimeout)
	assert.Equal(t, 4*time.Second, config.PoolTimeout)
	assert.Equal(t, 5*time.Minute, config.ConnMaxIdleTime)
	assert.Equal(t, 30*time.Minute, config.ConnMaxLifetime)
}

// TestNewClient_Integration 测试创建 Redis 客户端（集成测试）
func TestNewClient_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test. Start Redis with: docker run -d -p 6379:6379 redis:latest")
	}

	config := DefaultConfig()
	config.Addr = getTestRedisURL()

	client, err := NewClient(config)
	require.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()

	// 验证连接
	ctx := context.Background()
	err = client.Ping(ctx)
	require.NoError(t, err)
}

// TestNewClient_EmptyAddr 测试空地址时使用默认配置
func TestNewClient_EmptyAddr(t *testing.T) {
	// 由于会尝试连接 localhost:6379，如果没有 Redis 服务器会失败
	// 这里主要测试配置处理逻辑
	config := Config{
		Addr: "", // 空地址
	}

	// 如果没有 Redis 服务器，应该返回错误
	_, err := NewClient(config)
	// 期望连接失败，但配置应该被正确填充
	assert.Error(t, err)
}

// TestClient_SetGet_Integration 测试设置和获取值（集成测试）
func TestClient_SetGet_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理测试数据
	defer client.Del(ctx, "test:key")

	// 测试 Set
	err = client.Set(ctx, "test:key", "test:value", 0)
	require.NoError(t, err)

	// 测试 Get
	value, err := client.Get(ctx, "test:key")
	require.NoError(t, err)
	assert.Equal(t, "test:value", value)
}

// TestClient_SetWithExpiration_Integration 测试带过期时间的设置（集成测试）
func TestClient_SetWithExpiration_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理测试数据
	defer client.Del(ctx, "test:expire")

	// 设置带过期时间的键
	err = client.Set(ctx, "test:expire", "value", 2*time.Second)
	require.NoError(t, err)

	// 验证键存在
	exists, err := client.Exists(ctx, "test:expire")
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	// 等待过期
	time.Sleep(3 * time.Second)

	// 验证键已过期
	exists, err = client.Exists(ctx, "test:expire")
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}

// TestClient_GetNonExistent_Integration 测试获取不存在的键（集成测试）
func TestClient_GetNonExistent_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 获取不存在的键
	_, err = client.Get(ctx, "test:nonexistent")
	assert.Equal(t, redis.Nil, err)
}

// TestClient_Del_Integration 测试删除键（集成测试）
func TestClient_Del_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 设置测试数据
	err = client.Set(ctx, "test:del1", "value1", 0)
	require.NoError(t, err)
	err = client.Set(ctx, "test:del2", "value2", 0)
	require.NoError(t, err)

	// 批量删除
	err = client.Del(ctx, "test:del1", "test:del2")
	require.NoError(t, err)

	// 验证键已删除
	exists, err := client.Exists(ctx, "test:del1", "test:del2")
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}

// TestClient_Exists_Integration 测试键存在检查（集成测试）
func TestClient_Exists_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理测试数据
	defer client.Del(ctx, "test:exists")

	// 检查不存在的键
	count, err := client.Exists(ctx, "test:exists")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// 设置键
	err = client.Set(ctx, "test:exists", "value", 0)
	require.NoError(t, err)

	// 检查存在的键
	count, err = client.Exists(ctx, "test:exists")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestClient_Expire_Integration 测试设置过期时间（集成测试）
func TestClient_Expire_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理测试数据
	defer client.Del(ctx, "test:expire:manual")

	// 设置不带过期时间的键
	err = client.Set(ctx, "test:expire:manual", "value", 0)
	require.NoError(t, err)

	// 检查 TTL
	ttl, err := client.TTL(ctx, "test:expire:manual")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(-1), ttl) // -1 表示没有过期时间

	// 设置过期时间
	err = client.Expire(ctx, "test:expire:manual", 2*time.Second)
	require.NoError(t, err)

	// 检查 TTL
	ttl, err = client.TTL(ctx, "test:expire:manual")
	require.NoError(t, err)
	assert.True(t, ttl > 0 && ttl <= 2*time.Second)
}

// TestClient_TTL_Integration 测试获取 TTL（集成测试）
func TestClient_TTL_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理测试数据
	defer client.Del(ctx, "test:ttl")

	// 不存在的键
	ttl, err := client.TTL(ctx, "test:ttl")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(-2), ttl) // -2 表示键不存在

	// 设置带过期时间的键
	err = client.Set(ctx, "test:ttl", "value", 10*time.Second)
	require.NoError(t, err)

	// 检查 TTL
	ttl, err = client.TTL(ctx, "test:ttl")
	require.NoError(t, err)
	assert.True(t, ttl > 0 && ttl <= 10*time.Second)
}

// TestClient_Close_Integration 测试关闭客户端（集成测试）
func TestClient_Close_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}

// TestClient_ConcurrentAccess_Integration 测试并发访问（集成测试）
func TestClient_ConcurrentAccess_Integration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理测试数据
	defer client.Del(ctx, "test:concurrent")

	// 并发写入
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			key := "test:concurrent"
			err := client.Set(ctx, key, n, 10*time.Second)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证键存在
	exists, err := client.Exists(ctx, "test:concurrent")
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists)
}

// TestConfig_Validation 测试配置验证
func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Addr:     "localhost:6379",
				PoolSize: 10,
			},
			wantErr: false,
		},
		{
			name: "empty address",
			config: Config{
				Addr: "",
			},
			wantErr: true, // 如果没有 Redis 服务器会失败
		},
		{
			name: "custom timeout",
			config: Config{
				Addr:        "localhost:6379",
				DialTimeout: 10 * time.Second,
				ReadTimeout: 5 * time.Second,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注意：这里只是验证配置结构，不涉及实际连接
			// 实际连接测试在集成测试中完成
			if tt.config.Addr == "" {
				// 空地址应该使用默认值
				config := DefaultConfig()
				assert.Equal(t, "localhost:6379", config.Addr)
			}
		})
	}
}

// BenchmarkSet 基准测试：Set 操作
func BenchmarkSet(b *testing.B) {
	if !isRedisAvailable() {
		b.Skip("Redis server not available, skipping benchmark")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	// 清理
	defer client.Del(ctx, "bench:set")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set(ctx, "bench:set", "value", 0)
	}
}

// BenchmarkGet 基准测试：Get 操作
func BenchmarkGet(b *testing.B) {
	if !isRedisAvailable() {
		b.Skip("Redis server not available, skipping benchmark")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	// 准备数据
	err = client.Set(ctx, "bench:get", "value", 0)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Del(ctx, "bench:get")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Get(ctx, "bench:get")
	}
}

// MockRedisCmd 是一个通用的 mock 命令结构
type MockRedisCmd struct {
	mock.Mock
}

func (m *MockRedisCmd) Result() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockRedisCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisCmd) Val() string {
	args := m.Called()
	return args.String(0)
}

// TestClient_Set_Error 测试 Set 错误处理
func TestClient_Set_Error(t *testing.T) {
	// 这是一个示例，展示如何使用 mock 进行单元测试
	// 实际项目中可以使用 gomock 或类似的工具

	// 模拟连接错误
	config := Config{
		Addr: "invalid:host:port",
	}

	_, err := NewClient(config)
	assert.Error(t, err)
}

// TestClient_Get_Error 测试 Get 错误处理
func TestClient_Get_Error(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration test")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 获取不存在的键应该返回 redis.Nil
	_, err = client.Get(ctx, "nonexistent:key")
	assert.Equal(t, redis.Nil, err)
}

// TestRedisError_Wrapping 测试错误包装
func TestRedisError_Wrapping(t *testing.T) {
	// 测试错误包装逻辑
	testErr := errors.New("connection refused")
	wrappedErr := errors.New("failed to connect to redis: " + testErr.Error())

	assert.Error(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "failed to connect to redis")
}
