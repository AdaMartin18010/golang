# 压缩工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [压缩工具](#压缩工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

压缩工具提供了丰富的压缩解压函数，支持gzip、zlib等常见压缩格式，简化数据压缩和解压任务。

---

## 2. 功能特性

### 2.1 Gzip压缩

- `GzipCompress`: gzip压缩
- `GzipDecompress`: gzip解压
- `GzipCompressToFile`: gzip压缩到文件
- `GzipDecompressFromFile`: 从文件gzip解压
- `CompressLevel`: gzip压缩（指定压缩级别）
- `CompressBest`: gzip压缩（最佳压缩率）
- `CompressFast`: gzip压缩（最快速度）
- `CompressDefault`: gzip压缩（默认压缩率）
- `CompressNoCompression`: gzip压缩（不压缩）
- `IsGzip`: 检查数据是否为gzip格式

### 2.2 Zlib压缩

- `ZlibCompress`: zlib压缩
- `ZlibDecompress`: zlib解压

### 2.3 流式压缩

- `CompressStream`: gzip压缩流
- `DecompressStream`: gzip解压流

### 2.4 压缩统计

- `GetCompressionRatio`: 获取压缩率
- `GetCompressionSavings`: 获取压缩节省的字节数

---

## 3. 使用示例

### 3.1 Gzip压缩

```go
import "github.com/yourusername/golang/pkg/utils/compress"

// gzip压缩
data := []byte("hello world")
compressed, err := compress.GzipCompress(data)
if err != nil {
    // 处理错误
}

// gzip解压
decompressed, err := compress.GzipDecompress(compressed)
if err != nil {
    // 处理错误
}

// 压缩到文件
err := compress.GzipCompressToFile(data, "data.gz")

// 从文件解压
decompressed, err := compress.GzipDecompressFromFile("data.gz")
```

### 3.2 压缩级别

```go
// 最佳压缩率
compressed, err := compress.CompressBest(data)

// 最快速度
compressed, err := compress.CompressFast(data)

// 默认压缩率
compressed, err := compress.CompressDefault(data)

// 指定压缩级别（0-9）
compressed, err := compress.CompressLevel(data, 6)
```

### 3.3 Zlib压缩

```go
// zlib压缩
compressed, err := compress.ZlibCompress(data)

// zlib解压
decompressed, err := compress.ZlibDecompress(compressed)
```

### 3.4 流式压缩

```go
// 压缩流
reader := bytes.NewReader(data)
var buf bytes.Buffer
err := compress.CompressStream(reader, &buf)

// 解压流
compressedReader := bytes.NewReader(compressed)
var decompressedBuf bytes.Buffer
err := compress.DecompressStream(compressedReader, &decompressedBuf)
```

### 3.5 压缩统计

```go
// 检查是否为gzip格式
if compress.IsGzip(data) {
    // 是gzip格式
}

// 获取压缩率
ratio := compress.GetCompressionRatio(originalSize, compressedSize)

// 获取压缩节省的字节数
savings := compress.GetCompressionSavings(originalSize, compressedSize)
```

### 3.6 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/compress"
)

func main() {
    // 原始数据
    data := []byte("This is a test string that will be compressed")
    originalSize := len(data)

    // 压缩
    compressed, err := compress.GzipCompress(data)
    if err != nil {
        panic(err)
    }
    compressedSize := len(compressed)

    // 解压
    decompressed, err := compress.GzipDecompress(compressed)
    if err != nil {
        panic(err)
    }

    // 统计信息
    ratio := compress.GetCompressionRatio(originalSize, compressedSize)
    savings := compress.GetCompressionSavings(originalSize, compressedSize)

    fmt.Printf("Original size: %d bytes\n", originalSize)
    fmt.Printf("Compressed size: %d bytes\n", compressedSize)
    fmt.Printf("Compression ratio: %.2f%%\n", ratio)
    fmt.Printf("Space saved: %d bytes\n", savings)
    fmt.Printf("Decompressed matches original: %v\n",
        string(data) == string(decompressed))
}
```

---

**更新日期**: 2025-11-11
