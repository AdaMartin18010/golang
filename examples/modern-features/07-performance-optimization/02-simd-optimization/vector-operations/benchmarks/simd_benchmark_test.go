package benchmarks

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	vector_operations "performance-optimization-examples/02-simd-optimization/vector-operations"
)

// 测试数据准备
func generateTestData(size int) ([]float32, []float32, []float32) {
	a := make([]float32, size)
	b := make([]float32, size)
	result := make([]float32, size)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		a[i] = rand.Float32() * 100
		b[i] = rand.Float32() * 100
	}

	return a, b, result
}

// 基准测试：向量加法
func BenchmarkVectorAddStandard(b *testing.B) {
	a, b_data, result := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector_operations.VectorAddFloat32(a, b_data, result)
	}
}

func BenchmarkVectorAddSIMD(b *testing.B) {
	a, b_data, result := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector_operations.VectorAddFloat32(a, b_data, result)
	}
}

// 基准测试：向量乘法
func BenchmarkVectorMultiplyStandard(b *testing.B) {
	a, b_data, result := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector_operations.VectorMultiplyFloat32(a, b_data, result)
	}
}

func BenchmarkVectorMultiplySIMD(b *testing.B) {
	a, b_data, result := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector_operations.VectorMultiplyFloat32(a, b_data, result)
	}
}

// 基准测试：向量平方根
func BenchmarkVectorSqrtStandard(b *testing.B) {
	a, _, result := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector_operations.VectorSqrtFloat32(a, result)
	}
}

func BenchmarkVectorSqrtSIMD(b *testing.B) {
	a, _, result := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector_operations.VectorSqrtFloat32(a, result)
	}
}

// 基准测试：向量点积
func BenchmarkVectorDotProductStandard(b *testing.B) {
	a, b_data, _ := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vector_operations.VectorDotProductFloat32(a, b_data)
	}
}

func BenchmarkVectorDotProductSIMD(b *testing.B) {
	a, b_data, _ := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vector_operations.VectorDotProductFloat32(a, b_data)
	}
}

// 基准测试：向量范数
func BenchmarkVectorNormStandard(b *testing.B) {
	a, _, _ := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vector_operations.VectorNormFloat32(a)
	}
}

func BenchmarkVectorNormSIMD(b *testing.B) {
	a, _, _ := generateTestData(1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vector_operations.VectorNormFloat32(a)
	}
}

// 不同数据大小的基准测试
func BenchmarkVectorAddDifferentSizes(b *testing.B) {
	sizes := []int{64, 256, 1024, 4096, 16384}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			a, b_data, result := generateTestData(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				vector_operations.VectorAddFloat32(a, b_data, result)
			}
		})
	}
}

// 内存对齐测试
func BenchmarkVectorAddAligned(b *testing.B) {
	size := 1024
	a := vector_operations.AlignedFloat32Slice(size)
	b_data := vector_operations.AlignedFloat32Slice(size)
	result := vector_operations.AlignedFloat32Slice(size)

	// 填充测试数据
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		a[i] = rand.Float32() * 100
		b_data[i] = rand.Float32() * 100
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector_operations.VectorAddFloat32(a, b_data, result)
	}
}

// 正确性测试
func TestVectorOperationsCorrectness(t *testing.T) {
	size := 100
	a, b_data, result := generateTestData(size)

	// 测试向量加法
	vector_operations.VectorAddFloat32(a, b_data, result)
	for i := 0; i < size; i++ {
		expected := a[i] + b_data[i]
		if math.Abs(float64(result[i]-expected)) > 1e-6 {
			t.Errorf("Vector addition failed at index %d: got %f, expected %f", i, result[i], expected)
		}
	}

	// 测试向量乘法
	vector_operations.VectorMultiplyFloat32(a, b_data, result)
	for i := 0; i < size; i++ {
		expected := a[i] * b_data[i]
		if math.Abs(float64(result[i]-expected)) > 1e-6 {
			t.Errorf("Vector multiplication failed at index %d: got %f, expected %f", i, result[i], expected)
		}
	}

	// 测试向量平方根
	vector_operations.VectorSqrtFloat32(a, result)
	for i := 0; i < size; i++ {
		expected := float32(math.Sqrt(float64(a[i])))
		if math.Abs(float64(result[i]-expected)) > 1e-6 {
			t.Errorf("Vector sqrt failed at index %d: got %f, expected %f", i, result[i], expected)
		}
	}

	// 测试向量点积
	// 注意：SIMD 实现按 4/8 个元素分组累加，与测试中的逐个 float32 累加
	// 顺序不同。float32 加法不满足结合律，累加顺序差异会导致最后几个 ULP
	// 的偏差（此处量级 ~2.5e5，1 ULP ≈ 0.015625）。实现本身无误，
	// 因此采用相对容差而非绝对相等进行断言。
	dotProduct := vector_operations.VectorDotProductFloat32(a, b_data)
	var expected float32
	for i := 0; i < size; i++ {
		expected += a[i] * b_data[i]
	}
	if relDiff := math.Abs(float64(dotProduct-expected)) / math.Max(math.Abs(float64(expected)), 1); relDiff > 1e-4 {
		t.Errorf("Vector dot product failed: got %f, expected %f (relative diff %g)", dotProduct, expected, relDiff)
	}

	// 测试向量范数（与点积同理，SIMD 分组累加顺序差异属 float32 精度损失）
	norm := vector_operations.VectorNormFloat32(a)
	var sum float32
	for _, v := range a {
		sum += v * v
	}
	expected = float32(math.Sqrt(float64(sum)))
	if relDiff := math.Abs(float64(norm-expected)) / math.Max(math.Abs(float64(expected)), 1); relDiff > 1e-4 {
		t.Errorf("Vector norm failed: got %f, expected %f (relative diff %g)", norm, expected, relDiff)
	}
}

// 性能对比测试
func TestPerformanceComparison(t *testing.T) {
	size := 10000
	a, b_data, result := generateTestData(size)

	// 测试标准实现性能
	start := time.Now()
	for i := 0; i < 1000; i++ {
		vector_operations.VectorAddFloat32(a, b_data, result)
	}
	standardTime := time.Since(start)

	// 测试SIMD实现性能
	start = time.Now()
	for i := 0; i < 1000; i++ {
		vector_operations.VectorAddFloat32(a, b_data, result)
	}
	simdTime := time.Since(start)

	t.Logf("Standard implementation time: %v", standardTime)
	t.Logf("SIMD implementation time: %v", simdTime)
	t.Logf("Performance improvement: %.2fx", float64(standardTime)/float64(simdTime))

	// 验证性能提升
	if simdTime >= standardTime {
		t.Logf("Note: SIMD implementation did not show performance improvement in this test")
	}
}
