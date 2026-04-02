# LD-002: Go 编译器架构与 SSA (Go Compiler Architecture & SSA)

> **维度**: Language Design
> **级别**: S (20+ KB)
> **标签**: #go-compiler #ssa #ir #optimization
> **权威来源**: [Go Compiler Introduction](https://github.com/golang/go/blob/master/src/cmd/compile/README.md), [SSA in Go](https://go.dev/src/cmd/compile/internal/ssa/README), [Go 1.5 Compiler](https://talks.golang.org/2015/gogo.slide)

---

## 编译器架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go Compiler Pipeline                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Source Code                                                                  │
│       │                                                                       │
│       ▼                                                                       │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Frontend (cmd/compile)                          │   │
│  │  ──────────────────────────────────────────────────────────────────  │   │
│  │                                                                      │   │
│  │  1. Lexical Analysis (scanner)                                       │   │
│  │     └── Tokens: package, import, func, etc.                         │   │
│  │                                                                      │   │
│  │  2. Syntax Analysis (parser)                                         │   │
│  │     └── AST (Abstract Syntax Tree)                                  │   │
│  │                                                                      │   │
│  │  3. Type Checking                                                    │   │
│  │     └── Type inference, validation                                  │   │
│  │                                                                      │   │
│  │  4. Escape Analysis                                                  │   │
│  │     └── Stack vs heap allocation decisions                          │   │
│  │                                                                      │   │
│  │  5. IR Generation (cmd/compile/internal/ir)                         │   │
│  │     └── Static Single Assignment (SSA) form                         │   │
│  │                                                                      │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│       │                                                                       │
│       ▼                                                                       │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Middle End (SSA passes)                         │   │
│  │  ──────────────────────────────────────────────────────────────────  │   │
│  │                                                                      │   │
│  │  SSA Form ──► Optimization Passes ──► Lowering                      │   │
│  │                                                                      │   │
│  │  Optimizations:                                                      │   │
│  │  • Dead code elimination                                             │   │
│  │  • Constant folding                                                  │   │
│  │  • Bounds check elimination                                          │   │
│  │  • Nil check elimination                                             │   │
│  │  • Inlining                                                          │   │
│  │  • Escape analysis optimizations                                     │   │
│  │                                                                      │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│       │                                                                       │
│       ▼                                                                       │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      Backend (codegen)                               │   │
│  │  ──────────────────────────────────────────────────────────────────  │   │
│  │                                                                      │   │
│  │  Machine-specific lowering ──► Register allocation ──► Assembly     │   │
│  │                                                                      │   │
│  │  Architectures: amd64, arm64, wasm, etc.                             │   │
│  │                                                                      │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│       │                                                                       │
│       ▼                                                                       │
│  Machine Code                                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## SSA 中间表示

### 什么是 SSA

**Static Single Assignment**: 每个变量只被赋值一次。

```go
// 原始代码
func max(a, b int) int {
    m := a
    if b > a {
        m = b  // 重新赋值，违反 SSA
    }
    return m
}

// SSA 形式
func max_ssa(a, b int) int {
    m1 := a
    b_gt_a := b > a

    // Phi 函数：根据控制流选择值
    m2 := Phi(m1, b)  // if b_gt_a then b else m1

    return m2
}
```

### SSA 值定义

```go
// src/cmd/compile/internal/ssa/value.go

// Value 是 SSA 的基本单元
type Value struct {
    Op   Op          // 操作码
    Type *types.Type // 类型

    Args []*Value    // 参数（依赖的其他 Value）

    ID   int32       // 唯一 ID
    Block *Block     // 所属的 Basic Block

    // 辅助信息
    AuxInt int64     // 整数辅助信息
    Aux    interface{} // 通用辅助信息（如符号名）

    // 寄存器分配结果
    Reg int16

    // 位置信息（用于调试）
    Pos src.XPos
}

// 操作码示例
const (
    OpAdd64    Op = iota  // 64位加法
    OpSub64               // 64位减法
    OpMul64               // 64位乘法
    OpDiv64               // 64位除法

    OpLoad                // 内存加载
    OpStore               // 内存存储

    OpPhi                 // Phi 函数（SSA 核心）

    OpConst64             // 常量
    OpConstBool           // 布尔常量

    OpCall                // 函数调用
    OpTailCall            // 尾调用

    // ... 更多操作
)
```

### 基本块 (Basic Block)

```go
// src/cmd/compile/internal/ssa/block.go

// Block 是控制流图中的节点
type Block struct {
    Kind BlockKind    // 块类型

    // 控制流边
    Preds []Edge      // 前驱块
    Succs []Edge      // 后继块

    // SSA 值
    Values []*Value   // 该块中的指令

    // 控制值（决定分支）
    Control *Value    // 条件跳转的判断值

    ID int32          // 唯一 ID

    // 辅助信息
    Aux interface{}   // 跳转目标等
}

// 块类型
const (
    BlockPlain BlockKind = iota  // 直线执行（只有一个后继）
    BlockIf                      // 条件分支（两个后继）
    BlockCall                    // 函数调用
    BlockRet                     // 返回
    BlockRetJmp                  // 跳转返回
    BlockExit                    // 退出（panic 等）
    BlockDefer                   // defer
    BlockCheck                   // 边界检查失败等
)
```

---

## 编译器优化示例

### 1. 内联优化

```go
// 原始代码
func add(a, b int) int {
    return a + b
}

func main() {
    x := add(1, 2)
    println(x)
}

// 内联后
func main() {
    x := 1 + 2      // 直接替换函数调用
    println(x)
}

// 常量折叠后
func main() {
    x := 3          // 编译期计算
    println(x)
}
```

### 2. 逃逸分析

```go
// 原始代码
func newInt() *int {
    x := 42
    return &x  // x 逃逸到堆
}

func stackInt() int {
    x := 42
    return x   // x 在栈上
}

// 逃逸分析结果
// newInt: x escapes to heap
// stackInt: x does not escape
```

### 3. 边界检查消除

```go
// 原始代码
func sum(arr []int) int {
    s := 0
    for i := 0; i < len(arr); i++ {
        s += arr[i]  // 每次都需要边界检查
    }
    return s
}

// 优化后
func sum_opt(arr []int) int {
    s := 0
    n := len(arr)
    if n > 0 {
        // 证明 i < n 始终成立，消除检查
        ptr := &arr[0]
        end := &arr[n]
        for ptr < end {
            s += *ptr
            ptr++
        }
    }
    return s
}
```

---

## SSA 可视化

```bash
# 查看 SSA 各个阶段
go build -gcflags="-m -m" main.go        # 逃逸分析
go tool compile -S main.go               # 汇编输出
GOSSAFUNC=main go build main.go          # 生成 SSA HTML
```

### SSA HTML 结构

```html
<!-- ssa.html -->
<html>
<body>
    <ul>
        <li><a href="#start">start</a> - 初始 IR</li>
        <li><a href="#deadcode">deadcode</a> - 死代码消除</li>
        <li><a href="#opt">opt</a> - 通用优化</li>
        <li><a href="#lower">lower</a> - 机器无关优化</li>
        <li><a href="#regalloc">regalloc</a> - 寄存器分配</li>
        <li><a href="#genssa">genssa</a> - 最终 SSA</li>
    </ul>

    <pre id="start">
b1:
    v1 = InitMem <mem>
    v2 = SP <uintptr>
    v3 = SB <uintptr>
    ...
    </pre>
</body>
</html>
```

---

## 参考文献

1. [cmd/compile/README](https://github.com/golang/go/blob/master/src/cmd/compile/README.md) - Go Compiler Documentation
2. [SSA in Go](https://go.dev/src/cmd/compile/internal/ssa/README) - SSA Package Documentation
3. [Go 1.5 Compiler](https://talks.golang.org/2015/gogo.slide) - Rob Pike
4. [Efficiently Computing Static Single Assignment Form](https://www.cs.utexas.edu/~pingali/CS380C/2010/papers/ssaCytron.pdf) - Cytron et al.
5. [Static Single Assignment Book](https://pfalcon.github.io/ssabook/latest/book-full.pdf) - SSA Book
