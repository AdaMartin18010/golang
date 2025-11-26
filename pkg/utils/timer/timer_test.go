package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSimpleTimer(t *testing.T) {
	var count int64
	timer := NewSimpleTimer(100*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	
	timer.Start()
	time.Sleep(250 * time.Millisecond)
	timer.Stop()
	
	if atomic.LoadInt64(&count) < 2 {
		t.Error("Expected at least 2 executions")
	}
	
	if timer.IsRunning() {
		t.Error("Expected timer to be stopped")
	}
}

func TestOneShotTimer(t *testing.T) {
	var executed bool
	timer := NewOneShotTimer(100*time.Millisecond, func() {
		executed = true
	})
	
	timer.Start()
	time.Sleep(150 * time.Millisecond)
	
	if !executed {
		t.Error("Expected timer to execute")
	}
	
	if timer.IsRunning() {
		t.Error("Expected timer to be stopped after execution")
	}
}

func TestDebounceTimer(t *testing.T) {
	var count int64
	timer := NewDebounceTimer(100*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	
	// 快速触发多次
	for i := 0; i < 10; i++ {
		timer.Trigger()
		time.Sleep(10 * time.Millisecond)
	}
	
	time.Sleep(150 * time.Millisecond)
	
	// 应该只执行一次
	if atomic.LoadInt64(&count) != 1 {
		t.Errorf("Expected 1 execution, got %d", atomic.LoadInt64(&count))
	}
}

func TestThrottleTimer(t *testing.T) {
	var count int64
	timer := NewThrottleTimer(100*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	
	// 快速触发多次
	for i := 0; i < 10; i++ {
		timer.Trigger()
		time.Sleep(10 * time.Millisecond)
	}
	
	// 应该只执行一次（在第一次触发时）
	if atomic.LoadInt64(&count) != 1 {
		t.Errorf("Expected 1 execution, got %d", atomic.LoadInt64(&count))
	}
}

func TestIntervalTimer(t *testing.T) {
	var count int64
	timer := NewIntervalTimer(100*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	
	timer.Start()
	time.Sleep(250 * time.Millisecond)
	timer.Stop()
	
	if atomic.LoadInt64(&count) < 2 {
		t.Error("Expected at least 2 executions")
	}
	
	if timer.ExecutionCount() < 2 {
		t.Error("Expected execution count >= 2")
	}
}

