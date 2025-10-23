package main

// Pattern: Pipeline
// CSP Model: Pipeline = stage₁ >> stage₂ >> ... >> stageₙ
//
// Safety Properties:
//   - Deadlock-free: ✓ (Forward progress guaranteed)
//   - Race-free: ✓ (Sequential stages)
//
// Theory: 文档16 第1.4节

import (
	"context"
)

// Stage 表示管道的一个阶段
type Stage[In any, Out any] func(context.Context, <-chan In) <-chan Out

// Pipeline 创建一个多阶段处理管道
//
// 使用示例:
//   // 定义stages
//   stage1 := func(ctx context.Context, in <-chan int) <-chan int {
//       out := make(chan int)
//       go func() {
//           defer close(out)
//           for val := range in {
//               out <- val * 2
//           }
//       }()
//       return out
//   }
//
//   stage2 := func(ctx context.Context, in <-chan int) <-chan int {
//       out := make(chan int)
//       go func() {
//           defer close(out)
//           for val := range in {
//               out <- val + 1
//           }
//       }()
//       return out
//   }
//
//   // 创建pipeline
//   input := make(chan int)
//   output := Pipeline(context.Background(), input, stage1, stage2)
func Pipeline[T any](
	ctx context.Context,
	input <-chan T,
	stages ...Stage[T, T],
) <-chan T {
	current := input
	
	for _, stage := range stages {
		current = stage(ctx, current)
	}
	
	return current
}

// MapStage 创建一个map阶段
func MapStage[In any, Out any](fn func(In) Out) Stage[In, Out] {
	return func(ctx context.Context, input <-chan In) <-chan Out {
		output := make(chan Out)
		
		go func() {
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
		}()
		
		return output
	}
}

// FilterStage 创建一个filter阶段
func FilterStage[T any](predicate func(T) bool) Stage[T, T] {
	return func(ctx context.Context, input <-chan T) <-chan T {
		output := make(chan T)
		
		go func() {
			defer close(output)
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-input:
					if !ok {
						return
					}
					if predicate(val) {
						select {
						case output <- val:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
		
		return output
	}
}
