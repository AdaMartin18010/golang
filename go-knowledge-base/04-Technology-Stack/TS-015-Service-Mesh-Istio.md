# TS-015: Service Mesh 与 Istio (Service Mesh & Istio)

> **维度**: Technology Stack
> **级别**: S (17+ KB)
> **标签**: #service-mesh #istio #envoy #sidecar #microservices
> **权威来源**: [Istio Documentation](https://istio.io/latest/docs/), [Service Mesh Patterns](https://www.oreilly.com/library/view/service-mesh-patterns/9781492086449/)
> **版本**: Istio 1.25+

---

## Service Mesh 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Service Mesh Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  传统微服务 (无 Service Mesh):                                               │
│  ┌─────────┐      HTTP/TLS/mTLS      ┌─────────┐                            │
│  │ Service │ ─────────────────────── │ Service │                            │
│  │    A    │    (应用层处理)          │    B    │                            │
│  └─────────┘                         └─────────┘                            │
│                                                                              │
│  Service Mesh 架构:                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Pod/Service A                               │    │
│  │  ┌─────────────┐    localhost:15001   ┌─────────────┐              │    │
│  │  │ Application │◄────────────────────►│    Envoy    │◄────┐       │    │
│  │  │   (App)     │    (iptables 拦截)   │   (Sidecar) │     │       │    │
│  │  └─────────────┘                      └──────┬──────┘     │       │    │
│  │                                              │            │       │    │
│  └──────────────────────────────────────────────┼────────────┼───────┘    │
│                                                 │            │             │
│                         mTLS + Telemetry       │            │ mTLS        │
│                                                 │            │             │
│  ┌──────────────────────────────────────────────┼────────────┼───────┐    │
│  │                         Pod/Service B        │            │       │    │
│  │  ┌─────────────┐                      ┌──────┴──────┐     │       │    │
│  │  │ Application │◄────────────────────►│    Envoy    │◄────┘       │    │
│  │  │   (App)     │                      │   (Sidecar) │              │    │
│  │  └─────────────┘                      └─────────────┘              │    │
│  └────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  Control Plane (Istiod):                                                     │
│  - xDS API: 下发配置到 Envoy                                                  │
│  - Certificate Management: 自动 mTLS 证书                                     │
│  - Traffic Management: 路由、负载均衡                                          │
│  - Policy: 访问控制、限流                                                      │
│  - Telemetry: 指标、日志、追踪                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Istio 组件

### 控制平面 (Istiod)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Istio Control Plane                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Istiod                                      │    │
│  │                                                                      │    │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐           │    │
│  │  │   Pilot       │  │  Citadel      │  │  Galley       │           │    │
│  │  │               │  │               │  │               │           │    │
│  │  │ - xDS Server  │  │ - CA          │  │ - Config      │           │    │
│  │  │ - Service Dic │  │ - Cert Mgmt   │  │   Validation  │           │    │
│  │  │ - Routing     │  │ - mTLS        │  │ - MCP Server  │           │    │
│  │  └───────┬───────┘  └───────┬───────┘  └───────────────┘           │    │
│  │          │                  │                                       │    │
│  │          │ xDS API          │ Certificates                          │    │
│  │          ▼                  ▼                                       │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │                      Envoy Sidecars                         │    │    │
│  │  │  - LDS: Listener Discovery Service                          │    │    │
│  │  │  - RDS: Route Discovery Service                             │    │    │
│  │  │  - CDS: Cluster Discovery Service                           │    │    │
│  │  │  - EDS: Endpoint Discovery Service                          │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  xDS 协议:                                                                    │
│  - Envoy 通过 gRPC 流订阅配置                                                │
│  - 增量更新 (Delta xDS)                                                      │
│  - 配置热更新，无需重启                                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 数据平面 (Envoy)

```yaml
# Envoy 配置示例 (由 Istio 生成)
static_resources:
  listeners:
    - name: virtual_inbound
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 15006
      filter_chains:
        - filters:
            - name: envoy.filters.network.tcp_proxy
              typed_config:
                '@type': type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy
                cluster: inbound|8080||
                stat_prefix: inbound|8080||

  clusters:
    - name: outbound|443||api.example.com
      connect_timeout: 10s
      type: STRICT_DNS
      lb_policy: ROUND_ROBIN
      load_assignment:
        cluster_name: outbound|443||api.example.com
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: api.example.com
                      port_value: 443
      transport_socket:
        name: envoy.transport_sockets.tls
        typed_config:
          '@type': type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
          common_tls_context:
            tls_certificate_sds_secret_configs:
              - name: default
                sds_config:
                  path: /etc/istio/proxy/SDS
```

---

## 流量管理

### VirtualService & DestinationRule

```yaml
# VirtualService: 定义路由规则
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: reviews-route
spec:
  hosts:
    - reviews
  http:
    - match:
        - headers:
            end-user:
              exact: jason
      route:
        - destination:
            host: reviews
            subset: v2
    - route:
        - destination:
            host: reviews
            subset: v1
          weight: 90
        - destination:
            host: reviews
            subset: v2
          weight: 10

---
# DestinationRule: 定义服务子集和策略
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: reviews-destination
spec:
  host: reviews
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 50
        maxRequestsPerConnection: 2
    outlierDetection:
      consecutiveErrors: 5
      interval: 30s
      baseEjectionTime: 30s
  subsets:
    - name: v1
      labels:
        version: v1
    - name: v2
      labels:
        version: v2
```

### 熔断器与重试

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: service-retry
spec:
  hosts:
    - my-service
  http:
    - route:
        - destination:
            host: my-service
      retries:
        attempts: 3           # 重试次数
        perTryTimeout: 2s     # 单次超时
        retryOn: 5xx,gateway-error,connect-failure,refused-stream
      timeout: 10s            # 总超时
      fault:
        delay:
          percentage:
            value: 0.1        # 0.1% 延迟注入 (测试用)
          fixedDelay: 5s
```

---

## 安全 (mTLS)

### 自动 mTLS

```yaml
# 全局启用严格 mTLS
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: istio-system
spec:
  mtls:
    mode: STRICT  # 强制 mTLS

---
# 授权策略
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: service-access
  namespace: default
spec:
  selector:
    matchLabels:
      app: my-service
  action: ALLOW
  rules:
    - from:
        - source:
            principals: ["cluster.local/ns/default/sa/client-service"]
      to:
        - operation:
            methods: ["GET"]
            paths: ["/api/*"]
      when:
        - key: request.headers[x-api-key]
          values: ["valid-key"]
```

---

## 可观测性

### Kiali 服务拓扑

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Kiali Service Graph                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│     ┌─────────────┐                                                         │
│     │   Ingress   │                                                         │
│     └──────┬──────┘                                                         │
│            │ 100%                                                           │
│            ▼                                                                │
│     ┌─────────────┐    90%     ┌─────────────┐                             │
│     │  Frontend   │───────────►│  Reviews-v1 │                             │
│     └──────┬──────┘            └─────────────┘                             │
│            │                                                                 │
│            │ 10%                 ┌─────────────┐                            │
│            └────────────────────►│  Reviews-v2 │                            │
│                                  └──────┬──────┘                            │
│                                         │                                   │
│                                         ▼                                   │
│                                  ┌─────────────┐                            │
│                                  │  Ratings    │                            │
│                                  └─────────────┘                            │
│                                                                              │
│  颜色含义:                                                                   │
│  - 绿色: 健康                                                                │
│  - 红色: 错误率高                                                            │
│  - 黄色: 警告                                                                │
│  - 灰色: 无流量                                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 遥测配置

```yaml
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: default-metrics
spec:
  metrics:
    - providers:
        - name: prometheus
      overrides:
        - match:
            metric: REQUEST_COUNT
          tagOverrides:
            destination_service:
              value: "unknown"
  accessLogging:
    - providers:
        - name: envoy
      filter:
        expression: "response.code >= 500"
  tracing:
    - providers:
        - name: zipkin
      randomSamplingPercentage: 10.0
```

---

## 生产建议

| 场景 | 建议 |
|------|------|
| 渐进式采用 | 按命名空间启用，先观察模式 |
| 资源限制 | Sidecar 设置 CPU/Memory limits |
| 启动顺序 | 确保 Sidecar 先于应用就绪 |
| 调试 | 使用 `istioctl proxy-config` 查看配置 |
| 升级 | 使用金丝雀升级 (revision-based) |

---

## 参考文献

1. [Istio Documentation](https://istio.io/latest/docs/)
2. [Envoy Proxy Documentation](https://www.envoyproxy.io/docs)
3. [Service Mesh Patterns](https://www.oreilly.com/library/view/service-mesh-patterns/9781492086449/)
