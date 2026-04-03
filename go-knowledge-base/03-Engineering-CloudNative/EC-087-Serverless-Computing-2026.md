# EC-087-Serverless-Computing-2026

> **Dimension**: 03-Engineering-CloudNative
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: 2026 (AWS Lambda, Azure Functions, Knative, WebAssembly)
> **Size**: >20KB

---

## 1. Serverless概览

### 1.1 定义

Serverless是一种云计算执行模型，云提供商动态管理服务器资源分配。

**核心特征**:

- 无服务器管理
- 自动扩展
- 按使用付费
- 事件驱动

### 1.2 架构演变

```
传统架构 → IaaS → PaaS → Serverless

┌─────────────────────────────────────────┐
│           Serverless Architecture       │
├─────────────────────────────────────────┤
│                                         │
│  事件源 → 函数编排 → 函数执行 → 存储    │
│           │                              │
│           └── 自动扩展                   │
│                                         │
│  计费: 请求数 + 执行时间 + 内存          │
│                                         │
└─────────────────────────────────────────┘
```

---

## 2. AWS Lambda 2026

### 2.1 运行时支持

| 运行时 | 版本 | 状态 |
|--------|------|------|
| Node.js | 22.x | 推荐 |
| Python | 3.13 | 推荐 |
| Java | 21 | 推荐 |
| .NET | 9 | 推荐 |
| Go | 1.24 | 原生支持 |
| Ruby | 3.4 | 支持 |
| Rust | - | Custom Runtime |

### 2.2 性能特性

**SnapStart (Java)**:

- 冷启动时间减少90%
- 预初始化运行时
- 适合延迟敏感应用

**Provisioned Concurrency**:

- 预配置执行环境
- 消除冷启动
- 适合稳定流量

**Response Streaming**:

- 支持流式响应
- 最大6MB响应体
- 渐进式内容交付

### 2.3 Go示例

```go
// main.go
package main

import (
    "context"
    "encoding/json"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
    StatusCode int    `json:"statusCode"`
    Body       string `json:"body"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
    // 业务逻辑
    result := map[string]interface{}{
        "message": "Hello from Go Lambda",
        "path":    request.Path,
        "method":  request.HTTPMethod,
    }

    body, _ := json.Marshal(result)

    return Response{
        StatusCode: 200,
        Body:       string(body),
    }, nil
}

func main() {
    lambda.Start(handler)
}
```

### 2.4 部署

```bash
# 编译Linux二进制
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# 打包
zip function.zip bootstrap

# 部署
aws lambda create-function \
    --function-name go-function \
    --runtime provided.al2023 \
    --handler bootstrap \
    --zip-file fileb://function.zip \
    --role arn:aws:iam::ACCOUNT:role/lambda-role
```

---

## 3. Azure Functions 2026

### 3.1 托管选项

| 计划 | 特点 | 适用场景 |
|------|------|---------|
| Consumption | 自动扩展，按执行付费 | 可变流量 |
| Premium | 预配置，VNet支持 | 企业应用 |
| Dedicated | App Service计划 | 现有基础设施 |
| Container | Kubernetes托管 | 容器化工作负载 |

### 3.2 编程模型 v4

```go
// function.go
package main

import (
    "net/http"
    "github.com/Azure/azure-functions-go-worker/worker"
)

func init() {
    worker.RegisterHTTP(httpHandler)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "World"
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "Hello ` + name + `"}`))
}

func main() {
    worker.Start()
}
```

### 3.3 Durable Functions

```python
# Python Durable Functions示例
import azure.functions as func
import azure.durable_functions as df

myApp = df.DFApp(http_auth_level=func.AuthLevel.ANONYMOUS)

@myApp.route(route="orchestrators/{functionName}")
@myApp.durable_client_input(client_name="client")
async def http_start(req: func.HttpRequest, client):
    instance_id = await client.start_new(req.route_params["functionName"], None, None)
    return client.create_check_status_response(req, instance_id)

@myApp.orchestration_trigger(context_name="context")
def my_orchestrator(context):
    result1 = yield context.call_activity("activity1", "input1")
    result2 = yield context.call_activity("activity2", result1)
    return result2
```

---

## 4. Knative

### 4.1 架构

```
┌─────────────────────────────────────────┐
│           Knative Architecture          │
├─────────────────────────────────────────┤
│                                         │
│  ┌──────────────┐    ┌──────────────┐  │
│  │   Serving    │    │    Eventing  │  │
│  │              │    │              │  │
│  │ - Revision   │    │ - Source     │  │
│  │ - Route      │    │ - Broker     │  │
│  │ - Service    │    │ - Trigger    │  │
│  │ - Auto-scale │    │ - Channel    │  │
│  └──────────────┘    └──────────────┘  │
│                                         │
│  基于Kubernetes的Serverless平台        │
└─────────────────────────────────────────┘
```

### 4.2 安装

```bash
# 安装Knative Serving
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.16.0/serving-crds.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.16.0/serving-core.yaml

# 安装网络层 (Istio/Kourier)
kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.16.0/kourier.yaml
```

### 4.3 部署服务

```yaml
# service.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: hello-go
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/minScale: "0"
        autoscaling.knative.dev/maxScale: "100"
    spec:
      containers:
        - image: gcr.io/project/hello-go:latest
          ports:
            - containerPort: 8080
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "200m"
```

### 4.4 事件驱动

```yaml
# trigger.yaml
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: my-trigger
spec:
  broker: default
  filter:
    attributes:
      type: com.example.order.created
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: order-processor
```

---

## 5. WebAssembly Serverless

### 5.1 优势

| 特性 | WebAssembly | 传统容器 |
|------|-------------|---------|
| 启动时间 | <1ms | 100ms-1s |
| 包大小 | KB级 | MB-GB级 |
| 沙箱 | 强隔离 | 内核隔离 |
| 可移植性 | 跨平台 | 依赖镜像 |

### 5.2 Spin框架

```rust
// Rust示例
use spin_sdk::http::{Request, Response, Router};
use spin_sdk::http_component;

#[http_component]
fn handle_route(req: Request) -> Response {
    let mut router = Router::new();
    router.get("/api/hello", handle_hello);
    router.post("/api/echo", handle_echo);
    router.handle(req)
}

fn handle_hello(_req: Request) -> Response {
    Response::builder()
        .status(200)
        .header("content-type", "application/json")
        .body(Some(r#"{"message": "Hello from Wasm!"}"#.into()))
        .build()
}
```

### 5.3 WasmEdge

```go
// TinyGo编译为Wasm
package main

import (
    "github.com/second-state/WasmEdge-go/wasmedge"
)

func main() {
    vm := wasmedge.NewVM()
    vm.LoadWasmFile("function.wasm")
    vm.Validate()
    vm.Instantiate()

    // 执行函数
    res, err := vm.Execute("handler", []byte(`{"name": "world"}`))

    vm.Release()
}
```

---

## 6. 冷启动优化

### 6.1 冷启动来源

```
冷启动阶段:
1. 调度 (~100ms) - 容器分配
2. 下载镜像 (~500ms-5s) - 镜像拉取
3. 启动容器 (~100ms) - 运行时初始化
4. 运行时初始化 (~100ms-1s) - 代码加载
5. 业务逻辑 (~variable) - 连接池等

总冷启动时间: 1-10秒
```

### 6.2 优化策略

**代码层面**:

```go
// 延迟初始化
var db *sql.DB
var dbOnce sync.Once

func getDB() *sql.DB {
    dbOnce.Do(func() {
        db, _ = sql.Open("postgres", connStr)
        db.SetMaxIdleConns(2)
        db.SetMaxOpenConns(5)
    })
    return db
}

// 避免全局初始化
func handler(ctx context.Context, event Event) error {
    db := getDB()  // 首次调用时初始化
    // ...
}
```

**架构层面**:

- 使用Provisioned Concurrency
- 保持函数活跃 (定时ping)
- 精简依赖
- 使用Lambda Layers共享代码

---

## 7. Serverless模式

### 7.1 常见模式

```
1. 事件处理
Event Source → Lambda → Data Store

2. API后端
API Gateway → Lambda → Database

3. 工作流编排
Step Functions → 多个Lambda → 结果聚合

4. 实时处理
Kinesis Stream → Lambda → Analytics

5. 定时任务
EventBridge → Lambda → 业务逻辑
```

### 7.2 Saga模式

```python
# AWS Step Functions Saga模式
{
    "Comment": "Order Processing Saga",
    "StartAt": "ReserveInventory",
    "States": {
        "ReserveInventory": {
            "Type": "Task",
            "Resource": "arn:aws:lambda:...:reserve-inventory",
            "Next": "ProcessPayment",
            "Catch": [{
                "ErrorEquals": ["States.ALL"],
                "Next": "CompensateInventory"
            }]
        },
        "ProcessPayment": {
            "Type": "Task",
            "Resource": "arn:aws:lambda:...:process-payment",
            "Next": "ShipOrder",
            "Catch": [{
                "ErrorEquals": ["States.ALL"],
                "Next": "RefundPayment"
            }]
        },
        "CompensateInventory": {
            "Type": "Task",
            "Resource": "arn:aws:lambda:...:release-inventory"
        },
        "RefundPayment": {
            "Type": "Task",
            "Resource": "arn:aws:lambda:...:refund-payment",
            "Next": "CompensateInventory"
        }
    }
}
```

---

## 8. 监控和可观测性

### 8.1 AWS Lambda Insights

```yaml
# CloudWatch Lambda Insights
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: ./src
      Runtime: provided.al2023
      Architectures:
        - x86_64
      Policies:
        - CloudWatchLambdaInsightsExecutionRolePolicy
      Layers:
        - !Sub 'arn:aws:lambda:${AWS::Region}:580247275435:layer:LambdaInsightsExtension:38'
```

### 8.2 分布式追踪

```go
// OpenTelemetry追踪
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func handler(ctx context.Context, event Event) error {
    tracer := otel.Tracer("lambda-function")

    ctx, span := tracer.Start(ctx, "process-event")
    defer span.End()

    // 业务逻辑
    processData(ctx, event.Data)

    return nil
}
```

---

## 9. 安全最佳实践

### 9.1 最小权限原则

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "dynamodb:GetItem",
                "dynamodb:PutItem"
            ],
            "Resource": "arn:aws:dynamodb:*:*:table/Orders",
            "Condition": {
                "ForAllValues:StringEquals": {
                    "dynamodb:LeadingKeys": ["${aws:userid}"]
                }
            }
        }
    ]
}
```

### 9.2 秘密管理

```python
# AWS Secrets Manager
import boto3
import json

secrets = boto3.client('secretsmanager')

def get_secret():
    response = secrets.get_secret_value(SecretId='prod/db/password')
    return json.loads(response['SecretString'])

# 缓存秘密
_db_password = None

def get_db_password():
    global _db_password
    if _db_password is None:
        _db_password = get_secret()['password']
    return _db_password
```

---

## 10. 成本优化

### 10.1 计费模型

```
AWS Lambda计费:
- 请求数: $0.20 per 1M requests
- 执行时间: $0.0000166667 per GB-second

示例: 1000万次调用，每次128MB，200ms
- 请求: 10M × $0.20/M = $2.00
- 执行: 10M × 0.2s × 0.128GB × $0.0000166667/GB-s = $4.27
- 总计: ~$6.27
```

### 10.2 优化策略

| 策略 | 效果 |
|------|------|
| 内存优化 | 找到最优内存/性能平衡点 |
| 批量处理 | 减少调用次数 |
| 异步处理 | 使用队列缓冲 |
| 预留并发 | 适合稳定工作负载 |

---

## 11. 参考文献

1. AWS Lambda Documentation
2. Azure Functions Best Practices
3. Knative Documentation
4. WebAssembly Serverless Patterns
5. Serverless Framework Guide

---

*Last Updated: 2026-04-03*
