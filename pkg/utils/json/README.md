# JSON工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [JSON工具](#json工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

JSON工具提供了丰富的JSON操作函数，简化常见的JSON处理任务。

---

## 2. 功能特性

### 2.1 基础操作

- `Marshal`: 序列化为JSON
- `Unmarshal`: 反序列化JSON
- `MarshalString`: 序列化为JSON字符串
- `UnmarshalString`: 从字符串反序列化
- `PrettyPrint`: 美化打印JSON

### 2.2 验证和检查

- `IsValidJSON`: 检查字符串是否为有效的JSON

### 2.3 路径操作

- `Get`: 从JSON对象中获取值（使用点号分隔的路径）
- `Set`: 设置JSON对象中的值（使用点号分隔的路径）

### 2.4 合并和转换

- `Merge`: 合并多个JSON对象
- `Transform`: 转换JSON结构
- `Filter`: 过滤JSON对象

### 2.5 扁平化

- `Flatten`: 扁平化嵌套JSON对象
- `Unflatten`: 反扁平化JSON对象

### 2.6 文件操作

- `ReadFile`: 从文件读取JSON
- `WriteFile`: 将JSON写入文件

### 2.7 流操作

- `Decode`: 从Reader解码JSON
- `Encode`: 编码JSON到Writer

---

## 3. 使用示例

### 3.1 基础操作

```go
import "github.com/yourusername/golang/pkg/utils/json"

// 序列化为JSON字符串
data := map[string]interface{}{
    "name": "test",
    "age":  30,
}
jsonStr, err := json.MarshalString(data)

// 从字符串反序列化
var result map[string]interface{}
err := json.UnmarshalString(jsonStr, &result)

// 美化打印
pretty, err := json.PrettyPrint(data)
```

### 3.2 路径操作

```go
data := []byte(`{"user":{"name":"test","age":30}}`)

// 获取值
name, err := json.Get(data, "user.name")

// 设置值
newData, err := json.Set(data, "user.name", "new")
```

### 3.3 合并和转换

```go
json1 := []byte(`{"a":1,"b":2}`)
json2 := []byte(`{"b":3,"c":4}`)

// 合并
merged, err := json.Merge(json1, json2)

// 过滤
filtered, err := json.Filter(data, func(k string, v interface{}) bool {
    return k != "b"
})
```

### 3.4 扁平化

```go
nested := []byte(`{"user":{"name":"test","age":30}}`)

// 扁平化
flattened, err := json.Flatten(nested, ".")

// 反扁平化
unflattened, err := json.Unflatten(flattened, ".")
```

### 3.5 文件操作

```go
// 从文件读取
var data map[string]interface{}
err := json.ReadFile("data.json", &data)

// 写入文件
err := json.WriteFile("output.json", data, true) // true表示格式化
```

### 3.6 验证

```go
valid := json.IsValidJSON(`{"name":"test"}`)
if !valid {
    // 处理无效JSON
}
```

---

**更新日期**: 2025-11-11
