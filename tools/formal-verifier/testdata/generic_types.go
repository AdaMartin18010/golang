package main

import (
	"fmt"
)

// 示例1：简单的泛型函数
func Min[T ~int | ~float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// 示例2：泛型约束 - comparable
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// 示例3：泛型结构体
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// 示例4：泛型接口约束
type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](nums []T) T {
	var sum T
	for _, n := range nums {
		sum += n
	}
	return sum
}

// 示例5：复合约束
type Ordered interface {
	~int | ~int64 | ~float64 | ~string
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 示例6：泛型方法
type Container[T any] struct {
	value T
}

func (c Container[T]) Get() T {
	return c.value
}

func (c *Container[T]) Set(value T) {
	c.value = value
}

// 示例7：多类型参数
func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// 示例8：泛型类型别名 (Go 1.25.3+)
type MyInt[T ~int] T

func (m MyInt[T]) Double() T {
	return T(m) * 2
}

// 示例9：约束中的方法集
type Stringer interface {
	String() string
}

func PrintAll[T Stringer](items []T) {
	for _, item := range items {
		fmt.Println(item.String())
	}
}

// 示例10：类型集约束
type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Abs[T SignedInteger](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

func main() {
	// 使用泛型函数
	fmt.Println(Min(10, 20))
	fmt.Println(Min(3.14, 2.71))

	// 使用Contains
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(Contains(nums, 3))

	// 使用Stack
	stack := &Stack[int]{}
	stack.Push(1)
	stack.Push(2)
	val, ok := stack.Pop()
	if ok {
		fmt.Println(val)
	}

	// 使用Sum
	intSlice := []int{1, 2, 3, 4, 5}
	fmt.Println(Sum(intSlice))

	// 使用Max
	fmt.Println(Max(100, 200))
	fmt.Println(Max("hello", "world"))

	// 使用Container
	c := Container[string]{value: "hello"}
	fmt.Println(c.Get())
	c.Set("world")
	fmt.Println(c.Get())

	// 使用Map
	squared := Map(intSlice, func(n int) int {
		return n * n
	})
	fmt.Println(squared)

	// 使用MyInt
	var mi MyInt[int] = 21
	fmt.Println(mi.Double())

	// 使用Abs
	fmt.Println(Abs(-42))
	fmt.Println(Abs(int64(-100)))
}
