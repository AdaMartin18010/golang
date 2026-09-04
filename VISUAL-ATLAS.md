# VISUAL ATLAS — 形式化与可视化图谱索引

> **用途**：全仓形式化分析、可视化与语法语义专题文档的统一入口（由 view/ 实际内容索引生成，2026-09-04）。
> **权威来源**：概念解释以 `go-knowledge-base/` 五维权威页为准，本页仅作导航。

---

## 1. Go 形式化语义（view/formal/Go/）

| 专题 | 文档 |
| ------ | ------ |
| 总览 | [00-Overview](view/formal/Go/00-Overview.md) |
| 语法 | [EBNF Grammar](view/formal/Go/01-Syntax/EBNF-Grammar.md) |
| 静态语义 | [FG Calculus](view/formal/Go/02-Static-Semantics/FG-Calculus.md) |
| 动态语义 | [Small-Step Semantics](view/formal/Go/03-Dynamic-Semantics/Small-Step-Semantics.md) |
| 运行时 | [GMP Scheduler](view/formal/Go/04-Runtime-System/GMP-Scheduler.md) · [Channel Implementation](view/formal/Go/04-Runtime-System/Channel-Implementation.md) |
| 泛型扩展 | [FGG Calculus](view/formal/Go/05-Extension-Generics/FGG-Calculus.md) |
| 验证 | [Type Safety Proof](view/formal/Go/06-Verification/Type-Safety-Proof.md) |

## 2. Go 版本形式化分析（view/formal/）

| 版本 | 文档 |
| ------ | ------ |
| Go 1.25 | [Spec Changes](view/formal/Go-1.25-Spec-Changes.md) |
| Go 1.26.1 | [Comprehensive](view/formal/Go/Go-1.26.1-Comprehensive.md) · [Spec Changes](view/formal/Go/Go-1.26.1-Spec-Changes.md) · [Feature Interactions](view/formal/Go/Go-1.26.1-Feature-Interactions.md) |
| Go 1.27 | [Release](view/formal/Go/Go-1.27-Release.md) · [Preview](view/formal/Go/Go-1.27-Preview.md) |
| 泛型方法 | [Go Generic Methods](view/formal/Go-Generic-Methods.md) · [Type Inference](view/formal/Go-Generics-Type-Inference.md) |
| 内存模型 | [Formalization](view/formal/Go-Memory-Model-Formalization.md) · [Complete](view/formal/Go/Go-Memory-Model-Complete-Formalization.md) |
| 类型系统 | [Go Type System](view/formal/Go-Type-System.md) |
| 其他 | [Operational Semantics](view/formal/Go-Operational-Semantics.md) · [Runtime Execution Tree](view/formal/Go-Runtime-Execution-Tree.md) · [Syntax-Semantics Complete](view/formal/Go-Syntax-Semantics-Complete.md) |

## 3. 跨语言对比与验证工具

| 主题 | 文档 |
| ------ | ------ |
| CSP 精化 | [CSP Refinement (FDR)](view/formal/CSP-Refinement-FDR.md) |
| 三语言对比 | [Three Lines Comparison](view/formal/Three-Lines-Comparison.md) |
| WebAssembly | [Go WebAssembly Integration](view/formal/Go/Go-WebAssembly-Integration.md) |

## 4. 语法与反射专题（view/）

| 目录 | 内容 |
| ------ | ------ |
| [Go1.27语法/](view/Go1.27语法/README.md) | 1.27 语法/语义/工具链三篇分析 |
| [Go-1.26.1-Reflection-Panorama/](view/Go-1.26.1-Reflection-Panorama/) | 1.26.1 反射全景（7 篇） |
| [go1261_complete_analysis/](view/go1261_complete_analysis/) | 1.26.1 完整分析 |
