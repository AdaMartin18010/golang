# TS-005: Kubernetes Operator 模式 (K8s Operator Patterns)

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #kubernetes #operator #controller #crd
> **权威来源**: [Operator SDK](https://sdk.operatorframework.io/), [K8s Controller Concepts](https://kubernetes.io/docs/concepts/architecture/controller/)

---

## Operator 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Kubernetes Operator Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          Operator Pod                               │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                     Controller Manager                         │  │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │  │   │
│  │  │  │  Reconciler │  │   Watcher   │  │   Worker    │           │  │   │
│  │  │  │             │  │             │  │    Queue    │           │  │   │
│  │  │  │ - Compare   │  │ - Watch CR  │  │ - Rate      │           │  │   │
│  │  │  │ - Diff      │  │ - Enqueue   │  │   Limiter   │           │  │   │
│  │  │  │ - Apply     │  │ - Filter    │  │ - Retry     │           │  │   │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘           │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                              │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                       Client-Go                               │  │   │
│  │  │  - ListWatcher  - Informer  - WorkQueue                       │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                               │
│                              │ Watch/Update                                  │
│                              ▼                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                       Kubernetes API Server                         │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │   │
│  │  │   CustomResource │  │   Deployment    │  │    Service      │      │   │
│  │  │   (MyDatabase)   │  │                 │  │                 │      │   │
│  │  │                  │  │                 │  │                 │      │   │
│  │  │  spec:           │  │  spec:          │  │  spec:          │      │   │
│  │  │    replicas: 3   │  │    replicas: 3  │  │    ports:       │      │   │
│  │  │    storage: 100G │  │                 │  │                 │      │   │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────┘      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## CRD 定义

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: databases.example.com
spec:
  group: example.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                replicas:
                  type: integer
                  minimum: 1
                  maximum: 10
                storage:
                  type: string
                  pattern: '^[0-9]+(Gi|Mi)$'
                version:
                  type: string
                  enum: ["13", "14", "15"]
            status:
              type: object
              properties:
                phase:
                  type: string
                  enum: ["Pending", "Creating", "Running", "Failed"]
                readyReplicas:
                  type: integer
  scope: Namespaced
  names:
    plural: databases
    singular: database
    kind: Database
    shortNames:
      - db
```

---

## Go Controller 实现

```go
package controller

import (
    "context"
    "fmt"
    "time"

    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    examplev1 "github.com/example/api/v1"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=databases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    // 1. 获取 CR
    db := &examplev1.Database{}
    if err := r.Get(ctx, req.NamespacedName, db); err != nil {
        if errors.IsNotFound(err) {
            return ctrl.Result{}, nil // 已删除
        }
        return ctrl.Result{}, err
    }

    // 2. 创建/更新 Deployment
    if err := r.reconcileDeployment(ctx, db); err != nil {
        r.updateStatus(ctx, db, "Failed", err.Error())
        return ctrl.Result{RequeueAfter: 30 * time.Second}, err
    }

    // 3. 创建/更新 Service
    if err := r.reconcileService(ctx, db); err != nil {
        return ctrl.Result{}, err
    }

    // 4. 更新状态
    r.updateStatus(ctx, db, "Running", "")

    return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

func (r *DatabaseReconciler) reconcileDeployment(ctx context.Context, db *examplev1.Database) error {
    dep := &appsv1.Deployment{}
    depName := fmt.Sprintf("%s-db", db.Name)

    err := r.Get(ctx, client.ObjectKey{Name: depName, Namespace: db.Namespace}, dep)
    if err != nil && !errors.IsNotFound(err) {
        return err
    }

    // 创建新的 Deployment
    if errors.IsNotFound(err) {
        dep = &appsv1.Deployment{
            ObjectMeta: metav1.ObjectMeta{
                Name:      depName,
                Namespace: db.Namespace,
                OwnerReferences: []metav1.OwnerReference{
                    *metav1.NewControllerRef(db, examplev1.GroupVersion.WithKind("Database")),
                },
            },
            Spec: appsv1.DeploymentSpec{
                Replicas: &db.Spec.Replicas,
                Selector: &metav1.LabelSelector{
                    MatchLabels: map[string]string{"app": depName},
                },
                Template: corev1.PodTemplateSpec{
                    ObjectMeta: metav1.ObjectMeta{
                        Labels: map[string]string{"app": depName},
                    },
                    Spec: corev1.PodSpec{
                        Containers: []corev1.Container{{
                            Name:  "postgres",
                            Image: fmt.Sprintf("postgres:%s", db.Spec.Version),
                            Env: []corev1.EnvVar{
                                {Name: "POSTGRES_DB", Value: db.Name},
                            },
                            Resources: corev1.ResourceRequirements{
                                Requests: corev1.ResourceList{
                                    corev1.ResourceStorage: resource.MustParse(db.Spec.Storage),
                                },
                            },
                        }},
                    },
                },
            },
        }
        return r.Create(ctx, dep)
    }

    // 更新现有 Deployment
    if *dep.Spec.Replicas != db.Spec.Replicas {
        dep.Spec.Replicas = &db.Spec.Replicas
        return r.Update(ctx, dep)
    }

    return nil
}

func (r *DatabaseReconciler) reconcileService(ctx context.Context, db *examplev1.Database) error {
    // 类似逻辑创建 Service
    return nil
}

func (r *DatabaseReconciler) updateStatus(ctx context.Context, db *examplev1.Database, phase, message string) {
    db.Status.Phase = phase
    db.Status.Message = message
    if err := r.Status().Update(ctx, db); err != nil {
        log.FromContext(ctx).Error(err, "Failed to update status")
    }
}

func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplev1.Database{}).
        Owns(&appsv1.Deployment{}).
        Owns(&corev1.Service{}).
        Complete(r)
}
```

---

## 常用 Operator 框架

| 框架 | 特点 | 推荐场景 |
|------|------|---------|
| Operator SDK | Go, Ansible, Helm 支持 | 生产级 |
| Kubebuilder | Go, 官方推荐 | 复杂业务 |
| Helm Operator | 纯 Helm chart | 简单场景 |

---

## 参考文献

1. [Operator SDK](https://sdk.operatorframework.io/)
2. [Kubebuilder](https://book.kubebuilder.io/)
3. [Writing Controllers](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
