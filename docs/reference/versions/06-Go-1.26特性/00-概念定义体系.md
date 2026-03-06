# Go 1.26 概念定义体系

## 核心概念

### 1. 语言特性

| 概念 | 定义 | 相关术语 |
|------|------|----------|
| `new(expr)` | 内置函数 `new` 支持表达式操作数，允许直接初始化值并返回指针 | 复合字面量、指针、堆分配 |
| 递归泛型约束 | 泛型类型可以在自身类型参数列表中引用自己，允许自引用类型约束 | 类型参数、类型约束、不动点 |
| Modernizer | 自动代码现代化分析器，将代码更新为使用新特性和 API | go fix、代码分析、自动化 |
| Inline 指令 | `//go:fix inline` 标记函数可被自动内联替换 | API 迁移、代码生成 |

### 2. 运行时

| 概念 | 定义 | 相关术语 |
|------|------|----------|
| Green Tea GC | Go 1.26 默认启用的新一代垃圾回收器，提供更低延迟和更高吞吐 | GC、并发标记、写屏障 |
| 切片栈分配 | 编译器优化，将小容量切片的底层数组分配在栈上 | 逃逸分析、栈、堆 |
| cgo 开销 | Go 调用 C 代码的边界跨越成本，Go 1.26 减少约 30% | CGO、FFI、性能 |

### 3. 标准库

| 概念 | 定义 | 相关术语 |
|------|------|----------|
| HPKE | Hybrid Public Key Encryption，RFC 9180 标准混合加密 | KEM、KDF、AEAD、后量子 |
| SIMD | Single Instruction Multiple Data，单指令多数据并行处理 | 向量、AVX、SSE、NEON |
| Secure Memory | 安全内存管理，确保敏感数据被安全擦除 | 前向保密、内存清除 |
| errors.AsType | 泛型版本的 errors.As，提供类型安全的错误断言 | 类型断言、泛型、错误处理 |

### 4. 工具链

| 概念 | 定义 | 相关术语 |
|------|------|----------|
| Go Analysis Framework | 统一的代码分析框架，被 go vet 和 go fix 使用 | AST、SSA、静态分析 |
| Type Satisfaction | 类型满足约束的关系判断，递归约束的核心机制 | 类型系统、约束求解 |

## 概念关系图

```
Go 1.26 语言改进
├── 语法层
│   ├── new(expr) ──→ 简化指针创建
│   └── 递归约束 ──→ 自引用类型
├── 运行时层
│   ├── Green Tea GC ──→ 更低延迟
│   └── 栈分配优化 ──→ 更少堆分配
├── 标准库层
│   ├── crypto/hpke ──→ 后量子加密
│   ├── simd/archsimd ──→ SIMD 操作
│   ├── runtime/secret ──→ 安全内存
│   └── errors.AsType ──→ 泛型断言
└── 工具链层
    ├── go fix ──→ Modernizers
    └── //go:fix inline ──→ API 迁移
```

## 形式化定义

### new() 语义

```
new: Expression → *Type
────────────────────────────────
new(T(v)) ≡ &T(v)
new(T{...}) ≡ &T{...}
```

### 递归约束

```
Constraint[T Constraint[T]] interface { ... }
────────────────────────────────────────────
T satisfies C[T] ⟺ T 实现 C 的所有方法
                     ∧ T 满足 C[T] 的约束
```

### 类型满足性

```
satisfies(T, C) =
    methods(T) ⊇ methods(C) ∧
    ∀m ∈ methods(C), signature(T.m) <: signature(C.m)
```
