# gRPC 使用文档

> **版本**: v1.0
> **日期**: 2025-01-XX
> **位置**: `internal/interfaces/grpc/`

---

## 📋 目录

- [gRPC 使用文档](#grpc-使用文档)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
    - [特性](#特性)
  - [2. 快速开始](#2-快速开始)
    - [2.1 安装工具](#21-安装工具)
    - [2.2 生成代码](#22-生成代码)
  - [3. Proto 文件定义](#3-proto-文件定义)
    - [3.1 服务定义](#31-服务定义)
    - [3.2 消息定义](#32-消息定义)
  - [4. 代码生成](#4-代码生成)
    - [4.1 生成脚本](#41-生成脚本)
    - [4.2 Makefile 命令](#42-makefile-命令)
  - [5. Handler 实现](#5-handler-实现)
    - [5.1 用户服务 Handler](#51-用户服务-handler)
  - [6. 拦截器](#6-拦截器)
    - [6.1 日志拦截器](#61-日志拦截器)
    - [6.2 追踪拦截器](#62-追踪拦截器)
  - [7. 使用示例](#7-使用示例)
  - [📚 相关资源](#-相关资源)

---

## 1. 概述

gRPC 是一个高性能、开源的 RPC 框架，使用 Protocol Buffers 作为接口定义语言。

### 特性

- ✅ **高性能**: 使用 HTTP/2 和 Protocol Buffers
- ✅ **类型安全**: 通过 .proto 文件定义接口
- ✅ **流式支持**: 支持客户端流、服务端流、双向流
- ✅ **跨语言**: 支持多种编程语言

---

## 2. 快速开始

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

---

## 3. Proto 文件定义

### 3.1 服务定义

**文件**: `internal/interfaces/grpc/proto/user.proto`

```protobuf
syntax = "proto3";

package user.v1;

service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}
```

### 3.2 消息定义

```protobuf
message User {
  string id = 1;
  string name = 2;
  string email = 3;
}

message GetUserRequest {
  string id = 1;
}

message GetUserResponse {
  User user = 1;
}
```

---

## 4. 代码生成

### 4.1 生成脚本

**文件**: `scripts/grpc/generate.sh`

```bash
protoc \
  --go_out=internal/interfaces/grpc/proto \
  --go_opt=paths=source_relative \
  --go-grpc_out=internal/interfaces/grpc/proto \
  --go-grpc_opt=paths=source_relative \
  -I=internal/interfaces/grpc/proto \
  internal/interfaces/grpc/proto/*.proto
```

### 4.2 Makefile 命令

```bash
make generate-grpc
```

---

## 5. Handler 实现

### 5.1 用户服务 Handler

**文件**: `internal/interfaces/grpc/handlers/user_handler.go`

```go
type UserHandler struct {
    userpb.UnimplementedUserServiceServer
    service *user.Service
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
    u, err := h.service.GetByID(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.NotFound, err.Error())
    }
    return &userpb.GetUserResponse{User: toProtoUser(u)}, nil
}
```

---

## 6. 拦截器

### 6.1 日志拦截器

**文件**: `internal/interfaces/grpc/interceptors/logging.go`

记录所有 gRPC 请求和响应。

### 6.2 追踪拦截器

**文件**: `internal/interfaces/grpc/interceptors/tracing.go`

集成 OpenTelemetry 进行分布式追踪。

---

## 7. 使用示例

完整示例请参考：

- `examples/grpc/server.go`
- `examples/grpc/client.go`

---

## 📚 相关资源

- [gRPC 官方文档](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [代码实现](../internal/interfaces/grpc/)

---

**最后更新**: 2025-01-XX
