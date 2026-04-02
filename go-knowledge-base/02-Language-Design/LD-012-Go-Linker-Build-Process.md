# LD-012: Go 链接器与构建流程 (Go Linker & Build Process)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #go-build #go-linker #compiler #build-process
> **权威来源**: [Go Build](https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching), [Go Linker Internals](https://go.dev/src/cmd/link/doc.go)

---

## Go 构建流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Go Build Pipeline                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Go Source Files (.go)                                                       │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Compiler (cmd/compile)                      │    │
│  │  1. 词法分析 (Lexer)                                                 │    │
│  │  2. 语法分析 (Parser) → AST                                          │    │
│  │  3. 类型检查 (Type Checker)                                          │    │
│  │  4. SSA 生成 (Static Single Assignment)                              │    │
│  │  5. 机器码生成                                                       │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│       │                                                                      │
│       ▼                                                                      │
│  Object Files (.o / .a)                                                      │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Linker (cmd/link)                           │    │
│  │  1. 符号解析                                                         │    │
│  │  2. 重定位                                                           │    │
│  │  3. 死代码消除                                                       │    │
│  │  4. 生成可执行文件                                                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│       │                                                                      │
│       ▼                                                                      │
│  Executable Binary                                                           │
│                                                                              │
│  缓存: $GOCACHE (编译/链接结果缓存)                                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 编译器阶段详解

### 1. 词法与语法分析

```go
// 示例代码
package main

func Add(a, b int) int {
    return a + b
}

// AST 结构 (简化)
type FuncDecl struct {
    Name   *Ident    // "Add"
    Type   *FuncType // func(int, int) int
    Body   *BlockStmt
}

type BinaryExpr struct {
    Op    token.Token // +
    X     Expr        // a
    Y     Expr        // b
}
```

### 2. SSA 中间表示

```
Go 代码:
func Add(a, b int) int {
    return a + b
}

SSA (Static Single Assignment):
b1:
    v1 = Add <int> [a, b] → 参数
    v2 = Add64 <int> v1.a v1.b
    Ret v2

SSA 优势:
- 每个变量只赋值一次
- 便于优化 (常量传播、死代码消除)
- 便于寄存器分配
```

### 3. 编译优化

```bash
# 默认优化级别
go build

# 禁用优化 (调试用)
go build -gcflags="-N -l"
# -N: 禁用优化
# -l: 禁用内联

# 查看 SSA
go build -gcflags="-S" main.go

# 查看汇编
go tool objdump -S binary_name
```

---

## 链接器详解

### 链接过程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Linker Stages                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  输入:                                                                        │
│  - 主包对象文件 (main.o)                                                      │
│  - 依赖包对象文件 (fmt.a, net/http.a, ...)                                    │
│  - 运行时 (runtime.a)                                                         │
│                                                                              │
│  步骤:                                                                        │
│                                                                              │
│  1. 符号解析 (Symbol Resolution)                                              │
│     - 收集所有符号定义和引用                                                  │
│     - 解析跨包引用                                                           │
│                                                                              │
│  2. 重定位 (Relocation)                                                      │
│     - 计算最终地址                                                           │
│     - 修正代码中的地址引用                                                    │
│                                                                              │
│  3. 死代码消除 (Dead Code Elimination)                                        │
│     - 标记可达函数                                                           │
│     - 删除未使用代码                                                         │
│                                                                              │
│  4. 生成可执行文件                                                            │
│     - ELF/Mach-O/PE 格式                                                     │
│     - 包含代码段、数据段、符号表                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 构建缓存

```bash
# 缓存位置
echo $GOCACHE  # 默认 ~/.cache/go-build

# 缓存键基于:
# - 源文件内容哈希
# - 编译器版本
# - 编译参数
# - 环境变量 (GOOS, GOARCH, CGO_ENABLED)

# 清理缓存
go clean -cache

# 查看缓存大小
go clean -cache -n  # 只显示，不删除
```

---

## 构建模式

### 交叉编译

```bash
# Linux AMD64 上构建 Windows 可执行文件
GOOS=windows GOARCH=amd64 go build -o app.exe

# 常用组合
GOOS=linux   GOARCH=amd64     # Linux x86_64
GOOS=linux   GOARCH=arm64     # Linux ARM64
GOOS=darwin  GOARCH=amd64     # macOS Intel
GOOS=darwin  GOARCH=arm64     # macOS Apple Silicon
GOOS=windows GOARCH=amd64     # Windows x86_64

# 使用 buildx (Go 1.20+)
go build -o app .
```

### 条件编译

```go
// +build linux

package main

// Linux 特定代码
```

```go
//go:build linux && amd64
// +build linux,amd64

package main

// Linux AMD64 特定代码
```

### Build Tags

```bash
# 使用 build tags
go build -tags=debug .
go build -tags="feature_a feature_b" .
```

---

## 构建优化

| 技术 | 方法 | 效果 |
|------|------|------|
| 并行构建 | `go build -p 8` | 使用 8 个并行进程 |
| 增量构建 | 缓存 | 只编译变化的包 |
| 链接优化 | `go build -ldflags="-s -w"` | 减小二进制体积 |
| PGO | Profile-Guided Optimization | 基于运行时数据优化 |

### 减小二进制体积

```bash
# 去除符号表和调试信息
go build -ldflags="-s -w" -o small_app

# 压缩 (使用 upx)
upx -9 small_app

# 分析二进制大小
go build -o app && ls -lh app
go tool nm app | wc -l  # 符号数量
```

---

## 参考文献

1. [Go Build Cache](https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching)
2. [Go Command Documentation](https://golang.org/doc/cmd)
3. [Go Linker Internals](https://go.dev/src/cmd/link/doc.go)
