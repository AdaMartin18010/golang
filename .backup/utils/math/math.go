package math

import (
	"math"
	"math/rand"
	"time"
)

// Max 返回两个整数中的较大值
func Max(a, b int) int {
	if a > b {
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

// MaxFloat64 返回两个float64中的较大值
func MaxFloat64(a, b float64) float64 {
	return math.Max(a, b)
}

// Min 返回两个整数中的较小值
func Min(a, b int) int {
	if a < b {
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

// MinFloat64 返回两个float64中的较小值
func MinFloat64(a, b float64) float64 {
	return math.Min(a, b)
}

// MaxInts 返回整数切片中的最大值
func MaxInts(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
	}
	return max
}

// MinInts 返回整数切片中的最小值
func MinInts(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	min := nums[0]
	for _, num := range nums[1:] {
		if num < min {
			min = num
		}
	}
	return min
}

// Sum 计算整数切片的总和
func Sum(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// SumFloat64 计算float64切片的总和
func SumFloat64(nums []float64) float64 {
	sum := 0.0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// Average 计算整数切片的平均值
func Average(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	return float64(Sum(nums)) / float64(len(nums))
}

// AverageFloat64 计算float64切片的平均值
func AverageFloat64(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	return SumFloat64(nums) / float64(len(nums))
}

// Abs 返回整数的绝对值
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// AbsInt64 返回int64的绝对值
func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsFloat64 返回float64的绝对值
func AbsFloat64(x float64) float64 {
	return math.Abs(x)
}

// Pow 计算x的y次方
func Pow(x, y float64) float64 {
	return math.Pow(x, y)
}

// Sqrt 计算平方根
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

// Ceil 向上取整
func Ceil(x float64) float64 {
	return math.Ceil(x)
}

// Floor 向下取整
func Floor(x float64) float64 {
	return math.Floor(x)
}

// Round 四舍五入
func Round(x float64) float64 {
	return math.Round(x)
}

// RoundTo 四舍五入到指定小数位
func RoundTo(x float64, places int) float64 {
	multiplier := math.Pow(10, float64(places))
	return math.Round(x*multiplier) / multiplier
}

// Clamp 将值限制在[min, max]范围内
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat64 将值限制在[min, max]范围内
func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// IsInRange 检查值是否在范围内
func IsInRange(value, min, max int) bool {
	return value >= min && value <= max
}

// IsInRangeFloat64 检查值是否在范围内
func IsInRangeFloat64(value, min, max float64) bool {
	return value >= min && value <= max
}

// RandomInt 生成指定范围内的随机整数
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// RandomFloat64 生成指定范围内的随机浮点数
func RandomFloat64(min, max float64) float64 {
	if min >= max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// GCD 计算最大公约数
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return Abs(a)
}

// LCM 计算最小公倍数
func LCM(a, b int) int {
	return Abs(a*b) / GCD(a, b)
}

// Factorial 计算阶乘
func Factorial(n int) int64 {
	if n < 0 {
		return 0
	}
	if n == 0 || n == 1 {
		return 1
	}
	result := int64(1)
	for i := 2; i <= n; i++ {
		result *= int64(i)
	}
	return result
}

// IsPrime 检查是否为质数
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// NextPrime 返回下一个质数
func NextPrime(n int) int {
	if n < 2 {
		return 2
	}
	for {
		n++
		if IsPrime(n) {
			return n
		}
	}
}

// Fibonacci 计算斐波那契数列的第n项
func Fibonacci(n int) int64 {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return int64(n)
	}
	a, b := int64(0), int64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// Percent 计算百分比
func Percent(value, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

// PercentOf 计算某个值占另一个值的百分比
func PercentOf(value, of float64) float64 {
	if of == 0 {
		return 0
	}
	return (value / of) * 100
}

// PercentChange 计算百分比变化
func PercentChange(oldValue, newValue float64) float64 {
	if oldValue == 0 {
		return 0
	}
	return ((newValue - oldValue) / oldValue) * 100
}

// Lerp 线性插值
func Lerp(start, end, t float64) float64 {
	return start + (end-start)*t
}

// InverseLerp 反向线性插值
func InverseLerp(start, end, value float64) float64 {
	if start == end {
		return 0
	}
	return (value - start) / (end - start)
}

// Remap 重新映射值到新范围
func Remap(value, oldMin, oldMax, newMin, newMax float64) float64 {
	t := InverseLerp(oldMin, oldMax, value)
	return Lerp(newMin, newMax, t)
}

// Distance 计算两点之间的距离
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// Distance3D 计算三维空间中两点之间的距离
func Distance3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// DegToRad 角度转弧度
func DegToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RadToDeg 弧度转角度
func RadToDeg(radians float64) float64 {
	return radians * 180 / math.Pi
}

// Sin 正弦函数
func Sin(x float64) float64 {
	return math.Sin(x)
}

// Cos 余弦函数
func Cos(x float64) float64 {
	return math.Cos(x)
}

// Tan 正切函数
func Tan(x float64) float64 {
	return math.Tan(x)
}

// Log 自然对数
func Log(x float64) float64 {
	return math.Log(x)
}

// Log10 以10为底的对数
func Log10(x float64) float64 {
	return math.Log10(x)
}

// Exp 自然指数
func Exp(x float64) float64 {
	return math.Exp(x)
}

// Mod 取模运算
func Mod(x, y float64) float64 {
	return math.Mod(x, y)
}

// ModInt 整数取模运算
func ModInt(x, y int) int {
	return x % y
}

// IsEven 检查是否为偶数
func IsEven(n int) bool {
	return n%2 == 0
}

// IsOdd 检查是否为奇数
func IsOdd(n int) bool {
	return n%2 != 0
}

// IsDivisible 检查是否可整除
func IsDivisible(n, divisor int) bool {
	return n%divisor == 0
}
