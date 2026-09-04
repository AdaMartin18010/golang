# Go 1.27 语法专题文档集

> 本目录包含对 Go 1.27 / 1.27.1 语言特性的深入分析（2026-08 发布）。
> 形式化对应文档：[Go-1.27-Release.md](../formal/Go/Go-1.27-Release.md)（L5 级，含 BNF 与类型规则）
> 演进视角：[07-Go126-to-Go127.md](../../go-knowledge-base/02-Language-Design/03-Evolution/07-Go126-to-Go127.md)

## 文档列表

| 文档 | 内容 |
| ------ | ------ |
| [go127_syntax_analysis.md](go127_syntax_analysis.md) | 语法层变化：泛型方法、结构体字面量 key 泛化、泛型推断泛化 |
| [go127_semantic_analysis.md](go127_semantic_analysis.md) | 语义层变化：runtime 默认行为、JSON 栈切换、工具链语义 |
| [go127_toolchain_analysis.md](go127_toolchain_analysis.md) | 工具链：stdversion vet、go fix modernizer、go mod tidy、go doc |

## 🎯 Go 1.27 关键新特性速览

### 泛型方法（语言）

```go
// Go 1.26：方法不能有类型参数，只能写成泛型函数
func Drain[U any, T any](s *Stack[T], convert func(T) U) []U { ... }

// Go 1.27：方法可声明自己的类型参数
func (s *Stack[T]) Drain[U any](convert func(T) U) []U { ... }
strs := s.Drain(strconv.Itoa) // U 由实参推导
```

### 结构体字面量 key 泛化

```go
type Inner struct{ X int }
type Outer struct{ Inner }

o := Outer{Inner.X: 42} // Go 1.27 起合法（此前仅读取可提升）
```

### 容器感知 GOMAXPROCS（runtime，默认开启）

```go
// 容器中：1.26 默认 GOMAXPROCS = 宿主机核数
//         1.27 默认 GOMAXPROCS = min(宿主机核数, cgroup 配额)
// 回退：GODEBUG=containermaxprocs=0
```

### json/v2（标准库 GA）

```go
import json "encoding/json/v2"

err := json.Unmarshal(data, &v, json.RejectUnknownMembers(true))     // 语义 Options
err := json.Unmarshal(data, &v, jsontext.AllowDuplicateNames(true))  // 语法 Options
```

## 版本基线

- **当前版本**: Go 1.27.1（2026-08）
- **分析日期**: 2026-09-04
- **验证方式**: 本机 go1.27.1 `go doc` 与 GOROOT 源码交叉验证 + go.dev/doc/go1.27
