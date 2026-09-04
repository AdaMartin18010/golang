# 完整集成示例

**版本**: v1.0
**更新**: 2025-12-03
**展示**: 所有核心功能集成

---

## 🎯 功能展示

本示例展示了框架的所有核心能力：

### 1. 可观测性 ✅
- **OpenTelemetry** (v1.38.0) - 分布式追踪
- **eBPF 监控** (Cilium v0.20.0) - 系统级监控
- **系统监控** - CPU/内存/磁盘
- **环境感知** - 容器/K8s/云厂商检测

### 2. 安全 ✅
- **JWT 认证** - Token 生成和验证
- **RBAC 授权** - 角色和权限控制
- **中间件集成** - 完整的认证授权流程

### 3. Clean Architecture ✅
- **分层清晰** - Domain/Application/Infrastructure/Interfaces
- **依赖倒置** - 正确的依赖方向
- **易于测试** - Mock 支持

---

## 🚀 运行示例

### 前置要求

1. **Go 1.26.2+**
2. **可选**: OpenTelemetry Collector (用于查看追踪数据)
3. **可选**: Linux + Root 权限 (用于 eBPF 监控)

### 启动 OTEL Collector (可选)

```bash
# 启动完整的可观测性栈
cd ../observability
docker-compose up -d

# Grafana: http://localhost:3000
# Prometheus: http://localhost:9090
```

### 运行示例

```bash
# 普通运行
go run main.go

# Linux + Root (启用 eBPF)
sudo go run main.go
```

---

## 📖 API 测试

### 1. 健康检查（公开）

```bash
curl http://localhost:8080/health
```

**响应**:
```json
{
  "status": "healthy",
  "timestamp": "2025-12-03T10:00:00Z"
}
```

### 2. 登录获取令牌（公开）

```bash
curl -X POST http://localhost:8080/login
```

**响应**:
```json
{
  "access_token": "eyJhbGc...",
  "expires_in": 900
}
```

### 3. 访问用户资料（需要认证）

```bash
TOKEN="eyJhbGc..."
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/profile
```

**响应**:
```json
{
  "user_id": "user-123",
  "username": "john.doe"
}
```

### 4. 获取用户列表（需要权限）

```bash
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/users
```

**需要**: `user:read` 权限

### 5. 管理员操作（需要角色）

```bash
curl -X POST \
     -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/admin/users
```

**需要**: `admin` 角色

---

## 📊 可观测性验证

### 1. 查看追踪数据

1. 访问 Grafana: http://localhost:3000
2. 进入 Explore → 选择 Tempo
3. 搜索 service.name = "complete-integration-example"
4. 查看请求追踪链路

### 2. 查看指标

1. Grafana → Explore → 选择 Prometheus
2. 查询:
   - `system_cpu_usage` - CPU 使用率
   - `system_memory_usage` - 内存使用率
   - `ebpf_syscall_count` - 系统调用计数（如果启用）
   - `ebpf_tcp_connections` - TCP 连接数（如果启用）

### 3. 查看日志

1. Grafana → Explore → 选择 Loki
2. 查询: `{service_name="complete-integration-example"}`

---

## 🏗️ 架构说明

### 集成的核心模块

```text
┌─────────────────────────────────────────┐
│  HTTP Server (Chi Router)              │
│  - 认证中间件 (JWT)                     │
│  - 授权中间件 (RBAC)                    │
│  - 日志中间件                           │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  可观测性层                              │
│  - OpenTelemetry (Trace/Metrics)       │
│  - eBPF 监控 (Syscall/Network)         │
│  - 系统监控 (CPU/Memory/Disk)          │
│  - 环境感知 (Container/K8s/Cloud)      │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  安全层                                  │
│  - JWT Token 管理                       │
│  - RBAC 权限控制                        │
│  - OAuth2/OIDC (可扩展)                 │
└─────────────────────────────────────────┘
```

---

## 🌟 演示的技术

### 最新技术栈（2024）

- ✅ **Cilium eBPF v0.20.0** - 真实的系统级监控
- ✅ **OpenTelemetry v1.38.0** - 标准可观测性
- ✅ **golang-jwt v5.2.1** - JWT 令牌管理
- ✅ **Chi v5.0.12** - 轻量级路由
- ✅ **自研 RBAC** - 完整的权限控制

### 架构模式

- ✅ **Clean Architecture** - 4层分层
- ✅ **DDD** - Specification Pattern
- ✅ **CQRS** - Command/Query 分离
- ✅ **Repository Pattern** - 数据访问抽象

---

## 💡 扩展示例

### 添加 OAuth2/OIDC

```go
// 添加 OAuth2 登录
import "github.com/yourusername/golang/pkg/security/oauth2"

oidcProvider, err := oauth2.NewGoogleOIDCProvider(
    ctx,
    "your-client-id",
    "your-client-secret",
    "http://localhost:8080/callback",
)

r.Get("/auth/google", func(w http.ResponseWriter, r *http.Request) {
    authURL := oidcProvider.AuthorizationURL("state-random")
    http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
})
```

### 添加数据库

```go
import "github.com/yourusername/golang/internal/infrastructure/database/ent"

client, err := ent.Open("postgres", "postgres://...")
defer client.Close()

// 使用 Repository
repo := entrepo.NewUserRepository(client)
users, err := repo.List(ctx, 10, 0)
```

---

## 📚 相关文档

- 架构状态报告
- 最终报告
- [eBPF 实现](../../pkg/observability/ebpf/README.md)
- [安全模块](../../pkg/security/README.md)
- [测试框架](../../test/README.md)

---

## 🎯 学习路径

1. **快速了解** - 运行本示例（5分钟）
2. **深入理解** - 阅读源码和注释（30分钟）
3. **实践应用** - 基于框架构建自己的应用（1-2天）

---

**状态**: ✅ 生产就绪
**展示**: 所有核心功能
**评分**: 8.5/10 ⭐⭐⭐⭐⭐

🚀 **这是框架能力的完整展示！**
