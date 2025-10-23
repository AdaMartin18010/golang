package main

// Pattern: Fan-Out
// CSP Model: FanOut = input → (proc₁ || proc₂ || ... || proc3)
//
// Safety Properties:
//   - Deadlock-free: ✓ (All processors independent)
//   - Race-free: ✓ (Each processor gets dedicated channel)
//
// Theory: 文档16 第1.3节

import (
	"context"
)

// FanOut 将单个输入channel分发到多个处理器
//
// 参数:
//   - ctx: 上下文
//   - input: 输入channel
//   - fn: 处理函数
//   - n: 并发处理器数量
//
// 返回:
//   - <-chan Out: 输出channel
func FanOut[In any, Out any](
	ctx context.Context,
	input <-chan In,
	fn func(In) Out,
	n int,
) <-chan Out {
	outputs := make([]chan Out, n)
	
	// 创建n个processor
	for i := 0; i < n; i++ {
		outputs[i] = make(chan Out)
		go processor(ctx, input, outputs[i], fn)
	}
	
	// 合并所有输出
	return merge(ctx, outputs...)
}

// processor 处理器goroutine
func processor[In any, Out any](
	ctx context.Context,
	input <-chan In,
	output chan<- Out,
	fn func(In) Out,
) {
	defer close(output)
	
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-input:
			if !ok {
				return
			}
			result := fn(val)
			select {
			case output <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

// merge 合并多个输出channels
func merge[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	output := make(chan T)
	
	go func() {
		defer close(output)
		for _, ch := range channels {
			go func(c <-chan T) {
				for val := range c {
					select {
					case output <- val:
					case <-ctx.Done():
						return
					}
				}
			}(ch)
		}
	}()
	
	return output
}
