# 成本管理 (Cost Management)

> **分类**: 成熟应用领域
> **标签**: #cost #finops #optimization

---

## 资源标记

```go
type ResourceTagger struct {
    cloud CloudProvider
}

func (rt *ResourceTagger) TagResource(resourceID string, tags map[string]string) error {
    return rt.cloud.TagResource(resourceID, tags)
}

// 必需标签
var RequiredTags = []string{
    "Environment",  // prod, staging, dev
    "Team",         // 负责团队
    "Project",      // 所属项目
    "CostCenter",   // 成本中心
    "Owner",        // 负责人
}

func (rt *ResourceTagger) EnforceTags(ctx context.Context) error {
    resources, _ := rt.cloud.ListResources(ctx)

    for _, r := range resources {
        missing := rt.getMissingTags(r.Tags)
        if len(missing) > 0 {
            log.Printf("Resource %s missing tags: %v", r.ID, missing)
            // 发送告警或自动标记
        }
    }

    return nil
}
```

---

## 成本告警

```go
type CostAlert struct {
    Threshold   float64
    Period      time.Duration
    Notifier    Notifier
}

func (ca *CostAlert) Monitor(ctx context.Context) {
    ticker := time.NewTicker(ca.Period)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cost, _ := ca.getCurrentCost(ctx)

            if cost > ca.Threshold {
                ca.Notifier.Send(Notification{
                    Severity: "warning",
                    Message:  fmt.Sprintf("Cost $%.2f exceeds threshold $%.2f", cost, ca.Threshold),
                })
            }
        case <-ctx.Done():
            return
        }
    }
}
```

---

## 闲置资源检测

```go
func FindIdleResources(ctx context.Context) ([]Resource, error) {
    var idle []Resource

    // 查找低 CPU 使用的 VM
    vms, _ := cloud.ListVMs(ctx)
    for _, vm := range vms {
        metrics, _ := cloud.GetMetrics(ctx, vm.ID, "cpu", time.Hour*24*7)
        avgCPU := calculateAverage(metrics)

        if avgCPU < 5 {
            idle = append(idle, vm)
        }
    }

    // 查找未使用的磁盘
    disks, _ := cloud.ListDisks(ctx)
    for _, disk := range disks {
        if disk.AttachedTo == "" && time.Since(disk.Created) > time.Hour*24*30 {
            idle = append(idle, disk)
        }
    }

    // 查找未使用的 IP
    ips, _ := cloud.ListIPs(ctx)
    for _, ip := range ips {
        if !ip.InUse {
            idle = append(idle, ip)
        }
    }

    return idle, nil
}
```

---

## 自动优化

```go
type AutoOptimizer struct {
    policies []OptimizationPolicy
}

type OptimizationPolicy interface {
    Evaluate(ctx context.Context, resource Resource) (Recommendation, error)
}

func (ao *AutoOptimizer) Run(ctx context.Context) error {
    resources, _ := cloud.ListResources(ctx)

    for _, resource := range resources {
        for _, policy := range ao.policies {
            rec, err := policy.Evaluate(ctx, resource)
            if err != nil {
                continue
            }

            if rec.Action != NoAction {
                log.Printf("Recommendation for %s: %s (savings: $%.2f/month)",
                    resource.ID, rec.Action, rec.Savings)
            }
        }
    }

    return nil
}

// 右 size 策略
func RightSizePolicy(ctx context.Context, vm VM) (Recommendation, error) {
    metrics, _ := cloud.GetMetrics(ctx, vm.ID, "cpu", time.Hour*24*7)
    maxCPU := calculateMax(metrics)

    if maxCPU < 20 {
        return Recommendation{
            Action:  "downsize",
            Target:  getSmallerInstanceType(vm.Type),
            Savings: calculateSavings(vm.Type, getSmallerInstanceType(vm.Type)),
        }, nil
    }

    return Recommendation{Action: NoAction}, nil
}
```
