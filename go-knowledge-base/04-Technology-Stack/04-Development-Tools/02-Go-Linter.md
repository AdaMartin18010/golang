# Go Linter

> **分类**: 开源技术堆栈

---

## golangci-lint

### 安装

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 使用

```bash
# 运行所有 linter
golangci-lint run

# 指定目录
golangci-lint run ./...

# 快速模式
golangci-lint run --fast
```

---

## 配置

### .golangci.yml

```yaml
run:
  timeout: 5m
  skip-dirs:
    - vendor

linters:
  enable:
    - errcheck      # 检查未处理的错误
    - gosimple      # 简化代码建议
    - govet         # 标准 vet
    - ineffassign   # 无效赋值
    - staticcheck   # 静态分析
    - unused        # 未使用代码
    - gofmt         # 格式化
    - goimports     # 导入排序
    - gocyclo       # 圈复杂度
    - misspell      # 拼写检查

linters-settings:
  gocyclo:
    min-complexity: 15

issues:
  exclude-use-default: false
```

---

## 常用 Linter

| Linter | 用途 |
|--------|------|
| errcheck | 检查未处理的错误 |
| govet | 标准静态分析 |
| staticcheck | 高级静态分析 |
| gosimple | 简化建议 |
| gosec | 安全检查 |
| deadcode | 死代码检测 |
| structcheck | 未使用的结构体字段 |

---

## CI 集成

### GitHub Actions

```yaml
name: lint
on: [push, pull_request]

jobs:
  golangci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - uses: golangci/golangci-lint-action@v3
        with:
          version: latest
```
