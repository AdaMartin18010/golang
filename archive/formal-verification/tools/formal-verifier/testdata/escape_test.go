package main

import "fmt"

// 逃逸分析测试用例

// 场景1: 局部变量不逃逸（可栈分配）
func stackAlloc() {
	x := 42
	y := x + 1
	fmt.Println(y)
}

// 场景2: 返回局部变量指针（逃逸到堆）
func heapEscape1() *int {
	x := 42
	return &x // x逃逸到堆
}

// 场景3: new分配，不返回（可能栈分配）
func maybeStack() {
	p := new(int)
	*p = 100
	fmt.Println(*p)
}

// 场景4: new分配，返回（逃逸到堆）
func definitelyHeap() *int {
	p := new(int)
	*p = 200
	return p // p逃逸到堆
}

// 场景5: 闭包捕获（逃逸到堆）
func closureCapture() func() int {
	x := 10
	return func() int {
		x++ // x被闭包捕获，逃逸到堆
		return x
	}
}

// 场景6: make slice，小容量（可能栈分配）
func makeSliceSmall() {
	s := make([]int, 10)
	s[0] = 1
	fmt.Println(s)
}

// 场景7: make slice，返回（逃逸到堆）
func makeSliceEscape() []int {
	s := make([]int, 10)
	s[0] = 1
	return s // s逃逸到堆
}

// 场景8: 复合字面量，不逃逸
func compositeLiteral() {
	p := &struct {
		x, y int
	}{1, 2}
	fmt.Println(p.x, p.y)
}

// 场景9: 复合字面量，逃逸
func compositeLiteralEscape() *struct{ x, y int } {
	return &struct {
		x, y int
	}{3, 4} // 字面量逃逸到堆
}

// 场景10: interface逃逸
func interfaceEscape() interface{} {
	x := 100
	return x // x逃逸到堆（装箱）
}

// 场景11: 传递给goroutine（逃逸到堆）
func goroutineEscape() {
	x := 42
	go func() {
		fmt.Println(x) // x被goroutine捕获，逃逸
	}()
}

// 场景12: 数组，不逃逸
func arrayNoEscape() {
	var arr [100]int
	arr[0] = 1
	fmt.Println(arr[0])
}

// 场景13: 大数组，可能逃逸
func largeArrayEscape() [10000]int {
	var arr [10000]int
	arr[0] = 1
	return arr // 大数组可能逃逸
}

// 场景14: map逃逸
func mapEscape() {
	m := make(map[string]int)
	m["key"] = 1
	fmt.Println(m)
}

// 场景15: channel逃逸
func channelEscape() {
	ch := make(chan int, 1)
	ch <- 1
	close(ch)
}

func main() {
	// 调用各种场景
	stackAlloc()
	_ = heapEscape1()
	maybeStack()
	_ = definitelyHeap()
	f := closureCapture()
	_ = f()
	makeSliceSmall()
	_ = makeSliceEscape()
	compositeLiteral()
	_ = compositeLiteralEscape()
	_ = interfaceEscape()
	goroutineEscape()
	arrayNoEscape()
	_ = largeArrayEscape()
	mapEscape()
	channelEscape()
}
