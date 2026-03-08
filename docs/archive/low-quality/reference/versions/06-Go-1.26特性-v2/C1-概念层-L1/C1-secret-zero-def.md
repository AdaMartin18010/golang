# crypto/subtle - 秘密清零

> **文档层级**: C1-概念层 (Concept Layer L1)
> **文档类型**: 概念定义 (Concept Definition)
> **包路径**: `crypto/subtle`
> **最后更新**: 2026-03-06

---

## 一、概念定义

### 1.1 安全擦除

```
秘密清零 (Secure Zeroing):
  安全地从内存中擦除敏感数据，防止数据残留

威胁模型:
  - 内存转储攻击
  - 冷启动攻击
  - 交换分区泄露
  - 核心转储泄露

实现机制:
  - 使用volatile阻止编译器优化
  - 内存屏障确保写入执行
  - 覆盖后再释放
```

### 1.2 API设计

```go
package subtle

// ZeroBytes 安全清零字节切片
func ZeroBytes(b []byte)

// ZeroUint32 安全清零uint32数组
func ZeroUint32(x []uint32)

// ZeroUint64 安全清零uint64数组
func ZeroUint64(x []uint64)
```

---

## 二、使用场景

### 2.1 密钥管理

```go
func processKey(key []byte) {
    defer subtle.ZeroBytes(key)  // 确保密钥被清除

    // 使用密钥...
    result := encrypt(data, key)

    // 函数返回时密钥被安全清除
}
```

### 2.2 密码处理

```go
func checkPassword(input []byte) bool {
    // 复制输入以防止修改原始数据
    inputCopy := make([]byte, len(input))
    copy(inputCopy, input)
    defer subtle.ZeroBytes(inputCopy)

    // 验证密码...
    return verify(inputCopy, storedHash)
}
```

---

## 三、安全保证

```
性质:
  1. 编译器不会优化掉清零操作
  2. 内存实际被覆盖
  3. 不会被交换到磁盘（结合mlock）

限制:
  - 不防止硬件级别的内存扫描
  - 不防止操作系统层面的内存检查
```

---

**概念分类**: 标准库 - 安全工具
**Go版本**: 1.26+
**包路径**: `crypto/subtle`
