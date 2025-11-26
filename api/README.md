# API Definitions

API 定义目录，包含所有 API 的规范定义，用于代码生成、文档生成和 API 网关配置。

## 📋 目录结构

```text
api/
├── openapi/           # OpenAPI 3.1.0 规范
│   └── openapi.yaml   # RESTful API 规范定义
├── asyncapi/          # AsyncAPI 2.6.0 规范
│   └── asyncapi.yaml  # 异步消息 API 规范定义
├── graphql/           # GraphQL Schema
│   └── schema.graphql # GraphQL Schema 定义
└── README.md          # 本文件
```

## 📚 规范说明

### OpenAPI 3.1.0

**文件**: `openapi/openapi.yaml`

**用途**:

- RESTful API 规范定义
- API 文档自动生成
- 客户端代码生成
- API 网关配置
- API 测试用例生成

**特性**:

- ✅ 完整的请求/响应定义
- ✅ 错误响应标准化
- ✅ 分页和过滤支持
- ✅ 安全认证定义（BearerAuth）
- ✅ 健康检查端点
- ✅ 示例和描述

**相关文档**: [`docs/architecture/tech-stack/api/openapi.md`](../../docs/architecture/tech-stack/api/openapi.md)

**代码生成示例**:

```bash
# 使用 openapi-generator 生成 Go 客户端
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
  -i /local/api/openapi/openapi.yaml \
  -g go \
  -o /local/pkg/api/client

# 使用 oapi-codegen 生成 Go 服务器代码
oapi-codegen -generate types,server -package api api/openapi/openapi.yaml > pkg/api/server.go
```

### AsyncAPI 2.6.0

**文件**: `asyncapi/asyncapi.yaml`

**用途**:

- 异步消息 API 规范定义
- 事件驱动架构文档
- 消息队列配置
- 客户端代码生成

**特性**:

- ✅ 完整的事件定义
- ✅ 多协议绑定（Kafka、MQTT、NATS）
- ✅ 完整的 Schema 定义
- ✅ 分层通道命名
- ✅ 示例和描述

**相关文档**: [`docs/architecture/tech-stack/api/asyncapi.md`](../../docs/architecture/tech-stack/api/asyncapi.md)

**代码生成示例**:

```bash
# 使用 asyncapi-generator 生成 Go 客户端
docker run --rm -v ${PWD}:/local asyncapi/generator-cli \
  generate -g go -i /local/api/asyncapi/asyncapi.yaml \
  -o /local/pkg/api/async
```

### GraphQL Schema

**文件**: `graphql/schema.graphql`

**用途**:

- GraphQL API Schema 定义
- 类型系统定义
- 查询和变更定义
- 代码生成

**特性**:

- ✅ 类型定义
- ✅ 查询和变更
- ✅ 订阅支持
- ✅ 指令和标量类型

**相关文档**: [`docs/architecture/tech-stack/api/graphql.md`](../../docs/architecture/tech-stack/api/graphql.md)

**代码生成示例**:

```bash
# 使用 gqlgen 生成 Go 代码
go run github.com/99designs/gqlgen generate
```

## 🛠️ 工具和集成

### 验证规范

```bash
# 验证 OpenAPI 规范
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli validate \
  -i /local/api/openapi/openapi.yaml

# 验证 AsyncAPI 规范
docker run --rm -v ${PWD}:/local asyncapi/generator-cli \
  validate -i /local/api/asyncapi/asyncapi.yaml
```

### 文档生成

```bash
# 生成 OpenAPI 文档（HTML）
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate \
  -i /local/api/openapi/openapi.yaml \
  -g html \
  -o /local/docs/api/openapi

# 生成 AsyncAPI 文档（HTML）
docker run --rm -v ${PWD}:/local asyncapi/generator-cli \
  generate -g html -i /local/api/asyncapi/asyncapi.yaml \
  -o /local/docs/api/asyncapi
```

### CI/CD 集成

```yaml
# .github/workflows/api-validation.yml
name: API Validation

on:
  pull_request:
    paths:
      - 'api/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate OpenAPI
        run: |
          docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli \
            validate -i /local/api/openapi/openapi.yaml
      - name: Validate AsyncAPI
        run: |
          docker run --rm -v ${PWD}:/local asyncapi/generator-cli \
            validate -i /local/api/asyncapi/asyncapi.yaml
```

## 📖 最佳实践

### OpenAPI 最佳实践

1. **版本控制**: 使用语义化版本（如 `v1.0.0`）
2. **错误处理**: 定义标准的错误响应格式
3. **分页**: 使用标准的分页参数（`page`, `limit`, `offset`）
4. **安全**: 明确定义认证和授权机制
5. **示例**: 为每个端点和 Schema 提供示例

### AsyncAPI 最佳实践

1. **通道命名**: 使用分层命名（如 `user.created`, `order.paid`）
2. **Schema 定义**: 为每个消息定义完整的 Schema
3. **协议绑定**: 支持多种消息协议（Kafka、MQTT、NATS）
4. **版本控制**: 在消息头中包含版本信息
5. **示例**: 为每个消息提供示例

### GraphQL 最佳实践

1. **类型系统**: 使用强类型系统
2. **查询优化**: 避免 N+1 查询问题
3. **安全性**: 实现查询深度限制和复杂度限制
4. **版本控制**: 使用字段弃用而不是删除
5. **文档**: 为每个类型和字段提供描述

## 🔗 相关资源

- [OpenAPI 规范](https://spec.openapis.org/oas/v3.1.0)
- [AsyncAPI 规范](https://www.asyncapi.com/docs/specifications/v2.6.0)
- [GraphQL 规范](https://graphql.org/learn/)
- [API 设计最佳实践](../../docs/architecture/tech-stack/api/)

## 📝 更新日志

- **2025-11-11**: 增强 OpenAPI 和 AsyncAPI 规范，添加完整的 Schema 定义、错误处理和示例
- **2025-10-29**: 初始版本，包含基础 API 规范定义
