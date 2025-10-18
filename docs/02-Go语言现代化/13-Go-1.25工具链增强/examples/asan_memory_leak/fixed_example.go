package main

/*
#include <stdlib.h>
#include <string.h>

// ✅ 正确的字符串创建函数
char* create_string_safe(const char* str) {
    char* result = (char*)malloc(strlen(str) + 1);
    if (result) {
        strcpy(result, str);
    }
    return result;
}

// ✅ 正确的数据结构管理
typedef struct {
    int value;
} Data;

Data* create_data_safe(int val) {
    Data* d = (Data*)malloc(sizeof(Data));
    if (d) {
        d->value = val;
    }
    return d;
}

void free_data_safe(Data** d) {
    if (d && *d) {
        free(*d);
        *d = NULL;  // 防止双重释放
    }
}

int get_value_safe(Data* d) {
    if (d) {
        return d->value;
    }
    return -1;  // 错误值
}

// ✅ 正确的缓冲区处理
char* safe_copy(const char* str) {
    if (!str) return NULL;

    size_t len = strlen(str);
    char* buf = (char*)malloc(len + 1);
    if (buf) {
        strncpy(buf, str, len);
        buf[len] = '\0';
    }
    return buf;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("=== Go 1.25 AddressSanitizer - 正确示例 ===\n")

	// 示例 1: 正确的字符串管理
	fmt.Println("1. 正确的字符串管理:")
	{
		input := C.CString("Hello, Safe World!")
		defer C.free(unsafe.Pointer(input))

		cstr := C.create_string_safe(input)
		if cstr != nil {
			gostr := C.GoString(cstr)
			fmt.Printf("   创建字符串: %s\n", gostr)
			C.free(unsafe.Pointer(cstr)) // ✅ 正确释放
		}
	}
	fmt.Println("   ✅ 无内存泄漏\n")

	// 示例 2: 正确的数据结构管理
	fmt.Println("2. 正确的数据结构管理:")
	{
		data := C.create_data_safe(42)
		if data != nil {
			fmt.Printf("   Value: %d\n", int(C.get_value_safe(data)))
			C.free_data_safe(&data) // ✅ 正确释放并设置为 NULL

			// 尝试再次获取值 (安全的,因为会检查 NULL)
			value := int(C.get_value_safe(data))
			fmt.Printf("   Value after free (安全): %d\n", value)
		}
	}
	fmt.Println("   ✅ 无 Use-After-Free\n")

	// 示例 3: 正确的缓冲区处理
	fmt.Println("3. 正确的缓冲区处理:")
	{
		input := C.CString("This is a very long string that is safely handled")
		defer C.free(unsafe.Pointer(input))

		buf := C.safe_copy(input)
		if buf != nil {
			result := C.GoString(buf)
			fmt.Printf("   复制结果: %s\n", result)
			C.free(unsafe.Pointer(buf)) // ✅ 正确释放
		}
	}
	fmt.Println("   ✅ 无缓冲区溢出\n")

	// 示例 4: 批量处理 (正确管理资源)
	fmt.Println("4. 批量处理 (正确管理资源):")
	{
		const count = 100
		for i := 0; i < count; i++ {
			input := C.CString(fmt.Sprintf("Item %d", i))

			buf := C.safe_copy(input)
			if buf != nil {
				C.free(unsafe.Pointer(buf))
			}

			C.free(unsafe.Pointer(input))
		}
		fmt.Printf("   处理 %d 个项目\n", count)
	}
	fmt.Println("   ✅ 无累积泄漏\n")

	// 示例 5: 错误处理模式
	fmt.Println("5. 错误处理模式:")
	{
		processWithCleanup := func(input string) error {
			cInput := C.CString(input)
			defer C.free(unsafe.Pointer(cInput))

			buf := C.safe_copy(cInput)
			if buf == nil {
				return fmt.Errorf("failed to copy string")
			}
			defer C.free(unsafe.Pointer(buf))

			// 处理数据...
			result := C.GoString(buf)
			fmt.Printf("   处理结果: %s\n", result)

			return nil
		}

		if err := processWithCleanup("Test Data"); err != nil {
			fmt.Printf("   错误: %v\n", err)
		}
	}
	fmt.Println("   ✅ 错误时也正确清理\n")

	fmt.Println("=== 程序运行完成 ===")
	fmt.Println("\n💡 提示:")
	fmt.Println("   - 使用 'go build -asan -o fixed fixed_example.go' 编译")
	fmt.Println("   - 运行 './fixed' 不会检测到任何内存问题")
	fmt.Println("   - 这是正确的内存管理模式")
	fmt.Println("\n🎯 最佳实践:")
	fmt.Println("   1. 使用 defer 确保资源释放")
	fmt.Println("   2. 检查 NULL 指针")
	fmt.Println("   3. 释放后设置为 NULL")
	fmt.Println("   4. 使用边界检查的复制函数")
	fmt.Println("   5. 在错误路径中也要清理资源")
}
