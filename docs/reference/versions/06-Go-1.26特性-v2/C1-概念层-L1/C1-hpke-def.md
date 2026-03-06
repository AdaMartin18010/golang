# HPKE (混合公钥加密)

> **文档层级**: C1-概念层 (Concept Layer L1)
> **文档类型**: 概念定义 (Concept Definition)
> **标准**: RFC 9180
> **最后更新**: 2026-03-06

---

## 一、概念定义

### 1.1 形式化定义

```
HPKE (Hybrid Public Key Encryption) :
  RFC 9180标准的混合公钥加密方案

组成组件:
  - KEM: 密钥封装机制 (Key Encapsulation Mechanism)
  - KDF: 密钥派生函数 (Key Derivation Function)
  - AEAD: 认证加密算法 (Authenticated Encryption with Associated Data)

加密流程:
  1. KEM.Encap(pk) → (enc, shared_secret)
  2. KDF(shared_secret, info) → key
  3. AEAD.Seal(key, plaintext, aad) → ciphertext
```

### 1.2 设计目标

| 目标 | 说明 |
|------|------|
| 前向保密 | 长期密钥泄露不影响历史消息 |
| 后向保密 | 当前密钥泄露不影响未来消息（部分模式） |
| 认证性 | 可选的发送方认证 |
| 互操作性 | 标准算法套件，跨语言兼容 |

### 1.3 算法套件

```go
// 推荐的算法组合
const (
    // KEM
    DHKEM_X25519_HKDF_SHA256  = 0x0020
    DHKEM_P256_HKDF_SHA256    = 0x0010
    DHKEM_X448_HKDF_SHA512    = 0x0021

    // KDF
    HKDF_SHA256 = 0x0001
    HKDF_SHA384 = 0x0002
    HKDF_SHA512 = 0x0003

    // AEAD
    AES_128_GCM         = 0x0001
    AES_256_GCM         = 0x0002
    ChaCha20Poly1305    = 0x0003
)
```

---

## 二、操作模式

### 2.1 Base模式

```
基础模式：仅使用接收方公钥加密

特点:
  - 匿名发送方
  - 前向保密
  - 适用于广播场景

流程:
  Sender: enc, ct = Seal(pkR, pt, info)
  Receiver: pt = Open(skR, enc, ct, info)
```

### 2.2 PSK模式

```
预共享密钥模式：增加预共享密钥认证

特点:
  - 发送方和接收方共享PSK
  - 提供额外的认证层
  - 后向保密（如果PSK定期轮换）
```

### 2.3 Auth模式

```
认证模式：发送方使用自身私钥认证

特点:
  - 发送方身份可验证
  - 非否认性（发送方不可否认发送）

流程:
  Sender: enc, ct = AuthSeal(pkR, skS, pt, info)
  Receiver: pt = AuthOpen(skR, pkS, enc, ct, info)
```

### 2.4 AuthPSK模式

```
认证+PSK模式：结合Auth和PSK

特点:
  - 最强安全性保证
  - 同时使用私钥和PSK认证
  - 适用于高安全要求场景
```

---

## 三、Go API设计

### 3.1 简化API (Single-Shot)

```go
package main

import (
    "crypto/hpke"
)

func main() {
    // 生成密钥对
    pubKey, privKey, err := hpke.GenerateKeyPair(hpke.DHKEM_X25519_HKDF_SHA256)

    // 发送方：封装消息
    encapsulatedKey, ciphertext, err := hpke.SingleShotSeal(
        hpke.DHKEM_X25519_HKDF_SHA256,  // KEM
        hpke.HKDF_SHA256,                // KDF
        hpke.AES_256_GCM,                // AEAD
        pubKey,                          // 接收方公钥
        plaintext,                       // 明文
        []byte("application context"),   // info（可选上下文）
    )

    // 接收方：解封装
    plaintext, err := hpke.SingleShotOpen(
        hpke.DHKEM_X25519_HKDF_SHA256,
        hpke.HKDF_SHA256,
        hpke.AES_256_GCM,
        privKey,          // 接收方私钥
        encapsulatedKey,  // 封装的密钥
        ciphertext,       // 密文
        []byte("application context"),
    )
}
```

### 3.2 流式API (Sender/Receiver)

```go
// 流式加密（大数据量）
func streamEncrypt(pubKey []byte, reader io.Reader, writer io.Writer) error {
    sender, err := hpke.NewSender(
        hpke.DHKEM_X25519_HKDF_SHA256,
        hpke.HKDF_SHA256,
        hpke.AES_256_GCM,
        pubKey,
        []byte("stream context"),
    )

    encapsulatedKey, err := sender.Encapsulate()
    // 发送encapsulatedKey给接收方...

    // 流式加密
    stream := sender.NewStream()
    _, err = io.Copy(stream, reader)
    return err
}
```

---

## 四、安全特性

### 4.1 安全性分析 (Th3.1)

```
定理 Th3.1: HPKE满足IND-CCA2安全性

前提假设:
  - KEM是IND-CCA安全的
  - KDF是安全的密钥派生函数
  - AEAD是IND-CCA2安全的

证明概要:
  由组合定理，安全的KEM+KDF+AEAD组合提供IND-CCA2安全性。
```

### 4.2 安全保证

| 威胁模型 | Base | PSK | Auth | AuthPSK |
|----------|------|-----|------|---------|
| 被动窃听 | ✅ 安全 | ✅ 安全 | ✅ 安全 | ✅ 安全 |
| 主动攻击 | ✅ 安全 | ✅ 安全 | ✅ 安全 | ✅ 安全 |
| 长期密钥泄露 | ❌ 前向保密失效 | ❌ | ❌ | ❌ |
| 发送方冒充 | ❌ | ✅ PSK认证 | ✅ 签名认证 | ✅ 双认证 |

---

## 五、应用场景

### 5.1 TLS替代

```
场景: 需要加密通信但不使用TLS协议栈
示例: IoT设备、内部服务通信
```

### 5.2 消息加密

```
场景: 加密消息队列、事件总线
特点: 接收方可能有多个，使用Base模式广播
```

### 5.3 密钥封装

```
场景: 封装对称密钥进行存储或传输
优势: 比直接使用RSA加密更安全
```

### 5.4 混合加密系统

```
场景: 文件加密、数据库加密
模式: 用HPKE封装文件密钥，用AEAD加密文件内容
```

---

## 六、相关文档

- **标准**: RFC 9180
- **应用**: [C3-安全通信模式](../C3-实践层-L3/C3-安全通信模式.md)
- **定理**: [Th3.1](../R-参考层/R-定理索引.md#Th3.1)

---

**概念分类**: 标准库 - 加密安全
**Go版本**: 1.26+
**包路径**: `crypto/hpke`
