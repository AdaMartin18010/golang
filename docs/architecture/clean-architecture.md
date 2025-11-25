# Clean Architecture

> **简介**: 本文档介绍本项目采用的 Clean Architecture（整洁架构）设计，包括四层架构的职责、依赖关系和实现示例。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Clean Architecture](#clean-architecture)
  - [📋 目录](#-目录)
  - [1. 📚 架构概述](#1--架构概述)
  - [2. 🏗️ 架构层次](#2-️-架构层次)
    - [2.1 Domain Layer (领域层)](#21-domain-layer-领域层)
    - [2.2 Application Layer (应用层)](#22-application-layer-应用层)
    - [2.3 Infrastructure Layer (基础设施层)](#23-infrastructure-layer-基础设施层)
    - [2.4 Interfaces Layer (接口层)](#24-interfaces-layer-接口层)
  - [3. 🔄 依赖方向](#3--依赖方向)
  - [4. ✅ 优势](#4--优势)
  - [5. 📚 扩展阅读](#5--扩展阅读)

---

## 1. 📚 架构概述

本项目采用 **Clean Architecture（整洁架构）** 设计，将系统分为四个层次，确保业务逻辑与技术实现分离，提高代码的可测试性、可维护性和可扩展性。

### 核心原则

1. **依赖倒置**: 内层不依赖外层，外层依赖内层接口
2. **关注点分离**: 每层有明确的职责
3. **独立性**: 业务逻辑不依赖框架
4. **可测试性**: 每层都可以独立测试

---

## 2. 🏗️ 架构层次

### 2.1 Domain Layer (领域层)

**位置**: `internal/domain/`

**职责**:

- 核心业务逻辑
- 领域实体和值对象
- 领域服务接口
- 仓储接口

**规则**:

- 不依赖任何外部框架
- 不依赖 Infrastructure 或 Interfaces 层
- 只包含业务逻辑

**示例**:

```go
// internal/domain/user/entity.go
type User struct {
    ID    string
    Email string
    Name  string
}
```

### 2.2 Application Layer (应用层)

**位置**: `internal/application/`

**职责**:

- 用例编排
- 协调领域对象
- 应用服务
- DTO（数据传输对象）
- 工作流定义（Temporal Workflows）
- 活动定义（Temporal Activities）

**规则**:

- 只能导入 Domain 层
- 不依赖 Infrastructure 或 Interfaces 层

**示例**:

```go
// internal/application/user/service.go
type Service struct {
    repo domain.UserRepository
}

// internal/application/workflow/user_workflow.go
func UserWorkflow(ctx workflow.Context, input UserWorkflowInput) (UserWorkflowOutput, error) {
    // 工作流编排逻辑
}
```

### 2.3 Infrastructure Layer (基础设施层)

**位置**: `internal/infrastructure/`

**职责**:

- 技术实现细节
- 数据库访问
- 消息队列
- 外部服务集成
- 可观测性（OTLP）
- 工作流引擎（Temporal）

**规则**:

- 实现 Domain 层定义的接口
- 可以依赖外部库

**示例**:

```go
// internal/infrastructure/database/ent/repository/user_repository.go
type UserRepository struct {
    client *ent.Client
}

// internal/infrastructure/workflow/temporal/client.go
type Client struct {
    client client.Client
}
```

### 2.4 Interfaces Layer (接口层)

**位置**: `internal/interfaces/`

**职责**:

- 外部接口适配
- HTTP 处理器
- gRPC 服务
- GraphQL 解析器
- MQTT 处理器
- 工作流处理器（Temporal）

**规则**:

- 调用 Application 层
- 处理请求/响应转换

**示例**:

```go
// internal/interfaces/http/chi/handlers/user_handler.go
type UserHandler struct {
    service *application.UserService
}

// internal/interfaces/workflow/temporal/handler.go
type Handler struct {
    client client.Client
}
```

---

## 3. 🔄 依赖方向

依赖关系遵循以下规则：

```text
Interfaces → Application → Domain
     ↓            ↓
Infrastructure → Domain
```

### 依赖规则

1. **Domain Layer**: 不依赖任何其他层
2. **Application Layer**: 只能依赖 Domain Layer
3. **Infrastructure Layer**: 实现 Domain Layer 定义的接口
4. **Interfaces Layer**: 调用 Application Layer，处理外部请求

---

## 4. ✅ 优势

### 4.1 独立性

业务逻辑不依赖框架，可以轻松替换技术实现。

### 4.2 可测试性

每层都可以独立测试，通过 Mock 接口实现单元测试。

### 4.3 可维护性

清晰的职责分离，代码结构清晰，易于理解和维护。

### 4.4 可扩展性

易于添加新功能，新功能只需在相应层添加代码。

---

## 5. 📚 扩展阅读

- [领域模型设计](./domain-model.md) - 领域层设计详解
- [工作流架构设计](./workflow.md) - 工作流集成架构
- [Go 项目结构规范](../practices/engineering/06-项目结构规范.md)
- [Clean Architecture 原理解析](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

> 📚 **简介**
> 本文深入探讨 Clean Architecture 在本项目中的应用，系统讲解四层架构的职责、依赖关系和实现示例。通过本文，您将全面掌握项目的架构设计原则和实践方法。
