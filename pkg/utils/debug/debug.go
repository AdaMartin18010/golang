package debug

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// Stack 获取当前调用栈
func Stack() []byte {
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, false)
	return buf[:n]
}

// StackAll 获取所有goroutine的调用栈
func StackAll() []byte {
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	return buf[:n]
}

// Caller 获取调用者信息
func Caller(skip int) (file string, line int, function string) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", 0, ""
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, ""
	}
	return file, line, fn.Name()
}

// Callers 获取调用栈
func Callers(skip int, depth int) []string {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip+2, pcs)
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:n])
	var result []string
	for {
		frame, more := frames.Next()
		result = append(result, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return result
}

// PrintStack 打印调用栈
func PrintStack() {
	fmt.Print(string(Stack()))
}

// PrintStackAll 打印所有goroutine的调用栈
func PrintStackAll() {
	fmt.Print(string(StackAll()))
}

// FuncName 获取函数名
func FuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	return fn.Name()
}

// FileLine 获取文件和行号
func FileLine(skip int) (file string, line int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", 0
	}
	return file, line
}

// Dump 打印变量的详细信息
func Dump(v interface{}) {
	fmt.Printf("%+v\n", v)
}

// DumpWithLabel 带标签打印变量
func DumpWithLabel(label string, v interface{}) {
	fmt.Printf("%s: %+v\n", label, v)
}

// DumpType 打印变量类型
func DumpType(v interface{}) {
	fmt.Printf("Type: %T\n", v)
}

// DumpValue 打印变量值
func DumpValue(v interface{}) {
	fmt.Printf("Value: %v\n", v)
}

// DumpStruct 打印结构体详细信息
func DumpStruct(v interface{}) {
	fmt.Printf("Struct: %+v\n", v)
}

// Trace 跟踪函数执行
func Trace(name string) func() {
	start := time.Now()
	fmt.Printf("[TRACE] Entering %s\n", name)
	return func() {
		fmt.Printf("[TRACE] Exiting %s (took %v)\n", name, time.Since(start))
	}
}

// TraceFunc 跟踪函数执行（带返回值）
func TraceFunc(name string, fn func()) {
	defer Trace(name)()
	fn()
}

// TraceFuncWithResult 跟踪函数执行（带返回值）
func TraceFuncWithResult[T any](name string, fn func() T) T {
	defer Trace(name)()
	return fn()
}

// Measure 测量函数执行时间
func Measure(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// MeasureWithResult 测量函数执行时间（带返回值）
func MeasureWithResult[T any](fn func() T) (T, time.Duration) {
	start := time.Now()
	result := fn()
	return result, time.Since(start)
}

// Benchmark 基准测试
func Benchmark(name string, iterations int, fn func()) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	duration := time.Since(start)
	avg := duration / time.Duration(iterations)
	fmt.Printf("[BENCHMARK] %s: %d iterations, total: %v, avg: %v\n", name, iterations, duration, avg)
}

// BenchmarkWithResult 基准测试（带返回值）
func BenchmarkWithResult[T any](name string, iterations int, fn func() T) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = fn()
	}
	duration := time.Since(start)
	avg := duration / time.Duration(iterations)
	fmt.Printf("[BENCHMARK] %s: %d iterations, total: %v, avg: %v\n", name, iterations, duration, avg)
}

// Assert 断言
func Assert(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("Assertion failed: %s", message))
	}
}

// AssertEqual 断言相等
func AssertEqual[T comparable](expected, actual T, message string) {
	if expected != actual {
		panic(fmt.Sprintf("Assertion failed: %s (expected: %v, actual: %v)", message, expected, actual))
	}
}

// AssertNotEqual 断言不相等
func AssertNotEqual[T comparable](expected, actual T, message string) {
	if expected == actual {
		panic(fmt.Sprintf("Assertion failed: %s (expected not: %v, actual: %v)", message, expected, actual))
	}
}

// AssertNil 断言nil
func AssertNil(v interface{}, message string) {
	if v != nil {
		panic(fmt.Sprintf("Assertion failed: %s (expected nil, got: %v)", message, v))
	}
}

// AssertNotNil 断言非nil
func AssertNotNil(v interface{}, message string) {
	if v == nil {
		panic(fmt.Sprintf("Assertion failed: %s (expected not nil)", message))
	}
}

// LogCall 记录函数调用
func LogCall(function string, args ...interface{}) {
	fmt.Printf("[CALL] %s(", function)
	for i, arg := range args {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%v", arg)
	}
	fmt.Println(")")
}

// LogReturn 记录函数返回
func LogReturn(function string, result interface{}) {
	fmt.Printf("[RETURN] %s -> %v\n", function, result)
}

// LogError 记录错误
func LogError(err error, context ...string) {
	if err != nil {
		ctx := strings.Join(context, " ")
		if ctx != "" {
			fmt.Printf("[ERROR] %s: %v\n", ctx, err)
		} else {
			fmt.Printf("[ERROR] %v\n", err)
		}
	}
}

// LogInfo 记录信息
func LogInfo(message string, args ...interface{}) {
	fmt.Printf("[INFO] "+message+"\n", args...)
}

// LogWarning 记录警告
func LogWarning(message string, args ...interface{}) {
	fmt.Printf("[WARNING] "+message+"\n", args...)
}

// LogDebug 记录调试信息
func LogDebug(message string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+message+"\n", args...)
}

// GetGoroutineID 获取当前goroutine ID
func GetGoroutineID() int64 {
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	var id int64
	fmt.Sscanf(idField, "%d", &id)
	return id
}

// GetNumGoroutines 获取goroutine数量
func GetNumGoroutines() int {
	return runtime.NumGoroutine()
}

// GetMemStats 获取内存统计
func GetMemStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

// PrintMemStats 打印内存统计
func PrintMemStats() {
	m := GetMemStats()
	fmt.Printf("[MEM] Alloc: %d KB, TotalAlloc: %d KB, Sys: %d KB, NumGC: %d\n",
		m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}

// GC 执行GC并打印统计
func GC() {
	runtime.GC()
	PrintMemStats()
}

// PrintGoroutines 打印所有goroutine信息
func PrintGoroutines() {
	fmt.Print(string(StackAll()))
}

// IsDebug 检查是否在调试模式
var IsDebug = false

// SetDebug 设置调试模式
func SetDebug(debug bool) {
	IsDebug = debug
}

// DebugPrint 调试打印（仅在调试模式下）
func DebugPrint(message string, args ...interface{}) {
	if IsDebug {
		fmt.Printf("[DEBUG] "+message+"\n", args...)
	}
}

// DebugDump 调试转储（仅在调试模式下）
func DebugDump(v interface{}) {
	if IsDebug {
		Dump(v)
	}
}

// DebugTrace 调试跟踪（仅在调试模式下）
func DebugTrace(name string) func() {
	if !IsDebug {
		return func() {}
	}
	return Trace(name)
}
