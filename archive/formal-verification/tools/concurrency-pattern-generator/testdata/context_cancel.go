package main

import (
	"context"
	"time"
)

// WithCancel 创建可取消的context
func WithCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Do work
			}
		}
	}()
	
	time.Sleep(time.Second)
	cancel() // 取消所有goroutine
}
