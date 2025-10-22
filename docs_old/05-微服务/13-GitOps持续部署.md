# 13. 🔄 GitOps持续部署

> 📚 **简介**：本文档深入探讨GitOps理念及其在微服务持续部署中的应用，涵盖ArgoCD、Flux等主流工具的使用，以及如何构建自动化、可追溯、安全的部署流程。通过本文，读者将掌握现代化的GitOps实践方法。

<!-- TOC START -->
- [13. 🔄 GitOps持续部署](#13--gitops持续部署)
  - [13.1 📚 GitOps概述](#131--gitops概述)
  - [13.2 ⚙️ ArgoCD实践](#132-️-argocd实践)
    - [安装ArgoCD](#安装argocd)
    - [创建应用](#创建应用)
    - [Helm集成](#helm集成)
    - [多环境管理](#多环境管理)
  - [13.3 🌊 Flux CD实践](#133--flux-cd实践)
    - [安装Flux](#安装flux)
    - [GitRepository源](#gitrepository源)
    - [Kustomization部署](#kustomization部署)
    - [HelmRelease](#helmrelease)
  - [13.4 📋 部署策略](#134--部署策略)
    - [蓝绿部署](#蓝绿部署)
    - [金丝雀发布（Flagger）](#金丝雀发布flagger)
    - [渐进式交付（Argo Rollouts）](#渐进式交付argo-rollouts)
  - [13.5 🔐 安全最佳实践](#135--安全最佳实践)
    - [敏感信息管理](#敏感信息管理)
    - [RBAC配置](#rbac配置)
  - [13.6 📊 监控与告警](#136--监控与告警)
    - [Prometheus监控](#prometheus监控)
    - [通知集成](#通知集成)
  - [13.7 💻 实战案例](#137--实战案例)
    - [CI/CD完整流程](#cicd完整流程)
    - [多集群部署](#多集群部署)
  - [13.8 🎯 最佳实践](#138--最佳实践)
  - [13.9 ⚠️ 常见问题](#139-️-常见问题)
    - [Q1: ArgoCD应用一直OutOfSync？](#q1-argocd应用一直outofsync)
    - [Q2: 如何处理Helm Chart版本冲突？](#q2-如何处理helm-chart版本冲突)
    - [Q3: 如何回滚部署？](#q3-如何回滚部署)
    - [Q4: 多环境配置如何管理？](#q4-多环境配置如何管理)
  - [13.10 📚 扩展阅读](#1310--扩展阅读)
    - [官方文档](#官方文档)
    - [相关文档](#相关文档)
<!-- TOC END -->

## 13.1 📚 GitOps概述

**GitOps**: 以Git作为单一事实来源，通过声明式配置和自动化实现基础设施和应用的持续交付。

**核心原则**:

1. **声明式**: 所有配置以声明式方式存储在Git中
2. **版本化**: Git提供完整的变更历史和回滚能力
3. **自动化**: 系统自动同步Git状态到集群
4. **持续协调**: 持续监控并修正配置漂移

**传统部署 vs GitOps**:

| 维度 | 传统部署 | GitOps |
|------|---------|--------|
| 配置管理 | 手动/脚本 | Git仓库 |
| 部署触发 | CI管道推送 | CD系统拉取 |
| 状态管理 | 难以追踪 | Git历史 |
| 回滚 | 复杂 | Git revert |
| 审计 | 困难 | Git log |

## 13.2 ⚙️ ArgoCD实践

### 安装ArgoCD

```bash
# 1. 创建命名空间
kubectl create namespace argocd

# 2. 安装ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# 3. 暴露UI（LoadBalancer）
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'

# 4. 获取初始密码
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# 5. 登录CLI
argocd login <ARGOCD_SERVER>
argocd account update-password
```

### 创建应用

**通过CLI**:

```bash
argocd app create user-service \
  --repo https://github.com/myorg/k8s-manifests.git \
  --path apps/user-service \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace production \
  --sync-policy automated \
  --auto-prune \
  --self-heal
```

**通过YAML**:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: user-service
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/myorg/k8s-manifests.git
    targetRevision: main
    path: apps/user-service
  destination:
    server: https://kubernetes.default.svc
    namespace: production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
    - CreateNamespace=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

### Helm集成

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: user-service-helm
spec:
  source:
    repoURL: https://github.com/myorg/helm-charts.git
    targetRevision: main
    path: user-service
    helm:
      valueFiles:
      - values-production.yaml
      parameters:
      - name: image.tag
        value: v1.2.3
      - name: replicaCount
        value: "3"
```

### 多环境管理

```text
manifests/
├── base/                   # 基础配置
│   ├── deployment.yaml
│   ├── service.yaml
│   └── kustomization.yaml
├── overlays/
│   ├── dev/               # 开发环境
│   │   ├── kustomization.yaml
│   │   └── patch.yaml
│   ├── staging/           # 预发布环境
│   │   ├── kustomization.yaml
│   │   └── patch.yaml
│   └── production/        # 生产环境
│       ├── kustomization.yaml
│       └── patch.yaml
```

**Kustomization**:

```yaml
# base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- deployment.yaml
- service.yaml

---
# overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
patchesStrategicMerge:
- patch.yaml
replicas:
- name: user-service
  count: 5
images:
- name: myregistry/user-service
  newTag: v1.2.3
```

## 13.3 🌊 Flux CD实践

### 安装Flux

```bash
# 1. 安装Flux CLI
curl -s https://fluxcd.io/install.sh | sudo bash

# 2. 引导Flux
flux bootstrap github \
  --owner=myorg \
  --repository=k8s-cluster \
  --branch=main \
  --path=clusters/production \
  --personal
```

### GitRepository源

```yaml
apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: user-service
  namespace: flux-system
spec:
  interval: 1m
  url: https://github.com/myorg/user-service-manifests
  ref:
    branch: main
  secretRef:
    name: git-credentials
```

### Kustomization部署

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind:Kustomization
metadata:
  name: user-service
  namespace: flux-system
spec:
  interval: 5m
  path: ./kustomize/production
  prune: true
  sourceRef:
    kind: GitRepository
    name: user-service
  healthChecks:
  - apiVersion: apps/v1
    kind: Deployment
    name: user-service
    namespace: production
```

### HelmRelease

```yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: user-service
  namespace: production
spec:
  interval: 5m
  chart:
    spec:
      chart: user-service
      version: "1.2.3"
      sourceRef:
        kind: HelmRepository
        name: myorg
        namespace: flux-system
  values:
    replicaCount: 3
    image:
      tag: v1.2.3
```

## 13.4 📋 部署策略

### 蓝绿部署

```yaml
# 蓝色环境（当前）
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service-blue
  labels:
    app: user-service
    version: blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
      version: blue
  template:
    metadata:
      labels:
        app: user-service
        version: blue
    spec:
      containers:
      - name: user-service
        image: myregistry/user-service:v1

---
# 绿色环境（新版本）
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service-green
  labels:
    app: user-service
    version: green
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
      version: green
  template:
    metadata:
      labels:
        app: user-service
        version: green
    spec:
      containers:
      - name: user-service
        image: myregistry/user-service:v2

---
# Service切换
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
    version: blue  # 切换到green实现蓝绿切换
  ports:
  - port: 80
    targetPort: 8080
```

### 金丝雀发布（Flagger）

```yaml
apiVersion: flagger.app/v1beta1
kind: Canary
metadata:
  name: user-service
  namespace: production
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  service:
    port: 80
    targetPort: 8080
  analysis:
    interval: 1m
    threshold: 5
    maxWeight: 50
    stepWeight: 10
    metrics:
    - name: request-success-rate
      thresholdRange:
        min: 99
      interval: 1m
    - name: request-duration
      thresholdRange:
        max: 500
      interval: 1m
  webhooks:
  - name: load-test
    url: http://flagger-loadtester/
    timeout: 5s
    metadata:
      cmd: "hey -z 1m -q 10 -c 2 http://user-service-canary"
```

### 渐进式交付（Argo Rollouts）

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: user-service
spec:
  replicas: 5
  strategy:
    canary:
      steps:
      - setWeight: 20
      - pause: {duration: 5m}
      - setWeight: 40
      - pause: {duration: 5m}
      - setWeight: 60
      - pause: {duration: 5m}
      - setWeight: 80
      - pause: {duration: 5m}
      canaryService: user-service-canary
      stableService: user-service
      trafficRouting:
        istio:
          virtualService:
            name: user-service
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: myregistry/user-service:v2
```

## 13.5 🔐 安全最佳实践

### 敏感信息管理

**使用Sealed Secrets**:

```bash
# 1. 安装Sealed Secrets
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# 2. 创建加密Secret
echo -n "password123" | kubectl create secret generic db-secret \
  --dry-run=client \
  --from-file=password=/dev/stdin \
  -o yaml | \
  kubeseal -o yaml > sealed-secret.yaml

# 3. 提交到Git
git add sealed-secret.yaml
git commit -m "Add sealed secret"
```

**使用External Secrets Operator**:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: db-credentials
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: db-secret
  data:
  - secretKey: password
    remoteRef:
      key: prod/db/password
```

### RBAC配置

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: argocd-deployer
  namespace: production
rules:
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: [""]
  resources: ["services", "configmaps"]
  verbs: ["get", "list", "watch", "create", "update", "patch"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: argocd-deployer-binding
  namespace: production
subjects:
- kind: ServiceAccount
  name: argocd-application-controller
  namespace: argocd
roleRef:
  kind: Role
  name: argocd-deployer
  apiGroup: rbac.authorization.k8s.io
```

## 13.6 📊 监控与告警

### Prometheus监控

```yaml
# ArgoCD指标
argocd_app_sync_total
argocd_app_health_status
argocd_app_sync_status
```

**告警规则**:

```yaml
groups:
- name: argocd
  rules:
  - alert: ArgoCDAppOutOfSync
    expr: argocd_app_sync_status{sync_status!="Synced"} == 1
    for: 10m
    labels:
      severity: warning
    annotations:
      summary: "ArgoCD application {{ $labels.name }} is out of sync"
  
  - alert: ArgoCDAppUnhealthy
    expr: argocd_app_health_status{health_status!="Healthy"} == 1
    for: 15m
    labels:
      severity: critical
    annotations:
      summary: "ArgoCD application {{ $labels.name }} is unhealthy"
```

### 通知集成

```yaml
# ArgoCD通知配置
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-notifications-cm
  namespace: argocd
data:
  service.slack: |
    token: $slack-token
  template.app-deployed: |
    message: |
      Application {{.app.metadata.name}} is now running new version.
  trigger.on-deployed: |
    - description: Application is synced and healthy
      send:
      - app-deployed
      when: app.status.operationState.phase in ['Succeeded'] and app.status.health.status == 'Healthy'
```

## 13.7 💻 实战案例

### CI/CD完整流程

```yaml
# GitHub Actions
name: CI/CD Pipeline
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Build Docker Image
      run: |
        docker build -t myregistry/user-service:${{ github.sha }} .
        docker push myregistry/user-service:${{ github.sha }}
    
    - name: Update Manifest
      run: |
        git clone https://github.com/myorg/k8s-manifests.git
        cd k8s-manifests
        yq eval '.spec.template.spec.containers[0].image = "myregistry/user-service:${{ github.sha }}"' \
          -i apps/user-service/deployment.yaml
        git commit -am "Update image to ${{ github.sha }}"
        git push
```

### 多集群部署

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: user-service-multi-cluster
spec:
  generators:
  - list:
      elements:
      - cluster: us-east-1
        url: https://cluster-us-east-1
        namespace: production
      - cluster: eu-west-1
        url: https://cluster-eu-west-1
        namespace: production
  template:
    metadata:
      name: 'user-service-{{cluster}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/myorg/k8s-manifests.git
        targetRevision: main
        path: apps/user-service
      destination:
        server: '{{url}}'
        namespace: '{{namespace}}'
      syncPolicy:
        automated: {}
```

## 13.8 🎯 最佳实践

1. **Git为单一事实来源**: 所有配置变更必须通过Git
2. **环境分支策略**: dev/staging/production分支或目录
3. **自动化同步**: 启用auto-sync和self-heal
4. **渐进式部署**: 使用金丝雀或蓝绿部署
5. **监控同步状态**: 设置告警监控部署健康
6. **敏感信息加密**: 使用Sealed Secrets或External Secrets
7. **RBAC权限控制**: 最小权限原则
8. **版本标记**: 使用语义化版本标记
9. **回滚计划**: 准备快速回滚方案
10. **文档化**: 维护清晰的部署文档

## 13.9 ⚠️ 常见问题

### Q1: ArgoCD应用一直OutOfSync？

**A**: 检查：

```bash
# 查看差异
argocd app diff user-service

# 忽略特定字段
argocd app set user-service --ignore-difference 'group=apps,kind=Deployment,jsonPointers=/spec/replicas'
```

### Q2: 如何处理Helm Chart版本冲突？

**A**: 使用`targetRevision`固定版本：

```yaml
source:
  chart: user-service
  targetRevision: "1.2.3"  # 固定版本
```

### Q3: 如何回滚部署？

**A**:

```bash
# 查看历史
argocd app history user-service

# 回滚到指定版本
argocd app rollback user-service <revision>

# Git回滚
git revert <commit>
git push
```

### Q4: 多环境配置如何管理？

**A**: 使用Kustomize overlays或Helm values文件。

## 13.10 📚 扩展阅读

### 官方文档

- [ArgoCD文档](https://argo-cd.readthedocs.io/)
- [Flux CD文档](https://fluxcd.io/docs/)
- [GitOps Toolkit](https://toolkit.fluxcd.io/)

### 相关文档

- [11-Kubernetes微服务部署.md](./11-Kubernetes微服务部署.md)
- [12-Service Mesh集成.md](./12-Service-Mesh集成.md)
- [../06-云原生/07-GitHub Actions.md](../06-云原生/07-GitHub-Actions.md)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: ArgoCD 2.9+, Flux 2.2+, Kubernetes 1.27+
