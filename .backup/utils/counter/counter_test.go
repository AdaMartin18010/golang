package counter

import (
	"testing"
	"time"
)

func TestSimpleCounter(t *testing.T) {
	c := NewSimpleCounter()
	
	// 增加
	if c.Increment() != 1 {
		t.Error("Expected 1")
	}
	if c.Increment() != 2 {
		t.Error("Expected 2")
	}
	
	// 减少
	if c.Decrement() != 1 {
		t.Error("Expected 1")
	}
	
	// 增加指定值
	if c.Add(5) != 6 {
		t.Error("Expected 6")
	}
	
	// 获取值
	if c.Get() != 6 {
		t.Error("Expected 6")
	}
	
	// 重置
	c.Reset()
	if c.Get() != 0 {
		t.Error("Expected 0")
	}
}

func TestMaxCounter(t *testing.T) {
	c := NewMaxCounter()
	
	c.Increment()
	c.Increment()
	
	if c.Get() != 2 {
		t.Error("Expected 2")
	}
	
	// 尝试减少（应该无效）
	c.Add(-1)
	if c.Get() != 2 {
		t.Error("Expected 2")
	}
}

func TestMinCounter(t *testing.T) {
	c := NewMinCounter(10)
	
	c.Decrement()
	if c.Get() != 9 {
		t.Error("Expected 9")
	}
	
	c.Subtract(5)
	if c.Get() != 4 {
		t.Error("Expected 4")
	}
}

func TestRateCounter(t *testing.T) {
	rc := NewRateCounter(1*time.Second, 100*time.Millisecond)
	
	rc.Increment()
	rc.Increment()
	
	rate := rc.Get()
	if rate <= 0 {
		t.Error("Expected positive rate")
	}
}

func TestSlidingWindowCounter(t *testing.T) {
	swc := NewSlidingWindowCounter(1*time.Second, 100*time.Millisecond)
	
	swc.Increment()
	swc.Increment()
	
	total := swc.Get()
	if total != 2 {
		t.Errorf("Expected 2, got %d", total)
	}
}

func TestMultiCounter(t *testing.T) {
	mc := NewMultiCounter()
	
	mc.Increment("key1")
	mc.Increment("key1")
	mc.Increment("key2")
	
	if mc.Get("key1") != 2 {
		t.Error("Expected 2")
	}
	if mc.Get("key2") != 1 {
		t.Error("Expected 1")
	}
	
	all := mc.GetAll()
	if len(all) != 2 {
		t.Error("Expected 2 keys")
	}
}

