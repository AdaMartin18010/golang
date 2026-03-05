# 第十三章：云原生基础设施

> Go 在云原生生态系统中的核心项目和应用

---

## 13.1 CNCF 项目全景

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CNCF 项目全景 - Go 语言主导                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  层级                                                                        │
│  ────                                                                        │
│                                                                             │
│  应用定义与镜像  │ Helm(100%), OPA(50%), Kustomize(50%)                      │
│  ───────────────┼────────────────────────────────────────────────────────   │
│                                                                             │
│  编排与管理     │ Kubernetes(100%), etcd(100%), Consul(100%)                 │
│  ───────────────┼────────────────────────────────────────────────────────   │
│                                                                             │
│  运行时         │ containerd(100%), rkt(100%), CRI-O(100%)                   │
│  ───────────────┼────────────────────────────────────────────────────────   │
│                                                                             │
│  配置           │ Vitess(100%), gRPC(100%)                                   │
│  ───────────────┼────────────────────────────────────────────────────────   │
│                                                                             │
│  可观测性       │ Prometheus(100%), Jaeger(100%), Fluentd(30%)               │
│  ───────────────┼────────────────────────────────────────────────────────   │
│                                                                             │
│  平台           │ OpenShift(部分), Cloud Foundry(部分)                        │
│                                                                             │
│  注：百分比表示 Go 代码占比                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 13.2 Kubernetes 深度解析

### 13.2.1 控制器模式实现

```go
// Kubernetes 控制器模式

// 1. Informers - 高效监听资源变化
func NewPodInformer(client kubernetes.Interface) cache.SharedIndexInformer {
    listWatcher := cache.NewListWatchFromClient(
        client.CoreV1().RESTClient(),
        "pods",
        v1.NamespaceAll,
        fields.Everything(),
    )

    return cache.NewSharedIndexInformer(
        listWatcher,
        &v1.Pod{},
        0, // 不缓存
        cache.Indexers{},
    )
}

// 2. WorkQueue - 可靠的任务队列
func processQueue(queue workqueue.RateLimitingInterface, handler func(string) error) {
    for {
        key, shutdown := queue.Get()
        if shutdown {
            return
        }

        err := handler(key.(string))
        if err != nil {
            queue.AddRateLimited(key)
        } else {
            queue.Forget(key)
        }
        queue.Done(key)
    }
}

// 3. 完整控制器示例

type PodController struct {
    clientset kubernetes.Interface
    informer  cache.SharedIndexInformer
    queue     workqueue.RateLimitingInterface
}

func (c *PodController) Run(workers int, stopCh <-chan struct{}) {
    defer c.queue.ShutDown()

    go c.informer.Run(stopCh)
    if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
        return
    }

    for i := 0; i < workers; i++ {
        go wait.Until(c.runWorker, time.Second, stopCh)
    }

    <-stopCh
}

func (c *PodController) runWorker() {
    for c.processNextItem() {
    }
}

func (c *PodController) processNextItem() bool {
    key, quit := c.queue.Get()
    if quit {
        return false
    }
    defer c.queue.Done(key)

    err := c.syncHandler(key.(string))
    if err != nil {
        c.queue.AddRateLimited(key)
    } else {
        c.queue.Forget(key)
    }

    return true
}

func (c *PodController) syncHandler(key string) error {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return err
    }

    pod, err := c.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
    if err != nil {
        return err
    }

    // 业务逻辑处理
    return c.handlePod(pod)
}
```

### 13.2.2 Operator 开发

```go
// Operator SDK 示例

// 1. 定义 CRD
type DatabaseSpec struct {
    Size     int32  `json:"size"`
    Version  string `json:"version"`
    Storage  string `json:"storage"`
}

type DatabaseStatus struct {
    Nodes []string `json:"nodes"`
    Phase string   `json:"phase"`
}

// 2. 控制器实现
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("database", req.NamespacedName)

    // 获取 CR
    db := &databasev1.Database{}
    if err := r.Get(ctx, req.NamespacedName, db); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // 创建/更新 StatefulSet
    sts := r.desiredStatefulSet(db)
    if err := ctrl.SetControllerReference(db, sts, r.Scheme); err != nil {
        return ctrl.Result{}, err
    }

    // 应用资源
    if err := r.applyResource(ctx, sts); err != nil {
        return ctrl.Result{}, err
    }

    // 更新状态
    db.Status.Phase = "Running"
    if err := r.Status().Update(ctx, db); err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}
```

---

## 13.3 容器运行时

### 13.3.1 containerd 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      containerd 架构                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Client (ctr/cri)                                                           │
│       │                                                                     │
│       ▼                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          containerd                                  │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌───────────┐  │   │
│  │  │   Images    │  │   Content   │  │  Snapshots  │  │  Runtime  │  │   │
│  │  │   Service   │  │   Service   │  │   Service   │  │  Service  │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └───────────┘  │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │   │
│  │  │                     GRPC API                                     │ │   │
│  │  └─────────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│       │                                                                     │
│       ▼                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                                  │
│  │ runc     │  │ gVisor   │  │ Kata     │  (OCI Runtimes)                  │
│  └──────────┘  └──────────┘  └──────────┘                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 13.3.2 与 containerd 交互

```go
// containerd 客户端
import "github.com/containerd/containerd"

func containerdExample() {
    // 连接 containerd
    client, err := containerd.New("/run/containerd/containerd.sock")
    if err != nil {
        panic(err)
    }
    defer client.Close()

    ctx := context.Background()

    // 拉取镜像
    image, err := client.Pull(ctx, "docker.io/library/nginx:latest",
        containerd.WithPullUnpack)
    if err != nil {
        panic(err)
    }

    // 创建容器
    container, err := client.NewContainer(ctx, "nginx",
        containerd.WithImage(image),
        containerd.WithNewSnapshot("nginx-snapshot", image),
        containerd.WithNewSpec(oci.WithImageConfig(image)),
    )
    if err != nil {
        panic(err)
    }

    // 创建任务（进程）
    task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
    if err != nil {
        panic(err)
    }

    // 启动
    if err := task.Start(ctx); err != nil {
        panic(err)
    }

    // 等待退出
    exitStatus, err := task.Wait(ctx)
    <-exitStatus
}
```

---

## 13.4 可观测性栈

### 13.4.1 Prometheus 指标

```go
// 自定义指标
var (
    // Counter
    requestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "app_requests_total",
            Help: "Total number of requests",
        },
        []string{"method", "endpoint", "status"},
    )

    // Histogram
    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "app_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    // Gauge
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "app_active_connections",
            Help: "Number of active connections",
        },
    )

    // Summary
    requestSize = promauto.NewSummaryVec(
        prometheus.SummaryOpts{
            Name:       "app_request_size_bytes",
            Help:       "Request size in bytes",
            Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
        },
        []string{"method"},
    )
)

// 中间件
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        activeConnections.Inc()
        defer activeConnections.Dec()

        wrapped := &responseWriter{w, http.StatusOK}
        next.ServeHTTP(wrapped, r)

        duration := time.Since(start).Seconds()
        status := strconv.Itoa(wrapped.statusCode)

        requestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
        requestSize.WithLabelValues(r.Method).Observe(float64(r.ContentLength))
    })
}
```

### 13.4.2 分布式追踪

```go
// OpenTelemetry 集成

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
    "go.opentelemetry.io/otel/trace"
)

func initTracer() (*sdktrace.TracerProvider, error) {
    // 创建 Jaeger 导出器
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://localhost:14268/api/traces"),
    ))
    if err != nil {
        return nil, err
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName("my-service"),
            attribute.String("environment", "production"),
        )),
    )

    otel.SetTracerProvider(tp)
    return tp, nil
}

// 使用追踪
var tracer = otel.Tracer("my-service")

func processOrder(ctx context.Context, order Order) error {
    ctx, span := tracer.Start(ctx, "processOrder",
        trace.WithAttributes(
            attribute.String("order.id", order.ID),
            attribute.Float64("order.amount", order.Amount),
        ),
    )
    defer span.End()

    // 数据库操作
    if err := saveToDB(ctx, order); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "db error")
        return err
    }

    // 调用外部服务
    if err := callPaymentService(ctx, order); err != nil {
        span.RecordError(err)
        return err
    }

    return nil
}
```

### 13.4.3 日志聚合

```go
// 结构化日志与 Fluentd/Fluent Bit 集成

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

func initLogger() *zap.Logger {
    // 生产环境配置
    config := zap.NewProductionConfig()
    config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

    // 自定义编码
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    // 文件输出 + 轮转
    w := zapcore.AddSync(&lumberjack.Logger{
        Filename:   "/var/log/app/app.log",
        MaxSize:    100, // MB
        MaxBackups: 10,
        MaxAge:     30, // days
        Compress:   true,
    })

    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(config.EncoderConfig),
        w,
        config.Level,
    )

    return zap.New(core, zap.AddCaller())
}

// 使用
var logger *zap.Logger

func handleRequest(w http.ResponseWriter, r *http.Request) {
    logger.Info("request received",
        zap.String("method", r.Method),
        zap.String("path", r.URL.Path),
        zap.String("remote_addr", r.RemoteAddr),
        zap.String("trace_id", r.Header.Get("X-Trace-ID")),
    )
}
```

---

## 13.5 服务网格

### 13.5.1 Istio 与 Go

```go
// 使用 Istio 的 Go 应用最佳实践

// 1. 健康检查
func healthCheck(w http.ResponseWriter, r *http.Request) {
    // Istio 默认检查 /healthz
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}

// 2. 优雅关闭
func gracefulShutdown(server *http.Server) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

    <-quit
    log.Println("Shutting down server...")

    // Istio 需要延迟退出，让 envoy 完成转发
    time.Sleep(5 * time.Second)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(err)
    }
}

// 3. 分布式追踪上下文传递
func callService(ctx context.Context, service string) error {
    // Istio 自动注入追踪头
    req, _ := http.NewRequestWithContext(ctx, "GET", service, nil)

    // 确保传播追踪上下文
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}
```

---

## 13.6 GitOps 与持续交付

```go
// ArgoCD / Flux 风格的 GitOps 控制器

// 1. 监听 Git 仓库变化
type GitOpsReconciler struct {
    client.Client
    GitClient *git.Client
}

func (r *GitOpsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    app := &v1.Application{}
    if err := r.Get(ctx, req.NamespacedName, app); err != nil {
        return ctrl.Result{}, err
    }

    // 获取 Git 仓库最新状态
    commit, err := r.GitClient.Clone(app.Spec.RepoURL, app.Spec.Branch)
    if err != nil {
        return ctrl.Result{}, err
    }

    // 解析 manifest
    resources, err := r.parseManifests(commit)
    if err != nil {
        return ctrl.Result{}, err
    }

    // 同步到集群
    for _, res := range resources {
        if err := r.applyResource(ctx, res); err != nil {
            return ctrl.Result{}, err
        }
    }

    // 更新同步状态
    app.Status.SyncStatus = v1.SyncStatusSynced
    app.Status.LastSync = metav1.Now()

    return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// 2. 渐进式交付 (Argo Rollouts)
type CanaryStrategy struct {
    Steps []CanaryStep
}

type CanaryStep struct {
    SetWeight int
    Pause     PauseConfig
    Analysis  AnalysisConfig
}

func (s *CanaryStrategy) Execute(ctx context.Context, rollout *v1.Rollout) error {
    for _, step := range s.Steps {
        // 调整流量权重
        if err := s.setWeight(rollout, step.SetWeight); err != nil {
            return err
        }

        // 分析指标
        if step.Analysis.Enabled {
            if err := s.runAnalysis(ctx, rollout, step.Analysis); err != nil {
                // 回滚
                return s.rollback(rollout)
            }
        }

        // 暂停等待
        if step.Pause.Duration > 0 {
            time.Sleep(step.Pause.Duration)
        }
    }

    // 完成金丝雀，全量发布
    return s.promote(rollout)
}
```

---

## 13.7 安全实践

```go
// Pod 安全策略

// 1. 非 root 用户运行
securityContext := &v1.SecurityContext{
    RunAsNonRoot:             pointer.Bool(true),
    RunAsUser:                pointer.Int64(1000),
    RunAsGroup:               pointer.Int64(1000),
    ReadOnlyRootFilesystem:   pointer.Bool(true),
    AllowPrivilegeEscalation: pointer.Bool(false),
    Capabilities: &v1.Capabilities{
        Drop: []v1.Capability{"ALL"},
    },
}

// 2. 网络策略
type NetworkPolicy struct {
    // 只允许特定流量
    Ingress: []IngressRule{
        {
            From: []NetworkPolicyPeer{
                {NamespaceSelector: &metav1.LabelSelector{
                    MatchLabels: map[string]string{"name": "frontend"},
                }},
            },
            Ports: []NetworkPolicyPort{
                {Protocol: &tcp, Port: &intstr.IntOrString{IntVal: 8080}},
            },
        },
    },
}

// 3. 密钥管理 (Vault 集成)
func getSecretFromVault(path string) (map[string]string, error) {
    config := vault.DefaultConfig()
    config.Address = "http://vault:8200"

    client, err := vault.NewClient(config)
    if err != nil {
        return nil, err
    }

    // 使用 Kubernetes 认证
    client.SetAuthMethod(vault.KubernetesAuth("my-app"))

    secret, err := client.Logical().Read(path)
    if err != nil {
        return nil, err
    }

    return secret.Data, nil
}
```

---

## 13.8 云原生 Go 应用最佳实践

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    云原生 Go 应用 12 要素                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. 基准代码         一份代码库，多份部署                                    │
│  2. 依赖             显式声明依赖 (go.mod)                                   │
│  3. 配置             环境变量存储配置                                        │
│  4. 后端服务         把后端服务当作附加资源                                  │
│  5. 构建、发布、运行  严格分离构建和运行                                      │
│  6. 进程             以一个或多个无状态进程运行应用                          │
│  7. 端口绑定         通过端口绑定提供服务                                    │
│  8. 并发             通过进程模型进行扩展                                    │
│  9. 易处理           快速启动和优雅终止                                      │
│  10. 开发环境与线上  尽可能的保持开发、预发布、线上环境相同                    │
│  11. 日志            把日志当作事件流                                        │
│  12. 管理进程        后台管理任务当作一次性进程运行                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*本章涵盖了 Go 在云原生生态系统中的核心应用，从容器编排到可观测性，从服务网格到安全实践。*
