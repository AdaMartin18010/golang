package main

import "time"

// SelectPattern select多路复用
func SelectPattern(ch1, ch2 <-chan int, done <-chan struct{}) {
	for {
		select {
		case v1 := <-ch1:
			// Handle ch1
			_ = v1
		case v2 := <-ch2:
			// Handle ch2
			_ = v2
		case <-done:
			return
		case <-time.After(time.Second):
			// Timeout
			return
		}
	}
}
