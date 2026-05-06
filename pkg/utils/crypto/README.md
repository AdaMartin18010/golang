# 加密工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26.2

---

## 📋 目录

- [加密工具](#加密工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 密码哈希](#21-密码哈希)
    - [2.2 AES加密](#22-aes加密)
    - [2.3 密钥生成](#23-密钥生成)
    - [2.4 哈希函数](#24-哈希函数)
    - [2.5 HMAC](#25-hmac)
    - [2.6 密钥派生](#26-密钥派生)
  - [3. 使用示例](#3-使用示例)
    - [3.1 密码哈希](#31-密码哈希)
    - [3.2 AES加密](#32-aes加密)
    - [3.3 使用加密器](#33-使用加密器)
    - [3.4 密钥派生](#34-密钥派生)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 密码存储](#41-密码存储)
    - [4.2 密钥管理](#42-密钥管理)
    - [4.3 加密使用](#43-加密使用)
    - [4.4 密钥派生](#44-密钥派生)

---

## 1. 概述

加密工具提供了常用的加密和哈希功能：

- ✅ **密码哈希**: bcrypt密码哈希和验证
- ✅ **AES加密**: AES-256-GCM对称加密
- ✅ **密钥生成**: 随机密钥和盐值生成
- ✅ **哈希函数**: SHA256哈希
- ✅ **HMAC**: HMAC计算
- ✅ **密钥派生**: PBKDF2密钥派生

---

## 2. 功能特性

### 2.1 密码哈希

使用bcrypt进行密码哈希，安全可靠。

```go
hash, err := crypto.HashPassword("password123")
isValid := crypto.CheckPassword("password123", hash)
```

### 2.2 AES加密

AES-256-GCM模式加密，提供认证和加密。

```go
key, _ := crypto.GenerateAESKey()
ciphertext, _ := crypto.AESEncrypt("plaintext", key)
plaintext, _ := crypto.AESDecrypt(ciphertext, key)
```

### 2.3 密钥生成

生成随机密钥和盐值。

```go
key, _ := crypto.GenerateAESKey()
salt, _ := crypto.GenerateSalt(16)
```

### 2.4 哈希函数

SHA256哈希计算。

```go
hash := crypto.SHA256Hash("data")
```

### 2.5 HMAC

HMAC计算（使用SHA256）。

```go
hmac := crypto.HMAC("key", "message")
```

### 2.6 密钥派生

从密码派生密钥（PBKDF2）。

```go
key := crypto.DeriveKey("password", "salt", 10000)
```

---

## 3. 使用示例

### 3.1 密码哈希

```go
import "github.com/yourusername/golang/pkg/utils/crypto"

// 哈希密码
hash, err := crypto.HashPassword("userpassword")
if err != nil {
    // 处理错误
}

// 验证密码
isValid := crypto.CheckPassword("userpassword", hash)
if !isValid {
    // 密码错误
}
```

### 3.2 AES加密

```go
// 生成密钥
key, err := crypto.GenerateAESKey()
if err != nil {
    // 处理错误
}

// 加密
ciphertext, err := crypto.AESEncrypt("sensitive data", key)
if err != nil {
    // 处理错误
}

// 解密
plaintext, err := crypto.AESDecrypt(ciphertext, key)
if err != nil {
    // 处理错误
}
```

### 3.3 使用加密器

```go
// 创建加密器
key, _ := crypto.GenerateAESKey()
encryptor, err := crypto.NewAESEncryptor(key)
if err != nil {
    // 处理错误
}

// 加密
ciphertext, err := encryptor.Encrypt("data")
if err != nil {
    // 处理错误
}

// 解密
plaintext, err := encryptor.Decrypt(ciphertext)
if err != nil {
    // 处理错误
}
```

### 3.4 密钥派生

```go
// 生成盐值
salt, err := crypto.GenerateSalt(16)
if err != nil {
    // 处理错误
}

// 派生密钥
key := crypto.DeriveKey("password", salt, 10000)
```

---

## 4. 最佳实践

### 4.1 密码存储

- 使用bcrypt哈希密码，不要使用MD5或SHA256
- 设置合适的cost值（默认值通常足够）
- 永远不要存储明文密码

### 4.2 密钥管理

- 使用安全的随机数生成器生成密钥
- 妥善保管密钥，不要硬编码在代码中
- 使用密钥管理系统（如HashiCorp Vault、AWS KMS）

### 4.3 加密使用

- 使用AES-256-GCM模式（提供认证和加密）
- 每次加密使用不同的nonce
- 不要重复使用相同的密钥和nonce组合

### 4.4 密钥派生

- 使用PBKDF2从密码派生密钥
- 设置足够的迭代次数（建议10000+）
- 使用随机盐值

---

**更新日期**: 2025-11-11
