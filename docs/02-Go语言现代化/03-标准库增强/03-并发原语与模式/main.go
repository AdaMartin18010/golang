package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	demonstrateAfterFunc()
	demonstrateWithoutCancel()
	demonstrateOnceFunc()
	demonstrateAtomicTypes()
}

// --- 1. context.AfterFunc (Go 1.21+) ---
func demonstrateAfterFunc() {
	fmt.Println("--- 1. Demonstrating context.AfterFunc ---")

	// 创建一个100毫秒后超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 注册一个当 ctx 完成时执行的函数
	stop := context.AfterFunc(ctx, func() {
		fmt.Println("Cleanup function from AfterFunc executed.")
	})

	// 等待足够长的时间以确保 context 超时
	time.Sleep(200 * time.Millisecond)

	// stop 函数可以用于提前取消回调，这里我们只是演示调用它
	// 在这个例子中，由于ctx已完成，再次调用stop不会有任何效果
	if !stop() {
		fmt.Println("The cleanup function had already been called.")
	}
	fmt.Println("--- End of AfterFunc ---\n")
}

// --- 2. context.WithoutCancel (Go 1.21+) ---
func demonstrateWithoutCancel() {
	fmt.Println("--- 2. Demonstrating context.WithoutCancel ---")

	parentCtx, cancelParent := context.WithCancel(context.Background())

	// 立即取消父 context
	cancelParent()

	// 检查父 context 是否已取消
	if parentCtx.Err() != nil {
		fmt.Println("Parent context is canceled as expected.")
	}

	// 创建一个忽略父 context 取消信号的子 context
	childCtx := context.WithoutCancel(parentCtx)

	// 检查子 context 是否被取消
	if childCtx.Err() == nil {
		fmt.Println("Child context (created with WithoutCancel) is NOT canceled.")
	}
	fmt.Println("--- End of WithoutCancel ---\n")
}

// --- 3. sync.OnceFunc (Go 1.21+) ---
var setupCount int

func setupDatabase() {
	fmt.Println("Setting up the database connection...")
	setupCount++
}

func demonstrateOnceFunc() {
	fmt.Println("--- 3. Demonstrating sync.OnceFunc ---")

	// 使用新的 sync.OnceFunc
	initOnce := sync.OnceFunc(setupDatabase)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d is trying to initialize...\n", id)
			initOnce()
		}(i)
	}
	wg.Wait()

	fmt.Printf("The setup function was called %d time(s).\n", setupCount)
	fmt.Println("--- End of OnceFunc ---\n")
}

// --- 4. atomic Types (Go 1.19+) ---
func demonstrateAtomicTypes() {
	fmt.Println("--- 4. Demonstrating atomic Types ---")

	// 旧方式 (pre-Go 1.19)
	var oldCounter int64
	atomic.AddInt64(&oldCounter, 5)
	atomic.StoreInt64(&oldCounter, 10)
	fmt.Printf("Old atomic counter value: %d\n", atomic.LoadInt64(&oldCounter))

	// 新方式 (Go 1.19+)
	var newCounter atomic.Int64
	newCounter.Add(5)
	newCounter.Store(10)
	fmt.Printf("New atomic.Int64 counter value: %d\n", newCounter.Load())

	// 新方式提供了更强的类型安全和更清晰的 API
	// 例如，下面的代码无法编译，因为它试图对一个非原子类型进行原子操作：
	// var regularInt int32
	// regularInt.Add(1) // COMPILE ERROR

	fmt.Println("--- End of atomic Types ---\n")
}
