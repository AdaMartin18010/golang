# 反射/自解释能力

框架级别的反射和自解释能力，提供程序元数据、类型信息、函数信息等，使程序能够自我描述和解释。

## 📋 功能特性

- ✅ **类型检查**: 获取类型的完整信息（名称、包、方法、字段）
- ✅ **函数检查**: 获取函数的完整信息（名称、参数、返回值、位置）
- ✅ **结构体检查**: 获取结构体的完整信息（字段、标签）
- ✅ **自描述**: 提供对象的完整描述信息

## 🚀 快速开始

### 检查类型

```go
import "github.com/yourusername/golang/pkg/reflect"

inspector := reflect.NewInspector()

type User struct {
    ID    int    `json:"id" db:"id"`
    Name  string `json:"name" db:"name"`
    Email string `json:"email" db:"email"`
}

user := User{}
metadata := inspector.InspectType(user)

fmt.Printf("Type: %s\n", metadata.Name)
fmt.Printf("Package: %s\n", metadata.Package)
fmt.Printf("Fields: %d\n", len(metadata.Fields))
```

### 检查函数

```go
func Add(a, b int) int {
    return a + b
}

metadata := inspector.InspectFunction(Add)
fmt.Printf("Function: %s\n", metadata.Name)
fmt.Printf("Package: %s\n", metadata.Package)
fmt.Printf("File: %s:%d\n", metadata.File, metadata.Line)
fmt.Printf("Inputs: %v\n", metadata.Inputs)
fmt.Printf("Outputs: %v\n", metadata.Outputs)
```

### 检查结构体

```go
metadata := inspector.InspectStruct(user)
for _, field := range metadata.Fields {
    fmt.Printf("Field: %s, Type: %s, Tags: %v\n",
        field.Name, field.Type, field.Tags)
}
```

### 自描述

```go
description := inspector.Describe(user)
fmt.Println(description)
// 输出:
// Type: main.User
// Kind: struct
// Package: main
// Fields:
//   ID: int
//   Name: string
//   Email: string
```

## 📚 API 参考

### Inspector

```go
type Inspector struct{}

func NewInspector() *Inspector
func (i *Inspector) InspectType(v interface{}) TypeMetadata
func (i *Inspector) InspectFunction(fn interface{}) FunctionMetadata
func (i *Inspector) InspectStruct(v interface{}) StructMetadata
func (i *Inspector) Describe(v interface{}) string
```

## 🎯 使用场景

1. **API 文档生成**: 自动生成 API 文档
2. **数据验证**: 基于结构体标签进行验证
3. **序列化/反序列化**: 基于元数据进行序列化
4. **调试和诊断**: 运行时检查对象结构

## 🔗 相关文档

- [数据转换](../converter/README.md)
