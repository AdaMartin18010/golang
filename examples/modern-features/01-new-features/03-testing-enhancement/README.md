# Go 1.25.3 测试增强：基准测试最佳实践

<!-- TOC START -->
- [Go 1.25.3 测试增强：基准测试最佳实践](#go-1253-测试增强基准测试最佳实践)
  - [🎯 **核心概念**](#-核心概念)
  - [✨ **基准测试最佳实践**](#-基准测试最佳实践)
  - [💡 **最佳实践和使用场景**](#-最佳实践和使用场景)
  - [🚀 **运行示例**](#-运行示例)
    - [运行所有基准测试](#运行所有基准测试)
    - [运行特定基准测试](#运行特定基准测试)
    - [查看详细输出](#查看详细输出)
  - [📊 **总结**](#-总结)
<!-- TOC END -->

## 🎯 **核心概念**

Go 1.25.3 在测试框架方面继续优化，虽然移除了实验性的 `testing.B.Loop` API，但提供了更成熟的基准测试最佳实践。本指南将展示如何编写准确、高效的基准测试。

**注意**: Go 1.25 中移除了实验性的 `testing.Loop` API，我们使用传统的基准测试模式，但遵循最佳实践。

**基本语法**:

```go
func BenchmarkMyFunction(b *testing.B) {
    // 1. 在循环外执行所有昂贵的设置
    setupData := prepareExpensiveData()
    b.ResetTimer()  // 重置计时器，排除设置时间

    // 2. 使用传统的 for 循环执行基准测试
    for i := 0; i < b.N; i++ {
        // 3. 将要被测试的核心逻辑放在这里
        myFunction(setupData)
    }
}
```

## ✨ **基准测试最佳实践**

1. **清晰的职责分离**:
    - 将**"单次设置" (per-run setup)** 与 **"每次迭代的逻辑" (per-iteration logic)** 分离开来。所有昂贵的、只需执行一次的设置代码都放在循环之前，并在设置完成后调用 `b.ResetTimer()`。循环体内部只包含需要被精确计时的核心操作。

2. **正确使用计时器控制**:
    - `b.ResetTimer()`: 重置计时器，排除之前设置代码的时间
    - `b.StopTimer()`: 暂停计时器，用于不计入测试时间的操作
    - `b.StartTimer()`: 恢复计时器

3. **并行测试 (`b.RunParallel`)**:
    - 使用 `b.RunParallel` 进行并行基准测试，可以更好地利用多核 CPU
    - 使用 `sync.Pool` 来复用对象，减少分配开销

      ```go
      b.RunParallel(func(pb *testing.PB) {
          for pb.Next() {
              myFunction()
          }
      })
      ```

4. **处理"每次迭代的设置" (Per-iteration Setup)**:
    - 当每次迭代都需要新环境时，使用 `b.StopTimer()` 和 `b.StartTimer()` 来控制计时
    - 或者使用 `sync.Pool` 来复用对象

      ```go
      for i := 0; i < b.N; i++ {
          b.StopTimer()  // 暂停计时
          buffer := new(bytes.Buffer)
          b.StartTimer()  // 恢复计时

          myFunction(buffer)
      }
      ```

5. **内存分配跟踪**:
    - 使用 `b.ReportAllocs()` 来报告内存分配情况
    - 使用 `b.SetBytes()` 来设置每次操作处理的字节数

## 💡 **最佳实践和使用场景**

1. **简单场景**:
    - 对于没有复杂设置的基准测试，直接使用传统的 `for i := 0; i < b.N; i++` 循环。

2. **有昂贵设置的场景**:
    - 将所有昂贵的、与循环无关的设置代码放在循环之前
    - 在设置完成后调用 `b.ResetTimer()` 重置计时器

3. **每次迭代都需要新环境的场景**:
    - 使用 `b.StopTimer()` 和 `b.StartTimer()` 来控制计时
    - 或者使用 `sync.Pool` 来复用对象，减少分配开销

4. **并行测试场景**:
    - 使用 `b.RunParallel` 进行并行基准测试
    - 在并行测试中使用 `sync.Pool` 来管理资源

5. **内存分配跟踪**:
    - 使用 `b.ReportAllocs()` 来报告内存分配
    - 使用 `b.SetBytes()` 来设置每次操作处理的字节数

## 🚀 **运行示例**

### 运行所有基准测试

```bash
cd examples/modern-features/01-new-features/03-testing-enhancement
go test -bench=. -benchmem
```

### 运行特定基准测试

```bash
# 运行传统方式的基准测试
go test -bench=BenchmarkProcessData_Traditional

# 运行改进方式的基准测试
go test -bench=BenchmarkProcessData_Improved

# 运行并行测试
go test -bench=BenchmarkParallel
```

### 查看详细输出

```bash
# 显示内存分配信息
go test -bench=. -benchmem -v

# 设置基准测试运行时间
go test -bench=. -benchtime=5s
```

## 📊 **总结**

虽然 Go 1.25.3 移除了实验性的 `testing.B.Loop` API，但通过遵循最佳实践，我们仍然可以编写出准确、高效、易于维护的基准测试代码。关键要点：

1. ✅ **正确使用计时器控制** (`ResetTimer`, `StopTimer`, `StartTimer`)
2. ✅ **清晰的职责分离** (设置代码 vs 测试代码)
3. ✅ **合理使用并行测试** (`RunParallel`)
4. ✅ **内存分配跟踪** (`ReportAllocs`, `SetBytes`)
5. ✅ **对象复用** (`sync.Pool`)

这些实践确保了基准测试的准确性和可重复性，是编写高质量 Go 代码的重要组成部分。
