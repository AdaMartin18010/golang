package main

// Preamble: 紧邻 `import "C"` 上方的注释块。
// 这里可以编写 C 代码。
/*
#include <stdio.h>

// sum 是一个简单的 C 函数，计算两个整数的和。
int sum(int a, int b) {
    printf("[From C] Calculating sum of %d and %d\n", a, b);
    return a + b;
}
*/
import "C"
import "fmt"

func main() {
	fmt.Println("[From Go] Starting CGO basic demonstration.")

	// 定义两个 Go 整数
	a := 10
	b := 20

	fmt.Printf("[From Go] Calling C function 'sum' with arguments: %d, %d\n", a, b)

	// 调用 C 函数 `sum`。
	// Go 的 int 类型会自动转换为 C.int 类型。
	// C.sum 返回一个 C.int，它也会被自动转换回 Go 的 int。
	result := C.sum(C.int(a), C.int(b))

	fmt.Printf("[From Go] Received result from C function: %d\n", result)
}
