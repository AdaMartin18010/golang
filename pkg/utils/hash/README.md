# 哈希工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [哈希工具](#哈希工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

哈希工具提供了丰富的哈希算法函数，支持MD5、SHA1、SHA256、SHA512、CRC32、CRC64、FNV等常见哈希算法，简化数据哈希和校验任务。

---

## 2. 功能特性

### 2.1 MD5哈希

- `MD5`: MD5哈希
- `MD5String`: MD5哈希字符串
- `MD5File`: MD5哈希文件

### 2.2 SHA系列哈希

- `SHA1`: SHA1哈希
- `SHA1String`: SHA1哈希字符串
- `SHA1File`: SHA1哈希文件
- `SHA256`: SHA256哈希
- `SHA256String`: SHA256哈希字符串
- `SHA256File`: SHA256哈希文件
- `SHA512`: SHA512哈希
- `SHA512String`: SHA512哈希字符串
- `SHA512File`: SHA512哈希文件

### 2.3 CRC校验和

- `CRC32`: CRC32校验和
- `CRC32String`: CRC32校验和字符串
- `CRC32File`: CRC32校验和文件
- `CRC64`: CRC64校验和
- `CRC64String`: CRC64校验和字符串
- `CRC64File`: CRC64校验和文件

### 2.4 FNV哈希

- `FNV32`: FNV32哈希
- `FNV32String`: FNV32哈希字符串
- `FNV32a`: FNV32a哈希
- `FNV32aString`: FNV32a哈希字符串
- `FNV64`: FNV64哈希
- `FNV64String`: FNV64哈希字符串
- `FNV64a`: FNV64a哈希
- `FNV64aString`: FNV64a哈希字符串
- `FNV128`: FNV128哈希
- `FNV128String`: FNV128哈希字符串
- `FNV128a`: FNV128a哈希
- `FNV128aString`: FNV128a哈希字符串

### 2.5 通用哈希函数

- `Hash`: 通用哈希函数
- `HashString`: 通用哈希函数（字符串）
- `HashFile`: 通用哈希函数（文件）

### 2.6 哈希验证

- `CompareHash`: 比较哈希值
- `VerifyHash`: 验证哈希值
- `VerifyHashString`: 验证哈希值（字符串）
- `VerifyHashFile`: 验证哈希值（文件）

---

## 3. 使用示例

### 3.1 MD5哈希

```go
import "github.com/yourusername/golang/pkg/utils/hash"

// MD5哈希
data := []byte("hello world")
hash := hash.MD5(data)

// MD5哈希字符串
hash := hash.MD5String("hello world")

// MD5哈希文件
hash, err := hash.MD5File("file.txt")
```

### 3.2 SHA系列哈希

```go
// SHA1哈希
hash := hash.SHA1(data)
hash := hash.SHA1String("hello world")
hash, err := hash.SHA1File("file.txt")

// SHA256哈希
hash := hash.SHA256(data)
hash := hash.SHA256String("hello world")
hash, err := hash.SHA256File("file.txt")

// SHA512哈希
hash := hash.SHA512(data)
hash := hash.SHA512String("hello world")
hash, err := hash.SHA512File("file.txt")
```

### 3.3 CRC校验和

```go
// CRC32校验和
checksum := hash.CRC32(data)
checksum := hash.CRC32String("hello world")
checksum, err := hash.CRC32File("file.txt")

// CRC64校验和
checksum := hash.CRC64(data)
checksum := hash.CRC64String("hello world")
checksum, err := hash.CRC64File("file.txt")
```

### 3.4 FNV哈希

```go
// FNV32哈希
hash := hash.FNV32(data)
hash := hash.FNV32String("hello world")

// FNV64哈希
hash := hash.FNV64(data)
hash := hash.FNV64String("hello world")

// FNV128哈希
hash := hash.FNV128(data)
hash := hash.FNV128String("hello world")
```

### 3.5 通用哈希函数

```go
// 通用哈希函数
hash, err := hash.Hash(data, "md5")
hash, err := hash.Hash(data, "sha256")
hash, err := hash.Hash(data, "sha512")

// 通用哈希函数（字符串）
hash, err := hash.HashString("hello world", "md5")

// 通用哈希函数（文件）
hash, err := hash.HashFile("file.txt", "sha256")
```

### 3.6 哈希验证

```go
// 比较哈希值
if hash.CompareHash(hash1, hash2) {
    // 哈希值相同
}

// 验证哈希值
valid, err := hash.VerifyHash(data, "md5", expectedHash)
if valid {
    // 哈希值验证通过
}

// 验证哈希值（字符串）
valid, err := hash.VerifyHashString("hello world", "md5", expectedHash)

// 验证哈希值（文件）
valid, err := hash.VerifyHashFile("file.txt", "sha256", expectedHash)
```

### 3.7 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/hash"
)

func main() {
    data := []byte("hello world")

    // 计算各种哈希值
    md5Hash := hash.MD5(data)
    sha256Hash := hash.SHA256(data)
    crc32Checksum := hash.CRC32(data)

    fmt.Printf("MD5: %s\n", md5Hash)
    fmt.Printf("SHA256: %s\n", sha256Hash)
    fmt.Printf("CRC32: %d\n", crc32Checksum)

    // 验证哈希值
    valid, err := hash.VerifyHash(data, "md5", md5Hash)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hash verification: %v\n", valid)
}
```

---

**更新日期**: 2025-11-11
