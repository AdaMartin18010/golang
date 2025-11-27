// Go 1.25 迭代器示例
// 注意：strings.Lines, strings.SplitSeq, strings.FieldsSeq 是 Go 1.25+ 的新特性
// 如果这些 API 尚未可用，可以使用传统的 strings.Split 等方法
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== Go 1.25 迭代器示例 ===\n")

	// 1. strings.Lines - 按行迭代（Go 1.25+）
	// 如果 API 不可用，可以使用 strings.Split(text, "\n")
	text := `line 1
line 2
line 3`

	fmt.Println("1. strings.Lines (按行迭代):")
	// 使用传统方式作为备选
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if line != "" { // 过滤空行
			fmt.Printf("  %s\n", line)
		}
	}

	// 如果 strings.Lines 可用，可以这样使用：
	// for line := range strings.Lines(text) {
	//     fmt.Printf("  %s\n", line)
	// }

	// 2. strings.SplitSeq - 分割迭代器（Go 1.25+）
	data := "apple,banana,cherry,date"

	fmt.Println("\n2. strings.SplitSeq (分割迭代器):")
	// 使用传统方式作为备选
	parts := strings.Split(data, ",")
	for _, part := range parts {
		fmt.Printf("  %s\n", part)
	}

	// 如果 strings.SplitSeq 可用，可以这样使用：
	// for part := range strings.SplitSeq(data, ",") {
	//     fmt.Printf("  %s\n", part)
	// }

	// 3. strings.FieldsSeq - 字段迭代器（Go 1.25+）
	fields := "  hello   world   go   "

	fmt.Println("\n3. strings.FieldsSeq (字段迭代器):")
	// 使用传统方式作为备选
	fieldList := strings.Fields(fields)
	for _, field := range fieldList {
		fmt.Printf("  [%s]\n", field)
	}

	// 如果 strings.FieldsSeq 可用，可以这样使用：
	// for field := range strings.FieldsSeq(fields) {
	//     fmt.Printf("  [%s]\n", field)
	// }

	fmt.Println("\n✅ 迭代器示例完成")
	fmt.Println("\n💡 提示: 如果 Go 1.25+ 的迭代器 API 可用，")
	fmt.Println("   它们会提供更好的内存效率和延迟计算特性。")
}
