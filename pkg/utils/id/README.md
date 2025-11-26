# ID生成工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [ID生成工具](#id生成工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)
  - [4. 最佳实践](#4-最佳实践)

---

## 1. 概述

ID生成工具提供了多种ID生成策略，满足不同场景的需求：

- ✅ **UUID**: 标准UUID生成
- ✅ **短UUID**: 22字符的短UUID
- ✅ **NanoID**: 高性能的短ID生成
- ✅ **雪花ID**: 分布式唯一ID生成
- ✅ **时间戳ID**: 基于时间戳的ID
- ✅ **随机ID**: 随机十六进制/Base64 ID
- ✅ **顺序ID**: 顺序递增ID

---

## 2. 功能特性

### 2.1 UUID生成器

标准UUID（36字符，带连字符）

```go
id := id.UUID() // "550e8400-e29b-41d4-a716-446655440000"
```

### 2.2 短UUID生成器

短UUID（22字符，Base64编码）

```go
id := id.ShortUUID() // "dQw4w9WgXcQ"
```

### 2.3 NanoID生成器

高性能短ID（默认21字符，可自定义长度）

```go
id := id.NanoID() // "V1StGXR8_Z5jdHi6B-myT"
id := id.NanoIDWithSize(10) // "V1StGXR8_Z"
```

### 2.4 雪花ID生成器

分布式唯一ID（64位整数）

```go
generator := id.NewSnowflakeGenerator(1, 1)
id := generator.Generate() // "1234567890123456789"
```

### 2.5 时间戳ID生成器

基于时间戳的ID

```go
id := id.TimestampID() // "1699123456789012345"
```

### 2.6 随机ID生成器

随机十六进制或Base64 ID

```go
hexID := id.RandomHex() // "a1b2c3d4e5f6..."
hexID := id.RandomHexWithLength(16) // "a1b2c3d4e5f6g7h8"

base64ID := id.RandomBase64() // "dGVzdA=="
base64ID := id.RandomBase64WithLength(16) // "dGVzdA=="
```

### 2.7 顺序ID生成器

顺序递增ID（带前缀）

```go
id := id.SequentialID("ORDER") // "ORDER000000000001"
```

---

## 3. 使用示例

### 3.1 基本使用

```go
import "github.com/yourusername/golang/pkg/utils/id"

// UUID
uuid := id.UUID()

// 短UUID
shortUUID := id.ShortUUID()

// NanoID
nanoID := id.NanoID()

// 时间戳ID
timestampID := id.TimestampID()
```

### 3.2 使用生成器

```go
// UUID生成器
uuidGen := id.NewUUIDGenerator()
uuid := uuidGen.Generate()

// 短UUID生成器
shortUUIDGen := id.NewShortUUIDGenerator()
shortUUID := shortUUIDGen.Generate()

// NanoID生成器
nanoIDGen := id.NewNanoIDGenerator(21)
nanoID := nanoIDGen.Generate()

// 雪花ID生成器
snowflakeGen := id.NewSnowflakeGenerator(1, 1)
snowflakeID := snowflakeGen.Generate()

// 顺序ID生成器
sequentialGen := id.NewSequentialIDGenerator("ORDER")
orderID := sequentialGen.Generate()
```

### 3.3 自定义配置

```go
// 自定义长度的NanoID
nanoID := id.NanoIDWithSize(10)

// 自定义长度的随机十六进制ID
hexID := id.RandomHexWithLength(16)

// 自定义长度的随机Base64 ID
base64ID := id.RandomBase64WithLength(16)
```

---

## 4. 最佳实践

### 4.1 选择合适策略

- **UUID**: 适用于需要全局唯一性的场景
- **短UUID**: 适用于需要短ID但保持唯一性的场景
- **NanoID**: 适用于需要高性能短ID的场景
- **雪花ID**: 适用于分布式系统，需要有序ID的场景
- **时间戳ID**: 适用于需要时间信息的场景
- **随机ID**: 适用于需要随机性的场景
- **顺序ID**: 适用于需要顺序递增的场景

### 4.2 性能考虑

- **UUID**: 性能中等，适合大多数场景
- **短UUID**: 性能较好，适合需要短ID的场景
- **NanoID**: 性能最好，适合高并发场景
- **雪花ID**: 性能好，但需要配置worker ID和datacenter ID
- **时间戳ID**: 性能最好，但可能重复（需要额外处理）

### 4.3 分布式场景

在分布式系统中，推荐使用：
- **雪花ID**: 保证全局唯一且有序
- **UUID**: 保证全局唯一
- **NanoID**: 高性能且唯一性高

---

**更新日期**: 2025-11-11
