package main

import "fmt"

// 函数内联分析测试用例

// 场景1: 小函数，应该内联（cost < 10）
func add(a, b int) int {
	return a + b
}

// 场景2: 叶子函数，倾向内联
func square(x int) int {
	return x * x
}

// 场景3: 中等大小函数（cost < 80）
func compute(x, y int) int {
	temp := x + y
	result := temp * 2
	if result > 100 {
		result = 100
	}
	return result
}

// 场景4: 包含循环，成本高（cost += 30）
func sumArray(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

// 场景5: 递归函数，不应内联
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// 场景6: 包含多次调用，成本高
func complexFunc(a, b, c int) int {
	x := add(a, b)     // +20 调用开销
	y := square(x)     // +20 调用开销
	z := compute(y, c) // +20 调用开销
	return z
}

// 场景7: 包含defer，成本增加
func withDefer(x int) int {
	defer func() {
		fmt.Println("deferred")
	}()
	return x * 2
}

// 场景8: 包含goroutine，成本高
func withGoroutine(x int) {
	go func() {
		fmt.Println(x)
	}()
}

// 场景9: 包含闭包，成本很高
func withClosure(x int) func() int {
	return func() int {
		return x * 2
	}
}

// 场景10: 包含switch，成本增加
func withSwitch(x int) string {
	switch x {
	case 1:
		return "one"
	case 2:
		return "two"
	case 3:
		return "three"
	default:
		return "other"
	}
}

// 场景11: 空函数，应该内联
func noop() {
}

// 场景12: 简单getter，应该内联
type Point struct {
	X, Y int
}

func (p Point) GetX() int {
	return p.X
}

// 场景13: 复杂控制流，不应内联
func complexControl(x int) int {
	result := 0
	for i := 0; i < x; i++ {
		if i%2 == 0 {
			for j := 0; j < i; j++ {
				result += j
			}
		} else {
			result += i
		}
	}
	return result
}

// 场景14: 多返回值
func multiReturn(a, b int) (int, int) {
	return a + b, a - b
}

// 场景15: 有命名返回值
func namedReturn(x int) (result int) {
	result = x * 2
	return
}

// 场景16: panic, 成本增加
func mayPanic(x int) int {
	if x < 0 {
		panic("negative input")
	}
	return x * 2
}

// 场景17: recover, 成本增加
func withRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered")
		}
	}()
}

// 场景18: 大函数，不应内联
func bigFunction(a, b, c, d, e int) int {
	x1 := a + b
	x2 := c + d
	x3 := e + a
	y1 := x1 * x2
	y2 := x2 * x3
	y3 := x3 * x1
	z1 := y1 + y2
	z2 := y2 + y3
	z3 := y3 + y1

	if z1 > 100 {
		z1 = 100
	}
	if z2 > 100 {
		z2 = 100
	}
	if z3 > 100 {
		z3 = 100
	}

	result := z1 + z2 + z3

	for i := 0; i < 10; i++ {
		result += i
	}

	return result
}

func main() {
	// 调用各种函数
	_ = add(1, 2)
	_ = square(3)
	_ = compute(4, 5)
	_ = sumArray([]int{1, 2, 3})
	_ = factorial(5)
	_ = complexFunc(1, 2, 3)
	_ = withDefer(4)
	withGoroutine(5)
	f := withClosure(6)
	_ = f()
	_ = withSwitch(1)
	noop()
	p := Point{1, 2}
	_ = p.GetX()
	_ = complexControl(10)
	_, _ = multiReturn(1, 2)
	_ = namedReturn(3)
	_ = mayPanic(4)
	withRecover()
	_ = bigFunction(1, 2, 3, 4, 5)
}
