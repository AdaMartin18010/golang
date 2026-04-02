# 服务网格控制面 (Service Mesh Control Plane)

> **分类**: 成熟应用领域
> **标签**: #servicemesh #controlplane #xds

---

## xDS API 实现

### 控制面架构

```go
// Discovery Server
type DiscoveryServer struct {
    snapshots map[string]*cache.Snapshot
    cache     cache.SnapshotCache
}

func (s *DiscoveryServer) Start() {
    ctx := context.Background()

    // 启动 gRPC 服务
    grpcServer := grpc.NewServer()

    // 注册 xDS 服务
    discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, s)

    // 监听
    lis, _ := net.Listen("tcp", ":18000")
    grpcServer.Serve(lis)
}
```

---

## 配置推送

### CDS (Cluster Discovery Service)

```go
func makeCluster(clusterName string) *cluster.Cluster {
    return &cluster.Cluster{
        Name:                 clusterName,
        ConnectTimeout:       durationpb.New(5 * time.Second),
        ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
        EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
            ServiceName: clusterName,
            EdsConfig: &core.ConfigSource{
                ConfigSourceSpecifier: &core.ConfigSource_Ads{
                    Ads: &core.AggregatedConfigSource{},
                },
            },
        },
    }
}
```

### RDS (Route Discovery Service)

```go
func makeRoute(routeName, clusterName string) *route.RouteConfiguration {
    return &route.RouteConfiguration{
        Name: routeName,
        VirtualHosts: []*route.VirtualHost{{
            Name:    "local_service",
            Domains: []string{"*"},
            Routes: []*route.Route{{
                Match: &route.RouteMatch{
                    PathSpecifier: &route.RouteMatch_Prefix{
                        Prefix: "/",
                    },
                },
                Action: &route.Route_Route{
                    Route: &route.RouteAction{
                        ClusterSpecifier: &route.RouteAction_Cluster{
                            Cluster: clusterName,
                        },
                    },
                },
            }},
        }},
    }
}
```

---

## 增量更新

```go
func (s *DiscoveryServer) PushSnapshot(node string, snapshot *cache.Snapshot) {
    s.cache.SetSnapshot(context.Background(), node, snapshot)
}

// 当配置变更时
func (s *DiscoveryServer) OnConfigChange() {
    // 构建新的 snapshot
    snapshot := &cache.Snapshot{
        Resources: [7]cache.Resources{
            {Version: "1", Items: map[string]cache.Resource{
                "cluster1": makeCluster("cluster1"),
            }},
            // ... 其他资源
        },
    }

    // 推送到所有节点
    for node := range s.snapshots {
        s.PushSnapshot(node, snapshot)
    }
}
```
