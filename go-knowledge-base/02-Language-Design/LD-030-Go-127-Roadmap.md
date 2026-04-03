# LD-030-Go-127-Roadmap

> **Dimension**: 02-Language-Design
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: Go 1.27 Roadmap (Expected: August 2026)
> **Size**: >20KB

---

## 1. Go 1.27 概览

### 1.1 预期发布

- **预期日期**: 2026年8月
- **开发周期**: 6个月
- **主题**: 泛型方法、json/v2 GA、GC稳定化

### 1.2 路线图概览

| 特性 | 状态 | 优先级 |
|------|------|--------|
| 泛型方法 | 已接受 | 高 |
| Green Tea GC | 唯一GC | 高 |
| json/v2 | GA准备 | 高 |
| Goroutine泄漏检测 | 默认启用 | 中 |
| GODEBUG清理 | 移除旧设置 | 中 |

---

## 2. 泛型方法 (Generic Methods)

### 2.1 提案概述

**作者**: Robert Griesemer
**状态**: 已接受 (2025年12月)
**目标**: Go 1.27

**核心概念**:
允许具体类型拥有自己的泛型方法，这些方法的类型参数独立于接收者的类型参数。

### 2.2 语法与示例

```go
// 容器类型带泛型方法
type Container[T any] struct {
    items []T
}

// 泛型方法: 转换为其他类型
func (c Container[T]) MapTo[U any](fn func(T) U) Container[U] {
    result := make([]U, len(c.items))
    for i, item := range c.items {
        result[i] = fn(item)
    }
    return Container[U]{items: result}
}

// 使用
intContainer := Container[int]{items: []int{1, 2, 3}}
strContainer := intContainer.MapTo[string](func(i int) string {
    return strconv.Itoa(i)
})
// strContainer = Container[string]{"1", "2", "3"}
```

### 2.3 数据库客户端示例

```go
// 泛型数据库客户端
type DB struct {
    conn *sql.DB
}

// 泛型查询方法
func (db *DB) QueryOne[T any](query string, args ...any) (*T, error) {
    row := db.conn.QueryRow(query, args...)

    var result T
    // 使用反射或代码生成填充result
    if err := row.Scan(&result); err != nil {
        return nil, err
    }
    return &result, nil
}

// 使用
type User struct {
    ID   int
    Name string
}

user, err := db.QueryOne[User]("SELECT * FROM users WHERE id = ?", 1)
```

### 2.4 类型安全验证

```go
// 编译器确保类型安全
numbers := Container[int]{items: []int{1, 2, 3}}

// 编译错误: 类型不匹配
// strings := numbers.MapTo[string](strconv.Atoi)  // 错误签名

// 正确用法
strings := numbers.MapTo[string](func(i int) string {
    return fmt.Sprintf("num-%d", i)
})
```

### 2.5 与接口的交互

```go
// 泛型方法在接口中的限制
type Mapper[T any] interface {
    // 接口不能有泛型方法
    // MapTo[U any](func(T) U) Container[U]  // 非法
}

// 但可以在具体类型上定义
type IntMapper struct {
    value int
}

func (m IntMapper) MapTo[U any](fn func(int) U) U {
    return fn(m.value)
}
```

---

## 3. encoding/json/v2

### 3.1 当前状态

**Go 1.26**: 实验性 (`GOEXPERIMENT=jsonv2`)
**Go 1.27目标**: GA (一般可用)

### 3.2 新包结构

```
encoding/
├── json/           # v1 (保持兼容)
│   └── ...
├── json/v2/        # v2 API
│   ├── marshal.go
│   └── unmarshal.go
└── json/jsontext/  # 底层JSON处理
    ├── encode.go
    └── decode.go
```

### 3.3 主要改进

**流式处理**:

```go
import "encoding/json/jsontext"

// 流式编码
enc := jsontext.NewEncoder(writer)
enc.WriteToken(jsontext.ObjectStart)
enc.WriteToken(jsontext.String("name"))
enc.WriteToken(jsontext.String("Alice"))
enc.WriteToken(jsontext.ObjectEnd)

// 流式解码
dec := jsontext.NewDecoder(reader)
for dec.ReadToken() {
    tok := dec.Token()
    // 处理token
}
```

**新接口**:

```go
// 更高效的自定义序列化
type MarshalJSONTo interface {
    MarshalJSONTo(enc *jsontext.Encoder, opts jsonv2.Options) error
}

type UnmarshalJSONFrom interface {
    UnmarshalJSONFrom(dec *jsontext.Decoder, opts jsonv2.Options) error
}
```

**更好的错误信息**:

```go
import "encoding/json/v2"

var data struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// 错误包含JSON Pointer位置
err := jsonv2.Unmarshal(badJSON, &data)
// 错误: "at /user/age: cannot unmarshal string into Go struct field User.age of type int"
```

**严格模式**:

```go
// 默认更严格
// - 拒绝无效UTF-8
// - 拒绝重复字段名
// - 拒绝未知字段 (可选)

opts := jsonv2.Options{
    RejectUnknownFields: true,
    AllowDuplicateNames: false,
}

jsonv2.Unmarshal(data, &v, opts)
```

### 3.4 迁移路径

```go
// 渐进式迁移
import jsonv2 "encoding/json/v2"

// 使用V1兼容选项
opts := jsonv2.DefaultOptionsV1()
jsonv2.Unmarshal(data, &v, opts)  // 行为类似v1

// 逐步采用v2特性
opts := jsonv2.Options{
    // 混合v1和v2行为
}
```

---

## 4. Green Tea GC稳定化

### 4.1 移除opt-out

Go 1.27将完全移除`GOEXPERIMENT=nogreenteagc`选项。

```bash
# Go 1.26: 可选禁用
GOEXPERIMENT=nogreenteagc go run .  # 工作

# Go 1.27: 不识别此选项
GOEXPERIMENT=nogreenteagc go run .  # 警告: 未知实验
```

### 4.2 预期影响

- 所有工作负载使用Green Tea GC
- 进一步GC优化基于此架构
- 简化运行时维护

---

## 5. Goroutine泄漏检测

### 5.1 默认启用

Go 1.27可能将goroutine泄漏检测设为默认：

```go
// 当前 (Go 1.26): 实验性
// 需要显式启用

// Go 1.27: 可能默认启用
// 无需配置
```

### 5.2 使用场景

```go
// 测试时检测泄漏
func TestWithLeakDetection(t *testing.T) {
    // 自动检测未结束的goroutine
    defer goleak.VerifyNone(t)

    // 运行代码
    RunAsyncOperations()
}

// 生产环境监控
import _ "net/http/pprof"

// 访问 /debug/pprof/goroutine?debug=1
// 查看活跃goroutine
```

---

## 6. 标准库增强

### 6.1 预期改进

| 包 | 预期改进 |
|----|---------|
| sync | 新的并发原语 |
| context | 性能优化 |
| net/http | HTTP/3支持推进 |
| database/sql | 连接池优化 |
| crypto | 更多后量子算法 |

### 6.2 可能的API

```go
// sync: 新原语 (推测)
type WaitGroupMap[K comparable] struct {
    // 带key的WaitGroup
}

func (wg *WaitGroupMap[K]) Add(key K)
func (wg *WaitGroupMap[K]) Done(key K)
func (wg *WaitGroupMap[K]) Wait(key K)

// context: 取消原因
ctx, cancel := context.WithCause(parent)
cancel(context.Canceled)  // 带原因的取消

// 检查原因
if err := ctx.Err(); err != nil {
    cause := context.Cause(ctx)
    // 获取原始取消原因
}
```

---

## 7. 工具链改进

### 7.1 go mod增强

```bash
# 可能的新功能

# 工作区 vendoring
go work vendor

# 依赖图可视化
go mod graph --visual

# 漏洞扫描集成
go mod audit
```

### 7.2 编译器优化

```go
// 逃逸分析改进
// 更多堆分配到栈

// 内联优化
// 更大的函数可能被内联

// 边界检查消除
// 更智能的数组访问优化
```

---

## 8. 平台支持

### 8.1 新平台

| 平台 | 状态 |
|------|------|
| WASI 0.3 | 可能支持 |
| 新ARM版本 | 持续跟进 |

### 8.2 移除支持

| 平台 | 变更 |
|------|------|
| macOS 12 | Go 1.27需要macOS 13+ |
| PowerPC ELFv1 | 切换到ELFv2 |

---

## 9. 迁移准备

### 9.1 代码准备

```go
// 1. 清理GODEBUG使用
// 检查代码中依赖的旧行为

// 2. 准备json/v2迁移
// 测试GOEXPERIMENT=jsonv2

// 3. 泛型方法准备
// 审查可能需要泛型方法的代码
```

### 9.2 测试策略

```bash
# 在CI中使用Go tip版本
go install golang.org/dl/gotip@latest
gotip download

# 运行测试
gotip test ./...
```

---

## 10. 长期展望

### 10.1 Go 1.28+ 可能方向

| 方向 | 可能性 |
|------|--------|
| 轻量级线程 | 讨论中 |
| 结构化并发 | 提案阶段 |
| 错误处理改进 | 可能重新审视 |
| 标准库模块化 | 持续进行 |

### 10.2 生态系统趋势

- WebAssembly重要性增加
- AI/ML库增长
- 云原生工具深化
- 安全优先设计

---

## 11. 参考文献

1. Generic Methods Proposal (Robert Griesemer)
2. encoding/json/v2 Design
3. Go Release Cycle
4. Go 1.27 Milestone
5. Go Developer Survey 2025

---

*Last Updated: 2026-04-03*
