# Go 1.27 json/v2 迁移指南

**Go版本**: Go 1.27 / 1.27.1
**日期**: 2026-08
**权威来源**: [go.dev/doc/go1.27](https://go.dev/doc/go1.27)、`go doc encoding/json/v2`（本机 go1.27.1 验证）
**姊妹篇**: [Go-1.27-Release.md](../../view/formal/Go/Go-1.27-Release.md) §4

---

## 0. TL;DR

Go 1.27 中 `encoding/json/v2` 正式 GA，**且 v1 `encoding/json` 已改由 v2 实现**。多数 v1 代码无需改动即可编译运行，但以下三类必须回归：

1. 对错误消息文本做字符串匹配的测试；
2. 依赖重复 key 宽松处理、quoted null `,string` 等边角行为的代码（部分回归已在 1.27.1 修复）；
3. 使用 v1 `Decoder.Token` 流式处理的代码（1.27.1 修复 `ErrUnexpectedEOF` 缺失）。

紧急回退：`GOEXPERIMENT=nojsonv2`（后续版本将移除，仅作过渡）。

---

## 1. 包结构

| 包 | 角色 | 关键类型 |
|----|------|---------|
| `encoding/json`（v1） | 兼容门面，v2 实现 | 不变（`Marshal`、`Decoder`...） |
| `encoding/json/v2` | 语义层：Go 值 ↔ JSON 数据 | `Marshal/Unmarshal`、`Options`、`Marshaler` |
| `encoding/json/jsontext` | 语法层：JSON 文本流 | `Encoder/Decoder`、`Token`、`Value`、`Options` |

v2 函数签名接受两类 Options：**语义 Options**（v2 包，`RejectUnknownMembers`、`Deterministic` 等）与**语法 Options**（jsontext 包，`AllowDuplicateNames`、`Indent` 等）。

## 2. v1 → v2 用法对照

### 2.1 Marshal / Unmarshal

```go
import json "encoding/json/v2"

data, err := json.Marshal(v)                    // 同 v1
err := json.Unmarshal(data, &v)                 // 同 v1

// v2 特有：Writer/Reader 与 jsontext 集成
err := json.MarshalWrite(w, v)                  // 写入 io.Writer
err := json.UnmarshalRead(r, &v)                // 从 io.Reader
err := json.MarshalEncode(jsontext.NewEncoder(w), v)
```

### 2.2 Options（v2 核心增量）

```go
// 语义 Options（encoding/json/v2）
json.Unmarshal(data, &v, json.RejectUnknownMembers(true))  // 严格模式
json.Marshal(v, json.Deterministic(true))                  // 稳定输出（map 排序）

// 语法 Options（encoding/json/jsontext）
json.Unmarshal(data, &v, jsontext.AllowDuplicateNames(true)) // 放宽重复 key（v1 默认宽松）
json.Marshal(v, jsontext.WithIndent("  "))                   // 缩进
```

### 2.3 自定义序列化

| v1 | v2 |
|----|----|
| `json.Marshaler` / `Unmarshaler` | 兼容支持 + 新增 `MarshalerTo` / `UnmarshalerFrom`（jsontext 流式） |
| 仅方法定制 | 新增 `MarshalFunc[T]` / `UnmarshalFunc[T]` **函数级定制**（无需改类型） |

## 3. 行为差异清单

| 行为 | v1 | v2（默认） | 处置 |
|------|----|-----------|------|
| 对象重复 key | 静默取最后值 | **拒绝**（错误） | `jsontext.AllowDuplicateNames(true)` 恢复宽松 |
| 未知字段 | 忽略 | 忽略（`RejectUnknownMembers(true)` 可拒绝） | 严格模式可选开启 |
| 大小写匹配 | 优先精确，后 fold | 精确匹配 | fold 行为变化，相关测试需回归 |
| 错误消息文本 | v1 文案 | 可能不同 | **禁止对错误文本做字符串断言** |
| `,string` 选项 | 宽松 | 曾拒绝 quoted null（**1.27.1 已修复** #81083） | 升级 1.27.1 |
| `Decoder.Token` | 流式 | v1 门面由 v2 实现，曾缺失 `io.ErrUnexpectedEOF`（**1.27.1 已修复** #81012） | 升级 1.27.1 |
| 浮点格式化 | Ryū（ shortest ） | 语义等价 | 边界值测试回归 |
| HTML 转义（`<`/`>`） | 默认转义 | 默认行为对齐 v1 | `jsontext.EscapeHTML` 控制 |

## 4. 迁移步骤

1. **升级工具链到 1.27.1**（必须 ≥1.27.1，避开 #81083/#81012 回归）。
2. 全量回归 JSON 相关测试；把错误文本断言改为 `errors.Is` / 结构化判断。
3. 按需引入 v2：
   - 仅需稳定输出 → `json.Deterministic(true)`；
   - 需要严格输入校验 → `json.RejectUnknownMembers(true)`；
   - 需要流式/高性能 → `jsontext` 直编解码。
4. 新代码直接 `import json "encoding/json/v2"`；公共库保持 v1 门面以兼容下游。
5. 灰度验证后再移除 `GOEXPERIMENT=nojsonv2` 依赖（不要长期依赖回退开关）。

## 5. 最小可运行示例

见 [examples/go127-features](../../examples/go127-features) 的 `demonstrateJSONv2`（main.go）与 `TestJSONv2Options`（synctest_test.go），覆盖 Marshal/Unmarshal、重复 key 拒绝与放宽、RejectUnknownMembers。

## 6. 参考

- [encoding/json/v2 包文档](https://pkg.go.dev/encoding/json/v2)（`go doc encoding/json/v2`）
- Go 1.27.1 milestone：#81083、#81012
- [Go-1.27-Release.md](../../view/formal/Go/Go-1.27-Release.md) §4.3
