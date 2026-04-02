# 标准库概览 (Standard Library)

> **分类**: 开源技术堆栈

---

## 核心包

| 包名 | 用途 |
|------|------|
| `fmt` | 格式化 I/O |
| `os` | 操作系统功能 |
| `io` | 基本 I/O 接口 |
| `net/http` | HTTP 客户端/服务器 |
| `encoding/json` | JSON 编解码 |
| `database/sql` | SQL 数据库接口 |
| `sync` | 同步原语 |
| `context` | 上下文管理 |

---

## 包使用示例

### fmt

```go
fmt.Printf("Hello %s\n", "World")
fmt.Sprintf("value: %d", 42)
```

### os

```go
file, err := os.Open("file.txt")
defer file.Close()
```

### io

```go
io.Copy(dst, src)
io.ReadAll(reader)
```

---

## 设计哲学

- **小而美**: 每个包职责单一
- **接口优先**: 大量使用接口
- **显式错误**: 返回 error 而非 panic
