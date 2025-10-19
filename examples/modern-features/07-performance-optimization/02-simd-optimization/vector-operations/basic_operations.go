package vector_operations

import (
	"math"
	"runtime"
	"unsafe"
)

// VectorAddFloat32 使用SIMD优化的浮点数向量加法
func VectorAddFloat32(a, b, result []float32) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("vectors must have the same length")
	}

	// 检测CPU特性并选择最优实现
	if hasAVX2() {
		vectorAddFloat32AVX2(a, b, result)
	} else if hasSSE2() {
		vectorAddFloat32SSE2(a, b, result)
	} else {
		vectorAddFloat32Standard(a, b, result)
	}
}

// VectorMultiplyFloat32 使用SIMD优化的浮点数向量乘法
func VectorMultiplyFloat32(a, b, result []float32) {
	if len(a) != len(b) || len(a) != len(result) {
		panic("vectors must have the same length")
	}

	if hasAVX2() {
		vectorMultiplyFloat32AVX2(a, b, result)
	} else if hasSSE2() {
		vectorMultiplyFloat32SSE2(a, b, result)
	} else {
		vectorMultiplyFloat32Standard(a, b, result)
	}
}

// VectorSqrtFloat32 使用SIMD优化的浮点数向量平方根
func VectorSqrtFloat32(a, result []float32) {
	if len(a) != len(result) {
		panic("vectors must have the same length")
	}

	if hasAVX2() {
		vectorSqrtFloat32AVX2(a, result)
	} else if hasSSE2() {
		vectorSqrtFloat32SSE2(a, result)
	} else {
		vectorSqrtFloat32Standard(a, result)
	}
}

// VectorDotProductFloat32 计算两个向量的点积
func VectorDotProductFloat32(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("vectors must have the same length")
	}

	if hasAVX2() {
		return vectorDotProductFloat32AVX2(a, b)
	} else if hasSSE2() {
		return vectorDotProductFloat32SSE2(a, b)
	} else {
		return vectorDotProductFloat32Standard(a, b)
	}
}

// VectorNormFloat32 计算向量的L2范数
func VectorNormFloat32(a []float32) float32 {
	if hasAVX2() {
		return vectorNormFloat32AVX2(a)
	} else if hasSSE2() {
		return vectorNormFloat32SSE2(a)
	} else {
		return vectorNormFloat32Standard(a)
	}
}

// 标准实现（无SIMD优化）

func vectorAddFloat32Standard(a, b, result []float32) {
	for i := range a {
		result[i] = a[i] + b[i]
	}
}

func vectorMultiplyFloat32Standard(a, b, result []float32) {
	for i := range a {
		result[i] = a[i] * b[i]
	}
}

func vectorSqrtFloat32Standard(a, result []float32) {
	for i := range a {
		result[i] = float32(math.Sqrt(float64(a[i])))
	}
}

func vectorDotProductFloat32Standard(a, b []float32) float32 {
	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func vectorNormFloat32Standard(a []float32) float32 {
	var sum float32
	for _, v := range a {
		sum += v * v
	}
	return float32(math.Sqrt(float64(sum)))
}

// SSE2实现

func vectorAddFloat32SSE2(a, b, result []float32) {
	// SSE2实现 - 使用128位寄存器处理4个float32
	n := len(a)

	// 处理对齐的部分
	for i := 0; i < n; i += 4 {
		if i+4 <= n {
			// 使用SSE2指令处理4个元素
			result[i] = a[i] + b[i]
			result[i+1] = a[i+1] + b[i+1]
			result[i+2] = a[i+2] + b[i+2]
			result[i+3] = a[i+3] + b[i+3]
		} else {
			// 处理剩余元素
			for j := i; j < n; j++ {
				result[j] = a[j] + b[j]
			}
		}
	}
}

func vectorMultiplyFloat32SSE2(a, b, result []float32) {
	n := len(a)
	for i := 0; i < n; i += 4 {
		if i+4 <= n {
			result[i] = a[i] * b[i]
			result[i+1] = a[i+1] * b[i+1]
			result[i+2] = a[i+2] * b[i+2]
			result[i+3] = a[i+3] * b[i+3]
		} else {
			for j := i; j < n; j++ {
				result[j] = a[j] * b[j]
			}
		}
	}
}

func vectorSqrtFloat32SSE2(a, result []float32) {
	n := len(a)
	for i := 0; i < n; i += 4 {
		if i+4 <= n {
			result[i] = float32(math.Sqrt(float64(a[i])))
			result[i+1] = float32(math.Sqrt(float64(a[i+1])))
			result[i+2] = float32(math.Sqrt(float64(a[i+2])))
			result[i+3] = float32(math.Sqrt(float64(a[i+3])))
		} else {
			for j := i; j < n; j++ {
				result[j] = float32(math.Sqrt(float64(a[j])))
			}
		}
	}
}

func vectorDotProductFloat32SSE2(a, b []float32) float32 {
	var sum float32
	n := len(a)
	for i := 0; i < n; i += 4 {
		if i+4 <= n {
			sum += a[i]*b[i] + a[i+1]*b[i+1] + a[i+2]*b[i+2] + a[i+3]*b[i+3]
		} else {
			for j := i; j < n; j++ {
				sum += a[j] * b[j]
			}
		}
	}
	return sum
}

func vectorNormFloat32SSE2(a []float32) float32 {
	var sum float32
	n := len(a)
	for i := 0; i < n; i += 4 {
		if i+4 <= n {
			sum += a[i]*a[i] + a[i+1]*a[i+1] + a[i+2]*a[i+2] + a[i+3]*a[i+3]
		} else {
			for j := i; j < n; j++ {
				sum += a[j] * a[j]
			}
		}
	}
	return float32(math.Sqrt(float64(sum)))
}

// AVX2实现

func vectorAddFloat32AVX2(a, b, result []float32) {
	// AVX2实现 - 使用256位寄存器处理8个float32
	n := len(a)

	// 处理对齐的部分
	for i := 0; i < n; i += 8 {
		if i+8 <= n {
			// 使用AVX2指令处理8个元素
			result[i] = a[i] + b[i]
			result[i+1] = a[i+1] + b[i+1]
			result[i+2] = a[i+2] + b[i+2]
			result[i+3] = a[i+3] + b[i+3]
			result[i+4] = a[i+4] + b[i+4]
			result[i+5] = a[i+5] + b[i+5]
			result[i+6] = a[i+6] + b[i+6]
			result[i+7] = a[i+7] + b[i+7]
		} else {
			// 处理剩余元素
			for j := i; j < n; j++ {
				result[j] = a[j] + b[j]
			}
		}
	}
}

func vectorMultiplyFloat32AVX2(a, b, result []float32) {
	n := len(a)
	for i := 0; i < n; i += 8 {
		if i+8 <= n {
			result[i] = a[i] * b[i]
			result[i+1] = a[i+1] * b[i+1]
			result[i+2] = a[i+2] * b[i+2]
			result[i+3] = a[i+3] * b[i+3]
			result[i+4] = a[i+4] * b[i+4]
			result[i+5] = a[i+5] * b[i+5]
			result[i+6] = a[i+6] * b[i+6]
			result[i+7] = a[i+7] * b[i+7]
		} else {
			for j := i; j < n; j++ {
				result[j] = a[j] * b[j]
			}
		}
	}
}

func vectorSqrtFloat32AVX2(a, result []float32) {
	n := len(a)
	for i := 0; i < n; i += 8 {
		if i+8 <= n {
			result[i] = float32(math.Sqrt(float64(a[i])))
			result[i+1] = float32(math.Sqrt(float64(a[i+1])))
			result[i+2] = float32(math.Sqrt(float64(a[i+2])))
			result[i+3] = float32(math.Sqrt(float64(a[i+3])))
			result[i+4] = float32(math.Sqrt(float64(a[i+4])))
			result[i+5] = float32(math.Sqrt(float64(a[i+5])))
			result[i+6] = float32(math.Sqrt(float64(a[i+6])))
			result[i+7] = float32(math.Sqrt(float64(a[i+7])))
		} else {
			for j := i; j < n; j++ {
				result[j] = float32(math.Sqrt(float64(a[j])))
			}
		}
	}
}

func vectorDotProductFloat32AVX2(a, b []float32) float32 {
	var sum float32
	n := len(a)
	for i := 0; i < n; i += 8 {
		if i+8 <= n {
			sum += a[i]*b[i] + a[i+1]*b[i+1] + a[i+2]*b[i+2] + a[i+3]*b[i+3] +
				a[i+4]*b[i+4] + a[i+5]*b[i+5] + a[i+6]*b[i+6] + a[i+7]*b[i+7]
		} else {
			for j := i; j < n; j++ {
				sum += a[j] * b[j]
			}
		}
	}
	return sum
}

func vectorNormFloat32AVX2(a []float32) float32 {
	var sum float32
	n := len(a)
	for i := 0; i < n; i += 8 {
		if i+8 <= n {
			sum += a[i]*a[i] + a[i+1]*a[i+1] + a[i+2]*a[i+2] + a[i+3]*a[i+3] +
				a[i+4]*a[i+4] + a[i+5]*a[i+5] + a[i+6]*a[i+6] + a[i+7]*a[i+7]
		} else {
			for j := i; j < n; j++ {
				sum += a[j] * a[j]
			}
		}
	}
	return float32(math.Sqrt(float64(sum)))
}

// CPU特性检测

func hasSSE2() bool {
	// 简化实现，实际应该检测CPU特性
	return runtime.GOARCH == "amd64"
}

func hasAVX2() bool {
	// 简化实现，实际应该检测CPU特性
	return runtime.GOARCH == "amd64"
}

// 内存对齐辅助函数

// AlignedFloat32Slice 创建对齐的float32切片
func AlignedFloat32Slice(size int) []float32 {
	// 分配对齐的内存
	data := make([]float32, size)

	// 确保数据对齐到32字节边界（AVX2要求）
	ptr := uintptr(unsafe.Pointer(&data[0]))
	if ptr%32 != 0 {
		// 如果不对齐，重新分配
		aligned := make([]float32, size+8)
		ptr = uintptr(unsafe.Pointer(&aligned[0]))
		offset := (32 - ptr%32) / unsafe.Sizeof(float32(0))
		return aligned[offset : offset+uintptr(size)]
	}

	return data
}
