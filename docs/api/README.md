# API 文档

> **版本**: v1.0  
> **更新日期**: 2026-04-02  
> **状态**: 持续完善中

---

## 📚 API 概览

本项目提供以下 API 接口：

| 协议 | 路径 | 说明 |
|------|------|------|
| HTTP REST | `/api/v1/*` | 主 REST API |
| gRPC | `:50051` | 高性能 RPC 接口 |
| GraphQL | `/graphql` | 灵活查询接口 |
| WebSocket | `/ws` | 实时通信 |

---

## 🔗 接口文档

### HTTP REST API

**基础信息**:
- Base URL: `http://localhost:8080/api/v1`
- Content-Type: `application/json`
- 认证: Bearer Token (JWT)

**主要端点**:

| 方法 | 路径 | 说明 | 状态 |
|------|------|------|------|
| GET | `/health` | 健康检查 | ✅ 可用 |
| POST | `/users` | 创建用户 | ✅ 可用 |
| GET | `/users/{id}` | 获取用户 | ✅ 可用 |
| PUT | `/users/{id}` | 更新用户 | ✅ 可用 |
| DELETE | `/users/{id}` | 删除用户 | ✅ 可用 |
| GET | `/users` | 用户列表 | ✅ 可用 |

**请求示例**:
```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "email": "user@example.com",
    "name": "Test User"
  }'
```

---

### gRPC API

**服务定义**: `api/proto/v1/*.proto`

**主要服务**:

| 服务 | 说明 | 方法 |
|------|------|------|
| UserService | 用户服务 | CreateUser, GetUser, UpdateUser, DeleteUser, ListUsers |

**连接示例**:
```go
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := pb.NewUserServiceClient(conn)
resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
    Email: "user@example.com",
    Name:  "Test User",
})
```

---

## 🔧 生成 API 文档

```bash
# 生成 OpenAPI 文档
make generate-openapi

# 启动 Swagger UI
docker-compose up swagger-ui
```

---

*最后更新: 2026-04-02*
