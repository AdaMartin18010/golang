# Kubernetes Operators

> **分类**: 成熟应用领域

---

## 什么是 Operator

Operator 使用自定义资源定义 (CRD) 来扩展 Kubernetes API，实现复杂应用的自动化运维。

---

## 核心组件

```
┌─────────────────┐
│   Custom Resource │
│   (MyDatabase)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Controller   │
│  (Reconcile)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Kubernetes API │
└─────────────────┘
```

---

## Controller-runtime

```go
import (
    "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

type DatabaseReconciler struct {
    client.Client
}

func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    var db mygroupv1.Database
    if err := r.Get(ctx, req.NamespacedName, &db); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // 创建 Deployment
    deploy := r.desiredDeployment(db)
    if err := r.Create(ctx, deploy); err != nil {
        return ctrl.Result{}, err
    }

    // 创建 Service
    svc := r.desiredService(db)
    if err := r.Create(ctx, svc); err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}
```

---

## Kubebuilder

```bash
# 初始化项目
kubebuilder init --domain example.com --repo github.com/user/project

# 创建 API
kubebuilder create api --group database --version v1 --kind Database

# 构建并部署
make manifests
make install
make run
```
