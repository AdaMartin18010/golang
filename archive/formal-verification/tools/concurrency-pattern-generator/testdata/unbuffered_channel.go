package main

// UnbufferedChannel 无缓冲channel示例
func UnbufferedChannel() {
	ch := make(chan int)
	
	go func() {
		ch <- 42 // 阻塞直到接收
	}()
	
	val := <-ch // 阻塞直到发送
	_ = val
}
