# Kubernetes Operators in Go

> **Dimension**: Application Domains
> **Level**: S (18+ KB)
> **Tags**: #kubernetes #operator #controller #crd #kubebuilder

---

## 1. Domain Requirements Analysis

### 1.1 Operator vs Manual Management

| Aspect | Manual Management | Operator Pattern |
|--------|-------------------|------------------|
| Deployment | kubectl apply | Automated, GitOps-ready |
| Scaling | Manual replica changes | HPA integration, custom metrics |
| Upgrades | Manual version updates | Rolling upgrades with verification |
| Backup | External scripts | Automated scheduled backups |
| Recovery | Manual intervention | Self-healing based on conditions |
| Monitoring | External setup | Built-in metrics exporters |
| Configuration | ConfigMap/Secret edits | Validated CRD changes |

### 1.2 When to Build an Operator

**Ideal Scenarios:**

- Stateful applications requiring complex lifecycle management
- Applications needing domain-specific upgrade procedures
- Services requiring automated backup/restore
- Multi-cluster application coordination
- Custom resource orchestration needs

**Avoid When:**

- Simple stateless deployments suffice
- Helm charts meet all requirements
- Team lacks Kubernetes deep expertise
- Application lifecycle is straightforward

---

## 2. Architecture Formalization

### 2.1 Kubernetes Operator Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Kubernetes Operator Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Custom Resource Definition (CRD)                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                                                                     │    │
│  │  apiVersion: apiextensions.k8s.io/v1                                │    │
│  │  kind: CustomResourceDefinition                                     │    │
│  │  metadata:                                                          │    │
│  │    name: databases.example.com                                      │    │
│  │  spec:                                                              │    │
│  │    group: example.com                                               │    │
│  │    names:                                                           │    │
│  │      kind: Database                                                 │    │
│  │      plural: databases                                              │    │
│  │    scope: Namespaced                                                │    │
│  │    versions:                                                        │    │
│  │    - name: v1                                                       │    │
│  │      schema: ...                                                    │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Custom Resource                             │    │
│  │                                                                     │    │
│  │  apiVersion: example.com/v1                                         │    │
│  │  kind: Database                                                     │    │
│  │  metadata:                                                          │    │
│  │    name: production-db                                              │    │
│  │  spec:                                                              │    │
│  │    version: "13"                                                    │    │
│  │    storage: 100Gi                                                   │    │
│  │    replicas: 3                                                      │    │
│  │    backup:                                                          │    │
│  │      enabled: true                                                  │    │
│  │      schedule: "0 2 * * *"                                          │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                       Operator Controller                           │    │
│  │                                                                     │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │    │
│  │  │   Watcher   │  │  Reconciler │  │   Finalizer │  │  Webhook   │ │    │
│  │  │             │  │             │  │             │  │            │ │    │
│  │  │ Watch CR    │  │ Compare     │  │ Cleanup     │  │ Validate   │ │    │
│  │  │ changes     │  │ desired vs  │  │ resources   │  │ & mutate   │ │    │
│  │  │             │  │ actual      │  │ on delete   │  │ requests   │ │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘ │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Managed Resources                                │    │
│  │                                                                     │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │    │
│  │  │Deployment│  │  Service │  │ ConfigMap│  │  Secret  │            │    │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │    │
│  │  │   PVC    │  │ CronJob  │  │ Ingress  │  │ PodDisruptionBudget│  │    │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Controller Reconciliation Loop

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Controller Reconciliation Loop                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────┐                                                          │
│  │    Start      │                                                          │
│  └───────┬───────┘                                                          │
│          │                                                                  │
│          ▼                                                                  │
│  ┌───────────────┐     ┌───────────────┐     ┌───────────────┐             │
│  │   Watch CR    │────►│  Event Queue  │────►│  Dequeue      │             │
│  │   Changes     │     │  (Rate Limiter)│     │  Request      │             │
│  └───────────────┘     └───────────────┘     └───────┬───────┘             │
│                                                      │                       │
│                                                      ▼                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Reconcile()                                   │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  1. Fetch current state of CR                                 │  │   │
│  │  │     - Get Database object                                     │  │   │
│  │  │     - Handle not found (deleted)                              │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │                              ▼                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  2. Fetch actual state of managed resources                   │  │   │
│  │  │     - Get Deployment, Service, PVC status                     │  │   │
│  │  │     - Check pod health                                        │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │                              ▼                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  3. Compare desired vs actual state                           │  │   │
│  │  │     - Spec changed?                                           │  │   │
│  │  │     - Resources missing?                                      │  │   │
│  │  │     - Pods unhealthy?                                         │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │                              ▼                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  4. Apply changes                                             │  │   │
│  │  │     - Create missing resources                                │  │   │
│  │  │     - Update changed resources                                │  │   │
│  │  │     - Scale replicas                                          │  │   │
│  │  │     - Trigger backup if scheduled                             │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │                              ▼                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │  5. Update status                                             │  │   │
│  │  │     - Set phase: Pending/Creating/Running/Failed              │  │   │
│  │  │     - Update conditions                                       │  │   │
│  │  │     - Set observedGeneration                                  │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                               │
│              ┌───────────────┴───────────────┐                               │
│              │ Error?                        │ Success?                      │
│              ▼                               ▼                               │
│  ┌───────────────────┐          ┌───────────────────┐                       │
│  │ Requeue with      │          │ Return nil        │                       │
│  │ backoff           │          │ (no requeue)      │                       │
│  └───────────────────┘          └───────────────────┘                       │
│                                                                              │
│  Key Principles:                                                             │
│  - Idempotent: Same input = same result                                      │
│  - Level-triggered: React to current state, not events                       │
│  - Self-healing: Continuously converge to desired state                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Scalability and Performance Considerations

### 3.1 Controller Performance Optimization

| Technique | Before | After | Benefit |
|-----------|--------|-------|---------|
| Rate Limiting | Unbounded reconcile | Max 10/sec | Prevents thundering herd |
| Worker Pool | Single worker | 10 workers | Parallel processing |
| Predicate Filtering | All events processed | Filter irrelevant | Reduced CPU |
| Cache Indexing | O(n) lookups | O(1) lookups | Faster state fetch |
| Leader Election | Multiple active | Single leader | No conflicts |

### 3.2 Resource Management

```go
package controllers

import (
    "context"
    "sync"
    "time"

    "golang.org/x/time/rate"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/predicate"
    "sigs.k8s.io/controller-runtime/pkg/ratelimiter"
)

// PerformanceTunedController configures controller for production
type PerformanceTunedController struct {
    maxConcurrentReconciles int
    rateLimiter             ratelimiter.RateLimiter
    cacheSyncTimeout        time.Duration
}

func NewPerformanceTunedController() *PerformanceTunedController {
    return &PerformanceTunedController{
        maxConcurrentReconciles: 10,
        rateLimiter: workqueue.NewMaxOfRateLimiter(
            workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
            &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
        ),
        cacheSyncTimeout: 2 * time.Minute,
    }
}

// SetupWithManager configures the controller with performance tuning
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&databasev1.Database{}).
        WithEventFilter(predicate.GenerationChangedPredicate{}). // Skip status-only changes
        WithOptions(controller.Options{
            MaxConcurrentReconciles: 10,
            RateLimiter:             r.rateLimiter,
            CacheSyncTimeout:        2 * time.Minute,
        }).
        Owns(&appsv1.Deployment{}).
        Owns(&corev1.Service{}).
        Owns(&corev1.PersistentVolumeClaim{}).
        Watches(
            &source.Kind{Type: &corev1.Pod{}},
            handler.EnqueueRequestsFromMapFunc(r.findDatabasesForPod),
        ).
        Complete(r)
}

// predicate filters to reduce unnecessary reconciliations
type DatabasePredicate struct {
    predicate.Funcs
}

func (p DatabasePredicate) Update(e event.UpdateEvent) bool {
    oldDatabase := e.ObjectOld.(*databasev1.Database)
    newDatabase := e.ObjectNew.(*databasev1.Database)

    // Only reconcile if spec changed (ignore status updates)
    return oldDatabase.Generation != newDatabase.Generation
}

func (p DatabasePredicate) Create(e event.CreateEvent) bool {
    return true
}

func (p DatabasePredicate) Delete(e event.DeleteEvent) bool {
    return true
}
```

---

## 4. Technology Stack Recommendations

### 4.1 Operator Development Tools

| Tool | Purpose | Maturity | Learning Curve |
|------|---------|----------|----------------|
| Kubebuilder | Scaffold and build operators | Production | Medium |
| Operator SDK | Red Hat's operator framework | Production | Medium |
| Controller-runtime | Core controller library | Production | High |
| Kustomize | Kubernetes manifest management | Production | Low |
| Helm | Package management | Production | Low |
| OLM | Operator Lifecycle Manager | Production | Medium |

### 4.2 Recommended Project Structure

```
my-operator/
├── api/
│   └── v1/
│       ├── database_types.go          # CRD Go types
│       ├── groupversion_info.go       # API metadata
│       └── zz_generated.deepcopy.go   # Generated code
├── config/
│   ├── crd/                           # CRD manifests
│   ├── manager/                       # Operator deployment
│   ├── rbac/                          # RBAC configurations
│   └── samples/                       # Example CRs
├── controllers/
│   ├── database_controller.go         # Main controller
│   ├── suite_test.go                  # Integration tests
│   └── helpers.go                     # Controller utilities
├── internal/
│   ├── resources/
│   │   ├── deployment.go              # Deployment builder
│   │   ├── service.go                 # Service builder
│   │   └── pvc.go                     # PVC builder
│   └── utils/
│       ├── constants.go
│       └── labels.go
├── pkg/
│   └── backup/
│       └── backup.go                  # Backup logic
├── hack/
│   └── boilerplate.go.txt
├── Makefile
├── PROJECT                            # Kubebuilder project file
├── go.mod
└── main.go
```

---

## 5. Case Studies

### 5.1 Prometheus Operator

**Scale:** 10,000+ metrics endpoints managed

**Key Features:**

- Custom ServiceMonitor CRD for target discovery
- Prometheus CRD for instance configuration
- Alertmanager CRD for alert routing
- Rule CRD for recording/alerting rules

**Lessons Learned:**

- CRD versioning strategy (v1alpha1 → v1beta1 → v1)
- Validation webhooks prevent misconfiguration
- Status subresource for observed state

### 5.2 Strimzi Kafka Operator

**Scale:** 1000+ Kafka clusters managed

**Architecture Decisions:**

- StatefulSet for broker persistence
- Operator handles rolling upgrades
- User operator for ACL management
- Topic operator for topic lifecycle

**Key Patterns:**

- Multiple controllers in one operator
- Custom resource dependencies
- Event-driven reconciliation

---

## 6. Go Implementation Examples

### 6.1 Complete Database Operator

```go
package controllers

import (
    "context"
    "fmt"
    "time"

    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
    "sigs.k8s.io/controller-runtime/pkg/log"

    databasev1 "github.com/example/my-operator/api/v1"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=database.example.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.example.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.example.com,resources=databases/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete

// Reconcile is the main reconciliation loop
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)
    log.Info("Reconciling Database", "name", req.Name)

    // Fetch the Database CR
    database := &databasev1.Database{}
    if err := r.Get(ctx, req.NamespacedName, database); err != nil {
        if errors.IsNotFound(err) {
            // Object not found, could have been deleted
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, err
    }

    // Handle deletion
    if !database.DeletionTimestamp.IsZero() {
        return r.handleDeletion(ctx, database)
    }

    // Add finalizer if not present
    if !controllerutil.ContainsFinalizer(database, databaseFinalizer) {
        controllerutil.AddFinalizer(database, databaseFinalizer)
        if err := r.Update(ctx, database); err != nil {
            return ctrl.Result{}, err
        }
    }

    // Reconcile PVC for data persistence
    pvc, err := r.reconcilePVC(ctx, database)
    if err != nil {
        r.updateStatus(ctx, database, databasev1.PhaseFailed, err.Error())
        return ctrl.Result{}, err
    }

    // Reconcile ConfigMap for configuration
    configMap, err := r.reconcileConfigMap(ctx, database)
    if err != nil {
        r.updateStatus(ctx, database, databasev1.PhaseFailed, err.Error())
        return ctrl.Result{}, err
    }

    // Reconcile Deployment
    deployment, err := r.reconcileDeployment(ctx, database, pvc, configMap)
    if err != nil {
        r.updateStatus(ctx, database, databasev1.PhaseFailed, err.Error())
        return ctrl.Result{}, err
    }

    // Reconcile Service
    service, err := r.reconcileService(ctx, database)
    if err != nil {
        r.updateStatus(ctx, database, databasev1.PhaseFailed, err.Error())
        return ctrl.Result{}, err
    }

    // Check deployment status
    if deployment.Status.ReadyReplicas < database.Spec.Replicas {
        r.updateStatus(ctx, database, databasev1.PhaseCreating,
            fmt.Sprintf("Waiting for replicas: %d/%d", deployment.Status.ReadyReplicas, database.Spec.Replicas))
        return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
    }

    // Reconcile backup CronJob if enabled
    if database.Spec.Backup.Enabled {
        if err := r.reconcileBackupCronJob(ctx, database); err != nil {
            log.Error(err, "Failed to reconcile backup")
        }
    }

    // Update status to running
    r.updateStatus(ctx, database, databasev1.PhaseRunning, "Database is healthy")

    // Set service endpoint in status
    database.Status.Endpoint = fmt.Sprintf("%s.%s.svc.cluster.local:%d",
        service.Name, service.Namespace, database.Spec.Port)
    database.Status.StorageUsed = pvc.Status.Capacity.Storage()

    if err := r.Status().Update(ctx, database); err != nil {
        return ctrl.Result{}, err
    }

    // Schedule next reconciliation for health checks
    return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *DatabaseReconciler) reconcilePVC(ctx context.Context, database *databasev1.Database) (*corev1.PersistentVolumeClaim, error) {
    pvc := &corev1.PersistentVolumeClaim{
        ObjectMeta: metav1.ObjectMeta{
            Name:      database.Name + "-data",
            Namespace: database.Namespace,
            Labels:    r.labelsForDatabase(database),
        },
        Spec: corev1.PersistentVolumeClaimSpec{
            AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
            Resources: corev1.VolumeResourceRequirements{
                Requests: corev1.ResourceList{
                    corev1.ResourceStorage: database.Spec.Storage,
                },
            },
            StorageClassName: database.Spec.StorageClassName,
        },
    }

    // Set owner reference for garbage collection
    if err := controllerutil.SetControllerReference(database, pvc, r.Scheme); err != nil {
        return nil, err
    }

    // Create or update
    found := &corev1.PersistentVolumeClaim{}
    err := r.Get(ctx, client.ObjectKeyFromObject(pvc), found)
    if err != nil && errors.IsNotFound(err) {
        return pvc, r.Create(ctx, pvc)
    } else if err != nil {
        return nil, err
    }

    // PVC is immutable after creation, only update labels/annotations
    found.Labels = pvc.Labels
    return found, r.Update(ctx, found)
}

func (r *DatabaseReconciler) reconcileDeployment(ctx context.Context, database *databasev1.Database,
    pvc *corev1.PersistentVolumeClaim, configMap *corev1.ConfigMap) (*appsv1.Deployment, error) {

    replicas := database.Spec.Replicas
    if replicas == 0 {
        replicas = 1
    }

    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      database.Name,
            Namespace: database.Namespace,
            Labels:    r.labelsForDatabase(database),
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &replicas,
            Selector: &metav1.LabelSelector{
                MatchLabels: r.labelsForDatabase(database),
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: r.labelsForDatabase(database),
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "database",
                            Image: fmt.Sprintf("postgres:%s", database.Spec.Version),
                            Ports: []corev1.ContainerPort{
                                {ContainerPort: database.Spec.Port, Name: "postgres"},
                            },
                            Env: []corev1.EnvVar{
                                {
                                    Name: "POSTGRES_PASSWORD",
                                    ValueFrom: &corev1.EnvVarSource{
                                        SecretKeyRef: &corev1.SecretKeySelector{
                                            LocalObjectReference: corev1.LocalObjectReference{
                                                Name: database.Name + "-credentials",
                                            },
                                            Key: "password",
                                        },
                                    },
                                },
                            },
                            VolumeMounts: []corev1.VolumeMount{
                                {
                                    Name:      "data",
                                    MountPath: "/var/lib/postgresql/data",
                                },
                                {
                                    Name:      "config",
                                    MountPath: "/etc/postgresql/postgresql.conf",
                                    SubPath:   "postgresql.conf",
                                },
                            },
                            Resources: database.Spec.Resources,
                            LivenessProbe: &corev1.Probe{
                                Exec: &corev1.ExecAction{
                                    Command: []string{"pg_isready", "-U", "postgres"},
                                },
                                InitialDelaySeconds: 30,
                                PeriodSeconds:       10,
                            },
                            ReadinessProbe: &corev1.Probe{
                                Exec: &corev1.ExecAction{
                                    Command: []string{"pg_isready", "-U", "postgres"},
                                },
                                InitialDelaySeconds: 5,
                                PeriodSeconds:       5,
                            },
                        },
                    },
                    Volumes: []corev1.Volume{
                        {
                            Name: "data",
                            VolumeSource: corev1.VolumeSource{
                                PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                                    ClaimName: pvc.Name,
                                },
                            },
                        },
                        {
                            Name: "config",
                            VolumeSource: corev1.VolumeSource{
                                ConfigMap: &corev1.ConfigMapVolumeSource{
                                    LocalObjectReference: corev1.LocalObjectReference{
                                        Name: configMap.Name,
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Set owner reference
    if err := controllerutil.SetControllerReference(database, deployment, r.Scheme); err != nil {
        return nil, err
    }

    // Create or update
    found := &appsv1.Deployment{}
    err := r.Get(ctx, client.ObjectKeyFromObject(deployment), found)
    if err != nil && errors.IsNotFound(err) {
        return deployment, r.Create(ctx, deployment)
    } else if err != nil {
        return nil, err
    }

    // Update deployment
    found.Spec = deployment.Spec
    found.Labels = deployment.Labels
    return found, r.Update(ctx, found)
}

func (r *DatabaseReconciler) reconcileService(ctx context.Context, database *databasev1.Database) (*corev1.Service, error) {
    service := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      database.Name,
            Namespace: database.Namespace,
            Labels:    r.labelsForDatabase(database),
        },
        Spec: corev1.ServiceSpec{
            Selector: r.labelsForDatabase(database),
            Ports: []corev1.ServicePort{
                {
                    Port:       database.Spec.Port,
                    TargetPort: intstr.FromInt(int(database.Spec.Port)),
                    Name:       "postgres",
                },
            },
            Type: corev1.ServiceTypeClusterIP,
        },
    }

    if err := controllerutil.SetControllerReference(database, service, r.Scheme); err != nil {
        return nil, err
    }

    // Create or update
    found := &corev1.Service{}
    err := r.Get(ctx, client.ObjectKeyFromObject(service), found)
    if err != nil && errors.IsNotFound(err) {
        return service, r.Create(ctx, service)
    } else if err != nil {
        return nil, err
    }

    found.Spec = service.Spec
    found.Labels = service.Labels
    return found, r.Update(ctx, found)
}

func (r *DatabaseReconciler) updateStatus(ctx context.Context, database *databasev1.Database,
    phase databasev1.DatabasePhase, message string) {

    database.Status.Phase = phase
    database.Status.Message = message
    database.Status.LastUpdated = metav1.Now()

    // Update conditions
    condition := databasev1.DatabaseCondition{
        Type:               databasev1.ConditionReady,
        Status:             corev1.ConditionTrue,
        LastTransitionTime: metav1.Now(),
        Reason:             string(phase),
        Message:            message,
    }

    if phase == databasev1.PhaseFailed {
        condition.Status = corev1.ConditionFalse
    }

    database.Status.Conditions = updateConditions(database.Status.Conditions, condition)

    if err := r.Status().Update(ctx, database); err != nil {
        log.FromContext(ctx).Error(err, "Failed to update status")
    }
}

func (r *DatabaseReconciler) labelsForDatabase(database *databasev1.Database) map[string]string {
    return map[string]string{
        "app.kubernetes.io/name":       "database",
        "app.kubernetes.io/instance":   database.Name,
        "app.kubernetes.io/managed-by": "database-operator",
    }
}
```

### 6.2 Webhook Validation

```go
package v1

import (
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    logf "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/controller-runtime/pkg/webhook"
    "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var databaselog = logf.Log.WithName("database-resource")

func (r *Database) SetupWebhookWithManager(mgr ctrl.Manager) error {
    return ctrl.NewWebhookManagedBy(mgr).
        For(r).
        Complete()
}

//+kubebuilder:webhook:path=/mutate-database-example-com-v1-database,mutating=true,failurePolicy=fail,sideEffects=None,groups=database.example.com,resources=databases,verbs=create;update,versions=v1,name=mdatabase.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Database{}

// Default implements webhook.Defaulter
func (r *Database) Default() {
    databaselog.Info("default", "name", r.Name)

    // Set default values
    if r.Spec.Version == "" {
        r.Spec.Version = "15"
    }
    if r.Spec.Port == 0 {
        r.Spec.Port = 5432
    }
    if r.Spec.Replicas == 0 {
        r.Spec.Replicas = 1
    }
    if r.Spec.Storage.IsZero() {
        r.Spec.Storage = resource.MustParse("10Gi")
    }
}

//+kubebuilder:webhook:path=/validate-database-example-com-v1-database,mutating=false,failurePolicy=fail,sideEffects=None,groups=database.example.com,resources=databases,verbs=create;update,versions=v1,name=vdatabase.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Database{}

// ValidateCreate implements webhook.Validator
func (r *Database) ValidateCreate() (admission.Warnings, error) {
    databaselog.Info("validate create", "name", r.Name)

    if err := r.validateDatabase(); err != nil {
        return nil, err
    }
    return nil, nil
}

// ValidateUpdate implements webhook.Validator
func (r *Database) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
    databaselog.Info("validate update", "name", r.Name)

    if err := r.validateDatabase(); err != nil {
        return nil, err
    }

    // Check immutable fields
    oldDatabase := old.(*Database)
    if r.Spec.Storage.Cmp(oldDatabase.Spec.Storage) < 0 {
        return nil, fmt.Errorf("storage size cannot be reduced")
    }

    return nil, nil
}

// ValidateDelete implements webhook.Validator
func (r *Database) ValidateDelete() (admission.Warnings, error) {
    databaselog.Info("validate delete", "name", r.Name)

    // Could check for active connections, backups in progress, etc.
    return nil, nil
}

func (r *Database) validateDatabase() error {
    // Validate version
    validVersions := []string{"13", "14", "15", "16"}
    if !contains(validVersions, r.Spec.Version) {
        return fmt.Errorf("invalid version %s, must be one of: %v", r.Spec.Version, validVersions)
    }

    // Validate replicas
    if r.Spec.Replicas > 5 {
        return fmt.Errorf("replicas cannot exceed 5")
    }

    // Validate storage
    minStorage := resource.MustParse("1Gi")
    maxStorage := resource.MustParse("1Ti")
    if r.Spec.Storage.Cmp(minStorage) < 0 {
        return fmt.Errorf("storage must be at least %s", minStorage.String())
    }
    if r.Spec.Storage.Cmp(maxStorage) > 0 {
        return fmt.Errorf("storage cannot exceed %s", maxStorage.String())
    }

    // Validate backup schedule
    if r.Spec.Backup.Enabled && r.Spec.Backup.Schedule != "" {
        if _, err := cron.ParseStandard(r.Spec.Backup.Schedule); err != nil {
            return fmt.Errorf("invalid backup schedule: %v", err)
        }
    }

    return nil
}
```

---

## 7. Visual Representations

### 7.1 Operator Capability Model

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Operator Capability Model                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Level 5: Auto-Pilot                                                         │
│  ══════════════════                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  - Automatic vertical/horizontal scaling based on metrics          │    │
│  │  - Self-healing without human intervention                         │    │
│  │  - Predictive maintenance and proactive actions                    │    │
│  │  - Multi-cluster orchestration                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    ▲                                         │
│  Level 4: Deep Insights                                                      │
│  ═══════════════════════                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  - Advanced monitoring and alerting                                │    │
│  │  - Performance tuning recommendations                              │    │
│  │  - Capacity planning insights                                      │    │
│  │  - Cost optimization suggestions                                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    ▲                                         │
│  Level 3: Full Lifecycle                                                     │
│  ═════════════════════════                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  - Automated backups and point-in-time recovery                    │    │
│  │  - Zero-downtime upgrades                                          │    │
│  │  - Configuration management and drift detection                    │    │
│  │  - Disaster recovery automation                                    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    ▲                                         │
│  Level 2: Seamless Upgrades                                                  │
│  ══════════════════════════                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  - Rolling updates with health checks                              │    │
│  │  - Rollback capability                                             │    │
│  │  - Canary deployments                                              │    │
│  │  - Version management                                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    ▲                                         │
│  Level 1: Basic Install                                                      │
│  ═════════════════════                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  - Automated provisioning of application and resources             │    │
│  │  - Basic configuration management                                  │    │
│  │  - Health checks and status reporting                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Resource Ownership Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Resource Ownership Hierarchy                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Custom Resource (Owner)                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Database (mydb)                                                    │    │
│  │  OwnerReference: nil (top-level)                                    │    │
│  │  Finalizers: [database-finalizer]                                   │    │
│  └──────────────────┬──────────────────────────────────────────────────┘    │
│                     │                                                        │
│         ┌───────────┼───────────┬───────────────────┐                        │
│         │           │           │                   │                        │
│         ▼           ▼           ▼                   ▼                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐     ┌──────────┐                   │
│  │Deployment│ │ Service  │ │   PVC    │     │ ConfigMap│                   │
│  │          │ │          │ │          │     │          │                   │
│  │ownerRefs:│ │ownerRefs:│ │ownerRefs:│     │ownerRefs:│                   │
│  │- Database│ │- Database│ │- Database│     │- Database│                   │
│  │  mydb    │ │  mydb    │ │  mydb    │     │  mydb    │                   │
│  └────┬─────┘ └──────────┘ └──────────┘     └──────────┘                   │
│       │                                                                      │
│       ▼                                                                      │
│  ┌──────────┐                                                                │
│  │  ReplicaSet│                                                               │
│  │(created by │                                                               │
│  │Deployment) │                                                               │
│  └────┬─────┘                                                                │
│       │                                                                      │
│       ▼                                                                      │
│  ┌──────────┐                                                                │
│  │   Pods   │                                                                │
│  │          │                                                                │
│  └──────────┘                                                                │
│                                                                              │
│  Garbage Collection:                                                         │
│  - Deleting Database CR → cascades delete to all owned resources            │
│  - OwnerReference blocks deletion until dependent resources are handled     │
│  - Finalizer delays CR deletion until cleanup is complete                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Multi-Operator Coordination

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Multi-Operator Coordination                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Cluster Operators                                │    │
│  │                                                                     │    │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐           │    │
│  │  │   Cert-       │  │   Ingress     │  │   Monitoring  │           │    │
│  │  │   Manager     │  │   Controller  │  │   Stack       │           │    │
│  │  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘           │    │
│  │          │                  │                  │                    │    │
│  │          ▼                  ▼                  ▼                    │    │
│  │  ┌─────────────────────────────────────────────────────────────┐   │    │
│  │  │              Shared Infrastructure Resources                │   │    │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │   │    │
│  │  │  │ Certificate │  │   Ingress   │  │   ServiceMonitor    │  │   │    │
│  │  │  │   Secrets   │  │   Rules     │  │                     │  │   │    │
│  │  │  └─────────────┘  └─────────────┘  └─────────────────────┘  │   │    │
│  │  └─────────────────────────────────────────────────────────────┘   │    │
│  │                              │                                      │    │
│  └──────────────────────────────┼──────────────────────────────────────┘    │
│                                 │                                            │
│                                 ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                   Application Operators                             │    │
│  │                                                                     │    │
│  │  ┌───────────────────────────────────────────────────────────────┐  │    │
│  │  │  Database Operator                                            │  │    │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐   │  │    │
│  │  │  │ Database CR │  │ Database    │  │ Uses: Certificate   │   │  │    │
│  │  │  │             │  │ Resources   │  │       ServiceMonitor│   │  │    │
│  │  │  └─────────────┘  └─────────────┘  └─────────────────────┘   │  │    │
│  │  └───────────────────────────────────────────────────────────────┘  │    │
│  │                                                                     │    │
│  │  ┌───────────────────────────────────────────────────────────────┐  │    │
│  │  │  Cache Operator                                               │  │    │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐   │  │    │
│  │  │  │   Cache CR  │  │ Cache       │  │ Uses: Certificate   │   │  │    │
│  │  │  │             │  │ Resources   │  │       ServiceMonitor│   │  │    │
│  │  │  └─────────────┘  └─────────────┘  └─────────────────────┘   │  │    │
│  │  └───────────────────────────────────────────────────────────────┘  │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  Coordination Patterns:                                                      │
│  - Resource watches (cross-namespace)                                        │
│  - Status conditions indicating dependencies                                 │
│  - Events for communication between operators                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Security Requirements

### 8.1 Operator Security Checklist

| Category | Requirement | Implementation |
|----------|-------------|----------------|
| RBAC | Least privilege | Specific verbs on specific resources |
| Secrets | Encryption at rest | Kubernetes secret encryption |
| Images | Signed images | Cosign/Notary verification |
| Network | Network policies | Restrict pod-to-pod communication |
| Audit | Audit logging | API server audit configuration |
| Webhooks | TLS encryption | Cert-manager for certificates |

### 8.2 RBAC Configuration

```yaml
# config/rbac/role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: database-operator-role
rules:
# Database CRD management
- apiGroups: ["database.example.com"]
  resources: ["databases"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: ["database.example.com"]
  resources: ["databases/status"]
  verbs: ["get", "update", "patch"]
- apiGroups: ["database.example.com"]
  resources: ["databases/finalizers"]
  verbs: ["update"]

# Core resources
- apiGroups: [""]
  resources: ["configmaps", "secrets", "services", "persistentvolumeclaims", "events"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]

# Apps resources
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Batch resources (for backups)
- apiGroups: ["batch"]
  resources: ["cronjobs", "jobs"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

---

**Quality Rating**: S (18+ KB)
**Last Updated**: 2026-04-02
