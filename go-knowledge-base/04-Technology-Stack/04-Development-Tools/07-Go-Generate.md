# go:generate

> **分类**: 开源技术堆栈

---

## 基本用法

```go
//go:generate command arguments

// 示例
//go:generate go run gen.go
//go:generate protoc --go_out=. *.proto
```

---

## 常见用例

### 生成 Mock

```go
//go:generate mockgen -source=interface.go -destination=mock.go -package=mocks

type Store interface {
    Get(id string) (*User, error)
}
```

```bash
go generate ./...
```

### 生成 Stringer

```go
//go:generate stringer -type=Status

type Status int

const (
    Pending Status = iota
    Processing
    Completed
    Failed
)
```

### 生成 Protobuf

```go
//go:generate protoc --go_out=. --go-grpc_out=. service.proto
```

---

## 条件编译

```go
//go:build linux
// +build linux

package main

func init() {
    // Linux 特有代码
}
```

---

## 自定义生成器

```go
// gen.go
package main

import (
    "os"
    "text/template"
)

func main() {
    tmpl := `package main
    const Version = "{{.Version}}"
    `
    t := template.Must(template.New("").Parse(tmpl))
    t.Execute(os.Stdout, map[string]string{
        "Version": "1.0.0",
    })
}
```

```go
//go:generate go run gen.go
```
