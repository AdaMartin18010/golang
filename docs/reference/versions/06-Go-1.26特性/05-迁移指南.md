# Go 1.26 迁移指南

> 从 Go 1.24/1.25 迁移到 Go 1.26 的完整指南。

---

## 快速迁移清单

- [ ] 更新 `go.mod` 中的 Go 版本
- [ ] 运行 `go mod tidy` 更新依赖
- [ ] 运行 `go fix ./...` 应用 Modernizers
- [ ] 运行测试验证行为不变
- [ ] 检查加密包随机参数变化
- [ ] 检查 GODEBUG 设置变化
- [ ] 验证 JPEG 编解码输出 (如果相关)

---

## 1. 基本迁移步骤

### 1.1 更新 Go 版本

```bash
# 安装 Go 1.26
# https://go.dev/dl/

# 验证安装
go version
go version go1.26.0 linux/amd64
```

### 1.2 更新 go.mod

```bash
# 在项目根目录
go mod edit -go=1.26

# 或者手动编辑 go.mod
module example.com/myproject

go 1.26  # 更新这一行

require (
    // ...
)
```

### 1.3 整理依赖

```bash
# 下载并整理依赖
go mod tidy

# 验证依赖
go mod verify
```

### 1.4 应用自动修复

```bash
# 预览将要应用的修复
go fix -n ./...

# 应用所有修复
go fix ./...

# 整理导入 (如果使用 goimports)
goimports -w .
```

### 1.5 运行测试

```bash
# 运行所有测试
go test ./...

# 运行竞态检测
go test -race ./...

# 运行基准测试验证性能
go test -bench=. -benchmem ./...
```

---

## 2. 破坏性变化

### 2.1 加密包随机参数 (重要)

Go 1.26 中，`crypto` 包相关函数的随机参数被**忽略**，总是使用安全的随机源。

```go
// 受影响的包和函数:
// - crypto/ecdh: GenerateKey
// - crypto/ecdsa: GenerateKey, SignASN1, Sign
// - crypto/ed25519: GenerateKey
// - crypto/rsa: GenerateKey, GenerateMultiPrimeKey, EncryptPKCS1v15
// - crypto/rand: Prime

// ========== Go 1.25 及之前 ==========
import "crypto/rand"

// 显式指定随机源
privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

// ========== Go 1.26 ==========
// 随机参数被忽略，使用内部安全随机源
privateKey, err := ecdsa.GenerateKey(elliptic.P256(), nil)  // nil 也可以
// 或者
privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)  // 仍然兼容

// 内部实现现在总是使用 crypto/rand 的安全随机源
```

**迁移建议**: 代码无需修改，但注意：

1. 不要依赖特定的随机源进行确定性测试
2. 如果需要确定性测试，使用 `testing/cryptotest`

```go
// 测试时使用确定性随机
import "testing/cryptotest"

func TestWithDeterministicRandom(t *testing.T) {
    // 设置全局随机源
    cryptotest.SetGlobalRandom(yourTestRand)
    defer cryptotest.SetGlobalRandom(nil)  // 恢复

    // 现在 crypto 函数会使用你的随机源
    key, _ := ecdsa.GenerateKey(elliptic.P256(), nil)
    // ... 测试
}
```

### 2.2 GODEBUG 设置变化

以下 GODEBUG 设置将在 **Go 1.27** 中移除，在 1.26 中仍可临时使用：

```bash
# 这些设置将在 Go 1.27 移除
tlsunsafeekm=1      # 不安全 EKM
tlsrsakex=1         # 旧版 RSA 密钥交换
tls10server=1       # 允许 TLS 1.0
tls3des=1           # 允许 3DES
x509keypairleaf=0   # 不填充 Leaf 字段
```

**迁移建议**: 不要依赖这些设置，准备迁移到新的行为。

```go
// 检查 TLS 配置
// 确保使用 TLS 1.2 或更高版本
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
    },
}
```

### 2.3 JPEG 编解码器

Go 1.26 使用了全新的 JPEG 编解码器：

- **更快**: 编码/解码速度提升
- **更准确**: 更好的颜色处理
- **输出变化**: 编码后的文件可能有略微不同

**迁移建议**:

```go
// 如果你的测试依赖特定的 JPEG 输出
// 需要更新期望值

// 检查输出是否有效，而非字节级比较
func validateJPEG(data []byte) error {
    img, err := jpeg.Decode(bytes.NewReader(data))
    if err != nil {
        return err
    }
    // 验证图像尺寸等属性
    _ = img.Bounds()
    return nil
}

// 而非:
// expectedHash := "abc123..."
// actualHash := sha256.Sum256(data)
// if actualHash != expectedHash { ... }
```

---

## 3. 推荐的代码现代化

### 3.1 使用新的语言特性

```go
// ========== new(expr) 语法 ==========

// 旧代码
age := 25
person := &Person{
    Name: "Alice",
    Age:  &age,
}

// 新代码 (Go 1.26)
person := &Person{
    Name: "Alice",
    Age:  new(int(25)),  // 更简洁！
}

// ========== 递归泛型约束 ==========

// 旧代码: 无法实现自引用约束
// type Node[T any] interface { ... }  // 不够精确

// 新代码 (Go 1.26)
type Node[T Node[T]] interface {
    Children() []T
    Value() int
}

func Traverse[T Node[T]](root T, visit func(int)) {
    // 类型安全的通用遍历
}
```

### 3.2 使用新的标准库 API

```go
// ========== errors.AsType ==========

// 旧代码
var myErr *MyError
if errors.As(err, &myErr) {
    handle(myErr)
}

// 新代码 (Go 1.26)
if myErr, ok := errors.AsType[*MyError](err); ok {
    handle(myErr)
}

// ========== bytes.Buffer.Peek ==========

// 旧代码: 需要手动管理位置
buf := bytes.NewBuffer(data)
peeked := buf.Bytes()[:10]  // 不安全！

// 新代码 (Go 1.26)
peeked := buf.Peek(10)  // 安全，不移动读指针

// ========== log/slog.NewMultiHandler ==========

// 旧代码: 需要手动实现
// 新代码 (Go 1.26)
handler := slog.NewMultiHandler(
    slog.NewJSONHandler(file, nil),
    slog.NewTextHandler(os.Stdout, nil),
)
logger := slog.New(handler)
```

### 3.3 应用 Modernizers

```bash
# 运行所有 modernizers
go fix ./...

# 预期会看到以下类型的修复:
# - 使用 slices.Contains 替代循环查找
# - 使用 slices.Sort 替代 sort.Slice
# - 使用 maps.Clone 替代手动复制
# - 使用 strings.CutPrefix/CutSuffix
# - 使用内置 min/max
# - 使用 errors.AsType
```

---

## 4. 不同场景的迁移

### 4.1 Web 服务

```go
// 检查 TLS 配置
server := &http.Server{
    Addr: ":443",
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,  // 确保使用 TLS 1.2+
        // 移除旧版加密套件
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        },
    },
}
```

### 4.2 加密应用

```go
// 更新 HPKE 使用 (新包)
import "crypto/hpke"

// 使用新的后量子安全算法
kem, _ := hpke.GetKEM(hpke.KEM_MLKEM768)
kdf, _ := hpke.GetKDF(hpke.KDF_HKDF_SHA384)
aead, _ := hpke.GetAEAD(hpke.AEAD_AES256GCM)

suite, _ := hpke.NewSuite(kem, kdf, aead)
// ... 使用 suite 进行加密
```

### 4.3 CLI 工具

```bash
# 更新构建脚本
# 启用实验性特性 (如果需要)
export GOEXPERIMENT=simd
export GOEXPERIMENT=runtimesecret

go build -o myapp .
```

---

## 5. 故障排除

### 5.1 编译错误

```bash
# 问题: 类型参数循环引用
# 解决: 确保正确使用 Go 1.26 的递归约束语法

# 问题: 找不到实验性包
# 解决: 确认设置了 GOEXPERIMENT
GOEXPERIMENT=simd go build .

# 问题: 导入错误
# 解决: 运行 go mod tidy
go mod tidy
```

### 5.2 运行时问题

```bash
# 问题: GC 性能异常
# 解决: 检查 GC 统计
go test -bench=BenchmarkGC -gcflags="-m" ./...

# 问题: 切片栈分配导致问题
# 解决: 临时禁用优化
go build -gcflags=all=-d=variablemakehash=n ./...
```

### 5.3 测试失败

```bash
# 问题: 加密测试失败 (随机性)
# 解决: 使用 testing/cryptotest

# 问题: JPEG 测试失败 (输出变化)
# 解决: 更新期望值或使用图像属性验证

# 问题: 竞态检测失败
# 解决: 检查新的 crypto 包随机行为
```

---

## 6. 版本兼容性

### 6.1 Go 版本要求

```
项目版本    要求 Go 版本
────────    ────────────
Go 1.26     Go 1.24.6+ (引导)
Go 1.25     Go 1.22.0+ (引导)
Go 1.24     Go 1.20.0+ (引导)
```

### 6.2 模块兼容性

```go
// go.mod 中的 go 指令控制语言版本
module example.com/mymod

go 1.25  // 代码使用 Go 1.25 语言特性
toolchain go1.26.0  // 使用 Go 1.26 工具链构建

// 使用 //go:build 约束
//go:build go1.26
// +build go1.26

package mypkg

// 这段代码只在 Go 1.26+ 编译
```

---

## 7. 参考资源

- [Go 1.26 Release Notes](https://go.dev/doc/go1.26)
- [Go 1.26 迁移工具](https://go.dev/dl/)
- [Go Modules 参考](https://go.dev/ref/mod)
- [Go 兼容性承诺](https://go.dev/doc/go1compat)

---

## 迁移时间线建议

| 阶段 | 时间 | 任务 |
|------|------|------|
| 准备 | Day 1 | 备份代码，阅读 release notes |
| 升级 | Day 1-2 | 更新 Go 版本，运行 go fix |
| 验证 | Day 2-3 | 运行测试，检查性能 |
| 优化 | Day 3-5 | 应用新特性，优化代码 |
| 部署 | Day 5+ | 分阶段部署，监控指标 |
