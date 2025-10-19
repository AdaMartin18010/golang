package for_loop_semantics

import (
	"fmt"
	"sync"
)

/*
   ======================================================================
   for 循环变量语义变更 - 行为对比示例
   ======================================================================

   🎯 目的:
   通过具体代码示例，直观地对比 Go 1.22 前后 for 循环变量在
   并发场景下的行为差异。

   ⚙️ 如何运行:
   - 使用 Go 1.21 或更早版本编译运行，观察 `IncorrectBehaviorBeforeGo122` 的输出。
     $ go run .
   - 使用 Go 1.22 或更高版本编译运行，观察 `CorrectBehaviorInGo122` 的输出
     以及打印出的不同的内存地址。
     $ go run .
*/

// IncorrectBehaviorBeforeGo122 演示了在 Go 1.22 之前存在的经典 for 循环变量陷阱。
func IncorrectBehaviorBeforeGo122() {
	fmt.Println("--- 1. Demonstrating Incorrect Behavior (Pre-Go 1.22) ---")
	items := []string{"apple", "banana", "cherry"}
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 闭包捕获的是同一个'item'变量的引用。
			// 当goroutine执行时，循环已结束，'item'的值为"cherry"。
			fmt.Printf("Incorrectly captured item: %s\n", item)
		}()
	}
	wg.Wait()
	fmt.Println("Expected output: 'cherry' printed 3 times.")
	fmt.Println()
}

// CorrectBehaviorWithShadowing 演示了在 Go 1.22 之前修复此问题的传统方法。
func CorrectBehaviorWithShadowing() {
	fmt.Println("--- 2. Demonstrating Traditional Fix (Shadowing) ---")
	items := []string{"apple", "banana", "cherry"}
	var wg sync.WaitGroup

	for _, item := range items {
		// 创建一个循环体内的局部变量 `item` 来 "遮蔽" 循环变量。
		// 这个新的 `item` 在每次迭代时都是一个新变量。
		item := item
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 闭包现在捕获的是新创建的局部变量的副本。
			fmt.Printf("Correctly captured item with shadowing: %s\n", item)
		}()
	}
	wg.Wait()
	fmt.Println("Expected output: apple, banana, cherry (in any order).")
	fmt.Println()
}

// CorrectBehaviorInGo122 演示了 Go 1.22+ 中修正后的行为。
// 注意：此函数只有在用 Go 1.22+ 编译时才能展示出预期的行为。
func CorrectBehaviorInGo122() {
	fmt.Println("--- 3. Demonstrating Correct Behavior (Go 1.22+) ---")
	items := []string{"apple", "banana", "cherry"}
	var wg sync.WaitGroup

	for _, item := range items {
		// 在 Go 1.22+ 中，`item` 在每次迭代时都是一个全新的变量。
		// 我们打印它的内存地址来证明这一点。
		fmt.Printf("Loop iteration with item '%s' at address: %p\n", item, &item)
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 闭包捕获的是本次迭代中新创建的 `item` 变量。
			fmt.Printf("Correctly captured item in Go 1.22+: %s (from address: %p)\n", item, &item)
		}()
	}
	wg.Wait()
	fmt.Println("Expected output: apple, banana, cherry (in any order), with different addresses printed.")
	fmt.Println()
}

// main函数用于在一个地方调用所有示例。
func main() {
	IncorrectBehaviorBeforeGo122()
	CorrectBehaviorWithShadowing()
	CorrectBehaviorInGo122()
}
