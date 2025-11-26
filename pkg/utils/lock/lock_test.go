package lock

import (
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	m := NewMutex()
	
	// 加锁
	m.Lock()
	
	// 尝试加锁（应该失败）
	if m.TryLock() {
		t.Error("Expected TryLock to fail")
	}
	
	// 解锁
	m.Unlock()
	
	// 尝试加锁（应该成功）
	if !m.TryLock() {
		t.Error("Expected TryLock to succeed")
	}
	m.Unlock()
}

func TestRWMutex(t *testing.T) {
	rw := NewRWMutex()
	
	// 读锁
	rw.RLock()
	
	// 尝试读锁（应该成功）
	if !rw.TryRLock() {
		t.Error("Expected TryRLock to succeed")
	}
	rw.RUnlock()
	rw.RUnlock()
	
	// 写锁
	rw.Lock()
	
	// 尝试读锁（应该失败）
	if rw.TryRLock() {
		t.Error("Expected TryRLock to fail")
	}
	
	rw.Unlock()
}

func TestKeyedMutex(t *testing.T) {
	km := NewKeyedMutex()
	
	// 不同键可以同时加锁
	km.Lock("key1")
	km.Lock("key2")
	
	// 相同键不能同时加锁
	if km.TryLock("key1") {
		t.Error("Expected TryLock to fail for same key")
	}
	
	km.Unlock("key1")
	km.Unlock("key2")
}

func TestKeyedRWMutex(t *testing.T) {
	km := NewKeyedRWMutex()
	
	// 读锁
	km.RLock("key1")
	
	// 尝试读锁（应该成功）
	if !km.TryRLock("key1") {
		t.Error("Expected TryRLock to succeed")
	}
	km.RUnlock("key1")
	km.RUnlock("key1")
	
	// 写锁
	km.Lock("key1")
	
	// 尝试读锁（应该失败）
	if km.TryRLock("key1") {
		t.Error("Expected TryRLock to fail")
	}
	
	km.Unlock("key1")
}

func TestSpinLock(t *testing.T) {
	sl := NewSpinLock()
	
	// 加锁
	sl.Lock()
	
	// 尝试加锁（应该失败）
	if sl.TryLock() {
		t.Error("Expected TryLock to fail")
	}
	
	// 解锁
	sl.Unlock()
	
	// 尝试加锁（应该成功）
	if !sl.TryLock() {
		t.Error("Expected TryLock to succeed")
	}
	sl.Unlock()
}

func TestMutexConcurrency(t *testing.T) {
	m := NewMutex()
	var counter int64
	
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Lock()
			counter++
			m.Unlock()
		}()
	}
	
	wg.Wait()
	
	if counter != 100 {
		t.Errorf("Expected counter 100, got %d", counter)
	}
}

