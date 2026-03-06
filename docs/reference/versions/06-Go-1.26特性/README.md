# Go 1.26 特性完全指南

> Go 1.26 于 2026年2月 发布，带来了语言层面的重要改进、性能优化和标准库增强。

---

## 版本概览

| 特性类别 | 主要变化 |
|---------|---------|
| **语言变化** | `new()` 支持表达式、递归泛型约束 |
| **运行时** | Green Tea GC 默认启用、cgo 开销减少30% |
| **标准库** | 新增 `crypto/hpke`、`simd/archsimd`(实验性)、`runtime/secret`(实验性) |
| **工具链** | `go fix` 完全重写、Modernizers 自动化代码更新 |
| **编译器** | 切片栈分配优化、更多逃逸分析改进 |

---

## 快速导航

| 文档 | 内容 |
|------|------|
| [01-语言变化详解](./01-语言变化详解.md) | `new(expr)` 语法、递归泛型约束 |
| [02-标准库新增](./02-标准库新增.md) | crypto/hpke、simd/archsimd、runtime/secret |
| [03-工具改进](./03-工具改进.md) | go fix Modernizers、//go:fix inline |
| [04-性能优化](./04-性能优化.md) | Green Tea GC、cgo 优化、编译器改进 |
| [05-迁移指南](./05-迁移指南.md) | 从 Go 1.25/1.24 迁移到 1.26 |

---

## 语言变化详解

### 1. new() 支持表达式操作数

Go 1.26 中，内置的 `new` 函数现在允许其操作数为表达式，用于指定变量的初始值：

```go
// Go 1.25 及之前
x := int64(300)
ptr := &x

// Go 1.26 新方式
ptr := new(int64(300))  // 更简洁！
```

这在处理可选字段时特别有用：

```go
type Person struct {
    Name string `json:"name"`
    Age  *int   `json:"age"` // 可选字段
}

// 简洁地创建带有可选字段的结构
func createPerson(name string, born time.Time) ([]byte, error) {
    return json.Marshal(Person{
        Name: name,
        Age:  new(yearsSince(born)),  // Go 1.26!
    })
}
```

### 2. 递归泛型约束

泛型类型现在可以在自己的类型参数列表中引用自己：

```go
// Go 1.26 之前: 编译错误
// type Adder[A Adder[A]] interface { ... }

// Go 1.26: 合法！
type Adder[A Adder[A]] interface {
    Add(A) A
}

func algo[A Adder[A]](x, y A) A {
    return x.Add(y)
}
```

这在实现自引用数据结构（如树、图）时非常有用。

---

## 性能改进

### Green Tea GC

Go 1.26 默认启用之前实验性的 Green Tea 垃圾回收器：

- **更低延迟**: 10-40% 的延迟改善
- **更高吞吐**: 更好的内存利用效率
- **自动启用**: 无需配置，开箱即用

```go
// 监控 GC 性能
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("GC 次数: %d\n", m.NumGC)
fmt.Printf("GC 暂停: %d ns\n", m.PauseNs[(m.NumGC+255)%256])
```

### cgo 开销优化

cgo 的基础开销减少了约 30%，使得 Go 与 C 代码的互操作更加高效。

---

## 标准库新增

### crypto/hpke

实现了 RFC 9180 标准的混合公钥加密 (HPKE)，支持后量子混合 KEM：

```go
import "crypto/hpke"

// 使用 HPKE 进行加密
 kem, _ := hpke.GetKEM(hpke.KEM_P384_HKDF_SHA384)
 kdf, _ := hpke.GetKDF(hpke.KDF_HKDF_SHA384)
 aead, _ := hpke.GetAEAD(hpke.AEAD_AES256GCM)

 suite, _ := hpke.NewSuite(kem, kdf, aead)
 // ... 使用 suite 进行加密操作
```

### simd/archsimd (实验性)

通过 `GOEXPERIMENT=simd` 启用，提供架构特定的 SIMD 操作：

```go
// 需要设置 GOEXPERIMENT=simd
import "simd/archsimd"

// 256位整数向量
v1 := archsimd.Int8x32{...}
v2 := archsimd.Int8x32{...}
result := v1.Add(v2)  // SIMD 加法
```

**注意**: 该 API 目前不稳定，仅支持 amd64 架构。

### runtime/secret (实验性)

通过 `GOEXPERIMENT=runtimesecret` 启用，提供安全的临时数据擦除：

```go
// 需要设置 GOEXPERIMENT=runtimesecret
import "runtime/secret"

secret.WithSecrets(func() {
    // 在此函数内处理的敏感数据
    // 函数返回后会被安全擦除
})
```

---

## 工具改进

### 重写的 go fix

`go fix` 现在使用 Go 分析框架，提供可靠的自动化代码现代化：

```bash
# 运行所有 modernizers
go fix ./...

# 查看具体变化 (dry-run)
go fix -n ./...
```

包含的特性：

- **Modernizers**: 自动使用新语言和库特性
- **Inline 分析器**: 自动内联带有 `//go:fix inline` 注解的函数

### 使用 //go:fix inline

```go
//go:fix inline
func OldAPI() {
    NewAPI()
}

func NewAPI() {
    // 新实现
}
```

运行 `go fix` 后，所有对 `OldAPI()` 的调用都会被替换为 `NewAPI()`。

---

## 迁移指南

### 从 Go 1.25 迁移

Go 1.26 保持向后兼容，大多数项目可以直接升级：

```bash
# 更新 go.mod
go mod edit -go=1.26

# 下载依赖
go mod tidy

# 运行测试
go test ./...

# 使用 go fix 自动更新代码
go fix ./...
```

### 注意事项

1. **加密包随机参数**: `crypto/ecdh`、`crypto/ecdsa` 等包的随机参数现在被忽略，使用安全随机源
2. **GODEBUG 设置**: 一些旧的 GODEBUG 设置将在 Go 1.27 中移除
3. **JPEG 编解码器**: 新的实现可能有略微不同的输出

---

## 版本信息

- **发布日期**: 2026年2月
- **引导要求**: Go 1.24.6 或更高版本
- **支持平台**: 所有 Go 1.25 支持的平台
- **Windows/arm64**: 支持 cgo 内部链接模式

---

## 参考资源

- [Go 1.26 Release Notes](https://go.dev/doc/go1.26)
- [Go Blog: Go 1.26 is released](https://go.dev/blog/go1.26)
- [Effective Go](https://go.dev/doc/effective_go)
