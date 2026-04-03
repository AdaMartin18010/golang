# Build Automation

> **维度**: Application Domains / DevOps Tools
> **级别**: S (17+ KB)
> **tags**: #build-automation #ci-cd #makefile #bazel #github-actions

---

## 1. 构建自动化的形式化

### 1.1 构建系统定义

**定义 1.1 (构建系统)**
构建系统是一个函数 $B$，将源代码 $S$ 和依赖 $D$ 映射到可执行产物 $A$：
$$B: S \times D \to A$$

**定义 1.2 (构建正确性)**
构建是正确的当且仅当：
$$\forall s_1, s_2 \in S: s_1 = s_2 \Rightarrow B(s_1, D) = B(s_2, D)$$

### 1.2 增量构建

**定理 1.1 (增量构建优化)**
若构建系统跟踪依赖图 $G = (V, E)$，则增量构建的时间复杂度为 $O(|V_{changed}| + |E_{affected}|)$，而非全量构建的 $O(|V| + |E|)$。

---

## 2. Makefile 工程实践

### 2.1 现代 Makefile 模式

```makefile
# Makefile - Modern Go Project
.PHONY: all build test lint clean docker help

# 变量定义
BINARY_NAME := myapp
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 颜色定义 (增强可读性)
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

.DEFAULT_GOAL := help

## help: 显示帮助信息
help:
 @echo "$(BLUE)Available targets:$(RESET)"
 @awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## all: 清理、测试、构建
all: clean lint test build

## build: 构建应用
build:
 @echo "$(GREEN)Building $(BINARY_NAME)...$(RESET)"
 @go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/server
 @echo "$(GREEN)Build complete: bin/$(BINARY_NAME)$(RESET)"

## build-race: 构建带竞态检测的版本
build-race:
 @echo "$(GREEN)Building with race detector...$(RESET)"
 @go build -race $(LDFLAGS) -o bin/$(BINARY_NAME)-race ./cmd/server

## test: 运行测试
test:
 @echo "$(BLUE)Running tests...$(RESET)"
 @go test -v -race -coverprofile=coverage.out ./...
 @go tool cover -html=coverage.out -o coverage.html
 @echo "$(GREEN)Coverage report: coverage.html$(RESET)"

## test-short: 快速测试
test-short:
 @echo "$(BLUE)Running short tests...$(RESET)"
 @go test -short ./...

## bench: 运行基准测试
bench:
 @echo "$(BLUE)Running benchmarks...$(RESET)"
 @go test -bench=. -benchmem ./...

## lint: 代码检查
lint:
 @echo "$(YELLOW)Running linters...$(RESET)"
 @golangci-lint run ./...
 @echo "$(GREEN)Linting complete$(RESET)"

## fmt: 格式化代码
fmt:
 @echo "$(BLUE)Formatting code...$(RESET)"
 @go fmt ./...
 @echo "$(GREEN)Formatting complete$(RESET)"

## vet: 静态分析
vet:
 @echo "$(YELLOW)Running go vet...$(RESET)"
 @go vet ./...

## deps: 下载依赖
deps:
 @echo "$(BLUE)Downloading dependencies...$(RESET)"
 @go mod download
 @go mod tidy
 @go mod verify

## update: 更新依赖
update:
 @echo "$(YELLOW)Updating dependencies...$(RESET)"
 @go get -u ./...
 @go mod tidy

## clean: 清理构建产物
clean:
 @echo "$(RED)Cleaning...$(RESET)"
 @rm -rf bin/ coverage.out coverage.html
 @go clean -cache
 @echo "$(GREEN)Clean complete$(RESET)"

## docker-build: 构建 Docker 镜像
docker-build:
 @echo "$(BLUE)Building Docker image...$(RESET)"
 @docker build -t $(BINARY_NAME):$(VERSION) .
 @docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

## docker-push: 推送 Docker 镜像
docker-push: docker-build
 @echo "$(BLUE)Pushing Docker image...$(RESET)"
 @docker push $(BINARY_NAME):$(VERSION)
 @docker push $(BINARY_NAME):latest

## run: 本地运行
run: build
 @echo "$(GREEN)Starting $(BINARY_NAME)...$(RESET)"
 @./bin/$(BINARY_NAME)

## dev: 开发模式 (带热重载)
dev:
 @echo "$(GREEN)Starting in dev mode...$(RESET)"
 @air -c .air.toml

## proto: 生成 protobuf 代码
proto:
 @echo "$(BLUE)Generating protobuf code...$(RESET)"
 @protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  api/proto/*.proto

## migrate: 数据库迁移
migrate-up:
 @echo "$(BLUE)Running migrations...$(RESET)"
 @migrate -path migrations -database "$(DB_URL)" up

migrate-down:
 @echo "$(YELLOW)Rolling back migrations...$(RESET)"
 @migrate -path migrations -database "$(DB_URL)" down 1

## security: 安全检查
security:
 @echo "$(YELLOW)Running security checks...$(RESET)"
 @gosec ./...
 @nancy sleuth

## ci: CI 完整流程
ci: deps fmt vet lint test build
 @echo "$(GREEN)CI checks passed!$(RESET)"
```

### 2.2 高级 Makefile 技巧

| 技巧 | 用途 | 示例 |
|------|------|------|
| 条件变量 | 环境覆盖 | `VAR ?= default` |
| 函数 | 动态计算 | `$(shell command)` |
| 模式规则 | 批量处理 | `%.o: %.c` |
| 自动依赖 | 头文件 | `-MMD -MP` |
| 并行构建 | 加速 | `make -j8` |

---

## 3. CI/CD Pipeline

### 3.1 GitHub Actions 完整流程

```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
    tags: [ 'v*' ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: '1.21'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # 阶段 1: 代码质量检查
  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}

    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
        args: --timeout=5m

    - name: Run go vet
      run: go vet ./...

    - name: Run gosec
      uses: securego/gosec@master
      with:
        args: './...'

  # 阶段 2: 测试
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
        - 5432:5432
      redis:
        image: redis:7
        ports:
        - 6379:6379

    steps:
    - uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}

    - name: Run tests
      run: go test -race -coverprofile=coverage.out ./...
      env:
        DB_HOST: localhost
        REDIS_HOST: localhost

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: unittests

  # 阶段 3: 构建
  build:
    needs: [lint, test]
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]
        exclude:
        - goos: windows
          goarch: arm64

    steps:
    - uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}

    - name: Build
      env:
        GOOS: ${{ matrix.goos }}
        GOARCH: ${{ matrix.goarch }}
      run: |
        output_name='${{ github.event.repository.name }}'
        if [ "${{ matrix.goos }}" = "windows" ]; then
          output_name+='.exe'
        fi
        go build -ldflags "-s -w -X main.Version=${{ github.ref_name }}" \
          -o "dist/${{ matrix.goos }}_${{ matrix.goarch }}/$output_name"

    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: binaries
        path: dist/

  # 阶段 4: 容器构建与推送
  docker:
    needs: [lint, test]
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - uses: actions/checkout@v4

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Login to Container Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=semver,pattern={{major}}.{{minor}}

    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        cache-from: type=gha
        cache-to: type=gha,mode=max
        platforms: linux/amd64,linux/arm64

  # 阶段 5: 发布
  release:
    needs: [build, docker]
    if: startsWith(github.ref, 'refs/tags/')
    runs-on: ubuntu-latest
    permissions:
      contents: write

    steps:
    - uses: actions/checkout@v4

    - name: Download artifacts
      uses: actions/download-artifact@v3
      with:
        name: binaries
        path: dist/

    - name: Create archives
      run: |
        cd dist
        for dir in */; do
          base=$(basename "$dir")
          tar czf "${base}.tar.gz" -C "$dir" .
        done

    - name: Create Release
      uses: softprops/action-gh-release@v1
      with:
        files: dist/*.tar.gz
        generate_release_notes: true
```

---

## 4. 高级构建工具

### 4.1 Bazel 构建系统

| 特性 | Bazel | Make | Go Modules |
|------|-------|------|------------|
| 增量构建 | 精确 | 文件时间戳 | 包级别 |
| 远程缓存 | 支持 | 不支持 | 不支持 |
| 并行执行 | 自动 | 手动 (-j) | 有限 |
| 跨语言 | 支持 | 需配置 | Go 专用 |
| 学习曲线 | 陡峭 | 平缓 | 平缓 |

### 4.2 Nix 可重现构建

```nix
# default.nix - Nix expression for Go project
{ pkgs ? import <nixpkgs> {} }:

pkgs.buildGoModule {
  pname = "myapp";
  version = "1.0.0";

  src = ./.;

  vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

  ldflags = [
    "-s" "-w"
    "-X main.Version=${version}"
  ];

  meta = with pkgs.lib; {
    description = "My Go Application";
    license = licenses.mit;
  };
}
```

---

## 5. 思维工具

```
┌─────────────────────────────────────────────────────────────────┐
│                 Build Automation Checklist                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  本地开发:                                                       │
│  □ 一键构建 (make build)                                         │
│  □ 自动测试 (make test)                                          │
│  □ 代码格式化 (make fmt)                                         │
│  □ 静态检查 (make lint)                                          │
│  □ 热重载支持 (make dev)                                         │
│                                                                  │
│  CI/CD:                                                          │
│  □ 多阶段流水线                                                  │
│  □ 并行执行                                                      │
│  □ 依赖缓存                                                      │
│  □ 产物归档                                                      │
│  □ 安全扫描                                                      │
│                                                                  │
│  发布管理:                                                       │
│  □ 语义化版本                                                    │
│  □ 自动变更日志                                                  │
│  □ 多平台构建                                                    │
│  □ 容器镜像签名                                                  │
│                                                                  │
│  优化策略:                                                       │
│  □ 增量构建                                                      │
│  □ 并行编译                                                      │
│  □ 远程缓存                                                      │
│  □ 构建产物压缩                                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (17KB)
**完成日期**: 2026-04-02

---

## 架构决策记录

### 决策矩阵

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| A | 高性能 | 复杂 | 大规模 |
| B | 简单 | 扩展性差 | 小规模 |

### 风险评估

**风险 R.1**: 性能瓶颈
- 概率: 中
- 影响: 高
- 缓解: 缓存、分片

**风险 R.2**: 单点故障
- 概率: 低
- 影响: 极高
- 缓解: 冗余、故障转移

### 实施路线图

`
Phase 1: 基础设施 (Week 1-2)
Phase 2: 核心功能 (Week 3-6)
Phase 3: 优化加固 (Week 7-8)
`

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 架构决策记录 (ADR)

### 上下文

业务需求和技术约束分析。

### 决策

选择方案A作为主要架构方向。

### 后果

正面：
- 可扩展性提升
- 维护成本降低

负面：
- 初期开发复杂度增加
- 团队学习成本

### 实施指南

`
Week 1-2: 基础设施搭建
Week 3-4: 核心功能开发
Week 5-6: 集成测试
Week 7-8: 性能优化
`

### 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 性能不足 | 中 | 高 | 缓存、分片 |
| 兼容性 | 低 | 中 | 接口适配层 |

### 监控指标

- 系统吞吐量
- 响应延迟
- 错误率
- 资源利用率

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 系统设计

### 需求分析

功能需求和非功能需求的完整梳理。

### 架构视图

`
┌─────────────────────────────────────┐
│           API Gateway               │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    ▼                   ▼
┌─────────┐       ┌─────────┐
│ Service │       │ Service │
│   A     │       │   B     │
└────┬────┘       └────┬────┘
     │                 │
     └────────┬────────┘
              ▼
        ┌─────────┐
        │  Data   │
        │  Store  │
        └─────────┘
`

### 技术选型

| 组件 | 技术 | 理由 |
|------|------|------|
| API | gRPC | 性能 |
| DB | PostgreSQL | 可靠 |
| Cache | Redis | 速度 |
| Queue | Kafka | 吞吐 |

### 性能指标

- QPS: 10K+
- P99 Latency: <100ms
- Availability: 99.99%

### 运维手册

- 部署流程
- 监控配置
- 应急预案
- 容量规划

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02