# 反射工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [反射工具](#反射工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

反射工具提供了丰富的反射操作函数，简化类型检查、方法调用、字段访问等反射相关任务。

---

## 2. 功能特性

### 2.1 类型检查

- `GetType`: 获取值的类型名称
- `GetKind`: 获取值的Kind
- `IsNil`: 检查值是否为nil
- `IsZero`: 检查值是否为零值
- `IsPointer`: 检查值是否为指针
- `IsSlice`: 检查值是否为切片
- `IsMap`: 检查值是否为映射
- `IsStruct`: 检查值是否为结构体
- `IsInterface`: 检查值是否为接口
- `IsFunc`: 检查值是否为函数
- `IsChan`: 检查值是否为通道

### 2.2 指针操作

- `Dereference`: 解引用指针，如果不是指针则返回原值

### 2.3 结构体操作

- `GetField`: 获取结构体字段的值
- `SetField`: 设置结构体字段的值
- `HasField`: 检查结构体是否有指定字段
- `GetFieldNames`: 获取结构体的所有字段名
- `GetFieldTags`: 获取结构体字段的标签

### 2.4 方法操作

- `CallMethod`: 调用方法
- `HasMethod`: 检查值是否有指定方法
- `GetMethodNames`: 获取值的所有方法名

### 2.5 实例创建

- `NewInstance`: 创建类型的新实例
- `NewSlice`: 创建切片的新实例
- `NewMap`: 创建映射的新实例

### 2.6 类型转换

- `Convert`: 转换值的类型
- `IsAssignable`: 检查值是否可以赋值给目标类型
- `IsConvertible`: 检查值是否可以转换为目标类型

### 2.7 值比较和复制

- `DeepEqual`: 深度比较两个值是否相等
- `Copy`: 深度复制值

### 2.8 切片操作

- `GetSliceElement`: 获取切片元素
- `SetSliceElement`: 设置切片元素

### 2.9 映射操作

- `GetMapValue`: 获取映射的值
- `SetMapValue`: 设置映射的值

### 2.10 长度和容量

- `GetLength`: 获取值的长度（切片、映射、字符串、数组）
- `GetCapacity`: 获取值的容量（切片、数组）

---

## 3. 使用示例

### 3.1 类型检查

```go
import "github.com/yourusername/golang/pkg/utils/reflect"

// 获取类型名称
typeName := reflect.GetType(42) // "int"

// 获取Kind
kind := reflect.GetKind(42) // reflect.Int

// 检查是否为nil
if reflect.IsNil(ptr) {
    // 指针为nil
}

// 检查是否为零值
if reflect.IsZero(value) {
    // 值为零值
}

// 检查类型
if reflect.IsPointer(ptr) {
    // 是指针
}
if reflect.IsSlice(slice) {
    // 是切片
}
```

### 3.2 指针操作

```go
// 解引用指针
value := 42
ptr := &value
result := reflect.Dereference(ptr) // 42
```

### 3.3 结构体操作

```go
type User struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email" db:"user_email"`
}

user := User{Name: "test", Age: 30}

// 获取字段值
name, err := reflect.GetField(user, "Name")

// 设置字段值
err := reflect.SetField(&user, "Name", "new")

// 检查字段是否存在
if reflect.HasField(user, "Name") {
    // 字段存在
}

// 获取所有字段名
fields := reflect.GetFieldNames(user) // ["Name", "Age", "Email"]

// 获取字段标签
tags, err := reflect.GetFieldTags(user, "Email")
// tags: map[string]string{"json": "email", "db": "user_email"}
```

### 3.4 方法操作

```go
type Calculator struct {
    value int
}

func (c *Calculator) Add(n int) int {
    c.value += n
    return c.value
}

calc := &Calculator{value: 10}

// 调用方法
results, err := reflect.CallMethod(calc, "Add", 5)
// results: [15]

// 检查方法是否存在
if reflect.HasMethod(calc, "Add") {
    // 方法存在
}

// 获取所有方法名
methods := reflect.GetMethodNames(calc)
```

### 3.5 实例创建

```go
// 创建新实例
user := User{}
newUser := reflect.NewInstance(user).(*User)

// 创建切片
slice := reflect.NewSlice([]int{}, 0, 10).([]int)

// 创建映射
m := reflect.NewMap("", 0).(map[string]int)
```

### 3.6 类型转换

```go
// 转换类型
value := 42
int64Value, err := reflect.Convert(value, int64(0))

// 检查是否可以赋值
if reflect.IsAssignable(42, 0) {
    // 可以赋值
}

// 检查是否可以转换
if reflect.IsConvertible(42, int64(0)) {
    // 可以转换
}
```

### 3.7 值比较和复制

```go
// 深度比较
a := []int{1, 2, 3}
b := []int{1, 2, 3}
if reflect.DeepEqual(a, b) {
    // 相等
}

// 深度复制
original := []int{1, 2, 3}
copied := reflect.Copy(original)
```

### 3.8 切片和映射操作

```go
// 获取切片元素
slice := []int{1, 2, 3}
value, err := reflect.GetSliceElement(slice, 0) // 1

// 设置切片元素
err := reflect.SetSliceElement(&slice, 0, 10)

// 获取映射值
m := map[string]int{"a": 1}
value, ok := reflect.GetMapValue(m, "a") // 1, true

// 设置映射值
err := reflect.SetMapValue(&m, "b", 2)
```

### 3.9 长度和容量

```go
// 获取长度
slice := []int{1, 2, 3}
length, err := reflect.GetLength(slice) // 3

// 获取容量
capacity, err := reflect.GetCapacity(slice)
```

---

**更新日期**: 2025-11-11
