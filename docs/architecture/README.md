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
  - [🎯 快速导航](#-快速导航)
  - [📊 知识图谱](#-知识图谱)
  - [🔍 概念定义](#-概念定义)
  - [📖 对比矩阵](#-对比矩阵)

---

## 📚 文档列表

### 核心架构文档

1. **[Clean Architecture](./clean-architecture.md)** ⭐⭐⭐⭐⭐
   - 四层架构设计
   - 依赖关系
   - 实现示例

2. **[领域模型设计](./domain-model.md)** ⭐⭐⭐⭐
   - 用户领域模型
   - 实体和接口定义
   - 设计原则

3. **[工作流架构设计](./workflow.md)** ⭐⭐⭐⭐⭐
   - Temporal 工作流集成
   - 工作流模式
   - 最佳实践

### 导航文档

4. **[知识图谱](./00-知识图谱.md)** ⭐⭐⭐⭐
   - 架构全景图
   - 技术栈关系
   - 组件关系图

5. **[技术对比矩阵](./00-对比矩阵.md)** ⭐⭐⭐
   - Web 框架对比
   - ORM 对比
   - 工作流引擎对比

6. **[概念定义体系](./00-概念定义体系.md)** ⭐⭐⭐⭐
   - 架构概念定义
   - 层次概念定义
   - 技术概念定义

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

- **Temporal**: [工作流架构设计](./workflow.md)
- **Ent ORM**: [Clean Architecture](./clean-architecture.md) → Infrastructure Layer
- **Chi Router**: [Clean Architecture](./clean-architecture.md) → Interfaces Layer

---

## 📊 知识图谱

查看完整的架构知识图谱：

👉 **[架构知识图谱](./00-知识图谱.md)**

包括：
- 架构全景图
- 分层架构图
- 技术栈关系图
- 数据流程图
- 组件关系图

---

## 🔍 概念定义

查看完整的架构概念定义：

👉 **[概念定义体系](./00-概念定义体系.md)**

包括：
- 架构概念（Clean Architecture, DDD）
- 层次概念（Domain, Application, Infrastructure, Interfaces）
- 技术概念（Temporal, Ent, OpenTelemetry）
- 模式概念（Repository, Service, DTO）

---

## 📖 对比矩阵

查看技术选型对比：

👉 **[技术对比矩阵](./00-对比矩阵.md)**

包括：
- Web 框架对比（Chi vs Gin vs Echo）
- ORM 对比（Ent vs GORM）
- 工作流引擎对比（Temporal vs Airflow）
- 可观测性对比（OpenTelemetry vs Prometheus）
- 消息队列对比（Kafka vs MQTT）
- 数据库对比（PostgreSQL vs MySQL）

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

- [工作流使用指南](../guides/workflow.md) - 工作流使用指南
- [开发指南](../guides/development.md) - 开发指南
- [部署指南](../guides/deployment.md) - 部署指南

---

> 📚 **简介**
> 本文档提供了项目架构文档的完整索引，帮助快速定位所需的架构设计文档。通过本文档，您可以系统地了解项目的架构设计。
