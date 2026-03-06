# Go 1.26 快速参考卡片

> **文档层级**: T-工具层 (Tools Layer)
> **文档类型**: 快速参考 (Quick Reference)
> **适用场景**: 开发时快速查阅
> **最后更新**: 2026-03-06

---

## 一、new表达式速查

### 1.1 语法对比

| 场景 | Go ≤1.25 | Go 1.26 |
|------|----------|---------|
| 基本使用 | `ptr := &T{v}` | `ptr := new(T(v))` ✅ |
| 延迟初始化 | 需要中间变量 | `new(expr)`直接赋值 |
| 可选字段 | `&[]T{v}[0]` | `new(v)` ✅ |

### 1.2 常用模式

```go
// 可选字段配置
type Config struct {
    Timeout *int
}
cfg := Config{Timeout: new(30)}  // Go 1.26 ✨

// 构造者模式
builder.WithPort(new(8080))       // Go 1.26 ✨

// 延迟初始化
var ptr *int
ptr = new(42)                     // Go 1.26 ✨
```

---

## 二、递归泛型速查

### 2.1 基本语法

```go
// 递归约束定义
type Adder[A Adder[A]] interface {
    Add(A) A
}

// 树节点
type Node[T Node[T]] struct {
    Value T
    Children []T
}
```

### 2.2 典型应用

```go
// 树遍历
type Tree[T Tree[T]] interface {
    Children() []T
}

func Walk[T Tree[T]](root T, fn func(T)) {
    fn(root)
    for _, child := range root.Children() {
        Walk(child, fn)  // 递归调用
    }
}
```

---

## 三、GreenTeaGC速查

### 3.1 环境变量

| 变量 | 说明 | 推荐值 |
|------|------|--------|
| `GOGC` | GC触发百分比 | 100 (默认) |
| `GOMEMLIMIT` | 内存限制 | 根据容器限制设置 |
| `GOMAXPROCS` | 并行度 | 默认值即可 |

### 3.2 性能调优

```go
import "runtime/debug"

// 设置内存限制
debug.SetMemoryLimit(8 << 30)  // 8GB

// 设置GC目标百分比
debug.SetGCPercent(100)
```

### 3.3 监控指标

```go
import "runtime/metrics"

// 读取GC暂停时间
samples := []metrics.Sample{
    {Name: "/gc/pause:seconds"},
}
metrics.Read(samples)
```

---

## 四、HPKE加密速查

### 4.1 基础封装

```go
import "crypto/hpke"

// 生成密钥对
pubKey, privKey, err := hpke.GenerateKeyPair(hpke.DHKEM_X25519_HKDF_SHA256)

// 封装消息
encapsulatedKey, ciphertext, err := hpke.SingleShotSeal(
    hpke.DHKEM_X25519_HKDF_SHA256,
    hpke.HKDF_SHA256,
    hpke.AES_128_GCM,
    pubKey,
    plaintext,
    nil, // info
)

// 解封装
plaintext, err := hpke.SingleShotOpen(
    hpke.DHKEM_X25519_HKDF_SHA256,
    hpke.HKDF_SHA256,
    hpke.AES_128_GCM,
    privKey,
    encapsulatedKey,
    ciphertext,
    nil, // info
)
```

### 4.2 算法标识符

| KEM | KDF | AEAD |
|-----|-----|------|
| `DHKEM_X25519_HKDF_SHA256` | `HKDF_SHA256` | `AES_128_GCM` |
| `DHKEM_P256_HKDF_SHA256` | `HKDF_SHA384` | `AES_256_GCM` |
| `DHKEM_X448_HKDF_SHA512` | `HKDF_SHA512` | `ChaCha20Poly1305` |

---

## 五、SIMD加速速查

### 5.1 自动加速的包

```go
import (
    "bytes"   // SIMD加速的字符串操作
    "strings" // SIMD加速的搜索
)

// 以下操作自动使用SIMD（如果CPU支持）
bytes.Equal(a, b)
bytes.Index(haystack, needle)
strings.Index(haystack, needle)
```

### 5.2 性能提示

- ✅ 大批量数据处理才能体现优势
- ✅ 数据对齐可提高性能
- ❌ 小数据量可能无收益
- ❌ 不要手动调用SIMD（由runtime自动处理）

---

## 六、工具链速查

### 6.1 go fix现代化

```bash
# 自动现代化代码
go fix ./...

# 查看具体变更
go fix -diff ./...
```

### 6.2 编译优化

```bash
# 启用所有优化
go build -ldflags="-s -w" -gcflags="-B" .

# 查看逃逸分析
go build -gcflags="-m" .
```

---

## 七、兼容性速查

### 7.1 版本要求

| 特性 | 最低版本 | 向后兼容 |
|------|----------|----------|
| new(expr) | Go 1.26 | ❌ |
| 递归泛型 | Go 1.26 | ❌ |
| GreenTeaGC | Go 1.26 | ✅ (自动) |
| crypto/hpke | Go 1.26 | N/A |
| SIMD加速 | Go 1.26 | ✅ (自动) |

### 7.2 升级检查清单

```markdown
□ 更新go.mod: go 1.26
□ 运行go fix现代化
□ 测试new表达式使用
□ 测试递归泛型代码
□ 验证GC性能变化
□ 更新CI/CD环境
```

---

## 八、常见错误

### 8.1 new表达式常见错误

```go
// ❌ 错误：不能用于非值类型
ptr := new(make([]int, 10))

// ✅ 正确：只能用于值类型
ptr := new([]int{1, 2, 3})
```

### 8.2 递归泛型常见错误

```go
// ❌ 错误：递归无限展开
type Bad[T Bad[Bad[T]]] struct{}

// ✅ 正确：结构递归
type Good[T Good[T]] struct {
    Child *T
}
```

---

## 九、速查命令

```bash
# 查看Go版本
go version

# 查看1.26特性详情
go doc go1.26

# 运行测试
go test -v ./...

# 性能测试
go test -bench=. -benchmem ./...

# 查看GC统计
GODEBUG=gctrace=1 go run main.go
```

---

**打印建议**: 将此页面打印或保存为PDF，方便开发时快速查阅。
