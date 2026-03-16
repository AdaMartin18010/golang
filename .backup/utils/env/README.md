# 环境变量工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [环境变量工具](#环境变量工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

环境变量工具提供了丰富的环境变量操作函数，简化环境变量的读取、设置、验证和管理任务。

---

## 2. 功能特性

### 2.1 环境变量读取

- `Get`: 获取环境变量，如果不存在返回默认值
- `GetRequired`: 获取必需的环境变量，如果不存在则panic
- `GetInt`: 获取整数类型的环境变量
- `GetIntRequired`: 获取必需的整数类型的环境变量
- `GetInt64`: 获取64位整数类型的环境变量
- `GetInt64Required`: 获取必需的64位整数类型的环境变量
- `GetFloat64`: 获取浮点数类型的环境变量
- `GetFloat64Required`: 获取必需的浮点数类型的环境变量
- `GetBool`: 获取布尔类型的环境变量
- `GetBoolRequired`: 获取必需的布尔类型的环境变量
- `GetSlice`: 获取字符串切片类型的环境变量（使用逗号分隔）
- `GetSliceRequired`: 获取必需的字符串切片类型的环境变量
- `GetSliceWithSeparator`: 获取字符串切片类型的环境变量（使用指定分隔符）

### 2.2 环境变量操作

- `Set`: 设置环境变量
- `Unset`: 删除环境变量
- `Has`: 检查环境变量是否存在
- `IsSet`: 检查环境变量是否已设置（别名）
- `GetAll`: 获取所有环境变量
- `GetWithPrefix`: 获取所有以指定前缀开头的环境变量
- `Copy`: 复制环境变量到新map
- `Clear`: 清除所有环境变量

### 2.3 环境变量展开

- `Expand`: 展开环境变量（支持 ${VAR} 或 $VAR 格式）
- `ExpandMap`: 展开map中的环境变量

### 2.4 环境变量加载

- `LoadFromFile`: 从文件加载环境变量（.env格式）
- `MustLoadFromFile`: 从文件加载环境变量，如果失败则panic

### 2.5 环境变量验证

- `ValidateRequired`: 验证必需的环境变量是否都已设置

### 2.6 环境变量工具

- `Filter`: 过滤环境变量
- `Merge`: 合并环境变量

---

## 3. 使用示例

### 3.1 环境变量读取

```go
import "github.com/yourusername/golang/pkg/utils/env"

// 获取环境变量（带默认值）
value := env.Get("DATABASE_URL", "localhost:5432")

// 获取必需的环境变量
value := env.GetRequired("API_KEY")

// 获取整数类型的环境变量
port := env.GetInt("PORT", 8080)

// 获取64位整数类型的环境变量
timeout := env.GetInt64("TIMEOUT", 30)

// 获取浮点数类型的环境变量
rate := env.GetFloat64("RATE", 0.5)

// 获取布尔类型的环境变量
debug := env.GetBool("DEBUG", false)

// 获取字符串切片类型的环境变量
hosts := env.GetSlice("HOSTS", []string{"localhost"})
```

### 3.2 环境变量操作

```go
// 设置环境变量
err := env.Set("KEY", "value")

// 删除环境变量
err := env.Unset("KEY")

// 检查环境变量是否存在
if env.Has("KEY") {
    // 环境变量存在
}

// 获取所有环境变量
allEnv := env.GetAll()

// 获取所有以指定前缀开头的环境变量
dbEnv := env.GetWithPrefix("DB_")
```

### 3.3 环境变量展开

```go
// 展开环境变量
os.Setenv("NAME", "world")
result := env.Expand("Hello ${NAME}") // "Hello world"

// 展开map中的环境变量
m := map[string]string{
    "greeting": "Hello ${NAME}",
}
expanded := env.ExpandMap(m)
```

### 3.4 环境变量加载

```go
// 从文件加载环境变量
err := env.LoadFromFile(".env")

// 从文件加载环境变量（失败则panic）
env.MustLoadFromFile(".env")
```

### 3.5 环境变量验证

```go
// 验证必需的环境变量是否都已设置
err := env.ValidateRequired([]string{
    "DATABASE_URL",
    "API_KEY",
    "PORT",
})
if err != nil {
    // 处理错误
}
```

### 3.6 环境变量工具

```go
// 过滤环境变量
filtered := env.Filter(func(key, value string) bool {
    return strings.HasPrefix(key, "APP_")
})

// 合并环境变量
merged := env.Merge(env1, env2, env3)
```

### 3.7 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/env"
)

func main() {
    // 加载.env文件
    if err := env.LoadFromFile(".env"); err != nil {
        fmt.Printf("Warning: failed to load .env file: %v\n", err)
    }

    // 验证必需的环境变量
    if err := env.ValidateRequired([]string{
        "DATABASE_URL",
        "API_KEY",
    }); err != nil {
        panic(err)
    }

    // 读取配置
    dbURL := env.GetRequired("DATABASE_URL")
    apiKey := env.GetRequired("API_KEY")
    port := env.GetInt("PORT", 8080)
    debug := env.GetBool("DEBUG", false)

    fmt.Printf("Database URL: %s\n", dbURL)
    fmt.Printf("API Key: %s\n", apiKey)
    fmt.Printf("Port: %d\n", port)
    fmt.Printf("Debug: %v\n", debug)
}
```

---

**更新日期**: 2025-11-11
