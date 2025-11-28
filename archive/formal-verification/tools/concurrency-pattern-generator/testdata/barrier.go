package main

import "sync"

type Barrier struct {
	count int
	target int
	mu sync.Mutex
	cond *sync.Cond
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{target: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}
