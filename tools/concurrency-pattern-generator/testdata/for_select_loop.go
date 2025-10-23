package main

import "time"

// ForSelectLoop for-select事件循环
func ForSelectLoop(input <-chan int, done <-chan struct{}) {
	for {
		select {
		case val, ok := <-input:
			if !ok {
				return
			}
			// Process val
			_ = val
		case <-done:
			return
		}
	}
}
