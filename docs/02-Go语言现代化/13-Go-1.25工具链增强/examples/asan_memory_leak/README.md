# AddressSanitizer (ASan) 示例

> **Go 版本**: 1.25+  
> **平台**: Linux, macOS  
> **目的**: 演示如何使用 Go 1.25 的 AddressSanitizer 检测内存问题

---

## 📋 目录

- [快速开始](#快速开始)
- [示例文件](#示例文件)
- [编译和运行](#编译和运行)
- [预期输出](#预期输出)
- [常见错误类型](#常见错误类型)
- [故障排查](#故障排查)

---

## 快速开始

### 1. 编译有问题的示例

```bash
# 编译包含内存泄漏的示例
go build -asan -o leak leak_example.go

# 运行并查看 ASan 报告
./leak
```

### 2. 编译正确的示例

```bash
# 编译无内存问题的示例
go build -asan -o fixed fixed_example.go

# 运行 (不应该有 ASan 报告)
./fixed
```

---

## 示例文件

### `leak_example.go` - 有问题的代码

包含多种内存问题的示例:

1. **简单内存泄漏**: 分配内存但忘记释放
2. **字符串泄漏**: 返回 C 字符串但调用者忘记释放
3. **Use-After-Free**: 使用已释放的内存
4. **双重释放**: 释放同一块内存两次
5. **缓冲区溢出**: 写入超出分配的内存范围

### `fixed_example.go` - 正确的代码

演示正确的内存管理模式:

1. **使用 defer 释放**: 确保资源被释放
2. **NULL 检查**: 释放后设置为 NULL
3. **边界检查**: 使用安全的复制函数
4. **错误处理**: 错误路径也要清理资源
5. **批量处理**: 大量操作时的正确内存管理

---

## 编译和运行

### 编译选项

```bash
# 基本编译
go build -asan -o myapp main.go

# 编译所有文件
go build -asan ./...

# 运行测试
go test -asan ./...
```

### 环境变量配置

```bash
# 启用详细输出
export ASAN_OPTIONS='detect_leaks=1:log_path=./asan.log'

# 检测到错误时中止
export ASAN_OPTIONS='detect_leaks=1:abort_on_error=1'

# 禁用泄漏检测 (只检测其他错误)
export ASAN_OPTIONS='detect_leaks=0'
```

---

## 预期输出

### leak_example.go 输出

```text
=== Go 1.25 AddressSanitizer 示例 ===

1. 简单内存泄漏:
   ✅ 调用完成 (但有内存泄漏)

2. 字符串泄漏:
   创建字符串: Hello, World!
   ✅ 字符串创建完成 (但有内存泄漏)

3. Use-After-Free (已注释,取消注释会触发错误):
   ⚠️  已跳过 (会触发 ASan 错误)

4. 双重释放 (已注释,取消注释会触发错误):
   ⚠️  已跳过 (会触发 ASan 错误)

5. 缓冲区溢出 (已注释,取消注释会触发错误):
   ⚠️  已跳过 (会触发 ASan 错误)

=== 程序运行完成 ===

=================================================================
==12345==ERROR: LeakSanitizer: detected memory leaks

Direct leak of 1024 byte(s) in 1 object(s) allocated from:
    #0 0x... in malloc
    #1 0x... in simple_leak leak_example.go:7
    #2 0x... in main leak_example.go:55

Direct leak of 14 byte(s) in 1 object(s) allocated from:
    #0 0x... in malloc
    #1 0x... in create_string leak_example.go:13
    #2 0x... in main leak_example.go:61

SUMMARY: AddressSanitizer: 1038 byte(s) leaked in 2 allocation(s).
```

### fixed_example.go 输出

```text
=== Go 1.25 AddressSanitizer - 正确示例 ===

1. 正确的字符串管理:
   创建字符串: Hello, Safe World!
   ✅ 无内存泄漏

2. 正确的数据结构管理:
   Value: 42
   Value after free (安全): -1
   ✅ 无 Use-After-Free

3. 正确的缓冲区处理:
   复制结果: This is a very long string that is safely handled
   ✅ 无缓冲区溢出

4. 批量处理 (正确管理资源):
   处理 100 个项目
   ✅ 无累积泄漏

5. 错误处理模式:
   处理结果: Test Data
   ✅ 错误时也正确清理

=== 程序运行完成 ===

(无 ASan 错误报告)
```

---

## 常见错误类型

### 1. 内存泄漏 (Memory Leak)

**问题**:

```go
func leakExample() {
    cstr := C.CString("test")
    // 忘记: C.free(unsafe.Pointer(cstr))
}
```

**修复**:

```go
func fixedExample() {
    cstr := C.CString("test")
    defer C.free(unsafe.Pointer(cstr))  // ✅ 使用 defer
}
```

**ASan 输出**:

```text
Direct leak of 5 byte(s) in 1 object(s) allocated from:
    #0 in malloc
    #1 in _cgo_... 
    #2 in leakExample
```

---

### 2. Use-After-Free

**问题**:

```go
func useAfterFree() {
    data := C.create_data(42)
    C.free_data(data)
    value := C.get_value(data)  // ❌ 使用已释放的内存
}
```

**修复**:

```go
func fixed() {
    data := C.create_data(42)
    value := C.get_value(data)  // ✅ 使用前未释放
    C.free_data(data)
}
```

**ASan 输出**:

```text
ERROR: AddressSanitizer: heap-use-after-free
READ of size 4 at 0x... thread T0
    #0 in get_value
    #1 in useAfterFree
```

---

### 3. 缓冲区溢出 (Buffer Overflow)

**问题**:

```go
/*
void overflow() {
    char buf[10];
    strcpy(buf, "This is too long");  // ❌ 溢出
}
*/
```

**修复**:

```go
/*
void safe() {
    char buf[20];  // ✅ 足够大
    strncpy(buf, "This is safe", 19);
    buf[19] = '\0';
}
*/
```

**ASan 输出**:

```text
ERROR: AddressSanitizer: stack-buffer-overflow
WRITE of size 17 at 0x... thread T0
    #0 in strcpy
    #1 in overflow
```

---

### 4. 双重释放 (Double Free)

**问题**:

```go
func doubleFree() {
    ptr := C.malloc(100)
    C.free(ptr)
    C.free(ptr)  // ❌ 双重释放
}
```

**修复**:

```go
func fixed() {
    ptr := C.malloc(100)
    C.free(ptr)
    ptr = nil  // ✅ 设置为 nil 防止再次释放
}
```

**ASan 输出**:

```text
ERROR: AddressSanitizer: attempting double-free
    #0 in free
    #1 in doubleFree
```

---

## 故障排查

### 问题 1: 编译失败 "asan not supported"

**原因**: 平台不支持或 Go 版本过低

**解决**:

```bash
# 检查 Go 版本
go version  # 需要 1.25+

# 检查平台
uname -s  # 应该是 Linux 或 Darwin (macOS)

# 在 Windows 上需要 Clang
set CC=clang
go build -asan ./...
```

---

### 问题 2: 没有 ASan 报告

**可能原因**:

1. **没有内存问题**: 代码正确 ✅
2. **泄漏检测被禁用**: `ASAN_OPTIONS=detect_leaks=0`
3. **错误被抑制**: 检查 ASAN_OPTIONS 配置

**验证**:

```bash
# 确保泄漏检测启用
unset ASAN_OPTIONS
./myapp

# 或显式启用
ASAN_OPTIONS=detect_leaks=1 ./myapp
```

---

### 问题 3: CGO 编译问题

**错误**: `undefined reference to ...`

**解决**:

```bash
# 确保安装了 GCC/Clang
sudo apt-get install build-essential  # Ubuntu/Debian
sudo yum install gcc                   # CentOS/RHEL
brew install gcc                       # macOS

# 设置 CGO_ENABLED
export CGO_ENABLED=1

# 编译
go build -asan ./...
```

---

### 问题 4: 性能太慢

**原因**: ASan 有 ~2x 性能开销

**建议**:

1. **只在开发/测试环境使用**
2. **不要在性能测试中使用**
3. **禁用不需要的检查**:
   ```bash
   ASAN_OPTIONS=detect_leaks=0:fast_unwind_on_malloc=1 ./myapp
   ```

---

## 最佳实践

### 1. CI/CD 集成

```yaml
# .github/workflows/asan.yml
- name: Run ASan tests
  env:
    ASAN_OPTIONS: detect_leaks=1:abort_on_error=1
  run: go test -asan ./...
```

### 2. 本地开发

```bash
# 创建 Makefile
test-asan:
    @echo "Running ASan tests..."
    ASAN_OPTIONS=detect_leaks=1:log_path=./asan.log \
    go test -asan -v ./...
    @echo "Check asan.log for results"
```

### 3. CGO 内存管理

```go
// ✅ 推荐模式
func processData(data []byte) error {
    // 1. 转换为 C 类型
    cData := C.CBytes(data)
    defer C.free(cData)  // 立即 defer
    
    // 2. 使用 C 数据
    result := C.process(cData, C.int(len(data)))
    
    // 3. 检查结果
    if result != 0 {
        return fmt.Errorf("failed: %d", result)
    }
    
    return nil
}
```

---

## 相关资源

### 文档

- 📘 [Go ASan 技术文档](../01-go-build-asan内存泄漏检测.md)
- 📘 [AddressSanitizer Wiki](https://github.com/google/sanitizers/wiki/AddressSanitizer)
- 📘 [CGO Documentation](https://pkg.go.dev/cmd/cgo)

### 工具

- 🔧 [Valgrind](https://valgrind.org/) - 另一个内存检测工具
- 🔧 [Go Race Detector](https://go.dev/doc/articles/race_detector) - 数据竞争检测

---

**创建日期**: 2025年10月18日  
**更新日期**: 2025年10月18日  
**作者**: AI Assistant

---

<p align="center">
  <b>🔍 使用 ASan 让你的 CGO 代码更安全! 🛡️</b>
</p>

