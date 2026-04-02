# 包管理详解 (Package Management)

> **分类**: 语言设计
> **标签**: #package #module #import

---

## 包声明

```go
// 包名与目录名无关，但建议一致
package user

// main 包是可执行程序
package main

// 内部包（Go 1.4+）
// 放在 internal/ 目录下的包只能被父目录导入
```

---

## 导入模式

### 标准导入

```go
import (
    "fmt"
    "os"
)
```

### 别名导入

```go
import (
    f "fmt"           // 短别名
    myfmt "fmt"       // 自定义别名
    _ "github.com/lib/pq"  // 只执行 init()
    . "fmt"           // 点导入（不推荐）
)
```

### 条件导入

```go
//go:build linux
import "syscall"

//go:build windows
import "golang.org/x/sys/windows"
```

---

## init 函数

```go
package mypkg

var globalVar int

// 按文件名字母顺序执行
func init() {
    globalVar = 42
}

func init() {
    // 可以有多个 init
}
```

### init 执行顺序

```
1. 导入包的 init
2. 本包变量初始化
3. 本包 init 函数
4. main 函数
```

---

## 循环导入

### 问题

```go
// a.go
package a
import "b"

// b.go
package b
import "a"  // 循环导入！编译错误
```

### 解决

```go
// 提取公共接口到第三个包
// common/types.go
package common

type User interface {
    GetName() string
}

// a/a.go
package a
import "common"

// b/b.go
package b
import "common"
```

---

## 可见性规则

```go
package mypkg

// 大写开头 = 公开
func PublicFunc() {}
var PublicVar int

// 小写开头 = 私有
func privateFunc() {}
var privateVar int

// 结构体字段
 type MyStruct struct {
    PublicField  int  // 公开
    privateField int  // 私有
}
```

---

## 包设计最佳实践

### 1. 包名简洁

```go
// ✅ 好
package user
package order

// ❌ 不好
package userManagementService
package order_processing
```

### 2. 避免 util 包

```go
// ❌ 不要
package util
func ParseTime() {}
func FormatMoney() {}

// ✅ 要
package timeutil
package currency
```

### 3. 接口定义在消费者端

```go
// ✅ 好：消费者定义需要的方法
package storage
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 实现者在另一个包
package filesystem
type File struct{}
func (f *File) Read(p []byte) (n int, err error)
```
