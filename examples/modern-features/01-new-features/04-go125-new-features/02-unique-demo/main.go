// Go 1.25 unique包示例
// 注意：unique 包是 Go 1.25+ 的实验性功能
// 如果包不可用，本示例展示其概念和预期行为
package main

import (
	"fmt"
	"runtime"
	"sync"
)

// 如果 unique 包不可用，我们创建一个简单的模拟实现
// 实际使用中应该使用标准库的 unique 包

// UniqueHandle 模拟 unique.Handle
type UniqueHandle[T comparable] struct {
	value T
	mu    sync.Mutex
	pool  map[T]*T
}

var (
	uniquePool = make(map[interface{}]interface{})
	poolMu     sync.RWMutex
)

// MakeUnique 模拟 unique.Make，创建规范化值
func MakeUnique[T comparable](value T) T {
	poolMu.RLock()
	if existing, ok := uniquePool[value]; ok {
		poolMu.RUnlock()
		return existing.(T)
	}
	poolMu.RUnlock()

	poolMu.Lock()
	defer poolMu.Unlock()

	// 双重检查
	if existing, ok := uniquePool[value]; ok {
		return existing.(T)
	}

	uniquePool[value] = value
	return value
}

func main() {
	fmt.Println("=== Go 1.25 unique包示例 ===\n")
	fmt.Println("注意: 这是 unique 包概念的演示")
	fmt.Println("实际使用中应使用标准库的 unique 包\n")

	// 1. 字符串规范化
	s1 := MakeUnique("hello world")
	s2 := MakeUnique("hello world")
	s3 := MakeUnique("different")

	fmt.Println("1. 字符串规范化:")
	fmt.Printf("  s1 == s2: %v (相同值应该共享)\n", s1 == s2) // true
	fmt.Printf("  s1 == s3: %v (不同值)\n", s1 == s3)         // false
	fmt.Printf("  s1: %s\n", s1)

	// 2. 结构体规范化
	type Point struct {
		X, Y int
	}

	p1 := MakeUnique(Point{X: 1, Y: 2})
	p2 := MakeUnique(Point{X: 1, Y: 2})
	p3 := MakeUnique(Point{X: 3, Y: 4})

	fmt.Println("\n2. 结构体规范化:")
	fmt.Printf("  p1 == p2: %v (相同值应该共享)\n", p1 == p2) // true
	fmt.Printf("  p1 == p3: %v (不同值)\n", p1 == p3)         // false
	fmt.Printf("  p1: %+v\n", p1)

	// 3. 内存对比
	var m1, m2 runtime.MemStats

	// 不使用 unique（会产生重复）
	strings1 := make([]string, 10000)
	runtime.GC()
	runtime.ReadMemStats(&m1)
	for i := range strings1 {
		strings1[i] = "repeated string content" // 每个都是新分配
	}
	runtime.GC()
	runtime.ReadMemStats(&m2)
	alloc1 := m2.Alloc - m1.Alloc

	// 使用 unique（共享相同值）
	strings2 := make([]string, 10000)
	runtime.GC()
	runtime.ReadMemStats(&m1)
	for i := range strings2 {
		strings2[i] = MakeUnique("repeated string content") // 共享相同值
	}
	runtime.GC()
	runtime.ReadMemStats(&m2)
	alloc2 := m2.Alloc - m1.Alloc

	fmt.Println("\n3. 内存占用对比:")
	fmt.Printf("  不使用 unique: %d KB\n", alloc1/1024)
	fmt.Printf("  使用 unique: %d KB\n", alloc2/1024)
	fmt.Printf("  节省: %.1f%%\n", float64(alloc1-alloc2)/float64(alloc1)*100)

	fmt.Println("\n✅ unique包示例完成")
	fmt.Println("\n💡 提示: unique 包可以自动去重相同的内容，")
	fmt.Println("   在大量重复值的场景下可以显著节省内存。")
}
