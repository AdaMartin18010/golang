package main

import "time"

// 示例1：Goroutine泄露 - 无限循环且未被等待
func leakExample1() {
	go func() {
		for {
			// 无限循环，没有退出条件
			time.Sleep(1 * time.Second)
		}
	}()
	// 主函数返回，但goroutine继续运行
}

// 示例2：Goroutine泄露 - Channel阻塞
func leakExample2() {
	ch := make(chan int)

	go func() {
		// 永远阻塞在接收操作
		<-ch
	}()
	// 没有发送操作，goroutine永远阻塞
}

// 示例3：正确示例 - 使用WaitGroup
func noLeakExample() {
	// 这里应该使用WaitGroup等待
	done := make(chan bool)

	go func() {
		time.Sleep(100 * time.Millisecond)
		done <- true
	}()

	<-done
}

func main() {
	leakExample1()
	leakExample2()
	noLeakExample()
}
