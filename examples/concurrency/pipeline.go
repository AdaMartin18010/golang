package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Stage 表示管道中的一个阶段
type Stage func(<-chan int) <-chan int

// 生成器：生成数字序列
func generator(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			out <- i
		}
	}()
	return out
}

// 平方阶段：计算平方
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			result := n * n
			fmt.Printf("Square: %d^2 = %d\n", n, result)
			out <- result
		}
	}()
	return out
}

// 过滤阶段：过滤偶数
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				fmt.Printf("Filter: %d is even, passing through\n", n)
				out <- n
			} else {
				fmt.Printf("Filter: %d is odd, filtering out\n", n)
			}
		}
	}()
	return out
}

// 打印阶段：打印结果
func printer(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			fmt.Printf("Printer: %d\n", n)
			out <- n
		}
	}()
	return out
}

// 合并多个通道
func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// 为每个输入通道启动一个goroutine
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(ch)
	}

	// 启动一个goroutine来关闭输出通道
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// 扇出：将一个通道分发到多个通道
func fanOut(in <-chan int, numWorkers int) []<-chan int {
	channels := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		ch := make(chan int)
		channels[i] = ch

		go func(workerID int, out chan<- int) {
			defer close(out)
			for n := range in {
				// 模拟处理时间
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				fmt.Printf("Worker %d processing: %d\n", workerID, n)
				out <- n * workerID // 简单的处理逻辑
			}
		}(i+1, ch)
	}

	return channels
}

// 带缓冲的管道阶段
func bufferedStage(in <-chan int, bufferSize int) <-chan int {
	out := make(chan int, bufferSize)
	go func() {
		defer close(out)
		for n := range in {
			out <- n
		}
	}()
	return out
}

// 错误处理管道阶段
func errorHandlingStage(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n < 0 {
				fmt.Printf("Error: negative number %d, skipping\n", n)
				continue
			}
			out <- n
		}
	}()
	return out
}

// 示例1：基础管道
func basicPipeline() {
	fmt.Println("=== Basic Pipeline Example ===")

	// 创建管道：生成 -> 平方 -> 过滤 -> 打印
	pipeline := func(in <-chan int) <-chan int {
		return printer(filterEven(square(in)))
	}

	// 运行管道
	input := generator(10)
	output := pipeline(input)

	// 收集结果
	var results []int
	for n := range output {
		results = append(results, n)
	}

	fmt.Printf("Pipeline results: %v\n", results)
	fmt.Println()
}

// 示例2：扇入扇出管道
func fanInFanOutPipeline() {
	fmt.Println("=== Fan-In Fan-Out Pipeline Example ===")

	// 生成输入
	input := generator(20)

	// 扇出到3个工作协程
	workerOutputs := fanOut(input, 3)

	// 扇入合并结果
	merged := merge(workerOutputs...)

	// 收集结果
	var results []int
	for n := range merged {
		results = append(results, n)
	}

	fmt.Printf("Fan-In Fan-Out results: %v\n", results)
	fmt.Println()
}

// 示例3：带缓冲和错误处理的管道
func advancedPipeline() {
	fmt.Println("=== Advanced Pipeline Example ===")

	// 创建高级管道
	pipeline := func(in <-chan int) <-chan int {
		return printer(
			errorHandlingStage(
				bufferedStage(
					filterEven(
						square(in)), 5)))
	}

	// 运行管道
	input := generator(15)
	output := pipeline(input)

	// 收集结果
	var results []int
	for n := range output {
		results = append(results, n)
	}

	fmt.Printf("Advanced pipeline results: %v\n", results)
	fmt.Println()
}

// 示例4：动态管道构建
func dynamicPipeline() {
	fmt.Println("=== Dynamic Pipeline Example ===")

	// 动态构建管道阶段
	stages := []Stage{
		square,
		filterEven,
		func(in <-chan int) <-chan int {
			return bufferedStage(in, 3)
		},
		printer,
	}

	// 构建管道
	var pipeline func(<-chan int) <-chan int
	pipeline = func(in <-chan int) <-chan int {
		for _, stage := range stages {
			in = stage(in)
		}
		return in
	}

	// 运行管道
	input := generator(8)
	output := pipeline(input)

	// 收集结果
	var results []int
	for n := range output {
		results = append(results, n)
	}

	fmt.Printf("Dynamic pipeline results: %v\n", results)
	fmt.Println()
}

func main() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Go Pipeline Patterns Examples")
	fmt.Println("==============================")

	// 运行各种管道示例
	basicPipeline()
	fanInFanOutPipeline()
	advancedPipeline()
	dynamicPipeline()

	fmt.Println("All pipeline examples completed!")
}
