# Kubernetes Operators

> **分类**: 工程与云原生
> **标签**: #kubernetes #operators #crd #controller #automation
> **参考**: Kubernetes Operator Pattern, Operator SDK, CoreOS Operators

---

## 1. Formal Definition

### 1.1 What is a Kubernetes Operator?

A Kubernetes Operator is a method of packaging, deploying, and managing a Kubernetes application by extending the Kubernetes API through Custom Resource Definitions (CRDs) and custom controllers. Operators encode operational knowledge - the expertise required to run complex software - into software that automates lifecycle management.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Kubernetes Operator Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        KUBERNETES CLUSTER                            │   │
│   │                                                                       │   │
│   │  ┌───────────────────────────────────────────────────────────────┐   │   │
│   │  │                    CUSTOM RESOURCE DEFINITION                    │   │   │
│   │  │                                                                   │   │   │
│   │  │  apiVersion: databases.example.com/v1                             │   │   │
│   │  │  kind: PostgreSQL                                                 │   │   │
│   │  │  metadata:                                                        │   │   │
│   │  │    name: production-db                                            │   │   │
│   │  │  spec:                                                            │   │   │
│   │  │    version: "13"                                                  │   │   │
│   │  │    replicas: 3                                                    │   │   │
│   │  │    storage: 100Gi                                                 │   │   │
│   │  │    backup:                                                        │   │   │
│   │  │      enabled: true                                                │   │   │
│   │  │      schedule: "0 2 * * *"                                        │   │   │
│   │  └───────────────────────────────────────────────────────────────┘   │   │
│   │                              │                                        │   │
│   │                              │ WATCH                                  │   │
│   │                              ▼                                        │   │
│   │  ┌───────────────────────────────────────────────────────────────┐   │   │
│   │  │                    OPERATOR CONTROLLER                         │   │   │
│   │  │                                                                   │   │   │
│   │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │   │   │
│   │  │  │   Reconcile  │  │   Validate   │  │   Execute    │          │   │   │
│   │  │  │     Loop     │──│     CR       │──│   Actions    │          │   │   │
│   │  │  │              │  │              │  │              │          │   │   │
│   │  │  │ • Watch CR   │  │ • Schema     │  │ • Create     │          │   │   │
│   │  │  │ • Compare    │  │   validation │  │   resources  │          │   │   │
│   │  │  │   state      │  │ • Business   │  │ • Update     │          │   │   │
│   │  │  │ • Apply      │  │   rules      │  │   state      │          │   │   │
│   │  │  │   changes    │  │              │  │ • Cleanup    │          │   │   │
│   │  │  └──────────────┘  └──────────────┘  └──────────────┘          │   │   │
│   │  │                                                                   │   │   │
│   │  │  Operational Knowledge Encoded:                                   │   │   │
│   │  │  • How to deploy PostgreSQL cluster                               │   │   │
│   │  │  • How to handle failover                                         │   │   │
│   │  │  • How to perform backups                                         │   │   │
│   │  │  • How to upgrade versions                                        │   │   │
│   │  │  • How to handle errors                                           │   │   │
│   │  └───────────────────────────────────────────────────────────────┘   │   │
│   │                              │                                        │   │
│   │                              │ MANAGE                                 │   │
│   │                              ▼                                        │   │
│   │  ┌───────────────────────────────────────────────────────────────┐   │   │
│   │  │                    MANAGED RESOURCES                             │   │   │
│   │  │                                                                   │   │   │
│   │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐            │   │   │
│   │  │  │   Pod   │  │ Service │  │   PVC   │  │  Secret │            │   │   │
│   │  │  │  (DB)   │  │ (Access)│  │(Storage)│  │(Creds)  │            │   │   │
│   │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘            │   │   │
│   │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                        │   │   │
│   │  │  │ConfigMap│  │CronJob  │  │ Network │                        │   │   │
│   │  │  │(Config) │  │(Backup) │  │ Policy  │                        │   │   │
│   │  │  └─────────┘  └─────────┘  └─────────┘                        │   │   │
│   │  └───────────────────────────────────────────────────────────────┘   │   │
│   │                                                                       │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│   OPERATOR CAPABILITIES:                                                    │
│   • Automated provisioning and configuration                                │
│   • Self-healing and failover                                               │
│   • Backup and restore automation                                           │
│   • Rolling upgrades with zero downtime                                     │
│   • Metric collection and alerting integration                              │
│   • Custom business logic enforcement                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Operator Capability Levels

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Operator Capability Levels                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   LEVEL 5 - AUTO PILOT                    LEVEL 4 - DEEP INSIGHTS           │
│   ━━━━━━━━━━━━━━━━━━━━━━━                 ━━━━━━━━━━━━━━━━━━━━━━━━━         │
│                                                                             │
│   ✓ Horizontal/vertical auto-scaling      ✓ Automatic metric collection     │
│   ✓ Automatic configuration tuning        ✓ Custom metrics and alerting     │
│   ✓ Self-healing without intervention     ✓ Workload analysis and tuning    │
│   ✓ Predictive maintenance                ✓ Distributed tracing             │
│   ✓ Cost optimization                     ✓ Log aggregation                 │
│                                                                             │
│   LEVEL 3 - FULL LIFECYCLE                LEVEL 2 - SEAMLESS UPGRADES       │
│   ━━━━━━━━━━━━━━━━━━━━━━━━━━              ━━━━━━━━━━━━━━━━━━━━━━━━━         │
│                                                                             │
│   ✓ Backup and restore                    ✓ Minor version updates           │
│   ✓ Failure recovery                      ✓ Major version updates           │
│   ✓ Storage management                    ✓ Rolling updates                 │
│   ✓ Connection pooling                    ✓ Rollback support                │
│   ✓ Configuration management              ✓ Schema migration                │
│                                                                             │
│   LEVEL 1 - BASIC INSTALL                 LEVEL 0 - PLANNING                │
│   ━━━━━━━━━━━━━━━━━━━━━━━                 ━━━━━━━━━━━━━━━━━━━               │
│                                                                             │
│   ✓ Automated install                     ✓ Custom Resource Definition      │
│   ✓ Health checks                         ✓ Basic controller skeleton       │
│   ✓ Basic configuration                   ✓ Resource requirements analysis  │
│                                                                             │
│   ────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│   PROGRESSION PATH:                                                         │
│   Level 0 → Level 1 → Level 2 → Level 3 → Level 4 → Level 5                 │
│   (Planning)  (Basic)   (Updates)(Lifecycle)(Insights)(Auto-Pilot)         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 Basic Controller Pattern

```go
package controller

import (
    "context"
    "fmt"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/util/workqueue"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/handler"
    "sigs.k8s.io/controller-runtime/pkg/manager"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
    "sigs.k8s.io/controller-runtime/pkg/source"
)

// MyAppSpec defines the desired state of MyApp
type MyAppSpec struct {
    Replicas int32  `json:"replicas,omitempty"`
    Image    string `json:"image,omitempty"`
    Port     int32  `json:"port,omitempty"`
}

// MyAppStatus defines the observed state of MyApp
type MyAppStatus struct {
    Phase      string `json:"phase,omitempty"`
    Replicas   int32  `json:"replicas,omitempty"`
    ReadyReplicas int32 `json:"readyReplicas,omitempty"`
    Conditions []Condition `json:"conditions,omitempty"`
}

// Condition represents a condition of the resource
type Condition struct {
    Type               string    `json:"type"`
    Status             string    `json:"status"`
    LastTransitionTime time.Time `json:"lastTransitionTime"`
    Reason             string    `json:"reason,omitempty"`
    Message            string    `json:"message,omitempty"`
}

// MyApp is the Schema for the myapps API
type MyApp struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   MyAppSpec   `json:"spec,omitempty"`
    Status MyAppStatus `json:"status,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function
func (in *MyApp) DeepCopyInto(out *MyApp) {
    *out = *in
    out.TypeMeta = in.TypeMeta
    in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
    out.Spec = in.Spec
    in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function
func (in *MyApp) DeepCopy() *MyApp {
    if in == nil {
        return nil
    }
    out := new(MyApp)
    in.DeepCopyInto(out)
    return out
}

// DeepCopyObject is an autogenerated deepcopy function
func (in *MyApp) DeepCopyObject() runtime.Object {
    if c := in.DeepCopy(); c != nil {
        return c
    }
    return nil
}

// DeepCopyInto for MyAppStatus
func (in *MyAppStatus) DeepCopyInto(out *MyAppStatus) {
    *out = *in
    if in.Conditions != nil {
        in, out := &in.Conditions, &out.Conditions
        *out = make([]Condition, len(*in))
        for i := range *in {
            (*in)[i].LastTransitionTime = (*out)[i].LastTransitionTime
        }
    }
}

// Reconciler reconciles a MyApp object
type Reconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.example.com,resources=myapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=myapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile implements the reconciliation loop
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
    log := log.FromContext(ctx)

    // Fetch the MyApp instance
    myapp := &MyApp{}
    if err := r.Get(ctx, req.NamespacedName, myapp); err != nil {
        if errors.IsNotFound(err) {
            // Object not found, return
            return reconcile.Result{}, nil
        }
        return reconcile.Result{}, err
    }

    log.Info("Reconciling MyApp", "name", myapp.Name)

    // Create or update Deployment
    if err := r.reconcileDeployment(ctx, myapp); err != nil {
        r.updateStatus(ctx, myapp, "Failed", fmt.Sprintf("Deployment error: %v", err))
        return reconcile.Result{}, err
    }

    // Create or update Service
    if err := r.reconcileService(ctx, myapp); err != nil {
        r.updateStatus(ctx, myapp, "Failed", fmt.Sprintf("Service error: %v", err))
        return reconcile.Result{}, err
    }

    // Update status
    if err := r.updateStatus(ctx, myapp, "Running", "Successfully reconciled"); err != nil {
        return reconcile.Result{}, err
    }

    // Requeue for periodic status checks
    return reconcile.Result{RequeueAfter: 60 * time.Second}, nil
}

// reconcileDeployment ensures the Deployment exists and matches spec
func (r *Reconciler) reconcileDeployment(ctx context.Context, myapp *MyApp) error {
    deployment := &appsv1.Deployment{}
    deploymentName := myapp.Name

    err := r.Get(ctx, types.NamespacedName{
        Name:      deploymentName,
        Namespace: myapp.Namespace,
    }, deployment)

    if err != nil && errors.IsNotFound(err) {
        // Create new deployment
        deployment = r.buildDeployment(myapp)
        if err := r.Create(ctx, deployment); err != nil {
            return fmt.Errorf("failed to create deployment: %w", err)
        }
        return nil
    } else if err != nil {
        return fmt.Errorf("failed to get deployment: %w", err)
    }

    // Update existing deployment if needed
    desiredDeployment := r.buildDeployment(myapp)
    if !r.deploymentsEqual(deployment, desiredDeployment) {
        deployment.Spec = desiredDeployment.Spec
        if err := r.Update(ctx, deployment); err != nil {
            return fmt.Errorf("failed to update deployment: %w", err)
        }
    }

    return nil
}

// reconcileService ensures the Service exists and matches spec
func (r *Reconciler) reconcileService(ctx context.Context, myapp *MyApp) error {
    service := &corev1.Service{}
    serviceName := myapp.Name

    err := r.Get(ctx, types.NamespacedName{
        Name:      serviceName,
        Namespace: myapp.Namespace,
    }, service)

    if err != nil && errors.IsNotFound(err) {
        // Create new service
        service = r.buildService(myapp)
        if err := r.Create(ctx, service); err != nil {
            return fmt.Errorf("failed to create service: %w", err)
        }
        return nil
    } else if err != nil {
        return fmt.Errorf("failed to get service: %w", err)
    }

    // Update service if needed
    desiredService := r.buildService(myapp)
    if !r.servicesEqual(service, desiredService) {
        service.Spec = desiredService.Spec
        if err := r.Update(ctx, service); err != nil {
            return fmt.Errorf("failed to update service: %w", err)
        }
    }

    return nil
}

// buildDeployment creates a Deployment from MyApp spec
func (r *Reconciler) buildDeployment(myapp *MyApp) *appsv1.Deployment {
    labels := map[string]string{
        "app":     myapp.Name,
        "managed-by": "myapp-operator",
    }

    replicas := myapp.Spec.Replicas
    if replicas == 0 {
        replicas = 1
    }

    image := myapp.Spec.Image
    if image == "" {
        image = "nginx:latest"
    }

    port := myapp.Spec.Port
    if port == 0 {
        port = 8080
    }

    return &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      myapp.Name,
            Namespace: myapp.Namespace,
            OwnerReferences: []metav1.OwnerReference{
                *metav1.NewControllerRef(myapp, schema.GroupVersionKind{
                    Group:   "apps.example.com",
                    Version: "v1",
                    Kind:    "MyApp",
                }),
            },
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &replicas,
            Selector: &metav1.LabelSelector{
                MatchLabels: labels,
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: labels,
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "app",
                            Image: image,
                            Ports: []corev1.ContainerPort{
                                {
                                    ContainerPort: port,
                                    Protocol:      corev1.ProtocolTCP,
                                },
                            },
                            Resources: corev1.ResourceRequirements{
                                Requests: corev1.ResourceList{
                                    corev1.ResourceCPU:    resource.MustParse("100m"),
                                    corev1.ResourceMemory: resource.MustParse("128Mi"),
                                },
                                Limits: corev1.ResourceList{
                                    corev1.ResourceCPU:    resource.MustParse("500m"),
                                    corev1.ResourceMemory: resource.MustParse("256Mi"),
                                },
                            },
                        },
                    },
                },
            },
        },
    }
}

// buildService creates a Service from MyApp spec
func (r *Reconciler) buildService(myapp *MyApp) *corev1.Service {
    labels := map[string]string{
        "app":     myapp.Name,
        "managed-by": "myapp-operator",
    }

    port := myapp.Spec.Port
    if port == 0 {
        port = 8080
    }

    return &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      myapp.Name,
            Namespace: myapp.Namespace,
            OwnerReferences: []metav1.OwnerReference{
                *metav1.NewControllerRef(myapp, schema.GroupVersionKind{
                    Group:   "apps.example.com",
                    Version: "v1",
                    Kind:    "MyApp",
                }),
            },
        },
        Spec: corev1.ServiceSpec{
            Selector: labels,
            Ports: []corev1.ServicePort{
                {
                    Port:       port,
                    TargetPort: intstr.FromInt(int(port)),
                    Protocol:   corev1.ProtocolTCP,
                },
            },
            Type: corev1.ServiceTypeClusterIP,
        },
    }
}

// updateStatus updates the MyApp status
func (r *Reconciler) updateStatus(ctx context.Context, myapp *MyApp, phase, message string) error {
    myapp.Status.Phase = phase
    myapp.Status.Conditions = append(myapp.Status.Conditions, Condition{
        Type:               "Ready",
        Status:             phase,
        LastTransitionTime: time.Now(),
        Message:            message,
    })

    return r.Status().Update(ctx, myapp)
}

// SetupWithManager sets up the controller with the Manager
func (r *Reconciler) SetupWithManager(mgr manager.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&MyApp{}).
        Owns(&appsv1.Deployment{}).
        Owns(&corev1.Service{}).
        Complete(r)
}
```

### 2.2 Advanced Operator Patterns

```go
package controller

import (
    "context"
    "fmt"
    "time"

    "k8s.io/client-go/tools/record"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Finalizer for cleanup
type FinalizableReconciler struct {
    client.Client
    Scheme   *runtime.Scheme
    Recorder record.EventRecorder
}

// Finalizer name
const myAppFinalizer = "myapp.apps.example.com/finalizer"

// Reconcile with finalizer handling
func (r *FinalizableReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    myapp := &MyApp{}
    if err := r.Get(ctx, req.NamespacedName, myapp); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Check if being deleted
    if !myapp.ObjectMeta.DeletionTimestamp.IsZero() {
        return r.handleDeletion(ctx, myapp)
    }

    // Add finalizer if not present
    if !controllerutil.ContainsFinalizer(myapp, myAppFinalizer) {
        controllerutil.AddFinalizer(myapp, myAppFinalizer)
        if err := r.Update(ctx, myapp); err != nil {
            return ctrl.Result{}, err
        }
    }

    // Normal reconciliation
    return r.reconcileNormal(ctx, myapp)
}

// handleDeletion handles resource deletion
func (r *FinalizableReconciler) handleDeletion(ctx context.Context, myapp *MyApp) (ctrl.Result, error) {
    log := log.FromContext(ctx)
    log.Info("Handling deletion", "name", myapp.Name)

    // Perform cleanup
    if err := r.cleanup(ctx, myapp); err != nil {
        r.Recorder.Event(myapp, corev1.EventTypeWarning, "CleanupFailed", err.Error())
        return ctrl.Result{RequeueAfter: 10 * time.Second}, err
    }

    // Remove finalizer
    controllerutil.RemoveFinalizer(myapp, myAppFinalizer)
    if err := r.Update(ctx, myapp); err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}

// cleanup performs cleanup operations
func (r *FinalizableReconciler) cleanup(ctx context.Context, myapp *MyApp) error {
    // Example: Backup data before deletion
    // Example: Remove external resources
    // Example: Clean up persistent volumes
    return nil
}

// reconcileNormal handles normal reconciliation
func (r *FinalizableReconciler) reconcileNormal(ctx context.Context, myapp *MyApp) (ctrl.Result, error) {
    // Reconciliation logic
    return ctrl.Result{}, nil
}

// WebhookValidator validates MyApp resources
type WebhookValidator struct{}

// ValidateCreate validates on creation
func (v *WebhookValidator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
    myapp, ok := obj.(*MyApp)
    if !ok {
        return fmt.Errorf("expected MyApp but got %T", obj)
    }

    // Validate replicas
    if myapp.Spec.Replicas < 0 || myapp.Spec.Replicas > 100 {
        return fmt.Errorf("replicas must be between 0 and 100")
    }

    // Validate image
    if myapp.Spec.Image == "" {
        return fmt.Errorf("image is required")
    }

    // Validate port
    if myapp.Spec.Port < 1 || myapp.Spec.Port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }

    return nil
}

// ValidateUpdate validates on update
func (v *WebhookValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
    oldApp, newApp := oldObj.(*MyApp), newObj.(*MyApp)

    // Prevent certain immutable fields from changing
    // Example: Storage class changes

    return v.ValidateCreate(ctx, newApp)
}

// ValidateDelete validates on deletion
func (v *WebhookValidator) ValidateDelete(ctx context.Context, obj runtime.Object) error {
    return nil
}
```

---

## 3. Production-Ready Configurations

### 3.1 Operator Deployment

```yaml
# operator-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-operator
  namespace: operators
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp-operator
  template:
    metadata:
      labels:
        app: myapp-operator
    spec:
      serviceAccountName: myapp-operator
      containers:
      - name: manager
        image: myapp-operator:v1.0.0
        command:
        - /manager
        args:
        - --leader-elect
        - --metrics-bind-address=:8080
        - --health-probe-bind-address=:8081
        resources:
          limits:
            cpu: 500m
            memory: 256Mi
          requests:
            cpu: 100m
            memory: 128Mi
        ports:
        - containerPort: 8080
          name: metrics
        - containerPort: 8081
          name: health
        - containerPort: 9443
          name: webhook
          protocol: TCP
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
        volumeMounts:
        - mountPath: /tmp
          name: tmp
      volumes:
      - emptyDir: {}
        name: tmp
      securityContext:
        runAsNonRoot: true
        seccompProfile:
          type: RuntimeDefault
      terminationGracePeriodSeconds: 10

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: myapp-operator
  namespace: operators

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: myapp-operator-role
rules:
- apiGroups:
  - apps.example.com
  resources:
  - myapps
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - apps.example.com
  resources:
  - myapps/finalizers
  verbs:
  - update
- apiGroups:
  - apps.example.com
  resources:
  - myapps/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - services
  - configmaps
  - secrets
  - events
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: myapp-operator-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: myapp-operator-role
subjects:
- kind: ServiceAccount
  name: myapp-operator
  namespace: operators
```

---

## 4. Security Considerations

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Operator Security Best Practices                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  RBAC                                                                       │
│  • Principle of least privilege                                             │
│  • Separate service accounts per operator                                   │
│  • Regular review of permissions                                            │
│  • Avoid cluster-admin for operators                                        │
│                                                                             │
│  SECRETS                                                                    │
│  • Never log secrets                                                        │
│  • Use Kubernetes secrets or external vault                                 │
│  • Rotate credentials automatically                                         │
│  • Encrypt secrets at rest                                                  │
│                                                                             │
│  WEBHOOKS                                                                   │
│  • TLS for webhook servers                                                  │
│  • Certificate rotation                                                     │
│  • Input validation                                                         │
│  • Rate limiting                                                            │
│                                                                             │
│  CONTAINER SECURITY                                                         │
│  • Non-root execution                                                       │
│  • Read-only root filesystem                                                │
│  • Minimal base images                                                      │
│  • Regular security scans                                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Compliance Requirements

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Operator Compliance Requirements                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SOC 2                               ISO 27001                              │
│  ━━━━━━━━                            ━━━━━━━━━━━                            │
│                                                                             │
│  CC6.1 - Access controls             A.9.2.2 - User access provisioning     │
│  CC6.2 - Authentication              A.9.4.1 - Information access restriction│
│  CC7.2 - System monitoring           A.12.1.2 - Change management           │
│  CC7.3 - Incident detection          A.12.4.1 - Event logging               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Decision Matrices

### 6.1 When to Use Operators

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Operator Use Case Decision Matrix                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Use Case                          │  Recommendation                       │
├────────────────────────────────────┼───────────────────────────────────────│
│  Complex stateful application      │  Operator strongly recommended        │
│  (DB, message queue, cache)        │                                       │
│  ──────────────────────────────────┼───────────────────────────────────────│
│  Simple stateless app              │  Deployment + Helm is sufficient      │
│  ──────────────────────────────────┼───────────────────────────────────────│
│  Multi-cluster management          │  Operator with federation             │
│  ──────────────────────────────────┼───────────────────────────────────────│
│  Custom business logic             │  Operator for domain-specific         │
│  ──────────────────────────────────┼───────────────────────────────────────│
│  Infrastructure provisioning       │  Operator or external tool            │
│  (cloud resources)                 │  (Crossplane, Terraform)              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Operator Best Practices Summary                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DESIGN                                                                     │
│  ✓ Single responsibility per operator                                       │
│  ✓ Idempotent reconciliation                                                │
│  ✓ Declarative API design                                                   │
│  ✓ Clear status conditions                                                  │
│  ✓ Owner references for garbage collection                                  │
│                                                                             │
│  IMPLEMENTATION                                                             │
│  ✓ Use controller-runtime framework                                         │
│  ✓ Implement proper error handling and retry                                │
│  ✓ Add comprehensive metrics                                                │
│  ✓ Use finalizers for cleanup                                               │
│  ✓ Validate with webhooks                                                   │
│                                                                             │
│  TESTING                                                                    │
│  ✓ Unit tests for reconciliation logic                                      │
│  ✓ Integration tests with envtest                                           │
│  ✓ E2E tests on real clusters                                               │
│  ✓ Test upgrade scenarios                                                   │
│                                                                             │
│  DEPLOYMENT                                                                 │
│  ✓ OLM (Operator Lifecycle Manager) packaging                               │
│  ✓ Proper RBAC configuration                                                │
│  ✓ Resource limits and quotas                                               │
│  ✓ Multi-architecture support                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Kubernetes Operator Pattern Documentation
2. Operator SDK Documentation
3. CoreOS Operators
4. Operator Framework
5. Kubernetes Controller Runtime
