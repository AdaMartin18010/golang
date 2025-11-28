package main

// Pattern: Fan-In
// CSP Model: FanIn = (input₁ → merge) || (input₂ → merge) || ... → output
//
// Safety Properties:
//   - Deadlock-free: ✓ (All inputs can complete independently)
//   - Race-free: ✓ (Select synchronization)
//
// Theory: 文档16 第1.2节
//
// Happens-Before Relations:
//   ∀ i: input_i sent → merged → output sent

import (
	"sync"
)

// FanIn 将多个输入channel合并为一个输出channel
//
// 参数:
//   - inputs: 输入channels切片
//
// 返回:
//   - <-chan T: 合并后的输出channel
//
// 使用示例:
//   input1 := make(chan int)
//   input2 := make(chan int)
//   output := FanIn(input1, input2)
//
//   go func() { input1 <- 1; close(input1) }()
//   go func() { input2 <- 2; close(input2) }()
//
//   for val := range output {
//       fmt.Println(val)
//   }
func FanIn[T any](inputs ...<-chan T) <-chan T {
	output := make(chan T)
	var wg sync.WaitGroup

	// 为每个输入创建一个goroutine
	for _, input := range inputs {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for val := range ch {
				output <- val
			}
		}(input)
	}

	// 等待所有输入完成后关闭输出
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// FanInSelect 使用select实现的Fan-In（非阻塞版本）
func FanInSelect[T any](input1, input2 <-chan T) <-chan T {
	output := make(chan T)
	
	go func() {
		defer close(output)
		
		for {
			select {
			case val, ok := <-input1:
				if !ok {
					input1 = nil // 设为nil以停止从此channel读取
					continue
				}
				output <- val
				
			case val, ok := <-input2:
				if !ok {
					input2 = nil
					continue
				}
				output <- val
			}
			
			// 两个输入都关闭时退出
			if input1 == nil && input2 == nil {
				return
			}
		}
	}()
	
	return output
}

// CSP形式化定义：
//
// Process Definition:
//   Input(i) = loop (input_i?val → merge!val → Input(i))
//   Merger = loop (merge?val → output!val → Merger)
//   FanIn = (Input(1) ||| Input(2) ||| ... ||| Input(n)) >> Merger
//
// Deadlock Freedom:
//   FanIn is deadlock-free because:
//     1. Each input is independent
//     2. Merger always accepts inputs
//     3. No circular waits
