# 框架基础设施文档

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 概述

本目录包含框架基础设施的详细文档，说明框架的核心组件、使用方式和最佳实践。

## 📚 文档索引

### 核心文档

1. **[框架基础设施说明](00-框架基础设施说明.md)**
   - 框架核心组件概述
   - 开发工具和流程
   - 快速开始指南

2. **[API规范与代码生成](01-API规范与代码生成.md)**
   - OpenAPI 和 AsyncAPI 支持
   - 代码生成工具
   - 文档生成和验证

3. **[框架使用示例说明](02-框架使用示例说明.md)**
   - 示例代码说明
   - 使用方式
   - 最佳实践

## 🎯 快速导航

### 开发工具

- **日志系统**: `internal/framework/logger/`
- **热重载**: `.air.toml`
- **测试覆盖**: `scripts/test-coverage.sh`
- **API 生成**: `scripts/api/`

### 配置文件

- **Air 配置**: `.air.toml`
- **Codecov 配置**: `codecov.yml`
- **OpenAPI 规范**: `api/openapi/openapi.yaml`
- **AsyncAPI 规范**: `api/asyncapi/asyncapi.yaml`

### Makefile 命令

```bash
make help              # 显示所有命令
make run-dev           # 开发模式运行
make test-coverage     # 测试覆盖率
make generate-openapi  # 生成 OpenAPI 代码
make validate-api      # 验证 API 规范
```

## 📖 相关文档

- [项目定位与任务梳理](../00-项目定位与任务梳理.md)
- [架构说明](../architecture/00-架构模型与依赖注入完整说明.md)
- [开发指南](../guides/development.md)

---

**最后更新**: 2025-01-XX
