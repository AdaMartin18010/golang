# API 规范与代码生成

> **版本**: v1.0
> **日期**: 2025-01-XX
> **状态**: 持续完善中

---

## 📋 概述

框架支持 OpenAPI 和 AsyncAPI 规范，提供完整的代码生成、文档生成和验证工具。

## 🎯 支持的规范

### 1. OpenAPI 3.1.0

**位置**: `api/openapi/openapi.yaml`

**功能**:
- RESTful API 规范定义
- 代码生成（服务器、客户端、类型）
- 文档生成（HTML、Markdown）
- 规范验证
- Swagger UI 集成

### 2. AsyncAPI 3.0.0

**位置**: `api/asyncapi/asyncapi.yaml`

**功能**:
- 异步消息 API 规范定义
- 事件驱动架构文档
- 多协议支持（Kafka、MQTT、NATS）
- 代码生成
- 文档生成

---

## 🚀 快速开始

### 安装工具

```bash
# 安装 oapi-codegen（OpenAPI 代码生成）
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

# Docker 用于运行 OpenAPI Generator 和 AsyncAPI Generator
# 确保 Docker 已安装并运行
```

### 生成代码

```bash
# 生成 OpenAPI 代码
make generate-openapi

# 生成 AsyncAPI 代码
make generate-asyncapi

# 生成所有代码（包括 Ent、Wire、OpenAPI）
make generate
```

### 验证规范

```bash
# 验证 OpenAPI 规范
make validate-openapi

# 验证 AsyncAPI 规范
make validate-asyncapi

# 验证所有 API 规范
make validate-api
```

### 生成文档

```bash
# 生成 API 文档（HTML）
make generate-api-docs
```

---

## 📚 详细说明

### OpenAPI 代码生成

**生成服务器代码**:
```bash
oapi-codegen \
  -generate types,server,chi-server,spec \
  -package openapi \
  -o internal/interfaces/http/openapi/server.gen.go \
  api/openapi/openapi.yaml
```

**生成客户端代码**:
```bash
oapi-codegen \
  -generate types,client \
  -package client \
  -o pkg/api/client/client.gen.go \
  api/openapi/openapi.yaml
```

### AsyncAPI 代码生成

**使用 Docker**:
```bash
docker run --rm \
  -v ${PWD}:/local \
  asyncapi/generator-cli:latest \
  generate -g go \
  -i /local/api/asyncapi/asyncapi.yaml \
  -o /local/pkg/api/async
```

### Swagger UI 集成

框架提供了 Swagger UI 集成，可以在 HTTP 服务器中启用：

```go
import "github.com/yourusername/golang/internal/interfaces/http/openapi"

// 配置 Swagger UI
swaggerConfig := openapi.DefaultConfig()
swaggerConfig.OpenAPISpecPath = "api/openapi/openapi.yaml"
swaggerConfig.Title = "My API Documentation"

// 注册路由
router.Mount("/swagger", openapi.Handler(swaggerConfig))
```

访问 `http://localhost:8080/swagger/` 查看 API 文档。

---

## 🔧 配置说明

### OpenAPI 配置

**文件**: `api/openapi/openapi.yaml`

**关键配置**:
- `openapi: 3.1.0` - 规范版本
- `info` - API 信息（标题、版本、描述）
- `servers` - 服务器地址
- `paths` - API 路径定义
- `components` - 可复用组件（schemas、responses 等）

### AsyncAPI 配置

**文件**: `api/asyncapi/asyncapi.yaml`

**关键配置**:
- `asyncapi: 3.0.0` - 规范版本
- `info` - API 信息
- `servers` - 消息服务器配置
- `channels` - 消息通道定义
- `components` - 可复用组件（messages、schemas 等）

---

## 📝 最佳实践

1. **规范优先**: 先定义 API 规范，再生成代码
2. **版本控制**: 规范文件纳入版本控制
3. **持续验证**: 在 CI/CD 中验证规范
4. **文档同步**: 规范变更时更新文档
5. **代码生成**: 使用生成的代码，避免手动编写

---

## 🔗 相关文档

- [OpenAPI 深度解析](../../architecture/tech-stack/api/openapi.md)
- [AsyncAPI 深度解析](../../architecture/tech-stack/api/asyncapi.md)
- [API 定义目录](../../../api/README.md)

---

**最后更新**: 2025-01-XX
