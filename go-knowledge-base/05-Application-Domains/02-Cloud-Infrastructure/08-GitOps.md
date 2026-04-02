# GitOps 实践

> **分类**: 成熟应用领域
> **标签**: #gitops #argocd #flux

---

## GitOps 原则

1. **声明式**: 系统状态声明在 Git 中
2. **版本化**: Git 作为唯一事实来源
3. **自动同步**: 自动应用 Git 中的变更
4. **回滚**: 通过 Git 回滚

---

## Argo CD 集成

### Application 定义

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/org/repo.git
    targetRevision: HEAD
    path: k8s/overlays/production
  destination:
    server: https://kubernetes.default.svc
    namespace: production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

### Go 客户端

```go
import "github.com/argoproj/argo-cd/v2/pkg/apiclient"

client, err := apiclient.NewClient(&apiclient.ClientOptions{
    ServerAddr: "localhost:8080",
    AuthToken:  token,
})

// 创建应用
app, err := client.Create(context.Background(), &application.ApplicationCreateRequest{
    Application: &v1alpha1.Application{
        ObjectMeta: metav1.ObjectMeta{
            Name: "my-app",
        },
        Spec: v1alpha1.ApplicationSpec{
            Source: v1alpha1.ApplicationSource{
                RepoURL:        "https://github.com/org/repo",
                TargetRevision: "HEAD",
                Path:           "k8s/",
            },
        },
    },
})
```

---

## 结构

```
repo/
├── apps/
│   ├── my-app/
│   │   ├── base/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   └── kustomization.yaml
│   │   └── overlays/
│   │       ├── dev/
│   │       │   └── kustomization.yaml
│   │       └── prod/
│   │           └── kustomization.yaml
```

---

## 镜像更新自动化

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: my-apps
spec:
  generators:
  - git:
      repoURL: https://github.com/org/repo.git
      directories:
      - path: apps/*
  template:
    spec:
      source:
        repoURL: https://github.com/org/repo.git
        targetRevision: HEAD
```

---

## 健康检查

```go
// 自定义健康检查
func HealthCheck(ctx context.Context, app *v1alpha1.Application) error {
    if app.Status.Sync.Status != v1alpha1.SyncStatusCodeSynced {
        return fmt.Errorf("not synced: %s", app.Status.Sync.Status)
    }

    if app.Status.Health.Status != health.HealthStatusHealthy {
        return fmt.Errorf("not healthy: %s", app.Status.Health.Status)
    }

    return nil
}
```
