package main

import "sync"

type CountDownLatch struct {
	count int
	mu sync.Mutex
	cond *sync.Cond
}

func NewCountDownLatch(n int) *CountDownLatch {
	l := &CountDownLatch{count: n}
	l.cond = sync.NewCond(&l.mu)
	return l
}
