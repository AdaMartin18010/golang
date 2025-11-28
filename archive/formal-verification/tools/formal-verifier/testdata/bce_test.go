package main

import "fmt"

// 边界检查消除 (BCE) 测试用例

// 场景1: 常量索引，应该消除
func constantIndex() {
	arr := []int{1, 2, 3, 4, 5}
	x := arr[0] // 常量索引，BCE ✓
	y := arr[2] // 常量索引，BCE ✓
	fmt.Println(x, y)
}

// 场景2: range循环，应该消除
func rangeLoop() {
	arr := []int{1, 2, 3, 4, 5}
	for i := range arr {
		x := arr[i] // range循环索引，BCE ✓
		fmt.Println(x)
	}
}

// 场景3: 简单for循环，可能消除
func simpleForLoop() {
	arr := []int{1, 2, 3, 4, 5}
	for i := 0; i < len(arr); i++ {
		x := arr[i] // 循环归纳变量，BCE ✓
		fmt.Println(x)
	}
}

// 场景4: 条件保护，应该消除
func withCondition(arr []int, idx int) {
	if idx >= 0 && idx < len(arr) {
		x := arr[idx] // 条件已检查，BCE ✓
		fmt.Println(x)
	}
}

// 场景5: 重复访问，第二次应消除
func repeatedAccess() {
	arr := []int{1, 2, 3, 4, 5}
	x := arr[2] // 第一次访问
	y := arr[2] // 重复访问，BCE ✓
	fmt.Println(x, y)
}

// 场景6: 未知索引，不能消除
func unknownIndex(arr []int, idx int) {
	x := arr[idx] // 未知索引，BCE ✗
	fmt.Println(x)
}

// 场景7: 越界索引（编译时可检测）
func outOfBounds() {
	arr := [3]int{1, 2, 3}
	// x := arr[5]  // 编译错误，不包含在测试中
	x := arr[2] // 边界内，BCE ✓
	fmt.Println(x)
}

// 场景8: 二维数组，多次检查
func twoDimensional() {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	for i := range matrix {
		for j := range matrix[i] {
			x := matrix[i][j] // 两层range，BCE ✓
			fmt.Println(x)
		}
	}
}

// 场景9: 切片的切片，部分可消除
func sliceOfSlice() {
	arr := []int{1, 2, 3, 4, 5}
	sub := arr[1:4] // 切片操作
	x := sub[0]     // 常量索引，BCE ✓
	y := sub[1]     // 常量索引，BCE ✓
	fmt.Println(x, y)
}

// 场景10: 数组指针
func arrayPointer() {
	arr := [5]int{1, 2, 3, 4, 5}
	p := &arr
	x := p[2] // 常量索引，BCE ✓
	fmt.Println(x)
}

// 场景11: 字符串索引
func stringIndex() {
	s := "hello"
	c := s[0] // 常量索引，BCE ✓
	fmt.Println(string(c))
}

// 场景12: 复杂表达式索引
func complexIndexExpr() {
	arr := []int{1, 2, 3, 4, 5}
	i := 1
	j := 1
	x := arr[i+j] // 复杂索引，BCE ✗
	fmt.Println(x)
}

// 场景13: 循环中修改索引
func modifiedIndex() {
	arr := []int{1, 2, 3, 4, 5}
	for i := 0; i < len(arr); i++ {
		if i%2 == 0 {
			i++ // 修改索引
		}
		if i < len(arr) {
			x := arr[i] // 条件保护，BCE ✓
			fmt.Println(x)
		}
	}
}

// 场景14: 向后遍历
func reverseLoop() {
	arr := []int{1, 2, 3, 4, 5}
	for i := len(arr) - 1; i >= 0; i-- {
		x := arr[i] // 循环归纳变量，BCE ✓
		fmt.Println(x)
	}
}

// 场景15: 嵌套循环
func nestedLoop() {
	arr := []int{1, 2, 3, 4, 5}
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			idx := i + j
			if idx < len(arr) {
				x := arr[idx] // 条件保护，BCE ✓
				fmt.Println(x)
			}
		}
	}
}

// 场景16: append使用
func withAppend() {
	arr := []int{1, 2, 3}
	arr = append(arr, 4)
	x := arr[3] // 常量索引，BCE ✓
	fmt.Println(x)
}

// 场景17: make with capacity
func makeWithCap() {
	arr := make([]int, 5, 10)
	arr[0] = 1 // 常量索引，BCE ✓
	arr[4] = 5 // 常量索引，BCE ✓
	fmt.Println(arr)
}

// 场景18: 多个数组
func multipleArrays() {
	arr1 := []int{1, 2, 3}
	arr2 := []int{4, 5, 6}

	for i := range arr1 {
		x := arr1[i] // range索引，BCE ✓
		if i < len(arr2) {
			y := arr2[i] // 条件保护，BCE ✓
			fmt.Println(x, y)
		}
	}
}

func main() {
	constantIndex()
	rangeLoop()
	simpleForLoop()
	withCondition([]int{1, 2, 3}, 1)
	repeatedAccess()
	unknownIndex([]int{1, 2, 3}, 1)
	outOfBounds()
	twoDimensional()
	sliceOfSlice()
	arrayPointer()
	stringIndex()
	complexIndexExpr()
	modifiedIndex()
	reverseLoop()
	nestedLoop()
	withAppend()
	makeWithCap()
	multipleArrays()
}
