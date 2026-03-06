# Go 1.26 知识图谱

```mermaid
graph TB
    subgraph "语言特性"
        A[new 表达式] --> A1[简化指针创建]
        A --> A2[可选字段处理]
        B[递归泛型约束] --> B1[自引用类型]
        B --> B2[树/图抽象]
        B --> B3[通用算法]
    end

    subgraph "运行时"
        C[Green Tea GC] --> C1[更低延迟]
        C --> C2[更高吞吐]
        C --> C3[自动启用]
        D[栈分配优化] --> D1[小切片优化]
        D --> D2[逃逸分析改进]
        E[cgo 优化] --> E1[30% 开销减少]
    end

    subgraph "标准库"
        F[crypto/hpke] --> F1[后量子加密]
        F --> F2[混合 KEM]
        G[simd/archsimd] --> G1[向量操作]
        G --> G2[图像处理]
        H[runtime/secret] --> H1[安全擦除]
        H --> H2[前向保密]
        I[errors.AsType] --> I1[泛型断言]
    end

    subgraph "工具链"
        J[go fix] --> J1[Modernizers]
        J --> J2[自动内联]
        K[分析框架] --> K1[go vet]
        K --> K2[自定义分析器]
    end

    subgraph "应用场景"
        A2 --> L[API 请求构建]
        B2 --> M[二叉搜索树]
        B2 --> N[图遍历]
        F1 --> O[安全通信]
        G1 --> P[多媒体处理]
        H1 --> Q[密码学操作]
        J1 --> R[代码现代化]
    end

    style A fill:#4CAF50
    style B fill:#4CAF50
    style C fill:#2196F3
    style F fill:#FF9800
    style J fill:#9C27B0
```

## 学习路径

```
初学者路径:
├── 1. 了解 new(expr) 语法
├── 2. 运行 go fix 现代化代码
└── 3. 使用 errors.AsType

进阶路径:
├── 1. 掌握递归泛型约束
├── 2. 应用 Green Tea GC 调优
├── 3. 使用 crypto/hpke
└── 4. 开发自定义 Modernizer

专家路径:
├── 1. 深入 SIMD 优化
├── 2. 使用 runtime/secret
├── 3. 理解 GC 内部机制
└── 4. 贡献 Go 语言
```

## 特性依赖关系

```
Go 1.26 特性依赖:

new(expr)
  └── 无依赖 (独立特性)

递归泛型约束
  └── 泛型基础 (Go 1.18+)

Green Tea GC
  └── 无依赖 (自动启用)

crypto/hpke
  └── crypto 包基础

simd/archsimd
  └── GOEXPERIMENT=simd
  └── amd64 架构

runtime/secret
  └── GOEXPERIMENT=runtimesecret
  └── Linux amd64/arm64

go fix Modernizers
  └── Go 分析框架
```
