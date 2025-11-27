# Command Applications (应用入口)

> **简介**: 本文档说明项目中的应用入口程序，包括 HTTP 服务器、gRPC 服务器、GraphQL 服务器等。

**版本**: v1.0
**更新日期**: 2025-01-XX
**适用于**: Go 1.25.3

---

## 📋 目录

- [Command Applications (应用入口)](#command-applications-应用入口)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 应用入口结构](#2-应用入口结构)
  - [3. 启动流程说明](#3-启动流程说明)
    - [3.1 配置加载](#31-配置加载)
    - [3.2 日志初始化](#32-日志初始化)
    - [3.3 数据库初始化](#33-数据库初始化)
    - [3.4 依赖注入](#34-依赖注入)
    - [3.5 优雅关闭](#35-优雅关闭)
  - [4. 各应用入口说明](#4-各应用入口说明)
    - [4.1 HTTP 服务器 (server)](#41-http-服务器-server)
    - [4.2 gRPC 服务器 (grpc-server)](#42-grpc-服务器-grpc-server)
    - [4.3 GraphQL 服务器 (graphql-server)](#43-graphql-服务器-graphql-server)
    - [4.4 Temporal Worker (temporal-worker)](#44-temporal-worker-temporal-worker)
    - [4.5 MQTT 客户端 (mqtt-client)](#45-mqtt-客户端-mqtt-client)
    - [4.6 CLI 工具 (cli)](#46-cli-工具-cli)
  - [5. 依赖注入说明](#5-依赖注入说明)
    - [5.1 当前实现（手动注入）](#51-当前实现手动注入)
    - [5.2 推荐实现（Wire 依赖注入）](#52-推荐实现wire-依赖注入)
  - [6. 优雅关闭](#6-优雅关闭)
    - [6.1 信号监听](#61-信号监听)
    - [6.2 关闭超时](#62-关闭超时)
    - [6.3 资源清理](#63-资源清理)
  - [📚 相关资源](#-相关资源)

---

## 1. 概述

`cmd/` 目录包含所有可执行程序的入口点。每个子目录都是一个独立的应用程序，可以单独编译和运行。

**设计原则**：

1. **单一职责**：每个应用入口只负责一种类型的服务
2. **独立部署**：每个应用可以独立编译、部署和运行
3. **共享代码**：所有应用共享 `internal/` 和 `pkg/` 中的代码
4. **统一配置**：所有应用使用相同的配置系统

---

## 2. 应用入口结构

```text
cmd/
├── server/           # HTTP 服务器（REST API）
│   └── main.go       # HTTP 服务器入口
├── grpc-server/      # gRPC 服务器
│   └── main.go       # gRPC 服务器入口
├── graphql-server/   # GraphQL 服务器
│   └── main.go       # GraphQL 服务器入口
├── temporal-worker/  # Temporal 工作流 Worker
│   └── main.go       # Worker 入口
├── mqtt-client/      # MQTT 客户端
│   └── main.go       # MQTT 客户端入口
└── cli/              # CLI 命令行工具
    └── main.go       # CLI 工具入口
```

---

## 3. 启动流程说明

所有应用入口都遵循类似的启动流程：

```text
1. 加载配置
   ↓
2. 初始化日志
   ↓
3. 初始化可观测性（OpenTelemetry）
   ↓
4. 初始化数据库连接
   ↓
5. 初始化依赖（使用依赖注入或手动组装）
   ↓
6. 创建服务实例（HTTP/gRPC/GraphQL 等）
   ↓
7. 启动服务（在 goroutine 中）
   ↓
8. 等待中断信号（SIGINT/SIGTERM）
   ↓
9. 优雅关闭（等待请求完成，关闭连接）
```

**关键步骤说明**：

### 3.1 配置加载

```go
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal("Failed to load config", err)
}
```

- 从配置文件或环境变量加载配置
- 配置位置：`configs/config.yaml` 或环境变量

### 3.2 日志初始化

```go
logger := otlp.NewLogger()
slog.SetDefault(logger.Logger)
```

- 使用结构化日志（slog）
- 支持 JSON 和 Text 格式
- 集成 OpenTelemetry

### 3.3 数据库初始化

```go
entClient, err := entdb.NewClientFromConfig(...)
defer entClient.Close()
```

- 使用 Ent ORM 连接数据库
- 自动运行数据库迁移
- 支持连接池管理

### 3.4 依赖注入

**方式1：手动依赖注入（当前实现）**:

```go
userRepo := entrepo.NewUserRepository(entClient)
userService := appuser.NewService(userRepo)
router := chiRouter.NewRouter(userService)
```

**方式2：Wire 依赖注入（推荐）**:

```go
app, err := wire.InitializeApp(cfg)
if err != nil {
    log.Fatal(err)
}
```

### 3.5 优雅关闭

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

server.Shutdown(ctx)
```

- 监听系统信号（SIGINT、SIGTERM）
- 等待正在处理的请求完成
- 关闭数据库连接和其他资源

---

## 4. 各应用入口说明

### 4.1 HTTP 服务器 (server)

**功能**：提供 REST API 服务

**启动命令**：

```bash
go run cmd/server/main.go
# 或
go build -o bin/server cmd/server/main.go
./bin/server
```

**主要功能**：

- HTTP REST API 服务
- 使用 Chi Router
- 支持中间件（认证、限流、熔断等）
- 集成 OpenTelemetry 追踪
- 支持健康检查

**端口**：默认 `8080`（可在配置文件中修改）

**依赖**：

- 数据库（PostgreSQL）
- OpenTelemetry（可选）
- Temporal（可选）

### 4.2 gRPC 服务器 (grpc-server)

**功能**：提供 gRPC 服务

**启动命令**：

```bash
go run cmd/grpc-server/main.go
```

**主要功能**：

- gRPC 服务
- Protocol Buffers 序列化
- 支持流式 RPC
- 集成 OpenTelemetry

**端口**：默认 `8081`（Server.Port + 1）

**依赖**：

- 数据库（PostgreSQL）
- gRPC 服务定义（`.proto` 文件）

### 4.3 GraphQL 服务器 (graphql-server)

**功能**：提供 GraphQL API 服务

**启动命令**：

```bash
go run cmd/graphql-server/main.go
```

**主要功能**：

- GraphQL API 服务
- 支持查询（Query）和变更（Mutation）
- 支持订阅（Subscription）
- Schema 定义在 `api/graphql/schema.graphql`

**端口**：默认 `8082`（可在配置文件中修改）

### 4.4 Temporal Worker (temporal-worker)

**功能**：执行 Temporal 工作流和活动

**启动命令**：

```bash
go run cmd/temporal-worker/main.go
```

**主要功能**：

- 执行工作流（Workflow）
- 执行活动（Activity）
- 支持长时间运行的任务
- 支持重试和错误处理

**依赖**：

- Temporal 服务器
- 数据库（用于活动执行）

**配置**：

- `Temporal.Address`: Temporal 服务器地址
- `Temporal.TaskQueue`: 任务队列名称

### 4.5 MQTT 客户端 (mqtt-client)

**功能**：MQTT 消息客户端

**启动命令**：

```bash
go run cmd/mqtt-client/main.go
```

**主要功能**：

- 订阅 MQTT 主题
- 发布 MQTT 消息
- 处理消息回调

**配置**：

- `MQTT.Broker`: MQTT Broker 地址
- `MQTT.ClientID`: 客户端 ID
- `MQTT.Username`: 用户名（可选）
- `MQTT.Password`: 密码（可选）

### 4.6 CLI 工具 (cli)

**功能**：命令行工具

**启动命令**：

```bash
go run cmd/cli/main.go [command] [flags]
```

**主要功能**：

- 数据库迁移
- 数据导入/导出
- 系统管理任务
- 开发工具

---

## 5. 依赖注入说明

### 5.1 当前实现（手动注入）

当前代码使用手动依赖注入：

```go
// 1. 创建数据库客户端
entClient, err := entdb.NewClientFromConfig(...)

// 2. 创建仓储
userRepo := entrepo.NewUserRepository(entClient)

// 3. 创建应用服务
userService := appuser.NewService(userRepo)

// 4. 创建路由
router := chiRouter.NewRouter(userService)
```

**优点**：

- 简单直接
- 易于理解
- 不需要额外工具

**缺点**：

- 依赖关系分散
- 难以管理复杂依赖
- 容易出错

### 5.2 推荐实现（Wire 依赖注入）

推荐使用 Wire 进行依赖注入：

```go
// scripts/wire/wire.go
func InitializeApp(cfg *config.Config) (*App, error) {
    wire.Build(
        NewEntClient,
        NewUserRepository,
        NewUserService,
        NewRouter,
        NewApp,
    )
    return &App{}, nil
}

// cmd/server/main.go
app, err := wire.InitializeApp(cfg)
if err != nil {
    log.Fatal(err)
}
```

**优点**：

- 编译时检查依赖关系
- 类型安全
- 自动生成代码
- 易于维护

**迁移步骤**：

1. 在 `scripts/wire/wire.go` 中定义 Provider 函数
2. 运行 `go generate ./scripts/wire` 生成代码
3. 在 `main.go` 中使用生成的 `InitializeApp` 函数

详细说明请参考：[架构模型与依赖注入完整说明](../../docs/architecture/00-架构模型与依赖注入完整说明.md)

---

## 6. 优雅关闭

所有应用都实现了优雅关闭机制：

### 6.1 信号监听

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
```

- 监听 `SIGINT`（Ctrl+C）和 `SIGTERM`（kill 命令）
- 收到信号后开始关闭流程

### 6.2 关闭超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

- 设置 30 秒超时
- 如果超时，强制关闭

### 6.3 资源清理

```go
// 关闭 HTTP 服务器
httpServer.Shutdown(ctx)

// 关闭数据库连接
entClient.Close()

// 关闭 Temporal 客户端
temporalClient.Close()

// 关闭事件总线
eventBus.Stop()
```

**关闭顺序**：

1. 停止接收新请求
2. 等待正在处理的请求完成
3. 关闭数据库连接
4. 关闭其他资源

---

## 📚 相关资源

- [架构模型与依赖注入完整说明](../../docs/architecture/00-架构模型与依赖注入完整说明.md)
- [Wire 依赖注入文档](../../docs/architecture/tech-stack/config/wire.md)
- [配置管理文档](../../internal/config/README.md)
- [框架使用指南](../../docs/00-框架使用指南.md)

---

> 📚 **总结**
> `cmd/` 目录包含所有应用的入口点。每个应用都遵循统一的启动流程，支持配置管理、日志记录、依赖注入和优雅关闭。推荐使用 Wire 进行依赖注入，以提高代码的可维护性和可测试性。
