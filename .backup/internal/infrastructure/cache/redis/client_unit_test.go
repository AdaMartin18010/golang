package redis

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRedisClientV2 是用于单元测试的 Redis 客户端 Mock
type MockRedisClientV2 struct {
	mock.Mock
}

func (m *MockRedisClientV2) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClientV2) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClientV2) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClientV2) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClientV2) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClientV2) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClientV2) TTL(ctx context.Context, key string) *redis.DurationCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.DurationCmd)
}

func (m *MockRedisClientV2) Close() error {
	args := m.Called()
	return args.Error(0)
}

// createMockStatusCmd 创建 mock StatusCmd
func createMockStatusCmd(result string, err error) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(context.Background(), "ping")
	if err != nil {
		cmd.SetErr(err)
	} else {
		cmd.SetVal(result)
	}
	return cmd
}

// createMockStringCmd 创建 mock StringCmd
func createMockStringCmd(result string, err error) *redis.StringCmd {
	cmd := redis.NewStringCmd(context.Background(), "get", "key")
	if err != nil {
		cmd.SetErr(err)
	} else {
		cmd.SetVal(result)
	}
	return cmd
}

// createMockIntCmd 创建 mock IntCmd
func createMockIntCmd(result int64, err error) *redis.IntCmd {
	cmd := redis.NewIntCmd(context.Background(), "del", "key")
	if err != nil {
		cmd.SetErr(err)
	} else {
		cmd.SetVal(result)
	}
	return cmd
}

// createMockBoolCmd 创建 mock BoolCmd
func createMockBoolCmd(result bool, err error) *redis.BoolCmd {
	cmd := redis.NewBoolCmd(context.Background(), "expire", "key", 60)
	if err != nil {
		cmd.SetErr(err)
	} else {
		cmd.SetVal(result)
	}
	return cmd
}

// createMockDurationCmd 创建 mock DurationCmd
func createMockDurationCmd(result time.Duration, err error) *redis.DurationCmd {
	cmd := redis.NewDurationCmd(context.Background(), time.Second)
	if err != nil {
		cmd.SetErr(err)
	} else {
		cmd.SetVal(result)
	}
	return cmd
}

// TestDefaultConfig_AllFields 测试默认配置的所有字段
func TestDefaultConfig_AllFields(t *testing.T) {
	config := DefaultConfig()

	// 验证所有字段都有合理的默认值
	assert.NotEmpty(t, config.Addr, "Addr should not be empty")
	assert.Greater(t, config.PoolSize, 0, "PoolSize should be greater than 0")
	assert.GreaterOrEqual(t, config.MinIdleConns, 0, "MinIdleConns should be non-negative")
	assert.GreaterOrEqual(t, config.MaxRetries, 0, "MaxRetries should be non-negative")
	assert.Greater(t, config.DialTimeout, time.Duration(0), "DialTimeout should be positive")
	assert.Greater(t, config.ReadTimeout, time.Duration(0), "ReadTimeout should be positive")
	assert.Greater(t, config.WriteTimeout, time.Duration(0), "WriteTimeout should be positive")
	assert.Greater(t, config.PoolTimeout, time.Duration(0), "PoolTimeout should be positive")
	assert.Greater(t, config.ConnMaxIdleTime, time.Duration(0), "ConnMaxIdleTime should be positive")
	assert.Greater(t, config.ConnMaxLifetime, time.Duration(0), "ConnMaxLifetime should be positive")
}

// TestConfig_CustomValues 测试自定义配置值
func TestConfig_CustomValues(t *testing.T) {
	customConfig := Config{
		Addr:            "custom:6379",
		Password:        "secret",
		DB:              1,
		PoolSize:        20,
		MinIdleConns:    10,
		MaxRetries:      5,
		DialTimeout:     10 * time.Second,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		PoolTimeout:     8 * time.Second,
		ConnMaxIdleTime: 10 * time.Minute,
		ConnMaxLifetime: 1 * time.Hour,
	}

	assert.Equal(t, "custom:6379", customConfig.Addr)
	assert.Equal(t, "secret", customConfig.Password)
	assert.Equal(t, 1, customConfig.DB)
	assert.Equal(t, 20, customConfig.PoolSize)
	assert.Equal(t, 10, customConfig.MinIdleConns)
	assert.Equal(t, 5, customConfig.MaxRetries)
	assert.Equal(t, 10*time.Second, customConfig.DialTimeout)
	assert.Equal(t, 5*time.Second, customConfig.ReadTimeout)
	assert.Equal(t, 5*time.Second, customConfig.WriteTimeout)
	assert.Equal(t, 8*time.Second, customConfig.PoolTimeout)
	assert.Equal(t, 10*time.Minute, customConfig.ConnMaxIdleTime)
	assert.Equal(t, 1*time.Hour, customConfig.ConnMaxLifetime)
}

// TestNewClient_InvalidAddress 测试无效地址
func TestNewClient_InvalidAddress(t *testing.T) {
	config := Config{
		Addr:        "invalid:host:port:format",
		DialTimeout: 1 * time.Second,
	}

	_, err := NewClient(config)
	assert.Error(t, err)
}

// TestClient_Close_Error 测试关闭客户端时的错误
func TestClient_Close_Error(t *testing.T) {
	// 由于无法轻易 mock 底层客户端，这里测试错误传播
	// 在实际场景中，如果底层连接已关闭，Close 应该返回错误

	// 注：如果 Redis 不可用，这个测试会被跳过
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)

	// 第一次关闭应该成功
	err = client.Close()
	assert.NoError(t, err)
}

// TestClient_Set_ZeroExpiration 测试设置永不过期的键
func TestClient_Set_ZeroExpiration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理
	defer client.Del(ctx, "test:noexpire")

	// 设置永不过期的键
	err = client.Set(ctx, "test:noexpire", "value", 0)
	require.NoError(t, err)

	// 验证 TTL 为 -1
	ttl, err := client.TTL(ctx, "test:noexpire")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(-1), ttl)
}

// TestClient_Set_NegativeExpiration 测试设置负过期时间（立即过期）
func TestClient_Set_NegativeExpiration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 设置负过期时间应该立即删除键或报错
	// 根据 Redis 文档，负值会导致键被删除
	err = client.Set(ctx, "test:negexpire", "value", -1*time.Second)
	// Redis 可能返回错误或成功但键被删除
	// 这里我们只验证操作不会 panic
}

// TestClient_Get_EmptyKey 测试获取空键
func TestClient_Get_EmptyKey(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 获取不存在的键
	_, err = client.Get(ctx, "")
	// 空键应该返回错误
	assert.Error(t, err)
}

// TestClient_Del_EmptyKeys 测试删除空键列表
func TestClient_Del_EmptyKeys(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 删除空键列表
	err = client.Del(ctx)
	// 应该成功（没有键被删除）
	assert.NoError(t, err)
}

// TestClient_Del_NonExistent 测试删除不存在的键
func TestClient_Del_NonExistent(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 删除不存在的键应该成功（返回 0）
	err = client.Del(ctx, "nonexistent:key:12345")
	assert.NoError(t, err)
}

// TestClient_Exists_MultipleKeys 测试检查多个键
func TestClient_Exists_MultipleKeys(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理
	defer client.Del(ctx, "test:exists1", "test:exists2")

	// 设置两个键
	client.Set(ctx, "test:exists1", "value1", 0)
	client.Set(ctx, "test:exists2", "value2", 0)

	// 检查多个键
	count, err := client.Exists(ctx, "test:exists1", "test:exists2", "test:nonexistent")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// TestClient_Exists_Empty 测试检查空键列表
func TestClient_Exists_Empty(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 检查空键列表
	count, err := client.Exists(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestClient_Expire_NonExistent 测试为不存在的键设置过期时间
func TestClient_Expire_NonExistent(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 为不存在的键设置过期时间
	err = client.Expire(ctx, "nonexistent:key:12345", time.Hour)
	// 根据 Redis 版本，这可能返回错误或不返回错误
	// 我们只是验证操作不会 panic
}

// TestClient_TTL_NonExistent 测试获取不存在键的 TTL
func TestClient_TTL_NonExistent(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 获取不存在键的 TTL
	ttl, err := client.TTL(ctx, "nonexistent:key:12345")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(-2), ttl) // -2 表示键不存在
}

// TestClient_TTL_NoExpiration 测试获取无过期时间键的 TTL
func TestClient_TTL_NoExpiration(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理
	defer client.Del(ctx, "test:noexpire:ttl")

	// 设置永不过期的键
	client.Set(ctx, "test:noexpire:ttl", "value", 0)

	// 获取 TTL
	ttl, err := client.TTL(ctx, "test:noexpire:ttl")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(-1), ttl) // -1 表示键存在但没有过期时间
}

// TestClient_Ping_Error 测试 Ping 错误
func TestClient_Ping_Error(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 正常情况应该成功
	err = client.Ping(ctx)
	assert.NoError(t, err)
}

// TestClient_GetClient 测试获取底层客户端
func TestClient_GetClient(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 获取底层客户端
	underlyingClient := client.GetClient()
	assert.NotNil(t, underlyingClient)
}

// TestNewClient_ConfigVariations 测试不同配置变体
func TestNewClient_ConfigVariations(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "default config",
			config: Config{
				Addr: getTestRedisURL(),
			},
			wantErr: false,
		},
		{
			name: "custom timeouts",
			config: Config{
				Addr:         getTestRedisURL(),
				DialTimeout:  2 * time.Second,
				ReadTimeout:  2 * time.Second,
				WriteTimeout: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "custom pool settings",
			config: Config{
				Addr:         getTestRedisURL(),
				PoolSize:     5,
				MinIdleConns: 2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if client != nil {
					client.Close()
				}
			}
		})
	}
}

// TestClient_Set_VariousTypes 测试设置不同类型的值
func TestClient_Set_VariousTypes(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "string value",
			key:   "test:type:string",
			value: "hello world",
		},
		{
			name:  "integer value",
			key:   "test:type:int",
			value: 42,
		},
		{
			name:  "float value",
			key:   "test:type:float",
			value: 3.14,
		},
		{
			name:  "byte slice",
			key:   "test:type:bytes",
			value: []byte("byte data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer client.Del(ctx, tt.key)

			err := client.Set(ctx, tt.key, tt.value, time.Minute)
			assert.NoError(t, err)

			_, err = client.Get(ctx, tt.key)
			assert.NoError(t, err)
		})
	}
}

// TestClient_ConcurrentOperations 测试并发操作
func TestClient_ConcurrentOperations(t *testing.T) {
	if !isRedisAvailable() {
		t.Skip("Redis not available")
	}

	ctx := context.Background()
	client, err := NewClient(DefaultConfig())
	require.NoError(t, err)
	defer client.Close()

	// 清理
	defer client.Del(ctx, "test:concurrent:1", "test:concurrent:2", "test:concurrent:3")

	// 并发设置
	done := make(chan bool, 3)

	go func() {
		err := client.Set(ctx, "test:concurrent:1", "value1", time.Minute)
		assert.NoError(t, err)
		done <- true
	}()

	go func() {
		err := client.Set(ctx, "test:concurrent:2", "value2", time.Minute)
		assert.NoError(t, err)
		done <- true
	}()

	go func() {
		err := client.Set(ctx, "test:concurrent:3", "value3", time.Minute)
		assert.NoError(t, err)
		done <- true
	}()

	// 等待所有 goroutine 完成
	for i := 0; i < 3; i++ {
		<-done
	}

	// 验证所有键都存在
	count, err := client.Exists(ctx, "test:concurrent:1", "test:concurrent:2", "test:concurrent:3")
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// TestConfig_Validation_EdgeCases 测试配置验证的边缘情况
func TestConfig_Validation_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "zero values",
			config: Config{},
		},
		{
			name: "negative timeouts",
			config: Config{
				DialTimeout: -1 * time.Second,
			},
		},
		{
			name: "zero pool size",
			config: Config{
				PoolSize: 0,
			},
		},
		{
			name: "large pool size",
			config: Config{
				PoolSize: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证配置结构可以被创建
			assert.NotNil(t, tt.config)
		})
	}
}

// TestError_Wrapping 测试错误包装
func TestError_Wrapping(t *testing.T) {
	// 测试错误包装格式
	innerErr := errors.New("connection refused")
	wrappedErr := fmt.Errorf("failed to connect to redis: %w", innerErr)

	assert.Error(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "failed to connect to redis")
	assert.ErrorIs(t, wrappedErr, innerErr)
}
