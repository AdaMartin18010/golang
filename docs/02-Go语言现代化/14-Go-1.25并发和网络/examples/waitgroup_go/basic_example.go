package main

import (
	"fmt"
	"sync"
	"time"
)

// 示例 1: 基本使用
func basicExample() {
	fmt.Println("=== 示例 1: 基本使用 ===")

	var wg sync.WaitGroup

	// 使用 WaitGroup.Go() 启动 goroutine
	wg.Go(func() {
		fmt.Println("Task 1: Starting")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Task 1: Done")
	})

	wg.Go(func() {
		fmt.Println("Task 2: Starting")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Task 2: Done")
	})

	wg.Go(func() {
		fmt.Println("Task 3: Starting")
		time.Sleep(150 * time.Millisecond)
		fmt.Println("Task 3: Done")
	})

	// 等待所有任务完成
	wg.Wait()
	fmt.Println("All tasks completed!")
}

// 示例 2: 并行处理切片
func sliceProcessing() {
	fmt.Println("=== 示例 2: 并行处理切片 ===")

	items := []string{"apple", "banana", "cherry", "date", "elderberry"}
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Go(func() {
			// 模拟处理
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Processed: %s\n", item)
		})
	}

	wg.Wait()
	fmt.Println("All items processed!")
}

// 示例 3: 限制并发数
func limitedConcurrency() {
	fmt.Println("=== 示例 3: 限制并发数 ===")

	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	maxConcurrency := 3

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

	for _, item := range items {
		sem <- struct{}{} // 获取信号量

		wg.Go(func() {
			defer func() { <-sem }() // 释放信号量

			fmt.Printf("Processing item %d\n", item)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("Completed item %d\n", item)
		})
	}

	wg.Wait()
	fmt.Println("All items processed with limited concurrency!")
}

// 示例 4: 收集结果
func collectResults() {
	fmt.Println("=== 示例 4: 收集结果 ===")

	items := []int{1, 2, 3, 4, 5}
	results := make(chan int, len(items))

	var wg sync.WaitGroup

	for _, item := range items {
		wg.Go(func() {
			// 计算平方
			result := item * item
			results <- result
		})
	}

	// 等待所有任务完成
	wg.Wait()
	close(results)

	// 收集结果
	fmt.Print("Results: ")
	for result := range results {
		fmt.Printf("%d ", result)
	}
	fmt.Println("")
}

// 示例 5: 错误处理
func errorHandling() {
	fmt.Println("=== 示例 5: 错误处理 ===")

	items := []int{1, 2, 3, 4, 5}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	for _, item := range items {
		wg.Go(func() {
			// 模拟处理,可能失败
			if item%2 == 0 {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to process item %d", item))
				mu.Unlock()
			} else {
				fmt.Printf("Successfully processed item %d\n", item)
			}
		})
	}

	wg.Wait()

	if len(errors) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}
	fmt.Println()
}

// 示例 6: 对比传统方式
func comparison() {
	fmt.Println("=== 示例 6: 传统方式 vs WaitGroup.Go() ===")

	// 传统方式
	fmt.Println("传统方式:")
	{
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("  Task 1 (traditional)")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("  Task 2 (traditional)")
		}()

		wg.Wait()
	}

	// WaitGroup.Go() 方式
	fmt.Println("\nWaitGroup.Go() 方式:")
	{
		var wg sync.WaitGroup

		wg.Go(func() {
			fmt.Println("  Task 1 (WaitGroup.Go)")
		})

		wg.Go(func() {
			fmt.Println("  Task 2 (WaitGroup.Go)")
		})

		wg.Wait()
	}

	fmt.Println("\n代码更简洁,更易读!")
}

func main() {
	fmt.Println("Go 1.25 WaitGroup.Go() 示例")

	basicExample()
	sliceProcessing()
	limitedConcurrency()
	collectResults()
	errorHandling()
	comparison()

	fmt.Println("🎉 所有示例运行完成!")
}
