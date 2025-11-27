# 用户领域示例

> **注意**: 这是框架使用示例，展示如何使用框架构建领域模型和应用服务。

## 📋 说明

本示例展示了如何使用框架构建一个完整的用户管理功能，包括：

1. **领域层** (`domain/`) - 用户实体、仓储接口、领域错误
2. **应用层** (`application/`) - 用户应用服务、DTO
3. **接口层** (`interfaces/`) - HTTP Handler、gRPC Handler

## 🎯 示例目的

- 展示 Clean Architecture 在框架中的应用
- 展示领域模型设计
- 展示应用服务实现
- 展示接口层实现

## 📁 目录结构

```text
examples/framework-usage/user-domain/
├── domain/          # 领域层示例
│   ├── entity.go    # 用户实体
│   ├── repository.go # 仓储接口
│   └── errors.go    # 领域错误
├── application/     # 应用层示例
│   ├── service.go   # 应用服务
│   └── errors.go    # 应用错误
├── interfaces/      # 接口层示例
│   ├── http/        # HTTP Handler
│   └── grpc/        # gRPC Handler
└── README.md        # 本文件
```

## 🚀 使用方式

这些代码仅作为示例，展示框架的使用方式。在实际项目中，用户应该：

1. 定义自己的领域模型（通过 Ent Schema）
2. 实现自己的应用服务
3. 实现自己的接口层

## 📚 相关文档

- [框架快速开始指南](../../../docs/framework/05-快速开始指南.md) ⭐⭐⭐⭐⭐
- [框架使用示例说明](../../../docs/framework/02-框架使用示例说明.md)
- [框架最佳实践指南](../../../docs/framework/06-最佳实践指南.md)
- [架构说明](../../../docs/architecture/00-架构模型与依赖注入完整说明.md)

## 🔗 快速链接

- [运行示例](../README.md#2-运行示例)
- [框架基础设施说明](../../../docs/framework/00-框架基础设施说明.md)
- [API规范与代码生成](../../../docs/framework/01-API规范与代码生成.md)
