# Go Modules

> **分类**: 开源技术堆栈

---

## 初始化

```bash
go mod init example.com/project
go mod tidy
```

---

## go.mod

```go
module example.com/project

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4
)

require (
    github.com/davecgh/go-spew v1.1.1 // indirect
    github.com/pmezard/go-difflib v1.0.0 // indirect
)
```

---

## 常用命令

```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 查看依赖
go list -m all

# 更新依赖
go get -u ./...
go get -u github.com/gin-gonic/gin

# 清理未使用
go mod tidy

# 验证
go mod verify
```

---

## 版本管理

### 语义化版本

```go
// v1.2.3
// 主版本.次版本.补丁版本

go get github.com/pkg/errors@v0.9.1
go get github.com/pkg/errors@latest
go get github.com/pkg/errors@master
```

### 兼容版本

```go
// v0: 不保证兼容
go get github.com/pkg/errors@v0.9.1

// v1+: 保证向后兼容
go get github.com/gin-gonic/gin@v1.9.0
```

---

## 私有仓库

```bash
# 配置 GOPRIVATE
go env -w GOPRIVATE="*.example.com"

# 或
go env -w GOPRIVATE="github.com/mycompany/*"
```

---

## 工作区 (Go 1.18+)

```bash
# 创建工作区
go work init ./project1 ./project2

# go.work
go 1.22

use (
    ./project1
    ./project2
)

replace example.com/project1 => ./project1
```
