# 调试工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [调试工具](#调试工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

调试工具提供了各种调试功能，包括调用栈获取、函数跟踪、性能测量、断言、日志记录、内存统计等，帮助开发者快速定位和解决问题。

---

## 2. 功能特性

### 2.1 调用栈

- `Stack`: 获取当前调用栈
- `StackAll`: 获取所有goroutine的调用栈
- `Caller`: 获取调用者信息
- `Callers`: 获取调用栈
- `PrintStack`: 打印调用栈
- `PrintStackAll`: 打印所有goroutine的调用栈
- `FuncName`: 获取函数名
- `FileLine`: 获取文件和行号

### 2.2 变量转储

- `Dump`: 打印变量的详细信息
- `DumpWithLabel`: 带标签打印变量
- `DumpType`: 打印变量类型
- `DumpValue`: 打印变量值
- `DumpStruct`: 打印结构体详细信息

### 2.3 函数跟踪

- `Trace`: 跟踪函数执行
- `TraceFunc`: 跟踪函数执行（带返回值）
- `TraceFuncWithResult`: 跟踪函数执行（带返回值）

### 2.4 性能测量

- `Measure`: 测量函数执行时间
- `MeasureWithResult`: 测量函数执行时间（带返回值）
- `Benchmark`: 基准测试
- `BenchmarkWithResult`: 基准测试（带返回值）

### 2.5 断言

- `Assert`: 断言
- `AssertEqual`: 断言相等
- `AssertNotEqual`: 断言不相等
- `AssertNil`: 断言nil
- `AssertNotNil`: 断言非nil

### 2.6 日志记录

- `LogCall`: 记录函数调用
- `LogReturn`: 记录函数返回
- `LogError`: 记录错误
- `LogInfo`: 记录信息
- `LogWarning`: 记录警告
- `LogDebug`: 记录调试信息

### 2.7 运行时信息

- `GetGoroutineID`: 获取当前goroutine ID
- `GetNumGoroutines`: 获取goroutine数量
- `GetMemStats`: 获取内存统计
- `PrintMemStats`: 打印内存统计
- `GC`: 执行GC并打印统计
- `PrintGoroutines`: 打印所有goroutine信息

### 2.8 调试模式

- `IsDebug`: 检查是否在调试模式
- `SetDebug`: 设置调试模式
- `DebugPrint`: 调试打印（仅在调试模式下）
- `DebugDump`: 调试转储（仅在调试模式下）
- `DebugTrace`: 调试跟踪（仅在调试模式下）

---

## 3. 使用示例

### 3.1 调用栈

```go
import "github.com/yourusername/golang/pkg/utils/debug"

// 获取调用栈
stack := debug.Stack()
fmt.Print(string(stack))

// 获取所有goroutine的调用栈
allStack := debug.StackAll()
fmt.Print(string(allStack))

// 获取调用者信息
file, line, function := debug.Caller(0)
fmt.Printf("File: %s, Line: %d, Function: %s\n", file, line, function)

// 获取调用栈
callers := debug.Callers(0, 5)
for _, caller := range callers {
    fmt.Println(caller)
}

// 打印调用栈
debug.PrintStack()
```

### 3.2 变量转储

```go
// 打印变量详细信息
debug.Dump(variable)

// 带标签打印
debug.DumpWithLabel("variable", variable)

// 打印类型
debug.DumpType(variable)

// 打印值
debug.DumpValue(variable)

// 打印结构体
debug.DumpStruct(structVar)
```

### 3.3 函数跟踪

```go
// 跟踪函数执行
defer debug.Trace("myFunction")()
// ... 函数代码 ...

// 跟踪函数执行（带返回值）
debug.TraceFunc("myFunction", func() {
    // ... 函数代码 ...
})

// 跟踪函数执行（带返回值）
result := debug.TraceFuncWithResult("myFunction", func() int {
    return 42
})
```

### 3.4 性能测量

```go
// 测量函数执行时间
duration := debug.Measure(func() {
    // ... 代码 ...
})
fmt.Printf("Took: %v\n", duration)

// 测量函数执行时间（带返回值）
result, duration := debug.MeasureWithResult(func() int {
    return 42
})

// 基准测试
debug.Benchmark("myFunction", 1000, func() {
    // ... 代码 ...
})
```

### 3.5 断言

```go
// 断言
debug.Assert(condition, "condition must be true")

// 断言相等
debug.AssertEqual(expected, actual, "values must be equal")

// 断言不相等
debug.AssertNotEqual(expected, actual, "values must not be equal")

// 断言nil
debug.AssertNil(value, "value must be nil")

// 断言非nil
debug.AssertNotNil(value, "value must not be nil")
```

### 3.6 日志记录

```go
// 记录函数调用
debug.LogCall("myFunction", arg1, arg2)

// 记录函数返回
debug.LogReturn("myFunction", result)

// 记录错误
debug.LogError(err, "context")

// 记录信息
debug.LogInfo("message: %s", "info")

// 记录警告
debug.LogWarning("message: %s", "warning")

// 记录调试信息
debug.LogDebug("message: %s", "debug")
```

### 3.7 运行时信息

```go
// 获取goroutine ID
id := debug.GetGoroutineID()
fmt.Printf("Goroutine ID: %d\n", id)

// 获取goroutine数量
num := debug.GetNumGoroutines()
fmt.Printf("Goroutines: %d\n", num)

// 获取内存统计
stats := debug.GetMemStats()
fmt.Printf("Alloc: %d KB\n", stats.Alloc/1024)

// 打印内存统计
debug.PrintMemStats()

// 执行GC并打印统计
debug.GC()

// 打印所有goroutine信息
debug.PrintGoroutines()
```

### 3.8 调试模式

```go
// 设置调试模式
debug.SetDebug(true)

// 调试打印（仅在调试模式下）
debug.DebugPrint("message: %s", "debug")

// 调试转储（仅在调试模式下）
debug.DebugDump(variable)

// 调试跟踪（仅在调试模式下）
defer debug.DebugTrace("myFunction")()
```

### 3.9 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/debug"
)

func main() {
    // 设置调试模式
    debug.SetDebug(true)

    // 跟踪函数执行
    defer debug.Trace("main")()

    // 测量执行时间
    duration := debug.Measure(func() {
        // 执行一些操作
        for i := 0; i < 1000; i++ {
            _ = i * 2
        }
    })
    fmt.Printf("Execution took: %v\n", duration)

    // 打印内存统计
    debug.PrintMemStats()

    // 调试打印
    debug.DebugPrint("Debug message")
}
```

---

**更新日期**: 2025-11-11
