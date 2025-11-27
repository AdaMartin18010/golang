package main

import (
	"fmt"
	"strings"

	"github.com/yourusername/golang/examples/modern-features/01-new-features/01-generic-type-alias/advanced_patterns"
	"github.com/yourusername/golang/examples/modern-features/01-new-features/01-generic-type-alias/real_world_cases"
)

func main() {
	fmt.Println("=== Go 1.25.3 泛型类型别名示例 ===\n")

	// 运行基础示例
	fmt.Println("1. 基础示例:")
	runBasicExamples()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 运行高级模式示例
	fmt.Println("2. 高级模式示例:")
	advanced_patterns.PrintAdvancedUsage()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 运行实际应用案例
	fmt.Println("3. 实际应用案例:")
	real_world_cases.DemonstratePipeline()
}

// runBasicExamples 运行基础示例
func runBasicExamples() {
	// 1. 为泛型map创建别名
	type StringMap[V any] = map[string]V

	// 2. 为泛型slice创建别名
	type DataList[T any] = []T

	// 3. 结合约束（constraints）创建别名
	type Number interface {
		int | int64 | float32 | float64
	}

	type Pair[T Number] struct {
		First  T
		Second T
	}

	// 4. 为泛型函数类型创建别名
	type Calculator[T Number] func(T, T) T

	Add := func[T Number](a, b T) T {
		return a + b
	}

	// 使用 StringMap 别名
	strToIntMap := StringMap[int]{"one": 1, "two": 2}
	fmt.Printf("StringMap[int]: %v\n", strToIntMap)

	// 使用 DataList 别名
	stringList := DataList[string]{"hello", "world"}
	fmt.Printf("DataList[string]: %v\n", stringList)

	// 使用 Pair 别名
	intPair := Pair[int]{First: 10, Second: 20}
	fmt.Printf("Pair[int]: %+v\n", intPair)

	floatPair := Pair[float64]{First: 3.14, Second: 2.71}
	fmt.Printf("Pair[float64]: %+v\n", floatPair)

	// 使用 Calculator 别名
	var intCalculator Calculator[int] = Add[int]
	sum := intCalculator(intPair.First, intPair.Second)
	fmt.Printf("Sum of intPair using Calculator: %d\n", sum)
}
