package compare

import (
	"reflect"
	"time"
)

// Equal 检查两个值是否相等
func Equal(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// NotEqual 检查两个值是否不相等
func NotEqual(a, b interface{}) bool {
	return !Equal(a, b)
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

// CompareInt64 比较两个int64
func CompareInt64(a, b int64) int {
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

// CompareTime 比较两个时间
func CompareTime(a, b time.Time) int {
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}

// Less 检查a是否小于b
func Less(a, b interface{}) bool {
	switch aVal := a.(type) {
	case int:
		if bVal, ok := b.(int); ok {
			return aVal < bVal
		}
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal < bVal
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal < bVal
		}
	case string:
		if bVal, ok := b.(string); ok {
			return aVal < bVal
		}
	case time.Time:
		if bVal, ok := b.(time.Time); ok {
			return aVal.Before(bVal)
		}
	}
	return false
}

// Greater 检查a是否大于b
func Greater(a, b interface{}) bool {
	switch aVal := a.(type) {
	case int:
		if bVal, ok := b.(int); ok {
			return aVal > bVal
		}
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal > bVal
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal > bVal
		}
	case string:
		if bVal, ok := b.(string); ok {
			return aVal > bVal
		}
	case time.Time:
		if bVal, ok := b.(time.Time); ok {
			return aVal.After(bVal)
		}
	}
	return false
}

// LessOrEqual 检查a是否小于等于b
func LessOrEqual(a, b interface{}) bool {
	return Less(a, b) || Equal(a, b)
}

// GreaterOrEqual 检查a是否大于等于b
func GreaterOrEqual(a, b interface{}) bool {
	return Greater(a, b) || Equal(a, b)
}

// Min 返回两个值中的较小值
func Min(a, b interface{}) interface{} {
	if Less(a, b) {
		return a
	}
	return b
}

// Max 返回两个值中的较大值
func Max(a, b interface{}) interface{} {
	if Greater(a, b) {
		return a
	}
	return b
}

// MinInt 返回两个整数中的较小值
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt 返回两个整数中的较大值
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt64 返回两个int64中的较小值
func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// MaxInt64 返回两个int64中的较大值
func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat64 返回两个float64中的较小值
func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// MaxFloat64 返回两个float64中的较大值
func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinString 返回两个字符串中的较小值（字典序）
func MinString(a, b string) string {
	if a < b {
		return a
	}
	return b
}

// MaxString 返回两个字符串中的较大值（字典序）
func MaxString(a, b string) string {
	if a > b {
		return a
	}
	return b
}

// MinTime 返回两个时间中的较早时间
func MinTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// MaxTime 返回两个时间中的较晚时间
func MaxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// InRange 检查值是否在范围内
func InRange(value, min, max interface{}) bool {
	return GreaterOrEqual(value, min) && LessOrEqual(value, max)
}

// InRangeInt 检查整数是否在范围内
func InRangeInt(value, min, max int) bool {
	return value >= min && value <= max
}

// InRangeInt64 检查int64是否在范围内
func InRangeInt64(value, min, max int64) bool {
	return value >= min && value <= max
}

// InRangeFloat64 检查float64是否在范围内
func InRangeFloat64(value, min, max float64) bool {
	return value >= min && value <= max
}

// Clamp 将值限制在[min, max]范围内
func Clamp(value, min, max interface{}) interface{} {
	if Less(value, min) {
		return min
	}
	if Greater(value, max) {
		return max
	}
	return value
}

// ClampInt 将整数限制在[min, max]范围内
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt64 将int64限制在[min, max]范围内
func ClampInt64(value, min, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat64 将float64限制在[min, max]范围内
func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// IsZero 检查值是否为零值
func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Bool:
		return !val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return val.Complex() == 0
	case reflect.String:
		return val.String() == ""
	case reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	case reflect.Struct:
		return val.IsZero()
	default:
		return false
	}
}

// IsNil 检查值是否为nil
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

// IsEmpty 检查值是否为空（nil、零值或空集合）
func IsEmpty(v interface{}) bool {
	return IsNil(v) || IsZero(v)
}

// CompareFunc 比较函数类型
type CompareFunc[T any] func(a, b T) int

// LessFunc 小于函数类型
type LessFunc[T any] func(a, b T) bool

// EqualFunc 相等函数类型
type EqualFunc[T any] func(a, b T) bool

// Compare 使用比较函数比较两个值
func Compare[T any](a, b T, cmp CompareFunc[T]) int {
	return cmp(a, b)
}

// LessThan 使用小于函数比较两个值
func LessThan[T any](a, b T, less LessFunc[T]) bool {
	return less(a, b)
}

// EqualTo 使用相等函数比较两个值
func EqualTo[T any](a, b T, eq EqualFunc[T]) bool {
	return eq(a, b)
}

// CompareBy 根据键函数比较两个值
func CompareBy[T any, K comparable](a, b T, keyFunc func(T) K, cmp func(K, K) int) int {
	return cmp(keyFunc(a), keyFunc(b))
}

// LessBy 根据键函数检查a是否小于b
func LessBy[T any, K comparable](a, b T, keyFunc func(T) K, less func(K, K) bool) bool {
	return less(keyFunc(a), keyFunc(b))
}

// EqualBy 根据键函数检查a是否等于b
func EqualBy[T any, K comparable](a, b T, keyFunc func(T) K) bool {
	return keyFunc(a) == keyFunc(b)
}

// CompareSlice 比较两个切片
func CompareSlice[T comparable](a, b []T) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}
	return 0
}

// CompareSliceFunc 使用比较函数比较两个切片
func CompareSliceFunc[T any](a, b []T, cmp func(T, T) int) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		result := cmp(a[i], b[i])
		if result != 0 {
			return result
		}
	}
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}
	return 0
}

// CompareMap 比较两个映射
func CompareMap[K comparable, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || v != bv {
			return false
		}
	}
	return true
}

// CompareMapFunc 使用比较函数比较两个映射
func CompareMapFunc[K comparable, V any](a, b map[K]V, cmp func(V, V) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !cmp(v, bv) {
			return false
		}
	}
	return true
}
