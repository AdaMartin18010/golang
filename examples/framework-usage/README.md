# 框架使用示例

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [框架使用示例](#框架使用示例)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 运行示例](#2-运行示例)
    - [2.1 启动服务](#21-启动服务)
    - [2.2 测试登录](#22-测试登录)
    - [2.3 获取用户列表](#23-获取用户列表)
    - [2.4 创建用户](#24-创建用户)
  - [3. 功能说明](#3-功能说明)
    - [3.1 认证流程](#31-认证流程)
    - [3.2 权限控制](#32-权限控制)
    - [3.3 事件处理](#33-事件处理)
  - [4. API端点](#4-api端点)

---

## 1. 概述

这是一个完整的框架使用示例，展示了如何使用项目框架的各个组件：

- ✅ 错误处理和响应格式
- ✅ 日志记录
- ✅ JWT认证
- ✅ 请求验证
- ✅ 中间件（追踪、监控、认证）
- ✅ 健康检查
- ✅ 事件总线
- ✅ RBAC权限控制

---

## 2. 运行示例

### 2.1 启动服务

```bash
cd examples/framework-usage
go run main.go
```

### 2.2 测试登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'
```

### 2.3 获取用户列表

```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 2.4 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

---

## 3. 功能说明

### 3.1 认证流程

1. 用户通过 `/api/v1/auth/login` 登录
2. 服务器返回 `access_token` 和 `refresh_token`
3. 客户端在后续请求中使用 `Authorization: Bearer <token>` 头

### 3.2 权限控制

- `admin` 角色：可以读取和写入用户
- `user` 角色：只能读取用户

### 3.3 事件处理

创建用户时会发布 `user.created` 事件，订阅者会收到通知。

---

## 4. API端点

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/health` | 健康检查 | 否 |
| GET | `/metrics` | 性能指标 | 否 |
| POST | `/api/v1/auth/login` | 登录 | 否 |
| GET | `/api/v1/users` | 获取用户列表 | 是 |
| POST | `/api/v1/users` | 创建用户 | 是 |

---

**更新日期**: 2025-11-11
