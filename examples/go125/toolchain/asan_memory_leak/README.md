# AddressSanitizer (ASan) 内存泄漏检测示例

## 概述

本示例演示如何使用Go 1.23+的AddressSanitizer (ASan)功能检测内存泄漏和内存安全问题。

## ⚠️ 重要说明

**ASan需要CGO和C编译器支持**:

- 在Windows上需要: MinGW-w64或MSVC
- 在macOS上需要: Xcode Command Line Tools
- 在Linux上需要: gcc

如果您的环境中没有C编译器，可以使用我们提供的**模拟版本**来学习ASan的概念和使用方法。

## 运行方式

### 方式1: 使用CGO (需要C编译器)

```bash
# 使用ASan编译
go build -asan -o asan_example main.go

# 运行
./asan_example
```

### 方式2: 使用模拟版本 (无需C编译器)

```bash
# 编译模拟版本
go build -tags mock -o asan_mock main_mock.go

# 运行
./asan_mock
```

## 特性说明

### 1. 内存泄漏检测

检测未释放的内存分配:

```go
func memoryLeak() {
    // C内存分配但未释放
    ptr := C.malloc(C.size_t(1024))
    // 忘记 C.free(ptr)
}
```

### 2. 缓冲区溢出检测

检测越界访问:

```go
func bufferOverflow() {
    arr := (*[10]byte)(C.malloc(10))
    arr[10] = 42  // 越界访问！
}
```

### 3. Use-After-Free检测

检测释放后使用:

```go
func useAfterFree() {
    ptr := C.malloc(C.size_t(1024))
    C.free(ptr)
    // 使用已释放的内存
    *(*byte)(ptr) = 42  // 错误！
}
```

## ASan环境变量

```bash
# 控制ASan行为
export ASAN_OPTIONS="detect_leaks=1:halt_on_error=0"

# 查看更详细的错误信息
export ASAN_OPTIONS="verbosity=1"
```

## 性能影响

ASan会显著影响性能:

- **内存开销**: 2-3x
- **运行时间**: 2-5x
- **建议**: 仅在开发和测试环境使用

## 最佳实践

1. **CI/CD集成**: 在CI管道中添加ASan测试
2. **定期扫描**: 定期运行ASan检查内存问题
3. **修复问题**: 立即修复ASan发现的问题
4. **性能测试**: 不要在性能测试中启用ASan

## 参考资料

- [Go 1.23+ Release Notes - ASan](https://go.dev/doc/go1.23#asan)
- [AddressSanitizer Wiki](https://github.com/google/sanitizers/wiki/AddressSanitizer)
- [Go ASan Documentation](https://go.dev/doc/articles/asan)

## 故障排除

### 问题: "cgo: C compiler not found"

**解决方案**:

1. **Windows**: 安装MinGW-w64

   ```bash
   choco install mingw
   ```

2. **macOS**: 安装Xcode Command Line Tools

   ```bash
   xcode-select --install
   ```

3. **Linux**: 安装gcc

   ```bash
   sudo apt-get install build-essential  # Debian/Ubuntu
   sudo yum install gcc                  # CentOS/RHEL
   ```

4. **或者**: 使用模拟版本学习

   ```bash
   go build -tags mock -o asan_mock main_mock.go
   ```

---

**注意**: ASan是一个强大的内存安全工具，但它不能检测所有类型的内存问题。建议结合其他工具（如race detector、go vet）使用。
