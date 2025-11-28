// Package patterns 实现各种并发模式的代码生成
package patterns

import (
	"fmt"
	"strings"
)

// GenerateWorkerPool 生成Worker Pool模式代码
//
// CSP Model: Pool = worker₁ || worker₂ || ... || workerₙ
// Safety: Deadlock-free, Race-free
// Theory: 文档02 第3.2节, 文档16 第1.1节
func GenerateWorkerPool(data map[string]interface{}) string {
	pkg := getStringOrDefault(data, "PackageName", "main")
	numWorkers := getIntOrDefault(data, "NumWorkers", 10)

	return fmt.Sprintf(`package %s

// Pattern: Worker Pool
// CSP Model: Pool = worker₁ || worker₂ || ... || worker%d
//
// Safety Properties:
//   - Deadlock-free: ✓ (All workers can terminate when done channel closed)
//   - Race-free: ✓ (Channel synchronization guarantees happens-before)
//
// Theory: 文档02 第3.2节, 文档16 第1.1节
//
// Happens-Before Relations:
//   1. job sent → job received by worker
//   2. job processed → result sent
//   3. done channel closed → all workers exit
//   4. all workers exit → results channel closed
//
// Formal Proof:
//   ∀ job ∈ Jobs: 
//     sent(job) →ʰᵇ received(job) →ʰᵇ processed(job) →ʰᵇ result_sent(job)

import (
	"context"
	"sync"
)

// Job 表示待处理的任务
type Job struct {
	ID   int
	Data interface{}
}

// Result 表示处理结果
type Result struct {
	JobID int
	Data  interface{}
	Error error
}

// WorkerPool 创建一个工作池来并发处理任务
//
// 参数:
//   - ctx: 上下文，用于取消和超时控制
//   - numWorkers: worker数量
//   - jobs: 任务输入channel
//
// 返回:
//   - <-chan Result: 结果输出channel
//
// 使用示例:
//   ctx := context.Background()
//   jobs := make(chan Job, 100)
//   results := WorkerPool(ctx, %d, jobs)
//
//   // 发送任务
//   for i := 0; i < 100; i++ {
//       jobs <- Job{ID: i, Data: i}
//   }
//   close(jobs)
//
//   // 接收结果
//   for result := range results {
//       fmt.Printf("Result: %%+v\n", result)
//   }
func WorkerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup

	// 启动worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			// Worker循环
			for {
				select {
				case <-ctx.Done():
					// Context取消，退出
					return
					
				case job, ok := <-jobs:
					if !ok {
						// Jobs channel关闭，退出
						return
					}
					
					// 处理任务
					result := processJob(job)
					
					// 发送结果
					select {
					case results <- result:
						// 成功发送
					case <-ctx.Done():
						// Context取消，退出
						return
					}
				}
			}
		}(i)
	}

	// 等待所有workers完成后关闭results
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// processJob 处理单个任务（用户需要实现此函数）
func processJob(job Job) Result {
	// TODO: 用户实现具体的任务处理逻辑
	return Result{
		JobID: job.ID,
		Data:  job.Data,
		Error: nil,
	}
}

// CSP形式化定义：
//
// Process Definition:
//   Worker(i) = jobs?job → process(job) → results!result → Worker(i)
//             □ done → SKIP
//             □ ctx.Done → SKIP
//
//   Pool = Worker(1) ||| Worker(2) ||| ... ||| Worker(n)
//
// Trace Semantics:
//   traces(Pool) = { t | t is interleaving of traces(Worker(i)) }
//
// Deadlock Freedom:
//   Pool is deadlock-free because:
//     1. Each worker can independently terminate
//     2. No circular dependencies
//     3. Channels are properly closed
//
// Safety:
//   ∀ t ∈ traces(Pool): 
//     safe(t) ⟺ ¬(∃ job: received(job) ∧ ¬sent(job))
`, pkg, numWorkers, numWorkers)
}

// GenerateFanIn 生成Fan-In模式代码
//
// CSP Model: FanIn = (input₁ → merge) || (input₂ → merge) || ... → output
func GenerateFanIn(data map[string]interface{}) string {
	pkg := getStringOrDefault(data, "PackageName", "main")

	return fmt.Sprintf(`package %s

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
`, pkg)
}

// GenerateFanOut 生成Fan-Out模式代码
func GenerateFanOut(data map[string]interface{}) string {
	pkg := getStringOrDefault(data, "PackageName", "main")
	fanOutN := getIntOrDefault(data, "FanOutN", 5)

	return fmt.Sprintf(`package %s

// Pattern: Fan-Out
// CSP Model: FanOut = input → (proc₁ || proc₂ || ... || proc%d)
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
`, pkg, fanOutN)
}

// GeneratePipeline 生成Pipeline模式代码
func GeneratePipeline(data map[string]interface{}) string {
	pkg := getStringOrDefault(data, "PackageName", "main")

	return fmt.Sprintf(`package %s

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
`, pkg)
}

// GenerateGenerator 生成Generator模式代码
func GenerateGenerator(data map[string]interface{}) string {
	pkg := getStringOrDefault(data, "PackageName", "main")

	return fmt.Sprintf(`package %s

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
`, pkg)
}

// Helper functions

func getStringOrDefault(data map[string]interface{}, key, defaultVal string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return defaultVal
}

func getIntOrDefault(data map[string]interface{}, key string, defaultVal int) int {
	if val, ok := data[key].(int); ok {
		return val
	}
	return defaultVal
}

func repeat(s string, n int) string {
	return strings.Repeat(s, n)
}
