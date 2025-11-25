# Go应用开发

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go应用开发](#go应用开发)
  - [📋 目录](#-目录)
    - [核心模块](#核心模块)
  - [🎯 开发路径](#-开发路径)
    - [Web开发 (2-3周)](#web开发-2-3周)
    - [微服务 (3-4周)](#微服务-3-4周)
    - [云原生 (2-3周)](#云原生-2-3周)
  - [🚀 快速开始](#-快速开始)
    - [HTTP服务器](#http服务器)
    - [gRPC服务](#grpc服务)
  - [📖 系统文档](#-系统文档)
  - [🛠️ 常用技术栈](#️-常用技术栈)
    - [Web框架](#web框架)
    - [数据库](#数据库)
    - [微服务](#微服务)
  - [📚 推荐阅读顺序](#-推荐阅读顺序)

---

### 核心模块

1. **[Web开发](./web/README.md)** ⭐⭐⭐⭐⭐
   - HTTP服务器
   - RESTful API
   - Web框架 (Gin/Echo/Fiber)
   - 路由与中间件

2. **[微服务](./microservices/README.md)** ⭐⭐⭐⭐⭐
   - 微服务架构
   - gRPC与Protobuf
   - 服务发现与注册
   - API网关

3. **[数据库](./database/README.md)** ⭐⭐⭐⭐⭐
   - SQL数据库 (MySQL/PostgreSQL)
   - NoSQL数据库 (MongoDB/Redis)
   - ORM框架 (GORM)
   - 数据库设计

4. **[云原生](./cloud-native/README.md)** ⭐⭐⭐⭐⭐
   - Docker容器化
   - Kubernetes部署
   - 服务网格
   - 云原生实践

---

## 🎯 开发路径

### Web开发 (2-3周)

```
HTTP基础 → RESTful API → 数据库集成 → Web框架
```

### 微服务 (3-4周)

```
服务拆分 → gRPC → 服务发现 → API网关 → 监控
```

### 云原生 (2-3周)

```
Docker → Kubernetes → 配置管理 → CI/CD
```

---

## 🚀 快速开始

### HTTP服务器

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    r.Run(":8080")
}
```

### gRPC服务

```go
package main

import (
    "Context"
    "google.golang.org/grpc"
    pb "myapp/proto"
)

type server struct {
    pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx Context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
```

---

## 📖 系统文档

- **[知识图谱](./00-知识图谱.md)**: 开发知识体系全景图
- **[对比矩阵](./00-对比矩阵.md)**: 技术方案对比
- **[概念定义体系](./00-概念定义体系.md)**: 核心概念详解

---

## 🛠️ 常用技术栈

### Web框架

- Gin (⭐47K) - 高性能HTTP框架
- Echo (⭐27K) - 简洁优雅
- Fiber (⭐30K) - Express风格

### 数据库

- GORM (⭐34K) - 功能强大的ORM
- sqlx (⭐15K) - SQL扩展
- go-redis (⭐18K) - Redis客户端

### 微服务

- gRPC - RPC框架
- Consul/etcd - 服务发现
- Prometheus - 监控
- Jaeger - 链路追踪

---

## 📚 推荐阅读顺序
