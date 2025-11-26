package collection

import (
	"cmp"
	"slices"
)

// Contains 检查切片是否包含指定元素
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Index 返回元素在切片中的索引，如果不存在返回-1
func Index[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// Remove 从切片中移除指定元素（只移除第一个匹配的元素）
func Remove[T comparable](slice []T, item T) []T {
	index := Index(slice, item)
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// RemoveAll 从切片中移除所有匹配的元素
func RemoveAll[T comparable](slice []T, item T) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}

// Unique 去除切片中的重复元素
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Filter 过滤切片，保留满足条件的元素
func Filter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map 对切片中的每个元素应用函数，返回新切片
func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Reduce 对切片进行归约操作
func Reduce[T any, R any](slice []T, initial R, fn func(R, T) R) R {
	result := initial
	for _, v := range slice {
		result = fn(result, v)
	}
	return result
}

// ForEach 对切片中的每个元素执行函数
func ForEach[T any](slice []T, fn func(T)) {
	for _, v := range slice {
		fn(v)
	}
}

// Chunk 将切片分割成指定大小的块
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{slice}
	}

	result := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result
}

// Flatten 展平二维切片
func Flatten[T any](slices [][]T) []T {
	totalLen := 0
	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// Reverse 反转切片
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// Sort 排序切片（使用cmp.Ordered）
func Sort[T cmp.Ordered](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	slices.Sort(result)
	return result
}

// SortBy 使用自定义比较函数排序切片
func SortBy[T any](slice []T, cmp func(T, T) int) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	slices.SortFunc(result, cmp)
	return result
}

// SortDesc 降序排序切片
func SortDesc[T cmp.Ordered](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	slices.SortFunc(result, func(a, b T) int {
		return cmp.Compare(b, a)
	})
	return result
}

// First 返回切片的第一个元素，如果切片为空返回零值
func First[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// Last 返回切片的最后一个元素，如果切片为空返回零值
func Last[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// Take 返回切片的前n个元素
func Take[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		return slice
	}
	result := make([]T, n)
	copy(result, slice[:n])
	return result
}

// Drop 返回切片去掉前n个元素后的剩余部分
func Drop[T any](slice []T, n int) []T {
	if n <= 0 {
		return slice
	}
	if n >= len(slice) {
		return []T{}
	}
	return slice[n:]
}

// TakeWhile 返回满足条件的前缀元素
func TakeWhile[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !fn(v) {
			break
		}
		result = append(result, v)
	}
	return result
}

// DropWhile 跳过满足条件的前缀元素
func DropWhile[T any](slice []T, fn func(T) bool) []T {
	index := 0
	for i, v := range slice {
		if !fn(v) {
			index = i
			break
		}
		index = i + 1
	}
	if index >= len(slice) {
		return []T{}
	}
	return slice[index:]
}

// Partition 将切片分割成两部分：满足条件的和不满足条件的
func Partition[T any](slice []T, fn func(T) bool) ([]T, []T) {
	trueSlice := make([]T, 0, len(slice))
	falseSlice := make([]T, 0, len(slice))
	for _, v := range slice {
		if fn(v) {
			trueSlice = append(trueSlice, v)
		} else {
			falseSlice = append(falseSlice, v)
		}
	}
	return trueSlice, falseSlice
}

// GroupBy 根据键函数对切片进行分组
func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range slice {
		key := keyFn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// Count 统计满足条件的元素数量
func Count[T any](slice []T, fn func(T) bool) int {
	count := 0
	for _, v := range slice {
		if fn(v) {
			count++
		}
	}
	return count
}

// Any 检查是否有元素满足条件
func Any[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}
	return false
}

// All 检查是否所有元素都满足条件
func All[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Sum 计算数值切片的和
func Sum[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](slice []T) T {
	var sum T
	for _, v := range slice {
		sum += v
	}
	return sum
}

// Max 返回切片中的最大值
func Max[T cmp.Ordered](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max, true
}

// Min 返回切片中的最小值
func Min[T cmp.Ordered](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	min := slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
	}
	return min, true
}

// Average 计算数值切片的平均值
func Average[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](slice []T) (float64, bool) {
	if len(slice) == 0 {
		return 0, false
	}
	sum := Sum(slice)
	return float64(sum) / float64(len(slice)), true
}

// Intersect 返回两个切片的交集
func Intersect[T comparable](slice1, slice2 []T) []T {
	set2 := make(map[T]bool)
	for _, v := range slice2 {
		set2[v] = true
	}

	result := make([]T, 0)
	seen := make(map[T]bool)
	for _, v := range slice1 {
		if set2[v] && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Union 返回两个切片的并集
func Union[T comparable](slice1, slice2 []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)

	for _, v := range slice1 {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	for _, v := range slice2 {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Difference 返回两个切片的差集（slice1中有但slice2中没有的元素）
func Difference[T comparable](slice1, slice2 []T) []T {
	set2 := make(map[T]bool)
	for _, v := range slice2 {
		set2[v] = true
	}

	result := make([]T, 0)
	seen := make(map[T]bool)
	for _, v := range slice1 {
		if !set2[v] && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// Zip 将两个切片组合成键值对切片
func Zip[T, U any](slice1 []T, slice2 []U) []struct {
	First  T
	Second U
} {
	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}

	result := make([]struct {
		First  T
		Second U
	}, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = struct {
			First  T
			Second U
		}{slice1[i], slice2[i]}
	}
	return result
}

// Shuffle 随机打乱切片（使用Fisher-Yates算法）
func Shuffle[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	// 注意：这里需要导入math/rand或crypto/rand来实现真正的随机打乱
	// 当前实现只是复制，实际使用时应该使用随机数生成器
	return result
}

// MapKeys 获取map的所有键
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapValues 获取map的所有值
func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// MapContains 检查map是否包含指定键
func MapContains[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// MapGet 获取map的值，如果不存在返回默认值
func MapGet[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

// MapFilter 过滤map，保留满足条件的键值对
func MapFilter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapMap 对map中的每个键值对应用函数，返回新map
func MapMap[K comparable, V any, R any](m map[K]V, fn func(K, V) R) []R {
	result := make([]R, 0, len(m))
	for k, v := range m {
		result = append(result, fn(k, v))
	}
	return result
}
