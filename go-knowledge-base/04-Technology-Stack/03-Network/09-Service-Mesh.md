# 服务网格 (Service Mesh)

> **分类**: 开源技术堆栈
> **标签**: #servicemesh #istio #envoy

---

## Istio 集成

### 自动注入

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-app
  labels:
    istio-injection: enabled
```

### Go 应用适配

```go
// 无需修改代码，Istio 自动处理:
// - mTLS
// - 流量管理
// - 可观测性

func main() {
    // 普通 HTTP 服务即可
    r := gin.Default()
    r.GET("/api/data", handler)
    r.Run(":8080")
}
```

---

## Envoy Proxy

### 数据面代理

```go
// 使用 Envoy 的 xDS API
import (
    cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
    listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
)

func CreateListener() *listener.Listener {
    return &listener.Listener{
        Name: "listener_0",
        Address: &core.Address{
            Address: &core.Address_SocketAddress{
                SocketAddress: &core.SocketAddress{
                    Address: "0.0.0.0",
                    PortSpecifier: &core.SocketAddress_PortValue{
                        PortValue: 8080,
                    },
                },
            },
        },
    }
}
```

---

## 流量管理

### 金丝雀发布

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: my-app
spec:
  hosts:
  - my-app
  http:
  - match:
    - headers:
        canary:
          exact: "true"
    route:
    - destination:
        host: my-app
        subset: v2
      weight: 100
  - route:
    - destination:
        host: my-app
        subset: v1
      weight: 90
    - destination:
        host: my-app
        subset: v2
      weight: 10
```

---

## Go 客户端 Sidecar

```go
// 通过 localhost 访问 sidecar
conn, err := grpc.Dial("localhost:15003",  // Envoy 监听端口
    grpc.WithInsecure(),
)
```
