# 第三方技术栈文档索引

> **简介**: 本项目第三方技术栈文档索引，按技术栈分类组织，便于查找和维护。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [第三方技术栈文档索引](#第三方技术栈文档索引)
  - [📋 目录](#-目录)
  - [📚 文档概览](#-文档概览)
    - [通用文档](#通用文档)
  - [🌐 Web 框架层](#-web-框架层)
  - [🗄️ 数据访问层](#️-数据访问层)
  - [🔄 工作流层](#-工作流层)
  - [📊 可观测性层](#-可观测性层)
  - [💬 消息队列层](#-消息队列层)
  - [⚙️ 配置和工具层](#️-配置和工具层)
  - [🔌 API 协议层](#-api-协议层)
  - [🎯 快速导航](#-快速导航)
    - [按技术栈类型](#按技术栈类型)
    - [按文档类型](#按文档类型)
  - [📚 相关文档](#-相关文档)

---

## 📚 文档概览

### 通用文档

- **[技术栈概览](./00-技术栈概览.md)** - 技术栈分层和选型原则
- **[技术栈集成](./01-技术栈集成.md)** - 技术栈集成架构和最佳实践
- **[技术栈选型决策树](./02-技术栈选型决策树.md)** - 技术栈选型决策树

---

## 🌐 Web 框架层

- **[Chi Router](./web/chi-router.md)** - Chi Router 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

---

## 🗄️ 数据访问层

- **[Ent ORM](./data/ent-orm.md)** - Ent ORM 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[PostgreSQL (pgx)](./data/postgresql-pgx.md)** - PostgreSQL (pgx) 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

---

## 🔄 工作流层

- **[Temporal](./workflow/temporal.md)** - Temporal 工作流深度解析
  - 核心特性、选型论证、实际应用、最佳实践

---

## 📊 可观测性层

- **[OpenTelemetry](./observability/opentelemetry.md)** - OpenTelemetry 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[Prometheus](./observability/prometheus.md)** - Prometheus 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[Grafana](./observability/grafana.md)** - Grafana 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[Jaeger](./observability/jaeger.md)** - Jaeger 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

---

## 💬 消息队列层

- **[Kafka](./messaging/kafka.md)** - Kafka 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[MQTT](./messaging/mqtt.md)** - MQTT 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

---

## ⚙️ 配置和工具层

- **[Viper](./config/viper.md)** - Viper 配置管理深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[Slog](./config/slog.md)** - Slog 日志库深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[Wire](./config/wire.md)** - Wire 依赖注入深度解析
  - 核心特性、选型论证、实际应用、最佳实践

---

## 🔌 API 协议层

- **[gRPC](./api/grpc.md)** - gRPC 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

- **[GraphQL](./api/graphql.md)** - GraphQL 深度解析
  - 核心特性、选型论证、实际应用、最佳实践

---

## 🎯 快速导航

### 按技术栈类型

- **Web 框架**: [Chi Router](./web/chi-router.md)
- **数据访问**: [Ent ORM](./data/ent-orm.md) | [PostgreSQL](./data/postgresql-pgx.md)
- **工作流**: [Temporal](./workflow/temporal.md)
- **可观测性**: [OpenTelemetry](./observability/opentelemetry.md) | [Prometheus](./observability/prometheus.md) | [Grafana](./observability/grafana.md) | [Jaeger](./observability/jaeger.md)
- **消息队列**: [Kafka](./messaging/kafka.md) | [MQTT](./messaging/mqtt.md)
- **配置工具**: [Viper](./config/viper.md) | [Slog](./config/slog.md) | [Wire](./config/wire.md)
- **API 协议**: [gRPC](./api/grpc.md) | [GraphQL](./api/graphql.md)

### 按文档类型

- **概览**: [技术栈概览](./00-技术栈概览.md)
- **集成**: [技术栈集成](./01-技术栈集成.md)
- **选型**: [技术栈选型决策树](./02-技术栈选型决策树.md)

---

## 📚 相关文档

- [架构文档索引](../README.md) - 架构文档索引
- [项目文档索引](../../00-项目文档索引.md) - 完整文档索引

---

> 📚 **简介**
> 本文档提供了第三方技术栈文档的完整索引，按技术栈分类组织，便于快速查找和维护。
