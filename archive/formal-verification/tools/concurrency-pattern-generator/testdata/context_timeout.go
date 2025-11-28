package main

import (
	"context"
	"time"
)

// WithTimeout 创建带超时的context
func WithTimeout() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(3 * time.Second):
		return nil
	}
}

// WithDeadline 创建带截止时间的context
func WithDeadline() error {
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
