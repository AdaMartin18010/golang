package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// process_string 是一个 C 函数，它接收一个字符串，
// 在其后附加 " (processed by C)"，并返回一个*新分配*的字符串。
char* process_string(const char* input) {
    const char* suffix = " (processed by C)";
    // 分配足够的内存来存放原字符串、后缀和结尾的 '\0'
    char* result = (char*)malloc(strlen(input) + strlen(suffix) + 1);
    if (result == NULL) {
        return NULL;
    }
    strcpy(result, input);
    strcat(result, suffix);
    return result; // 调用者有责任释放这块内存
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("[From Go] Starting CGO string and memory management demo.")

	// 1. 准备一个 Go 字符串
	goString := "Hello from Go"
	fmt.Printf("[From Go] Original string: \"%s\"\n", goString)

	// 2. 将 Go 字符串转换为 C 字符串 (C.CString 在内部调用了 C.malloc)
	cInputString := C.CString(goString)
	// **关键**: 必须手动释放由 C.CString 分配的内存。
	// 使用 defer 确保在函数返回时执行释放操作。
	defer C.free(unsafe.Pointer(cInputString))
	fmt.Printf("[From Go] Converted to C string at address: %p\n", cInputString)

	// 3. 调用 C 函数，传递 C 字符串
	fmt.Println("[From Go] Calling C.process_string...")
	cResultString := C.process_string(cInputString)
	// **关键**: C 函数返回的也是一块新分配的内存，同样需要我们手动释放。
	defer C.free(unsafe.Pointer(cResultString))
	fmt.Printf("[From Go] Received new C string from C at address: %p\n", cResultString)

	// 4. 将返回的 C 字符串转换回 Go 字符串以便使用
	goResultString := C.GoString(cResultString)

	fmt.Printf("[From Go] Converted result back to Go string: \"%s\"\n", goResultString)
	fmt.Println("[From Go] Demo finished. All C memory has been scheduled for release by defer.")
}
