// Go 1.25 unique包示例
package main

import (
	"fmt"
	"runtime"
	"unique"
)

func main() {
	fmt.Println("=== Go 1.25 unique包示例 ===\n")

	// 1. 字符串规范化
	h1 := unique.Make("hello world")
	h2 := unique.Make("hello world")
	h3 := unique.Make("different")

	fmt.Println("1. 字符串规范化:")
	fmt.Printf("  h1 == h2: %v\n", h1 == h2) // true
	fmt.Printf("  h1 == h3: %v\n", h1 == h3) // false
	fmt.Printf("  h1.Value(): %s\n", h1.Value())

	// 2. 结构体规范化
	type Point struct {
		X, Y int
	}

	p1 := unique.Make(Point{X: 1, Y: 2})
	p2 := unique.Make(Point{X: 1, Y: 2})
	p3 := unique.Make(Point{X: 3, Y: 4})

	fmt.Println("\n2. 结构体规范化:")
	fmt.Printf("  p1 == p2: %v\n", p1 == p2) // true
	fmt.Printf("  p1 == p3: %v\n", p1 == p3) // false
	fmt.Printf("  p1.Value(): %+v\n", p1.Value())

	// 3. 内存对比
	var m runtime.MemStats

	handles := make([]unique.Handle[string], 100000)
	for i := range handles {
		handles[i] = unique.Make("repeated string content")
	}

	runtime.ReadMemStats(&m)
	fmt.Printf("\n3. 内存占用: %d KB\n", m.Alloc/1024)

	fmt.Println("\n✅ unique包示例完成")
}

