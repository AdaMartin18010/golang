# Interfaces Layer (接口层)

> **简介**: 本文档说明 Clean Architecture 的接口层，包含外部接口适配、请求处理和响应格式化。

**版本**: v1.0
**更新日期**: 2025-01-XX
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 概述](#1-概述)
- [2. 接口层结构](#2-接口层结构)
- [3. 设计原则](#3-设计原则)
- [4. 各协议实现说明](#4-各协议实现说明)
  - [4.1 HTTP 接口 (http/chi)](#41-http-接口-httpchi)
  - [4.2 gRPC 接口 (grpc)](#42-grpc-接口-grpc)
  - [4.3 GraphQL 接口 (graphql)](#43-graphql-接口-graphql)
  - [4.4 AsyncAPI (asyncapi)](#44-asyncapi-asyncapi)
- [5. 请求处理流程](#5-请求处理流程)
- [6. 错误处理](#6-错误处理)

---

## 1. 概述

接口层是 Clean Architecture 的最外层，负责外部接口适配、请求处理和响应格式化。

**核心职责**：

1. **协议适配**：适配不同的外部协议（HTTP、gRPC、GraphQL）
2. **请求处理**：接收和处理外部请求
3. **响应格式化**：格式化响应数据
4. **参数验证**：验证请求参数
5. **错误处理**：处理错误并格式化错误响应

**设计原则**：

- ✅ 调用 Application Layer，不直接访问 Domain Layer
- ✅ 处理外部请求，适配不同协议
- ✅ 负责请求/响应转换（DTO 转换）
- ✅ 处理协议相关的错误和状态码

---

## 2. 接口层结构

```text
interfaces/
├── http/          # HTTP 接口
│   ├── chi/       # Chi 路由实现
│   │   ├── router.go        # 路由配置
│   │   ├── middleware.go    # 中间件
│   │   ├── handlers/        # HTTP 处理器
│   │   └── middleware/      # 自定义中间件
│   ├── echo/      # Echo 路由（可选）
│   └── openapi/   # OpenAPI 规范定义
├── grpc/          # gRPC 接口
│   ├── handlers/  # gRPC 处理器
│   └── proto/     # Protocol Buffers 定义
├── graphql/       # GraphQL 接口
│   └── resolver.go # GraphQL 解析器
└── asyncapi/      # AsyncAPI 规范定义
    └── asyncapi.yaml
```

---

## 3. 设计原则

### 3.1 依赖方向

**依赖规则**：

```
Interfaces Layer → Application Layer → Domain Layer
```

- ✅ Interfaces Layer 只能调用 Application Layer
- ❌ Interfaces Layer 不能直接访问 Domain Layer
- ❌ Interfaces Layer 不能访问 Infrastructure Layer

### 3.2 职责分离

**接口层职责**：

1. **协议处理**：处理 HTTP、gRPC、GraphQL 等协议细节
2. **请求解析**：解析请求参数、请求体等
3. **响应格式化**：格式化响应数据（JSON、Protobuf 等）
4. **错误转换**：将业务错误转换为协议错误

**应用层职责**：

1. **业务编排**：协调领域对象完成业务用例
2. **DTO 转换**：领域对象和 DTO 之间的转换
3. **事务管理**：管理事务边界

### 3.3 协议隔离

**设计原理**：

- 接口层隔离协议变化
- 应用层不关心协议细节
- 可以轻松添加新协议支持

**示例**：

```go
// 同一个应用服务，可以被不同协议调用
// HTTP 接口
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 解析 HTTP 请求
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)

    // 调用应用服务
    user, err := h.service.CreateUser(r.Context(), req)

    // 格式化 HTTP 响应
    response.Success(w, http.StatusCreated, user)
}

// gRPC 接口
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
    // 调用应用服务
    user, err := s.service.CreateUser(ctx, toAppRequest(req))

    // 转换为 gRPC 响应
    return toGRPCResponse(user), err
}
```

---

## 4. 各协议实现说明

### 4.1 HTTP 接口 (http/chi)

**位置**：`internal/interfaces/http/chi/`

**功能**：

- REST API 服务
- 使用 Chi Router
- 支持中间件（认证、限流、熔断等）
- 集成 OpenTelemetry 追踪

**主要组件**：

1. **Router** (`router.go`)：路由配置和管理
2. **Middleware** (`middleware.go`)：中间件实现
3. **Handlers** (`handlers/`)：HTTP 处理器

**路由结构**：

```
/health                    # 健康检查
/api/v1/users             # 用户相关 API
/api/v1/workflows         # 工作流相关 API
```

**中间件**：

- RequestID：生成请求ID
- RealIP：获取真实IP
- Tracing：OpenTelemetry 追踪
- Logging：请求日志
- Recovery：Panic 恢复
- Timeout：请求超时
- CORS：跨域支持

### 4.2 gRPC 接口 (grpc)

**位置**：`internal/interfaces/grpc/`

**功能**：

- gRPC 服务
- Protocol Buffers 序列化
- 支持流式 RPC

**主要组件**：

1. **Handlers** (`handlers/`)：gRPC 处理器
2. **Proto** (`proto/`)：Protocol Buffers 定义

**使用流程**：

1. 定义 `.proto` 文件
2. 运行 `protoc` 生成代码
3. 实现 gRPC 服务接口
4. 注册服务到 gRPC 服务器

### 4.3 GraphQL 接口 (graphql)

**位置**：`internal/interfaces/graphql/`

**功能**：

- GraphQL API 服务
- 支持查询（Query）和变更（Mutation）
- 支持订阅（Subscription）

**主要组件**：

1. **Resolver** (`resolver.go`)：GraphQL 解析器
2. **Schema** (`api/graphql/schema.graphql`)：GraphQL Schema 定义

### 4.4 AsyncAPI (asyncapi)

**位置**：`internal/interfaces/asyncapi/` 和 `api/asyncapi/`

**功能**：

- 异步 API 规范定义
- 消息队列接口定义
- 事件驱动架构文档

**主要组件**：

1. **AsyncAPI 规范** (`asyncapi.yaml`)：异步 API 定义

---

## 5. 请求处理流程

**HTTP 请求处理流程**：

```
1. HTTP 请求到达
   ↓
2. 中间件链处理
   - RequestID 生成
   - 追踪上下文提取
   - 日志记录
   ↓
3. 路由匹配
   ↓
4. Handler 处理
   - 解析请求参数
   - 验证请求数据
   - 调用应用服务
   ↓
5. 应用服务处理
   - 业务逻辑编排
   - 调用领域对象
   ↓
6. 响应格式化
   - 转换为 DTO
   - 序列化为 JSON
   ↓
7. HTTP 响应返回
```

**错误处理流程**：

```
1. 业务错误发生
   ↓
2. 应用层返回错误
   ↓
3. Handler 捕获错误
   ↓
4. 错误类型判断
   - 业务错误 → 400/404/409
   - 系统错误 → 500
   ↓
5. 错误响应格式化
   - 转换为错误 DTO
   - 序列化为 JSON
   ↓
6. HTTP 错误响应返回
```

---

## 6. 错误处理

**错误处理原则**：

1. **错误分类**：区分业务错误和系统错误
2. **状态码映射**：将业务错误映射到合适的 HTTP 状态码
3. **错误信息**：返回明确的错误信息
4. **错误追踪**：记录错误日志和追踪信息

**错误类型映射**：

| 业务错误类型 | HTTP 状态码 | 说明 |
|------------|------------|------|
| 验证错误 | 400 Bad Request | 请求参数无效 |
| 未找到 | 404 Not Found | 资源不存在 |
| 冲突 | 409 Conflict | 资源冲突（如重复创建） |
| 未授权 | 401 Unauthorized | 未认证 |
| 禁止访问 | 403 Forbidden | 无权限 |
| 系统错误 | 500 Internal Server Error | 服务器内部错误 |

---

## 📚 相关资源

- [Clean Architecture 详解](../../docs/architecture/clean-architecture.md)
- [架构模型与依赖注入完整说明](../../docs/architecture/00-架构模型与依赖注入完整说明.md)
- [Chi Router 文档](../../docs/architecture/tech-stack/web/chi-router.md)
- [HTTP 接口实现示例](../../examples/framework-usage/)

---

> 📚 **总结**
> 接口层负责外部接口适配，调用应用层服务，处理协议细节。通过清晰的职责分离和协议隔离，可以轻松添加新的协议支持，同时保持应用层的稳定性。
