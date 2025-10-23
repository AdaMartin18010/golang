package main

import "context"

// Producer 生产者
func Producer(ctx context.Context, n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < n; i++ {
			select {
			case ch <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

// Consumer 消费者
func Consumer(ctx context.Context, ch <-chan int) {
	for {
		select {
		case val, ok := <-ch:
			if !ok {
				return
			}
			// Process val
			_ = val
		case <-ctx.Done():
			return
		}
	}
}
