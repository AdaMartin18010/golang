# 项目架构文档

> **简介**: 本项目架构相关文档索引，包括 Clean Architecture、领域模型、工作流等架构设计文档。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [项目架构文档](#项目架构文档)
  - [📋 目录](#-目录)
  - [📚 文档列表](#-文档列表)
    - [核心架构文档](#核心架构文档)
    - [导航文档](#导航文档)
  - [🎯 快速导航](#-快速导航)
    - [按主题导航](#按主题导航)
    - [按层次导航](#按层次导航)
    - [按技术导航](#按技术导航)
  - [📊 知识图谱](#-知识图谱)
  - [🔍 概念定义](#-概念定义)
  - [📖 对比矩阵](#-对比矩阵)
  - [🚀 快速开始](#-快速开始)
    - [1. 了解架构设计](#1-了解架构设计)
    - [2. 理解领域模型](#2-理解领域模型)
    - [3. 学习工作流](#3-学习工作流)
    - [4. 查看知识图谱](#4-查看知识图谱)
  - [📚 相关文档](#-相关文档)
    - [架构相关](#架构相关)
    - [使用指南](#使用指南)
    - [项目文档](#项目文档)

---

## 📚 文档列表

### 核心架构文档

1. **[Clean Architecture](./clean-architecture.md)** ⭐⭐⭐⭐⭐ ✅ 已增强
   - 四层架构设计
   - 依赖关系和依赖倒置原则
   - 设计决策论证（为什么选择 Clean Architecture、为什么选择四层架构）
   - 优势分析与量化收益（独立性、可测试性、可维护性、可扩展性）
   - 实际应用示例和代码说明

2. **[领域模型设计](./domain-model.md)** ⭐⭐⭐⭐ ✅ 已增强
   - 用户领域模型（实体、仓储接口、领域服务、领域错误）
   - DDD 核心概念说明（实体、值对象、聚合、仓储）
   - 设计原则论证（实体独立性、接口定义、业务规则封装、错误处理）
   - 实际应用示例和代码说明

3. **[工作流架构设计](./workflow.md)** ⭐⭐⭐⭐⭐ ✅ 已增强
   - Temporal 工作流集成（Server, Worker, Client）
   - 选型论证（为什么选择 Temporal）
   - 技术选型对比（Temporal vs Airflow vs Conductor vs Cadence）
   - 工作流模式（顺序执行、并行执行、条件执行、循环执行）
   - 应用场景分析（适用场景和不适用场景）
   - 最佳实践

### 导航文档

4. **[知识图谱](./00-知识图谱.md)** ⭐⭐⭐⭐ ✅ 已增强
   - 架构全景思维导图（mindmap 格式）
   - 四层架构关系图和详细说明
   - 组件交互关系图（包含端口、协议等详细信息）
   - 数据流向图（包含完整的请求处理流程）
   - 技术栈多维分析（技术选型决策树、技术栈关系网络）
   - 组件关系深度解析（依赖关系矩阵、交互时序图、设计模式应用）
   - 架构设计决策论证（为什么选择 Clean Architecture、为什么选择四层架构、为什么采用依赖倒置）

5. **[技术对比矩阵](./00-对比矩阵.md)** ⭐⭐⭐ ✅ 已增强
   - Web 框架多维对比（功能特性、性能基准、学习成本、选型决策论证）
   - ORM 多维对比（类型安全、性能、开发体验、选型决策论证）
   - 工作流引擎多维对比（功能特性、Go 支持、可观测性、选型决策论证）
   - 可观测性多维对比（功能覆盖、标准兼容、集成复杂度、选型决策论证）
   - 消息队列多维对比（协议特性、性能特性、使用场景、选型决策论证）
   - 数据库多维对比（功能特性、性能特性、生态系统、选型决策论证）

6. **[概念定义体系](./00-概念定义体系.md)** ⭐⭐⭐⭐ ✅ 已增强
   - 架构概念深度解析（Clean Architecture, DDD, 依赖倒置原则, 关注点分离）
   - 层次概念深度解析（Domain, Application, Infrastructure, Interfaces）
   - 技术概念深度解析（Temporal Workflow, Temporal Activity, Ent ORM, OpenTelemetry）
   - 模式概念深度解析（Repository Pattern, Service Pattern, DTO Pattern, Dependency Injection）
   - 设计原则深度解析（SOLID 原则, DRY 原则, KISS 原则, YAGNI 原则）

7. **[技术栈文档索引](./tech-stack/README.md)** ⭐⭐⭐⭐ ✅ 已拆分
   - 全面梳理和论证所有第三方技术栈
   - Chi Router、Ent ORM、Temporal、OpenTelemetry 等深度解析
   - 每个技术的核心特性、选型论证、实际应用、最佳实践
   - 技术栈集成架构和最佳实践

8. **[Go 1.25.3 技术栈对齐](./00-Go-1.25.3技术栈对齐.md)** ⭐⭐⭐ ✅ 新建
   - Go 1.25.3 新特性说明和应用
   - 技术版本对齐矩阵和策略
   - 新特性在项目中的应用场景
   - 技术栈演进路径和规划

---

## 🎯 快速导航

### 按主题导航

- **架构设计**: [Clean Architecture](./clean-architecture.md)
- **领域模型**: [领域模型设计](./domain-model.md)
- **工作流**: [工作流架构设计](./workflow.md)

### 按层次导航

- **领域层**: [领域模型设计](./domain-model.md) → Domain Layer
- **应用层**: [Clean Architecture](./clean-architecture.md) → Application Layer
- **基础设施层**: [Clean Architecture](./clean-architecture.md) → Infrastructure Layer
- **接口层**: [Clean Architecture](./clean-architecture.md) → Interfaces Layer

### 按技术导航

- **Temporal**: [工作流架构设计](./workflow.md) | [技术栈文档](./tech-stack/workflow/temporal.md)
- **Ent ORM**: [Clean Architecture](./clean-architecture.md) → Infrastructure Layer | [技术栈文档](./tech-stack/data/ent-orm.md)
- **Chi Router**: [Clean Architecture](./clean-architecture.md) → Interfaces Layer | [技术栈文档](./tech-stack/web/chi-router.md)
- **OpenTelemetry**: [技术栈文档](./tech-stack/observability/opentelemetry.md)
- **Prometheus**: [技术栈文档](./tech-stack/observability/prometheus.md)
- **Grafana**: [技术栈文档](./tech-stack/observability/grafana.md)
- **Jaeger**: [技术栈文档](./tech-stack/observability/jaeger.md)
- **PostgreSQL**: [技术栈文档](./tech-stack/data/postgresql-pgx.md)
- **Kafka**: [技术栈文档](./tech-stack/messaging/kafka.md)
- **MQTT**: [技术栈文档](./tech-stack/messaging/mqtt.md)
- **Viper**: [技术栈文档](./tech-stack/config/viper.md)
- **Slog**: [技术栈文档](./tech-stack/config/slog.md)
- **Wire**: [技术栈文档](./tech-stack/config/wire.md)
- **gRPC**: [技术栈文档](./tech-stack/api/grpc.md)
- **GraphQL**: [技术栈文档](./tech-stack/api/graphql.md)
- **Go 1.25.3**: [Go 1.25.3 技术栈对齐](./00-Go-1.25.3技术栈对齐.md)

---

## 📊 知识图谱

查看完整的架构知识图谱：

👉 **[架构知识图谱](./00-知识图谱.md)** ✅ 已增强

**内容**:

- 架构全景思维导图（mindmap 格式）
- 四层架构关系图和详细说明
- 组件交互关系图（包含端口、协议等详细信息）
- 数据流向图（包含完整的请求处理流程）
- 技术栈多维分析（技术选型决策树、技术栈关系网络）
- 组件关系深度解析（依赖关系矩阵、交互时序图、设计模式应用）
- 架构设计决策论证（为什么选择 Clean Architecture、为什么选择四层架构、为什么采用依赖倒置）

**增强特点**:

- ✅ 详细的思维导图（mindmap）
- ✅ 多维分析和详细说明
- ✅ 设计决策论证
- ✅ 性能分析和优化策略

---

## 🔍 概念定义

查看完整的架构概念定义：

👉 **[概念定义体系](./00-概念定义体系.md)** ✅ 已增强

**内容**:

- 架构概念深度解析（Clean Architecture, DDD, 依赖倒置原则, 关注点分离）
- 层次概念深度解析（Domain, Application, Infrastructure, Interfaces）
- 技术概念深度解析（Temporal Workflow, Temporal Activity, Ent ORM, OpenTelemetry）
- 模式概念深度解析（Repository Pattern, Service Pattern, DTO Pattern, Dependency Injection）
- 设计原则深度解析（SOLID 原则, DRY 原则, KISS 原则, YAGNI 原则）

**增强特点**:

- ✅ 详细的概念解释（定义、核心思想、特征）
- ✅ 设计原理论证（为什么需要、如何应用、优势）
- ✅ 实际应用示例和代码说明
- ✅ 在本项目中的应用说明

---

## 📖 对比矩阵

查看技术选型对比：

👉 **[技术对比矩阵](./00-对比矩阵.md)** ✅ 已增强

**内容**:

- Web 框架多维对比（功能特性、性能基准、学习成本、选型决策论证）
- ORM 多维对比（类型安全、性能、开发体验、选型决策论证）
- 工作流引擎多维对比（功能特性、Go 支持、可观测性、选型决策论证）
- 可观测性多维对比（功能覆盖、标准兼容、集成复杂度、选型决策论证）
- 消息队列多维对比（协议特性、性能特性、使用场景、选型决策论证）
- 数据库多维对比（功能特性、性能特性、生态系统、选型决策论证）

**增强特点**:

- ✅ 多维对比分析（功能、性能、学习成本、使用场景等）
- ✅ 性能基准测试数据
- ✅ 选型决策矩阵和详细论证
- ✅ 量化评估指标和权重分析

---

## 🚀 快速开始

### 1. 了解架构设计

从 [Clean Architecture](./clean-architecture.md) 开始，了解项目的整体架构设计。

### 2. 理解领域模型

阅读 [领域模型设计](./domain-model.md)，理解业务领域的建模方法。

### 3. 学习工作流

查看 [工作流架构设计](./workflow.md)，了解 Temporal 工作流的集成和使用。

### 4. 查看知识图谱

浏览 [知识图谱](./00-知识图谱.md)，获得架构的全局视图。

---

## 📚 相关文档

### 架构相关

- [技术栈文档索引](./tech-stack/README.md) - 第三方技术栈文档索引（按分类组织）
- [技术栈概览](./tech-stack/00-技术栈概览.md) - 技术栈分层和选型原则
- [技术栈思维导图](./tech-stack/00-技术栈思维导图.md) - 技术栈可视化（思维导图、关系网络图、决策流程图）
- [技术栈集成](./tech-stack/01-技术栈集成.md) - 技术栈集成架构和最佳实践
- [技术栈选型决策树](./tech-stack/02-技术栈选型决策树.md) - 技术栈选型决策树
- [第三方技术栈深度解析](./00-第三方技术栈深度解析.md) - 第三方技术栈全面解析（旧版，已拆分）
- [Go 1.25.3 技术栈对齐](./00-Go-1.25.3技术栈对齐.md) - Go 1.25.3 新特性应用
- [知识图谱](./00-知识图谱.md) - 架构知识图谱
- [概念定义体系](./00-概念定义体系.md) - 概念定义体系
- [技术对比矩阵](./00-对比矩阵.md) - 技术选型对比

### 使用指南

- [工作流使用指南](../guides/workflow.md) - 工作流使用指南
- [开发指南](../guides/development.md) - 开发指南
- [部署指南](../guides/deployment.md) - 部署指南

### 项目文档

- [项目文档索引](../00-项目文档索引.md) - 完整文档索引
- [文档结构规范](../00-项目文档结构规范.md) - 文档结构规范

---

> 📚 **简介**
> 本文档提供了项目架构文档的完整索引，帮助快速定位所需的架构设计文档。通过本文档，您可以系统地了解项目的架构设计。
