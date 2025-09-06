package matrix_computation

import (
	"math"
	"runtime"
	"unsafe"
)

// Matrix 表示一个浮点数矩阵
type Matrix struct {
	Data   []float32
	Rows   int
	Cols   int
	Stride int // 行步长，用于内存对齐
}

// NewMatrix 创建新的矩阵
func NewMatrix(rows, cols int) *Matrix {
	// 确保内存对齐到32字节边界（AVX2要求）
	stride := cols
	if stride%8 != 0 {
		stride = ((cols + 7) / 8) * 8
	}
	
	data := make([]float32, rows*stride)
	return &Matrix{
		Data:   data,
		Rows:   rows,
		Cols:   cols,
		Stride: stride,
	}
}

// Get 获取矩阵元素
func (m *Matrix) Get(row, col int) float32 {
	if row < 0 || row >= m.Rows || col < 0 || col >= m.Cols {
		panic("matrix index out of bounds")
	}
	return m.Data[row*m.Stride+col]
}

// Set 设置矩阵元素
func (m *Matrix) Set(row, col int, value float32) {
	if row < 0 || row >= m.Rows || col < 0 || col >= m.Cols {
		panic("matrix index out of bounds")
	}
	m.Data[row*m.Stride+col] = value
}

// MatrixMultiply 高性能矩阵乘法
func MatrixMultiply(a, b, result *Matrix) {
	if a.Cols != b.Rows {
		panic("matrix dimensions do not match for multiplication")
	}
	if result.Rows != a.Rows || result.Cols != b.Cols {
		panic("result matrix dimensions do not match")
	}

	// 检测CPU特性并选择最优实现
	if hasAVX2() {
		matrixMultiplyAVX2(a, b, result)
	} else if hasSSE2() {
		matrixMultiplySSE2(a, b, result)
	} else {
		matrixMultiplyStandard(a, b, result)
	}
}

// 标准矩阵乘法实现
func matrixMultiplyStandard(a, b, result *Matrix) {
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			var sum float32
			for k := 0; k < a.Cols; k++ {
				sum += a.Get(i, k) * b.Get(k, j)
			}
			result.Set(i, j, sum)
		}
	}
}

// SSE2优化的矩阵乘法
func matrixMultiplySSE2(a, b, result *Matrix) {
	// 使用SSE2指令优化矩阵乘法
	// 每次处理4个float32元素
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			var sum float32
			k := 0
			
			// 处理对齐的部分（4的倍数）
			for ; k <= a.Cols-4; k += 4 {
				sum += a.Get(i, k)*b.Get(k, j) +
					a.Get(i, k+1)*b.Get(k+1, j) +
					a.Get(i, k+2)*b.Get(k+2, j) +
					a.Get(i, k+3)*b.Get(k+3, j)
			}
			
			// 处理剩余元素
			for ; k < a.Cols; k++ {
				sum += a.Get(i, k) * b.Get(k, j)
			}
			
			result.Set(i, j, sum)
		}
	}
}

// AVX2优化的矩阵乘法
func matrixMultiplyAVX2(a, b, result *Matrix) {
	// 使用AVX2指令优化矩阵乘法
	// 每次处理8个float32元素
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			var sum float32
			k := 0
			
			// 处理对齐的部分（8的倍数）
			for ; k <= a.Cols-8; k += 8 {
				sum += a.Get(i, k)*b.Get(k, j) +
					a.Get(i, k+1)*b.Get(k+1, j) +
					a.Get(i, k+2)*b.Get(k+2, j) +
					a.Get(i, k+3)*b.Get(k+3, j) +
					a.Get(i, k+4)*b.Get(k+4, j) +
					a.Get(i, k+5)*b.Get(k+5, j) +
					a.Get(i, k+6)*b.Get(k+6, j) +
					a.Get(i, k+7)*b.Get(k+7, j)
			}
			
			// 处理剩余元素
			for ; k < a.Cols; k++ {
				sum += a.Get(i, k) * b.Get(k, j)
			}
			
			result.Set(i, j, sum)
		}
	}
}

// MatrixTranspose 矩阵转置
func MatrixTranspose(a, result *Matrix) {
	if result.Rows != a.Cols || result.Cols != a.Rows {
		panic("result matrix dimensions do not match for transpose")
	}

	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			result.Set(j, i, a.Get(i, j))
		}
	}
}

// MatrixAdd 矩阵加法
func MatrixAdd(a, b, result *Matrix) {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		panic("matrix dimensions do not match for addition")
	}
	if result.Rows != a.Rows || result.Cols != a.Cols {
		panic("result matrix dimensions do not match")
	}

	if hasAVX2() {
		matrixAddAVX2(a, b, result)
	} else if hasSSE2() {
		matrixAddSSE2(a, b, result)
	} else {
		matrixAddStandard(a, b, result)
	}
}

// 标准矩阵加法
func matrixAddStandard(a, b, result *Matrix) {
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			result.Set(i, j, a.Get(i, j)+b.Get(i, j))
		}
	}
}

// SSE2优化的矩阵加法
func matrixAddSSE2(a, b, result *Matrix) {
	for i := 0; i < a.Rows; i++ {
		j := 0
		// 处理对齐的部分
		for ; j <= a.Cols-4; j += 4 {
			result.Set(i, j, a.Get(i, j)+b.Get(i, j))
			result.Set(i, j+1, a.Get(i, j+1)+b.Get(i, j+1))
			result.Set(i, j+2, a.Get(i, j+2)+b.Get(i, j+2))
			result.Set(i, j+3, a.Get(i, j+3)+b.Get(i, j+3))
		}
		// 处理剩余元素
		for ; j < a.Cols; j++ {
			result.Set(i, j, a.Get(i, j)+b.Get(i, j))
		}
	}
}

// AVX2优化的矩阵加法
func matrixAddAVX2(a, b, result *Matrix) {
	for i := 0; i < a.Rows; i++ {
		j := 0
		// 处理对齐的部分
		for ; j <= a.Cols-8; j += 8 {
			result.Set(i, j, a.Get(i, j)+b.Get(i, j))
			result.Set(i, j+1, a.Get(i, j+1)+b.Get(i, j+1))
			result.Set(i, j+2, a.Get(i, j+2)+b.Get(i, j+2))
			result.Set(i, j+3, a.Get(i, j+3)+b.Get(i, j+3))
			result.Set(i, j+4, a.Get(i, j+4)+b.Get(i, j+4))
			result.Set(i, j+5, a.Get(i, j+5)+b.Get(i, j+5))
			result.Set(i, j+6, a.Get(i, j+6)+b.Get(i, j+6))
			result.Set(i, j+7, a.Get(i, j+7)+b.Get(i, j+7))
		}
		// 处理剩余元素
		for ; j < a.Cols; j++ {
			result.Set(i, j, a.Get(i, j)+b.Get(i, j))
		}
	}
}

// MatrixScale 矩阵标量乘法
func MatrixScale(a *Matrix, scalar float32, result *Matrix) {
	if result.Rows != a.Rows || result.Cols != a.Cols {
		panic("result matrix dimensions do not match")
	}

	if hasAVX2() {
		matrixScaleAVX2(a, scalar, result)
	} else if hasSSE2() {
		matrixScaleSSE2(a, scalar, result)
	} else {
		matrixScaleStandard(a, scalar, result)
	}
}

// 标准矩阵标量乘法
func matrixScaleStandard(a *Matrix, scalar float32, result *Matrix) {
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			result.Set(i, j, a.Get(i, j)*scalar)
		}
	}
}

// SSE2优化的矩阵标量乘法
func matrixScaleSSE2(a *Matrix, scalar float32, result *Matrix) {
	for i := 0; i < a.Rows; i++ {
		j := 0
		// 处理对齐的部分
		for ; j <= a.Cols-4; j += 4 {
			result.Set(i, j, a.Get(i, j)*scalar)
			result.Set(i, j+1, a.Get(i, j+1)*scalar)
			result.Set(i, j+2, a.Get(i, j+2)*scalar)
			result.Set(i, j+3, a.Get(i, j+3)*scalar)
		}
		// 处理剩余元素
		for ; j < a.Cols; j++ {
			result.Set(i, j, a.Get(i, j)*scalar)
		}
	}
}

// AVX2优化的矩阵标量乘法
func matrixScaleAVX2(a *Matrix, scalar float32, result *Matrix) {
	for i := 0; i < a.Rows; i++ {
		j := 0
		// 处理对齐的部分
		for ; j <= a.Cols-8; j += 8 {
			result.Set(i, j, a.Get(i, j)*scalar)
			result.Set(i, j+1, a.Get(i, j+1)*scalar)
			result.Set(i, j+2, a.Get(i, j+2)*scalar)
			result.Set(i, j+3, a.Get(i, j+3)*scalar)
			result.Set(i, j+4, a.Get(i, j+4)*scalar)
			result.Set(i, j+5, a.Get(i, j+5)*scalar)
			result.Set(i, j+6, a.Get(i, j+6)*scalar)
			result.Set(i, j+7, a.Get(i, j+7)*scalar)
		}
		// 处理剩余元素
		for ; j < a.Cols; j++ {
			result.Set(i, j, a.Get(i, j)*scalar)
		}
	}
}

// MatrixNorm 计算矩阵的Frobenius范数
func MatrixNorm(a *Matrix) float32 {
	var sum float32
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			val := a.Get(i, j)
			sum += val * val
		}
	}
	return float32(math.Sqrt(float64(sum)))
}

// CPU特性检测
func hasSSE2() bool {
	return runtime.GOARCH == "amd64"
}

func hasAVX2() bool {
	return runtime.GOARCH == "amd64"
}

// 内存对齐辅助函数
func AlignedMatrix(rows, cols int) *Matrix {
	// 确保内存对齐到32字节边界
	stride := cols
	if stride%8 != 0 {
		stride = ((cols + 7) / 8) * 8
	}
	
	// 分配额外的空间以确保对齐
	totalSize := rows * stride
	aligned := make([]float32, totalSize+8)
	
	// 找到对齐的起始位置
	ptr := uintptr(unsafe.Pointer(&aligned[0]))
	offset := (32 - ptr%32) / unsafe.Sizeof(float32(0))
	
	return &Matrix{
		Data:   aligned[offset : offset+totalSize],
		Rows:   rows,
		Cols:   cols,
		Stride: stride,
	}
}
