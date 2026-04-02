# LD-014: Go 汇编编程与底层接口 (Go Assembly Programming & Low-Level Interface)

> **维度**: Language Design
> **级别**: S (20+ KB)
> **标签**: #assembly #plan9-asm #runtime #syscall #inline-asm #low-level
> **权威来源**:
>
> - [A Quick Guide to Go's Assembler](https://go.dev/doc/asm) - Go Authors
> - [Plan 9 Assembler Manual](https://9p.io/sys/doc/asm.pdf) - Bell Labs
> - [Go Assembly by Example](https://github.com/teh-cmc/go-internals) - teh-cmc
> - [x86-64 ABI](https://github.com/hjl-tools/x86-psABI/wiki/X86-psABI) - System V AMD64 ABI
> - [ARM64 ABI](https://developer.arm.com/documentation/ihi0055/b/) - ARM Architecture

---

## 1. 形式化基础

### 1.1 汇编语言理论

**定义 1.1 (汇编语言)**
汇编语言是机器指令的符号表示：

$$\text{Assembly} = \{ \text{Instructions}, \text{Directives}, \text{Labels}, \text{Comments} \}$$

**定义 1.2 (指令格式)**

$$\text{Instruction} ::= \text{Opcode} \quad \text{Operands}$$

**定义 1.3 (Plan 9 汇编语法)**

$$\text{Destination} \leftarrow \text{Source}$$

与 Intel 语法相反：

- Plan 9: `MOVQ src, dst`
- Intel: `MOV dst, src`

### 1.2 寄存器约定

**定义 1.4 (AMD64 寄存器)**

| 寄存器 | 用途 | Callee-saved |
|--------|------|--------------|
| RAX | 返回值 | No |
| RBX | 通用 | Yes |
| RCX | 第4参数 | No |
| RDX | 第3参数 | No |
| RSI | 第2参数 | No |
| RDI | 第1参数 | No |
| RBP | 帧指针 | Yes |
| RSP | 栈指针 | Yes |
| R8-R11 | 第5-8参数 | No |
| R12-R15 | 通用 | Yes |

**定义 1.5 (ARM64 寄存器)**

| 寄存器 | 用途 |
|--------|------|
| X0-X7 | 参数/返回值 |
| X8 | 间接结果 |
| X9-X15 | 临时 |
| X19-X28 | 被调用者保存 |
| X29 | 帧指针 |
| X30 | 链接寄存器 |
| SP | 栈指针 |

---

## 2. Go 汇编语法

### 2.1 文件结构

**定义 2.1 (汇编文件格式)**

```asm
// 文件: example_amd64.s

#include "textflag.h"    // 包含标志定义

// 数据段
DATA ·pi+0(SB)/8, $3.14159265359
GLOBL ·pi(SB), RODATA, $8

// 代码段
TEXT ·FunctionName(SB), NOSPLIT, $frameSize-argumentSize
    // 指令
    RET
```

**定义 2.2 (TEXT 指令)**

```asm
TEXT symbol(SB), flags, $framesize-argumentsize
```

- `symbol`: 函数名（· 是包名分隔符）
- `flags`: NOSPLIT, WRAPPER, NEEDCTXT 等
- `framesize`: 栈帧大小
- `argumentsize`: 参数+返回值大小

### 2.2 伪寄存器

**定义 2.3 (Go 伪寄存器)**

| 伪寄存器 | 含义 | 用途 |
|----------|------|------|
| SB | Static Base | 全局符号地址 |
| FP | Frame Pointer | 函数参数/返回值 |
| PC | Program Counter | 当前指令地址 |
| SP | Stack Pointer | 栈顶（本地变量偏移） |

**定义 2.4 (内存引用语法)**

```asm
offset(DI)(index*scale)

// 示例:
0(FP)           // FP + 0
8(SP)           // SP + 8
16(DI)(AX*4)    // DI + 16 + AX*4
·global(SB)     // 全局符号
```

### 2.3 指令集

**定义 2.5 (数据移动)**

```asm
MOVB src, dst    // 8-bit move
MOVW src, dst    // 16-bit move
MOVL src, dst    // 32-bit move
MOVQ src, dst    // 64-bit move
MOVUPS src, dst  // 128-bit unaligned move (SSE)
MOVAPS src, dst  // 128-bit aligned move (SSE)
```

**定义 2.6 (算术运算)**

```asm
// 加减
ADDQ a, b        // b = b + a
SUBQ a, b        // b = b - a

// 乘除
IMULQ a, b       // b = b * a (signed)
MULQ a           // RDX:RAX = RAX * a (unsigned)
CQO; IDIVQ a     // RAX = RDX:RAX / a, RDX = RDX:RAX % a

// 位运算
ANDQ a, b
ORQ a, b
XORQ a, b
NOTQ a
SHLQ n, a        // a = a << n
SHRQ n, a        // a = a >> n (logical)
SARQ n, a        // a = a >> n (arithmetic)
```

**定义 2.7 (比较与跳转)**

```asm
// 比较
CMPQ a, b        // compare a and b
TESTQ a, b       // a & b

// 条件跳转
JE label         // jump if equal
JNE label        // jump if not equal
JL label         // jump if less (signed)
JLE label        // jump if less or equal
JG label         // jump if greater
JGE label        // jump if greater or equal
JB label         // jump if below (unsigned)
JA label         // jump if above (unsigned)

// 无条件跳转
JMP label
CALL symbol
RET
```

---

## 3. 函数实现

### 3.1 简单函数

**定义 3.1 (无栈帧函数)**

```go
// Go 声明
func Add(a, b int64) int64
```

```asm
// amd64 实现
TEXT ·Add(SB), NOSPLIT, $0-24
    // 参数位置: a at 0(FP), b at 8(FP)
    // 返回值位置: ret at 16(FP)

    MOVQ a+0(FP), AX    // AX = a
    MOVQ b+8(FP), BX    // BX = b
    ADDQ BX, AX         // AX = AX + BX
    MOVQ AX, ret+16(FP) // ret = AX
    RET
```

**定义 3.2 (带栈帧函数)**

```go
// Go 声明
func Sum(n int64) int64
```

```asm
// amd64 实现
TEXT ·Sum(SB), $16-16
    // $16: 栈帧大小 (用于局部变量)
    // -16: 参数8字节 + 返回值8字节

    MOVQ n+0(FP), CX    // CX = n
    MOVQ $0, AX         // AX = sum = 0
    MOVQ $1, DX         // DX = i = 1

loop:
    CMPQ DX, CX
    JGT done
    ADDQ DX, AX         // sum += i
    ADDQ $1, DX         // i++
    JMP loop

done:
    MOVQ AX, ret+8(FP)  // 返回 sum
    RET
```

### 3.3 调用约定

**定义 3.3 (Go 调用约定)**

```
参数传递:
- AMD64: 栈传递 (通过 FP)
- 返回值: 栈传递

寄存器使用:
- 调用者保存: AX, CX, DX, DI, SI, R8-R11, X0-X15
- 被调用者保存: BX, BP, R12-R15

栈对齐:
- AMD64: 16-byte aligned at function entry
```

**定义 3.4 (系统调用)**

```asm
// Linux syscall
TEXT ·Syscall(SB), NOSPLIT, $0-56
    // a1+8(FP), a2+16(FP), a3+24(FP), trap+0(FP)

    MOVQ a1+8(FP), DI
    MOVQ a2+16(FP), SI
    MOVQ a3+24(FP), DX
    MOVQ trap+0(FP), AX

    SYSCALL

    MOVQ AX, r1+32(FP)
    MOVQ DX, r2+40(FP)
    MOVQ CX, err+48(FP)  // -errno
    RET
```

---

## 4. 运行时集成

### 4.1 栈增长检查

**定义 4.1 (栈分裂)**

```asm
TEXT ·FunctionWithStackCheck(SB), $64-0
    // 栈增长检查
    MOVQ (TLS), CX          // g
    CMPQ SP, 16(CX)         // compare with g.stackguard0
    JLS morestack           // if SP <= stackguard, grow stack

    // 函数体
    // ...
    RET

morestack:
    CALL runtime·morestack(SB)
    JMP ·FunctionWithStackCheck(SB)
```

### 4.2 与 Go 代码互操作

**定义 4.2 (汇编调用 Go 函数)**

```asm
TEXT ·CallGoFunc(SB), NOSPLIT, $8-0
    // 准备参数
    MOVQ $42, 0(SP)

    // 调用 Go 函数
    CALL ·someGoFunc(SB)

    // 处理返回值
    // ...
    RET
```

**定义 4.3 (Go 调用汇编)**

```go
// Go 代码
package asm

//go:noescape
func AsmFunc(a, b int64) int64

//go:linkname runtimeMemmove runtime.memmove
func runtimeMemmove(dst, src unsafe.Pointer, n uintptr)
```

---

## 5. 多元表征

### 5.1 汇编语法对比

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Assembly Syntax Comparison                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  特性              │ Plan 9 (Go)    │ Intel          │ AT&T              │
│  ──────────────────┼────────────────┼────────────────┼───────────────────│
│  操作数顺序         │ DST, SRC       │ DST, SRC       │ SRC, DST          │
│  寄存器前缀         │ 无             │ 无             │ %                 │
│  立即数前缀         │ $              │ 无             │ $                 │
│  内存引用           │ 8(SP)          │ [rsp+8]        │ 8(%rsp)           │
│  注释               │ // 或 /* */   │ ;              │ # 或 /* */        │
│  指令后缀           │ Q, L, W, B     │ 无 (操作数决定) │ 无 (操作数决定)    │
│  符号引用           │ ·sym(SB)       │ sym            │ sym               │
│                                                                              │
│  示例: MOV 指令                                                               │
│  ─────────────                                                               │
│  Plan 9:   MOVQ $42, AX                                                      │
│  Intel:    mov rax, 42                                                       │
│  AT&T:     movq $42, %rax                                                    │
│                                                                              │
│  示例: 内存操作                                                               │
│  ─────────────                                                               │
│  Plan 9:   MOVQ 8(SP), AX                                                    │
│  Intel:    mov rax, [rsp+8]                                                  │
│  AT&T:     movq 8(%rsp), %rax                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 函数实现模板

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Assembly Function Templates                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  简单函数 (无栈帧)                                                            │
│  ═════════════════                                                               │
│  // func Add(a, b int64) int64                                               │
│  TEXT ·Add(SB), NOSPLIT, $0-24                                               │
│      MOVQ a+0(FP), AX                                                        │
│      MOVQ b+8(FP), BX                                                        │
│      ADDQ BX, AX                                                             │
│      MOVQ AX, ret+16(FP)                                                     │
│      RET                                                                     │
│                                                                              │
│  带局部变量的函数                                                             │
│  ════════════════════                                                        │
│  // func Sum(n int) int                                                      │
│  TEXT ·Sum(SB), $16-16                                                       │
│      MOVQ n+0(FP), CX      // 参数                                           │
│      MOVQ $0, AX           // sum                                            │
│      MOVQ $0, -8(SP)       // 局部变量 i                                     │
│                                                                              │
│  loop:                                                                       │
│      CMPQ -8(SP), CX                                                         │
│      JGE done                                                                │
│      ADDQ -8(SP), AX                                                         │
│      ADDQ $1, -8(SP)                                                         │
│      JMP loop                                                                │
│                                                                              │
│  done:                                                                       │
│      MOVQ AX, ret+8(FP)                                                      │
│      RET                                                                     │
│                                                                              │
│  循环结构                                                                    │
│  ═══════════                                                                │
│      MOVQ $0, CX           // i = 0                                          │
│                                                                              │
│  loop:                                                                       │
│      CMPQ CX, $10                                                            │
│      JGE done                                                                │
│      // loop body                                                            │
│      ADDQ $1, CX                                                             │
│      JMP loop                                                                │
│  done:                                                                       │
│                                                                              │
│  条件分支                                                                    │
│  ═══════════                                                                │
│      CMPQ AX, BX                                                             │
│      JE equal                                                                │
│      JG greater                                                              │
│      // less than                                                            │
│      JMP end                                                                 │
│  equal:                                                                      │
│      // equal case                                                           │
│      JMP end                                                                 │
│  greater:                                                                    │
│      // greater case                                                         │
│  end:                                                                        │
│                                                                              │
│  系统调用                                                                    │
│  ═══════════                                                                │
│  TEXT ·getpid(SB), NOSPLIT, $0-8                                             │
│      MOVQ $39, AX          // SYS_getpid                                     │
│      SYSCALL                                                                 │
│      MOVQ AX, ret+0(FP)                                                      │
│      RET                                                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 平台支持矩阵

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Platform Support Matrix                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  平台           │ 文件后缀      │ 寄存器大小 │ 对齐要求   │ 特殊特性       │
│  ───────────────┼───────────────┼────────────┼────────────┼────────────────│
│  amd64          │ _amd64.s      │ 64-bit     │ 16-byte    │ SSE, AVX       │
│  386            │ _386.s        │ 32-bit     │ 4-byte     │ x87 FPU        │
│  arm64          │ _arm64.s      │ 64-bit     │ 16-byte    │ NEON           │
│  arm            │ _arm.s        │ 32-bit     │ 4-byte     │ VFP            │
│  ppc64le        │ _ppc64le.s    │ 64-bit     │ 16-byte    │ Altivec        │
│  s390x          │ _s390x.s      │ 64-bit     │ 8-byte     │ Vector         │
│  wasm           │ _wasm.s       │ 32-bit     │ 8-byte     │ 线性内存       │
│                                                                              │
│  文件名规则:                                                                  │
│  • foo_amd64.s    - AMD64 特定                                              │
│  • foo_arm64.s    - ARM64 特定                                              │
│  • foo.s          - 通用（最后选择）                                         │
│  • 在 _* 前缀文件存在时，foo.s 会被忽略                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.4 调试与工具

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Assembly Debugging Tools                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  生成汇编                                                                    │
│  ═══════════                                                                │
│  go tool compile -S file.go          # 显示汇编                             │
│  go build -gcflags="-S"              # 编译并显示汇编                        │
│  go tool objdump -s "FuncName" bin   # 反汇编特定函数                        │
│                                                                              │
│  调试标志                                                                    │
│  ═══════════                                                                │
│  TEXT ·Func(SB), NOSPLIT, $0         # 禁用栈分裂检查                        │
│  TEXT ·Func(SB), $0-0                # 标准栈分裂检查                        │
│  TEXT ·Func(SB), WRAPPER, $0         # 包装器函数                            │
│  TEXT ·Func(SB), NEEDCTXT, $0        # 需要上下文                            │
│                                                                              │
│  编译器指令                                                                  │
│  ═══════════                                                                │
│  //go:noescape                       # 禁用逃逸分析                          │
│  //go:nosplit                        # 禁用栈分裂                            │
│  //go:noinline                       # 禁用内联                              │
│  //go:linkname local remote          # 链接到远程符号                        │
│  //go:cgo_import_static name         # CGO 静态导入                          │
│  //go:cgo_import_dynamic name        # CGO 动态导入                          │
│                                                                              │
│  常见错误                                                                    │
│  ═══════════                                                                │
│  • 参数偏移计算错误                                                          │
│  • 栈帧大小不足                                                              │
│  • 忘记保存/恢复被调用者保存寄存器                                            │
│  • 栈未对齐                                                                  │
│  • 忘记 RET                                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. 代码示例与基准测试

### 6.1 汇编函数实现

```go
package asm

import (
    "unsafe"
)

// 纯汇编实现的加法
//go:noescape
func AsmAdd(a, b int64) int64

// 字节数组比较
//go:noescape
func AsmBytesEqual(a, b []byte) bool

// 快速内存清零
//go:noescape
func AsmMemclr(p unsafe.Pointer, n uintptr)

// 计算汉明权重 (population count)
//go:noescape
func AsmPopcnt(x uint64) int
```

```asm
// add_amd64.s
#include "textflag.h"

TEXT ·AsmAdd(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    MOVQ AX, ret+16(FP)
    RET

// bytes_equal_amd64.s
TEXT ·AsmBytesEqual(SB), NOSPLIT, $0-50
    // a.len == b.len ?
    MOVQ a_len+8(FP), AX
    MOVQ b_len+32(FP), BX
    CMPQ AX, BX
    JNE notequal

    // 空切片相等
    CMPQ AX, $0
    JE equal

    // 比较字节
    MOVQ a_ptr+0(FP), SI
    MOVQ b_ptr+24(FP), DI
    MOVQ AX, CX

    CLD
    REP; CMPSB
    JNE notequal

equal:
    MOVB $1, ret+48(FP)
    RET

notequal:
    MOVB $0, ret+48(FP)
    RET

// memclr_amd64.s
TEXT ·AsmMemclr(SB), NOSPLIT, $0-16
    MOVQ p+0(FP), DI
    MOVQ n+8(FP), CX

    // 使用 REP STOSB 清零
    XORQ AX, AX
    CLD
    REP; STOSB
    RET

// popcnt_amd64.s
TEXT ·AsmPopcnt(SB), NOSPLIT, $0-16
    MOVQ x+0(FP), AX
    // 使用 POPCNT 指令 (SSE4.2)
    BYTE $0xF3; BYTE $0x48; BYTE $0x0F; BYTE $0xB8; BYTE $0xC0
    // POPCNT AX, AX

    MOVQ AX, ret+8(FP)
    RET
```

### 6.2 高级模式

```asm
// atomic_amd64.s
#include "textflag.h"

// CAS (Compare And Swap)
TEXT ·Cas64(SB), NOSPLIT, $0-25
    MOVQ ptr+0(FP), BX
    MOVQ old+8(FP), AX
    MOVQ new+16(FP), CX
    LOCK; CMPXCHGQ CX, 0(BX)
    SETE ret+24(FP)
    RET

// 内存屏障
TEXT ·MemoryBarrier(SB), NOSPLIT, $0-0
    MFENCE
    RET

// 使用 SIMD (AVX2)
TEXT ·SumFloat64AVX2(SB), NOSPLIT, $0-24
    MOVQ data+0(FP), AX
    MOVQ n+16(FP), CX

    VXORPD Y0, Y0, Y0      // 清零累加器

loop:
    CMPQ CX, $4
    JL tail

    VADDPD (AX), Y0, Y0    // 4 个 double 相加
    ADDQ $32, AX
    SUBQ $4, CX
    JMP loop

tail:
    // 处理剩余元素
    VEXTRACTF128 $1, Y0, X1
    VADDPD X1, X0, X0
    VHADDPD X0, X0, X0

    MOVSD X0, ret+8(FP)
    VZEROUPPER
    RET
```

### 6.3 性能基准测试

```go
package asm_test

import (
    "bytes"
    "testing"
    "unsafe"
    "asm"
)

func BenchmarkAsmAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = asm.AsmAdd(int64(i), int64(i))
    }
}

func BenchmarkGoAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = int64(i) + int64(i)
    }
}

func BenchmarkAsmBytesEqual(b *testing.B) {
    a := make([]byte, 1024)
    c := make([]byte, 1024)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = asm.AsmBytesEqual(a, c)
    }
}

func BenchmarkGoBytesEqual(b *testing.B) {
    a := make([]byte, 1024)
    c := make([]byte, 1024)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = bytes.Equal(a, c)
    }
}

func BenchmarkAsmMemclr(b *testing.B) {
    buf := make([]byte, 1024*1024)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        asm.AsmMemclr(unsafe.Pointer(&buf[0]), uintptr(len(buf)))
    }
}

func BenchmarkGoMemclr(b *testing.B) {
    buf := make([]byte, 1024*1024)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for j := range buf {
            buf[j] = 0
        }
    }
}

func BenchmarkAsmPopcnt(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = asm.AsmPopcnt(uint64(i))
    }
}

func BenchmarkGoPopcnt(b *testing.B) {
    popcnt := func(x uint64) int {
        count := 0
        for x != 0 {
            count++
            x &= x - 1
        }
        return count
    }

    for i := 0; i < b.N; i++ {
        _ = popcnt(uint64(i))
    }
}
```

---

## 7. 关系网络

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Go Assembly Context                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  汇编语言家族                                                                │
│  ├── x86/x64 (Intel, AT&T, Plan 9)                                          │
│  ├── ARM (ARM, Thumb, AArch64)                                              │
│  ├── RISC-V                                                                 │
│  ├── MIPS                                                                   │
│  ├── PowerPC                                                                │
│  └── WebAssembly                                                            │
│                                                                              │
│  相关工具                                                                    │
│  ├── Assembler (cmd/asm)                                                    │
│  ├── Disassembler (go tool objdump)                                         │
│  ├── Debugger (delve, gdb)                                                  │
│  └── Profiler (pprof)                                                       │
│                                                                              │
│  使用场景                                                                    │
│  ├── 性能关键代码                                                           │
│  ├── 底层系统调用                                                           │
│  ├── SIMD 优化                                                              │
│  ├── 原子操作                                                               │
│  └── 与 C 库互操作                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. 参考文献

### 汇编语言

1. **Pike, R. et al.** Plan 9 from Bell Labs.
2. **Intel.** Intel 64 and IA-32 Architectures Software Developer's Manual.
3. **ARM.** ARM Architecture Reference Manual.

### Go 汇编

1. **Go Authors**. A Quick Guide to Go's Assembler.
2. **Go Authors**. Plan 9 C Compilers.
3. **Clement, M.** Go Assembly by Example.

---

**质量评级**: S (20+ KB)
**完成日期**: 2026-04-02
