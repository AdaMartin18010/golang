package main

import "sync"

type BoundedQueue struct {
	mu sync.Mutex
	cond *sync.Cond
	items []interface{}
	max int
}

func NewBoundedQueue(n int) *BoundedQueue {
	q := &BoundedQueue{max: n}
	q.cond = sync.NewCond(&q.mu)
	return q
}
