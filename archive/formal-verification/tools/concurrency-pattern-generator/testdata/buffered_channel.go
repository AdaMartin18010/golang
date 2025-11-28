package main

// BufferedChannel 缓冲channel示例
func BufferedChannel() {
	// 创建容量为10的缓冲channel
	ch := make(chan int, 10)
	
	// 非阻塞发送（当缓冲区未满时）
	for i := 0; i < 10; i++ {
		ch <- i
	}
	
	// 接收数据
	for i := 0; i < 10; i++ {
		<-ch
	}
}
