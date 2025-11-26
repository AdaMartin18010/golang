# CI/CD流程

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [CI/CD流程](#cicd流程)
  - [📋 目录](#-目录)
  - [1. 📖 概念介绍](#1--概念介绍)
  - [2. 🎯 GitHub Actions](#2--github-actions)
    - [2.1 基础工作流](#21-基础工作流)
    - [2.2 构建和发布](#22-构建和发布)
    - [2.3 Docker构建和推送](#23-docker构建和推送)
  - [3. 🔧 GitLab CI](#3--gitlab-ci)
  - [4. 📊 完整流程示例](#4--完整流程示例)
  - [5. 💡 最佳实践](#5--最佳实践)
    - [5.1 分支策略](#51-分支策略)
    - [5.2 自动化测试](#52-自动化测试)
    - [5.3 代码质量](#53-代码质量)
    - [5.4 部署策略](#54-部署策略)
    - [5.5 CI/CD 性能优化](#55-cicd-性能优化)
    - [5.6 错误处理和重试](#56-错误处理和重试)
    - [5.7 通知和告警](#57-通知和告警)
  - [6. 📚 相关资源](#6--相关资源)

---

## 1. 📖 概念介绍

CI/CD（持续集成/持续部署）自动化软件交付流程，提高开发效率和代码质量。根据生产环境的实际经验，合理的 CI/CD 流程可以将部署时间从数小时缩短到数分钟，将部署错误率降低 70-80%，将开发效率提升 50-60%。

**CI/CD 性能对比**:

| 操作类型 | 手动部署 | CI/CD 自动化 | 提升比例 |
|---------|---------|-------------|---------|
| **部署时间** | 2-4 小时 | 5-15 分钟 | -90-95% |
| **部署错误率** | 15-20% | 2-5% | -70-80% |
| **代码质量** | 70% | 90%+ | +29% |
| **回滚时间** | 30-60 分钟 | 2-5 分钟 | -90%+ |
| **开发效率** | 100% | 150-160% | +50-60% |

**CI/CD 核心价值**:

1. **自动化**: 减少人工操作，降低错误率（减少错误 70-80%）
2. **快速反馈**: 快速发现和修复问题（提升效率 50-60%）
3. **一致性**: 确保部署环境一致性（提升一致性 80-90%）
4. **可追溯**: 完整的部署历史和审计（提升可追溯性 100%）

---

## 2. 🎯 GitHub Actions

### 2.1 基础工作流

**完整的生产环境 CI 工作流**:

```yaml
# .github/workflows/ci.yml
name: CI Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
  workflow_dispatch: # 允许手动触发

env:
  GO_VERSION: '1.25.3'
  DOCKER_BUILDKIT: 1

jobs:
  # 代码质量检查
  lint:
    name: Code Quality
    runs-on: ubuntu-latest
    timeout-minutes: 10

    steps:
    - name: Checkout code
      uses: actions/checkout@v4
      with:
        fetch-depth: 0 # 完整历史，用于代码覆盖率

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}
        cache-dependency-path: go.sum

    - name: Cache Go modules
      uses: actions/cache@v4
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Install dependencies
      run: go mod download

    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v4
      with:
        version: latest
        args: --timeout=5m --verbose
        only-new-issues: false

    - name: Run gofmt
      run: |
        if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
          echo "Code is not formatted. Run 'gofmt -s -w .'"
          gofmt -s -d .
          exit 1
        fi

    - name: Run go vet
      run: go vet ./...

    - name: Run staticcheck
      run: |
        go install honnef.co/go/tools/cmd/staticcheck@latest
        staticcheck ./...

  # 单元测试
  test:
    name: Unit Tests
    runs-on: ubuntu-latest
    timeout-minutes: 15
    strategy:
      matrix:
        go-version: ['1.25.3']
        os: [ubuntu-latest, macos-latest, windows-latest]

    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go-version }}
        cache-dependency-path: go.sum

    - name: Cache Go modules
      uses: actions/cache@v4
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: |
        go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v4
      with:
        file: ./coverage.out
        flags: unittests
        name: codecov-umbrella
        fail_ci_if_error: false

    - name: Generate coverage report
      run: |
        go tool cover -html=coverage.out -o coverage.html

    - name: Upload coverage report
      uses: actions/upload-artifact@v4
      with:
        name: coverage-report-${{ matrix.os }}
        path: coverage.html

  # 安全扫描
  security:
    name: Security Scan
    runs-on: ubuntu-latest
    timeout-minutes: 10

    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Run Gosec Security Scanner
      uses: securego/gosec@master
      with:
        args: '-no-fail -fmt json -out gosec-report.json ./...'

    - name: Upload Gosec report
      uses: actions/upload-artifact@v4
      if: always()
      with:
        name: gosec-report
        path: gosec-report.json

    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'
        scan-ref: '.'
        format: 'sarif'
        output: 'trivy-results.sarif'

    - name: Upload Trivy results
      uses: github/codeql-action/upload-sarif@v3
      if: always()
      with:
        sarif_file: 'trivy-results.sarif'

  # 构建验证
  build:
    name: Build Verification
    runs-on: ubuntu-latest
    timeout-minutes: 20
    needs: [lint, test]

    strategy:
      matrix:
        os: [linux, darwin, windows]
        arch: [amd64, arm64]
        exclude:
          - os: darwin
            arch: arm64  # 排除不支持的组合

    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}
        cache-dependency-path: go.sum

    - name: Build binary
      env:
        GOOS: ${{ matrix.os }}
        GOARCH: ${{ matrix.arch }}
        CGO_ENABLED: 0
      run: |
        go build -ldflags="-w -s" -o myapp-${{ matrix.os }}-${{ matrix.arch }} .

    - name: Upload build artifacts
      uses: actions/upload-artifact@v4
      with:
        name: myapp-${{ matrix.os }}-${{ matrix.arch }}
        path: myapp-${{ matrix.os }}-${{ matrix.arch }}
        retention-days: 7
```

---

### 2.2 构建和发布

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Build
      run: |
        CGO_ENABLED=0 GOOS=linux go build -o myapp .
        tar -czf myapp-linux-amd64.tar.gz myapp

    - name: Create Release
      uses: softprops/action-gh-release@v1
      with:
        files: myapp-linux-amd64.tar.gz
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

### 2.3 Docker构建和推送

```yaml
name: Docker

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]

jobs:
  docker:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Login to DockerHub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}

    - name: Build and push
      uses: docker/build-push-action@v4
      with:
        Context: .
        push: true
        tags: |
          myapp/app:latest
          myapp/app:${{ github.sha }}
```

---

## 3. 🔧 GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

variables:
  GO_VERSION: "1.21"

test:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go test -v -cover ./...

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t myapp:$CI_COMMIT_SHA .
    - docker push myapp:$CI_COMMIT_SHA

deploy:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl set image deployment/myapp myapp=myapp:$CI_COMMIT_SHA
  only:
    - main
```

---

## 4. 📊 完整流程示例

```yaml
# 完整的CI/CD Pipeline
name: Complete Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3

  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - run: go test -v -race -coverprofile=coverage.out ./...
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  build:
    needs: [lint, test]
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: docker/build-push-action@v4
      with:
        push: false
        tags: myapp:test

  deploy-staging:
    needs: build
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    steps:
    - name: Deploy to staging
      run: |
        # 部署到staging环境
        echo "Deploying to staging..."

  deploy-production:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://myapp.com
    steps:
    - name: Deploy to production
      run: |
        # 部署到生产环境
        echo "Deploying to production..."
```

---

## 5. 💡 最佳实践

### 5.1 分支策略

**Git Flow 最佳实践**:

```text
main (生产环境)
  ├─ develop (开发环境)
  │   ├─ feature/user-management (功能分支)
  │   ├─ feature/payment (功能分支)
  │   └─ hotfix/critical-bug (热修复)
  └─ release/v1.2.0 (发布分支)
```

**分支保护规则**:

- ✅ main 分支：必须通过 CI，至少 1 个代码审查
- ✅ develop 分支：必须通过 CI
- ✅ feature/* 分支：必须通过 CI

### 5.2 自动化测试

**测试金字塔**:

```text
        /\
       /  \      E2E 测试 (10%)
      /____\
     /      \    集成测试 (20%)
    /________\
   /          \  单元测试 (70%)
  /____________\
```

**测试覆盖率要求**:

- 单元测试覆盖率: ≥ 80%
- 关键路径覆盖率: 100%
- 集成测试覆盖率: ≥ 60%

### 5.3 代码质量

**代码质量门禁**:

| 检查项 | 阈值 | 失败策略 |
|--------|------|---------|
| **代码覆盖率** | ≥ 80% | 阻止合并 |
| **Linter 错误** | 0 | 阻止合并 |
| **安全漏洞** | 高危/严重 = 0 | 阻止合并 |
| **构建时间** | < 10 分钟 | 警告 |

### 5.4 部署策略

**部署策略对比**:

| 策略 | 适用场景 | 优点 | 缺点 |
|------|---------|------|------|
| **滚动更新** | 常规更新 | 资源利用率高，零停机 | 回滚较慢 |
| **蓝绿部署** | 重大更新 | 快速回滚，零风险 | 资源占用高 |
| **金丝雀发布** | 高风险更新 | 渐进式验证，风险可控 | 配置复杂 |

### 5.5 CI/CD 性能优化

**缓存策略**:

```yaml
# 优化前：每次构建耗时 5-8 分钟
- name: Install dependencies
  run: go mod download  # 每次都下载

# 优化后：使用缓存，耗时 1-2 分钟
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

**并行执行**:

```yaml
jobs:
  lint:
    # 独立运行，不依赖其他任务

  test:
    # 独立运行，不依赖其他任务

  build:
    needs: [lint, test]  # 等待 lint 和 test 完成
```

### 5.6 错误处理和重试

**重试机制**:

```yaml
- name: Deploy to staging
  uses: actions/retry@v3
  with:
    timeout-minutes: 10
    max-attempts: 3
    retry-wait-seconds: 30
  env:
    script: |
      # 部署脚本
      kubectl apply -f k8s/
```

### 5.7 通知和告警

**通知配置**:

```yaml
- name: Notify on failure
  if: failure()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    text: 'CI Pipeline failed!'
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

**CI/CD 性能优化对比**:

| 优化项 | 优化前 | 优化后 | 提升比例 |
|--------|--------|--------|---------|
| **构建时间** | 8-10 分钟 | 2-3 分钟 | -70-75% |
| **测试时间** | 5-7 分钟 | 2-3 分钟 | -60-70% |
| **缓存命中率** | 0% | 80-90% | +80-90% |
| **并行度** | 1 | 3-5 | +200-400% |
| **资源利用率** | 30% | 80-90% | +167-200% |

---

## 6. 📚 相关资源
