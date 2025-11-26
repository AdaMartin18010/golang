package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_GetSet(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	// 设置缓存
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 获取缓存
	value, ok := cache.Get(ctx, "key1")
	if !ok {
		t.Fatal("Expected cache hit")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	// 设置短期缓存
	err := cache.Set(ctx, "key1", "value1", 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 立即获取应该成功
	_, ok := cache.Get(ctx, "key1")
	if !ok {
		t.Fatal("Expected cache hit")
	}

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 获取应该失败
	_, ok = cache.Get(ctx, "key1")
	if ok {
		t.Fatal("Expected cache miss")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Minute)
	cache.Delete(ctx, "key1")

	_, ok := cache.Get(ctx, "key1")
	if ok {
		t.Fatal("Expected cache miss")
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Minute)
	cache.Set(ctx, "key2", "value2", time.Minute)

	cache.Clear(ctx)

	if cache.Size() != 0 {
		t.Errorf("Expected size 0, got %d", cache.Size())
	}
}

func TestGetOrSet(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	callCount := 0
	fn := func() (interface{}, error) {
		callCount++
		return "value1", nil
	}

	// 第一次调用应该执行函数
	value, err := GetOrSet(ctx, cache, "key1", time.Minute, fn)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// 第二次调用应该使用缓存
	value, err = GetOrSet(ctx, cache, "key1", time.Minute, fn)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}
