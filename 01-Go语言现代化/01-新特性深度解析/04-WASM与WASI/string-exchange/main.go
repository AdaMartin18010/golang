//go:build wasip1

package main

import "unsafe"

// greet 函数接收一个指向 WASM 线性内存的指针（由偏移量和长度表示），
// 将其内容读取为 Go 字符串，处理后返回一个新的字符串。
//
//go:wasmexport greet
func greet(ptr, size uint32) uint64 {
	// 使用 unsafe 将指针和长度转换为 Go 的 []byte 切片，这不会产生拷贝。
	inputBytes := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)

	// 将字节切片转换为字符串
	name := string(inputBytes)

	// 处理字符串
	result := "Hello, " + name + "!"

	// 将结果字符串转换为 []byte
	resultBytes := []byte(result)

	// 在 Go 的内存中分配新的空间来存放结果。
	// 注意：这块内存需要由宿主在读取后手动释放。
	buffer := malloc(uint32(len(resultBytes)))

	// 将结果拷贝到新分配的内存中。
	copy(unsafe.Slice((*byte)(buffer), len(resultBytes)), resultBytes)

	// 返回指向结果字符串的指针和其长度。
	// 我们将指针（32位）和长度（32位）打包成一个 64 位整数返回。
	// 这是在 WASM 中返回多个值的一种常见技巧。
	return (uint64(uintptr(unsafe.Pointer(buffer))) << 32) | uint64(len(resultBytes))
}

//go:wasmexport malloc
// 导出一个内存分配函数，供宿主环境使用。
// 它分配指定大小的字节数组，并返回指向其起始地址的指针。
func malloc(size uint32) *byte {
	// 创建一个足够大的字节切片
	buf := make([]byte, size)
	// 返回切片第一个元素的指针
	return &buf[0]
}

//go:wasmexport free
// 导出一个内存释放函数。
// 注意：在 Go 的垃圾回收机制下，这个函数实际上不是必需的，
// 因为只要没有指针引用 `buf`，它最终会被回收。
// 但导出它是一种好的实践，可以与 C/Rust 等需要手动内存管理的语言保持一致，
// 并明确地表达内存所有权的转移。
func free(ptr *byte, size uint32) {
	// 在 Go 中，我们实际上不需要做任何事。
	// 这个函数的存在主要是为了 API 的完整性。
}

func main() {
	// 模块加载时打印消息
	println("Go WASM module with string support loaded.")
}
