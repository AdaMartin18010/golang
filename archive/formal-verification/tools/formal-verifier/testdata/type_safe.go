package main

import "fmt"

// 示例1：类型安全的代码
func safeExample1() {
	var x int = 42
	var y int = 10
	z := x + y
	fmt.Println(z)
}

// 示例2：类型安全的函数
func safeAdd(a, b int) int {
	return a + b
}

// 示例3：类型安全的结构体
type Point struct {
	X int
	Y int
}

func (p Point) Distance() float64 {
	return float64(p.X*p.X + p.Y*p.Y)
}

// 示例4：类型安全的接口
type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 示例5：类型安全的切片操作
func safeSliceOps() {
	nums := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	fmt.Println(sum)
}

// 示例6：类型安全的映射操作
func safeMapOps() {
	m := make(map[string]int)
	m["one"] = 1
	m["two"] = 2

	if v, ok := m["one"]; ok {
		fmt.Println(v)
	}
}

// 示例7：类型安全的指针操作
func safePointerOps() {
	x := 42
	p := &x
	*p = 100
	fmt.Println(x)
}

// 示例8：类型安全的类型断言
func safeTypeAssertion(i interface{}) {
	if str, ok := i.(string); ok {
		fmt.Println(str)
	}
}

// 示例9：类型安全的类型转换
func safeTypeConversion() {
	var i int = 42
	var f float64 = float64(i)
	fmt.Println(f)
}

// 示例10：类型安全的闭包
func safeClosureExample() func() int {
	counter := 0
	return func() int {
		counter++
		return counter
	}
}

func main() {
	safeExample1()
	result := safeAdd(10, 20)
	fmt.Println(result)

	p := Point{X: 3, Y: 4}
	fmt.Println(p.Distance())

	r := Rectangle{Width: 10, Height: 5}
	fmt.Println(r.Area())

	safeSliceOps()
	safeMapOps()
	safePointerOps()
	safeTypeAssertion("hello")
	safeTypeConversion()

	counter := safeClosureExample()
	fmt.Println(counter())
	fmt.Println(counter())
}
