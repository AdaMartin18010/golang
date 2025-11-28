# CI/CD 流程

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 目录

- [1. CI/CD 架构](#1-cicd-架构)
- [2. GitHub Actions 工作流](#2-github-actions-工作流)
- [3. 构建流程](#3-构建流程)
- [4. 测试流程](#4-测试流程)
- [5. 部署流程](#5-部署流程)
- [6. 最佳实践](#6-最佳实践)

---

## 1. CI/CD 架构

### 1.1 CI/CD 流程图

```text
┌─────────────────────────────────────────────────────────────┐
│                    CI/CD 流程                                 │
└─────────────────────────────────────────────────────────────┘

开发者提交代码
    ↓
┌─────────────────┐
│  Pre-commit     │  Git Hooks（代码格式、检查）
│  Hooks          │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│  CI Pipeline    │
│  (GitHub Actions)│
│                 │
│  1. 代码质量检查 │  - gofmt, go vet, golangci-lint
│  2. 单元测试     │  - go test
│  3. 集成测试     │  - 集成测试套件
│  4. 安全扫描     │  - Gosec, Trivy
│  5. 构建镜像     │  - Docker build
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│  代码覆盖率      │  - Codecov
│  报告            │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│  CD Pipeline    │
│                 │
│  1. 推送镜像     │  - Docker Hub / Registry
│  2. 部署到环境   │  - Kubernetes / Docker
│  3. 健康检查     │  - 验证部署
│  4. 通知         │  - Slack / Email
└─────────────────┘
```

### 1.2 CI/CD 阶段

| 阶段 | 说明 | 触发条件 | 执行内容 |
|------|------|---------|---------|
| **Pre-commit** | 提交前检查 | Git commit | 代码格式、基础检查 |
| **CI** | 持续集成 | Push / PR | 代码质量、测试、构建 |
| **CD** | 持续部署 | Tag / Merge | 镜像推送、部署 |

---

## 2. GitHub Actions 工作流

### 2.1 代码质量检查工作流

```yaml
name: Code Quality

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

      - name: Check code format
        run: |
          if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
            echo "❌ 代码格式不符合规范"
            gofmt -s -d .
            exit 1
          fi

      - name: Run go vet
        run: go vet ./...
```

### 2.2 测试工作流

```yaml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go-version: ['1.21', '1.22']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Cache dependencies
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
          flags: unittests
          name: codecov-umbrella
```

### 2.3 构建和推送工作流

```yaml
name: Build and Push

on:
  push:
    tags:
      - 'v*'
    branches:
      - main
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ secrets.DOCKER_USERNAME }}/app
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ./deployments/docker/Dockerfile
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=registry,ref=${{ secrets.DOCKER_USERNAME }}/app:buildcache
          cache-to: type=registry,ref=${{ secrets.DOCKER_USERNAME }}/app:buildcache,mode=max
```

### 2.4 部署工作流

```yaml
name: Deploy

on:
  push:
    tags:
      - 'v*'
    branches:
      - main
  workflow_dispatch:
    inputs:
      environment:
        description: 'Environment to deploy'
        required: true
        default: 'staging'
        type: choice
        options:
          - staging
          - production

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: ${{ github.event.inputs.environment || 'staging' }}

    steps:
      - uses: actions/checkout@v4

      - name: Set up kubectl
        uses: azure/setup-kubectl@v3
        with:
          version: 'latest'

      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
          export KUBECONFIG=./kubeconfig

      - name: Deploy to Kubernetes
        run: |
          export KUBECONFIG=./kubeconfig
          kubectl apply -f deployments/kubernetes/
          kubectl rollout status deployment/app

      - name: Health check
        run: |
          export KUBECONFIG=./kubeconfig
          kubectl wait --for=condition=ready pod -l app=app --timeout=300s
```

---

## 3. 构建流程

### 3.1 构建阶段

| 阶段 | 说明 | 输出 |
|------|------|------|
| **依赖下载** | 下载 Go 模块依赖 | `go.mod` 缓存 |
| **代码生成** | 生成 Ent、Wire 代码 | 生成的代码文件 |
| **编译** | 编译 Go 代码 | 二进制文件 |
| **测试** | 运行测试套件 | 测试报告 |
| **构建镜像** | 构建 Docker 镜像 | Docker 镜像 |

### 3.2 多平台构建

```yaml
- name: Build for multiple platforms
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./deployments/docker/Dockerfile
    platforms: linux/amd64,linux/arm64
    push: true
    tags: ${{ steps.meta.outputs.tags }}
```

---

## 4. 测试流程

### 4.1 测试阶段

| 阶段 | 说明 | 工具 |
|------|------|------|
| **单元测试** | 测试单个函数/方法 | `go test` |
| **集成测试** | 测试组件集成 | `go test -tags=integration` |
| **E2E 测试** | 端到端测试 | 测试框架 |
| **性能测试** | 基准测试 | `go test -bench` |

### 4.2 测试矩阵

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    go-version: ['1.21', '1.22']
    test-type: [unit, integration]
```

---

## 5. 部署流程

### 5.1 部署策略

| 策略 | 说明 | 适用场景 |
|------|------|---------|
| **自动部署** | 合并到 main 后自动部署 | 开发/测试环境 |
| **手动部署** | 通过 workflow_dispatch 触发 | 生产环境 |
| **标签部署** | 推送 tag 后自动部署 | 版本发布 |

### 5.2 部署环境

| 环境 | 触发条件 | 部署目标 |
|------|---------|---------|
| **开发** | Push to develop | 开发集群 |
| **测试** | Push to main | 测试集群 |
| **生产** | Tag v* | 生产集群 |

---

## 6. 最佳实践

### 6.1 CI 最佳实践

1. **快速反馈**：CI 流程应该在 10 分钟内完成
2. **并行执行**：使用矩阵策略并行运行测试
3. **缓存依赖**：缓存 Go 模块和 Docker 层
4. **失败快速**：在早期阶段失败，避免浪费资源

### 6.2 CD 最佳实践

1. **蓝绿部署**：零停机部署
2. **金丝雀发布**：逐步流量切换
3. **回滚机制**：快速回滚到上一个版本
4. **健康检查**：部署后验证服务健康

### 6.3 安全最佳实践

1. **密钥管理**：使用 GitHub Secrets
2. **镜像扫描**：扫描 Docker 镜像漏洞
3. **最小权限**：使用最小权限原则
4. **审计日志**：记录所有部署操作

---

**最后更新**: 2025-01-XX
