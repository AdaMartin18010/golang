package main

// Pattern: Generator
// CSP Model: Generator = loop (output!value → Generator)
//
// Safety Properties:
//   - Deadlock-free: ✓ (Can be closed via context)
//   - Race-free: ✓ (Single producer)
//
// Theory: 文档16 第1.5节

import (
	"context"
)

// Generator 创建一个惰性生成器
//
// 参数:
//   - ctx: 上下文
//   - fn: 生成函数，返回(value, continue)
//
// 返回:
//   - <-chan T: 生成的值的channel
//
// 使用示例:
//   // 生成前10个自然数
//   i := 0
//   numbers := Generator(ctx, func() (int, bool) {
//       if i < 10 {
//           val := i
//           i++
//           return val, true
//       }
//       return 0, false
//   })
//
//   for num := range numbers {
//       fmt.Println(num)
//   }
func Generator[T any](ctx context.Context, fn func() (T, bool)) <-chan T {
	output := make(chan T)
	
	go func() {
		defer close(output)
		
		for {
			val, cont := fn()
			if !cont {
				return
			}
			
			select {
			case output <- val:
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return output
}

// RangeGenerator 生成一个范围内的数字
func RangeGenerator(ctx context.Context, start, end, step int) <-chan int {
	return Generator(ctx, func() (int, bool) {
		if start < end {
			val := start
			start += step
			return val, true
		}
		return 0, false
	})
}

// RepeatGenerator 重复生成相同的值
func RepeatGenerator[T any](ctx context.Context, value T, count int) <-chan T {
	i := 0
	return Generator(ctx, func() (T, bool) {
		if count < 0 || i < count {
			i++
			return value, true
		}
		var zero T
		return zero, false
	})
}

// TakeGenerator 从generator中取前n个值
func TakeGenerator[T any](ctx context.Context, input <-chan T, n int) <-chan T {
	output := make(chan T)
	
	go func() {
		defer close(output)
		
		count := 0
		for val := range input {
			if count >= n {
				return
			}
			
			select {
			case output <- val:
				count++
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return output
}
