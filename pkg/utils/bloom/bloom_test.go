package bloom

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter(1000, 3)
	
	// 添加元素
	bf.AddString("test1")
	bf.AddString("test2")
	bf.AddString("test3")
	
	// 检查存在的元素
	if !bf.ContainsString("test1") {
		t.Error("Expected 'test1' to be in filter")
	}
	if !bf.ContainsString("test2") {
		t.Error("Expected 'test2' to be in filter")
	}
	if !bf.ContainsString("test3") {
		t.Error("Expected 'test3' to be in filter")
	}
	
	// 检查不存在的元素（可能误判，但应该大部分情况下正确）
	if bf.ContainsString("nonexistent") {
		// 这是可能的假阳性，但概率应该很低
		t.Log("False positive detected (this is possible)")
	}
}

func TestBloomFilterBytes(t *testing.T) {
	bf := NewBloomFilter(1000, 3)
	
	// 添加字节数组
	bf.Add([]byte("test1"))
	bf.Add([]byte("test2"))
	
	// 检查
	if !bf.Contains([]byte("test1")) {
		t.Error("Expected 'test1' to be in filter")
	}
	if !bf.Contains([]byte("test2")) {
		t.Error("Expected 'test2' to be in filter")
	}
}

func TestBloomFilterClear(t *testing.T) {
	bf := NewBloomFilter(1000, 3)
	
	bf.AddString("test1")
	if !bf.ContainsString("test1") {
		t.Error("Expected 'test1' to be in filter")
	}
	
	bf.Clear()
	if bf.ContainsString("test1") {
		t.Error("Expected 'test1' not to be in filter after clear")
	}
}

func TestOptimalSize(t *testing.T) {
	size := OptimalSize(1000, 0.01)
	if size == 0 {
		t.Error("Expected non-zero size")
	}
}

func TestOptimalHashCount(t *testing.T) {
	count := OptimalHashCount(1000, 10000)
	if count == 0 {
		t.Error("Expected non-zero hash count")
	}
}

