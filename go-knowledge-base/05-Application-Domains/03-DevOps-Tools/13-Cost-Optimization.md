# 成本优化 (Cost Optimization)

> **分类**: 成熟应用领域
> **标签**: #cost #optimization #cloud

---

## 资源使用监控

```go
// 资源使用指标
type ResourceMetrics struct {
    CPUUsage    float64
    MemoryUsage float64
    DiskIO      float64
    NetworkIO   float64
}

func CollectMetrics() ResourceMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return ResourceMetrics{
        MemoryUsage: float64(m.Alloc) / 1024 / 1024,  // MB
    }
}

// Prometheus 导出
var (
    memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "app_memory_usage_mb",
        Help: "Current memory usage in MB",
    })
)

func recordMetrics() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    memoryUsage.Set(float64(m.Alloc) / 1024 / 1024)
}
```

---

## 自动伸缩

```go
type AutoScaler struct {
    minReplicas int
    maxReplicas int
    targetCPU   float64
    k8sClient   kubernetes.Interface
}

func (as *AutoScaler) Scale(ctx context.Context, deployment string) error {
    // 获取当前指标
    metrics, err := as.getMetrics(ctx, deployment)
    if err != nil {
        return err
    }

    // 计算所需副本数
    currentReplicas := metrics.CurrentReplicas
    cpuPercent := metrics.CPUUsage / float64(currentReplicas)

    desiredReplicas := int(cpuPercent / as.targetCPU * float64(currentReplicas))

    // 边界检查
    if desiredReplicas < as.minReplicas {
        desiredReplicas = as.minReplicas
    }
    if desiredReplicas > as.maxReplicas {
        desiredReplicas = as.maxReplicas
    }

    // 应用伸缩
    if desiredReplicas != currentReplicas {
        return as.applyScale(ctx, deployment, desiredReplicas)
    }

    return nil
}
```

---

## 资源配额管理

```go
type QuotaManager struct {
    quotas map[string]*Quota
    mu     sync.RWMutex
}

type Quota struct {
    CPU     int     // cores
    Memory  int     // MB
    Storage int     // GB
    Used    Usage
}

type Usage struct {
    CPU     int
    Memory  int
    Storage int
}

func (qm *QuotaManager) CheckQuota(team string, req Request) error {
    qm.mu.RLock()
    defer qm.mu.RUnlock()

    quota, ok := qm.quotas[team]
    if !ok {
        return fmt.Errorf("no quota for team %s", team)
    }

    if quota.Used.CPU+req.CPU > quota.CPU {
        return fmt.Errorf("CPU quota exceeded: %d/%d",
            quota.Used.CPU+req.CPU, quota.CPU)
    }

    if quota.Used.Memory+req.Memory > quota.Memory {
        return fmt.Errorf("memory quota exceeded: %d/%d",
            quota.Used.Memory+req.Memory, quota.Memory)
    }

    return nil
}
```

---

## Spot 实例利用

```go
type SpotInstanceManager struct {
    onDemandRatio float64  // 0.2 = 20% on-demand, 80% spot
}

func (sim *SpotInstanceManager) Allocate(nodes int) Allocation {
    onDemand := int(float64(nodes) * sim.onDemandRatio)
    spot := nodes - onDemand

    return Allocation{
        OnDemand: onDemand,
        Spot:     spot,
        Cost:     calculateCost(onDemand, spot),
    }
}

func calculateCost(onDemand, spot int) float64 {
    // Spot 通常便宜 60-90%
    spotPrice := 0.3  // 30% of on-demand

    onDemandCost := float64(onDemand) * 1.0
    spotCost := float64(spot) * spotPrice

    return onDemandCost + spotCost
}
```

---

## 成本告警

```go
type CostAlert struct {
    Threshold float64
    Current   float64
    Window    time.Duration
}

func (ca *CostAlert) Check(currentCost float64) (bool, string) {
    if currentCost > ca.Threshold {
        return true, fmt.Sprintf(
            "Cost alert: $%.2f exceeds threshold $%.2f",
            currentCost, ca.Threshold,
        )
    }
    return false, ""
}

func sendAlert(message string) {
    // 发送到 Slack/PagerDuty
    webhook.Send(Alert{
        Severity: "warning",
        Message:  message,
    })
}
```
