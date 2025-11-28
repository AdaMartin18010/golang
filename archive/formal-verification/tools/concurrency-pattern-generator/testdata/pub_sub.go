package main

import "sync"

// PubSub 发布订阅
type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan interface{}
}

// NewPubSub 创建PubSub
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Subscribe 订阅主题
func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	ch := make(chan interface{}, 10)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

// Publish 发布消息
func (ps *PubSub) Publish(topic string, msg interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	
	for _, ch := range ps.subscribers[topic] {
		select {
		case ch <- msg:
		default:
		}
	}
}
