package sort

import (
	"sort"
)

// Ints 对整数切片进行排序
func Ints(nums []int) {
	sort.Ints(nums)
}

// IntsAreSorted 检查整数切片是否已排序
func IntsAreSorted(nums []int) bool {
	return sort.IntsAreSorted(nums)
}

// SearchInts 在已排序的整数切片中搜索
func SearchInts(nums []int, x int) int {
	return sort.SearchInts(nums, x)
}

// Float64s 对float64切片进行排序
func Float64s(nums []float64) {
	sort.Float64s(nums)
}

// Float64sAreSorted 检查float64切片是否已排序
func Float64sAreSorted(nums []float64) bool {
	return sort.Float64sAreSorted(nums)
}

// SearchFloat64s 在已排序的float64切片中搜索
func SearchFloat64s(nums []float64, x float64) int {
	return sort.SearchFloat64s(nums, x)
}

// Strings 对字符串切片进行排序
func Strings(strs []string) {
	sort.Strings(strs)
}

// StringsAreSorted 检查字符串切片是否已排序
func StringsAreSorted(strs []string) bool {
	return sort.StringsAreSorted(strs)
}

// SearchStrings 在已排序的字符串切片中搜索
func SearchStrings(strs []string, x string) int {
	return sort.SearchStrings(strs, x)
}

// IntsReverse 对整数切片进行反向排序
func IntsReverse(nums []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
}

// Float64sReverse 对float64切片进行反向排序
func Float64sReverse(nums []float64) {
	sort.Sort(sort.Reverse(sort.Float64Slice(nums)))
}

// StringsReverse 对字符串切片进行反向排序
func StringsReverse(strs []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(strs)))
}

// SortBy 根据比较函数对切片进行排序
func SortBy[T any](slice []T, less func(i, j int) bool) {
	sort.Slice(slice, less)
}

// SortByFunc 根据比较函数对切片进行排序（使用元素比较）
func SortByFunc[T any](slice []T, cmp func(a, b T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return cmp(slice[i], slice[j])
	})
}

// SortStable 稳定排序
func SortStable[T any](slice []T, less func(i, j int) bool) {
	sort.SliceStable(slice, less)
}

// SortStableByFunc 稳定排序（使用元素比较）
func SortStableByFunc[T any](slice []T, cmp func(a, b T) bool) {
	sort.SliceStable(slice, func(i, j int) bool {
		return cmp(slice[i], slice[j])
	})
}

// Reverse 反转切片
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// IsSorted 检查切片是否已排序
func IsSorted[T any](slice []T, less func(i, j int) bool) bool {
	return sort.SliceIsSorted(slice, less)
}

// IsSortedFunc 检查切片是否已排序（使用元素比较）
func IsSortedFunc[T any](slice []T, cmp func(a, b T) bool) bool {
	return sort.SliceIsSorted(slice, func(i, j int) bool {
		return cmp(slice[i], slice[j])
	})
}

// Search 在已排序的切片中搜索
func Search[T any](n int, f func(int) bool) int {
	return sort.Search(n, f)
}

// SearchSlice 在已排序的切片中搜索元素
func SearchSlice[T comparable](slice []T, x T, less func(a, b T) bool) int {
	return sort.Search(len(slice), func(i int) bool {
		return !less(slice[i], x)
	})
}

// Unique 去重并排序
func Unique[T comparable](slice []T, less func(a, b T) bool) []T {
	if len(slice) == 0 {
		return slice
	}

	// 创建映射去重
	seen := make(map[T]bool)
	unique := make([]T, 0, len(slice))
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}

	// 排序
	SortByFunc(unique, less)
	return unique
}

// UniqueInts 去重并排序整数切片
func UniqueInts(nums []int) []int {
	return Unique(nums, func(a, b int) bool {
		return a < b
	})
}

// UniqueFloat64s 去重并排序float64切片
func UniqueFloat64s(nums []float64) []float64 {
	return Unique(nums, func(a, b float64) bool {
		return a < b
	})
}

// UniqueStrings 去重并排序字符串切片
func UniqueStrings(strs []string) []string {
	return Unique(strs, func(a, b string) bool {
		return a < b
	})
}

// TopN 返回前N个最大元素
func TopN[T any](slice []T, n int, less func(a, b T) bool) []T {
	if n <= 0 || len(slice) == 0 {
		return []T{}
	}
	if n >= len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		SortByFunc(result, less)
		return result
	}

	// 创建副本并排序
	result := make([]T, len(slice))
	copy(result, slice)
	SortByFunc(result, less)

	// 返回前N个
	return result[:n]
}

// BottomN 返回前N个最小元素
func BottomN[T any](slice []T, n int, less func(a, b T) bool) []T {
	if n <= 0 || len(slice) == 0 {
		return []T{}
	}
	if n >= len(slice) {
		result := make([]T, len(slice))
		copy(result, slice)
		SortByFunc(result, func(a, b T) bool {
			return !less(a, b)
		})
		return result
	}

	// 创建副本并排序
	result := make([]T, len(slice))
	copy(result, slice)
	SortByFunc(result, func(a, b T) bool {
		return !less(a, b)
	})

	// 返回前N个
	return result[:n]
}

// TopNInts 返回前N个最大整数
func TopNInts(nums []int, n int) []int {
	return TopN(nums, n, func(a, b int) bool {
		return a > b
	})
}

// BottomNInts 返回前N个最小整数
func BottomNInts(nums []int, n int) []int {
	return BottomN(nums, n, func(a, b int) bool {
		return a < b
	})
}

// TopNFloat64s 返回前N个最大float64
func TopNFloat64s(nums []float64, n int) []float64 {
	return TopN(nums, n, func(a, b float64) bool {
		return a > b
	})
}

// BottomNFloat64s 返回前N个最小float64
func BottomNFloat64s(nums []float64, n int) []float64 {
	return BottomN(nums, n, func(a, b float64) bool {
		return a < b
	})
}

// Shuffle 随机打乱切片
func Shuffle[T any](slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := int(float64(i+1) * float64(i) / float64(i+1))
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// ShuffleWithSeed 使用种子随机打乱切片
func ShuffleWithSeed[T any](slice []T, seed int64) {
	// 简单的线性同余生成器
	rng := seed
	for i := len(slice) - 1; i > 0; i-- {
		rng = (rng*1103515245 + 12345) & 0x7fffffff
		j := int(rng) % (i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// SortByKey 根据键函数对切片进行排序
func SortByKey[T any, K comparable](slice []T, keyFunc func(T) K, less func(a, b K) bool) {
	SortByFunc(slice, func(a, b T) bool {
		return less(keyFunc(a), keyFunc(b))
	})
}

// SortByKeyInt 根据整数键对切片进行排序
func SortByKeyInt[T any](slice []T, keyFunc func(T) int) {
	SortByKey(slice, keyFunc, func(a, b int) bool {
		return a < b
	})
}

// SortByKeyString 根据字符串键对切片进行排序
func SortByKeyString[T any](slice []T, keyFunc func(T) string) {
	SortByKey(slice, keyFunc, func(a, b string) bool {
		return a < b
	})
}

// SortByKeyFloat64 根据float64键对切片进行排序
func SortByKeyFloat64[T any](slice []T, keyFunc func(T) float64) {
	SortByKey(slice, keyFunc, func(a, b float64) bool {
		return a < b
	})
}

// MultiSort 多字段排序
func MultiSort[T any](slice []T, comparers ...func(a, b T) int) {
	SortByFunc(slice, func(a, b T) bool {
		for _, cmp := range comparers {
			result := cmp(a, b)
			if result != 0 {
				return result < 0
			}
		}
		return false
	})
}

// CompareInt 比较两个整数
func CompareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// CompareFloat64 比较两个float64
func CompareFloat64(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// CompareString 比较两个字符串
func CompareString(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
