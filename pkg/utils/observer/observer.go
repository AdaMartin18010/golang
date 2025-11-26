package observer

import (
	"sync"
)

// Observer 观察者接口
type Observer[T any] interface {
	Update(data T)
}

// ObserverFunc 观察者函数类型
type ObserverFunc[T any] func(data T)

// Update 实现Observer接口
func (f ObserverFunc[T]) Update(data T) {
	f(data)
}

// Subject 主题接口
type Subject[T any] interface {
	Subscribe(observer Observer[T]) func() // 返回取消订阅函数
	Unsubscribe(observer Observer[T])
	Notify(data T)
	Count() int
	Clear()
}

// SimpleSubject 简单主题实现
type SimpleSubject[T any] struct {
	observers []Observer[T]
	mu        sync.RWMutex
}

// NewSimpleSubject 创建简单主题
func NewSimpleSubject[T any]() *SimpleSubject[T] {
	return &SimpleSubject[T]{
		observers: make([]Observer[T], 0),
	}
}

// Subscribe 订阅
func (s *SimpleSubject[T]) Subscribe(observer Observer[T]) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.observers = append(s.observers, observer)
	
	// 返回取消订阅函数
	return func() {
		s.Unsubscribe(observer)
	}
}

// Unsubscribe 取消订阅
func (s *SimpleSubject[T]) Unsubscribe(observer Observer[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify 通知所有观察者
func (s *SimpleSubject[T]) Notify(data T) {
	s.mu.RLock()
	observers := make([]Observer[T], len(s.observers))
	copy(observers, s.observers)
	s.mu.RUnlock()
	
	for _, observer := range observers {
		observer.Update(data)
	}
}

// Count 获取观察者数量
func (s *SimpleSubject[T]) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.observers)
}

// Clear 清空所有观察者
func (s *SimpleSubject[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = s.observers[:0]
}

// AsyncSubject 异步主题实现
type AsyncSubject[T any] struct {
	observers []Observer[T]
	mu        sync.RWMutex
}

// NewAsyncSubject 创建异步主题
func NewAsyncSubject[T any]() *AsyncSubject[T] {
	return &AsyncSubject[T]{
		observers: make([]Observer[T], 0),
	}
}

// Subscribe 订阅
func (s *AsyncSubject[T]) Subscribe(observer Observer[T]) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.observers = append(s.observers, observer)
	
	return func() {
		s.Unsubscribe(observer)
	}
}

// Unsubscribe 取消订阅
func (s *AsyncSubject[T]) Unsubscribe(observer Observer[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify 异步通知所有观察者
func (s *AsyncSubject[T]) Notify(data T) {
	s.mu.RLock()
	observers := make([]Observer[T], len(s.observers))
	copy(observers, s.observers)
	s.mu.RUnlock()
	
	for _, observer := range observers {
		go observer.Update(data)
	}
}

// Count 获取观察者数量
func (s *AsyncSubject[T]) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.observers)
}

// Clear 清空所有观察者
func (s *AsyncSubject[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = s.observers[:0]
}

// FilteredSubject 过滤主题实现
type FilteredSubject[T any] struct {
	observers []filteredObserver[T]
	mu        sync.RWMutex
}

type filteredObserver[T any] struct {
	observer Observer[T]
	filter   func(T) bool
}

// NewFilteredSubject 创建过滤主题
func NewFilteredSubject[T any]() *FilteredSubject[T] {
	return &FilteredSubject[T]{
		observers: make([]filteredObserver[T], 0),
	}
}

// Subscribe 订阅（带过滤条件）
func (s *FilteredSubject[T]) Subscribe(observer Observer[T], filter func(T) bool) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.observers = append(s.observers, filteredObserver[T]{
		observer: observer,
		filter:   filter,
	})
	
	return func() {
		s.Unsubscribe(observer)
	}
}

// Unsubscribe 取消订阅
func (s *FilteredSubject[T]) Unsubscribe(observer Observer[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, fo := range s.observers {
		if fo.observer == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify 通知所有观察者（根据过滤条件）
func (s *FilteredSubject[T]) Notify(data T) {
	s.mu.RLock()
	observers := make([]filteredObserver[T], len(s.observers))
	copy(observers, s.observers)
	s.mu.RUnlock()
	
	for _, fo := range observers {
		if fo.filter == nil || fo.filter(data) {
			fo.observer.Update(data)
		}
	}
}

// Count 获取观察者数量
func (s *FilteredSubject[T]) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.observers)
}

// Clear 清空所有观察者
func (s *FilteredSubject[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = s.observers[:0]
}

// EventBus 事件总线
type EventBus struct {
	subjects map[string]Subject[interface{}]
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		subjects: make(map[string]Subject[interface{}]),
	}
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(event string, observer Observer[interface{}]) func() {
	eb.mu.Lock()
	subject, ok := eb.subjects[event]
	if !ok {
		subject = NewSimpleSubject[interface{}]()
		eb.subjects[event] = subject
	}
	eb.mu.Unlock()
	
	return subject.Subscribe(observer)
}

// Publish 发布事件
func (eb *EventBus) Publish(event string, data interface{}) {
	eb.mu.RLock()
	subject, ok := eb.subjects[event]
	eb.mu.RUnlock()
	
	if ok {
		subject.Notify(data)
	}
}

// Unsubscribe 取消订阅
func (eb *EventBus) Unsubscribe(event string, observer Observer[interface{}]) {
	eb.mu.RLock()
	subject, ok := eb.subjects[event]
	eb.mu.RUnlock()
	
	if ok {
		subject.Unsubscribe(observer)
	}
}

// Clear 清空指定事件的所有观察者
func (eb *EventBus) Clear(event string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	if subject, ok := eb.subjects[event]; ok {
		subject.Clear()
	}
}

// ClearAll 清空所有事件
func (eb *EventBus) ClearAll() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	eb.subjects = make(map[string]Subject[interface{}])
}

// Count 获取指定事件的观察者数量
func (eb *EventBus) Count(event string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	
	if subject, ok := eb.subjects[event]; ok {
		return subject.Count()
	}
	return 0
}

