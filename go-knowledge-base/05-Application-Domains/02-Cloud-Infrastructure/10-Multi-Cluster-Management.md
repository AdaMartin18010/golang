# 多集群管理 (Multi-Cluster Management)

> **分类**: 成熟应用领域
> **标签**: #multi-cluster #kubernetes #federation

---

## 集群联邦

### 集群注册

```go
type ClusterManager struct {
    clusters map[string]*Cluster
    mu       sync.RWMutex
}

type Cluster struct {
    Name       string
    KubeConfig *rest.Config
    Client     kubernetes.Interface
    Region     string
    Labels     map[string]string
}

func (cm *ClusterManager) Register(name string, kubeconfig []byte) error {
    config, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
    if err != nil {
        return err
    }

    restConfig, err := config.ClientConfig()
    if err != nil {
        return err
    }

    client, err := kubernetes.NewForConfig(restConfig)
    if err != nil {
        return err
    }

    cm.mu.Lock()
    defer cm.mu.Unlock()

    cm.clusters[name] = &Cluster{
        Name:       name,
        KubeConfig: restConfig,
        Client:     client,
    }

    return nil
}

func (cm *ClusterManager) GetCluster(name string) (*Cluster, error) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    cluster, ok := cm.clusters[name]
    if !ok {
        return nil, fmt.Errorf("cluster %s not found", name)
    }
    return cluster, nil
}
```

---

## 跨集群部署

```go
func (cm *ClusterManager) DeployToAll(ctx context.Context, deployment *appsv1.Deployment) error {
    cm.mu.RLock()
    clusters := make(map[string]*Cluster)
    for k, v := range cm.clusters {
        clusters[k] = v
    }
    cm.mu.RUnlock()

    var wg sync.WaitGroup
    errChan := make(chan error, len(clusters))

    for name, cluster := range clusters {
        wg.Add(1)
        go func(n string, c *Cluster) {
            defer wg.Done()

            _, err := c.Client.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
            if err != nil {
                errChan <- fmt.Errorf("cluster %s: %w", n, err)
            }
        }(name, cluster)
    }

    wg.Wait()
    close(errChan)

    var errors []error
    for err := range errChan {
        errors = append(errors, err)
    }

    if len(errors) > 0 {
        return fmt.Errorf("deployment failed: %v", errors)
    }

    return nil
}
```

---

## 全局负载均衡

```go
type GlobalLoadBalancer struct {
    clusters   []*Cluster
    healthCheck HealthChecker
}

func (glb *GlobalLoadBalancer) SelectCluster(ctx context.Context, req Request) (*Cluster, error) {
    healthy := glb.getHealthyClusters()
    if len(healthy) == 0 {
        return nil, errors.New("no healthy clusters")
    }

    // 根据地理位置选择
    if req.Location != "" {
        for _, c := range healthy {
            if c.Region == req.Location {
                return c, nil
            }
        }
    }

    // 根据负载选择
    var bestCluster *Cluster
    var minLoad float64 = math.MaxFloat64

    for _, c := range healthy {
        load, err := glb.getClusterLoad(ctx, c)
        if err != nil {
            continue
        }
        if load < minLoad {
            minLoad = load
            bestCluster = c
        }
    }

    return bestCluster, nil
}

func (glb *GlobalLoadBalancer) getClusterLoad(ctx context.Context, c *Cluster) (float64, error) {
    nodes, err := c.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
    if err != nil {
        return 0, err
    }

    var totalCPU, usedCPU int64
    for _, node := range nodes.Items {
        totalCPU += node.Status.Allocatable.Cpu().MilliValue()
        usedCPU += node.Status.Capacity.Cpu().MilliValue() - node.Status.Allocatable.Cpu().MilliValue()
    }

    if totalCPU == 0 {
        return 0, nil
    }

    return float64(usedCPU) / float64(totalCPU), nil
}
```

---

## 故障转移

```go
type FailoverManager struct {
    primary   *Cluster
    secondaries []*Cluster
}

func (fm *FailoverManager) MonitorAndFailover(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := fm.healthCheck(ctx, fm.primary); err != nil {
                log.Printf("Primary cluster unhealthy: %v", err)
                fm.performFailover(ctx)
            }
        case <-ctx.Done():
            return
        }
    }
}

func (fm *FailoverManager) performFailover(ctx context.Context) {
    for _, secondary := range fm.secondaries {
        if err := fm.healthCheck(ctx, secondary); err != nil {
            continue
        }

        log.Printf("Failing over to %s", secondary.Name)
        fm.promoteToPrimary(secondary)
        return
    }

    log.Fatal("No healthy secondary cluster available")
}
```
