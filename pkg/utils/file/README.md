# 文件操作工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [文件操作工具](#文件操作工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

文件操作工具提供了丰富的文件系统操作函数，简化常见的文件处理任务。

---

## 2. 功能特性

### 2.1 文件检查

- `Exists`: 检查文件或目录是否存在
- `IsFile`: 检查路径是否为文件
- `IsDir`: 检查路径是否为目录

### 2.2 文件读写

- `ReadFile`: 读取文件内容
- `ReadFileString`: 读取文件内容为字符串
- `WriteFile`: 写入文件内容
- `WriteFileString`: 写入字符串到文件
- `AppendFile`: 追加内容到文件
- `AppendFileString`: 追加字符串到文件
- `ReadLines`: 读取文件的所有行
- `WriteLines`: 写入多行到文件

### 2.3 文件操作

- `CopyFile`: 复制文件
- `MoveFile`: 移动文件
- `DeleteFile`: 删除文件
- `DeleteDir`: 删除目录（递归）

### 2.4 目录操作

- `CreateDir`: 创建目录
- `ListFiles`: 列出目录中的文件
- `ListDirs`: 列出目录中的子目录
- `ListAll`: 列出目录中的所有条目
- `WalkFiles`: 遍历目录中的所有文件
- `WalkDirs`: 遍历目录中的所有子目录
- `EnsureDir`: 确保目录存在
- `EnsureFileDir`: 确保文件所在目录存在

### 2.5 路径操作

- `GetExt`: 获取文件扩展名
- `GetBaseName`: 获取文件名（不含路径）
- `GetDirName`: 获取目录名
- `JoinPath`: 连接路径
- `CleanPath`: 清理路径
- `AbsPath`: 获取绝对路径
- `RelPath`: 获取相对路径
- `MatchPattern`: 匹配文件模式
- `Glob`: 匹配文件模式（支持通配符）

### 2.6 文件信息

- `GetFileSize`: 获取文件大小
- `GetFileMode`: 获取文件权限
- `Chmod`: 修改文件权限
- `Chown`: 修改文件所有者

---

## 3. 使用示例

### 3.1 文件检查

```go
import "github.com/yourusername/golang/pkg/utils/file"

// 检查文件是否存在
if file.Exists("test.txt") {
    // 文件存在
}

// 检查是否为文件
if file.IsFile("test.txt") {
    // 是文件
}

// 检查是否为目录
if file.IsDir("/path/to/dir") {
    // 是目录
}
```

### 3.2 文件读写

```go
// 读取文件
content, err := file.ReadFileString("test.txt")

// 写入文件
err := file.WriteFileString("test.txt", "content", 0644)

// 追加内容
err := file.AppendFileString("test.txt", "more content", 0644)

// 读取所有行
lines, err := file.ReadLines("test.txt")

// 写入多行
lines := []string{"line1", "line2", "line3"}
err := file.WriteLines("test.txt", lines, 0644)
```

### 3.3 文件操作

```go
// 复制文件
err := file.CopyFile("source.txt", "dest.txt")

// 移动文件
err := file.MoveFile("old.txt", "new.txt")

// 删除文件
err := file.DeleteFile("test.txt")

// 删除目录
err := file.DeleteDir("/path/to/dir")
```

### 3.4 目录操作

```go
// 创建目录
err := file.CreateDir("/path/to/dir", 0755)

// 列出文件
files, err := file.ListFiles("/path/to/dir")

// 列出子目录
dirs, err := file.ListDirs("/path/to/dir")

// 遍历所有文件
err := file.WalkFiles("/path/to/dir", func(path string) error {
    // 处理文件
    return nil
})

// 确保目录存在
err := file.EnsureDir("/path/to/dir", 0755)
```

### 3.5 路径操作

```go
// 获取扩展名
ext := file.GetExt("test.txt") // ".txt"

// 获取文件名
name := file.GetBaseName("/path/to/file.txt") // "file.txt"

// 获取目录名
dir := file.GetDirName("/path/to/file.txt") // "/path/to"

// 连接路径
path := file.JoinPath("/path", "to", "file.txt")

// 获取绝对路径
absPath, err := file.AbsPath("relative/path")

// 匹配文件模式
matches, err := file.Glob("*.txt")
```

---

**更新日期**: 2025-11-11
