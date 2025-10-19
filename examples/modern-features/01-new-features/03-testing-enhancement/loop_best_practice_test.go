package testing_enhancements

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
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
			parallelWork(rand.Intn(1000000))
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
			workerID := rand.Intn(1000000)
			parallelWork(workerID)
		}
	})
}
