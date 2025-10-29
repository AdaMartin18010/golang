# Go 1.23 testing包增强详解

> **难度**: ⭐⭐⭐⭐
> **标签**: #Go1.23 #testing #slogtest #并发测试

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---


---

## 📋 目录


- [1. testing包增强概述](#1.-testing包增强概述)
  - [1.1 Go 1.23的testing改进](#11-go-1.23的testing改进)
  - [1.2 核心价值](#12-核心价值)
- [2. testing/slogtest包详解](#2.-testingslogtest包详解)
  - [2.1 slogtest简介](#21-slogtest简介)
  - [2.2 基本用法](#22-基本用法)
  - [2.3 测试自定义Handler](#23-测试自定义handler)
  - [2.4 常见测试场景](#24-常见测试场景)
- [3. 测试输出改进](#3.-测试输出改进)
  - [3.1 更清晰的失败信息](#31-更清晰的失败信息)
  - [3.2 并行测试输出](#32-并行测试输出)
  - [3.3 子测试可视化](#33-子测试可视化)
- [4. 并发测试增强](#4.-并发测试增强)
  - [4.1 t.Parallel()改进](#41-t.parallel改进)
  - [4.2 并发测试最佳实践](#42-并发测试最佳实践)
  - [4.3 死锁检测](#43-死锁检测)
- [5. 基准测试改进](#5.-基准测试改进)
  - [5.1 内存分配报告](#51-内存分配报告)
  - [5.2 性能回归检测](#52-性能回归检测)
  - [5.3 benchstat集成](#53-benchstat集成)
- [6. Fuzzing增强](#6.-fuzzing增强)
  - [6.1 模糊测试改进](#61-模糊测试改进)
  - [6.2 语料库管理](#62-语料库管理)
  - [6.3 实战案例](#63-实战案例)
- [7. 测试覆盖率增强](#7.-测试覆盖率增强)
  - [7.1 更精确的覆盖率](#71-更精确的覆盖率)
  - [7.2 函数级覆盖率](#72-函数级覆盖率)
  - [7.3 HTML报告改进](#73-html报告改进)
- [8. 测试工具函数](#8.-测试工具函数)
  - [8.1 t.TempDir()最佳实践](#81-t.tempdir最佳实践)
  - [8.2 t.Setenv()使用](#82-t.setenv使用)
  - [8.3 t.Cleanup()模式](#83-t.cleanup模式)
- [9. 实战案例](#9.-实战案例)
  - [9.1 完整的日志Handler测试](#91-完整的日志handler测试)
  - [9.2 并发服务测试](#92-并发服务测试)
  - [9.3 性能基准测试套件](#93-性能基准测试套件)
- [10. 最佳实践](#10.-最佳实践)
  - [10.1 测试组织](#101-测试组织)
  - [10.2 测试命名](#102-测试命名)
  - [10.3 测试数据管理](#103-测试数据管理)
- [11. 参考资源](#11.-参考资源)
  - [官方文档](#官方文档)
  - [测试工具](#测试工具)
  - [博客文章](#博客文章)

## 1. testing包增强概述

### 1.1 Go 1.23的testing改进

**主要增强**:

1. **testing/slogtest包**（新增）
   - 用于测试slog.Handler实现
   - 验证日志处理器的正确性
   - 标准化的测试方法

2. **测试输出改进**
   - 更清晰的失败信息
   - 改进的并行测试输出
   - 更好的子测试可视化

3. **并发测试增强**
   - t.Parallel()的性能改进
   - 更好的死锁检测
   - 并发测试隔离

4. **基准测试改进**
   - 更详细的内存分配报告
   - 性能回归检测
   - benchstat工具增强

5. **Fuzzing增强**
   - 改进的语料库管理
   - 更快的模糊测试
   - 更好的错误报告

### 1.2 核心价值

| 改进 | 价值 |
|------|------|
| **testing/slogtest** | 标准化日志Handler测试 |
| **输出改进** | 更快定位问题 |
| **并发增强** | 更可靠的并发测试 |
| **基准测试** | 更准确的性能分析 |
| **Fuzzing** | 发现更多边界情况 |

---

## 2. testing/slogtest包详解

### 2.1 slogtest简介

**testing/slogtest**是Go 1.23新增的包，用于测试`log/slog.Handler`实现。

**核心函数**:

```go
package slogtest

// TestHandler测试Handler实现是否符合slog规范
func TestHandler(h slog.Handler, newHandler func() slog.Handler) error

// Run在testing.T中运行Handler测试
func Run(t *testing.T, newHandler func() slog.Handler, checks ...Check)
```

### 2.2 基本用法

**示例1：测试标准Handler**

```go
package mylog_test

import (
    "log/slog"
    "testing"
    "testing/slogtest"
)

func TestJSONHandler(t *testing.T) {
    var buf bytes.Buffer
    
    // 创建Handler工厂函数
    newHandler := func() slog.Handler {
        buf.Reset()
        return slog.NewJSONHandler(&buf, nil)
    }
    
    // 运行标准测试
    slogtest.Run(t, newHandler, slogtest.All)
}
```

**示例2：使用TestHandler**

```go
func TestCustomHandler(t *testing.T) {
    h := NewCustomHandler()
    
    newHandler := func() slog.Handler {
        return NewCustomHandler()
    }
    
    // 执行测试，返回错误
    if err := slogtest.TestHandler(h, newHandler); err != nil {
        t.Error(err)
    }
}
```

### 2.3 测试自定义Handler

**完整示例：自定义Handler测试**

```go
package customlog

import (
    "bytes"
    "context"
    "encoding/json"
    "log/slog"
    "testing"
    "testing/slogtest"
)

// CustomHandler自定义日志处理器
type CustomHandler struct {
    buf   *bytes.Buffer
    attrs []slog.Attr
    group string
}

func NewCustomHandler(buf *bytes.Buffer) *CustomHandler {
    return &CustomHandler{buf: buf}
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
    return true
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
    entry := make(map[string]interface{})
    
    // 添加基本字段
    entry["time"] = r.Time
    entry["level"] = r.Level.String()
    entry["msg"] = r.Message
    
    // 添加属性
    r.Attrs(func(a slog.Attr) bool {
        entry[a.Key] = a.Value.Any()
        return true
    })
    
    // 编码为JSON
    data, err := json.Marshal(entry)
    if err != nil {
        return err
    }
    
    h.buf.Write(data)
    h.buf.WriteByte('\n')
    return nil
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    newHandler := *h
    newHandler.attrs = append(newHandler.attrs, attrs...)
    return &newHandler
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
    newHandler := *h
    newHandler.group = name
    return &newHandler
}

// 测试
func TestCustomHandler(t *testing.T) {
    var buf bytes.Buffer
    
    newHandler := func() slog.Handler {
        buf.Reset()
        return NewCustomHandler(&buf)
    }
    
    // 运行所有标准测试
    slogtest.Run(t, newHandler, slogtest.All)
}

// 测试特定方面
func TestCustomHandlerWithAttrs(t *testing.T) {
    var buf bytes.Buffer
    h := NewCustomHandler(&buf)
    
    logger := slog.New(h)
    logger = logger.With("key1", "value1")
    logger.Info("test message", "key2", "value2")
    
    // 验证输出
    output := buf.String()
    if !strings.Contains(output, "key1") {
        t.Error("Missing key1")
    }
    if !strings.Contains(output, "key2") {
        t.Error("Missing key2")
    }
}
```

### 2.4 常见测试场景

**场景1：测试日志级别过滤**

```go
func TestHandlerLevelFilter(t *testing.T) {
    var buf bytes.Buffer
    
    h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
        Level: slog.LevelWarn,
    })
    
    logger := slog.New(h)
    
    // 应该被过滤
    logger.Debug("debug message")
    logger.Info("info message")
    
    // 应该输出
    logger.Warn("warn message")
    logger.Error("error message")
    
    output := buf.String()
    if strings.Contains(output, "debug") || strings.Contains(output, "info") {
        t.Error("Debug/Info messages should be filtered")
    }
    if !strings.Contains(output, "warn") || !strings.Contains(output, "error") {
        t.Error("Warn/Error messages should be present")
    }
}
```

**场景2：测试属性组**

```go
func TestHandlerGroups(t *testing.T) {
    var buf bytes.Buffer
    h := slog.NewJSONHandler(&buf, nil)
    
    logger := slog.New(h)
    logger = logger.WithGroup("request")
    logger.Info("handling request",
        "method", "GET",
        "path", "/api/users",
    )
    
    var result map[string]interface{}
    if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
        t.Fatal(err)
    }
    
    // 验证嵌套结构
    request, ok := result["request"].(map[string]interface{})
    if !ok {
        t.Fatal("Expected request group")
    }
    
    if request["method"] != "GET" {
        t.Error("Expected method=GET")
    }
}
```

**场景3：测试上下文处理**

```go
func TestHandlerContext(t *testing.T) {
    var buf bytes.Buffer
    
    // 自定义Handler，从context提取值
    h := NewContextAwareHandler(&buf)
    logger := slog.New(h)
    
    // 创建带值的context
    ctx := context.WithValue(context.Background(), "request_id", "req-123")
    
    logger.InfoContext(ctx, "processing request")
    
    output := buf.String()
    if !strings.Contains(output, "req-123") {
        t.Error("Request ID should be in output")
    }
}
```

---

## 3. 测试输出改进

### 3.1 更清晰的失败信息

**Go 1.23改进**:

```go
// Go 1.22及之前：失败信息可能不够清晰
// === RUN   TestExample
// --- FAIL: TestExample (0.00s)
//     example_test.go:10: assertion failed

// Go 1.23：更详细的上下文
// === RUN   TestExample
// --- FAIL: TestExample (0.00s)
//     example_test.go:10: assertion failed
//         Expected: 42
//         Got:      0
//         Diff:     +42
```

**最佳实践：使用t.Helper()**

```go
func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()  // 标记为辅助函数，错误指向调用者
    
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

func TestSomething(t *testing.T) {
    result := compute()
    assertEqual(t, result, 42)  // 错误会指向这一行
}
```

### 3.2 并行测试输出

**Go 1.23改进**:

```go
func TestParallelSuite(t *testing.T) {
    tests := []struct {
        name string
        fn   func(*testing.T)
    }{
        {"test1", testCase1},
        {"test2", testCase2},
        {"test3", testCase3},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // Go 1.23输出更有组织
            tt.fn(t)
        })
    }
}

// Go 1.23输出示例：
// === RUN   TestParallelSuite
// === PAUSE TestParallelSuite
// === CONT  TestParallelSuite
// === RUN   TestParallelSuite/test1
// === PAUSE TestParallelSuite/test1
// === RUN   TestParallelSuite/test2
// === PAUSE TestParallelSuite/test2
// === RUN   TestParallelSuite/test3
// === PAUSE TestParallelSuite/test3
// === CONT  TestParallelSuite/test1
// === CONT  TestParallelSuite/test2
// === CONT  TestParallelSuite/test3
// --- PASS: TestParallelSuite/test1 (0.10s)
// --- PASS: TestParallelSuite/test2 (0.15s)
// --- PASS: TestParallelSuite/test3 (0.20s)
// --- PASS: TestParallelSuite (0.20s)
```

### 3.3 子测试可视化

**改进的子测试输出**:

```go
func TestNestedSubtests(t *testing.T) {
    t.Run("group1", func(t *testing.T) {
        t.Run("case1", func(t *testing.T) {
            // 测试代码
        })
        t.Run("case2", func(t *testing.T) {
            // 测试代码
        })
    })
    
    t.Run("group2", func(t *testing.T) {
        t.Run("case1", func(t *testing.T) {
            // 测试代码
        })
    })
}

// Go 1.23输出：
// === RUN   TestNestedSubtests
// === RUN   TestNestedSubtests/group1
// === RUN   TestNestedSubtests/group1/case1
// --- PASS: TestNestedSubtests/group1/case1 (0.00s)
// === RUN   TestNestedSubtests/group1/case2
// --- PASS: TestNestedSubtests/group1/case2 (0.00s)
// --- PASS: TestNestedSubtests/group1 (0.00s)
// === RUN   TestNestedSubtests/group2
// === RUN   TestNestedSubtests/group2/case1
// --- PASS: TestNestedSubtests/group2/case1 (0.00s)
// --- PASS: TestNestedSubtests/group2 (0.00s)
// --- PASS: TestNestedSubtests (0.00s)
```

---

## 4. 并发测试增强

### 4.1 t.Parallel()改进

**Go 1.23性能改进**:

```go
func TestConcurrentOperations(t *testing.T) {
    // Go 1.23：t.Parallel()的调度更高效
    for i := 0; i < 100; i++ {
        i := i
        t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
            t.Parallel()  // 更好的并发控制
            
            // 测试代码
            result := expensiveOperation(i)
            if result != expected {
                t.Errorf("got %v, want %v", result, expected)
            }
        })
    }
}
```

### 4.2 并发测试最佳实践

**模式1：共享资源隔离**

```go
func TestConcurrentAccess(t *testing.T) {
    tests := []struct {
        name string
        data int
    }{
        {"test1", 1},
        {"test2", 2},
        {"test3", 3},
    }
    
    for _, tt := range tests {
        tt := tt  // 捕获循环变量
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            // 每个子测试有独立的资源
            resource := newTestResource()
            defer resource.Close()
            
            // 测试代码
            result := resource.Process(tt.data)
            assertEqual(t, result, tt.data*2)
        })
    }
}
```

**模式2：并发安全验证**

```go
func TestConcurrentMapAccess(t *testing.T) {
    m := &sync.Map{}
    
    // 并发写入
    t.Run("concurrent_writes", func(t *testing.T) {
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            go func(i int) {
                defer wg.Done()
                m.Store(i, i*2)
            }(i)
        }
        wg.Wait()
    })
    
    // 验证结果
    t.Run("verify_results", func(t *testing.T) {
        for i := 0; i < 100; i++ {
            val, ok := m.Load(i)
            if !ok {
                t.Errorf("key %d not found", i)
                continue
            }
            if val != i*2 {
                t.Errorf("key %d: got %v, want %v", i, val, i*2)
            }
        }
    })
}
```

### 4.3 死锁检测

**Go 1.23增强的超时检测**:

```go
func TestNoDeadlock(t *testing.T) {
    // Go 1.23会更快检测到死锁情况
    ch := make(chan int)
    
    done := make(chan bool)
    go func() {
        defer close(done)
        
        // 这会超时，Go 1.23会报告
        select {
        case v := <-ch:
            t.Logf("received %d", v)
        case <-time.After(1 * time.Second):
            t.Error("timeout waiting for value")
        }
    }()
    
    <-done
}
```

---

## 5. 基准测试改进

### 5.1 内存分配报告

**Go 1.23更详细的报告**:

```go
func BenchmarkStringConcat(b *testing.B) {
    b.ReportAllocs()  // Go 1.23提供更详细的分配信息
    
    for i := 0; i < b.N; i++ {
        s := "hello"
        s += " world"
        _ = s
    }
}

// Go 1.23输出示例：
// BenchmarkStringConcat-8   10000000   112 ns/op   32 B/op   2 allocs/op
//   Allocations by size:
//     16 bytes: 1 alloc
//     16 bytes: 1 alloc
```

### 5.2 性能回归检测

**使用benchstat检测回归**:

```bash
# 运行基准测试，保存结果
go test -bench=. -count=10 > old.txt

# 修改代码后再次运行
go test -bench=. -count=10 > new.txt

# 比较结果
benchstat old.txt new.txt
```

**示例输出**:

```text
name              old time/op    new time/op    delta
StringConcat-8     112ns ± 2%      98ns ± 1%   -12.50%  (p=0.000 n=10+10)

name              old alloc/op   new alloc/op   delta
StringConcat-8     32.0B ± 0%     16.0B ± 0%   -50.00%  (p=0.000 n=10+10)

name              old allocs/op  new allocs/op  delta
StringConcat-8      2.00 ± 0%      1.00 ± 0%   -50.00%  (p=0.000 n=10+10)
```

### 5.3 benchstat集成

**在CI中集成性能测试**:

```go
// benchmark_test.go
func BenchmarkCriticalPath(b *testing.B) {
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        result := criticalOperation()
        if result == nil {
            b.Fatal("unexpected nil")
        }
    }
}

// 在CI中运行
// go test -bench=. -benchmem -benchtime=10s
```

---

## 6. Fuzzing增强

### 6.1 模糊测试改进

**Go 1.23的Fuzzing增强**:

```go
func FuzzParseInput(f *testing.F) {
    // 添加种子语料
    f.Add("hello")
    f.Add("world")
    f.Add("12345")
    
    f.Fuzz(func(t *testing.T, input string) {
        // Go 1.23：更快的模糊测试
        result, err := ParseInput(input)
        
        if err != nil {
            // 预期的错误可以跳过
            if isExpectedError(err) {
                t.Skip()
            }
            t.Errorf("unexpected error: %v", err)
            return
        }
        
        // 验证结果
        if result == nil {
            t.Error("result should not be nil")
        }
    })
}
```

### 6.2 语料库管理

**改进的语料库组织**:

```text
testdata/
└── fuzz/
    └── FuzzParseInput/
        ├── corpus/
        │   ├── seed1
        │   ├── seed2
        │   └── seed3
        └── crashers/
            └── crash1
```

**添加自定义语料**:

```go
func FuzzJSON(f *testing.F) {
    // 从文件加载语料
    corpus, _ := os.ReadDir("testdata/json")
    for _, entry := range corpus {
        data, _ := os.ReadFile(filepath.Join("testdata/json", entry.Name()))
        f.Add(data)
    }
    
    f.Fuzz(func(t *testing.T, data []byte) {
        var v interface{}
        _ = json.Unmarshal(data, &v)
        // 不应该panic
    })
}
```

### 6.3 实战案例

**模糊测试URL解析器**:

```go
func FuzzURLParser(f *testing.F) {
    // 添加有效的URL种子
    f.Add("http://example.com")
    f.Add("https://example.com/path?query=value")
    f.Add("ftp://example.com:21/file.txt")
    
    f.Fuzz(func(t *testing.T, input string) {
        u, err := url.Parse(input)
        
        if err != nil {
            // 某些输入预期会失败
            return
        }
        
        // 验证解析结果的一致性
        reconstructed := u.String()
        u2, err2 := url.Parse(reconstructed)
        
        if err2 != nil {
            t.Errorf("re-parsing failed: %v", err2)
        }
        
        if u.Scheme != u2.Scheme || u.Host != u2.Host {
            t.Errorf("inconsistent parsing: %v vs %v", u, u2)
        }
    })
}
```

---

## 7. 测试覆盖率增强

### 7.1 更精确的覆盖率

**Go 1.23覆盖率改进**:

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out

# Go 1.23：更精确的覆盖率分析
go tool cover -func=coverage.out

# 输出示例：
# package/file.go:10:    FunctionA    100.0%
# package/file.go:20:    FunctionB     75.0%
# package/file.go:30:    FunctionC     50.0%
# total:                (statements)   80.0%
```

### 7.2 函数级覆盖率

**查看未覆盖的函数**:

```bash
# 显示未覆盖的函数
go tool cover -func=coverage.out | grep "0.0%"

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html
```

### 7.3 HTML报告改进

**Go 1.23的HTML报告增强**:

```html
<!-- 改进的特性 -->
<!-- 1. 更好的颜色对比 -->
<!-- 2. 行号导航 -->
<!-- 3. 函数跳转 -->
<!-- 4. 覆盖率百分比显示 -->
```

**生成并查看**:

```bash
go test -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out

# 在浏览器中自动打开
```

---

## 8. 测试工具函数

### 8.1 t.TempDir()最佳实践

**使用临时目录**:

```go
func TestFileOperations(t *testing.T) {
    // 自动清理的临时目录
    dir := t.TempDir()
    
    // 创建测试文件
    testFile := filepath.Join(dir, "test.txt")
    if err := os.WriteFile(testFile, []byte("test data"), 0644); err != nil {
        t.Fatal(err)
    }
    
    // 测试代码
    result, err := ProcessFile(testFile)
    if err != nil {
        t.Errorf("ProcessFile failed: %v", err)
    }
    
    // 不需要手动清理，t.TempDir()会自动处理
}
```

**并发测试中的TempDir**:

```go
func TestConcurrentFileOps(t *testing.T) {
    for i := 0; i < 10; i++ {
        i := i
        t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
            t.Parallel()
            
            // 每个并发测试有独立的临时目录
            dir := t.TempDir()
            
            // 测试代码
            testFile := filepath.Join(dir, "data.txt")
            // ...
        })
    }
}
```

### 8.2 t.Setenv()使用

**安全的环境变量设置**:

```go
func TestEnvironmentDependentCode(t *testing.T) {
    // t.Setenv()会在测试结束后自动恢复
    t.Setenv("API_KEY", "test-key-123")
    t.Setenv("DEBUG", "true")
    
    // 测试使用环境变量的代码
    client := NewAPIClient()  // 读取API_KEY
    result, err := client.FetchData()
    
    if err != nil {
        t.Errorf("FetchData failed: %v", err)
    }
    
    // 环境变量会自动恢复
}
```

**并行测试注意事项**:

```go
func TestParallelWithEnv(t *testing.T) {
    tests := []struct {
        name   string
        envVar string
        value  string
    }{
        {"test1", "VAR1", "value1"},
        {"test2", "VAR2", "value2"},
    }
    
    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            // ⚠️ t.Parallel() 和 t.Setenv() 要注意隔离
            t.Parallel()
            
            // 如果多个测试设置相同的环境变量，可能有问题
            // 最好使用不同的变量或避免并行
            t.Setenv(tt.envVar, tt.value)
            
            // 测试代码
        })
    }
}
```

### 8.3 t.Cleanup()模式

**注册清理函数**:

```go
func TestWithCleanup(t *testing.T) {
    // 创建资源
    db, err := sql.Open("postgres", testDSN)
    if err != nil {
        t.Fatal(err)
    }
    
    // 注册清理（类似defer，但更灵活）
    t.Cleanup(func() {
        db.Close()
    })
    
    // 创建更多资源
    conn, err := db.Conn(context.Background())
    if err != nil {
        t.Fatal(err)
    }
    
    // 再次注册清理（LIFO顺序）
    t.Cleanup(func() {
        conn.Close()
    })
    
    // 测试代码
    // ...
    
    // 清理会自动按LIFO顺序执行
}
```

**辅助函数中的cleanup**:

```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", testDSN)
    if err != nil {
        t.Fatal(err)
    }
    
    // 在辅助函数中注册清理
    t.Cleanup(func() {
        db.Close()
    })
    
    // 初始化数据库
    if err := initSchema(db); err != nil {
        t.Fatal(err)
    }
    
    return db
}

func TestDatabase(t *testing.T) {
    db := setupTestDB(t)
    // 使用db，不需要手动清理
}
```

---

## 9. 实战案例

### 9.1 完整的日志Handler测试

```go
package customlog_test

import (
    "bytes"
    "context"
    "encoding/json"
    "log/slog"
    "testing"
    "testing/slogtest"
    "time"
)

// JSONHandler自定义JSON处理器
type JSONHandler struct {
    buf   *bytes.Buffer
    opts  *slog.HandlerOptions
    attrs []slog.Attr
    group string
}

func NewJSONHandler(buf *bytes.Buffer, opts *slog.HandlerOptions) *JSONHandler {
    if opts == nil {
        opts = &slog.HandlerOptions{}
    }
    return &JSONHandler{
        buf:  buf,
        opts: opts,
    }
}

func (h *JSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
    minLevel := slog.LevelInfo
    if h.opts.Level != nil {
        minLevel = h.opts.Level.Level()
    }
    return level >= minLevel
}

func (h *JSONHandler) Handle(ctx context.Context, r slog.Record) error {
    entry := make(map[string]interface{})
    
    entry["time"] = r.Time.Format(time.RFC3339)
    entry["level"] = r.Level.String()
    entry["msg"] = r.Message
    
    // 添加handler的attrs
    for _, a := range h.attrs {
        entry[a.Key] = a.Value.Any()
    }
    
    // 添加record的attrs
    r.Attrs(func(a slog.Attr) bool {
        entry[a.Key] = a.Value.Any()
        return true
    })
    
    data, err := json.Marshal(entry)
    if err != nil {
        return err
    }
    
    h.buf.Write(data)
    h.buf.WriteByte('\n')
    return nil
}

func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    newHandler := *h
    newHandler.attrs = append([]slog.Attr{}, h.attrs...)
    newHandler.attrs = append(newHandler.attrs, attrs...)
    return &newHandler
}

func (h *JSONHandler) WithGroup(name string) slog.Handler {
    newHandler := *h
    newHandler.group = name
    return &newHandler
}

// 测试套件
func TestJSONHandler(t *testing.T) {
    var buf bytes.Buffer
    
    newHandler := func() slog.Handler {
        buf.Reset()
        return NewJSONHandler(&buf, nil)
    }
    
    // 运行标准测试
    slogtest.Run(t, newHandler, slogtest.All)
}

func TestJSONHandlerOutput(t *testing.T) {
    var buf bytes.Buffer
    h := NewJSONHandler(&buf, nil)
    logger := slog.New(h)
    
    logger.Info("test message",
        "key1", "value1",
        "key2", 42,
    )
    
    // 验证JSON输出
    var entry map[string]interface{}
    if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
        t.Fatalf("invalid JSON: %v", err)
    }
    
    if entry["msg"] != "test message" {
        t.Errorf("wrong message: %v", entry["msg"])
    }
    if entry["key1"] != "value1" {
        t.Errorf("wrong key1: %v", entry["key1"])
    }
    if entry["key2"] != float64(42) {  // JSON数字是float64
        t.Errorf("wrong key2: %v", entry["key2"])
    }
}

func TestJSONHandlerLevels(t *testing.T) {
    tests := []struct {
        name       string
        level      slog.Level
        shouldShow []slog.Level
        shouldHide []slog.Level
    }{
        {
            name:       "Debug",
            level:      slog.LevelDebug,
            shouldShow: []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError},
            shouldHide: []slog.Level{},
        },
        {
            name:       "Info",
            level:      slog.LevelInfo,
            shouldShow: []slog.Level{slog.LevelInfo, slog.LevelWarn, slog.LevelError},
            shouldHide: []slog.Level{slog.LevelDebug},
        },
        {
            name:       "Warn",
            level:      slog.LevelWarn,
            shouldShow: []slog.Level{slog.LevelWarn, slog.LevelError},
            shouldHide: []slog.Level{slog.LevelDebug, slog.LevelInfo},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var buf bytes.Buffer
            h := NewJSONHandler(&buf, &slog.HandlerOptions{
                Level: tt.level,
            })
            logger := slog.New(h)
            
            // 测试应该显示的级别
            for _, level := range tt.shouldShow {
                buf.Reset()
                logger.Log(context.Background(), level, "test")
                if buf.Len() == 0 {
                    t.Errorf("level %v should be shown", level)
                }
            }
            
            // 测试应该隐藏的级别
            for _, level := range tt.shouldHide {
                buf.Reset()
                logger.Log(context.Background(), level, "test")
                if buf.Len() > 0 {
                    t.Errorf("level %v should be hidden", level)
                }
            }
        })
    }
}
```

### 9.2 并发服务测试

```go
package server_test

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "sync"
    "testing"
    "time"
)

// Server并发HTTP服务器
type Server struct {
    mu    sync.RWMutex
    data  map[string]string
    calls int
}

func NewServer() *Server {
    return &Server{
        data: make(map[string]string),
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.mu.Lock()
    s.calls++
    s.mu.Unlock()
    
    switch r.Method {
    case http.MethodGet:
        s.handleGet(w, r)
    case http.MethodPost:
        s.handlePost(w, r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    
    s.mu.RLock()
    value, ok := s.data[key]
    s.mu.RUnlock()
    
    if !ok {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    
    w.Write([]byte(value))
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    value := r.URL.Query().Get("value")
    
    s.mu.Lock()
    s.data[key] = value
    s.mu.Unlock()
    
    w.WriteHeader(http.StatusCreated)
}

// 测试套件
func TestServerConcurrency(t *testing.T) {
    server := NewServer()
    ts := httptest.NewServer(server)
    defer ts.Close()
    
    // 并发写入
    t.Run("concurrent_writes", func(t *testing.T) {
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            go func(i int) {
                defer wg.Done()
                
                url := fmt.Sprintf("%s?key=key%d&value=value%d", ts.URL, i, i)
                req, _ := http.NewRequest(http.MethodPost, url, nil)
                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                    t.Errorf("request failed: %v", err)
                    return
                }
                resp.Body.Close()
                
                if resp.StatusCode != http.StatusCreated {
                    t.Errorf("unexpected status: %d", resp.StatusCode)
                }
            }(i)
        }
        wg.Wait()
    })
    
    // 验证数据
    t.Run("verify_writes", func(t *testing.T) {
        for i := 0; i < 100; i++ {
            url := fmt.Sprintf("%s?key=key%d", ts.URL, i)
            resp, err := http.Get(url)
            if err != nil {
                t.Errorf("GET failed: %v", err)
                continue
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != http.StatusOK {
                t.Errorf("key%d: unexpected status %d", i, resp.StatusCode)
            }
        }
    })
    
    // 并发读写
    t.Run("concurrent_read_write", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        
        var wg sync.WaitGroup
        
        // 启动读取goroutine
        for i := 0; i < 10; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for {
                    select {
                    case <-ctx.Done():
                        return
                    default:
                        url := fmt.Sprintf("%s?key=key0", ts.URL)
                        resp, _ := http.Get(url)
                        if resp != nil {
                            resp.Body.Close()
                        }
                    }
                }
            }()
        }
        
        // 启动写入goroutine
        for i := 0; i < 10; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for {
                    select {
                    case <-ctx.Done():
                        return
                    default:
                        url := fmt.Sprintf("%s?key=key0&value=updated", ts.URL)
                        req, _ := http.NewRequest(http.MethodPost, url, nil)
                        resp, _ := http.DefaultClient.Do(req)
                        if resp != nil {
                            resp.Body.Close()
                        }
                    }
                }
            }()
        }
        
        wg.Wait()
    })
}
```

### 9.3 性能基准测试套件

```go
package perf_test

import (
    "bytes"
    "encoding/json"
    "strings"
    "testing"
)

// 字符串拼接基准测试
func BenchmarkStringConcat(b *testing.B) {
    b.Run("plus_operator", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            s := "hello"
            s += " "
            s += "world"
            _ = s
        }
    })
    
    b.Run("sprintf", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            s := fmt.Sprintf("%s %s", "hello", "world")
            _ = s
        }
    })
    
    b.Run("strings_builder", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            var sb strings.Builder
            sb.WriteString("hello")
            sb.WriteString(" ")
            sb.WriteString("world")
            _ = sb.String()
        }
    })
    
    b.Run("bytes_buffer", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            var buf bytes.Buffer
            buf.WriteString("hello")
            buf.WriteString(" ")
            buf.WriteString("world")
            _ = buf.String()
        }
    })
}

// JSON序列化基准测试
func BenchmarkJSONMarshal(b *testing.B) {
    type Data struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
        Tags []string `json:"tags"`
    }
    
    data := Data{
        ID:   123,
        Name: "test",
        Tags: []string{"tag1", "tag2", "tag3"},
    }
    
    b.Run("marshal", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _, err := json.Marshal(data)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
    
    b.Run("marshal_indent", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            _, err := json.MarshalIndent(data, "", "  ")
            if err != nil {
                b.Fatal(err)
            }
        }
    })
    
    b.Run("encoder", func(b *testing.B) {
        b.ReportAllocs()
        var buf bytes.Buffer
        enc := json.NewEncoder(&buf)
        
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            buf.Reset()
            if err := enc.Encode(data); err != nil {
                b.Fatal(err)
            }
        }
    })
}

// 运行并保存结果
// go test -bench=. -benchmem -count=10 > bench.txt
// benchstat bench.txt
```

---

## 10. 最佳实践

### 10.1 测试组织

**按功能组织测试**:

```go
// user_test.go
package user_test

func TestUserCreation(t *testing.T) {
    t.Run("valid_user", func(t *testing.T) {
        // 测试有效用户创建
    })
    
    t.Run("invalid_email", func(t *testing.T) {
        // 测试无效邮箱
    })
    
    t.Run("duplicate_email", func(t *testing.T) {
        // 测试重复邮箱
    })
}

func TestUserUpdate(t *testing.T) {
    // 用户更新测试
}

func TestUserDelete(t *testing.T) {
    // 用户删除测试
}
```

### 10.2 测试命名

**清晰的测试命名**:

```go
// ✅ 好的命名
func TestUserCreate_WithValidData_Success(t *testing.T) {}
func TestUserCreate_WithInvalidEmail_ReturnsError(t *testing.T) {}
func TestUserCreate_WithDuplicateEmail_ReturnsConflictError(t *testing.T) {}

// ❌ 不好的命名
func TestUser1(t *testing.T) {}
func TestUser2(t *testing.T) {}
func TestUserFail(t *testing.T) {}
```

### 10.3 测试数据管理

**使用测试固件**:

```go
// testdata/
// ├── users.json
// ├── invalid_users.json
// └── test_config.yaml

func loadTestData(t *testing.T, filename string) []byte {
    t.Helper()
    
    data, err := os.ReadFile(filepath.Join("testdata", filename))
    if err != nil {
        t.Fatalf("failed to load test data: %v", err)
    }
    return data
}

func TestWithFixtures(t *testing.T) {
    data := loadTestData(t, "users.json")
    
    var users []User
    if err := json.Unmarshal(data, &users); err != nil {
        t.Fatal(err)
    }
    
    // 使用测试数据
    for _, user := range users {
        // 测试代码
    }
}
```

---

## 11. 参考资源

### 官方文档

- [testing Package](https://pkg.go.dev/testing)
- [testing/slogtest Package](https://pkg.go.dev/testing/slogtest)
- [Go 1.23 Release Notes - Testing](https://go.dev/doc/go1.23#testing)

### 测试工具

- [testify](https://github.com/stretchr/testify) - 测试断言库
- [gomock](https://github.com/golang/mock) - Mock框架
- [httptest](https://pkg.go.dev/net/http/httptest) - HTTP测试

### 博客文章

- [Go Blog - Testing](https://go.dev/blog/)
- [Advanced Testing in Go](https://about.sourcegraph.com/go/)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025-10-29  
**文档状态**: ✅ 完成  
**适用版本**: Go 1.23+

**贡献者**: 欢迎提交Issue和PR改进本文档
