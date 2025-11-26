# 编码工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [编码工具](#编码工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

编码工具提供了丰富的编码解码函数，支持Base64、十六进制、JSON等常见编码格式，以及类型转换功能。

---

## 2. 功能特性

### 2.1 Base64编码

- `Base64Encode`: Base64编码
- `Base64Decode`: Base64解码
- `Base64URLEncode`: Base64 URL安全编码
- `Base64URLDecode`: Base64 URL安全解码
- `Base64RawStdEncode`: Base64原始标准编码（无填充）
- `Base64RawStdDecode`: Base64原始标准解码（无填充）
- `Base64RawURLEncode`: Base64原始URL安全编码（无填充）
- `Base64RawURLDecode`: Base64原始URL安全解码（无填充）
- `Base64EncodeString`: Base64编码字符串
- `Base64DecodeString`: Base64解码字符串
- `IsBase64`: 检查字符串是否为有效的Base64编码

### 2.2 十六进制编码

- `HexEncode`: 十六进制编码
- `HexDecode`: 十六进制解码
- `HexEncodeUpper`: 十六进制编码（大写）
- `HexDecodeUpper`: 十六进制解码（大写）
- `HexEncodeString`: 十六进制编码字符串
- `HexDecodeString`: 十六进制解码字符串
- `IsHex`: 检查字符串是否为有效的十六进制编码

### 2.3 类型转换

- `StringToBytes`: 字符串转字节数组
- `BytesToString`: 字节数组转字符串
- `IntToString`: 整数转字符串
- `Int64ToString`: 64位整数转字符串
- `StringToInt`: 字符串转整数
- `StringToInt64`: 字符串转64位整数
- `Float64ToString`: 浮点数转字符串
- `StringToFloat64`: 字符串转浮点数
- `BoolToString`: 布尔值转字符串
- `StringToBool`: 字符串转布尔值
- `RuneToBytes`: 字符转字节数组
- `BytesToRunes`: 字节数组转字符数组
- `RunesToString`: 字符数组转字符串
- `StringToRunes`: 字符串转字符数组

### 2.4 JSON编码

- `JSONEncode`: JSON编码
- `JSONEncodePretty`: JSON编码（格式化）
- `JSONDecode`: JSON解码
- `JSONEncodeString`: JSON编码为字符串
- `JSONEncodePrettyString`: JSON编码为字符串（格式化）
- `JSONDecodeString`: JSON解码字符串
- `IsJSON`: 检查字符串是否为有效的JSON

### 2.5 字符串转义

- `EscapeString`: 转义字符串（HTML实体）
- `UnescapeString`: 反转义字符串（HTML实体）
- `EscapeURL`: 转义URL
- `UnescapeURL`: 反转义URL

---

## 3. 使用示例

### 3.1 Base64编码

```go
import "github.com/yourusername/golang/pkg/utils/encoding"

// Base64编码
data := []byte("hello world")
encoded := encoding.Base64Encode(data)
decoded, err := encoding.Base64Decode(encoded)

// Base64 URL安全编码
urlEncoded := encoding.Base64URLEncode(data)
urlDecoded, err := encoding.Base64URLDecode(urlEncoded)

// Base64编码字符串
strEncoded := encoding.Base64EncodeString("hello")
strDecoded, err := encoding.Base64DecodeString(strEncoded)

// 检查是否为有效的Base64
if encoding.IsBase64(encoded) {
    // 有效的Base64编码
}
```

### 3.2 十六进制编码

```go
// 十六进制编码
data := []byte("hello world")
encoded := encoding.HexEncode(data)
decoded, err := encoding.HexDecode(encoded)

// 十六进制编码（大写）
upperEncoded := encoding.HexEncodeUpper(data)

// 十六进制编码字符串
strEncoded := encoding.HexEncodeString("hello")
strDecoded, err := encoding.HexDecodeString(strEncoded)

// 检查是否为有效的十六进制
if encoding.IsHex(encoded) {
    // 有效的十六进制编码
}
```

### 3.3 类型转换

```go
// 字符串和字节数组转换
bytes := encoding.StringToBytes("hello")
str := encoding.BytesToString(bytes)

// 整数转换
str := encoding.IntToString(123)
num, err := encoding.StringToInt("123")

// 浮点数转换
str := encoding.Float64ToString(123.456)
num, err := encoding.StringToFloat64("123.456")

// 布尔值转换
str := encoding.BoolToString(true)
b, err := encoding.StringToBool("true")

// 字符数组转换
runes := encoding.StringToRunes("hello")
str := encoding.RunesToString(runes)
```

### 3.4 JSON编码

```go
// JSON编码
data := map[string]interface{}{
    "name": "test",
    "age":  30,
}
encoded, err := encoding.JSONEncode(data)

// JSON编码（格式化）
pretty, err := encoding.JSONEncodePretty(data)

// JSON解码
var decoded map[string]interface{}
err = encoding.JSONDecode(encoded, &decoded)

// JSON编码为字符串
str, err := encoding.JSONEncodeString(data)

// 检查是否为有效的JSON
if encoding.IsJSON(`{"name":"test"}`) {
    // 有效的JSON
}
```

### 3.5 字符串转义

```go
// HTML实体转义
escaped := encoding.EscapeString("<script>alert('xss')</script>")
unescaped := encoding.UnescapeString(escaped)

// URL转义
urlEscaped := encoding.EscapeURL("hello world")
urlUnescaped, err := encoding.UnescapeURL(urlEscaped)
```

---

**更新日期**: 2025-11-11
