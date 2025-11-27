package testing_enhancements

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	mathrand "math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// 注意：这个示例已更新为Go 1.25.3兼容版本
// Go 1.25中移除了实验性的testing.Loop API
// 我们使用传统的基准测试模式，但展示了最佳实践

// generateData 是一个模拟昂贵设置的函数。
func generateData() []byte {
	time.Sleep(50 * time.Millisecond) // 模拟耗时操作
	return []byte("this is some heavy data to process")
}

// processData 是我们要进行基准测试的核心函数。
func processData(data []byte) {
	_ = sha256.Sum256(data)
}

// --- 示例1: 传统 for b.N 写法 ---
// 缺点：如果设置代码（generateData）被误放入循环，会严重影响测试准确性。
func BenchmarkProcessData_Traditional(b *testing.B) {
	data := generateData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processData(data)
	}
}

// --- 示例2: 改进的基准测试写法 ---
// 优点：结构更清晰，设置代码和循环体分开。
// 注意：Go 1.25中testing.Loop API已被移除，我们使用传统方式
func BenchmarkProcessData_Improved(b *testing.B) {
	data := generateData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processData(data)
	}
}

// --- 示例3: 每次迭代都需要设置 (Per-iteration Setup) ---
// 场景：假设我们需要测试一个函数，它会修改传入的 buffer。
func processAndModify(buf *bytes.Buffer) {
	buf.WriteString("hello")
	// ... 其他操作
}

// 传统写法的困境：不得不在循环内创建对象，这会影响计时。
func BenchmarkModify_Traditional_Incorrect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 错误！对象创建的时间被计入了基准测试。
		buf := new(bytes.Buffer)
		processAndModify(buf)
	}
}

// 改进写法：手动处理Setup，但需要注意时间计入问题
// 注意：在Go 1.25中，如果需要每次迭代都重置对象，
// 可以使用sync.Pool或在循环外创建并重置
func BenchmarkModify_WithSetup(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf := new(bytes.Buffer)
		b.StartTimer()

		processAndModify(buf)
	}
}

// --- 示例4: B.Loop 与并行测试 (RunParallel) ---
var (
	pool = sync.Pool{
		New: func() any {
			// 创建一个可复用的、有一定容量的 builder
			return &strings.Builder{}
		},
	}
	workload = 100
)

// 并行测试的核心逻辑
func parallelWork(id int) {
	builder := pool.Get().(*strings.Builder)
	defer pool.Put(builder)
	builder.Reset()

	for i := 0; i < workload; i++ {
		builder.WriteString(strconv.Itoa(id + i))
	}
	_ = builder.String()
}

// 传统并行测试写法
func BenchmarkParallel_Traditional(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			parallelWork(mathrand.Intn(1000000))
		}
	})
}

// 改进的并行测试写法
func BenchmarkParallel_Improved(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			parallelWork(i)
			i++
		}
	})
}

// 改进的并行测试（带每次迭代设置）
func BenchmarkParallel_WithSetup(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 每次迭代生成新的worker ID
			workerID := mathrand.Intn(1000000)
			parallelWork(workerID)
		}
	})
}

// --- 示例5: 内存分配跟踪 ---

// 使用 ReportAllocs 来报告内存分配
func BenchmarkWithAllocTracking(b *testing.B) {
	b.ReportAllocs() // 报告内存分配情况

	data := generateData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 创建一个新的切片（会产生分配）
		result := make([]byte, len(data))
		copy(result, data)
		processData(result)
	}
}

// 使用 SetBytes 来设置每次操作处理的字节数
func BenchmarkWithBytesTracking(b *testing.B) {
	data := generateData()
	b.SetBytes(int64(len(data))) // 设置每次操作处理的字节数
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processData(data)
	}
}

// --- 示例6: 子基准测试 (Sub-benchmarks) ---

// 使用 Run 方法创建子基准测试，可以测试不同大小的输入
func BenchmarkProcessData_SubBenchmarks(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		size := size // 避免闭包问题
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			data := make([]byte, size)
			rand.Read(data)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				processData(data)
			}
		})
	}
}

// --- 示例7: 使用 sync.Pool 优化对象分配 ---

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func processWithBuffer(data []byte) {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

	buf.Reset()
	buf.Write(data)
	_ = buf.Bytes()
}

// 不使用 Pool 的版本（会产生大量分配）
func BenchmarkProcessWithoutPool(b *testing.B) {
	b.ReportAllocs()
	data := generateData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		buf.Write(data)
		_ = buf.Bytes()
	}
}

// 使用 Pool 的版本（减少分配）
func BenchmarkProcessWithPool(b *testing.B) {
	b.ReportAllocs()
	data := generateData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processWithBuffer(data)
	}
}

// --- 示例8: 基准测试的清理工作 ---

func BenchmarkWithCleanup(b *testing.B) {
	// 设置阶段
	data := generateData()
	b.ResetTimer()

	// 基准测试阶段
	for i := 0; i < b.N; i++ {
		processData(data)
	}

	// 清理阶段（不计入测试时间）
	b.StopTimer()
	// 执行清理工作
	_ = data
	b.StartTimer()
}

// --- 示例9: 比较不同实现的性能 ---

// 实现1: 使用字符串拼接
func concatStrings(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s
	}
	return result
}

// 实现2: 使用 strings.Builder
func concatWithBuilder(strs []string) string {
	var builder strings.Builder
	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()
}

func BenchmarkConcat_StringConcatenation(b *testing.B) {
	strs := []string{"hello", "world", "golang", "testing"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = concatStrings(strs)
	}
}

func BenchmarkConcat_StringBuilder(b *testing.B) {
	strs := []string{"hello", "world", "golang", "testing"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = concatWithBuilder(strs)
	}
}
