package main

type Semaphore struct {
	slots chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{slots: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire() {
	s.slots <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.slots
}
