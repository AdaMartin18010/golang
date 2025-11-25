# 工作流使用指南

> **简介**: 本文档介绍如何使用本项目中的 Temporal 工作流功能，包括快速开始、API 使用和工作流定义。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [工作流使用指南](#工作流使用指南)
  - [📋 目录](#-目录)
  - [1. 📚 概述](#1--概述)
  - [2. 🚀 快速开始](#2--快速开始)
  - [3. 💻 使用工作流](#3--使用工作流)
  - [4. 🔧 工作流定义](#4--工作流定义)
  - [5. ⚙️ 配置](#5-️-配置)
  - [6. 🔍 监控和调试](#6--监控和调试)
  - [7. 📚 扩展阅读](#7--扩展阅读)

---

## 1. 📚 概述

本项目使用 [Temporal](https://temporal.io/) 作为工作流编排引擎，用于构建可靠的分布式应用。

### 1.1 组件

1. **Temporal Server** - 工作流引擎
2. **Worker** - 执行工作流和活动
3. **Client** - 启动和查询工作流

### 1.2 工作流类型

- **UserWorkflow** - 用户相关操作工作流
  - 创建用户
  - 更新用户
  - 删除用户

---

## 2. 🚀 快速开始

### 2.1 启动 Temporal Server

使用 Docker Compose：

```bash
cd deployments/docker
docker-compose up -d temporal temporal-db temporal-ui
```

### 2.2 启动 Worker

```bash
go run ./cmd/temporal-worker
```

### 2.3 启动应用

```bash
go run ./cmd/server
```

### 2.4 访问 Temporal UI

打开浏览器访问：http://localhost:8088

---

## 3. 💻 使用工作流

### 3.1 启动工作流

```bash
curl -X POST http://localhost:8080/api/v1/workflows/user \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "email": "test@example.com",
    "name": "Test User",
    "action": "create"
  }'
```

### 3.2 查询工作流结果

```bash
curl http://localhost:8080/api/v1/workflows/user/{workflow_id}?run_id={run_id}
```

---

## 4. 🔧 工作流定义

### 4.1 UserWorkflow

用户工作流处理用户相关的操作，包括：

1. **创建用户**
   - 验证用户信息
   - 创建用户
   - 发送通知

2. **更新用户**
   - 更新用户信息
   - 发送通知

3. **删除用户**
   - 删除用户
   - 发送通知

### 4.2 活动（Activities）

- `ValidateUserActivity` - 验证用户信息
- `CreateUserActivity` - 创建用户
- `UpdateUserActivity` - 更新用户
- `DeleteUserActivity` - 删除用户
- `SendNotificationActivity` - 发送通知

---

## 5. ⚙️ 配置

在 `configs/config.yaml` 中配置：

```yaml
workflow:
  temporal:
    address: "localhost:7233"
    task_queue: "user-task-queue"
    namespace: "default"
```

### 5.1 配置文件

在 `configs/config.yaml` 中配置：

```yaml
workflow:
  temporal:
    address: "localhost:7233"
    task_queue: "user-task-queue"
    namespace: "default"
```

### 5.2 环境变量

也可以通过环境变量配置：

```bash
export TEMPORAL_ADDRESS=localhost:7233
export TEMPORAL_TASK_QUEUE=user-task-queue
```

---

## 6. 🔍 监控和调试

### 6.1 Temporal UI

访问 http://localhost:8088 查看：
- 工作流执行历史
- 活动执行状态
- 工作流查询和信号

### 6.2 指标

Temporal 提供丰富的指标，可以集成到 Prometheus 和 Grafana。

### 6.3 故障排除

#### Worker 无法连接

- 检查 Temporal Server 是否运行
- 验证连接地址配置
- 检查网络连接

#### 工作流执行失败

- 查看 Temporal UI 中的错误信息
- 检查活动日志
- 验证输入参数

---

## 7. 📚 扩展阅读

- [Temporal 官方文档](https://docs.temporal.io/)
- [Temporal Go SDK](https://docs.temporal.io/dev-guide/go)
- [工作流模式](https://docs.temporal.io/workflows)
- [工作流架构设计](../architecture/workflow.md) - 架构设计详解

---

> 📚 **简介**
> 本文深入探讨 Temporal 工作流的使用方法，系统讲解快速开始、API 使用、工作流定义和最佳实践。通过本文，您将全面掌握工作流的使用方法。
