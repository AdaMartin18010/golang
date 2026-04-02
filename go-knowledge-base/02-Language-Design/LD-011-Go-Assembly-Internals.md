# LD-011: Go 汇编与运行时内部 (Go Assembly & Runtime Internals)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #go-assembly #runtime #plan9-asm #go-internals
> **权威来源**: [A Quick Guide to Go's Assembler](https://go.dev/doc/asm), [Go Runtime](https://go.dev/src/runtime/)

---

## Go 汇编概述

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Go Assembly Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Go 使用 Plan 9 汇编语法，与 GNU 汇编不同                                     │
│                                                                              │
│  编译流程:                                                                    │
│  Go Source ──► Go Compiler ──► SSA/IR ──► Machine Code                      │
│      │                           │                                          │
│      │ (内联汇编)                  │ (go tool objdump)                       │
│      ▼                           ▼                                          │
│  .s 文件 (Plan 9 汇编)    汇编输出分析                                       │
│                                                                              │
│  文件命名:                                                                    │
│  - foo_amd64.s: AMD64 架构汇编                                               │
│  - foo_arm64.s: ARM64 架构汇编                                               │
│  - foo.s: 通用汇编 (Go 代码)                                                  │
│                                                                              │
│  寄存器命名 (伪寄存器):                                                        │
│  - FP: Frame pointer (参数/局部变量)                                          │
│  - PC: Program counter                                                        │
│  - SB: Static base (全局符号)                                                 │
│  - SP: Stack pointer                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 汇编基础

### 基本语法

```asm
// add_amd64.s
// func Add(a, b int64) int64

TEXT ·Add(SB), NOSPLIT, $0-24
    // 函数序言
    // $0-24: 栈帧大小 0，参数+返回值大小 24 字节

    MOVQ a+0(FP), AX    // 第一个参数 → AX
    MOVQ b+8(FP), BX    // 第二个参数 → BX

    ADDQ BX, AX         // AX = AX + BX

    MOVQ AX, ret+16(FP) // 结果 → 返回值位置

    RET                 // 返回

// 说明:
// ·Add: 包级符号 (中间点)
// (SB): Static base，表示这是符号地址
// a+0(FP): 参数 a 在 FP+0 的位置
// ret+16(FP): 返回值在 FP+16 的位置
```

### 栈帧布局

```
高地址
┌─────────────────────┐
│    返回地址         │  ← 调用者的返回地址
├─────────────────────┤
│    调用者 BP        │  ← 保存的基址指针 (如果 FRAME_SIZE > 0)
├─────────────────────┤
│    局部变量...      │  ← $FRAME_SIZE 指定的大小
├─────────────────────┤
│    参数 1           │  ← a+0(FP)
│    参数 2           │  ← b+8(FP)
├─────────────────────┤
│    返回值 1         │  ← ret+16(FP)
├─────────────────────┤
│    返回地址         │  ← 0(SP) 在函数内部视角
├─────────────────────┤
│    局部栈空间       │  ← 由 $FRAME_SIZE 分配
└─────────────────────┘
低地址
```

---

## 运行时原语

### Goroutine 切换

```asm
// runtime·gogo: 切换到 goroutine g
TEXT runtime·gogo(SB), NOSPLIT, $0-8
    MOVQ gg+0(FP), BX        // 加载 g
    MOVQ gobuf_g(BX), DX     // 加载 g 指针
    MOVQ DX, g(CX)           // 设置 TLS 中的 g

    MOVQ gobuf_sp(BX), SP    // 恢复 SP
    MOVQ gobuf_pc(BX), AX    // 恢复 PC
    MOVQ gobuf_bp(BX), BP    // 恢复 BP

    MOVQ $0, gobuf_sp(BX)    // 清空 gobuf
    MOVQ $0, gobuf_pc(BX)
    MOVQ $0, gobuf_bp(BX)

    JMP AX                   // 跳转到 goroutine 入口
```

### 系统调用

```asm
// Syscall 包装
TEXT ·Syscall(SB), NOSPLIT, $0-56
    MOVQ a1+8(FP), DI        // 参数 1
    MOVQ a2+16(FP), SI       // 参数 2
    MOVQ a3+24(FP), DX       // 参数 3
    MOVQ trap+0(FP), AX      // 系统调用号

    SYSCALL                  // 执行系统调用

    CMPQ AX, $0xfffffffffffff001
    JLS ok                   // 无错误

    // 错误处理
    NEGQ AX
    MOVQ AX, err+48(FP)
    MOVQ $-1, r1+32(FP)
    MOVQ $0, r2+40(FP)
    RET

ok:
    MOVQ AX, r1+32(FP)       // 返回值 1
    MOVQ DX, r2+40(FP)       // 返回值 2
    MOVQ $0, err+48(FP)      // 无错误
    RET
```

---

## 原子操作实现

```asm
// atomic_amd64.s
// 基于 x86-64 LOCK 前缀的原子操作

TEXT ·AddInt64(SB), NOSPLIT, $0-24
    MOVQ addr+0(FP), BP      // &val
    MOVQ delta+8(FP), AX     // 要加的值

    LOCK                     // LOCK 前缀保证原子性
    XADDQ AX, 0(BP)          // 交换并加

    MOVQ AX, ret+16(FP)      // 返回旧值
    RET

TEXT ·CompareAndSwapInt64(SB), NOSPLIT, $0-32
    MOVQ addr+0(FP), BP
    MOVQ old+8(FP), AX       // 期望值
    MOVQ new+16(FP), CX      // 新值

    LOCK
    CMPXCHGQ CX, 0(BP)       // 比较并交换

    SETEQ AL                 // 设置布尔结果
    MOVZX AL, AX
    MOVQ AX, ret+24(FP)
    RET
```

---

## 使用场景

| 场景 | 汇编用途 | 示例 |
|------|---------|------|
| 极致性能 | 关键路径优化 | crypto, math |
| 底层操作 | 直接硬件访问 | runtime 内存管理 |
| 原子操作 | CPU 原子指令 | sync/atomic |
| SIMD | 向量化计算 | image, crypto |
| 系统调用 | 无 C 库依赖 | syscall 包 |

---

## 参考文献

1. [A Quick Guide to Go's Assembler](https://go.dev/doc/asm)
2. [Plan 9 Assembler](https://9p.io/sys/doc/asm.pdf)
3. [Go Runtime Source](https://go.dev/src/runtime/)
