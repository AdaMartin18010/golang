// Go 1.25 迭代器示例
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== Go 1.25 迭代器示例 ===\n")

	// 1. strings.Lines - 按行迭代
	text := `line 1
line 2
line 3`

	fmt.Println("1. strings.Lines:")
	for line := range strings.Lines(text) {
		fmt.Printf("  %s\n", line)
	}

	// 2. strings.SplitSeq - 分割迭代器
	data := "apple,banana,cherry,date"

	fmt.Println("\n2. strings.SplitSeq:")
	for part := range strings.SplitSeq(data, ",") {
		fmt.Printf("  %s\n", part)
	}

	// 3. strings.FieldsSeq - 字段迭代器
	fields := "  hello   world   go   "

	fmt.Println("\n3. strings.FieldsSeq:")
	for field := range strings.FieldsSeq(fields) {
		fmt.Printf("  [%s]\n", field)
	}

	fmt.Println("\n✅ 迭代器示例完成")
}

