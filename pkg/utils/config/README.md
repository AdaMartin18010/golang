# 配置管理工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [配置管理工具](#配置管理工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 文件配置加载器](#21-文件配置加载器)
    - [2.2 环境变量加载器](#22-环境变量加载器)
    - [2.3 Map配置加载器](#23-map配置加载器)
    - [2.4 多源配置加载器](#24-多源配置加载器)
  - [3. 使用示例](#3-使用示例)
    - [3.1 从文件加载配置](#31-从文件加载配置)
    - [3.2 从环境变量加载配置](#32-从环境变量加载配置)
    - [3.3 从Map加载配置](#33-从map加载配置)
    - [3.4 多源配置加载](#34-多源配置加载)
    - [3.5 配置结构体标签](#35-配置结构体标签)

---

## 1. 概述

配置管理工具提供了灵活的配置加载机制，支持从多种源加载配置：

- ✅ **文件配置**: 从JSON文件加载配置
- ✅ **环境变量**: 从环境变量加载配置
- ✅ **Map配置**: 从Map加载配置
- ✅ **多源配置**: 支持多个配置源合并

---

## 2. 功能特性

### 2.1 文件配置加载器

从JSON文件加载配置。

```go
loader := config.NewFileLoader("config.json")
err := loader.Load(&config)
```

### 2.2 环境变量加载器

从环境变量加载配置，支持前缀。

```go
loader := config.NewEnvLoader("APP")
err := loader.Load(&config)
```

### 2.3 Map配置加载器

从Map加载配置。

```go
data := map[string]interface{}{
    "name": "test",
    "port": 8080,
}
loader := config.NewMapLoader(data)
err := loader.Load(&config)
```

### 2.4 多源配置加载器

支持从多个源加载配置，后面的会覆盖前面的。

```go
fileLoader := config.NewFileLoader("config.json")
envLoader := config.NewEnvLoader("APP")
multiLoader := config.NewMultiLoader(fileLoader, envLoader)
err := multiLoader.Load(&config)
```

---

## 3. 使用示例

### 3.1 从文件加载配置

```go
import "github.com/yourusername/golang/pkg/utils/config"

type Config struct {
    Name string `json:"name"`
    Port int    `json:"port"`
    Host string `json:"host"`
}

var cfg Config
err := config.Load("config.json", &cfg)
if err != nil {
    // 处理错误
}
```

### 3.2 从环境变量加载配置

```go
type Config struct {
    Name string `env:"name"`
    Port int    `env:"port"`
    Host string `env:"host"`
}

var cfg Config
err := config.LoadFromEnv("APP", &cfg)
// 会读取 APP_NAME, APP_PORT, APP_HOST 环境变量
```

### 3.3 从Map加载配置

```go
data := map[string]interface{}{
    "name": "test",
    "port": 8080,
    "host": "localhost",
}

var cfg Config
err := config.LoadFromMap(data, &cfg)
```

### 3.4 多源配置加载

```go
// 先加载文件配置
fileLoader := config.NewFileLoader("config.json")

// 再加载环境变量配置（会覆盖文件配置）
envLoader := config.NewEnvLoader("APP")

// 合并配置
multiLoader := config.NewMultiLoader(fileLoader, envLoader)
err := multiLoader.Load(&cfg)
```

### 3.5 配置结构体标签

```go
type Config struct {
    // JSON标签用于文件配置
    // env标签用于环境变量配置
    // map标签用于Map配置
    Name string `json:"name" env:"name" map:"name"`
    Port int    `json:"port" env:"port" map:"port"`
    Host string `json:"host" env:"host" map:"host"`
}
```

---

**更新日期**: 2025-11-11
