package benchmarks

import (
	"fmt"
	"testing"
	"time"

	"distributed-cache/internal/cache"
	"distributed-cache/internal/ring"
)

func BenchmarkCacheGet(b *testing.B) {
	config := &cache.Config{
		MaxSize:        1024 * 1024 * 100, // 100MB
		EvictionPolicy: cache.LRU,
	}
	c := cache.New(config)

	// Pre-populate cache
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		c.Set(key, value, time.Hour)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%10000)
			c.Get(key)
			i++
		}
	})
}

func BenchmarkCacheSet(b *testing.B) {
	config := &cache.Config{
		MaxSize:        1024 * 1024 * 100,
		EvictionPolicy: cache.LRU,
	}
	c := cache.New(config)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i)
			value := []byte(fmt.Sprintf("value-%d", i))
			c.Set(key, value, time.Hour)
			i++
		}
	})
}

func BenchmarkConsistentHashRing(b *testing.B) {
	r := ring.New(150)

	// Add nodes
	for i := 0; i < 10; i++ {
		node := &ring.Node{
			ID:      fmt.Sprintf("node-%d", i),
			Address: fmt.Sprintf("192.168.1.%d:8080", i),
			Weight:  1,
		}
		r.AddNode(node)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i)
			r.GetNode(key)
			i++
		}
	})
}

func BenchmarkCacheWithEviction(b *testing.B) {
	config := &cache.Config{
		MaxSize:        1024 * 1024, // 1MB - small to trigger eviction
		EvictionPolicy: cache.LRU,
	}
	c := cache.New(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := make([]byte, 1024) // 1KB values
		c.Set(key, value, time.Hour)
	}
}
