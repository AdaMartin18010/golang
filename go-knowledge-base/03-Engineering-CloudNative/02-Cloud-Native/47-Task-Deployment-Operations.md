# 任务部署运维 (Task Deployment Operations)

> **分类**: 工程与云原生
> **标签**: #deployment #operations #devops

---

## Kubernetes 部署

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-scheduler
spec:
  replicas: 3
  selector:
    matchLabels:
      app: task-scheduler
  template:
    metadata:
      labels:
        app: task-scheduler
    spec:
      containers:
      - name: scheduler
        image: task-scheduler:latest
        ports:
        - containerPort: 8080
        env:
        - name: TASK_WORKERS
          value: "50"
        - name: TASK_QUEUE_SIZE
          value: "10000"
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: task-scheduler-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: task-scheduler
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Pods
    pods:
      metric:
        name: task_queue_depth
      target:
        type: AverageValue
        averageValue: "1000"
```

---

## 健康检查

```go
// 健康检查处理器
type HealthChecker struct {
    checks map[string]HealthCheck
}

type HealthCheck func(ctx context.Context) error

func (hc *HealthChecker) Register(name string, check HealthCheck) {
    hc.checks[name] = check
}

func (hc *HealthChecker) LiveEndpoint(c *gin.Context) {
    // 简单存活检查
    c.JSON(200, gin.H{"status": "ok"})
}

func (hc *HealthChecker) ReadyEndpoint(c *gin.Context) {
    ctx := c.Request.Context()
    results := make(map[string]string)

    allReady := true
    for name, check := range hc.checks {
        if err := check(ctx); err != nil {
            results[name] = fmt.Sprintf("unhealthy: %v", err)
            allReady = false
        } else {
            results[name] = "healthy"
        }
    }

    if allReady {
        c.JSON(200, gin.H{
            "status": "ready",
            "checks": results,
        })
    } else {
        c.JSON(503, gin.H{
            "status": "not ready",
            "checks": results,
        })
    }
}

// 具体检查
func DatabaseHealthCheck(db *sql.DB) HealthCheck {
    return func(ctx context.Context) error {
        return db.PingContext(ctx)
    }
}

func QueueHealthCheck(queue TaskQueue) HealthCheck {
    return func(ctx context.Context) error {
        // 检查队列连接
        return queue.HealthCheck(ctx)
    }
}

func WorkerPoolHealthCheck(pool *WorkerPool) HealthCheck {
    return func(ctx context.Context) error {
        stats := pool.Stats()
        if stats.IdleWorkers == 0 && stats.PendingTasks > 100 {
            return fmt.Errorf("worker pool saturated")
        }
        return nil
    }
}
```

---

## 金丝雀部署

```go
type CanaryDeployer struct {
    k8sClient kubernetes.Interface
}

func (cd *CanaryDeployer) Deploy(ctx context.Context, deployment DeploymentConfig) error {
    // 1. 部署 Canary 版本 (10% 流量)
    canaryDeployment := deployment.Clone()
    canaryDeployment.Name = fmt.Sprintf("%s-canary", deployment.Name)
    canaryDeployment.Replicas = int32(float64(deployment.Replicas) * 0.1)
    canaryDeployment.Image = deployment.NewImage

    if err := cd.applyDeployment(ctx, canaryDeployment); err != nil {
        return err
    }

    // 2. 等待 Canary 就绪
    if err := cd.waitForReady(ctx, canaryDeployment.Name, 5*time.Minute); err != nil {
        cd.rollback(ctx, canaryDeployment.Name)
        return err
    }

    // 3. 监控指标
    if err := cd.monitorCanary(ctx, canaryDeployment.Name, 10*time.Minute); err != nil {
        cd.rollback(ctx, canaryDeployment.Name)
        return err
    }

    // 4. 全量部署
    deployment.Image = deployment.NewImage
    if err := cd.applyDeployment(ctx, deployment); err != nil {
        return err
    }

    // 5. 清理 Canary
    cd.deleteDeployment(ctx, canaryDeployment.Name)

    return nil
}

func (cd *CanaryDeployer) monitorCanary(ctx context.Context, name string, duration time.Duration) error {
    ctx, cancel := context.WithTimeout(ctx, duration)
    defer cancel()

    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            metrics, err := cd.getMetrics(ctx, name)
            if err != nil {
                return err
            }

            // 检查错误率
            if metrics.ErrorRate > 0.01 {
                return fmt.Errorf("canary error rate too high: %.2f%%", metrics.ErrorRate*100)
            }

            // 检查延迟
            if metrics.P99Latency > 500*time.Millisecond {
                return fmt.Errorf("canary p99 latency too high: %v", metrics.P99Latency)
            }
        }
    }
}
```

---

## 运维命令

```go
// 运维 CLI
var opsCmd = &cobra.Command{
    Use:   "ops",
    Short: "Operations commands",
}

func init() {
    opsCmd.AddCommand(
        drainCmd(),
        pauseCmd(),
        resumeCmd(),
        scaleCmd(),
        migrateCmd(),
    )
}

func drainCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "drain <node>",
        Short: "Drain tasks from a node",
        Run: func(cmd *cobra.Command, args []string) {
            node := args[0]
            client := NewOpsClient()

            // 优雅排空
            if err := client.DrainNode(cmd.Context(), node, DrainOptions{
                GracePeriod: 5 * time.Minute,
                Force:       false,
            }); err != nil {
                log.Fatal(err)
            }

            fmt.Printf("Node %s drained successfully\n", node)
        },
    }
}

func pauseCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "pause",
        Short: "Pause task scheduling",
        Run: func(cmd *cobra.Command, args []string) {
            client := NewOpsClient()

            if err := client.PauseScheduling(cmd.Context()); err != nil {
                log.Fatal(err)
            }

            fmt.Println("Task scheduling paused")
        },
    }
}

func scaleCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "scale <workers>",
        Short: "Scale worker pool",
        Run: func(cmd *cobra.Command, args []string) {
            workers, _ := strconv.Atoi(args[0])
            client := NewOpsClient()

            if err := client.ScaleWorkers(cmd.Context(), workers); err != nil {
                log.Fatal(err)
            }

            fmt.Printf("Worker pool scaled to %d\n", workers)
        },
    }
}
```
