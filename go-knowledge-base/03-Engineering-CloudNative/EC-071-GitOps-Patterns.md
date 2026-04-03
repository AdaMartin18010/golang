# GitOps Patterns

> **分类**: 工程与云原生
> **标签**: #gitops #argocd #flux #cicd #declarative
> **参考**: GitOps Principles, ArgoCD, Flux CD, Weaveworks

---

## 1. Formal Definition

### 1.1 What is GitOps?

GitOps is an operational framework that takes DevOps best practices used for application development (version control, collaboration, compliance) and applies them to infrastructure automation. It uses Git repositories as the single source of truth for declarative infrastructure and application configurations.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         GitOps Architecture                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                          GIT REPOSITORY                            │   │
│   │  (Single Source of Truth)                                           │   │
│   │                                                                     │   │
│   │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │   │
│   │  │  Application │  │  Infrastructure│  │   Policy    │              │   │
│   │  │   Configs    │  │    as Code     │  │   Rules     │              │   │
│   │  │              │  │                │  │             │              │   │
│   │  │ • Deployments│  │ • Terraform    │  │ • Security  │              │   │
│   │  │ • Services   │  │ • Ansible      │  │ • Compliance│              │   │
│   │  │ • ConfigMaps │  │ • CloudFormation│  │ • Cost      │              │   │
│   │  │ • Secrets    │  │ • Pulumi       │  │   Controls  │              │   │
│   │  └──────────────┘  └──────────────┘  └──────────────┘              │   │
│   │                                                                     │   │
│   │  Branches: main (production) ◄── staging ◄── development            │   │
│   │                                                                     │   │
│   └────────────────────────┬────────────────────────────────────────────┘   │
│                            │                                                │
│                            │ Git Push / PR Merge                            │
│                            │                                                │
│                            ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                     GITOPS CONTROLLER                              │   │
│   │                                                                     │   │
│   │   ┌─────────────────────────────────────────────────────────────┐   │   │
│   │   │                     RECONCILIATION LOOP                     │   │   │
│   │   │                                                             │   │   │
│   │   │  1. PULL ──► 2. COMPARE ──► 3. DETECT DRIFT ──► 4. APPLY   │   │   │
│   │   │                                                             │   │   │
│   │   │  • Watch Git repo    • Current state   • Differences   •    │   │   │
│   │   │  • Poll for changes  • vs desired      • detected      •    │   │   │
│   │   │                                                             │   │   │
│   │   └─────────────────────────────────────────────────────────────┘   │   │
│   │                                                                     │   │
│   │   Capabilities:                                                     │   │
│   │   • Automated sync                                                  │   │
│   │   • Self-healing (correct drift)                                    │   │
│   │   • Progressive delivery (canary, blue/green)                       │   │
│   │   • RBAC and multi-tenancy                                          │   │
│   │   • Observability and audit                                         │   │
│   │                                                                     │   │
│   └────────────────────────┬────────────────────────────────────────────┘   │
│                            │                                                │
│                            │ Apply / Sync                                   │
│                            │                                                │
│                            ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                     TARGET ENVIRONMENT                             │   │
│   │                                                                     │   │
│   │   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐           │   │
│   │   │ Cluster 1│  │ Cluster 2│  │ Cluster 3│  │ Cluster N│           │   │
│   │   │          │  │          │  │          │  │          │           │   │
│   │   │ Production│ │ Staging  │  │  Dev     │  │  DR      │           │   │
│   │   └──────────┘  └──────────┘  └──────────┘  └──────────┘           │   │
│   │                                                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│   KEY PRINCIPLES:                                                           │
│   1. Declarative: System defined by configuration in Git                    │
│   2. Versioned & Immutable: Git history provides audit trail                │
│   3. Pulled Automatically: Agents automatically apply changes               │
│   4. Continuously Reconciled: Drift detection and correction                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 GitOps Deployment Patterns

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       GitOps Deployment Patterns                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  PATTERN 1: MONO-REPO                                                       │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━                                               │
│                                                                             │
│  Repository Structure:                                                      │
│  gitops-repo/                                                               │
│  ├── apps/                                                                  │
│  │   ├── frontend/                                                          │
│  │   │   ├── base/                                                          │
│  │   │   └── overlays/                                                      │
│  │   │       ├── dev/                                                       │
│  │   │       ├── staging/                                                   │
│  │   │       └── prod/                                                      │
│  │   └── backend/                                                           │
│  │       └── ...                                                            │
│  ├── infrastructure/                                                        │
│  │   ├── networking/                                                        │
│  │   ├── storage/                                                           │
│  │   └── monitoring/                                                        │
│  └── policies/                                                              │
│                                                                             │
│  Pros: Single source of truth, atomic changes                               │
│  Cons: Scale challenges, blast radius, access control complexity            │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  PATTERN 2: REPO PER APP (Application-Centric)                              │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━                                │
│                                                                             │
│  Repository Structure:                                                      │
│  frontend-repo/                                                             │
│  ├── src/                                                                   │
│  ├── Dockerfile                                                             │
│  └── k8s/           ◄── Application manifests                               │
│      ├── deployment.yaml                                                    │
│      └── service.yaml                                                       │
│                                                                             │
│  backend-repo/                                                              │
│  ├── src/                                                                   │
│  ├── Dockerfile                                                             │
│  └── k8s/                                                                   │
│                                                                             │
│  gitops-config-repo/   ◄── Environment configurations                       │
│  ├── environments/                                                          │
│  │   ├── dev/                                                               │
│  │   ├── staging/                                                           │
│  │   └── prod/                                                              │
│  └── apps/          ◄── References to app repos                             │
│                                                                             │
│  Pros: Clear ownership, independent releases                                │
│  Cons: Cross-app changes span multiple repos                                │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  PATTERN 3: REPO PER ENVIRONMENT (Environment-Centric)                      │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━                           │
│                                                                             │
│  Repository Structure:                                                      │
│  gitops-dev/                                                                │
│  ├── apps/                                                                  │
│  ├── infrastructure/                                                        │
│  └── policies/                                                              │
│                                                                             │
│  gitops-staging/                                                            │
│  └── ...                                                                    │
│                                                                             │
│  gitops-production/                                                         │
│  └── ...                                                                    │
│                                                                             │
│  Pros: Strong isolation, environment-specific access control                │
│  Cons: Config duplication, promotion complexity                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns

### 2.1 ArgoCD Application Configuration

```yaml
# Application definition
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: myapp-production
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
  annotations:
    argocd.argoproj.io/sync-wave: "2"
spec:
  project: production
  source:
    repoURL: https://github.com/example/gitops-repo.git
    targetRevision: main
    path: apps/myapp/overlays/production
    helm:
      valueFiles:
        - values-production.yaml
    kustomize:
      namePrefix: prod-
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
      - PrunePropagationPolicy=foreground
      - PruneLast=true
      - Validate=true
      - ApplyOutOfSyncOnly=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
  ignoreDifferences:
    - group: apps
      kind: Deployment
      jsonPointers:
        - /spec/replicas
  revisionHistoryLimit: 10
```

```yaml
# ApplicationSet for multi-cluster
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: myapp
  namespace: argocd
spec:
  generators:
    - list:
        elements:
          - cluster: dev
            url: https://dev-cluster.example.com
            namespace: development
            valuesFile: values-dev.yaml
          - cluster: staging
            url: https://staging-cluster.example.com
            namespace: staging
            valuesFile: values-staging.yaml
          - cluster: prod
            url: https://prod-cluster.example.com
            namespace: production
            valuesFile: values-prod.yaml
  template:
    metadata:
      name: 'myapp-{{cluster}}'
    spec:
      project: default
      source:
        repoURL: https://github.com/example/gitops-repo.git
        targetRevision: HEAD
        path: apps/myapp
        helm:
          valueFiles:
            - '{{valuesFile}}'
      destination:
        server: '{{url}}'
        namespace: '{{namespace}}'
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
```

### 2.2 Flux Configuration

```yaml
# GitRepository
apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: myapp
  namespace: flux-system
spec:
  interval: 1m
  url: https://github.com/example/gitops-repo
  ref:
    branch: main
  secretRef:
    name: github-token
  ignore:
    - "/*.md"
    - "/docs/**"
---
# Kustomization
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: myapp
  namespace: flux-system
spec:
  interval: 10m
  path: ./apps/myapp/overlays/production
  prune: true
  sourceRef:
    kind: GitRepository
    name: myapp
  healthChecks:
    - apiVersion: apps/v1
      kind: Deployment
      name: myapp
      namespace: production
  timeout: 2m
  retryInterval: 1m
  wait: true
---
# ImageRepository
apiVersion: image.toolkit.fluxcd.io/v1beta2
kind: ImageRepository
metadata:
  name: myapp
  namespace: flux-system
spec:
  image: ghcr.io/example/myapp
  interval: 5m
  secretRef:
    name: ghcr-auth
---
# ImagePolicy
apiVersion: image.toolkit.fluxcd.io/v1beta2
kind: ImagePolicy
metadata:
  name: myapp
  namespace: flux-system
spec:
  imageRepositoryRef:
    name: myapp
  policy:
    semver:
      range: "1.x.x"
  filterTags:
    pattern: '^v(?P<version>.*)$'
    extract: '$version'
---
# ImageUpdateAutomation
apiVersion: image.toolkit.fluxcd.io/v1beta1
kind: ImageUpdateAutomation
metadata:
  name: myapp
  namespace: flux-system
spec:
  interval: 1m
  sourceRef:
    kind: GitRepository
    name: myapp
  git:
    checkout:
      ref:
        branch: main
    commit:
      author:
        name: Flux Bot
        email: flux@example.com
      messageTemplate: |
        Automated image update

        Images:
        {{ range .Updated.Images -}}
        - {{.}}
        {{ end }}
      signingKey:
        secretRef:
          name: flux-gpg-signing-key
    push:
      branch: main
  policy:
    semver:
      range: "1.x.x"
```

### 2.3 Progressive Delivery

```yaml
# Argo Rollout - Canary Deployment
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: myapp
  namespace: production
spec:
  replicas: 10
  strategy:
    canary:
      canaryService: myapp-canary
      stableService: myapp-stable
      trafficRouting:
        nginx:
          stableIngress: myapp-ingress
          annotationPrefix: nginx.ingress.kubernetes.io
      steps:
        - setWeight: 5
        - pause:
            duration: 10m
        - setWeight: 20
        - pause:
            duration: 10m
        - analysis:
            templates:
              - templateName: success-rate
            args:
              - name: service
                value: myapp-canary
        - setWeight: 50
        - pause:
            duration: 10m
        - setWeight: 100
      analysis:
        startingStep: 2
        args:
          - name: service
            value: myapp-canary
        templates:
          - templateName: success-rate
---
# Analysis Template
apiVersion: argoproj.io/v1alpha1
kind: AnalysisTemplate
metadata:
  name: success-rate
spec:
  args:
    - name: service
  metrics:
    - name: success-rate
      interval: 5m
      count: 3
      successCondition: result[0] >= 0.95
      provider:
        prometheus:
          address: http://prometheus:9090
          query: |
            sum(rate(http_requests_total{service="{{args.service}}",status=~"2.."}[5m]))
            /
            sum(rate(http_requests_total{service="{{args.service}}"}[5m]))
```

---

## 3. Production-Ready Configurations

### 3.1 Multi-Tenancy ArgoCD Setup

```yaml
# AppProject for team isolation
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: team-alpha
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  description: Team Alpha Project
  sourceRepos:
    - 'https://github.com/team-alpha/*'
    - 'https://charts.team-alpha.io'
  destinations:
    - namespace: team-alpha-*
      server: https://kubernetes.default.svc
    - namespace: team-alpha-*
      server: https://prod-cluster.example.com
  clusterResourceWhitelist:
    - group: ''
      kind: Namespace
  namespaceResourceBlacklist:
    - group: ''
      kind: ResourceQuota
  roles:
    - name: admin
      description: Team Alpha Admins
      policies:
        - p, proj:team-alpha:admin, applications, *, team-alpha/*, allow
        - p, proj:team-alpha:admin, exec, create, team-alpha/*, allow
      groups:
        - team-alpha-admins
    - name: readonly
      description: Team Alpha Read Only
      policies:
        - p, proj:team-alpha:readonly, applications, get, team-alpha/*, allow
      groups:
        - team-alpha-members
  syncWindows:
    - kind: deny
      schedule: '0 23 * * *'
      duration: 8h
      namespaces:
        - team-alpha-production
---
# RBAC Configuration
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  policy.default: role:readonly
  policy.csv: |
    p, role:org-admin, applications, *, *, allow
    p, role:org-admin, clusters, get, *, allow
    p, role:org-admin, repositories, *, *, allow

    g, admin-group, role:org-admin
```

---

## 4. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       GitOps Security Best Practices                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  REPOSITORY SECURITY                                                        │
│  ✓ Branch protection rules                                                  │
│  ✓ Required code review (2+ reviewers for production)                       │
│  ✓ Signed commits (GPG)                                                     │
│  ✓ Secret scanning in CI                                                    │
│  ✓ No secrets in Git (use sealed-secrets, SOPS, Vault)                      │
│                                                                             │
│  ACCESS CONTROL                                                             │
│  ✓ RBAC for GitOps controller                                               │
│  ✓ Project isolation (ArgoCD)                                               │
│  ✓ Namespace boundaries                                                     │
│  ✓ Service account per application                                          │
│                                                                             │
│  NETWORK SECURITY                                                           │
│  ✓ TLS for Git communication                                                │
│  ✓ Network policies for GitOps agents                                       │
│  ✓ Private Git repositories                                                 │
│  ✓ Webhook validation                                                       │
│                                                                             │
│  AUDIT & COMPLIANCE                                                         │
│  ✓ All changes tracked in Git                                               │
│  ✓ Audit logs for GitOps operations                                         │
│  ✓ Drift detection alerts                                                   │
│  ✓ Policy enforcement (OPA, Kyverno)                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Decision Matrices

### 5.1 GitOps Tool Selection

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     GitOps Tool Comparison Matrix                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Feature            │  ArgoCD    │  Flux CD   │  Jenkins X │  Spinnaker   │
├─────────────────────┼────────────┼────────────┼────────────┼──────────────│
│  User Interface     │  ★★★★★     │  ★★★☆☆     │  ★★★☆☆     │  ★★★★★       │
│  Multi-cluster      │  ★★★★★     │  ★★★★☆     │  ★★★☆☆     │  ★★★★★       │
│  Multi-tenancy      │  ★★★★★     │  ★★★☆☆     │  ★★☆☆☆     │  ★★★★☆       │
│  Progressive        │  ★★★★★     │  ★★★☆☆     │  ★★★☆☆     │  ★★★★★       │
│  Delivery           │            │  (Flagger) │            │              │
│  Image Automation   │  ★★☆☆☆     │  ★★★★★     │  ★★★★☆     │  ★★★★☆       │
│  Helm Support       │  ★★★★★     │  ★★★★★     │  ★★★★☆     │  ★★★★☆       │
│  Kustomize          │  ★★★★★     │  ★★★★★     │  ★★★☆☆     │  ★★★☆☆       │
│  Secrets Mgmt       │  ★★★☆☆     │  ★★★★☆     │  ★★★☆☆     │  ★★★★☆       │
│  Learning Curve     │  Medium    │  Medium    │  Steep     │  Steep       │
│  CNCF Status        │  Graduated │  Graduated │  Sandbox   │  -           │
│                                                                             │
│  Recommendation:                                                            │
│  • UI-first, enterprise: ArgoCD                                             │
│  • Git-native, image automation: Flux CD                                    │
│  • Progressive delivery focus: ArgoCD with Rollouts                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        GitOps Best Practices Summary                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  REPOSITORY STRUCTURE                                                       │
│  ✓ Clear directory structure (apps/, infrastructure/, policies/)            │
│  ✓ Environment separation (dev/, staging/, prod/)                          │
│  ✓ Kustomize bases and overlays for DRY                                     │
│  ✓ Consistent naming conventions                                            │
│                                                                             │
│  DEPLOYMENT                                                                 │
│  ✓ Automated sync with manual approval for prod                             │
│  ✓ Self-healing enabled                                                     │
│  ✓ Pruning of removed resources                                             │
│  ✓ Health checks and dependency management                                  │
│                                                                             │
│  SECURITY                                                                   │
│  ✓ RBAC for GitOps controller                                               │
│  ✓ Project/namespace isolation                                              │
│  ✓ Secret management outside Git                                            │
│  ✓ Policy enforcement (OPA, Kyverno)                                        │
│                                                                             │
│  OBSERVABILITY                                                              │
│  ✓ GitOps metrics and dashboards                                            │
│  ✓ Drift detection and alerts                                               │
│  ✓ Sync status monitoring                                                   │
│  ✓ Audit logging                                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. GitOps Working Group
2. ArgoCD Documentation
3. Flux CD Documentation
4. OpenGitOps
5. Kubernetes GitOps Best Practices
