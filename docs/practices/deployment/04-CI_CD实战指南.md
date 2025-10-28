# CI/CD实战指南

## 📋 目录

- [1. CI/CD概述](#1-cicd概述)
- [2. GitHub Actions](#2-github-actions)
- [3. GitLab CI](#3-gitlab-ci)
- [4. Jenkins](#4-jenkins)
- [5. 完整Pipeline](#5-完整pipeline)
- [6. 最佳实践](#6-最佳实践)
- [7. 性能优化](#7-性能优化)
- [8. 故障排查](#8-故障排查)

---

## 1. CI/CD概述

### CI/CD流程

```
代码提交 (git push)
    ↓
触发CI Pipeline
    ├─ 代码检查 (Linter)
    ├─ 单元测试 (go test)
    ├─ 集成测试
    ├─ 安全扫描
    └─ 构建镜像 (Docker)
    ↓
CD Pipeline
    ├─ 部署到测试环境
    ├─ 自动化测试
    ├─ 部署到预发布
    └─ 部署到生产环境
    ↓
监控和告警
```

---

## 2. GitHub Actions

### 基础工作流

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.25.3'
        cache: true
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
```

### 完整CI Pipeline

```yaml
# .github/workflows/complete-ci.yml
name: Complete CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25.3'
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --timeout=5m

  test:
    name: Test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.24', '1.25.3']
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go ${{ matrix.go-version }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html
      
      - name: Upload coverage
        uses: actions/upload-artifact@v3
        with:
          name: coverage
          path: coverage.html

  security:
    name: Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Gosec
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt json -out results.json ./...'
      
      - name: Run Trivy
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'

  build:
    name: Build
    needs: [lint, test, security]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25.3'
      
      - name: Build
        run: |
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-linux-amd64 ./cmd/app
          CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o app-darwin-amd64 ./cmd/app
          CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o app-windows-amd64.exe ./cmd/app
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: binaries
          path: app-*
```

### Docker构建与推送

```yaml
# .github/workflows/docker.yml
name: Docker Build and Push

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      
      - name: Setup Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Log in to Container Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          build-args: |
            VERSION=${{ github.ref_name }}
            COMMIT=${{ github.sha }}
```

### 自动部署到Kubernetes

```yaml
# .github/workflows/deploy.yml
name: Deploy to Kubernetes

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      
      - name: Setup kubectl
        uses: azure/setup-kubectl@v3
        with:
          version: 'v1.28.0'
      
      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
          export KUBECONFIG=./kubeconfig
      
      - name: Deploy to staging
        run: |
          kubectl set image deployment/myapp \
            myapp=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} \
            -n staging
          
          kubectl rollout status deployment/myapp -n staging --timeout=5m
      
      - name: Run smoke tests
        run: |
          ./scripts/smoke-test.sh staging
      
      - name: Deploy to production
        if: success()
        run: |
          kubectl set image deployment/myapp \
            myapp=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} \
            -n production
          
          kubectl rollout status deployment/myapp -n production --timeout=10m
```

---

## 3. GitLab CI

### 基础配置

```yaml
# .gitlab-ci.yml
image: golang:1.25.3

stages:
  - test
  - build
  - deploy

variables:
  GO111MODULE: "on"
  CGO_ENABLED: "0"

cache:
  paths:
    - .cache/go-build
    - .cache/go-mod

before_script:
  - mkdir -p .cache/go-build .cache/go-mod
  - export GOPATH=$CI_PROJECT_DIR/.cache/go-mod
  - export GOCACHE=$CI_PROJECT_DIR/.cache/go-build

test:
  stage: test
  script:
    - go fmt $(go list ./... | grep -v /vendor/)
    - go vet $(go list ./... | grep -v /vendor/)
    - go test -race -coverprofile=coverage.out ./...
  coverage: '/total:.*?(\d+\.\d+)%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml

lint:
  stage: test
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run -v

build:
  stage: build
  script:
    - go build -o myapp ./cmd/app
  artifacts:
    paths:
      - myapp
    expire_in: 1 week

docker-build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker tag $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA $CI_REGISTRY_IMAGE:latest
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - docker push $CI_REGISTRY_IMAGE:latest
  only:
    - main
```

### 多环境部署

```yaml
# .gitlab-ci.yml (完整版)
stages:
  - test
  - build
  - deploy-staging
  - test-staging
  - deploy-production

# 测试阶段
unit-test:
  stage: test
  script:
    - go test -v -race ./...

integration-test:
  stage: test
  services:
    - postgres:15
    - redis:7
  variables:
    POSTGRES_DB: test
    POSTGRES_USER: test
    POSTGRES_PASSWORD: test
  script:
    - go test -v -tags=integration ./...

# 构建阶段
build-binary:
  stage: build
  script:
    - CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app
  artifacts:
    paths:
      - app

build-docker:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

# 部署到Staging
deploy-staging:
  stage: deploy-staging
  image: bitnami/kubectl:latest
  script:
    - kubectl config set-cluster k8s --server="$KUBE_URL" --insecure-skip-tls-verify=true
    - kubectl config set-credentials admin --token="$KUBE_TOKEN"
    - kubectl config set-context default --cluster=k8s --user=admin
    - kubectl config use-context default
    - kubectl set image deployment/myapp myapp=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA -n staging
    - kubectl rollout status deployment/myapp -n staging
  environment:
    name: staging
    url: https://staging.example.com
  only:
    - develop

# Staging环境测试
smoke-test-staging:
  stage: test-staging
  script:
    - ./scripts/smoke-test.sh https://staging.example.com
  only:
    - develop

# 部署到Production
deploy-production:
  stage: deploy-production
  image: bitnami/kubectl:latest
  script:
    - kubectl set image deployment/myapp myapp=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA -n production
    - kubectl rollout status deployment/myapp -n production
  environment:
    name: production
    url: https://example.com
  when: manual
  only:
    - main
```

---

## 4. Jenkins

### Jenkinsfile

```groovy
// Jenkinsfile
pipeline {
    agent {
        docker {
            image 'golang:1.25.3'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
    
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = '0'
        DOCKER_REGISTRY = 'registry.example.com'
        IMAGE_NAME = 'myapp'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Dependencies') {
            steps {
                sh 'go mod download'
                sh 'go mod verify'
            }
        }
        
        stage('Lint') {
            steps {
                sh 'go fmt ./...'
                sh 'go vet ./...'
                sh 'golangci-lint run'
            }
        }
        
        stage('Test') {
            steps {
                sh 'go test -v -race -coverprofile=coverage.out ./...'
                sh 'go tool cover -html=coverage.out -o coverage.html'
            }
            post {
                always {
                    publishHTML([
                        reportDir: '.',
                        reportFiles: 'coverage.html',
                        reportName: 'Coverage Report'
                    ])
                }
            }
        }
        
        stage('Build') {
            steps {
                sh 'go build -o app ./cmd/app'
            }
        }
        
        stage('Docker Build') {
            when {
                branch 'main'
            }
            steps {
                script {
                    def imageTag = "${DOCKER_REGISTRY}/${IMAGE_NAME}:${env.BUILD_NUMBER}"
                    def latestTag = "${DOCKER_REGISTRY}/${IMAGE_NAME}:latest"
                    
                    sh "docker build -t ${imageTag} -t ${latestTag} ."
                    sh "docker push ${imageTag}"
                    sh "docker push ${latestTag}"
                }
            }
        }
        
        stage('Deploy to Staging') {
            when {
                branch 'develop'
            }
            steps {
                sh '''
                    kubectl set image deployment/myapp \
                        myapp=${DOCKER_REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} \
                        -n staging
                    kubectl rollout status deployment/myapp -n staging
                '''
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            steps {
                input message: 'Deploy to production?', ok: 'Deploy'
                sh '''
                    kubectl set image deployment/myapp \
                        myapp=${DOCKER_REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} \
                        -n production
                    kubectl rollout status deployment/myapp -n production
                '''
            }
        }
    }
    
    post {
        success {
            slackSend color: 'good', message: "Build ${env.BUILD_NUMBER} succeeded"
        }
        failure {
            slackSend color: 'danger', message: "Build ${env.BUILD_NUMBER} failed"
        }
        always {
            cleanWs()
        }
    }
}
```

---

## 5. 完整Pipeline

### 多阶段Dockerfile

```dockerfile
# Dockerfile (多阶段构建)
# 第1阶段: 构建
FROM golang:1.25.3-alpine AS builder

WORKDIR /build

# 安装依赖
RUN apk add --no-cache git make

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o app ./cmd/app

# 第2阶段: 运行
FROM alpine:3.19

WORKDIR /app

# 安全: 创建非root用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 安装CA证书
RUN apk --no-cache add ca-certificates tzdata

# 从builder复制二进制文件
COPY --from=builder /build/app .

# 切换到非root用户
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/app", "healthcheck"]

ENTRYPOINT ["/app/app"]
```

### 部署脚本

```bash
#!/bin/bash
# scripts/deploy.sh

set -e

ENVIRONMENT=$1
IMAGE_TAG=$2

if [ -z "$ENVIRONMENT" ] || [ -z "$IMAGE_TAG" ]; then
    echo "Usage: $0 <environment> <image-tag>"
    exit 1
fi

echo "Deploying to $ENVIRONMENT with image tag $IMAGE_TAG"

# 1. 更新Kubernetes配置
kubectl set image deployment/myapp \
    myapp=registry.example.com/myapp:$IMAGE_TAG \
    -n $ENVIRONMENT

# 2. 等待滚动更新完成
kubectl rollout status deployment/myapp -n $ENVIRONMENT --timeout=5m

# 3. 运行健康检查
echo "Running health check..."
HEALTH_URL="https://$ENVIRONMENT.example.com/health"
for i in {1..30}; do
    if curl -f $HEALTH_URL > /dev/null 2>&1; then
        echo "Health check passed!"
        break
    fi
    echo "Waiting for app to be ready... ($i/30)"
    sleep 2
done

# 4. 运行烟雾测试
echo "Running smoke tests..."
./scripts/smoke-test.sh $ENVIRONMENT

echo "Deployment completed successfully!"
```

### 烟雾测试

```bash
#!/bin/bash
# scripts/smoke-test.sh

ENVIRONMENT=$1
BASE_URL="https://$ENVIRONMENT.example.com"

echo "Running smoke tests against $BASE_URL"

# 测试健康检查
echo "Testing health endpoint..."
if ! curl -f $BASE_URL/health > /dev/null; then
    echo "❌ Health check failed"
    exit 1
fi
echo "✅ Health check passed"

# 测试API
echo "Testing API endpoints..."
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/api/v1/users)
if [ "$RESPONSE" != "200" ]; then
    echo "❌ API test failed: HTTP $RESPONSE"
    exit 1
fi
echo "✅ API test passed"

echo "All smoke tests passed!"
```

---

## 6. 最佳实践

### 1. 缓存优化

```yaml
# GitHub Actions
- uses: actions/cache@v3
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

### 2. 并行测试

```yaml
# 矩阵构建
strategy:
  matrix:
    go-version: ['1.24', '1.25']
    os: [ubuntu-latest, macos-latest, windows-latest]
runs-on: ${{ matrix.os }}
```

### 3. 安全扫描

```yaml
# Trivy扫描
- name: Run Trivy
  uses: aquasecurity/trivy-action@master
  with:
    image-ref: 'myapp:latest'
    format: 'sarif'
    output: 'trivy-results.sarif'

# Gosec扫描
- name: Run Gosec
  run: gosec -fmt json -out results.json ./...
```

### 4. 制品管理

```yaml
# 上传制品
- name: Upload artifacts
  uses: actions/upload-artifact@v3
  with:
    name: myapp-${{ github.sha }}
    path: |
      app
      config/
      migrations/
    retention-days: 30
```

---

## 7. 性能优化

### 构建时间优化

**优化前**:
```yaml
# 每次都重新下载依赖
- run: go mod download
- run: go build ./...
# 耗时: 5分钟
```

**优化后**:
```yaml
# 使用缓存
- uses: actions/cache@v3
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

# 使用buildx缓存
- name: Build
  uses: docker/build-push-action@v4
  with:
    cache-from: type=gha
    cache-to: type=gha,mode=max
# 耗时: 1分钟 (5x提升)
```

### Pipeline优化对比

| 优化项 | 优化前 | 优化后 | 提升 |
|--------|--------|--------|------|
| Go模块缓存 | 2min | 10s | 12x |
| Docker层缓存 | 3min | 30s | 6x |
| 并行测试 | 5min | 1min | 5x |
| 总Pipeline时间 | 15min | 3min | 5x |

---

## 8. 故障排查

### 常见问题

#### 1. 构建缓存失效

```yaml
# 问题: 缓存key不稳定
key: go-cache

# 解决: 使用go.sum哈希
key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

#### 2. Docker构建慢

```dockerfile
# ❌ 问题: 每次都复制所有文件
COPY . .
RUN go build

# ✅ 解决: 先复制依赖文件
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build
```

#### 3. 测试不稳定

```yaml
# 添加重试机制
- name: Run tests
  uses: nick-invision/retry@v2
  with:
    timeout_minutes: 10
    max_attempts: 3
    command: go test -v ./...
```

---

## 🔗 相关资源

- [Docker部署](./02-Docker部署.md)
- [Kubernetes部署](./03-Kubernetes部署.md)
- [容器化最佳实践](./05-容器化最佳实践.md)

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3  
**文档类型**: CI/CD实战指南 ✨

