# Gateway API实战指南

> **难度**: ⭐⭐⭐⭐
> **标签**: #GatewayAPI #Kubernetes #服务网格 #流量管理

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

## 📋 目录

- [Gateway API实战指南](#gateway-api实战指南)
  - [📋 目录](#-目录)
  - [1. Gateway API概述](#1-gateway-api概述)
    - [1.1 什么是Gateway API](#11-什么是gateway-api)
    - [1.2 核心资源](#12-核心资源)
    - [1.3 与Ingress对比](#13-与ingress对比)
  - [2. Gateway API安装](#2-gateway-api安装)
    - [2.1 安装CRD](#21-安装crd)
    - [2.2 安装Gateway控制器](#22-安装gateway控制器)
    - [2.3 验证安装](#23-验证安装)
  - [3. Gateway资源](#3-gateway资源)
    - [3.1 Gateway配置](#31-gateway配置)
    - [3.2 GatewayClass](#32-gatewayclass)
    - [3.3 监听器配置](#33-监听器配置)
  - [4. HTTPRoute配置](#4-httproute配置)
    - [4.1 基础路由](#41-基础路由)
    - [4.2 高级匹配](#42-高级匹配)
    - [4.3 流量分割](#43-流量分割)
  - [5. TLS与证书管理](#5-tls与证书管理)
    - [5.1 TLS终止](#51-tls终止)
    - [5.2 TLS透传](#52-tls透传)
    - [5.3 证书轮换](#53-证书轮换)
  - [6. 流量管理](#6-流量管理)
    - [6.1 流量分割与金丝雀](#61-流量分割与金丝雀)
    - [6.2 请求重定向](#62-请求重定向)
    - [6.3 请求镜像](#63-请求镜像)
  - [7. Go客户端实现](#7-go客户端实现)
    - [7.1 Gateway API客户端](#71-gateway-api客户端)
    - [7.2 动态路由管理](#72-动态路由管理)
  - [8. 实战案例](#8-实战案例)
    - [8.1 多租户API网关](#81-多租户api网关)
  - [9. 参考资源](#9-参考资源)
    - [官方文档](#官方文档)
    - [Go库](#go库)
    - [最佳实践](#最佳实践)

## 1. Gateway API概述

### 1.1 什么是Gateway API

**Gateway API** 是Kubernetes的下一代入口API，提供了更强大、更灵活的流量管理能力。

**核心优势**:

- ✅ **表达力强**：支持复杂的路由规则
- ✅ **角色导向**：基础设施管理员、集群运维、应用开发者职责分离
- ✅ **可扩展**：通过CRD扩展功能
- ✅ **可移植**：跨实现的一致性API
- ✅ **类型化**：强类型的Kubernetes资源

**应用场景**:

- API网关
- 服务网格
- 流量管理
- 金丝雀发布
- A/B测试

### 1.2 核心资源

**资源层次结构**:

```
GatewayClass (集群级别)
    ↓
Gateway (命名空间级别)
    ↓
Route (HTTPRoute, TCPRoute, etc.)
    ↓
Service (后端服务)
```

**核心资源说明**:

| 资源 | 角色 | 职责 |
|------|------|------|
| **GatewayClass** | 基础设施管理员 | 定义Gateway类型和配置 |
| **Gateway** | 集群运维 | 配置负载均衡器和监听器 |
| **Route** | 应用开发者 | 配置路由规则 |

### 1.3 与Ingress对比

**Ingress vs Gateway API**:

| 特性 | Ingress | Gateway API |
|------|---------|-------------|
| **表达力** | 基础 | 强大 |
| **协议支持** | HTTP/HTTPS | HTTP, HTTPS, TCP, UDP, gRPC |
| **流量分割** | ❌ 需要注解 | ✅ 原生支持 |
| **TLS配置** | 基础 | 高级 |
| **角色分离** | ❌ | ✅ |
| **可扩展性** | 通过注解 | 通过标准字段 |
| **成熟度** | GA | v1.0 GA |

---

## 2. Gateway API安装

### 2.1 安装CRD

**安装Gateway API CRDs**:

```bash
# 安装标准版本
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yaml

# 或安装实验性版本（包含更多特性）
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/experimental-install.yaml
```

**验证CRD安装**:

```bash
kubectl get crd | grep gateway
# 输出：
# gatewayclasses.gateway.networking.k8s.io
# gateways.gateway.networking.k8s.io
# httproutes.gateway.networking.k8s.io
# referencegrants.gateway.networking.k8s.io
```

### 2.2 安装Gateway控制器

**选择实现**:

Gateway API有多个实现可选：

1. **Istio** - 服务网格实现
2. **Envoy Gateway** - 轻量级网关
3. **Kong** - API网关
4. **Traefik** - 云原生边缘路由器
5. **Nginx Gateway Fabric** - Nginx实现

**安装Envoy Gateway示例**:

```bash
helm install eg oci://docker.io/envoyproxy/gateway-helm \
  --version v0.6.0 \
  -n envoy-gateway-system \
  --create-namespace
```

### 2.3 验证安装

**检查控制器状态**:

```bash
kubectl get pods -n envoy-gateway-system
kubectl get gatewayclass
```

---

## 3. Gateway资源

### 3.1 Gateway配置

**基础Gateway定义**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: my-gateway
  namespace: default
spec:
  gatewayClassName: envoy-gateway
  listeners:
  - name: http
    protocol: HTTP
    port: 80
    allowedRoutes:
      namespaces:
        from: All
  - name: https
    protocol: HTTPS
    port: 443
    tls:
      mode: Terminate
      certificateRefs:
      - name: my-tls-cert
    allowedRoutes:
      namespaces:
        from: All
```

### 3.2 GatewayClass

**GatewayClass定义**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: envoy-gateway
spec:
  controllerName: gateway.envoyproxy.io/gatewayclass-controller
  parametersRef:
    group: gateway.envoyproxy.io
    kind: EnvoyProxy
    name: custom-proxy-config
    namespace: envoy-gateway-system
```

### 3.3 监听器配置

**多协议监听器**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: multi-protocol-gateway
spec:
  gatewayClassName: envoy-gateway
  listeners:
  # HTTP监听器
  - name: http
    protocol: HTTP
    port: 80
    hostname: "*.example.com"

  # HTTPS监听器
  - name: https
    protocol: HTTPS
    port: 443
    hostname: "*.example.com"
    tls:
      mode: Terminate
      certificateRefs:
      - name: wildcard-cert

  # TCP监听器
  - name: tcp
    protocol: TCP
    port: 3306
```

---

## 4. HTTPRoute配置

### 4.1 基础路由

**简单HTTP路由**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: my-app-route
  namespace: default
spec:
  parentRefs:
  - name: my-gateway
    namespace: default

  hostnames:
  - "api.example.com"

  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api/v1
    backendRefs:
    - name: api-service
      port: 8080
```

### 4.2 高级匹配

**多条件匹配**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: advanced-route
spec:
  parentRefs:
  - name: my-gateway

  rules:
  # 匹配路径 + Header
  - matches:
    - path:
        type: PathPrefix
        value: /api/v2
      headers:
      - name: X-API-Version
        value: "2.0"
    backendRefs:
    - name: api-v2-service
      port: 8080

  # 匹配查询参数
  - matches:
    - path:
        type: PathPrefix
        value: /search
      queryParams:
      - name: version
        value: beta
    backendRefs:
    - name: search-beta-service
      port: 8080

  # 匹配HTTP方法
  - matches:
    - path:
        type: Exact
        value: /users
      method: POST
    backendRefs:
    - name: user-create-service
      port: 8080
```

### 4.3 流量分割

**金丝雀发布配置**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: canary-route
spec:
  parentRefs:
  - name: my-gateway

  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /app

    backendRefs:
    # 90% 流量到稳定版本
    - name: app-stable
      port: 8080
      weight: 90

    # 10% 流量到金丝雀版本
    - name: app-canary
      port: 8080
      weight: 10
```

---

## 5. TLS与证书管理

### 5.1 TLS终止

**HTTPS配置**:

```yaml
# 1. 创建TLS Secret
apiVersion: v1
kind: Secret
metadata:
  name: example-com-tls
  namespace: default
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-cert>
  tls.key: <base64-encoded-key>

---
# 2. Gateway配置TLS
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: tls-gateway
spec:
  gatewayClassName: envoy-gateway
  listeners:
  - name: https
    protocol: HTTPS
    port: 443
    hostname: "example.com"
    tls:
      mode: Terminate
      certificateRefs:
      - kind: Secret
        name: example-com-tls

---
# 3. HTTPRoute
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: secure-route
spec:
  parentRefs:
  - name: tls-gateway
  hostnames:
  - "example.com"
  rules:
  - backendRefs:
    - name: backend-service
      port: 8080
```

### 5.2 TLS透传

**TLS Passthrough配置**:

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: TLSRoute
metadata:
  name: tls-passthrough
spec:
  parentRefs:
  - name: my-gateway

  hostnames:
  - "secure.example.com"

  rules:
  - backendRefs:
    - name: secure-backend
      port: 443
```

### 5.3 证书轮换

**使用cert-manager自动轮换**:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com-cert
  namespace: default
spec:
  secretName: example-com-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
  - example.com
  - "*.example.com"
```

---

## 6. 流量管理

### 6.1 流量分割与金丝雀

**渐进式金丝雀发布**:

```yaml
# Phase 1: 5% 流量到金丝雀
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: gradual-rollout-phase1
spec:
  parentRefs:
  - name: my-gateway
  rules:
  - backendRefs:
    - name: app-stable
      port: 8080
      weight: 95
    - name: app-canary
      port: 8080
      weight: 5

---
# Phase 2: 25% 流量到金丝雀
# （更新weight即可）
```

### 6.2 请求重定向

**HTTP到HTTPS重定向**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: http-to-https
spec:
  parentRefs:
  - name: my-gateway

  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    filters:
    - type: RequestRedirect
      requestRedirect:
        scheme: https
        statusCode: 301
```

**路径重定向**:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: path-redirect
spec:
  parentRefs:
  - name: my-gateway

  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /old-api
    filters:
    - type: RequestRedirect
      requestRedirect:
        path:
          type: ReplaceFullPath
          replaceFullPath: /new-api
        statusCode: 301
```

### 6.3 请求镜像

**流量镜像到测试环境**:

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: HTTPRoute
metadata:
  name: traffic-mirror
spec:
  parentRefs:
  - name: my-gateway

  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api

    # 主要后端
    backendRefs:
    - name: production-service
      port: 8080

    # 镜像配置
    filters:
    - type: RequestMirror
      requestMirror:
        backendRef:
          name: test-service
          port: 8080
```

---

## 7. Go客户端实现

### 7.1 Gateway API客户端

**初始化客户端**:

```go
package main

import (
    "context"
    "fmt"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/clientcmd"
    gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
    gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

// GatewayAPIClient Gateway API客户端
type GatewayAPIClient struct {
    client *gatewayclient.Clientset
}

func NewGatewayAPIClient(kubeconfig string) (*GatewayAPIClient, error) {
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        return nil, err
    }

    client, err := gatewayclient.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    return &GatewayAPIClient{client: client}, nil
}

// CreateGateway 创建Gateway
func (c *GatewayAPIClient) CreateGateway(ctx context.Context, namespace string, gateway *gatewayv1.Gateway) error {
    _, err := c.client.GatewayV1().Gateways(namespace).Create(ctx, gateway, metav1.CreateOptions{})
    return err
}

// ListGateways 列出所有Gateway
func (c *GatewayAPIClient) ListGateways(ctx context.Context, namespace string) (*gatewayv1.GatewayList, error) {
    return c.client.GatewayV1().Gateways(namespace).List(ctx, metav1.ListOptions{})
}

// CreateHTTPRoute 创建HTTPRoute
func (c *GatewayAPIClient) CreateHTTPRoute(ctx context.Context, namespace string, route *gatewayv1.HTTPRoute) error {
    _, err := c.client.GatewayV1().HTTPRoutes(namespace).Create(ctx, route, metav1.CreateOptions{})
    return err
}

// 使用示例
func main() {
    client, err := NewGatewayAPIClient(clientcmd.RecommendedHomeFile)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // 列出Gateway
    gateways, err := client.ListGateways(ctx, "default")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d gateways\n", len(gateways.Items))
    for _, gw := range gateways.Items {
        fmt.Printf("- %s (Class: %s)\n", gw.Name, gw.Spec.GatewayClassName)
    }
}
```

### 7.2 动态路由管理

**路由管理器**:

```go
package gateway

import (
    "context"
    "fmt"

    gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RouteManager 路由管理器
type RouteManager struct {
    client *GatewayAPIClient
}

func NewRouteManager(client *GatewayAPIClient) *RouteManager {
    return &RouteManager{client: client}
}

// UpdateCanaryWeight 更新金丝雀权重
func (rm *RouteManager) UpdateCanaryWeight(ctx context.Context, namespace, routeName string, stableWeight, canaryWeight int32) error {
    // 获取现有HTTPRoute
    route, err := rm.client.client.GatewayV1().HTTPRoutes(namespace).Get(ctx, routeName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    // 更新权重
    if len(route.Spec.Rules) > 0 && len(route.Spec.Rules[0].BackendRefs) >= 2 {
        route.Spec.Rules[0].BackendRefs[0].Weight = &stableWeight
        route.Spec.Rules[0].BackendRefs[1].Weight = &canaryWeight
    }

    // 更新HTTPRoute
    _, err = rm.client.client.GatewayV1().HTTPRoutes(namespace).Update(ctx, route, metav1.UpdateOptions{})
    if err != nil {
        return fmt.Errorf("update route: %w", err)
    }

    fmt.Printf("Updated canary weight: stable=%d%%, canary=%d%%\n", stableWeight, canaryWeight)
    return nil
}

// CreateCanaryRoute 创建金丝雀路由
func (rm *RouteManager) CreateCanaryRoute(ctx context.Context, namespace, routeName, gatewayName string,
    stableService, canaryService string, stableWeight, canaryWeight int32) error {

    route := &gatewayv1.HTTPRoute{
        ObjectMeta: metav1.ObjectMeta{
            Name:      routeName,
            Namespace: namespace,
        },
        Spec: gatewayv1.HTTPRouteSpec{
            ParentRefs: []gatewayv1.ParentReference{
                {
                    Name: gatewayv1.ObjectName(gatewayName),
                },
            },
            Rules: []gatewayv1.HTTPRouteRule{
                {
                    BackendRefs: []gatewayv1.HTTPBackendRef{
                        {
                            BackendRef: gatewayv1.BackendRef{
                                BackendObjectReference: gatewayv1.BackendObjectReference{
                                    Name: gatewayv1.ObjectName(stableService),
                                    Port: portPtr(8080),
                                },
                                Weight: &stableWeight,
                            },
                        },
                        {
                            BackendRef: gatewayv1.BackendRef{
                                BackendObjectReference: gatewayv1.BackendObjectReference{
                                    Name: gatewayv1.ObjectName(canaryService),
                                    Port: portPtr(8080),
                                },
                                Weight: &canaryWeight,
                            },
                        },
                    },
                },
            },
        },
    }

    return rm.client.CreateHTTPRoute(ctx, namespace, route)
}

func portPtr(port int32) *gatewayv1.PortNumber {
    p := gatewayv1.PortNumber(port)
    return &p
}

// 渐进式金丝雀发布
func (rm *RouteManager) ProgressiveCanary(ctx context.Context, namespace, routeName string) error {
    stages := []struct {
        stable int32
        canary int32
    }{
        {95, 5},   // 5%
        {90, 10},  // 10%
        {75, 25},  // 25%
        {50, 50},  // 50%
        {25, 75},  // 75%
        {0, 100},  // 100%
    }

    for i, stage := range stages {
        fmt.Printf("Stage %d: %d%% canary\n", i+1, stage.canary)

        if err := rm.UpdateCanaryWeight(ctx, namespace, routeName, stage.stable, stage.canary); err != nil {
            return err
        }

        // 等待验证（实际应用中应该检查指标）
        // time.Sleep(5 * time.Minute)
    }

    fmt.Println("Canary deployment completed!")
    return nil
}
```

---

## 8. 实战案例

### 8.1 多租户API网关

**租户隔离配置**:

```yaml
# 租户A的Gateway
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: tenant-a-gateway
  namespace: tenant-a
spec:
  gatewayClassName: envoy-gateway
  listeners:
  - name: https
    protocol: HTTPS
    port: 443
    hostname: "api-a.example.com"
    tls:
      mode: Terminate
      certificateRefs:
      - name: tenant-a-cert

---
# 租户A的路由
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: tenant-a-routes
  namespace: tenant-a
spec:
  parentRefs:
  - name: tenant-a-gateway
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: tenant-a-service
      port: 8080
```

---

## 9. 参考资源

### 官方文档

- [Gateway API Documentation](https://gateway-api.sigs.k8s.io/)
- [Gateway API Spec](https://gateway-api.sigs.k8s.io/references/spec/)
- [Implementations](https://gateway-api.sigs.k8s.io/implementations/)

### Go库

- [gateway-api](https://github.com/kubernetes-sigs/gateway-api)
- [client-go](https://github.com/kubernetes/client-go)

### 最佳实践

- [Gateway API User Guides](https://gateway-api.sigs.k8s.io/guides/)
- [Gateway API Patterns](https://gateway-api.sigs.k8s.io/guides/http-routing/)

---

**文档维护者**: Go Documentation Team
**最后更新**: 2025-10-29
**文档状态**: ✅ 完成
**适用版本**: Go 1.21+ | Gateway API v1.0+

**贡献者**: 欢迎提交Issue和PR改进本文档
