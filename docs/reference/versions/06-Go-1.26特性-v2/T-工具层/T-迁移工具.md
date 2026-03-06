# Go 1.26 迁移工具指南

> **文档层级**: T-工具层 (Tools Layer)
> **文档类型**: 迁移指南 (Migration Guide)
> **最后更新**: 2026-03-06

---

## 一、自动化迁移工具

### 1.1 go fix

```bash
# 基本用法
go fix ./...

# 查看变更（不实际修改）
go fix -diff ./...

# 特定规则
go fix -r=newexpr ./...    # 仅new表达式规则
go fix -r=errorsas ./...   # 仅errors.AsType规则
```

### 1.2 支持的转换

| 规则 | 描述 | 示例 |
|------|------|------|
| `newexpr` | 转换为new(expr) | `&[]int{42}[0]` → `new(42)` |
| `errorsas` | 使用errors.AsType | `errors.As(err, &e)` → `errors.AsType[Error](err)` |

---

## 二、手动迁移步骤

### 2.1 项目升级流程

```bash
# 1. 更新Go版本
$ go version
go version go1.26

# 2. 更新go.mod
echo 'go 1.26' > go.mod
go mod tidy

# 3. 运行自动化修复
go fix ./...

# 4. 运行测试
go test ./...

# 5. 构建验证
go build ./...
```

### 2.2 常见迁移模式

```go
// 迁移前
type Config struct {
    Timeout *int
}
cfg := Config{
    Timeout: &[]int{30}[0],
}

// 迁移后
cfg := Config{
    Timeout: new(30),
}
```

---

## 三、兼容性处理

### 3.1 版本条件编译

```go
//go:build go1.26

package mypackage

import "errors"

func handleError(err error) (*MyError, bool) {
    return errors.AsType[*MyError](err)
}
```

### 3.2 回退方案

```go
// 兼容旧版本的封装
func asMyError(err error) (*MyError, bool) {
    var myErr *MyError
    if errors.As(err, &myErr) {
        return myErr, true
    }
    return nil, false
}
```

---

## 四、验证清单

```markdown
## 迁移验证
- [ ] 代码编译通过
- [ ] 所有测试通过
- [ ] 基准测试无性能回退
- [ ] 生产环境配置更新
```

---

**相关工具**: `go fix`, `go mod`, `go test`
