# GitHub Actions CI/CD

> **简介**: 使用GitHub Actions构建Go微服务的CI/CD流程，涵盖自动化测试、镜像构建和Kubernetes部署
> **版本**: Go 1.23+  
> **难度**: ⭐⭐⭐⭐  
> **标签**: #CI/CD #GitHub-Actions #自动化部署 #DevOps

<!-- TOC START -->
- [GitHub Actions CI/CD](#github-actions-cicd)
  - [7.1 📚 GitHub Actions概述](#71--github-actions概述)
  - [7.2 🎯 工作流配置](#72--工作流配置)
    - [基础工作流](#基础工作流)
  - [7.3 🧪 自动化测试](#73--自动化测试)
    - [单元测试与集成测试](#单元测试与集成测试)
    - [代码质量检查](#代码质量检查)
  - [7.4 🐳 Docker镜像构建](#74--docker镜像构建)
    - [多架构构建](#多架构构建)
    - [优化构建速度](#优化构建速度)
  - [7.5 🔐 安全扫描](#75--安全扫描)
    - [漏洞扫描](#漏洞扫描)
    - [依赖检查](#依赖检查)
  - [7.6 🚀 自动部署](#76--自动部署)
    - [Kubernetes部署](#kubernetes部署)
    - [ArgoCD同步](#argocd同步)
    - [Helm Chart部署](#helm-chart部署)
  - [7.7 📊 矩阵策略](#77--矩阵策略)
    - [多版本测试](#多版本测试)
    - [多环境部署](#多环境部署)
  - [7.8 🎯 最佳实践](#78--最佳实践)
  - [7.9 ⚠️ 常见问题](#79-️-常见问题)
    - [Q1: 如何加速GitHub Actions？](#q1-如何加速github-actions)
    - [Q2: Secret如何管理？](#q2-secret如何管理)
    - [Q3: 如何调试失败的工作流？](#q3-如何调试失败的工作流)
  - [7.10 📚 扩展阅读](#710--扩展阅读)
    - [官方文档](#官方文档)
    - [相关文档](#相关文档)
<!-- TOC END -->


## 📋 目录

- [1. 7.1 📚 GitHub Actions概述](#71--github-actions概述)
- [2. 7.2 🎯 工作流配置](#72--工作流配置)
- [3. 7.3 🧪 自动化测试](#73--自动化测试)
- [4. 7.4 🐳 Docker镜像构建](#74--docker镜像构建)
- [5. 7.5 🔐 安全扫描](#75--安全扫描)
- [6. 7.6 🚀 自动部署](#76--自动部署)
- [7. 7.7 📊 矩阵策略](#77--矩阵策略)
- [8. 7.8 🎯 最佳实践](#78--最佳实践)
- [9. 7.9 ⚠️ 常见问题](#79-️-常见问题)
- [10. 7.10 📚 扩展阅读](#710--扩展阅读)

---

## 7.1 📚 GitHub Actions概述

**GitHub Actions**: GitHub原生CI/CD平台，通过YAML配置自动化工作流。

**核心概念**:

- **Workflow**: 自动化流程
- **Job**: 工作流中的任务
- **Step**: Job中的具体步骤
- **Action**: 可复用的步骤单元
- **Runner**: 执行环境（GitHub托管或自托管）

## 7.2 🎯 工作流配置

### 基础工作流

```yaml
name: Go CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  workflow_dispatch:

env:
  GO_VERSION: '1.21'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - name: 检出代码
      uses: actions/checkout@v4
    
    - name: 设置Go环境
      uses: actions/setup-go@v5
      with:
        go-version: ${{ env.GO_VERSION }}
        cache: true
    
    - name: 下载依赖
      run: go mod download
    
    - name: 运行测试
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: 上传覆盖率
      uses: codecov/codecov-action@v4
      with:
        files: ./coverage.out
```

## 7.3 🧪 自动化测试

### 单元测试与集成测试

```yaml
jobs:
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
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v4
    
    - uses: actions/setup-go@v5
      with:
        go-version: '1.21'
    
    - name: 运行单元测试
      run: go test -v -short ./...
    
    - name: 运行集成测试
      env:
        DATABASE_URL: postgres://postgres:postgres@localhost:5432/test?sslmode=disable
        REDIS_URL: redis://localhost:6379
      run: go test -v -run Integration ./...
```

### 代码质量检查

```yaml
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - uses: actions/setup-go@v5
      with:
        go-version: '1.21'
    
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v4
      with:
        version: latest
        args: --timeout=5m
    
    - name: Go格式检查
      run: |
        if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
          echo "Go files must be formatted with gofmt -s"
          gofmt -s -l .
          exit 1
        fi
```

## 7.4 🐳 Docker镜像构建

### 多架构构建

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    
    steps:
    - uses: actions/checkout@v4
    
    - name: 设置QEMU
      uses: docker/setup-qemu-action@v3
    
    - name: 设置Docker Buildx
      uses: docker/setup-buildx-action@v3
    
    - name: 登录GHCR
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: 提取元数据
      id: meta
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=semver,pattern={{major}}.{{minor}}
          type=sha,prefix={{branch}}-
    
    - name: 构建并推送
      uses: docker/build-push-action@v5
      with:
        context: .
        platforms: linux/amd64,linux/arm64
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        cache-from: type=gha
        cache-to: type=gha,mode=max
```

### 优化构建速度

```yaml
- name: 构建（带缓存）
  uses: docker/build-push-action@v5
  with:
    context: .
    push: true
    tags: ${{ steps.meta.outputs.tags }}
    cache-from: type=registry,ref=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:buildcache
    cache-to: type=registry,ref=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:buildcache,mode=max
    build-args: |
      BUILDKIT_INLINE_CACHE=1
```

## 7.5 🔐 安全扫描

### 漏洞扫描

```yaml
jobs:
  security:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: 运行Trivy扫描
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'
        scan-ref: '.'
        format: 'sarif'
        output: 'trivy-results.sarif'
        severity: 'CRITICAL,HIGH'
    
    - name: 上传到GitHub Security
      uses: github/codeql-action/upload-sarif@v3
      with:
        sarif_file: 'trivy-results.sarif'
    
    - name: 扫描Docker镜像
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
        format: 'table'
        exit-code: '1'
        severity: 'CRITICAL,HIGH'
```

### 依赖检查

```yaml
- name: Go依赖审计
  run: |
    go list -json -m all | nancy sleuth
    
- name: Snyk依赖扫描
  uses: snyk/actions/golang@master
  env:
    SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
  with:
    command: test
```

## 7.6 🚀 自动部署

### Kubernetes部署

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    needs: [test, build, security]
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v4
    
    - name: 配置kubectl
      uses: azure/k8s-set-context@v3
      with:
        method: kubeconfig
        kubeconfig: ${{ secrets.KUBE_CONFIG }}
    
    - name: 部署到Kubernetes
      run: |
        kubectl set image deployment/user-service \
          user-service=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} \
          -n production
        
        kubectl rollout status deployment/user-service -n production
```

### ArgoCD同步

```yaml
- name: 更新Manifest
  run: |
    git clone https://${{ secrets.MANIFEST_TOKEN }}@github.com/myorg/k8s-manifests.git
    cd k8s-manifests
    
    # 使用yq更新镜像
    yq eval '.spec.template.spec.containers[0].image = "${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}"' \
      -i apps/user-service/deployment.yaml
    
    git config user.name "GitHub Actions"
    git config user.email "actions@github.com"
    git add apps/user-service/deployment.yaml
    git commit -m "Update image to ${{ github.sha }}"
    git push
```

### Helm Chart部署

```yaml
- name: Helm部署
  run: |
    helm upgrade --install user-service ./charts/user-service \
      --set image.tag=${{ github.sha }} \
      --set image.repository=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }} \
      --namespace production \
      --create-namespace \
      --wait
```

## 7.7 📊 矩阵策略

### 多版本测试

```yaml
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.20', '1.21', '1.22']
    
    steps:
    - uses: actions/checkout@v4
    
    - uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: 运行测试
      run: go test -v ./...
```

### 多环境部署

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        environment: [staging, production]
    
    environment:
      name: ${{ matrix.environment }}
      url: https://${{ matrix.environment }}.example.com
    
    steps:
    - name: 部署到${{ matrix.environment }}
      run: |
        kubectl set image deployment/user-service \
          user-service=${{ env.IMAGE }}:${{ github.sha }} \
          -n ${{ matrix.environment }}
```

## 7.8 🎯 最佳实践

1. **缓存依赖**: 使用`cache`加速构建
2. **并行执行**: 合理使用`matrix`和并行job
3. **条件执行**: 使用`if`避免不必要的步骤
4. **Secret管理**: 使用GitHub Secrets存储敏感信息
5. **环境保护**: 为生产环境配置审批流程
6. **可复用Actions**: 创建组合Action
7. **工作流触发**: 合理配置`on`事件
8. **权限最小化**: 使用`permissions`限制权限
9. **超时设置**: 防止工作流无限运行
10. **通知集成**: 集成Slack/钉钉通知

## 7.9 ⚠️ 常见问题

### Q1: 如何加速GitHub Actions？

**A**:

- 使用缓存（`actions/cache`, `cache-from`）
- 并行执行job
- 自托管Runner
- 减少不必要的步骤

### Q2: Secret如何管理？

**A**:

- 使用GitHub Secrets
- 环境级别Secret
- 组织级别Secret
- 集成Vault等外部工具

### Q3: 如何调试失败的工作流？

**A**:

```yaml
- name: 设置tmate调试
  if: ${{ failure() }}
  uses: mxschmitt/action-tmate@v3
  timeout-minutes: 15
```

## 7.10 📚 扩展阅读

### 官方文档

- [GitHub Actions文档](https://docs.github.com/en/actions)
- [Workflow语法](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [Actions Marketplace](https://github.com/marketplace?type=actions)

### 相关文档

- [06-GitOps部署.md](./06-GitOps部署.md)
- [08-GitLab-CI.md](./08-GitLab-CI.md)
- [09-多环境部署.md](./09-多环境部署.md)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: GitHub Actions, Go 1.21+, Kubernetes 1.27+
