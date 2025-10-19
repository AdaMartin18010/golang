package performance

/*
#include <stdlib.h>

// noop 是一个空的 C 函数，用于在基准测试中衡量 CGO 调用的纯粹开销。
// 它接收一个指针和大小，但内部不执行任何操作。
void noop(void* p, int size) {
    // Do nothing.
}
*/
import "C"
import (
	"testing"
	"unsafe"
)

const dataSize = 1024

// -- CGO 调用封装 --

// callCgoSafe 是一个标准的、安全的 CGO 调用封装。
// Go 运行时会对传递给 C 的指针进行检查。
func callCgoSafe(p unsafe.Pointer, size int) {
	C.noop(p, C.int(size))
}

// callCgoUnsafe 是一个使用 `//go:cgo_unsafe_args` 优化的 CGO 调用封装。
// 这个编译指令告诉编译器，可以跳过对输入参数的运行时安全检查。
//
// **注意**: 只有在确定传递给 C 的参数不包含任何 Go 指针时，
//
//	使用此指令才是安全的。这里我们传递的是由 C.malloc 分配的内存，
//	因此是安全的。
//
//go:cgo_unsafe_args
func callCgoUnsafe(p unsafe.Pointer, size int) {
	C.noop(p, C.int(size))
}

// -- 基准测试 --

// BenchmarkCgoSafe 测试标准 CGO 调用的性能。
func BenchmarkCgoSafe(b *testing.B) {
	// 在 Go 中调用 C.malloc 来分配一块不受 Go GC 管理的内存。
	// 这是 `//go:cgo_unsafe_args` 的一个安全使用场景。
	p := C.malloc(dataSize)
	defer C.free(p)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		callCgoSafe(p, dataSize)
	}
}

// BenchmarkCgoUnsafe 测试使用 `//go:cgo_unsafe_args` 优化后的 CGO 调用性能。
func BenchmarkCgoUnsafe(b *testing.B) {
	p := C.malloc(dataSize)
	defer C.free(p)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		callCgoUnsafe(p, dataSize)
	}
}

/*
   🔍 预期结果:
   运行 `go test -bench=.` 后，你会观察到 `BenchmarkCgoUnsafe` 的
   `ns/op` (纳秒/每次操作) 值明显低于 `BenchmarkCgoSafe`。

   这证明了 `//go:cgo_unsafe_args` 通过消除运行时指针检查，
   有效地降低了 CGO 调用的开销。

   示例输出:
   goos: linux
   goarch: amd64
   pkg: cgo-performance-demo
   cpu: Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz
   BenchmarkCgoSafe-12          33133887                35.94 ns/op
   BenchmarkCgoUnsafe-12        100000000               10.61 ns/op
   PASS
*/
