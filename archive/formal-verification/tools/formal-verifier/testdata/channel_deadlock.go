package main

// 示例1：无缓冲Channel死锁 - 发送但无接收
func deadlockExample1() {
	ch := make(chan int)

	// 发送操作会永远阻塞
	ch <- 42
	// 没有接收者
}

// 示例2：有缓冲Channel死锁 - 超过缓冲容量
func deadlockExample2() {
	ch := make(chan int, 2)

	// 前两次发送成功
	ch <- 1
	ch <- 2
	// 第三次发送会阻塞
	ch <- 3
	// 没有接收者
}

// 示例3：循环依赖死锁
func deadlockExample3() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		// 等待ch1，然后发送到ch2
		<-ch1
		ch2 <- 1
	}()

	go func() {
		// 等待ch2，然后发送到ch1
		<-ch2
		ch1 <- 1
	}()

	// 两个goroutine互相等待，形成死锁
}

// 示例4：正确示例 - 发送和接收平衡
func noDeadlockExample() {
	ch := make(chan int, 1)

	go func() {
		ch <- 42
	}()

	<-ch
}

func main() {
	// deadlockExample1() // 会死锁
	// deadlockExample2() // 会死锁
	// deadlockExample3() // 会死锁
	noDeadlockExample()
}
