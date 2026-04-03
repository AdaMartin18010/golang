# EC-090-API-Gateway-Design

> **Dimension**: 03-Engineering-CloudNative
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: 2026 (Kong, Envoy, AWS API Gateway, Traefik)
> **Size**: >20KB

---

## 1. API Gateway概览

### 1.1 核心功能

```
┌─────────────────────────────────────────┐
│           API Gateway                   │
├─────────────────────────────────────────┤
│                                         │
│  ┌──────────┐  ┌──────────┐  ┌────────┐│
│  │ 认证授权 │  │ 流量控制 │  │ 路由  ││
│  │ Auth     │  │ Rate     │  │ Route ││
│  │ OAuth2   │  │ Limit    │  │ Match ││
│  │ JWT      │  │ Throttle │  │       ││
│  └──────────┘  └──────────┘  └────────┘│
│                                         │
│  ┌──────────┐  ┌──────────┐  ┌────────┐│
│  │ 协议转换 │  │ 缓存加速 │  │ 日志  ││
│  │ HTTP/gRPC│  │ Redis    │  │ Trace ││
│  │ REST/SOAP│  │ CDN      │  │ Metric││
│  └──────────┘  └──────────┘  └────────┘│
│                                         │
│  ┌──────────┐  ┌──────────┐            │
│  │ 负载均衡 │  │ 安全加固 │            │
│  │ LB       │  │ WAF      │            │
│  │ Health   │  │ DDoS     │            │
│  └──────────┘  └──────────┘            │
│                                         │
└─────────────────────────────────────────┘
```

### 1.2 网关选型对比

| 网关 | 特点 | 适用场景 |
|------|------|---------|
| Kong | 插件丰富，企业级 | 大型微服务 |
| Envoy | 云原生，高性能 | Service Mesh |
| Traefik | 自动发现，易配置 | Kubernetes |
| Nginx | 成熟稳定 | 传统架构 |
| AWS APIG | 托管服务 | AWS生态 |
| APISIX | 开源高性能 | 云原生 |

---

## 2. Kong Gateway

### 2.1 架构

```
┌─────────────────────────────────────────┐
│           Kong Architecture             │
├─────────────────────────────────────────┤
│                                         │
│  Nginx/OpenResty                        │
│       │                                 │
│  Kong (Lua)                             │
│    ├── Router                           │
│    ├── Plugin System                    │
│    └── Admin API                        │
│       │                                 │
│  PostgreSQL/Cassandra (配置存储)        │
│                                         │
└─────────────────────────────────────────┘
```

### 2.2 核心配置

```yaml
# kong.yml 声明式配置
_format_version: "3.0"

services:
  - name: user-service
    url: http://user-service:8080
    routes:
      - name: user-routes
        paths:
          - /api/users
    plugins:
      - name: rate-limiting
        config:
          minute: 100
          policy: redis
          redis_host: redis
      - name: jwt
        config:
          uri_param_names: []
          cookie_names: []
          key_claim_name: iss
          secret_is_base64: false
          claims_to_verify:
            - exp
      - name: prometheus

  - name: order-service
    url: http://order-service:8080
    routes:
      - name: order-routes
        paths:
          - /api/orders
    plugins:
      - name: oauth2
        config:
          scopes:
            - email
            - profile
            - openid
          mandatory_scope: true
          enable_authorization_code: true
      - name: request-transformer
        config:
          add:
            headers:
              - X-Request-ID:$(request_id)
```

### 2.3 自定义插件

```lua
-- custom-auth.lua
local plugin = {
    PRIORITY = 1000,
    VERSION = "1.0.0",
}

function plugin:access(plugin_conf)
    local token = kong.request.get_header("Authorization")

    if not token then
        return kong.response.exit(401, { message = "Unauthorized" })
    end

    -- 调用认证服务
    local http = require "resty.http"
    local httpc = http.new()

    local res, err = httpc:request_uri("http://auth-service:8080/verify", {
        method = "POST",
        headers = {
            ["Authorization"] = token,
        },
    })

    if not res or res.status ~= 200 then
        return kong.response.exit(401, { message = "Invalid token" })
    end

    -- 解析用户信息
    local cjson = require "cjson"
    local user = cjson.decode(res.body)

    -- 设置header传递给后端
    kong.service.request.set_header("X-User-ID", user.id)
    kong.service.request.set_header("X-User-Role", user.role)
end

return plugin
```

---

## 3. Envoy Proxy

### 3.1 架构

```
┌─────────────────────────────────────────┐
│           Envoy Architecture            │
├─────────────────────────────────────────┤
│                                         │
│  Downstream (Client)                    │
│       │                                 │
│  ┌────┴────┐                            │
│  │ Listener│                            │
│  └────┬────┘                            │
│       │                                 │
│  ┌────┴────┐      ┌───────────────┐    │
│  │  Router │─────►│ Cluster Manager│    │
│  └────┬────┘      └───────┬───────┘    │
│       │                   │            │
│  Filter Chain         Clusters         │
│  - HTTP Connection    - Load Balance   │
│  - Rate Limit         - Health Check   │
│  - Auth               - Circuit Breaker│
│       │                   │            │
│  ┌────┴────┐      ┌───────┴───────┐    │
│  │ Upstream│      │   Upstream    │    │
│  └─────────┘      └───────────────┘    │
│                                         │
└─────────────────────────────────────────┘
```

### 3.2 配置示例

```yaml
# envoy.yaml
static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address:
          address: 0.0.0.0
          port_value: 8080
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                stat_prefix: ingress_http
                route_config:
                  name: local_route
                  virtual_hosts:
                    - name: backend
                      domains: ["*"]
                      routes:
                        - match:
                            prefix: "/api/users"
                          route:
                            cluster: user_service
                            timeout: 30s
                            retry_policy:
                              retry_on: "5xx,connect-failure"
                              num_retries: 3
                              per_try_timeout: 10s
                        - match:
                            prefix: "/api/orders"
                          route:
                            cluster: order_service
                            rate_limits:
                              - actions:
                                  - request_headers:
                                      header_name: x-user-id
                                      descriptor_key: user_id
                http_filters:
                  - name: envoy.filters.http.ext_authz
                    typed_config:
                      "@type": type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz
                      grpc_service:
                        envoy_grpc:
                          cluster_name: ext_authz
                        timeout: 5s
                      include_peer_certificate: true
                  - name: envoy.filters.http.local_ratelimit
                    typed_config:
                      "@type": type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit
                      stat_prefix: http_local_rate_limiter
                      token_bucket:
                        max_tokens: 1000
                        tokens_per_fill: 100
                        fill_interval: 1s
                      filter_enabled:
                        runtime_key: local_rate_limit_enabled
                        default_value:
                          numerator: 100
                          denominator: HUNDRED
                  - name: envoy.filters.http.router
                    typed_config:
                      "@type": type.googleapis.com/envoy.extensions.filters.http.router.v3.Router

  clusters:
    - name: user_service
      connect_timeout: 5s
      type: STRICT_DNS
      lb_policy: ROUND_ROBIN
      health_checks:
        - timeout: 5s
          interval: 10s
          unhealthy_threshold: 3
          healthy_threshold: 2
          http_health_check:
            path: "/health"
      load_assignment:
        cluster_name: user_service
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: user-service
                      port_value: 8080

    - name: order_service
      connect_timeout: 5s
      type: STRICT_DNS
      lb_policy: LEAST_REQUEST
      circuit_breakers:
        thresholds:
          - priority: DEFAULT
            max_connections: 1000
            max_pending_requests: 1000
            max_requests: 1000
            max_retries: 3
      load_assignment:
        cluster_name: order_service
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: order-service
                      port_value: 8080

    - name: ext_authz
      connect_timeout: 5s
      type: STRICT_DNS
      typed_extension_protocol_options:
        envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
          "@type": type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions
          upstream_protocol_options:
            explicit_http_config:
              http2_protocol_options: {}
      load_assignment:
        cluster_name: ext_authz
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: auth-service
                      port_value: 8080
```

---

## 4. Traefik

### 4.1 Kubernetes集成

```yaml
# IngressRoute
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: api-gateway
  namespace: default
spec:
  entryPoints:
    - web
    - websecure
  routes:
    - match: Host(`api.example.com`) && PathPrefix(`/users`)
      kind: Rule
      services:
        - name: user-service
          port: 8080
      middlewares:
        - name: rate-limit
        - name: auth-jwt

    - match: Host(`api.example.com`) && PathPrefix(`/orders`)
      kind: Rule
      services:
        - name: order-service
          port: 8080
          sticky:
            cookie:
              name: order_session
      middlewares:
        - name: circuit-breaker
        - name: retry

  tls:
    certResolver: letsencrypt
    options:
      name: tls-options

---
# Middleware - 限流
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: rate-limit
spec:
  rateLimit:
    average: 100
    burst: 200
    period: 1m
    sourceCriterion:
      requestHost: true

---
# Middleware - JWT认证
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: auth-jwt
spec:
  plugin:
    jwt:
      secret: ${JWT_SECRET}
      forwardHeaders:
        X-User-ID: sub
        X-User-Role: role

---
# Middleware - 熔断
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: circuit-breaker
spec:
  circuitBreaker:
    expression: LatencyAtQuantileMS(50.0) > 100
    fallbackDuration: 10s
    recoveryDuration: 30s
    responseCode: 503

---
# Middleware - 重试
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: retry
spec:
  retry:
    attempts: 3
    initialInterval: 100ms
```

---

## 5. 高级模式

### 5.1 GraphQL Federation

```go
// GraphQL Gateway
type Gateway struct {
    services []*Service
    schema   *graphql.Schema
}

type Service struct {
    Name   string
    URL    string
    Schema string
}

func (g *Gateway) Execute(ctx context.Context, query string) (*Response, error) {
    // 解析查询
    doc, err := parser.Parse(query)
    if err != nil {
        return nil, err
    }

    // 查询规划
    plan := g.planner.Plan(doc)

    // 并行执行子查询
    results := make(map[string]interface{})
    var wg sync.WaitGroup
    errChan := make(chan error, len(plan.Nodes))

    for _, node := range plan.Nodes {
        wg.Add(1)
        go func(n PlanNode) {
            defer wg.Done()

            result, err := g.executeNode(ctx, n)
            if err != nil {
                errChan <- err
                return
            }

            g.mergeResults(results, n.Path, result)
        }(node)
    }

    wg.Wait()
    close(errChan)

    if len(errChan) > 0 {
        return nil, <-errChan
    }

    return &Response{Data: results}, nil
}
```

### 5.2 BFF (Backend for Frontend)

```go
// Mobile BFF
type MobileBFF struct {
    userClient    pb.UserServiceClient
    orderClient   pb.OrderServiceClient
    productClient pb.ProductServiceClient
}

func (b *MobileBFF) GetHomePage(ctx context.Context, userID string) (*HomePage, error) {
    var wg sync.WaitGroup
    var user *User
    var orders []*Order
    var recommendations []*Product

    // 并行获取数据
    wg.Add(3)

    go func() {
        defer wg.Done()
        user, _ = b.userClient.GetUser(ctx, &pb.GetUserRequest{Id: userID})
    }()

    go func() {
        defer wg.Done()
        resp, _ := b.orderClient.ListOrders(ctx, &pb.ListOrdersRequest{UserId: userID, Limit: 5})
        orders = resp.Orders
    }()

    go func() {
        defer wg.Done()
        resp, _ := b.productClient.GetRecommendations(ctx, &pb.RecommendationRequest{UserId: userID})
        recommendations = resp.Products
    }()

    wg.Wait()

    return &HomePage{
        User:            user,
        RecentOrders:    orders,
        Recommendations: recommendations,
    }, nil
}
```

---

## 6. 可观测性

### 6.1 指标收集

```yaml
# Prometheus指标
- name: http_requests_total
  help: Total HTTP requests
  type: counter
  labels: [method, route, status]

- name: http_request_duration_seconds
  help: HTTP request duration
  type: histogram
  labels: [method, route]
  buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]

- name: active_connections
  help: Active connections
  type: gauge
```

### 6.2 分布式追踪

```go
// OpenTelemetry集成
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, span := tracer.Start(r.Context(), "gateway_request",
            trace.WithAttributes(
                attribute.String("http.method", r.Method),
                attribute.String("http.url", r.URL.String()),
                attribute.String("http.target", r.URL.Path),
            ),
        )
        defer span.End()

        // 传播trace context到后端
        r = r.WithContext(ctx)
        propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

        next.ServeHTTP(w, r)
    })
}
```

---

## 7. 安全最佳实践

### 7.1 认证架构

```
Client ──► API Gateway ──► Auth Server ──► Backend Services
              │                  │
              │              (OAuth2/OIDC)
              │
         JWT验证
         Token刷新
         Scope检查
```

### 7.2 防护策略

| 威胁 | 防护措施 |
|------|---------|
| DDoS | 速率限制、WAF、CDN |
| Injection | 输入验证、参数化查询 |
| 信息泄露 | 错误脱敏、Header清理 |
| 中间人 | TLS 1.3、证书固定 |
| 重放攻击 | 请求签名、时间戳检查 |

---

## 8. 参考文献

1. Kong Documentation
2. Envoy Proxy Documentation
3. Traefik Documentation
4. "Designing Web APIs" - Brenda Jin
5. API Gateway Patterns

---

*Last Updated: 2026-04-03*
