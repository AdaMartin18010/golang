# GitLab CI/CD

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.23+

---

## 📋 目录

- [GitLab CI/CD](#gitlab-cicd)
  - [📋 目录](#-目录)
  - [8.1 📚 GitLab CI/CD概述](#81--gitlab-cicd概述)
  - [8.2 🎯 Pipeline配置](#82--pipeline配置)
    - [基础Pipeline](#基础pipeline)
    - [多阶段Pipeline](#多阶段pipeline)
  - [8.3 🧪 测试阶段](#83--测试阶段)
    - [单元测试与集成测试](#单元测试与集成测试)
    - [代码质量](#代码质量)
  - [8.4 🐳 构建阶段](#84--构建阶段)
    - [Docker镜像构建](#docker镜像构建)
    - [使用Kaniko](#使用kaniko)
  - [8.5 🔐 安全扫描](#85--安全扫描)
    - [容器扫描](#容器扫描)
    - [Trivy扫描](#trivy扫描)
  - [8.6 🚀 部署阶段](#86--部署阶段)
    - [Kubernetes部署](#kubernetes部署)
    - [Helm部署](#helm部署)
  - [8.7 📊 高级特性](#87--高级特性)
    - [并行执行](#并行执行)
    - [动态子Pipeline](#动态子pipeline)
    - [自动回滚](#自动回滚)
  - [8.8 🎯 最佳实践](#88--最佳实践)
  - [8.9 ⚠️ 常见问题](#89-️-常见问题)
    - [Q1: 如何加速GitLab CI？](#q1-如何加速gitlab-ci)
    - [Q2: 如何调试Pipeline？](#q2-如何调试pipeline)
    - [Q3: 多项目如何共享配置？](#q3-多项目如何共享配置)
  - [8.10 📚 扩展阅读](#810--扩展阅读)
    - [官方文档](#官方文档)
    - [相关文档](#相关文档)

## 8.1 📚 GitLab CI/CD概述

**GitLab CI/CD**: GitLab内置的CI/CD平台，通过`.gitlab-ci.yml`配置Pipeline。

**核心概念**:

- **Pipeline**: 完整的CI/CD流程
- **Stage**: 阶段（如test、build、deploy）
- **Job**: 具体任务
- **Runner**: 执行环境（Shared或Self-hosted）
- **Artifact**: 产物

## 8.2 🎯 Pipeline配置

### 基础Pipeline

```yaml
# .gitlab-ci.yml
image: golang:1.21-alpine

variables:
  GO111MODULE: "on"
  CGO_ENABLED: 0
  DOCKER_DRIVER: overlay2

stages:
  - test
  - build
  - security
  - deploy

before_script:
  - go mod download

cache:
  paths:
    - .go/pkg/mod/

test:unit:
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
    expire_in: 30 days
```

### 多阶段Pipeline

```yaml
stages:
  - validate
  - test
  - build
  - security
  - review
  - deploy
  - cleanup

validate:lint:
  stage: validate
  image: golangci/golangci-lint:latest
  script:
    - golangci-lint run --timeout 5m
  only:
    - merge_requests
    - main

validate:format:
  stage: validate
  script:
    - gofmt -l . | tee /dev/stderr | wc -l | grep '^0$'
```

## 8.3 🧪 测试阶段

### 单元测试与集成测试

```yaml
test:unit:
  stage: test
  script:
    - go test -v -short ./...

test:integration:
  stage: test
  services:
    - postgres:15
    - redis:7-alpine
  variables:
    POSTGRES_DB: testdb
    POSTGRES_USER: testuser
    POSTGRES_PASSWORD: testpass
    DATABASE_URL: postgres://testuser:testpass@postgres:5432/testdb?sslmode=disable
    REDIS_URL: redis://redis:6379
  script:
    - go test -v -run Integration ./...

test:e2e:
  stage: test
  services:
    - docker:dind
  script:
    - docker-compose up -d
    - go test -v -tags=e2e ./tests/e2e/...
  after_script:
    - docker-compose down
```

### 代码质量

```yaml
code_quality:
  stage: test
  image: registry.gitlab.com/gitlab-org/ci-cd/codequality:latest
  services:
    - docker:dind
  script:
    - docker pull codeclimate/codeclimate:latest
    - docker run --env CODECLIMATE_CODE="$PWD" \
                 --volume "$PWD":/code \
                 --volume /var/run/docker.sock:/var/run/docker.sock \
                 --volume /tmp/cc:/tmp/cc \
                 codeclimate/codeclimate analyze -f json > gl-code-quality-report.json
  artifacts:
    reports:
      codequality: gl-code-quality-report.json
    expire_in: 1 week
```

## 8.4 🐳 构建阶段

### Docker镜像构建

```yaml
build:image:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  variables:
    DOCKER_TLS_CERTDIR: "/certs"
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA .
    - docker build -t $CI_REGISTRY_IMAGE:latest .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA
    - docker push $CI_REGISTRY_IMAGE:latest
  only:
    - main
    - tags

build:multiarch:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker run --privileged --rm tonistiigi/binfmt --install all
    - docker buildx create --use
  script:
    - docker buildx build --platform linux/amd64,linux/arm64 \
        -t $CI_REGISTRY_IMAGE:$CI_COMMIT_TAG \
        -t $CI_REGISTRY_IMAGE:latest \
        --push .
  only:
    - tags
```

### 使用Kaniko

```yaml
build:kaniko:
  stage: build
  image:
    name: gcr.io/kaniko-project/executor:latest
    entrypoint: [""]
  script:
    - mkdir -p /kaniko/.docker
    - echo "{\"auths\":{\"$CI_REGISTRY\":{\"auth\":\"$(echo -n $CI_REGISTRY_USER:$CI_REGISTRY_PASSWORD | base64)\"}}}" > /kaniko/.docker/config.json
    - /kaniko/executor \
        --context $CI_PROJECT_DIR \
        --dockerfile $CI_PROJECT_DIR/Dockerfile \
        --destination $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA \
        --destination $CI_REGISTRY_IMAGE:latest \
        --cache=true \
        --cache-ttl=24h
```

## 8.5 🔐 安全扫描

### 容器扫描

```yaml
container_scanning:
  stage: security
  image: registry.gitlab.com/security-products/container-scanning:latest
  variables:
    CI_APPLICATION_REPOSITORY: $CI_REGISTRY_IMAGE
    CI_APPLICATION_TAG: $CI_COMMIT_SHORT_SHA
  script:
    - gtcs scan
  artifacts:
    reports:
      container_scanning: gl-container-scanning-report.json

dependency_scanning:
  stage: security
  image: registry.gitlab.com/security-products/dependency-scanning:latest
  script:
    - /analyzer run
  artifacts:
    reports:
      dependency_scanning: gl-dependency-scanning-report.json

sast:
  stage: security
  image: registry.gitlab.com/security-products/sast:latest
  script:
    - /analyzer run
  artifacts:
    reports:
      sast: gl-sast-report.json
```

### Trivy扫描

```yaml
trivy:scan:
  stage: security
  image:
    name: aquasec/trivy:latest
    entrypoint: [""]
  script:
    - trivy image --exit-code 0 --severity LOW,MEDIUM $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA
    - trivy image --exit-code 1 --severity HIGH,CRITICAL $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA
```

## 8.6 🚀 部署阶段

### Kubernetes部署

```yaml
deploy:staging:
  stage: deploy
  image: bitnami/kubectl:latest
  environment:
    name: staging
    url: https://staging.example.com
  script:
    - kubectl config use-context staging-cluster
    - kubectl set image deployment/user-service \
        user-service=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA \
        -n staging
    - kubectl rollout status deployment/user-service -n staging
  only:
    - develop

deploy:production:
  stage: deploy
  image: bitnami/kubectl:latest
  environment:
    name: production
    url: https://example.com
  script:
    - kubectl config use-context production-cluster
    - kubectl set image deployment/user-service \
        user-service=$CI_REGISTRY_IMAGE:$CI_COMMIT_TAG \
        -n production
    - kubectl rollout status deployment/user-service -n production
  only:
    - tags
  when: manual
```

### Helm部署

```yaml
deploy:helm:
  stage: deploy
  image: alpine/helm:latest
  script:
    - helm upgrade --install user-service ./charts/user-service \
        --set image.tag=$CI_COMMIT_SHORT_SHA \
        --set image.repository=$CI_REGISTRY_IMAGE \
        --namespace production \
        --create-namespace \
        --wait \
        --timeout 5m
  environment:
    name: production
    kubernetes:
      namespace: production
  only:
    - main
```

## 8.7 📊 高级特性

### 并行执行

```yaml
test:parallel:
  stage: test
  parallel: 4
  script:
    - go test -v ./... -parallel 4

deploy:regions:
  stage: deploy
  parallel:
    matrix:
      - REGION: [us-east-1, eu-west-1, ap-southeast-1]
  script:
    - echo "Deploying to $REGION"
    - kubectl --context=$REGION deploy...
```

### 动态子Pipeline

```yaml
generate:pipeline:
  stage: build
  script:
    - ./scripts/generate-pipeline.sh > generated-pipeline.yml
  artifacts:
    paths:
      - generated-pipeline.yml

trigger:pipeline:
  stage: deploy
  trigger:
    include:
      - artifact: generated-pipeline.yml
        job: generate:pipeline
    strategy: depend
```

### 自动回滚

```yaml
deploy:production:
  stage: deploy
  script:
    - ./deploy.sh
  environment:
    name: production
    on_stop: rollback:production

rollback:production:
  stage: deploy
  script:
    - kubectl rollout undo deployment/user-service -n production
  environment:
    name: production
    action: stop
  when: manual
```

## 8.8 🎯 最佳实践

1. **使用模板**: 创建可复用的`.gitlab-ci.yml`模板
2. **缓存依赖**: 使用`cache`加速构建
3. **并行执行**: 使用`parallel`加速测试
4. **环境保护**: 为生产环境配置审批
5. **Secret管理**: 使用GitLab CI/CD变量
6. **Artifact管理**: 合理设置过期时间
7. **Runner优化**: 使用私有Runner提高速度
8. **Pipeline效率**: 避免不必要的job
9. **监控Pipeline**: 使用GitLab Insights
10. **文档化**: 维护Pipeline文档

## 8.9 ⚠️ 常见问题

### Q1: 如何加速GitLab CI？

**A**:

- 使用私有Runner
- 缓存依赖
- 并行执行
- Docker Layer缓存

### Q2: 如何调试Pipeline？

**A**:

```yaml
debug:
  stage: test
  script:
    - echo "Debugging..."
    - env
    - ls -la
  when: manual
```

### Q3: 多项目如何共享配置？

**A**: 使用`include`:

```yaml
include:
  - project: 'group/ci-templates'
    ref: main
    file: '/go-pipeline.yml'
```

## 8.10 📚 扩展阅读

### 官方文档

- [GitLab CI/CD文档](https://docs.gitlab.com/ee/ci/)
- [.gitlab-ci.yml参考](https://docs.gitlab.com/ee/ci/yaml/)
- [GitLab Runner](https://docs.gitlab.com/runner/)

### 相关文档

- [06-GitOps部署.md](./06-GitOps部署.md)
- [07-GitHub-Actions.md](./07-GitHub-Actions.md)
- [09-多环境部署.md](./09-多环境部署.md)

---

**文档维护者**: Go Documentation Team
**最后更新**: 2025-10-29
**文档状态**: 完成
**适用版本**: GitLab CE/EE 16+, Go 1.21+, Kubernetes 1.27+
