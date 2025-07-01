package real_world_cases

import (
	"fmt"
	"strings"
)

// --- 场景: 构建一个类型安全的数据处理管道 (Data Pipeline) ---

// 1. 定义管道中每个阶段的函数签名别名
// Stage[T] 代表管道中的一个处理阶段。
// 它接收一个类型为 T 的输入，并返回一个处理过的类型为 T 的输出。
type Stage[T any] func(T) T

// Pipeline[T] 代表由多个 Stage[T] 组成的完整管道。
type Pipeline[T any] struct {
	stages []Stage[T]
}

// NewPipeline 创建一个新的数据处理管道。
func NewPipeline[T any](stages ...Stage[T]) Pipeline[T] {
	return Pipeline[T]{stages: stages}
}

// Execute 方法按顺序执行管道中的所有阶段。
func (p Pipeline[T]) Execute(input T) T {
	currentData := input
	for _, stage := range p.stages {
		currentData = stage(currentData)
	}
	return currentData
}

// --- 示例：处理用户输入字符串的管道 ---

// 定义一个具体的数据类型
type UserInput string

// 定义具体的处理阶段
// Stage 1: 去除首尾空格
func TrimSpaceStage(input UserInput) UserInput {
	fmt.Println("Executing: TrimSpaceStage")
	return UserInput(strings.TrimSpace(string(input)))
}

// Stage 2: 转换为大写
func ToUpperStage(input UserInput) UserInput {
	fmt.Println("Executing: ToUpperStage")
	return UserInput(strings.ToUpper(string(input)))
}

// Stage 3: 添加问候语
func AddGreetingStage(input UserInput) UserInput {
	fmt.Println("Executing: AddGreetingStage")
	return UserInput("Hello, " + string(input) + "!")
}

// DemonstratePipeline 函数展示了如何构建和执行数据处理管道。
func DemonstratePipeline() {
	fmt.Println("--- Real World Case: Type-Safe Data Pipeline ---")

	// 使用别名来引用具体的处理阶段
	type StringStage = Stage[UserInput]

	// 创建一个处理用户输入的管道
	stringPipeline := NewPipeline[UserInput](
		StringStage(TrimSpaceStage),
		StringStage(ToUpperStage),
		StringStage(AddGreetingStage),
	)

	// 初始输入
	initialInput := UserInput("   gopher   ")
	fmt.Printf("Initial input: '%s'\n\n", initialInput)

	// 执行管道
	finalResult := stringPipeline.Execute(initialInput)

	// 打印最终结果
	fmt.Printf("\nFinal result: '%s'\n", finalResult)

	// 预期输出:
	// Initial input: '   gopher   '
	//
	// Executing: TrimSpaceStage
	// Executing: ToUpperStage
	// Executing: AddGreetingStage
	//
	// Final result: 'Hello, GOPHER!'
}
