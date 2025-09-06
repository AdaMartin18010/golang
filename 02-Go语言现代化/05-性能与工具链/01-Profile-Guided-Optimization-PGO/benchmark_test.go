package main

import "testing"

// BenchmarkHotFunction 对我们期望通过 PGO 优化的热点函数进行基准测试。
func BenchmarkHotFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 使用一个固定的、有一定计算量的值进行测试
		hotFunction(20)
	}
}
