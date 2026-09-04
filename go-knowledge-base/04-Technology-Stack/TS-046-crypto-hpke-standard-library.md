# Go 1.26 crypto/hpke 标准库深度解析

> **文档编号**: TS-046
> **分类**: 04-Technology-Stack / Security / Cryptography
> **Go 版本**: 1.26+
> **最后更新**: 2026-05-06

> **维度**: Technology Stack
> **级别**: B (3 KB)
> **标签**: #ts
---

## 一、概述

`crypto/hpke` 是 Go 1.26 引入的标准库包，实现 **RFC 9180** 定义的 Hybrid Public Key Encryption（HPKE）。

HPKE 结合非对称密钥封装（KEM）和对称加密（AEAD），提供高效的公钥加密，特别适用于：

- TLS 1.3 扩展（ECH, Encrypted Client Hello）
- 消息加密系统（如 MLS, Messaging Layer Security）
- 后量子安全通信（通过 ML-KEM 混合 KEM）

---

## 二、核心组件

HPKE 由三个独立可选的算法组件构成：

```text
HPKE 套件 (ciphersuite):
────────────────────────────────────────
KEM  (Key Encapsulation Mechanism):  DHKEM, MLKEM768, MLKEM1024
KDF  (Key Derivation Function):      HKDF-SHA256, HKDF-SHA384, HKDF-SHA512
AEAD (Authenticated Encryption):     AES-128-GCM, AES-256-GCM, ChaCha20Poly1305

示例套件:
├─ DHKEM(X25519, HKDF-SHA256) + HKDF-SHA256 + AES-128-GCM  (默认推荐)
├─ MLKEM768 + HKDF-SHA256 + AES-128-GCM                    (后量子)
└─ MLKEM768X25519 + HKDF-SHA256 + AES-128-GCM              (混合后量子)
```

---

## 三、Go 1.26 API

```go
package hpke // import "crypto/hpke"

// 高层 API: Seal / Open
func Seal(pk PublicKey, kdf KDF, aead AEAD, info, plaintext []byte) ([]byte, error)
func Open(k PrivateKey, kdf KDF, aead AEAD, info, ciphertext []byte) ([]byte, error)

// KEM 构造
func DHKEM(curve ecdh.Curve) KEM
func MLKEM768() KEM
func MLKEM768X25519() KEM

// KDF 构造
func HKDFSHA256() KDF
func HKDFSHA384() KDF
func HKDFSHA512() KDF

// AEAD 构造
func AES128GCM() AEAD
func AES256GCM() AEAD
func ChaCha20Poly1305() AEAD
```

---

## 四、使用模式

### 4.1 基础加密 (Base Mode)

```go
kem := hpke.DHKEM(ecdh.X25519())
kdf := hpke.HKDFSHA256()
aead := hpke.AES128GCM()

// 接收方生成密钥对
priv, _ := kem.GenerateKey()
pub := priv.PublicKey()

// 发送方加密
ciphertext, _ := hpke.Seal(pub, kdf, aead, info, plaintext)

// 接收方解密
plaintext, _ := hpke.Open(priv, kdf, aead, info, ciphertext)
```

### 4.2 后量子混合模式

```go
// MLKEM768X25519 结合经典 X25519 和后量子 ML-KEM-768
kem := hpke.MLKEM768X25519()
// 提供前向保密 + 后量子安全性
```

---

## 五、与 crypto/tls 集成

Go 1.26 在 `crypto/tls` 中默认启用后量子 TLS：

```text
TLS 1.3 密钥交换:
────────────────────────────────────────
SecP256r1MLKEM768  (默认启用)
SecP384r1MLKEM1024 (高安全场景)
```

---

## 六、交叉引用

- `examples/crypto-hpke/main.go` 可运行示例
- `SEC-2026-0526` Go 1.26.2 安全修复
- `AD-012` FinTech 系统安全架构
