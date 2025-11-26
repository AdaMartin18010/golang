package lru

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	lru := NewLRUCache[string, int](3)
	
	// 添加元素
	lru.Put("a", 1)
	lru.Put("b", 2)
	lru.Put("c", 3)
	
	// 获取元素
	val, ok := lru.Get("a")
	if !ok || val != 1 {
		t.Errorf("Expected Get('a') = 1, got %d", val)
	}
	
	// 添加新元素，应该删除最旧的
	lru.Put("d", 4)
	
	// "b"应该被删除
	if lru.Contains("b") {
		t.Error("Expected 'b' to be evicted")
	}
	
	// "a", "c", "d"应该存在
	if !lru.Contains("a") {
		t.Error("Expected 'a' to exist")
	}
	if !lru.Contains("c") {
		t.Error("Expected 'c' to exist")
	}
	if !lru.Contains("d") {
		t.Error("Expected 'd' to exist")
	}
}

func TestLRUCacheDelete(t *testing.T) {
	lru := NewLRUCache[string, int](3)
	
	lru.Put("a", 1)
	lru.Put("b", 2)
	
	if !lru.Delete("a") {
		t.Error("Expected Delete to return true")
	}
	
	if lru.Contains("a") {
		t.Error("Expected 'a' to be deleted")
	}
}

func TestLRUCacheClear(t *testing.T) {
	lru := NewLRUCache[string, int](3)
	
	lru.Put("a", 1)
	lru.Put("b", 2)
	
	lru.Clear()
	
	if lru.Size() != 0 {
		t.Error("Expected size 0 after clear")
	}
}

func TestLRUCacheSize(t *testing.T) {
	lru := NewLRUCache[string, int](3)
	
	if lru.Size() != 0 {
		t.Error("Expected initial size 0")
	}
	
	lru.Put("a", 1)
	if lru.Size() != 1 {
		t.Error("Expected size 1")
	}
	
	lru.Put("b", 2)
	if lru.Size() != 2 {
		t.Error("Expected size 2")
	}
}

func TestLRUCacheResize(t *testing.T) {
	lru := NewLRUCache[string, int](5)
	
	for i := 0; i < 5; i++ {
		lru.Put(string(rune('a'+i)), i)
	}
	
	lru.Resize(3)
	
	if lru.Size() != 3 {
		t.Errorf("Expected size 3 after resize, got %d", lru.Size())
	}
}

