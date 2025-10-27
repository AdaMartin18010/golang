# CI/CD流程

**难度**: 中级 | **预计阅读**: 20分钟 | **前置知识**: Git、Docker

---

## 📖 概念介绍

CI/CD（持续集成/持续部署）自动化软件交付流程，提高开发效率和代码质量。

---

## 🎯 GitHub Actions

### 基础工作流

```.yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -cover ./...
    
    - name: Run linter
      run: |
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        golangci-lint run
```

---

### 构建和发布

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

### Docker构建和推送

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
        context: .
        push: true
        tags: |
          myapp/app:latest
          myapp/app:${{ github.sha }}
```

---

## 🔧 GitLab CI

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

## 📊 完整流程示例

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

## 💡 最佳实践

1. **分支策略**
   - main: 生产环境
   - develop: 开发环境
   - feature/*: 功能分支

2. **自动化测试**
   - 单元测试
   - 集成测试
   - 代码覆盖率

3. **代码质量**
   - Linting
   - 格式检查
   - 安全扫描

4. **部署策略**
   - 蓝绿部署
   - 金丝雀发布
   - 滚动更新

---

## 📚 相关资源

- [GitHub Actions](https://docs.github.com/en/actions)
- [GitLab CI/CD](https://docs.gitlab.com/ee/ci/)

**下一步**: [05-监控与日志](./05-监控与日志.md)

---

**最后更新**: 2025-10-28

