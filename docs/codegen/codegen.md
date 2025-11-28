# 代码生成工具链文档

> **版本**: v1.0
> **日期**: 2025-01-XX
> **位置**: `scripts/`

---

## 📋 目录

- [1. 概述](#1-概述)
- [2. gRPC 代码生成](#2-grpc-代码生成)
- [3. OpenAPI 代码生成](#3-openapi-代码生成)
- [4. AsyncAPI 代码生成](#4-asyncapi-代码生成)
- [5. GraphQL 代码生成](#5-graphql-代码生成)
- [6. 使用说明](#6-使用说明)

---

## 1. 概述

项目提供了完整的代码生成工具链，支持从规范文件自动生成代码。

### 支持的工具

- ✅ **gRPC**: Protocol Buffers → Go 代码
- ✅ **OpenAPI**: OpenAPI 3.1.0 → Go 服务器/客户端代码
- ✅ **AsyncAPI**: AsyncAPI 3.0.0 → Go 代码（可选）
- ✅ **GraphQL**: GraphQL Schema → Go Resolver 代码（可选）

---

## 2. gRPC 代码生成

### 2.1 安装工具

```bash
# 安装 Protocol Buffers 编译器
# macOS
brew install protobuf

# Linux
sudo apt-get install protobuf-compiler

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2.2 生成代码

```bash
# 使用 Makefile
make generate-grpc

# 或直接运行脚本
bash scripts/grpc/generate.sh
```

### 2.3 输出文件

生成的文件位于 `internal/interfaces/grpc/proto/`:
- `user.pb.go` - 消息类型
- `user_grpc.pb.go` - 服务接口

---

## 3. OpenAPI 代码生成

### 3.1 安装工具

```bash
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

### 3.2 生成代码

```bash
# 使用 Makefile
make generate-openapi

# 或直接运行脚本
bash scripts/api/generate-openapi.sh
```

### 3.3 输出文件

- `internal/interfaces/http/openapi/server.gen.go` - 服务器代码
- `pkg/api/client/client.gen.go` - 客户端代码

---

## 4. AsyncAPI 代码生成

### 4.1 使用 Docker

```bash
# 使用 Makefile
make generate-asyncapi

# 或直接运行脚本
bash scripts/api/generate-asyncapi.sh
```

### 4.2 输出文件

- `pkg/api/async/` - 生成的异步 API 代码

---

## 5. GraphQL 代码生成（可选）

### 5.1 安装工具

```bash
go install github.com/99designs/gqlgen@latest
```

### 5.2 初始化配置

```bash
cd internal/interfaces/graphql
gqlgen init
```

### 5.3 生成代码

```bash
gqlgen generate
```

---

## 6. 使用说明

### 6.1 一键生成所有代码

```bash
make generate
```

这会生成：
- Ent 代码
- Wire 代码
- gRPC 代码
- OpenAPI 代码

### 6.2 验证规范

```bash
# 验证 OpenAPI
make validate-openapi

# 验证 AsyncAPI
make validate-asyncapi

# 验证所有 API
make validate-api
```

---

## 📚 相关资源

- [gRPC 文档](../grpc/grpc.md)
- [OpenAPI 规范](../../api/openapi/openapi.yaml)
- [AsyncAPI 规范](../../api/asyncapi/asyncapi.yaml)

---

**最后更新**: 2025-01-XX
