package random

import (
	"crypto/rand"
	"math/big"
	"math/rand/v2"
	"time"
	"unsafe"
)

const (
	// 字符集
	Letters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits      = "0123456789"
	LettersNum  = Letters + Digits
	HexChars    = "0123456789abcdef"
	HexCharsUpper = "0123456789ABCDEF"
	SpecialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	AllChars    = LettersNum + SpecialChars
)

var (
	// 默认随机数生成器（使用时间种子）
	defaultRand = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
)

// Int 生成随机整数 [0, max)
func Int(max int) int {
	if max <= 0 {
		return 0
	}
	return defaultRand.IntN(max)
}

// IntRange 生成指定范围内的随机整数 [min, max)
func IntRange(min, max int) int {
	if min >= max {
		return min
	}
	return min + defaultRand.IntN(max-min)
}

// Int64 生成随机64位整数 [0, max)
func Int64(max int64) int64 {
	if max <= 0 {
		return 0
	}
	return defaultRand.Int64N(max)
}

// Int64Range 生成指定范围内的随机64位整数 [min, max)
func Int64Range(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + defaultRand.Int64N(max-min)
}

// Float64 生成随机浮点数 [0.0, 1.0)
func Float64() float64 {
	return defaultRand.Float64()
}

// Float64Range 生成指定范围内的随机浮点数 [min, max)
func Float64Range(min, max float64) float64 {
	if min >= max {
		return min
	}
	return min + defaultRand.Float64()*(max-min)
}

// String 生成指定长度的随机字符串
func String(length int) string {
	return StringWithCharset(length, LettersNum)
}

// StringWithCharset 使用指定字符集生成随机字符串
func StringWithCharset(length int, charset string) string {
	if length <= 0 || len(charset) == 0 {
		return ""
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[defaultRand.IntN(len(charset))]
	}
	return string(b)
}

// Letters 生成随机字母字符串
func LettersString(length int) string {
	return StringWithCharset(length, Letters)
}

// Digits 生成随机数字字符串
func DigitsString(length int) string {
	return StringWithCharset(length, Digits)
}

// Hex 生成随机十六进制字符串
func Hex(length int) string {
	return StringWithCharset(length, HexChars)
}

// HexUpper 生成随机十六进制字符串（大写）
func HexUpper(length int) string {
	return StringWithCharset(length, HexCharsUpper)
}

// UUID 生成随机UUID字符串（简化版，不遵循标准格式）
func UUID() string {
	return Hex(32)
}

// UUIDWithDashes 生成带连字符的UUID字符串（简化版）
func UUIDWithDashes() string {
	uuid := Hex(32)
	return uuid[0:8] + "-" + uuid[8:12] + "-" + uuid[12:16] + "-" + uuid[16:20] + "-" + uuid[20:32]
}

// Bytes 生成指定长度的随机字节数组
func Bytes(length int) []byte {
	if length <= 0 {
		return []byte{}
	}
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return b
}

// SecureInt 使用加密安全的随机数生成器生成整数 [0, max)
func SecureInt(max int) (int, error) {
	if max <= 0 {
		return 0, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

// SecureIntRange 使用加密安全的随机数生成器生成指定范围内的整数 [min, max)
func SecureIntRange(min, max int) (int, error) {
	if min >= max {
		return min, nil
	}
	diff := max - min
	n, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		return 0, err
	}
	return min + int(n.Int64()), nil
}

// SecureString 使用加密安全的随机数生成器生成随机字符串
func SecureString(length int) (string, error) {
	return SecureStringWithCharset(length, LettersNum)
}

// SecureStringWithCharset 使用加密安全的随机数生成器和指定字符集生成随机字符串
func SecureStringWithCharset(length int, charset string) (string, error) {
	if length <= 0 || len(charset) == 0 {
		return "", nil
	}

	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// SecureBytes 使用加密安全的随机数生成器生成随机字节数组
func SecureBytes(length int) ([]byte, error) {
	if length <= 0 {
		return []byte{}, nil
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Choice 从切片中随机选择一个元素
func Choice[T any](slice []T) (T, bool) {
	var zero T
	if len(slice) == 0 {
		return zero, false
	}
	return slice[defaultRand.IntN(len(slice))], true
}

// Choices 从切片中随机选择n个元素（允许重复）
func Choices[T any](slice []T, n int) []T {
	if len(slice) == 0 || n <= 0 {
		return []T{}
	}
	result := make([]T, n)
	for i := range result {
		result[i] = slice[defaultRand.IntN(len(slice))]
	}
	return result
}

// Sample 从切片中随机选择n个不重复的元素
func Sample[T comparable](slice []T, n int) []T {
	if len(slice) == 0 || n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		// 如果n大于等于切片长度，返回打乱后的整个切片
		result := make([]T, len(slice))
		copy(result, slice)
		Shuffle(result)
		return result
	}

	// 使用map记录已选择的索引
	selected := make(map[int]bool)
	result := make([]T, 0, n)
	for len(result) < n {
		idx := defaultRand.IntN(len(slice))
		if !selected[idx] {
			selected[idx] = true
			result = append(result, slice[idx])
		}
	}
	return result
}

// Shuffle 随机打乱切片（Fisher-Yates算法）
func Shuffle[T any](slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := defaultRand.IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// WeightedChoice 根据权重随机选择元素
func WeightedChoice[T any](items []T, weights []float64) (T, bool) {
	var zero T
	if len(items) == 0 || len(weights) == 0 || len(items) != len(weights) {
		return zero, false
	}

	// 计算总权重
	totalWeight := 0.0
	for _, w := range weights {
		if w < 0 {
			return zero, false
		}
		totalWeight += w
	}

	if totalWeight == 0 {
		return zero, false
	}

	// 生成随机数
	r := defaultRand.Float64() * totalWeight

	// 根据权重选择
	currentWeight := 0.0
	for i, w := range weights {
		currentWeight += w
		if r < currentWeight {
			return items[i], true
		}
	}

	// 边界情况，返回最后一个元素
	return items[len(items)-1], true
}

// Bool 生成随机布尔值
func Bool() bool {
	return defaultRand.IntN(2) == 1
}

// Probability 根据概率返回true
func Probability(p float64) bool {
	if p <= 0 {
		return false
	}
	if p >= 1 {
		return true
	}
	return defaultRand.Float64() < p
}

// Duration 生成指定范围内的随机时间间隔
func Duration(min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}
	diff := max - min
	return min + time.Duration(defaultRand.Int64N(int64(diff)))
}

// Time 生成指定时间范围内的随机时间
func Time(start, end time.Time) time.Time {
	if start.After(end) {
		return start
	}
	diff := end.Sub(start)
	return start.Add(time.Duration(defaultRand.Int64N(int64(diff))))
}

// Seed 设置随机数生成器的种子
func Seed(seed int64) {
	defaultRand = rand.New(rand.NewPCG(uint64(seed), uint64(seed)))
}

// FastString 快速生成随机字符串（使用unsafe，性能更高）
func FastString(length int) string {
	if length <= 0 {
		return ""
	}
	const charset = LettersNum
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[defaultRand.IntN(len(charset))]
	}
	return *(*string)(unsafe.Pointer(&b))
}
