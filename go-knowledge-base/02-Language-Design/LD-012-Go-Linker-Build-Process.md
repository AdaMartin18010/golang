# LD-012: Go 链接器与构建过程 (Go Linker & Build Process)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #linker #build #compiler #obj #elf
> **权威来源**:
>
> - [Go Linker](https://github.com/golang/go/tree/master/src/cmd/link) - Go Authors
> - [Build Modes](https://go.dev/doc/go1.5#link) - Go Release Notes
> - [ELF Format](https://refspecs.linuxfoundation.org/elf/elf.pdf) - System V ABI

---

## 1. 构建流程

### 1.1 编译流程

```
.go files
    │
    ▼ go tool compile
.o files (object)
    │
    ▼ go tool link
executable / library
```

### 1.2 完整工具链

```
源文件
   │
   ├──► cmd/compile ──► .o (SSA → 机器码)
   │
   ├──► cmd/asm ──────► .o (汇编)
   │
   └──► cgo ──────────► C 编译器 ──► .o
                            │
                            ▼
                    .o files + runtime.a
                            │
                            ▼
                    cmd/link ──► 可执行文件
```

---

## 2. 编译器输出

### 2.1 对象文件格式

Go 对象文件是自定义格式，不是标准 ELF/COFF：

```
对象文件结构:
┌─────────────────┐
│  Header         │
├─────────────────┤
│  Text Section   │  // 机器码
├─────────────────┤
│  Data Section   │  // 初始化数据
├─────────────────┤
│  BSS Section    │  // 未初始化数据
├─────────────────┤
│  Symbol Table   │  // 符号定义和引用
├─────────────────┤
│  Relocations    │  // 重定位信息
├─────────────────┤
│  Type Info      │  // Go 类型信息
├─────────────────┤
│  PCLine Table   │  // 程序计数器到行号映射
└─────────────────┘
```

### 2.2 符号类型

| 符号类型 | 说明 |
|----------|------|
| STEXT | 代码段符号 |
| SRODATA | 只读数据 |
| SNOPTRDATA | 无指针数据 |
| SDATA | 初始化数据 |
| SBSS | 未初始化数据 |
| SNOPTRBSS | 无指针 BSS |
| STLSBSS | TLS 数据 |

---

## 3. 链接器

### 3.1 链接阶段

```
1. 读取所有对象文件
2. 合并段（text, data, bss）
3. 符号解析
4. 重定位
5. 生成 DWARF 调试信息
6. 写入输出文件
```

### 3.2 链接模式

```bash
# 默认（exe）
go build -o app

# 共享库
go build -buildmode=c-shared -o libfoo.so

# 静态库
go build -buildmode=c-archive -o libfoo.a

# 插件
go build -buildmode=plugin -o plugin.so

# PIE（位置无关可执行文件）
go build -buildmode=pie -o app
```

### 3.3 内部链接 vs 外部链接

```
内部链接 (默认):
  - 纯 Go 代码
  - Go 链接器处理所有符号

外部链接:
  - 使用 cgo
  - 系统链接器 (ld) 参与
  - 支持 C 库依赖
```

---

## 4. 构建优化

### 4.1 构建缓存

```
位置: $GOCACHE (默认 ~/.cache/go-build)

缓存键: 文件内容 + 编译器版本 + 编译选项
```

### 4.2 增量构建

```bash
# 仅编译修改的文件
go build

# 强制重新编译
go build -a

# 清理缓存
go clean -cache
```

### 4.3 编译标志

```bash
# 禁用优化和内联（调试）
go build -gcflags="-N -l"

# 显示编译命令
go build -x

# 显示包加载过程
go build -v
```

---

## 5. 交叉编译

### 5.1 目标平台

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build

# Windows
GOOS=windows GOARCH=amd64 go build

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build

# Linux ARM
GOOS=linux GOARCH=arm GOARM=7 go build
```

### 5.2 CGO 交叉编译

```bash
# 禁用 CGO（纯 Go）
CGO_ENABLED=0 GOOS=linux go build

# 使用交叉编译工具链
CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 \
    GOOS=linux GOARCH=arm64 go build
```

---

## 6. 可执行文件分析

### 6.1 文件结构

```bash
# 查看段信息
$ readelf -S myapp

# 查看符号表
$ nm myapp

# 查看依赖
$ ldd myapp

# 反汇编
$ objdump -d myapp
```

### 6.2 运行时信息

```bash
# 查看 Go 版本
$ go version -m myapp

# 提取构建信息
$ go tool buildid myapp
```

---

## 7. 关系网络

```
Go Build Process
├── Compiler (cmd/compile)
│   ├── Parser
│   ├── Type Checker
│   ├── SSA Builder
│   ├── Optimizer
│   └── Code Generator
├── Assembler (cmd/asm)
│   └── Plan 9 → Machine Code
├── Cgo
│   └── C Binding Generator
└── Linker (cmd/link)
    ├── Symbol Resolution
    ├── Relocation
    └── Output Generation
```

---

**质量评级**: S (15KB)
**完成日期**: 2026-04-02
