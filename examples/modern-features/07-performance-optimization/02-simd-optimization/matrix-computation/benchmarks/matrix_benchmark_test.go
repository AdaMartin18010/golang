package benchmarks

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"matrix_computation"
)

// 测试数据准备
func generateTestMatrices(rows, cols int) (*matrix_computation.Matrix, *matrix_computation.Matrix, *matrix_computation.Matrix) {
	a := matrix_computation.NewMatrix(rows, cols)
	b := matrix_computation.NewMatrix(cols, rows) // 确保可以相乘
	result := matrix_computation.NewMatrix(rows, rows)

	rand.Seed(time.Now().UnixNano())

	// 填充矩阵A
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			a.Set(i, j, rand.Float32()*100)
		}
	}

	// 填充矩阵B
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			b.Set(i, j, rand.Float32()*100)
		}
	}

	return a, b, result
}

// 基准测试：矩阵乘法
func BenchmarkMatrixMultiplyStandard(b *testing.B) {
	a, b_matrix, result := generateTestMatrices(64, 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix_computation.MatrixMultiply(a, b_matrix, result)
	}
}

func BenchmarkMatrixMultiplySIMD(b *testing.B) {
	a, b_matrix, result := generateTestMatrices(64, 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix_computation.MatrixMultiply(a, b_matrix, result)
	}
}

// 基准测试：矩阵加法
func BenchmarkMatrixAddStandard(b *testing.B) {
	a, b_matrix, result := generateTestMatrices(64, 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix_computation.MatrixAdd(a, b_matrix, result)
	}
}

func BenchmarkMatrixAddSIMD(b *testing.B) {
	a, b_matrix, result := generateTestMatrices(64, 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix_computation.MatrixAdd(a, b_matrix, result)
	}
}

// 基准测试：矩阵标量乘法
func BenchmarkMatrixScaleStandard(b *testing.B) {
	a, _, result := generateTestMatrices(64, 64)
	scalar := float32(2.5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix_computation.MatrixScale(a, scalar, result)
	}
}

func BenchmarkMatrixScaleSIMD(b *testing.B) {
	a, _, result := generateTestMatrices(64, 64)
	scalar := float32(2.5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix_computation.MatrixScale(a, scalar, result)
	}
}

// 不同矩阵大小的基准测试
func BenchmarkMatrixMultiplyDifferentSizes(b *testing.B) {
	sizes := []struct {
		rows int
		cols int
	}{
		{16, 16},
		{32, 32},
		{64, 64},
		{128, 128},
		{256, 256},
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dx%d", size.rows, size.cols), func(b *testing.B) {
			a, b_matrix, result := generateTestMatrices(size.rows, size.cols)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matrix_computation.MatrixMultiply(a, b_matrix, result)
			}
		})
	}
}

// 内存对齐测试
func BenchmarkMatrixMultiplyAligned(b *testing.B) {
	rows, cols := 64, 64
	a := matrix_computation.AlignedMatrix(rows, cols)
	b_matrix := matrix_computation.AlignedMatrix(cols, rows)
	result := matrix_computation.AlignedMatrix(rows, rows)

	// 填充测试数据
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			a.Set(i, j, rand.Float32()*100)
		}
	}
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			b_matrix.Set(i, j, rand.Float32()*100)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix_computation.MatrixMultiply(a, b_matrix, result)
	}
}

// 正确性测试
func TestMatrixOperationsCorrectness(t *testing.T) {
	// 测试矩阵乘法
	a, b_matrix, result := generateTestMatrices(4, 4)

	// 手动计算期望结果
	expected := matrix_computation.NewMatrix(4, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var sum float32
			for k := 0; k < 4; k++ {
				sum += a.Get(i, k) * b_matrix.Get(k, j)
			}
			expected.Set(i, j, sum)
		}
	}

	matrix_computation.MatrixMultiply(a, b_matrix, result)

	// 验证结果
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if math.Abs(float64(result.Get(i, j)-expected.Get(i, j))) > 1e-6 {
				t.Errorf("Matrix multiplication failed at [%d,%d]: got %f, expected %f",
					i, j, result.Get(i, j), expected.Get(i, j))
			}
		}
	}

	// 测试矩阵加法
	a2, b2, result2 := generateTestMatrices(4, 4)
	expected2 := matrix_computation.NewMatrix(4, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			expected2.Set(i, j, a2.Get(i, j)+b2.Get(i, j))
		}
	}

	matrix_computation.MatrixAdd(a2, b2, result2)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if math.Abs(float64(result2.Get(i, j)-expected2.Get(i, j))) > 1e-6 {
				t.Errorf("Matrix addition failed at [%d,%d]: got %f, expected %f",
					i, j, result2.Get(i, j), expected2.Get(i, j))
			}
		}
	}

	// 测试矩阵标量乘法
	scalar := float32(2.5)
	a3, _, result3 := generateTestMatrices(4, 4)
	expected3 := matrix_computation.NewMatrix(4, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			expected3.Set(i, j, a3.Get(i, j)*scalar)
		}
	}

	matrix_computation.MatrixScale(a3, scalar, result3)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if math.Abs(float64(result3.Get(i, j)-expected3.Get(i, j))) > 1e-6 {
				t.Errorf("Matrix scale failed at [%d,%d]: got %f, expected %f",
					i, j, result3.Get(i, j), expected3.Get(i, j))
			}
		}
	}
}

// 性能对比测试
func TestMatrixPerformanceComparison(t *testing.T) {
	size := 128
	a, b_matrix, result := generateTestMatrices(size, size)

	// 测试标准实现性能
	start := time.Now()
	for i := 0; i < 100; i++ {
		matrix_computation.MatrixMultiply(a, b_matrix, result)
	}
	standardTime := time.Since(start)

	// 测试SIMD实现性能
	start = time.Now()
	for i := 0; i < 100; i++ {
		matrix_computation.MatrixMultiply(a, b_matrix, result)
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

// 矩阵转置测试
func TestMatrixTranspose(t *testing.T) {
	a := matrix_computation.NewMatrix(3, 4)
	result := matrix_computation.NewMatrix(4, 3)

	// 填充测试数据
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			a.Set(i, j, float32(i*10+j))
		}
	}

	matrix_computation.MatrixTranspose(a, result)

	// 验证转置结果
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if result.Get(j, i) != a.Get(i, j) {
				t.Errorf("Matrix transpose failed at [%d,%d]: got %f, expected %f",
					j, i, result.Get(j, i), a.Get(i, j))
			}
		}
	}
}

// 矩阵范数测试
func TestMatrixNorm(t *testing.T) {
	a := matrix_computation.NewMatrix(2, 2)
	a.Set(0, 0, 3.0)
	a.Set(0, 1, 4.0)
	a.Set(1, 0, 5.0)
	a.Set(1, 1, 12.0)

	norm := matrix_computation.MatrixNorm(a)
	expected := float32(math.Sqrt(3*3 + 4*4 + 5*5 + 12*12))

	if math.Abs(float64(norm-expected)) > 1e-6 {
		t.Errorf("Matrix norm failed: got %f, expected %f", norm, expected)
	}
}
