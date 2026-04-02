# Go 构建模式 (Build Modes)

> **分类**: 开源技术堆栈  
> **标签**: #build #compilation #linking

---

## 常用构建命令

```bash
# 普通构建
go build

# 指定输出名称
go build -o myapp

# 交叉编译
go build -o myapp-linux GOOS=linux GOARCH=amd64

# 静态链接
go build -ldflags="-extldflags=-static" -o myapp

# 压缩二进制
go build -ldflags="-s -w" -o myapp
```

---

## 构建标签

### 条件编译

```go
//go:build linux
// +build linux

package main

func init() {
    println("Linux specific init")
}
```

### 多条件

```go
//go:build linux && amd64
//go:build (linux || darwin) && !windows
```

### 自定义标签

```go
//go:build debug

package main

import "log"

func debugLog(msg string) {
    log.Println("[DEBUG]", msg)
}
```

```bash
go build -tags debug
go build -tags "debug prod"
```

---

## 链接器标志

```bash
# 去除符号表和调试信息
go build -ldflags="-s -w"

# 版本信息注入
go build -ldflags="-X main.version=1.0.0 -X main.buildTime=$(date -u +%Y%m%d%H%M%S)"

# 减小二进制大小
go build -ldflags="-s -w -trimpath"

# 禁用 CGO
go build -ldflags="-linkmode external -extldflags=-static" -o myapp
```

---

## 构建模式

### 可执行文件

```bash
go build -buildmode=exe
```

### 共享库 (C Archive)

```go
// main.go
package main

import "C"

//export Add
func Add(a, b int) int {
    return a + b
}

func main() {}
```

```bash
go build -buildmode=c-archive -o libadd.a
```

### 共享库 (C Shared)

```bash
go build -buildmode=c-shared -o libadd.so
```

### 插件 (Plugin)

```go
// plugin.go
package main

type MyPlugin struct{}

func (p *MyPlugin) Do() string {
    return "hello from plugin"
}

var Plugin MyPlugin
```

```bash
go build -buildmode=plugin -o plugin.so
```

```go
// 主程序加载插件
p, err := plugin.Open("plugin.so")
sym, err := p.Lookup("Plugin")
myPlugin := sym.(*MyPlugin)
```

---

## 交叉编译矩阵

| GOOS | GOARCH | 说明 |
|------|--------|------|
| linux | amd64 | Linux x86_64 |
| linux | arm64 | Linux ARM64 |
| darwin | amd64 | macOS Intel |
| darwin | arm64 | macOS M1/M2 |
| windows | amd64 | Windows x64 |
| windows | 386 | Windows x86 |
| freebsd | amd64 | FreeBSD |
| js | wasm | WebAssembly |

```bash
# 构建所有平台
platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

for platform in "${platforms[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output="myapp-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output+=".exe"
    fi
    
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$output"
done
```

---

## 优化构建

### 使用 cache

```bash
# 自动使用 $GOCACHE
go build

# 清空缓存
go clean -cache
```

### 并行构建

```bash
# -p 并行度
go build -p 8
```

### Bazel 构建

```python
# BUILD
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "myapp",
    srcs = ["main.go"],
    deps = ["//pkg/mypackage"],
)
```
