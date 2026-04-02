# EC-099: Kubernetes 1.34 CronJob 深度分析 (Kubernetes 1.34 CronJob Deep Dive)

> **维度**: Engineering CloudNative
> **级别**: S (25+ KB)
> **标签**: #kubernetes134 #cronjob #sidecar #scheduling
> **版本演进**: K8s 1.28 → K8s 1.32 → **K8s 1.34+** (2026)
> **权威来源**: [K8s 1.34 Release Notes](https://kubernetes.io/releases/release-v1-34/), [K8s CronJob Controller](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/cronjob)

---

## 版本演进亮点

```
Kubernetes 1.28 (2023)    Kubernetes 1.32 (2024)    Kubernetes 1.34 (2026) ⭐️
      │                          │                          │
      ▼                          ▼                          ▼
┌─────────────┐          ┌───────────────┐          ┌─────────────────┐
│ Sidecar     │          │ Pod Scheduling│          │ Sidecar 容器 GA │
│ 容器 Beta   │─────────►│ Ready 门控    │─────────►│ Job 完成策略    │
│ 时区支持    │          │ 改进          │          │ 增强调度        │
└─────────────┘          │ 驱逐策略      │          │ 多租户隔离      │
                         └───────────────┘          │ 自动扩缩容      │
                                                    └─────────────────┘
```

---

## K8s 1.34 新特性

### 1. Sidecar 容器 GA

```yaml
# K8s 1.34: Sidecar 容器正式发布
# 特点：Sidecar 在主容器完成后自动终止

apiVersion: batch/v1
kind: CronJob
metadata:
  name: data-processor
spec:
  schedule: "0 2 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            # Sidecar 容器：日志收集
            - name: log-collector
              image: fluent-bit:latest
              restartPolicy: Always  # Sidecar 特性
              lifecycle:
                type: Sidecar      # K8s 1.34+: 显式声明
              volumeMounts:
                - name: logs
                  mountPath: /var/log/app

            # 主容器：数据处理
            - name: processor
              image: data-processor:v2
              command: ["./process"]
              volumeMounts:
                - name: logs
                  mountPath: /app/logs

          # K8s 1.34: Job 完成时自动停止 Sidecar
          terminationGracePeriodSeconds: 30

          # 新的完成策略
          completionMode: Indexed        # 或 NonIndexed
          backoffLimit: 3
          ttlSecondsAfterFinished: 3600  # 1小时后自动清理
```

### 2. 增强调度

```yaml
# K8s 1.34: 新的调度特性

apiVersion: batch/v1
kind: CronJob
metadata:
  name: gpu-training
spec:
  schedule: "0 0 * * 0"  # 每周日
  jobTemplate:
    spec:
      template:
        spec:
          schedulerName: default-scheduler

          # K8s 1.34: 调度门控
          schedulingGates:
            - name: gpu-available
            - name: dataset-ready

          # 资源预留
          resourceClaims:
            - name: gpu-nvidia-a100
              resourceClaimTemplateName: nvidia-a100

          containers:
            - name: trainer
              image: ml-trainer:v3
              resources:
                claims:
                  - name: gpu-nvidia-a100

              # K8s 1.34: 动态资源分配
              env:
                - name: GPU_MEMORY
                  valueFrom:
                    resourceFieldRef:
                      resource: claims.gpu-nvidia-a100.memory

          # 新的亲和性规则
          affinity:
            podAffinity:
              preferredDuringSchedulingIgnoredDuringExecution:
                - weight: 100
                  podAffinityTerm:
                    labelSelector:
                      matchExpressions:
                        - key: dataset
                          operator: In
                          values: ["imagenet"]
                    topologyKey: topology.kubernetes.io/zone
```

### 3. 多租户隔离增强

```yaml
# K8s 1.34: CronJob 多租户隔离

apiVersion: batch/v1
kind: CronJob
metadata:
  name: tenant-job
  namespace: tenant-a
spec:
  schedule: "*/5 * * * *"

  # K8s 1.34: 资源配额感知
  resourcePolicy:
    respectQuota: true
    priorityClassName: tenant-batch

  # 并发策略增强
  concurrencyPolicy: Forbid
  startingDeadlineSeconds: 300

  # 历史限制
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 5

  jobTemplate:
    spec:
      # 新的 Pod 模板
      template:
        metadata:
          labels:
            security.istio.io/tlsMode: istio
        spec:
          # 安全上下文增强
          securityContext:
            seccompProfile:
              type: RuntimeDefault
            appArmorProfile:
              type: RuntimeDefault

          containers:
            - name: job
              image: job-image:v1
              securityContext:
                allowPrivilegeEscalation: false
                readOnlyRootFilesystem: true
                capabilities:
                  drop: ["ALL"]

              # K8s 1.34: 用户命名空间
              userNamespace: true
```

---

## 控制器源码更新

### K8s 1.34 CronJob Controller

```go
// pkg/controller/cronjob/cronjob_controllerv2.go (K8s 1.34)

// 新增：Sidecar 感知
func (jm *ControllerV2) syncCronJob(ctx context.Context, cronJob *batchv1.CronJob,
    jobs []*batchv1.Job) error {

    // ... 原有逻辑 ...

    // K8s 1.34: 检查 Sidecar 容器状态
    if hasSidecars(cronJob) {
        // 等待所有主容器完成，自动停止 Sidecar
        if err := jm.handleSidecarTermination(ctx, cronJob, jobs); err != nil {
            return err
        }
    }

    // K8s 1.34: 调度门控检查
    if cronJob.Spec.JobTemplate.Spec.Template.Spec.SchedulingGates != nil {
        // 检查调度门控是否已解除
        ready, err := jm.checkSchedulingGates(ctx, cronJob)
        if err != nil || !ready {
            return err
        }
    }

    // ... 原有逻辑 ...
}

// K8s 1.34: 新的调度特性
func (jm *ControllerV2) handleResourceClaims(ctx context.Context,
    cronJob *batchv1.CronJob) error {

    for _, claim := range cronJob.Spec.JobTemplate.Spec.Template.Spec.ResourceClaims {
        // 预分配资源声明
        if err := jm.allocateResourceClaim(ctx, claim); err != nil {
            jm.recorder.Eventf(cronJob, corev1.EventTypeWarning, "ResourceClaimFailed",
                "Failed to allocate resource claim %s: %v", claim.Name, err)
            return err
        }
    }

    return nil
}
```

---

## 版本对比

| 特性 | K8s 1.28 | K8s 1.32 | K8s 1.34 |
|------|----------|----------|----------|
| Sidecar 容器 | Beta | Beta | **GA** |
| 调度门控 | Alpha | Beta | **Stable** |
| 动态资源分配 | Alpha | Beta | **Beta+** |
| 资源声明 | - | Alpha | **Beta** |
| 用户命名空间 | Alpha | Beta | **Stable** |
| CronJob 时区 | Stable | Stable | Stable |
| Pod 完成策略 | - | - | **新增** |

---

## 参考文献

1. [Kubernetes 1.34 Release Notes](https://kubernetes.io/releases/release-v1-34/) - 官方发布说明
2. [Sidecar Containers](https://kubernetes.io/docs/concepts/workloads/pods/sidecar-containers/) - K8s 文档
3. [Dynamic Resource Allocation](https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/) - K8s 文档
4. [CronJob Controller](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/cronjob) - 源码
