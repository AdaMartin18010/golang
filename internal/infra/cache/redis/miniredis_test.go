// Package redis provides unit tests using miniredis (in-memory Redis server).
// This file contains fast unit tests that don't require a real Redis server.
package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMiniredis creates a new miniredis server and returns its address
func setupMiniredis(t *testing.T) *miniredis.Miniredis {
	mr := miniredis.RunT(t)
	return mr
}

// TestDefaultConfig_Miniredis 测试默认配置
func TestDefaultConfig_Miniredis(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "localhost:6379", config.Addr)
	assert.Equal(t, "", config.Password)
	assert.Equal(t, 0, config.DB)
	assert.Equal(t, 10, config.PoolSize)
	assert.Equal(t, 5, config.MinIdleConns)
}

// TestNewClient_Miniredis 测试创建 Redis 客户端
func TestNewClient_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	config := Config{
		Addr:     mr.Addr(),
		Password: "",
		DB:       0,
	}

	client, err := NewClient(config)
	require.NoError(t, err)
	assert.NotNil(t, client)
	defer client.Close()

	// 验证连接
	ctx := context.Background()
	err = client.Ping(ctx)
	require.NoError(t, err)
}

// TestClient_SetGet_Miniredis 测试设置和获取值
func TestClient_SetGet_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	// 测试 Set
	err = client.Set(ctx, "test:key", "test:value", 0)
	require.NoError(t, err)

	// 测试 Get
	value, err := client.Get(ctx, "test:key")
	require.NoError(t, err)
	assert.Equal(t, "test:value", value)

	// 验证 miniredis 中确实有数据
	val, err := mr.Get("test:key")
	require.NoError(t, err)
	assert.Equal(t, "test:value", val)
}

// TestClient_SetWithExpiration_Miniredis 测试带过期时间的设置
func TestClient_SetWithExpiration_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	// 设置带过期时间的键
	err = client.Set(ctx, "test:expire", "value", 2*time.Second)
	require.NoError(t, err)

	// 验证键存在
	exists, err := client.Exists(ctx, "test:expire")
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	// 快进时间模拟过期
	mr.FastForward(3 * time.Second)

	// 验证键已过期
	exists, err = client.Exists(ctx, "test:expire")
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}

// TestClient_GetNonExistent_Miniredis 测试获取不存在的键
func TestClient_GetNonExistent_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	// 获取不存在的键
	_, err = client.Get(ctx, "test:nonexistent")
	assert.Equal(t, redis.Nil, err)
}

// TestClient_Del_Miniredis 测试删除键
func TestClient_Del_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
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

// TestClient_Exists_Miniredis 测试键存在检查
func TestClient_Exists_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	// 设置测试数据
	err = client.Set(ctx, "test:exists1", "value1", 0)
	require.NoError(t, err)
	err = client.Set(ctx, "test:exists2", "value2", 0)
	require.NoError(t, err)

	// 测试批量检查
	exists, err := client.Exists(ctx, "test:exists1", "test:exists2", "test:nonexistent")
	require.NoError(t, err)
	assert.Equal(t, int64(2), exists)
}

// TestClient_Expire_Miniredis 测试设置过期时间
func TestClient_Expire_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	// 设置测试数据
	err = client.Set(ctx, "test:expire", "value", 0)
	require.NoError(t, err)

	// 设置过期时间
	err = client.Expire(ctx, "test:expire", 2*time.Second)
	require.NoError(t, err)

	// 快进时间
	mr.FastForward(3 * time.Second)

	// 验证键已过期
	exists, err := client.Exists(ctx, "test:expire")
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}

// TestClient_TTL_Miniredis 测试获取剩余过期时间
func TestClient_TTL_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	// 设置带过期时间的测试数据
	err = client.Set(ctx, "test:ttl", "value", 10*time.Second)
	require.NoError(t, err)

	// 获取 TTL
	ttl, err := client.TTL(ctx, "test:ttl")
	require.NoError(t, err)
	// TTL 应该接近 10 秒（允许一定误差）
	assert.True(t, ttl > 8*time.Second && ttl <= 10*time.Second, "TTL should be around 10s, got %v", ttl)

	// 测试无过期时间的键
	err = client.Set(ctx, "test:noexpire", "value", 0)
	require.NoError(t, err)

	ttl, err = client.TTL(ctx, "test:noexpire")
	require.NoError(t, err)
	assert.Equal(t, time.Duration(-1), ttl) // -1 表示没有过期时间
}

// TestClient_DifferentTypes_Miniredis 测试不同类型的值
func TestClient_DifferentTypes_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name     string
		key      string
		value    interface{}
		expected string // Redis 存储的都是字符串
	}{
		{"string", "type:string", "hello", "hello"},
		{"integer", "type:integer", 42, "42"},
		{"float", "type:float", 3.14, "3.14"},
		{"json", "type:json", `{"name":"test","value":123}`, `{"name":"test","value":123}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.Set(ctx, tc.key, tc.value, 0)
			require.NoError(t, err)

			got, err := client.Get(ctx, tc.key)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// TestClient_ConcurrentOperations_Miniredis 测试并发操作
func TestClient_ConcurrentOperations_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client.Close()

	// 并发写入
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent:key:%d", i)
			value := fmt.Sprintf("value:%d", i)
			err := client.Set(ctx, key, value, 0)
			require.NoError(t, err)
		}(i)
	}

	// 等待所有 goroutine 完成
	time.Sleep(100 * time.Millisecond)

	// 验证所有键都存在
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent:key:%d", i)
		exists, err := client.Exists(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, int64(1), exists)
	}
}

// TestClient_Reconnect_Miniredis 测试重新连接
func TestClient_Reconnect_Miniredis(t *testing.T) {
	mr := setupMiniredis(t)

	ctx := context.Background()
	client, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)

	// 设置数据
	err = client.Set(ctx, "reconnect:test", "value", 0)
	require.NoError(t, err)

	// 关闭客户端
	err = client.Close()
	require.NoError(t, err)

	// 创建新客户端
	client2, err := NewClient(Config{Addr: mr.Addr()})
	require.NoError(t, err)
	defer client2.Close()

	// 验证数据仍然存在
	value, err := client2.Get(ctx, "reconnect:test")
	require.NoError(t, err)
	assert.Equal(t, "value", value)
}

// BenchmarkClient_Set_Miniredis Set 操作性能测试
func BenchmarkClient_Set_Miniredis(b *testing.B) {
	mr, err := miniredis.Run()
	if err != nil {
		b.Fatal(err)
	}
	defer mr.Close()

	ctx := context.Background()
	client, _ := NewClient(Config{Addr: mr.Addr()})
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set(ctx, "bench:key", "value", 0)
	}
}

// BenchmarkClient_Get_Miniredis Get 操作性能测试
func BenchmarkClient_Get_Miniredis(b *testing.B) {
	mr, err := miniredis.Run()
	if err != nil {
		b.Fatal(err)
	}
	defer mr.Close()

	ctx := context.Background()
	client, _ := NewClient(Config{Addr: mr.Addr()})
	defer client.Close()

	client.Set(ctx, "bench:key", "value", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Get(ctx, "bench:key")
	}
}
