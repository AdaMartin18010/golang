package pool

import (
	"strings"
	"testing"
)

func TestSimplePool(t *testing.T) {
	pool := NewSimplePool[string](func() string {
		return "new"
	})
	
	// 获取对象
	item := pool.Get()
	if item != "new" {
		t.Errorf("Expected 'new', got %s", item)
	}
	
	// 归还对象
	pool.Put("test")
	
	// 再次获取
	item = pool.Get()
	if item != "test" {
		t.Errorf("Expected 'test', got %s", item)
	}
	
	// 清空
	pool.Clear()
	if pool.Size() != 0 {
		t.Error("Expected size 0 after clear")
	}
}

func TestBoundedPool(t *testing.T) {
	pool := NewBoundedPool[string](3, func() string {
		return "new"
	})
	
	// 获取对象
	item := pool.Get()
	if item != "new" {
		t.Errorf("Expected 'new', got %s", item)
	}
	
	// 归还对象
	pool.Put("test")
	
	// 再次获取
	item = pool.Get()
	if item != "test" {
		t.Errorf("Expected 'test', got %s", item)
	}
	
	// 测试容量
	if pool.Capacity() != 3 {
		t.Errorf("Expected capacity 3, got %d", pool.Capacity())
	}
}

func TestBufferPool(t *testing.T) {
	pool := NewBufferPool()
	
	// 获取缓冲区
	buf := pool.Get()
	if len(buf) != 0 {
		t.Error("Expected empty buffer")
	}
	
	// 使用缓冲区
	buf = append(buf, []byte("test")...)
	
	// 归还缓冲区
	pool.Put(buf)
	
	// 再次获取
	buf = pool.Get()
	if len(buf) != 0 {
		t.Error("Expected reset buffer")
	}
}

func TestStringBuilderPool(t *testing.T) {
	pool := NewStringBuilderPool()
	
	// 获取字符串构建器
	sb := pool.Get()
	sb.WriteString("test")
	
	// 归还字符串构建器
	pool.Put(sb)
	
	// 再次获取
	sb = pool.Get()
	if sb.Len() != 0 {
		t.Error("Expected reset string builder")
	}
}

