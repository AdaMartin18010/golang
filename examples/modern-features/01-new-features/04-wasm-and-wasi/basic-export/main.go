//go:build wasip1

package main

/*
   ======================================================================
   Go WASM/WASI 函数导出基础示例
   ======================================================================

   🎯 目的:
   演示如何使用 `//go:wasmexport` (Go 1.24+) 编译指令将一个 Go 函数
   导出，使其可以被 WASM 宿主环境（Host）调用。

   ⚙️ 如何编译:
   使用 Go 1.24 或更高版本，在当前目录下运行以下命令：
   $ go build -o main.wasm .

   这条命令会生成一个 `main.wasm` 文件。`-o` 参数指定输出文件名。

   🚀 如何运行 (需要 WASM 运行时, 例如 Wasmtime):
   1. 首先，安装 Wasmtime: https://wasmtime.dev/
   2. 使用 Wasmtime 运行编译好的模块，并调用导出的 `add` 函数:
      $ wasmtime run --invoke add main.wasm 40 2

   🔍 预期输出:
   Wasmtime 会加载 `main.wasm` 模块，调用导出的 `add` 函数并传入参数 40 和 2，
   然后打印出函数的返回值。
   > 42

   你也可以看到这样的提示信息：
   > warning: using experimental feature: `wasi:cli/run`
*/

// `//go:wasmexport` 是一个编译器指令，它告诉 Go 编译器：
//  1. 将紧随其后的 `add` 函数标记为可导出的。
//  2. 在生成的 `.wasm` 模块的导出节（Export Section）中，
//     创建一个名为 "add" 的条目指向这个函数。
//
//go:wasmexport add
func add(a, b int) int {
	return a + b
}

// main 函数是必需的，即使它什么也不做。
// 在 `wasip1` 目标下，`main` 函数是程序的入口点，
// 如果 WASM 模块被作为独立的应用程序运行，`main` 将被执行。
// 当我们只是想调用导出的函数时，`main` 可以为空。
func main() {
	println("Go WASM module loaded. Ready to be invoked.")
}
