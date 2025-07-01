package basic_examples

import "fmt"

// 1. 为泛型map创建别名
// StringMap 是一个键为 string，值为任意类型 V 的 map 的别名。
type StringMap[V any] = map[string]V

// 2. 为泛型slice创建别名
// DataList 是一个任意类型 T 的切片的别名。
type DataList[T any] = []T

// 3. 结合约束（constraints）创建别名
// 定义一个约束，只允许整数和浮点数类型。
type Number interface {
	int | int64 | float32 | float64
}

// Pair 是一个包含两个相同数字类型元素的结构体。
type Pair[T Number] struct {
	First  T
	Second T
}

// 4. 为泛型函数类型创建别名
// Calculator 是一个接受两个相同数字类型参数并返回同类型结果的函数类型别名。
type Calculator[T Number] func(T, T) T

// Add 是一个符合 Calculator 签名以供演示的函数。
func Add[T Number](a, b T) T {
	return a + b
}

// PrintUsage 函数演示了如何使用上面定义的别名。
func PrintUsage() {
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
