// Go 1.25 testing.B.Loop 示例
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

// 传统方式
func BenchmarkTraditional(b *testing.B) {
	for i := 0; i < b.N; i++ {
		expensiveOperation()
	}
}

// Loop方式 (Go 1.25+)
func BenchmarkLoop(b *testing.B) {
	for b.Loop() {
		expensiveOperation()
	}
}

// 带内存分配报告
func BenchmarkLoopWithAllocs(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		data := make([]byte, 1024)
		_ = data
	}
}

