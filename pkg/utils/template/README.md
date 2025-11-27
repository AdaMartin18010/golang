# 模板工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [模板工具](#模板工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 文本模板](#21-文本模板)
    - [2.2 HTML模板](#22-html模板)
    - [2.3 快捷函数](#23-快捷函数)
    - [2.4 常用函数映射](#24-常用函数映射)
  - [3. 使用示例](#3-使用示例)
    - [3.1 文本模板](#31-文本模板)
    - [3.2 HTML模板](#32-html模板)
    - [3.3 从文件解析模板](#33-从文件解析模板)
    - [3.4 快捷渲染](#34-快捷渲染)
    - [3.5 自定义函数](#35-自定义函数)
    - [3.6 模板验证](#36-模板验证)
    - [3.7 模板克隆和查找](#37-模板克隆和查找)
    - [3.8 完整示例](#38-完整示例)

---

## 1. 概述

模板工具提供了text/template和html/template的便捷封装，简化模板解析、执行和管理任务。

---

## 2. 功能特性

### 2.1 文本模板

- `TextTemplate`: 文本模板结构体
- `NewTextTemplate`: 创建新的文本模板
- `Parse`: 解析模板字符串
- `ParseFiles`: 解析模板文件
- `ParseGlob`: 解析匹配的模板文件
- `Execute`: 执行模板
- `ExecuteToString`: 执行模板并返回字符串
- `ExecuteToBytes`: 执行模板并返回字节数组
- `ExecuteToFile`: 执行模板并写入文件
- `AddFunc`: 添加自定义函数
- `AddFuncs`: 添加多个自定义函数
- `Clone`: 克隆模板
- `Lookup`: 查找命名模板
- `DefinedTemplates`: 获取所有定义的模板名称

### 2.2 HTML模板

- `HTMLTemplate`: HTML模板结构体
- `NewHTMLTemplate`: 创建新的HTML模板
- `Parse`: 解析模板字符串
- `ParseFiles`: 解析模板文件
- `ParseGlob`: 解析匹配的模板文件
- `Execute`: 执行模板
- `ExecuteToString`: 执行模板并返回字符串
- `ExecuteToBytes`: 执行模板并返回字节数组
- `ExecuteToFile`: 执行模板并写入文件
- `AddFunc`: 添加自定义函数
- `AddFuncs`: 添加多个自定义函数
- `Clone`: 克隆模板
- `Lookup`: 查找命名模板
- `DefinedTemplates`: 获取所有定义的模板名称

### 2.3 快捷函数

- `Render`: 渲染文本模板
- `RenderHTML`: 渲染HTML模板
- `RenderFile`: 从文件渲染模板
- `RenderHTMLFile`: 从文件渲染HTML模板
- `Validate`: 验证模板是否有效
- `ValidateHTML`: 验证HTML模板是否有效

### 2.4 常用函数映射

- `CommonFuncMap`: 文本模板常用函数映射
- `HTMLCommonFuncMap`: HTML模板常用函数映射

---

## 3. 使用示例

### 3.1 文本模板

```go
import "github.com/yourusername/golang/pkg/utils/template"

// 创建文本模板
tmpl := template.NewTextTemplate("test")
tmpl, err := tmpl.Parse("Hello, {{.Name}}!")

// 执行模板
data := map[string]string{"Name": "World"}
result, err := tmpl.ExecuteToString(data)
// 结果: "Hello, World!"

// 执行模板到文件
err := tmpl.ExecuteToFile(data, "output.txt")
```

### 3.2 HTML模板

```go
// 创建HTML模板
tmpl := template.NewHTMLTemplate("test")
tmpl, err := tmpl.Parse("<h1>Hello, {{.Name}}!</h1>")

// 执行模板
data := map[string]string{"Name": "World"}
result, err := tmpl.ExecuteToString(data)
// 结果: "<h1>Hello, World!</h1>"
```

### 3.3 从文件解析模板

```go
// 解析模板文件
tmpl := template.NewTextTemplate("test")
tmpl, err := tmpl.ParseFiles("template.txt")

// 解析匹配的模板文件
tmpl, err := tmpl.ParseGlob("templates/*.txt")
```

### 3.4 快捷渲染

```go
// 渲染文本模板
templateText := "Hello, {{.Name}}!"
data := map[string]string{"Name": "World"}
result, err := template.Render(templateText, data)

// 渲染HTML模板
htmlTemplate := "<h1>Hello, {{.Name}}!</h1>"
result, err := template.RenderHTML(htmlTemplate, data)

// 从文件渲染
result, err := template.RenderFile("template.txt", data)
result, err := template.RenderHTMLFile("template.html", data)
```

### 3.5 自定义函数

```go
// 添加单个函数
tmpl := template.NewTextTemplate("test")
tmpl = tmpl.AddFunc("upper", func(s string) string {
    return strings.ToUpper(s)
})
tmpl, err := tmpl.Parse("Hello, {{upper .Name}}!")

// 添加多个函数
funcMap := template.FuncMap{
    "upper": strings.ToUpper,
    "lower": strings.ToLower,
}
tmpl = tmpl.AddFuncs(funcMap)

// 使用常用函数映射
tmpl = tmpl.AddFuncs(template.CommonFuncMap)
tmpl, err := tmpl.Parse("{{add 1 2}}") // 结果: "3"
```

### 3.6 模板验证

```go
// 验证模板是否有效
err := template.Validate("Hello, {{.Name}}!")
if err != nil {
    // 模板无效
}

// 验证HTML模板
err := template.ValidateHTML("<h1>Hello, {{.Name}}!</h1>")
```

### 3.7 模板克隆和查找

```go
// 克隆模板
cloned, err := tmpl.Clone()

// 查找命名模板
namedTmpl := tmpl.Lookup("template_name")

// 获取所有定义的模板名称
templates := tmpl.DefinedTemplates()
```

### 3.8 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/template"
)

func main() {
    // 创建模板
    tmpl := template.NewTextTemplate("greeting")
    tmpl = tmpl.AddFuncs(template.CommonFuncMap)

    // 解析模板
    tmpl, err := tmpl.Parse(`
Hello, {{.Name}}!
Your age is {{.Age}}.
{{if gt .Age 18}}
You are an adult.
{{else}}
You are a minor.
{{end}}
`)
    if err != nil {
        panic(err)
    }

    // 执行模板
    data := map[string]interface{}{
        "Name": "John",
        "Age":  25,
    }

    result, err := tmpl.ExecuteToString(data)
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
}
```

---

**更新日期**: 2025-11-11
