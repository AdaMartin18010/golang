package main

import "context"

// Message Actor消息
type Message struct {
	Type string
	Data interface{}
}

// Actor Actor模型
type Actor struct {
	mailbox chan Message
}

// NewActor 创建Actor
func NewActor() *Actor {
	return &Actor{
		mailbox: make(chan Message, 10),
	}
}

// Send 发送消息
func (a *Actor) Send(msg Message) {
	a.mailbox <- msg
}

// Start 启动Actor
func (a *Actor) Start(ctx context.Context) {
	for {
		select {
		case msg := <-a.mailbox:
			// Handle message
			_ = msg
		case <-ctx.Done():
			return
		}
	}
}
