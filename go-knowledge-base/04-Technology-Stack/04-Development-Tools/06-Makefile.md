# Makefile

> **分类**: 开源技术堆栈

---

## 基本结构

```makefile
.PHONY: build test clean

build:
    go build -o bin/app .

test:
    go test -v ./...

clean:
    rm -rf bin/
```

---

## 变量

```makefile
BINARY=app
VERSION=1.0.0
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

build:
    go build $(LDFLAGS) -o bin/$(BINARY) .
```

---

## 常用命令

```makefile
.DEFAULT_GOAL := help

help: ## 显示帮助
    @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dev: ## 开发模式
    air

lint: ## 代码检查
    golangci-lint run

fmt: ## 格式化
    go fmt ./...
```
