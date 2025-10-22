package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ContextPropagation 演示Context传播模式
// Context用于在goroutine之间传递取消信号、超时和请求范围的值

// WithTimeout 创建一个带超时的Context示例
func WithTimeout(parentCtx context.Context, timeout time.Duration, task func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- task(ctx)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("operation timed out: %w", ctx.Err())
	}
}

// WithCancel 创建一个可取消的Context示例
func WithCancel(parentCtx context.Context, task func(context.Context) error) (error, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parentCtx)

	done := make(chan error, 1)

	go func() {
		done <- task(ctx)
	}()

	go func() {
		select {
		case <-done:
			return
		case <-ctx.Done():
			return
		}
	}()

	return <-done, cancel
}

// WithValue 演示Context值传播
func WithValue(parentCtx context.Context, key, value interface{}, task func(context.Context) error) error {
	ctx := context.WithValue(parentCtx, key, value)
	return task(ctx)
}

// Pipeline 使用Context控制的管道
func ContextAwarePipeline(ctx context.Context, source <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-source:
				if !ok {
					return
				}
				select {
				case out <- v * 2:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

// WorkerPool 支持Context取消的Worker Pool
func ContextAwareWorkerPool(ctx context.Context, numWorkers int, jobs <-chan int) <-chan int {
	results := make(chan int, numWorkers)
	var wg sync.WaitGroup

	// 启动workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					select {
					case results <- job * job:
					case <-ctx.Done():
						return
					}
				}
			}
		}(w)
	}

	// 等待所有workers完成，然后关闭results
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
