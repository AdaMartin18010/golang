# 第十四章：微服务框架对比

> Go 1.26 时代的主流微服务框架全面对比

---

## 14.1 框架概览

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go 微服务框架生态图                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   基础设施层                                                                 │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│   │  Kubernetes │  │   Docker    │  │   Consul    │  │    etcd     │       │
│   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘       │
│                                                                             │
│   框架层                                                                     │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│   │   Encore    │  │   GoMicro   │  │    Kit      │  │    Kratos   │       │
│   │  (全功能)   │  │  (分布式)   │  │  (工具集)   │  │  (云原生)   │       │
│   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘       │
│                                                                             │
│   Web 层                                                                     │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│   │    Gin      │  │    Echo     │  │    Fiber    │  │     Chi     │       │
│   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 14.2 框架对比矩阵

| 特性 | Encore | GoMicro | Go kit | Kratos | Gin | Echo | Fiber |
|------|--------|---------|--------|--------|-----|------|-------|
| **定位** | 全栈框架 | 微服务框架 | 工具集 | 云原生框架 | Web框架 | Web框架 | Web框架 |
| **服务发现** | 内置 | 内置 | 需集成 | 内置 | 无 | 无 | 无 |
| **消息队列** | 内置 | 内置 | 需集成 | 内置 | 无 | 无 | 无 |
| **链路追踪** | 自动 | 需配置 | 需集成 | 内置 | 无 | 无 | 无 |
| **API文档** | 自动生成 | 需配置 | 需集成 | 自动生成 | 需配置 | 需配置 | 需配置 |
| **gRPC支持** | 是 | 是 | 是 | 是 | 需插件 | 需插件 | 需插件 |
| **学习曲线** | 低 | 中 | 高 | 中 | 低 | 低 | 低 |
| **GitHub Stars** | 5k+ | 22k+ | 27k+ | 24k+ | 81k+ | 31k+ | 35k+ |

---

## 14.3 Encore 框架

### 14.3.1 架构特点

```go
// Encore: Infrastructure from Code
// 代码即基础设施

package user

import (
    "context"
    "encore.dev/storage/sqldb"
)

// 数据库定义
var db = sqldb.NewDatabase("user", sqldb.DatabaseConfig{
    Migrations: "./migrations",
})

// API 定义 - 自动生成基础设施
//encore:api public path=/user/:id
type GetUserResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func GetUser(ctx context.Context, id int) (*GetUserResponse, error) {
    var user GetUserResponse
    err := db.QueryRow(ctx, `
        SELECT id, name FROM users WHERE id = $1
    `, id).Scan(&user.ID, &user.Name)
    return &user, err
}

// 服务间调用 - 类型安全
//encore:api private
type CreateUserParams struct {
    Name string `json:"name"`
    Email string `json:"email"`
}

func CreateUser(ctx context.Context, p *CreateUserParams) error {
    // 自动生成分布式追踪
    _, err := db.Exec(ctx, `
        INSERT INTO users (name, email) VALUES ($1, $2)
    `, p.Name, p.Email)
    return err
}

// Pub/Sub - 内置消息队列
import "encore.dev/pubsub"

type UserEvent struct {
    UserID int `json:"user_id"`
    Action string `json:"action"`
}

var UserEvents = pubsub.NewTopic[*UserEvent]("user-events", pubsub.TopicConfig{
    DeliveryGuarantee: pubsub.AtLeastOnce,
})

// 订阅者 - 自动部署
//encore:subscription service=notification
type NotificationService struct{}

func (s *NotificationService) SendWelcomeEmail(ctx context.Context, event *UserEvent) error {
    if event.Action == "created" {
        // 发送邮件
    }
    return nil
}
```

### 14.3.2 开发体验

```bash
# 一键本地开发环境
encore run

# 自动生成的架构图
encore graph

# 一键部署到云
encore deploy

# 自动生成的 API 文档
curl http://localhost:9400/
```

---

## 14.4 Go kit 工具集

### 14.4.1 分层架构

```go
// Go kit: 严格的分层架构

package main

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/transport/http"
)

// 1. Service 层 - 业务逻辑
type StringService interface {
    Uppercase(context.Context, string) (string, error)
    Count(context.Context, string) int
}

type stringService struct{}

func (s stringService) Uppercase(_ context.Context, str string) (string, error) {
    if str == "" {
        return "", ErrEmpty
    }
    return strings.ToUpper(str), nil
}

func (s stringService) Count(_ context.Context, str string) int {
    return len(str)
}

// 2. Endpoint 层 - 传输抽象
type uppercaseRequest struct {
    S string `json:"s"`
}

type uppercaseResponse struct {
    V   string `json:"v"`
    Err string `json:"err,omitempty"`
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(uppercaseRequest)
        v, err := svc.Uppercase(ctx, req.S)
        if err != nil {
            return uppercaseResponse{V: v, Err: err.Error()}, nil
        }
        return uppercaseResponse{V: v}, nil
    }
}

// 3. Transport 层 - HTTP
type httpError struct {
    Error string `json:"error"`
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req uppercaseRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return nil, err
    }
    return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response)
}

func main() {
    svc := stringService{}
    uppercaseHandler := http.NewServer(
        makeUppercaseEndpoint(svc),
        decodeUppercaseRequest,
        encodeResponse,
    )

    http.Handle("/uppercase", uppercaseHandler)
    http.ListenAndServe(":8080", nil)
}
```

### 14.4.2 中间件链

```go
// Go kit 中间件
func loggingMiddleware(logger log.Logger) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, request interface{}) (interface{}, error) {
            logger.Log("msg", "calling endpoint")
            defer logger.Log("msg", "called endpoint")
            return next(ctx, request)
        }
    }
}

func instrumentingMiddleware() endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, request interface{}) (interface{}, error) {
            // 指标收集
            return next(ctx, request)
        }
    }
}

// 应用中间件
endpoint := makeUppercaseEndpoint(svc)
endpoint = loggingMiddleware(logger)(endpoint)
endpoint = instrumentingMiddleware()(endpoint)
```

---

## 14.5 Kratos 框架

### 14.5.1 项目结构

```
myapp/
├── api/                    # API 定义 (Protobuf)
│   └── helloworld/
│       └── v1/
│           ├── helloworld.proto
│           └── helloworld.pb.go
├── cmd/                    # 入口
│   └── server/
│       └── main.go
├── internal/               # 内部实现
│   ├── biz/               # 业务逻辑
│   ├── data/              # 数据访问
│   ├── server/            # 服务注册
│   └── service/           # 服务实现
└── configs/               # 配置文件
```

### 14.5.2 代码示例

```go
// Kratos: 整洁架构

package service

import (
    "context"

    pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/middleware"
    "github.com/go-kratos/kratos/v2/middleware/recovery"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/go-kratos/kratos/v2/transport/http"
)

// GreeterService 实现
type GreeterService struct {
    pb.UnimplementedGreeterServer
    log *log.Helper
}

func NewGreeterService(logger log.Logger) *GreeterService {
    return &GreeterService{log: log.NewHelper(logger)}
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    s.log.WithContext(ctx).Infof("SayHello Received: %v", req.GetName())
    return &pb.HelloReply{Message: "Hello " + req.GetName()}, nil
}

// 服务启动
func main() {
    logger := log.DefaultLogger

    gs := grpc.NewServer(
        grpc.Address(":9000"),
        grpc.Middleware(
            recovery.Recovery(),
            logging.Server(logger),
        ),
    )

    hs := http.NewServer(
        http.Address(":8000"),
        http.Middleware(
            recovery.Recovery(),
            logging.Server(logger),
        ),
    )

    svc := NewGreeterService(logger)
    pb.RegisterGreeterServer(gs, svc)
    pb.RegisterGreeterHTTPServer(hs, svc)

    app := kratos.New(
        kratos.Name("helloworld"),
        kratos.Server(gs, hs),
    )

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## 14.6 框架选择决策树

```text
┌─────────────────────────────────────────────────────────────────────┐
│                    框架选择决策树                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  需要完整的微服务基础设施？                                          │
│       │                                                             │
│       ├── Yes ──▶ 需要自动部署和云原生支持？                         │
│       │              │                                               │
│       │              ├── Yes ──▶ Encore (最简单)                     │
│       │              │                                               │
│       │              └── No ──▶ 需要灵活的工具集？                   │
│       │                         │                                    │
│       │                         ├── Yes ──▶ Go kit                   │
│       │                         │                                    │
│       │                         └── No ──▶ Kratos                    │
│       │                                                              │
│       └── No ──▶ 只需要 Web 框架？                                   │
│                     │                                                │
│                     ├── Yes ──▶ 性能优先？                           │
│                     │              │                                 │
│                     │              ├── Yes ──▶ Fiber (基于 fasthttp) │
│                     │              │                                 │
│                     │              └── No ──▶ 社区支持优先？         │
│                     │                         │                      │
│                     │                         ├── Yes ──▶ Gin        │
│                     │                         │                      │
│                     │                         └── No ──▶ Echo        │
│                     │                                                │
│                     └── No ──▶ 标准库 net/http + 轻量路由           │
│                                  │                                   │
│                                  └── Chi                             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 14.7 性能基准测试

```text
┌─────────────────────────────────────────────────────────────────────────┐
│                    框架性能对比 (Requests/sec)                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Fiber      ████████████████████████████████████████████████  180,000  │
│  Gin        ████████████████████████████████████████          140,000  │
│  Echo       ██████████████████████████████████████            130,000  │
│  fasthttp   ███████████████████████████████████               120,000  │
│  Chi        ████████████████████████████████                  100,000  │
│  net/http   ████████████████████                               60,000  │
│                                                                         │
│  测试条件: CPU-bound, 并发 100, 无延迟                                  │
│  数据来源: github.com/smallnest/go-web-framework-benchmark             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 14.8 云原生集成

### 14.8.1 Kubernetes 部署

```yaml
# Kratos 部署示例
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-service
  template:
    metadata:
      labels:
        app: go-service
    spec:
      containers:
      - name: service
        image: myapp:latest
        ports:
        - containerPort: 8000
          name: http
        - containerPort: 9000
          name: grpc
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: go-service
spec:
  selector:
    app: go-service
  ports:
  - port: 80
    targetPort: 8000
    name: http
  - port: 9000
    targetPort: 9000
    name: grpc
```

---

## 14.9 小结

| 场景 | 推荐框架 | 理由 |
|------|----------|------|
| 快速原型/MVP | Encore | 基础设施自动生成 |
| 大型企业应用 | Kratos | 整洁架构，云原生支持 |
| 极致性能 | Fiber | 基于 fasthttp，最快 |
| 社区生态 | Gin | 最流行，插件丰富 |
| 学习/教育 | 标准库 | 理解底层原理 |
| 微服务工具集 | Go kit | 灵活，可控 |
