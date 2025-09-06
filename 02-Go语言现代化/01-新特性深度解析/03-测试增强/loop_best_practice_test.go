package testing_enhancements

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

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

// --- 示例2: B.Loop 基本用法 ---
// 优点：结构更清晰，强制将设置代码和循环体分开。
func BenchmarkProcessData_Loop_Basic(b *testing.B) {
	data := generateData()
	b.ResetTimer()

	b.Loop(func(i int) {
		processData(data)
	})
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

// B.Loop 正确写法：使用 testing.Loop 结构体，Setup 不会计入测试时间。
func BenchmarkModify_Loop_WithSetup(b *testing.B) {
	var buf *bytes.Buffer

	b.Loop(testing.Loop{
		// Setup 在每次 Body 执行前调用，且其执行时间不计入 b.N。
		Setup: func() {
			buf = new(bytes.Buffer)
		},
		Body: func(i int) {
			processAndModify(buf)
		},
	})
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
			parallelWork(rand.Int())
		}
	})
}

// B.Loop 并行测试写法：代码更统一、简洁。
func BenchmarkParallel_Loop(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		// 注意：这里不需要 pb.Next() 循环了
		b.Loop(func(i int) {
			// i 在并行测试中不是连续的，但可以用来获取唯一性
			parallelWork(i)
		})
	})
}

// B.Loop 并行测试（带每次迭代设置）：
// 虽然不常见，但展示了其组合能力。
func BenchmarkParallel_Loop_WithSetup(b *testing.B) {
	var workerID int

	b.RunParallel(func(pb *testing.PB) {
		b.Loop(testing.Loop{
			// Setup 在每个 goroutine 的每次迭代前执行
			Setup: func() {
				workerID = rand.Int()
			},
			Body: func(i int) {
				parallelWork(workerID)
			},
		})
	})
}
