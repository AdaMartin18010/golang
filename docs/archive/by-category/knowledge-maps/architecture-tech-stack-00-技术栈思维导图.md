# 技术栈思维导图

> **简介**: 本文档通过思维导图、关系网络图、决策流程图等多种可视化方式，全面展示项目技术栈的完整视图。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26

---

## 📋 目录

- [技术栈思维导图](#技术栈思维导图)
  - [📋 目录](#-目录)
  - [1. 🗺️ 技术栈全景思维导图](#1-️-技术栈全景思维导图)
    - [1.1 完整技术栈思维导图](#11-完整技术栈思维导图)
  - [2. 🔗 技术栈关系网络图](#2--技术栈关系网络图)
    - [2.1 技术栈依赖关系图](#21-技术栈依赖关系图)
    - [2.2 技术栈协作关系图](#22-技术栈协作关系图)
  - [3. 🎯 技术栈决策流程图](#3--技术栈决策流程图)
    - [3.1 Web 框架选型决策流程](#31-web-框架选型决策流程)
    - [3.2 数据访问层选型决策流程](#32-数据访问层选型决策流程)
    - [3.3 消息队列选型决策流程](#33-消息队列选型决策流程)
  - [4. 📊 技术栈分层架构图](#4--技术栈分层架构图)
    - [4.1 技术栈分层视图](#41-技术栈分层视图)
  - [5. 🔄 技术栈交互时序图](#5--技术栈交互时序图)
    - [5.1 HTTP 请求处理时序图](#51-http-请求处理时序图)
    - [5.2 工作流执行时序图](#52-工作流执行时序图)
  - [6. 📈 技术栈选型对比矩阵](#6--技术栈选型对比矩阵)
    - [6.1 Web 框架对比](#61-web-框架对比)
    - [6.2 数据访问对比](#62-数据访问对比)
    - [6.3 消息队列对比](#63-消息队列对比)
  - [7. 📚 扩展阅读](#7--扩展阅读)

---

## 1. 🗺️ 技术栈全景思维导图

### 1.1 完整技术栈思维导图

```mermaid
mindmap
  root((Go 技术栈))
    Web框架层
      Chi Router
        ✅ 标准库兼容
        ✅ 轻量级
        ✅ 高性能
        ✅ 中间件生态
      WebSocket
        ✅ 实时通信
        ✅ 双向通信
        ✅ 低延迟
    HTTP/3
        ✅ 最新协议
        ✅ 性能优化
   数据访问层
      Ent ORM
        ✅ 类型安全
        ✅ 代码生成
        ✅ Schema 定义
      PostgreSQL pgx
        ✅ 高性能
        ✅ 功能完整
        ✅ JSON 支持
      SQLite
        ✅ 轻量级
        ✅ 零配置
        ✅ 嵌入式
    工作流层
      Temporal
        ✅ Go SDK 官方
        ✅ 持久化
        ✅ 可观测性
    可观测性层
      OpenTelemetry
        ✅ 行业标准
        ✅ 三支柱支持
        ✅ 生态丰富
      Prometheus
        ✅ 监控标准
        ✅ 高性能
        ✅ 生态丰富
      Grafana
        ✅ 可视化强大
        ✅ 数据源丰富
        ✅ 社区活跃
      Jaeger
        ✅ OTLP 集成
        ✅ 可视化完善
        ✅ 开源免费
      eBPF
        ✅ 内核执行
        ✅ 高性能
        ✅ 安全性
    消息队列层
      Kafka
        ✅ 高吞吐量
        ✅ 持久化
        ✅ 分区支持
      MQTT
        ✅ IoT 适配
        ✅ 轻量级
        ✅ QoS 支持
      NATS
        ✅ 高性能
        ✅ 低延迟
        ✅ 云原生
    配置工具层
      Viper
        ✅ 功能完整
        ✅ 多格式支持
        ✅ 易用性
      Slog
        ✅ 标准库
        ✅ 结构化日志
        ✅ 性能优秀
      Wire
        ✅ 编译时注入
        ✅ 类型安全
        ✅ 无运行时开销
    工具层
      UUID
        ✅ 全局唯一
        ✅ 无需协调
        ✅ 标准化
      Cron
        ✅ 时间调度
        ✅ 并发安全
        ✅ 任务管理
     API协议层
      gRPC
        ✅ 高性能
        ✅ 类型安全
        ✅ 流式处理
      GraphQL
        ✅ 查询灵活
        ✅ 类型系统
        ✅ 客户端控制
      OpenAPI
        ✅ 标准化
        ✅ 代码生成
        ✅ 文档生成
      AsyncAPI
        ✅ 异步支持
        ✅ 多协议
        ✅ 事件驱动
      Protocol Buffers
        ✅ 高效
        ✅ 类型安全
        ✅ 版本兼容
      gRPC Gateway
        ✅ 协议转换
        ✅ 代码生成
        ✅ 统一接口
```

---

## 2. 🔗 技术栈关系网络图

### 2.1 技术栈依赖关系图

```mermaid
graph TB
    subgraph "Web 层"
        Chi[Chi Router]
        WS[WebSocket]
        HTTP3[HTTP/3]
    end

    subgraph "API 层"
        gRPC[gRPC]
        GraphQL[GraphQL]
        OpenAPI[OpenAPI]
        AsyncAPI[AsyncAPI]
        Protobuf[Protocol Buffers]
        Gateway[gRPC Gateway]
    end

    subgraph "数据层"
        Ent[Ent ORM]
        PGX[PostgreSQL pgx]
        SQLite[SQLite]
    end

    subgraph "消息层"
        Kafka[Kafka]
        MQTT[MQTT]
        NATS[NATS]
    end

    subgraph "工作流层"
        Temporal[Temporal]
    end

    subgraph "可观测性层"
        OTel[OpenTelemetry]
        Prom[Prometheus]
        Grafana[Grafana]
        Jaeger[Jaeger]
        eBPF[eBPF]
    end

    subgraph "配置层"
        Viper[Viper]
        Slog[Slog]
        Wire[Wire]
    end

    subgraph "工具层"
        UUID[UUID]
        Cron[Cron]
    end

    Chi --> gRPC
    Chi --> GraphQL
    Chi --> WS
    gRPC --> Protobuf
    Gateway --> gRPC
    Gateway --> Protobuf
    GraphQL --> OpenAPI
    AsyncAPI --> Kafka
    AsyncAPI --> MQTT
    AsyncAPI --> NATS

    Ent --> PGX
    Ent --> SQLite

    Temporal --> OTel
    Chi --> OTel
    gRPC --> OTel
    Kafka --> OTel

    OTel --> Prom
    OTel --> Jaeger
    Prom --> Grafana
    Jaeger --> Grafana

    Chi --> Viper
    Chi --> Slog
    Chi --> Wire

    Ent --> UUID
    Temporal --> Cron

    style Chi fill:#e1f5ff
    style gRPC fill:#fff4e1
    style Ent fill:#e8f5e9
    style OTel fill:#f3e5f5
```

### 2.2 技术栈协作关系图

```mermaid
graph LR
    subgraph "请求处理流程"
        A[客户端] -->|HTTP| B[Chi Router]
        B -->|路由| C[Handler]
        C -->|调用| D[Service]
        D -->|使用| E[Repository]
        E -->|查询| F[Ent ORM]
        F -->|SQL| G[PostgreSQL]
    end

    subgraph "可观测性流程"
        C -->|追踪| H[OpenTelemetry]
        D -->|追踪| H
        E -->|追踪| H
        H -->|指标| I[Prometheus]
        H -->|追踪| J[Jaeger]
        I -->|可视化| K[Grafana]
        J -->|可视化| K
    end

    subgraph "异步处理流程"
        D -->|发布| L[Kafka]
        L -->|消费| M[Consumer]
        M -->|处理| D
    end

    style B fill:#e1f5ff
    style D fill:#fff4e1
    style H fill:#f3e5f5
```

---

## 3. 🎯 技术栈决策流程图

### 3.1 Web 框架选型决策流程

```mermaid
flowchart TD
    Start([需要 Web 框架]) --> Q1{需要标准库兼容?}
    Q1 -->|是| Q2{需要轻量级?}
    Q1 -->|否| Q3{需要全功能框架?}

    Q2 -->|是| Chi[选择 Chi Router]
    Q2 -->|否| Q4{需要高性能?}

    Q3 -->|是| Gin[考虑 Gin/Echo]
    Q3 -->|否| Q5{需要最新特性?}

    Q4 -->|是| Chi
    Q4 -->|否| Q6{需要中间件生态?}

    Q5 -->|是| HTTP3[考虑 HTTP/3]
    Q5 -->|否| Standard[使用标准库]

    Q6 -->|是| Chi
    Q6 -->|否| Standard

    Chi --> End([Chi Router])
    Gin --> End2([Gin/Echo])
    HTTP3 --> End3([HTTP/3])
    Standard --> End4([标准库])
```

### 3.2 数据访问层选型决策流程

```mermaid
flowchart TD
    Start([需要数据访问]) --> Q1{需要类型安全?}
    Q1 -->|是| Q2{需要代码生成?}
    Q1 -->|否| Q3{需要原生 SQL?}

    Q2 -->|是| Ent[选择 Ent ORM]
    Q2 -->|否| GORM[考虑 GORM]

    Q3 -->|是| PGX[选择 pgx]
    Q3 -->|否| Q4{需要轻量级?}

    Q4 -->|是| SQLite[选择 SQLite]
    Q4 -->|否| PGX

    Ent --> Q5{数据库类型?}
    Q5 -->|PostgreSQL| PGX
    Q5 -->|SQLite| SQLite

    Ent --> End([Ent ORM])
    PGX --> End2([PostgreSQL pgx])
    SQLite --> End3([SQLite])
```

### 3.3 消息队列选型决策流程

```mermaid
flowchart TD
    Start([需要消息队列]) --> Q1{需要高吞吐量?}
    Q1 -->|是| Q2{需要持久化?}
    Q1 -->|否| Q3{需要低延迟?}

    Q2 -->|是| Kafka[选择 Kafka]
    Q2 -->|否| Q4{需要 IoT 支持?}

    Q3 -->|是| NATS[选择 NATS]
    Q3 -->|否| Q5{需要轻量级?}

    Q4 -->|是| MQTT[选择 MQTT]
    Q4 -->|否| NATS

    Q5 -->|是| NATS
    Q5 -->|否| Kafka

    Kafka --> End([Kafka])
    MQTT --> End2([MQTT])
    NATS --> End3([NATS])
```

---

## 4. 📊 技术栈分层架构图

### 4.1 技术栈分层视图

```mermaid
graph TB
    subgraph "接口层 - Interfaces Layer"
        HTTP[HTTP API<br/>Chi Router]
        gRPC[gRPC API<br/>Protocol Buffers]
        GraphQL[GraphQL API<br/>Schema]
        WS[WebSocket<br/>实时通信]
    end

    subgraph "应用层 - Application Layer"
        Service[Services<br/>业务逻辑]
        Workflow[Workflows<br/>Temporal]
        DTO[DTOs<br/>数据传输]
    end

    subgraph "领域层 - Domain Layer"
        Entity[Entities<br/>领域实体]
        Repo[Repository<br/>接口定义]
    end

    subgraph "基础设施层 - Infrastructure Layer"
        DB[Database<br/>Ent/PostgreSQL/SQLite]
        MQ[Message Queue<br/>Kafka/MQTT/NATS]
        Obs[Observability<br/>OTel/Prometheus/Grafana]
        Config[Configuration<br/>Viper/Slog/Wire]
    end

    HTTP --> Service
    gRPC --> Service
    GraphQL --> Service
    WS --> Service

    Service --> Entity
    Service --> Repo
    Workflow --> Service

    Repo --> DB
    Service --> MQ
    Service --> Obs
    Service --> Config

    style HTTP fill:#e1f5ff
    style Service fill:#fff4e1
    style Entity fill:#e8f5e9
    style DB fill:#f3e5f5
```

---

## 5. 🔄 技术栈交互时序图

### 5.1 HTTP 请求处理时序图

```mermaid
sequenceDiagram
    participant Client
    participant Chi as Chi Router
    participant Middleware
    participant Handler
    participant Service
    participant Repository
    participant Ent
    participant DB as PostgreSQL
    participant OTel as OpenTelemetry

    Client->>Chi: HTTP Request
    Chi->>Middleware: 中间件处理
    Middleware->>OTel: 开始追踪
    Middleware->>Handler: 调用 Handler
    Handler->>Service: 调用 Service
    Service->>Repository: 调用 Repository
    Repository->>Ent: 使用 Ent ORM
    Ent->>DB: SQL 查询
    DB-->>Ent: 返回数据
    Ent-->>Repository: 返回实体
    Repository-->>Service: 返回实体
    Service->>Service: 转换为 DTO
    Service-->>Handler: 返回 DTO
    Handler-->>Middleware: 返回响应
    Middleware->>OTel: 结束追踪
    Middleware-->>Chi: 响应
    Chi-->>Client: HTTP Response
```

### 5.2 工作流执行时序图

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Temporal as Temporal Client
    participant Server as Temporal Server
    participant Worker
    participant Activity
    participant Service

    Client->>Handler: 启动工作流请求
    Handler->>Temporal: 启动工作流
    Temporal->>Server: 创建工作流实例
    Server-->>Temporal: 返回工作流 ID
    Temporal-->>Handler: 返回工作流 ID
    Handler-->>Client: 返回工作流 ID

    Server->>Worker: 调度工作流
    Worker->>Activity: 执行 Activity
    Activity->>Service: 调用业务服务
    Service-->>Activity: 返回结果
    Activity-->>Worker: 返回结果
    Worker-->>Server: 更新工作流状态
```

---

## 6. 📈 技术栈选型对比矩阵

### 6.1 Web 框架对比

| 特性 | Chi Router | Gin | Echo | 标准库 |
|------|-----------|-----|------|--------|
| 性能 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 易用性 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ |
| 中间件 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐ |
| 标准库兼容 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| 社区 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

### 6.2 数据访问对比

| 特性 | Ent ORM | GORM | pgx | SQLite |
|------|---------|------|-----|--------|
| 类型安全 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 性能 | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 易用性 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 代码生成 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐ | ⭐ |
| 功能完整性 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |

### 6.3 消息队列对比

| 特性 | Kafka | MQTT | NATS |
|------|-------|------|------|
| 吞吐量 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| 延迟 | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 持久化 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ |
| 易用性 | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| IoT 支持 | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |

---

## 7. 📚 扩展阅读

- [技术栈概览](./00-技术栈概览.md)
- [技术栈集成](./01-技术栈集成.md)
- [技术栈选型决策树](./02-技术栈选型决策树.md)
- [技术栈文档索引](./README.md)
- [架构知识图谱](../00-知识图谱.md)

---

> 📚 **简介**
> 本文档通过多种可视化方式全面展示项目技术栈，包括思维导图、关系网络图、决策流程图等。
