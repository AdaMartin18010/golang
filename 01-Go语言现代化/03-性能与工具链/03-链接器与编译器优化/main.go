package main

import "fmt"

/*
   ======================================================================
   Go 编译器优化观察室
   ======================================================================

   🎯 目的:
   本文件包含一系列函数，旨在与 `go build -gcflags="-m"` 命令配合使用，
   以直观地观察 Go 编译器的关键优化决策，如函数内联、逃逸分析和
   边界检查消除。

   ⚙️ 如何使用:
   在当前目录下运行以下命令，并观察输出：
   $ go build -gcflags="-m" .

   (使用 `-m -m` 或 `-m=2` 可以获得更详细的信息)

   🔍 你会看到什么:
   - `can inline ...`: 编译器提示一个函数可以被内联。
   - `inlining call to ...`: 编译器决定将一个函数调用内联。
   - `does not escape`: 变量被成功分配在栈上。
   - `escapes to heap`: 变量因"逃逸"被分配在堆上。
   - `slice bounds check ... eliminated`: 编译器成功移除了边界检查。
*/

// --- 1. 函数内联 (Inlining) ---

// canInline 是一个足够简单的函数，符合被内联的条件。
func canInline(a, b int) int {
	return a + b
}

// cannotInline 因为包含复杂或当前编译器不支持内联的特性（如循环），
// 所以通常不会被内联。
func cannotInline(s []int) int {
	var sum int
	for _, v := range s {
		sum += v
	}
	return sum
}

// --- 2. 逃逸分析 (Escape Analysis) ---

// stackAlloc 函数返回一个 int 值。
// 内部创建的 User 结构体不会"逃逸"，因此会被分配在栈上。
func stackAlloc() int {
	user := User{ID: 1, Name: "on-stack"}
	return user.ID
}

// heapAlloc 函数返回一个指向 User 结构体的指针。
// 因为 `&user` 这个引用"逃逸"出了函数的作用域，
// 所以 `user` 变量必须被分配在堆上。
func heapAlloc() *User {
	user := User{ID: 2, Name: "on-heap"}
	return &user
}

type User struct {
	ID   int
	Name string
}

// --- 3. 边界检查消除 (Bounds Check Elimination) ---

// boundsCheckEliminated 演示了编译器如何消除不必要的边界检查。
func boundsCheckEliminated(s []int) {
	// 编译器知道 s[2] 的访问是安全的，因为它刚刚检查过 s 的长度。
	if len(s) >= 3 {
		_ = s[0] // 边界检查被消除
		_ = s[1] // 边界检查被消除
		_ = s[2] // 边界检查被消除
	}
}

// boundsCheckNeeded 演示了编译器无法消除边界检查的场景。
func boundsCheckNeeded(s []int, i int) {
	// 编译器无法在编译时确定 `i` 的值，
	// 所以必须在运行时保留对 s[i] 的边界检查。
	_ = s[i]
}

func main() {
	// Inlining Demo
	a := canInline(1, 2)
	b := cannotInline([]int{1, 2, 3})
	fmt.Println("Inlining demo:", a, b)

	// Escape Analysis Demo
	c := stackAlloc()
	d := heapAlloc()
	fmt.Println("Escape analysis demo:", c, d.Name)

	// Bounds Check Demo
	s := []int{10, 20, 30}
	boundsCheckEliminated(s)
	boundsCheckNeeded(s, 1)
	fmt.Println("Bounds check demo finished.")
}
