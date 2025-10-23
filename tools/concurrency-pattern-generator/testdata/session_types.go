package main

// Session 会话类型
type Session struct {
	send    chan<- interface{}
	receive <-chan interface{}
}

// NewSession 创建会话
func NewSession() (*Session, *Session) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	
	s1 := &Session{send: ch1, receive: ch2}
	s2 := &Session{send: ch2, receive: ch1}
	
	return s1, s2
}

// Send 发送数据
func (s *Session) Send(data interface{}) {
	s.send <- data
}

// Receive 接收数据
func (s *Session) Receive() interface{} {
	return <-s.receive
}
