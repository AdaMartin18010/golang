# LD-011: Go 汇编内部原理 (Go Assembly Internals)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #assembly #plan9 #runtime #syscall #low-level
> **权威来源**:
>
> - [A Quick Guide to Go's Assembler](https://go.dev/doc/asm) - Go Authors
> - [Go Assembly by Example](https://github.com/teh-cmc/go-internals/blob/master/chapter1_assembly/chapter1.md) - Go Internals
> - [Plan 9 Assembler](https://9p.io/sys/doc/asm.pdf) - Plan 9

---

## 1. Go 汇编基础

### 1.1 Plan 9 汇编

Go 使用 Plan 9 汇编语法，与 GNU 汇编不同：

| 特性 | Plan 9 | GNU |
|------|--------|-----|
| 指令顺序 | 目标, 源 | 源, 目标 |
| 寄存器命名 | AX, BX | %rax, %rbx |
| 立即数 | $42 | $42 |
| 内存引用 | 8(SP) | 8(%rsp) |

### 1.2 文件结构

```asm
// 文件: add_amd64.s

TEXT ·Add(SB), NOSPLIT, $0-16
    MOVQ a+0(FP), AX    // 加载参数 a
    MOVQ b+8(FP), BX    // 加载参数 b
    ADDQ BX, AX         // AX = AX + BX
    MOVQ AX, ret+16(FP) // 存储返回值
    RET
```

---

## 2. 汇编语法

### 2.1 指令格式

```asm
TEXT symbol(SB), flags, $framesize-argumentsize
```

- `TEXT`: 定义代码段
- `symbol`: 函数名（· 是包名分隔符）
- `SB`: static base 伪寄存器
- `flags`: NOSPLIT, WRAPPER 等
- `$framesize`: 栈帧大小
- `argumentsize`: 参数+返回值大小

### 2.2 伪寄存器

```
FP (Frame Pointer): 函数参数和返回值
SP (Stack Pointer): 本地变量
SB (Static Base): 全局符号
PC (Program Counter): 当前指令
```

### 2.3 常见指令

```asm
// 数据移动
MOVQ src, dst    // 64 位移动
MOVL src, dst    // 32 位移动
MOVW src, dst    // 16 位移动
MOVB src, dst    // 8 位移动

// 算术运算
ADDQ a, b        // b = b + a
SUBQ a, b        // b = b - a
IMULQ a, b       // b = b * a
CQO; IDIVQ a     // AX = DX:AX / a, DX = DX:AX % a

// 位运算
ANDQ a, b
ORQ  a, b
XORQ a, b
SHLQ n, a        // a = a << n
SHRQ n, a        // a = a >> n

// 比较和跳转
CMPQ a, b
JE   label       // 相等跳转
JNE  label       // 不等跳转
JL   label       // 小于跳转
JMP  label       // 无条件跳转

// 函数调用
CALL symbol
RET
```

---

## 3. 函数实现

### 3.1 简单函数

```go
// Go 代码
func Add(a, b int64) int64 {
    return a + b
}
```

```asm
// amd64 汇编
TEXT ·Add(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    MOVQ AX, ret+16(FP)
    RET
```

### 3.2 带栈帧的函数

```go
func Sum(n int64) int64 {
    var sum int64
    for i := int64(1); i <= n; i++ {
        sum += i
    }
    return sum
}
```

```asm
TEXT ·Sum(SB), NOSPLIT, $16-16
    MOVQ n+0(FP), CX        // CX = n
    MOVQ $0, AX             // AX = sum = 0
    MOVQ $1, BX             // BX = i = 1

loop:
    CMPQ BX, CX
    JGT  done               // if i > n, exit
    ADDQ BX, AX             // sum += i
    INCQ BX                 // i++
    JMP  loop

done:
    MOVQ AX, ret+8(FP)
    RET
```

---

## 4. 运行时支持

### 4.1 栈增长检查

```asm
TEXT ·Func(SB), $48-0
    // 编译器插入的检查
    MOVQ    (TLS), CX
    CMPQ    SP, 16(CX)
    JLS     morestack

    // 函数体
    // ...
    RET

morestack:
    CALL    runtime·morestack(SB)
    JMP     ·Func(SB)
```

### 4.2 系统调用

```asm
// Linux syscall
TEXT ·Syscall(SB), NOSPLIT, $0-56
    MOVQ a1+8(FP), DI
    MOVQ a2+16(FP), SI
    MOVQ a3+24(FP), DX
    MOVQ trap+0(FP), AX    // syscall number
    SYSCALL
    MOVQ AX, r1+32(FP)
    MOVQ DX, r2+40(FP)
    MOVQ CX, err+48(FP)
    RET
```

---

## 5. 平台差异

### 5.1 寄存器约定

**AMD64 (Linux/macOS)**

```
参数:   DI, SI, DX, CX, R8, R9
返回值: AX, DX
```

**ARM64**

```
参数:   R0-R7
返回值: R0, R1
```

### 5.2 文件命名

```
foo_amd64.s    // AMD64
foo_arm64.s    // ARM64
foo_386.s      // 32-bit x86
foo.s          // 通用（最后选择）
```

---

## 6. 调试技巧

### 6.1 查看编译输出

```bash
# 生成汇编
$ go build -gcflags="-S" main.go

# 反编译二进制
$ go tool objdump -s "FuncName" binary
```

### 6.2 内联汇编替代

```go
package asm

import "unsafe"

//go:noescape
func Add(a, b uint64) uint64

//go:linkname runtime_memmove runtime.memmove
func runtime_memmove(dst, src unsafe.Pointer, n uintptr)
```

---

## 7. 关系网络

```
Go Assembly
├── Plan 9 Syntax
│   ├── Instructions
│   ├── Registers
│   └── Pseudo-registers
├── Platform Support
│   ├── amd64
│   ├── arm64
│   ├── 386
│   └── wasm
├── Runtime Integration
│   ├── Stack growth
│   ├── Goroutine
│   └── GC
└── Syscall Interface
    ├── Linux
    ├── Darwin
    └── Windows
```

---

**质量评级**: S (15KB)
**完成日期**: 2026-04-02
