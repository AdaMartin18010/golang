package main

/*
#include <stdlib.h>
#include <string.h>

// 示例 1: 简单内存泄漏
void simple_leak() {
    void* ptr = malloc(1024);
    // 问题: 忘记 free(ptr)
}

// 示例 2: 字符串泄漏
char* create_string(const char* str) {
    char* result = (char*)malloc(strlen(str) + 1);
    strcpy(result, str);
    return result;
    // 问题: 调用者需要手动释放
}

// 示例 3: Use-After-Free
typedef struct {
    int value;
} Data;

Data* create_data(int val) {
    Data* d = (Data*)malloc(sizeof(Data));
    d->value = val;
    return d;
}

void free_data(Data* d) {
    free(d);
}

int get_value(Data* d) {
    return d->value;
}

// 示例 4: 双重释放
void double_free_example() {
    void* ptr = malloc(100);
    free(ptr);
    free(ptr);  // 问题: 双重释放
}

// 示例 5: 缓冲区溢出
void buffer_overflow() {
    char* buf = (char*)malloc(10);
    strcpy(buf, "This is a very long string");  // 问题: 缓冲区溢出
    free(buf);
}
*/
import "C"
import (
	"fmt"
)

func main() {
	fmt.Println("=== Go 1.25 AddressSanitizer 示例 ===\n")

	// 示例 1: 简单内存泄漏
	fmt.Println("1. 简单内存泄漏:")
	C.simple_leak()
	fmt.Println("   ✅ 调用完成 (但有内存泄漏)\n")

	// 示例 2: 字符串泄漏
	fmt.Println("2. 字符串泄漏:")
	cstr := C.create_string(C.CString("Hello, World!"))
	gostr := C.GoString(cstr)
	fmt.Printf("   创建字符串: %s\n", gostr)
	// 问题: 忘记 C.free(unsafe.Pointer(cstr))
	fmt.Println("   ✅ 字符串创建完成 (但有内存泄漏)\n")

	// 示例 3: Use-After-Free (取消注释会触发错误)
	fmt.Println("3. Use-After-Free (已注释,取消注释会触发错误):")
	// data := C.create_data(42)
	// fmt.Printf("   Value: %d\n", int(C.get_value(data)))
	// C.free_data(data)
	// fmt.Printf("   Value after free: %d\n", int(C.get_value(data)))  // 错误!
	fmt.Println("   ⚠️  已跳过 (会触发 ASan 错误)\n")

	// 示例 4: 双重释放 (取消注释会触发错误)
	fmt.Println("4. 双重释放 (已注释,取消注释会触发错误):")
	// C.double_free_example()  // 错误!
	fmt.Println("   ⚠️  已跳过 (会触发 ASan 错误)\n")

	// 示例 5: 缓冲区溢出 (取消注释会触发错误)
	fmt.Println("5. 缓冲区溢出 (已注释,取消注释会触发错误):")
	// C.buffer_overflow()  // 错误!
	fmt.Println("   ⚠️  已跳过 (会触发 ASan 错误)\n")

	fmt.Println("=== 程序运行完成 ===")
	fmt.Println("\n💡 提示:")
	fmt.Println("   - 使用 'go build -asan -o leak leak_example.go' 编译")
	fmt.Println("   - 运行 './leak' 会检测到内存泄漏")
	fmt.Println("   - 取消注释错误示例会触发 ASan 错误报告")
}
