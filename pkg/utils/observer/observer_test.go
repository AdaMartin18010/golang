package observer

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSimpleSubject(t *testing.T) {
	subject := NewSimpleSubject[string]()
	
	var count int64
	observer := ObserverFunc[string](func(data string) {
		atomic.AddInt64(&count, 1)
	})
	
	// 订阅
	unsubscribe := subject.Subscribe(observer)
	
	// 通知
	subject.Notify("test")
	
	if atomic.LoadInt64(&count) != 1 {
		t.Error("Expected 1 notification")
	}
	
	// 取消订阅
	unsubscribe()
	
	// 再次通知
	subject.Notify("test")
	
	if atomic.LoadInt64(&count) != 1 {
		t.Error("Expected still 1 notification")
	}
}

func TestAsyncSubject(t *testing.T) {
	subject := NewAsyncSubject[string]()
	
	var count int64
	observer := ObserverFunc[string](func(data string) {
		atomic.AddInt64(&count, 1)
	})
	
	subject.Subscribe(observer)
	subject.Notify("test")
	
	// 等待异步执行
	time.Sleep(100 * time.Millisecond)
	
	if atomic.LoadInt64(&count) != 1 {
		t.Error("Expected 1 notification")
	}
}

func TestFilteredSubject(t *testing.T) {
	subject := NewFilteredSubject[int]()
	
	var count int64
	observer := ObserverFunc[int](func(data int) {
		atomic.AddInt64(&count, 1)
	})
	
	// 只接收大于5的值
	subject.Subscribe(observer, func(data int) bool {
		return data > 5
	})
	
	subject.Notify(3)  // 不会触发
	subject.Notify(10) // 会触发
	
	if atomic.LoadInt64(&count) != 1 {
		t.Error("Expected 1 notification")
	}
}

func TestEventBus(t *testing.T) {
	bus := NewEventBus()
	
	var count int64
	observer := ObserverFunc[interface{}](func(data interface{}) {
		atomic.AddInt64(&count, 1)
	})
	
	bus.Subscribe("event1", observer)
	bus.Publish("event1", "data")
	
	if atomic.LoadInt64(&count) != 1 {
		t.Error("Expected 1 notification")
	}
}

