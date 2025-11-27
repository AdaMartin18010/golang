// Go 1.25 基准测试示例
// 注意：Go 1.25 中移除了实验性的 testing.B.Loop() API
// 本示例展示传统的基准测试最佳实践
package main

import (
	"strings"
	"testing"
)

func expensiveOperation() {
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString("test")
	}
	_ = sb.String()
}

// 传统方式 - 推荐使用
func BenchmarkTraditional(b *testing.B) {
	for i := 0; i < b.N; i++ {
		expensiveOperation()
	}
}

// 带设置的基准测试
func BenchmarkWithSetup(b *testing.B) {
	// 昂贵的设置操作
	setupData := strings.Repeat("test", 1000)
	b.ResetTimer() // 重置计时器，排除设置时间

	for i := 0; i < b.N; i++ {
		expensiveOperation()
		_ = setupData
	}
}

// 带内存分配报告
func BenchmarkWithAllocs(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		data := make([]byte, 1024)
		_ = data
	}
}

// 并行基准测试
func BenchmarkParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			expensiveOperation()
		}
	})
}
