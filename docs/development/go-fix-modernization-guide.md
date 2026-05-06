# Go 1.26 go fix 现代化实战手册

> **文档编号**: DEV-2026-GOFIX
> **适用版本**: Go 1.26+
> **最后更新**: 2026-05-06

---

## 一、go fix 重写了什么

Go 1.26 完全重写了 `go fix`，使其成为 Go 代码的**现代化工具**（modernizer）。

```bash
# 查看可用的 modernization 规则
go fix -list

# 常见输出:
#   modernize:
#     - slices.Sort 替换手搓排序
#     - maps.Keys/Values 替换手搓遍历
#     - min/max 内置函数替换条件判断
#     - range over integers
```

---

## 二、一键现代化工作流

### 2.1 本地开发

```bash
# 1. 执行现代化
go fix ./...

# 2. 查看变更
git diff

# 3. 运行测试确认无回归
go test ./...

# 4. 提交
git add -A && git commit -m "chore: apply go fix modernizations"
```

### 2.2 CI/CD 集成

```yaml
# .github/workflows/go-fix.yml（已集成到项目）
name: Go Fix Check
on: [push, pull_request]
jobs:
  go-fix:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26.2'
      - run: go fix ./...
      - run: git diff --exit-code
```

---

## 三、典型现代化场景

### 3.1 切片排序

```go
// 之前
sort.Slice(users, func(i, j int) bool {
    return users[i].Age < users[j].Age
})

// go fix 后
slices.SortFunc(users, func(a, b User) int {
    return cmp.Compare(a.Age, b.Age)
})
```

### 3.2 Map 键提取

```go
// 之前
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}

// go fix 后
keys := maps.Keys(m)
```

### 3.3 最小/最大值

```go
// 之前
if a < b {
    min = a
} else {
    min = b
}

// go fix 后
min = min(a, b) // Go 1.21 内置函数
```

---

## 四、注意事项

1. **go fix 会修改源代码**，提交前务必 review diff
2. **部分现代化可能改变语义**（如 `any` 替换 `interface{}`），需人工确认
3. **配合 `go vet`** 使用，捕获 fix 引入的潜在问题
4. **大型仓库建议分模块执行**，避免一次性变更过大

---

## 五、与 go fmt 的关系

```text
go fmt    → 仅格式化代码（不改变语义）
go vet    → 静态分析，报告可疑代码
go fix    → 自动应用现代化重写（改变源代码）

推荐顺序:
  go fix ./... → go vet ./... → go test ./... → gofmt -w .
```
