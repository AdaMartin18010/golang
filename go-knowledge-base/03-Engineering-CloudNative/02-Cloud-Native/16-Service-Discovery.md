# 服务发现 (Service Discovery)

> **分类**: 工程与云原生  
> **标签**: #service-discovery #consul #etcd

---

## Consul 集成

```go
import "github.com/hashicorp/consul/api"

func RegisterService(consulAddr, serviceName string, port int) error {
    config := api.DefaultConfig()
    config.Address = consulAddr
    
    client, err := api.NewClient(config)
    if err != nil {
        return err
    }
    
    // 获取本地 IP
    ip, _ := getLocalIP()
    
    // 注册服务
    registration := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s", serviceName, ip),
        Name:    serviceName,
        Address: ip,
        Port:    port,
        Tags:    []string{"go", "api"},
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", ip, port),
            Interval: "10s",
            Timeout:  "5s",
        },
    }
    
    return client.Agent().ServiceRegister(registration)
}

func DeregisterService(consulAddr, serviceID string) error {
    config := api.DefaultConfig()
    config.Address = consulAddr
    
    client, _ := api.NewClient(config)
    return client.Agent().ServiceDeregister(serviceID)
}

func DiscoverService(consulAddr, serviceName string) ([]*api.ServiceEntry, error) {
    config := api.DefaultConfig()
    config.Address = consulAddr
    
    client, _ := api.NewClient(config)
    
    services, _, err := client.Health().Service(serviceName, "", true, nil)
    return services, err
}
```

---

## 客户端负载均衡

```go
type ServiceResolver struct {
    consul *api.Client
    cache  map[string][]*api.ServiceEntry
    mu     sync.RWMutex
}

func (sr *ServiceResolver) Resolve(serviceName string) (*api.ServiceEntry, error) {
    sr.mu.RLock()
    entries := sr.cache[serviceName]
    sr.mu.RUnlock()
    
    if len(entries) == 0 {
        var err error
        entries, _, err = sr.consul.Health().Service(serviceName, "", true, nil)
        if err != nil {
            return nil, err
        }
        
        sr.mu.Lock()
        sr.cache[serviceName] = entries
        sr.mu.Unlock()
    }
    
    if len(entries) == 0 {
        return nil, fmt.Errorf("no healthy instances of %s", serviceName)
    }
    
    // 轮询
    idx := rand.Intn(len(entries))
    return entries[idx], nil
}
```

---

## 基于 DNS 的发现

```go
func DiscoverWithDNS(serviceName string) ([]string, error) {
    // 查询 Consul DNS
    records, err := net.LookupSRV("", "", fmt.Sprintf("%s.service.consul", serviceName))
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    for _, record := range records {
        addresses = append(addresses, fmt.Sprintf("%s:%d", record.Target, record.Port))
    }
    
    return addresses, nil
}
```

---

## Kubernetes 服务发现

```go
func DiscoverK8SService(namespace, serviceName string) (string, error) {
    // 集群内 DNS
    host := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)
    
    // 或使用环境变量
    // host := os.Getenv(fmt.Sprintf("%s_SERVICE_HOST", strings.ToUpper(serviceName)))
    // port := os.Getenv(fmt.Sprintf("%s_SERVICE_PORT", strings.ToUpper(serviceName)))
    
    return host, nil
}
```

---

## 健康检查集成

```go
type HealthAwareResolver struct {
    resolver *ServiceResolver
    health   map[string]bool
}

func (har *HealthAwareResolver) Watch(serviceName string) {
    ticker := time.NewTicker(10 * time.Second)
    
    for range ticker.C {
        entries, _, _ := har.resolver.consul.Health().Service(serviceName, "", true, nil)
        
        har.mu.Lock()
        for _, entry := range entries {
            har.health[entry.Service.ID] = true
        }
        har.mu.Unlock()
    }
}
```
