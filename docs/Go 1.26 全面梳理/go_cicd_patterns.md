# Go 1.23 项目CI/CD持续工作流全面指南

> 本文档系统性地梳理了 **Go 1.23** 语言项目从代码提交到生产部署的完整CI/CD工作流，涵盖GitHub Actions、GitLab CI、Docker、Kubernetes、Helm、Terraform等核心技术栈。
>
> **Go 1.23 更新**：
>
> - PGO（Profile Guided Optimization）编译时间开销大幅降低
> - 编译器栈帧重叠优化减少内存使用
> - 支持 `GOEXPERIMENT=aliastypeparams` 泛型类型别名
> - 新的 `go mod tidy -diff` 命令预览依赖变更

---

## 目录

- [Go 1.23 项目CI/CD持续工作流全面指南](#go-123-项目cicd持续工作流全面指南)
  - [目录](#目录)
  - [1. GitHub Actions工作流](#1-github-actions工作流)
    - [1.1 概念定义](#11-概念定义)
    - [1.2 流程图](#12-流程图)
    - [1.3 完整配置示例](#13-完整配置示例)
      - [基础Go项目CI配置](#基础go项目ci配置)
    - [1.4 触发器配置详解](#14-触发器配置详解)
    - [1.5 缓存优化配置](#15-缓存优化配置)
    - [1.6 反例说明](#16-反例说明)
    - [1.7 优化建议](#17-优化建议)
    - [1.8 最佳实践总结](#18-最佳实践总结)
  - [2. GitLab CI工作流](#2-gitlab-ci工作流)
    - [2.1 概念定义](#21-概念定义)
    - [2.2 流程图](#22-流程图)
    - [2.3 完整配置示例](#23-完整配置示例)
    - [2.4 作业依赖配置](#24-作业依赖配置)
    - [2.5 缓存与制品配置](#25-缓存与制品配置)
    - [2.6 Runner配置](#26-runner配置)
    - [2.7 反例说明](#27-反例说明)
    - [2.8 最佳实践总结](#28-最佳实践总结)
  - [3. Docker容器化](#3-docker容器化)
    - [3.1 概念定义](#31-概念定义)
    - [3.2 Dockerfile最佳实践](#32-dockerfile最佳实践)
    - [3.3 多阶段构建详解](#33-多阶段构建详解)
    - [3.4 镜像优化技巧](#34-镜像优化技巧)
    - [3.5 安全扫描集成](#35-安全扫描集成)
    - [3.6 Go项目Docker化完整示例](#36-go项目docker化完整示例)
    - [3.7 Docker Compose配置](#37-docker-compose配置)
    - [3.8 反例说明](#38-反例说明)
    - [3.9 最佳实践总结](#39-最佳实践总结)
  - [4. Kubernetes部署](#4-kubernetes部署)
    - [4.1 概念定义](#41-概念定义)
    - [4.2 流程图](#42-流程图)
    - [4.3 Deployment配置](#43-deployment配置)
    - [4.4 ConfigMap与Secret配置](#44-configmap与secret配置)
    - [4.5 HPA自动扩缩容](#45-hpa自动扩缩容)
    - [4.6 滚动更新策略](#46-滚动更新策略)
    - [4.7 完整K8s配置示例](#47-完整k8s配置示例)
    - [4.8 反例说明](#48-反例说明)
    - [4.9 最佳实践总结](#49-最佳实践总结)
  - [5. Helm Chart](#5-helm-chart)
    - [5.1 概念定义](#51-概念定义)
    - [5.2 Chart结构](#52-chart结构)
    - [5.3 Chart.yaml配置](#53-chartyaml配置)
    - [5.4 Values配置](#54-values配置)
    - [5.5 模板编写](#55-模板编写)
    - [5.6 依赖管理](#56-依赖管理)
    - [5.7 Go应用Helm示例](#57-go应用helm示例)
    - [5.8 反例说明](#58-反例说明)
    - [5.9 最佳实践总结](#59-最佳实践总结)
  - [6. Terraform基础设施](#6-terraform基础设施)
    - [6.1 概念定义](#61-概念定义)
    - [6.2 基础设施即代码流程](#62-基础设施即代码流程)
    - [6.3 Provider配置](#63-provider配置)
    - [6.4 资源管理](#64-资源管理)
    - [6.5 状态管理](#65-状态管理)
    - [6.6 Go项目Terraform示例](#66-go项目terraform示例)
    - [6.7 反例说明](#67-反例说明)
    - [6.8 最佳实践总结](#68-最佳实践总结)
  - [7. 持续集成流程](#7-持续集成流程)
    - [7.1 概念定义](#71-概念定义)
    - [7.2 流程图](#72-流程图)
    - [7.3 代码检查（Lint）](#73-代码检查lint)
    - [7.4 单元测试](#74-单元测试)
    - [7.5 代码覆盖率](#75-代码覆盖率)
    - [7.6 安全扫描](#76-安全扫描)
    - [7.7 完整CI流程示例](#77-完整ci流程示例)
    - [7.8 反例说明](#78-反例说明)
    - [7.9 最佳实践总结](#79-最佳实践总结)
  - [8. 持续部署流程](#8-持续部署流程)
    - [8.1 概念定义](#81-概念定义)
    - [8.2 流程图](#82-流程图)
    - [8.3 环境管理](#83-环境管理)
    - [8.4 蓝绿部署](#84-蓝绿部署)
    - [8.5 金丝雀发布](#85-金丝雀发布)
    - [8.6 功能开关](#86-功能开关)
    - [8.7 回滚策略](#87-回滚策略)
    - [8.8 完整CD流程示例](#88-完整cd流程示例)
    - [8.9 最佳实践总结](#89-最佳实践总结)
  - [9. 制品管理](#9-制品管理)
    - [9.1 概念定义](#91-概念定义)
    - [9.2 镜像仓库](#92-镜像仓库)
    - [9.3 二进制仓库](#93-二进制仓库)
    - [9.4 版本管理](#94-版本管理)
    - [9.5 签名验证](#95-签名验证)
    - [9.6 最佳实践总结](#96-最佳实践总结)
  - [10. 监控与可观测性集成](#10-监控与可观测性集成)
    - [10.1 概念定义](#101-概念定义)
    - [10.2 日志收集](#102-日志收集)
    - [10.3 指标监控](#103-指标监控)
    - [10.4 告警配置](#104-告警配置)
    - [10.5 Kubernetes监控集成](#105-kubernetes监控集成)
    - [10.6 Grafana仪表板](#106-grafana仪表板)
    - [10.7 分布式追踪](#107-分布式追踪)
    - [10.8 最佳实践总结](#108-最佳实践总结)
  - [附录：完整项目结构示例](#附录完整项目结构示例)
  - [总结](#总结)

---

## 1. GitHub Actions工作流

### 1.1 概念定义

**GitHub Actions** 是GitHub提供的持续集成和持续部署（CI/CD）平台，允许开发者自动化软件开发工作流。对于Go项目而言，GitHub Actions能够在代码推送、PR创建、发布等事件触发时自动执行构建、测试、部署等任务。

**核心组件：**

- **Workflow（工作流）**：定义在`.github/workflows/`目录下的YAML文件
- **Event（事件）**：触发工作流的GitHub活动（push、pull_request、release等）
- **Job（作业）**：工作流中的执行单元，可并行或串行执行
- **Step（步骤）**：作业中的单个任务，按顺序执行
- **Action（动作）**：可复用的自动化单元

### 1.2 流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                    GitHub Actions 工作流                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐      │
│  │  Push   │───▶│  PR     │───▶│  Build  │───▶│  Test   │      │
│  │  Event  │    │  Check  │    │  Job    │    │  Job    │      │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘      │
│       │                              │               │         │
│       │                              ▼               ▼         │
│       │                       ┌─────────┐    ┌─────────┐      │
│       │                       │  Cache  │    │  Lint   │      │
│       │                       │  Setup  │    │  Check  │      │
│       │                       └─────────┘    └─────────┘      │
│       │                                                         │
│       ▼                                                         │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐      │
│  │ Matrix  │───▶│  Cross  │───▶│ Docker  │───▶│ Deploy  │      │
│  │ Build   │    │ Platform│    │  Build  │    │  Job    │      │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 1.3 完整配置示例

#### 基础Go项目CI配置

```yaml
# .github/workflows/ci.yml
name: Go CI

on:
  push:
    branches: [main, develop]
    paths-ignore:
      - '**.md'
      - 'docs/**'
  pull_request:
    branches: [main]

env:
  GO_VERSION: '1.21'
  GOLANGCI_LINT_VERSION: v1.55

jobs:
  # 代码质量检查
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: ${{ env.GOLANGCI_LINT_VERSION }}
          args: --timeout=5m --config=.golangci.yml
          skip-cache: false

  # 单元测试
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Download dependencies
        run: go mod download

      - name: Run tests with coverage
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html

      - name: Upload coverage report
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: |
            coverage.out
            coverage.html

      - name: Upload to Codecov
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
          flags: unittests
          name: codecov-umbrella

  # 矩阵构建
  build:
    name: Build
    runs-on: ${{ matrix.os }}
    needs: [lint, test]
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.20', '1.21']
        include:
          - os: ubuntu-latest
            go: '1.21'
            target: linux-amd64
          - os: macos-latest
            go: '1.21'
            target: darwin-amd64
          - os: windows-latest
            go: '1.21'
            target: windows-amd64
      fail-fast: false

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
          cache: true

      - name: Build binary
        run: |
          go build -v -ldflags="-s -w -X main.version=${{ github.ref_name }}" \
            -o bin/myapp-${{ matrix.target }} ./cmd/myapp

      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: myapp-${{ matrix.target }}
          path: bin/myapp-${{ matrix.target }}

  # Docker构建与推送
  docker:
    name: Docker Build
    runs-on: ubuntu-latest
    needs: [build]
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ secrets.DOCKER_USERNAME }}/myapp
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha,prefix=,suffix=,format=short

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
```

### 1.4 触发器配置详解

```yaml
# .github/workflows/triggers.yml
name: Trigger Examples

on:
  # 分支推送触发
  push:
    branches:
      - main
      - 'release/**'
      - '!release/**-alpha'
    tags:
      - 'v*'
    paths:
      - '**.go'
      - 'go.mod'
      - 'go.sum'

  # PR触发
  pull_request:
    types: [opened, synchronize, reopened, ready_for_review]
    branches:
      - main
      - develop

  # 定时触发 (CRON表达式)
  schedule:
    - cron: '0 2 * * *'  # 每天凌晨2点
    - cron: '0 0 * * 0'  # 每周日午夜

  # 手动触发
  workflow_dispatch:
    inputs:
      environment:
        description: '部署环境'
        required: true
        default: 'staging'
        type: choice
        options:
          - staging
          - production
      debug_enabled:
        description: '启用调试'
        required: false
        default: false
        type: boolean

  # 发布触发
  release:
    types: [published, created]

  # 外部事件触发
  repository_dispatch:
    types: [deploy-command]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Echo inputs
        run: |
          echo "Environment: ${{ github.event.inputs.environment }}"
          echo "Debug: ${{ github.event.inputs.debug_enabled }}"
```

### 1.5 缓存优化配置

```yaml
# .github/workflows/cache-optimized.yml
name: Cache Optimized CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      # Go模块缓存
      - name: Go Module Cache
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      # 工具缓存 (golangci-lint等)
      - name: Tools Cache
        uses: actions/cache@v3
        with:
          path: |
            ~/go/bin
            ~/.cache/golangci-lint
          key: ${{ runner.os }}-tools-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-tools-

      # Docker层缓存
      - name: Docker Layer Cache
        uses: actions/cache@v3
        with:
          path: /tmp/.buildx-cache
          key: ${{ runner.os }}-buildx-${{ github.sha }}
          restore-keys: |
            ${{ runner.os }}-buildx-

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Download dependencies
        run: go mod download

      - name: Build
        run: go build -v ./...
```

### 1.6 反例说明

```yaml
# ❌ 错误配置示例
name: Bad CI Config

on: push  # 过于宽泛，每次推送都触发

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # ❌ 没有指定Go版本，使用默认版本
      - uses: actions/setup-go@v5

      # ❌ 没有缓存，每次都要重新下载依赖
      - run: go mod download

      # ❌ 没有并行化，串行执行耗时
      - run: go test ./...
      - run: go build ./...
      - run: golangci-lint run

      # ❌ 硬编码敏感信息
      - run: docker login -u admin -p password123

      # ❌ 没有条件判断，所有分支都部署
      - run: kubectl apply -f k8s/
```

**问题分析：**

1. **触发器过于宽泛**：没有分支过滤，导致所有推送都触发CI
2. **未指定Go版本**：使用默认版本可能导致构建不一致
3. **缺少缓存**：每次构建都重新下载依赖，浪费时间
4. **串行执行**：测试、构建、lint可以并行执行
5. **硬编码凭证**：安全风险，应使用Secrets
6. **无条件部署**：所有分支都执行部署操作

### 1.7 优化建议

```yaml
# ✅ 优化后的配置
name: Optimized CI

on:
  push:
    branches: [main, develop]
    paths-ignore: ['**.md', 'docs/**']
  pull_request:
    branches: [main]

env:
  GO_VERSION: '1.21'

jobs:
  # 并行执行lint、test、build
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      - uses: golangci/golangci-lint-action@v3

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      - run: go test -race -cover ./...

  build:
    runs-on: ubuntu-latest
    needs: [lint, test]  # 依赖前置任务
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      - run: go build -ldflags="-s -w" ./...

  deploy:
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'  # 仅main分支部署
    environment: production  # 需要审批
    steps:
      - uses: actions/checkout@v4
      - name: Deploy
        run: echo "Deploying..."
        env:
          KUBE_CONFIG: ${{ secrets.KUBE_CONFIG }}  # 使用Secrets
```

### 1.8 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 触发器 | 精确配置分支和路径过滤 | 使用宽泛的`on: push` |
| Go版本 | 显式指定版本，使用env变量 | 依赖默认版本 |
| 缓存 | 启用Go模块和构建缓存 | 每次重新下载依赖 |
| 并行化 | lint/test/build并行执行 | 串行执行所有任务 |
| 安全 | 使用Secrets管理凭证 | 硬编码敏感信息 |
| 部署 | 添加条件判断和环境保护 | 无条件部署 |
| 矩阵 | 使用矩阵测试多版本/平台 | 重复定义多个job |

---

## 2. GitLab CI工作流

### 2.1 概念定义

**GitLab CI/CD** 是GitLab内置的持续集成和持续部署系统，通过`.gitlab-ci.yml`文件定义流水线。与GitHub Actions相比，GitLab CI/CD更强调Pipeline的概念，所有任务按阶段（Stage）组织。

**核心概念：**

- **Pipeline（流水线）**：由多个阶段组成的完整CI/CD流程
- **Stage（阶段）**：流水线的执行阶段（build、test、deploy）
- **Job（作业）**：阶段中的具体任务
- **Runner（执行器）**：执行Job的代理程序
- **Artifact（制品）**：Job生成的可传递文件
- **Cache（缓存）**：跨Pipeline共享的数据

### 2.2 流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                    GitLab CI Pipeline                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│  │ .pre    │──▶│ build   │──▶│ test    │──▶│ deploy  │         │
│  │ (准备)  │   │ (构建)  │   │ (测试)  │   │ (部署)  │         │
│  └─────────┘   └─────────┘   └─────────┘   └─────────┘         │
│                    │              │              │              │
│                    ▼              ▼              ▼              │
│              ┌─────────┐   ┌─────────┐   ┌─────────┐           │
│              │ compile │   │ unit    │   │ staging │           │
│              │         │   │ test    │   │ deploy  │           │
│              └─────────┘   └─────────┘   └─────────┘           │
│                    │              │              │              │
│                    │         ┌─────────┐   ┌─────────┐         │
│                    │         │ integ   │   │ prod    │         │
│                    │         │ test    │   │ deploy  │         │
│                    │         └─────────┘   └─────────┘         │
│                    │                                             │
│              ┌─────────┐                                        │
│              │ docker  │                                        │
│              │ build   │                                        │
│              └─────────┘                                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2.3 完整配置示例

```yaml
# .gitlab-ci.yml
# Go项目完整CI/CD配置

# 定义变量
variables:
  GO_VERSION: "1.21"
  DOCKER_REGISTRY: "registry.gitlab.com"
  IMAGE_NAME: "$DOCKER_REGISTRY/$CI_PROJECT_PATH"
  GOLANGCI_LINT_VERSION: "v1.55.0"

# 定义阶段
stages:
  - prepare
  - build
  - test
  - security
  - package
  - deploy

# 全局默认配置
default:
  image: golang:${GO_VERSION}-alpine
  before_script:
    - apk add --no-cache git ca-certificates tzdata
    - go env -w GOPROXY=https://proxy.golang.org,direct
    - go env -w GOSUMDB=sum.golang.org

# ============ Prepare Stage ============
# 依赖下载（缓存优化）
download-deps:
  stage: prepare
  script:
    - go mod download
    - go mod verify
  cache:
    key: ${CI_COMMIT_REF_SLUG}
    paths:
      - ${GOPATH}/pkg/mod/

# ============ Build Stage ============
# 编译应用
compile:
  stage: build
  script:
    - CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=${CI_COMMIT_TAG:-$CI_COMMIT_SHORT_SHA}" -o bin/myapp ./cmd/myapp
  artifacts:
    paths:
      - bin/myapp
    expire_in: 1 hour
  needs:
    - job: download-deps
      artifacts: false

# 交叉编译（多平台）
cross-compile:
  stage: build
  parallel:
    matrix:
      - GOOS: linux
        GOARCH: [amd64, arm64]
      - GOOS: darwin
        GOARCH: [amd64, arm64]
      - GOOS: windows
        GOARCH: amd64
  script:
    - |
      OUTPUT="bin/myapp-${GOOS}-${GOARCH}"
      if [ "$GOOS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
      fi
      CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o ${OUTPUT} ./cmd/myapp
  artifacts:
    paths:
      - bin/myapp-*
    expire_in: 1 day

# ============ Test Stage ============
# 单元测试
unit-test:
  stage: test
  script:
    - go test -v -race -coverprofile=coverage.out ./...
    - go tool cover -func=coverage.out
  coverage: '/total:\s+\(statements\)\s+(\d+\.\d+)%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml
    paths:
      - coverage.out
  needs:
    - job: download-deps
      artifacts: false

# 集成测试
integration-test:
  stage: test
  services:
    - name: postgres:15-alpine
      alias: postgres
    - name: redis:7-alpine
      alias: redis
  variables:
    POSTGRES_DB: testdb
    POSTGRES_USER: testuser
    POSTGRES_PASSWORD: testpass
    REDIS_HOST: redis
  script:
    - go test -v -tags=integration ./tests/integration/...
  needs:
    - job: unit-test
      artifacts: false

# ============ Security Stage ============
# 依赖漏洞扫描
dependency-scan:
  stage: security
  image: securego/gosec:latest
  script:
    - gosec -fmt sarif -out gosec-report.sarif ./...
  artifacts:
    reports:
      sast: gosec-report.sarif
    paths:
      - gosec-report.sarif
  allow_failure: true

# 代码安全扫描
sast:
  stage: security
  image: returntocorp/semgrep
  script:
    - semgrep --config=auto --json --output=semgrep-report.json .
  artifacts:
    paths:
      - semgrep-report.json
  allow_failure: true

# ============ Package Stage ============
# Docker镜像构建
docker-build:
  stage: package
  image: docker:24-dind
  services:
    - docker:24-dind
  variables:
    DOCKER_DRIVER: overlay2
    DOCKER_TLS_CERTDIR: "/certs"
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - |
      docker build \
        --build-arg GO_VERSION=${GO_VERSION} \
        --cache-from ${IMAGE_NAME}:latest \
        -t ${IMAGE_NAME}:${CI_COMMIT_SHORT_SHA} \
        -t ${IMAGE_NAME}:latest \
        .
    - docker push ${IMAGE_NAME}:${CI_COMMIT_SHORT_SHA}
    - docker push ${IMAGE_NAME}:latest
  only:
    - main
    - tags

# ============ Deploy Stage ============
# 部署到Staging
deploy-staging:
  stage: deploy
  image: bitnami/kubectl:latest
  environment:
    name: staging
    url: https://staging.example.com
  script:
    - kubectl config use-context staging
    - kubectl set image deployment/myapp myapp=${IMAGE_NAME}:${CI_COMMIT_SHORT_SHA}
    - kubectl rollout status deployment/myapp
  only:
    - main

# 部署到Production（需要手动触发）
deploy-production:
  stage: deploy
  image: bitnami/kubectl:latest
  environment:
    name: production
    url: https://production.example.com
  script:
    - kubectl config use-context production
    - kubectl set image deployment/myapp myapp=${IMAGE_NAME}:${CI_COMMIT_TAG}
    - kubectl rollout status deployment/myapp
  when: manual
  only:
    - tags
```

### 2.4 作业依赖配置

```yaml
# 作业依赖示例
stages:
  - build
  - test
  - deploy

# 基础构建
build-base:
  stage: build
  script:
    - go build -o bin/base ./cmd/base
  artifacts:
    paths:
      - bin/base

# 依赖基础构建的扩展构建
build-extended:
  stage: build
  needs:
    - job: build-base
      artifacts: true  # 需要base的制品
  script:
    - go build -o bin/extended ./cmd/extended
  artifacts:
    paths:
      - bin/extended

# 测试需要两个构建完成
test-all:
  stage: test
  needs:
    - job: build-base
      artifacts: true
    - job: build-extended
      artifacts: true
  script:
    - go test ./...

# 部署需要测试通过
deploy:
  stage: deploy
  needs:
    - job: test-all
      artifacts: false
  script:
    - echo "Deploying..."
```

### 2.5 缓存与制品配置

```yaml
# 缓存和制品优化配置
variables:
  GOPATH: $CI_PROJECT_DIR/.go

# 依赖缓存Job
cache-deps:
  stage: .pre
  script:
    - go mod download
  cache:
    key:
      files:
        - go.mod
        - go.sum
    paths:
      - ${GOPATH}/pkg/mod/
    policy: pull-push

# 使用缓存的Job
build-with-cache:
  stage: build
  script:
    - go build ./...
  cache:
    key:
      files:
        - go.mod
        - go.sum
    paths:
      - ${GOPATH}/pkg/mod/
    policy: pull  # 只拉取，不更新

# 制品传递
generate-report:
  stage: build
  script:
    - go test -coverprofile=coverage.out ./...
    - go tool cover -html=coverage.out -o coverage.html
  artifacts:
    name: "coverage-$CI_COMMIT_REF_NAME"
    paths:
      - coverage.out
      - coverage.html
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml
    expire_in: 1 week
    when: always  # 即使失败也保存制品
```

### 2.6 Runner配置

```toml
# /etc/gitlab-runner/config.toml
# GitLab Runner配置示例

[[runners]]
  name = "docker-runner-01"
  url = "https://gitlab.com/"
  token = "YOUR_REGISTRATION_TOKEN"
  executor = "docker"

  [runners.docker]
    tls_verify = false
    image = "golang:1.21-alpine"
    privileged = true
    disable_cache = false
    volumes = [
      "/var/run/docker.sock:/var/run/docker.sock",
      "/cache"
    ]
    shm_size = 0

  [runners.cache]
    Type = "s3"
    Path = "gitlab-runner"
    Shared = true
    [runners.cache.s3]
      ServerAddress = "s3.amazonaws.com"
      BucketName = "gitlab-runner-cache"
      BucketLocation = "us-east-1"
      Insecure = false

  [runners.kubernetes]
    host = ""
    bearer_token = ""
    namespace = "gitlab-runner"
    image = "golang:1.21-alpine"
    privileged = true
```

### 2.7 反例说明

```yaml
# ❌ 错误配置示例
stages:
  - build
  - test
  - deploy

# ❌ 没有缓存，每次都要重新下载依赖
build:
  stage: build
  script:
    - go mod download  # 每次都执行
    - go build ./...

# ❌ 没有制品传递，后续Job需要重新构建
test:
  stage: test
  script:
    - go build ./...  # 重复构建
    - go test ./...

# ❌ 没有条件限制，所有分支都部署
deploy:
  stage: deploy
  script:
    - kubectl apply -f k8s/

# ❌ 没有资源限制，可能耗尽Runner资源
heavy-job:
  stage: build
  script:
    - go test -parallel 100 ./...
```

### 2.8 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 阶段划分 | 清晰的build/test/deploy阶段 | 所有任务在一个阶段 |
| 缓存策略 | 使用go.mod作为缓存key | 不使用缓存 |
| 制品传递 | 合理使用artifacts传递 | 每个Job重复构建 |
| 条件执行 | 使用only/when控制执行 | 无条件执行所有任务 |
| 资源限制 | 设置合理的并行度和超时 | 无限制的资源使用 |
| 安全扫描 | 集成SAST和依赖扫描 | 忽略安全检测 |
| 环境管理 | 使用environment配置 | 直接执行部署命令 |

---

## 3. Docker容器化

### 3.1 概念定义

**Docker容器化** 是将Go应用程序及其依赖打包成可移植容器镜像的过程。对于Go项目，容器化提供了环境一致性、快速部署和资源隔离的优势。

**核心优势：**

- **环境一致性**：开发、测试、生产环境完全一致
- **快速部署**：秒级启动，快速扩缩容
- **资源隔离**：进程、网络、文件系统隔离
- **可移植性**：一次构建，到处运行

### 3.2 Dockerfile最佳实践

```dockerfile
# ============================================
# 生产级Go应用Dockerfile（多阶段构建）
# ============================================

# ------------------ 阶段1：构建 ------------------
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置Go环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

# 先复制go.mod和go.sum（利用Docker缓存层）
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建二进制文件
# -ldflags参数说明：
#   -s: 去掉符号表
#   -w: 去掉调试信息
#   -X: 设置变量值（版本信息）
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN go build -ldflags="-s -w \
    -X main.version=${VERSION} \
    -X main.buildTime=${BUILD_TIME} \
    -X main.gitCommit=${GIT_COMMIT}" \
    -o /app/bin/server ./cmd/server

# ------------------ 阶段2：运行 ------------------
# 使用distroless镜像，最小化攻击面
FROM gcr.io/distroless/static:nonroot

# 设置时区
ENV TZ=UTC

# 从builder阶段复制证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 复制二进制文件
COPY --from=builder /app/bin/server /server

# 使用非root用户运行
USER nonroot:nonroot

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "-health-check"] || exit 1

# 启动应用
ENTRYPOINT ["/server"]
```

### 3.3 多阶段构建详解

```dockerfile
# ============================================
# 完整多阶段构建示例（开发/测试/生产）
# ============================================

# ------------------ 基础阶段 ------------------
FROM golang:1.21-alpine AS base
WORKDIR /app
RUN apk add --no-cache git ca-certificates
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download

# ------------------ 开发阶段 ------------------
FROM base AS development
RUN apk add --no-cache air  # 热重载工具
COPY . .
CMD ["air", "-c", ".air.toml"]

# ------------------ 测试阶段 ------------------
FROM base AS testing
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest
COPY . .
RUN go test -v ./...

# ------------------ 构建阶段 ------------------
FROM base AS build
COPY . .
RUN go build -ldflags="-s -w" -o /app/server ./cmd/server

# ------------------ 生产阶段 ------------------
FROM gcr.io/distroless/static:nonroot AS production
COPY --from=build /app/server /server
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/server"]
```

### 3.4 镜像优化技巧

```dockerfile
# ============================================
# 镜像优化配置
# ============================================

# 1. 使用最小的基础镜像
FROM alpine:3.18 AS minimal
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /server
ENTRYPOINT ["/server"]

# 2. 使用scratch镜像（完全空镜像）
FROM scratch AS scratch-image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /server
# 注意：scratch镜像没有shell，无法使用HEALTHCHECK
ENTRYPOINT ["/server"]

# 3. 使用distroless镜像（推荐）
FROM gcr.io/distroless/static-debian12:nonroot AS distroless
COPY --from=builder /app/server /server
USER nonroot
ENTRYPOINT ["/server"]

# 4. 多架构构建
# 使用docker buildx构建多平台镜像
# docker buildx build --platform linux/amd64,linux/arm64 -t myapp:latest .
```

### 3.5 安全扫描集成

```dockerfile
# Dockerfile with security considerations
FROM golang:1.21-alpine AS builder

# 创建非root用户
RUN adduser -D -g '' appuser

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/server

# 使用distroless镜像
FROM gcr.io/distroless/static:nonroot

# 只复制必要的文件
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/server /server

# 以非root用户运行
USER nonroot

# 只暴露必要的端口
EXPOSE 8080

# 使用exec格式避免shell
ENTRYPOINT ["/server"]
```

### 3.6 Go项目Docker化完整示例

```dockerfile
# ============================================
# 生产级Go Web服务Dockerfile
# ============================================

# 构建阶段
FROM golang:1.21-alpine AS builder

# 元数据标签
LABEL maintainer="dev@example.com"
LABEL version="1.0.0"
LABEL description="Go Web Service"

WORKDIR /build

# 安装依赖
RUN apk add --no-cache --virtual .build-deps \
    git \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# 设置环境变量
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct

# 复制并下载依赖（利用缓存）
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建参数
ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT
ARG GIT_BRANCH

# 构建应用
RUN go build -a -installsuffix cgo -ldflags="\
    -s -w \
    -extldflags '-static' \
    -X 'main.Version=${VERSION}' \
    -X 'main.BuildTime=${BUILD_TIME}' \
    -X 'main.GitCommit=${GIT_COMMIT}' \
    -X 'main.GitBranch=${GIT_BRANCH}' \
    " -o /build/bin/server ./cmd/server

# 运行阶段
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# 从builder复制必要文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/bin/server /server

# 非root用户
USER nonroot:nonroot

# 端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/server", "health"] || exit 1

# 启动
ENTRYPOINT ["/server"]
```

### 3.7 Docker Compose配置

```yaml
# docker-compose.yml
version: '3.8'

services:
  # Go应用服务
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: production
      args:
        VERSION: ${VERSION:-dev}
        BUILD_TIME: ${BUILD_TIME}
        GIT_COMMIT: ${GIT_COMMIT}
    image: myapp:${VERSION:-latest}
    container_name: myapp
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - LOG_LEVEL=info
      - DATABASE_URL=postgres://user:pass@postgres:5432/myapp?sslmode=disable
      - REDIS_URL=redis:6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "/server", "health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s

  # 数据库
  postgres:
    image: postgres:15-alpine
    container_name: myapp-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: myapp
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d myapp"]
      interval: 5s
      timeout: 5s
      retries: 5

  # 缓存
  redis:
    image: redis:7-alpine
    container_name: myapp-redis
    restart: unless-stopped
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  # 开发模式热重载
  app-dev:
    build:
      context: .
      dockerfile: Dockerfile
      target: development
    image: myapp:dev
    container_name: myapp-dev
    volumes:
      - .:/app
      - go_mod_cache:/go/pkg/mod
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      - LOG_LEVEL=debug
    profiles:
      - dev
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:
  go_mod_cache:

networks:
  app-network:
    driver: bridge
```

### 3.8 反例说明

```dockerfile
# ❌ 错误配置示例

# 1. 不使用多阶段构建
FROM golang:1.21
WORKDIR /app
COPY . .
RUN go build ./...
CMD ["./server"]
# 问题：镜像包含完整的Go SDK，体积巨大

# 2. 以root运行
FROM golang:1.21-alpine
COPY . .
RUN go build -o server ./cmd/server
CMD ["./server"]
# 问题：以root运行，安全风险

# 3. 复制不必要的文件
FROM golang:1.21-alpine
COPY . .  # 复制了.git, vendor等不需要的文件
RUN go build ./...
# 问题：增加了镜像体积，可能泄露敏感信息

# 4. 不使用.dockerignore
# 没有.dockerignore文件
# 问题：将本地开发文件复制到镜像

# 5. 硬编码敏感信息
ENV DATABASE_PASSWORD=mypassword123
# 问题：密码泄露在镜像中
```

### 3.9 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 基础镜像 | 使用distroless或alpine | 使用完整Go镜像 |
| 多阶段构建 | 分离构建和运行阶段 | 单阶段构建 |
| 用户权限 | 使用非root用户 | 以root运行 |
| 缓存优化 | 先复制go.mod再复制代码 | 一次性复制所有文件 |
| 安全扫描 | 集成Trivy/Clair扫描 | 忽略安全检测 |
| 健康检查 | 配置HEALTHCHECK | 无健康检查 |
| 镜像体积 | 使用-ldflags优化 | 包含调试信息 |

---

## 4. Kubernetes部署

### 4.1 概念定义

**Kubernetes（K8s）** 是用于自动部署、扩展和管理容器化应用程序的开源容器编排平台。对于Go应用，Kubernetes提供了声明式配置、自动扩缩容、服务发现和负载均衡等能力。

**核心概念：**

- **Pod**：K8s中最小的部署单元，包含一个或多个容器
- **Deployment**：管理Pod的副本集，支持滚动更新
- **Service**：提供Pod的网络访问和负载均衡
- **ConfigMap/Secret**：管理配置和敏感数据
- **HPA**：水平Pod自动扩缩容
- **Ingress**：管理外部访问的路由规则

### 4.2 流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                    Kubernetes 部署架构                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐                                                │
│  │   Ingress   │  外部流量入口                                   │
│  │  Controller │                                                │
│  └──────┬──────┘                                                │
│         │                                                       │
│         ▼                                                       │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐       │
│  │   Service   │────▶│   Service   │────▶│   Service   │       │
│  │  (ClusterIP)│     │  (ClusterIP)│     │  (ClusterIP)│       │
│  └──────┬──────┘     └──────┬──────┘     └──────┬──────┘       │
│         │                   │                   │              │
│         ▼                   ▼                   ▼              │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐       │
│  │ Deployment  │     │ Deployment  │     │ Deployment  │       │
│  │  (App Pod)  │     │  (App Pod)  │     │  (App Pod)  │       │
│  │  ┌───────┐  │     │  ┌───────┐  │     │  ┌───────┐  │       │
│  │  │Pod x3 │  │     │  │Pod x2 │  │     │  │Pod x2 │  │       │
│  │  └───────┘  │     │  └───────┘  │     │  └───────┘  │       │
│  └─────────────┘     └─────────────┘     └─────────────┘       │
│         │                   │                   │              │
│         ▼                   ▼                   ▼              │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐       │
│  │ ConfigMap   │     │   Secret    │     │   Secret    │       │
│  │  (配置)     │     │  (敏感数据) │     │  (TLS证书)  │       │
│  └─────────────┘     └─────────────┘     └─────────────┘       │
│                                                                 │
│  ┌─────────────┐                                                │
│  │    HPA      │  自动扩缩容                                     │
│  │ (Horizontal │                                                │
│  │Pod Autoscaler│                                               │
│  └─────────────┘                                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 4.3 Deployment配置

```yaml
# ============================================
# Deployment配置 - 生产级Go应用
# ============================================

apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
  namespace: production
  labels:
    app: go-app
    version: v1.0.0
    tier: backend
spec:
  # 副本数
  replicas: 3

  # 选择器
  selector:
    matchLabels:
      app: go-app

  # 滚动更新策略
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%        # 更新时最多可多出的Pod数
      maxUnavailable: 25%  # 更新时最多不可用的Pod数

  # Pod模板
  template:
    metadata:
      labels:
        app: go-app
        version: v1.0.0
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      # 亲和性配置
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - go-app
                topologyKey: kubernetes.io/hostname

      # 初始化容器
      initContainers:
        - name: migrate
          image: myregistry/go-app-migrate:latest
          command: ["/migrate", "up"]
          envFrom:
            - configMapRef:
                name: go-app-config
            - secretRef:
                name: go-app-secrets

      # 主容器
      containers:
        - name: go-app
          image: myregistry/go-app:v1.0.0
          imagePullPolicy: Always

          # 端口配置
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
            - name: metrics
              containerPort: 9090
              protocol: TCP

          # 环境变量
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP

          envFrom:
            - configMapRef:
                name: go-app-config
            - secretRef:
                name: go-app-secrets

          # 资源限制
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"

          # 健康检查
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 10
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3

          readinessProbe:
            httpGet:
              path: /ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3

          # 启动探针（慢启动应用）
          startupProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
            failureThreshold: 30

          # 安全上下文
          securityContext:
            runAsNonRoot: true
            runAsUser: 65534
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - ALL

          # 卷挂载
          volumeMounts:
            - name: tmp
              mountPath: /tmp
            - name: cache
              mountPath: /cache

      # 卷定义
      volumes:
        - name: tmp
          emptyDir: {}
        - name: cache
          emptyDir:
            sizeLimit: 1Gi

      # 镜像拉取密钥
      imagePullSecrets:
        - name: registry-credentials

      # 重启策略
      restartPolicy: Always

      # 终止优雅期
      terminationGracePeriodSeconds: 30

---
# ============================================
# Service配置
# ============================================
apiVersion: v1
kind: Service
metadata:
  name: go-app
  namespace: production
  labels:
    app: go-app
spec:
  type: ClusterIP
  ports:
    - name: http
      port: 80
      targetPort: 8080
      protocol: TCP
    - name: metrics
      port: 9090
      targetPort: 9090
      protocol: TCP
  selector:
    app: go-app
  sessionAffinity: None

---
# ============================================
# Headless Service（用于有状态应用）
# ============================================
apiVersion: v1
kind: Service
metadata:
  name: go-app-headless
  namespace: production
  labels:
    app: go-app
spec:
  type: ClusterIP
  clusterIP: None  # Headless
  ports:
    - name: http
      port: 80
      targetPort: 8080
  selector:
    app: go-app
```

### 4.4 ConfigMap与Secret配置

```yaml
# ============================================
# ConfigMap - 非敏感配置
# ============================================
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-app-config
  namespace: production
data:
  # 应用配置
  APP_ENV: "production"
  APP_NAME: "go-app"
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"

  # 服务器配置
  SERVER_PORT: "8080"
  SERVER_TIMEOUT: "30s"
  SERVER_MAX_CONNECTIONS: "1000"

  # 数据库配置（非敏感部分）
  DB_HOST: "postgres.production.svc.cluster.local"
  DB_PORT: "5432"
  DB_NAME: "myapp"
  DB_POOL_MAX: "20"
  DB_POOL_MIN: "5"

  # Redis配置
  REDIS_HOST: "redis.production.svc.cluster.local"
  REDIS_PORT: "6379"
  REDIS_DB: "0"

  # 缓存配置
  CACHE_TTL: "3600"

  # 功能开关
  FEATURE_FLAG_NEW_UI: "true"
  FEATURE_FLAG_BETA_API: "false"

---
# ============================================
# Secret - 敏感配置
# ============================================
apiVersion: v1
kind: Secret
metadata:
  name: go-app-secrets
  namespace: production
type: Opaque
stringData:
  # 数据库密码
  DB_PASSWORD: "your-secure-password"
  DB_USER: "app_user"

  # Redis密码
  REDIS_PASSWORD: "redis-password"

  # API密钥
  API_KEY: "your-api-key"
  JWT_SECRET: "your-jwt-secret-key"

  # 外部服务凭证
  EXTERNAL_API_TOKEN: "external-api-token"

  # TLS证书（base64编码）
  tls.crt: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  tls.key: |
    -----BEGIN PRIVATE KEY-----
    ...
    -----END PRIVATE KEY-----

---
# ============================================
# TLS Secret
# ============================================
apiVersion: v1
kind: Secret
metadata:
  name: go-app-tls
  namespace: production
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-certificate>
  tls.key: <base64-encoded-key>
```

### 4.5 HPA自动扩缩容

```yaml
# ============================================
# HorizontalPodAutoscaler - CPU和内存扩缩容
# ============================================
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-app-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-app

  # 副本数范围
  minReplicas: 3
  maxReplicas: 20

  # 扩缩容行为
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300  # 缩容前等待5分钟
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15
        - type: Pods
          value: 4
          periodSeconds: 15
      selectPolicy: Max

  # 指标配置
  metrics:
    # CPU使用率
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70

    # 内存使用率
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80

    # 自定义指标（需要Metrics Server）
    - type: Pods
      pods:
        metric:
          name: http_requests_per_second
        target:
          type: AverageValue
          averageValue: "1000"

    # 外部指标
    - type: External
      external:
        metric:
          name: queue_messages_ready
          selector:
            matchLabels:
              queue: worker_tasks
        target:
          type: AverageValue
          averageValue: "30"

---
# ============================================
# 自定义指标配置（Prometheus Adapter）
# ============================================
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-adapter-config
data:
  config.yaml: |
    rules:
      - seriesQuery: 'http_requests_total{namespace!="",pod!=""}'
        resources:
          overrides:
            namespace:
              resource: namespace
            pod:
              resource: pod
        name:
          matches: "^(.*)_total"
          as: "${1}_per_second"
        metricsQuery: 'sum(rate(<<.Series>>{<<.LabelMatchers>>}[2m])) by (<<.GroupBy>>)'
```

### 4.6 滚动更新策略

```yaml
# ============================================
# 滚动更新配置详解
# ============================================
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
spec:
  replicas: 10

  # 滚动更新策略
  strategy:
    type: RollingUpdate

    rollingUpdate:
      # 更新时最多可超出的Pod数量
      # 可以是绝对值（如：2）或百分比（如：25%）
      maxSurge: 25%

      # 更新时最多不可用的Pod数量
      # 可以是绝对值（如：2）或百分比（如：25%）
      maxUnavailable: 25%

  # 或者使用Recreate策略（先删除再创建）
  # strategy:
  #   type: Recreate

  template:
    spec:
      containers:
        - name: go-app
          image: myregistry/go-app:v2.0.0

          # 生命周期钩子
          lifecycle:
            # 启动前执行
            postStart:
              exec:
                command: ["/bin/sh", "-c", "echo 'Container started'"]

            # 终止前执行（优雅关闭）
            preStop:
              exec:
                command: ["/bin/sh", "-c", "sleep 15 && /server graceful-shutdown"]

      # 优雅终止时间（必须大于preStop执行时间）
      terminationGracePeriodSeconds: 60
```

### 4.7 完整K8s配置示例

```yaml
# ============================================
# 完整的Go应用K8s配置
# 包含：Namespace, Deployment, Service, ConfigMap, Secret, HPA, Ingress
# ============================================

---
# Namespace
apiVersion: v1
kind: Namespace
metadata:
  name: go-app-production
  labels:
    name: go-app-production
    environment: production

---
# ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-app-config
  namespace: go-app-production
data:
  APP_ENV: "production"
  LOG_LEVEL: "info"
  SERVER_PORT: "8080"
  DB_HOST: "postgres"
  DB_PORT: "5432"
  REDIS_HOST: "redis"

---
# Secret
apiVersion: v1
kind: Secret
metadata:
  name: go-app-secrets
  namespace: go-app-production
type: Opaque
stringData:
  DB_PASSWORD: "changeme"
  JWT_SECRET: "your-secret-key"

---
# Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
  namespace: go-app-production
  labels:
    app: go-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-app
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
        - name: go-app
          image: myregistry/go-app:latest
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef:
                name: go-app-config
            - secretRef:
                name: go-app-secrets
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "200m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5

---
# Service
apiVersion: v1
kind: Service
metadata:
  name: go-app
  namespace: go-app-production
spec:
  selector:
    app: go-app
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP

---
# HPA
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-app-hpa
  namespace: go-app-production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-app
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70

---
# Ingress
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-app-ingress
  namespace: go-app-production
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
    - hosts:
        - api.example.com
      secretName: go-app-tls
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: go-app
                port:
                  number: 80
```

### 4.8 反例说明

```yaml
# ❌ 错误配置示例

# 1. 没有资源限制
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: app
          image: myapp:latest
          # 缺少resources配置

# 2. 没有健康检查
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: app
          image: myapp:latest
          # 缺少livenessProbe和readinessProbe

# 3. 使用latest标签
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: app
          image: myapp:latest  # 不稳定，难以回滚

# 4. 没有安全上下文
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: app
          image: myapp:latest
          # 缺少securityContext，以root运行

# 5. 硬编码敏感信息
apiVersion: v1
kind: ConfigMap
data:
  DB_PASSWORD: "mypassword"  # 应该使用Secret
```

### 4.9 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 资源限制 | 设置requests和limits | 无限制配置 |
| 健康检查 | 配置liveness和readiness探针 | 无健康检查 |
| 镜像标签 | 使用语义化版本标签 | 使用latest标签 |
| 安全上下文 | 使用非root用户 | 以root运行 |
| 配置管理 | 使用ConfigMap和Secret | 硬编码配置 |
| 更新策略 | 配置合理的maxSurge/maxUnavailable | 无策略配置 |
| 自动扩缩容 | 配置HPA | 固定副本数 |

---

## 5. Helm Chart

### 5.1 概念定义

**Helm** 是Kubernetes的包管理工具，使用Chart（图表）来定义、安装和升级复杂的Kubernetes应用。对于Go项目，Helm提供了模板化配置、版本管理和依赖管理的能力。

**核心概念：**

- **Chart**：Helm包，包含一组K8s资源模板
- **Release**：Chart的运行实例
- **Values**：配置Chart的参数
- **Template**：使用Go模板语法定义的K8s资源
- **Hooks**：在Release生命周期中执行的任务

### 5.2 Chart结构

```
my-go-app/
├── Chart.yaml          # Chart元数据
├── values.yaml         # 默认配置值
├── values-production.yaml  # 生产环境配置
├── .helmignore         # 忽略文件
├── charts/             # 依赖Chart
│   └── postgresql-12.0.0.tgz
├── templates/          # K8s资源模板
│   ├── _helpers.tpl    # 辅助模板
│   ├── NOTES.txt       # 安装后说明
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── hpa.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── serviceaccount.yaml
│   └── tests/
│       └── test-connection.yaml
└── README.md           # 文档
```

### 5.3 Chart.yaml配置

```yaml
# Chart.yaml - Chart元数据
apiVersion: v2
name: my-go-app
description: A Helm chart for Go application
type: application
version: 1.2.0          # Chart版本
appVersion: "2.0.0"     # 应用版本
kubeVersion: ">=1.19.0-0"
keywords:
  - go
  - golang
  - web
  - api
home: https://github.com/example/my-go-app
sources:
  - https://github.com/example/my-go-app
maintainers:
  - name: DevOps Team
    email: devops@example.com
    url: https://example.com
icon: https://example.com/icon.png
annotations:
  category: Application
dependencies:
  - name: postgresql
    version: "12.x.x"
    repository: https://charts.bitnami.com/bitnami
    condition: postgresql.enabled
    tags:
      - database
  - name: redis
    version: "17.x.x"
    repository: https://charts.bitnami.com/bitnami
    condition: redis.enabled
```

### 5.4 Values配置

```yaml
# values.yaml - 默认配置值

# ============================================
# 全局配置
# ============================================
global:
  imageRegistry: ""
  imagePullSecrets: []
  storageClass: ""

# ============================================
# 镜像配置
# ============================================
image:
  registry: docker.io
  repository: mycompany/go-app
  tag: ""
  pullPolicy: IfNotPresent
  pullSecrets: []

# ============================================
# 副本与部署配置
# ============================================
replicaCount: 3

deployment:
  annotations: {}
  labels: {}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 0

# ============================================
# Pod配置
# ============================================
podAnnotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "8080"
  prometheus.io/path: "/metrics"

podLabels: {}

podSecurityContext:
  runAsNonRoot: true
  runAsUser: 65534
  fsGroup: 65534

securityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop:
      - ALL

# ============================================
# 服务配置
# ============================================
service:
  type: ClusterIP
  port: 80
  targetPort: 8080
  annotations: {}
  labels: {}

# ============================================
# Ingress配置
# ============================================
ingress:
  enabled: true
  className: nginx
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    cert-manager.io/cluster-issuer: "letsencrypt"
  hosts:
    - host: api.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: go-app-tls
      hosts:
        - api.example.com

# ============================================
# 资源限制
# ============================================
resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi

# ============================================
# 健康检查
# ============================================
livenessProbe:
  enabled: true
  path: /health
  port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
  successThreshold: 1

readinessProbe:
  enabled: true
  path: /ready
  port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
  successThreshold: 1

startupProbe:
  enabled: false
  path: /health
  port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 30

# ============================================
# HPA配置
# ============================================
autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
        - type: Percent
          value: 100
          periodSeconds: 15

# ============================================
# 应用配置
# ============================================
config:
  appEnv: production
  logLevel: info
  logFormat: json
  serverPort: 8080
  serverTimeout: 30s

  # 数据库配置
  database:
    host: ""
    port: 5432
    name: myapp
    sslMode: require
    poolMax: 20
    poolMin: 5

  # Redis配置
  redis:
    host: ""
    port: 6379
    db: 0

  # 功能开关
  featureFlags:
    newUI: true
    betaAPI: false

# ============================================
# Secret配置
# ============================================
secrets:
  dbPassword: ""
  jwtSecret: ""
  apiKey: ""

# ============================================
# 持久化存储
# ============================================
persistence:
  enabled: false
  storageClass: ""
  accessMode: ReadWriteOnce
  size: 10Gi
  mountPath: /data

# ============================================
# 服务账户
# ============================================
serviceAccount:
  create: true
  annotations: {}
  name: ""

# ============================================
# RBAC
# ============================================
rbac:
  create: true
  rules:
    - apiGroups: [""]
      resources: ["pods"]
      verbs: ["get", "list", "watch"]

# ============================================
# 依赖配置
# ============================================
postgresql:
  enabled: true
  auth:
    username: myapp
    password: changeme
    database: myapp
  primary:
    persistence:
      enabled: true
      size: 10Gi

redis:
  enabled: true
  auth:
    enabled: true
    password: changeme
  master:
    persistence:
      enabled: true
      size: 5Gi
```

### 5.5 模板编写

```yaml
# templates/_helpers.tpl - 辅助模板
{{/*
Expand the name of the chart.
*/}}
{{- define "my-go-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "my-go-app.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "my-go-app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "my-go-app.labels" -}}
helm.sh/chart: {{ include "my-go-app.chart" . }}
{{ include "my-go-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "my-go-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "my-go-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "my-go-app.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "my-go-app.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the image tag
*/}}
{{- define "my-go-app.imageTag" -}}
{{- default .Chart.AppVersion .Values.image.tag }}
{{- end }}
```

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "my-go-app.fullname" . }}
  labels:
    {{- include "my-go-app.labels" . | nindent 4 }}
    {{- with .Values.deployment.labels }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
  {{- with .Values.deployment.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "my-go-app.selectorLabels" . | nindent 6 }}
  {{- with .Values.deployment.strategy }}
  strategy:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  template:
    metadata:
      {{- with .Values.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "my-go-app.selectorLabels" . | nindent 8 }}
        {{- with .Values.podLabels }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
    spec:
      {{- with .Values.image.pullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}

      serviceAccountName: {{ include "my-go-app.serviceAccountName" . }}

      {{- with .Values.podSecurityContext }}
      securityContext:
        {{- toYaml . | nindent 8 }}
      {{- end }}

      containers:
        - name: {{ .Chart.Name }}
          {{- with .Values.securityContext }}
          securityContext:
            {{- toYaml . | nindent 12 }}
          {{- end }}

          image: "{{ .Values.image.registry }}/{{ .Values.image.repository }}:{{ include "my-go-app.imageTag" . }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}

          ports:
            - name: http
              containerPort: {{ .Values.config.serverPort }}
              protocol: TCP

          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace

          envFrom:
            - configMapRef:
                name: {{ include "my-go-app.fullname" . }}-config
            - secretRef:
                name: {{ include "my-go-app.fullname" . }}-secrets

          {{- with .Values.livenessProbe }}
          {{- if .enabled }}
          livenessProbe:
            httpGet:
              path: {{ .path }}
              port: {{ .port }}
            initialDelaySeconds: {{ .initialDelaySeconds }}
            periodSeconds: {{ .periodSeconds }}
            timeoutSeconds: {{ .timeoutSeconds }}
            failureThreshold: {{ .failureThreshold }}
            successThreshold: {{ .successThreshold }}
          {{- end }}
          {{- end }}

          {{- with .Values.readinessProbe }}
          {{- if .enabled }}
          readinessProbe:
            httpGet:
              path: {{ .path }}
              port: {{ .port }}
            initialDelaySeconds: {{ .initialDelaySeconds }}
            periodSeconds: {{ .periodSeconds }}
            timeoutSeconds: {{ .timeoutSeconds }}
            failureThreshold: {{ .failureThreshold }}
            successThreshold: {{ .successThreshold }}
          {{- end }}
          {{- end }}

          {{- with .Values.startupProbe }}
          {{- if .enabled }}
          startupProbe:
            httpGet:
              path: {{ .path }}
              port: {{ .port }}
            initialDelaySeconds: {{ .initialDelaySeconds }}
            periodSeconds: {{ .periodSeconds }}
            timeoutSeconds: {{ .timeoutSeconds }}
            failureThreshold: {{ .failureThreshold }}
          {{- end }}
          {{- end }}

          {{- with .Values.resources }}
          resources:
            {{- toYaml . | nindent 12 }}
          {{- end }}

          {{- if .Values.persistence.enabled }}
          volumeMounts:
            - name: data
              mountPath: {{ .Values.persistence.mountPath }}
          {{- end }}

      {{- if .Values.persistence.enabled }}
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: {{ include "my-go-app.fullname" . }}-data
      {{- end }}

      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}

      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}

      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
```

### 5.6 依赖管理

```yaml
# Chart.lock - 锁定依赖版本
generated: "2024-01-15T10:30:00Z"
digest: sha256:abc123...
dependencies:
  - name: postgresql
    repository: https://charts.bitnami.com/bitnami
    version: 12.12.10
  - name: redis
    repository: https://charts.bitnami.com/bitnami
    version: 17.15.0
```

```bash
# 依赖管理命令

# 添加依赖
helm dependency add mychart postgresql --repo https://charts.bitnami.com/bitnami --version 12.x.x

# 更新依赖
helm dependency update ./my-go-app

# 构建依赖（打包到charts目录）
helm dependency build ./my-go-app

# 列出依赖
helm dependency list ./my-go-app
```

### 5.7 Go应用Helm示例

```bash
# 创建Chart
helm create my-go-app

# 安装Chart
helm install my-app ./my-go-app

# 使用自定义values安装
helm install my-app ./my-go-app -f values-production.yaml

# 设置参数安装
helm install my-app ./my-go-app \
  --set image.tag=v2.0.0 \
  --set replicaCount=5 \
  --set config.logLevel=debug

# 升级Release
helm upgrade my-app ./my-go-app

# 回滚
helm rollback my-app 1

# 卸载
helm uninstall my-app

# 调试模板
helm template my-app ./my-go-app
helm template my-app ./my-go-app --debug

# 打包Chart
helm package ./my-go-app

# 推送到Chart仓库
helm push my-go-app-1.2.0.tgz oci://registry.example.com/charts
```

### 5.8 反例说明

```yaml
# ❌ 错误配置示例

# 1. 硬编码值
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 3  # 应该使用Values.replicaCount
  template:
    spec:
      containers:
        - name: app
          image: myapp:v1.0.0  # 应该使用模板

# 2. 缺少默认值
# values.yaml（缺少关键配置）
image:
  repository: myapp
  # 缺少tag和pullPolicy

# 3. 没有资源限制
# values.yaml
resources: {}  # 空配置

# 4. 模板语法错误
# templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.name }}  # 如果name未定义会报错
```

### 5.9 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 模板化 | 所有可配置项使用Values | 硬编码值 |
| 默认值 | 提供合理的默认值 | 缺少默认值 |
| 辅助模板 | 使用_helpers.tpl复用代码 | 重复代码 |
| 依赖管理 | 使用Chart.lock锁定版本 | 浮动版本 |
| 文档 | 提供README和NOTES.txt | 无文档 |
| 版本管理 | 语义化版本 | 随意版本号 |

---

## 6. Terraform基础设施

### 6.1 概念定义

**Terraform** 是HashiCorp开发的基础设施即代码（IaC）工具，允许开发者使用声明式配置语言（HCL）来定义和管理云基础设施。对于Go项目，Terraform可以管理K8s集群、数据库、缓存、网络等所有基础设施资源。

**核心概念：**

- **Provider**：云服务提供商插件（AWS、GCP、Azure、Kubernetes等）
- **Resource**：基础设施组件（VM、数据库、网络等）
- **Module**：可复用的资源配置单元
- **State**：Terraform管理的基础设施状态
- **Workspace**：隔离的环境（dev、staging、prod）

### 6.2 基础设施即代码流程

```
┌─────────────────────────────────────────────────────────────────┐
│                    Terraform 工作流程                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐  │
│  │  Write   │───▶│   Plan   │───▶│  Apply   │───▶│  Verify  │  │
│  │  Config  │    │  Changes │    │  Changes │    │  State   │  │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘  │
│       │               │               │               │        │
│       ▼               ▼               ▼               ▼        │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐  │
│  │ *.tf     │    │ Preview  │    │ Execute  │    │ tfstate  │  │
│  │ HCL代码  │    │ Changes  │    │ Actions  │    │ 状态文件 │  │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    状态管理                               │   │
│  │  Local State / Remote State (S3 + DynamoDB / Terraform Cloud) │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 6.3 Provider配置

```hcl
# providers.tf - Provider配置

terraform {
  required_version = ">= 1.5.0"

  required_providers {
    # AWS Provider
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }

    # Kubernetes Provider
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }

    # Helm Provider
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11"
    }

    # Docker Provider
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }

    # PostgreSQL Provider
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "~> 1.21"
    }

    # Random Provider
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }

  # 远程状态存储
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "go-app/terraform.tfstate"
    region         = "us-west-2"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}

# AWS Provider配置
provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

# Kubernetes Provider配置
provider "kubernetes" {
  host                   = module.eks.cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)

  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
  }
}

# Helm Provider配置
provider "helm" {
  kubernetes {
    host                   = module.eks.cluster_endpoint
    cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)

    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
    }
  }
}
```

### 6.4 资源管理

```hcl
# main.tf - 主资源配置

# ============================================
# 网络配置
# ============================================
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "${var.project_name}-${var.environment}"
  cidr = var.vpc_cidr

  azs             = var.availability_zones
  private_subnets = var.private_subnet_cidrs
  public_subnets  = var.public_subnet_cidrs

  enable_nat_gateway     = true
  single_nat_gateway     = var.environment != "production"
  enable_dns_hostnames   = true
  enable_dns_support     = true

  tags = {
    Name = "${var.project_name}-${var.environment}"
  }
}

# ============================================
# EKS集群
# ============================================
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "${var.project_name}-${var.environment}"
  cluster_version = "1.28"

  vpc_id                         = module.vpc.vpc_id
  subnet_ids                     = module.vpc.private_subnets
  control_plane_subnet_ids       = module.vpc.private_subnets

  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = true

  # 集群日志
  cluster_enabled_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  # EKS托管节点组
  eks_managed_node_groups = {
    general = {
      desired_size = var.node_desired_size
      min_size     = var.node_min_size
      max_size     = var.node_max_size

      instance_types = var.node_instance_types
      capacity_type  = var.environment == "production" ? "ON_DEMAND" : "SPOT"

      labels = {
        workload = "general"
      }

      tags = {
        Name = "${var.project_name}-${var.environment}-general"
      }
    }

    spot = {
      desired_size = 2
      min_size     = 0
      max_size     = 10

      instance_types = ["t3.medium", "t3a.medium"]
      capacity_type  = "SPOT"

      labels = {
        workload = "spot"
      }

      taints = [{
        key    = "spot"
        value  = "true"
        effect = "NO_SCHEDULE"
      }]
    }
  }

  # 集群访问条目
  access_entries = {
    admin = {
      principal_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/AdminRole"
      policy_associations = {
        admin = {
          policy_arn = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
          access_scope = {
            type = "cluster"
          }
        }
      }
    }
  }

  tags = {
    Name = "${var.project_name}-${var.environment}"
  }
}

# ============================================
# RDS PostgreSQL
# ============================================
module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 6.0"

  identifier = "${var.project_name}-${var.environment}"

  engine               = "postgres"
  engine_version       = "15.4"
  family               = "postgres15"
  major_engine_version = "15"
  instance_class       = var.db_instance_class

  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = var.db_max_allocated_storage
  storage_encrypted     = true

  db_name  = var.db_name
  username = var.db_username
  port     = 5432

  manage_master_user_password = true

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  # 安全组规则
  security_group_rules = {
    eks_ingress = {
      source_security_group_id = module.eks.node_security_group_id
      from_port                = 5432
      to_port                  = 5432
      ip_protocol              = "tcp"
      description              = "EKS nodes access"
    }
  }

  # 备份
  backup_retention_period = var.environment == "production" ? 30 : 7
  backup_window           = "03:00-04:00"
  maintenance_window      = "Mon:04:00-Mon:05:00"

  # 监控
  monitoring_interval    = 60
  monitoring_role_name   = "${var.project_name}-rds-monitoring"
  create_monitoring_role = true

  # 删除保护
  deletion_protection = var.environment == "production"

  tags = {
    Name = "${var.project_name}-${var.environment}"
  }
}

# ============================================
# ElastiCache Redis
# ============================================
module "elasticache" {
  source  = "terraform-aws-modules/elasticache/aws"
  version = "~> 1.0"

  cluster_id = "${var.project_name}-${var.environment}"

  engine               = "redis"
  engine_version       = "7.0"
  node_type            = var.redis_node_type
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379

  subnet_group_name  = "${var.project_name}-${var.environment}"
  subnet_group_names = module.vpc.private_subnets

  security_group_rules = {
    eks_ingress = {
      source_security_group_id = module.eks.node_security_group_id
      from_port                = 6379
      to_port                  = 6379
      ip_protocol              = "tcp"
    }
  }

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true

  tags = {
    Name = "${var.project_name}-${var.environment}"
  }
}

# ============================================
# Application Load Balancer
# ============================================
module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 9.0"

  name = "${var.project_name}-${var.environment}"

  load_balancer_type = "application"
  vpc_id             = module.vpc.vpc_id
  subnets            = module.vpc.public_subnets
  security_groups    = [aws_security_group.alb.id]

  # 访问日志
  access_logs = {
    bucket  = aws_s3_bucket.logs.id
    enabled = true
  }

  # 监听器和规则
  listeners = {
    https = {
      port            = 443
      protocol        = "HTTPS"
      certificate_arn = aws_acm_certificate.main.arn

      fixed_response = {
        content_type = "text/plain"
        message_body = "OK"
        status_code  = "200"
      }
    }
  }

  tags = {
    Name = "${var.project_name}-${var.environment}"
  }
}
```

### 6.5 状态管理

```hcl
# backend.tf - 远程状态配置

# S3 + DynamoDB后端配置
terraform {
  backend "s3" {
    bucket         = "mycompany-terraform-state"
    key            = "go-app/${var.environment}/terraform.tfstate"
    region         = "us-west-2"
    encrypt        = true
    kms_key_id     = "arn:aws:kms:us-west-2:123456789:key/terraform-state"
    dynamodb_table = "terraform-locks"

    # 工作区前缀
    workspace_key_prefix = "workspaces"
  }
}

# 或者使用Terraform Cloud
terraform {
  cloud {
    organization = "mycompany"

    workspaces {
      name = "go-app-${var.environment}"
    }
  }
}

# 状态锁定DynamoDB表
resource "aws_dynamodb_table" "terraform_locks" {
  name         = "terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name = "Terraform State Lock Table"
  }
}

# 状态存储S3桶
resource "aws_s3_bucket" "terraform_state" {
  bucket = "mycompany-terraform-state"

  tags = {
    Name = "Terraform State Bucket"
  }
}

resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.terraform_state.arn
    }
  }
}
```

### 6.6 Go项目Terraform示例

```hcl
# variables.tf - 变量定义

variable "project_name" {
  description = "项目名称"
  type        = string
  default     = "go-app"
}

variable "environment" {
  description = "环境名称"
  type        = string

  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "环境必须是 development, staging, 或 production。"
  }
}

variable "aws_region" {
  description = "AWS区域"
  type        = string
  default     = "us-west-2"
}

variable "vpc_cidr" {
  description = "VPC CIDR块"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "可用区列表"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

variable "node_instance_types" {
  description = "EKS节点实例类型"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "db_instance_class" {
  description = "RDS实例类型"
  type        = string
  default     = "db.t3.micro"
}

# locals.tf - 本地变量
locals {
  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
    CreatedAt   = timestamp()
  }

  cluster_name = "${var.project_name}-${var.environment}"
}

# outputs.tf - 输出值
output "cluster_endpoint" {
  description = "EKS集群端点"
  value       = module.eks.cluster_endpoint
}

output "cluster_name" {
  description = "EKS集群名称"
  value       = module.eks.cluster_name
}

output "database_endpoint" {
  description = "RDS数据库端点"
  value       = module.rds.db_instance_endpoint
  sensitive   = true
}

output "redis_endpoint" {
  description = "Redis端点"
  value       = module.elasticache.cluster_endpoint
  sensitive   = true
}
```

```bash
# Terraform工作流命令

# 初始化
terraform init

# 格式化代码
terraform fmt -recursive

# 验证配置
terraform validate

# 查看执行计划
terraform plan

# 应用配置
terraform apply

# 销毁资源
terraform destroy

# 工作区管理
terraform workspace new production
terraform workspace select production
terraform workspace list

# 导入现有资源
terraform import aws_instance.my_instance i-1234567890abcdef0

# 查看状态
terraform state list
terraform state show aws_instance.my_instance

# 刷新状态
terraform refresh
```

### 6.7 反例说明

```hcl
# ❌ 错误配置示例

# 1. 硬编码值
resource "aws_instance" "server" {
  ami           = "ami-12345678"  # 应该使用变量
  instance_type = "t2.micro"       # 应该使用变量
}

# 2. 没有状态管理
# 没有配置backend，使用本地状态
# 风险：状态丢失、团队协作困难

# 3. 没有资源依赖
resource "aws_db_instance" "db" {
  # 没有depends_on，可能在VPC创建前创建
}

# 4. 敏感信息明文存储
variable "db_password" {
  default = "mypassword123"  # 安全风险
}

# 5. 没有标签
resource "aws_s3_bucket" "logs" {
  bucket = "my-logs"
  # 缺少标签，难以管理
}
```

### 6.8 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 变量管理 | 使用variables.tf定义变量 | 硬编码值 |
| 状态管理 | 使用远程状态存储 | 本地状态 |
| 模块化 | 使用模块组织代码 | 单一大文件 |
| 敏感数据 | 使用Vault或KMS加密 | 明文存储 |
| 版本锁定 | 锁定Provider版本 | 浮动版本 |
| 标签管理 | 统一标签策略 | 无标签 |
| 工作区 | 使用workspace隔离环境 | 单工作区 |

---

## 7. 持续集成流程

### 7.1 概念定义

**持续集成（Continuous Integration, CI）** 是一种软件开发实践，开发人员频繁地将代码变更合并到主干分支，每次合并都通过自动化的构建和测试流程进行验证。对于Go项目，CI流程确保代码质量、功能正确性和安全性。

**核心目标：**

- **快速反馈**：尽早发现问题
- **自动化**：减少人工干预
- **质量保证**：确保代码符合标准
- **可追溯**：记录每次构建的详细信息

### 7.2 流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                    持续集成流程                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│  │  Push   │──▶│  Lint   │──▶│  Build  │──▶│  Test   │         │
│  │  Code   │   │  Check  │   │         │   │         │         │
│  └─────────┘   └─────────┘   └─────────┘   └─────────┘         │
│       │             │             │             │              │
│       │             ▼             ▼             ▼              │
│       │       ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│       │       │golangci │   │ Compile │   │  Unit   │         │
│       │       │ -lint   │   │ Binary  │   │  Test   │         │
│       │       └─────────┘   └─────────┘   └─────────┘         │
│       │                                         │              │
│       │                                         ▼              │
│       │                                   ┌─────────┐         │
│       │                                   │Coverage │         │
│       │                                   │ Report  │         │
│       │                                   └─────────┘         │
│       │                                         │              │
│       ▼                                         ▼              │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│  │ Security│──▶│  SAST   │──▶│  SCA    │──▶│  Image  │         │
│  │  Scan   │   │         │   │         │   │  Scan   │         │
│  └─────────┘   └─────────┘   └─────────┘   └─────────┘         │
│                                                                 │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐                       │
│  │ Artifact│──▶│  Push   │──▶│ Notify  │                       │
│  │  Build  │   │ Registry│   │         │                       │
│  └─────────┘   └─────────┘   └─────────┘                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 7.3 代码检查（Lint）

```yaml
# .golangci.yml - golangci-lint配置
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true
  skip-dirs:
    - vendor
    - third_party
    - testdata
  skip-files:
    - ".*\\.pb\\.go$"
    - ".*_test\\.go$"

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true

  govet:
    check-shadowing: true
    enable-all: true

  gocyclo:
    min-complexity: 15

  maligned:
    suggest-new: true

  dupl:
    threshold: 100

  goconst:
    min-len: 3
    min-occurrences: 3

  misspell:
    locale: US

  lll:
    line-length: 120

  goimports:
    local-prefixes: github.com/mycompany/myapp

  gocritic:
    enabled-tags:
      - performance
      - style
      - experimental
    disabled-checks:
      - wrapperFunc
      - dupImport

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

  disable:
    - maligned  # deprecated
    - prealloc  # can be noisy

issues:
  exclude-rules:
    # Exclude some linters from running on tests files
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
        - lll

    # Exclude known linter issues
    - text: "weak cryptographic primitive"
      linters:
        - gosec

    # Exclude shadow checking in test files
    - path: _test\.go
      text: "shadow"
      linters:
        - govet

  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

```yaml
# GitHub Actions Lint Job
lint:
  name: Lint
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'
        cache: true

    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: v1.55
        args: --timeout=5m --config=.golangci.yml
        skip-cache: false

    - name: Run go vet
      run: go vet ./...

    - name: Check formatting
      run: |
        if [ "$(gofmt -l . | wc -l)" -gt 0 ]; then
          echo "The following files need formatting:"
          gofmt -l .
          exit 1
        fi

    - name: Run goimports
      run: |
        go install golang.org/x/tools/cmd/goimports@latest
        if [ "$(goimports -l . | wc -l)" -gt 0 ]; then
          echo "The following files need import formatting:"
          goimports -l .
          exit 1
        fi
```

### 7.4 单元测试

```yaml
# GitHub Actions Test Job
test:
  name: Test
  runs-on: ubuntu-latest
  services:
    postgres:
      image: postgres:15-alpine
      env:
        POSTGRES_USER: test
        POSTGRES_PASSWORD: test
        POSTGRES_DB: testdb
      ports:
        - 5432:5432
      options: >-
        --health-cmd pg_isready
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5

    redis:
      image: redis:7-alpine
      ports:
        - 6379:6379
      options: >-
        --health-cmd "redis-cli ping"
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5

  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'
        cache: true

    - name: Download dependencies
      run: go mod download

    - name: Run unit tests
      run: go test -v -race -short ./...
      env:
        TEST_DATABASE_URL: postgres://test:test@localhost:5432/testdb?sslmode=disable
        TEST_REDIS_URL: localhost:6379

    - name: Run integration tests
      run: go test -v -race -tags=integration ./tests/integration/...
      env:
        TEST_DATABASE_URL: postgres://test:test@localhost:5432/testdb?sslmode=disable
        TEST_REDIS_URL: localhost:6379
```

```go
// 测试示例 - internal/service/user_test.go
package service

import (
 "context"
 "testing"

 "github.com/stretchr/testify/assert"
 "github.com/stretchr/testify/mock"
)

// Mock repository
type mockUserRepository struct {
 mock.Mock
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
 args := m.Called(ctx, id)
 if args.Get(0) == nil {
  return nil, args.Error(1)
 }
 return args.Get(0).(*User), args.Error(1)
}

func (m *mockUserRepository) Create(ctx context.Context, user *User) error {
 args := m.Called(ctx, user)
 return args.Error(0)
}

func TestUserService_GetUser(t *testing.T) {
 tests := []struct {
  name     string
  userID   int64
  mockFunc func(*mockUserRepository)
  want     *User
  wantErr  bool
 }{
  {
   name:   "success",
   userID: 1,
   mockFunc: func(m *mockUserRepository) {
    m.On("GetByID", mock.Anything, int64(1)).Return(&User{
     ID:    1,
     Name:  "John",
     Email: "john@example.com",
    }, nil)
   },
   want: &User{
    ID:    1,
    Name:  "John",
    Email: "john@example.com",
   },
   wantErr: false,
  },
  {
   name:   "not found",
   userID: 999,
   mockFunc: func(m *mockUserRepository) {
    m.On("GetByID", mock.Anything, int64(999)).Return(nil, ErrUserNotFound)
   },
   want:    nil,
   wantErr: true,
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   mockRepo := new(mockUserRepository)
   tt.mockFunc(mockRepo)

   svc := NewUserService(mockRepo)
   got, err := svc.GetUser(context.Background(), tt.userID)

   if tt.wantErr {
    assert.Error(t, err)
   } else {
    assert.NoError(t, err)
    assert.Equal(t, tt.want, got)
   }

   mockRepo.AssertExpectations(t)
  })
 }
}

// Benchmark测试
func BenchmarkUserService_GetUser(b *testing.B) {
 mockRepo := new(mockUserRepository)
 mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&User{
  ID:    1,
  Name:  "John",
  Email: "john@example.com",
 }, nil)

 svc := NewUserService(mockRepo)
 ctx := context.Background()

 b.ResetTimer()
 for i := 0; i < b.N; i++ {
  _, _ = svc.GetUser(ctx, 1)
 }
}
```

### 7.5 代码覆盖率

```yaml
# GitHub Actions Coverage Job
coverage:
  name: Coverage
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'
        cache: true

    - name: Run tests with coverage
      run: |
        go test -race -coverprofile=coverage.out -covermode=atomic ./...
        go tool cover -html=coverage.out -o coverage.html

    - name: Generate coverage report
      run: |
        go install github.com/axw/gocov/gocov@latest
        go install github.com/AlekSi/gocov-xml@latest
        gocov convert coverage.out | gocov-xml > coverage.xml

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
        flags: unittests
        name: codecov-umbrella
        fail_ci_if_error: true

    - name: Check coverage threshold
      run: |
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "Total coverage: $COVERAGE%"
        if (( $(echo "$COVERAGE < 80" | bc -l) )); then
          echo "Coverage $COVERAGE% is below threshold 80%"
          exit 1
        fi

    - name: Upload coverage artifacts
      uses: actions/upload-artifact@v4
      with:
        name: coverage-report
        path: |
          coverage.out
          coverage.html
          coverage.xml
```

### 7.6 安全扫描

```yaml
# GitHub Actions Security Scan Job
security:
  name: Security Scan
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    # SAST - 静态应用安全测试
    - name: Run Gosec Security Scanner
      uses: securego/gosec@master
      with:
        args: '-fmt sarif -out gosec-report.sarif ./...'

    - name: Upload Gosec report
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: gosec-report.sarif

    # 依赖漏洞扫描
    - name: Run Nancy (依赖漏洞扫描)
      run: |
        go install github.com/sonatypecommunity/nancy@latest
        go list -json -deps ./... | nancy sleuth

    # SCA - 软件成分分析
    - name: Run Snyk to check for vulnerabilities
      uses: snyk/actions/golang@master
      env:
        SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
      with:
        args: --severity-threshold=high

    # 容器镜像扫描
    - name: Build Docker image
      run: docker build -t myapp:test .

    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: 'myapp:test'
        format: 'sarif'
        output: 'trivy-report.sarif'

    - name: Upload Trivy report
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: trivy-report.sarif

    # 密钥扫描
    - name: Secret Detection
      uses: trufflesecurity/trufflehog@main
      with:
        path: ./
        base: main
        head: HEAD
```

### 7.7 完整CI流程示例

```yaml
# .github/workflows/complete-ci.yml
name: Complete CI Pipeline

on:
  push:
    branches: [main, develop]
    paths-ignore:
      - '**.md'
      - 'docs/**'
  pull_request:
    branches: [main]

env:
  GO_VERSION: '1.21'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # ============ 代码质量检查 ============
  lint:
    name: Code Quality
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: v1.55
          args: --timeout=5m

      - name: go vet
        run: go vet ./...

      - name: Check formatting
        run: test -z "$(gofmt -l .)"

  # ============ 安全扫描 ============
  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run Gosec
        uses: securego/gosec@master
        with:
          args: '-fmt sarif -out gosec.sarif ./...'

      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: gosec.sarif
        if: always()

  # ============ 单元测试 ============
  unit-test:
    name: Unit Tests
    runs-on: ubuntu-latest
    needs: [lint]
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Run tests
        run: go test -v -race -short ./...

  # ============ 集成测试 ============
  integration-test:
    name: Integration Tests
    runs-on: ubuntu-latest
    needs: [unit-test]
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: testdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Run integration tests
        run: go test -v -race -tags=integration ./tests/integration/...
        env:
          TEST_DATABASE_URL: postgres://test:test@localhost:5432/testdb?sslmode=disable

  # ============ 代码覆盖率 ============
  coverage:
    name: Coverage
    runs-on: ubuntu-latest
    needs: [unit-test]
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Generate coverage
        run: go test -race -coverprofile=coverage.out ./...

      - name: Upload to Codecov
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
          fail_ci_if_error: true

  # ============ 构建 ============
  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [lint, unit-test]
    strategy:
      matrix:
        os: [linux, darwin, windows]
        arch: [amd64, arm64]
        exclude:
          - os: windows
            arch: arm64
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Build binary
        run: |
          OUTPUT="bin/myapp-${{ matrix.os }}-${{ matrix.arch }}"
          if [ "${{ matrix.os }}" = "windows" ]; then
            OUTPUT="${OUTPUT}.exe"
          fi
          GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} \
            go build -ldflags="-s -w" -o ${OUTPUT} ./cmd/myapp

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: myapp-${{ matrix.os }}-${{ matrix.arch }}
          path: bin/

  # ============ Docker构建 ============
  docker:
    name: Docker Build
    runs-on: ubuntu-latest
    needs: [security, unit-test]
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
            type=sha,prefix=,suffix=,format=short

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          platforms: linux/amd64,linux/arm64
```

### 7.8 反例说明

```yaml
# ❌ 错误配置示例
name: Bad CI

on: push

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # ❌ 没有lint检查
      # ❌ 没有安全扫描

      - uses: actions/setup-go@v5
      # ❌ 没有指定Go版本

      - run: go test ./...
        # ❌ 没有-race检测
        # ❌ 没有覆盖率检查

      - run: go build ./...

      - run: docker build -t myapp .
        # ❌ 没有镜像扫描

      - run: docker push myapp
        # ❌ 没有条件判断
```

### 7.9 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| Lint检查 | 使用golangci-lint | 忽略代码风格 |
| 测试 | 单元测试+集成测试 | 仅单元测试 |
| 覆盖率 | 设置80%阈值 | 无覆盖率要求 |
| 安全扫描 | SAST+SCA+镜像扫描 | 忽略安全 |
| 并行执行 | lint/test并行 | 串行执行 |
| 缓存 | 启用Go模块缓存 | 每次重新下载 |

---

## 8. 持续部署流程

### 8.1 概念定义

**持续部署（Continuous Deployment, CD）** 是一种软件发布实践，将通过CI验证的代码自动部署到生产环境。对于Go项目，CD流程确保快速、可靠、可回滚的发布能力。

**部署策略：**

- **滚动部署（Rolling）**：逐步替换旧版本Pod
- **蓝绿部署（Blue-Green）**：同时维护两个环境，一键切换
- **金丝雀发布（Canary）**：小流量验证新版本
- **功能开关（Feature Flags）**：代码级发布控制

### 8.2 流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                    持续部署流程                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│  │  Build  │──▶│  Push   │──▶│ Deploy  │──▶│ Verify  │         │
│  │  Image  │   │ Registry│   │  Staging│   │  Smoke  │         │
│  └─────────┘   └─────────┘   └─────────┘   └─────────┘         │
│       │                              │             │           │
│       │                              ▼             ▼           │
│       │                       ┌─────────┐   ┌─────────┐       │
│       │                       │  Helm   │   │  Test   │       │
│       │                       │ Upgrade │   │  Suite  │       │
│       │                       └─────────┘   └─────────┘       │
│       │                              │             │           │
│       ▼                              ▼             ▼           │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│  │Promote  │──▶│ Deploy  │──▶│ Canary  │──▶│  Full   │         │
│  │  Image  │   │Production│   │  Test   │   │ Rollout │         │
│  └─────────┘   └─────────┘   └─────────┘   └─────────┘         │
│       │                              │             │           │
│       │                              ▼             ▼           │
│       │                       ┌─────────┐   ┌─────────┐       │
│       │                       │  5%     │   │ 100%    │       │
│       │                       │ Traffic │   │ Traffic │       │
│       │                       └─────────┘   └─────────┘       │
│       │                                                         │
│  ┌─────────┐   ┌─────────┐                                     │
│  │ Monitor │──▶│ Rollback│  （失败时回滚）                      │
│  │  Alert  │   │         │                                     │
│  └─────────┘   └─────────┘                                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 8.3 环境管理

```yaml
# GitHub Actions - 多环境部署
name: Deploy to Environments

on:
  push:
    branches: [main]
  release:
    types: [published]

jobs:
  # ============ 部署到Staging ============
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name staging-cluster

      - name: Deploy with Helm
        run: |
          helm upgrade --install myapp ./helm-chart \
            --namespace staging \
            --set image.tag=${{ github.sha }} \
            --set replicaCount=2 \
            --wait \
            --timeout 5m

      - name: Run smoke tests
        run: |
          curl -f https://staging.example.com/health || exit 1

  # ============ 部署到Production（需要审批） ============
  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: deploy-staging
    environment:
      name: production
      url: https://production.example.com
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name production-cluster

      - name: Deploy with Helm
        run: |
          helm upgrade --install myapp ./helm-chart \
            --namespace production \
            --set image.tag=${{ github.event.release.tag_name }} \
            --set replicaCount=5 \
            --wait \
            --timeout 10m

      - name: Verify deployment
        run: |
          kubectl rollout status deployment/myapp -n production
          curl -f https://production.example.com/health || exit 1
```

### 8.4 蓝绿部署

```yaml
# templates/deployment-bluegreen.yaml
# Helm Chart中的蓝绿部署配置

# 蓝色版本（当前生产）
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "myapp.fullname" . }}-blue
  labels:
    app: {{ include "myapp.name" . }}
    version: blue
    {{- include "myapp.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app: {{ include "myapp.name" . }}
      version: blue
  template:
    metadata:
      labels:
        app: {{ include "myapp.name" . }}
        version: blue
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.blue.tag }}"
          ports:
            - containerPort: 8080

---
# 绿色版本（新版本）
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "myapp.fullname" . }}-green
  labels:
    app: {{ include "myapp.name" . }}
    version: green
    {{- include "myapp.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app: {{ include "myapp.name" . }}
      version: green
  template:
    metadata:
      labels:
        app: {{ include "myapp.name" . }}
        version: green
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.green.tag }}"
          ports:
            - containerPort: 8080

---
# Service - 通过selector切换流量
apiVersion: v1
kind: Service
metadata:
  name: {{ include "myapp.fullname" . }}
  labels:
    {{- include "myapp.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: 8080
  selector:
    app: {{ include "myapp.name" . }}
    version: {{ .Values.activeVersion }}  # blue 或 green
```

```bash
#!/bin/bash
# blue-green-deploy.sh - 蓝绿部署脚本

set -e

APP_NAME="myapp"
NAMESPACE="production"
NEW_VERSION=$1

if [ -z "$NEW_VERSION" ]; then
    echo "Usage: $0 <new-version>"
    exit 1
fi

# 获取当前活跃版本
CURRENT_VERSION=$(kubectl get svc $APP_NAME -n $NAMESPACE -o jsonpath='{.spec.selector.version}')
echo "Current active version: $CURRENT_VERSION"

# 确定新版本的颜色
if [ "$CURRENT_VERSION" = "blue" ]; then
    NEW_COLOR="green"
else
    NEW_COLOR="blue"
fi
echo "Deploying to: $NEW_COLOR"

# 部署新版本
helm upgrade --install $APP_NAME ./helm-chart \
    --namespace $NAMESPACE \
    --set ${NEW_COLOR}.tag=$NEW_VERSION \
    --set activeVersion=$CURRENT_VERSION \
    --wait \
    --timeout 10m

# 等待新版本就绪
echo "Waiting for $NEW_COLOR deployment to be ready..."
kubectl rollout status deployment/$APP_NAME-$NEW_COLOR -n $NAMESPACE

# 健康检查
echo "Running health checks..."
HEALTH_URL="http://$APP_NAME-$NEW_COLOR.$NAMESPACE.svc.cluster.local/health"
for i in {1..5}; do
    if curl -sf $HEALTH_URL > /dev/null; then
        echo "Health check passed"
        break
    fi
    echo "Health check failed, retrying..."
    sleep 5
done

# 切换流量
echo "Switching traffic to $NEW_COLOR..."
kubectl patch svc $APP_NAME -n $NAMESPACE -p "{\"spec\":{\"selector\":{\"version\":\"$NEW_COLOR\"}}}"

# 验证
echo "Verifying traffic switch..."
sleep 10
NEW_ACTIVE=$(kubectl get svc $APP_NAME -n $NAMESPACE -o jsonpath='{.spec.selector.version}')
if [ "$NEW_ACTIVE" = "$NEW_COLOR" ]; then
    echo "Traffic successfully switched to $NEW_COLOR"
else
    echo "Traffic switch failed!"
    exit 1
fi

# 保留旧版本（可选：一段时间后删除）
echo "Deployment complete. Old version ($CURRENT_VERSION) is still running for rollback."
echo "To rollback: kubectl patch svc $APP_NAME -n $NAMESPACE -p '{\"spec\":{\"selector\":{\"version\":\"$CURRENT_VERSION\"}}}'"
```

### 8.5 金丝雀发布

```yaml
# templates/canary.yaml - 金丝雀发布配置

# 稳定版本（主流量）
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "myapp.fullname" . }}-stable
spec:
  replicas: {{ .Values.stable.replicas }}
  selector:
    matchLabels:
      app: {{ include "myapp.name" . }}
      canary: "false"
  template:
    metadata:
      labels:
        app: {{ include "myapp.name" . }}
        canary: "false"
    spec:
      containers:
        - name: app
          image: "{{ .Values.image.repository }}:{{ .Values.stable.tag }}"

---
# 金丝雀版本（小流量）
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "myapp.fullname" . }}-canary
spec:
  replicas: {{ .Values.canary.replicas }}
  selector:
    matchLabels:
      app: {{ include "myapp.name" . }}
      canary: "true"
  template:
    metadata:
      labels:
        app: {{ include "myapp.name" . }}
        canary: "true"
    spec:
      containers:
        - name: app
          image: "{{ .Values.image.repository }}:{{ .Values.canary.tag }}"

---
# Service - 同时选择两个版本
apiVersion: v1
kind: Service
metadata:
  name: {{ include "myapp.fullname" . }}
spec:
  selector:
    app: {{ include "myapp.name" . }}
  ports:
    - port: 80
      targetPort: 8080

---
# Ingress - 使用权重分流（需要NGINX Ingress）
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "myapp.fullname" . }}-canary
  annotations:
    nginx.ingress.kubernetes.io/canary: "true"
    nginx.ingress.kubernetes.io/canary-weight: "{{ .Values.canary.weight }}"
spec:
  rules:
    - host: {{ .Values.ingress.host }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: {{ include "myapp.fullname" . }}-canary
                port:
                  number: 80
```

```bash
#!/bin/bash
# canary-deploy.sh - 金丝雀发布脚本

set -e

APP_NAME="myapp"
NAMESPACE="production"
NEW_VERSION=$1
INITIAL_WEIGHT=${2:-10}  # 默认10%流量
STEP=${3:-10}            # 每次增加10%
INTERVAL=${4:-300}       # 每5分钟增加一次

if [ -z "$NEW_VERSION" ]; then
    echo "Usage: $0 <new-version> [initial-weight] [step] [interval]"
    exit 1
fi

echo "Starting canary deployment of $NEW_VERSION"
echo "Initial weight: $INITIAL_WEIGHT%, Step: $STEP%, Interval: ${INTERVAL}s"

# 部署金丝雀版本
echo "Deploying canary version..."
helm upgrade --install $APP_NAME ./helm-chart \
    --namespace $NAMESPACE \
    --set canary.tag=$NEW_VERSION \
    --set canary.replicas=1 \
    --set canary.weight=$INITIAL_WEIGHT \
    --wait

# 渐进式增加流量
CURRENT_WEIGHT=$INITIAL_WEIGHT
while [ $CURRENT_WEIGHT -lt 100 ]; do
    echo "Current canary weight: $CURRENT_WEIGHT%"

    # 检查指标
    echo "Checking metrics..."
    ERROR_RATE=$(curl -s "http://prometheus/api/v1/query?query=rate(http_requests_total{canary=\"true\",status=~\"5..\"}[5m])" | jq -r '.data.result[0].value[1] // 0')
    LATENCY=$(curl -s "http://prometheus/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket{canary=\"true\"}[5m]))" | jq -r '.data.result[0].value[1] // 0')

    echo "Error rate: $ERROR_RATE, P95 latency: $LATENCY"

    # 检查是否超过阈值
    if (( $(echo "$ERROR_RATE > 0.01" | bc -l) )); then
        echo "Error rate too high! Rolling back..."
        helm upgrade --install $APP_NAME ./helm-chart \
            --namespace $NAMESPACE \
            --set canary.weight=0 \
            --wait
        exit 1
    fi

    # 增加流量
    CURRENT_WEIGHT=$((CURRENT_WEIGHT + STEP))
    if [ $CURRENT_WEIGHT -gt 100 ]; then
        CURRENT_WEIGHT=100
    fi

    echo "Increasing canary weight to $CURRENT_WEIGHT%"
    helm upgrade --install $APP_NAME ./helm-chart \
        --namespace $NAMESPACE \
        --set canary.weight=$CURRENT_WEIGHT \
        --wait

    if [ $CURRENT_WEIGHT -lt 100 ]; then
        echo "Waiting ${INTERVAL}s before next increase..."
        sleep $INTERVAL
    fi
done

echo "Canary deployment complete! 100% traffic on new version."
```

### 8.6 功能开关

```go
// internal/feature/flags.go - 功能开关实现
package feature

import (
 "context"
 "sync"

 "github.com/redis/go-redis/v9"
)

// Flag 功能开关
type Flag string

const (
 FlagNewUI      Flag = "new_ui"
 FlagBetaAPI    Flag = "beta_api"
 FlagDarkMode   Flag = "dark_mode"
 FlagNewPayment Flag = "new_payment"
)

// Manager 功能开关管理器
type Manager struct {
 redis    *redis.Client
 defaults map[Flag]bool
 mu       sync.RWMutex
 local    map[Flag]bool
}

// NewManager 创建功能开关管理器
func NewManager(redis *redis.Client) *Manager {
 return &Manager{
  redis: redis,
  defaults: map[Flag]bool{
   FlagNewUI:      false,
   FlagBetaAPI:    false,
   FlagDarkMode:   true,
   FlagNewPayment: false,
  },
  local: make(map[Flag]bool),
 }
}

// IsEnabled 检查功能是否启用
func (m *Manager) IsEnabled(ctx context.Context, flag Flag) bool {
 // 1. 检查本地缓存
 m.mu.RLock()
 if val, ok := m.local[flag]; ok {
  m.mu.RUnlock()
  return val
 }
 m.mu.RUnlock()

 // 2. 从Redis获取
 val, err := m.redis.Get(ctx, string(flag)).Bool()
 if err == nil {
  m.mu.Lock()
  m.local[flag] = val
  m.mu.Unlock()
  return val
 }

 // 3. 返回默认值
 return m.defaults[flag]
}

// Enable 启用功能
func (m *Manager) Enable(ctx context.Context, flag Flag) error {
 if err := m.redis.Set(ctx, string(flag), true, 0).Err(); err != nil {
  return err
 }
 m.mu.Lock()
 m.local[flag] = true
 m.mu.Unlock()
 return nil
}

// Disable 禁用功能
func (m *Manager) Disable(ctx context.Context, flag Flag) error {
 if err := m.redis.Set(ctx, string(flag), false, 0).Err(); err != nil {
  return err
 }
 m.mu.Lock()
 m.local[flag] = false
 m.mu.Unlock()
 return nil
}

// IsEnabledForUser 为特定用户检查功能
func (m *Manager) IsEnabledForUser(ctx context.Context, flag Flag, userID int64) bool {
 // 获取全局开关状态
 if !m.IsEnabled(ctx, flag) {
  return false
 }

 // 检查用户是否在白名单中
 isInWhitelist, _ := m.redis.SIsMember(ctx, string(flag)+":whitelist", userID).Result()
 if isInWhitelist {
  return true
 }

 // 基于用户ID的百分比灰度
 percentage, _ := m.redis.Get(ctx, string(flag)+":percentage").Int()
 if percentage > 0 {
  return userID%100 < int64(percentage)
 }

 return false
}
```

```go
// 使用示例 - internal/handler/user.go
package handler

import (
 "net/http"

 "github.com/gin-gonic/gin"
 "myapp/internal/feature"
)

type UserHandler struct {
 featureManager *feature.Manager
}

func (h *UserHandler) GetUser(c *gin.Context) {
 userID := c.GetInt64("user_id")

 // 检查新UI功能开关
 if h.featureManager.IsEnabledForUser(c.Request.Context(), feature.FlagNewUI, userID) {
  // 返回新UI数据
  c.JSON(http.StatusOK, gin.H{
   "ui_version": "v2",
   "data":       newUIUserData(),
  })
  return
 }

 // 返回旧UI数据
 c.JSON(http.StatusOK, gin.H{
  "ui_version": "v1",
  "data":       oldUIUserData(),
 })
}
```

### 8.7 回滚策略

```yaml
# GitHub Actions - 自动回滚
name: Rollback

on:
  workflow_dispatch:
    inputs:
      environment:
        description: 'Environment to rollback'
        required: true
        default: 'staging'
        type: choice
        options:
          - staging
          - production
      revision:
        description: 'Revision to rollback to (empty for previous)'
        required: false
        type: string

jobs:
  rollback:
    name: Rollback ${{ github.event.inputs.environment }}
    runs-on: ubuntu-latest
    environment:
      name: ${{ github.event.inputs.environment }}
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name ${{ github.event.inputs.environment }}-cluster

      - name: Rollback deployment
        run: |
          if [ -n "${{ github.event.inputs.revision }}" ]; then
            # 回滚到指定版本
            helm rollback myapp ${{ github.event.inputs.revision }} -n ${{ github.event.inputs.environment }}
          else
            # 回滚到上一版本
            helm rollback myapp 0 -n ${{ github.event.inputs.environment }}
          fi

      - name: Verify rollback
        run: |
          kubectl rollout status deployment/myapp -n ${{ github.event.inputs.environment }}
          kubectl get deployment myapp -n ${{ github.event.inputs.environment }} -o yaml | grep image:

      - name: Notify rollback
        uses: slackapi/slack-github-action@v1
        with:
          payload: |
            {
              "text": "🔄 Rollback completed for ${{ github.event.inputs.environment }}"
            }
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
```

```bash
#!/bin/bash
# rollback.sh - 回滚脚本

set -e

APP_NAME="myapp"
NAMESPACE=${1:-production}
REVISION=${2:-0}  # 0表示上一版本

echo "Rolling back $APP_NAME in namespace $NAMESPACE to revision $REVISION"

# 查看历史版本
echo "Deployment history:"
helm history $APP_NAME -n $NAMESPACE

# 执行回滚
echo "Executing rollback..."
helm rollback $APP_NAME $REVISION -n $NAMESPACE

# 等待回滚完成
echo "Waiting for rollback to complete..."
kubectl rollout status deployment/$APP_NAME -n $NAMESPACE

# 验证
echo "Current deployment:"
kubectl get deployment $APP_NAME -n $NAMESPACE -o wide

echo "Rollback completed successfully!"
```

### 8.8 完整CD流程示例

```yaml
# .github/workflows/cd.yml
name: Continuous Deployment

on:
  push:
    branches: [main]
    tags: ['v*']

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    name: Build and Push Image
    runs-on: ubuntu-latest
    outputs:
      image_tag: ${{ steps.meta.outputs.tags }}
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Registry
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
            type=semver,pattern={{version}}
            type=sha

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}

  deploy-staging:
    name: Deploy to Staging
    needs: build-and-push
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Deploy
        run: |
          aws eks update-kubeconfig --name staging
          helm upgrade --install myapp ./helm \
            --namespace staging \
            --set image.tag=${{ github.sha }} \
            --wait

      - name: Smoke test
        run: curl -sf https://staging.example.com/health

  deploy-production:
    name: Deploy to Production
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://production.example.com
    if: startsWith(github.ref, 'refs/tags/v')
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Deploy with canary
        run: |
          aws eks update-kubeconfig --name production
          ./scripts/canary-deploy.sh ${{ github.ref_name }} 10 10 300

      - name: Verify
        run: |
          kubectl rollout status deployment/myapp -n production
          curl -sf https://production.example.com/health
```

### 8.9 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 环境隔离 | 使用GitHub Environments | 无环境管理 |
| 部署策略 | 蓝绿/金丝雀发布 | 直接全量发布 |
| 功能开关 | 使用Feature Flags | 代码分支发布 |
| 回滚能力 | 保留历史版本 | 无法回滚 |
| 验证测试 | 部署后smoke test | 无验证 |
| 通知机制 | 部署结果通知 | 静默部署 |

---

## 9. 制品管理

### 9.1 概念定义

**制品管理（Artifact Management）** 是CI/CD流程中存储、版本控制和分发构建产物的关键环节。对于Go项目，制品包括编译后的二进制文件、Docker镜像、Helm Chart等。

**核心目标：**

- **可追溯性**：每个制品都能追溯到源代码版本
- **安全性**：制品签名验证，防止篡改
- **可用性**：快速获取历史版本
- **生命周期管理**：自动清理过期制品

### 9.2 镜像仓库

```yaml
# GitHub Actions - 多镜像仓库推送
name: Push to Registries

on:
  push:
    tags: ['v*']

env:
  IMAGE_NAME: go-app

jobs:
  push:
    runs-on: ubuntu-latest
    permissions:
      packages: write
      contents: read
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      # GitHub Container Registry
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Docker Hub
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      # AWS ECR
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2

      # 构建并推送
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository_owner }}/${{ env.IMAGE_NAME }}:${{ github.ref_name }}
            ghcr.io/${{ github.repository_owner }}/${{ env.IMAGE_NAME }}:latest
            ${{ secrets.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:${{ github.ref_name }}
            ${{ secrets.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:latest
            ${{ steps.login-ecr.outputs.registry }}/${{ env.IMAGE_NAME }}:${{ github.ref_name }}
            ${{ steps.login-ecr.outputs.registry }}/${{ env.IMAGE_NAME }}:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

### 9.3 二进制仓库

```yaml
# GitHub Actions - 发布二进制文件到GitHub Releases
name: Release Binaries

on:
  release:
    types: [created]

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
          - goos: linux
            goarch: arm64
          - goos: darwin
            goarch: amd64
          - goos: darwin
            goarch: arm64
          - goos: windows
            goarch: amd64
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Build binary
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          OUTPUT="myapp-${{ matrix.goos }}-${{ matrix.goarch }}"
          if [ "$GOOS" = "windows" ]; then
            OUTPUT="${OUTPUT}.exe"
          fi
          go build -ldflags="-s -w -X main.version=${{ github.ref_name }}" -o ${OUTPUT} ./cmd/myapp

      - name: Create checksum
        run: |
          sha256sum myapp-${{ matrix.goos }}-${{ matrix.goarch }}* > checksums.txt

      - name: Upload to Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            myapp-${{ matrix.goos }}-${{ matrix.goarch }}*
            checksums.txt
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 9.4 版本管理

```yaml
# GitHub Actions - 语义化版本管理
name: Semantic Release

on:
  push:
    branches: [main]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install semantic-release
        run: |
          npm install -g semantic-release @semantic-release/git @semantic-release/changelog

      - name: Release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: semantic-release
```

```json
// .releaserc.json - 语义化发布配置
{
  "branches": ["main"],
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/changelog",
    "@semantic-release/github",
    [
      "@semantic-release/git",
      {
        "assets": ["CHANGELOG.md"],
        "message": "chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}"
      }
    ]
  ]
}
```

### 9.5 签名验证

```yaml
# GitHub Actions - 制品签名
name: Sign Artifacts

on:
  push:
    tags: ['v*']

jobs:
  sign:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: write
    steps:
      - uses: actions/checkout@v4

      # 安装cosign
      - name: Install cosign
        uses: sigstore/cosign-installer@v3

      # 构建镜像
      - name: Build image
        run: docker build -t ghcr.io/${{ github.repository }}:${{ github.ref_name }} .

      # 推送镜像
      - name: Push image
        run: docker push ghcr.io/${{ github.repository }}:${{ github.ref_name }}

      # 签名镜像
      - name: Sign image
        run: |
          cosign sign --yes \
            ghcr.io/${{ github.repository }}:${{ github.ref_name }}

      # 验证签名
      - name: Verify signature
        run: |
          cosign verify \
            --certificate-identity=${{ github.server_url }}/${{ github.repository }}/.github/workflows/sign.yml@refs/tags/${{ github.ref_name }} \
            --certificate-oidc-issuer=https://token.actions.githubusercontent.com \
            ghcr.io/${{ github.repository }}:${{ github.ref_name }}
```

```bash
#!/bin/bash
# verify-image.sh - 镜像签名验证脚本

IMAGE=$1

if [ -z "$IMAGE" ]; then
    echo "Usage: $0 <image>"
    exit 1
fi

echo "Verifying signature for $IMAGE"

cosign verify \
    --certificate-identity-regexp="https://github.com/.*/.github/workflows/.*" \
    --certificate-oidc-issuer="https://token.actions.githubusercontent.com" \
    $IMAGE

if [ $? -eq 0 ]; then
    echo "✅ Signature verified successfully!"
else
    echo "❌ Signature verification failed!"
    exit 1
fi
```

### 9.6 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 多仓库备份 | 推送到多个镜像仓库 | 单一仓库 |
| 版本标签 | 使用语义化版本 | 仅使用latest |
| 制品签名 | 使用cosign签名 | 无签名验证 |
| 生命周期 | 设置保留策略 | 无限增长 |
| 可追溯性 | 关联源代码版本 | 无版本关联 |

---

## 10. 监控与可观测性集成

### 10.1 概念定义

**可观测性（Observability）** 是通过系统的外部输出（日志、指标、追踪）来理解系统内部状态的能力。对于Go项目，可观测性包括日志收集、指标监控、分布式追踪和告警配置。

**三大支柱：**

- **日志（Logs）**：离散事件记录
- **指标（Metrics）**：可聚合的数值数据
- **追踪（Traces）**：请求在分布式系统中的完整路径

### 10.2 日志收集

```go
// internal/logger/logger.go - 结构化日志
package logger

import (
 "os"
 "time"

 "github.com/rs/zerolog"
 "github.com/rs/zerolog/log"
)

// Config 日志配置
type Config struct {
 Level      string
 Format     string // json 或 console
 Output     string // stdout 或 file
 FilePath   string
}

// Init 初始化日志
func Init(cfg Config) {
 // 设置日志级别
 level, err := zerolog.ParseLevel(cfg.Level)
 if err != nil {
  level = zerolog.InfoLevel
 }
 zerolog.SetGlobalLevel(level)

 // 设置时间格式
 zerolog.TimeFieldFormat = time.RFC3339

 // 配置输出
 var output zerolog.ConsoleWriter
 if cfg.Format == "console" {
  output = zerolog.ConsoleWriter{
   Out:        os.Stdout,
   TimeFormat: time.RFC3339,
  }
  log.Logger = zerolog.New(output).With().Timestamp().Logger()
 } else {
  log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
 }

 // 添加服务信息
 log.Logger = log.Logger.With().
  Str("service", os.Getenv("APP_NAME")).
  Str("version", os.Getenv("APP_VERSION")).
  Str("environment", os.Getenv("APP_ENV")).
  Logger()
}

// Logger 获取日志实例
func Logger() zerolog.Logger {
 return log.Logger
}

// WithContext 添加上下文信息
func WithContext(ctx context.Context) zerolog.Logger {
 logger := log.Logger

 if traceID := ctx.Value("trace_id"); traceID != nil {
  logger = logger.With().Str("trace_id", traceID.(string)).Logger()
 }

 if userID := ctx.Value("user_id"); userID != nil {
  logger = logger.With().Int64("user_id", userID.(int64)).Logger()
 }

 return logger
}
```

```go
// 使用示例
package main

import (
 "myapp/internal/logger"
)

func main() {
 logger.Init(logger.Config{
  Level:  "info",
  Format: "json",
 })

 log := logger.Logger()

 log.Info().
  Str("event", "server_start").
  Int("port", 8080).
  Msg("Server starting")

 // 错误日志
 if err := doSomething(); err != nil {
  log.Error().
   Err(err).
   Str("operation", "do_something").
   Msg("Operation failed")
 }
}
```

### 10.3 指标监控

```go
// internal/metrics/metrics.go - Prometheus指标
package metrics

import (
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
 // HTTP请求计数
 HTTPRequestsTotal = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "http_requests_total",
   Help: "Total number of HTTP requests",
  },
  []string{"method", "path", "status"},
 )

 // HTTP请求持续时间
 HTTPRequestDuration = promauto.NewHistogramVec(
  prometheus.HistogramOpts{
   Name:    "http_request_duration_seconds",
   Help:    "HTTP request duration in seconds",
   Buckets: prometheus.DefBuckets,
  },
  []string{"method", "path"},
 )

 // 活跃请求数
 HTTPRequestsInFlight = promauto.NewGauge(
  prometheus.GaugeOpts{
   Name: "http_requests_in_flight",
   Help: "Current number of HTTP requests being processed",
  },
 )

 // 数据库连接数
 DBConnections = promauto.NewGaugeVec(
  prometheus.GaugeOpts{
   Name: "db_connections",
   Help: "Number of database connections",
  },
  []string{"state"}, // open, in_use, idle
 )

 // 缓存命中率
 CacheHits = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "cache_hits_total",
   Help: "Total number of cache hits",
  },
  []string{"cache_name"},
 )

 CacheMisses = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "cache_misses_total",
   Help: "Total number of cache misses",
  },
  []string{"cache_name"},
 )

 // 业务指标
 UsersRegistered = promauto.NewCounter(
  prometheus.CounterOpts{
   Name: "users_registered_total",
   Help: "Total number of registered users",
  },
 )

 OrdersCreated = promauto.NewCounterVec(
  prometheus.CounterOpts{
   Name: "orders_created_total",
   Help: "Total number of orders created",
  },
  []string{"status"},
 )
)
```

```go
// Gin中间件 - 指标收集
package middleware

import (
 "strconv"
 "time"

 "github.com/gin-gonic/gin"
 "github.com/prometheus/client_golang/prometheus"
 "myapp/internal/metrics"
)

// MetricsMiddleware Prometheus指标中间件
func MetricsMiddleware() gin.HandlerFunc {
 return func(c *gin.Context) {
  start := time.Now()
  path := c.FullPath()
  if path == "" {
   path = "unknown"
  }

  // 增加活跃请求数
  metrics.HTTPRequestsInFlight.Inc()
  defer metrics.HTTPRequestsInFlight.Dec()

  // 处理请求
  c.Next()

  // 记录指标
  duration := time.Since(start).Seconds()
  status := strconv.Itoa(c.Writer.Status())
  method := c.Request.Method

  metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
  metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
 }
}
```

```go
// main.go - 启动指标服务器
package main

import (
 "net/http"

 "github.com/gin-gonic/gin"
 "github.com/prometheus/client_golang/prometheus/promhttp"
 "myapp/internal/middleware"
)

func main() {
 r := gin.New()

 // 指标中间件
 r.Use(middleware.MetricsMiddleware())

 // Prometheus指标端点
 r.GET("/metrics", gin.WrapH(promhttp.Handler()))

 // 健康检查
 r.GET("/health", func(c *gin.Context) {
  c.JSON(200, gin.H{"status": "ok"})
 })

 // 业务路由
 r.GET("/api/users", getUsers)
 r.POST("/api/users", createUser)

 r.Run(":8080")
}
```

### 10.4 告警配置

```yaml
# prometheus-rules.yml - Prometheus告警规则
groups:
  - name: go-app-alerts
    rules:
      # 高错误率告警
      - alert: HighErrorRate
        expr: |
          (
            sum(rate(http_requests_total{status=~"5.."}[5m]))
            /
            sum(rate(http_requests_total[5m]))
          ) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }} for the last 5 minutes"

      # 高延迟告警
      - alert: HighLatency
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
          ) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "P95 latency is {{ $value }}s for the last 5 minutes"

      # 服务宕机告警
      - alert: ServiceDown
        expr: up{job="go-app"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service is down"
          description: "Service {{ $labels.instance }} has been down for more than 1 minute"

      # 高CPU使用率告警
      - alert: HighCPUUsage
        expr: |
          100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage"
          description: "CPU usage is above 80% for the last 10 minutes"

      # 内存使用率告警
      - alert: HighMemoryUsage
        expr: |
          (
            node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes
          ) / node_memory_MemTotal_bytes * 100 > 85
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is above 85% for the last 10 minutes"

      # Pod重启告警
      - alert: PodRestarting
        expr: |
          increase(kube_pod_container_status_restarts_total[1h]) > 3
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Pod is restarting frequently"
          description: "Pod {{ $labels.pod }} has restarted {{ $value }} times in the last hour"
```

```yaml
# alertmanager.yml - Alertmanager配置
global:
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alerts@example.com'
  smtp_auth_username: 'alerts@example.com'
  smtp_auth_password: 'password'

route:
  receiver: 'default'
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  routes:
    - match:
        severity: critical
      receiver: 'pagerduty-critical'
      continue: true
    - match:
        severity: warning
      receiver: 'slack-warnings'

receivers:
  - name: 'default'
    email_configs:
      - to: 'devops@example.com'

  - name: 'pagerduty-critical'
    pagerduty_configs:
      - service_key: '<pagerduty-service-key>'
        severity: critical

  - name: 'slack-warnings'
    slack_configs:
      - api_url: '<slack-webhook-url>'
        channel: '#alerts'
        title: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```

### 10.5 Kubernetes监控集成

```yaml
# ServiceMonitor - Prometheus Operator配置
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: go-app-metrics
  namespace: monitoring
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app: go-app
  namespaceSelector:
    matchNames:
      - production
      - staging
  endpoints:
    - port: metrics
      path: /metrics
      interval: 15s
      scrapeTimeout: 10s
```

```yaml
# PodMonitor配置
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: go-app-pods
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app: go-app
  podMetricsEndpoints:
    - port: metrics
      path: /metrics
      interval: 15s
```

```yaml
# PrometheusRule - 告警规则
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: go-app-alerts
  namespace: monitoring
  labels:
    release: prometheus
spec:
  groups:
    - name: go-app
      rules:
        - alert: GoAppHighErrorRate
          expr: |
            sum(rate(http_requests_total{app="go-app",status=~"5.."}[5m]))
            /
            sum(rate(http_requests_total{app="go-app"}[5m])) > 0.05
          for: 5m
          labels:
            severity: critical
          annotations:
            summary: "Go App has high error rate"
```

### 10.6 Grafana仪表板

```json
{
  "dashboard": {
    "title": "Go Application Dashboard",
    "tags": ["go", "application"],
    "timezone": "UTC",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total[5m])) by (status)",
            "legendFormat": "{{status}}"
          }
        ],
        "yAxes": [
          {
            "label": "requests/sec"
          }
        ]
      },
      {
        "title": "Response Time (P95)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))",
            "legendFormat": "P95"
          }
        ],
        "yAxes": [
          {
            "label": "seconds"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))"
          }
        ],
        "format": "percentunit"
      },
      {
        "title": "Active Connections",
        "type": "singlestat",
        "targets": [
          {
            "expr": "http_requests_in_flight"
          }
        ]
      },
      {
        "title": "Go Routines",
        "type": "graph",
        "targets": [
          {
            "expr": "go_goroutines",
            "legendFormat": "goroutines"
          }
        ]
      },
      {
        "title": "Memory Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "go_memstats_heap_alloc_bytes",
            "legendFormat": "heap alloc"
          },
          {
            "expr": "go_memstats_sys_bytes",
            "legendFormat": "sys"
          }
        ]
      }
    ]
  }
}
```

### 10.7 分布式追踪

```go
// internal/tracing/tracing.go - OpenTelemetry追踪
package tracing

import (
 "context"

 "go.opentelemetry.io/otel"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/exporters/jaeger"
 "go.opentelemetry.io/otel/sdk/resource"
 sdktrace "go.opentelemetry.io/otel/sdk/trace"
 semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
 "go.opentelemetry.io/otel/trace"
)

// InitTracer 初始化追踪器
func InitTracer(serviceName, jaegerEndpoint string) (*sdktrace.TracerProvider, error) {
 // 创建Jaeger导出器
 exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
  jaeger.WithEndpoint(jaegerEndpoint),
 ))
 if err != nil {
  return nil, err
 }

 // 创建TracerProvider
 tp := sdktrace.NewTracerProvider(
  sdktrace.WithBatcher(exp),
  sdktrace.WithResource(resource.NewWithAttributes(
   semconv.SchemaURL,
   semconv.ServiceNameKey.String(serviceName),
   attribute.String("environment", "production"),
  )),
 )

 // 设置为全局TracerProvider
 otel.SetTracerProvider(tp)

 return tp, nil
}

// StartSpan 开始一个span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
 tracer := otel.Tracer("go-app")
 return tracer.Start(ctx, name, opts...)
}

// AddEvent 添加事件到span
func AddEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
 span.AddEvent(name, trace.WithAttributes(attrs...))
}
```

```go
// Gin中间件 - 追踪
package middleware

import (
 "github.com/gin-gonic/gin"
 "go.opentelemetry.io/otel/attribute"
 "go.opentelemetry.io/otel/propagation"
 "go.opentelemetry.io/otel/trace"
 "myapp/internal/tracing"
)

// TracingMiddleware 追踪中间件
func TracingMiddleware() gin.HandlerFunc {
 return func(c *gin.Context) {
  // 从请求头提取trace context
  ctx := propagation.TraceContext{}.Extract(c.Request.Context(),
   propagation.HeaderCarrier(c.Request.Header))

  // 开始新的span
  ctx, span := tracing.StartSpan(ctx, c.Request.Method+" "+c.FullPath())
  defer span.End()

  // 设置span属性
  span.SetAttributes(
   attribute.String("http.method", c.Request.Method),
   attribute.String("http.path", c.Request.URL.Path),
   attribute.String("http.host", c.Request.Host),
   attribute.String("http.user_agent", c.Request.UserAgent()),
  )

  // 将context保存到gin context
  c.Set("trace_ctx", ctx)

  // 处理请求
  c.Next()

  // 记录响应信息
  span.SetAttributes(
   attribute.Int("http.status_code", c.Writer.Status()),
   attribute.Int("http.response_size", c.Writer.Size()),
  )

  if c.Writer.Status() >= 500 {
   span.SetAttributes(attribute.Bool("error", true))
  }
 }
}
```

### 10.8 最佳实践总结

| 实践项 | 推荐做法 | 避免做法 |
|--------|----------|----------|
| 日志 | 结构化JSON日志 | 纯文本日志 |
| 指标 | 使用Prometheus客户端 | 自定义指标系统 |
| 告警 | 基于SLO设置阈值 | 过多无意义告警 |
| 追踪 | 使用OpenTelemetry | 厂商锁定方案 |
| 仪表板 | 统一Grafana仪表板 | 分散的监控工具 |
| 上下文 | 传递trace_id | 无关联的日志 |

---

## 附录：完整项目结构示例

```
my-go-app/
├── .github/
│   └── workflows/
│       ├── ci.yml              # 持续集成
│       ├── cd.yml              # 持续部署
│       └── release.yml         # 发布流程
├── .golangci.yml               # Lint配置
├── Dockerfile                  # 容器构建
├── docker-compose.yml          # 本地开发
├── Makefile                    # 构建脚本
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── myapp/
│       └── main.go             # 应用入口
├── internal/
│   ├── config/                 # 配置管理
│   ├── handler/                # HTTP处理器
│   ├── middleware/             # 中间件
│   ├── model/                  # 数据模型
│   ├── repository/             # 数据访问层
│   ├── service/                # 业务逻辑层
│   ├── logger/                 # 日志
│   ├── metrics/                # 指标
│   ├── tracing/                # 追踪
│   └── feature/                # 功能开关
├── pkg/                        # 公共库
├── api/                        # API定义
├── web/                        # 静态资源
├── scripts/
│   ├── build.sh
│   ├── test.sh
│   ├── deploy.sh
│   ├── canary-deploy.sh
│   └── rollback.sh
├── helm/
│   └── myapp/
│       ├── Chart.yaml
│       ├── values.yaml
│       ├── values-production.yaml
│       └── templates/
├── terraform/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── environments/
│       ├── dev/
│       ├── staging/
│       └── production/
├── k8s/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── hpa.yaml
│   └── configmap.yaml
├── monitoring/
│   ├── prometheus-rules.yml
│   ├── alertmanager.yml
│   └── grafana-dashboard.json
└── docs/
    ├── architecture.md
    ├── deployment.md
    └── operations.md
```

---

## 总结

本文档全面梳理了Go语言项目的CI/CD持续工作流，涵盖以下核心内容：

1. **GitHub Actions**：工作流配置、触发器、矩阵构建、缓存优化
2. **GitLab CI**：Pipeline配置、作业依赖、缓存与制品
3. **Docker容器化**：多阶段构建、镜像优化、安全扫描
4. **Kubernetes部署**：Deployment、Service、HPA、滚动更新
5. **Helm Chart**：Chart结构、模板编写、依赖管理
6. **Terraform**：基础设施即代码、状态管理
7. **持续集成**：代码检查、测试、覆盖率、安全扫描
8. **持续部署**：蓝绿部署、金丝雀发布、功能开关、回滚
9. **制品管理**：镜像仓库、版本管理、签名验证
10. **监控可观测性**：日志、指标、告警、追踪

每个工作流都包含了概念定义、流程图、完整配置、实现示例、反例说明、优化建议和最佳实践，为Go项目的DevOps实践提供了全面的参考指南。

---

*文档版本：1.0*
*最后更新：2024年*
