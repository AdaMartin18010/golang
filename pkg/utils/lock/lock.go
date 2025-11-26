package lock

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Locker 锁接口
type Locker interface {
	Lock()
	Unlock()
	TryLock() bool
	TryLockWithTimeout(timeout time.Duration) bool
}

// Mutex 互斥锁
type Mutex struct {
	mu sync.Mutex
}

// NewMutex 创建互斥锁
func NewMutex() *Mutex {
	return &Mutex{}
}

// Lock 加锁
func (m *Mutex) Lock() {
	m.mu.Lock()
}

// Unlock 解锁
func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// TryLock 尝试加锁
func (m *Mutex) TryLock() bool {
	return m.mu.TryLock()
}

// TryLockWithTimeout 带超时的尝试加锁
func (m *Mutex) TryLockWithTimeout(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	done := make(chan struct{})
	go func() {
		m.mu.Lock()
		close(done)
	}()
	
	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	}
}

// RWMutex 读写锁
type RWMutex struct {
	mu sync.RWMutex
}

// NewRWMutex 创建读写锁
func NewRWMutex() *RWMutex {
	return &RWMutex{}
}

// Lock 写锁
func (rw *RWMutex) Lock() {
	rw.mu.Lock()
}

// Unlock 写解锁
func (rw *RWMutex) Unlock() {
	rw.mu.Unlock()
}

// RLock 读锁
func (rw *RWMutex) RLock() {
	rw.mu.RLock()
}

// RUnlock 读解锁
func (rw *RWMutex) RUnlock() {
	rw.mu.RUnlock()
}

// TryLock 尝试写锁
func (rw *RWMutex) TryLock() bool {
	return rw.mu.TryLock()
}

// TryRLock 尝试读锁
func (rw *RWMutex) TryRLock() bool {
	return rw.mu.TryRLock()
}

// TryLockWithTimeout 带超时的尝试写锁
func (rw *RWMutex) TryLockWithTimeout(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	done := make(chan struct{})
	go func() {
		rw.mu.Lock()
		close(done)
	}()
	
	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	}
}

// TryRLockWithTimeout 带超时的尝试读锁
func (rw *RWMutex) TryRLockWithTimeout(timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	done := make(chan struct{})
	go func() {
		rw.mu.RLock()
		close(done)
	}()
	
	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	}
}

// KeyedMutex 键控互斥锁
type KeyedMutex struct {
	mutexes map[string]*sync.Mutex
	mu      sync.Mutex
}

// NewKeyedMutex 创建键控互斥锁
func NewKeyedMutex() *KeyedMutex {
	return &KeyedMutex{
		mutexes: make(map[string]*sync.Mutex),
	}
}

// getMutex 获取或创建指定键的互斥锁
func (km *KeyedMutex) getMutex(key string) *sync.Mutex {
	km.mu.Lock()
	defer km.mu.Unlock()
	
	if mutex, ok := km.mutexes[key]; ok {
		return mutex
	}
	
	mutex := &sync.Mutex{}
	km.mutexes[key] = mutex
	return mutex
}

// Lock 加锁
func (km *KeyedMutex) Lock(key string) {
	km.getMutex(key).Lock()
}

// Unlock 解锁
func (km *KeyedMutex) Unlock(key string) {
	km.mu.Lock()
	mutex, ok := km.mutexes[key]
	km.mu.Unlock()
	
	if ok {
		mutex.Unlock()
	}
}

// TryLock 尝试加锁
func (km *KeyedMutex) TryLock(key string) bool {
	return km.getMutex(key).TryLock()
}

// KeyedRWMutex 键控读写锁
type KeyedRWMutex struct {
	mutexes map[string]*sync.RWMutex
	mu      sync.Mutex
}

// NewKeyedRWMutex 创建键控读写锁
func NewKeyedRWMutex() *KeyedRWMutex {
	return &KeyedRWMutex{
		mutexes: make(map[string]*sync.RWMutex),
	}
}

// getMutex 获取或创建指定键的读写锁
func (km *KeyedRWMutex) getMutex(key string) *sync.RWMutex {
	km.mu.Lock()
	defer km.mu.Unlock()
	
	if mutex, ok := km.mutexes[key]; ok {
		return mutex
	}
	
	mutex := &sync.RWMutex{}
	km.mutexes[key] = mutex
	return mutex
}

// Lock 写锁
func (km *KeyedRWMutex) Lock(key string) {
	km.getMutex(key).Lock()
}

// Unlock 写解锁
func (km *KeyedRWMutex) Unlock(key string) {
	km.mu.Lock()
	mutex, ok := km.mutexes[key]
	km.mu.Unlock()
	
	if ok {
		mutex.Unlock()
	}
}

// RLock 读锁
func (km *KeyedRWMutex) RLock(key string) {
	km.getMutex(key).RLock()
}

// RUnlock 读解锁
func (km *KeyedRWMutex) RUnlock(key string) {
	km.mu.Lock()
	mutex, ok := km.mutexes[key]
	km.mu.Unlock()
	
	if ok {
		mutex.RUnlock()
	}
}

// TryLock 尝试写锁
func (km *KeyedRWMutex) TryLock(key string) bool {
	return km.getMutex(key).TryLock()
}

// TryRLock 尝试读锁
func (km *KeyedRWMutex) TryRLock(key string) bool {
	return km.getMutex(key).TryRLock()
}

// SpinLock 自旋锁
type SpinLock struct {
	locked int32
}

// NewSpinLock 创建自旋锁
func NewSpinLock() *SpinLock {
	return &SpinLock{}
}

// Lock 加锁
func (sl *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&sl.locked, 0, 1) {
		// 自旋等待
		runtime.Gosched()
	}
}

// Unlock 解锁
func (sl *SpinLock) Unlock() {
	atomic.StoreInt32(&sl.locked, 0)
}

// TryLock 尝试加锁
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&sl.locked, 0, 1)
}

// TryLockWithTimeout 带超时的尝试加锁
func (sl *SpinLock) TryLockWithTimeout(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if atomic.CompareAndSwapInt32(&sl.locked, 0, 1) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

