# CLI工具模板

**难度**: 入门 | **预计阅读**: 10分钟

---

## 📋 目录

- [1. 📖 CLI结构](#1--cli结构)
- [2. 📚 相关资源](#2--相关资源)

---

## 1. 📖 CLI结构

```
mycli/
├── cmd/
│   ├── root.go
│   ├── create.go
│   └── list.go
├── internal/
│   └── logic/
├── main.go
└── go.mod
```

---

## 🎯 Cobra CLI

```go
// cmd/root.go
package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "A brief description",
    Long:  `A longer description`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(createCmd)
    rootCmd.AddCommand(listCmd)
}

// cmd/create.go
var createCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new item",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        fmt.Printf("Creating: %s\n", name)
    },
}

// cmd/list.go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all items",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Listing items...")
    },
}

// main.go
func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

---

## 🚀 使用示例

```bash
# 构建
go build -o mycli

# 使用
./mycli create myitem
./mycli list
./mycli --help
```

---

## 📚 相关资源

- [Cobra](https://github.com/spf13/cobra)

**下一步**: [05-库项目模板](./05-库项目模板.md)

---

**最后更新**: 2025-10-28

